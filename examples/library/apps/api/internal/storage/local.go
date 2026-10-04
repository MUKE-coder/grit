package storage

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
	// localNamedSegment is where the default local disk serves a named local
	// disk's files: /files/_disks/<name>/<key>. No key on a local disk may start
	// with it.
	localNamedSegment = "_disks"
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
	for i, segment := range strings.Split(key, "/") {
		if segment == "" || strings.HasPrefix(segment, localTempPrefix) || (i == 0 && segment == localNamedSegment) {
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
	if err := checkVisibility(key, opts.Visibility); err != nil {
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
//
// A named local disk (STORAGE_DISKS) has no route of its own: its files are
// served here under _disks/<name>/, by that disk and with its own checks.
func (d *LocalDisk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, localFilesPrefix)
	if rest, ok := strings.CutPrefix(key, localNamedSegment+"/"); ok {
		name, namedKey, _ := strings.Cut(rest, "/")
		if name != "" && name != "default" {
			if named := Disks.Get(name); named != nil {
				if other, ok := named.Disk().(*LocalDisk); ok && other != d {
					other.serve(w, r, namedKey)
					return
				}
			}
		}
		writeFileError(w, http.StatusNotFound, "NOT_FOUND", "File not found")
		return
	}
	d.serve(w, r, key)
}

func (d *LocalDisk) serve(w http.ResponseWriter, r *http.Request, key string) {
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
