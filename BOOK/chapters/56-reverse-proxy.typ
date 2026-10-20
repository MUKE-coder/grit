#import "../template.typ": *

= Build your own reverse proxy and load balancer

Every production system has one of these in front of it, and most developers
treat it as a configuration file somebody else wrote. A proxy is a few hundred
lines, and writing one explains a list of things that are otherwise folklore:
why `X-Forwarded-For` can lie, why some headers must be stripped, why a retry
can corrupt data, and why your streaming endpoint buffers when it goes through
nginx.

== What a reverse proxy does

```
client  →  proxy  →  one of several backends
                 ←  response
```

Four jobs, and only the first is obvious:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Distribute], [Spread requests over healthy backends],
  [Protect], [Timeouts, limits, and TLS termination in one place],
  [Hide], [One public address for many internal services],
  [Survive], [Take a backend out when it fails, put it back when it recovers],
)

== The core, which is shorter than you expect

```go
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    backend, release, ok := p.pool.Pick(r)
    if !ok {
        http.Error(w, "no healthy backend", 503)
        return
    }
    defer release()

    outreq := r.Clone(r.Context())
    outreq.URL.Scheme = backend.URL.Scheme
    outreq.URL.Host = backend.URL.Host
    outreq.RequestURI = ""               // must be empty for a client
    stripHopByHop(outreq.Header)
    addForwarded(outreq, r)

    resp, err := p.transport.RoundTrip(outreq)
    if err != nil {
        p.pool.Fail(backend)
        http.Error(w, "upstream unavailable", 502)
        return
    }
    defer resp.Body.Close()

    stripHopByHop(resp.Header)
    copyHeader(w.Header(), resp.Header)
    w.WriteHeader(resp.StatusCode)
    p.copyBody(w, resp.Body)
}
```

Everything after this is the detail, and the detail is where proxies go wrong.

== Hop-by-hop headers, and why they must be removed

Some headers describe *this connection*, not the message. Forwarding them
onward is wrong because the next connection is a different connection:

```go
var hopByHop = []string{
    "Connection", "Keep-Alive", "Proxy-Authenticate",
    "Proxy-Authorization", "TE", "Trailer",
    "Transfer-Encoding", "Upgrade",
}
```

And the subtle part: `Connection` can *name* other headers that are also
hop-by-hop for this connection, so those must be stripped too.

#gotcha[
  Forwarding `Transfer-Encoding: chunked` to a backend that also sets a
  `Content-Length` is the smuggling setup from chapter fifty four, built by
  your own proxy.

  Strip, then let your HTTP client decide the framing for the connection it is
  making. Never copy framing headers across a hop.
]

== X-Forwarded-For, and the header that lies

```go
func addForwarded(out, in *http.Request) {
    ip, _, _ := net.SplitHostPort(in.RemoteAddr)

    prior := in.Header.Get("X-Forwarded-For")
    if prior != "" {
        ip = prior + ", " + ip
    }
    out.Header.Set("X-Forwarded-For", ip)
    out.Header.Set("X-Forwarded-Proto", scheme(in))
    out.Header.Set("X-Forwarded-Host", in.Host)
}
```

The chain grows as the request passes each proxy, so the *leftmost* entry is
the original client and each later one was added by a hop.

#gotcha[
  *A client can send `X-Forwarded-For` itself.* Nothing stops it. So the
  leftmost value is not the client's IP; it is whatever the client claimed,
  followed by the truth.

  The consequence is real: an application that rate limits, geolocates or
  audits by the leftmost entry can be fooled by anybody who sends a header.
  Logs attribute actions to an invented address.

  The rule: *count from the right, and trust exactly as many entries as you
  have proxies you control*. One proxy means the real client IP is the last
  entry your proxy appended, not the first in the list. Configure the count
  explicitly; do not infer it.
]

== Choosing a backend

Round robin is two lines and is the right default:

```go
func (p *Pool) roundRobin() *Backend {
    n := atomic.AddUint64(&p.next, 1)
    healthy := p.healthy()
    if len(healthy) == 0 {
        return nil
    }
    return healthy[n%uint64(len(healthy))]
}
```

Least connections is better when request durations vary a lot, which is most
real traffic:

```go
func (p *Pool) leastConns() *Backend {
    var best *Backend
    for _, b := range p.healthy() {
        if best == nil ||
            b.inflight.Load() < best.inflight.Load() {
            best = b
        }
    }
    return best
}
```

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Round robin], [Even request counts. Fine when requests cost the same],
  [Least connections], [Even *load*. Better with mixed fast and slow
    endpoints],
  [Weighted], [Backends of different sizes],
  [Consistent hash], [The same key goes to the same backend, for cache
    locality or sticky sessions],
)

Round robin's failure mode is worth knowing: one slow backend keeps receiving
its share, so a quarter of your users wait on the sick machine while three
healthy ones idle. Least connections drains away from it automatically, because
its in-flight count stays high.

== Health, actively and passively

*Active*: ask each backend, on a timer.

```go
func (p *Pool) probe(b *Backend) {
    ctx, cancel := context.WithTimeout(
        context.Background(), 2*time.Second)
    defer cancel()

    req, _ := http.NewRequestWithContext(ctx, "GET",
        b.URL.String()+"/health", nil)
    resp, err := p.probeClient.Do(req)

    switch {
    case err != nil || resp.StatusCode >= 500:
        b.markUnhealthy()
    default:
        resp.Body.Close()
        b.markHealthy()
    }
}
```

*Passive*: notice real traffic failing. A backend that refuses three
connections in a row is out, regardless of what its health endpoint says.

Recovery must be careful. A backend that has been out comes back *half open*:
it gets one request, and only if that succeeds does it rejoin.

#gotcha[
  Without the half-open step, a backend that is still broken rejoins the pool
  on a timer, takes its full share, fails them all, and leaves again. The pool
  oscillates and a steady fraction of users get errors for as long as the
  backend is sick.

  This is the circuit breaker pattern, and the half-open state is the part
  people leave out.
]

Thresholds matter too: two consecutive failures to eject, three consecutive
successes to restore. One failure is a blip, and treating it as an outage makes
the proxy less reliable than no proxy.

== Retries, and the one that corrupts data

Retrying a failed request sounds obviously good. It is conditional:

```go
func retryable(r *http.Request, err error, sentBytes bool) bool {
    if sentBytes {
        return false        // the backend may have acted on it
    }
    switch r.Method {
    case "GET", "HEAD", "OPTIONS", "PUT", "DELETE":
        return isConnectionError(err)
    default:
        return false        // POST and PATCH are not idempotent
    }
}
```

#gotcha[
  *Never retry a POST that was fully sent.* A connection that failed after the
  request reached the backend may have failed on the way back, and the work is
  already done. Retrying charges the card twice.

  A connection *refused* is safe to retry: nothing was delivered. A connection
  *reset mid-response* is not.

  And cap the total: a retry budget of ten percent of traffic. Without a cap, a
  failing backend turns into three times the traffic at the exact moment the
  system is least able to take it, which is how a small outage becomes a
  retry storm.
]

== Streaming, and the buffering that breaks it

```go
func (p *Proxy) copyBody(w http.ResponseWriter, body io.Reader) {
    if flusher, ok := w.(http.Flusher); ok && p.flushInterval > 0 {
        // Server-sent events, long polling, log tails: the client
        // needs each write as it happens, not when the buffer is
        // full or the response ends.
        copyAndFlush(w, body, flusher, p.flushInterval)
        return
    }
    io.Copy(w, body)
}
```

This is why your streaming endpoint works locally and arrives all at once
through a proxy: the default is to buffer, and a stream needs the opposite.
Detect it by `Content-Type: text/event-stream`, or by the absence of a
`Content-Length`, and flush.

WebSockets need the connection itself, not the body: on a `101 Switching
Protocols`, hijack the client connection and copy bytes in both directions
until either side closes.

== Draining a backend

Deployment is where a proxy earns its place. Taking a backend out *gracefully*
means: stop sending it new requests, let the in-flight ones finish, then stop
the process.

```go
func (b *Backend) Drain(ctx context.Context) error {
    b.accepting.Store(false)
    for b.inflight.Load() > 0 {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(50 * time.Millisecond):
        }
    }
    return nil
}
```

That loop, plus chapter twenty one's graceful shutdown in the backend, is a
zero-downtime deploy. There is nothing else to it.

#project(title: "Project 52: the proxy, measured under failure")[
  In `part6-from-scratch/56-proxy`, build it and then break things on
  purpose.

  *One: the proxy.* Header stripping, forwarded headers, round robin, and a
  configurable backend list.

  *Two: three backends.* Write a tiny backend that reports its own name and
  can be told to be slow or to fail:

  ```
  ./backend -name a -port 9001
  ./backend -name b -port 9002 -latency 500ms
  ./backend -name c -port 9003 -fail-rate 0.5
  ```

  *Three: the strategies.* Round robin and least connections, both
  implemented. Drive 10,000 requests through each with `hey` or `k6` and
  record p50, p95 and p99 for both. The difference with backend `b` slow is the
  whole lesson.

  *Four: health.* Active probes and passive ejection with thresholds, plus the
  half-open recovery. Kill backend `a` mid-load and record how many requests
  failed before it was ejected. Then restart it and time the recovery.

  *Five: the oscillation.* Remove the half-open step, make a backend fail its
  probe but fail real traffic, and watch the error rate pulse. Put it back.

  *Six: streaming.* Add an SSE endpoint to a backend and prove the events
  arrive one at a time through the proxy. Then set the flush interval to zero
  and watch them arrive in a lump.

  *Seven: the forged header.* Send a request with your own `X-Forwarded-For`
  and show what your backend logs. Then implement the trusted-hop count and
  show it logging the truth.

  *Eight: the drain.* Deploy a new version of a backend with zero failed
  requests under load. That is the proof.
]

#tryit[
  Put your proxy in front of the Grit API from Part V, with two copies of the
  API running.

  Then stop one mid-session in the admin panel. If you did the health checks
  and draining properly, the person using it sees nothing at all, which is the
  most satisfying negative result in this book.
]

#assignment(name: "Reverse Proxy - Assignment One")[
  Build the proxy and report its behaviour under failure, with numbers.

  *One.* The proxy with both balancing strategies, active and passive health
  checks, half-open recovery, conditional retries with a budget, and streaming
  support.

  *Two.* `BALANCING.md`: p50, p95 and p99 for both strategies, with all
  backends healthy and with one slow, as four columns of a table.

  *Three.* `FAILURE.md`: requests lost when a backend dies, before and after
  passive ejection; the oscillation experiment with and without half-open; and
  the drain with zero errors.

  *Four.* The forged `X-Forwarded-For`, before and after the trusted-hop fix,
  with both log lines.

  *Five.* One paragraph: which of nginx, Caddy, HAProxy or an application load
  balancer would you use instead of this, for what, and which one feature of
  yours would you miss.
]
