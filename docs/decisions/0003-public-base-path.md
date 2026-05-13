# Decision 0003: Browser-Facing Public Base Path

## Status

Accepted

## Context

Billtap may run behind a shared reverse proxy where the browser-visible URL is
path-prefixed, for example `https://localhost:8081/billtap`, while internal
container traffic still uses the service URL `http://billtap:8080`. Requiring
the proxy to rewrite HTML and JavaScript assets is brittle and makes local
multi-app stacks harder to reason about.

## Decision

Billtap supports a browser-facing public base path through:

- `PUBLIC_BASE_PATH`, a generic multi-app stack variable
- `BILLTAP_PUBLIC_BASE_PATH`, a Billtap-specific override
- `X-Forwarded-Prefix`, for proxies that strip the prefix before forwarding

The configured or forwarded path is used for:

- app links and static asset URLs
- dashboard API calls
- Stripe-like `/v1` browser calls
- redirect `Location` headers
- hosted checkout and portal URLs

`BILLTAP_PUBLIC_BASE_URL` remains the browser-visible origin. If both
`BILLTAP_PUBLIC_BASE_URL=https://localhost:8081` and
`PUBLIC_BASE_PATH=/billtap` are configured, generated hosted URLs use
`https://localhost:8081/billtap/...`.

Internal service-to-service traffic remains unprefixed and should continue to
use URLs such as `http://billtap:8080/v1`.

## Consequences

The Vite app is built with `/app/` under the public base path, so assets resolve
without proxy rewrites. Server routes accept both prefixed and unprefixed paths
when a configured public base path is present, which preserves internal compose
traffic while allowing browser requests under the prefix.

The public base path is validated as a URL path. Full URLs, query strings,
fragments, dot path segments, and empty path segments are rejected.
