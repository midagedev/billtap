package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	envConfigPath             = "BILLTAP_CONFIG"
	envAddr                   = "BILLTAP_ADDR"
	envDatabaseURL            = "BILLTAP_DATABASE_URL"
	envStaticDir              = "BILLTAP_STATIC_DIR"
	envEnvironment            = "BILLTAP_ENV"
	envRelayMode              = "BILLTAP_RELAY_MODE"
	envRawPayloadStorage      = "BILLTAP_RAW_PAYLOAD_STORAGE"
	envRetentionDays          = "BILLTAP_RETENTION_DAYS"
	envWebhookSignatureHeader = "BILLTAP_WEBHOOK_SIGNATURE_HEADER"

	RawPayloadStore        = "store"
	RawPayloadMetadataOnly = "metadata_only"

	DefaultWebhookSignatureHeader = "Billtap-Signature"
)

type LookupFunc func(string) (string, bool)

type Config struct {
	Addr                   string `json:"addr"`
	DatabaseURL            string `json:"database_url"`
	StaticDir              string `json:"static_dir"`
	Environment            string `json:"environment"`
	RelayMode              bool   `json:"relay_mode"`
	RawPayloadStorage      string `json:"raw_payload_storage"`
	RetentionDays          int    `json:"retention_days"`
	WebhookSignatureHeader string `json:"webhook_signature_header"`
}

func Default() Config {
	return Config{
		Addr:                   ":8080",
		DatabaseURL:            ".billtap/billtap.db",
		StaticDir:              "dist/app",
		Environment:            "development",
		RelayMode:              false,
		RawPayloadStorage:      RawPayloadStore,
		RetentionDays:          30,
		WebhookSignatureHeader: DefaultWebhookSignatureHeader,
	}
}

func Load(path string) (Config, error) {
	return LoadWithLookup(path, os.LookupEnv)
}

func LoadWithLookup(path string, lookup LookupFunc) (Config, error) {
	cfg := Default()
	if lookup == nil {
		lookup = os.LookupEnv
	}

	if path == "" {
		if value, ok := lookup(envConfigPath); ok {
			path = value
		}
	}

	if path != "" {
		fileCfg, err := loadFile(path)
		if err != nil {
			return Config{}, err
		}
		cfg = merge(cfg, fileCfg)
	}

	if value, ok := lookup(envAddr); ok {
		cfg.Addr = value
	}
	if value, ok := lookup(envDatabaseURL); ok {
		cfg.DatabaseURL = value
	}
	if value, ok := lookup(envStaticDir); ok {
		cfg.StaticDir = value
	}
	if value, ok := lookup(envEnvironment); ok {
		cfg.Environment = value
	}
	if value, ok := lookup(envRelayMode); ok {
		cfg.RelayMode = parseBool(value)
	}
	if value, ok := lookup(envRawPayloadStorage); ok {
		cfg.RawPayloadStorage = value
	}
	if value, ok := lookup(envRetentionDays); ok {
		days, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("parse %s: %w", envRetentionDays, err)
		}
		cfg.RetentionDays = days
	}
	if value, ok := lookup(envWebhookSignatureHeader); ok {
		cfg.WebhookSignatureHeader = value
	}
	if cfg.RelayMode {
		cfg.RawPayloadStorage = RawPayloadMetadataOnly
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Addr == "" {
		return errors.New("addr is required")
	}
	if c.DatabaseURL == "" {
		return errors.New("database_url is required")
	}
	if c.StaticDir == "" {
		return errors.New("static_dir is required")
	}
	if c.Environment == "" {
		return errors.New("environment is required")
	}
	switch c.RawPayloadStorage {
	case RawPayloadStore, RawPayloadMetadataOnly:
	default:
		return fmt.Errorf("raw_payload_storage must be %q or %q", RawPayloadStore, RawPayloadMetadataOnly)
	}
	if c.RetentionDays < 0 {
		return errors.New("retention_days must be zero or greater")
	}
	if c.WebhookSignatureHeader == "" {
		return errors.New("webhook_signature_header is required")
	}
	return nil
}

func loadFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open config file %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config file %q: %w", path, err)
	}
	return cfg, nil
}

func merge(base Config, override Config) Config {
	if override.Addr != "" {
		base.Addr = override.Addr
	}
	if override.DatabaseURL != "" {
		base.DatabaseURL = override.DatabaseURL
	}
	if override.StaticDir != "" {
		base.StaticDir = override.StaticDir
	}
	if override.Environment != "" {
		base.Environment = override.Environment
	}
	if override.RelayMode {
		base.RelayMode = true
	}
	if override.RawPayloadStorage != "" {
		base.RawPayloadStorage = override.RawPayloadStorage
	}
	if override.RetentionDays != 0 {
		base.RetentionDays = override.RetentionDays
	}
	if override.WebhookSignatureHeader != "" {
		base.WebhookSignatureHeader = override.WebhookSignatureHeader
	}
	return base
}

func parseBool(value string) bool {
	switch value {
	case "1", "true", "TRUE", "True", "yes", "YES", "on", "ON":
		return true
	default:
		return false
	}
}
