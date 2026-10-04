package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
)

type SecurityHandler struct {
	Bridge *services.SecObsBridge
}

// upstream runs one call to Sentinel or Pulse and records the ones that fail.
//
// Every call on these two dashboards used to throw its error away with a blank
// assignment. With the upstream down, each struct kept its zero value and the
// page rendered a 200 full of noughts: no banned IPs, no threats, no errors and
// no latency. That reads as "you are safe and fast", which is the opposite of
// what had happened.
//
// The names collected here go back in the response as "degraded", so the page
// can say which panels are guesses.
type upstream struct {
	degraded []string
}

func (u *upstream) get(name string, call func() error) {
	if err := call(); err != nil {
		u.degraded = append(u.degraded, name)
		log.Printf("dashboard: %s is unavailable: %v", name, err)
	}
}

// names returns the failures, never nil, so the JSON field is [] and not null.
func (u *upstream) names() []string {
	if u.degraded == nil {
		return []string{}
	}
	return u.degraded
}

// allFailed reports whether not one of total calls came back.
func (u *upstream) allFailed(total int) bool { return len(u.degraded) == total }

// sentinelBlocked is /sentinel/api/ip/blocked: the currently banned IPs.
type sentinelBlocked struct {
	Data []struct {
		IP        string     `json:"ip"`
		Reason    string     `json:"reason"`
		BlockedAt time.Time  `json:"blocked_at"`
		ExpiresAt *time.Time `json:"expires_at"`
	} `json:"data"`
}

// sentinelStats is /sentinel/api/analytics/summary: ThreatStats for a window.
type sentinelStats struct {
	Data struct {
		TotalThreats int64 `json:"total_threats"`
		BlockedCount int64 `json:"blocked_count"`
	} `json:"data"`
}

// sentinelThreats is /sentinel/api/threats: the most recent ThreatEvents.
type sentinelThreats struct {
	Data []struct {
		ID          string    `json:"id"`
		IP          string    `json:"ip"`
		Path        string    `json:"path"`
		ThreatTypes []string  `json:"threat_types"`
		Severity    string    `json:"severity"`
		Timestamp   time.Time `json:"timestamp"`
	} `json:"data"`
}

// Summary returns the flat security envelope the React dashboard reads.
// On a fresh app with no traffic, expect zeros and empty arrays — that's
// the truth, not a bug.
func (h *SecurityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodeSentinelOff, "Sentinel is not enabled")
		return
	}

	ctx := c.Request.Context()
	var up upstream

	get := func(name, path string, out interface{}) {
		up.get(name, func() error { return h.Bridge.SentinelGet(ctx, path, out) })
	}

	var blockedResp sentinelBlocked
	get("blocked_ips", "/sentinel/api/ip/blocked", &blockedResp)

	// blocked_count over the 24h window is the closest analogue Sentinel has
	// to "auto-bans in the last 24h".
	var statsResp sentinelStats
	get("threat_stats", "/sentinel/api/analytics/summary?window=24h", &statsResp)

	var threatsResp sentinelThreats
	get("recent_threats", "/sentinel/api/threats?limit=10", &threatsResp)

	// Nothing answered. Zeros here would be a claim that no IP is banned and
	// nothing has attacked this app, which is a worse answer than an error.
	if up.allFailed(3) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "SENTINEL_UNAVAILABLE",
				"message": "Sentinel is not answering",
				"details": gin.H{"degraded": up.names()},
			},
		})
		return
	}

	activeBans := make([]gin.H, 0, len(blockedResp.Data))
	for _, b := range blockedResp.Data {
		var expires string
		if b.ExpiresAt != nil {
			expires = b.ExpiresAt.Format(time.RFC3339)
		}
		activeBans = append(activeBans, gin.H{
			"ip":         b.IP,
			"reason":     b.Reason,
			"expires_at": expires,
			"level":      1, // Sentinel's BlockedIP doesn't carry an escalation level
		})
	}

	recentThreats := make([]gin.H, 0, len(threatsResp.Data))
	for _, t := range threatsResp.Data {
		threatType := ""
		if len(t.ThreatTypes) > 0 {
			threatType = t.ThreatTypes[0]
		}
		recentThreats = append(recentThreats, gin.H{
			"id":          t.ID,
			"type":        threatType,
			"ip":          t.IP,
			"description": t.Path,
			"created_at":  t.Timestamp.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"banned_ips_now":         len(blockedResp.Data),
		"auto_bans_24h":          statsResp.Data.BlockedCount,
		"rate_limited_last_hour": 0, // Not directly exposed by Sentinel; populated by rate-limit logs in a future release
		"active_bans":            activeBans,
		"rate_limit_hits_5min":   []gin.H{},
		"recent_threats":         recentThreats,
		// The panels whose numbers are missing rather than zero.
		"degraded": up.names(),
	})
}
