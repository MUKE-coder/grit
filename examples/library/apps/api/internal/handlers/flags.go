package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/flags"
	"library/apps/api/internal/models"
	"library/apps/api/internal/paginate"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
)

// FeatureFlagHandler exposes admin-side CRUD over feature flags +
// exposure analytics. Mounted under admin/* (admin role required).
type FeatureFlagHandler struct {
	DB     *gorm.DB
	Engine *flags.Engine
}

func NewFeatureFlagHandler(db *gorm.DB, engine *flags.Engine) *FeatureFlagHandler {
	return &FeatureFlagHandler{DB: db, Engine: engine}
}

// flags is the service every method below reads and writes through. What a flag
// row is, and what may be searched or sorted by, is a decision about data: see
// internal/services/feature_flag.go. When the engine's cache is stale is not,
// which is why RefreshAndBroadcast stays here.
func (h *FeatureFlagHandler) flags() *services.FeatureFlagService {
	return &services.FeatureFlagService{DB: h.DB}
}

// List returns all flags with the standard paginate envelope.
//
//	GET /api/admin/flags
func (h *FeatureFlagHandler) List(c *gin.Context) {
	res, err := h.flags().List(c.Request.Context(), paginate.Bind(c))
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, res)
}

// FeatureFlagRequest is the request shape for create/update. Rules is taken
// as a structured object — the handler encodes to JSON before hitting
// the DB so the wire format stays consistent.
type FeatureFlagRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Enabled     bool             `json:"enabled"`
	Rules       models.FlagRules `json:"rules"`
}

// Create adds a new flag. Name must be unique.
//
//	POST /api/admin/flags
//	{ "name": "new_dashboard", "enabled": true, "rules": { "rollout_percentage": 25 } }
func (h *FeatureFlagHandler) Create(c *gin.Context) {
	var body FeatureFlagRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	if body.Name == "" {
		respond.Fail(c, respond.CodeValidationError, "name is required")
		return
	}

	flag := models.FeatureFlag{
		Name:        body.Name,
		Description: body.Description,
		Enabled:     body.Enabled,
	}
	if err := flag.SetRules(body.Rules); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	if err := h.flags().Create(c.Request.Context(), &flag); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}

	h.Engine.RefreshAndBroadcast(flag.Name)
	c.JSON(http.StatusCreated, gin.H{"data": flag})
}

// Update modifies an existing flag.
//
//	PUT /api/admin/flags/:id
func (h *FeatureFlagHandler) Update(c *gin.Context) {
	id := c.Param("id")
	flag, err := h.flags().ByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.Fail(c, respond.CodeNotFound, "flag not found")
			return
		}
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}

	var body FeatureFlagRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	// Name is immutable post-create — too easy to break consumers.
	flag.Description = body.Description
	flag.Enabled = body.Enabled
	if err := flag.SetRules(body.Rules); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	if err := h.flags().Save(c.Request.Context(), flag); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}

	h.Engine.RefreshAndBroadcast(flag.Name)
	c.JSON(http.StatusOK, gin.H{"data": flag})
}

// Delete removes a flag. The cache refreshes immediately so app code
// stops seeing it on the next check.
//
//	DELETE /api/admin/flags/:id
func (h *FeatureFlagHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	flag, err := h.flags().ByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.Fail(c, respond.CodeNotFound, "flag not found")
			return
		}
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	if err := h.flags().Delete(c.Request.Context(), flag); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	h.Engine.RefreshAndBroadcast(flag.Name)
	c.JSON(http.StatusOK, gin.H{"message": "flag deleted"})
}

// Exposures returns aggregate counts per variant for one flag —
// powers the rollout-health view in the admin UI.
//
//	GET /api/admin/flags/:id/exposures
//	→ { "data": [{ "variant": "enabled", "count": 4231 }, ...] }
func (h *FeatureFlagHandler) Exposures(c *gin.Context) {
	rows, err := h.flags().Exposures(c.Request.Context(), c.Param("id"))
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
