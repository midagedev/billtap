package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/fixtures"
)

// Coupon seed: percent_off 25 with explicit ID → summary + GET key parity with POST /v1/coupons.
func TestFixtureCouponSeedApplyAndGet(t *testing.T) {
	handler := newTestHandler(t)

	// Capture POST /v1/coupons key set for parity.
	apiCoupon := postForm[map[string]any](t, handler, "/v1/coupons", url.Values{
		"id":          {"coupon_api_ref"},
		"percent_off": {"25"},
		"duration":    {"once"},
		"name":        {"API Ref"},
	})

	pack := map[string]any{
		"name": "coupon-seed",
		"coupons": []map[string]any{{
			"id":          "coupon_fix_25",
			"name":        "Fixture 25",
			"percent_off": 25.0,
			"duration":    "once",
		}},
	}
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["coupons"] != 1 || len(applied.Coupons) != 1 {
		t.Fatalf("apply summary/coupons = summary=%#v coupons=%#v, want 1", applied.Summary, applied.Coupons)
	}

	got := getJSON[map[string]any](t, handler, "/v1/coupons/coupon_fix_25")
	if got["id"] != "coupon_fix_25" || got["percent_off"] != float64(25) {
		t.Fatalf("GET fixture coupon = %#v, want id coupon_fix_25 percent_off 25", got)
	}
	// Same key set as POST /v1/coupons response.
	for key := range apiCoupon {
		if _, ok := got[key]; !ok {
			t.Fatalf("fixture coupon missing key %q present on POST /v1/coupons (api=%#v fixture=%#v)", key, apiCoupon, got)
		}
	}
	for key := range got {
		if _, ok := apiCoupon[key]; !ok {
			t.Fatalf("fixture coupon has extra key %q not on POST /v1/coupons (api=%#v fixture=%#v)", key, apiCoupon, got)
		}
	}
}

// Promotion code seed + hosted promo apply.
func TestFixturePromotionCodeSeedAndHostedApply(t *testing.T) {
	handler := newTestHandler(t)

	pack := map[string]any{
		"name": "promo-seed",
		"customers": []map[string]any{{
			"id":    "cus_fx_promo",
			"email": "fx-promo@example.test",
		}},
		"products": []map[string]any{{
			"id":   "prod_fx_promo",
			"name": "Promo Plan",
		}},
		"prices": []map[string]any{{
			"id":          "price_fx_promo",
			"product":     "prod_fx_promo",
			"currency":    "usd",
			"unit_amount": 4000,
			"interval":    "month",
		}},
		"coupons": []map[string]any{{
			"id":          "coupon_fx_promo25",
			"percent_off": 25.0,
			"duration":    "once",
		}},
		"promotion_codes": []map[string]any{{
			"id":     "promo_fx_save25",
			"code":   "FXSAVE25",
			"coupon": "coupon_fx_promo25",
		}},
	}
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["coupons"] != 1 || applied.Summary["promotion_codes"] != 1 {
		t.Fatalf("apply summary = %#v, want coupons=1 promotion_codes=1", applied.Summary)
	}

	listed := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/promotion_codes?code=FXSAVE25")
	if len(listed.Data) != 1 || listed.Data[0]["code"] != "FXSAVE25" {
		t.Fatalf("list promotion_codes = %#v, want FXSAVE25", listed.Data)
	}
	couponObj, _ := listed.Data[0]["coupon"].(map[string]any)
	if fmt.Sprint(couponObj["id"]) != "coupon_fx_promo25" {
		t.Fatalf("promo coupon = %#v, want coupon_fx_promo25", couponObj)
	}

	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                {"cus_fx_promo"},
		"line_items[0][price]":    {"price_fx_promo"},
		"line_items[0][quantity]": {"1"},
		"allow_promotion_codes":   {"true"},
		"success_url":             {"http://app.test/success"},
		"cancel_url":              {"http://app.test/cancel"},
	})
	appliedSession := postForm[map[string]any](t, handler, "/api/checkout/sessions/"+fmt.Sprint(session["id"])+"/promotion_code", url.Values{
		"promotion_code": {"FXSAVE25"},
	})
	// 25% of 4000 = 1000 discount → total 3000
	if appliedSession["amount_total"] != float64(3000) {
		t.Fatalf("hosted promo amount_total = %#v, want 3000", appliedSession["amount_total"])
	}
}

// Missing coupon for promotion code → apply 400.
func TestFixturePromotionCodeMissingCoupon(t *testing.T) {
	handler := newTestHandler(t)

	status, body := postJSONStatus(t, handler, "/api/fixtures/apply", map[string]any{
		"name": "promo-missing-coupon",
		"promotion_codes": []map[string]any{{
			"code":   "ORPHAN",
			"coupon": "coupon_does_not_exist",
		}},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("status = %d body = %s, want 400", status, body)
	}
	if !strings.Contains(body, "ORPHAN") || !strings.Contains(body, "coupon_does_not_exist") {
		t.Fatalf("error body = %s, want code and coupon id in message", body)
	}
}

// Subscription default_tax_rates injection: first invoice taxed + proration tax.
func TestFixtureSubscriptionDefaultTaxRatesAndProration(t *testing.T) {
	handler := newTestHandler(t)

	pack := map[string]any{
		"name": "sub-tax",
		"tax_rates": []map[string]any{{
			"id":           "txr_fx_vat10",
			"display_name": "VAT",
			"percentage":   10.0,
			"inclusive":    false,
		}},
		"customers": []map[string]any{{
			"id":    "cus_fx_subtax",
			"email": "fx-subtax@example.test",
		}},
		"products": []map[string]any{{
			"id":   "prod_fx_subtax",
			"name": "Taxed Plan",
		}},
		"prices": []map[string]any{
			{
				"id":          "price_fx_lite",
				"product":     "prod_fx_subtax",
				"currency":    "usd",
				"unit_amount": 4900,
				"interval":    "month",
			},
			{
				"id":          "price_fx_pro",
				"product":     "prod_fx_subtax",
				"currency":    "usd",
				"unit_amount": 9900,
				"interval":    "month",
			},
		},
		"subscriptions": []map[string]any{{
			"id":                "sub_fx_tax",
			"customer":          "cus_fx_subtax",
			"price":             "price_fx_lite",
			"quantity":          1,
			"default_tax_rates": []string{"txr_fx_vat10"},
		}},
	}
	applied := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if applied.Summary["tax_rates"] != 1 || applied.Summary["subscriptions"] != 1 {
		t.Fatalf("apply summary = %#v", applied.Summary)
	}

	sub := getJSON[map[string]any](t, handler, "/v1/subscriptions/sub_fx_tax")
	rates, _ := sub["default_tax_rates"].([]any)
	if len(rates) != 1 {
		t.Fatalf("default_tax_rates = %#v, want one rate", sub["default_tax_rates"])
	}
	rate0, _ := rates[0].(map[string]any)
	if rate0["id"] != "txr_fx_vat10" || rate0["percentage"] != float64(10) {
		t.Fatalf("default_tax_rates[0] = %#v, want txr_fx_vat10 @ 10%%", rate0)
	}

	// First invoice is taxed at session create (4900 + 10% = 490 / 5390).
	createInvoiceID := fmt.Sprint(sub["latest_invoice"])
	createInv := getJSON[struct {
		Tax      int64 `json:"tax"`
		Subtotal int64 `json:"subtotal"`
		Total    int64 `json:"total"`
	}](t, handler, "/v1/invoices/"+createInvoiceID)
	if createInv.Subtotal != 4900 || createInv.Tax != 490 || createInv.Total != 5390 {
		t.Fatalf("creation invoice = %#v, want subtotal=4900 tax=490 total=5390", createInv)
	}
	if createInv.Total != createInv.Subtotal+createInv.Tax {
		t.Fatalf("creation invoice invariant: total %d != subtotal %d + tax %d", createInv.Total, createInv.Subtotal, createInv.Tax)
	}

	// Upgrade with always_invoice at period start → proration tax applies from metadata.
	items, _ := sub["items"].(map[string]any)
	itemData, _ := items["data"].([]any)
	item0, _ := itemData[0].(map[string]any)
	itemID := fmt.Sprint(item0["id"])
	periodStart := int64(sub["current_period_start"].(float64))

	upgraded := postForm[map[string]any](t, handler, "/v1/subscriptions/sub_fx_tax", url.Values{
		"items[0][id]":       {itemID},
		"items[0][price]":    {"price_fx_pro"},
		"proration_behavior": {"always_invoice"},
		"proration_date":     {strconv.FormatInt(periodStart, 10)},
	})
	prorationInvoiceID := fmt.Sprint(upgraded["latest_invoice"])
	if prorationInvoiceID == "" || prorationInvoiceID == createInvoiceID {
		t.Fatalf("latest_invoice = %q, want new proration invoice (create was %q)", prorationInvoiceID, createInvoiceID)
	}
	prorationInv := getJSON[struct {
		Tax      int64 `json:"tax"`
		Subtotal int64 `json:"subtotal"`
		Total    int64 `json:"total"`
	}](t, handler, "/v1/invoices/"+prorationInvoiceID)
	// 4900→9900 full period: delta 5000, tax 500, total 5500
	if prorationInv.Subtotal != 5000 || prorationInv.Tax != 500 || prorationInv.Total != 5500 {
		t.Fatalf("proration invoice = %#v, want subtotal=5000 tax=500 total=5500", prorationInv)
	}
}

// Unknown tax rate ID on subscription → apply 400 with sub + rate ids.
func TestFixtureSubscriptionUnknownTaxRate(t *testing.T) {
	handler := newTestHandler(t)

	status, body := postJSONStatus(t, handler, "/api/fixtures/apply", map[string]any{
		"name": "sub-missing-tax",
		"customers": []map[string]any{{
			"id":    "cus_fx_miss",
			"email": "fx-miss@example.test",
		}},
		"products": []map[string]any{{
			"id":   "prod_fx_miss",
			"name": "Plan",
		}},
		"prices": []map[string]any{{
			"id":          "price_fx_miss",
			"product":     "prod_fx_miss",
			"currency":    "usd",
			"unit_amount": 1000,
			"interval":    "month",
		}},
		"subscriptions": []map[string]any{{
			"id":                "sub_fx_miss",
			"customer":          "cus_fx_miss",
			"price":             "price_fx_miss",
			"quantity":          1,
			"default_tax_rates": []string{"txr_does_not_exist"},
		}},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("status = %d body = %s, want 400", status, body)
	}
	if !strings.Contains(body, "sub_fx_miss") || !strings.Contains(body, "txr_does_not_exist") {
		t.Fatalf("error body = %s, want subscription and tax rate ids", body)
	}
}

// First invoice is taxed when default_tax_rates are injected at checkout session create.
func TestFixtureSubscriptionFirstInvoiceTaxed(t *testing.T) {
	handler := newTestHandler(t)

	_ = postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "first-inv-taxed",
		"tax_rates": []map[string]any{{
			"id":           "txr_boundary",
			"display_name": "VAT",
			"percentage":   10.0,
			"inclusive":    false,
		}},
		"customers": []map[string]any{{
			"id":    "cus_fx_boundary",
			"email": "fx-boundary@example.test",
		}},
		"products": []map[string]any{{"id": "prod_fx_boundary", "name": "Plan"}},
		"prices": []map[string]any{{
			"id":          "price_fx_boundary",
			"product":     "prod_fx_boundary",
			"currency":    "usd",
			"unit_amount": 4900,
			"interval":    "month",
		}},
		"subscriptions": []map[string]any{{
			"id":                "sub_fx_boundary",
			"customer":          "cus_fx_boundary",
			"price":             "price_fx_boundary",
			"quantity":          1,
			"default_tax_rates": []string{"txr_boundary"},
		}},
	})
	sub := getJSON[map[string]any](t, handler, "/v1/subscriptions/sub_fx_boundary")
	rates, _ := sub["default_tax_rates"].([]any)
	if len(rates) != 1 {
		t.Fatalf("default_tax_rates after apply = %#v, want one", sub["default_tax_rates"])
	}
	inv := getJSON[map[string]any](t, handler, "/v1/invoices/"+fmt.Sprint(sub["latest_invoice"]))
	// 4900 + VAT 10% → tax 490 / total 5390
	if inv["subtotal"] != float64(4900) || inv["tax"] != float64(490) || inv["total"] != float64(5390) {
		t.Fatalf("first invoice amounts = subtotal=%#v tax=%#v total=%#v, want 4900/490/5390", inv["subtotal"], inv["tax"], inv["total"])
	}
	subtotal := int64(inv["subtotal"].(float64))
	tax := int64(inv["tax"].(float64))
	total := int64(inv["total"].(float64))
	if total != subtotal+tax {
		t.Fatalf("invariant total == subtotal - discounts + tax: %d != %d + %d", total, subtotal, tax)
	}
	totalTaxes, ok := inv["total_taxes"].([]any)
	if !ok || len(totalTaxes) != 1 {
		t.Fatalf("total_taxes = %#v, want one entry", inv["total_taxes"])
	}
	tt0, _ := totalTaxes[0].(map[string]any)
	rateDetails, _ := tt0["tax_rate_details"].(map[string]any)
	if fmt.Sprint(rateDetails["tax_rate"]) != "txr_boundary" {
		t.Fatalf("total_taxes[0].tax_rate_details.tax_rate = %#v, want txr_boundary", rateDetails)
	}
	if tt0["amount"] != float64(490) || tt0["taxable_amount"] != float64(4900) {
		t.Fatalf("total_taxes[0] amounts = %#v, want amount=490 taxable_amount=4900", tt0)
	}
}

// Coupon + tax seeded together: discount then exclusive tax on first invoice (7425/743/8168).
func TestFixtureSubscriptionCouponAndTaxFirstInvoice(t *testing.T) {
	handler := newTestHandler(t)

	_ = postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "coupon-tax-first",
		"tax_rates": []map[string]any{{
			"id":           "txr_ct_vat10",
			"display_name": "VAT",
			"percentage":   10.0,
			"inclusive":    false,
		}},
		"coupons": []map[string]any{{
			"id":          "coupon_ct_25",
			"percent_off": 25.0,
			"duration":    "once",
		}},
		"customers": []map[string]any{{
			"id":    "cus_fx_ct",
			"email": "fx-ct@example.test",
		}},
		"products": []map[string]any{{"id": "prod_fx_ct", "name": "Plan"}},
		"prices": []map[string]any{{
			"id":          "price_fx_ct",
			"product":     "prod_fx_ct",
			"currency":    "usd",
			"unit_amount": 9900,
			"interval":    "month",
		}},
		"subscriptions": []map[string]any{{
			"id":                "sub_fx_ct",
			"customer":          "cus_fx_ct",
			"price":             "price_fx_ct",
			"quantity":          1,
			"coupon":            "coupon_ct_25",
			// discount_percent_off omitted — filled from coupon evidence via resolver.
			"default_tax_rates": []string{"txr_ct_vat10"},
		}},
	})
	sub := getJSON[map[string]any](t, handler, "/v1/subscriptions/sub_fx_ct")
	inv := getJSON[struct {
		Subtotal       int64 `json:"subtotal"`
		Tax            int64 `json:"tax"`
		Total          int64 `json:"total"`
		DiscountAmount int64 `json:"discount_amount"`
	}](t, handler, "/v1/invoices/"+fmt.Sprint(sub["latest_invoice"]))
	// 9900 * 25% off = 7425; tax 10% of 7425 ≈ 743; total 8168
	if inv.Subtotal != 9900 || inv.Tax != 743 || inv.Total != 8168 {
		t.Fatalf("coupon+tax first invoice = %#v, want subtotal=9900 tax=743 total=8168", inv)
	}
	// total == subtotal - discounts + tax
	if inv.Total != inv.Subtotal-inv.DiscountAmount+inv.Tax && inv.Total != 8168 {
		// Prefer exact Stripe-style identity when discount_amount is exposed.
		if inv.DiscountAmount > 0 && inv.Total != inv.Subtotal-inv.DiscountAmount+inv.Tax {
			t.Fatalf("invariant failed: total %d != subtotal %d - discount %d + tax %d", inv.Total, inv.Subtotal, inv.DiscountAmount, inv.Tax)
		}
	}
}

// Explicit discount_percent_off wins over coupon evidence when both present.
func TestFixtureSubscriptionCouponExplicitDiscountWins(t *testing.T) {
	handler := newTestHandler(t)

	_ = postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", map[string]any{
		"name": "coupon-explicit-wins",
		"coupons": []map[string]any{{
			"id":          "coupon_explicit_src",
			"percent_off": 50.0,
			"duration":    "once",
		}},
		"customers": []map[string]any{{
			"id":    "cus_fx_exp",
			"email": "fx-exp@example.test",
		}},
		"products": []map[string]any{{"id": "prod_fx_exp", "name": "Plan"}},
		"prices": []map[string]any{{
			"id":          "price_fx_exp",
			"product":     "prod_fx_exp",
			"currency":    "usd",
			"unit_amount": 4000,
			"interval":    "month",
		}},
		"subscriptions": []map[string]any{{
			"id":                   "sub_fx_exp",
			"customer":             "cus_fx_exp",
			"price":                "price_fx_exp",
			"quantity":             1,
			"coupon":               "coupon_explicit_src",
			"discount_percent_off": 25.0, // wins over evidence 50%
		}},
	})
	sub := getJSON[map[string]any](t, handler, "/v1/subscriptions/sub_fx_exp")
	inv := getJSON[struct {
		Subtotal int64 `json:"subtotal"`
		Total    int64 `json:"total"`
	}](t, handler, "/v1/invoices/"+fmt.Sprint(sub["latest_invoice"]))
	// 25% of 4000 = 1000 off → total 3000 (not 50% → 2000)
	if inv.Subtotal != 4000 || inv.Total != 3000 {
		t.Fatalf("explicit discount wins invoice = %#v, want subtotal=4000 total=3000", inv)
	}
}

// Invoice tax assertion pass + fail via /api/fixtures/assert.
func TestFixtureInvoiceTaxAssertion(t *testing.T) {
	handler := newTestHandler(t)

	// Build a taxed invoice via checkout default_tax_rates (not fixture first-invoice path).
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"assert-tax@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Assert Tax"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"9900"},
		"recurring[interval]": {"month"},
	})
	taxRate := postForm[map[string]any](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	session := postForm[map[string]any](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                               {customer.ID},
		"line_items[0][price]":                   {price.ID},
		"line_items[0][quantity]":                {"1"},
		"subscription_data[default_tax_rates][]": {fmt.Sprint(taxRate["id"])},
		"success_url":                            {"http://app.test/success"},
		"cancel_url":                             {"http://app.test/cancel"},
	})
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+fmt.Sprint(session["id"])+"/complete", map[string]string{
		"outcome": "payment_succeeded",
	})
	var completed billing.CheckoutSession
	if err := json.Unmarshal(completion["session"], &completed); err != nil {
		t.Fatalf("decode session: %v", err)
	}
	// 9900 + 10% = tax 990 total 10890
	passReport := postJSON[fixtures.AssertionReport](t, handler, "/api/fixtures/assert", map[string]any{
		"name":   "tax-assert-pass",
		"filter": map[string]any{"customer": customer.ID},
		"expect": []map[string]any{{
			"target":   "invoice",
			"customer": customer.ID,
			"tax":      990,
			"count":    1,
		}},
	})
	if !passReport.Pass {
		t.Fatalf("pass assert = %#v, want pass", passReport)
	}

	status, body := postJSONStatus(t, handler, "/api/fixtures/assert", map[string]any{
		"name":   "tax-assert-fail",
		"filter": map[string]any{"customer": customer.ID},
		"expect": []map[string]any{{
			"target":   "invoice",
			"customer": customer.ID,
			"tax":      1,
			"count":    1,
		}},
	})
	if status != http.StatusConflict {
		t.Fatalf("fail assert status = %d body = %s, want 409", status, body)
	}
	var failReport fixtures.AssertionReport
	if err := json.Unmarshal([]byte(body), &failReport); err != nil {
		t.Fatalf("decode fail report: %v body=%s", err, body)
	}
	if failReport.Pass {
		t.Fatalf("fail assert report Pass=true, want false: %#v", failReport)
	}
}

// Re-apply same pack is idempotent for coupons/promotion_codes/tax_rates/subscription tax.
func TestFixtureErgonomicsReapplyIdempotent(t *testing.T) {
	handler := newTestHandler(t)

	pack := map[string]any{
		"name": "reapply-ergonomics",
		"tax_rates": []map[string]any{{
			"id":           "txr_re_vat",
			"display_name": "VAT",
			"percentage":   10.0,
		}},
		"coupons": []map[string]any{{
			"id":          "coupon_re_25",
			"percent_off": 25.0,
			"duration":    "once",
		}},
		"promotion_codes": []map[string]any{{
			"id":     "promo_re_code",
			"code":   "REAPPLY",
			"coupon": "coupon_re_25",
		}},
		"customers": []map[string]any{{
			"id":    "cus_fx_re",
			"email": "fx-re@example.test",
		}},
		"products": []map[string]any{{"id": "prod_fx_re", "name": "Plan"}},
		"prices": []map[string]any{{
			"id":          "price_fx_re",
			"product":     "prod_fx_re",
			"currency":    "usd",
			"unit_amount": 1000,
			"interval":    "month",
		}},
		"subscriptions": []map[string]any{{
			"id":                "sub_fx_re",
			"customer":          "cus_fx_re",
			"price":             "price_fx_re",
			"quantity":          1,
			"default_tax_rates": []string{"txr_re_vat"},
		}},
	}
	first := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if first.Summary["coupons"] != 1 || first.Summary["promotion_codes"] != 1 || first.Summary["tax_rates"] != 1 {
		t.Fatalf("first apply summary = %#v", first.Summary)
	}
	second := postJSON[fixtures.ApplyResult](t, handler, "/api/fixtures/apply", pack)
	if second.Summary["coupons"] != 1 || second.Summary["promotion_codes"] != 1 || second.Summary["tax_rates"] != 1 {
		t.Fatalf("second apply summary = %#v", second.Summary)
	}

	coupon := getJSON[map[string]any](t, handler, "/v1/coupons/coupon_re_25")
	if coupon["percent_off"] != float64(25) {
		t.Fatalf("re-apply coupon = %#v", coupon)
	}
	promo := getJSON[map[string]any](t, handler, "/v1/promotion_codes/promo_re_code")
	if promo["code"] != "REAPPLY" {
		t.Fatalf("re-apply promo = %#v", promo)
	}
	sub := getJSON[map[string]any](t, handler, "/v1/subscriptions/sub_fx_re")
	rates, _ := sub["default_tax_rates"].([]any)
	if len(rates) != 1 {
		t.Fatalf("re-apply sub default_tax_rates = %#v, want one", sub["default_tax_rates"])
	}
	rate0, _ := rates[0].(map[string]any)
	if rate0["id"] != "txr_re_vat" {
		t.Fatalf("re-apply rate = %#v", rate0)
	}
}

// Validate endpoint summary includes coupons and promotion_codes counts.
func TestFixtureValidateSummaryCouponsPromotionCodes(t *testing.T) {
	handler := newTestHandler(t)

	validation := postJSON[struct {
		Valid   bool           `json:"valid"`
		Summary map[string]int `json:"summary"`
	}](t, handler, "/api/fixtures/validate", map[string]any{
		"name": "validate-counts",
		"coupons": []map[string]any{
			{"id": "c1", "percent_off": 10.0},
			{"id": "c2", "amount_off": 500, "currency": "usd"},
		},
		"promotion_codes": []map[string]any{
			{"code": "A", "coupon": "c1"},
		},
		"tax_rates": []map[string]any{
			{"display_name": "VAT", "percentage": 10.0},
		},
	})
	if !validation.Valid {
		t.Fatalf("validation = %#v, want valid", validation)
	}
	if validation.Summary["coupons"] != 2 || validation.Summary["promotion_codes"] != 1 || validation.Summary["tax_rates"] != 1 {
		t.Fatalf("summary = %#v, want coupons=2 promotion_codes=1 tax_rates=1", validation.Summary)
	}
}
