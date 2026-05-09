# Success Gates

## Gate Summary

| Gate | Name | Blocks |
| --- | --- | --- |
| G0 | Spec Readiness | implementation start |
| G1 | Runtime Contract | domain integration |
| G2 | Checkout MVP | webhook lab |
| G3 | Webhook Reliability | CI and release claims |
| G4 | Debuggability | MVP release |
| G5 | CI Contract | adoption by another repo |
| G6 | Portal Coverage | full app claim |
| G7 | SaaS Adoption Contract | SaaS repo adoption |
| G8 | Production Boundary | staging or production-adjacent use |
| G9 | Release Candidate | public release |
| G10 | sample-app Billtap Replacement Smoke | SaaS repo replacement proof |
| G11 | Fixture Assertion Ergonomics | developer integration test adoption |
| G12 | Public Release Readiness | community release |

## G0: Spec Readiness

Pass criteria:

- final goal documented
- full app scope documented
- non-goals documented
- Go + React stack accepted
- API, UI, scenario, and webhook contracts drafted
- agent orchestration documented

## G1: Runtime Contract

Pass criteria:

- Go server starts
- config loads from env and file
- SQLite migrations run
- in-memory storage works in tests
- health endpoint exists
- React asset serving path exists

## G2: Checkout MVP

Pass criteria:

- customer/product/price APIs work
- checkout session can be created
- hosted checkout page opens
- success and failure outcomes work
- subscription, invoice, and payment intent state is created
- dashboard can show checkout result

## G3: Webhook Reliability

Pass criteria:

- signed webhooks are delivered
- delivery attempts are recorded
- retries are configurable
- duplicate delivery is supported
- delayed delivery is supported
- out-of-order delivery is supported
- replay is supported
- handler response is visible

## G4: Debuggability

Pass criteria:

- dashboard shows object list and detail
- timeline connects billing objects and webhooks
- webhook detail shows signature, request, response, attempt count, and retry plan
- debug bundle export works
- failure reasons are copyable

## G5: CI Contract

Pass criteria:

- scenario runner works headlessly
- local clock advances deterministic scenarios
- JSON and Markdown reports are generated
- app callback assertions can fail the run
- exit codes are tested

## G6: Portal Coverage

Pass criteria:

- portal UI opens
- plan change works
- seat change works
- cancellation works
- resume works
- payment method update simulation works
- invoice history is visible

## G7: SaaS Adoption Contract

Pass criteria:

- `saas` profile can seed tenant, workspace, customer, plan, seats, and export fixtures
- workspace subscription checkout, upgrade preview, confirm, cancel, and retry scenarios pass
- additional seat estimate and purchase scenarios pass
- member invite and seat-limit scenarios pass
- export summary, export usage, export session, and extra export scenarios pass
- payment history includes invoice paid, payment intent succeeded, and charge refunded evidence
- platform and Connect webhook scenarios record duplicate, unresolved, and resolved outcomes
- support bundle includes workspace, subscription, seat, export, payment, webhook, workspace log, and app assertion evidence
- expected billing/export/webhook observability signals are defined for the scenario pack
- at least one application server local scenario can run without calling Stripe testmode

## G8: Production Boundary

Pass criteria:

- no real card processing path exists
- secrets are masked
- retention policy exists
- audit log records replay and delivery overrides
- production-adjacent relay mode is documented and tested
- raw payload storage is disabled by default in relay mode

## G9: Release Candidate

Pass criteria:

- all previous gates pass or have accepted deferral ADRs
- Docker image builds
- quickstart works on clean machine
- sample app works
- example scenarios pass
- docs match actual behavior

## G10: sample-app Billtap Replacement Smoke

Pass criteria:

- Billtap exposes the Stripe-compatible catalog, price, checkout, subscription, invoice preview, portal, and payment-method surfaces needed by sample-app smoke paths
- sample-app can select a Billtap-backed compose overlay without changing the default stripe-mock path
- SaaS catalog fixtures seed into Billtap with stable product and price IDs
- sample-app's existing Stripe search proxy can forward stateful calls to Billtap
- A compose smoke proves product search, price lookup, checkout completion, portal session creation, and payment-method listing through the sample-app overlay

## G11: Fixture Assertion Ergonomics

Pass criteria:

- Billtap can apply JSON/YAML fixture packs for customers, products, prices, and subscription graphs
- Fixture-applied objects carry run, namespace, fixture name, tenant, and ref metadata sufficient to isolate local/CI runs
- Billtap can return fixture-scoped snapshots with object counts and timeline evidence
- Billtap can assert expected customers, products, prices, subscriptions, invoices, payment intents, and timeline evidence with structured pass/fail reports
- sample-app's Billtap overlay uses a data-driven SaaS fixture pack rather than imperative catalog shell calls
- sample-app's Billtap overlay has a one-shot assert gate that checks both Billtap persisted state and the existing Stripe compatibility proxy surface before `application server` starts

## G12: Public Release Readiness

Pass criteria:

- public compatibility claims are tied to tests or release-blocking scorecard cases
- scorecard has `mismatch=0`, `error=0`, and `passed=true`
- supported request validation rejects unknown, wrong-type, missing, and invalid enum/quantity cases for the public subset
- deterministic payment-error simulation covers the documented checkout aliases
- release docs name source-only distribution state and artifact boundaries
- release checklist includes sample-app scenario preconditions
- public docs do not imply real payment processing, full Stripe parity, or production payment dependency
- `LICENSE` exists before public community release
