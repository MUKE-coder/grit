package scaffold

// Server-side sessions.
//
// A JWT is self-contained: once issued it is valid until it expires, and
// nothing the server does can take it back. That is fine for a short-lived
// access token and unacceptable for the refresh token behind it — it means
// "log out all devices", "revoke this laptop", and "kill every session when
// the password changes" are all impossible, which is the first thing an
// enterprise security review asks for.
//
// So every refresh token is now backed by a row. The token itself is never
// stored — only its SHA-256 — so a database leak doesn't hand out sessions.
// Refresh looks the row up, enforces both timeouts, and ROTATES the token.
//
// Two timeouts, because auditors ask for both and most apps ship one:
//   - idle:     no refresh for SESSION_IDLE_TIMEOUT  -> dead
//   - absolute: older than SESSION_ABSOLUTE_TIMEOUT  -> dead regardless of use
//
// Revocation applies at refresh time. An access token already in flight stays
// valid until it expires (minutes), which is the standard trade-off for not
// hitting the database on every request.

func apiSessionModelGo() string {
	return `package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)

// Session is a server-side record of one issued refresh token — i.e. one
// logged-in device. It exists so a refresh token can be revoked; a bare JWT
// cannot be.
type Session struct {
	ID     string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `

	// TokenHash is the SHA-256 of the current refresh token. The token itself is
	// never persisted, so a dump of this table cannot be replayed as a session.
	TokenHash string ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `

	// PrevTokenHash is the hash of the token this one replaced. If a client ever
	// presents THAT token again, the token was captured and replayed after
	// rotation — the session is killed rather than refreshed.
	PrevTokenHash string ` + "`" + `gorm:"size:64;index" json:"-"` + "`" + `

	UserAgent string ` + "`" + `gorm:"size:512" json:"user_agent"` + "`" + `
	IP        string ` + "`" + `gorm:"size:64" json:"ip"` + "`" + `

	CreatedAt  time.Time  ` + "`" + `json:"created_at"` + "`" + `
	LastSeenAt time.Time  ` + "`" + `gorm:"index" json:"last_seen_at"` + "`" + `
	ExpiresAt  time.Time  ` + "`" + `gorm:"index" json:"expires_at"` + "`" + `
	RevokedAt  *time.Time ` + "`" + `gorm:"index" json:"revoked_at,omitempty"` + "`" + `
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}

// HashSessionToken returns the storage form of a refresh token. Exported so the
// service and any test can derive the same value.
func HashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Active reports whether the session may still be used to refresh.
func (s *Session) Active(now time.Time, idleTimeout time.Duration) bool {
	if s.RevokedAt != nil {
		return false
	}
	if now.After(s.ExpiresAt) {
		return false // absolute timeout
	}
	if idleTimeout > 0 && now.Sub(s.LastSeenAt) > idleTimeout {
		return false // idle timeout
	}
	return true
}
`
}

func apiSessionServiceGo() string {
	return `package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// Session lifetimes. Overridable per deployment; the defaults are deliberately
// conservative — an unused login dies in a week, and any login dies in 30 days
// no matter how active it is.
var (
	SessionIdleTimeout     = 7 * 24 * time.Hour
	SessionAbsoluteTimeout = 30 * 24 * time.Hour
)

// ErrSessionInvalid is returned when a refresh token has no live session:
// unknown, revoked, idle-expired, absolutely expired, or replayed.
var ErrSessionInvalid = errors.New("session is not valid")

// CreateSession records a newly issued refresh token as a logged-in device.
func CreateSession(db *gorm.DB, c *gin.Context, userID, refreshToken string) (*models.Session, error) {
	now := time.Now()
	s := &models.Session{
		UserID:     userID,
		TokenHash:  models.HashSessionToken(refreshToken),
		UserAgent:  truncateStr(c.Request.UserAgent(), 512),
		IP:         c.ClientIP(),
		LastSeenAt: now,
		ExpiresAt:  now.Add(SessionAbsoluteTimeout),
	}
	if err := db.Create(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

// RotateSession validates a presented refresh token and swaps it for a new one.
//
// Rotation is what makes a stolen refresh token survivable: the thief and the
// victim cannot both use it, and whoever refreshes second presents an already-
// rotated token — which lands in the PrevTokenHash branch and kills the session
// for both, surfacing the theft instead of silently sharing the account.
func RotateSession(db *gorm.DB, c *gin.Context, oldToken, newToken string) (*models.Session, error) {
	oldHash := models.HashSessionToken(oldToken)

	var s models.Session
	err := db.Where("token_hash = ?", oldHash).First(&s).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// Not the current token. If it's the PREVIOUS one, this is a replay of a
		// rotated token — treat it as compromise and kill the session.
		var replayed models.Session
		if db.Where("prev_token_hash = ?", oldHash).First(&replayed).Error == nil {
			now := time.Now()
			db.Model(&replayed).Update("revoked_at", &now)
		}
		return nil, ErrSessionInvalid
	}

	if !s.Active(time.Now(), SessionIdleTimeout) {
		return nil, ErrSessionInvalid
	}

	now := time.Now()
	if err := db.Model(&s).Updates(map[string]interface{}{
		"token_hash":      models.HashSessionToken(newToken),
		"prev_token_hash": oldHash,
		"last_seen_at":    now,
		"ip":              c.ClientIP(),
		"user_agent":      truncateStr(c.Request.UserAgent(), 512),
	}).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// RevokeSessionByToken kills the session a refresh token belongs to (logout).
func RevokeSessionByToken(db *gorm.DB, refreshToken string) error {
	now := time.Now()
	return db.Model(&models.Session{}).
		Where("token_hash = ? AND revoked_at IS NULL", models.HashSessionToken(refreshToken)).
		Update("revoked_at", &now).Error
}

// RevokeSession kills one session by id, scoped to its owner so a user can
// never revoke someone else's device.
func RevokeSession(db *gorm.DB, userID, sessionID string) error {
	now := time.Now()
	res := db.Model(&models.Session{}).
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", sessionID, userID).
		Update("revoked_at", &now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrSessionInvalid
	}
	return nil
}

// RevokeAllUserSessions kills every session for a user. Call it on password
// change, on MFA change, and from "log out everywhere". exceptToken, when
// non-empty, spares the caller's own session.
func RevokeAllUserSessions(db *gorm.DB, userID, exceptToken string) error {
	now := time.Now()
	q := db.Model(&models.Session{}).Where("user_id = ? AND revoked_at IS NULL", userID)
	if exceptToken != "" {
		q = q.Where("token_hash <> ?", models.HashSessionToken(exceptToken))
	}
	return q.Update("revoked_at", &now).Error
}

// ListUserSessions returns a user's live sessions, newest activity first.
func ListUserSessions(db *gorm.DB, userID string) ([]models.Session, error) {
	var out []models.Session
	err := db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("last_seen_at desc").Find(&out).Error
	return out, err
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
`
}

func apiSessionHandlerGo() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// SessionHandler exposes a user's own logged-in devices so they can see and
// revoke them — the "where am I signed in?" screen every security review wants.
type SessionHandler struct {
	DB *gorm.DB
}

func NewSessionHandler(db *gorm.DB) *SessionHandler {
	return &SessionHandler{DB: db}
}

type sessionView struct {
	models.Session
	// Current marks the session making this request, so the UI can label it and
	// avoid offering a "revoke" that would log you out mid-click.
	Current bool ` + "`" + `json:"current"` + "`" + `
}

// List returns the caller's active sessions.
func (h *SessionHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	sessions, err := services.ListUserSessions(h.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to load sessions"},
		})
		return
	}

	currentHash := ""
	if rt, err := c.Cookie("grit_refresh"); err == nil && rt != "" {
		currentHash = models.HashSessionToken(rt)
	}

	out := make([]sessionView, 0, len(sessions))
	for _, s := range sessions {
		out = append(out, sessionView{Session: s, Current: currentHash != "" && s.TokenHash == currentHash})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// Revoke kills one of the caller's sessions.
func (h *SessionHandler) Revoke(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := services.RevokeSession(h.DB, userID, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Session not found"},
		})
		return
	}
	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:   "session.revoke",
		Severity: "warn",
		Summary:  "Revoked a session",
	})
	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

// RevokeAll signs the caller out everywhere else, keeping the current session.
func (h *SessionHandler) RevokeAll(c *gin.Context) {
	userID := c.GetString("user_id")
	current, _ := c.Cookie("grit_refresh")
	if err := services.RevokeAllUserSessions(h.DB, userID, current); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to revoke sessions"},
		})
		return
	}
	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:   "session.revoke_all",
		Severity: "warn",
		Summary:  "Signed out of all other devices",
	})
	c.JSON(http.StatusOK, gin.H{"message": "All other sessions revoked"})
}
`
}
