package scaffold

import (
	"strings"
	"testing"
)

// The hosts a seeded placeholder image actually loads from must be allowed.
//
// The seeder writes https://picsum.photos/seed/<name>/600/400 into every
// FileRef it creates. picsum answers that with a 302 to fastly.picsum.photos,
// and both a Content-Security-Policy and next/image's remotePatterns are
// checked against the host a redirect lands on rather than the one the markup
// asked for.
//
// So allowing only picsum.photos allowed nothing. Every scaffolded project with
// seed data showed broken images on its first run, with a console full of
// violations, for images the generator had written itself.
func TestPlaceholderImageHostsIncludeTheRedirectTarget(t *testing.T) {
	const (
		asked  = "picsum.photos"
		landed = "fastly.picsum.photos"
	)

	for name, src := range map[string]string{
		"the Next.js CSP img-src":     nextImageOrigins,
		"the Vite CSP img-src":        viteImageOrigins,
		"next/image's remotePatterns": nextSecurityHeaders(),
	} {
		// Comments stripped first. The comment above each of these explains the
		// redirect and so names the host, which made an earlier version of this
		// test pass against code that did not allow it.
		code := withoutLineComments(src)

		for _, host := range []string{asked, landed} {
			if !strings.Contains(code, `"`+host) && !strings.Contains(code, `'`+host) &&
				!strings.Contains(code, `"https://`+host) && !strings.Contains(code, `'https://`+host) &&
				!strings.Contains(code, `hostname: "`+host+`"`) {
				t.Errorf("%s does not allow %s", name, host)
			}
		}
	}

	// And only in development: a placeholder host has no business in a
	// production policy, which is the decision the original code made and this
	// must not quietly undo.
	for name, src := range map[string]string{
		"the Next.js CSP img-src":     nextImageOrigins,
		"the Vite CSP img-src":        viteImageOrigins,
		"next/image's remotePatterns": nextSecurityHeaders(),
	} {
		if !strings.Contains(withoutLineComments(src), "isDev") {
			t.Errorf("%s allows the placeholder hosts outside development", name)
		}
	}
}

// withoutLineComments drops // lines, so a test asserting on generated code
// cannot be satisfied by prose about that code.
func withoutLineComments(src string) string {
	var kept []string
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}
