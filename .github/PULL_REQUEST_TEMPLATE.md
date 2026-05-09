## Summary

-

## Scope

- [ ] This change stays within Billtap's sandbox boundary and does not add real payment processing.
- [ ] This change does not introduce live credentials, real card data, production customer data, or production payment payloads.
- [ ] I updated documentation when behavior, compatibility, support, or security expectations changed.

## Verification

Run the checks that match the change:

- [ ] `go test ./...`
- [ ] `npm run typecheck`
- [ ] `npm run build`
- [ ] `npm run smoke:sample`
- [ ] `go run ./cmd/billtap scenario run examples/subscription-payment-retry.yml`
- [ ] `go run ./cmd/billtap scenario run examples/saas-adoption-contract.yml`
- [ ] Not applicable, documentation-only change

## Compatibility And Boundary Notes

If this changes Stripe-like behavior, webhook semantics, billing state, fixtures, scenario execution, relay mode, masking, retention, or dashboard evidence, describe the fixture/test coverage and any documented limitations.

## Screenshots Or Evidence

Add sanitized screenshots, reports, logs, or dashboard evidence when useful. Do not include secrets or real payment data.
