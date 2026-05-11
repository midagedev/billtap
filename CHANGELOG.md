# Changelog

## Unreleased

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
