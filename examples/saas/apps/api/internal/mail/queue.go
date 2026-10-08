package mail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// TaskSend is the background job queued mail travels on. internal/jobs routes
// the same task type to the worker that sends it.
const TaskSend = "email:send"

// MaxQueuedAttachmentBytes caps the attachments a queued message may carry,
// 256 KB in all. The message waits in Redis until the worker sends it, and
// again for every retry, so a large file belongs in storage with a link to it
// in the message.
const MaxQueuedAttachmentBytes = 256 << 10

// ErrAttachmentsTooLarge is returned by Queue for a message whose attachments
// come to more than MaxQueuedAttachmentBytes.
var ErrAttachmentsTooLarge = errors.New("mail: attachments too large to queue")

// Enqueuer is the part of internal/jobs.Client that Queue uses. It is declared
// here because internal/jobs imports this package.
type Enqueuer interface {
	EnqueuePayload(ctx context.Context, taskType string, payload []byte) error
}

// QueuedMessage is the email:send payload Queue writes and the worker reads.
type QueuedMessage struct {
	Message *Message `json:"message"`
}

// Queue hands msg to the background worker, which sends it with the Mailer
// main.go built and retries a provider that is briefly down. The request that
// queues it does not wait for the provider.
//
// Mail somebody is waiting for, a password reset or a sign-in code, is better
// sent directly with Mailer.SendMessage, so a failure shows while they are
// still there.
func Queue(ctx context.Context, q Enqueuer, msg *Message) error {
	if msg == nil {
		return errors.New("mail: nil message")
	}
	if q == nil {
		return errors.New("mail: no job queue to send through (Redis is not configured); send with Mailer.SendMessage instead")
	}
	total := 0
	for _, a := range msg.Attachments {
		total += len(a.Content)
	}
	if total > MaxQueuedAttachmentBytes {
		return fmt.Errorf("%w: %q has %d bytes of attachments and a queued message may carry %d (256 KB). Store the file with internal/storage and send a link to it instead",
			ErrAttachmentsTooLarge, msg.Subject, total, MaxQueuedAttachmentBytes)
	}
	// Checked now, while the caller can still see the error. The worker would
	// find out after the request had already answered. From may be empty here:
	// the worker's Mailer fills it in.
	check := *msg
	if check.From == "" {
		check.From = "queued@example.com"
	}
	if err := check.Validate(); err != nil {
		return err
	}
	payload, err := json.Marshal(QueuedMessage{Message: msg})
	if err != nil {
		return fmt.Errorf("mail: encoding the queued message: %w", err)
	}
	if err := q.EnqueuePayload(ctx, TaskSend, payload); err != nil {
		return fmt.Errorf("mail: queueing %q: %w", msg.Subject, err)
	}
	return nil
}
