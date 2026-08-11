package generate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

func TestCheckReservedName(t *testing.T) {
	cases := []struct {
		name    string
		wantErr bool
	}{
		{"Ticket", true},         // the collision that started this
		{"TicketReply", true},    // the one that actually broke the build
		{"User", true},           //
		{"Backup", true},         //
		{"Product", false},       // ordinary names stay free
		{"Invoice", false},       //
		{"SupportTicket", false}, // the suggested alternative must not itself be reserved
		{"ticket", false},        // model names are PascalCase; lowercase is not the same symbol
	}

	for _, tc := range cases {
		err := CheckReservedName(tc.name)
		if tc.wantErr && err == nil {
			t.Errorf("CheckReservedName(%q) = nil, want an error", tc.name)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("CheckReservedName(%q) = %v, want nil", tc.name, err)
		}
	}
}

func TestReservedNameErrorSuggestsAWayForward(t *testing.T) {
	err := CheckReservedName("Ticket")
	if err == nil {
		t.Fatal("expected an error for Ticket")
	}
	msg := err.Error()
	for _, want := range []string{"support desk", "ticket.go", "--force", "SupportTicket"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message is missing %q:\n%s", want, msg)
		}
	}
}

// The reserved list is a hand-written copy of what the scaffold emits, so the
// thing most likely to go wrong is drift: a new built-in model lands, nobody
// adds it here, and `grit generate resource <NewModel>` silently overwrites it
// exactly the way Ticket was overwritten. So scaffold a real project and check
// the two agree.
func TestReservedListCoversEveryScaffoldedModel(t *testing.T) {
	if testing.Short() {
		t.Skip("scaffolds a project; skipped under -short")
	}

	// scaffold.Run writes into the working directory, so move there and back.
	root := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	if err := scaffold.Run(scaffold.Options{
		ProjectName:  "reservedcheck",
		Architecture: scaffold.ArchAPI,
	}); err != nil {
		t.Fatalf("scaffold: %v", err)
	}

	modelsDir := filepath.Join(root, "reservedcheck", "internal", "models")
	if _, err := os.Stat(modelsDir); err != nil {
		modelsDir = filepath.Join(root, "reservedcheck", "apps", "api", "internal", "models")
	}
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		t.Fatalf("reading %s: %v", modelsDir, err)
	}

	structRe := regexp.MustCompile(`(?m)^type ([A-Z][A-Za-z0-9]*) struct`)

	// Demo resources are meant to be removed and replaced, so they are
	// deliberately NOT reserved: `grit remove resource Blog` is a supported
	// first step on a new project.
	demo := map[string]bool{"Blog": true}

	// Not every struct in models/ is a persisted model — request/response
	// shapes live there too. Only names that AutoMigrate is given can be
	// clobbered in the way that matters, so filter to those.
	migrated := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		body, err := os.ReadFile(filepath.Join(modelsDir, e.Name()))
		if err != nil {
			t.Fatalf("reading %s: %v", e.Name(), err)
		}
		for _, m := range regexp.MustCompile(`&([A-Z][A-Za-z0-9]*)\{\},`).FindAllStringSubmatch(string(body), -1) {
			migrated[m[1]] = true
		}
	}

	var missing []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		body, err := os.ReadFile(filepath.Join(modelsDir, e.Name()))
		if err != nil {
			t.Fatalf("reading %s: %v", e.Name(), err)
		}
		for _, m := range structRe.FindAllStringSubmatch(string(body), -1) {
			name := m[1]
			if demo[name] || !migrated[name] {
				continue
			}
			if _, reserved := reservedModels[name]; !reserved {
				missing = append(missing, name+" ("+e.Name()+")")
			}
		}
	}

	if len(missing) > 0 {
		t.Errorf("these models are scaffolded but not reserved, so `grit generate resource <Name>`\n"+
			"would overwrite them the way Ticket was overwritten — add them to reservedModels:\n  %s",
			strings.Join(missing, "\n  "))
	}
}
