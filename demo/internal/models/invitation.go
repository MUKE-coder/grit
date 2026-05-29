package models

import "time"

type Invitation struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	BusinessID uint       `gorm:"not null;index" json:"business_id"`
	BranchID   *uint      `json:"branch_id"`
	Email      string     `gorm:"not null" json:"email"`
	Role       string     `gorm:"not null" json:"role"`
	Token      string     `gorm:"uniqueIndex;not null" json:"-"`
	AcceptedAt *time.Time `json:"accepted_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	Business   Business   `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	Branch     *Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
