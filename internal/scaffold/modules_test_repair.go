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
	if !strings.Contains(src, "func TestModules_DefaultOn(") {
		return src, nil, nil
	}
	out := src
	var fixed, warnings []string

	// APP_ENV, from v3.243.0: Load refuses production without real secrets.
	if !strings.Contains(out, `"APP_ENV"`) {
		if strings.Count(out, modulesTestSecretLine) != 1 {
			warnings = append(warnings, "config/modules_test.go is not the one Grit wrote: set APP_ENV=development before the tests call Load, which refuses production without real secrets")
		} else {
			out = strings.Replace(out, modulesTestSecretLine, modulesTestSecretLine+modulesTestEnvLines, 1)
			fixed = append(fixed, "the module flag tests run in development, where Load does not demand production secrets")
		}
	}

	// SSO: the one module that ships off, so a test asserting that every
	// module defaults on fails on every project the moment the flag arrives.
	if !strings.Contains(out, "defaultOff") {
		before := out
		ok := true
		for _, e := range modulesTestSSOEdits {
			if strings.Count(out, e.old) != 1 {
				ok = false
				break
			}
			out = strings.Replace(out, e.old, e.new, 1)
		}
		if !ok {
			out = before
			warnings = append(warnings, modulesTestSSOWarn)
		} else {
			fixed = append(fixed, "the module flag tests know that SSO ships disabled, which is what MODULE_SSO decides")
		}
	}

	if out == src {
		return src, nil, warnings
	}
	return out, fixed, warnings
}

const modulesTestSSOWarn = "config/modules_test.go is not the one Grit wrote, so it still " +
	"says every module defaults on: SSO does not (MODULE_SSO), and the test will fail until " +
	"it names the exception"

// modulesTestSSOEdits turns the all-on test into one that names its exception.
// Each old is the template's own text: see configModulesTestGo.
var modulesTestSSOEdits = []struct{ old, new string }{
	{
		old: "// Default ON: an existing app upgrading must behave exactly as before.\n",
		new: "// defaultOff names the modules that ship disabled.\n" +
			"//\n" +
			"// SSO is the only one, and it is here rather than absent so that adding\n" +
			"// another is a deliberate edit to this list. Its cost is four public\n" +
			"// endpoints, one of which parses a SAML assertion through a package with a\n" +
			"// called advisory and no fixed version, and most projects never configure an\n" +
			"// identity provider.\n" +
			"var defaultOff = map[string]bool{\"sso\": true}\n" +
			"\n" +
			"// Default ON: an existing app upgrading must behave exactly as before.\n",
	},
	{
		old: "\tfor _, k := range []string{\"MODULE_AI\", \"MODULE_JOBS\", \"MODULE_BACKUP\"} {\n",
		new: "\tfor _, k := range []string{\"MODULE_AI\", \"MODULE_JOBS\", \"MODULE_BACKUP\", \"MODULE_SSO\"} {\n",
	},
	{
		old: "\tfor name, on := range cfg.Modules.Map() {\n" +
			"\t\tif !on {\n" +
			"\t\t\tt.Errorf(\"module %q defaulted to OFF; upgrading an app must not disable features\", name)\n" +
			"\t\t}\n" +
			"\t}\n",
		new: "\tfor name, on := range cfg.Modules.Map() {\n" +
			"\t\tif defaultOff[name] {\n" +
			"\t\t\tif on {\n" +
			"\t\t\t\tt.Errorf(\"module %q defaulted to ON and is meant to ship disabled\", name)\n" +
			"\t\t\t}\n" +
			"\t\t\tcontinue\n" +
			"\t\t}\n" +
			"\t\tif !on {\n" +
			"\t\t\tt.Errorf(\"module %q defaulted to OFF; upgrading an app must not disable features\", name)\n" +
			"\t\t}\n" +
			"\t}\n",
	},
	{
		old: "\tif len(m) != 11 {\n",
		new: "\tif len(m) != 12 {\n",
	},
	{
		old: "\tfor name := range m {\n" +
			"\t\tif !cfg.Modules.Enabled(name) {\n" +
			"\t\t\tt.Errorf(\"Map key %q is not handled by Enabled()\", name)\n" +
			"\t\t}\n" +
			"\t}\n",
		new: "\t// Enabled must agree with Map, rather than be true: a module that ships\n" +
			"\t// disabled is false in both, and the thing that would break the admin is\n" +
			"\t// the two disagreeing.\n" +
			"\tfor name, on := range m {\n" +
			"\t\tif cfg.Modules.Enabled(name) != on {\n" +
			"\t\t\tt.Errorf(\"Map says %q is %v and Enabled() says %v; the admin hides nav \"+\n" +
			"\t\t\t\t\"from Map and the API gates on Enabled, so the two disagreeing shows \"+\n" +
			"\t\t\t\t\"a screen whose routes are not mounted\", name, on, cfg.Modules.Enabled(name))\n" +
			"\t\t}\n" +
			"\t}\n",
	},
}
