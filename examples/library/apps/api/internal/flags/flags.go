// Package flags is the feature flag + A/B testing engine.
//
// At a glance:
//
//	if flags.IsEnabled(c, "new_dashboard") {
//	    // … render the new dashboard
//	}
//
//	switch flags.Variant(c, "checkout_redesign") {
//	case "control":   /* old flow */
//	case "variant_a": /* new flow */
//	case "variant_b": /* alternate new flow */
//	}
//
// Mechanics:
//   - All flags are loaded into an in-memory map at boot. A background
//     goroutine refreshes every 30s. Flag checks never hit the DB.
//   - Bucketing: SHA-256(user_id || ":" || flag_name) % 100. Sticky
//     per (user, flag) — a user always gets the same bucket for a
//     given flag, so variant assignment doesn't flicker across sessions.
//   - Anonymous users (empty user_id) bucket on a random per-request
//     value, which is effectively random. For sticky anonymous flags
//     pass a stable identifier (session ID, device ID).
//   - Exposure tracking is fire-and-forget — flag checks never block
//     on the DB.
//   - When a flag is created/updated/deleted, the engine refreshes
//     immediately and broadcasts a "flag.updated" realtime event so
//     subscribed clients can refetch.
package flags

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/models"
	"library/apps/api/internal/realtime"
)

// DefaultRefreshInterval is how often the engine pulls fresh flag
// state from the DB. 30s is a reasonable middle ground — admin
// changes propagate quickly without hammering the DB.
const DefaultRefreshInterval = 30 * time.Second

// Engine owns the in-memory flag cache. One per process.
type Engine struct {
	db    *gorm.DB
	hub   *realtime.Hub // optional — when set, broadcasts on Refresh
	mu    sync.RWMutex
	flags map[string]*models.FeatureFlag
	stop  chan struct{}
	// exposures feeds one writer. Each flag check used to start a goroutine and
	// an INSERT of its own, which a busy page turned into thousands a second.
	exposures chan models.FlagExposure
}

// New returns an Engine with the cache pre-warmed. Call from
// routes.Setup. hub is optional — pass nil to disable broadcasts.
func New(db *gorm.DB, hub *realtime.Hub) *Engine {
	e := &Engine{
		db:        db,
		hub:       hub,
		flags:     make(map[string]*models.FeatureFlag),
		stop:      make(chan struct{}),
		exposures: make(chan models.FlagExposure, exposureQueueSize),
	}
	if err := e.Refresh(); err != nil {
		log.Printf("[flags] initial refresh failed: %v", err)
	}
	go e.refreshLoop()
	go e.writeExposures()
	setDefault(e)
	return e
}

// Stop terminates the background refresh goroutine. Call on graceful
// shutdown to avoid leaking goroutines in tests.
func (e *Engine) Stop() {
	close(e.stop)
}

// Refresh pulls all flags from the DB and replaces the cache. Called
// every DefaultRefreshInterval and immediately after admin writes.
func (e *Engine) Refresh() error {
	var rows []models.FeatureFlag
	if err := e.db.Find(&rows).Error; err != nil {
		return err
	}
	next := make(map[string]*models.FeatureFlag, len(rows))
	for i := range rows {
		f := rows[i]
		next[f.Name] = &f
	}
	e.mu.Lock()
	e.flags = next
	e.mu.Unlock()
	return nil
}

// RefreshAndBroadcast refreshes the cache and (if a hub was provided)
// emits a "flag.updated" realtime event so subscribed clients can
// refetch. Call after admin writes.
func (e *Engine) RefreshAndBroadcast(flagName string) {
	if err := e.Refresh(); err != nil {
		log.Printf("[flags] refresh after change failed: %v", err)
	}
	if e.hub != nil {
		e.hub.Broadcast(realtime.Event{
			Type:    "flag.updated",
			Payload: map[string]interface{}{"name": flagName},
		})
	}
}

func (e *Engine) refreshLoop() {
	t := time.NewTicker(DefaultRefreshInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := e.Refresh(); err != nil {
				log.Printf("[flags] periodic refresh failed: %v", err)
			}
		case <-e.stop:
			return
		}
	}
}

// Subject is who a flag is checked for: a user, and whatever attributes
// the rules can target (a business unit, a region, a plan).
type Subject struct {
	UserID     string
	Attributes map[string]string
}

// AttributesFor supplies the attributes of the user making a request, for
// rules that target them. The default knows the role the auth middleware
// set. Replace it at boot to target anything else your users carry:
//
//	flags.AttributesFor = func(c *gin.Context) map[string]string {
//	    return map[string]string{"business_unit": businessUnitOf(c)}
//	}
var AttributesFor = func(c *gin.Context) map[string]string {
	attrs := map[string]string{}
	if c == nil {
		return attrs
	}
	if v, ok := c.Get("user_role"); ok {
		if role, ok := v.(string); ok && role != "" {
			attrs["role"] = role
		}
	}
	return attrs
}

func subjectFrom(c *gin.Context) Subject {
	return Subject{UserID: userIDFrom(c), Attributes: AttributesFor(c)}
}

// IsEnabled returns true when the flag is on for the current user.
// Always returns false for unknown flags (fail closed).
func (e *Engine) IsEnabled(c *gin.Context, name string) bool {
	return e.evaluate(subjectFrom(c), name) == "enabled"
}

// Variant returns the assigned variant for an A/B flag. For boolean
// flags, returns "enabled" or "disabled". For unknown flags, returns
// the empty string.
func (e *Engine) Variant(c *gin.Context, name string) string {
	return e.evaluate(subjectFrom(c), name)
}

// IsEnabledForUser is the explicit form for backend code that has the
// user_id directly (e.g. cron jobs operating on a specific user). It has
// no attributes, so a flag with attribute rules is off for it: use
// IsEnabledFor with a Subject for those.
func (e *Engine) IsEnabledForUser(userID, name string) bool {
	return e.evaluate(Subject{UserID: userID}, name) == "enabled"
}

// VariantForUser is the explicit form of Variant.
func (e *Engine) VariantForUser(userID, name string) string {
	return e.evaluate(Subject{UserID: userID}, name)
}

// IsEnabledFor checks a flag for a subject built by the caller: a job that
// knows the user and their business unit, say.
func (e *Engine) IsEnabledFor(s Subject, name string) bool {
	return e.evaluate(s, name) == "enabled"
}

// VariantFor is the explicit form of Variant for a subject.
func (e *Engine) VariantFor(s Subject, name string) string {
	return e.evaluate(s, name)
}

// ── Package-level ────────────────────────────────────────────────────────
//
// The engine routes.Setup starts, reachable from any handler, service or job
// without threading it through. Before v3.223.0 the package comment showed
// flags.IsEnabled(c, ...) and no such function existed: the only engine was
// a local variable in routes.Setup, so application code could not check a
// flag at all.

var (
	defaultMu     sync.RWMutex
	defaultEngine *Engine
)

func setDefault(e *Engine) {
	defaultMu.Lock()
	defaultEngine = e
	defaultMu.Unlock()
}

func current() *Engine {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultEngine
}

// IsEnabled reports whether a flag is on for the request's user. False
// before the engine has started and for a flag that does not exist: a check
// that cannot be answered fails closed.
func IsEnabled(c *gin.Context, name string) bool {
	if e := current(); e != nil {
		return e.IsEnabled(c, name)
	}
	return false
}

// Variant returns the request user's variant, or "" when it cannot be answered.
func Variant(c *gin.Context, name string) string {
	if e := current(); e != nil {
		return e.Variant(c, name)
	}
	return ""
}

// IsEnabledForUser is IsEnabled for code with a user ID and no request.
func IsEnabledForUser(userID, name string) bool {
	if e := current(); e != nil {
		return e.IsEnabledForUser(userID, name)
	}
	return false
}

// VariantForUser is Variant for code with a user ID and no request.
func VariantForUser(userID, name string) string {
	if e := current(); e != nil {
		return e.VariantForUser(userID, name)
	}
	return ""
}

// IsEnabledFor is IsEnabled for a subject built by the caller.
func IsEnabledFor(s Subject, name string) bool {
	if e := current(); e != nil {
		return e.IsEnabledFor(s, name)
	}
	return false
}

// VariantFor is Variant for a subject built by the caller.
func VariantFor(s Subject, name string) string {
	if e := current(); e != nil {
		return e.VariantFor(s, name)
	}
	return ""
}

// evaluate is the core decision routine. Returns:
//
//	""           — unknown flag
//	"disabled"   — flag exists but rules deny the user
//	"enabled"    — boolean flag passed; user is in the rollout
//	"<variant>"  — A/B flag passed; the user's bucket maps to this variant
//
// Lock discipline: the read lock is held only long enough to copy the
// flag struct + ID. All decision logic (date checks, allowlist scans,
// bucketing) runs unlocked. Under sustained read load this turns the
// flag check into a near-zero-contention path.
func (e *Engine) evaluate(s Subject, name string) string {
	userID := s.UserID
	e.mu.RLock()
	cached, ok := e.flags[name]
	if !ok {
		e.mu.RUnlock()
		return ""
	}
	flagID := cached.ID
	enabled := cached.Enabled
	rulesJSON := cached.Rules
	e.mu.RUnlock()

	if !enabled {
		return "disabled"
	}

	// Decode rules outside the lock — JSON parsing is the slowest
	// part of the flag check and we don't want it serializing.
	flagForParse := models.FeatureFlag{Rules: rulesJSON}
	rules := flagForParse.ParsedRules()

	// Date window — out-of-window short-circuits before bucketing.
	now := time.Now()
	if rules.EnabledFrom != nil && now.Before(*rules.EnabledFrom) {
		return "disabled"
	}
	if rules.EnabledUntil != nil && now.After(*rules.EnabledUntil) {
		return "disabled"
	}

	// Blocklist always wins.
	for _, b := range rules.BlocklistUserIDs {
		if b == userID {
			return "disabled"
		}
	}

	// Attributes restrict the flag to matching subjects, every key listed.
	for key, allowed := range rules.Attributes {
		if !attributeMatches(s.Attributes[key], allowed) {
			return "disabled"
		}
	}

	// Allowlist (when non-empty) restricts to the listed users.
	// Skip the percentage roll for allowlisted users — they always
	// see it, that's the point.
	allowlistMode := len(rules.AllowlistUserIDs) > 0
	allowed := false
	for _, a := range rules.AllowlistUserIDs {
		if a == userID {
			allowed = true
			break
		}
	}
	if allowlistMode && !allowed {
		return "disabled"
	}

	bucket := bucketFor(userID, name)

	// A/B mode — assign variant by bucket.
	if len(rules.Variants) > 0 {
		v := rules.Variants[bucket%len(rules.Variants)]
		e.trackExposure(flagID, name, userID, v)
		return v
	}

	// Boolean mode — percentage rollout. Allowlisted users always
	// pass; everyone else is gated by the percentage.
	if allowed || bucket < rules.RolloutPercentage {
		e.trackExposure(flagID, name, userID, "enabled")
		return "enabled"
	}
	e.trackExposure(flagID, name, userID, "disabled")
	return "disabled"
}

// exposureQueueSize bounds the exposures waiting for the writer. Past it they
// are dropped: an exposure is analytics, never worth slowing a request for.
const exposureQueueSize = 4096

// trackExposure hands the flag check to the exposure writer without blocking.
func (e *Engine) trackExposure(flagID, flagName, userID, variant string) {
	if userID == "" {
		// Anonymous exposures pollute the table without buying us
		// anything (we can't link them to a user later). Skip.
		return
	}
	select {
	case e.exposures <- models.FlagExposure{FlagID: flagID, FlagName: flagName, UserID: userID, Variant: variant}:
	default:
	}
}

// writeExposures writes exposures in batches, and records each user, flag and
// variant once a UTC day: an exposure says who saw which variant, not how many
// times they reloaded the page.
func (e *Engine) writeExposures() {
	seen := map[string]bool{}
	day := ""
	for {
		var first models.FlagExposure
		select {
		case <-e.stop:
			return
		case first = <-e.exposures:
		}
		batch := []models.FlagExposure{first}
	drain:
		for len(batch) < 256 {
			select {
			case x := <-e.exposures:
				batch = append(batch, x)
			default:
				break drain
			}
		}

		today := time.Now().UTC().Format("2006-01-02")
		if today != day || len(seen) > 100000 {
			seen, day = map[string]bool{}, today
		}
		fresh := batch[:0]
		for _, x := range batch {
			k := x.UserID + "|" + x.FlagID + "|" + x.Variant
			if seen[k] {
				continue
			}
			seen[k] = true
			fresh = append(fresh, x)
		}
		if len(fresh) == 0 {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := e.db.WithContext(ctx).CreateInBatches(fresh, 256).Error; err != nil {
			log.Printf("[flags] recording %d exposures: %v", len(fresh), err)
		}
		cancel()
	}
}

// bucketFor hashes (userID || ":" || flagName) and returns the bucket
// 0..99. Same input always produces the same bucket — that's what
// makes the assignment sticky.
//
// We use SHA-256 (not Go's default hash) because it's stable across
// process restarts + Go versions. FNV would be faster but Grit isn't
// running flag checks in a hot loop — sub-microsecond cost is fine.
func bucketFor(userID, name string) int {
	if userID == "" {
		// Anonymous users get a uniform random bucket. We avoid
		// UnixNano%100 because nanosecond timing is biased toward
		// recent buckets under high QPS. crypto/rand gives us a
		// uniform draw without that artifact.
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			// rand should never fail on a healthy OS; if it does,
			// fall back to bucket 0 so behavior is deterministic.
			return 0
		}
		return int(binary.BigEndian.Uint32(b[:]) % 100)
	}
	h := sha256.Sum256([]byte(userID + ":" + name))
	return int(binary.BigEndian.Uint32(h[:4]) % 100)
}

// userIDFrom reads "user_id" from the gin context (set by the auth
// middleware). Empty string for anonymous requests.
func userIDFrom(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Get("user_id"); ok {
		s, _ := v.(string)
		return s
	}
	return ""
}

// attributeMatches reports whether a subject's value is one of a rule's.
// Case-insensitive, because these are typed in by hand. An empty value
// matches nothing.
func attributeMatches(value string, allowed []string) bool {
	if value == "" {
		return false
	}
	for _, a := range allowed {
		if strings.EqualFold(a, value) {
			return true
		}
	}
	return false
}
