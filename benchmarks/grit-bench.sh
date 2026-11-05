#!/usr/bin/env bash
#
# Rebuild the Grit side of the benchmark.
#
# grit-bench/ is gitignored, which was the right call and a bad one at the same
# time. Right, because it is scaffold output: committing it would let it drift
# from what the CLI emits today, and then the benchmark would be measuring a
# snapshot of a Grit that no longer exists. Bad, because the one side whose
# numbers the page makes a claim about was the one side nobody outside this
# machine could rebuild (grit#164).
#
# This script is the answer to both: the app is still not committed, and anyone
# can produce it, byte for byte, from a pinned version of the CLI.
#
# The whole point is that there is nothing bespoke in it. No tuned handler, no
# hand-written query, no benchmark-only shortcut. If the generated code is slow,
# the benchmark says so, which is the only way the number means anything.
#
# Usage:
#   ./grit-bench.sh                 # build with the pinned version
#   ./grit-bench.sh --version 3.465.0
#   ./grit-bench.sh --keep          # do not delete an existing grit-bench/
#
set -euo pipefail

# The version the numbers in results/ were produced with. Bump it, re-run the
# pair, and say in the results which version they came from. A benchmark whose
# subject is "latest" is a benchmark nobody can reproduce.
PINNED="3.465.0"

VERSION="$PINNED"
KEEP=0
while [ $# -gt 0 ]; do
  case "$1" in
    --version) VERSION="$2"; shift 2 ;;
    --keep)    KEEP=1; shift ;;
    -h|--help) sed -n '2,24p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) echo "unknown argument: $1" >&2; exit 2 ;;
  esac
done

cd "$(dirname "$0")"

if ! command -v grit >/dev/null 2>&1; then
  echo "grit is not on PATH. Install it: https://gritframework.com/docs/getting-started/installation" >&2
  exit 1
fi

HAVE="$(grit version | awk '{print $NF}')"
if [ "$HAVE" != "$VERSION" ]; then
  echo "WARNING: grit on PATH is $HAVE, the benchmark is pinned to $VERSION." >&2
  echo "         The numbers in results/ were produced with $VERSION. Either" >&2
  echo "         install that version, or re-run the whole pair and say which" >&2
  echo "         version produced the new numbers." >&2
  echo >&2
fi

if [ -d grit-bench ]; then
  if [ "$KEEP" = "1" ]; then
    echo "grit-bench/ exists and --keep was given; nothing to do."
    exit 0
  fi
  echo "Removing the existing grit-bench/ ..."
  rm -rf grit-bench
fi

# --api, because the benchmark drives HTTP and never opens a browser. A web and
# an admin frontend would add install time and nothing else: they are not in the
# request path, and Laravel's side has no frontend either.
echo "==> grit new grit-bench --api"
grit new grit-bench --api

# Column for column what laravel-bench's migration creates. price is a float
# whose name makes the generator store it as decimal(12,2), which is the same
# column Laravel's $table->decimal('price', 12, 2) produces: a benchmark where
# one side stores money as a double and the other as a numeric is measuring the
# column type.
echo "==> grit generate resource Product"
(
  cd grit-bench
  grit generate resource Product \
    --fields "name:string,sku:string,description:text,price:float,stock:int,active:bool"
)

cat <<'EOF'

Done. grit-bench/ is a plain scaffold with one generated resource, which is the
claim the benchmark rests on: nothing in it was written for the benchmark.

Next:

  docker compose up -d postgres redis
  docker compose exec -T postgres psql -U bench -d bench \
    -c "CREATE DATABASE bench_grit;" -c "CREATE DATABASE bench_laravel;"
  (cd grit-bench && grit migrate)

Then seed both databases from seed/products.sql and run ./pair.sh, as the
README describes.
EOF
