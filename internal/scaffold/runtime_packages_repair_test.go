package scaffold

import (
	"strings"
	"testing"
)

// A Dockerfile from before the marker gets it where a new one has it, and a
// new one is left alone.
func TestRuntimePackagesMarkerRepairMatchesTheTemplates(t *testing.T) {
	for name, tmpl := range map[string]string{"api": dockerfileAPI(), "single": dockerfileSingle()} {
		if strings.Count(tmpl, "# grit:runtime-packages") != 1 {
			t.Fatalf("%s Dockerfile template has no single runtime-packages marker", name)
		}
		if out, changes, _ := repairRuntimePackagesMarkerSource(tmpl); out != tmpl || len(changes) != 0 {
			t.Errorf("%s: a new Dockerfile was changed", name)
		}
		old := strings.Replace(tmpl, "\n"+runtimePackagesMarker, "", 1)
		if strings.Contains(old, "grit:runtime-packages") {
			t.Fatalf("%s: could not remove the marker to make an old Dockerfile", name)
		}
		fixed, changes, warnings := repairRuntimePackagesMarkerSource(old)
		if fixed != tmpl || len(changes) != 1 || len(warnings) != 0 {
			t.Errorf("%s: the repair did not restore the template (changes %v, warnings %v)", name, changes, warnings)
		}
	}
	if _, _, warnings := repairRuntimePackagesMarkerSource("FROM scratch\n"); len(warnings) != 1 {
		t.Error("a Dockerfile Grit did not write should get a warning")
	}
}
