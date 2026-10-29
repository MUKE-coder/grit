package scaffold

// apiSSOHooksGo emits internal/services/sso_hooks.go.
//
// Social sign-in has OnOAuthLogin, run after the account is resolved.
// Enterprise SSO had nothing equivalent: after resolveUser and
// applyGroupRoles the OIDC callback and the SAML ACS went straight to
// GenerateTokenPair, so a project had no supported place to act on an SSO
// sign-in. The only way in was to edit a handler Grit owns and rewrites on
// upgrade.
//
// The obvious thing to do there, in a multi-tenant app, is put the person in
// their company's organisation: an SSO connection belongs to one customer.
// applyGroupRoles only writes platform roles, and an SSOConnection has no
// organisation at all, so every SSO user landed with none, which is what grit
// doctor's tenant-sso-no-organization warning is about. See grit#127.
func apiSSOHooksGo() string {
	return `package services

import (
	"context"
	"log"
	"sync"

	"{{MODULE}}/internal/models"
)

// SSOLoginHook runs after an enterprise SSO sign-in, with the connection the
// person came through, the account it resolved to, and the groups the IdP
// asserted. Register one from routes.Setup:
//
//	services.OnSSOLogin(func(ctx context.Context, conn *models.SSOConnection, user *models.User, groups []string) error {
//	    // The connection belongs to one customer, so this is where a
//	    // multi-tenant app puts the person in their organisation.
//	    orgID, err := orgForConnection(ctx, conn.ID)
//	    if err != nil {
//	        return err
//	    }
//	    return ensureMember(ctx, orgID, user.ID, roleFromGroups(groups))
//	})
//
// The connection and the groups are both here because they are what an app
// needs and what the OAuth hook's goth.User does not carry: which customer
// this is, and what their directory says about the person.
//
// It runs on every SSO sign-in, not only the first, which is what makes it
// usable for a membership that should be kept in step with the directory.
type SSOLoginHook func(ctx context.Context, conn *models.SSOConnection, user *models.User, groups []string) error

var (
	ssoHookMu sync.RWMutex
	ssoHooks  []SSOLoginHook
)

// OnSSOLogin adds a hook. Call it at boot, from routes.Setup or main.
func OnSSOLogin(hook SSOLoginHook) {
	if hook == nil {
		return
	}
	ssoHookMu.Lock()
	defer ssoHookMu.Unlock()
	ssoHooks = append(ssoHooks, hook)
}

// RunSSOLoginHooks runs them in the order they were registered. The OIDC
// callback and the SAML ACS call this; an app does not.
//
// A hook that fails is logged and the rest still run, and the sign-in still
// succeeds, which is the same choice applyGroupRoles makes two lines earlier:
// somebody who authenticated correctly against their company's IdP should not
// be stranded because a membership row would not write. A hook that must be
// able to refuse a sign-in belongs in the callback itself.
//
// It is worth knowing what that means for a multi-tenant app: a failed hook
// leaves a person signed in with no organisation. If that is worse for you
// than a refused sign-in, have the hook record the failure somewhere your
// own middleware checks, rather than letting them through to a half-placed
// session.
func RunSSOLoginHooks(ctx context.Context, conn *models.SSOConnection, user *models.User, groups []string) {
	ssoHookMu.RLock()
	hooks := make([]SSOLoginHook, len(ssoHooks))
	copy(hooks, ssoHooks)
	ssoHookMu.RUnlock()

	for _, hook := range hooks {
		if err := hook(ctx, conn, user, groups); err != nil {
			slug := ""
			if conn != nil {
				slug = conn.Slug
			}
			log.Printf("sso %s: a login hook for %s failed: %v", slug, user.ID, err)
		}
	}
}

// ResetSSOLoginHooks empties the registry. For tests.
func ResetSSOLoginHooks() {
	ssoHookMu.Lock()
	defer ssoHookMu.Unlock()
	ssoHooks = nil
}
`
}

// apiSSOHooksTestGo emits internal/services/sso_hooks_test.go.
func apiSSOHooksTestGo() string {
	return `package services_test

import (
	"context"
	"errors"
	"testing"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// A failing hook must not stop the ones after it, and the connection and the
// asserted groups have to reach every hook: they are what tells an app which
// customer this is and what their directory says.
func TestSSOLoginHooksRunInOrderAndCarryTheConnection(t *testing.T) {
	services.ResetSSOLoginHooks()
	defer services.ResetSSOLoginHooks()

	var order []string
	var gotConn *models.SSOConnection
	var gotGroups []string

	services.OnSSOLogin(func(ctx context.Context, conn *models.SSOConnection, u *models.User, groups []string) error {
		order = append(order, "first")
		gotConn, gotGroups = conn, groups
		return errors.New("deliberate: the next hook must still run")
	})
	services.OnSSOLogin(func(ctx context.Context, conn *models.SSOConnection, u *models.User, groups []string) error {
		order = append(order, "second")
		return nil
	})
	services.OnSSOLogin(nil) // ignored, not a panic

	conn := &models.SSOConnection{ID: "conn-1", Slug: "acme"}
	user := &models.User{ID: "user-1"}
	services.RunSSOLoginHooks(context.Background(), conn, user, []string{"engineering", "admins"})

	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Errorf("hooks ran %v; a failing hook must not stop the rest", order)
	}
	if gotConn == nil || gotConn.Slug != "acme" {
		t.Error("the hook did not get the connection, so it cannot tell which customer this is")
	}
	if len(gotGroups) != 2 || gotGroups[0] != "engineering" {
		t.Errorf("the hook got groups %v, want the asserted ones", gotGroups)
	}
}

// Registering none is the common case and must not panic.
func TestSSOLoginHooksWithNoneRegistered(t *testing.T) {
	services.ResetSSOLoginHooks()
	services.RunSSOLoginHooks(context.Background(), nil, &models.User{ID: "u"}, nil)
}
`
}
