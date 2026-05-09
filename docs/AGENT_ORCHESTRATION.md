# Agent Orchestration Plan

## Purpose

Billtap is large enough to benefit from parallel agent development. The plan splits the app into stable ownership lanes and gates.

## Required Reading

Every agent reads:

1. `.specify/memory/constitution.md`
2. `docs/FINAL_GOAL.md`
3. `specs/000-product/spec.md`
4. `specs/000-product/gates.md`
5. this document

## Agent Roles

| Agent | Ownership | Primary outputs |
| --- | --- | --- |
| Spec Lead | specs, gates, ADRs | final scope, contracts, acceptance criteria |
| Runtime Core Agent | Go server, config, storage | API skeleton, SQLite, health, migrations |
| Billing Engine Agent | billing state machine | customer, product, price, subscription, invoice, payment intent transitions |
| Checkout UI Agent | hosted checkout | checkout React app, payment outcome flow |
| Portal UI Agent | hosted portal | plan change, cancellation, payment method, invoice history |
| Dashboard Agent | developer dashboard | timeline, object detail, webhook detail, scenario views |
| Webhook Agent | delivery system | signature, retries, duplicate, delay, out-of-order, replay |
| Scenario Agent | scenario runner | YAML DSL, local clock, reports |
| Assertion Agent | app-side checks | callback assertions and adapters |
| SaaS Profile Agent | workspace billing profile | workspace, seats, export, payment history, back-office fixtures, profile assertions |
| Compatibility Agent | Stripe-like API contracts | fixtures, endpoint behavior, SDK compatibility checks |
| Security Agent | secrets and boundaries | masking, retention, audit, production boundary tests |
| Release Agent | packaging | Docker image, examples, release checklist |

## Parallel Waves

### Wave 0: Scope Lock

- Spec Lead finalizes MVP and gates.
- Runtime Core Agent proposes module boundaries.
- UI agents define screen map.
- Compatibility Agent defines endpoint subset.

Gate:

- G0 Spec Readiness

### Wave 1: Runtime and Domain Core

- Runtime Core Agent builds Go skeleton and storage.
- Billing Engine Agent builds state model and transitions.
- Compatibility Agent adds fixtures and contract tests.

Gate:

- G1 Runtime Contract

### Wave 2: Checkout MVP

- Checkout UI Agent builds hosted checkout.
- Billing Engine Agent wires checkout to subscription/invoice/payment intent.
- Webhook Agent emits initial checkout and invoice events.

Gate:

- G2 Checkout MVP

### Wave 3: Webhook Reliability

- Webhook Agent implements signature, retry, duplicate, delay, out-of-order, replay.
- Dashboard Agent builds webhook delivery detail.
- Scenario Agent adds webhook scenarios.

Gate:

- G3 Webhook Reliability

### Wave 4: Dashboard and Debugging

- Dashboard Agent builds timeline and debug bundle.
- Assertion Agent adds app callback evidence.
- Scenario Agent adds scenario run view.

Gate:

- G4 Debuggability

### Wave 5: CI and Scenario Runner

- Scenario Agent implements YAML runner and local clock.
- Assertion Agent implements app-state assertion hooks.
- Release Agent adds CI examples.

Gate:

- G5 CI Contract

### Wave 6: Portal

- Portal UI Agent builds subscription management.
- Billing Engine Agent implements portal transitions.
- Webhook Agent emits related events.

Gate:

- G6 Portal Coverage

### Wave 7: SaaS Profile

- SaaS Profile Agent implements workspace, seats, export, extra export payment, payment history, support actions, and profile assertions.
- Scenario Agent adds the SaaS priority scenario pack.
- Dashboard Agent adds workspace entitlement, export, extra export, and support evidence views.
- Webhook Agent adds platform/connect profile evidence.

Gate:

- G7 SaaS Adoption Contract

### Wave 8: Production Boundaries and Release

- Security Agent locks masking, retention, audit, and boundary tests.
- Release Agent builds Docker image and examples.
- Spec Lead verifies docs match behavior.

Gate:

- G8 Production Boundary
- G9 Release Candidate

## Handoff Format

```text
Summary
- What changed

Files changed
- path

Verification
- command or evidence

Open risks
- risk and next step

Gate status
- passed / failed / not applicable
```

## Integration Rules

- Agents own disjoint files where possible.
- Shared contracts are changed only by Spec Lead or lead integration.
- Compatibility claims require fixtures.
- UI changes require at least one browser verification once implementation exists.
- Webhook behavior changes require retry/idempotency tests.
- SaaS profile changes require at least one workspace subscription, seat, export, and webhook scenario.
- Do not proceed past a failed gate except for independent research tasks.
