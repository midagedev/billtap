# Contributing

Billtap is currently source-only and pre-release.

Before contributing, read:

- `README.md` for the product boundary and adoption model
- `SUPPORT.md` for issue, support, and security routing
- `SECURITY.md` for vulnerability and production-boundary reports
- `CODE_OF_CONDUCT.md` for community expectations

## Development Setup

```bash
npm install
npm run build
go test ./...
```

## Issue And Pull Request Flow

Use the GitHub issue templates for bugs, feature requests, and public production-boundary questions. Use the private security process for vulnerabilities, secret exposure, webhook-signature bypasses, or real payment data handling regressions.

For pull requests:

- keep the change scoped
- describe the behavior change and the boundary impact
- add or update tests when billing state, webhook behavior, fixtures, assertions, relay mode, masking, retention, or compatibility behavior changes
- include sanitized screenshots, reports, or logs when they help review
- do not edit unrelated files or revert changes from other contributors

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

Documentation-only changes should still be reviewed for clear routing, no-real-payment language, and links that work from the public repository.

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

Out of scope:

- live payment processing
- storing or validating real card data
- production payment success paths that depend on Billtap
- full provider API parity without a documented, tested subset
- private company-specific adoption material in public docs
