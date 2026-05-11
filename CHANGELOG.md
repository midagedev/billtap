# Changelog

## Unreleased

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
