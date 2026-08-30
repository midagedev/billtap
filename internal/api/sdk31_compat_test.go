package api

import (
	"net/url"
	"testing"

	"github.com/hckim/billtap/internal/billing"
)

// SDK 31 moved a discount's origin into source.{coupon,promotion_code}; the
// legacy top-level coupon object stays for older SDKs.
func TestDiscountSourceShape(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"sdk31@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"SDK31 Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"2000"},
		"recurring[interval]": {"month"},
	})
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"percent_off": {"10"},
		"duration":    {"forever"},
	})
	subscription := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {price.ID},
		"items[0][quantity]":   {"1"},
		"discounts[0][coupon]": {coupon.ID},
	})

	type discountShape struct {
		Coupon struct {
			ID string `json:"id"`
		} `json:"coupon"`
		Source struct {
			Type   string `json:"type"`
			Coupon struct {
				ID     string `json:"id"`
				Object string `json:"object"`
			} `json:"coupon"`
			PromotionCode struct {
				ID string `json:"id"`
			} `json:"promotion_code"`
		} `json:"source"`
		PromotionCode any `json:"promotion_code"`
	}

	couponDiscount := getJSON[discountShape](t, handler, "/v1/subscriptions/"+subscription.ID+"/discount")
	if couponDiscount.Coupon.ID != coupon.ID {
		t.Fatalf("legacy coupon = %#v, want %s", couponDiscount.Coupon, coupon.ID)
	}
	if couponDiscount.Source.Type != "coupon" || couponDiscount.Source.Coupon.ID != coupon.ID || couponDiscount.Source.Coupon.Object != "coupon" {
		t.Fatalf("coupon source = %#v, want type=coupon with %s", couponDiscount.Source, coupon.ID)
	}

	// A promotion-code discount points the source at the promotion code.
	promo := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/promotion_codes", url.Values{
		"coupon": {coupon.ID},
		"code":   {"SDK31"},
	})
	promoSub := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":                     {customer.ID},
		"items[0][price]":              {price.ID},
		"items[0][quantity]":           {"1"},
		"discounts[0][promotion_code]": {promo.ID},
	})
	promoDiscount := getJSON[discountShape](t, handler, "/v1/subscriptions/"+promoSub.ID+"/discount")
	if promoDiscount.Source.Type != "promotion_code" || promoDiscount.Source.PromotionCode.ID != promo.ID {
		t.Fatalf("promotion_code source = %#v, want type=promotion_code with %s", promoDiscount.Source, promo.ID)
	}
	if promoDiscount.Source.Coupon.ID != coupon.ID {
		t.Fatalf("promotion_code source coupon = %s, want %s", promoDiscount.Source.Coupon.ID, coupon.ID)
	}
}

// transfer_lookup_key moves an in-use lookup key onto the updated price,
// clearing it from the price that held it.
func TestPriceTransferLookupKey(t *testing.T) {
	handler := newTestHandler(t)
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Lookup Plan"}})
	holder := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1000"},
		"lookup_key":          {"plan-standard"},
		"recurring[interval]": {"month"},
	})
	other := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"2000"},
		"recurring[interval]": {"month"},
	})

	transferred := postForm[billing.Price](t, handler, "/v1/prices/"+other.ID, url.Values{
		"lookup_key":          {"plan-standard"},
		"transfer_lookup_key": {"true"},
	})
	if transferred.LookupKey != "plan-standard" {
		t.Fatalf("transferred lookup_key = %q, want plan-standard", transferred.LookupKey)
	}
	afterHolder := getJSON[billing.Price](t, handler, "/v1/prices/"+holder.ID)
	if afterHolder.LookupKey != "" {
		t.Fatalf("previous holder lookup_key = %q, want cleared", afterHolder.LookupKey)
	}

	// Without the transfer flag the update still accepts the key but does not
	// clear other prices (bounded local semantics, documented in COMPATIBILITY.md).
	plain := postForm[billing.Price](t, handler, "/v1/prices/"+holder.ID, url.Values{
		"lookup_key": {"plan-standard"},
	})
	if plain.LookupKey != "plan-standard" {
		t.Fatalf("plain update lookup_key = %q, want plan-standard", plain.LookupKey)
	}
}
