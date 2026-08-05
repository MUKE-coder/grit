# Production-shaped Laravel: nginx + php-fpm + opcache, production-only vendor,
# and config, route, view and event caches warmed at container start. This is what a normal Laravel deployment
# looks like — not `php artisan serve`, which is a single-threaded dev server
# and would make the comparison meaningless.
#
# Deliberately NOT Octane/FrankenPHP. Octane keeps the framework booted between
# requests and is considerably faster, but it is opt-in and not what most
# Laravel apps run. If we benchmark Octane we should say so and label it
# separately rather than quietly using it as "Laravel".
# Vendor built in its own stage so the production flags are not optional:
# --no-dev keeps laravel/pail, collision, phpunit, mockery, faker and pint out of
# the autoloader and out of package discovery. Shipping those was costing real
# bootstrap time on every request and is not what anyone deploys.
FROM composer:2 AS vendor
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install       --no-dev       --no-interaction       --no-progress       --no-scripts       --ignore-platform-reqs

FROM php:8.4-fpm-alpine

RUN apk add --no-cache nginx supervisor postgresql-dev \
    && docker-php-ext-install pdo_pgsql opcache bcmath \
    && rm -rf /var/cache/apk/*

# Opcache on, with a JIT. Without this PHP recompiles every file on every
# request and the number you get is not one anyone would ship.
RUN { \
      echo 'opcache.enable=1'; \
      echo 'opcache.enable_cli=1'; \
      echo 'opcache.memory_consumption=256'; \
      echo 'opcache.interned_strings_buffer=32'; \
      echo 'opcache.max_accelerated_files=20000'; \
      echo 'opcache.validate_timestamps=0'; \
      echo 'opcache.jit_buffer_size=128M'; \
      echo 'opcache.jit=tracing'; \
      echo 'realpath_cache_size=4096K'; \
      echo 'realpath_cache_ttl=600'; \
      echo 'memory_limit=512M'; \
    } > /usr/local/etc/php/conf.d/zz-perf.ini

# Static worker pool. Sized against the container's CPU limit — too few and
# requests queue behind workers rather than behind the framework, which would
# read as a Laravel problem when it is a configuration problem.
RUN { \
      echo '[www]'; \
      echo 'user = www-data'; \
      echo 'group = www-data'; \
      echo 'listen = 127.0.0.1:9000'; \
      echo 'pm = static'; \
      echo 'pm.max_children = 32'; \
      echo 'pm.max_requests = 0'; \
      echo 'clear_env = no'; \
      echo 'catch_workers_output = yes'; \
    } > /usr/local/etc/php-fpm.d/zz-www.conf

COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/supervisord.conf /etc/supervisord.conf

WORKDIR /var/www/html
COPY . .

# The host's vendor/ is whatever was last installed there — usually with dev
# dependencies. Overwrite it with the production one built above.
COPY --from=vendor /app/vendor ./vendor
COPY --from=vendor /usr/bin/composer /usr/bin/composer

# Build the optimised classmap HERE, not in the vendor stage. That stage only
# has composer.json and composer.lock, so a classmap built there contains no
# App\* classes at all — and with --classmap-authoritative the autoloader stops
# looking anywhere else, so every controller resolves to "does not exist".
RUN composer dump-autoload       --no-dev       --optimize       --classmap-authoritative       --no-interaction       --no-scripts

# Belt and braces: .dockerignore keeps these out of the build context, but a
# stale manifest here names service providers that --no-dev removed, and Laravel
# cannot boot far enough to rebuild it.
RUN rm -f bootstrap/cache/*.php

RUN mkdir -p storage/framework/{cache,sessions,views} storage/logs bootstrap/cache \
    && chown -R www-data:www-data storage bootstrap/cache \
    && chmod -R 775 storage bootstrap/cache

EXPOSE 8080

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
