package middleware

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
