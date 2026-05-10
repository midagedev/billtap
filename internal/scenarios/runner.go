package scenarios

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/security"
	"github.com/hckim/billtap/internal/webhooks"
)

type BillingService interface {
	CreateCustomer(context.Context, billing.Customer) (billing.Customer, error)
	CreateProduct(context.Context, billing.Product) (billing.Product, error)
	CreatePrice(context.Context, billing.Price) (billing.Price, error)
	CreateCheckoutSession(context.Context, billing.CheckoutSession) (billing.CheckoutSession, error)
	CompleteCheckout(context.Context, string, string) (billing.CheckoutSession, error)
	CompleteCheckoutAt(context.Context, string, string, time.Time) (billing.CheckoutSession, error)
	GetSubscription(context.Context, string) (billing.Subscription, error)
	GetInvoice(context.Context, string) (billing.Invoice, error)
	GetPaymentIntent(context.Context, string) (billing.PaymentIntent, error)
	PayInvoice(context.Context, string, billing.InvoicePaymentOptions) (billing.InvoicePaymentResult, error)
	AdvanceClock(context.Context, time.Time) (billing.ClockAdvanceResult, error)
}

type WebhookService interface {
	CreateEvent(context.Context, webhooks.EventInput) (webhooks.Event, []webhooks.DeliveryAttempt, error)
	ReplayEvent(context.Context, string, webhooks.ReplayOptions) ([]webhooks.DeliveryAttempt, error)
}

type Runner struct {
	Billing    BillingService
	Webhooks   WebhookService
	HTTPClient *http.Client
}

func NewRunner(billing BillingService, webhooks WebhookService) *Runner {
	return &Runner{
		Billing:  billing,
		Webhooks: webhooks,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (r *Runner) Run(ctx context.Context, scenario Scenario) (Report, error) {
	report := Report{Name: scenario.Name, Status: "failed"}
	if err := Validate(scenario); err != nil {
		report.FailureType = FailureInvalidConfig
		report.Errors = append(report.Errors, err.Error())
		return report, err
	}

	clock, err := NewClock(scenario.Clock.Start)
	if err != nil {
		report.FailureType = FailureInvalidConfig
		report.Errors = append(report.Errors, err.Error())
		return report, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	report.ClockStart = clock.Now()
	report.ClockEnd = clock.Now()
	report.StartedAt = time.Now().UTC()

	state := newRunState(&clock)
	if err := r.seedCatalog(ctx, scenario, state); err != nil {
		report.FailureType = FailureRunner
		report.Errors = append(report.Errors, err.Error())
		report.FinishedAt = time.Now().UTC()
		report.Duration = report.FinishedAt.Sub(report.StartedAt).String()
		report.ClockEnd = clock.Now()
		return report, err
	}

	for _, step := range scenario.Steps {
		stepReport, runErr := r.runStep(ctx, scenario, step, state)
		report.Steps = append(report.Steps, stepReport)
		report.ClockEnd = clock.Now()
		if runErr != nil {
			report.FailureType = failureType(runErr)
			report.Errors = append(report.Errors, runErr.Error())
			report.FinishedAt = time.Now().UTC()
			report.Duration = report.FinishedAt.Sub(report.StartedAt).String()
			return report, runErr
		}
	}

	report.Status = "passed"
	report.FinishedAt = time.Now().UTC()
	report.Duration = report.FinishedAt.Sub(report.StartedAt).String()
	report.ClockEnd = clock.Now()
	return report, nil
}

type runState struct {
	clock      *Clock
	results    map[string]map[string]any
	productIDs map[string]string
	priceIDs   map[string]string
	sequence   int64
	saas       *saasState
}

func newRunState(clock *Clock) *runState {
	return &runState{
		clock:      clock,
		results:    map[string]map[string]any{},
		productIDs: map[string]string{},
		priceIDs:   map[string]string{},
		saas:       newSaaSState(),
	}
}

func (r *Runner) seedCatalog(ctx context.Context, scenario Scenario, state *runState) error {
	for _, product := range scenario.Catalog.Products {
		if r.Billing == nil {
			return errors.New("billing service is required for catalog.products")
		}
		active := true
		if product.Active != nil {
			active = *product.Active
		}
		created, err := r.Billing.CreateProduct(ctx, billing.Product{
			Name:        product.Name,
			Description: product.Description,
			Active:      active,
			Metadata:    product.Metadata,
		})
		if err != nil {
			return fmt.Errorf("seed product %q: %w", product.ID, err)
		}
		state.productIDs[product.ID] = created.ID
		state.results[product.ID] = map[string]any{"product": created}
	}
	for _, price := range scenario.Catalog.Prices {
		if r.Billing == nil {
			return errors.New("billing service is required for catalog.prices")
		}
		productID := state.productIDs[price.Product]
		if productID == "" {
			productID = price.Product
		}
		active := true
		if price.Active != nil {
			active = *price.Active
		}
		created, err := r.Billing.CreatePrice(ctx, billing.Price{
			ProductID:              productID,
			Currency:               price.Currency,
			UnitAmount:             price.UnitAmount,
			RecurringInterval:      price.Interval,
			RecurringIntervalCount: price.IntervalCount,
			Active:                 active,
			Metadata:               price.Metadata,
		})
		if err != nil {
			return fmt.Errorf("seed price %q: %w", price.ID, err)
		}
		state.priceIDs[price.ID] = created.ID
		state.results[price.ID] = map[string]any{"price": created}
	}
	return nil
}

func (r *Runner) runStep(ctx context.Context, scenario Scenario, step Step, state *runState) (StepReport, error) {
	stepReport := StepReport{
		ID:        step.ID,
		Action:    step.Action,
		Status:    "failed",
		Clock:     state.clock.Now(),
		StartedAt: time.Now().UTC(),
	}
	finish := func(output map[string]any, err error) (StepReport, error) {
		stepReport.FinishedAt = time.Now().UTC()
		stepReport.Duration = stepReport.FinishedAt.Sub(stepReport.StartedAt).String()
		if output != nil {
			stepReport.Output = output
		}
		if err != nil {
			stepReport.Error = err.Error()
			return stepReport, err
		}
		stepReport.Status = "passed"
		return stepReport, nil
	}

	params, err := state.resolveParams(step.Params)
	if err != nil {
		return finish(nil, err)
	}

	switch step.Action {
	case "customer.create":
		return finish(r.createCustomer(ctx, step, params, state))
	case "product.create":
		return finish(r.createProduct(ctx, step, params, state))
	case "price.create":
		return finish(r.createPrice(ctx, step, params, state))
	case "checkout.create":
		return finish(r.createCheckout(ctx, step, params, state))
	case "checkout.complete":
		return finish(r.completeCheckout(ctx, scenario, step, params, state))
	case "clock.advance":
		return finish(r.advanceClock(ctx, scenario, step, params, state))
	case "invoice.retry":
		return finish(r.retryInvoice(ctx, scenario, step, params, state))
	case "webhook.replay":
		return finish(r.replayWebhook(ctx, step, params, state))
	case "webhook.deliver_duplicate", "webhook.deliver_out_of_order":
		return finish(r.runSaaSStep(ctx, scenario, step, params, state))
	case "app.assert":
		output, assertion, err := r.assertApp(ctx, scenario, step, params, state)
		stepReport.Assertion = assertion
		return finish(output, err)
	default:
		if strings.HasPrefix(step.Action, "saas.") {
			return finish(r.runSaaSStep(ctx, scenario, step, params, state))
		}
		return finish(nil, fmt.Errorf("%w: unsupported action %q", ErrInvalidConfig, step.Action))
	}
}

func (r *Runner) createCustomer(ctx context.Context, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if r.Billing == nil {
		return nil, errors.New("billing service is required for customer.create")
	}
	customer, err := r.Billing.CreateCustomer(ctx, billing.Customer{
		Email: stringValue(params, "email"),
		Name:  stringValue(params, "name"),
	})
	if err != nil {
		return nil, err
	}
	output := map[string]any{"customer": customer}
	state.results[step.ID] = output
	return output, nil
}

func (r *Runner) createProduct(ctx context.Context, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if r.Billing == nil {
		return nil, errors.New("billing service is required for product.create")
	}
	product, err := r.Billing.CreateProduct(ctx, billing.Product{
		Name:        stringValue(params, "name"),
		Description: stringValue(params, "description"),
		Active:      boolValue(params, "active", true),
	})
	if err != nil {
		return nil, err
	}
	output := map[string]any{"product": product}
	state.results[step.ID] = output
	if alias := stringValue(params, "id"); alias != "" {
		state.productIDs[alias] = product.ID
		state.results[alias] = output
	}
	return output, nil
}

func (r *Runner) createPrice(ctx context.Context, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if r.Billing == nil {
		return nil, errors.New("billing service is required for price.create")
	}
	productID := stringValue(params, "product")
	if productID == "" {
		productID = stringValue(params, "productRef")
	}
	if mapped := state.productIDs[productID]; mapped != "" {
		productID = mapped
	}
	price, err := r.Billing.CreatePrice(ctx, billing.Price{
		ProductID:              productID,
		Currency:               stringDefault(params, "currency", "usd"),
		UnitAmount:             int64Value(params, "unitAmount"),
		RecurringInterval:      stringValue(params, "interval"),
		RecurringIntervalCount: intValue(params, "intervalCount"),
		Active:                 boolValue(params, "active", true),
	})
	if err != nil {
		return nil, err
	}
	output := map[string]any{"price": price}
	state.results[step.ID] = output
	if alias := stringValue(params, "id"); alias != "" {
		state.priceIDs[alias] = price.ID
		state.results[alias] = output
	}
	return output, nil
}

func (r *Runner) createCheckout(ctx context.Context, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if r.Billing == nil {
		return nil, errors.New("billing service is required for checkout.create")
	}
	customerID := firstString(params, "customerRef", "customer", "customer_id")
	lineItems := lineItems(params, state)
	if len(lineItems) == 0 {
		priceID := firstString(params, "priceRef", "price")
		if mapped := state.priceIDs[priceID]; mapped != "" {
			priceID = mapped
		}
		lineItems = []billing.LineItem{{PriceID: priceID, Quantity: int64ValueDefault(params, "quantity", 1)}}
	}
	session, err := r.Billing.CreateCheckoutSession(ctx, billing.CheckoutSession{
		CustomerID: customerID,
		Mode:       stringDefault(params, "mode", "subscription"),
		LineItems:  lineItems,
		SuccessURL: stringValue(params, "successUrl"),
		CancelURL:  stringValue(params, "cancelUrl"),
	})
	if err != nil {
		return nil, err
	}
	output := map[string]any{"session": session}
	state.results[step.ID] = output
	return output, nil
}

func (r *Runner) completeCheckout(ctx context.Context, scenario Scenario, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if r.Billing == nil {
		return nil, errors.New("billing service is required for checkout.complete")
	}
	sessionRef := firstString(params, "sessionRef", "session", "session_id")
	outcome := firstString(params, "outcome", "paymentMethod", "payment_method", "payment_method_id")
	if outcome == "" {
		outcome = "payment_succeeded"
	}
	session, err := r.Billing.CompleteCheckoutAt(ctx, sessionRef, outcome, state.clock.Now())
	if err != nil {
		return nil, err
	}
	output := map[string]any{"session": session}
	if session.SubscriptionID != "" {
		sub, err := r.Billing.GetSubscription(ctx, session.SubscriptionID)
		if err != nil {
			return nil, err
		}
		output["subscription"] = sub
	}
	if session.InvoiceID != "" {
		invoice, err := r.Billing.GetInvoice(ctx, session.InvoiceID)
		if err != nil {
			return nil, err
		}
		output["invoice"] = invoice
	}
	if session.PaymentIntentID != "" {
		intent, err := r.Billing.GetPaymentIntent(ctx, session.PaymentIntentID)
		if err != nil {
			return nil, err
		}
		output["payment_intent"] = intent
	}
	state.results[step.ID] = output
	state.aliasCheckoutCompletion(sessionRef, output)
	events, err := r.emitCheckoutWebhooks(ctx, scenario, output, state)
	if err != nil {
		return nil, err
	}
	if len(events) > 0 {
		output["events"] = events
		state.results[step.ID] = output
		state.aliasCheckoutCompletion(sessionRef, output)
	}
	return output, nil
}

func (r *Runner) advanceClock(ctx context.Context, scenario Scenario, step Step, params map[string]any, state *runState) (map[string]any, error) {
	raw := stringValue(params, "duration")
	before := state.clock.Now()
	d, err := state.clock.Advance(raw)
	if err != nil {
		return nil, err
	}
	output := map[string]any{
		"before":   before,
		"after":    state.clock.Now(),
		"duration": d.String(),
	}
	if r.Billing != nil {
		advance, err := r.Billing.AdvanceClock(ctx, state.clock.Now())
		if err != nil {
			return nil, err
		}
		output["billing"] = advance
		events, err := r.emitClockAdvanceWebhooks(ctx, scenario, advance, state)
		if err != nil {
			return nil, err
		}
		if len(events) > 0 {
			output["events"] = events
		}
	}
	state.results[step.ID] = output
	return output, nil
}

func (r *Runner) retryInvoice(ctx context.Context, scenario Scenario, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if r.Billing == nil {
		output := map[string]any{
			"subscription":  firstString(params, "subscriptionRef", "subscription", "subscription_id"),
			"invoice":       firstString(params, "invoiceRef", "invoice", "invoice_id"),
			"outcome":       stringDefault(params, "outcome", "payment_succeeded"),
			"clock":         state.clock.Now(),
			"deterministic": true,
			"note":          "invoice.retry recorded as scenario evidence because no billing service is configured",
		}
		state.results[step.ID] = output
		return output, nil
	}
	invoiceID := firstString(params, "invoiceRef", "invoice", "invoice_id")
	if invoiceID == "" {
		output := map[string]any{
			"subscription":  firstString(params, "subscriptionRef", "subscription", "subscription_id"),
			"invoice":       "",
			"outcome":       stringDefault(params, "outcome", "payment_succeeded"),
			"clock":         state.clock.Now(),
			"deterministic": true,
			"note":          "invoice.retry recorded as scenario evidence because no billing invoice was supplied",
		}
		state.results[step.ID] = output
		return output, nil
	}
	result, err := r.Billing.PayInvoice(ctx, invoiceID, billing.InvoicePaymentOptions{
		Outcome:         firstString(params, "outcome", "paymentMethod", "payment_method", "payment_method_id"),
		PaymentMethodID: firstString(params, "paymentMethod", "payment_method", "payment_method_id"),
		At:              state.clock.Now(),
	})
	if err != nil {
		return nil, err
	}
	output := map[string]any{
		"subscription":   result.Subscription,
		"invoice":        result.Invoice,
		"payment_intent": result.PaymentIntent,
		"outcome":        stringDefault(params, "outcome", stringDefault(params, "payment_method", "payment_succeeded")),
		"clock":          state.clock.Now(),
		"deterministic":  true,
	}
	events, err := r.emitInvoicePaymentWebhooks(ctx, scenario, result, state)
	if err != nil {
		return nil, err
	}
	if len(events) > 0 {
		output["events"] = events
	}
	state.results[step.ID] = output
	return output, nil
}

func (r *Runner) replayWebhook(ctx context.Context, step Step, params map[string]any, state *runState) (map[string]any, error) {
	if r.Webhooks == nil {
		return nil, errors.New("webhook service is required for webhook.replay")
	}
	eventID := firstString(params, "eventRef", "event", "event_id", "eventId")
	if eventID == "" {
		return nil, fmt.Errorf("%w: webhook.replay requires eventRef", ErrInvalidConfig)
	}
	attempts, err := r.Webhooks.ReplayEvent(ctx, eventID, replayOptions(params))
	if err != nil {
		return nil, err
	}
	output := map[string]any{
		"event_id":                eventID,
		"delivery_attempts":       attempts,
		"delivery_attempts_count": len(attempts),
	}
	state.results[step.ID] = output
	return output, nil
}

func (r *Runner) assertApp(ctx context.Context, scenario Scenario, step Step, params map[string]any, state *runState) (map[string]any, *Assertion, error) {
	target := stringValue(params, "target")
	expected := mapValue(params, "expected")
	assertionURL, err := assertionURL(scenario.App.Assertions.BaseURL, target)
	assertion := &Assertion{Target: target, Expected: expected, URL: assertionURL}
	if err != nil {
		assertion.Error = err.Error()
		return nil, assertion, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	payload := map[string]any{
		"target":   target,
		"expected": expected,
		"step": map[string]any{
			"id":     step.ID,
			"action": step.Action,
		},
		"context": state.results,
		"clock": map[string]any{
			"now": state.clock.Now(),
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		assertion.Error = err.Error()
		return nil, assertion, err
	}
	client := r.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, assertionURL, bytes.NewReader(body))
	if err != nil {
		assertion.Error = err.Error()
		return nil, assertion, fmt.Errorf("%w: %v", ErrAppCallbackFailure, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		assertion.Error = err.Error()
		return nil, assertion, fmt.Errorf("%w: %v", ErrAppCallbackFailure, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	assertion.ResponseStatus = resp.StatusCode
	assertion.ResponseBody = truncateEvidence(security.RedactText(string(respBody)), 8192)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		assertion.Error = resp.Status
		return nil, assertion, fmt.Errorf("%w: %s", ErrAppCallbackFailure, resp.Status)
	}
	if len(bytes.TrimSpace(respBody)) > 0 {
		var callback struct {
			Pass    *bool  `json:"pass"`
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		if err := json.Unmarshal(respBody, &callback); err == nil && callback.Pass != nil && !*callback.Pass {
			assertion.Error = firstNonEmpty(callback.Message, callback.Error, "assertion returned pass:false")
			return nil, assertion, fmt.Errorf("%w: %s", ErrAssertionFailed, assertion.Error)
		}
	}
	assertion.Pass = true
	output := map[string]any{"assertion": assertion}
	state.results[step.ID] = output
	return output, assertion, nil
}

func (r *Runner) emitCheckoutWebhooks(ctx context.Context, scenario Scenario, output map[string]any, state *runState) ([]webhooks.Event, error) {
	if r.Webhooks == nil {
		return nil, nil
	}
	events := []webhooks.Event{}
	for _, item := range checkoutWebhookPayloads(output) {
		state.sequence++
		event, _, err := r.Webhooks.CreateEvent(ctx, webhooks.EventInput{
			Type:           item.eventType,
			ObjectPayload:  item.payload,
			RequestID:      "req_" + item.objectID,
			IdempotencyKey: "billtap:" + item.eventType + ":" + item.objectID,
			ScenarioRunID:  scenario.Name,
			Source:         webhooks.SourceScenario,
			Sequence:       state.sequence,
		})
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *Runner) emitInvoicePaymentWebhooks(ctx context.Context, scenario Scenario, result billing.InvoicePaymentResult, state *runState) ([]webhooks.Event, error) {
	if r.Webhooks == nil {
		return nil, nil
	}
	var events []webhooks.Event
	for _, item := range invoicePaymentWebhookPayloads(result) {
		state.sequence++
		event, _, err := r.Webhooks.CreateEvent(ctx, webhooks.EventInput{
			Type:           item.eventType,
			ObjectPayload:  item.payload,
			RequestID:      "req_" + item.objectID,
			IdempotencyKey: "billtap:" + item.eventType + ":" + item.objectID,
			ScenarioRunID:  scenario.Name,
			Source:         webhooks.SourceScenario,
			Sequence:       state.sequence,
		})
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *Runner) emitClockAdvanceWebhooks(ctx context.Context, scenario Scenario, advance billing.ClockAdvanceResult, state *runState) ([]webhooks.Event, error) {
	if r.Webhooks == nil {
		return nil, nil
	}
	var events []webhooks.Event
	for _, renewal := range advance.Renewals {
		emitted, err := r.emitInvoicePaymentWebhooks(ctx, scenario, renewal, state)
		if err != nil {
			return events, err
		}
		events = append(events, emitted...)
	}
	for _, subscription := range advance.Canceled {
		raw, err := json.Marshal(subscription)
		if err != nil {
			continue
		}
		state.sequence++
		event, _, err := r.Webhooks.CreateEvent(ctx, webhooks.EventInput{
			Type:           "customer.subscription.deleted",
			ObjectPayload:  raw,
			RequestID:      "req_" + subscription.ID,
			IdempotencyKey: "billtap:customer.subscription.deleted:" + subscription.ID,
			ScenarioRunID:  scenario.Name,
			Source:         webhooks.SourceScenario,
			Sequence:       state.sequence,
		})
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

func replayOptions(params map[string]any) webhooks.ReplayOptions {
	delay := durationValue(params, "delay")
	if delay == 0 {
		delay = time.Duration(int64Value(params, "delay_seconds")) * time.Second
	}
	return webhooks.ReplayOptions{
		Duplicate:         int(int64ValueDefault(params, "duplicate", 1)),
		Delay:             delay,
		OutOfOrder:        boolValue(params, "outOfOrder", boolValue(params, "out_of_order", false)),
		ResponseStatus:    int(int64ValueDefault(params, "responseStatus", int64Value(params, "response_status"))),
		ResponseBody:      firstString(params, "responseBody", "response_body", "body"),
		SimulatedError:    firstString(params, "error", "simulatedError", "simulated_error"),
		SimulatedTimeout:  boolValue(params, "timeout", false),
		SignatureMismatch: boolValue(params, "signatureMismatch", boolValue(params, "signature_mismatch", false)),
	}
}

type webhookPayload struct {
	eventType string
	objectID  string
	payload   json.RawMessage
}

func checkoutWebhookPayloads(result map[string]any) []webhookPayload {
	session, _ := result["session"].(billing.CheckoutSession)
	subscription, _ := result["subscription"].(billing.Subscription)
	invoice, _ := result["invoice"].(billing.Invoice)
	paymentIntent, _ := result["payment_intent"].(billing.PaymentIntent)

	var out []webhookPayload
	appendPayload := func(eventType string, objectID string, value any) {
		raw, err := json.Marshal(value)
		if err == nil {
			out = append(out, webhookPayload{eventType: eventType, objectID: objectID, payload: raw})
		}
	}
	appendPayload(checkoutSessionEvent(session.Status), session.ID, session)
	if subscription.ID != "" {
		appendPayload("customer.subscription.created", subscription.ID, subscription)
	}
	if invoice.ID != "" {
		appendPayload("invoice.created", invoice.ID, invoice)
		appendPayload("invoice.finalized", invoice.ID, invoice)
	}
	if paymentIntent.ID != "" {
		appendPayload("payment_intent.created", paymentIntent.ID, paymentIntent)
		if eventType := paymentIntentTerminalEvent(paymentIntent.Status); eventType != "" {
			appendPayload(eventType, paymentIntent.ID, paymentIntent)
		}
	}
	if invoice.ID != "" {
		for _, eventType := range invoiceTerminalEvents(invoice.Status, paymentIntent.Status) {
			appendPayload(eventType, invoice.ID, invoice)
		}
	}
	if subscription.ID != "" {
		appendPayload("customer.subscription.updated", subscription.ID, subscription)
	}
	return out
}

func invoicePaymentWebhookPayloads(result billing.InvoicePaymentResult) []webhookPayload {
	var out []webhookPayload
	appendPayload := func(eventType string, objectID string, value any) {
		if eventType == "" || objectID == "" {
			return
		}
		raw, err := json.Marshal(value)
		if err == nil {
			out = append(out, webhookPayload{eventType: eventType, objectID: objectID, payload: raw})
		}
	}
	if result.PaymentIntent.ID != "" {
		appendPayload(paymentIntentTerminalEvent(result.PaymentIntent.Status), result.PaymentIntent.ID, result.PaymentIntent)
	}
	if result.Invoice.ID != "" {
		for _, eventType := range invoiceTerminalEvents(result.Invoice.Status, result.PaymentIntent.Status) {
			appendPayload(eventType, result.Invoice.ID, result.Invoice)
		}
	}
	if result.Subscription.ID != "" {
		appendPayload("customer.subscription.updated", result.Subscription.ID, result.Subscription)
	}
	return out
}

func paymentIntentTerminalEvent(status string) string {
	switch status {
	case "succeeded":
		return "payment_intent.succeeded"
	case "processing":
		return "payment_intent.processing"
	case "canceled":
		return "payment_intent.canceled"
	case "requires_payment_method":
		return "payment_intent.payment_failed"
	default:
		return "payment_intent.payment_failed"
	}
}

func checkoutSessionEvent(status string) string {
	if status == "expired" {
		return "checkout.session.expired"
	}
	return "checkout.session.completed"
}

func invoiceTerminalEvents(status string, paymentIntentStatus string) []string {
	switch status {
	case "paid":
		return []string{"invoice.payment_succeeded"}
	case "void":
		return []string{"invoice.voided"}
	case "open":
		if paymentIntentStatus == "processing" {
			return nil
		}
		return []string{"invoice.payment_failed"}
	default:
		return []string{"invoice.payment_failed"}
	}
}

func (s *runState) resolveParams(params map[string]any) (map[string]any, error) {
	if params == nil {
		return map[string]any{}, nil
	}
	out := make(map[string]any, len(params))
	for key, value := range params {
		resolved, err := s.resolveValue(value)
		if err != nil {
			return nil, fmt.Errorf("resolve params.%s: %w", key, err)
		}
		out[key] = resolved
	}
	return out, nil
}

func (s *runState) resolveValue(value any) (any, error) {
	switch typed := value.(type) {
	case string:
		return s.resolveString(typed)
	case []any:
		out := make([]any, len(typed))
		for i, item := range typed {
			resolved, err := s.resolveValue(item)
			if err != nil {
				return nil, err
			}
			out[i] = resolved
		}
		return out, nil
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, item := range typed {
			resolved, err := s.resolveValue(item)
			if err != nil {
				return nil, err
			}
			out[key] = resolved
		}
		return out, nil
	default:
		return value, nil
	}
}

func (s *runState) resolveString(value string) (any, error) {
	stepID, path, found := strings.Cut(value, ".")
	if !found {
		return value, nil
	}
	result, ok := s.results[stepID]
	if !ok {
		return value, nil
	}
	resolved, ok := valueAtPath(result, path)
	if !ok {
		return nil, fmt.Errorf("reference %q not found", value)
	}
	return resolved, nil
}

func (s *runState) aliasCheckoutCompletion(sessionRef string, output map[string]any) {
	for stepID, result := range s.results {
		session, ok := result["session"].(billing.CheckoutSession)
		if !ok || session.ID != sessionRef {
			continue
		}
		for key, value := range output {
			result[key] = value
		}
		s.results[stepID] = result
	}
}

func valueAtPath(root any, path string) (any, bool) {
	current := root
	for _, segment := range strings.Split(path, ".") {
		if segment == "" {
			return nil, false
		}
		var ok bool
		current, ok = fieldValue(current, segment)
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func fieldValue(value any, name string) (any, bool) {
	if value == nil {
		return nil, false
	}
	if m, ok := value.(map[string]any); ok {
		v, found := m[name]
		return v, found
	}
	raw, err := json.Marshal(value)
	if err == nil {
		var m map[string]any
		if err := json.Unmarshal(raw, &m); err == nil {
			if v, found := m[name]; found {
				return v, true
			}
		}
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, false
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		idx := -1
		if _, err := fmt.Sscan(name, &idx); err == nil && idx >= 0 && idx < rv.Len() {
			return rv.Index(idx).Interface(), true
		}
		return nil, false
	}
	if rv.Kind() != reflect.Struct {
		return nil, false
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonName == name || strings.EqualFold(field.Name, name) {
			return rv.Field(i).Interface(), true
		}
	}
	return nil, false
}

func assertionURL(baseURL, target string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("baseUrl must be absolute")
	}
	var escaped []string
	for _, segment := range strings.FieldsFunc(target, func(r rune) bool { return r == '.' || r == '/' }) {
		if segment == "" {
			continue
		}
		escaped = append(escaped, url.PathEscape(segment))
	}
	if len(escaped) == 0 {
		return "", fmt.Errorf("target is required")
	}
	basePath := strings.TrimRight(parsed.Path, "/")
	parsed.Path = basePath + "/" + strings.Join(escaped, "/")
	return parsed.String(), nil
}

func lineItems(params map[string]any, state *runState) []billing.LineItem {
	raw, ok := params["lineItems"]
	if !ok {
		raw, ok = params["line_items"]
	}
	if !ok {
		return nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]billing.LineItem, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		priceID := firstString(m, "priceRef", "price")
		if mapped := state.priceIDs[priceID]; mapped != "" {
			priceID = mapped
		}
		out = append(out, billing.LineItem{PriceID: priceID, Quantity: int64ValueDefault(m, "quantity", 1)})
	}
	return out
}

func mapValue(params map[string]any, key string) map[string]any {
	value, ok := params[key]
	if !ok {
		return nil
	}
	if m, ok := value.(map[string]any); ok {
		return m
	}
	return nil
}

func firstString(params map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := stringValue(params, key); value != "" {
			return value
		}
	}
	return ""
}

func stringValue(params map[string]any, key string) string {
	value, ok := params[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}

func stringDefault(params map[string]any, key, fallback string) string {
	if value := stringValue(params, key); value != "" {
		return value
	}
	return fallback
}

func intValue(params map[string]any, key string) int {
	return int(int64Value(params, key))
}

func int64Value(params map[string]any, key string) int64 {
	return int64ValueDefault(params, key, 0)
}

func int64ValueDefault(params map[string]any, key string, fallback int64) int64 {
	value, ok := params[key]
	if !ok || value == nil {
		return fallback
	}
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int64:
		return typed
	case int32:
		return int64(typed)
	case float64:
		return int64(typed)
	case float32:
		return int64(typed)
	case uint64:
		return int64(typed)
	case uint:
		return int64(typed)
	case string:
		var out int64
		if _, err := fmt.Sscan(typed, &out); err == nil {
			return out
		}
	}
	return fallback
}

func boolValue(params map[string]any, key string, fallback bool) bool {
	value, ok := params[key]
	if !ok || value == nil {
		return fallback
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		switch strings.ToLower(strings.TrimSpace(typed)) {
		case "true", "1", "yes":
			return true
		case "false", "0", "no":
			return false
		}
	}
	return fallback
}

func durationValue(params map[string]any, key string) time.Duration {
	value, ok := params[key]
	if !ok || value == nil {
		return 0
	}
	switch typed := value.(type) {
	case time.Duration:
		return typed
	case string:
		if parsed, err := time.ParseDuration(typed); err == nil {
			return parsed
		}
	case int:
		return time.Duration(typed) * time.Second
	case int64:
		return time.Duration(typed) * time.Second
	case int32:
		return time.Duration(typed) * time.Second
	case float64:
		return time.Duration(typed) * time.Second
	case float32:
		return time.Duration(typed) * time.Second
	}
	return 0
}

func failureType(err error) string {
	switch {
	case errors.Is(err, ErrInvalidConfig):
		return FailureInvalidConfig
	case errors.Is(err, ErrAppCallbackFailure):
		return FailureAppCallback
	case errors.Is(err, ErrAssertionFailed):
		return FailureAssertion
	default:
		return FailureRunner
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func truncateEvidence(value string, limit int) string {
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit] + "...[truncated]"
}
