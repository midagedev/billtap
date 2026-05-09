# Scenario Contract

Billtap scenarios are YAML files for deterministic billing tests.

## Example

```yaml
name: subscription-payment-retry
app:
  webhookUrl: http://localhost:3000/webhooks/stripe
  assertions:
    baseUrl: http://localhost:3000/test/assertions

catalog:
  products:
    - id: prod_pro
      name: Pro
  prices:
    - id: price_pro_monthly
      product: prod_pro
      currency: usd
      unitAmount: 4900
      interval: month

steps:
  - id: create-customer
    action: customer.create
    params:
      email: user@example.test

  - id: checkout
    action: checkout.create
    params:
      customerRef: create-customer.customer.id
      price: price_pro_monthly

  - id: complete-checkout
    action: checkout.complete
    params:
      sessionRef: checkout.session.id
      outcome: payment_failed

  - id: assert-grace
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: past_due

  - id: advance-clock
    action: clock.advance
    params:
      duration: 3d

  - id: retry-payment
    action: invoice.retry
    params:
      subscriptionRef: checkout.subscription.id
      outcome: payment_succeeded

  - id: assert-active
    action: app.assert
    params:
      target: workspace.subscription
      expected:
        status: active
```

## Required Sections

- `name`
- `steps`

## Optional Sections

- `app`
- `catalog`
- `clock`
- `defaults`
- `profile`
- `webhooks`
- `saas`

Webhook delivery knobs are defined in [webhooks.md](webhooks.md).

`profile: saas` enables the deterministic SaaS service profile. The
optional `saas` section seeds tenant rail policy and catalog preset data:

```yaml
profile: saas
saas:
  tenant:
    id: tenant_direct
    rail: CARD
  catalogPreset: saas-default
```

## Step Result References

Steps may reference previous step outputs with:

```text
<step-id>.<path>
```

Numeric path segments address list entries, for example
`complete-checkout.events.0.id`.

## Built-In Actions

- `customer.create`
- `product.create`
- `price.create`
- `checkout.create`
- `checkout.complete`
- `checkout.cancel`
- `subscription.update`
- `subscription.cancel`
- `invoice.fail_payment`
- `invoice.retry`
- `clock.advance`
- `webhook.replay`
- `webhook.deliver_duplicate`
- `webhook.deliver_out_of_order`
- `app.assert`

`webhook.replay` schedules delivery attempts for an existing event id. It can
use `duplicate`, `delay` or `delay_seconds`, `outOfOrder`, `responseStatus`,
`responseBody`, `timeout`, `error`, and `signatureMismatch` parameters.

## SaaS Profile Actions

The representative adoption scenario is
`examples/saas-adoption-contract.yml`. It seeds a workspace, upgrades a
subscription, exercises cancellation and resume, covers seat purchase and member
limits, exports design cases, records extra export payment evidence, captures
payment history, emits platform and Connect webhook outcomes, builds a support
bundle, defines observability expectations, and runs app assertion callbacks.

- `saas.tenant.configure`
- `saas.workspace.create`
- `saas.workspace.activate`
- `saas.subscription.get_current`
- `saas.subscription.checkout_upgrade`
- `saas.subscription.preview_upgrade`
- `saas.subscription.confirm_upgrade`
- `saas.subscription.cancel`
- `saas.subscription.stop_pending_cancellation`
- `saas.seat.estimate_purchase`
- `saas.seat.purchase`
- `saas.member.invite`
- `saas.member.delete`
- `saas.export.summary`
- `saas.export.usage`
- `saas.export.products`
- `saas.export_session.create`
- `saas.export_session.files`
- `saas.extra_export.preview`
- `saas.extra_export.create_payment_intent`
- `saas.extra_export.provide`
- `saas.payment.customer_portal`
- `saas.payment.history`
- `saas.backoffice.start_subscription`
- `saas.backoffice.change_plan`
- `saas.backoffice.change_seat`
- `saas.backoffice.change_period`
- `saas.backoffice.refund`
- `saas.backoffice.update_export_limit`
- `saas.webhook.platform`
- `saas.webhook.connect`
- `saas.support.bundle`
- `saas.observability.expect`

## SaaS Assertion Targets

- `saas.workspace.subscription`
- `saas.workspace.plan`
- `saas.workspace.seats`
- `saas.workspace.members`
- `saas.workspace.export.summary`
- `saas.workspace.export.usage`
- `saas.extra_export.payment`
- `saas.payment.history`
- `saas.webhook.event`
- `saas.workspace.log`
- `saas.backoffice.refund`
- `saas.observability.signal`

## Exit Policy

- any failed `app.assert` fails the scenario
- invalid scenario config exits with code `2`
- app callback failure exits with code `3`
