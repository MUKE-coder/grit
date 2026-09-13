"""Object storage: only uploads and thumbnails are public (H1, v3.246.0).

Before v3.246.0 the bucket policy granted anonymous s3:GetObject on every key,
so a database backup's key, seen once in a log or a Referer header, was a
permanent anonymous download of the whole database, and the 15-minute signed
link protected nothing: without its query string the URL still worked.

Needs the API, the MinIO it stores to, and the MinIO container's name, because
the probe objects are written with the mc client inside it:

    python3 tests/live/verify_storage.py http://localhost:8080 http://localhost:9002 app-uploads minio
"""
import json
import struct
import subprocess
import sys
import urllib.error
import urllib.request
import uuid
import zlib

api, minio, bucket, container = sys.argv[1].rstrip('/'), sys.argv[2].rstrip('/'), sys.argv[3], sys.argv[4]
results = []


def check(name, ok, detail=''):
    results.append((name, ok, detail))


def status_of(url, data=None, headers=None, method=None):
    req = urllib.request.Request(url, data=data, method=method)
    for k, v in (headers or {}).items():
        req.add_header(k, v)
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            return r.status, r.read()
    except urllib.error.HTTPError as e:
        return e.code, e.read()


def png():
    def chunk(kind, payload):
        return (struct.pack('>I', len(payload)) + kind + payload +
                struct.pack('>I', zlib.crc32(kind + payload) & 0xFFFFFFFF))
    header = struct.pack('>IIBBBBB', 1, 1, 8, 2, 0, 0, 0)
    return (b'\x89PNG\r\n\x1a\n' + chunk(b'IHDR', header) +
            chunk(b'IDAT', zlib.compress(b'\x00\xff\x00\x00')) + chunk(b'IEND', b''))


email = 'storage-%s@example.com' % uuid.uuid4().hex[:8]
body = json.dumps({'first_name': 'Store', 'last_name': 'Check', 'email': email, 'password': 'SuperSecret123!'}).encode()
status_of(api + '/api/v1/auth/register', body, {'Content-Type': 'application/json'}, 'POST')
status, raw = status_of(api + '/api/v1/auth/login',
                        json.dumps({'email': email, 'password': 'SuperSecret123!'}).encode(),
                        {'Content-Type': 'application/json'}, 'POST')
token = json.loads(raw)['data']['tokens']['access_token']

boundary = 'x' + uuid.uuid4().hex
payload = (('--%s\r\nContent-Disposition: form-data; name="file"; filename="pixel.png"\r\n'
            'Content-Type: image/png\r\n\r\n' % boundary).encode() + png() + ('\r\n--%s--\r\n' % boundary).encode())
status, raw = status_of(api + '/api/v1/uploads', payload,
                        {'Content-Type': 'multipart/form-data; boundary=' + boundary, 'Authorization': 'Bearer ' + token},
                        'POST')
url = ((json.loads(raw or b'{}').get('data') or {}).get('url')) if status in (200, 201) else None
check('an image uploads through the API', bool(url), (status, raw[:200]))
if url:
    status, _ = status_of(url)
    check('its URL loads without a signature, as an <img> needs', status == 200, (status, url))

probe = uuid.uuid4().hex[:10]
keys = {
    'backups/%s.zip' % probe: 403,
    'originals/2026/09/%s.jpg' % probe: 403,
    'exports/%s.csv' % probe: 403,
    'uploads/2026/09/%s.txt' % probe: 200,
    'thumbnails/2026/09/%s.txt' % probe: 200,
}
script = 'mc alias set local http://127.0.0.1:9000 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD" >/dev/null'
for key in keys:
    script += ' && echo probe | mc pipe local/%s/%s >/dev/null' % (bucket, key)
done = subprocess.run(['docker', 'exec', container, 'sh', '-c', script], capture_output=True, text=True)
check('probe objects written into the bucket', done.returncode == 0, done.stderr[-300:])

for key, want in keys.items():
    status, _ = status_of('%s/%s/%s' % (minio, bucket, key))
    prefix = key.split('/', 1)[0] + '/'
    if want == 403:
        check('%s is not readable without a signature (403)' % prefix, status == 403, status)
    else:
        check('%s is readable without a signature' % prefix, status == 200, status)

width = max(len(n) for n, _, _ in results)
failed = 0
for name, ok, detail in results:
    failed += 0 if ok else 1
    print('%-*s  %s  %s' % (width, name, 'PASS' if ok else 'FAIL', '' if ok else detail))
print('\n%d/%d passed' % (len(results) - failed, len(results)))
sys.exit(1 if failed else 0)
