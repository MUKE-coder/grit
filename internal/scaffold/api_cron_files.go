package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeCronFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "cron", "cron.go"):     cronSchedulerGo(),
		filepath.Join(apiRoot, "internal", "cron", "start.go"):    cronStartGo(),
		filepath.Join(apiRoot, "internal", "handlers", "cron.go"): cronHandlerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func cronSchedulerGo() string {
	return `package cron

import (
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// Task represents a registered cron task for display purposes.
type Task struct {
	Name     string ` + "`" + `json:"name"` + "`" + `
	Schedule string ` + "`" + `json:"schedule"` + "`" + `
	Type     string ` + "`" + `json:"type"` + "`" + `
}

// RegisteredTasks holds the list of cron tasks for the admin API.
var RegisteredTasks []Task

// Scheduler wraps asynq.Scheduler for cron-like job scheduling.
type Scheduler struct {
	scheduler *asynq.Scheduler
}

// New creates a new cron Scheduler connected to Redis.
func New(redisURL string) (*Scheduler, error) {
	redisOpt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL for cron: %w", err)
	}

	scheduler := asynq.NewScheduler(redisOpt, nil)

	// Register built-in cron tasks
	RegisteredTasks = []Task{}

	// Cleanup expired tokens — every hour
	_, err = scheduler.Register("0 * * * *", asynq.NewTask("tokens:cleanup", nil))
	if err != nil {
		return nil, fmt.Errorf("registering tokens cleanup: %w", err)
	}
	RegisteredTasks = append(RegisteredTasks, Task{
		Name:     "Cleanup expired tokens",
		Schedule: "0 * * * *",
		Type:     "tokens:cleanup",
	})

	// grit:cron-tasks

	return &Scheduler{scheduler: scheduler}, nil
}

// Start begins executing scheduled tasks.
func (s *Scheduler) Start() error {
	go func() {
		if err := s.scheduler.Run(); err != nil {
			log.Printf("Cron scheduler error: %v", err)
		}
	}()
	return nil
}

// Stop shuts down the scheduler gracefully.
func (s *Scheduler) Stop() {
	s.scheduler.Shutdown()
}
`
}

// cronStartGo returns a batteries-included cron.Start(cfg, cache) helper
// that wires up the asynq scheduler. This is what `cmd/server/main.go` for
// the single-app architecture imports — it lets the main.go stay tight and
// gives `grit add cron-task` a place to inject new tasks.
func cronStartGo() string {
	return `package cron

import (
	"log"

	"{{MODULE}}/internal/cache"
	"{{MODULE}}/internal/config"
)

// Start initializes the cron scheduler with the project's default tasks and
// begins executing scheduled jobs in the background.
//
// Returns the running Scheduler (so callers can Stop() it on shutdown) and
// any startup error. The cache argument is reserved for tasks that need
// shared cache access — it is currently unused but kept stable so adding a
// cache-backed task later is a drop-in change.
func Start(cfg *config.Config, _ *cache.Cache) (*Scheduler, error) {
	if cfg.RedisURL == "" {
		log.Println("Cron scheduler disabled: REDIS_URL not configured")
		return nil, nil
	}

	s, err := New(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	if err := s.Start(); err != nil {
		return nil, err
	}

	log.Println("Cron scheduler started")
	return s, nil
}
`
}

func cronHandlerGo() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/cron"
)

// CronHandler handles admin cron task endpoints.
type CronHandler struct{}

// ListTasks returns all registered cron tasks.
func (h *CronHandler) ListTasks(c *gin.Context) {
	tasks := cron.RegisteredTasks
	if tasks == nil {
		tasks = []cron.Task{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}
`
}
