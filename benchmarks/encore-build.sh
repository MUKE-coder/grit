#!/usr/bin/env bash
# Encore needs its own CLI to compile — there is no plain `tsc && node` path,
# because the framework generates a Rust-backed runtime around your handlers.
#
# The CLI is Linux-only, so it runs in a container here. The Docker socket is
# mounted so `encore build docker` can produce the image on the host daemon,
# and the app is bind-mounted so the build sees the real source.
set -euo pipefail
cd "$(dirname "$0")"
export MSYS_NO_PATHCONV=1

IMAGE="${IMAGE:-benchmarks-encore:latest}"

docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$PWD/encore-bench:/app" \
  -w /app \
  node:22-bookworm \
  bash -c '
    set -e
    apt-get update -qq && apt-get install -y -qq docker.io >/dev/null 2>&1
    curl -sL https://encore.dev/install.sh | bash >/dev/null 2>&1
    export PATH="$HOME/.encore/bin:$PATH"
    encore version
    npm install --no-audit --no-fund
    encore build docker --skip-config '"$IMAGE"'
  '

echo "built $IMAGE"
