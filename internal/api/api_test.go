package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/diagnostics"
	"github.com/hckim/billtap/internal/fixtures"
	"github.com/hckim/billtap/internal/scenarios"
	"github.com/hckim/billtap/internal/security"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
)

func TestCheckoutMVPFlow(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"buyer@example.test"},
		"name":  {"Buyer"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{
		"name": {"Team"},
	})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"9900"},
		"recurring[interval]": {"month"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"mode":                    {"subscription"},
		"line_items[0][price]":    {price.ID},
		"line_items[0][quantity]": {"2"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})

	if session.URL == "" || session.Status != "open" {
		t.Fatalf("session = %#v, want open session with hosted URL", session)
	}

	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})
	var completed billing.CheckoutSession
	if err := json.Unmarshal(completion["session"], &completed); err != nil {
		t.Fatalf("decode completed session: %v", err)
	}
	if completed.PaymentStatus != "paid" || completed.SubscriptionID == "" || completed.InvoiceID == "" || completed.PaymentIntentID == "" {
		t.Fatalf("completed session = %#v, want paid session with billing objects", completed)
	}

	var subscription billing.Subscription
	if err := json.Unmarshal(completion["subscription"], &subscription); err != nil {
		t.Fatalf("decode subscription: %v", err)
	}
	if subscription.Status != "active" {
		t.Fatalf("subscription status = %q, want active", subscription.Status)
	}

	timeline := getJSON[struct {
		Object string                  `json:"object"`
		Data   []billing.TimelineEntry `json:"data"`
	}](t, handler, "/api/timeline?checkoutSessionId="+session.ID)
	if got := len(timeline.Data); got < 4 {
		t.Fatalf("timeline entries = %d, want checkout completion evidence", got)
	}
}

func TestHostedURLsUseConfiguredPublicBaseURL(t *testing.T) {
	handler := newTestHandlerWithOptions(t, Options{PublicBaseURL: "http://127.0.0.1:18080/"})

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"public-url@example.test"},
		"name":  {"Public URL"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"9900"},
		"recurring[interval]": {"month"},
	})

	checkout := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {price.ID},
		"line_items[0][quantity]": {"1"},
	})
	if !strings.HasPrefix(checkout.URL, "http://127.0.0.1:18080/checkout/") {
		t.Fatalf("checkout url = %q, want configured public base URL", checkout.URL)
	}

	portal := postForm[struct {
		URL string `json:"url"`
	}](t, handler, "/v1/billing_portal/sessions", url.Values{"customer": {customer.ID}})
	if !strings.HasPrefix(portal.URL, "http://127.0.0.1:18080/portal?") || !strings.Contains(portal.URL, "customer_id="+customer.ID) {
		t.Fatalf("portal url = %q, want configured public base URL", portal.URL)
	}
}

func TestCheckoutFailureOutcome(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"buyer@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})

	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_failed",
	})
	var invoice billing.Invoice
	if err := json.Unmarshal(completion["invoice"], &invoice); err != nil {
		t.Fatalf("decode invoice: %v", err)
	}
	var paymentIntent billing.PaymentIntent
	if err := json.Unmarshal(completion["payment_intent"], &paymentIntent); err != nil {
		t.Fatalf("decode payment intent: %v", err)
	}
	if invoice.Status != "open" || paymentIntent.Status != "requires_payment_method" {
		t.Fatalf("invoice=%s payment_intent=%s, want failed checkout state", invoice.Status, paymentIntent.Status)
	}
}

func TestInvoicePayRetriesFailedCheckout(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"retry@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_failed",
	})
	var failedInvoice billing.Invoice
	if err := json.Unmarshal(completion["invoice"], &failedInvoice); err != nil {
		t.Fatalf("decode invoice: %v", err)
	}

	declined := postForm[struct {
		ID                 string `json:"id"`
		Status             string `json:"status"`
		AttemptCount       int    `json:"attempt_count"`
		AmountDue          int64  `json:"amount_due"`
		AmountPaid         int64  `json:"amount_paid"`
		NextPaymentAttempt *int64 `json:"next_payment_attempt"`
	}](t, handler, "/v1/invoices/"+failedInvoice.ID+"/pay", url.Values{
		"source":      {"pm_card_visa_chargeDeclined"},
		"off_session": {"true"},
		"mandate":     {"mandate_123"},
	})
	if declined.Status != "open" || declined.AttemptCount != failedInvoice.AttemptCount+1 || declined.AmountDue != 9900 || declined.AmountPaid != 0 || declined.NextPaymentAttempt == nil {
		t.Fatalf("declined retry invoice = %#v, want open invoice with retry evidence", declined)
	}

	paid := postForm[struct {
		ID                 string `json:"id"`
		Status             string `json:"status"`
		AttemptCount       int    `json:"attempt_count"`
		AmountDue          int64  `json:"amount_due"`
		AmountPaid         int64  `json:"amount_paid"`
		NextPaymentAttempt *int64 `json:"next_payment_attempt"`
	}](t, handler, "/v1/invoices/"+failedInvoice.ID+"/pay", url.Values{
		"payment_method": {"pm_card_visa"},
	})
	if paid.Status != "paid" || paid.AttemptCount != declined.AttemptCount+1 || paid.AmountDue != 0 || paid.AmountPaid != 9900 || paid.NextPaymentAttempt != nil {
		t.Fatalf("paid retry invoice = %#v, want paid invoice after retry", paid)
	}

	subscription := getJSON[struct {
		Status        string `json:"status"`
		LatestInvoice string `json:"latest_invoice"`
	}](t, handler, "/v1/subscriptions/"+failedInvoice.SubscriptionID)
	if subscription.Status != "active" || subscription.LatestInvoice != failedInvoice.ID {
		t.Fatalf("subscription = %#v, want active with retried invoice", subscription)
	}

	events := getJSON[struct {
		Data []struct {
			Type string `json:"type"`
		} `json:"data"`
	}](t, handler, "/v1/events")
	seen := map[string]bool{}
	for _, event := range events.Data {
		seen[event.Type] = true
	}
	for _, eventType := range []string{"payment_intent.payment_failed", "invoice.payment_failed", "payment_intent.succeeded", "invoice.payment_succeeded", "invoice.paid", "customer.subscription.updated"} {
		if !seen[eventType] {
			t.Fatalf("events missing %s: %#v", eventType, events.Data)
		}
	}
}

func TestCheckoutTerminalOutcomeVariants(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"variants@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})

	tests := []struct {
		name                string
		outcome             string
		sessionStatus       string
		paymentStatus       string
		subscriptionStatus  string
		invoiceStatus       string
		paymentIntentStatus string
		events              []string
		notEvents           []string
	}{
		{
			name:                "pending async payment",
			outcome:             "payment_pending",
			sessionStatus:       "complete",
			paymentStatus:       "unpaid",
			subscriptionStatus:  "incomplete",
			invoiceStatus:       "open",
			paymentIntentStatus: "processing",
			events:              []string{"checkout.session.completed", "payment_intent.processing", "customer.subscription.updated"},
			notEvents:           []string{"invoice.payment_failed"},
		},
		{
			name:                "canceled checkout",
			outcome:             "canceled",
			sessionStatus:       "expired",
			paymentStatus:       "unpaid",
			subscriptionStatus:  "incomplete_expired",
			invoiceStatus:       "void",
			paymentIntentStatus: "canceled",
			events:              []string{"checkout.session.expired", "payment_intent.canceled", "invoice.voided"},
			notEvents:           []string{"checkout.session.completed", "invoice.payment_failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
				"customer":             {customer.ID},
				"line_items[0][price]": {price.ID},
			})
			completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
				"outcome": tt.outcome,
			})

			var completed billing.CheckoutSession
			if err := json.Unmarshal(completion["session"], &completed); err != nil {
				t.Fatalf("decode checkout session: %v", err)
			}
			var subscription billing.Subscription
			if err := json.Unmarshal(completion["subscription"], &subscription); err != nil {
				t.Fatalf("decode subscription: %v", err)
			}
			var invoice billing.Invoice
			if err := json.Unmarshal(completion["invoice"], &invoice); err != nil {
				t.Fatalf("decode invoice: %v", err)
			}
			var paymentIntent billing.PaymentIntent
			if err := json.Unmarshal(completion["payment_intent"], &paymentIntent); err != nil {
				t.Fatalf("decode payment intent: %v", err)
			}

			if completed.Status != tt.sessionStatus || completed.PaymentStatus != tt.paymentStatus {
				t.Fatalf("session status=%s payment_status=%s, want %s/%s", completed.Status, completed.PaymentStatus, tt.sessionStatus, tt.paymentStatus)
			}
			if subscription.Status != tt.subscriptionStatus || invoice.Status != tt.invoiceStatus || paymentIntent.Status != tt.paymentIntentStatus {
				t.Fatalf("subscription=%s invoice=%s payment_intent=%s, want %s/%s/%s", subscription.Status, invoice.Status, paymentIntent.Status, tt.subscriptionStatus, tt.invoiceStatus, tt.paymentIntentStatus)
			}

			var events []webhooks.Event
			if err := json.Unmarshal(completion["webhook_events"], &events); err != nil {
				t.Fatalf("decode webhook events: %v", err)
			}
			types := map[string]bool{}
			for _, event := range events {
				types[event.Type] = true
			}
			for _, eventType := range tt.events {
				if !types[eventType] {
					t.Fatalf("webhook events missing %s in %#v", eventType, types)
				}
			}
			for _, eventType := range tt.notEvents {
				if types[eventType] {
					t.Fatalf("webhook events unexpectedly include %s in %#v", eventType, types)
				}
			}
		})
	}
}

func TestCheckoutPaymentErrorSimulation(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"errors@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})

	tests := []struct {
		name          string
		paymentMethod string
		status        string
		code          string
		declineCode   string
	}{
		{
			name:          "insufficient funds",
			paymentMethod: "pm_card_visa_chargeDeclinedInsufficientFunds",
			status:        "requires_payment_method",
			code:          "card_declined",
			declineCode:   "insufficient_funds",
		},
		{
			name:          "requires action",
			paymentMethod: "pm_card_threeDSecure2Required",
			status:        "requires_action",
			code:          "authentication_required",
			declineCode:   "authentication_required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
				"customer":             {customer.ID},
				"line_items[0][price]": {price.ID},
			})
			completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
				"payment_method": tt.paymentMethod,
			})
			var invoice billing.Invoice
			if err := json.Unmarshal(completion["invoice"], &invoice); err != nil {
				t.Fatalf("decode invoice: %v", err)
			}
			var paymentIntent billing.PaymentIntent
			if err := json.Unmarshal(completion["payment_intent"], &paymentIntent); err != nil {
				t.Fatalf("decode payment intent: %v", err)
			}
			if invoice.Status != "open" || paymentIntent.Status != tt.status {
				t.Fatalf("invoice=%s payment_intent=%s, want open/%s", invoice.Status, paymentIntent.Status, tt.status)
			}
			if paymentIntent.FailureCode != tt.code || paymentIntent.DeclineCode != tt.declineCode || paymentIntent.PaymentMethodID != tt.paymentMethod {
				t.Fatalf("payment intent = %#v, want code=%s decline=%s method=%s", paymentIntent, tt.code, tt.declineCode, tt.paymentMethod)
			}

			projected := getJSON[struct {
				Status           string `json:"status"`
				PaymentMethodID  string `json:"payment_method"`
				LastPaymentError struct {
					Type        string `json:"type"`
					Code        string `json:"code"`
					DeclineCode string `json:"decline_code"`
					Message     string `json:"message"`
				} `json:"last_payment_error"`
			}](t, handler, "/v1/payment_intents/"+paymentIntent.ID)
			if projected.Status != tt.status || projected.PaymentMethodID != tt.paymentMethod {
				t.Fatalf("projected payment intent = %#v", projected)
			}
			if projected.LastPaymentError.Type != "card_error" || projected.LastPaymentError.Code != tt.code || projected.LastPaymentError.DeclineCode != tt.declineCode || projected.LastPaymentError.Message == "" {
				t.Fatalf("last_payment_error = %#v, want card error code=%s decline=%s", projected.LastPaymentError, tt.code, tt.declineCode)
			}
		})
	}
}

func TestDirectPaymentIntentAndSetupIntentStateMachines(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"direct-intents@example.test"}})

	created := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents", url.Values{
		"amount":   {"4900"},
		"currency": {"usd"},
		"customer": {customer.ID},
	})
	if created.Status != "requires_payment_method" || created.CustomerID != customer.ID || created.Amount != 4900 {
		t.Fatalf("created payment intent = %#v, want requires_payment_method with customer/amount", created)
	}

	confirmed := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents/"+created.ID+"/confirm", url.Values{
		"payment_method": {"pm_card_visa"},
	})
	if confirmed.Status != "succeeded" || confirmed.PaymentMethodID != "pm_card_visa" {
		t.Fatalf("confirmed payment intent = %#v, want succeeded with payment method", confirmed)
	}
	reconfirmStatus, reconfirmBody := postFormStatus(t, handler, "/v1/payment_intents/"+confirmed.ID+"/confirm", url.Values{
		"payment_method": {"pm_card_visa"},
	})
	reconfirmErr := decodeErrorBody(t, reconfirmBody)
	if reconfirmStatus != http.StatusBadRequest || reconfirmErr.Error.Code != "parameter_invalid" || reconfirmErr.Error.Param != "status" {
		t.Fatalf("reconfirm status=%d error=%#v, want invalid terminal status", reconfirmStatus, reconfirmErr.Error)
	}
	cancelSucceededStatus, cancelSucceededBody := postFormStatus(t, handler, "/v1/payment_intents/"+confirmed.ID+"/cancel", url.Values{})
	cancelSucceededErr := decodeErrorBody(t, cancelSucceededBody)
	if cancelSucceededStatus != http.StatusBadRequest || cancelSucceededErr.Error.Code != "parameter_invalid" || cancelSucceededErr.Error.Param != "status" {
		t.Fatalf("cancel succeeded status=%d error=%#v, want invalid terminal status", cancelSucceededStatus, cancelSucceededErr.Error)
	}

	captureBeforeConfirmStatus, captureBeforeConfirmBody := postFormStatus(t, handler, "/v1/payment_intents/"+created.ID+"/capture", url.Values{})
	captureBeforeConfirmErr := decodeErrorBody(t, captureBeforeConfirmBody)
	if captureBeforeConfirmStatus != http.StatusBadRequest || captureBeforeConfirmErr.Error.Code != "parameter_invalid" || captureBeforeConfirmErr.Error.Param != "status" {
		t.Fatalf("capture before confirm status=%d error=%#v, want invalid status", captureBeforeConfirmStatus, captureBeforeConfirmErr.Error)
	}

	manual := postForm[struct {
		ID                string `json:"id"`
		Status            string `json:"status"`
		CaptureMethod     string `json:"capture_method"`
		AmountCapturable  int64  `json:"amount_capturable"`
		AmountReceived    int64  `json:"amount_received"`
		LastPaymentError  any    `json:"last_payment_error"`
		PaymentMethodID   string `json:"payment_method"`
		ClientSecret      string `json:"client_secret"`
		CreatedUnix       int64  `json:"created"`
		CreatedAtNotEmpty any    `json:"created_at"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"amount":         {"6600"},
		"currency":       {"usd"},
		"capture_method": {"manual"},
		"payment_method": {"pm_card_visa"},
		"confirm":        {"true"},
	})
	if manual.Status != "requires_capture" || manual.CaptureMethod != "manual" || manual.AmountCapturable != 6600 || manual.AmountReceived != 0 {
		t.Fatalf("manual payment intent = %#v, want requires_capture with capturable amount", manual)
	}
	partialCaptureStatus, partialCaptureBody := postFormStatus(t, handler, "/v1/payment_intents/"+manual.ID+"/capture", url.Values{
		"amount_to_capture": {"100"},
	})
	partialCaptureErr := decodeErrorBody(t, partialCaptureBody)
	if partialCaptureStatus != http.StatusBadRequest || partialCaptureErr.Error.Code != "parameter_invalid" || partialCaptureErr.Error.Param != "amount_to_capture" {
		t.Fatalf("partial capture status=%d error=%#v, want invalid amount_to_capture", partialCaptureStatus, partialCaptureErr.Error)
	}
	captured := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents/"+manual.ID+"/capture", url.Values{
		"amount_to_capture": {"6600"},
	})
	if captured.Status != "succeeded" {
		t.Fatalf("captured payment intent = %#v, want succeeded", captured)
	}

	declined := postForm[struct {
		ID               string `json:"id"`
		Status           string `json:"status"`
		LastPaymentError struct {
			Type        string `json:"type"`
			Code        string `json:"code"`
			DeclineCode string `json:"decline_code"`
			Message     string `json:"message"`
		} `json:"last_payment_error"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"amount":         {"4900"},
		"currency":       {"usd"},
		"payment_method": {"pm_card_visa_chargeDeclined"},
		"confirm":        {"true"},
	})
	if declined.Status != "requires_payment_method" || declined.LastPaymentError.Type != "card_error" || declined.LastPaymentError.Code != "card_declined" || declined.LastPaymentError.DeclineCode != "generic_decline" {
		t.Fatalf("declined payment intent = %#v, want card_declined generic_decline", declined)
	}

	canceled := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents/"+declined.ID+"/cancel", url.Values{
		"cancellation_reason": {"requested_by_customer"},
	})
	if canceled.Status != "canceled" {
		t.Fatalf("canceled payment intent = %#v, want canceled", canceled)
	}

	defaultOutcomeCustomer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"customer-default-outcome@example.test"},
		"metadata[billtap_default_payment_intent_outcome]": {"card_declined"},
	})
	customerDefaultDeclined := postForm[struct {
		Status           string            `json:"status"`
		Metadata         map[string]string `json:"metadata"`
		LastPaymentError struct {
			Code        string `json:"code"`
			DeclineCode string `json:"decline_code"`
		} `json:"last_payment_error"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"amount":   {"3100"},
		"currency": {"usd"},
		"customer": {defaultOutcomeCustomer.ID},
		"confirm":  {"true"},
	})
	if customerDefaultDeclined.Status != "requires_payment_method" || customerDefaultDeclined.Metadata[billing.MetadataPaymentIntentOutcome] != "card_declined" || customerDefaultDeclined.LastPaymentError.Code != "card_declined" {
		t.Fatalf("customer default declined intent = %#v, want customer default outcome applied", customerDefaultDeclined)
	}
	customerDefaultOverride := postForm[struct {
		Status        string            `json:"status"`
		Metadata      map[string]string `json:"metadata"`
		PaymentMethod string            `json:"payment_method"`
		NextAction    map[string]any    `json:"next_action"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"amount":   {"3100"},
		"currency": {"usd"},
		"customer": {defaultOutcomeCustomer.ID},
		"confirm":  {"true"},
		"metadata[billtap_payment_intent_outcome]": {"requires_action"},
	})
	if customerDefaultOverride.Status != "requires_action" || customerDefaultOverride.Metadata[billing.MetadataPaymentIntentOutcome] != "requires_action" || customerDefaultOverride.PaymentMethod != "pm_card_threeDSecure2Required" || customerDefaultOverride.NextAction["type"] != "use_stripe_sdk" {
		t.Fatalf("customer default override intent = %#v, want per-intent outcome to win", customerDefaultOverride)
	}

	deferredDecline := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents", url.Values{
		"amount":   {"2700"},
		"currency": {"usd"},
		"customer": {customer.ID},
		"metadata[billtap_payment_intent_outcome]": {"card_declined"},
		"metadata[fixture_ref]":                    {"deferred-decline"},
	})
	if deferredDecline.Status != "requires_payment_method" || deferredDecline.Metadata[billing.MetadataPaymentIntentOutcome] != "card_declined" {
		t.Fatalf("deferred decline intent = %#v, want stored outcome metadata without confirmation", deferredDecline)
	}
	deferredDeclined := postForm[struct {
		Status           string            `json:"status"`
		PaymentMethod    string            `json:"payment_method"`
		Metadata         map[string]string `json:"metadata"`
		LastPaymentError struct {
			Code        string `json:"code"`
			DeclineCode string `json:"decline_code"`
		} `json:"last_payment_error"`
	}](t, handler, "/v1/payment_intents/"+deferredDecline.ID+"/confirm", url.Values{
		"payment_method": {"pm_card_visa"},
	})
	if deferredDeclined.Status != "requires_payment_method" || deferredDeclined.Metadata[billing.MetadataPaymentIntentOutcome] != "card_declined" || deferredDeclined.LastPaymentError.Code != "card_declined" || deferredDeclined.LastPaymentError.DeclineCode != "generic_decline" {
		t.Fatalf("deferred declined intent = %#v, want configured card_declined outcome on confirm", deferredDeclined)
	}

	deferredSCA := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents", url.Values{
		"amount":   {"3200"},
		"currency": {"usd"},
		"customer": {customer.ID},
	})
	configuredSCA := postJSON[struct {
		Status   string            `json:"status"`
		Metadata map[string]string `json:"metadata"`
	}](t, handler, "/api/payment_intents/"+deferredSCA.ID+"/outcome", map[string]string{
		"outcome": "requires_action",
	})
	if configuredSCA.Status != "requires_payment_method" || configuredSCA.Metadata[billing.MetadataPaymentIntentOutcome] != "requires_action" {
		t.Fatalf("configured SCA intent = %#v, want deferred outcome metadata", configuredSCA)
	}
	confirmedSCA := postForm[struct {
		Status        string         `json:"status"`
		PaymentMethod string         `json:"payment_method"`
		NextAction    map[string]any `json:"next_action"`
	}](t, handler, "/v1/payment_intents/"+deferredSCA.ID+"/confirm", url.Values{
		"payment_method": {"pm_card_visa"},
	})
	if confirmedSCA.Status != "requires_action" || confirmedSCA.PaymentMethod != "pm_card_visa" || confirmedSCA.NextAction["type"] != "use_stripe_sdk" {
		t.Fatalf("confirmed SCA intent = %#v, want configured requires_action outcome", confirmedSCA)
	}

	setup := postForm[billing.SetupIntent](t, handler, "/v1/setup_intents", url.Values{
		"customer":       {customer.ID},
		"payment_method": {"pm_card_visa"},
		"usage":          {"off_session"},
		"confirm":        {"true"},
	})
	if setup.Status != "succeeded" || setup.CustomerID != customer.ID || setup.PaymentMethodID != "pm_card_visa" {
		t.Fatalf("setup intent = %#v, want succeeded with customer/payment_method", setup)
	}
	setupList := getJSON[struct {
		Object string                `json:"object"`
		Data   []billing.SetupIntent `json:"data"`
	}](t, handler, "/v1/setup_intents?customer="+customer.ID)
	if setupList.Object != "list" || len(setupList.Data) != 1 || setupList.Data[0].ID != setup.ID {
		t.Fatalf("setup intent list = %#v, want one setup intent for customer", setupList)
	}
	cancelableSetup := postForm[billing.SetupIntent](t, handler, "/v1/setup_intents", url.Values{
		"customer": {customer.ID},
		"usage":    {"off_session"},
	})
	setupCanceled := postForm[billing.SetupIntent](t, handler, "/v1/setup_intents/"+cancelableSetup.ID+"/cancel", url.Values{})
	if setupCanceled.Status != "canceled" {
		t.Fatalf("canceled setup intent = %#v, want canceled", setupCanceled)
	}
	setupCancelSucceededStatus, setupCancelSucceededBody := postFormStatus(t, handler, "/v1/setup_intents/"+setup.ID+"/cancel", url.Values{})
	setupCancelSucceededErr := decodeErrorBody(t, setupCancelSucceededBody)
	if setupCancelSucceededStatus != http.StatusBadRequest || setupCancelSucceededErr.Error.Code != "parameter_invalid" || setupCancelSucceededErr.Error.Param != "status" {
		t.Fatalf("setup cancel succeeded status=%d error=%#v, want invalid terminal status", setupCancelSucceededStatus, setupCancelSucceededErr.Error)
	}
	setupConfirmCanceledStatus, setupConfirmCanceledBody := postFormStatus(t, handler, "/v1/setup_intents/"+setupCanceled.ID+"/confirm", url.Values{
		"payment_method": {"pm_card_visa"},
	})
	setupConfirmCanceledErr := decodeErrorBody(t, setupConfirmCanceledBody)
	if setupConfirmCanceledStatus != http.StatusBadRequest || setupConfirmCanceledErr.Error.Code != "parameter_invalid" || setupConfirmCanceledErr.Error.Param != "status" {
		t.Fatalf("setup confirm canceled status=%d error=%#v, want invalid terminal status", setupConfirmCanceledStatus, setupConfirmCanceledErr.Error)
	}

	timeline := getJSON[struct {
		Object string                  `json:"object"`
		Data   []billing.TimelineEntry `json:"data"`
	}](t, handler, "/api/timeline?paymentIntentId="+created.ID)
	if len(timeline.Data) < 2 {
		t.Fatalf("direct payment intent timeline entries = %d, want create and confirm evidence", len(timeline.Data))
	}

	events := getJSON[struct {
		Object string           `json:"object"`
		Data   []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events")
	eventTypes := map[string]bool{}
	for _, event := range events.Data {
		eventTypes[event.Type] = true
	}
	for _, eventType := range []string{"payment_intent.created", "payment_intent.succeeded", "payment_intent.amount_capturable_updated", "payment_intent.canceled", "setup_intent.created", "setup_intent.succeeded", "setup_intent.canceled"} {
		if !eventTypes[eventType] {
			t.Fatalf("events missing %s in %#v", eventType, eventTypes)
		}
	}
}

func TestWebhookEndpointDeliveryAndReplay(t *testing.T) {
	var signatures []string
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signatures = append(signatures, r.Header.Get(webhooks.SignatureHeaderName))
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()

	handler := newTestHandler(t)
	endpoint := postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL},
		"enabled_events": {"checkout.session.completed,invoice.*"},
	})
	if endpoint.Secret != security.MaskedValue {
		t.Fatalf("webhook endpoint secret = %q, want masked", endpoint.Secret)
	}

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"buyer@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	_ = postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})

	events := getJSON[struct {
		Object string           `json:"object"`
		Data   []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events")
	if len(events.Data) < 4 {
		t.Fatalf("events = %d, want checkout webhook events", len(events.Data))
	}
	eventIDs := map[string]string{}
	for _, event := range events.Data {
		if eventIDs[event.Type] == "" {
			eventIDs[event.Type] = event.ID
		}
	}
	checkoutEventID := eventIDs["checkout.session.completed"]
	invoiceCreatedEventID := eventIDs["invoice.created"]
	invoiceFinalizedEventID := eventIDs["invoice.finalized"]
	invoicePaidEventID := eventIDs["invoice.paid"]
	if checkoutEventID == "" || invoiceCreatedEventID == "" || invoiceFinalizedEventID == "" || invoicePaidEventID == "" {
		t.Fatalf("events = %#v, want checkout and invoice events including invoice.paid", events.Data)
	}
	for _, event := range events.Data {
		if event.Type != "customer.subscription.created" {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal(event.Data.Object, &payload); err != nil {
			t.Fatalf("decode subscription webhook payload: %v", err)
		}
		if _, ok := payload["current_period_start"].(float64); !ok {
			t.Fatalf("subscription webhook current_period_start = %#v, want Stripe-style unix timestamp", payload["current_period_start"])
		}
	}

	attempts := getJSON[struct {
		Object string                     `json:"object"`
		Data   []webhooks.DeliveryAttempt `json:"data"`
	}](t, handler, "/api/delivery-attempts")
	if len(attempts.Data) == 0 {
		t.Fatal("no delivery attempts recorded")
	}
	if signatures[0] == "" || attempts.Data[0].RequestHeaders[webhooks.SignatureHeaderName] == "" {
		t.Fatalf("missing signature: received=%q attempt=%#v", signatures[0], attempts.Data[0].RequestHeaders)
	}

	replay := postJSON[struct {
		Message string                     `json:"message"`
		Data    []webhooks.DeliveryAttempt `json:"data"`
	}](t, handler, "/api/events/"+checkoutEventID+"/replay", map[string]any{
		"duplicate":     2,
		"delay_seconds": 30,
		"out_of_order":  true,
	})
	if len(replay.Data) != 2 {
		t.Fatalf("replay attempts = %d, want duplicate delayed attempts", len(replay.Data))
	}
	if replay.Data[0].Status != webhooks.StatusScheduled || replay.Data[0].Metadata["source"] != webhooks.SourceReplay || replay.Data[0].Metadata["out_of_order"] != "true" {
		t.Fatalf("replay attempt = %#v, want scheduled replay out-of-order metadata", replay.Data[0])
	}

	failure := postJSON[struct {
		Message string                     `json:"message"`
		Data    []webhooks.DeliveryAttempt `json:"data"`
	}](t, handler, "/api/events/"+checkoutEventID+"/replay", map[string]any{
		"response_status": 500,
		"response_body":   "receiver down",
	})
	if len(failure.Data) != 1 {
		t.Fatalf("failure replay attempts = %d, want 1 failed attempt", len(failure.Data))
	}
	failed := failure.Data[0]
	if failed.Status != webhooks.StatusFailed || failed.ResponseStatus != 500 || failed.ResponseBody != "receiver down" || failed.NextRetryAt == nil {
		t.Fatalf("failure replay attempt = %#v, want failed 500 with retry evidence", failed)
	}
	if failed.Metadata["response_status"] != "500" {
		t.Fatalf("failure replay metadata = %#v, want response_status evidence", failed.Metadata)
	}

	timeout := postJSON[struct {
		Message string                     `json:"message"`
		Data    []webhooks.DeliveryAttempt `json:"data"`
	}](t, handler, "/api/events/"+invoiceCreatedEventID+"/replay", map[string]any{
		"timeout": true,
	})
	if len(timeout.Data) != 1 || timeout.Data[0].Status != webhooks.StatusFailed || !strings.Contains(timeout.Data[0].Error, "timeout") || timeout.Data[0].NextRetryAt == nil {
		t.Fatalf("timeout replay = %#v, want failed timeout with retry evidence", timeout.Data)
	}

	signatureMismatch := postJSON[struct {
		Message string                     `json:"message"`
		Data    []webhooks.DeliveryAttempt `json:"data"`
	}](t, handler, "/api/events/"+invoiceFinalizedEventID+"/replay", map[string]any{
		"signature_mismatch": true,
		"response_status":    400,
		"response_body":      "bad signature",
	})
	if len(signatureMismatch.Data) != 1 || signatureMismatch.Data[0].Metadata["signature_mismatch"] != "true" {
		t.Fatalf("signature replay = %#v, want signature mismatch evidence", signatureMismatch.Data)
	}
	if signature, _ := signatureMismatch.Data[0].RequestHeaders[webhooks.SignatureHeaderName]; !strings.Contains(signature, "v1=****") {
		t.Fatalf("signature replay header = %q, want masked signature evidence", signature)
	}
}

func TestWebhookEndpointHistoricalReplay(t *testing.T) {
	handler := newTestHandler(t)
	since := time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"history@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"History"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	_ = postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})
	initialAttempts := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/api/delivery-attempts")
	if len(initialAttempts.Data) != 0 {
		t.Fatalf("initial attempts = %#v, want no attempts before endpoint registration", initialAttempts.Data)
	}

	var receivedTypes []string
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Type string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode webhook payload: %v", err)
		}
		receivedTypes = append(receivedTypes, payload.Type)
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()

	endpoint := postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL},
		"enabled_events": {"invoice.*"},
	})

	laterSession := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	_ = postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+laterSession.ID+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})
	if len(receivedTypes) != 4 {
		t.Fatalf("post-registration received types = %#v, want four normally delivered invoice events", receivedTypes)
	}

	catchup := postForm[struct {
		Object         string           `json:"object"`
		EndpointID     string           `json:"endpoint_id"`
		MatchedEvents  int              `json:"matched_events"`
		ReplayedEvents int              `json:"replayed_events"`
		SkippedEvents  int              `json:"skipped_events"`
		AttemptCount   int              `json:"attempt_count"`
		Events         []webhooks.Event `json:"events"`
		Data           []struct {
			EventType string            `json:"event_type"`
			Status    string            `json:"status"`
			Metadata  map[string]string `json:"metadata"`
		} `json:"data"`
	}](t, handler, "/api/webhooks/endpoints/"+endpoint.ID+"/replay-historical?since="+url.QueryEscape(since), url.Values{})
	if catchup.Object != "historical_webhook_replay" || catchup.EndpointID != endpoint.ID {
		t.Fatalf("catchup identity = %#v, want historical replay for endpoint", catchup)
	}
	if catchup.MatchedEvents != 4 || catchup.ReplayedEvents != 4 || catchup.SkippedEvents != 0 || catchup.AttemptCount != 4 || len(catchup.Data) != 4 || len(catchup.Events) != 4 {
		t.Fatalf("catchup = %#v, want four invoice events replayed", catchup)
	}
	for _, attempt := range catchup.Data {
		if !strings.HasPrefix(attempt.EventType, "invoice.") || attempt.Status != webhooks.StatusSucceeded || attempt.Metadata["source"] != webhooks.SourceReplay || attempt.Metadata["historical"] != "true" {
			t.Fatalf("catchup attempt = %#v, want successful historical invoice replay", attempt)
		}
	}
	if len(receivedTypes) != 8 {
		t.Fatalf("received types = %#v, want four normal and four historical invoice events", receivedTypes)
	}

	second := postForm[struct {
		MatchedEvents  int `json:"matched_events"`
		ReplayedEvents int `json:"replayed_events"`
		SkippedEvents  int `json:"skipped_events"`
		AttemptCount   int `json:"attempt_count"`
	}](t, handler, "/api/webhooks/endpoints/"+endpoint.ID+"/replay-historical?since="+url.QueryEscape(since), url.Values{})
	if second.MatchedEvents != 4 || second.ReplayedEvents != 0 || second.SkippedEvents != 4 || second.AttemptCount != 0 {
		t.Fatalf("second catchup = %#v, want idempotent skip of existing attempts", second)
	}
	if len(receivedTypes) != 8 {
		t.Fatalf("received types after second catchup = %#v, want no duplicate delivery", receivedTypes)
	}

	filtered := postForm[struct {
		MatchedEvents  int `json:"matched_events"`
		ReplayedEvents int `json:"replayed_events"`
		AttemptCount   int `json:"attempt_count"`
	}](t, handler, "/api/webhooks/endpoints/"+endpoint.ID+"/replay-historical?since="+url.QueryEscape(since)+"&type=invoice.paid&force=true", url.Values{})
	if filtered.MatchedEvents != 1 || filtered.ReplayedEvents != 1 || filtered.AttemptCount != 1 {
		t.Fatalf("filtered catchup = %#v, want forced invoice.paid replay", filtered)
	}
}

func TestProductionBoundaryRedactionAndAuditAPI(t *testing.T) {
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()

	handler := newTestHandler(t)
	endpoint := postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL + "?api_key=stripe_key_redaction_sample"},
		"secret":         {"webhook_secret_redaction_sample"},
		"enabled_events": {"*"},
	})
	if endpoint.Secret != security.MaskedValue {
		t.Fatalf("endpoint secret = %q, want masked", endpoint.Secret)
	}

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"buyer@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	_ = postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})

	attempts := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/api/delivery-attempts")
	if len(attempts.Data) == 0 {
		t.Fatal("no delivery attempts recorded")
	}
	first := attempts.Data[0]
	headers, _ := first["request_headers"].(map[string]any)
	signature, _ := headers[webhooks.SignatureHeaderName].(string)
	if !strings.Contains(signature, "v1=****") {
		t.Fatalf("signature = %q, want masked HMAC evidence", signature)
	}
	if strings.Contains(first["request_url"].(string), "stripe_key_redaction_sample") {
		t.Fatalf("request_url = %q, want sensitive query masked", first["request_url"])
	}

	events := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events")
	if len(events.Data) == 0 {
		t.Fatal("no events recorded")
	}
	replay := postJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/api/events/"+events.Data[0].ID+"/replay", map[string]any{
		"duplicate":     2,
		"delay_seconds": 30,
		"out_of_order":  true,
	})
	if len(replay.Data) != 2 {
		t.Fatalf("replay attempts = %d, want duplicate replay attempts", len(replay.Data))
	}
	replayHeaders, _ := replay.Data[0]["request_headers"].(map[string]any)
	if signature, _ := replayHeaders[webhooks.SignatureHeaderName].(string); !strings.Contains(signature, "v1=****") {
		t.Fatalf("replay signature = %q, want masked HMAC evidence", signature)
	}
	audit := getJSON[struct {
		Data []webhooks.AuditEntry `json:"data"`
	}](t, handler, "/api/audit-log?action=webhook.replay&targetId="+events.Data[0].ID)
	if len(audit.Data) != 1 || audit.Data[0].Metadata["out_of_order"] != "true" {
		t.Fatalf("audit = %#v, want replay audit evidence", audit.Data)
	}

	status, body := postJSONStatus(t, handler, "/v1/checkout/sessions", map[string]any{
		"payment_method_data": map[string]any{
			"card": map[string]any{"number": "4242424242424242"},
		},
	})
	if status != http.StatusBadRequest || !strings.Contains(body, "real card data") {
		t.Fatalf("status=%d body=%s, want real card data rejection", status, body)
	}
	errBody := decodeErrorBody(t, body)
	if errBody.Error.Type != "invalid_request_error" || errBody.Error.Code != "real_card_data_not_allowed" {
		t.Fatalf("error=%#v, want Stripe-like real-card-data error", errBody.Error)
	}
}

func TestStripeLikeAPIErrorEnvelope(t *testing.T) {
	handler := newTestHandler(t)

	t.Run("validation failure", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/products", url.Values{})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Type != "invalid_request_error" || errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "name" {
			t.Fatalf("error=%#v, want missing name invalid_request_error", errBody.Error)
		}
		if errBody.Error.Message != "Missing required param: name." {
			t.Fatalf("message=%q, want structured validation message", errBody.Error.Message)
		}
	})

	t.Run("resource missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/customers/cus_missing", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		errBody := decodeErrorBody(t, rec.Body.String())
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status=%d body=%s, want 404", rec.Code, rec.Body.String())
		}
		if errBody.Error.Type != "invalid_request_error" || errBody.Error.Code != "resource_missing" {
			t.Fatalf("error=%#v, want resource_missing invalid_request_error", errBody.Error)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/v1/products", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		errBody := decodeErrorBody(t, rec.Body.String())
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status=%d body=%s, want 405", rec.Code, rec.Body.String())
		}
		if rec.Header().Get("Allow") != "GET, POST" {
			t.Fatalf("Allow=%q, want GET, POST", rec.Header().Get("Allow"))
		}
		if errBody.Error.Type != "invalid_request_error" || errBody.Error.Code != "method_not_allowed" {
			t.Fatalf("error=%#v, want method_not_allowed invalid_request_error", errBody.Error)
		}
	})

	t.Run("nested parameter name", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer": {"cus_missing"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Type != "invalid_request_error" || errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "line_items" {
			t.Fatalf("error=%#v, want line_items parameter_missing", errBody.Error)
		}
	})

	t.Run("error category mapping", func(t *testing.T) {
		cases := []struct {
			status int
			want   string
		}{
			{status: http.StatusPaymentRequired, want: "card_error"},
			{status: http.StatusTooManyRequests, want: "invalid_request_error"},
			{status: http.StatusInternalServerError, want: "api_error"},
		}
		for _, tc := range cases {
			got := stripeErrorFor(tc.status, errors.New("simulated"))
			if got.Type != tc.want {
				t.Fatalf("status %d error type = %q, want %q", tc.status, got.Type, tc.want)
			}
		}
	})
}

func TestKnownStripeRouteUnsupportedFallback(t *testing.T) {
	handler := newTestHandler(t)

	status, body := postFormStatusWithHeaders(t, handler, "/v1/payment_intents/pi_123/apply_customer_balance", url.Values{}, map[string]string{
		"Request-Id": "req_known_unsupported",
	})
	errBody := decodeErrorBody(t, body)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", status, body)
	}
	if errBody.Error.Type != "invalid_request_error" || errBody.Error.Code != "unsupported_endpoint" {
		t.Fatalf("error=%#v, want unsupported_endpoint invalid_request_error", errBody.Error)
	}
	if !strings.Contains(errBody.Error.Message, "POST /v1/payment_intents/{id}/apply_customer_balance") {
		t.Fatalf("message=%q, want normalized known Stripe route", errBody.Error.Message)
	}

	searchReq := httptest.NewRequest(http.MethodGet, "/v1/customers/search?query=email:'buyer@example.test'", nil)
	searchReq.Header.Set("Request-Id", "req_known_search_unsupported")
	searchRec := httptest.NewRecorder()
	handler.ServeHTTP(searchRec, searchReq)
	searchErr := decodeErrorBody(t, searchRec.Body.String())
	if searchRec.Code != http.StatusBadRequest || searchErr.Error.Code != "unsupported_endpoint" {
		t.Fatalf("customer search status=%d error=%#v, want known-route unsupported", searchRec.Code, searchErr.Error)
	}

	v2Status, v2Body := postFormStatusWithHeaders(t, handler, "/v2/core/accounts", url.Values{}, map[string]string{
		"Request-Id": "req_known_v2_unsupported",
	})
	v2Err := decodeErrorBody(t, v2Body)
	if v2Status != http.StatusBadRequest || v2Err.Error.Code != "unsupported_endpoint" {
		t.Fatalf("v2 account status=%d error=%#v, want known-route unsupported", v2Status, v2Err.Error)
	}

	invalidLimitReq := httptest.NewRequest(http.MethodGet, "/v1/country_specs?limit=not-an-int", nil)
	invalidLimitReq.Header.Set("Request-Id", "req_known_schema_invalid_limit")
	invalidLimitRec := httptest.NewRecorder()
	handler.ServeHTTP(invalidLimitRec, invalidLimitReq)
	invalidLimitErr := decodeErrorBody(t, invalidLimitRec.Body.String())
	if invalidLimitRec.Code != http.StatusBadRequest || invalidLimitErr.Error.Code != "parameter_invalid" || invalidLimitErr.Error.Param != "limit" {
		t.Fatalf("country specs status=%d error=%#v, want OpenAPI-backed invalid limit", invalidLimitRec.Code, invalidLimitErr.Error)
	}

	unknownParamReq := httptest.NewRequest(http.MethodGet, "/v1/country_specs?nickname=legacy", nil)
	unknownParamRec := httptest.NewRecorder()
	handler.ServeHTTP(unknownParamRec, unknownParamReq)
	unknownParamErr := decodeErrorBody(t, unknownParamRec.Body.String())
	if unknownParamRec.Code != http.StatusBadRequest || unknownParamErr.Error.Code != "parameter_unknown" || unknownParamErr.Error.Param != "nickname" {
		t.Fatalf("country specs status=%d error=%#v, want OpenAPI-backed unknown parameter", unknownParamRec.Code, unknownParamErr.Error)
	}

	missingBodyStatus, missingBodyBody := postFormStatus(t, handler, "/v1/apps/secrets", url.Values{
		"payload":     {"secret"},
		"scope[type]": {"account"},
	})
	missingBodyErr := decodeErrorBody(t, missingBodyBody)
	if missingBodyStatus != http.StatusBadRequest || missingBodyErr.Error.Code != "parameter_missing" || missingBodyErr.Error.Param != "name" {
		t.Fatalf("apps secrets status=%d error=%#v, want OpenAPI-backed missing name", missingBodyStatus, missingBodyErr.Error)
	}

	invalidEnumStatus, invalidEnumBody := postFormStatus(t, handler, "/v1/apps/secrets", url.Values{
		"name":        {"token"},
		"payload":     {"secret"},
		"scope[type]": {"workspace"},
	})
	invalidEnumErr := decodeErrorBody(t, invalidEnumBody)
	if invalidEnumStatus != http.StatusBadRequest || invalidEnumErr.Error.Code != "parameter_invalid" || invalidEnumErr.Error.Param != "scope[type]" {
		t.Fatalf("apps secrets status=%d error=%#v, want OpenAPI-backed invalid nested enum", invalidEnumStatus, invalidEnumErr.Error)
	}

	deepNestedStatus, deepNestedBody := postFormStatus(t, handler, "/v1/account_sessions", url.Values{
		"account": {"acct_123"},
		"components[account_onboarding][enabled]": {"maybe"},
	})
	deepNestedErr := decodeErrorBody(t, deepNestedBody)
	if deepNestedStatus != http.StatusBadRequest || deepNestedErr.Error.Code != "parameter_invalid" || deepNestedErr.Error.Param != "components[account_onboarding][enabled]" {
		t.Fatalf("account sessions status=%d error=%#v, want OpenAPI-backed invalid deep nested boolean", deepNestedStatus, deepNestedErr.Error)
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/v1/products", nil)
	patchRec := httptest.NewRecorder()
	handler.ServeHTTP(patchRec, patchReq)
	if patchRec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PATCH /v1/products status=%d body=%s, want existing method_not_allowed", patchRec.Code, patchRec.Body.String())
	}

	unknownReq := httptest.NewRequest(http.MethodGet, "/v1/not_a_stripe_route", nil)
	unknownRec := httptest.NewRecorder()
	handler.ServeHTTP(unknownRec, unknownReq)
	if unknownRec.Code != http.StatusNotFound || strings.Contains(unknownRec.Body.String(), "unsupported_endpoint") {
		t.Fatalf("unknown route status=%d body=%s, want plain not found outside known catalog", unknownRec.Code, unknownRec.Body.String())
	}

	v1PrefixReq := httptest.NewRequest(http.MethodGet, "/v1", nil)
	v1PrefixRec := httptest.NewRecorder()
	handler.ServeHTTP(v1PrefixRec, v1PrefixReq)
	if v1PrefixRec.Code != http.StatusNotFound {
		t.Fatalf("GET /v1 status=%d body=%s, want 404 without ServeMux redirect", v1PrefixRec.Code, v1PrefixRec.Body.String())
	}
	v2PrefixReq := httptest.NewRequest(http.MethodGet, "/v2", nil)
	v2PrefixRec := httptest.NewRecorder()
	handler.ServeHTTP(v2PrefixRec, v2PrefixReq)
	if v2PrefixRec.Code != http.StatusNotFound {
		t.Fatalf("GET /v2 status=%d body=%s, want 404 without ServeMux redirect", v2PrefixRec.Code, v2PrefixRec.Body.String())
	}

	traces := getJSON[struct {
		Object string                     `json:"object"`
		Data   []diagnostics.RequestTrace `json:"data"`
	}](t, handler, "/api/request-traces?requestId=req_known_unsupported")
	if traces.Object != "list" || len(traces.Data) != 1 {
		t.Fatalf("request traces = %#v, want unsupported request trace", traces)
	}
	trace := traces.Data[0]
	if trace.Status != http.StatusBadRequest || trace.ErrorCode != "unsupported_endpoint" || trace.Path != "/v1/payment_intents/pi_123/apply_customer_balance" {
		t.Fatalf("trace = %#v, want unsupported endpoint evidence", trace)
	}

	validationTraces := getJSON[struct {
		Object string                     `json:"object"`
		Data   []diagnostics.RequestTrace `json:"data"`
	}](t, handler, "/api/request-traces?requestId=req_known_schema_invalid_limit")
	if validationTraces.Object != "list" || len(validationTraces.Data) != 1 {
		t.Fatalf("validation request traces = %#v, want schema-validation request trace", validationTraces)
	}
	if validationTraces.Data[0].Path != "/v1/country_specs" || validationTraces.Data[0].ErrorCode != "parameter_invalid" || validationTraces.Data[0].ErrorParam != "limit" {
		t.Fatalf("validation trace = %#v, want OpenAPI schema validation evidence", validationTraces.Data[0])
	}

	searchTraces := getJSON[struct {
		Object string                     `json:"object"`
		Data   []diagnostics.RequestTrace `json:"data"`
	}](t, handler, "/api/request-traces?requestId=req_known_search_unsupported")
	if searchTraces.Object != "list" || len(searchTraces.Data) != 1 {
		t.Fatalf("search request traces = %#v, want unsupported request trace", searchTraces)
	}
	if searchTraces.Data[0].Path != "/v1/customers/search" || searchTraces.Data[0].ErrorCode != "unsupported_endpoint" {
		t.Fatalf("search trace = %#v, want unsupported search route evidence", searchTraces.Data[0])
	}

	v2Traces := getJSON[struct {
		Object string                     `json:"object"`
		Data   []diagnostics.RequestTrace `json:"data"`
	}](t, handler, "/api/request-traces?requestId=req_known_v2_unsupported")
	if v2Traces.Object != "list" || len(v2Traces.Data) != 1 {
		t.Fatalf("v2 request traces = %#v, want unsupported request trace", v2Traces)
	}
	if v2Traces.Data[0].Path != "/v2/core/accounts" || v2Traces.Data[0].ErrorCode != "unsupported_endpoint" {
		t.Fatalf("v2 trace = %#v, want unsupported v2 endpoint evidence", v2Traces.Data[0])
	}

}

func TestConnectAccountSmokeCompatibility(t *testing.T) {
	handler := newTestHandlerWithOptions(t, Options{PublicBaseURL: "http://127.0.0.1:18080"})

	status, body := postFormStatusWithHeaders(t, handler, "/v1/accounts", url.Values{
		"type":                                   {"express"},
		"country":                                {"US"},
		"email":                                  {"seller@example.test"},
		"business_type":                          {"company"},
		"capabilities[card_payments][requested]": {"true"},
		"capabilities[transfers][requested]":     {"true"},
		"metadata[platform]":                     {"sample"},
	}, map[string]string{
		"Request-Id":     "req_connect_account_create",
		"Stripe-Account": "acct_platform_trace",
	})
	if status != http.StatusOK {
		t.Fatalf("account create status=%d body=%s", status, body)
	}
	var account struct {
		ID              string            `json:"id"`
		Object          string            `json:"object"`
		Type            string            `json:"type"`
		Country         string            `json:"country"`
		Email           string            `json:"email"`
		ChargesEnabled  bool              `json:"charges_enabled"`
		PayoutsEnabled  bool              `json:"payouts_enabled"`
		Capabilities    map[string]string `json:"capabilities"`
		Metadata        map[string]string `json:"metadata"`
		DefaultCurrency string            `json:"default_currency"`
	}
	if err := json.Unmarshal([]byte(body), &account); err != nil {
		t.Fatalf("decode account: %v", err)
	}
	if !strings.HasPrefix(account.ID, "acct_") || account.Object != "account" || account.Type != "express" || account.Country != "US" || !account.ChargesEnabled || !account.PayoutsEnabled {
		t.Fatalf("account = %#v, want Stripe-shaped connected account", account)
	}
	if account.Capabilities["card_payments"] != "active" || account.Capabilities["transfers"] != "active" || account.Metadata["platform"] != "sample" || account.DefaultCurrency != "usd" {
		t.Fatalf("account capabilities/metadata = %#v/%#v currency=%q", account.Capabilities, account.Metadata, account.DefaultCurrency)
	}

	retrieved := getJSON[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/accounts/"+account.ID)
	if retrieved.ID != account.ID {
		t.Fatalf("retrieved account = %#v, want %s", retrieved, account.ID)
	}

	updated := postForm[struct {
		ID           string            `json:"id"`
		Capabilities map[string]string `json:"capabilities"`
		Metadata     map[string]string `json:"metadata"`
	}](t, handler, "/v1/accounts/"+account.ID, url.Values{
		"capabilities[card_payments][status]": {"pending"},
		"metadata[stage]":                     {"onboarding"},
	})
	if updated.ID != account.ID || updated.Capabilities["card_payments"] != "pending" || updated.Capabilities["transfers"] != "active" || updated.Metadata["platform"] != "sample" || updated.Metadata["stage"] != "onboarding" {
		t.Fatalf("updated account = %#v, want merged capability and metadata", updated)
	}

	list := getJSON[struct {
		Object string `json:"object"`
		Data   []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/accounts?country=US&type=express")
	if list.Object != "list" || len(list.Data) != 1 || list.Data[0].ID != account.ID {
		t.Fatalf("accounts list = %#v, want created account", list)
	}

	link := postForm[struct {
		Object  string `json:"object"`
		Account string `json:"account"`
		Type    string `json:"type"`
		URL     string `json:"url"`
	}](t, handler, "/v1/account_links", url.Values{
		"account":     {account.ID},
		"type":        {"account_onboarding"},
		"refresh_url": {"http://app.example.test/refresh"},
		"return_url":  {"http://app.example.test/return"},
	})
	if link.Object != "account_link" || link.Account != account.ID || link.Type != "account_onboarding" || !strings.HasPrefix(link.URL, "http://127.0.0.1:18080/connect/accounts/"+account.ID) {
		t.Fatalf("account link = %#v, want local onboarding link", link)
	}

	session := postForm[struct {
		Object       string `json:"object"`
		Account      string `json:"account"`
		ClientSecret string `json:"client_secret"`
		ExpiresAt    int64  `json:"expires_at"`
		Components   map[string]struct {
			Enabled bool `json:"enabled"`
		} `json:"components"`
	}](t, handler, "/v1/account_sessions", url.Values{
		"account": {account.ID},
		"components[account_onboarding][enabled]":  {"true"},
		"components[notification_banner][enabled]": {"false"},
	})
	if session.Object != "account_session" || session.Account != account.ID || session.ClientSecret == "" || session.ExpiresAt == 0 || !session.Components["account_onboarding"].Enabled || session.Components["notification_banner"].Enabled {
		t.Fatalf("account session = %#v, want component flags, expiry, and client secret", session)
	}

	missingLinkTypeStatus, missingLinkTypeBody := postFormStatus(t, handler, "/v1/account_links", url.Values{
		"account":     {account.ID},
		"refresh_url": {"http://app.example.test/refresh"},
		"return_url":  {"http://app.example.test/return"},
	})
	missingLinkTypeErr := decodeErrorBody(t, missingLinkTypeBody)
	if missingLinkTypeStatus != http.StatusBadRequest || missingLinkTypeErr.Error.Code != "parameter_missing" || missingLinkTypeErr.Error.Param != "type" {
		t.Fatalf("account link missing type status=%d error=%#v, want missing type", missingLinkTypeStatus, missingLinkTypeErr.Error)
	}

	missingSessionComponentsStatus, missingSessionComponentsBody := postFormStatus(t, handler, "/v1/account_sessions", url.Values{
		"account": {account.ID},
	})
	missingSessionComponentsErr := decodeErrorBody(t, missingSessionComponentsBody)
	if missingSessionComponentsStatus != http.StatusBadRequest || missingSessionComponentsErr.Error.Code != "parameter_missing" || missingSessionComponentsErr.Error.Param != "components" {
		t.Fatalf("account session missing components status=%d error=%#v, want missing components", missingSessionComponentsStatus, missingSessionComponentsErr.Error)
	}

	traces := getJSON[struct {
		Data []diagnostics.RequestTrace `json:"data"`
	}](t, handler, "/api/request-traces?requestId=req_connect_account_create")
	if len(traces.Data) != 1 || traces.Data[0].RequestHeaders["Stripe-Account"] != "acct_platform_trace" || traces.Data[0].ResponseObject != "account" || traces.Data[0].ResponseObjectID != account.ID {
		t.Fatalf("connect trace = %#v, want Stripe-Account routing evidence", traces.Data)
	}
}

func TestConnectPlatformResourceCompatibility(t *testing.T) {
	handler := newTestHandlerWithOptions(t, Options{PublicBaseURL: "http://127.0.0.1:18080"})

	platformAccount := getJSON[struct {
		ID             string            `json:"id"`
		Object         string            `json:"object"`
		Type           string            `json:"type"`
		ChargesEnabled bool              `json:"charges_enabled"`
		PayoutsEnabled bool              `json:"payouts_enabled"`
		Capabilities   map[string]string `json:"capabilities"`
		Metadata       map[string]string `json:"metadata"`
	}](t, handler, "/v1/account")
	if platformAccount.ID != "acct_billtap_platform" || platformAccount.Object != "account" || platformAccount.Type != "standard" || !platformAccount.ChargesEnabled || !platformAccount.PayoutsEnabled {
		t.Fatalf("platform account = %#v, want deterministic local platform account", platformAccount)
	}
	if platformAccount.Capabilities["card_payments"] != "active" || platformAccount.Capabilities["transfers"] != "active" || platformAccount.Metadata["billtap_account"] != "platform" {
		t.Fatalf("platform account capabilities/metadata = %#v/%#v", platformAccount.Capabilities, platformAccount.Metadata)
	}

	account := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/accounts", url.Values{
		"type":    {"express"},
		"country": {"US"},
		"email":   {"platform-resource@example.test"},
	})

	person := postForm[struct {
		ID           string            `json:"id"`
		Object       string            `json:"object"`
		Account      string            `json:"account"`
		FirstName    string            `json:"first_name"`
		LastName     string            `json:"last_name"`
		Email        any               `json:"email"`
		DOB          map[string]int64  `json:"dob"`
		Relationship map[string]any    `json:"relationship"`
		Metadata     map[string]string `json:"metadata"`
		Verification map[string]any    `json:"verification"`
	}](t, handler, "/v1/accounts/"+account.ID+"/people", url.Values{
		"first_name":                       {"Ada"},
		"last_name":                        {"Lovelace"},
		"dob[day]":                         {"10"},
		"dob[month]":                       {"12"},
		"dob[year]":                        {"1815"},
		"relationship[owner]":              {"true"},
		"relationship[representative]":     {"1"},
		"relationship[percent_ownership]":  {"75"},
		"address[line1]":                   {"1 Platform Way"},
		"address[country]":                 {"US"},
		"metadata[platformRepresentative]": {"true"},
	})
	if !strings.HasPrefix(person.ID, "person_") || person.Object != "person" || person.Account != account.ID || person.FirstName != "Ada" || person.LastName != "Lovelace" {
		t.Fatalf("person = %#v, want local Connect person evidence", person)
	}
	if person.DOB["day"] != 10 || person.DOB["month"] != 12 || person.DOB["year"] != 1815 || person.Relationship["owner"] != true || person.Relationship["representative"] != true || person.Relationship["percent_ownership"] != float64(75) {
		t.Fatalf("person dob/relationship = %#v/%#v", person.DOB, person.Relationship)
	}
	if person.Metadata["platformRepresentative"] != "true" || person.Verification["status"] != "unverified" {
		t.Fatalf("person metadata/verification = %#v/%#v", person.Metadata, person.Verification)
	}
	peopleList := getJSON[struct {
		Object string `json:"object"`
		Data   []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/accounts/"+account.ID+"/people")
	if peopleList.Object != "list" || len(peopleList.Data) != 1 || peopleList.Data[0].ID != person.ID {
		t.Fatalf("people list = %#v, want created person", peopleList)
	}
	personAlias := getJSON[struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Account string `json:"account"`
	}](t, handler, "/v1/accounts/"+account.ID+"/persons/"+person.ID)
	if personAlias.ID != person.ID || personAlias.Object != "person" || personAlias.Account != account.ID {
		t.Fatalf("person alias = %#v, want persons alias to retrieve person", personAlias)
	}
	updatedPerson := postForm[struct {
		ID           string            `json:"id"`
		Email        string            `json:"email"`
		Relationship map[string]any    `json:"relationship"`
		Metadata     map[string]string `json:"metadata"`
	}](t, handler, "/v1/accounts/"+account.ID+"/people/"+person.ID, url.Values{
		"email":               {"ada@example.test"},
		"relationship[title]": {"Founder"},
		"metadata[stage]":     {"verified-by-test"},
	})
	if updatedPerson.ID != person.ID || updatedPerson.Email != "ada@example.test" || updatedPerson.Relationship["title"] != "Founder" || updatedPerson.Metadata["platformRepresentative"] != "true" || updatedPerson.Metadata["stage"] != "verified-by-test" {
		t.Fatalf("updated person = %#v, want merged metadata and person data", updatedPerson)
	}
	deletedPerson := deleteJSON[struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}](t, handler, "/v1/accounts/"+account.ID+"/persons/"+person.ID)
	if deletedPerson.ID != person.ID || deletedPerson.Object != "person" || !deletedPerson.Deleted {
		t.Fatalf("deleted person = %#v, want deleted marker", deletedPerson)
	}
	deletedPersonGetReq := httptest.NewRequest(http.MethodGet, "/v1/accounts/"+account.ID+"/people/"+person.ID, nil)
	deletedPersonGetRec := httptest.NewRecorder()
	handler.ServeHTTP(deletedPersonGetRec, deletedPersonGetReq)
	if deletedPersonGetRec.Code != http.StatusNotFound {
		t.Fatalf("deleted person get status=%d body=%s, want 404", deletedPersonGetRec.Code, deletedPersonGetRec.Body.String())
	}
	emptyPeopleList := getJSON[struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/accounts/"+account.ID+"/persons")
	if len(emptyPeopleList.Data) != 0 {
		t.Fatalf("persons list after delete = %#v, want deleted person filtered out", emptyPeopleList)
	}

	capability := postForm[struct {
		ID        string `json:"id"`
		Object    string `json:"object"`
		Account   string `json:"account"`
		Requested bool   `json:"requested"`
		Status    string `json:"status"`
	}](t, handler, "/v1/accounts/"+account.ID+"/capabilities/transfers", url.Values{
		"requested": {"false"},
	})
	if capability.ID != "transfers" || capability.Object != "capability" || capability.Account != account.ID || capability.Requested || capability.Status != "inactive" {
		t.Fatalf("capability = %#v, want inactive transfers capability", capability)
	}

	external := postForm[struct {
		ID                 string            `json:"id"`
		Object             string            `json:"object"`
		Account            string            `json:"account"`
		Country            string            `json:"country"`
		Currency           string            `json:"currency"`
		Last4              string            `json:"last4"`
		DefaultForCurrency bool              `json:"default_for_currency"`
		Metadata           map[string]string `json:"metadata"`
	}](t, handler, "/v1/accounts/"+account.ID+"/external_accounts", url.Values{
		"account_number":       {"000123456789"},
		"routing_number":       {"110000000"},
		"country":              {"US"},
		"currency":             {"usd"},
		"default_for_currency": {"true"},
		"metadata[purpose]":    {"payouts"},
	})
	if !strings.HasPrefix(external.ID, "ba_") || external.Object != "bank_account" || external.Account != account.ID || external.Country != "US" || external.Currency != "usd" || external.Last4 != "6789" || !external.DefaultForCurrency || external.Metadata["purpose"] != "payouts" {
		t.Fatalf("external account = %#v, want local bank account evidence", external)
	}
	externalList := getJSON[struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/accounts/"+account.ID+"/external_accounts")
	if len(externalList.Data) != 1 || externalList.Data[0].ID != external.ID {
		t.Fatalf("external account list = %#v, want created account", externalList)
	}

	transfer := postForm[struct {
		ID          string            `json:"id"`
		Object      string            `json:"object"`
		Amount      int64             `json:"amount"`
		Currency    string            `json:"currency"`
		Destination string            `json:"destination"`
		Metadata    map[string]string `json:"metadata"`
	}](t, handler, "/v1/transfers", url.Values{
		"amount":            {"2500"},
		"currency":          {"usd"},
		"destination":       {account.ID},
		"description":       {"Local platform transfer"},
		"transfer_group":    {"group-platform-smoke"},
		"metadata[orderId]": {"order_123"},
	})
	if !strings.HasPrefix(transfer.ID, "tr_") || transfer.Object != "transfer" || transfer.Amount != 2500 || transfer.Currency != "usd" || transfer.Destination != account.ID || transfer.Metadata["orderId"] != "order_123" {
		t.Fatalf("transfer = %#v, want local transfer", transfer)
	}

	reversal := postForm[struct {
		ID       string `json:"id"`
		Object   string `json:"object"`
		Amount   int64  `json:"amount"`
		Transfer string `json:"transfer"`
	}](t, handler, "/v1/transfers/"+transfer.ID+"/reversals", url.Values{
		"amount":                 {"500"},
		"refund_application_fee": {"true"},
	})
	if !strings.HasPrefix(reversal.ID, "trr_") || reversal.Object != "transfer_reversal" || reversal.Amount != 500 || reversal.Transfer != transfer.ID {
		t.Fatalf("transfer reversal = %#v, want reversal evidence", reversal)
	}

	payout := postForm[struct {
		ID          string `json:"id"`
		Object      string `json:"object"`
		Amount      int64  `json:"amount"`
		Currency    string `json:"currency"`
		Destination string `json:"destination"`
		Status      string `json:"status"`
	}](t, handler, "/v1/payouts", url.Values{
		"amount":      {"1800"},
		"currency":    {"usd"},
		"destination": {external.ID},
		"method":      {"standard"},
	})
	if !strings.HasPrefix(payout.ID, "po_") || payout.Object != "payout" || payout.Amount != 1800 || payout.Currency != "usd" || payout.Destination != external.ID || payout.Status != "paid" {
		t.Fatalf("payout = %#v, want payout evidence", payout)
	}
	canceledPayout := postForm[struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}](t, handler, "/v1/payouts/"+payout.ID+"/cancel", url.Values{})
	if canceledPayout.ID != payout.ID || canceledPayout.Status != "canceled" {
		t.Fatalf("canceled payout = %#v, want canceled payout", canceledPayout)
	}

	feeRefund := postForm[struct {
		ID       string `json:"id"`
		Object   string `json:"object"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
		Fee      string `json:"fee"`
	}](t, handler, "/v1/application_fees/fee_platform_smoke/refunds", url.Values{
		"amount": {"300"},
	})
	if !strings.HasPrefix(feeRefund.ID, "fr_") || feeRefund.Object != "fee_refund" || feeRefund.Amount != 300 || feeRefund.Currency != "usd" || feeRefund.Fee != "fee_platform_smoke" {
		t.Fatalf("fee refund = %#v, want application fee refund evidence", feeRefund)
	}
	fee := getJSON[struct {
		ID             string `json:"id"`
		Object         string `json:"object"`
		Amount         int64  `json:"amount"`
		AmountRefunded int64  `json:"amount_refunded"`
	}](t, handler, "/v1/application_fees/fee_platform_smoke")
	if fee.ID != "fee_platform_smoke" || fee.Object != "application_fee" || fee.Amount != 300 || fee.AmountRefunded != 300 {
		t.Fatalf("application fee = %#v, want auto-materialized fee evidence", fee)
	}

	rejected := postForm[struct {
		ID               string `json:"id"`
		ChargesEnabled   bool   `json:"charges_enabled"`
		PayoutsEnabled   bool   `json:"payouts_enabled"`
		DetailsSubmitted bool   `json:"details_submitted"`
	}](t, handler, "/v1/accounts/"+account.ID+"/reject", url.Values{
		"reason": {"fraud"},
	})
	if rejected.ID != account.ID || rejected.ChargesEnabled || rejected.PayoutsEnabled || rejected.DetailsSubmitted {
		t.Fatalf("rejected account = %#v, want disabled account evidence", rejected)
	}
	storedRejected := getJSON[struct {
		ID             string `json:"id"`
		ChargesEnabled bool   `json:"charges_enabled"`
		PayoutsEnabled bool   `json:"payouts_enabled"`
	}](t, handler, "/v1/accounts/"+account.ID)
	if storedRejected.ID != account.ID || storedRejected.ChargesEnabled || storedRejected.PayoutsEnabled {
		t.Fatalf("stored rejected account = %#v, want disabled state persisted", storedRejected)
	}
	deletedAccount := deleteJSON[struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}](t, handler, "/v1/accounts/"+account.ID)
	if deletedAccount.ID != account.ID || deletedAccount.Object != "account" || !deletedAccount.Deleted {
		t.Fatalf("deleted account = %#v, want local deletion marker", deletedAccount)
	}

	events := getJSON[struct {
		Data []struct {
			Type string `json:"type"`
		} `json:"data"`
	}](t, handler, "/v1/events")
	seen := map[string]bool{}
	for _, event := range events.Data {
		seen[event.Type] = true
	}
	for _, eventType := range []string{"transfer.created", "transfer.reversed", "payout.created", "payout.canceled", "application_fee.refunded"} {
		if !seen[eventType] {
			t.Fatalf("events missing %s: %#v", eventType, events.Data)
		}
	}
}

func TestSupportedEndpointRequestValidation(t *testing.T) {
	handler := newTestHandler(t)

	t.Run("unknown product parameter", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/products", url.Values{
			"name":     {"Team"},
			"nickname": {"unused"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_unknown" || errBody.Error.Param != "nickname" {
			t.Fatalf("error=%#v, want unknown nickname", errBody.Error)
		}
	})

	t.Run("unknown customer update parameter", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/customers/cus_missing", url.Values{
			"nickname": {"unused"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_unknown" || errBody.Error.Param != "nickname" {
			t.Fatalf("error=%#v, want unknown nickname", errBody.Error)
		}
	})

	t.Run("invalid price amount type", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/prices", url.Values{
			"product":     {"prod_missing"},
			"currency":    {"usd"},
			"unit_amount": {"not-an-int"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "unit_amount" {
			t.Fatalf("error=%#v, want invalid unit_amount", errBody.Error)
		}
	})

	t.Run("invalid JSON price amount type", func(t *testing.T) {
		status, body := postJSONStatus(t, handler, "/v1/prices", map[string]any{
			"product":     "prod_missing",
			"currency":    "usd",
			"unit_amount": 9.99,
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "unit_amount" {
			t.Fatalf("error=%#v, want invalid unit_amount", errBody.Error)
		}
	})

	t.Run("invalid price interval enum", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/prices", url.Values{
			"product":             {"prod_missing"},
			"currency":            {"usd"},
			"unit_amount":         {"9900"},
			"recurring[interval]": {"decade"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "recurring[interval]" {
			t.Fatalf("error=%#v, want invalid recurring interval", errBody.Error)
		}
	})

	t.Run("invalid price update active type", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/prices/price_missing", url.Values{
			"active": {"maybe"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "active" {
			t.Fatalf("error=%#v, want invalid active", errBody.Error)
		}
	})

	t.Run("missing price product", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/prices", url.Values{
			"currency":    {"usd"},
			"unit_amount": {"9900"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "product" {
			t.Fatalf("error=%#v, want missing product", errBody.Error)
		}
	})

	t.Run("price product must exist", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/prices", url.Values{
			"product":     {"prod_missing"},
			"currency":    {"usd"},
			"unit_amount": {"9900"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusNotFound {
			t.Fatalf("status=%d body=%s, want 404", status, body)
		}
		if errBody.Error.Code != "resource_missing" {
			t.Fatalf("error=%#v, want resource_missing", errBody.Error)
		}
	})

	t.Run("price search missing query", func(t *testing.T) {
		status, body := getStatus(t, handler, "/v1/prices/search")
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "query" {
			t.Fatalf("error=%#v, want missing query", errBody.Error)
		}
	})

	t.Run("price search unsupported clause", func(t *testing.T) {
		status, body := getStatus(t, handler, "/v1/prices/search?"+url.Values{
			"query": {"currency:'usd'"},
		}.Encode())
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "query" {
			t.Fatalf("error=%#v, want invalid query", errBody.Error)
		}
	})

	t.Run("price search invalid limit", func(t *testing.T) {
		status, body := getStatus(t, handler, "/v1/prices/search?"+url.Values{
			"query": {"active:'true'"},
			"limit": {"zero"},
		}.Encode())
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "limit" {
			t.Fatalf("error=%#v, want invalid limit", errBody.Error)
		}
	})

	t.Run("checkout invalid mode", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer":             {"cus_missing"},
			"mode":                 {"payment"},
			"line_items[0][price]": {"price_missing"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "mode" {
			t.Fatalf("error=%#v, want invalid mode", errBody.Error)
		}
	})

	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Validated Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"validator@example.test"}})

	t.Run("payment intent create requires amount", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/payment_intents", url.Values{
			"currency": {"usd"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "amount" {
			t.Fatalf("status=%d error=%#v, want missing amount", status, errBody.Error)
		}
	})

	t.Run("payment intent create validates amount", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/payment_intents", url.Values{
			"amount":   {"0"},
			"currency": {"usd"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "amount" {
			t.Fatalf("status=%d error=%#v, want invalid amount", status, errBody.Error)
		}
	})

	t.Run("payment intent create validates capture method", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/payment_intents", url.Values{
			"amount":         {"4900"},
			"currency":       {"usd"},
			"capture_method": {"later"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "capture_method" {
			t.Fatalf("status=%d error=%#v, want invalid capture_method", status, errBody.Error)
		}
	})

	t.Run("payment intent create confirm requires payment method or outcome", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/payment_intents", url.Values{
			"amount":   {"4900"},
			"currency": {"usd"},
			"confirm":  {"true"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "payment_method" {
			t.Fatalf("status=%d error=%#v, want missing payment_method", status, errBody.Error)
		}
	})

	t.Run("payment intent create confirm without customer default requires payment method or outcome", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/payment_intents", url.Values{
			"amount":   {"4900"},
			"currency": {"usd"},
			"customer": {customer.ID},
			"confirm":  {"true"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "payment_method" {
			t.Fatalf("status=%d error=%#v, want missing payment_method", status, errBody.Error)
		}
	})

	t.Run("payment intent confirm requires payment method or outcome", func(t *testing.T) {
		intent := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents", url.Values{
			"amount":   {"4900"},
			"currency": {"usd"},
		})
		status, body := postFormStatus(t, handler, "/v1/payment_intents/"+intent.ID+"/confirm", url.Values{})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "payment_method" {
			t.Fatalf("status=%d error=%#v, want missing payment_method", status, errBody.Error)
		}
	})

	t.Run("payment intent confirm rejects unknown parameter", func(t *testing.T) {
		intent := postForm[billing.PaymentIntent](t, handler, "/v1/payment_intents", url.Values{
			"amount":   {"4900"},
			"currency": {"usd"},
		})
		status, body := postFormStatus(t, handler, "/v1/payment_intents/"+intent.ID+"/confirm", url.Values{
			"nickname": {"unused"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_unknown" || errBody.Error.Param != "nickname" {
			t.Fatalf("status=%d error=%#v, want unknown nickname", status, errBody.Error)
		}
	})

	t.Run("setup intent create validates usage", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/setup_intents", url.Values{
			"usage": {"recurring"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "usage" {
			t.Fatalf("status=%d error=%#v, want invalid usage", status, errBody.Error)
		}
	})

	t.Run("setup intent create confirm requires payment method or outcome", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/setup_intents", url.Values{
			"confirm": {"true"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest || errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "payment_method" {
			t.Fatalf("status=%d error=%#v, want missing payment_method", status, errBody.Error)
		}
	})

	t.Run("checkout line item price required", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer":                {customer.ID},
			"line_items[0][quantity]": {"2"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "line_items[0][price]" {
			t.Fatalf("error=%#v, want missing line item price", errBody.Error)
		}
	})

	t.Run("checkout customer must exist", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer":             {"cus_missing"},
			"line_items[0][price]": {price.ID},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusNotFound {
			t.Fatalf("status=%d body=%s, want 404", status, body)
		}
		if errBody.Error.Code != "resource_missing" {
			t.Fatalf("error=%#v, want resource_missing", errBody.Error)
		}
	})

	t.Run("checkout price must exist", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer":             {customer.ID},
			"line_items[0][price]": {"price_missing"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusNotFound {
			t.Fatalf("status=%d body=%s, want 404", status, body)
		}
		if errBody.Error.Code != "resource_missing" {
			t.Fatalf("error=%#v, want resource_missing", errBody.Error)
		}
	})

	t.Run("checkout accepts Stripe SDK promotion and trial params", func(t *testing.T) {
		session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
			"customer":                             {customer.ID},
			"line_items[0][price]":                 {price.ID},
			"allow_promotion_codes":                {"true"},
			"subscription_data[trial_period_days]": {"14"},
		})
		if !session.AllowPromotionCodes || session.TrialPeriodDays != 14 {
			t.Fatalf("session = %#v, want promotion codes and 14-day trial", session)
		}

		completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
			"outcome": "payment_succeeded",
		})
		var completed billing.CheckoutSession
		if err := json.Unmarshal(completion["session"], &completed); err != nil {
			t.Fatalf("decode completed session: %v", err)
		}
		if completed.PaymentStatus != "no_payment_required" {
			t.Fatalf("completed payment_status = %q, want no_payment_required", completed.PaymentStatus)
		}
		var sub billing.Subscription
		if err := json.Unmarshal(completion["subscription"], &sub); err != nil {
			t.Fatalf("decode subscription: %v", err)
		}
		if sub.Status != "trialing" || sub.Metadata["trial_period_days"] != "14" {
			t.Fatalf("subscription = %#v, want trialing with trial metadata", sub)
		}
	})

	t.Run("checkout promotion flag must be boolean", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer":              {customer.ID},
			"line_items[0][price]":  {price.ID},
			"allow_promotion_codes": {"sometimes"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "allow_promotion_codes" {
			t.Fatalf("error=%#v, want invalid allow_promotion_codes", errBody.Error)
		}
	})

	t.Run("checkout trial period must be positive integer", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer":                             {customer.ID},
			"line_items[0][price]":                 {price.ID},
			"subscription_data[trial_period_days]": {"0"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "subscription_data[trial_period_days]" {
			t.Fatalf("error=%#v, want invalid trial_period_days", errBody.Error)
		}
	})

	t.Run("checkout quantity must be positive", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
			"customer":                {customer.ID},
			"line_items[0][price]":    {price.ID},
			"line_items[0][quantity]": {"0"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "line_items[0][quantity]" {
			t.Fatalf("error=%#v, want invalid line item quantity", errBody.Error)
		}
	})

	t.Run("subscription create quantity must be positive", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/subscriptions", url.Values{
			"customer":           {customer.ID},
			"items[0][price]":    {price.ID},
			"items[0][quantity]": {"0"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "items[0][quantity]" {
			t.Fatalf("error=%#v, want invalid subscription item quantity", errBody.Error)
		}
	})

	t.Run("subscription update cancel flag must be boolean", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/subscriptions/sub_missing", url.Values{
			"cancel_at_period_end": {"maybe"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "cancel_at_period_end" {
			t.Fatalf("error=%#v, want invalid cancel_at_period_end", errBody.Error)
		}
	})

	t.Run("subscription update billing cycle anchor must be now or timestamp", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/subscriptions/sub_missing", url.Values{
			"billing_cycle_anchor": {"soon"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "billing_cycle_anchor" {
			t.Fatalf("error=%#v, want invalid billing_cycle_anchor", errBody.Error)
		}
	})

	t.Run("subscription item create quantity must be positive", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/subscription_items", url.Values{
			"subscription": {"sub_missing"},
			"price":        {price.ID},
			"quantity":     {"0"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "quantity" {
			t.Fatalf("error=%#v, want invalid quantity", errBody.Error)
		}
	})

	t.Run("portal customer required", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/billing_portal/sessions", url.Values{})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "customer" {
			t.Fatalf("error=%#v, want missing customer", errBody.Error)
		}
	})

	t.Run("webhook url required", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/webhook_endpoints", url.Values{
			"enabled_events": {"invoice.*"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_missing" || errBody.Error.Param != "url" {
			t.Fatalf("error=%#v, want missing url", errBody.Error)
		}
	})

	t.Run("webhook update active must be boolean", func(t *testing.T) {
		status, body := postFormStatus(t, handler, "/v1/webhook_endpoints/we_missing", url.Values{
			"active": {"maybe"},
		})
		errBody := decodeErrorBody(t, body)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}
		if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "active" {
			t.Fatalf("error=%#v, want invalid active", errBody.Error)
		}
	})
}

func TestIdempotencyKeySimulation(t *testing.T) {
	handler := newTestHandler(t)

	t.Run("matching post replays first response", func(t *testing.T) {
		headers := map[string]string{"Idempotency-Key": "customer-create-same"}
		status, body := postFormStatusWithHeaders(t, handler, "/v1/customers", url.Values{
			"email": {"same@example.test"},
		}, headers)
		if status != http.StatusOK {
			t.Fatalf("status=%d body=%s, want 200", status, body)
		}

		replayStatus, replayBody := postFormStatusWithHeaders(t, handler, "/v1/customers", url.Values{
			"email": {"same@example.test"},
		}, headers)
		if replayStatus != http.StatusOK || replayBody != body {
			t.Fatalf("replay status=%d body=%s, want same body %s", replayStatus, replayBody, body)
		}

		customers := getJSON[struct {
			Data []billing.Customer `json:"data"`
		}](t, handler, "/v1/customers?email=same@example.test")
		if len(customers.Data) != 1 {
			t.Fatalf("customers = %#v, want exactly one created customer", customers.Data)
		}
	})

	t.Run("same key with different params conflicts", func(t *testing.T) {
		headers := map[string]string{"Idempotency-Key": "customer-create-conflict"}
		status, body := postFormStatusWithHeaders(t, handler, "/v1/customers", url.Values{
			"email": {"first@example.test"},
		}, headers)
		if status != http.StatusOK {
			t.Fatalf("status=%d body=%s, want 200", status, body)
		}

		conflictStatus, conflictBody := postFormStatusWithHeaders(t, handler, "/v1/customers", url.Values{
			"email": {"second@example.test"},
		}, headers)
		errBody := decodeErrorBody(t, conflictBody)
		if conflictStatus != http.StatusConflict {
			t.Fatalf("status=%d body=%s, want 409", conflictStatus, conflictBody)
		}
		if errBody.Error.Type != "idempotency_error" || errBody.Error.Code != "idempotency_key_in_use" {
			t.Fatalf("error=%#v, want idempotency conflict", errBody.Error)
		}
	})

	t.Run("validation errors are not cached", func(t *testing.T) {
		headers := map[string]string{"Idempotency-Key": "product-validation-not-cached"}
		status, body := postFormStatusWithHeaders(t, handler, "/v1/products", url.Values{}, headers)
		if status != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s, want 400", status, body)
		}

		retryStatus, retryBody := postFormStatusWithHeaders(t, handler, "/v1/products", url.Values{
			"name": {"Recovered Product"},
		}, headers)
		if retryStatus != http.StatusOK {
			t.Fatalf("status=%d body=%s, want 200 because validation errors are not cached", retryStatus, retryBody)
		}
	})
}

func TestStripeProtocolCompatibilityBaseline(t *testing.T) {
	handler := newTestHandler(t)

	status, body, headers := postFormStatusWithResponseHeaders(t, handler, "/v1/customers", url.Values{
		"email":    {"protocol@example.test"},
		"expand[]": {"subscriptions"},
	}, map[string]string{"Request-Id": "req_protocol_expand"})
	if status != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", status, body)
	}
	if got := headers.Get("Request-Id"); got != "req_protocol_expand" {
		t.Fatalf("Request-Id response header = %q, want caller request id", got)
	}

	status, body = postFormStatus(t, handler, "/v1/products", url.Values{
		"name":      {"Protocol Product"},
		"expand[0]": {"default_price"},
	})
	if status != http.StatusOK {
		t.Fatalf("status=%d body=%s, want expand[0] accepted and ignored", status, body)
	}

	list := getJSON[struct {
		Object  string `json:"object"`
		URL     string `json:"url"`
		HasMore bool   `json:"has_more"`
		Data    []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/customers?limit=1&expand[]=data.subscriptions")
	if list.Object != "list" || list.URL != "/v1/customers" || list.HasMore || len(list.Data) != 1 {
		t.Fatalf("list envelope = %#v, want Stripe-style list envelope", list)
	}

	traces := getJSON[struct {
		Object string                     `json:"object"`
		Data   []diagnostics.RequestTrace `json:"data"`
	}](t, handler, "/api/request-traces?requestId=req_protocol_expand")
	if traces.Object != "list" || len(traces.Data) != 1 {
		t.Fatalf("request traces = %#v, want trace filtered by request id", traces)
	}
	trace := traces.Data[0]
	if trace.RequestID != "req_protocol_expand" || trace.ResponseObject != "customer" || !strings.Contains(trace.RequestBody, "expand%5B%5D=subscriptions") {
		t.Fatalf("trace = %#v, want request id, response object, and expand request evidence", trace)
	}
}

func TestDashboardObjectsAndDebugBundle(t *testing.T) {
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()

	handler := newTestHandler(t)
	_ = postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL},
		"enabled_events": {"*"},
	})
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"buyer@example.test"},
		"name":  {"Buyer"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	_ = postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})

	objects := getJSON[struct {
		Object           string                    `json:"object"`
		Customers        []billing.Customer        `json:"customers"`
		CheckoutSessions []billing.CheckoutSession `json:"checkout_sessions"`
		Subscriptions    []billing.Subscription    `json:"subscriptions"`
		Invoices         []billing.Invoice         `json:"invoices"`
		PaymentIntents   []billing.PaymentIntent   `json:"payment_intents"`
		WebhookEvents    []webhooks.Event          `json:"webhook_events"`
	}](t, handler, "/api/objects")
	if objects.Object != "dashboard_objects" {
		t.Fatalf("object type = %q, want dashboard_objects", objects.Object)
	}
	if len(objects.Customers) != 1 || len(objects.CheckoutSessions) != 1 || len(objects.Subscriptions) != 1 || len(objects.Invoices) != 1 || len(objects.PaymentIntents) != 1 {
		t.Fatalf("dashboard objects = %#v, want completed checkout graph", objects)
	}
	if len(objects.WebhookEvents) == 0 {
		t.Fatal("dashboard objects did not include webhook events")
	}

	bundle := postJSON[struct {
		ID               string                     `json:"id"`
		Object           string                     `json:"object"`
		Target           map[string]string          `json:"target"`
		Filters          map[string]string          `json:"filters"`
		Timeline         []billing.TimelineEntry    `json:"timeline"`
		RequestTraces    []diagnostics.RequestTrace `json:"request_traces"`
		WebhookEvents    []webhooks.Event           `json:"webhook_events"`
		DeliveryAttempts []map[string]any           `json:"delivery_attempts"`
	}](t, handler, "/api/debug-bundles", map[string]string{
		"object_type": "checkoutSessions",
		"object_id":   session.ID,
	})
	if bundle.ID == "" || bundle.Object != "debug_bundle" {
		t.Fatalf("debug bundle = %#v, want identified debug_bundle", bundle)
	}
	if bundle.Target["id"] != session.ID || bundle.Filters["checkout_session_id"] != session.ID {
		t.Fatalf("bundle target=%#v filters=%#v, want checkout session filter", bundle.Target, bundle.Filters)
	}
	if len(bundle.Timeline) == 0 {
		t.Fatal("debug bundle timeline is empty")
	}
	if len(bundle.RequestTraces) == 0 {
		t.Fatal("debug bundle request_traces is empty")
	}
	if len(bundle.WebhookEvents) == 0 || len(bundle.DeliveryAttempts) == 0 {
		t.Fatalf("bundle webhook_events=%d delivery_attempts=%d, want webhook evidence", len(bundle.WebhookEvents), len(bundle.DeliveryAttempts))
	}
}

func TestRequestTraceRecordsStripeCalls(t *testing.T) {
	handler := newTestHandler(t)

	status, _ := postFormStatusWithHeaders(t, handler, "/v1/customers", url.Values{
		"email": {"trace@example.test"},
	}, map[string]string{
		"Authorization":   "Bearer sk_test_secret",
		"Idempotency-Key": "trace-key",
	})
	if status != http.StatusOK {
		t.Fatalf("create customer status = %d, want 200", status)
	}
	errorStatus, _ := postFormStatus(t, handler, "/v1/products", url.Values{})
	if errorStatus != http.StatusBadRequest {
		t.Fatalf("create product status = %d, want 400", errorStatus)
	}
	_ = getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/customers?email=trace@example.test&api_key=secret")

	traces := getJSON[struct {
		Object string                     `json:"object"`
		Data   []diagnostics.RequestTrace `json:"data"`
	}](t, handler, "/api/request-traces?limit=10")
	if traces.Object != "list" || len(traces.Data) != 3 {
		t.Fatalf("request traces = %#v, want three traces", traces)
	}

	var customerTrace, customerQueryTrace, errorTrace *diagnostics.RequestTrace
	for i := range traces.Data {
		trace := &traces.Data[i]
		switch trace.Path {
		case "/v1/customers":
			if trace.Method == http.MethodPost {
				customerTrace = trace
			} else {
				customerQueryTrace = trace
			}
		case "/v1/products":
			errorTrace = trace
		}
	}
	if customerTrace == nil || customerTrace.Status != http.StatusOK || customerTrace.ResponseObject != "customer" || customerTrace.ResponseObjectID == "" {
		t.Fatalf("customer trace = %#v, want successful customer response evidence", customerTrace)
	}
	if customerTrace.IdempotencyKey != "trace-key" || customerTrace.RequestHeaders["Authorization"] != security.MaskedValue {
		t.Fatalf("customer trace headers = %#v idempotency=%q, want masked authorization and idempotency key", customerTrace.RequestHeaders, customerTrace.IdempotencyKey)
	}
	if customerQueryTrace == nil || customerQueryTrace.Query != "api_key=%2A%2A%2A%2A&email=trace%40example.test" {
		t.Fatalf("customer query trace = %#v, want masked api_key", customerQueryTrace)
	}
	if errorTrace == nil || errorTrace.Status != http.StatusBadRequest || errorTrace.ErrorCode != "parameter_missing" || errorTrace.ErrorParam != "name" {
		t.Fatalf("error trace = %#v, want structured Stripe error evidence", errorTrace)
	}

	diagnostic := getJSON[struct {
		Object        string                     `json:"object"`
		Summary       map[string]int             `json:"summary"`
		RequestTraces []diagnostics.RequestTrace `json:"request_traces"`
	}](t, handler, "/api/diagnostics?limit=10")
	if diagnostic.Object != "diagnostic_bundle" || diagnostic.Summary["request_traces"] != 3 || len(diagnostic.RequestTraces) != 3 {
		t.Fatalf("diagnostic bundle = %#v, want request trace summary and data", diagnostic)
	}
}

func TestScenarioRunAPI(t *testing.T) {
	assertions := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/assertions/workspace/subscription" {
			t.Fatalf("assertion path = %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode assertion payload: %v", err)
		}
		if payload["target"] != "workspace.subscription" || payload["context"] == nil || payload["clock"] == nil {
			t.Fatalf("assertion payload = %#v", payload)
		}
		_, _ = w.Write([]byte(`{"pass":true}`))
	}))
	defer assertions.Close()

	handler := newTestHandler(t)
	report := postJSON[scenarios.Report](t, handler, "/api/scenarios/run", map[string]any{
		"name": "api-scenario",
		"app": map[string]any{
			"assertions": map[string]any{"baseUrl": assertions.URL + "/assertions"},
		},
		"catalog": map[string]any{
			"products": []map[string]any{{"id": "prod_pro", "name": "Pro"}},
			"prices": []map[string]any{{
				"id":         "price_pro_monthly",
				"product":    "prod_pro",
				"currency":   "usd",
				"unitAmount": 4900,
				"interval":   "month",
			}},
		},
		"steps": []map[string]any{
			{"id": "create-customer", "action": "customer.create", "params": map[string]any{"email": "api-scenario@example.test"}},
			{"id": "checkout", "action": "checkout.create", "params": map[string]any{"customerRef": "create-customer.customer.id", "price": "price_pro_monthly"}},
			{"id": "complete-checkout", "action": "checkout.complete", "params": map[string]any{"sessionRef": "checkout.session.id", "outcome": "payment_succeeded"}},
			{"id": "advance-clock", "action": "clock.advance", "params": map[string]any{"duration": "3d"}},
			{"id": "assert-active", "action": "app.assert", "params": map[string]any{"target": "workspace.subscription", "expected": map[string]any{"status": "active"}}},
		},
	})
	if report.Status != "passed" || report.ExitCode() != scenarios.ExitPass {
		t.Fatalf("scenario report = %#v, want passed", report)
	}
	if got := len(report.Steps); got != 5 {
		t.Fatalf("steps = %d, want 5", got)
	}
	if !report.ClockEnd.After(report.ClockStart) {
		t.Fatalf("clock did not advance: %s -> %s", report.ClockStart, report.ClockEnd)
	}
}

func TestPortalCoverageAPI(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"portal@example.test"},
		"name":  {"Portal User"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"9900"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})
	var sub billing.Subscription
	if err := json.Unmarshal(completion["subscription"], &sub); err != nil {
		t.Fatalf("decode subscription: %v", err)
	}

	state := getJSON[billing.PortalState](t, handler, "/api/portal?customer_id="+customer.ID)
	if state.Customer.ID != customer.ID || state.Subscription == nil || len(state.Invoices) != 1 {
		t.Fatalf("portal state = %#v, want customer subscription and invoice history", state)
	}

	plan := postJSON[struct {
		Subscription billing.Subscription `json:"subscription"`
		State        billing.PortalState  `json:"state"`
	}](t, handler, "/api/portal/subscriptions/"+sub.ID+"/plan-change", map[string]any{
		"plan":     "scale",
		"quantity": 9,
	})
	if plan.Subscription.Metadata["plan"] != "scale" || plan.Subscription.Items[0].Quantity != 9 {
		t.Fatalf("plan change subscription = %#v", plan.Subscription)
	}

	seats := postJSON[struct {
		Subscription billing.Subscription `json:"subscription"`
	}](t, handler, "/api/portal/subscriptions/"+sub.ID+"/seat-change", map[string]any{"quantity": 12})
	if seats.Subscription.Items[0].Quantity != 12 {
		t.Fatalf("seat quantity = %d, want 12", seats.Subscription.Items[0].Quantity)
	}

	canceled := postJSON[struct {
		Subscription billing.Subscription `json:"subscription"`
	}](t, handler, "/api/portal/subscriptions/"+sub.ID+"/cancel", map[string]string{"mode": "period"})
	if !canceled.Subscription.CancelAtPeriodEnd {
		t.Fatalf("cancelAtPeriodEnd = false, want scheduled cancellation")
	}

	resumed := postJSON[struct {
		Subscription billing.Subscription `json:"subscription"`
	}](t, handler, "/api/portal/subscriptions/"+sub.ID+"/resume", map[string]string{})
	if resumed.Subscription.CancelAtPeriodEnd || resumed.Subscription.Status != "active" {
		t.Fatalf("resumed subscription = %#v, want active without pending cancellation", resumed.Subscription)
	}

	immediate := postJSON[struct {
		Subscription billing.Subscription `json:"subscription"`
	}](t, handler, "/api/portal/subscriptions/"+sub.ID+"/cancel", map[string]string{"mode": "immediate"})
	if immediate.Subscription.Status != "canceled" || immediate.Subscription.CanceledAt == nil {
		t.Fatalf("immediate subscription = %#v, want canceled timestamp", immediate.Subscription)
	}

	payment := postJSON[struct {
		PaymentMethod billing.PaymentMethodSimulation `json:"payment_method"`
	}](t, handler, "/api/portal/customers/"+customer.ID+"/payment-method", map[string]string{"outcome": "fails"})
	if payment.PaymentMethod.Status != "failed" || payment.PaymentMethod.FailureCode == "" {
		t.Fatalf("payment method simulation = %#v, want failed card evidence", payment.PaymentMethod)
	}

	timeline := getJSON[struct {
		Data []billing.TimelineEntry `json:"data"`
	}](t, handler, "/api/timeline?customerId="+customer.ID)
	if len(timeline.Data) < 8 {
		t.Fatalf("timeline entries = %d, want portal transition evidence", len(timeline.Data))
	}
}

func TestBillingPortalSessionActionsEmitWebhooks(t *testing.T) {
	handler := newTestHandler(t)

	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_portal"},
		"frozen_time": {"1893456000"},
	})
	customer := postForm[struct {
		ID              string `json:"id"`
		InvoiceSettings struct {
			DefaultPaymentMethod *string `json:"default_payment_method"`
		} `json:"invoice_settings"`
	}](t, handler, "/v1/customers", url.Values{
		"id":         {"cus_portal_flow"},
		"email":      {"portal-flow@example.test"},
		"test_clock": {"clock_portal"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Portal Flow"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"9900"},
		"recurring[interval]": {"month"},
	})
	subscription := postForm[struct {
		ID               string `json:"id"`
		CurrentPeriodEnd int64  `json:"current_period_end"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {price.ID},
	})

	var receivedTypes []string
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Type string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode webhook payload: %v", err)
		}
		receivedTypes = append(receivedTypes, payload.Type)
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()
	postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL},
		"enabled_events": {"payment_method.attached,customer.updated,customer.subscription.updated,customer.subscription.deleted"},
	})

	portal := postForm[struct {
		Object        string         `json:"object"`
		Customer      string         `json:"customer"`
		Configuration *string        `json:"configuration"`
		ReturnURL     string         `json:"return_url"`
		URL           string         `json:"url"`
		Flow          map[string]any `json:"flow"`
	}](t, handler, "/v1/billing_portal/sessions", url.Values{
		"customer":                          {customer.ID},
		"return_url":                        {"https://app.example.test/billing/return"},
		"configuration":                     {"bpc_test"},
		"flow_data[type]":                   {"payment_method_update"},
		"flow_data[after_completion][type]": {"redirect"},
		"flow_data[after_completion][redirect][return_url]": {"https://app.example.test/billing/return"},
	})
	if portal.Object != "billing_portal.session" || portal.Customer != customer.ID || portal.Configuration == nil || *portal.Configuration != "bpc_test" || portal.Flow["type"] != "payment_method_update" {
		t.Fatalf("portal session = %#v, want Stripe-like session with flow/configuration", portal)
	}
	portalURL, err := url.Parse(portal.URL)
	if err != nil {
		t.Fatalf("parse portal URL: %v", err)
	}
	if portalURL.Query().Get("customer_id") != customer.ID || portalURL.Query().Get("return_url") != portal.ReturnURL || portalURL.Query().Get("redirect_on_action") != "true" || portalURL.Query().Get("flow") != "payment_method_update" {
		t.Fatalf("portal URL query = %s, want customer, return URL, action redirect, and flow", portalURL.RawQuery)
	}

	payment := postJSON[struct {
		PaymentMethod billing.PaymentMethodSimulation `json:"payment_method"`
	}](t, handler, "/api/portal/customers/"+customer.ID+"/payment-method", map[string]string{
		"outcome":        "succeeds",
		"payment_method": "pm_card_visa",
	})
	if payment.PaymentMethod.Status != "succeeded" || payment.PaymentMethod.PaymentMethodID != "pm_card_visa" {
		t.Fatalf("payment method update = %#v, want saved pm_card_visa", payment.PaymentMethod)
	}
	updatedCustomer := getJSON[struct {
		InvoiceSettings struct {
			DefaultPaymentMethod string `json:"default_payment_method"`
		} `json:"invoice_settings"`
	}](t, handler, "/v1/customers/"+customer.ID)
	if updatedCustomer.InvoiceSettings.DefaultPaymentMethod != "pm_card_visa" {
		t.Fatalf("default payment method = %q, want pm_card_visa", updatedCustomer.InvoiceSettings.DefaultPaymentMethod)
	}
	methods := getJSON[struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/payment_methods?customer="+customer.ID+"&type=card")
	if len(methods.Data) != 1 || methods.Data[0].ID != "pm_card_visa" {
		t.Fatalf("payment methods = %#v, want saved default method", methods.Data)
	}

	canceled := postJSON[struct {
		Subscription billing.Subscription `json:"subscription"`
	}](t, handler, "/api/portal/subscriptions/"+subscription.ID+"/cancel", map[string]string{"mode": "period"})
	if !canceled.Subscription.CancelAtPeriodEnd {
		t.Fatalf("portal cancel = %#v, want period-end cancellation", canceled.Subscription)
	}
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks/clock_portal/advance", url.Values{
		"frozen_time": {strconv.FormatInt(subscription.CurrentPeriodEnd+1, 10)},
	})

	wantTypes := []string{"payment_method.attached", "customer.updated", "customer.subscription.updated", "customer.subscription.deleted"}
	if strings.Join(receivedTypes, ",") != strings.Join(wantTypes, ",") {
		t.Fatalf("received webhook types = %#v, want %#v", receivedTypes, wantTypes)
	}

	filteredEvents := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?"+url.Values{
		"type": {"customer.updated"},
		"data.object.metadata[default_payment_method]": {"pm_card_visa"},
	}.Encode())
	if len(filteredEvents.Data) != 1 || filteredEvents.Data[0].Type != "customer.updated" {
		t.Fatalf("filtered customer.updated events = %#v, want one metadata-filtered event", filteredEvents.Data)
	}
}

func TestBillingPortalSessionStripeLikeShapeAndValidation(t *testing.T) {
	handler := newTestHandlerWithOptions(t, Options{PublicBaseURL: "http://127.0.0.1:18080"})
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"portal-shape@example.test"},
	})

	session := postForm[map[string]any](t, handler, "/v1/billing_portal/sessions", url.Values{
		"customer":     {customer.ID},
		"return_url":   {"https://app.example.test/account"},
		"locale":       {"auto"},
		"on_behalf_of": {"acct_connected"},
	})
	if session["object"] != "billing_portal.session" || session["customer"] != customer.ID || session["livemode"] != false {
		t.Fatalf("portal session = %#v, want Stripe-like object/customer/livemode fields", session)
	}
	if _, ok := session["flow"]; !ok || session["flow"] != nil {
		t.Fatalf("portal session flow = %#v, want explicit null flow field", session["flow"])
	}
	if session["locale"] != "auto" || session["on_behalf_of"] != "acct_connected" || session["return_url"] != "https://app.example.test/account" {
		t.Fatalf("portal session locale/on_behalf_of/return_url = %#v", session)
	}
	if urlValue, _ := session["url"].(string); !strings.HasPrefix(urlValue, "http://127.0.0.1:18080/portal?") {
		t.Fatalf("portal URL = %q, want hosted Billtap portal URL", urlValue)
	}

	tests := []struct {
		name       string
		values     url.Values
		wantParam  string
		wantCode   string
		wantStatus int
	}{
		{
			name:       "missing customer",
			values:     url.Values{"return_url": {"https://app.example.test/account"}},
			wantParam:  "customer",
			wantCode:   "parameter_missing",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid flow type",
			values:     url.Values{"customer": {customer.ID}, "flow_data[type]": {"not_a_flow"}},
			wantParam:  "flow_data[type]",
			wantCode:   "parameter_invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing subscription for cancel flow",
			values:     url.Values{"customer": {customer.ID}, "flow_data[type]": {"subscription_cancel"}},
			wantParam:  "flow_data[subscription_cancel][subscription]",
			wantCode:   "parameter_missing",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid locale",
			values:     url.Values{"customer": {customer.ID}, "locale": {"pirate"}},
			wantParam:  "locale",
			wantCode:   "parameter_invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown parameter",
			values:     url.Values{"customer": {customer.ID}, "made_up": {"true"}},
			wantParam:  "made_up",
			wantCode:   "parameter_unknown",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := postFormStatus(t, handler, "/v1/billing_portal/sessions", tt.values)
			errBody := decodeErrorBody(t, body)
			if status != tt.wantStatus || errBody.Error.Param != tt.wantParam || errBody.Error.Code != tt.wantCode || errBody.Error.Type != "invalid_request_error" {
				t.Fatalf("status=%d error=%#v, want status=%d param=%s code=%s", status, errBody.Error, tt.wantStatus, tt.wantParam, tt.wantCode)
			}
		})
	}
}

func TestFixtureCustomerPaymentMethodsCanBeEmptyUntilPortalSave(t *testing.T) {
	handler := newTestHandler(t)
	pack := map[string]any{
		"name":      "sample-payment-method-fixtures",
		"runId":     "run-payment-method-fixtures-1",
		"namespace": "sample-app",
		"customers": []map[string]any{
			{
				"id":              "cus_fixture_no_payment_methods",
				"email":           "no-payment-methods@example.test",
				"payment_methods": []map[string]any{},
				"metadata":        map[string]string{"tenantId": "sample"},
			},
			{
				"id":                      "cus_fixture_no_payment_methods_mode",
				"email":                   "no-payment-methods-mode@example.test",
				"payment_methods_fixture": "empty",
				"metadata":                map[string]string{"tenantId": "sample"},
			},
			{
				"id":    "cus_fixture_default_payment_method_projection",
				"email": "default-payment-method-projection@example.test",
			},
			{
				"id":    "cus_fixture_explicit_payment_methods",
				"email": "explicit-payment-methods@example.test",
				"payment_methods": []map[string]any{
					{"id": "pm_fixture_secondary"},
					{"id": "pm_fixture_primary", "default": true},
				},
			},
		},
	}
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["customers"] != 4 {
		t.Fatalf("apply summary = %#v, want four customers", applied.Summary)
	}

	assertPaymentMethodIDs := func(customerID string, want ...string) {
		t.Helper()
		methods := getJSON[struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}](t, handler, "/v1/payment_methods?customer="+customerID+"&type=card")
		var got []string
		for _, method := range methods.Data {
			got = append(got, method.ID)
		}
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Fatalf("payment methods for %s = %#v, want %#v", customerID, got, want)
		}

		nested := getJSON[struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}](t, handler, "/v1/customers/"+customerID+"/payment_methods?type=card")
		got = got[:0]
		for _, method := range nested.Data {
			got = append(got, method.ID)
		}
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Fatalf("nested payment methods for %s = %#v, want %#v", customerID, got, want)
		}
	}

	assertPaymentMethodIDs("cus_fixture_no_payment_methods")
	assertPaymentMethodIDs("cus_fixture_no_payment_methods_mode")
	assertPaymentMethodIDs("cus_fixture_default_payment_method_projection", "pm_cus_fixture_default_payment_method_projection")
	assertPaymentMethodIDs("cus_fixture_explicit_payment_methods", "pm_fixture_secondary", "pm_fixture_primary")

	noMethodCustomer := getJSON[struct {
		InvoiceSettings struct {
			DefaultPaymentMethod *string `json:"default_payment_method"`
		} `json:"invoice_settings"`
	}](t, handler, "/v1/customers/cus_fixture_no_payment_methods")
	if noMethodCustomer.InvoiceSettings.DefaultPaymentMethod != nil {
		t.Fatalf("default payment method = %v, want nil", *noMethodCustomer.InvoiceSettings.DefaultPaymentMethod)
	}
	explicitCustomer := getJSON[struct {
		InvoiceSettings struct {
			DefaultPaymentMethod string `json:"default_payment_method"`
		} `json:"invoice_settings"`
	}](t, handler, "/v1/customers/cus_fixture_explicit_payment_methods")
	if explicitCustomer.InvoiceSettings.DefaultPaymentMethod != "pm_fixture_primary" {
		t.Fatalf("explicit default payment method = %q, want pm_fixture_primary", explicitCustomer.InvoiceSettings.DefaultPaymentMethod)
	}

	_ = postJSON[struct {
		PaymentMethod billing.PaymentMethodSimulation `json:"payment_method"`
	}](t, handler, "/api/portal/customers/cus_fixture_no_payment_methods/payment-method", map[string]string{
		"outcome":        "succeeds",
		"payment_method": "pm_card_visa",
	})
	assertPaymentMethodIDs("cus_fixture_no_payment_methods", "pm_cus_fixture_no_payment_methods", "pm_card_visa")

	reapplied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if reapplied.Summary["customers"] != 4 {
		t.Fatalf("reapply summary = %#v, want four customers", reapplied.Summary)
	}
	assertPaymentMethodIDs("cus_fixture_no_payment_methods")
}

func TestFixtureCustomerDefaultPaymentIntentOutcome(t *testing.T) {
	handler := newTestHandler(t)
	pack := map[string]any{
		"name":      "sample-payment-intent-outcome-fixtures",
		"runId":     "run-payment-intent-outcomes-1",
		"namespace": "sample-app",
		"customers": []map[string]any{
			{
				"id":                             "cus_fixture_default_pi_decline",
				"email":                          "default-pi-decline@example.test",
				"default_payment_intent_outcome": "card_declined",
				"metadata":                       map[string]string{"tenantId": "sample"},
			},
		},
	}
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["customers"] != 1 || applied.Customers[0].Metadata[billing.MetadataDefaultPaymentIntentOutcome] != "card_declined" {
		t.Fatalf("apply result = %#v, want customer default payment intent outcome metadata", applied)
	}

	intent := postForm[struct {
		Status           string            `json:"status"`
		Metadata         map[string]string `json:"metadata"`
		LastPaymentError struct {
			Code string `json:"code"`
		} `json:"last_payment_error"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"amount":   {"4200"},
		"currency": {"usd"},
		"customer": {"cus_fixture_default_pi_decline"},
		"confirm":  {"true"},
	})
	if intent.Status != "requires_payment_method" || intent.Metadata[billing.MetadataPaymentIntentOutcome] != "card_declined" || intent.LastPaymentError.Code != "card_declined" {
		t.Fatalf("intent = %#v, want fixture customer default decline outcome", intent)
	}
}

func TestPaymentMethodsStripeLikeShapeFilteringAndValidation(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"payment-method-shape@example.test"},
	})

	methods := getJSON[struct {
		Object  string `json:"object"`
		URL     string `json:"url"`
		HasMore bool   `json:"has_more"`
		Data    []struct {
			ID             string         `json:"id"`
			Object         string         `json:"object"`
			AllowRedisplay string         `json:"allow_redisplay"`
			BillingDetails map[string]any `json:"billing_details"`
			Card           map[string]any `json:"card"`
			Customer       string         `json:"customer"`
			Livemode       bool           `json:"livemode"`
			Metadata       map[string]any `json:"metadata"`
			Redaction      any            `json:"redaction"`
			Type           string         `json:"type"`
		} `json:"data"`
	}](t, handler, "/v1/payment_methods?customer="+customer.ID+"&type=card&limit=3")
	if methods.Object != "list" || methods.URL != "/v1/payment_methods" || methods.HasMore || len(methods.Data) != 1 {
		t.Fatalf("payment methods list = %#v, want one Stripe-like list item", methods)
	}
	method := methods.Data[0]
	if method.Object != "payment_method" || method.Type != "card" || method.Customer != customer.ID || method.Livemode || method.AllowRedisplay != "unspecified" {
		t.Fatalf("payment method top-level fields = %#v", method)
	}
	if method.BillingDetails["address"] == nil || method.BillingDetails["email"] != nil || method.Metadata == nil || method.Redaction != nil {
		t.Fatalf("payment method billing/metadata/redaction = %#v/%#v/%#v", method.BillingDetails, method.Metadata, method.Redaction)
	}
	checks, _ := method.Card["checks"].(map[string]any)
	networks, _ := method.Card["networks"].(map[string]any)
	threeDSecure, _ := method.Card["three_d_secure_usage"].(map[string]any)
	if method.Card["brand"] != "visa" || method.Card["country"] != "US" || method.Card["funding"] != "credit" || method.Card["last4"] != "4242" || method.Card["fingerprint"] == "" {
		t.Fatalf("payment method card fields = %#v", method.Card)
	}
	if checks["cvc_check"] != "pass" || networks["preferred"] != nil || threeDSecure["supported"] != true {
		t.Fatalf("payment method nested card fields checks=%#v networks=%#v threeDS=%#v", checks, networks, threeDSecure)
	}

	nested := getJSON[struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/customers/"+customer.ID+"/payment_methods?type=card")
	if len(nested.Data) != 1 || nested.Data[0].ID != method.ID {
		t.Fatalf("nested payment methods = %#v, want same fixture method", nested.Data)
	}

	nonCard := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/payment_methods?customer="+customer.ID+"&type=us_bank_account")
	if len(nonCard.Data) != 0 {
		t.Fatalf("non-card payment methods = %#v, want empty list for unsupported local type", nonCard.Data)
	}

	tests := []struct {
		name      string
		path      string
		wantParam string
		wantCode  string
	}{
		{
			name:      "invalid type",
			path:      "/v1/payment_methods?customer=" + customer.ID + "&type=not_a_type",
			wantParam: "type",
			wantCode:  "parameter_invalid",
		},
		{
			name:      "invalid redisplay enum",
			path:      "/v1/customers/" + customer.ID + "/payment_methods?allow_redisplay=forever",
			wantParam: "allow_redisplay",
			wantCode:  "parameter_invalid",
		},
		{
			name:      "unknown query parameter",
			path:      "/v1/payment_methods?customer=" + customer.ID + "&made_up=true",
			wantParam: "made_up",
			wantCode:  "parameter_unknown",
		},
		{
			name:      "invalid limit",
			path:      "/v1/payment_methods?customer=" + customer.ID + "&limit=zero",
			wantParam: "limit",
			wantCode:  "parameter_invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := getStatus(t, handler, tt.path)
			errBody := decodeErrorBody(t, body)
			if status != http.StatusBadRequest || errBody.Error.Param != tt.wantParam || errBody.Error.Code != tt.wantCode || errBody.Error.Type != "invalid_request_error" {
				t.Fatalf("status=%d error=%#v, want param=%s code=%s", status, errBody.Error, tt.wantParam, tt.wantCode)
			}
		})
	}
}

func TestPriceSearchSupportsOneTimeLookupKey(t *testing.T) {
	handler := newTestHandler(t)
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{
		"id":              {"prod_extra_export"},
		"name":            {"Extra Export"},
		"metadata[scope]": {"exports"},
	})
	oneTime := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"id":                    {"price_extra_export_single_charge"},
		"product":               {product.ID},
		"currency":              {"usd"},
		"unit_amount":           {"1200"},
		"lookup_key":            {"extra_export_single_charge"},
		"metadata[scope]":       {"exports"},
		"metadata[environment]": {"e2e"},
	})
	recurring := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"id":                    {"price_extra_export_monthly"},
		"product":               {product.ID},
		"currency":              {"usd"},
		"unit_amount":           {"9900"},
		"lookup_key":            {"extra_export_single_charge"},
		"recurring[interval]":   {"month"},
		"metadata[scope]":       {"exports"},
		"metadata[environment]": {"e2e"},
	})
	if oneTime.ID == "" || recurring.ID == "" {
		t.Fatalf("seeded prices oneTime=%#v recurring=%#v", oneTime, recurring)
	}

	query := "active:'true' AND type:'one_time' AND lookup_key:'extra_export_single_charge'"
	search := getJSON[struct {
		Object   string  `json:"object"`
		URL      string  `json:"url"`
		HasMore  bool    `json:"has_more"`
		NextPage *string `json:"next_page"`
		Data     []struct {
			ID        string            `json:"id"`
			Object    string            `json:"object"`
			Active    bool              `json:"active"`
			LookupKey string            `json:"lookup_key"`
			Metadata  map[string]string `json:"metadata"`
			Type      string            `json:"type"`
		} `json:"data"`
	}](t, handler, "/v1/prices/search?"+url.Values{"query": {query}}.Encode())
	if search.Object != "search_result" || search.URL != "/v1/prices/search?query="+query || search.HasMore || search.NextPage != nil {
		t.Fatalf("price search envelope = %#v, want Stripe-like search result", search)
	}
	if len(search.Data) != 1 || search.Data[0].ID != oneTime.ID || search.Data[0].Object != "price" || !search.Data[0].Active || search.Data[0].Type != "one_time" || search.Data[0].LookupKey != "extra_export_single_charge" {
		t.Fatalf("price search data = %#v, want one matching one-time price", search.Data)
	}

	metadataSearch := getJSON[struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/prices/search?"+url.Values{
		"query": {"active:'true' AND type:'one_time' AND metadata['scope']:'exports' AND metadata['environment']:'e2e'"},
	}.Encode())
	if len(metadataSearch.Data) != 1 || metadataSearch.Data[0].ID != oneTime.ID {
		t.Fatalf("metadata price search = %#v, want one-time price", metadataSearch.Data)
	}

	empty := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/prices/search?"+url.Values{
		"query": {"active:'true' AND type:'one_time' AND lookup_key:'missing'"},
	}.Encode())
	if len(empty.Data) != 0 {
		t.Fatalf("empty price search = %#v, want no matches", empty.Data)
	}
}

func TestExpandedStripeSurfaceSimulations(t *testing.T) {
	handler := newTestHandler(t)

	var receivedTypes []string
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Type string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode webhook payload: %v", err)
		}
		receivedTypes = append(receivedTypes, payload.Type)
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()
	endpoint := postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL},
		"enabled_events": {"payment_intent.*,customer.subscription.updated,charge.dispute.*"},
	})

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"expanded@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Expanded"}})
	priceBasic := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1000"},
		"recurring[interval]": {"month"},
	})
	priceScale := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"2000"},
		"recurring[interval]": {"month"},
	})

	coupon := postForm[struct {
		ID         string `json:"id"`
		Object     string `json:"object"`
		PercentOff int64  `json:"percent_off"`
	}](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_launch"},
		"percent_off": {"25"},
		"duration":    {"once"},
	})
	if coupon.Object != "coupon" || coupon.ID != "coupon_launch" || coupon.PercentOff != 25 {
		t.Fatalf("coupon = %#v, want local coupon evidence", coupon)
	}
	promo := postForm[struct {
		Object string `json:"object"`
		Code   string `json:"code"`
		Coupon struct {
			ID string `json:"id"`
		} `json:"coupon"`
	}](t, handler, "/v1/promotion_codes", url.Values{
		"coupon": {coupon.ID},
		"code":   {"LAUNCH25"},
	})
	if promo.Object != "promotion_code" || promo.Code != "LAUNCH25" || promo.Coupon.ID != coupon.ID {
		t.Fatalf("promotion code = %#v, want coupon-backed promotion", promo)
	}

	sca := postForm[struct {
		ID            string         `json:"id"`
		Status        string         `json:"status"`
		PaymentMethod string         `json:"payment_method"`
		NextAction    map[string]any `json:"next_action"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"customer": {customer.ID},
		"amount":   {"3200"},
		"currency": {"usd"},
		"outcome":  {"requires_action"},
	})
	if sca.Status != "requires_action" || sca.PaymentMethod != "pm_card_threeDSecure2Required" || sca.NextAction["type"] != "use_stripe_sdk" {
		t.Fatalf("SCA payment intent = %#v, want requires_action with 3DS payment method and next_action", sca)
	}
	completedSCA := postJSON[struct {
		Status string `json:"status"`
	}](t, handler, "/api/payment_intents/"+sca.ID+"/complete_action", map[string]string{})
	if completedSCA.Status != "succeeded" {
		t.Fatalf("completed SCA = %#v, want succeeded", completedSCA)
	}

	bankTransfer := postForm[struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		PaymentMethod string `json:"payment_method"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"customer": {customer.ID},
		"amount":   {"4500"},
		"currency": {"usd"},
		"outcome":  {"bank_transfer"},
	})
	if bankTransfer.Status != "processing" || bankTransfer.PaymentMethod != "pm_bank_transfer" {
		t.Fatalf("bank transfer intent = %#v, want processing bank transfer", bankTransfer)
	}
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/customers/"+customer.ID+"/fund_cash_balance", url.Values{
		"amount": {"4500"},
	})
	settledBankTransfer := getJSON[struct {
		Status string `json:"status"`
	}](t, handler, "/v1/payment_intents/"+bankTransfer.ID)
	if settledBankTransfer.Status != "succeeded" {
		t.Fatalf("settled bank transfer = %#v, want succeeded after cash-balance funding", settledBankTransfer)
	}
	cashBalance := getJSON[struct {
		Available map[string]int64 `json:"available"`
	}](t, handler, "/v1/customers/"+customer.ID+"/cash_balance")
	if cashBalance.Available["usd"] != 4500 {
		t.Fatalf("cash balance = %#v, want funded usd balance", cashBalance.Available)
	}

	clock := postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_schedule"},
		"frozen_time": {"1893456000"},
	})
	_ = clock
	clockCustomer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":      {"schedule@example.test"},
		"test_clock": {"clock_schedule"},
	})
	subscription := postForm[struct {
		ID               string `json:"id"`
		CurrentPeriodEnd int64  `json:"current_period_end"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {clockCustomer.ID},
		"items[0][price]": {priceBasic.ID},
	})
	schedule := postForm[struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}](t, handler, "/v1/subscription_schedules", url.Values{
		"from_subscription":             {subscription.ID},
		"phases[0][start_date]":         {strconv.FormatInt(subscription.CurrentPeriodEnd, 10)},
		"phases[0][items][0][price]":    {priceScale.ID},
		"phases[0][items][0][quantity]": {"3"},
	})
	if schedule.ID == "" || schedule.Status != "active" {
		t.Fatalf("subscription schedule = %#v, want active schedule", schedule)
	}
	advance := postForm[struct {
		BilltapAdvanceResult struct {
			ScheduledCount int `json:"scheduled_count"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_schedule/advance", url.Values{
		"frozen_time": {strconv.FormatInt(subscription.CurrentPeriodEnd+1, 10)},
	})
	if advance.BilltapAdvanceResult.ScheduledCount != 1 {
		t.Fatalf("clock advance = %#v, want scheduled phase applied", advance)
	}
	updatedSubscription := getJSON[struct {
		Items struct {
			Data []struct {
				Price struct {
					ID string `json:"id"`
				} `json:"price"`
				Quantity int64 `json:"quantity"`
			} `json:"data"`
		} `json:"items"`
	}](t, handler, "/v1/subscriptions/"+subscription.ID)
	if len(updatedSubscription.Items.Data) != 1 || updatedSubscription.Items.Data[0].Price.ID != priceScale.ID || updatedSubscription.Items.Data[0].Quantity != 3 {
		t.Fatalf("updated subscription items = %#v, want scheduled price/quantity", updatedSubscription.Items.Data)
	}

	dispute := postForm[struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Charge string `json:"charge"`
	}](t, handler, "/v1/charges/ch_test_123/dispute", url.Values{
		"amount": {"1500"},
		"reason": {"fraudulent"},
	})
	if dispute.ID == "" || dispute.Status != "needs_response" || dispute.Charge != "ch_test_123" {
		t.Fatalf("dispute = %#v, want charge dispute evidence", dispute)
	}
	closedDispute := postForm[struct {
		Status string `json:"status"`
	}](t, handler, "/v1/disputes/"+dispute.ID+"/close", url.Values{})
	if closedDispute.Status != "won" {
		t.Fatalf("closed dispute = %#v, want won", closedDispute)
	}

	attempts := getJSON[struct {
		Data []struct {
			EventID string `json:"event_id"`
		} `json:"data"`
	}](t, handler, "/v1/webhook_endpoints/"+endpoint.ID+"/attempts")
	if len(attempts.Data) == 0 {
		t.Fatalf("endpoint attempts = %#v, want delivery attempt evidence", attempts.Data)
	}
	eventIDs := make([]string, 0, len(attempts.Data))
	for _, attempt := range attempts.Data {
		if attempt.EventID != "" {
			eventIDs = append(eventIDs, attempt.EventID)
		}
		if len(eventIDs) == 2 {
			break
		}
	}
	group := postForm[struct {
		Data []struct {
			Metadata map[string]string `json:"metadata"`
		} `json:"data"`
	}](t, handler, "/api/events/replay-group", url.Values{
		"event_ids":          {strings.Join(eventIDs, ",")},
		"mode":               {"out_of_order"},
		"signature_mismatch": {"true"},
	})
	if len(group.Data) != len(eventIDs) {
		t.Fatalf("replay group = %#v, want one replay attempt per event", group.Data)
	}
	for _, attempt := range group.Data {
		if attempt.Metadata["signature_mismatch"] != "true" || attempt.Metadata["out_of_order"] != "true" {
			t.Fatalf("replay group attempt metadata = %#v, want signature mismatch and out-of-order evidence", attempt.Metadata)
		}
	}

	expectedEvents := []string{"payment_intent.created", "payment_intent.requires_action", "payment_intent.succeeded", "payment_intent.created", "payment_intent.processing", "payment_intent.succeeded", "customer.subscription.updated", "charge.dispute.created", "charge.dispute.closed"}
	for _, eventType := range expectedEvents {
		if !containsString(receivedTypes, eventType) {
			t.Fatalf("received webhook types = %#v, missing %s", receivedTypes, eventType)
		}
	}
}

func TestStripeCompatCatalogAndPortalEndpoints(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"id":                 {"cus_e2e_pro"},
		"email":              {"stripe-compat@example.test"},
		"metadata[tenantId]": {"saas"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{
		"id":                                  {"prod_e2e_saas_premium"},
		"name":                                {"SaaS Pro Plan (E2E)"},
		"metadata[tenantId]":                  {"saas"},
		"metadata[tier]":                      {"PREMIUM"},
		"metadata[tierLevel]":                 {"3"},
		"metadata[basicSeat]":                 {"1"},
		"metadata[freeTrialPeriodDays]":       {"14"},
		"metadata[planExportLimit]":           {"-1"},
		"metadata[additionalSeatExportLimit]": {"-1"},
		"metadata[freeTrialExportLimit]":      {"100"},
		"metadata[productType]":               {"WORKSPACE_PLAN"},
		"metadata[version]":                   {"2"},
		"metadata[default_price]":             {"price_e2e_saas_premium_monthly"},
	})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"id":                        {"price_e2e_saas_premium_monthly"},
		"product":                   {product.ID},
		"currency":                  {"usd"},
		"unit_amount":               {"30000"},
		"lookup_key":                {"saas_plan_premium_monthly"},
		"recurring[interval]":       {"month"},
		"recurring[interval_count]": {"1"},
		"metadata[tenantId]":        {"saas"},
		"metadata[tier]":            {"PREMIUM"},
	})
	if customer.ID != "cus_e2e_pro" || product.ID != "prod_e2e_saas_premium" || price.ID != "price_e2e_saas_premium_monthly" {
		t.Fatalf("seeded IDs = customer:%s product:%s price:%s", customer.ID, product.ID, price.ID)
	}

	search := getJSON[struct {
		Object string `json:"object"`
		Data   []struct {
			ID       string            `json:"id"`
			Metadata map[string]string `json:"metadata"`
			Created  int64             `json:"created"`
		} `json:"data"`
	}](t, handler, "/v1/products/search?query=metadata['tenantId']:'saas'%20AND%20metadata['productType']:'WORKSPACE_PLAN'")
	if search.Object != "search_result" || len(search.Data) != 1 || search.Data[0].Metadata["tier"] != "PREMIUM" || search.Data[0].Created == 0 {
		t.Fatalf("search result = %#v", search)
	}

	prices := getJSON[struct {
		Data []struct {
			ID        string `json:"id"`
			LookupKey string `json:"lookup_key"`
			Recurring struct {
				Interval string `json:"interval"`
			} `json:"recurring"`
		} `json:"data"`
	}](t, handler, "/v1/prices?product="+product.ID+"&active=true&type=recurring")
	if len(prices.Data) != 1 || prices.Data[0].LookupKey != "saas_plan_premium_monthly" || prices.Data[0].Recurring.Interval != "month" {
		t.Fatalf("prices = %#v", prices)
	}

	portal := postForm[struct {
		Object string `json:"object"`
		URL    string `json:"url"`
	}](t, handler, "/v1/billing_portal/sessions", url.Values{
		"customer":   {customer.ID},
		"return_url": {"https://app.example.test/settings"},
	})
	if portal.Object != "billing_portal.session" || !strings.Contains(portal.URL, "/portal?customer_id="+customer.ID) {
		t.Fatalf("portal = %#v", portal)
	}

	methods := getJSON[struct {
		Data []struct {
			ID   string `json:"id"`
			Card struct {
				Last4 string `json:"last4"`
			} `json:"card"`
		} `json:"data"`
	}](t, handler, "/v1/payment_methods?customer="+customer.ID+"&type=card")
	if len(methods.Data) != 1 || methods.Data[0].Card.Last4 != "4242" {
		t.Fatalf("payment methods = %#v", methods)
	}

	subscription := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {price.ID},
	})

	canceled := postForm[struct {
		ID                  string `json:"id"`
		CancelAtPeriodEnd   bool   `json:"cancel_at_period_end"`
		CancelAt            *int64 `json:"cancel_at"`
		CanceledAt          *int64 `json:"canceled_at"`
		CancellationDetails struct {
			Comment  *string `json:"comment"`
			Feedback *string `json:"feedback"`
		} `json:"cancellation_details"`
	}](t, handler, "/v1/subscriptions/"+subscription.ID, url.Values{
		"cancel_at_period_end":           {"true"},
		"cancellation_details[comment]":  {"too expensive"},
		"cancellation_details[feedback]": {"too_expensive"},
	})
	if !canceled.CancelAtPeriodEnd || canceled.CancelAt == nil || canceled.CanceledAt == nil {
		t.Fatalf("canceled subscription = %#v, want pending cancellation timestamps", canceled)
	}
	if canceled.CancellationDetails.Comment == nil || *canceled.CancellationDetails.Comment != "too expensive" {
		t.Fatalf("cancellation comment = %#v, want preserved", canceled.CancellationDetails.Comment)
	}
	if canceled.CancellationDetails.Feedback == nil || *canceled.CancellationDetails.Feedback != "too_expensive" {
		t.Fatalf("cancellation feedback = %#v, want preserved", canceled.CancellationDetails.Feedback)
	}

	resumed := postForm[struct {
		CancelAtPeriodEnd   bool   `json:"cancel_at_period_end"`
		CancelAt            *int64 `json:"cancel_at"`
		CanceledAt          *int64 `json:"canceled_at"`
		CancellationDetails struct {
			Comment  *string `json:"comment"`
			Feedback *string `json:"feedback"`
		} `json:"cancellation_details"`
	}](t, handler, "/v1/subscriptions/"+subscription.ID, url.Values{
		"cancel_at_period_end": {"false"},
	})
	if resumed.CancelAtPeriodEnd || resumed.CancelAt != nil || resumed.CanceledAt != nil {
		t.Fatalf("resumed subscription = %#v, want cancellation cleared", resumed)
	}
	if resumed.CancellationDetails.Comment != nil || resumed.CancellationDetails.Feedback != nil {
		t.Fatalf("resumed cancellation details = %#v, want cleared", resumed.CancellationDetails)
	}
}

func TestSubscriptionUpdatePreservesItemsAndSupportsAdditiveSeatItems(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"paid-update@example.test"},
	})
	plan := postForm[billing.Product](t, handler, "/v1/products", url.Values{
		"id":                            {"prod_paid_update_plan"},
		"name":                          {"Paid Update Plan"},
		"metadata[productType]":         {"WORKSPACE_PLAN"},
		"metadata[tenantId]":            {"sample"},
		"metadata[tier]":                {"STANDARD"},
		"metadata[tierLevel]":           {"2"},
		"metadata[basicSeat]":           {"1"},
		"metadata[freeTrialPeriodDays]": {"14"},
	})
	seat := postForm[billing.Product](t, handler, "/v1/products", url.Values{
		"id":                    {"prod_paid_update_seat"},
		"name":                  {"Paid Update Seat"},
		"metadata[productType]": {"ADDITIONAL_SEAT"},
		"metadata[tenantId]":    {"sample"},
	})
	standardMonthly := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {plan.ID},
		"currency":            {"usd"},
		"unit_amount":         {"15000"},
		"lookup_key":          {"sample_plan_standard_monthly"},
		"recurring[interval]": {"month"},
	})
	proYearly := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {plan.ID},
		"currency":            {"usd"},
		"unit_amount":         {"288000"},
		"lookup_key":          {"sample_plan_premium_yearly"},
		"recurring[interval]": {"year"},
	})
	seatYearly := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {seat.ID},
		"currency":            {"usd"},
		"unit_amount":         {"96000"},
		"lookup_key":          {"sample_seat_yearly"},
		"recurring[interval]": {"year"},
	})

	type subscriptionItem struct {
		ID       string `json:"id"`
		Quantity int64  `json:"quantity"`
		Price    struct {
			ID string `json:"id"`
		} `json:"price"`
	}
	type subscriptionResponse struct {
		ID    string `json:"id"`
		Items struct {
			Data []subscriptionItem `json:"data"`
		} `json:"items"`
	}

	created := postForm[subscriptionResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {standardMonthly.ID},
	})
	if len(created.Items.Data) != 1 {
		t.Fatalf("created items = %#v, want one plan item", created.Items.Data)
	}

	upgraded := postForm[subscriptionResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {created.Items.Data[0].ID},
		"items[0][price]":      {proYearly.ID},
		"items[0][quantity]":   {"1"},
		"proration_behavior":   {"always_invoice"},
		"payment_behavior":     {"error_if_incomplete"},
		"billing_cycle_anchor": {"now"},
	})
	if len(upgraded.Items.Data) != 1 || upgraded.Items.Data[0].Price.ID != proYearly.ID {
		t.Fatalf("upgraded items = %#v, want existing item price replaced", upgraded.Items.Data)
	}

	withSeats := postForm[subscriptionResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][price]":      {seatYearly.ID},
		"items[0][quantity]":   {"2"},
		"proration_behavior":   {"always_invoice"},
		"payment_behavior":     {"error_if_incomplete"},
		"billing_cycle_anchor": {"now"},
	})
	if len(withSeats.Items.Data) != 2 {
		t.Fatalf("seat update items = %#v, want plan item plus additive seat item", withSeats.Items.Data)
	}
	if withSeats.Items.Data[0].Price.ID != proYearly.ID {
		t.Fatalf("plan item = %#v, want existing plan item preserved", withSeats.Items.Data[0])
	}
	if withSeats.Items.Data[1].Price.ID != seatYearly.ID || withSeats.Items.Data[1].Quantity != 2 {
		t.Fatalf("seat item = %#v, want additive seat item quantity 2", withSeats.Items.Data[1])
	}
}

func TestFixtureApplySnapshotAndAssertAPI(t *testing.T) {
	handler := newTestHandler(t)
	pack := map[string]any{
		"name":      "saas-basic",
		"runId":     "run-fixture-1",
		"namespace": "sample-app",
		"customers": []map[string]any{{
			"id":       "cus_fixture_pro",
			"email":    "fixture-pro@example.test",
			"metadata": map[string]string{"tenantId": "saas"},
		}},
		"catalog": map[string]any{
			"products": []map[string]any{{
				"id":   "prod_fixture_saas_premium",
				"name": "SaaS Pro Fixture",
				"metadata": map[string]string{
					"tenantId":      "saas",
					"productType":   "WORKSPACE_PLAN",
					"tier":          "PREMIUM",
					"default_price": "price_fixture_saas_premium_monthly",
				},
			}},
			"prices": []map[string]any{{
				"id":         "price_fixture_saas_premium_monthly",
				"product":    "prod_fixture_saas_premium",
				"currency":   "usd",
				"unitAmount": 30000,
				"lookupKey":  "saas_plan_premium_monthly",
				"interval":   "month",
				"metadata":   map[string]string{"tenantId": "saas", "tier": "PREMIUM"},
			}},
		},
		"subscriptions": []map[string]any{{
			"ref":      "pro-workspace",
			"customer": "cus_fixture_pro",
			"price":    "price_fixture_saas_premium_monthly",
			"quantity": 3,
			"outcome":  "payment_succeeded",
			"metadata": map[string]string{"tenantId": "saas"},
		}},
		"assertions": []map[string]any{
			{"target": "customer", "id": "cus_fixture_pro"},
			{"target": "price", "lookupKey": "saas_plan_premium_monthly"},
			{"target": "subscription", "customer": "cus_fixture_pro", "price": "price_fixture_saas_premium_monthly", "status": "active", "quantity": 3},
		},
	}
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Assertions == nil || !applied.Assertions.Pass {
		t.Fatalf("apply assertions = %#v, want pass", applied.Assertions)
	}
	if applied.Summary["customers"] != 1 || applied.Summary["products"] != 1 || applied.Summary["prices"] != 1 || applied.Summary["subscriptions"] != 1 {
		t.Fatalf("apply summary = %#v", applied.Summary)
	}

	appliedAgain := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if len(appliedAgain.Subscriptions) != 1 || appliedAgain.Subscriptions[0].ID != applied.Subscriptions[0].ID {
		t.Fatalf("second apply subscriptions = %#v, want idempotent fixture ref", appliedAgain.Subscriptions)
	}

	snapshot := getJSON[fixtures.Snapshot](t, handler, "/api/fixtures/snapshot?runId=run-fixture-1&tenantId=saas&namespace=sample-app")
	if snapshot.Summary["customers"] != 1 || snapshot.Summary["products"] != 1 || snapshot.Summary["prices"] != 1 || snapshot.Summary["subscriptions"] != 1 || snapshot.Summary["invoices"] != 1 || snapshot.Summary["payment_intents"] != 1 {
		t.Fatalf("snapshot summary = %#v", snapshot.Summary)
	}
	if snapshot.Subscriptions[0].Metadata[fixtures.MetadataFixtureRef] != "pro-workspace" {
		t.Fatalf("subscription metadata = %#v, want fixture ref", snapshot.Subscriptions[0].Metadata)
	}

	report := postJSON[fixtures.AssertionReport](t, handler, "/api/fixtures/assert", map[string]any{
		"name": "saas-basic-check",
		"filter": map[string]any{
			"runId":     "run-fixture-1",
			"tenantId":  "saas",
			"namespace": "sample-app",
		},
		"expect": []map[string]any{{
			"target":       "timeline",
			"customer":     "cus_fixture_pro",
			"countAtLeast": 4,
		}},
	})
	if !report.Pass || len(report.Results) != 1 {
		t.Fatalf("assert report = %#v, want pass", report)
	}

	status, body := postJSONStatus(t, handler, "/api/fixtures/assert", map[string]any{
		"name": "saas-basic-failing-check",
		"filter": map[string]any{
			"runId": "run-fixture-1",
		},
		"expect": []map[string]any{{
			"target": "subscription",
			"status": "canceled",
		}},
	})
	if status != http.StatusConflict || !strings.Contains(body, `"pass":false`) {
		t.Fatalf("status=%d body=%s, want assertion conflict report", status, body)
	}
}

func TestFixtureApplyPreservesStableBillingGraphIDs(t *testing.T) {
	handler := newTestHandler(t)
	pack := map[string]any{
		"name":      "sample-fixed-ids",
		"runId":     "run-fixed-ids-1",
		"namespace": "sample-e2e",
		"customers": []map[string]any{{
			"id":       "cus_fixture_fixed",
			"email":    "fixed@example.test",
			"metadata": map[string]string{"tenantId": "sample"},
		}},
		"products": []map[string]any{
			{
				"id":   "prod_fixture_plan",
				"name": "Fixture Plan",
				"metadata": map[string]string{
					"tenantId":    "sample",
					"productType": "WORKSPACE_PLAN",
					"tier":        "PREMIUM",
				},
			},
			{
				"id":   "prod_fixture_seat",
				"name": "Fixture Seat",
				"metadata": map[string]string{
					"tenantId":    "sample",
					"productType": "ADDITIONAL_SEAT",
				},
			},
		},
		"prices": []map[string]any{
			{
				"id":         "price_fixture_plan_monthly",
				"product":    "prod_fixture_plan",
				"currency":   "usd",
				"unitAmount": 30000,
				"interval":   "month",
				"metadata":   map[string]string{"tenantId": "sample", "tier": "PREMIUM"},
			},
			{
				"id":         "price_fixture_seat_monthly",
				"product":    "prod_fixture_seat",
				"currency":   "usd",
				"unitAmount": 10000,
				"interval":   "month",
				"metadata":   map[string]string{"tenantId": "sample"},
			},
		},
		"subscriptions": []map[string]any{{
			"id":              "sub_fixture_fixed",
			"checkoutSession": "cs_fixture_fixed",
			"invoice":         "in_fixture_fixed",
			"paymentIntent":   "pi_fixture_fixed",
			"customer":        "cus_fixture_fixed",
			"items": []map[string]any{
				{"price": "price_fixture_plan_monthly", "quantity": 1},
				{"price": "price_fixture_seat_monthly", "quantity": 2},
			},
			"metadata": map[string]string{
				"tenantId":               "sample",
				"manualExportLimitCount": "0",
			},
		}},
	}

	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if len(applied.CheckoutSessions) != 1 || applied.CheckoutSessions[0].ID != "cs_fixture_fixed" {
		t.Fatalf("checkout sessions = %#v, want fixed checkout session id", applied.CheckoutSessions)
	}
	if len(applied.Subscriptions) != 1 || applied.Subscriptions[0].ID != "sub_fixture_fixed" || applied.Subscriptions[0].LatestInvoiceID != "in_fixture_fixed" {
		t.Fatalf("subscriptions = %#v, want fixed subscription and invoice ids", applied.Subscriptions)
	}

	subscription := getJSON[struct {
		ID            string `json:"id"`
		LatestInvoice string `json:"latest_invoice"`
		Items         struct {
			Data []struct {
				ID       string `json:"id"`
				Quantity int64  `json:"quantity"`
				Price    struct {
					ID string `json:"id"`
				} `json:"price"`
			} `json:"data"`
		} `json:"items"`
	}](t, handler, "/v1/subscriptions/sub_fixture_fixed")
	if subscription.ID != "sub_fixture_fixed" || subscription.LatestInvoice != "in_fixture_fixed" || len(subscription.Items.Data) != 2 {
		t.Fatalf("subscription = %#v, want fixed id, latest invoice, and two items", subscription)
	}
	if subscription.Items.Data[1].Quantity != 2 || subscription.Items.Data[1].Price.ID != "price_fixture_seat_monthly" {
		t.Fatalf("seat item = %#v, want quantity 2 seat price", subscription.Items.Data[1])
	}

	list := getJSON[struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}](t, handler, "/v1/subscriptions?customer=cus_fixture_fixed&status=all")
	if len(list.Data) != 1 || list.Data[0].ID != "sub_fixture_fixed" {
		t.Fatalf("subscription list = %#v, want fixed subscription by customer", list)
	}

	invoice := getJSON[struct {
		ID            string `json:"id"`
		PaymentIntent string `json:"payment_intent"`
	}](t, handler, "/v1/invoices/in_fixture_fixed")
	if invoice.ID != "in_fixture_fixed" || invoice.PaymentIntent != "pi_fixture_fixed" {
		t.Fatalf("invoice = %#v, want fixed payment intent", invoice)
	}

	item := postForm[struct {
		ID       string `json:"id"`
		Quantity int64  `json:"quantity"`
		Price    struct {
			ID string `json:"id"`
		} `json:"price"`
	}](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {"sub_fixture_fixed"},
		"price":        {"price_fixture_seat_monthly"},
		"quantity":     {"3"},
	})
	if item.ID == "" || item.Price.ID != "price_fixture_seat_monthly" || item.Quantity != 3 {
		t.Fatalf("created subscription item = %#v", item)
	}

	deleted := deleteJSON[struct {
		ID      string `json:"id"`
		Deleted bool   `json:"deleted"`
	}](t, handler, "/v1/subscription_items/"+item.ID)
	if deleted.ID != item.ID || !deleted.Deleted {
		t.Fatalf("deleted item = %#v, want deleted marker", deleted)
	}

	createdSubscription := postForm[struct {
		ID               string `json:"id"`
		CollectionMethod string `json:"collection_method"`
		Items            struct {
			Data []struct {
				Quantity int64 `json:"quantity"`
				Price    struct {
					ID string `json:"id"`
				} `json:"price"`
			} `json:"data"`
		} `json:"items"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {"cus_fixture_fixed"},
		"items[0][price]":    {"price_fixture_plan_monthly"},
		"items[0][quantity]": {"2"},
		"collection_method":  {"send_invoice"},
	})
	if createdSubscription.ID == "" || createdSubscription.CollectionMethod != "send_invoice" || len(createdSubscription.Items.Data) != 1 {
		t.Fatalf("created subscription = %#v, want Stripe-style subscription create response", createdSubscription)
	}
	if createdSubscription.Items.Data[0].Quantity != 2 || createdSubscription.Items.Data[0].Price.ID != "price_fixture_plan_monthly" {
		t.Fatalf("created subscription item = %#v", createdSubscription.Items.Data[0])
	}
}

func TestFixtureResolveStatusAndTestClockAdvance(t *testing.T) {
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer receiver.Close()
	handler := newTestHandler(t)
	_ = postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL},
		"enabled_events": {"*"},
	})
	pack := map[string]any{
		"name":      "sample-lifecycle",
		"runId":     "run-lifecycle-1",
		"namespace": "sample-e2e",
		"test_clocks": []map[string]any{{
			"id":          "clock_e2e_lifecycle",
			"frozen_time": "2026-04-01T00:00:00Z",
		}},
		"customers": []map[string]any{
			{"id": "cus_e2e_trial", "email": "trial@example.test", "test_clock": "clock_e2e_lifecycle", "metadata": map[string]string{"tenantId": "sample"}},
			{"id": "cus_e2e_cancel", "email": "cancel@example.test", "metadata": map[string]string{"tenantId": "sample"}},
			{"id": "cus_e2e_renew", "email": "renew@example.test", "metadata": map[string]string{"tenantId": "sample"}},
			{"id": "cus_e2e_past_due", "email": "past-due@example.test", "metadata": map[string]string{"tenantId": "sample"}},
		},
		"products": []map[string]any{{
			"id":       "prod_e2e_plan",
			"name":     "Sample Pro",
			"metadata": map[string]string{"tenantId": "sample"},
		}},
		"prices": []map[string]any{{
			"id":          "price_e2e_plan_monthly",
			"product":     "prod_e2e_plan",
			"currency":    "usd",
			"unit_amount": 30000,
			"lookup_key":  "sample_plan_premium_monthly",
			"interval":    "month",
		}},
		"subscriptions": []map[string]any{
			{
				"id":               "sub_e2e_trial",
				"checkout_session": "cs_e2e_trial",
				"invoice":          "in_e2e_trial",
				"payment_intent":   "pi_e2e_trial",
				"ref":              "trial-to-active",
				"customer":         "cus_e2e_trial",
				"price":            "price_e2e_plan_monthly",
				"status":           "trialing",
				"trial_start":      "2026-04-01T00:00:00Z",
				"trial_end":        "2026-04-15T00:00:00Z",
				"test_clock":       "clock_e2e_lifecycle",
			},
			{
				"id":                   "sub_e2e_cancel",
				"ref":                  "cancel-at-period-end",
				"customer":             "cus_e2e_cancel",
				"price":                "price_e2e_plan_monthly",
				"status":               "active",
				"current_period_start": "2026-04-01T00:00:00Z",
				"current_period_end":   "2026-04-15T00:00:00Z",
				"cancel_at_period_end": true,
				"test_clock":           "clock_e2e_lifecycle",
			},
			{
				"id":                   "sub_e2e_renew",
				"ref":                  "renewal-success",
				"customer":             "cus_e2e_renew",
				"price":                "price_e2e_plan_monthly",
				"status":               "active",
				"current_period_start": "2026-04-01T00:00:00Z",
				"current_period_end":   "2026-04-15T00:00:00Z",
				"test_clock":           "clock_e2e_lifecycle",
			},
			{
				"id":                    "sub_e2e_past_due",
				"ref":                   "past-due-pro",
				"customer":              "cus_e2e_past_due",
				"price":                 "price_e2e_plan_monthly",
				"status":                "past_due",
				"latest_invoice_status": "open",
				"metadata":              map[string]string{"tenantId": "sample"},
			},
		},
	}
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["test_clocks"] != 1 || applied.Summary["subscriptions"] != 4 {
		t.Fatalf("apply summary = %#v", applied.Summary)
	}

	updated := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name":  "sample-lifecycle",
		"runId": "run-lifecycle-1",
		"test_clocks": []map[string]any{{
			"id":          "clock_e2e_lifecycle",
			"name":        "Lifecycle clock updated",
			"frozen_time": "2026-04-02T00:00:00Z",
		}},
	})
	if len(updated.TestClocks) != 1 || updated.TestClocks[0].Name != "Lifecycle clock updated" || updated.TestClocks[0].FrozenTime.Format("2006-01-02") != "2026-04-02" {
		t.Fatalf("updated test clock = %#v", updated.TestClocks)
	}
	directClockSub := postForm[struct {
		ID                 string `json:"id"`
		TestClock          string `json:"test_clock"`
		CurrentPeriodStart int64  `json:"current_period_start"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {"cus_e2e_trial"},
		"items[0][price]": {"price_e2e_plan_monthly"},
		"test_clock":      {"clock_e2e_lifecycle"},
	})
	if directClockSub.ID == "" || directClockSub.TestClock != "clock_e2e_lifecycle" || directClockSub.CurrentPeriodStart != 1775088000 {
		t.Fatalf("direct clock subscription = %#v, want creation at frozen clock time", directClockSub)
	}
	sessionsBeforeInvalidClock := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/checkout/sessions")
	status, _ := postFormStatus(t, handler, "/v1/subscriptions", url.Values{
		"customer":        {"cus_e2e_trial"},
		"items[0][price]": {"price_e2e_plan_monthly"},
		"test_clock":      {"clock_missing"},
	})
	if status != http.StatusNotFound {
		t.Fatalf("invalid test_clock subscription status = %d, want 404", status)
	}
	sessionsAfterInvalidClock := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/checkout/sessions")
	if len(sessionsAfterInvalidClock.Data) != len(sessionsBeforeInvalidClock.Data) {
		t.Fatalf("checkout sessions before=%d after=%d, invalid test_clock should not create session", len(sessionsBeforeInvalidClock.Data), len(sessionsAfterInvalidClock.Data))
	}

	createdEvents := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=customer.subscription.created")
	if len(createdEvents.Data) == 0 {
		t.Fatal("fixture apply did not create customer.subscription.created events")
	}
	replay := postJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/api/events/"+createdEvents.Data[0].ID+"/replay", map[string]any{
		"duplicate": 1,
	})
	if len(replay.Data) == 0 {
		t.Fatalf("fixture-created event replay = %#v, want delivery attempts", replay)
	}

	resolved := getJSON[fixtures.ResolveResult](t, handler, "/api/fixtures/resolve?ref=trial-to-active&runId=run-lifecycle-1")
	if resolved.CustomerID != "cus_e2e_trial" || resolved.SubscriptionID != "sub_e2e_trial" || resolved.InvoiceID != "in_e2e_trial" || resolved.PaymentIntentID != "pi_e2e_trial" || resolved.CheckoutSessionID != "cs_e2e_trial" {
		t.Fatalf("resolved = %#v, want stable fixture graph ids", resolved)
	}
	resolvedPrice := getJSON[fixtures.ResolveResult](t, handler, "/api/fixtures/resolve?lookup_key=sample_plan_premium_monthly&runId=run-lifecycle-1")
	if resolvedPrice.PriceID != "price_e2e_plan_monthly" || resolvedPrice.ProductID != "prod_e2e_plan" {
		t.Fatalf("resolved price = %#v, want lookup_key to resolve price/product", resolvedPrice)
	}

	filtered := getJSON[struct {
		Data []struct {
			ID       string            `json:"id"`
			Metadata map[string]string `json:"metadata"`
		} `json:"data"`
	}](t, handler, "/v1/subscriptions?metadata[billtap_fixture_ref]=trial-to-active")
	if len(filtered.Data) != 1 || filtered.Data[0].ID != "sub_e2e_trial" {
		t.Fatalf("metadata filtered subscriptions = %#v", filtered.Data)
	}

	invoiceCreatedBefore := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=invoice.created")
	advance := postForm[struct {
		ID      string `json:"id"`
		Status  string `json:"status"`
		Billtap struct {
			ActivatedCount int `json:"activated_count"`
			CanceledCount  int `json:"canceled_count"`
			Renewed        int `json:"renewed"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_e2e_lifecycle/advance", url.Values{
		"frozen_time": {"1776297600"},
	})
	if advance.ID != "clock_e2e_lifecycle" || advance.Status != "ready" || advance.Billtap.ActivatedCount != 1 || advance.Billtap.CanceledCount != 1 || advance.Billtap.Renewed != 1 {
		t.Fatalf("advance = %#v, want one trial activation, one renewal, and one cancellation", advance)
	}
	invoiceCreatedAfter := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=invoice.created")
	if len(invoiceCreatedAfter.Data) <= len(invoiceCreatedBefore.Data) {
		t.Fatalf("invoice.created events before=%d after=%d, want renewal event emitted", len(invoiceCreatedBefore.Data), len(invoiceCreatedAfter.Data))
	}
	active := getJSON[struct {
		Status string `json:"status"`
	}](t, handler, "/v1/subscriptions/sub_e2e_trial")
	if active.Status != "active" {
		t.Fatalf("trial subscription status = %q, want active", active.Status)
	}
	canceled := getJSON[struct {
		Status string `json:"status"`
	}](t, handler, "/v1/subscriptions/sub_e2e_cancel")
	if canceled.Status != "canceled" {
		t.Fatalf("canceled subscription status = %q, want canceled", canceled.Status)
	}
	pastDue := getJSON[struct {
		Status   string            `json:"status"`
		Metadata map[string]string `json:"metadata"`
	}](t, handler, "/v1/subscriptions/sub_e2e_past_due")
	if pastDue.Status != "past_due" || pastDue.Metadata["latest_invoice_status"] != "open" {
		t.Fatalf("past due subscription = %#v", pastDue)
	}
}

func TestFixtureSubscriptionStatusAndDatesOverrideOutcome(t *testing.T) {
	handler := newTestHandler(t)
	pack := map[string]any{
		"name":      "sample-status-override",
		"runId":     "run-status-override-1",
		"namespace": "sample-e2e",
		"test_clocks": []map[string]any{{
			"id":          "clock_fixture_trial",
			"frozen_time": "2030-01-10T00:00:00Z",
		}},
		"customers": []map[string]any{
			{"id": "cus_fixture_trial", "email": "trial-status@example.test", "test_clock": "clock_fixture_trial"},
			{"id": "cus_fixture_canceled", "email": "canceled-status@example.test"},
			{"id": "cus_fixture_past_due", "email": "past-due-status@example.test"},
			{"id": "cus_fixture_unpaid", "email": "unpaid-status@example.test"},
			{"id": "cus_fixture_incomplete", "email": "incomplete-status@example.test"},
		},
		"products": []map[string]any{{
			"id":   "prod_fixture_status_plan",
			"name": "Status Override Plan",
		}},
		"prices": []map[string]any{{
			"id":          "price_fixture_status_monthly",
			"product":     "prod_fixture_status_plan",
			"currency":    "usd",
			"unit_amount": 30000,
			"interval":    "month",
		}},
		"subscriptions": []map[string]any{
			{
				"id":                   "sub_fixture_trial",
				"ref":                  "trial-status-override",
				"customer":             "cus_fixture_trial",
				"price":                "price_fixture_status_monthly",
				"quantity":             1,
				"outcome":              "payment_succeeded",
				"status":               "trialing",
				"current_period_start": "2030-01-10T00:00:00Z",
				"current_period_end":   "2030-01-15T00:00:00Z",
				"trial_start":          "2030-01-10T00:00:00Z",
				"trial_end":            "2030-01-15T00:00:00Z",
				"test_clock":           "clock_fixture_trial",
			},
			{
				"id":                   "sub_fixture_canceled",
				"customer":             "cus_fixture_canceled",
				"price":                "price_fixture_status_monthly",
				"outcome":              "payment_succeeded",
				"status":               "canceled",
				"current_period_start": "2030-01-01T00:00:00Z",
				"current_period_end":   "2030-01-10T00:00:00Z",
				"canceled_at":          "2030-01-10T00:00:00Z",
				"ended_at":             "2030-01-10T00:00:00Z",
			},
			{
				"id":                    "sub_fixture_past_due",
				"customer":              "cus_fixture_past_due",
				"price":                 "price_fixture_status_monthly",
				"outcome":               "payment_succeeded",
				"status":                "past_due",
				"current_period_start":  "2030-01-01T00:00:00Z",
				"current_period_end":    "2030-02-01T00:00:00Z",
				"latest_invoice_status": "open",
			},
			{
				"id":                    "sub_fixture_unpaid",
				"customer":              "cus_fixture_unpaid",
				"price":                 "price_fixture_status_monthly",
				"outcome":               "payment_succeeded",
				"status":                "unpaid",
				"current_period_start":  "2030-01-01T00:00:00Z",
				"current_period_end":    "2030-02-01T00:00:00Z",
				"latest_invoice_status": "open",
			},
			{
				"id":                    "sub_fixture_incomplete",
				"customer":              "cus_fixture_incomplete",
				"price":                 "price_fixture_status_monthly",
				"outcome":               "payment_succeeded",
				"status":                "incomplete",
				"current_period_start":  "2030-01-01T00:00:00Z",
				"current_period_end":    "2030-02-01T00:00:00Z",
				"latest_invoice_status": "open",
			},
		},
	}

	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["test_clocks"] != 1 || applied.Summary["subscriptions"] != 5 {
		t.Fatalf("apply summary = %#v", applied.Summary)
	}

	trialList := getJSON[struct {
		Data []struct {
			ID                 string `json:"id"`
			Status             string `json:"status"`
			CurrentPeriodStart int64  `json:"current_period_start"`
			CurrentPeriodEnd   int64  `json:"current_period_end"`
			TrialStart         int64  `json:"trial_start"`
			TrialEnd           int64  `json:"trial_end"`
		} `json:"data"`
	}](t, handler, "/v1/subscriptions?customer=cus_fixture_trial&status=all")
	if len(trialList.Data) != 1 {
		t.Fatalf("trial list = %#v, want one subscription", trialList.Data)
	}
	trial := trialList.Data[0]
	if trial.ID != "sub_fixture_trial" || trial.Status != "trialing" || trial.CurrentPeriodStart != 1894233600 || trial.CurrentPeriodEnd != 1894665600 || trial.TrialStart != 1894233600 || trial.TrialEnd != 1894665600 {
		t.Fatalf("trial subscription = %#v, want fixture status and dates preserved", trial)
	}

	for _, tt := range []struct {
		id     string
		status string
	}{
		{"sub_fixture_canceled", "canceled"},
		{"sub_fixture_past_due", "past_due"},
		{"sub_fixture_unpaid", "unpaid"},
		{"sub_fixture_incomplete", "incomplete"},
	} {
		subscription := getJSON[struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}](t, handler, "/v1/subscriptions/"+tt.id)
		if subscription.ID != tt.id || subscription.Status != tt.status {
			t.Fatalf("%s = %#v, want status %q despite payment_succeeded outcome", tt.id, subscription, tt.status)
		}
	}

	if status := subscriptionCreatedEventStatus(t, handler, "sub_fixture_trial"); status != "trialing" {
		t.Fatalf("created event status = %q, want trialing", status)
	}

	advance := postForm[struct {
		Billtap struct {
			ActivatedCount int `json:"activated_count"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_fixture_trial/advance", url.Values{
		"frozen_time": {"1894752000"},
	})
	if advance.Billtap.ActivatedCount != 1 {
		t.Fatalf("first advance = %#v, want one trial activation", advance)
	}
	active := getJSON[struct {
		Status string `json:"status"`
	}](t, handler, "/v1/subscriptions/sub_fixture_trial")
	if active.Status != "active" {
		t.Fatalf("advanced trial status = %q, want active", active.Status)
	}

	_ = postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	reset := getJSON[struct {
		Status           string `json:"status"`
		CurrentPeriodEnd int64  `json:"current_period_end"`
		TrialEnd         int64  `json:"trial_end"`
	}](t, handler, "/v1/subscriptions/sub_fixture_trial")
	if reset.Status != "trialing" || reset.CurrentPeriodEnd != 1894665600 || reset.TrialEnd != 1894665600 {
		t.Fatalf("reset trial subscription = %#v, want fixture status and dates restored", reset)
	}
}

func TestFixtureApplyBackfillsCreatedEventForExistingSubscriptionSeed(t *testing.T) {
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
	billingService := billing.NewService(store)
	handler := New(Options{
		Billing:     billingService,
		Webhooks:    webhooks.NewService(store),
		Diagnostics: diagnostics.NewService(store),
	})

	pack := fixtures.Pack{
		Name:      "sample-existing-subscription",
		RunID:     "run-existing-created-1",
		Namespace: "sample-e2e",
		Customers: []fixtures.CustomerFixture{{
			ID:    "cus_fixture_existing_created",
			Email: "existing-created@example.test",
		}},
		Products: []fixtures.ProductFixture{{
			ID:   "prod_fixture_existing_plan",
			Name: "Existing Fixture Plan",
		}},
		Prices: []fixtures.PriceFixture{{
			ID:         "price_fixture_existing_monthly",
			Product:    "prod_fixture_existing_plan",
			Currency:   "usd",
			UnitAmount: 30000,
			LookupKey:  "sample_existing_monthly",
			Interval:   "month",
		}},
		Subscriptions: []fixtures.SubscriptionFixture{{
			ID:                 "sub_fixture_existing_created",
			CheckoutSessionID:  "cs_fixture_existing_created",
			InvoiceID:          "in_fixture_existing_created",
			PaymentIntentID:    "pi_fixture_existing_created",
			Ref:                "existing-created-ref",
			Customer:           "cus_fixture_existing_created",
			Price:              "price_fixture_existing_monthly",
			Status:             "active",
			CurrentPeriodStart: "2026-05-01T00:00:00Z",
			CurrentPeriodEnd:   "2026-06-01T00:00:00Z",
		}},
	}

	if _, err := fixtures.NewService(billingService).Apply(ctx, pack); err != nil {
		t.Fatalf("preseed fixture graph: %v", err)
	}
	if count, _ := countSubscriptionCreatedEvents(t, handler, "sub_fixture_existing_created"); count != 0 {
		t.Fatalf("created events before HTTP apply = %d, want 0", count)
	}

	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if len(applied.CheckoutSessions) != 0 || len(applied.Subscriptions) != 1 {
		t.Fatalf("apply result sessions=%#v subscriptions=%#v, want existing subscription path", applied.CheckoutSessions, applied.Subscriptions)
	}
	if count, source := countSubscriptionCreatedEvents(t, handler, "sub_fixture_existing_created"); count != 1 || source != webhooks.SourceFixture {
		t.Fatalf("created events count=%d source=%q, want one fixture backfill", count, source)
	}

	_ = postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if count, _ := countSubscriptionCreatedEvents(t, handler, "sub_fixture_existing_created"); count != 1 {
		t.Fatalf("created events after reapply = %d, want idempotent backfill", count)
	}
}

func countSubscriptionCreatedEvents(t *testing.T, handler http.Handler, subscriptionID string) (int, string) {
	t.Helper()
	events := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=customer.subscription.created")
	count := 0
	source := ""
	for _, event := range events.Data {
		var object struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(event.Data.Object, &object); err != nil {
			t.Fatalf("decode event object: %v", err)
		}
		if object.ID == subscriptionID {
			count++
			source = event.Billtap.Source
		}
	}
	return count, source
}

func subscriptionCreatedEventStatus(t *testing.T, handler http.Handler, subscriptionID string) string {
	t.Helper()
	events := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=customer.subscription.created")
	for _, event := range events.Data {
		var object struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal(event.Data.Object, &object); err != nil {
			t.Fatalf("decode event object: %v", err)
		}
		if object.ID == subscriptionID {
			return object.Status
		}
	}
	t.Fatalf("missing customer.subscription.created event for %s", subscriptionID)
	return ""
}

func TestRefundCreditNoteAPIsAndEvents(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"refund@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"30000"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	completion := postJSON[struct {
		Invoice       billing.Invoice       `json:"invoice"`
		PaymentIntent billing.PaymentIntent `json:"payment_intent"`
	}](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})

	refund := postForm[struct {
		ID            string `json:"id"`
		Charge        string `json:"charge"`
		Invoice       string `json:"invoice"`
		PaymentIntent string `json:"payment_intent"`
		Status        string `json:"status"`
	}](t, handler, "/v1/refunds", url.Values{
		"invoice": {completion.Invoice.ID},
		"amount":  {"15000"},
		"reason":  {"requested_by_customer"},
	})
	if refund.ID == "" || refund.Charge == "" || refund.Invoice != completion.Invoice.ID || refund.PaymentIntent != completion.PaymentIntent.ID || refund.Status != "succeeded" {
		t.Fatalf("refund = %#v, want linked refund evidence", refund)
	}

	note := postForm[struct {
		ID      string `json:"id"`
		Invoice string `json:"invoice"`
		Status  string `json:"status"`
	}](t, handler, "/v1/credit_notes", url.Values{
		"invoice": {completion.Invoice.ID},
		"amount":  {"15000"},
		"reason":  {"order_change"},
	})
	if note.ID == "" || note.Invoice != completion.Invoice.ID || note.Status != "issued" {
		t.Fatalf("credit note = %#v, want issued note", note)
	}

	refundEvents := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=charge.refunded")
	creditEvents := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=credit_note.created")
	if len(refundEvents.Data) != 1 || len(creditEvents.Data) != 1 {
		t.Fatalf("refund events=%d credit events=%d, want one each", len(refundEvents.Data), len(creditEvents.Data))
	}
	var chargePayload struct {
		ID     string `json:"id"`
		Object string `json:"object"`
	}
	if err := json.Unmarshal(refundEvents.Data[0].Data.Object, &chargePayload); err != nil {
		t.Fatalf("decode charge.refunded payload: %v", err)
	}
	if chargePayload.Object != "charge" || chargePayload.ID != refund.Charge {
		t.Fatalf("charge.refunded payload = %#v, want charge-shaped payload", chargePayload)
	}

	fixturePack := map[string]any{
		"name":  "refund-credit-fixture",
		"runId": "refund-run-1",
		"refunds": []map[string]any{{
			"id":      "re_e2e_partial_refund",
			"invoice": completion.Invoice.ID,
			"amount":  5000,
			"reason":  "requested_by_customer",
		}},
		"credit_notes": []map[string]any{{
			"id":      "cn_e2e_partial_credit",
			"invoice": completion.Invoice.ID,
			"amount":  5000,
			"reason":  "order_change",
		}},
	}
	firstApply := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", fixturePack)
	secondApply := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", fixturePack)
	if len(firstApply.Refunds) != 1 || len(secondApply.Refunds) != 1 || secondApply.Refunds[0].ID != "re_e2e_partial_refund" {
		t.Fatalf("fixture refunds first=%#v second=%#v, want idempotent stable refund", firstApply.Refunds, secondApply.Refunds)
	}
	if len(firstApply.CreditNotes) != 1 || len(secondApply.CreditNotes) != 1 || secondApply.CreditNotes[0].ID != "cn_e2e_partial_credit" {
		t.Fatalf("fixture credit notes first=%#v second=%#v, want idempotent stable credit note", firstApply.CreditNotes, secondApply.CreditNotes)
	}
}

func TestWebhookReplaySimulateAppFailureThenDeliver(t *testing.T) {
	calls := 0
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer receiver.Close()

	handler := newTestHandler(t)
	_ = postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {receiver.URL},
		"enabled_events": {"checkout.session.completed"},
	})
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"webhook-failure@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Team"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":     {product.ID},
		"currency":    {"usd"},
		"unit_amount": {"30000"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"line_items[0][price]": {price.ID},
	})
	_ = postJSON[map[string]any](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	events := getJSON[struct {
		Data []webhooks.Event `json:"data"`
	}](t, handler, "/v1/events?type=checkout.session.completed")
	if len(events.Data) == 0 {
		t.Fatal("checkout event not found")
	}
	beforeReplayCalls := calls
	replay := postJSON[struct {
		Data []webhooks.DeliveryAttempt `json:"data"`
	}](t, handler, "/api/events/"+events.Data[0].ID+"/replay", map[string]any{
		"simulate_app_failure": map[string]any{
			"status":                502,
			"fail_first_n_attempts": 1,
			"body":                  "Upstream timeout",
		},
	})
	if len(replay.Data) != 2 {
		t.Fatalf("replay attempts = %#v, want failed attempt followed by delivered attempt", replay.Data)
	}
	if replay.Data[0].Status != webhooks.StatusFailed || replay.Data[0].ResponseStatus != 502 || replay.Data[0].Metadata["simulate_app_failure"] != "true" {
		t.Fatalf("first replay attempt = %#v, want simulated 502 failure metadata", replay.Data[0])
	}
	if replay.Data[1].Status != webhooks.StatusSucceeded || replay.Data[1].ResponseStatus != http.StatusOK {
		t.Fatalf("second replay attempt = %#v, want actual receiver success", replay.Data[1])
	}
	if calls != beforeReplayCalls+1 {
		t.Fatalf("receiver calls = %d, before replay %d; simulated failure should not call receiver", calls, beforeReplayCalls)
	}
}

func TestWebhookEndpointPatchPreservesActiveUnlessEnabledChanges(t *testing.T) {
	handler := newTestHandler(t)
	endpoint := postForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints", url.Values{
		"url":            {"http://example.test/a"},
		"enabled_events": {"*"},
	})

	patched := patchForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints/"+endpoint.ID, url.Values{
		"enabled": {"false"},
	})
	if patched.Active {
		t.Fatalf("patched endpoint active = true, want disabled")
	}
	renamed := patchForm[webhooks.Endpoint](t, handler, "/v1/webhook_endpoints/"+endpoint.ID, url.Values{
		"url": {"http://example.test/b"},
	})
	if renamed.Active || renamed.URL != "http://example.test/b" {
		t.Fatalf("renamed endpoint = %#v, want inactive endpoint with updated url", renamed)
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	return newTestHandlerWithOptions(t, Options{})
}

func newTestHandlerWithOptions(t *testing.T, opts Options) http.Handler {
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
	opts.Billing = billing.NewService(store)
	opts.Webhooks = webhooks.NewService(store)
	opts.Diagnostics = diagnostics.NewService(store)
	return New(opts)
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func postForm[T any](t *testing.T, handler http.Handler, path string, values url.Values) T {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, stringsReader(values.Encode()))
	req.Host = "billtap.test"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return decodeResponse[T](t, rec)
}

func postJSON[T any](t *testing.T, handler http.Handler, path string, body any) T {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return decodeResponse[T](t, rec)
}

func postJSONStatus(t *testing.T, handler http.Handler, path string, body any) (int, string) {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func postFormStatus(t *testing.T, handler http.Handler, path string, values url.Values) (int, string) {
	t.Helper()
	return postFormStatusWithHeaders(t, handler, path, values, nil)
}

func postFormStatusWithHeaders(t *testing.T, handler http.Handler, path string, values url.Values, headers map[string]string) (int, string) {
	t.Helper()
	status, body, _ := postFormStatusWithResponseHeaders(t, handler, path, values, headers)
	return status, body
}

func postFormStatusWithResponseHeaders(t *testing.T, handler http.Handler, path string, values url.Values, headers map[string]string) (int, string, http.Header) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, stringsReader(values.Encode()))
	req.Host = "billtap.test"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String(), rec.Header()
}

func getStatus(t *testing.T, handler http.Handler, path string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func patchForm[T any](t *testing.T, handler http.Handler, path string, values url.Values) T {
	t.Helper()
	req := httptest.NewRequest(http.MethodPatch, path, stringsReader(values.Encode()))
	req.Host = "billtap.test"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return decodeResponse[T](t, rec)
}

func getJSON[T any](t *testing.T, handler http.Handler, path string) T {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return decodeResponse[T](t, rec)
}

func deleteJSON[T any](t *testing.T, handler http.Handler, path string) T {
	t.Helper()
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return decodeResponse[T](t, rec)
}

type stripeAPIErrorBody struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Param   string `json:"param,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

func decodeErrorBody(t *testing.T, body string) stripeAPIErrorBody {
	t.Helper()
	var out stripeAPIErrorBody
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		t.Fatalf("decode error body %s: %v", body, err)
	}
	return out
}

func decodeResponse[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}
	var out T
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response %s: %v", rec.Body.String(), err)
	}
	return out
}

func stringsReader(value string) *bytes.Reader {
	return bytes.NewReader([]byte(value))
}
