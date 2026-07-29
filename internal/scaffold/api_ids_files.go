package scaffold

// UUIDv7 identifier generation.
//
// Grit standardises on string UUID primary keys, which is what lets an offline
// desktop or mobile client mint an ID with no server contact — the property the
// sync engine is built on. UUIDv4 satisfies that but is fully random, which
// costs two things as tables grow:
//
//   - "the latest records" needs a separate created_at index, because primary
//     key order carries no meaning
//   - random inserts scatter across the B-tree, so index locality degrades and
//     more of the index has to stay hot in memory
//
// UUIDv7 fixes both while keeping the same 128-bit wire format and the same
// zero-coordination generation: the high bits are a millisecond timestamp, so
// ORDER BY id behaves like ORDER BY created_at. Nothing else changes — same
// column type, same GORM tags, same TypeScript and Zod types.
//
// It lives in its own package rather than in models so that services, jobs,
// handlers and the desktop sync engine can all reach it without importing
// models and risking an import cycle.

// apiIDsGo emits internal/ids/ids.go.
func apiIDsGo() string {
	return `// Package ids generates primary keys.
//
// Use ids.New() wherever a record needs an identifier. It returns a UUIDv7:
// still a standard 128-bit UUID that any client can generate offline with no
// coordination, but with a millisecond timestamp in the high bits, so IDs sort
// chronologically and insert with good index locality.
//
// Trade-off worth knowing: a v7 identifier encodes when it was created. If you
// expose raw IDs in public URLs, you are also publishing creation times. That
// is usually fine — and often useful — but if a particular resource must not
// leak that, give it its own opaque public token rather than switching the
// primary key back.
package ids

import "github.com/google/uuid"

// New returns a time-ordered UUIDv7 string.
//
// uuid.NewV7 only fails when the OS entropy source does, which is a situation
// where nothing else is going to work either. Rather than propagate an error
// into every BeforeCreate hook in the codebase, it falls back to a v4: an
// unordered identifier is a performance regression, an empty primary key is
// data corruption.
func New() string {
	if u, err := uuid.NewV7(); err == nil {
		return u.String()
	}
	return uuid.New().String()
}
`
}

// apiIDsTestGo emits internal/ids/ids_test.go.
func apiIDsTestGo() string {
	return `package ids

import (
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
)

// The whole point of v7: lexical order matches creation order, so ORDER BY id
// is ORDER BY created_at without a second index.
func TestNewIsTimeOrdered(t *testing.T) {
	var got []string
	for i := 0; i < 50; i++ {
		got = append(got, New())
		if i%10 == 0 {
			// Nudge past a millisecond boundary so ordering is observable
			// rather than dependent on within-tick counter behaviour.
			time.Sleep(2 * time.Millisecond)
		}
	}

	sorted := make([]string, len(got))
	copy(sorted, got)
	sort.Strings(sorted)

	for i := range got {
		if got[i] != sorted[i] {
			t.Fatalf("ids are not lexically time-ordered at index %d:\n  generated %s\n  sorted    %s", i, got[i], sorted[i])
		}
	}
}

// Still a valid UUID, and specifically version 7 — the column type, GORM tags
// and generated TypeScript all depend on it staying a standard UUID.
func TestNewIsValidUUIDv7(t *testing.T) {
	id := New()

	parsed, err := uuid.Parse(id)
	if err != nil {
		t.Fatalf("New() did not produce a parseable UUID: %v", err)
	}
	if v := parsed.Version(); v != 7 {
		t.Errorf("expected UUID version 7, got %d", v)
	}
	if len(id) != 36 {
		t.Errorf("expected the canonical 36-character form, got %d chars", len(id))
	}
}

// Uniqueness under a tight loop, which is where a naive time-based scheme
// collides: many IDs land inside the same millisecond.
func TestNewIsUniqueWithinAMillisecond(t *testing.T) {
	const n = 10000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		id := New()
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id generated within %d iterations: %s", i, id)
		}
		seen[id] = struct{}{}
	}
}

// An identifier must never come back empty — an empty primary key is worse
// than an unordered one, which is why the fallback exists.
func TestNewIsNeverEmpty(t *testing.T) {
	for i := 0; i < 100; i++ {
		if New() == "" {
			t.Fatal("New() returned an empty string")
		}
	}
}
`
}
