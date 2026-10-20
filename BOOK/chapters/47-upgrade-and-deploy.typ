#import "../template.typ": *

= Upgrading, deploying, and handing it over

Two things decide whether a framework is a help or a trap, and neither is
visible on the day you start: can you take a new version without a rewrite, and
can somebody else run this when you are not available.

This chapter closes Part V with both, plus the handover that the last project
of your career will thank you for.

== The upgrade problem, stated honestly

A framework that generates code into your project has a problem a library does
not: it wrote files, you edited them, and the next version wants to write them
again.

There are three possible answers, and two of them are bad.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Overwrite], [Your edits are gone. Nobody upgrades twice],
  [Never touch], [Framework fixes never reach existing projects. Everybody is
    stranded on the version they started with],
  [Know which], [Track what was generated, leave what changed],
)

```bash
grit upgrade
```

Files Grit wrote and you have not touched are regenerated. Files you have
edited are *left alone and reported*, because the upgrade knows the difference:
each generated file was recorded when it was written, and a file whose content
no longer matches is yours.

```bash
grit upgrade --diff    # show what it would change in the files it skipped
grit upgrade --plan    # write UPGRADE-PLAN.md, a runbook for those files
grit upgrade --force   # take everything, including your edits
```

#gotcha[
  `--force` is a real option and it is occasionally right, for instance when
  you edited a file by accident and want it back.

  Commit first. Then `--force` is a diff you can read and revert, rather than a
  decision you cannot undo. This is true of every generator command that writes
  into your project, and it is worth a habit rather than a rule.
]

== What an upgrade actually does

Three things, and the second is the one people do not expect.

*It regenerates framework-owned files.* The admin panel's shared components,
the configs, the web app shell.

*It repairs generated code it does not own.* Your services are yours, and
`upgrade` does not rewrite them. But a framework fix that only the *generator*
carries reaches new resources and leaves every existing one broken, so there
are targeted repairs that go and fix one specific thing in files the upgrade
otherwise leaves alone. "A record with a price cannot be saved" is not
something to leave broken on existing projects because of a rule about
ownership.

*It upgrades installed plugins.* Left to a command of its own, nobody runs it,
and a plugin fix reaches nobody.

#gotcha[
  A repair finds its target by matching an anchor in your file. If your version
  of that file has drifted far enough, the anchor is not found, and the repair
  warns and returns rather than guessing.

  So *read the warnings from an upgrade*. They are the list of places where the
  framework wanted to fix something and could not reach it, and each one is a
  manual edit somebody has to make.
]

== The CLI updates itself

```bash
grit update
```

Releases follow channels: a plain version is stable, and a version with an
`-rc`, `-beta` or `-alpha` suffix is not. `grit update` stays on stable unless
you ask otherwise, which is what you want on the machine you work from.

Upgrade the CLI and the project separately, in that order, and *test a fresh
project after upgrading the CLI*. A scaffold template change only shows up in
something newly generated.

== Deployment, in the order of effort

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [One binary], [`go build`, copy, run under systemd. Chapter twenty five],
  [`grit deploy`], [Build, upload over SSH, configure systemd, optionally
    Caddy with automatic TLS],
  [`grit deploy --railway`], [Link the project, push the variables from
    `.env`, upload the API, generate a URL],
  [Docker Compose], [`docker-compose.prod.yml`, for a single machine with
    Postgres and Redis beside it],
  [A platform], [Any host that runs a container. The API is one image],
)

Start at the top. A Go API on one server under systemd, with Caddy in front for
TLS, serves a great deal more traffic than most applications ever see, and it
is a stack you can hold in your head at three in the morning.

#bridge[
  Coming from Next.js, the thing that takes adjusting is that there is no
  build output to upload and no Node runtime to match. One file, built on your
  machine, copied to the server.

  The frontends are still Next.js and still need their own hosting or their own
  container. The simplification is on the API side, and it is the half that
  used to need the most care.
]

== Two deployment bugs worth not repeating

Both of these were real, and both have the same shape: a path guessed from the
wrong directory.

*The module is not where you are standing.* In a monorepo the Go module is in
`apps/api`, and reading `go.mod` from the repository root finds nothing. The
deploy announced itself, named the app `grit-app`, and then failed on "go.mod
file not found".

*A module's base name is not the application's name.* The module is
`invoices/apps/api`, whose base is `api`, so the systemd unit, the remote
directory and the binary were all called `api`. Two Grit projects on one server
would have overwritten each other.

#gotcha[
  The general lesson for any deployment script you write: *derive names from
  something that identifies the project*, and print every derived value before
  acting on it. A deploy that says "deploying api to /srv/api" gives somebody
  the chance to stop it; one that says "deploying" does not.
]

== The checklist before a first deploy

Run these and read the output rather than assuming:

```bash
grit doctor        # configuration, dependencies, common mistakes
grit routes        # every route, and which group it is on
grit test          # all suites
grit backup        # and then actually restore it somewhere
```

Then the environment:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`APP_URL`], [Correct, and matching the real port. Uploads build their URLs
    from it],
  [`JWT_SECRET`], [Generated for this environment, not copied from the
    example file],
  [Database], [Postgres, with a password that was not in a repository],
  [Sentinel storage], [Pointed at the real database, not the SQLite file that
    a redeploy deletes],
  [`MAIL_MAILER`], [A real driver, configured, and tested by sending one],
  [Backups], [Scheduled, and restored once to prove they work],
)

#gotcha[
  The one people skip is the restore. A backup that has never been restored is
  a file, not a backup, and the first restore attempt is not the moment to
  discover it is missing the uploads or the schema.

  Restore into a scratch database, open the application against it, and look at
  a record. Then you have a backup.
]

== Handing it over

The last thing a project needs is for somebody else to be able to run it. Four
files, none of them long:

*`README.md`*: what this is, the three commands to run it locally, and where
the environment variables come from.

*`ARCHITECTURE.md`*: the request path from route to query, and where rules
live. Chapter twenty eight's assignment, kept up to date.

*`RUNBOOK.md`*: how to deploy, how to roll back, how to take and restore a
backup, what to check when it is slow, who to call.

*`DECISIONS.md`*: why things are the way they are. The sort reason, the
currency choice, the reason that endpoint is not cached. Six lines each, and it
is the file that stops the next person undoing your work for reasons you
already considered.

#bridge[
  Writing these feels like overhead when you are the only person who will ever
  touch the project. You are not the only person: you in eight months is a
  different person with none of the context and all of the confidence.

  The test for whether they are good enough is simple. Hand somebody the
  repository, say nothing, and watch where they stop.
]

== Where to go next

You have the language, the standard library, the framework and five
applications. Three directions from here, and they are not equal.

*Deeper into Go.* Part VI builds things from scratch: a parser, a key-value
store, a Redis, a Git. That is where the language stops being a tool you use
and becomes one you understand.

*Deeper into the domain.* Pick the application in Part V that was closest to
real work and make it genuinely production ready for one real user. The last
ten percent is where the learning is.

*Deeper into operations.* Run something for six months. Watch it break, fix
it, and write down what you changed. Nothing in a book teaches what a pager
teaches.

#project(title: "Project 44: upgrade, deploy, hand over")[
  Take the application you are proudest of from Part V and put it in front of
  somebody.

  *One: upgrade.* Commit, run `grit upgrade --diff`, read every file it
  proposes to change, then upgrade. Record in `UPGRADE.md` which files it
  skipped and why each one was edited.

  *Two: the warnings.* If the upgrade warned about a repair it could not
  apply, apply it by hand and note what the anchor drift was.

  *Three: deploy.* To a real server or a real platform, with TLS and a domain.
  Not localhost. Write down every step you took that was not in a command.

  *Four: break it on purpose.* Stop the database in production, watch
  `/health`, and confirm the application fails the way you expected. Then fix
  it and time the recovery.

  *Five: restore.* Take a backup, restore it into a second environment, and
  open the application against the restored data. Confirm a record you created
  this week is there.

  *Six: hand it over.* Give the repository to somebody who has never seen it,
  with no verbal explanation, and ask them to get it running locally and deploy
  a change. Write down every question they asked. Those questions are the
  missing documentation, and nothing else will find them.
]

#tryit[
  Read `UPGRADE-PLAN.md` after running `grit upgrade --plan` on a project you
  have customised heavily.

  It is a list of the places where you and the framework have diverged, which
  is also a map of where your project is most itself. Some of those divergences
  are good decisions worth keeping, and some are things you did before you
  knew the framework had an answer.
]

#assignment(name: "Deployment - Assignment One")[
  In `part5-building/47-deploy`, document a real deployment of a real
  application.

  *One.* `RUNBOOK.md` with: deploy, roll back, backup, restore, what to check
  when it is slow, and the first three things to try during an outage.

  *Two.* `DECISIONS.md` with at least eight decisions and their reasons, in six
  lines each.

  *Three.* Evidence of a restore: the backup, the restored environment, and a
  record visible in both.

  *Four.* The upgrade record: which files were skipped, which repairs warned,
  what you did about each.

  *Five.* The handover test. One person, no explanation, and the list of
  questions they asked. Then fix the documentation so the next person asks
  none of them, and say which question was the most embarrassing.
]
