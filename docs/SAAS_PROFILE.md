# SaaS Workspace Profile

The public `saas` profile is a sanitized workspace-billing scenario pack for subscription apps.

It models common SaaS billing behaviors without naming a private product or depending on a private codebase:

- workspace owner and active workspace state
- tenant payment rail metadata
- plan tiers and payment cycles
- subscription upgrade, cancellation, resume, failure, and retry
- additional seat estimates and purchases
- member invite/delete flows
- included and manual export quota
- extra export payment and provision evidence
- customer portal evidence
- payment history and refund evidence
- platform/connect-style webhook evidence
- support bundle output
- app assertion targets

## Profile Usage

```yaml
name: saas-adoption-contract
profile: saas
saas:
  tenant:
    id: tenant_direct
    rail: CARD
  catalogPreset: saas-default
```

The representative scenario is:

```bash
go run ./cmd/billtap scenario run examples/saas-adoption-contract.yml
```

## Assertion Targets

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

## Boundary

This profile is a generic fixture-backed model for local and CI tests. It is not a contract for any private product, and it should not be treated as Stripe parity.
