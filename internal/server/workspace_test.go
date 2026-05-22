package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/hckim/billtap/internal/config"
	"github.com/hckim/billtap/internal/storage"
)

// newWorkspaceServer builds a SQLite-backed server whose configured
// DatabaseURL matches the default store, so named workspaces resolve to
// sibling files under <dir>/workspaces.
func newWorkspaceServer(t *testing.T) (*Server, string) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "billtap.db")
	store, err := storage.OpenSQLite(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	cfg := config.Config{
		Addr:        ":0",
		DatabaseURL: dbPath,
		StaticDir:   "web/dist",
		Environment: "test",
	}
	srv := New(Options{Config: cfg, Store: store})
	t.Cleanup(func() {
		_ = srv.Close()
		_ = store.Close()
	})
	return srv, dir
}

func countCustomers(t *testing.T, handler http.Handler, workspace string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	if workspace != "" {
		req.Header.Set(WorkspaceHeader, workspace)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list customers (workspace=%q) status = %d body = %s", workspace, rec.Code, rec.Body.String())
	}
	var out struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode customer list: %v body=%s", err, rec.Body.String())
	}
	return len(out.Data)
}

func TestWorkspacesIsolateBillingData(t *testing.T) {
	srv, _ := newWorkspaceServer(t)

	// Two customers in the default workspace, one in a named workspace.
	postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/v1/customers", map[string]string{"email": "default-1@example.test"})
	postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/v1/customers", map[string]string{"email": "default-2@example.test"})
	postFormWithHeaders[struct {
		ID string `json:"id"`
	}](t, srv, "/v1/customers", map[string]string{"email": "alt@example.test"},
		map[string]string{WorkspaceHeader: "test-a"})

	if got := countCustomers(t, srv, ""); got != 2 {
		t.Fatalf("default workspace customer count = %d, want 2", got)
	}
	if got := countCustomers(t, srv, "test-a"); got != 1 {
		t.Fatalf("test-a workspace customer count = %d, want 1", got)
	}
	if got := countCustomers(t, srv, "default"); got != 2 {
		t.Fatalf("explicit default workspace customer count = %d, want 2", got)
	}
	if got := countCustomers(t, srv, "test-b"); got != 0 {
		t.Fatalf("fresh workspace customer count = %d, want 0", got)
	}
}

func TestWorkspaceResolvedFromQueryParam(t *testing.T) {
	srv, _ := newWorkspaceServer(t)

	postFormWithHeaders[struct {
		ID string `json:"id"`
	}](t, srv, "/v1/customers?workspace=via-query", map[string]string{"email": "q@example.test"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/customers?workspace=via-query", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get(WorkspaceHeader); got != "via-query" {
		t.Fatalf("response %s = %q, want %q", WorkspaceHeader, got, "via-query")
	}
	if got := countCustomers(t, srv, ""); got != 0 {
		t.Fatalf("default workspace should stay empty, got %d", got)
	}
}

func TestWorkspaceHeaderEchoedAndInvalidRejected(t *testing.T) {
	srv, _ := newWorkspaceServer(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	req.Header.Set(WorkspaceHeader, "Mixed-Case")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get(WorkspaceHeader); got != "mixed-case" {
		t.Fatalf("resolved workspace = %q, want lowercased %q", got, "mixed-case")
	}

	for _, bad := range []string{"bad/name", "../escape", ".hidden", "with space"} {
		req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
		req.Header.Set(WorkspaceHeader, bad)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("workspace %q status = %d, want 400", bad, rec.Code)
		}
	}
}

func TestWorkspacesListingEndpoint(t *testing.T) {
	srv, _ := newWorkspaceServer(t)

	postFormWithHeaders[struct {
		ID string `json:"id"`
	}](t, srv, "/v1/customers", map[string]string{"email": "x@example.test"},
		map[string]string{WorkspaceHeader: "scenario-1"})

	req := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Data []struct {
			Name      string `json:"name"`
			IsDefault bool   `json:"is_default"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode workspace list: %v body=%s", err, rec.Body.String())
	}
	seen := make(map[string]bool)
	for _, ws := range out.Data {
		seen[ws.Name] = true
	}
	if !seen[DefaultWorkspace] || !seen["scenario-1"] {
		t.Fatalf("workspace list = %#v, want default and scenario-1", out.Data)
	}
}

func TestWorkspaceDSN(t *testing.T) {
	cases := []struct {
		base string
		name string
		want string
	}{
		{".billtap/billtap.db", "default", ".billtap/billtap.db"},
		{".billtap/billtap.db", "test-a", filepath.Join(".billtap", "workspaces", "test-a.db")},
		{"/data/billtap.db", "ci", filepath.Join("/data", "workspaces", "ci.db")},
		{":memory:", "iso", "file:billtap_ws_iso?mode=memory&cache=shared"},
	}
	for _, tc := range cases {
		if got := workspaceDSN(tc.base, tc.name); got != tc.want {
			t.Fatalf("workspaceDSN(%q, %q) = %q, want %q", tc.base, tc.name, got, tc.want)
		}
	}
}
