package models

import (
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/erasure"
	"library/apps/api/internal/ids"
	"library/apps/api/internal/jsontime"
)

// Loan represents a loan in the system.
type Loan struct {
	ID         string             `gorm:"primarykey;size:36" json:"id"`
	BookID     string             `gorm:"size:36;index" json:"book_id" binding:"required"`
	Book       *Book              `gorm:"foreignKey:BookID" json:"book,omitempty"`
	BorrowerID string             `gorm:"size:36;index" json:"borrower_id" binding:"required"`
	Borrower   *User              `gorm:"foreignKey:BorrowerID" json:"borrower,omitempty"`
	BorrowedAt *jsontime.DateTime `json:"borrowed_at"`
	DueAt      *jsontime.DateTime `json:"due_at"`
	ReturnedAt *jsontime.DateTime `json:"returned_at"`
	Status     string             `gorm:"size:255" json:"status" binding:"required"`
	Version    int                `gorm:"not null;default:1" json:"version"`
	CreatedAt  time.Time          `gorm:"index" json:"created_at"`
	UpdatedAt  time.Time          `gorm:"index" json:"updated_at"`
	DeletedAt  gorm.DeletedAt     `gorm:"index" json:"-"`
	// ArchivedAt is the "put this away without destroying it" state, and it is
	// deliberately not DeletedAt. A soft delete is invisible to every query and
	// means the row is gone as far as the app is concerned; an archived row is
	// still listable, still exportable and still restorable in one click. The
	// list endpoint hides archived rows unless ?archived=true asks for them.
	ArchivedAt *time.Time `gorm:"index" json:"archived_at,omitempty"`
}

// BeforeCreate generates a UUID before inserting.
func (m *Loan) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Loan) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// GetOwnerID identifies the row's owner for authz.MustOwnUnlessAdmin.
//
// This resource was generated with --owned-by borrower: a caller may only read or
// write rows where this matches their own user id, ADMIN excepted.
func (m *Loan) GetOwnerID() string {
	return m.BorrowerID
}

// init registers Loan for right-to-erasure: erasing a user deletes the rows
// they own here, by borrower_id. See internal/erasure.
func init() { erasure.Register(&Loan{}, "borrower_id") }
