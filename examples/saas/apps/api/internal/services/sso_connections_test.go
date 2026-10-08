package services

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"saas/apps/api/internal/models"
)

func ssoServiceDB(t *testing.T) (*gorm.DB, *SSOService) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.SSOConnection{}, &models.UserIdentity{}))
	return db, &SSOService{DB: db}
}

// Two connections claiming one domain made which identity provider a customer
// reaches depend on row order, so the second one to claim it is refused. The
// comparison is on the parsed list: the same domains in a different order, or
// with different spacing, are the same claim.
func TestDomainClaimedElsewhere(t *testing.T) {
	_, sso := ssoServiceDB(t)
	ctx := context.Background()
	acme := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com, acme.co.uk", Enabled: true}
	require.NoError(t, sso.Create(ctx, &acme))

	domain, owner, err := sso.DomainClaimedElsewhere(ctx, "acme.co.uk,example.com", "")
	require.NoError(t, err)
	assert.Equal(t, "acme.co.uk", domain)
	assert.Equal(t, "Acme", owner)

	// Editing that same connection is not a conflict with itself.
	domain, _, err = sso.DomainClaimedElsewhere(ctx, "acme.com", acme.ID)
	require.NoError(t, err)
	assert.Equal(t, "", domain, "a connection conflicted with its own domains")

	domain, _, err = sso.DomainClaimedElsewhere(ctx, "somewhere-else.com", "")
	require.NoError(t, err)
	assert.Equal(t, "", domain)

	// No domains at all is not a claim on everything.
	domain, _, err = sso.DomainClaimedElsewhere(ctx, "   ", "")
	require.NoError(t, err)
	assert.Equal(t, "", domain)
}

// A callback is answered by an enabled connection and nothing else, and a SAML
// assertion is not answered by a connection that has been switched to OIDC.
func TestEnabledBySlugRefusesTheWrongConnection(t *testing.T) {
	_, sso := ssoServiceDB(t)
	ctx := context.Background()
	require.NoError(t, sso.Create(ctx, &models.SSOConnection{Slug: "live", Name: "Live", Protocol: "oidc", Enabled: true}))
	require.NoError(t, sso.Create(ctx, &models.SSOConnection{Slug: "off", Name: "Off", Protocol: "oidc", Enabled: false}))

	conn, err := sso.EnabledBySlug(ctx, "live")
	require.NoError(t, err)
	assert.Equal(t, "Live", conn.Name)

	_, err = sso.EnabledBySlug(ctx, "off")
	assert.Error(t, err, "a disabled connection is still a way in")

	_, err = sso.EnabledBySlugAndProtocol(ctx, "live", "saml")
	assert.Error(t, err, "a SAML assertion was answered by an OIDC connection")
}

// The subject is what a returning sign-in matches on, because it is the only
// identifier an identity provider guarantees stable: somebody who changes their
// address there keeps their account.
func TestIdentityBySubjectAndItsLoginRecord(t *testing.T) {
	db, sso := ssoServiceDB(t)
	ctx := context.Background()
	link := models.UserIdentity{UserID: "u1", Provider: "acme", Subject: "sub-1", Email: "old@acme.com"}
	require.NoError(t, sso.LinkIdentity(ctx, &link))

	found, err := sso.IdentityBySubject(ctx, "acme", "sub-1")
	require.NoError(t, err)
	assert.Equal(t, "u1", found.UserID)

	_, err = sso.IdentityBySubject(ctx, "other", "sub-1")
	assert.Error(t, err, "one provider's subject matched another's")

	at := time.Now().Truncate(time.Second)
	require.NoError(t, sso.RecordIdentityLogin(ctx, found, at, "new@acme.com"))
	var stored models.UserIdentity
	require.NoError(t, db.First(&stored, "id = ?", found.ID).Error)
	assert.Equal(t, "new@acme.com", stored.Email, "the address the provider asserted was not recorded")
	require.NotNil(t, stored.LastLoginAt)
}

// Update returns the row as it now stands, because the admin screen renders what
// comes back and a connection carries fields the request never named.
func TestUpdateReturnsTheWholeConnection(t *testing.T) {
	_, sso := ssoServiceDB(t)
	ctx := context.Background()
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", GroupsClaim: "groups", Enabled: true}
	require.NoError(t, sso.Create(ctx, &conn))

	saved, err := sso.Update(ctx, &conn, map[string]interface{}{"name": "Acme Inc"})
	require.NoError(t, err)
	assert.Equal(t, "Acme Inc", saved.Name)
	assert.Equal(t, "groups", saved.GroupsClaim, "a field the request did not name came back empty")

	all, err := sso.All(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)

	require.NoError(t, sso.Delete(ctx, conn.ID))
	all, err = sso.All(ctx)
	require.NoError(t, err)
	assert.Empty(t, all)
}
