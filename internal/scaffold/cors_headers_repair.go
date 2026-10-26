package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The CORS allow list used to be a string literal written straight into the
// response:
//
//	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, ...")
//
// Which made every new request header a thing somebody had to remember. The
// multitenant plugin introduced X-Organization-ID, told clients to send it, and
// did not add it to that literal, so in a monorepo (web :3000, admin :3001, API
// :8080) the preflight succeeded and the browser then declined to send the real
// request: no status, no log line, an axios network error with nothing behind
// it, and the same call working perfectly under curl (grit#107).
//
// It is a list now, with AllowRequestHeader to add to it from an init(), so the
// header lives in the file that introduces it and leaves when that file is
// deleted. This brings an existing project's cors.go up to that shape.
//
// The headers already in the project's literal are carried across, including
// any added by hand: the repair reshapes the list, it does not replace its
// contents.
var corsAllowHeadersRe = regexp.MustCompile(
	`(?m)^(\s*)c\.Header\("Access-Control-Allow-Headers", "([^"]*)"\)$`)

// repairCORSHeaderList rewrites internal/middleware/cors.go.
func repairCORSHeaderList(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "middleware", "cors.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairCORSHeaderListSource)
}

func repairCORSHeaderListSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func AllowRequestHeader(") {
		return src, nil, nil
	}
	match := corsAllowHeadersRe.FindStringSubmatch(src)
	if match == nil {
		// Either the file has been rewritten by hand or it answers the header
		// some other way. Either is the developer's business, and guessing at
		// a rewrite of a security-relevant middleware is not worth the risk of
		// being wrong. Said out loud, because the consequence is invisible.
		return src, nil, []string{
			"no Access-Control-Allow-Headers line found, so AllowRequestHeader was not added: " +
				"a plugin that needs a request header of its own will not compile against this file"}
	}

	var headers []string
	for _, h := range strings.Split(match[2], ",") {
		if h = strings.TrimSpace(h); h != "" {
			headers = append(headers, h)
		}
	}
	if len(headers) == 0 {
		return src, nil, []string{
			"the Access-Control-Allow-Headers line is empty, so it was left alone"}
	}

	out := corsAllowHeadersRe.ReplaceAllString(src,
		"${1}c.Header(\"Access-Control-Allow-Headers\", allowedRequestHeadersValue)")

	// Before the first CORS constructor, so the list reads above the middleware
	// that writes it.
	anchor := "func CORS("
	if !strings.Contains(out, anchor) {
		anchor = "func CORSDynamic("
	}
	if !strings.Contains(out, anchor) {
		return src, nil, []string{
			"neither CORS nor CORSDynamic is declared here, so the header list was not added"}
	}
	out = strings.Replace(out, anchor, corsHeaderRegistrySource(headers)+anchor, 1)

	withImport, ok := addImportGroup(out, "strings")
	if !ok {
		return src, nil, []string{
			`the header list needs "strings" and this file has no import block to add it to: add the import by hand`}
	}
	return withImport, []string{fmt.Sprintf(
		"the CORS allow list is a list of %d headers with AllowRequestHeader to extend it, "+
			"so a plugin can add the header it introduces", len(headers))}, nil
}

// corsHeaderRegistrySource is the block the repair inserts. Deliberately the
// same text the scaffold writes, so an upgraded project and a new one read
// identically: two shapes of the same file is how the next repair starts
// guessing.
func corsHeaderRegistrySource(headers []string) string {
	var list strings.Builder
	for _, h := range headers {
		fmt.Fprintf(&list, "\t%q,\n", h)
	}
	return `// allowedRequestHeaders is what the browser is told the page may send.
//
// A header missing from this list is not refused with an error: it is stripped
// during preflight, so the request never arrives, no handler runs, nothing is
// logged, and the same call works perfectly under curl. That is the worst shape
// a bug can have, and it is why this is a list anything can add to rather than
// a string literal somebody has to remember to edit.
var allowedRequestHeaders = []string{
` + list.String() + `}

// allowedRequestHeadersValue is the joined header value, built once.
var allowedRequestHeadersValue = strings.Join(allowedRequestHeaders, ", ")

// AllowRequestHeader adds a header to that list.
//
// Call it from an init() in the file that introduces the header, so the feature
// and its CORS entry arrive and leave together instead of the second one being
// a thing to remember. internal/middleware/tenant.go does exactly this for
// X-Organization-ID, which is why removing the multitenant plugin also removes
// the header without touching this file.
//
// init() only. Package-level variables are initialised before any init() runs,
// and every init() finishes before main, so by the time a request can read
// allowedRequestHeadersValue the list is final and nothing synchronises it.
// Calling this from a handler would be a data race.
func AllowRequestHeader(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	for _, existing := range allowedRequestHeaders {
		if strings.EqualFold(existing, name) {
			return
		}
	}
	allowedRequestHeaders = append(allowedRequestHeaders, name)
	allowedRequestHeadersValue = strings.Join(allowedRequestHeaders, ", ")
}

`
}
