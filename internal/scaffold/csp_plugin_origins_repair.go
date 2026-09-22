package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// A plugin that loads a third party's script or frame had nowhere to say so:
// every CSP directive was one fixed string, and Stripe.js, which needs its own
// origin in script-src, frame-src and connect-src, was refused with nothing
// but a console violation. pluginOrigins is the list a plugin appends to, at
// the grit:csp-origins marker, and removal takes the lines out again. Found
// building the storefront blueprint.
const (
	cspPluginOrigins = `// Third-party origins a plugin needs, as [directive, origin]. Each is
// appended to its directive below. grit plugin add stripe puts Stripe.js here,
// and grit plugin remove takes it out again.
const pluginOrigins: [string, string][] = [
  // grit:csp-origins
];

`
	// frame-src 'self' is what default-src already allowed; it is spelled out
	// so a plugin has a directive to add a frame's origin to.
	cspFrameSrc = `"frame-src 'self'",`

	cspJoinOld         = "].join(\"; \");\n"
	cspJoinWithPlugins = `]
  .map((directive) => {
    const name = directive.split(" ")[0];
    const extra = pluginOrigins.filter(([d]) => d === name).map(([, origin]) => origin);
    return extra.length > 0 ? directive + " " + extra.join(" ") : directive;
  })
  .join("; ");
`
	cspFrameAncestors = "  \"frame-ancestors 'none'\",\n"
)

// repairCSPPluginOrigins gives an existing Next.js app the pluginOrigins list.
func repairCSPPluginOrigins(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, app := range []string{"web", "admin"} {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js"} {
			if path := filepath.Join(root, "apps", app, name); fileExists(path) {
				if err := repairTextFile(root, m, path, repairCSPPluginOriginsSource); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func repairCSPPluginOriginsSource(src string) (string, []string, []string) {
	if strings.Contains(src, "grit:csp-origins") || !strings.Contains(src, "const csp = [") {
		return src, nil, nil
	}
	if strings.Count(src, "\nconst csp = [\n") != 1 || strings.Count(src, cspJoinOld) != 1 ||
		strings.Count(src, cspFrameAncestors) != 1 || strings.Contains(src, "frame-src") {
		return src, nil, []string{"the Content-Security-Policy is not the one Grit wrote: a plugin that loads a third-party script (grit plugin add stripe) cannot add its origins, so add them to script-src, frame-src and connect-src by hand"}
	}
	out := strings.Replace(src, "\nconst csp = [\n", "\n"+cspPluginOrigins+"const csp = [\n", 1)
	out = strings.Replace(out, cspFrameAncestors, "  "+cspFrameSrc+"\n"+cspFrameAncestors, 1)
	out = strings.Replace(out, cspJoinOld, cspJoinWithPlugins, 1)
	return out, []string{"the Content-Security-Policy has a list plugins add third-party origins to"}, nil
}
