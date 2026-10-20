#import "../template.typ": *

= Build your own HTTP server, from a TCP socket

Chapter twenty one used `net/http`, which hides the protocol so completely that
most web developers have never seen a request line. This chapter removes the
hiding: a TCP listener, HTTP/1.1 parsed by hand, keep-alive, chunked bodies,
and the timeouts without which a server can be held open by one slow client.

The payoff is not a better server. It is that the next time you read about
request smuggling, a reverse proxy misbehaving or a mysterious hanging
connection, you will know exactly which bytes are involved.

== What a request actually is

```
POST /invoices?draft=1 HTTP/1.1\r\n
Host: api.example.com\r\n
Content-Type: application/json\r\n
Content-Length: 27\r\n
\r\n
{"client":1,"total":2500}
```

A request line, headers, a blank line, and optionally a body. The response has
the same shape with a status line instead:

```
HTTP/1.1 201 Created\r\n
Content-Type: application/json\r\n
Content-Length: 42\r\n
\r\n
{"data":{"id":7},"message":"created"}
```

That is HTTP/1.1. Everything else is rules about how to know where the body
ends and whether the connection stays open.

#bridge[
  This is why `curl -v` is such a good debugging tool: it prints exactly these
  lines. Once you have written the parser, `curl -v` stops being noise and
  becomes the clearest view of what your framework is actually sending.
]

== The accept loop, with the timeouts that matter

```go
func (s *Server) Serve(ln net.Listener) error {
    for {
        conn, err := ln.Accept()
        if err != nil {
            return err
        }
        go s.handleConn(conn)
    }
}

func (s *Server) handleConn(conn net.Conn) {
    defer conn.Close()
    r := bufio.NewReaderSize(conn, 8<<10)
    w := bufio.NewWriter(conn)

    for {
        conn.SetReadDeadline(time.Now().Add(s.ReadHeaderTimeout))

        req, err := readRequest(r)
        if err != nil {
            writeError(w, err)
            w.Flush()
            return
        }

        conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
        conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))

        keepAlive := s.serve(w, req)
        w.Flush()

        if !keepAlive {
            return
        }
        conn.SetReadDeadline(time.Now().Add(s.IdleTimeout))
    }
}
```

#gotcha[
  `ReadHeaderTimeout` is separate from `ReadTimeout`, and the reason has a
  name: *Slowloris*.

  The attack is to open a thousand connections and send one header byte every
  twenty seconds. Each connection is technically active, so a single overall
  timeout never fires, and the server runs out of connections while using
  almost no bandwidth.

  A deadline on the *headers* specifically is the defence, and it is why
  `net/http` grew that field. A server with no header deadline is held open by
  a laptop.
]

== Parsing the request

```go
func readRequest(r *bufio.Reader) (*Request, error) {
    line, err := readLine(r, maxLine)
    if err != nil {
        return nil, err
    }
    method, target, version, ok := splitRequestLine(line)
    if !ok {
        return nil, badRequest("malformed request line")
    }
    if version != "HTTP/1.1" && version != "HTTP/1.0" {
        return nil, &httpError{505, "version not supported"}
    }

    headers, err := readHeaders(r)
    if err != nil {
        return nil, err
    }
    body, err := bodyReader(r, headers)
    if err != nil {
        return nil, err
    }
    return &Request{
        Method: method, Target: target,
        Version: version, Headers: headers, Body: body,
    }, nil
}
```

Headers are a loop until the blank line, with two defences:

```go
func readHeaders(r *bufio.Reader) (Header, error) {
    h := Header{}
    for i := 0; ; i++ {
        if i > maxHeaders {            // 100 is plenty
            return nil, badRequest("too many headers")
        }
        line, err := readLine(r, maxLine)   // 8 KB each
        if err != nil {
            return nil, err
        }
        if len(line) == 0 {
            return h, nil
        }
        name, value, ok := strings.Cut(line, ":")
        if !ok {
            return nil, badRequest("malformed header")
        }
        h.Add(textproto.CanonicalMIMEHeaderKey(name),
            strings.TrimSpace(value))
    }
}
```

Header names are case insensitive, so they are canonicalised once on the way
in. Doing it at every lookup instead is both slower and a source of bugs where
`Content-Type` and `content-type` are two different entries.

The limits are not politeness. Without them, a client sending headers forever
is an out-of-memory error with no request ever completing.

== Where does the body end

This is the whole of HTTP/1.1's complexity, and there are exactly three
answers:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`Content-Length: n`], [Read exactly n bytes],
  [`Transfer-Encoding: chunked`], [Read chunks until a zero-length one],
  [Neither], [There is no body. A GET, usually],
)

Chunked is the one worth implementing, because it is how a server streams
something whose size it does not know yet:

```
7\r\n
Mozilla\r\n
9\r\n
Developer\r\n
0\r\n
\r\n
```

Each chunk is a hexadecimal length, CRLF, that many bytes, CRLF. A zero-length
chunk ends the body.

```go
func bodyReader(r *bufio.Reader, h Header) (io.Reader, error) {
    te := strings.ToLower(h.Get("Transfer-Encoding"))
    cl := h.Get("Content-Length")

    if te == "chunked" && cl != "" {
        // Both present: refuse. See below.
        return nil, badRequest("both Transfer-Encoding and " +
            "Content-Length")
    }
    switch {
    case te == "chunked":
        return newChunkedReader(r), nil
    case cl != "":
        n, err := strconv.ParseInt(cl, 10, 64)
        if err != nil || n < 0 {
            return nil, badRequest("bad Content-Length")
        }
        return io.LimitReader(r, n), nil
    default:
        return emptyReader{}, nil
    }
}
```

#gotcha[
  That refusal is not pedantry. A message with *both* framing headers is the
  basis of *HTTP request smuggling*, and it is one of the highest-impact web
  vulnerabilities there is.

  The attack works when two servers disagree: a front-end proxy uses
  `Content-Length` and your back-end uses `Transfer-Encoding`, or the reverse.
  The bytes one of them considers to be the end of request A, the other
  considers the start of request B. An attacker appends a request that your
  server attributes to the *next* user's connection.

  The defence is to refuse the ambiguity rather than resolve it. If both
  headers are present, the message is malformed: 400, and close the
  connection.
]

== Keep-alive, and the body you forgot to read

HTTP/1.1 connections are persistent by default. After responding, you loop and
read the next request from the same socket.

```go
func keepAlive(req *Request) bool {
    conn := strings.ToLower(req.Headers.Get("Connection"))
    if req.Version == "HTTP/1.0" {
        return conn == "keep-alive"
    }
    return conn != "close"
}
```

And the trap:

```go
// Before reading the next request, drain whatever the handler
// did not read of this one. Otherwise the leftover body bytes
// are parsed as the next request line, which produces a
// baffling "malformed request line: {"client":1..."
io.Copy(io.Discard, req.Body)
```

#gotcha[
  This bug is unforgettable once seen. A handler returns 400 early without
  reading the body, the next loop iteration reads the unread JSON as a request
  line, and the error message contains fragments of the previous request.

  Worse, it is intermittent: it only happens when a handler returns early, so
  it shows up under error conditions, which is when you are least able to
  reason about it.

  Drain or close. `net/http` drains a limited amount and closes beyond it,
  which is a good policy to copy.
]

== Writing the response

```go
func (w *responseWriter) WriteHeader(status int) {
    fmt.Fprintf(w.buf, "HTTP/1.1 %d %s\r\n",
        status, statusText(status))
    w.header.Set("Date", time.Now().UTC().
        Format(http.TimeFormat))

    if _, ok := w.header["Content-Length"]; !ok {
        w.header.Set("Transfer-Encoding", "chunked")
        w.chunked = true
    }
    for name, values := range w.header {
        for _, v := range values {
            fmt.Fprintf(w.buf, "%s: %s\r\n", name, v)
        }
    }
    w.buf.WriteString("\r\n")
}
```

If the handler knows the length, send `Content-Length`. If it does not, switch
to chunked. That is the decision `net/http` makes for you: it buffers a little,
and if the handler finishes within the buffer it sets a length, otherwise it
goes chunked.

Two responses have no body at all and must not have a length: 204 and 304. A
server that sends `Content-Length: 0` with a 304 confuses caches in ways that
take a day to find.

== 100-continue, the handshake nobody knows about

```
Expect: 100-continue
```

A client about to upload something large asks first. The server answers either:

```
HTTP/1.1 100 Continue\r\n\r\n
```

and the client sends the body, or it answers 413 or 401 immediately and the
client never uploads a gigabyte that was going to be rejected anyway.

`curl` sends this automatically for large bodies. A server that ignores it
makes those clients wait for a timeout before uploading, which presents as
"uploads are slow from curl and fine from the browser".

== Routing, which you have already written

The router from chapter forty eight drops straight in: match the method and the
path segments, collect parameters, call the handler. Nothing about routing
changes because you parsed the request yourself.

== What you are not building

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [TLS], [`crypto/tls` wraps your listener in about ten lines. Not hard,
    just not HTTP],
  [HTTP/2], [Binary framing, HPACK header compression, multiplexed streams.
    A different project of a similar size],
  [HTTP/3], [QUIC, over UDP. A much bigger project],
  [Compression], [`Accept-Encoding` and a gzip writer. A good exercise],
)

== Testing it

Against real clients, as in chapter fifty three:

```bash
curl -v http://localhost:8080/health
curl -v -X POST localhost:8080/echo -d '{"a":1}'
curl -v --http1.1 -H "Transfer-Encoding: chunked" \
  -T big.json localhost:8080/upload
```

And against raw bytes, which is where the interesting failures live:

```bash
printf 'GET / HTTP/1.1\r\nHost: x\r\n\r\n' | nc localhost 8080

# Two requests in one packet: pipelining
R='GET /a HTTP/1.1\r\nHost: x\r\n\r\n'
printf "$R$R" | nc localhost 8080
```

#project(title: "Project 50: the server, and the attacks it survives")[
  In `part6-from-scratch/54-http`, build it.

  *One: the parser.* Request line, headers with limits and canonical names,
  and all three body framings including chunked.

  *Two: the writer.* Status line, `Date`, `Content-Length` or chunked, and
  correct handling of 204 and 304.

  *Three: keep-alive.* Persistent connections, the drain, and `Connection:
  close` honoured in both directions.

  *Four: the timeouts.* Separate header, read, write and idle deadlines, each
  configurable.

  *Five: Slowloris.* Write a client that opens 200 connections and sends one
  header byte every two seconds. With the header timeout off, show the server
  stop accepting. Turn it on and show it survive. Record both.

  *Six: smuggling.* Send a request with both `Content-Length` and
  `Transfer-Encoding` and confirm a 400. Then read the two smuggling variants
  and write a paragraph on which one your refusal prevents.

  *Seven: pipelining.* Two requests in one packet, two responses in order.

  *Eight: interoperability.* Point a real HTTP client at it, Go's own
  `http.Client` and `curl`, and make both work. Then serve a file and open it
  in a browser.
]

#tryit[
  Run your server and the standard library's side by side, and compare the
  exact bytes with `curl -v` for the same handler.

  Every difference is a decision `net/http` made that you have not. Some will
  be headers you forgot; some will be behaviour you now understand well enough
  to disagree with, which is the point of the chapter.
]

#assignment(name: "HTTP Server - Assignment One")[
  Build the server, then attack it.

  *One.* The full implementation with chunked bodies, keep-alive, pipelining
  and the four timeouts.

  *Two.* `ATTACKS.md`: Slowloris with and without the header timeout and the
  connection counts; the smuggling refusal with the raw bytes you sent; and one
  oversized-header attempt.

  *Three.* The drain bug, demonstrated deliberately: remove the drain, send a
  POST whose handler returns early, and paste the baffling error from the next
  request on that connection.

  *Four.* Interoperability evidence: `curl -v` output, a Go `http.Client`
  test, and a screenshot of a browser loading a page from it.

  *Five.* One paragraph on what you would need before putting this on the
  internet, honestly ordered. TLS is on the list, but it is not first.
]
