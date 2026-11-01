# grit setup

`grit setup` is the guided installer. It creates a new Grit project, asks for a
starter and the services the project will use, then runs the steps that make
the result runnable: dependencies, migrations, sample data, git, and optionally
`grit start`.

It is separate from `grit new`. `grit new` is the fast, scriptable scaffold and
its behaviour is frozen. `grit setup` is the interactive path, and it only
creates new projects. Use `grit upgrade` for a project that already exists.

## Quick start

```bash
grit setup myapp
```

The wizard asks for:

- the architecture: single, double, triple, desktop, api, mobile or full
- the frontend: Next.js or TanStack Router, when the shape has one
- the starter: blank or blog
- the database: Postgres, MySQL, SQLite or in-memory
- Redis, storage and mail: local, external or off
- built-in plugins that match the chosen shape
- whether to run `grit start` at the end

## Non-interactive use

Pass `--yes` and the answers as flags. This is the shape to use in CI and
scripts. A project name is required.

```bash
grit setup myapp --yes \
  --arch triple --frontend next --starter blank \
  --db sqlite --redis off --storage local --mail log \
  --no-docker --no-git --no-start
```

Do not pass `--yes` in a terminal you expect to answer questions.

## Starters

A starter is the example content the new project begins with.

| Starter | What it does |
|---|---|
| `blank` | The same scaffold with the built-in Blog example removed. |
| `blog` | The Blog example kept exactly as `grit new` writes it. |

`api` and `mobile` projects always use `blank`, because they have no web or
admin pages for the Blog's screens.

## Services

Redis, storage and mail are configured in `.env`. Each one can point at a local
container, an external provider, or be disabled.

```bash
# Local services, started with docker compose
grit setup myapp --redis local --storage minio --mail mailhog

# External providers
grit setup myapp --redis external --redis-url redis://...
grit setup myapp --storage s3 --storage-endpoint ... --storage-access-key ... --storage-secret-key ...
grit setup myapp --mail resend --mail-api-key ...

# Off
grit setup myapp --redis off --storage local --mail log
```

Credentials are written to `.env`. Setup never prints an external secret back to
the terminal.

## Flags

| Flag | Meaning |
|---|---|
| `--arch` | `single`, `double`, `triple`, `desktop`, `api`, `mobile` or `full` |
| `--frontend` | `next` or `tanstack` |
| `--starter` | `blank` or `blog` |
| `--db` | `postgres`, `mysql`, `sqlite` or `memory` |
| `--theme` | `atlas`, `aurora`, `pulse`, `coral`, `amber`, `sky`, `mono` or `emerald` |
| `--redis` | `local`, `external` or `off` |
| `--redis-url` | External Redis URL |
| `--storage` | `minio`, `s3`, `r2`, `local` or `off` |
| `--storage-*` | S3 or R2 endpoint, keys, bucket, region and public URL |
| `--mail` | `mailhog`, `log`, `resend`, `mailgun`, `postmark`, `sendgrid`, `ses` or `smtp` |
| `--mail-*` | From address, API key, domain, endpoint and stream |
| `--smtp-*` | SMTP host, port, username, password and encryption |
| `--ses-region`, `--aws-*` | Amazon SES region and credentials |
| `--plugin` | Built-in plugin to install; repeat or comma-separate |
| `--i18n` | Add internationalisation |
| `--yes` | Use defaults and flags, with no prompts |
| `--here` | Scaffold into the current directory |
| `--force` | Allow scaffolding into a non-empty directory |
| `--no-install` | Skip `pnpm install` |
| `--no-docker` | Do not start local services |
| `--no-migrate` | Skip migrations |
| `--no-seed` | Skip sample data |
| `--no-git` | Skip git init and the first commit |
| `--star` | Star `MUKE-coder/grit` without asking |
| `--start` | Run `grit start` after setup |
| `--no-start` | Do not run `grit start` |

## What setup runs

1. Scaffolds the project with the chosen architecture, frontend and theme.
2. Applies the starter.
3. Configures `.env` for Redis, storage and mail.
4. Installs compatible built-in plugins.
5. Installs dependencies and tidies the Go module.
6. Starts local services with Docker, unless `--no-docker` is set.
7. Migrates the database and seeds sample data.
8. Initialises git and creates the first commit.
9. Offers to star the Grit repository.
10. Offers to run `grit start`.

Every step after the scaffold is a real command a person could run by hand. If
one fails, setup prints the command to retry.

## Credentials

The final output lists the routes the project exposes. The admin, GORM Studio
and Sentinel rows include the login created for that project:

```text
  GORM Studio   http://localhost:8080/studio       | admin / <generated>
  Sentinel      http://localhost:8080/sentinel/ui  | admin / <generated>
  Admin         http://localhost:3001              | admin@example.com / admin123
```

GORM Studio and Sentinel use generated local passwords. The admin password is
`admin123` in development unless `SEED_ADMIN_PASSWORD` is set. Change these
before production. When `SEED_ADMIN_PASSWORD` is set, setup does not echo the
secret.

If the seed step was skipped, the admin login is not displayed. Run `grit seed`
to create `admin@example.com` with the development password.
