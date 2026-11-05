package scaffold

import (
	"strings"
	"testing"
)

// Three pieces of configuration that used to be package-level variables.
//
// grit#155: the generated API is otherwise explicit about its dependencies, and
// these three were the exceptions that made "no hidden globals" need a
// footnote. They are fine for one server in one process and wrong for two: the
// instance that called Setup second owned all three, so the first instance
// wrote its refresh cookie at the second's path, closed the second's sockets on
// a revocation, and checked WebSocket origins against the second's CORS list.
// Two tests that wanted different values had the same problem and had to run
// one after the other, each restoring what the other also wrote.
//
// Each one now lives on the thing that uses it, with the variable kept as the
// fallback. This test reads routes.go, which is where all three were set, and
// holds it to writing none of them.
func TestSetupConfiguresInstancesNotPackageVariables(t *testing.T) {
	read := apiSource(t)
	routes := read("internal", "routes", "routes.go")

	for _, global := range []struct{ assignment, instead string }{
		{"services.RefreshCookiePath =", "RefreshCookiePath on the AuthService literal"},
		{"services.OnSessionsRevoked =", "services.OnSessionsRevokedFor(db, ...)"},
		{"realtime.AllowedOrigins =", "realtimeHub.SetAllowedOrigins(...)"},
	} {
		if strings.Contains(routes, global.assignment) {
			t.Errorf("routes.Setup still writes the package variable %s; use %s, "+
				"or a second instance in this process overwrites the first's",
				global.assignment, global.instead)
		}
	}

	for _, want := range []string{
		`RefreshCookiePath: "/api/" + APIVersion + "/auth",`,
		"services.OnSessionsRevokedFor(db, realtimeHub.DisconnectUser)",
		"realtimeHub.SetAllowedOrigins(corsOrigins)",
	} {
		if !strings.Contains(routes, want) {
			t.Errorf("routes.Setup is missing %q", want)
		}
	}
}

// The fallbacks have to survive, because a project that was written before this
// change sets the variables from its own routes.go and gets no repair for it.
func TestPackageLevelFallbacksSurvive(t *testing.T) {
	read := apiSource(t)

	auth := read("internal", "services", "auth.go")
	for _, want := range []string{
		"var RefreshCookiePath = ",
		"func (s *AuthService) refreshCookiePath() string {",
		"return RefreshCookiePath",
	} {
		if !strings.Contains(auth, want) {
			t.Errorf("services/auth.go is missing %q, so a project that sets the "+
				"package variable loses its refresh cookie path", want)
		}
	}

	guard := read("internal", "realtime", "guard.go")
	for _, want := range []string{
		"var AllowedOrigins func() []string",
		"func OriginAllowed(r *http.Request, allowed func() []string) bool {",
		"func (h *Hub) SetAllowedOrigins(allowed func() []string) {",
		"func (h *Hub) CheckOrigin(r *http.Request) bool {",
		"return OriginAllowed(r, AllowedOrigins)",
	} {
		if !strings.Contains(guard, want) {
			t.Errorf("realtime/guard.go is missing %q, so a project that sets the "+
				"package variable stops refusing cross-origin handshakes", want)
		}
	}

	// Where the origin list lives decides which files have to agree.
	// handlers/realtime_sse.go is framework code and is rewritten on every
	// upgrade, while handlers/realtime.go is yours and an upgrade holds it back
	// once you have edited it. Both check the origin, so the check has to be
	// something the rewritten file can always reach: a method on the hub, whose
	// declaration and whose registry are both in guard.go. Reached for a method
	// on the handler instead, and `grit upgrade` on a project at the older
	// shape left it with "h.checkOrigin undefined".
	for _, file := range []string{"realtime.go", "realtime_sse.go"} {
		src := read("internal", "handlers", file)
		if !strings.Contains(src, "h.Hub.CheckOrigin(c.Request)") {
			t.Errorf("handlers/%s does not check the origin through the hub, so a "+
				"project whose handlers/realtime.go is held back no longer compiles",
				file)
		}
	}

	handler := read("internal", "handlers", "realtime.go")
	// The backstop on the upgrader has to read the same origins Connect
	// checked. A package-level upgrader could not: it had no way to reach the
	// hub.
	if !strings.Contains(handler, "CheckOrigin:     hub.CheckOrigin,") {
		t.Error("the upgrader's origin backstop is not bound to the hub, so with " +
			"two instances it checks the other one's origins")
	}
	if strings.Contains(handler, "var upgrader = websocket.Upgrader{") {
		t.Error("handlers/realtime.go still has a package-level upgrader")
	}

	if !strings.Contains(guard, "return OriginAllowed(r, AllowedOrigins)") {
		t.Error("Hub.CheckOrigin does not fall back to realtime.AllowedOrigins, so " +
			"an upgraded project's own origin list is ignored and its admin " +
			"panel can no longer open a socket")
	}

	session := read("internal", "services", "session.go")
	if !strings.Contains(session, "hook = OnSessionsRevoked") {
		t.Error("services/session.go does not fall back to the package-level hook, " +
			"so an upgraded project's revocations stop closing sockets")
	}
}
