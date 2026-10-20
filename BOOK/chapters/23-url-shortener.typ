#import "../template.typ": *

= Project: a URL shortener

Everything from Parts I to III in one finished program: an HTTP server, a
database, concurrency, context, error design, tests. Around three hundred lines,
no framework, and it is genuinely deployable.

This chapter is one long project. Type it as you go.

== What it does

```
POST /shorten  {"url": "https://example.com/very/long"}
               ->  {"code":"aB3dQ"}
GET  /aB3dQ                               ->  302 redirect
GET  /api/stats/aB3dQ                     ->  {"hits": 17}
```

Plus: it refuses nonsense URLs, counts hits without slowing the redirect, and
shuts down cleanly.

== The schema and the store

```go
// store.go
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    _ "modernc.org/sqlite"
)

var (
    ErrNotFound = errors.New("short code not found")
    ErrInvalid  = errors.New("invalid url")
)

type Link struct {
    Code      string    `json:"code"`
    URL       string    `json:"url"`
    Hits      int64     `json:"hits"`
    CreatedAt time.Time `json:"created_at"`
}

type Store struct{ db *sql.DB }

func OpenStore(path string) (*Store, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("opening %s: %w", path, err)
    }
    db.SetMaxOpenConns(1)
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("connecting: %w", err)
    }
    if _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS links (
            code       TEXT PRIMARY KEY,
            url        TEXT    NOT NULL,
            hits       INTEGER NOT NULL DEFAULT 0,
            created_at TEXT    NOT NULL
        )`); err != nil {
        return nil, fmt.Errorf("creating schema: %w", err)
    }
    return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Save(ctx context.Context,
    code, url string) (*Link, error) {
    now := time.Now().UTC()
    _, err := s.db.ExecContext(ctx,
        `INSERT INTO links (code, url, created_at) VALUES (?, ?, ?)`,
        code, url, now.Format(time.RFC3339))
    if err != nil {
        return nil, fmt.Errorf("saving %s: %w", code, err)
    }
    return &Link{Code: code, URL: url, CreatedAt: now}, nil
}

func (s *Store) Lookup(ctx context.Context, code string) (*Link, error) {
    var l Link
    var created string
    err := s.db.QueryRowContext(ctx,
        `SELECT code, url, hits, created_at FROM links
         WHERE code = ?`, code,
    ).Scan(&l.Code, &l.URL, &l.Hits, &created)

    if errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("%q: %w", code, ErrNotFound)
    }
    if err != nil {
        return nil, fmt.Errorf("looking up %q: %w", code, err)
    }
    l.CreatedAt, _ = time.Parse(time.RFC3339, created)
    return &l, nil
}

// RecordHit is deliberately separate from Lookup so the redirect does not
// wait for a write.
func (s *Store) RecordHit(ctx context.Context, code string) error {
    _, err := s.db.ExecContext(ctx,
        `UPDATE links SET hits = hits + 1 WHERE code = ?`, code)
    if err != nil {
        return fmt.Errorf("recording hit for %q: %w", code, err)
    }
    return nil
}
```

== Generating codes

```go
// code.go
package main

import (
    "crypto/rand"
    "math/big"
)

// Base58: no 0, O, I or l, so a code read aloud or copied by hand
// survives.
const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZ" +
    "abcdefghijkmnopqrstuvwxyz"

func newCode(n int) (string, error) {
    b := make([]byte, n)
    max := big.NewInt(int64(len(alphabet)))
    for i := range b {
        k, err := rand.Int(rand.Reader, max)
        if err != nil {
            return "", err
        }
        b[i] = alphabet[k.Int64()]
    }
    return string(b), nil
}
```

#gotcha[
  `crypto/rand`, not `math/rand`. A short code generated from a predictable
  sequence can be guessed, which for a link shortener means enumerating
  everybody's private links. Use `crypto/rand` for anything a user should not
  be able to predict: codes, tokens, session IDs, password reset links.
]

== Validating the URL

```go
// validate.go
package main

import (
    "fmt"
    "net/url"
    "strings"
)

func validateURL(raw string) (string, error) {
    raw = strings.TrimSpace(raw)
    if raw == "" {
        return "", fmt.Errorf("url is empty: %w", ErrInvalid)
    }
    if len(raw) > 2048 {
        return "", fmt.Errorf("url too long: %w", ErrInvalid)
    }

    u, err := url.Parse(raw)
    if err != nil {
        return "", fmt.Errorf("parsing url: %w", ErrInvalid)
    }
    // Scheme must be http or https: javascript: and data: are XSS vectors
    // the moment this code appears in an href.
    if u.Scheme != "http" && u.Scheme != "https" {
        return "", fmt.Errorf("scheme %q not allowed: %w",
            u.Scheme, ErrInvalid)
    }
    if u.Host == "" {
        return "", fmt.Errorf("url has no host: %w", ErrInvalid)
    }
    return u.String(), nil
}
```

#gotcha[
  The scheme check is not optional. Without it, somebody shortens
  `javascript:alert(document.cookie)` and your redirect endpoint serves it. Any
  time user input ends up in an `href` or a `Location` header, allowlist the
  scheme.
]

== The server

```go
// main.go
package main

import (
    "context"
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Server struct {
    store *Store
    log   *log.Logger
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
    var in struct {
        URL string `json:"url"`
    }
    dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10))
    dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil {
        writeErr(w, http.StatusBadRequest, "invalid json body")
        return
    }

    clean, err := validateURL(in.URL)
    if err != nil {
        writeErr(w, http.StatusUnprocessableEntity, err.Error())
        return
    }

    // Retry on the astronomically unlikely collision rather than failing.
    for attempt := 0; attempt < 5; attempt++ {
        code, err := newCode(6)
        if err != nil {
            s.log.Printf("generating code: %v", err)
            writeErr(w, http.StatusInternalServerError, "internal error")
            return
        }
        link, err := s.store.Save(r.Context(), code, clean)
        if err == nil {
            writeJSON(w, http.StatusCreated, link)
            return
        }
        s.log.Printf("save attempt %d: %v", attempt+1, err)
    }
    writeErr(w, http.StatusInternalServerError,
        "could not allocate a code")
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
    code := r.PathValue("code")

    link, err := s.store.Lookup(r.Context(), code)
    if errors.Is(err, ErrNotFound) {
        writeErr(w, http.StatusNotFound, "no such link")
        return
    }
    if err != nil {
        s.log.Printf("lookup %q: %v", code, err)
        writeErr(w, http.StatusInternalServerError, "internal error")
        return
    }

    // Count the hit without making the user wait for the write. The
    // request
    // context is about to be cancelled, so the goroutine gets its own.
    go func() {
        ctx, cancel := context.WithTimeout(
            context.Background(), 5*time.Second)
        defer cancel()
        if err := s.store.RecordHit(ctx, code); err != nil {
            s.log.Printf("recording hit: %v", err)
        }
    }()

    http.Redirect(w, r, link.URL, http.StatusFound)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
    link, err := s.store.Lookup(r.Context(), r.PathValue("code"))
    if errors.Is(err, ErrNotFound) {
        writeErr(w, http.StatusNotFound, "no such link")
        return
    }
    if err != nil {
        writeErr(w, http.StatusInternalServerError, "internal error")
        return
    }
    writeJSON(w, http.StatusOK, link)
}

func (s *Server) routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /shorten", s.handleShorten)
    mux.HandleFunc("GET /api/stats/{code}", s.handleStats)
    mux.HandleFunc("GET /{code}", s.handleRedirect)
    mux.HandleFunc("GET /healthz",
        func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    return s.recoverer(s.logger(mux))
}

func (s *Server) logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        s.log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func (s *Server) recoverer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                s.log.Printf("panic on %s %s: %v",
                    r.Method, r.URL.Path, rec)
                writeErr(w, http.StatusInternalServerError,
                    "internal error")
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func main() {
    addr := os.Getenv("ADDR")
    if addr == "" {
        addr = ":8080"
    }
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "links.db"
    }

    logger := log.New(os.Stdout, "", log.LstdFlags)

    store, err := OpenStore(dbPath)
    if err != nil {
        logger.Fatalf("opening store: %v", err)
    }
    defer store.Close()

    srv := &http.Server{
        Addr:              addr,
        Handler:           (&Server{store: store, log: logger}).routes(),
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       15 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil &&
            !errors.Is(err, http.ErrServerClosed) {
            logger.Fatalf("listen: %v", err)
        }
    }()
    logger.Printf("listening on %s, database %s", addr, dbPath)

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    logger.Println("shutting down")
    ctx, cancel := context.WithTimeout(
        context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Printf("forced shutdown: %v", err)
    }
    logger.Println("stopped")
}
```

== Running it

```bash
go mod init github.com/you/shortener
go get modernc.org/sqlite
go run .
```

```bash
curl -X POST localhost:8080/shorten -d '{"url":"https://go.dev/doc/"}'
# {"code":"7fKpQ2","url":"https://go.dev/doc/","hits":0,...}

curl -i localhost:8080/7fKpQ2          # 302, Location: https://go.dev/doc/
curl localhost:8080/api/stats/7fKpQ2   # {"hits":1,...}

curl -X POST localhost:8080/shorten -d '{"url":"javascript:alert(1)"}'
# 422, scheme "javascript" not allowed
```

== The tests

```go
// server_test.go
package main

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "path/filepath"
    "strings"
    "testing"
)

func newTestServer(t *testing.T) *Server {
    t.Helper()
    store, err := OpenStore(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("opening store: %v", err)
    }
    t.Cleanup(func() { store.Close() })
    return &Server{store: store, log: testLogger(t)}
}

func TestShortenAndRedirect(t *testing.T) {
    s := newTestServer(t)
    h := s.routes()

    req := httptest.NewRequest("POST", "/shorten",
        strings.NewReader(`{"url":"https://example.com/x"}`))
    rec := httptest.NewRecorder()
    h.ServeHTTP(rec, req)

    if rec.Code != http.StatusCreated {
        t.Fatalf("status = %d, want 201: %s", rec.Code, rec.Body)
    }
    var link Link
    if err := json.NewDecoder(rec.Body).Decode(&link); err != nil {
        t.Fatalf("decoding: %v", err)
    }
    if len(link.Code) != 6 {
        t.Errorf("code = %q, want 6 characters", link.Code)
    }

    req = httptest.NewRequest("GET", "/"+link.Code, nil)
    rec = httptest.NewRecorder()
    h.ServeHTTP(rec, req)

    if rec.Code != http.StatusFound {
        t.Fatalf("redirect status = %d, want 302", rec.Code)
    }
    got := rec.Header().Get("Location")
    if got != "https://example.com/x" {
        t.Errorf("Location = %q", got)
    }
}

func TestRejectsBadSchemes(t *testing.T) {
    s := newTestServer(t)
    h := s.routes()

    for _, bad := range []string{
        `{"url":"javascript:alert(1)"}`,
        `{"url":"data:text/html,<script>"}`,
        `{"url":""}`,
        `{"url":"not a url at all"}`,
    } {
        req := httptest.NewRequest("POST", "/shorten",
          strings.NewReader(bad))
        rec := httptest.NewRecorder()
        h.ServeHTTP(rec, req)
        if rec.Code != http.StatusUnprocessableEntity {
            t.Errorf("%s: status = %d, want 422", bad, rec.Code)
        }
    }
}

func TestUnknownCodeIs404(t *testing.T) {
    s := newTestServer(t)
    rec := httptest.NewRecorder()
    s.routes().ServeHTTP(rec, httptest.NewRequest("GET", "/nope12", nil))
    if rec.Code != http.StatusNotFound {
        t.Errorf("status = %d, want 404", rec.Code)
    }
}

func TestLookupMissing(t *testing.T) {
    s := newTestServer(t)
    _, err := s.store.Lookup(context.Background(), "missing")
    if err == nil {
        t.Fatal("want an error for a missing code")
    }
}
```

```bash
go test -race -cover ./...
```

#project(title: "What this program demonstrates")[
  Walk back through it and name the chapter each decision came from.

  *Error design (11).* `ErrNotFound` and `ErrInvalid` are sentinels, wrapped
  with `%w`, and the handlers branch on `errors.Is`. No handler reads an error
  message.

  *Pointers and nil (8).* `Lookup` returns `*Link`, nil on failure, so a caller
  cannot accidentally use a half-built value.

  *Goroutines and context (13, 16).* The hit counter runs in its own goroutine
  with a *fresh* context, because `r.Context()` is cancelled the moment the
  redirect is written. Using the request context there would silently lose most
  of your statistics under load, and it is the kind of bug that never shows up
  in testing.

  *Interfaces (10).* `writeJSON` takes `http.ResponseWriter`, and the tests
  pass an `httptest.ResponseRecorder` instead. Nothing was designed for
  testability; it is testable because the standard library is built on small
  interfaces.

  *The database chapter (22).* `defer rows.Close()`, `rows.Err()`,
  placeholders, `sql.ErrNoRows` translated at the boundary.

  *Security.* `crypto/rand` for unguessable codes, a scheme allowlist against
  XSS, `MaxBytesReader` against oversized bodies, `DisallowUnknownFields`
  against typos, and timeouts on the server.

  Count the dependencies: one, the SQLite driver.
]

#tryit[
  Three extensions, in order of difficulty. Add an expiry column and a
  `GET /{code}` that returns 410 Gone for expired links. Add rate limiting to
  `POST /shorten` using a map of IP to a token bucket, protected by a mutex,
  which is chapter fifteen's material. Then add a custom alias
  (`{"url":..., "alias":"docs"}`), which forces you to handle a real
  `UNIQUE` constraint violation rather than the random collision retry, and
  that is a different error path entirely.
]
