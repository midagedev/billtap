# Testing Strategy

## Test Philosophy

Billing behavior is a state machine plus asynchronous delivery. Tests must cover both.

## Test Layers

### Unit tests

- billing state transitions
- invoice calculations
- local clock behavior
- webhook signature generation
- idempotency handling
- retry scheduling
- scenario parser
- assertion evaluation

### Contract tests

Fixture-backed compatibility tests for the Stripe-like API subset:

- customers
- products
- prices
- checkout sessions
- subscriptions
- invoices
- payment intents
- events
- webhook endpoints

Fixture ergonomics for integration tests:

- apply a data-driven fixture pack through `POST /api/fixtures/apply`
- snapshot fixture-scoped state through `GET /api/fixtures/snapshot`
- assert expected objects through `POST /api/fixtures/assert`
- keep fixture IDs stable for customer, product, and price setup
- use fixture `runId`, `namespace`, `tenantId`, and `ref` metadata to isolate repeated local/CI runs

### UI tests

Browser tests for:

- checkout success
- checkout failure
- checkout cancel
- portal plan change
- portal cancellation
- dashboard timeline
- webhook replay
- scenario detail

Run the release smoke for built UI routes with:

```bash
npm run smoke:web
```

The smoke command builds the React assets and Billtap Go binary, starts Billtap with an isolated temporary SQLite database, then checks `/app/dashboard/`, `/app/checkout/`, and `/app/portal/` in Chromium for key page text and JavaScript console errors. On a fresh CI runner, install the Chromium browser once with `npm run smoke:web:install`.

### Scenario tests

Run scenario YAML files end to end:

- create customer
- checkout subscription
- receive webhook
- app callback succeeds
- app assertion passes
- SaaS workspace subscription, seat, member, export, and payment history assertions pass

### Reliability tests

- duplicate webhook
- delayed webhook
- out-of-order webhook
- retry after 500
- idempotency conflict
- signature mismatch
- app timeout
- duplicate entitlement grant prevention
- out-of-order subscription event does not regress workspace status
- extra export provision retry does not double-count export quota

### SaaS profile tests

- free to paid checkout success
- abandoned checkout leaves subscription processing
- paid upgrade preview and confirm
- cancellation scheduled and resumed
- additional seat estimate and purchase
- member invite blocked by seat limit
- normal export consumes quota
- export blocked when quota is exhausted
- extra export payment and provision success
- payment failure then retry success
- back-office plan/seat/period changes
- refund changes payment history and export entitlement
- Connect webhook resolved and unresolved paths
- support debug bundle contains workspace, subscription, seat, export, payment, webhook, and assertion evidence
- expected billing/export observability event names and attributes are generated or exposed to a collector adapter

### Security tests

- secrets masked in dashboard
- webhook signing secret hidden
- request headers redacted
- no real card data persisted
- production relay mode does not store raw payload by default

## CI Gates

- backend unit tests pass
- frontend tests pass
- contract fixtures pass
- scenario replay tests pass
- dashboard smoke test passes
- production boundary tests pass before production-facing features are marked supported

## Exit Codes

- `0`: scenario passed
- `1`: scenario assertion failed
- `2`: invalid config
- `3`: app callback failed
- `4`: Billtap runtime failure
