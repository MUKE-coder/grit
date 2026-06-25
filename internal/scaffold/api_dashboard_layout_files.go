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
				UserID:       userID.(string),
				Cards:        datatypes.JSONSlice[string]{},
				Charts:       datatypes.JSONSlice[string]{},
				Tables:       datatypes.JSONSlice[string]{},
				Resources:    datatypes.JSONSlice[string]{},
				SectionOrder: datatypes.JSONSlice[string]{},
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

	var req struct {
		Cards        []string ` + "`json:\"cards\"`" + `
		Charts       []string ` + "`json:\"charts\"`" + `
		Tables       []string ` + "`json:\"tables\"`" + `
		Resources    []string ` + "`json:\"resources\"`" + `
		SectionOrder []string ` + "`json:\"section_order\"`" + `
		DatePreset   string   ` + "`json:\"date_preset\"`" + `
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
