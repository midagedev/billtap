# Billtap Documentation

This directory contains the public product, architecture, testing, and release notes for Billtap.

## Start Here

- `../README.md`: project overview, adoption model, quick start, and release state
- `FINAL_GOAL.md`: product north star and final release definition
- `ARCHITECTURE.md`: runtime components and module boundaries
- `COMPATIBILITY.md`: supported and unsupported Stripe-like behavior
- `API_VALIDATION_AND_ERROR_SIMULATION.md`: Stripe-like validation and
  deterministic error-simulation target
- `TESTING.md`: verification strategy and scenario coverage
- `PRODUCTION_BOUNDARIES.md`: real-payment and relay-mode safety boundaries

## Contracts

- `../specs/000-product/contracts/api.md`: Stripe-like API and Billtap dashboard API
- `../specs/000-product/contracts/scenario.md`: YAML scenario format and actions
- `../specs/000-product/contracts/webhooks.md`: event shape, signatures, retry, replay, and production boundary behavior

## Operations

- `runbooks/local-dev.md`: local development workflow
- `runbooks/production-boundary.md`: controlled relay and staging-adjacent guidance
- `RELEASE.md`: public release procedure for v0.1.0 and later pre-1.0 releases
- `RELEASE_CHECKLIST.md`: checks to run before publishing a release
- `GATE_STATUS.md`: public gate snapshot

## Product Notes

- `CONCEPT.md`: product concept and testing loop
- `PRD.md`: product requirements
- `UI_PRODUCT.md`: checkout, portal, and dashboard product notes
- `SAAS_PROFILE.md`: generic SaaS workspace profile used by public examples
- `ROADMAP.md`: forward-looking work

## Decisions

- `decisions/0001-project-shape.md`
- `decisions/0002-stack.md`

Internal adoption notes and original company-specific validation material are preserved locally under `.private/`, which is intentionally ignored by Git.
