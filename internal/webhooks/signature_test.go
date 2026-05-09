package webhooks

import (
	"strings"
	"testing"
	"time"
)

func TestSignatureHeaderVerifies(t *testing.T) {
	payload := []byte(`{"id":"evt_test"}`)
	now := time.Unix(1778233312, 0).UTC()
	header := SignatureHeader("webhook_secret_test", now, payload)

	if !strings.HasPrefix(header, "t=1778233312,v1=") {
		t.Fatalf("header = %q, want timestamp and v1 signature", header)
	}
	if err := VerifySignature(header, "webhook_secret_test", now, payload, DefaultTolerance); err != nil {
		t.Fatalf("VerifySignature returned error: %v", err)
	}
}

func TestSignatureRejectsMismatches(t *testing.T) {
	payload := []byte(`{"id":"evt_test"}`)
	now := time.Unix(1778233312, 0).UTC()
	header := SignatureHeader("webhook_secret_test", now, payload)

	if err := VerifySignature(header, "wrong", now, payload, DefaultTolerance); err == nil {
		t.Fatal("VerifySignature accepted wrong secret")
	}
	if err := VerifySignature(header, "webhook_secret_test", now.Add(10*time.Minute), payload, DefaultTolerance); err == nil {
		t.Fatal("VerifySignature accepted expired timestamp")
	}
}
