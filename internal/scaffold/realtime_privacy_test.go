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
