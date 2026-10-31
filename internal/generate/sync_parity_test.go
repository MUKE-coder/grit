package generate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Two writers produce the Zod schema for a resource, and nothing compared
// them.
//
// `grit generate resource` writes it from the definition; `grit sync` rewrites
// the same file from the Go model, and sync is the documented step after
// generating, so sync's file is the one projects keep. It was the weaker of the
// two in nine different ways at once, and a create form built from it could not
// be submitted (grit#108).
//
// This is the comparison. It generates a resource carrying every field type,
// then asks both writers for the create schema and fails on any field where
// they disagree. The point is not the nine bugs; it is that the tenth cannot
// be added to one writer alone.
func TestSyncAndGenerateAgreeOnEveryFieldType(t *testing.T) {
	root := t.TempDir()
	def := everyFieldTypeDefinition()
	g := newTestGenerator(root, "testapp/apps/api", def)
	names := BuildNames(def)

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}
	// The model first: the schema is written from it rather than from the
	// definition, which is the order Generator.Run uses.
	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}
	if err := g.writeZodSchema(names); err != nil {
		t.Fatalf("writeZodSchema: %v", err)
	}

	schemaPath := filepath.Join(g.SharedRoot(), "schemas", names.Kebab+".ts")
	generated, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("reading the generated schema: %v", err)
	}

	modelPath := filepath.Join(root, "apps", "api", "internal", "models", names.Snake+".go")
	structs, err := parseGoStructs(modelPath)
	if err != nil {
		t.Fatalf("parsing the generated model: %v", err)
	}
	var synced string
	for _, s := range structs {
		if s.Name == names.Pascal {
			synced = buildZodSchema(s)
		}
	}
	if synced == "" {
		t.Fatalf("no %s struct in the generated model", names.Pascal)
	}

	want := zodObjectFields(t, string(generated), "Create"+names.Pascal+"Schema")
	got := zodObjectFields(t, synced, "Create"+names.Pascal+"Schema")

	for key, wantZod := range want {
		gotZod, present := got[key]
		if !present {
			t.Errorf("%s: the generator validates it, sync drops it", key)
			continue
		}
		if gotZod != wantZod {
			t.Errorf("%s:\n  generate: %s\n  sync:     %s", key, wantZod, gotZod)
		}
	}
	for key := range got {
		if _, present := want[key]; !present {
			t.Errorf("%s: sync asks a create form for a field the generator knows "+
				"the server fills", key)
		}
	}

	// A server-set column is in neither schema, which a key-by-key comparison
	// of the two writers cannot tell from both of them being wrong.
	if _, present := got["closed_at"]; present {
		t.Error("closed_at is server-set and should not be asked of a caller")
	}

	// The conditional field's rule is a refinement, not a type, so a field-by
	// -field comparison cannot see it.
	if strings.Contains(string(generated), "superRefine") && !strings.Contains(synced, "superRefine") {
		t.Error("the generator enforces a conditional requirement that sync drops")
	}

	// And a schema that names MoneySchema has to import it, whichever writer
	// named it.
	for _, imp := range []string{"./money", "./file-ref", "./field-formats", "./phone"} {
		if strings.Contains(string(generated), imp) && !strings.Contains(synced, imp) {
			t.Errorf("sync's schema does not import %s", imp)
		}
	}
}

// The API drops it from every write, so neither writer may ask for it.
func TestVersionIsNeverOnACreateSchema(t *testing.T) {
	root := t.TempDir()
	def, err := ParseInlineFields("Note", "title:string")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, "testapp/apps/api", def)
	names := BuildNames(def)
	if err := g.writeGoModel(names); err != nil {
		t.Fatal(err)
	}
	structs, err := parseGoStructs(filepath.Join(root, "apps", "api", "internal", "models", "note.go"))
	if err != nil {
		t.Fatal(err)
	}
	synced := buildZodSchema(structs[0])

	create := zodObjectFields(t, synced, "CreateNoteSchema")
	if _, present := create["version"]; present {
		t.Error("version is on the create schema: it travels as an ETag, the " +
			"service drops it from the body, and a form cannot invent one")
	}
	// It does come back, though, so the response schema keeps it.
	entity := zodObjectFields(t, synced, "NoteSchema")
	if _, present := entity["version"]; !present {
		t.Error("version is missing from the entity schema, which is what you " +
			"validate a response against")
	}
}

var zodFieldRe = regexp.MustCompile(`^\s{2}([a-z_0-9]+):\s*(.+),$`)

// zodObjectFields reads one z.object({...}) out of a schema file, as
// key -> validator.
func zodObjectFields(t *testing.T, src, name string) map[string]string {
	t.Helper()
	open := "export const " + name + " = z.object({"
	i := strings.Index(src, open)
	if i < 0 {
		t.Fatalf("no %s in:\n%s", name, src)
	}
	rest := src[i+len(open):]
	end := strings.Index(rest, "\n})")
	if end < 0 {
		t.Fatalf("%s is never closed", name)
	}
	out := map[string]string{}
	for _, line := range strings.Split(rest[:end], "\n") {
		if m := zodFieldRe.FindStringSubmatch(line); m != nil {
			out[m[1]] = m[2]
		}
	}
	return out
}

// everyFieldTypeDefinition is one resource carrying every field type, required
// and optional, including the three the generator leaves out of a create
// schema: a slug, an owner, and a field required only while another field
// holds a value.
func everyFieldTypeDefinition() *ResourceDefinition {
	opt := func(values ...string) []FieldOption {
		out := make([]FieldOption, len(values))
		for i, v := range values {
			out[i] = FieldOption{Value: v, Label: strings.ToUpper(v)}
		}
		return out
	}
	return &ResourceDefinition{
		Name:    "Sweep",
		OwnedBy: "owner",
		Fields: []Field{
			{Name: "req_string", Type: "string", Required: true},
			{Name: "opt_string", Type: "string"},
			{Name: "body", Type: "text"},
			{Name: "article", Type: "richtext"},
			{Name: "handle", Type: "slug", SlugSource: "req_string"},
			{Name: "count", Type: "int", Required: true},
			{Name: "qty", Type: "uint"},
			{Name: "ratio", Type: "float"},
			{Name: "active", Type: "bool"},
			{Name: "flag", Type: "toggle"},
			{Name: "starts_at", Type: "datetime"},
			{Name: "due_on", Type: "date"},
			{Name: "price", Type: "money", Required: true},
			{Name: "tags_list", Type: "string_array"},
			{Name: "cover", Type: "file"},
			{Name: "gallery", Type: "files"},
			{Name: "stage", Type: "select", Required: true, Options: opt("a", "b")},
			// A default makes a required field optional in create, because the
			// server fills it in. The two writers disagreed about this and the
			// fixture had no defaulted field to notice with (grit#114).
			{Name: "mode", Type: "select", Required: true, Default: "draft",
				Options: opt("draft", "live")},
			// And a column the server fills belongs in neither schema.
			{Name: "closed_at", Type: "datetime", ServerSet: true},
			{Name: "pick", Type: "radio", Options: opt("x")},
			{Name: "perms", Type: "check", Options: opt("read", "write")},
			{Name: "email_addr", Type: "email", Required: true},
			{Name: "site", Type: "url"},
			{Name: "host", Type: "domain"},
			{Name: "phone", Type: "tel"},
			{Name: "nation", Type: "country"},
			{Name: "shade", Type: "color"},
			{Name: "done_pct", Type: "percent"},
			{Name: "score", Type: "rating", Max: 5},
			{Name: "at_time", Type: "time"},
			{Name: "payload", Type: "json"},
			{Name: "pipeline", Type: "belongs_to", RelatedModel: "Pipeline", Required: true},
			{Name: "contact", Type: "belongs_to", RelatedModel: "Contact"},
			{Name: "owner", Type: "belongs_to", RelatedModel: "User", Required: true},
			{Name: "lost_reason", Type: "text", Required: true,
				ShowWhen: &FieldCondition{Field: "stage", Value: "b"}},
		},
	}
}
