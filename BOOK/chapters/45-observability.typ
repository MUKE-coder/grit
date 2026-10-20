#import "../template.typ": *

= Observability: Studio, Pulse and Sentinel

B11 was "when something goes wrong, nobody can tell me what happened". The
honest version of that problem is that the information existed, in a log file,
on a server, in a format nobody could search, during an hour when the person
who could read it was asleep.

A generated project mounts three dashboards and a health endpoint. This chapter
is what each one answers, and the discipline that makes them useful rather than
decorative.

== Four questions, four tools

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`/studio`], [What is actually in the database?],
  [`/pulse`], [What is this request doing, and why is it slow?],
  [`/sentinel`], [Who is attacking us, and what did we do about it?],
  [`/health`], [Is everything this process depends on alive?],
)

They answer different questions, and reaching for the wrong one is how a
debugging session takes three hours. "The invoice shows the wrong total" is
Studio. "The invoice page takes four seconds" is Pulse.

== Studio: the database, as it is

The embedded database browser at `/studio` lists tables, browses rows, follows
foreign keys and runs queries. It is the fastest way to answer a question of
the form "did that actually save".

The habit worth forming: *before debugging the code, look at the row*. Half of
all "the API is broken" reports are a row that is exactly what the code wrote,
and the bug is upstream in what was sent.

#gotcha[
  Studio sees everything, which is the point and also the risk. The scoping
  that makes a user's query return only their rows is in the service layer, and
  Studio is underneath it.

  So a row you can see in Studio and cannot fetch through the API is usually
  ownership scoping doing its job rather than a bug. Chapter thirty five's 404
  is the same observation from the other side.
]

== Pulse: where the time went

Pulse mounts at `/pulse` and records request traces, database queries, runtime
metrics and errors. It is the tool for the question "why is this slow", and it
answers it in the only way that helps: by showing which part took the time.

The three views that earn their keep:

*The request trace.* One request, with its queries listed and timed. An
endpoint taking 800ms turns out to be forty queries of 20ms, which is the N+1
from chapter twenty two, now visible rather than suspected.

*The slow query list.* Sorted by total time rather than by worst case, because
a 5ms query run ten thousand times costs more than a 2 second report nobody
runs.

*Errors with their context.* Not just the stack, but the request that caused
it.

```go
pulse.WithExcludePaths("/studio/*", "/sentinel/*",
    "/docs/*", "/pulse/*")
```

The dashboards exclude each other, which sounds obvious and is a real bug
otherwise: an observability tool observing its own polling fills itself with
its own traffic.

#gotcha[
  A story worth carrying, because the shape of it is common.

  An earlier version of the request-body capture read the body to record it and
  put back only the first four kilobytes. Every request carrying a
  `Content-Length` reached the handler truncated. Uploads from a phone and from
  `curl` failed; uploads from a browser worked, because browsers send chunked
  bodies with no `Content-Length`.

  So the bug report was "uploads fail on mobile", and nothing in the upload
  code was wrong. *Middleware that reads a request body and puts it back is a
  place to look* whenever a request behaves differently depending on who sent
  it.
]

Pulse can keep its data in SQLite so traces survive a restart, or stay in an
in-memory ring buffer for peak throughput. The trade is the usual one: the ring
buffer costs nothing and forgets, the file costs writes and remembers.

== Sentinel: the attacks you were not watching

Sentinel is the security suite: a web application firewall, rate limiting, an
auth shield and anomaly detection, with a dashboard at `/sentinel`.

What it changes is not that you are attacked. It is that you can see it. A
login endpoint takes credential stuffing from the day it is public, and without
a dashboard the only evidence is a slightly busier log.

Two configuration details are worth knowing before you meet them.

*Development gets relaxed limits.* A rate limit tight enough for production
blocks you while you are testing, and a developer who has been blocked by their
own firewall usually turns the whole thing off.

*Its storage is its own.* Sentinel persists the threat log, blocked IPs and
audit trail through its own adapter, not the application's database handle.
Left unset it falls back to a local SQLite file, which is ephemeral inside a
container, so every redeploy drops the threat log and the blocked-IP list.
Point it at the real database in production.

#gotcha[
  A mount failure is logged rather than fatal. Sentinel refuses to start on a
  misconfiguration, and in development that should not take the API down with
  it.

  The consequence to watch for: *a security layer that failed to start looks
  exactly like one that is running, unless you read the log line*. Check the
  dashboard after a configuration change rather than assuming.
]

== Health: the one endpoint a load balancer reads

```json
{
  "status": "degraded",
  "components": {
    "api":      { "state": "ok" },
    "database": { "state": "ok", "latency_ms": 2 },
    "redis":    { "state": "off",
                  "detail": "REDIS_URL is empty, so caching, jobs
                             and cron are off on purpose" },
    "jobs":     { "state": "unknown" }
  }
}
```

The distinction that makes this endpoint useful is between *off* and *broken*.
Redis absent because nobody configured it is a working project without caching.
Redis configured and unreachable is an incident. A health check that reports
both as "red" gets ignored within a week, and a health check that reports both
as "green" is worse.

#bridge[
  The equivalent in most Node projects is `app.get('/health', () =>
  res.send('ok'))`, which answers "is the process running" and nothing else.

  The process is almost never the thing that broke. The database connection
  pool is exhausted, Redis went away, the disk is full. A health check that
  cannot say which is a status page that is always green during an outage.
]

== The logs, and the ID that joins everything

Every request carries a request ID, returned in a header and attached to every
log line it produces. That is what turns "it failed around two o'clock" into
one search.

```bash
grit logs                # the tail of this run
grit logs --errors       # only failures, with stack traces
grit logs --grep 500     # lines matching something
```

`grit start` records everything the API and the frontends print into
`.grit/logs/dev.log` and keeps the previous run beside it, which matters more
than it sounds: the interesting output is usually from the run that just
crashed and scrolled away.

#gotcha[
  Log the error and return a generic message, as chapter forty eight's handler
  did. The log gets the table name, the query and the stack. The client gets a
  code and a sentence.

  The request ID is what joins them, so it belongs in both: in the response
  header a user can copy, and in every log line. Without it, a support ticket
  with a timestamp is a search through a thousand lines around that minute.
]

== Making it a habit

Observability is not installed, it is practised. Four habits that cost nothing:

*Read Pulse after building a feature*, not after a complaint. The N+1 you find
on the day you write the code costs ten minutes.

*Check health after every deploy.* Not "the site loads", the endpoint.

*Look at Sentinel weekly.* You will learn what normal traffic looks like, which
is the only way to recognise abnormal traffic.

*Put the request ID in your error messages.* Then a screenshot from a user is
actionable.

#project(title: "Project 42: find three real problems")[
  Take any application you have built in Part V, put realistic data in it, and
  go looking.

  ```bash
  grit seed Invoice --count 5000
  grit start
  ```

  *One: an N+1 you did not know about.* Open the busiest list page, then open
  Pulse and read the trace. Count the queries. If it is one or two, add a
  relationship to the list and look again.

  *Two: the slowest query.* Sort Pulse's query list by total time. Explain in
  `OBSERVED.md` why that query is top: is it slow, or is it run too often?
  Those need different fixes.

  *Three: a real attack.* Leave the application running on a public port for an
  hour, or point a credential-stuffing script at your own login endpoint, and
  then read Sentinel. Record what it saw and what it did.

  *Four: health under failure.* Stop Redis and read `/health`. Stop the
  database and read it again. Confirm the states differ and that the detail
  says something a tired person could act on at 3am.

  *Five: the request ID.* Cause an error, find the ID in the response, and find
  every log line for that request with one command. Time yourself. That number
  is what observability is worth.

  *Six: the fix.* Fix the N+1 from step one, measure again, and put both
  numbers in the notes.
]

#tryit[
  Turn Pulse off in a staging environment and try to answer "why was the
  invoice page slow yesterday at 4pm" from logs alone.

  Ten minutes of that is the most persuasive argument for instrumentation you
  will ever get, and it is worth experiencing deliberately rather than during
  an incident.
]

#assignment(name: "Observability - Assignment One")[
  In `part5-building/45-observability`, instrument an application and report
  what you found.

  *One.* `OBSERVED.md` with three real findings from Pulse: a query count, a
  slow query with its total time, and one error with its request ID.

  *Two.* A before and after for one fix, with timings from Pulse rather than
  from a stopwatch.

  *Three.* The health endpoint in three states: everything up, a dependency
  deliberately off, and a dependency deliberately broken. Paste all three
  bodies and say which one should page somebody.

  *Four.* One paragraph on what you would add. Every application has a question
  these tools cannot answer, usually a business one: how many invoices were
  sent today, how many signups failed at the verification step. Name yours, and
  say where it would live.
]
