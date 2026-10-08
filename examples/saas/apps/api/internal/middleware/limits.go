package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"saas/apps/api/internal/respond"
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
			respond.Fail(c, respond.CodePayloadTooLarge, fmt.Sprintf("Request body exceeds %dMB limit", limit/(1<<20)))
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
