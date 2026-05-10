package storage

import (
	"context"
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
	if len(versions) != 9 || versions[0] != 1 || versions[1] != 2 || versions[2] != 3 || versions[3] != 4 || versions[4] != 5 || versions[5] != 6 || versions[6] != 7 || versions[7] != 8 || versions[8] != 9 {
		t.Fatalf("versions = %#v, want [1 2 3 4 5 6 7 8 9]", versions)
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
