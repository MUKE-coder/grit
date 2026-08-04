#!/usr/bin/env bash
# The published comparison: Grit as it now ships (idle pool == open pool)
# against production Laravel. Same session, back to back — cross-session numbers
# on this host are not comparable, the background load moves.
#
# Samples CPU DURING the run, not after — the first attempt at this sampled
# once k6 had already exited and dutifully reported 0.17%, which says nothing.
set -uo pipefail
cd "$(dirname "$0")"
export MSYS_NO_PATHCONV=1

# Refuse to start if another copy is already running. Two runners stopping each
# other's containers mid-test produced Laravel at 9 req/s and an interleaved
# result set that looked exactly like data until the file timestamps gave it
# away. A lock is cheaper than forensics.
LOCK="$(pwd)/.bench.lock"
if ! mkdir "$LOCK" 2>/dev/null; then
  echo "another run holds $LOCK — stop it first, or rm -rf $LOCK to override" >&2
  exit 1
fi
trap 'rm -rf "$LOCK"' EXIT

REPS="${REPS:-3}"
VUS="${VUS:-50}"
DURATION="${DURATION:-30s}"
SCENARIOS=(list show mixed write)
mkdir -p results/final

for app in grit laravel; do
  other=$([ "$app" = grit ] && echo laravel || echo grit)
  docker compose stop "$other" >/dev/null 2>&1
  docker compose up -d "$app" >/dev/null 2>&1
  for _ in $(seq 1 60); do
    docker run --rm --network benchmarks_default curlimages/curl -sf \
      "http://$app:8080/api/v1/products?page_size=1" >/dev/null 2>&1 && break
    sleep 1
  done

  for scenario in "${SCENARIOS[@]}"; do
    for rep in $(seq 1 "$REPS"); do
      # Reset BOTH tables to the identical 10,000-row seed before every run.
      # Without this the write scenario silently poisons the benchmark: it
      # inserts permanently, so list and mixed end up counting a table that grew
      # during the test. Left unchecked it reached 345,680 rows on the Grit side
      # against 30,255 on Laravel's — and because Grit writes faster it polluted
      # its own table harder, then paid for it on every COUNT. That reads as
      # Grit being slow at listing when it is Grit being fast at inserting.
      for db in bench_grit bench_laravel; do
        docker compose exec -T postgres psql -q -U bench -d "$db" \
          -c "TRUNCATE products;" >/dev/null 2>&1
        docker compose exec -T postgres psql -q -U bench -d "$db" \
          < seed/products.sql >/dev/null 2>&1
      done

      # warm-up, discarded
      docker run --rm --network benchmarks_default -v "$PWD:/bench" -w /bench \
        -e BASE="http://$app:8080" -e SCENARIO="$scenario" -e VUS="$VUS" -e DURATION=8s \
        grafana/k6 run -q k6/bench.js >/dev/null 2>&1

      ( sleep 14
        docker stats --no-stream --format "{{.Name}} {{.CPUPerc}}" 2>/dev/null \
          | grep -E "benchmarks-(grit|laravel|postgres)" \
          > "results/final/$app-$scenario-$rep.cpu" ) &

      docker run --rm --network benchmarks_default -v "$PWD:/bench" -w /bench \
        -e BASE="http://$app:8080" -e SCENARIO="$scenario" -e VUS="$VUS" -e DURATION="$DURATION" \
        grafana/k6 run --summary-export "results/final/$app-$scenario-$rep.json" k6/bench.js \
        > "results/final/$app-$scenario-$rep.log" 2>&1
      wait

      printf "%-12s %-6s r%s  " "$app" "$scenario" "$rep"
      grep -oE "http_reqs[.]+: [0-9]+ +[0-9.]+/s" "results/final/$app-$scenario-$rep.log" | head -1
      sed 's/^/               /' "results/final/$app-$scenario-$rep.cpu" 2>/dev/null
    done
  done
done
echo "done — results/final/"
