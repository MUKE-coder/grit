package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The two routes the email second factor needs, for a project that already
// exists.
//
// routes.go is a file every project edits, so an upgrade leaves it alone and
// the new handlers arrive with nothing pointing at them. What that looks like
// from the admin is a button that answers "no route matches
// POST /api/v1/auth/totp/email/send", which is exactly what it did before this
// repair existed.
const (
	emailTwoFactorAnchor = `		protected.POST("/auth/totp/enable", totpHandler.Enable)`

	emailTwoFactorRoutes = `
		// The second factor by email rather than by app. Two steps, like the
		// authenticator: send a code to prove the address works, then confirm
		// it, so nobody turns on a factor they cannot receive.
		protected.POST("/auth/totp/email/send", totpHandler.SendEmailSetupCode)
		protected.POST("/auth/totp/email/enable", totpHandler.EnableEmail)`
)

// The handler also needs something to send with. Without this it is built with
// a nil Config and a nil Mailer, and the first request for a code dereferences
// one of them: a 500 with an empty body and nothing in the log, which is the
// least useful failure a server can produce.
const (
	emailTwoFactorHandlerOld = `	totpHandler := &handlers.TOTPHandler{
		DB:          db,
		AuthService: authService,
		Issuer:      cfg.TOTPIssuer,
	}`
	emailTwoFactorHandlerNew = `	totpHandler := &handlers.TOTPHandler{
		DB:          db,
		AuthService: authService,
		Issuer:      cfg.TOTPIssuer,
		// For the email second factor: the same mailer and queue the auth
		// handler uses, so a code goes out the way every other email does.
		Config: cfg,
		Mailer: svc.Mailer,
		Jobs:   svc.Jobs,
	}`
)

func repairEmailTwoFactorRoutesSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "totpHandler") {
		// A project with no two-factor at all. Nothing to add to.
		return src, nil, nil
	}

	var changes []string
	var warnings []string
	out := src

	if !strings.Contains(out, "totpHandler.SendEmailSetupCode") {
		if strings.Count(out, emailTwoFactorAnchor) == 1 {
			out = strings.Replace(out, emailTwoFactorAnchor, emailTwoFactorAnchor+emailTwoFactorRoutes, 1)
			changes = append(changes, "the second factor can be a code by email (two new routes)")
		} else {
			warnings = append(warnings, "the two-factor routes in routes.go are not the ones Grit wrote, so codes by email have no endpoints: add protected.POST(\"/auth/totp/email/send\", totpHandler.SendEmailSetupCode) and .../email/enable beside them")
		}
	}

	if !strings.Contains(out, "Mailer: svc.Mailer,\n\t\tJobs:   svc.Jobs,\n\t}") && strings.Contains(out, "totpHandler := &handlers.TOTPHandler{") {
		if strings.Count(out, emailTwoFactorHandlerOld) == 1 {
			out = strings.Replace(out, emailTwoFactorHandlerOld, emailTwoFactorHandlerNew, 1)
			changes = append(changes, "the two-factor handler can send email, so asking for a code no longer panics")
		} else if !strings.Contains(out, "Mailer: svc.Mailer") {
			warnings = append(warnings, "the TOTPHandler in routes.go is not the one Grit wrote: give it Config, Mailer and Jobs, or a code by email cannot be sent")
		}
	}

	return out, changes, warnings
}

// repairEmailTwoFactorRoutes applies it to an existing project.
func repairEmailTwoFactorRoutes(root string) error {
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
	return repairSourceFile(root, m, path, repairEmailTwoFactorRoutesSource)
}
