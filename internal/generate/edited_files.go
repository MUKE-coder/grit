package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// checkEditedFiles refuses to overwrite this resource's own files once someone
// has edited them.
//
// Re-running a generate is not an unusual thing to do. It is the obvious way to
// add a field: you already know the command, so you run it again with one more
// entry in --fields. Doing that silently rewrote the model, the service and the
// handler, printing a green tick for each, and everything hand-written in them
// was gone. No prompt, no backup, no mention in the output that anything had
// been replaced rather than created.
//
// Which is the opposite of what the documentation promises. /docs/concepts/
// generated-files calls apps/api/internal/services/<name>.go "your custom
// business logic goes here" and says of regeneration: "Everything above is your
// code once generated, edit models, add service methods, restyle screens
// freely. The one exception is packages/shared/types/*". So the file a
// newcomer is told is the safe place for business logic was the file that lost
// it, and the promise that it would not was in writing.
//
// Found on a first week with the framework: a BeforeSave rule on the model and
// a PinnedCount method on the service, both gone after one re-run.
//
// The manifest already records who wrote every file and a hash of what was
// written, so it can tell an edited file from an untouched one. That is the
// whole check. Untouched files regenerate exactly as before, which keeps the
// common case, a resource you have not customised yet, entirely unchanged.
func (g *Generator) checkEditedFiles(names Names) error {
	m, err := manifest.Load(g.Root)
	if err != nil || len(m.Files) == 0 {
		// A project from before the manifest. There is nothing to compare
		// against, and refusing every regenerate on a project we cannot reason
		// about would be worse than the risk.
		return nil
	}

	mine := "resource:" + names.Pascal
	var edited []string
	for _, key := range m.Entries() {
		if m.Files[key].Generator != mine {
			continue
		}
		// The overlay is the one file that is meant to be edited, and the one
		// file generation never rewrites. Flagging it would refuse the command
		// for doing exactly what it was designed for.
		if strings.HasSuffix(key, ".custom.tsx") {
			continue
		}
		if m.StatusOf(g.Root, key) == manifest.Modified {
			edited = append(edited, key)
		}
	}
	if len(edited) == 0 {
		return nil
	}
	sort.Strings(edited)

	var b strings.Builder
	fmt.Fprintf(&b, "%s has %d file(s) you have edited since Grit wrote them:\n\n", names.Pascal, len(edited))
	for _, key := range edited {
		fmt.Fprintf(&b, "  %s\n", key)
	}
	b.WriteString("\nRegenerating rewrites those files whole, so anything you added to them\n")
	b.WriteString("would be gone with nothing to restore it from.\n\n")
	b.WriteString("To add a field without losing your work, add it in place. It reaches the\n")
	b.WriteString("model, the API, the types, the schemas and the admin:\n\n")
	fmt.Fprintf(&b, "  grit generate field %s <name:type>\n", names.Pascal)
	b.WriteString("  grit migrate\n\n")
	b.WriteString("That covers scalar, select and toggle fields. For a relationship or a file,\n")
	b.WriteString("commit first and pass --force to regenerate.")
	return fmt.Errorf("%s", b.String())
}
