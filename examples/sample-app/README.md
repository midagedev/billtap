# Billtap Sample App

This is a tiny local app for G9 release-candidate smoke checks. It has no npm
dependencies and uses Node's built-in HTTP server.

## Run

```bash
cd examples/sample-app
npm start
```

Useful environment variables:

```bash
PORT=3300
BILLTAP_WEBHOOK_SECRET=webhook_secret_local
BILLTAP_SIGNATURE_TOLERANCE_SECONDS=300
```

## Endpoints

- `POST /webhooks/stripe` receives Billtap/Stripe-style webhook events. If
  `BILLTAP_WEBHOOK_SECRET` is set, it verifies `Billtap-Signature` or
  `Stripe-Signature` with the `t=<unix>,v1=<hmac>` format.
- `POST /test/assertions/{target}` receives scenario-runner assertion payloads.
  The main target is `workspace.subscription`.
- `POST /test/reset` clears in-memory state.
- `GET /state` returns current in-memory state.
- `GET /healthz` returns a small health response.

The app stores state in memory only. Duplicate webhook event IDs are accepted
but not applied twice.

## Scenario Runner Smoke

From the repo root, start the app in one terminal:

```bash
cd examples/sample-app
npm start
```

Run the sample scenario in another terminal:

```bash
go run ./cmd/billtap scenario run examples/sample-app-checkout.yml
```

The scenario runner currently posts assertion context to the sample app. The
webhook endpoint is used when running Billtap as a server and registering this
app as a webhook endpoint.

## Webhook Smoke With Billtap Server

Start the sample app with a local signing secret:

```bash
cd examples/sample-app
BILLTAP_WEBHOOK_SECRET=webhook_secret_local npm start
```

Start Billtap in another terminal:

```bash
go run ./cmd/billtap
```

Register the sample app as a webhook endpoint:

```bash
curl -fsS -X POST http://localhost:8080/v1/webhook_endpoints \
  -d url=http://localhost:3300/webhooks/stripe \
  -d secret=webhook_secret_local \
  -d 'enabled_events[0]=*'
```

Then create and complete a checkout flow against Billtap. The sample app will
record received events, subscription state, and duplicate delivery evidence at
  `GET http://localhost:3300/state`.
