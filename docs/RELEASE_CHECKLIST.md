# Release Checklist

## Required Checks

These checks are automated by the `CI / Release gate` GitHub Actions workflow and should be configured as required for PRs and protected merges.

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
- `/tmp/billtap scenario run examples/subscription-payment-retry.yml`
- `/tmp/billtap scenario run examples/saas-adoption-contract.yml`

## Manual Smoke

- Start Billtap from the local binary with `BILLTAP_STATIC_DIR=dist/app`.
- Start the sample app with `PORT=3300 npm --prefix examples/sample-app start` before CLI scenario smoke checks.
- Open `/app/dashboard/`, `/app/checkout/`, and `/app/portal/`.
- Create customer, product, price, and checkout session through `/v1`.
- Complete checkout through `/api/checkout/sessions/{id}/complete`.
- Verify `/api/timeline`, `/api/objects`, and `/api/delivery-attempts`.
- Run relay mode with `BILLTAP_RELAY_MODE=true` and confirm endpoint secrets, signatures, and sensitive URL query values are masked.

## Release Notes Must State

- Billtap is not a real payment processor.
- Billtap must not be placed in a live payment success path.
- Relay mode is for controlled testmode or staging traffic only.
- Dashboard access controls are required before production-adjacent deployment outside a trusted local network.
