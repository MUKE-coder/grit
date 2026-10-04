package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"library/apps/api/internal/ids"
	"library/apps/api/internal/models"
)

// AuthService handles JWT token operations.
type AuthService struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration

	// DB is where sessions live. When it is set, an access token is refused as
	// soon as its session is revoked, rather than when the token expires.
	DB *gorm.DB
}

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// Token types. An access token lives for minutes and is what every route
// accepts; a refresh token lives for days and is only exchanged at /auth/refresh.
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// ErrWrongTokenType is returned when a token is presented where the other kind
// belongs.
var ErrWrongTokenType = errors.New("wrong token type")

// Claims represents JWT claims.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	// TokenType is "access" or "refresh". The two used to be the same shape, so
	// a seven-day refresh token was accepted as a bearer token everywhere.
	TokenType string `json:"typ,omitempty"`
	// SessionID is the sessions row both tokens of a pair belong to, so an access
	// token can be refused the moment that session is revoked.
	SessionID string `json:"sid,omitempty"`
	jwt.RegisteredClaims
}

// GenerateTokenPair creates a new access + refresh token pair for a new session.
// Record it with CreateSession: an access token whose session does not exist is
// refused.
func (s *AuthService) GenerateTokenPair(userID string, email, role string) (*TokenPair, error) {
	return s.GenerateSessionTokenPair(userID, email, role, ids.New())
}

// GenerateSessionTokenPair creates a token pair for an existing session. Refresh
// uses it, so the rotated tokens still name the session row they belong to.
func (s *AuthService) GenerateSessionTokenPair(userID, email, role, sessionID string) (*TokenPair, error) {
	accessToken, expiresAt, err := s.generateToken(userID, email, role, TokenTypeAccess, sessionID, s.AccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, _, err := s.generateToken(userID, email, role, TokenTypeRefresh, sessionID, s.RefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// ValidateAccessToken accepts an access token whose session is still live. It is
// what the auth middleware and the WebSocket handshake call.
func (s *AuthService) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeAccess {
		return nil, ErrWrongTokenType
	}
	if s.DB != nil && (claims.SessionID == "" || !SessionLive(s.DB, claims.SessionID)) {
		return nil, ErrSessionInvalid
	}
	return claims, nil
}

// ValidateRefreshToken accepts a refresh token. One minted before token types has
// none, and is accepted: the session row still has to exist and rotate, which is
// the check that matters, and it is how existing sessions survive the upgrade.
func (s *AuthService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeRefresh && claims.TokenType != "" {
		return nil, ErrWrongTokenType
	}
	return claims, nil
}

// ValidateToken checks a token's signature and expiry, and nothing about what
// kind of token it is. Prefer ValidateAccessToken or ValidateRefreshToken.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())

	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GenerateResetToken creates a random hex token for password resets.
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generating reset token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) generateToken(userID, email, role, tokenType, sessionID string, expiry time.Duration) (string, int64, error) {
	expiresAt := time.Now().Add(expiry)

	// Every token gets a unique jti. Without it, two tokens minted for the same
	// user in the same second are byte-identical — same claims, same
	// second-resolution exp, same key — so two different devices would share one
	// refresh token and could not be told apart or revoked independently.
	jti, err := GenerateResetToken()
	if err != nil {
		return "", 0, fmt.Errorf("generating token id: %w", err)
	}

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// RefreshCookiePath scopes the refresh cookie to the auth routes, so it is not
// sent on every request. routes.Setup builds it from APIVersion. It has to match
// where refresh and logout are mounted: written as /api/auth while they lived
// under /api/v1/auth, the browser never sent it to either, so refreshing from a
// cookie always failed and logout never revoked the session.
var RefreshCookiePath = "/api/v1/auth"

// legacyRefreshCookiePath is where projects before v3.244.0 set the cookie.
const legacyRefreshCookiePath = "/api/auth"

// SetAuthCookies writes the token pair as HttpOnly cookies so the browser
// holds the credentials out of JavaScript's reach. The native mobile and
// desktop clients keep using the Authorization: Bearer header, which is
// why the JSON body still includes the tokens — both paths work.
//
// Cookie names: grit_access (sent on every request) and grit_refresh
// (scoped to RefreshCookiePath so it isn't sent everywhere). Both are HttpOnly,
// Secure when on HTTPS, and SameSite=Lax so CSRF surface is limited to
// top-level navigations. The CSRF middleware adds defence in depth.
//
// Reference: docs/backend/authentication §"Token Storage on the Frontend".
func (s *AuthService) SetAuthCookies(c *gin.Context, pair *TokenPair) {
	secure := isRequestHTTPS(c)
	accessSeconds := int(s.AccessExpiry / time.Second)
	refreshSeconds := int(s.RefreshExpiry / time.Second)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("grit_access", pair.AccessToken, accessSeconds, "/", "", secure, true)
	c.SetCookie("grit_refresh", pair.RefreshToken, refreshSeconds, RefreshCookiePath, "", secure, true)
	// grit_signed_in carries no secret. It lives as long as the session so the
	// web app's middleware can tell a signed-in browser from a stranger on the
	// same host: grit_access expires with its token and grit_refresh is scoped
	// to the auth routes, so neither can be seen on an admin page.
	c.SetCookie("grit_signed_in", "1", refreshSeconds, "/", "", secure, true)
	// A cookie at the old path is never sent anywhere useful; clear it.
	c.SetCookie("grit_refresh", "", -1, legacyRefreshCookiePath, "", secure, true)
}

// ClearAuthCookies expires both auth cookies. Call this from the Logout
// handler so a stolen browser session is cut off as soon as the user
// signs out.
func (s *AuthService) ClearAuthCookies(c *gin.Context) {
	secure := isRequestHTTPS(c)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("grit_access", "", -1, "/", "", secure, true)
	c.SetCookie("grit_signed_in", "", -1, "/", "", secure, true)
	c.SetCookie("grit_refresh", "", -1, RefreshCookiePath, "", secure, true)
	c.SetCookie("grit_refresh", "", -1, legacyRefreshCookiePath, "", secure, true)
}

// The reads and writes behind signing in, so the handler runs none of its own.
//
// They live here rather than in internal/handlers/auth.go because each is a
// decision about data rather than about HTTP: which columns of a two-factor row
// are safe to load, what an unknown address answers, when a failure count is
// cleared. Generated resources were moved onto their services in v3.224 and the
// framework's own handlers were not, which left the first file a developer
// opens breaking the rule the generator enforces.

// EnabledTwoFactor returns the account's second factor, or an error when it has
// none to ask for.
//
// Only the id and the method are read. The secret is encrypted when
// FIELD_ENCRYPTION_KEY is set, and loading it here would let a secret that
// cannot be decrypted read as "this account has no second factor", which turns
// a key problem into a silently disabled second factor.
func (s *AuthService) EnabledTwoFactor(ctx context.Context, userID string) (*models.TwoFactorConfig, error) {
	var cfg models.TwoFactorConfig
	if err := s.DB.WithContext(ctx).Select("id", "method").
		Where("user_id = ? AND enabled = ?", userID, true).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

// StartPendingTOTP records a half-completed sign-in: the password is right and
// the second factor is still outstanding.
func (s *AuthService) StartPendingTOTP(ctx context.Context, pending *models.TOTPPendingToken) error {
	return s.DB.WithContext(ctx).Create(pending).Error
}

// UserByEmail is the sign-in lookup.
//
// The caller spends a bcrypt comparison whether or not this finds anything, so
// an address with no account takes the same time as one with a wrong password.
// Returning the error rather than a nil user keeps that decision at the call
// site, where the timing is managed.
func (s *AuthService) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UserByID reloads the account behind a refresh token.
//
// A refresh token is otherwise valid for its whole lifetime even after the
// account is deleted or deactivated. Re-reading the row closes that window and
// lets a role change take effect on the next refresh, which is partial
// revocation without a server-side token store.
func (s *AuthService) UserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := s.DB.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ClearLoginFailures resets the lockout counters.
//
// Called when a sign-in finishes, never when the password alone is accepted:
// clearing it at the password with a second factor still ahead handed a fresh
// set of guesses at the code to anybody who knew the password and simply signed
// in again.
func (s *AuthService) ClearLoginFailures(ctx context.Context, userID string) error {
	return s.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error
}

// CountLoginFailure counts one wrong credential against an account and returns
// the count as the database now holds it.
//
// The increment is in SQL, so parallel attempts cannot each read the same number
// and overwrite one another, and the count is read back for the same reason: the
// one the caller arrived with is already stale, and it is what a lockout decision
// is made from. Both the password step and the second factor had their own copy
// of this, written months apart.
func (s *AuthService) CountLoginFailure(ctx context.Context, userID string) (int, error) {
	if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("failed_login_count", gorm.Expr("failed_login_count + 1")).Error; err != nil {
		return 0, err
	}
	var fresh models.User
	if err := s.DB.WithContext(ctx).Select("id", "failed_login_count").
		First(&fresh, "id = ?", userID).Error; err != nil {
		return 0, err
	}
	return fresh.FailedLoginCount, nil
}

// LockAccount locks an account until until, and clears the counter and the
// pending two-factor tokens that are now useless to it.
//
// The pending tokens go because each one is a set of guesses at a code: leaving
// them is the difference between locking an account and pausing it.
func (s *AuthService) LockAccount(ctx context.Context, userID string, until time.Time) error {
	if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{"locked_until": until, "failed_login_count": 0}).Error; err != nil {
		return err
	}
	return s.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.TOTPPendingToken{}).Error
}

// isRequestHTTPS returns true when the request is on HTTPS (directly or
// via a trusted proxy that set X-Forwarded-Proto=https). We use it to flip
// the Secure cookie flag so the browser refuses to send these cookies
// over an unencrypted hop.
func isRequestHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	if strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
		return true
	}
	return false
}
