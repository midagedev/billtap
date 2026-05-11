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
- `GET /v1/prices/search`

Price search returns a Stripe-like `search_result` envelope. The supported
query subset covers `active:'true|false'`, `type:'one_time|recurring'`,
`lookup_key:'...'`, and `metadata['key']:'value'` clauses joined by `AND`.
This supports one-time price lookup paths while keeping unsupported query
clauses explicit validation errors.

### Coupons

- `POST /v1/coupons`
- `GET /v1/coupons/{id}`
- `GET /v1/coupons`
- `POST /v1/coupons/{id}`
- `DELETE /v1/coupons/{id}`

Coupons are local discount evidence. They support basic percent-off or
amount-off fields, metadata, and deletion markers, but do not apply discounts
to invoices or subscriptions.

### Promotion Codes

- `POST /v1/promotion_codes`
- `GET /v1/promotion_codes/{id}`
- `GET /v1/promotion_codes`
- `POST /v1/promotion_codes/{id}`

Promotion codes are local coupon-linked evidence only. Redemption, customer
discount state, expiration, and promotion-code analytics are not modeled.

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

### Subscription Schedules

- `POST /v1/subscription_schedules`
- `GET /v1/subscription_schedules/{id}`
- `GET /v1/subscription_schedules`
- `POST /v1/subscription_schedules/{id}`
- `POST /v1/subscription_schedules/{id}/cancel`
- `POST /v1/subscription_schedules/{id}/release`

Subscription schedules are local one-phase schedule evidence for existing
subscriptions. Test-clock advance applies a due phase by replacing subscription
items and emitting `customer.subscription.updated`. Multi-phase proration,
invoicing, and full schedule lifecycle parity are not modeled.

### Invoices

- `GET /v1/invoices/{id}`
- `GET /v1/invoices`
- `POST /v1/invoices/{id}/pay`
- `POST /v1/invoices/create_preview`
- `GET /v1/invoices/upcoming`
- `POST /v1/invoices/upcoming`

Direct invoice `pay` is a local retry mutation for open invoices created by
Billtap checkout and scenarios. It accepts deterministic sandbox
`payment_method` or legacy `source` aliases plus bounded protocol flags such as
`paid_out_of_band`, `forgive`, `off_session`, and `mandate`. Direct invoice
create, finalize, send, void, line mutation, collection, and full dunning
automation are not part of the current release-compatible subset.

Preview endpoints accept Stripe SDK-style `subscription`,
`subscription_details[items][0][price]`,
`subscription_details[items][0][quantity]`,
`subscription_details[proration_date]`, and
`subscription_details[proration_behavior]`. Billtap calculates a local
subscription-update proration line from the current period bounds and old/new
price totals. Taxes, discounts, pending invoice items, and collection behavior
are outside the modeled subset.

### Payment Intents

- `GET /v1/payment_intents/{id}`
- `GET /v1/payment_intents`
- `POST /v1/payment_intents`
- `POST /v1/payment_intents/{id}/confirm`
- `POST /v1/payment_intents/{id}/capture`
- `POST /v1/payment_intents/{id}/cancel`

Direct payment intents are local state-machine simulations. They support
deterministic sandbox aliases, manual capture, cancel, timeline evidence, and
local webhook events. `requires_action` returns a local
`next_action.use_stripe_sdk` shape that can be completed or canceled through
Billtap action callbacks, and local bank-transfer intents can move from
`processing` to `succeeded` when customer cash balance is funded. One-time
PaymentIntents can store a deferred per-intent outcome with
`metadata[billtap_payment_intent_outcome]`, `billtap_outcome`, or
`deferred_outcome`; the stored outcome is applied when the intent is confirmed.
If no per-intent outcome is present, customer metadata
`billtap_default_payment_intent_outcome` or `default_payment_intent_outcome`
can provide the default outcome for direct one-time PaymentIntents.

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

Returns deterministic sandbox card projections for known customers, including
the saved default payment method set by local portal payment-method simulation.
The response includes Stripe-like SDK fields such as `billing_details`,
`card.checks`, `card.networks`, `card.three_d_secure_usage`, `metadata`, and
`redaction`. Query validation covers `type`, `allow_redisplay`, `limit`, and
unknown parameters; valid non-card types return an empty local list.

### Customer Cash Balance

- `GET /v1/customers/{id}/cash_balance`
- `POST /v1/customers/{id}/cash_balance`
- `GET /v1/customers/{id}/cash_balance_transactions`
- `GET /v1/customers/{id}/cash_balance_transactions/{id}`
- `POST /v1/test_helpers/customers/{id}/fund_cash_balance`

Cash-balance APIs are local evidence for bank-transfer smoke tests. The
test-helper funding endpoint records a cash-balance transaction and settles
processing bank-transfer PaymentIntents for the customer.

### Refunds

- `POST /v1/refunds`
- `GET /v1/refunds/{id}`
- `GET /v1/refunds`
- `POST /v1/refunds/{id}`
- `POST /v1/refunds/{id}/cancel`

Refunds are local payment-history evidence. Creation accepts `charge`,
`payment_intent`, or `invoice`, plus `amount`, optional `reason`, and metadata.
It emits local `charge.refunded` and `charge.refund.updated` events. A fixture
or API call can keep a refund `pending`; when attached to a test clock with
`settle_at`/`available_on`, clock advance marks it `succeeded` and emits
`charge.refund.updated`. Balance transactions and processor settlement are
outside the modeled subset.

### Credit Notes

- `POST /v1/credit_notes`
- `GET /v1/credit_notes/{id}`
- `GET /v1/credit_notes`
- `POST /v1/credit_notes/{id}/void`

Credit notes are local invoice-history evidence. Creation accepts `invoice`,
`amount`, optional `reason`, and metadata. It emits `credit_note.created`; void
emits `credit_note.voided`. Line-level tax/discount accounting, PDF rendering,
and customer-balance math are outside the modeled subset.

### Disputes

- `GET /v1/disputes`
- `GET /v1/disputes/{id}`
- `POST /v1/disputes/{id}`
- `POST /v1/disputes/{id}/close`
- `GET /v1/charges/{id}/dispute`
- `POST /v1/charges/{id}/dispute`

Disputes are local chargeback-style evidence. Creating one emits
`charge.dispute.created`; updating evidence emits `charge.dispute.updated`;
fixture statuses can also emit `charge.dispute.funds_withdrawn`; closing one
emits `charge.dispute.closed`. Representment deadlines, balance movement, and
processor outcomes are outside the modeled subset.

### Test Clocks

- `POST /v1/test_helpers/test_clocks`
- `GET /v1/test_helpers/test_clocks/{id}`
- `GET /v1/test_helpers/test_clocks`
- `POST /v1/test_helpers/test_clocks/{id}/advance`

Test clocks are persisted local clocks for deterministic lifecycle simulation.
Customers, subscriptions, and pending refunds can be attached through
`test_clock`. Advancing a clock processes due trials, renewals, configured
renewal failures, period-end cancellations, and refund settlement for attached
objects.

### Billing Portal Sessions

- `POST /v1/billing_portal/sessions`

Returns a Stripe-like `billing_portal.session` object and Billtap hosted portal
URL for a known customer. The request accepts `customer`, `return_url`, optional
`configuration`, `locale`, `on_behalf_of`, and `flow_data`. The response
includes explicit `flow`, `locale`, `on_behalf_of`, `livemode`, and
`return_url` fields. Flow enum values and required nested flow fields are
validated with Stripe-style error envelopes. Hosted portal actions can save a
deterministic payment method, cancel a subscription, emit the matching local
webhook evidence, and redirect to `return_url`.

### Connect Platform Evidence

- `POST /v1/accounts`
- `GET /v1/account`
- `GET /v1/accounts/{id}`
- `GET /v1/accounts`
- `POST /v1/accounts/{id}`
- `DELETE /v1/accounts/{id}`
- `POST /v1/accounts/{id}/reject`
- `POST /v1/account_links`
- `POST /v1/account_sessions`
- `POST /v1/accounts/{id}/login_links`
- `GET /v1/accounts/{id}/capabilities`
- `GET /v1/accounts/{id}/capabilities/{capability}`
- `POST /v1/accounts/{id}/capabilities/{capability}`
- `POST /v1/accounts/{id}/people`
- `GET /v1/accounts/{id}/people`
- `GET /v1/accounts/{id}/people/{person}`
- `POST /v1/accounts/{id}/people/{person}`
- `DELETE /v1/accounts/{id}/people/{person}`
- `POST /v1/accounts/{id}/persons`
- `GET /v1/accounts/{id}/persons`
- `GET /v1/accounts/{id}/persons/{person}`
- `POST /v1/accounts/{id}/persons/{person}`
- `DELETE /v1/accounts/{id}/persons/{person}`
- `POST /v1/accounts/{id}/external_accounts`
- `GET /v1/accounts/{id}/external_accounts`
- `GET /v1/accounts/{id}/external_accounts/{external_account}`
- `POST /v1/accounts/{id}/external_accounts/{external_account}`
- `DELETE /v1/accounts/{id}/external_accounts/{external_account}`
- `POST /v1/accounts/{id}/bank_accounts`
- `GET /v1/accounts/{id}/bank_accounts/{bank_account}`
- `POST /v1/accounts/{id}/bank_accounts/{bank_account}`
- `DELETE /v1/accounts/{id}/bank_accounts/{bank_account}`
- `POST /v1/transfers`
- `GET /v1/transfers`
- `GET /v1/transfers/{id}`
- `POST /v1/transfers/{id}`
- `POST /v1/transfers/{id}/reversals`
- `GET /v1/transfers/{id}/reversals`
- `GET /v1/transfers/{id}/reversals/{reversal}`
- `POST /v1/transfers/{id}/reversals/{reversal}`
- `POST /v1/payouts`
- `GET /v1/payouts`
- `GET /v1/payouts/{id}`
- `POST /v1/payouts/{id}`
- `POST /v1/payouts/{id}/cancel`
- `POST /v1/payouts/{id}/reverse`
- `GET /v1/application_fees`
- `GET /v1/application_fees/{id}`
- `POST /v1/application_fees/{id}/refund`
- `POST /v1/application_fees/{id}/refunds`
- `GET /v1/application_fees/{id}/refunds`
- `GET /v1/application_fees/{id}/refunds/{refund}`
- `POST /v1/application_fees/{id}/refunds/{refund}`

Connect APIs are local smoke-test fixtures for platform-style routing. Platform
account retrieval returns deterministic local evidence. Account
create/list/retrieve/update persist connected-account profiles with metadata and
basic capability status, while account deletion returns a local marker rather
than closing a provider account. People/persons persist local representative
evidence. Account links, account sessions, and login links return local hosted
URLs/client secrets for onboarding or embedded-component tests. External
accounts, bank accounts, transfers, reversals, payouts, application fees, and
fee refunds are local evidence objects for platform-style integration tests.
Billtap records `Stripe-Account` request headers in redacted request traces and
can emit local Connect evidence webhooks, but it does not perform real
onboarding, KYC, identity verification, external-account verification, balance
movement, account closure, or settlement.

### Webhook Endpoints

- `POST /v1/webhook_endpoints`
- `GET /v1/webhook_endpoints/{id}`
- `GET /v1/webhook_endpoints`
- `POST /v1/webhook_endpoints/{id}`
- `PATCH /v1/webhook_endpoints/{id}`
- `DELETE /v1/webhook_endpoints/{id}`
- `GET /v1/webhook_endpoints/{id}/attempts`

### Events

- `GET /v1/events/{id}`
- `GET /v1/events`

Event list filters include `type`, `scenarioRunId`, created-time ranges,
`data.object.customer`, and `data.object.metadata[key]`.

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

### `POST /api/events/replay-group`

Replays multiple existing event IDs with ordered, out-of-order, delayed,
signature-mismatch, duplicate, or omitted delivery evidence.

### `POST /api/payment_intents/{id}/complete_action`

Completes a local `requires_action` PaymentIntent.

### `POST /api/payment_intents/{id}/cancel_action`

Cancels a local `requires_action` PaymentIntent.

### `POST /api/payment_intents/{id}/outcome`

Stores a local deferred PaymentIntent outcome before confirmation.

Request:

```json
{
  "outcome": "requires_action"
}
```

The outcome uses the same deterministic aliases accepted by direct
PaymentIntent confirmation, including `payment_succeeded`, `card_declined`,
`requires_action`, `payment_pending`, `bank_transfer`, and `canceled`.

### `POST /api/disputes`

Creates local dispute evidence when a test does not already have a charge-like
ID to use through `/v1/charges/{id}/dispute`.

### `POST /api/events/{id}/replay`

Replay a webhook event. Records `webhook.replay` in the audit log and returns
redacted delivery attempt evidence.

Replay accepts reliability controls such as duplicate delivery, out-of-order
delivery, signature mismatch, forced response status, and
`simulate_app_failure` with `status`, `fail_first_n_attempts`, and optional
`body`. Simulated app failures record failed delivery attempts without calling
the app endpoint for the injected failures, then deliver the real replay
attempt after the configured failures are exhausted.

### `POST /api/webhooks/endpoints/{id}/replay-historical`

Replay historical events to one webhook endpoint. This is a Billtap-specific
catchup API for fixture or startup flows where events already exist before the
application registers its webhook endpoint. It records
`webhook.replay_historical` in the audit log and returns redacted delivery
attempt evidence.

Query or form fields:

- `since`: optional RFC3339 timestamp, Unix timestamp, or `now`
- `until`: optional RFC3339 timestamp, Unix timestamp, or `now`; defaults to
  the endpoint creation time
- `type` / `types` / `event_type` / `event_types`: optional event type filters
  such as `invoice.paid` or `invoice.*`
- `limit`: optional positive replay count
- `force`: optional boolean; when false, events with existing delivery attempts
  for the endpoint are skipped

Historical replay respects the endpoint's `enabled_events` filters, keeps the
original event ID and payload, marks delivery attempts with replay and
historical metadata, and does not replay events created after the endpoint
registration time unless an explicit `until` is provided.

### `POST /api/debug-bundles`

Create a debug bundle.

### `POST /api/fixtures/apply`

Apply a developer-test fixture pack. Request body may be JSON or YAML.

### `POST /api/fixtures/validate`

Validate a fixture pack without mutating local billing state. Request body may
be JSON or YAML. The response contains `valid: true` plus counts for supported
sections when schema and local semantic checks pass.

Supported fixture sections:

- `customers`
- `connected_accounts`
- `catalog.products`
- `catalog.prices`
- `test_clocks`
- `subscriptions`
- `refunds`
- `credit_notes`
- `disputes`
- `assertions`

Customer fixtures can opt out of Billtap's default sandbox card projection with
`payment_methods_fixture: empty` or an explicit empty `payment_methods: []`
field. When either field is present, `GET /v1/payment_methods?customer=...`
returns an empty list until a portal payment-method update or explicit fixture
payment method adds a saved method. A `payment_methods` list may also provide
objects with `id` and optional `default: true` for multi-card local scenarios.

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
If a subscription already exists when the fixture is applied, Billtap backfills
a missing `customer.subscription.created` event for that subscription before
emitting fixture state updates. This keeps re-applied or pre-seeded fixture
graphs replayable without requiring a dedicated checkout flow.

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
