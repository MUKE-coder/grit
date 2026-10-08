package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"saas/apps/api/internal/ids"
)

// DashboardLayout stores the per-user customisation of the admin
// dashboard. One row per user (UserID is unique). The three slice
// columns hold widget keys -- the frontend resolves them against a
// catalog built from resources[].dashboard.widgets plus a fixed set
// of system widgets.
//
// Missing row = "show defaults". This means a fresh user gets the
// out-of-the-box dashboard without us having to seed a layout per
// user at registration. The PUT handler creates the row on first
// save.
//
// DatePreset persists the dashboard-wide date filter selection so a
// user who picks "Last 7 days" keeps that view across browser tabs
// and devices. Empty = "All time" (no filter).
type DashboardLayout struct {
	ID     string                      `gorm:"primarykey;size:36" json:"id"`
	UserID string                      `gorm:"size:36;uniqueIndex" json:"user_id"`
	Cards  datatypes.JSONSlice[string] `gorm:"type:json" json:"cards"`
	Charts datatypes.JSONSlice[string] `gorm:"type:json" json:"charts"`
	Tables datatypes.JSONSlice[string] `gorm:"type:json" json:"tables"`
	// v3.31.45 -- Resources holds enabled keys for the "By resource"
	// band. Convention: "<slug>:total" and "<slug>:latest" per
	// resource. Same semantics as Cards/Charts/Tables: empty list +
	// non-empty ID means "user hid everything"; missing row means
	// "show defaults". SectionOrder holds the section keys in render
	// order; default ["cards","charts","tables","by-resource"].
	Resources    datatypes.JSONSlice[string] `gorm:"type:json" json:"resources"`
	SectionOrder datatypes.JSONSlice[string] `gorm:"type:json" json:"section_order"`
	// v3.31.46 -- per-resource layout mode for the By Resource band.
	// Keys are resource slugs; values are "split" (Total left,
	// Latest right, the v3.31.44 default) or "tabs" (each widget
	// full-width in its own tab). Missing slugs fall back to
	// "split" at render time, so only non-default choices need to
	// be persisted.
	ResourceLayouts datatypes.JSON `gorm:"type:json" json:"resource_layouts"`
	// v3.31.47 -- Preset Chart builder. Array of user-defined chart
	// configurations rendered in the Charts section. Each entry has
	// resource slug + preset + field + viz + title. The handler
	// validates each entry on write so a single malformed chart
	// doesn't block the rest of the save.
	CustomCharts datatypes.JSON `gorm:"type:json" json:"custom_charts"`
	DatePreset   string         `gorm:"size:16" json:"date_preset"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (d *DashboardLayout) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = ids.New()
	}
	return nil
}
