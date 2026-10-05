package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
)

// APIKeyHandler manages a user's own machine credentials.
type APIKeyHandler struct {
	DB *gorm.DB
}

// List returns the caller's keys. The secret is not among them — there is
// nothing to return, since only its hash was ever stored.
func (h *APIKeyHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	var keys []models.APIKey
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID).
		Order("created_at desc").Find(&keys).Error; err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to load API keys")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": keys})
}

// What a caller sends to mint an API key.
type CreateAPIKeyRequest struct {
	Name string `json:"name" binding:"required,min=1,max=120"`
	// Kind is "publishable" or "secret". Empty means secret: a caller that has
	// not thought about it should get the careful one.
	Kind      string   `json:"kind"`
	Scopes    []string `json:"scopes"`
	Endpoints []string `json:"endpoints"`
	Origins   []string `json:"origins"`
	RateLimit int      `json:"rate_limit"`
	ExpiresIn int      `json:"expires_in_days"`
}

// Create issues a key and returns it once.

func (h *APIKeyHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		t := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &t
	}

	issued, err := services.GenerateAPIKey(h.DB, services.KeyOptions{
		UserID:    userID,
		Name:      req.Name,
		Kind:      req.Kind,
		Scopes:    req.Scopes,
		Endpoints: req.Endpoints,
		Origins:   req.Origins,
		RateLimit: req.RateLimit,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to create the API key")
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "api_key.create",
		Severity:     "warn",
		Summary:      "API key created: " + req.Name,
		ResourceType: "api_key",
		ResourceID:   issued.Record.ID,
	})

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"key":   issued.Record,
			"token": issued.Token,
		},
		"message": "Copy this key now: it will not be shown again.",
	})
}

// Revoke marks one of the caller's keys unusable.
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := services.RevokeAPIKey(h.DB, c.Param("id"), userID); err != nil {
		respond.Fail(c, respond.CodeNotFound, "Key not found, or already revoked")
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "api_key.revoke",
		Severity:     "warn",
		Summary:      "API key revoked",
		ResourceType: "api_key",
		ResourceID:   c.Param("id"),
	})

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked"})
}
