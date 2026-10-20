#import "../template.typ": *

= The event bus, realtime, and webhooks

An application grows a particular kind of mess: one thing happens, and five
systems have to hear about it. An order is placed, so write the audit record,
push to the open dashboard, send the webhook, queue the receipt, update the
cache. Every one of those is a line in the handler, and the sixth is the one
somebody forgets.

This chapter is the alternative: the handler says what happened, once, and
anything interested listens.

== The problem, in the framework's own history

Before the bus existed, a generated handler called the audit logger and nothing
else. Webhooks, realtime and background jobs all existed, and *none of them was
ever fired by generated code*. Dispatching a webhook meant hand-writing the
call in every handler, for every operation, on every resource.

Forgetting one produced no error at all. Just a subscriber that never heard
anything, which is the worst failure mode there is: silent, and only noticed
by the customer whose integration stopped working.

```go
events.Emit(ctx, events.Created{
    Resource: "order", ID: order.ID, Actor: actor,
})
```

One line in the handler. The audit log, outbound webhooks and realtime clients
are now subscribers rather than things somebody has to remember to call.

== Three deliveries, and choosing between them

```go
const (
    Sync Delivery = iota
    Async
    Durable
)
```

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`Sync`], [Runs inline, before the handler responds. For work that must be
    part of the request],
  [`Async`], [Runs on the bus's worker pool after the response. For anything
    with a network call],
  [`Durable`], [Through the transactional outbox: at least once, retried,
    surviving restarts],
)

*Sync* is for the audit record: it needs the request context and should be
written before the caller is told the write succeeded. A slow `Sync`
subscriber is a slow endpoint, and that is the whole cost model.

*Async* is for the webhook endpoint that takes four seconds. The API must not
take four seconds because somebody else's server is slow.

*Durable* is for work another module *depends on*. If losing the event would
leave two systems disagreeing, it is durable.

#gotcha[
  The async queue is bounded, and that is deliberate. An unbounded queue turns
  one slow subscriber into memory exhaustion, which fails later and worse than
  dropping events and saying so.

  When you see events being dropped in the logs, the fix is not a bigger
  queue. It is finding the subscriber that is slow, and usually moving it to
  `Durable` where retrying is the design rather than an accident.
]

== The outbox, and the choice it removes

`Durable` delivery exists because of a problem with two halves, and doing the
obvious thing gets one of them wrong whichever order you pick.

*Publish, then commit.* The publish succeeds, the commit fails, and a webhook
has announced an order that does not exist. Downstream systems act on it, and
there is nothing to reconcile against because the row was never written.

*Commit, then publish.* The commit succeeds, the process is killed before the
publish, and the order exists with nobody told. Nothing is logged anywhere,
because from the process's point of view nothing failed.

The outbox removes the choice:

```
BEGIN
  INSERT INTO orders ...
  INSERT INTO outbox_messages ...
COMMIT
```

The message is written in the *same transaction* as the business data, so it
commits or rolls back with it. A relay reads committed messages and delivers
them afterwards. Either both happened or neither did.

#bridge[
  This pattern is worth carrying with you whatever you build next, because it
  is not a framework feature. It is the standard answer to "I need to update my
  database and tell somebody else, atomically", and it works in any language
  against any queue.

  The name to look up is *transactional outbox*. Once you have seen it, you
  start noticing the places you got away without it.
]

The cost is honest and must be stated: delivery is *at least once*, not exactly
once. A message can be delivered twice if the process dies between the send and
the status update. Consumers have to be idempotent, which is what a dedup key
on the receiving end is for, and which is the same discipline chapter forty
three asked for.

== Realtime: channels, not just users

The hub can address a user, everybody, or a named channel. Channels are the one
that makes a real interface possible:

```js
socket.send(JSON.stringify({
  action: "subscribe",
  channel: "invoice.42",
}))
```

The server answers `subscribed`, `subscription_error` or `unsubscribed`, each
naming the channel, and server code publishes to the name:

```go
hub.Publish("invoice.42", payload)
```

Every subscriber on every replica receives it, which is the part that is
tedious to build: a second instance of your API has to reach clients connected
to the first, so the hub has a backplane.

#gotcha[
  A channel name is a permission question, and it must be answered server side.
  `invoice.42` means "whoever may read invoice 42", not "whoever can type that
  string".

  The hub has a registry of authorizers per channel pattern for this reason. A
  realtime layer where the client picks what it hears is an API with no access
  control, reachable from a browser console.
]

Without channels, a page showing one invoice had no way to hear about that
invoice without every event for that user being pushed at it, which is a real
cost on a busy account.

== When WebSockets do not arrive

Some platforms do not pass a WebSocket upgrade through to the application. One
managed host strips `Connection: upgrade` before the container, and it is not
alone: a corporate proxy, an older load balancer, or an API gateway configured
for plain HTTP all do the same.

The failure is unpleasant: the handshake fails, the client's reconnect loop
retries forever, and live updates silently never arrive. Nothing errors in a
way anybody sees.

So there is a Server-Sent Events fallback. SSE is plain HTTP with a response
that never ends, so it survives all of that.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [WebSocket], [Two-way, lower latency, needs the upgrade to survive the
    network],
  [SSE], [Server to client only, plain HTTP, reconnects by itself],
)

For dashboards, notifications and live lists, which is most of what an
application needs, SSE is enough. Reach for WebSockets when the client needs to
*send* continuously: a cursor position, a drawing, a chat typing indicator.

== Webhooks, outbound

```bash
grit plugin add webhooks
```

That wires outbound, Standard-Webhooks-signed deliveries with retries, a
delivery log and manual replay. Core already *verifies* incoming webhooks;
this *sends* outgoing ones, which is the opposite direction and a different
problem.

What matters when you are the one sending:

*Sign every delivery.* The receiver has to know it was you. Standard Webhooks
is a published convention for the headers and the signature, which beats
inventing your own and documenting it badly.

*Include an ID and a timestamp.* Together they let a receiver reject replays
and deduplicate retries, which they will need because you are delivering at
least once.

*Retry with backoff, and keep a log.* Somebody will ask "did you send it", and
the answer has to be a record with a status code, not a shrug.

*Make replay a button.* Integrations break on the other side. A delivery log
with a replay action turns a support conversation into one click.

#gotcha[
  A webhook subscriber must be `Async` or `Durable`, never `Sync`. Your API's
  response time should not depend on a receiver you do not control, and some
  receivers take seconds.

  The same goes for the timeout on the delivery itself. An HTTP client with no
  timeout, pointed at a stranger's server, is a slow leak of connections that
  ends in an API that will not answer anybody.
]

#project(title: "Project 41: wire one event to four listeners")[
  Using the invoicing application again, because an invoice being paid is a
  genuinely interesting event.

  *One: emit.* In `MarkPaid`, emit one event with the invoice ID, the amount
  and the actor. Nothing else changes in that method.

  *Two: four subscribers.*

  ```
  audit     Sync      writes the ledger entry
  realtime  Async     publishes to "invoice.<id>"
  mail      Async     queues the receipt
  webhook   Durable   notifies the accounting integration
  ```

  Each in its own file. None of them is called by the handler.

  *Three: watch the ordering.* Add a deliberate two second sleep to the
  realtime subscriber and time the endpoint. Then move it to `Sync` and time it
  again. That difference is the whole reason the delivery modes exist.

  *Four: the outbox.* Make the webhook durable, then kill the API process
  immediately after a payment, before the delivery. Restart and confirm the
  webhook still goes out. Then look in `outbox_messages` and read the row.

  *Five: the rollback.* Force the transaction to fail after the outbox write.
  Confirm there is no message and no delivery: that is the half of the problem
  publish-then-commit gets wrong.

  *Six: the browser.* Open an invoice page, subscribe to its channel, and mark
  it paid from a second window. The page should update with no polling. Then
  try to subscribe to a channel for an invoice your user may not read, and
  confirm the server refuses.
]

#tryit[
  Count the lines in your `MarkPaid` method before and after this chapter.

  It should be the same length or shorter, while doing four more things. That
  is the measure of whether an event bus is paying for itself: the handler
  describes what happened, and the list of interested parties is a
  configuration rather than a sequence of calls nobody can safely reorder.
]

#assignment(name: "Events - Assignment One")[
  In `part5-building/44-events`, wire one domain event to several listeners
  and prove the delivery guarantees.

  *One.* One `Emit`, four subscribers, one per file, with the delivery mode of
  each chosen and justified in a comment.

  *Two.* Timings: the endpoint with the slow subscriber on `Sync` and on
  `Async`, both recorded in `EVENTS.md`.

  *Three.* The crash test. Kill the process between commit and delivery, show
  the outbox row, restart, and show the delivery completing. Paste both states
  of the row.

  *Four.* The rollback test. A failed transaction leaves no message and no
  webhook.

  *Five.* A channel authorizer that refuses a subscription the user may not
  have, with the client's error response pasted in. One paragraph on what an
  unauthorised channel would expose in your own application.
]
