package middleware

import (
	"bytes"
	"context"
	"hash/fnv"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"

	"commerce/apps/api/internal/cache"
)

// cacheFlights collapses concurrent misses on one key into one handler run.
// Without it an entry expiring on a busy endpoint sent every request that
// arrived before the first one finished to the database at once, and each of
// them wrote the entry again.
var cacheFlights singleflight.Group

// CacheResponse caches GET request responses in Redis for the given duration.
// Only caches 200 OK responses. Skips caching if no cache service is available.
func CacheResponse(cacheService *cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cacheService == nil || c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Build cache key from URL + query params. FNV-1a rather than SHA-256:
		// cache keys don't need cryptographic strength and FNV is much faster
		// on the hot path of every cacheable request. v2 is the raw layout
		// below; entries in the earlier JSON layout are never read, and expire.
		h := fnv.New64a()
		h.Write([]byte(c.Request.URL.String()))
		key := "http:v2:" + strconv.FormatUint(h.Sum64(), 16)
		ctx := c.Request.Context()

		if cached, ok := loadCachedResponse(ctx, cacheService, key); ok {
			serveCachedResponse(c, cached)
			return
		}

		// One request per key runs the handler; the others wait for its
		// response. Only the one that ran it has written to its client.
		ran := false
		shared, _, _ := cacheFlights.Do(key, func() (interface{}, error) {
			// A flight that ended between the lookup above and this one has
			// stored the entry already.
			if cached, ok := loadCachedResponse(ctx, cacheService, key); ok {
				return cached, nil
			}
			ran = true
			return captureAndStore(ctx, c, cacheService, key, ttl), nil
		})
		if ran {
			return
		}
		if cached, ok := shared.(*cachedResponse); ok && cached != nil {
			serveCachedResponse(c, cached)
			return
		}
		// The handler's answer was not cacheable (an error, a redirect), so this
		// request gets its own.
		c.Next()
	}
}

// captureAndStore runs the rest of the chain, writing to the client as usual,
// and stores a 200 response. It returns the response, or nil when it was not
// one to cache.
func captureAndStore(ctx context.Context, c *gin.Context, cacheService *cache.Cache, key string, ttl time.Duration) *cachedResponse {
	// bytes.Buffer grows in chunks, so a 100 KB response takes a few
	// allocations instead of one per Write.
	writer := &responseCapture{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
	c.Writer = writer
	c.Header("X-Cache", "MISS")

	c.Next()

	if writer.Status() != http.StatusOK || writer.body.Len() == 0 {
		return nil
	}
	// And not if the response says one caller's copy must not be given to
	// another.
	if !shareableResponse(writer.Header()) {
		return nil
	}
	resp := &cachedResponse{
		Status:      http.StatusOK,
		ContentType: writer.Header().Get("Content-Type"),
		Body:        writer.body.Bytes(),
	}
	// A failed write only costs the next request a miss.
	_ = cacheService.Client().Set(ctx, key, resp.encode(), ttl).Err()
	return resp
}

// shareableResponse reports whether one caller's response may be served to
// another.
//
// The cache key is the URL and nothing else. That is exactly right for a
// catalogue, where every caller gets the same answer and one stored copy
// serves all of them. It is exactly wrong for anything that varies by who
// asked, and this middleware is mounted on the public group, where a cart, a
// recently-viewed list or a per-key quota is a reasonable thing to add later.
// Added there without this, the first shopper's cart is served to the second.
//
// The standard HTTP signals rather than a list of paths, so a handler opts out
// by saying what it is and does not have to know this middleware exists.
func shareableResponse(h http.Header) bool {
	for _, directive := range strings.Split(h.Get("Cache-Control"), ",") {
		switch strings.ToLower(strings.TrimSpace(directive)) {
		case "private", "no-store", "no-cache":
			return false
		}
	}

	// A response that sets a cookie is establishing or refreshing something
	// that identifies this caller, and replaying it hands that to the next one.
	if len(h.Values("Set-Cookie")) > 0 {
		return false
	}

	// Vary names what the response depends on. Anything the key does not
	// include cannot be cached under it, and "*" means never.
	for _, field := range strings.Split(h.Get("Vary"), ",") {
		switch http.CanonicalHeaderKey(strings.TrimSpace(field)) {
		case "":
			// No Vary header at all.
		case "Accept-Encoding":
			// Compression happens after this, per response, so a stored body
			// is the uncompressed one either way.
		default:
			return false
		}
	}

	return true
}

func serveCachedResponse(c *gin.Context, cached *cachedResponse) {
	c.Header("X-Cache", "HIT")
	c.Data(cached.Status, cached.ContentType, cached.Body)
	c.Abort()
}

// loadCachedResponse reads an entry. A missing key, a Redis error and an entry
// it cannot parse all mean the request is served uncached.
func loadCachedResponse(ctx context.Context, cacheService *cache.Cache, key string) (*cachedResponse, bool) {
	raw, err := cacheService.Client().Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return decodeCachedResponse(raw)
}

type cachedResponse struct {
	Status      int
	ContentType string
	Body        []byte
}

// encode lays an entry out as one line of status and content type, then the
// body as it was sent: "200 application/json; charset=utf-8\n{...}". JSON with a
// []byte field stored the body as base64, a third larger in Redis, and cost an
// encode on every miss and a decode on every hit.
func (r *cachedResponse) encode() []byte {
	head := strconv.Itoa(r.Status) + " " + r.ContentType + "\n"
	out := make([]byte, 0, len(head)+len(r.Body))
	out = append(out, head...)
	return append(out, r.Body...)
}

func decodeCachedResponse(raw []byte) (*cachedResponse, bool) {
	nl := bytes.IndexByte(raw, '\n')
	if nl < 0 {
		return nil, false
	}
	sp := bytes.IndexByte(raw[:nl], ' ')
	if sp < 0 {
		return nil, false
	}
	status, err := strconv.Atoi(string(raw[:sp]))
	if err != nil {
		return nil, false
	}
	return &cachedResponse{Status: status, ContentType: string(raw[sp+1 : nl]), Body: raw[nl+1:]}, true
}

type responseCapture struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString is captured too: gin's ResponseWriter has it, so a handler that
// writes a string would otherwise bypass Write and reach the client uncached.
func (w *responseCapture) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
