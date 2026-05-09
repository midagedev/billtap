# Billtap Concept

## One-Line Description

Billtap is a full-stack billing sandbox for Stripe-style checkout, subscriptions, webhooks, and app entitlement testing.

## Positioning

Billtap is not just a mock server. It is a billing lab.

It combines:

- API mock
- hosted checkout UI
- billing portal UI
- stateful billing engine
- webhook reliability lab
- scenario runner
- developer dashboard
- app-side assertion tool

## Why Existing Tools Are Not Enough

Stripe's official tools are valuable:

- `stripe-mock` provides API-shaped responses.
- Stripe CLI helps trigger and listen to webhooks.
- Test Clocks help test time-dependent billing in Stripe testmode.

But teams still struggle to test the app-side outcome:

- Did the workspace subscription change?
- Did seat count update?
- Did quota increase?
- Did payment failure move the account into the correct grace state?
- Did retry, duplicate, or out-of-order webhook delivery break idempotency?

Billtap focuses on those outcomes.

For SaaS, those outcomes are workspace-level outcomes: plan tier, subscription status, seat capacity, member invitations, export quota, extra export payment provision, payment history, and webhook evidence.

## Core Product Loop

1. Define or create a billing scenario.
2. Start checkout or portal UI.
3. Choose customer actions and payment outcomes.
4. Watch billing state and webhooks in the dashboard.
5. Assert the app-side result.
6. Save the scenario into CI.

## SaaS Service Loop

1. Pick a workspace scenario such as plan upgrade, seat purchase, or extra export.
2. Run Billtap with the `saas` profile.
3. Drive the frontend or application server against Billtap endpoints.
4. Watch workspace plan, seats, export quota, payment history, and webhook events in one timeline.
5. Assert the app-side state that support and QA actually care about.
6. Save the scenario as a regression test.

## Product Philosophy

### State matters

Billing is a state machine. Billtap must show object transitions, not only HTTP responses.

### Webhooks are the real integration

Most subscription bugs happen around async webhook delivery. Billtap should make webhook behavior explicit and controllable.

### UI matters

Developers and QA need to experience checkout and portal flows, not only call APIs.

### Determinism matters

CI scenarios must produce repeatable state transitions and reports.
