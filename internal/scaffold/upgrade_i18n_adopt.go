package scaffold

import (
	"os"
	"path/filepath"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// i18nWiredFiles are the files grit add i18n edits in a Next.js app.
var i18nWiredFiles = []string{
	"package.json",
	"next.config.ts",
	"app/layout.tsx",
	"components/chrome/PageHeader.tsx",
	"components/navbar.tsx",
}

// adoptI18nWired gives back to Grit the files whose only difference from what
// Grit last recorded is the i18n wiring Grit itself added.
//
// Before v3.222.0 grit add i18n (and grit new --i18n) edited these files
// without recording the result, so every upgrade of an i18n project reported
// them as edited by you and never updated them again. A file is adopted only
// when it is byte-identical to the template it is about to be replaced with
// plus that same wiring, worked out in a scratch copy of the app. That is proof
// there is no edit of anyone's in it to lose; anything else stays yours.
//
// Called before an app's files are written, with the templates about to be
// written, while upgrade's manifest recording is active.
func adoptI18nWired(appRoot string, templates map[string]string) int {
	if !fileContains(filepath.Join(appRoot, "package.json"), `"next-intl"`) {
		return 0
	}
	scratch, err := os.MkdirTemp("", "grit-i18n-")
	if err != nil {
		return 0
	}
	defer os.RemoveAll(scratch)
	mirror := filepath.Join(scratch, filepath.Base(appRoot))

	var candidates []string
	for _, rel := range i18nWiredFiles {
		path := filepath.Join(appRoot, filepath.FromSlash(rel))
		body, ok := templates[path]
		if !ok || !fileExists(path) || manifest.IsUnchanged(path) {
			continue
		}
		if writeScratch(filepath.Join(mirror, filepath.FromSlash(rel)), body) != nil {
			return 0
		}
		candidates = append(candidates, rel)
	}
	if len(candidates) == 0 {
		return 0
	}
	// wireI18n decides whether to feed the admin's catalogue by whether
	// lib/i18n.tsx exists, so the scratch copy has to agree with the real app.
	if fileExists(filepath.Join(appRoot, "lib", "i18n.tsx")) {
		if writeScratch(filepath.Join(mirror, "lib", "i18n.tsx"), adminI18nLib()) != nil {
			return 0
		}
	}
	// No routes.go under the scratch API root, so only the app is wired.
	if err := wireI18n(&i18nLayout{APIRoot: scratch, NextRoots: []string{mirror}}, &I18nResult{}); err != nil {
		return 0
	}

	adopted := 0
	for _, rel := range candidates {
		path := filepath.Join(appRoot, filepath.FromSlash(rel))
		current, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		expected, err := os.ReadFile(filepath.Join(mirror, filepath.FromSlash(rel)))
		if err != nil || string(current) != string(expected) {
			continue
		}
		manifest.Note(path, string(current))
		guardAdopt(path, string(current))
		adopted++
	}
	return adopted
}

func writeScratch(path, body string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}
