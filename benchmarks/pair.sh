#!/usr/bin/env bash
# ONE pair, one invocation: Grit against a single framework.
#
#   ./pair.sh bun
#   REPS=3 DURATION=30s ./pair.sh laravel
#
# Deliberately not a loop over all six. Looping meant a single long-lived
# process, and when one of those had to be killed it did not always die — three
# separate times a survivor kept running, fought the newer run over containers,
# and produced measurements that looked like data. One pair per invocation means
# a run that ends, ends.
#
# Everything is reset before the pair starts: every app container down, every
# database restored to the same 10,000 rows, connections cleared. Nothing is
# inherited from whatever ran before.

set -uo pipefail
cd "$(dirname "$0")"
export MSYS_NO_PATHCONV=1

OPPONENT="${1:-}"
if [ -z "$OPPONENT" ]; then
  echo "usage: ./pair.sh <bun|encore|express|nextjs|django|laravel>" >&2
  exit 2
fi

REPS="${REPS:-3}"
VUS="${VUS:-50}"
DURATION="${DURATION:-30s}"
SCENARIOS="${SCENARIOS:-show write list mixed}"
ALL_APPS="grit bun encore express nextjs django laravel"
OUT="results/pairs/$OPPONENT"

# ── Lock ────────────────────────────────────────────────────────────
# Released only by the process that took it. An unconditional trap meant a
# killed run tore down a lock a newer run already held, and the two then fought.
LOCK="$(pwd)/.bench.lock"
if ! mkdir "$LOCK" 2>/dev/null; then
  echo "another benchmark run holds $LOCK" >&2
  [ -f "$LOCK/owner" ] && echo "  started: $(cat "$LOCK/owner")" >&2
  echo "  is it alive?  ps aux | grep pair.sh" >&2
  echo "  if dead:      rm -rf $LOCK" >&2
  exit 1
fi
date "+%Y-%m-%d %H:%M:%S by pid $$" > "$LOCK/owner"
echo $$ > "$LOCK/pid"
release_lock() { [ "$(cat "$LOCK/pid" 2>/dev/null)" = "$$" ] && rm -rf "$LOCK"; }
trap release_lock EXIT

# ── Helpers ─────────────────────────────────────────────────────────
k6run() { # base scenario duration extra...
  local base="$1" scenario="$2" duration="$3"; shift 3
  docker run --rm --network benchmarks_default \
    -v "$PWD:/bench" -w /bench \
    -e BASE="$base" -e SCENARIO="$scenario" -e VUS="$VUS" -e DURATION="$duration" \
    grafana/k6 run "$@" k6/bench.js
}

healthy() {
  docker run --rm --network benchmarks_default curlimages/curl -sf --max-time 5 \
    "http://$1:8080/api/v1/products?page_size=1" >/dev/null 2>&1
}

seed_db() { # db
  local db="$1" count
  docker compose exec -T postgres psql -q -U bench -d postgres \
    -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity
        WHERE datname = '$db' AND pid <> pg_backend_pid();" >/dev/null 2>&1
  docker compose exec -T postgres psql -q -U bench -d "$db" \
    -c "SET lock_timeout = '10s';" -c "TRUNCATE products;" >/dev/null 2>&1
  docker compose exec -T postgres psql -q -U bench -d "$db" \
    < seed/products.sql >/dev/null 2>&1
  count=$(docker compose exec -T postgres psql -tA -U bench -d "$db" \
    -c "SELECT count(*) FROM products;" 2>/dev/null | tr -d '[:space:]')
  [ "$count" = "10000" ]
}

# ── Full reset, before anything is measured ─────────────────────────
echo "── resetting everything ──"

# Every app down. Only one may be up at a time: seven at 100 pooled connections
# each exceeds Postgres's max_connections and everything fails with "too many
# clients", which looks like a framework problem and is not.
for app in $ALL_APPS; do
  docker compose stop "$app" >/dev/null 2>&1
done
docker compose up -d postgres >/dev/null 2>&1

for _ in $(seq 1 60); do
  docker compose exec -T postgres pg_isready -U bench -d bench >/dev/null 2>&1 && break
  sleep 1
done

# Every database back to the same 10,000 rows — not just the two in this pair.
# The write scenario inserts permanently, and a table left at 300,000 rows makes
# the next COUNT(*) look like a slow framework.
for app in $ALL_APPS; do
  printf "   bench_%-9s " "$app"
  if seed_db "bench_$app"; then echo "10000 rows"; else echo "FAILED"; exit 1; fi
done

# A completed pair is never thrown away — it is moved aside with a timestamp.
# results/final once held a clean Laravel run and was deleted by an rm -rf during
# a reset, which meant redoing work that had already been done properly. Disk is
# cheaper than a re-run.
if [ -d "$OUT" ] && [ -n "$(ls -A "$OUT" 2>/dev/null)" ]; then
  archive="results/archive/$OPPONENT-$(date +%Y%m%d-%H%M%S)"
  mkdir -p "$(dirname "$archive")"
  mv "$OUT" "$archive"
  echo "   previous $OPPONENT results archived to $archive"
fi
mkdir -p "$OUT"

# ── Measure ─────────────────────────────────────────────────────────
measure() { # app scenario
  local app="$1" scenario="$2" rep
  for rep in $(seq 1 "$REPS"); do
    # Reset with the app STOPPED. Terminating its backends is not enough on its
    # own: Bun's SQL client reconnects instantly and re-blocks the exclusive lock
    # TRUNCATE needs, so every rep after the first was refused as a dirty
    # dataset — Bun ended up with one measurement per scenario instead of three.
    # Stopping first removes the contention entirely and costs a few seconds.
    docker compose stop "$app" >/dev/null 2>&1
    seed_db "bench_$app" || { echo "     !! dirty dataset — skipped" >&2; continue; }
    docker compose up -d "$app" >/dev/null 2>&1

    # Health re-checked before every measurement, not once at startup. A run
    # against a stopped container still writes a summary — 100% failed at 0s —
    # and that renders as a plausible-looking low number rather than as nothing.
    local ready=0 i
    for i in $(seq 1 60); do
      healthy "$app" && { ready=1; break; }
      sleep 1
    done
    [ "$ready" = 1 ] || { echo "     !! $app did not come back — skipping" >&2; continue; }

    k6run "http://$app:8080" "$scenario" 8s -q >/dev/null 2>&1   # warm-up, discarded

    ( sleep 14
      docker stats --no-stream --format "{{.Name}} {{.CPUPerc}}" 2>/dev/null \
        | grep -E "benchmarks-($app|postgres)" > "$OUT/$app-$scenario-$rep.cpu" ) &

    k6run "http://$app:8080" "$scenario" "$DURATION" \
      --summary-export "$OUT/$app-$scenario-$rep.json" \
      > "$OUT/$app-$scenario-$rep.log" 2>&1
    wait

    printf "     %-8s %-6s r%s  " "$app" "$scenario" "$rep"
    grep -oE "http_reqs[.]+: [0-9]+ +[0-9.]+/s" "$OUT/$app-$scenario-$rep.log" | head -1
    grep -oE "failed_requests[.]+: [0-9.]+%" "$OUT/$app-$scenario-$rep.log" \
      | grep -qE ": 0.00%" || echo "        ^^ NOT clean — failures above zero" >&2
  done
}

start_only() { # app
  local app="$1" other
  for other in $ALL_APPS; do
    [ "$other" = "$app" ] || docker compose stop "$other" >/dev/null 2>&1
  done
  docker compose up -d "$app" >/dev/null 2>&1
  for _ in $(seq 1 90); do
    healthy "$app" && return 0
    sleep 1
  done
  echo "  !! $app never became ready" >&2
  return 1
}

echo
echo "══ grit vs $OPPONENT — ${REPS} reps of ${DURATION} ══"

for scenario in $SCENARIOS; do
  echo "   -- $scenario"
  start_only grit       && measure grit       "$scenario"
  start_only "$OPPONENT" && measure "$OPPONENT" "$scenario"
done

# Leave nothing running, so the next pair starts from the same place this one did.
for app in $ALL_APPS; do
  docker compose stop "$app" >/dev/null 2>&1
done

echo
echo "done — $OUT"
