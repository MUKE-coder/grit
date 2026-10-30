package scaffold

import (
	"strings"
	"testing"
)

// A project that installed the multitenant plugin before v3.442.0 has the
// membership guards and no role API. The upgrade adds it, because the
// alternative is `grit plugin remove multitenant && grit plugin add
// multitenant`, which throws away whatever the project changed in those files.

// orgHandlerBefore is the handler as a project upgraded to v3.438.0 has it:
// guarded, with the three-segment key, and no roles.
const orgHandlerBefore = `package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/apps/api/internal/authz"
	"example/apps/api/internal/models"
	"example/apps/api/internal/tenant"
)

type OrganizationHandler struct {
	DB *gorm.DB
}

func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{DB: db}
}

func (h *OrganizationHandler) AddMember(c *gin.Context) {
	orgID := c.Param("id")
	if !h.canManageMembers(c, orgID) {
		return
	}
}

const PermissionMembersManage = "organizations.members.manage"

func (h *OrganizationHandler) canManageMembers(c *gin.Context, orgID string) bool {
	return authz.IsAdmin(c)
}

// isMember verifies the caller belongs to the organization, writing a 404 and
// returning false when they do not.
func (h *OrganizationHandler) isMember(c *gin.Context, orgID string) bool {
	return true
}
`

func TestUpgradeAddsTheOrgRoleAPI(t *testing.T) {
	out, changes, warnings := repairOrgMemberGuardsSource(orgHandlerBefore)
	if len(warnings) > 0 {
		t.Errorf("warnings on a project the repair should handle: %v", warnings)
	}
	if len(changes) == 0 {
		t.Fatal("nothing changed, so an upgraded project has no way to reach the role endpoints")
	}

	// The key that could never be authored.
	if strings.Contains(out, "organizations.members.manage") {
		t.Error("the three-segment key survived. authz.Keys() cannot produce it, " +
			"so no role can be authored holding it (grit#130)")
	}
	if !strings.Contains(out, `PermissionMembersManage = "orgmembers.manage"`) {
		t.Error("the member-management key was not renamed")
	}
	if !strings.Contains(out, `PermissionRolesManage = "orgroles.manage"`) {
		t.Error("the role-management key is missing, so canManageRoles guards on nothing")
	}

	for _, want := range []string{
		"func (h *OrganizationHandler) CreateRole(",
		"func (h *OrganizationHandler) UpdateRole(",
		"func (h *OrganizationHandler) DeleteRole(",
		"func (h *OrganizationHandler) Roles(",
		"func (h *OrganizationHandler) roles() *services.RoleService",
		"example/apps/api/internal/services",
		"var orgAdminKeys = []string{PermissionMembersManage, PermissionRolesManage}",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repair did not add %s", want)
		}
	}

	// Declaring each thing once. A repair that runs twice on the same project
	// is the normal case: people upgrade more than once.
	again, _, _ := repairOrgMemberGuardsSource(out)
	if again != out {
		t.Error("running the repair twice changed the file again, so an upgrade " +
			"duplicates declarations and the project stops compiling")
	}
	if n := strings.Count(out, "func (h *OrganizationHandler) CreateRole("); n != 1 {
		t.Errorf("CreateRole is declared %d times", n)
	}
	if n := strings.Count(out, "var orgAdminKeys"); n != 1 {
		t.Errorf("orgAdminKeys is declared %d times", n)
	}
}

// The handler methods are useless unmounted, and the routes live in another
// file: the repair that adds one without the other is the shape that looks
// like it worked.
func TestUpgradeMountsTheOrgRoleRoutes(t *testing.T) {
	const before = `package routes

func Register(r *gin.Engine) {
	{
		protected.GET("/organizations", orgHandler.MyOrganizations)
		protected.POST("/organizations/:id/members", orgHandler.AddMember)
		protected.DELETE("/organizations/:id/members/:userId", orgHandler.RemoveMember)
	}
}
`
	out, changes, warnings := repairOrgRoleRoutes(before)
	if len(warnings) > 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Fatal("the role endpoints were not mounted")
	}
	for _, want := range []string{
		`protected.GET("/organizations/:id/roles", orgHandler.Roles)`,
		`protected.POST("/organizations/:id/roles", orgHandler.CreateRole)`,
		`protected.PUT("/organizations/:id/roles/:roleId", orgHandler.UpdateRole)`,
		`protected.DELETE("/organizations/:id/roles/:roleId", orgHandler.DeleteRole)`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("not mounted: %s", want)
		}
	}

	again, _, _ := repairOrgRoleRoutes(out)
	if again != out {
		t.Error("the routes were mounted twice, which is a panic at startup: " +
			"gin refuses a duplicate route")
	}
}

// A project with no multitenant plugin must be left alone, not warned at.
func TestOrgRoleRepairIgnoresProjectsWithoutThePlugin(t *testing.T) {
	const plain = "package handlers\n\nfunc nothing() {}\n"
	out, changes, warnings := repairOrgMemberGuardsSource(plain)
	if out != plain || len(changes) > 0 || len(warnings) > 0 {
		t.Errorf("a project without the plugin was touched: changes=%v warnings=%v", changes, warnings)
	}

	const routes = "package routes\n\nfunc Register() {}\n"
	rout, rc, rw := repairOrgRoleRoutes(routes)
	if rout != routes || len(rc) > 0 || len(rw) > 0 {
		t.Errorf("routes without the plugin were touched: changes=%v warnings=%v", rc, rw)
	}
}

// A project from before v3.438.0 has neither the guards nor the roles, so the
// repair runs its whole length. That is the path that writes
// orgMemberGuardsSource into the project, and the path the two tests above
// skip, because a project that already has canManageMembers returns early.
//
// Idempotence matters most here. The rename is a ReplaceAll, so any text the
// repair writes that quotes the OLD key would be rewritten on the next
// upgrade, and every upgrade after that would report a change nobody made.
// That had happened twice by the time this test was written.
func TestUpgradeFromBeforeTheGuardsIsIdempotent(t *testing.T) {
	const before = `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/apps/api/internal/models"
	"example/apps/api/internal/tenant"
)

type OrganizationHandler struct {
	DB *gorm.DB
}

func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{DB: db}
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var admin models.Role
		tx.Where("name = ?", models.RoleAdmin).First(&admin) // grit:role-lookup-ok searched-for, not run
		return tenant.Unscoped(tx).Create(&models.OrganizationMember{
			UserID:         uid,
			OrganizationID: org.ID,
			RoleID:         admin.ID,
		}).Error
	})
	_ = err
}

// Members lists an organization's members.
func (h *OrganizationHandler) Members(c *gin.Context) {}

func (h *OrganizationHandler) AddMember(c *gin.Context) {
	orgID := c.Param("id")
	if !h.isMember(c, orgID) {
		return
	}
	member := models.OrganizationMember{UserID: in.UserID, OrganizationID: orgID, RoleID: in.RoleID}
	_ = member
}

func (h *OrganizationHandler) RemoveMember(c *gin.Context) {
	orgID := c.Param("id")
	if !h.isMember(c, orgID) {
		return
	}
}

// isMember verifies the caller belongs to the organization, writing a 404 and
// returning false when they do not.
func (h *OrganizationHandler) isMember(c *gin.Context, orgID string) bool {
	return true
}
`

	out, changes, warnings := repairOrgMemberGuardsSource(before)
	if len(warnings) > 0 {
		t.Errorf("warnings on a project the repair should handle whole: %v", warnings)
	}
	if len(changes) < 4 {
		t.Errorf("expected the escalation fixes, the rename and the role API; got %v", changes)
	}
	for _, want := range []string{
		"func (h *OrganizationHandler) canManageMembers(",
		"func (h *OrganizationHandler) roleBeyondCaller(",
		"func (h *OrganizationHandler) CreateRole(",
		`PermissionMembersManage = "orgmembers.manage"`,
		"example/apps/api/internal/services",
		"example/apps/api/internal/authz",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repair did not add %s", want)
		}
	}
	if strings.Contains(out, "RoleID:         admin.ID") {
		t.Error("creating an organization still assigns the platform ADMIN role")
	}

	// Nothing the repair writes may quote the old key, or the next upgrade
	// rewrites it.
	if strings.Contains(out, "organizations.members.manage") {
		t.Error("the repaired file still contains the old three-segment key, so " +
			"the next upgrade's rename will rewrite it and report a change that " +
			"nobody made (grit#130)")
	}

	again, moreChanges, _ := repairOrgMemberGuardsSource(out)
	if again != out {
		t.Error("a second upgrade changed the file again")
	}
	if len(moreChanges) > 0 {
		t.Errorf("a second upgrade reported changes with nothing to change: %v", moreChanges)
	}
}
