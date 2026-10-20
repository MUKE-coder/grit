#import "../template.typ": *

= Talking to a database with database/sql

Go's `database/sql` is not an ORM. It is a thin, driver-agnostic layer over
SQL: you write the queries, it manages connections and scans rows into your
variables. Coming from Prisma or Drizzle this feels primitive for about a day
and then feels like relief, because nothing is hidden.

== Opening a database

```go
import (
    "database/sql"
    // the driver, imported for side effects only
    _ "modernc.org/sqlite"
)

db, err := sql.Open("sqlite", "app.db")
if err != nil {
    return fmt.Errorf("opening database: %w", err)
}
defer db.Close()

if err := db.Ping(); err != nil {
    return fmt.Errorf("connecting: %w", err)
}
```

The blank import registers the driver with `database/sql`. You never call the
driver package directly, which is why it is imported as `_`.

#gotcha[
  *`sql.Open` does not connect.* It validates arguments and prepares a pool,
  lazily. A wrong password or a missing server is not reported until the first
  query. Always `Ping` at startup, or your service starts happily and fails on
  the first user request.
]

== The pool is the important part

`*sql.DB` is not a connection. It is a pool, safe for concurrent use, and you
should create exactly one and pass it around.

```go
db.SetMaxOpenConns(25)                   // hard ceiling on connections
db.SetMaxIdleConns(25)                   // keep them warm
db.SetConnMaxLifetime(5 * time.Minute)   // recycle, for load balancers
db.SetConnMaxIdleTime(5 * time.Minute)
```

#gotcha[
  The default `MaxOpenConns` is *unlimited*. Under load, Go will happily open
  more connections than your database allows, and Postgres will start refusing
  them. Set it. A common starting point is 25, and the right number is
  something like `(cores * 2) + effective spindles`, measured rather than
  guessed.

  Set `MaxIdleConns` to the same value as `MaxOpenConns`. If idle is lower,
  connections are closed and reopened constantly under steady load, which is
  pure latency.
]

== Querying one row

```go
func GetUser(ctx context.Context, db *sql.DB, id int) (*User, error) {
    var u User
    err := db.QueryRowContext(ctx,
        `SELECT id, name, email FROM users WHERE id = ?`, id,
    ).Scan(&u.ID, &u.Name, &u.Email)

    if errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)
    }
    if err != nil {
        return nil, fmt.Errorf("querying user %d: %w", id, err)
    }
    return &u, nil
}
```

`Scan` copies columns into the pointers you give it, in order. The number and
order must match the `SELECT`.

`sql.ErrNoRows` is the sentinel for "no row", and it is a normal outcome rather
than a failure. Translating it into your own `ErrNotFound`, as above, is
chapter eleven's design in practice: callers should not have to import
`database/sql` to understand your errors.

== Querying many rows

```go
func ListUsers(ctx context.Context, db *sql.DB) ([]User, error) {
    rows, err := db.QueryContext(ctx,
        `SELECT id, name, email FROM users ORDER BY id`)
    if err != nil {
        return nil, fmt.Errorf("listing users: %w", err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            return nil, fmt.Errorf("scanning user: %w", err)
        }
        users = append(users, u)
    }
    // The check everybody forgets: the loop ends the same way on error.
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("iterating users: %w", err)
    }
    return users, nil
}
```

Three things are mandatory and two of them are easy to forget.

*`defer rows.Close()`.* An unclosed `Rows` holds its connection. Leak enough
and the pool is exhausted and every query blocks forever.

*Check `rows.Err()` after the loop.* Exactly like `scanner.Err()` in chapter
twenty, and for the same reason.

*Check the `Scan` error inside the loop.*

== Writing, and what you get back

```go
res, err := db.ExecContext(ctx,
    `INSERT INTO users (name, email) VALUES (?, ?)`, name, email)
if err != nil {
    return fmt.Errorf("creating user: %w", err)
}

id, err := res.LastInsertId()     // SQLite and MySQL
n, err := res.RowsAffected()
```

#gotcha[
  `LastInsertId` is not supported by Postgres. There you use `RETURNING` and
  treat the insert as a query:

  ```go
  var id int
  err := db.QueryRowContext(ctx,
      `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
      name, email,
  ).Scan(&id)
  ```

  Note also the placeholders: SQLite and MySQL use `?`, Postgres uses `$1`,
  `$2`. This is the main thing that is not portable between drivers, and it is
  why switching databases is never quite free.
]

== Placeholders are not optional

```go
// NEVER do this.
query := "SELECT * FROM users WHERE email = '" + email + "'"

// Always this.
db.QueryRowContext(ctx, "SELECT ... WHERE email = ?", email)
```

A placeholder sends the value separately from the statement, so it can never be
parsed as SQL. String concatenation is SQL injection, and `'; DROP TABLE
users;--` is only the famous example rather than the dangerous one: the real
attacks read data rather than destroy it.

#bridge[
  Prisma and Drizzle protect you by generating the SQL. Here you write the SQL
  and the protection is the placeholder. The rule is simple: *values are
  placeholders, always*. Table and column names cannot be placeholders, so if
  those are dynamic, validate them against an allowlist you control. Grit does
  exactly this, which is why its generated services carry a `Sortable` map of
  permitted columns rather than interpolating whatever the client sent.
]

== NULL, and why your struct fields lie

A `NULL` column cannot scan into a `string`. It fails at runtime:

```
sql: Scan error on column index 2: converting NULL to string is unsupported
```

Three options, in order of preference:

```go
// 1. Best: make the column NOT NULL DEFAULT '' in the schema.

// 2. A pointer, when NULL is genuinely meaningful.
var bio *string
rows.Scan(&u.ID, &u.Name, &bio)

// 3. sql.NullString, when you want the distinction explicit.
var bio sql.NullString
rows.Scan(&u.ID, &u.Name, &bio)
if bio.Valid {
    u.Bio = bio.String
}
```

Since Go 1.22 there is also `sql.Null[T]`, the generic version:

```go
var bio sql.Null[string]
```

Most NULL problems are schema problems. If a column has no meaningful "absent"
state, make it `NOT NULL` and the whole question disappears.

== Transactions

```go
func Transfer(ctx context.Context, db *sql.DB,
    from, to int, amount int64) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("beginning transaction: %w", err)
    }
    // Safe to call after Commit: it becomes a no-op. This is the guard
    // that
    // makes every early return below correct.
    defer tx.Rollback()

    _, err = tx.ExecContext(ctx,
        `UPDATE accounts SET balance = balance - ? WHERE id = ?`,
        amount, from)
    if err != nil {
        return fmt.Errorf("debiting %d: %w", from, err)
    }

    var balance int64
    err = tx.QueryRowContext(ctx,
        `SELECT balance FROM accounts WHERE id = ?`, from).Scan(&balance)
    if err != nil {
        return fmt.Errorf("checking balance: %w", err)
    }
    if balance < 0 {
        return ErrInsufficientFunds      // the deferred Rollback undoes it
    }

    _, err = tx.ExecContext(ctx,
        `UPDATE accounts SET balance = balance + ? WHERE id = ?`,
           amount, to)
    if err != nil {
        return fmt.Errorf("crediting %d: %w", to, err)
    }

    return tx.Commit()
}
```

`defer tx.Rollback()` right after `BeginTx` is the pattern. After a successful
`Commit` the rollback does nothing, so one line covers every failure path
including a panic.

#gotcha[
  *Everything inside a transaction must use `tx`, not `db`.* A query on `db`
  inside a transaction takes a *different connection* from the pool, so it
  cannot see uncommitted changes and does not participate in the rollback.
  This compiles, passes casual testing, and corrupts data under concurrency.

  The usual defence is to pass something that could be either:

  ```go
  type Querier interface {
      QueryRowContext(context.Context, string, ...any) *sql.Row
      ExecContext(context.Context, string, ...any) (sql.Result, error)
  }
  ```

  Both `*sql.DB` and `*sql.Tx` satisfy it, so a function can be called inside
  or outside a transaction. This is chapter ten's consumer-declared interface,
  doing real work.
]

== The N+1 query problem

The most common performance bug in any application, in any language:

```go
posts, _ := ListPosts(ctx, db)              // 1 query
for i := range posts {
    posts[i].Author, _ = GetUser(ctx, db, posts[i].AuthorID)   // N queries
}
```

Fifty posts means fifty-one queries. Each is fast, so nothing looks slow, and
the page takes two seconds.

The fix is one more query, not fifty:

```go
ids := make([]int, 0, len(posts))
for _, p := range posts {
    ids = append(ids, p.AuthorID)
}

// Build the right number of placeholders.
ph := strings.Repeat("?,", len(ids))
ph = ph[:len(ph)-1]
args := make([]any, len(ids))
for i, id := range ids {
    args[i] = id
}

rows, err := db.QueryContext(ctx,
    `SELECT id, name FROM users WHERE id IN (`+ph+`)`, args...)
// ... build map[int]User, then attach to posts
```

Note that the *placeholders* are generated and the *values* are still passed as
arguments. That is safe. Interpolating the IDs themselves would not be.

#bridge[
  Prisma hides this behind `include: { author: true }`, and under the hood it
  does the same two-query dance. Here you do it yourself, which is more typing
  and means the query count is visible in your code rather than in a log you
  have to go looking for. In part four you will see Grit's profiler report
  "201 queries, 0 slow" for exactly this pattern.
]

== Should you use an ORM?

You will meet GORM in part four, because Grit uses it. The honest comparison:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [`database/sql`], [You write SQL. Nothing hidden, no surprises, more code],
  [`sqlc`], [You write SQL, it generates typed Go. Excellent, and a build step],
  [GORM], [Structs and methods. Fast to write, and you must know what it emits],
)

The position worth holding: an ORM is a productivity tool, not a replacement
for knowing SQL. Every team that has been burned by an ORM was burned because
nobody could say what query it was running. Learn this chapter first, then use
whatever your project uses, and keep the query log open.

#project(title: "Project 21: a database-backed note store")[
  Wire a real SQLite database behind the HTTP service from chapter twenty one.
  Pure standard library plus one driver.

  ```bash
  go get modernc.org/sqlite      # pure Go, no cgo, no compiler needed
  ```

  ```go
  package main

  import (
      "context"
      "database/sql"
      "errors"
      "fmt"
      "time"

      _ "modernc.org/sqlite"
  )

  var ErrNotFound = errors.New("not found")

  type Note struct {
      ID        int       `json:"id"`
      Text      string    `json:"text"`
      CreatedAt time.Time `json:"created_at"`
  }

  type Store struct{ db *sql.DB }

  func Open(path string) (*Store, error) {
      db, err := sql.Open("sqlite", path)
      if err != nil {
          return nil, fmt.Errorf("opening %s: %w", path, err)
      }
      db.SetMaxOpenConns(1)     // SQLite: one writer. For Postgres, 25.

      if err := db.Ping(); err != nil {
          return nil, fmt.Errorf("connecting: %w", err)
      }

      // Schema on startup. Real projects use migration files; chapter 31
      // covers how Grit does it.
      _, err = db.Exec(`
          CREATE TABLE IF NOT EXISTS notes (
              id         INTEGER PRIMARY KEY AUTOINCREMENT,
              text       TEXT    NOT NULL,
              created_at TEXT    NOT NULL
          )`)
      if err != nil {
          return nil, fmt.Errorf("creating schema: %w", err)
      }
      return &Store{db: db}, nil
  }

  func (s *Store) Close() error { return s.db.Close() }

  func (s *Store) Create(ctx context.Context, text string) (*Note, error) {
      now := time.Now().UTC()
      res, err := s.db.ExecContext(ctx,
          `INSERT INTO notes (text, created_at) VALUES (?, ?)`,
          text, now.Format(time.RFC3339))
      if err != nil {
          return nil, fmt.Errorf("creating note: %w", err)
      }
      id, err := res.LastInsertId()
      if err != nil {
          return nil, fmt.Errorf("reading new id: %w", err)
      }
      return &Note{ID: int(id), Text: text, CreatedAt: now}, nil
  }

  func (s *Store) Get(ctx context.Context, id int) (*Note, error) {
      var n Note
      var created string
      err := s.db.QueryRowContext(ctx,
          `SELECT id, text, created_at FROM notes WHERE id = ?`, id,
      ).Scan(&n.ID, &n.Text, &created)

      if errors.Is(err, sql.ErrNoRows) {
          return nil, fmt.Errorf("note %d: %w", id, ErrNotFound)
      }
      if err != nil {
          return nil, fmt.Errorf("getting note %d: %w", id, err)
      }
      n.CreatedAt, _ = time.Parse(time.RFC3339, created)
      return &n, nil
  }

  func (s *Store) List(ctx context.Context, limit int) ([]Note, error) {
      rows, err := s.db.QueryContext(ctx,
          `SELECT id, text, created_at FROM notes
           ORDER BY id DESC LIMIT ?`, limit)
      if err != nil {
          return nil, fmt.Errorf("listing notes: %w", err)
      }
      defer rows.Close()

      var out []Note
      for rows.Next() {
          var n Note
          var created string
          if err := rows.Scan(&n.ID, &n.Text, &created); err != nil {
              return nil, fmt.Errorf("scanning note: %w", err)
          }
          n.CreatedAt, _ = time.Parse(time.RFC3339, created)
          out = append(out, n)
      }
      if err := rows.Err(); err != nil {
          return nil, fmt.Errorf("iterating notes: %w", err)
      }
      return out, nil
  }

  func (s *Store) Delete(ctx context.Context, id int) error {
      res, err := s.db.ExecContext(ctx,
          `DELETE FROM notes WHERE id = ?`, id)
      if err != nil {
          return fmt.Errorf("deleting note %d: %w", id, err)
      }
      n, err := res.RowsAffected()
      if err != nil {
          return fmt.Errorf("checking delete: %w", err)
      }
      if n == 0 {
          return fmt.Errorf("note %d: %w", id, ErrNotFound)
      }
      return nil
  }
  ```

  Now swap this in behind the handlers from chapter twenty one and you have a
  persistent JSON API with two dependencies: the SQLite driver and nothing
  else.

  Two details worth dwelling on. `RowsAffected() == 0` is how `Delete` knows
  the row was missing, because a `DELETE` that matches nothing is not an error
  in SQL. And `SetMaxOpenConns(1)` is correct for SQLite, which allows one
  writer: leaving it unlimited gives you `database is locked` errors under
  concurrency, which is the single most common SQLite complaint and is a
  configuration problem rather than a SQLite problem.
]

#tryit[
  Write a test using `t.TempDir()` for the database file, covering create, get,
  list and a `Get` that returns `ErrNotFound`. Then add an `Update` and give it
  a transaction that refuses to blank out the text, returning a validation
  error, and prove with a test that the row is unchanged after the failure.
  That last test is the one that proves you understand `defer tx.Rollback()`.
]
