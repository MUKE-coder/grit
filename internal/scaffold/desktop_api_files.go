package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeDesktopAPIFiles makes the desktop app a HYBRID: a Wails shell that also
// runs a real Gin REST API in-process, backed by SQLite (offline default) or
// Postgres. That gives us three things the Wails bindings alone can't:
//
//	internal/files    — FileRef, the same JSON shape web/mobile store
//	internal/storage  — writes uploads into the OS app-data dir
//	internal/api      — POST /api/uploads, GET /uploads/*, GET /api/health
//
// The router is mounted twice on purpose:
//   - as Wails' AssetServer.Handler, so the webview can `fetch("/api/uploads")`
//     and `<img src="/uploads/x.jpg">` same-origin — no port, no CORS
//   - on 127.0.0.1:<APIPort>, so curl / another client can reach the same API
//
// Templates use "~" where a backtick belongs (Go raw strings can't hold one).
func writeDesktopAPIFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "internal", "files", "file_ref.go"):  desktopFileRefGo(),
		filepath.Join(root, "internal", "storage", "storage.go"): desktopStorageGo(),
		filepath.Join(root, "internal", "api", "router.go"):      desktopAPIRouterGo(),
		filepath.Join(root, "internal", "api", "uploads.go"):     desktopAPIUploadsGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "~", "`")
		content = strings.ReplaceAll(content, "<MODULE>", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// desktopFileRefGo mirrors the API scaffold's FileRef so a desktop model stores
// the same JSON shape as its web/mobile counterpart.
func desktopFileRefGo() string {
	return `package files

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// FileRef is a stored file, persisted as JSON in a single column. Same shape the
// web and mobile apps use, so a desktop record round-trips through the same code.
//
// URL is RELATIVE (e.g. "/uploads/abc.jpg"). The desktop API serves it, and the
// frontend renders it directly — same origin inside the Wails webview.
type FileRef struct {
	URL  string ~json:"url"~
	Key  string ~json:"key,omitempty"~
	Name string ~json:"name,omitempty"~
	MIME string ~json:"mime,omitempty"~
	Size int64  ~json:"size,omitempty"~
}

// Value implements driver.Valuer so GORM can persist the ref as JSON.
func (f FileRef) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan implements sql.Scanner. SQLite hands back a string, Postgres []byte.
func (f *FileRef) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("files: cannot scan FileRef")
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, f)
}

// FileRefs is a slice of FileRef for multi-file (~files~) fields, persisted as a
// JSON array in one column.
type FileRefs []FileRef

// Value implements driver.Valuer.
func (r FileRefs) Value() (driver.Value, error) {
	if r == nil {
		return "[]", nil
	}
	return json.Marshal(r)
}

// Scan implements sql.Scanner.
func (r *FileRefs) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("files: cannot scan FileRefs")
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, r)
}
`
}

// desktopStorageGo saves uploads into the OS app-data directory — not next to
// the binary, which may sit in a read-only Program Files / /Applications path.
func desktopStorageGo() string {
	return `package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"<MODULE>/internal/files"
)

// Storage writes uploaded files to disk under the OS app-data directory.
//
// Why not next to the binary? An installed app often lives somewhere read-only
// (Program Files, /Applications). App-data is writable, per-user, and survives
// upgrades — so a photo attached today is still there after the next release.
type Storage struct {
	dir string
}

// New resolves (and creates) <user-config-dir>/<appName>/uploads.
func New(appName string) (*Storage, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolving app data dir: %w", err)
	}
	dir := filepath.Join(base, appName, "uploads")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", dir, err)
	}
	return &Storage{dir: dir}, nil
}

// Dir is where the files live — the API serves /uploads/* straight from here.
func (s *Storage) Dir() string { return s.dir }

// Save streams an uploaded file to disk and returns the FileRef to persist.
// The stored name is a UUID, so two users uploading "photo.jpg" never collide;
// the original name is kept in the ref for display and downloads.
func (s *Storage) Save(file multipart.File, header *multipart.FileHeader) (*files.FileRef, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	name := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String(), ext)

	dst, err := os.Create(filepath.Join(s.dir, name))
	if err != nil {
		return nil, fmt.Errorf("creating file: %w", err)
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("writing file: %w", err)
	}

	return &files.FileRef{
		URL:  "/uploads/" + name,
		Key:  name,
		Name: header.Filename,
		MIME: header.Header.Get("Content-Type"),
		Size: size,
	}, nil
}

// Delete removes a stored file by its key. Missing files are not an error —
// the point is that it's gone.
func (s *Storage) Delete(key string) error {
	if key == "" {
		return nil
	}
	err := os.Remove(filepath.Join(s.dir, filepath.Base(key)))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
`
}

// desktopAPIRouterGo is the embedded REST API. It is mounted BOTH as Wails'
// AssetServer.Handler (same-origin for the webview) and on a loopback TCP port
// (so curl and other clients can reach it) — the "hybrid" part.
func desktopAPIRouterGo() string {
	return `package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"<MODULE>/internal/storage"
)

// Server holds the dependencies every route needs.
type Server struct {
	DB      *gorm.DB
	Storage *storage.Storage
}

// NewRouter builds the REST API.
//
// It's a plain http.Handler, which is what lets us mount it in two places:
//
//	AssetServer.Handler = router   → the webview fetches /api/... same-origin,
//	                                 and <img src="/uploads/x.jpg"> just works
//	http.ListenAndServe(":port")   → curl, scripts, a second client
//
// Anything the asset server can't resolve as a frontend file falls through to
// this router, so there's no port to discover and no CORS to configure.
func NewRouter(s *Server) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Uploaded files, served straight off disk from the app-data dir.
	r.Static("/uploads", s.Storage.Dir())

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"status": "ok"}})
		})

		api.POST("/uploads", s.CreateUpload)

		// grit:api-routes
	}

	return r
}
`
}

// desktopAPIUploadsGo is the multipart upload endpoint the generated file
// dropzone posts to.
func desktopAPIUploadsGo() string {
	return `package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// maxUploadSize caps a single file at 50MB.
const maxUploadSize = 50 << 20

// allowedMIME is the image allowlist. Widen it if you attach PDFs or video.
var allowedMIME = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// CreateUpload accepts multipart/form-data (field "file"), writes it to the
// app-data dir and returns the FileRef to store on your model.
//
//	POST /api/uploads   ->  { "data": { "url": "/uploads/…", "name": …, … } }
func (s *Server) CreateUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_FILE", "message": "No file provided"},
		})
		return
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "FILE_TOO_LARGE", "message": "File exceeds the 50MB limit"},
		})
		return
	}

	// Sniff the real content type from the first 512 bytes — the client-declared
	// Content-Type is spoofable, and http.DetectContentType reliably identifies
	// the allowed image formats. Validate against the sniffed type.
	sniff := make([]byte, 512)
	n, _ := io.ReadFull(file, sniff)
	if _, serr := file.Seek(0, io.SeekStart); serr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "UPLOAD_FAILED", "message": "Could not read the uploaded file"},
		})
		return
	}
	mime := strings.SplitN(http.DetectContentType(sniff[:n]), ";", 2)[0]
	if !allowedMIME[mime] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_FILE_TYPE", "message": "Only JPEG, PNG, GIF and WebP images are allowed"},
		})
		return
	}

	ref, err := s.Storage.Save(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "UPLOAD_FAILED", "message": "Failed to store the file"},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": ref, "message": "File uploaded"})
}
`
}
