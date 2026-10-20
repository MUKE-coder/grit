#import "../template.typ": *

= Handlers and the HTTP contract

A handler's job is small and the contract it serves is not. This chapter is
about the envelope every endpoint returns, why it is the same shape everywhere,
and the two status codes that are worth understanding properly.

== A handler, in full

```go
func (h *ContactHandler) List(c *gin.Context) {
    opts := paginate.FromQuery(c, contactListConfig)
    result, err := h.service().List(h.ctx(c), opts)
    if err != nil {
        respond.Fail(c, err)
        return
    }
    respond.Paginated(c, result)
}
```

Four lines of work. Parse the query against the whitelist, call one service
method, hand any error to the boundary, write the answer.

The `h.ctx(c)` is chapter thirty two's point: it converts the Gin context into
a `context.Context` carrying the actor and the organisation, which is what
makes the service's queries scoped.

== One envelope, everywhere

Every successful response has the same shape.

A single item:

```json
{
  "data": { "id": 1, "name": "Ada" },
  "message": "Contact created successfully"
}
```

A list:

```json
{
  "data": [ ... ],
  "meta": { "total": 100, "page": 1, "page_size": 20, "pages": 5 }
}
```

And every failure:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": { "email": "This field is required" }
  }
}
```

#bridge[
  The value of a fixed envelope is on the *client* side. One `fetch` wrapper
  knows how to find the data, one error handler knows how to find the message,
  and `meta` is always where pagination lives. Without it, every endpoint is a
  small negotiation and the frontend accumulates special cases.

  This is why the generated React hooks are so short: they can assume the
  shape.
]

== Status codes, and the ones that matter

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [200], [OK],
  [201], [Created. A `POST` that made something returns this, with the thing],
  [400], [The request was malformed. Bad JSON, a missing parameter],
  [401], [Not authenticated. You have no identity],
  [403], [Authenticated, and not allowed. You have one, it is insufficient],
  [404], [No such thing, or none you may see],
  [409], [Conflict. Covered below],
  [422], [Well formed, and the values are wrong. The validation case],
  [429], [Rate limited],
  [500], [We broke],
)

Two distinctions people get wrong.

*400 against 422.* A 400 means the request could not be understood: invalid
JSON, a parameter that is not a number. A 422 means it was understood perfectly
and the values are unacceptable: an empty name, an email without an `@`. The
frontend treats them differently, because only one of them can point at a
field.

*401 against 403.* A 401 says "I do not know who you are", and the client
should log in. A 403 says "I know exactly who you are and the answer is no",
and logging in again will not help. Returning 401 for a permissions failure
sends users into a login loop.

== 404 as a permissions answer

A deliberate choice worth stating: when you ask for a record you are not
allowed to see, you get *404, not 403*.

A 403 confirms the record exists. That is an information leak: an attacker can
enumerate IDs and learn which invoices, users or documents are real, even
without reading them. Answering 404 for both "absent" and "not yours" gives
nothing away.

#gotcha[
  This surprises developers during integration, because a record they can see
  in the database returns 404 through the API. The cause is almost always
  ownership scoping doing its job: the query is filtered by the actor in the
  context, so the row is genuinely not in the result set.

  The fix is not to loosen the scope. It is to check who the request is
  authenticated as.
]

== The 409, and a number that must not be a string

Optimistic locking is the `Version` column from chapter thirty one. A client
sends the version it is editing, and the server refuses if it has moved:

```
If-Match: 8
```

```json
{
  "error": {
    "code": "VERSION_CONFLICT",
    "message": "this invoice changed while you were editing it",
    "current_version": 9
  }
}
```

That response is doing something more useful than reporting a failure: it is
telling the interface *enough to act*. With the current version the client can
offer to reload, show the difference, or retry. A bare "conflict" leaves the
frontend with nowhere to go but an alert box.

#gotcha[
  `"current_version": 9`, not `"9"`.

  It shipped as a quoted string once. Every client that compared versions was
  comparing text, so `"10" < "9"` was true and the comparison was silently
  wrong in exactly the cases that matter.

  The general lesson: a number that crosses a wire as a string will be compared
  as a string by somebody, and the failure is silent. This is the same family
  of bug as `omitempty` on a bool from chapter nineteen.
]

== respond, and the 119 codes

The `respond` package is the boundary from chapter eleven, made concrete.

```go
respond.OK(c, data)
respond.Created(c, data, "Contact created successfully")
respond.Paginated(c, result)

respond.NotFound(c, "contact")
respond.Forbidden(c, "not your contact")
respond.Conflict(c, err)
respond.Fail(c, err)              // works it out from the error
```

`respond.Fail` is the one to use by default. It walks the error chain with
`errors.Is` and `errors.As`, finds the sentinel or typed error, and chooses the
status and code. The handler does not decide; the error already knows.

There are 119 distinct error codes, which sounds excessive until you consider
what a frontend does with them. `EMAIL_EXISTS` can focus the email field.
`ACCOUNT_LOCKED` can show a different screen from `INVALID_CREDENTIALS`.
`CSRF_INVALID` can retry after refreshing the token. A single
`VALIDATION_ERROR` for all of them makes every one of those a string match on
the message.

#bridge[
  A `Coded` interface lets an error carry its own wire shape:

  ```go
  // Implement it on a sentinel when the caller needs to act on it.
  type Coded interface {
      Code() string
      Status() int
  }
  ```

  So a package can define an error that knows it is a 409 called
  `STOCK_INSUFFICIENT`, and `respond.Fail` needs no case for it. That is
  chapter ten's consumer-declared interface in the framework's own code.
]

== Validation, in two places that agree

A generated handler binds the request into a struct with tags:

```go
type CreateContact struct {
    Name  string `json:"name"  binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Phone string `json:"phone" binding:"omitempty,e164"`
}
```

A binding failure becomes a 422 with per-field details, which is where the
`details` object in the envelope comes from.

The matching Zod schema in `packages/shared` was generated from the same field
definitions, so the browser refuses the same values for the same reasons.
Neither is trusted: the frontend validation is for the user, and the backend
validation is the control.

== Writing your own handler

Adding an endpoint is ordinary Gin, and it should stay thin:

```go
func (h *InvoiceHandler) MarkPaid(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        respond.BadRequest(c, "invalid id")
        return
    }

    inv, err := h.service().MarkPaid(h.ctx(c), uint(id), time.Now())
    if err != nil {
        respond.Fail(c, err)      // ErrAlreadyPaid becomes a 409
        return
    }
    respond.OK(c, inv)
}
```

Then register it in the routes file, next to the generated ones.

#gotcha[
  Resist the urge to put the rule here because the data is already to hand. The
  check for "already paid" belongs in the service, as it was written in chapter
  thirty two, so that the CSV importer and the CLI obey it too.

  The test for whether a line belongs in a handler: *would it still need to be
  true if this were called by a background job?* If yes, it belongs one layer
  down.
]

#project(title: "Project 30: exercise the whole contract")[
  Use the bookshop and drive every branch of the envelope from the command
  line, so the contract stops being documentation.

  ```bash
  grit start
  ```

  *The happy paths.*

  ```bash
  curl -s localhost:8080/api/v1/books | jq '.meta'
  curl -s -X POST localhost:8080/api/v1/books \
    -H 'Content-Type: application/json' \
    -d '{"title":"Go","price":2000}' | jq
  ```

  Note the 201 and the `message`, and that `meta` on the list has four keys.

  *400 against 422.* Send broken JSON, then valid JSON with a bad value:

  ```bash
  curl -s -X POST localhost:8080/api/v1/books \
    -d '{"title":' | jq '.error.code'
  curl -s -X POST localhost:8080/api/v1/books \
    -H 'Content-Type: application/json' -d '{"title":""}' | jq '.error'
  ```

  The second should carry `details` naming the field. That object is what lets
  a form highlight the right input rather than showing a banner.

  *401 against 403.* Call a protected route with no token, then with a token
  for a user who lacks the role. Two different codes, and the frontend should
  do two different things.

  *404 for someone else's record.* Create a book as one user, then fetch it as
  another. You get 404, not 403, and now you know why.

  *The 409.* Fetch a book and note its `version`. Update it twice, the second
  time sending the stale version:

  ```bash
  curl -s -X PATCH localhost:8080/api/v1/books/1 \
    -H 'If-Match: 1' -H 'Content-Type: application/json' \
    -d '{"title":"First"}' | jq '.data.version'

  curl -s -X PATCH localhost:8080/api/v1/books/1 \
    -H 'If-Match: 1' -H 'Content-Type: application/json' \
    -d '{"title":"Second"}' | jq
  ```

  The second returns 409 with `current_version`. Check the type in the raw
  JSON: it must be `9`, not `"9"`.

  *The whitelist.* Try `?sort=title`, then `?sort=secret_column`. One works and
  one falls back silently, which is the behaviour chapter thirty two warned
  about.
]

#tryit[
  Write a handler test with `httptest` that asserts the 409 body contains a
  numeric `current_version`, using a JSON decode into a struct with
  `CurrentVersion int`. If the server ever regresses to a string, the decode
  fails and the test catches it. That is a four line test guarding a bug that
  shipped once and was invisible.
]
