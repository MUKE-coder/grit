package files

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
	raw := `[{"url":"u1","key":"k1","name":"n1","mime":"image/jpeg","size":10}]`
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
