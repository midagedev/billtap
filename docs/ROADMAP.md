# Roadmap

## Final Target

Billtap should become a full-stack billing lab for Stripe-style subscription products: API sandbox, hosted checkout, hosted portal, webhook debugger, scenario runner, CI validator, developer dashboard, and app-specific billing profiles.

Each phase maps to gates in `specs/000-product/gates.md`.

## Phase 0: Product Definition

- concept
- PRD
- architecture
- UI product definition
- final goal
- agent orchestration
- gates

Gate:

- G0 Spec Readiness

## Phase 1: Runtime and Billing Core

- Go backend skeleton
- config
- SQLite storage
- billing state model
- customers, products, prices
- health endpoints

Gate:

- G1 Runtime Contract

## Phase 2: Checkout and Subscription MVP

- checkout session API
- hosted checkout UI
- subscription create flow
- invoice and payment intent creation
- payment success and failure outcomes

Gate:

- G2 Checkout MVP

## Phase 3: Webhook Lab

- webhook endpoints
- signed event delivery
- delivery attempts
- retries
- duplicate, delay, out-of-order
- replay

Gate:

- G3 Webhook Reliability

## Phase 4: Developer Dashboard

- object list/detail
- timeline
- webhook delivery detail
- scenario run view
- debug bundle

Gate:

- G4 Debuggability

## Phase 5: Scenario Runner and CI

- scenario DSL
- local clock
- app callback assertions
- JSON and Markdown reports
- exit code policy

Gate:

- G5 CI Contract

## Phase 6: Billing Portal

- portal UI
- plan change
- seat change
- cancellation
- payment method update
- invoice history

Gate:

- G6 Portal Coverage

## Phase 7: SaaS Service Profile

- workspace model
- SaaS plan catalog
- tenant payment rail policy
- additional seats and members
- export entitlement
- extra export payment and provision
- payment history
- back-office support operations
- platform and Connect webhook evidence
- SaaS scenario pack

Gate:

- G7 SaaS Adoption Contract

## Phase 8: Staging Relay and Production Boundaries

- testmode webhook relay
- safe replay
- secret masking
- retention policy
- operational runbook

Gate:

- G8 Production Boundary

## Phase 9: Release Candidate

- Docker image
- sample app
- examples
- public docs
- release notes

Gate:

- G9 Release Candidate

## Phase 10+: Broader Stripe API Compatibility

The current public claim remains a documented, stateful Stripe-like subset for
local subscription billing tests. Broader compatibility work should proceed
through `STRIPE_API_COMPATIBILITY_ROADMAP.md`:

- OpenAPI inventory and drift tracking
- protocol-level compatibility for errors, pagination, expand, idempotency, and
  request traces
- OpenAPI-backed validation and fixture-shaped responses
- deeper subscription, invoice, payment-intent, setup-intent, refund, dispute,
  entitlement, metering, and Connect simulation
- regression-driven simulation capacity tracked in
  `SIMULATION_CAPACITY_BACKLOG.md`
- official Stripe SDK smoke lanes and optional `stripe-mock`/live-testmode
  oracle jobs

Gate:

- new compatibility claims are levelled, scorecard-backed, and documented
  before release
