package main

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// upgradePassthrough hands the other binary the flags that change what the
// upgrade does. A flag added to `grit upgrade` and forgotten there is dropped
// in silence: `grit upgrade --to 3.455.0 --dry-run` would upgrade the project
// for real, which is the worst possible way for a forgotten line to show up.
//
// So this reads the command's own flag set rather than a list written here.
// A test with its own copy of the list would pass while the pair disagreed,
// which is the shape of bug this repository keeps finding.
func TestUpgradePassthroughCoversEveryFlag(t *testing.T) {
	// Every flag on the real command, and what it does about being passed on.
	handled := map[string]string{
		"force":   "passed through",
		"diff":    "passed through",
		"plan":    "passed through",
		"dry-run": "passed through",
		// --to is the flag being consumed. Passing it on would have the
		// fetched binary fetch itself.
		"to": "consumed by upgradeToVersion",
	}

	// With every boolean set, so a flag that is passed through shows up.
	all := upgradePassthrough(true, true, true, true)
	joined := strings.Join(all, " ")

	var missing []string
	upgradeCmd().Flags().VisitAll(func(f *pflag.Flag) {
		if _, known := handled[f.Name]; !known {
			missing = append(missing, f.Name)
			return
		}
		if handled[f.Name] != "passed through" {
			return
		}
		if !strings.Contains(joined, "--"+f.Name) {
			t.Errorf("--%s is a flag on `grit upgrade` and upgradePassthrough does not "+
				"pass it on, so `grit upgrade --to <v> --%s` silently ignores it",
				f.Name, f.Name)
		}
	})

	if len(missing) > 0 {
		t.Errorf("`grit upgrade` has flags this test has never heard of: %s\n\n"+
			"Decide what each one does when --to is given: passed through, or consumed "+
			"here, and add it to the map above. The failure mode for forgetting is "+
			"`--to <v> --dry-run` upgrading the project for real.",
			strings.Join(missing, ", "))
	}
}

// And nothing is passed on when nothing was asked for, so the other binary gets
// a bare `grit upgrade`.
func TestUpgradePassthroughIsEmptyByDefault(t *testing.T) {
	if got := upgradePassthrough(false, false, false, false); len(got) != 0 {
		t.Errorf("upgradePassthrough with no flags set = %v, want nothing", got)
	}
}

// --to naming the version this CLI already is should not download anything.
func TestUpgradeToThisVersionDoesNothing(t *testing.T) {
	// Both spellings, because somebody will write the v.
	for _, same := range []string{version, "v" + version} {
		if err := upgradeToVersion(same, nil); err != nil {
			t.Errorf("--to %s returned %v; it should fall through to the ordinary "+
				"upgrade rather than fetching this binary again", same, err)
		}
	}
}

func TestUpgradeToNeedsAVersion(t *testing.T) {
	for _, empty := range []string{"", "   ", "v"} {
		if err := upgradeToVersion(empty, nil); err == nil {
			t.Errorf("--to %q was accepted; it should say what it needs", empty)
		}
	}
}
