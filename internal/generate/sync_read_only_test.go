package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit#210: a model the server authors on a user's behalf had no mode.
//
// online_only means she cannot read her own report on a train. offline_first
// declares the rows a client's to create, so the outbox enqueues a write that
// is guaranteed to fail and parks it in front of her. Conflict: server_wins is
// not the same thing: it decides what happens when VERSIONS disagree, and a
// client-minted create has no version to disagree about.

func TestReadOnlyIsADeclarablePolicyMode(t *testing.T) {
	p := &SyncPolicy{Mode: "read_only"}
	if err := p.Validate("CostLedger"); err != nil {
		t.Fatalf("read_only was refused at parse time: %v", err)
	}
	if got := goConst("read_only"); got != "ModeReadOnly" {
		t.Errorf("goConst(read_only) = %q, want ModeReadOnly", got)
	}
}

// A typo is still refused, and the message lists what is allowed.
func TestAnUnknownModeIsStillRefused(t *testing.T) {
	p := &SyncPolicy{Mode: "readonly"} // the spelling somebody will try
	err := p.Validate("CostLedger")
	if err == nil {
		t.Fatal("readonly was accepted; only read_only is a mode")
	}
	if !strings.Contains(err.Error(), "read_only") {
		t.Errorf("the error does not name the real spelling: %v", err)
	}
}

// --------------------------------------------------------------- the doctor

// The check grit#210 asked for: the first one about what the SERVER will
// accept rather than whether the policy is internally coherent.
//
// Before it, a project with seven append-only tables declared offline_first
// was told "everything checks out", which is the state the sync docs
// themselves call the worst one: a declaration that reads as though it took
// effect.
func TestDoctorFlagsAnAppendOnlyModelDeclaredOfflineFirst(t *testing.T) {
	known := map[string]bool{"id": true, "version": true, "updated_at": true, "deleted_at": true}

	findings := checkPolicyLiteral("audit_entries",
		"sync.Policy{Mode: sync.ModeOfflineFirst, Conflict: sync.ConflictManual}", known, true)

	var found *SyncFinding
	for i := range findings {
		if strings.Contains(findings[i].Message, "append-only") {
			found = &findings[i]
		}
	}
	if found == nil {
		t.Fatalf("an append-only model declared offline_first was not flagged: %+v", findings)
	}
	if found.Level != "warning" {
		// Not an error: create is legal on an append-only model. A client
		// minting a ledger row is a real design, and only its updates and
		// deletes are certain to fail.
		t.Errorf("level = %q, want warning", found.Level)
	}
	if !strings.Contains(found.Fix, "read_only") {
		t.Errorf("the fix does not name the mode to use: %q", found.Fix)
	}
}

// The same model declared read_only is correct and must not be nagged about.
// A linter that cries wolf gets turned off.
func TestDoctorIsQuietAboutAnAppendOnlyModelDeclaredReadOnly(t *testing.T) {
	known := map[string]bool{"id": true, "version": true, "updated_at": true, "deleted_at": true}

	findings := checkPolicyLiteral("audit_entries",
		"sync.Policy{Mode: sync.ModeReadOnly, Conflict: sync.ConflictServerWins}", known, true)

	for _, f := range findings {
		if strings.Contains(f.Message, "append-only") {
			t.Errorf("a correctly declared model was still flagged: %q", f.Message)
		}
		if f.Level == "warning" || f.Level == "error" {
			t.Errorf("unexpected %s: %q", f.Level, f.Message)
		}
	}
	// It does say what read_only means, once, at info.
	var said bool
	for _, f := range findings {
		if strings.Contains(f.Message, "read_only") {
			said = true
		}
	}
	if !said {
		t.Error("read_only is not explained at all, so a reader of the report cannot tell " +
			"whether it is what they meant")
	}
}

// An ordinary model is unaffected, which is every existing project.
func TestDoctorSaysNothingNewAboutAnOrdinaryModel(t *testing.T) {
	known := map[string]bool{"id": true, "version": true, "updated_at": true, "deleted_at": true}

	findings := checkPolicyLiteral("grants",
		"sync.Policy{Mode: sync.ModeOfflineFirst, Conflict: sync.ConflictManual}", known, false)

	for _, f := range findings {
		t.Errorf("unexpected finding for a plain offline_first model: %s %q", f.Level, f.Message)
	}
}

// readAppendOnlyModels reads the registration the generator writes, because
// doctor runs against a project on disk where the definition that produced it
// is gone.
func TestReadAppendOnlyModelsFindsTheRegistration(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("audit_entry.go", "package models\n\nfunc init() { appendonly.Register(&AuditEntry{}) }\n")
	write("grant.go", "package models\n\ntype Grant struct{}\n")
	// A test file is not the project's model registry.
	write("audit_entry_test.go", "package models\n\nfunc init() { appendonly.Register(&NotAModel{}) }\n")

	got := readAppendOnlyModels(dir)
	if !got["AuditEntry"] {
		t.Error("the append-only registration was not found")
	}
	if got["Grant"] {
		t.Error("a plain model was reported as append-only")
	}
	if got["NotAModel"] {
		t.Error("a _test.go file was read as part of the model registry")
	}
}
