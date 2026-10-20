#import "../template.typ": *

= Build your own Redis: RESP over raw TCP

Everything you have written so far has spoken HTTP, which is a protocol
somebody else parsed for you. This chapter goes one layer down: a TCP listener,
a wire format you implement yourself, and a server that the real `redis-cli`
can connect to.

That last part is what makes this the best protocol exercise there is. You are
not testing against your own client. You are testing against a program written
by people who have never heard of you, and when it works, your implementation
is correct in a way no unit test proves.

== The protocol, in one page

RESP is deliberately simple. The first byte says the type:

```
+OK\r\n                    simple string
-ERR unknown command\r\n   error
:42\r\n                    integer
$5\r\nhello\r\n            bulk string: length, then bytes
$-1\r\n                    null bulk string
*2\r\n$4\r\nECHO\r\n$2\r\nhi\r\n   array of two bulk strings
```

A client sends commands as *arrays of bulk strings*. `SET name Ada` arrives as:

```
*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$3\r\nAda\r\n
```

That is the whole format. Five types, one of which is a container.

#bridge[
  Compare this with JSON, which you parsed in chapter fifty. RESP carries the
  *length* of every string before its bytes, so the reader never scans for a
  delimiter and never has to escape anything.

  That is the difference between a format designed for humans and one designed
  for throughput, and it is why `$5\r\nhello\r\n` looks redundant and is not: a
  string containing `\r\n` needs no special handling at all.
]

== The server loop

```go
func Serve(ln net.Listener, db *Store) error {
    for {
        conn, err := ln.Accept()
        if err != nil {
            return err
        }
        go handle(conn, db)
    }
}

func handle(conn net.Conn, db *Store) {
    defer conn.Close()
    r := bufio.NewReader(conn)
    w := bufio.NewWriter(conn)

    for {
        args, err := readCommand(r)
        if err != nil {
            if err != io.EOF {
                writeError(w, "ERR "+err.Error())
                w.Flush()
            }
            return
        }
        execute(db, args, w)

        // Flush only when nothing more is buffered: that is what
        // makes pipelining fast rather than chatty.
        if r.Buffered() == 0 {
            w.Flush()
        }
    }
}
```

Goroutine per connection, which in Go is the obvious design and in most
languages is not. Ten thousand connections is ten thousand goroutines and a few
hundred megabytes, which chapter thirteen promised and this is where you cash
it in.

#gotcha[
  The flush condition is the single most important line for performance.

  Flushing after every command turns a pipeline of a thousand commands into a
  thousand small writes and a thousand round trips. Flushing only when the read
  buffer is empty lets a client that sent a thousand commands in one packet
  receive one packet back.

  Measured with `redis-benchmark -P 100`, the difference is roughly an order of
  magnitude.
]

== Reading a command

```go
func readCommand(r *bufio.Reader) ([]string, error) {
    line, err := readLine(r)
    if err != nil {
        return nil, err
    }
    if len(line) == 0 {
        return nil, errors.New("protocol error: empty line")
    }
    if line[0] != '*' {
        // Inline command, what you get from telnet.
        return strings.Fields(string(line)), nil
    }

    n, err := strconv.Atoi(string(line[1:]))
    if err != nil || n < 0 {
        return nil, errors.New("protocol error: bad multibulk")
    }
    args := make([]string, 0, n)
    for i := 0; i < n; i++ {
        s, err := readBulk(r)
        if err != nil {
            return nil, err
        }
        args = append(args, s)
    }
    return args, nil
}

func readBulk(r *bufio.Reader) (string, error) {
    line, err := readLine(r)
    if err != nil {
        return "", err
    }
    if len(line) == 0 || line[0] != '$' {
        return "", errors.New("protocol error: expected $")
    }
    n, err := strconv.Atoi(string(line[1:]))
    if err != nil || n < 0 || n > maxBulk {
        return "", errors.New("protocol error: bad bulk length")
    }
    buf := make([]byte, n+2)          // payload plus CRLF
    if _, err := io.ReadFull(r, buf); err != nil {
        return "", err
    }
    return string(buf[:n]), nil
}
```

`io.ReadFull` is the heart of it, and it is the line that separates a working
server from a mysteriously broken one.

#gotcha[
  *TCP is a stream, not a sequence of messages.* One `Read` can return half a
  command, or two commands, or one and a half. Nothing preserves the boundaries
  the sender used.

  So a bulk string of five bytes is read with `io.ReadFull`, which loops until
  it has five, rather than with one `Read` that might return three. Every
  protocol bug that only appears over a real network and never over localhost
  is this, because localhost rarely splits packets and the internet does it
  constantly.
]

Note `n > maxBulk` as well. A client claiming a 4 GB bulk string must not cause
a 4 GB allocation, which is chapter forty two's bomb wearing a different hat.

== The store, and the lock

```go
type Store struct {
    mu   sync.RWMutex
    data map[string]entry
}

type entry struct {
    value   string
    expires time.Time         // zero means never
}

func (s *Store) Get(key string) (string, bool) {
    s.mu.RLock()
    e, ok := s.data[key]
    s.mu.RUnlock()

    if !ok {
        return "", false
    }
    if !e.expires.IsZero() && time.Now().After(e.expires) {
        s.Delete(key)              // lazy expiry
        return "", false
    }
    return e.value, true
}
```

Expiry is *lazy plus active*, which is what Redis itself does: a key is dropped
when somebody touches it after its deadline, and a background sweep samples
random keys to catch the ones nobody asks for. Without the sweep, keys nobody
reads live forever and the memory never comes back.

#gotcha[
  The real Redis is single threaded for command execution, so it needs no lock
  at all. Your Go version has a goroutine per connection and therefore does.

  That is a genuine trade: the lock costs you contention on a hot key, and
  Redis's single thread costs it one core. Before deciding yours is better,
  benchmark with `redis-benchmark` and look at the numbers, because a global
  `RWMutex` under write-heavy load is less impressive than it sounds. Sharding
  the map by key hash is the usual next step.
]

== Commands, as a table

```go
type command struct {
    arity   int               // -N means "at least N"
    handler func(*Store, []string, *bufio.Writer)
}

var commands = map[string]command{
    "ping":   {1, cmdPing},
    "echo":   {2, cmdEcho},
    "set":    {-3, cmdSet},
    "get":    {2, cmdGet},
    "del":    {-2, cmdDel},
    "exists": {-2, cmdExists},
    "incr":   {2, cmdIncr},
    "expire": {3, cmdExpire},
    "ttl":    {2, cmdTTL},
    "keys":   {2, cmdKeys},
    "command": {-1, cmdCommand},
}
```

Arity checked once, in `execute`, rather than in every handler. And `command`
is in the list for a practical reason: `redis-cli` sends `COMMAND DOCS` on
connect, and a server that errors on it looks broken before you have typed
anything.

== Incrementing, which is where atomicity shows up

```go
func cmdIncr(s *Store, args []string, w *bufio.Writer) {
    s.mu.Lock()
    defer s.mu.Unlock()

    e := s.data[args[1]]
    n, err := strconv.ParseInt(e.value, 10, 64)
    if e.value != "" && err != nil {
        writeError(w, "ERR value is not an integer or out of range")
        return
    }
    n++
    e.value = strconv.FormatInt(n, 10)
    s.data[args[1]] = e
    writeInteger(w, n)
}
```

The read, the increment and the write are inside *one* lock. Doing `Get` then
`Set` from the handler would be two operations with a window between them, and
two clients incrementing at once would lose a count. This is chapter fifteen's
lesson in its smallest possible form, and it is exactly why `INCR` exists as a
command rather than being left to the client.

== Persistence, which you have already built

An append-only file: every write command is appended to a log, and on startup
the log is replayed. That is chapter fifty one's write-ahead log, unchanged.

The questions are the same too. `fsync` on every write is durable and slow;
once a second is fast and loses up to a second. Redis makes that a
configuration option, and so should you, with the trade written down where the
person choosing can see it.

== Testing against the real thing

```bash
go run . &                      # your server on :6380

redis-cli -p 6380 ping
redis-cli -p 6380 set name Ada
redis-cli -p 6380 get name
redis-cli -p 6380 incr counter

redis-benchmark -p 6380 -t set,get -n 100000 -q
redis-benchmark -p 6380 -t set,get -n 100000 -P 100 -q
```

The benchmark with and without pipelining is the measurement that tells you
whether your flush logic is right.

#bridge[
  This is the most satisfying validation in the book. `redis-cli` knows nothing
  about your code, implements the protocol strictly, and will tell you
  immediately when a reply type is wrong.

  Interoperating with somebody else's client is a different standard of
  correctness from passing your own tests, and it is the standard that matters
  when you write a protocol anybody else will speak.
]

== What you are not building

Be clear about the gap, because "I built Redis" is a sentence that needs
qualifying:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Data types], [Real Redis has lists, sets, sorted sets, hashes, streams.
    This has strings],
  [Pub/sub], [A different connection mode, where the server pushes],
  [Scripting], [Lua, executed atomically on the server],
  [Replication], [A whole protocol of its own],
  [Cluster], [Key hashing, slots, redirects],
)

Strings plus expiry is perhaps five percent of Redis and a hundred percent of
the protocol work. That is the right trade for a chapter, and the remaining
list is a very good source of weekend projects.

#project(title: "Project 49: a server redis-cli cannot tell apart")[
  In `part6-from-scratch/53-redis`, build it.

  *One: the protocol.* Reader and writer for all five RESP types, with
  `io.ReadFull` for bulk payloads and a cap on declared lengths.

  *Two: the commands.* `PING`, `ECHO`, `SET` with `EX`, `GET`, `DEL`,
  `EXISTS`, `INCR`, `DECR`, `EXPIRE`, `TTL`, `KEYS`, `COMMAND`.

  *Three: expiry.* Lazy on read plus an active sampler. Prove with a test that
  a key nobody reads is eventually reclaimed.

  *Four: the real client.* Every command above verified through `redis-cli`,
  pasted into `SESSION.md`.

  *Five: the benchmark.* `redis-benchmark` with and without `-P 100`, both
  recorded. If pipelining does not help by a large factor, your flush logic is
  wrong and that is the finding.

  *Six: concurrency.* Run the benchmark with 50 clients under `-race`. Then
  shard the map into sixteen and measure again.

  *Seven: persistence.* An append-only log with a configurable sync policy.
  Kill the server mid-load, restart, and report what survived under each
  policy.

  *Eight: one data type.* Add lists with `LPUSH`, `RPUSH`, `LRANGE` and `LLEN`.
  Then look up how Redis actually stores a list and write a paragraph on why a
  Go slice is not what it uses.
]

#tryit[
  Connect to your server with `telnet localhost 6380` and type `PING` followed
  by Enter.

  That is the inline command path, and seeing a raw protocol answer a human
  typing into a socket is the moment networking stops being abstract. Then try
  `SET a b` and `GET a` the same way.
]

#assignment(name: "RESP Server - Assignment One")[
  Build the server and prove it against tools you did not write.

  *One.* The full implementation with the command table, expiry and
  persistence.

  *Two.* `SESSION.md`: a real `redis-cli` session exercising every command,
  including the error cases (wrong arity, `INCR` on a non-number).

  *Three.* `BENCH.md` with four numbers: get and set, pipelined and not. Plus
  the same four from a real Redis on your machine if you have one, and one
  sentence on the difference.

  *Four.* The framing test: a client that sends a command *one byte at a time*
  with a delay between bytes. Your server must handle it. If it does not, you
  found the bug this chapter warned about.

  *Five.* One paragraph on what you would have to build next to run this in
  production, in the order you would build it.
]
