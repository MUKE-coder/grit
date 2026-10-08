package models

import (
	"time"

	"gorm.io/gorm"
	"saas/apps/api/internal/ids"
)

// UserActivity is a semantic event written by handlers via
// activity.Log(...). Surfaced in the /system/activity dashboard so
// operators can see who did what, when, from where, and how severe.
//
// Separate from ActivityLog (the hash-chained HTTP audit) on purpose:
//   - This is human-readable, queryable, filter-friendly.
//   - ActivityLog is tamper-evident, immutable, mostly compliance fodder.
//
// Use the Severity column as the operator-facing impact level:
//   - info     → routine action (login, view, list)
//   - warn     → something noteworthy (failed login, role change)
//   - critical → high-impact action (user delete, billing change, mass
//     update). Surface these prominently in the dashboard.
type UserActivity struct {
	ID           string         `gorm:"primarykey;size:36" json:"id"`
	UserID       string         `gorm:"size:36;index" json:"user_id"`       // actor (empty = system)
	Action       string         `gorm:"size:64;index" json:"action"`        // dotted: "ticket.create", "user.delete"
	Severity     string         `gorm:"size:16;index" json:"severity"`      // info | warn | critical
	Summary      string         `gorm:"size:500" json:"summary"`            // single sentence for the dashboard row
	ResourceType string         `gorm:"size:64;index" json:"resource_type"` // e.g. "user", "ticket"
	ResourceID   string         `gorm:"size:64;index" json:"resource_id"`   // primary key of the target
	IPAddress    string         `gorm:"size:45" json:"ip_address"`
	UserAgent    string         `gorm:"size:500" json:"user_agent"`
	Metadata     string         `gorm:"type:text" json:"metadata"` // optional JSON blob for extra context
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *UserActivity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = ids.New()
	}
	if a.Severity == "" {
		a.Severity = "info"
	}
	return nil
}
