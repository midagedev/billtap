# Production Boundaries

## Principle

Billtap is not a production payment processor. Production-facing features must be limited to debugging, replaying, relaying, or validating testmode and controlled billing events.

## Allowed Production-Adjacent Uses

- Inspect staging billing events.
- Replay testmode webhook events.
- Validate webhook handler behavior.
- Run CI billing scenarios.
- Generate support-oriented debug bundles from non-sensitive metadata.

## Disallowed Uses

- Processing real card payments.
- Replacing Stripe Checkout for real customers.
- Replacing Stripe Billing as source of truth.
- Persisting real card data.
- Becoming mandatory for live payment success.

## Relay Mode Rules

Relay mode is enabled with:

```bash
BILLTAP_RELAY_MODE=true
```

When relay mode is enabled, Billtap forces `BILLTAP_RAW_PAYLOAD_STORAGE=metadata_only` even if the config file or environment requests raw payload storage.

Relay mode requirements:

- secrets must be masked
- raw payload storage is off by default
- endpoint failure behavior is explicit
- retry policy is bounded
- all replay actions are auditable

Supported boundary config:

- `BILLTAP_RELAY_MODE=true|false`
- `BILLTAP_RAW_PAYLOAD_STORAGE=store|metadata_only`
- `BILLTAP_RETENTION_DAYS=30`
- `BILLTAP_WEBHOOK_SIGNATURE_HEADER=Billtap-Signature|Stripe-Signature`

## Data Rules

Local:

- raw payloads allowed
- developer controlled

CI:

- artifacts redacted by default
- raw artifacts opt-in

Staging:

- redacted metadata default
- short TTL

Production:

- no raw payment payload storage by default
- no real card data
- no hidden dependency in the live payment path

## Implemented Controls

- API responses mask webhook endpoint secrets as `****`.
- Delivery evidence masks `Authorization`, `Cookie`, API key, token, secret, and webhook signature HMAC values.
- Delivery evidence masks sensitive query parameters in request URLs.
- JSON evidence bodies are redacted for secret-like fields and card-number-like fields.
- Scenario app assertion response bodies are redacted and truncated before reports are written.
- Requests containing real card fields such as `payment_method_data.card.number`, `card[cvc]`, or `card[exp_year]` are rejected.
- Relay mode still sends signed webhook payloads to configured endpoints but stores only metadata in `webhook_events.raw_payload` and `delivery_attempts.request_body`.
- `BILLTAP_WEBHOOK_SIGNATURE_HEADER` changes the delivery header name only; the
  HMAC payload format remains `t=<unix_seconds>,v1=<hex_hmac_sha256>`.
- Retention can redact old webhook raw payloads, delivery request bodies, and response bodies while preserving IDs, timestamps, statuses, event types, and metadata.
- Audit log records webhook replay and delivery override actions.

## Security Requirements

- Webhook signing secrets masked
- API keys masked
- Authorization headers masked
- cookies masked
- dashboard access controls before production-adjacent deployment
- audit log for replay and delivery override actions

## Non-Goals

- Billtap does not process real card payments.
- Billtap does not store PAN, CVC, or live payment credentials.
- Billtap does not replace Stripe Billing or Stripe Checkout in a live payment path.
