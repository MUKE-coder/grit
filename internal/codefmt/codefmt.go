// Package codefmt formats generated Go source before it lands on disk.
//
// Grit builds Go files by concatenating strings, which is fast to write and
// impossible to keep aligned by hand: struct tags drift out of column, import
// groups end up in whatever order the generator happened to append them, and
// every new field type is another chance to get it subtly wrong. Running the
// result through go/format makes that a non-issue — the generator only has to
// emit correct Go, not pretty Go.
package codefmt

import (
	"go/format"
	"path/filepath"
	"strings"
)

// Go formats Go source using the same algorithm as gofmt.
//
// Formatting is best-effort by design: if src does not parse, the original text
// is returned unchanged rather than an error. A generator that emits a syntax
// error should surface it at "go build" on the generated project, where the
// compiler points at the offending line — not here, where it would become an
// opaque scaffolding failure with nothing written to disk to inspect. Returning
// src unchanged preserves exactly the behaviour Grit had before formatting
// existed, so a parse failure can never be a regression.
func Go(src string) string {
	out, _ := GoChecked(src)
	return out
}

// GoChecked is Go, and says when src did not parse.
//
// Go discards that, which is right for a write: a file on disk that somebody
// can read beats a scaffolding failure with nothing written. It is wrong for a
// test. This function has already parsed the source and therefore already knows
// it is broken, and throwing that away was why a template emitting invalid Go
// was found by CI compiling a generated project ten minutes later, rather than
// by the package's own tests in a second (grit#159).
//
// On failure the returned string is src, unchanged, so a caller that ignores
// the error behaves exactly as Go does.
func GoChecked(src string) (string, error) {
	out, err := format.Source([]byte(src))
	if err != nil {
		return src, err
	}
	return string(out), nil
}

// File formats src when path names a Go file and returns it untouched
// otherwise, so callers can route every write through one helper instead of
// testing the extension at each site.
func File(path, src string) string {
	if !strings.EqualFold(filepath.Ext(path), ".go") {
		return src
	}
	return Go(src)
}
