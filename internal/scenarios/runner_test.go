package scenarios

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/storage"
)

func TestRunnerResolvesReferencesAndCheckoutAliases(t *testing.T) {
	runner := NewRunner(newTestBilling(t), nil)
	scenario := mustLoad(t, `
name: checkout-aliases
catalog:
  products:
    - id: prod_pro
      name: Pro
  prices:
    - id: price_pro_monthly
      product: prod_pro
      currency: usd
      unitAmount: 4900
      interval: month
clock:
  start: "2026-05-08T00:00:00Z"
steps:
  - id: create-customer
    action: customer.create
    params:
      email: user@example.test
  - id: checkout
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly
  - id: complete-checkout
    action: checkout.complete
    params:
      sessionRef: checkout.session.id
      outcome: payment_succeeded
  - id: retry-payment
    action: invoice.retry
    params:
      subscriptionRef: checkout.subscription.id
      invoiceRef: complete-checkout.invoice.id
      outcome: payment_succeeded
`)
	report, err := runner.Run(context.Background(), scenario)
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	if report.ExitCode() != ExitPass {
		t.Fatalf("ExitCode = %d, want pass", report.ExitCode())
	}
	if got := len(report.Steps); got != 4 {
		t.Fatalf("steps = %d, want 4", got)
	}
	retry := report.Steps[3].Output
	if retry["subscription"] == "" || retry["invoice"] == "" {
		t.Fatalf("retry output = %#v, want resolved subscription and invoice refs", retry)
	}
}

func TestRunnerAppAssertPass(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode assertion payload: %v", err)
		}
		if payload["target"] != "workspace.subscription" || payload["expected"] == nil || payload["context"] == nil || payload["clock"] == nil {
			t.Fatalf("payload = %#v, want assertion context", payload)
		}
		_, _ = w.Write([]byte(`{"pass":true}`))
	}))
	defer server.Close()

	report, err := NewRunner(nil, nil).Run(context.Background(), mustLoad(t, `
name: assertion-pass
app:
  assertions:
    baseUrl: `+server.URL+`/assertions
steps:
  - id: assert-subscription
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`))
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if gotPath != "/assertions/workspace/subscription" {
		t.Fatalf("callback path = %q, want slash-converted target", gotPath)
	}
	if report.ExitCode() != ExitPass {
		t.Fatalf("ExitCode = %d, want 0", report.ExitCode())
	}
}

func TestRunnerAppAssertFailureExitCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"pass":false,"message":"status was past_due"}`))
	}))
	defer server.Close()

	report, err := NewRunner(nil, nil).Run(context.Background(), mustLoad(t, `
name: assertion-fail
app:
  assertions:
    baseUrl: `+server.URL+`
steps:
  - id: assert-subscription
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`))
	if !errors.Is(err, ErrAssertionFailed) {
		t.Fatalf("Run error = %v, want ErrAssertionFailed", err)
	}
	if report.ExitCode() != ExitAssertionFailed {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitAssertionFailed)
	}
	if !strings.Contains(report.Markdown(), "status was past_due") {
		t.Fatalf("Markdown missing assertion detail:\n%s", report.Markdown())
	}
}

func TestRunnerAppCallbackErrorExitCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	report, err := NewRunner(nil, nil).Run(context.Background(), mustLoad(t, `
name: callback-error
app:
  assertions:
    baseUrl: `+server.URL+`
steps:
  - id: assert-subscription
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`))
	if !errors.Is(err, ErrAppCallbackFailure) {
		t.Fatalf("Run error = %v, want ErrAppCallbackFailure", err)
	}
	if report.ExitCode() != ExitAppCallbackFailure {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitAppCallbackFailure)
	}
}

func TestReportJSONMarkdownAndInvalidExitCode(t *testing.T) {
	report := Report{
		Name:        "bad",
		Status:      "failed",
		FailureType: FailureInvalidConfig,
		Errors:      []string{"name is required"},
	}
	body, err := report.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if !strings.Contains(string(body), `"failure_type": "invalid_config"`) {
		t.Fatalf("JSON = %s", body)
	}
	if !strings.Contains(report.Markdown(), "Exit code: `2`") {
		t.Fatalf("Markdown = %s", report.Markdown())
	}
	if report.ExitCode() != ExitInvalidConfig {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitInvalidConfig)
	}
}

func TestReportRuntimeFailureExitCode(t *testing.T) {
	report := Report{Name: "runtime", Status: "failed", FailureType: FailureRunner}
	if report.ExitCode() != ExitRuntimeFailure {
		t.Fatalf("ExitCode = %d, want %d", report.ExitCode(), ExitRuntimeFailure)
	}
	if !strings.Contains(report.Markdown(), "Exit code: `4`") {
		t.Fatalf("Markdown = %s", report.Markdown())
	}
}

func newTestBilling(t *testing.T) *billing.Service {
	t.Helper()
	store, err := storage.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})
	return billing.NewService(store)
}

func mustLoad(t *testing.T, raw string) Scenario {
	t.Helper()
	scenario, err := Load([]byte(raw))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return scenario
}
