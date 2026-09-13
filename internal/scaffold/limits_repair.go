package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// routesRequestLimitsBlock is the body cap and deadline table in routes.go. It
// replaced a global middleware.MaxBodySize(10 << 20), which wrapped the body
// before the upload handler could set its own, larger cap.
const routesRequestLimitsBlock = `	// Body size and deadlines, per route. A route not listed gets a 10 MB body
	// and the server's timeouts; a generated resource's /export and /import are
	// transfers without being listed. List a route here when it takes an upload,
	// streams a response or waits on something slow.
	r.Use(middleware.RequestLimits(map[string]middleware.Transfer{
		"POST /api/" + APIVersion + "/uploads":              {Body: 512 << 20, Read: true, Write: true},
		"POST /api/" + APIVersion + "/ai/complete":          {Write: true},
		"POST /api/" + APIVersion + "/ai/chat":              {Write: true},
		"POST /api/" + APIVersion + "/ai/stream":            {Write: true},
		"GET /api/" + APIVersion + "/users/:id/gdpr-export": {Write: true},
		"GET /api/" + APIVersion + "/audit/ocsf":            {Write: true},
		"GET /api/" + APIVersion + "/backups/:id/download":  {Write: true},
	}))
`

// serverTimeoutFields follows Addr and Handler in cmd/server/main.go. The server
// had 15 second read and write timeouts, each one deadline for a whole request
// or response.
const serverTimeoutFields = `
		// Short, because they are what stops a slow client holding a connection
		// open. Uploads, exports and streams get longer deadlines per route, from
		// middleware.RequestLimits in routes.go.
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
`

// repairRequestLimits brings a project up to the fix for H14 in the contact-app
// review: no upload over 10 MB could succeed, anything slower than 15 seconds
// failed, and a streamed export was cut off after it had answered 200.
//
// middleware/limits.go is framework code and arrives whole. routes.go and
// cmd/server/main.go are the developer's, so the changes there anchor on what
// Grit wrote.
func repairRequestLimits(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	for path, content := range map[string]string{
		filepath.Join(apiRoot, "internal", "middleware", "limits.go"):      middlewareLimitsGo(),
		filepath.Join(apiRoot, "internal", "middleware", "limits_test.go"): middlewareLimitsTestGo(),
	} {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if err := repairSourceFile(root, m, routes, repairRequestLimitsRoutesSource); err != nil {
		return err
	}
	server := filepath.Join(apiRoot, "cmd", "server", "main.go")
	if !fileExists(server) {
		return nil
	}
	return repairSourceFile(root, m, server, repairServerTimeoutsSource)
}

const oldBodyCapLine = "\tr.Use(middleware.MaxBodySize(10 << 20)) // 10MB max request body\n"

func repairRequestLimitsRoutesSource(src string) (string, []string, []string) {
	if strings.Contains(src, "middleware.RequestLimits(") || !strings.Contains(src, "middleware.MaxBodySize(") {
		return src, nil, nil
	}
	if strings.Count(src, oldBodyCapLine) != 1 || !strings.Contains(src, "APIVersion") {
		return src, nil, []string{"routes.go caps every body with middleware.MaxBodySize and is not the file Grit wrote: use middleware.RequestLimits, or no upload over the cap can succeed"}
	}
	return strings.Replace(src, oldBodyCapLine, routesRequestLimitsBlock, 1),
		[]string{"uploads, imports, exports and streams get their own body limit and deadlines; other routes keep 10 MB"}, nil
}

var (
	oldServerTimeoutsRe = regexp.MustCompile(`(?m)^\t\tAddr:[ \t]+(.+),\n\t\tHandler:[ \t]+router,\n\t\tReadTimeout:[ \t]+15 \* time\.Second,\n\t\tWriteTimeout:[ \t]+15 \* time\.Second,\n\t\tIdleTimeout:[ \t]+60 \* time\.Second,\n`)
	shortWriteTimeoutRe = regexp.MustCompile(`WriteTimeout:[ \t]+15 \* time\.Second`)
)

func repairServerTimeoutsSource(src string) (string, []string, []string) {
	if strings.Contains(src, "ReadHeaderTimeout") || !strings.Contains(src, "http.Server{") {
		return src, nil, nil
	}
	loc := oldServerTimeoutsRe.FindStringSubmatchIndex(src)
	if loc == nil {
		if shortWriteTimeoutRe.MatchString(src) {
			return src, nil, []string{"cmd/server/main.go is not the file Grit wrote: set ReadHeaderTimeout, and raise the 15 second timeouts, or slow uploads and long responses are cut off"}
		}
		return src, nil, nil
	}
	fields := "\t\tAddr:    " + src[loc[2]:loc[3]] + ",\n\t\tHandler: router,\n" + serverTimeoutFields
	return src[:loc[0]] + fields + src[loc[1]:],
		[]string{"the server gives a request 10 seconds for its headers, 30 to read and 60 to answer, unless its route is a transfer"}, nil
}
