package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeJobsFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "jobs", "client.go"):   jobsClientGo(),
		filepath.Join(apiRoot, "internal", "jobs", "workers.go"):  jobsWorkersGo(),
		filepath.Join(apiRoot, "internal", "handlers", "jobs.go"): jobsHandlerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func jobsClientGo() string {
	return `package jobs

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
	TypeEmailSend             = "email:send"
	TypeImageProcess          = "image:process"
	TypeTokensCleanup         = "tokens:cleanup"
	TypeUploadsOrphanCleanup  = "uploads:cleanup_orphans" // v3.31.33
	TypeBackupWeekly          = "backup:weekly"           // v3.31.77
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
	DefaultMaxRetries          = 5
	DefaultTaskTimeout         = 5 * time.Minute
	DefaultIdempotencyWindow   = 24 * time.Hour
	DefaultCompletedRetention  = 24 * time.Hour
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
	To       string                 ` + "`" + `json:"to"` + "`" + `
	Subject  string                 ` + "`" + `json:"subject"` + "`" + `
	Template string                 ` + "`" + `json:"template"` + "`" + `
	Data     map[string]interface{} ` + "`" + `json:"data"` + "`" + `
}

// ImagePayload holds the data for an image processing job.
type ImagePayload struct {
	UploadID string ` + "`" + `json:"upload_id"` + "`" + `
	Key      string ` + "`" + `json:"key"` + "`" + `
	MimeType string ` + "`" + `json:"mime_type"` + "`" + `
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
`
}

func jobsWorkersGo() string {
	return `package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"{{MODULE}}/internal/backup"
	"{{MODULE}}/internal/cache"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/files"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/storage"
)

// WorkerDeps holds dependencies needed by job handlers.
type WorkerDeps struct {
	DB      *gorm.DB
	Mailer  *mail.Mailer
	Storage *storage.Storage
	Cache   *cache.Cache
}

// ExponentialBackoff returns the delay before retry attempt n. The
// schedule is 1s, 2s, 4s, 8s, ... capped at 5 minutes — short enough to
// recover quickly from transient blips, long enough to be polite to
// downstream services that are actually overloaded.
//
// This is exported so callers can hand it to other queues if they roll
// their own, and so the docs page can link to a concrete reference.
func ExponentialBackoff(n int, _ error, _ *asynq.Task) time.Duration {
	if n < 0 {
		n = 0
	}
	// 1 << n overflows past 31 — clamp before shifting.
	if n > 30 {
		n = 30
	}
	d := time.Duration(1<<uint(n)) * time.Second
	if d > 5*time.Minute {
		d = 5 * time.Minute
	}
	return d
}

// StartWorker starts the asynq worker server in a goroutine.
// Returns a stop function and any startup error.
func StartWorker(redisURL string, deps WorkerDeps) (func(), error) {
	redisOpt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL for worker: %w", err)
	}

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 3,
			"default":  6,
			"low":      1,
		},
		// Explicit exponential backoff — see ExponentialBackoff for the
		// schedule. asynq's default is much slower (1min, 2min, ...) which
		// is too patient for the kind of HTTP-call jobs Grit apps run.
		RetryDelayFunc: ExponentialBackoff,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailSend, handleEmailSend(deps))
	mux.HandleFunc(TypeImageProcess, handleImageProcess(deps))
	mux.HandleFunc(TypeTokensCleanup, handleTokensCleanup(deps))
	mux.HandleFunc(TypeUploadsOrphanCleanup, handleUploadsOrphanCleanup(deps))
	mux.HandleFunc(TypeBackupWeekly, handleBackupWeekly(deps))

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Printf("Worker error: %v", err)
		}
	}()

	return func() {
		srv.Shutdown()
	}, nil
}

func handleEmailSend(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.Mailer == nil {
			return fmt.Errorf("mailer not configured")
		}

		var payload EmailPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return fmt.Errorf("unmarshaling email payload: %w", err)
		}

		log.Printf("Sending email to %s: %s", payload.To, payload.Subject)

		return deps.Mailer.Send(ctx, mail.SendOptions{
			To:       payload.To,
			Subject:  payload.Subject,
			Template: payload.Template,
			Data:     payload.Data,
		})
	}
}

func handleImageProcess(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.Storage == nil {
			return fmt.Errorf("storage not configured")
		}

		var payload ImagePayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return fmt.Errorf("unmarshaling image payload: %w", err)
		}

		log.Printf("Processing image: upload %s, key %s", payload.UploadID, payload.Key)

		// Download the original image
		reader, err := deps.Storage.Download(ctx, payload.Key)
		if err != nil {
			return fmt.Errorf("downloading image: %w", err)
		}
		defer reader.Close()

		// Generate thumbnail
		thumbBytes, err := storage.GenerateThumbnail(reader, payload.MimeType)
		if err != nil {
			return fmt.Errorf("generating thumbnail: %w", err)
		}

		// Upload thumbnail
		thumbKey := strings.Replace(payload.Key, "uploads/", "thumbnails/", 1)
		if err := deps.Storage.Upload(ctx, thumbKey, bytes.NewReader(thumbBytes), payload.MimeType); err != nil {
			return fmt.Errorf("uploading thumbnail: %w", err)
		}

		// Update the upload record with thumbnail URL
		thumbURL := deps.Storage.GetURL(thumbKey)
		if deps.DB != nil {
			deps.DB.Model(&models.Upload{}).Where("id = ?", payload.UploadID).Update("thumbnail_url", thumbURL)
		}

		log.Printf("Thumbnail created for upload %s", payload.UploadID)
		return nil
	}
}

func handleTokensCleanup(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.DB == nil {
			return fmt.Errorf("database not configured")
		}

		log.Println("Running token cleanup...")

		// Clean up soft-deleted records older than 30 days
		result := deps.DB.Exec("DELETE FROM users WHERE deleted_at IS NOT NULL AND deleted_at < NOW() - INTERVAL '30 days'")
		if result.Error != nil {
			return fmt.Errorf("cleaning up deleted users: %w", result.Error)
		}

		log.Printf("Token cleanup complete, removed %d records", result.RowsAffected)
		return nil
	}
}

// v3.31.33 -- orphan upload cleanup. Runs daily via the cron
// scheduler. Deletes Upload rows whose claimed_at IS NULL and
// created_at < 24h ago. The 24h grace period gives a user who
// uploaded a file plenty of time to finish the parent form before
// the cleanup considers their upload abandoned.
//
// minAge intentionally lives in code (not env) -- the value is
// load-bearing for correctness, not configuration. Bumping it
// would require thought about claim timing edge cases.
func handleUploadsOrphanCleanup(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.DB == nil {
			return fmt.Errorf("database not configured")
		}
		log.Println("Running orphan upload cleanup...")
		deleted, err := files.RunOrphanCleanup(ctx, deps.DB, deps.Storage, 24*time.Hour)
		if err != nil {
			return fmt.Errorf("orphan cleanup: %w", err)
		}
		log.Printf("Orphan upload cleanup complete, removed %d uploads", deleted)
		return nil
	}
}

// handleBackupWeekly takes the scheduled full-database backup: every registered
// model is dumped to a ZIP (CSV per table + dump.sql + metadata.json), uploaded
// to object storage, and old archives are pruned to the newest few.
//
// When object storage isn't configured (typical in local dev) we skip silently
// rather than fail the task forever — there's nowhere to put the archive.
func handleBackupWeekly(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.DB == nil {
			return fmt.Errorf("database not configured")
		}
		if deps.Storage == nil {
			log.Println("Weekly backup skipped: object storage is not configured")
			return nil
		}

		log.Println("Running weekly full-database backup...")
		svc := &backup.Service{DB: deps.DB, Storage: deps.Storage}
		rec, err := svc.Generate(ctx, "WEEKLY")
		if err != nil {
			return fmt.Errorf("weekly backup: %w", err)
		}
		log.Printf("Weekly backup %s complete — %d tables, %d rows, %.1f KB",
			rec.ID, rec.TableCount, rec.RowCount, float64(rec.SizeBytes)/1024)
		return nil
	}
}
`
}

func jobsHandlerGo() string {
	return `package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

// JobsHandler handles admin job queue endpoints.
type JobsHandler struct {
	RedisURL string
}

func (h *JobsHandler) getInspector() (*asynq.Inspector, error) {
	redisOpt, err := asynq.ParseRedisURI(h.RedisURL)
	if err != nil {
		return nil, err
	}
	return asynq.NewInspector(redisOpt), nil
}

// Stats returns queue statistics.
func (h *JobsHandler) Stats(c *gin.Context) {
	inspector, err := h.getInspector()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "REDIS_UNAVAILABLE",
				"message": "Job queue not available",
			},
		})
		return
	}
	defer inspector.Close()

	queues, err := inspector.Queues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch queue info",
			},
		})
		return
	}

	stats := make([]gin.H, 0)
	for _, q := range queues {
		info, err := inspector.GetQueueInfo(q)
		if err != nil {
			continue
		}
		stats = append(stats, gin.H{
			"queue":     info.Queue,
			"size":      info.Size,
			"active":    info.Active,
			"pending":   info.Pending,
			"completed": info.Completed,
			"failed":    info.Failed,
			"retry":     info.Retry,
			"scheduled": info.Scheduled,
			"processed": info.Processed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// ListByStatus returns jobs filtered by status.
func (h *JobsHandler) ListByStatus(c *gin.Context) {
	status := c.Param("status")
	queue := c.DefaultQuery("queue", "default")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	inspector, err := h.getInspector()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "REDIS_UNAVAILABLE",
				"message": "Job queue not available",
			},
		})
		return
	}
	defer inspector.Close()

	type jobInfo struct {
		ID        string ` + "`" + `json:"id"` + "`" + `
		Type      string ` + "`" + `json:"type"` + "`" + `
		Queue     string ` + "`" + `json:"queue"` + "`" + `
		MaxRetry  int    ` + "`" + `json:"max_retry"` + "`" + `
		Retried   int    ` + "`" + `json:"retried"` + "`" + `
		LastErr   string ` + "`" + `json:"last_error"` + "`" + `
	}

	var jobs []jobInfo

	opts := []asynq.ListOption{asynq.PageSize(pageSize), asynq.Page(page)}

	switch status {
	case "active":
		tasks, err := inspector.ListActiveTasks(queue, opts...)
		if err == nil {
			for _, t := range tasks {
				jobs = append(jobs, jobInfo{ID: t.ID, Type: t.Type, Queue: t.Queue, MaxRetry: t.MaxRetry, Retried: t.Retried, LastErr: t.LastErr})
			}
		}
	case "pending":
		tasks, err := inspector.ListPendingTasks(queue, opts...)
		if err == nil {
			for _, t := range tasks {
				jobs = append(jobs, jobInfo{ID: t.ID, Type: t.Type, Queue: t.Queue, MaxRetry: t.MaxRetry, Retried: t.Retried, LastErr: t.LastErr})
			}
		}
	case "completed":
		tasks, err := inspector.ListCompletedTasks(queue, opts...)
		if err == nil {
			for _, t := range tasks {
				jobs = append(jobs, jobInfo{ID: t.ID, Type: t.Type, Queue: t.Queue, MaxRetry: t.MaxRetry, Retried: t.Retried, LastErr: t.LastErr})
			}
		}
	case "failed":
		tasks, err := inspector.ListArchivedTasks(queue, opts...)
		if err == nil {
			for _, t := range tasks {
				jobs = append(jobs, jobInfo{ID: t.ID, Type: t.Type, Queue: t.Queue, MaxRetry: t.MaxRetry, Retried: t.Retried, LastErr: t.LastErr})
			}
		}
	case "retry":
		tasks, err := inspector.ListRetryTasks(queue, opts...)
		if err == nil {
			for _, t := range tasks {
				jobs = append(jobs, jobInfo{ID: t.ID, Type: t.Type, Queue: t.Queue, MaxRetry: t.MaxRetry, Retried: t.Retried, LastErr: t.LastErr})
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_STATUS",
				"message": "Status must be: active, pending, completed, failed, or retry",
			},
		})
		return
	}

	if jobs == nil {
		jobs = make([]jobInfo, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": jobs,
	})
}

// Retry re-enqueues a failed job.
func (h *JobsHandler) Retry(c *gin.Context) {
	id := c.Param("id")
	queue := c.DefaultQuery("queue", "default")

	inspector, err := h.getInspector()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "REDIS_UNAVAILABLE",
				"message": "Job queue not available",
			},
		})
		return
	}
	defer inspector.Close()

	if err := inspector.RunTask(queue, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "RETRY_FAILED",
				"message": "Failed to retry job: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job queued for retry",
	})
}

// ClearQueue deletes all tasks in a queue.
func (h *JobsHandler) ClearQueue(c *gin.Context) {
	queue := c.Param("queue")

	inspector, err := h.getInspector()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "REDIS_UNAVAILABLE",
				"message": "Job queue not available",
			},
		})
		return
	}
	defer inspector.Close()

	if _, err := inspector.DeleteAllCompletedTasks(queue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "CLEAR_FAILED",
				"message": "Failed to clear queue",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Queue cleared",
	})
}
`
}
