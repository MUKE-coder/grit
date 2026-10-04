package handlers

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

	"library/apps/api/internal/config"
	"library/apps/api/internal/services"
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
			"password":   "harbour-lamp-97",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
