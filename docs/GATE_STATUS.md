# Gate Status

This is the public gate snapshot. Internal adoption evidence and raw handoff notes are preserved locally under `.private/`.

| Gate                          | Status         | Public evidence                                                                                                                                             |
| ----------------------------- | -------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| G0 Spec Readiness             | Passed         | Constitution, product spec, plan, tasks, gates, and docs exist                                                                                              |
| G1 Runtime Contract           | Passed         | Go server, config loading, SQLite storage, in-memory test storage, health endpoints                                                                         |
| G2 Checkout MVP               | Passed         | Customer/product/price/checkout APIs, checkout UI, checkout completion state                                                                                |
| G3 Webhook Reliability        | Passed         | Signed delivery, attempts, retry, duplicate, delay, out-of-order, replay                                                                                    |
| G4 Debuggability              | Passed         | Dashboard objects, request traces, timeline, webhook detail, app response evidence, diagnostic/debug bundles                                                |
| G5 CI Contract                | Passed         | YAML scenario runner, local clock, app assertions, JSON/Markdown reports, exit codes                                                                        |
| G6 Portal Coverage            | Passed         | Portal UI, plan change, seat change, cancellation, resume, payment-method simulation, invoice history                                                       |
| G7 SaaS Workspace Profile     | Passed         | Generic `saas` scenario pack for workspace plans, seats, members, export quota, extra export payment, payment history, support bundle, and webhook evidence |
| G8 Production Boundary        | Passed         | Secret masking, card-data rejection, relay metadata-only raw payload storage, retention redaction, audit log                                                |
| G9 Release Candidate          | Passed locally | Dockerfile, sample app, public examples, release checklist                                                                                                  |
| G10 Fixture Integration Smoke | Passed locally | Fixture apply/snapshot/assert APIs support deterministic integration setup                                                                                  |
| G11 Assertion Ergonomics      | Passed locally | Structured pass/fail fixture assertions and fixture-scoped snapshots                                                                                        |
| G12 Public Release Readiness  | Passed locally | Public claims are tied to tests/scorecard cases; scorecard corpus has 30 release-blocking cases; Apache-2.0 `LICENSE` and `NOTICE` are present              |
| G13 Stripe API Expansion      | In progress    | Roadmap defines compatibility levels and endpoint-family priorities; OpenAPI inventory generator and optional workflow write JSON/Markdown coverage artifacts  |
| G14 Stripe API 90% Program    | In progress    | `docs/STRIPE_COMPATIBILITY_90_TARGET.md` defines 90% L1+ target, current 37/619 baseline, family thresholds, and chunk plan                              |

## Current Public Claim

Billtap is a source-first local billing sandbox with a GHCR container image. It
can be built, tested, run, and smoke-checked from this repository, and adoption
repos can pull `ghcr.io/midagedev/billtap:main` or a release tag. It is licensed
under Apache-2.0.

## Current Compatibility Evidence

- Scorecard version: `l3-public-readiness-v4`
- Release-blocking scorecard cases: 30
- Required scorecard release result: `mismatch=0`, `error=0`, `passed=true`
- Coverage focus: request validation, protocol parameter acceptance,
  idempotency mismatch, deterministic checkout payment-error aliases
- OpenAPI operation baseline: `37 / 619`, `6.0%`
- Long-running OpenAPI operation target: at least `558 / 619`, `90.0%`, at
  `L1+` with deeper P0/P1 behavior gates

## Last Local Code Verification

Verified on 2026-05-10 from branch `codex/stripe-protocol-baseline`:

- `go test ./...`
- `go run ./cmd/billtap compatibility scorecard --output-dir /tmp/billtap-compatibility`
  - result: `imported=30 skipped=1 unsupported=1 mismatch=0 error=0`
- `npm run typecheck`
- `npm run build`

Release verification should still be rerun on the final release branch or tag.

## Required Release Verification

```bash
go test ./...
go run ./cmd/billtap compatibility scorecard --output-dir /tmp/billtap-compatibility
npm run typecheck
npm run build
npm run smoke:sample
npm run smoke:sdk
npm run smoke:web:install
npm run smoke:web
go build -o /tmp/billtap ./cmd/billtap
docker build -t billtap:local .
PORT=3300 npm --prefix examples/sample-app start
go run ./cmd/billtap scenario run examples/subscription-payment-retry.yml
curl -fsS -X POST http://127.0.0.1:3300/test/reset
go run ./cmd/billtap scenario run examples/saas-adoption-contract.yml
```

The scenario commands require the sample app assertion endpoint to be running on port `3300`.
