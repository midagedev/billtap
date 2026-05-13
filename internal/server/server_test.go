package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hckim/billtap/internal/config"
	"github.com/hckim/billtap/internal/storage"
)

func TestHealthEndpoint(t *testing.T) {
	handler := New(Options{
		Config: config.Config{Addr: ":0", DatabaseURL: ":memory:", StaticDir: "web/dist", Environment: "test"},
		Store:  storage.NewMemoryStore(),
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body = %#v, want healthy process", body)
	}
	if _, ok := body["storage"]; ok {
		t.Fatalf("body = %#v, health endpoint should not include storage readiness", body)
	}
}

func TestReadyEndpointChecksStorage(t *testing.T) {
	handler := New(Options{
		Config: config.Config{Addr: ":0", DatabaseURL: ":memory:", StaticDir: "web/dist", Environment: "test"},
		Store:  storage.NewMemoryStore(),
	})

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "ok" || body["storage"] != "ok" {
		t.Fatalf("body = %#v, want healthy storage readiness", body)
	}
}

func TestReactAssetPathStub(t *testing.T) {
	handler := New(Options{
		Config: config.Config{Addr: ":0", DatabaseURL: ":memory:", StaticDir: "web/dist", Environment: "test"},
		Store:  storage.NewMemoryStore(),
	})

	for _, path := range []string{"/app/", "/assets/app.js"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", path, rec.Code, http.StatusOK)
		}
	}
}

func TestHostedRoutesWithoutTrailingSlashRedirectToApps(t *testing.T) {
	handler := New(Options{
		Config: config.Config{Addr: ":0", DatabaseURL: ":memory:", StaticDir: "web/dist", Environment: "test"},
		Store:  storage.NewMemoryStore(),
	})

	for _, tt := range []struct {
		path string
		want string
	}{
		{path: "/checkout?session_id=cs_test_123", want: "/app/checkout/?session_id=cs_test_123"},
		{path: "/portal?customer_id=cus_test_123", want: "/app/portal/?customer_id=cus_test_123"},
	} {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusFound {
			t.Fatalf("%s status = %d, want %d", tt.path, rec.Code, http.StatusFound)
		}
		if got := rec.Header().Get("Location"); got != tt.want {
			t.Fatalf("%s Location = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestPublicBasePathPrefixesBrowserRoutes(t *testing.T) {
	handler := New(Options{
		Config: config.Config{Addr: ":0", DatabaseURL: ":memory:", StaticDir: "web/dist", Environment: "test", PublicBasePath: "/billtap"},
		Store:  storage.NewMemoryStore(),
	})

	for _, tt := range []struct {
		path string
		want string
	}{
		{path: "/billtap", want: "/billtap/app/dashboard/"},
		{path: "/billtap/checkout?session_id=cs_test_123", want: "/billtap/app/checkout/?session_id=cs_test_123"},
		{path: "/checkout?session_id=cs_test_456", want: "/billtap/app/checkout/?session_id=cs_test_456"},
	} {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusFound {
			t.Fatalf("%s status = %d, want %d", tt.path, rec.Code, http.StatusFound)
		}
		if got := rec.Header().Get("Location"); got != tt.want {
			t.Fatalf("%s Location = %q, want %q", tt.path, got, tt.want)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/billtap/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("prefixed health status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestForwardedPrefixDrivesRedirects(t *testing.T) {
	handler := New(Options{
		Config: config.Config{Addr: ":0", DatabaseURL: ":memory:", StaticDir: "web/dist", Environment: "test"},
		Store:  storage.NewMemoryStore(),
	})

	req := httptest.NewRequest(http.MethodGet, "/portal?customer_id=cus_test_123", nil)
	req.Header.Set("X-Forwarded-Prefix", "/billtap")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
	}
	if got := rec.Header().Get("Location"); got != "/billtap/app/portal/?customer_id=cus_test_123" {
		t.Fatalf("Location = %q, want prefixed portal URL", got)
	}
}

func TestBuiltReactAppServing(t *testing.T) {
	staticDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(staticDir, "dashboard"), 0o755); err != nil {
		t.Fatalf("create dashboard dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(staticDir, "assets"), 0o755); err != nil {
		t.Fatalf("create assets dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(staticDir, "dashboard", "index.html"), []byte("dashboard-build"), 0o644); err != nil {
		t.Fatalf("write dashboard index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(staticDir, "assets", "app.js"), []byte("console.info('built')"), 0o644); err != nil {
		t.Fatalf("write built asset: %v", err)
	}

	handler := New(Options{
		Config: config.Config{Addr: ":0", DatabaseURL: ":memory:", StaticDir: staticDir, Environment: "test"},
		Store:  storage.NewMemoryStore(),
	})

	for _, path := range []string{"/app/dashboard/", "/app/assets/app.js"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", path, rec.Code, http.StatusOK)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/checkout?session_id=cs_test_123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("built checkout status = %d, want %d", rec.Code, http.StatusFound)
	}
	if got := rec.Header().Get("Location"); got != "/app/checkout/?session_id=cs_test_123" {
		t.Fatalf("built checkout Location = %q, want app redirect", got)
	}
}

func TestPublicBasePathPrefixesAPISessionURLs(t *testing.T) {
	handler := newSQLiteBackedServer(t, config.Config{
		Addr:           ":0",
		DatabaseURL:    ":memory:",
		StaticDir:      "web/dist",
		Environment:    "test",
		PublicBasePath: "/billtap",
		PublicBaseURL:  "https://localhost:8081",
	})

	customer := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/billtap/v1/customers", map[string]string{"email": "buyer@example.test"})
	product := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/billtap/v1/products", map[string]string{"name": "Team"})
	price := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/billtap/v1/prices", map[string]string{
		"product":             product.ID,
		"currency":            "usd",
		"unit_amount":         "9900",
		"recurring[interval]": "month",
	})
	session := postForm[struct {
		URL string `json:"url"`
	}](t, handler, "/billtap/v1/checkout/sessions", map[string]string{
		"customer":                customer.ID,
		"line_items[0][price]":    price.ID,
		"line_items[0][quantity]": "1",
	})
	if got, want := session.URL, "https://localhost:8081/billtap/checkout/"; len(got) < len(want) || got[:len(want)] != want {
		t.Fatalf("checkout session URL = %q, want prefix %q", got, want)
	}
}

func TestForwardedPrefixPrefixesAPISessionURLs(t *testing.T) {
	handler := newSQLiteBackedServer(t, config.Config{
		Addr:        ":0",
		DatabaseURL: ":memory:",
		StaticDir:   "web/dist",
		Environment: "test",
	})
	headers := map[string]string{"X-Forwarded-Prefix": "/billtap"}

	customer := postFormWithHeaders[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/customers", map[string]string{"email": "buyer@example.test"}, headers)
	product := postFormWithHeaders[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/products", map[string]string{"name": "Team"}, headers)
	price := postFormWithHeaders[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/prices", map[string]string{
		"product":             product.ID,
		"currency":            "usd",
		"unit_amount":         "9900",
		"recurring[interval]": "month",
	}, headers)
	session := postFormWithHeaders[struct {
		URL string `json:"url"`
	}](t, handler, "/v1/checkout/sessions", map[string]string{
		"customer":                customer.ID,
		"line_items[0][price]":    price.ID,
		"line_items[0][quantity]": "1",
	}, headers)
	if got, want := session.URL, "http://example.com/billtap/checkout/"; len(got) < len(want) || got[:len(want)] != want {
		t.Fatalf("checkout session URL = %q, want prefix %q", got, want)
	}
}

func postForm[T any](t *testing.T, handler http.Handler, path string, values map[string]string) T {
	return postFormWithHeaders[T](t, handler, path, values, nil)
}

func postFormWithHeaders[T any](t *testing.T, handler http.Handler, path string, values map[string]string, headers map[string]string) T {
	t.Helper()
	form := make(url.Values)
	for key, value := range values {
		form.Set(key, value)
	}
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code < 200 || rec.Code >= 300 {
		t.Fatalf("POST %s status = %d body = %s", path, rec.Code, rec.Body.String())
	}
	var out T
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode POST %s response: %v body=%s", path, err, rec.Body.String())
	}
	return out
}

func newSQLiteBackedServer(t *testing.T, cfg config.Config) http.Handler {
	t.Helper()
	store, err := storage.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})
	return New(Options{Config: cfg, Store: store})
}
