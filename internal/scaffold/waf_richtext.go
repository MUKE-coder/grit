package scaffold

import "fmt"

// WAFRichtextMarker is where wafExcludedRoutes takes additions from v3.335.0.
// grit generate injects before it; grit upgrade adds it to a project that has
// none.
const WAFRichtextMarker = "// grit:waf:richtext"

// wafRichtextMarkerBlock is the marker with the note that explains it.
const wafRichtextMarkerBlock = `		// grit generate resource adds a resource here when it emits a richtext
		// field. Add one by hand if a plain text column carries code or markup
		// that the heuristics read as an attack.
		` + WAFRichtextMarker

// WAFRichtextExclusion is one richtext resource's entry in wafExcludedRoutes'
// paths, comment included.
//
// Both prefixes, because the WAF inspects a body on the way in and it is the
// writes that carry the markup: the public routes are where a site reads from,
// and /admin is where the admin panel writes. Listing only the first is grit#90:
// a post whose body quoted a shell command could be published and then never
// edited, every PUT answered 403 with the WAF logging PathTraversal for the
// "../" inside a code block.
//
// One definition for both callers so that the text is identical whichever wrote
// it, and grit remove can take back out what grit upgrade put in.
func WAFRichtextExclusion(plural string) string {
	return fmt.Sprintf(`		// %s: a richtext body is <p>, <strong> and <img> by definition, and
		// a code block in one reads as an attack to the XSS, traversal and
		// command-injection heuristics. Body inspection only: auth, RBAC,
		// binding validation, sanitize:"html" and the rate limits all still run.
		"/%s", "/%s/*",
		"/admin/%s", "/admin/%s/*",`,
		plural, plural, plural, plural, plural)
}
