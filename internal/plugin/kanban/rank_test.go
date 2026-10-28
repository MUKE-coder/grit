package kanban

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
)

func TestFirstRank(t *testing.T) {
	if got := First(); got != "a0" {
		t.Fatalf("First() = %q, want %q", got, "a0")
	}
}

func TestBetweenEnds(t *testing.T) {
	first := First()

	before, err := Between("", first)
	if err != nil {
		t.Fatalf("Between(\"\", %q): %v", first, err)
	}
	if before >= first {
		t.Fatalf("Between(\"\", %q) = %q, which does not sort before it", first, before)
	}

	after, err := Between(first, "")
	if err != nil {
		t.Fatalf("Between(%q, \"\"): %v", first, err)
	}
	if after <= first {
		t.Fatalf("Between(%q, \"\") = %q, which does not sort after it", first, after)
	}

	mid, err := Between(before, after)
	if err != nil {
		t.Fatalf("Between(%q, %q): %v", before, after, err)
	}
	if !(before < mid && mid < after) {
		t.Fatalf("Between(%q, %q) = %q, which is not between them", before, after, mid)
	}
}

func TestBetweenRejectsBadOrder(t *testing.T) {
	a, b := "a1", "a0"
	if _, err := Between(a, b); err == nil {
		t.Fatalf("Between(%q, %q) succeeded; the arguments are in the wrong order", a, b)
	}
	if _, err := Between("a0", "a0"); err == nil {
		t.Fatal("Between with two equal ranks succeeded; there is no rank between a value and itself")
	}
}

func TestValidate(t *testing.T) {
	valid := []string{"a0", "a1", "Zz", "a0V"}
	for _, v := range valid {
		if err := Validate(v); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", v, err)
		}
	}

	invalid := map[string]string{
		"":                       "empty",
		"a":                      "integer part shorter than its head declares",
		"!0":                     "head outside the alphabet",
		"a0" + string(Digits[0]): "trailing zero",
		smallestInteger:          "the floor value itself",
	}
	for v, why := range invalid {
		if err := Validate(v); err == nil {
			t.Errorf("Validate(%q) = nil, want an error (%s)", v, why)
		}
	}
}

// TestAppendStaysShort is the claim the integer part exists to make true.
// Adding a card to the bottom of a column is the most common write a board
// takes, and a scheme that grows the key on every append would reach the
// column width long before a real board got large.
func TestAppendStaysShort(t *testing.T) {
	const appends = 10000

	rank := First()
	longest := len(rank)
	prev := rank
	for i := 0; i < appends; i++ {
		next, err := Between(rank, "")
		if err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
		if next <= prev {
			t.Fatalf("append %d produced %q, which does not sort after %q", i, next, prev)
		}
		prev = next
		rank = next
		if len(rank) > longest {
			longest = len(rank)
		}
	}

	if longest > 8 {
		t.Fatalf("after %d appends the longest rank was %d characters (%q); the integer part should have kept it short", appends, longest, rank)
	}
	t.Logf("%d appends, longest rank %d characters (%q)", appends, longest, rank)
}

// TestInsertBetweenSamePairAlwaysFits is the other half: dropping a card into
// the same slot over and over must keep working. Integer positions run out of
// gap here; a string rank grows one character at a time and never runs out.
func TestInsertBetweenSamePairAlwaysFits(t *testing.T) {
	lo := First()
	hi, err := Between(lo, "")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 500; i++ {
		mid, err := Between(lo, hi)
		if err != nil {
			t.Fatalf("insert %d between %q and %q: %v", i, lo, hi, err)
		}
		if !(lo < mid && mid < hi) {
			t.Fatalf("insert %d produced %q, which is not between %q and %q", i, mid, lo, hi)
		}
		if err := Validate(mid); err != nil {
			t.Fatalf("insert %d produced an invalid rank %q: %v", i, mid, err)
		}
		// Keep squeezing into the lower half, the worst case for the scheme.
		hi = mid
	}
}

// TestRandomInsertsKeepOrder is the property that matters to a user: whatever
// sequence of drags they perform, reading the column back by rank gives the
// order they arranged on screen.
func TestRandomInsertsKeepOrder(t *testing.T) {
	rng := rand.New(rand.NewSource(20261028))

	for trial := 0; trial < 50; trial++ {
		// ranks[i] is the rank of the card at visible position i.
		var ranks []string

		for step := 0; step < 120; step++ {
			at := 0
			if len(ranks) > 0 {
				at = rng.Intn(len(ranks) + 1)
			}
			prev, next := "", ""
			if at > 0 {
				prev = ranks[at-1]
			}
			if at < len(ranks) {
				next = ranks[at]
			}

			r, err := Between(prev, next)
			if err != nil {
				t.Fatalf("trial %d step %d: inserting at %d of %d: %v", trial, step, at, len(ranks), err)
			}
			if err := Validate(r); err != nil {
				t.Fatalf("trial %d step %d: produced invalid rank %q: %v", trial, step, r, err)
			}

			ranks = append(ranks, "")
			copy(ranks[at+1:], ranks[at:])
			ranks[at] = r

			if !sort.StringsAreSorted(ranks) {
				t.Fatalf("trial %d step %d: inserting at %d made the ranks unsorted: %v", trial, step, at, ranks)
			}
			for i := 1; i < len(ranks); i++ {
				if ranks[i-1] == ranks[i] {
					t.Fatalf("trial %d step %d: duplicate rank %q at %d", trial, step, ranks[i], i)
				}
			}
		}
	}
}

// TestMovesKeepOrder exercises the actual board operation: take a card out of
// one position and drop it into another, which is what a drag does.
func TestMovesKeepOrder(t *testing.T) {
	rng := rand.New(rand.NewSource(7))

	// Seed a column of 12 cards, identified by the label they keep across moves.
	type card struct {
		label string
		rank  string
	}
	var cards []card
	prev := ""
	for i := 0; i < 12; i++ {
		r, err := Between(prev, "")
		if err != nil {
			t.Fatal(err)
		}
		cards = append(cards, card{label: string(rune('A' + i)), rank: r})
		prev = r
	}

	for step := 0; step < 400; step++ {
		from := rng.Intn(len(cards))
		moved := cards[from]
		rest := append(append([]card{}, cards[:from]...), cards[from+1:]...)
		to := rng.Intn(len(rest) + 1)

		p, n := "", ""
		if to > 0 {
			p = rest[to-1].rank
		}
		if to < len(rest) {
			n = rest[to].rank
		}
		r, err := Between(p, n)
		if err != nil {
			t.Fatalf("step %d: moving %s from %d to %d: %v", step, moved.label, from, to, err)
		}
		moved.rank = r

		cards = append(rest[:to:to], append([]card{moved}, rest[to:]...)...)

		// The invariant: sorting by rank reproduces the on-screen order.
		want := make([]string, len(cards))
		for i, c := range cards {
			want[i] = c.label
		}
		byRank := append([]card{}, cards...)
		sort.Slice(byRank, func(i, j int) bool { return byRank[i].rank < byRank[j].rank })
		got := make([]string, len(byRank))
		for i, c := range byRank {
			got[i] = c.label
		}
		if strings.Join(got, "") != strings.Join(want, "") {
			t.Fatalf("step %d: screen order %s but rank order %s", step, strings.Join(want, ""), strings.Join(got, ""))
		}
	}
}

func TestNAfter(t *testing.T) {
	got, err := NAfter("", "", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 {
		t.Fatalf("NAfter(.., 5) returned %d ranks", len(got))
	}
	if !sort.StringsAreSorted(got) {
		t.Fatalf("NAfter returned unsorted ranks: %v", got)
	}

	lo := First()
	hi, err := Between(lo, "")
	if err != nil {
		t.Fatal(err)
	}
	mid, err := NAfter(lo, hi, 4)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range mid {
		if !(lo < r && r < hi) {
			t.Fatalf("NAfter produced %q outside (%q, %q)", r, lo, hi)
		}
	}
	if !sort.StringsAreSorted(mid) {
		t.Fatalf("NAfter returned unsorted ranks: %v", mid)
	}

	if got, err := NAfter("", "", 0); err != nil || got != nil {
		t.Fatalf("NAfter(.., 0) = %v, %v; want nil, nil", got, err)
	}
}

func TestDigitsAreAscending(t *testing.T) {
	// The database sorts ranks with a byte comparison. If the alphabet were
	// not in byte order the two orders would disagree and every ORDER BY in
	// the plugin would be subtly wrong.
	for i := 1; i < len(Digits); i++ {
		if Digits[i-1] >= Digits[i] {
			t.Fatalf("Digits is not ascending at %d: %q then %q", i, string(Digits[i-1]), string(Digits[i]))
		}
	}
	if len(Digits) != 62 {
		t.Fatalf("Digits has %d characters, want 62", len(Digits))
	}
}
