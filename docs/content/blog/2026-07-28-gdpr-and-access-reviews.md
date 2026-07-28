---
title: "Compliance you didn't have to build: GDPR erasure and access reviews"
subtitle: "Two admin surfaces every Grit app already ships with — a GDPR export/erase toolkit backed by a tamper-evident deletion journal, and access reviews that snapshot every role assignment for recertification. What they do, and the question that actually matters: how do they get populated?"
series: "The Daily Grit"
edition: 13
date: 2026-07-28
readingTime: "9 min"
author: "Muke JohnBaptist"
tags: [grit, security, gdpr, compliance, access-reviews, audit, go, react]
canonical: "https://gritframework.dev/blog/gdpr-and-access-reviews"
---

Two features you never look forward to building, and always end up needing, are
**GDPR data requests** and **access reviews**. They're not glamorous. But the
first day a user emails "delete my account and send me my data," or the first
time an auditor asks "who has admin, and who signed off on that?", you want them
already there.

In Grit they are. Every scaffolded app ships both, admin-only, under
*Security & Access* in the [System Hub](/docs/admin/overview). This post walks
through what each does — and spends most of its time on the question people
actually ask when they open these pages for the first time: **where does the
data come from? How do they get populated?**

The short answer, for both: **on demand, by an admin action.** Nothing is on a
schedule, and nothing is seeded. The tables start empty and stay empty until
someone does something. That's not a gap — for audit evidence, it's the point.

## Part 1 — The GDPR toolkit

Open `/system/gdpr`. You get a box for a user's UUID and two things you can do
with it: **export** their data, or **erase** them.

### Export reads; it never writes

`GET /api/users/:id/gdpr-export` gathers everything the app holds about one
person — profile, uploads, sessions, activity, dashboard layout, whether 2FA is
on — into a single JSON bundle and streams it as a download. A user can export
their own data; an admin can export anyone's.

The important thing about export is what it *doesn't* do: it creates no records.
It's a pure read. So it never appears in the journal we're about to talk about.
Populating the journal is the job of the *other* button.

### Erase writes exactly one journal row

`POST /api/users/:id/gdpr-erase` is admin-only and refuses to erase yourself. In
one transaction it:

1. **Hard-deletes the user's child PII** — uploads, sessions, password-reset
   tokens, role grants, 2FA configs, trusted devices, pending TOTP tokens,
   dashboard layouts, notifications — counting each table as it goes.
2. **Anonymizes the user row in place** instead of deleting it, so foreign keys
   in records you must keep still resolve. Name becomes `Erased User`, email
   becomes `erased-<id>@deleted.invalid`, secrets and device identifiers are
   blanked, the account is deactivated and demoted.
3. **Writes one `DeletionJournal` row** describing what just happened.

That third step is the only way the journal grows: **one row per erasure.** Not
per export, not on a timer, not from a seeder.

### Why the journal is "tamper-evident"

A journal row deliberately holds *no* PII — just the erased user's UUID, who ran
the erasure, a reason, per-table counts, and a timestamp. Plus two hashes,
because each row is chained to the one before it like a miniature blockchain:

```text
hash = sha256( prev_hash + canonical(entry) )

# genesis row's prev_hash is ""
# every later row's prev_hash = the previous row's hash
```

Rows are append-only — never updated, never deleted. `GET /api/gdpr/journal`
replays the whole chain, recomputes every hash, and checks each row links to the
last. Change or drop a row after the fact and the recompute won't match: the
replay reports the break, and the admin page flips its **"Chain verified"** pill
to **"Chain broken."** You can hand that to an auditor and say: not only did we
log every erasure, we can prove the log wasn't edited.

(One nice detail: the erasure leaves the [activity log](/docs/batteries/security)
alone — those rows carry only a UUID and have their *own* hash chain — but it
also drops a `user.gdpr_erase` event there, so the erasure still shows up in the
dashboard and any SIEM export.)

## Part 2 — Access reviews

Now open `/system/access-reviews`. Same idea, different compliance question:
*does everyone who has access still need it?* SOC 2 and ISO 27001 want you to ask
that on a cadence and keep the evidence. That's an **access review** (or access
recertification), and it's built on two records: a **campaign** and its **items**.

### Opening a review *is* the snapshot

Here's the populate step. An admin clicks **New review**, names it (say
"Q3 2026 quarterly"), and `POST /api/access-reviews` opens a campaign — then, in
the same transaction, **snapshots every current role assignment.** It reads the
`user_roles` table joined to users and roles, and creates one review *item* per
grant:

```text
AccessReview  (the campaign)   status: "open"
  └─ AccessReviewItem  (one per current user→role grant)
       user_email: "ada@acme.com"   ← copied in at open time
       role_name:  "ADMIN"          ← copied in at open time
       decision:   "pending"
```

The email and role name are **copied into the item**, not linked. So if someone
later renames the ADMIN role or deletes a user, the review still reads exactly as
it did the day you opened it. That's what makes it a *point-in-time* snapshot —
frozen evidence, not a live view.

And again: it exists only because an admin opened it. No schedule, no seeder. The
snapshot captures the full set of assignments as they stood that moment, every
item `pending`.

### Deciding an item changes real access

For each pending item the admin makes one call:

- **Keep** (`approved`) — certifies the grant. This *is* the attestation; there's
  no separate "attest" button.
- **Revoke** (`revoked`) — deletes the real `user_roles` grant in the same
  transaction, then records the decision. A revoke is **terminal** and is logged
  as `access_review.revoke`.

So a review isn't paperwork that sits next to reality — revoking an item *is* the
off-boarding. Once every item has a decision, **Complete review** signs the
campaign off (it won't let you while anything is still pending). A completed
review is immutable, and never reopened — next quarter is a fresh campaign.

### What the admin sees

A two-column page: campaigns on the left with a status pill and
`{pending} · {approved} · {revoked}` counts; on the right, the selected review's
items table (User, Role, Decision) with **Keep** / **Revoke** per pending row and
a **Complete review** button that unlocks once nothing is pending.

## The through-line

Both features answer their "how does it get populated?" the same way, and it's
worth saying plainly: **an admin populates them, on purpose, one action at a
time.** The GDPR journal gains a row when someone erases a user. An access
review's items appear the moment someone opens the campaign. Neither runs on a
schedule, and neither is faked with seed data.

That's deliberate. Audit evidence you can trust is evidence tied to a real human
decision, timestamped, and — in the journal's case — cryptographically sealed
against later edits. You didn't have to build any of it, but you can stand behind
all of it.

📖 **Docs:** [Privacy & Compliance — GDPR + Access Reviews](/docs/security/compliance) ·
[Roles & Permissions](/docs/security/authorization)
