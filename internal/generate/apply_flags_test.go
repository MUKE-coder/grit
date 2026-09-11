package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// owned_by: user in a --from file generated a resource every signed-in account
// could read and edit, because the unset --owned-by flag was assigned over it.
// Found generating a patient-referral resource from YAML: the handler had no
// owner scoping at all and demanded user_id in the request body.
func TestFlagsDoNotEraseTheDefinitionFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "note.yaml")
	body := `name: Note
owned_by: user
tree: true
public: true
append_only: true
tenant_owned: true
audit_reads: true
fields:
  - name: body
    type: text
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	def, err := LoadFromYAML(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	// Every flag at its default, as when only --from is passed.
	def.ApplyFlags(false, false, "", false, false, false)
	if def.OwnedBy != "user" || !def.Tree || !def.Public || !def.AppendOnly || !def.TenantOwned || !def.AuditReads {
		t.Errorf("a default flag erased the file's value: %+v", def)
	}

	// A flag still adds to a definition that did not say.
	plain := &ResourceDefinition{Name: "Post"}
	plain.ApplyFlags(true, true, "author", true, true, true)
	if plain.OwnedBy != "author" || !plain.Tree || !plain.Public || !plain.AppendOnly || !plain.TenantOwned || !plain.AuditReads {
		t.Errorf("a flag did not apply: %+v", plain)
	}
}
