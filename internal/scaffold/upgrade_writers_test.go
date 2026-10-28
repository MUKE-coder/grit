package scaffold

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Two framework files that must agree, and the writer that must reach both.
//
// internal/services/security_autoban.go calls methods on SecObsBridge, and
// both files come from writeSecurityObservabilityFiles. That writer was called
// when scaffolding a new project and nowhere else, so an existing project kept
// its original bridge while other paths refreshed the sweep that calls it:
//
//	security_autoban.go:240:22: s.Bridge.SentinelPost undefined
//	(type *SecObsBridge has no field or method SentinelPost)
//
// Two checks, because the bug needed both halves to go wrong. First that the
// templates agree with each other; second that upgrade delivers them.
func TestTheBanSweepOnlyCallsBridgeMethodsThatExist(t *testing.T) {
	src := readScaffoldSource(t, "api_security_observability_files.go")

	called := map[string]bool{}
	for _, m := range regexp.MustCompile(`Bridge\.([A-Z][A-Za-z0-9]*)\(`).FindAllStringSubmatch(src, -1) {
		called[m[1]] = true
	}
	if len(called) == 0 {
		t.Fatal("found no Bridge calls at all; this reader has gone stale, fix it rather than deleting it")
	}

	defined := map[string]bool{}
	for _, m := range regexp.MustCompile(`func \(b \*SecObsBridge\) ([A-Za-z0-9]*)\(`).FindAllStringSubmatch(src, -1) {
		defined[m[1]] = true
	}

	var missing []string
	for name := range called {
		if !defined[name] {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Errorf("the ban sweep calls SecObsBridge.%s, which the bridge template does not define",
			strings.Join(missing, ", SecObsBridge."))
	}
}

// TestUpgradeDeliversEveryFrameworkWriter.
//
// writeTOTPFiles and writeSecurityObservabilityFiles were called from
// scaffold.go only, so every file they own was frozen at whatever version a
// project was created with: the two-factor models and service, the security
// and observability handlers, the Sentinel/Pulse bridge, its poller, and the
// auto-ban sweep. New projects had them; no upgrade ever delivered one.
//
// Named rather than derived, because deriving it needs an allowlist of the
// twenty-odd writers the upgrade reaches by another path, and an allowlist
// nobody maintains is a test that passes for the wrong reason.
func TestUpgradeDeliversEveryFrameworkWriter(t *testing.T) {
	up := readScaffoldSource(t, "upgrade.go")

	for _, writer := range []struct {
		fn   string
		owns string
	}{
		{"writeTOTPFiles", "the two-factor models and service"},
		{"writeSecurityObservabilityFiles", "the Sentinel/Pulse bridge, its poller and the ban sweep"},
	} {
		if !strings.Contains(up, writer.fn+"(root, opts)") {
			t.Errorf("upgrade.go never calls %s, so %s stay at whatever version "+
				"a project was scaffolded with and the framework code that uses them does not compile",
				writer.fn, writer.owns)
		}
	}
}
