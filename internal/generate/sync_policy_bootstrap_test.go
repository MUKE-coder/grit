package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// routesGoWithMarker is the generated file, with the registrations the
// generator injects.
const routesGoWithMarker = `package routes

func RegisterSyncModels(syncRegistry *sync.Registry) {
	syncRegistry.Register("products", &models.Product{})
	syncRegistry.Register("orders", &models.Order{})
	// grit:sync
}
`

// policiesInTheirOwnFile is the pattern grit#214 uses and recommends: a
// hand-written file applied after the // grit:sync marker, so a regeneration
// cannot quietly put a table back to DefaultPolicy.
//
// It is also the pattern a reader that only opened routes.go would miss
// entirely, writing an empty bootstrap and looking like it worked.
const policiesInTheirOwnFile = `package routes

// Applied after the // grit:sync block.
func registerSyncPolicies(syncRegistry *sync.Registry) {
	syncRegistry.RegisterWithPolicy("sentence_origins", &models.SentenceOrigin{}, sync.Policy{
		Mode:     sync.ModeReadOnly,
		Conflict: sync.ConflictServerWins,
	})
	syncRegistry.RegisterWithPolicy("stock_levels", &models.StockLevel{}, sync.Policy{
		Mode:          sync.ModeOnlineOnly,
		MaxOfflineAge: 72 * time.Hour,
	})
	syncRegistry.RegisterWithPolicy("sales", &models.Sale{}, sync.Policy{
		Mode:     sync.ModeOfflineFirst,
		Conflict: sync.ConflictClientWins,
	})
}
`

func writePolicyProject(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	routesDir := filepath.Join(root, "apps", "api", "internal", "routes")
	if err := os.MkdirAll(routesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(routesDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "apps", "desktop", "sync"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

// The reader has to find a policy declared in a file of its own, because that
// is where the project that reported this declares them.
func TestPoliciesAreReadFromEveryRoutesFile(t *testing.T) {
	root := writePolicyProject(t, map[string]string{
		"routes.go":             routesGoWithMarker,
		"sync_policies.go":      policiesInTheirOwnFile,
		"sync_policies_test.go": "package routes // must be ignored",
	})

	declared, err := ReadDeclaredSyncPolicies(root)
	if err != nil {
		t.Fatal(err)
	}
	byTable := map[string]scaffold.SyncPolicyDeclaration{}
	for _, d := range declared {
		byTable[d.Table] = d
	}

	for table, want := range map[string]string{
		"sentence_origins": "read_only",
		"stock_levels":     "online_only",
		"sales":            "offline_first",
		"products":         "offline_first",
		"orders":           "offline_first",
	} {
		got, ok := byTable[table]
		if !ok {
			t.Errorf("%s was not found; a policy in a file of its own is still a policy", table)
			continue
		}
		if got.Mode != want {
			t.Errorf("%s: mode %q, want %q", table, got.Mode, want)
		}
	}
	if got := byTable["sentence_origins"].Conflict; got != "server_wins" {
		t.Errorf("sentence_origins conflict is %q, want server_wins", got)
	}
	if got := byTable["stock_levels"].MaxOfflineAgeSeconds; got != 72*3600 {
		t.Errorf("stock_levels max offline age is %d seconds, want %d", got, 72*3600)
	}
}

// The rendered file has to compile, hold the policies, and be stable: an
// unsorted map would produce a diff on every run.
func TestTheBootstrapIsSortedAndComplete(t *testing.T) {
	declared := []scaffold.SyncPolicyDeclaration{
		{Table: "stock_levels", Mode: "online_only"},
		{Table: "audit_entries", Mode: "read_only", Conflict: "server_wins"},
		{Table: "sales", Mode: "offline_first", Conflict: "manual"},
	}
	first := scaffold.RenderSyncPolicyBootstrap(declared)

	for _, want := range []string{
		`"audit_entries": {Mode: "read_only", Conflict: "server_wins"},`,
		`"sales": {Mode: "offline_first", Conflict: "manual"},`,
		`"stock_levels": {Mode: "online_only"},`,
	} {
		if !strings.Contains(first, want) {
			t.Errorf("the bootstrap is missing %s:\n%s", want, first)
		}
	}
	// Sorted: audit before sales before stock.
	if a, s := strings.Index(first, "audit_entries"), strings.Index(first, "sales"); a > s {
		t.Error("the entries are not sorted by table, so every regeneration is a diff")
	}
	// And the marker survives, or a second run has nowhere to write.
	if !strings.Contains(first, scaffold.SyncPolicyBootstrapMarker) {
		t.Error("the marker was consumed, so the file can only be generated once")
	}

	shuffled := []scaffold.SyncPolicyDeclaration{declared[2], declared[0], declared[1]}
	if second := scaffold.RenderSyncPolicyBootstrap(shuffled); second != first {
		t.Error("the same policies in a different order rendered differently")
	}
}

// End to end: a project's declarations reach the desktop app's file.
func TestSyncWritesTheBootstrapIntoTheDesktopApp(t *testing.T) {
	root := writePolicyProject(t, map[string]string{
		"routes.go":        routesGoWithMarker,
		"sync_policies.go": policiesInTheirOwnFile,
	})
	path := filepath.Join(root, "apps", "desktop", "sync", "policies_bootstrap.go")

	line, err := SyncDesktopPolicyBootstrap(root)
	if err != nil {
		t.Fatal(err)
	}
	if line == "" {
		t.Fatal("nothing was reported for a project that declares five tables")
	}
	if !strings.Contains(line, "2 of them refuse") {
		t.Errorf("the report does not say how many tables are guarded: %q", line)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the bootstrap was not written: %v", err)
	}
	if !strings.Contains(string(body), `"sentence_origins": {Mode: "read_only"`) {
		t.Errorf("the read_only table is not in the written file:\n%s", body)
	}

	// Second run, nothing changed: no write and nothing to report, so an
	// unchanged project does not show a line claiming work.
	line, err = SyncDesktopPolicyBootstrap(root)
	if err != nil {
		t.Fatal(err)
	}
	if line != "" {
		t.Errorf("a second run reported %q with nothing changed", line)
	}
}

// A project with no desktop app must be silent rather than failing or writing
// a stray directory.
func TestNoDesktopAppMeansNoBootstrap(t *testing.T) {
	root := t.TempDir()
	routesDir := filepath.Join(root, "apps", "api", "internal", "routes")
	if err := os.MkdirAll(routesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(routesDir, "routes.go"), []byte(routesGoWithMarker), 0o644); err != nil {
		t.Fatal(err)
	}

	line, err := SyncDesktopPolicyBootstrap(root)
	if err != nil {
		t.Errorf("a project with no desktop app returned %v", err)
	}
	if line != "" {
		t.Errorf("reported %q for a project with no desktop app", line)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "desktop")); err == nil {
		t.Error("a desktop directory was created in a project that has none")
	}
}

// Two registrations for one table: take the stricter. A guard that is wrongly
// strict shows a refusal; one that is wrongly lax queues a doomed write.
func TestADuplicateRegistrationTakesTheStricterPolicy(t *testing.T) {
	root := writePolicyProject(t, map[string]string{
		"a_routes.go": `package routes
func a() {
	syncRegistry.Register("ledger", &models.Ledger{})
}
`,
		"z_routes.go": `package routes
func z() {
	syncRegistry.RegisterWithPolicy("ledger", &models.Ledger{}, sync.Policy{
		Mode: sync.ModeReadOnly,
	})
}
`,
	})

	declared, err := ReadDeclaredSyncPolicies(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(declared) != 1 {
		t.Fatalf("read %d declarations, want 1: %+v", len(declared), declared)
	}
	if declared[0].Mode != "read_only" {
		t.Errorf("mode is %q, want read_only: the stricter of the two", declared[0].Mode)
	}
}

// An unrecognised mode constant passes through rather than being guessed at,
// and the client treats anything but offline_first as unwritable, so the
// unknown case errs toward refusing.
func TestAnUnknownModeConstantIsNotGuessed(t *testing.T) {
	root := writePolicyProject(t, map[string]string{
		"routes.go": `package routes
func r() {
	syncRegistry.RegisterWithPolicy("mystery", &models.Mystery{}, sync.Policy{
		Mode: sync.ModeSomethingNew,
	})
}
`,
	})

	declared, err := ReadDeclaredSyncPolicies(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(declared) != 1 || declared[0].Mode != "ModeSomethingNew" {
		t.Errorf("read %+v, want the constant name passed through unchanged", declared)
	}
}

// The spelling the generator itself emits, which is the one most likely to be
// in a real project: raw nanoseconds with the duration in a comment. Reading
// only `72h` returned zero here, and zero means "no limit", so a declared
// staleness limit silently stopped existing.
func TestMaxOfflineAgeReadsEverySpelling(t *testing.T) {
	cases := map[string]int{
		"MaxOfflineAge: 259200000000000 /* 72h */": 72 * 3600,
		"MaxOfflineAge: 72 * time.Hour":            72 * 3600,
		"MaxOfflineAge: 72*time.Hour":              72 * 3600,
		"MaxOfflineAge: time.Hour":                 3600,
		"MaxOfflineAge: 30 * time.Minute":          1800,
		"MaxOfflineAge: 45 * time.Second":          45,
		"Conflict: sync.ConflictManual":            0,
	}
	for literal, want := range cases {
		if got := readMaxOfflineAge(literal); got != want {
			t.Errorf("%q read as %d seconds, want %d", literal, got, want)
		}
	}
}

// The tie case, which the read_only duplicate test above does not reach and a
// real project found immediately: a table registered bare by the generator in
// routes.go and given a policy by hand in another file. Both are
// offline_first, so strictness cannot decide, and the explicit one has to win
// or the project's Conflict and MaxOfflineAge are silently replaced by
// defaults.
func TestAnExplicitPolicyBeatsABareRegistrationOfTheSameMode(t *testing.T) {
	root := writePolicyProject(t, map[string]string{
		// Sorts first, so before the fix this one won.
		"routes.go": `package routes
func r() {
	syncRegistry.Register("blogs", &models.Blog{})
	// grit:sync
}
`,
		"sync_policies.go": `package routes
func p() {
	syncRegistry.RegisterWithPolicy("blogs", &models.Blog{}, sync.Policy{
		Mode:          sync.ModeOfflineFirst,
		Conflict:      sync.ConflictClientWins,
		MaxOfflineAge: 72 * time.Hour,
	})
}
`,
	})

	declared, err := ReadDeclaredSyncPolicies(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(declared) != 1 {
		t.Fatalf("read %d declarations, want 1: %+v", len(declared), declared)
	}
	got := declared[0]
	if got.Conflict != "client_wins" {
		t.Errorf("conflict is %q, want client_wins: the bare registration's default replaced what the project wrote", got.Conflict)
	}
	if got.MaxOfflineAgeSeconds != 72*3600 {
		t.Errorf("max offline age is %d, want %d", got.MaxOfflineAgeSeconds, 72*3600)
	}
}
