package scaffold

// middlewareAuthLimitGo emits internal/middleware/authlimit.go: the rate limit
// on the routes that take a password, a code or an email address, which holds
// whether or not Sentinel is mounted.
//
// Enterprise review 2.8. Every limit on those routes lived inside
// mountSentinel, behind `if cfg.SentinelEnabled`. With SENTINEL_ENABLED=false,
// which is the default in development and a one-line change in production,
// /auth/login, /auth/forgot-password, /auth/totp/verify and the passkey routes
// had no limiter at all. The per-account lockout still barred an account after
// LOGIN_MAX_ATTEMPTS failures, so a single account was not brute-forceable, but
// nothing bounded credential stuffing across many accounts from one address, or
// a forgot-password loop pointed at one mailbox.
//
// This is the floor, not the ceiling. It is in-process and per-instance: with
// four replicas the effective limit is four times what the table says, and it
// resets on a restart. Sentinel keeps its counters in shared storage and is
// what you want in front of a fleet. The point of this file is that turning
// Sentinel off downgrades the protection instead of removing it.
func middlewareAuthLimitGo() string {
	return `package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthLimit is how many requests one client may make to one route in a window.
type AuthLimit struct {
	Requests int
	Window   time.Duration
}

// authCounter is one client's count for one route in the current window.
type authCounter struct {
	count int
	reset time.Time
}

// authLimiter holds the counts. One per middleware, so two groups with
// different tables do not share state.
type authLimiter struct {
	mu     sync.Mutex
	counts map[string]*authCounter
	limits map[string]AuthLimit
	// swept is when the map was last cleared of expired windows. Without it
	// the map grows by one entry per distinct address, forever, which is a
	// slow memory leak with an attacker's hand on the tap.
	swept time.Time
}

// AuthRateLimit limits the routes in the table by client address.
//
// Fixed window rather than a token bucket, deliberately: the thing being
// stopped is an automated client making hundreds of attempts, and a fixed
// window stops that with a counter and a timestamp. A burst of two at the
// window boundary is not the attack.
//
// The key is the client address plus the route path, so a limit on
// /auth/login is not spent by requests to /auth/register. The address comes
// from gin's ClientIP, which honours the trusted-proxy configuration this
// project already sets: behind a proxy that does not set it, every client
// looks like the proxy and the limit is shared, which fails closed.
//
// A path not in the table is not limited here. Add it to authRouteLimits in
// routes.go, which is the one table both this and Sentinel read.
func AuthRateLimit(limits map[string]AuthLimit) gin.HandlerFunc {
	l := &authLimiter{
		counts: make(map[string]*authCounter),
		limits: limits,
		swept:  time.Now(),
	}
	return func(c *gin.Context) {
		limit, ok := l.limits[c.Request.URL.Path]
		if !ok || limit.Requests <= 0 {
			c.Next()
			return
		}
		retry, allowed := l.take(c.ClientIP()+" "+c.Request.URL.Path, limit)
		if !allowed {
			// Retry-After in seconds, because a client that is told when to
			// come back is a client that does not retry in a loop.
			c.Header("Retry-After", itoa(int(retry.Seconds())+1))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "Too many attempts. Try again shortly.",
				},
			})
			return
		}
		c.Next()
	}
}

// take counts one request and says whether it is allowed, and if not, how long
// until the window resets.
func (l *authLimiter) take(key string, limit AuthLimit) (time.Duration, bool) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	// Every minute, drop the windows that have passed. Cheap, because it only
	// walks the map when a request arrives anyway.
	if now.Sub(l.swept) > time.Minute {
		for k, c := range l.counts {
			if now.After(c.reset) {
				delete(l.counts, k)
			}
		}
		l.swept = now
	}

	c, seen := l.counts[key]
	if !seen || now.After(c.reset) {
		l.counts[key] = &authCounter{count: 1, reset: now.Add(limit.Window)}
		return 0, true
	}
	if c.count >= limit.Requests {
		return c.reset.Sub(now), false
	}
	c.count++
	return 0, true
}

// itoa without importing strconv for one call.
func itoa(n int) string {
	if n <= 0 {
		return "1"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
`
}

// middlewareAuthLimitTestGo emits internal/middleware/authlimit_test.go.
func middlewareAuthLimitTestGo() string {
	return `package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// The limit holds, and it holds per route and per client.
func TestAuthRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limits := map[string]AuthLimit{
		"/api/v1/auth/login": {Requests: 3, Window: time.Minute},
	}
	r := gin.New()
	r.Use(AuthRateLimit(limits))
	r.POST("/api/v1/auth/login", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.POST("/api/v1/auth/other", func(c *gin.Context) { c.Status(http.StatusOK) })

	call := func(path, ip string) int {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.RemoteAddr = ip + ":1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	for i := 1; i <= 3; i++ {
		if code := call("/api/v1/auth/login", "1.2.3.4"); code != http.StatusOK {
			t.Fatalf("attempt %d answered %d, want 200", i, code)
		}
	}
	if code := call("/api/v1/auth/login", "1.2.3.4"); code != http.StatusTooManyRequests {
		t.Errorf("the fourth attempt answered %d, want 429", code)
	}

	// Another client is not limited by the first one's attempts.
	if code := call("/api/v1/auth/login", "5.6.7.8"); code != http.StatusOK {
		t.Errorf("a different address answered %d, want 200", code)
	}

	// A route not in the table is not limited.
	for i := 0; i < 10; i++ {
		if code := call("/api/v1/auth/other", "1.2.3.4"); code != http.StatusOK {
			t.Fatalf("an unlisted route answered %d, want 200", code)
		}
	}
}

// A refusal says when to come back, so a client does not retry in a loop.
func TestAuthRateLimitSaysWhenToRetry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(AuthRateLimit(map[string]AuthLimit{"/x": {Requests: 1, Window: 30 * time.Second}}))
	r.POST("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/x", nil)
		req.RemoteAddr = "9.9.9.9:1"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if i == 1 {
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("answered %d, want 429", w.Code)
			}
			if w.Header().Get("Retry-After") == "" {
				t.Error("a 429 with no Retry-After")
			}
		}
	}
}

// The window passes and the client is allowed again.
func TestAuthRateLimitWindowExpires(t *testing.T) {
	l := &authLimiter{counts: map[string]*authCounter{}, limits: map[string]AuthLimit{}, swept: time.Now()}
	limit := AuthLimit{Requests: 1, Window: 10 * time.Millisecond}

	if _, ok := l.take("k", limit); !ok {
		t.Fatal("the first request was refused")
	}
	if _, ok := l.take("k", limit); ok {
		t.Fatal("the second request inside the window was allowed")
	}
	time.Sleep(15 * time.Millisecond)
	if _, ok := l.take("k", limit); !ok {
		t.Error("the window did not reset")
	}
}

// Expired windows are dropped, so the map does not grow by one entry per
// address for the life of the process.
func TestAuthRateLimitSweepsExpiredWindows(t *testing.T) {
	l := &authLimiter{counts: map[string]*authCounter{}, limits: map[string]AuthLimit{}}
	limit := AuthLimit{Requests: 5, Window: time.Millisecond}

	for i := 0; i < 100; i++ {
		l.take(string(rune('a'+i%26))+itoa(i), limit)
	}
	if len(l.counts) != 100 {
		t.Fatalf("holding %d windows, want 100", len(l.counts))
	}

	// A sweep only runs a minute after the last one, so put that in the past.
	l.swept = time.Now().Add(-2 * time.Minute)
	time.Sleep(5 * time.Millisecond)
	l.take("fresh", limit)
	if len(l.counts) != 1 {
		t.Errorf("holding %d windows after the sweep, want 1", len(l.counts))
	}
}
`
}
