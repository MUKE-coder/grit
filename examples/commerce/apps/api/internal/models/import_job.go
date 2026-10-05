package models

import (
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// ImportJob records the progress of a background CSV import started by
// POST /<resource>/import. The handler inserts one row, launches a goroutine
// that streams the file, and updates Processed/Created/Skipped/Failed as it
// runs. Clients poll GET /imports/:id to drive a live progress bar, then read
// the final counts and per-row errors when Status is "completed". Errors holds
// up to the first 50 row failures as a JSON array string.
type ImportJob struct {
	ID        string `gorm:"primarykey;size:36" json:"id"`
	Resource  string `gorm:"size:255;index" json:"resource"`
	Status    string `gorm:"size:20;index" json:"status"` // processing | completed | failed
	Total     int    `json:"total"`
	Processed int    `json:"processed"`
	Created   int    `json:"created"`
	Skipped   int    `json:"skipped"`
	Failed    int    `json:"failed"`
	Errors    string `gorm:"type:text" json:"-"`
	Message   string `gorm:"size:500" json:"message"`
	// CreatedBy is the user who started the import. Only they, or an admin, may
	// read the job: its row errors quote the file's contents.
	CreatedBy string    `gorm:"size:36;index" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate assigns a UUID so the job id is opaque in poll URLs.
func (m *ImportJob) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Status == "" {
		m.Status = "processing"
	}
	return nil
}
