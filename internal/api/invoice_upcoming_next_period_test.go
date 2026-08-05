package api

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/fixtures"
)

// upcomingPreviewFields covers next-period upcoming / create_preview assertions.
type upcomingPreviewFields struct {
	Object               string `json:"object"`
	Subtotal             int64  `json:"subtotal"`
	Tax                  *int64 `json:"tax"`
	Total                int64  `json:"total"`
	TotalExcludingTax    int64  `json:"total_excluding_tax"`
	AmountDue            int64  `json:"amount_due"`
	AmountRemaining      int64  `json:"amount_remaining"`
	BillingReason        string `json:"billing_reason"`
	PeriodStart          int64  `json:"period_start"`
	PeriodEnd            int64  `json:"period_end"`
	TotalDiscountAmounts []struct {
		Amount int64 `json:"amount"`
	} `json:"total_discount_amounts"`
	TotalTaxes []struct {
		Amount        int64 `json:"amount"`
		TaxableAmount int64 `json:"taxable_amount"`
	} `json:"total_taxes"`
	BilltapPreview struct {
		ProrationSkippedReason *string `json:"proration_skipped_reason"`
	} `json:"billtap_preview"`
	Lines struct {
		Data []struct {
			Amount      int64  `json:"amount"`
			Proration   bool   `json:"proration"`
			Description string `json:"description"`
			Period      struct {
				Start int64 `json:"start"`
				End   int64 `json:"end"`
			} `json:"period"`
			Parent struct {
				Type                    string `json:"type"`
				SubscriptionItemDetails struct {
					Price     string `json:"price"`
					Proration bool   `json:"proration"`
				} `json:"subscription_item_details"`
			} `json:"parent"`
		} `json:"data"`
	} `json:"lines"`
}

func assertInvoiceInvariant(t *testing.T, p upcomingPreviewFields) {
	t.Helper()
	discountSum := int64(0)
	for _, d := range p.TotalDiscountAmounts {
		discountSum += d.Amount
	}
	tax := int64(0)
	if p.Tax != nil {
		tax = *p.Tax
	}
	want := p.Subtotal - discountSum + tax
	if p.Total != want {
		t.Fatalf("invariant total=%d != subtotal(%d) - discounts(%d) + tax(%d) = %d",
			p.Total, p.Subtotal, discountSum, tax, want)
	}
}

func taxVal(p upcomingPreviewFields) int64 {
	if p.Tax == nil {
		return 0
	}
	return *p.Tax
}

// TestUpcomingNextPeriodWithTaxNoParams: GET /v1/invoices/upcoming?subscription=…
// with default_tax_rates must return next-cycle invoice (not zero proration).
func TestUpcomingNextPeriodWithTaxNoParams(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, _, taxRateID := setupProrationPlans(t, handler)

	sub := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})

	preview := getJSON[upcomingPreviewFields](t, handler, "/v1/invoices/upcoming?subscription="+sub.ID)
	if preview.Object != "invoice" || preview.BillingReason != "upcoming" {
		t.Fatalf("preview = %#v, want object=invoice billing_reason=upcoming", preview)
	}
	if preview.Subtotal != 4900 || taxVal(preview) != 490 || preview.Total != 5390 || preview.AmountDue != 5390 {
		t.Fatalf("amounts = subtotal=%d tax=%d total=%d due=%d, want 4900/490/5390/5390",
			preview.Subtotal, taxVal(preview), preview.Total, preview.AmountDue)
	}
	if preview.TotalExcludingTax != 4900 || preview.AmountRemaining != 5390 {
		t.Fatalf("excl/remaining = %d/%d, want 4900/5390", preview.TotalExcludingTax, preview.AmountRemaining)
	}
	if len(preview.Lines.Data) != 1 || preview.Lines.Data[0].Amount != 4900 || preview.Lines.Data[0].Proration {
		t.Fatalf("lines = %#v, want 1 line amount=4900 proration=false", preview.Lines.Data)
	}
	if preview.PeriodStart != sub.CurrentPeriodEnd {
		t.Fatalf("period_start = %d, want current_period_end %d", preview.PeriodStart, sub.CurrentPeriodEnd)
	}
	if preview.PeriodEnd <= preview.PeriodStart {
		t.Fatalf("period_end = %d must be after period_start %d", preview.PeriodEnd, preview.PeriodStart)
	}
	if preview.Lines.Data[0].Period.Start != preview.PeriodStart || preview.Lines.Data[0].Period.End != preview.PeriodEnd {
		t.Fatalf("line period = %#v, want invoice period", preview.Lines.Data[0].Period)
	}
	if preview.BilltapPreview.ProrationSkippedReason != nil {
		t.Fatalf("proration_skipped_reason = %v, want null", preview.BilltapPreview.ProrationSkippedReason)
	}
	assertInvoiceInvariant(t, preview)
}

// TestCreatePreviewNextPeriodSubscriptionOnly: POST create_preview with only
// customer+subscription matches GET upcoming (no item overrides).
func TestCreatePreviewNextPeriodSubscriptionOnly(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, _, taxRateID := setupProrationPlans(t, handler)

	sub := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})

	viaGet := getJSON[upcomingPreviewFields](t, handler, "/v1/invoices/upcoming?subscription="+sub.ID)
	viaPost := postForm[upcomingPreviewFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"customer":     {customer.ID},
		"subscription": {sub.ID},
	})
	if viaPost.BillingReason != "upcoming" || viaPost.Subtotal != viaGet.Subtotal || viaPost.Total != viaGet.Total || taxVal(viaPost) != taxVal(viaGet) {
		t.Fatalf("create_preview=%#v upcoming=%#v, want identical next-period amounts", viaPost, viaGet)
	}
	if viaPost.PeriodStart != viaGet.PeriodStart || viaPost.PeriodEnd != viaGet.PeriodEnd {
		t.Fatalf("period mismatch post=%d/%d get=%d/%d", viaPost.PeriodStart, viaPost.PeriodEnd, viaGet.PeriodStart, viaGet.PeriodEnd)
	}
	if len(viaPost.Lines.Data) != 1 || viaPost.Lines.Data[0].Proration {
		t.Fatalf("lines = %#v, want 1 non-proration line", viaPost.Lines.Data)
	}
	assertInvoiceInvariant(t, viaPost)
}

// TestUpcomingItemOverrideKeepsProrationPath: item override → subscription_update + proration.
func TestUpcomingItemOverrideKeepsProrationPath(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})
	itemID := created.Items.Data[0].ID
	periodStart := created.CurrentPeriodStart

	preview := postForm[upcomingPreviewFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"customer":                                 {customer.ID},
		"subscription":                             {created.ID},
		"subscription_details[items][0][id]":       {itemID},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(periodStart, 10)},
	})
	if preview.BillingReason != "subscription_update" {
		t.Fatalf("billing_reason = %q, want subscription_update", preview.BillingReason)
	}
	if preview.Subtotal != 5000 || taxVal(preview) != 500 || preview.Total != 5500 {
		t.Fatalf("proration preview = %#v, want 5000/500/5500", preview)
	}
	if len(preview.Lines.Data) != 1 || !preview.Lines.Data[0].Proration {
		t.Fatalf("lines = %#v, want one proration line", preview.Lines.Data)
	}
	assertInvoiceInvariant(t, preview)
}

// TestUpcomingNextPeriodMultiItem: base 4900 + seat 1000×2 → subtotal 6900, 2 lines.
func TestUpcomingNextPeriodMultiItem(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"multi-upcoming@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Multi Upcoming"}})
	base := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"4900"},
		"recurring[interval]": {"month"},
	})
	seat := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1000"},
		"recurring[interval]": {"month"},
	})
	txr := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	sub := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {base.ID},
		"items[0][quantity]":   {"1"},
		"items[1][price]":      {seat.ID},
		"items[1][quantity]":   {"2"},
		"default_tax_rates[0]": {txr.ID},
	})

	preview := getJSON[upcomingPreviewFields](t, handler, "/v1/invoices/upcoming?subscription="+sub.ID)
	if preview.Subtotal != 6900 || taxVal(preview) != 690 || preview.Total != 7590 {
		t.Fatalf("multi amounts = %#v, want subtotal=6900 tax=690 total=7590", preview)
	}
	if len(preview.Lines.Data) != 2 {
		t.Fatalf("lines = %#v, want 2", preview.Lines.Data)
	}
	// Line amounts are per-item (4900 and 2000), order follows subscription items.
	got := map[int64]int{}
	for _, line := range preview.Lines.Data {
		if line.Proration {
			t.Fatalf("line proration true: %#v", line)
		}
		got[line.Amount]++
	}
	if got[4900] != 1 || got[2000] != 1 {
		t.Fatalf("line amounts = %#v, want 4900 and 2000", preview.Lines.Data)
	}
	assertInvoiceInvariant(t, preview)
}

// TestUpcomingNextPeriodDiscountAndTax: 25% coupon + 10% tax, tax after discount.
func TestUpcomingNextPeriodDiscountAndTax(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"disc-upcoming@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Disc Upcoming"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"4900"},
		"recurring[interval]": {"month"},
	})
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_upcoming_25"},
		"percent_off": {"25"},
		"duration":    {"forever"},
	})
	txr := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	sub := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {price.ID},
		"items[0][quantity]":   {"1"},
		"discounts[0][coupon]": {coupon.ID},
		"default_tax_rates[0]": {txr.ID},
	})

	// subtotal 4900, 25% off → 1225 discount → pretax 3675, tax 368, total 4043
	preview := getJSON[upcomingPreviewFields](t, handler, "/v1/invoices/upcoming?subscription="+sub.ID)
	if preview.Subtotal != 4900 {
		t.Fatalf("subtotal = %d, want 4900", preview.Subtotal)
	}
	if len(preview.TotalDiscountAmounts) != 1 || preview.TotalDiscountAmounts[0].Amount != 1225 {
		t.Fatalf("discounts = %#v, want 1225", preview.TotalDiscountAmounts)
	}
	if taxVal(preview) != 368 || preview.Total != 4043 || preview.TotalExcludingTax != 3675 {
		t.Fatalf("tax/total/excl = %d/%d/%d, want 368/4043/3675", taxVal(preview), preview.Total, preview.TotalExcludingTax)
	}
	if preview.BillingReason != "upcoming" {
		t.Fatalf("billing_reason = %q, want upcoming", preview.BillingReason)
	}
	assertInvoiceInvariant(t, preview)
}

// TestUpcomingNextPeriodIncludesPendingProration: create_prorations pending is
// added to upcoming subtotal but not consumed; renewal invoice matches upcoming.
func TestUpcomingNextPeriodIncludesPendingProration(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	periodStart := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2030, 1, 31, 0, 0, 0, 0, time.UTC)
	mid := time.Date(2030, 1, 16, 0, 0, 0, 0, time.UTC)

	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_upcoming_pending"},
		"frozen_time": {strconv.FormatInt(periodStart.Unix(), 10)},
	})
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "upcoming-pending",
		"customers": []map[string]any{{
			"id":         customer.ID,
			"email":      "upcoming-pending@example.test",
			"test_clock": "clock_upcoming_pending",
		}},
		"subscriptions": []map[string]any{{
			"id":                   "sub_upcoming_pending",
			"customer":             customer.ID,
			"price":                lite.ID,
			"status":               "active",
			"test_clock":           "clock_upcoming_pending",
			"current_period_start": periodStart.Format(time.RFC3339),
			"current_period_end":   periodEnd.Format(time.RFC3339),
			"default_tax_rates":    []string{taxRateID},
		}},
	})
	if len(applied.Subscriptions) != 1 {
		t.Fatalf("fixture = %#v", applied)
	}
	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/sub_upcoming_pending")

	// Remaining fraction from mid: 15/30 of (9900-4900)=5000 → pending 2500.
	updated := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+sub.ID, url.Values{
		"items[0][id]":         {sub.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"create_prorations"},
		"proration_date":       {strconv.FormatInt(mid.Unix(), 10)},
		"default_tax_rates[0]": {taxRateID},
	})
	pendingRaw := updated.Metadata[billing.MetadataPendingProrationAmount]
	pending, err := strconv.ParseInt(pendingRaw, 10, 64)
	if err != nil || pending == 0 {
		t.Fatalf("pending metadata = %q, want non-zero int", pendingRaw)
	}

	preview := getJSON[upcomingPreviewFields](t, handler, "/v1/invoices/upcoming?subscription="+updated.ID)
	wantSubtotal := 9900 + pending
	if preview.Subtotal != wantSubtotal {
		t.Fatalf("upcoming subtotal = %d, want %d (9900+pending %d)", preview.Subtotal, wantSubtotal, pending)
	}
	if preview.BillingReason != "upcoming" {
		t.Fatalf("billing_reason = %q, want upcoming", preview.BillingReason)
	}
	assertInvoiceInvariant(t, preview)

	// Preview must not consume pending metadata.
	still := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+updated.ID)
	if still.Metadata[billing.MetadataPendingProrationAmount] != pendingRaw {
		t.Fatalf("pending after preview = %q, want still %q", still.Metadata[billing.MetadataPendingProrationAmount], pendingRaw)
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
	}](t, handler, "/v1/test_helpers/test_clocks/clock_upcoming_pending/advance", url.Values{
		"frozen_time": {strconv.FormatInt(periodEnd.Unix(), 10)},
	})
	if len(advance.BilltapAdvanceResult.Renewals) != 1 {
		t.Fatalf("renewals = %#v, want 1", advance.BilltapAdvanceResult.Renewals)
	}
	renewal := advance.BilltapAdvanceResult.Renewals[0].Invoice
	if renewal.Subtotal != preview.Subtotal || renewal.Total != preview.Total || renewal.Tax != taxVal(preview) {
		t.Fatalf("renewal=%#v upcoming subtotal/total/tax=%d/%d/%d, want match",
			renewal, preview.Subtotal, preview.Total, taxVal(preview))
	}
}

// TestUpcomingNextPeriodTrialingUsesTrialEnd: period_start is trial_end.
func TestUpcomingNextPeriodTrialingUsesTrialEnd(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, _, taxRateID := setupProrationPlans(t, handler)

	trialEnd := time.Date(2030, 2, 15, 0, 0, 0, 0, time.UTC)
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "upcoming-trial",
		"subscriptions": []map[string]any{{
			"id":                   "sub_upcoming_trial",
			"customer":             customer.ID,
			"price":                lite.ID,
			"status":               "trialing",
			"current_period_start": "2030-01-15T00:00:00Z",
			"current_period_end":   trialEnd.Format(time.RFC3339),
			"trial_start":          "2030-01-15T00:00:00Z",
			"trial_end":            trialEnd.Format(time.RFC3339),
			"default_tax_rates":    []string{taxRateID},
		}},
	})
	if len(applied.Subscriptions) != 1 {
		t.Fatalf("fixture = %#v", applied)
	}

	preview := getJSON[upcomingPreviewFields](t, handler, "/v1/invoices/upcoming?subscription=sub_upcoming_trial")
	if preview.PeriodStart != trialEnd.Unix() {
		t.Fatalf("period_start = %d, want trial_end %d", preview.PeriodStart, trialEnd.Unix())
	}
	if preview.BillingReason != "upcoming" || preview.Subtotal != 4900 {
		t.Fatalf("preview = %#v, want upcoming with next-cycle amounts", preview)
	}
	if preview.PeriodEnd <= preview.PeriodStart {
		t.Fatalf("period_end = %d must be after start", preview.PeriodEnd)
	}
	assertInvoiceInvariant(t, preview)
}

// TestUpcomingNextPeriodNoPeriodReturnsZero: zero current_period_end → 0 + reason.
func TestUpcomingNextPeriodNoPeriodReturnsZero(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, _, _ := setupProrationPlans(t, handler)

	// Go zero time parses as IsZero() == true.
	zero := "0001-01-01T00:00:00Z"
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "upcoming-no-period",
		"subscriptions": []map[string]any{{
			"id":                   "sub_upcoming_no_period",
			"customer":             customer.ID,
			"price":                lite.ID,
			"status":               "active",
			"current_period_start": zero,
			"current_period_end":   zero,
		}},
	})
	if len(applied.Subscriptions) != 1 {
		t.Fatalf("fixture = %#v", applied)
	}

	preview := getJSON[upcomingPreviewFields](t, handler, "/v1/invoices/upcoming?subscription=sub_upcoming_no_period")
	if preview.Subtotal != 0 || preview.Total != 0 || preview.AmountDue != 0 {
		t.Fatalf("amounts = %#v, want zeros", preview)
	}
	if preview.BilltapPreview.ProrationSkippedReason == nil || *preview.BilltapPreview.ProrationSkippedReason != "no_period" {
		t.Fatalf("proration_skipped_reason = %v, want no_period", preview.BilltapPreview.ProrationSkippedReason)
	}
	if preview.BillingReason != "upcoming" {
		t.Fatalf("billing_reason = %q, want upcoming", preview.BillingReason)
	}
	if len(preview.Lines.Data) != 0 {
		t.Fatalf("lines = %#v, want empty", preview.Lines.Data)
	}
}
