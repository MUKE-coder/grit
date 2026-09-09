package generate

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// checkFileCollisions refuses to generate over a file the scaffold owns.
//
// reservedModels catches a name that collides with a built-in *model*. It does
// not catch a name that collides with a built-in *file*, and those are not the
// same set: the audit log's models are ActivityLog and UserActivity, both
// reserved, while the helpers that write to them live in services/activity.go,
// which "Activity" claims.
//
// So `grit generate resource Activity` reported success and overwrote twelve
// functions, including LogCreate, LogUpdate, LogDelete and DiffSummary, which
// every generated handler calls. The build then failed in access_review.go and
// event_subscribers.go, two files the developer had never opened, with nothing
// connecting the error to the command that caused it. Found by building a CRM,
// where "Activity" is what you call an interaction log.
//
// The manifest already records which files the scaffold wrote and which belong
// to a resource, so this asks it rather than growing a second hand-maintained
// list beside reservedModels. Whatever the scaffold gains next is protected on
// the day it lands.
func (g *Generator) checkFileCollisions(names Names) error {
	m, err := manifest.Load(g.Root)
	if err != nil || len(m.Files) == 0 {
		// No manifest means a project from before v3.147.0. Nothing to check
		// against, and refusing to generate would be worse than the risk.
		return nil
	}

	mine := "resource:" + names.Pascal
	var clashes []string
	seen := map[string]bool{}

	for _, abs := range g.targetFiles(names) {
		key, ok := manifest.Rel(g.Root, abs)
		if !ok {
			continue
		}
		entry, tracked := m.Files[key]
		if !tracked || entry.Generator == mine || seen[key] {
			continue
		}
		// Another resource's file is a different mistake with a different
		// message; this one is about the framework's own files.
		if strings.HasPrefix(entry.Generator, "resource:") {
			continue
		}
		seen[key] = true
		clashes = append(clashes, key)
	}
	if len(clashes) == 0 {
		return nil
	}
	sort.Strings(clashes)

	var b strings.Builder
	fmt.Fprintf(&b, "%q would overwrite %d file(s) the framework owns:\n\n", names.Pascal, len(clashes))
	for _, c := range clashes {
		fmt.Fprintf(&b, "  %s\n", c)
	}
	b.WriteString("\nThose files hold code other parts of the app call, so the build breaks\n")
	b.WriteString("somewhere else and the error does not mention this command. Pick another\n")
	b.WriteString("name:\n\n")
	fmt.Fprintf(&b, "  grit generate resource %s --fields \"...\"\n\n", suggestAlternative(names.Pascal))
	b.WriteString("If you really mean to replace them, pass --force.")
	return fmt.Errorf("%s", b.String())
}

// targetFiles lists the Go files a resource writes that the scaffold might
// already own.
//
// The API side only. A clash there stops the build, which is the failure worth
// refusing over; a frontend clash is visible in the app you are looking at.
func (g *Generator) targetFiles(names Names) []string {
	api := g.APIRoot()
	join := func(parts ...string) string {
		return filepath.Join(append([]string{api}, parts...)...)
	}
	return []string{
		join("internal", "models", names.Snake+".go"),
		join("internal", "services", names.Snake+".go"),
		join("internal", "handlers", names.Snake+".go"),
		join("internal", "handlers", names.Snake+"_import.go"),
		join("internal", "routes", names.Snake+"_routes.go"),
	}
}

// suggestAlternative offers a name that will not collide. Prefixing beats
// suffixing here: CrmActivity sorts next to nothing, ActivityRecord reads like
// a variant of the thing it must not be confused with.
func suggestAlternative(pascal string) string {
	for _, prefix := range []string{"Customer", "Team", "App"} {
		if !strings.HasPrefix(pascal, prefix) {
			return prefix + pascal
		}
	}
	return pascal + "Item"
}
