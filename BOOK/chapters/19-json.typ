#import "../template.typ": *

= JSON, struct tags, and the capital letter that bites

JSON is where your Go backend meets your React frontend, so this chapter earns
its place twice: once now, and again in part four when Grit generates the types
on both sides of that boundary.

== Encoding and decoding

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

u := User{ID: 1, Name: "Ada", Email: "ada@example.com"}

data, err := json.Marshal(u)
// {"id":1,"name":"Ada","email":"ada@example.com"}

var back User
err = json.Unmarshal(data, &back)     // note the &
```

`Unmarshal` takes a pointer, because it fills in the value you give it. Passing
`back` instead of `&back` compiles and silently does nothing useful, which is
the chapter's first trap.

For readable output:

```go
data, _ := json.MarshalIndent(u, "", "  ")
```

== The capital letter rule, again

This is the single most common Go and JSON bug, and you met the rule in chapter
two.

```go
type User struct {
    id    int       // lowercase: unexported
    name  string
}

json.Marshal(User{id: 1, name: "Ada"})   // {}
```

An empty object. No error, no warning. The `encoding/json` package lives in
another package and literally cannot see unexported fields.

#gotcha[
  *If your JSON comes out as `{}`, check for lowercase field names first.* It
  is the answer about nine times out of ten. The same applies on the way in: an
  unexported field is never populated by `Unmarshal`.
]

== Struct tags in full

```go
type Product struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    Notes     string    `json:"notes,omitempty"`   // omitted when ""
    Internal  string    `json:"-"`                 // never encoded
    CreatedAt time.Time `json:"created_at"`
    Count     int       `json:"count,string"`      // encoded as "42"
}
```

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [`json:"name"`], [Use this key],
  [`json:"-"`], [Never encode or decode this field],
  [`json:",omitempty"`], [Omit when the value is the zero value],
  [`json:",string"`], [Encode a number or bool as a JSON string],
  [no tag], [Use the Go field name, capital included],
)

#gotcha[
  The tag syntax is unforgiving and silent. All of these are wrong and none
  produce an error:

  ```go
  Name string `json: "name"`     // space after the colon: tag ignored
  Name string `json:'name'`      // single quotes: tag ignored
  Name string "json:\"name\""    // works, but use backticks
  ```

  The field then encodes as `Name`. If one key in your JSON is capitalised and
  the rest are not, you have found a stray space.

  `go vet` catches malformed tags. Run it.
]

== omitempty, and what it actually means

`omitempty` omits the *zero value*, which is rarely what people assume.

```go
type Filter struct {
    Page     int    `json:"page,omitempty"`      // 0 is omitted
    Active   bool   `json:"active,omitempty"`    // false is omitted!
    Query    string `json:"query,omitempty"`     // "" is omitted
}
```

`Active: false` disappears from the output. If false is a meaningful value and
not merely "unset", `omitempty` is wrong and will cause a real bug.

The fix is a pointer, which is chapter eight's lesson in its natural habitat:

```go
type Filter struct {
    Active *bool `json:"active,omitempty"`   // nil omitted, &false encoded
}
```

#bridge[
  JavaScript has no such ambiguity because `undefined` and `false` are
  different values. Go's zero values collapse them, so a `*bool` is how you get
  the distinction back. You will see `*string` and `*bool` all over the
  struct definitions of any API that supports partial updates, and now you know
  exactly why.
]

== Decoding unknown shapes

Without a struct, JSON decodes into `map[string]any`:

```go
var data map[string]any
json.Unmarshal(raw, &data)

name := data["name"].(string)       // panics if it is not a string
name, ok := data["name"].(string)   // safe
```

Every number becomes a `float64`, which surprises everyone:

```go
json.Unmarshal([]byte(`{"count": 42}`), &data)
n := data["count"].(int)        // PANIC: it is a float64
n := int(data["count"].(float64))   // correct
```

Prefer a struct whenever you know the shape. The type system is the point.

== Decoding is case-insensitive, and lenient

Two behaviours that catch people:

```go
type User struct {
    Name string `json:"name"`
}

json.Unmarshal([]byte(`{"NAME": "Ada"}`), &u)   // works: case-insensitive
// works: unknown fields ignored
json.Unmarshal([]byte(`{"nope": 1}`), &u)
```

Unknown fields are silently dropped. For an API that should reject nonsense,
use a decoder:

```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()
if err := dec.Decode(&input); err != nil {
    http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
    return
}
```

This turns a typo in a client's payload into a 400 instead of a silently
ignored field, which is usually what you want for a write endpoint.

== Streaming with Encoder and Decoder

For an HTTP body, do not read it all into memory first:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    var input CreateUser
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

`Decoder` reads from any `io.Reader` and `Encoder` writes to any `io.Writer`,
which is chapter ten's small interfaces paying off again.

== Times, and the format you will meet

`time.Time` marshals to RFC 3339:

```json
{"created_at": "2026-10-20T14:30:00Z"}
```

That parses directly with `new Date(...)` in JavaScript, so the default is the
right one. A zero time encodes as `"0001-01-01T00:00:00Z"`, which looks
alarming in output; use a `*time.Time` with `omitempty` when the field is
genuinely optional.

== Custom marshalling

Implement `MarshalJSON` to control a type's representation:

```go
type Money int64      // stored as cents

func (m Money) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf(`"%d.%02d"`, m/100, m%100)), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    // parse "12.34" into 1234 cents
    ...
}
```

Note the receivers: `MarshalJSON` on the value, `UnmarshalJSON` on the pointer,
because unmarshalling has to modify. Getting this backwards means your custom
decoder is silently never called, which is the second-most-common JSON bug
after lowercase fields.

#bridge[
  This is `toJSON()` and a reviver function, combined into an interface. Money
  is the classic case: a float64 cannot represent 0.10 exactly, so storing
  cents as an integer and formatting at the boundary is the correct design in
  any language. Grit has a `money` field type that does exactly this, and part
  four explains why it is a type rather than advice.
]

#project(title: "Project 18: a JSON reformatting CLI")[
  Build `jsonfmt`: read JSON on standard input or from a file, validate it,
  pretty print it, and optionally extract one field. This is a tool you will
  actually use, and it exercises streaming, unknown shapes and error handling.

  ```go
  package main

  import (
      "encoding/json"
      "flag"
      "fmt"
      "io"
      "os"
      "strings"
  )

  func main() {
      indent := flag.String("indent", "  ", "indent string")
      compact := flag.Bool("c", false,
          "compact output instead of indented")
      path := flag.String("get", "",
          "dotted path to extract, e.g. user.name")
      flag.Parse()

      // Read from a named file, or standard input when piped.
      var r io.Reader = os.Stdin
      if args := flag.Args(); len(args) > 0 {
          f, err := os.Open(args[0])
          if err != nil {
              fail("opening %s: %v", args[0], err)
          }
          defer f.Close()
          r = f
      }

      // any, not map[string]any: the document might be an array.
      var doc any
      dec := json.NewDecoder(r)
      dec.UseNumber()            // keep numbers exact instead of float64
      if err := dec.Decode(&doc); err != nil {
          fail("invalid json: %v", err)
      }

      if *path != "" {
          v, err := get(doc, *path)
          if err != nil {
              fail("%v", err)
          }
          doc = v
      }

      enc := json.NewEncoder(os.Stdout)
      if !*compact {
          enc.SetIndent("", *indent)
      }
      if err := enc.Encode(doc); err != nil {
          fail("writing: %v", err)
      }
  }

  // get walks a dotted path through decoded JSON.
  func get(doc any, path string) (any, error) {
      cur := doc
      for _, part := range strings.Split(path, ".") {
          m, ok := cur.(map[string]any)
          if !ok {
              return nil, fmt.Errorf("%q: not an object at %q", path, part)
          }
          cur, ok = m[part]
          if !ok {
              return nil, fmt.Errorf("%q: no such key %q", path, part)
          }
      }
      return cur, nil
  }

  func fail(format string, args ...any) {
      fmt.Fprintf(os.Stderr, "jsonfmt: "+format+"\n", args...)
      os.Exit(1)
  }
  ```

  Build it and try it:

  ```bash
  go build -o jsonfmt .
  echo '{"user":{"name":"Ada","id":1},"tags":["a","b"]}' | ./jsonfmt
  echo '{"user":{"name":"Ada"}}' | ./jsonfmt -get user.name
  echo '{"bad": }' | ./jsonfmt           # exits 1 with a message on stderr
  ```

  Two details worth understanding. `dec.UseNumber()` stops every number
  becoming a `float64`, so a large integer ID survives the round trip instead
  of losing precision: without it, `{"id": 12345678901234567}` comes back
  subtly wrong. And `doc` is declared as `any` rather than `map[string]any`, so
  a top-level JSON array works too.

  The `...any` in `fail` is a variadic parameter, Go's rest argument. Passing
  it on to `Fprintf` with `args...` spreads it again, which is the Go
  equivalent of `...args` in both directions.
]

#tryit[
  Add a `-keys` flag that prints the top-level keys one per line, sorted.
  You will need a type switch on `doc`, because the document might be an object
  or an array, and the two cases want different answers. Then add
  `DisallowUnknownFields` behind a `-strict` flag and find out what it does to
  a document with no struct to compare against. The answer tells you something
  real about what that option is actually for.
]
