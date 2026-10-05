package storage

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

	"commerce/apps/api/internal/respond"
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
		respond.Fail(c, respond.CodeRangeNotSatisfiable, "The requested range is outside the file")
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
		respond.Fail(c, respond.CodeNotFound, "File not found")
		return
	}
	_ = c.Error(err)
	respond.Fail(c, respond.CodeInternalError, "Could not read the file")
}
