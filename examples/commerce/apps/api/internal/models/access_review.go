package models

import (
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// AccessReview is one recertification campaign: a point-in-time snapshot of all
// role assignments, reviewed grant by grant.
type AccessReview struct {
	ID   string `gorm:"primarykey;size:36" json:"id"`
	Name string `gorm:"size:200;not null" json:"name" binding:"required"`
	Note string `gorm:"size:1000" json:"note"`

	// Status is "open" while items are being decided, "completed" once the
	// reviewer signs it off. A completed review is the artifact an auditor asks
	// for; it is never reopened — a new period is a new campaign.
	Status string `gorm:"size:16;index;default:open" json:"status"`

	CreatedBy      string `gorm:"size:36;index" json:"created_by"`
	CreatedByEmail string `gorm:"size:255" json:"created_by_email"`

	CompletedBy      string     `gorm:"size:36" json:"completed_by,omitempty"`
	CompletedByEmail string     `gorm:"size:255" json:"completed_by_email,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Items []AccessReviewItem `gorm:"foreignKey:ReviewID" json:"items,omitempty"`
}

func (r *AccessReview) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ids.New()
	}
	if r.Status == "" {
		r.Status = "open"
	}
	return nil
}

// AccessReviewItem is one role assignment under review. UserEmail and RoleName
// are snapshots taken when the campaign opens — the whole point of a review is
// that its record survives later changes to the things it reviewed.
type AccessReviewItem struct {
	ID       string `gorm:"primarykey;size:36" json:"id"`
	ReviewID string `gorm:"size:36;index;not null" json:"review_id"`

	UserID    string `gorm:"size:36;index" json:"user_id"`
	UserEmail string `gorm:"size:255" json:"user_email"`
	RoleID    string `gorm:"size:36;index" json:"role_id"`
	RoleName  string `gorm:"size:80" json:"role_name"`

	// Decision is "pending", "approved", or "revoked". A revoked item has already
	// had its grant removed — the decision and the effect are one transaction.
	Decision       string     `gorm:"size:16;index;default:pending" json:"decision"`
	DecidedBy      string     `gorm:"size:36" json:"decided_by,omitempty"`
	DecidedByEmail string     `gorm:"size:255" json:"decided_by_email,omitempty"`
	DecidedAt      *time.Time `json:"decided_at,omitempty"`
	Note           string     `gorm:"size:1000" json:"note"`

	CreatedAt time.Time `json:"created_at"`
}

func (i *AccessReviewItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = ids.New()
	}
	if i.Decision == "" {
		i.Decision = "pending"
	}
	return nil
}
