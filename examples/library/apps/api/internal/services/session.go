package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"library/apps/api/internal/models"
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
//
// The gin-shaped front door to CreateSessionCtx, kept because handlers call it
// with the request in hand. Everything it needs from that request is the three
// fields in RequestMeta.
func CreateSession(db *gorm.DB, c *gin.Context, userID, refreshToken string) (*models.Session, error) {
	return CreateSessionCtx(ContextOf(c), db, MetaOf(c), userID, refreshToken)
}

// CreateSessionCtx records a newly issued refresh token as a logged-in device.
//
// Contact-app review M29: this took a *gin.Context, so signing somebody in from
// a job, a CLI command or a test meant building a fake HTTP request first.
func CreateSessionCtx(ctx context.Context, db *gorm.DB, meta RequestMeta, userID, refreshToken string) (*models.Session, error) {
	now := time.Now()
	s := &models.Session{
		ID:         sessionIDFromToken(refreshToken),
		UserID:     userID,
		TokenHash:  models.HashSessionToken(refreshToken),
		UserAgent:  truncateStr(meta.UserAgent, 512),
		IP:         meta.IP,
		LastSeenAt: now,
		ExpiresAt:  now.Add(SessionAbsoluteTimeout),
	}
	if err := db.WithContext(ctx).Create(s).Error; err != nil {
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
	return RotateSessionCtx(ContextOf(c), db, MetaOf(c), oldToken, newToken)
}

// RotateSessionCtx is RotateSession without the gin context. See M29 on
// CreateSessionCtx for why both exist.
func RotateSessionCtx(ctx context.Context, db *gorm.DB, meta RequestMeta, oldToken, newToken string) (*models.Session, error) {
	db = db.WithContext(ctx)
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
			// A replay is the one time revocation must not quietly fail: the
			// token was captured, and the session it belongs to stays live.
			if err := db.Model(&replayed).Update("revoked_at", &now).Error; err != nil {
				return nil, err
			}
			forgetSession(replayed.ID)
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
		"ip":              meta.IP,
		"user_agent":      truncateStr(meta.UserAgent, 512),
	}).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// RevokeSessionByToken kills the session a refresh token belongs to (logout).
func RevokeSessionByToken(db *gorm.DB, refreshToken string) error {
	now := time.Now()
	hash := models.HashSessionToken(refreshToken)
	var sessionIDs []string
	if err := db.Model(&models.Session{}).Where("token_hash = ?", hash).Pluck("id", &sessionIDs).Error; err != nil {
		return err
	}
	if err := db.Model(&models.Session{}).
		Where("token_hash = ? AND revoked_at IS NULL", hash).
		Update("revoked_at", &now).Error; err != nil {
		return err
	}
	for _, id := range sessionIDs {
		forgetSession(id)
	}
	return nil
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
	forgetSession(sessionID)
	// Sockets are bound to a user, not to a session row, so the revoked
	// device cannot be singled out. Every connection this user holds is
	// closed and the ones still entitled reconnect.
	sessionsRevoked(userID)
	return nil
}

// OnSessionsRevoked runs after a user's sessions are revoked, with that user's
// id. routes.Setup points it at the realtime hub's DisconnectUser.
//
// Revoking a session marks a database row. It does not reach an already-open
// WebSocket, whose token was checked once at the handshake and is never
// consulted again, so without this hook "sign out of all devices" leaves the
// signed-out device receiving live events until its token expires. Closing the
// socket is the half of revocation that the row cannot do on its own.
//
// Nil by default so the services package stays free of a realtime dependency.
var OnSessionsRevoked func(userID string)

func sessionsRevoked(userID string) {
	if OnSessionsRevoked != nil {
		OnSessionsRevoked(userID)
	}
}

// RevokeAllUserSessions kills every session for a user. Call it on password
// change, on MFA change, and from "log out everywhere". exceptToken, when
// non-empty, spares the caller's own session.
//
// Note that exceptToken spares a session row but cannot spare a socket: every
// connection the user holds is closed and the surviving session simply
// reconnects. Dropping one live connection is a great deal better than leaving
// a revoked one open.
func RevokeAllUserSessions(db *gorm.DB, userID, exceptToken string) error {
	now := time.Now()
	q := db.Model(&models.Session{}).Where("user_id = ? AND revoked_at IS NULL", userID)
	if exceptToken != "" {
		q = q.Where("token_hash <> ?", models.HashSessionToken(exceptToken))
	}
	if err := q.Update("revoked_at", &now).Error; err != nil {
		return err
	}
	forgetUserSessions(userID)
	sessionsRevoked(userID)
	return nil
}

// ListUserSessions returns a user's live sessions, newest activity first.
func ListUserSessions(db *gorm.DB, userID string) ([]models.Session, error) {
	var out []models.Session
	err := db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("last_seen_at desc").Find(&out).Error
	return out, err
}

// sessionIDFromToken is the session id a refresh token carries, so the row's id
// is the one its access tokens name. A token without one gets a fresh id.
func sessionIDFromToken(token string) string {
	claims := &Claims{}
	if _, _, err := jwt.NewParser().ParseUnverified(token, claims); err != nil {
		return ""
	}
	if len(claims.SessionID) > 36 {
		return ""
	}
	return claims.SessionID
}

// SessionIDForToken finds the session a refresh token belongs to, for a token
// minted before sessions were named in the token itself.
func SessionIDForToken(db *gorm.DB, refreshToken string) string {
	var s models.Session
	if err := db.Select("id").Where("token_hash = ?", models.HashSessionToken(refreshToken)).First(&s).Error; err != nil {
		return ""
	}
	return s.ID
}

// SessionCheckTTL is how long a session's state is trusted between database
// reads. A session revoked through this process is refused at once; one revoked
// through another replica, within this long.
var SessionCheckTTL = 30 * time.Second

type sessionState struct {
	userID  string
	live    bool
	checked time.Time
}

var sessionStates = struct {
	sync.Mutex
	m map[string]sessionState
}{m: map[string]sessionState{}}

// SessionLive reports whether the session an access token names may still be
// used: it exists, is not revoked, and has not passed either timeout.
func SessionLive(db *gorm.DB, sessionID string) bool {
	now := time.Now()
	sessionStates.Lock()
	st, ok := sessionStates.m[sessionID]
	sessionStates.Unlock()
	if ok && now.Sub(st.checked) < SessionCheckTTL {
		return st.live
	}

	var s models.Session
	if err := db.Where("id = ?", sessionID).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rememberSession(sessionID, "", false, now)
		}
		return false
	}
	live := s.Active(now, SessionIdleTimeout)
	rememberSession(sessionID, s.UserID, live, now)
	return live
}

func rememberSession(sessionID, userID string, live bool, now time.Time) {
	sessionStates.Lock()
	defer sessionStates.Unlock()
	if len(sessionStates.m) > 10000 {
		for id, st := range sessionStates.m {
			if now.Sub(st.checked) >= SessionCheckTTL {
				delete(sessionStates.m, id)
			}
		}
	}
	sessionStates.m[sessionID] = sessionState{userID: userID, live: live, checked: now}
}

func forgetSession(sessionID string) {
	sessionStates.Lock()
	delete(sessionStates.m, sessionID)
	sessionStates.Unlock()
}

func forgetUserSessions(userID string) {
	sessionStates.Lock()
	for id, st := range sessionStates.m {
		if st.userID == userID {
			delete(sessionStates.m, id)
		}
	}
	sessionStates.Unlock()
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
