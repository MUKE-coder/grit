package scaffold

// The realtime socket's origin rule and connection caps.
//
// CheckOrigin returned true for every origin while the handshake authenticates
// with the grit_access cookie, so a page on another origin could open a socket
// as the signed-in user and read their events. MODULE_REALTIME=false hid the UI
// and left /api/ws mounted, and nothing bounded how many sockets one account
// could hold. The constants below are shared by the templates and by the
// upgrade repair, so both write the same text.

// apiRealtimeGuardGo emits internal/realtime/guard.go.
func apiRealtimeGuardGo() string {
	return `package realtime

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
`
}

// apiRealtimeGuardTestGo emits internal/realtime/guard_test.go.
func apiRealtimeGuardTestGo() string {
	return `package realtime

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func withAllowedOrigins(t *testing.T, origins ...string) {
	t.Helper()
	previous := AllowedOrigins
	AllowedOrigins = func() []string { return origins }
	t.Cleanup(func() { AllowedOrigins = previous })
}

func TestCheckOriginRefusesACookieHandshakeFromAnotherOrigin(t *testing.T) {
	withAllowedOrigins(t, "http://localhost:3001")
	r := httptest.NewRequest("GET", "http://api.example.com/api/ws", nil)
	r.Header.Set("Origin", "http://evil.example")
	if CheckOrigin(r) {
		t.Fatal("a page on another origin may open the socket with the user's cookie")
	}
	r.Header.Set("Origin", "null")
	if CheckOrigin(r) {
		t.Fatal("an opaque origin (sandboxed frame, file://) may open the socket with the user's cookie")
	}
}

func TestCheckOriginAcceptsAListedOrigin(t *testing.T) {
	withAllowedOrigins(t, "http://localhost:3000", "http://localhost:3001/")
	r := httptest.NewRequest("GET", "http://localhost:8080/api/ws", nil)
	r.Header.Set("Origin", "http://localhost:3001")
	if !CheckOrigin(r) {
		t.Fatal("the admin panel's own origin is refused")
	}
}

func TestCheckOriginAcceptsTheAPIsOwnHost(t *testing.T) {
	withAllowedOrigins(t)
	r := httptest.NewRequest("GET", "http://app.example.com/api/ws", nil)
	r.Header.Set("Origin", "https://app.example.com")
	if !CheckOrigin(r) {
		t.Fatal("a page served by the API itself is refused")
	}
}

func TestCheckOriginIgnoresAWildcard(t *testing.T) {
	withAllowedOrigins(t, "*")
	r := httptest.NewRequest("GET", "http://api.example.com/api/ws", nil)
	r.Header.Set("Origin", "http://evil.example")
	if CheckOrigin(r) {
		t.Fatal("CORS_ORIGINS=* hands the cookie socket to every site")
	}
}

func TestCheckOriginAcceptsATokenClient(t *testing.T) {
	withAllowedOrigins(t)
	native := httptest.NewRequest("GET", "http://api.example.com/api/ws", nil)
	native.Header.Set("Authorization", "Bearer abc")
	if !CheckOrigin(native) {
		t.Fatal("a native client with no Origin is refused")
	}
	desktop := httptest.NewRequest("GET", "http://api.example.com/api/ws?token=abc", nil)
	desktop.Header.Set("Origin", "wails://wails")
	if !CheckOrigin(desktop) {
		t.Fatal("a client that sends its own token is refused for its origin")
	}
}

func testClient(userID string) *Client {
	return &Client{UserID: userID, Send: make(chan []byte, 1)}
}

func TestAdmitRefusesTheEleventhSocketForOneUser(t *testing.T) {
	SetConnectionLimits(10, 0)
	t.Cleanup(func() { SetConnectionLimits(10, 10000) })
	hub := NewHub()
	for i := 0; i < 10; i++ {
		if err := hub.Admit(testClient("u1")); err != nil {
			t.Fatalf("socket %d refused: %v", i+1, err)
		}
	}
	if err := hub.Admit(testClient("u1")); !errors.Is(err, ErrTooManyConnections) {
		t.Fatalf("the 11th socket for one user was admitted (err=%v)", err)
	}
	if err := hub.Admit(testClient("u2")); err != nil {
		t.Fatalf("another user is refused because the first reached their cap: %v", err)
	}
}

func TestAdmitHonoursTheProcessCap(t *testing.T) {
	SetConnectionLimits(0, 3)
	t.Cleanup(func() { SetConnectionLimits(10, 10000) })
	hub := NewHub()
	for _, user := range []string{"a", "b", "c"} {
		if err := hub.Admit(testClient(user)); err != nil {
			t.Fatalf("%s refused under the cap: %v", user, err)
		}
	}
	if err := hub.Admit(testClient("d")); !errors.Is(err, ErrTooManyConnections) {
		t.Fatalf("a connection past the process cap was admitted (err=%v)", err)
	}
}
`
}

// apiRealtimeHandlerTestGo emits internal/handlers/realtime_test.go: the
// handshake end to end, over a real socket.
func apiRealtimeHandlerTestGo() string {
	return `package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"{{MODULE}}/internal/realtime"
)

// realtimeServer mounts Connect on a real listener, with CORS allowing only the
// admin panel's origin, and returns the socket URL and a valid access token.
func realtimeServer(t *testing.T) (string, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	auth := newTestAuthSvc(testCfg())
	pair, err := auth.GenerateTokenPair("user-1", "user@example.com", "USER")
	require.NoError(t, err)

	previous := realtime.AllowedOrigins
	realtime.AllowedOrigins = func() []string { return []string{"http://localhost:3001"} }
	t.Cleanup(func() { realtime.AllowedOrigins = previous })

	r := gin.New()
	r.GET("/api/ws", NewRealtimeHandler(realtime.NewHub(), auth).Connect)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws", pair.AccessToken
}

func dialRealtime(t *testing.T, url string, header http.Header) (*websocket.Conn, int) {
	t.Helper()
	conn, resp, err := websocket.DefaultDialer.Dial(url, header)
	status := 0
	if resp != nil {
		status = resp.StatusCode
	}
	if err != nil {
		return nil, status
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn, status
}

// readGreeting waits for system.connected, which Connect sends only once the
// hub has admitted the socket.
func readGreeting(t *testing.T, conn *websocket.Conn) error {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err == nil {
		require.Contains(t, string(msg), "system.connected")
	}
	return err
}

func TestRealtime_CrossOriginCookieHandshakeIsRefused(t *testing.T) {
	url, token := realtimeServer(t)
	header := http.Header{}
	header.Set("Origin", "http://evil.example")
	header.Set("Cookie", "grit_access="+token)
	conn, status := dialRealtime(t, url, header)
	require.Nil(t, conn, "a page on another origin opened the socket with the user's cookie")
	require.Equal(t, http.StatusForbidden, status)
}

func TestRealtime_ListedOriginCookieHandshakeIsAccepted(t *testing.T) {
	url, token := realtimeServer(t)
	header := http.Header{}
	header.Set("Origin", "http://localhost:3001")
	header.Set("Cookie", "grit_access="+token)
	conn, status := dialRealtime(t, url, header)
	require.NotNil(t, conn, "the admin panel's origin was refused (status %d)", status)
	require.NoError(t, readGreeting(t, conn))
}

func TestRealtime_BearerClientWithoutOriginIsAccepted(t *testing.T) {
	url, token := realtimeServer(t)
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	conn, status := dialRealtime(t, url, header)
	require.NotNil(t, conn, "a native bearer client was refused (status %d)", status)
	require.NoError(t, readGreeting(t, conn))
}

func TestRealtime_PerUserCapClosesTheEleventhSocket(t *testing.T) {
	realtime.SetConnectionLimits(10, 0)
	t.Cleanup(func() { realtime.SetConnectionLimits(10, 10000) })
	url, token := realtimeServer(t)
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	for i := 0; i < 10; i++ {
		conn, status := dialRealtime(t, url, header)
		require.NotNil(t, conn, "socket %d refused (status %d)", i+1, status)
		require.NoError(t, readGreeting(t, conn))
	}
	conn, _ := dialRealtime(t, url, header)
	require.NotNil(t, conn)
	err := readGreeting(t, conn)
	require.True(t, websocket.IsCloseError(err, websocket.CloseTryAgainLater),
		"the 11th socket was not closed with 1013 (err=%v)", err)
}
`
}

// Handler hunks. The Old text is exactly what Grit wrote before the fix.

const realtimeHandlerImportOld = "\t\"net/http\"\n\t\"time\"\n"
const realtimeHandlerImportNew = "\t\"net/http\"\n\t\"strings\"\n\t\"time\"\n"

const realtimeUpgraderOld = `// upgrader allows any origin — desktop clients use Wails (file://) and
// the API is mounted behind CORS that already restricts origins for
// regular HTTP traffic.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
`

const realtimeUpgraderNew = `// upgrader checks the origin again, as a backstop for any other caller of
// Upgrade. Connect has already refused a cross-origin cookie handshake with a
// 403 by the time it gets here. See realtime.CheckOrigin.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     realtime.CheckOrigin,
}
`

const realtimeTokenOld = `	tokenStr := c.Query("token")
	if tokenStr == "" {
		if cookie, err := c.Cookie("grit_access"); err == nil {
			tokenStr = cookie
		}
	}
`

const realtimeTokenNew = `	tokenStr := c.Query("token")
	if tokenStr == "" {
		if bearer := c.GetHeader("Authorization"); bearer != "" {
			tokenStr = strings.TrimPrefix(bearer, "Bearer ")
		} else if cookie, err := c.Cookie("grit_access"); err == nil {
			// The cookie is the one credential a page on another site can make
			// the browser send, so it counts only from an allowed origin.
			if !realtime.CheckOrigin(c.Request) {
				log.Printf("[ws] refused a cookie handshake from origin %q", c.GetHeader("Origin"))
				c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "this origin may not open the realtime socket with a cookie"}})
				return
			}
			tokenStr = cookie
		}
	}
`

const realtimeRegisterOld = "\th.Hub.Register(client)\n"

const realtimeRegisterNew = `	// Refused after the upgrade, not before, so the client reads a close code
	// (1013, try again later) instead of a failed handshake its reconnect loop
	// cannot tell apart from a network error.
	if err := h.Hub.Admit(client); err != nil {
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "too many connections"),
			time.Now().Add(wsWriteWait))
		_ = conn.Close()
		return
	}
`

// routes.go hunks.

const routesCORSOld = `	r.Use(middleware.CORSDynamic(func() []string {
		if stored := settings.String(context.Background(), "cors.origins"); strings.TrimSpace(stored) != "" {
			return splitOrigins(stored)
		}
		return cfg.CORSOrigins
	}))
`

const routesCORSNew = `	corsOrigins := func() []string {
		if stored := settings.String(context.Background(), "cors.origins"); strings.TrimSpace(stored) != "" {
			return splitOrigins(stored)
		}
		return cfg.CORSOrigins
	}
	r.Use(middleware.CORSDynamic(corsOrigins))
`

const routesRealtimeHubOld = "\trealtimeHub := realtime.NewHub(realtime.WithRedis(cfg.RedisURL, \"\"))\n"

const routesRealtimeHubNew = `	// MODULE_REALTIME=false turns realtime off: no Redis subscription here and
	// no /api/ws route below. The hub is still built, empty, so code that pushes
	// to it needs no nil check, and with nobody able to connect a push reaches
	// no one.
	var realtimeOptions []realtime.Option
	if cfg.Modules.Realtime {
		realtimeOptions = append(realtimeOptions, realtime.WithRedis(cfg.RedisURL, ""))
	}
	realtimeHub := realtime.NewHub(realtimeOptions...)
`

const routesRealtimeHandlerOld = "\trealtimeHandler := handlers.NewRealtimeHandler(realtimeHub, authService)\n"

const routesRealtimeHandlerNew = routesRealtimeHandlerOld + `	// A browser handshake authenticates with the grit_access cookie, which a
	// page on any site can make the browser send, so it is accepted only from
	// the origins CORS allows.
	realtime.AllowedOrigins = corsOrigins
`

const routesRealtimeRouteOld = `	// WebSocket: realtime hub. Auth via ?token=<jwt> on the handshake
	// because browsers can't set custom headers on WS upgrade.
	r.GET("/api/ws", realtimeHandler.Connect)
`

const routesRealtimeRouteNew = `	// WebSocket: realtime hub. Browsers authenticate with the grit_access
	// cookie from an allowed origin, native clients with ?token=<jwt> or an
	// Authorization header. Not mounted when MODULE_REALTIME=false, so the
	// path answers 404.
	if cfg.Modules.Realtime {
		r.GET("/api/ws", realtimeHandler.Connect)
	}
`
