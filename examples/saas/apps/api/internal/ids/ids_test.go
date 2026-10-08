package ids

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
