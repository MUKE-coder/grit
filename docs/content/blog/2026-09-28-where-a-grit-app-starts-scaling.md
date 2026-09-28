---
title: "Your app is already at stage four: what a framework can and cannot do about scaling"
subtitle: "I took a well-known scaling guide, ten stages from one server to sharding, and worked out which ones a Grit app is born past. Five, it turns out, and two of the remaining five were missing entirely. This is what I built, what I refused to build, and the one decision no framework will ever make for you."
series: "The Daily Grit"
edition: 17
date: 2026-09-28
readingTime: "18 min"
author: "Muke JohnBaptist"
tags: [grit, scaling, read-replicas, caching, postgres, performance, architecture]
canonical: "https://gritframework.dev/blog/where-a-grit-app-starts-scaling"
---

Somebody sent me a scaling guide. It is a good one: ten stages, from one server and one database up to sharding, each with a symptom, a diagnosis, the code, and what the change costs you. Its central rule is the best thing in it, and I will quote it, because everything below is downstream of it:

> Don't add a piece of infrastructure until something is actually breaking, and then add exactly one thing.

The question I was asked was: how hard is it to do all of this on a Grit app, and what should Grit do to make it easier?

I expected the answer to be a long list of things to build. It was not. I sat down with the guide open on one screen and a freshly generated project on the other, and went stage by stage asking a single question: **if a Grit app hit this symptom today, what would the developer actually have to do?**

Five of the ten stages: nothing. Two of them: a real gap, and I have now filled them. One of them Grit should never do, and I will explain why.

---

## Stage 1 is the one nobody talks about

Every scaling guide has a stage that reads like housekeeping. This one calls it "make Stage 1 as good as it can be", and lists four things:

- configuration comes from environment variables, not hard-coded values
- database access goes through one module, not scattered `new Client()` calls
- no local file storage for user uploads
- a health check endpoint exists

Then it says something that is easy to skim past: *"This single decision makes Stages 5, 6, and 9 dramatically easier."*

That is the whole ballgame. The reason scaling is hard for most teams is not that read replicas are hard. It is that by the time you need read replicas, you have four hundred call sites that each opened their own connection, uploads sitting on a disk that is about to become three disks, and config baked into a file. The scaling work is mostly undoing Stage 1 decisions that were made by accident in week one.

A generated Grit project has all four on the first commit. Not because I am clever about scaling, but because a code generator has no way to write the scattered version. There is one `database` package because the generator emits one. Uploads go to object storage because the storage service is what `grit generate resource` wires up. `/health` exists because the scaffold writes it.

That is the first real finding, and it is not a feature I can point at. **A framework's largest contribution to scaling is the shape it forces before anybody is thinking about scaling at all.**

---

## Stage 2 does not apply, and that is a language choice

The guide spends a page on a genuine Node.js problem: one Node process uses one CPU core, so a four-core machine running a single process wastes three of them. You reach for `cluster`, or PM2 with `-i max`, and the moment you do you have several copies of your app and every Stage 4 problem arrives early.

A Go binary uses every core in one process. There is no cluster module, no PM2, and no early arrival of the shared-state problem. Vertical scaling is "buy a bigger machine", and then it is finished.

I am not claiming Go is better than Node. I am saying that one page of that guide is a tax on a runtime choice, and picking Go means not paying it. When people ask why a framework would use Go for the backend, this is a more honest answer than benchmarks.

---

## Stage 4 is where I expected to find a bug, and did not

Stage 4 is stateless servers, and it contains my favourite line in the guide:

> That last row is sneaky: with 3 instances, your "send daily digest" cron sends **three** emails.

This is the bug that gets everybody. You scale to three instances on Friday, and on Monday the support inbox has three copies of every digest. It is not caught by tests, because it does not happen on one machine.

I went looking for it. Grit's cron uses asynq's scheduler, and I fully expected each replica to run its own. It does not:

```
internal/cron/leader.go
  leaderKey   = "grit:cron:leader"
  leaderLease = 30 * time.Second
  leaderRenew = 10 * time.Second
```

A lease in Redis, renewed every ten seconds, and only the holder schedules. The comment above it says *"Every replica used to run its own, so each job ran once per replica"*, which means somebody hit the three-emails bug and fixed it properly rather than telling people to run one instance.

The rest of Stage 4 was the same story. Sessions are rows in a `sessions` table, not process memory, which is why "log me out of every device" works at all. Uploads were never on local disk. There is a `cluster` package doing generation counters so a permission revoked on one replica reaches the others within a second, without Redis, using a table the replicas already share.

And the server drains on `SIGTERM` rather than dropping in-flight requests during a deploy.

Stage 4, in a framework that had not been asked about scaling: done.

---

## Stage 8 was the surprise

Queues I expected to be half-built. Most frameworks give you "here is a job runner" and leave retries, backoff, dead letters and scheduling to you.

Grit ships asynq with retries, exponential backoff, a jobs dashboard, cron, and the thing I did not expect: a **transactional outbox**. The guide mentions it as an advanced footnote:

> if the database write succeeds but enqueueing fails (Redis blip), the job is lost. For critical jobs, use the transactional outbox pattern

That is a real gap in most systems and it is genuinely fiddly. There is an `internal/outbox` package and an `outbox_messages` table. The job goes into the database in the same transaction as the row that caused it, and a relay moves it to the queue afterwards. A Redis blip delays the email; it does not lose it.

---

## So what was actually missing?

Two stages, and they were missing completely.

### Stage 6: read replicas

Nothing. No routing, no configuration, no guidance. A Grit app with a read-heavy database had exactly one option: get a bigger database.

### Stage 7: caching

There was a `cache` package with `Get`, `Set` and `Delete`. That is the hard half done and the easy half left undone, because cache-aside written by hand comes out wrong in the same three ways every time.

Plus one thing the guide puts at Stage 0 and everybody skips: **measuring**. Grit had request tracing through Pulse, but no percentiles, and nothing that would answer the question the guide is entirely about, which is *which stage am I at?*

---

## What I built

### `grit scale`, which tells you to do nothing

The temptation with a scaling feature is to build a dashboard. Dashboards are worse than useless here, because the guide's actual insight is that **choosing** is the hard part. A screen with fourteen graphs has handed the choosing back to you with extra steps.

So `grit scale` measures a running deployment and names exactly one thing:

```
  Requests
    p50      3ms   p95      7ms   p99     13ms   max    164ms
    65 requests in the window, 5.0/s since start

  Database
    11 of 100 connections in use (11%), pool max 25 per instance
    No read replicas

  Cache
    Not configured

  Healthy
    p99 13ms over 65 requests, 11 of 100 database connections in use.

  Do this next: Nothing. Adding infrastructure now buys complexity and no speed.

  Worth knowing, but not yet:
    No read replicas, and you should not add any until the database's CPU is
    pegged by reads with the indexes already right.
    No cache configured. Add one when the same expensive query shows up near
    the top of the slowest list.
```

Three things about that output are deliberate.

**The most common answer is "nothing".** A tool that nags you to add read replicas at forty requests a minute is a tool you learn to ignore, and then it is there when you need it and you do not read it.

**It checks connection exhaustion before anything else**, because that is the only failure in the whole guide that arrives as *errors* rather than slowness. `sorry, too many clients already` while every dashboard is green. Everything else makes you slow; that one makes you wrong.

**It checks for missing indexes before it will mention a replica.** A sequential scan on a large table looks exactly like a capacity problem, and an index is free where a replica is a monthly bill:

```
  Tables being read end to end
    orders                        14820 scans of 2400000 rows
    An index is free. A bigger database is not.
```

The measuring happens in the API, not the CLI. Percentiles and connection counts are only true where the load is; measuring on your laptop describes your laptop. The API exposes `GET /api/v1/scale` (admin only, since it reports connection counts and query shapes) and the CLI reads it over the network, so you can point it at the thing that is actually struggling.

### Read replicas: one environment variable

```bash
DATABASE_REPLICA_URLS=postgres://...replica-1,postgres://...replica-2
```

That is the whole change. No handler edits, no `db.read()` / `db.write()` split to remember at four hundred call sites.

I used GORM's `dbresolver`, which routes **per statement** rather than per call site: `SELECT` goes to a replica, writes go to the primary, and anything inside a transaction goes to the primary *including its reads*.

That last rule is the one that matters, and it is the one hand-rolled splits get wrong. A balance check inside the transaction that debits the balance must never read a replica. With per-statement routing it cannot, whatever the handler was written to do. With a hand-rolled `db.read()`, it is one tired afternoon away.

For reads outside a transaction whose answer decides a write, there is an explicit door:

```go
// Never stale: the answer decides whether to sell the seat.
database.Primary(h.DB).First(&seat, "id = ?", id)
```

### Read-your-own-writes, with a cookie

The classic replica bug: you post, the feed loads from a replica that has not caught up, your own post is missing, so you post it again.

The standard fix keeps "user X wrote at T" in Redis. That means a Redis lookup on every read, and a Redis outage taking the read path with it.

A cookie carries the same fact, on the one request that needs it:

```go
c.SetCookie("grit_rw", now, 5, "/", "", secure, true)
```

It costs no lookup, it needs no shared store, and it cannot itself be stale, because it travels with the person who wrote. It is not a security boundary and does not need to be: the worst a forged cookie achieves is reading from the primary, which is the correct answer served by a more expensive machine.

### `cache.Remember`

```go
var followers int64
err := cache.Remember(ctx, svc.Cache, "user:"+id+":followers", time.Minute, &followers,
    func() (int64, error) {
        var n int64
        err := h.DB.Model(&models.Follow{}).Where("followee_id = ?", id).Count(&n).Error
        return n, err
    })
```

Cache-aside, with the three things the hand-written version misses:

1. **A cache read failure falls through to the loader.** Redis being down makes the app slower, not broken. That difference is the entire justification for letting a cache be a dependency.
2. **TTLs are jittered.** Ten thousand keys written by one deploy expire in the same second without it, and the database gets all of them back at once.
3. **Stampede protection.** When a hot key expires, one caller rebuilds it while the rest wait briefly and read the result, instead of every concurrent request running the same expensive query.

Plus a hit rate, because a cache nobody can measure is a cache nobody can tune, and `grit scale` will tell you when yours is under 50% and therefore not earning its keep.

### The arithmetic, in `grit doctor`

Stage 5 is the only stage that is pure arithmetic, so a tool can just do it:

```
⚠ database  8 instances x DB_MAX_OPEN_CONNS=25 is 200 connections, and
            max_connections is 100. Spikes will fail with "sorry, too many
            clients already" while every dashboard looks healthy
            put pgbouncer in front, or set DB_MAX_OPEN_CONNS=10
```

The number people get wrong is *instances*. A pool of 25 is comfortable on one machine and fatal on eight, and nothing else in your config mentions the other seven. Tell Grit with `APP_INSTANCES` and it does the multiplication.

---

## What I refused to build

**Sharding.** Grit does nothing for Stage 9, deliberately.

Sharding is the most expensive tool in the box, it is close to irreversible, and the overwhelming majority of applications never need it. Before it there is archiving cold rows, Postgres native partitioning, moving analytics to a warehouse and search to a search engine, and databases like Citus and CockroachDB that shard for you.

A framework that shipped a sharding helper would be telling you sharding is a normal thing to reach for. It is not. The most useful thing Grit can do about Stage 9 is refuse to make it look easy.

**A scaling dashboard.** Same reason as above: the choosing is the work.

**Automatic anything.** No autoscaling triggers, no "we noticed your cache hit rate is low so we adjusted your TTLs". Every one of these stages trades correctness for speed, and a framework that silently makes that trade on your behalf has made a product decision it has no standing to make.

---

## The thing no framework will ever do for you

Here is the list I cannot generate:

| Cache hard | Cache carefully | Never cache |
|---|---|---|
| Follower counts | Event listings with prices | Account balances |
| Public profiles | Seat maps, for display | Seat availability at purchase |
| Trending feeds | Notification counts | Order status during checkout |

Replicas and caches make your application **faster and slightly wrong on purpose**. Queues make *done* mean *promised*. Which parts of your product are allowed to be a little wrong, for a little while, is a product decision, and it is the one part of scaling that is genuinely yours.

Grit can put the replica routing one environment variable away. It can make the cache-aside pattern one function. It can do the connection arithmetic and tell you your index is missing before you buy a machine. It can refuse to shard.

It cannot tell you whether a stale seat map is fine. Showing one is fine. *Selling* based on one is not. Nobody but you knows which of your screens is which.

---

## Where a Grit app actually starts

Counting honestly, on the first commit, before a single user:

- **Stage 1** done: env config, one database module, object storage, health endpoint
- **Stage 2** does not apply: Go uses every core
- **Stage 3** done: stateless, graceful shutdown, health checks for the balancer
- **Stage 4** done: DB sessions, cron leader election, cross-replica invalidation
- **Stage 5** ready: pgbouncer already in the production compose file
- **Stage 8** done: workers, retries, backoff, cron, transactional outbox

Stages 6 and 7 are now one environment variable and one function call.

That leaves you with the part that was always the real work: knowing which one you need, in what order, and what each one costs.

Start simple. Measure. Break it. Fix exactly one thing.

```bash
grit scale --api https://api.yourapp.com
```

The whole thing is at [gritframework.dev/docs/scaling](https://gritframework.dev/docs/scaling).
