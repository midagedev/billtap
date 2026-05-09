package scenarios

import (
	"errors"
	"testing"
	"time"
)

func TestLoadParsesScenarioSchema(t *testing.T) {
	scenario, err := Load([]byte(`
name: subscription-payment-retry
app:
  assertions:
    baseUrl: http://localhost:3000/test/assertions
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
  - id: assert-active
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
`))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if scenario.Name != "subscription-payment-retry" {
		t.Fatalf("scenario.Name = %q", scenario.Name)
	}
	if got := scenario.Catalog.Prices[0].UnitAmount; got != 4900 {
		t.Fatalf("unitAmount = %d, want 4900", got)
	}
}

func TestLoadRejectsInvalidConfig(t *testing.T) {
	_, err := Load([]byte(`
name: ""
steps:
  - id: bad
    action: app.assert
    params:
      target: workspace.subscription
`))
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("Load error = %v, want ErrInvalidConfig", err)
	}
}

func TestClockAdvancesDeterministically(t *testing.T) {
	clock, err := NewClock("2026-05-08T00:00:00Z")
	if err != nil {
		t.Fatalf("NewClock: %v", err)
	}
	if _, err := clock.Advance("3d"); err != nil {
		t.Fatalf("advance 3d: %v", err)
	}
	if _, err := clock.Advance("2h30m"); err != nil {
		t.Fatalf("advance 2h30m: %v", err)
	}
	want := time.Date(2026, 5, 11, 2, 30, 0, 0, time.UTC)
	if !clock.Now().Equal(want) {
		t.Fatalf("clock.Now = %s, want %s", clock.Now(), want)
	}
}
