package scaffold

// api_dashboard_layout_files.go scaffolds the v3.31.40 per-user
// dashboard customisation feature: model + handler. Routes injection
// happens in scaffold.go (the protected group), AutoMigrate
// registration happens in models.User.Models().

// dashboardLayoutModelGo emits models/dashboard_layout.go.
func dashboardLayoutModelGo() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	ID         string                      ` + "`gorm:\"primarykey;size:36\" json:\"id\"`" + `
	UserID     string                      ` + "`gorm:\"size:36;uniqueIndex\" json:\"user_id\"`" + `
	Cards      datatypes.JSONSlice[string] ` + "`gorm:\"type:json\" json:\"cards\"`" + `
	Charts     datatypes.JSONSlice[string] ` + "`gorm:\"type:json\" json:\"charts\"`" + `
	Tables     datatypes.JSONSlice[string] ` + "`gorm:\"type:json\" json:\"tables\"`" + `
	// v3.31.45 -- Resources holds enabled keys for the "By resource"
	// band. Convention: "<slug>:total" and "<slug>:latest" per
	// resource. Same semantics as Cards/Charts/Tables: empty list +
	// non-empty ID means "user hid everything"; missing row means
	// "show defaults". SectionOrder holds the section keys in render
	// order; default ["cards","charts","tables","by-resource"].
	Resources    datatypes.JSONSlice[string] ` + "`gorm:\"type:json\" json:\"resources\"`" + `
	SectionOrder datatypes.JSONSlice[string] ` + "`gorm:\"type:json\" json:\"section_order\"`" + `
	// v3.31.46 -- per-resource layout mode for the By Resource band.
	// Keys are resource slugs; values are "split" (Total left,
	// Latest right, the v3.31.44 default) or "tabs" (each widget
	// full-width in its own tab). Missing slugs fall back to
	// "split" at render time, so only non-default choices need to
	// be persisted.
	ResourceLayouts datatypes.JSON ` + "`gorm:\"type:json\" json:\"resource_layouts\"`" + `
	// v3.31.47 -- Preset Chart builder. Array of user-defined chart
	// configurations rendered in the Charts section. Each entry has
	// resource slug + preset + field + viz + title. The handler
	// validates each entry on write so a single malformed chart
	// doesn't block the rest of the save.
	CustomCharts datatypes.JSON ` + "`gorm:\"type:json\" json:\"custom_charts\"`" + `
	DatePreset string                      ` + "`gorm:\"size:16\" json:\"date_preset\"`" + `
	CreatedAt  time.Time                   ` + "`json:\"created_at\"`" + `
	UpdatedAt  time.Time                   ` + "`json:\"updated_at\"`" + `
}

func (d *DashboardLayout) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}
`
}

// dashboardLayoutHandlerGo emits handlers/dashboard_layout.go.
func dashboardLayoutHandlerGo() string {
	return `package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// DashboardLayoutHandler exposes GET + PUT for the per-user dashboard
// customisation row. Both endpoints require auth; the userID comes
// from the gin context set by middleware.Auth.
type DashboardLayoutHandler struct {
	DB *gorm.DB
}

// Get returns the current user's saved dashboard layout. If no row
// exists yet (fresh user, never saved), returns a zero-valued layout
// with the user_id populated -- the client treats that shape as
// "show defaults" so the dashboard works out of the box without us
// having to seed a row at registration time.
func (h *DashboardLayoutHandler) Get(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "UNAUTHENTICATED", "message": "Not signed in"},
		})
		return
	}

	var layout models.DashboardLayout
	err := h.DB.Where("user_id = ?", userID).First(&layout).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Empty layout = show every widget by default.
		c.JSON(http.StatusOK, gin.H{
			"data": models.DashboardLayout{
				UserID:          userID.(string),
				Cards:           datatypes.JSONSlice[string]{},
				Charts:          datatypes.JSONSlice[string]{},
				Tables:          datatypes.JSONSlice[string]{},
				Resources:       datatypes.JSONSlice[string]{},
				SectionOrder:    datatypes.JSONSlice[string]{},
				ResourceLayouts: datatypes.JSON([]byte("{}")),
				CustomCharts:    datatypes.JSON([]byte("[]")),
			},
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to load layout"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": layout})
}

// Put replaces the current user's dashboard layout. Whole-resource
// replace (not patch) because the layout payload is small (typically
// under a few KB even with hundreds of widgets) and the semantics are
// easier to reason about -- whatever you send is what you get.
func (h *DashboardLayoutHandler) Put(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "UNAUTHENTICATED", "message": "Not signed in"},
		})
		return
	}

	// v3.31.47 -- inline struct for one custom chart entry.
	type customChartReq struct {
		ID       string ` + "`json:\"id\"`" + `
		Title    string ` + "`json:\"title\"`" + `
		Resource string ` + "`json:\"resource\"`" + `
		Preset   string ` + "`json:\"preset\"`" + `
		Field    string ` + "`json:\"field\"`" + `
		Viz      string ` + "`json:\"viz\"`" + `
		Limit    int    ` + "`json:\"limit\"`" + `
		Grain    string ` + "`json:\"grain\"`" + `
	}

	var req struct {
		Cards           []string          ` + "`json:\"cards\"`" + `
		Charts          []string          ` + "`json:\"charts\"`" + `
		Tables          []string          ` + "`json:\"tables\"`" + `
		Resources       []string          ` + "`json:\"resources\"`" + `
		SectionOrder    []string          ` + "`json:\"section_order\"`" + `
		ResourceLayouts map[string]string ` + "`json:\"resource_layouts\"`" + `
		CustomCharts    []customChartReq  ` + "`json:\"custom_charts\"`" + `
		DatePreset      string            ` + "`json:\"date_preset\"`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	// Normalise nil -> empty slice so the JSON column never carries a
	// SQL NULL. The frontend always sees [], never null.
	if req.Cards == nil {
		req.Cards = []string{}
	}
	if req.Charts == nil {
		req.Charts = []string{}
	}
	if req.Tables == nil {
		req.Tables = []string{}
	}
	if req.Resources == nil {
		req.Resources = []string{}
	}
	if req.SectionOrder == nil {
		req.SectionOrder = []string{}
	}
	// v3.31.46 -- normalise the layout map and reject unknown values
	// rather than silently storing garbage. Only "split" and "tabs"
	// are recognised; missing entries fall back to "split" at render
	// time, so the map only needs to carry non-default choices.
	if req.ResourceLayouts == nil {
		req.ResourceLayouts = map[string]string{}
	}
	for k, v := range req.ResourceLayouts {
		if v != "split" && v != "tabs" {
			delete(req.ResourceLayouts, k)
		}
	}

	// v3.31.47 -- validate each custom chart entry. Drop malformed
	// rows rather than rejecting the whole save: an operator with
	// 12 valid charts and 1 garbage entry should still be able to
	// save the 12.
	validPresets := map[string]bool{
		"count_over_time": true,
		"group_by":        true,
		"sum_over_time":   true,
		"avg_over_time":   true,
	}
	validVizes := map[string]bool{
		"bar":   true,
		"line":  true,
		"area":  true,
		"pie":   true,
		"donut": true,
	}
	cleanCharts := make([]customChartReq, 0, len(req.CustomCharts))
	for _, ch := range req.CustomCharts {
		if ch.ID == "" || ch.Resource == "" {
			continue
		}
		if !validPresets[ch.Preset] || !validVizes[ch.Viz] {
			continue
		}
		if ch.Preset != "count_over_time" && ch.Field == "" {
			continue
		}
		if ch.Limit <= 0 || ch.Limit > 100 {
			ch.Limit = 10
		}
		if ch.Grain == "" {
			ch.Grain = "day"
		}
		cleanCharts = append(cleanCharts, ch)
	}
	req.CustomCharts = cleanCharts

	var layout models.DashboardLayout
	err := h.DB.Where("user_id = ?", userID).First(&layout).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		layout = models.DashboardLayout{UserID: userID.(string)}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to load layout"},
		})
		return
	}

	layout.Cards = datatypes.NewJSONSlice(req.Cards)
	layout.Charts = datatypes.NewJSONSlice(req.Charts)
	layout.Tables = datatypes.NewJSONSlice(req.Tables)
	layout.Resources = datatypes.NewJSONSlice(req.Resources)
	layout.SectionOrder = datatypes.NewJSONSlice(req.SectionOrder)
	if b, err := json.Marshal(req.ResourceLayouts); err == nil {
		layout.ResourceLayouts = datatypes.JSON(b)
	} else {
		// Map should always marshal cleanly; on the off chance it
		// fails we fall back to an empty object so the column is
		// never invalid JSON.
		layout.ResourceLayouts = datatypes.JSON([]byte("{}"))
	}
	if b, err := json.Marshal(req.CustomCharts); err == nil {
		layout.CustomCharts = datatypes.JSON(b)
	} else {
		layout.CustomCharts = datatypes.JSON([]byte("[]"))
	}
	layout.DatePreset = req.DatePreset

	if err := h.DB.Save(&layout).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to save layout"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    layout,
		"message": "Dashboard layout saved",
	})
}
`
}
