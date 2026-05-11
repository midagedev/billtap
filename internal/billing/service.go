package billing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnsupportedOutcome = errors.New("unsupported checkout outcome")
)

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
	RecordCheckoutCompletion(context.Context, CheckoutCompletion) (CheckoutSession, error)

	GetSubscription(context.Context, string) (Subscription, error)
	ListSubscriptions(context.Context) ([]Subscription, error)
	ListSubscriptionsByCustomer(context.Context, string) ([]Subscription, error)
	UpdateSubscription(context.Context, Subscription, []TimelineEntry) (Subscription, error)
	GetInvoice(context.Context, string) (Invoice, error)
	ListInvoices(context.Context) ([]Invoice, error)
	ListInvoicesFiltered(context.Context, InvoiceFilter) ([]Invoice, error)
	UpdateInvoicePayment(context.Context, Subscription, Invoice, PaymentIntent, []TimelineEntry) (Subscription, Invoice, PaymentIntent, error)
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
	if session.PaymentIntentID != "" {
		return session, nil
	}

	outcomeSpec, ok := checkoutOutcomeFor(outcome)
	if !ok {
		return CheckoutSession{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
	}

	total := int64(0)
	currency := "usd"
	for _, item := range session.LineItems {
		price, err := s.repo.GetPrice(ctx, item.PriceID)
		if err != nil {
			return CheckoutSession{}, err
		}
		if price.Currency != "" {
			currency = price.Currency
		}
		total += price.UnitAmount * item.Quantity
	}

	now := s.now()
	if !opts.At.IsZero() {
		now = opts.At
	}
	periodEnd := now.AddDate(0, 1, 0)
	paid := outcomeSpec.Paid
	trialing := paid && session.TrialPeriodDays > 0
	if trialing {
		periodEnd = now.AddDate(0, 0, int(session.TrialPeriodDays))
	}
	invoiceTotal := total
	if trialing {
		invoiceTotal = 0
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
	if trialing {
		sub.Status = "trialing"
		sub.Metadata["trial_period_days"] = fmt.Sprintf("%d", session.TrialPeriodDays)
		sub.Metadata["trial_start"] = now.Format(time.RFC3339Nano)
		sub.Metadata["trial_end"] = periodEnd.Format(time.RFC3339Nano)
	}
	invoice := Invoice{
		ID:             firstNonEmpty(opts.InvoiceID, id("in")),
		Object:         ObjectInvoice,
		CustomerID:     session.CustomerID,
		SubscriptionID: sub.ID,
		Status:         "paid",
		Currency:       currency,
		Subtotal:       invoiceTotal,
		Total:          invoiceTotal,
		AmountDue:      0,
		AmountPaid:     invoiceTotal,
		AttemptCount:   1,
		CreatedAt:      now,
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
		invoice.AmountDue = total
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

func (s *Service) ListInvoices(ctx context.Context) ([]Invoice, error) {
	return s.repo.ListInvoices(ctx)
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
	total, currency, err := s.subscriptionTotal(ctx, sub.Items)
	if err != nil {
		return InvoicePaymentResult{}, err
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
	sub.Metadata = copyMap(sub.Metadata)
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
		ID:             id("in"),
		Object:         ObjectInvoice,
		CustomerID:     sub.CustomerID,
		SubscriptionID: sub.ID,
		Status:         "paid",
		Currency:       currency,
		Subtotal:       total,
		Total:          total,
		AmountDue:      0,
		AmountPaid:     total,
		AttemptCount:   1,
		CreatedAt:      at,
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
	total := int64(0)
	currency := "usd"
	for _, item := range items {
		price, err := s.repo.GetPrice(ctx, item.PriceID)
		if err != nil {
			return 0, "", err
		}
		if price.Currency != "" {
			currency = price.Currency
		}
		quantity := item.Quantity
		if quantity <= 0 {
			quantity = 1
		}
		total += price.UnitAmount * quantity
	}
	return total, currency, nil
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
