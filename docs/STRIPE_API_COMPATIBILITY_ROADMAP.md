# Stripe API Compatibility Roadmap

Status date: 2026-05-10

This roadmap extends Billtap beyond one SaaS adoption path. The goal is not to
be a real payment processor or a complete Stripe clone. The goal is to make
Stripe compatibility measurable, versioned, and progressively useful across
common integration shapes while keeping Billtap's stateful local billing-lab
strength.

## Reference Model

Billtap compatibility work should be grounded in public Stripe references:

- Stripe API reference: resource-oriented REST, form-encoded request bodies,
  JSON responses, standard HTTP status codes, authentication, and API versioned
  behavior.
- `stripe/openapi`: machine-readable endpoint, schema, event, fixture,
  expandable-field, and resource-id reference. New inventory work should prefer
  the `/latest/` specs because they cover v1 and v2 endpoints/events together.
- `stripe/stripe-mock`: useful request route, parameter, and type validation
  oracle. It is stateless and fixture-backed, so it is not a behavioral oracle
  for billing state, declines, retries, webhooks, or app-side assertions.
- Stripe webhooks documentation: delivery failures, retries, duplicate
  delivery, signature verification, connected-account destinations, and thin vs
  snapshot event expectations.
- Stripe SDKs and samples: adoption smoke targets for real client behavior.

## Compatibility Levels

Every endpoint, event, or object family should have exactly one published
level. The level can move up only with tests, fixtures, and documentation.

| Level | Name              | Meaning                                                                                         |
| ----- | ----------------- | ----------------------------------------------------------------------------------------------- |
| L0    | Inventory only    | Known from OpenAPI/docs, tracked in the compatibility matrix, not implemented.                   |
| L1    | Schema validated  | Route, method, required params, types, enums, unknown params, and error envelopes are tested.    |
| L2    | Fixture response  | Returns Stripe-shaped objects from deterministic fixtures, with list/retrieve pagination basics. |
| L3    | Stateful local    | Mutates Billtap storage and can be retrieved/listed consistently across requests.                |
| L4    | Scenario capable  | Can be driven by fixtures/scenario runner and produces timeline/debug evidence.                  |
| L5    | Webhook modeled   | Emits documented event sequences and delivery attempts with replay/failure evidence.             |
| L6    | SDK/adoption pass | Covered by official Stripe SDK smoke and at least one app-style integration test.                |

Supported public claims require L3 or higher unless explicitly documented as
stateless/schema-only compatibility.

## Product Direction

Billtap should have two compatibility surfaces:

1. **Broad Stripe API surface:** OpenAPI-driven inventory, request validation,
   fixture responses, error envelopes, pagination, expand, and SDK smoke.
2. **Deep billing lab surface:** stateful subscription/payment behavior,
   deterministic failures, webhooks, local clock, scenarios, dashboard evidence,
   and app assertions.

This keeps the project honest: broad compatibility helps SDKs and integrations
start, while deep compatibility is reserved for flows where Billtap can provide
real local testing value.

## Endpoint Family Priorities

| Priority | Family                         | Target level | Rationale                                                                                 |
| -------- | ------------------------------ | ------------ | ----------------------------------------------------------------------------------------- |
| P0       | Protocol baseline              | L3           | Auth, errors, form parsing, metadata, pagination, expand, idempotency, version headers.   |
| P0       | Core billing catalog           | L4-L6        | Customers, products, prices, checkout, portal, subscriptions, invoices, payment intents.  |
| P0       | Webhooks and events            | L5-L6        | Billtap's strongest differentiator is delivery evidence and app-side debugging.           |
| P1       | Subscription lifecycle depth    | L4-L6        | Schedules, invoice items, trials, coupons, discounts, credit notes, test clocks, meters.  |
| P1       | Payment method lifecycle        | L4-L6        | PaymentMethods, SetupIntents, saved cards, failed setup, mandates, SCA-style transitions. |
| P1       | Payment/refund/dispute surface  | L3-L5        | Charges, refunds, balance transactions, disputes, payment history, refund webhooks.       |
| P1       | Entitlements and metering       | L4-L6        | Modern subscription apps use entitlements, features, usage events, and meter summaries.   |
| P2       | Connect/platform smoke          | L2-L5        | Accounts, transfers, application fees, connected-account webhook routing, payout evidence. |
| P2       | Tax and invoice rendering smoke | L2-L4        | Tax rates/codes/calculation-like fixtures enough for subscription app tests.              |
| P3       | Low-state auxiliary resources   | L1-L2        | Files, reporting, balance, country/spec resources, webhook destinations, search stubs.    |
| P3       | Risk, Issuing, Treasury, Atlas  | L0-L2        | Track inventory and schema fixtures; do not model real financial behavior.                |

## Roadmap Phases

### S0: OpenAPI Inventory And Coverage Matrix

Output:

- Import a pinned Stripe OpenAPI snapshot into generated test fixtures or
  checked metadata.
- Generate `dist/compatibility/stripe-api-inventory.json` and Markdown matrix.
- Track every path, method, resource id, event type, expandable field, and
  current Billtap level.
- Add drift detection when Stripe OpenAPI changes.

Gate:

- Inventory generation is deterministic.
- Every implemented Billtap `/v1` route maps to an OpenAPI path or a documented
  Billtap-specific exception.
- Unsupported endpoints are visible instead of implicit.

Current command:

```bash
go run ./cmd/billtap compatibility inventory --openapi path/to/openapi.spec3.json --output-dir dist/compatibility
```

### S1: Protocol Compatibility Baseline

Output:

- Common request parser for Stripe form params, JSON fallback, nested arrays,
  and `expand[]`.
- Shared Stripe-like error catalog for `invalid_request_error`, `card_error`,
  `idempotency_error`, `api_error`, and future rate-limit/server categories.
- Pagination and list envelope normalization.
- Request-id, API version, metadata, livemode, idempotency, and trace evidence.

Gate:

- Cross-endpoint protocol tests pass for every supported POST and list route.
- Scorecard includes protocol cases separate from resource behavior cases.

### S2: OpenAPI-Backed Validation And Fixture Responses

Output:

- Generated validators for broad L1 coverage, with hand-written overrides for
  stateful local behavior.
- Fixture responses for low-state retrieve/list/create smoke paths.
- Optional differential lane against `stripe-mock` for route/parameter/type
  sanity checks.

Gate:

- L1 claims have success, required-param, wrong-type, unknown-param, enum, and
  missing-resource cases.
- Differential mismatches are classified as Billtap bug, Stripe-mock limitation,
  or accepted Billtap boundary.

### S3: Billing Lifecycle Depth

Output:

- Subscription schedules, trial transitions, renewal invoices, proration
  previews, invoice items, coupons, promotion codes, discounts, credit notes,
  test clocks, usage records/meters, and entitlements.
- Local clock drives renewal, trial end, dunning, retry, cancellation, and
  scheduled phase changes.
- Scenario examples for upgrade, downgrade, renewal, failed renewal, retry,
  coupon application, discount expiration, credit note, and entitlement change.

Gate:

- Stateful billing graph remains explainable in the dashboard timeline.
- Every mutating billing behavior has event sequence tests and scenario reports.

### S4: Payment Intent And Payment Method Depth

Output:

- Direct `PaymentIntent` create/confirm/capture/cancel paths.
- Direct `SetupIntent` create/confirm/cancel paths.
- PaymentMethod create/attach/detach/list/update simulations without real card
  data.
- Deterministic outcomes for success, decline, processing, authentication
  required, missing payment method, and async completion.

Gate:

- SDK smoke covers direct payments and setup flows.
- No PAN/CVC storage path is introduced.

### S5: Refunds, Disputes, Balance, And Payment History

Output:

- Charges, refunds, credit notes, balance transactions, disputes, and related
  webhook events sufficient for SaaS support and accounting smoke tests.
- Scenario flows for refund created/succeeded/failed, charge refunded, dispute
  opened/closed, and invoice/payment history reconciliation.

Gate:

- Payment history can be asserted from API state, webhook evidence, and debug
  bundle output.

### S6: Connect And Platform Simulation

Output:

- Account, account link/session, external account, transfer, application fee,
  payout, and connected-account event fixtures.
- Connected-account webhook destination routing and signature evidence.
- Safe stubs for onboarding/KYC without modeling real identity verification.

Gate:

- Platform apps can test account routing, webhook scope, and reconciliation
  without depending on live Stripe testmode.

### S7: SDK And Sample Compatibility Lanes

Output:

- Official SDK smoke lanes for `stripe-node`, `stripe-go`, `stripe-java`,
  `stripe-python`, and `stripe-ruby`.
- Sample-app lanes for checkout, subscriptions, portal, webhooks, payment
  intents, setup intents, refunds, and Connect.
- Reports stored as JSON/Markdown artifacts with failing request replay data.

Gate:

- Each SDK lane names covered endpoint families and unsupported cases.
- Fast smoke remains suitable for CI; full matrix can be scheduled/manual.

### S8: Optional Live Drift And Oracle Runs

Output:

- Manual or scheduled jobs that compare selected cases with Stripe testmode,
  Stripe OpenAPI, and stripe-mock.
- No live Stripe calls in normal CI.
- Redacted artifacts that explain accepted divergences.

Gate:

- Drift reports update docs before compatibility claims change.
- Live credentials are never required for local development or public tests.

## Compatibility Matrix Schema

The generated matrix should be machine-readable so agents can answer coverage
questions quickly.

```json
{
  "api_version": "2025-12-15.clover",
  "source": "stripe/openapi latest",
  "generated_at": "2026-05-10T00:00:00Z",
  "resources": [
    {
      "family": "billing",
      "resource": "subscription",
      "path": "/v1/subscriptions/{subscription}",
      "method": "post",
      "stripe_resource_id": "subscription",
      "billtap_level": "L3",
      "stateful": true,
      "webhook_events": ["customer.subscription.updated"],
      "scorecard_cases": ["subscription.update.items.price"],
      "sdk_smoke": ["stripe-node"],
      "docs": "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
      "risks": ["no proration invoice yet"]
    }
  ]
}
```

## Scorecard Expansion

The current public scorecard remains the release gate for the documented
subset. Add broader scorecards without weakening that gate:

- `l3-public-readiness`: current release-blocking local subset.
- `stripe-api-inventory`: path/resource/event inventory; no runtime claim.
- `stripe-schema-validation`: L1 route/param/type validation.
- `stripe-fixture-shapes`: L2 response object and list envelope shapes.
- `stripe-stateful-billing`: L3-L5 billing state and event sequences.
- `stripe-sdk-smoke`: official SDK behavior across supported families.
- `stripe-oracle-optional`: manual/scheduled stripe-mock or live testmode
  comparison.

Each report should preserve the existing buckets:

- `imported`
- `skipped`
- `unsupported`
- `mismatch`
- `error`

## PR Chunk Plan

Each chunk should use the PR -> review -> fix -> merge workflow.

| Chunk | Ownership                 | Output                                                                           | Verification                                      |
| ----- | ------------------------- | -------------------------------------------------------------------------------- | ------------------------------------------------- |
| S0-A  | Compatibility inventory   | OpenAPI snapshot loader, inventory generator, initial matrix docs                | Go tests or Node script tests, generated diff     |
| S0-B  | Coverage docs             | Endpoint family matrix, compatibility levels, unsupported inventory              | Markdown review, link check                       |
| S1-A  | Protocol parser           | Common form parser, nested param tests, expand parsing                           | Go API tests                                      |
| S1-B  | Error/idempotency         | Error catalog hardening, request-id, idempotency trace expansion                 | API tests and scorecard cases                     |
| S2-A  | Generated validation P0   | OpenAPI-derived validators for catalog, customers, checkout, subscriptions       | Scorecard plus stripe-mock optional oracle        |
| S2-B  | Fixture response harness  | Low-state fixture response engine for inventory-only endpoints                   | Fixture-shape scorecard                           |
| S3-A  | Renewal/test clock        | Renewal invoices, trial end, retry mutation, local clock event sequences         | Scenario reports and webhook tests                |
| S3-B  | Discounts/credits/meters  | Coupons, promotion codes, discounts, credit notes, usage/meter events            | Scenario and billing engine tests                 |
| S4-A  | Direct payment intents    | Create/confirm/capture/cancel state machine and failure aliases                  | SDK smoke and API tests                           |
| S4-B  | Setup/payment methods     | SetupIntent and PaymentMethod attach/detach/update lifecycle                     | SDK smoke and no-card-data boundary tests         |
| S5-A  | Refund/dispute history    | Refunds, disputes, balance transactions, support/debug evidence                  | Scenario reports and dashboard/debug bundle tests |
| S6-A  | Connect smoke             | Connected-account fixtures, transfer/application fee/payout evidence             | Webhook routing tests and sample platform smoke   |
| S7-A  | SDK matrix                | Node, Go, Java, Python, Ruby smoke lanes and reports                             | CI/manual workflow artifacts                      |
| S8-A  | Optional oracle           | stripe-mock/live-testmode drift runner, redacted report format, accepted diffs   | Manual workflow only                              |

## Definition Of Done For A New Family

Before a family moves beyond inventory:

1. `docs/COMPATIBILITY.md` names the exact endpoint level and unsupported
   provider behavior.
2. A scorecard case proves the request and response shape.
3. Mutating behavior has storage tests and timeline/debug evidence.
4. Event-producing behavior has webhook envelope and delivery attempt tests.
5. Scenario-capable behavior has JSON/Markdown report output.
6. SDK-adoption behavior has at least one official SDK smoke path.
7. The public docs avoid implying real payment processing, live credential use,
   or complete Stripe parity.

## Near-Term Recommendation

Start with S0 and S1 before adding many endpoints. Without an inventory and
protocol baseline, broad compatibility will become a pile of ad hoc routes.
After S0/S1, the highest-value implementation order is:

1. Renewal/test clock/retry mutation for subscriptions and invoices.
2. Direct PaymentIntent and SetupIntent lifecycle.
3. Coupons, discounts, credit notes, refunds, and payment history.
4. OpenAPI-derived L1 validation for broad low-state endpoint coverage.
5. Official SDK matrix and optional stripe-mock oracle lane.
