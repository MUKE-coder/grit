package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"library/apps/api/internal/config"
	"library/apps/api/internal/mail/mailtest"
	"library/apps/api/internal/models"
	"library/apps/api/internal/services"
)

// newTestDB opens a fresh in-memory SQLite database and migrates what the auth
// handlers write to. It accepts testing.TB so both tests and benchmarks can
// call it.
//
// Registering does more than insert a user: it opens a session and writes an
// activity row. Those writes are best-effort and only logged, so a missing
// table does not fail the request outright — it surfaces later as a wrong
// status code, which is a confusing way to learn the test database was
// incomplete.
//
// email_verification_tokens is deliberately NOT migrated. That write happens on
// a goroutine (deliverVerificationEmail, so a slow SMTP call cannot hold the
// response open), and this connection is a single in-memory SQLite database.
// Creating the table turns a harmless early return into a background writer
// racing whatever the test does next. If you are testing verification itself,
// migrate it in that test and call the service directly rather than through the
// register handler.
func newTestDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(tb, err)

	// One connection, and only one.
	//
	// Every connection to ":memory:" gets its OWN empty database. The pool
	// opens a second one as soon as two queries overlap, which they do here
	// because registering writes its activity row on a goroutine: the login
	// that follows lands on a fresh, empty database, finds no user, and
	// returns 401. It passed or failed depending on timing, about six runs in
	// ten, which reads as a broken login rather than a test-harness problem.
	sqlDB, err := db.DB()
	require.NoError(tb, err)
	sqlDB.SetMaxOpenConns(1)

	require.NoError(tb, db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.UserActivity{},
	))
	return db
}

// db is passed explicitly rather than left nil, because the service reads the
// account on every sign-in. A test that passes nil is saying it never reaches
// the database, which the realtime tests do and these do not.
func newTestAuthSvc(cfg *config.Config, db *gorm.DB) *services.AuthService {
	return &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
		DB:            db,
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

// seedMailUser inserts a user directly. Registering would also start the
// verification email on a goroutine, which these tests do not want racing them.
func seedMailUser(t *testing.T, db *gorm.DB) models.User {
	t.Helper()
	user := models.User{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", Password: "not-a-real-hash", Role: "USER", Active: true}
	require.NoError(t, db.Create(&user).Error)
	return user
}

// The reset email goes through the Mailer, so the test reads it from the
// mailtest fake instead of a provider or a log line.
func TestAuthHandler_ForgotPassword_EmailsTheResetLink(t *testing.T) {
	cfg := testCfg()
	cfg.OAuthFrontendURL = "http://localhost:3001"
	db := newTestDB(t)
	require.NoError(t, db.AutoMigrate(&models.PasswordResetToken{}))
	seedMailUser(t, db)
	fake := mailtest.New()
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg, Mailer: fake.Mailer()}
	r := newAuthRouter(h)
	r.POST("/api/auth/forgot-password", h.ForgotPassword)

	w := postJSON(t, r, "/api/auth/forgot-password", map[string]string{"email": "jane@example.com"})
	assert.Equal(t, http.StatusOK, w.Code)
	msg := fake.AssertSentWithin(t, 5*time.Second, "jane@example.com", "Reset your password")
	assert.Contains(t, msg.HTML, "http://localhost:3001/reset-password?token=")

	// An unknown address gets the same answer and no email.
	fake.Reset()
	w = postJSON(t, r, "/api/auth/forgot-password", map[string]string{"email": "nobody@example.com"})
	assert.Equal(t, http.StatusOK, w.Code)
	time.Sleep(100 * time.Millisecond)
	fake.AssertNothingSent(t)
}

func TestAuthHandler_SendVerificationEmail_EmailsTheLink(t *testing.T) {
	cfg := testCfg()
	cfg.OAuthFrontendURL = "http://localhost:3001"
	db := newTestDB(t)
	require.NoError(t, db.AutoMigrate(&models.EmailVerificationToken{}))
	user := seedMailUser(t, db)
	fake := mailtest.New()
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg, Mailer: fake.Mailer()}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/auth/send-verification", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		h.SendVerificationEmail(c)
	})

	w := postJSON(t, r, "/api/auth/send-verification", map[string]string{})
	assert.Equal(t, http.StatusOK, w.Code)
	msg := fake.AssertSentWithin(t, 5*time.Second, "jane@example.com", "Confirm your email")
	assert.Contains(t, msg.HTML, "http://localhost:3001/verify-email?token=")
}

func TestAuthHandler_Register_Success(t *testing.T) {
	cfg := testCfg()
	db := newTestDB(t)
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg}
	r := newAuthRouter(h)

	w := postJSON(t, r, "/api/auth/register", map[string]string{
		"first_name": "Jane",
		"last_name":  "Doe",
		"email":      "jane@example.com",
		"password":   "harbour-lamp-97",
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
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg}
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
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg}
	r := newAuthRouter(h)

	body := map[string]string{
		"first_name": "Alice",
		"last_name":  "Smith",
		"email":      "alice@example.com",
		"password":   "harbour-lamp-97",
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
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg}
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
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg}
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
	h := &AuthHandler{DB: db, AuthService: newTestAuthSvc(cfg, db), Config: cfg}
	r := newAuthRouter(h)

	w := postJSON(t, r, "/api/auth/login", map[string]string{
		"email":    "nobody@example.com",
		"password": "anything",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
