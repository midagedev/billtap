package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hckim/billtap/internal/billing"
)

func setupDraftInvoiceWithLine(t *testing.T, handler http.Handler) (customer billing.Customer, invoiceID, lineID string) {
	t.Helper()
	customer = postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"inv-final@example.test"}})
	invoice := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/invoices", url.Values{
		"customer": {customer.ID},
	})
	line := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/invoiceitems", url.Values{
		"customer":    {customer.ID},
		"invoice":     {invoice.ID},
		"amount":      {"4000"},
		"currency":    {"usd"},
		"description": {"Original"},
	})
	return customer, invoice.ID, line.ID
}

func TestInvoiceSingleLineUpdate(t *testing.T) {
	handler := newTestHandler(t)
	_, invoiceID, lineID := setupDraftInvoiceWithLine(t, handler)

	// Unknown params are rejected.
	status, body := postFormStatus(t, handler, "/v1/invoices/"+invoiceID+"/lines/"+lineID, url.Values{
		"quantity": {"2"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("unknown line update param status = %d body = %s, want 400", status, body)
	}
	// Missing lines are 404.
	status, body = postFormStatus(t, handler, "/v1/invoices/"+invoiceID+"/lines/ii_missing", url.Values{
		"amount": {"1000"},
	})
	if status != http.StatusNotFound {
		t.Fatalf("missing line update status = %d body = %s, want 404", status, body)
	}

	updated := postForm[struct {
		Subtotal int64 `json:"subtotal"`
	}](t, handler, "/v1/invoices/"+invoiceID+"/lines/"+lineID, url.Values{
		"amount":      {"2500"},
		"description": {"Revised"},
		"metadata[k]": {"v"},
	})
	if updated.Subtotal != 2500 {
		t.Fatalf("single line update subtotal = %d, want 2500", updated.Subtotal)
	}
	lines := getJSON[struct {
		Data []struct {
			ID          string `json:"id"`
			Amount      int64  `json:"amount"`
			Description string `json:"description"`
		} `json:"data"`
	}](t, handler, "/v1/invoices/"+invoiceID+"/lines")
	if len(lines.Data) != 1 || lines.Data[0].ID != lineID || lines.Data[0].Amount != 2500 || lines.Data[0].Description != "Revised" {
		t.Fatalf("line after single update = %#v, want revised %s", lines.Data, lineID)
	}

	// Non-draft invoices reject line edits.
	_ = postJSON[map[string]json.RawMessage](t, handler, "/v1/invoices/"+invoiceID+"/finalize", map[string]string{})
	status, body = postFormStatus(t, handler, "/v1/invoices/"+invoiceID+"/lines/"+lineID, url.Values{
		"amount": {"1000"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("finalized invoice line update status = %d body = %s, want 400", status, body)
	}
}

func TestInvoiceAttachPayment(t *testing.T) {
	handler := newTestHandler(t)
	customer, invoiceID, _ := setupDraftInvoiceWithLine(t, handler)

	intent := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"customer": {customer.ID},
		"amount":   {"4000"},
		"currency": {"usd"},
	})

	// payment_intent or payment_record is required.
	status, body := postFormStatus(t, handler, "/v1/invoices/"+invoiceID+"/attach_payment", url.Values{})
	if status != http.StatusBadRequest {
		t.Fatalf("attach_payment without params status = %d body = %s, want 400", status, body)
	}
	// Unknown intents are 404.
	status, body = postFormStatus(t, handler, "/v1/invoices/"+invoiceID+"/attach_payment", url.Values{
		"payment_intent": {"pi_missing"},
	})
	if status != http.StatusNotFound {
		t.Fatalf("attach_payment unknown intent status = %d body = %s, want 404", status, body)
	}

	attached := postForm[struct {
		Metadata map[string]string `json:"metadata"`
	}](t, handler, "/v1/invoices/"+invoiceID+"/attach_payment", url.Values{
		"payment_intent": {intent.ID},
	})
	if attached.Metadata["billtap_attached_payment_intent"] != intent.ID {
		t.Fatalf("attach_payment metadata = %#v, want recorded %s", attached.Metadata, intent.ID)
	}
}

func TestSubscriptionMigrate(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"migrate@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Migrate Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1900"},
		"recurring[interval]": {"month"},
	})
	subscription := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {price.ID},
		"items[0][quantity]": {"1"},
	})

	// billing_mode[type] is required and enum-bound.
	status, body := postFormStatus(t, handler, "/v1/subscriptions/"+subscription.ID+"/migrate", url.Values{})
	if status != http.StatusBadRequest {
		t.Fatalf("migrate without billing_mode status = %d body = %s, want 400", status, body)
	}
	status, body = postFormStatus(t, handler, "/v1/subscriptions/"+subscription.ID+"/migrate", url.Values{
		"billing_mode[type]": {"legacy"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("migrate invalid billing_mode status = %d body = %s, want 400", status, body)
	}

	migrated := postForm[struct {
		ID       string            `json:"id"`
		Metadata map[string]string `json:"metadata"`
	}](t, handler, "/v1/subscriptions/"+subscription.ID+"/migrate", url.Values{
		"billing_mode[type]": {"flexible"},
	})
	if migrated.ID != subscription.ID || migrated.Metadata["billtap_billing_mode"] != "flexible" {
		t.Fatalf("migrated subscription = %#v, want billing_mode=flexible evidence", migrated)
	}
}

func TestPaymentIntentUpdate(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"pi-update@example.test"}})
	intent := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/payment_intents", url.Values{
		"customer": {customer.ID},
		"amount":   {"1200"},
		"currency": {"usd"},
	})

	// Unknown params are rejected.
	status, body := postFormStatus(t, handler, "/v1/payment_intents/"+intent.ID, url.Values{
		"amount": {"9999"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("intent amount update status = %d body = %s, want 400", status, body)
	}

	updated := postForm[struct {
		Description string            `json:"description"`
		Metadata    map[string]string `json:"metadata"`
	}](t, handler, "/v1/payment_intents/"+intent.ID, url.Values{
		"description":   {"Widget order"},
		"metadata[env]": {"ci"},
	})
	if updated.Description != "Widget order" {
		t.Fatalf("updated description = %q, want Widget order", updated.Description)
	}
	if updated.Metadata["env"] != "ci" {
		t.Fatalf("updated metadata = %#v, want env=ci", updated.Metadata)
	}
	fetched := getJSON[struct {
		Description string `json:"description"`
	}](t, handler, "/v1/payment_intents/"+intent.ID)
	if fetched.Description != "Widget order" {
		t.Fatalf("persisted description = %q, want Widget order", fetched.Description)
	}
}

func TestCreditNoteLines(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"cn-lines@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"CN Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"3000"},
		"recurring[interval]": {"month"},
	})
	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {price.ID},
		"line_items[0][quantity]": {"1"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var completed billing.CheckoutSession
	if err := json.Unmarshal(completion["session"], &completed); err != nil {
		t.Fatalf("decode completed session: %v", err)
	}

	note := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/credit_notes", url.Values{
		"invoice": {completed.InvoiceID},
		"amount":  {"1000"},
		"memo":    {"Partial refund credit"},
	})
	lines := getJSON[struct {
		Data []struct {
			ID          string `json:"id"`
			Object      string `json:"object"`
			Amount      int64  `json:"amount"`
			Description string `json:"description"`
		} `json:"data"`
	}](t, handler, "/v1/credit_notes/"+note.ID+"/lines")
	if len(lines.Data) != 1 {
		t.Fatalf("credit note lines = %d, want 1", len(lines.Data))
	}
	line := lines.Data[0]
	if line.Object != "credit_note_line_item" || line.Amount != -1000 || line.Description != "Partial refund credit" {
		t.Fatalf("credit note line = %#v, want negative 1000 with memo description", line)
	}

	missingReq := httptest.NewRequest(http.MethodGet, "/v1/credit_notes/cn_missing/lines", nil)
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing credit note lines status = %d body = %s, want 404", missingRec.Code, missingRec.Body.String())
	}
}
