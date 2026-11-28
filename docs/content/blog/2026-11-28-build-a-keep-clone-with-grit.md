---
title: "Build a Google Keep clone with Grit, organisations and billing included"
subtitle: "A note-taking SaaS in one sitting: teams, invitations, Stripe and a sign-in area you did not write, plus notes, reminders, archive and trash that are entirely yours. The point of this build is the line between the two, and where to draw it."
series: "The Daily Grit"
edition: 26
date: 2026-11-28
readingTime: "22 min"
author: "Muke JohnBaptist"
tags: [grit, saas, multitenancy, stripe, nextjs, golang, tutorial, portal]
canonical: "https://gritframework.dev/blog/build-a-keep-clone-with-grit"
---

Here is a thing I used to believe about code generators: they are for the boring half, and the interesting half is yours. That is right, and for a long time Grit got the split wrong anyway. It generated a complete admin screen for every resource and handed the customer a React Query hook. For an internal tool that is the whole job. For a product with customers it is half of one, and the missing half is the half people pay for.

v3.584.0 changes that, and the clearest way to show what changed is to build something small and real with it. So: Google Keep, with the parts a real product has. Notes belong to a workspace rather than a person. You can invite your team. There is a plan and a card. And the note grid, the thing that makes it Keep rather than a CRUD table, is written by hand, because of course it is.

Everything below is a real build. The screenshots are of it running.

## What we are making

| Screen | Where it comes from |
|---|---|
| Sign in, register, account, profile | Grit, from `--saas` |
| Organisation: members, invitations, roles | the multitenant plugin |
| Billing: plan, renewal, cancel, card | the Stripe plugin |
| Labels: list, create, edit | generated from the resource |
| **Notes, reminders, archive, trash** | **us** |

One of those rows is the product. The other four are the furniture that every SaaS has and nobody is paying you to design differently.

## One command

```bash
grit new keep --saas --saas-organizations --saas-billing
```

Drop the two long flags and it asks instead:

```
Does this app have customers?        > Yes, customers sign in
Do customers belong to organizations? > Yes, organizations
Do you charge them?                   > Yes, take payments
```

`--saas` is not a new scaffolder. It writes `"saas": true` into `grit.json` and then runs three things that already existed: `grit add web-auth` for the signed-in customer area, `grit plugin add multitenant` for organisations, `grit plugin add stripe` for payments. What it adds is being asked once, at the point where you know the answer, instead of discovering three commands separately over a fortnight.

The order they run in matters, which I know because I got it wrong first. The organisation screen and the billing screen are pages *inside* the customer area, so the customer area has to exist before either plugin looks for it. Installed the other way round, both land outside the shell: no sidebar, no header, no link from anywhere. More on that in a minute, because it was a real bug and not a hypothetical one.

I used SQLite and local disk so there is nothing to start:

```bash
grit new keep --triple --next --db sqlite --storage local \
  --saas --saas-organizations --saas-billing
cd keep && pnpm install && grit migrate && grit seed
```

## Before writing a line of our own

```bash
grit start
```

Sign up at `localhost:3000/register`, and the first thing that happens is the app refusing the password I tried to use:

```json
{"error":{"code":"VALIDATION_ERROR",
  "message":"Your password cannot contain your name or your email address."}}
```

I had typed `keep-password-123` for `ada@keep.test`. That check was already there. It is the sort of thing you would get round to in month three.

Then the organisation screen, which is where a tenant app sends you first:

![The organisation onboarding screen inside the customer shell: a sidebar with Notes, Overview, Profile, Organization, Billing and Labels, and a single Organization name field](/blog/keep/organization.png)

Note the sidebar. Overview and Profile came from `web-auth`, Organization from the multitenant plugin, Billing from Stripe, Labels from a resource we will generate in a moment. None of those entries were typed by hand.

## The model, and the one decision that matters

```bash
grit generate resource Note \
  --fields "title:string,body:text,colour:select:yellow|green|blue|pink,pinned:bool,archived:bool" \
  --tenant-owned --portal headless

grit generate resource Label \
  --fields "name:string,colour:select:yellow|green|blue|pink" \
  --tenant-owned

# Reminders and trash, added to the model after the fact
grit generate field Note remind_at:datetime
grit generate field Note trashed_at:datetime
grit migrate
```

Those last three lines are worth pausing on. `grit generate field` puts the
column in the Go model, the create and update Zod schemas, the TypeScript
type, the handler, the service's writable and filterable lists, the CSV
importer, the admin form and table, **and the portal descriptor**. One
command, nine places, and `grit migrate` adds the column because the model is
the source of truth.

Two flags are doing the work.

**`--tenant-owned`** says a note belongs to the organisation, not to the person who typed it. That is a product decision with teeth: every query is scoped by a GORM callback, so a forgotten `WHERE org_id` cannot leak another workspace's notes, and you cannot opt out by accident. The alternative is `--owned-by user`, which is right for a single-player product and wrong here, because the whole point of inviting your team is that they can read what you wrote.

Get this one wrong and the fix touches every query, every handler and every test. It is worth five minutes before you type the command.

**`--portal headless`** is the new one, and it is the subject of this post.

## screens, headless, off

In a `--saas` project, a resource that belongs to somebody gets a customer screen as well as an admin one. Three levels:

| Level | What you get |
|---|---|
| `screens` | A list, a create page, a detail page, a form, and an entry in the customer sidebar. The default |
| `headless` | The resource as data, plus the hooks. No JSX |
| `off` | The hook alone, which is what every project got before v3.584.0 |

The rule I would give anybody: **the furniture is the same in every SaaS and the product is not.**

Labels are furniture. A list, a name, a colour, an edit form. Nobody is going to redesign that, and here it is, generated, including the empty state:

![The generated Labels screen: a heading, a New button, and an empty state reading "No labels yet, when you create one it will show up here" with a Create the first one button](/blog/keep/labels-generated.png)

Notes are the product. A generated list of notes would be a table of titles, and Keep is not a table of titles, it is a ragged grid of coloured cards you can pin. So `--portal headless`: Grit writes the API, the hooks and a descriptor, and stops.

Here is the whole of what `headless` leaves behind:

```typescript
// apps/web/lib/portal/notes.ts, generated
export const notePortal = {
  resource: "notes",
  label: "Notes",
  path: "/notes",
  api: "/api/notes",
  scope: "organization",
  display: "title",
  subtitle: "body",
  fields: [
    { name: "title", label: "Title", kind: "text", required: true },
    { name: "body", label: "Body", kind: "textarea" },
    { name: "colour", label: "Colour", kind: "select",
      options: ["yellow", "green", "blue", "pink"] },
    { name: "pinned", label: "Pinned", kind: "checkbox" },
    { name: "archived", label: "Archived", kind: "checkbox" },
    { name: "remind_at", label: "Remind at", kind: "datetime" },
    { name: "trashed_at", label: "Trashed at", kind: "datetime" },
  ],
} as const;
```

That is not much, and the little it is buys one specific thing: the colours in our picker come from here rather than from a list typed into a component, so adding `purple` to the Go model adds it to the picker with no frontend edit. Hard-code the four colours in JSX and you have made a second source of truth that will be wrong by March.

Two things the descriptor is not. It is **not a permission check**: the API scopes every read and write by the active organisation, and a filter living in a page is a check an API client walks straight past. And it does not include the owner column, so a customer cannot hand their row to another workspace by editing a form.

## The part that is ours

`apps/web/app/(app)/notes/page.tsx`. Nothing scaffolded it and nothing will overwrite it.

```tsx
"use client";

import { useState } from "react";
import { Pin, PinOff, Archive, Trash2, Plus } from "lucide-react";
import type { Note } from "@repo/shared/types";
import type { CreateNoteInput, UpdateNoteInput } from "@repo/shared/schemas";
import { useNotes, useCreateNote, useUpdateNote, useDeleteNote } from "@/hooks/use-notes";
import { notePortal } from "@/lib/portal/notes";

// The colours come from the generated descriptor, so the Go model stays the
// one place they are declared.
const COLOURS = (notePortal.fields.find((f) => f.name === "colour")?.options ?? []) as readonly string[];

const SWATCH: Record<string, string> = {
  yellow: "bg-amber-100 border-amber-200 dark:bg-amber-500/10 dark:border-amber-500/20",
  green: "bg-emerald-100 border-emerald-200 dark:bg-emerald-500/10 dark:border-emerald-500/20",
  blue: "bg-sky-100 border-sky-200 dark:bg-sky-500/10 dark:border-sky-500/20",
  pink: "bg-pink-100 border-pink-200 dark:bg-pink-500/10 dark:border-pink-500/20",
};
```

The page itself is a filter and two sections:

```tsx
export default function NotesPage() {
  const search = useKeepSearch();
  const { pinned, others, active, isLoading } = useKeepViews(search);
  const actions = useNoteActions();
  const create = useCreateNote();
  // ...the compose box, then:
  return (
    <>
      {pinned.length > 0 ? <Section title="Pinned" notes={pinned} /> : null}
      {others.length > 0 ? <Section title="Others" notes={others} /> : null}
    </>
  );
}
```

The masonry is CSS columns, not a grid, because Keep's layout is ragged on
purpose and a grid lines the cards up:

```tsx
<div className="columns-1 gap-4 sm:columns-2 lg:columns-3 xl:columns-4">
  {notes.map((note) => (
    <article key={note.id}
      className={"mb-4 break-inside-avoid rounded-xl border p-4 " + cardColour(note)}>
      ...
    </article>
  ))}
</div>
```

And the one line worth copying into whatever you build next:

```tsx
function patch(note: Note, changes: Partial<Note>) {
  update.mutate({
    id: note.id,
    version: note.version,   // <- this
    ...changes,
  } as { id: string; version: number } & UpdateNoteInput);
}
```

`version` goes out as an `If-Match` header. Two tabs editing the same note do not silently overwrite each other: the second one gets a 409 with the current version. The generated hook makes that argument **required**, so a call site that forgets it does not compile. That is deliberate. Optional, every call site would have gone on sending nothing.

Here it is running:

![The Keep note grid: a rail of pill-shaped links down the left, a search bar across the top, a Pinned section with a yellow card and an Others section with pink, green and blue cards, one carrying a reminder chip](/blog/keep/notes-grid.png)

### Four views, one query

Notes, Reminders, Archive and Trash are the same rows under four filters, so
the app asks for them once:

```typescript
// apps/web/hooks/use-keep.ts, ours
export function useKeepViews(search = "") {
  const { data, isLoading, error } = useNotes({ pageSize: 200, sortBy: "created_at" });

  const all  = (data?.data ?? []).filter(matches(search));
  const live = all.filter((n) => !n.trashed_at);
  const active = live.filter((n) => !n.archived);

  return {
    pinned:    active.filter((n) => n.pinned),
    others:    active.filter((n) => !n.pinned),
    reminders: live.filter((n) => Boolean(n.remind_at)),
    archived:  live.filter((n) => n.archived),
    trashed:   all.filter((n) => Boolean(n.trashed_at)),
    isLoading, error,
  };
}
```

One request, one cache entry, and switching views costs nothing. The division
is in the client because these are small personal collections. A workspace
with ten thousand notes would want the filter in the query, and the generated
hook takes one: `useNotes({ filters: { archived: true } })`.

![The Reminders view: the same rail with Reminders active, showing the one note that has a reminder with its date chip](/blog/keep/reminders.png)

### Trash is a column, not a soft delete

This is the decision I would defend hardest in the whole build.

Grit has a Trash. It is good: soft-deleted rows, a restore endpoint, a
nightly purge past the retention window, and an admin screen at
`/admin/trash`. It is also **an operator tool**, mounted behind
`RequireRole("ADMIN")`, and the thing Keep calls Trash is a product feature a
customer uses every day.

So the note carries `trashed_at`, and the four actions are four writes:

```typescript
trash:    (note) => patch(note, { trashed_at: new Date().toISOString() }),
restore:  (note) => patch(note, { trashed_at: null }),
archive:  (note) => patch(note, { archived: true }),
remind:   (note, at) => patch(note, { remind_at: at }),
```

Only the Trash view can really delete, and that calls `useDeleteNote`, which
is the soft delete, which puts it in the operator's trash. Two layers, each
answering its own question: the customer's "can I get it back", and the
operator's "what was deleted from this account and when".

![The Trash view with one note in it, and a line saying notes in the trash can be restored while deleting one here is permanent](/blog/keep/trash.png)

### The one line worth copying

```tsx
function patch(note: Note, changes: Record<string, unknown>) {
  update.mutate({
    id: note.id,
    version: note.version,   // <- this
    ...changes,
  } as { id: string; version: number } & UpdateNoteInput);
}
```

`version` goes out as an `If-Match` header. Two tabs editing the same note do
not silently overwrite each other: the second gets a 409 with the current
version. The generated hook makes that argument **required**, so a call site
that forgets it does not compile. Optional, every call site would have gone
on sending nothing.

### The shell

Headless writes no sidebar entry, because it wrote no page to link to. The
rail, the search box and the pill-shaped active row are ours, in
`app/(app)/layout.tsx`:

```tsx
const PRODUCT = [
  { href: "/notes",     label: "Notes",     icon: Lightbulb },
  { href: "/reminders", label: "Reminders", icon: Bell },
  { href: "/labels",    label: "Labels",    icon: Tag },
  { href: "/archive",   label: "Archive",   icon: Archive },
  { href: "/trash",     label: "Trash",     icon: Trash2 },
  // grit:nav     <- generated entries land here
];

const ACCOUNT = [
  { href: "/organization",   label: "Organisation", icon: Building2 },
  { href: "/billing",        label: "Billing",      icon: CreditCard },
  { href: "/account",        label: "Overview",     icon: LayoutDashboard },
  { href: "/account/profile", label: "Profile",     icon: UserIcon },
];
```

Two things survived the redesign untouched, and they are the two worth
keeping: `ProtectedWebRoute`, which is the guard, and the `// grit:nav`
marker, which is how `--portal screens` added Labels. Everything else, the
rail, the search bar, the lightbulb, is a hundred lines of our own JSX
wrapped around them.

## The gotcha, which is the good kind

The first time I loaded `/notes` after signing up, every request came back 400. The API log:

```
[400] GET /api/v1/notes?page=1&page_size=100&sort_by=created_at&sort_order=desc
```

That is `--tenant-owned` doing its job. A note belongs to an organisation and Ada did not have one yet, so there was no scope to read within and the middleware refused rather than guessing. Create the organisation, come back, 200.

It is worth sitting with that for a second, because the alternative is worse in a way that is hard to see. A framework that answered "no notes" to a request with no organisation would look friendlier and would be lying: the same query from a user who *does* belong somewhere would be scoped, and the one from a user who does not would quietly read the whole table the day somebody removed the filter. Refusing is the behaviour you want from the layer that stops one customer reading another's data.

## Billing, which is mostly not your problem

![The billing screen: You are not subscribed, and a line explaining that no plan is configured yet and to set NEXT_PUBLIC_STRIPE_PRICE_ID](/blog/keep/billing.png)

That screen is generated by the Stripe plugin, and it is honest about the state it is in. Create a recurring price in Stripe, put the id in `apps/web/.env.local`, and the same screen grows a Subscribe button; once somebody subscribes it shows the plan, whether it is paid up, when it renews, and cancel and resume.

What it deliberately does not have is an invoice table. Stripe's billing portal already lists every invoice, lets a customer download a PDF and lets them change the card, and it is always right. A second copy of that in your app is a copy that drifts.

One detail that matters more than it looks: with the multitenant plugin installed, **billing belongs to the organisation**. One Stripe customer per workspace. The plan survives the person who subscribed leaving, and another admin can change the card. Without the plugin nothing changes, because the organisation id is simply never set.

## What Grit did not do

The note grid. The compose box. The pinning, the reminder chips, the archive
and trash flows, the four views, the rail. That is the product, and it is
about four hundred lines of TSX across six files that nobody could have
generated for me, written against an API, a type, five hooks and a descriptor
that were all generated.

That is the split I want from a framework, and it is the one the three levels exist to let you choose per resource. Labels got `screens` and cost nothing. Notes got `headless` and cost an afternoon. Billing, members and invitations cost a flag.

## Try it

```bash
grit update
grit new keep --saas
```

Answer the three questions, generate a resource with `--owned-by user`, and look at `/your-resource` in the web app rather than only in the admin. If the screen that comes out is the one you wanted, keep it. If it is not, regenerate with `--portal headless` and write the one you wanted, with the data layer already there.

The docs for all of it are at [the customer portal](/docs/frontend/portal), and the full changelog entry is [v3.584.0](/docs/changelog#v3.584.0).
