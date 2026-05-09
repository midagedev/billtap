# Final Goal

## North Star

Billtap should become the fastest way for a SaaS team to build, debug, and regression-test subscription billing flows without depending on live Stripe testmode state for every local or CI scenario.

## Final Product Statement

Billtap is a full-stack billing sandbox with Stripe-compatible APIs, hosted checkout and portal UIs, scenario-driven subscription lifecycle simulation, webhook delivery debugging, and app-side assertion tooling.

## End State Capabilities

### Local Development

- Start with one Docker command.
- Create products, prices, customers, checkout sessions, subscriptions, invoices, payment intents, and webhook endpoints.
- Open a hosted sandbox checkout page.
- Open a hosted sandbox billing portal page.
- Choose payment outcomes from UI: success, failure, requires action, canceled, retry later.
- Inspect all API calls, objects, state transitions, webhooks, delivery attempts, and handler responses.

### CI

- Run billing scenarios from YAML.
- Advance a local billing clock.
- Emit deterministic webhook sequences.
- Assert app callbacks and optional app-side state.
- Fail with JSON and Markdown reports.

### Developer Dashboard

- Show customer and subscription timelines.
- Show invoice and payment intent state machines.
- Show webhook delivery attempts, signatures, retries, duplicate deliveries, delays, and out-of-order events.
- Show idempotency behavior.
- Show app assertion results.
- Export a debug bundle for one customer, subscription, checkout session, invoice, or webhook.
- Export a SaaS workspace debug bundle that includes owner, plan, subscription status, seats, members, export quota, extra export payments, payment history, webhook evidence, and app assertion results.

### Staging and Controlled Integration

- Relay or replay testmode webhook events.
- Compare expected scenario state with app behavior.
- Debug webhook signature and idempotency handling.
- Avoid storing secrets or real card data.

## Target User Outcome

An engineer should be able to answer these questions in minutes:

- Which billing object changed?
- Which webhook was emitted?
- Was it delivered?
- Did the app return success?
- Did the app update the expected entitlement, plan, seat, quota, or workspace state?
- Did a SaaS workspace have the expected plan tier, subscription status, seat capacity, member state, and export entitlement?
- Would the app survive duplicate, delayed, or out-of-order webhook delivery?
- What exact scenario reproduces the issue?

## Strategic Constraints

- Billtap must not become a real PSP.
- Billtap must not require app code to depend on a custom SDK.
- Billtap should work with Stripe-like API and webhook semantics.
- Full Stripe compatibility is not required; documented useful subset is required.
- Webhook reliability and app-side assertions are the differentiators.

## Final Release Definition

Billtap reaches the final goal when all of these are true:

- Stripe-compatible API subset has fixture-backed contract tests.
- Checkout and portal UIs are usable for realistic subscription workflows.
- Scenario runner can exercise upgrade, downgrade, cancellation, payment failure, retry, seat change, and quota purchase flows.
- SaaS profile can exercise workspace subscription, seat purchase, member invite/delete, export summary, export session, extra export payment, payment history, back-office refund, and Connect webhook flows.
- Webhook delivery supports signature, retry, duplicate, delay, out-of-order, and replay.
- Dashboard can explain the full billing timeline.
- CI mode can validate app behavior and produce actionable reports.
- Production boundaries are documented and enforced.
