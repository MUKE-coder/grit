package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// isUserOwnedEmbeddedAdminFile is isUserOwnedAdminFile for a panel written into
// another app: its code under admin-panel/ and its routes under app/admin/ (a
// Next.js web app) or routes/admin/ (an SPA). host is the directory both hang off.
func isUserOwnedEmbeddedAdminFile(host, path string) bool {
	rel, err := filepath.Rel(host, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	switch {
	case strings.HasPrefix(rel, "admin-panel/"):
		rel = strings.TrimPrefix(rel, "admin-panel/")
	case strings.HasPrefix(rel, "app/admin/"):
		rel = "app/" + strings.TrimPrefix(rel, "app/admin/")
	case strings.HasPrefix(rel, "routes/admin/"):
		rel = "app/" + strings.TrimPrefix(rel, "routes/admin/")
	default:
		return false
	}
	return isUserOwnedAdminRel(rel)
}

var (
	registryImport = regexp.MustCompile(`(?m)^import \{\s*(\w+)\s*\} from ['"]\./`)
	resourceExport = regexp.MustCompile(`(?m)^export const (\w+Resource)\b`)
)

// templateRegistryImports are the definitions resources/index.ts imports as the
// scaffold writes it.
var templateRegistryImports = map[string]bool{"usersResource": true, "blogsResource": true}

// restoreResourceRegistrations puts back the registry entries an upgrade dropped.
//
// Before v3.272.0, upgrading a project whose admin panel lives inside its web
// app or SPA wrote the template over admin-panel/resources/index.ts, so every
// resource grit generate had registered disappeared from the panel while its
// definition and pages stayed on disk. A registry that imports nothing beyond
// the template's own resources, with other definitions beside it, is that state,
// and each of those definitions is registered again. A registry that imports
// anything else was edited by a person and is left as it is.
func restoreResourceRegistrations(root string, opts Options) ([]string, error) {
	var restored []string
	seen := map[string]bool{}
	for _, index := range []string{
		filepath.Join(root, "apps", "admin", "resources", "index.ts"),
		filepath.Join(root, "apps", "admin", "src", "resources", "index.ts"),
		filepath.Join(root, "apps", "web", "admin-panel", "resources", "index.ts"),
		filepath.Join(spaHostRoot(root, opts), "src", "admin-panel", "resources", "index.ts"),
	} {
		if seen[index] || !fileExists(index) {
			continue
		}
		seen[index] = true
		names, err := restoreRegistry(index)
		if err != nil {
			return restored, err
		}
		restored = append(restored, names...)
	}
	return restored, nil
}

type registryEntry struct{ name, from string }

func restoreRegistry(index string) ([]string, error) {
	raw, err := os.ReadFile(index)
	if err != nil {
		return nil, err
	}
	src := string(raw)
	crlf := strings.Contains(src, "\r\n")
	src = strings.ReplaceAll(src, "\r\n", "\n")
	if !hasMarkerLine(src, "// grit:resources") || !hasMarkerLine(src, "// grit:resource-list") {
		return nil, nil
	}
	for _, m := range registryImport.FindAllStringSubmatch(src, -1) {
		if !templateRegistryImports[m[1]] {
			return nil, nil
		}
	}

	dir := filepath.Dir(index)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var missing []registryEntry
	for _, e := range entries {
		var file, from string
		switch {
		case e.IsDir():
			file, from = filepath.Join(dir, e.Name(), e.Name()+".ts"), RegistryImportPath(e.Name())
		case strings.HasSuffix(e.Name(), ".ts") && e.Name() != "index.ts":
			file, from = filepath.Join(dir, e.Name()), "./"+strings.TrimSuffix(e.Name(), ".ts")
		default:
			continue
		}
		body, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		for _, m := range resourceExport.FindAllStringSubmatch(string(body), -1) {
			if !regexp.MustCompile(`\b` + m[1] + `\b`).MatchString(src) {
				missing = append(missing, registryEntry{name: m[1], from: from})
			}
		}
	}
	if len(missing) == 0 {
		return nil, nil
	}
	sort.Slice(missing, func(i, j int) bool { return missing[i].name < missing[j].name })

	var names []string
	for _, m := range missing {
		src = insertBeforeMarkerLine(src, "// grit:resources", fmt.Sprintf(`import { %s } from "%s";`, m.name, m.from))
		src = insertBeforeMarkerLine(src, "// grit:resource-list", "  "+m.name+",")
		names = append(names, m.name)
	}
	if crlf {
		src = strings.ReplaceAll(src, "\n", "\r\n")
	}
	if err := os.WriteFile(index, []byte(src), 0o644); err != nil {
		return nil, err
	}
	manifest.Refresh(index)
	return names, nil
}

func hasMarkerLine(src, marker string) bool {
	return markerLineStart(src, marker) >= 0
}

// markerLineStart is where the line that holds only marker starts, or -1.
func markerLineStart(src, marker string) int {
	at := 0
	for _, line := range strings.SplitAfter(src, "\n") {
		if strings.TrimSpace(line) == marker {
			return at
		}
		at += len(line)
	}
	return -1
}

// insertBeforeMarkerLine puts code on its own line above the marker's line.
func insertBeforeMarkerLine(src, marker, code string) string {
	at := markerLineStart(src, marker)
	if at < 0 {
		return src
	}
	return src[:at] + code + "\n" + src[at:]
}
