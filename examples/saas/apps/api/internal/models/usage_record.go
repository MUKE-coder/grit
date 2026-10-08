package models

import (
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/erasure"
	"saas/apps/api/internal/ids"
	"saas/apps/api/internal/jsontime"
)

// UsageRecord represents a usagerecord in the system.
type UsageRecord struct {
	ID             string         `gorm:"primarykey;size:36" json:"id"`
	SubscriptionID string         `gorm:"size:36;index" json:"subscription_id" binding:"required"`
	Subscription   *Subscription  `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
	Metric         string         `gorm:"size:255" json:"metric" binding:"required" search:"trigram"`
	Quantity       int            `json:"quantity"`
	RecordedOn     *jsontime.Date `gorm:"type:date" json:"recorded_on"`
	UserID         string         `gorm:"size:36;index" json:"user_id" binding:"required"`
	User           *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Version        int            `gorm:"not null;default:1" json:"version"`
	CreatedAt      time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	// ArchivedAt is the "put this away without destroying it" state, and it is
	// deliberately not DeletedAt. A soft delete is invisible to every query and
	// means the row is gone as far as the app is concerned; an archived row is
	// still listable, still exportable and still restorable in one click. The
	// list endpoint hides archived rows unless ?archived=true asks for them.
	ArchivedAt *time.Time `gorm:"index" json:"archived_at,omitempty"`
}

// BeforeCreate generates a UUID before inserting.
func (m *UsageRecord) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *UsageRecord) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// GetOwnerID identifies the row's owner for authz.MustOwnUnlessAdmin.
//
// This resource was generated with --owned-by user: a caller may only read or
// write rows where this matches their own user id, ADMIN excepted.
func (m *UsageRecord) GetOwnerID() string {
	return m.UserID
}

// init registers UsageRecord for right-to-erasure: erasing a user deletes the rows
// they own here, by user_id. See internal/erasure.
func init() { erasure.Register(&UsageRecord{}, "user_id") }
