package compatibility

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteInventoryArtifactsClassifiesOpenAPIOperations(t *testing.T) {
	dir := t.TempDir()
	inventory, err := WriteInventoryArtifacts(context.Background(), InventoryOptions{
		OpenAPIPath: filepath.Join("testdata", "stripe-openapi-minimal.json"),
		OutputDir:   dir,
		Source:      "stripe/openapi test fixture",
		Now:         fixedNow,
	})
	if err != nil {
		t.Fatalf("WriteInventoryArtifacts returned error: %v", err)
	}

	if inventory.InventoryVersion != InventoryVersion {
		t.Fatalf("inventory version = %q, want %q", inventory.InventoryVersion, InventoryVersion)
	}
	if inventory.Summary.TotalOperations != 14 {
		t.Fatalf("total operations = %d, want 14", inventory.Summary.TotalOperations)
	}
	if inventory.Summary.ImplementedOperations != 12 {
		t.Fatalf("implemented operations = %d, want 12", inventory.Summary.ImplementedOperations)
	}
	if inventory.Summary.SchemaValidatedOperations != 0 {
		t.Fatalf("schema validated operations = %d, want 0 because fixture has no parameter/requestBody schemas", inventory.Summary.SchemaValidatedOperations)
	}
	if inventory.Summary.InventoryOnlyOperations != 2 {
		t.Fatalf("inventory-only operations = %d, want 2", inventory.Summary.InventoryOnlyOperations)
	}
	if inventory.Summary.ImplementedPercent != 85.7 {
		t.Fatalf("implemented percent = %.1f, want 85.7", inventory.Summary.ImplementedPercent)
	}
	if inventory.Summary.BilltapOnlyRoutes != 1 {
		t.Fatalf("billtap-only routes = %d, want 1", inventory.Summary.BilltapOnlyRoutes)
	}
	if inventory.Summary.ByLevel["L0"] != 2 || inventory.Summary.ByLevel["L2"] != 2 || inventory.Summary.ByLevel["L3"] != 8 || inventory.Summary.ByLevel["L4"] != 1 || inventory.Summary.ByLevel["L5"] != 1 {
		t.Fatalf("by level = %#v, want L0=2 L2=2 L3=8 L4=1 L5=1", inventory.Summary.ByLevel)
	}

	customerCreate := findOperation(t, inventory, http.MethodPost, "/v1/customers")
	if !customerCreate.Implemented || customerCreate.BilltapLevel != "L3" || !customerCreate.Stateful {
		t.Fatalf("customer create coverage = %#v, want implemented L3 stateful", customerCreate)
	}
	if customerCreate.SchemaValidated {
		t.Fatalf("customer create schema_validated = true, want false because fixture has no parameter/requestBody schemas")
	}
	if customerCreate.Docs == "" {
		t.Fatalf("customer create docs = empty, want traceable compatibility docs")
	}
	customers := findFamily(t, inventory, "customers")
	if customers.Priority != "P1" || customers.TotalOperations != 4 || customers.ImplementedOperations != 4 || customers.ImplementedPercent != 100 {
		t.Fatalf("customers family = %#v, want P1 4/4 100%%", customers)
	}

	checkoutCreate := findOperation(t, inventory, http.MethodPost, "/v1/checkout/sessions")
	if !checkoutCreate.Implemented || checkoutCreate.BilltapLevel != "L4" || checkoutCreate.TargetLevel != "L4-L6" {
		t.Fatalf("checkout create coverage = %#v, want implemented L4 target L4-L6", checkoutCreate)
	}
	if !containsString(checkoutCreate.ScorecardCases, "checkout.sessions.create.java_sdk_optional_params") {
		t.Fatalf("checkout create scorecard cases = %#v, want checkout SDK params case", checkoutCreate.ScorecardCases)
	}

	webhookCreate := findOperation(t, inventory, http.MethodPost, "/v1/webhook_endpoints")
	if !webhookCreate.Implemented || webhookCreate.BilltapLevel != "L5" {
		t.Fatalf("webhook create coverage = %#v, want implemented L5", webhookCreate)
	}

	paymentIntentConfirm := findOperation(t, inventory, http.MethodPost, "/v1/payment_intents/{intent}/confirm")
	if !paymentIntentConfirm.Implemented || paymentIntentConfirm.BilltapLevel != "L3" || paymentIntentConfirm.TargetLevel != "L3-L6" {
		t.Fatalf("payment intent confirm coverage = %#v, want implemented L3 payment target", paymentIntentConfirm)
	}
	if !containsString(paymentIntentConfirm.ScorecardCases, "payment_intents.confirm.card_decline") {
		t.Fatalf("payment intent scorecard cases = %#v, want direct confirm decline case", paymentIntentConfirm.ScorecardCases)
	}

	accountsV2 := findOperation(t, inventory, http.MethodPost, "/v2/core/accounts")
	if accountsV2.Family != "connect" || accountsV2.BilltapLevel != "L0" || accountsV2.TargetLevel != "L2-L5" {
		t.Fatalf("accounts v2 coverage = %#v, want connect inventory target", accountsV2)
	}
	accountSession := findOperation(t, inventory, http.MethodPost, "/v1/account_sessions")
	if accountSession.Family != "connect" || accountSession.Resource != "account_session" || !accountSession.Implemented || accountSession.BilltapLevel != "L2" {
		t.Fatalf("account session coverage = %#v, want implemented connect account_session L2", accountSession)
	}
	account := findOperation(t, inventory, http.MethodGet, "/v1/account")
	if account.Family != "connect" || account.Resource != "account" || !account.Implemented || account.BilltapLevel != "L2" {
		t.Fatalf("account coverage = %#v, want implemented connect account L2", account)
	}
	applicationFeeRefund := findOperation(t, inventory, http.MethodGet, "/v1/application_fees/{fee}/refunds")
	if applicationFeeRefund.Family != "connect" || applicationFeeRefund.Resource != "application_fee_refund" || !applicationFeeRefund.Implemented || applicationFeeRefund.BilltapLevel != "L3" {
		t.Fatalf("application fee refund coverage = %#v, want implemented connect application fee refund L3", applicationFeeRefund)
	}
	financialConnectionAccount := findOperation(t, inventory, http.MethodGet, "/v1/financial_connections/accounts")
	if financialConnectionAccount.Family != "auxiliary" || financialConnectionAccount.Resource != "financial_connections.account" {
		t.Fatalf("financial connection account coverage = %#v, want auxiliary financial connection account", financialConnectionAccount)
	}
	connect := findFamily(t, inventory, "connect")
	if connect.Priority != "P1" || connect.TotalOperations != 4 || connect.ImplementedOperations != 3 || connect.ImplementedPercent != 75 {
		t.Fatalf("connect family = %#v, want P1 3/4 75%%", connect)
	}
	if !strings.Contains(connect.NextMilestone, "SDK/adoption") {
		t.Fatalf("connect milestone = %q, want Connect SDK/adoption follow-up", connect.NextMilestone)
	}

	jsonPath := filepath.Join(dir, "stripe-api-inventory.json")
	mdPath := filepath.Join(dir, "stripe-api-inventory.md")
	if !fileContains(t, jsonPath, `"inventory_version": "stripe-api-inventory-v2"`) {
		t.Fatalf("JSON inventory missing version")
	}
	if !fileContains(t, jsonPath, `"families"`) || !fileContains(t, jsonPath, `"implemented_percent": 85.7`) {
		t.Fatalf("JSON inventory missing measurable coverage fields")
	}
	if !fileContains(t, jsonPath, `"schema_validated_operations": 0`) || !fileContains(t, mdPath, "OpenAPI validation catalog") {
		t.Fatalf("inventory artifacts missing schema validation coverage fields")
	}
	if !fileContains(t, jsonPath, `"billtap_only_routes"`) || !fileContains(t, jsonPath, `/v1/checkout/sessions/{id}/complete`) {
		t.Fatalf("JSON inventory missing Billtap-specific route exception")
	}
	if !fileContains(t, mdPath, "# Stripe API Compatibility Inventory") || !fileContains(t, mdPath, "Family Coverage") || !fileContains(t, mdPath, "Billtap-Specific `/v1` Exceptions") {
		t.Fatalf("Markdown inventory missing expected sections")
	}
}

func TestGenerateInventorySchemaValidatedUsesInputOpenAPISurface(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "schema-surface.json")
	if err := os.WriteFile(specPath, []byte(`{
  "openapi": "3.0.0",
  "info": {"title": "schema surface", "version": "test"},
  "paths": {
    "/v1/country_specs": {
      "get": {
        "operationId": "GetCountrySpecs",
        "parameters": [
          {"name": "limit", "in": "query", "schema": {"type": "integer"}}
        ],
        "x-resourceId": "country_spec"
      }
    },
    "/v1/exchange_rates": {
      "get": {
        "operationId": "GetExchangeRates",
        "x-resourceId": "exchange_rate"
      }
    }
  }
}`), 0o644); err != nil {
		t.Fatalf("write schema surface spec: %v", err)
	}

	inventory, err := GenerateInventory(context.Background(), InventoryOptions{
		OpenAPIPath: specPath,
		Source:      "schema surface fixture",
		Now:         fixedNow,
	})
	if err != nil {
		t.Fatalf("GenerateInventory returned error: %v", err)
	}
	if inventory.Summary.TotalOperations != 2 || inventory.Summary.SchemaValidatedOperations != 1 || inventory.Summary.SchemaValidatedPercent != 50 {
		t.Fatalf("summary = %#v, want one of two operations schema-visible from input OpenAPI", inventory.Summary)
	}
	countrySpecs := findOperation(t, inventory, http.MethodGet, "/v1/country_specs")
	if !countrySpecs.SchemaValidated {
		t.Fatalf("country specs = %#v, want schema_validated true", countrySpecs)
	}
	exchangeRates := findOperation(t, inventory, http.MethodGet, "/v1/exchange_rates")
	if exchangeRates.SchemaValidated {
		t.Fatalf("exchange rates = %#v, want schema_validated false without input schema surface", exchangeRates)
	}
}

func TestGenerateInventoryRejectsInvalidOpenAPIInput(t *testing.T) {
	dir := t.TempDir()
	emptySpec := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(emptySpec, []byte(`{"openapi":"3.0.0","paths":{}}`), 0o644); err != nil {
		t.Fatalf("write empty spec: %v", err)
	}
	if _, err := GenerateInventory(context.Background(), InventoryOptions{OpenAPIPath: emptySpec}); err == nil {
		t.Fatalf("GenerateInventory returned nil error for empty paths")
	}

	if _, err := GenerateInventory(context.Background(), InventoryOptions{}); err == nil {
		t.Fatalf("GenerateInventory returned nil error for missing path")
	}
}

func TestInventoryJSONRoundTrip(t *testing.T) {
	inventory, err := GenerateInventory(context.Background(), InventoryOptions{
		OpenAPIPath: filepath.Join("testdata", "stripe-openapi-minimal.json"),
		Source:      "stripe/openapi test fixture",
		Now:         fixedNow,
	})
	if err != nil {
		t.Fatalf("GenerateInventory returned error: %v", err)
	}
	body, err := inventory.JSON()
	if err != nil {
		t.Fatalf("inventory JSON returned error: %v", err)
	}
	var decoded StripeAPIInventory
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("decode inventory JSON: %v", err)
	}
	if decoded.Summary.TotalOperations != inventory.Summary.TotalOperations || len(decoded.Operations) != len(inventory.Operations) {
		t.Fatalf("decoded inventory summary = %#v operations=%d", decoded.Summary, len(decoded.Operations))
	}
}

func findOperation(t *testing.T, inventory StripeAPIInventory, method string, path string) StripeOperationCoverage {
	t.Helper()
	for _, operation := range inventory.Operations {
		if operation.Method == method && operation.Path == path {
			return operation
		}
	}
	t.Fatalf("operation %s %s not found in inventory", method, path)
	return StripeOperationCoverage{}
}

func findFamily(t *testing.T, inventory StripeAPIInventory, family string) FamilyCoverage {
	t.Helper()
	for _, item := range inventory.Summary.Families {
		if item.Family == family {
			return item
		}
	}
	t.Fatalf("family %s not found in inventory summary", family)
	return FamilyCoverage{}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
