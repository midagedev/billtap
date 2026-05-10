package scenarios

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
)

func TestRunnerResolvesReferencesAndCheckoutAliases(t *testing.T) {
	runner := NewRunner(newTestBilling(t), nil)
	scenario := mustLoad(t, `
name: checkout-aliases
catalog:
  products:
    - id: prod_pro
      name: Pro
  prices:
    - id: price_pro_monthly
      product: prod_pro
      currency: usd
      unitAmount: 4900
      interval: month
clock:
  start: "2026-05-08T00:00:00Z"
steps:
  - id: create-customer
    action: customer.create
    params:
      email: user@example.test
  - id: checkout
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly
  - id: complete-checkout
    action: checkout.complete
    params:
      sessionRef: checkout.session.id
      outcome: payment_failed
  - id: retry-payment
    action: invoice.retry
    params:
      subscriptionRef: checkout.subscription.id
      invoiceRef: complete-checkout.invoice.id
      payment_method: pm_card_visa
`)
	report, err := runner.Run(context.Background(), scenario)
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	if report.ExitCode() != ExitPass {
		t.Fatalf("ExitCode = %d, want pass", report.ExitCode())
	}
	if got := len(report.Steps); got != 4 {
		t.Fatalf("steps = %d, want 4", got)
	}
	retry := report.Steps[3].Output
	if retry["subscription"] == "" || retry["invoice"] == "" {
		t.Fatalf("retry output = %#v, want resolved subscription and invoice refs", retry)
	}
}

func TestRunnerInvoiceRetryMutatesBillingState(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})

	report, err := NewRunner(billing.NewService(store), nil).Run(ctx, mustLoad(t, `
name: invoice-retry-mutates
clock:
  start: "2026-05-08T00:00:00Z"
catalog:
  products:
    - id: prod_pro
      name: Pro
  prices:
    - id: price_pro_monthly
      product: prod_pro
      currency: usd
      unitAmount: 4900
      interval: month
steps:
  - id: create-customer
    action: customer.create
    params:
      email: retry@example.test
  - id: checkout
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly
  - id: complete-checkout
    action: checkout.complete
    params:
      sessionRef: checkout.session.id
      outcome: payment_failed
  - id: retry-payment
    action: invoice.retry
    params:
      invoiceRef: complete-checkout.invoice.id
      payment_method: pm_card_visa
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	retry := report.Steps[3].Output
	invoice, ok := retry["invoice"].(billing.Invoice)
	if !ok || invoice.Status != "paid" || invoice.AmountPaid != 4900 || invoice.NextPaymentAttempt != nil {
		t.Fatalf("retry invoice = %#v, want paid invoice", retry["invoice"])
	}
	subscription, ok := retry["subscription"].(billing.Subscription)
	if !ok || subscription.Status != "active" {
		t.Fatalf("retry subscription = %#v, want active", retry["subscription"])
	}
	intent, ok := retry["payment_intent"].(billing.PaymentIntent)
	if !ok || intent.Status != "succeeded" {
		t.Fatalf("retry payment intent = %#v, want succeeded", retry["payment_intent"])
	}
}

func TestRunnerInvoiceRetryWithoutInvoiceKeepsEvidenceFallback(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})

	report, err := NewRunner(billing.NewService(store), nil).Run(ctx, mustLoad(t, `
name: invoice-retry-evidence
steps:
  - id: retry-payment
    action: invoice.retry
    params:
      subscriptionRef: sub_profile_123
      outcome: payment_succeeded
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	retry := report.Steps[0].Output
	if retry["subscription"] != "sub_profile_123" || retry["deterministic"] != true {
		t.Fatalf("retry output = %#v, want deterministic fallback evidence", retry)
	}
	if !strings.Contains(retry["note"].(string), "no billing invoice") {
		t.Fatalf("retry note = %#v, want missing invoice note", retry["note"])
	}
}

func TestRunnerClockAdvanceRenewsAndCancelsSubscriptions(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})

	report, err := NewRunner(billing.NewService(store), nil).Run(ctx, mustLoad(t, `
name: clock-renewal-cancel
clock:
  start: "2026-05-08T00:00:00Z"
catalog:
  products:
    - id: prod_pro
      name: Pro
  prices:
    - id: price_pro_monthly
      product: prod_pro
      currency: usd
      unitAmount: 4900
      interval: month
steps:
  - id: create-customer
    action: customer.create
    params:
      email: renew@example.test
  - id: checkout
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly
  - id: complete-checkout
    action: checkout.complete
    params:
      sessionRef: checkout.session.id
      outcome: payment_succeeded
  - id: advance-renewal
    action: clock.advance
    params:
      duration: 31d
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	advance := report.Steps[3].Output["billing"].(billing.ClockAdvanceResult)
	if advance.Renewed != 1 || len(advance.Renewals) != 1 {
		t.Fatalf("advance = %#v, want one renewal", advance)
	}
	if advance.Renewals[0].Invoice.Status != "paid" || advance.Renewals[0].Subscription.LatestInvoiceID != advance.Renewals[0].Invoice.ID {
		t.Fatalf("renewal = %#v, want paid invoice and latest invoice update", advance.Renewals[0])
	}

	sub := advance.Renewals[0].Subscription
	sub.CancelAtPeriodEnd = true
	canceledAt := time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC)
	sub.CanceledAt = &canceledAt
	if _, err := store.UpdateSubscription(ctx, sub, nil); err != nil {
		t.Fatalf("UpdateSubscription: %v", err)
	}
	result, err := billing.NewService(store).AdvanceClock(ctx, sub.CurrentPeriodEnd)
	if err != nil {
		t.Fatalf("AdvanceClock cancel: %v", err)
	}
	if result.CanceledCount != 1 || len(result.Canceled) != 1 || result.Canceled[0].Status != "canceled" {
		t.Fatalf("cancel advance = %#v, want canceled subscription", result)
	}
}

func TestRunnerAppAssertPass(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode assertion payload: %v", err)
		}
		if payload["target"] != "workspace.subscription" || payload["expected"] == nil || payload["context"] == nil || payload["clock"] == nil {
			t.Fatalf("payload = %#v, want assertion context", payload)
		}
		_, _ = w.Write([]byte(`{"pass":true}`))
	}))
	defer server.Close()

	report, err := NewRunner(nil, nil).Run(context.Background(), mustLoad(t, `
name: assertion-pass
app:
  assertions:
    baseUrl: `+server.URL+`/assertions
steps:
  - id: assert-subscription
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`))
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if gotPath != "/assertions/workspace/subscription" {
		t.Fatalf("callback path = %q, want slash-converted target", gotPath)
	}
	if report.ExitCode() != ExitPass {
		t.Fatalf("ExitCode = %d, want 0", report.ExitCode())
	}
}

func TestRunnerAppAssertFailureExitCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"pass":false,"message":"status was past_due"}`))
	}))
	defer server.Close()

	report, err := NewRunner(nil, nil).Run(context.Background(), mustLoad(t, `
name: assertion-fail
app:
  assertions:
    baseUrl: `+server.URL+`
steps:
  - id: assert-subscription
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`))
	if !errors.Is(err, ErrAssertionFailed) {
		t.Fatalf("Run error = %v, want ErrAssertionFailed", err)
	}
	if report.ExitCode() != ExitAssertionFailed {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitAssertionFailed)
	}
	if !strings.Contains(report.Markdown(), "status was past_due") {
		t.Fatalf("Markdown missing assertion detail:\n%s", report.Markdown())
	}
}

func TestRunnerWebhookReplayFailureSimulation(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})

	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()

	webhookService := webhooks.NewService(store)
	if _, err := webhookService.CreateEndpoint(ctx, webhooks.Endpoint{
		URL:           receiver.URL,
		EnabledEvents: []string{"checkout.session.completed"},
	}); err != nil {
		t.Fatalf("CreateEndpoint: %v", err)
	}

	report, err := NewRunner(billing.NewService(store), webhookService).Run(ctx, mustLoad(t, `
name: webhook-replay-failure
catalog:
  products:
    - id: prod_pro
      name: Pro
  prices:
    - id: price_pro_monthly
      product: prod_pro
      currency: usd
      unitAmount: 4900
      interval: month
steps:
  - id: create-customer
    action: customer.create
    params:
      email: user@example.test
  - id: checkout
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly
  - id: complete-checkout
    action: checkout.complete
    params:
      sessionRef: checkout.session.id
      outcome: payment_succeeded
  - id: replay-webhook
    action: webhook.replay
    params:
      eventRef: complete-checkout.events.0.id
      duplicate: 2
      responseStatus: 500
      responseBody: receiver down
      signatureMismatch: true
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	if report.ExitCode() != ExitPass {
		t.Fatalf("ExitCode = %d, want pass", report.ExitCode())
	}
	attempts, ok := report.Steps[3].Output["delivery_attempts"].([]webhooks.DeliveryAttempt)
	if !ok || len(attempts) != 2 {
		t.Fatalf("replay output = %#v, want two delivery attempts", report.Steps[3].Output)
	}
	for _, attempt := range attempts {
		if attempt.Status != webhooks.StatusFailed || attempt.ResponseStatus != 500 || attempt.Metadata["signature_mismatch"] != "true" {
			t.Fatalf("attempt = %#v, want failed 500 signature mismatch evidence", attempt)
		}
	}
}

func TestRunnerAppCallbackErrorExitCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	report, err := NewRunner(nil, nil).Run(context.Background(), mustLoad(t, `
name: callback-error
app:
  assertions:
    baseUrl: `+server.URL+`
steps:
  - id: assert-subscription
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`))
	if !errors.Is(err, ErrAppCallbackFailure) {
		t.Fatalf("Run error = %v, want ErrAppCallbackFailure", err)
	}
	if report.ExitCode() != ExitAppCallbackFailure {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitAppCallbackFailure)
	}
}

func TestReportJSONMarkdownAndInvalidExitCode(t *testing.T) {
	report := Report{
		Name:        "bad",
		Status:      "failed",
		FailureType: FailureInvalidConfig,
		Errors:      []string{"name is required"},
	}
	body, err := report.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if !strings.Contains(string(body), `"failure_type": "invalid_config"`) {
		t.Fatalf("JSON = %s", body)
	}
	if !strings.Contains(report.Markdown(), "Exit code: `2`") {
		t.Fatalf("Markdown = %s", report.Markdown())
	}
	if report.ExitCode() != ExitInvalidConfig {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitInvalidConfig)
	}
}

func TestReportRuntimeFailureExitCode(t *testing.T) {
	report := Report{Name: "runtime", Status: "failed", FailureType: FailureRunner}
	if report.ExitCode() != ExitRuntimeFailure {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitRuntimeFailure)
	}
	if !strings.Contains(report.Markdown(), "Exit code: `4`") {
		t.Fatalf("Markdown = %s", report.Markdown())
	}
}

func newTestBilling(t *testing.T) *billing.Service {
	t.Helper()
	store, err := storage.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})
	return billing.NewService(store)
}

func mustLoad(t *testing.T, raw string) Scenario {
	t.Helper()
	scenario, err := Load([]byte(raw))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return scenario
}
