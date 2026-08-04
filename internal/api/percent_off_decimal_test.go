package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
)

// TestPercentOffDecimalSupport covers stripe-node v22 number percent_off:
// create/GET round-trip, discount math (exact + round), applies_to, validation,
// tax composition, metadata renewal snapshot, and JSON integer backward-compat.
func TestPercentOffDecimalSupport(t *testing.T) {
	handler := newTestHandler(t)

	// --- 1. percent_off=12.5 create → GET 12.5; integer 50 stays 50 ---
	decimalCoupon := postForm[map[string]any](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_decimal_12_5"},
		"percent_off": {"12.5"},
		"duration":    {"once"},
	})
	if decimalCoupon["percent_off"] != float64(12.5) {
		t.Fatalf("created percent_off = %#v, want 12.5", decimalCoupon["percent_off"])
	}
	fetchedDecimal := getJSON[map[string]any](t, handler, "/v1/coupons/coupon_decimal_12_5")
	if fetchedDecimal["percent_off"] != float64(12.5) {
		t.Fatalf("GET percent_off = %#v, want 12.5", fetchedDecimal["percent_off"])
	}
	intCoupon := postForm[map[string]any](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_int_50"},
		"percent_off": {"50"},
		"duration":    {"once"},
	})
	if intCoupon["percent_off"] != float64(50) {
		t.Fatalf("integer coupon percent_off = %#v, want 50", intCoupon["percent_off"])
	}

	// Storage / JSON backward-compat: integer percent_off unmarshals into float64.
	var stored billing.Discount
	if err := json.Unmarshal([]byte(`{"coupon":"c_legacy","percent_off":50}`), &stored); err != nil {
		t.Fatalf("unmarshal integer percent_off: %v", err)
	}
	if stored.PercentOff != 50 {
		t.Fatalf("legacy integer percent_off = %v, want 50", stored.PercentOff)
	}
	// Metadata snapshot uses FormatFloat; integer 50 round-trips as "50".
	meta := billing.MergeDiscountMetadata(nil, []billing.Discount{{
		CouponID:   "c_meta",
		PercentOff: 12.5,
		Duration:   "once",
	}})
	if meta[billing.MetadataDiscountPercentOff] != "12.5" {
		t.Fatalf("metadata percent_off = %q, want 12.5", meta[billing.MetadataDiscountPercentOff])
	}
	restored := billing.DiscountsFromMetadata(meta)
	if len(restored) != 1 || restored[0].PercentOff != 12.5 {
		t.Fatalf("DiscountsFromMetadata = %#v, want percent_off 12.5", restored)
	}
	// Integer-string metadata (legacy rows) still parses.
	legacyMeta := map[string]string{
		billing.MetadataDiscountCouponID:   "c_legacy_meta",
		billing.MetadataDiscountPercentOff: "50",
	}
	legacyRestored := billing.DiscountsFromMetadata(legacyMeta)
	if len(legacyRestored) != 1 || legacyRestored[0].PercentOff != 50 {
		t.Fatalf("legacy metadata percent_off = %#v, want 50", legacyRestored)
	}

	// --- 5. validation: 0 / 100.5 / "abc" → 400 ---
	for _, tc := range []struct {
		name  string
		value string
	}{
		{"zero", "0"},
		{"above_100", "100.5"},
		{"non_numeric", "abc"},
	} {
		status, body := postFormStatus(t, handler, "/v1/coupons", url.Values{
			"percent_off": {tc.value},
			"duration":    {"once"},
		})
		if status != http.StatusBadRequest {
			t.Fatalf("percent_off=%s status = %d body = %s, want 400", tc.value, status, body)
		}
		errBody := decodeErrorBody(t, body)
		if errBody.Error.Param != "percent_off" {
			t.Fatalf("percent_off=%s error param = %#v, want percent_off", tc.value, errBody.Error)
		}
	}

	// --- 2. 12.5% of subtotal 8000 → discount 1000, total 7000 ---
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email": {"decimal-discount@example.test"},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Decimal Plan"}})
	price8k := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"8000"},
		"recurring[interval]": {"month"},
	})
	sessionExact := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {price8k.ID},
		"line_items[0][quantity]": {"1"},
		"discounts[0][coupon]":    {"coupon_decimal_12_5"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	if sessionExact["amount_subtotal"] != float64(8000) {
		t.Fatalf("exact amount_subtotal = %#v, want 8000", sessionExact["amount_subtotal"])
	}
	if sessionExact["amount_total"] != float64(7000) {
		t.Fatalf("exact amount_total = %#v, want 7000", sessionExact["amount_total"])
	}
	exactDetails, ok := sessionExact["total_details"].(map[string]any)
	if !ok || exactDetails["amount_discount"] != float64(1000) {
		t.Fatalf("exact total_details = %#v, want amount_discount=1000", sessionExact["total_details"])
	}
	completionExact := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+fmt.Sprint(sessionExact["id"])+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var completedExact billing.CheckoutSession
	if err := json.Unmarshal(completionExact["session"], &completedExact); err != nil {
		t.Fatalf("decode exact completed session: %v", err)
	}
	invoiceExact := getJSON[struct {
		Subtotal             int64            `json:"subtotal"`
		Total                int64            `json:"total"`
		TotalDiscountAmounts []map[string]any `json:"total_discount_amounts"`
	}](t, handler, "/v1/invoices/"+completedExact.InvoiceID)
	if invoiceExact.Subtotal != 8000 || invoiceExact.Total != 7000 || len(invoiceExact.TotalDiscountAmounts) != 1 {
		t.Fatalf("exact invoice = %#v, want subtotal=8000 total=7000 one discount", invoiceExact)
	}
	if amount, _ := invoiceExact.TotalDiscountAmounts[0]["amount"].(float64); amount != 1000 {
		t.Fatalf("exact discount amount = %#v, want 1000", invoiceExact.TotalDiscountAmounts[0])
	}

	// --- 3. 10.5% of 7000 → discount 735 (round); renewal same via metadata ---
	roundCoupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_decimal_10_5"},
		"percent_off": {"10.5"},
		"duration":    {"forever"},
	})
	price7k := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"7000"},
		"recurring[interval]": {"month"},
	})
	// Checkout amounts (no clock): Round(7000 * 10.5 / 100) = 735
	sessionRound := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {price7k.ID},
		"line_items[0][quantity]": {"1"},
		"discounts[0][coupon]":    {roundCoupon.ID},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	if sessionRound["amount_subtotal"] != float64(7000) || sessionRound["amount_total"] != float64(6265) {
		t.Fatalf("round session amounts subtotal=%#v total=%#v, want 7000/6265", sessionRound["amount_subtotal"], sessionRound["amount_total"])
	}
	roundDetails, ok := sessionRound["total_details"].(map[string]any)
	if !ok || roundDetails["amount_discount"] != float64(735) {
		t.Fatalf("round total_details = %#v, want amount_discount=735", sessionRound["total_details"])
	}
	// Subscription created on a test clock so a single advance yields one renewal invoice.
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_decimal_round"},
		"frozen_time": {strconv.FormatInt(time.Date(2032, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	roundCustomer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":      {"decimal-round@example.test"},
		"test_clock": {"clock_decimal_round"},
	})
	subRound := postForm[struct {
		ID       string            `json:"id"`
		Metadata map[string]string `json:"metadata"`
	}](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {roundCustomer.ID},
		"items[0][price]":      {price7k.ID},
		"items[0][quantity]":   {"1"},
		"discounts[0][coupon]": {roundCoupon.ID},
		"test_clock":           {"clock_decimal_round"},
	})
	if subRound.Metadata[billing.MetadataDiscountPercentOff] != "10.5" {
		t.Fatalf("subscription metadata percent_off = %q, want 10.5", subRound.Metadata[billing.MetadataDiscountPercentOff])
	}
	advanceRound := postForm[struct {
		BilltapAdvanceResult struct {
			Renewals []struct {
				Invoice struct {
					Subtotal       int64 `json:"subtotal"`
					Total          int64 `json:"total"`
					DiscountAmount int64 `json:"discount_amount"`
				} `json:"invoice"`
			} `json:"renewals"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_decimal_round/advance", url.Values{
		"frozen_time": {strconv.FormatInt(time.Date(2032, 2, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	if len(advanceRound.BilltapAdvanceResult.Renewals) != 1 {
		t.Fatalf("round renewal count = %d, want 1 (first invoice amounts: %#v)", len(advanceRound.BilltapAdvanceResult.Renewals), advanceRound.BilltapAdvanceResult.Renewals)
	}
	renewalRound := advanceRound.BilltapAdvanceResult.Renewals[0].Invoice
	if renewalRound.Subtotal != 7000 || renewalRound.DiscountAmount != 735 || renewalRound.Total != 6265 {
		t.Fatalf("round renewal invoice = %#v, want subtotal=7000 discount=735 total=6265", renewalRound)
	}

	// --- 4. percent_off=12.5 + applies_to product restriction ---
	productMatch := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Decimal Match"}})
	productOther := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Decimal Other"}})
	priceMatch := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {productMatch.ID},
		"currency":            {"usd"},
		"unit_amount":         {"8000"},
		"recurring[interval]": {"month"},
	})
	priceOther := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {productOther.ID},
		"currency":            {"usd"},
		"unit_amount":         {"2000"},
		"recurring[interval]": {"month"},
	})
	scopedCoupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"id":                     {"coupon_decimal_scoped"},
		"percent_off":            {"12.5"},
		"duration":               {"forever"},
		"applies_to[products][]": {productMatch.ID},
	})
	// eligible base 8000 * 12.5% = 1000; subtotal 10000; total 9000
	sessionScoped := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {customer.ID},
		"line_items[0][price]":    {priceMatch.ID},
		"line_items[0][quantity]": {"1"},
		"line_items[1][price]":    {priceOther.ID},
		"line_items[1][quantity]": {"1"},
		"discounts[0][coupon]":    {scopedCoupon.ID},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	if sessionScoped["amount_subtotal"] != float64(10000) {
		t.Fatalf("scoped amount_subtotal = %#v, want 10000", sessionScoped["amount_subtotal"])
	}
	if sessionScoped["amount_total"] != float64(9000) {
		t.Fatalf("scoped amount_total = %#v, want 9000", sessionScoped["amount_total"])
	}
	scopedDetails, ok := sessionScoped["total_details"].(map[string]any)
	if !ok || scopedDetails["amount_discount"] != float64(1000) {
		t.Fatalf("scoped total_details = %#v, want amount_discount=1000", sessionScoped["total_details"])
	}

	// --- 6. tax + 12.5% discount: after discount 7000, 10% tax → 700, total 7700 ---
	taxCustomer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":                 {"decimal-tax@example.test"},
		"metadata[tax_percent]": {"10"},
	})
	taxCoupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_decimal_tax"},
		"percent_off": {"12.5"},
		"duration":    {"once"},
	})
	sessionTax := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {taxCustomer.ID},
		"line_items[0][price]":    {price8k.ID},
		"line_items[0][quantity]": {"1"},
		"discounts[0][coupon]":    {taxCoupon.ID},
		"automatic_tax[enabled]":  {"true"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	// discount 1000, post-discount 7000, tax Round(7000*10/100)=700, total 7700
	if sessionTax["amount_subtotal"] != float64(8000) {
		t.Fatalf("tax amount_subtotal = %#v, want 8000", sessionTax["amount_subtotal"])
	}
	if sessionTax["amount_total"] != float64(7700) {
		t.Fatalf("tax amount_total = %#v, want 7700", sessionTax["amount_total"])
	}
	taxDetails, ok := sessionTax["total_details"].(map[string]any)
	if !ok || taxDetails["amount_discount"] != float64(1000) || taxDetails["amount_tax"] != float64(700) {
		t.Fatalf("tax total_details = %#v, want discount=1000 tax=700", sessionTax["total_details"])
	}
	completionTax := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+fmt.Sprint(sessionTax["id"])+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var completedTax billing.CheckoutSession
	if err := json.Unmarshal(completionTax["session"], &completedTax); err != nil {
		t.Fatalf("decode tax completed session: %v", err)
	}
	invoiceTax := getJSON[map[string]any](t, handler, "/v1/invoices/"+completedTax.InvoiceID)
	if invoiceTax["subtotal"] != float64(8000) || invoiceTax["total"] != float64(7700) || invoiceTax["tax"] != float64(700) {
		t.Fatalf("tax invoice = %#v, want subtotal=8000 total=7700 tax=700", invoiceTax)
	}
}

func TestApplyDiscountsPercentOffRounding(t *testing.T) {
	// Pure-node assertions for the math contract (no HTTP).
	// Exact: 50% of 2000 stays 1000 (integer path invariant).
	total, amount := billing.ApplyDiscountsWithEligibleBase(2000, 2000, "usd", []billing.Discount{{PercentOff: 50}})
	if total != 1000 || amount != 1000 {
		t.Fatalf("50%% of 2000 = total=%d amount=%d, want 1000/1000", total, amount)
	}
	// Exact decimal: 12.5% of 8000 = 1000.
	total, amount = billing.ApplyDiscountsWithEligibleBase(8000, 8000, "usd", []billing.Discount{{PercentOff: 12.5}})
	if total != 7000 || amount != 1000 {
		t.Fatalf("12.5%% of 8000 = total=%d amount=%d, want 7000/1000", total, amount)
	}
	// Round half away from zero: 10.5% of 7000 = 735.
	total, amount = billing.ApplyDiscountsWithEligibleBase(7000, 7000, "usd", []billing.Discount{{PercentOff: 10.5}})
	if total != 6265 || amount != 735 {
		t.Fatalf("10.5%% of 7000 = total=%d amount=%d, want 6265/735", total, amount)
	}
	// Clamp >100 to 100.
	total, amount = billing.ApplyDiscountsWithEligibleBase(1000, 1000, "usd", []billing.Discount{{PercentOff: 150}})
	if total != 0 || amount != 1000 {
		t.Fatalf("150%% clamp of 1000 = total=%d amount=%d, want 0/1000", total, amount)
	}
}
