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
	if len(versions) != 11 || versions[0] != 1 || versions[1] != 2 || versions[2] != 3 || versions[3] != 4 || versions[4] != 5 || versions[5] != 6 || versions[6] != 7 || versions[7] != 8 || versions[8] != 9 || versions[9] != 10 || versions[10] != 11 {
		t.Fatalf("versions = %#v, want [1 2 3 4 5 6 7 8 9 10 11]", versions)
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

func TestDirectIntentSchemaAllowsOptionalCustomerAndPreservesForeignKeys(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLite(ctx, filepath.Join(t.TempDir(), "billtap.db"))
	if err != nil {
		t.Fatalf("OpenSQLite returned error: %v", err)
	}
	defer store.Close()

	if _, err := store.CreatePaymentIntent(ctx, billing.PaymentIntent{
		ID:            "pi_direct_no_customer",
		Amount:        1000,
		Currency:      "usd",
		Status:        "requires_payment_method",
		CaptureMethod: "automatic",
		CreatedAt:     time.Now().UTC(),
	}); err != nil {
		t.Fatalf("CreatePaymentIntent without customer returned error: %v", err)
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
