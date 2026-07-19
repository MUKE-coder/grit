package plugin

import "strings"

// mtMiddleware emits the middleware that resolves the active organization.
func mtMiddleware(ctx Context) string {
	src := `package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/tenant"
)

// Tenant resolves the caller's active organization and puts it on the request
// context, where the GORM scoping callbacks pick it up.
//
// Resolution order:
//  1. the X-Organization-ID header (explicit; what SPA clients send)
//  2. the user's only organization, if they belong to exactly one
//
// No subdomain resolution: it forces DNS and TLS decisions on every deployment
// and makes local development awkward.
//
// Membership is ALWAYS verified against the database. Trusting the header alone
// would let any authenticated user read another tenant's data by editing one
// request header — the whole isolation model rests on this check.
func Tenant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.Next() // unauthenticated; nothing to scope
			return
		}
		uid, _ := userID.(string)

		requested := c.GetHeader("X-Organization-ID")

		var memberships []models.OrganizationMember
		if err := db.Where("user_id = ?", uid).Find(&memberships).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not resolve organization"},
			})
			c.Abort()
			return
		}

		active := ""
		switch {
		case requested != "":
			for _, m := range memberships {
				if m.OrganizationID == requested {
					active = requested
					break
				}
			}
			if active == "" {
				// Asked for an organization they do not belong to. 404 rather
				// than 403: a 403 confirms it exists and lets an attacker
				// enumerate tenant ids.
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{"code": "NOT_FOUND", "message": "Organization not found"},
				})
				c.Abort()
				return
			}
		case len(memberships) == 1:
			active = memberships[0].OrganizationID
		}

		if active != "" {
			c.Set("org_id", active)
			c.Set("org_role_id", roleFor(memberships, active))
			c.Request = c.Request.WithContext(tenant.WithOrg(c.Request.Context(), active))
		}
		c.Next()
	}
}

func roleFor(memberships []models.OrganizationMember, orgID string) string {
	for _, m := range memberships {
		if m.OrganizationID == orgID {
			return m.RoleID
		}
	}
	return ""
}

// RequireOrganization rejects a request with no active organization. Put it on
// routes that operate on tenant data.
func RequireOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get("org_id"); !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "NO_ORGANIZATION",
					"message": "Select an organization first (send the X-Organization-ID header)",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// mtHandler emits the organization + membership API.
func mtHandler(ctx Context) string {
	src := `package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/tenant"
)

// OrganizationHandler serves organization + membership endpoints.
type OrganizationHandler struct {
	DB *gorm.DB
}

func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{DB: db}
}

// MyOrganizations lists the organizations the caller belongs to.
//
// Reads memberships first and looks organizations up by id, rather than listing
// organizations and filtering — so a bug here can only ever return too FEW
// rows, never another tenant's.
func (h *OrganizationHandler) MyOrganizations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)

	var memberships []models.OrganizationMember
	if err := h.DB.Where("user_id = ?", uid).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not load organizations"}})
		return
	}
	if len(memberships) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Organization{}})
		return
	}

	ids := make([]string, 0, len(memberships))
	for _, m := range memberships {
		ids = append(ids, m.OrganizationID)
	}

	// Organizations are not themselves tenant-scoped rows, so this lookup is
	// exempt from the scoping callbacks.
	var orgs []models.Organization
	if err := tenant.Unscoped(h.DB).Where("id IN ?", ids).Find(&orgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not load organizations"}})
		return
	}

	active, _ := c.Get("org_id")
	c.JSON(http.StatusOK, gin.H{"data": orgs, "active": active})
}

type createOrgInput struct {
	Name        string ~json:"name" binding:"required"~
	Description string ~json:"description"~
}

// Create makes an organization and adds the caller as its owner.
func (h *OrganizationHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)

	var in createOrgInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "Name is required"}})
		return
	}

	org := models.Organization{
		Name:        strings.TrimSpace(in.Name),
		Slug:        slugify(in.Name),
		Description: in.Description,
		OwnerID:     uid,
		Active:      true,
	}

	// Both writes in one transaction: an organization with no members is
	// unreachable, because every read goes through membership.
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tenant.Unscoped(tx).Create(&org).Error; err != nil {
			return err
		}
		var admin models.Role
		tx.Where("name = ?", models.RoleAdmin).First(&admin)
		return tenant.Unscoped(tx).Create(&models.OrganizationMember{
			UserID:         uid,
			OrganizationID: org.ID,
			RoleID:         admin.ID,
		}).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not create organization"}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": org, "message": "Organization created"})
}

// Members lists an organization's members. Callers must belong to it.
func (h *OrganizationHandler) Members(c *gin.Context) {
	if !h.isMember(c, c.Param("id")) {
		return
	}
	var members []models.OrganizationMember
	if err := tenant.Unscoped(h.DB).Where("organization_id = ?", c.Param("id")).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not load members"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members})
}

type addMemberInput struct {
	UserID string ~json:"user_id" binding:"required"~
	RoleID string ~json:"role_id"~
}

// AddMember adds a user to an organization with a role.
func (h *OrganizationHandler) AddMember(c *gin.Context) {
	orgID := c.Param("id")
	if !h.isMember(c, orgID) {
		return
	}

	var in addMemberInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "user_id is required"}})
		return
	}

	member := models.OrganizationMember{UserID: in.UserID, OrganizationID: orgID, RoleID: in.RoleID}
	// Upsert: re-adding an existing member updates their role rather than
	// failing on the composite primary key.
	if err := tenant.Unscoped(h.DB).Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not add member"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": member, "message": "Member added"})
}

// RemoveMember removes a user from an organization.
func (h *OrganizationHandler) RemoveMember(c *gin.Context) {
	orgID := c.Param("id")
	if !h.isMember(c, orgID) {
		return
	}

	// The owner keeps access, or an organization can be orphaned with nobody
	// able to administer it.
	var org models.Organization
	if err := tenant.Unscoped(h.DB).Where("id = ?", orgID).First(&org).Error; err == nil {
		if org.OwnerID == c.Param("userId") {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "The owner cannot be removed"}})
			return
		}
	}

	if err := tenant.Unscoped(h.DB).
		Where("organization_id = ? AND user_id = ?", orgID, c.Param("userId")).
		Delete(&models.OrganizationMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not remove member"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// isMember verifies the caller belongs to the organization, writing a 404 and
// returning false when they do not.
//
// 404 rather than 403: a 403 confirms the organization exists, which lets an
// attacker enumerate tenants.
func (h *OrganizationHandler) isMember(c *gin.Context, orgID string) bool {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)

	var count int64
	tenant.Unscoped(h.DB).Model(&models.OrganizationMember{}).
		Where("user_id = ? AND organization_id = ?", uid, orgID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Organization not found"}})
		return false
	}
	return true
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// mtTenantTests emits isolation tests into the generated project.
//
// These matter more than most generated tests: they assert that a query for one
// organization cannot see another's rows, and that an unscoped context fails
// CLOSED. A regression here is a data breach, not a bug.
func mtTenantTests(ctx Context) string {
	src := `package tenant_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/tenant"
)

// Widget is a tenant-owned model for these tests.
type Widget struct {
	ID   string ~gorm:"primarykey;size:36"~
	Name string
	tenant.Owned
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Widget{}); err != nil {
		t.Fatal(err)
	}
	if err := tenant.RegisterScoping(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func seed(t *testing.T, db *gorm.DB, org, id, name string) {
	t.Helper()
	ctx := tenant.WithOrg(context.Background(), org)
	if err := db.WithContext(ctx).Create(&Widget{ID: id, Name: name}).Error; err != nil {
		t.Fatal(err)
	}
}

// The core guarantee: one organization cannot read another's rows.
func TestScoping_IsolatesOrganizations(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")
	seed(t, db, "org-b", "w2", "beta")

	var got []Widget
	ctx := tenant.WithOrg(context.Background(), "org-a")
	if err := db.WithContext(ctx).Find(&got).Error; err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("org-a saw %d widgets, want 1 — CROSS-TENANT LEAK", len(got))
	}
	if got[0].Name != "alpha" {
		t.Errorf("org-a saw %q", got[0].Name)
	}
}

// Create stamps the active organization, so a row cannot land in the wrong one.
func TestScoping_StampsOrgOnCreate(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")

	var w Widget
	if err := tenant.Unscoped(db).Where("id = ?", "w1").First(&w).Error; err != nil {
		t.Fatal(err)
	}
	if w.OrgID != "org-a" {
		t.Errorf("OrgID = %q, want org-a", w.OrgID)
	}
}

// No active organization must FAIL, not silently return everything. This is the
// property that makes the whole design safe.
func TestScoping_FailsClosedWithoutOrg(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")

	var got []Widget
	err := db.WithContext(context.Background()).Find(&got).Error
	if err == nil {
		t.Fatalf("query without an organization returned %d rows; it must error", len(got))
	}
}

// Updates and deletes are scoped too — otherwise one tenant could modify
// another's data even without reading it.
func TestScoping_UpdateAndDeleteAreScoped(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")
	seed(t, db, "org-b", "w2", "beta")

	ctxA := tenant.WithOrg(context.Background(), "org-a")
	if err := db.WithContext(ctxA).Model(&Widget{}).Where("1 = 1").Update("name", "hacked").Error; err != nil {
		t.Fatal(err)
	}
	var b Widget
	tenant.Unscoped(db).Where("id = ?", "w2").First(&b)
	if b.Name != "beta" {
		t.Errorf("org-a modified org-b's row: %q", b.Name)
	}

	if err := db.WithContext(ctxA).Where("1 = 1").Delete(&Widget{}).Error; err != nil {
		t.Fatal(err)
	}
	var count int64
	tenant.Unscoped(db).Model(&Widget{}).Where("id = ?", "w2").Count(&count)
	if count != 1 {
		t.Error("org-a deleted org-b's row")
	}
}

// Unscoped is the deliberate escape hatch for cross-tenant work.
func TestUnscoped_SeesEverything(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")
	seed(t, db, "org-b", "w2", "beta")

	var all []Widget
	if err := tenant.Unscoped(db).Find(&all).Error; err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Errorf("Unscoped saw %d widgets, want 2", len(all))
	}
}

// Models without an OrgID are untouched — users, roles and settings are global.
func TestScoping_IgnoresNonTenantModels(t *testing.T) {
	type Global struct {
		ID   string ~gorm:"primarykey;size:36"~
		Name string
	}
	db := newDB(t)
	if err := db.AutoMigrate(&Global{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&Global{ID: "g1", Name: "x"}).Error; err != nil {
		t.Fatalf("a non-tenant model must not require an organization: %v", err)
	}
	var got []Global
	if err := db.Find(&got).Error; err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("non-tenant query returned %d rows, want 1", len(got))
	}
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}
