package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"commerce/apps/api/internal/audit"
	"commerce/apps/api/internal/backup"
	"commerce/apps/api/internal/cache"
	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/mail"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/storage"
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
	mux.HandleFunc(TypeAuditPrune, handleAuditPrune(deps))
	mux.HandleFunc(TypeUploadsOrphanCleanup, handleUploadsOrphanCleanup(deps))
	mux.HandleFunc(TypeBackupWeekly, handleBackupWeekly(deps))
	mux.HandleFunc(TypeBackupScheduled, handleBackupScheduled(deps))
	// grit:jobs

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

		// mail.Queue carries the whole message. EnqueueSendEmail carries a
		// template name and its data, handled below.
		var queued mail.QueuedMessage
		if err := json.Unmarshal(task.Payload(), &queued); err != nil {
			return fmt.Errorf("unmarshaling email payload: %v: %w", err, asynq.SkipRetry)
		}
		if queued.Message != nil {
			log.Printf("Sending queued email: %s", queued.Message.Subject)
			return deps.Mailer.SendMessage(ctx, queued.Message)
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

		// Generate thumbnail. An image refused as too large, or one that cannot
		// be decoded, fails the same way on every attempt, so it is not retried.
		// A decompression bomb used to be decoded in this process and retried
		// until the retries ran out, restarting the API each time.
		thumbBytes, err := storage.GenerateThumbnail(reader, payload.MimeType)
		if errors.Is(err, storage.ErrImageTooLarge) || errors.Is(err, storage.ErrUnreadableImage) {
			return fmt.Errorf("generating thumbnail for %s: %v: %w", payload.Key, err, asynq.SkipRetry)
		}
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
			// Returned, so the job is retried: a thumbnail nobody can find was
			// generated for nothing.
			if err := deps.DB.WithContext(ctx).Model(&models.Upload{}).Where("id = ?", payload.UploadID).Update("thumbnail_url", thumbURL).Error; err != nil {
				return fmt.Errorf("recording the thumbnail for upload %s: %w", payload.UploadID, err)
			}
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

		// Clean up soft-deleted records older than 30 days.
		//
		// The cutoff is computed in Go and bound as a parameter rather than
		// written as NOW() - INTERVAL '30 days'. That syntax is Postgres-only:
		// on SQLite — a first-class target here, and the quick-start default —
		// it fails with a syntax error near the interval literal, so this job
		// errored on every run and the cleanup silently never happened.
		cutoff := time.Now().AddDate(0, 0, -30)
		result := deps.DB.WithContext(ctx).Exec(
			"DELETE FROM users WHERE deleted_at IS NOT NULL AND deleted_at < ?", cutoff)
		if result.Error != nil {
			return fmt.Errorf("cleaning up deleted users: %w", result.Error)
		}

		log.Printf("Token cleanup complete, removed %d records", result.RowsAffected)
		return nil
	}
}

// handleAuditPrune trims the tamper-evident activity log through audit.Prune,
// which deletes in chunks of audit.PruneChunk, one transaction each, and records
// every chunk in the chain so what remains still verifies.
//
// Retention is AUDIT_RETENTION_DAYS, default 365. Zero or negative disables
// pruning entirely, which is the right default for anyone who has not thought
// about it: silently discarding audit history is worse than a large table.
func handleAuditPrune(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.DB == nil {
			return fmt.Errorf("database not configured")
		}

		days := 365
		if v := os.Getenv("AUDIT_RETENTION_DAYS"); v != "" {
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("AUDIT_RETENTION_DAYS must be a whole number of days, got %q", v)
			}
			days = parsed
		}
		if days <= 0 {
			log.Println("Audit prune skipped: AUDIT_RETENTION_DAYS <= 0 (retain forever)")
			return nil
		}

		removed, err := audit.Prune(ctx, deps.DB, time.Now().AddDate(0, 0, -days))
		if err != nil {
			return fmt.Errorf("pruning the activity log after %d entries: %w", removed, err)
		}
		if removed > 0 {
			log.Printf("Audit prune complete: removed %d entries older than %d days", removed, days)
		}
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
		log.Printf("Weekly backup %s complete: %d tables, %d rows, %.1f KB",
			rec.ID, rec.TableCount, rec.RowCount, float64(rec.SizeBytes)/1024)
		return nil
	}
}

// handleBackupScheduled is the settings-driven auto-backup checker. It runs on a
// frequent cron tick and consults the BackupSchedule row: if a backup is due for
// the current period (daily/weekly/monthly/yearly at the configured time) and
// none has run yet, it takes one. This lets the period change at runtime without
// touching the cron registration.
func handleBackupScheduled(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.DB == nil {
			return fmt.Errorf("database not configured")
		}
		if deps.Storage == nil {
			return nil // no storage in local dev — silently skip
		}
		svc := &backup.Service{DB: deps.DB, Storage: deps.Storage}
		due, err := svc.DueNow(time.Now())
		if err != nil {
			return fmt.Errorf("backup due check: %w", err)
		}
		if !due {
			return nil
		}
		log.Println("Scheduled backup is due: running full-database backup...")
		rec, err := svc.Generate(ctx, "SCHEDULED")
		if err != nil {
			return fmt.Errorf("scheduled backup: %w", err)
		}
		log.Printf("Scheduled backup %s complete: %d tables, %d rows, %.1f KB",
			rec.ID, rec.TableCount, rec.RowCount, float64(rec.SizeBytes)/1024)
		return nil
	}
}
