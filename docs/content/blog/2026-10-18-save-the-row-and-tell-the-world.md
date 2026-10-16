---
title: "Save the row and tell the world, without losing one or the other"
subtitle: "Publish then commit, and a webhook announces an order that does not exist. Commit then publish, and the order exists with nobody told. There is no correct order, which is why the transactional outbox exists. Here is the whole pattern, including the lease that stops two replicas sending the same message twice."
series: "The Daily Grit"
edition: 21
date: 2026-10-18
readingTime: "15 min"
author: "Muke JohnBaptist"
tags: [grit, go, distributed-systems, outbox, postgres, webhooks, reliability]
canonical: "https://gritframework.dev/blog/save-the-row-and-tell-the-world"
---

Here is a bug that is almost impossible to reproduce and very easy to write.

An order is placed. You save it, then fire a webhook so the warehouse starts picking. Two lines of code, in the obvious order:

```go
db.Create(&order)
webhooks.Send("order.created", order)
```

That is wrong, and swapping the lines is also wrong. There is no third option. This post is about why, and about the thing you have to build instead.

## Both orders are broken

**Publish first, then commit.** The webhook goes out, the commit fails on a constraint or a dropped connection, and the warehouse is now picking an order that does not exist in your database. There is nothing to reconcile against, because the row was never written. Downstream systems have acted on a fact you do not have.

**Commit first, then publish.** The commit succeeds, the process is killed between the two lines, and the order exists with nobody told. The warehouse never picks it. No error is logged anywhere, because from the process's point of view nothing failed. You find out when the customer emails.

Neither is a rare race. A deploy restarts every process, by design, continuously. If you ship often enough, that gap gets hit.

The usual instinct is to make the gap smaller: retry the publish, defer it, wrap it in a goroutine with a timeout. All of that makes the window narrower and none of it closes it, because the problem is not the size of the window. The problem is that you are performing two operations that must either both happen or neither happen, and you have no mechanism that makes them atomic.

You do have exactly one such mechanism: the database transaction. So use it.

## The outbox

Write the message to a table in the same transaction as the business data.

```go
tx.Create(&order)
tx.Create(&outbox.Message{Topic: "order.created", Payload: payload})
```

One commit. Both rows or neither. There is no window, because there is no second operation to fail.

A separate relay then reads committed messages and delivers them. If delivery fails, the message is still sitting in the table and gets tried again. If the process dies mid-delivery, the message is still sitting in the table and gets tried again.

The package comment in `internal/outbox/outbox.go` puts the trade in one line, and it is worth being blunt about it because this is where people get surprised:

> The cost is delivery semantics: at-least-once, not exactly-once. A message can be delivered twice if the process dies between the send and the status update. Consumers have to be idempotent, which is what a dedup key on the receiving end is for.

**At-least-once is not a flaw in the design, it is the design.** You have moved the hard problem from "did this get sent at all" to "might this arrive twice", and the second problem has a cheap, local, verifiable answer: the receiver records message IDs it has processed and ignores repeats. The first problem has no answer at all.

Anyone who tells you their queue is exactly-once is describing at-least-once delivery plus deduplication, which is the same thing with the dedup hidden somewhere you cannot inspect.

## The relay, and the part people get wrong

A naive relay is a loop: select pending messages, send each one, mark it sent. That works on one process. Run two replicas and both select the same rows, and every message goes out twice. Not occasionally. Every time.

The fix is to claim work rather than read it. Here is the query the relay actually runs, lifted from a running project's query log:

```sql
SELECT * FROM outbox_messages
WHERE ((status = ? AND available_at <= ?)
    OR (status = ? AND claimed_at < ?))
  AND topic LIKE ?
ORDER BY created_at
LIMIT ?
```

Two conditions joined by OR, and the second one is the one that matters.

The first clause is ordinary pending work: messages whose status is pending and whose `available_at` has arrived. The `available_at` column is what makes backoff possible, which we will get to.

The second clause picks up messages that are marked *claimed* but whose `claimed_at` is older than some threshold. That is a **lease**. A relay claims a message by stamping `claimed_at` and marking it in-flight. If it delivers successfully, the message is marked sent. If the process dies holding the claim, nobody marks anything, and after the lease expires another relay sees a claim older than the threshold and takes it over.

Without a lease, a process that dies mid-delivery strands its claimed messages forever. They are not pending, so nobody picks them up, and they sit in the table looking fine. The lease is what turns "claimed" from a permanent state into a temporary one.

The lease duration is the usual availability-versus-duplication dial. Short, and a slow delivery gets taken over by a second relay while the first is still working, producing a duplicate. Long, and a crashed process's messages sit idle for that long. Since consumers must be idempotent anyway, erring short is usually right.

## Retries that do not hammer a dead endpoint

```go
if r.MaxAttempts <= 0 {
    r.MaxAttempts = 12
}
if r.BaseBackoff <= 0 {
    r.BaseBackoff = 5 * time.Second
}
if r.MaxBackoff <= 0 {
    r.MaxBackoff = 30 * time.Minute
}
```

Exponential backoff from five seconds, doubling, capped at thirty minutes, twelve attempts before the message is parked.

Those three numbers say something specific. Twelve attempts with doubling backoff capped at half an hour covers roughly five hours of a downstream system being unavailable. That is a deliberate statement: a webhook endpoint that has been down for five hours is not having a blip, and continuing to retry is not helping anybody. The message is parked for a human rather than retried forever.

Parked, and specifically not deleted:

> MaxAttempts is how many times to try before parking a message at "failed". Parked rather than deleted: a message nobody can deliver is evidence of a bug, and deleting the evidence is how the bug survives.

A dead-letter row is the only artefact that says what was supposed to happen and did not. Delete it and the system looks healthy while silently dropping work, which is the worst of both outcomes: no delivery and no trace.

The cap matters as much as the doubling. Uncapped exponential backoff reaches days between attempts surprisingly fast, which means a system that recovers at hour three waits until hour nine to find out. The cap bounds the detection delay at thirty minutes.

`available_at` is how backoff is stored. Rather than keeping a retry schedule in memory, which a restart loses, the next attempt time is written on the row. The relay's query only picks up messages whose `available_at` has passed. Restart the process and the schedule is still there, because it is in the database with everything else.

## What this costs, honestly

It is not free, and the costs are worth naming.

**Latency.** A message is delivered when the relay next polls, not at commit. That is a poll interval of delay on every message. Fine for webhooks, order events and email, and not fine for anything interactive.

**A table that grows.** Delivered messages accumulate and need sweeping. The alternative, deleting on delivery, loses the audit trail of what was sent and when, which is the thing you want most during an incident.

**Ordering is per-topic at best.** The relay orders by `created_at` and processes a batch, so two messages created in the same millisecond can be delivered in either order, and a retried message arrives after messages created later. If your consumer needs strict ordering, the outbox alone does not give it to you.

**Another moving part.** A relay is a background process that can be stopped, can lag, and needs monitoring. A pending count that climbs is the signal, and it needs to be on a dashboard, because the failure mode of a stopped relay is total silence.

## When you need it and when you do not

You need it when a database write and an external effect must agree: an order and its webhook, a payment and its receipt, a user and the welcome email, a subscription and the entitlement call to a third party.

You do not need it when the external effect is advisory. A cache invalidation, an analytics event, a log line: if losing one occasionally changes nothing, the obvious two lines are fine, and an outbox is machinery with no benefit.

The test is what happens the one time the publish is lost. If the answer is "nothing, it is a metric", skip it. If the answer is "a customer is owed something nobody knows about", you have found an outbox.

## The general lesson

The reason this pattern keeps reappearing under different names is that it is the answer to a structural fact: **you get exactly one atomic boundary, and it is the database transaction.** Anything that must agree with a database write has to be made to happen inside that boundary, and the only thing you can put inside it is another database write.

So you write down the intent, and let something else carry it out. The outbox is that idea applied to messages. The same shape appears whenever an effect must be as durable as the data: a job queue backed by the same database, a scheduled action that survives a restart, a two-phase workflow that records what it is about to do before it does it.

Once you see it, the question for any side effect becomes: is this inside the transaction boundary, or am I hoping?

---

*Grit is a full-stack Go and Next.js framework. `grit new my-app` scaffolds the outbox, its relay and the webhook delivery in this post. [Get started](/docs/start).*
