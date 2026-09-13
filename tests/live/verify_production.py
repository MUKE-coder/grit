"""Production mode, against a generated API started with APP_ENV=production.

Production is the default and the strict mode (v3.243.0). Before it, the login
and register rate limits and AuthShield were keyed to /api/auth/... while the
router mounts /api/v1, so they never fired; GORM Studio, a writable SQL console,
was mounted; and /docs published every route with a console to call them.

    python3 tests/live/verify_production.py http://localhost:8081
"""
import json
import sys
import urllib.error
import urllib.request
import uuid

base = (sys.argv[1] if len(sys.argv) > 1 else 'http://localhost:8081').rstrip('/')
results = []


def call(method, path, body=None, headers=None):
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(base + path, data=data, method=method)
    if data:
        req.add_header('Content-Type', 'application/json')
    for k, v in (headers or {}).items():
        req.add_header(k, v)
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            return r.status, r.read()
    except urllib.error.HTTPError as e:
        return e.code, e.read()


def check(name, ok, detail=''):
    results.append((name, ok, detail))


status, _ = call('GET', '/studio/')
check('Studio: the SQL console is not mounted in production', status == 404, status)
status, _ = call('GET', '/docs')
check('Docs: the API reference is not published in production', status == 404, status)

email = 'nobody-%s@example.com' % uuid.uuid4().hex[:8]
# A stale access cookie with no CSRF token: signing in is how you replace it, so
# the login route must not demand a token. The exemption never matched /api/v1.
status, body = call('POST', '/api/v1/auth/login', {'email': email, 'password': 'wrong-password-0'},
                    headers={'Cookie': 'grit_access=stale'})
check('CSRF: a sign-in with a stale session cookie is judged on its credentials', status == 401, (status, body[:160]))

codes = [status]
for i in range(6):
    status, _ = call('POST', '/api/v1/auth/login', {'email': email, 'password': 'wrong-password-%d' % (i + 1)})
    codes.append(status)
check('Rate limit: five failed logins are refused as credentials', codes[:5] == [401] * 5, codes)
check('Rate limit: the sixth inside fifteen minutes is 429', 429 in codes[5:], codes)

width = max(len(n) for n, _, _ in results)
failed = 0
for name, ok, detail in results:
    failed += 0 if ok else 1
    print('%-*s  %s  %s' % (width, name, 'PASS' if ok else 'FAIL', '' if ok else detail))
print('\n%d/%d passed' % (len(results) - failed, len(results)))
sys.exit(1 if failed else 0)
