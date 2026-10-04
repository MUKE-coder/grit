package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
	"library/apps/api/internal/paginate"
)

// ActivityService reads the two audit tables.
//
// Writing to them is what the LogX functions in this package do, and they stay as
// they are: a caller logging an action should not have to build a service to do
// it. Reading them is the other half, and it was in three handlers: the request
// log, the semantic feed and the OCSF export a collector polls.
//
// They are here because each is a decision about data rather than about a
// request: which columns a client may sort or filter by, that the severity panel
// bounds its own window rather than asking the database for an interval, and that
// the export is cursor-paginated on (created_at, id) so a collector resumes
// exactly where it stopped.
type ActivityService struct {
	DB *gorm.DB
}

func (s *ActivityService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// activityLogListConfig is what the request log may be sorted and filtered by.
// Whitelisted, because each name ends up in SQL.
var activityLogListConfig = paginate.Config{
	Sortable: map[string]bool{
		"created_at": true,
		"status":     true,
		"method":     true,
	},
	Filterable:   map[string]bool{"user_id": true, "method": true, "resource": true},
	DefaultSort:  "created_at",
	DefaultOrder: "desc",
}

// Logs returns one page of the request log.
//
// pathPrefix narrows to one area of the API. record answers the question an
// access review asks about one patient's chart: everybody who read or changed
// this row.
func (s *ActivityService) Logs(ctx context.Context, p paginate.Params, pathPrefix, record string) (paginate.Result[models.ActivityLog], error) {
	q := s.db(ctx).Model(&models.ActivityLog{})
	if pathPrefix != "" {
		q = q.Where("path LIKE ?", pathPrefix+"%")
	}
	if record != "" {
		q = q.Where("resource_ids LIKE ?", "%"+record+"%")
	}
	return paginate.List[models.ActivityLog](q, p, activityLogListConfig)
}

// userActivityListConfig is what the semantic feed may be sorted and filtered by.
var userActivityListConfig = paginate.Config{
	Sortable: map[string]bool{
		"created_at": true,
		"severity":   true,
		"action":     true,
	},
	Filterable:   map[string]bool{"user_id": true, "action": true, "severity": true, "resource_type": true},
	DefaultSort:  "created_at",
	DefaultOrder: "desc",
}

// Events returns one page of the semantic activity feed. search matches the
// summary, which is the line a person reads.
func (s *ActivityService) Events(ctx context.Context, p paginate.Params, search string) (paginate.Result[models.UserActivity], error) {
	q := s.db(ctx).Model(&models.UserActivity{})
	if search != "" {
		q = q.Where("summary LIKE ?", "%"+search+"%")
	}
	return paginate.List[models.UserActivity](q, p, userActivityListConfig)
}

// SeverityCounts counts events by severity since a cutoff, for the chips above
// the activity dashboard.
//
// The cutoff is bound as a value rather than written as NOW() - INTERVAL '24
// hours', which is Postgres-only syntax and errors on SQLite: that panel returned
// zeros on every SQLite project until it was changed. The error is returned, too.
// It used to be dropped, so a failure drew the same zeros as a quiet day.
func (s *ActivityService) SeverityCounts(ctx context.Context, since time.Time) (map[string]int64, error) {
	type bucket struct {
		Severity string
		Count    int64
	}
	var rows []bucket
	if err := s.db(ctx).Model(&models.UserActivity{}).
		Select("severity, COUNT(*) AS count").
		Where("created_at > ?", since).
		Group("severity").Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := map[string]int64{"info": 0, "warn": 0, "critical": 0, "total": 0}
	for _, r := range rows {
		out[r.Severity] = r.Count
		out["total"] += r.Count
	}
	return out, nil
}

// Export returns one page of the audit log for a collector, oldest first.
//
// since is a wall-clock floor, which is a collector's first poll. afterID is the
// exact cursor for every poll after that, and wins ties inside the same
// millisecond as since. An afterID that no longer exists falls through to since
// rather than erroring, so a collector that lost its place still makes progress
// instead of wedging.
func (s *ActivityService) Export(ctx context.Context, since time.Time, afterID string, limit int) ([]models.UserActivity, error) {
	q := s.db(ctx).Model(&models.UserActivity{}).Order("created_at asc, id asc").Limit(limit)
	if !since.IsZero() {
		q = q.Where("created_at >= ?", since)
	}
	if afterID != "" {
		var cursor models.UserActivity
		if err := s.db(ctx).Select("created_at", "id").First(&cursor, "id = ?", afterID).Error; err == nil {
			q = q.Where("(created_at, id) > (?, ?)", cursor.CreatedAt, cursor.ID)
		}
	}
	var rows []models.UserActivity
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
