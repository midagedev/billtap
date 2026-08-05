package billing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnsupportedOutcome = errors.New("unsupported checkout outcome")
)

const (
	MetadataDiscountCouponID        = "billtap_discount_coupon"
	MetadataDiscountPromotionCodeID = "billtap_discount_promotion_code"
	MetadataDiscountPercentOff      = "billtap_discount_percent_off"
	MetadataDiscountAmountOff       = "billtap_discount_amount_off"
	MetadataDiscountCurrency        = "billtap_discount_currency"
	MetadataDiscountDuration        = "billtap_discount_duration"
	MetadataDiscountCreated         = "billtap_discount_created"
	MetadataDiscountAppliesTo       = "billtap_discount_applies_to"
	MetadataAutomaticTax            = "billtap_automatic_tax"
	MetadataTaxPercent              = "billtap_tax_percent"
	MetadataDefaultTaxRates         = "billtap_default_tax_rates"
	// Pending proration accumulated by create_prorations; applied on next renewal.
	MetadataPendingProrationAmount = "billtap_pending_proration_amount"
	MetadataPendingProrationAt     = "billtap_pending_proration_at"
	// Invoice metadata keys for serialized billing_reason / credit evidence.
	MetadataBillingReason = "billtap_billing_reason"
	// Pre-discount unused old-cycle credit (anchor=now).
	MetadataProrationCredit = "billtap_proration_credit"
	// Discounted unused old-cycle credit when it differs from MetadataProrationCredit.
	MetadataProrationCreditDiscounted = "billtap_proration_credit_discounted"
)

// ErrPaymentRequired indicates a payment attempt failed under error_if_incomplete
// (HTTP 402 card_error at the API boundary).
var ErrPaymentRequired = errors.New("payment required")

// PaymentFailureError is a payment failure envelope for error_if_incomplete.
type PaymentFailureError struct {
	Code        string
	DeclineCode string
	Message     string
}

func (e *PaymentFailureError) Error() string {
	if e == nil || e.Message == "" {
		return ErrPaymentRequired.Error()
	}
	return e.Message
}

func (e *PaymentFailureError) Unwrap() error { return ErrPaymentRequired }

type Repository interface {
	CreateCustomer(context.Context, Customer) (Customer, error)
	GetCustomer(context.Context, string) (Customer, error)
	ListCustomers(context.Context) ([]Customer, error)
	UpdateCustomer(context.Context, string, Customer) (Customer, error)

	CreateProduct(context.Context, Product) (Product, error)
	GetProduct(context.Context, string) (Product, error)
	ListProducts(context.Context) ([]Product, error)
	UpdateProduct(context.Context, string, Product) (Product, error)

	CreatePrice(context.Context, Price) (Price, error)
	GetPrice(context.Context, string) (Price, error)
	ListPrices(context.Context) ([]Price, error)
	UpdatePrice(context.Context, string, Price) (Price, error)

	CreateCheckoutSession(context.Context, CheckoutSession) (CheckoutSession, error)
	GetCheckoutSession(context.Context, string) (CheckoutSession, error)
	ListCheckoutSessions(context.Context) ([]CheckoutSession, error)
	UpdateCheckoutSessionDiscounts(context.Context, string, []Discount) (CheckoutSession, error)
	RecordCheckoutCompletion(context.Context, CheckoutCompletion) (CheckoutSession, error)

	GetSubscription(context.Context, string) (Subscription, error)
	ListSubscriptions(context.Context) ([]Subscription, error)
	ListSubscriptionsByCustomer(context.Context, string) ([]Subscription, error)
	UpdateSubscription(context.Context, Subscription, []TimelineEntry) (Subscription, error)
	CreateInvoice(context.Context, Invoice, []TimelineEntry) (Invoice, error)
	GetInvoice(context.Context, string) (Invoice, error)
	ListInvoices(context.Context) ([]Invoice, error)
	ListInvoicesFiltered(context.Context, InvoiceFilter) ([]Invoice, error)
	UpdateInvoice(context.Context, Invoice, []TimelineEntry) (Invoice, error)
	CreateInvoiceItem(context.Context, InvoiceItem, Invoice, []TimelineEntry) (InvoiceItem, Invoice, error)
	ListInvoiceItemsFiltered(context.Context, InvoiceItemFilter) ([]InvoiceItem, error)
	FinalizeInvoice(context.Context, Invoice, PaymentIntent, []TimelineEntry) (Invoice, PaymentIntent, error)
	UpdateInvoicePayment(context.Context, Subscription, Invoice, PaymentIntent, []TimelineEntry) (Subscription, Invoice, PaymentIntent, error)
	UpdateManualInvoicePayment(context.Context, Invoice, PaymentIntent, []TimelineEntry) (Invoice, PaymentIntent, error)
	RecordSubscriptionRenewal(context.Context, Subscription, Invoice, PaymentIntent, []TimelineEntry) (Subscription, Invoice, PaymentIntent, error)
	GetPaymentIntent(context.Context, string) (PaymentIntent, error)
	CreatePaymentIntent(context.Context, PaymentIntent) (PaymentIntent, error)
	UpdatePaymentIntent(context.Context, PaymentIntent, []TimelineEntry) (PaymentIntent, error)
	ListPaymentIntents(context.Context) ([]PaymentIntent, error)
	ListPaymentIntentsFiltered(context.Context, PaymentIntentFilter) ([]PaymentIntent, error)
	CreateSetupIntent(context.Context, SetupIntent) (SetupIntent, error)
	GetSetupIntent(context.Context, string) (SetupIntent, error)
	UpdateSetupIntent(context.Context, SetupIntent, []TimelineEntry) (SetupIntent, error)
	ListSetupIntents(context.Context) ([]SetupIntent, error)
	CreateTestClock(context.Context, TestClock) (TestClock, error)
	GetTestClock(context.Context, string) (TestClock, error)
	ListTestClocks(context.Context) ([]TestClock, error)
	UpdateTestClock(context.Context, TestClock) (TestClock, error)
	CreateRefund(context.Context, Refund, []TimelineEntry) (Refund, error)
	GetRefund(context.Context, string) (Refund, error)
	ListRefundsFiltered(context.Context, RefundFilter) ([]Refund, error)
	UpdateRefund(context.Context, Refund, []TimelineEntry) (Refund, error)
	CreateCreditNote(context.Context, CreditNote, []TimelineEntry) (CreditNote, error)
	GetCreditNote(context.Context, string) (CreditNote, error)
	ListCreditNotesFiltered(context.Context, CreditNoteFilter) ([]CreditNote, error)
	UpdateCreditNote(context.Context, CreditNote, []TimelineEntry) (CreditNote, error)
	CreateAccount(context.Context, Account) (Account, error)
	GetAccount(context.Context, string) (Account, error)
	ListAccounts(context.Context) ([]Account, error)
	UpdateAccount(context.Context, string, Account) (Account, error)
	CreateConnectResource(context.Context, ConnectResource) (ConnectResource, error)
	GetConnectResource(context.Context, string, string) (ConnectResource, error)
	ListConnectResources(context.Context, ConnectResourceFilter) ([]ConnectResource, error)
	UpdateConnectResource(context.Context, string, string, ConnectResource) (ConnectResource, error)
	DeleteConnectResource(context.Context, string, string) (ConnectResource, error)

	Timeline(context.Context, TimelineFilter) ([]TimelineEntry, error)
	RecordTimeline(context.Context, TimelineEntry) error
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) CreateCustomer(ctx context.Context, in Customer) (Customer, error) {
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("cus")
	}
	in.Object = ObjectCustomer
	in.CreatedAt = now
	return s.repo.CreateCustomer(ctx, in)
}

func (s *Service) GetCustomer(ctx context.Context, id string) (Customer, error) {
	return s.repo.GetCustomer(ctx, id)
}

func (s *Service) ListCustomers(ctx context.Context) ([]Customer, error) {
	return s.repo.ListCustomers(ctx)
}

func (s *Service) UpdateCustomer(ctx context.Context, id string, in Customer) (Customer, error) {
	return s.repo.UpdateCustomer(ctx, id, in)
}

func (s *Service) CreateProduct(ctx context.Context, in Product) (Product, error) {
	if strings.TrimSpace(in.Name) == "" {
		return Product{}, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("prod")
	}
	in.Object = ObjectProduct
	if !in.Active {
		in.Active = true
	}
	if in.CreatedAt.IsZero() {
		in.CreatedAt = s.now()
	}
	return s.repo.CreateProduct(ctx, in)
}

func (s *Service) GetProduct(ctx context.Context, id string) (Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *Service) ListProducts(ctx context.Context) ([]Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *Service) UpdateProduct(ctx context.Context, id string, in Product) (Product, error) {
	return s.repo.UpdateProduct(ctx, id, in)
}

func (s *Service) CreatePrice(ctx context.Context, in Price) (Price, error) {
	if strings.TrimSpace(in.ProductID) == "" {
		return Price{}, fmt.Errorf("%w: product is required", ErrInvalidInput)
	}
	if strings.TrimSpace(in.Currency) == "" {
		return Price{}, fmt.Errorf("%w: currency is required", ErrInvalidInput)
	}
	if in.UnitAmount < 0 {
		return Price{}, fmt.Errorf("%w: unit_amount must be non-negative", ErrInvalidInput)
	}
	if in.RecurringIntervalCount == 0 {
		in.RecurringIntervalCount = 1
	}
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("price")
	}
	in.Object = ObjectPrice
	if !in.Active {
		in.Active = true
	}
	if in.CreatedAt.IsZero() {
		in.CreatedAt = s.now()
	}
	return s.repo.CreatePrice(ctx, in)
}

func (s *Service) GetPrice(ctx context.Context, id string) (Price, error) {
	return s.repo.GetPrice(ctx, id)
}

func (s *Service) ListPrices(ctx context.Context) ([]Price, error) {
	return s.repo.ListPrices(ctx)
}

func (s *Service) UpdatePrice(ctx context.Context, id string, in Price) (Price, error) {
	return s.repo.UpdatePrice(ctx, id, in)
}

func (s *Service) CreateAccount(ctx context.Context, in Account) (Account, error) {
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("acct")
	}
	if strings.TrimSpace(in.Type) == "" {
		in.Type = "express"
	}
	if strings.TrimSpace(in.Country) == "" {
		in.Country = "US"
	}
	if strings.TrimSpace(in.DefaultCurrency) == "" {
		in.DefaultCurrency = "usd"
	}
	if in.Capabilities == nil {
		in.Capabilities = map[string]string{
			"card_payments": "active",
			"transfers":     "active",
		}
	}
	in.Object = ObjectAccount
	now := s.now()
	if in.CreatedAt.IsZero() {
		in.CreatedAt = now
	}
	if in.UpdatedAt.IsZero() {
		in.UpdatedAt = in.CreatedAt
	}
	return s.repo.CreateAccount(ctx, in)
}

func (s *Service) GetAccount(ctx context.Context, id string) (Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *Service) ListAccounts(ctx context.Context) ([]Account, error) {
	return s.repo.ListAccounts(ctx)
}

func (s *Service) UpdateAccount(ctx context.Context, id string, in Account) (Account, error) {
	return s.repo.UpdateAccount(ctx, id, in)
}

func (s *Service) CreateConnectResource(ctx context.Context, in ConnectResource) (ConnectResource, error) {
	if strings.TrimSpace(in.Object) == "" {
		return ConnectResource{}, fmt.Errorf("%w: object is required", ErrInvalidInput)
	}
	if in.Amount < 0 {
		return ConnectResource{}, fmt.Errorf("%w: amount must be non-negative", ErrInvalidInput)
	}
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id(connectResourcePrefix(in.Object))
	}
	if strings.TrimSpace(in.Currency) != "" {
		in.Currency = strings.ToLower(in.Currency)
	}
	if strings.TrimSpace(in.Country) != "" {
		in.Country = strings.ToUpper(in.Country)
	}
	if strings.TrimSpace(in.Status) == "" {
		in.Status = defaultConnectResourceStatus(in.Object)
	}
	now := s.now()
	if in.CreatedAt.IsZero() {
		in.CreatedAt = now
	}
	if in.UpdatedAt.IsZero() {
		in.UpdatedAt = in.CreatedAt
	}
	return s.repo.CreateConnectResource(ctx, in)
}

func (s *Service) GetConnectResource(ctx context.Context, object string, id string) (ConnectResource, error) {
	return s.repo.GetConnectResource(ctx, object, id)
}

func (s *Service) ListConnectResources(ctx context.Context, filter ConnectResourceFilter) ([]ConnectResource, error) {
	return s.repo.ListConnectResources(ctx, filter)
}

func (s *Service) UpdateConnectResource(ctx context.Context, object string, id string, in ConnectResource) (ConnectResource, error) {
	return s.repo.UpdateConnectResource(ctx, object, id, in)
}

func (s *Service) DeleteConnectResource(ctx context.Context, object string, id string) (ConnectResource, error) {
	return s.repo.DeleteConnectResource(ctx, object, id)
}

func connectResourcePrefix(object string) string {
	switch object {
	case ObjectBankAccount:
		return "ba"
	case ObjectCard:
		return "card"
	case ObjectPerson:
		return "person"
	case ObjectTransfer:
		return "tr"
	case ObjectTransferReversal:
		return "trr"
	case ObjectPayout:
		return "po"
	case ObjectApplicationFee:
		return "fee"
	case ObjectFeeRefund:
		return "fr"
	default:
		return "conn"
	}
}

func defaultConnectResourceStatus(object string) string {
	switch object {
	case ObjectBankAccount:
		return "new"
	case ObjectCard:
		return "active"
	case ObjectTransfer:
		return "paid"
	case ObjectTransferReversal:
		return "succeeded"
	case ObjectPayout:
		return "paid"
	case ObjectFeeRefund:
		return "succeeded"
	default:
		return ""
	}
}

func (s *Service) CreateCheckoutSession(ctx context.Context, in CheckoutSession) (CheckoutSession, error) {
	if strings.TrimSpace(in.CustomerID) == "" {
		return CheckoutSession{}, fmt.Errorf("%w: customer is required", ErrInvalidInput)
	}
	if len(in.LineItems) == 0 {
		return CheckoutSession{}, fmt.Errorf("%w: at least one line item is required", ErrInvalidInput)
	}
	for idx, item := range in.LineItems {
		if strings.TrimSpace(item.PriceID) == "" {
			return CheckoutSession{}, fmt.Errorf("%w: line_items[%d].price is required", ErrInvalidInput, idx)
		}
		if item.Quantity <= 0 {
			in.LineItems[idx].Quantity = 1
		}
	}
	if in.Mode == "" {
		in.Mode = "subscription"
	}
	now := s.now()
	in.Discounts = normalizeDiscounts(in.Discounts, now)
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("cs")
	}
	in.Object = ObjectCheckoutSession
	in.URL = "/checkout/" + in.ID
	in.Status = "open"
	in.PaymentStatus = "unpaid"
	in.CreatedAt = now
	return s.repo.CreateCheckoutSession(ctx, in)
}

func (s *Service) GetCheckoutSession(ctx context.Context, id string) (CheckoutSession, error) {
	return s.repo.GetCheckoutSession(ctx, id)
}

func (s *Service) ListCheckoutSessions(ctx context.Context) ([]CheckoutSession, error) {
	return s.repo.ListCheckoutSessions(ctx)
}

// UpdateCheckoutSessionDiscounts replaces open-session discounts (promotion-code apply/remove).
// times_redeemed is not touched here — coupon create path records 0 and never increments on apply/complete.
func (s *Service) UpdateCheckoutSessionDiscounts(ctx context.Context, id string, discounts []Discount) (CheckoutSession, error) {
	discounts = normalizeDiscounts(discounts, s.now())
	return s.repo.UpdateCheckoutSessionDiscounts(ctx, id, discounts)
}

func (s *Service) CompleteCheckout(ctx context.Context, sessionID string, outcome string) (CheckoutSession, error) {
	return s.completeCheckout(ctx, sessionID, outcome, CheckoutCompletionOptions{})
}

func (s *Service) CompleteCheckoutAt(ctx context.Context, sessionID string, outcome string, at time.Time) (CheckoutSession, error) {
	return s.completeCheckout(ctx, sessionID, outcome, CheckoutCompletionOptions{At: at})
}

func (s *Service) CompleteCheckoutWithOptions(ctx context.Context, sessionID string, outcome string, opts CheckoutCompletionOptions) (CheckoutSession, error) {
	return s.completeCheckout(ctx, sessionID, outcome, opts)
}

func (s *Service) completeCheckout(ctx context.Context, sessionID string, outcome string, opts CheckoutCompletionOptions) (CheckoutSession, error) {
	session, err := s.repo.GetCheckoutSession(ctx, sessionID)
	if err != nil {
		return CheckoutSession{}, err
	}
	// Idempotent: already completed (including free payment with no PI).
	if session.Status != "open" || session.PaymentIntentID != "" {
		return session, nil
	}

	outcomeSpec, ok := checkoutOutcomeFor(outcome)
	if !ok {
		return CheckoutSession{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
	}

	subtotal := int64(0)
	currency := "usd"
	lineAmounts := make([]LineAmount, 0, len(session.LineItems))
	for _, item := range session.LineItems {
		price, err := s.repo.GetPrice(ctx, item.PriceID)
		if err != nil {
			return CheckoutSession{}, err
		}
		if price.Currency != "" {
			currency = price.Currency
		}
		quantity := item.Quantity
		if quantity <= 0 {
			quantity = 1
		}
		amount := price.UnitAmount * quantity
		subtotal += amount
		lineAmounts = append(lineAmounts, LineAmount{ProductID: price.ProductID, Amount: amount})
	}

	now := s.now()
	if !opts.At.IsZero() {
		now = opts.At
	}
	discounts := normalizeDiscounts(session.Discounts, now)
	eligibleBase := EligibleDiscountBase(subtotal, discounts, lineAmounts)
	discountedTotal, discountAmount := ApplyDiscountsWithEligibleBase(subtotal, eligibleBase, currency, discounts)

	if session.Mode == "payment" {
		return s.completePaymentCheckout(ctx, session, outcomeSpec, opts, currency, subtotal, discountedTotal, discountAmount, now)
	}

	periodEnd := now.AddDate(0, 1, 0)
	paid := outcomeSpec.Paid
	trialing := paid && session.TrialPeriodDays > 0
	if trialing {
		periodEnd = now.AddDate(0, 0, int(session.TrialPeriodDays))
	}
	invoiceSubtotal := subtotal
	invoiceTotal := discountedTotal
	invoiceDiscountAmount := discountAmount
	invoiceTax := int64(0)
	invoiceDefaultTaxRates := []AppliedTaxRate(nil)
	if !trialing {
		if len(session.DefaultTaxRates) > 0 {
			_, _, exclusiveTotal, taxTotal := ComputeTaxRateAmounts(discountedTotal, session.DefaultTaxRates)
			invoiceTax = taxTotal
			invoiceTotal = discountedTotal + exclusiveTotal
			invoiceDefaultTaxRates = session.DefaultTaxRates
		} else if session.AutomaticTax {
			invoiceTax = ExclusiveTaxAmount(discountedTotal, session.TaxPercent)
			invoiceTotal = discountedTotal + invoiceTax
		}
	}
	if trialing {
		invoiceSubtotal = 0
		invoiceTotal = 0
		invoiceDiscountAmount = 0
		invoiceTax = 0
		invoiceDefaultTaxRates = nil
	}
	invoiceAttemptCount := 1
	if outcomeSpec.InvoiceAttemptCount != nil {
		invoiceAttemptCount = *outcomeSpec.InvoiceAttemptCount
	}

	sub := Subscription{
		ID:                 firstNonEmpty(opts.SubscriptionID, id("sub")),
		Object:             ObjectSubscription,
		CustomerID:         session.CustomerID,
		Status:             "active",
		Items:              session.LineItems,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   periodEnd,
		Metadata:           map[string]string{"checkout_session": session.ID},
	}
	sub.Metadata = MergeDiscountMetadata(sub.Metadata, discounts)
	sub.Metadata = MergeTaxMetadata(sub.Metadata, session.AutomaticTax, session.TaxPercent)
	sub.Metadata = MergeDefaultTaxRatesMetadata(sub.Metadata, session.DefaultTaxRates)
	if trialing {
		sub.Status = "trialing"
		sub.Metadata["trial_period_days"] = fmt.Sprintf("%d", session.TrialPeriodDays)
		sub.Metadata["trial_start"] = now.Format(time.RFC3339Nano)
		sub.Metadata["trial_end"] = periodEnd.Format(time.RFC3339Nano)
	}
	invoice := Invoice{
		ID:              firstNonEmpty(opts.InvoiceID, id("in")),
		Object:          ObjectInvoice,
		CustomerID:      session.CustomerID,
		SubscriptionID:  sub.ID,
		Status:          "paid",
		Currency:        currency,
		Subtotal:        invoiceSubtotal,
		DiscountAmount:  invoiceDiscountAmount,
		Discounts:       discounts,
		AutomaticTax:    session.AutomaticTax,
		DefaultTaxRates: invoiceDefaultTaxRates,
		Tax:             invoiceTax,
		Total:           invoiceTotal,
		AmountDue:       0,
		AmountPaid:      invoiceTotal,
		AttemptCount:    1,
		CreatedAt:       now,
	}
	intent := PaymentIntent{
		ID:              firstNonEmpty(opts.PaymentIntentID, id("pi")),
		Object:          ObjectPaymentIntent,
		CustomerID:      session.CustomerID,
		InvoiceID:       invoice.ID,
		Amount:          invoiceTotal,
		Currency:        currency,
		Status:          outcomeSpec.PaymentIntentStatus,
		PaymentMethodID: outcomeSpec.PaymentMethodID,
		CreatedAt:       now,
	}
	if !paid {
		sub.Status = firstNonEmpty(outcomeSpec.SubscriptionStatus, "incomplete")
		invoice.Status = firstNonEmpty(outcomeSpec.InvoiceStatus, "open")
		invoice.AmountDue = invoiceTotal
		invoice.AmountPaid = 0
		invoice.AttemptCount = invoiceAttemptCount
		if invoice.Status == "void" {
			invoice.AmountDue = 0
		}
		if outcomeSpec.NextPaymentAttempt {
			nextAttempt := now.Add(24 * time.Hour)
			invoice.NextPaymentAttempt = &nextAttempt
		}
		intent.FailureCode = outcomeSpec.FailureCode
		intent.DeclineCode = outcomeSpec.DeclineCode
		intent.FailureMessage = outcomeSpec.FailureMessage
	}
	sub.LatestInvoiceID = invoice.ID
	invoice.PaymentIntentID = intent.ID

	return s.repo.RecordCheckoutCompletion(ctx, CheckoutCompletion{
		SessionID:     session.ID,
		SessionStatus: firstNonEmpty(outcomeSpec.SessionStatus, "complete"),
		PaymentStatus: outcomeSpec.PaymentStatus,
		CheckoutEvent: firstNonEmpty(outcomeSpec.CheckoutEvent, "checkout.session.completed"),
		Outcome:       outcomeSpec.Outcome,
		CompletedAt:   now,
		Subscription:  sub,
		Invoice:       invoice,
		PaymentIntent: intent,
	})
}

// completePaymentCheckout finishes a mode=payment checkout: no subscription/invoice;
// one PaymentIntent for the discounted (and taxed) total. Free totals skip the PI.
func (s *Service) completePaymentCheckout(
	ctx context.Context,
	session CheckoutSession,
	outcomeSpec checkoutOutcomeSpec,
	opts CheckoutCompletionOptions,
	currency string,
	subtotal int64,
	discountedTotal int64,
	_ int64,
	now time.Time,
) (CheckoutSession, error) {
	total := discountedTotal
	if len(session.DefaultTaxRates) > 0 {
		_, _, exclusiveTotal, _ := ComputeTaxRateAmounts(discountedTotal, session.DefaultTaxRates)
		total = discountedTotal + exclusiveTotal
	} else if session.AutomaticTax {
		total = discountedTotal + ExclusiveTaxAmount(discountedTotal, session.TaxPercent)
	}
	_ = subtotal

	paymentStatus := outcomeSpec.PaymentStatus
	sessionStatus := firstNonEmpty(outcomeSpec.SessionStatus, "complete")
	paid := outcomeSpec.Paid

	// Free amount: no PaymentIntent; Stripe reports no_payment_required.
	if total == 0 && paid {
		return s.repo.RecordCheckoutCompletion(ctx, CheckoutCompletion{
			SessionID:     session.ID,
			SessionStatus: sessionStatus,
			PaymentStatus: "no_payment_required",
			CheckoutEvent: firstNonEmpty(outcomeSpec.CheckoutEvent, "checkout.session.completed"),
			Outcome:       outcomeSpec.Outcome,
			CompletedAt:   now,
		})
	}

	if paymentStatus == "" {
		if paid {
			paymentStatus = "paid"
		} else {
			paymentStatus = "unpaid"
		}
	}

	// User payment_intent_data[metadata] only at the unprefixed surface; house
	// fields use billtap_* so they cannot collide with customer keys like
	// metadata[description]. Top-level PI fields are projected in stripePaymentIntent.
	meta := copyMap(session.PaymentIntentMetadata)
	if meta == nil {
		meta = map[string]string{}
	}
	meta["billtap_checkout_session"] = session.ID
	if session.SetupFutureUsage != "" {
		meta["billtap_setup_future_usage"] = session.SetupFutureUsage
	}
	if session.PaymentIntentDescription != "" {
		meta["billtap_description"] = session.PaymentIntentDescription
	}
	if session.ReceiptEmail != "" {
		meta["billtap_receipt_email"] = session.ReceiptEmail
	}

	intent := PaymentIntent{
		ID:              firstNonEmpty(opts.PaymentIntentID, id("pi")),
		Object:          ObjectPaymentIntent,
		CustomerID:      session.CustomerID,
		Amount:          total,
		Currency:        currency,
		Status:          outcomeSpec.PaymentIntentStatus,
		CaptureMethod:   firstNonEmpty(session.CaptureMethod, "automatic"),
		PaymentMethodID: outcomeSpec.PaymentMethodID,
		Metadata:        meta,
		CreatedAt:       now,
	}
	if !paid {
		intent.FailureCode = outcomeSpec.FailureCode
		intent.DeclineCode = outcomeSpec.DeclineCode
		intent.FailureMessage = outcomeSpec.FailureMessage
	}

	return s.repo.RecordCheckoutCompletion(ctx, CheckoutCompletion{
		SessionID:     session.ID,
		SessionStatus: sessionStatus,
		PaymentStatus: paymentStatus,
		CheckoutEvent: firstNonEmpty(outcomeSpec.CheckoutEvent, "checkout.session.completed"),
		Outcome:       outcomeSpec.Outcome,
		CompletedAt:   now,
		PaymentIntent: intent,
	})
}

func (s *Service) GetSubscription(ctx context.Context, id string) (Subscription, error) {
	return s.repo.GetSubscription(ctx, id)
}

func (s *Service) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	return s.repo.ListSubscriptions(ctx)
}

type SubscriptionPatch struct {
	Items              []LineItem
	ReplaceItems       bool
	Metadata           map[string]string
	CancelAtPeriodEnd  *bool
	Status             *string
	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	CanceledAt         *time.Time
	ClearCanceledAt    bool
	LatestInvoiceID    *string
	TimelineSource     string
	TimelineAction     string
	TimelineMessage    string
}

func (s *Service) PatchSubscription(ctx context.Context, subscriptionID string, patch SubscriptionPatch) (Subscription, error) {
	sub, err := s.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return Subscription{}, err
	}
	if patch.ReplaceItems {
		if len(patch.Items) == 0 {
			return Subscription{}, fmt.Errorf("%w: subscription items cannot be empty", ErrInvalidInput)
		}
		sub.Items = patch.Items
	}
	if patch.Metadata != nil {
		sub.Metadata = copyMap(sub.Metadata)
		for key, value := range patch.Metadata {
			if value == "" {
				delete(sub.Metadata, key)
			} else {
				sub.Metadata[key] = value
			}
		}
	}
	if patch.Status != nil {
		sub.Status = strings.ToLower(strings.TrimSpace(*patch.Status))
	}
	if patch.CurrentPeriodStart != nil {
		sub.CurrentPeriodStart = *patch.CurrentPeriodStart
	}
	if patch.CurrentPeriodEnd != nil {
		sub.CurrentPeriodEnd = *patch.CurrentPeriodEnd
	}
	if patch.LatestInvoiceID != nil {
		sub.LatestInvoiceID = strings.TrimSpace(*patch.LatestInvoiceID)
	}
	if patch.ClearCanceledAt {
		sub.CanceledAt = nil
	} else if patch.CanceledAt != nil {
		sub.CanceledAt = patch.CanceledAt
	}
	if patch.CancelAtPeriodEnd != nil {
		sub.Metadata = copyMap(sub.Metadata)
		sub.CancelAtPeriodEnd = *patch.CancelAtPeriodEnd
		if *patch.CancelAtPeriodEnd {
			if sub.CanceledAt == nil {
				canceledAt := s.now()
				sub.CanceledAt = &canceledAt
			}
			sub.Metadata["cancel_at"] = sub.CurrentPeriodEnd.Format(time.RFC3339Nano)
		} else {
			sub.CanceledAt = nil
			delete(sub.Metadata, "cancel_at")
			delete(sub.Metadata, "cancellation_details_comment")
			delete(sub.Metadata, "cancellation_details_feedback")
			if sub.Status == "canceled" {
				sub.Status = "active"
			}
		}
	}
	now := s.now()
	sub.Metadata = copyMap(sub.Metadata)
	sub.Metadata["stripe_compat_updated_at"] = now.Format(time.RFC3339Nano)
	action := firstNonEmpty(patch.TimelineAction, "customer.subscription.updated")
	message := firstNonEmpty(patch.TimelineMessage, "Stripe-compatible subscription updated")
	source := firstNonEmpty(patch.TimelineSource, "stripe_compat")
	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{portalTimeline(
		"stripe_compat_update_"+sub.ID+"_"+now.Format(time.RFC3339Nano),
		action,
		message,
		sub,
		map[string]string{"source": source, "status": sub.Status},
		now,
	)})
}

func (s *Service) GetInvoice(ctx context.Context, id string) (Invoice, error) {
	return s.repo.GetInvoice(ctx, id)
}

func (s *Service) CreateInvoice(ctx context.Context, in Invoice) (Invoice, error) {
	customerID := strings.TrimSpace(in.CustomerID)
	if customerID == "" {
		return Invoice{}, fmt.Errorf("%w: customer is required", ErrInvalidInput)
	}
	customer, err := s.repo.GetCustomer(ctx, customerID)
	if err != nil {
		return Invoice{}, err
	}
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("in")
	}
	in.Object = ObjectInvoice
	in.CustomerID = customer.ID
	in.SubscriptionID = strings.TrimSpace(in.SubscriptionID)
	in.Status = firstNonEmpty(strings.ToLower(strings.TrimSpace(in.Status)), "draft")
	in.Currency = strings.ToLower(firstNonEmpty(strings.TrimSpace(in.Currency), "usd"))
	in.Metadata = copyMap(in.Metadata)
	in.Subtotal = 0
	in.DiscountAmount = 0
	in.Total = 0
	in.AmountDue = 0
	in.AmountPaid = 0
	in.AttemptCount = 0
	in.PaymentIntentID = ""
	if in.CreatedAt.IsZero() {
		in.CreatedAt = now
	}
	return s.repo.CreateInvoice(ctx, in, []TimelineEntry{billingTimelineEntry(
		"invoice_created_"+in.ID,
		"invoice.created",
		"Invoice created",
		ObjectInvoice,
		in.ID,
		in.CustomerID,
		"",
		in.SubscriptionID,
		in.ID,
		"",
		map[string]string{"source": "invoice.create", "status": in.Status},
		in.CreatedAt,
	)})
}

func (s *Service) ListInvoices(ctx context.Context) ([]Invoice, error) {
	return s.repo.ListInvoices(ctx)
}

func (s *Service) CreateInvoiceItem(ctx context.Context, in InvoiceItem) (InvoiceItem, Invoice, error) {
	if strings.TrimSpace(in.InvoiceID) == "" {
		return InvoiceItem{}, Invoice{}, fmt.Errorf("%w: invoice is required", ErrInvalidInput)
	}
	if in.Amount == 0 {
		return InvoiceItem{}, Invoice{}, fmt.Errorf("%w: amount is required", ErrInvalidInput)
	}
	invoice, err := s.repo.GetInvoice(ctx, in.InvoiceID)
	if err != nil {
		return InvoiceItem{}, Invoice{}, err
	}
	customerID := firstNonEmpty(strings.TrimSpace(in.CustomerID), invoice.CustomerID)
	if customerID != invoice.CustomerID {
		return InvoiceItem{}, Invoice{}, fmt.Errorf("%w: customer must match invoice customer", ErrInvalidInput)
	}
	if _, err := s.repo.GetCustomer(ctx, customerID); err != nil {
		return InvoiceItem{}, Invoice{}, err
	}
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("ii")
	}
	in.Object = ObjectInvoiceItem
	in.CustomerID = customerID
	in.Currency = strings.ToLower(firstNonEmpty(strings.TrimSpace(in.Currency), invoice.Currency, "usd"))
	in.Metadata = copyMap(in.Metadata)
	if in.CreatedAt.IsZero() {
		in.CreatedAt = now
	}
	invoice.Subtotal += in.Amount
	invoice.Total += in.Amount
	if invoice.Total < 0 {
		invoice.Total = 0
	}
	invoice.AmountDue = invoice.Total - invoice.AmountPaid
	if invoice.AmountDue < 0 {
		invoice.AmountDue = 0
	}
	invoice.Currency = firstNonEmpty(invoice.Currency, in.Currency)
	createdItem, updatedInvoice, err := s.repo.CreateInvoiceItem(ctx, in, invoice, []TimelineEntry{billingTimelineEntry(
		"invoiceitem_created_"+in.ID,
		"invoiceitem.created",
		"Invoice item created",
		ObjectInvoiceItem,
		in.ID,
		in.CustomerID,
		"",
		invoice.SubscriptionID,
		invoice.ID,
		invoice.PaymentIntentID,
		map[string]string{"source": "invoiceitem.create", "amount": strconv.FormatInt(in.Amount, 10), "currency": in.Currency},
		in.CreatedAt,
	)})
	return createdItem, updatedInvoice, err
}

func (s *Service) ListInvoiceItems(ctx context.Context, filter InvoiceItemFilter) ([]InvoiceItem, error) {
	return s.repo.ListInvoiceItemsFiltered(ctx, filter)
}

func (s *Service) FinalizeInvoice(ctx context.Context, invoiceID string) (InvoicePaymentResult, error) {
	if strings.TrimSpace(invoiceID) == "" {
		return InvoicePaymentResult{}, fmt.Errorf("%w: invoice is required", ErrInvalidInput)
	}
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	if invoice.Status != "draft" {
		var intent PaymentIntent
		if invoice.PaymentIntentID != "" {
			intent, _ = s.repo.GetPaymentIntent(ctx, invoice.PaymentIntentID)
		}
		return InvoicePaymentResult{Invoice: invoice, PaymentIntent: intent}, nil
	}
	at := s.now()
	invoice.Status = "open"
	invoice.AmountDue = invoice.Total
	invoice.AmountPaid = 0
	invoice.NextPaymentAttempt = nil
	invoice.AttemptCount = 0
	intent := PaymentIntent{
		ID:            id("pi"),
		Object:        ObjectPaymentIntent,
		CustomerID:    invoice.CustomerID,
		InvoiceID:     invoice.ID,
		Amount:        invoice.Total,
		Currency:      invoice.Currency,
		Status:        "requires_payment_method",
		CaptureMethod: "automatic",
		Metadata:      copyMap(invoice.Metadata),
		CreatedAt:     at,
	}
	if paymentIntentConfiguredOutcome(intent.Metadata) == "" && invoice.CustomerID != "" {
		if customer, err := s.repo.GetCustomer(ctx, invoice.CustomerID); err == nil {
			if outcome := CustomerDefaultPaymentIntentOutcome(customer.Metadata); outcome != "" {
				if intent.Metadata == nil {
					intent.Metadata = map[string]string{}
				}
				intent.Metadata[MetadataPaymentIntentOutcome] = outcome
			}
		}
	}
	if outcome := paymentIntentConfiguredOutcome(intent.Metadata); outcome != "" && !IsSupportedPaymentIntentOutcome(outcome) {
		return InvoicePaymentResult{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
	}
	invoice.PaymentIntentID = intent.ID
	invoice.Metadata = copyMap(invoice.Metadata)
	invoice.Metadata["billtap_finalized_at"] = at.Format(time.RFC3339Nano)
	updatedInvoice, createdIntent, err := s.repo.FinalizeInvoice(ctx, invoice, intent, []TimelineEntry{
		billingTimelineEntry("invoice_finalized_"+invoice.ID, "invoice.finalized", "Invoice finalized", ObjectInvoice, invoice.ID, invoice.CustomerID, "", invoice.SubscriptionID, invoice.ID, intent.ID, map[string]string{"source": "invoice.finalize", "status": invoice.Status}, at),
		billingTimelineEntry("invoice_payment_intent_created_"+intent.ID, "payment_intent.created", "Invoice payment intent created", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", invoice.SubscriptionID, invoice.ID, intent.ID, map[string]string{"source": "invoice.finalize", "status": intent.Status}, at),
	})
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	return InvoicePaymentResult{Invoice: updatedInvoice, PaymentIntent: createdIntent}, nil
}

// SendInvoice records local email-send evidence for an open invoice.
// No real email is delivered; metadata and timeline/webhook consumers observe the attempt.
func (s *Service) SendInvoice(ctx context.Context, invoiceID string) (Invoice, error) {
	if strings.TrimSpace(invoiceID) == "" {
		return Invoice{}, fmt.Errorf("%w: invoice is required", ErrInvalidInput)
	}
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return Invoice{}, err
	}
	switch strings.ToLower(strings.TrimSpace(invoice.Status)) {
	case "open", "paid":
		// Stripe allows send for open and paid (paid emails omit payment reference).
	case "draft":
		return Invoice{}, fmt.Errorf("%w: Invoice must be finalized before it can be sent.", ErrInvalidInput)
	case "void", "uncollectible":
		return Invoice{}, fmt.Errorf("%w: Invoice is no longer open.", ErrInvalidInput)
	default:
		return Invoice{}, fmt.Errorf("%w: Invoice is no longer open.", ErrInvalidInput)
	}

	at := s.now().UTC()
	invoice.Metadata = copyMap(invoice.Metadata)
	if invoice.Metadata == nil {
		invoice.Metadata = map[string]string{}
	}
	invoice.Metadata["billtap_email_sent_at"] = at.Format(time.RFC3339)
	if invoice.CustomerID != "" {
		if customer, err := s.repo.GetCustomer(ctx, invoice.CustomerID); err == nil {
			if email := strings.TrimSpace(customer.Email); email != "" {
				invoice.Metadata["billtap_email_recipient"] = email
			}
		}
	}
	return s.repo.UpdateInvoice(ctx, invoice, []TimelineEntry{billingTimelineEntry(
		"invoice_sent_"+invoice.ID+"_"+at.Format(time.RFC3339Nano),
		"invoice.sent",
		"Invoice email evidence recorded",
		ObjectInvoice,
		invoice.ID,
		invoice.CustomerID,
		"",
		invoice.SubscriptionID,
		invoice.ID,
		invoice.PaymentIntentID,
		map[string]string{"source": "invoice.send", "status": invoice.Status},
		at,
	)})
}

func (s *Service) PayInvoice(ctx context.Context, invoiceID string, opts InvoicePaymentOptions) (InvoicePaymentResult, error) {
	if strings.TrimSpace(invoiceID) == "" {
		return InvoicePaymentResult{}, fmt.Errorf("%w: invoice is required", ErrInvalidInput)
	}
	at := opts.At
	if at.IsZero() {
		at = s.now()
	}
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	if invoice.SubscriptionID == "" {
		return s.payManualInvoice(ctx, invoice, opts, at)
	}
	subscription, err := s.repo.GetSubscription(ctx, invoice.SubscriptionID)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	intent, err := s.repo.GetPaymentIntent(ctx, invoice.PaymentIntentID)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	if invoice.Status != "open" {
		return InvoicePaymentResult{}, fmt.Errorf("%w: status must be open", ErrInvalidInput)
	}

	outcome := firstNonEmpty(opts.Outcome, opts.PaymentMethodID, "payment_succeeded")
	if opts.PaidOutOfBand {
		outcome = "payment_succeeded"
	}
	spec, ok := intentOutcomeSpec(outcome)
	if !ok {
		return InvoicePaymentResult{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
	}
	if opts.PaymentMethodID != "" {
		intent.PaymentMethodID = opts.PaymentMethodID
	}
	intent.PaymentMethodID = firstNonEmpty(intent.PaymentMethodID, spec.PaymentMethodID)
	intent.Status = spec.PaymentIntentStatus
	intent.FailureCode = spec.FailureCode
	intent.DeclineCode = spec.DeclineCode
	intent.FailureMessage = spec.FailureMessage
	invoice.AttemptCount++
	if invoice.AttemptCount <= 0 {
		invoice.AttemptCount = 1
	}

	success := intent.Status == "succeeded"
	if success {
		invoice.Status = "paid"
		invoice.AmountPaid = invoice.Total
		invoice.AmountDue = 0
		invoice.NextPaymentAttempt = nil
		intent.FailureCode = ""
		intent.DeclineCode = ""
		intent.FailureMessage = ""
		if subscription.Status != "canceled" {
			subscription.Status = "active"
		}
	} else {
		invoice.Status = "open"
		invoice.AmountPaid = 0
		invoice.AmountDue = invoice.Total
		nextAttempt := at.Add(24 * time.Hour)
		invoice.NextPaymentAttempt = &nextAttempt
		if subscription.Status != "canceled" {
			subscription.Status = "past_due"
		}
	}
	subscription.Metadata = copyMap(subscription.Metadata)
	subscription.Metadata["billtap_last_invoice_payment_attempt"] = at.Format(time.RFC3339Nano)
	subscription.Metadata["billtap_last_invoice_payment_outcome"] = outcome

	timeline := invoicePaymentTimeline(subscription, invoice, intent, success, at)
	subscription, invoice, intent, err = s.repo.UpdateInvoicePayment(ctx, subscription, invoice, intent, timeline)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return InvoicePaymentResult{}, fmt.Errorf("%w: status must be open or payment attempt is stale", ErrInvalidInput)
		}
		return InvoicePaymentResult{}, err
	}
	return InvoicePaymentResult{Invoice: invoice, Subscription: subscription, PaymentIntent: intent}, nil
}

func (s *Service) payManualInvoice(ctx context.Context, invoice Invoice, opts InvoicePaymentOptions, at time.Time) (InvoicePaymentResult, error) {
	if invoice.Status == "draft" {
		finalized, err := s.FinalizeInvoice(ctx, invoice.ID)
		if err != nil {
			return InvoicePaymentResult{}, err
		}
		invoice = finalized.Invoice
	}
	if invoice.Status != "open" {
		return InvoicePaymentResult{}, fmt.Errorf("%w: status must be open", ErrInvalidInput)
	}
	intent, err := s.repo.GetPaymentIntent(ctx, invoice.PaymentIntentID)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	configuredOutcome := paymentIntentConfiguredOutcome(intent.Metadata)
	outcome := firstNonEmpty(opts.Outcome, configuredOutcome)
	if opts.PaidOutOfBand {
		outcome = "payment_succeeded"
	}
	if outcome == "" && invoice.CustomerID != "" {
		if customer, err := s.repo.GetCustomer(ctx, invoice.CustomerID); err == nil {
			outcome = CustomerDefaultPaymentIntentOutcome(customer.Metadata)
		}
	}
	outcome = firstNonEmpty(outcome, opts.PaymentMethodID, "payment_succeeded")
	spec, ok := intentOutcomeSpec(outcome)
	if !ok {
		return InvoicePaymentResult{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
	}
	if opts.PaymentMethodID != "" {
		intent.PaymentMethodID = opts.PaymentMethodID
	}
	intent.PaymentMethodID = firstNonEmpty(intent.PaymentMethodID, spec.PaymentMethodID)
	intent.Status = spec.PaymentIntentStatus
	intent.FailureCode = spec.FailureCode
	intent.DeclineCode = spec.DeclineCode
	intent.FailureMessage = spec.FailureMessage
	invoice.AttemptCount++
	if invoice.AttemptCount <= 0 {
		invoice.AttemptCount = 1
	}
	if intent.Status == "succeeded" {
		invoice.Status = "paid"
		invoice.AmountPaid = invoice.Total
		invoice.AmountDue = 0
		invoice.NextPaymentAttempt = nil
		intent.FailureCode = ""
		intent.DeclineCode = ""
		intent.FailureMessage = ""
	} else {
		invoice.Status = "open"
		invoice.AmountPaid = 0
		invoice.AmountDue = invoice.Total
		nextAttempt := at.Add(24 * time.Hour)
		invoice.NextPaymentAttempt = &nextAttempt
	}
	invoice.Metadata = copyMap(invoice.Metadata)
	invoice.Metadata["billtap_last_invoice_payment_attempt"] = at.Format(time.RFC3339Nano)
	invoice.Metadata["billtap_last_invoice_payment_outcome"] = outcome
	timeline := []TimelineEntry{
		billingTimelineEntry("manual_invoice_payment_intent_"+intent.ID+"_"+at.Format(time.RFC3339Nano), paymentIntentEvent(intent.Status), "Invoice payment intent "+intent.Status, ObjectPaymentIntent, intent.ID, intent.CustomerID, "", "", invoice.ID, intent.ID, map[string]string{"source": "invoice.pay", "status": intent.Status, "outcome": outcome}, at),
	}
	for _, eventType := range invoiceEventTypesForPayment(invoice.Status, intent.Status) {
		timeline = append(timeline, billingTimelineEntry("manual_invoice_"+eventType+"_"+invoice.ID+"_"+at.Format(time.RFC3339Nano), eventType, "Invoice "+invoice.Status, ObjectInvoice, invoice.ID, invoice.CustomerID, "", "", invoice.ID, intent.ID, map[string]string{"source": "invoice.pay", "status": invoice.Status, "outcome": outcome}, at))
	}
	updatedInvoice, updatedIntent, err := s.repo.UpdateManualInvoicePayment(ctx, invoice, intent, timeline)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return InvoicePaymentResult{}, fmt.Errorf("%w: status must be open or payment attempt is stale", ErrInvalidInput)
		}
		return InvoicePaymentResult{}, err
	}
	return InvoicePaymentResult{Invoice: updatedInvoice, PaymentIntent: updatedIntent}, nil
}

func (s *Service) AdvanceClock(ctx context.Context, at time.Time) (ClockAdvanceResult, error) {
	return s.advanceClock(ctx, at, "")
}

func (s *Service) advanceClock(ctx context.Context, at time.Time, testClockID string) (ClockAdvanceResult, error) {
	if at.IsZero() {
		at = s.now()
	}
	result := ClockAdvanceResult{Object: "clock_advance", AdvancedTo: at, TestClockID: strings.TrimSpace(testClockID)}
	subscriptions, err := s.repo.ListSubscriptions(ctx)
	if err != nil {
		return result, err
	}
	for _, sub := range subscriptions {
		if result.TestClockID != "" && !s.subscriptionAttachedToClock(ctx, sub, result.TestClockID) {
			continue
		}
		if sub.Status == "canceled" {
			result.Skipped = append(result.Skipped, sub.ID)
			continue
		}
		current := sub
		for cycles := 0; !current.CurrentPeriodEnd.IsZero() && !current.CurrentPeriodEnd.After(at) && cycles < 24; cycles++ {
			result.Processed++
			if current.CancelAtPeriodEnd {
				canceled, err := s.cancelSubscriptionAtClock(ctx, current, current.CurrentPeriodEnd)
				if err != nil {
					return result, err
				}
				result.Canceled = append(result.Canceled, canceled)
				result.CanceledCount++
				break
			}
			if current.Status != "active" && current.Status != "trialing" {
				result.Skipped = append(result.Skipped, current.ID)
				break
			}
			if current.Status == "trialing" {
				activated, err := s.activateTrialSubscriptionAtClock(ctx, current, current.CurrentPeriodEnd)
				if err != nil {
					return result, err
				}
				result.Activated = append(result.Activated, activated)
				result.ActivatedCount++
				current = activated
				continue
			}
			renewal, err := s.renewSubscription(ctx, current, current.CurrentPeriodEnd)
			if err != nil {
				return result, err
			}
			result.Renewals = append(result.Renewals, renewal)
			result.Renewed++
			current = renewal.Subscription
		}
	}
	settled, err := s.settlePendingRefunds(ctx, at, result.TestClockID)
	if err != nil {
		return result, err
	}
	result.SettledRefunds = settled
	result.RefundCount = len(settled)
	result.Processed += len(settled)
	return result, nil
}

func (s *Service) CreateTestClock(ctx context.Context, in TestClock) (TestClock, error) {
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("clock")
	}
	in.Object = ObjectTestClock
	in.Status = firstNonEmpty(strings.TrimSpace(in.Status), "ready")
	if in.FrozenTime.IsZero() {
		in.FrozenTime = now
	}
	if in.CreatedAt.IsZero() {
		in.CreatedAt = now
	}
	if in.UpdatedAt.IsZero() {
		in.UpdatedAt = in.CreatedAt
	}
	return s.repo.CreateTestClock(ctx, in)
}

func (s *Service) GetTestClock(ctx context.Context, clockID string) (TestClock, error) {
	return s.repo.GetTestClock(ctx, clockID)
}

func (s *Service) ListTestClocks(ctx context.Context) ([]TestClock, error) {
	return s.repo.ListTestClocks(ctx)
}

func (s *Service) UpdateTestClock(ctx context.Context, clock TestClock) (TestClock, error) {
	clock.UpdatedAt = s.now()
	if strings.TrimSpace(clock.Status) == "" {
		clock.Status = "ready"
	}
	return s.repo.UpdateTestClock(ctx, clock)
}

func (s *Service) AdvanceTestClock(ctx context.Context, clockID string, frozenTime time.Time) (TestClock, ClockAdvanceResult, error) {
	clock, err := s.repo.GetTestClock(ctx, clockID)
	if err != nil {
		return TestClock{}, ClockAdvanceResult{}, err
	}
	if frozenTime.IsZero() {
		return TestClock{}, ClockAdvanceResult{}, fmt.Errorf("%w: frozen_time is required", ErrInvalidInput)
	}
	if frozenTime.Before(clock.FrozenTime) {
		return TestClock{}, ClockAdvanceResult{}, fmt.Errorf("%w: frozen_time must not move backwards", ErrInvalidInput)
	}
	clock.Status = "ready"
	clock.FrozenTime = frozenTime
	clock.UpdatedAt = s.now()
	updated, err := s.repo.UpdateTestClock(ctx, clock)
	if err != nil {
		return TestClock{}, ClockAdvanceResult{}, err
	}
	result, err := s.advanceClock(ctx, frozenTime, clock.ID)
	if err != nil {
		return updated, result, err
	}
	return updated, result, nil
}

func (s *Service) CreateRefund(ctx context.Context, in Refund) (Refund, error) {
	if in.Amount <= 0 {
		return Refund{}, fmt.Errorf("%w: amount must be at least 1", ErrInvalidInput)
	}
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("re")
	}
	in.Object = ObjectRefund
	in.Status = firstNonEmpty(strings.TrimSpace(in.Status), "succeeded")
	in.Currency = strings.ToLower(firstNonEmpty(strings.TrimSpace(in.Currency), "usd"))
	in.ChargeID = strings.TrimSpace(in.ChargeID)
	in.PaymentIntentID = strings.TrimSpace(in.PaymentIntentID)
	in.InvoiceID = strings.TrimSpace(in.InvoiceID)
	if in.PaymentIntentID != "" {
		intent, err := s.repo.GetPaymentIntent(ctx, in.PaymentIntentID)
		if err != nil {
			return Refund{}, err
		}
		in.CustomerID = firstNonEmpty(in.CustomerID, intent.CustomerID)
		in.InvoiceID = firstNonEmpty(in.InvoiceID, intent.InvoiceID)
		in.Currency = firstNonEmpty(in.Currency, intent.Currency)
	}
	if in.InvoiceID != "" {
		invoice, err := s.repo.GetInvoice(ctx, in.InvoiceID)
		if err != nil {
			return Refund{}, err
		}
		in.CustomerID = firstNonEmpty(in.CustomerID, invoice.CustomerID)
		in.Currency = firstNonEmpty(in.Currency, invoice.Currency)
		if in.PaymentIntentID == "" {
			in.PaymentIntentID = invoice.PaymentIntentID
		}
	}
	if in.ChargeID == "" && in.PaymentIntentID != "" {
		in.ChargeID = "ch_" + sanitizeID(in.PaymentIntentID)
	}
	if in.ChargeID == "" {
		return Refund{}, fmt.Errorf("%w: charge or payment_intent is required", ErrInvalidInput)
	}
	if in.CreatedAt.IsZero() {
		in.CreatedAt = now
	}
	return s.repo.CreateRefund(ctx, in, []TimelineEntry{billingTimelineEntry(
		"refund_"+in.ID,
		"charge.refunded",
		"Charge refunded",
		ObjectRefund,
		in.ID,
		in.CustomerID,
		"",
		"",
		in.InvoiceID,
		in.PaymentIntentID,
		map[string]string{"source": "refund.create", "charge": in.ChargeID, "status": in.Status, "reason": in.Reason},
		in.CreatedAt,
	)})
}

func (s *Service) GetRefund(ctx context.Context, refundID string) (Refund, error) {
	return s.repo.GetRefund(ctx, refundID)
}

func (s *Service) ListRefunds(ctx context.Context, filter RefundFilter) ([]Refund, error) {
	return s.repo.ListRefundsFiltered(ctx, filter)
}

func (s *Service) UpdateRefundStatus(ctx context.Context, refundID string, status string, at time.Time) (Refund, error) {
	refund, err := s.repo.GetRefund(ctx, refundID)
	if err != nil {
		return Refund{}, err
	}
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case "pending", "succeeded", "failed", "canceled":
	default:
		return Refund{}, fmt.Errorf("%w: status must be pending, succeeded, failed, or canceled", ErrInvalidInput)
	}
	if at.IsZero() {
		at = s.now()
	}
	refund.Status = status
	refund.Metadata = copyMap(refund.Metadata)
	refund.Metadata["billtap_last_status_update"] = at.Format(time.RFC3339Nano)
	return s.repo.UpdateRefund(ctx, refund, []TimelineEntry{billingTimelineEntry(
		"refund_status_"+refund.ID+"_"+status+"_"+at.Format(time.RFC3339Nano),
		"charge.refund.updated",
		"Refund "+status,
		ObjectRefund,
		refund.ID,
		refund.CustomerID,
		"",
		"",
		refund.InvoiceID,
		refund.PaymentIntentID,
		map[string]string{"source": "refund.update", "status": refund.Status, "charge": refund.ChargeID},
		at,
	)})
}

func (s *Service) CreateCreditNote(ctx context.Context, in CreditNote) (CreditNote, error) {
	if strings.TrimSpace(in.InvoiceID) == "" {
		return CreditNote{}, fmt.Errorf("%w: invoice is required", ErrInvalidInput)
	}
	if in.Amount <= 0 {
		return CreditNote{}, fmt.Errorf("%w: amount must be at least 1", ErrInvalidInput)
	}
	invoice, err := s.repo.GetInvoice(ctx, in.InvoiceID)
	if err != nil {
		return CreditNote{}, err
	}
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("cn")
	}
	in.Object = ObjectCreditNote
	in.Status = firstNonEmpty(strings.TrimSpace(in.Status), "issued")
	in.CustomerID = firstNonEmpty(in.CustomerID, invoice.CustomerID)
	in.Currency = strings.ToLower(firstNonEmpty(strings.TrimSpace(in.Currency), invoice.Currency, "usd"))
	if in.CreatedAt.IsZero() {
		in.CreatedAt = now
	}
	return s.repo.CreateCreditNote(ctx, in, []TimelineEntry{billingTimelineEntry(
		"credit_note_"+in.ID,
		"credit_note.created",
		"Credit note created",
		ObjectCreditNote,
		in.ID,
		in.CustomerID,
		"",
		invoice.SubscriptionID,
		in.InvoiceID,
		invoice.PaymentIntentID,
		map[string]string{"source": "credit_note.create", "status": in.Status, "reason": in.Reason},
		in.CreatedAt,
	)})
}

func (s *Service) GetCreditNote(ctx context.Context, creditNoteID string) (CreditNote, error) {
	return s.repo.GetCreditNote(ctx, creditNoteID)
}

func (s *Service) ListCreditNotes(ctx context.Context, filter CreditNoteFilter) ([]CreditNote, error) {
	return s.repo.ListCreditNotesFiltered(ctx, filter)
}

func (s *Service) VoidCreditNote(ctx context.Context, creditNoteID string) (CreditNote, error) {
	note, err := s.repo.GetCreditNote(ctx, creditNoteID)
	if err != nil {
		return CreditNote{}, err
	}
	if note.Status == "void" {
		return note, nil
	}
	if note.Status != "issued" {
		return CreditNote{}, fmt.Errorf("%w: status must be issued", ErrInvalidInput)
	}
	now := s.now()
	note.Status = "void"
	note.Metadata = copyMap(note.Metadata)
	note.Metadata["billtap_voided_at"] = now.Format(time.RFC3339Nano)
	return s.repo.UpdateCreditNote(ctx, note, []TimelineEntry{billingTimelineEntry(
		"credit_note_voided_"+note.ID+"_"+now.Format(time.RFC3339Nano),
		"credit_note.voided",
		"Credit note voided",
		ObjectCreditNote,
		note.ID,
		note.CustomerID,
		"",
		"",
		note.InvoiceID,
		"",
		map[string]string{"source": "credit_note.void", "status": note.Status},
		now,
	)})
}

func (s *Service) GetPaymentIntent(ctx context.Context, id string) (PaymentIntent, error) {
	return s.repo.GetPaymentIntent(ctx, id)
}

func (s *Service) CreatePaymentIntent(ctx context.Context, in PaymentIntent) (PaymentIntent, error) {
	if in.Amount <= 0 {
		return PaymentIntent{}, fmt.Errorf("%w: amount must be at least 1", ErrInvalidInput)
	}
	if strings.TrimSpace(in.Currency) == "" {
		return PaymentIntent{}, fmt.Errorf("%w: currency is required", ErrInvalidInput)
	}
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("pi")
	}
	in.Object = ObjectPaymentIntent
	in.Currency = strings.ToLower(strings.TrimSpace(in.Currency))
	in.CaptureMethod = firstNonEmpty(in.CaptureMethod, "automatic")
	in.Status = firstNonEmpty(in.Status, "requires_payment_method")
	if paymentIntentConfiguredOutcome(in.Metadata) == "" && in.CustomerID != "" && strings.TrimSpace(in.InvoiceID) == "" {
		customer, err := s.repo.GetCustomer(ctx, in.CustomerID)
		if err != nil {
			return PaymentIntent{}, err
		}
		if outcome := CustomerDefaultPaymentIntentOutcome(customer.Metadata); outcome != "" {
			if in.Metadata == nil {
				in.Metadata = map[string]string{}
			}
			in.Metadata[MetadataPaymentIntentOutcome] = outcome
		}
	}
	if outcome := paymentIntentConfiguredOutcome(in.Metadata); outcome != "" {
		if !IsSupportedPaymentIntentOutcome(outcome) {
			return PaymentIntent{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
		}
	}
	in.CreatedAt = now
	created, err := s.repo.CreatePaymentIntent(ctx, in)
	if err != nil {
		return PaymentIntent{}, err
	}
	return created, s.repo.RecordTimeline(ctx, timelineEntry("pi_"+created.ID+"_created", "payment_intent.created", "Payment intent created", ObjectPaymentIntent, created.ID, created.CustomerID, "", "", created.ID, map[string]string{"status": created.Status}, now))
}

func (s *Service) ConfirmPaymentIntent(ctx context.Context, id string, paymentMethodID string, outcome string) (PaymentIntent, error) {
	intent, err := s.repo.GetPaymentIntent(ctx, id)
	if err != nil {
		return PaymentIntent{}, err
	}
	if err := ensurePaymentIntentConfirmable(intent); err != nil {
		return PaymentIntent{}, err
	}
	configuredOutcome := paymentIntentConfiguredOutcome(intent.Metadata)
	effectiveOutcome := firstNonEmpty(outcome, configuredOutcome)
	if firstNonEmpty(paymentMethodID, intent.PaymentMethodID, effectiveOutcome) == "" {
		return PaymentIntent{}, fmt.Errorf("%w: payment_method is required", ErrInvalidInput)
	}
	if paymentMethodID != "" {
		intent.PaymentMethodID = paymentMethodID
	}
	spec, ok := intentOutcomeSpec(firstNonEmpty(effectiveOutcome, intent.PaymentMethodID))
	if !ok {
		return PaymentIntent{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, effectiveOutcome)
	}
	intent.PaymentMethodID = firstNonEmpty(intent.PaymentMethodID, spec.PaymentMethodID)
	intent.Status = spec.PaymentIntentStatus
	if intent.Status == "succeeded" && intent.CaptureMethod == "manual" {
		intent.Status = "requires_capture"
	}
	intent.FailureCode = spec.FailureCode
	intent.DeclineCode = spec.DeclineCode
	intent.FailureMessage = spec.FailureMessage
	if intent.Status == "succeeded" && intent.FailureCode == "" {
		intent.FailureMessage = ""
	}
	if intent.Status == "" {
		intent.Status = "succeeded"
	}
	now := s.now()
	return s.repo.UpdatePaymentIntent(ctx, intent, []TimelineEntry{
		timelineEntry("pi_"+intent.ID+"_confirmed_"+now.Format(time.RFC3339Nano), paymentIntentEvent(intent.Status), "Payment intent "+intent.Status, ObjectPaymentIntent, intent.ID, intent.CustomerID, "", "", intent.ID, map[string]string{"status": intent.Status, "outcome": firstNonEmpty(effectiveOutcome, intent.PaymentMethodID)}, now),
	})
}

func (s *Service) SetPaymentIntentOutcome(ctx context.Context, id string, outcome string) (PaymentIntent, error) {
	outcome = strings.TrimSpace(outcome)
	if outcome == "" {
		return PaymentIntent{}, fmt.Errorf("%w: outcome is required", ErrInvalidInput)
	}
	if _, ok := intentOutcomeSpec(outcome); !ok {
		return PaymentIntent{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
	}
	intent, err := s.repo.GetPaymentIntent(ctx, id)
	if err != nil {
		return PaymentIntent{}, err
	}
	if intent.Metadata == nil {
		intent.Metadata = map[string]string{}
	}
	intent.Metadata[MetadataPaymentIntentOutcome] = outcome
	now := s.now()
	return s.repo.UpdatePaymentIntent(ctx, intent, []TimelineEntry{
		timelineEntry("pi_"+intent.ID+"_outcome_"+now.Format(time.RFC3339Nano), "payment_intent.outcome_configured", "Payment intent outcome configured", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", "", intent.ID, map[string]string{"outcome": outcome}, now),
	})
}

func (s *Service) CapturePaymentIntent(ctx context.Context, id string, amountToCapture int64) (PaymentIntent, error) {
	intent, err := s.repo.GetPaymentIntent(ctx, id)
	if err != nil {
		return PaymentIntent{}, err
	}
	if intent.Status != "requires_capture" {
		return PaymentIntent{}, fmt.Errorf("%w: status must be requires_capture", ErrInvalidInput)
	}
	if amountToCapture != 0 && amountToCapture != intent.Amount {
		return PaymentIntent{}, fmt.Errorf("%w: amount_to_capture must be %d", ErrInvalidInput, intent.Amount)
	}
	intent.Status = "succeeded"
	intent.FailureCode = ""
	intent.DeclineCode = ""
	intent.FailureMessage = ""
	now := s.now()
	return s.repo.UpdatePaymentIntent(ctx, intent, []TimelineEntry{
		timelineEntry("pi_"+intent.ID+"_captured_"+now.Format(time.RFC3339Nano), "payment_intent.succeeded", "Payment intent captured", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", "", intent.ID, map[string]string{"status": intent.Status}, now),
	})
}

func (s *Service) CancelPaymentIntent(ctx context.Context, id string) (PaymentIntent, error) {
	intent, err := s.repo.GetPaymentIntent(ctx, id)
	if err != nil {
		return PaymentIntent{}, err
	}
	if intent.Status == "succeeded" || intent.Status == "canceled" {
		return PaymentIntent{}, fmt.Errorf("%w: status must be non-terminal", ErrInvalidInput)
	}
	intent.Status = "canceled"
	now := s.now()
	return s.repo.UpdatePaymentIntent(ctx, intent, []TimelineEntry{
		timelineEntry("pi_"+intent.ID+"_canceled_"+now.Format(time.RFC3339Nano), "payment_intent.canceled", "Payment intent canceled", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", "", intent.ID, map[string]string{"status": intent.Status}, now),
	})
}

func (s *Service) SettleBankTransferPaymentIntents(ctx context.Context, customerID string) ([]PaymentIntent, error) {
	if strings.TrimSpace(customerID) == "" {
		return nil, fmt.Errorf("%w: customer is required", ErrInvalidInput)
	}
	intents, err := s.repo.ListPaymentIntentsFiltered(ctx, PaymentIntentFilter{CustomerID: customerID})
	if err != nil {
		return nil, err
	}
	now := s.now()
	var settled []PaymentIntent
	for _, intent := range intents {
		if intent.Status != "processing" || !isBankTransferPaymentMethod(intent.PaymentMethodID) {
			continue
		}
		intent.Status = "succeeded"
		intent.FailureCode = ""
		intent.DeclineCode = ""
		intent.FailureMessage = ""
		updated, err := s.repo.UpdatePaymentIntent(ctx, intent, []TimelineEntry{
			timelineEntry("pi_"+intent.ID+"_bank_transfer_settled_"+now.Format(time.RFC3339Nano), "payment_intent.succeeded", "Bank transfer payment intent settled", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", "", intent.ID, map[string]string{"status": intent.Status, "source": "fund_cash_balance"}, now),
		})
		if err != nil {
			return settled, err
		}
		settled = append(settled, updated)
	}
	return settled, nil
}

func (s *Service) ListPaymentIntents(ctx context.Context) ([]PaymentIntent, error) {
	return s.repo.ListPaymentIntents(ctx)
}

func (s *Service) CreateSetupIntent(ctx context.Context, in SetupIntent) (SetupIntent, error) {
	now := s.now()
	if strings.TrimSpace(in.ID) == "" {
		in.ID = id("seti")
	}
	in.Object = ObjectSetupIntent
	in.Status = firstNonEmpty(in.Status, "requires_payment_method")
	in.Usage = firstNonEmpty(in.Usage, "off_session")
	in.CreatedAt = now
	created, err := s.repo.CreateSetupIntent(ctx, in)
	if err != nil {
		return SetupIntent{}, err
	}
	return created, s.repo.RecordTimeline(ctx, timelineEntry("seti_"+created.ID+"_created", "setup_intent.created", "Setup intent created", ObjectSetupIntent, created.ID, created.CustomerID, "", "", "", map[string]string{"status": created.Status}, now))
}

func (s *Service) GetSetupIntent(ctx context.Context, id string) (SetupIntent, error) {
	return s.repo.GetSetupIntent(ctx, id)
}

func (s *Service) ListSetupIntents(ctx context.Context) ([]SetupIntent, error) {
	return s.repo.ListSetupIntents(ctx)
}

func (s *Service) ConfirmSetupIntent(ctx context.Context, id string, paymentMethodID string, outcome string) (SetupIntent, error) {
	intent, err := s.repo.GetSetupIntent(ctx, id)
	if err != nil {
		return SetupIntent{}, err
	}
	if intent.Status == "succeeded" || intent.Status == "canceled" {
		return SetupIntent{}, fmt.Errorf("%w: status must be non-terminal", ErrInvalidInput)
	}
	if firstNonEmpty(paymentMethodID, intent.PaymentMethodID, outcome) == "" {
		return SetupIntent{}, fmt.Errorf("%w: payment_method is required", ErrInvalidInput)
	}
	if paymentMethodID != "" {
		intent.PaymentMethodID = paymentMethodID
	}
	spec, ok := intentOutcomeSpec(firstNonEmpty(outcome, intent.PaymentMethodID))
	if !ok {
		return SetupIntent{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
	}
	intent.PaymentMethodID = firstNonEmpty(intent.PaymentMethodID, spec.PaymentMethodID)
	intent.Status = "succeeded"
	intent.FailureCode = ""
	intent.DeclineCode = ""
	intent.FailureMessage = ""
	if spec.PaymentIntentStatus == "requires_action" {
		intent.Status = "requires_action"
		intent.FailureCode = spec.FailureCode
		intent.DeclineCode = spec.DeclineCode
		intent.FailureMessage = spec.FailureMessage
	}
	if spec.PaymentIntentStatus == "requires_payment_method" {
		intent.Status = "requires_payment_method"
		intent.FailureCode = spec.FailureCode
		intent.DeclineCode = spec.DeclineCode
		intent.FailureMessage = spec.FailureMessage
	}
	now := s.now()
	return s.repo.UpdateSetupIntent(ctx, intent, []TimelineEntry{
		timelineEntry("seti_"+intent.ID+"_confirmed_"+now.Format(time.RFC3339Nano), setupIntentEvent(intent.Status), "Setup intent "+intent.Status, ObjectSetupIntent, intent.ID, intent.CustomerID, "", "", "", map[string]string{"status": intent.Status}, now),
	})
}

func (s *Service) CancelSetupIntent(ctx context.Context, id string) (SetupIntent, error) {
	intent, err := s.repo.GetSetupIntent(ctx, id)
	if err != nil {
		return SetupIntent{}, err
	}
	if intent.Status == "succeeded" || intent.Status == "canceled" {
		return SetupIntent{}, fmt.Errorf("%w: status must be non-terminal", ErrInvalidInput)
	}
	intent.Status = "canceled"
	now := s.now()
	return s.repo.UpdateSetupIntent(ctx, intent, []TimelineEntry{
		timelineEntry("seti_"+intent.ID+"_canceled_"+now.Format(time.RFC3339Nano), "setup_intent.canceled", "Setup intent canceled", ObjectSetupIntent, intent.ID, intent.CustomerID, "", "", "", map[string]string{"status": intent.Status}, now),
	})
}

func (s *Service) Timeline(ctx context.Context, filter TimelineFilter) ([]TimelineEntry, error) {
	return s.repo.Timeline(ctx, filter)
}

func (s *Service) PortalState(ctx context.Context, customerID string) (PortalState, error) {
	if strings.TrimSpace(customerID) == "" {
		return PortalState{}, fmt.Errorf("%w: customer is required", ErrInvalidInput)
	}
	customer, err := s.repo.GetCustomer(ctx, customerID)
	if err != nil {
		return PortalState{}, err
	}
	subscriptions, err := s.repo.ListSubscriptionsByCustomer(ctx, customer.ID)
	if err != nil {
		return PortalState{}, err
	}
	subscription := currentPortalSubscription(subscriptions)

	invoiceFilter := InvoiceFilter{CustomerID: customer.ID}
	if subscription != nil {
		invoiceFilter.SubscriptionID = subscription.ID
	}
	invoices, err := s.repo.ListInvoicesFiltered(ctx, invoiceFilter)
	if err != nil {
		return PortalState{}, err
	}
	paymentIntentFilter := PaymentIntentFilter{CustomerID: customer.ID}
	if subscription != nil {
		paymentIntentFilter.InvoiceIDs = invoiceIDs(invoices)
	}
	paymentIntents, err := s.repo.ListPaymentIntentsFiltered(ctx, paymentIntentFilter)
	if err != nil {
		return PortalState{}, err
	}

	timeline, err := s.repo.Timeline(ctx, TimelineFilter{CustomerID: customer.ID})
	if err != nil {
		return PortalState{}, err
	}

	return PortalState{
		Object:         "portal_state",
		Customer:       customer,
		Subscription:   subscription,
		Invoices:       invoices,
		PaymentIntents: paymentIntents,
		Summary:        portalSummary(customer.ID, subscription, invoices, paymentIntents),
		Timeline:       timeline,
	}, nil
}

func (s *Service) ChangePortalPlan(ctx context.Context, subscriptionID string, change PortalPlanChange) (Subscription, error) {
	sub, err := s.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return Subscription{}, err
	}
	if sub.Status == "canceled" {
		return Subscription{}, fmt.Errorf("%w: canceled subscriptions must be resumed before plan changes", ErrInvalidInput)
	}

	priceID := strings.TrimSpace(change.PriceID)
	if priceID == "" && len(sub.Items) > 0 {
		priceID = sub.Items[0].PriceID
	}
	if priceID == "" {
		return Subscription{}, fmt.Errorf("%w: price is required", ErrInvalidInput)
	}
	if _, err := s.repo.GetPrice(ctx, priceID); err != nil {
		return Subscription{}, err
	}
	quantity := change.Quantity
	if quantity <= 0 && len(sub.Items) > 0 {
		quantity = sub.Items[0].Quantity
	}
	if quantity <= 0 {
		quantity = 1
	}

	previousPrice, previousQuantity := firstItem(sub.Items)
	sub.Items = []LineItem{{PriceID: priceID, Quantity: quantity}}
	sub.Metadata = copyMap(sub.Metadata)
	sub.Metadata["portal_last_action"] = "plan_change"
	sub.Metadata["portal_updated_at"] = s.now().Format(time.RFC3339Nano)
	if planID := strings.TrimSpace(change.PlanID); planID != "" {
		sub.Metadata["plan"] = planID
	}

	now := s.now()
	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{portalTimeline(
		"portal_plan_change_"+sub.ID+"_"+now.Format(time.RFC3339Nano),
		"customer.subscription.updated",
		"Portal plan changed",
		sub,
		map[string]string{
			"portal_action":     "plan_change",
			"plan":              strings.TrimSpace(change.PlanID),
			"price":             priceID,
			"previous_price":    previousPrice,
			"quantity":          fmt.Sprintf("%d", quantity),
			"previous_quantity": fmt.Sprintf("%d", previousQuantity),
		},
		now,
	)})
}

func (s *Service) ChangePortalSeats(ctx context.Context, subscriptionID string, change PortalSeatChange) (Subscription, error) {
	if change.Quantity <= 0 {
		return Subscription{}, fmt.Errorf("%w: quantity must be positive", ErrInvalidInput)
	}
	sub, err := s.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return Subscription{}, err
	}
	if sub.Status == "canceled" {
		return Subscription{}, fmt.Errorf("%w: canceled subscriptions must be resumed before seat changes", ErrInvalidInput)
	}
	if len(sub.Items) == 0 {
		return Subscription{}, fmt.Errorf("%w: subscription has no items", ErrInvalidInput)
	}
	_, previousQuantity := firstItem(sub.Items)
	for idx := range sub.Items {
		sub.Items[idx].Quantity = change.Quantity
	}
	sub.Metadata = copyMap(sub.Metadata)
	sub.Metadata["portal_last_action"] = "seat_change"
	sub.Metadata["portal_updated_at"] = s.now().Format(time.RFC3339Nano)

	now := s.now()
	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{portalTimeline(
		"portal_seat_change_"+sub.ID+"_"+now.Format(time.RFC3339Nano),
		"customer.subscription.updated",
		"Portal seat quantity changed",
		sub,
		map[string]string{
			"portal_action":     "seat_change",
			"quantity":          fmt.Sprintf("%d", change.Quantity),
			"previous_quantity": fmt.Sprintf("%d", previousQuantity),
		},
		now,
	)})
}

func (s *Service) CancelPortalSubscription(ctx context.Context, subscriptionID string, cancel PortalCancel) (Subscription, error) {
	mode := strings.ToLower(strings.TrimSpace(cancel.Mode))
	if mode == "" {
		mode = "period"
	}
	if mode != "period" && mode != "immediate" {
		return Subscription{}, fmt.Errorf("%w: mode must be period or immediate", ErrInvalidInput)
	}
	sub, err := s.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return Subscription{}, err
	}
	now := s.now()
	sub.Metadata = copyMap(sub.Metadata)
	if sub.Metadata["status_before_cancel"] == "" && sub.Status != "canceled" {
		sub.Metadata["status_before_cancel"] = sub.Status
	}
	sub.Metadata["portal_last_action"] = "cancel"
	sub.Metadata["portal_cancel_mode"] = mode
	sub.Metadata["portal_updated_at"] = now.Format(time.RFC3339Nano)
	data := map[string]string{"portal_action": "cancel", "mode": mode}
	action := "customer.subscription.updated"
	message := "Portal cancellation scheduled"
	if mode == "period" {
		sub.CancelAtPeriodEnd = true
		if sub.CanceledAt == nil {
			sub.CanceledAt = &now
		}
		sub.Metadata["cancel_at"] = sub.CurrentPeriodEnd.Format(time.RFC3339Nano)
	} else {
		sub.Status = "canceled"
		sub.CancelAtPeriodEnd = false
		delete(sub.Metadata, "cancel_at")
		if sub.CanceledAt == nil {
			sub.CanceledAt = &now
		}
		action = "customer.subscription.deleted"
		message = "Portal subscription canceled immediately"
		data["status"] = sub.Status
	}
	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{portalTimeline(
		"portal_cancel_"+sub.ID+"_"+mode+"_"+now.Format(time.RFC3339Nano),
		action,
		message,
		sub,
		data,
		now,
	)})
}

func (s *Service) ResumePortalSubscription(ctx context.Context, subscriptionID string) (Subscription, error) {
	sub, err := s.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return Subscription{}, err
	}
	now := s.now()
	statusBeforeCancel := sub.Metadata["status_before_cancel"]
	if statusBeforeCancel == "trialing" {
		sub.Status = "trialing"
	} else {
		sub.Status = "active"
	}
	sub.CancelAtPeriodEnd = false
	sub.CanceledAt = nil
	sub.Metadata = copyMap(sub.Metadata)
	delete(sub.Metadata, "cancel_at")
	delete(sub.Metadata, "cancellation_details_comment")
	delete(sub.Metadata, "cancellation_details_feedback")
	sub.Metadata["portal_last_action"] = "resume"
	sub.Metadata["portal_updated_at"] = now.Format(time.RFC3339Nano)

	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{portalTimeline(
		"portal_resume_"+sub.ID+"_"+now.Format(time.RFC3339Nano),
		"customer.subscription.updated",
		"Portal subscription resumed",
		sub,
		map[string]string{"portal_action": "resume", "status": sub.Status},
		now,
	)})
}

func (s *Service) SimulatePaymentMethodUpdate(ctx context.Context, customerID string, outcome string, paymentMethodID string) (PaymentMethodSimulation, error) {
	if strings.TrimSpace(customerID) == "" {
		return PaymentMethodSimulation{}, fmt.Errorf("%w: customer is required", ErrInvalidInput)
	}
	customer, err := s.repo.GetCustomer(ctx, customerID)
	if err != nil {
		return PaymentMethodSimulation{}, err
	}
	outcome = normalizePaymentMethodOutcome(outcome)
	if outcome != "succeeds" && outcome != "fails" {
		return PaymentMethodSimulation{}, fmt.Errorf("%w: outcome must be succeeds or fails", ErrInvalidInput)
	}
	now := s.now()
	result := PaymentMethodSimulation{
		ID:         id("pmupd"),
		Object:     "payment_method_update",
		CustomerID: customerID,
		Outcome:    outcome,
		Status:     "succeeded",
		CreatedAt:  now,
	}
	action := "payment_method.updated"
	message := "Portal payment method update succeeded"
	if outcome == "succeeds" {
		result.PaymentMethodID = strings.TrimSpace(paymentMethodID)
		if result.PaymentMethodID == "" {
			result.PaymentMethodID = id("pm")
		}
		customer.Metadata = copyMap(customer.Metadata)
		customer.Metadata[MetadataDefaultPaymentMethod] = result.PaymentMethodID
		customer.Metadata["payment_method_status"] = "saved"
		if _, err := s.repo.UpdateCustomer(ctx, customer.ID, Customer{Metadata: customer.Metadata}); err != nil {
			return PaymentMethodSimulation{}, err
		}
	} else {
		result.Status = "failed"
		result.FailureCode = "card_declined"
		result.FailureMessage = "Simulated payment method update failure"
		action = "payment_method.update_failed"
		message = "Portal payment method update failed"
	}
	err = s.repo.RecordTimeline(ctx, TimelineEntry{
		ID:         id("tl"),
		Object:     ObjectTimelineEntry,
		Action:     action,
		Message:    message,
		ObjectType: ObjectCustomer,
		ObjectID:   customerID,
		CustomerID: customerID,
		Data: map[string]string{
			"portal_action":  "payment_method_update",
			"outcome":        outcome,
			"status":         result.Status,
			"payment_method": result.PaymentMethodID,
			"failure_code":   result.FailureCode,
		},
		CreatedAt: now,
	})
	if err != nil {
		return PaymentMethodSimulation{}, err
	}
	return result, nil
}

func (s *Service) cancelSubscriptionAtClock(ctx context.Context, sub Subscription, at time.Time) (Subscription, error) {
	sub.Status = "canceled"
	sub.CancelAtPeriodEnd = false
	canceledAt := at
	sub.CanceledAt = &canceledAt
	sub.Metadata = copyMap(sub.Metadata)
	sub.Metadata["billtap_clock_canceled_at"] = at.Format(time.RFC3339Nano)
	delete(sub.Metadata, "cancel_at")
	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{billingTimelineEntry(
		"clock_cancel_"+sub.ID+"_"+at.Format(time.RFC3339Nano),
		"customer.subscription.deleted",
		"Subscription canceled at period end",
		ObjectSubscription,
		sub.ID,
		sub.CustomerID,
		"",
		sub.ID,
		sub.LatestInvoiceID,
		"",
		map[string]string{"source": "clock.advance", "status": sub.Status},
		at,
	)})
}

func (s *Service) activateTrialSubscriptionAtClock(ctx context.Context, sub Subscription, at time.Time) (Subscription, error) {
	periodStart := sub.CurrentPeriodEnd
	if periodStart.IsZero() {
		periodStart = at
	}
	periodEnd, err := s.nextPeriodEnd(ctx, sub.Items, periodStart)
	if err != nil {
		return Subscription{}, err
	}
	sub.Status = "active"
	sub.CurrentPeriodStart = periodStart
	sub.CurrentPeriodEnd = periodEnd
	sub.Metadata = copyMap(sub.Metadata)
	sub.Metadata["billtap_trial_activated_at"] = at.Format(time.RFC3339Nano)
	sub.Metadata["billtap_last_period_start"] = periodStart.Format(time.RFC3339Nano)
	sub.Metadata["billtap_last_period_end"] = periodEnd.Format(time.RFC3339Nano)
	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{billingTimelineEntry(
		"clock_trial_activate_"+sub.ID+"_"+at.Format(time.RFC3339Nano),
		"customer.subscription.updated",
		"Trial subscription activated",
		ObjectSubscription,
		sub.ID,
		sub.CustomerID,
		"",
		sub.ID,
		sub.LatestInvoiceID,
		"",
		map[string]string{"source": "clock.advance", "status": sub.Status, "previous_status": "trialing"},
		at,
	)})
}

func (s *Service) renewSubscription(ctx context.Context, sub Subscription, at time.Time) (InvoicePaymentResult, error) {
	subtotal, currency, lineAmounts, err := s.subscriptionLineAmounts(ctx, sub.Items)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	sub.Metadata = copyMap(sub.Metadata)
	// create_prorations pending amount is added to the renewal subtotal, then
	// discounts and tax run on the combined base (same flow as a normal renewal).
	if pendingRaw := strings.TrimSpace(sub.Metadata[MetadataPendingProrationAmount]); pendingRaw != "" {
		if pending, parseErr := strconv.ParseInt(pendingRaw, 10, 64); parseErr == nil {
			subtotal += pending
		}
		delete(sub.Metadata, MetadataPendingProrationAmount)
		delete(sub.Metadata, MetadataPendingProrationAt)
	}
	discounts := DiscountsFromMetadata(sub.Metadata)
	eligibleBase := EligibleDiscountBase(subtotal, discounts, lineAmounts)
	total, discountAmount := ApplyDiscountsWithEligibleBase(subtotal, eligibleBase, currency, discounts)
	automaticTax, taxPercent := AutomaticTaxFromMetadata(sub.Metadata)
	defaultTaxRates := DefaultTaxRatesFromMetadata(sub.Metadata)
	tax := int64(0)
	if len(defaultTaxRates) > 0 {
		_, _, exclusiveTotal, taxTotal := ComputeTaxRateAmounts(total, defaultTaxRates)
		tax = taxTotal
		total = total + exclusiveTotal
	} else if automaticTax {
		tax = ExclusiveTaxAmount(total, taxPercent)
		total = total + tax
	}
	periodStart := sub.CurrentPeriodEnd
	if periodStart.IsZero() {
		periodStart = at
	}
	periodEnd, err := s.nextPeriodEnd(ctx, sub.Items, periodStart)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	sub.Status = "active"
	sub.CurrentPeriodStart = periodStart
	sub.CurrentPeriodEnd = periodEnd
	sub.Metadata["billtap_last_renewal_at"] = at.Format(time.RFC3339Nano)
	sub.Metadata["billtap_last_renewal_period_start"] = periodStart.Format(time.RFC3339Nano)
	sub.Metadata["billtap_last_renewal_period_end"] = periodEnd.Format(time.RFC3339Nano)

	renewalOutcome := renewalOutcome(sub.Metadata)
	if renewalOutcome == "" && sub.CustomerID != "" {
		if customer, err := s.repo.GetCustomer(ctx, sub.CustomerID); err == nil {
			renewalOutcome = CustomerDefaultInvoiceOutcome(customer.Metadata)
		}
	}
	renewalFailed := renewalOutcome != ""
	invoice := Invoice{
		ID:              id("in"),
		Object:          ObjectInvoice,
		CustomerID:      sub.CustomerID,
		SubscriptionID:  sub.ID,
		Status:          "paid",
		Currency:        currency,
		Subtotal:        subtotal,
		DiscountAmount:  discountAmount,
		Discounts:       discounts,
		AutomaticTax:    automaticTax && len(defaultTaxRates) == 0,
		DefaultTaxRates: defaultTaxRates,
		Tax:             tax,
		Total:           total,
		AmountDue:       0,
		AmountPaid:      total,
		AttemptCount:    1,
		CreatedAt:       at,
	}
	if renewalFailed {
		invoice.Status = "open"
		invoice.AmountDue = total
		invoice.AmountPaid = 0
		nextPaymentAttempt := at.Add(24 * time.Hour)
		invoice.NextPaymentAttempt = &nextPaymentAttempt
		sub.Status = renewalFailureSubscriptionStatus(renewalOutcome)
		sub.Metadata["billtap_last_renewal_outcome"] = renewalOutcome
		sub.Metadata["billtap_next_retry_at"] = nextPaymentAttempt.Format(time.RFC3339Nano)
	}
	intent := PaymentIntent{
		ID:              id("pi"),
		Object:          ObjectPaymentIntent,
		CustomerID:      sub.CustomerID,
		InvoiceID:       invoice.ID,
		Amount:          total,
		Currency:        currency,
		Status:          "succeeded",
		CaptureMethod:   "automatic",
		PaymentMethodID: "pm_card_visa",
		CreatedAt:       at,
	}
	if renewalFailed {
		spec, ok := intentOutcomeSpec(renewalOutcome)
		if !ok {
			spec, _ = intentOutcomeSpec("card_declined")
		}
		intent.Status = spec.PaymentIntentStatus
		intent.PaymentMethodID = firstNonEmpty(spec.PaymentMethodID, "pm_card_declined")
		intent.FailureCode = spec.FailureCode
		intent.DeclineCode = spec.DeclineCode
		intent.FailureMessage = spec.FailureMessage
	}
	invoice.PaymentIntentID = intent.ID
	sub.LatestInvoiceID = invoice.ID

	timeline := []TimelineEntry{
		billingTimelineEntry("renewal_invoice_created_"+invoice.ID, "invoice.created", "Renewal invoice created", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": invoice.Status}, at),
		billingTimelineEntry("renewal_invoice_finalized_"+invoice.ID, "invoice.finalized", "Renewal invoice finalized", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": invoice.Status}, at),
		billingTimelineEntry("renewal_payment_intent_created_"+intent.ID, "payment_intent.created", "Renewal payment intent created", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": intent.Status}, at),
	}
	if renewalFailed {
		timeline = append(timeline,
			billingTimelineEntry("renewal_payment_intent_failed_"+intent.ID, paymentIntentEvent(intent.Status), "Renewal payment intent failed", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": intent.Status, "outcome": renewalOutcome}, at),
			billingTimelineEntry("renewal_invoice_payment_failed_"+invoice.ID, "invoice.payment_failed", "Renewal invoice payment failed", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": invoice.Status, "outcome": renewalOutcome}, at),
			billingTimelineEntry("renewal_subscription_past_due_"+sub.ID+"_"+invoice.ID, "customer.subscription.updated", "Subscription updated after renewal payment failure", ObjectSubscription, sub.ID, sub.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": sub.Status, "outcome": renewalOutcome}, at),
		)
	} else {
		timeline = append(timeline,
			billingTimelineEntry("renewal_payment_intent_succeeded_"+intent.ID, "payment_intent.succeeded", "Renewal payment intent succeeded", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": intent.Status}, at),
			billingTimelineEntry("renewal_invoice_payment_succeeded_"+invoice.ID, "invoice.payment_succeeded", "Renewal invoice payment succeeded", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": invoice.Status}, at),
			billingTimelineEntry("renewal_invoice_paid_"+invoice.ID, "invoice.paid", "Renewal invoice paid", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": invoice.Status}, at),
			billingTimelineEntry("renewal_subscription_updated_"+sub.ID+"_"+invoice.ID, "customer.subscription.updated", "Subscription renewed", ObjectSubscription, sub.ID, sub.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": "clock.advance", "status": sub.Status}, at),
		)
	}

	sub, invoice, intent, err = s.repo.RecordSubscriptionRenewal(ctx, sub, invoice, intent, timeline)
	if err != nil {
		return InvoicePaymentResult{}, err
	}
	return InvoicePaymentResult{Invoice: invoice, Subscription: sub, PaymentIntent: intent}, nil
}

// SubscriptionProrationRequest describes an item change that may bill immediately.
type SubscriptionProrationRequest struct {
	SubscriptionID     string
	NewItems           []LineItem
	ProrationBehavior  string // "none" | "create_prorations" | "always_invoice"
	ProrationDate      time.Time
	BillingCycleAnchor string // "" | "unchanged" | "now"
	PaymentBehavior    string // "" | "error_if_incomplete" | …
	DefaultTaxRates    []AppliedTaxRate // optional override; empty uses subscription metadata
	Metadata           map[string]string
	CancelAtPeriodEnd  *bool
}

// SubscriptionProrationResult is the outcome of a proration-aware subscription update.
type SubscriptionProrationResult struct {
	Subscription  Subscription
	Invoice       *Invoice
	PaymentIntent *PaymentIntent
	// PaymentResult is populated when an invoice was recorded (for emitRenewalWebhooks).
	PaymentResult InvoicePaymentResult
}

// UpdateSubscriptionItemsWithProration applies item changes and, depending on
// proration_behavior / billing_cycle_anchor, may create a subscription_update invoice.
func (s *Service) UpdateSubscriptionItemsWithProration(ctx context.Context, req SubscriptionProrationRequest) (SubscriptionProrationResult, error) {
	if strings.TrimSpace(req.SubscriptionID) == "" {
		return SubscriptionProrationResult{}, fmt.Errorf("%w: subscription is required", ErrInvalidInput)
	}
	if len(req.NewItems) == 0 {
		return SubscriptionProrationResult{}, fmt.Errorf("%w: subscription items cannot be empty", ErrInvalidInput)
	}
	sub, err := s.repo.GetSubscription(ctx, req.SubscriptionID)
	if err != nil {
		return SubscriptionProrationResult{}, err
	}

	at := req.ProrationDate
	if at.IsZero() {
		at = s.now()
	}

	oldTotal, currency, oldLineAmounts, err := s.subscriptionLineAmounts(ctx, sub.Items)
	if err != nil {
		return SubscriptionProrationResult{}, err
	}
	newTotal, newCurrency, newLineAmounts, err := s.subscriptionLineAmounts(ctx, req.NewItems)
	if err != nil {
		return SubscriptionProrationResult{}, err
	}
	if currency == "" {
		currency = newCurrency
	}
	if currency == "" {
		currency = "usd"
	}

	behavior := strings.TrimSpace(req.ProrationBehavior)
	if behavior == "" {
		behavior = "none"
	}
	anchor := strings.TrimSpace(req.BillingCycleAnchor)
	paymentBehavior := strings.TrimSpace(req.PaymentBehavior)

	// Build the updated subscription snapshot (not yet committed).
	updated := sub
	updated.Items = append([]LineItem{}, req.NewItems...)
	updated.Metadata = copyMap(sub.Metadata)
	for key, value := range req.Metadata {
		if value == "" {
			delete(updated.Metadata, key)
		} else {
			updated.Metadata[key] = value
		}
	}
	// Discounts: prefer post-merge metadata (request may attach a coupon); fall back to prior.
	discounts := DiscountsFromMetadata(updated.Metadata)
	if len(discounts) == 0 {
		discounts = DiscountsFromMetadata(sub.Metadata)
	}
	eligibleOld := EligibleDiscountBase(oldTotal, discounts, oldLineAmounts)
	oldDiscounted, _ := ApplyDiscountsWithEligibleBase(oldTotal, eligibleOld, currency, discounts)
	eligibleNew := EligibleDiscountBase(newTotal, discounts, newLineAmounts)
	newDiscounted, _ := ApplyDiscountsWithEligibleBase(newTotal, eligibleNew, currency, discounts)
	if req.CancelAtPeriodEnd != nil {
		updated.CancelAtPeriodEnd = *req.CancelAtPeriodEnd
		if *req.CancelAtPeriodEnd {
			if updated.CanceledAt == nil {
				canceledAt := at
				updated.CanceledAt = &canceledAt
			}
			updated.Metadata["cancel_at"] = updated.CurrentPeriodEnd.Format(time.RFC3339Nano)
		} else {
			updated.CanceledAt = nil
			delete(updated.Metadata, "cancel_at")
			delete(updated.Metadata, "cancellation_details_comment")
			delete(updated.Metadata, "cancellation_details_feedback")
			if updated.Status == "canceled" {
				updated.Status = "active"
			}
		}
	}
	updated.Metadata["stripe_compat_updated_at"] = at.Format(time.RFC3339Nano)

	rates := req.DefaultTaxRates
	if len(rates) == 0 {
		rates = DefaultTaxRatesFromMetadata(updated.Metadata)
	}
	automaticTax, taxPercent := AutomaticTaxFromMetadata(updated.Metadata)

	// none: items + metadata only, no invoice, period unchanged.
	if behavior == "none" {
		saved, err := s.repo.UpdateSubscription(ctx, updated, []TimelineEntry{portalTimeline(
			"stripe_compat_update_"+updated.ID+"_"+at.Format(time.RFC3339Nano),
			"customer.subscription.updated",
			"Stripe-compatible subscription updated",
			updated,
			map[string]string{"source": "stripe_compat", "status": updated.Status},
			at,
		)})
		if err != nil {
			return SubscriptionProrationResult{}, err
		}
		return SubscriptionProrationResult{Subscription: saved}, nil
	}

	remaining, periodSeconds, periodOK := ProrationFactor(sub.CurrentPeriodStart, sub.CurrentPeriodEnd, at)

	// create_prorations: accumulate pending delta, no invoice.
	if behavior == "create_prorations" {
		delta := int64(0)
		if periodOK {
			delta = ProrateDelta(newDiscounted-oldDiscounted, remaining, periodSeconds)
		}
		if delta != 0 {
			existing := int64(0)
			if raw := strings.TrimSpace(updated.Metadata[MetadataPendingProrationAmount]); raw != "" {
				existing, _ = strconv.ParseInt(raw, 10, 64)
			}
			updated.Metadata[MetadataPendingProrationAmount] = strconv.FormatInt(existing+delta, 10)
			updated.Metadata[MetadataPendingProrationAt] = at.Format(time.RFC3339)
		}
		saved, err := s.repo.UpdateSubscription(ctx, updated, []TimelineEntry{portalTimeline(
			"stripe_compat_update_"+updated.ID+"_"+at.Format(time.RFC3339Nano),
			"customer.subscription.updated",
			"Stripe-compatible subscription updated with pending proration",
			updated,
			map[string]string{"source": "stripe_compat", "status": updated.Status, "pending_proration": updated.Metadata[MetadataPendingProrationAmount]},
			at,
		)})
		if err != nil {
			return SubscriptionProrationResult{}, err
		}
		return SubscriptionProrationResult{Subscription: saved}, nil
	}

	// always_invoice
	resetAnchor := anchor == "now"
	var (
		invoiceSubtotal  int64
		invoiceDiscount  int64
		invoiceBase      int64 // discounted amount before exclusive tax
		creditSubtotal   int64 // pre-discount unused old-cycle credit (anchor=now)
		prorationCredit  int64 // post-discount unused old-cycle credit (anchor=now)
		shouldInvoice    bool
	)
	if resetAnchor {
		// Full new cycle minus unused old-cycle credit.
		// Subtotal uses pre-discount credit so serialization base (subtotal-discount)
		// matches the tax base (invoiceBase); see Stripe amount identity.
		if periodOK {
			creditSubtotal = ProrateDelta(oldTotal, remaining, periodSeconds)
			prorationCredit = ProrateDelta(oldDiscounted, remaining, periodSeconds)
		}
		invoiceSubtotal = newTotal - creditSubtotal
		if invoiceSubtotal < 0 {
			invoiceSubtotal = 0
		}
		invoiceBase = newDiscounted - prorationCredit
		if invoiceBase < 0 {
			invoiceBase = 0
		}
		invoiceDiscount = invoiceSubtotal - invoiceBase
		if invoiceDiscount < 0 {
			invoiceDiscount = 0
		}
		// Charge only when there is a positive post-credit amount to bill.
		shouldInvoice = invoiceBase > 0
		// Always reset period when anchor=now (even if no invoice).
		periodEnd, err := s.nextPeriodEnd(ctx, req.NewItems, at)
		if err != nil {
			return SubscriptionProrationResult{}, err
		}
		updated.CurrentPeriodStart = at
		updated.CurrentPeriodEnd = periodEnd
	} else {
		// Mid-cycle proration delta only; period unchanged.
		if !periodOK {
			// No usable period → no invoice, still apply items.
			shouldInvoice = false
		} else {
			delta := ProrateDelta(newDiscounted-oldDiscounted, remaining, periodSeconds)
			subtotalDelta := ProrateDelta(newTotal-oldTotal, remaining, periodSeconds)
			if delta <= 0 {
				// Downgrade / zero: no invoice (billtap does not model credit balance).
				shouldInvoice = false
			} else {
				invoiceBase = delta
				invoiceSubtotal = subtotalDelta
				invoiceDiscount = subtotalDelta - delta
				if invoiceDiscount < 0 {
					invoiceDiscount = 0
				}
				shouldInvoice = true
			}
		}
	}

	if !shouldInvoice {
		saved, err := s.repo.UpdateSubscription(ctx, updated, []TimelineEntry{portalTimeline(
			"stripe_compat_update_"+updated.ID+"_"+at.Format(time.RFC3339Nano),
			"customer.subscription.updated",
			"Stripe-compatible subscription updated without proration invoice",
			updated,
			map[string]string{"source": "stripe_compat", "status": updated.Status},
			at,
		)})
		if err != nil {
			return SubscriptionProrationResult{}, err
		}
		return SubscriptionProrationResult{Subscription: saved}, nil
	}

	// Tax on discounted base (same rule as renewSubscription).
	tax := int64(0)
	invoiceTotal := invoiceBase
	if len(rates) > 0 {
		_, _, exclusiveTotal, taxTotal := ComputeTaxRateAmounts(invoiceBase, rates)
		tax = taxTotal
		invoiceTotal = invoiceBase + exclusiveTotal
	} else if automaticTax {
		tax = ExclusiveTaxAmount(invoiceBase, taxPercent)
		invoiceTotal = invoiceBase + tax
	}

	// Payment outcome uses the same source as renewal (sub / customer metadata).
	renewalOutcome := renewalOutcome(updated.Metadata)
	if renewalOutcome == "" && updated.CustomerID != "" {
		if customer, err := s.repo.GetCustomer(ctx, updated.CustomerID); err == nil {
			renewalOutcome = CustomerDefaultInvoiceOutcome(customer.Metadata)
		}
	}
	renewalFailed := renewalOutcome != ""

	// error_if_incomplete + failure: do not commit item change or invoice.
	if renewalFailed && paymentBehavior == "error_if_incomplete" {
		spec, ok := intentOutcomeSpec(renewalOutcome)
		if !ok {
			spec, _ = intentOutcomeSpec("card_declined")
		}
		return SubscriptionProrationResult{}, &PaymentFailureError{
			Code:        firstNonEmpty(spec.FailureCode, "card_declined"),
			DeclineCode: spec.DeclineCode,
			Message:     firstNonEmpty(spec.FailureMessage, "Your card was declined."),
		}
	}

	invoice := Invoice{
		ID:              id("in"),
		Object:          ObjectInvoice,
		CustomerID:      updated.CustomerID,
		SubscriptionID:  updated.ID,
		Status:          "paid",
		Currency:        currency,
		Subtotal:        invoiceSubtotal,
		DiscountAmount:  invoiceDiscount,
		Discounts:       discounts,
		AutomaticTax:    automaticTax && len(rates) == 0,
		DefaultTaxRates: rates,
		Tax:             tax,
		Total:           invoiceTotal,
		AmountDue:       0,
		AmountPaid:      invoiceTotal,
		AttemptCount:    1,
		CreatedAt:       at,
		Metadata: map[string]string{
			MetadataBillingReason: "subscription_update",
		},
	}
	if creditSubtotal > 0 {
		// Primary key is pre-discount credit (matches subtotal reduction).
		invoice.Metadata[MetadataProrationCredit] = strconv.FormatInt(creditSubtotal, 10)
		if prorationCredit != creditSubtotal {
			invoice.Metadata[MetadataProrationCreditDiscounted] = strconv.FormatInt(prorationCredit, 10)
		}
	} else if prorationCredit > 0 {
		invoice.Metadata[MetadataProrationCredit] = strconv.FormatInt(prorationCredit, 10)
	}
	if renewalFailed {
		invoice.Status = "open"
		invoice.AmountDue = invoiceTotal
		invoice.AmountPaid = 0
		nextPaymentAttempt := at.Add(24 * time.Hour)
		invoice.NextPaymentAttempt = &nextPaymentAttempt
		updated.Status = renewalFailureSubscriptionStatus(renewalOutcome)
		updated.Metadata["billtap_last_renewal_outcome"] = renewalOutcome
		updated.Metadata["billtap_next_retry_at"] = nextPaymentAttempt.Format(time.RFC3339Nano)
	}

	intent := PaymentIntent{
		ID:              id("pi"),
		Object:          ObjectPaymentIntent,
		CustomerID:      updated.CustomerID,
		InvoiceID:       invoice.ID,
		Amount:          invoiceTotal,
		Currency:        currency,
		Status:          "succeeded",
		CaptureMethod:   "automatic",
		PaymentMethodID: "pm_card_visa",
		CreatedAt:       at,
	}
	if renewalFailed {
		spec, ok := intentOutcomeSpec(renewalOutcome)
		if !ok {
			spec, _ = intentOutcomeSpec("card_declined")
		}
		intent.Status = spec.PaymentIntentStatus
		intent.PaymentMethodID = firstNonEmpty(spec.PaymentMethodID, "pm_card_declined")
		intent.FailureCode = spec.FailureCode
		intent.DeclineCode = spec.DeclineCode
		intent.FailureMessage = spec.FailureMessage
	}
	invoice.PaymentIntentID = intent.ID
	updated.LatestInvoiceID = invoice.ID

	source := "subscription_update"
	timeline := []TimelineEntry{
		billingTimelineEntry("proration_invoice_created_"+invoice.ID, "invoice.created", "Subscription update invoice created", ObjectInvoice, invoice.ID, invoice.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status}, at),
		billingTimelineEntry("proration_invoice_finalized_"+invoice.ID, "invoice.finalized", "Subscription update invoice finalized", ObjectInvoice, invoice.ID, invoice.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status}, at),
		billingTimelineEntry("proration_payment_intent_created_"+intent.ID, "payment_intent.created", "Subscription update payment intent created", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": intent.Status}, at),
	}
	if renewalFailed {
		timeline = append(timeline,
			billingTimelineEntry("proration_payment_intent_failed_"+intent.ID, paymentIntentEvent(intent.Status), "Subscription update payment intent failed", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": intent.Status, "outcome": renewalOutcome}, at),
			billingTimelineEntry("proration_invoice_payment_failed_"+invoice.ID, "invoice.payment_failed", "Subscription update invoice payment failed", ObjectInvoice, invoice.ID, invoice.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status, "outcome": renewalOutcome}, at),
			billingTimelineEntry("proration_subscription_updated_"+updated.ID+"_"+invoice.ID, "customer.subscription.updated", "Subscription updated after proration payment failure", ObjectSubscription, updated.ID, updated.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": updated.Status, "outcome": renewalOutcome}, at),
		)
	} else {
		timeline = append(timeline,
			billingTimelineEntry("proration_payment_intent_succeeded_"+intent.ID, "payment_intent.succeeded", "Subscription update payment intent succeeded", ObjectPaymentIntent, intent.ID, intent.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": intent.Status}, at),
			billingTimelineEntry("proration_invoice_payment_succeeded_"+invoice.ID, "invoice.payment_succeeded", "Subscription update invoice payment succeeded", ObjectInvoice, invoice.ID, invoice.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status}, at),
			billingTimelineEntry("proration_invoice_paid_"+invoice.ID, "invoice.paid", "Subscription update invoice paid", ObjectInvoice, invoice.ID, invoice.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status}, at),
			billingTimelineEntry("proration_subscription_updated_"+updated.ID+"_"+invoice.ID, "customer.subscription.updated", "Subscription updated with proration invoice", ObjectSubscription, updated.ID, updated.CustomerID, "", updated.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": updated.Status}, at),
		)
	}

	savedSub, savedInvoice, savedIntent, err := s.repo.RecordSubscriptionRenewal(ctx, updated, invoice, intent, timeline)
	if err != nil {
		return SubscriptionProrationResult{}, err
	}
	paymentResult := InvoicePaymentResult{Invoice: savedInvoice, Subscription: savedSub, PaymentIntent: savedIntent}
	return SubscriptionProrationResult{
		Subscription:  savedSub,
		Invoice:       &savedInvoice,
		PaymentIntent: &savedIntent,
		PaymentResult: paymentResult,
	}, nil
}

func (s *Service) subscriptionAttachedToClock(ctx context.Context, sub Subscription, clockID string) bool {
	clockID = strings.TrimSpace(clockID)
	if clockID == "" {
		return true
	}
	for _, key := range []string{"test_clock", "testClock"} {
		if strings.TrimSpace(sub.Metadata[key]) == clockID {
			return true
		}
	}
	customer, err := s.repo.GetCustomer(ctx, sub.CustomerID)
	if err != nil {
		return false
	}
	for _, key := range []string{"test_clock", "testClock"} {
		if strings.TrimSpace(customer.Metadata[key]) == clockID {
			return true
		}
	}
	return false
}

func (s *Service) refundAttachedToClock(ctx context.Context, refund Refund, clockID string) bool {
	clockID = strings.TrimSpace(clockID)
	if clockID == "" {
		return true
	}
	if strings.TrimSpace(refund.Metadata["test_clock"]) == clockID || strings.TrimSpace(refund.Metadata["testClock"]) == clockID {
		return true
	}
	if refund.CustomerID == "" {
		return false
	}
	customer, err := s.repo.GetCustomer(ctx, refund.CustomerID)
	if err != nil {
		return false
	}
	return strings.TrimSpace(customer.Metadata["test_clock"]) == clockID || strings.TrimSpace(customer.Metadata["testClock"]) == clockID
}

func (s *Service) settlePendingRefunds(ctx context.Context, at time.Time, testClockID string) ([]Refund, error) {
	refunds, err := s.repo.ListRefundsFiltered(ctx, RefundFilter{})
	if err != nil {
		return nil, err
	}
	var settled []Refund
	for _, refund := range refunds {
		if refund.Status != "pending" || !s.refundAttachedToClock(ctx, refund, testClockID) {
			continue
		}
		settleAtRaw := firstNonEmpty(refund.Metadata["billtap_settle_at"], refund.Metadata["settle_at"], refund.Metadata["available_on"])
		if settleAtRaw == "" {
			continue
		}
		settleAt, err := parseMetadataTime(settleAtRaw)
		if err != nil {
			return settled, err
		}
		if settleAt.After(at) {
			continue
		}
		updated, err := s.UpdateRefundStatus(ctx, refund.ID, "succeeded", settleAt)
		if err != nil {
			return settled, err
		}
		settled = append(settled, updated)
	}
	return settled, nil
}

func renewalOutcome(metadata map[string]string) string {
	for _, key := range []string{"billtap_renewal_outcome", "renewal_outcome", "renewalOutcome"} {
		if value := normalizeInvoiceOutcome(metadata[key]); value != "" {
			return value
		}
	}
	return ""
}

// CustomerDefaultInvoiceOutcome returns the default renewal invoice outcome from customer metadata.
func CustomerDefaultInvoiceOutcome(metadata map[string]string) string {
	for _, key := range []string{MetadataDefaultInvoiceOutcome, "billtap_default_renewal_outcome", "default_invoice_outcome", "default_renewal_outcome"} {
		if value := normalizeInvoiceOutcome(metadata[key]); value != "" {
			return value
		}
	}
	return ""
}

func IsSupportedInvoiceOutcome(outcome string) bool {
	outcome = normalizeInvoiceOutcome(outcome)
	if outcome == "" {
		return true
	}
	return IsSupportedPaymentIntentOutcome(outcome)
}

func normalizeInvoiceOutcome(outcome string) string {
	value := strings.ToLower(strings.TrimSpace(outcome))
	switch value {
	case "", "payment_succeeded", "succeeded", "success":
		return ""
	case "payment_failed":
		return "card_declined"
	default:
		return value
	}
}

func parseMetadataTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("%w: time is required", ErrInvalidInput)
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(seconds, 0).UTC(), nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339Nano, value)
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: invalid time %q", ErrInvalidInput, value)
	}
	return parsed.UTC(), nil
}

func renewalFailureSubscriptionStatus(outcome string) string {
	switch strings.ToLower(strings.TrimSpace(outcome)) {
	case "unpaid":
		return "unpaid"
	default:
		return "past_due"
	}
}

func (s *Service) subscriptionTotal(ctx context.Context, items []LineItem) (int64, string, error) {
	total, currency, _, err := s.subscriptionLineAmounts(ctx, items)
	return total, currency, err
}

func (s *Service) subscriptionLineAmounts(ctx context.Context, items []LineItem) (int64, string, []LineAmount, error) {
	total := int64(0)
	currency := "usd"
	lineAmounts := make([]LineAmount, 0, len(items))
	for _, item := range items {
		price, err := s.repo.GetPrice(ctx, item.PriceID)
		if err != nil {
			return 0, "", nil, err
		}
		if price.Currency != "" {
			currency = price.Currency
		}
		quantity := item.Quantity
		if quantity <= 0 {
			quantity = 1
		}
		amount := price.UnitAmount * quantity
		total += amount
		lineAmounts = append(lineAmounts, LineAmount{ProductID: price.ProductID, Amount: amount})
	}
	return total, currency, lineAmounts, nil
}

func normalizeDiscounts(discounts []Discount, now time.Time) []Discount {
	if len(discounts) == 0 {
		return nil
	}
	out := make([]Discount, 0, len(discounts))
	for _, discount := range discounts {
		discount.CouponID = strings.TrimSpace(discount.CouponID)
		discount.PromotionCodeID = strings.TrimSpace(discount.PromotionCodeID)
		if discount.CouponID == "" && discount.PromotionCodeID == "" && discount.PercentOff == 0 && discount.AmountOff == 0 {
			continue
		}
		if discount.Object == "" {
			discount.Object = "discount"
		}
		if discount.ID == "" {
			source := firstNonEmpty(discount.PromotionCodeID, discount.CouponID, strconv.FormatInt(now.UnixNano(), 36))
			discount.ID = "di_" + sanitizeDiscountID(source)
		}
		if discount.Duration == "" {
			discount.Duration = "once"
		}
		if discount.CreatedAt.IsZero() {
			discount.CreatedAt = now
		}
		if discount.Metadata == nil {
			discount.Metadata = map[string]string{}
		}
		out = append(out, discount)
	}
	return out
}

func sanitizeDiscountID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "discount"
	}
	var b strings.Builder
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "discount"
	}
	return b.String()
}

// ApplyDiscounts applies Billtap's bounded Stripe-style discount subset.
// It supports a single effective coupon or promotion-code discount and never
// lets the invoice total become negative. Product-scoped coupons should use
// ApplyDiscountsWithEligibleBase with an eligible base computed from matching
// line items; this helper treats the full subtotal as eligible.
func ApplyDiscounts(subtotal int64, currency string, discounts []Discount) (int64, int64) {
	return ApplyDiscountsWithEligibleBase(subtotal, subtotal, currency, discounts)
}

// EligibleDiscountBase returns the amount a product-scoped discount may apply
// to. When the effective discount has no product restriction, it returns
// subtotal unchanged.
func EligibleDiscountBase(subtotal int64, discounts []Discount, lines []LineAmount) int64 {
	if len(discounts) == 0 || len(discounts[0].AppliesToProducts) == 0 {
		return subtotal
	}
	allowed := make(map[string]struct{}, len(discounts[0].AppliesToProducts))
	for _, productID := range discounts[0].AppliesToProducts {
		productID = strings.TrimSpace(productID)
		if productID != "" {
			allowed[productID] = struct{}{}
		}
	}
	if len(allowed) == 0 {
		return subtotal
	}
	var base int64
	for _, line := range lines {
		if _, ok := allowed[strings.TrimSpace(line.ProductID)]; ok {
			base += line.Amount
		}
	}
	return base
}

// ApplyDiscountsWithEligibleBase applies a single effective discount against
// eligibleBase (product-scoped or full subtotal) and never lets the total go
// negative relative to subtotal.
func ApplyDiscountsWithEligibleBase(subtotal int64, eligibleBase int64, currency string, discounts []Discount) (int64, int64) {
	if subtotal <= 0 || len(discounts) == 0 {
		return subtotal, 0
	}
	if eligibleBase < 0 {
		eligibleBase = 0
	}
	if eligibleBase > subtotal {
		eligibleBase = subtotal
	}
	discount := discounts[0]
	amount := int64(0)
	if discount.PercentOff > 0 {
		percentOff := discount.PercentOff
		if percentOff > 100 {
			percentOff = 100
		}
		// Stripe-compatible: round half away from zero (math.Round).
		// Integer percents that divide evenly keep the prior integer-division result.
		amount = int64(math.Round(float64(eligibleBase) * percentOff / 100.0))
	} else if discount.AmountOff > 0 {
		if discount.Currency == "" || strings.EqualFold(discount.Currency, currency) {
			amount = discount.AmountOff
			if amount > eligibleBase {
				amount = eligibleBase
			}
		}
	}
	if amount > subtotal {
		amount = subtotal
	}
	return subtotal - amount, amount
}

func MergeDiscountMetadata(metadata map[string]string, discounts []Discount) map[string]string {
	if len(discounts) == 0 {
		return metadata
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	discount := normalizeDiscounts(discounts, time.Now().UTC())[0]
	metadata[MetadataDiscountCouponID] = discount.CouponID
	metadata[MetadataDiscountPromotionCodeID] = discount.PromotionCodeID
	metadata[MetadataDiscountPercentOff] = strconv.FormatFloat(discount.PercentOff, 'f', -1, 64)
	metadata[MetadataDiscountAmountOff] = strconv.FormatInt(discount.AmountOff, 10)
	metadata[MetadataDiscountCurrency] = strings.ToLower(discount.Currency)
	metadata[MetadataDiscountDuration] = discount.Duration
	metadata[MetadataDiscountCreated] = discount.CreatedAt.Format(time.RFC3339Nano)
	if len(discount.AppliesToProducts) > 0 {
		metadata[MetadataDiscountAppliesTo] = strings.Join(discount.AppliesToProducts, ",")
	} else {
		delete(metadata, MetadataDiscountAppliesTo)
	}
	return metadata
}

// ParseCustomerTaxPercent reads customer metadata tax_percent. Invalid or
// negative values are treated as 0.
func ParseCustomerTaxPercent(metadata map[string]string) float64 {
	if metadata == nil {
		return 0
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(metadata["tax_percent"]), 64)
	if err != nil || value < 0 {
		return 0
	}
	return value
}

// ExclusiveTaxAmount applies exclusive tax after discounts.
func ExclusiveTaxAmount(amountAfterDiscount int64, taxPercent float64) int64 {
	if amountAfterDiscount <= 0 || taxPercent <= 0 {
		return 0
	}
	return int64(math.Round(float64(amountAfterDiscount) * taxPercent / 100.0))
}

// TaxRateAmount is one computed tax line.
type TaxRateAmount struct {
	Rate   AppliedTaxRate
	Amount int64
}

// ComputeTaxRateAmounts applies rate snapshots to a discounted base.
// pretax = base - sum(inclusive amounts); total = base + sum(exclusive amounts).
// Amounts are returned in the same order as rates.
func ComputeTaxRateAmounts(base int64, rates []AppliedTaxRate) (amounts []TaxRateAmount, pretax, exclusiveTotal, taxTotal int64) {
	if base < 0 {
		base = 0
	}
	if len(rates) == 0 {
		return nil, base, 0, 0
	}
	inclusiveSum := 0.0
	for _, rate := range rates {
		if rate.Inclusive && rate.Percentage > 0 {
			inclusiveSum += rate.Percentage
		}
	}
	// First pass: inclusive amounts (base already includes them).
	inclusiveByIndex := make([]int64, len(rates))
	inclusiveTotal := int64(0)
	for i, rate := range rates {
		if !rate.Inclusive {
			continue
		}
		amount := int64(0)
		if base > 0 && rate.Percentage > 0 && inclusiveSum > 0 {
			amount = int64(math.Round(float64(base) * rate.Percentage / (100.0 + inclusiveSum)))
		}
		inclusiveByIndex[i] = amount
		inclusiveTotal += amount
	}
	pretax = base - inclusiveTotal
	if pretax < 0 {
		pretax = 0
	}
	// Second pass: exclusive amounts on pretax; emit full list in rate order.
	amounts = make([]TaxRateAmount, 0, len(rates))
	for i, rate := range rates {
		amount := inclusiveByIndex[i]
		if !rate.Inclusive {
			if pretax > 0 && rate.Percentage > 0 {
				amount = int64(math.Round(float64(pretax) * rate.Percentage / 100.0))
			} else {
				amount = 0
			}
			exclusiveTotal += amount
		}
		amounts = append(amounts, TaxRateAmount{Rate: rate, Amount: amount})
		taxTotal += amount
	}
	return amounts, pretax, exclusiveTotal, taxTotal
}

// MergeTaxMetadata snapshots automatic tax flags onto subscription metadata.
func MergeTaxMetadata(metadata map[string]string, automaticTax bool, taxPercent float64) map[string]string {
	if !automaticTax {
		return metadata
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata[MetadataAutomaticTax] = "true"
	metadata[MetadataTaxPercent] = strconv.FormatFloat(taxPercent, 'f', -1, 64)
	return metadata
}

// AutomaticTaxFromMetadata restores the automatic-tax snapshot from subscription metadata.
func AutomaticTaxFromMetadata(metadata map[string]string) (bool, float64) {
	if metadata == nil {
		return false, 0
	}
	if strings.TrimSpace(metadata[MetadataAutomaticTax]) != "true" {
		return false, 0
	}
	percent, err := strconv.ParseFloat(strings.TrimSpace(metadata[MetadataTaxPercent]), 64)
	if err != nil || percent < 0 {
		percent = 0
	}
	return true, percent
}

// MergeDefaultTaxRatesMetadata snapshots applied default tax rates onto subscription metadata.
// An empty rates slice removes any existing snapshot.
func MergeDefaultTaxRatesMetadata(metadata map[string]string, rates []AppliedTaxRate) map[string]string {
	if len(rates) == 0 {
		if metadata == nil {
			return metadata
		}
		delete(metadata, MetadataDefaultTaxRates)
		return metadata
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	raw, err := json.Marshal(rates)
	if err != nil {
		return metadata
	}
	metadata[MetadataDefaultTaxRates] = string(raw)
	return metadata
}

// DefaultTaxRatesFromMetadata restores applied default tax rate snapshots from subscription metadata.
func DefaultTaxRatesFromMetadata(metadata map[string]string) []AppliedTaxRate {
	if metadata == nil {
		return nil
	}
	raw := strings.TrimSpace(metadata[MetadataDefaultTaxRates])
	if raw == "" {
		return nil
	}
	var rates []AppliedTaxRate
	if err := json.Unmarshal([]byte(raw), &rates); err != nil {
		return nil
	}
	if len(rates) == 0 {
		return nil
	}
	return rates
}

func ClearDiscountMetadata(metadata map[string]string) map[string]string {
	if metadata == nil {
		return metadata
	}
	for _, key := range []string{
		MetadataDiscountCouponID,
		MetadataDiscountPromotionCodeID,
		MetadataDiscountPercentOff,
		MetadataDiscountAmountOff,
		MetadataDiscountCurrency,
		MetadataDiscountDuration,
		MetadataDiscountCreated,
		MetadataDiscountAppliesTo,
	} {
		delete(metadata, key)
	}
	return metadata
}

func DiscountsFromMetadata(metadata map[string]string) []Discount {
	if metadata == nil {
		return nil
	}
	couponID := strings.TrimSpace(metadata[MetadataDiscountCouponID])
	promotionCodeID := strings.TrimSpace(metadata[MetadataDiscountPromotionCodeID])
	percentOff, _ := strconv.ParseFloat(strings.TrimSpace(metadata[MetadataDiscountPercentOff]), 64)
	amountOff, _ := strconv.ParseInt(strings.TrimSpace(metadata[MetadataDiscountAmountOff]), 10, 64)
	if couponID == "" && promotionCodeID == "" && percentOff == 0 && amountOff == 0 {
		return nil
	}
	createdAt, _ := parseMetadataTime(metadata[MetadataDiscountCreated])
	var appliesTo []string
	if raw := strings.TrimSpace(metadata[MetadataDiscountAppliesTo]); raw != "" {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				appliesTo = append(appliesTo, part)
			}
		}
	}
	return normalizeDiscounts([]Discount{{
		CouponID:          couponID,
		PromotionCodeID:   promotionCodeID,
		PercentOff:        percentOff,
		AmountOff:         amountOff,
		Currency:          strings.ToLower(strings.TrimSpace(metadata[MetadataDiscountCurrency])),
		Duration:          firstNonEmpty(metadata[MetadataDiscountDuration], "once"),
		AppliesToProducts: appliesTo,
		CreatedAt:         createdAt,
	}}, time.Now().UTC())
}

func (s *Service) nextPeriodEnd(ctx context.Context, items []LineItem, start time.Time) (time.Time, error) {
	if len(items) == 0 {
		return start.AddDate(0, 1, 0), nil
	}
	price, err := s.repo.GetPrice(ctx, items[0].PriceID)
	if err != nil {
		return time.Time{}, err
	}
	count := price.RecurringIntervalCount
	if count <= 0 {
		count = 1
	}
	switch price.RecurringInterval {
	case "day":
		return start.AddDate(0, 0, count), nil
	case "week":
		return start.AddDate(0, 0, 7*count), nil
	case "year":
		return start.AddDate(count, 0, 0), nil
	default:
		return start.AddDate(0, count, 0), nil
	}
}

// NextPeriodEnd is the public wrapper for next billing-period end calculation.
// API preview uses this so period math stays in the billing service (no API-layer duplication).
func (s *Service) NextPeriodEnd(ctx context.Context, items []LineItem, start time.Time) (time.Time, error) {
	return s.nextPeriodEnd(ctx, items, start)
}

func invoicePaymentTimeline(sub Subscription, invoice Invoice, intent PaymentIntent, success bool, at time.Time) []TimelineEntry {
	source := "invoice.pay"
	attemptSuffix := fmt.Sprintf("_attempt_%d_%s", invoice.AttemptCount, at.Format(time.RFC3339Nano))
	entries := []TimelineEntry{
		billingTimelineEntry("invoice_pay_payment_intent_"+intent.ID+attemptSuffix, paymentIntentEvent(intent.Status), "Invoice payment intent "+intent.Status, ObjectPaymentIntent, intent.ID, intent.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": intent.Status}, at),
	}
	if success {
		entries = append(entries,
			billingTimelineEntry("invoice_pay_succeeded_"+invoice.ID+attemptSuffix, "invoice.payment_succeeded", "Invoice payment succeeded", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status}, at),
			billingTimelineEntry("invoice_paid_"+invoice.ID+attemptSuffix, "invoice.paid", "Invoice paid", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status}, at),
		)
	} else {
		entries = append(entries, billingTimelineEntry("invoice_pay_failed_"+invoice.ID+attemptSuffix, "invoice.payment_failed", "Invoice payment failed", ObjectInvoice, invoice.ID, invoice.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": invoice.Status}, at))
	}
	entries = append(entries, billingTimelineEntry("invoice_pay_subscription_"+sub.ID+attemptSuffix, "customer.subscription.updated", "Subscription updated after invoice payment attempt", ObjectSubscription, sub.ID, sub.CustomerID, "", sub.ID, invoice.ID, intent.ID, map[string]string{"source": source, "status": sub.Status}, at))
	return entries
}

func billingTimelineEntry(seed, action, message, objectType, objectID, customerID, checkoutSessionID, subscriptionID, invoiceID, paymentIntentID string, data map[string]string, at time.Time) TimelineEntry {
	return TimelineEntry{
		ID:                "tl_" + sanitizeID(seed),
		Object:            ObjectTimelineEntry,
		Action:            action,
		Message:           message,
		ObjectType:        objectType,
		ObjectID:          objectID,
		CustomerID:        customerID,
		CheckoutSessionID: checkoutSessionID,
		SubscriptionID:    subscriptionID,
		InvoiceID:         invoiceID,
		PaymentIntentID:   paymentIntentID,
		Data:              compactMetadata(data),
		CreatedAt:         at,
	}
}

// CheckoutCompletion records terminal checkout state.
// Subscription and Invoice use zero-value (empty ID) when absent — payment mode
// omits both; free payment also omits PaymentIntent. Callers and storage skip
// INSERT/timeline when the relevant ID is empty (avoids pointer churn at use sites).
type CheckoutCompletion struct {
	SessionID     string
	SessionStatus string
	PaymentStatus string
	CheckoutEvent string
	Outcome       string
	CompletedAt   time.Time
	Subscription  Subscription
	Invoice       Invoice
	PaymentIntent PaymentIntent
}

type CheckoutCompletionOptions struct {
	SubscriptionID  string
	InvoiceID       string
	PaymentIntentID string
	At              time.Time
}

type TimelineFilter struct {
	CustomerID        string
	CheckoutSessionID string
	SubscriptionID    string
	InvoiceID         string
	PaymentIntentID   string
	ObjectType        string
	ObjectID          string
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func currentPortalSubscription(subscriptions []Subscription) *Subscription {
	for idx := range subscriptions {
		if subscriptions[idx].Status == "active" || subscriptions[idx].Status == "trialing" {
			return &subscriptions[idx]
		}
	}
	if len(subscriptions) == 0 {
		return nil
	}
	return &subscriptions[0]
}

func portalSummary(customerID string, subscription *Subscription, invoices []Invoice, paymentIntents []PaymentIntent) PortalStateSummary {
	summary := PortalStateSummary{
		CustomerID:         customerID,
		InvoiceCount:       len(invoices),
		PaymentIntentCount: len(paymentIntents),
	}
	for _, invoice := range invoices {
		if invoice.Status == "open" {
			summary.OpenInvoiceCount++
		}
	}
	if subscription == nil {
		return summary
	}
	periodEnd := subscription.CurrentPeriodEnd
	summary.SubscriptionID = subscription.ID
	summary.SubscriptionStatus = subscription.Status
	summary.Active = subscription.Status == "active" || subscription.Status == "trialing"
	summary.PendingCancellation = subscription.CancelAtPeriodEnd
	summary.CancelAtPeriodEnd = subscription.CancelAtPeriodEnd
	summary.CurrentPeriodEnd = &periodEnd
	summary.LatestInvoiceID = subscription.LatestInvoiceID
	return summary
}

func invoiceIDs(invoices []Invoice) []string {
	out := make([]string, 0, len(invoices))
	for _, invoice := range invoices {
		if invoice.ID != "" {
			out = append(out, invoice.ID)
		}
	}
	return out
}

func firstItem(items []LineItem) (string, int64) {
	if len(items) == 0 {
		return "", 0
	}
	return items[0].PriceID, items[0].Quantity
}

func portalTimeline(seed, action, message string, sub Subscription, data map[string]string, at time.Time) TimelineEntry {
	return TimelineEntry{
		ID:             "tl_" + sanitizeID(seed),
		Object:         ObjectTimelineEntry,
		Action:         action,
		Message:        message,
		ObjectType:     ObjectSubscription,
		ObjectID:       sub.ID,
		CustomerID:     sub.CustomerID,
		SubscriptionID: sub.ID,
		InvoiceID:      sub.LatestInvoiceID,
		Data:           compactMetadata(data),
		CreatedAt:      at,
	}
}

func ensurePaymentIntentConfirmable(intent PaymentIntent) error {
	switch intent.Status {
	case "", "requires_payment_method", "requires_action":
		return nil
	default:
		return fmt.Errorf("%w: status must be requires_payment_method", ErrInvalidInput)
	}
}

func timelineEntry(seed, action, message, objectType, objectID, customerID, checkoutSessionID, subscriptionID, paymentIntentID string, data map[string]string, at time.Time) TimelineEntry {
	return TimelineEntry{
		ID:                "tl_" + sanitizeID(seed),
		Object:            ObjectTimelineEntry,
		Action:            action,
		Message:           message,
		ObjectType:        objectType,
		ObjectID:          objectID,
		CustomerID:        customerID,
		CheckoutSessionID: checkoutSessionID,
		SubscriptionID:    subscriptionID,
		PaymentIntentID:   paymentIntentID,
		Data:              compactMetadata(data),
		CreatedAt:         at,
	}
}

func copyMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		out[key] = value
	}
	return out
}

func compactMetadata(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		if strings.TrimSpace(value) != "" {
			out[key] = value
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func normalizePaymentMethodOutcome(outcome string) string {
	switch strings.ToLower(strings.TrimSpace(outcome)) {
	case "", "success", "succeeded", "succeeds":
		return "succeeds"
	case "failure", "failed", "fails", "card_declined":
		return "fails"
	default:
		return outcome
	}
}

func sanitizeID(raw string) string {
	replacer := strings.NewReplacer("/", "_", ".", "_", " ", "_", ":", "_", "+", "_")
	return replacer.Replace(raw)
}

func id(prefix string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%s_%d", prefix, time.Now().UTC().UnixNano())
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}

type checkoutOutcomeSpec struct {
	Outcome             string
	Paid                bool
	SessionStatus       string
	PaymentStatus       string
	CheckoutEvent       string
	SubscriptionStatus  string
	InvoiceStatus       string
	InvoiceAttemptCount *int
	NextPaymentAttempt  bool
	PaymentIntentStatus string
	PaymentMethodID     string
	FailureCode         string
	DeclineCode         string
	FailureMessage      string
}

func checkoutOutcomeFor(outcome string) (checkoutOutcomeSpec, bool) {
	normalized := strings.ToLower(strings.TrimSpace(outcome))
	switch normalized {
	case "", "success", "payment_succeeded", "succeeded", "paid", "pm_card_visa":
		return checkoutOutcomeSpec{
			Outcome:             "success",
			Paid:                true,
			PaymentIntentStatus: "succeeded",
			PaymentMethodID:     paymentMethodID(normalized),
		}, true
	case "failure", "payment_failed", "failed", "card_declined", "generic_decline", "pm_card_visa_chargedeclined":
		return failedCheckoutOutcome("card_declined", paymentMethodID(normalized), "card_declined", "generic_decline", "Your card was declined."), true
	case "insufficient_funds", "pm_card_visa_chargedeclinedinsufficientfunds":
		return failedCheckoutOutcome("insufficient_funds", paymentMethodID(normalized), "card_declined", "insufficient_funds", "Your card has insufficient funds."), true
	case "customer_payment_method_failed", "pm_card_chargecustomerfail":
		return failedCheckoutOutcome("customer_payment_method_failed", paymentMethodID(normalized), "card_declined", "do_not_honor", "The customer's payment method was declined."), true
	case "expired_card":
		return failedCheckoutOutcome("expired_card", "", "expired_card", "expired_card", "Your card has expired."), true
	case "incorrect_cvc":
		return failedCheckoutOutcome("incorrect_cvc", "", "incorrect_cvc", "incorrect_cvc", "Your card's security code is incorrect."), true
	case "processing_error":
		return failedCheckoutOutcome("processing_error", "", "processing_error", "processing_error", "An error occurred while processing your card. Try again later."), true
	case "missing_payment_method", "payment_method_missing":
		return failedCheckoutOutcome("missing_payment_method", "", "payment_method_missing", "", "No payment method is available to complete this checkout."), true
	case "authentication_required", "requires_action", "pm_card_threedsecure2required":
		methodID := paymentMethodID(normalized)
		if methodID == "" {
			methodID = "pm_card_threeDSecure2Required"
		}
		spec := failedCheckoutOutcome("authentication_required", methodID, "authentication_required", "authentication_required", "This payment requires authentication.")
		spec.PaymentIntentStatus = "requires_action"
		return spec, true
	case "payment_pending", "pending", "processing", "async_payment_pending":
		return checkoutOutcomeSpec{
			Outcome:             "payment_pending",
			Paid:                false,
			PaymentStatus:       "unpaid",
			SubscriptionStatus:  "incomplete",
			InvoiceStatus:       "open",
			PaymentIntentStatus: "processing",
			PaymentMethodID:     paymentMethodID(normalized),
		}, true
	case "bank_transfer", "bank_transfer_processing", "pm_bank_transfer":
		return checkoutOutcomeSpec{
			Outcome:             "bank_transfer",
			Paid:                false,
			PaymentStatus:       "unpaid",
			SubscriptionStatus:  "incomplete",
			InvoiceStatus:       "open",
			PaymentIntentStatus: "processing",
			PaymentMethodID:     "pm_bank_transfer",
		}, true
	case "canceled", "cancelled", "cancel", "payment_canceled", "pm_card_chargecustomercancel":
		zeroAttempts := 0
		return checkoutOutcomeSpec{
			Outcome:             "canceled",
			Paid:                false,
			SessionStatus:       "expired",
			PaymentStatus:       "unpaid",
			CheckoutEvent:       "checkout.session.expired",
			SubscriptionStatus:  "incomplete_expired",
			InvoiceStatus:       "void",
			InvoiceAttemptCount: &zeroAttempts,
			PaymentIntentStatus: "canceled",
			PaymentMethodID:     paymentMethodID(normalized),
		}, true
	default:
		return checkoutOutcomeSpec{}, false
	}
}

func failedCheckoutOutcome(outcome string, paymentMethodID string, code string, declineCode string, message string) checkoutOutcomeSpec {
	return checkoutOutcomeSpec{
		Outcome:             outcome,
		Paid:                false,
		PaymentStatus:       "unpaid",
		SubscriptionStatus:  "incomplete",
		InvoiceStatus:       "open",
		NextPaymentAttempt:  true,
		PaymentIntentStatus: "requires_payment_method",
		PaymentMethodID:     paymentMethodID,
		FailureCode:         code,
		DeclineCode:         declineCode,
		FailureMessage:      message,
	}
}

func isBankTransferPaymentMethod(paymentMethodID string) bool {
	value := strings.ToLower(strings.TrimSpace(paymentMethodID))
	return value == "pm_bank_transfer" || strings.Contains(value, "bank_transfer")
}

func intentOutcomeSpec(outcome string) (checkoutOutcomeSpec, bool) {
	raw := strings.TrimSpace(outcome)
	if spec, ok := checkoutOutcomeFor(raw); ok {
		if spec.PaymentMethodID == "" && strings.HasPrefix(strings.ToLower(raw), "pm_") {
			spec.PaymentMethodID = raw
		}
		return spec, true
	}
	if strings.HasPrefix(strings.ToLower(raw), "pm_") {
		return checkoutOutcomeSpec{
			Outcome:             "success",
			Paid:                true,
			PaymentIntentStatus: "succeeded",
			PaymentMethodID:     raw,
		}, true
	}
	return checkoutOutcomeSpec{}, false
}

func paymentIntentConfiguredOutcome(metadata map[string]string) string {
	for _, key := range []string{MetadataPaymentIntentOutcome, "billtap_outcome"} {
		if value := strings.TrimSpace(metadata[key]); value != "" {
			return value
		}
	}
	return ""
}

// CustomerDefaultPaymentIntentOutcome returns the default direct PaymentIntent outcome from customer metadata.
func CustomerDefaultPaymentIntentOutcome(metadata map[string]string) string {
	for _, key := range []string{MetadataDefaultPaymentIntentOutcome, "default_payment_intent_outcome"} {
		if value := strings.TrimSpace(metadata[key]); value != "" {
			return value
		}
	}
	return ""
}

// IsSupportedPaymentIntentOutcome reports whether Billtap can apply outcome to a local PaymentIntent.
func IsSupportedPaymentIntentOutcome(outcome string) bool {
	_, ok := intentOutcomeSpec(outcome)
	return ok
}

func paymentIntentEvent(status string) string {
	switch status {
	case "succeeded":
		return "payment_intent.succeeded"
	case "processing":
		return "payment_intent.processing"
	case "canceled":
		return "payment_intent.canceled"
	case "requires_action":
		return "payment_intent.requires_action"
	case "requires_capture":
		return "payment_intent.amount_capturable_updated"
	default:
		return "payment_intent.payment_failed"
	}
}

func invoiceEventTypesForPayment(status string, paymentIntentStatus string) []string {
	switch status {
	case "paid":
		return []string{"invoice.payment_succeeded", "invoice.paid"}
	case "void":
		return []string{"invoice.voided"}
	case "open":
		if paymentIntentStatus == "processing" {
			return nil
		}
		return []string{"invoice.payment_failed"}
	default:
		return []string{"invoice.payment_failed"}
	}
}

func setupIntentEvent(status string) string {
	switch status {
	case "succeeded":
		return "setup_intent.succeeded"
	case "canceled":
		return "setup_intent.canceled"
	case "requires_action":
		return "setup_intent.requires_action"
	default:
		return "setup_intent.setup_failed"
	}
}

func paymentMethodID(normalizedOutcome string) string {
	switch normalizedOutcome {
	case "pm_card_visa":
		return "pm_card_visa"
	case "pm_card_visa_chargedeclined":
		return "pm_card_visa_chargeDeclined"
	case "pm_card_visa_chargedeclinedinsufficientfunds":
		return "pm_card_visa_chargeDeclinedInsufficientFunds"
	case "pm_card_chargecustomerfail":
		return "pm_card_chargeCustomerFail"
	case "pm_card_threedsecure2required":
		return "pm_card_threeDSecure2Required"
	default:
		return ""
	}
}
