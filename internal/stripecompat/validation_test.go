package stripecompat

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultValidationCatalogContainsOpenAPIRouteSchemas(t *testing.T) {
	catalog := DefaultValidationCatalog()
	operations := catalog.Operations()
	if len(operations) != 619 {
		t.Fatalf("default validation operations = %d, want 619", len(operations))
	}

	countrySpecs, ok := catalog.Lookup(http.MethodGet, "/v1/country_specs")
	if !ok || countrySpecs.OperationID != "GetCountrySpecs" {
		t.Fatalf("country specs validation = %#v ok=%t, want OpenAPI operation", countrySpecs, ok)
	}
	if !hasParam(countrySpecs, "limit", "integer") || !hasParam(countrySpecs, "expand", "array") {
		t.Fatalf("country specs params = %#v, want limit integer and expand array", countrySpecs.Params)
	}

	appSecret, ok := catalog.Lookup(http.MethodPost, "/v1/apps/secrets")
	if !ok || appSecret.OperationID != "PostAppsSecrets" {
		t.Fatalf("app secrets validation = %#v ok=%t, want OpenAPI operation", appSecret, ok)
	}
	if !hasRequiredParam(appSecret, "name") || !hasRequiredParam(appSecret, "payload") || !hasRequiredParam(appSecret, "scope") {
		t.Fatalf("app secrets params = %#v, want required name, payload, and scope", appSecret.Params)
	}
	if !hasNestedEnum(appSecret, "scope", "type", "account") {
		t.Fatalf("app secrets scope param = %#v, want scope[type] enum", appSecret.Params)
	}

	confirm, ok := catalog.Lookup(http.MethodPost, "/v1/payment_intents/pi_123/confirm")
	if !ok || confirm.Path != "/v1/payment_intents/{id}/confirm" {
		t.Fatalf("payment intent confirm validation = %#v ok=%t, want normalized lookup", confirm, ok)
	}
}

func TestValidationCatalogRejectsUnknownTypeMissingAndEnumParams(t *testing.T) {
	catalog := DefaultValidationCatalog()

	req := httptest.NewRequest(http.MethodGet, "/v1/country_specs?limit=not-an-int", nil)
	if err := catalog.Validate(req); err == nil || err.Code != ValidationCodeParamInvalid || err.Param != "limit" {
		t.Fatalf("limit validation err = %#v, want invalid limit", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/country_specs?nickname=legacy", nil)
	if err := catalog.Validate(req); err == nil || err.Code != ValidationCodeParamUnknown || err.Param != "nickname" {
		t.Fatalf("unknown validation err = %#v, want unknown nickname", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/v1/apps/secrets", strings.NewReader("payload=secret&scope[type]=account"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := catalog.Validate(req); err == nil || err.Code != ValidationCodeParamMissing || err.Param != "name" {
		t.Fatalf("missing validation err = %#v, want missing name", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/v1/apps/secrets", strings.NewReader("name=token&payload=secret&scope[type]=workspace"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := catalog.Validate(req); err == nil || err.Code != ValidationCodeParamInvalid || err.Param != "scope[type]" {
		t.Fatalf("enum validation err = %#v, want invalid scope[type]", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/v1/account_sessions", strings.NewReader("account=acct_123&components[account_onboarding][enabled]=maybe"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := catalog.Validate(req); err == nil || err.Code != ValidationCodeParamInvalid || err.Param != "components[account_onboarding][enabled]" {
		t.Fatalf("deep nested validation err = %#v, want invalid components[account_onboarding][enabled]", err)
	}
}

func hasParam(operation OperationValidation, name string, typ string) bool {
	for _, param := range operation.Params {
		if param.Name == name && param.Type == typ {
			return true
		}
	}
	return false
}

func hasRequiredParam(operation OperationValidation, name string) bool {
	for _, param := range operation.Params {
		if param.Name == name && param.Required {
			return true
		}
	}
	return false
}

func hasNestedEnum(operation OperationValidation, name string, childName string, enumValue string) bool {
	for _, param := range operation.Params {
		if param.Name != name {
			continue
		}
		for _, child := range param.Children {
			if child.Name != childName {
				continue
			}
			for _, value := range child.Enum {
				if value == enumValue {
					return true
				}
			}
		}
	}
	return false
}
