# Contributing

Billtap is currently source-only and pre-release.

## Development Setup

```bash
npm install
npm run build
go test ./...
```

## Before Sending Changes

Run:

```bash
go test ./...
npm run typecheck
npm run build
npm run smoke:sample
go run ./cmd/billtap scenario run examples/subscription-payment-retry.yml
```

The scenario command expects the sample app to be running on port `3300`.

For changes touching webhooks, signatures, retries, idempotency, fixture behavior, or billing state, add or update tests.

## Public Repo Hygiene

Do not commit:

- credentials or API keys
- real customer data
- private company names or internal repository names
- internal hostnames or private IP addresses
- local absolute paths
- raw evidence artifacts

Company-specific adoption notes should stay under `.private/`, which is intentionally ignored.

## Scope

Billtap is a local billing sandbox. Contributions should keep the product boundary clear: no real payment processing, no production payment dependency, and no claim of full Stripe parity without fixture-backed tests and documentation.
