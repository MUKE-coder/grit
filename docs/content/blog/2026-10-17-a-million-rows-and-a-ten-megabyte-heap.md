---
title: "A million rows, sixty seconds, and a ten megabyte heap"
subtitle: "Grit's seeder is a two-pool pipeline: goroutines building rows on one side, goroutines writing batches on the other, a bounded channel between them. Here is why it is two pools and not one, why the batch size is computed rather than chosen, and why the heap stays flat at ten megabytes whether you ask for a thousand rows or a million."
series: "The Daily Grit"
edition: 20
date: 2026-10-17
readingTime: "14 min"
author: "Muke JohnBaptist"
tags: [grit, go, concurrency, performance, database, seeding, pipelines]
canonical: "https://gritframework.dev/blog/a-million-rows-and-a-ten-megabyte-heap"
---

Somebody sent me [Kennedy Karoko's two-part series on processing a million transactions in under a second](https://blog.karoko.dev/processing-1-million-transactions-in-under-a-second-using-go-part-1-98c9079375ac) and asked whether it is the same thing Grit does when you run `grit seed Product --count 1000000`.

It is the same family, and the differences are the interesting part. So here is the honest comparison, followed by the parts of Grit's seeder that go further, and the one idea in that article that Grit does not use anywhere and probably should.

First, the number, because a post about performance that does not measure anything is an opinion:

```
$ grit seed Product --count 1000000

products: 756000 of 999800 rows (16434 rows/s, about 14s left)
products: 968000 of 999800 rows (16689 rows/s, about 1s left)
products: seeded 999800 rows in 59.718s (16742 rows/s), now 1000000; heap in use 10 MB
```

Sixty seconds, on SQLite, on a laptop. The number that matters is not 16,742 rows per second. **It is the 10 MB.**

## The shape both pipelines share

Karoko's article builds up to it from the naive version, which is the right way to teach it. The naive version reads a file line by line and writes each row as it goes, and gets about 350 transactions per second. The final version:

1. One goroutine reads, pushing into a buffered channel.
2. A worker pool consumes, parsing concurrently.
3. Each worker batches, and aggregates into its own private map.
4. The maps merge at the end, and the result is flushed in bulk.

Grit's seeder, in `internal/database/seedbatch.go`, has the same skeleton. A bounded channel, producers on one end, consumers on the other, batched writes:

```go
chunks := make(chan []T, writers*2)
```

And the two headline conclusions are identical, so let me state them plainly because they are the whole lesson:

**Batching the writes is the win, not the goroutines.** Going from one insert per row to one insert per batch is the change that moves the number by an order of magnitude. Parallelism on top of that is worth maybe five times. Parallelism without batching is worth almost nothing, because a million round trips is a million round trips no matter how many goroutines are waiting on them.

**Bounded memory is a property of the channel, not of discipline.** A buffered channel of fixed capacity means a fast producer blocks when consumers fall behind. That backpressure is what keeps the heap flat. It is not something you remember to do; it is something the data structure does whether you remember or not.

## Where Grit's version differs

### Two pools, not one

The article has one worker pool that parses *and* writes. Grit has two pools with the channel between them:

```go
makers := runtime.NumCPU()
if makers > 4 {
    makers = 4
}
// ... makers build rows and push []T into chunks ...

for w := 0; w < writers; w++ {
    go func() {
        for rows := range chunks {
            db.Transaction(func(tx *gorm.DB) error {
                return tx.CreateInBatches(rows, len(rows)).Error
            })
        }
    }()
}
```

The reason is that building a row and writing a row are different kinds of work, and the right number of goroutines for each is different. Building rows is CPU-bound: `gofakeit.ProductName()`, a money value, a random number. More than a few cores of that is wasted, so makers cap at four. Writing is I/O-bound and bounded by what the database will accept concurrently, which is a completely different number.

Fuse them into one pool and you have to pick a single worker count that is wrong for one half of the job. Split them and each side gets sized on its own terms, with the channel absorbing the mismatch.

### The writer count is a property of the database

```go
writers = 3
if db.Dialector.Name() == "sqlite" {
    writers = 1
}
```

SQLite takes a single writer lock. Pointing eight goroutines at it does not produce eight times the throughput; it produces eight goroutines taking turns, plus lock contention and `SQLITE_BUSY` retries. One writer is not a limitation being worked around, it is the correct number for that engine.

This is the general form of the article's own closing observation, that pushing the worker count higher stops helping past a point. The point is not a universal constant. It is a property of what sits at the end of the pipeline, and the code should know which one it is talking to.

### The batch size is computed, not chosen

The article suggests 500 to 1,000 and says to benchmark. Reasonable advice. Grit computes it:

```go
func seedBatchSize(db *gorm.DB, model any) int {
    limit := 65535 // Postgres and MySQL
    if db.Dialector.Name() == "sqlite" {
        limit = 32766
    }
    cols := 32
    if err := stmt.Parse(model); err == nil && len(stmt.Schema.DBNames) > 0 {
        cols = len(stmt.Schema.DBNames)
    }
    size := limit / cols
    if size > 1000 {
        size = 1000
    }
    return size
}
```

There is a hard ceiling here that benchmarking will not reveal gently. Postgres and MySQL accept 65,535 bind parameters in one statement; SQLite accepts 32,766. A multi-row insert uses one parameter per column per row, so the real limit is `65535 / columns`. A model with 12 columns tops out at 5,461 rows per insert; a wide model with 60 columns tops out at 1,092.

Pick 1,000 by benchmarking on a narrow table and it works. Add columns to the model over the next year and one day the insert fails with a parameter-limit error that has nothing obviously to do with the column you added. Deriving it means the ceiling moves when the model does.

The `if size > 1000` cap on top is the benchmarking answer: past about a thousand rows the gains flatten and the cost of a failed batch grows, so there is no reason to go near the ceiling even when it is allowed.

### Work is claimed, not dispatched

The article's producer pushes each line into the channel and the workers pull. Grit's makers pull their own work from a counter:

```go
var next atomic.Int64

for {
    k := next.Add(1) - 1
    if k >= chunkCount {
        return
    }
    start := have + k*int64(batch)
    // ... build rows [start, end) ...
}
```

`next.Add(1)` is a single atomic increment. There is no dispatcher goroutine, no mutex, and no coordination between makers beyond that one instruction. Each maker computes which rows it owns from the index it claimed, which means it can build its chunk without reading any shared state at all.

This is the same insight as the article's worker-local maps, applied to work distribution instead of to aggregation: **the fastest synchronisation is the synchronisation you do not perform.** The article removes lock contention from the aggregation hot path. This removes it from the scheduling hot path. Same idea, different place.

### It resumes, which changes what the index means

```go
var have int64
db.Model(new(T)).Count(&have)
if have >= plan.Target {
    return nil // nothing to do
}
need := plan.Target - have
```

The seeder counts what is in the table and inserts only the difference. Run it twice and the second run does nothing. Kill it halfway through a million rows and run it again and it picks up where it stopped.

That is why `Make` takes an index:

```go
Make: func(n int64) models.Product {
    return models.Product{
        Name:  gofakeit.ProductName(),
        Sku:   gofakeit.Word(),
        // ...
    }
},
```

The index `n` runs from the number of rows already present, not from zero. Mix it into a unique column and a resumed run cannot collide with the rows the previous run wrote. Start from zero every time and the second run fails on a unique constraint, which is exactly the moment people give up and add `ON CONFLICT DO NOTHING`, which quietly turns a seeder into something that reports success over a half-empty table.

### The first error stops everything, carefully

```go
var firstErr error
var errOnce sync.Once
fail := func(err error) {
    errOnce.Do(func() { firstErr = err; cancel() })
}
```

and in the writer loop:

```go
for rows := range chunks {
    if ctx.Err() != nil {
        continue // drain so the makers can finish
    }
    // ...
}
```

That `continue` is the subtle line, and it is the kind of thing a pipeline tutorial usually skips. When a write fails, the context is cancelled, and the natural instinct is to `return` from the writer immediately. Do that and the makers block forever pushing into a channel nobody is reading, and `makersDone.Wait()` never returns. The process hangs instead of reporting the error.

Draining instead of returning means the makers finish, the channel closes, the writers exit, and the first error surfaces. A cancelled pipeline has to unwind in the same direction the data flows, and the failure mode for getting it wrong is a deadlock rather than a wrong answer, which is why it is worth saying out loud.

## The idea Grit does not use, and where it should

The strongest technique in Karoko's article is worker-local aggregation. Rather than every worker updating one shared `map[string]*UserAnalytics` behind a mutex, each worker owns a private map and the maps merge once at the end. Zero lock contention in the hot path, one coordination cost at the finish.

**Grit's seeder does not use it, and should not, because seeding aggregates nothing.** Each row is independent. There is no shared state to contend over, so there is nothing for the technique to improve.

But Grit has places that do aggregate under concurrency, and the one that stands out is **import**. `POST /products/import` reads a CSV or XLSX and inserts rows concurrently, and anything that wants per-file totals, per-column validity counts, or a tally of rejected rows is exactly the workload the article describes. The same goes for any insight or report that sweeps a table and groups as it goes.

If you are writing that kind of aggregation in a Grit project, take the pattern from the article directly: give each worker its own map, merge at the end, and do not reach for `sync.Map` or a mutex just because the word "concurrent" appeared. A `sync.Map` is for a map with a stable key set read far more than written, which is the opposite of an accumulator.

## What transfers to your own code

The seeder is one application of a shape that shows up constantly: a producer, a bounded channel, a consumer pool, batched writes. Four things are worth carrying over.

**Size the pools separately when the work differs in kind.** CPU-bound and I/O-bound stages want different worker counts, and a single pool forces one number on both.

**Let the channel do the backpressure.** The reason the heap stays at 10 MB for a million rows is that `chunks` holds `writers*2` batches and no more. No cleanup, no high-water mark, no tuning. If you find yourself writing memory-management code inside a pipeline, the channel is probably unbounded.

**Derive limits from the system instead of picking them.** Bind-parameter ceilings, single-writer engines, column counts: these are facts the code can read. A constant is a guess that stops being true when something else changes, and the failure arrives long after the change that caused it.

**Decide what cancellation does before you need it.** Cancelling a pipeline is where the deadlocks live, and every one of them is a stage that stopped reading while an earlier stage was still writing.

The question that started this was whether Grit's seeder is the same concept as the article. Yes, and the agreement is the point: a bounded channel, a worker pool and batched writes is just what this problem looks like when you solve it properly. The parts worth arguing about are the ones underneath, where the right answer depends on which database is at the end of the pipe and what happens when one batch fails.

---

*Grit is a full-stack Go and Next.js framework. `grit new my-app` gives you an API, an admin panel and a web app, with the seeder in this post generated for every resource you create. [Get started](/docs/start).*
