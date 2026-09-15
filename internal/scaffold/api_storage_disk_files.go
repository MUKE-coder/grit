package scaffold

// The storage Disk interface and its local driver (issue 003, C1).
//
// storageDiskGo is the interface every driver implements, storageLocalGo the
// driver that keeps files in a directory and serves them from the API, and
// storageDiskTestGo one suite run against the local disk always and against a
// bucket when STORAGE_TEST_S3_* names one. storage.go (storageServiceGo) holds
// the S3 driver and the Storage type handlers have always used.

func storageDiskGo() string {
	return `package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

// Disk is a place files are kept: a bucket, or a directory on this machine.
//
// Every driver answers the same way. A missing file is ErrNotFound, a key that
// could name something outside the store is ErrInvalidKey, and Delete is not
// an error for a key that is already gone. Code written against Disk runs
// unchanged on the local disk in development and on a bucket in production.
type Disk interface {
	// Put stores r at key, replacing any file already there.
	Put(ctx context.Context, key string, r io.Reader, opts PutOptions) error
	// Get opens the file at key. The caller closes it.
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Exists reports whether a file is stored at key.
	Exists(ctx context.Context, key string) (bool, error)
	// Stat describes the file at key without reading it.
	Stat(ctx context.Context, key string) (Object, error)
	// Delete removes every key given.
	Delete(ctx context.Context, keys ...string) error
	// Copy stores a copy of the file at from under to.
	Copy(ctx context.Context, from, to string) error
	// Move renames the file at from to to.
	Move(ctx context.Context, from, to string) error
	// List returns every file whose key starts with prefix.
	List(ctx context.Context, prefix string) ([]Object, error)
	// URL is where a browser loads a public file from.
	URL(key string) string
	// TemporaryURL is a link to any file, public or not, that stops working
	// after ttl.
	TemporaryURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

// PutOptions describe a file being stored.
type PutOptions struct {
	// ContentType is recorded with the file where the driver records one.
	// Empty means application/octet-stream.
	ContentType string
}

// Object describes one stored file.
type Object struct {
	Key  string
	Size int64
	// ContentType is empty in a List, where a bucket does not return it.
	ContentType  string
	LastModified time.Time
}

// Presigner is a Disk that can give a browser a URL to PUT one file to
// directly, so the bytes never pass through the API.
type Presigner interface {
	PresignPut(ctx context.Context, key, contentType string, size int64, ttl time.Duration) (string, error)
}

var (
	// ErrNotFound is returned, wrapped, for a key with no file.
	ErrNotFound = errors.New("storage: file not found")
	// ErrInvalidKey is returned, wrapped, for a key that is empty, absolute
	// or climbs out of the store with a ".." segment.
	ErrInvalidKey = errors.New("storage: invalid key")
	// ErrPresignUnsupported is returned by PresignPutURL when the driver takes
	// uploads through the API instead (STORAGE_DRIVER=local).
	ErrPresignUnsupported = errors.New("storage: this driver cannot presign uploads")
)

// checkKey refuses a key that no driver should accept. A bucket would store
// "../x" as an ordinary name, but the same key on the local disk is a path out
// of the root, so it is refused everywhere and code behaves the same on both.
func checkKey(key string) error {
	if key == "" || strings.HasPrefix(key, "/") || strings.ContainsRune(key, 0) {
		return fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	for _, segment := range strings.Split(key, "/") {
		if segment == "." || segment == ".." {
			return fmt.Errorf("%w: %q", ErrInvalidKey, key)
		}
	}
	return nil
}

// escapeKey escapes each segment of a key for a URL path, keeping the slashes.
func escapeKey(key string) string {
	segments := strings.Split(key, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}

// exists is Exists for any driver: Stat, with a missing file as false.
func exists(ctx context.Context, disk Disk, key string) (bool, error) {
	_, err := disk.Stat(ctx, key)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
`
}

func storageLocalGo() string {
	return `package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// LocalFilesRoute serves the files of a local disk. routes.Setup mounts it
// when STORAGE_DRIVER=local.
const LocalFilesRoute = "/files/*key"

const (
	localFilesPrefix = "/files/"
	// localTempPrefix names a write in progress. List skips these, and no key
	// may start with it.
	localTempPrefix = ".grit-tmp-"
)

// LocalConfig configures a LocalDisk.
type LocalConfig struct {
	// Root is the directory files are kept in (STORAGE_LOCAL_ROOT). It is
	// created if it does not exist.
	Root string
	// PublicURL is where LocalFilesRoute is reached, for example
	// http://localhost:8080/files.
	PublicURL string
	// Secret signs temporary URLs (STORAGE_URL_SECRET, or JWT_SECRET).
	Secret string
}

// LocalDisk keeps files in a directory on this machine and serves them from
// the API.
//
// Keys under PublicPrefixes are served to anyone, as a public bucket would
// serve them. Every other key, backups included, is served only through a
// TemporaryURL, whose signature and expiry are checked on each request.
//
// Files live on one machine, so a second replica cannot see them and a
// container replaced on deploy loses them unless Root is on a volume. That is
// why production refuses this driver unless ALLOW_LOCAL_STORAGE_IN_PRODUCTION
// is set.
type LocalDisk struct {
	root      string
	publicURL string
	secret    []byte
}

// NewLocalDisk opens a local disk rooted at cfg.Root, creating the directory.
func NewLocalDisk(cfg LocalConfig) (*LocalDisk, error) {
	if strings.TrimSpace(cfg.Root) == "" {
		return nil, errors.New("storage: the local driver needs a directory, set STORAGE_LOCAL_ROOT")
	}
	if strings.TrimSpace(cfg.PublicURL) == "" {
		return nil, errors.New("storage: the local driver needs the URL its files are served from")
	}
	root, err := filepath.Abs(cfg.Root)
	if err != nil {
		return nil, fmt.Errorf("storage: resolving %q: %w", cfg.Root, err)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("storage: creating %q: %w", root, err)
	}
	return &LocalDisk{
		root:      root,
		publicURL: strings.TrimRight(cfg.PublicURL, "/"),
		secret:    []byte(cfg.Secret),
	}, nil
}

// Root is the absolute directory files are kept in.
func (d *LocalDisk) Root() string { return d.root }

// path maps a key to a file under the root, refusing any key that could land
// anywhere else.
func (d *LocalDisk) path(key string) (string, error) {
	if err := checkKey(key); err != nil {
		return "", err
	}
	// A backslash is a separator on Windows and a colon names a drive or an
	// alternate data stream there, so neither reaches the filesystem.
	if strings.Contains(key, "\\") || (runtime.GOOS == "windows" && strings.Contains(key, ":")) {
		return "", fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	for _, segment := range strings.Split(key, "/") {
		if segment == "" || strings.HasPrefix(segment, localTempPrefix) {
			return "", fmt.Errorf("%w: %q", ErrInvalidKey, key)
		}
	}
	full := filepath.Join(d.root, filepath.FromSlash(key))
	rel, err := filepath.Rel(d.root, full)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	return full, nil
}

// failed wraps a filesystem error, reporting a missing file as ErrNotFound.
func failed(op, key string, err error) error {
	if errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTDIR) {
		return fmt.Errorf("%s %q: %w", op, key, ErrNotFound)
	}
	return fmt.Errorf("%s %q: %w", op, key, err)
}

// Put writes to a temporary file beside the target and renames it into place,
// so a reader never sees half a file and a failed write leaves the old one.
func (d *LocalDisk) Put(ctx context.Context, key string, r io.Reader, opts PutOptions) error {
	full, err := d.path(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("storing %q: %w", key, err)
	}
	if err := writeAtomically(full, r); err != nil {
		return fmt.Errorf("storing %q: %w", key, err)
	}
	return nil
}

func writeAtomically(target string, r io.Reader) (err error) {
	tmp, err := os.CreateTemp(filepath.Dir(target), localTempPrefix+"*")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
		}
	}()
	if _, err = io.Copy(tmp, r); err != nil {
		return err
	}
	if err = tmp.Chmod(0o644); err != nil {
		return err
	}
	if err = tmp.Sync(); err != nil {
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), target)
}

// Get opens the file at key.
func (d *LocalDisk) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	full, err := d.path(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(full)
	if err != nil {
		return nil, failed("reading", key, err)
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, failed("reading", key, err)
	}
	if info.IsDir() {
		_ = f.Close()
		return nil, fmt.Errorf("reading %q: %w", key, ErrNotFound)
	}
	return f, nil
}

// Exists reports whether a file is stored at key.
func (d *LocalDisk) Exists(ctx context.Context, key string) (bool, error) {
	return exists(ctx, d, key)
}

// Stat describes the file at key. The content type comes from the extension,
// or from the first bytes when the extension says nothing.
func (d *LocalDisk) Stat(ctx context.Context, key string) (Object, error) {
	full, err := d.path(key)
	if err != nil {
		return Object{}, err
	}
	info, err := os.Stat(full)
	if err != nil {
		return Object{}, failed("stat", key, err)
	}
	if info.IsDir() {
		return Object{}, fmt.Errorf("stat %q: %w", key, ErrNotFound)
	}
	contentType, err := localContentType(full)
	if err != nil {
		return Object{}, failed("stat", key, err)
	}
	return Object{Key: key, Size: info.Size(), ContentType: contentType, LastModified: info.ModTime()}, nil
}

func localContentType(full string) (string, error) {
	if byExtension := mime.TypeByExtension(filepath.Ext(full)); byExtension != "" {
		return strings.TrimSpace(strings.SplitN(byExtension, ";", 2)[0]), nil
	}
	f, err := os.Open(full)
	if err != nil {
		return "", err
	}
	defer f.Close()
	head := make([]byte, 512)
	n, err := io.ReadFull(f, head)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}
	return strings.SplitN(http.DetectContentType(head[:n]), ";", 2)[0], nil
}

// Delete removes every key given. A key already gone is not an error.
func (d *LocalDisk) Delete(ctx context.Context, keys ...string) error {
	paths := make([]string, 0, len(keys))
	for _, key := range keys {
		full, err := d.path(key)
		if err != nil {
			return err
		}
		paths = append(paths, full)
	}
	for i, full := range paths {
		if err := os.Remove(full); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("deleting %q: %w", keys[i], err)
		}
	}
	return nil
}

// Copy stores a copy of the file at from under to.
func (d *LocalDisk) Copy(ctx context.Context, from, to string) error {
	if _, err := d.path(to); err != nil {
		return err
	}
	src, err := d.Get(ctx, from)
	if err != nil {
		return err
	}
	defer src.Close()
	return d.Put(ctx, to, src, PutOptions{})
}

// Move renames the file at from to to, in one step on the same filesystem.
func (d *LocalDisk) Move(ctx context.Context, from, to string) error {
	fromPath, err := d.path(from)
	if err != nil {
		return err
	}
	toPath, err := d.path(to)
	if err != nil {
		return err
	}
	info, err := os.Stat(fromPath)
	if err != nil {
		return failed("moving", from, err)
	}
	if info.IsDir() {
		return fmt.Errorf("moving %q: %w", from, ErrNotFound)
	}
	if err := os.MkdirAll(filepath.Dir(toPath), 0o755); err != nil {
		return fmt.Errorf("moving %q: %w", from, err)
	}
	if err := os.Rename(fromPath, toPath); err != nil {
		return failed("moving", from, err)
	}
	return nil
}

// List returns every file whose key starts with prefix, in key order.
func (d *LocalDisk) List(ctx context.Context, prefix string) ([]Object, error) {
	start := d.root
	if i := strings.LastIndex(prefix, "/"); i > 0 {
		dir, err := d.path(prefix[:i])
		if err != nil {
			return nil, err
		}
		start = dir
	}
	var objects []Object
	err := filepath.WalkDir(start, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			if p == start && errors.Is(err, fs.ErrNotExist) {
				return fs.SkipAll
			}
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() || strings.HasPrefix(entry.Name(), localTempPrefix) {
			return nil
		}
		rel, err := filepath.Rel(d.root, p)
		if err != nil {
			return err
		}
		key := filepath.ToSlash(rel)
		if !strings.HasPrefix(key, prefix) {
			return nil
		}
		info, err := entry.Info()
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		objects = append(objects, Object{Key: key, Size: info.Size(), LastModified: info.ModTime()})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("listing %q: %w", prefix, err)
	}
	return objects, nil
}

// URL is the public link to key, served by LocalFilesRoute.
func (d *LocalDisk) URL(key string) string {
	return d.publicURL + "/" + escapeKey(key)
}

// TemporaryURL signs a link to key that stops working after ttl.
func (d *LocalDisk) TemporaryURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	if _, err := d.path(key); err != nil {
		return "", err
	}
	if len(d.secret) == 0 {
		return "", errors.New("storage: temporary URLs need a secret, set STORAGE_URL_SECRET or JWT_SECRET")
	}
	if ttl <= 0 {
		return "", fmt.Errorf("storage: a temporary URL needs a lifetime above zero, not %s", ttl)
	}
	return d.signedURL(key, time.Now().Add(ttl)), nil
}

func (d *LocalDisk) signedURL(key string, expires time.Time) string {
	unix := strconv.FormatInt(expires.Unix(), 10)
	return d.URL(key) + "?expires=" + unix + "&signature=" + hex.EncodeToString(d.signature(key, unix))
}

func (d *LocalDisk) signature(key, expires string) []byte {
	mac := hmac.New(sha256.New, d.secret)
	mac.Write([]byte("grit-storage-url\n" + key + "\n" + expires))
	return mac.Sum(nil)
}

// validSignature reports whether a signature is this key's and has not expired.
func (d *LocalDisk) validSignature(key, expires, signature string) bool {
	if len(d.secret) == 0 || expires == "" || signature == "" {
		return false
	}
	unix, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return false
	}
	given, err := hex.DecodeString(signature)
	if err != nil || !hmac.Equal(given, d.signature(key, expires)) {
		return false
	}
	return time.Now().Unix() <= unix
}

// ServeHTTP serves one file under LocalFilesRoute.
//
// The files come from the API's own origin, where its cookies live, so each
// is sent with a sandboxing Content-Security-Policy: an uploaded page or SVG
// opened directly cannot run script as the API.
func (d *LocalDisk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, localFilesPrefix)
	full, err := d.path(key)
	if err != nil {
		writeFileError(w, http.StatusNotFound, "NOT_FOUND", "File not found")
		return
	}
	public := IsPublicKey(key)
	if !public {
		query := r.URL.Query()
		if query.Get("signature") == "" {
			writeFileError(w, http.StatusNotFound, "NOT_FOUND", "File not found")
			return
		}
		if !d.validSignature(key, query.Get("expires"), query.Get("signature")) {
			writeFileError(w, http.StatusForbidden, "FORBIDDEN", "This link has expired or is not valid")
			return
		}
	}
	f, err := os.Open(full)
	if err != nil {
		writeFileError(w, http.StatusNotFound, "NOT_FOUND", "File not found")
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		writeFileError(w, http.StatusNotFound, "NOT_FOUND", "File not found")
		return
	}

	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	// The admin and web apps load these from their own origins.
	header.Set("Cross-Origin-Resource-Policy", "cross-origin")
	if strings.EqualFold(filepath.Ext(full), ".pdf") {
		// A sandboxed document cannot use the browser's PDF viewer.
		header.Del("Content-Security-Policy")
	} else {
		header.Set("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; media-src 'self'; style-src 'unsafe-inline'; sandbox")
	}
	if !public {
		header.Set("Cache-Control", "private, no-store")
	}
	http.ServeContent(w, r, info.Name(), info.ModTime(), f)
}

func writeFileError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	body := map[string]map[string]string{"error": {"code": code, "message": message}}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return
	}
}
`
}

func storageDiskTestGo() string {
	return `package storage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"{{MODULE}}/internal/config"
)

// The same checks run against every driver, so code written against Disk
// behaves the same on the local disk in development and on a bucket in
// production.

func TestLocalDisk(t *testing.T) {
	var disk *LocalDisk
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		disk.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)

	root := filepath.Join(t.TempDir(), "files")
	d, err := NewLocalDisk(LocalConfig{Root: root, PublicURL: server.URL + "/files", Secret: "a-test-secret-long-enough-to-sign-with"})
	if err != nil {
		t.Fatal(err)
	}
	disk = d

	runDiskSuite(t, d, func(key string) string {
		return d.signedURL(key, time.Now().Add(-time.Minute))
	})

	t.Run("a refused key writes nothing outside the root", func(t *testing.T) {
		_ = d.Put(context.Background(), "../escape.txt", strings.NewReader("x"), PutOptions{})
		if _, err := os.Stat(filepath.Join(filepath.Dir(root), "escape.txt")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("a file was written outside the root: %v", err)
		}
	})

	t.Run("a signature for one key does not open another", func(t *testing.T) {
		ctx := context.Background()
		for _, key := range []string{"private/a.txt", "private/b.txt"} {
			if err := d.Put(ctx, key, strings.NewReader(key), PutOptions{}); err != nil {
				t.Fatal(err)
			}
		}
		signed, err := d.TemporaryURL(ctx, "private/a.txt", time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		forged := strings.Replace(signed, "private/a.txt", "private/b.txt", 1)
		if status, _ := fetch(t, forged); status != http.StatusForbidden {
			t.Fatalf("a signature for a.txt opened b.txt: status %d", status)
		}
	})
}

// TestS3Disk runs the same suite against a bucket, when one is named:
// STORAGE_TEST_S3_ENDPOINT, STORAGE_TEST_S3_ACCESS_KEY, STORAGE_TEST_S3_SECRET_KEY
// and STORAGE_TEST_S3_BUCKET. A local MinIO works.
func TestS3Disk(t *testing.T) {
	endpoint := os.Getenv("STORAGE_TEST_S3_ENDPOINT")
	if endpoint == "" {
		t.Skip("set STORAGE_TEST_S3_ENDPOINT, STORAGE_TEST_S3_ACCESS_KEY, STORAGE_TEST_S3_SECRET_KEY and STORAGE_TEST_S3_BUCKET to run against a bucket")
	}
	d, err := NewS3Disk(config.StorageConfig{
		Endpoint:  endpoint,
		AccessKey: os.Getenv("STORAGE_TEST_S3_ACCESS_KEY"),
		SecretKey: os.Getenv("STORAGE_TEST_S3_SECRET_KEY"),
		Bucket:    os.Getenv("STORAGE_TEST_S3_BUCKET"),
		Region:    "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	runDiskSuite(t, d, func(key string) string {
		signed, err := d.TemporaryURL(context.Background(), key, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Second)
		return signed
	})
}

// The methods handlers have always called still work, over any driver.
func TestStorageKeepsItsMethodNames(t *testing.T) {
	ctx := context.Background()
	d, err := NewLocalDisk(LocalConfig{Root: t.TempDir(), PublicURL: "http://localhost:8080/files", Secret: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	s := Wrap(d)

	if err := s.Upload(ctx, "uploads/a.txt", strings.NewReader("hello"), "text/plain"); err != nil {
		t.Fatal(err)
	}
	r, err := s.Download(ctx, "uploads/a.txt")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(r)
	r.Close()
	if err != nil || string(body) != "hello" {
		t.Fatalf("Download = %q, %v", body, err)
	}
	if size, ctype, err := s.Stat(ctx, "uploads/a.txt"); err != nil || size != 5 || ctype != "text/plain" {
		t.Fatalf("Stat = %d, %q, %v", size, ctype, err)
	}
	if got := s.GetURL("uploads/a.txt"); got != "http://localhost:8080/files/uploads/a.txt" {
		t.Fatalf("GetURL = %q", got)
	}
	if _, err := s.GetSignedURL(ctx, "uploads/a.txt", time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PresignPutURL(ctx, "uploads/b.txt", "text/plain", 5); !errors.Is(err, ErrPresignUnsupported) {
		t.Fatalf("PresignPutURL on the local disk = %v, want ErrPresignUnsupported", err)
	}
	if s.FileServer() == nil {
		t.Fatal("the local disk has no file server")
	}
	if err := s.DeleteMany(ctx, []string{"uploads/a.txt", "uploads/missing.txt"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(ctx, "uploads/a.txt"); err != nil {
		t.Fatal(err)
	}
}

func runDiskSuite(t *testing.T, disk Disk, expiredURL func(key string) string) {
	ctx := context.Background()
	run := strconv.FormatInt(time.Now().UnixNano(), 36)
	public := "uploads/disk-suite-" + run + "/"
	private := "private/disk-suite-" + run + "/"
	t.Cleanup(func() {
		for _, prefix := range []string{public, private} {
			objects, err := disk.List(ctx, prefix)
			if err != nil {
				t.Errorf("cleaning up %s: %v", prefix, err)
				continue
			}
			keys := make([]string, 0, len(objects))
			for _, obj := range objects {
				keys = append(keys, obj.Key)
			}
			if len(keys) > 0 {
				if err := disk.Delete(ctx, keys...); err != nil {
					t.Errorf("cleaning up %s: %v", prefix, err)
				}
			}
		}
	})

	put := func(t *testing.T, key, body string) {
		t.Helper()
		if err := disk.Put(ctx, key, strings.NewReader(body), PutOptions{ContentType: "text/plain"}); err != nil {
			t.Fatalf("Put(%q): %v", key, err)
		}
	}
	read := func(t *testing.T, key string) string {
		t.Helper()
		r, err := disk.Get(ctx, key)
		if err != nil {
			t.Fatalf("Get(%q): %v", key, err)
		}
		defer r.Close()
		body, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("reading %q: %v", key, err)
		}
		return string(body)
	}
	present := func(t *testing.T, key string) bool {
		t.Helper()
		ok, err := disk.Exists(ctx, key)
		if err != nil {
			t.Fatalf("Exists(%q): %v", key, err)
		}
		return ok
	}

	t.Run("put, get, exists and stat", func(t *testing.T) {
		key := public + "hello.txt"
		put(t, key, "hello")
		if got := read(t, key); got != "hello" {
			t.Fatalf("Get = %q, want hello", got)
		}
		if !present(t, key) || present(t, public+"missing.txt") {
			t.Fatal("Exists is wrong")
		}
		obj, err := disk.Stat(ctx, key)
		if err != nil {
			t.Fatal(err)
		}
		if obj.Size != 5 || obj.ContentType != "text/plain" || obj.LastModified.IsZero() {
			t.Fatalf("Stat = %+v", obj)
		}
		if _, err := disk.Stat(ctx, public+"missing.txt"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Stat of a missing file = %v, want ErrNotFound", err)
		}
		if _, err := disk.Get(ctx, public+"missing.txt"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Get of a missing file = %v, want ErrNotFound", err)
		}
		put(t, key, "replaced")
		if got := read(t, key); got != "replaced" {
			t.Fatalf("a second Put left %q", got)
		}
	})

	t.Run("copy and move", func(t *testing.T) {
		src := public + "copy/source.txt"
		put(t, src, "copied")
		if err := disk.Copy(ctx, src, public+"copy/target.txt"); err != nil {
			t.Fatal(err)
		}
		if read(t, public+"copy/target.txt") != "copied" || !present(t, src) {
			t.Fatal("Copy did not leave two files")
		}
		if err := disk.Move(ctx, public+"copy/target.txt", private+"moved/file.txt"); err != nil {
			t.Fatal(err)
		}
		if read(t, private+"moved/file.txt") != "copied" || present(t, public+"copy/target.txt") {
			t.Fatal("Move did not rename the file")
		}
		if err := disk.Move(ctx, public+"copy/missing.txt", public+"copy/x.txt"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Move of a missing file = %v, want ErrNotFound", err)
		}
	})

	t.Run("list", func(t *testing.T) {
		dir := public + "list/"
		for _, name := range []string{"a.txt", "b.txt", "nested/c.txt"} {
			put(t, dir+name, name)
		}
		put(t, public+"list-sibling.txt", "not in the directory")
		objects, err := disk.List(ctx, dir)
		if err != nil {
			t.Fatal(err)
		}
		var keys []string
		for _, obj := range objects {
			keys = append(keys, obj.Key)
			if obj.Size == 0 {
				t.Errorf("%s has no size", obj.Key)
			}
		}
		sort.Strings(keys)
		want := []string{dir + "a.txt", dir + "b.txt", dir + "nested/c.txt"}
		if strings.Join(keys, ",") != strings.Join(want, ",") {
			t.Fatalf("List = %v, want %v", keys, want)
		}
		if missing, err := disk.List(ctx, public+"nothing-here/"); err != nil || len(missing) != 0 {
			t.Fatalf("List of an empty prefix = %v, %v", missing, err)
		}
	})

	t.Run("delete many", func(t *testing.T) {
		keys := []string{public + "del/1.txt", public + "del/2.txt", public + "del/3.txt"}
		for _, key := range keys {
			put(t, key, key)
		}
		if err := disk.Delete(ctx, append(keys, public+"del/never-there.txt")...); err != nil {
			t.Fatal(err)
		}
		for _, key := range keys {
			if present(t, key) {
				t.Fatalf("%s survived Delete", key)
			}
		}
	})

	t.Run("URL serves a public file", func(t *testing.T) {
		key := public + "page one.txt"
		put(t, key, "public body")
		link := disk.URL(key)
		if !strings.Contains(link, "page%20one.txt") {
			t.Fatalf("URL = %q, want the space escaped", link)
		}
		if status, body := fetch(t, link); status != http.StatusOK || body != "public body" {
			t.Fatalf("GET %s = %d %q", link, status, body)
		}
	})

	t.Run("temporary URL opens a private file until it expires", func(t *testing.T) {
		key := private + "secret.txt"
		put(t, key, "private body")
		if status, _ := fetch(t, disk.URL(key)); status == http.StatusOK {
			t.Fatal("a private file is served without a signature")
		}
		signed, err := disk.TemporaryURL(ctx, key, time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		if status, body := fetch(t, signed); status != http.StatusOK || body != "private body" {
			t.Fatalf("GET signed URL = %d %q", status, body)
		}
		if status, _ := fetch(t, expiredURL(key)); status != http.StatusForbidden {
			t.Fatalf("GET expired URL = %d, want 403", status)
		}
	})

	t.Run("keys that climb out of the store are refused", func(t *testing.T) {
		for _, key := range []string{"../escape.txt", "uploads/../../escape.txt", "/etc/passwd", ""} {
			if err := disk.Put(ctx, key, strings.NewReader("x"), PutOptions{}); !errors.Is(err, ErrInvalidKey) {
				t.Errorf("Put(%q) = %v, want ErrInvalidKey", key, err)
			}
			if _, err := disk.Get(ctx, key); !errors.Is(err, ErrInvalidKey) {
				t.Errorf("Get(%q) = %v, want ErrInvalidKey", key, err)
			}
			if err := disk.Delete(ctx, key); !errors.Is(err, ErrInvalidKey) {
				t.Errorf("Delete(%q) = %v, want ErrInvalidKey", key, err)
			}
		}
	})
}

func fetch(t *testing.T, link string) (int, string) {
	t.Helper()
	res, err := http.Get(link)
	if err != nil {
		t.Fatalf("GET %s: %v", link, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", link, err)
	}
	return res.StatusCode, string(body)
}
`
}
