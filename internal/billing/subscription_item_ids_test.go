package billing

import (
	"testing"
)

func TestAssignSubscriptionItemIDs_EmptyAndNil(t *testing.T) {
	if got := AssignSubscriptionItemIDs("sub_x", nil); got != nil {
		t.Fatalf("nil items = %#v, want nil", got)
	}
	if got := AssignSubscriptionItemIDs("sub_x", []LineItem{}); len(got) != 0 {
		t.Fatalf("empty items = %#v, want empty", got)
	}
}

func TestAssignSubscriptionItemIDs_AssignsLowestUnused(t *testing.T) {
	items := []LineItem{
		{PriceID: "price_a", Quantity: 1},
		{PriceID: "price_b", Quantity: 2},
	}
	got := AssignSubscriptionItemIDs("sub_abc", items)
	if got[0].ID != FormatSubscriptionItemID("sub_abc", 0) {
		t.Fatalf("item0 id = %q, want index 0", got[0].ID)
	}
	if got[1].ID != FormatSubscriptionItemID("sub_abc", 1) {
		t.Fatalf("item1 id = %q, want index 1", got[1].ID)
	}
	// Original slice must not be mutated.
	if items[0].ID != "" || items[1].ID != "" {
		t.Fatalf("input mutated: %#v", items)
	}
}

func TestAssignSubscriptionItemIDs_PreservesExisting(t *testing.T) {
	existing0 := FormatSubscriptionItemID("sub_abc", 0)
	existing2 := FormatSubscriptionItemID("sub_abc", 2)
	items := []LineItem{
		{ID: existing0, PriceID: "price_a", Quantity: 1},
		{PriceID: "price_b", Quantity: 1}, // should take unused min index = 1
		{ID: existing2, PriceID: "price_c", Quantity: 1},
	}
	got := AssignSubscriptionItemIDs("sub_abc", items)
	if got[0].ID != existing0 {
		t.Fatalf("preserved id0 = %q, want %q", got[0].ID, existing0)
	}
	if got[1].ID != FormatSubscriptionItemID("sub_abc", 1) {
		t.Fatalf("new id = %q, want unused index 1", got[1].ID)
	}
	if got[2].ID != existing2 {
		t.Fatalf("preserved id2 = %q, want %q", got[2].ID, existing2)
	}
}

func TestAssignSubscriptionItemIDs_ReusesDeletedIndex(t *testing.T) {
	// After deleting middle (index 1), remaining are 0 and 2; next empty gets 1.
	items := []LineItem{
		{ID: FormatSubscriptionItemID("sub_abc", 0), PriceID: "price_a", Quantity: 1},
		{ID: FormatSubscriptionItemID("sub_abc", 2), PriceID: "price_c", Quantity: 1},
		{PriceID: "price_new", Quantity: 1},
	}
	got := AssignSubscriptionItemIDs("sub_abc", items)
	if got[2].ID != FormatSubscriptionItemID("sub_abc", 1) {
		t.Fatalf("reused id = %q, want index 1", got[2].ID)
	}
	if got[0].ID != FormatSubscriptionItemID("sub_abc", 0) || got[1].ID != FormatSubscriptionItemID("sub_abc", 2) {
		t.Fatalf("existing ids changed: %#v", got)
	}
}

func TestFormatSubscriptionItemID_MatchesHistoricalShape(t *testing.T) {
	got := FormatSubscriptionItemID("sub_ff7ddbe0c6b176b5", 2)
	if got != "si_sub_ff7ddbe0c6b176b5_2" {
		t.Fatalf("id = %q, want si_sub_ff7ddbe0c6b176b5_2", got)
	}
}

func TestResolvedSubscriptionItemID_LegacyFallback(t *testing.T) {
	item := LineItem{PriceID: "price_a", Quantity: 1}
	if got := ResolvedSubscriptionItemID("sub_x", item, 3); got != FormatSubscriptionItemID("sub_x", 3) {
		t.Fatalf("legacy resolved = %q", got)
	}
	item.ID = "si_custom"
	if got := ResolvedSubscriptionItemID("sub_x", item, 3); got != "si_custom" {
		t.Fatalf("stored resolved = %q", got)
	}
}
