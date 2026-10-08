package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"saas/apps/api/internal/realtime"
	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
)

const (
	wsWriteWait      = 10 * time.Second
	wsPongWait       = 60 * time.Second
	wsPingPeriod     = (wsPongWait * 9) / 10
	wsMaxMessageSize = 2048 // subscribe messages, and client events with a payload of up to realtime.MaxClientEventBytes
)

// upgrader checks the origin again, as a backstop for any other caller of
// Upgrade. Connect has already refused a cross-origin cookie handshake with a
// 403 by the time it gets here. See realtime.CheckOrigin.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     realtime.CheckOrigin,
}

// RealtimeHandler upgrades an HTTP request to a WebSocket and registers
// it with the hub. Authentication uses a query-string JWT (?token=...)
// because browsers can't set custom Authorization headers on WebSocket
// handshakes — there is no other portable way to pass the JWT.
type RealtimeHandler struct {
	Hub  *realtime.Hub
	Auth *services.AuthService
}

// NewRealtimeHandler wires the handler to the global Hub and AuthService.
func NewRealtimeHandler(hub *realtime.Hub, auth *services.AuthService) *RealtimeHandler {
	return &RealtimeHandler{Hub: hub, Auth: auth}
}

// Connect upgrades the request to a WebSocket connection.
//
//	GET /api/ws?token=<jwt>
func (h *RealtimeHandler) Connect(c *gin.Context) {
	// A native client holds a real token and passes it explicitly. A browser
	// cannot: login stores the JWT in the HttpOnly grit_access cookie so that
	// scripts cannot read it, which also means a script cannot put it in this
	// query string. The cookie rides along with the handshake GET, so read it
	// from there instead of inventing a way to hand the token to JavaScript.
	tokenStr := c.Query("token")
	if tokenStr == "" {
		if bearer := c.GetHeader("Authorization"); bearer != "" {
			tokenStr = strings.TrimPrefix(bearer, "Bearer ")
		} else if cookie, err := c.Cookie("grit_access"); err == nil {
			// The cookie is the one credential a page on another site can make
			// the browser send, so it counts only from an allowed origin.
			if !realtime.CheckOrigin(c.Request) {
				log.Printf("[ws] refused a cookie handshake from origin %q", c.GetHeader("Origin"))
				c.JSON(http.StatusForbidden, gin.H{
					"error": gin.H{"code": "FORBIDDEN", "message": "this origin may not open the realtime socket with a cookie"},
				})
				return
			}
			tokenStr = cookie
		}
	}
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "MISSING_TOKEN", "message": "a ?token query or a grit_access cookie is required"},
		})
		return
	}
	claims, err := h.Auth.ValidateAccessToken(tokenStr)
	if err != nil {
		// The parser's reason (expired, bad signature, malformed) is for the log,
		// not for whoever sent the token.
		log.Printf("[ws] refused a handshake token: %v", err)
		respond.Fail(c, respond.CodeInvalidToken, "Invalid or expired token")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}

	client := &realtime.Client{
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 32),
	}
	// The handshake is the only point at which the token is verified, so the
	// connection has to carry its own deadline. Without one, a socket opened
	// with a 15 minute token keeps delivering events indefinitely, including
	// after the session behind it has been revoked.
	if exp := claims.ExpiresAt; exp != nil {
		client.ExpiresAt = exp.Time
	}
	// Greeting so the client knows the link is live. Queued before Admit, while
	// nothing else can reach this client: once the hub holds it, DisconnectUser
	// may close Send at any moment, and a send on a closed channel panics.
	// writePump starts after Admit, so the greeting still arrives only once the
	// socket is admitted.
	greeting, _ := json.Marshal(realtime.Event{
		Type:    "system.connected",
		Payload: gin.H{"user_id": claims.UserID},
	})
	select {
	case client.Send <- greeting:
	default:
	}

	// Refused after the upgrade, not before, so the client reads a close code
	// (1013, try again later) instead of a failed handshake its reconnect loop
	// cannot tell apart from a network error.
	if err := h.Hub.Admit(client); err != nil {
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "too many connections"),
			time.Now().Add(wsWriteWait))
		_ = conn.Close()
		return
	}

	go writePump(client)
	go readPump(h.Hub, client)
}

// readPump pumps messages from the client to the hub. Mutations go through the
// REST API; what a client sends here is channel subscribe and unsubscribe, which
// the hub answers on the same socket. It also services ping/pong and cleans up
// on disconnect.
func readPump(hub *realtime.Hub, c *realtime.Client) {
	defer func() {
		hub.Unregister(c)
		_ = c.Conn.Close()
	}()
	c.Conn.SetReadLimit(wsMaxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(wsPongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		hub.HandleClientMessage(c, msg)
	}
}

// writePump pumps messages from the hub → client and emits keepalive pings.
func writePump(c *realtime.Client) {
	ticker := time.NewTicker(wsPingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			// Retire the connection once the token that authorised it has
			// expired. The client reconnects with a fresh one, which is how
			// a revoked session stops receiving: the same bound REST has.
			if !c.ExpiresAt.IsZero() && time.Now().After(c.ExpiresAt) {
				_ = c.Conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
				_ = c.Conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "token expired"))
				return
			}
			_ = c.Conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
