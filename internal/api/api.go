package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/fixtures"
	"github.com/hckim/billtap/internal/scenarios"
	"github.com/hckim/billtap/internal/security"
	"github.com/hckim/billtap/internal/webhooks"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Billing  *billing.Service
	Webhooks *webhooks.Service
}

type Handler struct {
	billing  *billing.Service
	webhooks *webhooks.Service
	mux      *http.ServeMux
}

func New(opts Options) http.Handler {
	h := &Handler{
		billing:  opts.Billing,
		webhooks: opts.Webhooks,
		mux:      http.NewServeMux(),
	}
	h.routes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) routes() {
	h.mux.HandleFunc("/v1/customers", h.handleCustomers)
	h.mux.HandleFunc("/v1/customers/", h.handleCustomer)
	h.mux.HandleFunc("/v1/products", h.handleProducts)
	h.mux.HandleFunc("/v1/products/search", h.handleProductSearch)
	h.mux.HandleFunc("/v1/products/", h.handleProduct)
	h.mux.HandleFunc("/v1/prices", h.handlePrices)
	h.mux.HandleFunc("/v1/prices/", h.handlePrice)
	h.mux.HandleFunc("/v1/checkout/sessions", h.handleCheckoutSessions)
	h.mux.HandleFunc("/v1/checkout/sessions/", h.handleCheckoutSession)
	h.mux.HandleFunc("/v1/billing_portal/sessions", h.handleBillingPortalSessions)
	h.mux.HandleFunc("/v1/subscriptions", h.handleSubscriptions)
	h.mux.HandleFunc("/v1/subscriptions/", h.handleSubscription)
	h.mux.HandleFunc("/v1/invoices/create_preview", h.handleInvoicePreview)
	h.mux.HandleFunc("/v1/invoices", h.handleInvoices)
	h.mux.HandleFunc("/v1/invoices/", h.handleInvoice)
	h.mux.HandleFunc("/v1/payment_intents", h.handlePaymentIntents)
	h.mux.HandleFunc("/v1/payment_intents/", h.handlePaymentIntent)
	h.mux.HandleFunc("/v1/payment_methods", h.handlePaymentMethods)
	h.mux.HandleFunc("/v1/webhook_endpoints", h.handleWebhookEndpoints)
	h.mux.HandleFunc("/v1/webhook_endpoints/", h.handleWebhookEndpoint)
	h.mux.HandleFunc("/v1/events", h.handleEvents)
	h.mux.HandleFunc("/v1/events/", h.handleEvent)
	h.mux.HandleFunc("/api/checkout/sessions/", h.handleCheckoutCompletion)
	h.mux.HandleFunc("/api/objects", h.handleObjects)
	h.mux.HandleFunc("/api/debug-bundles", h.handleDebugBundles)
	h.mux.HandleFunc("/api/fixtures/apply", h.handleFixtureApply)
	h.mux.HandleFunc("/api/fixtures/snapshot", h.handleFixtureSnapshot)
	h.mux.HandleFunc("/api/fixtures/assert", h.handleFixtureAssert)
	h.mux.HandleFunc("/api/portal", h.handlePortal)
	h.mux.HandleFunc("/api/portal/customers/", h.handlePortalCustomer)
	h.mux.HandleFunc("/api/portal/subscriptions/", h.handlePortalSubscription)
	h.mux.HandleFunc("/api/scenarios/run", h.handleScenarioRun)
	h.mux.HandleFunc("/api/timeline", h.handleTimeline)
	h.mux.HandleFunc("/api/delivery-attempts", h.handleDeliveryAttempts)
	h.mux.HandleFunc("/api/audit-log", h.handleAuditLog)
	h.mux.HandleFunc("/api/retention/apply", h.handleRetentionApply)
	h.mux.HandleFunc("/api/events/", h.handleEventAction)
}

func (h *Handler) handleCustomers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		customer, err := h.billing.CreateCustomer(r.Context(), billing.Customer{
			ID:       p.string("id"),
			Email:    p.string("email"),
			Name:     p.string("name"),
			Metadata: p.metadata(),
		})
		writeResult(w, stripeCustomer(customer), err)
	case http.MethodGet:
		customers, err := h.billing.ListCustomers(r.Context())
		data := make([]map[string]any, 0, len(customers))
		for _, customer := range customers {
			if email := r.URL.Query().Get("email"); email != "" && customer.Email != email {
				continue
			}
			data = append(data, stripeCustomer(customer))
			if limit := queryInt(r, "limit"); limit > 0 && len(data) >= limit {
				break
			}
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handleCustomer(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/customers/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		customer, err := h.billing.GetCustomer(r.Context(), id)
		writeResult(w, stripeCustomer(customer), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		customer, err := h.billing.UpdateCustomer(r.Context(), id, billing.Customer{
			Email:    p.string("email"),
			Name:     p.string("name"),
			Metadata: p.metadata(),
		})
		writeResult(w, stripeCustomer(customer), err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		product, err := h.billing.CreateProduct(r.Context(), billing.Product{
			ID:          p.string("id"),
			Name:        p.string("name"),
			Description: p.string("description"),
			Active:      p.boolDefault("active", true),
			Metadata:    p.metadata(),
		})
		writeResult(w, stripeProduct(product), err)
	case http.MethodGet:
		products, err := h.billing.ListProducts(r.Context())
		writeResult(w, stripeList(r.URL.Path, stripeProducts(products)), err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handleProductSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	products, err := h.billing.ListProducts(r.Context())
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	query := r.URL.Query().Get("query")
	filtered := filterProducts(products, query)
	writeResult(w, stripeSearchResult(r.URL.Path, query, stripeProducts(filtered)), nil)
}

func (h *Handler) handleProduct(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/products/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		product, err := h.billing.GetProduct(r.Context(), id)
		writeResult(w, stripeProduct(product), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		product, err := h.billing.UpdateProduct(r.Context(), id, billing.Product{
			Name:        p.string("name"),
			Description: p.string("description"),
			Active:      p.boolDefault("active", true),
			Metadata:    p.metadata(),
		})
		writeResult(w, stripeProduct(product), err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handlePrices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		price, err := h.billing.CreatePrice(r.Context(), billing.Price{
			ID:                     p.string("id"),
			ProductID:              p.first("product", "product_id"),
			Currency:               p.stringDefault("currency", "usd"),
			UnitAmount:             p.int64("unit_amount"),
			LookupKey:              p.first("lookup_key", "lookupKey"),
			RecurringInterval:      p.first("recurring[interval]", "recurring_interval", "interval"),
			RecurringIntervalCount: int(p.int64Default("recurring[interval_count]", 1)),
			Active:                 p.boolDefault("active", true),
			Metadata:               p.metadata(),
		})
		writeResult(w, stripePrice(price), err)
	case http.MethodGet:
		prices, err := h.billing.ListPrices(r.Context())
		writeResult(w, stripeList(r.URL.Path, stripePrices(filterPrices(prices, r))), err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handlePrice(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/prices/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		price, err := h.billing.GetPrice(r.Context(), id)
		writeResult(w, stripePrice(price), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		price, err := h.billing.UpdatePrice(r.Context(), id, billing.Price{
			ProductID:              p.first("product", "product_id"),
			Currency:               p.string("currency"),
			UnitAmount:             p.int64("unit_amount"),
			LookupKey:              p.first("lookup_key", "lookupKey"),
			RecurringInterval:      p.first("recurring[interval]", "recurring_interval", "interval"),
			RecurringIntervalCount: int(p.int64("recurring[interval_count]")),
			Active:                 p.boolDefault("active", true),
			Metadata:               p.metadata(),
		})
		writeResult(w, stripePrice(price), err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handleCheckoutSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		session, err := h.billing.CreateCheckoutSession(r.Context(), billing.CheckoutSession{
			CustomerID: p.first("customer", "customer_id"),
			Mode:       p.stringDefault("mode", "subscription"),
			LineItems:  p.lineItems(),
			SuccessURL: p.string("success_url"),
			CancelURL:  p.string("cancel_url"),
		})
		if err == nil {
			session.URL = absoluteURL(r, session.URL)
		}
		writeResult(w, stripeCheckoutSession(session), err)
	case http.MethodGet:
		sessions, err := h.billing.ListCheckoutSessions(r.Context())
		data := make([]map[string]any, 0, len(sessions))
		for i := range sessions {
			sessions[i].URL = absoluteURL(r, sessions[i].URL)
			data = append(data, stripeCheckoutSession(sessions[i]))
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handleCheckoutSession(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/v1/checkout/sessions/")
	if strings.HasSuffix(rest, "/complete") {
		id := strings.TrimSuffix(rest, "/complete")
		h.completeCheckout(w, r, id)
		return
	}
	if rest == "" || strings.Contains(rest, "/") {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	session, err := h.billing.GetCheckoutSession(r.Context(), rest)
	if err == nil {
		session.URL = absoluteURL(r, session.URL)
	}
	writeResult(w, stripeCheckoutSession(session), err)
}

func (h *Handler) handleBillingPortalSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	customerID := p.first("customer", "customer_id")
	if customerID == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("%w: customer is required", billing.ErrInvalidInput))
		return
	}
	if _, err := h.billing.GetCustomer(r.Context(), customerID); err != nil {
		writeResult(w, nil, err)
		return
	}
	session := map[string]any{
		"id":         "bps_" + sanitizeID(customerID+"_"+time.Now().UTC().Format(time.RFC3339Nano)),
		"object":     "billing_portal.session",
		"customer":   customerID,
		"return_url": p.string("return_url"),
		"url":        absoluteURL(r, "/portal?customer_id="+customerID),
		"created":    time.Now().UTC().Unix(),
		"livemode":   false,
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *Handler) handleCheckoutCompletion(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/checkout/sessions/")
	id, suffix, _ := strings.Cut(rest, "/")
	if id == "" || suffix != "complete" {
		http.NotFound(w, r)
		return
	}
	h.completeCheckout(w, r, id)
}

func (h *Handler) completeCheckout(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	previous, _ := h.billing.GetCheckoutSession(r.Context(), id)
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	session, err := h.billing.CompleteCheckout(r.Context(), id, p.stringDefault("outcome", "payment_succeeded"))
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	session.URL = absoluteURL(r, session.URL)
	result := map[string]any{"session": session}
	if session.SubscriptionID != "" {
		if sub, err := h.billing.GetSubscription(r.Context(), session.SubscriptionID); err == nil {
			result["subscription"] = sub
		}
	}
	if session.InvoiceID != "" {
		if invoice, err := h.billing.GetInvoice(r.Context(), session.InvoiceID); err == nil {
			result["invoice"] = invoice
		}
	}
	if session.PaymentIntentID != "" {
		if pi, err := h.billing.GetPaymentIntent(r.Context(), session.PaymentIntentID); err == nil {
			result["payment_intent"] = pi
		}
	}
	if h.webhooks != nil && previous.PaymentIntentID == "" {
		result["webhook_events"] = h.emitCheckoutWebhooks(r, result)
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) handleSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	items, err := h.billing.ListSubscriptions(r.Context())
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	filtered := filterSubscriptions(items, r)
	data := make([]map[string]any, 0, len(filtered))
	for _, item := range filtered {
		data = append(data, h.stripeSubscription(r, item))
		if limit := queryInt(r, "limit"); limit > 0 && len(data) >= limit {
			break
		}
	}
	writeResult(w, stripeList(r.URL.Path, data), nil)
}

func (h *Handler) handleSubscription(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/subscriptions/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		subscription, err := h.billing.GetSubscription(r.Context(), id)
		writeResult(w, h.stripeSubscription(r, subscription), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		replaceItems := hasSubscriptionItemPatch(p)
		var items []billing.LineItem
		if replaceItems {
			current, err := h.billing.GetSubscription(r.Context(), id)
			if err != nil {
				writeResult(w, nil, err)
				return
			}
			items = subscriptionItemsFromParams(p, current)
		}
		subscription, err := h.billing.PatchSubscription(r.Context(), id, billing.SubscriptionPatch{
			Items:             items,
			ReplaceItems:      replaceItems,
			Metadata:          p.metadata(),
			CancelAtPeriodEnd: p.boolPtr("cancel_at_period_end"),
		})
		writeResult(w, h.stripeSubscription(r, subscription), err)
	case http.MethodDelete:
		subscription, err := h.billing.CancelPortalSubscription(r.Context(), id, billing.PortalCancel{Mode: "immediate"})
		writeResult(w, h.stripeSubscription(r, subscription), err)
	default:
		methodNotAllowed(w, "GET, POST, DELETE")
	}
}

func (h *Handler) handleInvoicePreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	now := time.Now().UTC().Unix()
	writeJSON(w, http.StatusOK, map[string]any{
		"id":          "upcoming_in_" + strconv.FormatInt(now, 10),
		"object":      "invoice",
		"amount_due":  0,
		"subtotal":    0,
		"total":       0,
		"currency":    "usd",
		"created":     now,
		"status":      "draft",
		"lines":       stripeList("/v1/invoices/create_preview/lines", []map[string]any{}),
		"livemode":    false,
		"description": "Billtap preview uses zero-value proration for local smoke tests",
	})
}

func (h *Handler) handleInvoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	items, err := h.billing.ListInvoices(r.Context())
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	filtered := filterInvoices(items, r)
	data := make([]map[string]any, 0, len(filtered))
	for _, item := range filtered {
		data = append(data, stripeInvoice(item))
	}
	writeResult(w, stripeList(r.URL.Path, data), nil)
}

func (h *Handler) handleInvoice(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/invoices/")
	if id == "" || strings.Contains(id, "/") || r.Method != http.MethodGet {
		if r.Method != http.MethodGet {
			methodNotAllowed(w, "GET")
			return
		}
		http.NotFound(w, r)
		return
	}
	invoice, err := h.billing.GetInvoice(r.Context(), id)
	writeResult(w, stripeInvoice(invoice), err)
}

func (h *Handler) handlePaymentIntents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	items, err := h.billing.ListPaymentIntents(r.Context())
	data := make([]map[string]any, 0, len(items))
	for _, item := range items {
		data = append(data, stripePaymentIntent(item))
	}
	writeResult(w, stripeList(r.URL.Path, data), err)
}

func (h *Handler) handlePaymentIntent(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/payment_intents/")
	if id == "" || strings.Contains(id, "/") || r.Method != http.MethodGet {
		if r.Method != http.MethodGet {
			methodNotAllowed(w, "GET")
			return
		}
		http.NotFound(w, r)
		return
	}
	paymentIntent, err := h.billing.GetPaymentIntent(r.Context(), id)
	writeResult(w, stripePaymentIntent(paymentIntent), err)
}

func (h *Handler) handlePaymentMethods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	customerID := r.URL.Query().Get("customer")
	if customerID == "" {
		writeResult(w, stripeList(r.URL.Path, []map[string]any{}), nil)
		return
	}
	if _, err := h.billing.GetCustomer(r.Context(), customerID); err != nil {
		writeResult(w, nil, err)
		return
	}
	method := map[string]any{
		"id":       "pm_" + sanitizeID(customerID),
		"object":   "payment_method",
		"customer": customerID,
		"type":     "card",
		"card": map[string]any{
			"brand":     "visa",
			"last4":     "4242",
			"exp_month": 12,
			"exp_year":  2035,
		},
		"created":  time.Now().UTC().Unix(),
		"livemode": false,
	}
	writeResult(w, stripeList(r.URL.Path, []map[string]any{method}), nil)
}

func (h *Handler) handleObjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	objects, err := h.collectObjects(r)
	writeResult(w, objects, err)
}

func (h *Handler) handleWebhookEndpoints(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		endpoint, err := h.webhooks.CreateEndpoint(r.Context(), webhooks.Endpoint{
			URL:              p.string("url"),
			Secret:           p.string("secret"),
			EnabledEvents:    p.list("enabled_events"),
			RetryMaxAttempts: int(p.int64("retry_max_attempts")),
			RetryBackoff:     p.list("retry_backoff"),
		})
		writeResult(w, maskEndpoint(endpoint), err)
	case http.MethodGet:
		endpoints, err := h.webhooks.ListEndpoints(r.Context(), webhooks.EndpointFilter{})
		writeResult(w, map[string]any{"object": "list", "data": maskEndpoints(endpoints)}, err)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) handleWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/webhook_endpoints/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		endpoint, err := h.webhooks.GetEndpoint(r.Context(), id)
		writeResult(w, maskEndpoint(endpoint), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		endpoint, err := h.webhooks.UpdateEndpoint(r.Context(), id, webhooks.Endpoint{
			URL:              p.string("url"),
			Secret:           p.string("secret"),
			EnabledEvents:    p.list("enabled_events"),
			RetryMaxAttempts: int(p.int64("retry_max_attempts")),
			RetryBackoff:     p.list("retry_backoff"),
			Active:           p.boolDefault("active", true),
		})
		writeResult(w, maskEndpoint(endpoint), err)
	case http.MethodDelete:
		endpoint, err := h.webhooks.DeleteEndpoint(r.Context(), id)
		writeResult(w, maskEndpoint(endpoint), err)
	default:
		methodNotAllowed(w, "GET, POST, DELETE")
	}
}

func (h *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeJSON(w, http.StatusOK, map[string]any{"object": "list", "data": []any{}})
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	q := r.URL.Query()
	events, err := h.webhooks.ListEvents(r.Context(), webhooks.EventFilter{
		Type:          q.Get("type"),
		ScenarioRunID: q.Get("scenarioRunId"),
	})
	writeResult(w, map[string]any{"object": "list", "data": events}, err)
}

func (h *Handler) handleEvent(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/events/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	event, err := h.webhooks.GetEvent(r.Context(), id)
	writeResult(w, event, err)
}

func (h *Handler) handleTimeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	q := r.URL.Query()
	entries, err := h.billing.Timeline(r.Context(), billing.TimelineFilter{
		CustomerID:        q.Get("customerId"),
		CheckoutSessionID: q.Get("checkoutSessionId"),
		SubscriptionID:    q.Get("subscriptionId"),
		InvoiceID:         q.Get("invoiceId"),
		PaymentIntentID:   q.Get("paymentIntentId"),
	})
	writeResult(w, map[string]any{"object": "list", "data": entries}, err)
}

func (h *Handler) handleDebugBundles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	filter := debugBundleTimelineFilter(p)
	timeline, err := h.billing.Timeline(r.Context(), filter)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	objects, err := h.collectObjects(r)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	var events []webhooks.Event
	var attempts []webhooks.DeliveryAttempt
	if h.webhooks != nil {
		events, _ = h.webhooks.ListEvents(r.Context(), webhooks.EventFilter{})
		attempts, _ = h.webhooks.ListDeliveryAttempts(r.Context(), webhooks.DeliveryAttemptFilter{})
	}
	targetType := p.first("objectType", "object_type", "targetType", "target_type", "type")
	targetID := p.first("objectId", "object_id", "targetId", "target_id", "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"id":                "dbg_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		"object":            "debug_bundle",
		"target":            compactStringMap(map[string]string{"type": targetType, "id": targetID}),
		"filters":           timelineFilterMap(filter),
		"objects":           objects,
		"timeline":          timeline,
		"webhook_events":    events,
		"delivery_attempts": h.deliveryAttemptResponses(r, attempts),
		"created_at":        time.Now().UTC(),
	})
}

func (h *Handler) handleFixtureApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	body, err := readSafeFixtureBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	pack, err := fixtures.LoadPack(body, r.Header.Get("Content-Type"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("%w: %v", billing.ErrInvalidInput, err))
		return
	}
	result, err := fixtures.NewService(h.billing).Apply(r.Context(), pack)
	if err != nil {
		if errors.Is(err, fixtures.ErrAssertionFailed) {
			writeJSON(w, http.StatusConflict, result)
			return
		}
		writeResult(w, result, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) handleFixtureSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	snapshot, err := fixtures.NewService(h.billing).Snapshot(r.Context(), fixtureSnapshotFilter(r))
	writeResult(w, snapshot, err)
}

func (h *Handler) handleFixtureAssert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	body, err := readSafeFixtureBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req, err := fixtures.LoadAssertionRequest(body, r.Header.Get("Content-Type"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("%w: %v", billing.ErrInvalidInput, err))
		return
	}
	report, err := fixtures.NewService(h.billing).Assert(r.Context(), req)
	if err != nil {
		if errors.Is(err, fixtures.ErrAssertionFailed) {
			writeJSON(w, http.StatusConflict, report)
			return
		}
		writeResult(w, report, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func readSafeFixtureBody(r *http.Request) ([]byte, error) {
	body := mustReadRequestBody(r)
	if security.ContainsCardDataAny(decodeLooseBody(body, r.Header.Get("Content-Type"))) {
		return nil, fmt.Errorf("%w: real card data is not accepted by Billtap", webhooks.ErrInvalidInput)
	}
	return body, nil
}

func decodeLooseBody(body []byte, contentType string) any {
	var decoded any
	if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "yml") {
		if err := yaml.Unmarshal(body, &decoded); err == nil {
			return decoded
		}
		return nil
	}
	if err := json.Unmarshal(body, &decoded); err == nil {
		return decoded
	}
	return nil
}

func fixtureSnapshotFilter(r *http.Request) fixtures.SnapshotFilter {
	return fixtures.SnapshotFilter{
		CustomerID:  firstQuery(r, "customer", "customerId", "customer_id"),
		RunID:       firstQuery(r, "runId", "run_id"),
		TenantID:    firstQuery(r, "tenantId", "tenant_id"),
		FixtureName: firstQuery(r, "fixture", "fixtureName", "fixture_name", "name"),
		Namespace:   firstQuery(r, "namespace", "ns"),
	}
}

func debugBundleTimelineFilter(p params) billing.TimelineFilter {
	filter := billing.TimelineFilter{
		CustomerID:        p.first("customerId", "customer_id"),
		CheckoutSessionID: p.first("checkoutSessionId", "checkout_session_id"),
		SubscriptionID:    p.first("subscriptionId", "subscription_id"),
		InvoiceID:         p.first("invoiceId", "invoice_id"),
		PaymentIntentID:   p.first("paymentIntentId", "payment_intent_id"),
	}

	objectID := p.first("objectId", "object_id", "targetId", "target_id", "id")
	if objectID == "" {
		return filter
	}
	switch dashboardObjectType(p.first("objectType", "object_type", "targetType", "target_type", "type")) {
	case "customer":
		if filter.CustomerID == "" {
			filter.CustomerID = objectID
		}
	case "checkout_session":
		if filter.CheckoutSessionID == "" {
			filter.CheckoutSessionID = objectID
		}
	case "subscription":
		if filter.SubscriptionID == "" {
			filter.SubscriptionID = objectID
		}
	case "invoice":
		if filter.InvoiceID == "" {
			filter.InvoiceID = objectID
		}
	case "payment_intent":
		if filter.PaymentIntentID == "" {
			filter.PaymentIntentID = objectID
		}
	default:
		switch {
		case strings.HasPrefix(objectID, "cus_") && filter.CustomerID == "":
			filter.CustomerID = objectID
		case strings.HasPrefix(objectID, "cs_") && filter.CheckoutSessionID == "":
			filter.CheckoutSessionID = objectID
		case strings.HasPrefix(objectID, "sub_") && filter.SubscriptionID == "":
			filter.SubscriptionID = objectID
		case strings.HasPrefix(objectID, "in_") && filter.InvoiceID == "":
			filter.InvoiceID = objectID
		case strings.HasPrefix(objectID, "pi_") && filter.PaymentIntentID == "":
			filter.PaymentIntentID = objectID
		}
	}
	return filter
}

func dashboardObjectType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "-", "_")
	switch value {
	case "customers", "customer":
		return "customer"
	case "checkout_sessions", "checkoutsessions", "checkout_session", "checkout.session":
		return "checkout_session"
	case "subscriptions", "subscription":
		return "subscription"
	case "invoices", "invoice":
		return "invoice"
	case "payment_intents", "paymentintents", "payment_intent", "payment.intent":
		return "payment_intent"
	default:
		return value
	}
}

func timelineFilterMap(filter billing.TimelineFilter) map[string]string {
	return compactStringMap(map[string]string{
		"customer_id":         filter.CustomerID,
		"checkout_session_id": filter.CheckoutSessionID,
		"subscription_id":     filter.SubscriptionID,
		"invoice_id":          filter.InvoiceID,
		"payment_intent_id":   filter.PaymentIntentID,
	})
}

func compactStringMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range in {
		if value != "" {
			out[key] = value
		}
	}
	return out
}

func (h *Handler) handlePortal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	customerID := firstQuery(r, "customerId", "customer_id", "customer", "id")
	if customerID == "" {
		customers, err := h.billing.ListCustomers(r.Context())
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		if len(customers) == 0 {
			writeResult(w, nil, fmt.Errorf("%w: customer is required", billing.ErrInvalidInput))
			return
		}
		customerID = customers[0].ID
	}
	state, err := h.billing.PortalState(r.Context(), customerID)
	writeResult(w, state, err)
}

func (h *Handler) handlePortalCustomer(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/portal/customers/")
	customerID, action, hasAction := strings.Cut(rest, "/")
	if customerID == "" {
		http.NotFound(w, r)
		return
	}
	if !hasAction {
		if r.Method != http.MethodGet {
			methodNotAllowed(w, "GET")
			return
		}
		state, err := h.billing.PortalState(r.Context(), customerID)
		writeResult(w, state, err)
		return
	}
	if action != "payment-method" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := h.billing.SimulatePaymentMethodUpdate(r.Context(), customerID, p.first("outcome", "simulate", "status"))
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	state, stateErr := h.billing.PortalState(r.Context(), customerID)
	writeResult(w, map[string]any{"object": "portal_action", "action": "payment_method", "payment_method": result, "state": state}, stateErr)
}

func (h *Handler) handlePortalSubscription(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/portal/subscriptions/")
	subscriptionID, action, found := strings.Cut(rest, "/")
	if subscriptionID == "" || !found {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var sub billing.Subscription
	switch action {
	case "plan-change":
		sub, err = h.billing.ChangePortalPlan(r.Context(), subscriptionID, billing.PortalPlanChange{
			PlanID:   p.first("plan", "plan_id", "planId"),
			PriceID:  p.first("price", "price_id", "priceId"),
			Quantity: p.int64Default("quantity", p.int64("seats")),
		})
	case "seat-change":
		sub, err = h.billing.ChangePortalSeats(r.Context(), subscriptionID, billing.PortalSeatChange{
			Quantity: p.int64Default("quantity", p.int64("seats")),
		})
	case "cancel":
		sub, err = h.billing.CancelPortalSubscription(r.Context(), subscriptionID, billing.PortalCancel{
			Mode: p.stringDefault("mode", "period"),
		})
	case "resume":
		sub, err = h.billing.ResumePortalSubscription(r.Context(), subscriptionID)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	state, stateErr := h.billing.PortalState(r.Context(), sub.CustomerID)
	writeResult(w, map[string]any{"object": "portal_action", "action": action, "subscription": sub, "state": state}, stateErr)
}

func firstQuery(r *http.Request, keys ...string) string {
	q := r.URL.Query()
	for _, key := range keys {
		if value := strings.TrimSpace(q.Get(key)); value != "" {
			return value
		}
	}
	return ""
}

func (h *Handler) handleDeliveryAttempts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	if h.webhooks == nil {
		writeJSON(w, http.StatusOK, map[string]any{"object": "list", "data": []any{}})
		return
	}
	q := r.URL.Query()
	attempts, err := h.webhooks.ListDeliveryAttempts(r.Context(), webhooks.DeliveryAttemptFilter{
		EventID:    q.Get("eventId"),
		EndpointID: q.Get("endpointId"),
		Status:     q.Get("status"),
	})
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	writeResult(w, map[string]any{"object": "list", "data": h.deliveryAttemptResponses(r, attempts)}, nil)
}

func (h *Handler) handleAuditLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	if h.webhooks == nil {
		writeJSON(w, http.StatusOK, map[string]any{"object": "list", "data": []any{}})
		return
	}
	entries, err := h.webhooks.ListAuditEntries(r.Context(), webhooks.AuditFilter{
		Action:   r.URL.Query().Get("action"),
		TargetID: r.URL.Query().Get("targetId"),
	})
	writeResult(w, map[string]any{"object": "list", "data": entries}, err)
}

func (h *Handler) handleRetentionApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	result, err := h.webhooks.ApplyRetention(r.Context())
	writeResult(w, result, err)
}

func (h *Handler) handleScenarioRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var scenario scenarios.Scenario
	if strings.Contains(r.Header.Get("Content-Type"), "yaml") || strings.Contains(r.Header.Get("Content-Type"), "yml") {
		loaded, err := scenarios.Load(mustReadRequestBody(r))
		if err != nil {
			report := scenarioErrorReport("request", scenarios.FailureInvalidConfig, err)
			writeJSON(w, http.StatusBadRequest, report)
			return
		}
		scenario = loaded
	} else if err := json.NewDecoder(r.Body).Decode(&scenario); err != nil {
		report := scenarioErrorReport("request", scenarios.FailureInvalidConfig, err)
		writeJSON(w, http.StatusBadRequest, report)
		return
	}

	runner := scenarios.NewRunner(h.billing, h.webhooks)
	report, err := runner.Run(r.Context(), scenario)
	status := http.StatusOK
	if err != nil {
		switch report.ExitCode() {
		case scenarios.ExitInvalidConfig:
			status = http.StatusBadRequest
		case scenarios.ExitAppCallbackFailure:
			status = http.StatusBadGateway
		default:
			status = http.StatusConflict
		}
	}
	writeJSON(w, status, report)
}

func scenarioErrorReport(name string, failureType string, err error) scenarios.Report {
	now := time.Now().UTC()
	return scenarios.Report{
		Name:        name,
		Status:      "failed",
		FailureType: failureType,
		StartedAt:   now,
		FinishedAt:  now,
		Errors:      []string{err.Error()},
	}
}

func (h *Handler) handleEventAction(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/api/events/")
	id, action, found := strings.Cut(rest, "/")
	if id == "" || !found || action != "replay" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	delay := time.Duration(p.int64("delay_seconds")) * time.Second
	if delay == 0 {
		delay = time.Duration(p.int64("delay")) * time.Second
	}
	attempts, err := h.webhooks.ReplayEvent(r.Context(), id, webhooks.ReplayOptions{
		Duplicate:  int(p.int64Default("duplicate", 1)),
		Delay:      delay,
		OutOfOrder: p.boolDefault("out_of_order", false),
	})
	writeResult(w, map[string]any{"message": "replay scheduled", "object": "list", "data": h.deliveryAttemptResponses(r, attempts)}, err)
}

type params struct {
	values map[string]string
}

func mustReadRequestBody(r *http.Request) []byte {
	body, _ := io.ReadAll(r.Body)
	return body
}

func parseParams(r *http.Request) (params, error) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var raw map[string]any
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			return params{}, err
		}
		if security.ContainsCardDataAny(raw) {
			return params{}, fmt.Errorf("%w: real card data is not accepted by Billtap", webhooks.ErrInvalidInput)
		}
		values := map[string]string{}
		for key, value := range raw {
			switch v := value.(type) {
			case string:
				values[key] = v
			case float64:
				values[key] = strconv.FormatInt(int64(v), 10)
			case bool:
				values[key] = strconv.FormatBool(v)
			case map[string]any:
				for nestedKey, nested := range v {
					values[fmt.Sprintf("%s[%s]", key, nestedKey)] = fmt.Sprint(nested)
				}
			case []any:
				for idx, item := range v {
					if record, ok := item.(map[string]any); ok {
						for k, nested := range record {
							values[fmt.Sprintf("%s[%d][%s]", key, idx, k)] = fmt.Sprint(nested)
						}
					}
				}
			}
		}
		if security.ContainsCardData(values) {
			return params{}, fmt.Errorf("%w: real card data is not accepted by Billtap", webhooks.ErrInvalidInput)
		}
		return params{values: values}, nil
	}
	if err := r.ParseForm(); err != nil {
		return params{}, err
	}
	values := map[string]string{}
	for key, value := range r.Form {
		if len(value) > 0 {
			values[key] = value[0]
		}
	}
	if security.ContainsCardData(values) {
		return params{}, fmt.Errorf("%w: real card data is not accepted by Billtap", webhooks.ErrInvalidInput)
	}
	return params{values: values}, nil
}

func (p params) string(key string) string {
	return strings.TrimSpace(p.values[key])
}

func (p params) stringDefault(key string, fallback string) string {
	if value := p.string(key); value != "" {
		return value
	}
	return fallback
}

func (p params) first(keys ...string) string {
	for _, key := range keys {
		if value := p.string(key); value != "" {
			return value
		}
	}
	return ""
}

func (p params) int64(key string) int64 {
	value, _ := strconv.ParseInt(p.string(key), 10, 64)
	return value
}

func (p params) int64Default(key string, fallback int64) int64 {
	if value := p.int64(key); value != 0 {
		return value
	}
	return fallback
}

func (p params) boolDefault(key string, fallback bool) bool {
	value := p.string(key)
	if value == "" {
		return fallback
	}
	return value == "true" || value == "1"
}

func (p params) boolPtr(key string) *bool {
	value := p.string(key)
	if value == "" {
		return nil
	}
	result := value == "true" || value == "1"
	return &result
}

var metadataParamPattern = regexp.MustCompile(`^metadata\[([^\]]+)\]$`)

func (p params) metadata() map[string]string {
	out := map[string]string{}
	for key, value := range p.values {
		matches := metadataParamPattern.FindStringSubmatch(key)
		if len(matches) != 2 {
			continue
		}
		out[matches[1]] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func (p params) list(key string) []string {
	var out []string
	for rawKey, value := range p.values {
		if rawKey == key || strings.HasPrefix(rawKey, key+"[") {
			for _, part := range strings.Split(value, ",") {
				part = strings.TrimSpace(part)
				if part != "" {
					out = append(out, part)
				}
			}
		}
	}
	return out
}

func (p params) lineItems() []billing.LineItem {
	var out []billing.LineItem
	for i := 0; i < 100; i++ {
		price := p.first(fmt.Sprintf("line_items[%d][price]", i), fmt.Sprintf("lineItems[%d][price]", i))
		if price == "" {
			if i == 0 {
				price = p.string("price")
			}
			if price == "" {
				continue
			}
		}
		quantity := p.int64Default(fmt.Sprintf("line_items[%d][quantity]", i), 1)
		out = append(out, billing.LineItem{PriceID: price, Quantity: quantity})
	}
	return out
}

func queryInt(r *http.Request, key string) int {
	value, _ := strconv.Atoi(r.URL.Query().Get(key))
	return value
}

func stripeList(urlPath string, data any) map[string]any {
	return map[string]any{
		"object":   "list",
		"url":      urlPath,
		"has_more": false,
		"data":     data,
	}
}

func stripeSearchResult(urlPath string, query string, data any) map[string]any {
	return map[string]any{
		"object":    "search_result",
		"url":       urlPath + "?query=" + query,
		"has_more":  false,
		"next_page": nil,
		"data":      data,
	}
}

func stripeCustomer(customer billing.Customer) map[string]any {
	return map[string]any{
		"id":         customer.ID,
		"object":     billing.ObjectCustomer,
		"email":      customer.Email,
		"name":       customer.Name,
		"metadata":   nonNilMap(customer.Metadata),
		"created":    unix(customer.CreatedAt),
		"created_at": customer.CreatedAt,
		"invoice_settings": map[string]any{
			"default_payment_method": nil,
		},
		"livemode": false,
	}
}

func stripeProduct(product billing.Product) map[string]any {
	return map[string]any{
		"id":            product.ID,
		"object":        billing.ObjectProduct,
		"active":        product.Active,
		"created":       unix(product.CreatedAt),
		"created_at":    product.CreatedAt,
		"default_price": product.Metadata["default_price"],
		"description":   product.Description,
		"images":        []string{},
		"livemode":      false,
		"metadata":      nonNilMap(product.Metadata),
		"name":          product.Name,
		"type":          "service",
		"updated":       unix(product.CreatedAt),
	}
}

func stripeProducts(products []billing.Product) []map[string]any {
	out := make([]map[string]any, 0, len(products))
	for _, product := range products {
		out = append(out, stripeProduct(product))
	}
	return out
}

func stripePrice(price billing.Price) map[string]any {
	priceType := "one_time"
	var recurring any
	if price.RecurringInterval != "" {
		priceType = "recurring"
		recurring = map[string]any{
			"interval":       price.RecurringInterval,
			"interval_count": price.RecurringIntervalCount,
			"usage_type":     "licensed",
		}
	}
	lookupKey := price.LookupKey
	if lookupKey == "" {
		lookupKey = price.Metadata["lookup_key"]
	}
	return map[string]any{
		"id":                       price.ID,
		"object":                   billing.ObjectPrice,
		"active":                   price.Active,
		"billing_scheme":           "per_unit",
		"created":                  unix(price.CreatedAt),
		"created_at":               price.CreatedAt,
		"currency":                 strings.ToLower(price.Currency),
		"livemode":                 false,
		"lookup_key":               lookupKey,
		"metadata":                 nonNilMap(price.Metadata),
		"product":                  price.ProductID,
		"recurring":                recurring,
		"recurring_interval":       price.RecurringInterval,
		"recurring_interval_count": price.RecurringIntervalCount,
		"tax_behavior":             "unspecified",
		"type":                     priceType,
		"unit_amount":              price.UnitAmount,
		"unit_amount_decimal":      strconv.FormatInt(price.UnitAmount, 10),
	}
}

func stripePrices(prices []billing.Price) []map[string]any {
	out := make([]map[string]any, 0, len(prices))
	for _, price := range prices {
		out = append(out, stripePrice(price))
	}
	return out
}

func stripeCheckoutSession(session billing.CheckoutSession) map[string]any {
	return map[string]any{
		"id":             session.ID,
		"object":         billing.ObjectCheckoutSession,
		"customer":       session.CustomerID,
		"mode":           session.Mode,
		"line_items":     session.LineItems,
		"success_url":    session.SuccessURL,
		"cancel_url":     session.CancelURL,
		"url":            session.URL,
		"status":         session.Status,
		"payment_status": session.PaymentStatus,
		"subscription":   emptyToNil(session.SubscriptionID),
		"invoice":        emptyToNil(session.InvoiceID),
		"payment_intent": emptyToNil(session.PaymentIntentID),
		"created":        unix(session.CreatedAt),
		"created_at":     session.CreatedAt,
		"completed_at":   session.CompletedAt,
		"livemode":       false,
	}
}

func (h *Handler) stripeSubscription(r *http.Request, sub billing.Subscription) map[string]any {
	items := make([]map[string]any, 0, len(sub.Items))
	for idx, item := range sub.Items {
		items = append(items, h.stripeSubscriptionItem(r, sub, item, idx))
	}
	return map[string]any{
		"id":                   sub.ID,
		"object":               billing.ObjectSubscription,
		"customer":             sub.CustomerID,
		"status":               sub.Status,
		"items":                stripeList("/v1/subscription_items?subscription="+sub.ID, items),
		"current_period_start": unix(sub.CurrentPeriodStart),
		"current_period_end":   unix(sub.CurrentPeriodEnd),
		"start_date":           unix(sub.CurrentPeriodStart),
		"cancel_at_period_end": sub.CancelAtPeriodEnd,
		"canceled_at":          optionalUnix(sub.CanceledAt),
		"cancel_at":            nil,
		"ended_at":             nil,
		"trial_start":          nil,
		"trial_end":            nil,
		"latest_invoice":       emptyToNil(sub.LatestInvoiceID),
		"metadata":             nonNilMap(sub.Metadata),
		"collection_method":    stringDefault(sub.Metadata["collection_method"], "charge_automatically"),
		"billing_cycle_anchor": unix(sub.CurrentPeriodStart),
		"currency":             "usd",
		"livemode":             false,
		"pause_collection":     nil,
		"pending_update":       nil,
		"cancellation_details": map[string]any{"comment": nil, "feedback": nil, "reason": nil},
	}
}

func (h *Handler) stripeSubscriptionItem(r *http.Request, sub billing.Subscription, item billing.LineItem, idx int) map[string]any {
	price, err := h.billing.GetPrice(r.Context(), item.PriceID)
	var priceObject map[string]any
	if err == nil {
		priceObject = stripePrice(price)
		if product, productErr := h.billing.GetProduct(r.Context(), price.ProductID); productErr == nil {
			priceObject["product"] = stripeProduct(product)
		}
	} else {
		priceObject = map[string]any{
			"id":                  item.PriceID,
			"object":              billing.ObjectPrice,
			"currency":            "usd",
			"unit_amount":         0,
			"unit_amount_decimal": "0",
			"recurring":           map[string]any{"interval": "month", "interval_count": 1, "usage_type": "licensed"},
			"type":                "recurring",
			"product":             "",
			"livemode":            false,
		}
	}
	return map[string]any{
		"id":                   "si_" + sanitizeID(sub.ID+"_"+strconv.Itoa(idx)),
		"object":               "subscription_item",
		"subscription":         sub.ID,
		"price":                priceObject,
		"quantity":             item.Quantity,
		"created":              unix(sub.CurrentPeriodStart),
		"current_period_start": unix(sub.CurrentPeriodStart),
		"current_period_end":   unix(sub.CurrentPeriodEnd),
		"metadata":             map[string]string{},
	}
}

func stripeInvoice(invoice billing.Invoice) map[string]any {
	return map[string]any{
		"id":             invoice.ID,
		"object":         billing.ObjectInvoice,
		"customer":       invoice.CustomerID,
		"subscription":   emptyToNil(invoice.SubscriptionID),
		"parent":         map[string]any{"subscription_details": map[string]any{"subscription": emptyToNil(invoice.SubscriptionID)}},
		"status":         invoice.Status,
		"currency":       invoice.Currency,
		"subtotal":       invoice.Subtotal,
		"total":          invoice.Total,
		"amount_due":     invoice.AmountDue,
		"amount_paid":    invoice.AmountPaid,
		"attempt_count":  invoice.AttemptCount,
		"payment_intent": emptyToNil(invoice.PaymentIntentID),
		"payments":       stripeList("/v1/invoices/"+invoice.ID+"/payments", []map[string]any{}),
		"lines":          stripeList("/v1/invoices/"+invoice.ID+"/lines", []map[string]any{}),
		"created":        unix(invoice.CreatedAt),
		"created_at":     invoice.CreatedAt,
		"status_transitions": map[string]any{
			"paid_at": optionalPaidAt(invoice),
		},
		"hosted_invoice_url": "",
		"livemode":           false,
	}
}

func stripePaymentIntent(intent billing.PaymentIntent) map[string]any {
	return map[string]any{
		"id":                 intent.ID,
		"object":             billing.ObjectPaymentIntent,
		"customer":           intent.CustomerID,
		"invoice":            emptyToNil(intent.InvoiceID),
		"amount":             intent.Amount,
		"currency":           intent.Currency,
		"status":             intent.Status,
		"last_payment_error": paymentIntentError(intent),
		"client_secret":      intent.ID + "_secret_billtap",
		"created":            unix(intent.CreatedAt),
		"created_at":         intent.CreatedAt,
		"livemode":           false,
	}
}

func filterPrices(prices []billing.Price, r *http.Request) []billing.Price {
	query := r.URL.Query()
	out := make([]billing.Price, 0, len(prices))
	for _, price := range prices {
		if product := query.Get("product"); product != "" && price.ProductID != product {
			continue
		}
		if active := query.Get("active"); active != "" && price.Active != (active == "true" || active == "1") {
			continue
		}
		if priceType := query.Get("type"); priceType == "recurring" && price.RecurringInterval == "" {
			continue
		}
		out = append(out, price)
	}
	return out
}

func filterSubscriptions(items []billing.Subscription, r *http.Request) []billing.Subscription {
	query := r.URL.Query()
	out := make([]billing.Subscription, 0, len(items))
	for _, item := range items {
		if customer := query.Get("customer"); customer != "" && item.CustomerID != customer {
			continue
		}
		status := strings.ToLower(query.Get("status"))
		if status != "" && status != "all" && item.Status != status {
			continue
		}
		out = append(out, item)
	}
	return out
}

func filterInvoices(items []billing.Invoice, r *http.Request) []billing.Invoice {
	query := r.URL.Query()
	out := make([]billing.Invoice, 0, len(items))
	for _, item := range items {
		if customer := query.Get("customer"); customer != "" && item.CustomerID != customer {
			continue
		}
		if subscription := query.Get("subscription"); subscription != "" && item.SubscriptionID != subscription {
			continue
		}
		if status := query.Get("status"); status != "" && item.Status != status {
			continue
		}
		out = append(out, item)
	}
	return out
}

var (
	searchMetadataPattern = regexp.MustCompile(`metadata\['([^']+)'\]:'([^']*)'`)
	searchActivePattern   = regexp.MustCompile(`active:'(true|false)'`)
)

func filterProducts(products []billing.Product, query string) []billing.Product {
	metadata := map[string]string{}
	for _, match := range searchMetadataPattern.FindAllStringSubmatch(query, -1) {
		if len(match) == 3 {
			metadata[match[1]] = match[2]
		}
	}
	var active *bool
	if match := searchActivePattern.FindStringSubmatch(query); len(match) == 2 {
		value := match[1] == "true"
		active = &value
	}
	out := make([]billing.Product, 0, len(products))
	for _, product := range products {
		if active != nil && product.Active != *active {
			continue
		}
		matched := true
		for key, value := range metadata {
			if product.Metadata[key] != value {
				matched = false
				break
			}
		}
		if matched {
			out = append(out, product)
		}
	}
	return out
}

func hasSubscriptionItemPatch(p params) bool {
	for key := range p.values {
		if strings.HasPrefix(key, "items[") {
			return true
		}
	}
	return false
}

func subscriptionItemsFromParams(p params, current billing.Subscription) []billing.LineItem {
	var out []billing.LineItem
	for i := 0; i < 100; i++ {
		priceID := p.first(fmt.Sprintf("items[%d][price]", i), fmt.Sprintf("items[%d][price_id]", i))
		quantity := p.int64(fmt.Sprintf("items[%d][quantity]", i))
		if i < len(current.Items) {
			if priceID == "" {
				priceID = current.Items[i].PriceID
			}
			if quantity == 0 {
				quantity = current.Items[i].Quantity
			}
		}
		if priceID == "" {
			continue
		}
		if quantity <= 0 {
			quantity = 1
		}
		out = append(out, billing.LineItem{PriceID: priceID, Quantity: quantity})
	}
	return out
}

func unix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}

func optionalUnix(t *time.Time) any {
	if t == nil || t.IsZero() {
		return nil
	}
	return t.Unix()
}

func optionalPaidAt(invoice billing.Invoice) any {
	if invoice.Status != "paid" {
		return nil
	}
	return unix(invoice.CreatedAt)
}

func paymentIntentError(intent billing.PaymentIntent) any {
	if intent.FailureCode == "" && intent.FailureMessage == "" {
		return nil
	}
	return map[string]any{
		"code":    intent.FailureCode,
		"message": intent.FailureMessage,
	}
}

func nonNilMap(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	return in
}

func emptyToNil(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func stringDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func sanitizeID(value string) string {
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "billtap"
	}
	return b.String()
}

func writeResult(w http.ResponseWriter, value any, err error) {
	if err != nil {
		switch {
		case errors.Is(err, billing.ErrNotFound):
			writeError(w, http.StatusNotFound, err)
		case errors.Is(err, billing.ErrInvalidInput), errors.Is(err, billing.ErrUnsupportedOutcome):
			writeError(w, http.StatusBadRequest, err)
		case errors.Is(err, fixtures.ErrInvalidFixture):
			writeError(w, http.StatusBadRequest, err)
		case errors.Is(err, webhooks.ErrNotFound):
			writeError(w, http.StatusNotFound, err)
		case errors.Is(err, webhooks.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func absoluteURL(r *http.Request, path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + path
}

func (h *Handler) emitCheckoutWebhooks(r *http.Request, result map[string]any) []webhooks.Event {
	if h.webhooks == nil {
		return nil
	}
	ctx := r.Context()
	sequence := time.Now().UTC().UnixNano()
	var emitted []webhooks.Event
	for idx, item := range checkoutWebhookPayloads(result) {
		event, _, err := h.webhooks.CreateEvent(ctx, webhooks.EventInput{
			Type:           item.eventType,
			ObjectPayload:  item.payload,
			RequestID:      "req_" + item.objectID,
			IdempotencyKey: "billtap:" + item.eventType + ":" + item.objectID,
			Source:         webhooks.SourceCheckout,
			Sequence:       sequence + int64(idx),
		})
		if err == nil {
			emitted = append(emitted, event)
		}
	}
	return emitted
}

func (h *Handler) deliveryAttemptResponses(r *http.Request, attempts []webhooks.DeliveryAttempt) []map[string]any {
	out := make([]map[string]any, 0, len(attempts))
	eventTypes := map[string]string{}
	for _, attempt := range attempts {
		eventType := eventTypes[attempt.EventID]
		if eventType == "" {
			if event, err := h.webhooks.GetEvent(r.Context(), attempt.EventID); err == nil {
				eventType = event.Type
				eventTypes[attempt.EventID] = eventType
			}
		}
		out = append(out, map[string]any{
			"id":              attempt.ID,
			"object":          attempt.Object,
			"event_id":        attempt.EventID,
			"event_type":      eventType,
			"endpoint_id":     attempt.EndpointID,
			"attempt_number":  attempt.AttemptNumber,
			"attempts":        attempt.AttemptNumber,
			"status":          attempt.Status,
			"scheduled_at":    attempt.ScheduledAt,
			"delivered_at":    attempt.DeliveredAt,
			"request_url":     security.RedactURL(attempt.RequestURL),
			"request_headers": security.RedactHeaders(attempt.RequestHeaders),
			"request_body":    security.RedactText(string(attempt.RequestBody)),
			"response_status": attempt.ResponseStatus,
			"response_body":   security.RedactText(attempt.ResponseBody),
			"error":           security.RedactText(attempt.Error),
			"next_retry_at":   attempt.NextRetryAt,
			"metadata":        attempt.Metadata,
		})
	}
	return out
}

func maskEndpoint(endpoint webhooks.Endpoint) webhooks.Endpoint {
	if endpoint.Secret != "" {
		endpoint.Secret = security.MaskSecret(endpoint.Secret)
	}
	return endpoint
}

func maskEndpoints(endpoints []webhooks.Endpoint) []webhooks.Endpoint {
	out := make([]webhooks.Endpoint, len(endpoints))
	for i, endpoint := range endpoints {
		out[i] = maskEndpoint(endpoint)
	}
	return out
}

func (h *Handler) collectObjects(r *http.Request) (map[string]any, error) {
	customers, err := h.billing.ListCustomers(r.Context())
	if err != nil {
		return nil, err
	}
	products, err := h.billing.ListProducts(r.Context())
	if err != nil {
		return nil, err
	}
	prices, err := h.billing.ListPrices(r.Context())
	if err != nil {
		return nil, err
	}
	checkoutSessions, err := h.billing.ListCheckoutSessions(r.Context())
	if err != nil {
		return nil, err
	}
	for i := range checkoutSessions {
		checkoutSessions[i].URL = absoluteURL(r, checkoutSessions[i].URL)
	}
	subscriptions, err := h.billing.ListSubscriptions(r.Context())
	if err != nil {
		return nil, err
	}
	invoices, err := h.billing.ListInvoices(r.Context())
	if err != nil {
		return nil, err
	}
	paymentIntents, err := h.billing.ListPaymentIntents(r.Context())
	if err != nil {
		return nil, err
	}
	result := map[string]any{
		"object":            "dashboard_objects",
		"customers":         customers,
		"products":          products,
		"prices":            prices,
		"checkout_sessions": checkoutSessions,
		"subscriptions":     subscriptions,
		"invoices":          invoices,
		"payment_intents":   paymentIntents,
	}
	if h.webhooks != nil {
		if endpoints, err := h.webhooks.ListEndpoints(r.Context(), webhooks.EndpointFilter{}); err == nil {
			result["webhook_endpoints"] = maskEndpoints(endpoints)
		}
		if events, err := h.webhooks.ListEvents(r.Context(), webhooks.EventFilter{}); err == nil {
			result["webhook_events"] = events
		}
	}
	return result, nil
}

type webhookPayload struct {
	eventType string
	objectID  string
	payload   json.RawMessage
}

func checkoutWebhookPayloads(result map[string]any) []webhookPayload {
	session, _ := result["session"].(billing.CheckoutSession)
	subscription, _ := result["subscription"].(billing.Subscription)
	invoice, _ := result["invoice"].(billing.Invoice)
	paymentIntent, _ := result["payment_intent"].(billing.PaymentIntent)

	var out []webhookPayload
	appendPayload := func(eventType string, objectID string, value any) {
		raw, err := json.Marshal(value)
		if err == nil {
			out = append(out, webhookPayload{eventType: eventType, objectID: objectID, payload: raw})
		}
	}
	appendPayload("checkout.session.completed", session.ID, session)
	if subscription.ID != "" {
		appendPayload("customer.subscription.created", subscription.ID, subscription)
	}
	if invoice.ID != "" {
		appendPayload("invoice.created", invoice.ID, invoice)
		appendPayload("invoice.finalized", invoice.ID, invoice)
	}
	if paymentIntent.ID != "" {
		appendPayload("payment_intent.created", paymentIntent.ID, paymentIntent)
		if paymentIntent.Status == "succeeded" {
			appendPayload("payment_intent.succeeded", paymentIntent.ID, paymentIntent)
		} else {
			appendPayload("payment_intent.payment_failed", paymentIntent.ID, paymentIntent)
		}
	}
	if invoice.ID != "" {
		if invoice.Status == "paid" {
			appendPayload("invoice.payment_succeeded", invoice.ID, invoice)
		} else {
			appendPayload("invoice.payment_failed", invoice.ID, invoice)
		}
	}
	if subscription.ID != "" {
		appendPayload("customer.subscription.updated", subscription.ID, subscription)
	}
	return out
}
