package scaffold

import "strings"

// apiSSOTestGo emits internal/handlers/sso_test.go.
//
// These cover the part of SSO that a browser test can't reach: how a login
// resolves to an account, and how IdP groups become roles. They exist because
// the first version of this code had two real bugs that only a test caught —
// a `default:true` bool meant an operator turning OFF just-in-time
// provisioning silently got it ON, and the same for a disabled connection.
func apiSSOTestGo() string {
	src := `package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
)

func ssoTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&models.User{}, &models.SSOConnection{}, &models.UserIdentity{},
		&models.Role{}, &models.UserRole{},
	))
	return db
}

func ssoCtx() *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c
}

func ssoHandler(db *gorm.DB) *SSOHandler { return &SSOHandler{DB: db} }

func ssoRequest(t *testing.T, method, body string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func extUser(sub, email, first string) externalIdentity {
	return externalIdentity{Subject: sub, Email: email, FirstName: first}
}

// JIT provisioning creates the account AND links the external identity.
func TestSSO_JITProvisioning(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true}
	require.NoError(t, db.Create(&conn).Error)

	h := ssoHandler(db)
	u, err := h.resolveUser(ssoCtx(), &conn, extUser("sub-1", "Bob@Acme.com", "Bob"))
	require.NoError(t, err)
	assert.Equal(t, "bob@acme.com", u.Email, "email should be normalized to lowercase")
	assert.True(t, u.Active)
	require.NotNil(t, u.EmailVerifiedAt, "the IdP asserted the address")

	var id models.UserIdentity
	require.NoError(t, db.Where("provider = ? AND subject = ?", "acme", "sub-1").First(&id).Error)
	assert.Equal(t, u.ID, id.UserID)
}

// The second login matches on the IdP subject even after the address changes —
// the user must keep their account, not silently acquire a second one.
func TestSSO_MatchesBySubjectAfterEmailChange(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true}
	require.NoError(t, db.Create(&conn).Error)
	h := ssoHandler(db)

	first, err := h.resolveUser(ssoCtx(), &conn, extUser("sub-1", "bob@acme.com", "Bob"))
	require.NoError(t, err)

	second, err := h.resolveUser(ssoCtx(), &conn, extUser("sub-1", "robert@acme.com", "Robert"))
	require.NoError(t, err)

	assert.Equal(t, first.ID, second.ID, "same subject must resolve to the same account")

	var count int64
	db.Model(&models.User{}).Count(&count)
	assert.EqualValues(t, 1, count, "a renamed address must not create a second user")
}

// An existing password user is adopted by email the first time SSO is used.
func TestSSO_LinksExistingUserByEmail(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: false, Enabled: true}
	require.NoError(t, db.Create(&conn).Error)

	existing := models.User{FirstName: "Bob", LastName: "B", Email: "bob@acme.com", Active: true}
	require.NoError(t, db.Create(&existing).Error)

	h := ssoHandler(db)
	u, err := h.resolveUser(ssoCtx(), &conn, extUser("sub-9", "bob@acme.com", "Bob"))
	require.NoError(t, err)
	assert.Equal(t, existing.ID, u.ID, "should adopt the existing account, not create one")

	var id models.UserIdentity
	require.NoError(t, db.Where("subject = ?", "sub-9").First(&id).Error)
	assert.Equal(t, existing.ID, id.UserID)
}

// With JIT off, a valid authentication for an unknown person is refused.
func TestSSO_JITDisabledRefusesUnknownUser(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: false, Enabled: true}
	require.NoError(t, db.Create(&conn).Error)

	h := ssoHandler(db)
	_, err := h.resolveUser(ssoCtx(), &conn, extUser("sub-x", "stranger@acme.com", "Stranger"))
	require.Error(t, err)

	var count int64
	db.Model(&models.User{}).Count(&count)
	assert.EqualValues(t, 0, count, "no account should be created")
}

// A disabled account cannot sign in through SSO.
func TestSSO_DisabledUserRefused(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true}
	require.NoError(t, db.Create(&conn).Error)
	// Create then disable: User.Active carries gorm:"default:true", so a struct
	// create with Active:false silently stores true. The admin disables via an
	// update, which is what this mirrors.
	gone := models.User{FirstName: "Gone", LastName: "Away", Email: "gone@acme.com"}
	require.NoError(t, db.Create(&gone).Error)
	require.NoError(t, db.Model(&gone).Update("active", false).Error)

	h := ssoHandler(db)
	_, err := h.resolveUser(ssoCtx(), &conn, extUser("sub-g", "gone@acme.com", "Gone"))
	require.Error(t, err)
}

// Group mapping must write BOTH the user_roles join and the legacy users.role
// string — the JWT is minted from the latter, so if they disagree the token and
// the grants disagree.
func TestSSO_GroupMappingWritesBothRoleStores(t *testing.T) {
	db := ssoTestDB(t)
	require.NoError(t, db.Create(&models.Role{Name: "ADMIN"}).Error)
	require.NoError(t, db.Create(&models.Role{Name: "EDITOR"}).Error)

	conn := models.SSOConnection{
		Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true,
		GroupsClaim:   "groups",
		GroupMappings: ~{"it-admins":"ADMIN","engineering":"EDITOR"}~,
	}
	require.NoError(t, db.Create(&conn).Error)

	h := ssoHandler(db)
	ext := extUser("sub-1", "bob@acme.com", "Bob")
	ext.Groups = []string{"it-admins", "unrelated"}

	u, err := h.resolveUser(ssoCtx(), &conn, ext)
	require.NoError(t, err)
	require.NoError(t, h.applyGroupRoles(&conn, u, ext))

	var links []models.UserRole
	require.NoError(t, db.Where("user_id = ?", u.ID).Find(&links).Error)
	assert.Len(t, links, 1, "only the mapped group should grant a role")

	var reloaded models.User
	require.NoError(t, db.Where("id = ?", u.ID).First(&reloaded).Error)
	assert.Equal(t, "ADMIN", reloaded.Role, "legacy role string must track the join table")
}

// Removing someone from a group at the IdP must revoke the role here.
func TestSSO_GroupRemovalRevokesRole(t *testing.T) {
	db := ssoTestDB(t)
	require.NoError(t, db.Create(&models.Role{Name: "ADMIN"}).Error)
	require.NoError(t, db.Create(&models.Role{Name: "EDITOR"}).Error)

	conn := models.SSOConnection{
		Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true,
		GroupsClaim:   "groups",
		GroupMappings: ~{"it-admins":"ADMIN","engineering":"EDITOR"}~,
	}
	require.NoError(t, db.Create(&conn).Error)
	h := ssoHandler(db)

	ext := extUser("sub-1", "bob@acme.com", "Bob")
	ext.Groups = []string{"it-admins"}
	u, err := h.resolveUser(ssoCtx(), &conn, ext)
	require.NoError(t, err)
	require.NoError(t, h.applyGroupRoles(&conn, u, ext))

	// Next login: demoted at the IdP.
	ext.Groups = []string{"engineering"}
	require.NoError(t, h.applyGroupRoles(&conn, u, ext))

	var reloaded models.User
	require.NoError(t, db.Where("id = ?", u.ID).First(&reloaded).Error)
	assert.Equal(t, "EDITOR", reloaded.Role, "losing the admin group must drop admin")

	var links []models.UserRole
	require.NoError(t, db.Where("user_id = ?", u.ID).Find(&links).Error)
	assert.Len(t, links, 1)
}

// A malformed mapping must grant nothing rather than everything.
func TestSSO_MalformedMappingFailsClosed(t *testing.T) {
	db := ssoTestDB(t)
	require.NoError(t, db.Create(&models.Role{Name: "ADMIN"}).Error)

	conn := models.SSOConnection{
		Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true,
		GroupsClaim: "groups", GroupMappings: "{not json",
	}
	require.NoError(t, db.Create(&conn).Error)
	h := ssoHandler(db)

	ext := extUser("sub-1", "bob@acme.com", "Bob")
	ext.Groups = []string{"it-admins"}
	u, err := h.resolveUser(ssoCtx(), &conn, ext)
	require.NoError(t, err)
	require.NoError(t, h.applyGroupRoles(&conn, u, ext))

	var links []models.UserRole
	require.NoError(t, db.Where("user_id = ?", u.ID).Find(&links).Error)
	assert.Len(t, links, 0, "a broken mapping must not grant roles")
}

// A customer's IdP is trusted for the customer's domains and nothing else. It
// used to be able to assert any address, and linking by email then signed it
// in as whoever held that address, the app's own administrator included.
func TestSSO_RefusesAnAddressOutsideTheDomains(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true}
	require.NoError(t, db.Create(&conn).Error)
	admin := models.User{FirstName: "Ada", LastName: "Admin", Email: "ada@yourapp.com", Role: "ADMIN", Active: true}
	require.NoError(t, db.Create(&admin).Error)
	h := ssoHandler(db)

	_, err := h.resolveUser(ssoCtx(), &conn, extUser("attacker", "ada@yourapp.com", "Mallory"))
	require.Error(t, err, "the IdP signed in as an account outside its domains")
	_, err = h.resolveUser(ssoCtx(), &conn, extUser("stranger", "someone@elsewhere.com", "S"))
	require.Error(t, err, "the IdP provisioned an account outside its domains")

	var identities, users int64
	db.Model(&models.UserIdentity{}).Count(&identities)
	db.Model(&models.User{}).Count(&users)
	assert.EqualValues(t, 0, identities, "nothing may be linked")
	assert.EqualValues(t, 1, users, "nothing may be provisioned")
}

// A link made before the domain check existed must not keep working when it
// points outside the connection's domains.
func TestSSO_RefusesAnOldLinkOutsideTheDomains(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true}
	require.NoError(t, db.Create(&conn).Error)
	admin := models.User{FirstName: "Ada", LastName: "Admin", Email: "ada@yourapp.com", Role: "ADMIN", Active: true}
	require.NoError(t, db.Create(&admin).Error)
	require.NoError(t, db.Create(&models.UserIdentity{UserID: admin.ID, Provider: "acme", Subject: "attacker"}).Error)

	_, err := ssoHandler(db).resolveUser(ssoCtx(), &conn, extUser("attacker", "ada@yourapp.com", "Mallory"))
	require.Error(t, err)
}

// Leaving every mapped group drops the role, rather than keeping what it
// granted until someone notices.
func TestSSO_LeavingEveryMappedGroupRevokes(t *testing.T) {
	db := ssoTestDB(t)
	require.NoError(t, db.Create(&models.Role{Name: "EDITOR"}).Error)
	require.NoError(t, db.Create(&models.Role{Name: models.RoleUser}).Error)
	conn := models.SSOConnection{
		Slug: "acme", Name: "Acme", Domains: "acme.com", JITProvisioning: true, Enabled: true,
		GroupsClaim: "groups", GroupMappings: ~{"engineering":"EDITOR"}~,
	}
	require.NoError(t, db.Create(&conn).Error)
	h := ssoHandler(db)

	ext := extUser("sub-1", "bob@acme.com", "Bob")
	ext.Groups = []string{"engineering"}
	u, err := h.resolveUser(ssoCtx(), &conn, ext)
	require.NoError(t, err)
	require.NoError(t, h.applyGroupRoles(&conn, u, ext))

	ext.Groups = nil
	require.NoError(t, h.applyGroupRoles(&conn, u, ext))
	var reloaded models.User
	require.NoError(t, db.Where("id = ?", u.ID).First(&reloaded).Error)
	assert.Equal(t, models.RoleUser, reloaded.Role, "leaving the group must revoke what it granted")
}

// Turning IdP-initiated sign-in off failed on a column GORM had named
// differently, and an update that left fields out cleared them.
func TestSSO_UpdateTogglesIdPInitiatedAndKeepsTheRest(t *testing.T) {
	db := ssoTestDB(t)
	conn := models.SSOConnection{
		Slug: "acme", Name: "Acme", Domains: "acme.com", Protocol: "saml", Enabled: true,
		MetadataURL: "https://idp.acme.com/metadata", GroupMappings: ~{"eng":"EDITOR"}~, AllowIDPInitiated: true,
	}
	require.NoError(t, db.Create(&conn).Error)

	c, w := ssoRequest(t, "PUT", ~{"allow_idp_initiated": false}~)
	c.Params = gin.Params{{Key: "id", Value: conn.ID}}
	ssoHandler(db).Update(c)
	require.Equal(t, 200, w.Code, w.Body.String())

	var got models.SSOConnection
	require.NoError(t, db.First(&got, "id = ?", conn.ID).Error)
	assert.False(t, got.AllowIDPInitiated, "IdP-initiated sign-in is still on")
	assert.Equal(t, "https://idp.acme.com/metadata", got.MetadataURL, "an update that left it out cleared it")
	assert.Equal(t, ~{"eng":"EDITOR"}~, got.GroupMappings)
}

// Discovery sends an address to the first connection claiming its domain, so a
// second claim made sign-in depend on row order.
func TestSSO_ADomainBelongsToOneConnection(t *testing.T) {
	db := ssoTestDB(t)
	require.NoError(t, db.Create(&models.SSOConnection{Slug: "acme", Name: "Acme", Domains: "acme.com", Enabled: true}).Error)

	c, w := ssoRequest(t, "POST",
		~{"slug":"acme2","name":"Acme again","protocol":"saml","domains":"ACME.com","metadata_url":"https://idp.example/m"}~)
	ssoHandler(db).Create(c)
	assert.Equal(t, 400, w.Code, "a second connection claimed acme.com")
}
`
	return strings.ReplaceAll(src, "~", "`")
}
