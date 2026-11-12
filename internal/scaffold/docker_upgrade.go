package scaffold

import (
	"fmt"
	"path/filepath"
)

// Bringing an existing project's Docker files up to date.
//
// writeDockerFiles was called from Run and RunSingle and from nowhere else, so
// every fix written into a Dockerfile template reached new projects only. An
// existing project kept whatever it was scaffolded with, for as many releases
// as it took somebody to try a production deploy and read a forty-line type
// error (grit#216). That is the same shape the desktop app had before its file
// map was added to the upgrade.
//
// These go through writeUpgradeFiles rather than writeFile, so a Dockerfile a
// project has edited is held back and named rather than overwritten. A
// deployment file is one people genuinely do edit: an extra RUN for a system
// library, a different base image, a proxy's certificates.

// upgradeDockerFiles rewrites the framework-owned Docker files that an
// existing project has.
//
// Only files that are already there. A project with no apps/web does not grow
// one, and a project that deleted its docker-compose.prod.yml on purpose does
// not get it back: upgrading is not the same as scaffolding, and a file
// somebody removed is a decision.
func upgradeDockerFiles(root string, opts Options, force bool) (int, error) {
	candidates := map[string]string{}

	add := func(path, body string) {
		if fileExists(path) {
			candidates[path] = body
		}
	}

	add(filepath.Join(root, "docker-compose.yml"), dockerCompose(opts))

	if opts.Architecture == ArchSingle {
		add(filepath.Join(root, "Dockerfile"), dockerfileSingle())
	} else {
		add(filepath.Join(root, "docker-compose.prod.yml"), dockerComposeProd(opts))
		add(filepath.Join(root, "docker-compose.shared-network.yml"), dockerComposeSharedNetwork(opts))
		add(filepath.Join(root, "apps", "api", "Dockerfile"), dockerfileAPI())
		add(filepath.Join(root, "apps", "api", "railway.json"), railwayJSON())
		add(filepath.Join(root, "apps", "api", ".dockerignore"), apiDockerIgnore)

		// Which image a frontend wants is decided by what the project is, not
		// by what the options say: a project scaffolded with Vite still has a
		// Vite app after an upgrade, and handing it the Next.js Dockerfile
		// would copy .next/standalone from a build that produces dist/.
		webRoot := webAppRoot(root, opts)
		if fileExists(filepath.Join(webRoot, "Dockerfile")) {
			if dirExists(filepath.Join(webRoot, "src", "routes")) || fileExists(filepath.Join(webRoot, "vite.config.ts")) {
				add(filepath.Join(webRoot, "Dockerfile"), dockerfileVite("web"))
				add(filepath.Join(webRoot, "nginx.conf"), viteNginxConf("web"))
			} else {
				add(filepath.Join(webRoot, "Dockerfile"), dockerfileNextJS("web"))
			}
		}
		adminRoot := filepath.Join(root, "apps", "admin")
		if fileExists(filepath.Join(adminRoot, "Dockerfile")) {
			if fileExists(filepath.Join(adminRoot, "vite.config.ts")) {
				add(filepath.Join(adminRoot, "Dockerfile"), dockerfileVite("admin"))
				add(filepath.Join(adminRoot, "nginx.conf"), viteNginxConf("admin"))
			} else {
				add(filepath.Join(adminRoot, "Dockerfile"), dockerfileNextJS("admin"))
			}
		}
		add(filepath.Join(root, "apps", "docs", "Dockerfile"), dockerfileNextJS("docs"))
	}

	if opts.Architecture != ArchAPI {
		add(filepath.Join(root, ".dockerignore"), dockerIgnore())
	}

	n, err := writeUpgradeFiles(candidates, force)
	if err != nil {
		return 0, fmt.Errorf("updating the Docker files: %w", err)
	}
	return n, nil
}
