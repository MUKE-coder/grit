#!/usr/bin/env bash
# Run the docs blocks that claim to work.
#
# A block marked verify="<flow>" in a docs page is a claim that the commands in
# it work. This turns the claim into a test: for each flow, scaffold a project,
# migrate it, and run the flow's blocks in page order, failing on the first
# non-zero exit.
#
# "The docs said to do X and X did not work" is the most repeated root cause in
# this project's changelog, and every instance was found by a person following a
# page. This is the cheaper way.
#
# Usage:
#   tests/docs/run-flows.sh                    # every flow
#   tests/docs/run-flows.sh migrate-rollback   # one flow
#
# Needs: a grit on PATH built from this checkout, a reachable Postgres, and
# DATABASE_URL_BASE pointing at it without a database name, e.g.
#   DATABASE_URL_BASE=postgres://grit:grit@localhost:5432
# Set WORK to keep the scratch projects somewhere you can look at them.
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DOCS="$REPO/docs/app"
WORK="${WORK:-$(mktemp -d)}"
: "${DATABASE_URL_BASE:?set DATABASE_URL_BASE to a Postgres server, with no database name}"
PSQL="${PSQL:-psql}"
ADMIN_URL="${ADMIN_URL:-$DATABASE_URL_BASE/postgres}"

flows=("$@")
if [ ${#flows[@]} -eq 0 ]; then
  # shellcheck disable=SC2207
  flows=($(cd "$REPO" && go run ./tools/docsamples -docs "$DOCS" -list | cut -f1))
fi
if [ ${#flows[@]} -eq 0 ]; then
  echo "no verify= blocks found in the docs: has the marker changed?" >&2
  exit 1
fi

echo "running ${#flows[@]} docs flow(s) in $WORK"
index=0
for flow in "${flows[@]}"; do
  index=$((index + 1))
  echo
  echo "=============================================================="
  echo "  flow: $flow"
  echo "=============================================================="

  database="docsflow_${index}"
  # PSQL and ADMIN_URL are overridable so this runs against a Postgres in Docker
  # with no client installed on the host: PSQL="docker exec -i pg psql" with an
  # ADMIN_URL the container can reach. Unquoted on purpose, so a wrapper command
  # splits into words.
  # shellcheck disable=SC2086
  $PSQL "$ADMIN_URL" -v ON_ERROR_STOP=1 -q \
    -c "DROP DATABASE IF EXISTS $database" -c "CREATE DATABASE $database"
  export DATABASE_URL="$DATABASE_URL_BASE/$database?sslmode=disable"

  directory="$WORK/$flow"
  rm -rf "$directory"
  mkdir -p "$directory"
  cd "$directory"

  # A flow named fresh-* starts from nothing, because its first command is the
  # one that creates the project. Everything else starts from a project that has
  # been scaffolded and migrated, which is where a reader of that page is.
  if [[ "$flow" == fresh-* ]]; then
    echo "  (starting from an empty directory)"
  else
    grit new app --api > "$directory/scaffold.log" 2>&1 ||
      { echo "scaffolding failed:"; tail -30 "$directory/scaffold.log"; exit 1; }
    cd "$directory/app"
    grit migrate > "$directory/migrate.log" 2>&1 ||
      { echo "the first migrate failed:"; tail -30 "$directory/migrate.log"; exit 1; }
  fi

  (cd "$REPO" && go run ./tools/docsamples -docs "$DOCS" -flow "$flow") > "$directory/flow.sh"
  bash "$directory/flow.sh"
  echo "  flow $flow: ok"
done

echo
echo "all docs flows passed"
