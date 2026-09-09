package plugin

import "strings"

// pairingModelGo is the in-flight handshake. One row per QR shown.
func pairingModelGo(ctx Context) string {
	src := "package models\n\n" + `import (
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
)

// PairingRequest is one "link a new device" handshake in flight.
//
// The browser asks for a code, renders it as a QR, and polls. A signed-in
// device scans the code, is shown what it is about to trust, and approves. The
// browser's next poll returns a real token pair.
//
// The row is short-lived by design and pruned once claimed or expired: it is a
// handshake, not a record of anything worth keeping. What survives is the
// session the approval creates.
type PairingRequest struct {
	ID string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `

	// Code is the secret in the QR. Never returned by any endpoint except the
	// one that mints it, and never logged.
	Code string ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `

	// Filled in when a signed-in device approves.
	UserID     string     ` + "`" + `gorm:"size:36;index" json:"user_id,omitempty"` + "`" + `
	ApprovedAt *time.Time ` + "`" + `json:"approved_at,omitempty"` + "`" + `
	DeniedAt   *time.Time ` + "`" + `json:"denied_at,omitempty"` + "`" + `

	// ClaimedAt marks the single poll allowed to collect the tokens. Without
	// it a leaked code stays spendable for the rest of its TTL.
	ClaimedAt *time.Time ` + "`" + `json:"claimed_at,omitempty"` + "`" + `

	// What the browser said about itself. Shown to the approver before they
	// approve, which is the whole point: approving something you cannot see is
	// not consent.
	UserAgent string ` + "`" + `gorm:"size:512" json:"user_agent"` + "`" + `
	IP        string ` + "`" + `gorm:"size:64;index" json:"ip"` + "`" + `

	ExpiresAt time.Time ` + "`" + `gorm:"index;not null" json:"expires_at"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
}

func (m *PairingRequest) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// Pending reports whether this request can still be approved or denied.
func (m *PairingRequest) Pending(now time.Time) bool {
	return m.ApprovedAt == nil && m.DeniedAt == nil &&
		m.ClaimedAt == nil && now.Before(m.ExpiresAt)
}

// Live reports whether the browser should keep polling.
func (m *PairingRequest) Live(now time.Time) bool {
	return m.ClaimedAt == nil && m.DeniedAt == nil && now.Before(m.ExpiresAt)
}
`
	return strings.ReplaceAll(src, "{{MODULE}}", ctx.Module)
}

// pairingHandlerGo is the handshake itself.
func pairingHandlerGo(ctx Context) string {
	src := "package handlers\n\n" + `import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

const (
	// pairingTTL is deliberately short. A QR on a screen in a public place is
	// visible to everyone in the room, so the window in which a passer-by can
	// use it has to be small enough that they would have to be waiting for it.
	pairingTTL = 2 * time.Minute

	// pairingPerIP caps how many codes one address can have in flight. The
	// start endpoint is anonymous, so without a cap it mints rows for free.
	pairingPerIP = 5
)

// PairingHandler implements "scan a QR to sign this browser in".
type PairingHandler struct {
	DB   *gorm.DB
	Auth *services.AuthService
}

func NewPairingHandler(db *gorm.DB, auth *services.AuthService) *PairingHandler {
	return &PairingHandler{DB: db, Auth: auth}
}

// pairingError keeps every failure shaped the same. An endpoint that answers
// "unknown code" differently from "expired code" tells an attacker which
// guesses were close.
func pairingError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

// Start mints a code for an anonymous browser and returns it with a QR PNG.
//
//	POST /api/v1/pair/start
func (h *PairingHandler) Start(c *gin.Context) {
	now := time.Now()
	ip := c.ClientIP()

	// Sweep this address's dead rows first, so a caller that keeps abandoning
	// codes is not locked out by its own litter.
	h.DB.Where("ip = ? AND (expires_at < ? OR claimed_at IS NOT NULL OR denied_at IS NOT NULL)", ip, now).
		Delete(&models.PairingRequest{})

	var inFlight int64
	h.DB.Model(&models.PairingRequest{}).
		Where("ip = ? AND expires_at >= ?", ip, now).Count(&inFlight)
	if inFlight >= pairingPerIP {
		pairingError(c, http.StatusTooManyRequests, "TOO_MANY_REQUESTS",
			"Too many pairing attempts. Wait a moment and try again.")
		return
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		pairingError(c, http.StatusInternalServerError, "INTERNAL", "Could not generate a pairing code")
		return
	}
	code := hex.EncodeToString(raw)

	pr := models.PairingRequest{
		Code:      code,
		UserAgent: truncate(c.Request.UserAgent(), 512),
		IP:        ip,
		ExpiresAt: now.Add(pairingTTL),
	}
	if err := h.DB.Create(&pr).Error; err != nil {
		pairingError(c, http.StatusInternalServerError, "INTERNAL", "Could not start pairing")
		return
	}

	// The payload a scanner reads. A URL rather than a bare code so a phone
	// camera that is not the app still does something sensible with it.
	png, err := qrcode.Encode(pairingURL(c, code), qrcode.Medium, 320)
	if err != nil {
		pairingError(c, http.StatusInternalServerError, "INTERNAL", "Could not render the QR code")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"code":       code,
		"qr_png":     "data:image/png;base64," + base64.StdEncoding.EncodeToString(png),
		"expires_at": pr.ExpiresAt,
		"expires_in": int(pairingTTL.Seconds()),
	}})
}

// Status is polled by the browser.
//
//	GET /api/v1/pair/:code
func (h *PairingHandler) Status(c *gin.Context) {
	var pr models.PairingRequest
	if err := h.DB.Where("code = ?", c.Param("code")).First(&pr).Error; err != nil {
		pairingError(c, http.StatusNotFound, "NOT_FOUND", "This pairing code is no longer valid")
		return
	}
	now := time.Now()

	if pr.DeniedAt != nil {
		pairingError(c, http.StatusForbidden, "DENIED", "The request was refused on the other device")
		return
	}
	if !pr.Live(now) {
		pairingError(c, http.StatusGone, "EXPIRED", "This pairing code is no longer valid")
		return
	}
	if pr.ApprovedAt == nil {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"status": "pending"}})
		return
	}

	// Approved. Claim it under a conditional update so two polls arriving
	// together cannot both walk away with a token pair.
	res := h.DB.Model(&models.PairingRequest{}).
		Where("id = ? AND claimed_at IS NULL", pr.ID).
		Update("claimed_at", now)
	if res.Error != nil || res.RowsAffected == 0 {
		pairingError(c, http.StatusGone, "EXPIRED", "This pairing code is no longer valid")
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", pr.UserID).Error; err != nil {
		pairingError(c, http.StatusInternalServerError, "INTERNAL", "The approving account no longer exists")
		return
	}

	pair, err := h.Auth.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		pairingError(c, http.StatusInternalServerError, "INTERNAL", "Could not issue tokens")
		return
	}

	// Back the refresh token with a real session row, so the paired browser
	// appears in Active Sessions and is revoked by the machinery already there.
	if _, err := services.CreateSession(h.DB, c, user.ID, pair.RefreshToken); err != nil {
		pairingError(c, http.StatusInternalServerError, "INTERNAL", "Could not open a session")
		return
	}
	// Cookies as well as the body: the web app authenticates by HttpOnly
	// cookie, and a token in a JSON body it cannot store is no use to it.
	h.Auth.SetAuthCookies(c, pair)

	// The handshake is over. Nothing here is worth keeping.
	h.DB.Delete(&models.PairingRequest{}, "id = ?", pr.ID)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"status": "approved",
		"tokens": pair,
		"user":   user,
	}})
}

// Describe tells the approving device what it is about to trust.
//
// This exists so the phone can show "Chrome on Windows, 41.210.x.x" before it
// asks. Without it the user is approving an opaque code, which is not consent
// and turns a shoulder-surfed QR into a silent takeover.
//
//	GET /api/v1/pair/:code/request
func (h *PairingHandler) Describe(c *gin.Context) {
	var pr models.PairingRequest
	if err := h.DB.Where("code = ?", c.Param("code")).First(&pr).Error; err != nil {
		pairingError(c, http.StatusNotFound, "NOT_FOUND", "This pairing code is no longer valid")
		return
	}
	if !pr.Pending(time.Now()) {
		pairingError(c, http.StatusGone, "EXPIRED", "This pairing code is no longer valid")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"user_agent": pr.UserAgent,
		"ip":         pr.IP,
		"created_at": pr.CreatedAt,
		"expires_at": pr.ExpiresAt,
	}})
}

// Approve binds the request to the caller's account.
//
//	POST /api/v1/pair/:code/approve
func (h *PairingHandler) Approve(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		pairingError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Sign in first")
		return
	}

	now := time.Now()
	// Conditional update, so two devices racing on one code cannot both bind
	// it, and an expired one cannot be revived.
	res := h.DB.Model(&models.PairingRequest{}).
		Where("code = ? AND approved_at IS NULL AND denied_at IS NULL AND claimed_at IS NULL AND expires_at >= ?",
			c.Param("code"), now).
		Updates(map[string]interface{}{"user_id": userID, "approved_at": now})
	if res.Error != nil {
		pairingError(c, http.StatusInternalServerError, "INTERNAL", "Could not approve the request")
		return
	}
	if res.RowsAffected == 0 {
		pairingError(c, http.StatusGone, "EXPIRED", "This pairing code is no longer valid")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"status": "approved"},
		"message": "Device linked",
	})
}

// Deny refuses the request and burns the code immediately.
//
// A first-class endpoint rather than letting it lapse: "that wasn't me" has to
// be one tap, and leaving the code live for the rest of its TTL after the user
// has said no is the wrong answer to the one question that matters here.
//
//	POST /api/v1/pair/:code/deny
func (h *PairingHandler) Deny(c *gin.Context) {
	if c.GetString("user_id") == "" {
		pairingError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Sign in first")
		return
	}
	now := time.Now()
	res := h.DB.Model(&models.PairingRequest{}).
		Where("code = ? AND approved_at IS NULL AND denied_at IS NULL", c.Param("code")).
		Update("denied_at", now)
	if res.Error != nil || res.RowsAffected == 0 {
		pairingError(c, http.StatusGone, "EXPIRED", "This pairing code is no longer valid")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"status": "denied"},
		"message": "Request refused",
	})
}

// pairingURL is what the QR encodes. Absolute, so a phone camera outside the
// app opens something meaningful instead of showing a hex string.
func pairingURL(c *gin.Context, code string) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host + "/api/v1/pair/" + code
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
`
	return strings.ReplaceAll(src, "{{MODULE}}", ctx.Module)
}

// pairingHandlerTestGo covers the three races that make this worth shipping as
// a plugin rather than leaving to every project.
func pairingHandlerTestGo(ctx Context) string {
	src := "package handlers\n\n" + `import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

func pairingTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// One connection: every new connection to :memory: is its own empty
	// database, so a pool turns this into a flaky test.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.PairingRequest{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { _ = db.Migrator().DropTable(&models.PairingRequest{}) })
	return db
}

func seedPairing(t *testing.T, db *gorm.DB, mutate func(*models.PairingRequest)) models.PairingRequest {
	t.Helper()
	pr := models.PairingRequest{
		Code:      "test-code",
		UserAgent: "Chrome on Windows",
		IP:        "203.0.113.7",
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}
	if mutate != nil {
		mutate(&pr)
	}
	if err := db.Create(&pr).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	return pr
}

func pairingCtx(w *httptest.ResponseRecorder, code, userID string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Params = gin.Params{{Key: "code", Value: code}}
	if userID != "" {
		c.Set("user_id", userID)
	}
	return c
}

// Two devices scanning the same QR must not both bind it.
//
// A read-then-write would let both see "pending" and both write their own
// user id, and the browser would be signed in as whichever landed second.
func TestPairingApproveIsSingleWinner(t *testing.T) {
	db := pairingTestDB(t)
	seedPairing(t, db, nil)
	h := NewPairingHandler(db, nil)

	const racers = 8
	var wg sync.WaitGroup
	results := make([]int, racers)
	start := make(chan struct{})

	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w := httptest.NewRecorder()
			<-start
			h.Approve(pairingCtx(w, "test-code", "user-"+string(rune('a'+i))))
			results[i] = w.Code
		}(i)
	}
	close(start)
	wg.Wait()

	won := 0
	for _, code := range results {
		if code == http.StatusOK {
			won++
		}
	}
	if won != 1 {
		t.Errorf("%d of %d concurrent approvals succeeded; exactly one may win", won, racers)
	}
}

// An expired request cannot be approved, however the clock is read.
func TestPairingExpiredCannotBeApproved(t *testing.T) {
	db := pairingTestDB(t)
	seedPairing(t, db, func(pr *models.PairingRequest) {
		pr.ExpiresAt = time.Now().Add(-time.Second)
	})
	h := NewPairingHandler(db, nil)

	w := httptest.NewRecorder()
	h.Approve(pairingCtx(w, "test-code", "user-1"))
	if w.Code != http.StatusGone {
		t.Errorf("approving an expired request returned %d, want 410", w.Code)
	}
}

// Denying burns the code straight away rather than letting it lapse.
func TestPairingDenyClosesTheWindow(t *testing.T) {
	db := pairingTestDB(t)
	seedPairing(t, db, nil)
	h := NewPairingHandler(db, nil)

	w := httptest.NewRecorder()
	h.Deny(pairingCtx(w, "test-code", "user-1"))
	if w.Code != http.StatusOK {
		t.Fatalf("deny returned %d, want 200", w.Code)
	}

	w = httptest.NewRecorder()
	h.Approve(pairingCtx(w, "test-code", "user-1"))
	if w.Code == http.StatusOK {
		t.Error("a denied request was approved afterwards: 'that wasn't me' has " +
			"to close the window, not just pause it")
	}
}

// Describe must not hand out anything that helps someone use the code.
func TestPairingDescribeLeaksNoSecret(t *testing.T) {
	db := pairingTestDB(t)
	seedPairing(t, db, nil)
	h := NewPairingHandler(db, nil)

	w := httptest.NewRecorder()
	c := pairingCtx(w, "test-code", "user-1")
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.Describe(c)

	if w.Code != http.StatusOK {
		t.Fatalf("describe returned %d, want 200", w.Code)
	}
	if body := w.Body.String(); contains(body, "test-code") {
		t.Error("the response echoes the pairing code back")
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
`
	return strings.ReplaceAll(src, "{{MODULE}}", ctx.Module)
}
