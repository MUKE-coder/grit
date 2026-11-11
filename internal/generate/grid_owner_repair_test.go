package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The hole the repair closes, and the proof it is a hole: these are the exact
// lines v3.504.0 generated for `--owned-by user`.
//
// Built from the generator rather than typed out, so the test cannot pass
// against a shape the generator stopped producing.
func ownedService(t *testing.T, confidential bool) string {
	t.Helper()
	def := &ResourceDefinition{
		Name: "Memo",
		Fields: []Field{
			{Name: "user", Type: "belongs_to", RelatedModel: "User"},
			{Name: "body", Type: "string"},
		},
		OwnedBy:      "user",
		Confidential: confidential,
	}
	g := &Generator{Definition: def, Module: "app"}
	return g.serviceSource(g.Names())
}

// ownedServiceAsOf504 is the same service as v3.504.0 wrote it: owned, with a
// grid that ignores ownership. This is what a real project has on disk.
//
// Built by putting the pre-fix blocks back, using the same gridRowsLoad the
// generator used to take them out. Spelling the old text out by hand here
// would make this fixture the thing that drifts.
func ownedServiceAsOf504(t *testing.T, confidential bool) string {
	t.Helper()
	src := ownedService(t, confidential)

	scoped := gridRowsLoad("Memo", "memos", ownerGridScopeSource("user", "memo", gridScopeCallFor(confidential)))
	if !strings.Contains(src, scoped) {
		t.Fatal("the generator no longer writes the scoped row load this fixture undoes")
	}
	src = strings.Replace(src, scoped, gridRowsLoad("Memo", "memos", ""), 1)

	stamp := ownerGridStampSource("user", "user_id", "UserID")
	if !strings.Contains(src, stamp) {
		t.Fatal("the generator no longer writes the owner stamp this fixture undoes")
	}
	src = strings.Replace(src, stamp, "", 1)
	return src
}

func gridScopeCallFor(confidential bool) string {
	if confidential {
		return `authz.ScopeConfidential(ctx, rows, "user_id", "memos.confidential")`
	}
	return `authz.ScopeOwned(ctx, rows, "user_id")`
}

// The one that matters: repairing an old service produces the file the
// generator would write today, to the byte.
//
// Not "contains the scope". A repaired service that differs from a generated
// one by a comment is a service nobody can reason about, and the two drift the
// moment either text is written twice: that is what
// TestRepairedBulkEditIsTheGeneratedBulkEdit caught about this very change.
func TestARepairedServiceIsAGeneratedService(t *testing.T) {
	for _, confidential := range []bool{false, true} {
		name := "owned"
		if confidential {
			name = "confidential"
		}
		t.Run(name, func(t *testing.T) {
			repaired, note := repairGridOwnerSource(ownedServiceAsOf504(t, confidential), "Memo", "memo")
			if note != "" {
				t.Fatalf("the repair declined a service Grit wrote: %s", note)
			}
			want := ownedService(t, confidential)
			if repaired != want {
				t.Errorf("a repaired service is not a generated one.\n%s", firstDifference(repaired, want))
			}
		})
	}
}

// firstDifference is the line the two first disagree on, with context: a diff
// of two 800-line files says nothing useful.
func firstDifference(got, want string) string {
	g, w := strings.Split(got, "\n"), strings.Split(want, "\n")
	for i := 0; i < len(g) && i < len(w); i++ {
		if g[i] != w[i] {
			return fmt.Sprintf("line %d:\n  repaired:  %q\n  generated: %q", i+1, g[i], w[i])
		}
	}
	return fmt.Sprintf("same %d lines, then one is longer: repaired %d lines, generated %d", min(len(g), len(w)), len(g), len(w))
}

// A service generated before v3.505.0 gets both halves back.
func TestTheRepairScopesAnOldGridAndStampsItsCreates(t *testing.T) {
	before := ownedServiceAsOf504(t, false)

	after, note := repairGridOwnerSource(before, "Memo", "memo")
	if note != "" {
		t.Fatalf("the repair declined a service Grit itself wrote: %s", note)
	}
	if !strings.Contains(after, `rows = authz.ScopeOwned(ctx, rows, "user_id")`) {
		t.Error("BulkEdit still loads the rows the request named with no ownership scope, so a " +
			"caller holding memos.edit can rewrite another user's row by its id")
	}
	if !strings.Contains(after, "\t\titem.UserID = authz.UserIDFrom(ctx)") {
		t.Error("BulkCreate still takes the owner from the row, so the grid can file rows under " +
			"another account")
	}
	// The scope has to be before the read, not after it.
	scope := strings.Index(after, "rows = authz.ScopeOwned(")
	read := strings.Index(after, "rows.Find(&allowed)")
	if scope < 0 || read < 0 || scope > read {
		t.Errorf("the scope is not applied before the rows are read (scope at %d, read at %d)", scope, read)
	}
}

// A --confidential resource gets the stricter scope, not the plain one. The
// plain one exempts an ADMIN, which is the whole thing that flag removes.
func TestTheRepairKeepsAConfidentialResourceConfidential(t *testing.T) {
	after, note := repairGridOwnerSource(ownedServiceAsOf504(t, true), "Memo", "memo")
	if note != "" {
		t.Fatalf("declined: %s", note)
	}
	if strings.Contains(after, "rows = authz.ScopeOwned(") {
		t.Error("the repair gave a --confidential resource the scope that exempts an ADMIN, which " +
			"hands every row to the operator the flag exists to keep out")
	}
	if !strings.Contains(after, `rows = authz.ScopeConfidential(ctx, rows, "user_id", "memos.confidential")`) {
		t.Error("the confidential scope is missing from BulkEdit")
	}
}

// Running twice changes nothing the second time. An upgrade is run repeatedly,
// and a repair that stacks its own output corrupts the file it is fixing.
func TestTheRepairIsIdempotent(t *testing.T) {
	once, _ := repairGridOwnerSource(ownedServiceAsOf504(t, false), "Memo", "memo")
	twice, note := repairGridOwnerSource(once, "Memo", "memo")
	if note != "" {
		t.Fatalf("the second pass declined its own output: %s", note)
	}
	if twice != once {
		t.Error("a second pass changed the file again, so every upgrade adds another scope line")
	}
}

// A shared resource has no owner to scope by, and must come back untouched
// rather than with a scope naming a column it has not got.
func TestTheRepairLeavesASharedResourceAlone(t *testing.T) {
	g := &Generator{
		Definition: &ResourceDefinition{Name: "Memo", Fields: []Field{{Name: "body", Type: "string"}}},
		Module:     "app",
	}
	src := g.serviceSource(g.Names())
	after, note := repairGridOwnerSource(src, "Memo", "memo")
	if note != "" {
		t.Errorf("a shared resource was reported as a problem: %s", note)
	}
	if after != src {
		t.Error("a shared resource was rewritten: it has no owner column, so any scope added here " +
			"names a column that does not exist and the project stops compiling")
	}
}

// A service somebody has rewritten is reported, not patched. The message has to
// say the service is the vulnerable one: a warning nobody acts on is worse here
// than a change nobody asked for.
func TestTheRepairReportsAServiceItCannotRead(t *testing.T) {
	before := ownedServiceAsOf504(t, false)
	before = strings.Replace(before,
		`s.db(ctx).Model(&models.Memo{}).Where("id IN ?", ids).Find(&allowed)`,
		`s.rowsFor(ctx, ids).Find(&allowed)`, 1)

	after, note := repairGridOwnerSource(before, "Memo", "memo")
	if note == "" {
		t.Fatal("a hand-written BulkEdit was reported as fine")
	}
	if !strings.Contains(note, "BulkEdit") || !strings.Contains(note, "by hand") {
		t.Errorf("the note does not say what to check: %q", note)
	}
	if after != before {
		t.Error("the repair edited a service it said it could not read")
	}
}

// End to end over a directory, which is what grit upgrade calls.
func TestRepairGridOwnerScopingWalksTheServices(t *testing.T) {
	apiRoot := t.TempDir()
	dir := filepath.Join(apiRoot, "internal", "services")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "memo.go")
	if err := os.WriteFile(path, []byte(ownedServiceAsOf504(t, false)), 0o644); err != nil {
		t.Fatal(err)
	}

	changed, blocked, err := RepairGridOwnerScoping(apiRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocked) > 0 {
		t.Fatalf("blocked: %v", blocked)
	}
	if len(changed) != 1 {
		t.Fatalf("changed %d services, want 1: %v", len(changed), changed)
	}
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "rows = authz.ScopeOwned(") {
		t.Error("the file on disk was not scoped")
	}

	// And nothing on the second run.
	changed, _, err = RepairGridOwnerScoping(apiRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 0 {
		t.Errorf("a second upgrade reported %d changes: %v", len(changed), changed)
	}
}

// A missing services directory is a project shape, not a failure.
func TestRepairGridOwnerScopingIgnoresAProjectWithNoServices(t *testing.T) {
	changed, blocked, err := RepairGridOwnerScoping(t.TempDir())
	if err != nil || len(changed) > 0 || len(blocked) > 0 {
		t.Errorf("changed=%v blocked=%v err=%v", changed, blocked, err)
	}
}

// ownerScopeOf is what decides the column, so it has to read both shapes and
// nothing else.
func TestOwnerScopeOfReadsBothScopes(t *testing.T) {
	cases := []struct {
		name, src, col, perm string
	}{
		{"owned", `query = authz.ScopeOwned(ctx, query, "user_id")`, "user_id", ""},
		{"confidential", `query = authz.ScopeConfidential(ctx, query, "patient_id", "records.confidential")`, "patient_id", "records.confidential"},
		{"shared", `query := s.db(ctx).Model(&models.Memo{})`, "", ""},
		// The bulk and grid calls scope a different variable. Reading one of
		// those would be reading the repair's own output.
		{"not the bulk call", `scope = authz.ScopeOwned(ctx, scope, "user_id")`, "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col, perm := ownerScopeOf(c.src)
			if col != c.col || perm != c.perm {
				t.Errorf("got (%q, %q), want (%q, %q)", col, perm, c.col, c.perm)
			}
		})
	}
}

func TestPascalFromSnake(t *testing.T) {
	for in, want := range map[string]string{
		"user_id":          "UserID",
		"patient_id":       "PatientID",
		"owner_account_id": "OwnerAccountID",
		"author":           "Author",
	} {
		if got := pascalFromSnake(in); got != want {
			t.Errorf("pascalFromSnake(%q) = %q, want %q", in, got, want)
		}
	}
}
