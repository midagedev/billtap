# Production Boundary Runbook

## Purpose

Use this runbook when Billtap is connected to staging or production-adjacent webhook traffic for debugging, replay validation, or support evidence. Billtap must not be placed in the live payment success path.

## Safe Relay Configuration

```bash
export BILLTAP_ENV=staging
export BILLTAP_RELAY_MODE=true
export BILLTAP_RAW_PAYLOAD_STORAGE=metadata_only
export BILLTAP_RETENTION_DAYS=7
```

Relay mode forces metadata-only raw payload storage even when another value is supplied for `BILLTAP_RAW_PAYLOAD_STORAGE`.

## Startup Check

1. Start Billtap with relay mode enabled.
2. Confirm `/readyz` returns `storage: ok`.
3. Create webhook endpoints only for controlled testmode or staging destinations.
4. Confirm `/v1/webhook_endpoints` returns `secret: "****"`.
5. Trigger a controlled testmode event.
6. Confirm `/api/delivery-attempts` masks signatures, credentials, and sensitive URL query values.
7. Confirm `/api/audit-log` records any replay or delivery override action.

## Retention

Run retention after debugging sessions:

```bash
curl -fsS -X POST http://127.0.0.1:8080/api/retention/apply
```

Retention redacts old webhook raw payloads and delivery request/response bodies while preserving IDs, event types, timestamps, statuses, and audit metadata.

## Hard Stops

- Do not submit real card numbers, CVC, or live payment credentials to Billtap.
- Do not make Billtap required for customer payment success.
- Do not run relay mode without endpoint access controls in front of the dashboard.
