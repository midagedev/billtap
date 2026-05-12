package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/diagnostics"
	"github.com/hckim/billtap/internal/fixtures"
	"github.com/hckim/billtap/internal/scenarios"
	"github.com/hckim/billtap/internal/security"
	"github.com/hckim/billtap/internal/stripecompat"
	"github.com/hckim/billtap/internal/webhooks"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Billing       *billing.Service
	Webhooks      *webhooks.Service
	Diagnostics   *diagnostics.Service
	PublicBaseURL string
}

type Handler struct {
	billing     *billing.Service
	webhooks    *webhooks.Service
	diagnostics *diagnostics.Service
	publicBase  string
	mux         *http.ServeMux
	idem        *idempotencyStore
	local       *localEvidenceStore
	compat      stripecompat.Registry
	knownRoutes stripecompat.RouteCatalog
	validation  stripecompat.ValidationCatalog
}

func New(opts Options) http.Handler {
	h := &Handler{
		billing:     opts.Billing,
		webhooks:    opts.Webhooks,
		diagnostics: opts.Diagnostics,
		publicBase:  strings.TrimRight(opts.PublicBaseURL, "/"),
		mux:         http.NewServeMux(),
		idem:        newIdempotencyStore(),
		local:       newLocalEvidenceStore(),
		compat:      stripecompat.DefaultRegistry(),
		knownRoutes: stripecompat.DefaultRouteCatalog(),
		validation:  stripecompat.DefaultValidationCatalog(),
	}
	h.routes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.diagnostics != nil && isStripeAPIPath(r.URL.Path) {
		h.serveWithRequestTrace(w, r)
		return
	}
	h.serveWithIdempotency(w, r)
}

func (h *Handler) routes() {
	h.mux.HandleFunc("/v1/customers", h.handleCustomers)
	h.mux.HandleFunc("/v1/customers/", h.handleCustomer)
	h.mux.HandleFunc("/v1/products", h.handleProducts)
	h.mux.HandleFunc("/v1/products/search", h.handleProductSearch)
	h.mux.HandleFunc("/v1/products/", h.handleProduct)
	h.mux.HandleFunc("/v1/coupons", h.handleCoupons)
	h.mux.HandleFunc("/v1/coupons/", h.handleCoupon)
	h.mux.HandleFunc("/v1/promotion_codes", h.handlePromotionCodes)
	h.mux.HandleFunc("/v1/promotion_codes/", h.handlePromotionCode)
	h.mux.HandleFunc("/v1/prices/search", h.handlePriceSearch)
	h.mux.HandleFunc("/v1/prices", h.handlePrices)
	h.mux.HandleFunc("/v1/prices/", h.handlePrice)
	h.mux.HandleFunc("/v1/account", h.handlePlatformAccount)
	h.mux.HandleFunc("/v1/accounts", h.handleAccounts)
	h.mux.HandleFunc("/v1/accounts/", h.handleAccount)
	h.mux.HandleFunc("/v1/account_links", h.handleAccountLinks)
	h.mux.HandleFunc("/v1/account_sessions", h.handleAccountSessions)
	h.mux.HandleFunc("/v1/application_fees", h.handleApplicationFees)
	h.mux.HandleFunc("/v1/application_fees/", h.handleApplicationFee)
	h.mux.HandleFunc("/v1/transfers", h.handleTransfers)
	h.mux.HandleFunc("/v1/transfers/", h.handleTransfer)
	h.mux.HandleFunc("/v1/payouts", h.handlePayouts)
	h.mux.HandleFunc("/v1/payouts/", h.handlePayout)
	h.mux.HandleFunc("/v1/checkout/sessions", h.handleCheckoutSessions)
	h.mux.HandleFunc("/v1/checkout/sessions/", h.handleCheckoutSession)
	h.mux.HandleFunc("/v1/billing_portal/sessions", h.handleBillingPortalSessions)
	h.mux.HandleFunc("/v1/subscriptions", h.handleSubscriptions)
	h.mux.HandleFunc("/v1/subscriptions/", h.handleSubscription)
	h.mux.HandleFunc("/v1/subscription_schedules", h.handleSubscriptionSchedules)
	h.mux.HandleFunc("/v1/subscription_schedules/", h.handleSubscriptionSchedule)
	h.mux.HandleFunc("/v1/subscription_items", h.handleSubscriptionItems)
	h.mux.HandleFunc("/v1/subscription_items/", h.handleSubscriptionItem)
	h.mux.HandleFunc("/v1/invoices/create_preview", h.handleInvoicePreview)
	h.mux.HandleFunc("/v1/invoices/upcoming", h.handleInvoicePreview)
	h.mux.HandleFunc("/v1/invoices", h.handleInvoices)
	h.mux.HandleFunc("/v1/invoices/", h.handleInvoice)
	h.mux.HandleFunc("/v1/refunds", h.handleRefunds)
	h.mux.HandleFunc("/v1/refunds/", h.handleRefund)
	h.mux.HandleFunc("/v1/credit_notes", h.handleCreditNotes)
	h.mux.HandleFunc("/v1/credit_notes/", h.handleCreditNote)
	h.mux.HandleFunc("/v1/payment_intents", h.handlePaymentIntents)
	h.mux.HandleFunc("/v1/payment_intents/", h.handlePaymentIntent)
	h.mux.HandleFunc("/v1/setup_intents", h.handleSetupIntents)
	h.mux.HandleFunc("/v1/setup_intents/", h.handleSetupIntent)
	h.mux.HandleFunc("/v1/disputes", h.handleDisputes)
	h.mux.HandleFunc("/v1/disputes/", h.handleDispute)
	h.mux.HandleFunc("/v1/charges/", h.handleChargeSubresource)
	h.mux.HandleFunc("/v1/test_helpers/customers/", h.handleTestHelperCustomer)
	h.mux.HandleFunc("/v1/test_helpers/test_clocks", h.handleTestClocks)
	h.mux.HandleFunc("/v1/test_helpers/test_clocks/", h.handleTestClock)
	h.mux.HandleFunc("/v1/payment_methods", h.handlePaymentMethods)
	h.mux.HandleFunc("/v1/webhook_endpoints", h.handleWebhookEndpoints)
	h.mux.HandleFunc("/v1/webhook_endpoints/", h.handleWebhookEndpoint)
	h.mux.HandleFunc("/v1/events", h.handleEvents)
	h.mux.HandleFunc("/v1/events/", h.handleEvent)
	h.mux.HandleFunc("/v1", h.handleStripeFallback)
	h.mux.HandleFunc("/v1/", h.handleStripeFallback)
	h.mux.HandleFunc("/v2", h.handleStripeFallback)
	h.mux.HandleFunc("/v2/", h.handleStripeFallback)
	h.mux.HandleFunc("/api/checkout/sessions/", h.handleCheckoutCompletion)
	h.mux.HandleFunc("/api/objects", h.handleObjects)
	h.mux.HandleFunc("/api/debug-bundles", h.handleDebugBundles)
	h.mux.HandleFunc("/api/diagnostics", h.handleDiagnostics)
	h.mux.HandleFunc("/api/request-traces", h.handleRequestTraces)
	h.mux.HandleFunc("/api/fixtures/apply", h.handleFixtureApply)
	h.mux.HandleFunc("/api/fixtures/validate", h.handleFixtureValidate)
	h.mux.HandleFunc("/api/fixtures/snapshot", h.handleFixtureSnapshot)
	h.mux.HandleFunc("/api/fixtures/resolve", h.handleFixtureResolve)
	h.mux.HandleFunc("/api/fixtures/assert", h.handleFixtureAssert)
	h.mux.HandleFunc("/api/portal", h.handlePortal)
	h.mux.HandleFunc("/api/portal/customers/", h.handlePortalCustomer)
	h.mux.HandleFunc("/api/portal/subscriptions/", h.handlePortalSubscription)
	h.mux.HandleFunc("/api/scenarios/run", h.handleScenarioRun)
	h.mux.HandleFunc("/api/timeline", h.handleTimeline)
	h.mux.HandleFunc("/api/delivery-attempts", h.handleDeliveryAttempts)
	h.mux.HandleFunc("/api/audit-log", h.handleAuditLog)
	h.mux.HandleFunc("/api/retention/apply", h.handleRetentionApply)
	h.mux.HandleFunc("/api/payment_intents/", h.handlePaymentIntentAction)
	h.mux.HandleFunc("/api/events/replay-group", h.handleEventReplayGroup)
	h.mux.HandleFunc("/api/webhooks/endpoints/", h.handleWebhookEndpointAction)
	h.mux.HandleFunc("/api/events/", h.handleEventAction)
	h.mux.HandleFunc("/api/disputes", h.handleDisputeSimulation)
}

func isStripeAPIPath(path string) bool {
	return strings.HasPrefix(path, "/v1/") || strings.HasPrefix(path, "/v2/")
}

func (h *Handler) compatibilityClaim(r *http.Request) (stripecompat.Claim, bool) {
	return h.compat.Lookup(r.Method, r.URL.Path)
}

func (h *Handler) handleStripeFallback(w http.ResponseWriter, r *http.Request) {
	if h.writeKnownUnsupportedRoute(w, r) {
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	if h.writeKnownUnsupportedRoute(w, r) {
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) methodNotAllowed(w http.ResponseWriter, r *http.Request, allow string) {
	if h.writeKnownUnsupportedRoute(w, r) {
		return
	}
	methodNotAllowed(w, allow)
}

func (h *Handler) writeKnownUnsupportedRoute(w http.ResponseWriter, r *http.Request) bool {
	route, ok := h.knownRoutes.Lookup(r.Method, r.URL.Path)
	if !ok {
		return false
	}
	if claim, ok := h.compatibilityClaim(r); ok {
		claimPath := stripecompat.NormalizePath(claim.Path)
		routePath := stripecompat.NormalizePath(route.Path)
		requestPath := stripecompat.NormalizePath(r.URL.Path)
		if claimPath == routePath || (!strings.Contains(claim.Path, "{") && claimPath == requestPath) {
			return false
		}
	}
	if validationErr := h.validation.Validate(r); validationErr != nil {
		writeStripeError(w, http.StatusBadRequest, stripeAPIError{
			Type:    stripeErrorInvalidReq,
			Message: validationErr.Message,
			Param:   validationErr.Param,
			Code:    validationErr.Code,
		})
		return true
	}
	message := fmt.Sprintf("Billtap knows this Stripe API route from %s but does not implement it yet: %s %s.", route.Source, route.Method, route.Path)
	writeStripeError(w, http.StatusBadRequest, stripeAPIError{
		Type:    stripeErrorInvalidReq,
		Message: message,
		Code:    "unsupported_endpoint",
	})
	return true
}

func (h *Handler) handleCustomers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateCustomerCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		metadata := p.metadata()
		if p.string("test_clock") != "" {
			if metadata == nil {
				metadata = map[string]string{}
			}
			metadata["test_clock"] = p.string("test_clock")
		}
		if discounts, err := h.discountsFromParams(p); err != nil {
			writeResult(w, nil, err)
			return
		} else if len(discounts) > 0 {
			metadata = billing.MergeDiscountMetadata(metadata, discounts)
		}
		customer, err := h.billing.CreateCustomer(r.Context(), billing.Customer{
			ID:       p.string("id"),
			Email:    p.string("email"),
			Name:     p.string("name"),
			Metadata: metadata,
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
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleCustomer(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/v1/customers/")
	id, subresource, hasSubresource := strings.Cut(rest, "/")
	if id == "" {
		h.notFound(w, r)
		return
	}
	if hasSubresource {
		h.handleCustomerSubresource(w, r, id, subresource)
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
		if err := validateCustomerUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		metadata := p.metadata()
		if metadata != nil || p.has("test_clock") || hasDiscountParams(p) {
			current, err := h.billing.GetCustomer(r.Context(), id)
			if err != nil {
				writeResult(w, nil, err)
				return
			}
			merged := copyStringMap(current.Metadata)
			for key, value := range metadata {
				merged[key] = value
			}
			metadata = merged
		}
		if p.string("test_clock") != "" {
			if metadata == nil {
				metadata = map[string]string{}
			}
			metadata["test_clock"] = p.string("test_clock")
		}
		if discounts, err := h.discountsFromParams(p); err != nil {
			writeResult(w, nil, err)
			return
		} else if len(discounts) > 0 {
			metadata = billing.MergeDiscountMetadata(metadata, discounts)
		}
		customer, err := h.billing.UpdateCustomer(r.Context(), id, billing.Customer{
			Email:    p.string("email"),
			Name:     p.string("name"),
			Metadata: metadata,
		})
		writeResult(w, stripeCustomer(customer), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleCustomerSubresource(w http.ResponseWriter, r *http.Request, customerID string, subresource string) {
	switch subresource {
	case "cash_balance":
		h.handleCustomerCashBalance(w, r, customerID)
	case "cash_balance_transactions":
		h.handleCustomerCashBalanceTransactions(w, r, customerID, "")
	case "payment_methods":
		h.handleCustomerPaymentMethods(w, r, customerID)
	case "discount":
		h.handleCustomerDiscount(w, r, customerID)
	default:
		if strings.HasPrefix(subresource, "cash_balance_transactions/") {
			h.handleCustomerCashBalanceTransactions(w, r, customerID, strings.TrimPrefix(subresource, "cash_balance_transactions/"))
			return
		}
		if strings.HasPrefix(subresource, "subscriptions/") {
			rest := strings.TrimPrefix(subresource, "subscriptions/")
			subscriptionID, nested, ok := strings.Cut(rest, "/")
			if ok && nested == "discount" && subscriptionID != "" {
				h.handleSubscriptionDiscount(w, r, subscriptionID)
				return
			}
		}
		h.notFound(w, r)
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
		if err := validateProductCreate(p); err != nil {
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
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleProductSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
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
		h.notFound(w, r)
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
		if err := validateProductUpdate(p); err != nil {
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
		h.methodNotAllowed(w, r, "GET, POST")
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
		if err := validatePriceCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateProductExists(h.billing.GetProduct(r.Context(), p.first("product", "product_id"))); err != nil {
			writeResult(w, nil, err)
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
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handlePriceSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	p := paramsFromValues(r.URL.Query())
	if err := validatePriceSearch(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	criteria, err := parsePriceSearchQuery(p.string("query"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	prices, err := h.billing.ListPrices(r.Context())
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	filtered := filterPriceSearchResults(prices, criteria)
	if limit := queryInt(r, "limit"); limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	writeResult(w, stripeSearchResult(r.URL.Path, p.string("query"), stripePrices(filtered)), nil)
}

func (h *Handler) handlePrice(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/prices/")
	if id == "" || strings.Contains(id, "/") {
		h.notFound(w, r)
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
		if err := validatePriceUpdate(p); err != nil {
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
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleAccounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateAccountCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		account, err := h.billing.CreateAccount(r.Context(), billing.Account{
			ID:               p.string("id"),
			Type:             p.stringDefault("type", "express"),
			Country:          strings.ToUpper(p.stringDefault("country", "US")),
			Email:            p.string("email"),
			BusinessType:     p.string("business_type"),
			DefaultCurrency:  p.stringDefault("default_currency", "usd"),
			ChargesEnabled:   true,
			PayoutsEnabled:   true,
			DetailsSubmitted: p.boolDefault("details_submitted", true),
			Capabilities:     accountCapabilities(p),
			Metadata:         p.metadata(),
		})
		writeResult(w, stripeAccount(account), err)
	case http.MethodGet:
		accounts, err := h.billing.ListAccounts(r.Context())
		writeResult(w, stripeList(r.URL.Path, stripeAccounts(filterAccounts(accounts, r))), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handlePlatformAccount(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/account" {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	writeJSON(w, http.StatusOK, stripeAccount(billing.Account{
		ID:               "acct_billtap_platform",
		Object:           billing.ObjectAccount,
		Type:             "standard",
		Country:          "US",
		Email:            "platform@example.test",
		BusinessType:     "company",
		DefaultCurrency:  "usd",
		ChargesEnabled:   true,
		PayoutsEnabled:   true,
		DetailsSubmitted: true,
		Capabilities:     map[string]string{"card_payments": "active", "transfers": "active"},
		Metadata:         map[string]string{"billtap_account": "platform"},
		CreatedAt:        time.Unix(0, 0).UTC(),
		UpdatedAt:        time.Unix(0, 0).UTC(),
	}))
}

func (h *Handler) handleAccount(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/accounts/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	parts := strings.Split(rest, "/")
	id := parts[0]
	if len(parts) > 1 {
		h.handleAccountSubresource(w, r, id, parts[1:])
		return
	}
	switch r.Method {
	case http.MethodGet:
		account, err := h.billing.GetAccount(r.Context(), id)
		writeResult(w, stripeAccount(account), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateAccountUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		account, err := h.billing.UpdateAccount(r.Context(), id, billing.Account{
			Email:           p.string("email"),
			BusinessType:    p.string("business_type"),
			DefaultCurrency: p.string("default_currency"),
			Capabilities:    accountCapabilities(p),
			Metadata:        p.metadata(),
		})
		writeResult(w, stripeAccount(account), err)
	case http.MethodDelete:
		_, err := h.billing.GetAccount(r.Context(), id)
		writeResult(w, stripeDeleted(id, billing.ObjectAccount), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST, DELETE")
	}
}

func (h *Handler) handleAccountSubresource(w http.ResponseWriter, r *http.Request, accountID string, parts []string) {
	if accountID == "" || len(parts) == 0 {
		h.notFound(w, r)
		return
	}
	switch parts[0] {
	case "capabilities":
		h.handleAccountCapabilities(w, r, accountID, parts[1:])
	case "external_accounts", "bank_accounts":
		h.handleAccountExternalAccounts(w, r, accountID, parts[0], parts[1:])
	case "people", "persons":
		h.handleAccountPeople(w, r, accountID, parts[0], parts[1:])
	case "login_links":
		h.handleAccountLoginLinks(w, r, accountID, parts[1:])
	case "reject":
		h.handleAccountReject(w, r, accountID, parts[1:])
	default:
		h.notFound(w, r)
	}
}

func (h *Handler) handleAccountCapabilities(w http.ResponseWriter, r *http.Request, accountID string, parts []string) {
	account, err := h.billing.GetAccount(r.Context(), accountID)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	if len(parts) == 0 {
		if r.Method != http.MethodGet {
			h.methodNotAllowed(w, r, "GET")
			return
		}
		data := make([]map[string]any, 0, len(account.Capabilities))
		for capability, status := range account.Capabilities {
			data = append(data, stripeCapability(account.ID, capability, status))
		}
		writeJSON(w, http.StatusOK, stripeList(r.URL.Path, data))
		return
	}
	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}
	capability := parts[0]
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, stripeCapability(account.ID, capability, stringDefault(account.Capabilities[capability], "inactive")))
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateAccountCapabilityUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		status := p.string("status")
		if status == "" {
			if p.boolDefault("requested", true) {
				status = "active"
			} else {
				status = "inactive"
			}
		}
		updated, err := h.billing.UpdateAccount(r.Context(), account.ID, billing.Account{Capabilities: map[string]string{capability: status}})
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		writeJSON(w, http.StatusOK, stripeCapability(updated.ID, capability, updated.Capabilities[capability]))
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleAccountExternalAccounts(w http.ResponseWriter, r *http.Request, accountID string, collection string, parts []string) {
	if _, err := h.billing.GetAccount(r.Context(), accountID); err != nil {
		writeResult(w, nil, err)
		return
	}
	if len(parts) == 0 {
		switch r.Method {
		case http.MethodGet:
			resources, err := h.billing.ListConnectResources(r.Context(), billing.ConnectResourceFilter{Object: billing.ObjectBankAccount, AccountID: accountID})
			writeResult(w, stripeList(r.URL.Path, stripeConnectResources(resources)), err)
		case http.MethodPost:
			p, err := parseParams(r)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			if err := validateExternalAccountCreate(p); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			resource, err := h.billing.CreateConnectResource(r.Context(), billing.ConnectResource{
				ID:            p.string("id"),
				Object:        billing.ObjectBankAccount,
				AccountID:     accountID,
				Country:       p.stringDefault("country", "US"),
				Currency:      p.stringDefault("currency", "usd"),
				BankName:      p.stringDefault("bank_name", "Billtap Bank"),
				Last4:         last4(p.first("account_number", "external_account", "token")),
				RoutingNumber: p.string("routing_number"),
				Status:        "new",
				Metadata:      p.metadata(),
				Data: map[string]string{
					"default_for_currency": strconv.FormatBool(p.boolDefault("default_for_currency", false)),
					"account_holder_name":  p.string("account_holder_name"),
					"account_holder_type":  p.string("account_holder_type"),
				},
			})
			writeResult(w, stripeExternalAccount(resource), err)
		default:
			h.methodNotAllowed(w, r, "GET, POST")
		}
		return
	}
	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}
	id := parts[0]
	switch r.Method {
	case http.MethodGet:
		resource, err := h.billing.GetConnectResource(r.Context(), billing.ObjectBankAccount, id)
		if err == nil && resource.AccountID != accountID {
			err = billing.ErrNotFound
		}
		writeResult(w, stripeExternalAccount(resource), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateExternalAccountUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		current, err := h.billing.GetConnectResource(r.Context(), billing.ObjectBankAccount, id)
		if err == nil && current.AccountID != accountID {
			err = billing.ErrNotFound
		}
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		data := map[string]string{}
		if p.has("default_for_currency") {
			data["default_for_currency"] = strconv.FormatBool(p.boolDefault("default_for_currency", false))
		}
		if p.has("account_holder_name") {
			data["account_holder_name"] = p.string("account_holder_name")
		}
		if p.has("account_holder_type") {
			data["account_holder_type"] = p.string("account_holder_type")
		}
		resource, err := h.billing.UpdateConnectResource(r.Context(), billing.ObjectBankAccount, id, billing.ConnectResource{
			BankName: p.string("bank_name"),
			Metadata: p.metadata(),
			Data:     data,
		})
		writeResult(w, stripeExternalAccount(resource), err)
	case http.MethodDelete:
		current, err := h.billing.GetConnectResource(r.Context(), billing.ObjectBankAccount, id)
		if err == nil && current.AccountID != accountID {
			err = billing.ErrNotFound
		}
		if err == nil {
			_, err = h.billing.DeleteConnectResource(r.Context(), billing.ObjectBankAccount, id)
		}
		writeResult(w, stripeDeleted(id, billing.ObjectBankAccount), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST, DELETE")
	}
	_ = collection
}

func (h *Handler) handleAccountPeople(w http.ResponseWriter, r *http.Request, accountID string, collection string, parts []string) {
	if _, err := h.billing.GetAccount(r.Context(), accountID); err != nil {
		writeResult(w, nil, err)
		return
	}
	if len(parts) == 0 {
		switch r.Method {
		case http.MethodGet:
			resources, err := h.billing.ListConnectResources(r.Context(), billing.ConnectResourceFilter{Object: billing.ObjectPerson, AccountID: accountID})
			writeResult(w, stripeList(r.URL.Path, stripeConnectResources(resources)), err)
		case http.MethodPost:
			p, err := parseParams(r)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			if err := validateAccountPersonCreate(p); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			resource, err := h.billing.CreateConnectResource(r.Context(), billing.ConnectResource{
				ID:        p.string("id"),
				Object:    billing.ObjectPerson,
				AccountID: accountID,
				Metadata:  p.metadata(),
				Data:      personDataFromParams(p),
			})
			writeResult(w, stripePerson(resource), err)
		default:
			h.methodNotAllowed(w, r, "GET, POST")
		}
		return
	}
	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}
	personID := parts[0]
	switch r.Method {
	case http.MethodGet:
		resource, err := h.billing.GetConnectResource(r.Context(), billing.ObjectPerson, personID)
		if err == nil && resource.AccountID != accountID {
			err = billing.ErrNotFound
		}
		writeResult(w, stripePerson(resource), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateAccountPersonUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		current, err := h.billing.GetConnectResource(r.Context(), billing.ObjectPerson, personID)
		if err == nil && current.AccountID != accountID {
			err = billing.ErrNotFound
		}
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		resource, err := h.billing.UpdateConnectResource(r.Context(), billing.ObjectPerson, personID, billing.ConnectResource{
			Metadata: p.metadata(),
			Data:     personDataFromParams(p),
		})
		writeResult(w, stripePerson(resource), err)
	case http.MethodDelete:
		current, err := h.billing.GetConnectResource(r.Context(), billing.ObjectPerson, personID)
		if err == nil && current.AccountID != accountID {
			err = billing.ErrNotFound
		}
		if err == nil {
			_, err = h.billing.DeleteConnectResource(r.Context(), billing.ObjectPerson, personID)
		}
		writeResult(w, stripeDeleted(personID, billing.ObjectPerson), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST, DELETE")
	}
	_ = collection
}

func (h *Handler) handleAccountLoginLinks(w http.ResponseWriter, r *http.Request, accountID string, parts []string) {
	if len(parts) != 0 {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	if _, err := h.billing.GetAccount(r.Context(), accountID); err != nil {
		writeResult(w, nil, err)
		return
	}
	now := time.Now().UTC()
	writeJSON(w, http.StatusOK, map[string]any{
		"id":       "ll_" + sanitizeID(accountID+"_"+now.Format(time.RFC3339Nano)),
		"object":   "login_link",
		"created":  now.Unix(),
		"url":      h.absoluteURL(r, "/connect/accounts/"+accountID+"/login"),
		"livemode": false,
	})
}

func (h *Handler) handleAccountReject(w http.ResponseWriter, r *http.Request, accountID string, parts []string) {
	if len(parts) != 0 {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateAccountReject(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	account, err := h.billing.UpdateAccount(r.Context(), accountID, billing.Account{
		Metadata: map[string]string{"rejected_reason": p.stringDefault("reason", "other")},
	})
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	account.ChargesEnabled = false
	account.PayoutsEnabled = false
	writeJSON(w, http.StatusOK, stripeAccount(account))
}

func (h *Handler) handleAccountLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateAccountLinkCreate(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	accountID := p.string("account")
	if _, err := h.billing.GetAccount(r.Context(), accountID); err != nil {
		writeResult(w, nil, err)
		return
	}
	linkType := p.stringDefault("type", "account_onboarding")
	link := map[string]any{
		"id":          "link_" + sanitizeID(accountID+"_"+linkType+"_"+time.Now().UTC().Format(time.RFC3339Nano)),
		"object":      billing.ObjectAccountLink,
		"created":     time.Now().UTC().Unix(),
		"expires_at":  time.Now().UTC().Add(30 * time.Minute).Unix(),
		"url":         h.absoluteURL(r, "/connect/accounts/"+accountID+"/"+linkType),
		"account":     accountID,
		"type":        linkType,
		"livemode":    false,
		"refresh_url": emptyToNil(p.string("refresh_url")),
		"return_url":  emptyToNil(p.string("return_url")),
	}
	writeJSON(w, http.StatusOK, link)
}

func (h *Handler) handleAccountSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateAccountSessionCreate(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	accountID := p.string("account")
	if _, err := h.billing.GetAccount(r.Context(), accountID); err != nil {
		writeResult(w, nil, err)
		return
	}
	sessionID := "accsess_" + sanitizeID(accountID+"_"+time.Now().UTC().Format(time.RFC3339Nano))
	expiresAt := time.Now().UTC().Add(30 * time.Minute).Unix()
	writeJSON(w, http.StatusOK, map[string]any{
		"id":            sessionID,
		"object":        billing.ObjectAccountSession,
		"account":       accountID,
		"client_secret": sessionID + "_secret_billtap",
		"components":    accountSessionComponents(p),
		"created":       time.Now().UTC().Unix(),
		"expires_at":    expiresAt,
		"livemode":      false,
	})
}

func (h *Handler) handleTransfers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resources, err := h.billing.ListConnectResources(r.Context(), billing.ConnectResourceFilter{
			Object:      billing.ObjectTransfer,
			Destination: r.URL.Query().Get("destination"),
		})
		writeResult(w, stripeList(r.URL.Path, stripeConnectResources(resources)), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateTransferCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		resource, err := h.billing.CreateConnectResource(r.Context(), billing.ConnectResource{
			ID:                p.string("id"),
			Object:            billing.ObjectTransfer,
			AccountID:         r.Header.Get("Stripe-Account"),
			Amount:            p.int64("amount"),
			Currency:          p.string("currency"),
			Destination:       p.string("destination"),
			SourceTransaction: p.string("source_transaction"),
			Description:       p.string("description"),
			Metadata:          p.metadata(),
			Data: map[string]string{
				"transfer_group": p.string("transfer_group"),
			},
		})
		if err == nil {
			h.emitGenericWebhook(r, "transfer.created", resource.ID, stripeTransfer(resource), webhooks.SourceAPI)
		}
		writeResult(w, stripeTransfer(resource), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleTransfer(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/transfers/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	parts := strings.Split(rest, "/")
	transferID := parts[0]
	if len(parts) > 1 {
		if parts[1] != "reversals" {
			h.notFound(w, r)
			return
		}
		h.handleTransferReversals(w, r, transferID, parts[2:])
		return
	}
	switch r.Method {
	case http.MethodGet:
		resource, err := h.billing.GetConnectResource(r.Context(), billing.ObjectTransfer, transferID)
		writeResult(w, stripeTransfer(resource), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateTransferUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		resource, err := h.billing.UpdateConnectResource(r.Context(), billing.ObjectTransfer, transferID, billing.ConnectResource{
			Description: p.string("description"),
			Metadata:    p.metadata(),
		})
		writeResult(w, stripeTransfer(resource), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleTransferReversals(w http.ResponseWriter, r *http.Request, transferID string, parts []string) {
	transfer, err := h.billing.GetConnectResource(r.Context(), billing.ObjectTransfer, transferID)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	if len(parts) == 0 {
		switch r.Method {
		case http.MethodGet:
			resources, err := h.billing.ListConnectResources(r.Context(), billing.ConnectResourceFilter{Object: billing.ObjectTransferReversal, ParentID: transferID})
			writeResult(w, stripeList(r.URL.Path, stripeConnectResources(resources)), err)
		case http.MethodPost:
			p, err := parseParams(r)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			if err := validateTransferReversalCreate(p); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			amount := p.int64Default("amount", transfer.Amount)
			resource, err := h.billing.CreateConnectResource(r.Context(), billing.ConnectResource{
				ID:        p.string("id"),
				Object:    billing.ObjectTransferReversal,
				ParentID:  transferID,
				AccountID: transfer.AccountID,
				Amount:    amount,
				Currency:  transfer.Currency,
				Metadata:  p.metadata(),
				Data: map[string]string{
					"refund_application_fee": strconv.FormatBool(p.boolDefault("refund_application_fee", false)),
				},
			})
			if err == nil {
				_, _ = h.billing.UpdateConnectResource(r.Context(), billing.ObjectTransfer, transferID, billing.ConnectResource{
					Status: "reversed",
					Data:   map[string]string{"amount_reversed": strconv.FormatInt(amount, 10)},
				})
				h.emitGenericWebhook(r, "transfer.reversed", resource.ID, stripeTransferReversal(resource), webhooks.SourceAPI)
			}
			writeResult(w, stripeTransferReversal(resource), err)
		default:
			h.methodNotAllowed(w, r, "GET, POST")
		}
		return
	}
	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}
	reversalID := parts[0]
	switch r.Method {
	case http.MethodGet:
		resource, err := h.billing.GetConnectResource(r.Context(), billing.ObjectTransferReversal, reversalID)
		if err == nil && resource.ParentID != transferID {
			err = billing.ErrNotFound
		}
		writeResult(w, stripeTransferReversal(resource), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateTransferReversalUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		current, err := h.billing.GetConnectResource(r.Context(), billing.ObjectTransferReversal, reversalID)
		if err == nil && current.ParentID != transferID {
			err = billing.ErrNotFound
		}
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		resource, err := h.billing.UpdateConnectResource(r.Context(), billing.ObjectTransferReversal, reversalID, billing.ConnectResource{Metadata: p.metadata()})
		writeResult(w, stripeTransferReversal(resource), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handlePayouts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resources, err := h.billing.ListConnectResources(r.Context(), billing.ConnectResourceFilter{
			Object: billing.ObjectPayout,
			Status: r.URL.Query().Get("status"),
		})
		writeResult(w, stripeList(r.URL.Path, stripeConnectResources(resources)), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validatePayoutCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		arrival := time.Now().UTC().Add(24 * time.Hour)
		resource, err := h.billing.CreateConnectResource(r.Context(), billing.ConnectResource{
			ID:          p.string("id"),
			Object:      billing.ObjectPayout,
			AccountID:   r.Header.Get("Stripe-Account"),
			Amount:      p.int64("amount"),
			Currency:    p.string("currency"),
			Destination: p.string("destination"),
			Description: p.string("description"),
			ArrivalDate: arrival,
			Metadata:    p.metadata(),
			Data: map[string]string{
				"method":               p.stringDefault("method", "standard"),
				"statement_descriptor": p.string("statement_descriptor"),
			},
		})
		if err == nil {
			h.emitGenericWebhook(r, "payout.created", resource.ID, stripePayout(resource), webhooks.SourceAPI)
		}
		writeResult(w, stripePayout(resource), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handlePayout(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/payouts/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	parts := strings.Split(rest, "/")
	id := parts[0]
	if len(parts) == 2 {
		h.handlePayoutAction(w, r, id, parts[1])
		return
	}
	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		resource, err := h.billing.GetConnectResource(r.Context(), billing.ObjectPayout, id)
		writeResult(w, stripePayout(resource), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validatePayoutUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		resource, err := h.billing.UpdateConnectResource(r.Context(), billing.ObjectPayout, id, billing.ConnectResource{
			Description: p.string("description"),
			Metadata:    p.metadata(),
		})
		writeResult(w, stripePayout(resource), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handlePayoutAction(w http.ResponseWriter, r *http.Request, id string, action string) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	status := ""
	eventType := ""
	switch action {
	case "cancel":
		status = "canceled"
		eventType = "payout.canceled"
	case "reverse":
		status = "reversed"
		eventType = "payout.reversed"
	default:
		h.notFound(w, r)
		return
	}
	resource, err := h.billing.UpdateConnectResource(r.Context(), billing.ObjectPayout, id, billing.ConnectResource{Status: status})
	if err == nil {
		h.emitGenericWebhook(r, eventType, resource.ID, stripePayout(resource), webhooks.SourceAPI)
	}
	writeResult(w, stripePayout(resource), err)
}

func (h *Handler) handleApplicationFees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	resources, err := h.billing.ListConnectResources(r.Context(), billing.ConnectResourceFilter{Object: billing.ObjectApplicationFee})
	writeResult(w, stripeList(r.URL.Path, stripeConnectResources(resources)), err)
}

func (h *Handler) handleApplicationFee(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/application_fees/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	parts := strings.Split(rest, "/")
	feeID := parts[0]
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			h.methodNotAllowed(w, r, "GET")
			return
		}
		fee, err := h.billing.GetConnectResource(r.Context(), billing.ObjectApplicationFee, feeID)
		writeResult(w, stripeApplicationFee(fee, nil), err)
		return
	}
	if parts[1] == "refund" && len(parts) == 2 {
		h.handleApplicationFeeRefundCreate(w, r, feeID)
		return
	}
	if parts[1] == "refunds" {
		h.handleApplicationFeeRefunds(w, r, feeID, parts[2:])
		return
	}
	h.notFound(w, r)
}

func (h *Handler) handleApplicationFeeRefunds(w http.ResponseWriter, r *http.Request, feeID string, parts []string) {
	if len(parts) == 0 {
		switch r.Method {
		case http.MethodGet:
			resources, err := h.billing.ListConnectResources(r.Context(), billing.ConnectResourceFilter{Object: billing.ObjectFeeRefund, ParentID: feeID})
			writeResult(w, stripeList(r.URL.Path, stripeConnectResources(resources)), err)
		case http.MethodPost:
			h.handleApplicationFeeRefundCreate(w, r, feeID)
		default:
			h.methodNotAllowed(w, r, "GET, POST")
		}
		return
	}
	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}
	refundID := parts[0]
	switch r.Method {
	case http.MethodGet:
		resource, err := h.billing.GetConnectResource(r.Context(), billing.ObjectFeeRefund, refundID)
		if err == nil && resource.ParentID != feeID {
			err = billing.ErrNotFound
		}
		writeResult(w, stripeApplicationFeeRefund(resource), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateApplicationFeeRefundUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		current, err := h.billing.GetConnectResource(r.Context(), billing.ObjectFeeRefund, refundID)
		if err == nil && current.ParentID != feeID {
			err = billing.ErrNotFound
		}
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		resource, err := h.billing.UpdateConnectResource(r.Context(), billing.ObjectFeeRefund, refundID, billing.ConnectResource{Metadata: p.metadata()})
		writeResult(w, stripeApplicationFeeRefund(resource), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleApplicationFeeRefundCreate(w http.ResponseWriter, r *http.Request, feeID string) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateApplicationFeeRefundCreate(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	fee, err := h.ensureApplicationFee(r, feeID, p)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	amount := p.int64Default("amount", fee.Amount)
	resource, err := h.billing.CreateConnectResource(r.Context(), billing.ConnectResource{
		ID:        p.string("id"),
		Object:    billing.ObjectFeeRefund,
		AccountID: fee.AccountID,
		ParentID:  fee.ID,
		Amount:    amount,
		Currency:  fee.Currency,
		Metadata:  p.metadata(),
	})
	if err == nil {
		refunded := amount
		if current := fee.Data["amount_refunded"]; current != "" {
			refunded += parseInt64(current)
		}
		_, _ = h.billing.UpdateConnectResource(r.Context(), billing.ObjectApplicationFee, fee.ID, billing.ConnectResource{
			Data: map[string]string{"amount_refunded": strconv.FormatInt(refunded, 10)},
		})
		h.emitGenericWebhook(r, "application_fee.refunded", resource.ID, stripeApplicationFeeRefund(resource), webhooks.SourceAPI)
	}
	writeResult(w, stripeApplicationFeeRefund(resource), err)
}

func (h *Handler) ensureApplicationFee(r *http.Request, feeID string, p params) (billing.ConnectResource, error) {
	fee, err := h.billing.GetConnectResource(r.Context(), billing.ObjectApplicationFee, feeID)
	if err == nil {
		return fee, nil
	}
	if !errors.Is(err, billing.ErrNotFound) {
		return billing.ConnectResource{}, err
	}
	amount := p.int64("amount")
	if amount == 0 {
		amount = 1
	}
	return h.billing.CreateConnectResource(r.Context(), billing.ConnectResource{
		ID:        feeID,
		Object:    billing.ObjectApplicationFee,
		AccountID: r.Header.Get("Stripe-Account"),
		Amount:    amount,
		Currency:  "usd",
		Data: map[string]string{
			"charge": p.string("charge"),
		},
	})
}

func (h *Handler) handleCheckoutSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateCheckoutSessionCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		customer, err := h.billing.GetCustomer(r.Context(), p.first("customer", "customer_id"))
		if err := validateCustomerExists(customer, err); err != nil {
			writeResult(w, nil, err)
			return
		}
		for _, item := range p.lineItems() {
			if err := validatePriceExists(h.billing.GetPrice(r.Context(), item.PriceID)); err != nil {
				writeResult(w, nil, err)
				return
			}
		}
		discounts, err := h.discountsFromParamsOrCustomer(r, p, customer)
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		session, err := h.billing.CreateCheckoutSession(r.Context(), billing.CheckoutSession{
			CustomerID:          p.first("customer", "customer_id"),
			Mode:                p.stringDefault("mode", "subscription"),
			LineItems:           p.lineItems(),
			Discounts:           discounts,
			SuccessURL:          p.string("success_url"),
			CancelURL:           p.string("cancel_url"),
			AllowPromotionCodes: p.boolDefault("allow_promotion_codes", false),
			TrialPeriodDays:     p.int64("subscription_data[trial_period_days]"),
		})
		if err == nil {
			session.URL = h.absoluteURL(r, session.URL)
		}
		writeResult(w, stripeCheckoutSession(session), err)
	case http.MethodGet:
		sessions, err := h.billing.ListCheckoutSessions(r.Context())
		data := make([]map[string]any, 0, len(sessions))
		for i := range sessions {
			sessions[i].URL = h.absoluteURL(r, sessions[i].URL)
			data = append(data, stripeCheckoutSession(sessions[i]))
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
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
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	session, err := h.billing.GetCheckoutSession(r.Context(), rest)
	if err == nil {
		session.URL = h.absoluteURL(r, session.URL)
	}
	writeResult(w, stripeCheckoutSession(session), err)
}

func (h *Handler) handleBillingPortalSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateBillingPortalSessionCreate(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	customerID := p.first("customer", "customer_id")
	if _, err := h.billing.GetCustomer(r.Context(), customerID); err != nil {
		writeResult(w, nil, err)
		return
	}
	sessionID := "bps_" + sanitizeID(customerID+"_"+time.Now().UTC().Format(time.RFC3339Nano))
	returnURL := billingPortalSessionReturnURL(p)
	flowType := billingPortalSessionFlowType(p)
	session := map[string]any{
		"id":            sessionID,
		"object":        "billing_portal.session",
		"customer":      customerID,
		"configuration": emptyToNil(p.string("configuration")),
		"flow":          nil,
		"locale":        emptyToNil(p.string("locale")),
		"on_behalf_of":  emptyToNil(p.string("on_behalf_of")),
		"return_url":    returnURL,
		"url":           h.absoluteURL(r, billingPortalSessionPath(customerID, sessionID, returnURL, flowType, billingPortalSessionSubscriptionID(p))),
		"created":       time.Now().UTC().Unix(),
		"livemode":      false,
	}
	if flowType != "" {
		session["flow"] = map[string]any{"type": flowType}
	}
	writeJSON(w, http.StatusOK, session)
}

func billingPortalSessionReturnURL(p params) string {
	return p.first("return_url", "flow_data[after_completion][redirect][return_url]")
}

func billingPortalSessionFlowType(p params) string {
	return p.first("flow_data[type]")
}

func billingPortalSessionSubscriptionID(p params) string {
	return p.first(
		"flow_data[subscription_cancel][subscription]",
		"flow_data[subscription_update][subscription]",
		"flow_data[subscription_update_confirm][subscription]",
	)
}

func billingPortalSessionPath(customerID string, sessionID string, returnURL string, flowType string, subscriptionID string) string {
	query := url.Values{}
	query.Set("customer_id", customerID)
	query.Set("session_id", sessionID)
	if returnURL != "" {
		query.Set("return_url", returnURL)
		query.Set("redirect_on_action", "true")
	}
	if flowType != "" {
		query.Set("flow", flowType)
	}
	if subscriptionID != "" {
		query.Set("subscription_id", subscriptionID)
	}
	return "/portal?" + query.Encode()
}

func (h *Handler) handleCheckoutCompletion(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/checkout/sessions/")
	id, suffix, _ := strings.Cut(rest, "/")
	if id == "" || suffix != "complete" {
		h.notFound(w, r)
		return
	}
	h.completeCheckout(w, r, id)
}

func (h *Handler) completeCheckout(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	previous, _ := h.billing.GetCheckoutSession(r.Context(), id)
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	outcome := p.first("outcome", "payment_method", "payment_method_id", "paymentMethod")
	if outcome == "" {
		outcome = "payment_succeeded"
	}
	session, err := h.billing.CompleteCheckout(r.Context(), id, outcome)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	session.URL = h.absoluteURL(r, session.URL)
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
	switch r.Method {
	case http.MethodGet:
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
	case http.MethodPost:
		subscription, err := h.createSubscriptionFromParams(r)
		writeResult(w, h.stripeSubscription(r, subscription), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) createSubscriptionFromParams(r *http.Request) (billing.Subscription, error) {
	p, err := parseParams(r)
	if err != nil {
		return billing.Subscription{}, err
	}
	if err := validateSubscriptionCreate(p); err != nil {
		return billing.Subscription{}, err
	}
	customerID := p.first("customer", "customer_id")
	items := subscriptionCreateItemsFromParams(p)
	if customerID == "" {
		return billing.Subscription{}, fmt.Errorf("%w: customer is required", billing.ErrInvalidInput)
	}
	if len(items) == 0 {
		return billing.Subscription{}, fmt.Errorf("%w: at least one item is required", billing.ErrInvalidInput)
	}
	customer, err := h.billing.GetCustomer(r.Context(), customerID)
	if err := validateCustomerExists(customer, err); err != nil {
		return billing.Subscription{}, err
	}
	for _, item := range items {
		if err := validatePriceExists(h.billing.GetPrice(r.Context(), item.PriceID)); err != nil {
			return billing.Subscription{}, err
		}
	}
	discounts, err := h.discountsFromParamsOrCustomer(r, p, customer)
	if err != nil {
		return billing.Subscription{}, err
	}
	testClockID := p.string("test_clock")
	if testClockID == "" {
		testClockID = customer.Metadata["test_clock"]
	}
	completionOptions := billing.CheckoutCompletionOptions{}
	if testClockID != "" {
		clock, err := h.billing.GetTestClock(r.Context(), testClockID)
		if err != nil {
			return billing.Subscription{}, err
		}
		completionOptions.At = clock.FrozenTime
	}
	session, err := h.billing.CreateCheckoutSession(r.Context(), billing.CheckoutSession{
		CustomerID: customerID,
		Mode:       "subscription",
		LineItems:  items,
		Discounts:  discounts,
	})
	if err != nil {
		return billing.Subscription{}, err
	}
	completed, err := h.billing.CompleteCheckoutWithOptions(r.Context(), session.ID, p.stringDefault("outcome", "payment_succeeded"), completionOptions)
	if err != nil {
		return billing.Subscription{}, err
	}
	subscription, err := h.billing.GetSubscription(r.Context(), completed.SubscriptionID)
	if err != nil {
		return billing.Subscription{}, err
	}
	metadata := p.metadata()
	if testClockID != "" {
		if metadata == nil {
			metadata = map[string]string{}
		}
		metadata["test_clock"] = testClockID
	}
	for _, key := range []string{"collection_method", "days_until_due", "cancel_at", "billing_cycle_anchor"} {
		if value := p.string(key); value != "" {
			if metadata == nil {
				metadata = map[string]string{}
			}
			metadata[key] = value
		}
	}
	if metadata == nil {
		return subscription, nil
	}
	return h.billing.PatchSubscription(r.Context(), subscription.ID, billing.SubscriptionPatch{Metadata: metadata})
}

func (h *Handler) handleSubscription(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/subscriptions/")
	if strings.HasSuffix(id, "/discount") {
		subscriptionID := strings.TrimSuffix(id, "/discount")
		if subscriptionID == "" || strings.Contains(subscriptionID, "/") {
			h.notFound(w, r)
			return
		}
		h.handleSubscriptionDiscount(w, r, subscriptionID)
		return
	}
	if id == "" || strings.Contains(id, "/") {
		h.notFound(w, r)
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
		if err := validateSubscriptionUpdate(p); err != nil {
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
		metadata := subscriptionUpdateMetadata(p)
		discounts, err := h.discountsFromParams(p)
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		if len(discounts) > 0 {
			metadata = billing.MergeDiscountMetadata(metadata, discounts)
		}
		subscription, err := h.billing.PatchSubscription(r.Context(), id, billing.SubscriptionPatch{
			Items:             items,
			ReplaceItems:      replaceItems,
			Metadata:          metadata,
			CancelAtPeriodEnd: p.boolPtr("cancel_at_period_end"),
		})
		if err == nil {
			h.emitSubscriptionWebhook(r, "customer.subscription.updated", subscription, webhooks.SourceAPI)
			if len(discounts) > 0 {
				h.emitGenericWebhook(r, "customer.discount.created", discounts[0].ID, h.stripeDiscount(discounts[0], subscription.CustomerID, subscription.ID, ""), webhooks.SourceAPI)
			}
		}
		writeResult(w, h.stripeSubscription(r, subscription), err)
	case http.MethodDelete:
		subscription, err := h.billing.CancelPortalSubscription(r.Context(), id, billing.PortalCancel{Mode: "immediate"})
		if err == nil {
			h.emitSubscriptionWebhook(r, "customer.subscription.deleted", subscription, webhooks.SourceAPI)
		}
		writeResult(w, h.stripeSubscription(r, subscription), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST, DELETE")
	}
}

func (h *Handler) handleSubscriptionItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateSubscriptionItemCreate(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	subscriptionID := p.string("subscription")
	priceID := p.first("price", "price_id")
	if subscriptionID == "" {
		writeResult(w, nil, fmt.Errorf("%w: subscription is required", billing.ErrInvalidInput))
		return
	}
	if priceID == "" {
		writeResult(w, nil, fmt.Errorf("%w: price is required", billing.ErrInvalidInput))
		return
	}
	if err := validatePriceExists(h.billing.GetPrice(r.Context(), priceID)); err != nil {
		writeResult(w, nil, err)
		return
	}
	current, err := h.billing.GetSubscription(r.Context(), subscriptionID)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	items := append([]billing.LineItem{}, current.Items...)
	items = append(items, billing.LineItem{PriceID: priceID, Quantity: p.int64Default("quantity", 1)})
	subscription, err := h.billing.PatchSubscription(r.Context(), subscriptionID, billing.SubscriptionPatch{
		Items:        items,
		ReplaceItems: true,
	})
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, h.stripeSubscriptionItem(r, subscription, subscription.Items[len(subscription.Items)-1], len(subscription.Items)-1))
}

func (h *Handler) handleSubscriptionItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/subscription_items/")
	if id == "" || strings.Contains(id, "/") {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodDelete {
		h.methodNotAllowed(w, r, "DELETE")
		return
	}
	subscription, idx, found, err := h.findSubscriptionItem(r, id)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	if !found {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	items := append([]billing.LineItem{}, subscription.Items[:idx]...)
	items = append(items, subscription.Items[idx+1:]...)
	if len(items) == 0 {
		writeResult(w, nil, fmt.Errorf("%w: subscription items cannot be empty", billing.ErrInvalidInput))
		return
	}
	if _, err := h.billing.PatchSubscription(r.Context(), subscription.ID, billing.SubscriptionPatch{Items: items, ReplaceItems: true}); err != nil {
		writeResult(w, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "object": "subscription_item", "deleted": true})
}

func (h *Handler) handleInvoicePreview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/v1/invoices/create_preview" && r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	if r.URL.Path == "/v1/invoices/upcoming" && r.Method != http.MethodGet && r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "GET, POST")
		return
	}
	var (
		p   params
		err error
	)
	if r.Method == http.MethodGet {
		p = params{values: firstValues(r.URL.Query())}
	} else {
		p, err = parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
	}
	if err := validateInvoicePreview(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	preview, err := h.invoicePreview(r.Context(), r.URL.Path, p)
	writeResult(w, preview, err)
}

func (h *Handler) invoicePreview(ctx context.Context, path string, p params) (map[string]any, error) {
	now := time.Now().UTC()
	subscriptionID := p.string("subscription")
	customerID := p.string("customer")
	var subscription billing.Subscription
	var err error
	if subscriptionID != "" {
		subscription, err = h.billing.GetSubscription(ctx, subscriptionID)
		if err != nil {
			return nil, err
		}
		if customerID == "" {
			customerID = subscription.CustomerID
		}
	}
	items := invoicePreviewLineItems(p)
	if len(items) == 0 && subscription.ID != "" {
		items = append([]billing.LineItem{}, subscription.Items...)
	}
	newTotal, currency, prices, err := h.lineItemTotal(ctx, items)
	if err != nil {
		return nil, err
	}
	if currency == "" {
		currency = strings.ToLower(p.stringDefault("currency", "usd"))
	}
	discounts, err := h.invoicePreviewDiscounts(ctx, p, subscription, customerID)
	if err != nil {
		return nil, err
	}
	discountedTotal, discountAmount := billing.ApplyDiscounts(newTotal, currency, discounts)
	behavior := invoicePreviewProrationBehavior(p)
	createdAt := invoicePreviewProrationDate(p, now)
	billingCycleAnchor := invoicePreviewBillingCycleAnchor(p, createdAt)
	amount := discountedTotal
	subtotal := newTotal
	totalDiscountAmount := discountAmount
	lines := []map[string]any{}
	description := "Upcoming invoice preview"
	if subscription.ID != "" {
		oldTotal, oldCurrency, _, err := h.lineItemTotal(ctx, subscription.Items)
		if err != nil {
			return nil, err
		}
		if currency == "" {
			currency = oldCurrency
		}
		oldDiscountedTotal, _ := billing.ApplyDiscounts(oldTotal, currency, discounts)
		amount = 0
		subtotal = 0
		totalDiscountAmount = 0
		if behavior != "none" && !subscription.CurrentPeriodStart.IsZero() && !subscription.CurrentPeriodEnd.IsZero() && subscription.CurrentPeriodEnd.After(createdAt) && subscription.CurrentPeriodEnd.After(subscription.CurrentPeriodStart) {
			periodSeconds := subscription.CurrentPeriodEnd.Unix() - subscription.CurrentPeriodStart.Unix()
			remainingSeconds := subscription.CurrentPeriodEnd.Unix() - createdAt.Unix()
			if periodSeconds > 0 && remainingSeconds > 0 {
				subtotal = (newTotal - oldTotal) * remainingSeconds / periodSeconds
				amount = (discountedTotal - oldDiscountedTotal) * remainingSeconds / periodSeconds
				totalDiscountAmount = subtotal - amount
				if totalDiscountAmount < 0 {
					totalDiscountAmount = 0
				}
			}
		}
		description = "Subscription update preview"
		if amount != 0 {
			priceID := ""
			quantity := int64(1)
			var pricePayload any
			if len(items) > 0 {
				priceID = items[0].PriceID
				quantity = items[0].Quantity
				pricePayload = prices[priceID]
			}
			lines = append(lines, map[string]any{
				"id":               "il_preview_" + sanitizeID(subscription.ID),
				"object":           "line_item",
				"amount":           amount,
				"currency":         currency,
				"description":      "Proration for subscription update",
				"discount_amounts": discountAmounts(discounts, totalDiscountAmount),
				"discountable":     false,
				"period": map[string]any{
					"start": createdAt.Unix(),
					"end":   subscription.CurrentPeriodEnd.Unix(),
				},
				"proration":    true,
				"price":        pricePayload,
				"quantity":     quantity,
				"subscription": subscription.ID,
				"type":         "invoiceitem",
				"parent": map[string]any{
					"type": "subscription_item_details",
					"subscription_item_details": map[string]any{
						"price":        priceID,
						"proration":    true,
						"subscription": subscription.ID,
					},
				},
			})
		}
	}
	if behavior == "none" {
		amount = 0
		lines = nil
	}
	return map[string]any{
		"id":                     "upcoming_in_" + strconv.FormatInt(now.Unix(), 10),
		"object":                 "invoice",
		"customer":               emptyToNil(customerID),
		"subscription":           emptyToNil(subscriptionID),
		"amount_due":             amount,
		"amount_paid":            0,
		"amount_remaining":       amount,
		"subtotal":               subtotal,
		"total":                  amount,
		"discount":               firstDiscountObject(h, discounts, customerID, subscriptionID, ""),
		"discounts":              stripeList(path+"/discounts", discountObjects(h, discounts, customerID, subscriptionID, "")),
		"total_discount_amounts": discountAmounts(discounts, totalDiscountAmount),
		"currency":               currency,
		"created":                now.Unix(),
		"status":                 "draft",
		"lines":                  stripeList(path+"/lines", lines),
		"livemode":               false,
		"description":            description,
		"billtap_preview": map[string]any{
			"proration_behavior":   behavior,
			"proration_date":       createdAt.Unix(),
			"billing_cycle_anchor": billingCycleAnchor.Unix(),
		},
	}, nil
}

func (h *Handler) lineItemTotal(ctx context.Context, items []billing.LineItem) (int64, string, map[string]map[string]any, error) {
	total := int64(0)
	currency := ""
	prices := map[string]map[string]any{}
	for _, item := range items {
		quantity := item.Quantity
		if quantity <= 0 {
			quantity = 1
		}
		price, err := h.billing.GetPrice(ctx, item.PriceID)
		if err != nil {
			return 0, "", nil, err
		}
		total += price.UnitAmount * quantity
		if currency == "" {
			currency = price.Currency
		}
		prices[price.ID] = stripePrice(price)
	}
	return total, currency, prices, nil
}

func invoicePreviewProrationBehavior(p params) string {
	behavior := p.first("subscription_details[proration_behavior]", "subscriptionDetails[prorationBehavior]", "proration_behavior")
	switch behavior {
	case "none", "create_prorations", "always_invoice":
		return behavior
	default:
		return "create_prorations"
	}
}

func invoicePreviewProrationDate(p params, fallback time.Time) time.Time {
	raw := p.first("subscription_details[proration_date]", "subscriptionDetails[prorationDate]", "proration_date")
	if raw == "" {
		return fallback
	}
	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return time.Unix(seconds, 0).UTC()
}

func invoicePreviewBillingCycleAnchor(p params, fallback time.Time) time.Time {
	raw := p.first("subscription_details[billing_cycle_anchor]", "subscriptionDetails[billingCycleAnchor]")
	if raw == "" || raw == "now" {
		return fallback
	}
	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return time.Unix(seconds, 0).UTC()
}

func invoicePreviewLineItems(p params) []billing.LineItem {
	var out []billing.LineItem
	for i := 0; i < 100; i++ {
		price := p.first(
			fmt.Sprintf("subscription_details[items][%d][price]", i),
			fmt.Sprintf("subscriptionDetails[items][%d][price]", i),
			fmt.Sprintf("subscription_items[%d][price]", i),
			fmt.Sprintf("items[%d][price]", i),
		)
		if price == "" && i == 0 {
			price = p.string("price")
		}
		if price == "" {
			continue
		}
		quantity := p.int64Default(fmt.Sprintf("subscription_details[items][%d][quantity]", i), 0)
		if quantity == 0 {
			quantity = p.int64Default(fmt.Sprintf("subscription_items[%d][quantity]", i), 0)
		}
		if quantity == 0 {
			quantity = p.int64Default(fmt.Sprintf("items[%d][quantity]", i), 1)
		}
		out = append(out, billing.LineItem{PriceID: price, Quantity: quantity})
	}
	return out
}

func (h *Handler) handleInvoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
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
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/invoices/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	parts := strings.Split(rest, "/")
	id := parts[0]
	if len(parts) == 2 && parts[1] == "pay" {
		if r.Method != http.MethodPost {
			h.methodNotAllowed(w, r, "POST")
			return
		}
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateInvoicePay(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		paymentMethodID := p.string("payment_method")
		if paymentMethodID == "" {
			paymentMethodID = p.string("source")
		}
		outcome := paymentMethodID
		forgive := p.boolDefault("forgive", false)
		result, err := h.billing.PayInvoice(r.Context(), id, billing.InvoicePaymentOptions{
			Outcome:         outcome,
			PaymentMethodID: paymentMethodID,
			PaidOutOfBand:   p.boolDefault("paid_out_of_band", false) || forgive,
		})
		if err == nil {
			h.emitInvoicePaymentWebhooks(r, result, webhooks.SourceAPI)
		}
		writeResult(w, stripeInvoice(result.Invoice), err)
		return
	}
	if len(parts) != 1 {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	invoice, err := h.billing.GetInvoice(r.Context(), id)
	writeResult(w, stripeInvoice(invoice), err)
}

func (h *Handler) handleRefunds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		refunds, err := h.billing.ListRefunds(r.Context(), billing.RefundFilter{
			ChargeID:        r.URL.Query().Get("charge"),
			PaymentIntentID: r.URL.Query().Get("payment_intent"),
			InvoiceID:       r.URL.Query().Get("invoice"),
			CustomerID:      r.URL.Query().Get("customer"),
		})
		data := make([]map[string]any, 0, len(refunds))
		for _, refund := range refunds {
			data = append(data, stripeRefund(refund))
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateRefundCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		refund, err := h.billing.CreateRefund(r.Context(), billing.Refund{
			ID:              p.string("id"),
			ChargeID:        p.string("charge"),
			PaymentIntentID: p.string("payment_intent"),
			InvoiceID:       p.string("invoice"),
			CustomerID:      p.string("customer"),
			Amount:          p.int64("amount"),
			Currency:        p.string("currency"),
			Reason:          p.string("reason"),
			Status:          p.string("status"),
			Metadata:        p.metadata(),
		})
		if err == nil {
			h.emitRefundWebhooks(r, refund, webhooks.SourceAPI)
		}
		writeResult(w, stripeRefund(refund), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleRefund(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/refunds/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	id, action, hasAction := strings.Cut(rest, "/")
	if hasAction && action == "cancel" {
		if r.Method != http.MethodPost {
			h.methodNotAllowed(w, r, "POST")
			return
		}
		refund, err := h.billing.UpdateRefundStatus(r.Context(), id, "canceled", time.Time{})
		if err == nil {
			h.emitGenericWebhook(r, "charge.refund.updated", refund.ID, stripeRefund(refund), webhooks.SourceAPI)
		}
		writeResult(w, stripeRefund(refund), err)
		return
	}
	if hasAction {
		h.notFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		refund, err := h.billing.GetRefund(r.Context(), id)
		writeResult(w, stripeRefund(refund), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateRefundUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		status := p.stringDefault("status", "succeeded")
		refund, err := h.billing.UpdateRefundStatus(r.Context(), id, status, time.Time{})
		if err == nil {
			h.emitGenericWebhook(r, "charge.refund.updated", refund.ID, stripeRefund(refund), webhooks.SourceAPI)
		}
		writeResult(w, stripeRefund(refund), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
		return
	}
}

func (h *Handler) handleCreditNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		notes, err := h.billing.ListCreditNotes(r.Context(), billing.CreditNoteFilter{
			InvoiceID:  r.URL.Query().Get("invoice"),
			CustomerID: r.URL.Query().Get("customer"),
		})
		data := make([]map[string]any, 0, len(notes))
		for _, note := range notes {
			data = append(data, stripeCreditNote(note))
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateCreditNoteCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		note, err := h.billing.CreateCreditNote(r.Context(), billing.CreditNote{
			ID:         p.string("id"),
			InvoiceID:  p.string("invoice"),
			CustomerID: p.string("customer"),
			Amount:     p.int64("amount"),
			Currency:   p.string("currency"),
			Reason:     p.string("reason"),
			Status:     p.string("status"),
			Metadata:   p.metadata(),
		})
		if err == nil {
			h.emitGenericWebhook(r, "credit_note.created", note.ID, stripeCreditNote(note), webhooks.SourceAPI)
		}
		writeResult(w, stripeCreditNote(note), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleCreditNote(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/credit_notes/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	id, action, hasAction := strings.Cut(rest, "/")
	if hasAction && action == "void" {
		if r.Method != http.MethodPost {
			h.methodNotAllowed(w, r, "POST")
			return
		}
		note, err := h.billing.VoidCreditNote(r.Context(), id)
		if err == nil {
			h.emitGenericWebhook(r, "credit_note.voided", note.ID, stripeCreditNote(note), webhooks.SourceAPI)
		}
		writeResult(w, stripeCreditNote(note), err)
		return
	}
	if hasAction && action != "" {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	note, err := h.billing.GetCreditNote(r.Context(), id)
	writeResult(w, stripeCreditNote(note), err)
}

func (h *Handler) handlePaymentIntents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validatePaymentIntentCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		customerID := p.string("customer")
		var customer billing.Customer
		if customerID != "" {
			var customerErr error
			customer, customerErr = h.billing.GetCustomer(r.Context(), customerID)
			if err := validateCustomerExists(customer, customerErr); err != nil {
				writeResult(w, nil, err)
				return
			}
		}
		if p.boolDefault("confirm", false) && !p.hasAny("payment_method", "outcome") && !hasPaymentIntentDeferredOutcome(p) && billing.CustomerDefaultPaymentIntentOutcome(customer.Metadata) == "" {
			writeError(w, http.StatusBadRequest, missingParam("payment_method"))
			return
		}
		metadata := paymentIntentMetadata(p)
		intent, err := h.billing.CreatePaymentIntent(r.Context(), billing.PaymentIntent{
			ID:              p.string("id"),
			CustomerID:      customerID,
			Amount:          p.int64("amount"),
			Currency:        p.string("currency"),
			CaptureMethod:   p.stringDefault("capture_method", "automatic"),
			PaymentMethodID: p.string("payment_method"),
			Metadata:        metadata,
		})
		if err == nil {
			h.emitPaymentIntentWebhook(r, "payment_intent.created", intent)
		}
		if err == nil && (p.boolDefault("confirm", false) || p.has("outcome")) {
			intent, err = h.billing.ConfirmPaymentIntent(r.Context(), intent.ID, p.string("payment_method"), p.string("outcome"))
			if err == nil {
				h.emitPaymentIntentWebhook(r, paymentIntentTerminalEvent(intent.Status), intent)
			}
		}
		writeResult(w, stripePaymentIntent(intent), err)
	case http.MethodGet:
		items, err := h.billing.ListPaymentIntents(r.Context())
		data := make([]map[string]any, 0, len(items))
		for _, item := range items {
			if customer := r.URL.Query().Get("customer"); customer != "" && item.CustomerID != customer {
				continue
			}
			data = append(data, stripePaymentIntent(item))
			if limit := queryInt(r, "limit"); limit > 0 && len(data) >= limit {
				break
			}
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handlePaymentIntent(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/v1/payment_intents/")
	id, action, hasAction := strings.Cut(rest, "/")
	if id == "" {
		h.notFound(w, r)
		return
	}
	if !hasAction {
		if r.Method != http.MethodGet {
			h.methodNotAllowed(w, r, "GET")
			return
		}
		paymentIntent, err := h.billing.GetPaymentIntent(r.Context(), id)
		writeResult(w, stripePaymentIntent(paymentIntent), err)
		return
	}
	if strings.Contains(action, "/") || r.Method != http.MethodPost {
		if r.Method != http.MethodPost {
			h.methodNotAllowed(w, r, "POST")
			return
		}
		h.notFound(w, r)
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var intent billing.PaymentIntent
	switch action {
	case "confirm":
		if err := validatePaymentIntentConfirm(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		intent, err = h.billing.ConfirmPaymentIntent(r.Context(), id, p.string("payment_method"), p.string("outcome"))
		if err == nil {
			h.emitPaymentIntentWebhook(r, paymentIntentTerminalEvent(intent.Status), intent)
		}
	case "capture":
		if err := validatePaymentIntentCapture(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		intent, err = h.billing.CapturePaymentIntent(r.Context(), id, p.int64("amount_to_capture"))
		if err == nil {
			h.emitPaymentIntentWebhook(r, "payment_intent.succeeded", intent)
		}
	case "cancel":
		if err := validatePaymentIntentCancel(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		intent, err = h.billing.CancelPaymentIntent(r.Context(), id)
		if err == nil {
			h.emitPaymentIntentWebhook(r, "payment_intent.canceled", intent)
		}
	default:
		h.notFound(w, r)
		return
	}
	writeResult(w, stripePaymentIntent(intent), err)
}

func (h *Handler) handleSetupIntents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateSetupIntentCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		customerID := p.string("customer")
		if customerID != "" {
			if err := validateCustomerExists(h.billing.GetCustomer(r.Context(), customerID)); err != nil {
				writeResult(w, nil, err)
				return
			}
		}
		intent, err := h.billing.CreateSetupIntent(r.Context(), billing.SetupIntent{
			ID:              p.string("id"),
			CustomerID:      customerID,
			Usage:           p.stringDefault("usage", "off_session"),
			PaymentMethodID: p.string("payment_method"),
		})
		if err == nil {
			h.emitSetupIntentWebhook(r, "setup_intent.created", intent)
		}
		if err == nil && p.boolDefault("confirm", false) {
			intent, err = h.billing.ConfirmSetupIntent(r.Context(), intent.ID, p.string("payment_method"), p.string("outcome"))
			if err == nil {
				h.emitSetupIntentWebhook(r, setupIntentTerminalEvent(intent.Status), intent)
			}
		}
		writeResult(w, stripeSetupIntent(intent), err)
	case http.MethodGet:
		items, err := h.billing.ListSetupIntents(r.Context())
		data := make([]map[string]any, 0, len(items))
		for _, item := range items {
			if customer := r.URL.Query().Get("customer"); customer != "" && item.CustomerID != customer {
				continue
			}
			data = append(data, stripeSetupIntent(item))
			if limit := queryInt(r, "limit"); limit > 0 && len(data) >= limit {
				break
			}
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleSetupIntent(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/v1/setup_intents/")
	id, action, hasAction := strings.Cut(rest, "/")
	if id == "" {
		h.notFound(w, r)
		return
	}
	if !hasAction {
		if r.Method != http.MethodGet {
			h.methodNotAllowed(w, r, "GET")
			return
		}
		intent, err := h.billing.GetSetupIntent(r.Context(), id)
		writeResult(w, stripeSetupIntent(intent), err)
		return
	}
	if strings.Contains(action, "/") || r.Method != http.MethodPost {
		if r.Method != http.MethodPost {
			h.methodNotAllowed(w, r, "POST")
			return
		}
		h.notFound(w, r)
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var intent billing.SetupIntent
	switch action {
	case "confirm":
		if err := validateSetupIntentConfirm(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		intent, err = h.billing.ConfirmSetupIntent(r.Context(), id, p.string("payment_method"), p.string("outcome"))
		if err == nil {
			h.emitSetupIntentWebhook(r, setupIntentTerminalEvent(intent.Status), intent)
		}
	case "cancel":
		if err := validateSetupIntentCancel(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		intent, err = h.billing.CancelSetupIntent(r.Context(), id)
		if err == nil {
			h.emitSetupIntentWebhook(r, "setup_intent.canceled", intent)
		}
	default:
		h.notFound(w, r)
		return
	}
	writeResult(w, stripeSetupIntent(intent), err)
}

func (h *Handler) handleTestClocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		clocks, err := h.billing.ListTestClocks(r.Context())
		data := make([]map[string]any, 0, len(clocks))
		for _, clock := range clocks {
			data = append(data, stripeTestClock(clock))
		}
		writeResult(w, stripeList(r.URL.Path, data), err)
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateTestClockCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		frozenTime, err := parseTimestampParam(p.first("frozen_time", "frozenTime"))
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		clock, err := h.billing.CreateTestClock(r.Context(), billing.TestClock{
			ID:         p.string("id"),
			Name:       p.string("name"),
			FrozenTime: frozenTime,
		})
		writeResult(w, stripeTestClock(clock), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleTestClock(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/test_helpers/test_clocks/"), "/")
	if rest == "" {
		h.notFound(w, r)
		return
	}
	id, action, hasAction := strings.Cut(rest, "/")
	if !hasAction {
		if r.Method != http.MethodGet {
			h.methodNotAllowed(w, r, "GET")
			return
		}
		clock, err := h.billing.GetTestClock(r.Context(), id)
		writeResult(w, stripeTestClock(clock), err)
		return
	}
	if action != "advance" || strings.Contains(action, "/") {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateTestClockAdvance(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	frozenTime, err := parseTimestampParam(p.first("frozen_time", "frozenTime"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	clock, advance, err := h.billing.AdvanceTestClock(r.Context(), id, frozenTime)
	if err == nil {
		h.emitClockAdvanceWebhooks(r, advance)
		if scheduled := h.applyDueSubscriptionSchedules(r, id, frozenTime); len(scheduled) > 0 {
			advance.Scheduled = append(advance.Scheduled, scheduled...)
			advance.ScheduledCount += len(scheduled)
			advance.Processed += len(scheduled)
		}
	}
	response := stripeTestClock(clock)
	if err == nil {
		response["billtap_advance_result"] = advance
	}
	writeResult(w, response, err)
}

func (h *Handler) handlePaymentMethods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	customerID := r.URL.Query().Get("customer")
	h.writeCustomerPaymentMethods(w, r, customerID)
}

func (h *Handler) handleCustomerPaymentMethods(w http.ResponseWriter, r *http.Request, customerID string) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	h.writeCustomerPaymentMethods(w, r, customerID)
}

func (h *Handler) writeCustomerPaymentMethods(w http.ResponseWriter, r *http.Request, customerID string) {
	p := paramsFromValues(r.URL.Query())
	if err := validatePaymentMethodList(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if customerID == "" {
		writeResult(w, stripeList(r.URL.Path, []map[string]any{}), nil)
		return
	}
	customer, err := h.billing.GetCustomer(r.Context(), customerID)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	if paymentMethodType := strings.TrimSpace(r.URL.Query().Get("type")); paymentMethodType != "" && paymentMethodType != "card" {
		writeResult(w, stripeList(r.URL.Path, []map[string]any{}), nil)
		return
	}
	writeResult(w, stripeList(r.URL.Path, stripePaymentMethods(customer)), nil)
}

func (h *Handler) handleObjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
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
		if err := validateWebhookEndpointCreate(p); err != nil {
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
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/v1/webhook_endpoints/")
	id, action, hasAction := strings.Cut(rest, "/")
	if id == "" {
		h.notFound(w, r)
		return
	}
	if hasAction {
		if action != "attempts" {
			h.notFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			h.methodNotAllowed(w, r, "GET")
			return
		}
		attempts, err := h.webhooks.ListDeliveryAttempts(r.Context(), webhooks.DeliveryAttemptFilter{EndpointID: id, Status: r.URL.Query().Get("status")})
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"object": "list", "data": h.deliveryAttemptResponses(r, attempts)})
		return
	}
	switch r.Method {
	case http.MethodGet:
		endpoint, err := h.webhooks.GetEndpoint(r.Context(), id)
		writeResult(w, maskEndpoint(endpoint), err)
	case http.MethodPost, http.MethodPatch:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateWebhookEndpointUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		current, err := h.webhooks.GetEndpoint(r.Context(), id)
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		active := current.Active
		if p.has("active") {
			active = p.boolDefault("active", current.Active)
		}
		if p.has("enabled") {
			active = p.boolDefault("enabled", current.Active)
		}
		endpoint, err := h.webhooks.UpdateEndpoint(r.Context(), id, webhooks.Endpoint{
			URL:              p.string("url"),
			Secret:           p.string("secret"),
			EnabledEvents:    p.list("enabled_events"),
			RetryMaxAttempts: int(p.int64("retry_max_attempts")),
			RetryBackoff:     p.list("retry_backoff"),
			Active:           active,
		})
		writeResult(w, maskEndpoint(endpoint), err)
	case http.MethodDelete:
		endpoint, err := h.webhooks.DeleteEndpoint(r.Context(), id)
		writeResult(w, maskEndpoint(endpoint), err)
	default:
		h.methodNotAllowed(w, r, "GET, POST, PATCH, DELETE")
	}
}

func (h *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeJSON(w, http.StatusOK, map[string]any{"object": "list", "data": []any{}})
		return
	}
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	q := r.URL.Query()
	events, err := h.webhooks.ListEvents(r.Context(), webhooks.EventFilter{
		Type:          q.Get("type"),
		ScenarioRunID: q.Get("scenarioRunId"),
	})
	if err == nil {
		events = filterEventsByQuery(events, q)
	}
	writeResult(w, map[string]any{"object": "list", "data": events}, err)
}

func filterEventsByQuery(events []webhooks.Event, q url.Values) []webhooks.Event {
	createdGTE := queryInt64Values(q, "created[gte]", "created_gte")
	createdGT := queryInt64Values(q, "created[gt]", "created_gt")
	createdLTE := queryInt64Values(q, "created[lte]", "created_lte")
	createdLT := queryInt64Values(q, "created[lt]", "created_lt")
	customerID := firstValue(q, "data.object.customer", "customer")
	metadataFilters := eventMetadataFilters(q)
	limit := queryIntFromValues(q, "limit")

	out := make([]webhooks.Event, 0, len(events))
	for _, event := range events {
		if createdGTE != 0 && event.Created < createdGTE {
			continue
		}
		if createdGT != 0 && event.Created <= createdGT {
			continue
		}
		if createdLTE != 0 && event.Created > createdLTE {
			continue
		}
		if createdLT != 0 && event.Created >= createdLT {
			continue
		}
		object := eventObjectMap(event)
		if customerID != "" && stringFromObject(object, "customer") != customerID {
			continue
		}
		if !eventMetadataMatches(object, metadataFilters) {
			continue
		}
		out = append(out, event)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

func queryInt64Values(q url.Values, keys ...string) int64 {
	for _, key := range keys {
		value := strings.TrimSpace(q.Get(key))
		if value == "" {
			continue
		}
		parsed, _ := strconv.ParseInt(value, 10, 64)
		return parsed
	}
	return 0
}

func queryIntFromValues(q url.Values, key string) int {
	parsed, _ := strconv.Atoi(strings.TrimSpace(q.Get(key)))
	return parsed
}

func eventObjectMap(event webhooks.Event) map[string]any {
	var object map[string]any
	_ = json.Unmarshal(event.Data.Object, &object)
	return object
}

func stringFromObject(object map[string]any, key string) string {
	if value, ok := object[key].(string); ok {
		return value
	}
	return ""
}

func eventMetadataFilters(q url.Values) map[string]string {
	out := map[string]string{}
	for key, values := range q {
		if !strings.HasPrefix(key, "data.object.metadata[") && !strings.HasPrefix(key, "metadata[") {
			continue
		}
		inner := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(key, "data.object.metadata["), "metadata["), "]")
		if inner != "" && len(values) > 0 {
			out[inner] = values[0]
		}
	}
	return out
}

func eventMetadataMatches(object map[string]any, filters map[string]string) bool {
	if len(filters) == 0 {
		return true
	}
	metadata, _ := object["metadata"].(map[string]any)
	for key, value := range filters {
		if fmt.Sprint(metadata[key]) != value {
			return false
		}
	}
	return true
}

func (h *Handler) handleEvent(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/events/")
	if id == "" || strings.Contains(id, "/") {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	event, err := h.webhooks.GetEvent(r.Context(), id)
	writeResult(w, event, err)
}

func (h *Handler) handleTimeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	entries, err := h.billing.Timeline(r.Context(), diagnosticTimelineFilter(r))
	writeResult(w, map[string]any{"object": "list", "data": entries}, err)
}

func (h *Handler) handleDebugBundles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
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
	requestTraces, _ := h.listRequestTraces(r.Context(), diagnostics.RequestTraceFilter{ObjectID: targetID, Limit: 100})
	writeJSON(w, http.StatusOK, map[string]any{
		"id":                "dbg_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		"object":            "debug_bundle",
		"target":            compactStringMap(map[string]string{"type": targetType, "id": targetID}),
		"filters":           timelineFilterMap(filter),
		"objects":           objects,
		"timeline":          timeline,
		"request_traces":    requestTraces,
		"webhook_events":    events,
		"delivery_attempts": h.deliveryAttemptResponses(r, attempts),
		"created_at":        time.Now().UTC(),
	})
}

func (h *Handler) handleRequestTraces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	traces, err := h.listRequestTraces(r.Context(), requestTraceFilter(r))
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"object": "list", "data": traces})
}

func (h *Handler) handleDiagnostics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	filter := diagnosticTimelineFilter(r)
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
	requestTraces, err := h.listRequestTraces(r.Context(), requestTraceFilter(r))
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	var events []webhooks.Event
	var attempts []map[string]any
	if h.webhooks != nil {
		events, _ = h.webhooks.ListEvents(r.Context(), webhooks.EventFilter{})
		rawAttempts, _ := h.webhooks.ListDeliveryAttempts(r.Context(), webhooks.DeliveryAttemptFilter{})
		attempts = h.deliveryAttemptResponses(r, rawAttempts)
	}
	snapshot, _ := fixtures.NewService(h.billing).Snapshot(r.Context(), fixtureSnapshotFilter(r))
	writeJSON(w, http.StatusOK, map[string]any{
		"id":                "diag_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		"object":            "diagnostic_bundle",
		"summary":           diagnosticsSummary(objects, timeline, requestTraces, events, attempts),
		"filters":           timelineFilterMap(filter),
		"objects":           objects,
		"fixture_snapshot":  snapshot,
		"timeline":          timeline,
		"request_traces":    requestTraces,
		"webhook_events":    events,
		"delivery_attempts": attempts,
		"created_at":        time.Now().UTC(),
	})
}

func (h *Handler) handleFixtureApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
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
	if disputes, err := h.applyFixtureDisputes(r, pack); err != nil {
		writeResult(w, result, err)
		return
	} else if len(disputes) > 0 {
		result.Disputes = disputes
		result.Summary["disputes"] = len(disputes)
	}
	h.emitFixtureApplyWebhooks(r, result)
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) handleFixtureValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
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
	if err := fixtures.NewService(h.billing).Validate(pack); err != nil {
		writeResult(w, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":     "fxval_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 36),
		"object": "fixture_validation",
		"valid":  true,
		"summary": map[string]int{
			"customers":          len(pack.Customers),
			"products":           len(pack.Products) + len(pack.Catalog.Products),
			"prices":             len(pack.Prices) + len(pack.Catalog.Prices),
			"connected_accounts": len(pack.ConnectedAccounts),
			"test_clocks":        len(pack.TestClocks),
			"subscriptions":      len(pack.Subscriptions),
			"refunds":            len(pack.Refunds),
			"credit_notes":       len(pack.CreditNotes),
			"disputes":           len(pack.Disputes),
		},
	})
}

func (h *Handler) applyFixtureDisputes(r *http.Request, pack fixtures.Pack) ([]map[string]any, error) {
	if len(pack.Disputes) == 0 {
		return nil, nil
	}
	out := make([]map[string]any, 0, len(pack.Disputes))
	for _, fixture := range pack.Disputes {
		dispute := disputeFixturePayload(fixture)
		h.local.mu.Lock()
		h.local.disputes[fmt.Sprint(dispute["id"])] = dispute
		h.local.mu.Unlock()
		out = append(out, cloneEvidence(dispute))
		h.emitGenericWebhook(r, "charge.dispute.created", fmt.Sprint(dispute["id"]), dispute, webhooks.SourceFixture)
		if fmt.Sprint(dispute["status"]) != "needs_response" {
			h.emitGenericWebhook(r, "charge.dispute.updated", fmt.Sprint(dispute["id"]), dispute, webhooks.SourceFixture)
			h.emitGenericWebhook(r, "charge.dispute.funds_withdrawn", fmt.Sprint(dispute["id"]), dispute, webhooks.SourceFixture)
		}
		if status := fmt.Sprint(dispute["status"]); status == "won" || status == "lost" {
			h.emitGenericWebhook(r, "charge.dispute.closed", fmt.Sprint(dispute["id"]), dispute, webhooks.SourceFixture)
		}
	}
	return out, nil
}

func disputeFixturePayload(fixture fixtures.DisputeFixture) map[string]any {
	now := time.Now().UTC()
	id := strings.TrimSpace(fixture.ID)
	if id == "" {
		id = "dp_" + strconv.FormatInt(now.UnixNano(), 36)
	}
	amount := fixture.Amount
	if amount <= 0 {
		amount = 1000
	}
	currency := strings.ToLower(strings.TrimSpace(fixture.Currency))
	if currency == "" {
		currency = "usd"
	}
	status := strings.ToLower(strings.TrimSpace(fixture.Status))
	if status == "" {
		status = "needs_response"
	}
	return map[string]any{
		"id":                   id,
		"object":               "dispute",
		"charge":               strings.TrimSpace(fixture.Charge),
		"amount":               amount,
		"currency":             currency,
		"reason":               firstNonEmptyString(fixture.Reason, "general"),
		"status":               status,
		"evidence":             map[string]any{},
		"evidence_details":     map[string]any{"has_evidence": false, "submission_count": 0, "past_due": false},
		"balance_transactions": []map[string]any{},
		"is_charge_refundable": true,
		"metadata":             nonNilMap(fixture.Metadata),
		"created":              now.Unix(),
		"livemode":             false,
	}
}

func (h *Handler) handleFixtureSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	snapshot, err := fixtures.NewService(h.billing).Snapshot(r.Context(), fixtureSnapshotFilter(r))
	writeResult(w, snapshot, err)
}

func (h *Handler) handleFixtureResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	result, err := fixtures.NewService(h.billing).Resolve(r.Context(), fixtures.ResolveFilter{
		Ref:         firstQuery(r, "ref", "id", "lookup_key", "lookupKey"),
		RunID:       firstQuery(r, "runId", "run_id"),
		FixtureName: firstQuery(r, "fixtureName", "fixture_name", "name"),
		Namespace:   firstQuery(r, "namespace"),
		TenantID:    firstQuery(r, "tenantId", "tenant_id"),
	})
	writeResult(w, result, err)
}

func (h *Handler) handleFixtureAssert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
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
	objectType := dashboardObjectType(p.first("objectType", "object_type", "targetType", "target_type", "type"))
	objectID := p.first("objectId", "object_id", "targetId", "target_id", "id")
	filter := billing.TimelineFilter{
		CustomerID:        p.first("customerId", "customer_id"),
		CheckoutSessionID: p.first("checkoutSessionId", "checkout_session_id"),
		SubscriptionID:    p.first("subscriptionId", "subscription_id"),
		InvoiceID:         p.first("invoiceId", "invoice_id"),
		PaymentIntentID:   p.first("paymentIntentId", "payment_intent_id"),
		ObjectType:        objectType,
	}

	if objectID == "" {
		return filter
	}
	switch objectType {
	case "customer":
		filter.ObjectType = ""
		if filter.CustomerID == "" {
			filter.CustomerID = objectID
		}
	case "checkout_session":
		filter.ObjectType = ""
		if filter.CheckoutSessionID == "" {
			filter.CheckoutSessionID = objectID
		}
	case "subscription":
		filter.ObjectType = ""
		if filter.SubscriptionID == "" {
			filter.SubscriptionID = objectID
		}
	case "invoice":
		filter.ObjectType = ""
		if filter.InvoiceID == "" {
			filter.InvoiceID = objectID
		}
	case "payment_intent":
		filter.ObjectType = ""
		if filter.PaymentIntentID == "" {
			filter.PaymentIntentID = objectID
		}
	default:
		if objectType != "" {
			filter.ObjectID = objectID
			return filter
		}
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

func diagnosticTimelineFilter(r *http.Request) billing.TimelineFilter {
	q := r.URL.Query()
	return debugBundleTimelineFilter(params{values: map[string]string{
		"customerId":        firstValue(q, "customerId", "customer_id", "customer"),
		"checkoutSessionId": firstValue(q, "checkoutSessionId", "checkout_session_id", "checkoutSession"),
		"subscriptionId":    firstValue(q, "subscriptionId", "subscription_id", "subscription"),
		"invoiceId":         firstValue(q, "invoiceId", "invoice_id", "invoice"),
		"paymentIntentId":   firstValue(q, "paymentIntentId", "payment_intent_id", "paymentIntent"),
		"objectType":        firstValue(q, "objectType", "object_type", "targetType", "target_type", "type"),
		"objectId":          firstValue(q, "objectId", "object_id", "targetId", "target_id", "id"),
	}})
}

func firstValue(values url.Values, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(values.Get(key)); value != "" {
			return value
		}
	}
	return ""
}

func requestTraceFilter(r *http.Request) diagnostics.RequestTraceFilter {
	status, _ := strconv.Atoi(r.URL.Query().Get("status"))
	limit := queryInt(r, "limit")
	if limit == 0 {
		limit = 100
	}
	return diagnostics.RequestTraceFilter{
		Method:         strings.ToUpper(firstQuery(r, "method")),
		Path:           firstQuery(r, "path"),
		Status:         status,
		RequestID:      firstQuery(r, "requestId", "request_id"),
		IdempotencyKey: firstQuery(r, "idempotencyKey", "idempotency_key"),
		ObjectID:       firstQuery(r, "objectId", "object_id", "targetId", "target_id", "id"),
		Limit:          limit,
	}
}

func (h *Handler) listRequestTraces(ctx context.Context, filter diagnostics.RequestTraceFilter) ([]diagnostics.RequestTrace, error) {
	if h.diagnostics == nil {
		return []diagnostics.RequestTrace{}, nil
	}
	return h.diagnostics.ListRequestTraces(ctx, filter)
}

func diagnosticsSummary(objects map[string]any, timeline []billing.TimelineEntry, traces []diagnostics.RequestTrace, events []webhooks.Event, attempts []map[string]any) map[string]any {
	summary := map[string]any{
		"timeline":          len(timeline),
		"request_traces":    len(traces),
		"webhook_events":    len(events),
		"delivery_attempts": len(attempts),
	}
	for _, key := range []string{"customers", "products", "prices", "checkout_sessions", "subscriptions", "invoices", "payment_intents", "webhook_endpoints"} {
		if count, ok := diagnosticObjectCount(objects[key]); ok {
			summary[key] = count
		}
	}
	return summary
}

func diagnosticObjectCount(value any) (int, bool) {
	switch v := value.(type) {
	case []billing.Customer:
		return len(v), true
	case []billing.Product:
		return len(v), true
	case []billing.Price:
		return len(v), true
	case []billing.CheckoutSession:
		return len(v), true
	case []billing.Subscription:
		return len(v), true
	case []billing.Invoice:
		return len(v), true
	case []billing.PaymentIntent:
		return len(v), true
	case []webhooks.Endpoint:
		return len(v), true
	default:
		return 0, false
	}
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
		"object_type":         filter.ObjectType,
		"object_id":           filter.ObjectID,
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
		h.methodNotAllowed(w, r, "GET")
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
		h.notFound(w, r)
		return
	}
	if !hasAction {
		if r.Method != http.MethodGet {
			h.methodNotAllowed(w, r, "GET")
			return
		}
		state, err := h.billing.PortalState(r.Context(), customerID)
		writeResult(w, state, err)
		return
	}
	if action != "payment-method" {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := h.billing.SimulatePaymentMethodUpdate(r.Context(), customerID, p.first("outcome", "simulate", "status"), p.first("payment_method", "payment_method_id", "paymentMethod", "paymentMethodID"))
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	if result.Status == "succeeded" && result.PaymentMethodID != "" {
		h.emitGenericWebhook(r, "payment_method.attached", result.PaymentMethodID, stripePaymentMethod(customerID, result.PaymentMethodID), webhooks.SourcePortal)
		if customer, err := h.billing.GetCustomer(r.Context(), customerID); err == nil {
			h.emitGenericWebhook(r, "customer.updated", customer.ID, stripeCustomer(customer), webhooks.SourcePortal)
		}
	}
	state, stateErr := h.billing.PortalState(r.Context(), customerID)
	writeResult(w, map[string]any{"object": "portal_action", "action": "payment_method", "payment_method": result, "state": state}, stateErr)
}

func (h *Handler) handlePortalSubscription(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/portal/subscriptions/")
	subscriptionID, action, found := strings.Cut(rest, "/")
	if subscriptionID == "" || !found {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
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
		h.notFound(w, r)
		return
	}
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	eventType := "customer.subscription.updated"
	if action == "cancel" && sub.Status == "canceled" {
		eventType = "customer.subscription.deleted"
	}
	h.emitSubscriptionWebhook(r, eventType, sub, webhooks.SourcePortal)
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
		h.methodNotAllowed(w, r, "GET")
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
		h.methodNotAllowed(w, r, "GET")
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
		h.methodNotAllowed(w, r, "POST")
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
		h.methodNotAllowed(w, r, "POST")
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
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
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
	simulateAppFailure := replaySimulatedAppFailure(p)
	responseStatus := int(p.int64("response_status"))
	if responseStatus == 0 {
		responseStatus = int(p.int64("responseStatus"))
	}
	if responseStatus == 0 {
		responseStatus = int(p.int64("force_response_status"))
	}
	attempts, err := h.webhooks.ReplayEvent(r.Context(), id, webhooks.ReplayOptions{
		Duplicate:          int(p.int64Default("duplicate", 1)),
		Delay:              delay,
		OutOfOrder:         p.boolDefault("out_of_order", false),
		ResponseStatus:     responseStatus,
		ResponseBody:       p.first("response_body", "responseBody", "body"),
		SimulatedError:     p.first("error", "simulated_error"),
		SimulatedTimeout:   p.boolDefault("timeout", false),
		SignatureMismatch:  p.boolDefault("signature_mismatch", false),
		SimulateAppFailure: simulateAppFailure,
	})
	writeResult(w, map[string]any{"message": "replay scheduled", "object": "list", "data": h.deliveryAttemptResponses(r, attempts)}, err)
}

func (h *Handler) handlePaymentIntentAction(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/payment_intents/")
	id, action, found := strings.Cut(rest, "/")
	if id == "" || !found {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	var (
		intent billing.PaymentIntent
		err    error
	)
	switch action {
	case "complete_action":
		current, currentErr := h.billing.GetPaymentIntent(r.Context(), id)
		if currentErr != nil {
			writeResult(w, nil, currentErr)
			return
		}
		intent, err = h.billing.ConfirmPaymentIntent(r.Context(), id, firstNonEmptyString(current.PaymentMethodID, "pm_card_visa"), "payment_succeeded")
		if err == nil {
			h.emitPaymentIntentWebhook(r, "payment_intent.succeeded", intent)
		}
	case "cancel_action":
		intent, err = h.billing.CancelPaymentIntent(r.Context(), id)
		if err == nil {
			h.emitPaymentIntentWebhook(r, "payment_intent.canceled", intent)
		}
	case "outcome":
		p, parseErr := parseParams(r)
		if parseErr != nil {
			writeError(w, http.StatusBadRequest, parseErr)
			return
		}
		if validateErr := validatePaymentIntentOutcomeUpdate(p); validateErr != nil {
			writeError(w, http.StatusBadRequest, validateErr)
			return
		}
		intent, err = h.billing.SetPaymentIntentOutcome(r.Context(), id, p.string("outcome"))
	default:
		h.notFound(w, r)
		return
	}
	writeResult(w, stripePaymentIntent(intent), err)
}

func (h *Handler) handleEventReplayGroup(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	eventIDs := historicalReplayEventIDs(p)
	if len(eventIDs) == 0 {
		writeError(w, http.StatusBadRequest, missingParam("event_ids"))
		return
	}
	mode := strings.ToLower(p.stringDefault("mode", "ordered"))
	if mode == "out_of_order" || mode == "reverse" {
		for left, right := 0, len(eventIDs)-1; left < right; left, right = left+1, right-1 {
			eventIDs[left], eventIDs[right] = eventIDs[right], eventIDs[left]
		}
	}
	omit := map[string]bool{}
	for _, id := range p.list("omit_event_ids") {
		omit[id] = true
	}
	delaySeconds := p.int64("delay_seconds")
	var attempts []webhooks.DeliveryAttempt
	for idx, eventID := range eventIDs {
		if omit[eventID] {
			continue
		}
		delay := time.Duration(delaySeconds) * time.Second
		if mode == "delayed" && delay == 0 {
			delay = time.Duration(idx+1) * time.Second
		}
		next, err := h.webhooks.ReplayEvent(r.Context(), eventID, webhooks.ReplayOptions{
			Duplicate:         int(p.int64Default("duplicate", 1)),
			Delay:             delay,
			OutOfOrder:        mode == "out_of_order" || mode == "reverse",
			SignatureMismatch: p.boolDefault("signature_mismatch", false) || p.string("simulate") == "invalid_signature",
		})
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		attempts = append(attempts, next...)
	}
	writeJSON(w, http.StatusOK, map[string]any{"object": "event_replay_group", "mode": mode, "data": h.deliveryAttemptResponses(r, attempts)})
}

func historicalReplayEventIDs(p params) []string {
	var out []string
	for _, key := range []string{"event_ids", "eventIds", "events"} {
		out = append(out, p.list(key)...)
	}
	if id := p.string("event_id"); id != "" {
		out = append(out, id)
	}
	return out
}

func (h *Handler) handleWebhookEndpointAction(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		writeResult(w, nil, webhooks.ErrNotFound)
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/api/webhooks/endpoints/")
	endpointID, action, found := strings.Cut(rest, "/")
	if endpointID == "" || !found || action != "replay-historical" {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateHistoricalReplay(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	since, err := parseOptionalReplayTime(p.first("since", "created_after", "createdAfter"))
	if err != nil {
		writeError(w, http.StatusBadRequest, invalidParam("since", "Expected RFC3339 timestamp, Unix timestamp, or now."))
		return
	}
	until, err := parseOptionalReplayTime(p.first("until", "created_before", "createdBefore"))
	if err != nil {
		writeError(w, http.StatusBadRequest, invalidParam("until", "Expected RFC3339 timestamp, Unix timestamp, or now."))
		return
	}
	limit := int(p.int64("limit"))
	result, err := h.webhooks.ReplayHistoricalForEndpoint(r.Context(), endpointID, webhooks.HistoricalReplayOptions{
		Since:      since,
		Until:      until,
		EventTypes: historicalReplayEventTypes(p),
		Limit:      limit,
		Force:      p.boolDefault("force", false),
	})
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"object":          result.Object,
		"endpoint_id":     result.EndpointID,
		"since":           result.Since,
		"until":           result.Until,
		"matched_events":  result.MatchedEvents,
		"replayed_events": result.ReplayedEvents,
		"skipped_events":  result.SkippedEvents,
		"attempt_count":   result.AttemptCount,
		"events":          result.Events,
		"data":            h.deliveryAttemptResponses(r, result.Attempts),
	})
}

func replaySimulatedAppFailure(p params) *webhooks.SimulatedAppFailure {
	status := p.int64("simulate_app_failure[status]")
	if status == 0 {
		status = p.int64("simulateAppFailure[status]")
	}
	if status == 0 {
		return nil
	}
	failFirst := p.int64Default("simulate_app_failure[fail_first_n_attempts]", 1)
	if failFirst == 1 {
		failFirst = p.int64Default("simulateAppFailure[fail_first_n_attempts]", failFirst)
	}
	if failFirst == 1 {
		failFirst = p.int64Default("simulate_app_failure[failFirstNAttempts]", failFirst)
	}
	return &webhooks.SimulatedAppFailure{
		Status:             int(status),
		FailFirstNAttempts: int(failFirst),
		Body:               p.first("simulate_app_failure[body]", "simulateAppFailure[body]"),
	}
}

func parseOptionalReplayTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}
	if value == "now" {
		return time.Now().UTC(), nil
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed.UTC(), nil
	}
	unix, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(unix, 0).UTC(), nil
}

func historicalReplayEventTypes(p params) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, key := range []string{"type", "types", "event_type", "event_types"} {
		for _, value := range p.list(key) {
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
			out = append(out, value)
		}
	}
	return out
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
		decoder := json.NewDecoder(r.Body)
		decoder.UseNumber()
		if err := decoder.Decode(&raw); err != nil {
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
			case json.Number:
				values[key] = v.String()
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
	values := firstValues(r.Form)
	if security.ContainsCardData(values) {
		return params{}, fmt.Errorf("%w: real card data is not accepted by Billtap", webhooks.ErrInvalidInput)
	}
	return params{values: values}, nil
}

func paramsFromValues(values url.Values) params {
	return params{values: firstValues(values)}
}

func firstValues(values url.Values) map[string]string {
	out := map[string]string{}
	for key, value := range values {
		if len(value) > 0 {
			out[key] = value[0]
		}
	}
	return out
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

func paymentIntentMetadata(p params) map[string]string {
	metadata := p.metadata()
	outcome := firstNonEmptyString(
		metadataValue(metadata, billing.MetadataPaymentIntentOutcome),
		p.first("billtap_outcome", "deferred_outcome", "payment_intent_outcome"),
		metadataValue(metadata, "billtap_outcome"),
	)
	if actionType := p.first("billtap_next_action_type", "next_action_type"); actionType != "" {
		if metadata == nil {
			metadata = map[string]string{}
		}
		metadata["billtap_next_action_type"] = actionType
	}
	if returnURL := p.string("return_url"); returnURL != "" {
		if metadata == nil {
			metadata = map[string]string{}
		}
		metadata["billtap_return_url"] = returnURL
	}
	if outcome == "" {
		return metadata
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata[billing.MetadataPaymentIntentOutcome] = outcome
	return metadata
}

func metadataValue(metadata map[string]string, key string) string {
	if metadata == nil {
		return ""
	}
	return strings.TrimSpace(metadata[key])
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

var capabilityParamPattern = regexp.MustCompile(`^capabilities\[([^\]]+)\]\[(requested|status)\]$`)

func accountCapabilities(p params) map[string]string {
	out := map[string]string{}
	for key, value := range p.values {
		matches := capabilityParamPattern.FindStringSubmatch(key)
		if len(matches) != 3 {
			continue
		}
		capability := matches[1]
		switch matches[2] {
		case "requested":
			if value == "true" || value == "1" {
				out[capability] = "active"
			} else {
				out[capability] = "inactive"
			}
		case "status":
			out[capability] = value
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

var personDataParamPattern = regexp.MustCompile(`^(first_name|last_name|email|phone|id_number|ssn_last_4|relationship\[(owner|director|executive|representative|title|percent_ownership)\]|dob\[(day|month|year)\]|address\[(line1|line2|city|state|postal_code|country)\])$`)

func personDataFromParams(p params) map[string]string {
	out := map[string]string{}
	for key, value := range p.values {
		if personDataParamPattern.MatchString(key) {
			out[key] = strings.TrimSpace(value)
		}
	}
	return out
}

var accountSessionComponentPattern = regexp.MustCompile(`^components\[([^\]]+)\]\[enabled\]$`)

func accountSessionComponents(p params) map[string]any {
	out := map[string]any{}
	for key, value := range p.values {
		matches := accountSessionComponentPattern.FindStringSubmatch(key)
		if len(matches) != 2 {
			continue
		}
		out[matches[1]] = map[string]any{"enabled": value == "true" || value == "1"}
	}
	if len(out) == 0 {
		return map[string]any{"account_onboarding": map[string]any{"enabled": true}}
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

func parseInt64(value string) int64 {
	out, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	return out
}

func last4(value string) string {
	cleaned := sanitizeID(value)
	if len(cleaned) <= 4 {
		if cleaned == "" {
			return "6789"
		}
		return cleaned
	}
	return cleaned[len(cleaned)-4:]
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
	defaultPaymentMethod := customer.Metadata[billing.MetadataDefaultPaymentMethod]
	discounts := billing.DiscountsFromMetadata(customer.Metadata)
	return map[string]any{
		"id":         customer.ID,
		"object":     billing.ObjectCustomer,
		"email":      customer.Email,
		"name":       customer.Name,
		"metadata":   nonNilMap(customer.Metadata),
		"created":    unix(customer.CreatedAt),
		"created_at": customer.CreatedAt,
		"invoice_settings": map[string]any{
			"default_payment_method": emptyToNil(defaultPaymentMethod),
		},
		"test_clock": emptyToNil(customer.Metadata["test_clock"]),
		"discount":   firstDiscountObject(nil, discounts, customer.ID, "", ""),
		"livemode":   false,
	}
}

func stripePaymentMethod(customerID string, paymentMethodID string) map[string]any {
	if strings.TrimSpace(paymentMethodID) == "" {
		paymentMethodID = "pm_" + sanitizeID(customerID)
	}
	brand := paymentMethodBrand(paymentMethodID)
	return map[string]any{
		"id":              paymentMethodID,
		"object":          "payment_method",
		"allow_redisplay": "unspecified",
		"billing_details": map[string]any{
			"address": map[string]any{
				"city":        nil,
				"country":     nil,
				"line1":       nil,
				"line2":       nil,
				"postal_code": nil,
				"state":       nil,
			},
			"email": nil,
			"name":  nil,
			"phone": nil,
		},
		"customer": customerID,
		"type":     "card",
		"card": map[string]any{
			"brand":          brand,
			"checks":         map[string]any{"address_line1_check": nil, "address_postal_code_check": nil, "cvc_check": "pass"},
			"country":        "US",
			"exp_month":      12,
			"exp_year":       2035,
			"fingerprint":    paymentMethodFingerprint(paymentMethodID),
			"funding":        paymentMethodFunding(paymentMethodID),
			"generated_from": nil,
			"last4":          paymentMethodLast4(paymentMethodID),
			"networks":       map[string]any{"available": []string{brand}, "preferred": nil},
			"three_d_secure_usage": map[string]any{
				"supported": true,
			},
			"wallet": nil,
		},
		"created":   time.Now().UTC().Unix(),
		"livemode":  false,
		"metadata":  map[string]string{},
		"redaction": nil,
	}
}

func stripePaymentMethods(customer billing.Customer) []map[string]any {
	ids := customerPaymentMethodIDs(customer)
	data := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		data = append(data, stripePaymentMethod(customer.ID, id))
	}
	return data
}

func customerPaymentMethodIDs(customer billing.Customer) []string {
	defaultPaymentMethod := strings.TrimSpace(customer.Metadata[billing.MetadataDefaultPaymentMethod])
	explicitIDs := splitPaymentMethodIDs(customer.Metadata[billing.MetadataPaymentMethodIDs])
	autoPaymentMethod := "pm_" + sanitizeID(customer.ID)

	switch strings.ToLower(strings.TrimSpace(customer.Metadata[billing.MetadataPaymentMethodsFixture])) {
	case billing.PaymentMethodsFixtureEmpty:
		if defaultPaymentMethod == "" {
			return explicitIDs
		}
		return uniquePaymentMethodIDs([]string{autoPaymentMethod, defaultPaymentMethod})
	case billing.PaymentMethodsFixtureExplicit:
		ids := append([]string{}, explicitIDs...)
		if defaultPaymentMethod != "" {
			ids = append(ids, defaultPaymentMethod)
		}
		return uniquePaymentMethodIDs(ids)
	default:
		if defaultPaymentMethod != "" {
			return []string{defaultPaymentMethod}
		}
		return []string{autoPaymentMethod}
	}
}

func splitPaymentMethodIDs(raw string) []string {
	var ids []string
	for _, value := range strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n'
	}) {
		value = strings.TrimSpace(value)
		if value != "" {
			ids = append(ids, value)
		}
	}
	return uniquePaymentMethodIDs(ids)
}

func uniquePaymentMethodIDs(ids []string) []string {
	out := make([]string, 0, len(ids))
	seen := map[string]bool{}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func paymentMethodLast4(paymentMethodID string) string {
	switch strings.ToLower(strings.TrimSpace(paymentMethodID)) {
	case "pm_card_visa", "pm_card_visa_debit":
		return "4242"
	case "pm_card_three_secure", "pm_card_threesecure2required", "pm_card_authenticationrequired":
		return "3220"
	case "pm_card_chargecustomerfail", "pm_card_charge_declined", "pm_card_chargecustomerfail_attach":
		return "0002"
	default:
		return "4242"
	}
}

func paymentMethodBrand(paymentMethodID string) string {
	id := strings.ToLower(strings.TrimSpace(paymentMethodID))
	switch {
	case strings.Contains(id, "mastercard"):
		return "mastercard"
	case strings.Contains(id, "amex"):
		return "amex"
	case strings.Contains(id, "discover"):
		return "discover"
	default:
		return "visa"
	}
}

func paymentMethodFunding(paymentMethodID string) string {
	if strings.Contains(strings.ToLower(strings.TrimSpace(paymentMethodID)), "debit") {
		return "debit"
	}
	return "credit"
}

func paymentMethodFingerprint(paymentMethodID string) string {
	value := sanitizeID(paymentMethodID)
	if len(value) > 16 {
		value = value[:16]
	}
	return "bt_" + value
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
	var recurring any
	if price.RecurringInterval != "" {
		recurring = map[string]any{
			"interval":       price.RecurringInterval,
			"interval_count": price.RecurringIntervalCount,
			"usage_type":     "licensed",
		}
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
		"lookup_key":               priceLookupKey(price),
		"metadata":                 nonNilMap(price.Metadata),
		"product":                  price.ProductID,
		"recurring":                recurring,
		"recurring_interval":       price.RecurringInterval,
		"recurring_interval_count": price.RecurringIntervalCount,
		"tax_behavior":             "unspecified",
		"type":                     priceSearchType(price),
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

func stripeAccount(account billing.Account) map[string]any {
	return map[string]any{
		"id":                  account.ID,
		"object":              billing.ObjectAccount,
		"type":                account.Type,
		"country":             strings.ToUpper(account.Country),
		"email":               emptyToNil(account.Email),
		"business_type":       emptyToNil(account.BusinessType),
		"charges_enabled":     account.ChargesEnabled,
		"payouts_enabled":     account.PayoutsEnabled,
		"details_submitted":   account.DetailsSubmitted,
		"default_currency":    strings.ToLower(account.DefaultCurrency),
		"capabilities":        nonNilMap(account.Capabilities),
		"metadata":            nonNilMap(account.Metadata),
		"requirements":        accountRequirements(account),
		"future_requirements": accountRequirements(account),
		"settings": map[string]any{
			"dashboard": map[string]any{"display_name": nil, "timezone": "Etc/UTC"},
			"payouts":   map[string]any{"schedule": map[string]any{"interval": "manual"}},
		},
		"created":  unix(account.CreatedAt),
		"livemode": false,
	}
}

func stripeAccounts(accounts []billing.Account) []map[string]any {
	out := make([]map[string]any, 0, len(accounts))
	for _, account := range accounts {
		out = append(out, stripeAccount(account))
	}
	return out
}

func accountRequirements(account billing.Account) map[string]any {
	currentlyDue := []string{}
	if !account.DetailsSubmitted {
		currentlyDue = append(currentlyDue, "business_profile.url", "external_account")
	}
	requirements := emptyRequirements()
	requirements["currently_due"] = currentlyDue
	requirements["eventually_due"] = currentlyDue
	return requirements
}

func emptyRequirements() map[string]any {
	return map[string]any{
		"alternatives":         []map[string]any{},
		"currently_due":        []string{},
		"eventually_due":       []string{},
		"past_due":             []string{},
		"pending_verification": []string{},
		"disabled_reason":      nil,
		"current_deadline":     nil,
		"errors":               []map[string]any{},
	}
}

func stripeCapability(accountID string, capability string, status string) map[string]any {
	if status == "" {
		status = "inactive"
	}
	return map[string]any{
		"id":        capability,
		"object":    billing.ObjectCapability,
		"account":   accountID,
		"requested": status != "inactive",
		"status":    status,
		"livemode":  false,
		"requirements": map[string]any{
			"alternatives":         []map[string]any{},
			"currently_due":        []string{},
			"eventually_due":       []string{},
			"past_due":             []string{},
			"pending_verification": []string{},
		},
	}
}

func stripeConnectResources(resources []billing.ConnectResource) []map[string]any {
	out := make([]map[string]any, 0, len(resources))
	for _, resource := range resources {
		out = append(out, stripeConnectResource(resource))
	}
	return out
}

func stripeConnectResource(resource billing.ConnectResource) map[string]any {
	switch resource.Object {
	case billing.ObjectBankAccount, billing.ObjectCard:
		return stripeExternalAccount(resource)
	case billing.ObjectPerson:
		return stripePerson(resource)
	case billing.ObjectTransfer:
		return stripeTransfer(resource)
	case billing.ObjectTransferReversal:
		return stripeTransferReversal(resource)
	case billing.ObjectPayout:
		return stripePayout(resource)
	case billing.ObjectApplicationFee:
		return stripeApplicationFee(resource, nil)
	case billing.ObjectFeeRefund:
		return stripeApplicationFeeRefund(resource)
	default:
		return map[string]any{
			"id":       resource.ID,
			"object":   resource.Object,
			"metadata": nonNilMap(resource.Metadata),
			"livemode": false,
		}
	}
}

func stripePerson(resource billing.ConnectResource) map[string]any {
	return map[string]any{
		"id":         resource.ID,
		"object":     billing.ObjectPerson,
		"account":    resource.AccountID,
		"first_name": emptyToNil(resource.Data["first_name"]),
		"last_name":  emptyToNil(resource.Data["last_name"]),
		"email":      emptyToNil(resource.Data["email"]),
		"phone":      emptyToNil(resource.Data["phone"]),
		"created":    unix(resource.CreatedAt),
		"dob": map[string]any{
			"day":   optionalIntString(resource.Data["dob[day]"]),
			"month": optionalIntString(resource.Data["dob[month]"]),
			"year":  optionalIntString(resource.Data["dob[year]"]),
		},
		"address": map[string]any{
			"line1":       emptyToNil(resource.Data["address[line1]"]),
			"line2":       emptyToNil(resource.Data["address[line2]"]),
			"city":        emptyToNil(resource.Data["address[city]"]),
			"state":       emptyToNil(resource.Data["address[state]"]),
			"postal_code": emptyToNil(resource.Data["address[postal_code]"]),
			"country":     emptyToNil(resource.Data["address[country]"]),
		},
		"relationship": map[string]any{
			"owner":             truthy(resource.Data["relationship[owner]"]),
			"director":          truthy(resource.Data["relationship[director]"]),
			"executive":         truthy(resource.Data["relationship[executive]"]),
			"representative":    truthy(resource.Data["relationship[representative]"]),
			"title":             emptyToNil(resource.Data["relationship[title]"]),
			"percent_ownership": optionalIntString(resource.Data["relationship[percent_ownership]"]),
		},
		"requirements":        emptyRequirements(),
		"future_requirements": emptyRequirements(),
		"verification": map[string]any{
			"status":              "unverified",
			"document":            nil,
			"additional_document": nil,
			"details":             nil,
			"details_code":        nil,
		},
		"metadata": nonNilMap(resource.Metadata),
		"livemode": false,
	}
}

func stripeExternalAccount(resource billing.ConnectResource) map[string]any {
	return map[string]any{
		"id":                   resource.ID,
		"object":               stringDefault(resource.Object, billing.ObjectBankAccount),
		"account":              resource.AccountID,
		"account_holder_name":  emptyToNil(resource.Data["account_holder_name"]),
		"account_holder_type":  emptyToNil(resource.Data["account_holder_type"]),
		"bank_name":            stringDefault(resource.BankName, "Billtap Bank"),
		"country":              strings.ToUpper(stringDefault(resource.Country, "US")),
		"currency":             strings.ToLower(stringDefault(resource.Currency, "usd")),
		"default_for_currency": resource.Data["default_for_currency"] == "true",
		"fingerprint":          "bt_" + sanitizeID(resource.ID),
		"last4":                stringDefault(resource.Last4, "6789"),
		"metadata":             nonNilMap(resource.Metadata),
		"routing_number":       emptyToNil(resource.RoutingNumber),
		"status":               stringDefault(resource.Status, "new"),
		"created":              unix(resource.CreatedAt),
		"livemode":             false,
	}
}

func stripeTransfer(resource billing.ConnectResource) map[string]any {
	amountReversed := parseInt64(resource.Data["amount_reversed"])
	return map[string]any{
		"id":                  resource.ID,
		"object":              billing.ObjectTransfer,
		"amount":              resource.Amount,
		"amount_reversed":     amountReversed,
		"balance_transaction": "txn_" + resource.ID,
		"created":             unix(resource.CreatedAt),
		"currency":            strings.ToLower(resource.Currency),
		"description":         emptyToNil(resource.Description),
		"destination":         resource.Destination,
		"destination_payment": "py_" + resource.ID,
		"livemode":            false,
		"metadata":            nonNilMap(resource.Metadata),
		"reversed":            resource.Status == "reversed" || amountReversed >= resource.Amount,
		"source_transaction":  emptyToNil(resource.SourceTransaction),
		"source_type":         "card",
		"transfer_group":      emptyToNil(resource.Data["transfer_group"]),
	}
}

func stripeTransferReversal(resource billing.ConnectResource) map[string]any {
	return map[string]any{
		"id":                         resource.ID,
		"object":                     billing.ObjectTransferReversal,
		"amount":                     resource.Amount,
		"balance_transaction":        "txn_" + resource.ID,
		"created":                    unix(resource.CreatedAt),
		"currency":                   strings.ToLower(resource.Currency),
		"destination_payment_refund": nil,
		"metadata":                   nonNilMap(resource.Metadata),
		"source_refund":              nil,
		"transfer":                   resource.ParentID,
	}
}

func stripePayout(resource billing.ConnectResource) map[string]any {
	return map[string]any{
		"id":                          resource.ID,
		"object":                      billing.ObjectPayout,
		"amount":                      resource.Amount,
		"arrival_date":                unix(resource.ArrivalDate),
		"automatic":                   false,
		"balance_transaction":         "txn_" + resource.ID,
		"created":                     unix(resource.CreatedAt),
		"currency":                    strings.ToLower(resource.Currency),
		"description":                 emptyToNil(resource.Description),
		"destination":                 emptyToNil(resource.Destination),
		"failure_balance_transaction": nil,
		"failure_code":                nil,
		"failure_message":             nil,
		"livemode":                    false,
		"metadata":                    nonNilMap(resource.Metadata),
		"method":                      stringDefault(resource.Data["method"], "standard"),
		"reconciliation_status":       "not_applicable",
		"source_type":                 "card",
		"statement_descriptor":        emptyToNil(resource.Data["statement_descriptor"]),
		"status":                      stringDefault(resource.Status, "paid"),
		"type":                        "bank_account",
	}
}

func stripeApplicationFee(resource billing.ConnectResource, refunds []billing.ConnectResource) map[string]any {
	refundData := stripeConnectResources(refunds)
	amountRefunded := parseInt64(resource.Data["amount_refunded"])
	return map[string]any{
		"id":                      resource.ID,
		"object":                  billing.ObjectApplicationFee,
		"account":                 emptyToNil(resource.AccountID),
		"amount":                  resource.Amount,
		"amount_refunded":         amountRefunded,
		"application":             "ca_billtap",
		"balance_transaction":     "txn_" + resource.ID,
		"charge":                  emptyToNil(resource.Data["charge"]),
		"created":                 unix(resource.CreatedAt),
		"currency":                strings.ToLower(stringDefault(resource.Currency, "usd")),
		"livemode":                false,
		"originating_transaction": emptyToNil(resource.ParentID),
		"refunded":                amountRefunded >= resource.Amount && resource.Amount > 0,
		"refunds":                 stripeList("/v1/application_fees/"+resource.ID+"/refunds", refundData),
	}
}

func stripeApplicationFeeRefund(resource billing.ConnectResource) map[string]any {
	return map[string]any{
		"id":                  resource.ID,
		"object":              billing.ObjectFeeRefund,
		"amount":              resource.Amount,
		"balance_transaction": "txn_" + resource.ID,
		"created":             unix(resource.CreatedAt),
		"currency":            strings.ToLower(stringDefault(resource.Currency, "usd")),
		"fee":                 resource.ParentID,
		"metadata":            nonNilMap(resource.Metadata),
	}
}

func stripeDeleted(id string, object string) map[string]any {
	return map[string]any{"id": id, "object": object, "deleted": true}
}

func stripeCheckoutSession(session billing.CheckoutSession) map[string]any {
	return map[string]any{
		"id":                    session.ID,
		"object":                billing.ObjectCheckoutSession,
		"customer":              session.CustomerID,
		"mode":                  session.Mode,
		"line_items":            nil,
		"success_url":           session.SuccessURL,
		"cancel_url":            session.CancelURL,
		"allow_promotion_codes": session.AllowPromotionCodes,
		"trial_period_days":     session.TrialPeriodDays,
		"url":                   session.URL,
		"status":                session.Status,
		"payment_status":        session.PaymentStatus,
		"subscription":          emptyToNil(session.SubscriptionID),
		"invoice":               emptyToNil(session.InvoiceID),
		"payment_intent":        emptyToNil(session.PaymentIntentID),
		"created":               unix(session.CreatedAt),
		"created_at":            session.CreatedAt,
		"completed_at":          session.CompletedAt,
		"livemode":              false,
	}
}

func (h *Handler) stripeSubscription(r *http.Request, sub billing.Subscription) map[string]any {
	items := make([]map[string]any, 0, len(sub.Items))
	for idx, item := range sub.Items {
		items = append(items, h.stripeSubscriptionItem(r, sub, item, idx))
	}
	discounts := billing.DiscountsFromMetadata(sub.Metadata)
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
		"cancel_at":            subscriptionCancelAt(sub),
		"ended_at":             metadataUnix(sub.Metadata["ended_at"]),
		"trial_start":          metadataUnix(sub.Metadata["trial_start"]),
		"trial_end":            metadataUnix(sub.Metadata["trial_end"]),
		"latest_invoice":       emptyToNil(sub.LatestInvoiceID),
		"discount":             firstDiscountObject(h, discounts, sub.CustomerID, sub.ID, ""),
		"discounts":            stripeList("/v1/subscriptions/"+sub.ID+"/discounts", discountObjects(h, discounts, sub.CustomerID, sub.ID, "")),
		"test_clock":           emptyToNil(sub.Metadata["test_clock"]),
		"metadata":             nonNilMap(sub.Metadata),
		"collection_method":    stringDefault(sub.Metadata["collection_method"], "charge_automatically"),
		"billing_cycle_anchor": unix(sub.CurrentPeriodStart),
		"currency":             "usd",
		"livemode":             false,
		"pause_collection":     nil,
		"pending_update":       nil,
		"cancellation_details": subscriptionCancellationDetails(sub),
	}
}

func subscriptionCancelAt(sub billing.Subscription) any {
	if !sub.CancelAtPeriodEnd {
		return nil
	}
	if value := metadataUnix(sub.Metadata["cancel_at"]); value != nil {
		return value
	}
	return unix(sub.CurrentPeriodEnd)
}

func subscriptionCancellationDetails(sub billing.Subscription) map[string]any {
	return map[string]any{
		"comment":  emptyToNil(sub.Metadata["cancellation_details_comment"]),
		"feedback": emptyToNil(sub.Metadata["cancellation_details_feedback"]),
		"reason":   nil,
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
		"id":                   subscriptionItemID(sub, idx),
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

func subscriptionItemID(sub billing.Subscription, idx int) string {
	return "si_" + sanitizeID(sub.ID+"_"+strconv.Itoa(idx))
}

func stripeInvoice(invoice billing.Invoice) map[string]any {
	return map[string]any{
		"id":                     invoice.ID,
		"object":                 billing.ObjectInvoice,
		"customer":               invoice.CustomerID,
		"subscription":           emptyToNil(invoice.SubscriptionID),
		"parent":                 map[string]any{"subscription_details": map[string]any{"subscription": emptyToNil(invoice.SubscriptionID)}},
		"status":                 invoice.Status,
		"currency":               invoice.Currency,
		"subtotal":               invoice.Subtotal,
		"discount":               firstDiscountObject(nil, invoice.Discounts, invoice.CustomerID, invoice.SubscriptionID, invoice.ID),
		"discounts":              stripeList("/v1/invoices/"+invoice.ID+"/discounts", discountObjects(nil, invoice.Discounts, invoice.CustomerID, invoice.SubscriptionID, invoice.ID)),
		"total_discount_amounts": discountAmounts(invoice.Discounts, invoice.DiscountAmount),
		"total":                  invoice.Total,
		"amount_due":             invoice.AmountDue,
		"amount_paid":            invoice.AmountPaid,
		"attempt_count":          invoice.AttemptCount,
		"next_payment_attempt":   optionalUnix(invoice.NextPaymentAttempt),
		"payment_intent":         emptyToNil(invoice.PaymentIntentID),
		"payments":               stripeList("/v1/invoices/"+invoice.ID+"/payments", []map[string]any{}),
		"lines":                  stripeList("/v1/invoices/"+invoice.ID+"/lines", []map[string]any{}),
		"created":                unix(invoice.CreatedAt),
		"created_at":             invoice.CreatedAt,
		"status_transitions": map[string]any{
			"paid_at": optionalPaidAt(invoice),
		},
		"hosted_invoice_url": "",
		"livemode":           false,
	}
}

func stripePaymentIntent(intent billing.PaymentIntent) map[string]any {
	captureMethod := stringDefault(intent.CaptureMethod, "automatic")
	return map[string]any{
		"id":                 intent.ID,
		"object":             billing.ObjectPaymentIntent,
		"customer":           emptyToNil(intent.CustomerID),
		"invoice":            emptyToNil(intent.InvoiceID),
		"amount":             intent.Amount,
		"amount_capturable":  paymentIntentCapturableAmount(intent),
		"amount_received":    paymentIntentReceivedAmount(intent),
		"currency":           intent.Currency,
		"status":             intent.Status,
		"capture_method":     captureMethod,
		"payment_method":     emptyToNil(intent.PaymentMethodID),
		"metadata":           nonNilMap(intent.Metadata),
		"last_payment_error": paymentIntentError(intent),
		"next_action":        paymentIntentNextAction(intent),
		"client_secret":      intent.ID + "_secret_billtap",
		"created":            unix(intent.CreatedAt),
		"created_at":         intent.CreatedAt,
		"livemode":           false,
	}
}

func paymentIntentNextAction(intent billing.PaymentIntent) any {
	if intent.Status != "requires_action" {
		return nil
	}
	if strings.TrimSpace(intent.Metadata["billtap_next_action_type"]) == "redirect_to_url" {
		returnURL := firstNonEmptyString(intent.Metadata["billtap_return_url"], "http://127.0.0.1:18080/payment_intents/"+intent.ID+"/return")
		return map[string]any{
			"type": "redirect_to_url",
			"redirect_to_url": map[string]any{
				"url":        "/api/payment_intents/" + intent.ID + "/complete_action?return_url=" + url.QueryEscape(returnURL),
				"return_url": returnURL,
			},
		}
	}
	return map[string]any{
		"type": "use_stripe_sdk",
		"use_stripe_sdk": map[string]any{
			"type":                    "three_d_secure_redirect",
			"stripe_js":               "https://js.stripe.com/v3/",
			"source":                  intent.PaymentMethodID,
			"server_transaction_id":   "billtap_" + sanitizeID(intent.ID),
			"hosted_voucher_url":      "",
			"billtap_completion_url":  "/api/payment_intents/" + intent.ID + "/complete_action",
			"billtap_abandonment_url": "/api/payment_intents/" + intent.ID + "/cancel_action",
		},
	}
}

func stripeSetupIntent(intent billing.SetupIntent) map[string]any {
	return map[string]any{
		"id":               intent.ID,
		"object":           billing.ObjectSetupIntent,
		"customer":         emptyToNil(intent.CustomerID),
		"status":           intent.Status,
		"usage":            stringDefault(intent.Usage, "off_session"),
		"payment_method":   emptyToNil(intent.PaymentMethodID),
		"last_setup_error": setupIntentError(intent),
		"client_secret":    intent.ID + "_secret_billtap",
		"created":          unix(intent.CreatedAt),
		"created_at":       intent.CreatedAt,
		"livemode":         false,
	}
}

func stripeTestClock(clock billing.TestClock) map[string]any {
	return map[string]any{
		"id":          clock.ID,
		"object":      billing.ObjectTestClock,
		"name":        clock.Name,
		"status":      clock.Status,
		"frozen_time": unix(clock.FrozenTime),
		"created":     unix(clock.CreatedAt),
		"livemode":    false,
	}
}

func stripeRefund(refund billing.Refund) map[string]any {
	return map[string]any{
		"id":             refund.ID,
		"object":         billing.ObjectRefund,
		"amount":         refund.Amount,
		"currency":       refund.Currency,
		"charge":         emptyToNil(refund.ChargeID),
		"payment_intent": emptyToNil(refund.PaymentIntentID),
		"invoice":        emptyToNil(refund.InvoiceID),
		"customer":       emptyToNil(refund.CustomerID),
		"reason":         emptyToNil(refund.Reason),
		"status":         refund.Status,
		"metadata":       nonNilMap(refund.Metadata),
		"created":        unix(refund.CreatedAt),
		"livemode":       false,
	}
}

func stripeChargeFromRefund(refund billing.Refund) map[string]any {
	refunded := refund.Amount
	return map[string]any{
		"id":              refund.ChargeID,
		"object":          "charge",
		"amount":          refund.Amount,
		"amount_refunded": refunded,
		"currency":        refund.Currency,
		"customer":        emptyToNil(refund.CustomerID),
		"invoice":         emptyToNil(refund.InvoiceID),
		"payment_intent":  emptyToNil(refund.PaymentIntentID),
		"paid":            true,
		"refunded":        true,
		"refunds": stripeList("/v1/refunds?charge="+url.QueryEscape(refund.ChargeID), []map[string]any{
			stripeRefund(refund),
		}),
		"metadata": nonNilMap(refund.Metadata),
		"created":  unix(refund.CreatedAt),
		"livemode": false,
	}
}

func stripeCreditNote(note billing.CreditNote) map[string]any {
	return map[string]any{
		"id":       note.ID,
		"object":   billing.ObjectCreditNote,
		"invoice":  note.InvoiceID,
		"customer": emptyToNil(note.CustomerID),
		"amount":   note.Amount,
		"currency": note.Currency,
		"reason":   emptyToNil(note.Reason),
		"status":   note.Status,
		"metadata": nonNilMap(note.Metadata),
		"created":  unix(note.CreatedAt),
		"livemode": false,
		"lines":    stripeList("/v1/credit_notes/"+note.ID+"/lines", []map[string]any{}),
	}
}

func paymentIntentCapturableAmount(intent billing.PaymentIntent) int64 {
	if intent.Status == "requires_capture" {
		return intent.Amount
	}
	return 0
}

func paymentIntentReceivedAmount(intent billing.PaymentIntent) int64 {
	if intent.Status == "succeeded" {
		return intent.Amount
	}
	return 0
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
		if priceType := query.Get("type"); priceType != "" && priceSearchType(price) != priceType {
			continue
		}
		out = append(out, price)
	}
	return out
}

func filterAccounts(accounts []billing.Account, r *http.Request) []billing.Account {
	query := r.URL.Query()
	metadataFilters := queryMetadataFilters(query)
	out := make([]billing.Account, 0, len(accounts))
	for _, account := range accounts {
		if country := query.Get("country"); country != "" && !strings.EqualFold(account.Country, country) {
			continue
		}
		if accountType := query.Get("type"); accountType != "" && account.Type != accountType {
			continue
		}
		if !metadataMatches(account.Metadata, metadataFilters) {
			continue
		}
		out = append(out, account)
	}
	return out
}

func filterSubscriptions(items []billing.Subscription, r *http.Request) []billing.Subscription {
	query := r.URL.Query()
	metadataFilters := queryMetadataFilters(query)
	out := make([]billing.Subscription, 0, len(items))
	for _, item := range items {
		if customer := query.Get("customer"); customer != "" && item.CustomerID != customer {
			continue
		}
		status := strings.ToLower(query.Get("status"))
		if status != "" && status != "all" && item.Status != status {
			continue
		}
		if !metadataMatches(item.Metadata, metadataFilters) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func queryMetadataFilters(query url.Values) map[string]string {
	out := map[string]string{}
	for key, values := range query {
		matches := metadataParamPattern.FindStringSubmatch(key)
		if len(matches) != 2 || len(values) == 0 {
			continue
		}
		out[matches[1]] = values[0]
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func metadataMatches(metadata map[string]string, filters map[string]string) bool {
	for key, value := range filters {
		if metadata[key] != value {
			return false
		}
	}
	return true
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

type priceSearchCriteria struct {
	Active    *bool
	Type      string
	LookupKey string
	Metadata  map[string]string
}

var (
	priceSearchActiveClause   = regexp.MustCompile(`^active:'(true|false)'$`)
	priceSearchTypeClause     = regexp.MustCompile(`^type:'(one_time|recurring)'$`)
	priceSearchLookupClause   = regexp.MustCompile(`^lookup_key:'([^']*)'$`)
	priceSearchMetadataClause = regexp.MustCompile(`^metadata\['([^']+)'\]:'([^']*)'$`)
	priceSearchANDPattern     = regexp.MustCompile(`(?i)\s+AND\s+`)
)

func parsePriceSearchQuery(query string) (priceSearchCriteria, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return priceSearchCriteria{}, missingParam("query")
	}
	criteria := priceSearchCriteria{Metadata: map[string]string{}}
	for _, rawClause := range priceSearchANDPattern.Split(query, -1) {
		clause := strings.TrimSpace(rawClause)
		if clause == "" {
			return priceSearchCriteria{}, invalidParam("query", "Expected field:'value' clauses joined by AND.")
		}
		switch {
		case priceSearchActiveClause.MatchString(clause):
			match := priceSearchActiveClause.FindStringSubmatch(clause)
			active := match[1] == "true"
			criteria.Active = &active
		case priceSearchTypeClause.MatchString(clause):
			match := priceSearchTypeClause.FindStringSubmatch(clause)
			criteria.Type = match[1]
		case priceSearchLookupClause.MatchString(clause):
			match := priceSearchLookupClause.FindStringSubmatch(clause)
			criteria.LookupKey = match[1]
		case priceSearchMetadataClause.MatchString(clause):
			match := priceSearchMetadataClause.FindStringSubmatch(clause)
			criteria.Metadata[match[1]] = match[2]
		default:
			return priceSearchCriteria{}, invalidParam("query", "Unsupported prices search clause: "+clause+".")
		}
	}
	if len(criteria.Metadata) == 0 {
		criteria.Metadata = nil
	}
	return criteria, nil
}

func filterPriceSearchResults(prices []billing.Price, criteria priceSearchCriteria) []billing.Price {
	out := make([]billing.Price, 0, len(prices))
	for _, price := range prices {
		if criteria.Active != nil && price.Active != *criteria.Active {
			continue
		}
		if criteria.Type != "" && priceSearchType(price) != criteria.Type {
			continue
		}
		if criteria.LookupKey != "" && priceLookupKey(price) != criteria.LookupKey {
			continue
		}
		if !metadataMatches(price.Metadata, criteria.Metadata) {
			continue
		}
		out = append(out, price)
	}
	return out
}

func priceSearchType(price billing.Price) string {
	if price.RecurringInterval != "" {
		return "recurring"
	}
	return "one_time"
}

func priceLookupKey(price billing.Price) string {
	if price.LookupKey != "" {
		return price.LookupKey
	}
	return price.Metadata["lookup_key"]
}

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
	out := append([]billing.LineItem{}, current.Items...)
	for i := 0; i < 100; i++ {
		itemID := p.string(fmt.Sprintf("items[%d][id]", i))
		priceID := p.first(fmt.Sprintf("items[%d][price]", i), fmt.Sprintf("items[%d][price_id]", i))
		quantity := p.int64(fmt.Sprintf("items[%d][quantity]", i))
		if itemID != "" {
			idx := subscriptionItemIndexByID(current, itemID)
			if idx < 0 {
				continue
			}
			if priceID == "" {
				priceID = out[idx].PriceID
			}
			if quantity <= 0 {
				quantity = out[idx].Quantity
			}
			out[idx] = billing.LineItem{PriceID: priceID, Quantity: quantity}
			continue
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

func subscriptionItemIndexByID(sub billing.Subscription, itemID string) int {
	for idx := range sub.Items {
		if subscriptionItemID(sub, idx) == itemID {
			return idx
		}
	}
	return -1
}

func subscriptionCreateItemsFromParams(p params) []billing.LineItem {
	var out []billing.LineItem
	for i := 0; i < 100; i++ {
		priceID := p.first(fmt.Sprintf("items[%d][price]", i), fmt.Sprintf("items[%d][price_id]", i))
		if priceID == "" {
			if i == 0 {
				priceID = p.first("price", "price_id")
			}
			if priceID == "" {
				continue
			}
		}
		quantity := p.int64Default(fmt.Sprintf("items[%d][quantity]", i), 1)
		if quantity <= 0 {
			quantity = 1
		}
		out = append(out, billing.LineItem{PriceID: priceID, Quantity: quantity})
	}
	return out
}

func subscriptionUpdateMetadata(p params) map[string]string {
	metadata := p.metadata()
	for _, item := range []struct {
		param string
		key   string
	}{
		{param: "cancellation_details[comment]", key: "cancellation_details_comment"},
		{param: "cancellation_details[feedback]", key: "cancellation_details_feedback"},
	} {
		if value := p.string(item.param); value != "" {
			if metadata == nil {
				metadata = map[string]string{}
			}
			metadata[item.key] = value
		}
	}
	return metadata
}

func (h *Handler) findSubscriptionItem(r *http.Request, itemID string) (billing.Subscription, int, bool, error) {
	subscriptions, err := h.billing.ListSubscriptions(r.Context())
	if err != nil {
		return billing.Subscription{}, 0, false, err
	}
	for _, subscription := range subscriptions {
		for idx := range subscription.Items {
			if subscriptionItemID(subscription, idx) == itemID {
				return subscription, idx, true, nil
			}
		}
	}
	return billing.Subscription{}, 0, false, nil
}

func unix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}

func parseTimestampParam(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("%w: frozen_time is required", billing.ErrInvalidInput)
	}
	if value == "now" {
		return time.Now().UTC(), nil
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(seconds, 0).UTC(), nil
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("%w: invalid timestamp %q", billing.ErrInvalidInput, value)
}

func optionalUnix(t *time.Time) any {
	if t == nil || t.IsZero() {
		return nil
	}
	return t.Unix()
}

func metadataUnix(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil || parsed.IsZero() {
		return nil
	}
	return parsed.Unix()
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
	out := map[string]any{
		"type":    stripeErrorCard,
		"code":    intent.FailureCode,
		"message": intent.FailureMessage,
	}
	if intent.DeclineCode != "" {
		out["decline_code"] = intent.DeclineCode
	}
	return out
}

func setupIntentError(intent billing.SetupIntent) any {
	if intent.FailureCode == "" && intent.FailureMessage == "" {
		return nil
	}
	out := map[string]any{
		"type":    stripeErrorCard,
		"code":    intent.FailureCode,
		"message": intent.FailureMessage,
	}
	if intent.DeclineCode != "" {
		out["decline_code"] = intent.DeclineCode
	}
	return out
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

func optionalIntString(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return parseInt64(value)
}

func truthy(value string) bool {
	return value == "true" || value == "1"
}

func stringDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
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
		var validationErr *validationError
		if errors.As(err, &validationErr) {
			writeError(w, http.StatusBadRequest, err)
			return
		}
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

func (h *Handler) absoluteURL(r *http.Request, path string) string {
	return absoluteURL(r, path, h.publicBase)
}

func absoluteURL(r *http.Request, path string, publicBase string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if publicBase != "" {
		if path == "" {
			return publicBase
		}
		if strings.HasPrefix(path, "/") {
			return publicBase + path
		}
		return publicBase + "/" + path
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
	for idx, item := range h.checkoutWebhookPayloads(r, result) {
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

func (h *Handler) emitSubscriptionWebhook(r *http.Request, eventType string, subscription billing.Subscription, source string) []webhooks.Event {
	if h.webhooks == nil || subscription.ID == "" {
		return nil
	}
	raw, err := json.Marshal(h.stripeSubscription(r, subscription))
	if err != nil {
		return nil
	}
	sequence := time.Now().UTC().UnixNano()
	event, _, err := h.webhooks.CreateEvent(r.Context(), webhooks.EventInput{
		Type:           eventType,
		ObjectPayload:  raw,
		RequestID:      "req_" + subscription.ID,
		IdempotencyKey: fmt.Sprintf("billtap:%s:%s:%d", eventType, subscription.ID, sequence),
		Source:         source,
		Sequence:       sequence,
	})
	if err != nil {
		return nil
	}
	return []webhooks.Event{event}
}

func (h *Handler) emitPaymentIntentWebhook(r *http.Request, eventType string, intent billing.PaymentIntent) []webhooks.Event {
	if h.webhooks == nil || intent.ID == "" || eventType == "" {
		return nil
	}
	raw, err := json.Marshal(stripePaymentIntent(intent))
	if err != nil {
		return nil
	}
	sequence := time.Now().UTC().UnixNano()
	event, _, err := h.webhooks.CreateEvent(r.Context(), webhooks.EventInput{
		Type:           eventType,
		ObjectPayload:  raw,
		RequestID:      "req_" + intent.ID,
		IdempotencyKey: fmt.Sprintf("billtap:%s:%s:%d", eventType, intent.ID, sequence),
		Source:         webhooks.SourceAPI,
		Sequence:       sequence,
	})
	if err != nil {
		return nil
	}
	return []webhooks.Event{event}
}

func (h *Handler) emitInvoicePaymentWebhooks(r *http.Request, result billing.InvoicePaymentResult, source string) []webhooks.Event {
	if h.webhooks == nil || result.Invoice.ID == "" {
		return nil
	}
	sequence := time.Now().UTC().UnixNano()
	var emitted []webhooks.Event
	for idx, item := range h.invoicePaymentWebhookPayloads(r, result) {
		event, _, err := h.webhooks.CreateEvent(r.Context(), webhooks.EventInput{
			Type:           item.eventType,
			ObjectPayload:  item.payload,
			RequestID:      "req_" + item.objectID,
			IdempotencyKey: fmt.Sprintf("billtap:%s:%s:%d", item.eventType, item.objectID, sequence),
			Source:         source,
			Sequence:       sequence + int64(idx),
		})
		if err == nil {
			emitted = append(emitted, event)
		}
	}
	return emitted
}

func (h *Handler) emitRenewalWebhooks(r *http.Request, result billing.InvoicePaymentResult, source string) []webhooks.Event {
	if h.webhooks == nil {
		return nil
	}
	sequence := time.Now().UTC().UnixNano()
	var emitted []webhooks.Event
	for idx, item := range h.renewalWebhookPayloads(r, result) {
		event, _, err := h.webhooks.CreateEvent(r.Context(), webhooks.EventInput{
			Type:           item.eventType,
			ObjectPayload:  item.payload,
			RequestID:      "req_" + item.objectID,
			IdempotencyKey: fmt.Sprintf("billtap:%s:%s:%d", item.eventType, item.objectID, sequence),
			Source:         source,
			Sequence:       sequence + int64(idx),
		})
		if err == nil {
			emitted = append(emitted, event)
		}
	}
	return emitted
}

func (h *Handler) emitClockAdvanceWebhooks(r *http.Request, advance billing.ClockAdvanceResult) []webhooks.Event {
	if h.webhooks == nil {
		return nil
	}
	var emitted []webhooks.Event
	for _, subscription := range advance.Activated {
		emitted = append(emitted, h.emitSubscriptionWebhook(r, "customer.subscription.updated", subscription, webhooks.SourceAPI)...)
	}
	for _, renewal := range advance.Renewals {
		emitted = append(emitted, h.emitRenewalWebhooks(r, renewal, webhooks.SourceAPI)...)
	}
	for _, subscription := range advance.Canceled {
		emitted = append(emitted, h.emitSubscriptionWebhook(r, "customer.subscription.deleted", subscription, webhooks.SourceAPI)...)
	}
	for _, refund := range advance.SettledRefunds {
		emitted = append(emitted, h.emitGenericWebhook(r, "charge.refund.updated", refund.ID, stripeRefund(refund), webhooks.SourceAPI)...)
	}
	return emitted
}

func (h *Handler) emitFixtureApplyWebhooks(r *http.Request, result fixtures.ApplyResult) []webhooks.Event {
	if h.webhooks == nil {
		return nil
	}
	subscriptions := make(map[string]billing.Subscription, len(result.Subscriptions))
	for _, subscription := range result.Subscriptions {
		subscriptions[subscription.ID] = subscription
	}
	seenSubscriptions := map[string]bool{}
	var emitted []webhooks.Event
	for _, session := range result.CheckoutSessions {
		subscription := subscriptions[session.SubscriptionID]
		if subscription.ID != "" {
			seenSubscriptions[subscription.ID] = true
		}
		invoice := billing.Invoice{}
		if session.InvoiceID != "" {
			if found, err := h.billing.GetInvoice(r.Context(), session.InvoiceID); err == nil {
				invoice = found
			}
		}
		intent := billing.PaymentIntent{}
		if session.PaymentIntentID != "" {
			if found, err := h.billing.GetPaymentIntent(r.Context(), session.PaymentIntentID); err == nil {
				intent = found
			}
		}
		emitted = append(emitted, h.emitCheckoutWebhooks(r, map[string]any{
			"session":        session,
			"subscription":   subscription,
			"invoice":        invoice,
			"payment_intent": intent,
		})...)
	}
	for _, subscription := range result.Subscriptions {
		if seenSubscriptions[subscription.ID] {
			continue
		}
		emitted = append(emitted, h.ensureFixtureSubscriptionCreatedWebhook(r, subscription)...)
		emitted = append(emitted, h.emitSubscriptionWebhook(r, "customer.subscription.updated", subscription, webhooks.SourceAPI)...)
		if subscription.Status == "canceled" {
			emitted = append(emitted, h.emitSubscriptionWebhook(r, "customer.subscription.deleted", subscription, webhooks.SourceAPI)...)
		}
		if subscription.Status == "past_due" || subscription.Status == "unpaid" || subscription.Status == "incomplete" {
			invoice, intent := h.subscriptionInvoicePaymentEvidence(r, subscription)
			emitted = append(emitted, h.emitInvoicePaymentWebhooks(r, billing.InvoicePaymentResult{
				Invoice:       invoice,
				PaymentIntent: intent,
				Subscription:  subscription,
			}, webhooks.SourceAPI)...)
		}
	}
	for _, refund := range result.Refunds {
		emitted = append(emitted, h.emitRefundWebhooks(r, refund, webhooks.SourceAPI)...)
	}
	for _, note := range result.CreditNotes {
		emitted = append(emitted, h.emitGenericWebhook(r, "credit_note.created", note.ID, stripeCreditNote(note), webhooks.SourceAPI)...)
		if note.Status == "void" {
			emitted = append(emitted, h.emitGenericWebhook(r, "credit_note.voided", note.ID, stripeCreditNote(note), webhooks.SourceAPI)...)
		}
	}
	return emitted
}

func (h *Handler) ensureFixtureSubscriptionCreatedWebhook(r *http.Request, subscription billing.Subscription) []webhooks.Event {
	if h.webhooks == nil || subscription.ID == "" || h.webhookEventExistsForObject(r.Context(), "customer.subscription.created", subscription.ID) {
		return nil
	}
	return h.emitSubscriptionWebhook(r, "customer.subscription.created", subscription, webhooks.SourceFixture)
}

func (h *Handler) webhookEventExistsForObject(ctx context.Context, eventType string, objectID string) bool {
	if h.webhooks == nil || eventType == "" || objectID == "" {
		return false
	}
	events, err := h.webhooks.ListEvents(ctx, webhooks.EventFilter{Type: eventType})
	if err != nil {
		return false
	}
	for _, event := range events {
		var object struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(event.Data.Object, &object); err == nil && object.ID == objectID {
			return true
		}
	}
	return false
}

func (h *Handler) subscriptionInvoicePaymentEvidence(r *http.Request, subscription billing.Subscription) (billing.Invoice, billing.PaymentIntent) {
	invoice := billing.Invoice{}
	if subscription.LatestInvoiceID != "" {
		if found, err := h.billing.GetInvoice(r.Context(), subscription.LatestInvoiceID); err == nil {
			invoice = found
		}
	}
	intent := billing.PaymentIntent{}
	if invoice.PaymentIntentID != "" {
		if found, err := h.billing.GetPaymentIntent(r.Context(), invoice.PaymentIntentID); err == nil {
			intent = found
		}
	}
	return invoice, intent
}

func (h *Handler) emitRefundWebhooks(r *http.Request, refund billing.Refund, source string) []webhooks.Event {
	var emitted []webhooks.Event
	emitted = append(emitted, h.emitGenericWebhook(r, "charge.refunded", refund.ChargeID, stripeChargeFromRefund(refund), source)...)
	emitted = append(emitted, h.emitGenericWebhook(r, "charge.refund.updated", refund.ID, stripeRefund(refund), source)...)
	return emitted
}

func (h *Handler) emitGenericWebhook(r *http.Request, eventType string, objectID string, payload any, source string) []webhooks.Event {
	if h.webhooks == nil || eventType == "" || objectID == "" {
		return nil
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	sequence := time.Now().UTC().UnixNano()
	event, _, err := h.webhooks.CreateEvent(r.Context(), webhooks.EventInput{
		Type:           eventType,
		ObjectPayload:  raw,
		RequestID:      "req_" + objectID,
		IdempotencyKey: fmt.Sprintf("billtap:%s:%s:%d", eventType, objectID, sequence),
		Source:         source,
		Sequence:       sequence,
	})
	if err != nil {
		return nil
	}
	return []webhooks.Event{event}
}

func (h *Handler) emitSetupIntentWebhook(r *http.Request, eventType string, intent billing.SetupIntent) []webhooks.Event {
	if h.webhooks == nil || intent.ID == "" || eventType == "" {
		return nil
	}
	raw, err := json.Marshal(stripeSetupIntent(intent))
	if err != nil {
		return nil
	}
	sequence := time.Now().UTC().UnixNano()
	event, _, err := h.webhooks.CreateEvent(r.Context(), webhooks.EventInput{
		Type:           eventType,
		ObjectPayload:  raw,
		RequestID:      "req_" + intent.ID,
		IdempotencyKey: fmt.Sprintf("billtap:%s:%s:%d", eventType, intent.ID, sequence),
		Source:         webhooks.SourceAPI,
		Sequence:       sequence,
	})
	if err != nil {
		return nil
	}
	return []webhooks.Event{event}
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
		checkoutSessions[i].URL = h.absoluteURL(r, checkoutSessions[i].URL)
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
	setupIntents, err := h.billing.ListSetupIntents(r.Context())
	if err != nil {
		return nil, err
	}
	testClocks, err := h.billing.ListTestClocks(r.Context())
	if err != nil {
		return nil, err
	}
	refunds, err := h.billing.ListRefunds(r.Context(), billing.RefundFilter{})
	if err != nil {
		return nil, err
	}
	creditNotes, err := h.billing.ListCreditNotes(r.Context(), billing.CreditNoteFilter{})
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
		"setup_intents":     setupIntents,
		"test_clocks":       testClocks,
		"refunds":           refunds,
		"credit_notes":      creditNotes,
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

func (h *Handler) checkoutWebhookPayloads(r *http.Request, result map[string]any) []webhookPayload {
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
	appendPayload(checkoutSessionEvent(session.Status), session.ID, stripeCheckoutSession(session))
	if subscription.ID != "" {
		appendPayload("customer.subscription.created", subscription.ID, h.stripeSubscription(r, subscription))
		if discounts := billing.DiscountsFromMetadata(subscription.Metadata); len(discounts) > 0 {
			appendPayload("customer.discount.created", discounts[0].ID, h.stripeDiscount(discounts[0], subscription.CustomerID, subscription.ID, ""))
		}
	}
	if invoice.ID != "" {
		appendPayload("invoice.created", invoice.ID, stripeInvoice(invoice))
		appendPayload("invoice.finalized", invoice.ID, stripeInvoice(invoice))
	}
	if paymentIntent.ID != "" {
		appendPayload("payment_intent.created", paymentIntent.ID, stripePaymentIntent(paymentIntent))
		if eventType := paymentIntentTerminalEvent(paymentIntent.Status); eventType != "" {
			appendPayload(eventType, paymentIntent.ID, stripePaymentIntent(paymentIntent))
		}
	}
	if invoice.ID != "" {
		for _, eventType := range invoiceTerminalEvents(invoice.Status, paymentIntent.Status) {
			appendPayload(eventType, invoice.ID, stripeInvoice(invoice))
		}
	}
	if subscription.ID != "" {
		appendPayload("customer.subscription.updated", subscription.ID, h.stripeSubscription(r, subscription))
	}
	return out
}

func (h *Handler) invoicePaymentWebhookPayloads(r *http.Request, result billing.InvoicePaymentResult) []webhookPayload {
	var out []webhookPayload
	appendPayload := func(eventType string, objectID string, value any) {
		if eventType == "" || objectID == "" {
			return
		}
		raw, err := json.Marshal(value)
		if err == nil {
			out = append(out, webhookPayload{eventType: eventType, objectID: objectID, payload: raw})
		}
	}
	if result.PaymentIntent.ID != "" {
		appendPayload(paymentIntentTerminalEvent(result.PaymentIntent.Status), result.PaymentIntent.ID, stripePaymentIntent(result.PaymentIntent))
	}
	if result.Invoice.ID != "" {
		for _, eventType := range invoiceTerminalEvents(result.Invoice.Status, result.PaymentIntent.Status) {
			appendPayload(eventType, result.Invoice.ID, stripeInvoice(result.Invoice))
		}
	}
	if result.Subscription.ID != "" {
		appendPayload("customer.subscription.updated", result.Subscription.ID, h.stripeSubscription(r, result.Subscription))
	}
	return out
}

func (h *Handler) renewalWebhookPayloads(r *http.Request, result billing.InvoicePaymentResult) []webhookPayload {
	var out []webhookPayload
	appendPayload := func(eventType string, objectID string, value any) {
		if eventType == "" || objectID == "" {
			return
		}
		raw, err := json.Marshal(value)
		if err == nil {
			out = append(out, webhookPayload{eventType: eventType, objectID: objectID, payload: raw})
		}
	}
	if result.Invoice.ID != "" {
		appendPayload("invoice.created", result.Invoice.ID, stripeInvoice(result.Invoice))
		appendPayload("invoice.finalized", result.Invoice.ID, stripeInvoice(result.Invoice))
	}
	if result.PaymentIntent.ID != "" {
		appendPayload("payment_intent.created", result.PaymentIntent.ID, stripePaymentIntent(result.PaymentIntent))
		appendPayload(paymentIntentTerminalEvent(result.PaymentIntent.Status), result.PaymentIntent.ID, stripePaymentIntent(result.PaymentIntent))
	}
	if result.Invoice.ID != "" {
		for _, eventType := range invoiceTerminalEvents(result.Invoice.Status, result.PaymentIntent.Status) {
			appendPayload(eventType, result.Invoice.ID, stripeInvoice(result.Invoice))
		}
	}
	if result.Subscription.ID != "" {
		appendPayload("customer.subscription.updated", result.Subscription.ID, h.stripeSubscription(r, result.Subscription))
	}
	return out
}

func paymentIntentTerminalEvent(status string) string {
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
	case "requires_payment_method":
		return "payment_intent.payment_failed"
	default:
		return "payment_intent.payment_failed"
	}
}

func setupIntentTerminalEvent(status string) string {
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

func checkoutSessionEvent(status string) string {
	if status == "expired" {
		return "checkout.session.expired"
	}
	return "checkout.session.completed"
}

func invoiceTerminalEvents(status string, paymentIntentStatus string) []string {
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
