package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/markbates/goth"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/services"
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
