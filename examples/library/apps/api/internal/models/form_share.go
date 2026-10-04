package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"library/apps/api/internal/ids"
)

// FormShare is a public link to a resource's create form. Operators
// generate one of these in the admin to expose a single Grit resource
// (e.g. Contact) without requiring auth — useful for lead forms,
// applications, public submissions.
//
// Optional bcrypt password adds a gate: if PasswordHash is set, the
// public submit endpoint requires the visitor to verify the password
// first.
type FormShare struct {
	ID           string `gorm:"primarykey;size:36" json:"id"`
	ResourceName string `gorm:"size:64;not null;index" json:"resource_name"`
	// Token is the URL-safe identifier (32 chars). Used as
	// /forms/<token> in the public web app and
	// /api/public/forms/:token/* in the API.
	Token string `gorm:"size:64;not null;uniqueIndex" json:"token"`
	// PasswordHash is empty for open-access shares. When set (bcrypt
	// cost 10), visitors must POST the plaintext password to
	// /check-password before /submit succeeds.
	PasswordHash string `gorm:"size:255" json:"-"`
	HasPassword  bool   `gorm:"-" json:"has_password"` // computed, not stored
	// not null, but no default: a default made Enabled:false unstorable on
	// create, so a share link meant to start disabled went live instead.
	Enabled         bool   `gorm:"not null" json:"enabled"`
	SubmissionCount int    `gorm:"not null;default:0" json:"submission_count"`
	CreatedByUserID string `gorm:"size:36;index" json:"created_by_user_id"`
	Label           string `gorm:"size:200" json:"label"` // optional operator-facing label

	// v3.31.50 — operator-customisable surface for the public form.
	// CustomTitle / CustomDescription replace the default heading +
	// subtitle when set; HiddenFields is a list of field keys
	// (matching the json tag on the model) to omit from the
	// rendered form -- useful for optional columns the operator
	// doesn't want anonymous visitors filling in.
	CustomTitle       string                      `gorm:"size:200" json:"custom_title"`
	CustomDescription string                      `gorm:"size:500" json:"custom_description"`
	HiddenFields      datatypes.JSONSlice[string] `gorm:"type:json" json:"hidden_fields"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *FormShare) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}

// AfterFind computes the HasPassword virtual field so the admin UI can
// show a lock icon without exposing the hash itself.
func (s *FormShare) AfterFind(tx *gorm.DB) error {
	s.HasPassword = s.PasswordHash != ""
	return nil
}
