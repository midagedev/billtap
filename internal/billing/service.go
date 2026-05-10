package billing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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
	GetPaymentIntent(context.Context, string) (PaymentIntent, error)
	CreatePaymentIntent(context.Context, PaymentIntent) (PaymentIntent, error)
	UpdatePaymentIntent(context.Context, PaymentIntent, []TimelineEntry) (PaymentIntent, error)
	ListPaymentIntents(context.Context) ([]PaymentIntent, error)
	ListPaymentIntentsFiltered(context.Context, PaymentIntentFilter) ([]PaymentIntent, error)
	CreateSetupIntent(context.Context, SetupIntent) (SetupIntent, error)
	GetSetupIntent(context.Context, string) (SetupIntent, error)
	UpdateSetupIntent(context.Context, SetupIntent, []TimelineEntry) (SetupIntent, error)
	ListSetupIntents(context.Context) ([]SetupIntent, error)

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
	Items             []LineItem
	ReplaceItems      bool
	Metadata          map[string]string
	CancelAtPeriodEnd *bool
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
			sub.Metadata[key] = value
		}
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
	return s.repo.UpdateSubscription(ctx, sub, []TimelineEntry{portalTimeline(
		"stripe_compat_update_"+sub.ID+"_"+now.Format(time.RFC3339Nano),
		"customer.subscription.updated",
		"Stripe-compatible subscription updated",
		sub,
		map[string]string{"source": "stripe_compat"},
		now,
	)})
}

func (s *Service) GetInvoice(ctx context.Context, id string) (Invoice, error) {
	return s.repo.GetInvoice(ctx, id)
}

func (s *Service) ListInvoices(ctx context.Context) ([]Invoice, error) {
	return s.repo.ListInvoices(ctx)
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
	if firstNonEmpty(paymentMethodID, intent.PaymentMethodID, outcome) == "" {
		return PaymentIntent{}, fmt.Errorf("%w: payment_method is required", ErrInvalidInput)
	}
	if paymentMethodID != "" {
		intent.PaymentMethodID = paymentMethodID
	}
	spec, ok := intentOutcomeSpec(firstNonEmpty(outcome, intent.PaymentMethodID))
	if !ok {
		return PaymentIntent{}, fmt.Errorf("%w: %s", ErrUnsupportedOutcome, outcome)
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
		timelineEntry("pi_"+intent.ID+"_confirmed_"+now.Format(time.RFC3339Nano), paymentIntentEvent(intent.Status), "Payment intent "+intent.Status, ObjectPaymentIntent, intent.ID, intent.CustomerID, "", "", intent.ID, map[string]string{"status": intent.Status}, now),
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

func (s *Service) SimulatePaymentMethodUpdate(ctx context.Context, customerID string, outcome string) (PaymentMethodSimulation, error) {
	if strings.TrimSpace(customerID) == "" {
		return PaymentMethodSimulation{}, fmt.Errorf("%w: customer is required", ErrInvalidInput)
	}
	if _, err := s.repo.GetCustomer(ctx, customerID); err != nil {
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
		result.PaymentMethodID = id("pm")
	} else {
		result.Status = "failed"
		result.FailureCode = "card_declined"
		result.FailureMessage = "Simulated payment method update failure"
		action = "payment_method.update_failed"
		message = "Portal payment method update failed"
	}
	err := s.repo.RecordTimeline(ctx, TimelineEntry{
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
}

type TimelineFilter struct {
	CustomerID        string
	CheckoutSessionID string
	SubscriptionID    string
	InvoiceID         string
	PaymentIntentID   string
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
		spec := failedCheckoutOutcome("authentication_required", paymentMethodID(normalized), "authentication_required", "authentication_required", "This payment requires authentication.")
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
