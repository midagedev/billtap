package api

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/hckim/billtap/internal/billing"
)

// previewTaxFields is the subset compared for preview == confirmed invoice tax parity.
type previewTaxFields struct {
	Subtotal          int64  `json:"subtotal"`
	Tax               *int64 `json:"tax"`
	Total             int64  `json:"total"`
	TotalExcludingTax int64  `json:"total_excluding_tax"`
	AmountDue         int64  `json:"amount_due"`
	DefaultTaxRates   []struct {
		ID         string  `json:"id"`
		Percentage float64 `json:"percentage"`
		Inclusive  bool    `json:"inclusive"`
	} `json:"default_tax_rates"`
	TotalTaxes []struct {
		Amount         int64 `json:"amount"`
		TaxableAmount  int64 `json:"taxable_amount"`
		TaxRateDetails struct {
			TaxRate string `json:"tax_rate"`
		} `json:"tax_rate_details"`
	} `json:"total_taxes"`
	TotalTaxAmounts []struct {
		Amount    int64  `json:"amount"`
		Inclusive bool   `json:"inclusive"`
		TaxRate   string `json:"tax_rate"`
	} `json:"total_tax_amounts"`
}

func taxPtr(v int64) *int64 { return &v }

// TestInvoicePreviewTaxMatchesConfirmedInvoice is the core contract:
// create_preview numbers must equal the confirmed always_invoice proration invoice.
func TestInvoicePreviewTaxMatchesConfirmedInvoice(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})
	periodStart := created.CurrentPeriodStart
	itemID := created.Items.Data[0].ID

	preview := postForm[previewTaxFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"customer":                                 {customer.ID},
		"subscription":                             {created.ID},
		"subscription_details[items][0][id]":       {itemID},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(periodStart, 10)},
	})
	// Full remaining period: subtotal 5000, tax 500, total 5500.
	if preview.Subtotal != 5000 || preview.Tax == nil || *preview.Tax != 500 || preview.Total != 5500 || preview.TotalExcludingTax != 5000 {
		t.Fatalf("preview = %#v, want subtotal=5000 tax=500 total=5500 total_excluding_tax=5000", preview)
	}
	if len(preview.TotalTaxes) != 1 || preview.TotalTaxes[0].Amount != 500 || preview.TotalTaxes[0].TaxableAmount != 5000 {
		t.Fatalf("preview total_taxes = %#v, want amount=500 taxable=5000", preview.TotalTaxes)
	}
	if preview.TotalTaxes[0].TaxRateDetails.TaxRate != taxRateID {
		t.Fatalf("preview tax_rate = %q, want %s", preview.TotalTaxes[0].TaxRateDetails.TaxRate, taxRateID)
	}
	if len(preview.DefaultTaxRates) != 1 || preview.DefaultTaxRates[0].ID != taxRateID {
		t.Fatalf("preview default_tax_rates = %#v, want %s", preview.DefaultTaxRates, taxRateID)
	}

	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {itemID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"proration_date":       {strconv.FormatInt(periodStart, 10)},
		"default_tax_rates[0]": {taxRateID},
	})
	invoice := getJSON[previewTaxFields](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	if invoice.Subtotal != preview.Subtotal || invoice.Total != preview.Total || invoice.TotalExcludingTax != preview.TotalExcludingTax {
		t.Fatalf("preview vs invoice subtotal/total/excl: preview=%d/%d/%d invoice=%d/%d/%d",
			preview.Subtotal, preview.Total, preview.TotalExcludingTax,
			invoice.Subtotal, invoice.Total, invoice.TotalExcludingTax)
	}
	if invoice.Tax == nil || preview.Tax == nil || *invoice.Tax != *preview.Tax {
		t.Fatalf("preview tax=%v invoice tax=%v", preview.Tax, invoice.Tax)
	}
	if len(invoice.TotalTaxes) != 1 || invoice.TotalTaxes[0].Amount != preview.TotalTaxes[0].Amount ||
		invoice.TotalTaxes[0].TaxableAmount != preview.TotalTaxes[0].TaxableAmount ||
		invoice.TotalTaxes[0].TaxRateDetails.TaxRate != preview.TotalTaxes[0].TaxRateDetails.TaxRate {
		t.Fatalf("preview total_taxes=%#v invoice total_taxes=%#v", preview.TotalTaxes, invoice.TotalTaxes)
	}
}

// TestInvoicePreviewTaxDecimalPercentMatchesConfirmed: 8.875% mid-period must round identically.
func TestInvoicePreviewTaxDecimalPercentMatchesConfirmed(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, _ := setupProrationPlans(t, handler)
	txr := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"NY"},
		"percentage":   {"8.875"},
		"inclusive":    {"false"},
	})

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {txr.ID},
	})
	// Mid-period: half remaining of the live subscription period.
	half := (created.CurrentPeriodEnd - created.CurrentPeriodStart) / 2
	prorationDate := created.CurrentPeriodStart + half

	remaining := created.CurrentPeriodEnd - prorationDate
	periodSec := created.CurrentPeriodEnd - created.CurrentPeriodStart
	wantSubtotal := (9900 - 4900) * remaining / periodSec
	_, _, _, wantTax := billing.ComputeTaxRateAmounts(wantSubtotal, []billing.AppliedTaxRate{{Percentage: 8.875, Inclusive: false}})
	wantTotal := wantSubtotal + wantTax

	preview := postForm[previewTaxFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"subscription":                             {created.ID},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(prorationDate, 10)},
	})
	if preview.Subtotal != wantSubtotal || preview.Tax == nil || *preview.Tax != wantTax || preview.Total != wantTotal {
		t.Fatalf("preview = subtotal=%d tax=%v total=%d, want %d/%d/%d",
			preview.Subtotal, preview.Tax, preview.Total, wantSubtotal, wantTax, wantTotal)
	}

	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {created.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"proration_date":       {strconv.FormatInt(prorationDate, 10)},
		"default_tax_rates[0]": {txr.ID},
	})
	invoice := getJSON[previewTaxFields](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	if invoice.Tax == nil || preview.Tax == nil || *invoice.Tax != *preview.Tax || invoice.Total != preview.Total || invoice.Subtotal != preview.Subtotal {
		t.Fatalf("decimal preview tax=%v total=%d subtotal=%d; invoice tax=%v total=%d subtotal=%d",
			preview.Tax, preview.Total, preview.Subtotal, invoice.Tax, invoice.Total, invoice.Subtotal)
	}
}

// TestInvoicePreviewInclusiveTaxMatchesConfirmed: inclusive 10% → total stays at base.
func TestInvoicePreviewInclusiveTaxMatchesConfirmed(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, _ := setupProrationPlans(t, handler)
	txr := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT incl"},
		"percentage":   {"10"},
		"inclusive":    {"true"},
	})

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {txr.ID},
	})
	periodStart := created.CurrentPeriodStart

	// base 5000 inclusive 10% → tax Round(5000*10/110)=455, pretax 4545, total 5000
	_, wantPretax, _, wantTax := billing.ComputeTaxRateAmounts(5000, []billing.AppliedTaxRate{{Percentage: 10, Inclusive: true}})

	preview := postForm[previewTaxFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"subscription":                             {created.ID},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(periodStart, 10)},
	})
	if preview.Subtotal != 5000 || preview.Total != 5000 || preview.Tax == nil || *preview.Tax != wantTax {
		t.Fatalf("inclusive preview = %#v, want subtotal=total=5000 tax=%d", preview, wantTax)
	}
	if preview.TotalExcludingTax != wantPretax || preview.TotalExcludingTax >= preview.Subtotal {
		t.Fatalf("inclusive total_excluding_tax = %d, want %d < subtotal", preview.TotalExcludingTax, wantPretax)
	}

	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {created.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"proration_date":       {strconv.FormatInt(periodStart, 10)},
		"default_tax_rates[0]": {txr.ID},
	})
	invoice := getJSON[previewTaxFields](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	if invoice.Total != preview.Total || invoice.Tax == nil || *invoice.Tax != *preview.Tax ||
		invoice.TotalExcludingTax != preview.TotalExcludingTax || invoice.Subtotal != preview.Subtotal {
		t.Fatalf("inclusive mismatch preview=%#v invoice=%#v", preview, invoice)
	}
}

// TestInvoicePreviewDiscountAndTaxMatchesConfirmed: 25% coupon then 10% VAT on discounted base.
func TestInvoicePreviewDiscountAndTaxMatchesConfirmed(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"percent_off": {"25"},
		"duration":    {"forever"},
	})

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"coupon":               {coupon.ID},
		"default_tax_rates[0]": {taxRateID},
	})
	periodStart := created.CurrentPeriodStart

	// subtotal delta 5000, discounted base 3750, tax 375, total 4125
	wantSubtotal := int64(5000)
	wantBase := int64(3750)
	_, _, _, wantTax := billing.ComputeTaxRateAmounts(wantBase, []billing.AppliedTaxRate{{Percentage: 10, Inclusive: false}})
	wantTotal := wantBase + wantTax

	preview := postForm[previewTaxFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"subscription":                             {created.ID},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(periodStart, 10)},
	})
	if preview.Subtotal != wantSubtotal || preview.Tax == nil || *preview.Tax != wantTax || preview.Total != wantTotal {
		t.Fatalf("discount+tax preview = %#v, want subtotal=%d tax=%d total=%d", preview, wantSubtotal, wantTax, wantTotal)
	}

	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {created.Items.Data[0].ID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"proration_date":       {strconv.FormatInt(periodStart, 10)},
		"default_tax_rates[0]": {taxRateID},
	})
	invoice := getJSON[previewTaxFields](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	if invoice.Subtotal != preview.Subtotal || invoice.Total != preview.Total ||
		invoice.Tax == nil || *invoice.Tax != *preview.Tax ||
		invoice.TotalExcludingTax != preview.TotalExcludingTax {
		t.Fatalf("discount+tax mismatch preview=%#v invoice=%#v", preview, invoice)
	}
}

// TestInvoicePreviewNoTaxRatesKeepsNilTax: untaxed subscription keeps tax nil / empty arrays.
func TestInvoicePreviewNoTaxRatesKeepsNilTax(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, _ := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {lite.ID},
		"items[0][quantity]": {"1"},
	})
	periodStart := created.CurrentPeriodStart

	preview := postForm[previewTaxFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"subscription":                             {created.ID},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(periodStart, 10)},
	})
	if preview.Subtotal != 5000 || preview.Total != 5000 || preview.AmountDue != 5000 {
		t.Fatalf("no-tax preview amounts = %#v, want 5000", preview)
	}
	if preview.Tax != nil {
		t.Fatalf("tax = %v, want nil", preview.Tax)
	}
	if len(preview.TotalTaxes) != 0 || len(preview.DefaultTaxRates) != 0 {
		t.Fatalf("total_taxes/default_tax_rates = %#v / %#v, want empty", preview.TotalTaxes, preview.DefaultTaxRates)
	}
}

// TestInvoicePreviewPriceIDFormMatchesPrice: [price_id] must equal [price] for proration+tax.
func TestInvoicePreviewPriceIDFormMatchesPrice(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})
	periodStart := created.CurrentPeriodStart

	viaPrice := postForm[previewTaxFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"subscription":                             {created.ID},
		"subscription_details[items][0][price]":    {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(periodStart, 10)},
	})
	viaPriceID := postForm[previewTaxFields](t, handler, "/v1/invoices/create_preview", url.Values{
		"subscription": {created.ID},
		"subscription_details[items][0][price_id]": {pro.ID},
		"subscription_details[proration_behavior]": {"always_invoice"},
		"subscription_details[proration_date]":     {strconv.FormatInt(periodStart, 10)},
	})
	if viaPrice.Subtotal != viaPriceID.Subtotal || viaPrice.Total != viaPriceID.Total ||
		viaPrice.Tax == nil || viaPriceID.Tax == nil || *viaPrice.Tax != *viaPriceID.Tax ||
		viaPrice.TotalExcludingTax != viaPriceID.TotalExcludingTax {
		t.Fatalf("price vs price_id mismatch: price=%#v price_id=%#v", viaPrice, viaPriceID)
	}
	if viaPrice.Subtotal != 5000 || *viaPrice.Tax != 500 || viaPrice.Total != 5500 {
		t.Fatalf("price_id path amounts = %#v, want 5000/500/5500", viaPriceID)
	}
}
