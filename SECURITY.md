# Security Policy

## Reporting a vulnerability

**Please do not open a public issue for a security vulnerability.**

Report it privately through
[GitHub Security Advisories](https://github.com/MUKE-coder/grit/security/advisories/new),
which lets us discuss and fix the issue before any of it becomes public. If you
cannot use that, email **gmukejohnbaptist@gmail.com** with `SECURITY` in the
subject.

Please include:

- what the issue is and which component it affects (CLI, scaffolded API,
  admin panel, web app, a `grit-plugins` package)
- the version — `grit version`, and the tag in the affected project's `grit.json`
- steps to reproduce, ideally against a freshly scaffolded project
- what an attacker gains

### What to expect

| | |
|---|---|
| First response | within 3 business days |
| Triage decision (accepted / not-a-vuln / needs info) | within 7 days |
| Fix for a confirmed high or critical issue | targeted within 30 days |
| Credit | in the advisory and the changelog, unless you'd rather not be named |

We will tell you when a fix ships and coordinate disclosure timing with you. If
we disagree that something is a vulnerability, we will say so and explain why
rather than letting the report go quiet.

## Supported versions

Grit moves fast and does not maintain long-term support branches. **Fixes land
on the latest minor release.** If you are on an older version, `grit update`
and re-run your tests.

| Version | Supported |
|---|---|
| latest `3.x` | yes |
| older `3.x` | upgrade to latest |
| `2.x` and earlier | no |

## Scope

**In scope**

- The `grit` CLI itself.
- Code Grit *generates* — a vulnerability in a scaffolded API, admin panel, or
  web app is a vulnerability in Grit, because every user gets that code. This
  is the category we care most about.
- The `grit-plugins` packages.
- The release and distribution pipeline.

**Out of scope**

- Vulnerabilities in third-party dependencies with no Grit-specific exposure —
  report those upstream. If Grit's *use* of a dependency is what creates the
  exposure, that is in scope.
- Findings that require an attacker to already have the victim's credentials,
  server access, or database access.
- Missing hardening with no demonstrable impact (for example "header X is not
  set" without an attack that header would have stopped).
- Anything in the demo/example projects under `docs/`.
- Rate limiting on a locally-run development server.

## Verifying a release

Every release is signed and carries a Software Bill of Materials and SLSA build
provenance. Nothing here requires trusting us — the checks below verify against
public transparency logs.

**Checksums**

```bash
curl -LO https://github.com/MUKE-coder/grit/releases/download/vX.Y.Z/SHA256SUMS
sha256sum -c SHA256SUMS --ignore-missing
```

**Signature** (keyless — no key to steal, identity recorded in Rekor):

```bash
cosign verify-blob \
  --certificate SHA256SUMS.pem \
  --signature SHA256SUMS.sig \
  --certificate-identity-regexp '^https://github.com/MUKE-coder/grit/.github/workflows/release.yml@refs/tags/v' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  SHA256SUMS
```

**Build provenance** — proves the binary was built by our release workflow from
this repository, not uploaded by hand:

```bash
gh attestation verify grit-linux-amd64.tar.gz --repo MUKE-coder/grit
```

**SBOM** — `grit-sbom.spdx.json` is attached to every release for dependency
scanners.

## What Grit generates, security-wise

Scaffolded projects are hardened by default rather than as an opt-in. The
relevant behaviour is documented under
[Backend → Authentication](https://gritframework.com/docs/backend/authentication)
and [Security](https://gritframework.com/docs/security):

- bcrypt password hashing; tokens in HttpOnly cookies, never `localStorage`
- server-side sessions — refresh tokens are revocable, rotate on use, and a
  replayed token kills the session
- single-use, hashed, expiring password reset tokens; resetting signs out every
  device
- CSRF protection on cookie-authenticated requests, strict security headers,
  brute-force lockout and rate limiting on auth routes
- SSRF-safe outbound HTTP for user-supplied URLs

If you find a way around any of these in generated code, that is exactly the
kind of report we want.
