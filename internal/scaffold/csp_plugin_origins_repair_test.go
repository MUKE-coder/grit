package scaffold

import (
	"strings"
	"testing"
)

// A new config has the list; an old one gets it exactly where a new one has it.
func TestCSPPluginOriginsRepairMatchesTheTemplate(t *testing.T) {
	tmpl := nextSecurityHeaders()
	if !strings.Contains(tmpl, "// grit:csp-origins") || !strings.Contains(tmpl, cspFrameSrc) {
		t.Fatal("the Next.js CSP has no plugin origins")
	}
	if out, changes, _ := repairCSPPluginOriginsSource(tmpl); out != tmpl || len(changes) != 0 {
		t.Error("a current config was changed")
	}
	old := strings.Replace(tmpl, cspPluginOrigins, "", 1)
	old = strings.Replace(old, "  "+cspFrameSrc+"\n", "", 1)
	old = strings.Replace(old, cspJoinWithPlugins, cspJoinOld, 1)
	fixed, changes, warnings := repairCSPPluginOriginsSource(old)
	if fixed != tmpl || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("the repair did not restore the template (%v, %v)", changes, warnings)
	}
	if _, _, warnings := repairCSPPluginOriginsSource("const csp = [\"default-src *\"].join(\";\");\n"); len(warnings) != 1 {
		t.Error("a hand-written CSP should get a warning")
	}
}
