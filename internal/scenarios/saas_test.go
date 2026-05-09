package scenarios

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestSaaSExampleScenariosPass(t *testing.T) {
	assertions := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode assertion payload: %v", err)
		}
		if payload["target"] == "" || payload["context"] == nil {
			t.Fatalf("assertion payload = %#v, want target and context", payload)
		}
		_, _ = w.Write([]byte(`{"pass":true}`))
	}))
	defer assertions.Close()

	for _, name := range []string{
		"saas-adoption-contract.yml",
		"saas-connect-webhook.yml",
		"saas-extra-export.yml",
		"saas-seat-and-member.yml",
		"saas-workspace-upgrade.yml",
	} {
		t.Run(name, func(t *testing.T) {
			scenario, err := LoadFile(filepath.Join("..", "..", "examples", name))
			if err != nil {
				t.Fatalf("LoadFile: %v", err)
			}
			scenario.App.Assertions.BaseURL = assertions.URL
			report, err := NewRunner(nil, nil).Run(context.Background(), scenario)
			if err != nil {
				t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
			}
			if report.ExitCode() != ExitPass {
				t.Fatalf("ExitCode = %d, want pass", report.ExitCode())
			}
		})
	}
}

func TestSaaSSupportBundleContainsAdoptionEvidence(t *testing.T) {
	assertions := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"pass":true}`))
	}))
	defer assertions.Close()

	scenario, err := LoadFile(filepath.Join("..", "..", "examples", "saas-adoption-contract.yml"))
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	scenario.App.Assertions.BaseURL = assertions.URL

	report, err := NewRunner(nil, nil).Run(context.Background(), scenario)
	if err != nil {
		t.Fatalf("Run returned error: %v\n%s", err, report.Markdown())
	}
	var bundle map[string]any
	for _, step := range report.Steps {
		if step.ID == "support-bundle" {
			bundle, _ = step.Output["supportBundle"].(map[string]any)
		}
	}
	if bundle == nil {
		t.Fatal("support bundle step missing")
	}
	for _, key := range []string{"workspace", "subscription", "seats", "exportSummary", "payments", "webhooks", "workspaceLogs"} {
		if bundle[key] == nil {
			t.Fatalf("support bundle missing %s: %#v", key, bundle)
		}
	}
	appAssertions, _ := bundle["appAssertions"].([]map[string]any)
	if len(appAssertions) == 0 {
		t.Fatalf("support bundle missing app assertion evidence: %#v", bundle)
	}
}
