package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// pruneMovedUserRoleSyncTest removes handlers/user_role_sync_test.go from a
// project that is getting the user service.
//
// Its two cases moved to services/user_test.go, beside the code they are about.
// The file has to go rather than stay, because what it called,
// syncUserRoleAssignment, is no longer in the handler the upgrade just
// delivered: left behind, it is a package that does not compile under go test,
// which is worse than a missing test and would be blamed on the upgrade rather
// than on the file.
//
// Only when the file is still as Grit wrote it. Somebody who added their own
// cases to it keeps the file and gets told what to change, because deleting
// somebody's tests to tidy up a rename is not a trade an upgrade gets to make.
func pruneMovedUserRoleSyncTest(root string, opts Options) error {
	handlers := filepath.Join(opts.APIRoot(root), "internal", "handlers")
	path := filepath.Join(handlers, "user_role_sync_test.go")
	if !fileExists(path) {
		return nil
	}
	// A handler somebody edited is kept rather than replaced, and it still
	// declares what this test calls. Removing the test there would leave that
	// function with no caller, which their own linter reports: so the test stays,
	// and it still passes.
	if fileContains(filepath.Join(handlers, "user.go"), "func syncUserRoleAssignment(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	key, inside := manifest.Rel(root, path)
	if !inside {
		return nil
	}
	if m.StatusOf(root, key) != manifest.Unchanged {
		fmt.Println("  ⚠ apps/api/internal/handlers/user_role_sync_test.go: you have edited this, and " +
			"syncUserRoleAssignment has moved to services.UserService.SyncRoleAssignment: call that, " +
			"or delete the file (its cases are now in internal/services/user_test.go)")
		return nil
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("removing %s: %w", path, err)
	}
	fmt.Println("  ✓ apps/api/internal/handlers/user_role_sync_test.go: removed, its cases are now in internal/services/user_test.go")
	return nil
}
