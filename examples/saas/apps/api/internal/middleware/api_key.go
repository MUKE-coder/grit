package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/authz"
	"saas/apps/api/internal/cache"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
)

// APIKeyOrAuth accepts either an API key or the normal JWT.
//
// It runs before the JWT middleware and, on a valid key, populates exactly the
// same context values (user_id, user_role) before calling Next. Everything
// downstream — RequireRole, the activity logger, handlers reading user_id —
// therefore works unchanged, which is the point: machine callers should not
// need a parallel set of handlers.
//
// A request with no key at all falls through untouched, so the JWT middleware
// behind this one handles browsers as before.
func APIKeyOrAuth(db *gorm.DB, jwtAuth gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractAPIKey(c)
		if token == "" {
			jwtAuth(c)
			return
		}

		key, err := services.VerifyAPIKey(db.WithContext(c.Request.Context()), token)
		if err != nil {
			respond.Fail(c, respond.CodeInvalidAPIKey, "That API key is not valid, has expired, or was revoked")
			return
		}

		// The safety property of the whole design, in five lines.
		//
		// A publishable key lives in a browser bundle or a phone binary, where
		// anyone can read it. It authenticates nobody. Reaching a protected
		// route with one is refused on the kind alone, before permissions are
		// even consulted, so no combination of scopes can talk its way past
		// this. If a publishable key leaks, and it will, the worst it opens is
		// what you already marked public.
		if key.Publishable() {
			respond.Fail(c, respond.CodePublishableKeyNotAllowed, "This is a protected endpoint. A publishable key reaches "+"public endpoints only; use a secret key from a server, or sign in.")
			return
		}

		if !key.AllowsEndpoint(c.Request.Method, c.FullPath()) {
			respond.Fail(c, respond.CodeEndpointNotAllowed, "That key is not allowed to call this endpoint")
			return
		}
		if !key.AllowsOrigin(c.GetHeader("Origin")) {
			respond.Fail(c, respond.CodeOriginNotAllowed, "That key is not allowed from this origin")
			return
		}

		var user models.User
		if err := db.First(&user, "id = ?", key.UserID).Error; err != nil {
			respond.Fail(c, respond.CodeInvalidAPIKey, "The owner of that key no longer exists")
			return
		}
		if !user.Active {
			respond.Fail(c, respond.CodeAccountDisabled, "The owner of that key is disabled")
			return
		}

		services.TouchAPIKey(db, key.ID)

		// This must set exactly what middleware.Auth sets — see auth.go. Handlers
		// read c.Get("user") for the whole record, and RequireRole reads
		// user_grants for permission checks.
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)
		if grants, err := authz.GrantsFor(db, user.ID); err == nil {
			c.Set("user_grants", grants)
		}
		// Marked so handlers that must refuse machine callers — anything
		// changing credentials, say — can tell the difference.
		c.Set("auth_via", "api_key")
		if len(key.Scopes) > 0 {
			c.Set("api_key_scopes", []string(key.Scopes))
		}
		c.Next()
	}
}

// extractAPIKey reads a key from either header. X-API-Key is what most
// integrations reach for; Authorization: Bearer is what an OpenAPI client
// generates. Supporting both costs four lines.
// RequireAPIKey guards a public endpoint.
//
// Public does not mean unauthenticated. It means no user is required: any
// valid key gets in, publishable or secret, and no JWT is needed. What that
// buys is not secrecy, because a publishable key is readable by anyone with
// your app. It buys identification, a rate-limit bucket per key, per-endpoint
// and per-origin narrowing, and the ability to turn one client off without
// deploying anything.
//
// Deliberately no user is set on the context. A public handler that quietly
// depended on c.Get("user_id") would then behave differently depending on
// which credential turned up, which is the kind of difference nobody notices
// until it is a data leak.
func RequireAPIKey(db *gorm.DB, limiter *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractAPIKey(c)
		if token == "" {
			respond.Fail(c, respond.CodeAPIKeyRequired, "Send your publishable key as X-API-Key")
			return
		}

		key, err := services.VerifyAPIKey(db.WithContext(c.Request.Context()), token)
		if err != nil {
			respond.Fail(c, respond.CodeInvalidAPIKey, "That API key is not valid, has expired, or was revoked")
			return
		}
		if !key.AllowsEndpoint(c.Request.Method, c.FullPath()) {
			respond.Fail(c, respond.CodeEndpointNotAllowed, "That key is not allowed to call this endpoint")
			return
		}
		if !key.AllowsOrigin(c.GetHeader("Origin")) {
			respond.Fail(c, respond.CodeOriginNotAllowed, "That key is not allowed from this origin")
			return
		}

		if !allowKeyRequest(c, limiter, key) {
			return
		}

		services.TouchAPIKey(db, key.ID)
		c.Set("api_key_id", key.ID)
		c.Set("api_key_kind", key.Kind)
		c.Set("auth_via", "api_key")
		c.Next()
	}
}

// allowKeyRequest enforces a key's own rate limit. Reports whether the request
// may continue; writes the 429 itself when it may not.
//
// A fixed window in Redis: one INCR per request against a key that carries the
// current minute, with a two minute expiry so the bucket cleans itself up. A
// sliding window would be fairer at the boundary and needs a sorted set per
// key; a fixed window is one round trip and is the right trade for throttling
// a misbehaving client rather than metering billing.
//
// No Redis means no per-key limiting, and that is deliberate rather than a
// fallback to something in-process: an in-memory counter would be per-instance,
// so the effective limit would silently multiply by however many API
// containers happen to be running.
func allowKeyRequest(c *gin.Context, limiter *cache.Cache, key *models.APIKey) bool {
	if key.RateLimit <= 0 || limiter == nil {
		return true
	}

	window := time.Now().UTC().Format("200601021504")
	bucket := "apikey:rl:" + key.ID + ":" + window

	ctx := c.Request.Context()
	count, err := limiter.Client().Incr(ctx, bucket).Result()
	if err != nil {
		// Redis being unreachable must not close the API. Failing open on a
		// throttle is the right direction: the global IP limit still applies,
		// and refusing every request because a counter is down would turn a
		// cache outage into an outage.
		return true
	}
	if count == 1 {
		limiter.Client().Expire(ctx, bucket, 2*time.Minute)
	}

	remaining := key.RateLimit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	c.Header("X-RateLimit-Limit", strconv.Itoa(key.RateLimit))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

	if int(count) > key.RateLimit {
		c.Header("Retry-After", "60")
		respond.Fail(c, respond.CodeRateLimited, "This API key is limited to "+strconv.Itoa(key.RateLimit)+" requests per minute")
		return false
	}
	return true
}

func extractAPIKey(c *gin.Context) string {
	if v := strings.TrimSpace(c.GetHeader("X-API-Key")); v != "" {
		return v
	}
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	const bearer = "Bearer "
	if strings.HasPrefix(auth, bearer) {
		token := strings.TrimSpace(auth[len(bearer):])
		// Only claim it if it looks like one of ours — otherwise a JWT would
		// be sent to the key verifier and rejected before the JWT middleware
		// ever sees it.
		if strings.HasPrefix(token, services.KeyTokenPrefix+"_") {
			return token
		}
	}
	return ""
}
