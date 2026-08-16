package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// ensureSyncPolicySupport makes sure a project can actually honour a sync
// policy before the generator writes one into its routes.
//
// Two files are involved and they need opposite treatment. internal/sync/
// policy.go is new, so it is simply written when absent. internal/handlers/
// sync.go already exists and is where enforcement happens, so it is replaced
// only when the manifest proves nobody has edited it, and otherwise left alone
// with a warning.
//
// The failure being avoided is specific: a declared policy that nothing
// enforces. routes.go would say server_wins, the handler would keep returning
// VERSION_CONFLICT, and the app would prompt users about conflicts its author
// had explicitly said not to ask about.
func (g *Generator) ensureSyncPolicySupport() error {
	if g.Definition.Sync.GoLiteral() == "" {
		return nil // no policy declared, nothing to support
	}

	apiRoot := g.APIRoot()
	policyPath := filepath.Join(apiRoot, "internal", "sync", "policy.go")
	if !fileExists(policyPath) {
		if err := writeFileWithDirs(policyPath, strings.ReplaceAll(scaffold.APISyncPolicyGo(), "{{MODULE}}", g.Module)); err != nil {
			return fmt.Errorf("writing sync policy support: %w", err)
		}
		fmt.Println("  ✓ Added internal/sync/policy.go")
	}

	handlerPath := filepath.Join(apiRoot, "internal", "handlers", "sync.go")
	data, err := os.ReadFile(handlerPath)
	if err != nil {
		// No sync handler at all. The resource still generates; the doctor
		// reports the missing endpoint.
		return nil
	}
	// The route is checked independently of the handler. Gating it behind the
	// refresh meant a second run, where the handler was already current,
	// skipped the route entirely and left the policy undiscoverable.
	if strings.Contains(string(data), "PolicyFor") {
		g.ensurePolicyRoute()
		return nil
	}

	if g.refreshSyncHandler(handlerPath) {
		g.ensurePolicyRoute()
		return nil
	}

	yellow := color.New(color.FgHiYellow)
	yellow.Printf("\n  ⚠ %s predates sync policies and will not enforce the one you just declared.\n", filepath.Join("internal", "handlers", "sync.go"))
	fmt.Println("    Conflicts will stay manual whatever the policy says.")
	fmt.Println("    Replace that file with the current template, then run: grit sync doctor")
	fmt.Println()
	return nil
}

// refreshSyncHandler brings the sync handler forward so it enforces policies.
func (g *Generator) refreshSyncHandler(path string) bool {
	body := strings.ReplaceAll(scaffold.APISyncHandlerGo(), "{{MODULE}}", g.Module)
	if !g.refreshIfUnchanged(path, body) {
		return false
	}
	fmt.Println("  ✓ Updated internal/handlers/sync.go so it enforces sync policies")
	return true
}

// refreshIfUnchanged rewrites a file only when the manifest shows it is
// exactly what Grit last wrote there. Reports whether it did.
//
// An untracked file is not overwritten either. Untracked means the project
// predates the manifest, so there is no evidence in either direction, and the
// rule that shipped with the manifest was that no evidence is not permission.
func (g *Generator) refreshIfUnchanged(path, body string) bool {
	m, err := manifest.Load(g.Root)
	if err != nil {
		return false
	}
	rel, inside := manifest.Rel(g.Root, path)
	if !inside || m.StatusOf(g.Root, rel) != manifest.Unchanged {
		return false
	}
	return writeFileWithDirs(path, body) == nil
}

// ensurePolicyRoute mounts GET /sync/policy in a project whose routes.go
// predates it.
//
// Anchored on the pull route rather than a marker, because routes.go has no
// marker there and the pull route is the one line that must exist for any of
// this to be reachable. Idempotent: a project that already has the route is
// left alone.
func (g *Generator) ensurePolicyRoute() {
	path := filepath.Join(g.APIRoot(), "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(data)
	if strings.Contains(content, "syncHandler.Policy") {
		return
	}

	const anchor = "protected.GET(\"/sync/pull\", syncHandler.Pull)"
	idx := strings.Index(content, anchor)
	if idx < 0 {
		return // routes.go has been restructured; the doctor will report it
	}
	insertAt := idx + len(anchor)
	updated := content[:insertAt] +
		"\n\t\tprotected.GET(\"/sync/policy\", syncHandler.Policy)" +
		content[insertAt:]

	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Mounted GET /api/sync/policy")
}
