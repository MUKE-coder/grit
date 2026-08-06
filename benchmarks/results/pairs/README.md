# Measured pairs

Each directory is one head-to-head run: Grit against a single framework, back to
back, three repetitions per scenario, zero failed requests.

All six were re-measured on v3.134.0, which cut a generated write from seven
statements to one. Every earlier set of figures is withdrawn; they measured a
Grit that wrapped a single INSERT in a transaction, re-read the row it had just
written, and wrote an audit row for requests with no authenticated user.

| pair | runs | notes |
|---|---|---|
| `bun/` | 24/24 | complete |
| `encore/` | 24/24 | complete |
| `express/` | 24/24 | complete |
| `nextjs/` | 24/24 | complete |
| `django/` | 24/24 | complete |
| `laravel/` | 24/24 | complete; the app itself was fixed earlier, see the guide page |

Regenerate the published table with:

```bash
python pair-report.py bun encore express nextjs django laravel
```

## Reading these honestly

The absolute figures are only meaningful beside the Grit baseline measured in
the same run. Even with all six run in one sitting, Grit's single-row read
ranged from 4,536 req/s in the Bun pair to 8,509 in the Express pair, from an
identical binary, because the machine drifts as write scenarios accumulate in
Postgres. **Ratios compare within a pair; absolutes do not compare across
pairs.**
