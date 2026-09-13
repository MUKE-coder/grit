"""Body size and deadlines, against a running generated API (H14, v3.252.0).

Before v3.252.0 a global 10 MB body cap wrapped every request before the upload
handler set its own, larger one, so no upload over 10 MB could succeed. The
server's 15 second timeouts were one deadline for the whole request and the
whole response, so a slow upload failed and a streamed export to a slow client
was cut off after it had already answered 200.

    python3 tests/live/verify_limits.py http://localhost:8080 path/to/app.db [--storage]

The database path is used to seed rows and to promote the test account to ADMIN.
--storage adds the upload checks, which need object storage the API can reach.
"""
import json
import socket
import sqlite3
import sys
import time
import urllib.parse
import uuid

base = sys.argv[1].rstrip('/')
db_path = sys.argv[2]
with_storage = '--storage' in sys.argv
PASSWORD = 'SuperSecret123!'
ROWS = 300000
url = urllib.parse.urlparse(base)
results = []


def check(name, ok, detail=''):
    results.append((name, ok, detail))


def request(method, path, payload=b'', content_type='application/json', bearer=None, spread=0.0):
    """One HTTP/1.1 request on a raw socket, returning (status, body).

    Expect: 100-continue means a request refused on its headers is refused before
    the body is sent. spread sends the body in 40 pieces over that many seconds.
    """
    sock = socket.create_connection((url.hostname, url.port or 80), timeout=120)
    head = '%s %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n' % (method, path, url.netloc)
    if bearer:
        head += 'Authorization: Bearer %s\r\n' % bearer
    if payload:
        head += 'Content-Type: %s\r\nContent-Length: %d\r\nExpect: 100-continue\r\n' % (content_type, len(payload))
    sock.sendall((head + '\r\n').encode())
    reply = b''

    def read_head():
        nonlocal reply
        while b'\r\n\r\n' not in reply:
            part = sock.recv(65536)
            if not part:
                break
            reply += part

    if payload:
        read_head()
        if reply.startswith(b'HTTP/1.1 100'):
            reply = reply.split(b'\r\n\r\n', 1)[1]
            step = max(1, len(payload) // 40) if spread else len(payload)
            try:
                for i in range(0, len(payload), step):
                    sock.sendall(payload[i:i + step])
                    if spread:
                        time.sleep(spread / 40)
            except OSError as e:
                sock.close()
                return 0, str(e).encode()
    try:
        while True:
            part = sock.recv(1 << 20)
            if not part:
                break
            reply += part
    except OSError:
        pass
    sock.close()
    status_line, _, rest = reply.partition(b'\r\n')
    try:
        status = int(status_line.split()[1])
    except (IndexError, ValueError):
        return 0, reply[:200]
    return status, rest.partition(b'\r\n\r\n')[2]


def api(method, path, body=None, bearer=None):
    status, raw = request(method, path, json.dumps(body).encode() if body is not None else b'', bearer=bearer)
    try:
        return status, json.loads(raw or b'{}')
    except ValueError:
        return status, {}


def code_of(body):
    if isinstance(body, bytes):
        try:
            body = json.loads(body or b'{}')
        except ValueError:
            return None
    return ((body or {}).get('error') or {}).get('code')


email = 'limits-%s@example.com' % uuid.uuid4().hex[:8]
status, body = api('POST', '/api/v1/auth/register',
                   {'first_name': 'Big', 'last_name': 'Body', 'email': email, 'password': PASSWORD})
if status not in (200, 201):
    sys.exit('register: %s %s' % (status, body))

con = sqlite3.connect(db_path)
con.execute("UPDATE users SET role = 'ADMIN' WHERE email = ?", (email,))
con.commit()
status, body = api('POST', '/api/v1/auth/login', {'email': email, 'password': PASSWORD})
token = ((body.get('data') or {}).get('tokens') or {}).get('access_token')
if not token:
    sys.exit('login: %s %s' % (status, body))

# An ordinary route keeps the 10 MB cap.
status, raw = request('POST', '/api/v1/contacts', b'{"name":"' + b'x' * (11 << 20) + b'"}', bearer=token)
check('an 11 MB JSON body to an ordinary route is refused (413)', status == 413, (status, code_of(raw)))

# Seed enough rows that the CSV export is far bigger than the socket buffers.
cols = con.execute('PRAGMA table_info(contacts)').fetchall()
have = con.execute('SELECT COUNT(*) FROM contacts WHERE deleted_at IS NULL').fetchone()[0]
if have < ROWS:
    now = time.strftime('%Y-%m-%d %H:%M:%S')
    names = [c[1] for c in cols if not (c[5] and c[2].upper().startswith('INTEGER'))]

    def row(i):
        values = []
        for _, name, ctype, notnull, _, pk in cols:
            ctype = ctype.upper()
            if pk and ctype.startswith('INTEGER'):
                continue
            if name == 'id':
                values.append(str(uuid.uuid4()))
            elif name == 'name':
                values.append('Contact %d' % i)
            elif name == 'email':
                values.append('contact%d-%s@example.com' % (i, uuid.uuid4().hex[:6]))
            elif name in ('created_at', 'updated_at'):
                values.append(now)
            elif name == 'deleted_at':
                values.append(None)
            elif notnull:
                values.append(0 if any(t in ctype for t in ('INT', 'REAL', 'NUM')) else '')
            else:
                values.append(None)
        return values

    con.executemany('INSERT INTO contacts (%s) VALUES (%s)' % (','.join(names), ','.join('?' for _ in names)),
                    (row(i) for i in range(ROWS - have)))
    con.commit()
total = con.execute('SELECT COUNT(*) FROM contacts WHERE deleted_at IS NULL').fetchone()[0]
con.close()


def slow_export(slow_seconds):
    """GET the CSV export, reading slowly for slow_seconds, then at full speed."""
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_RCVBUF, 4096)
    sock.connect((url.hostname, url.port or 80))
    sock.sendall(('GET /api/v1/contacts/export?format=csv HTTP/1.1\r\nHost: %s\r\n'
                  'Authorization: Bearer %s\r\nConnection: close\r\n\r\n' % (url.netloc, token)).encode())
    sock.settimeout(120)
    started, chunks = time.time(), []
    while True:
        slow = time.time() - started < slow_seconds
        try:
            part = sock.recv(2048 if slow else 1 << 20)
        except OSError:
            break
        if not part:
            break
        chunks.append(part)
        if slow:
            time.sleep(0.1)
    sock.close()
    raw = b''.join(chunks)
    head, _, rest = raw.partition(b'\r\n\r\n')
    status_line = head.split(b'\r\n', 1)[0].decode(errors='replace')
    if b'transfer-encoding: chunked' in head.lower():
        body, complete = b'', False
        while rest:
            size_line, _, rest = rest.partition(b'\r\n')
            try:
                size = int(size_line.split(b';')[0], 16)
            except ValueError:
                break
            if size == 0:
                complete = True
                break
            body, rest = body + rest[:size], rest[size + 2:]
    else:
        body, complete = rest, True
    return status_line, complete, body.count(b'\n'), time.time() - started


status_line, complete, lines, took = slow_export(20)
check('a streamed export to a slow reader is not cut off',
      ' 200 ' in status_line + ' ' and complete and lines >= total,
      '%s, complete=%s, %d lines for %d rows, %.0fs' % (status_line, complete, lines, total, took))

if with_storage:
    boundary = uuid.uuid4().hex
    # A 20 MB MP4 for a video field, the case the review found could never succeed:
    # an ftyp box so the bytes sniff as video/mp4, then padding.
    mp4 = b'\x00\x00\x00\x18ftypmp42\x00\x00\x00\x00mp42isom'
    payload = (('--%s\r\nContent-Disposition: form-data; name="file"; filename="big.mp4"\r\n'
                'Content-Type: video/mp4\r\n\r\n' % boundary).encode()
               + mp4 + b'\0' * ((20 << 20) - len(mp4)) + ('\r\n--%s--\r\n' % boundary).encode())
    multipart = 'multipart/form-data; boundary=' + boundary
    status, raw = request('POST', '/api/v1/uploads?accepts=video', payload, multipart, bearer=token)
    check('a 20 MB upload succeeds', status in (200, 201), (status, code_of(raw), raw[:160]))
    status, raw = request('POST', '/api/v1/uploads?accepts=video', payload, multipart, bearer=token, spread=25)
    check('a 20 MB upload sent over 25 seconds succeeds', status in (200, 201), (status, code_of(raw), raw[:160]))

width = max(len(n) for n, _, _ in results)
failed = 0
for name, ok, detail in results:
    failed += 0 if ok else 1
    print('%-*s  %s  %s' % (width, name, 'PASS' if ok else 'FAIL', '' if ok else detail))
print('\n%d/%d passed' % (len(results) - failed, len(results)))
sys.exit(1 if failed else 0)
