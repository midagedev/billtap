package billing

import "time"

const (
	ObjectCustomer        = "customer"
	ObjectProduct         = "product"
	ObjectPrice           = "price"
	ObjectCheckoutSession = "checkout.session"
	ObjectSubscription    = "subscription"
	ObjectInvoice         = "invoice"
	ObjectPaymentIntent   = "payment_intent"
	ObjectSetupIntent     = "setup_intent"
	ObjectTestClock       = "test_helpers.test_clock"
	ObjectRefund          = "refund"
	ObjectCreditNote      = "credit_note"
	ObjectTimelineEntry   = "timeline_entry"
	ObjectAccount         = "account"
	ObjectAccountLink     = "account_link"
	ObjectAccountSession  = "account_session"
)

type Customer struct {
	ID        string            `json:"id"`
	Object    string            `json:"object"`
	Email     string            `json:"email,omitempty"`
	Name      string            `json:"name,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type Product struct {
	ID          string            `json:"id"`
	Object      string            `json:"object"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Active      bool              `json:"active"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
}

type Price struct {
	ID                     string            `json:"id"`
	Object                 string            `json:"object"`
	ProductID              string            `json:"product"`
	Currency               string            `json:"currency"`
	UnitAmount             int64             `json:"unit_amount"`
	LookupKey              string            `json:"lookup_key,omitempty"`
	RecurringInterval      string            `json:"recurring_interval,omitempty"`
	RecurringIntervalCount int               `json:"recurring_interval_count,omitempty"`
	Active                 bool              `json:"active"`
	Metadata               map[string]string `json:"metadata,omitempty"`
	CreatedAt              time.Time         `json:"created_at,omitempty"`
}

type LineItem struct {
	PriceID  string `json:"price"`
	Quantity int64  `json:"quantity"`
}

type CheckoutSession struct {
	ID                  string     `json:"id"`
	Object              string     `json:"object"`
	CustomerID          string     `json:"customer"`
	Mode                string     `json:"mode"`
	LineItems           []LineItem `json:"line_items"`
	SuccessURL          string     `json:"success_url,omitempty"`
	CancelURL           string     `json:"cancel_url,omitempty"`
	AllowPromotionCodes bool       `json:"allow_promotion_codes,omitempty"`
	TrialPeriodDays     int64      `json:"trial_period_days,omitempty"`
	URL                 string     `json:"url"`
	Status              string     `json:"status"`
	PaymentStatus       string     `json:"payment_status"`
	SubscriptionID      string     `json:"subscription,omitempty"`
	InvoiceID           string     `json:"invoice,omitempty"`
	PaymentIntentID     string     `json:"payment_intent,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
}

type Subscription struct {
	ID                 string            `json:"id"`
	Object             string            `json:"object"`
	CustomerID         string            `json:"customer"`
	Status             string            `json:"status"`
	Items              []LineItem        `json:"items"`
	CurrentPeriodStart time.Time         `json:"current_period_start"`
	CurrentPeriodEnd   time.Time         `json:"current_period_end"`
	CancelAtPeriodEnd  bool              `json:"cancel_at_period_end"`
	CanceledAt         *time.Time        `json:"canceled_at,omitempty"`
	LatestInvoiceID    string            `json:"latest_invoice,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

type Invoice struct {
	ID                 string     `json:"id"`
	Object             string     `json:"object"`
	CustomerID         string     `json:"customer"`
	SubscriptionID     string     `json:"subscription,omitempty"`
	Status             string     `json:"status"`
	Currency           string     `json:"currency"`
	Subtotal           int64      `json:"subtotal"`
	Total              int64      `json:"total"`
	AmountDue          int64      `json:"amount_due"`
	AmountPaid         int64      `json:"amount_paid"`
	AttemptCount       int        `json:"attempt_count"`
	NextPaymentAttempt *time.Time `json:"next_payment_attempt,omitempty"`
	PaymentIntentID    string     `json:"payment_intent,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

type InvoicePaymentOptions struct {
	PaymentMethodID string
	Outcome         string
	PaidOutOfBand   bool
	At              time.Time
}

type InvoicePaymentResult struct {
	Invoice       Invoice       `json:"invoice"`
	Subscription  Subscription  `json:"subscription,omitempty"`
	PaymentIntent PaymentIntent `json:"payment_intent,omitempty"`
}

type ClockAdvanceResult struct {
	Object         string                 `json:"object"`
	AdvancedTo     time.Time              `json:"advanced_to"`
	TestClockID    string                 `json:"test_clock,omitempty"`
	Activated      []Subscription         `json:"activated,omitempty"`
	Renewals       []InvoicePaymentResult `json:"renewals,omitempty"`
	Canceled       []Subscription         `json:"canceled,omitempty"`
	Skipped        []string               `json:"skipped,omitempty"`
	Processed      int                    `json:"processed"`
	ActivatedCount int                    `json:"activated_count"`
	Renewed        int                    `json:"renewed"`
	CanceledCount  int                    `json:"canceled_count"`
}

type PaymentIntent struct {
	ID              string    `json:"id"`
	Object          string    `json:"object"`
	CustomerID      string    `json:"customer"`
	InvoiceID       string    `json:"invoice,omitempty"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	CaptureMethod   string    `json:"capture_method,omitempty"`
	FailureCode     string    `json:"failure_code,omitempty"`
	DeclineCode     string    `json:"decline_code,omitempty"`
	FailureMessage  string    `json:"failure_message,omitempty"`
	PaymentMethodID string    `json:"payment_method,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type SetupIntent struct {
	ID              string    `json:"id"`
	Object          string    `json:"object"`
	CustomerID      string    `json:"customer,omitempty"`
	Status          string    `json:"status"`
	Usage           string    `json:"usage,omitempty"`
	FailureCode     string    `json:"failure_code,omitempty"`
	DeclineCode     string    `json:"decline_code,omitempty"`
	FailureMessage  string    `json:"failure_message,omitempty"`
	PaymentMethodID string    `json:"payment_method,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type TestClock struct {
	ID         string    `json:"id"`
	Object     string    `json:"object"`
	Name       string    `json:"name,omitempty"`
	Status     string    `json:"status"`
	FrozenTime time.Time `json:"frozen_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Refund struct {
	ID              string            `json:"id"`
	Object          string            `json:"object"`
	ChargeID        string            `json:"charge,omitempty"`
	PaymentIntentID string            `json:"payment_intent,omitempty"`
	InvoiceID       string            `json:"invoice,omitempty"`
	CustomerID      string            `json:"customer,omitempty"`
	Amount          int64             `json:"amount"`
	Currency        string            `json:"currency"`
	Reason          string            `json:"reason,omitempty"`
	Status          string            `json:"status"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

type CreditNote struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	InvoiceID  string            `json:"invoice"`
	CustomerID string            `json:"customer,omitempty"`
	Amount     int64             `json:"amount"`
	Currency   string            `json:"currency"`
	Reason     string            `json:"reason,omitempty"`
	Status     string            `json:"status"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

type Account struct {
	ID               string            `json:"id"`
	Object           string            `json:"object"`
	Type             string            `json:"type"`
	Country          string            `json:"country"`
	Email            string            `json:"email,omitempty"`
	BusinessType     string            `json:"business_type,omitempty"`
	DefaultCurrency  string            `json:"default_currency"`
	ChargesEnabled   bool              `json:"charges_enabled"`
	PayoutsEnabled   bool              `json:"payouts_enabled"`
	DetailsSubmitted bool              `json:"details_submitted"`
	Capabilities     map[string]string `json:"capabilities,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type TimelineEntry struct {
	ID                string            `json:"id"`
	Object            string            `json:"object"`
	Action            string            `json:"action"`
	Message           string            `json:"message"`
	ObjectType        string            `json:"object_type"`
	ObjectID          string            `json:"object_id"`
	CustomerID        string            `json:"customer,omitempty"`
	CheckoutSessionID string            `json:"checkout_session,omitempty"`
	SubscriptionID    string            `json:"subscription,omitempty"`
	InvoiceID         string            `json:"invoice,omitempty"`
	PaymentIntentID   string            `json:"payment_intent,omitempty"`
	Data              map[string]string `json:"data,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
}

type SubscriptionFilter struct {
	CustomerID string
}

type InvoiceFilter struct {
	CustomerID     string
	SubscriptionID string
}

type PaymentIntentFilter struct {
	CustomerID string
	InvoiceIDs []string
}

type RefundFilter struct {
	ChargeID        string
	PaymentIntentID string
	InvoiceID       string
	CustomerID      string
}

type CreditNoteFilter struct {
	InvoiceID  string
	CustomerID string
}

type PortalState struct {
	Object         string             `json:"object"`
	Customer       Customer           `json:"customer"`
	Subscription   *Subscription      `json:"subscription,omitempty"`
	Invoices       []Invoice          `json:"invoices"`
	PaymentIntents []PaymentIntent    `json:"payment_intents"`
	Summary        PortalStateSummary `json:"summary"`
	Timeline       []TimelineEntry    `json:"timeline,omitempty"`
}

type PortalStateSummary struct {
	CustomerID          string     `json:"customer"`
	SubscriptionID      string     `json:"subscription,omitempty"`
	SubscriptionStatus  string     `json:"subscription_status,omitempty"`
	Active              bool       `json:"active"`
	PendingCancellation bool       `json:"pending_cancellation"`
	CancelAtPeriodEnd   bool       `json:"cancel_at_period_end"`
	CurrentPeriodEnd    *time.Time `json:"current_period_end,omitempty"`
	LatestInvoiceID     string     `json:"latest_invoice,omitempty"`
	InvoiceCount        int        `json:"invoice_count"`
	OpenInvoiceCount    int        `json:"open_invoice_count"`
	PaymentIntentCount  int        `json:"payment_intent_count"`
}

type PortalPlanChange struct {
	PlanID   string `json:"plan,omitempty"`
	PriceID  string `json:"price,omitempty"`
	Quantity int64  `json:"quantity,omitempty"`
}

type PortalSeatChange struct {
	Quantity int64 `json:"quantity"`
}

type PortalCancel struct {
	Mode string `json:"mode"`
}

type PaymentMethodSimulation struct {
	ID              string    `json:"id"`
	Object          string    `json:"object"`
	CustomerID      string    `json:"customer"`
	PaymentMethodID string    `json:"payment_method,omitempty"`
	Outcome         string    `json:"outcome"`
	Status          string    `json:"status"`
	FailureCode     string    `json:"failure_code,omitempty"`
	FailureMessage  string    `json:"failure_message,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type List[T any] struct {
	Object string `json:"object"`
	Data   []T    `json:"data"`
}
