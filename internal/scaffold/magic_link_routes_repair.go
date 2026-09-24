package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The two routes a sign-in link needs, for a project that already exists.
//
// routes.go is a file every project edits, so an upgrade leaves it alone: the
// handler arrives with nothing pointing at it, and the admin's "Email me a
// sign-in link" answers 404 with no clue why. The same was true of the email
// second factor, which is why this looks like its neighbour.
const (
	magicLinkRouteAnchor = `		auth.POST("/verify-email", authHandler.VerifyEmail)`

	magicLinkRoutes = `
		// Signing in with an emailed link. Consuming is a POST from the page
		// the link opens, never the GET that opens it: mail scanners follow
		// every URL in a message and would spend the token first.
		auth.POST("/magic-link", authHandler.RequestMagicLink)
		auth.POST("/magic-link/consume", authHandler.ConsumeMagicLink)`
)

// And the one route that reads the activity back, which is behind the session
// rather than public.
const (
	magicLinkRecentAnchor = `		protected.GET("/auth/sessions", sessionHandler.List)`

	magicLinkRecentRoute = `
		// Who has been asking for sign-in links to this account. Read-only,
		// and never the token itself.
		protected.GET("/auth/magic-link/recent", authHandler.RecentMagicLinks)`
)

func repairMagicLinkRoutesSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "authHandler") {
		return src, nil, nil
	}

	var changes []string
	var warnings []string
	out := src

	if !strings.Contains(out, "authHandler.RequestMagicLink") {
		if strings.Count(out, magicLinkRouteAnchor) == 1 {
			out = strings.Replace(out, magicLinkRouteAnchor, magicLinkRouteAnchor+magicLinkRoutes, 1)
			changes = append(changes, "signing in with an emailed link (two new routes)")
		} else {
			warnings = append(warnings, "the auth routes in routes.go are not the ones Grit wrote, so sign-in links have no endpoints: add auth.POST(\"/magic-link\", authHandler.RequestMagicLink) and .../magic-link/consume beside the other auth routes")
		}
	}

	if !strings.Contains(out, "authHandler.RecentMagicLinks") {
		if strings.Count(out, magicLinkRecentAnchor) == 1 {
			out = strings.Replace(out, magicLinkRecentAnchor, magicLinkRecentAnchor+magicLinkRecentRoute, 1)
			changes = append(changes, "the account screen can show who asked for a sign-in link")
		} else {
			warnings = append(warnings, "the protected routes in routes.go are not the ones Grit wrote, so the Sign-in links card on the account screen will stay empty: add protected.GET(\"/auth/magic-link/recent\", authHandler.RecentMagicLinks)")
		}
	}

	return out, changes, warnings
}

// repairMagicLinkRoutes applies it to an existing project.
func repairMagicLinkRoutes(root string) error {
	path := filepath.Join(root, "apps", "api", "internal", "routes", "routes.go")
	if !fileExists(path) {
		path = filepath.Join(root, "api", "internal", "routes", "routes.go")
		if !fileExists(path) {
			return nil
		}
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairMagicLinkRoutesSource)
}
