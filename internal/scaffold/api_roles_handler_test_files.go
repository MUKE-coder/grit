package scaffold

import "strings"

// roleHandlerTestGo emits internal/handlers/role_test.go.
//
// Shipped with the app because these pin behaviour teams depend on and can
// easily break while editing the catalog:
//   - grants are STORED unexpanded (so wildcards keep inheriting) but SERVED
//     expanded (so the UI never reimplements matching)
//   - unknown permission keys are rejected, instead of stored and silently
//     never matching
//   - built-in roles cannot be renamed or deleted via the API, not merely
//     greyed out in the UI
//   - assigning roles keeps the legacy users.role string in step, or routes
//     still guarded by role name start returning spurious 403s
func roleHandlerTestGo() string {
	src := `package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
)

func roleTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}))
	require.NoError(t, models.SeedRoles(db))
	return db
}

func roleRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewRoleHandler(db)
	r := gin.New()
	r.GET("/permissions", h.Catalog)
	r.GET("/roles", h.List)
	r.POST("/roles", h.Create)
	r.PUT("/roles/:id", h.Update)
	r.DELETE("/roles/:id", h.Delete)
	r.PUT("/users/:id/roles", h.AssignUserRoles)
	return r
}

func do(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer = bytes.NewBuffer(nil)
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRoleAPI_CatalogAndList(t *testing.T) {
	r := roleRouter(roleTestDB(t))

	w := do(t, r, "GET", "/permissions", nil)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "users.")

	w = do(t, r, "GET", "/roles", nil)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	for _, want := range []string{"ADMIN", "EDITOR", "USER", "expanded", "is_system"} {
		assert.Contains(t, body, want)
	}
}

// A wildcard grant must come back expanded, so the UI never reimplements matching.
func TestRoleAPI_ExpandsWildcards(t *testing.T) {
	db := roleTestDB(t)
	r := roleRouter(db)

	w := do(t, r, "POST", "/roles", map[string]any{
		"name": "Uploader", "grants": []string{"uploads.*"},
	})
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var resp struct {
		Data struct {
			Grants   []string ~json:"grants"~
			Expanded []string ~json:"expanded"~
		} ~json:"data"~
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	// Stored unexpanded, so future permissions are inherited...
	assert.Equal(t, []string{"uploads.*"}, resp.Data.Grants)
	// ...but served expanded for rendering.
	assert.Contains(t, resp.Data.Expanded, "uploads.delete")
	assert.NotContains(t, resp.Data.Expanded, "users.delete")
}

// Typos must be rejected, not silently stored and never matched.
func TestRoleAPI_RejectsUnknownPermissions(t *testing.T) {
	r := roleRouter(roleTestDB(t))
	w := do(t, r, "POST", "/roles", map[string]any{
		"name": "Bad", "grants": []string{"nope.create"},
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "nope.create")
}

// System roles are protected SERVER-side, not just in the UI.
func TestRoleAPI_SystemRolesProtected(t *testing.T) {
	db := roleTestDB(t)
	r := roleRouter(db)

	var admin models.Role
	require.NoError(t, db.Where("name = ?", models.RoleAdmin).First(&admin).Error)

	w := do(t, r, "PUT", "/roles/"+admin.ID, map[string]any{
		"name": "SuperBoss", "grants": []string{"*"},
	})
	assert.Equal(t, http.StatusForbidden, w.Code, "renaming a built-in role must be refused")

	w = do(t, r, "DELETE", "/roles/"+admin.ID, nil)
	assert.Equal(t, http.StatusForbidden, w.Code, "deleting a built-in role must be refused")

	// Its permissions ARE editable — only name/delete are locked.
	w = do(t, r, "PUT", "/roles/"+admin.ID, map[string]any{
		"name": admin.Name, "grants": []string{"users.view"},
	})
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func TestRoleAPI_DeleteBlockedWhileAssigned(t *testing.T) {
	db := roleTestDB(t)
	r := roleRouter(db)

	// A custom role assigned to a user via the legacy role STRING (how the admin
	// user form assigns) must not be deletable — deleting it would silently
	// strip that user of its grants.
	w := do(t, r, "POST", "/roles", map[string]any{"name": "support", "grants": []string{"users.view"}})
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var role models.Role
	require.NoError(t, db.Where("name = ?", "support").First(&role).Error)

	u := models.User{ID: "sup1", Email: "sup1@x.com", FirstName: "S", LastName: "1", Role: "support"}
	require.NoError(t, db.Create(&u).Error)

	w = do(t, r, "DELETE", "/roles/"+role.ID, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code, "a role assigned via the role string must not be deletable")

	// Reassign the user, and the role becomes deletable — then its name is free
	// to reuse (hard delete, not a soft-delete tombstone that trips the unique
	// index).
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", "sup1").Update("role", models.RoleUser).Error)
	w = do(t, r, "DELETE", "/roles/"+role.ID, nil)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w = do(t, r, "POST", "/roles", map[string]any{"name": "support", "grants": []string{"users.view"}})
	assert.Equal(t, http.StatusCreated, w.Code, "the name must be reusable after delete")
}

func TestRoleAPI_AssignUserRoles(t *testing.T) {
	db := roleTestDB(t)
	r := roleRouter(db)

	u := models.User{ID: "u1", Email: "u1@x.com", FirstName: "U", LastName: "1", Role: models.RoleUser}
	require.NoError(t, db.Create(&u).Error)

	var editor models.Role
	require.NoError(t, db.Where("name = ?", models.RoleEditor).First(&editor).Error)

	w := do(t, r, "PUT", "/users/u1/roles", map[string]any{"role_ids": []string{editor.ID}})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var n int64
	db.Model(&models.UserRole{}).Where("user_id = ?", "u1").Count(&n)
	assert.Equal(t, int64(1), n)

	// The legacy role string must be kept in step, or routes still guarded by
	// role NAME start returning spurious 403s.
	var got models.User
	db.Where("id = ?", "u1").First(&got)
	assert.Equal(t, models.RoleEditor, got.Role)

	// Unknown role ids are rejected rather than silently dropped.
	w = do(t, r, "PUT", "/users/u1/roles", map[string]any{"role_ids": []string{"does-not-exist"}})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// userRoleSyncTestGo emits internal/handlers/user_role_sync_test.go.
//
// Pins the trap that makes the simple Role dropdown honest: grant resolution
// prefers user_roles, so changing users.role WITHOUT syncing the assignment
// updates the string and changes nothing about what the user can actually do.
func userRoleSyncTestGo() string {
	return `package handlers

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/models"
)

func syncDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}))
	require.NoError(t, models.SeedRoles(db))
	return db
}

// Changing the role dropdown must actually change what the user can do, even
// when they already hold an explicit assignment. Grant resolution prefers
// user_roles, so without the sync this update is a silent no-op.
func TestSyncUserRoleAssignment_ChangesEffectivePermissions(t *testing.T) {
	db := syncDB(t)
	u := models.User{ID: "u1", Email: "u1@x.com", FirstName: "U", LastName: "1", Role: models.RoleAdmin}
	require.NoError(t, db.Create(&u).Error)

	// Start as ADMIN via an explicit assignment.
	require.NoError(t, syncUserRoleAssignment(db, u.ID, models.RoleAdmin))
	g, _ := authz.GrantsFor(db, u.ID)
	require.True(t, authz.Granted(g, "users.delete"), "admin should hold users.delete")

	// Demote to USER — the effective permissions must follow.
	require.NoError(t, syncUserRoleAssignment(db, u.ID, models.RoleUser))
	g, _ = authz.GrantsFor(db, u.ID)
	assert.False(t, authz.Granted(g, "users.delete"), "demotion did not take effect")

	// Exactly one assignment, not an accumulation.
	var n int64
	db.Model(&models.UserRole{}).Where("user_id = ?", u.ID).Count(&n)
	assert.Equal(t, int64(1), n)
}

// An unrecognised legacy role name must not fail the update; resolution falls
// back to the users.role string.
func TestSyncUserRoleAssignment_UnknownNameIsTolerated(t *testing.T) {
	db := syncDB(t)
	u := models.User{ID: "u2", Email: "u2@x.com", FirstName: "U", LastName: "2", Role: "LEGACY_THING"}
	require.NoError(t, db.Create(&u).Error)

	require.NoError(t, syncUserRoleAssignment(db, u.ID, "LEGACY_THING"))

	var n int64
	db.Model(&models.UserRole{}).Where("user_id = ?", u.ID).Count(&n)
	assert.Equal(t, int64(0), n, "unknown role should leave no assignment")
}
`
}
