package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"commerce/apps/api/internal/realtime"
)

func TestRealtime_ChannelsOverTheSocket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := newTestAuthSvc(testCfg(), nil)
	pair, err := auth.GenerateTokenPair("user-7", "seven@example.com", "USER")
	require.NoError(t, err)
	realtime.Channel("socket-test-notes.{id}", func(c realtime.ChannelContext) bool {
		return c.UserID == "user-7" && c.Param("id") == "1"
	})

	hub := realtime.NewHub()
	r := gin.New()
	r.GET("/api/ws", NewRealtimeHandler(hub, auth).Connect)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	header := http.Header{}
	header.Set("Authorization", "Bearer "+pair.AccessToken)
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/api/ws", header)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	type message struct {
		Type    string          `json:"type"`
		Channel string          `json:"channel"`
		Payload json.RawMessage `json:"payload"`
	}
	read := func() (message, error) {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		var msg message
		err := conn.ReadJSON(&msg)
		return msg, err
	}
	send := func(kind, channel string) {
		require.NoError(t, conn.WriteJSON(map[string]string{"type": kind, "channel": channel}))
	}

	msg, err := read()
	require.NoError(t, err)
	require.Equal(t, "system.connected", msg.Type)

	send("subscribe", "private-socket-test-notes.2")
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "subscription_error", msg.Type, "a channel the authorizer rejects was subscribed")

	send("subscribe", "private-socket-test-notes.1")
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "subscribed", msg.Type)

	hub.Publish("private-socket-test-notes.1", realtime.Event{Type: "notes.updated", Payload: map[string]string{"id": "1"}})
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "notes.updated", msg.Type)
	require.Equal(t, "private-socket-test-notes.1", msg.Channel)

	send("unsubscribe", "private-socket-test-notes.1")
	msg, err = read()
	require.NoError(t, err)
	require.Equal(t, "unsubscribed", msg.Type)

	hub.Publish("private-socket-test-notes.1", realtime.Event{Type: "notes.updated"})
	_ = conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, _, err = conn.ReadMessage()
	require.Error(t, err, "a publish arrived after unsubscribe")
}
