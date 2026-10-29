package scaffold

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The k6 suite and the router are two things that have to agree, and nothing
// compared them. The shipped journey asked for /health, which no Grit API has
// ever served, and /api/blogs, which only a project with a Blog resource
// serves. Run against a real project every request 404'd, every threshold
// breached, and all six load tests failed before measuring anything.
//
// Worse than failing: one of its checks was "body non-empty", which an error
// response satisfies. A check that passes on a 500 is not a check.
//
// So: every path the journey asks for, and the line in the routes template
// that proves the API answers it. A new request in the journey fails here
// until it is added, and a route that moves fails here too.
var k6JourneyPaths = map[string]string{
	`/api/health`:        `r.GET("/api/health", healthCheck)`,
	`/api/profile`:       `profile := protected.Group("/profile")`,
	`/api/auth/login`:    `auth.POST("/login", authHandler.Login)`,
	`/api/auth/register`: `auth.POST("/register", authHandler.Register)`,
}

// http.get(BASE_URL + '/api/health') and the like. The literal has to be the
// whole path, so the environment-driven resource read, which continues with
// " + __ENV.RESOURCE", is not matched here: there is no fixed path to check.
var k6PathLiteral = regexp.MustCompile(`http\.(?:get|post|put|del)\(\s*BASE_URL \+ '([^']+)'\s*[,)]`)

func TestK6JourneyAsksForPathsTheAPIServes(t *testing.T) {
	journey := k6CommonJS()
	routes := apiRoutesGo()

	found := map[string]bool{}
	for _, m := range k6PathLiteral.FindAllStringSubmatch(journey, -1) {
		path := strings.SplitN(m[1], "?", 2)[0]
		found[path] = true
	}
	if len(found) == 0 {
		t.Fatal("no BASE_URL + '<path>' requests found in the k6 journey: " +
			"either the journey stopped making requests, or it started " +
			"building paths in a way this test cannot see, which is how the " +
			"last set of wrong paths went unnoticed")
	}

	for path := range found {
		anchor, ok := k6JourneyPaths[path]
		if !ok {
			t.Errorf("the k6 journey requests %s, which is not in "+
				"k6JourneyPaths. Add it there with the line in the routes "+
				"template that shows the API serves it, or stop requesting it.",
				path)
			continue
		}
		if !strings.Contains(routes, anchor) {
			t.Errorf("the k6 journey requests %s, and the routes template no "+
				"longer contains %q. Either the route moved and the journey "+
				"now 404s, or the anchor needs updating.", path, anchor)
		}
	}

	for path := range k6JourneyPaths {
		if !found[path] {
			t.Errorf("k6JourneyPaths lists %s but the journey does not "+
				"request it; drop it from the map.", path)
		}
	}
}

// The resource read is the one request worth load-testing, and it has to come
// from the environment. A hardcoded resource name is a name some projects do
// not have, which is the bug this replaced.
func TestK6JourneyTakesItsResourceFromTheEnvironment(t *testing.T) {
	journey := withoutComments(k6CommonJS())

	if !strings.Contains(journey, "__ENV.RESOURCE") {
		t.Error("the k6 journey no longer reads RESOURCE from the environment")
	}
	for _, hardcoded := range []string{"/api/blogs", "/api/posts", "/api/products"} {
		if strings.Contains(journey, hardcoded) {
			t.Errorf("the k6 journey hardcodes %s, which is a resource most "+
				"projects do not have. Use -e RESOURCE=<name>.", hardcoded)
		}
	}
	// A body that merely exists proves nothing: an error envelope has one.
	if strings.Contains(journey, "body.length > 0") {
		t.Error("the k6 journey checks that a body is non-empty, which an " +
			"error response passes")
	}
}

// k6 calls setup() only when the entry script exports it, so a script that
// imports the journey without re-exporting setup gets an undefined token and
// every authenticated request 401s.
func TestEveryK6ScriptReExportsSetup(t *testing.T) {
	scripts := map[string]string{
		"smoke.js":        k6SmokeJS(),
		"average-load.js": k6AverageLoadJS(),
		"stress.js":       k6StressJS(),
		"spike.js":        k6SpikeJS(),
		"soak.js":         k6SoakJS(),
		"breakpoint.js":   k6BreakpointJS(),
	}

	names := make([]string, 0, len(scripts))
	for name := range scripts {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		src := scripts[name]
		if !strings.Contains(src, "export { setup } from './lib/common.js'") {
			t.Errorf("%s does not re-export setup, so k6 never signs in and "+
				"every authenticated request in the journey gets a 401", name)
		}
		if !strings.Contains(src, "export default userJourney") {
			t.Errorf("%s no longer runs the shared journey", name)
		}
	}
}
