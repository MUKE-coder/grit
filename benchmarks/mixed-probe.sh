#!/usr/bin/env bash
# Why does Grit lose the mixed scenario when it wins every part of it?
#
# Hypothesis: mixed is 5% inserts, and the faster writer grows its own table
# faster. list runs COUNT(*) over that table on every request, so within a
# single 30-second run the framework that writes quicker makes its own reads
# progressively more expensive. If true, the row count at the end of a run is
# the tell, and the scenario is measuring the wrong thing.
#
# This runs mixed alone, three times per side, and records the row count each
# run finished with alongside the throughput.
set -u

# Git Bash rewrites /bench into C:/Program Files/Git/bench before docker sees
# it. pair.sh has the same line for the same reason.
export MSYS_NO_PATHCONV=1

REPS=${REPS:-3}
DURATION=${DURATION:-30s}
OUT=results/mixed-probe
mkdir -p "$OUT"

rows() { # db
  docker compose exec -T postgres psql -tA -U bench -d "$1" \
    -c "SELECT count(*) FROM products;" 2>/dev/null | tr -d '[:space:]'
}

seed() { # db
  docker compose exec -T postgres psql -q -U bench -d postgres \
    -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity
        WHERE datname = '$1' AND pid <> pg_backend_pid();" >/dev/null 2>&1
  docker compose exec -T postgres psql -q -U bench -d "$1" \
    -c "SET lock_timeout = '10s';" -c "TRUNCATE products;" >/dev/null 2>&1
  docker compose exec -T postgres psql -q -U bench -d "$1" < seed/products.sql >/dev/null 2>&1
}

k6run() { # base scenario duration
  docker run --rm --network benchmarks_default     -v "$PWD:/bench" -w /bench     -e BASE="$1" -e SCENARIO="$2" -e VUS="${VUS:-50}" -e DURATION="$3"     grafana/k6 run k6/bench.js
}

printf '%-8s %-4s %10s %12s %10s\n' app rep "req/s" "rows after" "grew by"
for app in grit bun; do
  for rep in $(seq 1 "$REPS"); do
    docker compose stop "$app" >/dev/null 2>&1
    seed "bench_$app"
    before=$(rows "bench_$app")
    docker compose up -d "$app" >/dev/null 2>&1
    sleep 8

    log="$OUT/$app-mixed-$rep.log"
    k6run "http://$app:8080" mixed "$DURATION" > "$log" 2>&1

    after=$(rows "bench_$app")
    rps=$(grep -oE "http_reqs[.]+: [0-9]+ +[0-9.]+/s" "$log" | grep -oE "[0-9.]+/s" | tr -d '/s')
    printf '%-8s %-4s %10s %12s %10s\n' "$app" "$rep" "${rps:-?}" "${after:-?}" "$(( ${after:-0} - ${before:-0} ))"
  done
done
