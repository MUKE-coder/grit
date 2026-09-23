package scaffold

// apiOAuthHooksGo emits internal/services/oauth_hooks.go.
//
// The OAuth callback learns things about a person that only the provider
// knows: their GitHub login, their Google picture, whatever the provider puts
// in the profile. Grit used one field of it and threw the rest away, so an app
// that wanted the GitHub handle had to edit the callback, which Grit owns and
// rewrites on upgrade. Building a shop that invites people to a repository by
// their GitHub username, the username it already had in its hand at sign-in
// was asked for again on a form, and a form is where a typo comes from.
func apiOAuthHooksGo() string {
	return `package services

import (
	"context"
	"log"
	"sync"

	"github.com/markbates/goth"

	"{{MODULE}}/internal/models"
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
`
}

// apiOAuthHooksTestGo emits internal/services/oauth_hooks_test.go.
func apiOAuthHooksTestGo() string {
	return `package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/markbates/goth"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

func TestOAuthLoginHooksRunInOrderWithTheProfile(t *testing.T) {
	services.ResetOAuthLoginHooks()
	t.Cleanup(services.ResetOAuthLoginHooks)

	var order []string
	services.OnOAuthLogin(func(ctx context.Context, user *models.User, profile goth.User) error {
		order = append(order, "first:"+profile.NickName)
		return nil
	})
	services.OnOAuthLogin(func(ctx context.Context, user *models.User, profile goth.User) error {
		order = append(order, "second:"+user.Email)
		return nil
	})

	services.RunOAuthLoginHooks(context.Background(),
		&models.User{Email: "ada@example.com"},
		goth.User{Provider: "github", NickName: "ada-okello"})

	if len(order) != 2 || order[0] != "first:ada-okello" || order[1] != "second:ada@example.com" {
		t.Errorf("hooks ran as %v; they should run in order, with the account and the profile", order)
	}
}

// A hook that fails must not take the login down with it, or a bad avatar URL
// becomes somebody who cannot sign in.
func TestAFailingOAuthHookDoesNotStopTheOthers(t *testing.T) {
	services.ResetOAuthLoginHooks()
	t.Cleanup(services.ResetOAuthLoginHooks)

	ran := false
	services.OnOAuthLogin(func(ctx context.Context, user *models.User, profile goth.User) error {
		return errors.New("the database said no")
	})
	services.OnOAuthLogin(func(ctx context.Context, user *models.User, profile goth.User) error {
		ran = true
		return nil
	})

	services.RunOAuthLoginHooks(context.Background(), &models.User{}, goth.User{Provider: "github"})

	if !ran {
		t.Error("a hook that failed stopped the ones after it")
	}
}

// Nothing registered is the ordinary case, and it must not panic.
func TestNoOAuthHooksIsFine(t *testing.T) {
	services.ResetOAuthLoginHooks()
	services.OnOAuthLogin(nil)
	services.RunOAuthLoginHooks(context.Background(), &models.User{}, goth.User{})
}
`
}
