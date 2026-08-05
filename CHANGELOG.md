# Changelog

## Unreleased

- Subscription item create and delete now accept `proration_behavior` and
  `proration_date` (previously `parameter_unknown` on create and silently
  ignored on delete), routing through the same proration path as subscription
  update so a seat add with `always_invoice` bills the prorated delta with the
  subscription's tax rates. Deleting the last item is rejected, delete accepts
  its parameters from the query string as well as the body, and item-level
  `tax_rates` are kept as evidence only — totals still come from the
  subscription's `default_tax_rates`.
- Invoice previews now apply tax: `create_preview` and `upcoming` read the
  subscription's `default_tax_rates` (or `automatic_tax`) and tax the
  post-discount proration base through the same helper the confirmed invoice
  serialization uses, so preview and confirmed amounts match field for field
  including decimal-rate rounding and inclusive rates. Preview item parsing
  also accepts `[price_id]`, which the validator already allowed but the
  parser ignored, and previews that cannot prorate now report
  `billtap_preview.proration_skipped_reason` instead of a silent zero.
- Subscription item changes now bill instead of only recording proration
  parameters as metadata: `proration_behavior=always_invoice` issues a paid
  `subscription_update` invoice and repoints `latest_invoice`,
  `billing_cycle_anchor=now` also resets the period and bills the new cycle
  net of the unused old-cycle credit, and `create_prorations` defers the delta
  to the next renewal invoice. Proration reuses the invoice preview's
  calculator so preview and actual agree, applies `default_tax_rates` after
  discounts, keeps `total == subtotal - discounts + tax` on every proration
  invoice, emits the full invoice/payment-intent webhook sequence, and returns
  HTTP 402 without committing the item change when
  `payment_behavior=error_if_incomplete` meets a failing outcome.
- Checkout session create now accepts session-level `metadata[...]` and
  round-trips it through retrieval and completion, restoring Stripe parity for
  SDK callers that attach the same map to both the session and
  `payment_intent_data[metadata]`. The two maps stay independent: session
  metadata is never promoted into the completed PaymentIntent.
- Added a promotion-code input to the hosted checkout page for
  `allow_promotion_codes` sessions, backed by Billtap-specific
  `POST/DELETE /api/checkout/sessions/{id}/promotion_code`: applying a valid
  code attaches the coupon discount to the open session (refreshing the
  Subtotal/Discount/Total breakdown before completion), invalid, inactive,
  expired, product-restricted-no-match, and duplicate applications return
  inline errors, and removal restores the original totals.
- Added a top-level `tax_rates` fixture pack key that seeds local tax-rate
  evidence (same path as disputes): explicit IDs are stored as-is so seeded
  products and checkout `default_tax_rates` can reference them, re-applying
  a pack overwrites by ID, and apply/validate summaries report the count.
- Added checkout `mode=payment` for one-time payments: `line_items` accept
  inline `price_data` (creating real local product/price evidence),
  `payment_intent_data` (`setup_future_usage`, `description`,
  `receipt_email`, `capture_method`, `metadata`) and `client_reference_id`
  are accepted, completion creates a single PaymentIntent (no subscription,
  no invoice) with `payment_status=paid` or `no_payment_required` for free
  totals, PaymentIntents expose `description`/`receipt_email`/
  `setup_future_usage` as top-level fields, and the hosted checkout page
  renders payment-mode sessions without subscription/invoice rows. `setup`
  mode remains rejected.
- Wired tax rates into real billing totals: checkout sessions accept
  `subscription_data[default_tax_rates][]` and subscriptions accept
  `default_tax_rates` on create/update (empty string clears), resolving
  `txr_*` IDs against local tax-rate evidence and snapshotting them onto the
  session, subscription, and invoices. Inclusive and exclusive rates use
  Stripe-style math after discounts through completion and renewal invoices,
  subscriptions serialize `default_tax_rates` as TaxRate objects, invoices
  emit per-rate `total_taxes`/`total_tax_amounts` with real rate IDs, and
  `automatic_tax` remains mutually exclusive with `default_tax_rates`.
- Added a simulated Stripe Tax surface in stripe-node v22 shapes: checkout
  sessions accept `automatic_tax[enabled]` and `tax_id_collection[enabled]`,
  the tax rate snapshots from customer metadata `tax_percent` at creation
  (absent metadata means 0% with status `complete`, matching Stripe's
  no-registration behavior) and applies exclusively after discounts through
  completion and renewal invoices. Sessions, invoices, and subscriptions
  serialize `automatic_tax` in their v22 shapes, invoices emit `total_taxes`
  alongside legacy `total_tax_amounts`, and `/v1/tax_rates` plus customer
  `/tax_ids` ship as local evidence stores. The Stripe Tax API family
  (`/v1/tax/*`) remains unimplemented.
- Added product-scoped coupons: `applies_to[products]`, `duration_in_months`,
  `redeem_by`, `max_redemptions`, and `times_redeemed` persist and round-trip;
  product-restricted discounts reject sessions/subscriptions with no matching
  line-item product and apply only to matching items through completion,
  renewal, and preview. `percent_off` now accepts decimals (e.g. `12.5`) with
  Stripe-compatible rounding.
- Checkout sessions expose `currency`, `amount_subtotal`, `amount_total`,
  `total_details`, and array-shaped `discounts`, and the hosted checkout page
  renders a Subtotal / Discount / Tax / Total breakdown from those fields.
  Money fields across the hosted pages now always parse as Stripe minor
  units, fixing sub-$10 amounts rendering 100x too large.
- Added invoice email evidence: `POST /v1/invoices/{id}/send` records
  metadata evidence and emits `invoice.sent` for open and paid invoices,
  finalizing a `send_invoice` invoice emits the same evidence in
  finalized-then-sent order, and `hosted_invoice_url` redirects to the
  invoice API resource. No real email is delivered.
- Hosted pages now repoint caller-provided localhost redirect targets at the
  run's configured public origin: when a run has a `public_base_url`, the
  hosted checkout "Return to app" link and billing portal return
  link/redirect swap the scheme/host/port of `localhost`/`127.0.0.1`
  `success_url`/`return_url` values for the run origin (path and query kept),
  surfaced via the `billtap_return_url` extension field and the portal URL
  query while stored sessions, `success_url`, and portal `return_url`
  response fields keep the caller's original values. External domains and
  unconfigured runs are untouched.
- Added run-scoped public base URLs so several proxied stacks can share one
  Billtap server: `POST /runs/<runId>/v1/config` pins `public_base_url` and an
  optional `public_base_path` per run, and absolute URLs (checkout
  `session.url`, billing portal URLs) prefer the run's base over an
  `X-Billtap-Public-Base-Url` request header, the forwarded proxy origin
  (`X-Forwarded-Proto`/`Host`/`Prefix`) on run-scoped requests, and the global
  `BILLTAP_PUBLIC_BASE_URL`. The default run keeps its previous behaviour.
- Added multi-workspace support so one running server can hold several fully
  isolated billing datasets. Requests select a workspace with the
  `X-Billtap-Workspace` header or `workspace` query parameter, unselected
  requests keep using the backward-compatible `default` workspace, named
  workspaces open their own SQLite database lazily under `workspaces/`, and
  `GET /workspaces` lists the known workspaces.
- Added a manual invoice-backed one-time payment flow for local SaaS usage
  charges, including `POST /v1/invoices`, `POST /v1/invoiceitems`,
  `POST /v1/invoices/{id}/finalize`, metadata preservation, expanded
  payment-intent evidence, and customer-level default outcomes for invoice pay.
- Fixed invoice pay requests so a submitted payment method no longer overrides
  a configured invoice-backed PaymentIntent outcome.
- Accepted Stripe-compatible `proration_date` on subscription updates and
  retained subscription update billing/proration parameters as local evidence.
- Added customer-level default PaymentIntent outcomes so fixture-seeded
  customers can drive confirmed one-time payment failures without changing
  app-created PaymentIntent metadata.
- Added deferred per-PaymentIntent outcome controls for one-time payment flows
  through `metadata[billtap_payment_intent_outcome]`, local create aliases, and
  `POST /api/payment_intents/{id}/outcome`.
- Hardened Stripe-like shape and validation for billing portal sessions and
  customer payment-method lists, including portal flow enum checks,
  PaymentMethod SDK fields, non-card filtering, and Stripe-style validation
  error envelopes.
- Added `GET /v1/prices/search` for the measured Stripe Search Query Language
  subset used by one-time price lookup paths: `active`, `type`, `lookup_key`,
  and metadata equality clauses joined by `AND`.
- Added customer fixture controls for empty or explicit payment-method lists so
  no-card local billing scenarios can remain deterministic until portal save.
- Expanded the Stripe-like simulation surface with hosted billing portal
  sessions/actions, local coupons and promotion codes, subscription schedules,
  SCA-required PaymentIntent callbacks, customer cash-balance funding for
  bank-transfer intents, dispute evidence, event filtering, endpoint-scoped
  delivery attempts, and grouped webhook replay controls.
- Added endpoint-scoped historical webhook replay so apps can catch up fixture
  events emitted before webhook endpoint registration.
- Licensed Billtap under Apache-2.0 and added a top-level `NOTICE`.
- Hardened public-readiness validation: JSON numeric request values now
  preserve decimal input for wrong-type rejection, subscription quantities no
  longer silently normalize invalid values, and update endpoints have explicit
  parameter validation.
- Expanded the compatibility scorecard to `l3-public-readiness-v2` with 28
  release-blocking cases covering request validation, idempotency mismatch, and
  deterministic checkout payment-error aliases.
- Added public release readiness documentation and clarified release evidence
  requirements.
- Added public compatibility and release process docs for the v0.1.0 release
  path.
- Expanded Connect compatibility with platform account retrieval, account
  deletion markers, people/persons evidence, and updated OpenAPI inventory
  tracking to `110 / 587` operations.
- Documented the supported Stripe-like subset, Billtap-specific APIs,
  unsupported provider behavior, and fixture/scenario/webhook claim boundaries.
- Prepared the repository for a public source-only release.
- Reworked public documentation around the Billtap default testing lane and Stripe testmode fallback lane.
- Moved company-specific adoption material and raw validation notes into ignored `.private/` storage.
- Sanitized the public workspace billing profile as `saas`.
- Added public docs for contribution, security, gate status, and documentation navigation.

## 0.0.0

- Initial source state with Go backend, React checkout/portal/dashboard apps, scenario runner, fixture APIs, webhook reliability controls, and Dockerfile.
