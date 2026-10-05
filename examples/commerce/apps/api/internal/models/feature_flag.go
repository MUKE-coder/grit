package models

import (
	"encoding/json"
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FeatureFlag is one rollout switch. Two flavors:
//   - Boolean (Variants empty) → IsEnabled returns true/false
//   - A/B (Variants set)       → Variant returns one of the listed
//     strings, sticky per (user, flag).
//
// Rules JSON shape (FlagRules): rollout_percentage, allowlist_user_ids,
// blocklist_user_ids, enabled_from, enabled_until, variants. The
// percentage and variant assignment both bucket users by
// SHA-256(user_id || ":" || flag_name) % 100 so the same user always
// lands in the same slot for a given flag — no flicker between sessions.
type FeatureFlag struct {
	ID          string         `gorm:"primarykey;size:36" json:"id"`
	Name        string         `gorm:"size:100;uniqueIndex;not null" json:"name"` // e.g. "new_dashboard"
	Description string         `gorm:"type:text" json:"description"`
	Enabled     bool           `gorm:"not null;default:false" json:"enabled"` // master switch — false short-circuits all rules
	Rules       datatypes.JSON `json:"rules"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Version     int            `gorm:"not null;default:1" json:"version"`
}

func (f *FeatureFlag) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = ids.New()
	}
	return nil
}

func (f *FeatureFlag) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// FlagRules is the structured form of FeatureFlag.Rules. Use
// (*FeatureFlag).ParsedRules() to decode; (*FeatureFlag).SetRules() to
// encode + assign.
type FlagRules struct {
	RolloutPercentage int        `json:"rollout_percentage,omitempty"` // 0..100; 0 = off, 100 = full rollout
	AllowlistUserIDs  []string   `json:"allowlist_user_ids,omitempty"` // when non-empty, ONLY these users get the flag
	BlocklistUserIDs  []string   `json:"blocklist_user_ids,omitempty"` // always-deny set; runs before allowlist + percentage
	EnabledFrom       *time.Time `json:"enabled_from,omitempty"`       // before this, flag is off (date window)
	EnabledUntil      *time.Time `json:"enabled_until,omitempty"`      // after this, flag is off
	Variants          []string   `json:"variants,omitempty"`           // when set, A/B mode — Variant() returns one of these

	// Attributes restrict the flag to subjects whose attributes match:
	// {"business_unit": ["eu", "ke"]} is on only for a subject whose
	// business_unit is eu or ke. Every key listed must match, ignoring case.
	// A subject without the attribute does not match, so a rule on something
	// the app never supplies fails closed. See flags.AttributesFor.
	Attributes map[string][]string `json:"attributes,omitempty"`
}

// ParsedRules decodes the Rules JSON. Returns a zero FlagRules on
// missing or malformed JSON — callers shouldn't error out for
// misconfigured flags; they should fail closed (return false).
func (f *FeatureFlag) ParsedRules() FlagRules {
	var r FlagRules
	if len(f.Rules) > 0 {
		_ = json.Unmarshal(f.Rules, &r)
	}
	return r
}

// SetRules encodes a FlagRules and assigns it. Errors propagate.
func (f *FeatureFlag) SetRules(r FlagRules) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	f.Rules = b
	return nil
}

// FlagExposure records that a user was checked against a flag and what
// outcome they got. Used by the admin UI to show rollout health
// ("4,231 users saw variant_a, 4,189 saw variant_b") and to power
// downstream A/B analytics joins.
//
// Insert is fire-and-forget — exposure tracking should never block a
// flag check. We persist async in a goroutine.
type FlagExposure struct {
	ID        string    `gorm:"primarykey;size:36" json:"id"`
	FlagID    string    `gorm:"size:36;index;not null" json:"flag_id"`
	FlagName  string    `gorm:"size:100;index" json:"flag_name"` // denormalized for join-free analytics
	UserID    string    `gorm:"size:36;index" json:"user_id"`
	Variant   string    `gorm:"size:50" json:"variant"` // "enabled" / "disabled" / "control" / "variant_a" / etc.
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (e *FlagExposure) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = ids.New()
	}
	return nil
}
