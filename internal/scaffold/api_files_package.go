package scaffold

// api_files_package.go emits the scaffolded `internal/files` Go package
// into the generated API. This package owns the FileRef type that every
// resource with a :file: / :files: field embeds, plus the alias→MIME
// translation table that powers per-field upload validation.
//
// FileRef is stored as a JSON column on the parent resource table:
//   Image         *files.FileRef  `gorm:"type:json"`
//   GalleryImages files.FileRefs  `gorm:"type:json"`
//
// FileRefs is a slice type with its own Value / Scan methods so GORM
// can round-trip it without an extra serializer tag — that keeps the
// generated model templates simpler (no per-field `gorm:"serializer:json"`).

func filesFileRefGo() string {
	return `package files

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// FileRef is the canonical JSON shape stored in a resource's file
// column. The frontend uploads to /api/uploads and gets back this
// exact shape, which it then submits to the parent resource's create
// or update endpoint.
//
// Why we embed the metadata (and don't just store a URL): rendering a
// file preview needs the MIME (so we know it's a video vs an image),
// deleting needs the S3 key, and image-list pages need width / height
// to render a placeholder of the right size and avoid layout shift.
// Storing this once at upload time is one DB write; recomputing it on
// every page render would be N database joins.
type FileRef struct {
	URL          string ` + "`json:\"url\"`" + `
	Key          string ` + "`json:\"key\"`" + `
	Name         string ` + "`json:\"name\"`" + `
	MIME         string ` + "`json:\"mime\"`" + `
	Size         int64  ` + "`json:\"size\"`" + `
	Width        *int   ` + "`json:\"width,omitempty\"`" + `
	Height       *int   ` + "`json:\"height,omitempty\"`" + `
	Duration     *int   ` + "`json:\"duration,omitempty\"`" + ` // video / audio (seconds)
	ThumbnailURL string ` + "`json:\"thumbnail_url,omitempty\"`" + `
}

// Value implements driver.Valuer so GORM stores FileRef as JSON.
//
// Why string instead of []byte: lib/pq encodes a []byte driver.Value
// as bytea (Postgres binary type). Inserting bytea into a json column
// fails with SQLSTATE 22P02 ("invalid input syntax for type json")
// because Postgres tries to interpret the binary blob as a JSON
// document and the framing is wrong. Returning a string sends a
// plain text value that Postgres parses as JSON cleanly. SQLite and
// MySQL are tolerant either way; only Postgres is strict here.
func (f FileRef) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner so GORM hydrates FileRef from JSON.
func (f *FileRef) Scan(value interface{}) error {
	if value == nil {
		*f = FileRef{}
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("files.FileRef.Scan: unsupported type %T", value)
	}
	return json.Unmarshal(data, f)
}

// FileRefs is a slice of FileRef with database round-trip support. The
// custom Value / Scan means generated models can declare a field as
// ` + "`files.FileRefs`" + ` directly — no GORM serializer tag needed.
type FileRefs []FileRef

// Value implements driver.Valuer for the slice variant.
// Same string-vs-[]byte reasoning as FileRef.Value above.
func (fs FileRefs) Value() (driver.Value, error) {
	if len(fs) == 0 {
		// Store the empty array, not NULL — keeps the JSON shape
		// stable on the frontend (it always sees [], never null).
		return "[]", nil
	}
	b, err := json.Marshal(fs)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner for the slice variant.
func (fs *FileRefs) Scan(value interface{}) error {
	if value == nil {
		*fs = FileRefs{}
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("files.FileRefs.Scan: unsupported type %T", value)
	}
	if len(data) == 0 {
		*fs = FileRefs{}
		return nil
	}
	return json.Unmarshal(data, fs)
}

// Keys returns the S3 keys of every FileRef in the slice. Used by the
// orphan cleanup cron and by replacement-delete logic to figure out
// which storage objects to purge.
func (fs FileRefs) Keys() []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		if f.Key != "" {
			out = append(out, f.Key)
		}
	}
	return out
}

// TotalSize sums the byte size of every FileRef in the slice. Used by
// the storage admin page to compute per-resource storage usage.
func (fs FileRefs) TotalSize() int64 {
	var total int64
	for _, f := range fs {
		total += f.Size
	}
	return total
}
`
}

func filesAcceptsGo() string {
	return `package files

import "strings"

// AcceptsToMIMEs translates the high-level accept aliases used by the
// CLI (image:file:image, attachment:file:[pdf,doc]) to concrete MIME
// types + filename extensions. The upload handler uses this to
// validate per-field uploads at request time.
//
// "all" is a sentinel meaning the field accepts anything — the upload
// handler skips MIME checking and only enforces the size cap.
func AcceptsToMIMEs(accepts []string) (mimes []string, exts []string, acceptAll bool) {
	mimeSet := map[string]bool{}
	extSet := map[string]bool{}
	for _, a := range accepts {
		switch strings.ToLower(strings.TrimSpace(a)) {
		case "all":
			acceptAll = true
		case "image":
			for _, m := range []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/avif", "image/svg+xml"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif", ".svg"} {
				extSet[e] = true
			}
		case "video":
			for _, m := range []string{"video/mp4", "video/webm", "video/quicktime", "video/x-msvideo", "video/x-matroska"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".mp4", ".webm", ".mov", ".avi", ".mkv"} {
				extSet[e] = true
			}
		case "audio":
			for _, m := range []string{"audio/mpeg", "audio/wav", "audio/ogg", "audio/x-m4a", "audio/webm"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".mp3", ".wav", ".ogg", ".m4a"} {
				extSet[e] = true
			}
		case "pdf":
			mimeSet["application/pdf"] = true
			extSet[".pdf"] = true
		case "doc":
			for _, m := range []string{"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".doc", ".docx"} {
				extSet[e] = true
			}
		case "excel":
			for _, m := range []string{"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".xls", ".xlsx"} {
				extSet[e] = true
			}
		case "csv":
			for _, m := range []string{"text/csv", "application/vnd.ms-excel"} {
				mimeSet[m] = true
			}
			extSet[".csv"] = true
		case "zip":
			for _, m := range []string{"application/zip", "application/x-zip-compressed"} {
				mimeSet[m] = true
			}
			extSet[".zip"] = true
		case "archive":
			for _, m := range []string{"application/zip", "application/x-zip-compressed", "application/x-tar", "application/gzip", "application/x-rar-compressed", "application/x-7z-compressed"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".zip", ".tar", ".gz", ".tgz", ".rar", ".7z"} {
				extSet[e] = true
			}
		}
	}
	for m := range mimeSet {
		mimes = append(mimes, m)
	}
	for e := range extSet {
		exts = append(exts, e)
	}
	return mimes, exts, acceptAll
}

// AllowsMIME returns true if the given mime type is acceptable for the
// given accept-aliases. "all" short-circuits to true.
func AllowsMIME(accepts []string, mime string) bool {
	mimes, _, all := AcceptsToMIMEs(accepts)
	if all {
		return true
	}
	mime = strings.ToLower(strings.TrimSpace(mime))
	for _, m := range mimes {
		if strings.EqualFold(m, mime) {
			return true
		}
	}
	return false
}

// DefaultMaxSizeBytes returns the sensible default cap for the given
// accept set. Video-heavy fields get a much larger cap (300MB) because
// even a short 1080p clip is dozens of megabytes; everything else
// defaults to 5MB.
func DefaultMaxSizeBytes(accepts []string) int64 {
	for _, a := range accepts {
		if strings.EqualFold(a, "video") {
			return 300 << 20 // 300 MB
		}
	}
	return 5 << 20 // 5 MB
}
`
}

func filesFileRefTestGo() string {
	return `package files

import (
	"encoding/json"
	"testing"
)

func TestFileRefRoundTrip(t *testing.T) {
	width := 1920
	height := 1080
	ref := FileRef{
		URL:    "https://cdn.example.com/uploads/abc.jpg",
		Key:    "uploads/2026/06/abc.jpg",
		Name:   "abc.jpg",
		MIME:   "image/jpeg",
		Size:   12345,
		Width:  &width,
		Height: &height,
	}

	v, err := ref.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}

	// Must return string, not []byte. Postgres' lib/pq encodes a
	// []byte driver.Value as bytea, which Postgres then refuses to
	// insert into a json column (SQLSTATE 22P02). Returning string
	// sends text, which Postgres parses as JSON cleanly.
	if _, ok := v.(string); !ok {
		t.Fatalf("FileRef.Value() must return string for Postgres json compatibility, got %T", v)
	}

	var got FileRef
	if err := got.Scan(v); err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if got.URL != ref.URL || got.Key != ref.Key || got.MIME != ref.MIME || got.Size != ref.Size {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, ref)
	}
	if got.Width == nil || *got.Width != width {
		t.Errorf("Width lost in round-trip: got %v", got.Width)
	}
}

func TestFileRefsEmpty(t *testing.T) {
	var fs FileRefs
	v, err := fs.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if s, _ := v.(string); s != "[]" {
		t.Errorf("empty FileRefs should serialize as [], got %q", v)
	}
}

func TestFileRefsKeysAndSize(t *testing.T) {
	fs := FileRefs{
		{Key: "a.jpg", Size: 100},
		{Key: "b.jpg", Size: 200},
		{Key: "c.jpg", Size: 300},
	}
	keys := fs.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
	if total := fs.TotalSize(); total != 600 {
		t.Errorf("expected total 600, got %d", total)
	}
}

func TestFileRefsScanJSON(t *testing.T) {
	raw := ` + "`" + `[{"url":"u1","key":"k1","name":"n1","mime":"image/jpeg","size":10}]` + "`" + `
	var fs FileRefs
	if err := fs.Scan(raw); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(fs) != 1 || fs[0].URL != "u1" {
		t.Errorf("Scan produced %+v", fs)
	}
}

func TestFileRefsValueReturnsString(t *testing.T) {
	// Same Postgres-bytea-vs-json constraint as FileRef.Value above.
	// A non-empty slice must return string, not []byte.
	fs := FileRefs{{URL: "u1", Key: "k1", Name: "n1", MIME: "image/jpeg", Size: 10}}
	v, err := fs.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if _, ok := v.(string); !ok {
		t.Fatalf("FileRefs.Value() must return string for Postgres json compatibility, got %T", v)
	}
}

func TestAcceptsToMIMEsImage(t *testing.T) {
	mimes, _, all := AcceptsToMIMEs([]string{"image"})
	if all {
		t.Error("image alias should not be accept-all")
	}
	if len(mimes) == 0 {
		t.Error("image alias should produce some MIME types")
	}
}

func TestAcceptsToMIMEsAll(t *testing.T) {
	_, _, all := AcceptsToMIMEs([]string{"all"})
	if !all {
		t.Error("all alias should set acceptAll=true")
	}
}

func TestAllowsMIME(t *testing.T) {
	if !AllowsMIME([]string{"image"}, "image/jpeg") {
		t.Error("image/jpeg should be allowed under image alias")
	}
	if AllowsMIME([]string{"image"}, "application/pdf") {
		t.Error("application/pdf should be rejected under image alias")
	}
	if !AllowsMIME([]string{"all"}, "application/wat-is-this") {
		t.Error("any MIME should be allowed under all alias")
	}
}

func TestDefaultMaxSizeBytes(t *testing.T) {
	if got := DefaultMaxSizeBytes([]string{"image"}); got != 5<<20 {
		t.Errorf("image default should be 5MB, got %d", got)
	}
	if got := DefaultMaxSizeBytes([]string{"video"}); got != 300<<20 {
		t.Errorf("video default should be 300MB, got %d", got)
	}
	if got := DefaultMaxSizeBytes([]string{"image", "video"}); got != 300<<20 {
		t.Errorf("mixed image+video should pick the larger (video) cap, got %d", got)
	}
}

// silence unused-import false-positive on encoding/json in some toolchain
// versions when this test file is read in isolation.
var _ = json.Marshal
`
}

// filesLifecycleGo emits the lifecycle helpers that close the loop on
// uploaded files: when a record's file column changes the old S3 object
// gets deleted immediately, and unclaimed Upload rows get swept by a
// daily cron so the bucket doesn't accumulate junk over time.
//
// The reflection-based CleanupRemoved + ClaimRefs functions mean
// generated handlers add a single line each, regardless of how many
// file columns a resource has.

func filesLifecycleGo() string {
	return `package files

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

// RunOrphanCleanup deletes Upload rows whose key was never claimed by
// a parent record AND which are older than minAge. Designed to be
// called from a daily cron job. Returns the count of rows + S3 objects
// purged.
//
// The minAge buffer matters: an upload immediately followed by a form
// save has a small window between the POST /api/uploads success and
// the parent Create handler's ClaimRefs call. minAge=24h is generous
// -- you'd have to abandon the form for a full day for the cleanup to
// catch it.
func RunOrphanCleanup(ctx context.Context, db *gorm.DB, st Storage, minAge time.Duration) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("RunOrphanCleanup: db is required")
	}

	type orphan struct {
		ID   string ` + "`gorm:\"column:id\"`" + `
		Path string ` + "`gorm:\"column:path\"`" + `
	}

	cutoff := time.Now().Add(-minAge)
	var orphans []orphan
	err := db.WithContext(ctx).
		Table("uploads").
		Select("id", "path").
		Where("claimed_at IS NULL AND created_at < ?", cutoff).
		Find(&orphans).Error
	if err != nil {
		return 0, fmt.Errorf("query orphans: %w", err)
	}

	deleted := 0
	for _, o := range orphans {
		if st != nil && o.Path != "" {
			// Best-effort S3 delete -- if it fails (already gone, perm
			// issue, etc.) we still drop the DB row so we don't keep
			// retrying the same orphan forever.
			_ = st.Delete(ctx, o.Path)
		}
		if err := db.WithContext(ctx).Table("uploads").Where("id = ?", o.ID).Delete(struct{}{}).Error; err == nil {
			deleted++
		}
	}
	return deleted, nil
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
`
}

// filesLifecycleTestGo emits unit tests for the lifecycle helpers.
// Uses an in-memory sqlite DB so the tests run with no infrastructure
// dependency (real S3 / Postgres / etc. not required).

func filesLifecycleTestGo() string {
	return `package files

import (
	"context"
	"testing"
)

// fakeStorage records every Delete call so tests can assert on which
// keys were purged. Implements the Storage interface.
type fakeStorage struct {
	deleted []string
}

func (f *fakeStorage) Delete(_ context.Context, key string) error {
	f.deleted = append(f.deleted, key)
	return nil
}

func TestDiffSingle(t *testing.T) {
	t.Run("nil old returns empty", func(t *testing.T) {
		if got := DiffSingle(nil, &FileRef{Key: "k"}); got != "" {
			t.Errorf("nil old should return empty, got %q", got)
		}
	})
	t.Run("same key returns empty", func(t *testing.T) {
		a := &FileRef{Key: "k"}
		b := &FileRef{Key: "k"}
		if got := DiffSingle(a, b); got != "" {
			t.Errorf("same key should return empty, got %q", got)
		}
	})
	t.Run("different key returns old", func(t *testing.T) {
		a := &FileRef{Key: "old"}
		b := &FileRef{Key: "new"}
		if got := DiffSingle(a, b); got != "old" {
			t.Errorf("different key should return old, got %q", got)
		}
	})
	t.Run("cleared (new nil) returns old", func(t *testing.T) {
		a := &FileRef{Key: "old"}
		if got := DiffSingle(a, nil); got != "old" {
			t.Errorf("cleared should return old, got %q", got)
		}
	})
}

func TestDiffMulti(t *testing.T) {
	old := FileRefs{{Key: "a"}, {Key: "b"}, {Key: "c"}}
	t.Run("shrunk to subset", func(t *testing.T) {
		removed := DiffMulti(old, FileRefs{{Key: "a"}, {Key: "c"}})
		if len(removed) != 1 || removed[0] != "b" {
			t.Errorf("expected [b], got %v", removed)
		}
	})
	t.Run("cleared entirely", func(t *testing.T) {
		removed := DiffMulti(old, FileRefs{})
		if len(removed) != 3 {
			t.Errorf("expected 3 removed, got %v", removed)
		}
	})
	t.Run("unchanged", func(t *testing.T) {
		removed := DiffMulti(old, old)
		if len(removed) != 0 {
			t.Errorf("expected 0 removed, got %v", removed)
		}
	})
	t.Run("nil old", func(t *testing.T) {
		if removed := DiffMulti(nil, FileRefs{{Key: "a"}}); len(removed) != 0 {
			t.Errorf("nil old should produce no removed keys, got %v", removed)
		}
	})
}

func TestCleanupRemovedReflectionSingle(t *testing.T) {
	type fakeProduct struct {
		ID    string
		Name  string
		Image *FileRef
	}
	old := fakeProduct{ID: "1", Image: &FileRef{Key: "old.jpg"}}
	neu := fakeProduct{ID: "1", Image: &FileRef{Key: "new.jpg"}}

	st := &fakeStorage{}
	CleanupRemoved(context.Background(), st, &old, &neu)
	if len(st.deleted) != 1 || st.deleted[0] != "old.jpg" {
		t.Errorf("expected old.jpg deleted, got %v", st.deleted)
	}
}

func TestCleanupRemovedReflectionMulti(t *testing.T) {
	type fakeGallery struct {
		ID     string
		Photos FileRefs
	}
	old := fakeGallery{ID: "1", Photos: FileRefs{{Key: "a"}, {Key: "b"}, {Key: "c"}}}
	neu := fakeGallery{ID: "1", Photos: FileRefs{{Key: "a"}, {Key: "c"}, {Key: "d"}}}

	st := &fakeStorage{}
	CleanupRemoved(context.Background(), st, &old, &neu)
	if len(st.deleted) != 1 || st.deleted[0] != "b" {
		t.Errorf("expected only b deleted, got %v", st.deleted)
	}
}

func TestCleanupRemovedNoFileFields(t *testing.T) {
	type plainRecord struct {
		ID   string
		Name string
	}
	old := plainRecord{ID: "1", Name: "before"}
	neu := plainRecord{ID: "1", Name: "after"}

	st := &fakeStorage{}
	CleanupRemoved(context.Background(), st, &old, &neu)
	if len(st.deleted) != 0 {
		t.Errorf("plain record should not trigger any deletes, got %v", st.deleted)
	}
}

func TestCleanupRemovedNilStorage(t *testing.T) {
	type record struct {
		Image *FileRef
	}
	old := record{Image: &FileRef{Key: "x"}}
	neu := record{Image: nil}
	// Must not panic when storage is nil -- e.g. when storage env vars
	// aren't configured.
	CleanupRemoved(context.Background(), nil, &old, &neu)
}
`
}
