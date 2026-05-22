package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/hckim/billtap/internal/config"
	"github.com/hckim/billtap/internal/storage"
)

const (
	// DefaultWorkspace is the implicit workspace used when a request does not
	// name one. It is backed by the configured DatabaseURL so existing
	// integrations keep working unchanged.
	DefaultWorkspace = "default"

	// WorkspaceHeader carries the target workspace name on a request and is
	// echoed back on the response so callers can confirm the resolved value.
	WorkspaceHeader = "X-Billtap-Workspace"

	// WorkspaceQueryParam is an alternative to WorkspaceHeader for callers
	// that cannot easily set headers.
	WorkspaceQueryParam = "workspace"

	maxWorkspaceNameLength = 63
)

// workspaceNamePattern keeps names safe to use as SQLite filenames: an
// alphanumeric lead character followed by alphanumerics, dot, dash, or
// underscore. The leading-character rule rejects "" and dotted paths ("..").
var workspaceNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// apiHandlerBuilder constructs the Stripe-like API handler for one store.
type apiHandlerBuilder func(storage.Store) (http.Handler, error)

// workspaceManager owns one isolated billing store (and API handler) per
// workspace name. The default workspace reuses the externally-owned store;
// named workspaces open their own SQLite database lazily on first use.
type workspaceManager struct {
	cfg   config.Config
	build apiHandlerBuilder

	mu       sync.Mutex
	handlers map[string]http.Handler  // name -> API handler
	stores   map[string]storage.Store // name -> store (lazily opened only)

	// apiEnabled is false when the default store cannot back the API (for
	// example a non-billing store). It preserves the previous behaviour of
	// not mounting /v1/ at all in that case.
	apiEnabled bool
}

func newWorkspaceManager(cfg config.Config, defaultStore storage.Store, build apiHandlerBuilder) *workspaceManager {
	m := &workspaceManager{
		cfg:      cfg,
		build:    build,
		handlers: make(map[string]http.Handler),
		stores:   make(map[string]storage.Store),
	}
	if defaultStore == nil {
		return m
	}
	handler, err := build(defaultStore)
	if err != nil {
		return m
	}
	m.handlers[DefaultWorkspace] = handler
	m.apiEnabled = true
	return m
}

// handler returns the dispatcher mounted on /v1/ and /api/. It resolves the
// workspace for each request, lazily provisioning isolated storage as needed.
func (m *workspaceManager) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := resolveWorkspace(r)
		if err != nil {
			writeWorkspaceError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		apiHandler, err := m.get(r.Context(), name)
		if err != nil {
			writeWorkspaceError(w, r, http.StatusInternalServerError,
				fmt.Sprintf("could not open workspace %q: %v", name, err))
			return
		}
		w.Header().Set(WorkspaceHeader, name)
		// The workspace selector is consumed here; strip it so the strict
		// API parameter validation downstream never sees it.
		apiHandler.ServeHTTP(w, stripWorkspaceQuery(r))
	})
}

// stripWorkspaceQuery returns a request with the workspace query parameter
// removed, leaving the original untouched when it carries no such parameter.
func stripWorkspaceQuery(r *http.Request) *http.Request {
	query := r.URL.Query()
	if _, ok := query[WorkspaceQueryParam]; !ok {
		return r
	}
	query.Del(WorkspaceQueryParam)
	clone := r.Clone(r.Context())
	cloned := *r.URL
	cloned.RawQuery = query.Encode()
	cloned.RawPath = ""
	clone.URL = &cloned
	return clone
}

// get returns the API handler for name, opening its store on first use.
func (m *workspaceManager) get(ctx context.Context, name string) (http.Handler, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if handler, ok := m.handlers[name]; ok {
		return handler, nil
	}
	if !m.apiEnabled {
		return nil, fmt.Errorf("billing API is not available for this storage backend")
	}

	// Decouple the store lifetime from the triggering request: the store is
	// reused for every later request, so a cancelled first request must not
	// tear it down.
	store, err := storage.OpenSQLite(context.WithoutCancel(ctx), workspaceDSN(m.cfg.DatabaseURL, name))
	if err != nil {
		return nil, err
	}
	handler, err := m.build(store)
	if err != nil {
		_ = store.Close()
		return nil, err
	}
	m.handlers[name] = handler
	m.stores[name] = store
	return handler, nil
}

// list reports the known workspaces: the default, any opened this session,
// and any whose database file already exists on disk.
func (m *workspaceManager) list() []string {
	set := map[string]bool{DefaultWorkspace: true}

	m.mu.Lock()
	for name := range m.handlers {
		set[name] = true
	}
	m.mu.Unlock()

	if dir := workspacesDir(m.cfg.DatabaseURL); dir != "" {
		ext := workspaceDBExt(m.cfg.DatabaseURL)
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ext) {
					continue
				}
				set[strings.TrimSuffix(entry.Name(), ext)] = true
			}
		}
	}

	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Close releases every lazily-opened workspace store. The default store is
// owned by the caller and is left untouched.
func (m *workspaceManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error
	for name, store := range m.stores {
		if err := store.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(m.stores, name)
		delete(m.handlers, name)
	}
	return firstErr
}

// resolveWorkspace extracts and validates the workspace name from a request,
// falling back to DefaultWorkspace when none is supplied.
func resolveWorkspace(r *http.Request) (string, error) {
	raw := strings.TrimSpace(r.Header.Get(WorkspaceHeader))
	if raw == "" {
		raw = strings.TrimSpace(r.URL.Query().Get(WorkspaceQueryParam))
	}
	if raw == "" {
		return DefaultWorkspace, nil
	}

	// Filenames on macOS/Windows are case-insensitive; normalise so "Foo"
	// and "foo" cannot resolve to two handlers over one file.
	name := strings.ToLower(raw)
	if name == DefaultWorkspace {
		return DefaultWorkspace, nil
	}
	if len(name) > maxWorkspaceNameLength {
		return "", fmt.Errorf("workspace name must be at most %d characters", maxWorkspaceNameLength)
	}
	if !workspaceNamePattern.MatchString(name) {
		return "", fmt.Errorf("workspace name %q is invalid: use letters, digits, '.', '-', '_' and a leading alphanumeric", raw)
	}
	return name, nil
}

// workspaceDSN derives the SQLite DSN for a named workspace from the base
// (default) DSN. The default workspace returns the base DSN unchanged.
func workspaceDSN(baseDSN, name string) string {
	if name == "" || name == DefaultWorkspace {
		return baseDSN
	}
	if isMemoryDSN(baseDSN) {
		// Each in-memory workspace needs a distinct shared-cache name so it
		// stays isolated yet survives across pooled connections.
		return fmt.Sprintf("file:billtap_ws_%s?mode=memory&cache=shared", name)
	}

	path, query := splitDSN(baseDSN)
	ext := filepath.Ext(path)
	if ext == "" {
		ext = ".db"
	}
	wsPath := filepath.Join(filepath.Dir(path), "workspaces", name+ext)
	if query == "" {
		return wsPath
	}
	return "file:" + wsPath + query
}

// workspacesDir returns the directory that holds named workspace databases,
// or "" when the base DSN is in-memory.
func workspacesDir(baseDSN string) string {
	if isMemoryDSN(baseDSN) {
		return ""
	}
	path, _ := splitDSN(baseDSN)
	return filepath.Join(filepath.Dir(path), "workspaces")
}

func workspaceDBExt(baseDSN string) string {
	path, _ := splitDSN(baseDSN)
	if ext := filepath.Ext(path); ext != "" {
		return ext
	}
	return ".db"
}

// splitDSN separates a SQLite DSN into its filesystem path and trailing
// query/fragment, dropping any leading "file:" scheme.
func splitDSN(dsn string) (path string, query string) {
	path = strings.TrimPrefix(dsn, "file:")
	if idx := strings.IndexAny(path, "?#"); idx >= 0 {
		return path[:idx], path[idx:]
	}
	return path, ""
}

func isMemoryDSN(dsn string) bool {
	return dsn == ":memory:" ||
		strings.HasPrefix(dsn, "file::memory:") ||
		strings.Contains(dsn, "mode=memory")
}

func writeWorkspaceError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if r.Method == http.MethodHead {
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"type":    "invalid_request_error",
			"message": message,
		},
	})
}
