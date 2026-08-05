package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
)

// Seat add/remove via POST|DELETE /v1/subscription_items with proration_behavior.
// Amounts are integer cents; VAT 10% exclusive. Item IDs are stored at creation (stable).

type subItemResponse struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Quantity int64  `json:"quantity"`
	Price    struct {
		ID string `json:"id"`
	} `json:"price"`
	TaxRates []struct {
		ID string `json:"id"`
	} `json:"tax_rates"`
	Metadata map[string]string `json:"metadata"`
}

func setupSeatPlans(t *testing.T, handler http.Handler) (customer billing.Customer, base, seat billing.Price, taxRateID string) {
	t.Helper()
	customer = postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"seats@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Seat Plan"}})
	base = postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"9900"},
		"recurring[interval]": {"month"},
	})
	seat = postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"1000"},
		"recurring[interval]": {"month"},
	})
	txr := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"VAT"},
		"percentage":   {"10"},
		"inclusive":    {"false"},
	})
	return customer, base, seat, txr.ID
}

// 1. Seat add + always_invoice: 9900 base + VAT10%, add 1000×2 → delta 2000 / tax 200 / total 2200.
func TestSubscriptionItemCreateAlwaysInvoiceSeatAdd(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, taxRateID := setupSeatPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {base.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})
	createInvoiceID := created.LatestInvoice

	item := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"2"},
		"proration_behavior": {"always_invoice"},
		"proration_date":     {strconv.FormatInt(created.CurrentPeriodStart, 10)},
	})
	if item.Object != "subscription_item" || item.Quantity != 2 || item.Price.ID != seat.ID {
		t.Fatalf("created item = %#v", item)
	}

	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if sub.LatestInvoice == "" || sub.LatestInvoice == createInvoiceID {
		t.Fatalf("latest_invoice = %q, want new proration invoice (create was %q)", sub.LatestInvoice, createInvoiceID)
	}
	if len(sub.Items.Data) != 2 {
		t.Fatalf("items len = %d, want 2", len(sub.Items.Data))
	}

	invoice := getJSON[struct {
		prorationInvoiceResponse
		TotalExcludingTax    int64 `json:"total_excluding_tax"`
		TotalDiscountAmounts []struct {
			Amount int64 `json:"amount"`
		} `json:"total_discount_amounts"`
		TotalTaxes []struct {
			Amount        int64 `json:"amount"`
			TaxableAmount int64 `json:"taxable_amount"`
		} `json:"total_taxes"`
	}](t, handler, "/v1/invoices/"+sub.LatestInvoice)
	// Full remaining period: delta 2000, tax 200, total 2200
	if invoice.Subtotal != 2000 || invoice.Tax != 200 || invoice.Total != 2200 || invoice.AmountPaid != 2200 {
		t.Fatalf("seat proration invoice = %#v, want subtotal=2000 tax=200 total=2200", invoice)
	}
	if invoice.BillingReason != "subscription_update" {
		t.Fatalf("billing_reason = %q, want subscription_update", invoice.BillingReason)
	}
	assertInvoiceAmountIdentity(t, invoice.Subtotal, invoice.TotalDiscountAmounts, invoice.Tax, invoice.Total, invoice.TotalExcludingTax, invoice.TotalTaxes)
}

// 2. Seat add + create_prorations: no invoice; pending accumulates; renewal absorbs it.
func TestSubscriptionItemCreateCreateProrationsDefersToRenewal(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, taxRateID := setupSeatPlans(t, handler)

	// Fixed clock so proration_date mid-period and renewal advance are deterministic.
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_seat_create_prorations"},
		"frozen_time": {strconv.FormatInt(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {base.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
		"test_clock":           {"clock_seat_create_prorations"},
	})
	periodStart := created.CurrentPeriodStart
	periodEnd := created.CurrentPeriodEnd
	// Mid-period pin for remaining fraction.
	mid := periodStart + (periodEnd-periodStart)/2
	remaining := periodEnd - mid
	periodSec := periodEnd - periodStart
	// delta = (9900+2000) - 9900 = 2000 prorated
	wantPending := 2000 * remaining / periodSec

	before := listSubscriptionInvoices(t, handler, created.ID)
	_ = postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"2"},
		"proration_behavior": {"create_prorations"},
		"proration_date":     {strconv.FormatInt(mid, 10)},
	})
	after := listSubscriptionInvoices(t, handler, created.ID)
	if len(after) != len(before) {
		t.Fatalf("create_prorations created invoice: before=%d after=%d", len(before), len(after))
	}
	updated := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if updated.Metadata[billing.MetadataPendingProrationAmount] != strconv.FormatInt(wantPending, 10) {
		t.Fatalf("pending = %q, want %d", updated.Metadata[billing.MetadataPendingProrationAmount], wantPending)
	}

	advance := postForm[struct {
		BilltapAdvanceResult struct {
			Renewals []struct {
				Invoice struct {
					Subtotal int64 `json:"subtotal"`
					Total    int64 `json:"total"`
					Tax      int64 `json:"tax"`
				} `json:"invoice"`
			} `json:"renewals"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_seat_create_prorations/advance", url.Values{
		"frozen_time": {strconv.FormatInt(periodEnd, 10)},
	})
	if len(advance.BilltapAdvanceResult.Renewals) != 1 {
		t.Fatalf("renewals = %#v, want 1", advance.BilltapAdvanceResult.Renewals)
	}
	renewal := advance.BilltapAdvanceResult.Renewals[0].Invoice
	// New period base = 9900 + 2000 = 11900 + pending
	wantSubtotal := 11900 + wantPending
	_, _, _, wantTax := billing.ComputeTaxRateAmounts(wantSubtotal, []billing.AppliedTaxRate{{Percentage: 10, Inclusive: false}})
	if renewal.Subtotal != wantSubtotal {
		t.Fatalf("renewal subtotal = %d, want %d", renewal.Subtotal, wantSubtotal)
	}
	if renewal.Tax != wantTax {
		t.Fatalf("renewal tax = %d, want %d", renewal.Tax, wantTax)
	}
	if renewal.Total != wantSubtotal+wantTax {
		t.Fatalf("renewal total = %d, want %d", renewal.Total, wantSubtotal+wantTax)
	}
}

// 3. Seat add + none / unspecified: no invoice.
func TestSubscriptionItemCreateNoneNoInvoice(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {base.ID},
		"items[0][quantity]": {"1"},
	})
	before := listSubscriptionInvoices(t, handler, created.ID)

	_ = postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"2"},
		"proration_behavior": {"none"},
	})
	afterNone := listSubscriptionInvoices(t, handler, created.ID)
	if len(afterNone) != len(before) {
		t.Fatalf("none created invoices: before=%d after=%d", len(before), len(afterNone))
	}

	_ = postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {created.ID},
		"price":        {seat.ID},
		"quantity":     {"1"},
	})
	afterUnspec := listSubscriptionInvoices(t, handler, created.ID)
	if len(afterUnspec) != len(before) {
		t.Fatalf("unspecified created invoices: before=%d after=%d", len(before), len(afterUnspec))
	}
	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if len(sub.Items.Data) != 3 {
		t.Fatalf("items = %d, want 3", len(sub.Items.Data))
	}
	if sub.LatestInvoice != created.LatestInvoice {
		t.Fatalf("latest_invoice changed under none/unspecified")
	}
}

// 4. quantity=0 → 400 parameter_invalid (scorecard case).
func TestSubscriptionItemCreateInvalidQuantity(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {base.ID},
	})
	status, body := postFormStatus(t, handler, "/v1/subscription_items", url.Values{
		"subscription": {created.ID},
		"price":        {seat.ID},
		"quantity":     {"0"},
	})
	errBody := decodeErrorBody(t, body)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", status, body)
	}
	if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "quantity" {
		t.Fatalf("error=%#v, want invalid quantity", errBody.Error)
	}
}

// 5. Invalid proration_behavior enum → 400.
func TestSubscriptionItemCreateInvalidProrationBehavior(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {base.ID},
	})
	status, body := postFormStatus(t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"1"},
		"proration_behavior": {"bogus"},
	})
	errBody := decodeErrorBody(t, body)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", status, body)
	}
	if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "proration_behavior" {
		t.Fatalf("error=%#v, want invalid proration_behavior", errBody.Error)
	}
}

// 6. DELETE + always_invoice on seat: delta ≤ 0 → no invoice; item removed; deleted:true.
func TestSubscriptionItemDeleteAlwaysInvoiceNoInvoiceOnDowngrade(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, taxRateID := setupSeatPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {base.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})
	item := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"2"},
		"proration_behavior": {"none"},
	})
	before := listSubscriptionInvoices(t, handler, created.ID)
	subBefore := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)

	deleted := deleteForm[struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}](t, handler, "/v1/subscription_items/"+item.ID, url.Values{
		"proration_behavior": {"always_invoice"},
		"proration_date":     {strconv.FormatInt(created.CurrentPeriodStart, 10)},
	})
	if deleted.ID != item.ID || !deleted.Deleted || deleted.Object != "subscription_item" {
		t.Fatalf("deleted = %#v", deleted)
	}
	after := listSubscriptionInvoices(t, handler, created.ID)
	if len(after) != len(before) {
		t.Fatalf("always_invoice delete created invoice: before=%d after=%d", len(before), len(after))
	}
	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if len(sub.Items.Data) != 1 {
		t.Fatalf("items after delete = %d, want 1", len(sub.Items.Data))
	}
	if sub.LatestInvoice != subBefore.LatestInvoice {
		t.Fatalf("latest_invoice changed on downgrade delete")
	}
}

// 7. DELETE + create_prorations: negative pending; next renewal deducts it.
func TestSubscriptionItemDeleteCreateProrationsNegativePending(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, taxRateID := setupSeatPlans(t, handler)

	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_seat_delete_prorations"},
		"frozen_time": {strconv.FormatInt(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), 10)},
	})
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {base.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
		"test_clock":           {"clock_seat_delete_prorations"},
	})
	// Append seat without proration so items are base+seat.
	item := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"2"},
		"proration_behavior": {"none"},
	})
	periodStart := created.CurrentPeriodStart
	periodEnd := created.CurrentPeriodEnd
	mid := periodStart + (periodEnd-periodStart)/2
	remaining := periodEnd - mid
	periodSec := periodEnd - periodStart
	// Remove 1000×2 seat → delta = -2000 prorated (negative pending).
	wantPending := -2000 * remaining / periodSec

	_ = deleteForm[map[string]any](t, handler, "/v1/subscription_items/"+item.ID, url.Values{
		"proration_behavior": {"create_prorations"},
		"proration_date":     {strconv.FormatInt(mid, 10)},
	})
	updated := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if updated.Metadata[billing.MetadataPendingProrationAmount] != strconv.FormatInt(wantPending, 10) {
		t.Fatalf("pending = %q, want %d (negative)", updated.Metadata[billing.MetadataPendingProrationAmount], wantPending)
	}
	if wantPending >= 0 {
		t.Fatalf("expected negative pending, got %d", wantPending)
	}

	advance := postForm[struct {
		BilltapAdvanceResult struct {
			Renewals []struct {
				Invoice struct {
					Subtotal int64 `json:"subtotal"`
					Total    int64 `json:"total"`
					Tax      int64 `json:"tax"`
				} `json:"invoice"`
			} `json:"renewals"`
		} `json:"billtap_advance_result"`
	}](t, handler, "/v1/test_helpers/test_clocks/clock_seat_delete_prorations/advance", url.Values{
		"frozen_time": {strconv.FormatInt(periodEnd, 10)},
	})
	if len(advance.BilltapAdvanceResult.Renewals) != 1 {
		t.Fatalf("renewals = %#v, want 1", advance.BilltapAdvanceResult.Renewals)
	}
	renewal := advance.BilltapAdvanceResult.Renewals[0].Invoice
	// After delete: base only 9900 + negative pending
	wantSubtotal := 9900 + wantPending
	_, _, _, wantTax := billing.ComputeTaxRateAmounts(wantSubtotal, []billing.AppliedTaxRate{{Percentage: 10, Inclusive: false}})
	if renewal.Subtotal != wantSubtotal {
		t.Fatalf("renewal subtotal = %d, want %d (9900+pending %d)", renewal.Subtotal, wantSubtotal, wantPending)
	}
	if renewal.Tax != wantTax {
		t.Fatalf("renewal tax = %d, want %d", renewal.Tax, wantTax)
	}
	if renewal.Total != wantSubtotal+wantTax {
		t.Fatalf("renewal total = %d, want %d", renewal.Total, wantSubtotal+wantTax)
	}
}

// 8. Last item delete → 400.
func TestSubscriptionItemDeleteLastItemRejected(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, _, _ := setupSeatPlans(t, handler)
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {base.ID},
	})
	if len(created.Items.Data) != 1 {
		t.Fatalf("items = %#v, want 1", created.Items.Data)
	}
	itemID := created.Items.Data[0].ID
	status, body := deleteFormStatus(t, handler, "/v1/subscription_items/"+itemID, nil)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", status, body)
	}
	errBody := decodeErrorBody(t, body)
	if !strings.Contains(errBody.Error.Message, "last subscription item") {
		t.Fatalf("message=%q, want last subscription item rejection", errBody.Error.Message)
	}
}

// 9. Middle item delete shifts subsequent item IDs (index-based; existing constraint).
func TestSubscriptionItemDeleteKeepsStableIDs(t *testing.T) {
	// 2026-08-05: subscription item IDs are now stored at creation; deleting an item no longer shifts the IDs of later items (was asserted as shifting before).
	handler := newTestHandler(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)
	extra := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {base.ProductID},
		"currency":            {"usd"},
		"unit_amount":         {"500"},
		"recurring[interval]": {"month"},
	})

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {base.ID},
		"items[0][quantity]": {"1"},
	})
	item0 := created.Items.Data[0].ID
	item1 := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {created.ID},
		"price":        {seat.ID},
		"quantity":     {"1"},
	})
	item2 := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {created.ID},
		"price":        {extra.ID},
		"quantity":     {"1"},
	})
	// IDs are si_<sub>_<idx>; capture before middle delete.
	before := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if len(before.Items.Data) != 3 {
		t.Fatalf("items before = %d, want 3", len(before.Items.Data))
	}
	idsBefore := []string{before.Items.Data[0].ID, before.Items.Data[1].ID, before.Items.Data[2].ID}
	if idsBefore[0] != item0 || idsBefore[1] != item1.ID || idsBefore[2] != item2.ID {
		t.Fatalf("ids before = %v, want %s,%s,%s", idsBefore, item0, item1.ID, item2.ID)
	}

	_ = deleteForm[map[string]any](t, handler, "/v1/subscription_items/"+item1.ID, url.Values{
		"proration_behavior": {"none"},
	})
	after := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if len(after.Items.Data) != 2 {
		t.Fatalf("items after = %d, want 2", len(after.Items.Data))
	}
	// Stable stored IDs: remaining items keep their creation-time IDs (no index shift).
	if after.Items.Data[0].ID != item0 {
		t.Fatalf("item0 id = %s, want stable %s", after.Items.Data[0].ID, item0)
	}
	if after.Items.Data[1].ID != item2.ID {
		t.Fatalf("remaining item id = %s, want stable former last id %s (no shift)", after.Items.Data[1].ID, item2.ID)
	}
	if after.Items.Data[1].Price.ID != extra.ID {
		t.Fatalf("remaining item price = %s, want extra %s", after.Items.Data[1].Price.ID, extra.ID)
	}
}

// After middle delete, a new item reuses the lowest unused index; remaining IDs stay put.
func TestSubscriptionItemDeleteThenAddReusesLowestIndex(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)
	extra := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {base.ProductID},
		"currency":            {"usd"},
		"unit_amount":         {"500"},
		"recurring[interval]": {"month"},
	})
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {base.ID},
		"items[0][quantity]": {"1"},
	})
	item0 := created.Items.Data[0].ID
	item1 := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {created.ID},
		"price":        {seat.ID},
		"quantity":     {"1"},
	})
	item2 := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {created.ID},
		"price":        {extra.ID},
		"quantity":     {"1"},
	})
	_ = deleteForm[map[string]any](t, handler, "/v1/subscription_items/"+item1.ID, url.Values{
		"proration_behavior": {"none"},
	})
	added := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"1"},
		"proration_behavior": {"none"},
	})
	// Lowest unused index was 1 (deleted); new item reuses that slot.
	wantReused := billing.FormatSubscriptionItemID(created.ID, 1)
	if added.ID != wantReused {
		t.Fatalf("added id = %s, want reused lowest unused %s", added.ID, wantReused)
	}
	after := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if len(after.Items.Data) != 3 {
		t.Fatalf("items after add = %d, want 3", len(after.Items.Data))
	}
	if after.Items.Data[0].ID != item0 {
		t.Fatalf("item0 id drifted to %s", after.Items.Data[0].ID)
	}
	if after.Items.Data[1].ID != item2.ID {
		t.Fatalf("former last id drifted to %s, want %s", after.Items.Data[1].ID, item2.ID)
	}
	if after.Items.Data[2].ID != wantReused {
		t.Fatalf("new item id = %s, want %s", after.Items.Data[2].ID, wantReused)
	}
}

// 10. DELETE query-string params are parsed (bogus enum → 400).
func TestSubscriptionItemDeleteQueryStringParamsParsed(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {base.ID},
	})
	item := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {created.ID},
		"price":        {seat.ID},
		"quantity":     {"1"},
	})
	status, body := deleteStatus(t, handler, "/v1/subscription_items/"+item.ID+"?proration_behavior=bogus")
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", status, body)
	}
	errBody := decodeErrorBody(t, body)
	if errBody.Error.Code != "parameter_invalid" || errBody.Error.Param != "proration_behavior" {
		t.Fatalf("error=%#v, want invalid proration_behavior from query", errBody.Error)
	}
}

// 11. Item tax_rates evidence on create; totals use subscription default_tax_rates only.
func TestSubscriptionItemCreateTaxRatesEvidenceOnly(t *testing.T) {
	handler := newTestHandler(t)
	customer, base, seat, subTaxID := setupSeatPlans(t, handler)
	// Different item-level rate (20%) — must NOT affect totals if only evidence.
	itemTxr := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/tax_rates", url.Values{
		"display_name": {"ItemVAT"},
		"percentage":   {"20"},
		"inclusive":    {"false"},
	})

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {base.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {subTaxID},
	})

	item := postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"2"},
		"proration_behavior": {"always_invoice"},
		"proration_date":     {strconv.FormatInt(created.CurrentPeriodStart, 10)},
		"tax_rates[0]":       {itemTxr.ID},
	})
	if len(item.TaxRates) != 1 || item.TaxRates[0].ID != itemTxr.ID {
		t.Fatalf("item tax_rates = %#v, want evidence %s", item.TaxRates, itemTxr.ID)
	}

	sub := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	invoice := getJSON[prorationInvoiceResponse](t, handler, "/v1/invoices/"+sub.LatestInvoice)
	// Subscription VAT 10% on delta 2000 → tax 200, not item 20% (400).
	if invoice.Subtotal != 2000 || invoice.Tax != 200 || invoice.Total != 2200 {
		t.Fatalf("invoice = %#v, want subtotal=2000 tax=200 total=2200 (subscription rate, not item 20%%)", invoice)
	}
}

func deleteForm[T any](t *testing.T, handler http.Handler, path string, values url.Values) T {
	t.Helper()
	var body *strings.Reader
	if values != nil {
		body = strings.NewReader(values.Encode())
	} else {
		body = strings.NewReader("")
	}
	req := httptest.NewRequest(http.MethodDelete, path, body)
	req.Host = "billtap.test"
	if values != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return decodeResponse[T](t, rec)
}

func deleteFormStatus(t *testing.T, handler http.Handler, path string, values url.Values) (int, string) {
	t.Helper()
	var body *strings.Reader
	if values != nil {
		body = strings.NewReader(values.Encode())
	} else {
		body = strings.NewReader("")
	}
	req := httptest.NewRequest(http.MethodDelete, path, body)
	req.Host = "billtap.test"
	if values != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
