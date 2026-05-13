package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hckim/billtap/internal/api"
	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/config"
	"github.com/hckim/billtap/internal/diagnostics"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
)

type Options struct {
	Config config.Config
	Store  storage.Store
}

type Server struct {
	cfg   config.Config
	store storage.Store
	mux   *http.ServeMux
}

func New(opts Options) http.Handler {
	s := &Server{
		cfg:   opts.Config,
		store: opts.Store,
		mux:   http.NewServeMux(),
	}
	s.cfg.PublicBasePath = config.NormalizePublicBasePath(s.cfg.PublicBasePath)
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	basePath := s.requestBasePath(r)
	if basePath == "" {
		s.mux.ServeHTTP(w, r)
		return
	}
	r2 := r.Clone(r.Context())
	r2.Header = r.Header.Clone()
	if r2.Header.Get("X-Forwarded-Prefix") == "" {
		r2.Header.Set("X-Forwarded-Prefix", basePath)
	}
	u := *r.URL
	if stripped, ok := stripBasePath(u.Path, basePath); ok {
		u.Path = stripped
		u.RawPath = ""
	}
	r2.URL = &u
	s.mux.ServeHTTP(w, r2)
}

func (s *Server) routes() {
	if repo, ok := s.store.(billing.Repository); ok {
		var webhookService *webhooks.Service
		if webhookRepo, ok := s.store.(webhooks.Repository); ok {
			webhookService = webhooks.NewServiceWithOptions(webhookRepo, webhooks.ServiceOptions{
				StoreRawPayloads:    s.cfg.RawPayloadStorage != config.RawPayloadMetadataOnly,
				RetentionDays:       s.cfg.RetentionDays,
				SignatureHeaderName: s.cfg.WebhookSignatureHeader,
				APIVersion:          s.cfg.WebhookAPIVersion,
			})
		}
		var diagnosticsService *diagnostics.Service
		if diagnosticsRepo, ok := s.store.(diagnostics.Repository); ok {
			diagnosticsService = diagnostics.NewService(diagnosticsRepo)
		}
		handler := api.New(api.Options{
			Billing:       billing.NewService(repo),
			Webhooks:      webhookService,
			Diagnostics:   diagnosticsService,
			PublicBaseURL: publicBaseURLWithPath(s.cfg.PublicBaseURL, s.cfg.PublicBasePath),
		})
		s.mux.Handle("/v1/", handler)
		s.mux.Handle("/api/", handler)
	}
	s.mux.HandleFunc("/", s.handleRoot)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/healthz", s.handleHealth)
	s.mux.HandleFunc("/readyz", s.handleReady)
	s.mux.HandleFunc("/app/", s.handleApp)
	s.mux.HandleFunc("/checkout", s.handleHostedCheckout)
	s.mux.HandleFunc("/checkout/", s.handleHostedCheckout)
	s.mux.HandleFunc("/portal", s.handleHostedPortal)
	s.mux.HandleFunc("/portal/", s.handleHostedPortal)
	s.mux.HandleFunc("/assets/", s.handleAssets)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w)
		return
	}
	http.Redirect(w, r, s.prefixedPath(r, "/app/dashboard/"), http.StatusFound)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w)
		return
	}

	writeJSON(w, r, http.StatusOK, map[string]string{
		"status":      "ok",
		"environment": s.cfg.Environment,
	})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w)
		return
	}

	status := http.StatusOK
	body := map[string]string{
		"status":      "ok",
		"environment": s.cfg.Environment,
	}

	if s.store == nil {
		status = http.StatusServiceUnavailable
		body["status"] = "degraded"
		body["storage"] = "missing"
	} else if err := s.store.Ping(r.Context()); err != nil {
		status = http.StatusServiceUnavailable
		body["status"] = "degraded"
		body["storage"] = "error"
	} else {
		body["storage"] = "ok"
	}

	writeJSON(w, r, status, body)
}

func (s *Server) handleApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w)
		return
	}

	if s.serveBuiltApp(w, r) {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	assetPath := s.prefixedPath(r, "/assets/app.js")
	_, _ = w.Write([]byte(`<!doctype html><html><head><title>Billtap</title></head><body><div id="root" data-billtap-app="stub"></div><script type="module" src="` + assetPath + `"></script></body></html>`))
}

func (s *Server) handleHostedCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w)
		return
	}
	if s.cfg.StaticDir != "" && serveFileIfExists(w, r, filepath.Join(s.cfg.StaticDir, "checkout", "index.html")) {
		return
	}
	sessionID := strings.Trim(strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/checkout"), "/"), "/")
	target := s.prefixedPath(r, "/app/checkout/")
	if sessionID != "" {
		target += "?session_id=" + sessionID
	} else if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (s *Server) handleHostedPortal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w)
		return
	}
	if s.cfg.StaticDir != "" && serveFileIfExists(w, r, filepath.Join(s.cfg.StaticDir, "portal", "index.html")) {
		return
	}
	customerID := strings.Trim(strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/portal"), "/"), "/")
	target := s.prefixedPath(r, "/app/portal/")
	if customerID != "" {
		target += "?customer_id=" + customerID
	} else if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (s *Server) requestBasePath(r *http.Request) string {
	if prefix := forwardedPrefix(r); prefix != "" {
		return prefix
	}
	return config.NormalizePublicBasePath(s.cfg.PublicBasePath)
}

func forwardedPrefix(r *http.Request) string {
	raw := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Prefix"), ",")[0])
	return config.NormalizePublicBasePath(raw)
}

func stripBasePath(path string, basePath string) (string, bool) {
	if basePath == "" {
		return path, false
	}
	if path == basePath {
		return "/", true
	}
	if strings.HasPrefix(path, basePath+"/") {
		stripped := strings.TrimPrefix(path, basePath)
		if stripped == "" {
			return "/", true
		}
		return stripped, true
	}
	return path, false
}

func (s *Server) prefixedPath(r *http.Request, path string) string {
	return joinURLPath(s.requestBasePath(r), path)
}

func joinURLPath(basePath string, path string) string {
	basePath = config.NormalizePublicBasePath(basePath)
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if basePath == "" {
		return path
	}
	return basePath + path
}

func publicBaseURLWithPath(baseURL string, basePath string) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	basePath = config.NormalizePublicBasePath(basePath)
	if baseURL == "" || basePath == "" || strings.HasSuffix(baseURL, basePath) {
		return baseURL
	}
	return baseURL + basePath
}

func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/assets/")
	if name == "" || strings.Contains(name, "..") {
		http.NotFound(w, r)
		return
	}

	if s.cfg.StaticDir != "" {
		assetPath := filepath.Join(s.cfg.StaticDir, "assets", name)
		if info, err := os.Stat(assetPath); err == nil && !info.IsDir() {
			http.ServeFile(w, r, assetPath)
			return
		}
	}

	if name == "app.js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`console.info("Billtap React asset path stub");`))
		}
		return
	}

	http.NotFound(w, r)
}

func (s *Server) serveBuiltApp(w http.ResponseWriter, r *http.Request) bool {
	if s.cfg.StaticDir == "" {
		return false
	}

	rel := strings.TrimPrefix(r.URL.Path, "/app/")
	if rel == "" {
		return serveFileIfExists(w, r, filepath.Join(s.cfg.StaticDir, "dashboard", "index.html"))
	}
	if strings.Contains(rel, "..") {
		http.NotFound(w, r)
		return true
	}

	cleanRel := filepath.Clean(rel)
	if cleanRel == "." {
		return serveFileIfExists(w, r, filepath.Join(s.cfg.StaticDir, "dashboard", "index.html"))
	}

	if strings.HasPrefix(cleanRel, "assets"+string(filepath.Separator)) || strings.Contains(filepath.Base(cleanRel), ".") {
		if serveFileIfExists(w, r, filepath.Join(s.cfg.StaticDir, cleanRel)) {
			return true
		}
	}

	parts := strings.Split(cleanRel, string(filepath.Separator))
	switch parts[0] {
	case "checkout", "dashboard", "portal":
		if len(parts) == 1 || parts[1] == "" || r.URL.Path == "/app/"+parts[0]+"/" {
			return serveFileIfExists(w, r, filepath.Join(s.cfg.StaticDir, parts[0], "index.html"))
		}
	}

	return false
}

func serveFileIfExists(w http.ResponseWriter, r *http.Request, path string) bool {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		http.ServeFile(w, r, path)
		return true
	}
	return false
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if r.Method == http.MethodHead {
		return
	}
	_ = json.NewEncoder(w).Encode(value)
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Allow", "GET, HEAD")
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
