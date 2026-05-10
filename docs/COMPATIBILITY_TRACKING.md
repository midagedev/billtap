# Compatibility Tracking

Status date: 2026-05-10

This document defines how Billtap tracks Stripe API compatibility as a
measurable, long-running body of work. The public claim is the generated
inventory plus the tests behind each implemented level, not a hand-maintained
checklist.

## Source Of Truth

Generate the matrix from a pinned or downloaded Stripe OpenAPI spec:

```bash
go run ./cmd/billtap compatibility inventory \
  --openapi path/to/openapi.spec3.json \
  --output-dir dist/compatibility \
  --source stripe/openapi-latest-public
```

The command writes:

- `stripe-api-inventory.json`: machine-readable source for agents and CI.
- `stripe-api-inventory.md`: reviewable compatibility table for humans.

`dist/` is intentionally ignored. PRs should quote the before/after summary
instead of committing generated artifacts.

## Measured Fields

The inventory summary is the compatibility scoreboard:

- `summary.total_operations`: Stripe OpenAPI operations seen by the generator.
- `summary.implemented_operations`: operations with a tested Billtap runtime
  claim.
- `summary.inventory_only_operations`: known Stripe operations still at `L0`.
- `summary.implemented_percent`: implemented operations divided by total
  operations.
- `summary.families[]`: family-level totals, implemented counts, percentage,
  priority, target level, and next milestone.

Family coverage is the operational unit for long-running work. A family moves
up only when the route claim, fixtures, tests, docs, and scorecard evidence are
merged together.

## Fill Loop

Each compatibility PR should follow this loop:

1. Pick one family row from `summary.families[]`, prioritizing `P0`, then
   adoption-critical `P1` gaps.
2. Add or raise one small endpoint cluster, not an entire Stripe family.
3. Add tests proving the new level: validation for `L1`, fixtures for `L2`,
   storage behavior for `L3`, scenario evidence for `L4`, webhook delivery for
   `L5`, and SDK/adoption smoke for `L6`.
4. Update `currentStripeRouteCoverage()` only after the test evidence exists.
5. Regenerate the inventory and include the before/after family delta in the PR
   description.
6. Review, fix, and merge before opening the next compatibility chunk.

## Current Priority Queue

| Priority | Family | First measurable chunk |
| --- | --- | --- |
| P0 | webhooks | Thin/snapshot event shape audit, connected-account routing evidence, replay diagnostics. |
| P0 | billing | Renewal, trial, cancellation, retry, proration preview, and subscription schedule scenarios. |
| P0 | checkout | Close subscription checkout optional-param gaps and SDK smoke. |
| P0 | billing_portal | Portal configuration fixtures, payment-method update, cancellation reason coverage. |
| P1 | connect | `GET /v1/accounts/{id}`, account list fixture, `Stripe-Account` tracing, connected webhook smoke. |
| P1 | payments | PaymentIntent and SetupIntent create/confirm/capture/cancel state machines. |
| P1 | payment_history | Charges, refunds, balance transactions, disputes, and debug-bundle reconciliation. |
| P1 | catalog/customers | OpenAPI-backed validation and low-state fixture response breadth. |

Connect is tracked as `P1` because platform-style SaaS integrations commonly
need account routing and connected-account webhooks before broad auxiliary API
coverage matters.

## PR Gate

A compatibility PR is not complete until it can answer:

- Which generated family row changed?
- Which operations moved from `L0` or a lower level?
- Which tests prove the new level?
- What remains inventory-only?
- What behavior is explicitly out of scope?

