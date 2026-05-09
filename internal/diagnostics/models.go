package diagnostics

import (
	"context"
	"time"
)

const ObjectRequestTrace = "request_trace"

type RequestTrace struct {
	ID               string              `json:"id"`
	Object           string              `json:"object"`
	Method           string              `json:"method"`
	Path             string              `json:"path"`
	Query            string              `json:"query,omitempty"`
	Status           int                 `json:"status"`
	DurationMS       int64               `json:"duration_ms"`
	RequestID        string              `json:"request_id,omitempty"`
	IdempotencyKey   string              `json:"idempotency_key,omitempty"`
	RequestHeaders   map[string]string   `json:"request_headers,omitempty"`
	RequestBody      string              `json:"request_body,omitempty"`
	ResponseBody     string              `json:"response_body,omitempty"`
	ResponseObject   string              `json:"response_object,omitempty"`
	ResponseObjectID string              `json:"response_object_id,omitempty"`
	ErrorType        string              `json:"error_type,omitempty"`
	ErrorCode        string              `json:"error_code,omitempty"`
	ErrorParam       string              `json:"error_param,omitempty"`
	RelatedIDs       map[string][]string `json:"related_ids,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
}

type RequestTraceFilter struct {
	Method         string
	Path           string
	Status         int
	RequestID      string
	IdempotencyKey string
	ObjectID       string
	Limit          int
}

type Repository interface {
	RecordRequestTrace(context.Context, RequestTrace) error
	ListRequestTraces(context.Context, RequestTraceFilter) ([]RequestTrace, error)
}
