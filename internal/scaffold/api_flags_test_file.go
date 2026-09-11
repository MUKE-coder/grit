package scaffold

// apiFlagsTestGo emits internal/flags/flags_test.go, so every project checks
// the two things v3.223.0 fixed: a flag can target a business unit, and the
// package-level flags.IsEnabled that the package documented exists.
func apiFlagsTestGo() string {
	return `package flags

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
)

// testEngine starts an engine over one flag in an in-memory database.
func testEngine(t *testing.T, flag models.FeatureFlag) *Engine {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	// Every new connection to ":memory:" is a new, empty database, and the
	// engine reads on its own goroutines: one connection keeps them all on
	// the database the flag was written to.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.FeatureFlag{}, &models.FlagExposure{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&flag).Error; err != nil {
		t.Fatal(err)
	}
	e := New(db, nil)
	t.Cleanup(e.Stop)
	return e
}

// A flag for some business units: on for them, off for the rest, and off for
// a subject the app never gave a business unit.
func TestAttributesTargetABusinessUnit(t *testing.T) {
	flag := models.FeatureFlag{Name: "eu_approvals", Enabled: true}
	if err := flag.SetRules(models.FlagRules{
		RolloutPercentage: 100,
		Attributes:        map[string][]string{"business_unit": {"eu", "KE"}},
	}); err != nil {
		t.Fatal(err)
	}
	e := testEngine(t, flag)

	for _, tc := range []struct {
		unit string
		want bool
	}{
		{"eu", true},
		{"ke", true},
		{"us", false},
		{"", false},
	} {
		s := Subject{UserID: "user-" + tc.unit, Attributes: map[string]string{}}
		if tc.unit != "" {
			s.Attributes["business_unit"] = tc.unit
		}
		if got := e.IsEnabledFor(s, "eu_approvals"); got != tc.want {
			t.Errorf("business_unit %q: got %v, want %v", tc.unit, got, tc.want)
		}
	}
	if e.IsEnabledForUser("user-1", "eu_approvals") {
		t.Error("a check with no attributes passed an attribute rule")
	}
}

// The package documented flags.IsEnabled(c, ...) before it existed.
func TestPackageLevelChecksReachTheEngine(t *testing.T) {
	flag := models.FeatureFlag{Name: "everyone", Enabled: true}
	if err := flag.SetRules(models.FlagRules{RolloutPercentage: 100}); err != nil {
		t.Fatal(err)
	}
	testEngine(t, flag)
	if !IsEnabledForUser("user-1", "everyone") {
		t.Error("the package-level check did not reach the engine")
	}
	if IsEnabledForUser("user-1", "no_such_flag") {
		t.Error("an unknown flag was on")
	}
}
`
}
