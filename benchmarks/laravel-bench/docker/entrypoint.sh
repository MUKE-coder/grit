#!/bin/sh
set -e

# Warm Laravel's production caches HERE rather than at image build time.
#
# `artisan config:cache` bakes the resolved config — including DB_HOST,
# DB_DATABASE and friends — into bootstrap/cache/config.php. At build time those
# environment variables do not exist yet, so a build-time cache would freeze the
# wrong database in and Laravel would ignore what compose passes at runtime.
#
# Doing it at start means the cache is built once per container, from the real
# environment, before the first request arrives. That is also what a real
# deployment does (release step, not build step).
#
# This is not optional dressing. Without it Laravel re-parses every config file
# and re-registers every route on each request, which costs far more than the
# thing being benchmarked.
php artisan optimize:clear >/dev/null 2>&1 || true

# Re-run discovery against the production vendor. It cannot happen at build time
# because storage/ and the environment are not in place yet.
php artisan package:discover --no-ansi
php artisan config:cache
php artisan route:cache
php artisan view:cache
php artisan event:cache

exec "$@"
