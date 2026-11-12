package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Keeping the desktop screen's columns current when a model gains a field.
//
// grit#87: the mobile and desktop screens are written once, at generation, and
// no command updates them. `grit sync` reports the drift and cannot fix it, so
// a field added to a model is in the API, in the shared types, in the admin,
// and absent from the screen somebody is actually looking at on a laptop.
//
// The admin solved this at v3.31.16 with marker fences and an insert-only
// pass: a new field is appended inside the fence, and a row somebody has
// customised is left exactly as it is because the field is already mentioned.
// This is that, for apps/desktop.
//
// # What it deliberately does not do
//
// It does not delete. A field removed from the Go model keeps its column here,
// along with whatever was done to it. A sync that removed things would be a
// sync nobody could leave running.
//
// It does not touch the Expo list. That screen renders a header block and a
// row block that must line up, with a TABLE_WIDTH constant summed from the
// column widths: three things that have to agree, where an insert that gets
// one of them wrong produces a table whose headers sit over the wrong columns.
// `grit generate resource <Name> --force` rebuilds it from the model instead,
// which is why --force stopped needing --fields in the same release.

const (
	desktopColsStart   = "// grit:cols:auto-start"
	desktopColsEnd     = "// grit:cols:auto-end"
	desktopSearchStart = "// grit:search-keys:start"
	desktopSearchEnd   = "// grit:search-keys:end"
)

var desktopSearchKeysRe = regexp.MustCompile(`const SEARCH_KEYS = \[([^\]]*)\];`)

// SyncDesktopScreen appends model fields the desktop screen has never seen to
// its COLUMNS array, and adds the searchable ones to SEARCH_KEYS.
//
// Returns the number of columns added and any warnings for the CLI.
func SyncDesktopScreen(root string, s GoStruct) (added int, warnings []string, err error) {
	// Only a resource the generator wrote. A hand-written screen was never
	// derived from the model, so its "missing" fields were never meant to be
	// there and inserting them would be vandalism.
	if !generatedResource(root, s.Name) {
		return 0, nil, nil
	}

	pluralKebab := strings.ReplaceAll(Pluralize(toSnakeCase(s.Name)), "_", "-")
	path := filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "app", pluralKebab+".index.tsx")
	raw, readErr := os.ReadFile(path)
	if os.IsNotExist(readErr) {
		return 0, nil, nil // this project has no desktop app, or no screen for this resource
	}
	if readErr != nil {
		return 0, nil, fmt.Errorf("reading %s: %w", path, readErr)
	}
	content := string(raw)

	// Markers arrived with this feature. A screen generated before it does not
	// have them, and the honest answer is to say so rather than guess at where
	// the array ends.
	if !strings.Contains(content, desktopColsEnd) {
		return 0, []string{fmt.Sprintf(
			"  ⚠ %s: no grit:cols:auto-end marker, so new columns cannot be added in place.\n"+
				"     Rebuild the screen with `grit generate resource %s --force`, which now reads the field list from the model.",
			filepath.ToSlash(filepath.Join("apps", "desktop", "frontend", "src", "routes", "app", pluralKebab+".index.tsx")),
			s.Name)}, nil
	}

	updated := content
	var newKeys []string

	for _, f := range s.Fields {
		key := f.JSONName
		if key == "" {
			key = toSnakeCase(f.Name)
		}
		if isAutoField(key) || frameworkColumn(key) || key == "-" || f.Name == "Version" {
			continue
		}
		// An embedding: 1536 numbers in a table cell. The generator leaves it
		// out and so does this, or sync would add back what generate omitted.
		if f.GoType == "vectors.Vector" {
			continue
		}
		if isAssociationField(f) {
			continue
		}
		// Mentioned anywhere in the file, inside the fence or in a row
		// somebody moved out of it. Either way it is known about.
		if strings.Contains(updated, `"`+key+`"`) {
			continue
		}

		line := "  { key: \"" + key + "\", label: \"" + humanizeLabel(f.Name) + "\", format: \"" + desktopColumnFormat(f) + "\" },"
		updated = strings.Replace(updated, "  "+desktopColsEnd, line+"\n  "+desktopColsEnd, 1)
		added++
		if desktopFieldIsSearchable(f) {
			newKeys = append(newKeys, key)
		}
	}

	if added == 0 {
		return 0, nil, nil
	}

	// A column you can see and cannot find is half a feature, so the search
	// keys move with the columns.
	if len(newKeys) > 0 {
		if strings.Contains(updated, desktopSearchEnd) {
			updated = addDesktopSearchKeys(updated, newKeys)
		} else {
			warnings = append(warnings, fmt.Sprintf(
				"  ⚠ %s: added %d column(s) but SEARCH_KEYS has no marker, so they are not searchable yet",
				pluralKebab, added))
		}
	}

	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return 0, warnings, fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	return added, warnings, nil
}

// addDesktopSearchKeys appends keys to the SEARCH_KEYS literal, keeping what
// is already there.
func addDesktopSearchKeys(src string, keys []string) string {
	m := desktopSearchKeysRe.FindStringSubmatch(src)
	if m == nil {
		return src
	}
	existing := strings.TrimSpace(m[1])
	quoted := make([]string, 0, len(keys))
	for _, k := range keys {
		if strings.Contains(existing, `"`+k+`"`) {
			continue
		}
		quoted = append(quoted, `"`+k+`"`)
	}
	if len(quoted) == 0 {
		return src
	}
	joined := strings.Join(quoted, ", ")
	if existing != "" {
		joined = existing + ", " + joined
	}
	return strings.Replace(src, m[0], "const SEARCH_KEYS = ["+joined+"];", 1)
}

// desktopColumnFormat is how a column renders in the desktop table. It mirrors
// the generator's own choice so a synced column looks like a generated one.
func desktopColumnFormat(f GoField) string {
	switch {
	case f.GoType == "bool":
		return "boolean"
	case f.GoType == "money.Money":
		return "money"
	case f.GoType == "time.Time" || f.GoType == "*time.Time":
		return "date"
	case f.GoType == "files.FileRef" || f.GoType == "*files.FileRef":
		return "image"
	case f.Enum != "":
		return "badge"
	case strings.HasPrefix(f.GoType, "int") || strings.HasPrefix(f.GoType, "float") ||
		strings.HasPrefix(f.GoType, "uint"):
		return "number"
	default:
		return "text"
	}
}

// desktopFieldIsSearchable reports a column worth matching a typed query
// against. Searching a boolean or a timestamp by substring finds nothing and
// makes the search look broken.
func desktopFieldIsSearchable(f GoField) bool {
	if f.GoType != "string" {
		return false
	}
	return f.Enum == "" && f.Format != "tel"
}
