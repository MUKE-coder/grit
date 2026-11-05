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

	// The old token is worthless once the grace window has passed. Inside it,
	// presenting the old token is a race rather than a replay and is allowed
	// on purpose: see TestTwoConcurrentRefreshesDoNotKillTheSession.
	stale := time.Now().Add(-services.SessionRotationGrace - time.Minute)
	require.NoError(t, db.Model(&stored).Update("rotated_at", &stale).Error)

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

	// Past the grace window, which is the replay rotation exists to catch: a
	// token captured and used later. One used in the same second is a thief
	// racing the victim, and no amount of rotation ever stopped that.
	stale := time.Now().Add(-services.SessionRotationGrace - time.Minute)
	require.NoError(t, db.Model(created).Update("rotated_at", &stale).Error)

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

// Two tabs whose access tokens expire together both present the token they
// hold, and the second one presents what rotation has already replaced. That
// is not a replay, it is a race, and before the grace window it signed the
// user out of every device (grit#153).
func TestTwoConcurrentRefreshesDoNotKillTheSession(t *testing.T) {
	db := newSessionDB(t)
	created, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)

	// The first tab refreshes.
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-2")
	require.NoError(t, err)

	// The second tab's request was already in flight with the same token.
	second, err := services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-3")
	require.NoError(t, err, "a refresh racing another one must not be read as a replay")
	require.Equal(t, created.ID, second.ID, "the race started a different session")

	var stored models.Session
	require.NoError(t, db.First(&stored, "id = ?", created.ID).Error)
	assert.Nil(t, stored.RevokedAt, "the session was revoked by an honest collision")

	// The token the second caller was handed is the live one.
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-3", "token-4")
	assert.NoError(t, err)
}

// The window is anchored to the first rotation. Presenting the old token over
// and over must not keep pushing the deadline out, or a captured token never
// expires as long as it is used.
func TestTheGraceWindowDoesNotSlide(t *testing.T) {
	db := newSessionDB(t)
	created, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-2")
	require.NoError(t, err)

	var afterFirst models.Session
	require.NoError(t, db.First(&afterFirst, "id = ?", created.ID).Error)
	require.NotNil(t, afterFirst.RotatedAt, "rotation did not record when it happened")

	// A second use of the old token, inside the window.
	_, err = services.RotateSession(db, sessionCtx("Chrome"), "token-1", "token-3")
	require.NoError(t, err)

	var afterSecond models.Session
	require.NoError(t, db.First(&afterSecond, "id = ?", created.ID).Error)
	require.NotNil(t, afterSecond.RotatedAt)
	assert.True(t, afterSecond.RotatedAt.Equal(*afterFirst.RotatedAt),
		"the grace window moved, so the old token outlives its deadline")
	assert.Equal(t, afterFirst.PrevTokenHash, afterSecond.PrevTokenHash,
		"the previous token changed, so the window now protects a different token")
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

// Two instances in one process keep their own revocation hook.
//
// The hook closes the user's live WebSockets, and the hub it has to reach
// belongs to one instance. While there was a single package-level hook, the
// instance that started second owned it: signing out everywhere on the first
// instance closed the second instance's sockets and left the first instance's
// open, and nothing anywhere reported it. The hook is registered against the
// database handle now, which is the one thing the revoke functions are given.
// The refresh cookie goes where this service says, not where the package
// variable says.
//
// Asserted against a path the fallback could not have produced, because the
// fallback's value and what routes.Setup computes are the same string for a v1
// API: a test that only checked the path would pass with the field ignored.
func TestRefreshCookieTakesItsPathFromTheService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest("POST", "/api/v2/auth/login", nil)

	auth := &services.AuthService{
		AccessExpiry:      time.Minute,
		RefreshExpiry:     time.Hour,
		RefreshCookiePath: "/api/v2/auth",
	}
	auth.SetAuthCookies(c, &services.TokenPair{AccessToken: "access", RefreshToken: "refresh"})

	var written *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		// SetAuthCookies also expires a cookie at the pre-v3.244.0 path, which
		// is the same name with an empty value.
		if cookie.Name == "grit_refresh" && cookie.Value != "" {
			written = cookie
		}
	}
	require.NotNil(t, written, "no grit_refresh cookie was written")
	assert.Equal(t, "/api/v2/auth", written.Path)
	assert.NotEqual(t, services.RefreshCookiePath, written.Path,
		"the package variable decided the path, so the field is not being read")
}

func TestRevocationHookBelongsToItsDatabaseHandle(t *testing.T) {
	first, second := newSessionDB(t), newSessionDB(t)

	var firstSaw, secondSaw []string
	services.OnSessionsRevokedFor(first, func(userID string) { firstSaw = append(firstSaw, userID) })
	services.OnSessionsRevokedFor(second, func(userID string) { secondSaw = append(secondSaw, userID) })
	t.Cleanup(func() {
		services.OnSessionsRevokedFor(first, nil)
		services.OnSessionsRevokedFor(second, nil)
	})

	_, err := services.CreateSession(first, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)
	require.NoError(t, services.RevokeAllUserSessions(first, "user-1", ""))

	assert.Equal(t, []string{"user-1"}, firstSaw)
	assert.Empty(t, secondSaw, "the other instance's hook ran for a revocation on this one")
}

// A Sessions value takes its hook directly, for code that holds neither a
// registration nor the package variable.
func TestSessionsCarriesItsOwnHook(t *testing.T) {
	db := newSessionDB(t)
	mine, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", "token-1")
	require.NoError(t, err)

	var called []string
	store := services.Sessions{DB: db, OnRevoked: func(userID string) { called = append(called, userID) }}
	require.NoError(t, store.Revoke("user-1", mine.ID))
	assert.Equal(t, []string{"user-1"}, called)
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

func tokenService(db *gorm.DB) *services.AuthService {
	return &services.AuthService{
		Secret:        "a-test-secret-that-is-longer-than-32-characters",
		AccessExpiry:  time.Minute,
		RefreshExpiry: time.Hour,
		DB:            db,
	}
}

// A refresh token is not an access token, and an access token dies with its
// session rather than when it expires.
func TestAccessTokensFollowTheirSession(t *testing.T) {
	db := newSessionDB(t)
	auth := tokenService(db)
	pair, err := auth.GenerateTokenPair("user-1", "a@example.com", "USER")
	require.NoError(t, err)
	s, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", pair.RefreshToken)
	require.NoError(t, err)

	claims, err := auth.ValidateAccessToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, s.ID, claims.SessionID, "the row's id is the one the tokens name")

	_, err = auth.ValidateAccessToken(pair.RefreshToken)
	assert.Error(t, err, "a refresh token was accepted as an access token")
	_, err = auth.ValidateRefreshToken(pair.AccessToken)
	assert.Error(t, err, "an access token was accepted as a refresh token")

	require.NoError(t, services.RevokeSessionByToken(db, pair.RefreshToken))
	_, err = auth.ValidateAccessToken(pair.AccessToken)
	assert.Error(t, err, "an access token outlived its revoked session")
}

// A token minted without CreateSession names a session that does not exist.
func TestAnAccessTokenWithNoSessionIsRefused(t *testing.T) {
	db := newSessionDB(t)
	auth := tokenService(db)
	pair, err := auth.GenerateTokenPair("user-1", "a@example.com", "USER")
	require.NoError(t, err)
	_, err = auth.ValidateAccessToken(pair.AccessToken)
	assert.Error(t, err)
}

// Rotation keeps the session id, so refreshed access tokens still name the row.
func TestRefreshKeepsTheSessionID(t *testing.T) {
	db := newSessionDB(t)
	auth := tokenService(db)
	pair, err := auth.GenerateTokenPair("user-1", "a@example.com", "USER")
	require.NoError(t, err)
	s, err := services.CreateSession(db, sessionCtx("Chrome"), "user-1", pair.RefreshToken)
	require.NoError(t, err)

	next, err := auth.GenerateSessionTokenPair("user-1", "a@example.com", "USER", services.SessionIDForToken(db, pair.RefreshToken))
	require.NoError(t, err)
	_, err = services.RotateSession(db, sessionCtx("Chrome"), pair.RefreshToken, next.RefreshToken)
	require.NoError(t, err)

	claims, err := auth.ValidateAccessToken(next.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, s.ID, claims.SessionID)
}
`
}
