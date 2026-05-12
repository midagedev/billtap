# Simulation Capacity Backlog

Status date: 2026-05-12

This backlog turns production-regression learnings into public, product-neutral
Billtap capability work. It is intentionally broader than the current
implementation. A capability should move from backlog to compatibility claim
only after it has tests, fixture examples, documented boundaries, and, where
useful, SDK or app-style smoke evidence.

Do not commit private customer names, service names, internal issue IDs, or
incident evidence here. Adoption repositories can map private incidents to
these public capability rows in their own private docs.

## Priority Model

| Tier | Meaning | Merge rule |
| --- | --- | --- |
| P0 | Blocks high-value billing regression automation or closes a known class of payment/subscription production bugs. | Implement with API or scenario tests plus docs. |
| P1 | Expands near-term billing matrix depth for subscriptions, invoices, payment methods, webhooks, or history. | Implement as one bounded endpoint/scenario cluster per PR. |
| P2 | Useful for broader product coverage, tax/compliance, multi-region, or hosted UI parity. | Implement only with a clear adoption scenario. |
| P3 | Observability, performance, or developer-experience support. | Prefer additive diagnostics and generated evidence. |

## Current Baseline

Billtap currently has a stateful Stripe-like local subset with:

- customers, products, prices, price search, coupons, promotion codes, checkout
  sessions, subscriptions, subscription schedules, invoices, invoice preview,
  payment intents, setup intents, payment methods, refunds, credit notes,
  disputes, cash balance, test clocks, Connect evidence, webhook endpoints, and
  events
- fixture apply/snapshot/assert APIs with stable fixture IDs and fixture refs
- local checkout and portal UI flows
- webhook retry, duplicate, delay, out-of-order, replay, historical replay, and
  delivery attempt evidence
- diagnostic bundles, request traces, timeline evidence, and dashboard views
- OpenAPI inventory baseline: `144 / 587` implemented operations, `24.5%`

## P0 Regression-Paired Backlog

These are the highest-value next chunks because they make app-side billing
regression specs easier to close with deterministic Billtap evidence.

| Capability | Current state | Next implementation shape | Verification |
| --- | --- | --- | --- |
| Trial-used customer fixture | Trialing subscriptions, canceled subscriptions, stable fixture IDs, and test-clock trial activation exist. A dedicated "trial history only" customer fixture pattern is not documented or asserted. | Add fixture examples and assertions for a customer with canceled trial history but no active subscription. Preserve customer metadata and subscription history lookup behavior. | API tests for seeded history, scenario using fresh checkout after trial history, docs fixture snippet. |
| Trial to cancel history | Immediate cancellation and period-end cancellation exist; trial cancellation history needs a direct scenario contract. | Add scenario action/fixture pattern that creates trialing subscription, cancels it before activation, and leaves retrievable canceled history. | Scenario report and webhook sequence assertions. |
| Proration invariant matrix | Subscription-update invoice preview calculates bounded proration and Stripe SDK-friendly invoice fields. | Add a matrix of small positive, large positive, negative/downgrade, zero, and full-credit preview cases. | API tests for exact `amount_due`, `subtotal`, line periods, and discount/tax-neutral totals. |
| Renewal failure depth | Renewal failure can move subscriptions to `past_due` with invoice/payment-intent evidence. | Add configurable retry/dunning phases: first failure, retry failure, unpaid, grace window, and final cancellation. | Test-clock scenario reports and event ordering assertions. |
| Customer subscription history filters | Top-level subscription listing exists; customer nested subscription routes remain inventory-only. | Add `/v1/customers/{customer}/subscriptions` list with `status` filter coverage for `trialing`, `active`, `past_due`, `canceled`, `incomplete`, `unpaid`, and `all`. | Stripe-shaped list response tests and fixture history scenarios. |

## Capability Matrix

### 1. Subscription Lifecycle

| Capability | Current state | Priority |
| --- | --- | --- |
| trial to active | Implemented with test clocks and fixture dates. | Done |
| trial to cancel history | Partial through cancellation APIs; needs dedicated scenario pattern. | P0 |
| past_due and unpaid transitions | Partial renewal failure; needs retry/dunning phases. | P0 |
| pause and resume API | Not modeled. | P1 |
| subscription schedule plan change | One due phase can replace items; full scheduled downgrade/proration remains bounded. | P1 |
| proration behavior `always_invoice`, `create_prorations`, `none` | Preview supports the three values; invoice creation side effects remain bounded. | P0 |
| cancel at period end vs immediate cancel | Implemented through API/portal/test-clock paths. | Done |
| grace period simulation | Not modeled as a separate state. | P1 |

### 2. Customer History And Metadata

| Capability | Current state | Priority |
| --- | --- | --- |
| customer subscription list status filter | Nested route inventory exists; implementation needed. | P0 |
| arbitrary customer metadata seed | Implemented through fixtures and customer APIs. | Done |
| customer search by email/metadata | Inventory-only. | P1 |
| customer delete and recreate history policy | Not modeled. | P1 |
| cross-merchant customer migration | Connect evidence exists; migration policy not modeled. | P2 |
| tax IDs, shipping, preferred locale | Partially represented in validation inventory; stateful coverage needed. | P2 |

### 3. Invoice Correctness

| Capability | Current state | Priority |
| --- | --- | --- |
| proration exactness matrix | Bounded proration exists; matrix coverage needed. | P0 |
| invoice discounts and coupons | Bounded single-discount subset implemented. | Done |
| automatic tax and regional rates | Response fields exist; tax calculation not modeled. | P2 |
| invoice line period and quantity diversity | Basic proration lines exist; broader cases needed. | P1 |
| amount due/paid/remaining consistency | Implemented for checkout, retry, renewal, and preview paths. | Done |
| zero-dollar invoice | Needs explicit full-credit scenario. | P1 |
| negative balance and ending balance | Not modeled beyond default fields. | P1 |
| hosted invoice URL and receipt | Placeholder fields exist; hosted receipt UI not modeled. | P2 |
| invoice payment settings diversity | Default response shape exists; behavior breadth needed. | P2 |

### 4. Payment Intent And Payment Method

| Capability | Current state | Priority |
| --- | --- | --- |
| 3D Secure / requires_action | Direct PaymentIntent and SetupIntent require-action simulation exists. | Done |
| saved card vs new card | Customer payment-method fixtures and portal save exist; create/attach/detach breadth remains. | P1 |
| wallet token simulation | Not modeled. | P2 |
| decline reason breadth | Common aliases exist; complete card-decline catalog remains. | P1 |
| provider outage mode | Not modeled. | P2 |
| mandate flows | Not modeled beyond shape fields. | P2 |

### 5. Webhook Reliability

| Capability | Current state | Priority |
| --- | --- | --- |
| retry after delivery failure | Implemented through replay and simulated failure evidence. | Done |
| duplicate delivery | Implemented. | Done |
| out-of-order delivery | Implemented. | Done |
| signature timestamp tolerance cases | Signature evidence exists; explicit tolerance matrix needed. | P1 |
| API version variation | Configurable webhook API version exists; compatibility matrix needed. | P1 |
| dead-letter simulation | Not modeled. | P2 |
| historical replay automation | Explicit historical replay exists; automatic policy remains bounded. | P1 |

### 6. Test Clocks

| Capability | Current state | Priority |
| --- | --- | --- |
| multi-clock customer simulation | Implemented through customer/subscription test-clock metadata. | Done |
| past period replay | Partial through fixture dates and clock advance. | P1 |
| mid-cycle proration on clock advance | Not modeled as a clock side effect. | P1 |
| trial expiry auto advance | Implemented. | Done |

### 7. Connect And Multi-Account

| Capability | Current state | Priority |
| --- | --- | --- |
| connected accounts | Implemented for local evidence and request traces. | Done |
| transfers, payouts, application fees | Implemented as local evidence. | Done |
| application fee reversal | Implemented for fee refunds/reversals evidence. | Done |
| platform vs connected customer separation | Request trace and Connect evidence exist; deeper isolation scenarios needed. | P1 |

### 8. Coupon And Promotion Code

| Capability | Current state | Priority |
| --- | --- | --- |
| duration once/repeating/forever | Shape exists; redemption lifecycle not modeled. | P1 |
| coupon stacking | Not modeled; single effective discount only. | P2 |
| promotion redemption count and max redemptions | Not modeled. | P1 |
| coupon currency restrictions | Amount-off currency guard exists; broader restrictions needed. | P1 |
| coupon expiration | Not modeled. | P1 |

### 9. Refund And Credit Note

| Capability | Current state | Priority |
| --- | --- | --- |
| partial refund and remaining amount | Implemented for local refund evidence; richer remaining-balance assertions needed. | P1 |
| refund reason classification | Implemented for common local reasons; UI/payment-history scenarios needed. | P1 |
| credit note simulation | Implemented for create/retrieve/void local evidence. | Done |
| credit note vs refund branch | Needs scenario matrix. | P1 |
| bank-transfer refund | Pending/settlement simulation exists for refunds; bank-transfer-specific behavior remains. | P2 |
| refund webhooks | Implemented for local refund paths; expand event matrix as needed. | P1 |

### 10. Charge And Dispute

| Capability | Current state | Priority |
| --- | --- | --- |
| dispute creation | Implemented as local evidence. | Done |
| evidence submission | Basic update evidence exists; structured evidence fields need expansion. | P1 |
| dispute outcome won/lost/warning closed | Not fully modeled. | P1 |
| dispute fees and adjustments | Not modeled. | P2 |

### 11. Tax, VAT, And Compliance

| Capability | Current state | Priority |
| --- | --- | --- |
| tax ID collection | Inventory-visible; local stateful flow needed. | P2 |
| automatic tax option | Invoice response shape exists; calculation not modeled. | P2 |
| regional tax rates | Not modeled. | P2 |
| invoice receipt PDF/HTML | Not modeled beyond placeholder fields. | P2 |
| VAT refund cases | Not modeled. | P2 |

### 12. SetupIntent And Card Registration

| Capability | Current state | Priority |
| --- | --- | --- |
| SetupIntent simulation | Implemented for create/list/retrieve/confirm/cancel. | Done |
| card add/delete/default switch | Portal save exists; PaymentMethod attach/detach/update breadth needed. | P1 |
| card expiry notification | Not modeled. | P2 |
| payment method attach/detach | Inventory-visible; stateful breadth needed. | P1 |
| off-session setup for renewal | Partial through default payment methods; mandate/off-session detail needed. | P2 |

### 13. Multi-Region

| Capability | Current state | Priority |
| --- | --- | --- |
| region-specific webhook routing | Connect/webhook evidence exists; explicit region router simulation needed. | P2 |
| cross-region customer migration | Not modeled. | P2 |
| region-specific currency | Multi-currency price fields exist; scenario matrix needed. | P2 |
| region-specific tax rate | Not modeled. | P2 |
| wrong-region forwarding | Not modeled. | P3 |

### 14. Bank Transfer And Payment Diversity

| Capability | Current state | Priority |
| --- | --- | --- |
| ACH credit transfer | Cash-balance funding and bank-transfer PaymentIntent processing exist in bounded form. | P1 |
| SEPA debit | Not modeled. | P2 |
| local virtual account rails | Not modeled. | P3 |
| invoice `payment_method_types` diversity | Shape exists; behavior breadth needed. | P2 |
| manual vs automatic confirmation | PaymentIntent confirmation modes need broader coverage. | P1 |

### 15. Subscription, Price, And Product Metadata

| Capability | Current state | Priority |
| --- | --- | --- |
| subscription metadata seed | Implemented through fixtures and APIs. | Done |
| price lookup key collision | Search exists; collision validation needed. | P1 |
| product statement descriptor | Shape/metadata coverage needed. | P2 |
| tax behavior inclusive/exclusive | Inventory-visible; invoice math not modeled. | P2 |
| tiered pricing | Not modeled. | P2 |

### 16. Replay And Simulation Tools

| Capability | Current state | Priority |
| --- | --- | --- |
| specific event replay | Implemented. | Done |
| bulk seed/dump | Fixture apply/snapshot exists; dump/export format can be expanded. | P1 |
| diagnostic snapshots | Implemented; JSON schema standardization needed. | P3 |
| time-travel rewind | Not modeled. | P3 |
| cross-test isolation | Fixture namespace/ref metadata exists; runner isolation policy can improve. | P1 |

### 17. Edge Cases

| Capability | Current state | Priority |
| --- | --- | --- |
| partial payment and remaining balance | Limited invoice retry evidence; partial payment not modeled. | P1 |
| subscription quantity zero invalid | Validation exists for supported quantity paths. | Done |
| promotion plus proration | Bounded single-discount proration preview exists. | Done |
| zero-dollar invoice | Needs explicit scenario. | P1 |
| negative balance | Not modeled. | P1 |
| concurrent checkout and cancel | Not modeled. | P2 |
| webhook before sync API race | Can be simulated manually with out-of-order replay; scenario helper needed. | P1 |

### 18. Hosted Checkout And Portal UI

| Capability | Current state | Priority |
| --- | --- | --- |
| hosted checkout UI selector guard | Local UI exists; Stripe selector parity is not claimed. | P2 |
| billing portal manage/cancel/payment method | Local portal actions exist. | Done |
| invoice download UI | Not modeled. | P2 |
| customer portal plan change | Local portal plan change exists; deeper scenario coverage needed. | P1 |

### 19. Rate Limit And API Error Simulation

| Capability | Current state | Priority |
| --- | --- | --- |
| rate-limited 429 | Not modeled. | P1 |
| API connection/auth/API errors | Stripe-shaped errors exist; injectable outage/error modes needed. | P1 |
| idempotency key consistency | Implemented for same-process POST replay and mismatch. | Done |
| partial outage mode | Not modeled. | P2 |

### 20. Observability And Diagnostics

| Capability | Current state | Priority |
| --- | --- | --- |
| request log and replay timeline | Request traces and timelines exist. | Done |
| webhook delivery dashboard | Dashboard and delivery attempts exist; retry summary UX can improve. | P2 |
| diagnostic export JSON schema | Export exists; schema standardization needed. | P3 |
| fixture seed assertions | Fixture assertion API exists; more domain assertions needed. | P1 |
| latency/performance metrics | Not modeled. | P3 |

## Execution Rules

1. Keep broad capability tracking in this file, but implement behavior in
   bounded PRs that each name the compatibility family and verification.
2. Do not raise `internal/stripecompat` claim levels until tests and docs prove
   the new behavior.
3. Any change affecting billing state, webhook order, signature handling,
   idempotency, or request validation must include tests.
4. Any simulation that deliberately diverges from Stripe must be explicit in
   `docs/COMPATIBILITY.md`.
5. Avoid private adoption names, internal issue keys, and company-specific
   terms in public docs, code, fixtures, and PR descriptions.
