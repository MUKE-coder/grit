package services

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// v3.31.49 -- ResolveClientIP returns the IP we should log for an
// activity event. It honours the browser-supplied X-Public-IP-Hint
// header, but only when the TCP peer is loopback -- so production
// traffic from real proxies (which sets X-Forwarded-For for
// gin.Context.ClientIP to consume) keeps using the trusted path
// and can't be spoofed by a client header.
//
// The hint exists for one reason: in local dev the operator runs
// admin + API on the same machine, so c.ClientIP() returns "::1"
// for every action. The activity feed showing "::1" on every row
// is technically correct but useless. The admin client looks up
// the operator's public IP once per session (via api.ipify.org)
// and attaches it as X-Public-IP-Hint; this function reads it.
func ResolveClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if isLoopback(ip) {
		if hint := strings.TrimSpace(c.GetHeader("X-Public-IP-Hint")); hint != "" {
			// Cap at 64 chars matching the column width so a
			// malformed header can't bloat the audit row.
			if len(hint) > 64 {
				hint = hint[:64]
			}
			return hint
		}
	}
	return ip
}

// isLoopback covers the two loopback forms gin.Context.ClientIP can
// return on a local-only deployment: "::1" (IPv6) and "127.0.0.1"
// (IPv4). "0.0.0.0" shouldn't appear as a client IP but we include
// it for symmetry with the frontend prettyIP helper.
func isLoopback(ip string) bool {
	return ip == "::1" || ip == "127.0.0.1" || ip == "0.0.0.0"
}
