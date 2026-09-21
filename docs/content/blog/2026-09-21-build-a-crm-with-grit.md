---
title: "Build a CRM with Grit: phone numbers from every country, a million contacts, and who else is looking"
subtitle: "One project that uses everything new in Grit this month: ten field types that check and clean their own values, a seeder that fills a table with a million rows and picks up where it stopped, a live 'who else has this contact open' banner, import and export, and lint that passes on day one. Every command was run on a real project while writing this, and every number is the one it printed."
series: "The Daily Grit"
edition: 16
date: 2026-09-21
readingTime: "22 min"
author: "Muke JohnBaptist"
tags: [grit, crm, field-types, seeding, realtime, presence, biome, tutorial]
canonical: "https://gritframework.dev/blog/build-a-crm-with-grit"
---

A CRM is a good test of a framework, because it is mostly boring data done right. A contact has a phone number, and the phone number has to be real, in one format, from any country. A company has a website and a domain and nobody should be able to type "blue" into its brand colour. There are a lot of rows, so the list has to stay fast and the test data has to be big enough to show you where it will not. And two sales reps must not phone the same lead at the same time.

This guide builds that CRM. It is also a tour of what landed in Grit this month, because every one of these problems now has a feature behind it:

- **Ten field types** (`email`, `url`, `domain`, `tel`, `country`, `color`, `percent`, `rating`, `time`, `json`), each with its own admin input, its own checks on the API, and realistic seed data. Grit v3.294.0.
- **A seeder that fills a million rows** in batches, shows its progress and resumes where it stopped: `grit seed Contact --count 1000000`. Grit v3.295.0.
- **Realtime presence and client events**, which is how the "who else is looking" banner works. Grit v3.281.0.
- **Biome** for lint and format, so `pnpm lint` passes on a new project. Grit v3.293.0.

You need Grit v3.295.0 or later. Check with `grit version`, and run `grit update` if it is older.

---

## 1. A new project, no database server needed

```bash
grit new crm --triple --next --db sqlite
cd crm
```

`--triple` gives you the Go API, a Next.js web app and a separate admin panel. `--db sqlite` means you need nothing running to follow along: the database is a file. Everything in this guide works the same on Postgres, and the numbers at the end show both.

---

## 2. Companies: six field types in one command

```bash
grit generate resource Company --fields "name:string,website:url,domain:domain:unique,brand:color,country:country:UG,industry:select:tech=Technology|retail=Retail|health=Health|finance=Finance,profile:json" --faker --count 50
```

Here is what each field became:

| Field | Type | Stored as | In the admin |
|---|---|---|---|
| `website` | `url` | `http` or `https` only | a URL input; a link that opens in a new tab |
| `domain` | `domain:unique` | a bare host, lowercase, international names as punycode (`münchen.de` is stored as `xn--mnchen-3ya.de`) | an input that strips a pasted `https://` and path |
| `brand` | `color` | `#rrggbb`, lowercase, so `#ABC` becomes `#aabbcc` | a colour picker plus a hex box; a swatch in the table |
| `country` | `country:UG` | the two-letter ISO code, defaulting to `UG` | a searchable picker with flags, 249 countries |
| `industry` | `select` | one of the four values | a dropdown that cannot offer anything else |
| `profile` | `json` | JSON (a text column on SQLite, `jsonb` on Postgres) | an editor that checks the JSON as you type |

The third part of a field is its option: `country:UG` sets the default, `:unique` adds a unique index. `--faker --count 50` also writes a seeder that fills the table with 50 fake companies, and we will use it in a moment.

---

## 3. Contacts: phone numbers from every country

```bash
grit generate resource Contact --fields "first_name:string,last_name:string,email:email:unique,phone:tel:UG:unique,company:belongs_to:Company,rating:rating,best_time:time,win_chance:percent,stage:select:lead=Lead|qualified=Qualified|customer=Customer|lost=Lost" --faker --count 100
grit migrate
```

The interesting one is `phone:tel:UG:unique`.

A `tel` field is stored in **E.164**, the international format: a plus, the country code, the number, no spaces. `0772 123456` typed by a rep in Kampala and `+256 772 123 456` pasted from an email signature are the same number, so they are stored the same way, `+256772123456`, and your unique index can actually do its job. The `UG` is the default country: it is what a number without a `+` is read as.

Every number is checked against libphonenumber's metadata, Google's database of how numbers look in each country, on both sides: `libphonenumber-js` in the admin and a Go port of it in the API. They use the same data, so the form never accepts a number the API will then refuse.

In the admin the field is a country picker (flag, name and dial code, searchable by any of them) next to an input that formats the number as you type. Type `+44` and the picker switches to the United Kingdom by itself.

The other new types in this resource:

- `email:unique` is stored lowercase, so `Ada@Example.COM` and `ada@example.com` cannot both exist.
- `rating` is a whole number from 1 to 5, shown as stars you can set with the keyboard. `rating:10` would give ten.
- `best_time` is a time of day, stored as `HH:MM`: `9:05` is stored as `09:05`.
- `win_chance` is a `percent`, a number from 0 to 100 shown as `62.5%`.

---

## 4. A million contacts

The generated seeders are ready. Fill the companies first, because every contact needs one to belong to:

```bash
grit seed Company --count 1000
grit seed Contact --count 1000000
```

This is what the second command printed, trimmed:

```text
contacts: 23000 of 1000000 rows (2875 rows/s, about 5m39s left)
...
contacts: seeded 1000000 rows in 14m37.251s (1140 rows/s), now 1000000; heap in use 14 MB
```

A million contacts, each linked to a real company, each with a valid phone number, in under fifteen minutes on SQLite on a laptop, using 14 MB of memory the whole way. Before v3.295.0 the same seeder inserted one row per statement, and on the same machine that works out to hours.

Four things are going on:

**It inserts in batches.** Rows go in with multi-row inserts, one transaction per batch, as many rows per statement as the database accepts for the table's column count, up to 1,000. GORM hooks still run for every row, so the IDs and anything your model's `BeforeCreate` does work as usual.

**It tops up.** `--count` is the number of rows the table should hold, not the number to add. The seeder counts what is there and inserts the rest. Run the same command again and it says so and does nothing:

```text
contacts: 1000000 rows, target 1000000, nothing to seed
```

That also makes it safe to stop. Press Ctrl+C halfway through, run the command again, and it carries on from where it stopped, to exactly the number you asked for.

**Unique columns stay unique.** At a million rows, random values collide. The seeder builds each row from its number in the table, so `email` and `phone` are distinct for every one of the million, and a column you marked unique yourself, a SKU say, seeds as `SKU-0000001`, `SKU-0000002` and so on.

**It fails loudly.** The first batch that fails stops the run with an error that says which rows. The old seeder logged a warning for each failed row and reported success at the end, which is how you get a "seeded" table that is half empty.

A note on the numbers. SQLite takes one writer at a time, and this table carries two unique indexes (email and phone) plus a foreign key, so it starts near 2,900 rows a second and slows as the indexes grow. On Postgres, which the seeder writes three batches at a time, a simpler contact table seeded a million rows in about 50 seconds. If you are going to seed at this size, do it on the database you will run.

The seeder is a normal Go file, `apps/api/internal/database/contacts_seeder.go`, and you can change any value in it. The part you will edit looks like this:

```go
Make: func(n int64) models.Contact {
    i := int(n)
    return models.Contact{
        FirstName: gofakeit.FirstName(),
        LastName:  gofakeit.LastName(),
        Email:     fieldtypes.SampleEmail(gofakeit.FirstName(), gofakeit.LastName(), gofakeit.DomainName(), i),
        Phone:     phone.Sample(i),
        CompanyID: pickID(companyIDs), // a real, existing company
        // ...
    }
},
```

`i` is the row's number. The `Sample` helpers use it to keep values unique and valid: `phone.Sample(i)` spreads the numbers across 20 countries and only ever returns one libphonenumber accepts.

---

## 5. What the API does with a value

Start the project with `grit start`, sign in to the admin, and open Contacts. The list says 1,000,000. Now try it from the API, as a signed-in admin, because this is where the field types earn their keep.

Create a contact the way a person would type it:

```json
POST /api/v1/contacts
{
  "first_name": "Ada",
  "last_name": "Nakato",
  "email": "Ada.Nakato@Example.COM",
  "phone": "0772 123456",
  "company_id": "…",
  "rating": 4,
  "best_time": "9:05",
  "win_chance": 62.5,
  "stage": "lead"
}
```

It answers `201`, and what it stored is the clean version:

```json
{ "email": "ada.nakato@example.com", "phone": "+256772123456", "best_time": "09:05", "win_chance": 62.5 }
```

Now a number that is too short to be a Ugandan number:

```json
{ "phone": "+256 12", ... }
```

```json
422
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Phone is not a valid phone number for UG",
    "details": { "phone": "Phone is not a valid phone number for UG" }
  }
}
```

The check lives on the model, not the form, so it also applies to updates, bulk edits, CSV imports and offline sync. There is no path into the table that skips it.

---

## 6. Who else has this contact open

Two reps phoning the same lead ten minutes apart is how a CRM loses a customer. We will add a banner to the contact page that shows who else has the contact open, and tells everyone else the moment one of them starts a call.

This uses two realtime features Grit already has:

- A **presence channel** keeps a live list of who is subscribed. The server attaches each person's name.
- A **client event**, or whisper, is a message one browser sends to the others on a channel. It goes through the server and is never stored. It is made for exactly this: "is typing", "is calling", "moved the cursor".

**First, the API.** A private or presence channel is refused unless you register who may join it. Add this in `apps/api/internal/routes/routes.go`, right after the line `realtime.AllowedOrigins = corsOrigins`:

```go
// presence-contacts.<id>: who has this contact open right now. Anyone
// signed in to the admin can open a contact, so the check is that the
// contact exists; the viewer's name travels with them to everyone else.
realtime.Channel("contacts.{id}", func(c realtime.ChannelContext) bool {
    var exists int64
    if err := db.Model(&models.Contact{}).Where("id = ?", c.Param("id")).Count(&exists).Error; err != nil || exists == 0 {
        return false
    }
    var viewer models.User
    if err := db.Select("first_name", "last_name").First(&viewer, "id = ?", c.UserID).Error; err != nil {
        return false
    }
    c.SetInfo(map[string]any{"name": strings.TrimSpace(viewer.FirstName + " " + viewer.LastName)})
    return true
})
```

One registration covers both `private-contacts.42` and `presence-contacts.42`. Whatever you pass to `SetInfo` is shown to the other members, so only put in it what they are allowed to see.

**Then the banner.** Create `apps/admin/components/contact-presence.tsx`:

```tsx
"use client";

import { useState } from "react";
import { useMe } from "@/hooks/use-auth";
import { useChannel, usePresence, useWhisper } from "@/hooks/use-realtime";
import type { ClientEvent } from "@/lib/realtime";

type Viewer = { name: string };
type Calling = { name: string; on: boolean };

export function ContactPresence({ id }: { id: string }) {
  const channel = `presence-contacts.${id}`;
  const { data: me } = useMe();
  const viewers = usePresence<Viewer>(channel).filter((v) => v.user_id !== me?.id);
  const whisper = useWhisper(channel);
  const [caller, setCaller] = useState<string | null>(null);
  const [onCall, setOnCall] = useState(false);

  useChannel(channel, {
    "client-event:calling": (p: ClientEvent<Calling>) => setCaller(p.data?.on ? (p.data.name ?? "Someone") : null),
  });

  const toggleCall = () => {
    const next = !onCall;
    setOnCall(next);
    whisper("calling", { name: me?.first_name ?? "Someone", on: next });
  };

  return (
    <div className="mb-4 flex flex-wrap items-center gap-3 rounded-lg border border-border bg-bg-secondary px-4 py-3 text-sm">
      {viewers.length === 0 ? (
        <span className="text-text-secondary">Only you have this contact open.</span>
      ) : (
        <span className="text-text-secondary">
          Also viewing: <strong className="text-foreground">{viewers.map((v) => v.info?.name ?? "a teammate").join(", ")}</strong>
        </span>
      )}
      {caller && <span className="rounded-md bg-warning/15 px-2 py-1 text-warning">{caller} is calling this contact now</span>}
      <button type="button" onClick={toggleCall} className="ml-auto rounded-md border border-border px-3 py-1.5 hover:bg-bg-hover">
        {onCall ? "End call" : "Log a call"}
      </button>
    </div>
  );
}
```

**And put it on the page.** Open `apps/admin/app/(dashboard)/resources/contacts/[id]/page.tsx` and render the banner above the detail view:

```tsx
"use client";

import { use } from "react";
import { ContactPresence } from "@/components/contact-presence";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { contactResource } from "@/resources/contacts/contacts";

export default function ContactsDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return (
    <>
      <ContactPresence id={id} />
      <ResourceDetailPage resource={contactResource} id={id} />
    </>
  );
}
```

Open the same contact as two different users in two browsers. Each sees the other's name. Click **Log a call** in one, and the other shows "Ada is calling this contact now" at once.

A few rules the server enforces for you, so you do not have to: a client event only goes to `private-` and `presence-` channels, only from a connection already subscribed to that channel, at most ten a second per connection, with a payload under 1 KB. The `user_id` on a received event is filled in by the server from the sender's session, so you can trust who sent it. The rest is whatever the other browser sent, so treat it as user input.

This is the same mechanism that powers typing indicators in a chat. The realtime docs cover channels, presence and client events in full.

---

## 7. Import, export and PDFs

You did not write any of this, and it is already there. Every generated resource gets:

- **Export** from the admin toolbar as CSV or Excel, with the current filters applied: `GET /api/v1/contacts/export`.
- **Import** from CSV: `POST /api/v1/contacts/import`, with a template to fill in at `GET /api/v1/contacts/import/template`. Imported rows go through the same field checks, so a bad phone number is reported with its row number rather than stored.
- **A PDF of one record**: `GET /api/v1/contacts/:id/pdf`.
- **Bulk actions**: select rows in the table and change or delete them together.

Big imports run in the background as a job with progress, so a large CSV does not tie up the request that uploaded it.

---

## 8. Lint passes on day one

```bash
cd apps/admin && pnpm lint
```

On this project it checks all 208 files in the admin and reports nothing. That includes the banner we just wrote. Grit projects lint with Biome now: one tool for lint and format, one config file (`biome.jsonc`) at the root, and every rule that is switched off has its reason written next to it. `pnpm format` formats the whole project in the style the templates use, if you want to turn formatting checks on too.

---

## 9. Before you put it in front of customers

A CRM holds personal data, so a few defaults are worth knowing about:

- **Sign-in says nothing about accounts.** A wrong password for an unknown, locked, disabled or social-login account gets the same answer in the same time, so nobody can use your sign-in page to find out who your users are.
- **Two-factor secrets are encrypted.** `grit new` wrote a `FIELD_ENCRYPTION_KEY` into `.env`. Back it up somewhere other than your database backups: if you lose it, every user with two-factor is locked out. The deployment checklist in the docs covers this.
- **Forwarded headers are trusted only from your proxy.** If you deploy behind a load balancer that is not on a private network, set `TRUSTED_PROXIES` so client IPs in your audit log and rate limits are real.
- **Seed data stays out of production.** The admin account `grit seed` creates exists only with `APP_ENV=development`. Anywhere else the seeder needs a `SEED_ADMIN_PASSWORD` of at least 12 characters.

---

## What you built, and what it cost

About six commands, one Go block and one React component, and you have:

- Companies and contacts with values that are checked and stored one way: real phone numbers from any country, lowercase emails, bare domains, proper colours, times and percentages
- An admin with an input made for each of those, including a country picker and a phone field that formats as you type
- A million realistic contacts, seeded in minutes, with a command you can stop and resume
- A live banner that stops two reps calling the same lead
- Import, export, PDFs and bulk actions you did not write
- Lint that passes

The parts you would still add for a real team: a deals pipeline (`grit generate resource Deal --fields "..."` with a `belongs_to:Contact` and a `stage` select, and the same banner), activity logging per contact, and email. For the email, `grit generate mail FollowUp` writes a typed template with a `Queue` helper, so the send goes through the background job queue and survives a restart.

If you build something with it, tell us about it on GitHub. The next post in this series builds the first Grit UI blueprint: a WhatsApp clone, on web, mobile and desktop, on the same realtime machinery as the banner above.
