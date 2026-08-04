# Grit vs Laravel — benchmark harness

Two apps, the same CRUD resource, the same database, the same load script.
Everything here exists so the result is reproducible and so neither side is
quietly handicapped.

## What is being compared

A public `Product` CRUD API — no auth on either side, because with a token in
play part of what you measure is JWT parsing rather than the request path.

| | Grit | Laravel |
|---|---|---|
| Version | see `grit --version` | 13.x |
| Runtime | Go 1.26, single binary | PHP 8.4, nginx + php-fpm |
| ORM | GORM | Eloquent |
| Built by | `grit new` + `grit generate resource Product` | `composer create-project` + hand-written controller |
| Mode | `GIN_MODE=release` | `APP_ENV=production`, `APP_DEBUG=false`, `artisan optimize` |
| Opcache | n/a | on, JIT tracing, `validate_timestamps=0` |

The Laravel controller is written to match Grit's generated handler rather than
to be idiomatic-at-any-cost: same default page size (20) and cap (100), same
searchable columns, same sortable allow-list, same `{data, meta}` envelope, same
`{error:{code,message}}` shape, and the same `version` bump on update. Plain
Eloquent, not API Resources — Resources would add a transformation layer Grit's
handler has no equivalent of, and that cost would show up as a Laravel tax that
is really a difference in what the two are doing.

The schemas are column-for-column identical, including the index on
`deleted_at`. A benchmark where one side has an index the other lacks measures
the index.

## Fairness rules

- **One Postgres**, tuned once, shared. `fsync=off` and `synchronous_commit=off`
  so disk does not dominate. Separate databases (`bench_grit`, `bench_laravel`)
  only because both frameworks want a `users` table.
- **Identical container limits**: 4 CPUs, 2 GB each. Postgres gets 8 CPUs —
  deliberately more than either app — so the database is not the ceiling.
- **Identical data**: 10,000 rows generated from a fixed seed, so both sides get
  the same UUIDs, names and prices. The `show` scenario reads the same 500 ids
  on both sides.
- **One app runs at a time.** The other is stopped, not idle.
- **Warm-up is discarded.** The first requests pay for Grit's connection pool
  filling and Laravel's opcache/JIT warming.
- **The dataset is reset before every single run.** This one bit hard: the write
  scenario inserts permanently, so `list` and `mixed` end up counting a table
  that grew during the test. Left unchecked it reached 345,680 rows on the Grit
  side against 30,255 on Laravel's — and because Grit writes faster it polluted
  its own table harder and then paid for it on every `COUNT(*)`. That reads as
  "Grit is slow at listing" when it is really "Grit is fast at inserting".
  Every run now starts from the same 10,000 rows.
- **Three repetitions, median reported.** Not the best run — publishing a best
  run is how a benchmark becomes unreproducible.
- **Redis is running.** Not for Laravel's benefit (it uses `QUEUE_CONNECTION=sync`)
  but because Grit's asynq worker has no env switch to disable it, and without a
  Redis to reach it retries in a tight loop and burns CPU that has nothing to do
  with serving requests.

### k6 runs inside the Docker network

This one changes the numbers materially. On Docker Desktop for Windows a
published port goes through a userland proxy in the VM, and that proxy becomes
the ceiling before either framework does:

| | through the published port | on the container network |
|---|---|---|
| Grit, single-row read | 432 req/s, 3.94 ms floor | 740 req/s, 1.11 ms floor |

Everything in `results/` is container-to-container.

## Running it

```bash
docker compose up -d postgres redis
docker compose exec -T postgres psql -U bench -d bench \
  -c "CREATE DATABASE bench_grit;" -c "CREATE DATABASE bench_laravel;"

# schemas
(cd grit-bench && grit migrate)
docker compose up -d laravel && docker compose exec -T laravel php artisan migrate --force
docker compose exec -T laravel php artisan optimize

# identical rows on both sides
python - <<'PY'
# see the generator inline in this repo's history, or reuse seed/products.sql
PY
for db in bench_grit bench_laravel; do
  docker compose exec -T postgres psql -q -U bench -d $db < seed/products.sql
done

./run.sh              # 3 reps x 4 scenarios x 2 apps
python aggregate.py   # medians + which rows are DB-bound
```

## Scenarios

| Scenario | Request | What it exercises |
|---|---|---|
| `list` | `GET /products?page=N&page_size=20` | routing, ORM hydration of 20 rows, JSON encode, plus a `COUNT(*)` |
| `show` | `GET /products/:id` | routing, one indexed lookup, JSON encode — the cleanest framework-overhead signal |
| `write` | `POST /products` | body parse, validation, insert |
| `mixed` | 85% list / 10% show / 5% write | something closer to a real read-heavy API |

## Reading the results honestly

`aggregate.py` prints the CPU of the app container and of Postgres for each row,
and labels each scenario **app-bound** or **DB-bound**.

That label is the important part. A row is only a true ceiling for a framework
when that framework's container is the thing that saturated. Where Postgres
saturated while the app container still had headroom, the number is a **floor** —
the framework could go faster given a bigger database, and the gap between the
two is wider than the row shows.

Do not quote a DB-bound row as "framework X does N req/s". Quote it as "at least
N req/s, with the database as the limit".
