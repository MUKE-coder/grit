#import "../template.typ": *

= grit new: what lands, and the shape of a project

One command, and the answers to most of chapter twenty six are already in your
directory. This chapter is a tour of what arrives, so that nothing in part five
is a surprise.

== Installing

```bash
curl -fsSL https://gritframework.dev/install.sh | sh    # macOS, Linux
# Windows PowerShell
irm https://gritframework.dev/install.ps1 | iex
go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
```

```bash
grit version
```

== Creating a project

```bash
grit new clinic
```

It asks a few questions, or you can answer them up front:

```bash
grit new clinic --api              # API only, no frontend
grit new clinic --architecture triple
```

The *architecture* decides how many frontends you get. It is a word this
framework uses in a specific way, and it is worth pinning down now because it
is not the same thing as an architecture style:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`triple`], [Web app, admin panel and API, as three applications],
  [`double`], [Web app and API, with admin as role-protected routes in the web app],
  [`single`], [One binary serving an embedded SPA],
  [`api`], [The Go API alone],
  [`mobile`], [API plus an Expo application],
)

#gotcha[
  In Grit, *architecture* means that choice and nothing else. The arrangement
  of the code inside, the modular monolith with its layers, is called the
  architecture *style*, and the two are documented separately for exactly this
  reason. When someone asks "which architecture does Grit use", find out which
  question they are asking.
]

== What arrives

```
clinic/
├── grit.json                 the project's own description
├── docker-compose.yml        Postgres, Redis, MinIO, Mailhog
├── docker-compose.prod.yml
├── apps/
│   ├── api/                  the Go backend
│   ├── web/                  the public Next.js app
│   └── admin/                the admin panel
├── packages/
│   └── shared/               Zod schemas and TypeScript types
└── tests/
```

The Go side is where this book's attention goes:

```
apps/api/
├── go.mod
├── Dockerfile
├── cmd/
│   ├── server/               the API
│   ├── migrate/              schema changes
│   ├── seed/                 realistic data
│   ├── backup/               and
│   └── restore/              the half that matters
└── internal/                 about fifty packages
```

Five commands rather than one. `cmd/restore` existing beside `cmd/backup` is
B5 from the last chapter answered honestly: a backup you have never restored is
a rumour, so the restore is a first-class command rather than a runbook.

== The internal packages

Around fifty, and the names tell you what has already been decided for you.
Grouped by what they answer:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*Your layers*], [`models`, `services`, `handlers`, `routes`, `policies`],
  [*Identity*], [`authz`, `password`, `totp`, `access`, `crypto`],
  [*Data*], [`database`, `migrate`, `paginate`, `concurrency`, `money`, `ids`],
  [*Work*], [`jobs`, `cron`, `outbox`, `saga`, `workflow`, `events`],
  [*Outside*], [`storage`, `files`, `media`, `mail`, `sms`, `webhooks`, `safefetch`],
  [*Operations*], [`audit`, `appendonly`, `backup`, `health`, `tracing`, `flags`],
  [*Compliance*], [`erasure`, `gdpr`-adjacent pieces, `settings`],
  [*Delivery*], [`respond`, `sanitize`, `imports`, `export`, `pdf`, `realtime`, `sync`],
)

You will not read all of these. The point of listing them is that each one is a
thing you did not write, and most of them correspond to a lettered problem from
chapter twenty six.

== Running it

```bash
cd clinic
grit migrate        # create the schema
grit seed           # realistic data to develop against
grit start          # everything, together
```

`grit start` brings up the API, the frontends and the infrastructure at once,
which is D2 gone. What is now listening:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`:8080`], [The Go API],
  [`:3000`], [The web app],
  [`:3001`], [The admin panel],
  [`/studio`], [A database browser: tables, queries, the ERD],
  [`/docs`], [The API reference, generated from the registered routes],
  [`/pulse`], [Latency, query counts, profiles],
  [`/sentinel`], [Errors, with the request attached],
)

The second group is the interesting one. `/studio` is D6, `/docs` is D7 and
`/pulse` is D5, and all three are paths that were already there rather than
tools you installed.

#bridge[
  `/docs` deserves a note because it is the kind of thing that is usually a
  lie. It is assembled from the routes actually registered at runtime, not from
  comments or a hand-written specification. It cannot drift, because there is
  no second description to drift from. That is the only way an API document
  stays true, and it is why Grit generates one rather than asking you to
  maintain one.
]

== grit.json, and the seeding that matters

`grit.json` records what the project is: its architecture, its module
toggles, its name. The CLI reads it so that later commands know the shape of
the thing they are editing.

The seed is worth insisting on, because it is the habit people skip:

```bash
grit seed Patient --count 50000
```

Fifty thousand rows is where you discover whether the list page works. At ten
rows the sort is instant, pagination has one page, and an N+1 costs ten queries
instead of ten thousand. That is D9, and the seeder is built as the pipeline
from chapter fourteen: rows built in chunks, written in batches, with the
number in flight bounded by a channel, so a million rows runs in about a minute
with the heap flat.

== Modules you can switch off

Eleven subsystems are toggles rather than permanent fixtures:

```bash
MODULE_AI=false
MODULE_AUDIT=false
MODULE_BACKUP=false
MODULE_CRON=false
MODULE_FILES=false
MODULE_FLAGS=false
MODULE_JOBS=false
MODULE_MAIL=false
MODULE_REALTIME=false
MODULE_TWOFACTOR=false
MODULE_WEBHOOKS=false
```

Set one to `false` in `.env` and that subsystem does not start. Delete its
package and the rest still builds. This is the "batteries included, individually
removable" claim, and it is checkable: the packages are separate and the seams
are real.

#gotcha[
  The default project is the *hardened* one. Auth, sessions, the audit chain,
  backups and the upload pipeline are present before you write a line, rather
  than being plugins you remember to enable.

  That is deliberate and it is worth stating as a principle: an opt-in security
  feature is a security feature most projects do not have. The cost is a
  bigger project on day one, which is a real cost and the trade this framework
  makes on purpose.
]

== It is still just Go

Open `apps/api/internal/services/user.go`. It is a struct with a `*gorm.DB`
field and methods on it. Open a handler: it reads a request, calls one service
method and writes a response. There is no framework type in your code and
nothing to learn before you can read it.

```go
// The shape of every generated service.
type ContactService struct {
    DB *gorm.DB
}
```

That is the whole dependency. If Grit disappeared tonight, this project
compiles tomorrow, because there is nothing of Grit inside it. Everything you
learned in parts one to three applies directly, which was the reason for
putting those parts first.

#project(title: "Project 24: make a project and look around")[
  No code to write. Make one and read it, because the rest of part four assumes
  you have seen these files.

  ```bash
  grit new bookshop --api
  cd bookshop
  grit migrate
  grit seed
  grit start
  ```

  Then answer these by reading, not by searching the documentation.

  *One.* Open `apps/api/internal/handlers/user.go` and find where it calls the
  service. How many lines of that handler are about HTTP, and how many are
  about users? That ratio is the thin-handler rule from chapter eleven,
  enforced by the generator.

  *Two.* Open `apps/api/internal/services/user.go` and find `paginate.Config`.
  It lists which columns a client may search, sort and filter by. Why is that a
  whitelist rather than passing the client's column name through? The answer is
  in chapter twenty two, under placeholders.

  *Three.* Visit `/studio` and look at the tables. Then visit `/docs` and find
  the route for listing users. Compare the parameters it documents with the
  `paginate.Config` you just read: they are the same list, because one is
  derived from the other.

  *Four.* Count the entries in `apps/api/go.mod`. Then count what a comparable
  Express project needs in `package.json`. Neither number is a virtue on its
  own, and the comparison is worth seeing.

  *Five.* Set `MODULE_MAIL=false` in `.env`, restart, and confirm the API still
  starts. Then look at what disappeared from the startup log.
]

#tryit[
  Run `grit new` twice with the same name into different directories and
  compare the two trees byte for byte. They should be identical: the generator
  is deterministic, which is the property chapter twenty six said a model
  cannot give you. Then run `grit version` and note it, because the next
  chapters describe this version and the framework moves quickly.
]
