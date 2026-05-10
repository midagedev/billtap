# Compatibility

Billtap is a fixture-backed billing lab for local development and CI. It is not
a full Stripe clone, a Stripe Dashboard replacement, or a payment processor.

The compatibility promise is intentionally narrow:

- Stripe-like request and response shapes where they help subscription apps run
  deterministic local tests.
- Hosted sandbox checkout and portal flows for exercising app integration code.
- Webhook delivery evidence for retries, duplicate delivery, delay,
  out-of-order delivery, and replay.
- Scenario and fixture APIs that make the supported subset repeatable.

Anything outside this document should be treated as unsupported until it has a
fixture, a test, and an explicit compatibility note.

Known Stripe OpenAPI routes that are not implemented by Billtap first run
OpenAPI-derived parameter validation. Malformed requests return Stripe-shaped
`parameter_unknown`, `parameter_missing`, or `parameter_invalid` errors; valid
but unimplemented requests return `unsupported_endpoint` instead of silently
approximating provider behavior. When diagnostics are enabled, the request trace
captures that error code and the original path so agents can distinguish bad
test setup, unsupported coverage gaps, missing local data, and webhook failures.

## Compatibility Levels

| Level            | Meaning                                                                                                                                                                                                     |
| ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Supported        | Implemented in the local runtime and covered by automated tests or release-blocking compatibility scorecard cases. Examples and fixtures are supporting evidence, not the sole basis for a supported claim. |
| Billtap-specific | Public Billtap API, not intended to match Stripe.                                                                                                                                                           |
| Partial          | Useful for smoke tests, but not a full provider behavior model.                                                                                                                                             |
| Unsupported      | Not implemented, not claimed, or intentionally out of scope.                                                                                                                                                |

## Compatibility Scorecard

The offline compatibility scorecard can be generated without external services:

```bash
go run ./cmd/billtap compatibility scorecard --output-dir dist/compatibility
```

It writes:

- `compatibility-scorecard.json`
- `compatibility-scorecard.md`
- `replay-bundles/*.json` for any `mismatch` or `error` case

Scorecard statuses are:

- `imported`: case ran against Billtap and matched the normalized expectation.
- `skipped`: case is in the corpus but intentionally not run by the offline lane.
- `unsupported`: case documents unsupported behavior from this compatibility contract.
- `mismatch`: case ran but normalized actual behavior differed from expectation.
- `error`: the scorecard runner or Billtap returned an unexpected internal error.

Current public-readiness corpus:

- Scorecard version: `l3-public-readiness-v6`
- Release-blocking cases: 44
- Covered categories: request validation, protocol parameter acceptance,
  OpenAPI-backed fallback validation, idempotency mismatch, deterministic
  checkout payment-error aliases, and direct PaymentIntent/SetupIntent state
  transitions, including invalid terminal-state transitions
- Required release result: `mismatch=0`, `error=0`, and `passed=true`

The scorecard is intentionally a release contract for Billtap's documented
local subset. It is not a claim of broad Stripe API parity.

Broader Stripe API compatibility work is tracked separately in
`STRIPE_API_COMPATIBILITY_ROADMAP.md`. New endpoint families should move through
inventory, schema validation, fixture response, stateful local behavior,
scenario coverage, webhook modeling, and SDK smoke levels before they become
public claims.

The generated Stripe API inventory also reports
`summary.schema_validated_operations`, which counts operations from the input
OpenAPI file that expose parameter or request-body schemas and also match
Billtap's bundled OpenAPI-derived validation catalog. This is a diagnostic and
planning metric only. It does not increase `summary.implemented_operations`; an
endpoint still needs an explicit runtime claim, tests, fixtures, and
documentation before it counts as implemented.

## Supported Stripe-Like API Subset

Base path: `/v1`

| Resource                | Endpoints                                                                                                                                                           | Level            | Scope                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Customers               | `POST /v1/customers`, `GET /v1/customers`, `GET /v1/customers/{id}`, `POST /v1/customers/{id}`                                                                      | Supported        | Create, list, retrieve, and update `email`, `name`, and metadata. List supports `email` and `limit` filters.                                                                                                                                                                                                                                                                                                                                                            |
| Products                | `POST /v1/products`, `GET /v1/products`, `GET /v1/products/{id}`, `POST /v1/products/{id}`                                                                          | Supported        | Create, list, retrieve, and update local service products with metadata.                                                                                                                                                                                                                                                                                                                                                                                                |
| Product search          | `GET /v1/products/search`                                                                                                                                           | Partial          | Supports metadata equality filters such as `metadata['tenantId']:'saas'` and `active:'true'`. This is not Stripe Search Query Language parity.                                                                                                                                                                                                                                                                                                                          |
| Prices                  | `POST /v1/prices`, `GET /v1/prices`, `GET /v1/prices/{id}`, `POST /v1/prices/{id}`                                                                                  | Supported        | Create, list, retrieve, and update prices. Supports `product`, `currency`, `unit_amount`, `lookup_key`, recurring interval fields, `active`, and metadata. List supports `product`, `active`, and `type=recurring`.                                                                                                                                                                                                                                                     |
| Checkout sessions       | `POST /v1/checkout/sessions`, `GET /v1/checkout/sessions`, `GET /v1/checkout/sessions/{id}`                                                                         | Supported        | Creates subscription-mode sandbox checkout sessions from request line items and hosted Billtap URLs. The Stripe-style session response leaves `line_items` unexpanded. Accepts Stripe SDK form params `allow_promotion_codes` and `subscription_data[trial_period_days]`; trial checkout creates local `trialing` subscription evidence. Hosted URLs use the request host by default, or `BILLTAP_PUBLIC_BASE_URL` when configured for container-to-host browser flows. |
| Checkout completion     | `POST /v1/checkout/sessions/{id}/complete`, `POST /api/checkout/sessions/{id}/complete`                                                                             | Billtap-specific | Completes a sandbox checkout and creates subscription, invoice, payment intent, timeline, and checkout webhook evidence. Supports success plus deterministic failure aliases such as `card_declined`, `insufficient_funds`, `expired_card`, `incorrect_cvc`, `processing_error`, `authentication_required`, `payment_pending`, `canceled`, and documented Stripe test PaymentMethod IDs such as `pm_card_visa_chargeDeclined`.                                      |
| Billing portal sessions | `POST /v1/billing_portal/sessions`                                                                                                                                  | Partial          | Returns a Billtap portal URL for a known customer. Portal configuration and full Stripe-hosted portal behavior are not modeled, but the hosted portal includes Stripe-like selectors for common local Page Object flows.                                                                                                                                                                                                                                                |
| Subscriptions           | `POST /v1/subscriptions`, `GET /v1/subscriptions`, `GET /v1/subscriptions/{id}`, `POST /v1/subscriptions/{id}`, `DELETE /v1/subscriptions/{id}`                     | Partial          | Create/list/retrieve subscriptions through the local checkout-completion state path. Update supports item replacement, metadata merge, and `cancel_at_period_end`. Delete performs immediate sandbox cancellation.                                                                                                                                                                                                                                                      |
| Subscription items      | `POST /v1/subscription_items`, `DELETE /v1/subscription_items/{id}`                                                                                                 | Partial          | Add or remove local subscription items for integration smoke paths. Billing proration and invoice recalculation are not modeled.                                                                                                                                                                                                                                                                                                                                        |
| Invoices                | `GET /v1/invoices`, `GET /v1/invoices/{id}`, `POST /v1/invoices/create_preview`                                                                                     | Partial          | List/retrieve invoices created by checkout. Preview returns a zero-value local smoke-test invoice.                                                                                                                                                                                                                                                                                                                                                                      |
| Payment intents         | `POST /v1/payment_intents`, `GET /v1/payment_intents`, `GET /v1/payment_intents/{id}`, `POST /v1/payment_intents/{id}/confirm`, `POST /v1/payment_intents/{id}/capture`, `POST /v1/payment_intents/{id}/cancel` | Partial          | Create/list/retrieve and mutate local payment intents. `confirm` supports deterministic sandbox PaymentMethod aliases such as `pm_card_visa`, `pm_card_visa_chargeDeclined`, and `pm_card_threeDSecure2Required`; manual capture moves through `requires_capture` before `capture` succeeds. This is a local state machine, not card processing or full PaymentIntent parameter parity.                                                                                 |
| Setup intents           | `POST /v1/setup_intents`, `GET /v1/setup_intents`, `GET /v1/setup_intents/{id}`, `POST /v1/setup_intents/{id}/confirm`, `POST /v1/setup_intents/{id}/cancel`        | Partial          | Create/list/retrieve and mutate local setup intents with deterministic success, decline, and authentication-required aliases. Mandates, bank-account verification, and full SCA behavior are not modeled.                                                                                                                                                                                                                                                               |
| Payment methods         | `GET /v1/payment_methods?customer={id}&type=card`                                                                                                                   | Partial          | Returns a deterministic sandbox card projection for a known customer. Create, attach, detach, and update are not supported.                                                                                                                                                                                                                                                                                                                                             |
| Webhook endpoints       | `POST /v1/webhook_endpoints`, `GET /v1/webhook_endpoints`, `GET /v1/webhook_endpoints/{id}`, `POST /v1/webhook_endpoints/{id}`, `DELETE /v1/webhook_endpoints/{id}` | Supported        | Manage local webhook endpoints. Secrets are generated when omitted and masked in API responses. `enabled_events` supports exact event names, `*`, and prefix wildcards such as `invoice.*`.                                                                                                                                                                                                                                                                             |
| Events                  | `GET /v1/events`, `GET /v1/events/{id}`                                                                                                                             | Supported        | List and retrieve Billtap-created events. Filters include `type` and `scenarioRunId`.                                                                                                                                                                                                                                                                                                                                                                                   |

## Billtap APIs

Base path: `/api`

| Area               | Endpoints                                                                                                                                                                                                                                                                                               | Scope                                                                                                                                        |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| Checkout           | `POST /api/checkout/sessions/{id}/complete`                                                                                                                                                                                                                                                             | Billtap-only checkout completion endpoint used by the hosted checkout UI and local tests.                                                    |
| Portal             | `GET /api/portal`, `GET /api/portal/customers/{id}`, `POST /api/portal/subscriptions/{id}/plan-change`, `POST /api/portal/subscriptions/{id}/seat-change`, `POST /api/portal/subscriptions/{id}/cancel`, `POST /api/portal/subscriptions/{id}/resume`, `POST /api/portal/customers/{id}/payment-method` | Sandbox portal state and actions. These update local billing state and timeline evidence; they do not claim Stripe Billing Portal parity.    |
| Dashboard evidence | `GET /api/objects`, `GET /api/timeline`, `GET /api/delivery-attempts`, `POST /api/debug-bundles`                                                                                                                                                                                                        | Object lists, timelines, delivery evidence, and debug bundle data for local investigation.                                                   |
| Webhook operations | `POST /api/events/{id}/replay`                                                                                                                                                                                                                                                                          | Replays an existing event and can schedule duplicate, delayed, and out-of-order attempts.                                                    |
| Fixtures           | `POST /api/fixtures/apply`, `GET /api/fixtures/snapshot`, `POST /api/fixtures/assert`                                                                                                                                                                                                                   | Data-driven setup and assertion APIs for customers, products, prices, subscription graphs, invoices, payment intents, and timeline evidence. |
| Scenarios          | `POST /api/scenarios/run`                                                                                                                                                                                                                                                                               | Runs a scenario JSON object or YAML payload and returns the scenario report.                                                                 |
| Boundary controls  | `GET /api/audit-log`, `POST /api/retention/apply`                                                                                                                                                                                                                                                       | Audit and retention controls for replay, delivery overrides, and raw evidence redaction.                                                     |

## Webhook Compatibility Claim

Billtap emits Stripe-style event envelopes for the supported checkout sequence
and stores delivery attempts for local debugging. The envelope includes event
IDs, event type, created time, `livemode: false`, `data.object`, request
metadata, and Billtap metadata.

Supported generic event types:

- `checkout.session.completed`
- `checkout.session.expired`
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`
- `invoice.created`
- `invoice.finalized`
- `invoice.payment_succeeded`
- `invoice.paid`
- `invoice.payment_failed`
- `invoice.voided`
- `payment_intent.created`
- `payment_intent.succeeded`
- `payment_intent.processing`
- `payment_intent.canceled`
- `payment_intent.payment_failed`
- `payment_intent.requires_action`
- `payment_intent.amount_capturable_updated`
- `setup_intent.created`
- `setup_intent.succeeded`
- `setup_intent.canceled`
- `setup_intent.setup_failed`
- `setup_intent.requires_action`

Current event boundaries:

- Checkout completion emits the generic checkout, subscription, invoice, and
  payment-intent sequence. Async-pending checkout emits
  `payment_intent.processing` without an invoice failure event; canceled
  checkout emits `checkout.session.expired`, `payment_intent.canceled`, and
  `invoice.voided`.
- Direct PaymentIntent and SetupIntent APIs emit local intent events for create,
  confirm, capture, cancel, failure, and requires-action states. These events
  are for local webhook/debug evidence and do not imply real payment processing.
- Portal subscription actions update local billing state and timeline evidence,
  and enqueue `customer.subscription.updated` or
  `customer.subscription.deleted` when the subscription changes. Portal payment
  method simulation remains Billtap evidence and does not claim Stripe Billing
  Portal parity.
- Replay keeps the original event ID and payload, then creates new delivery
  attempts with replay metadata.
- Duplicate delivery reuses the event ID and payload.
- Delay changes delivery scheduling, not event creation time.
- Out-of-order delivery changes attempt ordering evidence, not canonical event
  sequence.

Delivery headers use:

```text
Billtap-Signature: t=<unix_seconds>,v1=<hex_hmac_sha256>
```

The default header name is Billtap-specific. Set
`BILLTAP_WEBHOOK_SIGNATURE_HEADER=Stripe-Signature` when an application already
verifies Stripe's standard webhook header and should consume Billtap through the
same receiver path.

Webhook envelopes emit a Stripe API version in `api_version` so Stripe SDK
webhook deserializers can hydrate `data.object` into typed models. The default
is `2025-12-15.clover`; set `BILLTAP_WEBHOOK_API_VERSION` when the application
under test pins a different Stripe SDK/API version.

## Scenario Claim

Billtap scenarios are deterministic YAML/JSON flows for local and CI tests.

Supported generic actions:

- `customer.create`
- `product.create`
- `price.create`
- `checkout.create`
- `checkout.complete`
- `clock.advance`
- `invoice.retry`
- `app.assert`

Supported SaaS profile webhook evidence actions:

- `webhook.deliver_duplicate`
- `webhook.deliver_out_of_order`

Supported SaaS profile actions are documented in
`specs/000-product/contracts/scenario.md` and exercised by
`examples/saas-adoption-contract.yml`.

Current scenario boundaries:

- `checkout.complete` mutates local billing state and emits checkout-related
  webhook evidence.
- `invoice.retry` is currently deterministic scenario evidence; it does not yet
  mutate the generic billing invoice/payment-intent state.
- Generic webhook replay is available through `POST /api/events/{id}/replay`
  and the `webhook.replay` scenario action. Replay can schedule duplicate,
  delayed, out-of-order, simulated endpoint status, timeout, generic transport
  error, and signature-mismatch delivery attempts.
- `webhook.deliver_duplicate` and `webhook.deliver_out_of_order` currently
  update SaaS profile webhook evidence. Use `webhook.replay` for generic HTTP
  delivery attempts.
- App assertions call the configured app assertion endpoint and can fail the
  run with a non-zero exit code.
- Scenario reports are JSON and Markdown capable from the CLI.

## Fixture Claim

Fixture packs support repeatable local setup and assertions for:

- customers
- catalog products
- catalog prices
- subscription graphs created through the normal checkout-completion path
- optional stable checkout session, subscription, invoice, and payment intent
  IDs for provider-replacement tests that need exact fixture IDs
- fixture-scoped snapshots
- assertion reports for customers, products, prices, checkout sessions,
  subscriptions, invoices, payment intents, and timeline entries

Fixture metadata is written to created objects:

- `billtap_fixture_name`
- `billtap_fixture_run_id`
- `billtap_fixture_namespace`
- `billtap_fixture_ref`

Fixtures are intended for local and CI setup. They are not migration tooling and
must not contain real card data, live credentials, or production customer data.

## Adoption Smoke Claim

`npm run smoke:sdk` exercises the documented Stripe-like subset with the
official `stripe-node` SDK. The lane covers customer, product, price, checkout
session, event, webhook endpoint, and related retrieve/list flows against an
isolated local Billtap server by default. It can target an existing Billtap
server with `BILLTAP_STRIPE_SDK_SMOKE_BASE_URL`.

## Unsupported Stripe Behavior

Billtap does not support or claim:

- Full Stripe API coverage.
- Real payment processing or live payment success paths.
- Stripe Dashboard behavior.
- Real card data, PAN, CVC, expiration fields, or live credentials.
- Full Stripe request idempotency-key semantics across all endpoints. Billtap
  only caches same-process `POST` responses for supported API simulation and
  rejects same-key parameter mismatches.
- Direct charges, refunds, disputes, payouts, transfers, balance transactions,
  accounts, Connect onboarding, mandates, tax, coupons,
  promotion-code redemption, discounts, subscriptions schedules, quotes, or
  usage-based metering.
- PaymentIntent customer-balance application, incremental authorization, bank
  microdeposit verification, and full payment-method option parity.
- Stripe-hosted Checkout or Billing Portal parity.
- Provider-specific settlement, risk, tax, invoice rendering, fraud, account,
  payout, or dispute behavior.
- Complete webhook event coverage.
- Direct invoice finalize/pay/void endpoints.

Use Stripe testmode or the real provider sandbox as the fallback lane for these
behaviors.

## Compatibility Change Rules

Before adding a new compatibility claim:

1. Add or update a fixture, scenario, or contract test.
2. Document the endpoint, event, action, or boundary here.
3. State unsupported provider behavior instead of silently approximating it.
4. Run the release verification in `docs/RELEASE.md`.
