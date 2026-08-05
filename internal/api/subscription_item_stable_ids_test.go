package api

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/diagnostics"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
)

// newTestHandlerWithStore returns handler + raw store for legacy row injection.
func newTestHandlerWithStore(t *testing.T) (http.Handler, *storage.SQLiteStore) {
	t.Helper()
	store, err := storage.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})
	opts := Options{
		Billing:     billing.NewService(store),
		Webhooks:    webhooks.NewService(store),
		Diagnostics: diagnostics.NewService(store),
	}
	return New(opts), store
}

// Legacy rows with empty item IDs must still expose index-derived IDs and accept DELETE by them.
func TestLegacySubscriptionItemIDsIndexDerived(t *testing.T) {
	handler, store := newTestHandlerWithStore(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {base.ID},
		"items[0][quantity]": {"1"},
	})
	_ = postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"1"},
		"proration_behavior": {"none"},
	})

	// Strip stored IDs to simulate a pre-migration subscription row.
	ctx := context.Background()
	sub, err := store.GetSubscription(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSubscription: %v", err)
	}
	if len(sub.Items) != 2 {
		t.Fatalf("items = %d, want 2", len(sub.Items))
	}
	legacyItems := []billing.LineItem{
		{PriceID: sub.Items[0].PriceID, Quantity: sub.Items[0].Quantity},
		{PriceID: sub.Items[1].PriceID, Quantity: sub.Items[1].Quantity},
	}
	sub.Items = legacyItems
	if _, err := store.UpdateSubscription(ctx, sub, nil); err != nil {
		t.Fatalf("UpdateSubscription strip ids: %v", err)
	}

	// Response IDs must match index derivation (byte-stable vs pre-fix exposure).
	listed := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	want0 := billing.FormatSubscriptionItemID(created.ID, 0)
	want1 := billing.FormatSubscriptionItemID(created.ID, 1)
	if len(listed.Items.Data) != 2 || listed.Items.Data[0].ID != want0 || listed.Items.Data[1].ID != want1 {
		t.Fatalf("legacy exposed ids = %#v, want %s,%s", listed.Items.Data, want0, want1)
	}

	// DELETE by derived ID must resolve the second item.
	_ = deleteForm[map[string]any](t, handler, "/v1/subscription_items/"+want1, url.Values{
		"proration_behavior": {"none"},
	})
	after := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if len(after.Items.Data) != 1 || after.Items.Data[0].ID != want0 {
		t.Fatalf("after legacy delete = %#v, want single %s", after.Items.Data, want0)
	}
}

// Backfill on update must not change IDs that were already exposed for legacy rows.
func TestLegacySubscriptionItemBackfillKeepsExposedIDs(t *testing.T) {
	handler, store := newTestHandlerWithStore(t)
	customer, base, seat, _ := setupSeatPlans(t, handler)
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {base.ID},
		"items[0][quantity]": {"1"},
	})
	_ = postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {seat.ID},
		"quantity":           {"2"},
		"proration_behavior": {"none"},
	})

	ctx := context.Background()
	sub, err := store.GetSubscription(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSubscription: %v", err)
	}
	// Legacy: clear stored IDs.
	sub.Items = []billing.LineItem{
		{PriceID: sub.Items[0].PriceID, Quantity: sub.Items[0].Quantity},
		{PriceID: sub.Items[1].PriceID, Quantity: sub.Items[1].Quantity},
	}
	if _, err := store.UpdateSubscription(ctx, sub, nil); err != nil {
		t.Fatalf("strip ids: %v", err)
	}

	before := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	id0 := before.Items.Data[0].ID
	id1 := before.Items.Data[1].ID

	// Quantity change triggers PatchSubscription → AssignSubscriptionItemIDs backfill.
	updated := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":       {id0},
		"items[0][quantity]": {"1"},
		"items[1][id]":       {id1},
		"items[1][quantity]": {"3"},
		"proration_behavior": {"none"},
	})
	if len(updated.Items.Data) != 2 {
		t.Fatalf("items after update = %d", len(updated.Items.Data))
	}
	if updated.Items.Data[0].ID != id0 || updated.Items.Data[1].ID != id1 {
		t.Fatalf("backfill changed ids: before=%s,%s after=%s,%s",
			id0, id1, updated.Items.Data[0].ID, updated.Items.Data[1].ID)
	}

	// Confirm IDs were persisted (not only derived on read).
	stored, err := store.GetSubscription(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSubscription after backfill: %v", err)
	}
	if stored.Items[0].ID != id0 || stored.Items[1].ID != id1 {
		t.Fatalf("stored ids = %#v, want %s,%s", stored.Items, id0, id1)
	}
}

// P1 path: items[0][id]=stored ID price swap still issues subscription_update invoice.
func TestSubscriptionUpdateByStoredItemIDProration(t *testing.T) {
	handler := newTestHandler(t)
	customer, lite, pro, taxRateID := setupProrationPlans(t, handler)

	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {lite.ID},
		"items[0][quantity]":   {"1"},
		"default_tax_rates[0]": {taxRateID},
	})
	itemID := created.Items.Data[0].ID
	// Add a second item so the first is not the only index, then keep using stored ID for swap.
	_ = postForm[subItemResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription":       {created.ID},
		"price":              {lite.ID},
		"quantity":           {"1"},
		"proration_behavior": {"none"},
	})
	// Delete second so first retains index-0 ID (stable stored).
	listed := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID)
	if len(listed.Items.Data) != 2 {
		t.Fatalf("items = %d", len(listed.Items.Data))
	}
	_ = deleteForm[map[string]any](t, handler, "/v1/subscription_items/"+listed.Items.Data[1].ID, url.Values{
		"proration_behavior": {"none"},
	})

	upgraded := postForm[prorationSubResponse](t, handler, "/v1/subscriptions/"+created.ID, url.Values{
		"items[0][id]":         {itemID},
		"items[0][price]":      {pro.ID},
		"proration_behavior":   {"always_invoice"},
		"proration_date":       {strconv.FormatInt(created.CurrentPeriodStart, 10)},
		"payment_behavior":     {"error_if_incomplete"},
		"default_tax_rates[0]": {taxRateID},
	})
	if upgraded.Items.Data[0].ID != itemID {
		t.Fatalf("item id after price swap = %s, want stable %s", upgraded.Items.Data[0].ID, itemID)
	}
	if upgraded.Items.Data[0].Price.ID != pro.ID {
		t.Fatalf("price after swap = %s, want %s", upgraded.Items.Data[0].Price.ID, pro.ID)
	}
	invoice := getJSON[prorationInvoiceResponse](t, handler, "/v1/invoices/"+upgraded.LatestInvoice)
	if invoice.BillingReason != "subscription_update" {
		t.Fatalf("billing_reason = %q, want subscription_update", invoice.BillingReason)
	}
	// Full-period upgrade 4900→9900: delta 5000 + VAT10% 500 = 5500.
	if invoice.Subtotal != 5000 || invoice.Tax != 500 || invoice.Total != 5500 {
		t.Fatalf("proration invoice = %#v, want 5000/500/5500", invoice)
	}
}

// Renewal invoices must report billing_reason=subscription_cycle; create stays subscription_create.
func TestRenewalInvoiceBillingReasonSubscriptionCycle(t *testing.T) {
	handler := newTestHandler(t)
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Cycle Plan"}})
	price := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"9900"},
		"recurring[interval]": {"month"},
	})
	frozen := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"id":          {"clock_billing_reason"},
		"frozen_time": {strconv.FormatInt(frozen.Unix(), 10)},
	})
	customer := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{
		"email":      {"cycle@example.test"},
		"test_clock": {"clock_billing_reason"},
	})
	created := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":        {customer.ID},
		"items[0][price]": {price.ID},
		"test_clock":      {"clock_billing_reason"},
	})
	createInvoice := getJSON[prorationInvoiceResponse](t, handler, "/v1/invoices/"+created.LatestInvoice)
	if createInvoice.BillingReason != "subscription_create" {
		t.Fatalf("create billing_reason = %q, want subscription_create", createInvoice.BillingReason)
	}

	// Advance past period end to trigger renewal.
	_ = postForm[map[string]any](t, handler, "/v1/test_helpers/test_clocks/clock_billing_reason/advance", url.Values{
		"frozen_time": {strconv.FormatInt(created.CurrentPeriodEnd+1, 10)},
	})
	invoices := listSubscriptionInvoices(t, handler, created.ID)
	if len(invoices) < 2 {
		t.Fatalf("invoices = %#v, want create + renewal", invoices)
	}
	// Newest first (Stripe list default is reverse chrono in billtap).
	var createCount, cycleCount int
	for _, inv := range invoices {
		switch inv.BillingReason {
		case "subscription_create":
			createCount++
		case "subscription_cycle":
			cycleCount++
		}
	}
	if createCount != 1 {
		t.Fatalf("subscription_create count = %d, want 1 among %#v", createCount, invoices)
	}
	if cycleCount < 1 {
		t.Fatalf("subscription_cycle count = %d, want ≥1 among %#v", cycleCount, invoices)
	}
}
