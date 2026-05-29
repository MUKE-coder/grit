package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// SecObsPoller turns Sentinel + Pulse findings into in-app notifications.
// Started from main.go alongside the cron scheduler when both Sentinel
// and Pulse are enabled. Runs once a minute on its own goroutine; the
// dedup unique index means repeated firings bump Count + UpdatedAt
// instead of spamming the bell.
type SecObsPoller struct {
	DB     *gorm.DB
	Bridge *SecObsBridge
	stop   chan struct{}
}

func NewSecObsPoller(db *gorm.DB, bridge *SecObsBridge) *SecObsPoller {
	return &SecObsPoller{DB: db, Bridge: bridge, stop: make(chan struct{})}
}

func (p *SecObsPoller) Start() {
	if p.Bridge == nil || p.DB == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		p.tick()
		for {
			select {
			case <-p.stop:
				return
			case <-ticker.C:
				p.tick()
			}
		}
	}()
}

func (p *SecObsPoller) Stop() {
	select {
	case <-p.stop:
		// already closed
	default:
		close(p.stop)
	}
}

func (p *SecObsPoller) tick() {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	// — Sentinel: high-CVSS unresolved threats
	var threats struct {
		Data []struct {
			ID       string  `json:"id"`
			Type     string  `json:"type"`
			Severity string  `json:"severity"`
			CVSS     float64 `json:"cvss"`
			SourceIP string  `json:"source_ip"`
			Route    string  `json:"route"`
		} `json:"data"`
	}
	if err := p.Bridge.SentinelGet(ctx, "/sentinel/api/threats?limit=20&resolved=false", &threats); err == nil {
		for _, t := range threats.Data {
			if t.CVSS < 7.0 && t.Severity != "critical" && t.Severity != "high" {
				continue
			}
			sev := t.Severity
			if sev == "" {
				if t.CVSS >= 9.0 {
					sev = "critical"
				} else {
					sev = "high"
				}
			}
			p.upsert(&models.Notification{
				Source:   "sentinel",
				Severity: sev,
				Title:    fmt.Sprintf("Threat — %s (CVSS %.1f)", t.Type, t.CVSS),
				Body:     fmt.Sprintf("Source %s · route %s", t.SourceIP, t.Route),
				Link:     "/sentinel/ui/threats?focus=" + t.ID,
				Dedup:    dedup("sentinel", "threat", t.ID),
			})
		}
	}

	// — Pulse: firing alerts
	var alerts struct {
		Data []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Severity string `json:"severity"`
			Message  string `json:"message"`
		} `json:"data"`
	}
	if err := p.Bridge.PulseGet(ctx, "/pulse/api/alerts?state=firing&limit=20", &alerts); err == nil {
		for _, a := range alerts.Data {
			if a.Severity != "critical" && a.Severity != "high" {
				continue
			}
			p.upsert(&models.Notification{
				Source:   "pulse",
				Severity: a.Severity,
				Title:    "Alert — " + a.Name,
				Body:     a.Message,
				Link:     "/pulse/ui/alerts?focus=" + a.ID,
				Dedup:    dedup("pulse", "alert", a.ID),
			})
		}
	}

	// — Pulse: top-impact N+1 (single notification, dedup'd by hour)
	var n1 struct {
		Data []struct {
			Route               string  `json:"route"`
			Pattern             string  `json:"pattern"`
			ImpactScore         float64 `json:"impact_score"`
			Occurrences         int     `json:"occurrences"`
			AvgQueriesPerReq    int     `json:"avg_queries_per_request"`
		} `json:"data"`
	}
	if err := p.Bridge.PulseGet(ctx, "/pulse/api/database/n1/ranked?range=1h&limit=1", &n1); err == nil && len(n1.Data) > 0 {
		top := n1.Data[0]
		if top.ImpactScore > 100 {
			p.upsert(&models.Notification{
				Source:   "pulse",
				Severity: "medium",
				Title:    "N+1 query detected — " + top.Route,
				Body:     fmt.Sprintf("%d occurrences · ~%d queries/req · impact %.0f", top.Occurrences, top.AvgQueriesPerReq, top.ImpactScore),
				Link:     "/pulse/ui/database#n1",
				Dedup:    dedup("pulse", "n1", time.Now().UTC().Format("2006-01-02-15"), top.Route),
			})
		}
	}
}

func (p *SecObsPoller) upsert(n *models.Notification) {
	res := p.DB.Model(&models.Notification{}).
		Where("dedup = ?", n.Dedup).
		Updates(map[string]interface{}{
			"count":      gorm.Expr("count + 1"),
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		log.Printf("secobs_poller: update %s: %v", n.Dedup, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		if err := p.DB.Create(n).Error; err != nil {
			log.Printf("secobs_poller: create %s: %v", n.Dedup, err)
		}
	}
}

func dedup(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:32]
}
