package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/sync"
)

/*
 * What a soft delete leaves behind, and what happens to it.
 *
 * Every generated model carries gorm.DeletedAt, so Delete hides the row rather
 * than removing it. That has always been true and has never been visible: the
 * row vanished from every list, every export and every count, and the only way
 * to see it again was a SQL console. "Deleted" and "gone" meant the same thing
 * to anybody using the app, which makes the soft delete a database detail
 * rather than a feature.
 *
 * This is the other half: the rows are listed, restorable for a retention
 * window, and removed for good when it runs out.
 *
 * The resources come from the sync registry rather than a list maintained here.
 * It is already keyed by the plural table name, already populated by the
 * generator for every resource, and already the answer to "which models does
 * this app own". A second registry would be a second thing to forget to add to.
 */

// TrashRetentionDays is how long a deleted row can be brought back.
//
// Thirty days is the window the rest of the industry has taught people to
// expect, and it is long enough that somebody returning from leave can still
// undo a Friday afternoon. Change it here and the sweep below follows.
const TrashRetentionDays = 30

// TrashService lists, restores and removes soft-deleted rows.
type TrashService struct {
	DB       *gorm.DB
	Registry *sync.Registry
}

func (s *TrashService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// TrashItem is one deleted row, flattened to what a list can show without
// knowing the model.
type TrashItem struct {
	Table     string         `json:"table"`
	ID        string         `json:"id"`
	Label     string         `json:"label"`
	DeletedAt time.Time      `json:"deleted_at"`
	ExpiresAt time.Time      `json:"expires_at"`
	Fields    map[string]any `json:"fields,omitempty"`
}

// TrashBucket is one resource and how much of it is in the bin.
type TrashBucket struct {
	Table string `json:"table"`
	Count int64  `json:"count"`
}

// Buckets counts the deleted rows of every registered resource.
//
// A count per table rather than one total, because "14 things deleted" is not
// something anybody can act on and "9 contacts, 5 groups" is.
func (s *TrashService) Buckets(ctx context.Context) ([]TrashBucket, error) {
	out := make([]TrashBucket, 0)
	for _, table := range s.tables() {
		proto, err := s.Registry.New(table)
		if err != nil {
			continue
		}
		var n int64
		if err := s.deleted(ctx, proto).Count(&n).Error; err != nil {
			// A table the registry knows and the database does not is a
			// migration that has not run yet, not a reason to fail the page.
			continue
		}
		out = append(out, TrashBucket{Table: table, Count: n})
	}
	return out, nil
}

// List returns the deleted rows of one resource, newest deletion first.
func (s *TrashService) List(ctx context.Context, table string, limit, offset int) ([]TrashItem, int64, error) {
	proto, err := s.Registry.New(table)
	if err != nil {
		return nil, 0, fmt.Errorf("unknown resource %q", table)
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var total int64
	if err := s.deleted(ctx, proto).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Into maps rather than the model: this endpoint serves every resource, and
	// the page only needs an id, a name and a date.
	var rows []map[string]any
	if err := s.deleted(ctx, proto).
		Order("deleted_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]TrashItem, 0, len(rows))
	for _, row := range rows {
		deletedAt := asTime(row["deleted_at"])
		items = append(items, TrashItem{
			Table:     table,
			ID:        fmt.Sprint(row["id"]),
			Label:     rowLabel(row),
			DeletedAt: deletedAt,
			ExpiresAt: deletedAt.AddDate(0, 0, TrashRetentionDays),
			Fields:    summaryFields(row),
		})
	}
	return items, total, nil
}

// Restore brings one row back, by clearing the column that hid it.
func (s *TrashService) Restore(ctx context.Context, table, id string) error {
	proto, err := s.Registry.New(table)
	if err != nil {
		return fmt.Errorf("unknown resource %q", table)
	}
	res := s.db(ctx).Unscoped().Model(proto).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("nothing deleted with id %q in %s", id, table)
	}
	return nil
}

// Purge removes one row for good.
func (s *TrashService) Purge(ctx context.Context, table, id string) error {
	proto, err := s.Registry.New(table)
	if err != nil {
		return fmt.Errorf("unknown resource %q", table)
	}
	// Unscoped on a row that is already soft-deleted: the WHERE keeps this from
	// hard-deleting a live row if an id is sent by hand.
	res := s.db(ctx).Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(proto)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("nothing deleted with id %q in %s", id, table)
	}
	return nil
}

// Empty purges every deleted row of one resource, or of all of them when table
// is empty.
func (s *TrashService) Empty(ctx context.Context, table string) (int64, error) {
	tables := s.tables()
	if table != "" {
		if _, err := s.Registry.New(table); err != nil {
			return 0, fmt.Errorf("unknown resource %q", table)
		}
		tables = []string{table}
	}

	var removed int64
	for _, name := range tables {
		proto, err := s.Registry.New(name)
		if err != nil {
			continue
		}
		res := s.db(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(proto)
		if res.Error != nil {
			return removed, res.Error
		}
		removed += res.RowsAffected
	}
	return removed, nil
}

// PurgeExpired removes what the retention window has run out on.
//
// Run from cron. Returns what it removed per table so the job log says
// something other than "done".
func (s *TrashService) PurgeExpired(ctx context.Context) (map[string]int64, error) {
	cutoff := time.Now().AddDate(0, 0, -TrashRetentionDays)
	out := map[string]int64{}
	for _, table := range s.tables() {
		proto, err := s.Registry.New(table)
		if err != nil {
			continue
		}
		res := s.db(ctx).Unscoped().
			Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoff).
			Delete(proto)
		if res.Error != nil {
			return out, res.Error
		}
		if res.RowsAffected > 0 {
			out[table] = res.RowsAffected
		}
	}
	return out, nil
}

// deleted is the query every method here starts from: this model's table, only
// the rows a soft delete hid.
func (s *TrashService) deleted(ctx context.Context, proto any) *gorm.DB {
	return s.db(ctx).Unscoped().Model(proto).Where("deleted_at IS NOT NULL")
}

// tables is the registered resources, in a stable order so two calls to the
// same page do not reshuffle it.
func (s *TrashService) tables() []string {
	if s.Registry == nil {
		return nil
	}
	tables := s.Registry.Tables()
	sort.Strings(tables)
	return tables
}

// rowLabel is what to call a row whose type is not known here.
//
// The same reasoning as the admin's relationship pickers: a person recognises a
// name, a title or a reference, and never a UUID. The id is the last resort
// rather than the first.
func rowLabel(row map[string]any) string {
	for _, key := range []string{"name", "title", "label", "subject", "number", "code", "reference", "email"} {
		if v, ok := row[key]; ok {
			if s := strings.TrimSpace(fmt.Sprint(v)); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return fmt.Sprint(row["id"])
}

// summaryFields is the handful of columns worth showing beside the label.
//
// Everything would mean password hashes and encrypted columns on a page about
// undoing a mistake, so this is an allowlist of the shapes that identify a row
// rather than everything the row holds.
func summaryFields(row map[string]any) map[string]any {
	out := map[string]any{}
	for _, key := range []string{"email", "status", "reference", "number", "code"} {
		if v, ok := row[key]; ok && v != nil {
			if s := strings.TrimSpace(fmt.Sprint(v)); s != "" {
				out[key] = s
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// asTime reads a timestamp back out of a map scan, which hands it over as a
// time.Time on Postgres and as a string on SQLite.
func asTime(v any) time.Time {
	switch t := v.(type) {
	case time.Time:
		return t
	case *time.Time:
		if t != nil {
			return *t
		}
	case string:
		for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05"} {
			if parsed, err := time.Parse(layout, t); err == nil {
				return parsed
			}
		}
	}
	return time.Time{}
}
