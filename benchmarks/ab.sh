#!/usr/bin/env bash
# A/B the connection-pool change: grit (MaxIdleConns=10, the scaffold default)
# against grit-pooled (MaxIdleConns=100). Identical binaries otherwise.
#
# Samples CPU DURING the run, not after — the first attempt at this sampled
# once k6 had already exited and dutifully reported 0.17%, which says nothing.
set -uo pipefail
cd "$(dirname "$0")"
export MSYS_NO_PATHCONV=1

REPS="${REPS:-3}"
VUS="${VUS:-50}"
DURATION="${DURATION:-30s}"
SCENARIOS=(list show mixed write)
mkdir -p results/ab

for app in grit grit-pooled; do
  other=$([ "$app" = grit ] && echo grit-pooled || echo grit)
  docker compose stop "$other" laravel >/dev/null 2>&1
  docker compose up -d "$app" >/dev/null 2>&1
  for _ in $(seq 1 60); do
    docker run --rm --network benchmarks_default curlimages/curl -sf \
      "http://$app:8080/api/v1/products?page_size=1" >/dev/null 2>&1 && break
    sleep 1
  done

  for scenario in "${SCENARIOS[@]}"; do
    for rep in $(seq 1 "$REPS"); do
      # warm-up, discarded
      docker run --rm --network benchmarks_default -v "$PWD:/bench" -w /bench \
        -e BASE="http://$app:8080" -e SCENARIO="$scenario" -e VUS="$VUS" -e DURATION=8s \
        grafana/k6 run -q k6/bench.js >/dev/null 2>&1

      ( sleep 14
        docker stats --no-stream --format "{{.Name}} {{.CPUPerc}}" 2>/dev/null \
          | grep -E "benchmarks-(grit|grit-pooled|postgres)" \
          > "results/ab/$app-$scenario-$rep.cpu" ) &

      docker run --rm --network benchmarks_default -v "$PWD:/bench" -w /bench \
        -e BASE="http://$app:8080" -e SCENARIO="$scenario" -e VUS="$VUS" -e DURATION="$DURATION" \
        grafana/k6 run --summary-export "results/ab/$app-$scenario-$rep.json" k6/bench.js \
        > "results/ab/$app-$scenario-$rep.log" 2>&1
      wait

      printf "%-12s %-6s r%s  " "$app" "$scenario" "$rep"
      grep -oE "http_reqs[.]+: [0-9]+ +[0-9.]+/s" "results/ab/$app-$scenario-$rep.log" | head -1
      sed 's/^/               /' "results/ab/$app-$scenario-$rep.cpu" 2>/dev/null
    done
  done
done
echo "done — results/ab/"
