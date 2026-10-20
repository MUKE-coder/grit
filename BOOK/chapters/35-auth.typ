#import "../template.typ": *

= Authentication, roles and per-record ownership

B1 was "can different people see different things", and it is the question
clients ask first. The answer has three layers that people routinely confuse,
so this chapter separates them before showing what is generated.

== Three questions, not one

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*Authentication*], [Who are you? Produces an identity, or a 401],
  [*Authorization*], [May people like you do this? Roles and permissions, or a 403],
  [*Ownership*], [May *you* do it to *this record*? Per-row, or a 404],
)

Most teams build the first, assume the second is the whole of the problem, and
discover the third when a salesperson can open another salesperson's accounts
by changing a number in the URL.

#bridge[
  That third one has a name: *insecure direct object reference*, and it is
  perennially near the top of the OWASP list. It is not an exotic attack. It is
  changing `/invoices/41` to `/invoices/42`, and an application that checks
  only the role will serve it happily, because the role was right.
]

== What arrives

A new project already has: registration, login, logout, password reset, email
verification, refresh tokens, server-side sessions, account lockout, optional
two-factor, and OAuth and SAML behind flags. Three roles exist by default,
`ADMIN`, `EDITOR` and `USER`, and you can add more.

None of this is a plugin you enable. B2's answer is present before you write a
line, for the reason chapter twenty seven gave: an opt-in security feature is a
security feature most projects do not have.

== Sessions, not just tokens

A JWT alone cannot be revoked: it is valid until it expires, so "log this
device out" is not expressible. Every refresh token here is backed by a
`sessions` row, stored as a SHA-256 hash rather than the token itself.

That buys four things a bare JWT cannot:

- revoke one device, from the account page
- detect replay: a refresh token used twice means it was stolen, and the whole
  session chain is killed
- idle and absolute timeouts
- a password change signs out every device

#gotcha[
  Storing only the hash matters. A database read, a leaked backup or an
  over-broad log should not hand someone a usable token. The same reasoning as
  password hashing, applied to the thing that is just as good as a password.

  The refresh cookie is also scoped to `/api/auth`, so it is not sent to every
  endpoint. The consequence to know: a handler outside that path cannot see the
  refresh cookie, so it cannot identify the *session*, only the user.
]

== Four layers on the sign-in form

The login endpoint is the most attacked route in any application, and it
carries four defences:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [1], [A rate limit on the route],
  [2], [A lockout on the account after repeated failures],
  [3], [A growing delay on repeated failure],
  [4], [A locked account answers exactly like a wrong password],
)

The fourth is the one that gets missed, and it is the interesting one.

If a locked account says "this account is locked" and an unknown one says "no
such user", you have built an oracle for discovering who has an account.
Combined with a breach list, that is a targeted attack. The polite, helpful
error message *is* the vulnerability.

#gotcha[
  The same trap hides in the timing. If the locked path returns instantly while
  the real check spends eighty milliseconds hashing a password, an attacker
  reads the clock instead of the message and learns the same thing.

  So the delay has to be on both paths. This is the kind of detail that is
  obvious once stated and absent from almost every hand-rolled login, and it is
  a decent argument for not hand-rolling one.
]

== Roles, and adding your own

```bash
grit add role MODERATOR
```

That updates seven places: the Go constants, the Zod schema, the TypeScript
union, the sidebar visibility, the form options and the seeded permissions.
Doing it by hand means finding all seven, and the one you miss is usually the
frontend, where the role exists but nothing shows for it.

Protecting a route:

```go
admin := api.Group("/admin")
admin.Use(middleware.RequireRole("ADMIN", "EDITOR"))
```

And in a service, where it also holds for jobs and imports:

```go
if !authz.Can(ctx, "invoices.delete") {
    return respond.ErrForbidden
}
```

== Ownership is the one that bites

Roles say "editors may delete invoices". Ownership says "this editor may delete
*this* invoice".

```bash
grit generate resource Invoice --owned-by user \
  --fields "number:string, total:money"
```

That adds a `UserID` column and, more importantly, scopes every generated query
by the actor in the context. Listing returns only yours. Fetching someone
else's returns 404, not 403, for the reason chapter thirty three gave.

The scoping lives in the service rather than the handler, so it applies to the
CSV import, the export, the background job and the CLI, which is the whole
argument of chapter thirty two meeting the problem it was written for.

#gotcha[
  This is why threading the request's context matters so much. A
  `context.Background()` created halfway through a request has no actor, so the
  scope cannot be applied, and the query returns *everybody's* rows.

  The failure is silent and it is a data breach rather than an error. When a
  background job genuinely must see everything, there is an explicit way to say
  so, and the explicitness is the point: an unscoped query should be a sentence
  someone wrote deliberately.
]

== Multitenancy is the same idea one level up

B9 asked "can two of our clients use one system". Ownership scopes to a user;
multitenancy scopes to an organisation, through the same context.

The property to hold onto: *isolation is not a feature you add, it is something
every query must have*. One forgotten `WHERE` in one place, in a background job
or an export rather than a route, and one client sees another's data.

That is why it travels in the context rather than being a parameter each method
remembers to apply. A parameter can be forgotten at one call site. A context
that carries the organisation is applied by the layer that builds the query.

== What is not solved for you

Two honest limits.

*Permission design is yours.* The framework gives roles, permissions and
per-record ownership. Deciding that a regional manager sees their region, and
what happens when someone moves region, is domain modelling and nothing
generates it.

*Approval workflows are not authorization.* "An invoice over ten thousand needs
a second approver" is a business rule that lives in a service, possibly a saga.
It looks like permissions and it is not, and modelling it as a role produces a
mess.

#project(title: "Project 32: prove the isolation holds")[
  The only way to trust ownership scoping is to try to break it.

  ```bash
  grit new tenancy-lab --api
  cd tenancy-lab
  grit generate resource Note --owned-by user \
    --fields "title:string, body:text"
  grit migrate
  grit start
  ```

  *One: two users.* Register Ada and Grace through the API and keep both
  tokens.

  *Two: a note each.* Create one note as each user. Note the IDs.

  *Three: the attack.* As Grace, fetch Ada's note by ID.

  ```bash
  curl -s -H "Authorization: Bearer $GRACE" \
    localhost:8080/api/v1/notes/1 | jq
  ```

  You should get 404. Not 403, and not the note. Check `/studio` to confirm the
  row really is there, so you know the 404 is scoping rather than absence.

  *Four: the list.* `GET /notes` as each user returns one note each, and
  `meta.total` says 1 rather than 2. A count that leaks the true total is a
  subtler version of the same bug, so it is worth checking.

  *Five: the back door.* Try the same through every other surface. Export the
  notes as Grace and see whose appear. Run the CSV import as Grace with a row
  referencing Ada's note. If the scope only lived in the handler, one of these
  would leak; it does not, because it lives in the service.

  *Six: break it deliberately.* In `services/note.go`, change one method from
  `s.db(ctx)` to `s.DB` directly. Rebuild, and repeat step three. You will now
  get Ada's note.

  That one character is the difference between an application and an incident,
  and having caused it on purpose once is the best way to recognise it in
  review.
]

#tryit[
  Write the step-three test as a real Go test, with two users and an assertion
  that the second gets `ErrNotFound`. Then put it in CI.

  An isolation test is the single highest-value test in a multi-user
  application, because the bug it catches is silent, is not a crash, and is the
  one you cannot apologise your way out of. Grit's own repository keeps one for
  exactly this reason.
]

#assignment(name: "Auth - Assignment One")[
  Try to break your own isolation, then stop yourself from ever regressing it.

  *One.* Generate a resource with `--owned-by user`, and create one record as
  each of two users.

  *Two.* As the second user, request the first user's record by ID. You want a
  404. Confirm in `/studio` that the row is there.

  *Three.* Try the same through every other door: the list total, the export,
  the import, a background job. One of them is where a hand-written
  application would have leaked.

  *Four.* Break it on purpose, in one character, and watch the data appear.
  Then put it back.

  *Five.* Write the isolation test as a Go test and add it to CI. In a
  multi-user application it is the highest-value test you will ever write,
  because the bug it catches is silent.]
