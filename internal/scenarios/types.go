package scenarios

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ExitPass               = 0
	ExitAssertionFailed    = 1
	ExitInvalidConfig      = 2
	ExitAppCallbackFailure = 3
	ExitRuntimeFailure     = 4

	FailureInvalidConfig = "invalid_config"
	FailureAssertion     = "assertion_failed"
	FailureAppCallback   = "app_callback_failure"
	FailureRunner        = "runner_error"
)

var (
	ErrInvalidConfig      = errors.New("invalid scenario config")
	ErrAssertionFailed    = errors.New("scenario assertion failed")
	ErrAppCallbackFailure = errors.New("app callback failure")
)

type InvalidConfigError struct {
	Problems []string
}

func (e *InvalidConfigError) Error() string {
	if e == nil || len(e.Problems) == 0 {
		return ErrInvalidConfig.Error()
	}
	return fmt.Sprintf("%s: %s", ErrInvalidConfig, strings.Join(e.Problems, "; "))
}

func (e *InvalidConfigError) Unwrap() error {
	return ErrInvalidConfig
}

type Scenario struct {
	Name     string         `yaml:"name" json:"name"`
	Profile  string         `yaml:"profile" json:"profile,omitempty"`
	App      AppConfig      `yaml:"app" json:"app,omitempty"`
	Catalog  Catalog        `yaml:"catalog" json:"catalog,omitempty"`
	Clock    ClockConfig    `yaml:"clock" json:"clock,omitempty"`
	Defaults map[string]any `yaml:"defaults" json:"defaults,omitempty"`
	Webhooks map[string]any `yaml:"webhooks" json:"webhooks,omitempty"`
	SaaS     SaaSConfig     `yaml:"saas" json:"saas,omitempty"`
	Steps    []Step         `yaml:"steps" json:"steps"`
}

type AppConfig struct {
	WebhookURL string           `yaml:"webhookUrl" json:"webhookUrl,omitempty"`
	Assertions AssertionsConfig `yaml:"assertions" json:"assertions,omitempty"`
}

type AssertionsConfig struct {
	BaseURL string `yaml:"baseUrl" json:"baseUrl,omitempty"`
}

type ClockConfig struct {
	Start string `yaml:"start" json:"start,omitempty"`
}

type Catalog struct {
	Products []CatalogProduct `yaml:"products" json:"products,omitempty"`
	Prices   []CatalogPrice   `yaml:"prices" json:"prices,omitempty"`
}

type SaaSConfig struct {
	Tenant        SaaSTenant `yaml:"tenant" json:"tenant,omitempty"`
	CatalogPreset string     `yaml:"catalogPreset" json:"catalogPreset,omitempty"`
}

type SaaSTenant struct {
	ID                 string `yaml:"id" json:"id,omitempty"`
	Rail               string `yaml:"rail" json:"rail,omitempty"`
	ConnectedAccountID string `yaml:"connectedAccountId" json:"connectedAccountId,omitempty"`
}

type CatalogProduct struct {
	ID          string            `yaml:"id" json:"id,omitempty"`
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description,omitempty"`
	Active      *bool             `yaml:"active" json:"active,omitempty"`
	Metadata    map[string]string `yaml:"metadata" json:"metadata,omitempty"`
}

type CatalogPrice struct {
	ID            string            `yaml:"id" json:"id,omitempty"`
	Product       string            `yaml:"product" json:"product"`
	Currency      string            `yaml:"currency" json:"currency"`
	UnitAmount    int64             `yaml:"unitAmount" json:"unitAmount"`
	Interval      string            `yaml:"interval" json:"interval,omitempty"`
	IntervalCount int               `yaml:"intervalCount" json:"intervalCount,omitempty"`
	Active        *bool             `yaml:"active" json:"active,omitempty"`
	Metadata      map[string]string `yaml:"metadata" json:"metadata,omitempty"`
}

type Step struct {
	ID     string         `yaml:"id" json:"id"`
	Action string         `yaml:"action" json:"action"`
	Params map[string]any `yaml:"params" json:"params,omitempty"`
}
