# Tasks: Billtap Product Foundation

## Phase 0: Repository Foundation

- [x] T000 Create spec-driven repository structure
- [x] T001 Write project constitution
- [x] T002 Write product specification
- [x] T003 Write implementation plan
- [x] T004 Write final goal and PRD
- [x] T005 Write UI product, testing, roadmap, and production boundary docs
- [x] T006 Write agent orchestration and gates
- [x] T007 Write SaaS service feature profile

Gate:

- [x] G0 Spec Readiness review

## Phase 1: Runtime Core

- [x] T010 Create Go module
- [x] T011 Create server skeleton
- [x] T012 Add config loading
- [x] T013 Add SQLite storage and migrations
- [x] T014 Add in-memory test storage
- [x] T015 Add health endpoints
- [x] T016 Add frontend workspace

Suggested agents:

- Runtime Core Agent
- Release Agent

Gate:

- [x] G1 Runtime Contract

## Phase 2: Billing Core and Checkout MVP

- [x] T020 Implement customer model and APIs
- [x] T021 Implement product and price models and APIs
- [x] T022 Implement checkout session model and APIs
- [x] T023 Implement subscription model
- [x] T024 Implement invoice model
- [x] T025 Implement payment intent model
- [x] T026 Build hosted checkout UI
- [x] T027 Wire checkout outcomes to billing state

Suggested agents:

- Billing Engine Agent
- Checkout UI Agent
- Compatibility Agent

Gate:

- [x] G2 Checkout MVP

## Phase 3: Webhook Lab

- [x] T030 Implement webhook endpoint model and APIs
- [x] T031 Implement event model
- [x] T032 Implement HMAC signature generation
- [x] T033 Implement delivery attempts
- [x] T034 Implement retry policy
- [x] T035 Implement duplicate delivery
- [x] T036 Implement delayed delivery
- [x] T037 Implement out-of-order delivery
- [x] T038 Implement replay

Suggested agents:

- Webhook Agent
- Security Agent

Gate:

- [x] G3 Webhook Reliability

## Phase 4: Developer Dashboard

- [x] T040 Build dashboard shell
- [x] T041 Add object list/detail views
- [x] T042 Add customer/subscription timeline
- [x] T043 Add webhook delivery detail
- [x] T044 Add app response view
- [x] T045 Add debug bundle export

Suggested agents:

- Dashboard Agent

Gate:

- [x] G4 Debuggability

## Phase 5: Scenario Runner and CI

- [x] T050 Define scenario YAML schema
- [x] T051 Implement scenario parser
- [x] T052 Implement local clock
- [x] T053 Implement scenario executor
- [x] T054 Implement app callback assertions
- [x] T055 Implement JSON and Markdown reports
- [x] T056 Implement exit code policy

Suggested agents:

- Scenario Agent
- Assertion Agent
- CI Agent

Gate:

- [x] G5 CI Contract

## Phase 6: Billing Portal

- [x] T060 Build hosted portal UI
- [x] T061 Implement plan change
- [x] T062 Implement seat change
- [x] T063 Implement cancellation
- [x] T064 Implement resume
- [x] T065 Implement payment method update simulation
- [x] T066 Implement invoice history

Suggested agents:

- Portal UI Agent
- Billing Engine Agent

Gate:

- [x] G6 Portal Coverage

## Phase 7: SaaS Service Profile

- [x] T070 Implement SaaS workspace model and seed fixtures
- [x] T071 Implement tenant payment rail policy for card, out-of-band, and connect
- [x] T072 Implement SaaS plan catalog preset
- [x] T073 Implement subscription checkout, upgrade, cancel, failure, and retry scenarios
- [x] T074 Implement additional seat estimate and purchase behavior
- [x] T075 Implement member invite, signed-up, and deletion behavior
- [x] T076 Implement export summary, usage, session, and file behavior
- [x] T077 Implement extra export preview, payment intent, and provision behavior
- [x] T078 Implement payment history and customer portal evidence
- [x] T079 Implement support/back-office actions and debug bundle
- [x] T080 Implement platform/connect webhook evidence
- [x] T081 Add SaaS app assertion targets
- [x] T082 Add SaaS priority scenario examples
- [x] T083 Add SaaS observability expectation fixtures

Suggested agents:

- SaaS Profile Agent
- Scenario Agent
- Dashboard Agent
- Webhook Agent

Gate:

- [x] G7 SaaS Adoption Contract

## Phase 8: Production Boundaries

- [x] T084 Implement secret masking
- [x] T085 Implement retention policy
- [x] T086 Implement audit log
- [x] T087 Implement safe testmode relay spike
- [x] T088 Add production boundary tests
- [x] T089 Add deployment runbook

Suggested agents:

- Security Agent
- Release Agent

Gate:

- [x] G8 Production Boundary

## Phase 9: Release Candidate

- [x] T090 Add Dockerfile
- [x] T091 Add sample app
- [x] T092 Add example scenarios
- [x] T093 Add release checklist
- [x] T094 Run all gates

Suggested agents:

- Release Agent
- Spec Lead

Gate:

- [x] G9 Release Candidate

## Phase 10: SaaS Repo Replacement Smoke

- [x] T100 Add Stripe-compatible product search and filtered catalog responses
- [x] T101 Add Stripe-compatible price, checkout session, subscription, invoice, portal, and payment-method response projections
- [x] T102 Add metadata and lookup-key persistence for seeded catalog data
- [x] T103 Add sample-app Billtap compose overlay and SaaS catalog seed script
- [x] T104 Smoke Billtap through sample-app's Stripe compatibility proxy

Suggested agents:

- Compatibility Agent
- SaaS Integration Agent

Gate:

- [x] G10 sample-app Billtap Replacement Smoke

## Phase 11: Fixture and Assertion Ergonomics

- [x] T110 Add data-driven fixture pack apply API
- [x] T111 Add fixture-scoped snapshot API
- [x] T112 Add built-in fixture assertion API with structured pass/fail reports
- [x] T113 Make fixture-applied subscription graphs use normal checkout completion state
- [x] T114 Replace sample-app Billtap shell seeding with a YAML fixture pack
- [x] T115 Add sample-app Billtap fixture/assert one-shot compose gate

Suggested agents:

- Compatibility Agent
- Assertion Agent
- SaaS Integration Agent

Gate:

- [x] G11 Fixture Assertion Ergonomics

## Phase 12: Public Release Readiness

- [x] T120 Tie public compatibility claims to automated tests and release-blocking scorecard cases
- [x] T121 Fix JSON numeric validation so decimal integer fields are rejected instead of truncated
- [x] T122 Add explicit validation for supported update endpoints and subscription item quantities
- [x] T123 Expand deterministic payment-error scorecard cases for documented checkout aliases
- [x] T124 Document source-only release state, release blockers, and public claim boundaries
- [x] T125 Add project-owner-selected `LICENSE` before public community release

Suggested agents:

- Compatibility Agent
- Release Agent
- Security Agent

Gate:

- [x] G12 Public Release Readiness

## Phase 13: Broader Stripe API Compatibility Roadmap

- [x] T130 Define compatibility levels for inventory, schema validation, fixture responses, stateful behavior, scenarios, webhooks, and SDK smoke
- [x] T131 Document endpoint family priorities beyond the DS/SaaS adoption path
- [x] T132 Define OpenAPI inventory, scorecard expansion, and optional oracle lanes
- [x] T133 Implement Stripe OpenAPI inventory generator and coverage matrix
- [x] T134 Add protocol baseline tests for pagination, expand, request IDs, API version, and idempotency traces across supported endpoints
- [x] T135 Add OpenAPI-backed validation for broad L1 route/parameter/type coverage
- [x] T136 Add direct PaymentIntent and SetupIntent state machines
- [x] T137 Add renewal/test-clock/retry mutation for subscriptions and invoices
- [ ] T138 Add coupons, discounts, credit notes, refunds, disputes, and payment-history simulation
- [ ] T139 Add Connect/platform smoke fixtures and connected-account webhook routing
- [ ] T140 Add official Stripe SDK matrix for Node, Go, Java, Python, and Ruby
- [x] T141 Add optional scheduled/manual Stripe OpenAPI inventory workflow artifacts

Suggested agents:

- Compatibility Agent
- Billing Engine Agent
- Webhook Agent
- Scenario Agent
- Release Agent

Gate:

- [ ] G13 Stripe API Compatibility Expansion Readiness
