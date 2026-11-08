package scaffold

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// grit#196: the desktop's ColumnFormat did not hold what the generator writes
// into it.
//
//	src/routes/app/ai_runs.index.tsx(32,33): Type '"money"' is not assignable
//	  to type 'ColumnFormat | undefined'.
//
// Ten errors across nine resources, and quiet until CI: vite builds fine, only
// tsc -b fails. The generator emits one format per field type and writes it
// into both the admin and the desktop, and the desktop had nine of the
// twenty-three.
//
// This is the two-frontends shape: the generator has one source of truth for
// the format and each frontend has its own union, so one of them silently
// stops accepting what the generator emits.

var tsUnionMember = regexp.MustCompile(`"([a-z-]+)"`)

func TestBothColumnFormatUnionsAgree(t *testing.T) {
	admin := columnFormatMembers(t, adminResourceTypes(), "export type ColumnFormat =")
	desktop := columnFormatMembers(t, desktopClientDataTable(), "export type ColumnFormat =")

	if len(admin) < 15 {
		t.Fatalf("found only %d members in the admin union; the regex has stopped matching", len(admin))
	}

	// Everything the admin accepts, the desktop accepts. The generator writes
	// the same format string into both.
	have := map[string]bool{}
	for _, m := range desktop {
		have[m] = true
	}
	var missing []string
	for _, m := range admin {
		if !have[m] {
			missing = append(missing, m)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Errorf("the desktop ColumnFormat is missing %s.\n\n"+
			"    The generator writes one format per field type into both frontends, so a\n"+
			"    resource using one of these fails tsc -b in apps/desktop while building\n"+
			"    fine in the admin. vite build does not notice, so it is quiet until CI.",
			strings.Join(missing, ", "))
	}
}

// And the formats the generator actually emits are in both, which is the thing
// the unions exist for. money is the one that started this: the desktop had
// currency, which is not a field type and not what the generator writes.
func TestEveryEmittedColumnFormatIsAccepted(t *testing.T) {
	desktop := columnFormatMembers(t, desktopClientDataTable(), "export type ColumnFormat =")
	have := map[string]bool{}
	for _, m := range desktop {
		have[m] = true
	}

	// One per field type whose ColumnFormat is not simply "text".
	for _, format := range []string{
		"money", "percent", "json", "country", "rating", "time",
		"boolean", "date", "relative", "badge", "email", "link", "image",
	} {
		if !have[format] {
			t.Errorf("the generator emits format %q and the desktop union does not accept it", format)
		}
	}
}

// columnFormatMembers reads the string literals out of a TypeScript union.
func columnFormatMembers(t *testing.T, src, decl string) []string {
	t.Helper()
	at := strings.Index(src, decl)
	if at == -1 {
		t.Fatalf("no %q in the source", decl)
	}
	body := src[at+len(decl):]
	// To the semicolon that ends the type alias.
	if end := strings.Index(body, ";"); end != -1 {
		body = body[:end]
	}
	var out []string
	for _, m := range tsUnionMember.FindAllStringSubmatch(body, -1) {
		out = append(out, m[1])
	}
	return out
}
