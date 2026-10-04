package storage

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func testLocalDisk(t *testing.T) *LocalDisk {
	t.Helper()
	d, err := NewLocalDisk(LocalConfig{Root: t.TempDir(), PublicURL: "http://localhost:8080/files", Secret: "a-test-secret-long-enough-to-sign-with"})
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// formFile builds the *multipart.FileHeader a handler gets from a form.
func formFile(t *testing.T, filename, contentType string, body []byte) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", contentType)
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(body); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	form, err := multipart.NewReader(&buf, w.Boundary()).ReadForm(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })
	return form.File["file"][0]
}

func pngBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.Set(1, 1, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestStoreGeneratesTheKeyAndSniffsTheType(t *testing.T) {
	ctx := context.Background()
	disk := testLocalDisk(t)

	// Declared as a PDF, named to climb out of the store: the key is generated
	// and the type comes from the bytes.
	key, err := Store(ctx, disk, "uploads", formFile(t, "../../evil.PNG", "application/pdf", pngBytes(t)), StoreOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^uploads/\d{4}/\d{2}/[0-9a-f-]{36}\.png$`).MatchString(key) {
		t.Fatalf("key = %q, want uploads/<yyyy>/<mm>/<uuid>.png", key)
	}
	obj, err := disk.Stat(ctx, key)
	if err != nil || obj.ContentType != "image/png" {
		t.Fatalf("Stat = %+v, %v", obj, err)
	}

	named, err := StoreAs(ctx, disk, "avatars", formFile(t, "me.png", "image/png", pngBytes(t)), "user-1.png", StoreOptions{})
	if err != nil || named != "avatars/user-1.png" {
		t.Fatalf("StoreAs = %q, %v", named, err)
	}
	if _, err := StoreAs(ctx, disk, "avatars", formFile(t, "me.png", "image/png", pngBytes(t)), "../x.png", StoreOptions{}); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("StoreAs with a climbing name = %v, want ErrInvalidKey", err)
	}
}

func TestStoreRefuses(t *testing.T) {
	ctx := context.Background()
	disk := testLocalDisk(t)
	cases := []struct {
		name string
		file *multipart.FileHeader
		opts StoreOptions
		want error
	}{
		{"a file over the limit", formFile(t, "a.txt", "text/plain", []byte("12345")), StoreOptions{MaxSize: 4}, ErrFileTooLarge},
		{"HTML, whatever it claims", formFile(t, "a.txt", "text/plain", []byte("<!DOCTYPE html><script>alert(1)</script>")), StoreOptions{}, ErrFileTypeNotAllowed},
		{"an image that is not one", formFile(t, "a.png", "image/png", []byte("MZ not a png")), StoreOptions{}, ErrContentMismatch},
		{"a type Allow refuses", formFile(t, "a.txt", "text/plain", []byte("hello")), StoreOptions{Allow: func(ct string) bool { return strings.HasPrefix(ct, "image/") }}, ErrFileTypeNotAllowed},
		{"a private file under a public prefix", formFile(t, "a.txt", "text/plain", []byte("hello")), StoreOptions{Visibility: VisibilityPrivate}, ErrVisibilityMismatch},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := Store(ctx, disk, "uploads", c.file, c.opts); !errors.Is(err, c.want) {
				t.Fatalf("Store = %v, want %v", err, c.want)
			}
		})
	}
	if objects, err := disk.List(ctx, "uploads/"); err != nil || len(objects) != 0 {
		t.Fatalf("a refused file was stored: %v, %v", objects, err)
	}
}

func TestVisibilityFollowsThePublicPrefixes(t *testing.T) {
	saved := PublicPrefixes
	t.Cleanup(func() { PublicPrefixes = saved })

	ctx := context.Background()
	disk := testLocalDisk(t)
	if err := disk.Put(ctx, "backups/a.zip", strings.NewReader("x"), PutOptions{Visibility: VisibilityPublic}); !errors.Is(err, ErrVisibilityMismatch) {
		t.Fatalf("a public Put under backups/ = %v, want ErrVisibilityMismatch", err)
	}
	if err := disk.Put(ctx, "backups/a.zip", strings.NewReader("x"), PutOptions{Visibility: VisibilityPrivate}); err != nil {
		t.Fatal(err)
	}

	SetPublicPrefixes([]string{"media", " avatars/ ", ""})
	if !IsPublicKey("media/a.png") || !IsPublicKey("avatars/a.png") || IsPublicKey("uploads/a.png") {
		t.Fatalf("PublicPrefixes = %v", PublicPrefixes)
	}
	if !strings.Contains(BucketPolicy("b"), "arn:aws:s3:::b/media/*") || strings.Contains(BucketPolicy("b"), "uploads/") {
		t.Fatalf("the bucket policy does not follow the prefixes: %s", BucketPolicy("b"))
	}
	SetPublicPrefixes(nil)
	if !IsPublicKey("media/a.png") {
		t.Fatal("an empty list replaced the prefixes")
	}
}

func serve(t *testing.T, disk Disk, key string, disposition Disposition, method string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(method, "/download", nil)
	for k, v := range headers {
		c.Request.Header.Set(k, v)
	}
	ServeFileAs(c, disk, key, disposition, "report final.txt")
	return rec
}

// onlyReader hides the file's Seek, as a bucket's body would.
type onlyReader struct{ Disk }

func (d onlyReader) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	body, err := d.Disk.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return struct {
		io.Reader
		io.Closer
	}{body, body}, nil
}

func TestServeFile(t *testing.T) {
	ctx := context.Background()
	local := testLocalDisk(t)
	if err := local.Put(ctx, "private/report.txt", strings.NewReader("0123456789"), PutOptions{}); err != nil {
		t.Fatal(err)
	}

	for name, disk := range map[string]Disk{"a seekable body": local, "a body that cannot seek": onlyReader{local}} {
		t.Run(name, func(t *testing.T) {
			full := serve(t, disk, "private/report.txt", Attachment, http.MethodGet, nil)
			if full.Code != http.StatusOK || full.Body.String() != "0123456789" {
				t.Fatalf("GET = %d %q", full.Code, full.Body.String())
			}
			h := full.Header()
			if !strings.HasPrefix(h.Get("Content-Type"), "text/plain") || h.Get("Content-Length") != "10" || h.Get("Last-Modified") == "" || h.Get("Accept-Ranges") != "bytes" {
				t.Fatalf("headers = %v", h)
			}
			if got := h.Get("Content-Disposition"); got != `attachment; filename="report final.txt"` {
				t.Fatalf("Content-Disposition = %q", got)
			}

			part := serve(t, disk, "private/report.txt", Inline, http.MethodGet, map[string]string{"Range": "bytes=2-5"})
			if part.Code != http.StatusPartialContent || part.Body.String() != "2345" || part.Header().Get("Content-Range") != "bytes 2-5/10" || part.Header().Get("Content-Length") != "4" {
				t.Fatalf("Range 2-5 = %d %q %v", part.Code, part.Body.String(), part.Header())
			}
			if !strings.HasPrefix(part.Header().Get("Content-Disposition"), "inline") {
				t.Fatalf("Content-Disposition = %q", part.Header().Get("Content-Disposition"))
			}
			if tail := serve(t, disk, "private/report.txt", Inline, http.MethodGet, map[string]string{"Range": "bytes=-3"}); tail.Body.String() != "789" {
				t.Fatalf("Range -3 = %d %q", tail.Code, tail.Body.String())
			}
			if past := serve(t, disk, "private/report.txt", Inline, http.MethodGet, map[string]string{"Range": "bytes=10-"}); past.Code != http.StatusRequestedRangeNotSatisfiable || past.Header().Get("Content-Range") != "bytes */10" {
				t.Fatalf("Range 10- = %d %v", past.Code, past.Header())
			}
			stale := serve(t, disk, "private/report.txt", Inline, http.MethodGet, map[string]string{"Range": "bytes=2-5", "If-Range": "Mon, 01 Jan 2001 00:00:00 GMT"})
			if stale.Code != http.StatusOK || stale.Body.Len() != 10 {
				t.Fatalf("a stale If-Range = %d, %d bytes, want the whole file", stale.Code, stale.Body.Len())
			}
			if head := serve(t, disk, "private/report.txt", Inline, http.MethodHead, nil); head.Code != http.StatusOK || head.Body.Len() != 0 || head.Header().Get("Content-Length") != "10" {
				t.Fatalf("HEAD = %d, %d bytes, %v", head.Code, head.Body.Len(), head.Header())
			}
			since := time.Now().Add(time.Hour).UTC().Format(http.TimeFormat)
			if cached := serve(t, disk, "private/report.txt", Inline, http.MethodGet, map[string]string{"If-Modified-Since": since}); cached.Code != http.StatusNotModified {
				t.Fatalf("If-Modified-Since = %d, want 304", cached.Code)
			}
			if missing := serve(t, disk, "private/missing.txt", Inline, http.MethodGet, nil); missing.Code != http.StatusNotFound {
				t.Fatalf("a missing file = %d, want 404", missing.Code)
			}
		})
	}
}

func TestParseRange(t *testing.T) {
	cases := []struct {
		header        string
		start, length int64
		partial, ok   bool
	}{
		{"", 0, 100, false, true},
		{"bytes=0-0", 0, 1, true, true},
		{"bytes=90-", 90, 10, true, true},
		{"bytes=90-500", 90, 10, true, true},
		{"bytes=-10", 90, 10, true, true},
		{"bytes=-500", 0, 100, true, true},
		{"bytes=100-", 0, 0, false, false},
		{"bytes=5-1", 0, 100, false, true},
		{"bytes=0-1,5-6", 0, 100, false, true},
		{"items=0-1", 0, 100, false, true},
	}
	for _, c := range cases {
		start, length, partial, ok := parseRange(c.header, 100)
		if start != c.start || length != c.length || partial != c.partial || ok != c.ok {
			t.Errorf("parseRange(%q) = %d, %d, %v, %v", c.header, start, length, partial, ok)
		}
	}
}

func TestRegistry(t *testing.T) {
	r := &Registry{}
	def := Wrap(testLocalDisk(t))
	backups := Wrap(testLocalDisk(t))
	r.SetDefault(def)
	r.Add("Backups", backups)
	if r.Get("") != def || r.Get("default") != def || r.Get("backups") != backups {
		t.Fatal("Get does not return what was added")
	}
	if r.Get("archive") != nil {
		t.Fatal("a name nobody configured fell back to a store")
	}
	if names := r.Names(); len(names) != 1 || names[0] != "backups" {
		t.Fatalf("Names = %v", names)
	}
}

// A named local disk has no route of its own: the default local disk serves
// it under /files/_disks/<name>/, with the same signature check.
func TestANamedLocalDiskIsServedUnderTheDefaultRoute(t *testing.T) {
	ctx := context.Background()
	var def *LocalDisk
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { def.ServeHTTP(w, r) }))
	t.Cleanup(server.Close)

	var err error
	def, err = NewLocalDisk(LocalConfig{Root: t.TempDir(), PublicURL: server.URL + "/files", Secret: "secret-one"})
	if err != nil {
		t.Fatal(err)
	}
	backups, err := NewLocalDisk(LocalConfig{Root: t.TempDir(), PublicURL: server.URL + "/files/_disks/backups", Secret: "secret-one"})
	if err != nil {
		t.Fatal(err)
	}
	saved := Disks
	Disks = &Registry{}
	t.Cleanup(func() { Disks = saved })
	Disks.SetDefault(Wrap(def))
	Disks.Add("backups", Wrap(backups))

	if err := backups.Put(ctx, "backups/a.zip", strings.NewReader("archive"), PutOptions{Visibility: VisibilityPrivate}); err != nil {
		t.Fatal(err)
	}
	signed, err := backups.TemporaryURL(ctx, "backups/a.zip", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if status, body := fetch(t, signed); status != http.StatusOK || body != "archive" {
		t.Fatalf("GET signed named-disk URL = %d %q", status, body)
	}
	if status, _ := fetch(t, backups.URL("backups/a.zip")); status == http.StatusOK {
		t.Fatal("a private file on a named disk is served without a signature")
	}
	if err := def.Put(ctx, "_disks/backups/x.txt", strings.NewReader("x"), PutOptions{}); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("a default-disk key under _disks/ = %v, want ErrInvalidKey", err)
	}
}
