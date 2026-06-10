# Changelog

## Unreleased

- Hosted pages now repoint caller-provided localhost redirect targets at the
  run's configured public origin: when a run has a `public_base_url`, the
  hosted checkout "Return to app" link and billing portal return
  link/redirect swap the scheme/host/port of `localhost`/`127.0.0.1`
  `success_url`/`return_url` values for the run origin (path and query kept),
  surfaced via the `billtap_return_url` extension field and the portal URL
  query while stored sessions, `success_url`, and portal `return_url`
  response fields keep the caller's original values. External domains and
  unconfigured runs are untouched.
- Added run-scoped public base URLs so several proxied stacks can share one
  Billtap server: `POST /runs/<runId>/v1/config` pins `public_base_url` and an
  optional `public_base_path` per run, and absolute URLs (checkout
  `session.url`, billing portal URLs) prefer the run's base over an
  `X-Billtap-Public-Base-Url` request header, the forwarded proxy origin
  (`X-Forwarded-Proto`/`Host`/`Prefix`) on run-scoped requests, and the global
  `BILLTAP_PUBLIC_BASE_URL`. The default run keeps its previous behaviour.
- Added multi-workspace support so one running server can hold several fully
  isolated billing datasets. Requests select a workspace with the
  `X-Billtap-Workspace` header or `workspace` query parameter, unselected
  requests keep using the backward-compatible `default` workspace, named
  workspaces open their own SQLite database lazily under `workspaces/`, and
  `GET /workspaces` lists the known workspaces.
- Added a manual invoice-backed one-time payment flow for local SaaS usage
  charges, including `POST /v1/invoices`, `POST /v1/invoiceitems`,
  `POST /v1/invoices/{id}/finalize`, metadata preservation, expanded
  payment-intent evidence, and customer-level default outcomes for invoice pay.
- Fixed invoice pay requests so a submitted payment method no longer overrides
  a configured invoice-backed PaymentIntent outcome.
- Accepted Stripe-compatible `proration_date` on subscription updates and
  retained subscription update billing/proration parameters as local evidence.
- Added customer-level default PaymentIntent outcomes so fixture-seeded
  customers can drive confirmed one-time payment failures without changing
  app-created PaymentIntent metadata.
- Added deferred per-PaymentIntent outcome controls for one-time payment flows
  through `metadata[billtap_payment_intent_outcome]`, local create aliases, and
  `POST /api/payment_intents/{id}/outcome`.
- Hardened Stripe-like shape and validation for billing portal sessions and
  customer payment-method lists, including portal flow enum checks,
  PaymentMethod SDK fields, non-card filtering, and Stripe-style validation
  error envelopes.
- Added `GET /v1/prices/search` for the measured Stripe Search Query Language
  subset used by one-time price lookup paths: `active`, `type`, `lookup_key`,
  and metadata equality clauses joined by `AND`.
- Added customer fixture controls for empty or explicit payment-method lists so
  no-card local billing scenarios can remain deterministic until portal save.
- Expanded the Stripe-like simulation surface with hosted billing portal
  sessions/actions, local coupons and promotion codes, subscription schedules,
  SCA-required PaymentIntent callbacks, customer cash-balance funding for
  bank-transfer intents, dispute evidence, event filtering, endpoint-scoped
  delivery attempts, and grouped webhook replay controls.
- Added endpoint-scoped historical webhook replay so apps can catch up fixture
  events emitted before webhook endpoint registration.
- Licensed Billtap under Apache-2.0 and added a top-level `NOTICE`.
- Hardened public-readiness validation: JSON numeric request values now
  preserve decimal input for wrong-type rejection, subscription quantities no
  longer silently normalize invalid values, and update endpoints have explicit
  parameter validation.
- Expanded the compatibility scorecard to `l3-public-readiness-v2` with 28
  release-blocking cases covering request validation, idempotency mismatch, and
  deterministic checkout payment-error aliases.
- Added public release readiness documentation and clarified release evidence
  requirements.
- Added public compatibility and release process docs for the v0.1.0 release
  path.
- Expanded Connect compatibility with platform account retrieval, account
  deletion markers, people/persons evidence, and updated OpenAPI inventory
  tracking to `110 / 587` operations.
- Documented the supported Stripe-like subset, Billtap-specific APIs,
  unsupported provider behavior, and fixture/scenario/webhook claim boundaries.
- Prepared the repository for a public source-only release.
- Reworked public documentation around the Billtap default testing lane and Stripe testmode fallback lane.
- Moved company-specific adoption material and raw validation notes into ignored `.private/` storage.
- Sanitized the public workspace billing profile as `saas`.
- Added public docs for contribution, security, gate status, and documentation navigation.

## 0.0.0

- Initial source state with Go backend, React checkout/portal/dashboard apps, scenario runner, fixture APIs, webhook reliability controls, and Dockerfile.
