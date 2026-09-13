package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// oldGzipMiddleware is Gzip and its writer as middleware/logger.go carried them
// before v3.254.0, with the blank line after them. Only this exact text is
// removed.
const oldGzipMiddleware = `// Gzip compresses responses using gzip encoding when the client supports it.
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.Next()
			return
		}
		defer gz.Close()

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipResponseWriter{ResponseWriter: c.Writer, Writer: gz}
		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func (g *gzipResponseWriter) WriteString(s string) (int, error) {
	return g.Writer.Write([]byte(s))
}

`

// repairGzipMiddleware brings a project up to the fix for H16 in the
// contact-app review: the gzip middleware built a compressor per request,
// compressed responses that already were, and swallowed Flush, so the AI stream
// did not stream.
//
// The new middleware lives in middleware/gzip.go, which is framework code. The
// old one is in logger.go, next to code the developer may have changed, so it
// is removed only when it is exactly what Grit wrote, and gzip.go arrives only
// once nothing else in the package declares Gzip.
func repairGzipMiddleware(root string, opts Options) error {
	dir := filepath.Join(opts.APIRoot(root), "internal", "middleware")
	logger := filepath.Join(dir, "logger.go")
	if !fileExists(logger) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if err := repairSourceFile(root, m, logger, repairGzipLoggerSource); err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		name := e.Name()
		if name != "gzip.go" && strings.HasSuffix(name, ".go") && fileContains(filepath.Join(dir, name), "func Gzip(") {
			return nil // still declared elsewhere; logger.go's repair has said so
		}
	}
	for path, content := range map[string]string{
		filepath.Join(dir, "gzip.go"):      middlewareGzipGo(),
		filepath.Join(dir, "gzip_test.go"): middlewareGzipTestGo(),
	} {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func repairGzipLoggerSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func Gzip() gin.HandlerFunc {") {
		return src, nil, nil
	}
	if strings.Count(src, oldGzipMiddleware) != 1 {
		return src, nil, []string{"middleware/logger.go has a Gzip that is not the one Grit wrote: pool the gzip writers, skip compressed types and event streams, and pass Flush through, or streamed responses arrive in one piece"}
	}
	out := strings.Replace(src, oldGzipMiddleware, "", 1)
	if !strings.Contains(out, "gzip.") {
		out = strings.Replace(out, "\t\"compress/gzip\"\n", "", 1)
	}
	return out, []string{"gzip moved to middleware/gzip.go: writers are pooled, compressed types and event streams are left alone, and Flush reaches the client"}, nil
}
