#import "../template.typ": *

= Background jobs, cron and transactional email

Some work must not happen inside a request. Resizing twelve images, emailing
four hundred people, reconciling yesterday's payments: if the user is waiting
for it, you have built a timeout.

This chapter is the queue, the scheduler and the mailer, and most of it is
about the two questions that separate a job system from a `go func()`: what
happens when it fails, and what happens when it runs twice.

== Why not just a goroutine

```go
go sendTheEmails(users)   // please do not
```

Three things are wrong with it, and the third is the one that hurts.

*It dies with the process.* A deploy, a crash, a scale-down, and the work is
gone with no record that it existed.

*It cannot be retried.* The error goes to a log, if anywhere.

*It cannot be seen.* Nobody can answer "did the invoice reminders run last
night", which is the question you will actually be asked.

#bridge[
  This is the same reason a Node application reaches for BullMQ rather than
  `setTimeout`. The difference in Go is that the goroutine is so cheap and so
  easy that the temptation is stronger, and it works in development every
  single time.

  The framework does keep a small escape hatch for genuinely minor work: a
  detached task with its own context and a hard cap of thirty two concurrent
  slots, deliberately small, because *anything needing more than thirty two
  workers is a job*.
]

== A job, generated

```bash
grit generate job SendInvoiceReminders --cron "0 9 * * *"
```

That writes three things: the job file, a handler registration in the worker's
mux, and, with `--cron`, a schedule entry that the admin's cron page lists.

The job file is yours from the moment it is written, and a second `generate`
refuses rather than overwrites.

```go
const TypeSendInvoiceReminders = "send_invoice_reminders"

type SendInvoiceRemindersPayload struct {
    Day string `json:"day,omitempty"`
}

func handleSendInvoiceReminders(deps WorkerDeps) func(
    context.Context, *asynq.Task) error {

    return func(ctx context.Context, task *asynq.Task) error {
        var p SendInvoiceRemindersPayload
        if len(task.Payload()) > 0 {
            if err := json.Unmarshal(task.Payload(), &p); err != nil {
                return fmt.Errorf("decoding: %v: %w",
                    err, asynq.SkipRetry)
            }
        }
        // deps.DB, deps.Mailer, deps.Storage, deps.Cache
        return nil
    }
}
```

Two conventions in that template are the whole discipline.

*The payload is IDs and small values.* It is serialised into Redis, and a retry
sees exactly what the first attempt saw. Putting a whole struct in the payload
means a retry works from a snapshot that may now be stale.

*An undecodable payload is wrapped with `asynq.SkipRetry`.* A payload that does
not parse now will not parse in four minutes either. Retrying a permanent
failure is how a queue fills up with work that can never succeed.

== Return an error to retry, return nil to stop

```go
return nil   // done, or can never succeed
return err   // retry with backoff
```

That is the entire protocol, and the second half of the first line is the part
people get wrong. A task that fails for a reason that will never change, a
deleted record, a malformed input, a file that is not an image, should return
`nil` and log loudly, or wrap with `SkipRetry`. Returning the error buys you
the same failure five more times with exponential backoff.

#gotcha[
  The practical rule: *retry anything that is about the world, do not retry
  anything that is about the data*.

  A network timeout, a rate limit, a database that is briefly down: retry.
  Invalid JSON, a missing user, an unsupported format: do not. If you cannot
  tell which one you have, you can log it and return `nil`, and the work stays
  visible in the logs rather than invisible in a retry loop.
]

== The two kinds of "it ran twice"

A queue guarantees *at least once*, not exactly once. The retry that happens
after a worker was killed mid-task is the normal case, not the pathological
one.

So every job has to answer: what if this runs again?

*Idempotency at the enqueue.* The same work queued twice is refused:

```go
queue.EnqueueSendInvoiceReminders(ctx, payload,
    jobs.EnqueueOption{
        IdempotencyKey: "reminders:" + day,
        Window:         24 * time.Hour,
    })
```

A second enqueue inside the window returns `ErrDuplicateTask` rather than
running. This is how a nightly job that is triggered twice, by a scheduler and
by an impatient human, stays a single run.

*Idempotency in the work.* Even one enqueue can execute twice. So the work
itself should be safe to repeat: mark the invoice as reminded in the same
transaction that records the send, and skip the ones already marked.

#gotcha[
  The order matters, and the wrong order is the one that feels natural.

  Send the email, then mark it sent: a crash between the two sends it again
  tomorrow. Mark it sent, then send: a crash between them means nobody is ever
  reminded.

  Neither is free, so choose the failure you prefer. For a reminder, a
  duplicate is better than a silence. For a payout, a silence is very much
  better than a duplicate, and the answer is usually a record written first
  with a status that a later pass reconciles.
]

== Cron, and the five fields

```bash
grit generate job NightlyReconcile --cron "30 23 * * *"
grit generate job Heartbeat --cron "@every 15m"
grit generate job Cleanup --cron "@daily"
```

The spec is validated when you generate, not at three in the morning:

```
--cron wants five fields (minute hour day-of-month month
day-of-week), e.g. "30 23 * * *" for 23:30 every day
```

Scheduled runs arrive with an empty payload, which is why the template handles
that case. The schedule is registered with a retry limit, and every registered
task appears on the admin's cron page with its name, its spec and its type, so
"did it run last night" is a screen rather than a grep.

#gotcha[
  Cron runs in the *server's* timezone. A daily report at `0 9 * * *` lands at
  nine in UTC, which is not nine for anybody you are reporting to.

  Set the timezone deliberately and write it down. And for anything that must
  happen on a date rather than at an hour, a daily job that checks for the date
  beats a schedule that tries to express it.
]

== Mail, and the drivers

```
MAIL_MAILER = smtp | resend | mailgun | postmark |
              sendgrid | ses | log | failover
```

Three behaviours in the configuration are worth copying into whatever you build
next.

*An unset driver does something sensible.* In development the mail goes to
Mailhog from Docker Compose, falling back to the log when Mailhog is not
running, so a password reset link is never lost.

*A named but misconfigured driver is an error, not a fallback.* An application
that believes it sends through SES must not quietly write to a log instead.

*Production with no mailer configured returns nothing*, and the code that
sends says so, rather than pretending.

#tryit[
  Set `MAIL_MAILER=log`, trigger a password reset, and read the link out of
  your terminal.

  It is the fastest way to work on an auth flow, and it is the setting to reach
  for whenever an email is in the way of what you are actually testing.
]

== Queued or inline, decided once

Transactional mail goes through one function, and the decision is made there:

```go
func dispatchMail(ctx context.Context, mailer *mail.Mailer,
    queue *jobs.Client, key string,
    opts mail.SendOptions) error {

    if queue != nil {
        return queue.EnqueueSendEmail(ctx, opts.To,
            opts.Subject, opts.Template, opts.Data,
            jobs.EnqueueOption{IdempotencyKey: key})
    }
    if mailer == nil {
        return ErrNoMailer
    }
    sendCtx, cancel := context.WithTimeout(
        context.WithoutCancel(ctx), 15*time.Second)
    defer cancel()
    return mailer.Send(sendCtx, opts)
}
```

With Redis, the mail is queued and retried. Without it, the mail is sent inline
with a timeout. The caller does not care, which is why a project can start
without Redis and gain it later without touching the sending code.

`context.WithoutCancel` is the detail worth stealing. The send is allowed to
finish after the request it came from has been answered, but not to run
forever: the request's context is cancelled the moment the response is written,
and inheriting that cancellation would kill the send mid-flight.

#gotcha[
  The `key` is the idempotency key, and the thing to pass is something stable
  for the *business action*, such as `"verify:" + user.ID`.

  A client retrying a request, or a proxy replaying it, then does not send the
  same person the same verification email twice. A key built from a timestamp
  or a random value deduplicates nothing at all, which is the common mistake.
]

== Seeing it work

Three places, and knowing they exist saves an afternoon:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Admin: Jobs], [Queued, active, retrying and dead tasks, with payloads],
  [Admin: Cron], [Every registered schedule and when it last ran],
  [Admin: Mail], [A preview of each template, rendered, without sending],
)

The mail preview is the one that changes how you work. Email templates are
miserable to iterate on through actual sends, and a page that renders them with
sample data turns a twenty minute loop into a reload.

#project(title: "Project 40: make the reminders trustworthy")[
  Using the invoicing application from chapter forty one.

  ```bash
  grit generate job SendInvoiceReminders --cron "0 9 * * *"
  ```

  *One: the work.* Find invoices that are sent, unpaid and past due, and email
  each client once. Use the mail dispatch function rather than the mailer
  directly.

  *Two: idempotency, twice over.* An `IdempotencyKey` of
  `"reminders:" + day` on the enqueue, and a `reminded_at` column checked and
  written in the same transaction as the send record. Then run the job three
  times in a row and prove one email.

  *Three: the failure.* Point the mailer at a broken SMTP host. Watch the retry
  with backoff in the admin's Jobs page, and find the task in the dead queue
  when it gives up. Then fix the host and replay it from the admin.

  *Four: a permanent failure.* Feed the job an invoice whose client has no
  email. It must not retry five times. Decide between `SkipRetry` and
  returning `nil` with a log, and write one sentence on why.

  *Five: the clock.* Change the cron spec to run in two minutes, deploy
  nothing, restart the worker, and confirm it fires. Note the timezone the
  server is in and what that means for `0 9 * * *`.

  *Six: the preview.* Open the mail preview page, find your reminder template,
  and fix the thing that is wrong with it. There is always something wrong with
  it, usually the subject line when the client name is empty.
]

#tryit[
  Kill the worker process while a job is running, then start it again.

  Watch what happens to the task. That is the "at least once" guarantee in
  front of you, and it is the single best argument for the idempotency work
  this chapter asks for.
]

#assignment(name: "Jobs and Mail - Assignment One")[
  In `part5-building/43-jobs`, build the reminder job and prove its behaviour
  under failure.

  *One.* The job, scheduled, with both layers of idempotency, and a test that
  runs the handler twice and asserts one email.

  *Two.* `JOBS.md` with four short sections: what this job retries, what it
  refuses to retry, what happens if it runs twice, and what happens if it never
  runs at all.

  *Three.* Evidence from the admin: a screenshot or a paste of the Jobs page
  showing a retry, and the dead queue entry after it exhausted its attempts.

  *Four.* The ordering question, answered for your own job. Which do you write
  first, the send or the mark, and which failure did you choose to live with?

  *Five.* One job of your own that is not mail: a nightly export, a cleanup of
  expired rows, a reconciliation. Schedule it, run it, and say in one line how
  you would know tomorrow morning whether it ran.
]
