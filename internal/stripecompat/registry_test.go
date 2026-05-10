package stripecompat

import (
	"net/http"
	"testing"
)

func TestDefaultRegistryContainsCurrentPublicClaims(t *testing.T) {
	registry := DefaultRegistry()
	claims := registry.Claims()
	if len(claims) != 46 {
		t.Fatalf("default claims = %d, want 46", len(claims))
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

	if _, ok := registry.Lookup(http.MethodGet, "/v1/accounts/{account}"); ok {
		t.Fatalf("connect account should not be claimed before T6")
	}

	confirm, ok := registry.Lookup(http.MethodPost, "/v1/payment_intents/pi_123/confirm")
	if !ok || confirm.Level != "L3" || !confirm.Stateful {
		t.Fatalf("payment intent confirm claim = %#v ok=%t, want L3 stateful", confirm, ok)
	}
	setupConfirm, ok := registry.Lookup(http.MethodPost, "/v1/setup_intents/seti_123/confirm")
	if !ok || setupConfirm.Level != "L3" || !setupConfirm.Stateful {
		t.Fatalf("setup intent confirm claim = %#v ok=%t, want L3 stateful", setupConfirm, ok)
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
