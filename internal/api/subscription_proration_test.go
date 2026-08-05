package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/fixtures"
)

// subscription proration helpers — amounts are integer cents; VAT 10% exclusive.

type prorationSubResponse struct {
	ID                 string            `json:"id"`
	LatestInvoice      string            `json:"latest_invoice"`
	CurrentPeriodStart int64             `json:"current_period_start"`
	CurrentPeriodEnd   int64             `json:"current_period_end"`
	Metadata           map[string]string `json:"metadata"`
	Items              struct {
		Data []struct {
			ID    string `json:"id"`
			Price struct {
				ID string `json:"id"`
			} `json:"price"`
		} `json:"data"`
	} `json:"items"`
}

type prorationInvoiceResponse struct {
	ID            string `json:"id"`
	Subtotal      int64  `json:"subtotal"`
	Total         int64  `json:"total"`
	Tax           int64  `json:"tax"`
	AmountPaid    int64  `json:"amount_paid"`
	AmountDue     int64  `json:"amount_due"`
	BillingReason string `json:"billing_reason"`
	Status        string `json:"status"`
}

func setupProrationPlans(t *testing.T, handler http.Handler) (customer billing.Customer, lite, pro billing.Price, taxRateID string) {
	t.Helper()
	customer = postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"proration@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Proration Plan"}})
	lite = postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"4900"},
		"recurring[interval]": {"month"},
	})
	pro = postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"9900"},
		"recurring[interval]": {"month"},
	})
	txr := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	return customer, lite, pro, txr.ID
}

func listSubscriptionInvoices(t *testing.T, handler http.Handler, subID string) []prorationInvoiceResponse {
	t.Helper()
	listed := getJSON[struct {
		Data []prorationInvoiceResponse `json:"data"`
	}](t, handler, "/v1/invoices?subscription="+subID)
	return listed.Data
}

// TestSubscriptionUpdateAlwaysInvoiceUpgrade: 4900→9900 at period start →
// delta 5000, tax 500, total 5500, billing_reason=subscription_update.
func TestSubscriptionUpdateAlwaysInvoiceUpgrade(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})
	periodStart := created.CurrentPeriodStart
	periodEnd := created.CurrentPeriodEnd
	createInvoiceID := created.LatestInvoice

	// Pin proration_date to period start so remaining == full period (no clock skew).
	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {created.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"proration_date":       {strconv.FormatInt(periodStart, 10)},
		"payment_behavior":     {"error_if_incomplete"},
		"default_tax_rates[0]": {taxRateID},
	})
	if upgraded.LatestInvoice == "" || upgraded.LatestInvoice == createInvoiceID {
		t.Fatalf("latest_invoice = %q, want new proration invoice (create was %q)", upgraded.LatestInvoice, createInvoiceID)
	}
	if upgraded.CurrentPeriodStart != periodStart || upgraded.CurrentPeriodEnd != periodEnd {
		t.Fatalf("period changed = %d-%d, want unchanged %d-%d", upgraded.CurrentPeriodStart, upgraded.CurrentPeriodEnd, periodStart, periodEnd)
	}
	if upgraded.Items.Data[0].Price.ID != pro.ID {
		t.Fatalf("price = %s, want pro", upgraded.Items.Data[0].Price.ID)
	}

	invoice := getJSON[prorationInvoiceResponse](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	// remaining ≈ full period: subtotal 5000, tax 500, total 5500
	if invoice.Subtotal != 5000 || invoice.Tax != 500 || invoice.Total != 5500 || invoice.AmountPaid != 5500 || invoice.AmountDue != 0 {
		t.Fatalf("proration invoice = %#v, want subtotal=5000 tax=500 total=5500 paid=5500", invoice)
	}
	if invoice.BillingReason != "subscription_update" {
		t.Fatalf("billing_reason = %q, want subscription_update", invoice.BillingReason)
	}

	invoices := listSubscriptionInvoices(t, handler, created.ID)
	if len(invoices) != 2 {
		t.Fatalf("invoice count = %d, want 2", len(invoices))
	}
}

// TestSubscriptionUpdateAlwaysInvoiceBillingCycleAnchorNow: net charge after pre-discount
// credit (subtotal = newTotal - creditSubtotal). Full remaining period → 5000/0/500/5500.
// (Previously subtotal was wrongly newTotal 9900 with credit only in metadata, which broke
// Stripe's subtotal - discount + tax = total identity and made total_taxes recompute as 990.)
func TestSubscriptionUpdateAlwaysInvoiceBillingCycleAnchorNow(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	// Fixed period; proration_date = period start so remaining == full period.
	periodStart := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2030, 1, 31, 0, 0, 0, 0, time.UTC)
	at := periodStart
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name":  "anchor-now-proration",
		"runId": "anchor-now-1",
		"subscriptions": []map[string]any{{
			"id":                   "sub_anchor_now",
			"customer":             customer.ID,
			"price":                lite.ID,
			"status":               "active",
			"current_period_start": periodStart.Format(time.RFC3339),
			"current_period_end":   periodEnd.Format(time.RFC3339),
		}},
	})
	if len(applied.Subscriptions) != 1 {
		t.Fatalf("fixture = %#v", applied)
	}
	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/sub_anchor_now")

	// creditSubtotal = 4900, subtotal = 9900−4900 = 5000, base = 5000, tax 500, total 5500
	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+sub.ID, url.Values{
		"items[0][id]":         {sub.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"billing_cycle_anchor": {"now"},
		"proration_date":       {strconv.FormatInt(at.Unix(), 10)},
		"default_tax_rates[0]": {taxRateID},
	})
	if upgraded.CurrentPeriodStart != at.Unix() {
		t.Fatalf("period start = %d, want proration_date %d", upgraded.CurrentPeriodStart, at.Unix())
	}
	wantEnd := at.AddDate(0, 1, 0)
	if upgraded.CurrentPeriodEnd != wantEnd.Unix() {
		t.Fatalf("period end = %d, want %d", upgraded.CurrentPeriodEnd, wantEnd.Unix())
	}
	if upgraded.CurrentPeriodEnd == periodEnd.Unix() {
		t.Fatal("period end still old bound; anchor=now did not reset")
	}

	invoice := getJSON[struct {
		prorationInvoiceResponse
		TotalExcludingTax int64             `json:"total_excluding_tax"`
		Metadata          map[string]string `json:"metadata"`
		TotalDiscountAmounts []struct {
			Amount int64 `json:"amount"`
		} `json:"total_discount_amounts"`
		TotalTaxes []struct {
			Amount         int64 `json:"amount"`
			TaxableAmount  int64 `json:"taxable_amount"`
		} `json:"total_taxes"`
	}](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	if invoice.Subtotal != 5000 || invoice.Tax != 500 || invoice.Total != 5500 || invoice.AmountPaid != 5500 {
		t.Fatalf("anchor=now invoice = %#v, want subtotal=5000 tax=500 total=5500 paid=5500", invoice)
	}
	if invoice.BillingReason != "subscription_update" {
		t.Fatalf("billing_reason = %q", invoice.BillingReason)
	}
	if invoice.Metadata[billing.MetadataProrationCredit] != "4900" {
		t.Fatalf("proration credit metadata = %q, want 4900 (pre-discount)", invoice.Metadata[billing.MetadataProrationCredit])
	}
	assertInvoiceAmountIdentity(t, invoice.Subtotal, invoice.TotalDiscountAmounts, invoice.Tax, invoice.Total, invoice.TotalExcludingTax, invoice.TotalTaxes)
}

// TestSubscriptionUpdateMidPeriodMatchesPreview: proration_date mid-period must match create_preview.
func TestSubscriptionUpdateMidPeriodMatchesPreview(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, _ := setupProrationPlans(t, handler)

	periodStart := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2030, 1, 31, 0, 0, 0, 0, time.UTC)
	mid := time.Date(2030, 1, 16, 0, 0, 0, 0, time.UTC)
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name":  "mid-period-proration",
		"runId": "mid-period-1",
		"subscriptions": []map[string]any{{
			"id":                   "sub_mid_proration",
			"customer":             customer.ID,
			"price":                lite.ID,
			"status":               "active",
			"current_period_start": periodStart.Format(time.RFC3339),
			"current_period_end":   periodEnd.Format(time.RFC3339),
		}},
	})
	if len(applied.Subscriptions) != 1 {
		t.Fatalf("fixture = %#v", applied)
	}
	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/sub_mid_proration")

	// (9900-4900)*remaining/period via truncating integer division
	remaining := periodEnd.Unix() - mid.Unix()
	periodSec := periodEnd.Unix() - periodStart.Unix()
	wantDelta := (9900 - 4900) * remaining / periodSec

	preview := postForm[struct {
		AmountDue int64 `json:"amount_due"`
		Subtotal  int64 `json:"subtotal"`
	}](t, handler, "/v1/invoices/create_preview", url.Values{
		"subscription":                             {"sub_mid_proration"},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[items][0][quantity]": {"1"},
		"subscription_details[proration_behavior]": {"create_prorations"},
		"subscription_details[proration_date]":     {strconv.FormatInt(mid.Unix(), 10)},
	})
	if preview.AmountDue != wantDelta || preview.Subtotal != wantDelta {
		t.Fatalf("preview amount_due/subtotal = %d/%d, want %d", preview.AmountDue, preview.Subtotal, wantDelta)
	}

	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+sub.ID, url.Values{
		"items[0][id]":       {sub.Items.Data[0].ID},
		"items[0][price]":    {pro.ID},
		"proration_behavior": {"always_invoice"},
		"proration_date":     {strconv.FormatInt(mid.Unix(), 10)},
	})
	invoice := getJSON[prorationInvoiceResponse](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	if invoice.Subtotal != wantDelta || invoice.Total != wantDelta || invoice.AmountPaid != wantDelta {
		t.Fatalf("invoice = %#v, want total/subtotal=%d matching preview", invoice, wantDelta)
	}
	if invoice.Subtotal != preview.Subtotal || invoice.Total != preview.AmountDue {
		t.Fatalf("invoice vs preview: invoice subtotal=%d total=%d preview subtotal=%d amount_due=%d",
			invoice.Subtotal, invoice.Total, preview.Subtotal, preview.AmountDue)
	}
}

// TestSubscriptionUpdateCreateProrationsDefersToRenewal: pending metadata then renewal absorbs it.
func TestSubscriptionUpdateCreateProrationsDefersToRenewal(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	periodStart := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2030, 2, 1, 0, 0, 0, 0, time.UTC)
	// mid: Jan 16 → remaining = 16 days of 31
	mid := time.Date(2030, 1, 16, 0, 0, 0, 0, time.UTC)
	remaining := periodEnd.Unix() - mid.Unix()
	periodSec := periodEnd.Unix() - periodStart.Unix()
	// discounted delta (no coupon): (9900-4900)*rem/period
	wantPending := (9900 - 4900) * remaining / periodSec

	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "create-prorations-renewal",
		"test_clocks": []map[string]any{{
			"id":          "clock_create_prorations",
			"frozen_time": periodStart.Format(time.RFC3339),
		}},
		"subscriptions": []map[string]any{{
			"id":                   "sub_create_prorations",
			"customer":             customer.ID,
			"price":                lite.ID,
			"status":               "active",
			"test_clock":           "clock_create_prorations",
			"current_period_start": periodStart.Format(time.RFC3339),
			"current_period_end":   periodEnd.Format(time.RFC3339),
			"default_tax_rates":    []string{taxRateID},
		}},
	})
	if len(applied.Subscriptions) != 1 {
		t.Fatalf("fixture = %#v", applied)
	}
	// Attach tax rate snapshot if fixture did not.
	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/sub_create_prorations")
	beforeInvoices := listSubscriptionInvoices(t, handler, sub.ID)

	updated := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+sub.ID, url.Values{
		"items[0][id]":         {sub.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"create_prorations"},
		"proration_date":       {strconv.FormatInt(mid.Unix(), 10)},
		"default_tax_rates[0]": {taxRateID},
	})
	afterUpdate := listSubscriptionInvoices(t, handler, sub.ID)
	if len(afterUpdate) != len(beforeInvoices) {
		t.Fatalf("create_prorations created invoice: before=%d after=%d", len(beforeInvoices), len(afterUpdate))
	}
	if updated.Metadata[billing.MetadataPendingProrationAmount] != strconv.FormatInt(wantPending, 10) {
		t.Fatalf("pending amount = %q, want %d", updated.Metadata[billing.MetadataPendingProrationAmount], wantPending)
	}
	if updated.Metadata[billing.MetadataPendingProrationAt] == "" {
		t.Fatal("pending at metadata missing")
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
	}](t, handler, "/v1/test_helpers/test_clocks/clock_create_prorations/advance", url.Values{
		"frozen_time": {strconv.FormatInt(periodEnd.Unix(), 10)},
	})
	if len(advance.BilltapAdvanceResult.Renewals) != 1 {
		t.Fatalf("renewals = %#v, want 1", advance.BilltapAdvanceResult.Renewals)
	}
	// renewal subtotal = new plan 9900 + pending; tax on discounted base (=subtotal without coupon)
	renewal := advance.BilltapAdvanceResult.Renewals[0].Invoice
	wantSubtotal := 9900 + wantPending
	wantTax := int64(float64(wantSubtotal)*0.10 + 0.5) // ComputeTaxRateAmounts uses Round
	// Match ComputeTaxRateAmounts exactly:
	_, _, _, wantTaxExact := billing.ComputeTaxRateAmounts(wantSubtotal, []billing.AppliedTaxRate{{Percentage: 10, Inclusive: false}})
	if renewal.Subtotal != wantSubtotal {
		t.Fatalf("renewal subtotal = %d, want %d (9900+pending %d)", renewal.Subtotal, wantSubtotal, wantPending)
	}
	if renewal.Tax != wantTaxExact {
		t.Fatalf("renewal tax = %d, want %d (approx %d)", renewal.Tax, wantTaxExact, wantTax)
	}
	if renewal.Total != wantSubtotal+wantTaxExact {
		t.Fatalf("renewal total = %d, want %d", renewal.Total, wantSubtotal+wantTaxExact)
	}

	renewed := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/sub_create_prorations")
	if renewed.Metadata[billing.MetadataPendingProrationAmount] != "" {
		t.Fatalf("pending metadata still present: %v", renewed.Metadata)
	}
}

// TestSubscriptionUpdateProrationNone: no invoice, period unchanged.
func TestSubscriptionUpdateProrationNone(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, _ := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {lite.ID},
	})
	before := listSubscriptionInvoices(t, handler, created.ID)
	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":       {created.Items.Data[0].ID},
		"items[0][price]":    {pro.ID},
		"proration_behavior": {"none"},
	})
	after := listSubscriptionInvoices(t, handler, created.ID)
	if len(after) != len(before) {
		t.Fatalf("none created invoices: before=%d after=%d", len(before), len(after))
	}
	if upgraded.LatestInvoice != created.LatestInvoice {
		t.Fatalf("latest_invoice changed under none")
	}
	if upgraded.CurrentPeriodStart != created.CurrentPeriodStart || upgraded.CurrentPeriodEnd != created.CurrentPeriodEnd {
		t.Fatal("period changed under none")
	}
	if upgraded.Items.Data[0].Price.ID != pro.ID {
		t.Fatal("item not updated")
	}
}

// TestSubscriptionUpdateDowngradeNoInvoice: delta ≤ 0 → items change, no new invoice.
func TestSubscriptionUpdateDowngradeNoInvoice(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, _ := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {pro.ID},
	})
	before := listSubscriptionInvoices(t, handler, created.ID)
	downgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":       {created.Items.Data[0].ID},
		"items[0][price]":    {lite.ID},
		"proration_behavior": {"always_invoice"},
	})
	after := listSubscriptionInvoices(t, handler, created.ID)
	if len(after) != len(before) {
		t.Fatalf("downgrade created invoice: before=%d after=%d", len(before), len(after))
	}
	if downgraded.Items.Data[0].Price.ID != lite.ID {
		t.Fatal("item not downgraded")
	}
	if downgraded.LatestInvoice != created.LatestInvoice {
		t.Fatal("latest_invoice should stay on create invoice for downgrade")
	}
}

// TestSubscriptionUpdateErrorIfIncompleteRollsBack: 402, items unchanged, no invoice.
func TestSubscriptionUpdateErrorIfIncompleteRollsBack(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"fail-upgrade@example.test"},
		"metadata[billtap_default_invoice_outcome]": {"card_declined"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Fail Plan"}})
	lite := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product": {product.ID}, "currency": {"usd"}, "unit_amount": {"4900"}, "recurring[interval]": {"month"},
	})
	pro := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product": {product.ID}, "currency": {"usd"}, "unit_amount": {"9900"}, "recurring[interval]": {"month"},
	})

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {lite.ID},
	})
	before := listSubscriptionInvoices(t, handler, created.ID)

	status, body := postFormStatus(t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":       {created.Items.Data[0].ID},
		"items[0][price]":    {pro.ID},
		"proration_behavior": {"always_invoice"},
		"payment_behavior":   {"error_if_incomplete"},
	})
	if status != http.StatusPaymentRequired {
		t.Fatalf("status = %d body=%s, want 402", status, body)
	}
	var errBody struct {
		Error struct {
			Type string `json:"type"`
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(body), &errBody); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Error.Type != "card_error" || errBody.Error.Code != "card_declined" {
		t.Fatalf("error = %#v, want card_error/card_declined", errBody.Error)
	}

	got := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if got.Items.Data[0].Price.ID != lite.ID {
		t.Fatalf("items changed to %s despite error_if_incomplete", got.Items.Data[0].Price.ID)
	}
	after := listSubscriptionInvoices(t, handler, created.ID)
	if len(after) != len(before) {
		t.Fatalf("invoice created on rollback path: before=%d after=%d", len(before), len(after))
	}
}

// TestSubscriptionUpdateAlwaysInvoiceWebhookOrder: renewal-style event sequence.
func TestSubscriptionUpdateAlwaysInvoiceWebhookOrder(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, _ := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {lite.ID},
	})
	// Drain existing events by noting count.
	beforeEvents := getJSON[struct {
		Data []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"data"`
	}](t, handler, "/v1/events?limit=100")
	beforeIDs := map[string]bool{}
	for _, e := range beforeEvents.Data {
		beforeIDs[e.ID] = true
	}

	_ = postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":       {created.Items.Data[0].ID},
		"items[0][price]":    {pro.ID},
		"proration_behavior": {"always_invoice"},
	})

	events := getJSON[struct {
		Data []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"data"`
	}](t, handler, "/v1/events?limit=100")

	// Events are typically newest-first; collect new ones in reverse chronological then reverse.
	var newTypes []string
	for i := len(events.Data) - 1; i >= 0; i-- {
		e := events.Data[i]
		if beforeIDs[e.ID] {
			continue
		}
		newTypes = append(newTypes, e.Type)
	}
	// Also try chronological from API order (newest first): rebuild from newest-first scan.
	var newestFirst []string
	for _, e := range events.Data {
		if beforeIDs[e.ID] {
			continue
		}
		newestFirst = append(newestFirst, e.Type)
	}
	// Prefer chronological sequence (oldest first) for order assertion.
	want := []string{
		"invoice.created",
		"invoice.finalized",
		"payment_intent.created",
		"payment_intent.succeeded",
		"invoice.payment_succeeded",
		"invoice.paid",
		"customer.subscription.updated",
	}
	// Find a contiguous subsequence matching want in either chronological direction.
	if !containsEventOrder(newTypes, want) && !containsEventOrder(reverseStrings(newestFirst), want) && !containsEventOrder(newestFirst, want) {
		t.Fatalf("new event types (chronological) = %v (newest-first %v), want order %v", newTypes, newestFirst, want)
	}
}

func containsEventOrder(got, want []string) bool {
	if len(want) == 0 {
		return true
	}
	wi := 0
	for _, g := range got {
		if g == want[wi] {
			wi++
			if wi == len(want) {
				return true
			}
		}
	}
	return false
}

func reverseStrings(in []string) []string {
	out := make([]string, len(in))
	for i := range in {
		out[len(in)-1-i] = in[i]
	}
	return out
}

// TestSubscriptionUpdateAlwaysInvoiceWithDiscountAndTax: 25% coupon, then 10% VAT on discounted delta.
func TestSubscriptionUpdateAlwaysInvoiceWithDiscountAndTax(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"percent_off": {"25"},
		"duration":    {"forever"},
	})

	// Create with coupon so subscription metadata carries the discount.
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"coupon":               {coupon.ID},
		"default_tax_rates[0]": {taxRateID},
	})

	// old discounted 4900*0.75=3675, new 9900*0.75=7425, delta 3750, tax 375, total 4125
	// Pin proration_date to period start for full remaining fraction.
	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {created.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"proration_date":       {strconv.FormatInt(created.CurrentPeriodStart, 10)},
		"default_tax_rates[0]": {taxRateID},
	})
	invoice := getJSON[struct {
		prorationInvoiceResponse
		TotalExcludingTax    int64 `json:"total_excluding_tax"`
		TotalDiscountAmounts []struct {
			Amount int64 `json:"amount"`
		} `json:"total_discount_amounts"`
		TotalTaxes []struct {
			Amount        int64 `json:"amount"`
			TaxableAmount int64 `json:"taxable_amount"`
		} `json:"total_taxes"`
	}](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	// subtotal is pre-discount delta 5000; discounted base 3750; tax Round(3750*0.10)=375; total 4125
	wantSubtotal := int64(5000)
	wantBase := int64(3750)
	_, _, _, wantTax := billing.ComputeTaxRateAmounts(wantBase, []billing.AppliedTaxRate{{Percentage: 10, Inclusive: false}})
	wantTotal := wantBase + wantTax
	if invoice.Subtotal != wantSubtotal || invoice.Tax != wantTax || invoice.Total != wantTotal || invoice.AmountPaid != wantTotal {
		t.Fatalf("discounted proration invoice = %#v, want subtotal=%d tax=%d total=%d", invoice, wantSubtotal, wantTax, wantTotal)
	}
	assertInvoiceAmountIdentity(t, invoice.Subtotal, invoice.TotalDiscountAmounts, invoice.Tax, invoice.Total, invoice.TotalExcludingTax, invoice.TotalTaxes)
}

// invoiceAmountView is the serialized invoice fields needed for Stripe amount identity.
type invoiceAmountView struct {
	Subtotal             int64 `json:"subtotal"`
	Tax                  int64 `json:"tax"`
	Total                int64 `json:"total"`
	TotalExcludingTax    int64 `json:"total_excluding_tax"`
	BillingReason        string `json:"billing_reason"`
	TotalDiscountAmounts []struct {
		Amount int64 `json:"amount"`
	} `json:"total_discount_amounts"`
	TotalTaxes []struct {
		Amount        int64 `json:"amount"`
		TaxableAmount int64 `json:"taxable_amount"`
	} `json:"total_taxes"`
	Metadata map[string]string `json:"metadata"`
}

// assertInvoiceAmountIdentity checks:
//   total == subtotal - sum(total_discount_amounts) + tax
//   tax == sum(total_taxes[].amount)
//   total_excluding_tax == pretax base (subtotal - discount for exclusive rates)
func assertInvoiceAmountIdentity(t *testing.T, subtotal int64, discounts []struct {
	Amount int64 `json:"amount"`
}, tax, total, totalExcludingTax int64, totalTaxes []struct {
	Amount        int64 `json:"amount"`
	TaxableAmount int64 `json:"taxable_amount"`
}) {
	t.Helper()
	discountSum := int64(0)
	for _, d := range discounts {
		discountSum += d.Amount
	}
	if got := subtotal - discountSum + tax; got != total {
		t.Fatalf("amount identity broken: subtotal(%d) - discount(%d) + tax(%d) = %d, want total %d",
			subtotal, discountSum, tax, got, total)
	}
	taxSum := int64(0)
	for _, entry := range totalTaxes {
		taxSum += entry.Amount
	}
	if taxSum != tax {
		t.Fatalf("stored tax %d != sum(total_taxes)=%d (entries=%#v)", tax, taxSum, totalTaxes)
	}
	// Exclusive-tax path: total_excluding_tax is pretax base = subtotal - discount.
	wantPretax := subtotal - discountSum
	if wantPretax < 0 {
		wantPretax = 0
	}
	if totalExcludingTax != wantPretax {
		t.Fatalf("total_excluding_tax = %d, want pretax base %d", totalExcludingTax, wantPretax)
	}
	for _, entry := range totalTaxes {
		if entry.TaxableAmount != wantPretax {
			t.Fatalf("total_taxes taxable_amount = %d, want pretax %d", entry.TaxableAmount, wantPretax)
		}
	}
}

// TestSubscriptionProrationInvoiceAmountIdentity covers mid-cycle delta, anchor=now,
// coupon, and fractional tax rate — each must satisfy Stripe amount identity and
// tax serialization consistency.
func TestSubscriptionProrationInvoiceAmountIdentity(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, vat10ID := setupProrationPlans(t, handler)
	// Fractional tax rate (8.875%) to stress rounding vs stored tax.
	vatFrac := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"NYC"},
		"percentage":   {"8.875"},
		"inclusive":    {"false"},
	})
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"percent_off": {"25"},
		"duration":    {"forever"},
	})

	periodStart := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2030, 1, 31, 0, 0, 0, 0, time.UTC)
	mid := time.Date(2030, 1, 16, 0, 0, 0, 0, time.UTC)

	type caseSpec struct {
		name     string
		subID    string
		price    string
		coupon   string
		taxRate  string
		params   url.Values
		// optional numeric checks
		wantSubtotal *int64
		wantTotal    *int64
	}
	i64 := func(v int64) *int64 { return &v }

	// Seed four independent fixtures so cases don't share state.
	for _, seed := range []struct {
		id, price, coupon string
	}{
		{"sub_id_mid", lite.ID, ""},
		{"sub_id_anchor", lite.ID, ""},
		{"sub_id_coupon", lite.ID, coupon.ID},
		{"sub_id_frac", lite.ID, ""},
		{"sub_id_anchor_coupon", lite.ID, coupon.ID},
	} {
		subFixture := map[string]any{
			"id":                   seed.id,
			"customer":             customer.ID,
			"price":                seed.price,
			"status":               "active",
			"current_period_start": periodStart.Format(time.RFC3339),
			"current_period_end":   periodEnd.Format(time.RFC3339),
		}
		if seed.coupon != "" {
			subFixture["coupon"] = seed.coupon
			subFixture["discount_percent_off"] = 25
		}
		applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
			"name":          "identity-" + seed.id,
			"runId":         "identity-run",
			"subscriptions": []map[string]any{subFixture},
		})
		if len(applied.Subscriptions) != 1 {
			t.Fatalf("seed %s: %#v", seed.id, applied)
		}
	}

	cases := []caseSpec{
		{
			name:  "mid-cycle delta VAT10",
			subID: "sub_id_mid",
			params: url.Values{
				"proration_behavior":   {"always_invoice"},
				"proration_date":       {strconv.FormatInt(mid.Unix(), 10)},
				"default_tax_rates[0]": {vat10ID},
			},
			// (9900-4900)*rem/period
			wantSubtotal: i64((9900 - 4900) * (periodEnd.Unix() - mid.Unix()) / (periodEnd.Unix() - periodStart.Unix())),
		},
		{
			name:  "anchor=now full remaining VAT10",
			subID: "sub_id_anchor",
			params: url.Values{
				"proration_behavior":   {"always_invoice"},
				"billing_cycle_anchor": {"now"},
				"proration_date":       {strconv.FormatInt(periodStart.Unix(), 10)},
				"default_tax_rates[0]": {vat10ID},
			},
			// subtotal = 9900 - 4900 = 5000 (pre-discount credit netted into subtotal)
			wantSubtotal: i64(5000),
			wantTotal:    i64(5500),
		},
		{
			name:  "mid-cycle with 25% coupon VAT10",
			subID: "sub_id_coupon",
			params: url.Values{
				"proration_behavior":   {"always_invoice"},
				"proration_date":       {strconv.FormatInt(periodStart.Unix(), 10)},
				"default_tax_rates[0]": {vat10ID},
			},
			// pre-discount delta 5000, base 3750, tax 375, total 4125
			wantSubtotal: i64(5000),
			wantTotal:    i64(4125),
		},
		{
			name:  "mid-cycle fractional tax 8.875%",
			subID: "sub_id_frac",
			params: url.Values{
				"proration_behavior":   {"always_invoice"},
				"proration_date":       {strconv.FormatInt(periodStart.Unix(), 10)},
				"default_tax_rates[0]": {vatFrac.ID},
			},
			wantSubtotal: i64(5000),
		},
		{
			name:  "anchor=now with 25% coupon VAT10",
			subID: "sub_id_anchor_coupon",
			params: url.Values{
				"proration_behavior":   {"always_invoice"},
				"billing_cycle_anchor": {"now"},
				"proration_date":       {strconv.FormatInt(periodStart.Unix(), 10)},
				"default_tax_rates[0]": {vat10ID},
			},
			// creditSubtotal=4900, subtotal=max(0,9900-4900)=5000
			// creditDiscounted=3675, base=max(0,7425-3675)=3750, tax=375, total=4125
			// discount = 5000-3750 = 1250
			wantSubtotal: i64(5000),
			wantTotal:    i64(4125),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+tc.subID)
			params := url.Values{}
			for k, v := range tc.params {
				params[k] = v
			}
			params.Set("items[0][id]", sub.Items.Data[0].ID)
			params.Set("items[0][price]", pro.ID)

			upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+tc.subID, params)
			if upgraded.LatestInvoice == "" {
				t.Fatal("expected proration invoice")
			}
			invoice := getJSON[invoiceAmountView](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
			if invoice.BillingReason != "subscription_update" {
				t.Fatalf("billing_reason = %q", invoice.BillingReason)
			}
			if tc.wantSubtotal != nil && invoice.Subtotal != *tc.wantSubtotal {
				t.Fatalf("subtotal = %d, want %d", invoice.Subtotal, *tc.wantSubtotal)
			}
			if tc.wantTotal != nil && invoice.Total != *tc.wantTotal {
				t.Fatalf("total = %d, want %d", invoice.Total, *tc.wantTotal)
			}
			assertInvoiceAmountIdentity(t, invoice.Subtotal, invoice.TotalDiscountAmounts, invoice.Tax, invoice.Total, invoice.TotalExcludingTax, invoice.TotalTaxes)

			// anchor=now credit metadata: pre-discount primary key
			if tc.params.Get("billing_cycle_anchor") == "now" {
				if invoice.Metadata[billing.MetadataProrationCredit] == "" {
					t.Fatal("missing billtap_proration_credit on anchor=now invoice")
				}
				if tc.subID == "sub_id_anchor_coupon" {
					// discounted credit differs from pre-discount
					if invoice.Metadata[billing.MetadataProrationCreditDiscounted] == "" {
						t.Fatalf("want billtap_proration_credit_discounted when coupon applies; meta=%v", invoice.Metadata)
					}
				}
			}
		})
	}
}
