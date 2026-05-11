# API Contract

Billtap exposes a Stripe-like API subset under `/v1` and Billtap-specific dashboard APIs under `/api`.
The public compatibility snapshot in `docs/COMPATIBILITY.md` is the release
source of truth for supported and unsupported provider behavior.

Webhook event shape, signature, retry, duplicate, delay, out-of-order, and replay behavior are defined in [webhooks.md](webhooks.md).

## Health

### `GET /healthz`

Process health.

### `GET /readyz`

Storage and worker readiness.

## Stripe-like API

### Customers

- `POST /v1/customers`
- `GET /v1/customers/{id}`
- `GET /v1/customers`
- `POST /v1/customers/{id}`

### Products

- `POST /v1/products`
- `GET /v1/products/{id}`
- `GET /v1/products`
- `GET /v1/products/search`

### Prices

- `POST /v1/prices`
- `GET /v1/prices/{id}`
- `GET /v1/prices`

### Checkout Sessions

- `POST /v1/checkout/sessions`
- `GET /v1/checkout/sessions/{id}`
- `GET /v1/checkout/sessions`

Response includes:

- `id`
- `object`
- `url`
- `status`
- `payment_status`

### Subscriptions

- `GET /v1/subscriptions/{id}`
- `GET /v1/subscriptions`
- `POST /v1/subscriptions/{id}`
- `DELETE /v1/subscriptions/{id}`

### Invoices

- `GET /v1/invoices/{id}`
- `GET /v1/invoices`
- `POST /v1/invoices/{id}/pay`
- `POST /v1/invoices/create_preview`

Direct invoice `pay` is a local retry mutation for open invoices created by
Billtap checkout and scenarios. It accepts deterministic sandbox
`payment_method` or legacy `source` aliases plus bounded protocol flags such as
`paid_out_of_band`, `forgive`, `off_session`, and `mandate`. Direct invoice
create, finalize, send, void, line mutation, collection, and full dunning
automation are not part of the current release-compatible subset.

### Payment Intents

- `GET /v1/payment_intents/{id}`
- `GET /v1/payment_intents`
- `POST /v1/payment_intents`
- `POST /v1/payment_intents/{id}/confirm`
- `POST /v1/payment_intents/{id}/capture`
- `POST /v1/payment_intents/{id}/cancel`

Direct payment intents are local state-machine simulations. They support
deterministic sandbox aliases, manual capture, cancel, timeline evidence, and
local webhook events; they do not process real cards or claim full Stripe
PaymentIntent parameter parity.

### Setup Intents

- `GET /v1/setup_intents/{id}`
- `GET /v1/setup_intents`
- `POST /v1/setup_intents`
- `POST /v1/setup_intents/{id}/confirm`
- `POST /v1/setup_intents/{id}/cancel`

Setup intents are local state-machine simulations for saved-payment-method
smoke tests. Mandates, bank verification, and full SCA behavior are not part of
the current release-compatible subset.

### Payment Methods

- `GET /v1/payment_methods?customer={id}&type=card`

Returns deterministic sandbox card projections for known customers.

### Refunds

- `POST /v1/refunds`
- `GET /v1/refunds/{id}`
- `GET /v1/refunds`

Refunds are local payment-history evidence. Creation accepts `charge`,
`payment_intent`, or `invoice`, plus `amount`, optional `reason`, and metadata.
It emits local `charge.refunded` and `charge.refund.updated` events, but does
not model settlement, balance transactions, disputes, or failed refund
processing.

### Credit Notes

- `POST /v1/credit_notes`
- `GET /v1/credit_notes/{id}`
- `GET /v1/credit_notes`

Credit notes are local invoice-history evidence. Creation accepts `invoice`,
`amount`, optional `reason`, and metadata. It emits `credit_note.created`, but
does not model line-level tax/discount accounting, PDF rendering, or voiding.

### Test Clocks

- `POST /v1/test_helpers/test_clocks`
- `GET /v1/test_helpers/test_clocks/{id}`
- `GET /v1/test_helpers/test_clocks`
- `POST /v1/test_helpers/test_clocks/{id}/advance`

Test clocks are persisted local clocks for deterministic lifecycle simulation.
Customers and subscriptions can be attached through `test_clock`. Advancing a
clock processes due trials, renewals, configured renewal failures, and
period-end cancellations for attached objects.

### Billing Portal Sessions

- `POST /v1/billing_portal/sessions`

Returns a Billtap hosted portal URL for a known customer.

### Webhook Endpoints

- `POST /v1/webhook_endpoints`
- `GET /v1/webhook_endpoints/{id}`
- `GET /v1/webhook_endpoints`
- `POST /v1/webhook_endpoints/{id}`
- `PATCH /v1/webhook_endpoints/{id}`
- `DELETE /v1/webhook_endpoints/{id}`

### Events

- `GET /v1/events/{id}`
- `GET /v1/events`

## Hosted UI

### `GET /checkout/{sessionId}`

Hosted sandbox checkout page.

### `GET /portal/{customerId}`

Hosted sandbox billing portal page.

## Dashboard API

### `GET /api/timeline`

Filters:

- customerId
- subscriptionId
- checkoutSessionId
- invoiceId
- paymentIntentId
- eventId
- scenarioRunId

### `GET /api/delivery-attempts`

Webhook delivery attempts. Response evidence masks endpoint credentials,
sensitive headers, sensitive request URL query parameters, and webhook
signature HMAC values.

### `POST /api/events/{id}/replay`

Replay a webhook event. Records `webhook.replay` in the audit log and returns
redacted delivery attempt evidence.

Replay accepts reliability controls such as duplicate delivery, out-of-order
delivery, signature mismatch, forced response status, and
`simulate_app_failure` with `status`, `fail_first_n_attempts`, and optional
`body`. Simulated app failures record failed delivery attempts without calling
the app endpoint for the injected failures, then deliver the real replay
attempt after the configured failures are exhausted.

### `POST /api/debug-bundles`

Create a debug bundle.

### `POST /api/fixtures/apply`

Apply a developer-test fixture pack. Request body may be JSON or YAML.

Supported fixture sections:

- `customers`
- `catalog.products`
- `catalog.prices`
- `test_clocks`
- `subscriptions`
- `refunds`
- `credit_notes`
- `assertions`

Billtap tags created objects with fixture metadata:

- `billtap_fixture_name`
- `billtap_fixture_run_id`
- `billtap_fixture_namespace`
- `billtap_fixture_ref`

Subscription fixtures are created through the same checkout-completion path as
normal billing flows, so subscriptions, invoices, payment intents, checkout
sessions, and timeline evidence remain consistent.
When the HTTP fixture apply API is used, Billtap also creates seeded webhook
events for the local checkout, subscription, invoice, payment-intent, refund,
and credit-note evidence so tests can list and replay those events through the
same `/v1/events` and `/api/events/{id}/replay` paths.

Fixture-provided IDs are preserved for seeded objects. Fixtures also tag
objects with `billtap_fixture_ref`, and the resolve endpoint below can map a
fixture ref to the generated or stable customer, subscription, invoice, payment
intent, checkout session, product, and price IDs.

### `GET /api/fixtures/resolve`

Resolve fixture-backed objects. Query fields:

- `ref`
- `id`
- `lookup_key` / `lookupKey`
- `runId`
- `fixture` / `fixtureName`
- `namespace`

Returns the matching local IDs for the seeded object graph.

### `GET /api/fixtures/snapshot`

Return a filtered billing snapshot for fixture-driven tests. Query fields:

- `customer` / `customerId`
- `runId`
- `tenantId`
- `fixture` / `fixtureName`
- `namespace`

Response includes customers, products, prices, checkout sessions,
subscriptions, invoices, payment intents, timeline entries, and object counts.

### `POST /api/fixtures/assert`

Assert expected fixture state. Request body may be JSON or YAML.

Supported assertion targets:

- `customer`
- `product`
- `price`
- `checkout_session`
- `subscription`
- `invoice`
- `payment_intent`
- `timeline`

Returns `200` when all assertions pass and `409` with a structured assertion
report when any assertion fails.

### `GET /api/audit-log`

List audit log entries. Query fields:

- `action`
- `targetId`

### `POST /api/retention/apply`

Apply the configured retention policy. Old webhook raw payloads and delivery
request/response bodies are redacted while IDs, statuses, timestamps, metadata,
and audit records are preserved.

### `GET /api/portal`

Load portal state. Accepts `customer_id` or `customerId`.

### `GET /api/portal/customers/{id}`

Load portal state for a customer.

### `POST /api/portal/subscriptions/{id}/plan-change`

Apply a sandbox plan change. Body fields:

- `plan`
- `price`
- `quantity`

### `POST /api/portal/subscriptions/{id}/seat-change`

Apply a sandbox seat quantity change. Body fields:

- `quantity`

### `POST /api/portal/subscriptions/{id}/cancel`

Cancel a sandbox subscription. Body fields:

- `mode`: `period` or `immediate`

### `POST /api/portal/subscriptions/{id}/resume`

Resume a subscription from pending or immediate cancellation state.

### `POST /api/portal/customers/{id}/payment-method`

Simulate payment method update. Body fields:

- `outcome`: `succeeds` or `fails`

### `POST /api/scenarios/run`

Run a scenario.

Request body may be a scenario JSON object or YAML content.

Response is a scenario report with:

- `name`
- `status`
- `failure_type`
- `clock_start`
- `clock_end`
- `steps`
- `errors`
