#!/usr/bin/env python3
"""Tenancy, custom roles and impersonation, together, against a running app.

This is the reviewer's Project 1 from grit-hardening-and-devx-plan.md, kept as a
fixture rather than run once: every check here crosses two subsystems that had only
been exercised apart, on the theory that intersections are where the next bugs
live. It was right. The first run found three:

  * the tenant middleware was mounted on the protected group alone, so every
    DELETE, every bulk route and every admin-only route touching a tenant-owned
    model answered 500: the request never resolved an organization at all;
  * "no active organization" arrived as that 500 rather than as something a client
    could act on;
  * a role held through an organization membership granted nothing, because the
    middleware put it on the request context and nothing read it, so being an
    administrator of one organization meant being one in all of them.

Needs a project scaffolded with the multitenant and impersonate plugins, a Deal
(--tenant-owned) and a Contact (--tenant-owned --owned-by user). See
.github/workflows/live.yml, which builds exactly that.

    python3 tests/live/verify_tenancy.py http://localhost:8080 \\
        --db grit --db-user grit --db-password grit
"""
import argparse
import http.cookiejar
import json
import os
import random
import subprocess
import sys
import urllib.error
import urllib.request

args = None
results = []
RUN = "%06d" % random.randint(0, 999999)


def check(name, ok, detail=""):
    results.append((name, bool(ok), detail))


def csrf_of(jar):
    """The grit_csrf cookie. An unsafe request has to echo it back in a header."""
    for cookie in jar or []:
        if cookie.name == "grit_csrf":
            return cookie.value
    return None


def call(method, path, token=None, body=None, headers=None, jar=None):
    data = json.dumps(body).encode() if body is not None else None
    request = urllib.request.Request(args.base + path, data=data, method=method)
    if data is not None:
        request.add_header("Content-Type", "application/json")
    if token:
        request.add_header("Authorization", "Bearer " + token)
    if jar is not None and method not in ("GET", "HEAD", "OPTIONS"):
        value = csrf_of(jar)
        if value:
            request.add_header("X-CSRF-Token", value)
    for key, value in (headers or {}).items():
        request.add_header(key, value)
    opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar)) \
        if jar is not None else urllib.request.build_opener()
    try:
        with opener.open(request, timeout=30) as response:
            raw = response.read()
            return response.status, (json.loads(raw) if raw else {})
    except urllib.error.HTTPError as error:
        raw = error.read()
        try:
            return error.code, (json.loads(raw) if raw else {})
        except ValueError:
            return error.code, {"raw": raw.decode(errors="replace")[:160]}


def code(body):
    return (body or {}).get("error", {}).get("code")


def psql(sql):
    """One statement, for the only thing HTTP cannot do: promote an ADMIN."""
    command = ["psql", "-h", args.db_host, "-p", str(args.db_port), "-U", args.db_user,
               "-d", args.db, "-v", "ON_ERROR_STOP=1", "-q", "-c", sql]
    if args.psql_container:
        command = ["docker", "exec", "-i", args.psql_container, "psql", "-U", args.db_user,
                   "-d", args.db, "-v", "ON_ERROR_STOP=1", "-q", "-c", sql]
    environment = dict(os.environ, PGPASSWORD=args.db_password) if args.db_password else None
    return subprocess.run(command, capture_output=True, text=True, env=environment)


def register(tag):
    email = "%s-%s@example.com" % (tag, RUN)
    status, body = call("POST", "/api/v1/auth/register", body={
        "first_name": tag.capitalize(), "last_name": "Seams",
        "email": email, "password": "SuperSecret123!"})
    if status not in (200, 201):
        sys.exit("register %s failed: %s %s" % (tag, status, body))
    data = body["data"]
    return data["tokens"]["access_token"], data["user"]["id"], email


def login(email, jar=None):
    status, body = call("POST", "/api/v1/auth/login",
                        body={"email": email, "password": "SuperSecret123!"}, jar=jar)
    if status != 200:
        sys.exit("login %s failed: %s %s" % (email, status, body))
    return body["data"]["tokens"]["access_token"]


def main():
    global args
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("base", nargs="?", default="http://localhost:8080")
    parser.add_argument("--db", default="grit")
    parser.add_argument("--db-user", default="grit")
    parser.add_argument("--db-host", default="localhost")
    parser.add_argument("--db-port", default="5432")
    parser.add_argument("--db-password", default="")
    parser.add_argument("--psql-container", default="")
    args = parser.parse_args()
    args.base = args.base.rstrip("/")

    # ── the cast ──────────────────────────────────────────────────────────────
    alice, alice_id, alice_email = register("alice")   # owns organization A
    bob, bob_id, bob_email = register("bob")           # owns organization B
    carol, carol_id, carol_email = register("carol")   # a member of both
    dave, dave_id, dave_email = register("dave")       # a member through a role
    admin, admin_id, admin_email = register("admin")   # platform ADMIN

    status, body = call("POST", "/api/v1/organizations", alice, {"name": "Org A " + RUN})
    org_a = (body.get("data") or {}).get("id")
    check("Tenancy: an organization can be created", status in (200, 201) and org_a, (status, body))

    status, body = call("POST", "/api/v1/organizations", bob, {"name": "Org B " + RUN})
    org_b = (body.get("data") or {}).get("id")
    check("Tenancy: and a second one, by somebody else", status in (200, 201) and org_b, (status, body))

    status, body = call("POST", "/api/v1/organizations/%s/members" % org_a, alice,
                        {"user_id": carol_id})
    check("Tenancy: a member can be added", status in (200, 201), (status, body))

    # ── rows belong to an organization ────────────────────────────────────────
    status, created = call("POST", "/api/v1/deals", alice,
                           {"title": "A deal " + RUN,
                            "amount": {"amount": 5000, "currency": "USD"}, "stage": "new"},
                           headers={"X-Organization-ID": org_a})
    deal_a = (created.get("data") or {}).get("id")
    check("Tenancy: a tenant-owned row is created in the active organization",
          status in (200, 201) and deal_a, (status, created))

    status, body = call("POST", "/api/v1/deals", bob,
                        {"title": "B deal " + RUN,
                         "amount": {"amount": 7000, "currency": "USD"}, "stage": "new"},
                        headers={"X-Organization-ID": org_b})
    deal_b = (body.get("data") or {}).get("id")
    check("Tenancy: the other organization has its own", status in (200, 201), (status, body))

    status, body = call("GET", "/api/v1/deals", alice, headers={"X-Organization-ID": org_a})
    check("Tenancy: a list is scoped to the active organization",
          status == 200 and len(body.get("data") or []) == 1, (status, body.get("data")))

    # ── the header is not trusted ─────────────────────────────────────────────
    status, body = call("GET", "/api/v1/deals", bob, headers={"X-Organization-ID": org_a})
    check("Tenancy: naming an organization you are not in is 404",
          status == 404 and code(body) == "NOT_FOUND", (status, code(body)))

    status, body = call("GET", "/api/v1/deals/%s" % deal_a, bob,
                        headers={"X-Organization-ID": org_b})
    check("Tenancy: another organization's row by id is 404", status == 404, (status, code(body)))

    # ── no active organization is an answer, not a crash ──────────────────────
    status, body = call("GET", "/api/v1/deals", dave)
    check("Tenancy: no organization is 400 NO_ORGANIZATION, not a 500",
          status == 400 and code(body) == "NO_ORGANIZATION", (status, code(body)))

    status, body = call("POST", "/api/v1/deals", dave,
                        {"title": "x", "amount": {"amount": 1, "currency": "USD"}, "stage": "new"})
    check("Tenancy: and the same on a write",
          status == 400 and code(body) == "NO_ORGANIZATION", (status, code(body)))

    call("POST", "/api/v1/organizations/%s/members" % org_b, bob, {"user_id": carol_id})
    carol = login(carol_email)
    status, body = call("GET", "/api/v1/deals", carol)
    check("Tenancy: belonging to two and naming neither is NO_ORGANIZATION too",
          status == 400 and code(body) == "NO_ORGANIZATION", (status, code(body)))

    # ── a platform ADMIN is exempt from ownership, not from tenancy ───────────
    psql("update users set role = 'ADMIN' where id = '%s'" % admin_id)
    admin = login(admin_email)
    status, body = call("GET", "/api/v1/deals", admin)
    check("Tenancy: an ADMIN outside every organization is told so, not given a 500",
          status == 400 and code(body) == "NO_ORGANIZATION", (status, code(body)))

    call("POST", "/api/v1/organizations/%s/members" % org_a, alice, {"user_id": admin_id})
    admin = login(admin_email)
    status, body = call("GET", "/api/v1/deals", admin, headers={"X-Organization-ID": org_a})
    check("Tenancy: an ADMIN inside one organization sees only its rows",
          status == 200 and len(body.get("data") or []) == 1, (status, body.get("data")))

    # ── an organization and an owner on the same row ──────────────────────────
    carol = login(carol_email)
    status, body = call("POST", "/api/v1/contacts", carol,
                        {"name": "Carol's contact", "email": "c-%s@example.com" % RUN},
                        headers={"X-Organization-ID": org_a})
    contact_id = (body.get("data") or {}).get("id")
    check("Tenancy: an owned, tenant-scoped row can be created",
          status in (200, 201) and contact_id, (status, body))

    status, body = call("GET", "/api/v1/contacts", alice, headers={"X-Organization-ID": org_a})
    check("Tenancy: another member of the organization does not see somebody else's owned row",
          status == 200 and len(body.get("data") or []) == 0, (status, body.get("data")))

    status, body = call("GET", "/api/v1/contacts/%s" % contact_id, alice,
                        headers={"X-Organization-ID": org_a})
    check("Tenancy: and by id it is 404, not 403", status == 404, (status, code(body)))

    status, body = call("GET", "/api/v1/contacts", admin, headers={"X-Organization-ID": org_a})
    check("Tenancy: an ADMIN in the organization does see it",
          status == 200 and len(body.get("data") or []) == 1, (status, body.get("data")))

    # ── a role held through a membership grants inside that organization only ──
    status, body = call("POST", "/api/v1/roles", admin,
                        {"name": "Org A Deleter " + RUN, "grants": ["deals.delete"]},
                        headers={"X-Organization-ID": org_a})
    org_role = (body.get("data") or {}).get("id")
    check("Roles: a role to hold through a membership can be created",
          status in (200, 201) and org_role, (status, body))

    call("POST", "/api/v1/organizations/%s/members" % org_a, alice,
         {"user_id": dave_id, "role_id": org_role})
    call("POST", "/api/v1/organizations/%s/members" % org_b, bob, {"user_id": dave_id})
    dave = login(dave_email)

    status, body = call("DELETE", "/api/v1/deals/%s" % deal_a, dave,
                        headers={"X-Organization-ID": org_a})
    check("Roles: the membership role grants the delete in its own organization",
          status in (200, 204), (status, body))

    status, body = call("DELETE", "/api/v1/deals/%s" % deal_b, dave,
                        headers={"X-Organization-ID": org_b})
    check("Roles: and grants nothing in the other one",
          status == 403 and code(body) == "FORBIDDEN", (status, code(body)))

    status, body = call("GET", "/api/v1/deals/%s" % deal_b, bob,
                        headers={"X-Organization-ID": org_b})
    check("Roles: so the other organization's row is still there", status == 200,
          (status, code(body)))

    # ── impersonation, inside a tenant ───────────────────────────────────────
    jar = http.cookiejar.CookieJar()
    call("GET", "/api/health", jar=jar)       # issues the CSRF cookie
    login(admin_email, jar=jar)
    status, body = call("POST", "/api/v1/admin/impersonate/%s" % carol_id, admin, jar=jar)
    check("Impersonation: an ADMIN can start", status == 200, (status, body))

    status, body = call("GET", "/api/v1/contacts", headers={"X-Organization-ID": org_a}, jar=jar)
    check("Impersonation: the owned row belongs to the impersonated user, not the admin",
          status == 200 and len(body.get("data") or []) == 1, (status, body.get("data")))

    status, body = call("GET", "/api/v1/deals", headers={"X-Organization-ID": org_b}, jar=jar)
    check("Impersonation: and their organizations are the ones reachable", status == 200,
          (status, code(body)))

    status, body = call("POST", "/api/v1/auth/impersonate/stop", jar=jar)
    check("Impersonation: it can be stopped", status == 200, (status, body))

    status, body = call("GET", "/api/v1/users", jar=jar)
    check("Impersonation: and the admin is themselves again", status == 200,
          (status, code(body)))

    width = max(len(name) for name, _, _ in results)
    failed = 0
    for name, ok, detail in results:
        if not ok:
            failed += 1
        print("%-*s  %s  %s" % (width, name, "PASS" if ok else "FAIL", "" if ok else detail))
    print("\n%d/%d passed" % (len(results) - failed, len(results)))
    if failed:
        sys.exit(1)


if __name__ == "__main__":
    main()
