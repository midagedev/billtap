package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	SignatureHeaderName       = "Billtap-Signature"
	StripeSignatureHeaderName = "Stripe-Signature"
	DefaultTolerance          = 300 * time.Second
)

var (
	ErrInvalidSignatureHeader = errors.New("invalid webhook signature header")
	ErrSignatureMismatch      = errors.New("webhook signature mismatch")
	ErrTimestampOutsideWindow = errors.New("webhook signature timestamp outside tolerance")
)

func Sign(secret string, timestamp time.Time, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(signedPayload(timestamp, payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func SignatureHeader(secret string, timestamp time.Time, payload []byte) string {
	return fmt.Sprintf("t=%d,v1=%s", timestamp.Unix(), Sign(secret, timestamp, payload))
}

func NormalizeSignatureHeaderName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return SignatureHeaderName
	}
	return name
}

func SignatureHeaderValue(headers map[string]string) string {
	if headers == nil {
		return ""
	}
	if value := headers[SignatureHeaderName]; value != "" {
		return value
	}
	if value := headers[StripeSignatureHeaderName]; value != "" {
		return value
	}
	for key, value := range headers {
		if strings.EqualFold(key, SignatureHeaderName) || strings.EqualFold(key, StripeSignatureHeaderName) {
			return value
		}
	}
	return ""
}

func VerifySignature(header string, secret string, now time.Time, payload []byte, tolerance time.Duration) error {
	if tolerance == 0 {
		tolerance = DefaultTolerance
	}
	timestamp, signatures, err := parseSignatureHeader(header)
	if err != nil {
		return err
	}
	signedAt := time.Unix(timestamp, 0).UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	age := now.UTC().Sub(signedAt)
	if age < 0 {
		age = -age
	}
	if age > tolerance {
		return ErrTimestampOutsideWindow
	}

	expected := Sign(secret, signedAt, payload)
	for _, candidate := range signatures {
		if hmac.Equal([]byte(candidate), []byte(expected)) {
			return nil
		}
	}
	return ErrSignatureMismatch
}

func signedPayload(timestamp time.Time, payload []byte) []byte {
	prefix := strconv.FormatInt(timestamp.Unix(), 10) + "."
	out := make([]byte, 0, len(prefix)+len(payload))
	out = append(out, prefix...)
	out = append(out, payload...)
	return out
}

func parseSignatureHeader(header string) (int64, []string, error) {
	parts := strings.Split(header, ",")
	var timestamp int64
	var signatures []string
	for _, part := range parts {
		key, value, found := strings.Cut(strings.TrimSpace(part), "=")
		if !found {
			continue
		}
		switch key {
		case "t":
			parsed, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return 0, nil, ErrInvalidSignatureHeader
			}
			timestamp = parsed
		case "v1":
			if value != "" {
				signatures = append(signatures, value)
			}
		}
	}
	if timestamp == 0 || len(signatures) == 0 {
		return 0, nil, ErrInvalidSignatureHeader
	}
	return timestamp, signatures, nil
}
