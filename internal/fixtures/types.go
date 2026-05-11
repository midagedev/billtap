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
	Name          string                `json:"name" yaml:"name"`
	RunID         string                `json:"runId" yaml:"runId"`
	Namespace     string                `json:"namespace" yaml:"namespace"`
	Customers     []CustomerFixture     `json:"customers" yaml:"customers"`
	Catalog       CatalogFixture        `json:"catalog" yaml:"catalog"`
	Products      []ProductFixture      `json:"products" yaml:"products"`
	Prices        []PriceFixture        `json:"prices" yaml:"prices"`
	TestClocks    []TestClockFixture    `json:"test_clocks" yaml:"test_clocks"`
	Subscriptions []SubscriptionFixture `json:"subscriptions" yaml:"subscriptions"`
	Refunds       []RefundFixture       `json:"refunds" yaml:"refunds"`
	CreditNotes   []CreditNoteFixture   `json:"credit_notes" yaml:"credit_notes"`
	Assertions    []Expectation         `json:"assertions" yaml:"assertions"`
}

type CatalogFixture struct {
	Products []ProductFixture `json:"products" yaml:"products"`
	Prices   []PriceFixture   `json:"prices" yaml:"prices"`
}

type CustomerFixture struct {
	ID                         string                 `json:"id" yaml:"id"`
	Email                      string                 `json:"email" yaml:"email"`
	Name                       string                 `json:"name" yaml:"name"`
	Ref                        string                 `json:"ref" yaml:"ref"`
	TestClock                  string                 `json:"test_clock" yaml:"test_clock"`
	PaymentMethodsFixture      string                 `json:"payment_methods_fixture" yaml:"payment_methods_fixture"`
	PaymentMethodsFixtureCamel string                 `json:"paymentMethodsFixture" yaml:"paymentMethodsFixture"`
	PaymentMethods             []PaymentMethodFixture `json:"payment_methods" yaml:"payment_methods"`
	PaymentMethodsCamel        []PaymentMethodFixture `json:"paymentMethods" yaml:"paymentMethods"`
	Metadata                   map[string]string      `json:"metadata" yaml:"metadata"`
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
	Metadata      map[string]string `json:"metadata" yaml:"metadata"`
}

type CreditNoteFixture struct {
	ID       string            `json:"id" yaml:"id"`
	Invoice  string            `json:"invoice" yaml:"invoice"`
	Customer string            `json:"customer" yaml:"customer"`
	Amount   int64             `json:"amount" yaml:"amount"`
	Currency string            `json:"currency" yaml:"currency"`
	Reason   string            `json:"reason" yaml:"reason"`
	Metadata map[string]string `json:"metadata" yaml:"metadata"`
}

type ApplyResult struct {
	ID               string                    `json:"id"`
	Object           string                    `json:"object"`
	Name             string                    `json:"name"`
	RunID            string                    `json:"runId,omitempty"`
	Namespace        string                    `json:"namespace,omitempty"`
	AppliedAt        time.Time                 `json:"appliedAt"`
	Customers        []billing.Customer        `json:"customers,omitempty"`
	Products         []billing.Product         `json:"products,omitempty"`
	Prices           []billing.Price           `json:"prices,omitempty"`
	CheckoutSessions []billing.CheckoutSession `json:"checkoutSessions,omitempty"`
	Subscriptions    []billing.Subscription    `json:"subscriptions,omitempty"`
	TestClocks       []billing.TestClock       `json:"testClocks,omitempty"`
	Refunds          []billing.Refund          `json:"refunds,omitempty"`
	CreditNotes      []billing.CreditNote      `json:"creditNotes,omitempty"`
	Assertions       *AssertionReport          `json:"assertions,omitempty"`
	Summary          map[string]int            `json:"summary"`
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
