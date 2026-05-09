# Billtap

Full-stack Stripe-style billing sandbox for local development, CI scenarios, and controlled staging checks.

`billtap` is designed for the fast path in subscription billing work:

- use `billtap` as the default lane for deterministic checkout, portal, webhook, fixture, and scenario tests
- keep Stripe testmode or provider sandboxes as the high-fidelity fallback lane

It provides a Go backend, React checkout/portal/dashboard UIs, a practical Stripe-compatible API subset, webhook reliability controls, YAML scenarios, fixture apply/snapshot/assert APIs, and a generic `saas` workspace profile. It is not a real payment processor and must not be used in a live payment path.

## Why This Project Exists

Stripe testmode is useful, but routine billing tests often need a different tradeoff:

- deterministic object state instead of shared remote testmode state
- controllable webhook order, retry, duplication, delay, and replay
- local checkout and portal flows that work in CI
- fixture-backed setup and structured assertions
- one timeline that explains API calls, billing state, webhooks, delivery attempts, and app callbacks

Billtap makes billing flows reproducible before they reach real payment infrastructure.

For the exact supported surface, see `docs/COMPATIBILITY.md`. Anything not
listed there should be treated as unsupported until it has fixture-backed tests
and documentation.

## Recommended Adoption Model

Use two lanes instead of forcing one tool to satisfy every billing test:

- default lane: `billtap` for local development, CI regression scenarios, fixture setup, and webhook reliability tests
- fallback lane: Stripe testmode or the real provider sandbox for provider-specific behavior, dashboard behavior, hosted-payment-provider parity, and final compatibility checks

This keeps the common suite fast without pretending a local sandbox is full Stripe parity.

## Value For Subscription Teams

Typical gains when replacing ad hoc mocks or always-on provider sandbox tests for supported scenarios:

- faster local and CI feedback
- deterministic webhook and retry behavior
- fixture packs that are easy to apply and assert
- reproducible subscription lifecycle scenarios
- visible billing timelines for debugging app entitlement bugs

This is most useful for SaaS products with checkout, billing portal, subscriptions, invoices, payment failures, seats, workspace entitlements, quota purchases, and webhook handlers.

## Billtap vs Stripe Testmode Only

| Topic | `billtap` | Stripe testmode only |
| --- | --- | --- |
| Startup/runtime overhead | Local process or Docker image | Remote provider dependency |
| Determinism | High for supported local semantics | Depends on remote account state and async provider behavior |
| Webhook failure controls | Built-in duplicate, delay, out-of-order, replay, retry | Requires extra tooling and manual setup |
| Fixture setup | JSON/YAML fixture apply/snapshot/assert APIs | Usually custom scripts against remote state |
| Provider fidelity | Partial by design | High for Stripe-specific behavior |
| Best fit | Fast development and CI regression lane | High-fidelity provider compatibility lane |

Recommended strategy:

- default profile: Billtap-backed local or CI scenarios
- high-fidelity profile: Stripe testmode/provider sandbox for unsupported or provider-specific behavior

## Why This Works As The Default Lane

Billtap is optimized around test ergonomics, not full payment-provider emulation:

- webhook replay, duplicate delivery, delay, and out-of-order controls are explicit Billtap operations
- fixture packs can seed customers, catalog entries, subscriptions, and expected assertions
- reports are structured for CI failures
- dashboard and API evidence is redacted by default where secrets can appear
- unsupported provider behavior stays documented instead of being silently approximated

## When To Use / Not Use

| Decision | Use when... |
| --- | --- |
| Use | You need deterministic subscription billing tests in local dev or CI |
| Use | You need to validate webhook idempotency, retries, duplicate delivery, delays, or replay |
| Use | You want fixture-driven setup with easy snapshot/assert APIs |
| Use | You need a local checkout/portal/dashboard loop for app integration work |
| Do not use | You need full Stripe API or hosted Stripe Dashboard parity |
| Do not use | You are handling real card data, real customer payment data, or production payment success paths |
| Do not use | Your test must prove provider-specific settlement, risk, tax, invoice rendering, or account behavior |

## Quick Start

Current distribution state: source-only. No package, Homebrew formula, or published Docker image is released yet.

Requirements:

- Go 1.25+
- Node.js and npm
- Docker, optional for image smoke checks

Build and run locally:

```bash
npm install
npm run build
go run ./cmd/billtap
```

Open:

```text
http://localhost:8080
```

Run a scenario with the sample app assertion endpoint:

```bash
PORT=3300 npm --prefix examples/sample-app start
```

In another terminal:

```bash
go run ./cmd/billtap scenario run examples/subscription-payment-retry.yml \
  --report-json billtap-report.json \
  --report-md billtap-report.md
```

Run the generic SaaS workspace scenario:

```bash
go run ./cmd/billtap scenario run examples/saas-adoption-contract.yml
```

Build a local image:

```bash
docker build -t billtap:local .
docker run --rm -p 8080:8080 -v billtap-data:/data billtap:local
```

## Fixture And Assertion APIs

Billtap includes local integration-test helpers:

- `POST /api/fixtures/apply`: apply JSON/YAML customers, catalog, subscriptions, and assertions
- `GET /api/fixtures/snapshot`: read a filtered fixture-scoped billing snapshot
- `POST /api/fixtures/assert`: assert expected customer, product, price, subscription, invoice, payment intent, and timeline state

Fixture-applied subscriptions use the normal checkout-completion path so invoices, payment intents, checkout sessions, and timeline evidence stay consistent.

## Compatibility Snapshot

| Area | Current level | Notes |
| --- | --- | --- |
| Runtime | Go server with SQLite local default | In-memory storage exists for tests |
| Frontend | React checkout, portal, and dashboard apps | Built with Vite into `dist/app` |
| Stripe-like API | Practical local subset | Customers, products, prices, checkout sessions, subscriptions, invoices, payment intents, webhook endpoints, events, search/list projections used by tests |
| Webhooks | Signed delivery with reliability controls | Retry, duplicate, delay, out-of-order, replay, delivery evidence, redaction |
| Scenarios | YAML runner | Local clock, app assertions, JSON/Markdown reports, exit-code policy |
| Fixtures | Apply/snapshot/assert APIs | JSON/YAML input, fixture metadata isolation, structured pass/fail reports |
| SaaS profile | Generic workspace billing profile | Plans, seats, members, export quota, extra export, payment history, support bundle, platform/connect-style webhook evidence |
| Release state | Source-only | Local Docker image builds; no published image/package yet |

Detailed compatibility matrix: `docs/COMPATIBILITY.md`.

## Known Limitations

- This is not full Stripe compatibility.
- It does not process real payments and rejects real card-data fields.
- Hosted UI behavior is a sandbox approximation, not a Stripe-hosted UI clone.
- Dashboard access control is not a production security boundary.
- Relay mode is only for controlled testmode or staging-adjacent debugging and stores raw payloads as metadata-only.
- Published release automation is not present yet.

## Development Commands

```bash
go test ./...
npm run typecheck
npm run build
npm run smoke:sample
go build -o /tmp/billtap ./cmd/billtap
docker build -t billtap:local .
```

Scenario smoke, with `PORT=3300 npm --prefix examples/sample-app start` running:

```bash
go run ./cmd/billtap scenario run examples/subscription-payment-retry.yml
go run ./cmd/billtap scenario run examples/saas-adoption-contract.yml
```

## Documentation

- Documentation index: `docs/README.md`
- Product goal: `docs/FINAL_GOAL.md`
- Architecture: `docs/ARCHITECTURE.md`
- Compatibility: `docs/COMPATIBILITY.md`
- SaaS profile: `docs/SAAS_PROFILE.md`
- Testing: `docs/TESTING.md`
- Production boundaries: `docs/PRODUCTION_BOUNDARIES.md`
- Release process: `docs/RELEASE.md`
- Release checklist: `docs/RELEASE_CHECKLIST.md`
- Roadmap: `docs/ROADMAP.md`
- Scenario contract: `specs/000-product/contracts/scenario.md`
- API contract: `specs/000-product/contracts/api.md`
- Webhook contract: `specs/000-product/contracts/webhooks.md`
- Changelog: `CHANGELOG.md`
- Contributing: `CONTRIBUTING.md`
- Code of conduct: `CODE_OF_CONDUCT.md`
- Support: `SUPPORT.md`
- Security: `SECURITY.md`

## Community And Support

Before opening an issue or pull request, read `CONTRIBUTING.md` and `SUPPORT.md`. Security reports and production-boundary bypasses should follow `SECURITY.md` rather than public issues.

Public reports and examples must be sanitized. Do not include real card data, live credentials, production customer data, private company data, or production payment payloads.

## Release Process

Public release procedure: `docs/RELEASE.md`.

Maintainer checklist summary:

1. Run the development commands above.
2. Run the scenario smoke commands above.
3. Build and smoke the local Docker image.
4. Confirm `docs/COMPATIBILITY.md` matches the implemented API surface.
5. Confirm the public-surface scan is clean and `.private/` is ignored.
6. Tag the release as `vX.Y.Z`.
7. Publish the Docker image or package only after release automation is explicitly added.

## License

License is not declared yet. Add a `LICENSE` file before publishing or accepting external contributions.
