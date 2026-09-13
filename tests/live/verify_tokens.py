"""Token types and revocation, against a running generated API (H5, v3.244.0).

Before v3.244.0 access and refresh tokens were the same shape, so a 7-day refresh
token was accepted as a bearer token and on the WebSocket. The auth middleware
never looked at the sessions table, so logout, "sign out everywhere" and a
password change left every issued access token working. And the refresh cookie
was scoped to /api/auth while the routes live under /api/v1/auth, so a browser
never sent it: cookie refresh always failed and logout never revoked anything.

    python3 tests/live/verify_tokens.py http://localhost:8080
"""
import http.cookiejar
import json
import sys
import urllib.error
import urllib.request
import uuid

base = (sys.argv[1] if len(sys.argv) > 1 else 'http://localhost:8080').rstrip('/')
PASSWORD = 'SuperSecret123!'
results = []


def check(name, ok, detail=''):
    results.append((name, ok, detail))


class Client:
    """One device: a cookie jar, as a browser has."""

    def __init__(self):
        self.jar = http.cookiejar.CookieJar()
        self.opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(self.jar))

    def cookie(self, name):
        for c in self.jar:
            if c.name == name:
                return c
        return None

    def call(self, method, path, body=None, bearer=None, csrf=False):
        data = json.dumps(body).encode() if body is not None else None
        req = urllib.request.Request(base + path, data=data, method=method)
        if data is not None:
            req.add_header('Content-Type', 'application/json')
        if bearer:
            req.add_header('Authorization', 'Bearer ' + bearer)
        if csrf and self.cookie('grit_csrf'):
            req.add_header('X-CSRF-Token', self.cookie('grit_csrf').value)
        try:
            with self.opener.open(req, timeout=30) as r:
                return r.status, json.loads(r.read() or b'{}')
        except urllib.error.HTTPError as e:
            raw = e.read()
            try:
                return e.code, json.loads(raw or b'{}')
            except ValueError:
                return e.code, {'raw': raw[:120].decode(errors='replace')}


def bare(method, path, bearer=None, body=None):
    return Client().call(method, path, body=body, bearer=bearer)


email = 'tokens-%s@example.com' % uuid.uuid4().hex[:8]
status, body = bare('POST', '/api/v1/auth/register',
                    body={'first_name': 'Token', 'last_name': 'Check', 'email': email, 'password': PASSWORD})
if status not in (200, 201):
    sys.exit('register: %s %s' % (status, body))


def login(client=None):
    client = client or Client()
    status, body = client.call('POST', '/api/v1/auth/login', {'email': email, 'password': PASSWORD})
    tokens = ((body.get('data') or {}).get('tokens') or {})
    if status != 200 or not tokens.get('access_token'):
        sys.exit('login: %s %s' % (status, body))
    return client, tokens['access_token'], tokens['refresh_token']


# 1. A refresh token is not an access token, on the API or the WebSocket.
_, access, refresh = login()
status, _ = bare('GET', '/api/v1/auth/me', bearer=access)
check('an access token reaches /auth/me', status == 200, status)
status, _ = bare('GET', '/api/v1/auth/me', bearer=refresh)
check('a refresh token is refused as a bearer token (401)', status == 401, status)
status, _ = bare('GET', '/api/ws?token=' + refresh)
check('a refresh token is refused on the WebSocket (401)', status == 401, status)
status, _ = bare('POST', '/api/v1/auth/refresh', body={'refresh_token': access})
check('an access token is refused by /auth/refresh (401)', status == 401, status)

# 2. The browser flow: the refresh cookie reaches the refresh and logout routes.
browser, _, _ = login()
cookie = browser.cookie('grit_refresh')
check('the refresh cookie is scoped to the versioned auth routes',
      cookie is not None and cookie.path == '/api/v1/auth', cookie and cookie.path)
before = browser.cookie('grit_access').value
status, body = browser.call('POST', '/api/v1/auth/refresh')
check('a browser refreshes with the cookie alone', status == 200, (status, body))
after = browser.cookie('grit_access').value
check('and receives a new access token', after != before)
status, _ = bare('GET', '/api/v1/auth/me', bearer=after)
check('the refreshed access token works', status == 200, status)

# 3. Logout ends the session, and the access token with it.
browser.call('GET', '/api/v1/auth/me')  # picks up the CSRF cookie
status, body = browser.call('POST', '/api/v1/auth/logout', csrf=True)
check('logout answers', status == 200, (status, body))
status, _ = bare('GET', '/api/v1/auth/me', bearer=after)
check('an access token stops working when its session logs out (401)', status == 401, status)

# 4. Sign out everywhere reaches the other devices' access tokens.
_, laptop_access, _ = login()
_, phone_access, phone_refresh = login()
status, body = bare('POST', '/api/v1/auth/sessions/revoke-all', bearer=laptop_access)
check('sign out everywhere answers', status == 200, (status, body))
status, _ = bare('GET', '/api/v1/auth/me', bearer=phone_access)
check('another device\'s access token stops working (401)', status == 401, status)
status, _ = bare('POST', '/api/v1/auth/refresh', body={'refresh_token': phone_refresh})
check('and its refresh token too (401)', status == 401, status)

width = max(len(n) for n, _, _ in results)
failed = 0
for name, ok, detail in results:
    failed += 0 if ok else 1
    print('%-*s  %s  %s' % (width, name, 'PASS' if ok else 'FAIL', '' if ok else detail))
print('\n%d/%d passed' % (len(results) - failed, len(results)))
sys.exit(1 if failed else 0)
