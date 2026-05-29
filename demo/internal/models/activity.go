package models

import (
	"time"

	"gorm.io/datatypes"
)

// Activity is an immutable audit-log entry. Every meaningful user action
// (login, sale created, loan approved, etc.) creates one row. We never
// update or delete activities — they're append-only history.
//
// Description is the human-readable, already-formatted story of what
// happened ("Sold 3 items for UGX 52,500"). Metadata holds machine-
// readable structured context for future drill-downs.
type Activity struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	BusinessID  uint           `gorm:"not null;index:idx_activity_biz_created" json:"business_id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Action      string         `gorm:"type:varchar(64);not null;index" json:"action"`     // login, create, update, delete, approve, verify, etc.
	Resource    string         `gorm:"type:varchar(64);not null;index" json:"resource"`   // sale, loan, product, etc.
	ResourceID  *uint          `gorm:"index" json:"resource_id,omitempty"`                // nullable — login has no resource id
	Description string         `gorm:"type:text" json:"description"`
	Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	IPAddress   string         `gorm:"type:varchar(64)" json:"ip_address,omitempty"`
	UserAgent   string         `gorm:"type:varchar(512)" json:"user_agent,omitempty"`
	CreatedAt   time.Time      `gorm:"index:idx_activity_biz_created" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Common action codes — keep these stable; the frontend uses them for
// icon/color mapping.
const (
	ActionLogin    = "login"
	ActionLogout   = "logout"
	ActionCreate   = "create"
	ActionUpdate   = "update"
	ActionDelete   = "delete"
	ActionApprove  = "approve"
	ActionReject   = "reject"
	ActionVerify   = "verify"
	ActionDisburse = "disburse"
	ActionTransfer = "transfer"
	ActionImport   = "import"
	ActionReturn   = "return"
)
