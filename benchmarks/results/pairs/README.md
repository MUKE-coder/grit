# Measured pairs

Each directory is one head-to-head run: Grit against a single framework, back to
back, three repetitions per scenario, zero failed requests.

| pair | runs | notes |
|---|---|---|
| `bun/` | 24/24 | complete |
| `encore/` | 24/24 | complete |
| `express/` | 23/24 | `express mixed` has two repetitions, not three — the health check refused the third rather than measure against an app that was not answering |
| `nextjs/` | 24/24 | complete; run a day later than the rest, after the Docker daemon that killed the first attempt was restarted |
| `django/` | 24/24 | complete |

`laravel` is not here. Its figures were measured before this harness existed and
survive only in the published table, which is why its rows carry an explicit
ceiling classification rather than raw Postgres CPU percentages.

Regenerate the published table with:

```bash
python pair-report.py bun encore express nextjs django
```

## Reading these honestly

The absolute figures are only meaningful beside the Grit baseline measured in
the same run. Grit's single-row read measured 6,600 req/s in the Bun pair, 4,651
in the Django pair, 4,392 in the Encore pair, 1,911 in the Next.js pair and
1,635 in the Express pair — from an identical binary. The machine drifts as write
scenarios accumulate in Postgres across a session, and recovers after a restart.
**Ratios compare within a pair; absolutes do not compare across pairs.**
