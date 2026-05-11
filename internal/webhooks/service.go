package webhooks

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Service struct {
	repo                Repository
	now                 func() time.Time
	client              *http.Client
	storeRawPayloads    bool
	retentionDays       int
	signatureHeaderName string
	apiVersion          string
}

func NewService(repo Repository) *Service {
	return NewServiceWithOptions(repo, ServiceOptions{StoreRawPayloads: true, RetentionDays: 30})
}

func NewServiceWithOptions(repo Repository, opts ServiceOptions) *Service {
	if opts.RetentionDays == 0 {
		opts.RetentionDays = 30
	}
	if strings.TrimSpace(opts.APIVersion) == "" {
		opts.APIVersion = DefaultAPIVersion
	}
	return &Service{
		repo:                repo,
		now:                 func() time.Time { return time.Now().UTC() },
		storeRawPayloads:    opts.StoreRawPayloads,
		retentionDays:       opts.RetentionDays,
		signatureHeaderName: NormalizeSignatureHeaderName(opts.SignatureHeaderName),
		apiVersion:          opts.APIVersion,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *Service) CreateEndpoint(ctx context.Context, in Endpoint) (Endpoint, error) {
	if strings.TrimSpace(in.URL) == "" {
		return Endpoint{}, fmt.Errorf("%w: endpoint URL is required", ErrInvalidInput)
	}
	now := s.now()
	in.ID = id("we")
	in.Object = ObjectEndpoint
	if in.Secret == "" {
		in.Secret = "btsec_" + randomHex(16)
	}
	if in.RetryMaxAttempts == 0 {
		in.RetryMaxAttempts = 5
	}
	if len(in.RetryBackoff) == 0 {
		in.RetryBackoff = []string{"10s", "30s", "2m", "10m"}
	}
	in.Active = true
	in.CreatedAt = now
	in.UpdatedAt = now
	return s.repo.CreateWebhookEndpoint(ctx, in)
}

func (s *Service) GetEndpoint(ctx context.Context, endpointID string) (Endpoint, error) {
	return s.repo.GetWebhookEndpoint(ctx, endpointID)
}

func (s *Service) ListEndpoints(ctx context.Context, filter EndpointFilter) ([]Endpoint, error) {
	return s.repo.ListWebhookEndpoints(ctx, filter)
}

func (s *Service) UpdateEndpoint(ctx context.Context, endpointID string, in Endpoint) (Endpoint, error) {
	in.UpdatedAt = s.now()
	endpoint, err := s.repo.UpdateWebhookEndpoint(ctx, endpointID, in)
	if err != nil {
		return Endpoint{}, err
	}
	if in.Secret != "" {
		_ = s.audit(ctx, "webhook.endpoint_secret_update", "webhook_endpoint", endpointID, nil)
	}
	return endpoint, nil
}

func (s *Service) DeleteEndpoint(ctx context.Context, endpointID string) (Endpoint, error) {
	return s.repo.DeleteWebhookEndpoint(ctx, endpointID)
}

func (s *Service) CreateEvent(ctx context.Context, in EventInput) (Event, []DeliveryAttempt, error) {
	if strings.TrimSpace(in.Type) == "" {
		return Event{}, nil, fmt.Errorf("%w: event type is required", ErrInvalidInput)
	}
	if len(in.ObjectPayload) == 0 {
		in.ObjectPayload = json.RawMessage(`{}`)
	}
	if in.Source == "" {
		in.Source = SourceAPI
	}
	if in.Sequence == 0 {
		in.Sequence = s.now().UnixNano()
	}

	endpoints, err := s.repo.ListWebhookEndpoints(ctx, EndpointFilter{ActiveOnly: true, EventType: in.Type})
	if err != nil {
		return Event{}, nil, err
	}

	now := s.now()
	event := Event{
		ID:              id("evt"),
		Object:          ObjectEvent,
		Type:            in.Type,
		Created:         now.Unix(),
		Livemode:        false,
		APIVersion:      s.apiVersion,
		PendingWebhooks: len(endpoints),
		Request: EventRequest{
			ID:             in.RequestID,
			IdempotencyKey: in.IdempotencyKey,
		},
		Data: EventData{Object: in.ObjectPayload},
		Billtap: EventMetadata{
			ScenarioRunID: in.ScenarioRunID,
			Source:        in.Source,
			Sequence:      in.Sequence,
		},
		CreatedAt: now,
	}
	raw, err := json.Marshal(event)
	if err != nil {
		return Event{}, nil, err
	}
	event.RawPayload = raw

	storedEvent := event
	if !s.storeRawPayloads {
		storedEvent.RawPayload = nil
	}
	event, err = s.repo.CreateEvent(ctx, storedEvent)
	if err != nil {
		return Event{}, nil, err
	}
	event.RawPayload = raw

	attempts, err := s.createAttempts(ctx, event, endpoints, in.DeliveryOptions)
	if err != nil {
		return Event{}, nil, err
	}
	if hasDeliveryOverride(in.DeliveryOptions) {
		_ = s.audit(ctx, "webhook.delivery_override", "event", event.ID, deliveryMetadata(in.DeliveryOptions))
	}
	return event, attempts, nil
}

func (s *Service) GetEvent(ctx context.Context, eventID string) (Event, error) {
	return s.repo.GetEvent(ctx, eventID)
}

func (s *Service) ListEvents(ctx context.Context, filter EventFilter) ([]Event, error) {
	return s.repo.ListEvents(ctx, filter)
}

func (s *Service) ListDeliveryAttempts(ctx context.Context, filter DeliveryAttemptFilter) ([]DeliveryAttempt, error) {
	return s.repo.ListDeliveryAttempts(ctx, filter)
}

func (s *Service) ReplayEvent(ctx context.Context, eventID string, opts ReplayOptions) ([]DeliveryAttempt, error) {
	event, err := s.repo.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	event, err = s.ensureEventRawPayload(event)
	if err != nil {
		return nil, err
	}
	endpoints, err := s.repo.ListWebhookEndpoints(ctx, EndpointFilter{ActiveOnly: true, EventType: event.Type})
	if err != nil {
		return nil, err
	}
	deliveryOpts := DeliveryOptions{
		Duplicate:          opts.Duplicate,
		Delay:              opts.Delay,
		OutOfOrder:         opts.OutOfOrder,
		Replay:             true,
		ResponseStatus:     opts.ResponseStatus,
		ResponseBody:       opts.ResponseBody,
		SimulatedError:     opts.SimulatedError,
		SimulatedTimeout:   opts.SimulatedTimeout,
		SignatureMismatch:  opts.SignatureMismatch,
		SimulateAppFailure: opts.SimulateAppFailure,
	}
	attempts, err := s.createAttempts(ctx, event, endpoints, deliveryOpts)
	if err != nil {
		return nil, err
	}
	_ = s.audit(ctx, "webhook.replay", "event", event.ID, deliveryMetadata(deliveryOpts))
	return attempts, nil
}

func (s *Service) ReplayHistoricalForEndpoint(ctx context.Context, endpointID string, opts HistoricalReplayOptions) (HistoricalReplayResult, error) {
	endpoint, err := s.repo.GetWebhookEndpoint(ctx, endpointID)
	if err != nil {
		return HistoricalReplayResult{}, err
	}
	if !endpoint.Active || endpoint.DeletedAt != nil {
		return HistoricalReplayResult{}, fmt.Errorf("%w: endpoint is not active", ErrInvalidInput)
	}
	until := opts.Until
	if until.IsZero() {
		until = endpoint.CreatedAt
	}
	if until.IsZero() {
		until = s.now()
	}
	if !opts.Since.IsZero() && opts.Since.After(until) {
		return HistoricalReplayResult{}, fmt.Errorf("%w: since must be before until", ErrInvalidInput)
	}

	events, err := s.repo.ListEvents(ctx, EventFilter{})
	if err != nil {
		return HistoricalReplayResult{}, err
	}

	deliveryOpts := opts.DeliveryOptions
	deliveryOpts.Replay = true
	deliveryOpts.Historical = true
	if deliveryOpts.Duplicate <= 0 {
		deliveryOpts.Duplicate = 1
	}

	result := HistoricalReplayResult{
		Object:     "historical_webhook_replay",
		EndpointID: endpoint.ID,
	}
	if !opts.Since.IsZero() {
		since := opts.Since
		result.Since = &since
	}
	result.Until = &until

	typeFilter := Endpoint{Active: true, EnabledEvents: opts.EventTypes}
	for _, event := range events {
		if !opts.Since.IsZero() && event.CreatedAt.Before(opts.Since) {
			continue
		}
		if event.CreatedAt.After(until) {
			continue
		}
		if len(opts.EventTypes) > 0 && !EndpointMatches(typeFilter, event.Type) {
			continue
		}
		if !EndpointMatches(endpoint, event.Type) {
			continue
		}
		result.MatchedEvents++
		if opts.Limit > 0 && result.ReplayedEvents >= opts.Limit {
			result.SkippedEvents++
			continue
		}
		existing, err := s.repo.ListDeliveryAttempts(ctx, DeliveryAttemptFilter{EventID: event.ID, EndpointID: endpoint.ID})
		if err != nil {
			return result, err
		}
		if len(existing) > 0 && !opts.Force {
			result.SkippedEvents++
			continue
		}
		event, err = s.ensureEventRawPayload(event)
		if err != nil {
			return result, err
		}
		attempts, err := s.createAttempts(ctx, event, []Endpoint{endpoint}, deliveryOpts)
		if err != nil {
			return result, err
		}
		result.Events = append(result.Events, event)
		result.Attempts = append(result.Attempts, attempts...)
		result.ReplayedEvents++
		result.AttemptCount += len(attempts)
	}
	_ = s.audit(ctx, "webhook.replay_historical", "webhook_endpoint", endpoint.ID, map[string]string{
		"since":           timeString(opts.Since),
		"until":           until.Format(time.RFC3339Nano),
		"matched_events":  fmt.Sprintf("%d", result.MatchedEvents),
		"replayed_events": fmt.Sprintf("%d", result.ReplayedEvents),
		"skipped_events":  fmt.Sprintf("%d", result.SkippedEvents),
	})
	return result, nil
}

func (s *Service) ensureEventRawPayload(event Event) (Event, error) {
	if !s.storeRawPayloads || len(bytes.TrimSpace(event.RawPayload)) == 0 || bytes.Equal(bytes.TrimSpace(event.RawPayload), []byte(`{}`)) {
		raw, err := json.Marshal(event)
		if err != nil {
			return Event{}, err
		}
		event.RawPayload = raw
	}
	return event, nil
}

func (s *Service) ListAuditEntries(ctx context.Context, filter AuditFilter) ([]AuditEntry, error) {
	return s.repo.ListAuditEntries(ctx, filter)
}

func (s *Service) ApplyRetention(ctx context.Context) (RetentionResult, error) {
	if s.retentionDays <= 0 {
		return RetentionResult{Cutoff: s.now()}, nil
	}
	cutoff := s.now().Add(-time.Duration(s.retentionDays) * 24 * time.Hour)
	result, err := s.repo.ApplyRetention(ctx, cutoff)
	if err != nil {
		return RetentionResult{}, err
	}
	_ = s.audit(ctx, "retention.apply", "webhook_evidence", "retention", map[string]string{
		"cutoff": cutoff.Format(time.RFC3339),
	})
	return result, nil
}

func (s *Service) createAttempts(ctx context.Context, event Event, endpoints []Endpoint, opts DeliveryOptions) ([]DeliveryAttempt, error) {
	if err := validateDeliveryOptions(opts); err != nil {
		return nil, err
	}
	if opts.SimulateAppFailure != nil {
		return s.createFailureInjectedAttempts(ctx, event, endpoints, opts)
	}
	duplicates := opts.Duplicate
	if duplicates <= 0 {
		duplicates = 1
	}

	var attempts []DeliveryAttempt
	for _, endpoint := range endpoints {
		for i := 0; i < duplicates; i++ {
			attempt, err := s.createAttempt(ctx, event, endpoint, opts, i)
			if err != nil {
				return attempts, err
			}
			attempts = append(attempts, attempt)
		}
	}
	return attempts, nil
}

func validateDeliveryOptions(opts DeliveryOptions) error {
	if opts.ResponseStatus != 0 && (opts.ResponseStatus < 100 || opts.ResponseStatus > 599) {
		return fmt.Errorf("%w: response status must be between 100 and 599", ErrInvalidInput)
	}
	if opts.SimulateAppFailure != nil {
		if opts.SimulateAppFailure.Status < 500 || opts.SimulateAppFailure.Status > 599 {
			return fmt.Errorf("%w: simulate_app_failure.status must be between 500 and 599", ErrInvalidInput)
		}
	}
	return nil
}

func (s *Service) createFailureInjectedAttempts(ctx context.Context, event Event, endpoints []Endpoint, opts DeliveryOptions) ([]DeliveryAttempt, error) {
	failure := *opts.SimulateAppFailure
	if failure.FailFirstNAttempts <= 0 {
		failure.FailFirstNAttempts = 1
	}
	var attempts []DeliveryAttempt
	for _, endpoint := range endpoints {
		for i := 0; i < failure.FailFirstNAttempts; i++ {
			failOpts := opts
			failOpts.ResponseStatus = failure.Status
			failOpts.ResponseBody = failure.Body
			failOpts.SimulatedError = ""
			failOpts.SimulatedTimeout = false
			failOpts.NoRetrySchedule = true
			attempt, err := s.createAttempt(ctx, event, endpoint, failOpts, i)
			if err != nil {
				return attempts, err
			}
			attempts = append(attempts, attempt)
		}
		successOpts := opts
		successOpts.ResponseStatus = 0
		successOpts.ResponseBody = ""
		successOpts.SimulatedError = ""
		successOpts.SimulatedTimeout = false
		successOpts.NoRetrySchedule = false
		successOpts.SimulateAppFailure = nil
		attempt, err := s.createAttempt(ctx, event, endpoint, successOpts, failure.FailFirstNAttempts)
		if err != nil {
			return attempts, err
		}
		attempts = append(attempts, attempt)
	}
	return attempts, nil
}

func (s *Service) createAttempt(ctx context.Context, event Event, endpoint Endpoint, opts DeliveryOptions, duplicateIndex int) (DeliveryAttempt, error) {
	now := s.now()
	scheduledAt := now.Add(opts.Delay)
	existing, err := s.repo.ListDeliveryAttempts(ctx, DeliveryAttemptFilter{EventID: event.ID, EndpointID: endpoint.ID})
	if err != nil {
		return DeliveryAttempt{}, err
	}

	metadata := map[string]string{}
	if opts.Replay {
		metadata["source"] = SourceReplay
	}
	if opts.Historical {
		metadata["historical"] = "true"
	}
	if duplicateIndex > 0 {
		metadata["duplicate"] = "true"
	}
	if opts.OutOfOrder {
		metadata["out_of_order"] = "true"
	}
	if opts.ResponseStatus != 0 {
		metadata["response_status"] = fmt.Sprintf("%d", opts.ResponseStatus)
	}
	if opts.SimulatedTimeout {
		metadata["timeout"] = "true"
	}
	if opts.SimulatedError != "" {
		metadata["error"] = opts.SimulatedError
	}
	if opts.SignatureMismatch {
		metadata["signature_mismatch"] = "true"
	}
	if opts.SimulateAppFailure != nil {
		metadata["simulate_app_failure"] = "true"
		metadata["fail_first_n_attempts"] = fmt.Sprintf("%d", opts.SimulateAppFailure.FailFirstNAttempts)
		if duplicateIndex < opts.SimulateAppFailure.FailFirstNAttempts {
			metadata["simulated_attempt"] = "true"
		}
	}

	signatureHeader := SignatureHeader(endpoint.Secret, scheduledAt, event.RawPayload)
	if opts.SignatureMismatch {
		signatureHeader = fmt.Sprintf("t=%d,v1=%s", scheduledAt.Unix(), strings.Repeat("0", 64))
	}

	attempt := DeliveryAttempt{
		ID:            id("delatt"),
		Object:        ObjectDeliveryAttempt,
		EventID:       event.ID,
		EndpointID:    endpoint.ID,
		AttemptNumber: len(existing) + 1,
		Status:        StatusScheduled,
		ScheduledAt:   scheduledAt,
		RequestURL:    endpoint.URL,
		RequestHeaders: map[string]string{
			s.signatureHeaderName: signatureHeader,
			"Content-Type":        "application/json",
		},
		RequestBody: event.RawPayload,
		Metadata:    metadata,
		CreatedAt:   now,
	}

	if opts.Delay > 0 {
		return s.createDeliveryAttempt(ctx, attempt)
	}
	if hasSimulatedDeliveryResult(opts) {
		return s.recordSimulatedAttempt(ctx, endpoint, attempt, opts)
	}

	return s.deliverAttempt(ctx, endpoint, attempt)
}

func hasSimulatedDeliveryResult(opts DeliveryOptions) bool {
	return opts.ResponseStatus != 0 || opts.SimulatedError != "" || opts.SimulatedTimeout
}

func (s *Service) recordSimulatedAttempt(ctx context.Context, endpoint Endpoint, attempt DeliveryAttempt, opts DeliveryOptions) (DeliveryAttempt, error) {
	deliveredAt := s.now()
	attempt.DeliveredAt = &deliveredAt
	if opts.SimulatedTimeout {
		attempt.Status = StatusFailed
		attempt.Error = "simulated timeout"
		if opts.NoRetrySchedule {
			return s.createDeliveryAttempt(ctx, attempt)
		}
		return s.recordRetry(ctx, endpoint, attempt)
	}
	if opts.SimulatedError != "" {
		attempt.Status = StatusFailed
		attempt.Error = opts.SimulatedError
		if opts.NoRetrySchedule {
			return s.createDeliveryAttempt(ctx, attempt)
		}
		return s.recordRetry(ctx, endpoint, attempt)
	}

	attempt.ResponseStatus = opts.ResponseStatus
	attempt.ResponseBody = opts.ResponseBody
	if attempt.ResponseStatus >= 200 && attempt.ResponseStatus < 300 {
		attempt.Status = StatusSucceeded
		return s.createDeliveryAttempt(ctx, attempt)
	}
	attempt.Status = StatusFailed
	attempt.Error = http.StatusText(attempt.ResponseStatus)
	if attempt.Error == "" {
		attempt.Error = fmt.Sprintf("HTTP %d", attempt.ResponseStatus)
	}
	if opts.NoRetrySchedule {
		return s.createDeliveryAttempt(ctx, attempt)
	}
	return s.recordRetry(ctx, endpoint, attempt)
}

func (s *Service) deliverAttempt(ctx context.Context, endpoint Endpoint, attempt DeliveryAttempt) (DeliveryAttempt, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.URL, bytes.NewReader(attempt.RequestBody))
	if err != nil {
		attempt.Status = StatusFailed
		attempt.Error = err.Error()
		return s.recordRetry(ctx, endpoint, attempt)
	}
	for key, value := range attempt.RequestHeaders {
		req.Header.Set(key, value)
	}

	deliveredAt := s.now()
	resp, err := s.client.Do(req)
	attempt.DeliveredAt = &deliveredAt
	if err != nil {
		attempt.Status = StatusFailed
		attempt.Error = err.Error()
		return s.recordRetry(ctx, endpoint, attempt)
	}
	defer resp.Body.Close()

	attempt.ResponseStatus = resp.StatusCode
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	attempt.ResponseBody = string(body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		attempt.Status = StatusSucceeded
		return s.createDeliveryAttempt(ctx, attempt)
	}
	attempt.Status = StatusFailed
	attempt.Error = resp.Status
	return s.recordRetry(ctx, endpoint, attempt)
}

func (s *Service) recordRetry(ctx context.Context, endpoint Endpoint, attempt DeliveryAttempt) (DeliveryAttempt, error) {
	maxAttempts := endpoint.RetryMaxAttempts
	if maxAttempts == 0 {
		maxAttempts = 5
	}
	if attempt.AttemptNumber >= maxAttempts {
		attempt.Status = StatusAbandoned
		return s.createDeliveryAttempt(ctx, attempt)
	}

	nextRetry := attempt.ScheduledAt.Add(DefaultBackoff(endpoint, attempt.AttemptNumber))
	attempt.NextRetryAt = &nextRetry
	created, err := s.createDeliveryAttempt(ctx, attempt)
	if err != nil {
		return DeliveryAttempt{}, err
	}
	retry := DeliveryAttempt{
		ID:             id("delatt"),
		Object:         ObjectDeliveryAttempt,
		EventID:        attempt.EventID,
		EndpointID:     attempt.EndpointID,
		AttemptNumber:  attempt.AttemptNumber + 1,
		Status:         StatusScheduled,
		ScheduledAt:    nextRetry,
		RequestURL:     attempt.RequestURL,
		RequestHeaders: attempt.RequestHeaders,
		RequestBody:    attempt.RequestBody,
		Metadata: map[string]string{
			"retry_for": attempt.ID,
		},
		CreatedAt: s.now(),
	}
	if _, err := s.createDeliveryAttempt(ctx, retry); err != nil {
		return DeliveryAttempt{}, err
	}
	return created, nil
}

func (s *Service) createDeliveryAttempt(ctx context.Context, attempt DeliveryAttempt) (DeliveryAttempt, error) {
	stored := attempt
	if !s.storeRawPayloads {
		stored.RequestBody = nil
	}
	return s.repo.CreateDeliveryAttempt(ctx, stored)
}

func (s *Service) audit(ctx context.Context, action, targetType, targetID string, metadata map[string]string) error {
	_, err := s.repo.CreateAuditEntry(ctx, AuditEntry{
		ID:         id("audit"),
		Object:     ObjectAuditEntry,
		Action:     action,
		Actor:      "system",
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   metadata,
		CreatedAt:  s.now(),
	})
	return err
}

func EndpointMatches(endpoint Endpoint, eventType string) bool {
	if !endpoint.Active || endpoint.DeletedAt != nil {
		return false
	}
	if len(endpoint.EnabledEvents) == 0 {
		return true
	}
	for _, enabled := range endpoint.EnabledEvents {
		enabled = strings.TrimSpace(enabled)
		if enabled == "*" || enabled == eventType {
			return true
		}
		if strings.HasSuffix(enabled, ".*") && strings.HasPrefix(eventType, strings.TrimSuffix(enabled, "*")) {
			return true
		}
	}
	return false
}

func DefaultBackoff(endpoint Endpoint, attemptNumber int) time.Duration {
	if attemptNumber > 0 && attemptNumber <= len(endpoint.RetryBackoff) {
		if parsed, err := time.ParseDuration(endpoint.RetryBackoff[attemptNumber-1]); err == nil {
			return parsed
		}
	}
	switch attemptNumber {
	case 1:
		return 10 * time.Second
	case 2:
		return 30 * time.Second
	case 3:
		return 2 * time.Minute
	default:
		return 10 * time.Minute
	}
}

func hasDeliveryOverride(opts DeliveryOptions) bool {
	return opts.Duplicate > 1 || opts.Delay > 0 || opts.OutOfOrder || hasSimulatedDeliveryResult(opts) || opts.SignatureMismatch || opts.SimulateAppFailure != nil
}

func deliveryMetadata(opts DeliveryOptions) map[string]string {
	metadata := map[string]string{}
	if opts.Replay {
		metadata["source"] = SourceReplay
	}
	if opts.Historical {
		metadata["historical"] = "true"
	}
	if opts.Duplicate > 1 {
		metadata["duplicate"] = fmt.Sprintf("%d", opts.Duplicate)
	}
	if opts.Delay > 0 {
		metadata["delay"] = opts.Delay.String()
	}
	if opts.OutOfOrder {
		metadata["out_of_order"] = "true"
	}
	if opts.ResponseStatus != 0 {
		metadata["response_status"] = fmt.Sprintf("%d", opts.ResponseStatus)
	}
	if opts.SimulatedTimeout {
		metadata["timeout"] = "true"
	}
	if opts.SimulatedError != "" {
		metadata["error"] = opts.SimulatedError
	}
	if opts.SignatureMismatch {
		metadata["signature_mismatch"] = "true"
	}
	if opts.SimulateAppFailure != nil {
		metadata["simulate_app_failure"] = "true"
		metadata["response_status"] = fmt.Sprintf("%d", opts.SimulateAppFailure.Status)
		metadata["fail_first_n_attempts"] = fmt.Sprintf("%d", opts.SimulateAppFailure.FailFirstNAttempts)
	}
	return metadata
}

func timeString(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339Nano)
}

func id(prefix string) string {
	return prefix + "_" + randomHex(8)
}

func randomHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}
	return hex.EncodeToString(buf)
}
