# Support

Billtap is source-only and pre-release. Support is best effort unless a paid or private arrangement is documented outside this repository.

## Choose The Right Channel

Use GitHub issues for:

- reproducible bugs in the Billtap API, dashboard, checkout, portal, scenario runner, fixtures, or webhook delivery
- documentation bugs
- feature requests that fit the local/CI billing sandbox scope
- compatibility gaps in the documented Stripe-like subset

Use pull requests for:

- focused fixes with tests or clear verification
- documentation improvements
- fixture-backed compatibility improvements
- safety improvements around masking, retention, relay mode, or production boundaries

Use the security process in `SECURITY.md` for:

- secret exposure
- real card-data handling regressions
- webhook signature bypasses
- relay-mode boundary bypasses
- vulnerabilities that should not be disclosed publicly before a fix is available

## What Billtap Can Help With

- deterministic local and CI billing scenarios
- checkout, portal, subscription, invoice, payment intent, and webhook test flows
- webhook retries, duplicate delivery, delay, out-of-order delivery, and replay
- fixture apply/snapshot/assert workflows
- generic SaaS workspace billing fixtures and assertions
- debugging app-side billing behavior with redacted evidence

## Out Of Scope

Billtap is not support for live payment operations. Do not ask the project to:

- process real card payments
- handle live customer payment data
- replace Stripe Billing, Stripe Checkout, or another payment provider in production
- guarantee full Stripe API compatibility
- debug provider-specific settlement, tax, risk, dispute, payout, or account behavior

If you are unsure whether something belongs in Billtap, open a feature request or discussion-style issue with a sanitized example and no real payment data.

## Before Opening An Issue

Include:

- Billtap version or commit SHA
- operating system and runtime versions
- command or scenario file used
- expected behavior
- actual behavior
- sanitized logs, reports, or screenshots
- whether the issue affects local development, CI, or controlled staging-adjacent relay mode

Do not include:

- API keys, webhook signing secrets, tokens, cookies, or credentials
- real customer identifiers
- real card data
- production payment payloads
- private company names or internal hostnames

## Response Expectations

Maintainers prioritize issues that are reproducible, fixture-backed, and aligned with the documented product boundary. Pre-release APIs may change, and maintainers may close issues that request real payment processing, production payment dependency, or undocumented full-provider parity.
