package deploy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A dry run says what it would do, and does none of it.
//
// --dry-run reached the Railway path only. On the SSH path it was accepted and
// ignored: a flag whose help reads "print every command that would run, and run
// none of them" built the binary and opened an SSH connection to the host, then
// printed "Deployment successful!" when that failed on an unreachable name.
//
// The plan is checked rather than the absence of side effects, because a test
// that asserts nothing happened passes just as well when the function returns
// early for the wrong reason.
func TestDryRunPrintsThePlanForAMonorepo(t *testing.T) {
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api")
	if err := os.MkdirAll(filepath.Join(api, "cmd", "server"), 0o755); err != nil {
		t.Fatal(err)
	}

	out := capture(t, func() error {
		return Run(Config{
			Host:    "user@example.com",
			AppName: "contacts",
			Domain:  "contacts.example.com",
			APIDir:  api,
			DryRun:  true,
		})
	})

	for _, want := range []string{
		"Nothing below has been run.",
		"go build -o bin/contacts ./cmd/server",
		"ssh user@example.com 'mkdir -p /opt/contacts'",
		"scp bin/contacts user@example.com:/opt/contacts/",
		"contacts.service",
		"contacts.example.com -> :8080",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the plan does not mention %q:\n%s", want, out)
		}
	}

	// A monorepo's frontends are Next.js apps with their own hosting, so there
	// is nothing to build here for the binary being uploaded.
	if strings.Contains(out, "pnpm build") {
		t.Errorf("the plan builds a frontend the binary does not carry:\n%s", out)
	}
}

// A single project's binary carries its frontend, so the plan builds it first:
// //go:embed reads what is on disk, and a binary built before the frontend
// carries the placeholder.
func TestDryRunBuildsTheFrontendForAnEmbeddedSingle(t *testing.T) {
	root := t.TempDir()
	api := filepath.Join(root, "api")
	if err := os.MkdirAll(filepath.Join(api, "web"), 0o755); err != nil {
		t.Fatal(err)
	}

	out := capture(t, func() error {
		return Run(Config{
			Host:    "user@example.com",
			AppName: "contacts",
			APIDir:  api,
			WebDir:  root,
			DryRun:  true,
		})
	})

	if !strings.Contains(out, "pnpm build") {
		t.Errorf("the plan does not build the frontend the binary embeds:\n%s", out)
	}
	// No cmd/server here: a single project's main package is the top of its
	// module, and "./cmd/server" is a package that does not exist.
	if !strings.Contains(out, "go build -o bin/contacts .") {
		t.Errorf("the plan builds the wrong package:\n%s", out)
	}
}

// The host is still required, dry run or not: a plan naming no server is a plan
// for nothing, and the error is the useful answer.
func TestDryRunStillRequiresAHost(t *testing.T) {
	err := Run(Config{AppName: "contacts", APIDir: t.TempDir(), DryRun: true})
	if err == nil {
		t.Error("a deploy with no host was accepted")
	}
}

// capture runs fn with stdout redirected, and returns what it printed.
func capture(t *testing.T, fn func() error) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 0, 4096)
		chunk := make([]byte, 1024)
		for {
			n, readErr := r.Read(chunk)
			buf = append(buf, chunk[:n]...)
			if readErr != nil {
				break
			}
		}
		done <- string(buf)
	}()

	runErr := fn()
	w.Close()
	os.Stdout = old
	out := <-done

	if runErr != nil {
		t.Fatalf("dry run returned an error: %v\n%s", runErr, out)
	}
	return out
}
