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
	// DefaultRun is the implicit run used when a request does not name one. It
	// is backed by the configured DatabaseURL so existing integrations keep
	// working unchanged.
	DefaultRun = "default"

	// DefaultWorkspace is kept as a legacy alias for callers and tests that
	// still use workspace terminology.
	DefaultWorkspace = DefaultRun

	// WorkspaceHeader is a legacy run selector. It is echoed with the resolved
	// run so old integrations can confirm the selected partition.
	WorkspaceHeader = "X-Billtap-Workspace"

	// WorkspaceQueryParam is a legacy alternative to WorkspaceHeader for callers
	// that cannot easily set headers. Prefer /runs/<runId> for new code.
	WorkspaceQueryParam = "workspace"

	// RunHeader carries the path-scoped run ID resolved from /runs/<runId>.
	RunHeader = "X-Billtap-Run-Id"

	// RunPrefixHeader carries the browser path prefix for a path-scoped run.
	RunPrefixHeader = "X-Billtap-Run-Prefix"

	maxRunIDLength = 63
)

// runIDPattern keeps names safe to use as SQLite filenames: an
// alphanumeric lead character followed by alphanumerics, dot, dash, or
// underscore. The leading-character rule rejects "" and dotted paths ("..").
var runIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// apiHandlerBuilder constructs the Stripe-like API handler for one store.
type apiHandlerBuilder func(storage.Store) (http.Handler, error)

type runContextKey struct{}

// runManager owns one isolated billing store and API handler per run ID. The
// default run reuses the externally-owned store; named runs open their own
// SQLite database lazily on first use.
type runManager struct {
	cfg   config.Config
	build apiHandlerBuilder

	mu           sync.Mutex
	handlers     map[string]http.Handler  // name -> API handler
	stores       map[string]storage.Store // name -> store (lazily opened only)
	defaultStore storage.Store

	// apiEnabled is false when the default store cannot back the API (for
	// example a non-billing store). It preserves the previous behaviour of
	// not mounting /v1/ at all in that case.
	apiEnabled bool
}

func newRunManager(cfg config.Config, defaultStore storage.Store, build apiHandlerBuilder) *runManager {
	m := &runManager{
		cfg:          cfg,
		build:        build,
		handlers:     make(map[string]http.Handler),
		stores:       make(map[string]storage.Store),
		defaultStore: defaultStore,
	}
	if defaultStore == nil {
		return m
	}
	handler, err := build(defaultStore)
	if err != nil {
		return m
	}
	m.handlers[DefaultRun] = handler
	m.apiEnabled = true
	return m
}

type runSummary struct {
	Name      string
	IsDefault bool
	Open      bool
	Storage   string
	Summary   map[string]int
	Error     string
}

// apiHandler returns the dispatcher mounted on /v1/ and /api/. It resolves the
// run for each request, lazily provisioning isolated storage as needed.
func (m *runManager) apiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := resolveRunSelector(r)
		if err != nil {
			writeRunError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		apiHandler, err := m.get(r.Context(), name)
		if err != nil {
			writeRunError(w, r, http.StatusInternalServerError,
				fmt.Sprintf("could not open run %q: %v", name, err))
			return
		}
		w.Header().Set(RunHeader, name)
		w.Header().Set(WorkspaceHeader, name)
		apiHandler.ServeHTTP(w, requestForRunAPI(r, name))
	})
}

// requestForRunAPI returns the request seen by the Stripe-compatible API after
// the server has resolved isolation. It overwrites any client-supplied run
// header with the canonical run and strips the legacy workspace query so strict
// parameter validation never sees it.
func requestForRunAPI(r *http.Request, name string) *http.Request {
	clone := r.Clone(r.Context())
	clone.Header = r.Header.Clone()
	clone.Header.Set(RunHeader, name)

	cloned := *r.URL
	if query := r.URL.Query(); query.Has(WorkspaceQueryParam) {
		query.Del(WorkspaceQueryParam)
		cloned.RawQuery = query.Encode()
		cloned.RawPath = ""
	}
	clone.URL = &cloned
	return clone
}

// get returns the API handler for name, opening its store on first use.
func (m *runManager) get(ctx context.Context, name string) (http.Handler, error) {
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
	store, err := storage.OpenSQLite(context.WithoutCancel(ctx), runDSN(m.cfg.DatabaseURL, name))
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

// list reports the known runs: the default, any opened this session, and any
// whose database file already exists on disk.
func (m *runManager) list() []string {
	set := map[string]bool{DefaultRun: true}

	m.mu.Lock()
	for name := range m.handlers {
		set[name] = true
	}
	m.mu.Unlock()

	if dir := runStoreDir(m.cfg.DatabaseURL); dir != "" {
		ext := runDBExt(m.cfg.DatabaseURL)
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

func (m *runManager) summaries(ctx context.Context) []runSummary {
	names := m.list()
	out := make([]runSummary, 0, len(names))
	for _, name := range names {
		summary := runSummary{
			Name:      name,
			IsDefault: name == DefaultRun,
			Open:      m.isOpen(name),
			Storage:   runDSN(m.cfg.DatabaseURL, name),
			Summary:   map[string]int{},
		}
		store, closeStore, err := m.storeForSummary(ctx, name)
		if err != nil {
			summary.Error = err.Error()
			out = append(out, summary)
			continue
		}
		counts, err := storage.SQLiteTableCounts(ctx, store)
		if closeStore != nil {
			closeStore()
		}
		if err != nil {
			summary.Error = err.Error()
		} else {
			summary.Summary = counts
		}
		out = append(out, summary)
	}
	return out
}

func (m *runManager) isOpen(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if name == DefaultRun {
		return m.defaultStore != nil
	}
	_, ok := m.handlers[name]
	return ok
}

func (m *runManager) storeForSummary(ctx context.Context, name string) (storage.Store, func(), error) {
	m.mu.Lock()
	if name == DefaultRun {
		store := m.defaultStore
		m.mu.Unlock()
		if store == nil {
			return nil, nil, fmt.Errorf("default run storage is not open")
		}
		return store, nil, nil
	}
	if store, ok := m.stores[name]; ok {
		m.mu.Unlock()
		return store, nil, nil
	}
	m.mu.Unlock()

	if !runDBExists(m.cfg.DatabaseURL, name) {
		return nil, nil, fmt.Errorf("run storage does not exist")
	}
	store, err := storage.OpenSQLite(context.WithoutCancel(ctx), runDSN(m.cfg.DatabaseURL, name))
	if err != nil {
		return nil, nil, err
	}
	return store, func() { _ = store.Close() }, nil
}

func (m *runManager) delete(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("run id is required")
	}
	if name == DefaultRun {
		m.mu.Lock()
		store := m.defaultStore
		m.mu.Unlock()
		if store == nil {
			return fmt.Errorf("default run storage is not open")
		}
		return storage.ResetSQLiteData(ctx, store)
	}

	var store storage.Store
	m.mu.Lock()
	if existing, ok := m.stores[name]; ok {
		store = existing
		delete(m.stores, name)
	}
	delete(m.handlers, name)
	m.mu.Unlock()
	if store != nil {
		if err := store.Close(); err != nil {
			return err
		}
	}
	if isMemoryDSN(m.cfg.DatabaseURL) {
		return nil
	}
	path := runDBPath(m.cfg.DatabaseURL, name)
	for _, candidate := range []string{path, path + "-wal", path + "-shm"} {
		if candidate == "" {
			continue
		}
		if err := os.Remove(candidate); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// Close releases every lazily-opened run store. The default store is owned by
// the caller and is left untouched.
func (m *runManager) Close() error {
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

// resolveRunSelector returns the run chosen by path scope or, for unprefixed
// requests, by the legacy workspace selector.
func resolveRunSelector(r *http.Request) (string, error) {
	if scoped, ok := r.Context().Value(runContextKey{}).(string); ok && scoped != "" {
		return scoped, nil
	}
	raw := strings.TrimSpace(r.Header.Get(WorkspaceHeader))
	if raw == "" {
		raw = strings.TrimSpace(r.URL.Query().Get(WorkspaceQueryParam))
	}
	return normalizeRunID(raw)
}

func normalizeRunID(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return DefaultRun, nil
	}

	// Filenames on macOS/Windows are case-insensitive; normalise so "Foo"
	// and "foo" cannot resolve to two handlers over one file.
	name := strings.ToLower(raw)
	if name == DefaultRun {
		return DefaultRun, nil
	}
	if len(name) > maxRunIDLength {
		return "", fmt.Errorf("run id must be at most %d characters", maxRunIDLength)
	}
	if !runIDPattern.MatchString(name) {
		return "", fmt.Errorf("run id %q is invalid: use letters, digits, '.', '-', '_' and a leading alphanumeric", raw)
	}
	return name, nil
}

// NormalizeRunID validates a user-supplied run ID using the same rules as
// /runs/<runId> routing and legacy workspace aliases.
func NormalizeRunID(raw string) (string, error) {
	return normalizeRunID(raw)
}

// runDSN derives the SQLite DSN for a named run from the base DSN. The default
// run returns the base DSN unchanged.
func runDSN(baseDSN, name string) string {
	if name == "" || name == DefaultRun {
		return baseDSN
	}
	if isMemoryDSN(baseDSN) {
		// Each in-memory run needs a distinct shared-cache name so it
		// stays isolated yet survives across pooled connections.
		return fmt.Sprintf("file:billtap_run_%s?mode=memory&cache=shared", name)
	}

	_, query := splitDSN(baseDSN)
	wsPath := runDBPath(baseDSN, name)
	if query == "" {
		return wsPath
	}
	return "file:" + wsPath + query
}

// RunDSN returns the SQLite DSN used for a path-scoped run ID.
func RunDSN(baseDSN, runID string) string {
	return runDSN(baseDSN, runID)
}

func runDBPath(baseDSN string, name string) string {
	path, _ := splitDSN(baseDSN)
	ext := filepath.Ext(path)
	if ext == "" {
		ext = ".db"
	}
	// Keep the existing on-disk directory so previously-created isolated stores
	// remain visible after the run terminology cleanup.
	return filepath.Join(filepath.Dir(path), "workspaces", name+ext)
}

func runDBExists(baseDSN string, name string) bool {
	if name == "" || name == DefaultRun {
		return true
	}
	if isMemoryDSN(baseDSN) {
		return false
	}
	if _, err := os.Stat(runDBPath(baseDSN, name)); err == nil {
		return true
	}
	return false
}

// runStoreDir returns the directory that holds named run databases, or "" when
// the base DSN is in-memory.
func runStoreDir(baseDSN string) string {
	if isMemoryDSN(baseDSN) {
		return ""
	}
	path, _ := splitDSN(baseDSN)
	return filepath.Join(filepath.Dir(path), "workspaces")
}

func runDBExt(baseDSN string) string {
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

func writeRunError(w http.ResponseWriter, r *http.Request, status int, message string) {
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
