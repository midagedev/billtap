# Webhook Contract

Billtap webhooks emulate a practical Stripe-style subset. The goal is deterministic local and CI testing, not complete Stripe parity.

## Event Envelope

Every emitted event is stored and delivered with this envelope:

```json
{
  "id": "evt_...",
  "object": "event",
  "type": "checkout.session.completed",
  "created": 1710000000,
  "livemode": false,
  "api_version": "2025-12-15.clover",
  "pending_webhooks": 1,
  "request": {
    "id": "req_...",
    "idempotency_key": "..."
  },
  "data": {
    "object": {}
  },
  "billtap": {
    "scenario_run_id": "sr_...",
    "source": "api|checkout|portal|scenario|replay",
    "sequence": 1
  }
}
```

`api_version` defaults to `2025-12-15.clover` and can be overridden with
`BILLTAP_WEBHOOK_API_VERSION` for applications pinned to a different Stripe SDK
API version.

Required fields for delivery are `id`, `object`, `type`, `created`, `livemode`, `data.object`, and `billtap.sequence`.

## MVP Event Types

| Event type                      | Emitted when                                                     | `data.object`             |
| ------------------------------- | ---------------------------------------------------------------- | ------------------------- |
| `customer.created`              | a customer is created                                            | customer snapshot         |
| `customer.updated`              | a customer is updated                                            | customer snapshot         |
| `product.created`               | a product is created                                             | product snapshot          |
| `price.created`                 | a price is created                                               | price snapshot            |
| `checkout.session.created`      | a checkout session is created                                    | checkout session snapshot |
| `checkout.session.completed`    | checkout succeeds or enters async pending state                  | checkout session snapshot |
| `checkout.session.expired`      | checkout is canceled or expires                                  | checkout session snapshot |
| `customer.subscription.created` | a subscription is created                                        | subscription snapshot     |
| `customer.subscription.updated` | subscription plan, period, status, or cancellation state changes | subscription snapshot     |
| `customer.subscription.deleted` | subscription is canceled immediately                             | subscription snapshot     |
| `invoice.created`               | an invoice is created                                            | invoice snapshot          |
| `invoice.finalized`             | an invoice is finalized                                          | invoice snapshot          |
| `invoice.payment_succeeded`     | invoice payment succeeds                                         | invoice snapshot          |
| `invoice.payment_failed`        | invoice payment fails                                            | invoice snapshot          |
| `invoice.voided`                | invoice is voided after checkout cancellation                    | invoice snapshot          |
| `payment_intent.created`        | a payment intent is created                                      | payment intent snapshot   |
| `payment_intent.succeeded`      | a payment intent succeeds                                        | payment intent snapshot   |
| `payment_intent.processing`     | an async payment remains pending                                 | payment intent snapshot   |
| `payment_intent.canceled`       | payment intent is canceled                                       | payment intent snapshot   |
| `payment_intent.payment_failed` | a payment intent fails                                           | payment intent snapshot   |

SaaS profile events may add `saas.*` profile evidence events, but generic Stripe-like billing events remain the compatibility surface.

## Default Checkout Sequences

Successful subscription checkout emits, in order:

1. `checkout.session.completed`
2. `customer.subscription.created`
3. `invoice.created`
4. `invoice.finalized`
5. `payment_intent.created`
6. `payment_intent.succeeded`
7. `invoice.payment_succeeded`
8. `customer.subscription.updated`

Failed subscription checkout emits, in order:

1. `checkout.session.completed`
2. `customer.subscription.created`
3. `invoice.created`
4. `invoice.finalized`
5. `payment_intent.created`
6. `payment_intent.payment_failed`
7. `invoice.payment_failed`
8. `customer.subscription.updated`

Canceled checkout emits:

1. `checkout.session.expired`
2. `customer.subscription.created`
3. `invoice.created`
4. `invoice.finalized`
5. `payment_intent.created`
6. `payment_intent.canceled`
7. `invoice.voided`
8. `customer.subscription.updated`

Async-pending checkout emits `payment_intent.processing` and no invoice terminal payment event. Requires-action outcomes keep the subscription incomplete and emit the failure sequence with a `requires_action` payment intent status.

## Signatures

Deliveries include a Stripe-compatible-style header:

```text
Billtap-Signature: t=<unix_seconds>,v1=<hex_hmac_sha256>
```

The signed payload is:

```text
<timestamp>.<raw-json-body>
```

The signing secret is the webhook endpoint secret. Verification tolerance defaults to 300 seconds. API and dashboard responses must mask endpoint secrets and signature HMAC values while preserving enough timestamp evidence for debugging.

## Delivery Attempts

Each endpoint/event delivery creates one `delivery_attempts` row.

Statuses:

- `scheduled`: attempt is queued for a future time
- `delivering`: worker has started the HTTP request
- `succeeded`: endpoint returned HTTP 2xx
- `failed`: endpoint returned non-2xx or request failed
- `abandoned`: retry policy is exhausted
- `skipped`: endpoint is inactive or event type is not enabled

Stored fields:

- event id
- endpoint id
- attempt number
- status
- scheduled time
- delivered time when available
- request URL
- masked request URL query parameters
- masked request headers
- request body
- response status
- response body truncated to the configured limit
- error message when available
- next retry time when available

## Retry Policy

Default local retry policy:

- max attempts: 5
- backoff: 10 seconds, 30 seconds, 2 minutes, 10 minutes
- terminal status after exhaustion: `abandoned`

Scenario files may override retry timing for deterministic CI tests. Retry scheduling uses the Billtap local clock when a scenario run is active.

## Reliability Controls

### Duplicate

Duplicate delivery reuses the same event id and payload and creates another delivery attempt. It must not create a new billing transition.

### Delay

Delayed delivery schedules the attempt in the future without changing event creation time or billing object state.

### Out Of Order

Out-of-order delivery changes delivery order only. Event `created` timestamps and `billtap.sequence` values remain the canonical transition order.

### Replay

Replay sends an existing stored event again. The replay attempt references the original event id and records `billtap.source = "replay"` in delivery metadata, not in the original event payload.

Replay actions must create `webhook.replay` audit log entries. Duplicate,
delayed, and out-of-order delivery overrides must create
`webhook.delivery_override` audit log entries.

## Production Boundary

Relay mode is enabled with `BILLTAP_RELAY_MODE=true`. In relay mode,
`raw_payload_storage` is forced to `metadata_only`.

Raw payload policy:

- local default: raw event payloads and delivery request bodies may be stored
- relay mode: signed payloads may be delivered, but raw event payloads and delivery request bodies are not persisted
- retention: old raw event payloads, delivery request bodies, and response bodies are redacted while preserving event ids, types, sequence, statuses, timestamps, metadata, and audit records

Billtap rejects real card data fields such as card number, CVC, and expiration fields.

## Scenario Knobs

Scenario steps may configure webhook behavior with:

```yaml
webhooks:
  retry:
    maxAttempts: 3
    backoff: [1s, 5s]
  delivery:
    delay: 30s
    duplicate: 2
    outOfOrder: true
```

Action-level knobs may override scenario defaults for `webhook.replay`, `webhook.deliver_duplicate`, and `webhook.deliver_out_of_order`.

`webhook.replay` accepts deterministic failure evidence knobs:

- `responseStatus` or `response_status`: records a simulated endpoint HTTP
  status and response body without calling the endpoint.
- `responseBody` or `response_body`: response evidence for the simulated
  endpoint status.
- `timeout`: records a failed delivery attempt with timeout evidence and retry
  scheduling.
- `error` or `simulatedError`: records a failed transport error and retry
  scheduling.
- `signatureMismatch` or `signature_mismatch`: signs the delivery with an
  intentionally invalid HMAC while preserving masked signature evidence.

## Compatibility Boundaries

- Billtap does not claim full Stripe event coverage.
- Event snapshots must be internally consistent with Billtap state.
- Fixture-backed tests are required before claiming compatibility for any endpoint or event type.
- Webhook behavior that affects order, signature, retry, idempotency, or billing state must have tests.
