/*
 * What is being worked on next, what is waiting, and what Grit has decided
 * not to do.
 *
 * Somebody deciding whether to build on a framework wants three things: who
 * maintains it, what is planned, and how decisions get made. Until this page
 * existed all three had to be inferred from the changelog (grit#246).
 *
 * ── Why this is a snapshot and not a live query ───────────────────────────
 *
 * The site is statically built, so a live GitHub query would mean a client
 * fetch on a marketing page and a spinner where the roadmap should be. The
 * date below says when it was last reviewed, and every item links to its
 * issue, where the real state is. A roadmap is a statement of intent anyway:
 * the thing that goes stale is the intent, and that only a human can refresh.
 *
 * ── The rule for a bucket ─────────────────────────────────────────────────
 *
 * next     work started, or the next thing to start. Small enough to say when.
 * later    agreed it should exist, not scheduled. Most of the list.
 * question open because a decision is owed, usually by the maintainer, and
 *          naming what the decision is so somebody can argue with it.
 *
 * "Not planned" is separate, below, and holds only things already decided
 * against with a reason. An issue nobody has got to is not "not planned",
 * and filing it under that would be using the roadmap to close the backlog.
 */

/** When a human last went through the open issues and re-bucketed them. */
export const ROADMAP_REVIEWED = '2026-11-22'

export type Bucket = 'next' | 'later' | 'question'

export const bucketLabel: Record<Bucket, string> = {
  next: 'Next',
  later: 'Later',
  question: 'Waiting on a decision',
}

export const bucketStyle: Record<Bucket, string> = {
  next: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  later: 'bg-sky-500/10 text-sky-400 border-sky-500/20',
  question: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
}

export const bucketBlurb: Record<Bucket, string> = {
  next: 'Started, or the next thing to start.',
  later:
    'Agreed it should exist, not scheduled. Say so on the issue if you need one sooner: that is how several of these moved up.',
  question:
    'Open because a decision is owed before any code is worth writing. The decision is named, so it can be argued with.',
}

export interface RoadmapItem {
  /** The GitHub issue, which is where the real state lives. */
  issue: number
  /** Said in the project's own words, not the issue title. */
  title: string
  bucket: Bucket
  /** Why it is in this bucket. One sentence. */
  note: string
}

export const ROADMAP: RoadmapItem[] = [
  // ── Next ────────────────────────────────────────────────────────────────
  {
    issue: 255,
    title: 'Fail CI when the shared types drift from the Go models',
    bucket: 'next',
    note: 'The generator already proves generate and sync agree. This turns that into a gate a project can run on its own code.',
  },
  {
    issue: 58,
    title: 'Migration safety: a dry-run plan and a lock-aware advisor',
    bucket: 'next',
    note: 'grit migrate already refuses to roll back a baseline. Seeing what a migration will do before it does it is the missing half.',
  },
  {
    issue: 249,
    title: 'grit eject: remove a battery from a project cleanly',
    bucket: 'next',
    note: 'The inverse of the removal work already done for resources and plugins, which is where the mechanism comes from.',
  },
  {
    issue: 248,
    title: 'What you can delete from a new project, and a --lean preset',
    bucket: 'next',
    note: 'A full project is a lot of files. Saying which are safe to delete is cheap and is the most common question after the first scaffold.',
  },
  {
    issue: 245,
    title: 'A stable tag for teams that cannot take a release every day',
    bucket: 'next',
    note: 'Every merged change ships, which is right for the project and wrong for a team with a change-control process.',
  },
  {
    issue: 266,
    title: 'A monthly digest of the changelog',
    bucket: 'next',
    note: 'The changelog is thorough and nobody reads 754 entries. A monthly summary is generated from the same source.',
  },

  // ── Later ───────────────────────────────────────────────────────────────
  {
    issue: 41,
    title: 'OpenTelemetry: traces, metrics and correlated logs',
    bucket: 'later',
    note: 'Pulse answers the same questions inside one app. OTel is what you need once the app is not the only thing in the request path.',
  },
  {
    issue: 261,
    title: 'A History tab on every admin record page',
    bucket: 'later',
    note: 'The audit log already records the changes. This is the screen that reads it per record.',
  },
  {
    issue: 260,
    title: 'Per-resource data retention',
    bucket: 'later',
    note: 'A deletion policy expressed on the resource, which is where the rest of the resource policy already lives.',
  },
  {
    issue: 251,
    title: 'Ranked full-text search, and hybrid keyword plus vector',
    bucket: 'later',
    note: 'Vector search ships today. Ranked text search and a hybrid of the two is a bigger piece of query work.',
  },
  {
    issue: 252,
    title: 'An :embed modifier that keeps a vector in step with a text field',
    bucket: 'later',
    note: 'Pairs with the search work above, and should land with it rather than before it.',
  },
  {
    issue: 250,
    title: 'Stripe: per-seat and usage-based billing for organizations',
    bucket: 'later',
    note: 'The Stripe plugin handles payments. Subscription billing against the multitenant organization model is its own design.',
  },
  {
    issue: 253,
    title: 'Mobile money: M-Pesa, MTN MoMo, Airtel Money',
    bucket: 'later',
    note: 'The payment rails most of this project’s own users actually need, and each one is a separate integration with a separate sandbox.',
  },
  {
    issue: 224,
    title: 'The rest of the MCP plugin: grit mcp sync, doctor checks, admin screens',
    bucket: 'later',
    note: 'The server, the connections and a generated tool set per resource shipped in v3.548.0. What is left is the resync command, the doctor checks for a tools file that has drifted from its model, and making a connection from the admin rather than over the API.',
  },
  {
    issue: 16,
    title: 'A typed TypeScript SDK generated from the OpenAPI spec',
    bucket: 'later',
    note: 'The monorepo already shares types. This is for the clients that are not in the monorepo.',
  },
  {
    issue: 17,
    title: 'A notifications module: in-app, email and push',
    bucket: 'later',
    note: 'The pieces exist separately. The module is the one place that decides which channel a given notification uses.',
  },
  {
    issue: 59,
    title: 'One-command deploy to Dokploy, Fly and Railway',
    bucket: 'later',
    note: 'grit deploy covers several targets already; the rest is one provider at a time, and each needs testing against a real account.',
  },
  {
    issue: 259,
    title: 'Kubernetes output: manifests or a Helm chart',
    bucket: 'later',
    note: 'Docker Compose and a VPS cover the projects Grit sees. This is for the ones that outgrow that.',
  },
  {
    issue: 254,
    title: 'Preview environments per branch',
    bucket: 'later',
    note: 'A GitHub Action over the deploy work above, so it follows it.',
  },
  {
    issue: 66,
    title: 'grit generate suite: a whole management system in one command',
    bucket: 'later',
    note: 'The largest single item on the list. The blueprints cover the same ground today, by hand.',
  },
  {
    issue: 63,
    title: 'A double-entry accounting plugin',
    bucket: 'later',
    note: 'Money is exact in Grit already, which is the hard part. A ledger on top of it is a real piece of work and a real domain.',
  },
  {
    issue: 64,
    title: 'A hardware plugin: thermal printers, cash drawers, scanners',
    bucket: 'later',
    note: 'The desktop app is where this belongs, and it needs the hardware in the room to test.',
  },
  {
    issue: 65,
    title: 'An MQTT plugin for IoT',
    bucket: 'later',
    note: 'Companion to the websockets plugin. Waiting for somebody with a real device fleet to want it.',
  },
  {
    issue: 52,
    title: 'A pluggable adapter library for Stripe, Twilio, Resend and S3',
    bucket: 'later',
    note: 'Each of those works today behind its own interface. One common shape across them is a refactor, not a feature.',
  },
  {
    issue: 53,
    title: 'grit describe: propose models and handlers from a sentence',
    bucket: 'later',
    note: 'The agent skills do this well already with a human in the loop, which is the part this would remove.',
  },
  {
    issue: 54,
    title: 'An AI reviewer that knows the framework conventions',
    bucket: 'later',
    note: 'grit doctor covers the mechanical half offline. The rest is judgement, and a reviewer that cries wolf gets switched off.',
  },
  {
    issue: 60,
    title: 'A component playground and a design-token surface',
    bucket: 'later',
    note: 'The admin components are copied into each project rather than imported, so a playground has to be generated too.',
  },
  {
    issue: 256,
    title: 'Try grit generate in a browser, with nothing installed',
    bucket: 'later',
    note: 'The generator is Go and compiles to WebAssembly, so this is possible. It is a day of work and a surface to maintain.',
  },
  {
    issue: 258,
    title: 'A server-rendered frontend option: templ and htmx',
    bucket: 'later',
    note: 'A fifth frontend is a lot of surface. Worth it only if enough people want a Go-only stack to justify testing it.',
  },
  {
    issue: 92,
    title: 'Laravel Cloud as a deploy target',
    bucket: 'later',
    note: 'Five of the ten pieces shipped. The start command is still not settable through their API, which blocks the rest.',
  },
  {
    issue: 265,
    title: 'A Built-with-Grit case study page',
    bucket: 'later',
    note: 'Needs real numbers from a real team. Inventing them is the one thing a page like that must not do.',
  },
  {
    issue: 264,
    title: 'A video walkthrough series',
    bucket: 'later',
    note: 'Scripts and chapters are written. The recording is the part nobody else can do.',
  },
  {
    issue: 199,
    title: 'Multitenant follow-ups found on a real project',
    bucket: 'later',
    note: 'A list of small corrections from somebody using invitations in anger, each one its own fix.',
  },
  {
    issue: 179,
    title: 'Ten proposals from a business suite built in ten days',
    bucket: 'later',
    note: 'A field report rather than a feature. Several of its proposals have already shipped under their own issues.',
  },
  // ── Waiting on a decision ───────────────────────────────────────────────
  {
    issue: 47,
    title: 'Rate limits per permission tier, and for API keys',
    bucket: 'question',
    note: 'The design question is whether a tier is a property of the role, of the key, or of the route. A proposal is on the issue and the decision is open.',
  },
]
