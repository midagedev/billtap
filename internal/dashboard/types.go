package dashboard

import (
	"encoding/json"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/webhooks"
)

const (
	ObjectList          = "dashboard.object_list"
	ObjectDetailObject  = "dashboard.object_detail"
	ObjectTimeline      = "dashboard.timeline"
	ObjectWebhookDetail = "dashboard.webhook_detail"
	ObjectDebugBundle   = "dashboard.debug_bundle"

	TypeWebhookEvent    = "webhook.event"
	TypeWebhookEndpoint = "webhook.endpoint"
)

type ObjectRef struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type ObjectListQuery struct {
	Type string `json:"type"`
}

type ObjectListResult struct {
	Object string          `json:"object"`
	Type   string          `json:"type"`
	Data   []ObjectSummary `json:"data"`
}

type ObjectSummary struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	Display    string            `json:"display,omitempty"`
	Status     string            `json:"status,omitempty"`
	CustomerID string            `json:"customer,omitempty"`
	Related    map[string]string `json:"related,omitempty"`
	CreatedAt  time.Time         `json:"created_at,omitempty"`
}

type ObjectDetail struct {
	Object  string            `json:"object"`
	Ref     ObjectRef         `json:"ref"`
	Data    any               `json:"data"`
	Related map[string]string `json:"related,omitempty"`
}

type TimelineQuery struct {
	CustomerID        string `json:"customer,omitempty"`
	CheckoutSessionID string `json:"checkout_session,omitempty"`
	SubscriptionID    string `json:"subscription,omitempty"`
	InvoiceID         string `json:"invoice,omitempty"`
	PaymentIntentID   string `json:"payment_intent,omitempty"`
	EventID           string `json:"event,omitempty"`
	EventType         string `json:"event_type,omitempty"`
	ScenarioRunID     string `json:"scenario_run_id,omitempty"`
	IncludeWebhooks   bool   `json:"include_webhooks,omitempty"`
}

type TimelineResult struct {
	Object string         `json:"object"`
	Filter TimelineQuery  `json:"filter"`
	Data   []TimelineItem `json:"data"`
}

type TimelineItem struct {
	ID                string            `json:"id"`
	Kind              string            `json:"kind"`
	Action            string            `json:"action"`
	Message           string            `json:"message"`
	ObjectType        string            `json:"object_type,omitempty"`
	ObjectID          string            `json:"object_id,omitempty"`
	CustomerID        string            `json:"customer,omitempty"`
	CheckoutSessionID string            `json:"checkout_session,omitempty"`
	SubscriptionID    string            `json:"subscription,omitempty"`
	InvoiceID         string            `json:"invoice,omitempty"`
	PaymentIntentID   string            `json:"payment_intent,omitempty"`
	EventID           string            `json:"event,omitempty"`
	DeliveryAttemptID string            `json:"delivery_attempt,omitempty"`
	WebhookStatus     string            `json:"webhook_status,omitempty"`
	Data              map[string]string `json:"data,omitempty"`
	At                time.Time         `json:"at"`
}

type WebhookDetail struct {
	Object          string                    `json:"object"`
	Event           webhooks.Event            `json:"event"`
	AttemptCount    int                       `json:"attempt_count"`
	LatestStatus    string                    `json:"latest_status,omitempty"`
	SignatureHeader string                    `json:"signature_header,omitempty"`
	RequestURL      string                    `json:"request_url,omitempty"`
	RequestBody     json.RawMessage           `json:"request_body,omitempty"`
	Attempts        []DeliveryAttemptEvidence `json:"attempts"`
	RetryPlan       []RetryEvidence           `json:"retry_plan,omitempty"`
	Flags           WebhookFlags              `json:"flags"`
}

type DeliveryAttemptEvidence struct {
	ID              string            `json:"id"`
	EndpointID      string            `json:"endpoint_id"`
	AttemptNumber   int               `json:"attempt_number"`
	Status          string            `json:"status"`
	ScheduledAt     time.Time         `json:"scheduled_at"`
	DeliveredAt     *time.Time        `json:"delivered_at,omitempty"`
	RequestURL      string            `json:"request_url"`
	SignatureHeader string            `json:"signature_header,omitempty"`
	RequestHeaders  map[string]string `json:"request_headers,omitempty"`
	RequestBody     json.RawMessage   `json:"request_body,omitempty"`
	ResponseStatus  int               `json:"response_status,omitempty"`
	ResponseBody    string            `json:"response_body,omitempty"`
	Error           string            `json:"error,omitempty"`
	NextRetryAt     *time.Time        `json:"next_retry_at,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type RetryEvidence struct {
	AttemptID  string    `json:"attempt_id"`
	EndpointID string    `json:"endpoint_id"`
	RetryAt    time.Time `json:"retry_at"`
}

type WebhookFlags struct {
	Duplicate  bool `json:"duplicate"`
	OutOfOrder bool `json:"out_of_order"`
	Replay     bool `json:"replay"`
}

type DebugBundle struct {
	Object      string          `json:"object"`
	Target      ObjectRef       `json:"target"`
	GeneratedAt time.Time       `json:"generated_at"`
	Detail      ObjectDetail    `json:"detail"`
	Timeline    TimelineResult  `json:"timeline"`
	Webhooks    []WebhookDetail `json:"webhooks,omitempty"`
}

func billingTimelineFilter(q TimelineQuery) billing.TimelineFilter {
	return billing.TimelineFilter{
		CustomerID:        q.CustomerID,
		CheckoutSessionID: q.CheckoutSessionID,
		SubscriptionID:    q.SubscriptionID,
		InvoiceID:         q.InvoiceID,
		PaymentIntentID:   q.PaymentIntentID,
	}
}
