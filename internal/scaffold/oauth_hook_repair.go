package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The call site for services.OnOAuthLogin, in an existing project.
//
// handlers/auth_oauth.go is the project's file: people change what a new
// account is given, where it is redirected, which providers are allowed. So
// the one line that hands the provider's profile to the app's hooks is
// inserted rather than the file rewritten, and a file that no longer looks
// like Grit's is told what to add instead.
const (
	oauthHookAnchor = `	// Generate JWT tokens
	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)`

	oauthHookCall = `	// What the provider knew, handed to whatever this app wants to do with it:
	// the GitHub login, the avatar, the locale. See services.OnOAuthLogin.
	services.RunOAuthLoginHooks(c.Request.Context(), &user, gothUser)

`
)

func repairOAuthHookSource(src string) (string, []string, []string) {
	if strings.Contains(src, "services.RunOAuthLoginHooks(") {
		return src, nil, nil
	}
	if strings.Count(src, oauthHookAnchor) != 1 || !strings.Contains(src, "gothic.CompleteUserAuth") {
		return src, nil, []string{"internal/handlers/auth_oauth.go is not the one Grit wrote, so the social-login hooks are never run: call services.RunOAuthLoginHooks(c.Request.Context(), &user, gothUser) before the tokens are generated"}
	}
	return strings.Replace(src, oauthHookAnchor, oauthHookCall+oauthHookAnchor, 1),
		[]string{"a social login hands the provider's profile to services.OnOAuthLogin hooks"}, nil
}

// repairOAuthHook applies it.
func repairOAuthHook(root string) error {
	path := filepath.Join(root, "apps", "api", "internal", "handlers", "auth_oauth.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairOAuthHookSource)
}
