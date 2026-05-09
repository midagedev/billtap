# Local Development Runbook

This runbook covers the local runtime shell and Docker image.

## Start Billtap

```bash
npm install
npm run build
go run ./cmd/billtap
```

Open:

```text
http://localhost:8080
```

Current app paths:

```text
http://localhost:8080/app/dashboard/
http://localhost:8080/app/checkout/
http://localhost:8080/app/portal/
```

Health checks:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

## Optional Config

```bash
BILLTAP_ADDR=127.0.0.1:18080 \
BILLTAP_DATABASE_URL=:memory: \
BILLTAP_ENV=smoke \
BILLTAP_STATIC_DIR=dist/app \
go run ./cmd/billtap
```

## Docker

```bash
docker build -t billtap:local .
docker run --rm -p 8080:8080 -v billtap-data:/data billtap:local
```

To use the published image instead of building locally:

```bash
docker run --rm -p 8080:8080 -v billtap-data:/data ghcr.io/midagedev/billtap:main
```

## Configure App Under Test

Set the app's Stripe base URL to:

```text
http://localhost:8080/v1
```

Set webhook endpoint in Billtap:

```text
http://host.docker.internal:3000/webhooks/stripe
```

## Run a Checkout Flow

```bash
customer_id=$(curl -fsS -X POST http://localhost:8080/v1/customers \
  -d email=buyer@example.test | jq -r .id)

product_id=$(curl -fsS -X POST http://localhost:8080/v1/products \
  -d name=Team | jq -r .id)

price_id=$(curl -fsS -X POST http://localhost:8080/v1/prices \
  -d product="$product_id" \
  -d currency=usd \
  -d unit_amount=9900 \
  -d 'recurring[interval]=month' | jq -r .id)

session_json=$(curl -fsS -X POST http://localhost:8080/v1/checkout/sessions \
  -d customer="$customer_id" \
  -d mode=subscription \
  -d "line_items[0][price]=$price_id" \
  -d success_url=http://localhost:3000/success \
  -d cancel_url=http://localhost:3000/cancel)

session_id=$(echo "$session_json" | jq -r .id)
echo "$session_json" | jq -r .url
```

Open the returned URL, or complete through the API:

```bash
curl -fsS -X POST "http://localhost:8080/api/checkout/sessions/$session_id/complete" \
  -H 'Content-Type: application/json' \
  -d '{"outcome":"payment_succeeded"}'
```

Inspect timeline:

```bash
curl -fsS "http://localhost:8080/api/timeline?checkoutSessionId=$session_id"
```

## Capture A Diagnostic Bundle

When an app under test is configured to use Billtap, capture the diagnostic
bundle before changing settings or restarting containers:

```bash
curl -fsS "http://localhost:8080/api/diagnostics?limit=200" \
  -o billtap-diagnostics.json
```

For a single object:

```bash
curl -fsS -X POST "http://localhost:8080/api/debug-bundles" \
  -d object_type=checkout_session \
  -d object_id="$session_id" \
  -o billtap-debug-$session_id.json
```

Use the bundle as the first artifact for failed local dev or isolated e2e runs:

- `request_traces`: proves whether the app called Billtap, which `/v1` path it
  used, what query/body was sent, the status Billtap returned, and the
  Stripe-like error code/param when validation failed
- `fixture_snapshot`: proves whether seed data was applied and whether expected
  fixture IDs are present
- `timeline`: explains checkout, subscription, invoice, and payment-intent state
  transitions
- `webhook_events` and `delivery_attempts`: prove whether Billtap emitted a
  webhook, where it sent it, what response the app returned, and whether a retry
  is pending

Common triage checks:

```bash
# Did the app route Stripe API calls to Billtap?
curl -fsS "http://localhost:8080/api/request-traces?limit=50" | jq '.data[].path'

# Did a specific checkout session have matching API and state evidence?
curl -fsS "http://localhost:8080/api/diagnostics?objectId=$session_id&limit=100"

# Did the webhook handler return success?
curl -fsS "http://localhost:8080/api/delivery-attempts" \
  | jq '.data[] | {event_id, request_url, status, response_status, error}'
```
