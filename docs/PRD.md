# Product Requirements Document

## Product

Billtap

## Date

2026-05-08

## Stack

- Backend: Go
- Frontend: React + TypeScript
- Storage: SQLite local default, in-memory test option

## Background

Subscription billing is hard to test because it mixes API calls, hosted checkout, asynchronous webhooks, time-based state transitions, retries, and app-specific entitlement logic. Existing mock tools cover pieces of this problem but do not give teams a complete local and CI billing lab.

Billtap provides a full-stack sandbox that simulates Stripe-style billing flows and makes the full integration timeline debuggable.

For SaaS adoption, Billtap must go beyond generic customer/subscription objects. It must model workspace plans, additional seats, export limits, extra export payments, payment history, back-office operations, and webhook evidence well enough that application server and frontend billing flows can be tested locally and in CI.

## Goals

- Provide local checkout and portal UI for subscription flows.
- Provide a useful Stripe-compatible API subset.
- Provide deterministic webhook delivery and failure simulation.
- Provide a scenario runner for CI.
- Provide app-side assertion hooks.
- Provide a dashboard that explains the billing lifecycle.
- Provide a SaaS service profile for workspace-centered billing, seats, export quota, and support debugging.

## Non-Goals

- Process real card payments.
- Become a production payment processor.
- Implement every Stripe API.
- Replace Stripe Billing in production.
- Store real card numbers or production secrets.

## Personas

### SaaS backend engineer

Needs to verify webhook handlers, idempotency, subscription state, invoice state, and entitlement updates.

### Frontend engineer

Needs a checkout and portal UI to test user-facing billing flows without connecting to real Stripe testmode every time.

### QA engineer

Needs deterministic billing scenarios for upgrade, downgrade, cancellation, payment failure, retry, and quota purchase.

### Support engineer

Needs a timeline that explains why a customer has a specific billing or entitlement state.

### SaaS platform engineer

Needs to reproduce workspace subscription, seat, export, and webhook problems without depending on live Stripe testmode state.

## MVP Requirements

### API

- Customers
- Products
- Prices
- Checkout Sessions
- Subscriptions
- Invoices
- PaymentIntents
- Webhook Endpoints
- Events
- Workspaces
- Workspace members
- Entitlement snapshots
- Export sessions
- Extra export payments
- Support operations

### UI

- Hosted checkout page
- Hosted billing portal page
- Developer dashboard
- Scenario run view
- Webhook delivery detail

### Scenario Runner

- YAML scenario format
- Local clock
- Expected webhook sequence
- Expected app callback response
- Optional app state assertion adapter
- SaaS profile actions for workspace subscription, seats, members, export, payment history, and webhook evidence

### Webhooks

- HMAC signature generation
- Delivery attempts
- Retry policy
- Duplicate delivery
- Delay
- Out-of-order delivery
- Replay

## Success Metrics

- A developer can run Billtap and complete a checkout flow in under 10 minutes.
- A CI scenario can catch a missing entitlement update.
- A webhook signature failure can be diagnosed from dashboard evidence.
- A subscription payment failure and retry flow can be reproduced deterministically.
- A debug bundle can explain one customer's billing timeline.
- A SaaS workspace debug bundle can explain subscription status, plan tier, seats, export quota, extra export payment, payment history, webhook delivery, and app assertion evidence.

## SaaS Profile Requirements

Billtap must provide a `saas` profile with fixture-backed semantics for:

- workspace ownership and active workspace state
- tenant payment rails: card, out-of-band, and connect
- plan tiers: free, standard, premium
- payment cycles: monthly, yearly, custom
- subscription statuses: none, processing, completed, failed, canceled
- additional seat estimate and purchase
- workspace member invite and deletion
- export summary and design-case export usage
- export session creation and file listing
- extra export preview, payment intent creation, and provision
- customer portal URL generation with masked signed URL evidence
- payment history from invoice, payment intent, and refund events
- platform and Connect webhook event handling
- back-office plan, seat, period, export limit, and refund actions
- expected observability events for billing and export scenarios

SaaS profile is successful when a developer can reproduce the common support questions:

- Why is this workspace still free, processing, failed, or canceled?
- Why can this user not invite a member?
- Why is export blocked?
- Did extra export payment succeed but provision fail?
- Did a duplicate webhook double-grant entitlement?
- Did Connect identity resolution fail?
- Which event should support paste into a Jira or Datadog investigation?
