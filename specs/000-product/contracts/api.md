# API Contract

Billtap exposes a Stripe-like API subset under `/v1` and Billtap-specific dashboard APIs under `/api`.
The public compatibility snapshot in `docs/COMPATIBILITY.md` is the release
source of truth for supported and unsupported provider behavior.

Webhook event shape, signature, retry, duplicate, delay, out-of-order, and replay behavior are defined in [webhooks.md](webhooks.md).

## Health

### `GET /healthz`

Process health.

### `GET /readyz`

Storage and worker readiness.

## Runs

One running server can host several isolated billing datasets. Every `/v1`
and `/api` request resolves a run before dispatch:

- A request with no selector uses the `default` run, backed by the
  configured `database_url`. This keeps existing integrations unchanged.
- A request may select a named run with the `/runs/<runId>` path prefix, for
  example `/runs/ci-123/v1/customers` or `/runs/ci-123/api/diagnostics`.
  Named runs are created on first use, have their own storage, and are isolated
  from each other and from `default`.
- The resolved run ID is returned on `X-Billtap-Run-Id`. The legacy
  `X-Billtap-Workspace` response header is also echoed for compatibility.
- The legacy `X-Billtap-Workspace` request header and `workspace` query
  parameter remain supported as aliases for unprefixed requests. An invalid
  run ID returns `400`.

### `GET /admin/runs`

Lists known runs (the default, any opened this session, and any whose database
file already exists). Returns a `list` envelope of `run` objects with `runId`,
`is_default`, `open`, `storage`, and table row-count `summary`.

### `DELETE /runs/<runId>`

Deletes a named run store. For `default`, clears user data while preserving
schema metadata.

### `GET /workspaces`

Legacy alias that lists run partitions as `workspace` objects with `name` and
`is_default`.

## Stripe-like API

### Customers

- `POST /v1/customers`
- `GET /v1/customers/{id}`
- `GET /v1/customers`
- `POST /v1/customers/{id}`

### Products

- `POST /v1/products`
- `GET /v1/products/{id}`
- `GET /v1/products`
- `GET /v1/products/search`

### Prices

- `POST /v1/prices`
- `GET /v1/prices/{id}`
- `GET /v1/prices`
- `GET /v1/prices/search`

Price search returns a Stripe-like `search_result` envelope. The supported
query subset covers `active:'true|false'`, `type:'one_time|recurring'`,
`lookup_key:'...'`, and `metadata['key']:'value'` clauses joined by `AND`.
This supports one-time price lookup paths while keeping unsupported query
clauses explicit validation errors.

### Coupons

- `POST /v1/coupons`
- `GET /v1/coupons/{id}`
- `GET /v1/coupons`
- `POST /v1/coupons/{id}`
- `DELETE /v1/coupons/{id}`

Coupons are local discount evidence. They support basic percent-off or
amount-off fields, metadata, and deletion markers. Billtap applies a bounded
single-discount subset to customer defaults, checkout sessions, subscriptions,
invoice previews, and renewal invoices.

Revised 2026-08-04 (product-scoped coupon adoption): coupons also persist
`applies_to[products]`, `duration_in_months`, `redeem_by`, `max_redemptions`,
and `times_redeemed` in their stripe-node v22 shapes. Product-scoped coupons
are enforced end to end: creation rejects sessions/subscriptions with no
matching line-item product, and discount math applies only to matching items
through completion, renewal, and preview. `redeem_by` expiry is enforced at
use time; `times_redeemed`/`max_redemptions` are recorded but not enforced.
Integer `percent_off` only (decimal percent remains outside the subset).

### Promotion Codes

- `POST /v1/promotion_codes`
- `GET /v1/promotion_codes/{id}`
- `GET /v1/promotion_codes`
- `POST /v1/promotion_codes/{id}`

Promotion codes are local coupon-linked evidence. They can be listed by
`code`, `coupon`, `customer`, and `active`, and can be applied through
`discounts[0][promotion_code]`. Redemption limits, expiration rules, and
promotion-code analytics are not modeled.

### Checkout Sessions

- `POST /v1/checkout/sessions`
- `GET /v1/checkout/sessions/{id}`
- `GET /v1/checkout/sessions`

Response includes:

- `id`
- `object`
- `url`
- `status`
- `payment_status`

Revised 2026-08-05 (consumer session metadata parity): session create accepts
session-level `metadata[...]` and every session response carries `metadata`
(an object; `{}` when unset). Session metadata is stored, survives retrieval
and completion, and is kept strictly independent of
`payment_intent_data[metadata]` — identical keys on both sides do not
overwrite each other, and session metadata is never promoted onto the
PaymentIntent created at completion (SDK callers that attach the same map to
both places get exactly what they sent in each place).

Revised 2026-08-04 (SaaS tax/discount adoption): responses also include
`currency`, `amount_subtotal`, `amount_total`, `total_details`
(`amount_discount`/`amount_shipping`/`amount_tax`), array-shaped `discounts`,
`automatic_tax` (`enabled`/`liability`/`provider`/`status`), and
`tax_id_collection` (`enabled`/`required`) in their stripe-node v22 shapes.
`automatic_tax[enabled]` and `tax_id_collection[enabled]` are accepted at
creation. Automatic tax is a deterministic local simulation: the rate
snapshots from customer metadata `tax_percent` at session creation and
applies exclusively after discounts (absent metadata means 0% with status
`complete`, matching Stripe's no-registration behavior). Address- or
jurisdiction-based calculation is not modeled.

Revised 2026-08-05 (consumer one-time payment adoption): `mode=payment` is
now supported alongside `mode=subscription` (`setup` still rejected — the
release-blocking invalid-mode case moved to `setup`). Payment-mode sessions
accept `line_items[i][price_data][...]` (`currency` required; one of
`product` / `product_data[name]`; one of `unit_amount` /
`unit_amount_decimal`, integer minor units only; `recurring` rejected in
payment mode), `payment_intent_data[setup_future_usage|description|
receipt_email|capture_method|metadata[...]]`, and `client_reference_id`
(now always serialized as `string | null` in both modes). `price_data`
creates real local Product/Price evidence. Completion creates no
subscription and no invoice: it creates one PaymentIntent for the
discounted (and taxed) total, sets `payment_status=paid` (free totals:
`no_payment_required`, no PaymentIntent), emits `checkout.session.completed`,
and exposes `description` / `receipt_email` / `setup_future_usage` as
top-level PaymentIntent fields (`string | null`); caller-supplied
`payment_intent_data[metadata]` round-trips without billtap key injection
(internal snapshots use `billtap_*` keys). `payment_intent_data` is rejected
in subscription mode and `subscription_data` in payment mode. The hosted
checkout page renders payment-mode sessions without subscription/invoice
rows and labels the plan slot "One-time payment" when no nickname exists.

Revised 2026-08-05 (hosted promotion-code entry): the hosted checkout page
shows an "Add promotion code" control on open `allow_promotion_codes`
sessions, backed by Billtap-specific
`POST/DELETE /api/checkout/sessions/{id}/promotion_code` (form
`promotion_code=<code>`; `promo_` IDs also accepted). Apply validates the
code through the same path as creation-time `discounts[0][promotion_code]`
(existence, active, coupon `redeem_by`, `applies_to[products]` line-item
match) and stores the discount on the open session so serialized totals and
subsequent completion reflect it; only one code may be applied at a time and
removal restores the original totals. Sessions created without
`allow_promotion_codes` reject application. `times_redeemed` remains
untouched, matching the creation path.

### Tax Rates and Customer Tax IDs

- `POST /v1/tax_rates`
- `GET /v1/tax_rates/{id}`
- `GET /v1/tax_rates`
- `POST /v1/tax_rates/{id}`
- `POST /v1/customers/{id}/tax_ids`
- `GET /v1/customers/{id}/tax_ids/{id}`
- `GET /v1/customers/{id}/tax_ids`
- `DELETE /v1/customers/{id}/tax_ids/{id}`

Added 2026-08-04: local evidence stores in stripe-node v22 shapes. Tax rates
were originally pure evidence and were not wired into totals (automatic tax
uses the customer-metadata simulation above). Tax IDs carry `verified` status
without provider verification. The Stripe Tax API family
(`/v1/tax/calculations|registrations|settings|transactions`) remains
unimplemented and returns `unsupported_endpoint`.

Revised 2026-08-04 (consumer default_tax_rates adoption): tax rates now apply
to real billing totals. Checkout sessions accept
`subscription_data[default_tax_rates][]`, and subscriptions accept
`default_tax_rates` on create and update (Emptyable: a single empty string
clears previously set rates), all as arrays of `txr_*` IDs resolved against
the local tax-rate evidence at creation time (`resource_missing` when
unknown). Rates snapshot onto the session/subscription/invoice
(`billtap_default_tax_rates` subscription metadata) and drive
inclusive/exclusive math after discounts through completion and renewal
invoices: inclusive amounts extract as `base × pct / (100 + Σ inclusive pct)`,
exclusive amounts add on the pretax base. Subscriptions serialize
`default_tax_rates` as full TaxRate objects; invoices populate
`default_tax_rates` and per-rate `total_taxes`/`total_tax_amounts` entries
with real `txr_*` IDs and `taxability_reason: null`. `automatic_tax[enabled]`
and `default_tax_rates` are mutually exclusive (400 `parameter_invalid`).
Invoice previews did not apply either tax path at the time of this revision.

Revised 2026-08-05 (consumer preview parity): invoice previews now apply tax.
`POST /v1/invoices/create_preview` and `GET,POST /v1/invoices/upcoming` read
the subscription's `default_tax_rates` snapshot (or the `automatic_tax`
simulation) and tax the post-discount proration base through the same helper
the confirmed invoice serialization uses, so a preview and the invoice
produced by actually applying the same change are identical field for field —
`subtotal`, `tax`, `total`, `total_excluding_tax`, and per-rate `total_taxes`
/ `total_tax_amounts` — including decimal-rate rounding and inclusive rates.
Previews for a customer with no subscription keep `tax: null` and empty tax
arrays. Preview item parsing also accepts `[price_id]` alongside `[price]`
(previously the validator allowed `price_id` while the parser ignored it,
silently yielding a zero proration), and a preview that cannot prorate
reports why in `billtap_preview.proration_skipped_reason` instead of
returning a silent zero.

Revised 2026-08-05 (upcoming-invoice parity): a preview that names a
subscription but overrides no items now returns that subscription's **next
billing cycle** instead of a zero-amount proration. The response carries one
non-proration line per subscription item, the next period
(`current_period_end` onward, or `trial_end` for a trialing subscription),
discounts and tax over the item total, and `billing_reason: upcoming`. Any
amount deferred by an earlier `create_prorations` update is included without
being consumed, so the preview equals the invoice the next renewal produces.
Previews that do override items keep the proration behaviour described above.

### Subscriptions

- `GET /v1/subscriptions/{id}`
- `GET /v1/subscriptions`
- `POST /v1/subscriptions/{id}`
- `DELETE /v1/subscriptions/{id}`
- `GET /v1/customers/{id}/discount`
- `DELETE /v1/customers/{id}/discount`
- `GET /v1/subscriptions/{id}/discount`
- `DELETE /v1/subscriptions/{id}/discount`

Revised 2026-08-05 (consumer plan-change adoption): item changes on
`POST /v1/subscriptions/{id}` now bill, instead of only recording
`proration_behavior` as metadata evidence.

- `proration_behavior=always_invoice` issues a paid `subscription_update`
  invoice immediately and repoints `latest_invoice` at it. Mid-cycle the
  invoice bills the prorated delta between the old and new item totals,
  scaled by the unused fraction of the period (truncating integer division,
  the same calculator the invoice preview uses, so preview and actual agree).
  A non-positive delta (downgrade) issues no invoice — credit balances are not
  modeled.
- `billing_cycle_anchor=now` additionally resets the period to
  `now .. now + interval` and bills the new full cycle net of the unused
  old-cycle credit. The credit is netted out of `subtotal` (and recorded as
  `billtap_proration_credit` invoice metadata) so every proration invoice
  satisfies `total == subtotal - discounts + tax`, and the serialized
  `total_taxes` recomputation agrees with the stored `tax`.
- `proration_behavior=create_prorations` issues no invoice; the delta
  accumulates in `billtap_pending_proration_amount` subscription metadata and
  is added to the next renewal invoice's subtotal before discounts and tax.
- `proration_behavior=none` (the default) keeps the previous
  items-and-metadata-only behavior.
- Tax follows the same rule as completion and renewal: the subscription's
  `default_tax_rates` (or the `automatic_tax` simulation) apply exclusively to
  the post-discount base.
- `payment_behavior=error_if_incomplete` with a configured failing outcome
  returns HTTP 402 in Stripe's card-error shape and leaves the subscription
  items unchanged (nothing is committed). Other `payment_behavior` values
  leave the invoice `open` and still apply the item change, matching the
  renewal-failure path.
- Invoicing paths emit `invoice.created`, `invoice.finalized`,
  `payment_intent.created`, the PaymentIntent terminal event,
  `invoice.payment_succeeded`/`invoice.paid` (or `invoice.payment_failed`),
  then `customer.subscription.updated`.
- `billing_cycle_anchor` accepts `now`, `unchanged`, or a Unix timestamp.
- `POST /v1/subscription_items` and `DELETE /v1/subscription_items/{id}` accept
  `proration_behavior` and `proration_date` and route through the same
  proration path, so adding seats with `always_invoice` bills the prorated
  delta under the subscription's tax rates. Delete reads its parameters from
  the query string as well as the body, rejects deleting the last item, and
  accepts `clear_usage` as evidence only. Item-level `tax_rates` round-trip as
  evidence but do not affect totals — billtap taxes at the subscription level.
  Subscription item IDs remain index-derived (`si_<subscription>_<index>`), so
  deleting a middle item shifts the IDs of later items.

### Subscription Schedules

- `POST /v1/subscription_schedules`
- `GET /v1/subscription_schedules/{id}`
- `GET /v1/subscription_schedules`
- `POST /v1/subscription_schedules/{id}`
- `POST /v1/subscription_schedules/{id}/cancel`
- `POST /v1/subscription_schedules/{id}/release`

Subscription schedules are local one-phase schedule evidence for existing
subscriptions. Test-clock advance applies a due phase by replacing subscription
items and emitting `customer.subscription.updated`. Multi-phase proration,
invoicing, and full schedule lifecycle parity are not modeled.

### Invoices

- `POST /v1/invoices`
- `GET /v1/invoices/{id}`
- `GET /v1/invoices`
- `POST /v1/invoices/{id}/finalize`
- `POST /v1/invoices/{id}/pay`
- `GET /v1/invoices/{id}/lines`
- `GET /v1/invoices/{id}/payments`
- `POST /v1/invoiceitems`
- `GET /v1/invoiceitems`
- `GET /v1/invoiceitems/{id}`
- `POST /v1/invoices/create_preview`
- `GET /v1/invoices/upcoming`
- `POST /v1/invoices/upcoming`

Direct invoice `pay` is a local retry mutation for open invoices created by
Billtap checkout and scenarios. It accepts deterministic sandbox
`payment_method` or legacy `source` aliases plus bounded protocol flags such as
`paid_out_of_band`, `forgive`, `off_session`, and `mandate`.

Manual one-time invoices support a bounded local flow:

1. `POST /v1/invoices` creates a draft invoice for a customer and accepts
   Stripe SDK-style fields including `currency`, `collection_method`,
   `default_payment_method`, `description`, `auto_advance`,
   `pending_invoice_items_behavior`, `payment_settings[payment_method_types]`,
   and `metadata[...]`.
2. `POST /v1/invoiceitems` attaches a one-time amount and preserves
   `metadata[...]`.
3. `POST /v1/invoices/{id}/finalize` opens the invoice and creates local
   PaymentIntent evidence.
4. `POST /v1/invoices/{id}/pay` applies deterministic PaymentIntent outcomes,
   including customer-level `default_payment_intent_outcome` fixture metadata.

Invoice responses include local fallback fields for SDK adoption paths:
`hosted_invoice_url`, `invoice_pdf`, `confirmation_secret.client_secret`, and
`payments.data.payment.payment_intent` with id, status, client secret, metadata,
and `last_payment_error` where applicable. Full invoice rendering, send, void,
line mutation, automatic collection, tax, and dunning automation are not part of
the current release-compatible subset. (Original boundary retained for
reference; partially revised 2026-08-04 below.)

Revised 2026-08-04 (SaaS tax/email adoption): two of those boundaries moved.
Invoices now carry the automatic-tax simulation (`automatic_tax` in its v22
shape, `tax`, `total_taxes` alongside legacy `total_tax_amounts`, and
tax-exclusive `total_excluding_tax`/`subtotal_excluding_tax`) computed from
the checkout/subscription snapshot. `POST /v1/invoices/{id}/send` records
local email evidence (`billtap_email_sent_at`/`billtap_email_recipient`
metadata, an `invoice.sent` event, and a timeline entry) and finalizing a
`send_invoice` invoice emits the same evidence in finalized → sent order; no
real email is delivered. `hosted_invoice_url` redirects to the invoice API
resource. Rendering, PDF, void, line mutation, automatic collection, and
dunning remain outside the subset.

Preview endpoints accept Stripe SDK-style `subscription`,
`subscription_details[items][0][price]`,
`subscription_details[items][0][quantity]`,
`subscription_details[proration_date]`, and
`subscription_details[proration_behavior]`, and
`subscription_details[billing_cycle_anchor]`. Billtap calculates a local
subscription-update proration line from the current period bounds and old/new
price totals. The bounded single-discount subset updates preview `subtotal`,
`total`, and `total_discount_amounts`. Preview responses include the common
Stripe Invoice defaults needed by generated SDK models, including
`attempt_count`, `attempted`, `auto_advance`, `automatic_tax`,
`billing_reason`, `collection_method`, `metadata`, `paid`, `payment_settings`,
`period_start`, `period_end`, `status_transitions`, `subtotal_excluding_tax`,
`total_excluding_tax`, `total_tax_amounts`, and array-shaped `discounts`.
Taxes, pending invoice items, and collection behavior are outside the modeled
subset. (Revised 2026-08-04: realized invoices now model the automatic-tax
simulation described above; previews themselves remain tax-free.)

### Payment Intents

- `GET /v1/payment_intents/{id}`
- `GET /v1/payment_intents`
- `POST /v1/payment_intents`
- `POST /v1/payment_intents/{id}/confirm`
- `POST /v1/payment_intents/{id}/capture`
- `POST /v1/payment_intents/{id}/cancel`

Direct payment intents are local state-machine simulations. They support
deterministic sandbox aliases, manual capture, cancel, timeline evidence, and
local webhook events. `requires_action` returns a local
`next_action.use_stripe_sdk` shape that can be completed or canceled through
Billtap action callbacks, and local bank-transfer intents can move from
`processing` to `succeeded` when customer cash balance is funded. One-time
PaymentIntents can store a deferred per-intent outcome with
`metadata[billtap_payment_intent_outcome]`, `billtap_outcome`, or
`deferred_outcome`; the stored outcome is applied when the intent is confirmed.
If no per-intent outcome is present, customer metadata
`billtap_default_payment_intent_outcome` or `default_payment_intent_outcome`
can provide the default outcome for direct and invoice-backed one-time
PaymentIntents.

### Setup Intents

- `GET /v1/setup_intents/{id}`
- `GET /v1/setup_intents`
- `POST /v1/setup_intents`
- `POST /v1/setup_intents/{id}/confirm`
- `POST /v1/setup_intents/{id}/cancel`

Setup intents are local state-machine simulations for saved-payment-method
smoke tests. Mandates, bank verification, and full SCA behavior are not part of
the current release-compatible subset.

### Payment Methods

- `GET /v1/payment_methods?customer={id}&type=card`

Returns deterministic sandbox card projections for known customers, including
the saved default payment method set by local portal payment-method simulation.
The response includes Stripe-like SDK fields such as `billing_details`,
`card.checks`, `card.networks`, `card.three_d_secure_usage`, `metadata`, and
`redaction`. Query validation covers `type`, `allow_redisplay`, `limit`, and
unknown parameters; valid non-card types return an empty local list.

### Customer Cash Balance

- `GET /v1/customers/{id}/cash_balance`
- `POST /v1/customers/{id}/cash_balance`
- `GET /v1/customers/{id}/cash_balance_transactions`
- `GET /v1/customers/{id}/cash_balance_transactions/{id}`
- `POST /v1/test_helpers/customers/{id}/fund_cash_balance`

Cash-balance APIs are local evidence for bank-transfer smoke tests. The
test-helper funding endpoint records a cash-balance transaction and settles
processing bank-transfer PaymentIntents for the customer.

### Refunds

- `POST /v1/refunds`
- `GET /v1/refunds/{id}`
- `GET /v1/refunds`
- `POST /v1/refunds/{id}`
- `POST /v1/refunds/{id}/cancel`

Refunds are local payment-history evidence. Creation accepts `charge`,
`payment_intent`, or `invoice`, plus `amount`, optional `reason`, and metadata.
It emits local `charge.refunded` and `charge.refund.updated` events. A fixture
or API call can keep a refund `pending`; when attached to a test clock with
`settle_at`/`available_on`, clock advance marks it `succeeded` and emits
`charge.refund.updated`. Balance transactions and processor settlement are
outside the modeled subset.

### Credit Notes

- `POST /v1/credit_notes`
- `GET /v1/credit_notes/{id}`
- `GET /v1/credit_notes`
- `POST /v1/credit_notes/{id}/void`

Credit notes are local invoice-history evidence. Creation accepts `invoice`,
`amount`, optional `reason`, and metadata. It emits `credit_note.created`; void
emits `credit_note.voided`. Line-level tax/discount accounting, PDF rendering,
and customer-balance math are outside the modeled subset.

### Disputes

- `GET /v1/disputes`
- `GET /v1/disputes/{id}`
- `POST /v1/disputes/{id}`
- `POST /v1/disputes/{id}/close`
- `GET /v1/charges/{id}/dispute`
- `POST /v1/charges/{id}/dispute`

Disputes are local chargeback-style evidence. Creating one emits
`charge.dispute.created`; updating evidence emits `charge.dispute.updated`;
fixture statuses can also emit `charge.dispute.funds_withdrawn`; closing one
emits `charge.dispute.closed`. Representment deadlines, balance movement, and
processor outcomes are outside the modeled subset.

### Test Clocks

- `POST /v1/test_helpers/test_clocks`
- `GET /v1/test_helpers/test_clocks/{id}`
- `GET /v1/test_helpers/test_clocks`
- `POST /v1/test_helpers/test_clocks/{id}/advance`

Test clocks are persisted local clocks for deterministic lifecycle simulation.
Customers, subscriptions, and pending refunds can be attached through
`test_clock`. Advancing a clock processes due trials, renewals, configured
renewal failures, period-end cancellations, and refund settlement for attached
objects.

### Billing Portal Sessions

- `POST /v1/billing_portal/sessions`

Returns a Stripe-like `billing_portal.session` object and Billtap hosted portal
URL for a known customer. The request accepts `customer`, `return_url`, optional
`configuration`, `locale`, `on_behalf_of`, and `flow_data`. The response
includes explicit `flow`, `locale`, `on_behalf_of`, `livemode`, and
`return_url` fields. Flow enum values and required nested flow fields are
validated with Stripe-style error envelopes. Hosted portal actions can save a
deterministic payment method, cancel a subscription, emit the matching local
webhook evidence, and redirect to `return_url`.

### Connect Platform Evidence

- `POST /v1/accounts`
- `GET /v1/account`
- `GET /v1/accounts/{id}`
- `GET /v1/accounts`
- `POST /v1/accounts/{id}`
- `DELETE /v1/accounts/{id}`
- `POST /v1/accounts/{id}/reject`
- `POST /v1/account_links`
- `POST /v1/account_sessions`
- `POST /v1/accounts/{id}/login_links`
- `GET /v1/accounts/{id}/capabilities`
- `GET /v1/accounts/{id}/capabilities/{capability}`
- `POST /v1/accounts/{id}/capabilities/{capability}`
- `POST /v1/accounts/{id}/people`
- `GET /v1/accounts/{id}/people`
- `GET /v1/accounts/{id}/people/{person}`
- `POST /v1/accounts/{id}/people/{person}`
- `DELETE /v1/accounts/{id}/people/{person}`
- `POST /v1/accounts/{id}/persons`
- `GET /v1/accounts/{id}/persons`
- `GET /v1/accounts/{id}/persons/{person}`
- `POST /v1/accounts/{id}/persons/{person}`
- `DELETE /v1/accounts/{id}/persons/{person}`
- `POST /v1/accounts/{id}/external_accounts`
- `GET /v1/accounts/{id}/external_accounts`
- `GET /v1/accounts/{id}/external_accounts/{external_account}`
- `POST /v1/accounts/{id}/external_accounts/{external_account}`
- `DELETE /v1/accounts/{id}/external_accounts/{external_account}`
- `POST /v1/accounts/{id}/bank_accounts`
- `GET /v1/accounts/{id}/bank_accounts/{bank_account}`
- `POST /v1/accounts/{id}/bank_accounts/{bank_account}`
- `DELETE /v1/accounts/{id}/bank_accounts/{bank_account}`
- `POST /v1/transfers`
- `GET /v1/transfers`
- `GET /v1/transfers/{id}`
- `POST /v1/transfers/{id}`
- `POST /v1/transfers/{id}/reversals`
- `GET /v1/transfers/{id}/reversals`
- `GET /v1/transfers/{id}/reversals/{reversal}`
- `POST /v1/transfers/{id}/reversals/{reversal}`
- `POST /v1/payouts`
- `GET /v1/payouts`
- `GET /v1/payouts/{id}`
- `POST /v1/payouts/{id}`
- `POST /v1/payouts/{id}/cancel`
- `POST /v1/payouts/{id}/reverse`
- `GET /v1/application_fees`
- `GET /v1/application_fees/{id}`
- `POST /v1/application_fees/{id}/refund`
- `POST /v1/application_fees/{id}/refunds`
- `GET /v1/application_fees/{id}/refunds`
- `GET /v1/application_fees/{id}/refunds/{refund}`
- `POST /v1/application_fees/{id}/refunds/{refund}`

Connect APIs are local smoke-test fixtures for platform-style routing. Platform
account retrieval returns deterministic local evidence. Account
create/list/retrieve/update persist connected-account profiles with metadata and
basic capability status, while account deletion returns a local marker rather
than closing a provider account. People/persons persist local representative
evidence. Account links, account sessions, and login links return local hosted
URLs/client secrets for onboarding or embedded-component tests. External
accounts, bank accounts, transfers, reversals, payouts, application fees, and
fee refunds are local evidence objects for platform-style integration tests.
Billtap records `Stripe-Account` request headers in redacted request traces and
can emit local Connect evidence webhooks, but it does not perform real
onboarding, KYC, identity verification, external-account verification, balance
movement, account closure, or settlement.

### Webhook Endpoints

- `POST /v1/webhook_endpoints`
- `GET /v1/webhook_endpoints/{id}`
- `GET /v1/webhook_endpoints`
- `POST /v1/webhook_endpoints/{id}`
- `PATCH /v1/webhook_endpoints/{id}`
- `DELETE /v1/webhook_endpoints/{id}`
- `GET /v1/webhook_endpoints/{id}/attempts`

### Events

- `GET /v1/events/{id}`
- `GET /v1/events`

Event list filters include `type`, `scenarioRunId`, created-time ranges,
`data.object.customer`, and `data.object.metadata[key]`.

## Hosted UI

### `GET /checkout/{sessionId}`

Hosted sandbox checkout page.

### `GET /portal/{customerId}`

Hosted sandbox billing portal page.

## Dashboard API

### `GET /api/timeline`

Filters:

- customerId
- subscriptionId
- checkoutSessionId
- invoiceId
- paymentIntentId
- eventId
- scenarioRunId

### `GET /api/delivery-attempts`

Webhook delivery attempts. Response evidence masks endpoint credentials,
sensitive headers, sensitive request URL query parameters, and webhook
signature HMAC values.

### `POST /api/events/replay-group`

Replays multiple existing event IDs with ordered, out-of-order, delayed,
signature-mismatch, duplicate, or omitted delivery evidence.

### `POST /api/payment_intents/{id}/complete_action`

Completes a local `requires_action` PaymentIntent.

### `POST /api/payment_intents/{id}/cancel_action`

Cancels a local `requires_action` PaymentIntent.

### `POST /api/payment_intents/{id}/outcome`

Stores a local deferred PaymentIntent outcome before confirmation.

Request:

```json
{
  "outcome": "requires_action"
}
```

The outcome uses the same deterministic aliases accepted by direct
PaymentIntent confirmation, including `payment_succeeded`, `card_declined`,
`requires_action`, `payment_pending`, `bank_transfer`, and `canceled`.

### `POST /api/disputes`

Creates local dispute evidence when a test does not already have a charge-like
ID to use through `/v1/charges/{id}/dispute`.

### `POST /api/events/{id}/replay`

Replay a webhook event. Records `webhook.replay` in the audit log and returns
redacted delivery attempt evidence.

Replay accepts reliability controls such as duplicate delivery, out-of-order
delivery, signature mismatch, forced response status, and
`simulate_app_failure` with `status`, `fail_first_n_attempts`, and optional
`body`. Simulated app failures record failed delivery attempts without calling
the app endpoint for the injected failures, then deliver the real replay
attempt after the configured failures are exhausted.

### `POST /api/webhooks/endpoints/{id}/replay-historical`

Replay historical events to one webhook endpoint. This is a Billtap-specific
catchup API for fixture or startup flows where events already exist before the
application registers its webhook endpoint. It records
`webhook.replay_historical` in the audit log and returns redacted delivery
attempt evidence.

Query or form fields:

- `since`: optional RFC3339 timestamp, Unix timestamp, or `now`
- `until`: optional RFC3339 timestamp, Unix timestamp, or `now`; defaults to
  the endpoint creation time
- `type` / `types` / `event_type` / `event_types`: optional event type filters
  such as `invoice.paid` or `invoice.*`
- `limit`: optional positive replay count
- `force`: optional boolean; when false, events with existing delivery attempts
  for the endpoint are skipped

Historical replay respects the endpoint's `enabled_events` filters, keeps the
original event ID and payload, marks delivery attempts with replay and
historical metadata, and does not replay events created after the endpoint
registration time unless an explicit `until` is provided.

### `POST /api/debug-bundles`

Create a debug bundle.

### `POST /api/fixtures/apply`

Apply a developer-test fixture pack. Request body may be JSON or YAML.

### `POST /api/fixtures/validate`

Validate a fixture pack without mutating local billing state. Request body may
be JSON or YAML. The response contains `valid: true` plus counts for supported
sections when schema and local semantic checks pass.

Supported fixture sections:

- `customers`
- `connected_accounts`
- `catalog.products`
- `catalog.prices`
- `test_clocks`
- `subscriptions`
- `refunds`
- `credit_notes`
- `disputes`
- `assertions`

Customer fixtures can opt out of Billtap's default sandbox card projection with
`payment_methods_fixture: empty` or an explicit empty `payment_methods: []`
field. When either field is present, `GET /v1/payment_methods?customer=...`
returns an empty list until a portal payment-method update or explicit fixture
payment method adds a saved method. A `payment_methods` list may also provide
objects with `id` and optional `default: true` for multi-card local scenarios.

Billtap tags created objects with fixture metadata:

- `billtap_fixture_name`
- `billtap_fixture_run_id`
- `billtap_fixture_namespace`
- `billtap_fixture_ref`

Subscription fixtures are created through the same checkout-completion path as
normal billing flows, so subscriptions, invoices, payment intents, checkout
sessions, and timeline evidence remain consistent.
When the HTTP fixture apply API is used, Billtap also creates seeded webhook
events for the local checkout, subscription, invoice, payment-intent, refund,
and credit-note evidence so tests can list and replay those events through the
same `/v1/events` and `/api/events/{id}/replay` paths.
If a subscription already exists when the fixture is applied, Billtap backfills
a missing `customer.subscription.created` event for that subscription before
emitting fixture state updates. This keeps re-applied or pre-seeded fixture
graphs replayable without requiring a dedicated checkout flow.

Fixture-provided IDs are preserved for seeded objects. Fixtures also tag
objects with `billtap_fixture_ref`, and the resolve endpoint below can map a
fixture ref to the generated or stable customer, subscription, invoice, payment
intent, checkout session, product, and price IDs.

Revised 2026-08-05 (consumer seeding ergonomics): packs can declare the tax and
discount evidence their subscriptions need, so an E2E seed no longer needs
side-channel API calls.

- Top-level `coupons` and `promotion_codes` seed local coupon and
  promotion-code evidence in the same shapes `POST /v1/coupons` and
  `POST /v1/promotion_codes` produce. A promotion code whose coupon is neither
  in the pack nor already seeded fails the apply.
- `subscriptions[].default_tax_rates` attaches seeded `txr_*` rates to a
  subscription. Tax-rate, coupon, and promotion-code evidence is seeded
  *before* the subscription checkout runs and the resolved rates are passed
  into the checkout session, so the subscription's **first** invoice is taxed
  like any other billing flow (no post-hoc metadata patching), and renewal,
  proration, and preview inherit the same snapshot. Unknown rate IDs fail the
  apply, naming the subscription and rate.
- A subscription fixture that references a seeded coupon no longer has to
  repeat its `discount_percent_off` / `discount_amount_off`: the values are
  read from the coupon evidence, and an explicit value still wins.
- Assertions accept `subtotal`, `tax`, and `total` on the `invoice` target
  (rejected on other targets), so a seed can prove its own tax math.

### `GET /api/fixtures/resolve`

Resolve fixture-backed objects. Query fields:

- `ref`
- `id`
- `lookup_key` / `lookupKey`
- `runId`
- `fixture` / `fixtureName`
- `namespace`

Returns the matching local IDs for the seeded object graph.

### `GET /api/fixtures/snapshot`

Return a filtered billing snapshot for fixture-driven tests. Query fields:

- `customer` / `customerId`
- `runId`
- `tenantId`
- `fixture` / `fixtureName`
- `namespace`

Response includes customers, products, prices, checkout sessions,
subscriptions, invoices, payment intents, timeline entries, and object counts.

### `POST /api/fixtures/assert`

Assert expected fixture state. Request body may be JSON or YAML.

Supported assertion targets:

- `customer`
- `product`
- `price`
- `checkout_session`
- `subscription`
- `invoice`
- `payment_intent`
- `timeline`

Returns `200` when all assertions pass and `409` with a structured assertion
report when any assertion fails.

### `GET /api/audit-log`

List audit log entries. Query fields:

- `action`
- `targetId`

### `POST /api/retention/apply`

Apply the configured retention policy. Old webhook raw payloads and delivery
request/response bodies are redacted while IDs, statuses, timestamps, metadata,
and audit records are preserved.

### `GET /api/portal`

Load portal state. Accepts `customer_id` or `customerId`.

### `GET /api/portal/customers/{id}`

Load portal state for a customer.

### `POST /api/portal/subscriptions/{id}/plan-change`

Apply a sandbox plan change. Body fields:

- `plan`
- `price`
- `quantity`

### `POST /api/portal/subscriptions/{id}/seat-change`

Apply a sandbox seat quantity change. Body fields:

- `quantity`

### `POST /api/portal/subscriptions/{id}/cancel`

Cancel a sandbox subscription. Body fields:

- `mode`: `period` or `immediate`

### `POST /api/portal/subscriptions/{id}/resume`

Resume a subscription from pending or immediate cancellation state.

### `POST /api/portal/customers/{id}/payment-method`

Simulate payment method update. Body fields:

- `outcome`: `succeeds` or `fails`

### `POST /api/scenarios/run`

Run a scenario.

Request body may be a scenario JSON object or YAML content.

Response is a scenario report with:

- `name`
- `status`
- `failure_type`
- `clock_start`
- `clock_end`
- `steps`
- `errors`
