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
func (f FileRef) Value() (driver.Value, error) {
	return json.Marshal(f)
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
func (fs FileRefs) Value() (driver.Value, error) {
	if len(fs) == 0 {
		// Store the empty array, not NULL — keeps the JSON shape
		// stable on the frontend (it always sees [], never null).
		return "[]", nil
	}
	return json.Marshal(fs)
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
