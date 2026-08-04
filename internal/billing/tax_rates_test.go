package billing

import "testing"

func TestComputeTaxRateAmountsExclusiveMatchesExclusiveTaxAmount(t *testing.T) {
	rates := []AppliedTaxRate{{ID: "txr_ex", DisplayName: "Sales", Percentage: 10, Inclusive: false}}
	amounts, pretax, exclusiveTotal, taxTotal := ComputeTaxRateAmounts(6000, rates)
	want := ExclusiveTaxAmount(6000, 10)
	if taxTotal != want || exclusiveTotal != want {
		t.Fatalf("taxTotal/exclusive=%d/%d, want ExclusiveTaxAmount=%d", taxTotal, exclusiveTotal, want)
	}
	if pretax != 6000 {
		t.Fatalf("pretax=%d, want 6000", pretax)
	}
	if len(amounts) != 1 || amounts[0].Amount != want {
		t.Fatalf("amounts=%#v, want one line of %d", amounts, want)
	}
}

func TestComputeTaxRateAmountsInclusive(t *testing.T) {
	// base 1100 inclusive 10% → tax 100, pretax 1000, exclusive 0
	rates := []AppliedTaxRate{{ID: "txr_in", DisplayName: "VAT", Percentage: 10, Inclusive: true}}
	amounts, pretax, exclusiveTotal, taxTotal := ComputeTaxRateAmounts(1100, rates)
	if taxTotal != 100 || pretax != 1000 || exclusiveTotal != 0 {
		t.Fatalf("got pretax=%d exclusive=%d taxTotal=%d, want 1000/0/100", pretax, exclusiveTotal, taxTotal)
	}
	if len(amounts) != 1 || amounts[0].Amount != 100 {
		t.Fatalf("amounts=%#v", amounts)
	}
}

func TestComputeTaxRateAmountsMixed(t *testing.T) {
	// base 1100, inclusive 10% + exclusive 5% → incl 100, pretax 1000, excl 50, taxTotal 150
	rates := []AppliedTaxRate{
		{ID: "txr_in", DisplayName: "VAT", Percentage: 10, Inclusive: true},
		{ID: "txr_ex", DisplayName: "Local", Percentage: 5, Inclusive: false},
	}
	amounts, pretax, exclusiveTotal, taxTotal := ComputeTaxRateAmounts(1100, rates)
	if pretax != 1000 || exclusiveTotal != 50 || taxTotal != 150 {
		t.Fatalf("got pretax=%d exclusive=%d taxTotal=%d, want 1000/50/150", pretax, exclusiveTotal, taxTotal)
	}
	if len(amounts) != 2 || amounts[0].Amount != 100 || amounts[1].Amount != 50 {
		t.Fatalf("amounts=%#v", amounts)
	}
}

func TestComputeTaxRateAmountsDecimalPercent(t *testing.T) {
	// 7.25% of 10000 = 725
	rates := []AppliedTaxRate{{ID: "txr_dec", DisplayName: "CA", Percentage: 7.25, Inclusive: false}}
	_, _, exclusiveTotal, taxTotal := ComputeTaxRateAmounts(10000, rates)
	want := ExclusiveTaxAmount(10000, 7.25)
	if taxTotal != want || exclusiveTotal != want || want != 725 {
		t.Fatalf("7.25%% of 10000 = taxTotal=%d exclusive=%d ExclusiveTaxAmount=%d, want 725", taxTotal, exclusiveTotal, want)
	}
}

func TestDefaultTaxRatesMetadataRoundTrip(t *testing.T) {
	rates := []AppliedTaxRate{
		{ID: "txr_1", DisplayName: "VAT", Percentage: 10, Inclusive: false},
	}
	meta := MergeDefaultTaxRatesMetadata(nil, rates)
	got := DefaultTaxRatesFromMetadata(meta)
	if len(got) != 1 || got[0].ID != "txr_1" || got[0].Percentage != 10 || got[0].Inclusive {
		t.Fatalf("round-trip = %#v", got)
	}
	cleared := MergeDefaultTaxRatesMetadata(meta, nil)
	if DefaultTaxRatesFromMetadata(cleared) != nil {
		t.Fatalf("clear left residual = %#v", cleared)
	}
}
