package models

import (
	"encoding/json"
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// Role is a named bag of permission grants.
//
// Grants hold whatever the operator authored, wildcards included — see
// internal/authz for the key format and matcher. Storing "products.*" rather
// than the expanded leaves is deliberate: the role keeps working when a new
// action shows up in the catalog.
type Role struct {
	ID          string `gorm:"primarykey;size:36" json:"id"`
	Name        string `gorm:"size:80;uniqueIndex;not null" json:"name" binding:"required"`
	Description string `gorm:"size:500" json:"description"`

	// Grants is a JSON array of permission keys. Use GrantsList/SetGrants
	// rather than touching it directly.
	Grants string `gorm:"type:text" json:"-"`

	// IsSystem marks a role the app seeds and depends on. System roles can be
	// edited but not deleted or renamed — enforced in the handler, not just the
	// UI, so a direct API call can't rename ADMIN out from under the guards.
	IsSystem bool `gorm:"default:false;index" json:"is_system"`

	Version   int            `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ids.New()
	}
	return nil
}

func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// GrantsList decodes the stored grants.
//
// A malformed value yields NO permissions rather than an error: this is an
// authorization path, and the safe failure mode is to deny.
func (r *Role) GrantsList() []string {
	if r.Grants == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(r.Grants), &out); err != nil {
		return nil
	}
	return out
}

func (r *Role) SetGrants(grants []string) error {
	if grants == nil {
		grants = []string{}
	}
	b, err := json.Marshal(grants)
	if err != nil {
		return err
	}
	r.Grants = string(b)
	return nil
}

// MarshalJSON exposes grants as a real array to API clients while keeping the
// column a string internally.
func (r Role) MarshalJSON() ([]byte, error) {
	type alias Role
	return json.Marshal(struct {
		alias
		Grants []string `json:"grants"`
	}{alias(r), r.GrantsList()})
}

// UserRole assigns a Role to a User.
//
// No DeletedAt on purpose — see the note on roleModelGo. Re-granting a
// previously removed role must be idempotent, and a tombstone row makes the
// composite PK reject it.
type UserRole struct {
	UserID    string    `gorm:"primaryKey;size:36;index" json:"user_id"`
	RoleID    string    `gorm:"primaryKey;size:36;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string { return "user_roles" }

// DefaultRoles are seeded on first boot.
//
// The grants mirror what the routes enforced BEFORE permissions existed, so
// upgrading an app doesn't change who can do what:
//
//	ADMIN  — everything, via the "*" superuser grant
//	EDITOR — content, matching the role's existing meaning
//	USER   — nothing privileged; ordinary users hold no admin permissions
//
// Note this is NOT "every role gets every permission". The route guard passes if
// the caller matches the legacy role name OR holds the permission, so granting
// USER the full catalog would hand every signed-up user the admin panel.
func DefaultRoles() []Role {
	return []Role{
		{
			Name:        RoleAdmin,
			Description: "Full access to every feature and action.",
			IsSystem:    true,
			Grants:      `["*"]`,
		},
		{
			Name:        RoleEditor,
			Description: "Manages content and uploads. No user or system administration.",
			IsSystem:    true,
			Grants:      `["uploads.create","uploads.view","uploads.delete","users.view"]`,
		},
		{
			Name:        RoleUser,
			Description: "Standard account. No administrative permissions.",
			IsSystem:    true,
			Grants:      `[]`,
		},
	}
}

// SeedRoles inserts the default roles if they are missing.
//
// Idempotent, and deliberately NON-destructive: an existing role's grants are
// left alone. Re-applying defaults on every boot would silently revert an
// operator's edits — a real bug in the system this was modelled on, where
// seeding ran on every roles-page load and quietly undid their changes.
// "Reset to defaults" is an explicit action in the UI instead.
func SeedRoles(db *gorm.DB) error {
	for _, want := range DefaultRoles() {
		var existing Role
		// Find, not First: on a fresh database every role is missing, and First
		// logs each miss as "record not found" — three scary-looking lines on
		// the first migrate even though creating them is exactly the point.
		// Find treats zero rows as normal and stays quiet.
		if err := db.Where("name = ?", want.Name).Limit(1).Find(&existing).Error; err != nil {
			return err
		}
		if existing.ID != "" {
			continue // already present — leave the operator's grants intact
		}
		role := want
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}
	return nil
}
