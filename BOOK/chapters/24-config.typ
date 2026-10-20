#import "../template.typ": *

= Configuration, environment and secrets

A program that works on your laptop and nowhere else usually has its
configuration hard-coded. This short chapter is how to avoid that, and how not
to leak a database password into your git history on the way.

== Environment variables, the baseline

```go
addr := os.Getenv("ADDR")
if addr == "" {
    addr = ":8080"       // a default, so the program runs with no setup
}
```

`os.Getenv` returns `""` for a variable that is unset, so you cannot tell unset
from set-to-empty. When that matters:

```go
addr, ok := os.LookupEnv("ADDR")
if !ok {
    // genuinely not set
}
```

#bridge[
  This is `process.env.ADDR`, with the same string-only limitation. There is no
  `dotenv` in the standard library, so a `.env` file needs a package, usually
  `joho/godotenv`. Many Go teams skip it entirely and export variables from a
  shell script or a compose file, because production almost never uses a `.env`
  file anyway.
]

== Parse once, into a struct

Scattered `os.Getenv` calls are the thing to avoid. Read everything at startup,
validate it, and pass the result down.

```go
type Config struct {
    Addr        string
    DatabaseURL string
    LogLevel    string
    Debug       bool
    Timeout     time.Duration
    MaxWorkers  int
}

func Load() (*Config, error) {
    c := &Config{
        Addr:       getEnv("ADDR", ":8080"),
        LogLevel:   getEnv("LOG_LEVEL", "info"),
        Debug:      getEnvBool("DEBUG", false),
        Timeout:    getEnvDuration("TIMEOUT", 30*time.Second),
        MaxWorkers: getEnvInt("MAX_WORKERS", 4),
    }

    // Required: no sensible default exists.
    c.DatabaseURL = os.Getenv("DATABASE_URL")
    if c.DatabaseURL == "" {
        return nil, errors.New("DATABASE_URL is required")
    }
    if c.MaxWorkers < 1 {
        return nil, fmt.Errorf("MAX_WORKERS must be at least 1, got %d",
            c.MaxWorkers)
    }
    return c, nil
}

func getEnv(key, def string) string {
    if v, ok := os.LookupEnv(key); ok {
        return v
    }
    return def
}

func getEnvInt(key string, def int) int {
    if v, ok := os.LookupEnv(key); ok {
        if n, err := strconv.Atoi(v); err == nil {
            return n
        }
    }
    return def
}

func getEnvBool(key string, def bool) bool {
    if v, ok := os.LookupEnv(key); ok {
        if b, err := strconv.ParseBool(v); err == nil {
            return b
        }
    }
    return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
    if v, ok := os.LookupEnv(key); ok {
        if d, err := time.ParseDuration(v); err == nil {
            return d
        }
    }
    return def
}
```

`time.ParseDuration` accepts `"30s"`, `"5m"`, `"1h30m"`, which is far better
than an integer whose unit lives in a comment.

#gotcha[
  Those helpers silently fall back to the default when parsing fails, so
  `MAX_WORKERS=eight` quietly becomes 4. For anything that matters, fail
  loudly instead:

  ```go
  func mustInt(key string, def int) (int, error) {
      v, ok := os.LookupEnv(key)
      if !ok {
          return def, nil
      }
      n, err := strconv.Atoi(v)
      if err != nil {
          return 0, fmt.Errorf("%s=%q: not a number", key, v)
      }
      return n, nil
  }
  ```

  A typo in a deployment variable should stop the program at startup, not
  change its behaviour silently three weeks later.
]

== Fail at startup, not on the first request

```go
func main() {
    cfg, err := Load()
    if err != nil {
        log.Fatalf("configuration: %v", err)
    }
    // From here, cfg is valid. Nothing downstream re-checks it.
}
```

This is the single most valuable habit in this chapter. A missing variable
should crash the container immediately, so the deploy fails visibly, rather
than letting a half-configured process serve errors.

== Secrets

Three rules, and the first one is the one people break.

*Never commit a secret.* Not in a `.env` file, not in a test fixture, not
"temporarily". Git keeps it forever, and a repository that goes public later
publishes your history too.

```
# .gitignore
.env
.env.local
*.pem
```

Commit a `.env.example` with the *keys* and no values, so a new developer knows
what to set.

*Never log one.* This is easy to do by accident:

```go
log.Printf("config: %+v", cfg)    // prints DatabaseURL, password and all
```

The defence is a `String` method that redacts:

```go
func (c Config) String() string {
    return fmt.Sprintf("Config{Addr:%s LogLevel:%s DatabaseURL:%s}",
        c.Addr, c.LogLevel, redact(c.DatabaseURL))
}

func redact(s string) string {
    if s == "" {
        return ""
    }
    return "[redacted]"
}
```

Because `fmt` looks for `String()`, this covers `%v`, `%s` and `Println`. It
does not cover `%+v` on a struct that embeds yours, so the surer answer is a
dedicated type:

```go
type Secret string

func (Secret) String() string          { return "[redacted]" }
func (Secret) GoString() string        { return "[redacted]" }
func (s Secret) Reveal() string        { return string(s) }
```

Now the only way to see the value is to call `Reveal()`, which is greppable in
review.

*Never put one in a URL you log.* `postgres://user:pass@host/db` is a secret in
its entirety, which is why `redact` above hides the whole string rather than
trying to mask the password portion.

#bridge[
  The mechanics are the same as Node: environment variables, a gitignored
  `.env`, a secret manager in production. The Go-specific part is that a
  compiled binary feels like a sealed box and is not: strings compiled into it
  are readable with the `strings` command. Never embed a key in source, even in
  a binary you only ship internally.
]

== Flags, for tools rather than services

For a CLI, `flag` is in the standard library:

```go
func main() {
    var (
        verbose = flag.Bool("v", false, "verbose output")
        out     = flag.String("o", "", "output file (default stdout)")
        limit   = flag.Int("n", 10, "how many results")
    )
    flag.Parse()

    for _, path := range flag.Args() {   // the non-flag arguments
        ...
    }
}
```

`flag` gives you `-h` for free, printing every flag with its default. It does
not support subcommands or long-and-short pairs; for those, `spf13/cobra` is
the near-universal choice and is what Grit's CLI is built on.

The usual division: *flags for tools, environment for services*. A service
configured by flags is awkward in a container, and a CLI configured by
environment variables is awkward to use.

#project(title: "Project 22: configure the shortener properly")[
  Take the URL shortener from the last chapter and give it a real
  configuration layer, with validation that fails at startup.

  ```go
  // config.go
  package main

  import (
      "errors"
      "fmt"
      "os"
      "strconv"
      "time"
  )

  type Config struct {
      Addr         string
      DBPath       string
      CodeLength   int
      ReadTimeout  time.Duration
      WriteTimeout time.Duration
      BaseURL      string
  }

  func LoadConfig() (*Config, error) {
      c := &Config{
          Addr:    envOr("ADDR", ":8080"),
          DBPath:  envOr("DB_PATH", "links.db"),
          BaseURL: envOr("BASE_URL", "http://localhost:8080"),
      }

      var err error
      if c.CodeLength, err = envInt("CODE_LENGTH", 6); err != nil {
          return nil, err
      }
      c.ReadTimeout, err = envDur("READ_TIMEOUT", 15*time.Second)
      if err != nil {
          return nil, err
      }
      c.WriteTimeout, err = envDur("WRITE_TIMEOUT", 15*time.Second)
      if err != nil {
          return nil, err
      }

      // Validation. Every one of these is a mistake worth catching before
      // the first request rather than after.
      if c.CodeLength < 4 || c.CodeLength > 32 {
          return nil, fmt.Errorf("CODE_LENGTH must be 4 to 32, got %d",
              c.CodeLength)
      }
      if c.DBPath == "" {
          return nil, errors.New("DB_PATH must not be empty")
      }
      if c.ReadTimeout <= 0 || c.WriteTimeout <= 0 {
          return nil, errors.New("timeouts must be positive")
      }
      return c, nil
  }

  // String keeps the config loggable without leaking. Nothing here
  // is secret today, and the method exists so that it stays safe when
  // somebody adds an API key next month.
  func (c Config) String() string {
      return fmt.Sprintf(
          "Config{Addr:%s DBPath:%s CodeLength:%d BaseURL:%s}",
          c.Addr, c.DBPath, c.CodeLength, c.BaseURL)
  }

  func envOr(key, def string) string {
      if v, ok := os.LookupEnv(key); ok && v != "" {
          return v
      }
      return def
  }

  func envInt(key string, def int) (int, error) {
      v, ok := os.LookupEnv(key)
      if !ok || v == "" {
          return def, nil
      }
      n, err := strconv.Atoi(v)
      if err != nil {
          return 0, fmt.Errorf("%s=%q is not a number", key, v)
      }
      return n, nil
  }

  func envDur(key string, def time.Duration) (time.Duration, error) {
      v, ok := os.LookupEnv(key)
      if !ok || v == "" {
          return def, nil
      }
      d, err := time.ParseDuration(v)
      if err != nil {
          return 0, fmt.Errorf("%s=%q is not a duration, try 15s", key, v)
      }
      return d, nil
  }
  ```

  And in `main`:

  ```go
  cfg, err := LoadConfig()
  if err != nil {
      fmt.Fprintln(os.Stderr, "configuration:", err)
      os.Exit(1)
  }
  logger.Printf("starting with %v", cfg)
  ```

  Try it:

  ```bash
  go run .                               # defaults
  ADDR=:9000 CODE_LENGTH=8 go run .      # overrides
  CODE_LENGTH=2 go run .                 # exits 1 with a clear message
  CODE_LENGTH=eight go run .             # exits 1, names the bad value
  ```

  The last two are the point. A deployment with a typo'd variable now fails
  immediately and says which variable and what it received, instead of
  starting up and generating two-character codes that collide constantly.

  Write the test too, since config is easy to test and almost nobody does:

  ```go
  func TestLoadConfigRejectsShortCodes(t *testing.T) {
      t.Setenv("CODE_LENGTH", "2")      // restored automatically
      if _, err := LoadConfig(); err == nil {
          t.Fatal("want an error for CODE_LENGTH=2")
      }
  }
  ```

  `t.Setenv` sets a variable for one test and restores it afterwards, and it
  refuses to run in a parallel test, because environment variables are process
  global and that would be a race.
]

#tryit[
  Add an `API_KEY` using a `Secret` type with a redacting `String()`, and prove
  with a test that `fmt.Sprintf("%v", cfg)` does not contain the value. Then
  add `-config` as a flag pointing at a JSON file, with environment variables
  taking precedence over the file and the file over the defaults. Getting that
  precedence order right, and writing a test for each layer, is most of what
  configuration libraries do for you.
]
