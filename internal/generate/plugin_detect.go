package generate

import "github.com/MUKE-coder/grit/v3/internal/project"

// PluginInstalled reports whether a plugin is recorded in the project's
// lockfile.
//
// Kept as a name in this package because callers here read naturally without
// the prefix, but the answer comes from internal/project, which internal/
// scaffold reads too: the generator decides what to emit for a new resource
// and the upgrade decides what to repair in an existing one, and those two
// must agree about whether a project is multitenant.
func PluginInstalled(root, name string) bool {
	return project.PluginInstalled(root, name)
}
