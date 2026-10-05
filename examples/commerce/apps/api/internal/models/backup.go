package models

import (
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// Backup indexes one full-database snapshot. The archive lives in object storage
// (R2 / S3 / MinIO) under the "backups/" prefix; this row is what the admin UI,
// the REST API and the CLI read. Lifecycle:
//
//	RUNNING -> READY | FAILED -> PURGED
//
// PURGED means rolling retention deleted the object but we kept the row, so the
// audit trail still shows a backup existed on that date.
type Backup struct {
	ID          string     `gorm:"primarykey;size:36" json:"id"`
	Kind        string     `gorm:"size:20;index" json:"kind"`   // WEEKLY | MANUAL | CLI
	Status      string     `gorm:"size:20;index" json:"status"` // RUNNING | READY | FAILED | PURGED
	StorageKey  string     `gorm:"size:512" json:"-"`
	SizeBytes   int64      `json:"size_bytes"`
	TableCount  int        `json:"table_count"`
	RowCount    int        `json:"row_count"`
	RowCounts   string     `gorm:"type:text" json:"-"` // JSON map of table -> rows
	Error       string     `gorm:"size:1000" json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// BeforeCreate assigns a UUID so backup ids are opaque in download URLs.
func (m *Backup) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Status == "" {
		m.Status = "RUNNING"
	}
	return nil
}

// BackupSchedule is the single-row configuration for automatic backups. The
// scheduler ticks frequently and consults this row to decide whether a backup
// is due — so the period can be changed at runtime without re-registering cron.
//
//	Frequency: daily | weekly | monthly | yearly (default weekly)
//	Time:      "HH:MM" server-local time-of-day the run should land at
//	Enabled:   master switch for automatic backups (manual backups still work)
type BackupSchedule struct {
	ID        uint   `gorm:"primarykey" json:"-"` // singleton, id = 1
	Frequency string `gorm:"size:20;default:weekly" json:"frequency"`
	Time      string `gorm:"size:5;default:02:00" json:"time"`
	// No default: it would make Enabled:false unstorable on create, so a
	// schedule created as disabled would still run backups.
	Enabled   bool      `json:"enabled"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BackupSchedule) TableName() string { return "backup_schedules" }
