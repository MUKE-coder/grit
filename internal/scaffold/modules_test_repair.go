package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairModulesTest brings internal/config/modules_test.go up to the production
// default of v3.243.0. config.Load refuses to start in production without real
// secrets, so the module flag tests, which call Load, have failed in every
// project since, and CI reported tail's exit code instead of the test's.
//
// The file is written when a project is scaffolded, and upgrade does not deliver
// it, so the one line anchors on what Grit wrote.
func repairModulesTest(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "config", "modules_test.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairModulesTestSource)
}

const (
	modulesTestSecretLine = "\tos.Setenv(\"JWT_SECRET\", \"test-secret-key-for-module-flag-tests-only\")\n"
	modulesTestEnvLines   = "\t// Load refuses to start in production without real secrets, and production\n" +
		"\t// is the default.\n" +
		"\tos.Setenv(\"APP_ENV\", \"development\")\n"
)

func repairModulesTestSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func TestModules_DefaultOn(") || strings.Contains(src, `"APP_ENV"`) {
		return src, nil, nil
	}
	if strings.Count(src, modulesTestSecretLine) != 1 {
		return src, nil, []string{"config/modules_test.go is not the one Grit wrote: set APP_ENV=development before the tests call Load, which refuses production without real secrets"}
	}
	out := strings.Replace(src, modulesTestSecretLine, modulesTestSecretLine+modulesTestEnvLines, 1)
	return out, []string{"the module flag tests run in development, where Load does not demand production secrets"}, nil
}
