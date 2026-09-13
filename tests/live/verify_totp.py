"""Two-factor sign-in: guessing and replay, against a running generated API (H6, v3.245.0).

Before v3.245.0 a pending 2FA token accepted unlimited guesses for five minutes,
and a correct password cleared the failure count, so signing in again handed
out fresh guesses. A code could be used again for as long as it was valid, and
enabling 2FA silently replaced a secret that was already enabled.

    python3 tests/live/verify_totp.py http://localhost:8080
"""
import base64
import hashlib
import hmac
import json
import struct
import sys
import time
import urllib.error
import urllib.request
import uuid

base = (sys.argv[1] if len(sys.argv) > 1 else 'http://localhost:8080').rstrip('/')
PASSWORD = 'SuperSecret123!'
results = []


def check(name, ok, detail=''):
    results.append((name, ok, detail))


def call(method, path, body=None, bearer=None):
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(base + path, data=data, method=method)
    if data is not None:
        req.add_header('Content-Type', 'application/json')
    if bearer:
        req.add_header('Authorization', 'Bearer ' + bearer)
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            return r.status, json.loads(r.read() or b'{}')
    except urllib.error.HTTPError as e:
        raw = e.read()
        try:
            return e.code, json.loads(raw or b'{}')
        except ValueError:
            return e.code, {}


def code_of(body):
    return ((body or {}).get('error') or {}).get('code')


def totp(secret, step):
    key = base64.b32decode(secret.upper() + '=' * (-len(secret) % 8))
    digest = hmac.new(key, struct.pack('>Q', step), hashlib.sha1).digest()
    offset = digest[-1] & 0x0F
    value = struct.unpack('>I', digest[offset:offset + 4])[0] & 0x7FFFFFFF
    return '%06d' % (value % 1000000)


def step_now():
    return int(time.time()) // 30


# Keep clear of a step boundary, so "the current step" means the same thing to
# the script and the server for the length of the test.
while int(time.time()) % 30 > 20:
    time.sleep(1)

email = 'totp-%s@example.com' % uuid.uuid4().hex[:8]
status, body = call('POST', '/api/v1/auth/register',
                    {'first_name': 'Two', 'last_name': 'Factor', 'email': email, 'password': PASSWORD})
if status not in (200, 201):
    sys.exit('register: %s %s' % (status, body))


def login():
    return call('POST', '/api/v1/auth/login', {'email': email, 'password': PASSWORD})


status, body = login()
access = ((body.get('data') or {}).get('tokens') or {}).get('access_token')
status, body = call('POST', '/api/v1/auth/totp/setup', bearer=access)
secret = (body.get('data') or {}).get('secret')
if not secret:
    sys.exit('setup: %s %s' % (status, body))

s0 = step_now()
status, body = call('POST', '/api/v1/auth/totp/enable', {'secret': secret, 'code': totp(secret, s0)}, bearer=access)
check('2FA enables with a valid code', status == 200, (status, body))
status, body = call('POST', '/api/v1/auth/totp/enable', {'secret': secret, 'code': totp(secret, s0 + 1)}, bearer=access)
check('enabling again does not replace the secret (409)', status == 409, (status, code_of(body)))

# Replay: the code used to enable is spent.
status, body = login()
pending = (body.get('data') or {}).get('pending_token')
check('sign-in asks for the second factor', status == 200 and pending, (status, body))
status, body = call('POST', '/api/v1/auth/totp/verify', {'pending_token': pending, 'code': totp(secret, s0)})
check('the code used to enable cannot be used to sign in (401)', status == 401, (status, code_of(body)))
status, body = call('POST', '/api/v1/auth/totp/verify', {'pending_token': pending, 'code': totp(secret, s0 + 1)})
check('the next code signs in', status == 200, (status, code_of(body)))

status, body = login()
pending = (body.get('data') or {}).get('pending_token')
status, body = call('POST', '/api/v1/auth/totp/verify', {'pending_token': pending, 'code': totp(secret, s0 + 1)})
check('a code that signed in cannot sign in again (401)', status == 401, (status, code_of(body)))

# Guessing: four more wrong codes spend this pending token.
for _ in range(4):
    call('POST', '/api/v1/auth/totp/verify', {'pending_token': pending, 'code': '000000'})
status, body = call('POST', '/api/v1/auth/totp/verify', {'pending_token': pending, 'code': '111111'})
check('after five wrong codes the pending token is spent', status == 401 and code_of(body) == 'INVALID_PENDING_TOKEN',
      (status, code_of(body)))

# Signing in again does not reset the count: ten wrong codes lock the account.
# The successful sign-in above cleared it, and this account has five wrong codes
# since then, so five more on a new sign-in make ten.
status, body = login()
pending = (body.get('data') or {}).get('pending_token')
check('the password still works after the pending token is spent', status == 200 and pending, (status, body))
for _ in range(5):
    call('POST', '/api/v1/auth/totp/verify', {'pending_token': pending, 'code': '222222'})
status, body = login()
check('ten wrong codes across sign-ins lock the account (429)', status == 429 and code_of(body) == 'ACCOUNT_LOCKED',
      (status, code_of(body)))

width = max(len(n) for n, _, _ in results)
failed = 0
for name, ok, detail in results:
    failed += 0 if ok else 1
    print('%-*s  %s  %s' % (width, name, 'PASS' if ok else 'FAIL', '' if ok else detail))
print('\n%d/%d passed' % (len(results) - failed, len(results)))
sys.exit(1 if failed else 0)
