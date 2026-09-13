package scaffold

// middlewareGzipGo emits internal/middleware/gzip.go.
//
// The gzip middleware it replaced (H16 in the contact-app review) built a new
// compressor for every request, compressed responses that already were (xlsx,
// PDF, images, zip), and did not pass Flush through its buffer, so a
// server-sent event stream arrived in one piece when the handler returned.
func middlewareGzipGo() string {
	return `package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// gzipMinLength is the smallest declared body worth compressing: below it the
// gzip header and footer cost more than they save.
const gzipMinLength = 1024

// gzipWriters are reused: a compressor holds hundreds of kilobytes of state, and
// building one per request made every response pay for it.
var gzipWriters = sync.Pool{New: func() any {
	gz, err := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
	if err != nil {
		log.Printf("gzip: building a writer: %v", err)
		return gzip.NewWriter(io.Discard)
	}
	return gz
}}

// Gzip compresses a response when the client accepts gzip and the response is
// worth compressing.
//
// The decision is made at the handler's first write, when its status and
// headers are known. A response is sent as it is when it has no body, is
// already encoded, is a server-sent event stream, declares a length under
// gzipMinLength, or is a type that is already compressed (images, video, PDF,
// zip, xlsx): only text-like types are compressed. Flush reaches the client, so
// a streamed response streams.
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodHead || c.GetHeader("Upgrade") != "" ||
			!acceptsGzip(c.GetHeader("Accept-Encoding")) {
			c.Next()
			return
		}
		w := &gzipResponseWriter{ResponseWriter: c.Writer}
		c.Writer = w
		defer w.finish()
		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	gz      *gzip.Writer
	decided bool
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	g.decide(data)
	if g.gz == nil {
		return g.ResponseWriter.Write(data)
	}
	return g.gz.Write(data)
}

func (g *gzipResponseWriter) WriteString(s string) (int, error) {
	return g.Write([]byte(s))
}

// Flush sends what has been compressed so far, then flushes the connection.
func (g *gzipResponseWriter) Flush() {
	if g.gz != nil {
		if err := g.gz.Flush(); err != nil {
			return
		}
	}
	g.ResponseWriter.Flush()
}

// decide chooses, once, whether this response is compressed.
func (g *gzipResponseWriter) decide(first []byte) {
	if g.decided {
		return
	}
	g.decided = true
	h := g.ResponseWriter.Header()
	if h.Get("Content-Type") == "" {
		// Sniffed from the plain bytes: net/http would otherwise sniff the
		// compressed ones.
		h.Set("Content-Type", http.DetectContentType(first))
	}
	status := g.ResponseWriter.Status()
	if status < http.StatusOK || status == http.StatusNoContent || status == http.StatusNotModified ||
		h.Get("Content-Encoding") != "" || !compressible(h.Get("Content-Type")) {
		return
	}
	if n, err := strconv.Atoi(h.Get("Content-Length")); err == nil && n < gzipMinLength {
		return
	}
	h.Del("Content-Length")
	h.Set("Content-Encoding", "gzip")
	h.Add("Vary", "Accept-Encoding")
	gz, ok := gzipWriters.Get().(*gzip.Writer)
	if !ok {
		return
	}
	gz.Reset(g.ResponseWriter)
	g.gz = gz
}

func (g *gzipResponseWriter) finish() {
	if g.gz == nil {
		return
	}
	if err := g.gz.Close(); err != nil {
		// The client went away mid-response; there is nobody to tell.
		log.Printf("gzip: finishing a response: %v", err)
	}
	g.gz.Reset(io.Discard)
	gzipWriters.Put(g.gz)
	g.gz = nil
}

// acceptsGzip reports whether an Accept-Encoding header allows gzip.
func acceptsGzip(header string) bool {
	for _, part := range strings.Split(header, ",") {
		name, params, _ := strings.Cut(strings.TrimSpace(part), ";")
		name = strings.ToLower(strings.TrimSpace(name))
		if name != "gzip" && name != "*" {
			continue
		}
		// "gzip;q=0" is a refusal, not an offer.
		if value, found := strings.CutPrefix(strings.ReplaceAll(strings.ToLower(params), " ", ""), "q="); found {
			q, err := strconv.ParseFloat(value, 64)
			return err == nil && q > 0
		}
		return true
	}
	return false
}

// compressible reports whether a Content-Type is text-like. Everything else is
// either already compressed or not worth the CPU.
func compressible(contentType string) bool {
	mediaType, _, _ := strings.Cut(contentType, ";")
	mediaType = strings.ToLower(strings.TrimSpace(mediaType))
	switch {
	case mediaType == "text/event-stream":
		return false
	case strings.HasPrefix(mediaType, "text/"),
		strings.HasSuffix(mediaType, "+json"),
		strings.HasSuffix(mediaType, "+xml"):
		return true
	}
	switch mediaType {
	case "application/json", "application/javascript", "application/xml",
		"application/x-ndjson", "application/graphql-response+json", "application/wasm":
		return true
	}
	return false
}
`
}

// middlewareGzipTestGo emits internal/middleware/gzip_test.go.
func middlewareGzipTestGo() string {
	return `package middleware

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
	if perRequest := (after.TotalAlloc - before.TotalAlloc) / n; perRequest > 256<<10 {
		t.Errorf("%d KB allocated per compressed response; the compressor is not being reused", perRequest>>10)
	}
}
`
}
