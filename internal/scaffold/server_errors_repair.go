package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// A 500 used to send err.Error() to the client in about 35 of Grit's handlers,
// so a failed query answered with the table and column it failed on, and the
// public /api/health endpoint told anonymous callers why Redis or the database
// was down. respond.Internal, meanwhile, discarded the error it was given, so
// those 500s left nothing in the log. respond.ServerError logs the cause with
// the request id and sends the client a message that is safe to read.
//
// hideServerErrors is the rewrite. The templates went through the same rewrite,
// and a test holds them to it; repairServerErrors applies it to the handlers of
// an existing project that upgrade does not rewrite whole.

var (
	serverErrorJSON    = regexp.MustCompile(`c\.JSON\(http\.Status(InternalServerError|BadRequest), gin\.H\{`)
	serverErrorCode    = regexp.MustCompile(`"code":\s*"([A-Z_]+)"`)
	serverErrorMessage = regexp.MustCompile(`"message":\s*(?:"([^"]*)"\s*\+\s*)?err\.Error\(\)`)
	healthErrorLine    = regexp.MustCompile(`(?m)^(\t+)(dbStatus|redisStatus)\.Error = err\.Error\(\)$`)
)

// serverFault400 are the 400s that passed a server's error through, with what
// the client is told instead. Every other 400 describes the request, and is
// left alone.
var serverFault400 = map[string]string{
	"CHART_FAILED":      "Could not compute this chart",
	"STATS_FAILED":      "Could not compute the stats for this resource",
	"SUBMISSION_FAILED": "Your submission could not be saved",
}

// serverErrorHandlers are Grit's handlers that sent err.Error() with a 500.
var serverErrorHandlers = []string{
	"access_review.go", "activity.go", "chart.go", "flags.go", "form_share.go", "gdpr.go",
	"jobs.go", "notification.go", "resource_stats.go", "settings.go", "sync.go",
	"ticket.go", "user_activity.go", "webhooks.go",
}

// hideServerErrors replaces each c.JSON(500, ...) whose message is err.Error()
// with respond.ServerError, keeping the error code. A message with a literal
// prefix ("Failed to retry job: " + err.Error()) keeps the prefix as its text.
func hideServerErrors(src string) (string, int) {
	var b strings.Builder
	last, n := 0, 0
	for _, loc := range serverErrorJSON.FindAllStringSubmatchIndex(src, -1) {
		if loc[0] < last {
			continue
		}
		// An example in a comment is not code: rewriting it would add an import
		// nothing uses.
		if strings.Contains(src[strings.LastIndex(src[:loc[0]], "\n")+1:loc[0]], "//") {
			continue
		}
		end := goStatementEnd(src, loc[0])
		if end < 0 {
			continue
		}
		stmt := src[loc[0]:end]
		code := serverErrorCode.FindStringSubmatch(stmt)
		msg := serverErrorMessage.FindStringSubmatch(stmt)
		if code == nil || msg == nil || strings.Count(stmt, "err.Error()") != 1 || strings.Contains(stmt, "details") {
			continue
		}
		public := "Internal server error"
		if src[loc[2]:loc[3]] == "BadRequest" {
			p, ok := serverFault400[code[1]]
			if !ok {
				continue
			}
			public = p
		} else if prefix := strings.TrimSuffix(strings.TrimSpace(msg[1]), ":"); prefix != "" {
			public = prefix
		}
		b.WriteString(src[last:loc[0]])
		fmt.Fprintf(&b, "respond.ServerError(c, %q, err, %q)", code[1], public)
		last = end
		n++
	}
	b.WriteString(src[last:])
	return b.String(), n
}

// goStatementEnd returns the offset just past the parenthesis that closes the
// call starting at start, skipping string literals, or -1.
func goStatementEnd(src string, start int) int {
	depth, inString := 0, false
	for i := start; i < len(src); i++ {
		switch ch := src[i]; {
		case inString:
			if ch == '\\' {
				i++
			} else if ch == '"' {
				inString = false
			}
		case ch == '"':
			inString = true
		case ch == '(':
			depth++
		case ch == ')':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

// repairServerErrors brings an existing project's handlers and health endpoint
// up to the fix. It runs after writeRespondFiles, which delivers
// respond.ServerError.
func repairServerErrors(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if !fileContains(filepath.Join(apiRoot, "internal", "respond", "respond.go"), "func ServerError(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()
	for _, name := range serverErrorHandlers {
		path := filepath.Join(apiRoot, "internal", "handlers", name)
		if !fileExists(path) {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return repairServerErrorsSource(src, module, name)
		}); err != nil {
			return err
		}
	}
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if fileExists(routes) {
		if err := repairSourceFile(root, m, routes, repairHealthErrorsSource); err != nil {
			return err
		}
	}
	return nil
}

func repairServerErrorsSource(src, module, name string) (string, []string, []string) {
	out, n := hideServerErrors(src)
	if n == 0 {
		return src, nil, nil
	}
	var ok bool
	if out, ok = addImportGroup(out, module+"/internal/respond"); !ok {
		return src, nil, []string{"could not add the respond import to " + name}
	}
	return out, []string{fmt.Sprintf("%s logs %d server errors instead of sending their text to the client", name, n)}, nil
}

func repairHealthErrorsSource(src string) (string, []string, []string) {
	if !healthErrorLine.MatchString(src) {
		return src, nil, nil
	}
	out := healthErrorLine.ReplaceAllStringFunc(src, func(line string) string {
		parts := healthErrorLine.FindStringSubmatch(line)
		label := "database"
		if parts[2] == "redisStatus" {
			label = "redis"
		}
		return parts[1] + `log.Printf("health: ` + label + ` ping failed: %v", err)`
	})
	var ok bool
	if out, ok = addImportGroup(out, "log"); !ok {
		return src, nil, []string{"could not add the log import to routes.go"}
	}
	return out, []string{"/api/health logs why a dependency is down instead of telling anonymous callers"}, nil
}
