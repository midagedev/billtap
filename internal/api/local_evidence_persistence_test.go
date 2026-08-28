package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/diagnostics"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
	"github.com/hckim/billtap/internal/webhooks/webhookstest"
)

// handlerOnStore builds a handler over an already-open store, which is how a
// restart is expressed in a test: same store, second handler.
func handlerOnStore(t *testing.T, store *storage.SQLiteStore) http.Handler {
	t.Helper()
	webhookService := webhooks.NewService(store)
	webhookstest.RegisterStoreCleanup(t, webhookService, store)
	return New(Options{
		Billing:       billing.NewService(store),
		Webhooks:      webhookService,
		Diagnostics:   diagnostics.NewService(store),
		LocalEvidence: store,
	})
}

// A restart used to drop every evidence object while the rest of the store
// survived, so lookups kept returning 200 with the sub-record missing.
func TestLocalEvidenceSurvivesHandlerRestart(t *testing.T) {
	store, err := storage.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	first := handlerOnStore(t, store)

	// Path 1: the API. Path 2: fixture apply, which writes explicit IDs and is the
	// path a seeded environment actually uses.
	created := postForm[map[string]any](t, first, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	apiID, _ := created["id"].(string)
	if apiID == "" {
		t.Fatalf("tax rate create returned no id: %#v", created)
	}
	postJSON[map[string]any](t, first, "/api/fixtures/apply", map[string]any{
		"tax_rates": []map[string]any{
			{"id": "txr_fixture_vat", "display_name": "Fixture VAT", "percentage": 10},
		},
	})
	postForm[map[string]any](t, first, "/v1/coupons", url.Values{
		"id":          {"cpn_keepme"},
		"percent_off": {"25"},
		"duration":    {"forever"},
	})

	second := handlerOnStore(t, store)
	for _, id := range []string{apiID, "txr_fixture_vat"} {
		got := getJSON[map[string]any](t, second, "/v1/tax_rates/"+id)
		if got["id"] != id {
			t.Fatalf("tax rate %s missing after restart: %#v", id, got)
		}
	}
	if got := getJSON[map[string]any](t, second, "/v1/coupons/cpn_keepme"); got["id"] != "cpn_keepme" {
		t.Fatalf("coupon missing after restart: %#v", got)
	}

	list := getJSON[map[string]any](t, second, "/v1/tax_rates")
	data, _ := list["data"].([]any)
	if len(data) != 2 {
		t.Fatalf("tax rate list after restart = %d rows, want 2", len(data))
	}
}

// Deletes must persist too, or a restart resurrects the object.
func TestLocalEvidenceDeleteSurvivesHandlerRestart(t *testing.T) {
	store, err := storage.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	first := handlerOnStore(t, store)
	postForm[map[string]any](t, first, "/v1/coupons", url.Values{
		"id":          {"cpn_gone"},
		"percent_off": {"25"},
		"duration":    {"forever"},
	})
	deleteJSON[map[string]any](t, first, "/v1/coupons/cpn_gone")

	second := handlerOnStore(t, store)
	req := httptest.NewRequest(http.MethodGet, "/v1/coupons/cpn_gone", nil)
	rec := httptest.NewRecorder()
	second.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("deleted coupon came back after restart: status %d", rec.Code)
	}
}

// A store that cannot hold evidence keeps the old in-memory behavior.
func TestLocalEvidenceWithoutRepositoryStaysInMemory(t *testing.T) {
	s := newLocalEvidenceStore(nil)
	if err := s.save(kindTaxRate, "txr_1", map[string]any{"id": "txr_1"}); err != nil {
		t.Fatalf("save without repo: %v", err)
	}
	if _, ok := s.taxRates["txr_1"]; !ok {
		t.Fatal("in-memory save did not record the object")
	}
}
