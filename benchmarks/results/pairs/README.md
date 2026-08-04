# Measured pairs

Each directory is one head-to-head run: Grit against a single framework, back to
back, three repetitions per scenario, zero failed requests.

| pair | runs | notes |
|---|---|---|
| `bun/` | 24/24 | complete |
| `encore/` | 24/24 | complete |
| `express/` | 23/24 | `express mixed` has two repetitions, not three — the health check refused the third rather than measure against an app that was not answering |

`nextjs`, `django` and `laravel` are not here yet. Next.js was interrupted when
the Docker daemon started returning 500s mid-run; Django and Laravel have not
been run as dedicated pairs.

Regenerate the published table with:

```bash
python pair-report.py bun encore express
```

## Reading these honestly

The absolute figures are only meaningful beside the Grit baseline measured in
the same run. Grit's single-row read measured 6,600 req/s in the Bun pair, 4,392
in the Encore pair and 1,635 in the Express pair — from an identical binary, as
hours of write scenarios accumulated in Postgres. **Ratios compare within a
pair; absolutes do not compare across pairs.**
