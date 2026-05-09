package compatibility

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteArtifactsGeneratesPassingScorecard(t *testing.T) {
	dir := t.TempDir()
	scorecard, err := WriteArtifacts(context.Background(), Options{
		OutputDir: dir,
		Now:       fixedNow,
	})
	if err != nil {
		t.Fatalf("WriteArtifacts returned error: %v", err)
	}
	if !scorecard.Summary.Passed {
		t.Fatalf("scorecard did not pass: %#v", scorecard.Summary)
	}
	if scorecard.Summary.Imported == 0 || scorecard.Summary.Skipped == 0 || scorecard.Summary.Unsupported == 0 {
		t.Fatalf("summary = %#v, want imported, skipped, and unsupported statuses", scorecard.Summary)
	}
	if scorecard.Summary.Mismatch != 0 || scorecard.Summary.Error != 0 {
		t.Fatalf("summary = %#v, want no mismatches or errors", scorecard.Summary)
	}
	if _, ok := scorecard.StatusMap[string(StatusImported)]; !ok {
		t.Fatalf("status_map missing imported: %#v", scorecard.StatusMap)
	}

	jsonPath := filepath.Join(dir, "compatibility-scorecard.json")
	mdPath := filepath.Join(dir, "compatibility-scorecard.md")
	if !fileContains(t, jsonPath, `"scorecard_version": "l3-foundation-v1"`) {
		t.Fatalf("JSON scorecard missing version")
	}
	if !fileContains(t, jsonPath, `"mismatch": 0`) || !fileContains(t, jsonPath, `"error": 0`) {
		t.Fatalf("JSON scorecard missing zero mismatch/error counts")
	}
	if !fileContains(t, mdPath, "# Compatibility Scorecard") || !fileContains(t, mdPath, "`unsupported`") {
		t.Fatalf("Markdown scorecard missing expected content")
	}

	entries, err := os.ReadDir(filepath.Join(dir, "replay-bundles"))
	if err != nil {
		t.Fatalf("read replay-bundles: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("replay bundles = %d, want none for passing scorecard", len(entries))
	}
}

func TestMismatchWritesReplayBundle(t *testing.T) {
	dir := t.TempDir()
	scorecard, err := WriteArtifacts(context.Background(), Options{
		OutputDir: dir,
		Now:       fixedNow,
		cases: []caseSpec{{
			ID:              "products.create.intentional_mismatch",
			Name:            "Intentional mismatch",
			Category:        "request-validation",
			Level:           "L3",
			ReleaseBlocking: true,
			Steps: []requestSpec{{
				Name:   "create product",
				Method: http.MethodPost,
				Path:   "/v1/products",
				Params: map[string]string{"name": "Mismatch Product"},
			}},
			Expect: Observation{HTTPStatus: http.StatusCreated, Object: "product"},
		}},
	})
	if err != nil {
		t.Fatalf("WriteArtifacts returned error: %v", err)
	}
	if scorecard.Summary.Mismatch != 1 || scorecard.Summary.Passed {
		t.Fatalf("summary = %#v, want one blocking mismatch", scorecard.Summary)
	}
	if len(scorecard.Cases) != 1 || scorecard.Cases[0].ReplayBundle == "" {
		t.Fatalf("case result missing replay bundle: %#v", scorecard.Cases)
	}

	bundlePath := filepath.Join(dir, scorecard.Cases[0].ReplayBundle)
	body, err := os.ReadFile(bundlePath)
	if err != nil {
		t.Fatalf("read replay bundle: %v", err)
	}
	var bundle ReplayBundle
	if err := json.Unmarshal(body, &bundle); err != nil {
		t.Fatalf("decode replay bundle: %v", err)
	}
	if bundle.CaseID != "products.create.intentional_mismatch" || bundle.Status != StatusMismatch {
		t.Fatalf("bundle = %#v, want mismatch bundle", bundle)
	}
	if len(bundle.Steps) != 1 || bundle.Steps[0].Method != http.MethodPost || bundle.Steps[0].Path != "/v1/products" {
		t.Fatalf("bundle steps = %#v, want replayable request data", bundle.Steps)
	}
	if bundle.Steps[0].Expected.HTTPStatus != http.StatusCreated || bundle.Steps[0].Actual.HTTPStatus != http.StatusOK {
		t.Fatalf("bundle observations = expected %#v actual %#v", bundle.Steps[0].Expected, bundle.Steps[0].Actual)
	}
}

func fixedNow() time.Time {
	return time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
}

func fileContains(t *testing.T, path string, want string) bool {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return strings.Contains(string(body), want)
}
