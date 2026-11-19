package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// wailsVersionRef finds every place a Wails version is named.
//
// Both spellings: the module requirement in a go.mod and the @version on a
// CLI install line. A doc table and a README use the second.
var wailsVersionRef = regexp.MustCompile(`wailsapp/wails/v2(?:/cmd/wails@|\s+)(v2\.[0-9]+\.[0-9]+)`)

// Every Wails version in the templates is the same version.
//
// TestTheWailsVersionsAgree compares two of them, the desktop client's go.mod
// and the release workflow. That is how desktopGoMod, the go.mod for a
// standalone desktop project, sat at v2.9.1 while the others were on v2.16.0:
// eight minor versions behind, and the test said the versions agreed.
//
// It mattered more than a warning on a build. v2.17.0 closed a dev-server hole
// where a page in the developer's browser could reach ws://localhost:34115 and
// invoke the whole bound Go API, and v2.17.1 fixed GHSA-r5j2-ghh9-2jqw, high
// severity at CVSS 8.3, where a Linux iframe bypassed OriginValidator. A
// project scaffolded from the forgotten template had both.
//
// So this counts rather than compares a pair: a sixth site added later is held
// to the same version without anybody remembering to add it here.
func TestEveryWailsVersionInTheTemplatesAgrees(t *testing.T) {
	found := map[string][]string{} // version -> where

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		name := info.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		for _, match := range wailsVersionRef.FindAllStringSubmatch(string(body), -1) {
			found[match[1]] = append(found[match[1]], path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the templates: %v", err)
	}

	if len(found) == 0 {
		t.Fatal("no Wails version found anywhere, so this guard is watching nothing")
	}

	versions := make([]string, 0, len(found))
	for version := range found {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	if len(versions) > 1 {
		for _, version := range versions {
			t.Errorf("Wails %s in %v", version, found[version])
		}
		t.Errorf("the templates name %d different Wails versions. Wails warns on "+
			"every build when the CLI and the module differ, and the version left "+
			"behind is the one that keeps a security fix out of a project.",
			len(versions))
	}

	// And the count, because a site that stops matching the pattern is as bad
	// as one that disagrees: it drifts silently instead of loudly.
	const expected = 5
	total := 0
	for _, places := range found {
		total += len(places)
	}
	if total != expected {
		t.Errorf("found %d places naming a Wails version, expected %d. If you added "+
			"or removed one, change the number here and say why; if this dropped on "+
			"its own, a template stopped matching and is no longer held to the "+
			"version.", total, expected)
	}
}

// The version is at or past the one that carries the security fixes.
//
// A floor rather than an exact match, so a later bump does not fail. It names
// the advisories because "upgrade Wails" is not an argument and
// "GHSA-r5j2-ghh9-2jqw, 8.3" is.
func TestTheWailsVersionIsNotVulnerable(t *testing.T) {
	const floor = "v2.17.1"

	body, err := os.ReadFile(filepath.Join("desktop_client_files.go"))
	if err != nil {
		t.Fatal(err)
	}
	match := wailsVersionRef.FindStringSubmatch(string(body))
	if match == nil {
		t.Fatal("the desktop go.mod names no Wails version")
	}

	if compareSemver(match[1], floor) < 0 {
		t.Errorf("Wails %s is pinned, and %s is the first release without "+
			"GHSA-r5j2-ghh9-2jqw (high, CVSS 8.3: a Linux iframe bypasses "+
			"OriginValidator) or the v2.17.0 dev-server hole where a page in the "+
			"developer's browser can invoke the whole bound Go API.",
			match[1], floor)
	}
}

// compareSemver is enough for v2.x.y: -1, 0 or 1.
func compareSemver(a, b string) int {
	parse := func(v string) [3]int {
		var out [3]int
		for i, part := range strings.SplitN(strings.TrimPrefix(v, "v"), ".", 3) {
			if i > 2 {
				break
			}
			n := 0
			for _, ch := range part {
				if ch < '0' || ch > '9' {
					break
				}
				n = n*10 + int(ch-'0')
			}
			out[i] = n
		}
		return out
	}
	left, right := parse(a), parse(b)
	for i := 0; i < 3; i++ {
		switch {
		case left[i] < right[i]:
			return -1
		case left[i] > right[i]:
			return 1
		}
	}
	return 0
}
