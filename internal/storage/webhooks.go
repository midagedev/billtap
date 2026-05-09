package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/webhooks"
)

var _ webhooks.Repository = (*SQLiteStore)(nil)

func (s *SQLiteStore) CreateWebhookEndpoint(ctx context.Context, endpoint webhooks.Endpoint) (webhooks.Endpoint, error) {
	if _, err := s.db.ExecContext(ctx, `INSERT INTO webhook_endpoints (id, url, secret, enabled_events, active, retry_max_attempts, retry_backoff, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		endpoint.ID, endpoint.URL, endpoint.Secret, encodeStrings(endpoint.EnabledEvents), boolInt(endpoint.Active), endpoint.RetryMaxAttempts, encodeStrings(endpoint.RetryBackoff), encodeTime(endpoint.CreatedAt), encodeTime(endpoint.UpdatedAt), encodeOptionalTime(endpoint.DeletedAt)); err != nil {
		return webhooks.Endpoint{}, err
	}
	return s.GetWebhookEndpoint(ctx, endpoint.ID)
}

func (s *SQLiteStore) GetWebhookEndpoint(ctx context.Context, id string) (webhooks.Endpoint, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, url, secret, enabled_events, active, retry_max_attempts, retry_backoff, created_at, updated_at, deleted_at FROM webhook_endpoints WHERE id = ?`, id)
	endpoint, err := scanWebhookEndpoint(row)
	if errors.Is(err, sql.ErrNoRows) {
		return webhooks.Endpoint{}, webhooks.ErrNotFound
	}
	return endpoint, err
}

func (s *SQLiteStore) ListWebhookEndpoints(ctx context.Context, filter webhooks.EndpointFilter) ([]webhooks.Endpoint, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if !filter.IncludeDeleted {
		clauses = append(clauses, "deleted_at IS NULL")
	}
	if filter.ActiveOnly {
		clauses = append(clauses, "active = 1")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, url, secret, enabled_events, active, retry_max_attempts, retry_backoff, created_at, updated_at, deleted_at
		FROM webhook_endpoints WHERE `+strings.Join(clauses, " AND ")+` ORDER BY created_at DESC, id DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []webhooks.Endpoint
	for rows.Next() {
		endpoint, err := scanWebhookEndpoint(rows)
		if err != nil {
			return nil, err
		}
		if filter.EventType == "" || webhooks.EndpointMatches(endpoint, filter.EventType) {
			out = append(out, endpoint)
		}
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpdateWebhookEndpoint(ctx context.Context, id string, in webhooks.Endpoint) (webhooks.Endpoint, error) {
	current, err := s.GetWebhookEndpoint(ctx, id)
	if err != nil {
		return webhooks.Endpoint{}, err
	}
	if in.URL != "" {
		current.URL = in.URL
	}
	if in.Secret != "" {
		current.Secret = in.Secret
	}
	if in.EnabledEvents != nil {
		current.EnabledEvents = in.EnabledEvents
	}
	if in.RetryMaxAttempts > 0 {
		current.RetryMaxAttempts = in.RetryMaxAttempts
	}
	if in.RetryBackoff != nil {
		current.RetryBackoff = in.RetryBackoff
	}
	current.Active = in.Active
	current.UpdatedAt = time.Now().UTC()
	if _, err := s.db.ExecContext(ctx, `UPDATE webhook_endpoints SET url = ?, secret = ?, enabled_events = ?, active = ?, retry_max_attempts = ?, retry_backoff = ?, updated_at = ? WHERE id = ?`,
		current.URL, current.Secret, encodeStrings(current.EnabledEvents), boolInt(current.Active), current.RetryMaxAttempts, encodeStrings(current.RetryBackoff), encodeTime(current.UpdatedAt), id); err != nil {
		return webhooks.Endpoint{}, err
	}
	return s.GetWebhookEndpoint(ctx, id)
}

func (s *SQLiteStore) DeleteWebhookEndpoint(ctx context.Context, id string) (webhooks.Endpoint, error) {
	now := time.Now().UTC()
	if _, err := s.db.ExecContext(ctx, `UPDATE webhook_endpoints SET active = 0, deleted_at = ?, updated_at = ? WHERE id = ?`, encodeTime(now), encodeTime(now), id); err != nil {
		return webhooks.Endpoint{}, err
	}
	return s.GetWebhookEndpoint(ctx, id)
}

func (s *SQLiteStore) CreateEvent(ctx context.Context, event webhooks.Event) (webhooks.Event, error) {
	if len(event.RawPayload) == 0 {
		event.RawPayload = json.RawMessage(`{}`)
	}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO webhook_events (id, type, created, livemode, api_version, pending_webhooks, request_id, idempotency_key, object_payload, raw_payload, source, sequence, scenario_run_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.Type, event.Created, boolInt(event.Livemode), event.APIVersion, event.PendingWebhooks, event.Request.ID, event.Request.IdempotencyKey, string(event.Data.Object), string(event.RawPayload), event.Billtap.Source, event.Billtap.Sequence, event.Billtap.ScenarioRunID, encodeTime(event.CreatedAt)); err != nil {
		return webhooks.Event{}, err
	}
	return s.GetEvent(ctx, event.ID)
}

func (s *SQLiteStore) GetEvent(ctx context.Context, id string) (webhooks.Event, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, type, created, livemode, api_version, pending_webhooks, request_id, idempotency_key, object_payload, raw_payload, source, sequence, scenario_run_id, created_at FROM webhook_events WHERE id = ?`, id)
	event, err := scanWebhookEvent(row)
	if errors.Is(err, sql.ErrNoRows) {
		return webhooks.Event{}, webhooks.ErrNotFound
	}
	return event, err
}

func (s *SQLiteStore) ListEvents(ctx context.Context, filter webhooks.EventFilter) ([]webhooks.Event, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.Type != "" {
		clauses = append(clauses, "type = ?")
		args = append(args, filter.Type)
	}
	if filter.ScenarioRunID != "" {
		clauses = append(clauses, "scenario_run_id = ?")
		args = append(args, filter.ScenarioRunID)
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, type, created, livemode, api_version, pending_webhooks, request_id, idempotency_key, object_payload, raw_payload, source, sequence, scenario_run_id, created_at
		FROM webhook_events WHERE `+strings.Join(clauses, " AND ")+` ORDER BY sequence ASC, created_at ASC, id ASC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []webhooks.Event
	for rows.Next() {
		event, err := scanWebhookEvent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, event)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) CreateDeliveryAttempt(ctx context.Context, attempt webhooks.DeliveryAttempt) (webhooks.DeliveryAttempt, error) {
	if _, err := s.db.ExecContext(ctx, `INSERT INTO delivery_attempts (id, event_id, endpoint_id, attempt_number, status, scheduled_at, delivered_at, request_url, request_headers, request_body, response_status, response_body, error, next_retry_at, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		attempt.ID, attempt.EventID, attempt.EndpointID, attempt.AttemptNumber, attempt.Status, encodeTime(attempt.ScheduledAt), encodeOptionalTime(attempt.DeliveredAt), attempt.RequestURL, encodeMap(attempt.RequestHeaders), string(attempt.RequestBody), attempt.ResponseStatus, attempt.ResponseBody, attempt.Error, encodeOptionalTime(attempt.NextRetryAt), encodeMap(attempt.Metadata), encodeTime(attempt.CreatedAt)); err != nil {
		return webhooks.DeliveryAttempt{}, err
	}
	attempts, err := s.ListDeliveryAttempts(ctx, webhooks.DeliveryAttemptFilter{EventID: attempt.EventID, EndpointID: attempt.EndpointID})
	if err != nil {
		return webhooks.DeliveryAttempt{}, err
	}
	for _, existing := range attempts {
		if existing.ID == attempt.ID {
			return existing, nil
		}
	}
	return attempt, nil
}

func (s *SQLiteStore) ListDeliveryAttempts(ctx context.Context, filter webhooks.DeliveryAttemptFilter) ([]webhooks.DeliveryAttempt, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.EventID != "" {
		clauses = append(clauses, "event_id = ?")
		args = append(args, filter.EventID)
	}
	if filter.EndpointID != "" {
		clauses = append(clauses, "endpoint_id = ?")
		args = append(args, filter.EndpointID)
	}
	if filter.Status != "" {
		clauses = append(clauses, "status = ?")
		args = append(args, filter.Status)
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, event_id, endpoint_id, attempt_number, status, scheduled_at, delivered_at, request_url, request_headers, request_body, response_status, response_body, error, next_retry_at, metadata, created_at
		FROM delivery_attempts WHERE `+strings.Join(clauses, " AND ")+` ORDER BY scheduled_at ASC, attempt_number ASC, id ASC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []webhooks.DeliveryAttempt
	for rows.Next() {
		attempt, err := scanDeliveryAttempt(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, attempt)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) CreateAuditEntry(ctx context.Context, entry webhooks.AuditEntry) (webhooks.AuditEntry, error) {
	if entry.Object == "" {
		entry.Object = webhooks.ObjectAuditEntry
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO audit_log (id, action, actor, target_type, target_id, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Action, entry.Actor, entry.TargetType, entry.TargetID, encodeMap(entry.Metadata), encodeTime(entry.CreatedAt)); err != nil {
		return webhooks.AuditEntry{}, err
	}
	return entry, nil
}

func (s *SQLiteStore) ListAuditEntries(ctx context.Context, filter webhooks.AuditFilter) ([]webhooks.AuditEntry, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.Action != "" {
		clauses = append(clauses, "action = ?")
		args = append(args, filter.Action)
	}
	if filter.TargetID != "" {
		clauses = append(clauses, "target_id = ?")
		args = append(args, filter.TargetID)
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, action, actor, target_type, target_id, metadata, created_at
		FROM audit_log WHERE `+strings.Join(clauses, " AND ")+` ORDER BY created_at ASC, id ASC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []webhooks.AuditEntry
	for rows.Next() {
		var entry webhooks.AuditEntry
		var metadata, createdAt string
		if err := rows.Scan(&entry.ID, &entry.Action, &entry.Actor, &entry.TargetType, &entry.TargetID, &metadata, &createdAt); err != nil {
			return nil, err
		}
		entry.Object = webhooks.ObjectAuditEntry
		entry.Metadata = decodeMap(metadata)
		entry.CreatedAt = decodeTime(createdAt)
		out = append(out, entry)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ApplyRetention(ctx context.Context, cutoff time.Time) (webhooks.RetentionResult, error) {
	result := webhooks.RetentionResult{Cutoff: cutoff}
	events, err := s.db.ExecContext(ctx, `UPDATE webhook_events SET raw_payload = '{}' WHERE created_at < ? AND raw_payload <> '{}'`, encodeTime(cutoff))
	if err != nil {
		return result, err
	}
	result.WebhookEventsRedacted, _ = events.RowsAffected()

	attempts, err := s.db.ExecContext(ctx, `UPDATE delivery_attempts SET request_body = '', response_body = '' WHERE created_at < ? AND (request_body <> '' OR response_body <> '')`, encodeTime(cutoff))
	if err != nil {
		return result, err
	}
	result.DeliveryAttemptsRedacted, _ = attempts.RowsAffected()
	return result, nil
}

func scanWebhookEndpoint(row scanner) (webhooks.Endpoint, error) {
	var endpoint webhooks.Endpoint
	var enabledEvents, retryBackoff, createdAt, updatedAt string
	var active int
	var deletedAt sql.NullString
	if err := row.Scan(&endpoint.ID, &endpoint.URL, &endpoint.Secret, &enabledEvents, &active, &endpoint.RetryMaxAttempts, &retryBackoff, &createdAt, &updatedAt, &deletedAt); err != nil {
		return endpoint, err
	}
	endpoint.Object = webhooks.ObjectEndpoint
	endpoint.EnabledEvents = decodeStrings(enabledEvents)
	endpoint.Active = active != 0
	endpoint.RetryBackoff = decodeStrings(retryBackoff)
	endpoint.CreatedAt = decodeTime(createdAt)
	endpoint.UpdatedAt = decodeTime(updatedAt)
	if deletedAt.Valid {
		t := decodeTime(deletedAt.String)
		endpoint.DeletedAt = &t
	}
	return endpoint, nil
}

func scanWebhookEvent(row scanner) (webhooks.Event, error) {
	var event webhooks.Event
	var objectPayload, rawPayload, source, scenarioRunID, createdAt string
	var livemode int
	if err := row.Scan(&event.ID, &event.Type, &event.Created, &livemode, &event.APIVersion, &event.PendingWebhooks, &event.Request.ID, &event.Request.IdempotencyKey, &objectPayload, &rawPayload, &source, &event.Billtap.Sequence, &scenarioRunID, &createdAt); err != nil {
		return event, err
	}
	event.Object = webhooks.ObjectEvent
	event.Livemode = livemode != 0
	event.Data.Object = json.RawMessage(objectPayload)
	event.RawPayload = json.RawMessage(rawPayload)
	event.Billtap.Source = source
	event.Billtap.ScenarioRunID = scenarioRunID
	event.CreatedAt = decodeTime(createdAt)
	return event, nil
}

func scanDeliveryAttempt(row scanner) (webhooks.DeliveryAttempt, error) {
	var attempt webhooks.DeliveryAttempt
	var scheduledAt, requestHeaders, requestBody, metadata, createdAt string
	var deliveredAt, nextRetryAt sql.NullString
	if err := row.Scan(&attempt.ID, &attempt.EventID, &attempt.EndpointID, &attempt.AttemptNumber, &attempt.Status, &scheduledAt, &deliveredAt, &attempt.RequestURL, &requestHeaders, &requestBody, &attempt.ResponseStatus, &attempt.ResponseBody, &attempt.Error, &nextRetryAt, &metadata, &createdAt); err != nil {
		return attempt, err
	}
	attempt.Object = webhooks.ObjectDeliveryAttempt
	attempt.ScheduledAt = decodeTime(scheduledAt)
	if deliveredAt.Valid {
		t := decodeTime(deliveredAt.String)
		attempt.DeliveredAt = &t
	}
	attempt.RequestHeaders = decodeMap(requestHeaders)
	attempt.RequestBody = json.RawMessage(requestBody)
	if nextRetryAt.Valid {
		t := decodeTime(nextRetryAt.String)
		attempt.NextRetryAt = &t
	}
	attempt.Metadata = decodeMap(metadata)
	attempt.CreatedAt = decodeTime(createdAt)
	return attempt, nil
}

func encodeStrings(values []string) string {
	if values == nil {
		return "[]"
	}
	raw, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(raw)
}

func decodeStrings(raw string) []string {
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}
