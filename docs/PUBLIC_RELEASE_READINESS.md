# Public Release Readiness

Status date: 2026-05-09

Billtap has a clear community-facing strength: it is a stateful local billing
lab for deterministic subscription, fixture, scenario, and webhook reliability
tests. The public claim must stay narrower than full Stripe parity.

## Public Claim

Billtap can be described as:

> A source-first local billing sandbox, also available as a GHCR container
> image, with a documented Stripe-like API subset, stateful
> checkout/subscription fixtures, deterministic payment-error simulation,
> webhook reliability controls, scenario reports, and a dashboard for
> explaining billing evidence.

Do not describe Billtap as:

- a real payment processor
- full Stripe compatibility
- a replacement for Stripe testmode
- a hosted Stripe Checkout or Billing Portal clone
- safe for real card data or live payment credentials

## Current Strengths

- Stateful subscription graphs: customers, products, prices, checkout sessions,
  subscriptions, invoices, payment intents, and timeline evidence can persist
  locally.
- Provider injection: apps that already support `STRIPE_API_BASE_URL` and
  `Stripe-Signature` can use Billtap without service-code branches.
- Fixture workflow: JSON/YAML fixture packs can apply, snapshot, and assert
  deterministic billing state.
- Error simulation: checkout completion can simulate success, card declines,
  insufficient funds, expired cards, incorrect CVC, processing errors, missing
  payment methods, and authentication-required states.
- Webhook lab: delivery attempts record signature evidence, retry scheduling,
  duplicate delivery, delay, out-of-order metadata, replay, endpoint failure,
  timeout, and signature-mismatch evidence.
- Scorecard: the offline compatibility scorecard writes JSON, Markdown, and
  replay bundles and must pass with `mismatch=0` and `error=0`.

## Release Blockers

| Blocker | Status | Resolution |
| --- | --- | --- |
| License | Resolved | Apache-2.0 `LICENSE` and top-level `NOTICE` are present. |
| Public repository URL | Confirm before release | `docs/RELEASE.md` currently assumes `https://github.com/midagedev/billtap.git`. |
| Release verification | Must pass per release | Run `docs/RELEASE_CHECKLIST.md` on the release branch or tag. |

## Required Evidence

Every public release candidate should include:

- `go test ./...`
- `go run ./cmd/billtap compatibility scorecard --output-dir /tmp/billtap-compatibility`
- `npm run typecheck`
- `npm run build`
- `npm run smoke:sample`
- `npm run smoke:sdk`
- `npm run smoke:web:install`
- `npm run smoke:web`
- `go build -o /tmp/billtap ./cmd/billtap`
- `docker build -t billtap:local .`
- post-publish `docker pull ghcr.io/midagedev/billtap:<tag>` and health smoke
- sample app scenario smoke for `examples/subscription-payment-retry.yml`
- sample app scenario smoke for `examples/saas-adoption-contract.yml`

The scorecard evidence should state:

- scorecard version
- total cases
- release-blocking cases
- imported/skipped/unsupported/mismatch/error counts
- whether `passed` is `true`

Current local evidence on 2026-05-09:

- Scorecard version: `l3-public-readiness-v2`
- Scorecard result: `imported=28 skipped=1 unsupported=1 mismatch=0 error=0`
- Local release checks listed in `docs/GATE_STATUS.md` passed on branch
  `codex/apache-license-release-readiness`.
- Apache-2.0 `LICENSE`, `NOTICE`, `package.json`, and `package-lock.json`
  metadata are aligned.

## Compatibility Positioning

Use Stripe references as compatibility input, but keep the Billtap claim tied to
Billtap-owned tests:

- Stripe API errors define the shape and categories Billtap mirrors for local
  errors.
- Stripe idempotent request behavior informs Billtap's same-process POST replay
  and parameter mismatch behavior.
- Stripe webhook docs inform duplicate, retry, raw-body signature, and failure
  handling simulations.
- `stripe/stripe-mock` is useful as a route/parameter sanity oracle, but its
  own README states it is stateless, hardcoded, and not meant for specific
  response/error testing. Billtap's differentiator is the stateful local
  behavior around that gap.

## Next Public-Quality Work

- Add a short screencast or screenshot set after dashboard UI verification.
- Add an optional external oracle lane that runs selected route/parameter cases
  against `stripe-mock` when Docker/network access is available.
- Add signed binary, Homebrew, or npm artifact automation only when the project
  is ready to claim those distribution channels.
