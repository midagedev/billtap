package fixtures

import (
	"strings"
	"testing"

	"github.com/hckim/billtap/internal/billing"
)

func TestCountInvoicesTaxTotalSubtotalFilters(t *testing.T) {
	invoices := []billing.Invoice{
		{ID: "in_a", CustomerID: "cus_1", Tax: 990, Total: 10890, Subtotal: 9900, Status: "paid"},
		{ID: "in_b", CustomerID: "cus_1", Tax: 0, Total: 4900, Subtotal: 4900, Status: "paid"},
		{ID: "in_c", CustomerID: "cus_2", Tax: 990, Total: 10890, Subtotal: 9900, Status: "paid"},
	}
	tax990 := int64(990)
	tax0 := int64(0)
	total10890 := int64(10890)
	sub9900 := int64(9900)
	count1 := 1

	if got := countInvoices(invoices, Expectation{Customer: "cus_1", Tax: &tax990}); got != 1 {
		t.Fatalf("tax=990 customer cus_1 matched %d, want 1", got)
	}
	if got := countInvoices(invoices, Expectation{Customer: "cus_1", Tax: &tax0}); got != 1 {
		t.Fatalf("tax=0 customer cus_1 matched %d, want 1", got)
	}
	if got := countInvoices(invoices, Expectation{Tax: &tax990, Total: &total10890, Subtotal: &sub9900}); got != 2 {
		t.Fatalf("tax+total+subtotal matched %d, want 2", got)
	}

	result := evaluateExpectation(Snapshot{Invoices: invoices}, Expectation{
		Target:   "invoice",
		Customer: "cus_1",
		Tax:      &tax990,
		Count:    &count1,
	})
	if !result.Pass || result.Matched != 1 {
		t.Fatalf("evaluateExpectation pass case = %#v", result)
	}

	wrong := int64(1)
	fail := evaluateExpectation(Snapshot{Invoices: invoices}, Expectation{
		Target:   "invoice",
		Customer: "cus_1",
		Tax:      &wrong,
		Count:    &count1,
	})
	if fail.Pass {
		t.Fatalf("evaluateExpectation fail case passed: %#v", fail)
	}
}

func TestExpectationTaxRejectedOnNonInvoiceTarget(t *testing.T) {
	tax := int64(100)
	result := evaluateExpectation(Snapshot{}, Expectation{
		Target: "subscription",
		Tax:    &tax,
	})
	if result.Pass {
		t.Fatalf("expected fail for tax on subscription, got %#v", result)
	}
	if !strings.Contains(result.Message, "invoice") {
		t.Fatalf("message = %q, want invoice mention", result.Message)
	}
}

func TestValidatePackCouponsAndPromotionCodes(t *testing.T) {
	// Valid pack.
	if err := validatePack(Pack{
		Coupons: []CouponFixture{
			{ID: "c1", PercentOff: 25},
			{ID: "c2", AmountOff: 500, Currency: "usd"},
		},
		PromotionCodes: []PromotionCodeFixture{
			{Code: "SAVE", Coupon: "c1"},
		},
	}); err != nil {
		t.Fatalf("valid pack error = %v", err)
	}

	// Neither percent nor amount.
	err := validatePack(Pack{Coupons: []CouponFixture{{ID: "c"}}})
	if err == nil || !strings.Contains(err.Error(), "coupons[0]") {
		t.Fatalf("neither discount error = %v", err)
	}
	// Both percent and amount.
	err = validatePack(Pack{Coupons: []CouponFixture{{PercentOff: 10, AmountOff: 100, Currency: "usd"}}})
	if err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("both discounts error = %v", err)
	}
	// percent_off out of range.
	err = validatePack(Pack{Coupons: []CouponFixture{{PercentOff: 150}}})
	if err == nil || !strings.Contains(err.Error(), "percent_off") {
		t.Fatalf("percent range error = %v", err)
	}
	// amount_off without currency.
	err = validatePack(Pack{Coupons: []CouponFixture{{AmountOff: 100}}})
	if err == nil || !strings.Contains(err.Error(), "currency") {
		t.Fatalf("currency required error = %v", err)
	}
	// promotion code missing fields.
	err = validatePack(Pack{PromotionCodes: []PromotionCodeFixture{{Code: "X"}}})
	if err == nil || !strings.Contains(err.Error(), "coupon") {
		t.Fatalf("promo coupon required error = %v", err)
	}
	// tax/total/subtotal only on invoice.
	tax := int64(1)
	err = validatePack(Pack{Assertions: []Expectation{{Target: "subscription", Tax: &tax}}})
	if err == nil || !strings.Contains(err.Error(), "invoice") {
		t.Fatalf("assertion tax target error = %v", err)
	}
}
