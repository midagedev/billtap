package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/webhooks"
)

func TestWebhookPersistenceAndDeliveryAttempt(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	endpoint, err := store.CreateWebhookEndpoint(ctx, webhooks.Endpoint{
		ID:            "we_test",
		Object:        webhooks.ObjectEndpoint,
		URL:           "http://example.test/webhook",
		Secret:        "webhook_secret_test",
		EnabledEvents: []string{"invoice.*"},
		Active:        true,
		CreatedAt:     mustTime("2026-05-08T00:00:00Z"),
		UpdatedAt:     mustTime("2026-05-08T00:00:00Z"),
	})
	if err != nil {
		t.Fatalf("CreateWebhookEndpoint returned error: %v", err)
	}
	if !webhooks.EndpointMatches(endpoint, "invoice.payment_succeeded") {
		t.Fatal("endpoint should match invoice wildcard")
	}

	event, err := store.CreateEvent(ctx, webhooks.Event{
		ID:              "evt_test",
		Object:          webhooks.ObjectEvent,
		Type:            "invoice.payment_succeeded",
		Created:         1778233312,
		APIVersion:      webhooks.DefaultAPIVersion,
		PendingWebhooks: 1,
		Data:            webhooks.EventData{Object: json.RawMessage(`{"id":"in_test"}`)},
		Billtap:         webhooks.EventMetadata{Source: webhooks.SourceCheckout, Sequence: 1},
		RawPayload:      json.RawMessage(`{"id":"evt_test"}`),
		CreatedAt:       mustTime("2026-05-08T00:00:01Z"),
	})
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}
	if event.Type != "invoice.payment_succeeded" {
		t.Fatalf("event type = %q", event.Type)
	}

	_, err = store.CreateDeliveryAttempt(ctx, webhooks.DeliveryAttempt{
		ID:             "delatt_test",
		Object:         webhooks.ObjectDeliveryAttempt,
		EventID:        event.ID,
		EndpointID:     endpoint.ID,
		AttemptNumber:  1,
		Status:         webhooks.StatusScheduled,
		ScheduledAt:    mustTime("2026-05-08T00:00:02Z"),
		RequestURL:     endpoint.URL,
		RequestHeaders: map[string]string{webhooks.SignatureHeaderName: "t=1,v1=abc"},
		RequestBody:    event.RawPayload,
		CreatedAt:      mustTime("2026-05-08T00:00:02Z"),
	})
	if err != nil {
		t.Fatalf("CreateDeliveryAttempt returned error: %v", err)
	}
	attempts, err := store.ListDeliveryAttempts(ctx, webhooks.DeliveryAttemptFilter{EventID: event.ID})
	if err != nil {
		t.Fatalf("ListDeliveryAttempts returned error: %v", err)
	}
	if len(attempts) != 1 || attempts[0].RequestHeaders[webhooks.SignatureHeaderName] == "" {
		t.Fatalf("attempts = %#v, want signed persisted attempt", attempts)
	}
}

func mustTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestWebhookServiceRecordsSignedDelivery(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	var gotSignature string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSignature = r.Header.Get(webhooks.SignatureHeaderName)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := webhooks.NewService(store)
	if _, err := service.CreateEndpoint(ctx, webhooks.Endpoint{URL: server.URL, EnabledEvents: []string{"checkout.session.completed"}}); err != nil {
		t.Fatalf("CreateEndpoint returned error: %v", err)
	}
	_, attempts, err := service.CreateEvent(ctx, webhooks.EventInput{
		Type:          "checkout.session.completed",
		ObjectPayload: json.RawMessage(`{"id":"cs_test"}`),
		Source:        webhooks.SourceCheckout,
		Sequence:      1,
	})
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}
	if len(attempts) != 1 || attempts[0].Status != webhooks.StatusSucceeded {
		t.Fatalf("attempts = %#v, want one successful delivery", attempts)
	}
	if gotSignature == "" {
		t.Fatal("test server did not receive Billtap-Signature")
	}
}

func TestWebhookServiceAsyncDeliveryDoesNotWaitForReceiver(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	var received atomic.Int32
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.Add(1)
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	defer close(release)

	service := webhooks.NewService(store)
	if _, err := service.CreateEndpoint(ctx, webhooks.Endpoint{URL: server.URL, EnabledEvents: []string{"checkout.session.completed"}}); err != nil {
		t.Fatalf("CreateEndpoint returned error: %v", err)
	}

	start := time.Now()
	_, attempts, err := service.CreateEvent(ctx, webhooks.EventInput{
		Type:          "checkout.session.completed",
		ObjectPayload: json.RawMessage(`{"id":"cs_async"}`),
		Source:        webhooks.SourceCheckout,
		Sequence:      1,
		AsyncDelivery: true,
	})
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Fatalf("CreateEvent took %s, want async return before receiver is released", elapsed)
	}
	if len(attempts) != 0 {
		t.Fatalf("attempts = %#v, want none returned for async delivery", attempts)
	}

	deadline := time.Now().Add(2 * time.Second)
	for received.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if received.Load() == 0 {
		t.Fatal("async delivery did not reach receiver")
	}
}

func TestWebhookServiceLimitsConcurrentDelivery(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	var received atomic.Int32
	var inFlight atomic.Int32
	var maxInFlight atomic.Int32
	release := make(chan struct{})
	var released atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := inFlight.Add(1)
		for {
			max := maxInFlight.Load()
			if current <= max || maxInFlight.CompareAndSwap(max, current) {
				break
			}
		}
		received.Add(1)
		<-release
		inFlight.Add(-1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	defer func() {
		if released.CompareAndSwap(false, true) {
			close(release)
		}
	}()

	service := webhooks.NewServiceWithOptions(store, webhooks.ServiceOptions{
		StoreRawPayloads:    true,
		RetentionDays:       30,
		DeliveryConcurrency: 1,
	})
	if _, err := service.CreateEndpoint(ctx, webhooks.Endpoint{URL: server.URL, EnabledEvents: []string{"checkout.session.completed"}}); err != nil {
		t.Fatalf("CreateEndpoint returned error: %v", err)
	}

	done := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func(sequence int64) {
			_, _, err := service.CreateEvent(ctx, webhooks.EventInput{
				Type:          "checkout.session.completed",
				ObjectPayload: json.RawMessage(`{"id":"cs_limit"}`),
				Source:        webhooks.SourceCheckout,
				Sequence:      sequence,
			})
			done <- err
		}(int64(i + 1))
	}

	deadline := time.Now().Add(2 * time.Second)
	for received.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if received.Load() == 0 {
		t.Fatal("first delivery did not start")
	}
	time.Sleep(50 * time.Millisecond)
	if maxInFlight.Load() != 1 {
		t.Fatalf("max in-flight delivery = %d, want 1", maxInFlight.Load())
	}
	if released.CompareAndSwap(false, true) {
		close(release)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("CreateEvent returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("first CreateEvent did not complete after release")
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("CreateEvent returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("second CreateEvent did not complete after release")
	}
	if received.Load() != 2 {
		t.Fatalf("received deliveries = %d, want 2", received.Load())
	}
}

func TestWebhookServiceCanUseStripeSignatureHeader(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	var gotStripeSignature string
	var gotAPIVersion string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotStripeSignature = r.Header.Get(webhooks.StripeSignatureHeaderName)
		if r.Header.Get(webhooks.SignatureHeaderName) != "" {
			t.Fatalf("received legacy signature header in Stripe mode")
		}
		var payload struct {
			APIVersion string `json:"api_version"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode webhook payload: %v", err)
		}
		gotAPIVersion = payload.APIVersion
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := webhooks.NewServiceWithOptions(store, webhooks.ServiceOptions{
		StoreRawPayloads:    true,
		SignatureHeaderName: webhooks.StripeSignatureHeaderName,
		APIVersion:          "2025-03-31.basil",
	})
	if _, err := service.CreateEndpoint(ctx, webhooks.Endpoint{URL: server.URL, EnabledEvents: []string{"checkout.session.completed"}}); err != nil {
		t.Fatalf("CreateEndpoint returned error: %v", err)
	}
	_, attempts, err := service.CreateEvent(ctx, webhooks.EventInput{
		Type:          "checkout.session.completed",
		ObjectPayload: json.RawMessage(`{"id":"cs_test"}`),
		Source:        webhooks.SourceCheckout,
		Sequence:      1,
	})
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}
	if len(attempts) != 1 || attempts[0].RequestHeaders[webhooks.StripeSignatureHeaderName] == "" {
		t.Fatalf("attempts = %#v, want Stripe-Signature evidence", attempts)
	}
	if gotStripeSignature == "" {
		t.Fatal("test server did not receive Stripe-Signature")
	}
	if gotAPIVersion != "2025-03-31.basil" {
		t.Fatalf("api_version = %q, want configured Stripe API version", gotAPIVersion)
	}
	if webhooks.SignatureHeaderValue(attempts[0].RequestHeaders) == "" {
		t.Fatal("SignatureHeaderValue did not find configured Stripe-Signature")
	}
}

func TestRelayModeDoesNotPersistRawPayloadsAndAuditsOverrides(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode received webhook: %v", err)
		}
		receivedBody = body["id"].(string)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := webhooks.NewServiceWithOptions(store, webhooks.ServiceOptions{StoreRawPayloads: false, RetentionDays: 30})
	if _, err := service.CreateEndpoint(ctx, webhooks.Endpoint{URL: server.URL, EnabledEvents: []string{"checkout.session.completed"}}); err != nil {
		t.Fatalf("CreateEndpoint returned error: %v", err)
	}
	event, _, err := service.CreateEvent(ctx, webhooks.EventInput{
		Type:          "checkout.session.completed",
		ObjectPayload: json.RawMessage(`{"id":"cs_relay"}`),
		Source:        webhooks.SourceCheckout,
		Sequence:      1,
	})
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}
	if receivedBody == "" {
		t.Fatal("webhook receiver did not get signed payload")
	}
	persistedEvent, err := store.GetEvent(ctx, event.ID)
	if err != nil {
		t.Fatalf("GetEvent returned error: %v", err)
	}
	if string(persistedEvent.RawPayload) != "{}" {
		t.Fatalf("persisted raw payload = %s, want metadata-only payload", persistedEvent.RawPayload)
	}
	attempts, err := store.ListDeliveryAttempts(ctx, webhooks.DeliveryAttemptFilter{EventID: event.ID})
	if err != nil {
		t.Fatalf("ListDeliveryAttempts returned error: %v", err)
	}
	if len(attempts) != 1 || len(attempts[0].RequestBody) != 0 {
		t.Fatalf("attempts = %#v, want request body omitted in relay mode", attempts)
	}

	if _, err := service.ReplayEvent(ctx, event.ID, webhooks.ReplayOptions{Duplicate: 2, Delay: time.Minute, OutOfOrder: true}); err != nil {
		t.Fatalf("ReplayEvent returned error: %v", err)
	}
	replayAudit, err := store.ListAuditEntries(ctx, webhooks.AuditFilter{Action: "webhook.replay", TargetID: event.ID})
	if err != nil {
		t.Fatalf("ListAuditEntries returned error: %v", err)
	}
	if len(replayAudit) != 1 || replayAudit[0].Metadata["out_of_order"] != "true" {
		t.Fatalf("replay audit = %#v, want replay override evidence", replayAudit)
	}

	overrideEvent, _, err := service.CreateEvent(ctx, webhooks.EventInput{
		Type:          "checkout.session.completed",
		ObjectPayload: json.RawMessage(`{"id":"cs_override"}`),
		Source:        webhooks.SourceScenario,
		Sequence:      2,
		DeliveryOptions: webhooks.DeliveryOptions{
			Duplicate:  2,
			Delay:      time.Minute,
			OutOfOrder: true,
		},
	})
	if err != nil {
		t.Fatalf("CreateEvent with override returned error: %v", err)
	}
	overrideAudit, err := store.ListAuditEntries(ctx, webhooks.AuditFilter{Action: "webhook.delivery_override", TargetID: overrideEvent.ID})
	if err != nil {
		t.Fatalf("ListAuditEntries returned error: %v", err)
	}
	if len(overrideAudit) != 1 || overrideAudit[0].Metadata["duplicate"] != "2" {
		t.Fatalf("override audit = %#v, want delivery override evidence", overrideAudit)
	}
}

func TestRetentionRedactsOldWebhookEvidence(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	old := mustTime("2026-05-01T00:00:00Z")
	event, err := store.CreateEvent(ctx, webhooks.Event{
		ID:              "evt_old",
		Object:          webhooks.ObjectEvent,
		Type:            "invoice.paid",
		Created:         old.Unix(),
		APIVersion:      webhooks.DefaultAPIVersion,
		PendingWebhooks: 1,
		Data:            webhooks.EventData{Object: json.RawMessage(`{"id":"in_old"}`)},
		Billtap:         webhooks.EventMetadata{Source: webhooks.SourceCheckout, Sequence: 1},
		RawPayload:      json.RawMessage(`{"id":"evt_old","client_secret":"secret"}`),
		CreatedAt:       old,
	})
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}
	if _, err := store.CreateWebhookEndpoint(ctx, webhooks.Endpoint{
		ID:        "we_old",
		Object:    webhooks.ObjectEndpoint,
		URL:       "http://example.test/webhook",
		Secret:    "webhook_secret_test",
		Active:    true,
		CreatedAt: old,
		UpdatedAt: old,
	}); err != nil {
		t.Fatalf("CreateWebhookEndpoint returned error: %v", err)
	}
	if _, err := store.CreateDeliveryAttempt(ctx, webhooks.DeliveryAttempt{
		ID:             "delatt_old",
		Object:         webhooks.ObjectDeliveryAttempt,
		EventID:        event.ID,
		EndpointID:     "we_old",
		AttemptNumber:  1,
		Status:         webhooks.StatusFailed,
		ScheduledAt:    old,
		RequestURL:     "http://example.test/webhook",
		RequestHeaders: map[string]string{webhooks.SignatureHeaderName: "t=1,v1=abc"},
		RequestBody:    json.RawMessage(`{"id":"evt_old"}`),
		ResponseBody:   `{"client_secret":"secret"}`,
		CreatedAt:      old,
	}); err != nil {
		t.Fatalf("CreateDeliveryAttempt returned error: %v", err)
	}

	result, err := store.ApplyRetention(ctx, mustTime("2026-05-08T00:00:00Z"))
	if err != nil {
		t.Fatalf("ApplyRetention returned error: %v", err)
	}
	if result.WebhookEventsRedacted != 1 || result.DeliveryAttemptsRedacted != 1 {
		t.Fatalf("retention result = %#v, want event and attempt redacted", result)
	}
	persistedEvent, _ := store.GetEvent(ctx, event.ID)
	if string(persistedEvent.RawPayload) != "{}" {
		t.Fatalf("raw payload = %s, want redacted", persistedEvent.RawPayload)
	}
	attempts, _ := store.ListDeliveryAttempts(ctx, webhooks.DeliveryAttemptFilter{EventID: event.ID})
	if len(attempts) != 1 || len(attempts[0].RequestBody) != 0 || attempts[0].ResponseBody != "" {
		t.Fatalf("attempts = %#v, want request and response bodies redacted", attempts)
	}
}
