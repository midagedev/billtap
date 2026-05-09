package webhooks

import (
	"context"
	"encoding/json"
	"time"
)

const (
	ObjectEndpoint        = "webhook_endpoint"
	ObjectEvent           = "event"
	ObjectDeliveryAttempt = "delivery_attempt"
	ObjectAuditEntry      = "audit_log_entry"

	APIVersion = "billtap-2026-05-08"

	SourceAPI      = "api"
	SourceCheckout = "checkout"
	SourcePortal   = "portal"
	SourceScenario = "scenario"
	SourceReplay   = "replay"

	StatusScheduled  = "scheduled"
	StatusDelivering = "delivering"
	StatusSucceeded  = "succeeded"
	StatusFailed     = "failed"
	StatusAbandoned  = "abandoned"
	StatusSkipped    = "skipped"
)

type Endpoint struct {
	ID               string     `json:"id"`
	Object           string     `json:"object"`
	URL              string     `json:"url"`
	Secret           string     `json:"secret,omitempty"`
	EnabledEvents    []string   `json:"enabled_events"`
	Active           bool       `json:"active"`
	RetryMaxAttempts int        `json:"retry_max_attempts,omitempty"`
	RetryBackoff     []string   `json:"retry_backoff,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type Event struct {
	ID              string          `json:"id"`
	Object          string          `json:"object"`
	Type            string          `json:"type"`
	Created         int64           `json:"created"`
	Livemode        bool            `json:"livemode"`
	APIVersion      string          `json:"api_version"`
	PendingWebhooks int             `json:"pending_webhooks"`
	Request         EventRequest    `json:"request,omitempty"`
	Data            EventData       `json:"data"`
	Billtap         EventMetadata   `json:"billtap"`
	RawPayload      json.RawMessage `json:"-"`
	CreatedAt       time.Time       `json:"-"`
}

type EventRequest struct {
	ID             string `json:"id,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

type EventData struct {
	Object json.RawMessage `json:"object"`
}

type EventMetadata struct {
	ScenarioRunID string `json:"scenario_run_id,omitempty"`
	Source        string `json:"source"`
	Sequence      int64  `json:"sequence"`
}

type DeliveryAttempt struct {
	ID             string            `json:"id"`
	Object         string            `json:"object"`
	EventID        string            `json:"event_id"`
	EndpointID     string            `json:"endpoint_id"`
	AttemptNumber  int               `json:"attempt_number"`
	Status         string            `json:"status"`
	ScheduledAt    time.Time         `json:"scheduled_at"`
	DeliveredAt    *time.Time        `json:"delivered_at,omitempty"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers,omitempty"`
	RequestBody    json.RawMessage   `json:"request_body,omitempty"`
	ResponseStatus int               `json:"response_status,omitempty"`
	ResponseBody   string            `json:"response_body,omitempty"`
	Error          string            `json:"error,omitempty"`
	NextRetryAt    *time.Time        `json:"next_retry_at,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
}

type AuditEntry struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	Action     string            `json:"action"`
	Actor      string            `json:"actor"`
	TargetType string            `json:"target_type"`
	TargetID   string            `json:"target_id"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

type EndpointFilter struct {
	IncludeDeleted bool
	ActiveOnly     bool
	EventType      string
}

type EventFilter struct {
	Type          string
	ScenarioRunID string
}

type DeliveryAttemptFilter struct {
	EventID    string
	EndpointID string
	Status     string
}

type AuditFilter struct {
	Action   string
	TargetID string
}

type RetentionResult struct {
	Cutoff                   time.Time `json:"cutoff"`
	WebhookEventsRedacted    int64     `json:"webhook_events_redacted"`
	DeliveryAttemptsRedacted int64     `json:"delivery_attempts_redacted"`
}

type EventInput struct {
	Type            string
	ObjectPayload   json.RawMessage
	RequestID       string
	IdempotencyKey  string
	ScenarioRunID   string
	Source          string
	Sequence        int64
	DeliverNow      bool
	DeliveryOptions DeliveryOptions
}

type DeliveryOptions struct {
	Duplicate  int
	Delay      time.Duration
	OutOfOrder bool
	Replay     bool
}

type ReplayOptions struct {
	Duplicate  int
	Delay      time.Duration
	OutOfOrder bool
}

type ServiceOptions struct {
	StoreRawPayloads bool
	RetentionDays    int
}

type Repository interface {
	CreateWebhookEndpoint(ctx context.Context, endpoint Endpoint) (Endpoint, error)
	GetWebhookEndpoint(ctx context.Context, id string) (Endpoint, error)
	ListWebhookEndpoints(ctx context.Context, filter EndpointFilter) ([]Endpoint, error)
	UpdateWebhookEndpoint(ctx context.Context, id string, endpoint Endpoint) (Endpoint, error)
	DeleteWebhookEndpoint(ctx context.Context, id string) (Endpoint, error)

	CreateEvent(ctx context.Context, event Event) (Event, error)
	GetEvent(ctx context.Context, id string) (Event, error)
	ListEvents(ctx context.Context, filter EventFilter) ([]Event, error)

	CreateDeliveryAttempt(ctx context.Context, attempt DeliveryAttempt) (DeliveryAttempt, error)
	ListDeliveryAttempts(ctx context.Context, filter DeliveryAttemptFilter) ([]DeliveryAttempt, error)

	CreateAuditEntry(ctx context.Context, entry AuditEntry) (AuditEntry, error)
	ListAuditEntries(ctx context.Context, filter AuditFilter) ([]AuditEntry, error)
	ApplyRetention(ctx context.Context, cutoff time.Time) (RetentionResult, error)
}
