package storage

import (
	"context"
	"errors"
	"path/filepath"
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
	want := make([]int, len(versions))
	for i := range want {
		want[i] = i + 1
	}
	if len(versions) != 20 {
		t.Fatalf("versions = %#v, want contiguous 1..%d (migration 020 = invoiceitem pending + credit note out_of_band)", versions, len(want))
	}
	for i := range want {
		if versions[i] != want[i] {
			t.Fatalf("versions = %#v, want %#v", versions, want)
		}
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
