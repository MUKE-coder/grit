// Package ids generates primary keys.
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
