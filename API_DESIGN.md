# API design

The HTTP contract every Grit API honours, and the contract every generated
resource and hand-written handler has to keep.

This is a **contract, not a style guide.** Three clients are built against it in
this repository alone (the admin panel, the web app, the Expo app), plus the
desktop client, the generated SDK and whatever an operator writes next. A
handler that answers in a different shape does not fail a test; it fails one
screen, usually the one nobody opened before release.

Where this document and the code disagree, the code named under **Source of
truth** wins, and this document is wrong and should be fixed.

---

## Source of truth

| Thing | Lives in |
|---|---|
| The envelope, and every helper that writes it | `internal/scaffold/api_files.go`, the `respond` package it emits |
| Every error code, and the status each one carries | `internal/errorcodes/catalog.go`, which generates `internal/respond/codes.go` |
| Which routes exist and what each demands of a caller | `internal/scaffold/api_access.go`, generated into `internal/access/access.go` |
| The version prefix | `APIVersion` in `internal/scaffold/api_files.go` |

**There is one table of code-to-status, and it is the catalogue.** There were
two for a while: the named helpers carried a status each, and the generated
`codes.go` carried the same pairs. That is how `VALIDATION_ERROR` came back as
422 from one helper and 400 from a handler that wrote its own envelope. Do not
reintroduce a second table, and do not write `c.JSON` with an error body by
hand.

---

## The envelope

Every response is an object. Never a bare array, never a bare string: an array
at the top level cannot gain a `meta` later without breaking every caller.

### One record

```json
{
  "data": { "id": "...", "name": "..." },
  "message": "User created successfully"
}
```

`message` is optional and is for a person. A client must never branch on it.

### A list

```json
{
  "data": [ ... ],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}
```

`total`, `page`, `page_size` and `pages` are **always present**, including on an
empty result. Zero is an answer: a response that omits `total` leaves every
client doing `meta.total` with `undefined`, which renders as a blank stat card
rather than a nought and turns any arithmetic into `NaN`.

A keyset (cursor) response adds `next_cursor`, `has_more` and `mode: "cursor"`,
and omits them otherwise.

### An error

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": { "email": "This field is required" }
  }
}
```

- `code` is a stable, screaming-snake-case string from the catalogue. Clients
  branch on this.
- `message` is one sentence for a person. It may change without notice.
- `details` maps a field name to its own sentence, so a form can put the message
  under the input it belongs to. Omitted when there is nothing per-field.

---

## Status codes

| Code | When |
|---|---|
| 200 | Read, update, delete, and anything else that worked |
| 201 | A record was created |
| 204 | Nothing to return, and the client knows it (rare: prefer 200 with `data`) |
| 400 | The request is malformed. The client cannot fix it without sending something different |
| 401 | No credentials, or they are not valid |
| 403 | Valid credentials, insufficient permission |
| 404 | No such record, or none this caller may see |
| 409 | A conflict with a record that exists: a duplicate, a version mismatch |
| 422 | The request is well-formed and the values are wrong, or a business rule says no |
| 429 | Rate limited |
| 500 | We broke |

**400 against 422** is the distinction that gets muddled. 400 is "I cannot
parse this" or "a required parameter is absent". 422 is "I understood you and
the answer is no": a failed validation, a password that breaks the rules, debits
that do not balance.

### 404 over 403 for records

A record the caller may not see answers **404**, not 403. A 403 confirms the
record exists, which is an information leak on any endpoint keyed by a guessable
id. Reserve 403 for an *action* the caller is not permitted, on a record they
can already see.

### 500 never explains itself

`ServerError` logs the error with the method, path and request id, and sends a
generic message. The text of a driver error can describe the schema and
sometimes contains SQL. The caller gets a code and a sentence; the operator gets
the cause and the `X-Request-ID` that ties a report to its log line.

---

## Writing a response

Use the helpers. They are each one line over `Fail`, and they read the status
out of the catalogue.

```go
respond.OK(c, user)                        // 200 { data }
respond.OK(c, user, "Saved")               // 200 { data, message }
respond.Created(c, user, "User created")   // 201

respond.BadRequest(c, "page must be a number")        // 400
respond.Unauthorized(c, "")                           // 401
respond.Forbidden(c, "You cannot publish")            // 403
respond.NotFound(c, "User not found")                 // 404
respond.Conflict(c, "That email is already taken")    // 409
respond.Validation(c, "Check the form", fields)       // 422
respond.Internal(c, err)                              // 500, logged
```

A paginated list has **no helper yet**, and is written out:

```go
c.JSON(http.StatusOK, gin.H{"data": out, "meta": res.Meta})
```

That is the one place the envelope is spelled by hand, and it is the most
common response in the product, which is the wrong way round. A `respond.List`
would close it. Until then: `data` and `meta`, those names, that order, and
`respond.Meta` for the struct so the counts cannot be omitted.

From inside a model hook or a service, where there is no `*gin.Context`, return
a rule error and let the handler translate it:

```go
func (e *JournalEntry) BeforeCreate(tx *gorm.DB) error {
    if !balanced(e.Lines) {
        return respond.Rule("debits and credits do not balance")
    }
    return nil
}
```

`respond.WriteError` turns that into 422 with that sentence, and anything it
does not recognise into a logged 500. It also recognises a driver's duplicate
key and answers 409 naming the column, so a form can mark the field.

### Adding an error code

Add it to `internal/errorcodes/catalog.go` with its status, and regenerate.
Never invent a code at the call site: an uncatalogued code is one no client can
branch on and no document lists.

The one exception is `ServerError`, which takes a plain string, so a 500 can
name the subsystem that broke (`DB_ERROR`, `STORAGE_ERROR`) without every one of
them being a catalogued constant. The status is 500 either way.

---

## Routes

```
/api/v1/<plural-resource>      generated resources
/api/auth/...                  sign in, register, refresh, reset, OAuth, TOTP
/api/profile                   the caller's own account
/api/admin/...                 operator-only endpoints
/api/health                    liveness and the probe registry
/api/uploads/...               presigned uploads
/api/public/...                deliberately unauthenticated
```

- **Plural, lowercase, hyphenated.** `/api/v1/blog-posts`, never
  `/api/v1/blogPost` or `/api/v1/blog_posts`.
- **Generated resources are versioned**, under `/api/v1`. The framework's own
  endpoints are not, and that asymmetry is deliberate: a project's own resources
  are the surface its clients pin, and `auth` and `health` are ours to move.
- **The version is pinned by the client, once.** `lib/api-core.ts` rewrites
  `/api` to `/api/v1` in an interceptor, so no call site in any frontend spells
  the version. Do not spell it at a call site.
- **A verb in a path is a smell, not a sin.** `POST /api/v1/orders/:id/refund`
  is better than modelling a refund as a `PATCH` with a magic field. Prefer a
  noun where one exists.

### Query parameters for a list

| Parameter | Means |
|---|---|
| `page` | 1-based page number |
| `page_size` | Rows per page. Clamped server-side; a client cannot ask for everything |
| `sort` | A column, optionally `-` prefixed for descending |
| `search` | Free text, matched against the resource's searchable columns |
| `<field>` | An exact filter on a filterable column |

Snake_case, matching the JSON field names, matching the column names.

---

## Authentication

- The API issues **HttpOnly** `grit_access` and `grit_refresh` cookies on sign
  in, register, refresh and the OAuth callback. JavaScript never reads or writes
  the access token.
- `grit_refresh` is **scoped to `/api/auth`**. A handler outside that prefix
  cannot see it, and therefore cannot identify the current session. If an
  endpoint needs to, it belongs under `/api/auth` or it needs another signal.
- Every refresh token is backed by a `sessions` row, stored as SHA-256 only, and
  rotated on use with replay detection.
- The CSRF token rides a **non-HttpOnly** `grit_csrf` cookie, echoed into
  `X-CSRF-Token` on every state-changing method.
- `Authorization: Bearer <token>` is accepted as well, for the Expo app, the
  desktop client and any API client. Cookies are for browsers.

---

## Headers

| Header | Direction | Means |
|---|---|---|
| `X-Request-ID` | response | Ties a report of a failure to its log line. Echoed from the request when sent |
| `X-CSRF-Token` | request | Required on POST, PUT, PATCH, DELETE from a browser |
| `Idempotency-Key` | request | A client-generated key. A replay returns the first response rather than acting twice |
| `X-Forwarded-For`, `X-Forwarded-Proto`, `X-Real-IP` | request | Honoured behind a proxy only, from the configured trusted set |

---

## The rules that are easy to break

1. **Never write an error body by hand.** `c.JSON(400, gin.H{"error": ...})`
   compiles, ships, and puts a second shape on the wire.
2. **Handlers stay thin.** A handler parses, calls a service, and responds. A
   query in a handler is a query no other caller can reuse and no test can
   cover, and twenty of them were pulled out of handlers for exactly that
   reason.
3. **Every error is handled.** No `_ = err` on a path that can fail. A dropped
   error is a 200 describing a write that did not happen.
4. **A public endpoint that writes is not cacheable.** `--public` is read-only;
   a writing public endpoint is hand-written and must not be mounted on the
   URL-cached `publicAPI` group.
5. **Pagination is not optional.** A list endpoint with no `page_size` clamp is
   a denial of service with extra steps.
6. **The response shape is part of the change.** If a handler's output changes,
   the Zod schema in `packages/shared` and the TypeScript type change with it,
   in the same commit. `grit sync` exists for this.
