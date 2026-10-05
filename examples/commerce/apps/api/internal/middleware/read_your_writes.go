package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Read-your-own-writes, across replicas, without a shared store.
//
// The bug this exists for: somebody posts, the feed loads from a replica that
// has not received the post yet, and their own post is missing. They post it
// again. It is the most reported bug of Stage 6 and the least obvious, because
// it only happens to the person who wrote and only for a second.
//
// The usual fix keeps "user X wrote at T" in Redis, which means a Redis lookup
// on every read and a Redis outage taking the read path with it. A cookie
// carries the same fact on the one request that needs it, costs nothing, and
// cannot be stale: it travels with the person who wrote.
//
// It is not a security boundary and does not need to be. The worst a forged
// cookie achieves is reading from the primary, which is the correct answer
// served by a more expensive machine.
const (
	readYourWritesCookie = "grit_rw"
	// Long enough to cover normal replication, short enough that a user who
	// writes once does not pin every later read to the primary. Replication
	// under a second is typical; five is the number to raise if your lag
	// monitor says otherwise.
	readYourWritesWindow = 5 * time.Second
)

// ReadYourWrites marks a person who has just written, and flags the requests
// that follow so reads can be pinned to the primary.
//
// Mount it once, globally. It does no work on a project with no replicas: the
// flag is read by database.ForRequest, which returns the same handle either way
// when nothing is configured.
func ReadYourWrites() gin.HandlerFunc {
	return func(c *gin.Context) {
		if raw, err := c.Cookie(readYourWritesCookie); err == nil {
			if ms, err := strconv.ParseInt(raw, 10, 64); err == nil {
				if time.Since(time.UnixMilli(ms)) < readYourWritesWindow {
					c.Set("read_your_writes", true)
				}
			}
		}

		c.Next()

		// Only a write that worked. A failed POST changed nothing, so pinning
		// later reads to the primary buys nothing and costs the primary.
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			if c.Writer.Status() < 300 {
				secure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie(
					readYourWritesCookie,
					strconv.FormatInt(time.Now().UnixMilli(), 10),
					int(readYourWritesWindow.Seconds()),
					"/", "", secure, true,
				)
			}
		}
	}
}
