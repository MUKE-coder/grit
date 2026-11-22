package selfupdate

import (
	"os"
	"testing"
)

// The bytes the site actually serves, decoded by the struct the CLI uses.
//
// Set GRIT_CHANNEL_FIXTURE to a file holding a real /stable.json response.
// Skipped without it: this exists to check one against the other by hand
// before a release, not to make the suite need a server.
func TestTheServedChannelDecodes(t *testing.T) {
	path := os.Getenv("GRIT_CHANNEL_FIXTURE")
	if path == "" {
		t.Skip("set GRIT_CHANNEL_FIXTURE to a captured /stable.json")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the fixture: %v", err)
	}
	withChannelBytes(t, raw)

	got, err := ResolveChannel("stable")
	if err != nil {
		t.Fatalf("ResolveChannel: %v", err)
	}
	if got == "" {
		t.Fatal("the served channel decoded to no version")
	}
	t.Logf("the served channel resolves to v%s", got)
}
