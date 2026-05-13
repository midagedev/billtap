package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := LoadWithLookup("", func(string) (string, bool) { return "", false })
	if err != nil {
		t.Fatalf("LoadWithLookup returned error: %v", err)
	}

	if cfg.Addr != ":8080" {
		t.Fatalf("Addr = %q, want %q", cfg.Addr, ":8080")
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("DatabaseURL is empty")
	}
	if cfg.RawPayloadStorage != RawPayloadStore || cfg.RetentionDays != 30 {
		t.Fatalf("boundary defaults raw=%q retention=%d", cfg.RawPayloadStorage, cfg.RetentionDays)
	}
	if cfg.WebhookSignatureHeader != DefaultWebhookSignatureHeader {
		t.Fatalf("WebhookSignatureHeader = %q, want %q", cfg.WebhookSignatureHeader, DefaultWebhookSignatureHeader)
	}
	if cfg.WebhookAPIVersion != DefaultWebhookAPIVersion {
		t.Fatalf("WebhookAPIVersion = %q, want %q", cfg.WebhookAPIVersion, DefaultWebhookAPIVersion)
	}
}

func TestLoadConfigFileThenEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "billtap.json")
	body := `{"addr":":9000","database_url":"file.db","static_dir":"static","public_base_path":"/file-prefix","public_base_url":"https://billtap.example.test","environment":"test","raw_payload_storage":"metadata_only","retention_days":7,"webhook_signature_header":"Billtap-Signature","webhook_api_version":"2025-03-31.basil"}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	env := map[string]string{
		envAddr:                   ":9100",
		envDatabaseURL:            ":memory:",
		envPublicBasePath:         "generic-prefix",
		envBilltapPublicBasePath:  "/billtap/",
		envPublicBaseURL:          "http://127.0.0.1:18080",
		envWebhookSignatureHeader: "Stripe-Signature",
		envWebhookAPIVersion:      "2025-12-15.clover",
	}
	cfg, err := LoadWithLookup(path, func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	})
	if err != nil {
		t.Fatalf("LoadWithLookup returned error: %v", err)
	}

	if cfg.Addr != ":9100" {
		t.Fatalf("Addr = %q, want env override", cfg.Addr)
	}
	if cfg.DatabaseURL != ":memory:" {
		t.Fatalf("DatabaseURL = %q, want env override", cfg.DatabaseURL)
	}
	if cfg.StaticDir != "static" {
		t.Fatalf("StaticDir = %q, want file value", cfg.StaticDir)
	}
	if cfg.PublicBasePath != "/billtap" {
		t.Fatalf("PublicBasePath = %q, want normalized app-specific env override", cfg.PublicBasePath)
	}
	if cfg.PublicBaseURL != "http://127.0.0.1:18080" {
		t.Fatalf("PublicBaseURL = %q, want env override", cfg.PublicBaseURL)
	}
	if cfg.RawPayloadStorage != RawPayloadMetadataOnly || cfg.RetentionDays != 7 {
		t.Fatalf("boundary config raw=%q retention=%d", cfg.RawPayloadStorage, cfg.RetentionDays)
	}
	if cfg.WebhookSignatureHeader != "Stripe-Signature" {
		t.Fatalf("WebhookSignatureHeader = %q, want env override", cfg.WebhookSignatureHeader)
	}
	if cfg.WebhookAPIVersion != "2025-12-15.clover" {
		t.Fatalf("WebhookAPIVersion = %q, want env override", cfg.WebhookAPIVersion)
	}
}

func TestLoadRejectsInvalidPublicBasePath(t *testing.T) {
	_, err := LoadWithLookup("", func(key string) (string, bool) {
		if key == envPublicBasePath {
			return "https://example.test/billtap", true
		}
		return "", false
	})
	if err == nil {
		t.Fatal("LoadWithLookup succeeded, want public_base_path validation error")
	}
}

func TestRelayModeDisablesRawPayloadStorage(t *testing.T) {
	env := map[string]string{
		envRelayMode:         "true",
		envRawPayloadStorage: RawPayloadStore,
	}
	cfg, err := LoadWithLookup("", func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	})
	if err != nil {
		t.Fatalf("LoadWithLookup returned error: %v", err)
	}
	if !cfg.RelayMode || cfg.RawPayloadStorage != RawPayloadMetadataOnly {
		t.Fatalf("relay cfg = %#v, want relay mode with metadata-only raw payload storage", cfg)
	}
}

func TestLoadRejectsUnknownConfigFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "billtap.json")
	if err := os.WriteFile(path, []byte(`{"unknown":true}`), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	if _, err := LoadWithLookup(path, nil); err == nil {
		t.Fatal("LoadWithLookup succeeded, want unknown field error")
	}
}
