# Billtap Constitution

## Article 1: Billtap Is a Billing Lab, Not a Payment Processor

Billtap must never process real payments. It simulates billing flows, checkout, portal actions, webhook delivery, and app-side assertions.

## Article 2: Full App Means API, UI, Scenarios, and Debugging

Billtap must not stop at a mock API server. The product includes a hosted checkout UI, hosted billing portal UI, developer dashboard, scenario runner, webhook delivery system, and debugging tools.

## Article 3: Stripe Compatibility Must Be Practical and Tested

Billtap should emulate a useful Stripe subset, not the entire Stripe API. Compatibility claims require fixtures, contract tests, and documented limitations.

## Article 4: Billing State Must Be Explainable

Every customer, checkout session, subscription, invoice, payment intent, webhook event, and app callback must be explainable in a timeline.

## Article 5: Webhook Reliability Is Core Product Value

Duplicate delivery, retry, delay, out-of-order delivery, signature verification, idempotency, and replay must be first-class concepts.

## Article 6: Scenario Tests Are Product Features

Billtap exists to make billing flows deterministic. Scenario DSL, local clock, replay, and assertions are core features.

## Article 7: Production Boundaries Must Stay Clear

Production-facing modes may inspect or relay test/staging billing events, but Billtap must not become a hidden production payment dependency.

## Article 8: Agent Work Requires Gates

Parallel agent development is encouraged, but every lane needs explicit ownership, fixtures, tests, and gate status.

## Article 9: Decisions Must Be Written Down

Architecture, compatibility, storage, security, production boundary, and API-surface decisions must be recorded under `docs/decisions/`.

