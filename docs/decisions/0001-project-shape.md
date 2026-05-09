# Decision 0001: Build Billtap as a Full Billing Lab

## Status

Accepted

## Context

A plain Stripe mock server would compete directly with `stripe-mock` and other existing tools. The stronger product opportunity is a full app that combines API compatibility, checkout UI, billing portal UI, webhooks, scenarios, and app-side assertions.

## Decision

Billtap will be a full billing lab, not just a mock API server.

## Consequences

Positive:

- Stronger product differentiation
- Better fit for real subscription debugging
- More useful for QA and CI

Negative:

- Larger implementation scope
- Requires stronger orchestration and gates
- Needs UI quality, not only API correctness

