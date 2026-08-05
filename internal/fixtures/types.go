package fixtures

import (
	"time"

	"github.com/hckim/billtap/internal/billing"
)

const (
	MetadataFixtureName      = "billtap_fixture_name"
	MetadataFixtureRunID     = "billtap_fixture_run_id"
	MetadataFixtureNamespace = "billtap_fixture_namespace"
	MetadataFixtureRef       = "billtap_fixture_ref"
)

type Pack struct {
	Name              string                `json:"name" yaml:"name"`
	RunID             string                `json:"runId" yaml:"runId"`
	Namespace         string                `json:"namespace" yaml:"namespace"`
	Customers         []CustomerFixture     `json:"customers" yaml:"customers"`
	Catalog           CatalogFixture        `json:"catalog" yaml:"catalog"`
	Products          []ProductFixture      `json:"products" yaml:"products"`
	Prices            []PriceFixture        `json:"prices" yaml:"prices"`
	ConnectedAccounts []AccountFixture      `json:"connected_accounts" yaml:"connected_accounts"`
	TestClocks        []TestClockFixture    `json:"test_clocks" yaml:"test_clocks"`
	Subscriptions     []SubscriptionFixture `json:"subscriptions" yaml:"subscriptions"`
	Refunds           []RefundFixture       `json:"refunds" yaml:"refunds"`
	CreditNotes       []CreditNoteFixture   `json:"credit_notes" yaml:"credit_notes"`
	Disputes          []DisputeFixture         `json:"disputes" yaml:"disputes"`
	TaxRates          []TaxRateFixture         `json:"tax_rates" yaml:"tax_rates"`
	Coupons           []CouponFixture          `json:"coupons" yaml:"coupons"`
	PromotionCodes    []PromotionCodeFixture   `json:"promotion_codes" yaml:"promotion_codes"`
	Assertions        []Expectation            `json:"assertions" yaml:"assertions"`
}

type CatalogFixture struct {
	Products []ProductFixture `json:"products" yaml:"products"`
	Prices   []PriceFixture   `json:"prices" yaml:"prices"`
}

type CustomerFixture struct {
	ID                          string                 `json:"id" yaml:"id"`
	Email                       string                 `json:"email" yaml:"email"`
	Name                        string                 `json:"name" yaml:"name"`
	Ref                         string                 `json:"ref" yaml:"ref"`
	TestClock                   string                 `json:"test_clock" yaml:"test_clock"`
	PaymentMethodsFixture       string                 `json:"payment_methods_fixture" yaml:"payment_methods_fixture"`
	PaymentMethodsFixtureCamel  string                 `json:"paymentMethodsFixture" yaml:"paymentMethodsFixture"`
	PaymentMethods              []PaymentMethodFixture `json:"payment_methods" yaml:"payment_methods"`
	PaymentMethodsCamel         []PaymentMethodFixture `json:"paymentMethods" yaml:"paymentMethods"`
	DefaultPaymentIntentOutcome string                 `json:"default_payment_intent_outcome" yaml:"default_payment_intent_outcome"`
	DefaultPIOutcomeCamel       string                 `json:"defaultPaymentIntentOutcome" yaml:"defaultPaymentIntentOutcome"`
	DefaultInvoiceOutcome       string                 `json:"default_invoice_outcome" yaml:"default_invoice_outcome"`
	DefaultInvoiceOutcomeCamel  string                 `json:"defaultInvoiceOutcome" yaml:"defaultInvoiceOutcome"`
	DefaultRenewalOutcome       string                 `json:"default_renewal_outcome" yaml:"default_renewal_outcome"`
	DefaultRenewalOutcomeCamel  string                 `json:"defaultRenewalOutcome" yaml:"defaultRenewalOutcome"`
	Metadata                    map[string]string      `json:"metadata" yaml:"metadata"`
}

type AccountFixture struct {
	ID               string            `json:"id" yaml:"id"`
	Type             string            `json:"type" yaml:"type"`
	Country          string            `json:"country" yaml:"country"`
	Email            string            `json:"email" yaml:"email"`
	BusinessType     string            `json:"business_type" yaml:"business_type"`
	DefaultCurrency  string            `json:"default_currency" yaml:"default_currency"`
	ChargesEnabled   *bool             `json:"charges_enabled" yaml:"charges_enabled"`
	PayoutsEnabled   *bool             `json:"payouts_enabled" yaml:"payouts_enabled"`
	DetailsSubmitted *bool             `json:"details_submitted" yaml:"details_submitted"`
	Capabilities     map[string]string `json:"capabilities" yaml:"capabilities"`
	Metadata         map[string]string `json:"metadata" yaml:"metadata"`
}

type PaymentMethodFixture struct {
	ID      string `json:"id" yaml:"id"`
	Type    string `json:"type" yaml:"type"`
	Brand   string `json:"brand" yaml:"brand"`
	Last4   string `json:"last4" yaml:"last4"`
	Default bool   `json:"default" yaml:"default"`
}

type ProductFixture struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Description string            `json:"description" yaml:"description"`
	Active      *bool             `json:"active" yaml:"active"`
	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
}

type PriceFixture struct {
	ID                 string            `json:"id" yaml:"id"`
	Product            string            `json:"product" yaml:"product"`
	Currency           string            `json:"currency" yaml:"currency"`
	UnitAmount         int64             `json:"unitAmount" yaml:"unitAmount"`
	UnitAmountSnake    int64             `json:"unit_amount" yaml:"unit_amount"`
	LookupKey          string            `json:"lookupKey" yaml:"lookupKey"`
	LookupKeySnake     string            `json:"lookup_key" yaml:"lookup_key"`
	Interval           string            `json:"interval" yaml:"interval"`
	IntervalCount      int               `json:"intervalCount" yaml:"intervalCount"`
	IntervalCountSnake int               `json:"interval_count" yaml:"interval_count"`
	Active             *bool             `json:"active" yaml:"active"`
	Metadata           map[string]string `json:"metadata" yaml:"metadata"`
}

type TestClockFixture struct {
	ID         string `json:"id" yaml:"id"`
	Name       string `json:"name" yaml:"name"`
	FrozenTime string `json:"frozen_time" yaml:"frozen_time"`
}

type SubscriptionFixture struct {
	ID                     string                    `json:"id" yaml:"id"`
	CheckoutSessionID      string                    `json:"checkoutSession" yaml:"checkoutSession"`
	CheckoutSessionIDSnake string                    `json:"checkout_session" yaml:"checkout_session"`
	InvoiceID              string                    `json:"invoice" yaml:"invoice"`
	PaymentIntentID        string                    `json:"paymentIntent" yaml:"paymentIntent"`
	PaymentIntentIDSnake   string                    `json:"payment_intent" yaml:"payment_intent"`
	Ref                    string                    `json:"ref" yaml:"ref"`
	Customer               string                    `json:"customer" yaml:"customer"`
	Price                  string                    `json:"price" yaml:"price"`
	Quantity               int64                     `json:"quantity" yaml:"quantity"`
	Items                  []SubscriptionItemFixture `json:"items" yaml:"items"`
	Outcome                string                    `json:"outcome" yaml:"outcome"`
	Metadata               map[string]string         `json:"metadata" yaml:"metadata"`
	CancelAtPeriodEnd      *bool                     `json:"cancelAtPeriodEnd" yaml:"cancelAtPeriodEnd"`
	CancelAtPeriodEndSnake *bool                     `json:"cancel_at_period_end" yaml:"cancel_at_period_end"`
	Status                 string                    `json:"status" yaml:"status"`
	CurrentPeriodStart     string                    `json:"current_period_start" yaml:"current_period_start"`
	CurrentPeriodEnd       string                    `json:"current_period_end" yaml:"current_period_end"`
	TrialStart             string                    `json:"trial_start" yaml:"trial_start"`
	TrialEnd               string                    `json:"trial_end" yaml:"trial_end"`
	CancelAt               string                    `json:"cancel_at" yaml:"cancel_at"`
	CanceledAt             string                    `json:"canceled_at" yaml:"canceled_at"`
	EndedAt                string                    `json:"ended_at" yaml:"ended_at"`
	LatestInvoiceStatus    string                    `json:"latest_invoice_status" yaml:"latest_invoice_status"`
	TestClock              string                    `json:"test_clock" yaml:"test_clock"`
	RenewalOutcome         string                    `json:"renewal_outcome" yaml:"renewal_outcome"`
	Coupon                 string                    `json:"coupon" yaml:"coupon"`
	PromotionCode          string                    `json:"promotion_code" yaml:"promotion_code"`
	DiscountPercentOff     float64                   `json:"discount_percent_off" yaml:"discount_percent_off"`
	DiscountAmountOff      int64                     `json:"discount_amount_off" yaml:"discount_amount_off"`
	DiscountCurrency       string                    `json:"discount_currency" yaml:"discount_currency"`
	// DefaultTaxRates attaches local tax_rate evidence IDs to the subscription after apply
	// (API-layer hook). Applied to renewals/prorations/previews; the creation invoice is not backfilled.
	DefaultTaxRates []string `json:"default_tax_rates,omitempty" yaml:"default_tax_rates,omitempty"`
}

type SubscriptionItemFixture struct {
	Price    string `json:"price" yaml:"price"`
	Quantity int64  `json:"quantity" yaml:"quantity"`
}

type RefundFixture struct {
	ID            string            `json:"id" yaml:"id"`
	Charge        string            `json:"charge" yaml:"charge"`
	PaymentIntent string            `json:"payment_intent" yaml:"payment_intent"`
	Invoice       string            `json:"invoice" yaml:"invoice"`
	Customer      string            `json:"customer" yaml:"customer"`
	Amount        int64             `json:"amount" yaml:"amount"`
	Currency      string            `json:"currency" yaml:"currency"`
	Reason        string            `json:"reason" yaml:"reason"`
	Status        string            `json:"status" yaml:"status"`
	TestClock     string            `json:"test_clock" yaml:"test_clock"`
	SettleAt      string            `json:"settle_at" yaml:"settle_at"`
	AvailableOn   string            `json:"available_on" yaml:"available_on"`
	Metadata      map[string]string `json:"metadata" yaml:"metadata"`
}

type CreditNoteFixture struct {
	ID       string            `json:"id" yaml:"id"`
	Invoice  string            `json:"invoice" yaml:"invoice"`
	Customer string            `json:"customer" yaml:"customer"`
	Amount   int64             `json:"amount" yaml:"amount"`
	Currency string            `json:"currency" yaml:"currency"`
	Reason   string            `json:"reason" yaml:"reason"`
	Status   string            `json:"status" yaml:"status"`
	Metadata map[string]string `json:"metadata" yaml:"metadata"`
}

type DisputeFixture struct {
	ID       string            `json:"id" yaml:"id"`
	Charge   string            `json:"charge" yaml:"charge"`
	Amount   int64             `json:"amount" yaml:"amount"`
	Currency string            `json:"currency" yaml:"currency"`
	Reason   string            `json:"reason" yaml:"reason"`
	Status   string            `json:"status" yaml:"status"`
	Metadata map[string]string `json:"metadata" yaml:"metadata"`
}

// TaxRateFixture seeds a local tax_rate evidence object (same path as disputes).
// Explicit ID is stored as-is so product/checkout metadata can reference it.
type TaxRateFixture struct {
	ID          string            `json:"id,omitempty" yaml:"id,omitempty"`
	DisplayName string            `json:"display_name" yaml:"display_name"`
	Percentage  float64           `json:"percentage" yaml:"percentage"`
	Inclusive   bool              `json:"inclusive,omitempty" yaml:"inclusive,omitempty"`
	Country     string            `json:"country,omitempty" yaml:"country,omitempty"`
	State       string            `json:"state,omitempty" yaml:"state,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Active      *bool             `json:"active,omitempty" yaml:"active,omitempty"` // default true
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// CouponFixture seeds a local coupon evidence object (same path as tax_rates).
type CouponFixture struct {
	ID               string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	PercentOff       float64           `json:"percent_off,omitempty" yaml:"percent_off,omitempty"`
	AmountOff        int64             `json:"amount_off,omitempty" yaml:"amount_off,omitempty"`
	Currency         string            `json:"currency,omitempty" yaml:"currency,omitempty"`
	Duration         string            `json:"duration,omitempty" yaml:"duration,omitempty"` // default "once"
	DurationInMonths int64             `json:"duration_in_months,omitempty" yaml:"duration_in_months,omitempty"`
	MaxRedemptions   int64             `json:"max_redemptions,omitempty" yaml:"max_redemptions,omitempty"`
	RedeemBy         int64             `json:"redeem_by,omitempty" yaml:"redeem_by,omitempty"`
	AppliesTo        []string          `json:"applies_to_products,omitempty" yaml:"applies_to_products,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// PromotionCodeFixture seeds a local promotion_code evidence object linked to a coupon.
type PromotionCodeFixture struct {
	ID            string            `json:"id,omitempty" yaml:"id,omitempty"`
	Code          string            `json:"code" yaml:"code"`
	Coupon        string            `json:"coupon" yaml:"coupon"`
	Customer      string            `json:"customer,omitempty" yaml:"customer,omitempty"`
	Active        *bool             `json:"active,omitempty" yaml:"active,omitempty"` // default true
	MaxRedemptions int64            `json:"max_redemptions,omitempty" yaml:"max_redemptions,omitempty"`
	ExpiresAt     int64             `json:"expires_at,omitempty" yaml:"expires_at,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

type ApplyResult struct {
	ID                string                    `json:"id"`
	Object            string                    `json:"object"`
	Name              string                    `json:"name"`
	RunID             string                    `json:"runId,omitempty"`
	Namespace         string                    `json:"namespace,omitempty"`
	AppliedAt         time.Time                 `json:"appliedAt"`
	Customers         []billing.Customer        `json:"customers,omitempty"`
	Products          []billing.Product         `json:"products,omitempty"`
	Prices            []billing.Price           `json:"prices,omitempty"`
	ConnectedAccounts []billing.Account         `json:"connectedAccounts,omitempty"`
	CheckoutSessions  []billing.CheckoutSession `json:"checkoutSessions,omitempty"`
	Subscriptions     []billing.Subscription    `json:"subscriptions,omitempty"`
	TestClocks        []billing.TestClock       `json:"testClocks,omitempty"`
	Refunds           []billing.Refund          `json:"refunds,omitempty"`
	CreditNotes       []billing.CreditNote      `json:"creditNotes,omitempty"`
	Disputes          []map[string]any          `json:"disputes,omitempty"`
	TaxRates          []map[string]any          `json:"tax_rates,omitempty"`
	Coupons           []map[string]any          `json:"coupons,omitempty"`
	PromotionCodes    []map[string]any          `json:"promotion_codes,omitempty"`
	Assertions        *AssertionReport          `json:"assertions,omitempty"`
	Summary           map[string]int            `json:"summary"`
}

type ResolveFilter struct {
	Ref         string `json:"ref,omitempty"`
	RunID       string `json:"runId,omitempty"`
	FixtureName string `json:"fixtureName,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	TenantID    string `json:"tenantId,omitempty"`
}

type ResolveResult struct {
	Object            string            `json:"object"`
	Ref               string            `json:"ref"`
	CustomerID        string            `json:"customerId,omitempty"`
	SubscriptionID    string            `json:"subscriptionId,omitempty"`
	InvoiceID         string            `json:"invoiceId,omitempty"`
	PaymentIntentID   string            `json:"paymentIntentId,omitempty"`
	CheckoutSessionID string            `json:"checkoutSessionId,omitempty"`
	PriceID           string            `json:"priceId,omitempty"`
	ProductID         string            `json:"productId,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

type SnapshotFilter struct {
	CustomerID  string `json:"customer,omitempty" yaml:"customer,omitempty"`
	RunID       string `json:"runId,omitempty" yaml:"runId,omitempty"`
	TenantID    string `json:"tenantId,omitempty" yaml:"tenantId,omitempty"`
	FixtureName string `json:"fixtureName,omitempty" yaml:"fixtureName,omitempty"`
	Namespace   string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type Snapshot struct {
	Object           string                    `json:"object"`
	Filter           SnapshotFilter            `json:"filter"`
	Customers        []billing.Customer        `json:"customers"`
	Products         []billing.Product         `json:"products"`
	Prices           []billing.Price           `json:"prices"`
	CheckoutSessions []billing.CheckoutSession `json:"checkoutSessions"`
	Subscriptions    []billing.Subscription    `json:"subscriptions"`
	Invoices         []billing.Invoice         `json:"invoices"`
	PaymentIntents   []billing.PaymentIntent   `json:"paymentIntents"`
	Timeline         []billing.TimelineEntry   `json:"timeline"`
	Summary          map[string]int            `json:"summary"`
	CapturedAt       time.Time                 `json:"capturedAt"`
}

type AssertionRequest struct {
	Name   string         `json:"name" yaml:"name"`
	Filter SnapshotFilter `json:"filter" yaml:"filter"`
	Expect []Expectation  `json:"expect" yaml:"expect"`
}

type Expectation struct {
	Target       string            `json:"target" yaml:"target"`
	ID           string            `json:"id" yaml:"id"`
	Customer     string            `json:"customer" yaml:"customer"`
	Email        string            `json:"email" yaml:"email"`
	Product      string            `json:"product" yaml:"product"`
	Price        string            `json:"price" yaml:"price"`
	LookupKey    string            `json:"lookupKey" yaml:"lookupKey"`
	Status       string            `json:"status" yaml:"status"`
	Metadata     map[string]string `json:"metadata" yaml:"metadata"`
	Exists       *bool             `json:"exists" yaml:"exists"`
	Count        *int              `json:"count" yaml:"count"`
	CountAtLeast *int              `json:"countAtLeast" yaml:"countAtLeast"`
	Quantity     *int64            `json:"quantity" yaml:"quantity"`
	// Total/Tax/Subtotal filter invoice amounts (smallest currency unit). Invoice target only.
	Total    *int64 `json:"total" yaml:"total"`
	Tax      *int64 `json:"tax" yaml:"tax"`
	Subtotal *int64 `json:"subtotal" yaml:"subtotal"`
}

type AssertionReport struct {
	Object    string            `json:"object"`
	Name      string            `json:"name,omitempty"`
	Pass      bool              `json:"pass"`
	Results   []AssertionResult `json:"results"`
	CheckedAt time.Time         `json:"checkedAt"`
}

type AssertionResult struct {
	Target   string         `json:"target"`
	Pass     bool           `json:"pass"`
	Matched  int            `json:"matched"`
	Expected map[string]any `json:"expected"`
	Message  string         `json:"message"`
}
