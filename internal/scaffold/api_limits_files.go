package scaffold

// middlewareLimitsGo emits internal/middleware/limits.go: the request body cap
// and the deadlines, chosen per route.
//
// A single global cap and one pair of server timeouts cannot suit an upload, a
// streamed export and a JSON request at once (H14 in the contact-app review): a
// 10 MB cap wrapped the body before the upload handler could set its own, larger
// one, and a 15 second WriteTimeout cut a streamed export off mid-file after it
// had already answered 200.
func middlewareLimitsGo() string {
	return `package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultBodyLimit caps the body of any route not listed as a transfer.
	DefaultBodyLimit int64 = 10 << 20

	// ImportBodyLimit caps a generated resource's CSV import.
	ImportBodyLimit int64 = 100 << 20

	// TransferTimeout is how long a transfer route may spend reading its body
	// or writing its response. The server's own timeouts stay short, because
	// they are what stops a slow client holding a connection open.
	TransferTimeout = 30 * time.Minute
)

// Transfer is a route whose body or response does not fit the defaults.
type Transfer struct {
	// Body is the largest body the route accepts. Zero keeps DefaultBodyLimit.
	Body int64
	// Read extends the deadline for reading the body: uploads, imports.
	Read bool
	// Write extends the deadline for writing the response: exports, downloads,
	// streams, and anything waiting on a slow upstream.
	Write bool
}

// RequestLimits caps the request body and sets the deadlines for each route.
//
// transfers is keyed by method and route pattern, as gin reports it:
// "POST /api/v1/uploads", "GET /api/v1/backups/:id/download". A generated
// resource's GET .../export and POST .../import are transfers without being
// listed. Anything else gets DefaultBodyLimit and the server's timeouts.
func RequestLimits(transfers map[string]Transfer) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, ok := transfers[c.Request.Method+" "+c.FullPath()]
		if !ok {
			t, ok = generatedTransfer(c.Request.Method, c.FullPath())
		}
		limit := DefaultBodyLimit
		if ok && t.Body > 0 {
			limit = t.Body
		}
		if c.Request.ContentLength > limit {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": gin.H{
					"code":    "PAYLOAD_TOO_LARGE",
					"message": fmt.Sprintf("Request body exceeds %dMB limit", limit/(1<<20)),
				},
			})
			return
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		}
		if ok {
			extendDeadlines(c, t)
		}
		c.Next()
	}
}

// generatedTransfer recognises the export and import routes grit generate
// writes for every resource.
func generatedTransfer(method, path string) (Transfer, bool) {
	switch {
	case method == http.MethodGet && strings.HasSuffix(path, "/export"):
		return Transfer{Write: true}, true
	case method == http.MethodPost && strings.HasSuffix(path, "/import"):
		return Transfer{Body: ImportBodyLimit, Read: true, Write: true}, true
	}
	return Transfer{}, false
}

func extendDeadlines(c *gin.Context, t Transfer) {
	rc := http.NewResponseController(c.Writer)
	until := time.Now().Add(TransferTimeout)
	if t.Read {
		if err := rc.SetReadDeadline(until); err != nil && !errors.Is(err, http.ErrNotSupported) {
			log.Printf("limits: extending the read deadline for %s: %v", c.FullPath(), err)
		}
	}
	if t.Write {
		if err := rc.SetWriteDeadline(until); err != nil && !errors.Is(err, http.ErrNotSupported) {
			log.Printf("limits: extending the write deadline for %s: %v", c.FullPath(), err)
		}
	}
}
`
}

// middlewareLimitsTestGo emits internal/middleware/limits_test.go.
func middlewareLimitsTestGo() string {
	return `package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func limitsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestLimits(map[string]Transfer{
		"POST /uploads": {Body: 64 << 20, Read: true, Write: true},
	}))
	read := func(c *gin.Context) {
		n, err := io.Copy(io.Discard, c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.JSON(http.StatusOK, gin.H{"read": n})
	}
	r.POST("/uploads", read)
	r.POST("/notes", read)
	r.POST("/notes/import", read)
	return r
}

func post(r http.Handler, path string, size int, chunked bool) int {
	var body io.Reader = bytes.NewReader(make([]byte, size))
	if chunked {
		body = io.MultiReader(body) // hides the length, so no Content-Length is sent
	}
	req := httptest.NewRequest(http.MethodPost, path, body)
	if chunked {
		req.ContentLength = -1
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

// The body cap a route gets is its own, not the default wrapped around it.
func TestTransferRouteTakesItsOwnBodyLimit(t *testing.T) {
	r := limitsRouter()
	if code := post(r, "/uploads", 20<<20, false); code != http.StatusOK {
		t.Errorf("a 20 MB upload answered %d; the default cap must not apply to a transfer", code)
	}
	if code := post(r, "/uploads", 65<<20, false); code != http.StatusRequestEntityTooLarge {
		t.Errorf("a body over the route's own limit answered %d, want 413", code)
	}
}

func TestOrdinaryRouteKeepsTheDefault(t *testing.T) {
	r := limitsRouter()
	if code := post(r, "/notes", 11<<20, false); code != http.StatusRequestEntityTooLarge {
		t.Errorf("an 11 MB body answered %d, want 413", code)
	}
	// Without a Content-Length the cap is enforced while reading.
	if code := post(r, "/notes", 11<<20, true); code != http.StatusRequestEntityTooLarge {
		t.Errorf("an 11 MB body with no Content-Length answered %d, want 413", code)
	}
}

func TestGeneratedImportIsATransfer(t *testing.T) {
	r := limitsRouter()
	if code := post(r, "/notes/import", 20<<20, false); code != http.StatusOK {
		t.Errorf("a 20 MB import answered %d", code)
	}
}
`
}
