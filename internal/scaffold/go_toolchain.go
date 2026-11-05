package scaffold

// The Go toolchain versions Grit knows to be free of the standard-library
// vulnerabilities a generated project reaches.
//
// Two things need these and used to carry their own copies: the go directive
// in a generated go.mod, and `grit doctor`, which had no opinion about the
// toolchain at all. A project pinned to a safe Go in its Dockerfile was still
// running every local test and every desktop build on whatever `go` was on
// PATH, and nothing said so (grit#178).

// GoToolchainPin is the go directive a generated go.mod carries.
//
// 1.26.10, and not an earlier patch. Earlier ones carry standard library
// vulnerabilities this code reaches before anyone has signed in: Go 1.26.4 and
// below, eight of them, among them net/http and the encoding/xml behind the
// SAML sign-in path; Go 1.26.8 and below, GO-2026-6617, an HPACK encoder race
// that crashes an HTTP/2 server.
//
// It is also what CI, setup-go and the release workflow install, and an older
// Go downloads this one, so raising it here raises it everywhere.
const GoToolchainPin = "1.26.10"

// GoPatchFloors is the lowest patch of each Go minor that is free of the
// standard-library advisories a generated project is exposed to.
//
// A MINOR a reader is not on is not listed, and a minor newer than this
// release knows about is not warned on: claiming a vulnerability in a Go that
// did not exist when this shipped would be a guess, and a check that guesses
// is a check people switch off.
//
// Refreshing this is part of the release ritual when a Go security release
// lands. The newest patch of each supported minor is at
// https://go.dev/dl/?mode=json, and every entry here is the newest patch at
// the time of the release rather than the oldest safe one: Go's security
// releases do not backport past the two supported minors, so "newest" and
// "safe" are the same answer.
var GoPatchFloors = map[string]int{
	// go1.25.14, 2026. The oldest minor still receiving security fixes.
	"1.25": 14,
	// go1.26.10, which is GoToolchainPin above.
	"1.26": 10,
	// go1.27.3 fixes GO-2026-6809 in crypto/x509, which 1.27.2 and below carry.
	"1.27": 3,
}
