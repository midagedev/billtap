package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hckim/billtap/internal/config"
	"github.com/hckim/billtap/internal/storage"
)

// newRunServer builds a SQLite-backed server whose configured DatabaseURL
// matches the default store, so named runs resolve to sibling SQLite files.
func newRunServer(t *testing.T) (*Server, string) {
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

func countCustomers(t *testing.T, handler http.Handler, legacyWorkspace string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	if legacyWorkspace != "" {
		req.Header.Set(WorkspaceHeader, legacyWorkspace)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list customers (legacy workspace=%q) status = %d body = %s", legacyWorkspace, rec.Code, rec.Body.String())
	}
	var out struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode customer list: %v body=%s", err, rec.Body.String())
	}
	return len(out.Data)
}

func countCustomersPath(t *testing.T, handler http.Handler, path string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list customers %s status = %d body = %s", path, rec.Code, rec.Body.String())
	}
	var out struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode customer list: %v body=%s", err, rec.Body.String())
	}
	return len(out.Data)
}

func TestLegacyWorkspaceSelectorIsolatesBillingData(t *testing.T) {
	srv, _ := newRunServer(t)

	// Two customers in the default run, one through the legacy workspace alias.
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
		t.Fatalf("default run customer count = %d, want 2", got)
	}
	if got := countCustomers(t, srv, "test-a"); got != 1 {
		t.Fatalf("test-a run customer count = %d, want 1", got)
	}
	if got := countCustomers(t, srv, "default"); got != 2 {
		t.Fatalf("explicit default run customer count = %d, want 2", got)
	}
	if got := countCustomers(t, srv, "test-b"); got != 0 {
		t.Fatalf("fresh run customer count = %d, want 0", got)
	}
}

func TestLegacyWorkspaceResolvedFromQueryParam(t *testing.T) {
	srv, _ := newRunServer(t)

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
		t.Fatalf("default run should stay empty, got %d", got)
	}
}

func TestLegacyWorkspaceHeaderEchoedAndInvalidRejected(t *testing.T) {
	srv, _ := newRunServer(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	req.Header.Set(WorkspaceHeader, "Mixed-Case")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get(WorkspaceHeader); got != "mixed-case" {
		t.Fatalf("resolved run = %q, want lowercased %q", got, "mixed-case")
	}

	for _, bad := range []string{"bad/name", "../escape", ".hidden", "with space"} {
		req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
		req.Header.Set(WorkspaceHeader, bad)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("run %q status = %d, want 400", bad, rec.Code)
		}
	}
}

func TestLegacyWorkspacesListingEndpoint(t *testing.T) {
	srv, _ := newRunServer(t)

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
	if !seen[DefaultRun] || !seen["scenario-1"] {
		t.Fatalf("workspace list = %#v, want default and scenario-1", out.Data)
	}
}

func TestRunDSN(t *testing.T) {
	cases := []struct {
		base string
		name string
		want string
	}{
		{".billtap/billtap.db", "default", ".billtap/billtap.db"},
		{".billtap/billtap.db", "test-a", filepath.Join(".billtap", "workspaces", "test-a.db")},
		{"/data/billtap.db", "ci", filepath.Join("/data", "workspaces", "ci.db")},
		{":memory:", "iso", "file:billtap_run_iso?mode=memory&cache=shared"},
	}
	for _, tc := range cases {
		if got := runDSN(tc.base, tc.name); got != tc.want {
			t.Fatalf("runDSN(%q, %q) = %q, want %q", tc.base, tc.name, got, tc.want)
		}
	}
}

func TestRunPathPrefixIsolatesBillingData(t *testing.T) {
	srv, _ := newRunServer(t)

	postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/run-a/v1/customers", map[string]string{"id": "cus_shared", "email": "a@example.test"})
	postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/run-b/v1/customers", map[string]string{"id": "cus_shared", "email": "b@example.test"})
	postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/v1/customers", map[string]string{"id": "cus_shared", "email": "default@example.test"})

	if got := countCustomersPath(t, srv, "/runs/run-a/v1/customers"); got != 1 {
		t.Fatalf("run-a customer count = %d, want 1", got)
	}
	if got := countCustomersPath(t, srv, "/runs/run-b/v1/customers"); got != 1 {
		t.Fatalf("run-b customer count = %d, want 1", got)
	}
	if got := countCustomers(t, srv, ""); got != 1 {
		t.Fatalf("default customer count = %d, want 1", got)
	}

	req := httptest.NewRequest(http.MethodGet, "/runs/run-a/v1/customers", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if got := rec.Header().Get(RunHeader); got != "run-a" {
		t.Fatalf("%s = %q, want run-a", RunHeader, got)
	}
}

func TestRunPathScopeWinsOverLegacyWorkspaceHeader(t *testing.T) {
	srv, _ := newRunServer(t)

	req := httptest.NewRequest(http.MethodPost, "/runs/run-a/v1/customers", strings.NewReader("email=a@example.test"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(WorkspaceHeader, "run-b")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("create scoped customer status = %d body = %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get(RunHeader); got != "run-a" {
		t.Fatalf("%s = %q, want run-a", RunHeader, got)
	}
	if got := countCustomersPath(t, srv, "/runs/run-a/v1/customers"); got != 1 {
		t.Fatalf("run-a customers = %d, want 1", got)
	}
	if got := countCustomersPath(t, srv, "/runs/run-b/v1/customers"); got != 0 {
		t.Fatalf("run-b customers = %d, want 0", got)
	}
}

func TestRunPathPrefixesHostedCheckoutURL(t *testing.T) {
	srv, _ := newRunServer(t)

	customer := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/checkout-run/v1/customers", map[string]string{"email": "buyer@example.test"})
	product := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/checkout-run/v1/products", map[string]string{"name": "Team"})
	price := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/checkout-run/v1/prices", map[string]string{
		"product":             product.ID,
		"currency":            "usd",
		"unit_amount":         "9900",
		"recurring[interval]": "month",
	})
	session := postForm[struct {
		URL string `json:"url"`
	}](t, srv, "/runs/checkout-run/v1/checkout/sessions", map[string]string{
		"customer":                customer.ID,
		"line_items[0][price]":    price.ID,
		"line_items[0][quantity]": "1",
	})
	if want := "http://example.com/runs/checkout-run/checkout/"; len(session.URL) < len(want) || session.URL[:len(want)] != want {
		t.Fatalf("checkout URL = %q, want run prefix %q", session.URL, want)
	}

	sessionID := session.URL[strings.LastIndex(session.URL, "/")+1:]
	req := httptest.NewRequest(http.MethodGet, "/runs/checkout-run/checkout/"+sessionID, nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("run checkout redirect status = %d, want 302", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/runs/checkout-run/app/checkout/?session_id="+sessionID {
		t.Fatalf("run checkout redirect = %q", got)
	}
}

func TestRunWebhookEndpointsOnlyReceiveRunEvents(t *testing.T) {
	srv, _ := newRunServer(t)

	endpointA := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/webhook-a/v1/webhook_endpoints", map[string]string{
		"url":            "https://app-a.example.test/webhook",
		"enabled_events": "checkout.session.completed",
	})
	endpointB := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/webhook-b/v1/webhook_endpoints", map[string]string{
		"url":            "https://app-b.example.test/webhook",
		"enabled_events": "checkout.session.completed",
	})

	customer := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/webhook-a/v1/customers", map[string]string{"email": "buyer@example.test"})
	product := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/webhook-a/v1/products", map[string]string{"name": "Team"})
	price := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/webhook-a/v1/prices", map[string]string{
		"product":             product.ID,
		"currency":            "usd",
		"unit_amount":         "9900",
		"recurring[interval]": "month",
	})
	session := postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/webhook-a/v1/checkout/sessions", map[string]string{
		"customer":                customer.ID,
		"line_items[0][price]":    price.ID,
		"line_items[0][quantity]": "1",
	})
	_ = postForm[map[string]any](t, srv, "/runs/webhook-a/v1/checkout/sessions/"+session.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})

	attemptsA := getRunList(t, srv, "/runs/webhook-a/v1/webhook_endpoints/"+endpointA.ID+"/attempts")
	if len(attemptsA.Data) == 0 {
		t.Fatalf("run-a endpoint attempts = 0, want checkout delivery attempts")
	}
	attemptsB := getRunList(t, srv, "/runs/webhook-b/v1/webhook_endpoints/"+endpointB.ID+"/attempts")
	if len(attemptsB.Data) != 0 {
		t.Fatalf("run-b endpoint attempts = %d, want 0", len(attemptsB.Data))
	}
}

func TestRunAdminAndCleanup(t *testing.T) {
	srv, _ := newRunServer(t)
	postForm[struct {
		ID string `json:"id"`
	}](t, srv, "/runs/cleanup-run/v1/customers", map[string]string{"email": "cleanup@example.test"})

	before := getRunSummaries(t, srv)
	if got := before["cleanup-run"]["customers"]; got != 1 {
		t.Fatalf("cleanup-run customers before cleanup = %d, want 1; summaries=%#v", got, before)
	}

	req := httptest.NewRequest(http.MethodDelete, "/runs/cleanup-run", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("cleanup status = %d body = %s", rec.Code, rec.Body.String())
	}
	if got := countCustomersPath(t, srv, "/runs/cleanup-run/v1/customers"); got != 0 {
		t.Fatalf("cleanup-run customers after cleanup = %d, want 0", got)
	}
}

func getRunList(t *testing.T, handler http.Handler, path string) struct {
	Data []json.RawMessage `json:"data"`
} {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s status = %d body = %s", path, rec.Code, rec.Body.String())
	}
	var out struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode %s: %v body=%s", path, err, rec.Body.String())
	}
	return out
}

func getRunSummaries(t *testing.T, handler http.Handler) map[string]map[string]int {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/admin/runs", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin runs status = %d body = %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Data []struct {
			RunID   string         `json:"runId"`
			Summary map[string]int `json:"summary"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode admin runs: %v body=%s", err, rec.Body.String())
	}
	summaries := map[string]map[string]int{}
	for _, item := range out.Data {
		summaries[item.RunID] = item.Summary
	}
	return summaries
}
