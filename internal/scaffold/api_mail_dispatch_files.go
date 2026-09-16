package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Contact-app review M30: mail left this API from bare goroutines.
//
// Every email a handler sent was sent from a `go func()` the request started.
// Nothing bounded them, so a burst of sign-ups was a burst of goroutines each
// holding a database connection and a provider socket. Nothing retried them, so
// a provider down for a minute lost the mail for good. A deploy in the middle
// dropped whatever was in flight. Meanwhile jobs.Client.EnqueueSendEmail, which
// does all three properly, was never called from anywhere.
//
// internal/handlers/mail_dispatch.go is where mail leaves a handler now.

// writeMailDispatchFiles writes internal/handlers/mail_dispatch.go.
//
// Its own writer so grit upgrade can deliver it whole: it is a new file, and
// the handlers that were repaired to call it do not compile without it.
func writeMailDispatchFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	path := filepath.Join(apiRoot, "internal", "handlers", "mail_dispatch.go")
	content := strings.ReplaceAll(apiMailDispatchGo(), "{{MODULE}}", opts.Module())
	if err := writeFile(path, content); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// apiMailDispatchGo returns internal/handlers/mail_dispatch.go.
func apiMailDispatchGo() string {
	return `// Mail leaves a handler through this file, and never from a bare goroutine.
//
// Two things live here. dispatchMail puts a message on the background queue
// when the project has one, so it is retried and survives a restart, and sends
// it inline when it does not, which is what happens in development without
// Redis. detach runs the little that genuinely has to stay off the request
// path, with a cap on how much of it can be in flight at once.
package handlers

import (
	"context"
	"errors"
	"log"
	"time"

	"{{MODULE}}/internal/jobs"
	"{{MODULE}}/internal/mail"
)

// ErrNoMailer is returned when a project has neither a queue nor a mailer, so
// the caller can log the link instead of pretending the mail went out.
var ErrNoMailer = errors.New("handlers: no mail queue and no mailer configured")

// mailSendTimeout caps an inline send. A provider that hangs must not hold a
// request, or a detached slot, open forever.
const mailSendTimeout = 15 * time.Second

// dispatchMail queues opts when there is a queue, and sends it inline when
// there is not.
//
// key deduplicates the enqueue: a client retrying a request, or a proxy
// replaying it, does not send the same person the same email twice. Pass
// something stable for the business action, such as "verify:" + user.ID.
func dispatchMail(ctx context.Context, mailer *mail.Mailer, queue *jobs.Client, key string, opts mail.SendOptions) error {
	if queue != nil {
		return queue.EnqueueSendEmail(ctx, opts.To, opts.Subject, opts.Template, opts.Data, jobs.EnqueueOption{
			IdempotencyKey: key,
		})
	}
	if mailer == nil {
		return ErrNoMailer
	}
	// WithoutCancel: the send is allowed to finish even when the request it
	// came from is already answered, but not to run forever.
	sendCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), mailSendTimeout)
	defer cancel()
	return mailer.Send(sendCtx, opts)
}

// detachedSlots caps how much work requests may push off their own path.
//
// 32 is deliberately small. Anything that needs more than 32 concurrent workers
// is a job, and belongs on the queue with the mail.
var detachedSlots = make(chan struct{}, 32)

// detachTimeout is how long a detached task may run before its context is
// cancelled. It owns its context: the request's is cancelled the moment the
// response is written.
const detachTimeout = 30 * time.Second

// detach runs fn off the request path, with at most 32 in flight at once.
//
// Over the cap it runs fn inline instead of starting a 33rd goroutine or
// dropping the work: a slow request is a better answer than an unbounded one,
// and much better than an email nobody ever receives.
//
// Use it only where finishing inside the request would leak something, such as
// whether an email address has an account here. Everything else belongs on the
// queue, where it is retried.
func detach(fn func(ctx context.Context)) {
	run := func() {
		defer func() {
			if r := recover(); r != nil {
				// A panic in a detached task used to take the process with it,
				// because nothing was watching the goroutine it ran in.
				log.Printf("detached task panicked: %v", r)
			}
		}()
		ctx, cancel := context.WithTimeout(context.Background(), detachTimeout)
		defer cancel()
		fn(ctx)
	}
	select {
	case detachedSlots <- struct{}{}:
		go func() {
			defer func() { <-detachedSlots }()
			run()
		}()
	default:
		run()
	}
}
`
}
