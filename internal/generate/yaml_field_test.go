package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// A YAML resource takes the same field syntax the command line does.
func TestYAMLAcceptsBothFieldForms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "project.yaml")
	body := `name: Project
fields:
  - title:string:unique
  - kind:select[one_off, recurring]
  - cycle:select[monthly, annually] if kind == recurring
  - name: notes
    type: richtext
    required: false
  - name: plan
    type: belongs_to
    related_model: Plan
    if: kind == recurring
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	def, err := LoadFromYAML(path)
	if err != nil {
		t.Fatalf("LoadFromYAML: %v", err)
	}
	if len(def.Fields) != 5 {
		t.Fatalf("got %d fields", len(def.Fields))
	}
	if !def.Fields[0].Unique {
		t.Error("the one-line spec lost its modifier")
	}
	if len(def.Fields[1].Options) != 2 || def.Fields[1].Options[0].Label != "One Off" {
		t.Errorf("options = %+v", def.Fields[1].Options)
	}
	cycle := def.Fields[2]
	if cycle.ShowWhen == nil || cycle.ShowWhen.Field != "kind" || cycle.ShowWhen.Value != "recurring" {
		t.Errorf("one-line condition = %+v", cycle.ShowWhen)
	}
	if def.Fields[3].Type != "richtext" || def.Fields[3].Required {
		t.Errorf("the mapping form broke: %+v", def.Fields[3])
	}
	plan := def.Fields[4]
	if plan.ShowWhen == nil || plan.ShowWhen.Field != "kind" {
		t.Errorf("the if: key did not reach the field: %+v", plan.ShowWhen)
	}
	if plan.RelatedModel != "Plan" {
		t.Errorf("related model = %q", plan.RelatedModel)
	}
}
