package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
}
