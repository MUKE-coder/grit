package middleware

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
