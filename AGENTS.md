# Billtap Agent Instructions

Billtap uses spec-driven development and parallel agent orchestration.

## Required Reading Order

1. `.specify/memory/constitution.md`
2. `docs/FINAL_GOAL.md`
3. `docs/AGENT_ORCHESTRATION.md`
4. `specs/000-product/spec.md`
5. `specs/000-product/plan.md`
6. `specs/000-product/tasks.md`
7. `specs/000-product/gates.md`
8. Relevant docs under `docs/`

## Development Rules

- Build a full app, not just a mock API.
- Keep Stripe compatibility practical and fixture-backed.
- Do not build real payment processing.
- Use Go for backend runtime and React + TypeScript for UI.
- Any behavior affecting webhook order, signature, retry, idempotency, or billing state must have tests.
- Any production-facing feature must be optional, bounded, and safe.
- Agent work must declare ownership, verification, and gate status.

## Agent Handoff

Every agent handoff should include:

- Summary
- Files changed
- Verification
- Open risks
- Gate status

