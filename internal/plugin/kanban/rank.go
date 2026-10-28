// Package kanban orders rows by a string key, so that moving one row writes
// one row.
//
// The obvious way to order a board is an integer position column. It is also
// wrong, and wrong in a way that only shows up with two people on the board:
//
//	UPDATE cards SET position = position + 1 WHERE column_id = ? AND position >= ?
//
// That rewrites every card below the insertion point. Two people dragging at
// once interleave those rewrites and the board ends up with duplicate or
// skipped positions, and on Postgres the two statements can deadlock outright
// because they touch the same rows in opposite orders. Gap numbering (100,
// 200, 300) only delays it: the gap closes after a few drags between the same
// pair and you are renumbering again.
//
// So a rank here is a string, compared lexicographically, and there is always
// a string strictly between any two others. Moving a card is:
//
//	UPDATE cards SET column_id = ?, rank = ? WHERE id = ?
//
// One row. No renumbering, no transaction spanning other people's cards, and
// no lock held while a user decides where to drop. The same ORDER BY works
// unchanged on Postgres, MySQL and SQLite, which matters because Grit targets
// all three and a scheme that needs dialect-specific SQL is three bugs.
//
// Two people who drop a card in the same slot at the same time compute the
// same rank. That is a tie, not corruption: order by (rank, id) and the board
// is stable and identical for both of them. Every query in this package's
// callers does exactly that.
//
// The algorithm is fractional indexing, as described by Figma and implemented
// in the `fractional-indexing` npm package. A key has an integer part whose
// first character encodes its own length, then an optional fractional part.
// The integer part is what keeps appending cheap: adding to the end of a
// column increments an integer rather than growing the string, so a thousand
// appends leave the key short. Only repeated insertion between the same two
// neighbours lengthens a key, by one character at a time.
package kanban

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Digits is the rank alphabet. It must be in ascending byte order, because
// the database sorts these strings with a plain byte comparison and the
// algorithm assumes the two orders agree.
const Digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// smallestInteger is the floor of the integer part. Reaching it means somebody
// has prepended to the same list about 62^26 times.
var smallestInteger = "A" + strings.Repeat(string(Digits[0]), 26)

// ErrOutOfRange is returned when the integer part cannot move any further.
// It is unreachable in practice and reported rather than ignored so that if it
// ever does happen, it says so instead of silently producing a bad key.
var ErrOutOfRange = errors.New("kanban: rank space exhausted")

// Between returns a rank that sorts strictly after prev and strictly before
// next.
//
// An empty prev means "before every existing rank", an empty next means
// "after every existing rank", and both empty means the list is empty. A rank
// is never itself the empty string, so there is no ambiguity.
//
// Callers pass the ranks of the two cards the dragged card was dropped
// between, which is the information the client already has.
func Between(prev, next string) (string, error) {
	if prev != "" {
		if err := Validate(prev); err != nil {
			return "", err
		}
	}
	if next != "" {
		if err := Validate(next); err != nil {
			return "", err
		}
	}
	if prev != "" && next != "" && prev >= next {
		return "", fmt.Errorf("kanban: %q is not before %q", prev, next)
	}

	if prev == "" {
		if next == "" {
			return "a" + string(Digits[0]), nil
		}
		ib, err := integerPart(next)
		if err != nil {
			return "", err
		}
		fb := next[len(ib):]
		if ib == smallestInteger {
			mid, err := midpoint("", fb)
			if err != nil {
				return "", err
			}
			return ib + mid, nil
		}
		if ib < next {
			return ib, nil
		}
		return decrementInteger(ib)
	}

	if next == "" {
		ia, err := integerPart(prev)
		if err != nil {
			return "", err
		}
		fa := prev[len(ia):]
		i, err := incrementInteger(ia)
		if errors.Is(err, ErrOutOfRange) {
			mid, mErr := midpoint(fa, "")
			if mErr != nil {
				return "", mErr
			}
			return ia + mid, nil
		}
		if err != nil {
			return "", err
		}
		return i, nil
	}

	ia, err := integerPart(prev)
	if err != nil {
		return "", err
	}
	fa := prev[len(ia):]
	ib, err := integerPart(next)
	if err != nil {
		return "", err
	}
	fb := next[len(ib):]

	if ia == ib {
		mid, err := midpoint(fa, fb)
		if err != nil {
			return "", err
		}
		return ia + mid, nil
	}
	i, err := incrementInteger(ia)
	if err != nil {
		return "", err
	}
	if i < next {
		return i, nil
	}
	mid, err := midpoint(fa, "")
	if err != nil {
		return "", err
	}
	return ia + mid, nil
}

// First returns the rank for the first row in an empty list.
func First() string {
	r, err := Between("", "")
	if err != nil {
		// Between("", "") is a constant; an error here is impossible.
		panic(err)
	}
	return r
}

// NAfter returns n ranks in ascending order, all after prev and before next.
//
// Useful for seeding a board and for bulk import, where calling Between in a
// loop would otherwise need each previous result fed back in.
func NAfter(prev, next string, n int) ([]string, error) {
	if n <= 0 {
		return nil, nil
	}
	out := make([]string, 0, n)
	cur := prev
	for i := 0; i < n; i++ {
		r, err := Between(cur, next)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
		cur = r
	}
	return out, nil
}

// Validate reports whether a string is a well-formed rank. Ranks come back
// from the database and from clients, so they are checked rather than
// assumed: a malformed one would otherwise produce a key that sorts somewhere
// nobody intended.
func Validate(key string) error {
	if key == "" {
		return errors.New("kanban: rank is empty")
	}
	if key == smallestInteger {
		return fmt.Errorf("kanban: invalid rank %q", key)
	}
	i, err := integerPart(key)
	if err != nil {
		return err
	}
	if strings.HasSuffix(key[len(i):], string(Digits[0])) {
		// A trailing zero has two spellings for one position, and midpoint
		// relies on there being exactly one.
		return fmt.Errorf("kanban: invalid rank %q (trailing zero)", key)
	}
	return nil
}

// integerLength decodes the length an integer part declares in its first
// character. 'a' through 'z' count upward from 2, 'A' through 'Z' count
// downward from 'Z', so the integer part can grow in both directions without
// a separate length field.
func integerLength(head byte) (int, error) {
	switch {
	case head >= 'a' && head <= 'z':
		return int(head-'a') + 2, nil
	case head >= 'A' && head <= 'Z':
		return int('Z'-head) + 2, nil
	default:
		return 0, fmt.Errorf("kanban: invalid rank head %q", string(head))
	}
}

func integerPart(key string) (string, error) {
	if key == "" {
		return "", errors.New("kanban: rank is empty")
	}
	n, err := integerLength(key[0])
	if err != nil {
		return "", err
	}
	if n > len(key) {
		return "", fmt.Errorf("kanban: invalid rank %q", key)
	}
	return key[:n], nil
}

func validateInteger(i string) error {
	if i == "" {
		return errors.New("kanban: empty integer part")
	}
	n, err := integerLength(i[0])
	if err != nil {
		return err
	}
	if len(i) != n {
		return fmt.Errorf("kanban: invalid integer part %q", i)
	}
	return nil
}

func incrementInteger(x string) (string, error) {
	if err := validateInteger(x); err != nil {
		return "", err
	}
	head := x[0]
	digs := []byte(x[1:])
	carry := true
	for i := len(digs) - 1; carry && i >= 0; i-- {
		d := strings.IndexByte(Digits, digs[i]) + 1
		if d == len(Digits) {
			digs[i] = Digits[0]
		} else {
			digs[i] = Digits[d]
			carry = false
		}
	}
	if !carry {
		return string(head) + string(digs), nil
	}
	if head == 'Z' {
		return "a" + string(Digits[0]), nil
	}
	if head == 'z' {
		return "", ErrOutOfRange
	}
	h := head + 1
	if h > 'a' {
		digs = append(digs, Digits[0])
	} else {
		digs = digs[:len(digs)-1]
	}
	return string(h) + string(digs), nil
}

func decrementInteger(x string) (string, error) {
	if err := validateInteger(x); err != nil {
		return "", err
	}
	head := x[0]
	digs := []byte(x[1:])
	borrow := true
	for i := len(digs) - 1; borrow && i >= 0; i-- {
		d := strings.IndexByte(Digits, digs[i]) - 1
		if d == -1 {
			digs[i] = Digits[len(Digits)-1]
		} else {
			digs[i] = Digits[d]
			borrow = false
		}
	}
	if !borrow {
		return string(head) + string(digs), nil
	}
	if head == 'a' {
		return "Z" + string(Digits[len(Digits)-1]), nil
	}
	if head == 'A' {
		return "", ErrOutOfRange
	}
	h := head - 1
	if h < 'Z' {
		digs = append(digs, Digits[len(Digits)-1])
	} else {
		digs = digs[:len(digs)-1]
	}
	return string(h) + string(digs), nil
}

// midpoint returns a string strictly between the fractional parts a and b,
// where an empty b means "no upper bound". It is the part that guarantees
// there is always room between two ranks: when no shorter string fits, it
// appends a digit and tries again one place further along.
func midpoint(a, b string) (string, error) {
	zero := Digits[0]
	if b != "" && a >= b {
		return "", fmt.Errorf("kanban: %q is not before %q", a, b)
	}
	if strings.HasSuffix(a, string(zero)) || (b != "" && strings.HasSuffix(b, string(zero))) {
		return "", fmt.Errorf("kanban: trailing zero in %q or %q", a, b)
	}

	if b != "" {
		// Skip the common prefix and recurse on what differs, so the result
		// keeps the shared leading characters.
		n := 0
		for n < len(b) {
			ca := zero
			if n < len(a) {
				ca = a[n]
			}
			if ca != b[n] {
				break
			}
			n++
		}
		if n > 0 {
			var tail string
			if n < len(a) {
				tail = a[n:]
			}
			mid, err := midpoint(tail, b[n:])
			if err != nil {
				return "", err
			}
			return b[:n] + mid, nil
		}
	}

	digitA := 0
	if a != "" {
		digitA = strings.IndexByte(Digits, a[0])
		if digitA < 0 {
			return "", fmt.Errorf("kanban: invalid digit %q", string(a[0]))
		}
	}
	digitB := len(Digits)
	if b != "" {
		digitB = strings.IndexByte(Digits, b[0])
		if digitB < 0 {
			return "", fmt.Errorf("kanban: invalid digit %q", string(b[0]))
		}
	}

	if digitB-digitA > 1 {
		mid := int(math.Round(0.5 * float64(digitA+digitB)))
		return string(Digits[mid]), nil
	}
	if b != "" && len(b) > 1 {
		return b[:1], nil
	}
	// The two first digits are adjacent, so nothing fits between them at this
	// position. Keep a's digit and find room one character further along.
	var tail string
	if len(a) > 0 {
		tail = a[1:]
	}
	mid, err := midpoint(tail, "")
	if err != nil {
		return "", err
	}
	return string(Digits[digitA]) + mid, nil
}
