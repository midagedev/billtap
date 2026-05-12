# Quickstart

Billtap can run as a local binary or Docker image.

## Build Locally

```bash
npm ci
npm run build
go build -o /tmp/billtap ./cmd/billtap
```

Start the server:

```bash
BILLTAP_ADDR=127.0.0.1:8080 \
BILLTAP_DATABASE_URL=.billtap/billtap.db \
BILLTAP_STATIC_DIR=dist/app \
/tmp/billtap
```

## Start With Docker

```bash
docker build -t billtap:local .
docker run --rm \
  -p 8080:8080 \
  -v billtap-data:/data \
  billtap:local
```

Open:

```text
http://localhost:8080
```

## Configure App

```bash
export STRIPE_API_BASE_URL=http://localhost:8080/v1
export STRIPE_WEBHOOK_SECRET=webhook_secret_local
```

## Create Checkout Session

```bash
customer_id=$(curl -fsS -X POST http://localhost:8080/v1/customers \
  -d email=buyer@example.test | jq -r .id)

product_id=$(curl -fsS -X POST http://localhost:8080/v1/products \
  -d name=Team | jq -r .id)

price_id=$(curl -fsS -X POST http://localhost:8080/v1/prices \
  -d product="$product_id" \
  -d currency=usd \
  -d unit_amount=9900 | jq -r .id)

curl -fsS -X POST http://localhost:8080/v1/checkout/sessions \
  -d customer="$customer_id" \
  -d "line_items[0][price]=$price_id" \
  -d success_url=http://localhost:3000/success \
  -d cancel_url=http://localhost:3000/cancel
```

Open the returned `url`.

## Run Scenario

Start the sample app assertion endpoint in one terminal:

```bash
PORT=3300 npm --prefix examples/sample-app start
```

Run the scenario in another terminal:

```bash
/tmp/billtap scenario run examples/subscription-payment-retry.yml \
  --report-json billtap-report.json \
  --report-md billtap-report.md
```

## Production Boundary Smoke

```bash
BILLTAP_RELAY_MODE=true \
BILLTAP_RAW_PAYLOAD_STORAGE=metadata_only \
BILLTAP_RETENTION_DAYS=7 \
/tmp/billtap
```

Relay mode masks endpoint secrets and delivery evidence, rejects real card fields, stores webhook metadata without raw payloads, and records replay audit entries.
