# Measured pairs

Each directory is one head-to-head run: Grit against a single framework, back to
back, three repetitions per scenario, zero failed requests.

| pair | runs | notes |
|---|---|---|
| `bun/` | 24/24 | complete; re-measured after v3.133.0 turned GORM's prepared-statement cache on, which the first run predated |
| `encore/` | 24/24 | complete |
| `express/` | 23/24 | `express mixed` has two repetitions, not three — the health check refused the third rather than measure against an app that was not answering |
| `nextjs/` | 24/24 | complete; run a day later than the rest, after the Docker daemon that killed the first attempt was restarted |
| `django/` | 24/24 | complete |
| `laravel/` | 24/24 | complete; re-measured after three setup faults were found — see below |

All six pairs are here. Laravel's original figures were withdrawn: it had been
measured with its dev dependencies in the autoloader (`composer install` without
`--no-dev`), opening a fresh Postgres connection per request while every other
framework pooled, and attempting an SSL handshake on each of those because
Laravel defaults to `sslmode=prefer`. Corrected, its single-row read went from
113 to 175 req/s. The Dockerfile, entrypoint and `config/database.php` in
`laravel-bench/` carry the fixes.

Regenerate the published table with:

```bash
python pair-report.py bun encore express nextjs django laravel
```

## Reading these honestly

The absolute figures are only meaningful beside the Grit baseline measured in
the same run. Grit's single-row read measured 10,655 req/s in the Bun pair, 4,651
in the Django pair, 4,392 in the Encore pair, 4,340 in the Laravel pair, 1,911
in the Next.js pair and 1,635 in the Express pair — from an identical binary. The machine drifts as write
scenarios accumulate in Postgres across a session, and recovers after a restart.
**Ratios compare within a pair; absolutes do not compare across pairs.**
