package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairRoleModel delivers internal/models/role.go on an upgrade.
//
// Nothing did. writeSecurityFiles writes it, and only Run and RunSingle call
// that, so the Role model has never changed in a project that already exists.
// It went unnoticed while nothing else depended on the file: when
// internal/authz/ceiling.go started calling models.PlatformRoles, which
// repairGrantCeiling does deliver, an upgraded project got the caller without
// the function and stopped compiling.
//
// So this runs BEFORE repairGrantCeiling. The two files have to travel
// together now, and the one that declares goes first.
//
// Held back the usual way when the project has edited it, and said out loud,
// because a project that keeps its own role.go and takes the new ceiling.go
// is the state that does not build.
func repairRoleModel(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	path := filepath.Join(apiRoot, "internal", "models", "role.go")
	if !fileExists(path) {
		return nil
	}

	content := strings.ReplaceAll(roleModelGo(), "{{MODULE}}", opts.Module())
	if err := writeFile(path, content); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	// writeFile goes through the upgrade guard, so an edited file is kept.
	// Whether that happened is worth knowing here specifically: the rest of
	// the upgrade assumes PlatformRoles exists.
	if !fileContains(path, "func PlatformRoles(") {
		m, err := manifest.Load(root)
		if err == nil {
			if rel, ok := manifest.Rel(root, path); ok && m.StatusOf(root, rel) != manifest.Untracked {
				fmt.Println("  ⚠ internal/models/role.go is yours and does not declare " +
					"PlatformRoles. internal/authz and internal/services now call it, " +
					"so add it from the template or the API will not compile (grit#130).")
			}
		}
	}
	return nil
}
