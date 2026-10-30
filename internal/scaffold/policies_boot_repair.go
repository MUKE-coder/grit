package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairPoliciesBoot makes an existing project call policies.Register().
//
// It never did. The call had no site anywhere in a generated project: not in
// main, not in routes, nowhere. So every rule anybody wrote in
// internal/policies was registered into a map that nothing ever filled, and
// no rule ran on any path. A Patch on a row a rule denied answered 200.
//
// That is the worst shape an authorization feature can fail in, because
// everything about it looks finished: the package is generated with the
// resource, it is documented, it carries a worked example, and each resource
// injects its own registration into a marker in it. The only missing piece
// was the one nobody looks at twice.
//
// Found while testing the enterprise review's finding 2.4, which reported
// that policies were evadable through a bulk action. They were evadable
// through everything.
//
// Every project that has ever run `grit generate resource` is affected, so
// this cannot wait for the next generate.
func repairPoliciesBoot(root string, opts Options) error {
	api := opts.APIRoot(root)

	// No policies package means no resources generated yet, and the generator
	// wires the call when it creates one.
	if !fileExists(filepath.Join(api, "internal", "policies", "policies.go")) {
		return nil
	}
	main := filepath.Join(api, "cmd", "server", "main.go")
	if !fileExists(main) {
		return nil
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()
	return repairSourceFile(root, m, main, func(src string) (string, []string, []string) {
		return repairPoliciesBootSource(src, module)
	})
}

func repairPoliciesBootSource(src, module string) (string, []string, []string) {
	if strings.Contains(src, "policies.Register()") {
		return src, nil, nil
	}

	// Before the router is built. The rules are a package-level map and need
	// nothing from the database, so anywhere before the server serves is
	// correct; here is where a reader looks for wiring.
	const anchor = "\t// Setup router\n"
	if !strings.Contains(src, anchor) {
		return src, nil, []string{
			"could not find where to call policies.Register() in main.go, so " +
				"every rule in internal/policies stays dead: nothing registers " +
				"them and no rule runs on any path. Add policies.Register() " +
				"before routes.Setup by hand"}
	}

	out := strings.Replace(src, anchor,
		"\t// grit:main:boot\n\tpolicies.Register()\n\n"+anchor, 1)

	importPath := module + "/internal/policies"
	if !strings.Contains(out, importPath) {
		added, ok := AddImportGroup(out, importPath)
		if !ok {
			return src, nil, []string{
				"main.go has no import block to add " + importPath + " to, so " +
					"the policies stay unregistered: add the import and " +
					"policies.Register() by hand"}
		}
		out = added
	}

	return out, []string{"main calls policies.Register(), so the rules in " +
		"internal/policies actually run (they never have)"}, nil
}
