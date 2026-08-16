package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncPolicyValidateFillsDefaults(t *testing.T) {
	p := &SyncPolicy{}
	if err := p.Validate("Sale"); err != nil {
		t.Fatalf("an empty policy should be valid: %v", err)
	}
	if p.Mode != "offline_first" || p.Conflict != "manual" {
		t.Fatalf("defaults not applied: %+v", p)
	}
}

func TestSyncPolicyValidateRejectsTypos(t *testing.T) {
	for _, tc := range []struct {
		name   string
		policy SyncPolicy
		want   string
	}{
		{"bad mode", SyncPolicy{Mode: "offline"}, "sync.mode"},
		{"bad conflict", SyncPolicy{Conflict: "server-wins"}, "sync.conflict"},
		{"bad duration", SyncPolicy{MaxOfflineAge: "72 hours"}, "max_offline_age"},
		{"negative duration", SyncPolicy{MaxOfflineAge: "-1h"}, "must be positive"},
		{"version local_only", SyncPolicy{LocalOnly: []string{"version"}}, "cannot include"},
		{"id local_only", SyncPolicy{LocalOnly: []string{"id"}}, "cannot include"},
		{
			"same field in both lists",
			SyncPolicy{Fields: []string{"total"}, LocalOnly: []string{"total"}},
			"both sync.fields and sync.local_only",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.policy
			err := p.Validate("Sale")
			if err == nil {
				t.Fatalf("expected an error for %+v", tc.policy)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error %q does not mention %q", err, tc.want)
			}
			if !strings.Contains(err.Error(), "Sale") {
				t.Errorf("error %q does not name the resource", err)
			}
		})
	}
}

// The silent one. sync.fields naming a column that does not exist errors
// nowhere: the allowlist simply excludes the real column, and every client
// mirrors rows with the value missing.
func TestSyncPolicyCatchesFieldsThatDoNotExist(t *testing.T) {
	fields := []Field{
		{Name: "total", Type: "float"},
		{Name: "payment_method", Type: "string"},
	}
	p := &SyncPolicy{
		Fields:    []string{"totl", "payment_method"},
		LocalOnly: []string{"draft_not"},
	}

	unknown := p.CheckAgainstFields("Sale", fields)
	if len(unknown) != 2 {
		t.Fatalf("got %v, want the two misspelled names", unknown)
	}
	if unknown[0] != "totl" || unknown[1] != "draft_not" {
		t.Fatalf("got %v", unknown)
	}

	// The bookkeeping columns are always available even though no Field
	// declares them, so naming one must not be reported as a typo.
	ok := &SyncPolicy{Fields: []string{"total", "id", "version", "updated_at"}}
	if got := ok.CheckAgainstFields("Sale", fields); len(got) != 0 {
		t.Fatalf("bookkeeping columns reported as unknown: %v", got)
	}
}

func TestSyncPolicyGoLiteral(t *testing.T) {
	// A resource with no real policy keeps the plain Register call it always
	// had, so a project without policies produces identical routes.go.
	if got := (&SyncPolicy{Mode: "offline_first", Conflict: "manual"}).GoLiteral(); got != "" {
		t.Errorf("an all-defaults policy should emit nothing, got %q", got)
	}
	var nilPolicy *SyncPolicy
	if got := nilPolicy.GoLiteral(); got != "" {
		t.Errorf("a nil policy should emit nothing, got %q", got)
	}

	p := &SyncPolicy{
		Mode:          "offline_first",
		Conflict:      "server_wins",
		Fields:        []string{"total", "payment_method"},
		LocalOnly:     []string{"draft_note"},
		MaxOfflineAge: "72h",
	}
	got := p.GoLiteral()
	for _, want := range []string{
		"sync.Policy{",
		"Mode: sync.ModeOfflineFirst",
		"Conflict: sync.ConflictServerWins",
		`Fields: []string{"total", "payment_method"}`,
		`LocalOnly: []string{"draft_note"}`,
		"MaxOfflineAge: 259200000000000",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("literal missing %q:\n%s", want, got)
		}
	}
	// The duration is emitted as nanoseconds with the source written beside
	// it, so a reader of routes.go does not have to divide by a billion.
	if !strings.Contains(got, "/* 72h */") {
		t.Errorf("literal should keep the human duration as a comment:\n%s", got)
	}
}

func TestLoadFromYAMLParsesAndValidatesSyncPolicy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sale.yaml")
	body := `name: Sale
fields:
  - name: total
    type: float
  - name: payment_method
    type: string
  - name: draft_note
    type: text
sync:
  mode: offline_first
  conflict: server_wins
  fields: [total, payment_method]
  local_only: [draft_note]
  max_offline_age: 72h
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	def, err := LoadFromYAML(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if def.Sync == nil {
		t.Fatal("the sync block was not parsed")
	}
	if def.Sync.Conflict != "server_wins" {
		t.Errorf("conflict = %q", def.Sync.Conflict)
	}
	if len(def.Sync.Fields) != 2 || len(def.Sync.LocalOnly) != 1 {
		t.Errorf("field lists wrong: %+v", def.Sync)
	}
}

func TestLoadFromYAMLRejectsAPolicyNamingAMissingField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sale.yaml")
	body := `name: Sale
fields:
  - name: total
    type: float
sync:
  fields: [totl]
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadFromYAML(path)
	if err == nil {
		t.Fatal("a policy naming a field the resource does not have should fail to load")
	}
	if !strings.Contains(err.Error(), "totl") {
		t.Errorf("the error should name the offending field, got %q", err)
	}
}

func TestLoadFromYAMLWithoutASyncBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "post.yaml")
	if err := os.WriteFile(path, []byte("name: Post\nfields:\n  - name: title\n    type: string\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	def, err := LoadFromYAML(path)
	if err != nil {
		t.Fatalf("a definition with no sync block must still load: %v", err)
	}
	if def.Sync != nil {
		t.Fatal("no sync block should stay nil, so the generator emits the plain Register call")
	}
}
