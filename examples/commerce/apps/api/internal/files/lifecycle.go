package files

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// Storage is the subset of the storage package we need here -- declared
// as an interface so this file doesn't pull in the heavy AWS SDK just
// for the cleanup path. The concrete *storage.Storage satisfies it.
type Storage interface {
	Delete(ctx context.Context, key string) error
}

// DiffSingle returns the S3 key of an old FileRef if it was replaced
// or cleared by the new value, or nil if both refer to the same file.
// Used during Update handlers to figure out which S3 object (if any)
// to purge after a single-file field is reassigned.
func DiffSingle(old, neu *FileRef) string {
	if old == nil || old.Key == "" {
		return ""
	}
	if neu != nil && neu.Key == old.Key {
		return ""
	}
	return old.Key
}

// DiffMulti returns the S3 keys present in old but not in new. Used by
// multi-file fields when the gallery shrinks or some entries get
// swapped out -- those S3 objects can be deleted right away.
func DiffMulti(old, neu FileRefs) []string {
	if len(old) == 0 {
		return nil
	}
	newKeys := make(map[string]struct{}, len(neu))
	for _, r := range neu {
		if r.Key != "" {
			newKeys[r.Key] = struct{}{}
		}
	}
	var removed []string
	for _, r := range old {
		if r.Key == "" {
			continue
		}
		if _, kept := newKeys[r.Key]; !kept {
			removed = append(removed, r.Key)
		}
	}
	return removed
}

// CleanupRemoved walks the fields of old + neu (must be pointers to
// the same struct type), finds *FileRef + FileRefs columns, computes
// the diff, and deletes the removed S3 objects via storage.
//
// Errors per key are swallowed -- a delete failure for one orphan
// shouldn't fail the parent UPDATE. Real failures get logged via the
// returned err string for upstream observability.
func CleanupRemoved(ctx context.Context, st Storage, old, neu interface{}) {
	if st == nil || old == nil || neu == nil {
		return
	}
	oldVal := derefPtr(reflect.ValueOf(old))
	newVal := derefPtr(reflect.ValueOf(neu))
	if !oldVal.IsValid() || !newVal.IsValid() {
		return
	}
	if oldVal.Type() != newVal.Type() {
		return
	}

	t := oldVal.Type()
	for i := 0; i < t.NumField(); i++ {
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		// *FileRef -- single-file column.
		if oldField.Kind() == reflect.Ptr && oldField.Type() == reflect.TypeOf((*FileRef)(nil)) {
			var oldRef, newRef *FileRef
			if !oldField.IsNil() {
				oldRef = oldField.Interface().(*FileRef)
			}
			if !newField.IsNil() {
				newRef = newField.Interface().(*FileRef)
			}
			if k := DiffSingle(oldRef, newRef); k != "" {
				_ = st.Delete(ctx, k)
			}
			continue
		}

		// FileRefs (slice) -- multi-file column.
		if oldField.Type() == reflect.TypeOf(FileRefs{}) {
			oldRefs := oldField.Interface().(FileRefs)
			newRefs := newField.Interface().(FileRefs)
			for _, k := range DiffMulti(oldRefs, newRefs) {
				_ = st.Delete(ctx, k)
			}
			continue
		}
	}
}

// ClaimRefs walks a record's file columns and marks each referenced
// Upload row as claimed (claimed_at = now()). The orphan-cleanup cron
// uses claimed_at to distinguish in-use uploads from abandoned ones
// (a user uploaded then closed the browser before saving the parent
// form, for example).
//
// Safe to call on records that have no file columns -- the function
// just iterates zero times.
func ClaimRefs(ctx context.Context, db *gorm.DB, record interface{}) {
	if db == nil || record == nil {
		return
	}
	v := derefPtr(reflect.ValueOf(record))
	if !v.IsValid() {
		return
	}
	t := v.Type()

	var keys []string
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Ptr && f.Type() == reflect.TypeOf((*FileRef)(nil)) {
			if !f.IsNil() {
				keys = append(keys, f.Interface().(*FileRef).Key)
			}
			continue
		}
		if f.Type() == reflect.TypeOf(FileRefs{}) {
			for _, r := range f.Interface().(FileRefs) {
				if r.Key != "" {
					keys = append(keys, r.Key)
				}
			}
			continue
		}
	}
	if len(keys) == 0 {
		return
	}

	// Update claimed_at on all matched Upload rows. We don't care about
	// rows that don't match (uploads from another source could be in
	// the same table) -- they stay unclaimed and will be considered
	// orphans, which is correct.
	_ = db.WithContext(ctx).
		Table("uploads").
		Where("path IN ?", keys).
		Update("claimed_at", time.Now()).
		Error
}

// orphanPage is how many orphans one pass of RunOrphanCleanup reads, deletes
// from storage in one request and deletes from the table in one statement.
const orphanPage = 1000

// batchDeleter is a Storage that deletes many keys in one request, as
// storage.Storage does with S3's DeleteObjects.
type batchDeleter interface {
	DeleteMany(ctx context.Context, keys []string) error
}

// RunOrphanCleanup deletes Upload rows whose key was never claimed by
// a parent record AND which are older than minAge, with their stored
// objects. Designed to be called from a daily cron job. Returns the
// count of rows purged.
//
// The minAge buffer matters: an upload immediately followed by a form
// save has a small window between the POST /api/uploads success and
// the parent Create handler's ClaimRefs call. minAge=24h is generous
// -- you'd have to abandon the form for a full day for the cleanup to
// catch it.
//
// Orphans go orphanPage at a time. The cleanup used to load every orphan
// at once and delete each with a storage request and a statement of its own.
func RunOrphanCleanup(ctx context.Context, db *gorm.DB, st Storage, minAge time.Duration) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("RunOrphanCleanup: db is required")
	}

	type orphan struct {
		ID   string `gorm:"column:id"`
		Path string `gorm:"column:path"`
	}

	cutoff := time.Now().Add(-minAge)
	deleted, after := 0, ""
	for {
		var page []orphan
		err := db.WithContext(ctx).
			Table("uploads").
			Select("id", "path").
			Where("claimed_at IS NULL AND created_at < ? AND id > ?", cutoff, after).
			Order("id").
			Limit(orphanPage).
			Find(&page).Error
		if err != nil {
			return deleted, fmt.Errorf("query orphans: %w", err)
		}
		if len(page) == 0 {
			return deleted, nil
		}
		after = page[len(page)-1].ID

		ids := make([]string, 0, len(page))
		keys := make([]string, 0, len(page))
		for _, o := range page {
			ids = append(ids, o.ID)
			if o.Path != "" {
				keys = append(keys, o.Path)
			}
		}
		// Best-effort, as it always was: a row whose object could not be
		// deleted (already gone, a permissions issue) is still dropped, so the
		// same orphan is not retried forever.
		if st != nil && len(keys) > 0 {
			if many, ok := st.(batchDeleter); ok {
				_ = many.DeleteMany(ctx, keys)
			} else {
				for _, key := range keys {
					_ = st.Delete(ctx, key)
				}
			}
		}
		res := db.WithContext(ctx).Table("uploads").Where("id IN ?", ids).Delete(struct{}{})
		if res.Error != nil {
			return deleted, fmt.Errorf("delete orphans: %w", res.Error)
		}
		deleted += int(res.RowsAffected)
		if len(page) < orphanPage {
			return deleted, nil
		}
	}
}

// derefPtr unwraps a single level of pointer indirection so callers
// can pass either a struct or a pointer-to-struct to the cleanup
// helpers without writing two code paths.
func derefPtr(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}
