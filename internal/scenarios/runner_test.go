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

func TestRunnerInvoiceRetryCanDeclineThenSucceedAtSameClock(t *testing.T) {
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
name: invoice-retry-same-clock
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
  - id: decline-retry
    action: invoice.retry
    params:
      invoiceRef: complete-checkout.invoice.id
      payment_method: pm_card_visa_chargeDeclined
  - id: successful-retry
    action: invoice.retry
    params:
      invoiceRef: complete-checkout.invoice.id
      payment_method: pm_card_visa
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	declined := report.Steps[3].Output["invoice"].(billing.Invoice)
	paid := report.Steps[4].Output["invoice"].(billing.Invoice)
	if declined.Status != "open" || paid.Status != "paid" || paid.AttemptCount != declined.AttemptCount+1 {
		t.Fatalf("declined=%#v paid=%#v, want same-clock decline then paid retry", declined, paid)
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

	report, err := NewRunner(billing.NewService(store), webhooks.NewService(store)).Run(ctx, mustLoad(t, `
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
  - id: schedule-cancel
    action: subscription.update
    params:
      subscriptionRef: advance-renewal.billing.renewals.0.subscription.id
      cancel_at_period_end: true
  - id: advance-cancel
    action: clock.advance
    params:
      duration: 60d
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
	events, ok := report.Steps[3].Output["events"].([]webhooks.Event)
	if !ok {
		t.Fatalf("advance events = %#v, want webhook events", report.Steps[3].Output["events"])
	}
	seen := map[string]bool{}
	for _, event := range events {
		seen[event.Type] = true
	}
	for _, eventType := range []string{"invoice.created", "invoice.finalized", "payment_intent.created", "payment_intent.succeeded", "invoice.payment_succeeded", "invoice.paid", "customer.subscription.updated"} {
		if !seen[eventType] {
			t.Fatalf("advance events missing %s: %#v", eventType, events)
		}
	}

	scheduled := report.Steps[4].Output["subscription"].(billing.Subscription)
	if !scheduled.CancelAtPeriodEnd {
		t.Fatalf("scheduled subscription = %#v, want cancel_at_period_end", scheduled)
	}
	boundary := scheduled.CurrentPeriodEnd
	result := report.Steps[5].Output["billing"].(billing.ClockAdvanceResult)
	if result.CanceledCount != 1 || len(result.Canceled) != 1 || result.Canceled[0].Status != "canceled" {
		t.Fatalf("cancel advance = %#v, want canceled subscription", result)
	}
	if result.Canceled[0].CanceledAt == nil || !result.Canceled[0].CanceledAt.Equal(boundary) {
		t.Fatalf("cancel advance canceled_at = %v, want boundary %v", result.Canceled[0].CanceledAt, boundary)
	}
}

func TestRunnerCheckoutCancelAndSubscriptionLifecycleActions(t *testing.T) {
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

	report, err := NewRunner(billing.NewService(store), webhooks.NewService(store)).Run(ctx, mustLoad(t, `
name: lifecycle-actions
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
      email: lifecycle@example.test
  - id: checkout-to-cancel
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly
  - id: cancel-checkout
    action: checkout.cancel
    params:
      sessionRef: checkout-to-cancel.session.id
  - id: checkout-to-complete
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly
  - id: complete-checkout
    action: checkout.complete
    params:
      sessionRef: checkout-to-complete.session.id
      outcome: payment_succeeded
  - id: schedule-cancel
    action: subscription.cancel
    params:
      subscriptionRef: complete-checkout.subscription.id
      mode: period
  - id: resume-subscription
    action: subscription.resume
    params:
      subscriptionRef: schedule-cancel.subscription.id
  - id: immediate-cancel
    action: subscription.cancel
    params:
      subscriptionRef: resume-subscription.subscription.id
      mode: immediate
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}

	canceledCheckout := report.Steps[2].Output
	session := canceledCheckout["session"].(billing.CheckoutSession)
	if session.Status != "expired" || session.PaymentStatus != "unpaid" {
		t.Fatalf("canceled checkout session = %#v, want expired unpaid", session)
	}
	canceledInvoice := canceledCheckout["invoice"].(billing.Invoice)
	if canceledInvoice.Status != "void" {
		t.Fatalf("canceled checkout invoice = %#v, want void", canceledInvoice)
	}
	canceledIntent := canceledCheckout["payment_intent"].(billing.PaymentIntent)
	if canceledIntent.Status != "canceled" {
		t.Fatalf("canceled checkout intent = %#v, want canceled", canceledIntent)
	}
	checkoutEvents, ok := canceledCheckout["events"].([]webhooks.Event)
	if !ok {
		t.Fatalf("canceled checkout events = %#v, want webhook events", canceledCheckout["events"])
	}
	seenCheckoutExpired := false
	for _, event := range checkoutEvents {
		if event.Type == "checkout.session.expired" {
			seenCheckoutExpired = true
		}
	}
	if !seenCheckoutExpired {
		t.Fatalf("canceled checkout events = %#v, want checkout.session.expired", checkoutEvents)
	}

	scheduled := report.Steps[5].Output["subscription"].(billing.Subscription)
	if !scheduled.CancelAtPeriodEnd || scheduled.Status != "active" {
		t.Fatalf("scheduled subscription = %#v, want active cancel_at_period_end", scheduled)
	}
	resumed := report.Steps[6].Output["subscription"].(billing.Subscription)
	if resumed.CancelAtPeriodEnd || resumed.Status != "active" || resumed.CanceledAt != nil {
		t.Fatalf("resumed subscription = %#v, want active without cancellation", resumed)
	}
	immediate := report.Steps[7].Output["subscription"].(billing.Subscription)
	if immediate.Status != "canceled" || immediate.CancelAtPeriodEnd {
		t.Fatalf("immediate cancel subscription = %#v, want canceled immediately", immediate)
	}
	events, ok := report.Steps[7].Output["events"].([]webhooks.Event)
	if !ok || len(events) != 1 || events[0].Type != "customer.subscription.deleted" {
		t.Fatalf("immediate cancel events = %#v, want customer.subscription.deleted", report.Steps[7].Output["events"])
	}
}

func TestRunnerInvoiceFailPaymentActionMutatesOpenInvoice(t *testing.T) {
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

	report, err := NewRunner(billing.NewService(store), webhooks.NewService(store)).Run(ctx, mustLoad(t, `
name: invoice-fail-payment
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
      email: fail@example.test
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
  - id: fail-payment
    action: invoice.fail_payment
    params:
      invoiceRef: complete-checkout.invoice.id
      payment_method: pm_card_visa_chargeDeclinedInsufficientFunds
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}

	output := report.Steps[3].Output
	if output["failure_simulation"] != true {
		t.Fatalf("fail output = %#v, want failure_simulation", output)
	}
	invoice := output["invoice"].(billing.Invoice)
	if invoice.Status != "open" || invoice.NextPaymentAttempt == nil {
		t.Fatalf("invoice = %#v, want open invoice with next payment attempt", invoice)
	}
	subscription := output["subscription"].(billing.Subscription)
	if subscription.Status != "past_due" {
		t.Fatalf("subscription = %#v, want past_due", subscription)
	}
	intent := output["payment_intent"].(billing.PaymentIntent)
	if intent.Status != "requires_payment_method" || intent.DeclineCode != "insufficient_funds" {
		t.Fatalf("payment intent = %#v, want insufficient funds failure", intent)
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

func TestRunnerWebhookDeliveryConvenienceActionsUseGenericReplay(t *testing.T) {
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
name: webhook-delivery-convenience-actions
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
      email: webhook@example.test
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
  - id: duplicate-delivery
    action: webhook.deliver_duplicate
    params:
      eventRef: complete-checkout.events.0.id
  - id: out-of-order-delivery
    action: webhook.deliver_out_of_order
    params:
      eventRef: complete-checkout.events.0.id
`))
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}

	duplicateAttempts, ok := report.Steps[3].Output["delivery_attempts"].([]webhooks.DeliveryAttempt)
	if !ok || len(duplicateAttempts) != 2 {
		t.Fatalf("duplicate output = %#v, want two replay attempts", report.Steps[3].Output)
	}
	if duplicateAttempts[1].Metadata["duplicate"] != "true" || duplicateAttempts[1].Metadata["source"] != webhooks.SourceReplay {
		t.Fatalf("duplicate attempt metadata = %#v, want duplicate replay evidence", duplicateAttempts[1].Metadata)
	}

	outOfOrderAttempts, ok := report.Steps[4].Output["delivery_attempts"].([]webhooks.DeliveryAttempt)
	if !ok || len(outOfOrderAttempts) != 1 {
		t.Fatalf("out-of-order output = %#v, want one replay attempt", report.Steps[4].Output)
	}
	if outOfOrderAttempts[0].Metadata["out_of_order"] != "true" || outOfOrderAttempts[0].Metadata["source"] != webhooks.SourceReplay {
		t.Fatalf("out-of-order attempt metadata = %#v, want out-of-order replay evidence", outOfOrderAttempts[0].Metadata)
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
