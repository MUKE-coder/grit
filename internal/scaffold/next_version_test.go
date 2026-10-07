package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// Next is pinned exactly, like React and pnpm.
//
// "^16.1.6" is any 16.x, so a project installed whatever had shipped that
// morning. One created in October 2026 got 16.4.0, three minors past the
// template, and landed on an App Router rewrite nothing in the project chose.
func TestNextIsPinnedExactly(t *testing.T) {
	opts := Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendNext, Theme: "emerald"}
	caret := regexp.MustCompile(`"next":\s*"[\^~]`)

	for name, pkg := range map[string]string{
		"web":   webPackageJSON(opts),
		"admin": adminPackageJSON(opts),
		"docs":  docsPackageJSON(opts),
	} {
		if caret.MatchString(pkg) {
			t.Errorf("the %s package.json lets Next float, so a project installs whatever is newest", name)
		}
		if !strings.Contains(pkg, `"next": "`+nextVersion+`"`) {
			t.Errorf("the %s package.json does not pin Next to %s", name, nextVersion)
		}
	}

	// A bare version, not a range: this is the thing being asserted.
	if strings.ContainsAny(nextVersion, "^~*x") {
		t.Errorf("nextVersion %q is a range", nextVersion)
	}
}
