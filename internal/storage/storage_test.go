package storage

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
)

func TestSQLiteMigrationsRun(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	versions, err := store.MigrationVersions(ctx)
	if err != nil {
		t.Fatalf("MigrationVersions returned error: %v", err)
	}
	// Expect every embedded migration, numbered 1..N with no gap. Spelling the list
	// out by hand meant editing this test on every migration, and a gap — the failure
	// that actually matters, since a skipped number means a file never ran — read the
	// same as "the count moved".
	files, err := fs.ReadDir(migrations, "migrations")
	if err != nil {
		t.Fatalf("read migrations: %v", err)
	}
	want := make([]int, 0, len(files))
	for i := range files {
		want = append(want, i+1)
	}
	if !slices.Equal(versions, want) {
		t.Fatalf("versions = %#v, want %#v", versions, want)
	}
}

func TestMemoryStoreWorksInTests(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, Options{Driver: DriverMemory})
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	if err := store.Ping(ctx); err != nil {
		t.Fatalf("Ping returned error: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if err := store.Ping(ctx); err == nil {
		t.Fatal("Ping after Close succeeded, want error")
	}
}

func TestPendingInvoiceItemAllowsNullInvoiceID(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	service := billing.NewService(store)
	customer, err := service.CreateCustomer(ctx, billing.Customer{ID: "cus_pending", Email: "pending@example.test"})
	if err != nil {
		t.Fatalf("CreateCustomer: %v", err)
	}
	item, invoice, err := service.CreateInvoiceItem(ctx, billing.InvoiceItem{
		CustomerID:  customer.ID,
		Amount:      500,
		Currency:    "usd",
		Description: "pending usage",
	})
	if err != nil {
		t.Fatalf("CreateInvoiceItem pending: %v", err)
	}
	if invoice.ID != "" {
		t.Fatalf("pending create returned invoice %#v, want none", invoice)
	}
	if item.InvoiceID != "" {
		t.Fatalf("pending item invoice_id = %q, want empty", item.InvoiceID)
	}
	pending, err := service.ListInvoiceItems(ctx, billing.InvoiceItemFilter{CustomerID: customer.ID, Pending: true})
	if err != nil {
		t.Fatalf("ListInvoiceItems pending: %v", err)
	}
	if len(pending) != 1 || pending[0].ID != item.ID {
		t.Fatalf("pending list = %#v, want the created item", pending)
	}

	_, err = store.db.ExecContext(ctx, `INSERT INTO invoice_items (`+invoiceItemColumns+`)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"ii_missing_invoice", customer.ID, "in_missing", int64(100), "usd", "", "{}", time.Now().UTC().Format(time.RFC3339Nano), "", "", int64(0), "")
	if err == nil {
		t.Fatal("insert with unknown invoice_id succeeded, want foreign key error")
	}
}

func TestCheckoutSessionMetadataDefaultAndRoundTrip(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	service := billing.NewService(store)
	customer, err := service.CreateCustomer(ctx, billing.Customer{ID: "cus_meta", Email: "meta@example.test"})
	if err != nil {
		t.Fatalf("CreateCustomer: %v", err)
	}
	product, err := service.CreateProduct(ctx, billing.Product{ID: "prod_meta", Name: "Meta"})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}
	price, err := service.CreatePrice(ctx, billing.Price{ID: "price_meta", ProductID: product.ID, Currency: "usd", UnitAmount: 1500})
	if err != nil {
		t.Fatalf("CreatePrice: %v", err)
	}

	// Pre-019-style insert: omit metadata column so DEFAULT '{}' applies; decode must not break.
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := store.db.ExecContext(ctx, `INSERT INTO checkout_sessions (
		id, customer_id, mode, line_items, discounts, success_url, cancel_url, status, payment_status,
		allow_promotion_codes, trial_period_days, automatic_tax, tax_id_collection, tax_percent,
		default_tax_rates, client_reference_id, payment_intent_data, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0, 0, 0, 0, '[]', '', '{}', ?)`,
		"cs_default_meta", customer.ID, "payment", `[{"price":"price_meta","quantity":1}]`, "[]",
		"https://x.test/ok", "", "open", "unpaid", now); err != nil {
		t.Fatalf("insert without metadata column: %v", err)
	}
	var rawMeta string
	if err := store.db.QueryRowContext(ctx, `SELECT metadata FROM checkout_sessions WHERE id = ?`, "cs_default_meta").Scan(&rawMeta); err != nil {
		t.Fatalf("SELECT metadata: %v", err)
	}
	if rawMeta != "{}" {
		t.Fatalf("DEFAULT metadata = %q, want {}", rawMeta)
	}
	legacy, err := store.GetCheckoutSession(ctx, "cs_default_meta")
	if err != nil {
		t.Fatalf("GetCheckoutSession legacy: %v", err)
	}
	if len(legacy.Metadata) != 0 {
		t.Fatalf("legacy Metadata = %#v, want empty", legacy.Metadata)
	}

	// Explicit metadata round-trip + survives completion UPDATE (metadata column not touched).
	created, err := store.CreateCheckoutSession(ctx, billing.CheckoutSession{
		ID:            "cs_with_meta",
		CustomerID:    customer.ID,
		Mode:          "payment",
		LineItems:     []billing.LineItem{{PriceID: price.ID, Quantity: 1}},
		Status:        "open",
		PaymentStatus: "unpaid",
		Metadata:      map[string]string{"paymentType": "EXTRA_EXPORT", "accountId": "acc_1"},
		CreatedAt:     time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("CreateCheckoutSession: %v", err)
	}
	if created.Metadata["paymentType"] != "EXTRA_EXPORT" || created.Metadata["accountId"] != "acc_1" {
		t.Fatalf("created metadata = %#v", created.Metadata)
	}

	completed, err := store.RecordCheckoutCompletion(ctx, billing.CheckoutCompletion{
		SessionID:     created.ID,
		SessionStatus: "complete",
		Outcome:       "payment_succeeded",
		CompletedAt:   time.Now().UTC(),
		PaymentIntent: billing.PaymentIntent{
			ID:         "pi_meta_keep",
			CustomerID: customer.ID,
			Amount:     1500,
			Currency:   "usd",
			Status:     "succeeded",
			Metadata:   map[string]string{"paymentType": "PI_SIDE"},
			CreatedAt:  time.Now().UTC(),
		},
	})
	if err != nil {
		t.Fatalf("RecordCheckoutCompletion: %v", err)
	}
	if completed.Metadata["paymentType"] != "EXTRA_EXPORT" || completed.Metadata["accountId"] != "acc_1" {
		t.Fatalf("completed session metadata = %#v, want preserved", completed.Metadata)
	}
}

func TestDirectIntentSchemaAllowsOptionalCustomerAndPreservesForeignKeys(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	created, err := store.CreatePaymentIntent(ctx, billing.PaymentIntent{
		ID:            "pi_direct_no_customer",
		Amount:        1000,
		Currency:      "usd",
		Status:        "requires_payment_method",
		CaptureMethod: "automatic",
		Metadata:      map[string]string{billing.MetadataPaymentIntentOutcome: "card_declined"},
		CreatedAt:     time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("CreatePaymentIntent without customer returned error: %v", err)
	}
	if created.Metadata[billing.MetadataPaymentIntentOutcome] != "card_declined" {
		t.Fatalf("created payment intent metadata = %#v, want deferred outcome metadata", created.Metadata)
	}

	if _, err := store.CreatePaymentIntent(ctx, billing.PaymentIntent{
		ID:            "pi_direct_bad_customer",
		CustomerID:    "cus_missing",
		Amount:        1000,
		Currency:      "usd",
		Status:        "requires_payment_method",
		CaptureMethod: "automatic",
		CreatedAt:     time.Now().UTC(),
	}); err == nil {
		t.Fatalf("CreatePaymentIntent with missing customer succeeded, want FK error")
	}

	if _, err := store.CreateSetupIntent(ctx, billing.SetupIntent{
		ID:         "seti_direct_bad_customer",
		CustomerID: "cus_missing",
		Status:     "requires_payment_method",
		Usage:      "off_session",
		CreatedAt:  time.Now().UTC(),
	}); err == nil {
		t.Fatalf("CreateSetupIntent with missing customer succeeded, want FK error")
	}
}

func TestUpdateInvoicePaymentRejectsStaleOpenAttempt(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	service := billing.NewService(store)
	customer, err := service.CreateCustomer(ctx, billing.Customer{ID: "cus_guard", Email: "guard@example.test"})
	if err != nil {
		t.Fatalf("CreateCustomer: %v", err)
	}
	product, err := service.CreateProduct(ctx, billing.Product{ID: "prod_guard", Name: "Guard"})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}
	price, err := service.CreatePrice(ctx, billing.Price{ID: "price_guard", ProductID: product.ID, Currency: "usd", UnitAmount: 4900})
	if err != nil {
		t.Fatalf("CreatePrice: %v", err)
	}
	session, err := service.CreateCheckoutSession(ctx, billing.CheckoutSession{
		ID:         "cs_guard",
		CustomerID: customer.ID,
		LineItems:  []billing.LineItem{{PriceID: price.ID, Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("CreateCheckoutSession: %v", err)
	}
	completed, err := service.CompleteCheckout(ctx, session.ID, "payment_failed")
	if err != nil {
		t.Fatalf("CompleteCheckout: %v", err)
	}

	staleInvoice, err := store.GetInvoice(ctx, completed.InvoiceID)
	if err != nil {
		t.Fatalf("GetInvoice: %v", err)
	}
	staleSub, err := store.GetSubscription(ctx, completed.SubscriptionID)
	if err != nil {
		t.Fatalf("GetSubscription: %v", err)
	}
	staleIntent, err := store.GetPaymentIntent(ctx, completed.PaymentIntentID)
	if err != nil {
		t.Fatalf("GetPaymentIntent: %v", err)
	}
	if _, err := service.PayInvoice(ctx, staleInvoice.ID, billing.InvoicePaymentOptions{PaymentMethodID: "pm_card_visa"}); err != nil {
		t.Fatalf("PayInvoice: %v", err)
	}

	staleSub.Status = "past_due"
	staleInvoice.Status = "open"
	staleInvoice.AttemptCount++
	staleInvoice.AmountPaid = 0
	staleInvoice.AmountDue = staleInvoice.Total
	nextAttempt := time.Now().UTC().Add(24 * time.Hour)
	staleInvoice.NextPaymentAttempt = &nextAttempt
	staleIntent.Status = "requires_payment_method"
	staleIntent.FailureCode = "card_declined"
	staleIntent.DeclineCode = "generic_decline"
	if _, _, _, err := store.UpdateInvoicePayment(ctx, staleSub, staleInvoice, staleIntent, nil); !errors.Is(err, billing.ErrInvalidInput) {
		t.Fatalf("UpdateInvoicePayment stale err = %v, want ErrInvalidInput", err)
	}

	current, err := store.GetInvoice(ctx, staleInvoice.ID)
	if err != nil {
		t.Fatalf("GetInvoice current: %v", err)
	}
	if current.Status != "paid" {
		t.Fatalf("current invoice = %#v, want paid invoice preserved", current)
	}
}

// A timeline id is a deterministic seed for a billing event, so the same event can be recorded
// twice — a repeated clock advance replays the activation for the same subscription and period.
// That must be a no-op. Before this guard it raised a UNIQUE constraint error that surfaced as a
// 500 and wedged the fixture: the clock had already moved, so every later advance collided again.
func TestRecordTimelineIsIdempotentForRepeatedEventIDs(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	entry := billing.TimelineEntry{
		ID:         "tl_clock_trial_activate_sub_1_2030-01-15T00:00:00Z",
		Object:     billing.ObjectTimelineEntry,
		Action:     "customer.subscription.updated",
		Message:    "Trial subscription activated",
		ObjectType: billing.ObjectSubscription,
		ObjectID:   "sub_1",
		CreatedAt:  time.Date(2030, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	if err := store.RecordTimeline(ctx, entry); err != nil {
		t.Fatalf("first RecordTimeline returned error: %v", err)
	}
	if err := store.RecordTimeline(ctx, entry); err != nil {
		t.Fatalf("repeated RecordTimeline returned error: %v", err)
	}

	entries, err := store.Timeline(ctx, billing.TimelineFilter{})
	if err != nil {
		t.Fatalf("Timeline returned error: %v", err)
	}
	matching := 0
	for _, got := range entries {
		if got.ID == entry.ID {
			matching++
		}
	}
	if matching != 1 {
		t.Fatalf("timeline holds %d copies of %s, want exactly 1", matching, entry.ID)
	}
}

func TestLocalEvidenceRoundTrip(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	if err := store.SaveLocalEvidence(ctx, "tax_rate", "txr_1", `{"id":"txr_1"}`); err != nil {
		t.Fatalf("SaveLocalEvidence: %v", err)
	}
	// Same key twice must update, not fail on the primary key.
	if err := store.SaveLocalEvidence(ctx, "tax_rate", "txr_1", `{"id":"txr_1","active":false}`); err != nil {
		t.Fatalf("SaveLocalEvidence (update): %v", err)
	}
	if err := store.SaveLocalEvidence(ctx, "coupon", "txr_1", `{"id":"txr_1","object":"coupon"}`); err != nil {
		t.Fatalf("SaveLocalEvidence (other kind, same id): %v", err)
	}

	all, err := store.LoadLocalEvidence(ctx)
	if err != nil {
		t.Fatalf("LoadLocalEvidence: %v", err)
	}
	if got := all["tax_rate"]["txr_1"]; got != `{"id":"txr_1","active":false}` {
		t.Fatalf("tax_rate row = %q, want the updated document", got)
	}
	if _, ok := all["coupon"]["txr_1"]; !ok {
		t.Fatal("an id shared across kinds must not collide")
	}

	if err := store.DeleteLocalEvidence(ctx, "tax_rate", "txr_1"); err != nil {
		t.Fatalf("DeleteLocalEvidence: %v", err)
	}
	// Deleting what is already gone is not an error.
	if err := store.DeleteLocalEvidence(ctx, "tax_rate", "txr_1"); err != nil {
		t.Fatalf("DeleteLocalEvidence (absent): %v", err)
	}
	all, err = store.LoadLocalEvidence(ctx)
	if err != nil {
		t.Fatalf("LoadLocalEvidence after delete: %v", err)
	}
	if _, ok := all["tax_rate"]["txr_1"]; ok {
		t.Fatal("deleted row came back")
	}
}
