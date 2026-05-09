# Implementation Plan

## Status

Draft

## Stack

- Backend: Go
- Frontend: React + TypeScript
- Storage: SQLite local default, in-memory test option
- Packaging: Docker

## Architecture Direction

Build a single Go server that serves:

- Stripe-compatible API subset
- Billtap dashboard API
- hosted checkout UI assets
- hosted portal UI assets
- developer dashboard assets

The server includes an embedded worker for webhook delivery and scenario clock events.

## Module Plan

### Go modules

```text
cmd/billtap
internal/api
internal/billing
internal/checkout
internal/portal
internal/webhooks
internal/scenarios
internal/assertions
internal/storage
internal/dashboard
internal/config
```

### React apps

```text
web/checkout
web/portal
web/dashboard
web/shared
```

## Milestones and Gates

### M0: Product definition

Gate:

- G0 Spec Readiness

### M1: Runtime contract

- Go server
- config
- SQLite
- migrations
- health
- base dashboard shell

Gate:

- G1 Runtime Contract

### M2: Checkout MVP

- customer/product/price APIs
- checkout session API
- checkout UI
- subscription/invoice/payment intent creation

Gate:

- G2 Checkout MVP

### M3: Webhook lab

- endpoint config
- signed delivery
- attempts
- retry/duplicate/delay/out-of-order/replay

Gate:

- G3 Webhook Reliability

### M4: Dashboard

- object detail
- timeline
- webhook detail
- debug bundle

Gate:

- G4 Debuggability

### M5: Scenario runner

- YAML DSL
- local clock
- app callback assertions
- reports

Gate:

- G5 CI Contract

### M6: Portal

- portal UI
- plan change
- seat change
- cancellation
- payment method update

Gate:

- G6 Portal Coverage

### M7: SaaS service profile

- workspace model
- tenant payment rail policy
- SaaS plan catalog
- seats and members
- export entitlement
- extra export payment and provision
- payment history and customer portal evidence
- back-office support actions
- platform/connect webhook evidence
- SaaS scenario pack

Gate:

- G7 SaaS Adoption Contract

### M8: Production boundaries

- masking
- audit
- retention
- relay boundary docs and tests

Gate:

- G8 Production Boundary
- G9 Release Candidate

## Key Risks

| Risk | Mitigation |
| --- | --- |
| Scope becomes full Stripe clone | keep useful subset and document limitations |
| UI takes too much time | split checkout, portal, dashboard into agent lanes |
| Webhook semantics are wrong | fixture-backed contract tests |
| Real payment confusion | clear copy and production boundaries |
| Scenario DSL becomes too complex | start with linear workflows and local clock |
| SaaS profile becomes too app-specific | isolate under profile modules and keep generic billing APIs independent |
