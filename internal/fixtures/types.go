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
	Subscriptions []SubscriptionFixture `json:"subscriptions" yaml:"subscriptions"`
	Assertions    []Expectation         `json:"assertions" yaml:"assertions"`
}

type CatalogFixture struct {
	Products []ProductFixture `json:"products" yaml:"products"`
	Prices   []PriceFixture   `json:"prices" yaml:"prices"`
}

type CustomerFixture struct {
	ID       string            `json:"id" yaml:"id"`
	Email    string            `json:"email" yaml:"email"`
	Name     string            `json:"name" yaml:"name"`
	Metadata map[string]string `json:"metadata" yaml:"metadata"`
}

type ProductFixture struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Description string            `json:"description" yaml:"description"`
	Active      *bool             `json:"active" yaml:"active"`
	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
}

type PriceFixture struct {
	ID            string            `json:"id" yaml:"id"`
	Product       string            `json:"product" yaml:"product"`
	Currency      string            `json:"currency" yaml:"currency"`
	UnitAmount    int64             `json:"unitAmount" yaml:"unitAmount"`
	LookupKey     string            `json:"lookupKey" yaml:"lookupKey"`
	Interval      string            `json:"interval" yaml:"interval"`
	IntervalCount int               `json:"intervalCount" yaml:"intervalCount"`
	Active        *bool             `json:"active" yaml:"active"`
	Metadata      map[string]string `json:"metadata" yaml:"metadata"`
}

type SubscriptionFixture struct {
	ID                string                    `json:"id" yaml:"id"`
	CheckoutSessionID string                    `json:"checkoutSession" yaml:"checkoutSession"`
	InvoiceID         string                    `json:"invoice" yaml:"invoice"`
	PaymentIntentID   string                    `json:"paymentIntent" yaml:"paymentIntent"`
	Ref               string                    `json:"ref" yaml:"ref"`
	Customer          string                    `json:"customer" yaml:"customer"`
	Price             string                    `json:"price" yaml:"price"`
	Quantity          int64                     `json:"quantity" yaml:"quantity"`
	Items             []SubscriptionItemFixture `json:"items" yaml:"items"`
	Outcome           string                    `json:"outcome" yaml:"outcome"`
	Metadata          map[string]string         `json:"metadata" yaml:"metadata"`
	CancelAtPeriodEnd *bool                     `json:"cancelAtPeriodEnd" yaml:"cancelAtPeriodEnd"`
}

type SubscriptionItemFixture struct {
	Price    string `json:"price" yaml:"price"`
	Quantity int64  `json:"quantity" yaml:"quantity"`
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
	Assertions       *AssertionReport          `json:"assertions,omitempty"`
	Summary          map[string]int            `json:"summary"`
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
