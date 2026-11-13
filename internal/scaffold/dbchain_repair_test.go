package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The repair has to find its anchor in the file the generator actually
// writes, not in a fixture written to match the repair.
//
// A hand-written fixture agrees with whatever the repair looks for, including
// an anchor the generator stopped emitting a year ago, and the repair then
// no-ops on every real project while this passes.
func TestTheConnectRepairFindsItsAnchorInTheRealFile(t *testing.T) {
	src := strings.ReplaceAll(apiDatabaseGo(), "{{MODULE}}", "myapp/apps/api")
	// The generated file already has the call: that is the other writer, and
	// the two have to agree. Take it out to get the shape of an older project.
	before := removeApplyCall(t, src)
	if strings.Contains(before, "dbchain.Apply(db)") {
		t.Fatal("the test could not produce a file without the call")
	}

	out, fixed, warn := repairDBChainConnect(before, "myapp/apps/api")
	if len(warn) > 0 {
		t.Fatalf("the repair warned instead of repairing: %v", warn)
	}
	if len(fixed) == 0 {
		t.Fatal("the repair reported no change")
	}
	if !strings.Contains(out, "dbchain.Apply(db)") {
		t.Fatal("the repair did not add the call")
	}
	if !strings.Contains(out, `"myapp/apps/api/internal/dbchain"`) {
		t.Fatal("the repair added the call without the import")
	}
	// Byte for byte what the generator writes, which is the only version of
	// this assertion worth anything. An upgraded project that merely has the
	// call somewhere, with the import in the standard-library group, differs
	// from a generated one in a way nothing else looks at: gofmt sorts within
	// a group and never moves a line between groups.
	if out != src {
		wantAt, gotAt := firstDiff(src, out)
		t.Errorf("the repaired database.go is not what the template holds.\nwant:\n%s\ngot:\n%s", wantAt, gotAt)
	}
	// Twice is once.
	again, _, _ := repairDBChainConnect(out, "myapp/apps/api")
	if again != out {
		t.Error("running the repair twice changed the file again")
	}
}

// Apply has to be the last thing Connect installs.
//
// Not a style point. Apply registers every link with GORM, and GORM
// silently ignores an ordering constraint naming a callback it has not
// seen yet. With Apply before the field-type install, a link anchored to
// fieldtypes:create was registered against nothing and ran wherever GORM
// happened to put it, which is precisely what the Anchor check exists to
// prevent.
//
// Found by a test about something else: the byte comparison between the
// field-types repair and the template. This asserts the thing itself.
func TestApplyIsTheLastInstallInConnect(t *testing.T) {
	src := strings.ReplaceAll(apiDatabaseGo(), "{{MODULE}}", "myapp/apps/api")
	apply := strings.Index(src, "dbchain.Apply(db)")
	if apply < 0 {
		t.Fatal("Connect does not apply the chain at all")
	}
	for _, install := range []string{
		"crypto.Install(db)",
		"sanitize.Install(db)",
		"appendonly.Install(db)",
		"fieldtypes.Install(db)",
	} {
		at := strings.Index(src, install)
		if at < 0 {
			t.Errorf("Connect no longer calls %s, so this test is reading for the wrong thing", install)
			continue
		}
		if at > apply {
			t.Errorf("%s is installed after dbchain.Apply, so a link anchored to its "+
				"callback is registered before that callback exists and GORM ignores "+
				"the ordering", install)
		}
	}
}

func TestTheRequestIDRepairFindsItsAnchorInTheRealFile(t *testing.T) {
	src := strings.ReplaceAll(apiLoggerMiddlewareGo(Options{}), "{{MODULE}}", "myapp/apps/api")
	before := strings.Replace(src, requestContextLine, "", 1)
	if strings.Contains(before, "dbchain.WithRequest(") {
		t.Fatalf("the test could not produce a file without the call:\n%s", before)
	}

	out, fixed, warn := repairDBChainRequestID(before, "myapp/apps/api")
	if len(warn) > 0 {
		t.Fatalf("the repair warned instead of repairing: %v", warn)
	}
	if len(fixed) == 0 {
		t.Fatal("the repair reported no change")
	}
	if !strings.Contains(out, "dbchain.WithRequest(") {
		t.Fatal("the repair did not put the request on the context")
	}
	if !strings.Contains(out, `"myapp/apps/api/internal/dbchain"`) {
		t.Fatal("the repair added the call without the import")
	}
	again, _, _ := repairDBChainRequestID(out, "myapp/apps/api")
	if again != out {
		t.Error("running the repair twice changed the file again")
	}
}

// A file that has neither the anchor nor the call is reported, not silently
// left: a repair that warns is a repair somebody can act on, and one that
// returns quietly is a project that will never get the fix.
func TestTheRepairsWarnWhenTheyCannotFindTheAnchor(t *testing.T) {
	if _, _, warn := repairDBChainConnect("package database\n", "m"); len(warn) == 0 {
		t.Error("the Connect repair was silent about a file it could not repair")
	}
	if _, _, warn := repairDBChainRequestID("package middleware\n", "m"); len(warn) == 0 {
		t.Error("the RequestID repair was silent about a file it could not repair")
	}
}

func TestTheQueryLogIsGitIgnored(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	if err := os.WriteFile(path, []byte("node_modules\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := ignoreQueryLog(root); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(raw), ".grit/queries.jsonl"); n != 1 {
		t.Errorf("the log is ignored %d times, want 1:\n%s", n, raw)
	}
	if !strings.Contains(string(raw), "node_modules") {
		t.Error("the repair lost what was already in the file")
	}
}

// removeApplyCall takes the block back out, so the test starts from the shape
// of a project generated before it existed.
func removeApplyCall(t *testing.T, src string) string {
	t.Helper()
	out := strings.Replace(src, dbChainConnectHook, "", 1)
	if out == src {
		t.Fatalf("dbChainConnectHook is not in the generated database.go verbatim, so the two "+
			"writers have drifted and the repair is reading for something the "+
			"generator no longer writes:\n%s", dbChainConnectHook)
	}
	return out
}

// firstDiff is the two files from the first byte they differ at, so a
// failure shows the difference rather than two thousand identical lines.
func firstDiff(want, got string) (string, string) {
	n := min(len(want), len(got))
	for i := 0; i < n; i++ {
		if want[i] != got[i] {
			from := max(0, i-200)
			return want[from:min(len(want), i+300)], got[from:min(len(got), i+300)]
		}
	}
	return want[n:], got[n:]
}
