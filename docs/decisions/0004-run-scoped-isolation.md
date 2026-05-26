# 0004 Run-Scoped Isolation

## Status

Accepted

## Context

Parallel billing test suites need to share one Billtap server without leaking
Stripe-compatible objects, webhook endpoint registrations, idempotency state, or
test clocks across jobs. Existing integrations also need unprefixed `/v1/...`
requests to keep using the default dataset.

## Decision

Billtap exposes run-scoped routing at `/runs/<runId>/v1/...` and
`/runs/<runId>/api/...`. Unscoped requests use the `default` run.

The current implementation maps each named run to one isolated SQLite store.
The default run uses the configured `database_url`; named runs use sibling
SQLite databases under the existing `workspaces/` directory for on-disk
compatibility with earlier Billtap builds. This preserves duplicate Stripe
object IDs across runs and isolates webhook fan-out, test clocks, traces, and
fixture state without forcing a risky all-table primary-key migration in the
same change.

`GET /admin/runs` reports known runs and row-count summaries. `DELETE
/runs/<runId>` removes a named run store; `DELETE /runs/default` clears user data
from the default store while retaining schema metadata.

## Consequences

- Stripe SDK users can set the API base to `http://billtap:8080/runs/<runId>`
  and keep normal SDK paths.
- Hosted checkout and portal URLs generated from a run-scoped API request retain
  the `/runs/<runId>` prefix.
- Existing `X-Billtap-Workspace` and `workspace` query selectors remain
  supported as compatibility aliases.
- A future row-level `run_id` schema can be added if a single physical SQLite
  file becomes required, but the public isolation contract does not depend on
  that migration.
