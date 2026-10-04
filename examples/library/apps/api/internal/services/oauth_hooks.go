package services

import (
	"context"
	"log"
	"sync"

	"github.com/markbates/goth"

	"library/apps/api/internal/models"
)

// OAuthProfileHook runs after a social login, with the account the login
// resolved to and the profile the provider returned. Register one from
// routes.Setup:
//
//	services.OnOAuthLogin(func(ctx context.Context, user *models.User, profile goth.User) error {
//	    if profile.Provider != "github" || user.GithubUsername != "" {
//	        return nil
//	    }
//	    // NickName is the provider's handle: the GitHub login, the Twitter @.
//	    return db.WithContext(ctx).Model(user).Update("github_username", profile.NickName).Error
//	})
//
// The user row already exists and has been linked to the provider, so a hook
// may write to it. It runs on every social login, not only the first, which is
// what makes it usable for filling in something that was added to the schema
// later.
type OAuthProfileHook func(ctx context.Context, user *models.User, profile goth.User) error

var (
	oauthHookMu sync.RWMutex
	oauthHooks  []OAuthProfileHook
)

// OnOAuthLogin adds a hook. Call it at boot, from routes.Setup or main.
func OnOAuthLogin(hook OAuthProfileHook) {
	if hook == nil {
		return
	}
	oauthHookMu.Lock()
	defer oauthHookMu.Unlock()
	oauthHooks = append(oauthHooks, hook)
}

// RunOAuthLoginHooks runs them in the order they were registered. The OAuth
// callback calls this; an app does not.
//
// A hook that fails is logged and the rest still run, and the login still
// succeeds. That is a deliberate choice: these hooks decorate an account with
// what the provider knew, and refusing somebody a session because their avatar
// URL would not save is a worse failure than the one it reports. A hook that
// must be able to refuse a login belongs in the callback itself.
func RunOAuthLoginHooks(ctx context.Context, user *models.User, profile goth.User) {
	oauthHookMu.RLock()
	hooks := make([]OAuthProfileHook, len(oauthHooks))
	copy(hooks, oauthHooks)
	oauthHookMu.RUnlock()

	for _, hook := range hooks {
		if err := hook(ctx, user, profile); err != nil {
			log.Printf("oauth: a login hook for %s failed: %v", profile.Provider, err)
		}
	}
}

// ResetOAuthLoginHooks empties the registry. For tests.
func ResetOAuthLoginHooks() {
	oauthHookMu.Lock()
	defer oauthHookMu.Unlock()
	oauthHooks = nil
}
