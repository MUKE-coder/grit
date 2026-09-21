package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A new project's Next.js CSP names the API's socket origin, and the Vite
// frontends' nginx CSP allows sockets at all.
func TestNewFrontendCSPAdmitsTheRealtimeSocket(t *testing.T) {
	next := nextSecurityHeaders()
	for _, want := range []string{nextAPIWSOrigin, nextConnectSrcNew} {
		if !strings.Contains(next, want) {
			t.Errorf("the Next.js CSP is missing %q", want)
		}
	}
	if strings.Contains(next, nextConnectSrcOld) {
		t.Error("the Next.js CSP still has the connect-src without the socket origin")
	}
	if nginx := viteNginxConf("web"); !strings.Contains(nginx, nginxConnectSrcNew) {
		t.Error("the nginx CSP does not allow ws: and wss:")
	}
}

// The repair turns the config Grit wrote before into the one it writes now,
// once, and leaves a config it did not write alone with a warning.
func TestCSPWebSocketRepair(t *testing.T) {
	now := nextSecurityHeaders()
	before := strings.Replace(strings.Replace(now, nextAPIWSOrigin, "", 1), nextConnectSrcNew, nextConnectSrcOld, 1)
	if before == now {
		t.Fatal("could not rebuild the old config from the template")
	}
	got, fixed, warn := repairCSPWebSocketNextSource(before)
	if got != now || len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("repair did not produce the current config (fixed %v, warned %v)", fixed, warn)
	}
	if again, fixed, _ := repairCSPWebSocketNextSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}

	custom := "const csp = [\"connect-src 'self' https://mine.example\"];\n"
	if out, _, warn := repairCSPWebSocketNextSource(custom); out != custom || len(warn) != 1 {
		t.Errorf("a hand-written CSP was changed or not warned about: %q %v", out, warn)
	}

	oldNginx := strings.Replace(viteNginxConf("web"), nginxConnectSrcNew, nginxConnectSrcOld, 1)
	if out, fixed, _ := repairCSPWebSocketNginxSource(oldNginx); out != viteNginxConf("web") || len(fixed) != 1 {
		t.Error("the nginx repair did not produce the current config")
	}
}

func TestBuildArtefactsAreIgnoredOnce(t *testing.T) {
	if !strings.Contains(rootGitignore(), "*.tsbuildinfo\nnext-env.d.ts\n") {
		t.Error("a new project's .gitignore does not ignore *.tsbuildinfo and next-env.d.ts")
	}
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	if err := os.WriteFile(path, []byte("node_modules/\r\n.env\r\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := ignoreBuildArtefacts(root); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(got), "*.tsbuildinfo"); n != 1 {
		t.Errorf("*.tsbuildinfo appears %d times after two upgrades, want 1", n)
	}
	if strings.Contains(strings.ReplaceAll(string(got), "\r\n", ""), "\n") {
		t.Error("a CRLF .gitignore got LF lines")
	}
}
