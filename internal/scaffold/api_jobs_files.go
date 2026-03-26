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
		filepath.Join(apiRoot, "internal", "jobs", "client.go"):    jobsClientGo(),
		filepath.Join(apiRoot, "internal", "jobs", "workers.go"):   jobsWorkersGo(),
		filepath.Join(apiRoot, "internal", "handlers", "jobs.go"):  jobsHandlerGo(),
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
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Task type constants.
const (
	TypeEmailSend     = "email:send"
	TypeImageProcess  = "image:process"
	TypeTokensCleanup = "tokens:cleanup"
)

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

// EmailPayload holds the data for an email send job.
type EmailPayload struct {
	To       string                 ` + "`" + `json:"to"` + "`" + `
	Subject  string                 ` + "`" + `json:"subject"` + "`" + `
	Template string                 ` + "`" + `json:"template"` + "`" + `
	Data     map[string]interface{} ` + "`" + `json:"data"` + "`" + `
}

// ImagePayload holds the data for an image processing job.
type ImagePayload struct {
	UploadID uint   ` + "`" + `json:"upload_id"` + "`" + `
	Key      string ` + "`" + `json:"key"` + "`" + `
	MimeType string ` + "`" + `json:"mime_type"` + "`" + `
}

// EnqueueSendEmail enqueues an email send job.
func (c *Client) EnqueueSendEmail(to, subject, template string, data map[string]interface{}) error {
	payload, err := json.Marshal(EmailPayload{
		To:       to,
		Subject:  subject,
		Template: template,
		Data:     data,
	})
	if err != nil {
		return fmt.Errorf("marshaling email payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailSend, payload)
	_, err = c.client.Enqueue(task, asynq.MaxRetry(3))
	if err != nil {
		return fmt.Errorf("enqueuing email job: %w", err)
	}
	return nil
}

// EnqueueProcessImage enqueues an image processing job.
func (c *Client) EnqueueProcessImage(uploadID uint, key, mimeType string) error {
	payload, err := json.Marshal(ImagePayload{
		UploadID: uploadID,
		Key:      key,
		MimeType: mimeType,
	})
	if err != nil {
		return fmt.Errorf("marshaling image payload: %w", err)
	}

	task := asynq.NewTask(TypeImageProcess, payload)
	_, err = c.client.Enqueue(task, asynq.MaxRetry(2))
	if err != nil {
		return fmt.Errorf("enqueuing image job: %w", err)
	}
	return nil
}

// EnqueueTokensCleanup enqueues a token cleanup job.
func (c *Client) EnqueueTokensCleanup() error {
	task := asynq.NewTask(TypeTokensCleanup, nil)
	_, err := c.client.Enqueue(task, asynq.MaxRetry(1))
	if err != nil {
		return fmt.Errorf("enqueuing cleanup job: %w", err)
	}
	return nil
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

	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"{{MODULE}}/internal/cache"
	"{{MODULE}}/internal/mail"
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
			"default":  6,
			"critical": 3,
			"low":      1,
		},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailSend, handleEmailSend(deps))
	mux.HandleFunc(TypeImageProcess, handleImageProcess(deps))
	mux.HandleFunc(TypeTokensCleanup, handleTokensCleanup(deps))

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

		log.Printf("Processing image: upload %d, key %s", payload.UploadID, payload.Key)

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

		log.Printf("Thumbnail created for upload %d", payload.UploadID)
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
