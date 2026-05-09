# Release Checklist

## Required Checks

Most checks are automated by the `CI / Release gate` GitHub Actions workflow and
should be configured as required for PRs and protected merges. GHCR pull/run
checks are post-publish release checks and require the `Container Image`
workflow to complete first.

- `LICENSE` exists and matches the project owner's intended public license
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
- `docker pull ghcr.io/midagedev/billtap:main`
- `docker run --rm -d --name billtap-release-smoke -p 18080:8080 ghcr.io/midagedev/billtap:main`
- `curl -fsS http://127.0.0.1:18080/healthz`
- `docker rm -f billtap-release-smoke`
- Start the sample app with `PORT=3300 npm --prefix examples/sample-app start`
  and wait for `http://127.0.0.1:3300/healthz`
- `/tmp/billtap scenario run examples/subscription-payment-retry.yml`
- `curl -fsS -X POST http://127.0.0.1:3300/test/reset`
- `/tmp/billtap scenario run examples/saas-adoption-contract.yml`

## Manual Smoke

- Start Billtap from the local binary with `BILLTAP_STATIC_DIR=dist/app`.
- Start the sample app with `PORT=3300 npm --prefix examples/sample-app start` before CLI scenario smoke checks.
- Open `/app/dashboard/`, `/app/checkout/`, and `/app/portal/`.
- Create customer, product, price, and checkout session through `/v1`.
- Complete checkout through `/api/checkout/sessions/{id}/complete`.
- Verify `/api/timeline`, `/api/objects`, and `/api/delivery-attempts`.
- Run relay mode with `BILLTAP_RELAY_MODE=true` and confirm endpoint secrets, signatures, and sensitive URL query values are masked.
- Confirm the GHCR image tag used by downstream compose overlays is available
  before updating adoption repos.

## Release Notes Must State

- Billtap is not a real payment processor.
- Billtap must not be placed in a live payment success path.
- Relay mode is for controlled testmode or staging traffic only.
- Dashboard access controls are required before production-adjacent deployment outside a trusted local network.
