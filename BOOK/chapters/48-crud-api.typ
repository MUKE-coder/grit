#import "../template.typ": *

= Build your own CRUD API from scratch

Every backend job, in the end, asks for the same thing: create a record, read
one, list many with a filter, change one, delete one. Parts I to III gave you
every piece. This chapter assembles them into a complete API with no framework
and no ORM, and then counts what it cost, because that number is the honest
case both for and against the rest of this book.

Nothing here imports anything outside the standard library except a database
driver.

== What we are building

A tasks API, in one binary:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`POST /tasks`], [Create, with validation, returning 201],
  [`GET /tasks`], [List, paginated, filtered, sorted, with `meta`],
  [`GET /tasks/{id}`], [Read one, or 404],
  [`PATCH /tasks/{id}`], [Partial update],
  [`DELETE /tasks/{id}`], [Delete, returning 204],
)

Plus the things that separate an API from a demo: one response envelope, real
status codes, a request ID in the logs, timeouts on the server, a recover
boundary, graceful shutdown and tests that run without a port.

The envelope is the one from chapter thirty three, on purpose. You are about to
write by hand the contract a generated handler gives you.

== Three layers, decided before any code

```
main.go        wiring: config, database, routes, server
http/          router, handlers, middleware, the envelope
task/          the type, its rules, and the store
```

The rule, which is the whole of chapter thirty two said again in your own
project: *a handler parses and responds, and decides nothing*. The store talks
to the database and knows no HTTP. In between sits the thing that knows the
rules.

#bridge[
  In an Express application this layering is a convention you keep by
  discipline, and the reason it slips is that `req` and `res` are in scope
  everywhere, so any file can reach for them.

  Here, the store package cannot see an HTTP request even if it wants to: it
  does not import `net/http`. The architecture is held by the import graph
  rather than by review, which is why it survives a busy week.
]

== The type and its rules

```go
package task

type Task struct {
    ID        int64      `json:"id"`
    Title     string     `json:"title"`
    Notes     string     `json:"notes,omitempty"`
    Done      bool       `json:"done"`
    DueAt     *time.Time `json:"due_at"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

`DueAt` is a pointer because a task may genuinely have no due date, and
chapter nineteen's lesson applies: a zero `time.Time` serialises as the year
one, which is not "none".

Validation lives with the type, returning field errors rather than one string,
because the frontend needs to know which input to highlight:

```go
type FieldErrors map[string]string

func (e FieldErrors) Error() string { return "validation failed" }

func (t *Task) Validate() error {
    errs := FieldErrors{}
    if strings.TrimSpace(t.Title) == "" {
        errs["title"] = "this field is required"
    }
    if len(t.Title) > 200 {
        errs["title"] = "must be 200 characters or fewer"
    }
    if t.DueAt != nil && t.DueAt.Before(time.Now()) {
        errs["due_at"] = "must be in the future"
    }
    if len(errs) > 0 {
        return errs
    }
    return nil
}
```

== The store

The schema, applied at startup. For a project this size that is enough, and
chapter thirty one explains what you give up by not having real migrations:

```go
const schema = `
CREATE TABLE IF NOT EXISTS tasks (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  title      TEXT    NOT NULL,
  notes      TEXT    NOT NULL DEFAULT '',
  done       BOOLEAN NOT NULL DEFAULT 0,
  due_at     DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS tasks_done ON tasks(done);`
```

Create, returning the assigned ID:

```go
func (s *Store) Create(ctx context.Context, t *Task) error {
    now := time.Now().UTC()
    t.CreatedAt, t.UpdatedAt = now, now

    res, err := s.db.ExecContext(ctx, `
        INSERT INTO tasks (title, notes, done, due_at,
                           created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)`,
        t.Title, t.Notes, t.Done, t.DueAt, t.CreatedAt, t.UpdatedAt)
    if err != nil {
        return fmt.Errorf("insert task: %w", err)
    }

    t.ID, err = res.LastInsertId()
    if err != nil {
        return fmt.Errorf("task id: %w", err)
    }
    return nil
}
```

Read one, with the sentinel the rest of the program will match on:

```go
var ErrNotFound = errors.New("task not found")

func (s *Store) Get(ctx context.Context, id int64) (*Task, error) {
    var t Task
    err := s.db.QueryRowContext(ctx, `
        SELECT id, title, notes, done, due_at, created_at, updated_at
        FROM tasks WHERE id = ?`, id).
        Scan(&t.ID, &t.Title, &t.Notes, &t.Done, &t.DueAt,
            &t.CreatedAt, &t.UpdatedAt)

    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("get task %d: %w", id, err)
    }
    return &t, nil
}
```

#gotcha[
  `sql.ErrNoRows` must not escape this package. It is a `database/sql` detail,
  and if a handler matches on it then swapping the store for an in-memory one
  breaks the handler, which is exactly the coupling the layering was for.

  Translate it once, at the edge of the package that produced it. That is
  chapter eleven's boundary rule, and this is the line where it is easiest to
  get right and easiest to skip.
]

== Listing, which is where the work actually is

Create and read are five lines each. Listing is where every API grows hair:
pagination, a filter, a search, a sort, and a total count that has to agree
with the filter.

```go
type ListOptions struct {
    Page, PageSize int
    Search         string
    Done           *bool
    Sort           string
    Desc           bool
}
```

Build the `WHERE` once, and use it for both queries:

```go
func (s *Store) List(ctx context.Context, o ListOptions) (
    []Task, int, error) {

    where := []string{"1 = 1"}
    args := []any{}

    if o.Search != "" {
        where = append(where, "title LIKE ?")
        args = append(args, "%"+o.Search+"%")
    }
    if o.Done != nil {
        where = append(where, "done = ?")
        args = append(args, *o.Done)
    }
    clause := strings.Join(where, " AND ")
```

The count comes from the same clause, or your pagination lies:

```go
    var total int
    err := s.db.QueryRowContext(ctx,
        "SELECT COUNT(*) FROM tasks WHERE "+clause, args...).
        Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("count tasks: %w", err)
    }
```

And then the page itself. Note what is concatenated and what is a placeholder:

```go
    sortCols := map[string]string{
        "title":      "title",
        "due_at":     "due_at",
        "created_at": "created_at",
    }
    col, ok := sortCols[o.Sort]
    if !ok {
        col = "created_at"
    }
    dir := "ASC"
    if o.Desc {
        dir = "DESC"
    }

    q := fmt.Sprintf(`
        SELECT id, title, notes, done, due_at, created_at, updated_at
        FROM tasks WHERE %s ORDER BY %s %s LIMIT ? OFFSET ?`,
        clause, col, dir)

    args = append(args, o.PageSize, (o.Page-1)*o.PageSize)
    rows, err := s.db.QueryContext(ctx, q, args...)
```

#gotcha[
  A column name cannot be a placeholder. `ORDER BY ?` is not valid SQL, so the
  sort column has to be concatenated into the string, and that is precisely how
  a sort parameter becomes an injection.

  The map is the defence, and the important part is that it is a *whitelist of
  allowed values*, not an escape or a check for suspicious characters. Anything
  not in the map is not a column, so the worst an attacker achieves is the
  default sort.

  This is the same whitelist a generated Grit service carries, and now you know
  why it exists rather than being a configuration chore.
]

Finish the rows properly, because the loop is not the end:

```go
    defer rows.Close()

    tasks := []Task{}
    for rows.Next() {
        var t Task
        if err := rows.Scan(&t.ID, &t.Title, &t.Notes, &t.Done,
            &t.DueAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
            return nil, 0, fmt.Errorf("scan task: %w", err)
        }
        tasks = append(tasks, t)
    }
    return tasks, total, rows.Err()
}
```

`rows.Err()` is the line people omit for years without noticing. A connection
dropped halfway through a result set ends the loop quietly, and without that
check you return half a page and call it success.

== A router, written rather than imported

Go 1.22 taught `ServeMux` method and wildcard patterns, so in a real project
you would use it. Write one once anyway: a router is about forty lines, and
knowing that number changes how you think about frameworks.

```go
type route struct {
    method  string
    parts   []string       // "tasks", ":id"
    handler http.HandlerFunc
}

type Router struct{ routes []route }

func (r *Router) Handle(method, pattern string,
    h http.HandlerFunc) {
    r.routes = append(r.routes, route{
        method:  method,
        parts:   strings.Split(strings.Trim(pattern, "/"), "/"),
        handler: h,
    })
}
```

Matching is a comparison, segment by segment, with the parameters collected
into the request context as they are found:

```go
func (r *Router) ServeHTTP(w http.ResponseWriter,
    req *http.Request) {

    got := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
    allowed := []string{}

    for _, rt := range r.routes {
        params, ok := match(rt.parts, got)
        if !ok {
            continue
        }
        if rt.method != req.Method {
            allowed = append(allowed, rt.method)
            continue
        }
        ctx := context.WithValue(req.Context(), paramsKey, params)
        rt.handler(w, req.WithContext(ctx))
        return
    }

    if len(allowed) > 0 {
        w.Header().Set("Allow", strings.Join(allowed, ", "))
        Fail(w, req, StatusError{405, "METHOD_NOT_ALLOWED",
            "method not allowed"})
        return
    }
    Fail(w, req, StatusError{404, "NOT_FOUND", "no such route"})
}
```

```go
func match(pattern, got []string) (map[string]string, bool) {
    if len(pattern) != len(got) {
        return nil, false
    }
    params := map[string]string{}
    for i, p := range pattern {
        if strings.HasPrefix(p, ":") {
            params[p[1:]] = got[i]
            continue
        }
        if p != got[i] {
            return nil, false
        }
    }
    return params, true
}
```

The 405 is the detail worth copying. A path that exists for another method is
not a 404, and an API that cannot tell the difference sends every client on a
debugging trip through its own URL building.

== One envelope, one place

```go
func JSON(w http.ResponseWriter, status int, body any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(body); err != nil {
        slog.Error("write response", "err", err)
    }
}

func OK(w http.ResponseWriter, data any) {
    JSON(w, 200, map[string]any{"data": data})
}

func Created(w http.ResponseWriter, data any, msg string) {
    JSON(w, 201, map[string]any{"data": data, "message": msg})
}

func Paginated(w http.ResponseWriter, data any, m Meta) {
    JSON(w, 200, map[string]any{"data": data, "meta": m})
}
```

And one failure path that works out the status from the error, so no handler
decides it twice:

```go
func Fail(w http.ResponseWriter, r *http.Request, err error) {
    var fe task.FieldErrors
    var se StatusError

    switch {
    case errors.As(err, &fe):
        write(w, 422, "VALIDATION_ERROR", "validation failed", fe)
    case errors.Is(err, task.ErrNotFound):
        write(w, 404, "NOT_FOUND", "task not found", nil)
    case errors.As(err, &se):
        write(w, se.Status, se.Code, se.Message, nil)
    default:
        slog.Error("unhandled", "err", err,
            "request_id", RequestID(r.Context()))
        write(w, 500, "INTERNAL", "something went wrong", nil)
    }
}
```

#gotcha[
  The `default` branch logs the real error and returns a generic one. That
  split is deliberate: the message a stranger sees must not contain the table
  name, the query or the file path, and the message you see must contain all
  three.

  The request ID is what joins the two halves. Return it in a header, log it
  with every line, and a support message saying "it failed at 14:05" becomes a
  one-command search.
]

== Handlers, which are now boring

```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var t task.Task
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        Fail(w, r, StatusError{400, "BAD_JSON", "invalid JSON"})
        return
    }
    if err := t.Validate(); err != nil {
        Fail(w, r, err)
        return
    }
    if err := h.store.Create(r.Context(), &t); err != nil {
        Fail(w, r, err)
        return
    }
    Created(w, t, "Task created successfully")
}
```

Four steps, in the same order, in every handler in the file: parse, validate,
call, respond. When a fifth thing appears in a handler it is nearly always a
rule, and it belongs one layer down.

The list handler is the one that parses the query, and it clamps rather than
trusts:

```go
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    o := task.ListOptions{
        Page:     atoiDefault(q.Get("page"), 1),
        PageSize: atoiDefault(q.Get("page_size"), 20),
        Search:   q.Get("search"),
        Sort:     q.Get("sort"),
        Desc:     q.Get("order") == "desc",
    }
    if o.Page < 1 {
        o.Page = 1
    }
    if o.PageSize < 1 || o.PageSize > 100 {
        o.PageSize = 20
    }
    ...
}
```

#gotcha[
  The clamp on `page_size` is not tidiness. Without it, `?page_size=1000000` is
  a one-request denial of service that your own API happily performs, and it
  will be found by a crawler rather than by an attacker.
]

== Middleware, as a chain

```go
func Chain(h http.Handler, mw ...func(http.Handler) http.Handler) (
    http.Handler) {
    for i := len(mw) - 1; i >= 0; i-- {
        h = mw[i](h)
    }
    return h
}
```

The recover boundary is the one that earns its place on the first bad day:

```go
func Recover(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter,
        r *http.Request) {
        defer func() {
            if v := recover(); v != nil {
                slog.Error("panic", "value", v,
                    "stack", string(debug.Stack()))
                write(w, 500, "INTERNAL", "something went wrong",
                    nil)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

Chapter twelve said `recover` belongs only at a boundary. This is the boundary:
one panic takes down one request rather than the process, and the stack reaches
the log instead of the client.

== Wiring it up

```go
func main() {
    cfg := config.Load()

    db, err := sql.Open("sqlite", cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("open database: %v", err)
    }
    defer db.Close()
    db.SetMaxOpenConns(cfg.MaxConns)

    if _, err := db.Exec(schema); err != nil {
        log.Fatalf("migrate: %v", err)
    }

    r := &Router{}
    h := &Handler{store: task.NewStore(db)}
    r.Handle("POST", "/tasks", h.Create)
    r.Handle("GET", "/tasks", h.List)
    r.Handle("GET", "/tasks/:id", h.Get)
    r.Handle("PATCH", "/tasks/:id", h.Update)
    r.Handle("DELETE", "/tasks/:id", h.Delete)

    srv := &http.Server{
        Addr:         cfg.Addr,
        Handler:      Chain(r, Recover, RequestIDMW, LogMW),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    serveUntilSignal(srv)
}
```

`serveUntilSignal` is chapter twenty one's graceful shutdown, unchanged.

== Testing the whole thing without a port

```go
func newTestAPI(t *testing.T) http.Handler {
    t.Helper()
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("open: %v", err)
    }
    t.Cleanup(func() { db.Close() })
    if _, err := db.Exec(schema); err != nil {
        t.Fatalf("schema: %v", err)
    }
    return routes(&Handler{store: task.NewStore(db)})
}

func TestCreateRejectsEmptyTitle(t *testing.T) {
    api := newTestAPI(t)
    req := httptest.NewRequest("POST", "/tasks",
        strings.NewReader(`{"title":""}`))
    rec := httptest.NewRecorder()

    api.ServeHTTP(rec, req)

    if rec.Code != 422 {
        t.Fatalf("status = %d, want 422", rec.Code)
    }
    if !strings.Contains(rec.Body.String(), "title") {
        t.Errorf("body should name the field, got %s", rec.Body)
    }
}
```

An in-memory database per test, no fixtures to clean up, no ports to collide,
and the tests run in parallel. This is the shape every API test in your career
should have.

== Now count it

When it works, count the lines:

```bash
find . -name '*.go' -not -name '*_test.go' | xargs wc -l | tail -1
```

For one resource, done properly, it lands somewhere around six hundred lines,
plus tests. Most of it is not the task at all: it is the envelope, the router,
the middleware, the pagination and the whitelist, and all of that is written
once and then reused.

Which is the interesting part. *The second resource costs you maybe a hundred
and fifty lines*, because the machinery is already there. The tenth costs the
same hundred and fifty. That is D1 from chapter twenty six, measured in your
own editor rather than taken on trust.

#bridge[
  Now put this chapter next to chapter twenty nine. `grit generate resource
  Task` writes the equivalent of these six hundred lines, plus the TypeScript
  types, the Zod schema, the React Query hooks and an admin screen, from one
  description.

  The reason to have written this chapter first is not to make that look
  impressive. It is so that when the generated service does something you did
  not expect, you are reading familiar code rather than magic, and so you can
  tell, specifically, which parts of it you could do without.
]

#project(title: "Project 40: take it past CRUD")[
  Build the API above, then add the three things that make it a real service.

  *One: a transaction.* Add `POST /tasks/{id}/complete` that marks the task
  done and writes a row to a `task_events` table. Both or neither, in one
  transaction, with a test that proves a failure on the second write leaves no
  event and an unfinished task.

  *Two: concurrency.* Add a `version` column. `PATCH` requires an `If-Match`
  header, and a stale version returns 409 with the current one as a *number*.
  Chapter thirty three explains what the quoted string cost.

  *Three: a bulk endpoint.* `POST /tasks/bulk` takes an array and reports
  per-row results rather than failing the whole request on row seventeen.
  Decide deliberately whether a partial success is a 200 or a 207, and write
  the reason in a comment.

  *Four: measure it.* Put a thousand rows in and time the list endpoint with
  and without the index on `done`. Keep both numbers.
]

#tryit[
  Delete the sort whitelist and point the API at a throwaway database. Then
  send `?sort=title; DROP TABLE tasks --` and watch what happens to your data.

  Do it once, on purpose, in a directory you do not mind losing. An injection
  you have performed yourself stops being a checklist item and becomes
  something you can see in a diff.
]

#assignment(name: "CRUD API - Assignment One")[
  Build the whole API in `part6-from-scratch/48-crud-api`, with the three
  additions from Project 40 and the tests passing.

  *One.* All five endpoints, the envelope, the router, the middleware chain and
  graceful shutdown. No third-party imports except the database driver.

  *Two.* Tests covering: a 201 with the body, a 422 naming the field, a 404,
  a 405 with an `Allow` header, and a list whose `meta.total` respects the
  filter.

  *Three.* `wc -l` on your non-test Go files, recorded in `NOTES.md` with the
  time it took you. Then add a second resource and record both again. The ratio
  between the two is the number this book has been arguing about since chapter
  twenty six.

  *Four.* One paragraph: after writing it, which parts would you want a
  framework to hand you, and which would you rather keep writing yourself?
  Keep the answer. Chapter twenty nine is more interesting if you disagree with
  it from a position of having done the work.
]
