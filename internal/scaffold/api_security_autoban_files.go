package scaffold

// The auto-ban sweep: the thing that turns a threat log into a defence.
//
// Sentinel's own AutoBlock (v2.8.0) catches a burst: ten high-severity threats
// from one address inside ten minutes. That is a scanner going fast. It does
// not catch a scanner going slowly, and a slow one is the more common shape:
// one deployment's log had 18,405 rows, nearly all of them env-file probes
// from a handful of addresses, every one detected and recorded, and not one of
// them stopped, because no single address ever tripped ten in ten minutes.
//
// This sweep reads the window rather than the burst. Every few hours it counts
// each address's high-severity threats over the last day and bans the ones
// over the line. The two are complements: the burst rule answers in seconds
// and this one answers for everything patient enough to avoid it.
//
// It counts high severity only, deliberately. ThreatCount on a Sentinel actor
// is every severity, so a client tripping a low-severity heuristic forty times
// would be banned by a count that said nothing about intent. High is what the
// scanner patterns score.

func securityAutoBanGo() string {
	return `package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"{{MODULE}}/internal/config"
)

// AutoBanConfig is the sweep's policy, read from the environment.
type AutoBanConfig struct {
	// Enabled turns the sweep off without removing it from the schedule.
	Enabled bool
	// MinThreats is the count an address has to exceed, not meet. At the
	// default of 3, an address is banned on its fourth high-severity threat
	// inside the window.
	MinThreats int
	// Window is how far back the count looks.
	Window time.Duration
	// Permanent bans until somebody lifts it. Otherwise Duration applies.
	Permanent bool
	Duration  time.Duration
	// MaxBans caps one sweep. A log that has gone wrong should not be able to
	// ban a thousand addresses in one pass before anybody sees it.
	MaxBans int
}

// AutoBanResult is what one sweep did, for the log and for the dashboard.
type AutoBanResult struct {
	Scanned   int               ` + "`json:\"scanned\"`" + `
	Matched   int               ` + "`json:\"matched\"`" + `
	Banned    []string          ` + "`json:\"banned\"`" + `
	Skipped   []string          ` + "`json:\"skipped\"`" + `
	Counts    map[string]int    ` + "`json:\"counts\"`" + `
	Failed    map[string]string ` + "`json:\"failed,omitempty\"`" + `
	DryRun    bool              ` + "`json:\"dry_run\"`" + `
	Truncated bool              ` + "`json:\"truncated\"`" + `
}

// AutoBanReason is the prefix every ban this sweep writes begins with. The
// dashboard tells an automatic ban from a hand-made one by it, so changing it
// changes that split. Sentinel's BlockedIP carries no flag for this.
const AutoBanReason = "auto-ban"

// IsAutoBan reports whether a ban reason was written by this sweep or by
// Sentinel's own burst rule.
func IsAutoBan(reason string) bool {
	r := strings.ToLower(reason)
	return strings.HasPrefix(r, AutoBanReason) || strings.HasPrefix(r, "auto-blocked")
}

// autoBanPageSize is how many threats are read per request. The sweep pages
// until the window is covered or the cap is hit.
const autoBanPageSize = 500

// autoBanMaxPages bounds one sweep. At 500 a page that is 25,000 threats,
// which is more than a day of even a bad week, and it stops a runaway log from
// turning into an unbounded loop of HTTP calls.
const autoBanMaxPages = 50

// SecurityAutoBan counts high-severity threats per address and bans the
// addresses over the line.
type SecurityAutoBan struct {
	Bridge *SecObsBridge
	Config AutoBanConfig
}

// NewSecurityAutoBan builds the sweep from the project's configuration, or
// returns nil when Sentinel is not mounted, because a sweep with nothing to
// read is not a sweep that should be scheduled.
func NewSecurityAutoBan(bridge *SecObsBridge, cfg *config.Config) *SecurityAutoBan {
	if bridge == nil {
		return nil
	}
	return &SecurityAutoBan{
		Bridge: bridge,
		Config: AutoBanConfig{
			Enabled:    cfg.SecurityAutoBanEnabled,
			MinThreats: cfg.SecurityAutoBanMinThreats,
			Window:     cfg.SecurityAutoBanWindow,
			Permanent:  cfg.SecurityAutoBanPermanent,
			Duration:   7 * 24 * time.Hour,
			MaxBans:    200,
		},
	}
}

// severities counted. Critical is included because an address that scores it
// is not a milder case than High.
//
// Capitalised, because Sentinel stores and filters on the exact string: its
// constants are "Critical", "High", "Medium", "Low". Asking for "high" matches
// nothing and returns an empty page, which reads exactly like a quiet week.
// The first version of this swept 0 threats on a log holding eleven.
var autoBanSeverities = []string{SeverityCritical, SeverityHigh}

// The severity words Sentinel stores. Everything here that names a severity
// uses these rather than a literal, because the mismatch is silent.
const (
	SeverityCritical = "Critical"
	SeverityHigh     = "High"
	SeverityMedium   = "Medium"
	SeverityLow      = "Low"
)

// CanonicalSeverity turns what a filter sends ("high", "HIGH") into what
// Sentinel stores ("High"), and reports whether it is a severity at all.
func CanonicalSeverity(input string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "critical":
		return SeverityCritical, true
	case "high":
		return SeverityHigh, true
	case "medium":
		return SeverityMedium, true
	case "low":
		return SeverityLow, true
	}
	return "", false
}

type autoBanThreatPage struct {
	Data []struct {
		IP        string    ` + "`json:\"ip\"`" + `
		Severity  string    ` + "`json:\"severity\"`" + `
		Timestamp time.Time ` + "`json:\"timestamp\"`" + `
	} ` + "`json:\"data\"`" + `
	Meta struct {
		Total int ` + "`json:\"total\"`" + `
	} ` + "`json:\"meta\"`" + `
}

type autoBanBlockedList struct {
	Data []struct {
		IP     string ` + "`json:\"ip\"`" + `
		Reason string ` + "`json:\"reason\"`" + `
	} ` + "`json:\"data\"`" + `
}

// Run performs one sweep. With dryRun it reports who would be banned and bans
// nobody, which is what the dashboard's button shows before it is pressed.
func (s *SecurityAutoBan) Run(ctx context.Context, dryRun bool) (*AutoBanResult, error) {
	if s == nil || s.Bridge == nil {
		return nil, fmt.Errorf("security auto-ban: Sentinel is not configured")
	}
	cfg := s.Config
	if cfg.MinThreats < 1 {
		cfg.MinThreats = 3
	}
	if cfg.Window <= 0 {
		cfg.Window = 24 * time.Hour
	}
	if cfg.MaxBans < 1 {
		cfg.MaxBans = 200
	}

	counts, scanned, truncated, err := s.count(ctx, cfg.Window)
	if err != nil {
		return nil, err
	}

	// Already banned, so there is nothing to do and nothing to report as done.
	var blocked autoBanBlockedList
	if err := s.Bridge.SentinelGet(ctx, "/sentinel/api/ip/blocked", &blocked); err != nil {
		return nil, fmt.Errorf("security auto-ban: reading the ban list: %w", err)
	}
	already := make(map[string]bool, len(blocked.Data))
	for _, b := range blocked.Data {
		already[b.IP] = true
	}

	result := &AutoBanResult{
		Scanned:   scanned,
		Counts:    map[string]int{},
		Banned:    []string{},
		Skipped:   []string{},
		DryRun:    dryRun,
		Truncated: truncated,
	}

	// Worst first, so a capped sweep bans the loudest rather than whichever
	// address a map happened to yield first.
	type candidate struct {
		ip    string
		count int
	}
	var candidates []candidate
	for ip, n := range counts {
		if n <= cfg.MinThreats {
			continue
		}
		candidates = append(candidates, candidate{ip, n})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].count != candidates[j].count {
			return candidates[i].count > candidates[j].count
		}
		return candidates[i].ip < candidates[j].ip
	})
	result.Matched = len(candidates)

	for _, c := range candidates {
		result.Counts[c.ip] = c.count
		if already[c.ip] {
			result.Skipped = append(result.Skipped, c.ip)
			continue
		}
		// Never ban the address the app itself speaks from. Behind a reverse
		// proxy that Sentinel has not been told to trust, every request
		// carries the proxy's address, and banning it bans everybody. The
		// same reasoning as the AutoBlock exemptions in routes.go.
		if IsPrivateOrLoopbackIP(c.ip) {
			result.Skipped = append(result.Skipped, c.ip)
			continue
		}
		if len(result.Banned) >= cfg.MaxBans {
			result.Truncated = true
			break
		}
		if dryRun {
			result.Banned = append(result.Banned, c.ip)
			continue
		}
		body := map[string]interface{}{
			"ip": c.ip,
			"reason": fmt.Sprintf("%s: %d high-severity threats in %s",
				AutoBanReason, c.count, humanWindow(cfg.Window)),
		}
		if cfg.Permanent {
			body["permanent"] = true
		} else if cfg.Duration > 0 {
			body["duration"] = cfg.Duration.String()
		}
		if err := s.Bridge.SentinelPost(ctx, "/sentinel/api/ip/block", body, nil); err != nil {
			if result.Failed == nil {
				result.Failed = map[string]string{}
			}
			result.Failed[c.ip] = err.Error()
			continue
		}
		result.Banned = append(result.Banned, c.ip)
	}
	return result, nil
}

// count reads the threat log for the window and totals high-severity threats
// per address.
//
// Sentinel can sort actors by threat_count, which would be one request, but
// that count is every severity. A client tripping a low-severity heuristic is
// not a scanner, and banning by that number is how a dashboard starts refusing
// real users.
func (s *SecurityAutoBan) count(ctx context.Context, window time.Duration) (map[string]int, int, bool, error) {
	counts := map[string]int{}
	since := time.Now().Add(-window).UTC().Format(time.RFC3339)
	scanned := 0
	truncated := false

	for _, severity := range autoBanSeverities {
		for page := 1; page <= autoBanMaxPages; page++ {
			path := fmt.Sprintf("/sentinel/api/threats?severity=%s&start_time=%s&page=%d&page_size=%d",
				url.QueryEscape(severity), url.QueryEscape(since), page, autoBanPageSize)
			var got autoBanThreatPage
			if err := s.Bridge.SentinelGet(ctx, path, &got); err != nil {
				return nil, scanned, truncated, fmt.Errorf("security auto-ban: reading %s threats: %w", severity, err)
			}
			for _, t := range got.Data {
				if t.IP == "" {
					continue
				}
				counts[t.IP]++
				scanned++
			}
			if len(got.Data) < autoBanPageSize {
				break
			}
			if page == autoBanMaxPages {
				// The window holds more than this sweep will read. The
				// addresses already counted are the busiest, so the sweep is
				// still useful; saying so beats pretending it was complete.
				truncated = true
			}
		}
	}
	return counts, scanned, truncated, nil
}

// IsPrivateOrLoopbackIP reports whether an address is one the app's own
// infrastructure speaks from. Exported because the hand-ban endpoint refuses
// the same addresses the sweep skips, and two copies of that list would drift.
func IsPrivateOrLoopbackIP(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		// Not an address this sweep understands, so not one it will ban.
		return true
	}
	return parsed.IsLoopback() || parsed.IsPrivate() || parsed.IsLinkLocalUnicast() ||
		parsed.IsUnspecified() || isCarrierNAT(parsed)
}

// isCarrierNAT covers 100.64.0.0/10, which is shared address space: a carrier,
// or an ingress controller handing out addresses from it.
func isCarrierNAT(ip net.IP) bool {
	v4 := ip.To4()
	if v4 == nil {
		return false
	}
	return v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127
}

// humanWindow prints a duration the way the ban reason should read.
func humanWindow(d time.Duration) string {
	if d%(24*time.Hour) == 0 && d >= 24*time.Hour {
		days := int(d / (24 * time.Hour))
		if days == 1 {
			return "24 hours"
		}
		return fmt.Sprintf("%d days", days)
	}
	if d%time.Hour == 0 {
		return fmt.Sprintf("%d hours", int(d/time.Hour))
	}
	return d.String()
}

// LogResult prints one line a human can act on, which is what a scheduled job
// leaves behind when nobody is watching it run.
func (r *AutoBanResult) LogResult(cfg AutoBanConfig) {
	if r == nil {
		return
	}
	if len(r.Banned) == 0 {
		log.Printf("Auto-ban sweep: %d high-severity threats from %d address(es), none over %d",
			r.Scanned, len(r.Counts), cfg.MinThreats)
		return
	}
	log.Printf("Auto-ban sweep: banned %d address(es) over %d high-severity threats in %s: %s",
		len(r.Banned), cfg.MinThreats, humanWindow(cfg.Window), strings.Join(r.Banned, ", "))
	for ip, err := range r.Failed {
		log.Printf("Auto-ban sweep: %s could not be banned: %s", ip, err)
	}
	if r.Truncated {
		log.Printf("Auto-ban sweep: the window held more threats than one sweep reads; " +
			"the busiest addresses were counted")
	}
}
`
}
