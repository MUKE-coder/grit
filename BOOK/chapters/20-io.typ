#import "../template.typ": *

= Files, streams and the io interfaces

Two interfaces carry most of Go's I/O, and you met them in chapter ten:

```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```

One method each. Files, network connections, HTTP bodies, gzip streams,
in-memory buffers, standard input and standard output all satisfy them, which
means a function written against `io.Reader` works with every one of those
without knowing which it has.

== Reading a whole file

When the file is small and you want it all:

```go
data, err := os.ReadFile("config.json")   // []byte
if err != nil {
    return fmt.Errorf("reading config: %w", err)
}
text := string(data)
```

And writing:

```go
err := os.WriteFile("out.txt", []byte("hello"), 0o644)
```

`0o644` is the Unix permission mode: owner read and write, everyone else read.
On Windows it is mostly ignored. `0o600` for anything secret.

#gotcha[
  `os.ReadFile` loads the entire file into memory. For a 2GB log file that is
  2GB of RAM, and your process is killed. Use it for configuration and small
  documents; stream anything that could be large, which is the rest of this
  chapter.
]

== Streaming line by line

```go
f, err := os.Open("access.log")
if err != nil {
    return err
}
defer f.Close()

scanner := bufio.NewScanner(f)
for scanner.Scan() {
    line := scanner.Text()
    // ...
}
if err := scanner.Err(); err != nil {
    return fmt.Errorf("reading access.log: %w", err)
}
```

Memory stays constant regardless of file size.

#gotcha[
  *You must check `scanner.Err()`.* The loop ends identically whether the file
  finished or the disk failed, and without the check a read error looks exactly
  like a successful end. This is the most commonly skipped error check in Go.

  Second trap: `Scanner` refuses lines longer than 64KB by default, returning
  `bufio.Scanner: token too long`. For long lines, enlarge the buffer:

  ```go
  scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)   // 1MB lines
  ```
]

== os.Open against os.Create

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [`os.Open`], [Read only. Fails if the file is missing],
  [`os.Create`], [Write only. Creates, or truncates an existing file to empty],
  [`os.OpenFile`], [Everything else: append, read-write, custom permissions],
)

`os.Create` truncating is worth remembering, because it means opening a file to
"check something" with `Create` destroys it. To append:

```go
f, err := os.OpenFile("app.log",
    os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
```

== Buffer your writes

Every `Write` on a file is a system call. Writing a million lines one at a time
is a million system calls.

```go
f, err := os.Create("out.txt")
if err != nil {
    return err
}
defer f.Close()

w := bufio.NewWriter(f)
defer w.Flush()            // without this, you lose the tail of the file

for _, line := range lines {
    fmt.Fprintln(w, line)
}
```

#gotcha[
  *A buffered writer must be flushed.* The last few kilobytes sit in memory
  until you do. `defer w.Flush()` immediately after creating it is the habit,
  and note the order: `defer f.Close()` was registered first so it runs *last*,
  which is correct, because flushing into a closed file fails.

  Better still, check the flush error, because a failed flush means data did
  not reach the disk:

  ```go
  defer func() {
      if ferr := w.Flush(); ferr != nil && err == nil {
          err = fmt.Errorf("flushing: %w", ferr)
      }
  }()
  ```
]

== io.Copy, and the pattern it enables

```go
io.Copy(dst, src)      // stream everything from src to dst
```

That one line downloads a file to disk, with constant memory:

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()

f, err := os.Create("download.bin")
if err != nil {
    return err
}
defer f.Close()

written, err := io.Copy(f, resp.Body)
```

An HTTP response body is an `io.Reader`. A file is an `io.Writer`. Neither
knows about the other.

Useful companions:

```go
io.Discard                       // a Writer that throws everything away
// read at most 1MB: use on untrusted input
io.LimitReader(r, 1<<20)
io.MultiWriter(f, os.Stdout)     // write to both at once, like tee
io.TeeReader(r, logFile)         // read from r while copying to logFile
```

#gotcha[
  `io.LimitReader` on anything a user controls is a real defence. Without it, a
  request body or an uploaded file can be as large as the client wants, and
  reading it into memory is a denial of service that costs the attacker
  nothing. For HTTP specifically there is `http.MaxBytesReader`, which also
  returns a proper error.
]

== Writing testable code against interfaces

This is the practical payoff. A function that takes concrete types is hard to
test:

```go
// Hard: needs real files on disk.
func Process(inPath, outPath string) error
```

The same function over interfaces is trivial to test:

```go
func Process(r io.Reader, w io.Writer) error {
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        fmt.Fprintln(w, strings.ToUpper(scanner.Text()))
    }
    return scanner.Err()
}
```

```go
func TestProcess(t *testing.T) {
    in := strings.NewReader("hello\nworld\n")
    var out bytes.Buffer

    if err := Process(in, &out); err != nil {
        t.Fatal(err)
    }
    if got, want := out.String(), "HELLO\nWORLD\n"; got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

No temp files, no cleanup, no filesystem. `strings.NewReader` and
`bytes.Buffer` are the two types you will use most in tests.

#bridge[
  This is why Node's stream API exists, and Go's version is far smaller: one
  method per direction instead of an event emitter with `data`, `end`, `error`
  and backpressure semantics you have to learn. The tradeoff is that Go's
  streams are synchronous and blocking, which is fine because blocking a
  goroutine is cheap.
]

== Paths, portably

```go
import "path/filepath"

// data/users/ada.json or data\users\ada.json
filepath.Join("data", "users", "ada.json")
filepath.Ext("report.pdf")                   // ".pdf"
filepath.Base("/tmp/report.pdf")             // "report.pdf"
filepath.Dir("/tmp/report.pdf")              // "/tmp"
```

Use `filepath`, not `path`, and never build a path with string concatenation.
`path` is for URLs and always uses forward slashes; `filepath` adapts to the
operating system.

#gotcha[
  Never join a user-supplied name onto a directory without cleaning it:

  ```go
  // Directory traversal: "../../etc/passwd" escapes the directory.
  f, _ := os.Open(filepath.Join("uploads", userInput))
  ```

  `filepath.Join` does call `Clean`, which resolves the `..`, but the result
  can still point outside `uploads`. The check is explicit:

  ```go
  full := filepath.Join(base, filepath.Clean("/"+userInput))
  if !strings.HasPrefix(full, base+string(os.PathSeparator)) {
      return errors.New("invalid path")
  }
  ```

  Go 1.24 added `os.Root` for this, which opens a directory and refuses to
  escape it, and is the better answer where available.
]

== Temporary files and directories

```go
dir, err := os.MkdirTemp("", "myapp-*")
if err != nil {
    return err
}
defer os.RemoveAll(dir)

f, err := os.CreateTemp(dir, "upload-*.tmp")
```

In tests, `t.TempDir()` does this and cleans up automatically:

```go
func TestWrite(t *testing.T) {
    dir := t.TempDir()        // removed when the test ends, pass or fail
    path := filepath.Join(dir, "out.txt")
    ...
}
```

#project(title: "Project 19: a log analyser")[
  Bring together everything from Parts I and II: stream files that may be
  large, parse them, aggregate with maps, sort, and process several files
  concurrently.

  ```go
  package main

  import (
      "bufio"
      "flag"
      "fmt"
      "io"
      "os"
      "slices"
      "strconv"
      "strings"
      "sync"
  )

  // A line looks like:
  // 192.168.1.1 GET /api/users 200 134ms
  type stats struct {
      requests int
      byStatus map[int]int
      byPath   map[string]int
      totalMS  int
      slowest  string
      sloweste int
  }

  func newStats() *stats {
      return &stats{byStatus: map[int]int{}, byPath: map[string]int{}}
  }

  // analyse takes an io.Reader, not a filename, so it is testable with a
  // strings.Reader and works on stdin, a file or a network stream alike.
  func analyse(r io.Reader) (*stats, error) {
      s := newStats()
      sc := bufio.NewScanner(r)
      sc.Buffer(make([]byte, 0, 256*1024), 256*1024)

      for sc.Scan() {
          f := strings.Fields(sc.Text())
          if len(f) < 5 {
              // skip malformed lines rather than failing
              continue
          }
          status, err := strconv.Atoi(f[3])
          if err != nil {
              continue
          }
          ms, err := strconv.Atoi(strings.TrimSuffix(f[4], "ms"))
          if err != nil {
              continue
          }

          s.requests++
          s.byStatus[status]++
          s.byPath[f[2]]++
          s.totalMS += ms
          if ms > s.sloweste {
              s.sloweste, s.slowest = ms, f[2]
          }
      }
      // The check that is always skipped and always matters.
      if err := sc.Err(); err != nil {
          return nil, fmt.Errorf("scanning: %w", err)
      }
      return s, nil
  }

  func (s *stats) merge(o *stats) {
      s.requests += o.requests
      s.totalMS += o.totalMS
      for k, v := range o.byStatus {
          s.byStatus[k] += v
      }
      for k, v := range o.byPath {
          s.byPath[k] += v
      }
      if o.sloweste > s.sloweste {
          s.sloweste, s.slowest = o.sloweste, o.slowest
      }
  }

  func main() {
      top := flag.Int("top", 5, "how many paths to show")
      flag.Parse()

      paths := flag.Args()
      if len(paths) == 0 {
          // No files given: read standard input, so the tool pipes.
          s, err := analyse(os.Stdin)
          if err != nil {
              fail(err)
          }
          report(s, *top)
          return
      }

      // One goroutine per file, each with its own stats to avoid any
      // shared
      // state, merged at the end under a mutex.
      total := newStats()
      var mu sync.Mutex
      var wg sync.WaitGroup
      var failed int

      for _, p := range paths {
          wg.Go(func() {
              f, err := os.Open(p)
              if err != nil {
                  mu.Lock()
                  failed++
                  mu.Unlock()
                  fmt.Fprintf(os.Stderr, "skipping %s: %v\n", p, err)
                  return
              }
              defer f.Close()

              s, err := analyse(f)
              if err != nil {
                  fmt.Fprintf(os.Stderr, "reading %s: %v\n", p, err)
                  return
              }
              mu.Lock()
              total.merge(s)
              mu.Unlock()
          })
      }
      wg.Wait()

      report(total, *top)
      if failed > 0 {
          os.Exit(1)
      }
  }

  func report(s *stats, top int) {
      if s.requests == 0 {
          fmt.Println("no requests found")
          return
      }
      fmt.Printf("requests   %d\n", s.requests)
      fmt.Printf("average    %dms\n", s.totalMS/s.requests)
      fmt.Printf("slowest    %s (%dms)\n\n", s.slowest, s.sloweste)

      codes := make([]int, 0, len(s.byStatus))
      for c := range s.byStatus {
          codes = append(codes, c)
      }
      slices.Sort(codes)
      for _, c := range codes {
          fmt.Printf("  %d  %d\n", c, s.byStatus[c])
      }

      type pair struct {
          path string
          n    int
      }
      paths := make([]pair, 0, len(s.byPath))
      for p, n := range s.byPath {
          paths = append(paths, pair{p, n})
      }
      slices.SortFunc(paths, func(a, b pair) int { return b.n - a.n })

      fmt.Println()
      for _, p := range paths[:min(top, len(paths))] {
          fmt.Printf("  %-30s %d\n", p.path, p.n)
      }
  }

  func fail(err error) {
      fmt.Fprintln(os.Stderr, "loganalyse:", err)
      os.Exit(1)
  }
  ```

  Make a test file and run it:

  ```bash
  cat > access.log <<'EOF'
  1.1.1.1 GET /api/users 200 134ms
  1.1.1.2 GET /api/users 200 89ms
  1.1.1.3 POST /api/login 401 12ms
  1.1.1.4 GET /api/reports 200 2400ms
  EOF
  go run . access.log
  cat access.log | go run .
  ```

  Three design choices worth noticing. `analyse` takes an `io.Reader`, so the
  same function serves a file, standard input, and a `strings.Reader` in a
  test: write that test, it is four lines. Each goroutine builds its own
  `stats` and merges once at the end, so the lock is held briefly rather than
  on every line. And malformed lines are skipped rather than fatal, because a
  log analyser that dies on one bad line is useless against real logs.
]

#tryit[
  Write the test for `analyse` using `strings.NewReader`, covering a good line,
  a malformed line and an empty input. Then add a `-json` flag that emits the
  stats as JSON instead of text, which needs the capital-letter rule and struct
  tags from the last chapter. Finally, feed it a file with a 100KB single line
  and watch where the buffer limit bites if you remove the `sc.Buffer` call.
]
