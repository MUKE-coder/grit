package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// apiSource scaffolds a triple-tier API once and hands back a reader for the
// generated Go files.
func apiSource(t *testing.T) func(parts ...string) string {
	t.Helper()
	root := t.TempDir()
	opts := Options{ProjectName: "chat", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAPIFiles(root, opts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}
	// upload.go comes from the storage writer, not the API writer.
	if err := writeStorageFiles(root, opts); err != nil {
		t.Fatalf("writeStorageFiles: %v", err)
	}
	return func(parts ...string) string {
		full := append([]string{root, "apps", "api"}, parts...)
		b, err := os.ReadFile(filepath.Join(full...))
		if err != nil {
			t.Fatalf("read %s: %v", filepath.Join(parts...), err)
		}
		return string(b)
	}
}

// A resource event must never go to every connected client.
//
// registerRealtime used to subscribe to "*" and call hub.Broadcast, so every
// create, update and delete of every resource was pushed to every open socket
// in the process. The payload carries e.Label, the row's human-readable title,
// so an unrelated signed-in user received things like
//
//	{"type":"conversations.created","label":"Chemo results - family only", ...}
//
// for records they had no right to read, across tenants. Reproduced against a
// stock scaffold with two accounts and one websocket before this was changed.
func TestResourceEventsAreNotBroadcastToEveryone(t *testing.T) {
	read := apiSource(t)
	src := read("internal", "services", "event_subscribers.go")

	if strings.Contains(src, "hub.Broadcast(") {
		t.Error("registerRealtime broadcasts resource events to every connected " +
			"client; every signed-in user then learns the label of every row " +
			"written anywhere in the system")
	}
	if !strings.Contains(src, "RealtimeAudience") {
		t.Error("no RealtimeAudience hook: there is no way for an app to say who " +
			"is allowed to see a resource event")
	}
	if !strings.Contains(src, "hub.SendToUsers(audience") {
		t.Error("registerRealtime should deliver to the audience, not to everyone")
	}
}

// A websocket must not outlive the token that authorised it.
//
// The handshake is the only point the JWT is checked. Without a deadline on
// the connection, a socket opened with a 15 minute token kept delivering
// events indefinitely, including after the session behind it had been revoked.
func TestWebSocketConnectionCarriesItsTokenExpiry(t *testing.T) {
	read := apiSource(t)

	hub := read("internal", "realtime", "hub.go")
	if !strings.Contains(hub, "ExpiresAt time.Time") {
		t.Error("realtime.Client has no ExpiresAt, so nothing bounds the life of " +
			"a connection to the life of its token")
	}

	h := read("internal", "handlers", "realtime.go")
	if !strings.Contains(h, "client.ExpiresAt = exp.Time") {
		t.Error("Connect does not copy the token expiry onto the client")
	}
	if !strings.Contains(h, "time.Now().After(c.ExpiresAt)") {
		t.Error("writePump never checks ExpiresAt, so an expired token keeps its " +
			"socket open forever")
	}
}

// Revoking a session has to close that user's live sockets.
//
// "Sign out of all devices" marked rows revoked and nothing else, so the
// signed-out device kept receiving message bodies over its already-open
// websocket. Verified before the fix: after a successful revoke-all the
// revoked connection still received every subsequent event.
func TestRevokingSessionsClosesWebSockets(t *testing.T) {
	read := apiSource(t)

	hub := read("internal", "realtime", "hub.go")
	if !strings.Contains(hub, "func (h *Hub) DisconnectUser(") {
		t.Fatal("the hub cannot close a user's connections, so revocation has no " +
			"way to reach an open socket")
	}

	sess := read("internal", "services", "session.go")
	for _, want := range []string{
		"var OnSessionsRevoked func(userID string)",
		"sessionsRevoked(userID)",
	} {
		if !strings.Contains(sess, want) {
			t.Errorf("session service is missing %q, so revocation cannot notify "+
				"the realtime layer", want)
		}
	}
	// Both revocation paths, not just the bulk one.
	if n := strings.Count(sess, "sessionsRevoked(userID)"); n < 2 {
		t.Errorf("sessionsRevoked is called %d time(s); both RevokeSession and "+
			"RevokeAllUserSessions need it", n)
	}

	routes := read("internal", "routes", "routes.go")
	if !strings.Contains(routes, "services.OnSessionsRevoked = realtimeHub.DisconnectUser") {
		t.Error("routes.Setup never wires the revocation hook to the hub, so the " +
			"hook exists but nothing calls into it")
	}
}

// A file field declared accepts:"audio" must actually accept audio.
//
// The admin's file-accepts lib has had an "audio" group listing audio/mpeg,
// audio/wav, audio/ogg, audio/x-m4a and audio/webm for a long time, and the
// storage dashboard buckets uploads by mime_type LIKE 'audio/%'. The upload
// allowlist had no audio entry at all, so the picker offered audio files and
// every one of them came back INVALID_FILE_TYPE. Voice notes were impossible
// without editing framework code. This is the same trap the archive types
// were added to close.
func TestAudioIsUploadable(t *testing.T) {
	read := apiSource(t)
	src := read("internal", "handlers", "upload.go")

	for _, mime := range []string{
		"audio/webm", // what MediaRecorder produces in Chrome
		"audio/mp4",  // what it produces in Safari
		"audio/ogg",
		"audio/mpeg",
		"audio/wav",
		"audio/x-m4a", // listed by the admin's audio accept group
	} {
		if !strings.Contains(src, `"`+mime+`"`) {
			t.Errorf("%s is not uploadable, so a field with accepts:\"audio\" "+
				"offers a picker whose files the server then refuses", mime)
		}
	}
}

// A browser must be able to open the socket, which means the cookie.
//
// The documented way to connect is ?token=<jwt> read from storage. The web app
// has no readable token: login sets grit_access HttpOnly precisely so scripts
// cannot read it, so the documented call could not be written in the app Grit
// generates, and nothing in apps/web ever connected. The handshake is a plain
// GET, so the cookie arrives with it.
func TestWebSocketAcceptsTheAuthCookie(t *testing.T) {
	read := apiSource(t)
	h := read("internal", "handlers", "realtime.go")

	if !strings.Contains(h, `c.Cookie("grit_access")`) {
		t.Error("Connect only reads ?token=, which a browser cannot supply: the " +
			"access token is in an HttpOnly cookie by design")
	}
	// The query parameter still comes first, for Expo and service clients that
	// hold a real token and have no cookie jar.
	if !strings.Contains(h, `tokenStr := c.Query("token")`) {
		t.Error("explicit tokens are no longer accepted, which breaks every " +
			"non-browser client")
	}
}

// The client half has to exist, in every frontend the project has.
func TestRealtimeClientShipsWithTheFrontends(t *testing.T) {
	for _, frontend := range []Frontend{FrontendNext, FrontendTanStack} {
		root := t.TempDir()
		opts := Options{
			ProjectName: "chat", Architecture: ArchTriple,
			Frontend: frontend, IncludeExpo: true,
		}
		if err := createDirectories(root, opts); err != nil {
			t.Fatalf("%s: createDirectories: %v", frontend, err)
		}
		if err := writeRealtimeClientFiles(root, opts); err != nil {
			t.Fatalf("%s: %v", frontend, err)
		}

		adminHooks := filepath.Join("apps", "admin", "hooks", "use-realtime.ts")
		if frontend == FrontendTanStack {
			adminHooks = filepath.Join("apps", "admin", "src", "hooks", "use-realtime.ts")
		}
		for _, rel := range []string{
			filepath.Join("apps", "web", "lib", "realtime.ts"),
			filepath.Join("apps", "web", "hooks", "use-realtime.ts"),
			adminHooks,
			filepath.Join("apps", "expo", "lib", "realtime.ts"),
		} {
			if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
				t.Errorf("%s: missing %s", frontend, filepath.ToSlash(rel))
			}
		}

		// Reconnection is the part a copy-pasted snippet never has, and the
		// part that decides whether live updates survive a laptop lid.
		b, err := os.ReadFile(filepath.Join(root, "apps", "web", "lib", "realtime.ts"))
		if err != nil {
			t.Fatalf("%s: %v", frontend, err)
		}
		src := string(b)
		for _, want := range []string{"backoffDelay", "Math.random()", "scheduleReconnect"} {
			if !strings.Contains(src, want) {
				t.Errorf("%s: the client has no %s, so it does not reconnect "+
					"with jitter and every client retries in lockstep", frontend, want)
			}
		}

		// "use client" belongs to Next only; Vite and Metro warn about a stray one.
		expo, _ := os.ReadFile(filepath.Join(root, "apps", "expo", "hooks", "use-realtime.ts"))
		if strings.Contains(string(expo), `"use client"`) {
			t.Errorf("%s: the Expo hook carries a Next-only directive", frontend)
		}
	}
}

// Realtime must survive a second replica.
//
// The Hub is an in-process registry, so without a backplane a user connected
// to replica A never receives an event published on replica B: the push
// succeeds into a registry that does not contain them, and nothing errors or
// logs. That is invisible in development, invisible on one instance, and shows
// up as "messages sometimes do not arrive" during the first rolling deploy
// where two versions overlap.
func TestRealtimeHasACrossProcessBackplane(t *testing.T) {
	read := apiSource(t)

	bp := read("internal", "realtime", "backplane.go")
	for _, want := range []string{
		"type Backplane interface",
		"func RedisBackplane(",
		"func (h *Hub) publish(",
		"func (h *Hub) receive(",
	} {
		if !strings.Contains(bp, want) {
			t.Errorf("the backplane is missing %q", want)
		}
	}
	// A node has to ignore its own messages, or every user sees their own
	// events twice: once locally, once off the backplane.
	if !strings.Contains(bp, "f.Origin == h.nodeID") {
		t.Error("no origin check: a node would re-deliver its own published events")
	}
	// Revocation has to cross the wire, or "sign out of all devices" only
	// reaches the replica that served the request.
	if !strings.Contains(bp, "case f.Kick != \"\":") {
		t.Error("a kick does not cross the backplane, so a revoked socket on " +
			"another replica stays open and keeps receiving")
	}

	hub := read("internal", "realtime", "hub.go")
	for _, want := range []string{"func WithBackplane(", "func WithRedis(", "h.publish(fanout{Kick: userID})"} {
		if !strings.Contains(hub, want) {
			t.Errorf("hub.go is missing %q", want)
		}
	}
	// Local delivery must not depend on the backplane, so a Redis outage
	// degrades realtime to one replica rather than breaking it.
	if !strings.Contains(hub, "h.deliverLocal(userIDs, bytes)") {
		t.Error("SendToUsers does not deliver locally before publishing")
	}

	// And the wiring, so a project with Redis gets it without being asked.
	routes := read("internal", "routes", "routes.go")
	if !strings.Contains(routes, "realtime.NewHub(realtime.WithRedis(cfg.RedisURL") {
		t.Error("routes.Setup does not pass Redis to the hub, so every project " +
			"stays single-process no matter what it has configured")
	}
}
