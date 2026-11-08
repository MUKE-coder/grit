package scaffold

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The two govulncheck allowlists have to agree.
//
// .github/workflows/scan.yml scaffolds a project and runs govulncheck over it,
// against THIS repository's .github/govulncheck-allow.txt. A generated project
// ships its own, from govulncheckAllowTXT. So the same finding reaches both,
// and an entry added to one leaves the other red.
//
// That is what happened: v3.490.0 shipped the crewjam/saml entry in the
// template, which turned a new project's CI green, and the Scan workflow here
// stayed red on the identical advisory. It had been red since v3.486.0 and
// nobody was reading it, which is the same way grit#190 survived 30 releases.
//
// The direction matters. Everything the template accepts must also be accepted
// here, because this job scans the thing the template produces. The reverse is
// not required: this repository's own code may carry something a generated app
// does not.
var osvID = regexp.MustCompile(`(?m)^(GO-[0-9]{4}-[0-9]+)`)

func TestTheTwoAllowlistsAgree(t *testing.T) {
	shipped := osvIDs(govulncheckAllowTXT())

	raw, err := os.ReadFile("../../.github/govulncheck-allow.txt")
	if err != nil {
		t.Fatalf("reading this repository's allowlist: %v", err)
	}
	ours := osvIDs(string(raw))

	for _, id := range shipped {
		var found bool
		for _, m := range ours {
			if m == id {
				found = true
			}
		}
		if !found {
			t.Errorf("a generated project accepts %s and .github/govulncheck-allow.txt does not, "+
				"so the Scan workflow fails on a finding this repository has already decided to "+
				"accept. Add the same line, with the same reasoning.", id)
		}
	}

	// And neither file may carry a bare id: the format is the reason and the
	// issue, because an allowlist entry nobody can justify is one nobody will
	// remove when the fix lands.
	for _, src := range []struct {
		name string
		body string
	}{
		{"the template's allowlist", govulncheckAllowTXT()},
		{".github/govulncheck-allow.txt", string(raw)},
	} {
		for _, line := range strings.Split(src.body, "\n") {
			if !strings.HasPrefix(line, "GO-") {
				continue
			}
			if strings.Count(line, "|") < 2 {
				t.Errorf("%s has an entry with no reason and no issue: %q", src.name, line)
			}
		}
	}
}

// osvIDs is every advisory an allowlist accepts, in the form the workflow's own
// grep matches: anchored at the start of a line.
func osvIDs(body string) []string {
	var out []string
	for _, m := range osvID.FindAllStringSubmatch(body, -1) {
		out = append(out, m[1])
	}
	sort.Strings(out)
	return out
}
