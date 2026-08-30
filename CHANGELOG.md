# Changelog

## Unreleased

- `DELETE /v1/test_helpers/test_clocks/{id}` removes a test clock and detaches
  it from the customers and subscriptions that referenced it (Stripe deletes
  attached objects; Billtap keeps them detached so fixture packs can re-seed
  against a recreated clock reusing the same id).
- Re-applying a fixture pack now restores a clock that ran past the pack's
  declared `frozen_time`: the clock is recreated at the declared time before
  subscriptions re-seed, so a consumed premise — for example a trial that
  already activated — comes back instead of surviving until the database is
  rebuilt. Backwards advance stays rejected.
- Discount objects now serialize a Stripe SDK 31 `source` sub-object
  (`source.type` plus `source.coupon` or `source.promotion_code`) alongside the
  legacy top-level `coupon`, so SDK 31 clients can read the discount origin.
- `POST /v1/prices/{id}` accepts `transfer_lookup_key`: moving a lookup key
  that another price holds clears it there first, matching Stripe's transfer
  behavior. Without the flag, updates keep the existing set-in-place semantics.
- The OpenAPI inventory moves from `200 / 587` (`34.1%`) to `201 / 587`
  (`34.2%`).
- The billing family is now full inventory (`39 / 39` OpenAPI operations):
  all four P0 families (webhooks, checkout, billing, billing_portal) are at
  100%.
- `POST /v1/invoices/{id}/lines/{line_item_id}` updates one draft-invoice line
  (`amount`, `description`, metadata) with the same totals recomputation as
  `update_lines`.
- `POST /v1/invoices/{id}/attach_payment` records an attached PaymentIntent
  (customer-matched) or payment record on a draft invoice as evidence;
  collection still runs through `finalize`/`pay`.
- `POST /v1/subscriptions/{id}/migrate` records
  `billing_mode[type]=flexible` (plus optional `proration_discounts`) as
  subscription metadata evidence; flexible-billing proration recalculation is
  not modeled.
- `POST /v1/payment_intents/{id}` updates intent metadata and evidenced
  description; amount, currency, and status stay immutable.
- `GET /v1/credit_notes/{id}/lines` returns one line derived from the stored
  credit-note amount and memo/reason.
- The OpenAPI inventory moves from `195 / 587` (`33.2%`) to `200 / 587`
  (`34.1%`): billing `39 / 39`, payments `23 / 41`, payment_history `14 / 30`.
- `GET /v1/checkout/sessions/{id}/line_items` returns the session's line items
  as expanded Stripe `item` objects (pre-discount amounts; per-line
  discount/tax splits are not modeled).
- `POST /v1/checkout/sessions/{id}` updates an open session: metadata merge
  plus `line_items[N][quantity]` overrides with immutable prices. Non-open
  sessions, out-of-range indexes, and non-positive quantities are rejected.
  Checkout sessions are now full inventory: `6 / 6` operations.
- `POST /v1/invoices/{id}` updates draft invoices (description, `days_until_due`,
  `default_payment_method`, metadata merge) and `DELETE /v1/invoices/{id}`
  deletes a draft together with its attached lines and timeline evidence.
- `POST /v1/invoices/{id}/add_lines`, `update_lines`, and `remove_lines` mutate
  draft-invoice lines through the invoice-item path, recomputing
  `subtotal`/`total`/`amount_due` per change.
- `DELETE /v1/products/{id}` removes local product evidence; existing prices
  keep referencing the deleted product id.
- The OpenAPI inventory moves from `187 / 587` (`31.9%`) to `195 / 587`
  (`33.2%`): checkout closes to `6 / 6`, billing moves to `36 / 39`, and
  catalog to `28 / 54`.
- Added billing portal configurations: `GET/POST
  /v1/billing_portal/configurations` and `GET/POST
  /v1/billing_portal/configurations/{id}` store local portal-configuration
  evidence with Stripe-shaped `business_profile`, `default_return_url`,
  `login_page`, `features[...]`, and `metadata` fields, list filters for
  `active`/`is_default`, and the first configuration in a run becoming the
  default. The hosted portal does not render configuration features.
- `GET /v1/subscription_items` lists resolved items across subscriptions with
  the `subscription` filter and the standard list envelope, and
  `GET /v1/subscription_items/{id}` retrieves one item.
- `POST /v1/subscription_items/{id}` updates an item's `price` and `quantity`
  (plus metadata and evidence-only `tax_rates`), routing
  `proration_behavior=always_invoice/create_prorations` through the same
  proration path as item create so a quantity bump with `always_invoice`
  issues the prorated `subscription_update` invoice.
- `GET/DELETE /v1/customers/{id}/subscriptions/{id}/discount` now mirror the
  top-level subscription discount routes and are customer-scoped: a mismatched
  customer gets `404`.
- Fixed subscription discount deletion: `DELETE .../subscriptions/{id}/discount`
  built a metadata copy with the discount keys removed, but subscription
  metadata patches merge, so the stored discount survived its own deletion and
  a following `GET` still returned it. The patch now zeroes the discount keys,
  which is how merge-mode patches delete.
- Registered compat claims for already-tested routes: invoice `void`,
  invoice `mark_uncollectible`, and checkout `expire`. The OpenAPI inventory
  moves from `175 / 587` (`29.8%`) to `187 / 587` (`31.9%`): billing
  `31 / 39`, billing_portal `5 / 5`, checkout `4 / 6`.
- Local evidence objects — coupons, promotion codes, subscription schedules,
  disputes, tax rates, tax IDs and customer cash balances — are now stored in
  the run's own database instead of process memory. They were the only objects
  that did not survive a restart, so a restarted server kept answering with the
  rest of its data while every lookup that needed one of them failed: a
  subscription whose `default_tax_rates` referenced a tax rate created before
  the restart returned a `resource_missing` error, with nothing in the
  surviving data to suggest why. Runs backed by memory stay ephemeral, which
  is what they were always for. Idempotency keys remain in memory on purpose —
  losing them on restart is the Stripe-like behaviour.
- `TestSQLiteMigrationsRun` now derives the expected versions from the embedded
  migration files rather than a hand-written list, so it no longer needs an edit
  per migration and it fails on a gap or a duplicated number.
- `POST /v1/invoices/{id}/void` moves an `open` invoice to `void`, records
  `billtap_voided_at`, and emits `invoice.voided`. Other statuses return
  `invalid_request_error` with `status must be open`.
- `POST /v1/invoices/{id}/mark_uncollectible` moves an `open` invoice to
  `uncollectible`, records `billtap_marked_uncollectible_at`, and emits
  `invoice.marked_uncollectible`.
- `POST /v1/checkout/sessions/{id}/expire` moves an `open` session to
  `expired` and emits `checkout.session.expired`. Non-open sessions return
  `invalid_request_error` with `status must be open`.
- `POST /v1/subscriptions` accepts `proration_behavior`
  (`none` / `create_prorations` / `always_invoice`) and nested
  `payment_settings[...]` keys, storing the received values in metadata.
- `POST /v1/subscriptions/{id}` accepts `cancel_at` as a unix timestamp,
  stores it, and echoes `cancel_at` without changing status immediately.
- `POST /v1/customers/{id}` accepts
  `invoice_settings[default_payment_method]` and echoes it on
  `invoice_settings.default_payment_method`. Attaching a payment method sets
  that default when the customer has none.
- `GET /v1/subscriptions` honors `current_period_end[gte]` and
  `current_period_end[lt]` as unix seconds.
- List endpoints honor `starting_after` and set `has_more` when a `limit`
  truncates remaining items.
- Invoice `payment_settings[payment_method_options][customer_balance][...]`
  is stored and echoed on
  `payment_settings.payment_method_options.customer_balance`.
- `POST /v1/invoiceitems` can omit `invoice`. The item is stored as a pending
  customer item (`invoice` is null) and `subscription` is accepted and echoed.
- `POST /v1/invoices` `pending_invoice_items_behavior=include` attaches that
  customer's pending items of the same currency and adds them to
  `subtotal`/`total` once at create. `exclude` (and omit) leaves them
  unattached.
- `POST /v1/invoices` accepts `subscription` (stored and echoed) and
  `days_until_due` (response `due_date` is created plus that many days).
- Invoice responses populate `lines.data` from the invoice's items using the
  same serialization as `GET /v1/invoices/{id}/lines`, including Wave 0
  `pricing` and `quantity`.
- Paid invoice responses set `post_payment_credit_notes_amount` to the sum of
  issued (non-void) credit notes' `credit_amount` for that invoice.
- Invoice preview accepts `subscription_details[trial_end]` (`now` or a unix
  timestamp). For a `trialing` subscription, `trial_end=now` (or a timestamp
  that is not in the future) previews the first paid cycle instead of a
  zero-amount remaining-trial window. `POST /v1/subscriptions/{id}` with the
  same `trial_end` moves the subscription to `active` and invoices that same
  first-cycle amount so preview and confirm agree.
- `POST /v1/refunds` treats a missing `amount` as a full refund of the
  remaining refundable balance on the invoice or payment intent. Requested
  amounts still cannot exceed that remaining balance.
- `POST /v1/credit_notes` accepts `memo`, `out_of_band_amount`, and
  `refund_amount`, persists them, and echoes them (plus derived
  `credit_amount`) on create and retrieve. `out_of_band_amount` is external
  settlement and does not change customer cash balance.
- `POST /v1/invoiceitems` now accepts `pricing[price]` and `quantity` as an
  alternative to `amount` (`amount` cannot be sent with `pricing` or
  `quantity`). The line amount is the price's `unit_amount` times `quantity`
  (default 1), currency comes from the price, a missing price is
  `resource_missing`, and the create response includes `quantity` plus a
  `pricing` object (`type=price_details`) when a price id was used.
- Asynchronous webhook deliveries are now waitable. They ran on a detached
  context with no lifecycle handle, so they could still be writing delivery
  records after a caller tore down its storage — which surfaced as intermittent
  `TempDir ... directory not empty` test failures. `WaitForAsyncDeliveries`
  drains in-flight deliveries, the server drains every run's deliveries on
  `Close` before releasing run storage, and tests wait before closing their
  store. Async delivery still never blocks the request that triggered it.
- Subscription item IDs are now stored when the item is created instead of
  being derived from its array position, so deleting an item no longer shifts
  the IDs of the items after it. The `si_<subscription>_<n>` shape is
  unchanged, subscriptions stored before this release keep exposing their
  position-derived IDs and are backfilled with those same values on their next
  write, and a later add reuses the lowest unused index.
- Renewal invoices now report `billing_reason: subscription_cycle`. They were
  labelled `subscription_create` because the renewal check looked for an
  `in_renewal_` ID prefix that renewals never had, leaving first charges and
  renewals indistinguishable to webhook consumers.
- Invoice previews without item overrides now return the subscription's next
  billing cycle instead of a zero-amount preview: `create_preview` and
  `upcoming` with just a subscription report the upcoming items, discounts, and
  tax with `billing_reason: upcoming`, include any deferred
  `create_prorations` amount without consuming it, and match the invoice the
  next renewal actually produces. Previews that pass item overrides keep the
  existing proration behaviour.
- Fixture packs can now seed the tax and discount evidence their subscriptions
  need: top-level `coupons` and `promotion_codes` keys, and
  `subscriptions[].default_tax_rates`. Tax-rate and coupon evidence is seeded
  before the subscription checkout runs and the resolved rates are passed into
  the checkout session, so a seeded subscription's first invoice is taxed and
  renewal/proration/preview inherit the snapshot. Subscription fixtures that
  reference a seeded coupon no longer need to repeat its percent or amount,
  and assertions accept `subtotal`/`tax`/`total` on the `invoice` target.
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
