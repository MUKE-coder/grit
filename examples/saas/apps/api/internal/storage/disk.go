package storage

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
	// Visibility says who may read the file, and empty follows the key.
	//
	// Every driver decides visibility by key prefix: a key under PublicPrefixes
	// is readable by anyone, every other key only through TemporaryURL. On S3
	// and MinIO the bucket policy says so, on the local disk the file route
	// does. Naming a visibility is a check, not a switch: Put refuses a key whose
	// prefix disagrees with ErrVisibilityMismatch, so a file you meant to keep
	// private never lands where anyone can read it.
	//
	// Cloudflare R2 and Backblaze B2 have no bucket policies. There the bucket
	// is public or it is not, whatever the key, so keep private files in a
	// private bucket (a named disk) and serve them with TemporaryURL.
	Visibility Visibility
}

// Visibility is who may read a stored file.
type Visibility string

const (
	// VisibilityPublic is a file anyone with its URL may read.
	VisibilityPublic Visibility = "public"
	// VisibilityPrivate is a file read only through a TemporaryURL, or through
	// the API with ServeFile.
	VisibilityPrivate Visibility = "private"
)

// VisibilityOf reports who may read the file at key.
func VisibilityOf(key string) Visibility {
	if IsPublicKey(key) {
		return VisibilityPublic
	}
	return VisibilityPrivate
}

// checkVisibility refuses a Put whose visibility disagrees with its key.
func checkVisibility(key string, visibility Visibility) error {
	switch visibility {
	case "":
		return nil
	case VisibilityPublic, VisibilityPrivate:
		if actual := VisibilityOf(key); actual != visibility {
			return fmt.Errorf("%w: %q is %s, not %s, because storage.PublicPrefixes (STORAGE_PUBLIC_PREFIXES) decides", ErrVisibilityMismatch, key, actual, visibility)
		}
		return nil
	default:
		return fmt.Errorf("storage: unknown visibility %q, use public or private", visibility)
	}
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
	// ErrVisibilityMismatch is returned, wrapped, by Put for a visibility the
	// key's prefix does not give.
	ErrVisibilityMismatch = errors.New("storage: the key does not have that visibility")
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
