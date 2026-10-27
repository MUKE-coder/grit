package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Sentinel filters threats on the exact stored word, and its words are
// capitalised: "Critical", "High", "Medium", "Low". Everything on the security
// dashboard that names a severity has to spell it that way.
//
// Getting it wrong is silent in the worst way. A sweep asking for "high" gets
// an empty page, which is indistinguishable from a quiet week, so the job runs
// on schedule, bans nobody, and logs that nobody was over the threshold. The
// first version of the sweep did exactly that against a log holding eleven
// high-severity threats.
//
// This reads the words out of Sentinel's own source in the module cache when
// it is there, and falls back to the known values when it is not, so the test
// still means something on a machine that has not downloaded the module.
func TestTheSeverityWordsMatchSentinel(t *testing.T) {
	want := sentinelSeverityWords(t)

	src := securityAutoBanGo()
	for _, word := range want {
		decl := `= "` + word + `"`
		if !strings.Contains(src, decl) {
			t.Errorf("the severity constants do not carry %q as Sentinel spells it: "+
				"a filter asking for the wrong case matches nothing and reads as a quiet week", word)
		}
	}

	// And the sweep counts the two that mean somebody is attacking.
	if !strings.Contains(src, "autoBanSeverities = []string{SeverityCritical, SeverityHigh}") {
		t.Error("the sweep no longer counts Critical and High through the constants")
	}
	// A literal lowercase severity in a query is the bug this test exists for.
	for _, bad := range []string{`severity=high`, `severity=critical`, `"high", "critical"`, `{"critical", "high"}`} {
		if strings.Contains(src, bad) {
			t.Errorf("a lowercase severity literal is back (%q): Sentinel matches nothing on it", bad)
		}
	}
}

// The handler maps whatever the filter tabs send onto those words rather than
// passing them through, because the tabs are lowercase and the store is not.
func TestTheThreatFilterMapsSeverityRatherThanPassingItThrough(t *testing.T) {
	src := securityHandlerGo()
	if !strings.Contains(src, "services.CanonicalSeverity(severity)") {
		t.Error("the threat log sends the filter's own spelling to Sentinel, " +
			"so every tab but All returns an empty page")
	}
	if strings.Contains(src, `query += "&severity=" + severity`) {
		t.Error("the raw query value is forwarded to Sentinel unmapped")
	}
}

var severityDecl = regexp.MustCompile(`Severity(Low|Medium|High|Critical)\s+Severity\s*=\s*"([A-Za-z]+)"`)

// sentinelSeverityWords reads the constants from the Sentinel module in the
// local module cache. It is the real source of truth; the fallback is what it
// said when this was written.
func sentinelSeverityWords(t *testing.T) []string {
	t.Helper()
	fallback := []string{"Low", "Medium", "High", "Critical"}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fallback
		}
		gopath = filepath.Join(home, "go")
	}
	// The module cache escapes an upper-case letter as "!" + its lower case.
	matches, err := filepath.Glob(filepath.Join(gopath, "pkg", "mod",
		"github.com", "!m!u!k!e-coder", "sentinel", "v2@*", "core", "constants.go"))
	if err != nil || len(matches) == 0 {
		t.Log("Sentinel is not in the module cache; checking against the known words")
		return fallback
	}
	raw, err := os.ReadFile(matches[len(matches)-1])
	if err != nil {
		return fallback
	}
	var words []string
	for _, m := range severityDecl.FindAllStringSubmatch(string(raw), -1) {
		words = append(words, m[2])
	}
	if len(words) == 0 {
		return fallback
	}
	return words
}
