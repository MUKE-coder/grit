"""The production compose stack hands each side container only its own settings (H10, v3.248.0).

Before v3.248.0 Postgres and MinIO were given the whole .env through env_file:
the JWT secret, the Sentinel keys and every dashboard password. MinIO's root
credentials fell back to minioadmin, and it sits behind the public proxy for
presigned uploads.

    docker compose -f docker-compose.prod.yml config --format json > compose.json
    python3 tests/live/verify_compose.py compose.json
"""
import json
import sys

API_SECRETS = {'JWT_SECRET', 'SENTINEL_SECRET_KEY', 'SENTINEL_PASSWORD', 'SENTINEL_AUDIT_KEY',
               'PULSE_PASSWORD', 'GORM_STUDIO_PASSWORD', 'FIELD_ENCRYPTION_KEY'}
results = []


def check(name, ok, detail=''):
    results.append((name, ok, detail))


services = json.load(open(sys.argv[1], encoding='utf-8'))['services']

for name in ('postgres', 'pgbouncer', 'redis', 'minio'):
    if name not in services:
        continue
    env = services[name].get('environment') or {}
    leaked = sorted(API_SECRETS & set(env))
    check('%s receives none of the API\'s secrets' % name, not leaked, leaked)

minio = (services.get('minio') or {}).get('environment') or {}
if services.get('minio'):
    check('MinIO\'s root password is set and not minioadmin',
          minio.get('MINIO_ROOT_PASSWORD') not in (None, '', 'minioadmin'), bool(minio.get('MINIO_ROOT_PASSWORD')))

api = (services.get('api') or {}).get('environment') or {}
check('the api still receives its own secrets', 'JWT_SECRET' in api, sorted(api)[:5])

width = max(len(n) for n, _, _ in results)
failed = 0
for name, ok, detail in results:
    failed += 0 if ok else 1
    print('%-*s  %s  %s' % (width, name, 'PASS' if ok else 'FAIL', '' if ok else detail))
print('\n%d/%d passed' % (len(results) - failed, len(results)))
sys.exit(1 if failed else 0)
