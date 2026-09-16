package scaffold

// Storage helpers, named disks and visibility (issue 003, C2 to C4).
//
// storageStoreGo is Store, StoreAs and ServeFile: what a handler needs to take a
// file from a form and hand one back. storageDisksGo is Open and the Disks
// registry, so a backup can go to a private bucket of its own.
// storageStoreTestGo runs in the generated project against the local disk.

func storageStoreGo() string {
	return `package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DefaultMaxFileSize is Store's size limit when StoreOptions.MaxSize is zero.
const DefaultMaxFileSize = 50 << 20

var (
	// ErrFileTooLarge is returned, wrapped, for a file over the size limit.
	ErrFileTooLarge = errors.New("storage: file is over the size limit")
	// ErrFileTypeNotAllowed is returned, wrapped, for a type Store refuses: HTML
	// and SVG always, and anything StoreOptions.Allow says no to.
	ErrFileTypeNotAllowed = errors.New("storage: file type not allowed")
	// ErrContentMismatch is returned for a file that claims to be an image and
	// is not one.
	ErrContentMismatch = errors.New("storage: file content does not match its declared type")
)

// StoreOptions limit what Store accepts.
type StoreOptions struct {
	// MaxSize is the most bytes the file may have. Zero is DefaultMaxFileSize.
	MaxSize int64
	// Allow, when set, is asked about the sniffed content type.
	Allow func(contentType string) bool
	// Visibility is passed to Disk.Put. Empty follows the key's prefix.
	Visibility Visibility
}

// Store saves an uploaded file under dir with a generated name,
// <dir>/<yyyy>/<mm>/<uuid><ext>, and returns its key.
//
// The name the file arrived with never reaches the key, so a file called
// "../../x.html" or "invoice (final).pdf" is stored as a UUID. Keep the
// original name on your own row. The content type is sniffed from the bytes
// rather than taken from the form: see DetectContentType.
func Store(ctx context.Context, disk Disk, dir string, file *multipart.FileHeader, opts StoreOptions) (string, error) {
	return StoreAs(ctx, disk, dir, file, "", opts)
}

// StoreAs is Store with a name of your choosing: <dir>/<name>. An empty name
// is generated as Store generates one.
func StoreAs(ctx context.Context, disk Disk, dir string, file *multipart.FileHeader, name string, opts StoreOptions) (string, error) {
	if file == nil {
		return "", errors.New("storage: no file to store")
	}
	limit := opts.MaxSize
	if limit <= 0 {
		limit = DefaultMaxFileSize
	}
	if file.Size > limit {
		return "", fmt.Errorf("%w: %s is %d bytes, the limit is %d", ErrFileTooLarge, file.Filename, file.Size, limit)
	}

	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("opening %q: %w", file.Filename, err)
	}
	defer f.Close()

	contentType, err := DetectContentType(f, file.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	if opts.Allow != nil && !opts.Allow(contentType) {
		return "", fmt.Errorf("%w: %s", ErrFileTypeNotAllowed, contentType)
	}

	var key string
	if name == "" {
		key = NewKey(dir, Extension(file.Filename, contentType))
	} else {
		if err := checkKey(name); err != nil {
			return "", err
		}
		key = strings.TrimPrefix(path.Join(dir, name), "/")
	}
	if err := disk.Put(ctx, key, f, PutOptions{ContentType: contentType, Visibility: opts.Visibility}); err != nil {
		return "", err
	}
	return key, nil
}

// DetectContentType reads the first bytes of r to learn what it really is,
// reconciles that with the type the client declared, and rewinds r.
//
// The declared type is trivially spoofed. HTML and SVG are refused whatever they
// claim, because served from your origin they run script; a declared image must
// sniff as one; a sniffed image wins over the declaration. Anything else keeps
// the declared type, because many valid documents sniff only as
// application/octet-stream. Parameters such as ";codecs=opus" are dropped.
func DetectContentType(r io.ReadSeeker, declared string) (string, error) {
	head := make([]byte, 512)
	n, err := io.ReadFull(r, head)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", fmt.Errorf("reading the file: %w", err)
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("rewinding the file: %w", err)
	}
	detected := bareType(http.DetectContentType(head[:n]))
	declared = bareType(declared)

	if detected == "text/html" || detected == "image/svg+xml" {
		return "", fmt.Errorf("%w: %s", ErrFileTypeNotAllowed, detected)
	}
	if strings.HasPrefix(declared, "image/") && !strings.HasPrefix(detected, "image/") {
		return "", fmt.Errorf("%w: declared %s, the bytes are %s", ErrContentMismatch, declared, detected)
	}
	if strings.HasPrefix(detected, "image/") || declared == "" {
		return detected, nil
	}
	return declared, nil
}

func bareType(contentType string) string {
	return strings.ToLower(strings.TrimSpace(strings.SplitN(contentType, ";", 2)[0]))
}

// NewKey is a fresh key under dir: <dir>/<yyyy>/<mm>/<uuid><ext>.
func NewKey(dir, ext string) string {
	name := time.Now().UTC().Format("2006/01") + "/" + uuid.NewString() + ext
	if dir = strings.Trim(dir, "/"); dir != "" {
		return dir + "/" + name
	}
	return name
}

var preferredExtensions = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/gif":       ".gif",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
	"video/mp4":       ".mp4",
	"text/plain":      ".txt",
	"text/csv":        ".csv",
}

// Extension is the extension a stored file keeps: the uploaded name's, when it
// is a plain one, or the usual one for the content type.
func Extension(filename, contentType string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if plainExtension(ext) {
		return ext
	}
	if known, ok := preferredExtensions[contentType]; ok {
		return known
	}
	if exts, err := mime.ExtensionsByType(contentType); err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ""
}

func plainExtension(ext string) bool {
	if len(ext) < 2 || len(ext) > 16 || ext[0] != '.' {
		return false
	}
	for _, r := range ext[1:] {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// Disposition tells a browser whether to show a served file or save it.
type Disposition string

const (
	// Inline shows the file in the browser where it can.
	Inline Disposition = "inline"
	// Attachment saves the file under its name.
	Attachment Disposition = "attachment"
)

// RangeGetter is a Disk that can open part of a file, so ServeFile answers a
// Range request without reading the bytes before it. S3Disk is one. A disk that
// is not is still served, by skipping to the range.
type RangeGetter interface {
	GetRange(ctx context.Context, key string, offset, length int64) (io.ReadCloser, error)
}

// ServeFile streams the file at key through the API, named after the last
// segment of its key. See ServeFileAs.
func ServeFile(c *gin.Context, disk Disk, key string, disposition Disposition) {
	ServeFileAs(c, disk, key, disposition, path.Base(key))
}

// ServeFileAs streams the file at key with Content-Type, Content-Length,
// Last-Modified and Content-Disposition, answers If-Modified-Since with 304
// and a single byte range with 206, so a video seeks and a large download
// resumes.
//
// It is how a private file reaches a user you have already authorised, on any
// driver: the bytes pass through the API, so nothing about the file is public.
// A missing file is 404. The response carries a sandboxing
// Content-Security-Policy, so an uploaded page opened inline cannot run script
// as your API.
func ServeFileAs(c *gin.Context, disk Disk, key string, disposition Disposition, filename string) {
	r := c.Request
	obj, err := disk.Stat(r.Context(), key)
	if err != nil {
		serveFileError(c, err)
		return
	}

	contentType := obj.ContentType
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(key))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if disposition == "" {
		disposition = Attachment
	}
	header := c.Writer.Header()
	header.Set("Content-Type", contentType)
	header.Set("Accept-Ranges", "bytes")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Disposition", contentDisposition(disposition, filename))
	if bareType(contentType) != "application/pdf" {
		// A sandboxed document cannot use the browser's PDF viewer.
		header.Set("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; media-src 'self'; style-src 'unsafe-inline'; sandbox")
	}
	lastModified := ""
	if !obj.LastModified.IsZero() {
		lastModified = obj.LastModified.UTC().Format(http.TimeFormat)
		header.Set("Last-Modified", lastModified)
	}
	if notModified(r, obj.LastModified) {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	start, length, partial, ok := parseRange(r.Header.Get("Range"), obj.Size)
	if !ok {
		header.Set("Content-Range", fmt.Sprintf("bytes */%d", obj.Size))
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{
			"error": gin.H{"code": "RANGE_NOT_SATISFIABLE", "message": "The requested range is outside the file"},
		})
		return
	}
	// If-Range names the version a client resumed from. Any other version gets
	// the whole file, or the client would splice two versions together.
	if ifRange := r.Header.Get("If-Range"); partial && ifRange != "" && ifRange != lastModified {
		start, length, partial = 0, obj.Size, false
	}

	status := http.StatusOK
	header.Set("Content-Length", strconv.FormatInt(length, 10))
	if partial {
		status = http.StatusPartialContent
		header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, start+length-1, obj.Size))
	}
	if r.Method == http.MethodHead {
		c.Status(status)
		c.Writer.WriteHeaderNow()
		return
	}

	body, err := openRange(r.Context(), disk, key, start, length, partial)
	if err != nil {
		header.Del("Content-Length")
		header.Del("Content-Range")
		serveFileError(c, err)
		return
	}
	defer body.Close()
	// Written now rather than with the first byte, so an empty file still sends
	// its status.
	c.Status(status)
	c.Writer.WriteHeaderNow()
	if _, err := io.CopyN(c.Writer, body, length); err != nil {
		// The status is already sent, so the client sees a short body.
		_ = c.Error(fmt.Errorf("serving %q: %w", key, err))
	}
}

func contentDisposition(disposition Disposition, filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" || filename == "." || filename == "/" {
		return string(disposition)
	}
	if value := mime.FormatMediaType(string(disposition), map[string]string{"filename": filename}); value != "" {
		return value
	}
	return string(disposition)
}

func notModified(r *http.Request, lastModified time.Time) bool {
	if lastModified.IsZero() || r.Header.Get("If-None-Match") != "" || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
		return false
	}
	since, err := http.ParseTime(r.Header.Get("If-Modified-Since"))
	if err != nil {
		return false
	}
	return !lastModified.Truncate(time.Second).After(since)
}

// parseRange reads a Range header for a file of size bytes. A header that is
// absent, malformed or asks for several ranges is served as the whole file,
// which RFC 9110 allows; a range starting past the end is not satisfiable.
func parseRange(value string, size int64) (start, length int64, partial, ok bool) {
	spec, found := strings.CutPrefix(strings.TrimSpace(value), "bytes=")
	if !found || strings.Contains(spec, ",") {
		return 0, size, false, true
	}
	first, last, found := strings.Cut(strings.TrimSpace(spec), "-")
	if !found {
		return 0, size, false, true
	}
	if first == "" {
		// The last n bytes.
		n, err := strconv.ParseInt(last, 10, 64)
		if err != nil || n < 0 {
			return 0, size, false, true
		}
		if n == 0 || size == 0 {
			return 0, 0, false, false
		}
		if n > size {
			n = size
		}
		return size - n, n, true, true
	}
	from, err := strconv.ParseInt(first, 10, 64)
	if err != nil || from < 0 {
		return 0, size, false, true
	}
	if from >= size {
		return 0, 0, false, false
	}
	to := size - 1
	if last != "" {
		parsed, err := strconv.ParseInt(last, 10, 64)
		if err != nil || parsed < from {
			return 0, size, false, true
		}
		if parsed < to {
			to = parsed
		}
	}
	return from, to - from + 1, true, true
}

type limitedBody struct {
	io.Reader
	io.Closer
}

func openRange(ctx context.Context, disk Disk, key string, start, length int64, partial bool) (io.ReadCloser, error) {
	if !partial {
		return disk.Get(ctx, key)
	}
	if ranged, ok := disk.(RangeGetter); ok {
		return ranged.GetRange(ctx, key, start, length)
	}
	body, err := disk.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if seeker, ok := body.(io.Seeker); ok {
		if _, err := seeker.Seek(start, io.SeekStart); err != nil {
			body.Close()
			return nil, fmt.Errorf("seeking in %q: %w", key, err)
		}
	} else if _, err := io.CopyN(io.Discard, body, start); err != nil {
		body.Close()
		return nil, fmt.Errorf("skipping to the range in %q: %w", key, err)
	}
	return limitedBody{Reader: io.LimitReader(body, length), Closer: body}, nil
}

func serveFileError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrInvalidKey) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "File not found"},
		})
		return
	}
	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not read the file"},
	})
}
`
}

func storageDisksGo() string {
	return `package storage

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"{{MODULE}}/internal/config"
)

// ErrNotConfigured is returned, wrapped, by Open for a bucket driver without
// the settings it needs to connect.
var ErrNotConfigured = errors.New("storage: not configured")

// Open builds the store a driver names: a directory for "local", a bucket for
// "minio", "s3", "r2" and "b2".
//
// AWS S3 needs no endpoint, and with no access key the SDK's own credential
// chain supplies one, so an IAM role on EC2, ECS or Lambda works with only
// S3_BUCKET set. MinIO, R2 and B2 need an endpoint and a key.
func Open(driver string, cfg config.StorageConfig, local LocalConfig) (*Storage, error) {
	switch driver {
	case "local":
		return NewLocal(local)
	case "s3":
		return New(cfg)
	case "minio", "r2", "b2", "":
		if cfg.Endpoint == "" || cfg.AccessKey == "" {
			return nil, fmt.Errorf("%w: STORAGE_DRIVER=%s needs an endpoint and an access key", ErrNotConfigured, driver)
		}
		return New(cfg)
	default:
		return nil, fmt.Errorf("storage: unknown STORAGE_DRIVER %q, use local, minio, s3, r2 or b2", driver)
	}
}

// Disks holds every store the app opened: the default one, and each named disk
// in STORAGE_DISKS. main.go fills it at startup, before anything is served.
//
//	backups := storage.Disks.Get("backups") // nil when STORAGE_DISKS has no backups
var Disks = &Registry{}

// Registry is a set of stores by name.
type Registry struct {
	mu    sync.RWMutex
	def   *Storage
	named map[string]*Storage
}

// SetDefault records the store the app uses when it names none.
func (r *Registry) SetDefault(s *Storage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.def = s
}

// Add records a named store. "" and "default" set the default one.
func (r *Registry) Add(name string, s *Storage) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || name == "default" {
		r.SetDefault(s)
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.named == nil {
		r.named = map[string]*Storage{}
	}
	r.named[name] = s
}

// Default is the store the app uses when it names none, or nil.
func (r *Registry) Default() *Storage {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.def
}

// Get returns the store called name, or nil when none is configured. "" and
// "default" are the default store. A missing name does not fall back to the
// default: code asking for a private bucket should not quietly write to the
// public one.
func (r *Registry) Get(name string) *Storage {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || name == "default" {
		return r.Default()
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.named[name]
}

// Names lists the named stores, sorted.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.named))
	for name := range r.named {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
`
}

func storageStoreTestGo() string {
	return `package storage

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
	h.Set("Content-Disposition", ` + "`" + `form-data; name="file"; filename="` + "`" + `+filename+` + "`" + `"` + "`" + `)
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
	if !regexp.MustCompile(` + "`" + `^uploads/\d{4}/\d{2}/[0-9a-f-]{36}\.png$` + "`" + `).MatchString(key) {
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
			if got := h.Get("Content-Disposition"); got != ` + "`" + `attachment; filename="report final.txt"` + "`" + ` {
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
`
}
