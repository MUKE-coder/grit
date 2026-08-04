#!/usr/bin/env bash
# Grit vs Laravel — full matrix. Writes results/<app>-<scenario>.{json,log,cpu}.
#
# k6 runs INSIDE the docker network, not on the host. On Docker Desktop for
# Windows a published port goes through a userland proxy in the VM, and that
# proxy — not the framework — becomes the ceiling: measured through it, Grit's
# single-row read capped at 432 rps with a 3.94ms floor; on the container
# network the same test does 740 rps with a 1.11ms floor. Everything below runs
# container-to-container so the number belongs to the framework.
#
# One app at a time, the other stopped rather than idle. Each scenario gets a
# warm-up that is discarded — the first requests pay for Grit's connection pool
# filling and Laravel's opcache/JIT warming, and counting those would flatter
# whichever ran second.
set -uo pipefail
cd "$(dirname "$0")"

VUS="${VUS:-50}"
DURATION="${DURATION:-30s}"
WARMUP="${WARMUP:-10s}"
NETWORK="${NETWORK:-benchmarks_default}"
HOST_DIR="${HOST_DIR:-D:/LEARNING/Grit Framework/Benchmarks}"
SCENARIOS=(list show mixed write)
REPS="${REPS:-3}"

mkdir -p results
export MSYS_NO_PATHCONV=1

k6run() { # app scenario duration extra-args...
  local base="$1" scenario="$2" duration="$3"; shift 3
  docker run --rm --network "$NETWORK" \
    -v "$HOST_DIR:/bench" -w /bench \
    -e BASE="$base" -e SCENARIO="$scenario" -e VUS="$VUS" -e DURATION="$duration" \
    grafana/k6 run "$@" k6/bench.js
}

for app in grit laravel; do
  other=$([ "$app" = grit ] && echo laravel || echo grit)
  base="http://$app:8080"

  echo "=== $app ==="
  docker compose stop "$other" >/dev/null 2>&1
  docker compose up -d "$app" >/dev/null 2>&1

  for _ in $(seq 1 60); do
    docker run --rm --network "$NETWORK" curlimages/curl -sf \
      "$base/api/v1/products?page_size=1" >/dev/null 2>&1 && break
    sleep 1
  done

  for scenario in "${SCENARIOS[@]}"; do
   for rep in $(seq 1 "$REPS"); do
    echo "--- $app / $scenario (rep $rep/$REPS)"

    k6run "$base" "$scenario" "$WARMUP" -q >/dev/null 2>&1

    # Sample container CPU mid-run. If the database is pegged while the app is
    # not, that scenario is measuring Postgres and the result understates the
    # faster framework. The write-up has to say so rather than quietly report a
    # number it knows is a floor.
    ( sleep 14; docker stats --no-stream --format "{{.Name}} {{.CPUPerc}}" \
        2>/dev/null | grep benchmarks > "results/$app-$scenario-$rep.cpu" ) &

    k6run "$base" "$scenario" "$DURATION" \
      --summary-export "results/$app-$scenario-$rep.json" \
      > "results/$app-$scenario-$rep.log" 2>&1

    wait
    grep -E "http_reqs\." "results/$app-$scenario-$rep.log" | head -1
    cat "results/$app-$scenario-$rep.cpu" 2>/dev/null | sed 's/^/    cpu: /'
   done
  done
done

echo
echo "done — results/ has raw json, logs and cpu samples"
