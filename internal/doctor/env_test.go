package doctor

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/clierr"
)

// A port that something is listening on must be reported as busy.
//
// The check is the one thing here that could plausibly be written backwards:
// dialling a port instead of listening on it reports "busy" whenever anything
// refuses the connection and "free" whenever a server is slow to accept. So
// the test opens a real listener and asks.
func TestCheckPortSeesABusyPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not open a listener to test against: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	busy := checkPort(PortCheck{Port: port, What: "the test"})
	if busy.Level != EnvWarn {
		t.Fatalf("a port with a listener on it reported %q, want %q", busy.Level, EnvWarn)
	}
	if busy.Problem == nil || busy.Problem.Slug != "port-in-use" {
		t.Errorf("a busy port must point at the port-in-use page, got %+v", busy.Problem)
	}
	if !strings.Contains(busy.Detail, "the test") {
		t.Errorf("the detail should name what wants the port, got %q", busy.Detail)
	}

	// And the same port, once nothing is listening, must be free. Closing
	// and re-checking is what proves the check reads the port rather than
	// returning a constant.
	if err := ln.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	free := checkPort(PortCheck{Port: port, What: "the test"})
	if free.Level != EnvOK {
		t.Fatalf("after closing the listener the port reported %q, want %q: %s",
			free.Level, EnvOK, free.Detail)
	}
}

// Every finding that is not OK must carry a problem the catalogue knows, or
// the CLI prints a link to a page that is not there.
func TestEveryEnvFailureHasAKnownProblem(t *testing.T) {
	report := CheckEnv(context.Background(), EnvOptions{
		Ports: []PortCheck{{Port: busyPort(t), What: "the test"}},
	})

	for _, f := range report.Findings {
		if f.Level == EnvOK {
			continue
		}
		if f.Problem == nil {
			// A warning with no page is allowed only where the detail is
			// already the whole story: an unreadable version string.
			if strings.Contains(f.Detail, "could not read the version") {
				continue
			}
			t.Errorf("%s is %q with no problem to read: %q", f.Name, f.Level, f.Detail)
			continue
		}
		if _, ok := clierr.BySlug(f.Problem.Slug); !ok {
			t.Errorf("%s points at %q, which is not in the catalogue", f.Name, f.Problem.Slug)
		}
	}
}

// The version comparison, which is easy to get backwards.
func TestOlder(t *testing.T) {
	for _, tc := range []struct {
		a, b string
		want bool
	}{
		{"1.26.10", "1.26.11", true},
		{"1.26.11", "1.26.11", false},
		{"1.26.12", "1.26.11", false},
		{"1.25.99", "1.26.0", true},
		{"2.0.0", "1.99.99", false},
	} {
		aMaj, aMin, aPatch := parseVersion(tc.a)
		bMaj, bMin, bPatch := parseVersion(tc.b)
		if got := older(aMaj, aMin, aPatch, bMaj, bMin, bPatch); got != tc.want {
			t.Errorf("older(%s, %s) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func busyPort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	return ln.Addr().(*net.TCPAddr).Port
}

// Every problem in the catalogue has the parts a message needs. An entry with
// an empty Fix prints a blank line under the failure, which reads as a bug in
// the CLI rather than as a gap in the catalogue.
func TestCatalogueEntriesAreComplete(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range clierr.All {
		if seen[p.Slug] {
			t.Errorf("duplicate slug %q", p.Slug)
		}
		seen[p.Slug] = true

		for _, field := range []struct{ name, value string }{
			{"Slug", p.Slug}, {"Title", p.Title}, {"Cause", p.Cause}, {"Fix", p.Fix},
		} {
			if strings.TrimSpace(field.value) == "" {
				t.Errorf("%s: %s is empty", p.Slug, field.name)
			}
		}
		if want := fmt.Sprintf("https://gritframework.dev/e/%s", p.Slug); p.URL() != want {
			t.Errorf("%s: URL() = %q, want %q", p.Slug, p.URL(), want)
		}
	}
}
