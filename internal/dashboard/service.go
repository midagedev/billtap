package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/webhooks"
)

var (
	ErrInvalidInput            = errors.New("invalid input")
	ErrUnsupportedObjectType   = errors.New("unsupported object type")
	ErrMissingBillingSource    = errors.New("missing billing source")
	ErrMissingWebhookSource    = errors.New("missing webhook source")
	ErrCheckoutListUnsupported = errors.New("checkout session listing is unsupported by billing source")
)

type BillingSource interface {
	GetCustomer(context.Context, string) (billing.Customer, error)
	ListCustomers(context.Context) ([]billing.Customer, error)
	GetProduct(context.Context, string) (billing.Product, error)
	ListProducts(context.Context) ([]billing.Product, error)
	GetPrice(context.Context, string) (billing.Price, error)
	ListPrices(context.Context) ([]billing.Price, error)
	GetCheckoutSession(context.Context, string) (billing.CheckoutSession, error)
	GetSubscription(context.Context, string) (billing.Subscription, error)
	ListSubscriptions(context.Context) ([]billing.Subscription, error)
	GetInvoice(context.Context, string) (billing.Invoice, error)
	ListInvoices(context.Context) ([]billing.Invoice, error)
	GetPaymentIntent(context.Context, string) (billing.PaymentIntent, error)
	ListPaymentIntents(context.Context) ([]billing.PaymentIntent, error)
	Timeline(context.Context, billing.TimelineFilter) ([]billing.TimelineEntry, error)
}

type CheckoutSessionLister interface {
	ListCheckoutSessions(context.Context) ([]billing.CheckoutSession, error)
}

type WebhookSource interface {
	GetEndpoint(context.Context, string) (webhooks.Endpoint, error)
	ListEndpoints(context.Context, webhooks.EndpointFilter) ([]webhooks.Endpoint, error)
	GetEvent(context.Context, string) (webhooks.Event, error)
	ListEvents(context.Context, webhooks.EventFilter) ([]webhooks.Event, error)
	ListDeliveryAttempts(context.Context, webhooks.DeliveryAttemptFilter) ([]webhooks.DeliveryAttempt, error)
}

type Service struct {
	billing  BillingSource
	webhooks WebhookSource
	now      func() time.Time
}

func NewService(b BillingSource, w WebhookSource) *Service {
	return &Service{
		billing:  b,
		webhooks: w,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) ListObjects(ctx context.Context, q ObjectListQuery) (ObjectListResult, error) {
	typ := normalizeType(q.Type)
	if typ == "" {
		return ObjectListResult{}, fmt.Errorf("%w: object type is required", ErrInvalidInput)
	}
	result := ObjectListResult{Object: ObjectList, Type: typ}

	switch typ {
	case billing.ObjectCustomer:
		if s.billing == nil {
			return result, ErrMissingBillingSource
		}
		items, err := s.billing.ListCustomers(ctx)
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, customerSummary(item))
		}
	case billing.ObjectProduct:
		if s.billing == nil {
			return result, ErrMissingBillingSource
		}
		items, err := s.billing.ListProducts(ctx)
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, productSummary(item))
		}
	case billing.ObjectPrice:
		if s.billing == nil {
			return result, ErrMissingBillingSource
		}
		items, err := s.billing.ListPrices(ctx)
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, priceSummary(item))
		}
	case billing.ObjectCheckoutSession:
		if s.billing == nil {
			return result, ErrMissingBillingSource
		}
		lister, ok := s.billing.(CheckoutSessionLister)
		if !ok {
			return result, ErrCheckoutListUnsupported
		}
		items, err := lister.ListCheckoutSessions(ctx)
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, checkoutSessionSummary(item))
		}
	case billing.ObjectSubscription:
		if s.billing == nil {
			return result, ErrMissingBillingSource
		}
		items, err := s.billing.ListSubscriptions(ctx)
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, subscriptionSummary(item))
		}
	case billing.ObjectInvoice:
		if s.billing == nil {
			return result, ErrMissingBillingSource
		}
		items, err := s.billing.ListInvoices(ctx)
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, invoiceSummary(item))
		}
	case billing.ObjectPaymentIntent:
		if s.billing == nil {
			return result, ErrMissingBillingSource
		}
		items, err := s.billing.ListPaymentIntents(ctx)
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, paymentIntentSummary(item))
		}
	case TypeWebhookEvent, webhooks.ObjectEvent:
		if s.webhooks == nil {
			return result, ErrMissingWebhookSource
		}
		items, err := s.webhooks.ListEvents(ctx, webhooks.EventFilter{})
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, eventSummary(item))
		}
	case TypeWebhookEndpoint, webhooks.ObjectEndpoint:
		if s.webhooks == nil {
			return result, ErrMissingWebhookSource
		}
		items, err := s.webhooks.ListEndpoints(ctx, webhooks.EndpointFilter{})
		if err != nil {
			return result, err
		}
		for _, item := range items {
			result.Data = append(result.Data, endpointSummary(item))
		}
	default:
		return result, fmt.Errorf("%w: %s", ErrUnsupportedObjectType, q.Type)
	}
	return result, nil
}

func (s *Service) ObjectDetail(ctx context.Context, ref ObjectRef) (ObjectDetail, error) {
	ref.Type = normalizeType(ref.Type)
	if ref.Type == "" || strings.TrimSpace(ref.ID) == "" {
		return ObjectDetail{}, fmt.Errorf("%w: object type and id are required", ErrInvalidInput)
	}

	detail := ObjectDetail{Object: ObjectDetailObject, Ref: ref}
	switch ref.Type {
	case billing.ObjectCustomer:
		if s.billing == nil {
			return detail, ErrMissingBillingSource
		}
		data, err := s.billing.GetCustomer(ctx, ref.ID)
		detail.Data = data
		return detail, err
	case billing.ObjectProduct:
		if s.billing == nil {
			return detail, ErrMissingBillingSource
		}
		data, err := s.billing.GetProduct(ctx, ref.ID)
		detail.Data = data
		return detail, err
	case billing.ObjectPrice:
		if s.billing == nil {
			return detail, ErrMissingBillingSource
		}
		data, err := s.billing.GetPrice(ctx, ref.ID)
		detail.Data = data
		detail.Related = map[string]string{billing.ObjectProduct: data.ProductID}
		return detail, err
	case billing.ObjectCheckoutSession:
		if s.billing == nil {
			return detail, ErrMissingBillingSource
		}
		data, err := s.billing.GetCheckoutSession(ctx, ref.ID)
		detail.Data = data
		detail.Related = checkoutSessionRelated(data)
		return detail, err
	case billing.ObjectSubscription:
		if s.billing == nil {
			return detail, ErrMissingBillingSource
		}
		data, err := s.billing.GetSubscription(ctx, ref.ID)
		detail.Data = data
		detail.Related = subscriptionRelated(data)
		return detail, err
	case billing.ObjectInvoice:
		if s.billing == nil {
			return detail, ErrMissingBillingSource
		}
		data, err := s.billing.GetInvoice(ctx, ref.ID)
		detail.Data = data
		detail.Related = invoiceRelated(data)
		return detail, err
	case billing.ObjectPaymentIntent:
		if s.billing == nil {
			return detail, ErrMissingBillingSource
		}
		data, err := s.billing.GetPaymentIntent(ctx, ref.ID)
		detail.Data = data
		detail.Related = paymentIntentRelated(data)
		return detail, err
	case TypeWebhookEvent, webhooks.ObjectEvent:
		if s.webhooks == nil {
			return detail, ErrMissingWebhookSource
		}
		data, err := s.webhooks.GetEvent(ctx, ref.ID)
		detail.Data = data
		return detail, err
	case TypeWebhookEndpoint, webhooks.ObjectEndpoint:
		if s.webhooks == nil {
			return detail, ErrMissingWebhookSource
		}
		data, err := s.webhooks.GetEndpoint(ctx, ref.ID)
		detail.Data = data
		return detail, err
	default:
		return detail, fmt.Errorf("%w: %s", ErrUnsupportedObjectType, ref.Type)
	}
}

func (s *Service) Timeline(ctx context.Context, q TimelineQuery) (TimelineResult, error) {
	result := TimelineResult{Object: ObjectTimeline, Filter: q}
	if s.billing != nil {
		entries, err := s.billing.Timeline(ctx, billingTimelineFilter(q))
		if err != nil {
			return result, err
		}
		for _, entry := range entries {
			result.Data = append(result.Data, billingTimelineItem(entry))
		}
	}

	if q.IncludeWebhooks || q.EventID != "" || q.EventType != "" || q.ScenarioRunID != "" {
		if s.webhooks == nil {
			return result, ErrMissingWebhookSource
		}
		events, err := s.webhooks.ListEvents(ctx, webhooks.EventFilter{Type: q.EventType, ScenarioRunID: q.ScenarioRunID})
		if err != nil {
			return result, err
		}
		for _, event := range events {
			if !eventMatchesTimeline(event, q) {
				continue
			}
			result.Data = append(result.Data, eventTimelineItem(event))
			attempts, err := s.webhooks.ListDeliveryAttempts(ctx, webhooks.DeliveryAttemptFilter{EventID: event.ID})
			if err != nil {
				return result, err
			}
			for _, attempt := range attempts {
				result.Data = append(result.Data, attemptTimelineItem(attempt))
			}
		}
	}

	sort.SliceStable(result.Data, func(i, j int) bool {
		if result.Data[i].At.Equal(result.Data[j].At) {
			return result.Data[i].ID < result.Data[j].ID
		}
		return result.Data[i].At.Before(result.Data[j].At)
	})
	return result, nil
}

func (s *Service) WebhookDetail(ctx context.Context, eventID string) (WebhookDetail, error) {
	if strings.TrimSpace(eventID) == "" {
		return WebhookDetail{}, fmt.Errorf("%w: event id is required", ErrInvalidInput)
	}
	if s.webhooks == nil {
		return WebhookDetail{}, ErrMissingWebhookSource
	}

	event, err := s.webhooks.GetEvent(ctx, eventID)
	if err != nil {
		return WebhookDetail{}, err
	}
	attempts, err := s.webhooks.ListDeliveryAttempts(ctx, webhooks.DeliveryAttemptFilter{EventID: event.ID})
	if err != nil {
		return WebhookDetail{}, err
	}
	return buildWebhookDetail(event, attempts), nil
}

func (s *Service) DebugBundle(ctx context.Context, ref ObjectRef) (DebugBundle, error) {
	detail, err := s.ObjectDetail(ctx, ref)
	if err != nil {
		return DebugBundle{}, err
	}
	query := timelineQueryForDetail(detail)
	query.IncludeWebhooks = true

	timeline, err := s.Timeline(ctx, query)
	if err != nil {
		return DebugBundle{}, err
	}

	bundle := DebugBundle{
		Object:      ObjectDebugBundle,
		Target:      ObjectRef{Type: detail.Ref.Type, ID: detail.Ref.ID},
		GeneratedAt: s.now(),
		Detail:      detail,
		Timeline:    timeline,
	}
	if s.webhooks != nil {
		webhookDetails, err := s.webhookDetailsForTimeline(ctx, query)
		if err != nil {
			return DebugBundle{}, err
		}
		bundle.Webhooks = webhookDetails
	}
	return bundle, nil
}

func (s *Service) DebugBundleJSON(ctx context.Context, ref ObjectRef) ([]byte, error) {
	bundle, err := s.DebugBundle(ctx, ref)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(bundle, "", "  ")
}

func (s *Service) webhookDetailsForTimeline(ctx context.Context, q TimelineQuery) ([]WebhookDetail, error) {
	if q.EventID != "" {
		detail, err := s.WebhookDetail(ctx, q.EventID)
		if err != nil {
			return nil, err
		}
		return []WebhookDetail{detail}, nil
	}

	events, err := s.webhooks.ListEvents(ctx, webhooks.EventFilter{Type: q.EventType, ScenarioRunID: q.ScenarioRunID})
	if err != nil {
		return nil, err
	}
	var out []WebhookDetail
	for _, event := range events {
		if !eventMatchesTimeline(event, q) {
			continue
		}
		attempts, err := s.webhooks.ListDeliveryAttempts(ctx, webhooks.DeliveryAttemptFilter{EventID: event.ID})
		if err != nil {
			return nil, err
		}
		out = append(out, buildWebhookDetail(event, attempts))
	}
	return out, nil
}

func buildWebhookDetail(event webhooks.Event, attempts []webhooks.DeliveryAttempt) WebhookDetail {
	detail := WebhookDetail{
		Object:       ObjectWebhookDetail,
		Event:        event,
		AttemptCount: len(attempts),
		Flags:        WebhookFlags{},
	}
	sort.SliceStable(attempts, func(i, j int) bool {
		if attempts[i].ScheduledAt.Equal(attempts[j].ScheduledAt) {
			if attempts[i].AttemptNumber == attempts[j].AttemptNumber {
				return attempts[i].ID < attempts[j].ID
			}
			return attempts[i].AttemptNumber < attempts[j].AttemptNumber
		}
		return attempts[i].ScheduledAt.Before(attempts[j].ScheduledAt)
	})

	endpointAttempts := map[string]int{}
	for _, attempt := range attempts {
		endpointAttempts[attempt.EndpointID]++
		evidence := deliveryAttemptEvidence(attempt)
		detail.Attempts = append(detail.Attempts, evidence)
		detail.LatestStatus = attempt.Status
		if detail.SignatureHeader == "" {
			detail.SignatureHeader = evidence.SignatureHeader
		}
		if detail.RequestURL == "" {
			detail.RequestURL = attempt.RequestURL
		}
		if len(detail.RequestBody) == 0 {
			detail.RequestBody = attempt.RequestBody
		}
		if attempt.NextRetryAt != nil {
			detail.RetryPlan = append(detail.RetryPlan, RetryEvidence{
				AttemptID:  attempt.ID,
				EndpointID: attempt.EndpointID,
				RetryAt:    *attempt.NextRetryAt,
			})
		}
		if attempt.Metadata["duplicate"] == "true" {
			detail.Flags.Duplicate = true
		}
		if attempt.Metadata["out_of_order"] == "true" {
			detail.Flags.OutOfOrder = true
		}
		if attempt.Metadata["source"] == webhooks.SourceReplay || attempt.Metadata["retry_for"] != "" {
			detail.Flags.Replay = detail.Flags.Replay || attempt.Metadata["source"] == webhooks.SourceReplay
		}
	}
	for _, count := range endpointAttempts {
		if count > 1 {
			detail.Flags.Duplicate = true
		}
	}
	return detail
}

func deliveryAttemptEvidence(attempt webhooks.DeliveryAttempt) DeliveryAttemptEvidence {
	return DeliveryAttemptEvidence{
		ID:              attempt.ID,
		EndpointID:      attempt.EndpointID,
		AttemptNumber:   attempt.AttemptNumber,
		Status:          attempt.Status,
		ScheduledAt:     attempt.ScheduledAt,
		DeliveredAt:     attempt.DeliveredAt,
		RequestURL:      attempt.RequestURL,
		SignatureHeader: attempt.RequestHeaders[webhooks.SignatureHeaderName],
		RequestHeaders:  attempt.RequestHeaders,
		RequestBody:     attempt.RequestBody,
		ResponseStatus:  attempt.ResponseStatus,
		ResponseBody:    attempt.ResponseBody,
		Error:           attempt.Error,
		NextRetryAt:     attempt.NextRetryAt,
		Metadata:        attempt.Metadata,
	}
}

func billingTimelineItem(entry billing.TimelineEntry) TimelineItem {
	return TimelineItem{
		ID:                entry.ID,
		Kind:              "billing",
		Action:            entry.Action,
		Message:           entry.Message,
		ObjectType:        entry.ObjectType,
		ObjectID:          entry.ObjectID,
		CustomerID:        entry.CustomerID,
		CheckoutSessionID: entry.CheckoutSessionID,
		SubscriptionID:    entry.SubscriptionID,
		InvoiceID:         entry.InvoiceID,
		PaymentIntentID:   entry.PaymentIntentID,
		Data:              entry.Data,
		At:                entry.CreatedAt,
	}
}

func eventTimelineItem(event webhooks.Event) TimelineItem {
	related := eventRelatedIDs(event)
	return TimelineItem{
		ID:                event.ID,
		Kind:              "webhook_event",
		Action:            event.Type,
		Message:           "Webhook " + event.Type + " emitted",
		ObjectType:        TypeWebhookEvent,
		ObjectID:          event.ID,
		CustomerID:        related.CustomerID,
		CheckoutSessionID: related.CheckoutSessionID,
		SubscriptionID:    related.SubscriptionID,
		InvoiceID:         related.InvoiceID,
		PaymentIntentID:   related.PaymentIntentID,
		EventID:           event.ID,
		Data: map[string]string{
			"source":           event.Billtap.Source,
			"scenario_run_id":  event.Billtap.ScenarioRunID,
			"pending_webhooks": fmt.Sprintf("%d", event.PendingWebhooks),
		},
		At: event.CreatedAt,
	}
}

func attemptTimelineItem(attempt webhooks.DeliveryAttempt) TimelineItem {
	at := attempt.ScheduledAt
	if attempt.DeliveredAt != nil {
		at = *attempt.DeliveredAt
	}
	return TimelineItem{
		ID:                attempt.ID,
		Kind:              "webhook_attempt",
		Action:            "webhook.delivery." + attempt.Status,
		Message:           fmt.Sprintf("Webhook delivery attempt %d %s", attempt.AttemptNumber, attempt.Status),
		ObjectType:        webhooks.ObjectDeliveryAttempt,
		ObjectID:          attempt.ID,
		EventID:           attempt.EventID,
		DeliveryAttemptID: attempt.ID,
		WebhookStatus:     attempt.Status,
		Data:              attempt.Metadata,
		At:                at,
	}
}

func eventMatchesTimeline(event webhooks.Event, q TimelineQuery) bool {
	if q.EventID != "" && event.ID != q.EventID {
		return false
	}
	if q.EventType != "" && event.Type != q.EventType {
		return false
	}
	if q.ScenarioRunID != "" && event.Billtap.ScenarioRunID != q.ScenarioRunID {
		return false
	}
	if q.CustomerID == "" && q.CheckoutSessionID == "" && q.SubscriptionID == "" && q.InvoiceID == "" && q.PaymentIntentID == "" {
		return true
	}
	related := eventRelatedIDs(event)
	return (q.CustomerID == "" || q.CustomerID == related.CustomerID) &&
		(q.CheckoutSessionID == "" || q.CheckoutSessionID == related.CheckoutSessionID) &&
		(q.SubscriptionID == "" || q.SubscriptionID == related.SubscriptionID) &&
		(q.InvoiceID == "" || q.InvoiceID == related.InvoiceID) &&
		(q.PaymentIntentID == "" || q.PaymentIntentID == related.PaymentIntentID)
}

type relatedIDs struct {
	CustomerID        string
	CheckoutSessionID string
	SubscriptionID    string
	InvoiceID         string
	PaymentIntentID   string
}

func eventRelatedIDs(event webhooks.Event) relatedIDs {
	fields := rawFields(event.Data.Object)
	out := relatedIDs{
		CustomerID:        fields["customer"],
		CheckoutSessionID: fields["checkout_session"],
		SubscriptionID:    fields["subscription"],
		InvoiceID:         fields["invoice"],
		PaymentIntentID:   fields["payment_intent"],
	}
	objectID := fields["id"]
	objectType := fields["object"]
	switch {
	case strings.HasPrefix(event.Type, "checkout.session.") || objectType == billing.ObjectCheckoutSession:
		out.CheckoutSessionID = first(out.CheckoutSessionID, objectID)
	case strings.HasPrefix(event.Type, "customer.subscription.") || objectType == billing.ObjectSubscription:
		out.SubscriptionID = first(out.SubscriptionID, objectID)
	case strings.HasPrefix(event.Type, "invoice.") || objectType == billing.ObjectInvoice:
		out.InvoiceID = first(out.InvoiceID, objectID)
	case strings.HasPrefix(event.Type, "payment_intent.") || objectType == billing.ObjectPaymentIntent:
		out.PaymentIntentID = first(out.PaymentIntentID, objectID)
	case strings.HasPrefix(event.Type, "customer.") || objectType == billing.ObjectCustomer:
		out.CustomerID = first(out.CustomerID, objectID)
	}
	return out
}

func rawFields(raw json.RawMessage) map[string]string {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil
	}
	out := map[string]string{}
	for key, value := range obj {
		if str, ok := value.(string); ok {
			out[key] = str
		}
	}
	return out
}

func timelineQueryForDetail(detail ObjectDetail) TimelineQuery {
	switch data := detail.Data.(type) {
	case billing.Customer:
		return TimelineQuery{CustomerID: data.ID}
	case billing.CheckoutSession:
		return TimelineQuery{CustomerID: data.CustomerID, CheckoutSessionID: data.ID}
	case billing.Subscription:
		return TimelineQuery{CustomerID: data.CustomerID, SubscriptionID: data.ID}
	case billing.Invoice:
		return TimelineQuery{CustomerID: data.CustomerID, SubscriptionID: data.SubscriptionID, InvoiceID: data.ID}
	case billing.PaymentIntent:
		return TimelineQuery{CustomerID: data.CustomerID, InvoiceID: data.InvoiceID, PaymentIntentID: data.ID}
	case webhooks.Event:
		return TimelineQuery{EventID: data.ID, EventType: data.Type, ScenarioRunID: data.Billtap.ScenarioRunID}
	default:
		return TimelineQuery{}
	}
}

func normalizeType(typ string) string {
	typ = strings.TrimSpace(strings.ToLower(typ))
	typ = strings.TrimSuffix(typ, "s")
	switch typ {
	case "customer":
		return billing.ObjectCustomer
	case "product":
		return billing.ObjectProduct
	case "price":
		return billing.ObjectPrice
	case "checkout_session", "checkout.session", "session":
		return billing.ObjectCheckoutSession
	case "subscription":
		return billing.ObjectSubscription
	case "invoice":
		return billing.ObjectInvoice
	case "payment_intent", "payment.intent":
		return billing.ObjectPaymentIntent
	case "event", "webhook_event", "webhook.event":
		return TypeWebhookEvent
	case "webhook_endpoint", "endpoint", "webhook.endpoint":
		return TypeWebhookEndpoint
	default:
		return typ
	}
}

func customerSummary(c billing.Customer) ObjectSummary {
	return ObjectSummary{ID: c.ID, Object: c.Object, Display: first(c.Name, c.Email, c.ID), CreatedAt: c.CreatedAt}
}

func productSummary(p billing.Product) ObjectSummary {
	return ObjectSummary{ID: p.ID, Object: p.Object, Display: first(p.Name, p.ID), Status: activeStatus(p.Active)}
}

func priceSummary(p billing.Price) ObjectSummary {
	return ObjectSummary{ID: p.ID, Object: p.Object, Display: fmt.Sprintf("%s %d", strings.ToUpper(p.Currency), p.UnitAmount), Status: activeStatus(p.Active), Related: map[string]string{billing.ObjectProduct: p.ProductID}}
}

func checkoutSessionSummary(cs billing.CheckoutSession) ObjectSummary {
	return ObjectSummary{ID: cs.ID, Object: cs.Object, Display: cs.ID, Status: cs.Status, CustomerID: cs.CustomerID, Related: checkoutSessionRelated(cs), CreatedAt: cs.CreatedAt}
}

func subscriptionSummary(sub billing.Subscription) ObjectSummary {
	return ObjectSummary{ID: sub.ID, Object: sub.Object, Display: sub.ID, Status: sub.Status, CustomerID: sub.CustomerID, Related: subscriptionRelated(sub), CreatedAt: sub.CurrentPeriodStart}
}

func invoiceSummary(inv billing.Invoice) ObjectSummary {
	return ObjectSummary{ID: inv.ID, Object: inv.Object, Display: inv.ID, Status: inv.Status, CustomerID: inv.CustomerID, Related: invoiceRelated(inv), CreatedAt: inv.CreatedAt}
}

func paymentIntentSummary(pi billing.PaymentIntent) ObjectSummary {
	return ObjectSummary{ID: pi.ID, Object: pi.Object, Display: pi.ID, Status: pi.Status, CustomerID: pi.CustomerID, Related: paymentIntentRelated(pi), CreatedAt: pi.CreatedAt}
}

func eventSummary(event webhooks.Event) ObjectSummary {
	related := eventRelatedIDs(event)
	return ObjectSummary{
		ID:         event.ID,
		Object:     TypeWebhookEvent,
		Display:    event.Type,
		Status:     fmt.Sprintf("pending:%d", event.PendingWebhooks),
		CustomerID: related.CustomerID,
		Related: compactMap(map[string]string{
			billing.ObjectCheckoutSession: related.CheckoutSessionID,
			billing.ObjectSubscription:    related.SubscriptionID,
			billing.ObjectInvoice:         related.InvoiceID,
			billing.ObjectPaymentIntent:   related.PaymentIntentID,
			"scenario_run":                event.Billtap.ScenarioRunID,
		}),
		CreatedAt: event.CreatedAt,
	}
}

func endpointSummary(endpoint webhooks.Endpoint) ObjectSummary {
	return ObjectSummary{ID: endpoint.ID, Object: TypeWebhookEndpoint, Display: endpoint.URL, Status: activeStatus(endpoint.Active), CreatedAt: endpoint.CreatedAt}
}

func checkoutSessionRelated(cs billing.CheckoutSession) map[string]string {
	return compactMap(map[string]string{
		billing.ObjectSubscription:  cs.SubscriptionID,
		billing.ObjectInvoice:       cs.InvoiceID,
		billing.ObjectPaymentIntent: cs.PaymentIntentID,
	})
}

func subscriptionRelated(sub billing.Subscription) map[string]string {
	return compactMap(map[string]string{
		billing.ObjectInvoice:         sub.LatestInvoiceID,
		billing.ObjectCheckoutSession: sub.Metadata["checkout_session"],
	})
}

func invoiceRelated(inv billing.Invoice) map[string]string {
	return compactMap(map[string]string{
		billing.ObjectSubscription:  inv.SubscriptionID,
		billing.ObjectPaymentIntent: inv.PaymentIntentID,
	})
}

func paymentIntentRelated(pi billing.PaymentIntent) map[string]string {
	return compactMap(map[string]string{billing.ObjectInvoice: pi.InvoiceID})
}

func compactMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		if value != "" {
			out[key] = value
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func activeStatus(active bool) string {
	if active {
		return "active"
	}
	return "inactive"
}

func first(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
