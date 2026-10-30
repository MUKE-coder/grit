package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Review 2.18 and what it would have broken.
//
// A response that reflects the Origin header has to say so, or a shared cache
// in front of the API can hand one origin's Access-Control-Allow-Origin to
// another. With credentials enabled that header is what decides whether the
// browser lets the page read the body.
//
// And the URL cache refuses to store a response whose Vary names anything its
// key does not include, which is correct and which means adding that header
// switched the cache off for every response at once. Silently. These two
// tests exist together because the pair is the thing that has to hold.

// Read from the template source rather than a scaffold, so the test is fast
// and does not depend on a working toolchain.
func TestCORSSetsVaryOriginBeforeTheDecision(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("api_files.go"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)

	vary := strings.Index(s, `c.Header("Vary", "Origin")`)
	if vary < 0 {
		t.Fatal("the CORS middleware does not set Vary: Origin, so a shared " +
			"cache can serve one origin's allow header to another (review 2.18)")
	}
	allow := strings.Index(s, `c.Header("Access-Control-Allow-Origin", origin)`)
	if allow < 0 {
		t.Fatal("the CORS middleware no longer reflects the origin; this test " +
			"checks nothing")
	}
	// Before the branch, not inside it: a refusal is also a response that
	// varies by origin, and caching it without Vary would replay the refusal
	// to an origin that is allowed.
	if vary > allow {
		t.Error("Vary is set after the allow decision, so a refused origin's " +
			"response is cacheable without it")
	}
}

func TestTheURLCacheToleratesVaryOrigin(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("pool_cache_repair.go"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)

	// Both copies: the one new projects get and the one the upgrade writes
	// into projects that exist. One without the other is a cache that works
	// in a fresh project and silently stops in an upgraded one.
	if n := strings.Count(s, `case "Origin":`); n != 2 {
		t.Errorf(`shareableResponse allows Vary: Origin in %d of 2 copies. `+
			`Every CORS response carries that header, so a copy without this `+
			`case never caches anything again`, n)
	}
}
