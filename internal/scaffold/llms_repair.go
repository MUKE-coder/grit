package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// /llms.txt and /llms-full.txt, mounted in a project that predates them.
//
// internal/llms arrives with the framework-owned files. What does not arrive is
// the wiring, because routes.go is the developer's file: it carries their
// handlers, their middleware and their routes, and an upgrade rewrites parts of
// it rather than replacing it.
//
// Two edits, both anchored on text Grit wrote: the mount goes after the API
// reference's own gate, because the two files describe the same surface and
// belong behind the same condition, and the handler goes beside eventBusStatus,
// where the other function routes.go owns already lives.

// repairLLMSFiles mounts the two files in an existing project.
func repairLLMSFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	// Without the package this would write a mount that does not compile, so it
	// does nothing instead. writeFrameworkOwnedFiles runs earlier in the upgrade
	// and delivers it.
	if !fileExists(routes) || !fileContains(filepath.Join(apiRoot, "internal", "llms", "llms.go"), "func Index(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, routes, repairLLMSFilesSource(opts.Module()))
}

func repairLLMSFilesSource(module string) func(string) (string, []string, []string) {
	return func(src string) (string, []string, []string) {
		if strings.Contains(src, "registerLLMSFiles(r, cfg)") {
			return src, nil, nil
		}
		// The gate the reference is behind. A project that has turned its
		// reference off, or never had one, is not a project to mount these in:
		// they say the same things about the same routes.
		if strings.Count(src, routesDocsBlock) != 1 {
			return src, nil, nil
		}
		// Written beside eventBusStatus, so a routes.go that has not taken the
		// four-state repair yet waits one more upgrade rather than having a
		// function dropped at a line number.
		anchor := "func eventBusStatus() health.Report {"
		if !strings.Contains(src, anchor) {
			return src, nil, nil
		}
		end := strings.Index(src[strings.Index(src, anchor):], "\n}\n")
		if end < 0 {
			return src, nil, []string{"could not find the end of eventBusStatus in routes.go, so /llms.txt was not mounted"}
		}
		at := strings.Index(src, anchor) + end + len("\n}\n")

		out := src[:at] + llmsRegisterFunc + src[at:]
		out = strings.Replace(out, routesDocsBlock, routesDocsBlock+llmsRoutesBlock, 1)
		var ok bool
		for _, path := range []string{module + "/internal/llms", module + "/internal/paginate"} {
			if out, ok = addImportGroup(out, path); !ok {
				return src, nil, []string{"could not add the imports /llms.txt needs to routes.go, so it was not mounted"}
			}
		}
		return out, []string{"/llms.txt and /llms-full.txt are served beside the API reference, for an agent pointed at this service"}, nil
	}
}
