package scaffold

import "strings"

// apiBoolFlagTestGo emits internal/models/bool_flags_test.go.
//
// This guards one specific, recurring GORM trap: a bool column declared with
// gorm:"default:true" can never be stored as false through a struct Create.
// GORM omits zero-valued fields from the INSERT when the column carries a
// default, so the database default wins and an operator's explicit "off"
// silently becomes "on".
//
// It has bitten this codebase five times — SSO's Enabled and JITProvisioning,
// User.Active, FormShare.Enabled and BackupSchedule.Enabled — each time
// inverting a security-relevant switch (provision accounts automatically,
// deactivate a user, publish a public form link). None of them fail a build,
// and none are visible without asserting on a round-trip, which is why this
// test exists rather than a code-review rule.
func apiBoolFlagTestGo() string {
	src := `package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func boolFlagDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}, &FormShare{}, &BackupSchedule{}))
	return db
}

// A user created as inactive must stay inactive.
func TestUserActiveFalseSurvivesCreate(t *testing.T) {
	db := boolFlagDB(t)

	u := User{FirstName: "In", LastName: "Active", Email: "inactive@example.com", Active: false}
	require.NoError(t, db.Create(&u).Error)

	var got User
	require.NoError(t, db.Where("email = ?", u.Email).First(&got).Error)
	require.False(t, got.Active,
		"Active:false was not persisted — check for a gorm default on the column")
}

// A share link created disabled must not be live. This one is security
// relevant: the token grants public form submission.
func TestFormShareEnabledFalseSurvivesCreate(t *testing.T) {
	db := boolFlagDB(t)

	fs := FormShare{Token: "tok-disabled", ResourceName: "invoices", Enabled: false}
	require.NoError(t, db.Create(&fs).Error)

	var got FormShare
	require.NoError(t, db.Where("token = ?", fs.Token).First(&got).Error)
	require.False(t, got.Enabled,
		"Enabled:false was not persisted — a disabled share link would be publicly live")
}

// A backup schedule created disabled must not run backups.
func TestBackupScheduleEnabledFalseSurvivesCreate(t *testing.T) {
	db := boolFlagDB(t)

	sc := BackupSchedule{ID: 1, Frequency: "weekly", Time: "02:00", Enabled: false}
	require.NoError(t, db.Create(&sc).Error)

	var got BackupSchedule
	require.NoError(t, db.First(&got, 1).Error)
	require.False(t, got.Enabled,
		"Enabled:false was not persisted — automatic backups would keep running")
}

// True still round-trips, so removing the column defaults didn't invert them.
func TestBoolFlagsTrueStillPersists(t *testing.T) {
	db := boolFlagDB(t)

	u := User{FirstName: "A", LastName: "B", Email: "active@example.com", Active: true}
	require.NoError(t, db.Create(&u).Error)

	var got User
	require.NoError(t, db.Where("email = ?", u.Email).First(&got).Error)
	require.True(t, got.Active)
}
`
	return strings.ReplaceAll(src, "~", "`")
}
