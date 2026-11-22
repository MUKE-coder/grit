---
title: "What's new in Grit: October 2026"
subtitle: "179 releases in one month, which is too many to read. Here are the ten that change something for a project you already have, what to run, and the six security fixes you should not skip."
series: "What's new in Grit"
date: 2026-11-22
readingTime: "7 min"
author: "Muke JohnBaptist"
tags: [grit, release-notes, changelog]
canonical: "https://gritframework.dev/blog/whats-new-in-grit-october-2026"
---

Grit shipped 179 releases in October. That number is the point of this post and also the problem with it: the [changelog](/docs/changelog) says why each one happened, which is exactly what you want at two in the morning and exactly nothing you want on a Monday.

So this is the first of a monthly digest. Ten changes that affect a project you already have, what to run, and what you can ignore. Nothing in October was marked breaking.

## Run this first

```bash
grit update        # bring the CLI up to date
grit upgrade       # bring the project's framework files up to date
grit migrate       # the new tables the month added
grit doctor        # what is now checkable that was not before
```

If you only do one thing, do `grit upgrade`. Six of October's fixes are security fixes in code that lives in your repository, and nothing reaches you until you run it.

## The security fixes, because they are the ones that matter

**The audit chain protected against nobody (v3.456.0).** Every row of the activity log is hashed over the row before it, so changing one row breaks every hash after it. The hash was `SHA-256(prev_hash || row)`, unkeyed, over values anybody can read off the source. The attack is: edit a row, recompute the rest, and verification passes. The package comment said it defended against a DBA covering tracks, which is precisely the one person it did not defend against.

It is `HMAC-SHA256` now. The key is `AUDIT_CHAIN_KEY`, derived from `JWT_SECRET` through SHA-256 when unset, so you get a keyed chain without setting a second secret. Set it explicitly when the chain should outlive a `JWT_SECRET` rotation, and keep it somewhere your database administrator is not, because that separation is the entire mechanism.

**A provider sign-in trusted an address the provider never verified (v3.453.0).** goth puts the provider's userinfo response in `RawData`, `email_verified` is in there, and nothing read it. Three things followed from that one missing check: a new account was created with `EmailVerifiedAt` already set; an existing account had the provider linked to it, which is account takeover with one more step; and an existing account whose address was never confirmed had its password removed and its sessions revoked, which is account takeover outright. Any provider without an explicit rule now counts as unverified, so adding one cannot quietly reopen it.

**Every policy rule was dead (v3.451.0).** `policies.Register()` had no call site anywhere in a generated project. Every rule in `internal/policies` registered into a map nothing ever filled, so a `PATCH` on a row a rule denied answered 200. That is the worst shape an authorization feature can fail in, because everything about it looks finished: the package is generated, documented, and carries a worked example.

**Two privilege escalations in the multitenant plugin (v3.438.0).** One needed no attack at all. It was the happy path: creating an organization assigned the platform ADMIN role as your membership role, and ADMIN holds `*`. Anybody who created an organization became a platform administrator.

**Cross-tenant reads through `Scan` (v3.429.0).** The organization filter was installed on GORM's Query, Update, Delete and Create callbacks. There is a fifth, Row, and it is the one `Scan` goes through, so `db.Model(&Order{}).Scan(&out)` on a tenant-owned model returned every organization's rows. Run `grit plugin update multitenant`: the plugin reports 1.0.1 once you have the fix.

**The admin panel loaded for people who were not signed in (v3.459.0).** The whole page arrived, then React asked the API who you were, got a 401 and replaced the route with the login. Nothing leaked, because the API refuses every request behind that panel and that is the check that decides anything. But the panel was sent to a stranger and its shape was readable in the response.

## What is new that you might want

**`grit new --multitenant` (v3.448.0)**, with organisations, invitations, per-organisation roles and a switcher in the admin. Multitenancy stopped being only a plugin you bolt on afterwards and became a shape you can scaffold.

**`grit plugin new` (v3.461.0)**, which scaffolds a plugin with its manifest, its install and its removal. The removal is the half people forget: a plugin that cannot be removed cleanly is one nobody tries.

**`grit upgrade --dry-run` (v3.439.0).** `--diff` was never a dry run: it upgraded and showed the diff afterwards. Now there is one that changes nothing.

**`grit update` verifies what it is about to run (v3.450.0).** The binary it downloads is checked against the release's signed checksums before it replaces the one you have.

## Two bugs worth knowing about even now

**`grit sync` undid what `grit generate` had just written (v3.460.0).** Two writers produced `packages/shared/types/<model>.ts` and the second to run won, which was usually sync, because the generator prints "run grit sync" as a next step. What each got wrong was the opposite of the other, which is why neither looked broken alone: the generator knew `Money`, `FileRef` and an enum's union and sync replaced them with `unknown` and `string`, while sync knew six columns the generator had never heard of. There is one writer now. If you generated a resource in early October and a shared type looks thinner than its model, regenerate.

**An injection recorded your edit as Grit's (v3.454.0).** The manifest exists so `grit upgrade` can overwrite the files nobody has touched and leave the rest alone. After a generator injected into a file, the manifest re-hashed it, which is right for a file Grit owns and wrong for one you had already edited: your work was recorded as Grit's, and the next upgrade saw an untouched file and overwrote it without a word. If you edited `routes.go` or another injected file in October and later lost those edits to an upgrade, that was this.

**Generated tests passed on SQLite and failed on Postgres (v3.440.0).** Every resource with a `belongs_to` had a handler test that failed on the database the project actually runs on, with a foreign key violation: the factory leaves the parent column empty on purpose, the generated test did not pass a parent in, and an empty string is not NULL. The tests create the parent first now.

## The rest

179 releases is 179 entries, and most of them are a fix with a one-line story. The [full changelog for October](/docs/changelog) has all of them, newest first, each saying what was wrong and why.

A note on the number itself, since it comes up: every merged change ships as a release, which is why the version is in the three thousands. From v3.567.0 there is a [stable channel](/docs/versioning#channels) for teams that cannot take a release every day:

```bash
grit update --channel stable
grit upgrade --to stable
```

stable is the newest release that is at least a week old and had every nightly canary run in that week come back green. No backports: a fix lands on edge and reaches stable when it has aged.

Next digest on the first of the month.
