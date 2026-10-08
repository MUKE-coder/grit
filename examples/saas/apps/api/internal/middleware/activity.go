package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/audit"
	"saas/apps/api/internal/models"
)

// ActivityLogger records every successful authenticated mutation
// (POST/PUT/PATCH/DELETE) into models.ActivityLog. Skips:
//   - safe methods (GET/HEAD/OPTIONS)
//   - non-2xx responses (errors aren't audit-relevant)
//   - unauthenticated requests (no user_id ⇒ nothing to attribute)
//
// The payload digest is a SHA-256 hash of the request body — enough to
// prove "this exact payload was sent" without persisting plain-text
// passwords / secrets / PII. Buffered in memory, so MaxBodySize earlier
// in the chain still bounds it.
//
// Insert is fire-and-forget via a bounded channel + single writer
// goroutine. The single-writer design eliminates lock contention on
// the hash chain — only one goroutine ever appends — and the bounded
// channel caps memory + goroutine count under traffic spikes.
func ActivityLogger(db *gorm.DB) gin.HandlerFunc {
	// One chain writer per process, shared with the security-event log.
	audit.Start(db)
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			recordRead(c)
			return
		}

		// Capture the body so we can hash it after the handler runs.
		// gin reads from c.Request.Body, so we tee it through a
		// bytes.Buffer and put a fresh ReadCloser back.
		//
		// Skip multipart/form-data (file uploads): buffering the whole file
		// into memory is pointless for an audit digest, and re-reading it here
		// can leave the handler's ParseMultipartForm with nothing to parse
		// ("No file provided"). Uploads are logged by path/actor, not payload.
		var bodyBytes []byte
		if c.Request.Body != nil &&
			!strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		started := time.Now()
		c.Next()

		// Only log successful mutations — failed ones can be diagnosed
		// from request logs without polluting the audit trail.
		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
			return
		}

		userID, _ := c.Get("user_id")
		uid, _ := userID.(string)
		if uid == "" {
			return // unauthenticated — nothing to audit
		}

		entry := models.ActivityLog{
			UserID:        uid,
			Method:        c.Request.Method,
			Path:          c.FullPath(),
			Status:        c.Writer.Status(),
			PayloadDigest: digestBody(bodyBytes),
			IPAddress:     resolveClientIP(c),
			UserAgent:     c.Request.UserAgent(),
			DurationMS:    time.Since(started).Milliseconds(),
			// Which record the write touched, so the log answers "who changed
			// this record" as well as "who read it". Empty on routes without one.
			ResourceIDs: c.Param("id"),
		}
		// Non-blocking. The writer stamps created_at and chains the entry;
		// with a full backlog it is dropped rather than stalling the request.
		audit.Enqueue(entry)
	}
}

// recordRead runs a read and, when the handler marked what it served, records it.
//
// Only marked reads: handlers generated with --audit-reads mark the rows they
// returned, and nothing else is recorded, because reads are most of all
// traffic and logging every page load would bury the writes.
func recordRead(c *gin.Context) {
	started := time.Now()
	c.Next()

	mark, ok := audit.ReadMarkOf(c)
	if !ok || c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	audit.Enqueue(models.ActivityLog{
		UserID: uid,
		Method: c.Request.Method,
		Path:   c.FullPath(),
		Status: c.Writer.Status(),
		// The query can hold what was searched for, a patient's name typed
		// into a search box, so only its digest is kept.
		PayloadDigest: digestBody([]byte(c.Request.URL.RawQuery)),
		IPAddress:     resolveClientIP(c),
		UserAgent:     c.Request.UserAgent(),
		DurationMS:    time.Since(started).Milliseconds(),
		Resource:      mark.Resource,
		ResourceIDs:   strings.Join(mark.IDs, ","),
		RecordCount:   mark.Count,
	})
}

// v3.31.49 -- mirror of services.ResolveClientIP. Inlined here
// (rather than imported) because middleware is a leaf dep that the
// services package itself relies on through the request chain;
// duplicating ten lines avoids the cycle and keeps the audit path
// allocation-free.
func resolveClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" || ip == "127.0.0.1" || ip == "0.0.0.0" {
		if hint := strings.TrimSpace(c.GetHeader("X-Public-IP-Hint")); hint != "" {
			if len(hint) > 64 {
				hint = hint[:64]
			}
			return hint
		}
	}
	return ip
}

// AuditDroppedCount returns the number of audit entries dropped because the
// writer's backlog was full. Read it from a /healthz or admin endpoint to spot
// sustained back-pressure.
func AuditDroppedCount() uint64 {
	return audit.Dropped()
}

func digestBody(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
