package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeReactPkg(t *testing.T, path, version string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "{\n  \"dependencies\": {\n    \"react\": \"" + version + "\",\n    \"react-dom\": \"" + version + "\"\n  }\n}\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func reactPinOf(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, `"react-dom": "`) {
		t.Fatalf("%s lost react-dom", path)
	}
	react := s[strings.Index(s, `"react": "`)+len(`"react": "`):]
	dom := s[strings.Index(s, `"react-dom": "`)+len(`"react-dom": "`):]
	react, dom = react[:strings.Index(react, `"`)], dom[:strings.Index(dom, `"`)]
	if react != dom {
		t.Fatalf("%s pins react %s and react-dom %s", path, react, dom)
	}
	return react
}

// With the Expo app in the monorepo, every web frontend uses Expo's React, so
// the hoisted node_modules holds one copy and component tests render.
func TestAlignReactVersionsFollowsExpo(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web", "package.json")
	admin := filepath.Join(root, "apps", "admin", "package.json")
	desktop := filepath.Join(root, "apps", "desktop", "frontend", "package.json")
	expo := filepath.Join(root, "apps", "expo", "package.json")
	for _, p := range []string{web, admin, desktop} {
		writeReactPkg(t, p, webReactVersion)
	}
	writeReactPkg(t, expo, expoReactVersion)

	if changed := alignReactVersions(root); len(changed) != 3 {
		t.Fatalf("changed %d files, want 3: %v", len(changed), changed)
	}
	for _, p := range []string{web, admin, desktop, expo} {
		if got := reactPinOf(t, p); got != expoReactVersion {
			t.Errorf("%s pins %s, want %s", p, got, expoReactVersion)
		}
	}
	if changed := alignReactVersions(root); len(changed) != 0 {
		t.Errorf("a second run changed %v", changed)
	}
}

// Without Expo, the web frontends stay on (or return to) the web version.
func TestAlignReactVersionsWithoutExpo(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web", "package.json")
	admin := filepath.Join(root, "apps", "admin", "package.json")
	writeReactPkg(t, web, expoReactVersion)
	writeReactPkg(t, admin, "^19.0.0")

	alignReactVersions(root)
	for _, p := range []string{web, admin} {
		if got := reactPinOf(t, p); got != webReactVersion {
			t.Errorf("%s pins %s, want %s", p, got, webReactVersion)
		}
	}
}
