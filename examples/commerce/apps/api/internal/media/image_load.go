package media

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"commerce/apps/api/internal/safefetch"
)

// Disk is where FromDisk reads and Store writes: pass a storage.Disk, from
// storage.Storage's Disk().
//
// It is not storage.Disk itself because the storage package imports media for
// its image helpers, and Go refuses an import cycle. So media names the
// methods it can name, which storage.Disk satisfies at compile time, and calls
// Put, whose options type lives in storage, through reflection in putFile.
type Disk interface {
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	URL(key string) string
}

// Visibility says who may read a stored image.
type Visibility string

const (
	// PublicFile is readable by anyone with the URL: an avatar, a product
	// photo.
	PublicFile Visibility = "public"
	// PrivateFile is read only through a temporary URL.
	PrivateFile Visibility = "private"
)

// Key prefixes for each visibility. The public one matches the first of
// storage.PublicPrefixes, which is what the local driver serves without a
// signature and what the bucket policy opens to anonymous reads; everything
// outside those prefixes is private on every driver.
const (
	publicKeyPrefix  = "uploads/"
	privateKeyPrefix = "private/"
)

// FromDisk opens the image stored at key, reading at most the byte limit.
func FromDisk(ctx context.Context, disk Disk, key string, opts ...Option) (*Image, error) {
	if disk == nil {
		return nil, errors.New("media: FromDisk needs a disk")
	}
	o := collectOptions(opts)
	body, err := disk.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("media: opening %q: %w", key, err)
	}
	data, err := readLimited(body, o.maxBytes)
	if closeErr := body.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf("media: closing %q: %w", key, closeErr)
	}
	if err != nil {
		return nil, err
	}
	return decodeBytes(data, o)
}

// FromURL downloads and opens an image from a URL somebody else chose.
//
// The request goes through safefetch, so a loopback, private or cloud metadata
// address is refused, including one a public hostname resolves to or
// redirects to. The body is read up to the byte limit and the whole fetch is
// bounded by the timeout (WithMaxBytes, WithTimeout).
func FromURL(ctx context.Context, rawURL string, opts ...Option) (*Image, error) {
	o := collectOptions(opts)
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	resp, err := safefetch.Get(ctx, rawURL)
	if err != nil {
		return nil, fmt.Errorf("media: fetching image: %w", err)
	}
	data, err := readResponse(resp, o.maxBytes)
	if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf("media: closing the response: %w", closeErr)
	}
	if err != nil {
		return nil, err
	}
	return decodeBytes(data, o)
}

func readResponse(resp *http.Response, limit int64) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("media: fetching image: the server answered %s", resp.Status)
	}
	if resp.ContentLength > limit {
		return nil, fmt.Errorf("%w: the server says %d bytes", ErrTooLarge, resp.ContentLength)
	}
	return readLimited(resp.Body, limit)
}

// Store encodes the image and writes it to disk under dir, returning the key.
//
//	key, err := img.Orient().Cover(150, 150).ToWebP().
//	    Store(ctx, svc.Storage.Disk(), "avatars", media.PublicFile)
//	url := svc.Storage.Disk().URL(key)
//
// A PublicFile lands under uploads/<dir>/, a PrivateFile under private/<dir>/,
// named with 24 random hex characters and the extension of what was encoded.
func (i *Image) Store(ctx context.Context, disk Disk, dir string, visibility Visibility) (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if disk == nil {
		return "", errors.New("media: Store needs a disk")
	}
	if ctx == nil {
		// Refused rather than replaced with context.Background(): a stand-in
		// context outlives the request it was meant to belong to.
		return "", errors.New("media: Store needs a context: the request's, or context.Background() outside one")
	}
	var prefix string
	switch visibility {
	case PublicFile:
		prefix = publicKeyPrefix
	case PrivateFile:
		prefix = privateKeyPrefix
	default:
		return "", fmt.Errorf("media: visibility must be media.PublicFile or media.PrivateFile, got %q", visibility)
	}
	dir, err := cleanDir(dir)
	if err != nil {
		return "", err
	}
	data, err := i.Encode()
	if err != nil {
		return "", err
	}
	name := make([]byte, 12)
	if _, err := rand.Read(name); err != nil {
		return "", fmt.Errorf("media: naming the file: %w", err)
	}
	format := i.outputFormat()
	key := prefix + dir + hex.EncodeToString(name) + format.ext()
	if err := putFile(ctx, disk, key, data, format.mime(), visibility); err != nil {
		return "", err
	}
	return key, nil
}

// cleanDir returns dir with one trailing slash, or "" for none, refusing a
// directory that could climb out of its prefix.
func cleanDir(dir string) (string, error) {
	dir = strings.Trim(dir, "/")
	if dir == "" {
		return "", nil
	}
	for _, segment := range strings.Split(dir, "/") {
		if segment == "" || segment == "." || segment == ".." || strings.ContainsAny(segment, "\\\x00") {
			return "", fmt.Errorf("media: %q is not a directory Store can write to", dir)
		}
	}
	return dir + "/", nil
}

var (
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	readerType  = reflect.TypeOf((*io.Reader)(nil)).Elem()
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
)

// putFile calls disk.Put(ctx, key, body, storage.PutOptions{ContentType: ...}).
//
// Reflection, because media cannot import storage (see Disk). The shape is
// checked before the call, so a disk without that method is a clear error
// rather than a panic. If the options gain a Visibility string field, it is
// set too.
func putFile(ctx context.Context, disk Disk, key string, data []byte, contentType string, visibility Visibility) error {
	method := reflect.ValueOf(disk).MethodByName("Put")
	if !method.IsValid() {
		return fmt.Errorf("media: %T has no Put method, pass a storage.Disk", disk)
	}
	t := method.Type()
	if t.NumIn() != 4 || t.NumOut() != 1 || t.In(0) != contextType || t.In(1).Kind() != reflect.String ||
		t.In(2) != readerType || t.In(3).Kind() != reflect.Struct || t.Out(0) != errorType {
		return fmt.Errorf("media: %T.Put is not Put(ctx, key, reader, options) error, pass a storage.Disk", disk)
	}
	opts := reflect.New(t.In(3)).Elem()
	if field := opts.FieldByName("ContentType"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(contentType)
	}
	if field := opts.FieldByName("Visibility"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(string(visibility))
	}
	out := method.Call([]reflect.Value{
		reflect.ValueOf(&ctx).Elem(),
		reflect.ValueOf(key).Convert(t.In(1)),
		reflect.ValueOf(io.Reader(bytes.NewReader(data))),
		opts,
	})
	if err, ok := out[0].Interface().(error); ok && err != nil {
		return fmt.Errorf("media: storing %q: %w", key, err)
	}
	return nil
}
