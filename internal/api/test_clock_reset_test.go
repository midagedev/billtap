package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/fixtures"
)

func TestTestClockDeleteDetachesReferences(t *testing.T) {
	handler := newTestHandler(t)
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":      {"clock-delete@example.test"},
		"test_clock": {"clock_detached"},
	})
	base := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_detached"},
		"frozen_time": {fmt.Sprint(base.Unix())},
	})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Clock Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1500"},
		"recurring[interval]": {"month"},
	})
	subscription := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {price.ID},
		"items[0][quantity]": {"1"},
		"test_clock":         {"clock_detached"},
	})

	deleted := deleteForm[struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_detached", url.Values{})
	if !deleted.Deleted || deleted.Object != "test_clock" {
		t.Fatalf("clock delete = %#v, want deleted test_clock marker", deleted)
	}
	missingReq := httptest.NewRequest(http.MethodGet, "/v1/test_helpers/test_clocks/clock_detached", nil)
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("deleted clock GET status = %d, want 404", missingRec.Code)
	}
	// Re-delete is 404.
	status, body := deleteFormStatus(t, handler, "/v1/test_helpers/test_clocks/clock_detached", url.Values{})
	if status != http.StatusNotFound {
		t.Fatalf("re-delete clock status = %d body = %s, want 404", status, body)
	}
	// Objects survive, detached from the clock.
	subscriptionAfter := getJSON[struct {
		ID       string            `json:"id"`
		Metadata map[string]string `json:"metadata"`
	}](t, handler, "/v1/subscriptions/"+subscription.ID)
	if subscriptionAfter.ID != subscription.ID || subscriptionAfter.Metadata["test_clock"] != "" {
		t.Fatalf("subscription after clock delete = %#v, want kept but detached", subscriptionAfter)
	}
	customerAfter := getJSON[billing.Customer](t, handler, "/v1/customers/"+customer.ID)
	if customerAfter.ID != customer.ID || customerAfter.Metadata["test_clock"] != "" {
		t.Fatalf("customer after clock delete = %#v, want kept but detached", customerAfter)
	}
}

// The issue #76 scenario: advancing a trial clock consumes the pack's premise;
// re-applying the same pack must restore it.
func TestFixtureReapplyRestoresAdvancedTrialClock(t *testing.T) {
	handler := newTestHandler(t)
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Trial Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1200"},
		"recurring[interval]": {"month"},
	})

	base := time.Date(2030, 3, 1, 0, 0, 0, 0, time.UTC)
	trialEnd := base.Add(14 * 24 * time.Hour)
	advanced := trialEnd.Add(60 * time.Second)
	pack := map[string]any{
		"name": "trial-clock-reset",
		"test_clocks": []map[string]any{{
			"id":          "clock_trial_reset",
			"frozen_time": fmt.Sprint(base.Unix()),
		}},
		"customers": []map[string]any{{
			"id":         "cus_clock_trial_reset",
			"email":      "trial-reset@example.test",
			"test_clock": "clock_trial_reset",
		}},
		"subscriptions": []map[string]any{{
			"ref":                  "trial_sub",
			"customer":             "cus_clock_trial_reset",
			"price":                price.ID,
			"outcome":              "payment_succeeded",
			"status":               "trialing",
			"test_clock":           "clock_trial_reset",
			"current_period_start": fmt.Sprint(base.Unix()),
			"current_period_end":   fmt.Sprint(trialEnd.Unix()),
			"trial_start":          fmt.Sprint(base.Unix()),
			"trial_end":            fmt.Sprint(trialEnd.Unix()),
		}},
	}

	applyFixturePack(t, handler, pack)
	first := getJSON[subscriptionStatusList](t, handler, "/v1/subscriptions?customer=cus_clock_trial_reset")
	if len(first.Data) != 1 || first.Data[0].Status != "trialing" {
		t.Fatalf("first apply statuses = %#v, want one trialing subscription", first.Data)
	}

	// Advance past trial_end: the trial premise is consumed (issue #76).
	postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks/clock_trial_reset/advance", url.Values{
		"frozen_time": {fmt.Sprint(advanced.Unix())},
	})
	afterAdvance := getJSON[subscriptionStatusList](t, handler, "/v1/subscriptions?customer=cus_clock_trial_reset")
	if len(afterAdvance.Data) != 1 || afterAdvance.Data[0].Status != "active" {
		t.Fatalf("after advance statuses = %#v, want one active subscription", afterAdvance.Data)
	}
	// Backwards advance stays rejected.
	status, body := postFormStatus(t, handler, "/v1/test_helpers/test_clocks/clock_trial_reset/advance", url.Values{
		"frozen_time": {fmt.Sprint(base.Unix())},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("backwards advance status = %d body = %s, want 400", status, body)
	}

	// Re-applying the pack recreates the clock at the declared time and
	// restores the trialing premise.
	applyFixturePack(t, handler, pack)
	clock := getJSON[struct {
		ID         string `json:"id"`
		FrozenTime int64  `json:"frozen_time"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_trial_reset")
	if clock.ID != "clock_trial_reset" || clock.FrozenTime != base.Unix() {
		t.Fatalf("clock after re-apply = %#v, want recreated at %d", clock, base.Unix())
	}
	restored := getJSON[subscriptionStatusList](t, handler, "/v1/subscriptions?customer=cus_clock_trial_reset")
	if len(restored.Data) != 1 || restored.Data[0].Status != "trialing" {
		t.Fatalf("restored statuses = %#v, want one trialing subscription", restored.Data)
	}
}

type subscriptionStatusList struct {
	Data []struct {
		Status string `json:"status"`
	} `json:"data"`
}

func applyFixturePack(t *testing.T, handler http.Handler, pack map[string]any) fixtures.ApplyResult {
	t.Helper()
	status, body := postJSONStatus(t, handler, "/api/fixtures/apply", pack)
	if status != http.StatusOK {
		t.Fatalf("fixture apply status = %d body = %s", status, body)
	}
	var result fixtures.ApplyResult
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatalf("decode fixture apply result: %v", err)
	}
	return result
}
