package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hckim/billtap/internal/billing"
)

func setupCheckoutLineFixtures(t *testing.T, handler http.Handler) (customer billing.Customer, price billing.Price) {
	t.Helper()
	customer = postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"cs-lines@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Line Item Plan"}})
	price = postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"2500"},
		"recurring[interval]": {"month"},
	})
	return customer, price
}

func TestCheckoutSessionLineItems(t *testing.T) {
	handler := newTestHandler(t)
	customer, price := setupCheckoutLineFixtures(t, handler)

	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {price.ID},
		"line_items[0][quantity]": {"1"},
		"line_items[1][price]":    {price.ID},
		"line_items[1][quantity]": {"2"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})

	listed := getJSON[struct {
		Data []struct {
			ID             string `json:"id"`
			Object         string `json:"object"`
			Quantity       int64  `json:"quantity"`
			AmountTotal    int64  `json:"amount_total"`
			AmountSubtotal int64  `json:"amount_subtotal"`
			Currency       string `json:"currency"`
			Price          struct {
				ID string `json:"id"`
			} `json:"price"`
		} `json:"data"`
	}](t, handler, "/v1/checkout/sessions/"+session.ID+"/line_items")
	if len(listed.Data) != 2 {
		t.Fatalf("session line_items = %d, want 2", len(listed.Data))
	}
	for idx, item := range listed.Data {
		if item.Object != "item" || item.Price.ID != price.ID {
			t.Fatalf("line item %d = %#v, want expanded price %s", idx, item, price.ID)
		}
	}
	if listed.Data[0].Quantity != 1 || listed.Data[0].AmountTotal != 2500 || listed.Data[0].AmountSubtotal != 2500 {
		t.Fatalf("first line item = %#v, want quantity 1 amount 2500", listed.Data[0])
	}
	if listed.Data[1].Quantity != 2 || listed.Data[1].AmountTotal != 5000 {
		t.Fatalf("second line item = %#v, want quantity 2 amount 5000", listed.Data[1])
	}

	missingReq := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions/cs_missing/line_items", nil)
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing session line_items status = %d body = %s, want 404", missingRec.Code, missingRec.Body.String())
	}
}

func TestCheckoutSessionUpdateMetadataAndQuantities(t *testing.T) {
	handler := newTestHandler(t)
	customer, price := setupCheckoutLineFixtures(t, handler)

	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {price.ID},
		"line_items[0][quantity]": {"1"},
		"line_items[1][price]":    {price.ID},
		"line_items[1][quantity]": {"1"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})

	// Unknown update params are rejected.
	status, body := postFormStatus(t, handler, "/v1/checkout/sessions/"+session.ID, url.Values{
		"shipping_options[0][shipping_rate]": {"shr_1"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("unknown update param status = %d body = %s, want 400", status, body)
	}
	// Out-of-range line index is rejected.
	status, body = postFormStatus(t, handler, "/v1/checkout/sessions/"+session.ID, url.Values{
		"line_items[5][quantity]": {"2"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("out-of-range line_items status = %d body = %s, want 400", status, body)
	}
	// Quantities must stay positive.
	status, body = postFormStatus(t, handler, "/v1/checkout/sessions/"+session.ID, url.Values{
		"line_items[0][quantity]": {"0"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("zero quantity status = %d body = %s, want 400", status, body)
	}

	updated := postForm[struct {
		Metadata    map[string]string `json:"metadata"`
		AmountTotal int64             `json:"amount_total"`
	}](t, handler, "/v1/checkout/sessions/"+session.ID, url.Values{
		"metadata[env]":           {"ci"},
		"line_items[1][quantity]": {"3"},
	})
	if updated.Metadata["env"] != "ci" {
		t.Fatalf("updated metadata = %#v, want env=ci", updated.Metadata)
	}
	if updated.AmountTotal != 2500*4 {
		t.Fatalf("updated amount_total = %d, want %d", updated.AmountTotal, 2500*4)
	}

	// Non-open sessions reject updates.
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	if completion["session"] == nil {
		t.Fatalf("completion returned no session")
	}
	status, body = postFormStatus(t, handler, "/v1/checkout/sessions/"+session.ID, url.Values{
		"metadata[env]": {"done"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("completed session update status = %d body = %s, want 400", status, body)
	}
}

func TestInvoiceUpdateAndDelete(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"inv-update@example.test"}})
	invoice := postForm[struct {
		ID      string `json:"id"`
		Created int64  `json:"created"`
	}](t, handler, "/v1/invoices", url.Values{
		"customer": {customer.ID},
	})

	// Unknown params are rejected.
	status, body := postFormStatus(t, handler, "/v1/invoices/"+invoice.ID, url.Values{
		"auto_advance": {"false"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("unknown invoice update param status = %d body = %s, want 400", status, body)
	}

	updated := postForm[struct {
		Description string            `json:"description"`
		DueDate     int64             `json:"due_date"`
		Metadata    map[string]string `json:"metadata"`
	}](t, handler, "/v1/invoices/"+invoice.ID, url.Values{
		"description":    {"Consulting"},
		"days_until_due": {"14"},
		"metadata[env]":  {"ci"},
	})
	if updated.Description != "Consulting" {
		t.Fatalf("updated description = %q, want Consulting", updated.Description)
	}
	if updated.Metadata["env"] != "ci" {
		t.Fatalf("updated metadata = %#v, want env=ci", updated.Metadata)
	}
	if want := invoice.Created + 14*86400; updated.DueDate != want {
		t.Fatalf("updated due_date = %d, want %d (created %d + 14d)", updated.DueDate, want, invoice.Created)
	}

	// Attached lines are deleted with the draft invoice.
	item := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/invoiceitems", url.Values{
		"customer":    {customer.ID},
		"invoice":     {invoice.ID},
		"amount":      {"5000"},
		"currency":    {"usd"},
		"description": {"Initial"},
	})
	deleted := deleteForm[struct {
		ID      string `json:"id"`
		Deleted bool   `json:"deleted"`
	}](t, handler, "/v1/invoices/"+invoice.ID, url.Values{})
	if !deleted.Deleted || deleted.ID != invoice.ID {
		t.Fatalf("invoice delete = %#v, want deleted marker", deleted)
	}
	afterReq := httptest.NewRequest(http.MethodGet, "/v1/invoices/"+invoice.ID, nil)
	afterRec := httptest.NewRecorder()
	handler.ServeHTTP(afterRec, afterReq)
	if afterRec.Code != http.StatusNotFound {
		t.Fatalf("deleted invoice GET status = %d body = %s, want 404", afterRec.Code, afterRec.Body.String())
	}
	itemReq := httptest.NewRequest(http.MethodGet, "/v1/invoiceitems/"+item.ID, nil)
	itemRec := httptest.NewRecorder()
	handler.ServeHTTP(itemRec, itemReq)
	if itemRec.Code != http.StatusNotFound {
		t.Fatalf("attached item after invoice delete status = %d, want 404", itemRec.Code)
	}

	// Non-draft invoices cannot be deleted.
	other := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/invoices", url.Values{
		"customer": {customer.ID},
	})
	_ = postJSON[map[string]json.RawMessage](t, handler, "/v1/invoices/"+other.ID+"/finalize", map[string]string{})
	status, body = deleteFormStatus(t, handler, "/v1/invoices/"+other.ID, url.Values{})
	if status != http.StatusBadRequest {
		t.Fatalf("finalized invoice delete status = %d body = %s, want 400", status, body)
	}
}

func TestInvoiceLineMutation(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"inv-lines@example.test"}})
	invoice := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/invoices", url.Values{
		"customer": {customer.ID},
	})
	first := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/invoiceitems", url.Values{
		"customer":    {customer.ID},
		"invoice":     {invoice.ID},
		"amount":      {"5000"},
		"currency":    {"usd"},
		"description": {"First"},
	})

	type invoiceState struct {
		Subtotal  int64 `json:"subtotal"`
		Total     int64 `json:"total"`
		AmountDue int64 `json:"amount_due"`
	}

	// add_lines appends attached lines and recomputes totals.
	added := postForm[invoiceState](t, handler, "/v1/invoices/"+invoice.ID+"/add_lines", url.Values{
		"line_items[0][amount]":      {"2000"},
		"line_items[0][description]": {"Second"},
		"line_items[0][currency]":    {"usd"},
	})
	if added.Subtotal != 7000 || added.Total != 7000 || added.AmountDue != 7000 {
		t.Fatalf("after add_lines = %#v, want 7000", added)
	}
	// add_lines requires amounts.
	status, body := postFormStatus(t, handler, "/v1/invoices/"+invoice.ID+"/add_lines", url.Values{
		"line_items[0][description]": {"No amount"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("add_lines without amount status = %d body = %s, want 400", status, body)
	}

	// update_lines patches a line by id and recomputes totals.
	updated := postForm[invoiceState](t, handler, "/v1/invoices/"+invoice.ID+"/update_lines", url.Values{
		"line_items[0][id]":          {first.ID},
		"line_items[0][amount]":      {"3000"},
		"line_items[0][description]": {"First (revised)"},
	})
	if updated.Subtotal != 5000 || updated.Total != 5000 || updated.AmountDue != 5000 {
		t.Fatalf("after update_lines = %#v, want 5000", updated)
	}
	lines := getJSON[struct {
		Data []struct {
			ID          string `json:"id"`
			Amount      int64  `json:"amount"`
			Description string `json:"description"`
		} `json:"data"`
	}](t, handler, "/v1/invoices/"+invoice.ID+"/lines")
	if len(lines.Data) != 2 {
		t.Fatalf("invoice lines = %d, want 2", len(lines.Data))
	}
	for _, line := range lines.Data {
		if line.ID == first.ID && (line.Amount != 3000 || line.Description != "First (revised)") {
			t.Fatalf("updated line = %#v, want amount 3000 revised description", line)
		}
	}
	// update_lines without id is rejected.
	status, body = postFormStatus(t, handler, "/v1/invoices/"+invoice.ID+"/update_lines", url.Values{
		"line_items[0][amount]": {"1000"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("update_lines without id status = %d body = %s, want 400", status, body)
	}

	// remove_lines detaches a line and recomputes totals.
	removed := postForm[invoiceState](t, handler, "/v1/invoices/"+invoice.ID+"/remove_lines", url.Values{
		"line_items[0][id]": {first.ID},
	})
	if removed.Subtotal != 2000 || removed.Total != 2000 || removed.AmountDue != 2000 {
		t.Fatalf("after remove_lines = %#v, want 2000", removed)
	}
	// Line mutations require draft invoices.
	_ = postJSON[map[string]json.RawMessage](t, handler, "/v1/invoices/"+invoice.ID+"/finalize", map[string]string{})
	status, body = postFormStatus(t, handler, "/v1/invoices/"+invoice.ID+"/add_lines", url.Values{
		"line_items[0][amount]": {"1000"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("add_lines on finalized invoice status = %d body = %s, want 400", status, body)
	}
}

func TestProductDelete(t *testing.T) {
	handler := newTestHandler(t)
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Doomed Plan"}})

	deleted := deleteForm[struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}](t, handler, "/v1/products/"+product.ID, url.Values{})
	if !deleted.Deleted || deleted.Object != "product" || deleted.ID != product.ID {
		t.Fatalf("product delete = %#v, want deleted product marker", deleted)
	}
	afterReq := httptest.NewRequest(http.MethodGet, "/v1/products/"+product.ID, nil)
	afterRec := httptest.NewRecorder()
	handler.ServeHTTP(afterRec, afterReq)
	if afterRec.Code != http.StatusNotFound {
		t.Fatalf("deleted product GET status = %d body = %s, want 404", afterRec.Code, afterRec.Body.String())
	}
	status, body := deleteFormStatus(t, handler, "/v1/products/"+product.ID, url.Values{})
	if status != http.StatusNotFound {
		t.Fatalf("re-delete product status = %d body = %s, want 404", status, body)
	}
}
