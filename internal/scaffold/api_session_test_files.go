package scaffold

// Tests for server-side sessions.
//
// These exercise the service layer directly against an in-memory SQLite
// database, because the interesting behaviour is all in the state machine —
// rotation, replay detection, the two timeouts, and owner scoping — not in the
// HTTP plumbing around it.

func apiSessionTestGo() string {
	return `package services_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

func newSessionDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(tb, err)
	require.NoError(tb, db.AutoMigrate(&models.Session{}))
	return db
}

// sessionCtx fakes the request a session is created from, so UserAgent and IP
// land on the row the way they would in production.
func sessionCtx(userAgent string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest("POST", "/api/auth/login", nil)
	req.Header.Set("User-Agent", userAgent)
	c.Request = req
	return c
}

func TestCreateSessionStoresOnlyTheHash(t *testing.T) {
	db := newSessionDB(t)
	s, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "refresh-token-abc")
	require.NoError(t, err)

	assert.Equal(t, models.HashSessionToken("refresh-token-abc"), s.TokenHash)
	assert.NotContains(t, s.TokenHash, "refresh-token-abc", "the raw token must never be persisted")
	assert.Equal(t, "Chrome", s.UserAgent)
	assert.Nil(t, s.RevokedAt)
}

func TestRotateSessionSwapsTheToken(t *testing.T) {
	db := newSessionDB(t)
	_, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "old-token")
	require.NoError(t, err)

	rotated, err := services.RotateSession(db, sessionCtx("Chrome"), "old-token", "new-token")
	require.NoError(t, err)

	var stored models.Session
	require.NoError(t, db.First(&stored, "id = ?", rotated.ID).Error)
	assert.Equal(t, models.HashSessionToken("new-token"), stored.TokenHash)
	assert.Equal(t, models.HashSessionToken("old-token"), stored.PrevTokenHash)

	// The old token is now worthless.
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "old-token", "another-token")
	assert.ErrorIs(t, err, services.ErrSessionInvalid)
}

// Replaying a rotated token is the signature of a stolen refresh token: the
// thief and the victim cannot both hold the current one, so whoever is second
// presents a rotated one. That must kill the session, not just fail the call.
func TestRotateSessionReplayKillsTheSession(t *testing.T) {
	db := newSessionDB(t)
	created, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-2")
	require.NoError(t, err)

	// The attacker replays the token they captured before rotation.
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-3")
	assert.ErrorIs(t, err, services.ErrSessionInvalid)

	var stored models.Session
	require.NoError(t, db.First(&stored, "id = ?", created.ID).Error)
	require.NotNil(t, stored.RevokedAt, "a replayed token must revoke the session")

	// And the legitimate holder is locked out too — by design. Better both
	// re-authenticate than let a thief ride along unnoticed.
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-2", "token-4")
	assert.ErrorIs(t, err, services.ErrSessionInvalid)
}

func TestSessionIdleTimeout(t *testing.T) {
	db := newSessionDB(t)
	s, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)

	// Backdate last activity past the idle window.
	stale := time.Now().Add(-services.SessionIdleTimeout - time.Hour)
	require.NoError(t, db.Model(s).Update("last_seen_at", stale).Error)

	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-2")
	assert.ErrorIs(t, err, services.ErrSessionInvalid, "an idle session must not refresh")
}

func TestSessionAbsoluteTimeout(t *testing.T) {
	db := newSessionDB(t)
	s, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)

	// Still actively used, but older than the absolute cap.
	require.NoError(t, db.Model(s).Updates(map[string]interface{}{
		"expires_at":   time.Now().Add(-time.Minute),
		"last_seen_at": time.Now(),
	}).Error)

	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-2")
	assert.ErrorIs(t, err, services.ErrSessionInvalid, "an expired session must not refresh regardless of activity")
}

func TestRevokeSessionIsScopedToItsOwner(t *testing.T) {
	db := newSessionDB(t)
	mine, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-mine")
	require.NoError(t, err)

	// Another user must not be able to revoke it, even knowing the id.
	err = services.RevokeSession(db, "user-2", mine.ID)
	assert.ErrorIs(t, err, services.ErrSessionInvalid)

	var stored models.Session
	require.NoError(t, db.First(&stored, "id = ?", mine.ID).Error)
	assert.Nil(t, stored.RevokedAt, "someone else's revoke must not touch my session")

	// The owner can.
	require.NoError(t, services.RevokeSession(db, "user-1", mine.ID))
	require.NoError(t, db.First(&stored, "id = ?", mine.ID).Error)
	assert.NotNil(t, stored.RevokedAt)
}

func TestRevokeAllUserSessions(t *testing.T) {
	db := newSessionDB(t)
	for _, tok := range []string{"t1", "t2", "t3"} {
		_, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", tok)
		require.NoError(t, err)
	}
	// A different user's session must survive.
	_, err := services.CreateSession(db, sessionCtx("Chrome"), "user-2", "other")
	require.NoError(t, err)

	// Spare the caller's own token.
	require.NoError(t, services.RevokeAllUserSessions(db, "user-1", "t2"))

	live, err := services.ListUserSessions(db, "user-1")
	require.NoError(t, err)
	require.Len(t, live, 1)
	assert.Equal(t, models.HashSessionToken("t2"), live[0].TokenHash)

	others, err := services.ListUserSessions(db, "user-2")
	require.NoError(t, err)
	assert.Len(t, others, 1, "revoking one user's sessions must not touch another's")
}

func TestRevokeSessionByToken(t *testing.T) {
	db := newSessionDB(t)
	_, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)

	require.NoError(t, services.RevokeSessionByToken(db, "token-1"))

	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-2")
	assert.ErrorIs(t, err, services.ErrSessionInvalid, "a logged-out session must not refresh")
}

func TestListUserSessionsExcludesRevoked(t *testing.T) {
	db := newSessionDB(t)
	keep, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "keep")
	require.NoError(t, err)
	drop, err := services.CreateSession(db, sessionCtx("Safari"), "user-1", "drop")
	require.NoError(t, err)

	require.NoError(t, services.RevokeSession(db, "user-1", drop.ID))

	live, err := services.ListUserSessions(db, "user-1")
	require.NoError(t, err)
	require.Len(t, live, 1)
	assert.Equal(t, keep.ID, live[0].ID)
}
`
}
