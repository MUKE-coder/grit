# Production-shaped Laravel: nginx + php-fpm + opcache, with config, route and
# view caches warmed at build time. This is what a normal Laravel deployment
# looks like — not `php artisan serve`, which is a single-threaded dev server
# and would make the comparison meaningless.
#
# Deliberately NOT Octane/FrankenPHP. Octane keeps the framework booted between
# requests and is considerably faster, but it is opt-in and not what most
# Laravel apps run. If we benchmark Octane we should say so and label it
# separately rather than quietly using it as "Laravel".
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

RUN mkdir -p storage/framework/{cache,sessions,views} storage/logs bootstrap/cache \
    && chown -R www-data:www-data storage bootstrap/cache \
    && chmod -R 775 storage bootstrap/cache

EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
