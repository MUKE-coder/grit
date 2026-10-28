package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

/*
 * The Google scope, in an existing project.
 *
 * goth's google provider falls back to []string{"email"} when New is called
 * with no scopes, and with only that Google returns no name and no picture.
 * Every account created through Google got FirstName "User", a blank last
 * name and no avatar, so anything listing people showed a column of "User".
 * The callback was always written to copy the name; nothing ever asked for
 * it. Reported as grit#120.
 *
 * cmd/server/main.go carries the project's own code, so this inserts the two
 * scopes rather than rewriting the file, and a main.go that no longer looks
 * like Grit's is told what to add instead of being touched.
 */

const (
	// The call as every release before this one wrote it. Matched without
	// leading indentation so it survives gofmt moving the block.
	googleScopeAnchor = `cfg.AppURL+"/api/auth/oauth/google/callback",
		))`

	googleScopeFixed = `cfg.AppURL+"/api/auth/oauth/google/callback",
			// goth defaults to ["email"] when no scope is given, and with only
			// that Google returns no name and no picture: every account signed
			// up through Google was created as "User" with a blank avatar.
			// "profile" is what asks for them.
			"email", "profile",
		))`

	googleScopeAdvice = `apps/api/cmd/server/main.go is not the one Grit wrote, so Google sign-in still asks for only the email scope and every new account is named "User": add "email", "profile" as the last arguments to google.New`
)

// repairGoogleScopeSource is the edit, separated from the file handling so a
// test can drive it with a string.
func repairGoogleScopeSource(src string) (string, []string, []string) {
	// Nothing to do: either the project has no Google provider, or the scopes
	// are already there.
	if !strings.Contains(src, "google.New(") {
		return src, nil, nil
	}
	if strings.Contains(src, `"email", "profile"`) {
		return src, nil, nil
	}
	if strings.Count(src, googleScopeAnchor) != 1 {
		return src, nil, []string{googleScopeAdvice}
	}
	return strings.Replace(src, googleScopeAnchor, googleScopeFixed, 1),
		[]string{`Google sign-in asks for the profile scope, so new accounts get the person's name and photo instead of "User"`},
		nil
}

// repairGoogleScope applies it.
func repairGoogleScope(root string) error {
	path := filepath.Join(root, "apps", "api", "cmd", "server", "main.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairGoogleScopeSource)
}
