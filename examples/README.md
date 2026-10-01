# Grit architecture guides — the Job Portal, six ways

**These folders hold guides, not checked-in applications.** Each one is the
exact sequence of commands that builds the same Job Portal in a different Grit
architecture, with the decisions explained and the shape of the result
described. Nothing here is a repository you clone and run.

That is deliberate. Six copies of a generated application, checked in, would
be six copies to keep green against a framework that releases most days, and
the day they fell behind they would be teaching the wrong thing while looking
authoritative. The guides stay correct because the commands in them are checked
against the real command tree on every push, and a renamed flag fails the docs
build.

**If you want an application to read, clone and run, it is [`demo/`](../demo).**
That one is real code, it is the app behind the live demo, and it is built on
every release.

## The six shapes

| Folder | Architecture | What it is |
|---|---|---|
| [`job-portal-api-only`](job-portal-api-only) | `--api` | A Go API and nothing else. Bring your own frontend. |
| [`job-portal-single-vite`](job-portal-single-vite) | `--single --vite` | One Go binary serving a TanStack Router SPA. |
| [`job-portal-double-vite`](job-portal-double-vite) | `--double --vite` | API plus one Vite frontend with the admin inside it. |
| [`job-portal-triple-vite`](job-portal-triple-vite) | `--triple --vite` | API, public TanStack site, separate TanStack admin. |
| [`job-portal-triple-next`](job-portal-triple-next) | `--triple --next` | API, public Next.js site, separate Next.js admin. |
| [`job-portal-mobile-expo`](job-portal-mobile-expo) | `--triple --next --mobile` | The above plus an Expo app sharing the types. |

Each folder has:

- `README.md` — the architecture, and who should pick it
- `GUIDE.md` — the commands in order, from `grit new` to a running app
- `.env.example` — the environment the guide assumes
- `docker-compose.prod.yml` — how that shape deploys

## Features the guides build

The same application each time, so the only variable is the architecture:

- Auth: email and password, OAuth (Google, GitHub), JWT, TOTP 2FA
- Resources: Job, Company, Application, Category
- A dashboard with stat cards
- Data tables with sorting, filtering, pagination and export
- File uploads through presigned URLs (company logos, résumés)
- A generated admin panel with roles and permissions

## Start here

```bash
# Read the shape you want, then follow its GUIDE.md.
# Or let the CLI ask you, which is the same decision with fewer tabs open:
grit new job-portal
```
