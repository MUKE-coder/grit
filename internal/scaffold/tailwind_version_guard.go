package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

// stylesheetIsV3 reports whether an app's stylesheet still uses the Tailwind v3
// directives.
//
// The v3 to v4 move spans three files that only work as a set: package.json
// names the engine, postcss.config.js names the plugin, and globals.css uses
// either the v3 @tailwind directives or the v4 @import. Upgrade writes all
// three, and the manifest guard holds back whichever ones the reader has
// edited, individually. Nobody edits all three or none of them.
//
// So the reachable state is a v4 plugin parsing a v3 stylesheet, and it does
// not degrade: the first utility built from tailwind.config.ts becomes
// "Cannot apply unknown utility" and the app stops building. The upgrade
// reports the held-back files, but nothing connects "your package.json was
// left alone" to "your web app no longer compiles".
//
// So the postcss config is left on v3 whenever the stylesheet is still on v3.
// A consistent v3 app keeps working, and the migration stays available as
// something to do deliberately rather than something an upgrade does to two
// files out of three.
func stylesheetIsV3(appRoot string) bool {
	for _, rel := range []string{
		filepath.Join("app", "globals.css"),
		filepath.Join("src", "globals.css"),
		filepath.Join("src", "styles", "globals.css"),
		filepath.Join("app", "styles", "globals.css"),
	} {
		body, err := os.ReadFile(filepath.Join(appRoot, rel))
		if err != nil {
			continue
		}
		css := string(body)
		if strings.Contains(css, `@import "tailwindcss"`) {
			return false
		}
		if strings.Contains(css, "@tailwind base") ||
			strings.Contains(css, "@tailwind utilities") {
			return true
		}
	}
	// No stylesheet found, or one that says neither: a fresh scaffold, where
	// the file is about to be written as v4 anyway.
	return false
}

// postCSSConfigFor returns the PostCSS config an app should have, which is the
// v3 one only while its stylesheet is still v3.
func postCSSConfigFor(appRoot string) string {
	if stylesheetIsV3(appRoot) {
		return legacyPostCSSConfigV3()
	}
	return webPostCSSConfig()
}

// legacyPostCSSConfigV3 is the pre-v3.183.0 config, kept for apps whose
// stylesheet has not been migrated.
//
// Not dead code and not a fallback nobody reaches: every project scaffolded
// before v3.183.0 that customised its globals.css lands here on upgrade, and
// stays here until somebody migrates the stylesheet on purpose.
func legacyPostCSSConfigV3() string {
	return `// Tailwind v3.
//
// Your globals.css still uses the v3 directives (@tailwind base/components/
// utilities), so this stays on the v3 plugin to match. Swapping this file
// alone to @tailwindcss/postcss makes every utility built from
// tailwind.config.ts an unknown utility, and the build stops on the first one.
//
// To move to v4, change all three together:
//
//   1. package.json:  tailwindcss ^4, add @tailwindcss/postcss, drop autoprefixer
//   2. globals.css:   replace the three @tailwind lines with
//                       @import "tailwindcss";
//                       @config "../tailwind.config.ts";
//                     The @config line keeps your existing theme file working,
//                     so nothing has to move into an @theme block on day one.
//   3. this file:     "@tailwindcss/postcss": {}
//
// Then delete this comment.
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`
}
