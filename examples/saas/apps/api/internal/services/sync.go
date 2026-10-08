package services

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SyncService is the row store behind the offline sync protocol.
//
// Every other service in this package knows what it is reading. This one does
// not: a mobile client pushes changes to whatever models the project registered
// for sync, so each method takes the destination the caller built by reflection
// and the registry decided the type of. The handler keeps that reflection, and
// the decisions about the protocol, which is what a push and a pull mean.
//
// What is here is the part that is a decision about data rather than about the
// request: that a pull is keyset-paginated on (updated_at, id), that it reads
// soft-deleted rows because those are the tombstones a client needs, and that a
// save does not touch associations or created_at.
type SyncService struct {
	DB *gorm.DB
}

func (s *SyncService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// ByIDs reads the rows a push's updates and deletes name, into dest.
//
// One query per model rather than one per change. A row it does not load is read
// by its own change, so a failure here costs queries, not correctness.
func (s *SyncService) ByIDs(ctx context.Context, dest interface{}, ids []string) error {
	return s.db(ctx).Where("id IN ?", ids).Find(dest).Error
}

// ByID reads one row into dest.
func (s *SyncService) ByID(ctx context.Context, dest interface{}, id string) error {
	return s.db(ctx).First(dest, "id = ?", id).Error
}

// Create inserts a row a client created offline.
func (s *SyncService) Create(ctx context.Context, obj interface{}) error {
	return s.db(ctx).Create(obj).Error
}

// Save writes a row a client changed offline.
//
// Associations are omitted, because a sync payload carries columns and not
// relations, and GORM would otherwise upsert whatever an association field
// happens to hold. created_at is omitted for the same reason: a client that
// round-trips a row should not be able to rewrite when it was made.
func (s *SyncService) Save(ctx context.Context, obj interface{}) error {
	return s.db(ctx).Omit(clause.Associations, "CreatedAt").Save(obj).Error
}

// Delete soft-deletes a row, which is what leaves the tombstone a pull sends to
// every other device.
func (s *SyncService) Delete(ctx context.Context, obj interface{}, id string) error {
	return s.db(ctx).Delete(obj, "id = ?", id).Error
}

// Page reads one page of changes into dest, for a pull.
//
// A row's change time is its updated_at, which a soft delete sets too (see
// internal/sync), so one indexed column carries edits and deletes alike. Ordering
// by the later of updated_at and deleted_at sorted the whole table on every pull,
// and MySQL has no two-argument MAX to write it with.
//
// Keyset on (updated_at, id), so rows sharing a timestamp are not lost at a page
// boundary, as they were with a cursor of the time alone.
//
// Unscoped, because the soft-deleted rows are the tombstones: a client that never
// hears about a deletion keeps the row forever.
func (s *SyncService) Page(ctx context.Context, model, dest interface{}, since time.Time, afterID string, limit int) error {
	q := s.db(ctx).Unscoped().Model(model)
	if !since.IsZero() {
		if afterID == "" {
			q = q.Where("updated_at > ?", since)
		} else {
			q = q.Where("updated_at > ? OR (updated_at = ? AND id > ?)", since, since, afterID)
		}
	}
	return q.Order("updated_at asc, id asc").Limit(limit).Find(dest).Error
}
