package scaffold

import (
	"strings"
	"testing"
)

// A static export has to have a file for every URL it serves.
//
// A dynamic segment cannot be built, because its values are rows in a database
// the build never sees. So the pages that had one read an identifier from the
// query string: /resources/users/view?id= rather than /resources/users/123.
// Same screen, same component, a URL the build can produce.
func TestStaticExportHasNoDynamicSegments(t *testing.T) {
	opts := singleNextOptions()

	if seg := adminDetailSegment(opts); seg != "view" {
		t.Errorf("a resource's detail page is in %q, which the export cannot build", seg)
	}
	if seg := adminDetailSegment(Options{ProjectName: "a", Architecture: ArchTriple, Frontend: FrontendNext}); seg != "[id]" {
		t.Errorf("a monorepo lost its dynamic segment: %q", seg)
	}

	page := adminResourceDetailRoute("users", "users", "Users", opts)
	if !strings.Contains(page, `useSearchParams().get("id")`) {
		t.Error("the detail page does not read the id from the query string")
	}
	// Prerendering a page that reads the query string needs a boundary, and Next
	// fails the whole export without one.
	if !strings.Contains(page, "<Suspense") {
		t.Error("no Suspense boundary, so the export fails on useSearchParams")
	}

	// And the links agree with the routes, or every row opens a 404.
	if !strings.Contains(adminDetailHrefTS(opts), `base + "/view?id=" +`) {
		t.Error("detailHref still builds a path segment, which has no page in the export")
	}
	if !strings.Contains(adminDetailHrefTS(Options{Architecture: ArchTriple}), `base + "/" +`) {
		t.Error("a monorepo's detailHref changed shape")
	}
}

// The binary serving the export must allow what Next needs to start.
//
// Next delivers its bootstrap as an inline <script>. Under the API's strict
// script-src 'self' the browser blocks it, React never hydrates, and the result
// is a page that renders perfectly and does nothing: the sign-in form submits
// natively and puts the password in the URL. Nothing in the console says so,
// because a blocked inline script is a CSP violation report, not an error.
//
// The Vite single is why this went unnoticed: its bundle loads from a src, so
// the strict policy was never in its way.
func TestTheBinarySendsACSPItsOwnFrontendCanRunUnder(t *testing.T) {
	single := apiLoggerMiddlewareGo(Options{ProjectName: "app", Architecture: ArchSingle})
	if !strings.Contains(single, `scriptSrc := "script-src 'self' 'unsafe-inline'; "`) {
		t.Error("a single project's binary serves a frontend under a policy that blocks it")
	}

	// An API with no frontend keeps the strict one: it has no inline script to
	// allow, and allowing one would be a loosening with nothing in return.
	api := apiLoggerMiddlewareGo(Options{ProjectName: "app", Architecture: ArchAPI})
	if !strings.Contains(api, `scriptSrc := "script-src 'self'; "`) {
		t.Error("an API-only project loosened its script-src for a frontend it does not have")
	}
}

// A single project's frontend and API are the same origin, so the address is
// not configuration.
//
// Baking one in is wrong rather than merely unnecessary: the value is fixed when
// the frontend is built and the port is chosen when the binary is run, so a
// project built with the default and started on another port has a frontend
// calling a server that is not there.
func TestASingleCallsItsOwnOrigin(t *testing.T) {
	single := apiCoreTS(Options{ProjectName: "app", Architecture: ArchSingle, Frontend: FrontendNext})
	if !strings.Contains(single, "window.location.origin") {
		t.Error("a single project bakes in an API address instead of using the origin it was served from")
	}

	mono := apiCoreTS(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext})
	if !strings.Contains(mono, `process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"`) {
		t.Error("a monorepo stopped reading its configured API address, which it needs: different origins")
	}
	if strings.Contains(mono, "window.location.origin") {
		t.Error("a monorepo's frontend would call itself rather than the API")
	}
}
