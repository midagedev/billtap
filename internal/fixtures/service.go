package fixtures

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
)

var (
	ErrInvalidFixture  = errors.New("invalid fixture")
	ErrAssertionFailed = errors.New("fixture assertion failed")
)

type Service struct {
	billing *billing.Service
	now     func() time.Time
}

func NewService(billingService *billing.Service) *Service {
	return &Service{billing: billingService, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) Apply(ctx context.Context, pack Pack) (ApplyResult, error) {
	if s.billing == nil {
		return ApplyResult{}, fmt.Errorf("%w: billing service is required", ErrInvalidFixture)
	}
	pack = normalizePack(pack)
	if err := validatePack(pack); err != nil {
		return ApplyResult{}, err
	}

	result := ApplyResult{
		ID:        "fxapp_" + sanitizeID(pack.Name+"_"+pack.RunID+"_"+time.Now().UTC().Format(time.RFC3339Nano)),
		Object:    "fixture_apply_result",
		Name:      pack.Name,
		RunID:     pack.RunID,
		Namespace: pack.Namespace,
		AppliedAt: s.now(),
		Summary:   map[string]int{},
	}

	for _, fixture := range pack.Customers {
		customer, err := s.upsertCustomer(ctx, pack, fixture)
		if err != nil {
			return result, err
		}
		result.Customers = append(result.Customers, customer)
	}
	for _, fixture := range pack.Products {
		product, err := s.upsertProduct(ctx, pack, fixture)
		if err != nil {
			return result, err
		}
		result.Products = append(result.Products, product)
	}
	for _, fixture := range pack.Prices {
		price, err := s.upsertPrice(ctx, pack, fixture)
		if err != nil {
			return result, err
		}
		result.Prices = append(result.Prices, price)
	}
	for _, fixture := range pack.TestClocks {
		clock, err := s.upsertTestClock(ctx, fixture)
		if err != nil {
			return result, err
		}
		result.TestClocks = append(result.TestClocks, clock)
	}
	for _, fixture := range pack.Subscriptions {
		session, subscription, err := s.upsertSubscription(ctx, pack, fixture)
		if err != nil {
			return result, err
		}
		if session.ID != "" {
			result.CheckoutSessions = append(result.CheckoutSessions, session)
		}
		result.Subscriptions = append(result.Subscriptions, subscription)
	}
	for _, fixture := range pack.Refunds {
		refund, err := s.createRefund(ctx, pack, fixture)
		if err != nil {
			return result, err
		}
		result.Refunds = append(result.Refunds, refund)
	}
	for _, fixture := range pack.CreditNotes {
		note, err := s.createCreditNote(ctx, pack, fixture)
		if err != nil {
			return result, err
		}
		result.CreditNotes = append(result.CreditNotes, note)
	}

	result.Summary = map[string]int{
		"customers":         len(result.Customers),
		"products":          len(result.Products),
		"prices":            len(result.Prices),
		"test_clocks":       len(result.TestClocks),
		"checkout_sessions": len(result.CheckoutSessions),
		"subscriptions":     len(result.Subscriptions),
		"refunds":           len(result.Refunds),
		"credit_notes":      len(result.CreditNotes),
	}

	if len(pack.Assertions) > 0 {
		report, err := s.Assert(ctx, AssertionRequest{
			Name:   pack.Name,
			Filter: pack.defaultFilter(),
			Expect: pack.Assertions,
		})
		result.Assertions = &report
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (s *Service) Snapshot(ctx context.Context, filter SnapshotFilter) (Snapshot, error) {
	if s.billing == nil {
		return Snapshot{}, fmt.Errorf("%w: billing service is required", ErrInvalidFixture)
	}
	filter = normalizeFilter(filter)
	customers, err := s.billing.ListCustomers(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	products, err := s.billing.ListProducts(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	prices, err := s.billing.ListPrices(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	sessions, err := s.billing.ListCheckoutSessions(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	subscriptions, err := s.billing.ListSubscriptions(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	invoices, err := s.billing.ListInvoices(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	paymentIntents, err := s.billing.ListPaymentIntents(ctx)
	if err != nil {
		return Snapshot{}, err
	}

	filteredCustomers := filterCustomers(customers, filter)
	customerIDs := customerSet(filteredCustomers)
	filteredProducts := filterProducts(products, filter)
	productIDs := productSet(filteredProducts)
	filteredPrices := filterPrices(prices, filter, productIDs)
	priceIDs := priceSet(filteredPrices)
	filteredSubscriptions := filterSubscriptions(subscriptions, filter, customerIDs, priceIDs)
	subscriptionIDs := subscriptionSet(filteredSubscriptions)
	filteredSessions := filterCheckoutSessions(sessions, filter, customerIDs, subscriptionIDs)
	checkoutIDs := checkoutSessionSet(filteredSessions)
	filteredInvoices := filterInvoices(invoices, filter, customerIDs, subscriptionIDs)
	invoiceIDs := invoiceSet(filteredInvoices)
	filteredPaymentIntents := filterPaymentIntents(paymentIntents, filter, customerIDs, invoiceIDs)
	paymentIntentIDs := paymentIntentSet(filteredPaymentIntents)
	timeline, err := s.filteredTimeline(ctx, filter, customerIDs, checkoutIDs, subscriptionIDs, invoiceIDs, paymentIntentIDs)
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		Object:           "fixture_snapshot",
		Filter:           filter,
		Customers:        filteredCustomers,
		Products:         filteredProducts,
		Prices:           filteredPrices,
		CheckoutSessions: filteredSessions,
		Subscriptions:    filteredSubscriptions,
		Invoices:         filteredInvoices,
		PaymentIntents:   filteredPaymentIntents,
		Timeline:         timeline,
		Summary: map[string]int{
			"customers":         len(filteredCustomers),
			"products":          len(filteredProducts),
			"prices":            len(filteredPrices),
			"checkout_sessions": len(filteredSessions),
			"subscriptions":     len(filteredSubscriptions),
			"invoices":          len(filteredInvoices),
			"payment_intents":   len(filteredPaymentIntents),
			"timeline":          len(timeline),
		},
		CapturedAt: s.now(),
	}, nil
}

func (s *Service) Assert(ctx context.Context, req AssertionRequest) (AssertionReport, error) {
	snapshot, err := s.Snapshot(ctx, req.Filter)
	if err != nil {
		return AssertionReport{}, err
	}
	report := AssertionReport{
		Object:    "fixture_assertion_report",
		Name:      req.Name,
		Pass:      true,
		CheckedAt: s.now(),
	}
	for _, expectation := range req.Expect {
		result := evaluateExpectation(snapshot, expectation)
		if !result.Pass {
			report.Pass = false
		}
		report.Results = append(report.Results, result)
	}
	if !report.Pass {
		return report, ErrAssertionFailed
	}
	return report, nil
}

func (s *Service) Resolve(ctx context.Context, filter ResolveFilter) (ResolveResult, error) {
	ref := strings.TrimSpace(filter.Ref)
	if ref == "" {
		return ResolveResult{}, fmt.Errorf("%w: ref is required", ErrInvalidFixture)
	}
	snapshot, err := s.Snapshot(ctx, SnapshotFilter{
		RunID:       filter.RunID,
		FixtureName: filter.FixtureName,
		Namespace:   filter.Namespace,
		TenantID:    filter.TenantID,
	})
	if err != nil {
		return ResolveResult{}, err
	}
	result := ResolveResult{Object: "fixture_resolve", Ref: ref}
	for _, subscription := range snapshot.Subscriptions {
		if subscription.ID != ref && subscription.Metadata[MetadataFixtureRef] != ref {
			continue
		}
		result.SubscriptionID = subscription.ID
		result.CustomerID = subscription.CustomerID
		result.InvoiceID = subscription.LatestInvoiceID
		result.Metadata = subscription.Metadata
		if len(subscription.Items) > 0 {
			result.PriceID = subscription.Items[0].PriceID
		}
		break
	}
	for _, session := range snapshot.CheckoutSessions {
		if result.SubscriptionID != "" && session.SubscriptionID == result.SubscriptionID {
			result.CheckoutSessionID = session.ID
			break
		}
		if session.ID == ref {
			result.CheckoutSessionID = session.ID
			result.CustomerID = session.CustomerID
			result.SubscriptionID = session.SubscriptionID
			result.InvoiceID = session.InvoiceID
			result.PaymentIntentID = session.PaymentIntentID
			break
		}
	}
	if result.InvoiceID != "" {
		for _, invoice := range snapshot.Invoices {
			if invoice.ID == result.InvoiceID {
				result.PaymentIntentID = invoice.PaymentIntentID
				result.CustomerID = firstFixtureNonEmpty(result.CustomerID, invoice.CustomerID)
				break
			}
		}
	}
	for _, customer := range snapshot.Customers {
		if customer.ID == ref || customer.Metadata[MetadataFixtureRef] == ref {
			result.CustomerID = customer.ID
			result.Metadata = customer.Metadata
			break
		}
	}
	for _, price := range snapshot.Prices {
		if price.ID == ref || price.LookupKey == ref || price.Metadata[MetadataFixtureRef] == ref {
			result.PriceID = price.ID
			result.ProductID = price.ProductID
			result.Metadata = price.Metadata
			break
		}
	}
	for _, product := range snapshot.Products {
		if product.ID == ref || product.Metadata[MetadataFixtureRef] == ref {
			result.ProductID = product.ID
			result.Metadata = product.Metadata
			break
		}
	}
	if result.CustomerID == "" && result.SubscriptionID == "" && result.InvoiceID == "" && result.PaymentIntentID == "" && result.CheckoutSessionID == "" && result.PriceID == "" && result.ProductID == "" {
		return ResolveResult{}, billing.ErrNotFound
	}
	return result, nil
}

func (s *Service) upsertCustomer(ctx context.Context, pack Pack, fixture CustomerFixture) (billing.Customer, error) {
	metadata := fixtureMetadata(fixture.Metadata, pack, firstFixtureNonEmpty(fixture.Ref, fixture.ID))
	if strings.TrimSpace(fixture.TestClock) != "" {
		if metadata == nil {
			metadata = map[string]string{}
		}
		metadata["test_clock"] = strings.TrimSpace(fixture.TestClock)
	}
	var err error
	metadata, err = applyCustomerPaymentMethodFixture(metadata, fixture)
	if err != nil {
		return billing.Customer{}, err
	}
	metadata = applyCustomerDefaultPaymentIntentOutcome(metadata, fixture)
	if fixture.ID != "" {
		current, err := s.billing.GetCustomer(ctx, fixture.ID)
		if err == nil {
			metadata = mergeStringMap(current.Metadata, metadata)
			metadata, err = applyCustomerPaymentMethodFixture(metadata, fixture)
			if err != nil {
				return billing.Customer{}, err
			}
			metadata = applyCustomerDefaultPaymentIntentOutcome(metadata, fixture)
			return s.billing.UpdateCustomer(ctx, fixture.ID, billing.Customer{
				Email:    fixture.Email,
				Name:     fixture.Name,
				Metadata: metadata,
			})
		}
		if !errors.Is(err, billing.ErrNotFound) {
			return billing.Customer{}, err
		}
	}
	return s.billing.CreateCustomer(ctx, billing.Customer{
		ID:       fixture.ID,
		Email:    fixture.Email,
		Name:     fixture.Name,
		Metadata: metadata,
	})
}

func (s *Service) upsertProduct(ctx context.Context, pack Pack, fixture ProductFixture) (billing.Product, error) {
	active := true
	if fixture.Active != nil {
		active = *fixture.Active
	}
	metadata := fixtureMetadata(fixture.Metadata, pack, fixture.ID)
	if fixture.ID != "" {
		current, err := s.billing.GetProduct(ctx, fixture.ID)
		if err == nil {
			metadata = mergeStringMap(current.Metadata, metadata)
			return s.billing.UpdateProduct(ctx, fixture.ID, billing.Product{
				Name:        fixture.Name,
				Description: fixture.Description,
				Active:      active,
				Metadata:    metadata,
			})
		}
		if !errors.Is(err, billing.ErrNotFound) {
			return billing.Product{}, err
		}
	}
	return s.billing.CreateProduct(ctx, billing.Product{
		ID:          fixture.ID,
		Name:        fixture.Name,
		Description: fixture.Description,
		Active:      active,
		Metadata:    metadata,
	})
}

func (s *Service) upsertPrice(ctx context.Context, pack Pack, fixture PriceFixture) (billing.Price, error) {
	active := true
	if fixture.Active != nil {
		active = *fixture.Active
	}
	metadata := fixtureMetadata(fixture.Metadata, pack, fixture.ID)
	if fixture.ID != "" {
		current, err := s.billing.GetPrice(ctx, fixture.ID)
		if err == nil {
			metadata = mergeStringMap(current.Metadata, metadata)
			return s.billing.UpdatePrice(ctx, fixture.ID, billing.Price{
				ProductID:              fixture.Product,
				Currency:               fixture.Currency,
				UnitAmount:             fixtureUnitAmount(fixture),
				LookupKey:              fixtureLookupKey(fixture),
				RecurringInterval:      fixture.Interval,
				RecurringIntervalCount: fixtureIntervalCount(fixture),
				Active:                 active,
				Metadata:               metadata,
			})
		}
		if !errors.Is(err, billing.ErrNotFound) {
			return billing.Price{}, err
		}
	}
	return s.billing.CreatePrice(ctx, billing.Price{
		ID:                     fixture.ID,
		ProductID:              fixture.Product,
		Currency:               fixture.Currency,
		UnitAmount:             fixtureUnitAmount(fixture),
		LookupKey:              fixtureLookupKey(fixture),
		RecurringInterval:      fixture.Interval,
		RecurringIntervalCount: fixtureIntervalCount(fixture),
		Active:                 active,
		Metadata:               metadata,
	})
}

func (s *Service) upsertTestClock(ctx context.Context, fixture TestClockFixture) (billing.TestClock, error) {
	frozenTime, err := parseFixtureTime(firstFixtureNonEmpty(fixture.FrozenTime, time.Now().UTC().Format(time.RFC3339Nano)))
	if err != nil {
		return billing.TestClock{}, err
	}
	if strings.TrimSpace(fixture.ID) != "" {
		current, err := s.billing.GetTestClock(ctx, fixture.ID)
		if err == nil {
			changed := false
			if !frozenTime.Equal(current.FrozenTime) {
				current.FrozenTime = frozenTime
				changed = true
			}
			if strings.TrimSpace(fixture.Name) != "" {
				current.Name = strings.TrimSpace(fixture.Name)
				changed = true
			}
			if changed {
				return s.billing.UpdateTestClock(ctx, current)
			}
			return current, nil
		}
		if !errors.Is(err, billing.ErrNotFound) {
			return billing.TestClock{}, err
		}
	}
	return s.billing.CreateTestClock(ctx, billing.TestClock{
		ID:         strings.TrimSpace(fixture.ID),
		Name:       strings.TrimSpace(fixture.Name),
		FrozenTime: frozenTime,
	})
}

func (s *Service) upsertSubscription(ctx context.Context, pack Pack, fixture SubscriptionFixture) (billing.CheckoutSession, billing.Subscription, error) {
	ref := subscriptionRef(fixture)
	items := subscriptionLineItems(fixture)
	existing, found, err := s.findSubscription(ctx, pack, fixture, ref)
	if err != nil {
		return billing.CheckoutSession{}, billing.Subscription{}, err
	}
	metadata := fixtureMetadata(fixture.Metadata, pack, ref)
	if strings.TrimSpace(fixture.TestClock) != "" {
		metadata = ensureStringMap(metadata)
		metadata["test_clock"] = strings.TrimSpace(fixture.TestClock)
	}
	if strings.TrimSpace(fixture.RenewalOutcome) != "" {
		metadata = ensureStringMap(metadata)
		metadata["billtap_renewal_outcome"] = strings.TrimSpace(fixture.RenewalOutcome)
	}
	if strings.TrimSpace(fixture.TrialStart) != "" {
		metadata = ensureStringMap(metadata)
		metadata["trial_start"] = strings.TrimSpace(fixture.TrialStart)
	}
	if strings.TrimSpace(fixture.TrialEnd) != "" {
		metadata = ensureStringMap(metadata)
		metadata["trial_end"] = strings.TrimSpace(fixture.TrialEnd)
	}
	if found {
		patch, err := subscriptionStatePatch(fixture, metadata, items)
		if err != nil {
			return billing.CheckoutSession{}, billing.Subscription{}, err
		}
		patch.Items = items
		patch.ReplaceItems = true
		if patch.CancelAtPeriodEnd == nil {
			patch.CancelAtPeriodEnd = fixtureCancelAtPeriodEnd(fixture)
		}
		subscription, err := s.billing.PatchSubscription(ctx, existing.ID, patch)
		return billing.CheckoutSession{}, subscription, err
	}

	trialDays, completionAt, err := fixtureTrialCompletionTiming(fixture)
	if err != nil {
		return billing.CheckoutSession{}, billing.Subscription{}, err
	}
	session, err := s.billing.CreateCheckoutSession(ctx, billing.CheckoutSession{
		ID:              strings.TrimSpace(fixtureCheckoutSessionID(fixture)),
		CustomerID:      fixture.Customer,
		Mode:            "subscription",
		LineItems:       items,
		TrialPeriodDays: trialDays,
	})
	if err != nil {
		return billing.CheckoutSession{}, billing.Subscription{}, err
	}
	completed, err := s.billing.CompleteCheckoutWithOptions(ctx, session.ID, fixtureOutcome(fixture), billing.CheckoutCompletionOptions{
		SubscriptionID:  strings.TrimSpace(fixture.ID),
		InvoiceID:       strings.TrimSpace(fixture.InvoiceID),
		PaymentIntentID: strings.TrimSpace(fixturePaymentIntentID(fixture)),
		At:              completionAt,
	})
	if err != nil {
		return billing.CheckoutSession{}, billing.Subscription{}, err
	}
	subscription, err := s.billing.GetSubscription(ctx, completed.SubscriptionID)
	if err != nil {
		return completed, billing.Subscription{}, err
	}
	patch, err := subscriptionStatePatch(fixture, metadata, nil)
	if err != nil {
		return billing.CheckoutSession{}, billing.Subscription{}, err
	}
	if patch.CancelAtPeriodEnd == nil {
		patch.CancelAtPeriodEnd = fixtureCancelAtPeriodEnd(fixture)
	}
	subscription, err = s.billing.PatchSubscription(ctx, subscription.ID, patch)
	return completed, subscription, err
}

func (s *Service) findSubscription(ctx context.Context, pack Pack, fixture SubscriptionFixture, ref string) (billing.Subscription, bool, error) {
	if id := strings.TrimSpace(fixture.ID); id != "" {
		subscription, err := s.billing.GetSubscription(ctx, id)
		if err == nil {
			return subscription, true, nil
		}
		if !errors.Is(err, billing.ErrNotFound) {
			return billing.Subscription{}, false, err
		}
	}
	subscriptions, err := s.billing.ListSubscriptions(ctx)
	if err != nil {
		return billing.Subscription{}, false, err
	}
	filter := pack.defaultFilter()
	for _, subscription := range subscriptions {
		if subscription.Metadata[MetadataFixtureRef] != ref {
			continue
		}
		if fixtureMetadataMatches(subscription.Metadata, filter) {
			return subscription, true, nil
		}
	}
	return billing.Subscription{}, false, nil
}

func (s *Service) createRefund(ctx context.Context, pack Pack, fixture RefundFixture) (billing.Refund, error) {
	if strings.TrimSpace(fixture.ID) != "" {
		current, err := s.billing.GetRefund(ctx, fixture.ID)
		if err == nil {
			return current, nil
		}
		if !errors.Is(err, billing.ErrNotFound) {
			return billing.Refund{}, err
		}
	}
	metadata := fixtureMetadata(fixture.Metadata, pack, firstFixtureNonEmpty(fixture.ID, fixture.Charge, fixture.PaymentIntent, fixture.Invoice))
	return s.billing.CreateRefund(ctx, billing.Refund{
		ID:              strings.TrimSpace(fixture.ID),
		ChargeID:        strings.TrimSpace(fixture.Charge),
		PaymentIntentID: strings.TrimSpace(fixture.PaymentIntent),
		InvoiceID:       strings.TrimSpace(fixture.Invoice),
		CustomerID:      strings.TrimSpace(fixture.Customer),
		Amount:          fixture.Amount,
		Currency:        strings.TrimSpace(fixture.Currency),
		Reason:          strings.TrimSpace(fixture.Reason),
		Metadata:        metadata,
	})
}

func (s *Service) createCreditNote(ctx context.Context, pack Pack, fixture CreditNoteFixture) (billing.CreditNote, error) {
	if strings.TrimSpace(fixture.ID) != "" {
		current, err := s.billing.GetCreditNote(ctx, fixture.ID)
		if err == nil {
			return current, nil
		}
		if !errors.Is(err, billing.ErrNotFound) {
			return billing.CreditNote{}, err
		}
	}
	metadata := fixtureMetadata(fixture.Metadata, pack, firstFixtureNonEmpty(fixture.ID, fixture.Invoice))
	return s.billing.CreateCreditNote(ctx, billing.CreditNote{
		ID:         strings.TrimSpace(fixture.ID),
		InvoiceID:  strings.TrimSpace(fixture.Invoice),
		CustomerID: strings.TrimSpace(fixture.Customer),
		Amount:     fixture.Amount,
		Currency:   strings.TrimSpace(fixture.Currency),
		Reason:     strings.TrimSpace(fixture.Reason),
		Metadata:   metadata,
	})
}

func (s *Service) filteredTimeline(ctx context.Context, filter SnapshotFilter, customerIDs, checkoutIDs, subscriptionIDs, invoiceIDs, paymentIntentIDs map[string]bool) ([]billing.TimelineEntry, error) {
	if filter.isZero() {
		return s.billing.Timeline(ctx, billing.TimelineFilter{})
	}
	seen := map[string]bool{}
	var out []billing.TimelineEntry
	appendEntries := func(entries []billing.TimelineEntry) {
		for _, entry := range entries {
			if seen[entry.ID] {
				continue
			}
			seen[entry.ID] = true
			out = append(out, entry)
		}
	}
	for id := range customerIDs {
		entries, err := s.billing.Timeline(ctx, billing.TimelineFilter{CustomerID: id})
		if err != nil {
			return nil, err
		}
		appendEntries(entries)
	}
	for id := range checkoutIDs {
		entries, err := s.billing.Timeline(ctx, billing.TimelineFilter{CheckoutSessionID: id})
		if err != nil {
			return nil, err
		}
		appendEntries(entries)
	}
	for id := range subscriptionIDs {
		entries, err := s.billing.Timeline(ctx, billing.TimelineFilter{SubscriptionID: id})
		if err != nil {
			return nil, err
		}
		appendEntries(entries)
	}
	for id := range invoiceIDs {
		entries, err := s.billing.Timeline(ctx, billing.TimelineFilter{InvoiceID: id})
		if err != nil {
			return nil, err
		}
		appendEntries(entries)
	}
	for id := range paymentIntentIDs {
		entries, err := s.billing.Timeline(ctx, billing.TimelineFilter{PaymentIntentID: id})
		if err != nil {
			return nil, err
		}
		appendEntries(entries)
	}
	return out, nil
}

func normalizePack(pack Pack) Pack {
	if strings.TrimSpace(pack.Name) == "" {
		pack.Name = "fixture"
	}
	pack.Name = strings.TrimSpace(pack.Name)
	pack.RunID = strings.TrimSpace(pack.RunID)
	pack.Namespace = strings.TrimSpace(pack.Namespace)
	pack.Products = append([]ProductFixture{}, append(pack.Catalog.Products, pack.Products...)...)
	pack.Prices = append([]PriceFixture{}, append(pack.Catalog.Prices, pack.Prices...)...)
	for i := range pack.Customers {
		pack.Customers[i].ID = strings.TrimSpace(pack.Customers[i].ID)
	}
	return pack
}

func validatePack(pack Pack) error {
	var problems []string
	for idx, customer := range pack.Customers {
		if strings.TrimSpace(customer.ID) == "" && strings.TrimSpace(customer.Email) == "" {
			problems = append(problems, fmt.Sprintf("customers[%d] requires id or email", idx))
		}
		if _, _, _, configured, err := customerPaymentMethodFixtureConfig(customer); configured && err != nil {
			problems = append(problems, fmt.Sprintf("customers[%d].payment_methods is invalid: %v", idx, err))
		}
		if outcome := customerDefaultPaymentIntentOutcomeFixture(customer); outcome != "" && !billing.IsSupportedPaymentIntentOutcome(outcome) {
			problems = append(problems, fmt.Sprintf("customers[%d].default_payment_intent_outcome is invalid", idx))
		}
	}
	for idx, product := range pack.Products {
		if strings.TrimSpace(product.Name) == "" {
			problems = append(problems, fmt.Sprintf("products[%d].name is required", idx))
		}
	}
	for idx, price := range pack.Prices {
		if strings.TrimSpace(price.Product) == "" {
			problems = append(problems, fmt.Sprintf("prices[%d].product is required", idx))
		}
		if strings.TrimSpace(price.Currency) == "" {
			problems = append(problems, fmt.Sprintf("prices[%d].currency is required", idx))
		}
		if fixtureUnitAmount(price) < 0 {
			problems = append(problems, fmt.Sprintf("prices[%d].unitAmount must be non-negative", idx))
		}
	}
	for idx, clock := range pack.TestClocks {
		if strings.TrimSpace(clock.ID) == "" {
			problems = append(problems, fmt.Sprintf("test_clocks[%d].id is required", idx))
		}
		if strings.TrimSpace(clock.FrozenTime) == "" {
			problems = append(problems, fmt.Sprintf("test_clocks[%d].frozen_time is required", idx))
		} else if _, err := parseFixtureTime(clock.FrozenTime); err != nil {
			problems = append(problems, fmt.Sprintf("test_clocks[%d].frozen_time is invalid", idx))
		}
	}
	for idx, subscription := range pack.Subscriptions {
		if strings.TrimSpace(subscription.Customer) == "" {
			problems = append(problems, fmt.Sprintf("subscriptions[%d].customer is required", idx))
		}
		if strings.TrimSpace(subscription.Price) == "" && len(subscription.Items) == 0 {
			problems = append(problems, fmt.Sprintf("subscriptions[%d].price or items is required", idx))
		}
		for itemIdx, item := range subscription.Items {
			if strings.TrimSpace(item.Price) == "" {
				problems = append(problems, fmt.Sprintf("subscriptions[%d].items[%d].price is required", idx, itemIdx))
			}
		}
	}
	for idx, refund := range pack.Refunds {
		if strings.TrimSpace(refund.Charge) == "" && strings.TrimSpace(refund.PaymentIntent) == "" && strings.TrimSpace(refund.Invoice) == "" {
			problems = append(problems, fmt.Sprintf("refunds[%d].charge, payment_intent, or invoice is required", idx))
		}
		if refund.Amount <= 0 {
			problems = append(problems, fmt.Sprintf("refunds[%d].amount must be positive", idx))
		}
	}
	for idx, note := range pack.CreditNotes {
		if strings.TrimSpace(note.Invoice) == "" {
			problems = append(problems, fmt.Sprintf("credit_notes[%d].invoice is required", idx))
		}
		if note.Amount <= 0 {
			problems = append(problems, fmt.Sprintf("credit_notes[%d].amount must be positive", idx))
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidFixture, strings.Join(problems, "; "))
	}
	return nil
}

func (pack Pack) defaultFilter() SnapshotFilter {
	return SnapshotFilter{
		RunID:       pack.RunID,
		FixtureName: pack.Name,
		Namespace:   pack.Namespace,
	}
}

func normalizeFilter(filter SnapshotFilter) SnapshotFilter {
	filter.CustomerID = strings.TrimSpace(filter.CustomerID)
	filter.RunID = strings.TrimSpace(filter.RunID)
	filter.TenantID = strings.TrimSpace(filter.TenantID)
	filter.FixtureName = strings.TrimSpace(filter.FixtureName)
	filter.Namespace = strings.TrimSpace(filter.Namespace)
	return filter
}

func (filter SnapshotFilter) isZero() bool {
	return filter.CustomerID == "" && filter.RunID == "" && filter.TenantID == "" && filter.FixtureName == "" && filter.Namespace == ""
}

func (filter SnapshotFilter) hasFixtureCriteria() bool {
	return filter.RunID != "" || filter.TenantID != "" || filter.FixtureName != "" || filter.Namespace != ""
}

func filterCustomers(items []billing.Customer, filter SnapshotFilter) []billing.Customer {
	if filter.isZero() {
		return items
	}
	var out []billing.Customer
	for _, item := range items {
		if filter.CustomerID != "" && item.ID != filter.CustomerID {
			continue
		}
		if filter.hasFixtureCriteria() && !fixtureMetadataMatches(item.Metadata, filter) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func filterProducts(items []billing.Product, filter SnapshotFilter) []billing.Product {
	if filter.isZero() || !filter.hasFixtureCriteria() {
		return items
	}
	var out []billing.Product
	for _, item := range items {
		if fixtureMetadataMatches(item.Metadata, filter) {
			out = append(out, item)
		}
	}
	return out
}

func filterPrices(items []billing.Price, filter SnapshotFilter, productIDs map[string]bool) []billing.Price {
	if filter.isZero() || !filter.hasFixtureCriteria() {
		return items
	}
	var out []billing.Price
	for _, item := range items {
		if productIDs[item.ProductID] || fixtureMetadataMatches(item.Metadata, filter) {
			out = append(out, item)
		}
	}
	return out
}

func filterCheckoutSessions(items []billing.CheckoutSession, filter SnapshotFilter, customerIDs, subscriptionIDs map[string]bool) []billing.CheckoutSession {
	if filter.isZero() {
		return items
	}
	var out []billing.CheckoutSession
	for _, item := range items {
		if customerIDs[item.CustomerID] || subscriptionIDs[item.SubscriptionID] {
			out = append(out, item)
		}
	}
	return out
}

func filterSubscriptions(items []billing.Subscription, filter SnapshotFilter, customerIDs, priceIDs map[string]bool) []billing.Subscription {
	if filter.isZero() {
		return items
	}
	var out []billing.Subscription
	for _, item := range items {
		if customerIDs[item.CustomerID] || fixtureMetadataMatches(item.Metadata, filter) || lineItemHasPrice(item.Items, priceIDs) {
			out = append(out, item)
		}
	}
	return out
}

func filterInvoices(items []billing.Invoice, filter SnapshotFilter, customerIDs, subscriptionIDs map[string]bool) []billing.Invoice {
	if filter.isZero() {
		return items
	}
	var out []billing.Invoice
	for _, item := range items {
		if customerIDs[item.CustomerID] || subscriptionIDs[item.SubscriptionID] {
			out = append(out, item)
		}
	}
	return out
}

func filterPaymentIntents(items []billing.PaymentIntent, filter SnapshotFilter, customerIDs, invoiceIDs map[string]bool) []billing.PaymentIntent {
	if filter.isZero() {
		return items
	}
	var out []billing.PaymentIntent
	for _, item := range items {
		if customerIDs[item.CustomerID] || invoiceIDs[item.InvoiceID] {
			out = append(out, item)
		}
	}
	return out
}

func fixtureMetadata(base map[string]string, pack Pack, ref string) map[string]string {
	out := mergeStringMap(nil, base)
	if out == nil {
		out = map[string]string{}
	}
	if pack.Name != "" {
		out[MetadataFixtureName] = pack.Name
	}
	if pack.RunID != "" {
		out[MetadataFixtureRunID] = pack.RunID
	}
	if pack.Namespace != "" {
		out[MetadataFixtureNamespace] = pack.Namespace
	}
	if ref != "" {
		out[MetadataFixtureRef] = ref
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func applyCustomerPaymentMethodFixture(metadata map[string]string, fixture CustomerFixture) (map[string]string, error) {
	mode, ids, defaultID, configured, err := customerPaymentMethodFixtureConfig(fixture)
	if err != nil || !configured {
		return metadata, err
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata[billing.MetadataPaymentMethodsFixture] = mode
	delete(metadata, billing.MetadataPaymentMethodIDs)
	delete(metadata, billing.MetadataDefaultPaymentMethod)
	if len(ids) > 0 {
		metadata[billing.MetadataPaymentMethodIDs] = strings.Join(ids, ",")
	}
	if defaultID != "" {
		metadata[billing.MetadataDefaultPaymentMethod] = defaultID
	}
	return metadata, nil
}

func applyCustomerDefaultPaymentIntentOutcome(metadata map[string]string, fixture CustomerFixture) map[string]string {
	outcome := customerDefaultPaymentIntentOutcomeFixture(fixture)
	if outcome == "" {
		return metadata
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata[billing.MetadataDefaultPaymentIntentOutcome] = outcome
	return metadata
}

func customerDefaultPaymentIntentOutcomeFixture(fixture CustomerFixture) string {
	return firstFixtureNonEmpty(
		fixture.DefaultPaymentIntentOutcome,
		fixture.DefaultPIOutcomeCamel,
		fixture.Metadata[billing.MetadataDefaultPaymentIntentOutcome],
		fixture.Metadata["default_payment_intent_outcome"],
	)
}

func customerPaymentMethodFixtureConfig(fixture CustomerFixture) (string, []string, string, bool, error) {
	mode := strings.ToLower(firstFixtureNonEmpty(fixture.PaymentMethodsFixture, fixture.PaymentMethodsFixtureCamel))
	if mode != "" && mode != billing.PaymentMethodsFixtureEmpty && mode != billing.PaymentMethodsFixtureExplicit {
		return "", nil, "", true, fmt.Errorf("unsupported payment_methods_fixture %q", mode)
	}
	methodsConfigured := fixture.PaymentMethods != nil || fixture.PaymentMethodsCamel != nil
	methods := append([]PaymentMethodFixture{}, fixture.PaymentMethods...)
	methods = append(methods, fixture.PaymentMethodsCamel...)
	if !methodsConfigured {
		if mode == "" {
			return "", nil, "", false, nil
		}
		return mode, nil, "", true, nil
	}
	if len(methods) == 0 {
		return billing.PaymentMethodsFixtureEmpty, nil, "", true, nil
	}

	ids := make([]string, 0, len(methods))
	defaultID := ""
	for idx, method := range methods {
		if strings.TrimSpace(method.ID) == "" {
			return "", nil, "", true, fmt.Errorf("payment_methods[%d].id is required", idx)
		}
		id := strings.TrimSpace(method.ID)
		ids = append(ids, id)
		if method.Default {
			defaultID = id
		}
	}
	return billing.PaymentMethodsFixtureExplicit, uniqueFixtureStrings(ids), defaultID, true, nil
}

func fixtureMetadataMatches(metadata map[string]string, filter SnapshotFilter) bool {
	if filter.RunID != "" && metadata[MetadataFixtureRunID] != filter.RunID {
		return false
	}
	if filter.TenantID != "" && metadata["tenantId"] != filter.TenantID {
		return false
	}
	if filter.FixtureName != "" && metadata[MetadataFixtureName] != filter.FixtureName {
		return false
	}
	if filter.Namespace != "" && metadata[MetadataFixtureNamespace] != filter.Namespace {
		return false
	}
	return true
}

func mergeStringMap(base map[string]string, overlay map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range base {
		out[key] = value
	}
	for key, value := range overlay {
		out[key] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func positiveQuantity(value int64) int64 {
	if value <= 0 {
		return 1
	}
	return value
}

func subscriptionRef(fixture SubscriptionFixture) string {
	if strings.TrimSpace(fixture.Ref) != "" {
		return strings.TrimSpace(fixture.Ref)
	}
	if strings.TrimSpace(fixture.ID) != "" {
		return strings.TrimSpace(fixture.ID)
	}
	return strings.TrimSpace(fixture.Customer + "_" + fixture.Price)
}

func subscriptionLineItems(fixture SubscriptionFixture) []billing.LineItem {
	if len(fixture.Items) == 0 {
		return []billing.LineItem{{PriceID: strings.TrimSpace(fixture.Price), Quantity: positiveQuantity(fixture.Quantity)}}
	}
	out := make([]billing.LineItem, 0, len(fixture.Items))
	for _, item := range fixture.Items {
		out = append(out, billing.LineItem{
			PriceID:  strings.TrimSpace(item.Price),
			Quantity: positiveQuantity(item.Quantity),
		})
	}
	return out
}

func subscriptionStatePatch(fixture SubscriptionFixture, metadata map[string]string, items []billing.LineItem) (billing.SubscriptionPatch, error) {
	patch := billing.SubscriptionPatch{
		Metadata:        metadata,
		TimelineSource:  "fixture",
		TimelineMessage: "Fixture subscription state applied",
	}
	if len(items) > 0 {
		patch.Items = items
		patch.ReplaceItems = true
	}
	status := strings.ToLower(strings.TrimSpace(fixture.Status))
	if status != "" {
		switch status {
		case "active", "trialing", "past_due", "unpaid", "incomplete", "incomplete_expired", "canceled":
			patch.Status = &status
		default:
			return patch, fmt.Errorf("%w: unsupported subscription status %q", ErrInvalidFixture, fixture.Status)
		}
	}
	if value := firstFixtureNonEmpty(fixture.CurrentPeriodStart, fixture.TrialStart); value != "" {
		parsed, err := parseFixtureTime(value)
		if err != nil {
			return patch, err
		}
		patch.CurrentPeriodStart = &parsed
	}
	if value := firstFixtureNonEmpty(fixture.CurrentPeriodEnd, fixture.TrialEnd, fixture.CancelAt); value != "" {
		parsed, err := parseFixtureTime(value)
		if err != nil {
			return patch, err
		}
		patch.CurrentPeriodEnd = &parsed
	}
	if value := strings.TrimSpace(fixture.TrialStart); value != "" {
		if _, err := parseFixtureTime(value); err != nil {
			return patch, err
		}
		metadata = ensureStringMap(metadata)
		metadata["trial_start"] = value
		patch.Metadata = metadata
	}
	if value := strings.TrimSpace(fixture.TrialEnd); value != "" {
		if _, err := parseFixtureTime(value); err != nil {
			return patch, err
		}
		metadata = ensureStringMap(metadata)
		metadata["trial_end"] = value
		patch.Metadata = metadata
	}
	if value := strings.TrimSpace(fixture.CancelAt); value != "" {
		if _, err := parseFixtureTime(value); err != nil {
			return patch, err
		}
		metadata = ensureStringMap(metadata)
		metadata["cancel_at"] = value
		patch.Metadata = metadata
		trueValue := true
		patch.CancelAtPeriodEnd = &trueValue
	}
	if value := strings.TrimSpace(fixture.CanceledAt); value != "" {
		parsed, err := parseFixtureTime(value)
		if err != nil {
			return patch, err
		}
		patch.CanceledAt = &parsed
	}
	if value := strings.TrimSpace(fixture.EndedAt); value != "" {
		if _, err := parseFixtureTime(value); err != nil {
			return patch, err
		}
		metadata = ensureStringMap(metadata)
		metadata["ended_at"] = value
		patch.Metadata = metadata
		if patch.CanceledAt == nil {
			parsed, _ := parseFixtureTime(value)
			patch.CanceledAt = &parsed
		}
	}
	if status == "canceled" && patch.CanceledAt == nil {
		canceledAt := time.Now().UTC()
		patch.CanceledAt = &canceledAt
		falseValue := false
		patch.CancelAtPeriodEnd = &falseValue
	}
	if strings.TrimSpace(fixture.LatestInvoiceStatus) != "" {
		metadata = ensureStringMap(metadata)
		metadata["latest_invoice_status"] = strings.TrimSpace(fixture.LatestInvoiceStatus)
		patch.Metadata = metadata
	}
	return patch, nil
}

func fixtureTrialCompletionTiming(fixture SubscriptionFixture) (int64, time.Time, error) {
	status := strings.ToLower(strings.TrimSpace(fixture.Status))
	if status != "trialing" && strings.TrimSpace(fixture.TrialEnd) == "" {
		return 0, time.Time{}, nil
	}
	trialStartRaw := firstFixtureNonEmpty(fixture.TrialStart, fixture.CurrentPeriodStart)
	trialEndRaw := firstFixtureNonEmpty(fixture.TrialEnd, fixture.CurrentPeriodEnd)
	if trialStartRaw == "" || trialEndRaw == "" {
		return 0, time.Time{}, nil
	}
	trialStart, err := parseFixtureTime(trialStartRaw)
	if err != nil {
		return 0, time.Time{}, err
	}
	trialEnd, err := parseFixtureTime(trialEndRaw)
	if err != nil {
		return 0, time.Time{}, err
	}
	if !trialEnd.After(trialStart) {
		return 0, time.Time{}, fmt.Errorf("%w: trial_end must be after trial_start", ErrInvalidFixture)
	}
	days := int64(trialEnd.Sub(trialStart).Hours() / 24)
	if days <= 0 {
		days = 1
	}
	return days, trialStart, nil
}

func fixtureOutcome(fixture SubscriptionFixture) string {
	status := strings.ToLower(strings.TrimSpace(fixture.Status))
	switch status {
	case "trialing":
		return "payment_succeeded"
	case "canceled", "incomplete_expired":
		return "canceled"
	case "past_due", "unpaid":
		return "card_declined"
	case "incomplete":
		return "payment_pending"
	}
	if strings.TrimSpace(fixture.Outcome) != "" {
		return strings.TrimSpace(fixture.Outcome)
	}
	return "payment_succeeded"
}

func fixtureCancelAtPeriodEnd(fixture SubscriptionFixture) *bool {
	if fixture.CancelAtPeriodEnd != nil {
		return fixture.CancelAtPeriodEnd
	}
	return fixture.CancelAtPeriodEndSnake
}

func fixtureCheckoutSessionID(fixture SubscriptionFixture) string {
	return firstFixtureNonEmpty(fixture.CheckoutSessionID, fixture.CheckoutSessionIDSnake)
}

func fixturePaymentIntentID(fixture SubscriptionFixture) string {
	return firstFixtureNonEmpty(fixture.PaymentIntentID, fixture.PaymentIntentIDSnake)
}

func fixtureUnitAmount(fixture PriceFixture) int64 {
	if fixture.UnitAmount != 0 {
		return fixture.UnitAmount
	}
	return fixture.UnitAmountSnake
}

func fixtureLookupKey(fixture PriceFixture) string {
	return firstFixtureNonEmpty(fixture.LookupKey, fixture.LookupKeySnake)
}

func fixtureIntervalCount(fixture PriceFixture) int {
	if fixture.IntervalCount > 0 {
		return fixture.IntervalCount
	}
	return fixture.IntervalCountSnake
}

func parseFixtureTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}
	if unix, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return unix.UTC(), nil
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(seconds, 0).UTC(), nil
	}
	return time.Time{}, fmt.Errorf("%w: invalid fixture timestamp %q", ErrInvalidFixture, value)
}

func firstFixtureNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func uniqueFixtureStrings(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func ensureStringMap(in map[string]string) map[string]string {
	if in != nil {
		return in
	}
	return map[string]string{}
}

func lineItemHasPrice(items []billing.LineItem, priceIDs map[string]bool) bool {
	for _, item := range items {
		if priceIDs[item.PriceID] {
			return true
		}
	}
	return false
}

func customerSet(items []billing.Customer) map[string]bool {
	out := map[string]bool{}
	for _, item := range items {
		out[item.ID] = true
	}
	return out
}

func productSet(items []billing.Product) map[string]bool {
	out := map[string]bool{}
	for _, item := range items {
		out[item.ID] = true
	}
	return out
}

func priceSet(items []billing.Price) map[string]bool {
	out := map[string]bool{}
	for _, item := range items {
		out[item.ID] = true
	}
	return out
}

func checkoutSessionSet(items []billing.CheckoutSession) map[string]bool {
	out := map[string]bool{}
	for _, item := range items {
		out[item.ID] = true
	}
	return out
}

func subscriptionSet(items []billing.Subscription) map[string]bool {
	out := map[string]bool{}
	for _, item := range items {
		out[item.ID] = true
	}
	return out
}

func invoiceSet(items []billing.Invoice) map[string]bool {
	out := map[string]bool{}
	for _, item := range items {
		out[item.ID] = true
	}
	return out
}

func paymentIntentSet(items []billing.PaymentIntent) map[string]bool {
	out := map[string]bool{}
	for _, item := range items {
		out[item.ID] = true
	}
	return out
}

func sanitizeID(value string) string {
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "fixture"
	}
	return b.String()
}
