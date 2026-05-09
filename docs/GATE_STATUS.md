# Gate Status

This is the public gate snapshot. Internal adoption evidence and raw handoff notes are preserved locally under `.private/`.

| Gate | Status | Public evidence |
| --- | --- | --- |
| G0 Spec Readiness | Passed | Constitution, product spec, plan, tasks, gates, and docs exist |
| G1 Runtime Contract | Passed | Go server, config loading, SQLite storage, in-memory test storage, health endpoints |
| G2 Checkout MVP | Passed | Customer/product/price/checkout APIs, checkout UI, checkout completion state |
| G3 Webhook Reliability | Passed | Signed delivery, attempts, retry, duplicate, delay, out-of-order, replay |
| G4 Debuggability | Passed | Dashboard objects, timeline, webhook detail, app response evidence, debug bundle |
| G5 CI Contract | Passed | YAML scenario runner, local clock, app assertions, JSON/Markdown reports, exit codes |
| G6 Portal Coverage | Passed | Portal UI, plan change, seat change, cancellation, resume, payment-method simulation, invoice history |
| G7 SaaS Workspace Profile | Passed | Generic `saas` scenario pack for workspace plans, seats, members, export quota, extra export payment, payment history, support bundle, and webhook evidence |
| G8 Production Boundary | Passed | Secret masking, card-data rejection, relay metadata-only raw payload storage, retention redaction, audit log |
| G9 Release Candidate | Passed locally | Dockerfile, sample app, public examples, release checklist |
| G10 Fixture Integration Smoke | Passed locally | Fixture apply/snapshot/assert APIs support deterministic integration setup |
| G11 Assertion Ergonomics | Passed locally | Structured pass/fail fixture assertions and fixture-scoped snapshots |

## Current Public Claim

Billtap is a source-only local billing sandbox. It can be built, tested, run, and smoke-checked from this repository, but it is not yet published as a package or Docker image.

## Required Release Verification

```bash
go test ./...
npm run typecheck
npm run build
npm run smoke:sample
go build -o /tmp/billtap ./cmd/billtap
docker build -t billtap:local .
go run ./cmd/billtap scenario run examples/subscription-payment-retry.yml
go run ./cmd/billtap scenario run examples/saas-adoption-contract.yml
```

The scenario commands require the sample app assertion endpoint to be running on port `3300`.
