package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hckim/billtap/internal/scenarios"
)

func TestParseScenarioRunArgsAllowsFlagsAfterFile(t *testing.T) {
	path, reportJSON, reportMD, databaseURL, err := parseScenarioRunArgs([]string{
		"scenario.yml",
		"--report-json",
		"report.json",
		"--report-md=report.md",
		"--database-url",
		"state.db",
	})
	if err != nil {
		t.Fatalf("parseScenarioRunArgs returned error: %v", err)
	}
	if path != "scenario.yml" || reportJSON != "report.json" || reportMD != "report.md" || databaseURL != "state.db" {
		t.Fatalf("path=%q reportJSON=%q reportMD=%q databaseURL=%q", path, reportJSON, reportMD, databaseURL)
	}
}

func TestRunScenarioWritesReportsAndReturnsPass(t *testing.T) {
	assertions := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/assertions/workspace/subscription" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pass":true}`))
	}))
	defer assertions.Close()

	dir := t.TempDir()
	scenarioPath := filepath.Join(dir, "scenario.yml")
	reportJSON := filepath.Join(dir, "report.json")
	reportMD := filepath.Join(dir, "report.md")
	writeFile(t, scenarioPath, `
name: cli-pass
app:
  assertions:
    baseUrl: `+assertions.URL+`/assertions
steps:
  - id: assert-active
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`)
	code := runScenario([]string{"run", scenarioPath, "--report-json", reportJSON, "--report-md", reportMD, "--database-url", filepath.Join(dir, "state.db")})
	if code != scenarios.ExitPass {
		t.Fatalf("exit code = %d, want %d", code, scenarios.ExitPass)
	}
	if !fileContains(t, reportJSON, `"status": "passed"`) {
		t.Fatalf("JSON report missing passed status")
	}
	if !fileContains(t, reportMD, "Exit code: `0`") {
		t.Fatalf("Markdown report missing exit code")
	}
}

func TestRunCompatibilityWritesScorecard(t *testing.T) {
	dir := t.TempDir()
	code := runCompatibility([]string{"scorecard", "--output-dir", dir})
	if code != scenarios.ExitPass {
		t.Fatalf("exit code = %d, want %d", code, scenarios.ExitPass)
	}
	if !fileContains(t, filepath.Join(dir, "compatibility-scorecard.json"), `"mismatch": 0`) {
		t.Fatalf("JSON scorecard missing zero mismatch count")
	}
	if !fileContains(t, filepath.Join(dir, "compatibility-scorecard.json"), `"error": 0`) {
		t.Fatalf("JSON scorecard missing zero error count")
	}
	if !fileContains(t, filepath.Join(dir, "compatibility-scorecard.json"), `"passed": true`) {
		t.Fatalf("JSON scorecard missing passed true")
	}
	if !fileContains(t, filepath.Join(dir, "compatibility-scorecard.json"), `"release_blocking": 30`) {
		t.Fatalf("JSON scorecard missing expected release-blocking count")
	}
	if !fileContains(t, filepath.Join(dir, "compatibility-scorecard.md"), "# Compatibility Scorecard") {
		t.Fatalf("Markdown scorecard missing heading")
	}
}

func TestRunCompatibilityWritesInventory(t *testing.T) {
	dir := t.TempDir()
	openAPIPath := filepath.Join("..", "..", "internal", "compatibility", "testdata", "stripe-openapi-minimal.json")
	code := runCompatibility([]string{"inventory", "--openapi", openAPIPath, "--output-dir", dir, "--source", "stripe/openapi test fixture"})
	if code != scenarios.ExitPass {
		t.Fatalf("exit code = %d, want %d", code, scenarios.ExitPass)
	}
	if !fileContains(t, filepath.Join(dir, "stripe-api-inventory.json"), `"inventory_version": "stripe-api-inventory-v1"`) {
		t.Fatalf("JSON inventory missing version")
	}
	if !fileContains(t, filepath.Join(dir, "stripe-api-inventory.json"), `"implemented_operations": 7`) {
		t.Fatalf("JSON inventory missing implemented operation count")
	}
	if !fileContains(t, filepath.Join(dir, "stripe-api-inventory.md"), "# Stripe API Compatibility Inventory") {
		t.Fatalf("Markdown inventory missing heading")
	}
}

func TestRunCompatibilityInventoryRequiresOpenAPIPath(t *testing.T) {
	code := runCompatibility([]string{"inventory", "--output-dir", t.TempDir()})
	if code != scenarios.ExitInvalidConfig {
		t.Fatalf("exit code = %d, want %d", code, scenarios.ExitInvalidConfig)
	}
}

func TestRunScenarioReturnsInvalidConfigExitCode(t *testing.T) {
	dir := t.TempDir()
	scenarioPath := filepath.Join(dir, "bad.yml")
	reportJSON := filepath.Join(dir, "report.json")
	writeFile(t, scenarioPath, `
name: ""
steps: []
`)
	code := runScenario([]string{"run", scenarioPath, "--report-json", reportJSON})
	if code != scenarios.ExitInvalidConfig {
		t.Fatalf("exit code = %d, want %d", code, scenarios.ExitInvalidConfig)
	}
	if !fileContains(t, reportJSON, `"failure_type": "invalid_config"`) {
		t.Fatalf("JSON report missing invalid_config failure")
	}
}

func TestRunScenarioReturnsAppCallbackFailureExitCode(t *testing.T) {
	assertions := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusServiceUnavailable)
	}))
	defer assertions.Close()

	dir := t.TempDir()
	scenarioPath := filepath.Join(dir, "callback.yml")
	writeFile(t, scenarioPath, `
name: cli-callback
app:
  assertions:
    baseUrl: `+assertions.URL+`
steps:
  - id: assert-active
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`)
	code := runScenario([]string{"run", scenarioPath, "--database-url", filepath.Join(dir, "state.db")})
	if code != scenarios.ExitAppCallbackFailure {
		t.Fatalf("exit code = %d, want %d", code, scenarios.ExitAppCallbackFailure)
	}
}

func writeFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func fileContains(t *testing.T, path string, want string) bool {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return strings.Contains(string(body), want)
}
