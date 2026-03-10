package scaffold

func apiAuthTestGo() string {
	return `package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// newTestDB opens a fresh in-memory SQLite database and migrates the User model.
// It accepts testing.TB so both tests and benchmarks can call it.
func newTestDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(tb, err)
	require.NoError(tb, db.AutoMigrate(&models.User{}))
	return db
}

func newTestAuthSvc(cfg *config.Config) *services.AuthService {
	return &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
	}
}

func testCfg() *config.Config {
	return &config.Config{
		JWTSecret:        "test-secret-key-for-testing-only",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
	}
}

func newAuthRouter(h *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/auth/register", h.Register)
	r.POST("/api/auth/login", h.Login)
	return r
}

func postJSON(tb testing.TB, r *gin.Engine, path string, body map[string]string) *httptest.ResponseRecorder {
	tb.Helper()
	raw, err := json.Marshal(body)
	require.NoError(tb, err)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthHandler_Register_Success(t *testing.T) {
	cfg := testCfg()
	db := newTestDB(t)
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg), Config: cfg}
	r := newAuthRouter(h)

	w := postJSON(t, r, "/api/auth/register", map[string]string{
		"first_name": "Jane",
		"last_name":  "Doe",
		"email":      "jane@example.com",
		"password":   "password123",
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "User registered successfully", resp["message"])

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "response should have data field")
	assert.NotNil(t, data["tokens"])
	assert.NotNil(t, data["user"])
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	cfg := testCfg()
	db := newTestDB(t)
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg), Config: cfg}
	r := newAuthRouter(h)

	// Missing required fields (first_name, last_name, password)
	w := postJSON(t, r, "/api/auth/register", map[string]string{
		"email": "incomplete@example.com",
	})

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj, ok := resp["error"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "VALIDATION_ERROR", errObj["code"])
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	cfg := testCfg()
	db := newTestDB(t)
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg), Config: cfg}
	r := newAuthRouter(h)

	body := map[string]string{
		"first_name": "Alice",
		"last_name":  "Smith",
		"email":      "alice@example.com",
		"password":   "password123",
	}

	// First registration — should succeed
	w1 := postJSON(t, r, "/api/auth/register", body)
	require.Equal(t, http.StatusCreated, w1.Code)

	// Second registration same email — should conflict
	w2 := postJSON(t, r, "/api/auth/register", body)
	assert.Equal(t, http.StatusConflict, w2.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp))
	errObj, ok := resp["error"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "EMAIL_EXISTS", errObj["code"])
}

func TestAuthHandler_Login_Success(t *testing.T) {
	cfg := testCfg()
	db := newTestDB(t)
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg), Config: cfg}
	r := newAuthRouter(h)

	// Register first
	postJSON(t, r, "/api/auth/register", map[string]string{
		"first_name": "Bob",
		"last_name":  "Jones",
		"email":      "bob@example.com",
		"password":   "mypassword99",
	})

	// Then login
	w := postJSON(t, r, "/api/auth/login", map[string]string{
		"email":    "bob@example.com",
		"password": "mypassword99",
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Logged in successfully", resp["message"])

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	tokens, ok := data["tokens"].(map[string]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, tokens["access_token"])
	assert.NotEmpty(t, tokens["refresh_token"])
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	cfg := testCfg()
	db := newTestDB(t)
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg), Config: cfg}
	r := newAuthRouter(h)

	postJSON(t, r, "/api/auth/register", map[string]string{
		"first_name": "Carol",
		"last_name":  "White",
		"email":      "carol@example.com",
		"password":   "correct-password",
	})

	w := postJSON(t, r, "/api/auth/login", map[string]string{
		"email":    "carol@example.com",
		"password": "wrong-password",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj, ok := resp["error"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "INVALID_CREDENTIALS", errObj["code"])
}

func TestAuthHandler_Login_UnknownEmail(t *testing.T) {
	cfg := testCfg()
	db := newTestDB(t)
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg), Config: cfg}
	r := newAuthRouter(h)

	w := postJSON(t, r, "/api/auth/login", map[string]string{
		"email":    "nobody@example.com",
		"password": "anything",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
`
}

func apiUserTestGo() string {
	return `package handlers

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

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/middleware"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

func newUserTestSetup(t *testing.T) (*gin.Engine, *services.AuthService, *config.Config) {
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

	return r, authSvc, cfg
}

func TestUserHandler_List_RequiresAuth(t *testing.T) {
	r, _, _ := newUserTestSetup(t)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_List_AdminAccess(t *testing.T) {
	r, authSvc, _ := newUserTestSetup(t)

	// Fetch the admin we seeded — just generate a token directly
	tokens, err := authSvc.GenerateTokenPair(1, "admin@example.com", models.RoleAdmin)
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
	r, authSvc, _ := newUserTestSetup(t)

	tokens, err := authSvc.GenerateTokenPair(1, "admin@example.com", models.RoleAdmin)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/users/99999", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	r, authSvc, _ := newUserTestSetup(t)

	tokens, err := authSvc.GenerateTokenPair(1, "admin@example.com", models.RoleAdmin)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp["data"])
}
`
}

func apiBenchTestGo() string {
	return `package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/services"
)

func newBenchRouter(b *testing.B) *gin.Engine {
	b.Helper()
	gin.SetMode(gin.TestMode)

	db := newTestDB(b)
	cfg := &config.Config{
		JWTSecret:        "bench-secret-key",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
	}
	authSvc := &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
	}
	h := &AuthHandler{DB: db, AuthService: authSvc, Config: cfg}

	r := gin.New()
	r.POST("/api/auth/register", h.Register)
	r.POST("/api/auth/login", h.Login)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

// BenchmarkHealthCheck measures the raw HTTP overhead for a minimal JSON endpoint.
func BenchmarkHealthCheck(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

// BenchmarkAuthLogin measures login throughput including JWT generation.
// A user is pre-registered before the timer starts.
func BenchmarkAuthLogin(b *testing.B) {
	r := newBenchRouter(b)

	body, _ := json.Marshal(map[string]string{
		"first_name": "Login",
		"last_name":  "Bench",
		"email":      "loginbench@example.com",
		"password":   "benchpassword99",
	})
	req0 := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req0.Header.Set("Content-Type", "application/json")
	w0 := httptest.NewRecorder()
	r.ServeHTTP(w0, req0)
	require.Equal(b, http.StatusCreated, w0.Code, "pre-registration must succeed")

	loginBody, _ := json.Marshal(map[string]string{
		"email":    "loginbench@example.com",
		"password": "benchpassword99",
	})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

// BenchmarkAuthRegister measures full registration throughput including bcrypt hashing.
// Note: each iteration uses a unique email to avoid UNIQUE constraint failures.
func BenchmarkAuthRegister(b *testing.B) {
	gin.SetMode(gin.TestMode)
	db := newTestDB(b)
	cfg := &config.Config{
		JWTSecret:        "bench-secret-key",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
	}
	h := &AuthHandler{
		DB: db,
		AuthService: &services.AuthService{
			Secret:        cfg.JWTSecret,
			AccessExpiry:  cfg.JWTAccessExpiry,
			RefreshExpiry: cfg.JWTRefreshExpiry,
		},
		Config: cfg,
	}
	r := gin.New()
	r.POST("/api/auth/register", h.Register)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		email := "user" + fmt.Sprint(i) + "@bench.dev"
		body, _ := json.Marshal(map[string]string{
			"first_name": "Bench",
			"last_name":  "User",
			"email":      email,
			"password":   "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
`
}
