---
title: "Your own post is missing: read-your-writes in a cookie"
subtitle: "You add a read replica, latency drops, and one bug appears that nobody can reproduce: right after someone posts, their own post is gone from the list. The textbook fix puts Redis on every read. There is a smaller answer that costs nothing and cannot be stale, because it travels with the person who wrote."
series: "The Daily Grit"
edition: 22
date: 2026-10-19
readingTime: "12 min"
author: "Muke JohnBaptist"
tags: [grit, go, replication, consistency, scaling, postgres, caching]
canonical: "https://gritframework.dev/blog/your-own-post-is-missing"
---

You add a read replica. Reads go to the replica, writes go to the primary, and the primary stops being the bottleneck. Everything is faster and nothing breaks.

Except one thing, and it is the single most reported bug of that change: **somebody posts, the page reloads, and their own post is not there.** So they post it again.

It is maddening to debug because it does not reproduce. You post, you reload, it is there. The replica caught up in the time it took you to move your hand to the keyboard. It happens to a fraction of users, for a fraction of a second, and only to the person who wrote. Everyone else sees exactly what they expect, which is the list as it was a moment ago, and a moment ago is fine for them.

## Why this is not a bug you can fix with care

The sequence is just what replication is:

```
t=0ms    POST /api/posts      → primary, committed
t=5ms    GET  /api/posts      → replica, has not received it yet
t=40ms   replication catches up
```

There is no mistake in that trace. The write went to the right place, the read went to the right place, and the replica is doing exactly what an asynchronous replica does. The data is not lost; it is in flight.

What is wrong is the expectation. Eventual consistency is fine for reading other people's data, where "a moment ago" is indistinguishable from "now". It is not fine for reading your own, because you have a memory of what you just did, and a system that contradicts it looks broken.

The property you want has a name: **read-your-own-writes consistency**. Not global consistency, which is expensive and which you do not need. Just the guarantee that one person's reads reflect that person's writes.

## The textbook answer, and its bill

The standard solution is a sticky-write store. When a user writes, record "user X wrote at T" in Redis. On every read, look up that user, and if they wrote within the replication window, read from the primary.

It works. It also means:

- **A Redis round trip on every read.** Not on the writes, which are rare. On the reads, which are everything. You have added a network hop to the hot path of your entire application to fix a one-second edge case.
- **Redis is now on the critical path for reading.** If Redis is down, every read either fails or silently loses the guarantee. You have taken a cache, which is supposed to be optional, and made it load-bearing.
- **It is one more thing to run**, size, monitor and pay for.

There is a reflex in systems work that reaches for shared state whenever two requests need to agree about something. Sometimes that is right. It is worth asking first whether the two requests have anything in common that could carry the fact for you.

Here they do. **They are the same browser.**

## Put the fact on the person

The whole mechanism, from `internal/middleware/read_your_writes.go`:

```go
const (
    readYourWritesCookie = "grit_rw"
    readYourWritesWindow = 5 * time.Second
)

func ReadYourWrites() gin.HandlerFunc {
    return func(c *gin.Context) {
        if raw, err := c.Cookie(readYourWritesCookie); err == nil {
            if ms, err := strconv.ParseInt(raw, 10, 64); err == nil {
                if time.Since(time.UnixMilli(ms)) < readYourWritesWindow {
                    c.Set("read_your_writes", true)
                }
            }
        }

        c.Next()

        switch c.Request.Method {
        case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
            if c.Writer.Status() < 300 {
                // set grit_rw to now
            }
        }
    }
}
```

A successful write sets a cookie holding the current timestamp. A read that arrives with a recent timestamp in that cookie is flagged, and the flag tells the database layer to use the primary for this one request.

That is it. No store, no lookup, no new dependency.

The properties fall out of where the fact lives:

**It cannot be stale.** The cookie travels with the browser that did the write. There is no replication delay on the thing that measures replication delay, which is the circularity a shared store has to work around.

**It costs nothing on the read path.** Reading a cookie that is already on the request is free. Compare that with a Redis round trip per read.

**It fails in the right direction.** No cookie, expired cookie, nonsense cookie: you read from the replica, which is the normal path. The failure mode of the mechanism is the behaviour you had before the mechanism.

**It is per-browser, which is per-person in practice.** Exactly the scope the guarantee needs.

## Three details that are each a bug if you get them wrong

### Only on a write that worked

```go
if c.Writer.Status() < 300 {
```

A `POST` that returned 422 changed nothing. Setting the cookie for it pins that user's next five seconds of reads to the primary to protect a write that never happened. On a form people fill in badly a few times before getting it right, those are the reads you least want on the primary.

### Check the method, not the route

A `GET` never sets the cookie, however it is spelled. Routes with names like `/items/:id/archive` that are secretly writes are a different problem, and the honest fix is to make them a `POST`. Deciding by HTTP method keeps the rule mechanical.

### Five seconds, and why it is not fifty

```go
// Long enough to cover normal replication, short enough that a user who
// writes once does not pin every later read to the primary. Replication
// under a second is typical; five is the number to raise if your lag
// monitor says otherwise.
readYourWritesWindow = 5 * time.Second
```

This is the only real dial, and it is a direct trade.

Too short and the bug comes back for the tail of your replication distribution. Too long and every user who writes once sends all their reads to the primary for that long, which on a write-heavy application is most of your traffic, and you have quietly undone the reason you added the replica.

Five seconds is roughly five to ten times typical sub-second replication lag, which covers the tail without being generous. The number to tune it against is your actual lag metric, not a guess, which is why the comment says so.

## It is not a security boundary, and does not need to be

A cookie the client can forge decides which database a read goes to. That sounds alarming for about ten seconds, until you ask what forging it achieves.

> The worst a forged cookie achieves is reading from the primary, which is the correct answer served by a more expensive machine.

There is no privilege here. The flag does not change *what* you can read, only *which replica answers*. Both hold the same data. A user who forges the cookie permanently has opted themselves out of the read replica and is getting slightly fresher data at slightly higher cost to you.

The only abuse worth a thought is a client forging it at scale to push load onto the primary, and that is a rate-limiting problem, not an authentication one. If someone can send you enough traffic for that to matter, they can hurt you without the cookie.

Not every mechanism needs integrity. The question is what an attacker gains, and here the answer is a marginally slower response for themselves.

### And it is free when you have no replicas

> It does no work on a project with no replicas: the flag is read by `database.ForRequest`, which returns the same handle either way when nothing is configured.

The middleware is always mounted. With one database, `ForRequest` returns that database whether the flag is set or not, so the whole thing is a cookie read and a map write. When you add a replica later, the consistency bug you would otherwise have discovered in production is already handled.

That is worth more than the mechanism itself: **the time to build this is before you need it**, because the moment you need it is the moment you have just added a replica and are busy debugging something that does not reproduce.

## The general shape

Strip out the specifics and the move is: *there was a fact two requests needed to share, and rather than put it in a shared store, we noticed the requests had a common carrier and wrote it there.*

That carrier is not always a cookie. It can be a header a client echoes back, a token's claims, a URL the first response hands out, an ETag. The question to ask before adding a store is whether the two parties that need to agree already have a channel between them.

When they do, the fact rides along for free, it cannot be stale, and there is nothing new to operate. When they do not, you build the store, having confirmed you need it.

---

*Grit is a full-stack Go and Next.js framework. `grit new my-app` mounts the middleware in this post by default, so adding a read replica later is a configuration change. [Get started](/docs/start), or read about [where a Grit app starts scaling](/blog/where-a-grit-app-starts-scaling).*
