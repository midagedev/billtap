package stripecompat

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const DefaultDocs = "docs/COMPATIBILITY.md#supported-stripe-like-api-subset"

type Claim struct {
	Method         string
	Path           string
	Level          string
	Stateful       bool
	Docs           string
	WebhookEvents  []string
	ScorecardCases []string
	SDKSmoke       []string
	Risks          []string
}

type Registry struct {
	claims  map[string]Claim
	ordered []Claim
}

func NewRegistry(claims []Claim) (Registry, error) {
	registry := Registry{claims: map[string]Claim{}}
	for _, claim := range claims {
		claim.Method = strings.ToUpper(strings.TrimSpace(claim.Method))
		claim.Path = NormalizePath(claim.Path)
		if claim.Method == "" {
			return Registry{}, fmt.Errorf("claim for %s has empty method", claim.Path)
		}
		if claim.Path == "" {
			return Registry{}, fmt.Errorf("claim for %s has empty path", claim.Method)
		}
		if claim.Level == "" {
			return Registry{}, fmt.Errorf("claim for %s %s has empty level", claim.Method, claim.Path)
		}
		if !validLevel(claim.Level) {
			return Registry{}, fmt.Errorf("claim for %s %s has unsupported level %q", claim.Method, claim.Path, claim.Level)
		}
		if claim.Docs == "" {
			claim.Docs = DefaultDocs
		}
		key := RouteKey(claim.Method, claim.Path)
		if _, exists := registry.claims[key]; exists {
			return Registry{}, fmt.Errorf("duplicate compatibility claim: %s", key)
		}
		copied := cloneClaim(claim)
		registry.claims[key] = copied
		registry.ordered = append(registry.ordered, copied)
	}
	return registry, nil
}

func MustRegistry(claims []Claim) Registry {
	registry, err := NewRegistry(claims)
	if err != nil {
		panic(err)
	}
	return registry
}

func DefaultRegistry() Registry {
	return MustRegistry(DefaultClaims())
}

func (r Registry) Lookup(method string, path string) (Claim, bool) {
	claim, ok := r.claims[RouteKey(method, NormalizePath(path))]
	if ok {
		return cloneClaim(claim), true
	}
	normalizedMethod := strings.ToUpper(strings.TrimSpace(method))
	for _, claim := range r.ordered {
		if claim.Method == normalizedMethod && pathMatches(claim.Path, path) {
			return cloneClaim(claim), true
		}
	}
	return Claim{}, false
}

func (r Registry) Claims() []Claim {
	claims := make([]Claim, 0, len(r.ordered))
	for _, claim := range r.ordered {
		claims = append(claims, cloneClaim(claim))
	}
	return claims
}

func RouteKey(method string, normalizedPath string) string {
	return strings.ToUpper(strings.TrimSpace(method)) + " " + NormalizePath(normalizedPath)
}

var pathParamPattern = regexp.MustCompile(`\{[^}/]+\}`)

func NormalizePath(path string) string {
	return pathParamPattern.ReplaceAllString(strings.TrimSpace(path), "{id}")
}

func pathMatches(template string, candidate string) bool {
	template = NormalizePath(template)
	candidate = strings.TrimSpace(strings.Split(candidate, "?")[0])
	if NormalizePath(candidate) == template {
		return true
	}

	templateParts := splitPath(template)
	candidateParts := splitPath(candidate)
	if len(templateParts) != len(candidateParts) {
		return false
	}
	for i := range templateParts {
		if templateParts[i] == "{id}" {
			if candidateParts[i] == "" {
				return false
			}
			continue
		}
		if templateParts[i] != candidateParts[i] {
			return false
		}
	}
	return true
}

func splitPath(path string) []string {
	path = strings.Trim(strings.TrimSpace(path), "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

func validLevel(level string) bool {
	switch level {
	case "L1", "L2", "L3", "L4", "L5", "L6":
		return true
	default:
		return false
	}
}

func DefaultClaims() []Claim {
	var claims []Claim
	add := func(method string, path string, claim Claim) {
		claim.Method = method
		claim.Path = path
		claims = append(claims, claim)
	}
	statefulL3 := Claim{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}}
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		add(method, "/v1/customers", statefulL3)
		add(method, "/v1/customers/{id}", statefulL3)
		add(method, "/v1/products", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"products.create.success"}, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/products/{id}", statefulL3)
		add(method, "/v1/prices", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"prices.create.invalid_json_amount_type"}, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/prices/{id}", statefulL3)
		add(method, "/v1/subscriptions", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted"}})
		add(method, "/v1/subscriptions/{id}", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.updated", "customer.subscription.deleted"}})
	}
	add(http.MethodGet, "/v1/products/search", Claim{Level: "L2", Risks: []string{"metadata equality filters only; no Stripe Search Query Language parity"}})
	add(http.MethodGet, "/v1/prices/search", Claim{Level: "L3", Stateful: true, Risks: []string{"supports a measured prices search subset for active, type, lookup_key, and metadata equality clauses joined by AND"}})

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		add(method, "/v1/accounts", Claim{Level: "L3", Stateful: true, Risks: []string{"local Connect smoke only; onboarding, KYC, external accounts, and platform settlement are not modeled"}})
		add(method, "/v1/accounts/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"local Connect smoke only; onboarding, KYC, external accounts, and platform settlement are not modeled"}})
	}
	add(http.MethodGet, "/v1/account", Claim{Level: "L2", Stateful: false, Risks: []string{"returns deterministic local platform-account evidence only"}})
	add(http.MethodDelete, "/v1/accounts/{id}", Claim{Level: "L2", Stateful: false, Risks: []string{"returns local account deletion evidence only; provider-side account closure is not modeled"}})
	add(http.MethodPost, "/v1/account_links", Claim{Level: "L2", Stateful: false, Risks: []string{"returns local hosted onboarding/update URLs only; no real Connect onboarding session"}})
	add(http.MethodPost, "/v1/account_sessions", Claim{Level: "L2", Stateful: false, Risks: []string{"returns deterministic local client secrets for embedded-component smoke only"}})
	add(http.MethodGet, "/v1/accounts/{id}/capabilities", Claim{Level: "L2", Stateful: true, Risks: []string{"local account capability projection only; requirements are not provider-verified"}})
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		add(method, "/v1/accounts/{id}/capabilities/{id}", Claim{Level: "L2", Stateful: true, Risks: []string{"local account capability projection only; requirements are not provider-verified"}})
		add(method, "/v1/accounts/{id}/people", Claim{Level: "L3", Stateful: true, Risks: []string{"local person evidence only; identity verification and KYC are not modeled"}})
		add(method, "/v1/accounts/{id}/people/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"local person evidence only; identity verification and KYC are not modeled"}})
		add(method, "/v1/accounts/{id}/persons", Claim{Level: "L3", Stateful: true, Risks: []string{"alias for local person evidence only; identity verification and KYC are not modeled"}})
		add(method, "/v1/accounts/{id}/persons/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"alias for local person evidence only; identity verification and KYC are not modeled"}})
		add(method, "/v1/accounts/{id}/external_accounts", Claim{Level: "L3", Stateful: true, Risks: []string{"local external-account evidence only; bank verification and payouts settlement are not modeled"}})
		add(method, "/v1/accounts/{id}/external_accounts/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"local external-account evidence only; bank verification and payouts settlement are not modeled"}})
		add(method, "/v1/accounts/{id}/bank_accounts/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"alias for local external bank-account evidence; verification is not modeled"}})
		add(method, "/v1/transfers", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"transfer.created", "transfer.reversed"}, Risks: []string{"local transfer evidence only; balance movement and settlement are not modeled"}})
		add(method, "/v1/transfers/{id}", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"transfer.created", "transfer.reversed"}, Risks: []string{"local transfer evidence only; balance movement and settlement are not modeled"}})
		add(method, "/v1/transfers/{id}/reversals", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"transfer.reversed"}, Risks: []string{"local transfer reversal evidence only; balance movement is not modeled"}})
		add(method, "/v1/transfers/{id}/reversals/{id}", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"transfer.reversed"}, Risks: []string{"local transfer reversal evidence only; balance movement is not modeled"}})
		add(method, "/v1/payouts", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"payout.created", "payout.canceled", "payout.reversed"}, Risks: []string{"local payout evidence only; banking rails and settlement are not modeled"}})
		add(method, "/v1/payouts/{id}", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"payout.created", "payout.canceled", "payout.reversed"}, Risks: []string{"local payout evidence only; banking rails and settlement are not modeled"}})
		add(method, "/v1/application_fees/{id}/refunds", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"application_fee.refunded"}, Risks: []string{"local application-fee refund evidence only; ledger and balance transactions are not modeled"}})
		add(method, "/v1/application_fees/{id}/refunds/{id}", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"application_fee.refunded"}, Risks: []string{"local application-fee refund evidence only; ledger and balance transactions are not modeled"}})
	}
	add(http.MethodPost, "/v1/accounts/{id}/bank_accounts", Claim{Level: "L3", Stateful: true, Risks: []string{"alias for local external bank-account evidence; verification is not modeled"}})
	add(http.MethodDelete, "/v1/accounts/{id}/external_accounts/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"local external-account deletion evidence only"}})
	add(http.MethodDelete, "/v1/accounts/{id}/bank_accounts/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"local external-account deletion evidence only"}})
	add(http.MethodDelete, "/v1/accounts/{id}/people/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"local person deletion evidence only"}})
	add(http.MethodDelete, "/v1/accounts/{id}/persons/{id}", Claim{Level: "L3", Stateful: true, Risks: []string{"alias for local person deletion evidence only"}})
	add(http.MethodPost, "/v1/accounts/{id}/login_links", Claim{Level: "L2", Stateful: false, Risks: []string{"returns local login URLs only; no Express dashboard session"}})
	add(http.MethodPost, "/v1/accounts/{id}/reject", Claim{Level: "L2", Stateful: true, Risks: []string{"records local rejection metadata only; no provider risk review"}})
	add(http.MethodPost, "/v1/payouts/{id}/cancel", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"payout.canceled"}, Risks: []string{"local payout cancellation evidence only"}})
	add(http.MethodPost, "/v1/payouts/{id}/reverse", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"payout.reversed"}, Risks: []string{"local payout reversal evidence only"}})
	add(http.MethodGet, "/v1/application_fees", Claim{Level: "L2", Stateful: true, Risks: []string{"local application-fee evidence only; fees are materialized by local simulation paths"}})
	add(http.MethodGet, "/v1/application_fees/{id}", Claim{Level: "L2", Stateful: true, Risks: []string{"local application-fee evidence only; ledger behavior is not modeled"}})
	add(http.MethodPost, "/v1/application_fees/{id}/refund", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"application_fee.refunded"}, Risks: []string{"legacy local fee-refund evidence only"}})

	add(http.MethodPost, "/v1/checkout/sessions", Claim{Level: "L4", Stateful: true, ScorecardCases: []string{"checkout.sessions.create.java_sdk_optional_params"}, SDKSmoke: []string{"stripe-node"}, Risks: []string{"subscription mode only"}})
	add(http.MethodGet, "/v1/checkout/sessions", Claim{Level: "L4", Stateful: true, SDKSmoke: []string{"stripe-node"}, Risks: []string{"subscription mode only"}})
	add(http.MethodGet, "/v1/checkout/sessions/{id}", Claim{Level: "L4", Stateful: true, SDKSmoke: []string{"stripe-node"}, Risks: []string{"subscription mode only"}})
	add(http.MethodPost, "/v1/billing_portal/sessions", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.updated", "customer.subscription.deleted", "payment_method.attached", "customer.updated"}, Risks: []string{"hosted portal is a local stub; portal configuration rendering and full Stripe-hosted portal behavior are not modeled"}})

	add(http.MethodGet, "/v1/coupons", Claim{Level: "L2", Stateful: true, Risks: []string{"local coupon evidence only; discount accounting is bounded"}})
	add(http.MethodPost, "/v1/coupons", Claim{Level: "L2", Stateful: true, Risks: []string{"local coupon evidence only; discount accounting is bounded"}})
	add(http.MethodGet, "/v1/coupons/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/coupons/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodDelete, "/v1/coupons/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodGet, "/v1/promotion_codes", Claim{Level: "L2", Stateful: true, Risks: []string{"local promotion-code evidence only; redemption limits are not modeled"}})
	add(http.MethodPost, "/v1/promotion_codes", Claim{Level: "L2", Stateful: true, Risks: []string{"local promotion-code evidence only; redemption limits are not modeled"}})
	add(http.MethodGet, "/v1/promotion_codes/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/promotion_codes/{id}", Claim{Level: "L2", Stateful: true})

	add(http.MethodPost, "/v1/subscription_items", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"subscription_items.create.invalid_quantity"}})
	add(http.MethodDelete, "/v1/subscription_items/{id}", Claim{Level: "L3", Stateful: true})
	add(http.MethodDelete, "/v1/subscriptions/{id}", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.deleted"}})
	add(http.MethodGet, "/v1/subscription_schedules", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"customer.subscription.updated"}, Risks: []string{"local phase evidence only; complex proration and phase billing are not modeled"}})
	add(http.MethodPost, "/v1/subscription_schedules", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"customer.subscription.updated"}, Risks: []string{"local phase evidence only; complex proration and phase billing are not modeled"}})
	add(http.MethodGet, "/v1/subscription_schedules/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/subscription_schedules/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/subscription_schedules/{id}/cancel", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/subscription_schedules/{id}/release", Claim{Level: "L2", Stateful: true})

	add(http.MethodGet, "/v1/invoices", statefulL3)
	add(http.MethodGet, "/v1/invoices/{id}", statefulL3)
	add(http.MethodPost, "/v1/invoices/{id}/pay", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"invoices.pay.failed_invoice_succeeds", "invoices.pay.failed_invoice_declines_again"}, WebhookEvents: []string{"payment_intent.succeeded", "payment_intent.payment_failed", "invoice.payment_succeeded", "invoice.payment_failed", "invoice.paid", "customer.subscription.updated"}, Risks: []string{"local retry/payment mutation only; finalize, send, void, collection, and dunning automation are not modeled"}})
	add(http.MethodPost, "/v1/invoices/create_preview", Claim{Level: "L3", Stateful: true, Risks: []string{"local subscription-update proration subset; taxes, discounts, pending invoice items, and full invoice preview parity are not modeled"}})
	add(http.MethodGet, "/v1/invoices/upcoming", Claim{Level: "L3", Stateful: true, Risks: []string{"Stripe-compatible upcoming preview alias backed by the same local proration subset"}})
	add(http.MethodPost, "/v1/invoices/upcoming", Claim{Level: "L3", Stateful: true, Risks: []string{"local form-compatible upcoming preview convenience"}})
	add(http.MethodGet, "/v1/refunds", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"charge.refunded", "charge.refund.updated"}, Risks: []string{"local refund evidence only; charges, balances, payouts, and processor accounting are not modeled"}})
	add(http.MethodPost, "/v1/refunds", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"charge.refunded", "charge.refund.updated"}, Risks: []string{"local refund evidence only; charges, balances, payouts, and processor accounting are not modeled"}})
	add(http.MethodGet, "/v1/refunds/{id}", Claim{Level: "L3", Stateful: true})
	add(http.MethodPost, "/v1/refunds/{id}", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"charge.refund.updated"}})
	add(http.MethodPost, "/v1/refunds/{id}/cancel", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"charge.refund.updated"}})
	add(http.MethodGet, "/v1/credit_notes", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"credit_note.created", "credit_note.voided"}, Risks: []string{"local credit note evidence only; line/tax/customer-balance math is not modeled"}})
	add(http.MethodPost, "/v1/credit_notes", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"credit_note.created", "credit_note.voided"}, Risks: []string{"local credit note evidence only; line/tax/customer-balance math is not modeled"}})
	add(http.MethodGet, "/v1/credit_notes/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/credit_notes/{id}/void", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"credit_note.voided"}})

	add(http.MethodGet, "/v1/payment_intents", statefulL3)
	add(http.MethodPost, "/v1/payment_intents", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"payment_intents.create.confirm.succeeds", "payment_intents.confirm.card_decline"}, Risks: []string{"local state machine only; no card processing or full PaymentIntent parameter parity"}})
	add(http.MethodGet, "/v1/payment_intents/{id}", statefulL3)
	add(http.MethodPost, "/v1/payment_intents/{id}/confirm", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"payment_intents.confirm.card_decline"}, Risks: []string{"local deterministic outcome aliases only"}})
	add(http.MethodPost, "/v1/payment_intents/{id}/capture", Claim{Level: "L3", Stateful: true, Risks: []string{"local capture marks the intent succeeded; partial capture accounting is not modeled"}})
	add(http.MethodPost, "/v1/payment_intents/{id}/cancel", Claim{Level: "L3", Stateful: true})
	add(http.MethodGet, "/v1/setup_intents", statefulL3)
	add(http.MethodPost, "/v1/setup_intents", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"setup_intents.create.confirm.succeeds"}, Risks: []string{"local state machine only; mandates and full SCA flows are not modeled"}})
	add(http.MethodGet, "/v1/setup_intents/{id}", statefulL3)
	add(http.MethodPost, "/v1/setup_intents/{id}/confirm", Claim{Level: "L3", Stateful: true, ScorecardCases: []string{"setup_intents.create.confirm.succeeds"}, Risks: []string{"local deterministic outcome aliases only"}})
	add(http.MethodPost, "/v1/setup_intents/{id}/cancel", Claim{Level: "L3", Stateful: true})
	add(http.MethodGet, "/v1/payment_methods", Claim{Level: "L2", Risks: []string{"deterministic sandbox card projection only"}})
	add(http.MethodGet, "/v1/customers/{customer}/payment_methods", Claim{Level: "L2", Risks: []string{"deterministic sandbox card projection only"}})

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		add(method, "/v1/webhook_endpoints", Claim{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/webhook_endpoints/{id}", Claim{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	}
	add(http.MethodPatch, "/v1/webhook_endpoints/{id}", Claim{Level: "L5", Stateful: true, Risks: []string{"Billtap accepts PATCH as a local mutation convenience in addition to Stripe-style POST update"}})
	add(http.MethodDelete, "/v1/webhook_endpoints/{id}", Claim{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/webhook_endpoints/{id}/attempts", Claim{Level: "L5", Stateful: true, Risks: []string{"Billtap-specific endpoint-scoped delivery evidence; not a Stripe API route"}})
	add(http.MethodGet, "/v1/events", Claim{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/events/{id}", Claim{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/test_helpers/test_clocks", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.updated", "customer.subscription.deleted", "invoice.created", "invoice.paid", "invoice.payment_failed"}, Risks: []string{"deterministic local clock subset; not full Stripe test-clock attachment semantics"}})
	add(http.MethodPost, "/v1/test_helpers/test_clocks", Claim{Level: "L3", Stateful: true, Risks: []string{"deterministic local clock subset; not full Stripe test-clock attachment semantics"}})
	add(http.MethodGet, "/v1/test_helpers/test_clocks/{id}", Claim{Level: "L3", Stateful: true})
	add(http.MethodPost, "/v1/test_helpers/test_clocks/{id}/advance", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.updated", "customer.subscription.deleted", "invoice.created", "invoice.paid", "invoice.payment_failed"}})
	add(http.MethodGet, "/v1/customers/{id}/cash_balance", Claim{Level: "L2", Stateful: true, Risks: []string{"local customer cash-balance evidence only"}})
	add(http.MethodPost, "/v1/customers/{id}/cash_balance", Claim{Level: "L2", Stateful: true, Risks: []string{"local customer cash-balance evidence only"}})
	add(http.MethodGet, "/v1/customers/{id}/cash_balance_transactions", Claim{Level: "L2", Stateful: true})
	add(http.MethodGet, "/v1/customers/{id}/cash_balance_transactions/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/test_helpers/customers/{id}/fund_cash_balance", Claim{Level: "L3", Stateful: true, WebhookEvents: []string{"payment_intent.succeeded"}, Risks: []string{"local test-helper funding only; banking rails are not modeled"}})
	add(http.MethodGet, "/v1/disputes", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"charge.dispute.created", "charge.dispute.updated", "charge.dispute.funds_withdrawn", "charge.dispute.closed"}, Risks: []string{"local dispute evidence only; representment workflow and balance movements are not modeled"}})
	add(http.MethodGet, "/v1/disputes/{id}", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/disputes/{id}", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"charge.dispute.updated"}})
	add(http.MethodPost, "/v1/disputes/{id}/close", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"charge.dispute.closed"}})
	add(http.MethodGet, "/v1/charges/{id}/dispute", Claim{Level: "L2", Stateful: true})
	add(http.MethodPost, "/v1/charges/{id}/dispute", Claim{Level: "L2", Stateful: true, WebhookEvents: []string{"charge.dispute.created"}})

	return claims
}

func cloneClaim(claim Claim) Claim {
	claim.WebhookEvents = append([]string(nil), claim.WebhookEvents...)
	claim.ScorecardCases = append([]string(nil), claim.ScorecardCases...)
	claim.SDKSmoke = append([]string(nil), claim.SDKSmoke...)
	claim.Risks = append([]string(nil), claim.Risks...)
	return claim
}
