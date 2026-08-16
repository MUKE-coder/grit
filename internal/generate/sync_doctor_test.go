package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// syncProject writes a minimal project: one routes.go and the model files it
// registers.
func syncProject(t *testing.T, routes string, models map[string]string) string {
	t.Helper()
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api", "internal")

	routesDir := filepath.Join(api, "routes")
	if err := os.MkdirAll(routesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(routesDir, "routes.go"), []byte(routes), 0o644); err != nil {
		t.Fatal(err)
	}

	modelsDir := filepath.Join(api, "models")
	if err := os.MkdirAll(modelsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range models {
		if err := os.WriteFile(filepath.Join(modelsDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const healthySale = `package models

import (
	"time"

	"gorm.io/gorm"
)

type Sale struct {
	ID            string         ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Total         float64        ` + "`" + `json:"total"` + "`" + `
	PaymentMethod string         ` + "`" + `json:"payment_method"` + "`" + `
	DraftNote     string         ` + "`" + `json:"draft_note"` + "`" + `
	Version       int            ` + "`" + `json:"version"` + "`" + `
	CreatedAt     time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt     time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt     gorm.DeletedAt ` + "`" + `json:"deleted_at"` + "`" + `
}
`

func findingsFor(report *SyncReport, level string) []SyncFinding {
	var out []SyncFinding
	for _, f := range report.Findings {
		if f.Level == level {
			out = append(out, f)
		}
	}
	return out
}

func TestSyncDoctorPassesAHealthyProject(t *testing.T) {
	root := syncProject(t,
		`package routes

func Setup() {
	syncRegistry.Register("sales", &models.Sale{})
}
`,
		map[string]string{"sale.go": healthySale})

	report, err := SyncDoctor(root)
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if report.Errors() != 0 {
		t.Fatalf("healthy project reported errors: %+v", findingsFor(report, "error"))
	}
	if len(report.Registered) != 1 || report.Registered[0] != "sales" {
		t.Fatalf("registered = %v", report.Registered)
	}
	// Without this the test passes when no model file is read at all, which is
	// exactly how it passed the first time it was run.
	if len(report.Models) == 0 {
		t.Fatal("no models were parsed, so this passed by finding nothing to check")
	}
}

// The failure this whole command exists for. An allowlist entry that matches
// no column does not error anywhere: it excludes the real column, and every
// client mirrors rows with the value missing.
func TestSyncDoctorCatchesAnAllowlistTypo(t *testing.T) {
	root := syncProject(t,
		`package routes

func Setup() {
	syncRegistry.RegisterWithPolicy("sales", &models.Sale{}, sync.Policy{Mode: sync.ModeOfflineFirst, Conflict: sync.ConflictServerWins, Fields: []string{"totl", "payment_method"}})
}
`,
		map[string]string{"sale.go": healthySale})

	report, err := SyncDoctor(root)
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}

	errors := findingsFor(report, "error")
	if len(errors) != 1 {
		t.Fatalf("want exactly one error, got %+v", errors)
	}
	if !strings.Contains(errors[0].Message, "totl") {
		t.Errorf("the error should name the typo: %q", errors[0].Message)
	}
	if errors[0].Model != "sales" {
		t.Errorf("the error should name the model, got %q", errors[0].Model)
	}

	// server_wins is worth saying out loud, since the user is never asked.
	var sawConflictNote bool
	for _, f := range findingsFor(report, "info") {
		if strings.Contains(f.Message, "server_wins") {
			sawConflictNote = true
		}
	}
	if !sawConflictNote {
		t.Error("server_wins should be reported, because nobody is prompted when it fires")
	}
}

// A model with no Version column cannot detect a conflict at all: every push
// is accepted and the last writer wins, silently.
func TestSyncDoctorCatchesAModelThatCannotVersion(t *testing.T) {
	root := syncProject(t,
		`package routes

func Setup() {
	syncRegistry.Register("notes", &models.Note{})
}
`,
		map[string]string{"note.go": `package models

import "time"

type Note struct {
	ID        string    ` + "`" + `json:"id"` + "`" + `
	Body      string    ` + "`" + `json:"body"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}
`})

	report, err := SyncDoctor(root)
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}

	var sawVersion bool
	for _, f := range findingsFor(report, "error") {
		if strings.Contains(f.Message, "Version") {
			sawVersion = true
		}
	}
	if !sawVersion {
		t.Fatalf("a model with no Version should be an error, got %+v", report.Findings)
	}

	// And no soft delete means deletes never reach a mirror as tombstones.
	var sawTombstone bool
	for _, f := range findingsFor(report, "warning") {
		if strings.Contains(f.Message, "tombstone") {
			sawTombstone = true
		}
	}
	if !sawTombstone {
		t.Error("a model with no soft delete should warn about tombstones")
	}
}

func TestSyncDoctorReportsAProjectWithNothingRegistered(t *testing.T) {
	root := syncProject(t, "package routes\n\nfunc Setup() {}\n", nil)

	report, err := SyncDoctor(root)
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if report.Errors() != 0 {
		t.Errorf("nothing registered is not an error: %+v", findingsFor(report, "error"))
	}
	if report.Warnings() == 0 {
		t.Error("nothing registered should be worth mentioning")
	}
}

func TestSyncDoctorNeedsAProject(t *testing.T) {
	if _, err := SyncDoctor(t.TempDir()); err == nil {
		t.Fatal("expected an error outside a Grit project")
	}
}
