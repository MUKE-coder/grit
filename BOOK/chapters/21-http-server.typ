#import "../template.typ": *

= net/http: a server with no framework

Go's standard library includes a production-grade HTTP server. Not a toy, not a
reference implementation: the thing you actually deploy. This chapter builds a
real one with no dependencies at all, because understanding this makes every
framework afterwards, Gin and Grit included, read as a convenience rather than
magic.

== Five lines

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "hello")
    })
    http.ListenAndServe(":8080", nil)
}
```

```bash
go run . &
curl localhost:8080/hello
```

No `npm install`, no framework, no config. That server handles concurrent
requests across every core, because `net/http` starts a goroutine per request.

#bridge[
  Express needs `npm install express` and gives you `app.get('/hello', (req,
  res) => res.send('hello'))`. The shapes are similar. The differences:
  Go's server is in the standard library and versioned with the language, and
  it is concurrent by default rather than on one event loop.
]

== The two core types

```go
func handler(w http.ResponseWriter, r *http.Request)
```

`*http.Request` is the incoming request: method, URL, headers, body, context.
`http.ResponseWriter` is where you write the response. Note the asymmetry: the
request is a pointer to a struct you read, the writer is an interface you call.

```go
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Method)                 // "GET"
    fmt.Println(r.URL.Path)               // "/users/42"
    fmt.Println(r.URL.Query().Get("q"))   // ?q=search
    fmt.Println(r.Header.Get("Authorization"))

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)     // 201
    w.Write([]byte(`{"ok":true}`))
}
```

#gotcha[
  *Order matters and the compiler will not help you.* Headers must be set
  before `WriteHeader`, and `WriteHeader` before `Write`. Writing the body
  first implicitly sends a 200, and your later `WriteHeader(404)` is ignored
  with a log line:

  ```
  http: superfluous response.WriteHeader call
  ```

  The rule: set headers, then status, then body. Always that order.
]

== Routing with ServeMux

Since Go 1.22 the standard router handles methods and path parameters, which
removed the main reason people reached for a third-party router.

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /users", listUsers)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)

http.ListenAndServe(":8080", mux)
```

Path parameters come out by name:

```go
func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    ...
}
```

A request that matches the path but not the method gets a 405 automatically.

#bridge[
  This is `app.get('/users/:id', ...)` with `req.params.id`. Two differences:
  the method is part of the pattern string rather than the function name, and
  there is no `next()`, because middleware composes by wrapping rather than by
  calling through.
]

== Always configure a Server

`http.ListenAndServe` uses a default server with *no timeouts at all*. A single
slow client can hold a connection open indefinitely, and enough of them
exhaust your file descriptors. This is a real denial of service and it is the
default.

```go
srv := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,   // slowloris defence
    ReadTimeout:       15 * time.Second,
    WriteTimeout:      15 * time.Second,
    IdleTimeout:       60 * time.Second,
    MaxHeaderBytes:    1 << 20,
}
log.Fatal(srv.ListenAndServe())
```

#gotcha[
  Set `ReadHeaderTimeout` at minimum. Without it, a client that opens a
  connection and sends one byte per minute occupies a goroutine and a file
  descriptor forever. This attack is called slowloris, it is twenty years old,
  and Go's defaults do not defend against it because the standard library
  cannot know your deployment.
]

== Middleware is a function that wraps a handler

There is no middleware API. A middleware is a function taking a handler and
returning a handler.

```go
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

Compose by nesting:

```go
handler := Logger(Recover(Auth(mux)))
```

The outermost runs first. Reading that line right to left gives you the request
path, which is why many codebases add a small helper:

```go
func Chain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
    for i := len(m) - 1; i >= 0; i-- {
        h = m[i](h)
    }
    return h
}

handler := Chain(mux, Logger, Recover, Auth)   // now reads in order
```

#bridge[
  Express middleware takes `(req, res, next)` and calls `next()` to continue.
  Go's takes a handler and *returns* one, so "continue" is calling
  `next.ServeHTTP(w, r)` and "stop" is simply not calling it:

  ```go
  func Auth(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          if r.Header.Get("Authorization") == "" {
              http.Error(w, "unauthorized", http.StatusUnauthorized)
              return                      // next is never called: stop here
          }
          next.ServeHTTP(w, r)
      })
  }
  ```

  No `next(err)`, no forgetting to call `next()` and hanging the request.
]

== Passing values down: context, not globals

Middleware that authenticates a user needs to hand it to the handler. The
mechanism is the request context from chapter sixteen.

```go
type ctxKey string
const userKey ctxKey = "user"

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        u, err := userFromToken(r.Header.Get("Authorization"))
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), userKey, u)
        next.ServeHTTP(w, r.WithContext(ctx))   // a new request, with the value
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    u, ok := r.Context().Value(userKey).(*User)
    if !ok {
        http.Error(w, "no user", http.StatusInternalServerError)
        return
    }
    ...
}
```

`r.WithContext` returns a *copy* of the request. Forgetting to pass the new one
to `next.ServeHTTP` is a common bug: the value is set on a request nobody sees.

== Graceful shutdown

Killing a server mid-request drops responses. Production code waits.

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: handler}

    go func() {
        if err := srv.ListenAndServe(); err != nil &&
            !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("listen: %v", err)
        }
    }()
    log.Println("listening on :8080")

    // Wait for an interrupt.
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    log.Println("shutting down")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("forced shutdown: %v", err)
    }
    log.Println("stopped")
}
```

`Shutdown` stops accepting new connections and waits for in-flight requests,
up to the context's deadline. `ListenAndServe` returns `http.ErrServerClosed`
when this happens, which is why that error is checked for and ignored.

This matters for every container deployment: Kubernetes sends `SIGTERM` and
waits before `SIGKILL`, and a server without this drops every request in
flight on each deploy.

== Testing handlers without a server

`net/http/httptest` runs a handler in memory.

```go
func TestGetUser(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/42", nil)
    req.SetPathValue("id", "42")
    rec := httptest.NewRecorder()

    getUser(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d, want 200", rec.Code)
    }
    var got User
    if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
        t.Fatalf("decoding body: %v", err)
    }
    if got.ID != 42 {
        t.Errorf("ID = %d, want 42", got.ID)
    }
}
```

No port, no network, no cleanup, and it runs in microseconds. For testing a
whole mux including middleware, `httptest.NewServer` starts a real one on a
random port.

#project(title: "Project 20: a complete HTTP service, no dependencies")[
  Everything above in one runnable file: routing, middleware, JSON, timeouts,
  graceful shutdown.

  ```go
  package main

  import (
      "context"
      "encoding/json"
      "errors"
      "log"
      "net/http"
      "os"
      "os/signal"
      "sync"
      "syscall"
      "time"
  )

  type Note struct {
      ID   int    `json:"id"`
      Text string `json:"text"`
  }

  type store struct {
      mu     sync.RWMutex
      notes  map[int]Note
      nextID int
  }

  func newStore() *store {
      return &store{notes: map[int]Note{}, nextID: 1}
  }

  func (s *store) list() []Note {
      s.mu.RLock()
      defer s.mu.RUnlock()
      out := make([]Note, 0, len(s.notes))
      for _, n := range s.notes {
          out = append(out, n)
      }
      return out
  }

  func (s *store) add(text string) Note {
      s.mu.Lock()
      defer s.mu.Unlock()
      n := Note{ID: s.nextID, Text: text}
      s.notes[n.ID] = n
      s.nextID++
      return n
  }

  // One helper for every JSON response, so the header is never forgotten.
  func writeJSON(w http.ResponseWriter, status int, v any) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(status)
      if err := json.NewEncoder(w).Encode(v); err != nil {
          log.Printf("encoding response: %v", err)
      }
  }

  func logger(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          start := time.Now()
          next.ServeHTTP(w, r)
          log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
      })
  }

  func recoverer(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          defer func() {
              if rec := recover(); rec != nil {
                  log.Printf("panic: %v", rec)
                  writeJSON(w, http.StatusInternalServerError,
                      map[string]string{"error": "internal error"})
              }
          }()
          next.ServeHTTP(w, r)
      })
  }

  func main() {
      s := newStore()
      mux := http.NewServeMux()

      mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
          writeJSON(w, http.StatusOK, s.list())
      })

      mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
          var in struct {
              Text string `json:"text"`
          }
          dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
          dec.DisallowUnknownFields()
          if err := dec.Decode(&in); err != nil {
              writeJSON(w, http.StatusBadRequest,
                  map[string]string{"error": "invalid json"})
              return
          }
          if in.Text == "" {
              writeJSON(w, http.StatusUnprocessableEntity,
                  map[string]string{"text": "must not be empty"})
              return
          }
          writeJSON(w, http.StatusCreated, s.add(in.Text))
      })

      mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
          panic("deliberate")      // proves the recoverer works
      })

      srv := &http.Server{
          Addr:              ":8080",
          Handler:           logger(recoverer(mux)),
          ReadHeaderTimeout: 5 * time.Second,
          ReadTimeout:       15 * time.Second,
          WriteTimeout:      15 * time.Second,
          IdleTimeout:       60 * time.Second,
      }

      go func() {
          if err := srv.ListenAndServe(); err != nil &&
              !errors.Is(err, http.ErrServerClosed) {
              log.Fatalf("listen: %v", err)
          }
      }()
      log.Println("listening on :8080")

      stop := make(chan os.Signal, 1)
      signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
      <-stop

      ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
      defer cancel()
      if err := srv.Shutdown(ctx); err != nil {
          log.Printf("forced: %v", err)
      }
      log.Println("stopped cleanly")
  }
  ```

  Exercise it:

  ```bash
  curl localhost:8080/notes
  curl -X POST localhost:8080/notes -d '{"text":"buy milk"}'
  curl -X POST localhost:8080/notes -d '{"txt":"typo"}'     # 400, unknown field
  curl -X POST localhost:8080/notes -d '{"text":""}'        # 422
  curl -X DELETE localhost:8080/notes                        # 405, automatic
  curl localhost:8080/panic                                  # 500, server survives
  ```

  Then press Ctrl+C and watch it log `stopped cleanly` rather than dying.

  Count the dependencies in `go.mod`: there are none. This is a concurrent,
  panic-resistant, gracefully shutting down JSON API with request logging and
  body size limits, in one file, with nothing installed. That is the baseline
  Go gives you, and it is why Go web frameworks are thin.
]

#tryit[
  Add `GET /notes/{id}` using `r.PathValue("id")` and `strconv.Atoi`, returning
  404 when absent. Then write the handler test with `httptest` and no running
  server. Finally, remove `ReadHeaderTimeout` and leave a connection open with
  `telnet localhost 8080` without sending anything, to see for yourself what
  the default costs you.
]
