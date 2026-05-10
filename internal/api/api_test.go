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
	"strings"
	"testing"

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
	if portal.URL != "http://127.0.0.1:18080/portal?customer_id="+customer.ID {
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
		"name":      "ds5-fixed-ids",
		"runId":     "run-fixed-ids-1",
		"namespace": "ds5",
		"customers": []map[string]any{{
			"id":       "cus_fixture_fixed",
			"email":    "fixed@example.test",
			"metadata": map[string]string{"tenantId": "dentbird"},
		}},
		"products": []map[string]any{
			{
				"id":   "prod_fixture_plan",
				"name": "Fixture Plan",
				"metadata": map[string]string{
					"tenantId":    "dentbird",
					"productType": "WORKSPACE_PLAN",
					"tier":        "PREMIUM",
				},
			},
			{
				"id":   "prod_fixture_seat",
				"name": "Fixture Seat",
				"metadata": map[string]string{
					"tenantId":    "dentbird",
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
				"metadata":   map[string]string{"tenantId": "dentbird", "tier": "PREMIUM"},
			},
			{
				"id":         "price_fixture_seat_monthly",
				"product":    "prod_fixture_seat",
				"currency":   "usd",
				"unitAmount": 10000,
				"interval":   "month",
				"metadata":   map[string]string{"tenantId": "dentbird"},
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
				"tenantId":               "dentbird",
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
	req := httptest.NewRequest(http.MethodPost, path, stringsReader(values.Encode()))
	req.Host = "billtap.test"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
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
