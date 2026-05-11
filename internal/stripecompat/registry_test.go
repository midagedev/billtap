package stripecompat

import (
	"net/http"
	"testing"
)

func TestDefaultRegistryContainsCurrentPublicClaims(t *testing.T) {
	registry := DefaultRegistry()
	claims := registry.Claims()
	if len(claims) != 140 {
		t.Fatalf("default claims = %d, want 140", len(claims))
	}

	checkout, ok := registry.Lookup(http.MethodPost, "/v1/checkout/sessions")
	if !ok || checkout.Level != "L4" || !checkout.Stateful {
		t.Fatalf("checkout claim = %#v ok=%t, want L4 stateful", checkout, ok)
	}
	if len(checkout.ScorecardCases) != 1 || checkout.ScorecardCases[0] != "checkout.sessions.create.java_sdk_optional_params" {
		t.Fatalf("checkout scorecard cases = %#v", checkout.ScorecardCases)
	}

	subscription, ok := registry.Lookup(http.MethodGet, "/v1/subscriptions/{subscription}")
	if !ok || subscription.Level != "L3" {
		t.Fatalf("subscription normalized lookup = %#v ok=%t, want L3", subscription, ok)
	}

	concreteSubscription, ok := registry.Lookup(http.MethodGet, "/v1/subscriptions/sub_123")
	if !ok || concreteSubscription.Level != "L3" {
		t.Fatalf("subscription concrete lookup = %#v ok=%t, want L3", concreteSubscription, ok)
	}

	account, ok := registry.Lookup(http.MethodGet, "/v1/accounts/{account}")
	if !ok || account.Level != "L3" || !account.Stateful {
		t.Fatalf("connect account claim = %#v ok=%t, want L3 stateful", account, ok)
	}
	accountSession, ok := registry.Lookup(http.MethodPost, "/v1/account_sessions")
	if !ok || accountSession.Level != "L2" || accountSession.Stateful {
		t.Fatalf("connect account session claim = %#v ok=%t, want L2 smoke", accountSession, ok)
	}
	platformAccount, ok := registry.Lookup(http.MethodGet, "/v1/account")
	if !ok || platformAccount.Level != "L2" || platformAccount.Stateful {
		t.Fatalf("platform account claim = %#v ok=%t, want L2 non-stateful", platformAccount, ok)
	}
	personCreate, ok := registry.Lookup(http.MethodPost, "/v1/accounts/acct_123/people")
	if !ok || personCreate.Level != "L3" || !personCreate.Stateful {
		t.Fatalf("connect people claim = %#v ok=%t, want L3 stateful", personCreate, ok)
	}
	personDelete, ok := registry.Lookup(http.MethodDelete, "/v1/accounts/acct_123/persons/person_123")
	if !ok || personDelete.Level != "L3" || !personDelete.Stateful {
		t.Fatalf("connect persons delete claim = %#v ok=%t, want L3 stateful", personDelete, ok)
	}
	transfer, ok := registry.Lookup(http.MethodPost, "/v1/transfers")
	if !ok || transfer.Level != "L3" || !transfer.Stateful {
		t.Fatalf("connect transfer claim = %#v ok=%t, want L3 stateful", transfer, ok)
	}
	payoutCancel, ok := registry.Lookup(http.MethodPost, "/v1/payouts/po_123/cancel")
	if !ok || payoutCancel.Level != "L3" || !payoutCancel.Stateful {
		t.Fatalf("connect payout cancel claim = %#v ok=%t, want L3 stateful", payoutCancel, ok)
	}

	confirm, ok := registry.Lookup(http.MethodPost, "/v1/payment_intents/pi_123/confirm")
	if !ok || confirm.Level != "L3" || !confirm.Stateful {
		t.Fatalf("payment intent confirm claim = %#v ok=%t, want L3 stateful", confirm, ok)
	}
	setupConfirm, ok := registry.Lookup(http.MethodPost, "/v1/setup_intents/seti_123/confirm")
	if !ok || setupConfirm.Level != "L3" || !setupConfirm.Stateful {
		t.Fatalf("setup intent confirm claim = %#v ok=%t, want L3 stateful", setupConfirm, ok)
	}
	customerPaymentMethods, ok := registry.Lookup(http.MethodGet, "/v1/customers/cus_123/payment_methods")
	if !ok || customerPaymentMethods.Level != "L2" {
		t.Fatalf("customer payment methods claim = %#v ok=%t, want L2", customerPaymentMethods, ok)
	}
	invoicePay, ok := registry.Lookup(http.MethodPost, "/v1/invoices/in_123/pay")
	if !ok || invoicePay.Level != "L3" || !invoicePay.Stateful {
		t.Fatalf("invoice pay claim = %#v ok=%t, want L3 stateful", invoicePay, ok)
	}
	portal, ok := registry.Lookup(http.MethodPost, "/v1/billing_portal/sessions")
	if !ok || portal.Level != "L3" || len(portal.WebhookEvents) < 2 {
		t.Fatalf("portal claim = %#v ok=%t, want webhook-backed L3 portal session", portal, ok)
	}
	schedule, ok := registry.Lookup(http.MethodPost, "/v1/subscription_schedules/sub_sched_123/release")
	if !ok || schedule.Level != "L2" || !schedule.Stateful {
		t.Fatalf("subscription schedule release claim = %#v ok=%t, want L2 stateful", schedule, ok)
	}
	cashBalance, ok := registry.Lookup(http.MethodPost, "/v1/test_helpers/customers/cus_123/fund_cash_balance")
	if !ok || cashBalance.Level != "L3" || !cashBalance.Stateful {
		t.Fatalf("cash-balance funding claim = %#v ok=%t, want L3 stateful", cashBalance, ok)
	}
	dispute, ok := registry.Lookup(http.MethodPost, "/v1/charges/ch_123/dispute")
	if !ok || dispute.Level != "L2" || !dispute.Stateful {
		t.Fatalf("dispute claim = %#v ok=%t, want L2 stateful", dispute, ok)
	}
	attempts, ok := registry.Lookup(http.MethodGet, "/v1/webhook_endpoints/we_123/attempts")
	if !ok || attempts.Level != "L5" || !attempts.Stateful {
		t.Fatalf("webhook attempts claim = %#v ok=%t, want L5 stateful", attempts, ok)
	}
}

func TestNewRegistryRejectsDuplicateClaims(t *testing.T) {
	_, err := NewRegistry([]Claim{
		{Method: http.MethodGet, Path: "/v1/customers/{customer}", Level: "L1"},
		{Method: http.MethodGet, Path: "/v1/customers/{id}", Level: "L1"},
	})
	if err == nil {
		t.Fatalf("NewRegistry returned nil error for duplicate normalized claims")
	}
}

func TestNewRegistryRejectsInvalidLevel(t *testing.T) {
	_, err := NewRegistry([]Claim{
		{Method: http.MethodGet, Path: "/v1/customers/{customer}", Level: "L7"},
	})
	if err == nil {
		t.Fatalf("NewRegistry returned nil error for unsupported level")
	}
}

func TestRegistryDefensivelyCopiesClaims(t *testing.T) {
	seed := []Claim{
		{Method: http.MethodGet, Path: "/v1/customers/{customer}", Level: "L1", Risks: []string{"seed"}},
	}
	registry, err := NewRegistry(seed)
	if err != nil {
		t.Fatalf("NewRegistry returned error: %v", err)
	}

	seed[0].Risks[0] = "mutated input"
	claim, ok := registry.Lookup(http.MethodGet, "/v1/customers/{customer}")
	if !ok {
		t.Fatalf("registered claim not found")
	}
	if claim.Risks[0] != "seed" {
		t.Fatalf("registry retained mutable input slice: %#v", claim.Risks)
	}

	claim.Risks[0] = "mutated lookup"
	fresh, ok := registry.Lookup(http.MethodGet, "/v1/customers/{customer}")
	if !ok {
		t.Fatalf("registered claim not found after lookup mutation")
	}
	if fresh.Risks[0] != "seed" {
		t.Fatalf("Lookup exposed mutable claim slice: %#v", fresh.Risks)
	}

	defaultRegistry := DefaultRegistry()
	claims := defaultRegistry.Claims()
	for i := range claims {
		if len(claims[i].Risks) > 0 {
			claims[i].Risks[0] = "mutated claims"
			break
		}
	}

	freshClaims := defaultRegistry.Claims()
	for _, claim := range freshClaims {
		for _, risk := range claim.Risks {
			if risk == "mutated claims" {
				t.Fatalf("registry claims exposed mutable risk slices")
			}
		}
	}
}

func TestDefaultRouteCatalogContainsLatestOpenAPIRoutes(t *testing.T) {
	catalog := DefaultRouteCatalog()
	routes := catalog.Routes()
	if len(routes) != 619 {
		t.Fatalf("default known routes = %d, want 619", len(routes))
	}

	confirm, ok := catalog.Lookup(http.MethodPost, "/v1/payment_intents/pi_123/confirm")
	if !ok || confirm.Path != "/v1/payment_intents/{id}/confirm" {
		t.Fatalf("payment intent confirm route = %#v ok=%t, want normalized known route", confirm, ok)
	}
	if confirm.Source == "" {
		t.Fatalf("payment intent confirm source = empty")
	}

	if _, ok := catalog.Lookup(http.MethodGet, "/v1/not_a_stripe_route"); ok {
		t.Fatalf("unknown route unexpectedly matched known catalog")
	}
}

func TestNewRouteCatalogRejectsDuplicateRoutes(t *testing.T) {
	_, err := NewRouteCatalog([]Route{
		{Method: http.MethodGet, Path: "/v1/customers/{customer}"},
		{Method: http.MethodGet, Path: "/v1/customers/{id}"},
	})
	if err == nil {
		t.Fatalf("NewRouteCatalog returned nil error for duplicate normalized route")
	}
}
