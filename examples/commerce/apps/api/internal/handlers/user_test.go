package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"commerce/apps/api/internal/config"
	"commerce/apps/api/internal/middleware"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/services"
)

func newUserTestSetup(t *testing.T) (*gin.Engine, *services.AuthService, *config.Config, models.User) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db := newTestDB(t)
	cfg := &config.Config{
		JWTSecret:        "test-secret-key-for-testing-only",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
	}
	authSvc := &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
	}
	userHandler := &UserHandler{DB: db}

	// Seed an admin user for tests that need auth
	admin := models.User{
		FirstName: "Admin",
		LastName:  "Test",
		Email:     "admin@example.com",
		Password:  "adminpass123",
		Role:      models.RoleAdmin,
		Active:    true,
	}
	require.NoError(t, db.Create(&admin).Error)

	r := gin.New()
	protected := r.Group("/api", middleware.Auth(db, authSvc))
	protected.GET("/users", userHandler.List)
	protected.GET("/users/:id", userHandler.GetByID)

	return r, authSvc, cfg, admin
}

func TestUserHandler_List_RequiresAuth(t *testing.T) {
	r, _, _, _ := newUserTestSetup(t)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_List_AdminAccess(t *testing.T) {
	r, authSvc, _, admin := newUserTestSetup(t)

	// Fetch the admin we seeded — just generate a token directly
	tokens, err := authSvc.GenerateTokenPair(admin.ID, admin.Email, admin.Role)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp["data"])
	assert.NotNil(t, resp["meta"])
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	r, authSvc, _, admin := newUserTestSetup(t)

	tokens, err := authSvc.GenerateTokenPair(admin.ID, admin.Email, admin.Role)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/users/00000000-0000-0000-0000-000000000000", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	r, authSvc, _, admin := newUserTestSetup(t)

	tokens, err := authSvc.GenerateTokenPair(admin.ID, admin.Email, admin.Role)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/users/"+admin.ID, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp["data"])
}
