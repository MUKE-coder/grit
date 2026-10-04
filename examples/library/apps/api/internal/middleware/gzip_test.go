package middleware

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func gzipRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Gzip())
	big := strings.Repeat("grit ", 2000)
	r.GET("/json", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"text": big}) })
	r.GET("/zip", func(c *gin.Context) { c.Data(http.StatusOK, "application/zip", []byte(big)) })
	r.GET("/empty", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	r.GET("/sized", func(c *gin.Context) {
		c.DataFromReader(http.StatusOK, int64(len(big)), "text/csv", strings.NewReader(big), nil)
	})
	return r
}

func get(r http.Handler, path, encoding string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if encoding != "" {
		req.Header.Set("Accept-Encoding", encoding)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGzipCompressesText(t *testing.T) {
	w := get(gzipRouter(), "/json", "gzip, deflate, br")
	if w.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("JSON was not compressed: %v", w.Header())
	}
	zr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(zr)
	if err != nil || !bytes.Contains(body, []byte("grit grit")) {
		t.Fatalf("the body does not decompress: %v", err)
	}
}

func TestGzipLeavesCompressedTypesAlone(t *testing.T) {
	if w := get(gzipRouter(), "/zip", "gzip"); w.Header().Get("Content-Encoding") != "" {
		t.Error("a zip was compressed again")
	}
}

func TestGzipRespectsTheClient(t *testing.T) {
	for _, enc := range []string{"", "br", "gzip;q=0", "identity"} {
		if w := get(gzipRouter(), "/json", enc); w.Header().Get("Content-Encoding") != "" {
			t.Errorf("Accept-Encoding %q got a gzip response", enc)
		}
	}
	if w := get(gzipRouter(), "/json", "gzip;q=0.5"); w.Header().Get("Content-Encoding") != "gzip" {
		t.Error("gzip;q=0.5 did not get a gzip response")
	}
}

func TestGzipNoBodyStaysEmpty(t *testing.T) {
	w := get(gzipRouter(), "/empty", "gzip")
	if w.Body.Len() != 0 || w.Header().Get("Content-Encoding") != "" {
		t.Errorf("a 204 got %d body bytes and Content-Encoding %q", w.Body.Len(), w.Header().Get("Content-Encoding"))
	}
}

// A handler that declares a length must not send a compressed body under it.
func TestGzipDropsADeclaredLength(t *testing.T) {
	w := get(gzipRouter(), "/sized", "gzip")
	if cl := w.Header().Get("Content-Length"); cl != "" && w.Header().Get("Content-Encoding") == "gzip" {
		t.Errorf("a gzip body went out under the uncompressed Content-Length %s", cl)
	}
}

// Each flushed event reaches the client while the handler is still running.
func TestGzipStreamsServerSentEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Gzip())
	release := make(chan struct{})
	r.GET("/events", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.SSEvent("message", "first")
		c.Writer.Flush()
		<-release
		c.SSEvent("message", "second")
	})
	srv := httptest.NewServer(r)
	defer srv.Close()
	defer close(release)

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/events", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := (&http.Transport{DisableCompression: true}).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	got := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(resp.Body).ReadString('\n')
		got <- line
	}()
	select {
	case line := <-got:
		if !strings.Contains(line, "event:message") {
			t.Errorf("first line %q; the stream was compressed or mangled", line)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the first event did not arrive before the handler finished: Flush is not reaching the client")
	}
}

// Compressors are pooled: a thousand compressed responses must not each build
// one.
//
// An unpooled compressor costs about 1.2 MB. The limit sits at half that, not
// near zero, because go test -race makes sync.Pool drop a quarter of what is
// put back, on purpose: pooled responses then average about 300 KB, and a
// tighter limit failed every project's CI, which runs with -race.
func TestGzipReusesCompressors(t *testing.T) {
	r := gzipRouter()
	get(r, "/json", "gzip")
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	const n = 1000
	for i := 0; i < n; i++ {
		get(r, "/json", "gzip")
	}
	runtime.ReadMemStats(&after)
	if perRequest := (after.TotalAlloc - before.TotalAlloc) / n; perRequest > 640<<10 {
		t.Errorf("%d KB allocated per compressed response; the compressor is not being reused", perRequest>>10)
	}
}
