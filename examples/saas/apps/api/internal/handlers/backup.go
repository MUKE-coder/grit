package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/backup"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/storage"
)

const (
	// manualBackupWindow rate-limits on-demand backups. The weekly cron bypasses it.
	manualBackupWindow = 24 * time.Hour
	// downloadURLTTL: long enough for a slow phone, short enough that a leaked
	// link stops working before anyone can use it.
	downloadURLTTL = 15 * time.Minute
	// backupTimeout bounds a single run so a hung upload can't wedge the worker.
	backupTimeout = 30 * time.Minute
)

// BackupHandler serves the full-database backup index.
type BackupHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
}

func (h *BackupHandler) svc() *backup.Service {
	return &backup.Service{DB: h.DB, Storage: h.Storage}
}

// List returns backups newest-first. Poll it while one is RUNNING.
func (h *BackupHandler) List(c *gin.Context) {
	var items []models.Backup
	if err := h.DB.WithContext(c.Request.Context()).Order("created_at desc").Limit(50).Find(&items).Error; err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to list backups")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// Generate starts a manual backup in the background and returns the RUNNING row
// immediately — a full dump can take a while. Poll List until it flips to READY.
func (h *BackupHandler) Generate(c *gin.Context) {
	if h.svc().Store() == nil {
		respond.Fail(c, respond.CodeStorageUnavailable, "Object storage is not configured")
		return
	}

	limited, err := h.svc().ManualRateLimited(manualBackupWindow)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to check rate limit")
		return
	}
	if limited {
		respond.Fail(c, respond.CodeRateLimited, "A manual backup was already taken in the last 24 hours")
		return
	}

	rec, err := h.svc().Start("MANUAL")
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to start backup")
		return
	}

	go func(r models.Backup) {
		ctx, cancel := context.WithTimeout(context.Background(), backupTimeout)
		defer cancel()
		_ = h.svc().Run(ctx, &r)
	}(*rec)

	c.JSON(http.StatusAccepted, gin.H{"data": rec, "message": "Backup started"})
}

// Download mints a short-lived pre-signed URL so the client pulls the archive
// straight from object storage — no proxying a multi-hundred-MB file through the API.
func (h *BackupHandler) Download(c *gin.Context) {
	if h.svc().Store() == nil {
		respond.Fail(c, respond.CodeStorageUnavailable, "Object storage is not configured")
		return
	}

	var b models.Backup
	if err := h.DB.WithContext(c.Request.Context()).First(&b, "id = ?", c.Param("id")).Error; err != nil {
		respond.Fail(c, respond.CodeNotFound, "Backup not found")
		return
	}
	if b.Status != "READY" || b.StorageKey == "" {
		respond.Fail(c, respond.CodeNotAvailable, "This backup is not available for download")
		return
	}

	// On the backups disk, or on the default store for an archive taken before
	// STORAGE_DISKS named one.
	url, err := h.svc().Locate(c.Request.Context(), b.StorageKey).GetSignedURL(c.Request.Context(), b.StorageKey, downloadURLTTL)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to sign download URL")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"url": url, "expires_in": int(downloadURLTTL.Seconds())},
	})
}

// GetSettings returns the automatic-backup schedule (frequency, time, enabled).
func (h *BackupHandler) GetSettings(c *gin.Context) {
	sc, err := h.svc().GetSchedule()
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to load backup schedule")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sc})
}

// The automatic-backup schedule.
type BackupSettingsRequest struct {
	Frequency string `json:"frequency"`
	Time      string `json:"time"`
	Enabled   bool   `json:"enabled"`
}

// UpdateSettings changes the automatic-backup schedule. The scheduler picks up
// the new period on its next tick — no restart needed.

func (h *BackupHandler) UpdateSettings(c *gin.Context) {
	var req BackupSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeInvalidBody, err.Error())
		return
	}
	sc, err := h.svc().SaveSchedule(req.Frequency, req.Time, req.Enabled)
	if err != nil {
		respond.Fail(c, respond.CodeInvalidSchedule, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sc, "message": "Backup schedule updated"})
}
