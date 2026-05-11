package fixtures

import "testing"

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
