# Security Policy

## Supported Versions

Billtap is pre-release. No supported version line exists yet.

## Reporting A Vulnerability

Open a private security advisory in the public hosting provider when available. If advisories are not configured yet, contact the maintainer through the repository owner profile.

Do not include real credentials, real customer data, or production payment payloads in reports.

## Product Boundary

Billtap must not process real payments. It is intended for local development, CI, and controlled testmode or staging-adjacent debugging.

Relay mode must be treated as bounded and optional:

- enable with `BILLTAP_RELAY_MODE=true`
- raw payload storage is forced to `metadata_only`
- dashboard/API evidence should mask secrets and signature HMAC values
- real card-data fields are rejected

## Secret Handling

Billtap redacts sensitive headers, URL query values, endpoint secrets, and signature HMAC values in dashboard/API evidence. Keep new integrations on the same masking path.
