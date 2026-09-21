package scaffold

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The React every web frontend pins, and the one Expo pins.
//
// Expo SDK 54 needs React 19.1.0 exactly: React 19 checks that react and
// React Native's renderer are the same version. The web apps pin 19.2.7. In a
// monorepo that has both, the hoisted node_modules holds Expo's copy at the
// root and every other app gets a nested one, so a test loaded React more than
// once: Testing Library's act() flushed a different React from the one
// rendering, render() left the container empty, and components failed with
// "Invalid hook call". Every web and admin component test failed in a new
// --expo project. With Expo in the project, every app uses its version, and
// there is one React.
const (
	webReactVersion  = "19.2.7"
	expoReactVersion = "19.1.0"
)

// reactVersionFor is the React the package.json templates pin for this shape.
func reactVersionFor(opts Options) string {
	if opts.ShouldIncludeExpo() {
		return expoReactVersion
	}
	return webReactVersion
}

// monorepoReactVersion is the React a project's web frontends should pin, read
// from the project on disk.
func monorepoReactVersion(root string) string {
	if dirExists(filepath.Join(root, "apps", "expo")) {
		return expoReactVersion
	}
	return webReactVersion
}

// reactPackageJSONs are the frontends that pin React. Expo's own is left out:
// Expo decides its version.
func reactPackageJSONs(root string) []string {
	return []string{
		filepath.Join(root, "package.json"), // single-app scaffold
		filepath.Join(root, "apps", "admin", "package.json"),
		filepath.Join(root, "apps", "web", "package.json"),
		filepath.Join(root, "apps", "docs", "package.json"),
		filepath.Join(root, "apps", "desktop", "frontend", "package.json"),
	}
}

// alignReactVersions pins react and react-dom to one exact version in every
// web frontend, leaving the rest of each file untouched, and returns the files
// it changed. Older scaffolds used "^19.0.0", which drifts, or a "19.1.0" pin
// pnpm 10 applied to react only, and a react/react-dom mismatch white-screens
// the app.
func alignReactVersions(root string) []string {
	target := monorepoReactVersion(root)
	var pairs []string
	for _, from := range []string{`^19.0.0`, expoReactVersion, webReactVersion} {
		if from == target {
			continue
		}
		pairs = append(pairs,
			`"react": "`+from+`"`, `"react": "`+target+`"`,
			`"react-dom": "`+from+`"`, `"react-dom": "`+target+`"`,
		)
	}
	replacer := strings.NewReplacer(pairs...)
	var changed []string
	for _, path := range reactPackageJSONs(root) {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		out := replacer.Replace(string(data))
		if out == string(data) {
			continue
		}
		if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
			continue
		}
		// Grit's own edit: a new project's package.json is not "edited by you".
		manifest.Refresh(path)
		changed = append(changed, path)
	}
	return changed
}
