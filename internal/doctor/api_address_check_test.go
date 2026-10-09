package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeEnvLocal puts a .env.local under apps/<app>.
func writeEnvLocal(t *testing.T, root, app, body string) {
	t.Helper()
	dir := filepath.Join(root, "apps", app)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env.local"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAPIAddressAgreesIsQuietWhenEverythingMatches(t *testing.T) {
	root := t.TempDir()
	writeEnvLocal(t, root, "admin", "NEXT_PUBLIC_API_URL=http://localhost:8099\n")
	p := &project{
		root: root,
		env: map[string]string{
			"APP_PORT": "8099",
			"APP_URL":  "http://localhost:8099",
			"API_URL":  "http://localhost:8099",
		},
	}
	if findings := checkAPIAddressAgrees(p); len(findings) != 0 {
		t.Fatalf("expected nothing to report, got %v", findings)
	}
}

// The case that cost an afternoon: APP_PORT and APP_URL moved together, so the
// existing check passed, and API_URL stayed behind.
func TestAPIURLLeftBehindIsReported(t *testing.T) {
	p := &project{
		root: t.TempDir(),
		env: map[string]string{
			"APP_PORT": "8099",
			"APP_URL":  "http://localhost:8099",
			"API_URL":  "http://localhost:8080",
		},
	}
	findings := checkAPIAddressAgrees(p)
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d: %v", len(findings), findings)
	}
	if findings[0].Resource != "API_URL" {
		t.Errorf("the finding names %q, not API_URL", findings[0].Resource)
	}
	if !strings.Contains(findings[0].Fix, "8099") {
		t.Errorf("the fix does not name the port the API listens on: %q", findings[0].Fix)
	}
}

// The seeder writes .env.local once and never again, so this is the copy that
// goes stale without anybody editing a file.
func TestAStaleEnvLocalIsReported(t *testing.T) {
	root := t.TempDir()
	writeEnvLocal(t, root, "web", "# written by the seeder\nNEXT_PUBLIC_API_URL=http://localhost:8080\n")
	writeEnvLocal(t, root, "admin", "NEXT_PUBLIC_API_URL=http://localhost:8080\n")
	p := &project{
		root: root,
		env: map[string]string{
			"APP_PORT": "8099",
			"APP_URL":  "http://localhost:8099",
			"API_URL":  "http://localhost:8099",
		},
	}
	findings := checkAPIAddressAgrees(p)
	if len(findings) != 2 {
		t.Fatalf("expected both apps reported, got %d: %v", len(findings), findings)
	}
	for _, f := range findings {
		if !strings.HasSuffix(f.Resource, ".env.local") {
			t.Errorf("unexpected resource %q", f.Resource)
		}
	}
}

// A deployed origin is behind a proxy and has no port to compare.
func TestADeployedURLIsNotComparedToAPPPORT(t *testing.T) {
	p := &project{
		root: t.TempDir(),
		env: map[string]string{
			"APP_PORT": "8099",
			"API_URL":  "https://api.example.com",
		},
	}
	if findings := checkAPIAddressAgrees(p); len(findings) != 0 {
		t.Fatalf("a production URL was reported: %v", findings)
	}
}

// Nothing to compare against means nothing to say.
func TestNoAppPortMeansNoFinding(t *testing.T) {
	p := &project{
		root: t.TempDir(),
		env:  map[string]string{"API_URL": "http://localhost:8080"},
	}
	if findings := checkAPIAddressAgrees(p); len(findings) != 0 {
		t.Fatalf("expected nothing without APP_PORT, got %v", findings)
	}
}
