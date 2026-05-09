# UI Product

## UI Surfaces

Billtap has three UI surfaces.

## 1. Hosted Checkout UI

Purpose: simulate a Stripe Checkout-like user flow.

Required screens:

- checkout session summary
- customer identity
- plan and price
- payment method selector
- outcome selector
- confirmation
- failure and retry
- cancel return

Outcome options:

- payment succeeds
- card declined
- insufficient funds
- expired card
- requires action
- user cancels
- async payment pending

Developer value:

- QA can manually exercise billing flows.
- Engineers can choose failure modes without writing scripts.
- Checkout completion emits realistic events.

## 2. Hosted Billing Portal UI

Purpose: simulate customer self-service billing management.

Required screens:

- current subscription
- plan change
- seat quantity change
- cancel subscription
- resume subscription
- payment method update
- invoice history

Outcome options:

- immediate upgrade
- scheduled downgrade
- cancel at period end
- immediate cancellation
- payment method update succeeds
- payment method update fails

## 3. Developer Dashboard

Purpose: debug the full billing integration.

Required views:

- overview
- workspaces
- workspace entitlement snapshot
- seats and members
- export entitlement
- extra export payments
- customers
- subscriptions
- invoices
- payment intents
- checkout sessions
- webhook timeline
- delivery attempts
- scenario runs
- app assertions
- debug bundles
- settings

## Key Dashboard Interactions

### Timeline

Shows object changes and webhook deliveries in one ordered stream.

Each timeline entry should include:

- timestamp
- object
- transition
- event type
- delivery status
- app response
- linked scenario step

### Webhook Detail

Shows:

- event JSON
- signature header
- endpoint URL
- attempt count
- response status
- response body
- retry schedule
- idempotency key hints

### Scenario View

Shows:

- scenario file
- current step
- emitted events
- clock advances
- assertions
- pass/fail status

### Debug Bundle

Exportable bundle for:

- customer
- subscription
- workspace
- entitlement snapshot
- export session
- extra export payment
- checkout session
- invoice
- payment intent
- webhook event
- scenario run

Bundle includes:

- object timeline
- emitted webhooks
- delivery attempts
- app responses
- assertions
- replay instructions

### SaaS Workspace View

Shows:

- workspace id, owner, tenant route, and payment rail
- checkout/admin management availability
- plan tier, cycle, subscription status, and period
- basic seats, additional seats, used seats, and pending invitations
- export default limit, total limit, total remaining, manual limit, and renewal date
- extra export payment status and provision status
- latest payment history item
- latest platform/connect webhook event
- latest workspace log

Primary actions:

- open hosted checkout
- open hosted portal
- run upgrade preview
- confirm upgrade
- schedule cancellation
- stop pending cancellation
- estimate seat purchase
- purchase seats
- invite member
- delete member
- preview extra export
- create extra export payment intent
- provide extra export
- replay webhook
- export support bundle

## UX Principles

- Operational, dense, and scan-friendly.
- Avoid marketing-style layouts.
- Show state transitions clearly.
- Make failure reasons copyable.
- Always link webhook events back to billing objects.
