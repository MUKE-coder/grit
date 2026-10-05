package authz_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/models"
)

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}); err != nil {
		t.Fatal(err)
	}
	if err := models.SeedRoles(db); err != nil {
		t.Fatal(err)
	}
	return db
}

// Roles must exist automatically, with sane grants.
func TestSeedRoles(t *testing.T) {
	db := newDB(t)

	var roles []models.Role
	db.Find(&roles)
	if len(roles) != 3 {
		t.Fatalf("seeded %d roles, want 3", len(roles))
	}

	var admin models.Role
	db.Where("name = ?", models.RoleAdmin).First(&admin)
	if got := admin.GrantsList(); len(got) != 1 || got[0] != "*" {
		t.Errorf("ADMIN grants = %v, want [*]", got)
	}
	if !admin.IsSystem {
		t.Error("ADMIN must be a system role")
	}

	// Re-seeding must not duplicate or clobber.
	admin.SetGrants([]string{"users.view"})
	db.Save(&admin)
	if err := models.SeedRoles(db); err != nil {
		t.Fatal(err)
	}
	db.Where("name = ?", models.RoleAdmin).First(&admin)
	if got := admin.GrantsList(); len(got) != 1 || got[0] != "users.view" {
		t.Errorf("re-seeding clobbered an operator edit: %v", got)
	}
}

// The legacy users.role string must still authorise, so apps upgrading from
// role-only auth don't lose access before anyone is assigned a Role row.
func TestGrantsFor_LegacyRoleFallback(t *testing.T) {
	db := newDB(t)
	admin := models.User{ID: "u-admin", Email: "a@x.com", FirstName: "A", LastName: "B", Role: models.RoleAdmin}
	plain := models.User{ID: "u-plain", Email: "p@x.com", FirstName: "P", LastName: "B", Role: models.RoleUser}
	db.Create(&admin)
	db.Create(&plain)

	g, err := authz.GrantsFor(db, admin.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !authz.Granted(g, "users.delete") {
		t.Errorf("legacy ADMIN lost access; grants=%v", g)
	}

	g, err = authz.GrantsFor(db, plain.ID)
	if err != nil {
		t.Fatal(err)
	}
	if authz.Granted(g, "users.delete") {
		t.Errorf("plain USER must not hold users.delete; grants=%v", g)
	}
}

// An explicit user_roles assignment takes precedence over the legacy string.
func TestGrantsFor_ExplicitAssignmentWins(t *testing.T) {
	db := newDB(t)
	u := models.User{ID: "u1", Email: "u1@x.com", FirstName: "U", LastName: "1", Role: models.RoleUser}
	db.Create(&u)

	custom := models.Role{Name: "Support"}
	custom.SetGrants([]string{"uploads.*"})
	db.Create(&custom)
	db.Create(&models.UserRole{UserID: u.ID, RoleID: custom.ID})
	authz.Invalidate()

	g, _ := authz.GrantsFor(db, u.ID)
	if !authz.Granted(g, "uploads.delete") {
		t.Errorf("uploads.* should cover uploads.delete; grants=%v", g)
	}
	if authz.Granted(g, "users.delete") {
		t.Errorf("assignment must not leak other permissions; grants=%v", g)
	}
}

// Revoking must take effect immediately, not after a cache TTL.
func TestInvalidate(t *testing.T) {
	db := newDB(t)
	u := models.User{ID: "u2", Email: "u2@x.com", FirstName: "U", LastName: "2", Role: models.RoleAdmin}
	db.Create(&u)

	if g, _ := authz.GrantsFor(db, u.ID); !authz.Granted(g, "users.delete") {
		t.Fatal("expected initial admin access")
	}

	db.Model(&models.User{}).Where("id = ?", u.ID).Update("role", models.RoleUser)
	authz.Invalidate()

	if g, _ := authz.GrantsFor(db, u.ID); authz.Granted(g, "users.delete") {
		t.Error("demotion did not take effect after Invalidate")
	}
}
