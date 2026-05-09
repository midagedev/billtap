package security

import (
	"strings"
	"testing"
)

func TestRedactHeadersMasksSensitiveValues(t *testing.T) {
	headers := RedactHeaders(map[string]string{
		"Authorization":     "Bearer stripe_key_redaction_sample",
		"Cookie":            "session=secret",
		"Billtap-Signature": "t=1,v1=abc",
	})
	if headers["Authorization"] != MaskedValue || headers["Cookie"] != MaskedValue {
		t.Fatalf("headers = %#v, want sensitive headers masked", headers)
	}
	if headers["Billtap-Signature"] != "t=1,v1=****" {
		t.Fatalf("signature header = %q, want masked signature evidence", headers["Billtap-Signature"])
	}
}

func TestRedactURLMasksSensitiveQuery(t *testing.T) {
	got := RedactURL("https://example.test/webhook?api_key=stripe_key_redaction_sample&workspace=ok")
	if got != "https://example.test/webhook?api_key=%2A%2A%2A%2A&workspace=ok" {
		t.Fatalf("RedactURL = %q", got)
	}
}

func TestRedactTextMasksSensitiveJSON(t *testing.T) {
	got := RedactText(`{"client_secret":"pi_secret","nested":{"card":{"number":"4242"}}}`)
	if got != `{"client_secret":"****","nested":{"card":{"number":"****"}}}` {
		t.Fatalf("RedactText = %s", got)
	}
}

func TestRedactTextMasksSensitiveURLInFreeformError(t *testing.T) {
	got := RedactText(`Post "http://127.0.0.1/webhook?api_key=stripe_key_redaction_sample": connection refused`)
	if strings.Contains(got, "stripe_key_redaction_sample") {
		t.Fatalf("RedactText leaked sensitive URL query: %s", got)
	}
}

func TestContainsCardDataAny(t *testing.T) {
	payload := map[string]any{
		"payment_method_data": map[string]any{
			"card": map[string]any{
				"number": "4242424242424242",
			},
		},
	}
	if !ContainsCardDataAny(payload) {
		t.Fatal("ContainsCardDataAny returned false for nested card data")
	}
	if ContainsCardData(map[string]string{"outcome": "payment_succeeded"}) {
		t.Fatal("ContainsCardData returned true for non-card fields")
	}
}
