package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Task type constants.
const (
	TypeEmailSend            = "email:send"
	TypeImageProcess         = "image:process"
	TypeTokensCleanup        = "tokens:cleanup"
	TypeUploadsOrphanCleanup = "uploads:cleanup_orphans" // v3.31.33
	TypeBackupWeekly         = "backup:weekly"           // v3.31.77 (legacy)
	TypeBackupScheduled      = "backup:scheduled"        // settings-driven auto backup
	TypeAuditPrune           = "audit:prune"             // activity-log retention
)

// Default per-task settings used when a caller doesn't override via
// EnqueueOption. Tuned for production safety, not raw throughput:
//
//   - MaxRetries: 5 attempts is enough to ride out a 5-10 min downstream
//     blip given the exponential backoff schedule, without flooding the
//     dead queue with hours-old work.
//   - Timeout: 5 min cap so a stuck job can't lock up its worker slot
//     forever. Long jobs should pass a higher value explicitly.
//   - DefaultIdempotencyWindow: 24h is long enough that a same-key retry
//     hours later is still deduped, short enough not to hoard memory.
const (
	DefaultMaxRetries         = 5
	DefaultTaskTimeout        = 5 * time.Minute
	DefaultIdempotencyWindow  = 24 * time.Hour
	DefaultCompletedRetention = 24 * time.Hour
)

// ErrDuplicateTask is returned by Enqueue when the same IdempotencyKey
// is enqueued twice within its Window. Callers should treat it as a
// successful enqueue — the original task is on its way.
var ErrDuplicateTask = errors.New("jobs: duplicate task within idempotency window")

// EnqueueOption configures how a single task is queued. Zero-value fields
// fall back to package defaults.
type EnqueueOption struct {
	// IdempotencyKey, if non-empty, deduplicates enqueues with the same
	// key within Window. A retry of the same business operation (e.g.,
	// "send receipt for order 1234") can pass the same key and be safe.
	// Implementation note: this maps to asynq.TaskID — duplicate enqueues
	// return ErrDuplicateTask from this package.
	IdempotencyKey string

	// Window is how long the IdempotencyKey blocks re-enqueues. After
	// Window passes, the same key can be enqueued again. Defaults to
	// DefaultIdempotencyWindow.
	Window time.Duration

	// MaxRetries is the number of attempts (including the first) before
	// the task lands in the dead-letter queue. Defaults to
	// DefaultMaxRetries. Set to 0 for "no retries".
	MaxRetries int

	// Timeout caps how long a single attempt can run. Defaults to
	// DefaultTaskTimeout. Long-running jobs (PDF render, large export)
	// should set this explicitly.
	Timeout time.Duration

	// Queue overrides the queue (default, critical, low). Defaults to
	// "default".
	Queue string

	// Delay schedules the task in the future. Useful for "send reminder
	// in 24h" patterns. Defaults to immediate.
	Delay time.Duration

	// Retention is how long a SUCCEEDED task stays visible in the
	// inspector. Defaults to DefaultCompletedRetention.
	Retention time.Duration
}

// Client wraps asynq.Client for enqueuing background jobs.
type Client struct {
	client *asynq.Client
}

// NewClient creates a new job queue client connected to Redis.
func NewClient(redisURL string) (*Client, error) {
	redisOpt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL for jobs: %w", err)
	}

	client := asynq.NewClient(redisOpt)
	return &Client{client: client}, nil
}

// Close shuts down the client connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Enqueue submits a task with framework defaults plus any caller-supplied
// EnqueueOption. payload may be any JSON-marshalable value; if it's
// already []byte it's used as-is so callers that hand-encode can skip the
// reflection cost.
//
// Returns ErrDuplicateTask when an IdempotencyKey is supplied and a task
// with the same key is still within its Window — the original is already
// queued, so the caller's intent is satisfied.
func (c *Client) Enqueue(ctx context.Context, taskType string, payload any, opts ...EnqueueOption) error {
	var raw []byte
	switch v := payload.(type) {
	case nil:
		raw = nil
	case []byte:
		raw = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshaling task payload: %w", err)
		}
		raw = b
	}

	var opt EnqueueOption
	if len(opts) > 0 {
		opt = opts[0]
	}
	asynqOpts := []asynq.Option{
		asynq.MaxRetry(orDefaultInt(opt.MaxRetries, DefaultMaxRetries)),
		asynq.Timeout(orDefaultDuration(opt.Timeout, DefaultTaskTimeout)),
		asynq.Retention(orDefaultDuration(opt.Retention, DefaultCompletedRetention)),
	}
	if opt.Queue != "" {
		asynqOpts = append(asynqOpts, asynq.Queue(opt.Queue))
	}
	if opt.Delay > 0 {
		asynqOpts = append(asynqOpts, asynq.ProcessIn(opt.Delay))
	}
	if opt.IdempotencyKey != "" {
		// asynq.TaskID + Unique together give us "exactly one of this key
		// in flight at a time". Re-enqueues during the window return
		// asynq.ErrDuplicateTask, which we translate to our typed
		// ErrDuplicateTask so callers don't have to import asynq.
		asynqOpts = append(asynqOpts,
			asynq.TaskID(opt.IdempotencyKey),
			asynq.Unique(orDefaultDuration(opt.Window, DefaultIdempotencyWindow)),
		)
	}

	task := asynq.NewTask(taskType, raw)
	_, err := c.client.EnqueueContext(ctx, task, asynqOpts...)
	if err != nil {
		if errors.Is(err, asynq.ErrDuplicateTask) || errors.Is(err, asynq.ErrTaskIDConflict) {
			return ErrDuplicateTask
		}
		return fmt.Errorf("enqueuing %s: %w", taskType, err)
	}
	return nil
}

func orDefaultInt(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

func orDefaultDuration(v, def time.Duration) time.Duration {
	if v <= 0 {
		return def
	}
	return v
}

// EmailPayload holds the data for an email send job.
type EmailPayload struct {
	To       string                 `json:"to"`
	Subject  string                 `json:"subject"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

// ImagePayload holds the data for an image processing job.
type ImagePayload struct {
	UploadID string `json:"upload_id"`
	Key      string `json:"key"`
	MimeType string `json:"mime_type"`
}

// EnqueueSendEmail enqueues an email send job.
//
// The optional opts argument lets callers supply an IdempotencyKey so
// that a retry of the same business action (e.g., "send order receipt
// for order 1234") doesn't fire the email twice. ErrDuplicateTask from
// the underlying Enqueue is treated as success here.
func (c *Client) EnqueueSendEmail(ctx context.Context, to, subject, template string, data map[string]interface{}, opts ...EnqueueOption) error {
	payload := EmailPayload{To: to, Subject: subject, Template: template, Data: data}
	err := c.Enqueue(ctx, TypeEmailSend, payload, opts...)
	if errors.Is(err, ErrDuplicateTask) {
		return nil
	}
	return err
}

// EnqueuePayload queues a task whose payload is already encoded, with the
// default options. mail.Queue sends through it: internal/mail cannot import
// this package, so it asks for this one method.
func (c *Client) EnqueuePayload(ctx context.Context, taskType string, payload []byte) error {
	if c == nil || c.client == nil {
		return errors.New("jobs: no queue client, Redis is not configured")
	}
	err := c.Enqueue(ctx, taskType, payload)
	if errors.Is(err, ErrDuplicateTask) {
		return nil
	}
	return err
}

// EnqueueProcessImage enqueues an image processing job. Image uploads
// have a natural idempotency key (the upload's UUID): callers should set
// EnqueueOption.IdempotencyKey = uploadID to prevent double-processing
// after a retry.
func (c *Client) EnqueueProcessImage(ctx context.Context, uploadID string, key, mimeType string, opts ...EnqueueOption) error {
	payload := ImagePayload{UploadID: uploadID, Key: key, MimeType: mimeType}
	err := c.Enqueue(ctx, TypeImageProcess, payload, opts...)
	if errors.Is(err, ErrDuplicateTask) {
		return nil
	}
	return err
}

// EnqueueTokensCleanup enqueues a token cleanup job. The cleanup is
// naturally idempotent (DELETE WHERE deleted_at < ...) so it's safe to
// re-run; we still cap retries low and pin the queue to "low".
func (c *Client) EnqueueTokensCleanup(ctx context.Context) error {
	err := c.Enqueue(ctx, TypeTokensCleanup, nil, EnqueueOption{
		MaxRetries: 1,
		Queue:      "low",
	})
	if errors.Is(err, ErrDuplicateTask) {
		return nil
	}
	return err
}
