// Package webhookstest holds test-only helpers for webhook-backed stores. It
// lives outside package webhooks so the production binary never links testing.
package webhookstest

import (
	"context"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/webhooks"
)

// RegisterStoreCleanup registers t.Cleanup that waits for in-flight async
// deliveries and then closes the store. Call after t.TempDir() so LIFO order
// runs this before TempDir removal (avoids "directory not empty" races when
// async deliveries write SQLite WAL/SHM after Close).
//
// Wait failures are logged, not fatal: a hang should not flip a green test, but
// leftover files still surface as TempDir cleanup failures.
func RegisterStoreCleanup(t testing.TB, service *webhooks.Service, closer interface{ Close() error }) {
	t.Helper()
	t.Cleanup(func() {
		if service != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := service.WaitForAsyncDeliveries(ctx); err != nil {
				t.Logf("WaitForAsyncDeliveries: %v", err)
			}
			cancel()
		}
		if closer != nil {
			if err := closer.Close(); err != nil {
				t.Fatalf("close store: %v", err)
			}
		}
	})
}
