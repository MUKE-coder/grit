package realtime

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
