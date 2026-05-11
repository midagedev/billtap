package fixtures

import (
	"testing"

	"github.com/hckim/billtap/internal/billing"
)

func TestFixtureOutcomeStatusOverridesSuccessfulOutcome(t *testing.T) {
	tests := []struct {
		name    string
		fixture SubscriptionFixture
		want    string
	}{
		{
			name: "trialing keeps trial seed path",
			fixture: SubscriptionFixture{
				Status:  "trialing",
				Outcome: "payment_succeeded",
			},
			want: "payment_succeeded",
		},
		{
			name: "canceled ignores successful payment outcome",
			fixture: SubscriptionFixture{
				Status:  "canceled",
				Outcome: "payment_succeeded",
			},
			want: "canceled",
		},
		{
			name: "past due ignores successful payment outcome",
			fixture: SubscriptionFixture{
				Status:  "past_due",
				Outcome: "payment_succeeded",
			},
			want: "card_declined",
		},
		{
			name: "unpaid ignores successful payment outcome",
			fixture: SubscriptionFixture{
				Status:  "unpaid",
				Outcome: "payment_succeeded",
			},
			want: "card_declined",
		},
		{
			name: "incomplete ignores successful payment outcome",
			fixture: SubscriptionFixture{
				Status:  "incomplete",
				Outcome: "payment_succeeded",
			},
			want: "payment_pending",
		},
		{
			name: "active keeps explicit failed outcome for payment evidence",
			fixture: SubscriptionFixture{
				Status:  "active",
				Outcome: "card_declined",
			},
			want: "card_declined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixtureOutcome(tt.fixture); got != tt.want {
				t.Fatalf("fixtureOutcome() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCustomerPaymentMethodFixtureConfig(t *testing.T) {
	yamlPack := []byte(`
customers:
  - id: cus_fixture_empty_methods
    payment_methods_fixture: empty
  - id: cus_fixture_empty_methods_list
    payment_methods: []
  - id: cus_fixture_explicit_methods
    payment_methods:
      - id: pm_fixture_secondary
      - id: pm_fixture_primary
        default: true
`)
	pack, err := LoadPack(yamlPack, "application/yaml")
	if err != nil {
		t.Fatalf("LoadPack() error = %v", err)
	}

	mode, ids, defaultID, configured, err := customerPaymentMethodFixtureConfig(pack.Customers[0])
	if err != nil || !configured || mode != billing.PaymentMethodsFixtureEmpty || len(ids) != 0 || defaultID != "" {
		t.Fatalf("mode fixture config = mode %q ids %#v default %q configured %v err %v, want empty", mode, ids, defaultID, configured, err)
	}

	mode, ids, defaultID, configured, err = customerPaymentMethodFixtureConfig(pack.Customers[1])
	if err != nil || !configured || mode != billing.PaymentMethodsFixtureEmpty || len(ids) != 0 || defaultID != "" {
		t.Fatalf("empty list config = mode %q ids %#v default %q configured %v err %v, want empty", mode, ids, defaultID, configured, err)
	}

	mode, ids, defaultID, configured, err = customerPaymentMethodFixtureConfig(pack.Customers[2])
	if err != nil || !configured || mode != billing.PaymentMethodsFixtureExplicit || len(ids) != 2 || ids[0] != "pm_fixture_secondary" || ids[1] != "pm_fixture_primary" || defaultID != "pm_fixture_primary" {
		t.Fatalf("explicit config = mode %q ids %#v default %q configured %v err %v, want explicit methods with primary default", mode, ids, defaultID, configured, err)
	}
}
