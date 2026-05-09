# Product Specification: Billtap

## Status

Draft

## Summary

Billtap is a full-stack billing sandbox for Stripe-style subscription apps. It provides a stateful API subset, hosted checkout UI, hosted billing portal UI, webhook delivery lab, scenario runner, dashboard, and app-side assertion tooling.

## Problem

Billing integrations fail at the boundary between payment provider state and application state. Mock HTTP responses are insufficient because real billing behavior includes hosted UI, asynchronous webhooks, retries, duplicate delivery, time-based subscription transitions, idempotency, and app-specific entitlement updates.

## Users

- backend engineer
- frontend engineer
- QA engineer
- support engineer
- platform engineer

## Core Scenarios

### Scenario 1: Local checkout success

A developer creates a checkout session, opens Billtap hosted checkout, selects success, and sees the subscription, invoice, payment intent, webhook events, delivery attempts, and app response in one timeline.

Acceptance criteria:

- checkout UI is usable
- subscription is created
- invoice is paid
- payment intent succeeds
- expected webhooks are emitted
- dashboard shows the timeline

### Scenario 2: Payment failure and retry

A QA engineer runs a scenario where the first invoice payment fails, the app receives failure webhook, local clock advances, retry succeeds, and entitlement updates.

Acceptance criteria:

- failure webhook emitted
- retry schedule visible
- retry can be advanced by local clock
- success webhook emitted after retry
- app assertion can verify final state

### Scenario 3: Webhook reliability

A backend engineer sends duplicate, delayed, and out-of-order webhooks to verify idempotent handling.

Acceptance criteria:

- delivery pattern is configurable
- attempts are signed
- app response is recorded
- dashboard identifies duplicate and out-of-order events

### Scenario 4: Billing portal

A frontend engineer opens the portal, changes plan, cancels at period end, and updates payment method.

Acceptance criteria:

- portal UI supports core actions
- subscription state changes
- correct webhooks emitted
- timeline explains changes

### Scenario 5: CI scenario

CI runs a YAML scenario and fails when the app does not update expected entitlement.

Acceptance criteria:

- scenario is deterministic
- report includes failing step
- report includes app response and expected assertion
- exit code is non-zero on failure

### Scenario 6: SaaS workspace entitlement

A platform engineer runs a SaaS profile scenario that upgrades a workspace plan, purchases seats, invites members, consumes export quota, buys extra export, and verifies app-side workspace state.

Acceptance criteria:

- workspace subscription status is visible
- plan tier and cycle are visible
- seat capacity and used seats are visible
- export summary and usage are visible
- extra export payment and provision status are visible
- app assertions can validate the workspace state

### Scenario 7: SaaS support investigation

A support engineer opens one workspace timeline and sees subscription, seat, export, payment history, webhook, Connect identity, and workspace log evidence.

Acceptance criteria:

- support bundle includes workspace, subscription, seats, export, payment, webhook, and assertion evidence
- duplicate webhooks are identified
- unresolved Connect webhook events are recorded
- failure reasons are copyable

## Functional Requirements

### API

- FR-001: Create, retrieve, list, and update customers.
- FR-002: Create, retrieve, and list products.
- FR-003: Create, retrieve, and list prices.
- FR-004: Create and retrieve checkout sessions.
- FR-005: Create, retrieve, update, cancel, and list subscriptions.
- FR-006: Create, retrieve, finalize, pay, void, and list invoices.
- FR-007: Create, retrieve, confirm, fail, and list payment intents.
- FR-008: Create, retrieve, update, delete, and list webhook endpoints.
- FR-009: Create, retrieve, and list events.

### Hosted UI

- FR-020: Hosted checkout page for subscription checkout.
- FR-021: Hosted checkout supports success, failure, requires-action, pending, and cancellation outcomes.
- FR-022: Hosted portal page for subscription management.
- FR-023: Hosted portal supports plan change, seat change, cancellation, resume, payment method update, and invoice history.

### Webhooks

- FR-030: Generate signed webhook events.
- FR-031: Record delivery attempts.
- FR-032: Support retry policy.
- FR-033: Support duplicate, delayed, out-of-order, and replay delivery.
- FR-034: Show handler status and response body.

### Scenario Runner

- FR-040: Parse YAML scenarios.
- FR-041: Support local clock advancement.
- FR-042: Support expected event sequence assertions.
- FR-043: Support app callback assertions.
- FR-044: Produce JSON and Markdown reports.

### Dashboard

- FR-050: Show billing object lists and detail pages.
- FR-051: Show customer/subscription timeline.
- FR-052: Show webhook delivery detail.
- FR-053: Show scenario runs.
- FR-054: Export debug bundles.

### SaaS Profile

- FR-060: Provide a `saas` profile preset.
- FR-061: Model workspaces with owner, members, customer, subscription, and tenant payment rail.
- FR-062: Model SaaS plan tiers, product types, payment cycles, billing types, free trial, seat limits, and export limits.
- FR-063: Support subscription checkout, upgrade preview, upgrade confirm, cancellation, stop pending cancellation, payment failure, and retry.
- FR-064: Support additional seat estimate and purchase.
- FR-065: Support workspace member invitation, signed-up state, and deletion.
- FR-066: Support export summary, design-case export usage, export session creation, and session file lookup.
- FR-067: Support extra export preview, payment intent creation, payment state, and provision state.
- FR-068: Support customer portal URL generation with masked signed URL evidence.
- FR-069: Support payment history from invoice paid, payment intent succeeded, and charge refunded events.
- FR-070: Support platform and Connect webhook evidence, including duplicate, unresolved identity, and cache/workspace log outcomes.
- FR-071: Support back-office plan, seat, period, refund, member, extra export payment, and manual export limit actions.
- FR-072: Provide SaaS app assertion targets for subscription, plan, seats, members, export, extra export payment, payment history, webhook event, workspace log, and back-office refund.
- FR-073: Provide SaaS dashboard views for workspace overview, entitlement snapshot, seats/members, export entitlement, extra export payments, webhook timeline, payment history, and support operations.
- FR-074: Provide SaaS priority scenario examples.
- FR-075: Provide expected observability signal definitions for SaaS billing, export, and webhook scenarios.

## Non-Functional Requirements

- NFR-001: Local app starts with one Docker command.
- NFR-002: Backend runtime is Go.
- NFR-003: Frontend is React + TypeScript.
- NFR-004: Local persistent storage uses SQLite.
- NFR-005: No real card data is stored.
- NFR-006: Contract behavior is fixture-backed.
- NFR-007: Profile-specific behavior is fixture-backed and does not require production payment credentials.

## Non-Goals

- Real payment processing
- Complete Stripe API
- Production PSP replacement
- Long-term billing data warehouse
