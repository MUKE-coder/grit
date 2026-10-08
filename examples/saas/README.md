# saas

A subscription product with a **customer portal** and an **operator console**,
built with Grit and committed so the claim can be checked rather than believed.

It answers one question people ask before starting a SaaS with Grit: *do my
customers log into the admin panel, or do they get their own area?* They get
their own area. This is what that looks like.

---

## The shape

One users table. One login. One API. Two front doors.

| | Lives in | Who signs in | What they see |
|---|---|---|---|
| **Customer portal** | `apps/web`, the `(app)` route group | a customer | their plan, their invoices, their usage |
| **Operator console** | `apps/admin` | you and your team | every plan, every subscription, every customer |

The portal is not a second API. `/account/billing` reads
`GET /api/subscriptions` and `GET /api/invoices`: the same endpoints the admin
panel reads, with the same JWT and the same permission checks.

**What makes it the customer's data and not everybody's is one flag.** The three
customer-facing resources were generated with `--owned-by user`:

```bash
grit generate resource Subscription --owned-by user --fields "..."
grit generate resource Invoice      --owned-by user --fields "..."
grit generate resource UsageRecord  --owned-by user --fields "..."
```

That puts `authz.ScopeOwned(ctx, query, "user_id")` on every list and
`authz.Owns(ctx, &item)` on every read, inside the service, where a job and a
command get it too. A customer asking for the whole table gets their own rows
back. ADMIN is exempt, which is what gives the operator console its view.

A customer is a plain `USER` holding three grants and nothing else:

```go
role.SetGrants([]string{"subscriptions.view", "invoices.view", "usage_records.view"})
```

The permission gets them through the door. The owner scope decides which rows
are behind it. Those are two separate mechanisms and the example needs both, so
it tests both.

---

## Run it

```bash
pnpm install
grit migrate
grit seed
grit start
```

The seed makes two customers so the isolation is visible rather than asserted:

| Sign in at `localhost:3000/login` | Plan | What they should see |
|---|---|---|
| `ada@example.com` / `customer123` | Team, $99/mo, 9 of 15 seats | three $99 invoices, three metrics |
| `basil@example.com` / `customer123` | Starter, $19/mo, trialing | three $19 invoices, two metrics |

Sign in as each and compare. Neither can see the other's anything. Then sign in
to the admin at `localhost:3001` as `admin@example.com` / `admin123` and you see
both, because that is the other door.

---

## What is generated and what is not

Four resources are generated and untouched: `Plan`, `Subscription`, `Invoice`,
`UsageRecord`. Everything in `apps/admin` is generated. So is every model,
service, handler, migration and schema.

Four files are hand-written:

| File | Why it is not generated |
|---|---|
| `apps/web/app/(app)/account/billing/page.tsx` | The customer's view of their plan and invoices |
| `apps/web/app/(app)/account/usage/page.tsx` | Metered usage, counted by the server with `?breakdown=metric` |
| `apps/api/internal/database/portal_demo_seeder.go` | Two customers whose rows agree across five tables |
| `apps/api/internal/handlers/portal_isolation_test.go` | **The point of the example** |

Plus two lines in `apps/web/app/(app)/layout.tsx` adding the sections to the
customer nav, which is what that file's own comment tells you to do.

---

## The test is the deliverable

`portal_isolation_test.go` is why this is checked in rather than described.
Data isolation in a README is a claim. These are the requests that make it a
fact:

- a customer listing subscriptions sees theirs and not the other customer's
- the same for invoices, which are the ones carrying money
- **a guessed id is still not yours**: reading another customer's row by id
  answers 404, not 403, because 403 confirms the row exists
- the operator sees both customers, so the test cannot pass on an API that
  returns nothing to anybody
- a signed-in stranger without the grants gets 403

Run them:

```bash
cd apps/api && go test ./internal/handlers/ -run 'TestThePortal|TestAGuessedID|TestTheOperator|TestASignedIn' -v
```

They are verified to fail when the protection is removed, not merely to pass.
Deleting the `authz.ScopeOwned` line from `subscriptions.List` turns the first
one red with `Ada can see Basil's subscription: the owner scope is not on the
list query`.

---

## What building it found

Examples earn their place by finding things, and this one found a framework
bug before it was finished.

**A date column could be filtered but not sorted.** `Sortable` was built by
testing a field's Go type against `string`, `int` and `uint`. A date's Go type
is `*jsontime.Date`, so `issued_on` was filterable and not sortable, and
`paginate` drops an unrecognised sort column rather than refusing it. So
`?sort=issued_on` returned 200 with the rows in whatever order the database
chose, and the invoices table looked like it had a sorting bug. An invoice list
that cannot be ordered by its issue date is the obvious case of it. Fixed in
v3.393.0; dates and floats are sortable now.

---

## Taking it further

Nothing here is Stripe. There is no payment provider, because the portal is
the part that is the same whichever one you pick: a webhook writes a row, and
these pages read it. Add the provider, keep the shape.

The two pages are deliberately plain. Copy them, rename the resource, and the
confinement comes along for free, because it never lived in the page.
