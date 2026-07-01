---
title: "Auth that just works: JWT, OAuth & 2FA out of the box"
subtitle: "The feature everyone rebuilds and nobody enjoys — already wired in every Grit project, endpoints and all."
series: "The Daily Grit"
edition: 6
date: 2026-07-06
readingTime: "5 min"
author: "Muke JohnBaptist"
tags: [grit, auth, jwt, oauth, 2fa, security]
canonical: "https://gritframework.dev/blog/auth-that-just-works"
---

Authentication is the most re-implemented feature in software. Every project: hash
passwords, issue tokens, add refresh, bolt on social login, eventually add 2FA. In
Grit, it's already there the moment you scaffold.

## What ships on day one

- **JWT** access + refresh tokens with middleware
- **OAuth2 social login** (Google, GitHub)
- **Two-factor auth** — TOTP authenticator apps, 10 backup codes, trusted devices
- **RBAC** — role-based access control baked into route registration
- **Password reset + email verification** flows with templates

## The endpoints you already have

```
POST /api/auth/register        → JWT tokens
POST /api/auth/login           → JWT tokens (or totp_required + pending_token)
POST /api/auth/totp/verify     → exchange a TOTP code for JWT
GET  /api/auth/me              → current user (protected)
GET  /api/auth/oauth/:provider → Google / GitHub social login
```

Login is 2FA-aware out of the box: if a user has 2FA on, `login` returns
`totp_required` + a short-lived pending token instead of full JWTs, and the client
finishes the exchange at `/api/auth/totp/verify`.

## Protecting a route

On the API, the auth middleware guards a route group, and the current user id is
on the context — so your handlers stay thin and ownership-scoped:

```go
protected := v1.Group("/", middleware.Auth())
protected.GET("/products", handler.ListProducts) // c.GetString("user_id") available
```

Need a role gate? Generate the resource with `--roles`, or add
`middleware.RequireRole("ADMIN")` to the group.

## Protecting a page (frontend)

For customer-facing pages, Grit gives you two helpers with one command:

```bash
grit add web-auth
```

That drops in `middleware.ts` (an SSR cookie check that redirects to `/login`) and
`<ProtectedWebRoute>` (a client guard using the `useMe()` hook). Wrap a page or
edit the matcher — done.

## Turning on 2FA (the flow)

```
1. POST /api/auth/2fa/enable   → QR code + secret + 10 backup codes
2. user scans the QR in any authenticator app
3. POST /api/auth/2fa/verify { code } → 2FA active
4. next login: password → then POST /api/auth/2fa/validate { code } → JWT
```

Backup codes are hashed like passwords and single-use. Trusted devices skip 2FA
for 30 days via a sliding cookie.

## Why this matters

Auth is where security bugs love to hide — token handling, enumeration, brute
force, session fixation. Grit's auth follows OWASP-2025 patterns and pairs with
the rest of the security battery (rate limiting, CSRF for OAuth flows, a
tamper-evident audit log). You get a reviewed baseline instead of a from-scratch
attempt under deadline.

**Tomorrow — the finale of week one:** from idea to production with a single
command, `grit deploy`.

*Go + React. Built with Grit.*
