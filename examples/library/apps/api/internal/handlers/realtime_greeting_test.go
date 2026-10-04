package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"library/apps/api/internal/realtime"
)

// disconnectOnRegister signs the user out the moment the hub logs that it
// admitted their socket, which is the instant Connect used to send its greeting.
type disconnectOnRegister struct {
	hub *realtime.Hub
	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *disconnectOnRegister) Write(p []byte) (int, error) {
	if i := bytes.Index(p, []byte("client registered user=")); i >= 0 {
		user := strings.Fields(string(p[i+len("client registered user="):]))[0]
		go w.hub.DisconnectUser(user)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

// A session revoked while a socket is being admitted must not panic Connect.
//
// Connect sent its greeting to client.Send after Admit. From Admit on, the hub
// holds the client, and DisconnectUser closes Send; a send on a closed channel
// panics. The greeting is now queued before Admit, while nothing else can reach
// the client.
func TestRealtime_GreetingDoesNotRaceADisconnect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	realtime.SetConnectionLimits(0, 0)
	t.Cleanup(func() { realtime.SetConnectionLimits(10, 10000) })
	auth := newTestAuthSvc(testCfg(), nil)
	pair, err := auth.GenerateTokenPair("user-g", "greeting@example.com", "USER")
	require.NoError(t, err)

	hub := realtime.NewHub()
	log.SetOutput(&disconnectOnRegister{hub: hub})
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	var panics atomic.Int64
	r := gin.New()
	r.Use(func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				panics.Add(1)
			}
		}()
		c.Next()
	})
	r.GET("/api/ws", NewRealtimeHandler(hub, auth).Connect)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws"
	header := http.Header{}
	header.Set("Authorization", "Bearer "+pair.AccessToken)
	const sockets = 300
	for i := 0; i < sockets; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)
		// Read until the server closes it: the greeting, then the sign-out. A
		// Connect that panicked leaves the socket open, hence the short deadline.
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
		_ = conn.Close()
	}
	t.Logf("%d sockets signed out as they were admitted: %d panics", sockets, panics.Load())
	require.Zero(t, panics.Load(), "Connect panicked sending to a connection DisconnectUser had closed")
}
