package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/diagnostics"
	"github.com/hckim/billtap/internal/security"
)

func (h *Handler) serveWithRequestTrace(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UTC()
	ensureTraceRequestID(w, r)
	var rawBody []byte
	if r.Body != nil {
		var err error
		rawBody, err = io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
	}
	r.Body = io.NopCloser(bytes.NewReader(rawBody))

	rec := &recordingResponseWriter{ResponseWriter: w, status: http.StatusOK}
	h.serveWithIdempotency(rec, r)
	h.recordRequestTrace(r, rawBody, rec.status, rec.body.Bytes(), start)
}

func (h *Handler) recordRequestTrace(r *http.Request, rawBody []byte, status int, responseBody []byte, start time.Time) {
	if h.diagnostics == nil {
		return
	}
	bodyEvidence := redactedBodyEvidence(r.Header.Get("Content-Type"), rawBody)
	responseEvidence := truncateEvidence(security.RedactText(string(responseBody)))
	object, objectID, errType, errCode, errParam := responseEvidenceFields(responseBody)
	relatedIDs := relatedIDsForTrace(r.URL.Path, string(rawBody), string(responseBody))
	if objectID != "" {
		relatedIDs = addRelatedID(relatedIDs, objectID)
	}

	_ = h.diagnostics.RecordRequestTrace(r.Context(), diagnostics.RequestTrace{
		ID:               "rt_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		Method:           r.Method,
		Path:             r.URL.Path,
		Query:            redactedQueryEvidence(r.URL.RawQuery),
		Status:           status,
		DurationMS:       time.Since(start).Milliseconds(),
		RequestID:        requestIDForTrace(r),
		IdempotencyKey:   strings.TrimSpace(r.Header.Get(idempotencyKeyHeader)),
		RequestHeaders:   redactedRequestHeaders(r.Header),
		RequestBody:      bodyEvidence,
		ResponseBody:     responseEvidence,
		ResponseObject:   object,
		ResponseObjectID: objectID,
		ErrorType:        errType,
		ErrorCode:        errCode,
		ErrorParam:       errParam,
		RelatedIDs:       relatedIDs,
		CreatedAt:        start,
	})
}

func ensureTraceRequestID(w http.ResponseWriter, r *http.Request) {
	requestID := requestIDForTrace(r)
	if r.Header.Get("Request-Id") == "" && r.Header.Get("X-Request-Id") == "" && r.Header.Get("X-Correlation-Id") == "" {
		r.Header.Set("Request-Id", requestID)
	}
	w.Header().Set("Request-Id", requestID)
}

func requestIDForTrace(r *http.Request) string {
	for _, key := range []string{"Request-Id", "X-Request-Id", "X-Correlation-Id"} {
		if value := strings.TrimSpace(r.Header.Get(key)); value != "" {
			return value
		}
	}
	return "req_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

func redactedRequestHeaders(headers http.Header) map[string]string {
	out := map[string]string{}
	for key, values := range headers {
		if len(values) == 0 {
			continue
		}
		out[key] = strings.Join(values, ",")
	}
	return security.RedactHeaders(out)
}

func redactedQueryEvidence(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return truncateEvidence(security.RedactText(rawQuery))
	}
	for key, items := range values {
		if security.IsSensitiveKey(key) || security.IsCardDataKey(key) {
			for i := range items {
				items[i] = security.MaskedValue
			}
			values[key] = items
		}
	}
	return truncateEvidence(values.Encode())
}

func redactedBodyEvidence(contentType string, raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		values, err := url.ParseQuery(string(raw))
		if err == nil {
			for key, items := range values {
				if security.IsSensitiveKey(key) || security.IsCardDataKey(key) {
					for i := range items {
						items[i] = security.MaskedValue
					}
					values[key] = items
				}
			}
			return truncateEvidence(values.Encode())
		}
	}
	return truncateEvidence(security.RedactText(string(raw)))
}

func truncateEvidence(value string) string {
	const maxEvidenceBytes = 32 * 1024
	if len(value) <= maxEvidenceBytes {
		return value
	}
	return value[:maxEvidenceBytes] + "\n...truncated..."
}

func responseEvidenceFields(raw []byte) (object string, objectID string, errorType string, errorCode string, errorParam string) {
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", "", "", "", ""
	}
	if errPayload, ok := payload["error"].(map[string]any); ok {
		return "", "", stringField(errPayload, "type"), stringField(errPayload, "code"), stringField(errPayload, "param")
	}
	return stringField(payload, "object"), stringField(payload, "id"), "", "", ""
}

func stringField(values map[string]any, key string) string {
	if value, ok := values[key].(string); ok {
		return value
	}
	return ""
}

var traceIDPattern = regexp.MustCompile(`\b(cus|prod|price|cs|sub|si|in|pi|evt|we|bpc|pm)_[A-Za-z0-9_]+\b`)

func relatedIDsForTrace(parts ...string) map[string][]string {
	out := map[string][]string{}
	seen := map[string]bool{}
	for _, part := range parts {
		for _, id := range traceIDPattern.FindAllString(part, -1) {
			if seen[id] {
				continue
			}
			seen[id] = true
			out = addRelatedID(out, id)
		}
	}
	return out
}

func addRelatedID(out map[string][]string, id string) map[string][]string {
	if id == "" {
		return out
	}
	if out == nil {
		out = map[string][]string{}
	}
	prefix, _, found := strings.Cut(id, "_")
	if !found {
		return out
	}
	kind := traceIDKind(prefix)
	for _, existing := range out[kind] {
		if existing == id {
			return out
		}
	}
	out[kind] = append(out[kind], id)
	sort.Strings(out[kind])
	return out
}

func traceIDKind(prefix string) string {
	switch prefix {
	case "cus":
		return "customers"
	case "prod":
		return "products"
	case "price":
		return "prices"
	case "cs":
		return "checkout_sessions"
	case "sub":
		return "subscriptions"
	case "si":
		return "subscription_items"
	case "in":
		return "invoices"
	case "pi":
		return "payment_intents"
	case "evt":
		return "events"
	case "we":
		return "webhook_endpoints"
	case "bpc":
		return "billing_portal_sessions"
	case "pm":
		return "payment_methods"
	default:
		return prefix
	}
}
