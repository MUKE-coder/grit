#!/usr/bin/env python3
"""Live checks against a running generated project.

Every check here is a bug that was found by hand once, in a real application,
and then fixed. Running them by hand found each one; keeping them here is what
stops the next release putting one back.

The suite talks HTTP to a running server and reads Postgres through psql, so it
tests the thing a user deploys rather than the code that writes it. It needs a
project scaffolded and migrated, with these resources generated (see
tests/live/README.md, and .github/workflows/live.yml, which does it in CI):

    grit generate resource Lot     --fields "title:string,current_bid:int"
    grit generate resource Note    --fields "body:text" --owned-by user
    grit generate resource Category --fields "name:string,slug:slug" --tree --public
    grit generate resource Product --fields "name:string,slug:slug,price:money,secret:text:encrypted,category:belongs_to:Category,tags:many_to_many:Tag" --public
    grit generate resource Tag     --fields "name:string"
    grit generate resource Invoice --fields "number:string" --items "InvoiceItem:description:string,qty:int,unit_rate:float"

Usage:
    python verify.py [BASE_URL] [--db thin] [--psql-container NAME]

Exits non-zero if any check fails, so CI fails with it.
"""
import argparse
import json
import subprocess
import sys
import threading
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid

PASSWORD = 'Str0ng!Passw0rd#'
results = []
args = None
RUN = uuid.uuid4().hex


def check(name, ok, detail=''):
    results.append((name, bool(ok), detail))


def call(method, path, token=None, body=None, headers=None, raw=None,
         ctype='application/json'):
    # The public group caches by URL for a minute, so each run reads its own
    # keys rather than the previous run's answers.
    if path.startswith('/api/v1/public/'):
        path += ('&' if '?' in path else '?') + 'run=' + RUN
    data = raw if raw is not None else (json.dumps(body).encode() if body is not None else None)
    req = urllib.request.Request(args.base + path, data=data, method=method)
    if data is not None:
        req.add_header('Content-Type', ctype)
    if token:
        req.add_header('Authorization', 'Bearer ' + token)
    for key, value in (headers or {}).items():
        req.add_header(key, value)
    try:
        with urllib.request.urlopen(req, timeout=30) as response:
            return response.status, {k.lower(): v for k, v in response.headers.items()}, response.read()
    except urllib.error.HTTPError as e:
        return e.code, {k.lower(): v for k, v in e.headers.items()}, e.read()


def data(body):
    try:
        return json.loads(body).get('data')
    except Exception:
        return None


def err(body):
    try:
        return json.loads(body).get('error') or {}
    except Exception:
        return {}


def ids(rows):
    return [r.get('id') for r in rows or []]


def psql(sql):
    """One value out of the database, whichever engine is behind the API.

    Named psql for the Postgres history; it dispatches on --engine now, because a
    suite that can only read Postgres can only verify Postgres, and Grit claims
    three engines.
    """
    if args.engine == "mysql":
        return mysql_query(sql)
    if args.engine == "sqlite":
        return sqlite_query(sql)
    return postgres_query(sql)


def mysql_query(sql):
    command = ["mysql", "-N", "-B", "-h", args.db_host, "-P", str(args.db_port),
               "-u", args.db_user, args.db, "-e", sql]
    if args.db_password:
        command.insert(1, "-p" + args.db_password)
    if args.psql_container:
        command = ["docker", "exec", "-i", args.psql_container, "mysql", "-N", "-B",
                   "-u", args.db_user] + (["-p" + args.db_password] if args.db_password else []) + \
            [args.db, "-e", sql]
    out = subprocess.run(command, capture_output=True, text=True)
    return (out.stdout or "").strip()


def sqlite_query(sql):
    """Read the file with Python's own sqlite3, not a client binary.

    --db carries the path for this engine. Read-only and with a timeout, so a
    check never blocks on the API's writer, and no sqlite3 CLI has to exist on
    the machine running the suite.
    """
    import sqlite3
    # Read-write, not read-only: the suite promotes one account to ADMIN through
    # this, and a read-only handle made that update silently do nothing, which
    # turned every staff check into a 403.
    uri = 'file:%s' % args.db.replace('?', '%3f')
    try:
        connection = sqlite3.connect(uri, uri=True, timeout=10)
    except sqlite3.Error as error:
        return 'sqlite: %s' % error
    try:
        cursor = connection.execute(sql)
        row = cursor.fetchone()
        connection.commit()
        return '' if row is None or row[0] is None else str(row[0])
    except sqlite3.Error as error:
        return 'sqlite: %s' % error
    finally:
        connection.close()


def postgres_query(sql):
    """One value out of Postgres. The database is the only witness for some of
    these: whether a column is ciphertext, whether a path was rewritten."""
    if args.psql_container:
        cmd = ['docker', 'exec', args.psql_container, 'psql', '-U', args.db_user, '-d', args.db, '-tAc', sql]
    else:
        cmd = ['psql', '-h', args.db_host, '-p', args.db_port, '-U', args.db_user, '-d', args.db, '-tAc', sql]
    out = subprocess.run(cmd, capture_output=True, text=True, env=args.psql_env)
    return out.stdout.strip()


def login(email):
    status, _, body = call('POST', '/api/auth/login', body={'email': email, 'password': PASSWORD})
    d = data(body) or {}
    token = (d.get('tokens') or {}).get('access_token') or d.get('access_token')
    if not token:
        sys.exit('login %s: %s %s' % (email, status, body[:300]))
    return token, (d.get('user') or {}).get('id')


def register(tag):
    email = '%s-%s@example.com' % (tag, uuid.uuid4().hex[:8])
    status, _, body = call('POST', '/api/auth/register', body={
        'first_name': tag, 'last_name': 'Live', 'email': email, 'password': PASSWORD})
    if status not in (200, 201):
        sys.exit('register %s: %s %s' % (tag, status, body[:300]))
    token, uid = login(email)
    return token, uid, email


def wait_for_server():
    for _ in range(180):
        try:
            urllib.request.urlopen(args.base + '/api/health', timeout=2)
            return True
        except urllib.error.HTTPError:
            return True  # answering at all is enough
        except Exception:
            time.sleep(1)
    return False


def wait_job(job_id, token):
    for _ in range(120):
        _, _, body = call('GET', '/api/v1/imports/' + job_id, token)
        job = data(body) or {}
        if job.get('status') in ('completed', 'failed'):
            return job
        time.sleep(0.5)
    return {}


def upload_csv(path, token, text):
    boundary = 'x' + uuid.uuid4().hex
    raw = (('--%s\r\nContent-Disposition: form-data; name="file"; filename="rows.csv"\r\n'
            'Content-Type: text/csv\r\n\r\n' % boundary).encode() + text.encode() +
           ('\r\n--%s--\r\n' % boundary).encode())
    return call('POST', path, token, raw=raw,
                ctype='multipart/form-data; boundary=' + boundary)


# --------------------------------------------------------------- the checks

def check_optimistic_locking(alice):
    """Two people saving the same record: the second silently overwrote the
    first. Models carried a version nothing checked."""
    lots = '/api/v1/lots'
    title = 'Lot ' + uuid.uuid4().hex[:6]
    status, _, body = call('POST', lots, alice, {'title': title, 'current_bid': 100})
    lot = data(body) or {}
    check('If-Match: a new row starts at version 1', status == 201 and lot.get('version') == 1, (status, body[:200]))
    lot_id = lot.get('id', 'missing')

    _, headers, _ = call('GET', lots + '/' + lot_id, alice)
    check('If-Match: a read carries the version as an ETag', headers.get('etag') == 'W/"1"', headers.get('etag'))

    status, headers, body = call('PATCH', lots + '/' + lot_id, alice, {'current_bid': 150}, {'If-Match': 'W/"1"'})
    check('If-Match: a write on the current version lands', status == 200 and headers.get('etag') == 'W/"2"',
          (status, headers.get('etag'), body[:200]))

    status, _, body = call('PATCH', lots + '/' + lot_id, alice, {'current_bid': 160}, {'If-Match': 'W/"1"'})
    detail = err(body)
    check('If-Match: a stale version is 409 naming the current one',
          status == 409 and detail.get('code') == 'VERSION_CONFLICT'
          and (detail.get('details') or {}).get('current_version') == 2, (status, body[:200]))

    status, _, body = call('GET', lots + '/' + lot_id, alice)
    check('If-Match: the refused write changed nothing', (data(body) or {}).get('current_bid') == 150, body[:200])

    # Twenty writers on one version: exactly one may land. Counted per process,
    # a second replica would let a second one through.
    outcomes, lock, gate = [], threading.Lock(), threading.Barrier(20)

    def bid(i):
        gate.wait()
        status, _, _ = call('PATCH', lots + '/' + lot_id, alice, {'current_bid': 200 + i}, {'If-Match': 'W/"2"'})
        with lock:
            outcomes.append(status)

    threads = [threading.Thread(target=bid, args=(i,)) for i in range(20)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    check('If-Match: of 20 simultaneous writes on one version, one lands',
          outcomes.count(200) == 1 and outcomes.count(409) == 19, sorted(outcomes))
    return lot_id


def check_lists_and_exports(alice, admin):
    """Every CSV export was an empty 200: FindInBatches hands its callback a
    fresh session, and the rows were re-read from it. A sort ahead of the key it
    pages by repeated rows."""
    lots = '/api/v1/lots'
    title = 'Export ' + uuid.uuid4().hex[:6]
    status, _, body = call('POST', lots, alice, {'title': title, 'current_bid': 7})
    lot_id = (data(body) or {}).get('id', 'missing')

    status, _, body = call('GET', lots + '?search=' + urllib.parse.quote(title.split()[1]), alice)
    check('List: search finds the row', status == 200 and lot_id in ids(data(body)), (status, body[:200]))
    status, _, body = call('GET', lots + '?title=' + urllib.parse.quote(title), alice)
    check('List: a column filter narrows to it', status == 200 and ids(data(body)) == [lot_id], (status, body[:200]))
    status, _, body = call('GET', lots + '?sort_by=password_hash', alice)
    check('List: a column outside the whitelist is not sorted by', status != 500, (status, body[:160]))

    status, _, body = call('GET', lots + '/export?format=csv', alice)
    check('Export: the CSV has the row', status == 200 and title.encode() in body, (status, body[:200]))
    status, _, body = call('GET', lots + '/export?format=xlsx', alice)
    check('Export: the XLSX is a workbook', status == 200 and body[:2] == b'PK', status)

    status, _, body = call('POST', lots + '/bulk', admin,
                           {'action': 'patch', 'ids': [lot_id], 'patch': {'current_bid': 999}})
    _, _, after = call('GET', lots + '/' + lot_id, alice)
    check('Bulk: patch writes a whitelisted column',
          status == 200 and (data(after) or {}).get('current_bid') == 999, (status, body[:200]))
    status, _, body = call('POST', lots + '/bulk', admin,
                           {'action': 'patch', 'ids': [lot_id], 'patch': {'version': 1}})
    check('Bulk: patch refuses a column that is not writable', status in (400, 422), (status, body[:200]))
    status, _, body = call('POST', lots + '/bulk', admin, {'action': 'delete', 'ids': [lot_id]})
    after_status, _, _ = call('GET', lots + '/' + lot_id, alice)
    check('Bulk: delete removes the row', status == 200 and after_status == 404, (status, after_status))


def check_ownership(alice, alice_id, bob, bob_id, staff, admin):
    """An owned resource checked ownership on list, read, update and delete, and
    on nothing else: a second account printed another user's record as a PDF and
    rewrote it with PATCH. The export and the bulk route were the same hole."""
    notes = '/api/v1/notes'
    status, _, body = call('POST', notes, alice, {'body': 'alice private', 'user_id': bob_id})
    note = data(body) or {}
    note_id = note.get('id', 'missing')
    check('Owned: the owner is the caller, whatever the body says',
          status == 201 and note.get('user_id') == alice_id, (status, body[:200]))

    for label, (method, path, payload) in {
        'read it': ('GET', '/' + note_id, None),
        'print it': ('GET', '/' + note_id + '/pdf', None),
        'patch it': ('PATCH', '/' + note_id, {'body': 'bob was here'}),
        'replace it': ('PUT', '/' + note_id, {'body': 'bob was here'}),
    }.items():
        status, _, body = call(method, notes + path, bob, payload)
        check('Owned: another account cannot %s (404)' % label, status == 404, (status, body[:160]))

    status, _, body = call('GET', notes, bob)
    check('Owned: it is not in another account\'s list', status == 200 and b'alice private' not in body, status)
    status, _, body = call('GET', notes + '/export?format=csv', bob)
    check('Owned: it is not in another account\'s export', status == 200 and b'alice private' not in body, status)

    # Staff holds notes.delete and is not ADMIN: the permission lets it reach
    # the route, and ownership still decides the row.
    status, _, body = call('DELETE', notes + '/' + note_id, staff)
    check('Owned: a staff account with the permission still cannot delete another user\'s row',
          status == 404, (status, body[:160]))
    status, _, body = call('POST', notes + '/bulk', staff, {'action': 'delete', 'ids': [note_id]})
    after_status, _, after = call('GET', notes + '/' + note_id, alice)
    check('Owned: a staff bulk delete skips rows it does not own',
          status == 200 and after_status == 200 and (data(after) or {}).get('body') == 'alice private',
          (status, after_status))

    status, _, body = call('GET', notes + '/' + note_id + '/pdf', alice)
    check('Owned: the owner can print it', status == 200 and body[:4] == b'%PDF', status)
    status, _, body = call('GET', notes, admin)
    check('Owned: ADMIN sees every owner\'s rows', status == 200 and note_id in ids(data(body)), status)


def check_relations_and_money(alice):
    """A mistyped many-to-many id used to be dropped, so a PUT naming one
    answered 200 and left the row with no tags. Money has to stay exact, and an
    encrypted column must not be readable in the database."""
    status, _, body = call('POST', '/api/v1/categories', alice, {'name': 'Hardware ' + uuid.uuid4().hex[:4]})
    category = (data(body) or {}).get('id', 'missing')
    tags = []
    for name in ('red', 'blue'):
        _, _, body = call('POST', '/api/v1/tags', alice, {'name': name + '-' + uuid.uuid4().hex[:4]})
        tags.append((data(body) or {}).get('id'))

    products = '/api/v1/products'
    status, _, body = call('POST', products, alice, {
        'name': 'Hammer ' + uuid.uuid4().hex[:4], 'price': {'amount': 1999, 'currency': 'USD'},
        'secret': 'supplier-cost-7', 'category_id': category, 'tag_ids': tags})
    product = data(body) or {}
    product_id = product.get('id', 'missing')
    check('Relations: a create with two tags returns both, preloaded',
          status == 201 and len(product.get('tags') or []) == 2, (status, body[:300]))
    check('Money: the amount and currency round-trip',
          (product.get('price') or {}).get('amount') == 1999
          and (product.get('price') or {}).get('currency') == 'USD', product.get('price'))
    check('Encryption: the value reads back through the API', product.get('secret') == 'supplier-cost-7',
          product.get('secret'))
    stored = psql("select secret from products where id = '%s'" % product_id)
    check('Encryption: and is not readable in the database', stored and 'supplier-cost-7' not in stored, stored[:60])

    status, _, body = call('PATCH', products + '/' + product_id, alice, {'tag_ids': [tags[1]]})
    check('Relations: PATCH replaces the set', status == 200 and ids((data(body) or {}).get('tags')) == [tags[1]],
          (status, body[:250]))
    status, _, body = call('PUT', products + '/' + product_id, alice,
                           {'name': 'Renamed', 'tag_ids': [str(uuid.uuid4())]})
    _, _, after = call('GET', products + '/' + product_id, alice)
    got = data(after) or {}
    check('Relations: an id that matches nothing is refused, and nothing changes',
          status == 422 and ids(got.get('tags')) == [tags[1]] and got.get('name') != 'Renamed',
          (status, body[:200], ids(got.get('tags')), got.get('name')))
    return category, product_id


def check_line_items(alice):
    """A replace used to append: the old lines stayed, so an invoice grew every
    time it was saved."""
    invoices = '/api/v1/invoices'
    status, _, body = call('POST', invoices, alice, {'number': 'INV-' + uuid.uuid4().hex[:6], 'items': [
        {'description': 'nails', 'qty': 10, 'unit_rate': 0.5},
        {'description': 'glue', 'qty': 1, 'unit_rate': 4.25}]})
    invoice = data(body) or {}
    invoice_id = invoice.get('id', 'missing')
    check('Items: a create writes both lines', status == 201 and len(invoice.get('items') or []) == 2,
          (status, body[:300]))

    status, _, body = call('PUT', invoices + '/' + invoice_id, alice,
                           {'number': invoice.get('number'), 'items': [
                               {'description': 'screws', 'qty': 3, 'unit_rate': 1.0}]},
                           {'If-Match': 'W/"1"'})
    got = data(body) or {}
    check('Items: a replace replaces, and does not append',
          status == 200 and [i.get('description') for i in got.get('items') or []] == ['screws'],
          (status, body[:300]))
    live = psql("select count(*) from invoice_items where invoice_id = '%s' and deleted_at is null" % invoice_id)
    check('Items: one live line is left in the table', live == '1', live)

    status, _, body = call('PUT', invoices + '/' + invoice_id, alice,
                           {'number': invoice.get('number'), 'items': [
                               {'description': 'lost', 'qty': 1, 'unit_rate': 1.0}]},
                           {'If-Match': 'W/"1"'})
    _, _, after = call('GET', invoices + '/' + invoice_id, alice)
    check('Items: a stale write is refused and leaves the lines alone',
          status == 409 and [i.get('description') for i in (data(after) or {}).get('items') or []] == ['screws'],
          (status, after[:250]))


def check_tree(alice):
    """A move has to carry its subtree, and a reorder has to keep the parent the
    node has. Dragging a node that predated --tree once added one to the depth of
    every row in the table."""
    categories = '/api/v1/categories'

    def node(name, parent=None):
        payload = {'name': name + '-' + uuid.uuid4().hex[:4]}
        if parent:
            payload['parent_id'] = parent
        _, _, body = call('POST', categories, alice, payload)
        return (data(body) or {}).get('id', 'missing')

    top, other = node('Top'), node('Other')
    middle = node('Middle', top)
    bottom = node('Bottom', middle)

    status, _, body = call('GET', categories + '/tree', alice)
    tree = data(body) or []
    root = next((n for n in tree if n.get('id') == top), {})
    check('Tree: the whole tree comes back nested',
          status == 200 and (root.get('children') or [{}])[0].get('id') == middle, (status, body[:250]))

    status, _, body = call('GET', categories + '/' + bottom + '/breadcrumbs', alice)
    check('Tree: breadcrumbs run root to node', status == 200 and ids(data(body)) == [top, middle, bottom],
          (status, body[:250]))

    status, _, body = call('PATCH', categories + '/' + middle + '/move', alice, {'parent_id': other, 'position': 0})
    moved = psql("select concat(path, '|', depth) from categories where id = '%s'" % bottom)
    check('Tree: a move carries the subtree with it',
          status == 200 and moved == '/%s/%s/%s/|2' % (other, middle, bottom), (status, moved))

    status, _, body = call('PATCH', categories + '/' + bottom + '/move', alice, {'position': 5})
    kept = psql("select concat(parent_id, '|', position) from categories where id = '%s'" % bottom)
    check('Tree: a move with no parent keeps the one it has', status == 200 and kept == '%s|5' % middle,
          (status, kept))

    status, _, body = call('PATCH', categories + '/' + other + '/move', alice, {'parent_id': bottom, 'position': 0})
    check('Tree: a node cannot be moved under its own descendant', status == 422, (status, body[:200]))


def check_public_surface(alice, admin, category, product_id):
    """The public endpoints are anonymous and guarded by an API key, and they
    publish an allowlist rather than the model."""
    status, _, body = call('POST', '/api/v1/api-keys', alice, {'name': 'live-' + RUN[:6], 'kind': 'publishable'})
    issued = data(body) or {}
    key = issued.get('token') if isinstance(issued.get('token'), str) else ''
    check('Public: a publishable key is issued', status in (200, 201) and key, (status, body[:200]))
    headers = {'X-API-Key': key}

    status, _, body = call('GET', '/api/v1/public/products')
    check('Public: no key is 401', status == 401, status)
    status, _, body = call('GET', '/api/v1/public/products', headers=headers)
    rows = data(body) or []
    check('Public: the list answers with a key', status == 200 and rows, (status, body[:200]))
    check('Public: it publishes the allowlist, not the model',
          all('secret' not in r and 'category_id' not in r for r in rows), body[:250])

    status, _, body = call('GET', '/api/v1/public/products?category_id=' + category, headers=headers)
    check('Public: a foreign key filters the list', status == 200 and product_id in ids(data(body)),
          (status, body[:200]))
    status, _, body = call('GET', '/api/v1/public/categories/tree', headers=headers)
    check('Public: the category tree is published', status == 200 and data(body), (status, body[:200]))
    status, _, body = call('GET', '/api/v1/public/products/no-such-slug', headers=headers)
    check('Public: an unknown slug is 404', status == 404, status)

    # Archived rows leave the public surface. The list is cached by URL for a
    # minute, so this reads a fresh key.
    status, _, body = call('POST', '/api/v1/products/bulk', admin, {'action': 'archive', 'ids': [product_id]})
    _, _, after = call('GET', '/api/v1/public/products?fresh=' + uuid.uuid4().hex, headers=headers)
    check('Public: an archived row leaves the list', status == 200 and product_id not in ids(data(after)),
          (status, body[:160]))
    call('POST', '/api/v1/products/bulk', admin, {'action': 'restore', 'ids': [product_id]})


def check_csv_import(alice, admin, alice_id, alice_email, bob, bob_id):
    """The importer took the owner from a CSV column, so any account could file
    rows under another name, and it created a user for every unknown email."""
    status, _, body = call('GET', '/api/v1/lots/import/template', alice)
    check('Import: a template names the columns', status == 200 and b'title' in body, (status, body[:120]))

    marker = 'imported-' + uuid.uuid4().hex[:6]
    status, _, body = upload_csv('/api/v1/lots/import', alice,
                                 'title,current_bid\r\n%s,10\r\n%s-2,20\r\n' % (marker, marker))
    job = wait_job((data(body) or {}).get('job_id', 'missing'), alice)
    check('Import: the upload is a job that completes with both rows',
          status == 202 and job.get('status') == 'completed' and job.get('created') == 2, (status, job))
    count = psql("select count(*) from lots where title like '%s%%'" % marker)
    check('Import: the rows are in the table', count == '2', count)

    # An owned resource: the rows belong to whoever imports them, and only an
    # ADMIN may name another owner.
    _, _, body = call('GET', '/api/v1/notes/import/template', bob)
    header = body.decode(errors='replace').strip()
    column = 'user' if 'user' in header.split(',') else 'user_id'
    value = alice_email if column == 'user' else alice_id

    body_text = 'filed by another account as alice ' + RUN[:6]
    _, _, response = upload_csv('/api/v1/notes/import', bob,
                                'body,%s\r\n%s,%s\r\n' % (column, body_text, value))
    job = wait_job((data(response) or {}).get('job_id', 'missing'), bob)
    owner = psql("select user_id from notes where body = '%s'" % body_text)
    check('Import (owned): an ordinary account\'s rows are its own, whatever the CSV says',
          job.get('status') == 'completed' and owner == bob_id, (job, owner, bob_id))

    admin_text = 'filed for alice by an admin ' + RUN[:6]
    _, _, response = upload_csv('/api/v1/notes/import', admin,
                                'body,%s\r\n%s,%s\r\n' % (column, admin_text, value))
    job = wait_job((data(response) or {}).get('job_id', 'missing'), admin)
    owner = psql("select user_id from notes where body = '%s'" % admin_text)
    check('Import (owned): an ADMIN may name the owner',
          job.get('status') == 'completed' and owner == alice_id, (job, owner))


def check_error_codes(alice, alice_email):
    """One code, one status, in the running API.

    Each of these was inconsistent before v3.232.0: VALIDATION_ERROR arrived as
    422 from some handlers and 400 from others, INVALID_TOKEN as both 401 and 400,
    and an expired one-time link borrowed a 401 it had no business with. A scan of
    the templates cannot see any of that; this can.
    """
    def code_of(body):
        return (err(body) or {}).get('code')

    # A field that fails validation, and a body that is not JSON at all. Both are
    # 422 from a generated handler: gin binds and validates in one step.
    status, _, body = call('POST', '/api/v1/lots', alice, {'current_bid': 5})
    check('Codes: a missing field is 422 VALIDATION_ERROR',
          status == 422 and code_of(body) == 'VALIDATION_ERROR', (status, code_of(body)))
    status, _, body = call('POST', '/api/v1/lots', alice, raw=b'{not json')
    check('Codes: an unparseable body is 422 VALIDATION_ERROR',
          status == 422 and code_of(body) == 'VALIDATION_ERROR', (status, code_of(body)))

    # An id nobody has.
    status, _, body = call('GET', '/api/v1/lots/' + str(uuid.uuid4()), alice)
    check('Codes: an unknown id is 404 NOT_FOUND',
          status == 404 and code_of(body) == 'NOT_FOUND', (status, code_of(body)))

    # No credentials, and credentials that are not credentials.
    status, _, body = call('GET', '/api/v1/lots')
    check('Codes: no token is 401 UNAUTHORIZED',
          status == 401 and code_of(body) == 'UNAUTHORIZED', (status, code_of(body)))
    status, _, body = call('GET', '/api/v1/lots', 'not-a-jwt')
    check('Codes: a bad token is 401 UNAUTHORIZED',
          status == 401 and code_of(body) == 'UNAUTHORIZED', (status, code_of(body)))

    # A one-time link is not a credential: 400, and a code of its own, so a client
    # offers a new link rather than a sign-in form.
    status, _, body = call('POST', '/api/v1/auth/verify-email', body={'token': 'nope'})
    check('Codes: a bad verification link is 400 INVALID_LINK',
          status == 400 and code_of(body) == 'INVALID_LINK', (status, code_of(body)))
    status, _, body = call('POST', '/api/v1/auth/reset-password',
                           body={'token': 'nope', 'password': 'SuperSecret123!'})
    check('Codes: a bad reset link is 400 INVALID_LINK',
          status == 400 and code_of(body) == 'INVALID_LINK', (status, code_of(body)))

    # Credentials that are wrong, and an address that is taken.
    status, _, body = call('POST', '/api/v1/auth/login',
                           body={'email': alice_email, 'password': 'definitely-wrong'})
    check('Codes: a wrong password is 401 INVALID_CREDENTIALS',
          status == 401 and code_of(body) == 'INVALID_CREDENTIALS', (status, code_of(body)))
    status, _, body = call('POST', '/api/v1/auth/register',
                           body={'first_name': 'Taken', 'last_name': 'Address', 'email': alice_email,
                                 'password': 'SuperSecret123!'})
    check('Codes: a taken address is 409 EMAIL_EXISTS',
          status == 409 and code_of(body) == 'EMAIL_EXISTS', (status, code_of(body)))

    # The envelope itself: the same three keys whatever answered, because a client
    # reads one shape.
    status, _, body = call('GET', '/api/v1/lots/' + str(uuid.uuid4()), alice)
    envelope = err(body)
    check('Codes: the envelope is code, message and optional details',
          isinstance(envelope, dict) and envelope.get('code') and envelope.get('message')
          and set(envelope) <= {'code', 'message', 'details'}, envelope)

    # And details carries the fields, which is what a form needs to mark inputs.
    status, _, body = call('POST', '/api/v1/invoices', alice, {'number': ''})
    if status == 422:
        check('Codes: a validation failure may carry per-field details',
              isinstance(err(body).get('details', {}), dict), err(body))


def check_write_errors(alice):
    """Five write paths bound the error and returned a constant 500, so a rule
    enforced in a GORM hook reached neither the client nor the log."""
    status, _, body = call('POST', '/api/v1/lots', alice, {'current_bid': 5})
    check('Errors: a missing required field is 422, not 500', status == 422, (status, body[:200]))
    status, _, body = call('PATCH', '/api/v1/lots/' + str(uuid.uuid4()), alice, {'title': 'x'})
    check('Errors: an unknown id is 404, not 500', status == 404, (status, body[:200]))
    status, _, body = call('GET', '/api/v1/lots/' + str(uuid.uuid4()), alice)
    check('Errors: so is a read of one', status == 404, status)


def main():
    global args
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('base', nargs='?', default='http://localhost:8080',
                        help='the running API, e.g. http://localhost:8080')
    parser.add_argument('--db', default='grit', help='database name')
    parser.add_argument('--db-user', default='grit')
    parser.add_argument('--db-host', default='localhost')
    parser.add_argument('--db-port', default='5432')
    parser.add_argument('--db-password', default='')
    parser.add_argument('--psql-container', default='',
                        help='run the database client inside this docker container instead of on the host')
    parser.add_argument('--engine', default='postgres', choices=['postgres', 'mysql', 'sqlite'],
                        help='which engine is behind the API: it decides how the checks that read the '
                             'database directly are run. For sqlite, --db is the file path.')
    args = parser.parse_args()
    args.base = args.base.rstrip('/')
    args.psql_env = None
    if args.db_password:
        import os
        args.psql_env = dict(os.environ, PGPASSWORD=args.db_password)

    if not wait_for_server():
        sys.exit('the API at %s never answered /api/health' % args.base)

    alice, alice_id, alice_email = register('alice')
    bob, bob_id, _ = register('bob')
    _, admin_id, admin_email = register('admin')
    psql("update users set role = 'ADMIN' where id = '%s'" % admin_id)
    admin, _ = login(admin_email)

    # A staff account that holds one permission and is not ADMIN: the role gate
    # lets it reach the route, and ownership still decides the row.
    staff_token, staff_id, staff_email = register('staff')
    status, _, body = call('POST', '/api/v1/roles', admin,
                           {'name': 'Note Cleaner ' + RUN[:6], 'grants': ['notes.delete']})
    role_id = (data(body) or {}).get('id')
    call('PUT', '/api/v1/users/%s/roles' % staff_id, admin, {'role_ids': [role_id]})
    staff, _ = login(staff_email)
    check('Setup: a staff account holds notes.delete and nothing else', status in (200, 201) and role_id,
          (status, body[:160]))

    check_optimistic_locking(alice)
    check_lists_and_exports(alice, admin)
    check_ownership(alice, alice_id, bob, bob_id, staff, admin)
    category, product_id = check_relations_and_money(alice)
    check_line_items(alice)
    check_tree(alice)
    check_public_surface(alice, admin, category, product_id)
    check_csv_import(alice, admin, alice_id, alice_email, bob, bob_id)
    check_write_errors(alice)
    check_error_codes(alice, alice_email)

    width = max(len(name) for name, _, _ in results)
    failed = 0
    for name, ok, detail in results:
        if not ok:
            failed += 1
        print('%-*s  %s  %s' % (width, name, 'PASS' if ok else 'FAIL', '' if ok else detail))
    print('\n%d/%d passed' % (len(results) - failed, len(results)))
    if failed:
        sys.exit(1)


if __name__ == '__main__':
    main()
