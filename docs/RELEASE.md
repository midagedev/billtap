# Release Process

This document describes the public release procedure for Billtap `v0.1.0` and
later pre-1.0 releases.

Billtap is currently a source-first project. Unless release automation is added
in the same release, `v0.1.0` should be published as source code plus
maintainer-verified local build instructions. Do not claim a published package,
Homebrew formula, or official Docker image until that artifact exists.

This document does not choose a license.

## Release Goals For v0.1.0

The `v0.1.0` release should make three public claims:

- Billtap is a local and CI billing sandbox, not a real payment processor.
- Billtap supports a fixture-backed practical Stripe-like subset documented in
  `docs/COMPATIBILITY.md`.
- Billtap can be built, tested, run, and smoke-checked from a clean checkout.

## Required Public Statements

Every release note must state:

- Billtap must not be used in a live payment path.
- Billtap does not process real payments or accept real card data.
- Billtap is not full Stripe compatibility.
- Unsupported provider-specific behavior must be tested with Stripe testmode or
  the real provider sandbox.
- Relay mode is only for controlled testmode or staging-adjacent debugging.

## Release Branch And PR Flow

1. Create a release preparation branch from `main`.

   ```bash
   git switch main
   git pull --ff-only
   git switch -c release/v0.1.0
   ```

2. Update release notes.

   - Move relevant `CHANGELOG.md` entries under `## 0.1.0 - YYYY-MM-DD`.
   - Keep future work under `## Unreleased`.
   - Confirm `README.md`, `docs/README.md`, `docs/COMPATIBILITY.md`, and this
     document agree on release state.

3. Run local release verification.

4. Open a PR from `release/v0.1.0` to `main`.

   The PR must include:

   - Summary of release claim.
   - Verification command results.
   - Compatibility and boundary notes.
   - Any accepted deferrals.

5. Wait for CI and review.

   - Required CI must pass.
   - Review feedback must be addressed before merge.
   - Advisory AI review can be used as review assistance, but CI remains the
     release gate.

6. Merge the PR.

7. Tag from the merged `main`.

   ```bash
   git switch main
   git pull --ff-only
   git tag -a v0.1.0 -m "Billtap v0.1.0"
   git push origin v0.1.0
   ```

8. Create the GitHub release from the tag.

   The release description should link to:

   - `README.md`
   - `docs/COMPATIBILITY.md`
   - `docs/PRODUCTION_BOUNDARIES.md`
   - `docs/RELEASE_CHECKLIST.md`
   - `CHANGELOG.md`

## Local Release Verification

Run from the repository root:

```bash
npm ci
go test ./...
go run ./cmd/billtap compatibility scorecard --output-dir /tmp/billtap-compatibility
npm run typecheck
npm run build
npm run smoke:sample
npm run smoke:web:install
npm run smoke:web
go build -o /tmp/billtap ./cmd/billtap
docker build -t billtap:local .
```

Run scenario smoke with the sample app assertion endpoint:

```bash
PORT=3300 npm --prefix examples/sample-app start
```

In another terminal:

```bash
/tmp/billtap scenario run examples/subscription-payment-retry.yml
/tmp/billtap scenario run examples/saas-adoption-contract.yml
```

Stop the sample app after both scenarios pass.

## Manual Smoke

Start Billtap from the local binary:

```bash
BILLTAP_STATIC_DIR=dist/app /tmp/billtap
```

Verify:

- `GET /healthz`
- `GET /readyz`
- dashboard loads at `/app/dashboard/`
- checkout loads at `/app/checkout/`
- portal loads at `/app/portal/`
- customer, product, price, and checkout session can be created through `/v1`
- checkout can be completed through `/api/checkout/sessions/{id}/complete`
- `/api/timeline`, `/api/objects`, and `/api/delivery-attempts` return local
  evidence
- `BILLTAP_RELAY_MODE=true` masks endpoint secrets, webhook signatures, and
  sensitive URL query values

## Fresh Checkout Smoke

Before publishing the release, verify a clean checkout can build from source:

```bash
tmpdir=$(mktemp -d)
git clone https://github.com/midagedev/billtap.git "$tmpdir/billtap"
cd "$tmpdir/billtap"
git checkout v0.1.0
npm ci
npm run build
go test ./...
docker build -t billtap:v0.1.0 .
```

If the tag has not been pushed yet, run the same commands against the release
branch before tagging.

## Artifact Policy

For `v0.1.0`:

- Source archive from the GitHub release is expected.
- Local binary builds are expected.
- Local Docker image builds are expected.
- Official Docker image publishing is optional and must not be claimed unless
  the image is actually pushed and documented.
- npm package publishing is out of scope while `package.json` remains private.

Future releases may add signed binaries, container provenance, SBOMs, and
registry publishing. Add those steps here before claiming them.

## Release Blockers

Block the release if any of these are true:

- Required verification fails without a documented accepted deferral.
- Public docs imply full Stripe parity.
- Public docs imply real payment processing.
- Public docs include live credentials, real customer data, private hostnames,
  private company-specific adoption material, or raw evidence artifacts.
- `docs/COMPATIBILITY.md` does not match the implemented API surface.
- `docs/PRODUCTION_BOUNDARIES.md` does not match relay, masking, retention, or
  card-data rejection behavior.

## Post-Release Checks

After publishing:

- Confirm the GitHub release page links to compatibility and production
  boundaries.
- Confirm issue templates and support docs route compatibility gaps and
  security reports correctly.
- Open a follow-up issue for any release automation that was intentionally
  deferred.
