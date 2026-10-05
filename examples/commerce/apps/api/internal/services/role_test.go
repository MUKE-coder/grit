package services

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/models"
)

func roleServiceDB(t *testing.T) (*gorm.DB, *RoleService) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}))
	require.NoError(t, models.SeedRoles(db))
	return db, &RoleService{DB: db}
}

func roleNamed(t *testing.T, db *gorm.DB, name string) models.Role {
	t.Helper()
	var role models.Role
	require.NoError(t, db.Where("name = ?", name).First(&role).Error)
	return role
}

// Replacing a user's roles is one change: the assignments and the legacy
// users.role column that name-guarded routes still read.
//
// A stale column there is how a user with the right roles got a 403 anyway, so
// this checks the column and not only the join table.
func TestReplaceUserRolesKeepsTheLegacyColumnInStep(t *testing.T) {
	db, roles := roleServiceDB(t)
	ctx := context.Background()
	admin := roleNamed(t, db, models.RoleAdmin)
	u := models.User{ID: "r1", Email: "r1@x.com", FirstName: "R", LastName: "1", Role: models.RoleUser}
	require.NoError(t, db.Create(&u).Error)

	require.NoError(t, roles.ReplaceUserRoles(ctx, u.ID, []string{admin.ID}))

	var n int64
	db.Model(&models.UserRole{}).Where("user_id = ?", u.ID).Count(&n)
	assert.Equal(t, int64(1), n, "the assignment was not written")
	var after models.User
	require.NoError(t, db.First(&after, "id = ?", u.ID).Error)
	assert.Equal(t, models.RoleAdmin, after.Role, "users.role was left behind, so a name-guarded route still refuses this user")
	g, _ := authz.GrantsFor(db, u.ID)
	assert.True(t, authz.Granted(g, "users.delete"), "the grant cache was not dropped after the write")

	// An empty set clears both halves rather than leaving the old primary role
	// in the column.
	require.NoError(t, roles.ReplaceUserRoles(ctx, u.ID, nil))
	db.Model(&models.UserRole{}).Where("user_id = ?", u.ID).Count(&n)
	assert.Equal(t, int64(0), n)
	require.NoError(t, db.First(&after, "id = ?", u.ID).Error)
	assert.Equal(t, "", after.Role, "users.role still names a role the user no longer holds")
}

// A role reaches a user two ways, and the delete guard has to count both: the
// join table, which is the real model, and the users.role string, which the
// admin's user form still writes. Counting only the first let a role be deleted
// out from under the users assigned to it by name.
//
// It counts users rather than rows, because the number it produces is what the
// operator is told to go and reassign.
func TestAssignmentsCountsBothWaysARoleIsHeld(t *testing.T) {
	db, roles := roleServiceDB(t)
	ctx := context.Background()
	editor := roleNamed(t, db, models.RoleEditor)

	assigned, err := roles.Assignments(ctx, &editor)
	require.NoError(t, err)
	assert.Equal(t, int64(0), assigned)

	// By name only, the way the user form assigns one.
	require.NoError(t, db.Create(&models.User{ID: "r2", Email: "r2@x.com", FirstName: "R", LastName: "2", Role: editor.Name}).Error)
	assigned, err = roles.Assignments(ctx, &editor)
	require.NoError(t, err)
	assert.Equal(t, int64(1), assigned, "a user holding the role by name was not counted")

	// And through the join table. ReplaceUserRoles writes the legacy column too,
	// so this user holds the role both ways and is still one user.
	require.NoError(t, db.Create(&models.User{ID: "r3", Email: "r3@x.com", FirstName: "R", LastName: "3", Role: models.RoleUser}).Error)
	require.NoError(t, roles.ReplaceUserRoles(ctx, "r3", []string{editor.ID}))
	assigned, err = roles.Assignments(ctx, &editor)
	require.NoError(t, err)
	assert.Equal(t, int64(2), assigned, "a user holding the role both ways was counted twice")

	// A deleted account is nobody to reassign.
	require.NoError(t, db.Delete(&models.User{}, "id = ?", "r3").Error)
	assigned, err = roles.Assignments(ctx, &editor)
	require.NoError(t, err)
	assert.Equal(t, int64(1), assigned, "a soft-deleted user still blocks deleting the role")
}

// The list page draws every role with its holder count, from one grouped count
// rather than a query per role.
func TestHolderCounts(t *testing.T) {
	db, roles := roleServiceDB(t)
	ctx := context.Background()
	admin := roleNamed(t, db, models.RoleAdmin)
	for _, id := range []string{"r4", "r5"} {
		require.NoError(t, db.Create(&models.User{ID: id, Email: id + "@x.com", FirstName: "R", LastName: id, Role: models.RoleUser}).Error)
		require.NoError(t, roles.ReplaceUserRoles(ctx, id, []string{admin.ID}))
	}

	counts, err := roles.HolderCounts(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(2), counts[admin.ID])

	all, err := roles.All(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, all)
	assert.True(t, all[0].IsSystem, "built-in roles are listed first, which is the order the admin page shows")
}

// A role deleted outright, because roles.name is unique and a soft delete leaves
// a tombstone that makes that name impossible to use again.
func TestDeleteLeavesNoTombstone(t *testing.T) {
	_, roles := roleServiceDB(t)
	ctx := context.Background()
	mine := models.Role{Name: "AUDITOR"}
	require.NoError(t, roles.Create(ctx, &mine))

	require.NoError(t, roles.Delete(ctx, &mine))

	taken, err := roles.NameTaken(ctx, "AUDITOR")
	require.NoError(t, err)
	assert.False(t, taken, "the name is still taken, so it can never be created again")
	again := models.Role{Name: "AUDITOR"}
	assert.NoError(t, roles.Create(ctx, &again), "a role with the deleted name could not be created")
}
