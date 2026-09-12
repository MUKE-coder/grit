# Live checks

`verify.py` drives a running generated application over HTTP and checks the
database behind it. Every check in it is a bug that was found by hand once, in a
real application, and then fixed: an export that returned an empty file with a
200, a list that handed one user another user's rows, a save that silently
overwrote a newer one, a replace that appended instead of replacing, a mistyped
relation id that emptied a set.

Unit tests over generated source cannot catch those. They are behaviours of the
running thing, and several are behaviours of Postgres.

## In CI

`.github/workflows/live.yml` runs this on Postgres 15, 16 and 17 on every push:
it scaffolds a project, generates the resource shapes the bugs lived in,
migrates, starts the server, and runs this suite against it.

## By hand

Scaffold a project and generate the fixture resources:

```bash
grit new app --api && cd app
grit generate resource Tag --fields "name:string"
grit generate resource Lot --fields "title:string,current_bid:int"
grit generate resource Note --fields "body:text" --owned-by user
grit generate resource Category --fields "name:string,slug:slug" --tree --public
grit generate resource Product \
  --fields "name:string,slug:slug,price:money,secret:text:encrypted,category:belongs_to:Category,tags:many_to_many:Tag" \
  --public
grit generate resource Invoice --fields "number:string" \
  --items "InvoiceItem:description:string,qty:int,unit_rate:float"
```

Migrate, then start the API with a field-encryption key, or the
"not readable in the database" check passes for the wrong reason:

```bash
export FIELD_ENCRYPTION_KEY=$(openssl rand -base64 32)
grit migrate && grit start server
```

Then, from this repository:

```bash
# psql on the host
python tests/live/verify.py http://localhost:8080 \
  --db myapp --db-user grit --db-password "$PGPASSWORD"

# or psql inside the compose container
python tests/live/verify.py http://localhost:8080 \
  --db myapp --psql-container myapp-postgres
```

It prints one line per check and exits non-zero if any failed.

## Adding a check

Add one whenever a bug is found by using the application rather than by a test.
Name the behaviour, not the mechanism ("an archived row leaves the public list",
not "archived_at is filtered"), and write the comment so the next reader knows
what went wrong once: that sentence is the reason the check is allowed to cost
CI time.
