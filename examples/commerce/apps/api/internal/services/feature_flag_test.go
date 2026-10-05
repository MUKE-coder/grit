package services

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/paginate"
)

func flagServiceDB(t *testing.T) (*gorm.DB, *FeatureFlagService) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.FeatureFlag{}, &models.FlagExposure{}))
	return db, &FeatureFlagService{DB: db}
}

// The rollout-health view counts people, not checks. A flag read on every request
// would otherwise report a rollout reaching millions it never left one page for.
func TestExposuresCountsPeopleNotChecks(t *testing.T) {
	db, flags := flagServiceDB(t)
	ctx := context.Background()
	flag := models.FeatureFlag{Name: "new_dashboard", Enabled: true}
	require.NoError(t, flags.Create(ctx, &flag))

	// One person who saw the new dashboard forty times, one who saw it once, and
	// one who got the old one.
	for i := 0; i < 40; i++ {
		require.NoError(t, db.Create(&models.FlagExposure{FlagID: flag.ID, UserID: "u1", Variant: "enabled"}).Error)
	}
	require.NoError(t, db.Create(&models.FlagExposure{FlagID: flag.ID, UserID: "u2", Variant: "enabled"}).Error)
	require.NoError(t, db.Create(&models.FlagExposure{FlagID: flag.ID, UserID: "u3", Variant: "disabled"}).Error)
	// And somebody on another flag entirely.
	other := models.FeatureFlag{Name: "other", Enabled: true}
	require.NoError(t, flags.Create(ctx, &other))
	require.NoError(t, db.Create(&models.FlagExposure{FlagID: other.ID, UserID: "u9", Variant: "enabled"}).Error)

	rows, err := flags.Exposures(ctx, flag.ID)
	require.NoError(t, err)
	require.Len(t, rows, 2)
	// Largest first, which is what the view draws.
	assert.Equal(t, "enabled", rows[0].Variant)
	assert.Equal(t, int64(2), rows[0].Count, "forty checks by one person counted as forty people")
	assert.Equal(t, int64(1), rows[1].Count)
}

// A flag is searched and sorted by a whitelist, because each of those names ends
// up in SQL.
func TestFeatureFlagListIsWhitelisted(t *testing.T) {
	_, flags := flagServiceDB(t)
	ctx := context.Background()
	for _, name := range []string{"beta_banner", "new_dashboard"} {
		require.NoError(t, flags.Create(ctx, &models.FeatureFlag{Name: name, Description: name + " rollout"}))
	}

	page, err := flags.List(ctx, paginate.Params{Page: 1, PageSize: 10, Search: "dashboard"})
	require.NoError(t, err)
	require.Len(t, page.Data, 1)
	assert.Equal(t, "new_dashboard", page.Data[0].Name)

	// A column nobody whitelisted is ignored rather than reaching SQL.
	page, err = flags.List(ctx, paginate.Params{Page: 1, PageSize: 10, SortBy: "rules; drop table feature_flags"})
	require.NoError(t, err)
	assert.Len(t, page.Data, 2)
}

// Create, read, change and remove, which is all the admin screen does.
func TestFeatureFlagRoundTrip(t *testing.T) {
	_, flags := flagServiceDB(t)
	ctx := context.Background()
	flag := models.FeatureFlag{Name: "round_trip"}
	require.NoError(t, flags.Create(ctx, &flag))

	loaded, err := flags.ByID(ctx, flag.ID)
	require.NoError(t, err)
	assert.False(t, loaded.Enabled)

	loaded.Enabled = true
	require.NoError(t, flags.Save(ctx, loaded))
	again, err := flags.ByID(ctx, flag.ID)
	require.NoError(t, err)
	assert.True(t, again.Enabled)

	require.NoError(t, flags.Delete(ctx, again))
	_, err = flags.ByID(ctx, flag.ID)
	assert.Error(t, err)
}
