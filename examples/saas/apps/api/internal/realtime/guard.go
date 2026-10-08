package realtime

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// AllowedOrigins lists the browser origins that may open a socket with the
// grit_access cookie. routes.Setup points it at the list CORS uses, so an origin
// added in the admin reaches both. While it is nil, only a page served from the
// API's own host qualifies.
var AllowedOrigins func() []string

// CheckOrigin reports whether a WebSocket handshake may go ahead.
//
// A browser sends the grit_access cookie with a handshake whichever page starts
// it, and the same-origin policy does not cover WebSockets. Accepting every
// origin let a page on another origin open a socket as the signed-in user and
// read their events. So a handshake that authenticates with the cookie has to
// come from an allowed origin, or from the API's own host.
//
// A handshake that carries its own token (?token= or an Authorization header)
// is not at risk: a page that already holds the token has nothing to hijack,
// and Connect never falls back to the cookie when one is sent. Native clients
// send no Origin at all, and the desktop webview sends its own scheme with a
// token, so both connect as before.
func CheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" || r.URL.Query().Get("token") != "" || r.Header.Get("Authorization") != "" {
		return true
	}
	if u, err := url.Parse(origin); err == nil && u.Host != "" && strings.EqualFold(u.Host, r.Host) {
		return true
	}
	if AllowedOrigins == nil {
		return false
	}
	for _, allowed := range AllowedOrigins() {
		// A wildcard in CORS_ORIGINS is not a reason to hand the cookie to every
		// site. The origin has to be named.
		if allowed != "*" && strings.EqualFold(strings.TrimRight(allowed, "/"), origin) {
			return true
		}
	}
	return false
}

// ErrTooManyConnections is what Admit returns when a connection cap is reached.
var ErrTooManyConnections = errors.New("realtime: too many connections")

// Connection caps, read once from REALTIME_MAX_CONNECTIONS_PER_USER (default 10)
// and REALTIME_MAX_CONNECTIONS (default 10000, per process). Zero or less turns
// a cap off. Without them one account, or one script holding a stolen token,
// can open sockets until the process runs out of file descriptors, and realtime
// goes down for everyone.
var (
	limitsOnce sync.Once
	maxPerUser atomic.Int64
	maxTotal   atomic.Int64
)

// loadLimits reads the caps on first use rather than at package init, which
// runs before main has loaded .env.
func loadLimits() {
	limitsOnce.Do(func() {
		maxPerUser.Store(envLimit("REALTIME_MAX_CONNECTIONS_PER_USER", 10))
		maxTotal.Store(envLimit("REALTIME_MAX_CONNECTIONS", 10000))
	})
}

// SetConnectionLimits overrides the caps, for tests and for code that reads them
// from somewhere other than the environment.
func SetConnectionLimits(perUser, total int) {
	loadLimits()
	maxPerUser.Store(int64(perUser))
	maxTotal.Store(int64(total))
}

func envLimit(key string, fallback int64) int64 {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		log.Printf("[realtime] %s=%q is not a number, using %d", key, raw, fallback)
		return fallback
	}
	return n
}

// Admit registers a client unless a connection cap is reached, in which case it
// registers nothing and returns ErrTooManyConnections. The check and the insert
// happen under one lock, so a burst of handshakes cannot all pass the check
// before any of them is counted.
func (h *Hub) Admit(c *Client) error {
	loadLimits()
	perUser, total := int(maxPerUser.Load()), int(maxTotal.Load())

	h.mu.Lock()
	defer h.mu.Unlock()
	set := h.clients[c.UserID]
	if perUser > 0 && len(set) >= perUser {
		log.Printf("[realtime] refused a connection for user=%s: the cap of %d per user is reached", c.UserID, perUser)
		return ErrTooManyConnections
	}
	if total > 0 {
		open := 0
		for _, s := range h.clients {
			open += len(s)
		}
		if open >= total {
			log.Printf("[realtime] refused a connection for user=%s: the cap of %d connections is reached", c.UserID, total)
			return ErrTooManyConnections
		}
	}
	if set == nil {
		set = make(map[*Client]struct{})
		h.clients[c.UserID] = set
	}
	set[c] = struct{}{}
	log.Printf("[realtime] client registered user=%s total=%d", c.UserID, len(set))
	return nil
}
