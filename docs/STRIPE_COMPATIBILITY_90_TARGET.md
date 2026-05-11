# Stripe Compatibility 90% Target

Status date: 2026-05-11

Billtap's long-running Stripe API compatibility target is measurable coverage
of at least 90% of the public Stripe OpenAPI operation inventory, without
claiming that every operation has deep payment-processing behavior.

## Target Definition

The 90% target is based on generated `stripe-api-inventory.json`:

- **Overall target:** `summary.implemented_percent >= 90.0`.
- **Current baseline:** `98 / 587` operations, `16.7%`, using Stripe OpenAPI
  `2026-04-22.dahlia` from `stripe/openapi` master on 2026-05-11.
- **Minimum target count:** `529 / 587` operations at `L1` or higher.
- **Remaining inventory-only budget:** at most `58 / 587` operations at `L0`.

This is a broad compatibility target. The target does not mean that 90% of
Stripe has local state machines, real financial behavior, or complete Stripe
parity.

## Level Targets

| Scope | Required target |
| --- | --- |
| Overall OpenAPI operations | At least 90% at `L1+`. |
| P0 billing lab families | At least 85% at `L3+`; critical flows at `L4-L6`. |
| P1 adoption families | At least 75% at `L2+`; adoption-critical flows at `L3-L6`. |
| P3 auxiliary families | At least 90% at `L1-L2`; fixture/schema smoke only. |
| Explicitly unsafe financial domains | May remain below 90% when documented as out of scope. |

## What Counts

An operation counts toward the 90% total only when it has a tested Billtap
runtime claim:

- `L1`: route and request validation with Stripe-shaped error envelopes.
- `L2`: deterministic Stripe-shaped fixture response and list envelope basics.
- `L3`: stateful local mutation and retrieval.
- `L4`: scenario runner and timeline/debug evidence.
- `L5`: webhook emission and delivery evidence.
- `L6`: official SDK or adoption smoke.

Inventory-only documentation does not count.

The generated OpenAPI validation catalog is separate from implementation
coverage. `summary.schema_validated_operations` can rise independently for a
Stripe OpenAPI input with parameter/request-body schemas while
`summary.implemented_operations` remains unchanged; an operation only counts
toward the 90% target after it has an explicit tested claim at `L1+`.

## Baseline By Family

Latest measured baseline from `stripe/openapi` master on 2026-05-11:

| Priority | Family | Total | Implemented | Coverage | 90% target count | First target |
| --- | --- | ---: | ---: | ---: | ---: | --- |
| P0 | webhooks | 7 | 7 | 100.0% | 7 | Expand connected-account routing, thin event fixtures, and replay evidence. |
| P0 | checkout | 6 | 3 | 50.0% | 6 | Close checkout route gaps and SDK smoke. |
| P0 | billing | 39 | 11 | 28.2% | 36 | Add remaining invoice, schedule, coupon, discount, and deeper credit-note/test-clock behavior. |
| P0 | billing_portal | 5 | 1 | 20.0% | 5 | Add portal configurations and session retrieval fixtures. |
| P1 | catalog | 54 | 9 | 16.7% | 49 | Add low-state product/price/coupon/promotion/tax fixtures and validation. |
| P1 | customers | 31 | 4 | 12.9% | 28 | Add search, sources, tax ids, cash balance, and validation fixtures. |
| P1 | payments | 41 | 12 | 29.3% | 37 | Add PaymentMethod lifecycle breadth and remaining PaymentIntent/SetupIntent adjunct routes. |
| P1 | connect | 53 | 41 | 77.4% | 48 | Add account self/delete and people/person fixtures to close the remaining Connect inventory routes. |
| P1 | payment_history | 30 | 6 | 20.0% | 27 | Add charges, deeper refunds, balance transactions, disputes, credit history. |
| P3 | auxiliary | 321 | 4 | 1.2% | 289 | Add generic L1/L2 schema and fixture smoke for safe low-state endpoints. |

## PR Chunk Plan

Each chunk follows PR -> review -> fix -> merge.

| Wave | Chunk | Target delta | Output |
| --- | --- | ---: | --- |
| T0 | 90% target and scoreboard gate | 0 ops | This target doc, thresholds, and PR queue. |
| T1 | Runtime claim registry | 0 ops | Add `internal/stripecompat` claim registry and route matching shared by runtime and inventory. |
| T2 | Known-route unsupported fallback | 0 ops | Known OpenAPI routes return Stripe-shaped unsupported errors and traces without counting as implemented. |
| T3 | Auxiliary L1 route validation | +155 ops | Generic OpenAPI route registry, auth/error/list envelope validation, no fake financial state. |
| T4 | Auxiliary L2 fixture responses | +155 ops | Deterministic fixture responses for safe low-state retrieve/list endpoints. |
| T5 | Catalog/customers L1-L2 breadth | +72 ops | Products, prices, coupons, promotion codes, tax rates, customer adjunct resources. |
| T6 | Connect platform evidence | +48 ops | Accounts, account links/sessions, capabilities, external accounts, application fees, transfers, payouts, `Stripe-Account` traces. |
| T7 | Payments and setup breadth | +29 ops | PaymentMethod create/retrieve/list/update/cancel fixture/state coverage and remaining PaymentIntent/SetupIntent adjunct routes. |
| T8 | Payment history L2-L3 | +30 ops | Charges, refunds, balance transactions, disputes, debug bundle evidence. |
| T9 | Billing lifecycle depth | +37 ops | Schedules, invoice items, trials, renewals, coupons, discounts, credit notes, test clocks. |
| T10 | SDK/adoption matrix | 0 ops | Node/Go/Java/Python/Ruby smoke and adoption-style reports that promote existing operations to `L6`. |

The deltas are planning targets, not claims. Each PR records the actual
before/after `summary.implemented_percent` and changed family rows.

T1 and T2 intentionally do not increase coverage. They prevent long-term
coverage inflation from becoming a hardcoded list in
`currentStripeRouteCoverage()` and ensure inventory-only routes are visible at
runtime without being counted as implemented.

T135 adds the broad OpenAPI validation catalog and schema-validation fallback
for known routes, but it does not by itself increase
`summary.implemented_operations`. T3-T9 still need to opt specific operations
into tested `L1+` claims.

T139 adds the first Connect smoke slice: account create/list/retrieve/update,
account links, account sessions, and `Stripe-Account` request trace evidence.
This raises the generated inventory from `57 / 587` (`9.7%`) to `63 / 587`
(`10.7%`), with Connect moving from `0 / 53` to `6 / 53`.

T140 expands Connect platform evidence across account capabilities, external
accounts, transfers/reversals, payouts, application fees/refunds, and local
Connect webhook evidence. This raises the generated inventory from `63 / 587`
(`10.7%`) to `98 / 587` (`16.7%`), with Connect moving from `6 / 53` to
`41 / 53`.

T10 also does not increase `summary.implemented_operations` by itself. It
raises confidence and levels for already counted operations; new operation
coverage must come from T3-T9. The planned T3-T9 delta is intentionally larger
than the `+431` operations needed to move the current `98 / 587` baseline to
the `529 / 587` target.

## Derived Gate Checks

The overall `summary.implemented_percent` check is necessary but not sufficient.
G14 also uses `summary.families[].by_level` to calculate family depth:

- P0 depth: for each P0 family, `(L3 + L4 + L5 + L6) / total >= 85%`.
- P1 depth: for each P1 family, `(L2 + L3 + L4 + L5 + L6) / total >= 75%`.
- P3 breadth: for auxiliary families, `(L1 + L2) / total >= 90%`, unless the
  family is explicitly documented as unsafe or out of scope.

## Non-Goals

- Do not process real payments.
- Do not store real card data.
- Do not model regulated Stripe domains deeply unless there is a local testing
  use case and a production-boundary review.
- Do not mark an endpoint implemented solely because it appears in OpenAPI.
