package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBillingPortalConfigurationsCRUD(t *testing.T) {
	handler := newTestHandler(t)

	// Unknown param → 400 parameter_unknown.
	status, body := postFormStatus(t, handler, "/v1/billing_portal/configurations", url.Values{
		"business_profile[name]": {"Nope"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("unknown business_profile key status = %d body = %s, want 400", status, body)
	}
	// Invalid feature enum → 400.
	status, body = postFormStatus(t, handler, "/v1/billing_portal/configurations", url.Values{
		"features[subscription_cancel][mode]": {"whenever"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("invalid cancel mode status = %d body = %s, want 400", status, body)
	}
	// Invalid allowed_updates member → 400.
	status, body = postFormStatus(t, handler, "/v1/billing_portal/configurations", url.Values{
		"features[customer_update][allowed_updates][0]": {"address"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("invalid allowed_updates status = %d body = %s, want 400", status, body)
	}

	created := postForm[map[string]any](t, handler, "/v1/billing_portal/configurations", url.Values{
		"business_profile[headline]":                        {"Acme Billing"},
		"business_profile[privacy_policy_url]":              {"https://acme.test/privacy"},
		"default_return_url":                                {"https://acme.test/return"},
		"login_page[logo_url]":                              {"https://acme.test/logo.png"},
		"features[customer_update][enabled]":                {"true"},
		"features[customer_update][allowed_updates][0]":     {"email"},
		"features[subscription_cancel][mode]":               {"immediately"},
		"features[subscription_update][proration_behavior]": {"always_invoice"},
		"metadata[env]":                                     {"ci"},
	})
	id := fmt.Sprint(created["id"])
	if id == "" || created["object"] != "billing_portal.configuration" {
		t.Fatalf("created configuration = %#v, want billing_portal.configuration with id", created)
	}
	if created["active"] != true || created["is_default"] != true || created["livemode"] != false {
		t.Fatalf("created flags = %#v, want active default non-live", created)
	}
	profile, _ := created["business_profile"].(map[string]any)
	if profile["headline"] != "Acme Billing" || profile["privacy_policy_url"] != "https://acme.test/privacy" || profile["terms_of_service_url"] != nil {
		t.Fatalf("created business_profile = %#v", profile)
	}
	features, _ := created["features"].(map[string]any)
	cancel, _ := features["subscription_cancel"].(map[string]any)
	if cancel["mode"] != "immediately" {
		t.Fatalf("created subscription_cancel = %#v, want mode immediately", cancel)
	}
	customerUpdate, _ := features["customer_update"].(map[string]any)
	allowed, _ := customerUpdate["allowed_updates"].([]any)
	if customerUpdate["enabled"] != true || len(allowed) != 1 || fmt.Sprint(allowed[0]) != "email" {
		t.Fatalf("created customer_update = %#v, want enabled with email only", customerUpdate)
	}
	update, _ := features["subscription_update"].(map[string]any)
	if update["proration_behavior"] != "always_invoice" {
		t.Fatalf("created subscription_update = %#v, want always_invoice", update)
	}
	if _, ok := features["invoice_history"]; !ok {
		t.Fatalf("created features missing default invoice_history: %#v", features)
	}
	loginPage, _ := created["login_page"].(map[string]any)
	if fmt.Sprint(loginPage["logo_url"]) != "https://acme.test/logo.png" {
		t.Fatalf("created login_page = %#v", loginPage)
	}
	meta, _ := created["metadata"].(map[string]any)
	if fmt.Sprint(meta["env"]) != "ci" {
		t.Fatalf("created metadata = %#v, want env=ci", created["metadata"])
	}

	fetched := getJSON[map[string]any](t, handler, "/v1/billing_portal/configurations/"+id)
	if fetched["id"] != id || fetched["default_return_url"] != "https://acme.test/return" {
		t.Fatalf("GET configuration = %#v, want %s with return url", fetched, id)
	}

	// Second configuration is not the default; list filters narrow by flag.
	second := postForm[map[string]any](t, handler, "/v1/billing_portal/configurations", url.Values{
		"business_profile[headline]": {"Secondary"},
	})
	if second["is_default"] != false {
		t.Fatalf("second configuration is_default = %#v, want false", second["is_default"])
	}
	listed := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/billing_portal/configurations")
	if len(listed.Data) != 2 {
		t.Fatalf("list configurations = %d items, want 2", len(listed.Data))
	}
	defaults := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/billing_portal/configurations?is_default=true")
	if len(defaults.Data) != 1 || fmt.Sprint(defaults.Data[0]["id"]) != id {
		t.Fatalf("is_default filter = %#v, want the first configuration only", defaults.Data)
	}

	// Retrieve misses are 404.
	missingReq := httptest.NewRequest(http.MethodGet, "/v1/billing_portal/configurations/bpc_missing", nil)
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing configuration status = %d body = %s, want 404", missingRec.Code, missingRec.Body.String())
	}

	updated := postForm[map[string]any](t, handler, "/v1/billing_portal/configurations/"+id, url.Values{
		"active":                              {"false"},
		"business_profile[headline]":          {"Acme Billing v2"},
		"features[subscription_cancel][mode]": {"at_period_end"},
		"features[invoice_history][enabled]":  {"false"},
		"metadata[env]":                       {"ci2"},
	})
	if updated["active"] != false {
		t.Fatalf("updated active = %#v, want false", updated["active"])
	}
	profile, _ = updated["business_profile"].(map[string]any)
	if profile["headline"] != "Acme Billing v2" || profile["privacy_policy_url"] != "https://acme.test/privacy" {
		t.Fatalf("updated business_profile = %#v, want new headline keeping privacy url", profile)
	}
	features, _ = updated["features"].(map[string]any)
	cancel, _ = features["subscription_cancel"].(map[string]any)
	if cancel["mode"] != "at_period_end" {
		t.Fatalf("updated subscription_cancel = %#v, want at_period_end", cancel)
	}
	invoiceHistory, _ := features["invoice_history"].(map[string]any)
	if invoiceHistory["enabled"] != false {
		t.Fatalf("updated invoice_history = %#v, want disabled", invoiceHistory)
	}
	meta, _ = updated["metadata"].(map[string]any)
	if fmt.Sprint(meta["env"]) != "ci2" {
		t.Fatalf("updated metadata = %#v, want merged env=ci2", updated["metadata"])
	}
	activeFiltered := getJSON[struct {
		Data []map[string]any `json:"data"`
	}](t, handler, "/v1/billing_portal/configurations?active=false")
	if len(activeFiltered.Data) != 1 || fmt.Sprint(activeFiltered.Data[0]["id"]) != id {
		t.Fatalf("active=false filter = %#v, want the updated configuration", activeFiltered.Data)
	}
}
