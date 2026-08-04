package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/fixtures"
)

func TestTaxRatesEvidenceCRUD(t *testing.T) {
	handler := newTestHandler(t)

	// Missing required params → 400.
	status, body := postFormStatus(t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("missing percentage status = %d body = %s, want 400", status, body)
	}
	missing := decodeErrorBody(t, body)
	if missing.Error.Param != "percentage" {
		t.Fatalf("missing percentage error = %#v, want percentage param", missing.Error)
	}
	status, body = postFormStatus(t, handler, "/v1/tax_rates", url.Values{
		"percentage": {"10"},
		"inclusive":  {"false"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("missing display_name status = %d body = %s, want 400", status, body)
	}

	created := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"Sales Tax"},
		"percentage":   {"9.75"},
		"inclusive":    {"false"},
		"country":      {"US"},
		"state":        {"CA"},
		"description":  {"CA sales tax"},
		"metadata[k]":  {"v"},
	})
	id := fmt.Sprint(created["id"])
	if id == "" || created["object"] != "tax_rate" {
		t.Fatalf("created tax_rate = %#v, want tax_rate object with id", created)
	}
	if created["display_name"] != "Sales Tax" || created["percentage"] != 9.75 || created["inclusive"] != false {
		t.Fatalf("created tax_rate fields = %#v", created)
	}
	if created["active"] != true {
		t.Fatalf("created active = %#v, want true default", created["active"])
	}
	if created["country"] != "US" || created["state"] != "CA" {
		t.Fatalf("created jurisdiction fields = %#v", created)
	}
	if created["rate_type"] != "percentage" || created["flat_amount"] != nil || created["effective_percentage"] != nil || created["jurisdiction_level"] != nil || created["tax_type"] != nil {
		t.Fatalf("created tax_rate v22 fields = %#v, want rate_type=percentage and null flat/effective/jurisdiction/tax_type", created)
	}

	fetched := getJSON[map[string]any](t, handler, "/v1/tax_rates/"+id)
	if fetched["id"] != id || fetched["display_name"] != "Sales Tax" {
		t.Fatalf("GET tax_rate = %#v, want id %s", fetched, id)
	}

	listed := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/tax_rates")
	if len(listed.Data) != 1 || listed.Data[0]["id"] != id {
		t.Fatalf("list tax_rates = %#v, want one item %s", listed.Data, id)
	}

	updated := postForm[map[string]any](t, handler, "/v1/tax_rates/"+id, url.Values{
		"active":       {"false"},
		"display_name": {"Updated Tax"},
		"description":  {"updated"},
		"metadata[k]":  {"v2"},
	})
	if updated["active"] != false || updated["display_name"] != "Updated Tax" || updated["description"] != "updated" {
		t.Fatalf("updated tax_rate = %#v", updated)
	}
	meta, _ := updated["metadata"].(map[string]any)
	if fmt.Sprint(meta["k"]) != "v2" {
		t.Fatalf("updated metadata = %#v, want k=v2", updated["metadata"])
	}
	// percentage is not updatable and must remain.
	if updated["percentage"] != 9.75 {
		t.Fatalf("updated percentage = %#v, want 9.75", updated["percentage"])
	}
}

func TestCustomerTaxIDsEvidenceCRUD(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"tax-id@example.test"},
	})

	// Missing customer → 404.
	status, body := postFormStatus(t, handler, "/v1/customers/cus_missing/tax_ids", url.Values{
		"type":  {"eu_vat"},
		"value": {"DE123"},
	})
	if status != http.StatusNotFound {
		t.Fatalf("missing customer tax_id status = %d body = %s, want 404", status, body)
	}

	// Missing required params → 400.
	status, body = postFormStatus(t, handler, "/v1/customers/"+customer.ID+"/tax_ids", url.Values{
		"type": {"eu_vat"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("missing value status = %d body = %s, want 400", status, body)
	}

	created := postForm[map[string]any](t, handler, "/v1/customers/"+customer.ID+"/tax_ids", url.Values{
		"type":  {"eu_vat"},
		"value": {"DE123456789"},
	})
	id := fmt.Sprint(created["id"])
	if id == "" || created["object"] != "tax_id" || created["customer"] != customer.ID {
		t.Fatalf("created tax_id = %#v", created)
	}
	if created["type"] != "eu_vat" || created["value"] != "DE123456789" {
		t.Fatalf("created tax_id fields = %#v", created)
	}
	if created["customer_account"] != nil {
		t.Fatalf("created tax_id customer_account = %#v, want null", created["customer_account"])
	}
	owner, ok := created["owner"].(map[string]any)
	if !ok || owner["type"] != "customer" || owner["customer"] != customer.ID || owner["customer_account"] != nil {
		t.Fatalf("created tax_id owner = %#v, want type=customer and customer id", created["owner"])
	}
	verification, ok := created["verification"].(map[string]any)
	if !ok || verification["status"] != "verified" {
		t.Fatalf("verification = %#v, want verified", created["verification"])
	}

	fetched := getJSON[map[string]any](t, handler, "/v1/customers/"+customer.ID+"/tax_ids/"+id)
	if fetched["id"] != id {
		t.Fatalf("GET tax_id = %#v, want %s", fetched, id)
	}

	listed := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/customers/"+customer.ID+"/tax_ids")
	if len(listed.Data) != 1 || listed.Data[0]["id"] != id {
		t.Fatalf("list tax_ids = %#v, want one item %s", listed.Data, id)
	}

	deleted := deleteJSON[map[string]any](t, handler, "/v1/customers/"+customer.ID+"/tax_ids/"+id)
	if deleted["id"] != id || deleted["object"] != "tax_id" || deleted["deleted"] != true {
		t.Fatalf("deleted tax_id = %#v", deleted)
	}
	status, body = getStatus(t, handler, "/v1/customers/"+customer.ID+"/tax_ids/"+id)
	if status != http.StatusNotFound {
		t.Fatalf("GET deleted tax_id status = %d body = %s, want 404", status, body)
	}
}

func TestAutomaticTaxCheckoutInvoiceAndRenewal(t *testing.T) {
	handler := newTestHandler(t)

	// Price 7000, amount-off coupon 1000 → discounted 6000, tax 10% = 600, total 6600.
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":                 {"taxed@example.test"},
		"metadata[tax_percent]": {"10"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Taxed Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"7000"},
		"recurring[interval]": {"month"},
	})
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"id":         {"coupon_tax_off_1000"},
		"amount_off": {"1000"},
		"currency":   {"usd"},
		"duration":   {"forever"},
	})

	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                   {customer.ID},
		"line_items[0][price]":       {price.ID},
		"line_items[0][quantity]":    {"1"},
		"discounts[0][coupon]":       {coupon.ID},
		"automatic_tax[enabled]":     {"true"},
		"tax_id_collection[enabled]": {"true"},
		"success_url":                {"http://app.test/success"},
		"cancel_url":                 {"http://app.test/cancel"},
	})
	if session["amount_subtotal"] != float64(7000) {
		t.Fatalf("amount_subtotal = %#v, want 7000", session["amount_subtotal"])
	}
	if session["amount_total"] != float64(6600) {
		t.Fatalf("amount_total = %#v, want 6600 (6000 + 10%% tax)", session["amount_total"])
	}
	totalDetails, ok := session["total_details"].(map[string]any)
	if !ok || totalDetails["amount_discount"] != float64(1000) || totalDetails["amount_tax"] != float64(600) {
		t.Fatalf("total_details = %#v, want discount=1000 tax=600", session["total_details"])
	}
	autoTax, ok := session["automatic_tax"].(map[string]any)
	if !ok || autoTax["enabled"] != true || autoTax["status"] != "complete" || autoTax["provider"] != "stripe" || autoTax["liability"] != nil {
		t.Fatalf("session automatic_tax = %#v, want enabled/status/provider/liability v22 shape", session["automatic_tax"])
	}
	if _, hasDisabled := autoTax["disabled_reason"]; hasDisabled {
		t.Fatalf("session automatic_tax must not include disabled_reason: %#v", autoTax)
	}
	taxIDCollection, ok := session["tax_id_collection"].(map[string]any)
	if !ok || taxIDCollection["enabled"] != true || taxIDCollection["required"] != "never" {
		t.Fatalf("session tax_id_collection = %#v, want enabled true required never", session["tax_id_collection"])
	}

	sessionID := fmt.Sprint(session["id"])
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+sessionID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var completed billing.CheckoutSession
	if err := json.Unmarshal(completion["session"], &completed); err != nil {
		t.Fatalf("decode completed session: %v", err)
	}

	invoice := getJSON[map[string]any](t, handler, "/v1/invoices/"+completed.InvoiceID)
	if invoice["subtotal"] != float64(7000) || invoice["total"] != float64(6600) || invoice["amount_paid"] != float64(6600) {
		t.Fatalf("invoice amounts = %#v, want subtotal=7000 total/paid=6600", invoice)
	}
	if invoice["tax"] != float64(600) {
		t.Fatalf("invoice tax = %#v, want 600", invoice["tax"])
	}
	if invoice["total_excluding_tax"] != float64(6000) {
		t.Fatalf("total_excluding_tax = %#v, want 6000", invoice["total_excluding_tax"])
	}
	invAutoTax, ok := invoice["automatic_tax"].(map[string]any)
	if !ok || invAutoTax["enabled"] != true || invAutoTax["status"] != "complete" || invAutoTax["provider"] != "stripe" || invAutoTax["liability"] != nil || invAutoTax["disabled_reason"] != nil {
		t.Fatalf("invoice automatic_tax = %#v, want v22 invoice shape", invoice["automatic_tax"])
	}
	taxAmounts, ok := invoice["total_tax_amounts"].([]any)
	if !ok || len(taxAmounts) != 1 {
		t.Fatalf("total_tax_amounts = %#v, want one entry", invoice["total_tax_amounts"])
	}
	taxEntry, _ := taxAmounts[0].(map[string]any)
	if taxEntry["amount"] != float64(600) || taxEntry["inclusive"] != false || taxEntry["tax_rate"] != "txr_billtap_simulated" {
		t.Fatalf("tax amount entry = %#v", taxEntry)
	}
	totalTaxes, ok := invoice["total_taxes"].([]any)
	if !ok || len(totalTaxes) != 1 {
		t.Fatalf("total_taxes = %#v, want one entry", invoice["total_taxes"])
	}
	totalTax, _ := totalTaxes[0].(map[string]any)
	if totalTax["amount"] != float64(600) || totalTax["tax_behavior"] != "exclusive" || totalTax["taxability_reason"] != "standard_rated" || totalTax["taxable_amount"] != float64(6000) || totalTax["type"] != "tax_rate_details" {
		t.Fatalf("total_taxes[0] = %#v, want exclusive 600 on taxable 6000", totalTax)
	}
	rateDetails, _ := totalTax["tax_rate_details"].(map[string]any)
	if rateDetails["tax_rate"] != "txr_billtap_simulated" {
		t.Fatalf("total_taxes tax_rate_details = %#v", totalTax["tax_rate_details"])
	}

	subscription := getJSON[struct {
		ID       string            `json:"id"`
		Metadata map[string]string `json:"metadata"`
		AutoTax  map[string]any    `json:"automatic_tax"`
	}](t, handler, "/v1/subscriptions/"+completed.SubscriptionID)
	if subscription.Metadata[billing.MetadataAutomaticTax] != "true" || subscription.Metadata[billing.MetadataTaxPercent] != "10" {
		t.Fatalf("subscription tax metadata = %#v", subscription.Metadata)
	}
	if subscription.AutoTax["enabled"] != true || subscription.AutoTax["disabled_reason"] != nil || subscription.AutoTax["liability"] != nil {
		t.Fatalf("subscription automatic_tax = %#v, want enabled with null liability/disabled_reason", subscription.AutoTax)
	}
	if _, hasStatus := subscription.AutoTax["status"]; hasStatus {
		t.Fatalf("subscription automatic_tax must not include status: %#v", subscription.AutoTax)
	}
	if _, hasProvider := subscription.AutoTax["provider"]; hasProvider {
		t.Fatalf("subscription automatic_tax must not include provider: %#v", subscription.AutoTax)
	}

	// Tax-id collection echo already covered on session create above.
	// automatic_tax without tax_percent → tax=0, status complete, total unchanged.
	noTaxCustomer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"no-tax-meta@example.test"},
	})
	noTaxSession := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {noTaxCustomer.ID},
		"line_items[0][price]":    {price.ID},
		"line_items[0][quantity]": {"1"},
		"automatic_tax[enabled]":  {"true"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	if noTaxSession["amount_subtotal"] != float64(7000) || noTaxSession["amount_total"] != float64(7000) {
		t.Fatalf("no tax_percent session amounts = %#v, want total=7000", noTaxSession)
	}
	noTaxDetails, _ := noTaxSession["total_details"].(map[string]any)
	if noTaxDetails["amount_tax"] != float64(0) {
		t.Fatalf("no tax_percent amount_tax = %#v, want 0", noTaxDetails["amount_tax"])
	}
	noTaxAuto, _ := noTaxSession["automatic_tax"].(map[string]any)
	if noTaxAuto["enabled"] != true || noTaxAuto["status"] != "complete" {
		t.Fatalf("no tax_percent automatic_tax = %#v", noTaxSession["automatic_tax"])
	}

	// Renewal via test clock uses metadata snapshot (not live customer tax_percent).
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_tax"},
		"frozen_time": {strconv.FormatInt(time.Date(2032, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	clockCustomer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":                 {"tax-renewal@example.test"},
		"test_clock":            {"clock_tax"},
		"metadata[tax_percent]": {"10"},
	})
	sub := postForm[struct {
		ID            string            `json:"id"`
		LatestInvoice string            `json:"latest_invoice"`
		Metadata      map[string]string `json:"metadata"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":               {clockCustomer.ID},
		"items[0][price]":        {price.ID},
		"items[0][quantity]":     {"1"},
		"discounts[0][coupon]":   {coupon.ID},
		"automatic_tax[enabled]": {"true"},
		"test_clock":             {"clock_tax"},
	})
	if sub.Metadata[billing.MetadataAutomaticTax] != "true" || sub.Metadata[billing.MetadataTaxPercent] != "10" {
		t.Fatalf("direct subscription tax metadata = %#v", sub.Metadata)
	}
	firstInvoice := getJSON[map[string]any](t, handler, "/v1/invoices/"+sub.LatestInvoice)
	if firstInvoice["total"] != float64(6600) || firstInvoice["tax"] != float64(600) {
		t.Fatalf("direct subscription first invoice = %#v, want total=6600 tax=600", firstInvoice)
	}

	// Change customer tax_percent after snapshot — renewal must keep 10%.
	_ = postForm[billing.Customer](t, handler, "/v1/customers/"+clockCustomer.ID, url.Values{
		"metadata[tax_percent]": {"25"},
	})
	advance := postForm[struct {
		BilltapAdvanceResult struct {
			Renewals []struct {
				Invoice struct {
					Subtotal int64 `json:"subtotal"`
					Total    int64 `json:"total"`
					Tax      int64 `json:"tax"`
				} `json:"invoice"`
			} `json:"renewals"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_tax/advance", url.Values{
		"frozen_time": {strconv.FormatInt(time.Date(2032, 2, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	if len(advance.BilltapAdvanceResult.Renewals) != 1 {
		t.Fatalf("tax renewal = %#v, want one renewal", advance.BilltapAdvanceResult.Renewals)
	}
	renewal := advance.BilltapAdvanceResult.Renewals[0].Invoice
	if renewal.Subtotal != 7000 || renewal.Tax != 600 || renewal.Total != 6600 {
		t.Fatalf("tax renewal invoice = %#v, want subtotal=7000 tax=600 total=6600 (snapshot 10%%, not 25%%)", renewal)
	}
}

func TestCheckoutWithoutAutomaticTaxKeepsZeroTax(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":                 {"plain@example.test"},
		"metadata[tax_percent]": {"10"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Plain"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"5000"},
		"recurring[interval]": {"month"},
	})
	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {price.ID},
		"line_items[0][quantity]": {"1"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	if session["amount_total"] != float64(5000) {
		t.Fatalf("amount_total = %#v, want 5000 without automatic_tax", session["amount_total"])
	}
	details, _ := session["total_details"].(map[string]any)
	if details["amount_tax"] != float64(0) {
		t.Fatalf("amount_tax = %#v, want 0 without automatic_tax", details["amount_tax"])
	}
	autoTax, ok := session["automatic_tax"].(map[string]any)
	if !ok || autoTax["enabled"] != false || autoTax["status"] != nil || autoTax["provider"] != nil || autoTax["liability"] != nil {
		t.Fatalf("automatic_tax = %#v, want disabled v22 shape", session["automatic_tax"])
	}
	taxIDCollection, ok := session["tax_id_collection"].(map[string]any)
	if !ok || taxIDCollection["enabled"] != false || taxIDCollection["required"] != "never" {
		t.Fatalf("tax_id_collection = %#v, want disabled required never", session["tax_id_collection"])
	}

	sessionID := fmt.Sprint(session["id"])
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+sessionID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var completed billing.CheckoutSession
	if err := json.Unmarshal(completion["session"], &completed); err != nil {
		t.Fatalf("decode completed session: %v", err)
	}
	invoice := getJSON[map[string]any](t, handler, "/v1/invoices/"+completed.InvoiceID)
	if invoice["total"] != float64(5000) || invoice["tax"] != nil {
		t.Fatalf("plain invoice = %#v, want total=5000 tax=null", invoice)
	}
	invAuto, _ := invoice["automatic_tax"].(map[string]any)
	if invAuto["enabled"] != false || invAuto["status"] != nil || invAuto["provider"] != nil || invAuto["disabled_reason"] != nil {
		t.Fatalf("plain invoice automatic_tax = %#v, want disabled v22 shape", invoice["automatic_tax"])
	}
	taxAmounts, _ := invoice["total_tax_amounts"].([]any)
	if len(taxAmounts) != 0 {
		t.Fatalf("plain total_tax_amounts = %#v, want empty", invoice["total_tax_amounts"])
	}
	totalTaxes, _ := invoice["total_taxes"].([]any)
	if len(totalTaxes) != 0 {
		t.Fatalf("plain total_taxes = %#v, want empty", invoice["total_taxes"])
	}
}

func TestDefaultTaxRatesCheckoutInvoiceAndSubscription(t *testing.T) {
	handler := newTestHandler(t)

	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"default-tax@example.test"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Taxed Default"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"5000"},
		"recurring[interval]": {"month"},
	})
	taxRate := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	txrID := fmt.Sprint(taxRate["id"])

	// Single-value [] form (firstValues leaves key as subscription_data[default_tax_rates][]).
	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {customer.ID},
		"line_items[0][price]":                   {price.ID},
		"line_items[0][quantity]":                {"1"},
		"subscription_data[default_tax_rates][]": {txrID},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	if session["amount_total"] != float64(5500) {
		t.Fatalf("amount_total = %#v, want 5500", session["amount_total"])
	}
	details, _ := session["total_details"].(map[string]any)
	if details["amount_tax"] != float64(500) {
		t.Fatalf("amount_tax = %#v, want 500", details["amount_tax"])
	}

	sessionID := fmt.Sprint(session["id"])
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+sessionID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var completed billing.CheckoutSession
	if err := json.Unmarshal(completion["session"], &completed); err != nil {
		t.Fatalf("decode completed session: %v", err)
	}

	invoice := getJSON[map[string]any](t, handler, "/v1/invoices/"+completed.InvoiceID)
	if invoice["tax"] != float64(500) || invoice["total"] != float64(5500) {
		t.Fatalf("invoice tax/total = %#v, want 500/5500", invoice)
	}
	defaultRates, ok := invoice["default_tax_rates"].([]any)
	if !ok || len(defaultRates) != 1 {
		t.Fatalf("invoice default_tax_rates = %#v", invoice["default_tax_rates"])
	}
	rate0, _ := defaultRates[0].(map[string]any)
	if rate0["id"] != txrID {
		t.Fatalf("default_tax_rates[0].id = %#v, want %s", rate0["id"], txrID)
	}
	totalTaxes, ok := invoice["total_taxes"].([]any)
	if !ok || len(totalTaxes) != 1 {
		t.Fatalf("total_taxes = %#v", invoice["total_taxes"])
	}
	totalTax, _ := totalTaxes[0].(map[string]any)
	rateDetails, _ := totalTax["tax_rate_details"].(map[string]any)
	if rateDetails["tax_rate"] != txrID {
		t.Fatalf("total_taxes tax_rate = %#v, want %s", rateDetails["tax_rate"], txrID)
	}
	if totalTax["taxability_reason"] != nil {
		t.Fatalf("manual tax taxability_reason = %#v, want nil", totalTax["taxability_reason"])
	}

	subscription := getJSON[map[string]any](t, handler, "/v1/subscriptions/"+completed.SubscriptionID)
	subRates, ok := subscription["default_tax_rates"].([]any)
	if !ok || len(subRates) != 1 {
		t.Fatalf("subscription default_tax_rates = %#v", subscription["default_tax_rates"])
	}
	subRate, _ := subRates[0].(map[string]any)
	if subRate["percentage"] != float64(10) || subRate["id"] != txrID {
		t.Fatalf("subscription default_tax_rates[0] = %#v", subRate)
	}
}

func TestDefaultTaxRatesInclusiveAndMixed(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"incl-tax@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Incl Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1100"},
		"recurring[interval]": {"month"},
	})
	incl := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT incl"},
		"percentage":   {"10"},
		"inclusive":    {"true"},
	})
	excl := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"Local excl"},
		"percentage":   {"5"},
		"inclusive":    {"false"},
	})

	// Inclusive alone: base 1100 → tax 100, amount_total 1100, pretax 1000.
	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {customer.ID},
		"line_items[0][price]":                   {price.ID},
		"line_items[0][quantity]":                {"1"},
		"subscription_data[default_tax_rates][]": {fmt.Sprint(incl["id"])},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	if session["amount_total"] != float64(1100) {
		t.Fatalf("inclusive amount_total = %#v, want 1100", session["amount_total"])
	}
	details, _ := session["total_details"].(map[string]any)
	if details["amount_tax"] != float64(100) {
		t.Fatalf("inclusive amount_tax = %#v, want 100", details["amount_tax"])
	}
	sessionID := fmt.Sprint(session["id"])
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+sessionID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var completed billing.CheckoutSession
	if err := json.Unmarshal(completion["session"], &completed); err != nil {
		t.Fatalf("decode: %v", err)
	}
	invoice := getJSON[map[string]any](t, handler, "/v1/invoices/"+completed.InvoiceID)
	if invoice["total"] != float64(1100) || invoice["tax"] != float64(100) || invoice["total_excluding_tax"] != float64(1000) {
		t.Fatalf("inclusive invoice = %#v, want total=1100 tax=100 pretax=1000", invoice)
	}

	// Mixed via multi-value [] form (firstValues expands to [0]/[1]).
	mixedCustomer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"mixed-tax@example.test"}})
	mixedSession := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {mixedCustomer.ID},
		"line_items[0][price]":                   {price.ID},
		"line_items[0][quantity]":                {"1"},
		"subscription_data[default_tax_rates][]": {fmt.Sprint(incl["id"]), fmt.Sprint(excl["id"])},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	// base 1100 → incl 100, pretax 1000, excl 50, total 1150, taxTotal 150
	if mixedSession["amount_total"] != float64(1150) {
		t.Fatalf("mixed amount_total = %#v, want 1150", mixedSession["amount_total"])
	}
	mixedDetails, _ := mixedSession["total_details"].(map[string]any)
	if mixedDetails["amount_tax"] != float64(150) {
		t.Fatalf("mixed amount_tax = %#v, want 150", mixedDetails["amount_tax"])
	}
	mixedID := fmt.Sprint(mixedSession["id"])
	mixedCompletion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+mixedID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var mixedCompleted billing.CheckoutSession
	if err := json.Unmarshal(mixedCompletion["session"], &mixedCompleted); err != nil {
		t.Fatalf("decode mixed: %v", err)
	}
	mixedInvoice := getJSON[map[string]any](t, handler, "/v1/invoices/"+mixedCompleted.InvoiceID)
	if mixedInvoice["total"] != float64(1150) || mixedInvoice["tax"] != float64(150) || mixedInvoice["total_excluding_tax"] != float64(1000) {
		t.Fatalf("mixed invoice = %#v, want total=1150 tax=150 pretax=1000", mixedInvoice)
	}
	mixedTaxes, _ := mixedInvoice["total_taxes"].([]any)
	if len(mixedTaxes) != 2 {
		t.Fatalf("mixed total_taxes = %#v, want 2 entries", mixedInvoice["total_taxes"])
	}
}

func TestDefaultTaxRatesWithCoupon(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"coupon-tax@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Coupon Tax"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"8000"},
		"recurring[interval]": {"month"},
	})
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_tax_25"},
		"percent_off": {"25"},
		"duration":    {"forever"},
	})
	taxRate := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	// subtotal 8000, 25% off → 6000, 10% exclusive → tax 600, total 6600
	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {customer.ID},
		"line_items[0][price]":                   {price.ID},
		"line_items[0][quantity]":                {"1"},
		"discounts[0][coupon]":                   {coupon.ID},
		"subscription_data[default_tax_rates][]": {fmt.Sprint(taxRate["id"])},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	if session["amount_total"] != float64(6600) {
		t.Fatalf("coupon+tax amount_total = %#v, want 6600", session["amount_total"])
	}
	details, _ := session["total_details"].(map[string]any)
	if details["amount_discount"] != float64(2000) || details["amount_tax"] != float64(600) {
		t.Fatalf("coupon+tax total_details = %#v, want discount=2000 tax=600", details)
	}
}

func TestDefaultTaxRatesRenewal(t *testing.T) {
	handler := newTestHandler(t)
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_default_tax"},
		"frozen_time": {strconv.FormatInt(time.Date(2033, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":      {"renew-tax@example.test"},
		"test_clock": {"clock_default_tax"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Renew Tax"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"5000"},
		"recurring[interval]": {"month"},
	})
	taxRate := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	sub := postForm[struct {
		ID            string `json:"id"`
		LatestInvoice string `json:"latest_invoice"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":            {customer.ID},
		"items[0][price]":     {price.ID},
		"items[0][quantity]":  {"1"},
		"default_tax_rates[]": {fmt.Sprint(taxRate["id"])},
		"test_clock":          {"clock_default_tax"},
	})
	first := getJSON[map[string]any](t, handler, "/v1/invoices/"+sub.LatestInvoice)
	if first["total"] != float64(5500) || first["tax"] != float64(500) {
		t.Fatalf("first invoice = %#v, want total=5500 tax=500", first)
	}
	advance := postForm[struct {
		BilltapAdvanceResult struct {
			Renewals []struct {
				Invoice struct {
					Subtotal int64 `json:"subtotal"`
					Total    int64 `json:"total"`
					Tax      int64 `json:"tax"`
				} `json:"invoice"`
			} `json:"renewals"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_default_tax/advance", url.Values{
		"frozen_time": {strconv.FormatInt(time.Date(2033, 2, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	if len(advance.BilltapAdvanceResult.Renewals) != 1 {
		t.Fatalf("renewals = %#v, want 1", advance.BilltapAdvanceResult.Renewals)
	}
	renewal := advance.BilltapAdvanceResult.Renewals[0].Invoice
	if renewal.Subtotal != 5000 || renewal.Tax != 500 || renewal.Total != 5500 {
		t.Fatalf("renewal invoice = %#v, want subtotal=5000 tax=500 total=5500", renewal)
	}
}

func TestDefaultTaxRatesSubscriptionUpdateAndClear(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"update-tax@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Update Tax"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"4000"},
		"recurring[interval]": {"month"},
	})
	rateA := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"A"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	rateB := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"B"},
		"percentage":   {"20"},
		"inclusive":    {"false"},
	})
	sub := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":            {customer.ID},
		"items[0][price]":     {price.ID},
		"items[0][quantity]":  {"1"},
		"default_tax_rates[]": {fmt.Sprint(rateA["id"])},
	})
	updated := postForm[map[string]any](t, handler, "/v1/subscriptions/"+sub.ID, url.Values{
		"default_tax_rates[]": {fmt.Sprint(rateB["id"])},
	})
	rates, _ := updated["default_tax_rates"].([]any)
	if len(rates) != 1 {
		t.Fatalf("updated rates = %#v, want one", updated["default_tax_rates"])
	}
	rate0, _ := rates[0].(map[string]any)
	if rate0["id"] != rateB["id"] || rate0["percentage"] != float64(20) {
		t.Fatalf("updated rate = %#v, want B 20%%", rate0)
	}
	// Emptyable clear with single empty string.
	cleared := postForm[map[string]any](t, handler, "/v1/subscriptions/"+sub.ID, url.Values{
		"default_tax_rates": {""},
	})
	clearedRates, _ := cleared["default_tax_rates"].([]any)
	if len(clearedRates) != 0 {
		t.Fatalf("cleared default_tax_rates = %#v, want []", cleared["default_tax_rates"])
	}
}

func TestDefaultTaxRatesErrors(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":                 {"err-tax@example.test"},
		"metadata[tax_percent]": {"10"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Err Tax"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"3000"},
		"recurring[interval]": {"month"},
	})

	// Missing tax rate → 400 resource_missing.
	status, body := postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {customer.ID},
		"line_items[0][price]":                   {price.ID},
		"line_items[0][quantity]":                {"1"},
		"subscription_data[default_tax_rates][]": {"txr_missing"},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("missing tax rate status = %d body = %s, want 400", status, body)
	}
	errBody := decodeErrorBody(t, body)
	if errBody.Error.Code != "resource_missing" {
		t.Fatalf("missing tax rate error = %#v, want resource_missing", errBody.Error)
	}
	if errBody.Error.Message != "No such tax rate: 'txr_missing'" {
		t.Fatalf("missing tax rate message = %q", errBody.Error.Message)
	}

	// automatic_tax + default_tax_rates → 400 parameter_invalid.
	taxRate := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	status, body = postFormStatus(t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {customer.ID},
		"line_items[0][price]":                   {price.ID},
		"line_items[0][quantity]":                {"1"},
		"automatic_tax[enabled]":                 {"true"},
		"subscription_data[default_tax_rates][]": {fmt.Sprint(taxRate["id"])},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("conflict status = %d body = %s, want 400", status, body)
	}
	conflict := decodeErrorBody(t, body)
	if conflict.Error.Code != "parameter_invalid" {
		t.Fatalf("conflict error = %#v, want parameter_invalid", conflict.Error)
	}
	if conflict.Error.Message != "You cannot specify both automatic_tax[enabled]=true and default_tax_rates." {
		t.Fatalf("conflict message = %q", conflict.Error.Message)
	}

	// Indexed form [0] also accepted (same as multi-value expansion).
	okSession := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                                {customer.ID},
		"line_items[0][price]":                    {price.ID},
		"line_items[0][quantity]":                 {"1"},
		"subscription_data[default_tax_rates][0]": {fmt.Sprint(taxRate["id"])},
		"success_url":                             {"http://app.test/success"},
		"cancel_url":                              {"http://app.test/cancel"},
	})
	if okSession["amount_total"] != float64(3300) {
		t.Fatalf("indexed form amount_total = %#v, want 3300", okSession["amount_total"])
	}
}

func TestFixtureTaxRatesApplyGetAndValidate(t *testing.T) {
	handler := newTestHandler(t)

	pack := map[string]any{
		"name": "tax-rate-fixture",
		"tax_rates": []map[string]any{
			{
				"id":           "txr_fixture_ex",
				"display_name": "Sales Tax",
				"percentage":   10.0,
				"inclusive":    false,
				"country":      "US",
				"state":        "CA",
			},
			{
				// Auto-generated txr_ id; inclusive for field round-trip.
				"display_name": "VAT Inclusive",
				"percentage":   20.0,
				"inclusive":    true,
			},
		},
	}

	// validate endpoint summary includes tax_rates count.
	validation := postJSON[struct {
		Valid   bool           `json:"valid"`
		Summary map[string]int `json:"summary"`
	}](t, handler, "/api/fixtures/validate", pack)
	if !validation.Valid || validation.Summary["tax_rates"] != 2 {
		t.Fatalf("fixture validation = %#v, want valid with tax_rates=2", validation)
	}

	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["tax_rates"] != 2 || len(applied.TaxRates) != 2 {
		t.Fatalf("fixture apply summary/tax_rates = summary=%#v tax_rates=%#v, want 2", applied.Summary, applied.TaxRates)
	}

	// Explicit ID round-trip.
	explicit := getJSON[map[string]any](t, handler, "/v1/tax_rates/txr_fixture_ex")
	if explicit["id"] != "txr_fixture_ex" || explicit["percentage"] != float64(10) || explicit["inclusive"] != false {
		t.Fatalf("GET explicit tax_rate = %#v, want id/percentage/inclusive round-trip", explicit)
	}
	if explicit["display_name"] != "Sales Tax" || explicit["active"] != true {
		t.Fatalf("GET explicit tax_rate fields = %#v", explicit)
	}

	// Auto ID: one of the applied rates is not the explicit id and has inclusive=true.
	var autoID string
	for _, rate := range applied.TaxRates {
		id := fmt.Sprint(rate["id"])
		if id != "txr_fixture_ex" {
			autoID = id
			if rate["inclusive"] != true || rate["percentage"] != float64(20) {
				t.Fatalf("auto tax_rate payload = %#v, want inclusive 20%%", rate)
			}
			if !strings.HasPrefix(id, "txr_") {
				t.Fatalf("auto tax_rate id = %q, want txr_ prefix", id)
			}
			break
		}
	}
	if autoID == "" {
		t.Fatalf("applied tax_rates = %#v, want one auto-generated id", applied.TaxRates)
	}
	autoFetched := getJSON[map[string]any](t, handler, "/v1/tax_rates/"+autoID)
	if autoFetched["id"] != autoID || autoFetched["inclusive"] != true || autoFetched["percentage"] != float64(20) {
		t.Fatalf("GET auto tax_rate = %#v", autoFetched)
	}

	// Re-apply same pack: idempotent overwrite for explicit ID, no error.
	reapplied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if reapplied.Summary["tax_rates"] != 2 {
		t.Fatalf("re-apply summary tax_rates = %#v, want 2", reapplied.Summary)
	}
	again := getJSON[map[string]any](t, handler, "/v1/tax_rates/txr_fixture_ex")
	if again["percentage"] != float64(10) || again["display_name"] != "Sales Tax" {
		t.Fatalf("re-apply GET = %#v, want stable explicit tax rate", again)
	}
}

func TestFixtureTaxRatesCheckoutDefaultTaxRates(t *testing.T) {
	handler := newTestHandler(t)

	// Seed exclusive 10% tax rate via fixture pack.
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "tax-rate-checkout",
		"tax_rates": []map[string]any{{
			"id":           "txr_fixture_checkout",
			"display_name": "Sales",
			"percentage":   10.0,
			"inclusive":    false,
		}},
		"products": []map[string]any{{
			"id":   "prod_fixture_tax",
			"name": "Taxed Plan",
		}},
		"prices": []map[string]any{{
			"id":          "price_fixture_tax",
			"product":     "prod_fixture_tax",
			"currency":    "usd",
			"unit_amount": 1000,
			"interval":    "month",
		}},
		"customers": []map[string]any{{
			"id":    "cus_fixture_tax",
			"email": "fixture-tax@example.test",
		}},
	})
	if applied.Summary["tax_rates"] != 1 {
		t.Fatalf("apply summary = %#v, want tax_rates=1", applied.Summary)
	}

	// 10% exclusive on subtotal 1000 → tax 100.
	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {"cus_fixture_tax"},
		"line_items[0][price]":                   {"price_fixture_tax"},
		"line_items[0][quantity]":                {"1"},
		"subscription_data[default_tax_rates][]": {"txr_fixture_checkout"},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	details, _ := session["total_details"].(map[string]any)
	if details["amount_tax"] != float64(100) {
		t.Fatalf("amount_tax = %#v, want 100 (10%% of 1000)", details["amount_tax"])
	}
	if session["amount_total"] != float64(1100) {
		t.Fatalf("amount_total = %#v, want 1100", session["amount_total"])
	}
}
