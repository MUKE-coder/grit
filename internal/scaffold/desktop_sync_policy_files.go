package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
)

// The compiled-in sync policy table a desktop app has before it has ever
// reached the server.
//
// grit#214 asked for three things, in increasing order of how much the
// reporter wanted them: a typed error, a cached policy that survives a
// restart, and a sensible default for a client that has never reached the
// server. The third one is the one that decides whether this is usable,
// because the refusal only matters offline, and a client that fetches the
// policy at startup and gives up when the fetch fails has implemented the half
// that never fires.
//
// A compiled-in copy is the only way to refuse from install. The reporter
// maintained one by hand, with a test that parses the server's source and
// fails if the two lists disagree. That test is the right instinct and it
// should not have to exist: one generator can write both sides.
//
// So `grit sync` rewrites the bootstrap from the project's own registrations.
// A table absent from it is offline_first, which is why generating a resource
// does not break the desktop app before anybody runs grit sync.

// SyncPolicyBootstrapMarker is the anchor the generated map is written inside.
const SyncPolicyBootstrapMarker = "// grit:sync-policies"

// DesktopSyncPolicyBootstrapPath is where the bootstrap lives, relative to the
// desktop app root.
var DesktopSyncPolicyBootstrapPath = filepath.Join("sync", "policies_bootstrap.go")

// SyncPolicyDeclaration is one table's declared policy, as read from a
// project's sync registrations.
type SyncPolicyDeclaration struct {
	Table                string
	Mode                 string
	Conflict             string
	MaxOfflineAgeSeconds int

	// Explicit marks a declaration that came from a RegisterWithPolicy call
	// rather than from a bare Register. A table registered both ways has to
	// resolve to what somebody wrote, not to the default that happened to be
	// read first.
	Explicit bool
}

// ensureSyncPolicyBootstrap writes the bootstrap file if it is not there.
//
// Not in writeDesktopClientFiles' map, because that map is rewritten on every
// upgrade and this file's contents belong to the project: overwriting a
// generated policy table with an empty one would turn the guard off on exactly
// the projects that declared a policy.
func ensureSyncPolicyBootstrap(desktopRoot string) error {
	path := filepath.Join(desktopRoot, DesktopSyncPolicyBootstrapPath)
	if fileExists(path) {
		return nil
	}
	return writeFile(path, tmpl("desktop/sync/policies_bootstrap.go"))
}

// RenderSyncPolicyBootstrap builds the file's contents for a set of declared
// policies. Sorted by table so a regeneration with no change produces no diff.
func RenderSyncPolicyBootstrap(declared []SyncPolicyDeclaration) string {
	base := tmpl("desktop/sync/policies_bootstrap.go")

	sorted := make([]SyncPolicyDeclaration, len(declared))
	copy(sorted, declared)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Table < sorted[j].Table })

	var b strings.Builder
	for _, d := range sorted {
		// Only what a client acts on, and only when it is not the default.
		// A literal full of zero values reads as though somebody chose them.
		var fields []string
		if d.Mode != "" {
			fields = append(fields, fmt.Sprintf("Mode: %q", d.Mode))
		}
		if d.Conflict != "" {
			fields = append(fields, fmt.Sprintf("Conflict: %q", d.Conflict))
		}
		if d.MaxOfflineAgeSeconds > 0 {
			fields = append(fields, fmt.Sprintf("MaxOfflineAgeSeconds: %d", d.MaxOfflineAgeSeconds))
		}
		if len(fields) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("\t%q: {%s},\n", d.Table, strings.Join(fields, ", ")))
	}

	if b.Len() == 0 {
		return base
	}
	// The marker stays, so a later regeneration finds its place again.
	return strings.Replace(base,
		"\t"+SyncPolicyBootstrapMarker+"\n",
		"\t"+SyncPolicyBootstrapMarker+"\n"+b.String(), 1)
}

// WriteSyncPolicyBootstrap rewrites the bootstrap for a project that has a
// desktop app. It is a no-op for a project that does not.
//
// Returns the number of tables written, and whether the file changed, so the
// caller can say something true rather than claiming a write that did nothing.
func WriteSyncPolicyBootstrap(desktopRoot string, declared []SyncPolicyDeclaration) (int, bool, error) {
	path := filepath.Join(desktopRoot, DesktopSyncPolicyBootstrapPath)
	if !dirExists(filepath.Dir(path)) {
		return 0, false, nil
	}

	want := RenderSyncPolicyBootstrap(declared)
	// Compare against what writeFile will actually put on disk, not against
	// the template's own text. writeFile runs the Go formatter on the way out,
	// and gofmt aligns the keys of a map literal, so a byte comparison with
	// the unformatted string reports a change on every single run. Found by a
	// test asserting the second run was quiet.
	final := codefmt.File(path, want)
	existing, err := os.ReadFile(path)
	if err == nil && normaliseNewlines(string(existing)) == normaliseNewlines(final) {
		return countBootstrapEntries(final), false, nil
	}
	if err := writeFile(path, want); err != nil {
		return 0, false, fmt.Errorf("writing %s: %w", path, err)
	}
	return countBootstrapEntries(final), true, nil
}

// countBootstrapEntries counts the map entries in a rendered bootstrap.
func countBootstrapEntries(src string) int {
	n := 0
	inside := false
	for _, line := range strings.Split(normaliseNewlines(src), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "var bootstrapPolicies"):
			inside = true
		case inside && trimmed == "}":
			inside = false
		case inside && strings.HasPrefix(trimmed, `"`):
			n++
		}
	}
	return n
}

// normaliseNewlines makes a comparison independent of line endings, which on
// this platform is the difference between "no change" and a rewrite on every
// run.
func normaliseNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
