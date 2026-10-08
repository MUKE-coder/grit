package services

import (
	"context"

	"gorm.io/gorm"

	"saas/apps/api/internal/models"
	"saas/apps/api/internal/paginate"
)

// FeatureFlagService owns the flag rows and the exposure counts drawn from them.
//
// The engine in internal/flags decides what a flag answers for a given user and
// caches it; this is the table behind it. The handler validates the request and
// tells the engine to refresh, which is the one thing that cannot live here: a
// service does not know when a cache that other replicas hold is stale.
type FeatureFlagService struct {
	DB *gorm.DB
}

func (s *FeatureFlagService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// featureFlagListConfig is what the admin's flags page may search, sort and
// filter by. Whitelisted, because each name ends up in SQL.
var featureFlagListConfig = paginate.Config{
	Searchable:   []string{"name", "description"},
	Sortable:     map[string]bool{"name": true, "created_at": true, "enabled": true},
	Filterable:   map[string]bool{"enabled": true},
	DefaultSort:  "name",
	DefaultOrder: "asc",
}

// List returns one page of flags.
func (s *FeatureFlagService) List(ctx context.Context, p paginate.Params) (paginate.Result[models.FeatureFlag], error) {
	return paginate.List[models.FeatureFlag](s.db(ctx).Model(&models.FeatureFlag{}), p, featureFlagListConfig)
}

// ByID reads one flag.
func (s *FeatureFlagService) ByID(ctx context.Context, id string) (*models.FeatureFlag, error) {
	var flag models.FeatureFlag
	if err := s.db(ctx).First(&flag, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &flag, nil
}

// Create saves a new flag. The unique index on the name decides a collision.
func (s *FeatureFlagService) Create(ctx context.Context, flag *models.FeatureFlag) error {
	return s.db(ctx).Create(flag).Error
}

// Save writes a flag the caller has changed.
func (s *FeatureFlagService) Save(ctx context.Context, flag *models.FeatureFlag) error {
	return s.db(ctx).Save(flag).Error
}

// Delete removes a flag.
func (s *FeatureFlagService) Delete(ctx context.Context, flag *models.FeatureFlag) error {
	return s.db(ctx).Delete(flag).Error
}

// FlagExposure is one variant of a flag and how many people saw it.
type FlagExposure struct {
	Variant string `json:"variant"`
	Count   int64  `json:"count"`
}

// Exposures counts the people who saw each variant of one flag, which is what
// the rollout-health view draws.
//
// Distinct users, not rows: a flag checked on every request would otherwise
// report a rollout reaching millions of people it never left the first page for.
func (s *FeatureFlagService) Exposures(ctx context.Context, flagID string) ([]FlagExposure, error) {
	var rows []FlagExposure
	if err := s.db(ctx).Model(&models.FlagExposure{}).
		Select("variant, COUNT(DISTINCT user_id) as count").
		Where("flag_id = ?", flagID).
		Group("variant").
		Order("count desc").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
