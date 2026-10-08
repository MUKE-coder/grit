// Package audit owns the tamper-evident hash chain over the activity log.
//
// Each row's Hash = SHA-256(PrevHash || canonical(row)) where canonical
// is a stable JSON serialization of the audit-relevant fields. Any
// mutation to a row breaks every Hash from that row forward, which
// VerifyChain detects.
//
// Insert is serialized via a row-level FOR UPDATE lock on the latest
// row inside the same transaction that does the INSERT — concurrent
// inserts queue cleanly without forking the chain. Verification walks
// the chain in created_at + id order; ties broken by id.
//
// What this defends against:
//   - Direct SQL UPDATE / DELETE on activity_logs (most common attack
//     vector — DBA covering tracks).
//   - Out-of-band insertion of forged history.
//
// What this does NOT defend against:
//   - Compromise of the running server itself (an attacker with code
//     execution can rewrite the whole chain). External anchoring
//     (publishing the daily root hash to a public ledger) is the
//     follow-up — see #48.
package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"saas/apps/api/internal/models"
)

// Canonical returns the stable JSON bytes of an entry for hashing.
// We exclude ID / PrevHash / Hash from the canonical form: ID is
// random and uncorrelated with content; PrevHash + Hash are derived
// values, not inputs to the hash.
func Canonical(e *models.ActivityLog) ([]byte, error) {
	c := canonicalEntry{
		UserID:        e.UserID,
		Method:        e.Method,
		Path:          e.Path,
		Status:        e.Status,
		PayloadDigest: e.PayloadDigest,
		IPAddress:     e.IPAddress,
		UserAgent:     e.UserAgent,
		DurationMS:    e.DurationMS,
		// Use unix-nano so the canonical bytes are stable across tz
		// changes / TIMESTAMPTZ formatting differences.
		CreatedAtUnixNano: e.CreatedAt.UTC().UnixNano(),
		Resource:          e.Resource,
		ResourceIDs:       e.ResourceIDs,
		RecordCount:       e.RecordCount,
	}
	return json.Marshal(c)
}

// canonicalEntry's field order is the wire format for hashing —
// reorder ONLY in a major version bump (verify breaks otherwise).
type canonicalEntry struct {
	UserID            string `json:"user_id"`
	Method            string `json:"method"`
	Path              string `json:"path"`
	Status            int    `json:"status"`
	PayloadDigest     string `json:"payload_digest"`
	IPAddress         string `json:"ip_address"`
	UserAgent         string `json:"user_agent"`
	DurationMS        int64  `json:"duration_ms"`
	CreatedAtUnixNano int64  `json:"created_at_unix_nano"`
	// Appended, and omitted when empty, so every entry written before they
	// existed has the same bytes and old chains still verify.
	Resource    string `json:"resource,omitempty"`
	ResourceIDs string `json:"resource_ids,omitempty"`
	RecordCount int    `json:"record_count,omitempty"`
}

// ComputeHash returns hex(sha256(prevHash || canonical)) — the prev
// hash is included as a hex string (not raw bytes) so the input is
// trivially auditable: cat prev_hash | xxd; cat canonical.json.
func ComputeHash(prevHash string, canonical []byte) string {
	h := sha256.New()
	h.Write([]byte(prevHash))
	h.Write(canonical)
	return hex.EncodeToString(h.Sum(nil))
}

// MaxReadIDs caps how many ids one read entry lists. A page is a few hundred
// rows at most; an export is recorded as a count instead.
const MaxReadIDs = 1000

const readKey = "audit.read"

// ReadMark is what a handler served, for the activity middleware to record.
type ReadMark struct {
	Resource string
	IDs      []string
	Count    int
}

// Read marks the request as having served these rows of resource. The
// activity middleware records it in the chain once the response has gone out
// with a 2xx. Handlers generated with --audit-reads call it; any handler can.
func Read(c *gin.Context, resource string, ids ...string) {
	mark := ReadMark{Resource: resource, Count: len(ids)}
	if len(ids) > MaxReadIDs {
		ids = ids[:MaxReadIDs]
	}
	mark.IDs = append([]string(nil), ids...)
	c.Set(readKey, mark)
}

// ReadCount marks a read too large to list row by row: an export.
func ReadCount(c *gin.Context, resource string, n int) {
	c.Set(readKey, ReadMark{Resource: resource, Count: n})
}

// ReadMarkOf returns what the handler marked, if it marked anything.
func ReadMarkOf(c *gin.Context) (ReadMark, bool) {
	v, ok := c.Get(readKey)
	if !ok {
		return ReadMark{}, false
	}
	mark, ok := v.(ReadMark)
	return mark, ok
}

// Precision is what every supported database keeps of a timestamp: Postgres
// stores microseconds and MySQL, as GORM creates the column, milliseconds.
// The hash covers created_at, so a stamp finer than the column gives a hash
// nobody can recompute from the stored row. Before v3.215.0 every entry carried
// nanoseconds, and the chain failed verification on its first row on both.
const Precision = time.Millisecond

// chainLockKey is the Postgres advisory lock every chain writer takes, so two
// API replicas never read the same latest hash and fork the chain.
const chainLockKey int64 = 0x677269745f617564

var (
	queue     = make(chan models.ActivityLog, 4096)
	startOnce sync.Once
	dropped   atomic.Uint64
)

// Start runs this process's chain writer. Safe to call more than once.
//
// One writer fed by a bounded channel, so a burst of requests never waits on
// the database or spawns a goroutine per entry. It is not what keeps the chain
// whole across processes: every batch takes the chain lock and reads the latest
// hash from the database, so each replica can run one.
func Start(db *gorm.DB) {
	startOnce.Do(func() { go writer(db) })
}

// Enqueue hands an entry to the writer without blocking, and reports false
// when the backlog was full and the entry was dropped. Losing an audit row is
// better than stalling the request path; Dropped says how often it happened.
func Enqueue(entry models.ActivityLog) bool {
	select {
	case queue <- entry:
		return true
	default:
		dropped.Add(1)
		return false
	}
}

// Dropped is how many entries Enqueue has dropped since the process started.
// Sustained growth means the writer cannot keep up.
func Dropped() uint64 { return dropped.Load() }

func writer(db *gorm.DB) {
	for first := range queue {
		batch := []models.ActivityLog{first}
	drain:
		for len(batch) < 256 {
			select {
			case e := <-queue:
				batch = append(batch, e)
			default:
				break drain
			}
		}
		if err := appendBatch(db, batch); err != nil {
			// One bad entry must not take the rest of the batch with it.
			for _, e := range batch {
				if err := appendBatch(db, []models.ActivityLog{e}); err != nil {
					log.Printf("[audit] could not record %s %s: %v", e.Method, e.Path, err)
				}
			}
		}
	}
}

// AppendChained writes one entry and returns once it is stored, for a caller
// that must know it landed: a security event, a reseal. Same lock and same
// stamping as the writer, so the two never fork the chain between them.
func AppendChained(db *gorm.DB, entry *models.ActivityLog) error {
	batch := []models.ActivityLog{*entry}
	if err := appendBatch(db, batch); err != nil {
		return err
	}
	*entry = batch[0]
	return nil
}

// appendBatch chains entries onto the latest stored one, in one transaction
// holding the chain lock. created_at is stamped here rather than by the caller:
// at the precision every database keeps, and strictly after the entry before,
// so VerifyChain's (created_at, id) order is the order the chain was written.
func appendBatch(db *gorm.DB, entries []models.ActivityLog) error {
	return db.Transaction(func(tx *gorm.DB) error {
		head, err := lockAndReadHead(tx)
		if err != nil {
			return err
		}
		prevHash, last := head.Hash, head.CreatedAt
		for i := range entries {
			e := &entries[i]
			clipToColumns(e)
			e.CreatedAt = nextStamp(last)
			canonical, err := Canonical(e)
			if err != nil {
				return fmt.Errorf("canonicalize: %w", err)
			}
			e.PrevHash = prevHash
			e.Hash = ComputeHash(prevHash, canonical)
			prevHash, last = e.Hash, e.CreatedAt
		}
		// One multi-row INSERT for the batch. An INSERT per entry kept the chain
		// lock, which every replica's writer waits on, for 256 round trips.
		return tx.CreateInBatches(entries, 256).Error
	})
}

// lockAndReadHead takes the chain lock and returns the latest entry, or a zero
// entry for an empty log. Postgres takes an advisory lock held until the
// transaction ends; MySQL a locking read, which sees the latest committed row
// once granted; SQLite serialises writers on its own.
func lockAndReadHead(tx *gorm.DB) (models.ActivityLog, error) {
	var head models.ActivityLog
	q := tx.Order("created_at desc, id desc").Limit(1)
	switch tx.Dialector.Name() {
	case "postgres":
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", chainLockKey).Error; err != nil {
			return head, fmt.Errorf("taking the audit chain lock: %w", err)
		}
	case "mysql":
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.Find(&head).Error; err != nil {
		return head, fmt.Errorf("reading the chain head: %w", err)
	}
	return head, nil
}

// nextStamp is now at the chain's precision, and strictly after prev.
func nextStamp(prev time.Time) time.Time {
	now := time.Now().UTC().Truncate(Precision)
	if !prev.IsZero() {
		if floor := prev.UTC().Truncate(Precision).Add(Precision); now.Before(floor) {
			now = floor
		}
	}
	return now
}

// clipToColumns trims what could overflow its column, before hashing, so the
// stored row is the hashed row. An oversized user agent used to fail the insert
// and lose the entry.
func clipToColumns(e *models.ActivityLog) {
	e.Method = clip(e.Method, 10)
	e.Path = clip(e.Path, 500)
	e.IPAddress = clip(e.IPAddress, 45)
	e.UserAgent = clip(e.UserAgent, 500)
	e.Resource = clip(e.Resource, 100)
}

// clip shortens s to at most n bytes without splitting a character, which a
// database would refuse as invalid UTF-8.
func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	s = s[:n]
	for len(s) > 0 && !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// PruneChunk is how many entries one prune transaction deletes.
const PruneChunk = 5000

// prunePath marks the SECURITY entry a prune appends.
const prunePath = "audit.chain.pruned"

// Prune deletes the entries created before cutoff, oldest first and PruneChunk
// at a time, and returns how many it deleted.
//
// The entries that remain keep their hashes. The oldest of them now chains
// from an entry that is gone, so each chunk appends a SECURITY entry naming the
// hash it chains from, and VerifyChain accepts a log that starts there only
// when that entry exists: old entries deleted any other way still fail
// verification. Each chunk is one transaction holding the chain lock, so the
// log verifies after every commit and a prune that stops part way keeps what
// it did.
//
// Pruning used to rewrite the oldest remaining entry's hash, which broke the
// link from the entry after it, so every prune that deleted anything left a
// log that failed verification.
func Prune(ctx context.Context, db *gorm.DB, cutoff time.Time) (int64, error) {
	var total int64
	for {
		n, err := pruneChunk(db.WithContext(ctx), cutoff)
		total += n
		if err != nil || n < PruneChunk {
			return total, err
		}
	}
}

func pruneChunk(db *gorm.DB, cutoff time.Time) (int64, error) {
	var removed int64
	err := db.Transaction(func(tx *gorm.DB) error {
		if _, err := lockAndReadHead(tx); err != nil {
			return err
		}
		// The chunk ends at its last entry, so it is deleted as an index range
		// rather than a list of 5,000 ids.
		var last models.ActivityLog
		if err := tx.Where("created_at < ?", cutoff).
			Order("created_at asc, id asc").
			Offset(PruneChunk - 1).Limit(1).
			Find(&last).Error; err != nil {
			return fmt.Errorf("finding the end of the chunk: %w", err)
		}
		del := tx.Where("created_at < ?", cutoff)
		if last.ID != "" {
			del = del.Where("created_at < ? OR (created_at = ? AND id <= ?)", last.CreatedAt, last.CreatedAt, last.ID)
		}
		res := del.Delete(&models.ActivityLog{})
		if res.Error != nil {
			return fmt.Errorf("deleting entries: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return nil
		}
		var first models.ActivityLog
		if err := tx.Order("created_at asc, id asc").Limit(1).Find(&first).Error; err != nil {
			return fmt.Errorf("reading the oldest remaining entry: %w", err)
		}
		removed = res.RowsAffected
		return appendBatch(tx, []models.ActivityLog{{
			Method:      "SECURITY",
			Path:        prunePath,
			Status:      200,
			Resource:    "activity_logs",
			ResourceIDs: first.PrevHash,
			RecordCount: int(removed),
		}})
	})
	if err != nil {
		return 0, err
	}
	return removed, nil
}

// pruneRecorded reports whether a prune recorded that the log starts after the
// entry whose hash is prevHash.
func pruneRecorded(ctx context.Context, db *gorm.DB, prevHash string) (bool, error) {
	var n int64
	err := db.WithContext(ctx).Model(&models.ActivityLog{}).
		Where("method = ? AND path = ? AND resource_ids = ?", "SECURITY", prunePath, prevHash).
		Count(&n).Error
	return n > 0, err
}

// ErrChainIntact is returned by Reseal when the chain verifies.
var ErrChainIntact = errors.New("the chain verifies: there is nothing to reseal")

// ErrNotTheBreak is returned by Reseal when the entry named is not the first
// one that fails verification.
var ErrNotTheBreak = errors.New("that is not the first entry that fails verification")

// Reseal recomputes the chain from its first bad entry, and records that it did.
//
// For a log that fails through no one's tampering: before v3.215.0 every entry
// was hashed with a timestamp finer than Postgres and MySQL store, so those
// chains fail on their first row and always will. A reseal trusts the rows as
// they stand now, which is exactly what makes it dangerous, so it is never
// automatic. The caller names the first bad entry, which is checked against a
// fresh verification, and the reseal appends a SECURITY entry naming who did
// it, from which entry, how many it covered, and a digest of every hash it
// replaced. Changing that entry afterwards breaks the chain like any other.
func Reseal(ctx context.Context, db *gorm.DB, fromID, userID, ip, userAgent string) (int, error) {
	status, err := VerifyChain(ctx, db)
	if err != nil {
		return 0, err
	}
	if status.Valid {
		return 0, ErrChainIntact
	}
	if status.BrokenAtID != fromID {
		return 0, fmt.Errorf("%w: the chain first fails at %s", ErrNotTheBreak, status.BrokenAtID)
	}

	resealed := 0
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := lockAndReadHead(tx); err != nil {
			return err
		}
		var from models.ActivityLog
		if err := tx.First(&from, "id = ?", fromID).Error; err != nil {
			return fmt.Errorf("loading %s: %w", fromID, err)
		}
		// The entry before the break is the last one that verified.
		var before models.ActivityLog
		if err := tx.Where("(created_at, id) < (?, ?)", from.CreatedAt, from.ID).
			Order("created_at desc, id desc").Limit(1).Find(&before).Error; err != nil {
			return fmt.Errorf("loading the entry before %s: %w", fromID, err)
		}

		prevHash := before.Hash
		replaced := sha256.New()
		cond, args := "(created_at, id) >= (?, ?)", []interface{}{from.CreatedAt, from.ID}
		for {
			var batch []models.ActivityLog
			if err := tx.Where(cond, args...).Order("created_at asc, id asc").
				Limit(verifyBatchSize).Find(&batch).Error; err != nil {
				return err
			}
			for i := range batch {
				e := &batch[i]
				canonical, err := Canonical(e)
				if err != nil {
					return err
				}
				hash := ComputeHash(prevHash, canonical)
				replaced.Write([]byte(e.Hash))
				if err := tx.Model(&models.ActivityLog{}).Where("id = ?", e.ID).
					Updates(map[string]interface{}{"prev_hash": prevHash, "hash": hash}).Error; err != nil {
					return fmt.Errorf("resealing %s: %w", e.ID, err)
				}
				prevHash = hash
				resealed++
			}
			if len(batch) < verifyBatchSize {
				break
			}
			last := batch[len(batch)-1]
			cond, args = "(created_at, id) > (?, ?)", []interface{}{last.CreatedAt, last.ID}
		}

		return appendBatch(tx, []models.ActivityLog{{
			UserID:        userID,
			Method:        "SECURITY",
			Path:          "audit.chain.resealed",
			Status:        200,
			PayloadDigest: hex.EncodeToString(replaced.Sum(nil)),
			IPAddress:     ip,
			UserAgent:     userAgent,
			Resource:      "activity_logs",
			ResourceIDs:   fromID,
			RecordCount:   resealed,
		}})
	})
	if err != nil {
		return 0, err
	}
	return resealed, nil
}

// ChainStatus is the result of VerifyChain.
type ChainStatus struct {
	Valid        bool   `json:"valid"`
	TotalEntries int    `json:"total_entries"`
	BrokenAtID   string `json:"broken_at_id,omitempty"`
	BrokenAt     int    `json:"broken_at,omitempty"` // zero-indexed position
	Expected     string `json:"expected,omitempty"`
	Got          string `json:"got,omitempty"`
	Message      string `json:"message,omitempty"`
}

// VerifyChain walks the entire activity log in (created_at, id) order
// and recomputes every Hash. The first mismatch is reported with the
// position and offending row's ID — everything before that position
// is trustworthy.
//
// Memory-bounded: iterates in batches of verifyBatchSize so a 100M-row
// log doesn't OOM the process. Honours context cancellation so the
// caller can attach a deadline (the admin endpoint should pass
// c.Request.Context() with a 30s timeout).
//
// Cost is O(n) — about a second per million rows on a warm cache.
// Wire to a nightly cron + a /api/admin/activity/integrity endpoint.
const verifyBatchSize = 1000

func VerifyChain(ctx context.Context, db *gorm.DB) (ChainStatus, error) {
	prevHash := ""
	total := 0
	var lastCreatedAt time.Time
	var lastID string

	for {
		select {
		case <-ctx.Done():
			return ChainStatus{TotalEntries: total}, ctx.Err()
		default:
		}

		var batch []models.ActivityLog
		q := db.Order("created_at asc, id asc").Limit(verifyBatchSize)
		if total > 0 {
			// Cursor on (created_at, id) so we don't re-read rows
			// already verified in the previous batch.
			q = q.Where("(created_at, id) > (?, ?)", lastCreatedAt, lastID)
		}
		if err := q.Find(&batch).Error; err != nil {
			return ChainStatus{TotalEntries: total}, err
		}
		if len(batch) == 0 {
			break
		}

		for i := range batch {
			e := &batch[i]
			canonical, err := Canonical(e)
			if err != nil {
				return ChainStatus{TotalEntries: total}, err
			}
			if total+i == 0 && e.PrevHash != "" {
				// The oldest entry chains from one that is gone, which is what a
				// prune leaves. It verifies only if the prune recorded that hash.
				recorded, err := pruneRecorded(ctx, db, e.PrevHash)
				if err != nil {
					return ChainStatus{TotalEntries: total}, err
				}
				if !recorded {
					return ChainStatus{
						Valid:      false,
						BrokenAtID: e.ID,
						Got:        e.PrevHash,
						Message:    "the log starts after entries that were deleted, and no prune recorded deleting them",
					}, nil
				}
				prevHash = e.PrevHash
			}
			expected := ComputeHash(prevHash, canonical)
			if expected != e.Hash {
				return ChainStatus{
					Valid:        false,
					TotalEntries: total + i,
					BrokenAtID:   e.ID,
					BrokenAt:     total + i,
					Expected:     expected,
					Got:          e.Hash,
					Message:      "hash mismatch — row was modified, deleted, or inserted out of order",
				}, nil
			}
			if e.PrevHash != prevHash {
				return ChainStatus{
					Valid:        false,
					TotalEntries: total + i,
					BrokenAtID:   e.ID,
					BrokenAt:     total + i,
					Expected:     prevHash,
					Got:          e.PrevHash,
					Message:      "prev_hash mismatch — chain link broken",
				}, nil
			}
			prevHash = e.Hash
		}

		last := &batch[len(batch)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
		total += len(batch)

		if len(batch) < verifyBatchSize {
			break // last page
		}
	}

	return ChainStatus{
		Valid:        true,
		TotalEntries: total,
	}, nil
}
