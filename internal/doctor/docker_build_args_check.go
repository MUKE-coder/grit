package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Two things a Docker build gets wrong that nothing says out loud until a
// first production deploy, which is the worst moment to find either (grit#216).
//
// Both are offline file reads, like every other check here: a doctor that
// needs a registry or a daemon is one CI cannot run.

var doctorNextPublicRe = regexp.MustCompile(`NEXT_PUBLIC_[A-Z0-9_]+`)

// checkNextPublicBuildArgs finds a NEXT_PUBLIC_* variable the app reads and
// its Dockerfile cannot accept.
//
// Next inlines these at build time, so one that is not a build arg takes its
// development default in the deployed image. Setting it in the platform's
// environment panel does not help: by then the value is already in the bundle.
//
// NEXT_PUBLIC_STORAGE_URL is the one that costs most. apps/web/next.config.ts
// builds the Content-Security-Policy from it and falls back to a MinIO on
// localhost, so the deployed site permits only the developer's laptop as a
// storage origin. Uploads are presigned PUTs the browser makes straight to
// object storage: they are blocked, and so is every stored image. There is no
// HTTP status and no server log line, only a console violation.
func checkNextPublicBuildArgs(p *project) []Finding {
	var findings []Finding

	for _, app := range []string{"web", "admin"} {
		appDir := filepath.Join(p.root, "apps", app)
		dockerfile := filepath.Join(appDir, "Dockerfile")
		raw, err := os.ReadFile(dockerfile)
		if err != nil {
			continue // no such app, or no Dockerfile: nothing to be wrong about
		}
		accepted := map[string]bool{}
		for _, line := range strings.Split(string(raw), "\n") {
			if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "ARG "); ok {
				accepted[strings.TrimSpace(strings.SplitN(rest, "=", 2)[0])] = true
			}
		}

		missing := map[string]bool{}
		_ = filepath.Walk(appDir, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			if info.IsDir() {
				switch info.Name() {
				case "node_modules", ".next", "dist":
					return filepath.SkipDir
				}
				return nil
			}
			switch filepath.Ext(path) {
			case ".ts", ".tsx", ".js", ".mjs", ".cjs":
			default:
				return nil
			}
			body, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil
			}
			for _, name := range doctorNextPublicRe.FindAllString(string(body), -1) {
				if !accepted[name] {
					missing[name] = true
				}
			}
			return nil
		})

		if len(missing) == 0 {
			continue
		}
		names := make([]string, 0, len(missing))
		for name := range missing {
			names = append(names, name)
		}
		sort.Strings(names)

		level := "warning"
		why := "so each takes its development default in the built image"
		if missing["NEXT_PUBLIC_STORAGE_URL"] {
			// Not a preference: the deployed app cannot upload.
			level = "error"
			why = "NEXT_PUBLIC_STORAGE_URL among them. That leaves a Content-Security-Policy " +
				"permitting only localhost, so every upload and every stored image is blocked, " +
				"with no HTTP status and nothing in any log"
		}

		findings = append(findings, Finding{
			Level:    level,
			Resource: "apps/" + app + "/Dockerfile",
			Message: fmt.Sprintf("the app reads %s and this Dockerfile declares no ARG for them: %s",
				strings.Join(names, ", "), why),
			Fix: "Add an ARG and a matching ENV for each, and pass them in docker-compose.prod.yml's " +
				"build.args. Next inlines NEXT_PUBLIC_* at build time, so a runtime environment " +
				"variable is too late and changing one needs a rebuild rather than a restart. " +
				"`grit upgrade` writes both files.",
		})
	}
	return findings
}

// checkDockerCopiesWorkspacePackages finds a builder stage that will not build
// a workspace whose packages have their own node_modules.
//
// pnpm's hoisted linker puts one copy of each dependency at the root, except
// where two packages need different versions, and then it nests one. A
// packages/* member whose dependency pins a peer more tightly than the root
// does therefore has its own node_modules, and a builder stage that does not
// copy it resolves the root's older copy instead.
//
// The failure is forty lines of structural type mismatch naming two paths that
// are the same library at two versions, in a file nobody changed. It reads as
// a bug in your own code.
func checkDockerCopiesWorkspacePackages(p *project) []Finding {
	// Only a project that actually has one. A workspace where every package
	// is flat is not affected and should not be told about it.
	var nested []string
	entries, err := os.ReadDir(filepath.Join(p.root, "packages"))
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if info, statErr := os.Stat(filepath.Join(p.root, "packages", e.Name(), "node_modules")); statErr == nil && info.IsDir() {
			nested = append(nested, "packages/"+e.Name())
		}
	}
	if len(nested) == 0 {
		return nil
	}
	sort.Strings(nested)

	var findings []Finding
	for _, app := range []string{"web", "admin"} {
		path := filepath.Join(p.root, "apps", app, "Dockerfile")
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		if strings.Contains(string(raw), "COPY --from=deps /app/packages ./packages") {
			continue
		}
		findings = append(findings, Finding{
			Level:    "error",
			Resource: "apps/" + app + "/Dockerfile",
			Message: fmt.Sprintf(
				"%s %s its own node_modules and the builder stage does not copy it, so the image will not build",
				strings.Join(nested, ", "), plural(len(nested), "has", "have")),
			Fix: "Add `COPY --from=deps /app/packages ./packages` after `COPY . .` in the builder " +
				"stage. It must come after, or the context copy puts a source tree with no " +
				"node_modules back over it. `grit upgrade` writes it.",
		})
	}
	return findings
}

// plural picks the verb for a count.
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
