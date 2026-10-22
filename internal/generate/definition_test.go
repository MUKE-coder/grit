package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// ── ParseInlineFields ────────────────────────────────────────────────────────

func TestParseInlineFields(t *testing.T) {
	t.Run("basic fields", func(t *testing.T) {
		def, err := ParseInlineFields("Post", "title:string,content:text,published:bool")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Name != "Post" {
			t.Errorf("Name = %q, want %q", def.Name, "Post")
		}
		if len(def.Fields) != 3 {
			t.Fatalf("len(Fields) = %d, want 3", len(def.Fields))
		}
		if def.Fields[0].Name != "title" || def.Fields[0].Type != "string" {
			t.Errorf("Fields[0] = %+v, want {title string}", def.Fields[0])
		}
		if def.Fields[1].Name != "content" || def.Fields[1].Type != "text" {
			t.Errorf("Fields[1] = %+v, want {content text}", def.Fields[1])
		}
		if def.Fields[2].Name != "published" || def.Fields[2].Type != "bool" {
			t.Errorf("Fields[2] = %+v, want {published bool}", def.Fields[2])
		}
	})

	t.Run("slug field infers source", func(t *testing.T) {
		def, err := ParseInlineFields("Post", "title:string,slug:slug:title")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		slugField := def.Fields[1]
		if slugField.Type != "slug" {
			t.Errorf("expected slug type, got %q", slugField.Type)
		}
		if slugField.SlugSource != "title" {
			t.Errorf("SlugSource = %q, want %q", slugField.SlugSource, "title")
		}
		if !slugField.Unique {
			t.Error("slug field should be Unique")
		}
	})

	t.Run("belongs_to without explicit model", func(t *testing.T) {
		def, err := ParseInlineFields("Post", "category:belongs_to")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		f := def.Fields[0]
		if f.Type != "belongs_to" {
			t.Errorf("Type = %q, want belongs_to", f.Type)
		}
		if !f.Required {
			t.Error("belongs_to field should be Required")
		}
	})

	t.Run("belongs_to with explicit model", func(t *testing.T) {
		def, err := ParseInlineFields("Post", "author:belongs_to:User")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		f := def.Fields[0]
		if f.RelatedModel != "User" {
			t.Errorf("RelatedModel = %q, want %q", f.RelatedModel, "User")
		}
	})

	t.Run("many_to_many", func(t *testing.T) {
		def, err := ParseInlineFields("Post", "tags:many_to_many:Tag")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		f := def.Fields[0]
		if f.Type != "many_to_many" {
			t.Errorf("Type = %q, want many_to_many", f.Type)
		}
		if f.RelatedModel != "Tag" {
			t.Errorf("RelatedModel = %q, want Tag", f.RelatedModel)
		}
	})

	t.Run("many_to_many without model fails", func(t *testing.T) {
		_, err := ParseInlineFields("Post", "tags:many_to_many")
		if err == nil {
			t.Error("expected error for many_to_many without model, got nil")
		}
	})

	t.Run("unique modifier", func(t *testing.T) {
		def, err := ParseInlineFields("User", "email:string:unique")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !def.Fields[0].Unique {
			t.Error("field should be Unique")
		}
	})

	t.Run("optional modifier overrides string default required", func(t *testing.T) {
		def, err := ParseInlineFields("Post", "subtitle:string:optional")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Fields[0].Required {
			t.Error("field should not be Required when optional modifier is set")
		}
	})

	t.Run("empty fields error", func(t *testing.T) {
		_, err := ParseInlineFields("Post", "")
		if err == nil {
			t.Error("expected error for empty fields, got nil")
		}
	})

	t.Run("invalid type error", func(t *testing.T) {
		_, err := ParseInlineFields("Post", "title:notavalidtype")
		if err == nil {
			t.Error("expected error for invalid type, got nil")
		}
	})

	t.Run("invalid modifier error", func(t *testing.T) {
		_, err := ParseInlineFields("Post", "title:string:badmodifier")
		if err == nil {
			t.Error("expected error for invalid modifier, got nil")
		}
	})

	t.Run("missing name error", func(t *testing.T) {
		_, err := ParseInlineFields("Post", ":string")
		if err == nil {
			t.Error("expected error for empty field name, got nil")
		}
	})

	// v3.31.30 — file/files types
	t.Run("single file with image alias", func(t *testing.T) {
		def, err := ParseInlineFields("Product", "image:file:image")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(def.Fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(def.Fields))
		}
		f := def.Fields[0]
		if f.Type != "file" {
			t.Errorf("expected type file, got %q", f.Type)
		}
		if len(f.FileAccepts) != 1 || f.FileAccepts[0] != "image" {
			t.Errorf("expected FileAccepts=[image], got %v", f.FileAccepts)
		}
	})

	t.Run("single file with all alias", func(t *testing.T) {
		def, err := ParseInlineFields("Product", "attachment:file:all")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Fields[0].FileAccepts[0] != "all" {
			t.Errorf("expected accepts=[all], got %v", def.Fields[0].FileAccepts)
		}
	})

	t.Run("file with default accept list (omitted)", func(t *testing.T) {
		def, err := ParseInlineFields("Product", "doc:file")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Fields[0].FileAccepts[0] != "all" {
			t.Errorf("default accepts should be [all], got %v", def.Fields[0].FileAccepts)
		}
	})

	t.Run("file with bracketed multi-accept list", func(t *testing.T) {
		def, err := ParseInlineFields("Product", "attachment:file:[pdf,doc,image,video,zip]")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := def.Fields[0].FileAccepts
		want := []string{"pdf", "doc", "image", "video", "zip"}
		if len(got) != len(want) {
			t.Fatalf("expected %d accepts, got %d (%v)", len(want), len(got), got)
		}
		for i, v := range want {
			if got[i] != v {
				t.Errorf("accepts[%d]: want %q, got %q", i, v, got[i])
			}
		}
	})

	t.Run("files (plural) with bracket list", func(t *testing.T) {
		def, err := ParseInlineFields("Product", "gallery:files:[image,video]")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Fields[0].Type != "files" {
			t.Errorf("expected files (plural), got %q", def.Fields[0].Type)
		}
	})

	t.Run("file mixed with other fields", func(t *testing.T) {
		def, err := ParseInlineFields("Product", "name:string,image:file:image,price:float")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(def.Fields) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(def.Fields))
		}
	})

	t.Run("file with bracket list does not split fields incorrectly", func(t *testing.T) {
		// Bracket-awareness check: the inner commas in [pdf,doc] must
		// not break the top-level field split.
		def, err := ParseInlineFields("Product", "name:string,attachments:files:[pdf,doc,image],price:float")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(def.Fields) != 3 {
			t.Fatalf("expected 3 fields (not %d) — bracket-aware split failed", len(def.Fields))
		}
		if len(def.Fields[1].FileAccepts) != 3 {
			t.Errorf("expected 3 accepts, got %d", len(def.Fields[1].FileAccepts))
		}
	})

	t.Run("invalid file accept alias", func(t *testing.T) {
		_, err := ParseInlineFields("Product", "x:file:wat")
		if err == nil {
			t.Error("expected error for invalid accept alias")
		}
	})

	t.Run("empty bracket list", func(t *testing.T) {
		_, err := ParseInlineFields("Product", "x:file:[]")
		if err == nil {
			t.Error("expected error for empty bracket accept-list")
		}
	})
}

// ── LoadFromYAML ─────────────────────────────────────────────────────────────

func TestLoadFromYAML(t *testing.T) {
	t.Run("valid YAML", func(t *testing.T) {
		content := `name: Product
fields:
  - name: title
    type: string
    required: true
  - name: price
    type: float
  - name: active
    type: bool
`
		f := writeTempYAML(t, content)
		def, err := LoadFromYAML(f)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Name != "Product" {
			t.Errorf("Name = %q, want Product", def.Name)
		}
		if len(def.Fields) != 3 {
			t.Fatalf("len(Fields) = %d, want 3", len(def.Fields))
		}
	})

	t.Run("missing name error", func(t *testing.T) {
		content := `fields:
  - name: title
    type: string
`
		f := writeTempYAML(t, content)
		_, err := LoadFromYAML(f)
		if err == nil {
			t.Error("expected error for missing name, got nil")
		}
	})

	t.Run("no fields error", func(t *testing.T) {
		content := `name: Product
fields: []
`
		f := writeTempYAML(t, content)
		_, err := LoadFromYAML(f)
		if err == nil {
			t.Error("expected error for empty fields, got nil")
		}
	})

	t.Run("invalid field type error", func(t *testing.T) {
		content := `name: Product
fields:
  - name: price
    type: currency
`
		f := writeTempYAML(t, content)
		_, err := LoadFromYAML(f)
		if err == nil {
			t.Error("expected error for invalid field type, got nil")
		}
	})

	t.Run("missing field name error", func(t *testing.T) {
		content := `name: Product
fields:
  - type: string
`
		f := writeTempYAML(t, content)
		_, err := LoadFromYAML(f)
		if err == nil {
			t.Error("expected error for missing field name, got nil")
		}
	})

	t.Run("file not found error", func(t *testing.T) {
		_, err := LoadFromYAML("/tmp/does-not-exist-grit-test.yaml")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})
}

// ── ValidFieldTypes ───────────────────────────────────────────────────────────

func TestValidFieldTypes(t *testing.T) {
	types := ValidFieldTypes()
	if len(types) == 0 {
		t.Fatal("ValidFieldTypes() returned empty slice")
	}

	required := []string{"string", "text", "int", "uint", "float", "bool", "datetime", "date", "slug", "richtext", "belongs_to", "many_to_many", "string_array"}
	typeSet := make(map[string]bool, len(types))
	for _, tp := range types {
		typeSet[tp] = true
	}
	for _, req := range required {
		if !typeSet[req] {
			t.Errorf("ValidFieldTypes() missing %q", req)
		}
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "resource.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempYAML: %v", err)
	}
	return path
}

// A field that only applies when another field holds a particular value.
//
// One modifier produces both halves: the admin predicates that hide the input,
// and the model tag internal/fieldtypes enforces for every other writer. A
// condition naming a field that does not exist is refused here, because the
// predicate it would produce reads an undefined property and the input simply
// never appears, with nothing saying why.
func TestWhenModifier(t *testing.T) {
	def, err := ParseInlineFields("Project",
		"kind:select:one_off=One-off|recurring=Recurring, "+
			"cycle:select:monthly=Monthly|annually=Annually:when=kind=recurring")
	if err != nil {
		t.Fatalf("a valid condition was refused: %v", err)
	}
	cycle := def.Fields[1]
	if cycle.ShowWhen == nil {
		t.Fatal("when= did not reach the field")
	}
	if cycle.ShowWhen.Field != "kind" || cycle.ShowWhen.Value != "recurring" {
		t.Errorf("parsed %+v", *cycle.ShowWhen)
	}

	// Modifiers used to be ignored on a select, which returned as soon as it
	// had read its options.
	def, err = ParseInlineFields("Project",
		"kind:select:a=A|b=B, note:select:x=X|y=Y:optional:when=kind=a")
	if err != nil {
		t.Fatalf("optional plus when= on a select: %v", err)
	}
	if def.Fields[1].Required {
		t.Error("optional was ignored on a select")
	}

	for name, spec := range map[string]string{
		"unknown field": "kind:select:a=A|b=B, cycle:string:when=nosuch=a",
		"itself":        "cycle:string:when=cycle=a",
		"no value":      "kind:select:a=A|b=B, cycle:string:when=kind",
		"empty value":   "kind:select:a=A|b=B, cycle:string:when=kind=",
	} {
		if _, err := ParseInlineFields("Project", spec); err == nil {
			t.Errorf("%s: accepted %q", name, spec)
		}
	}
}

// Modifiers after a type's own argument, on every type that takes one.
//
// Five branches read their argument and returned, so a belongs_to could not be
// made optional and when= did nothing on a relationship, with no error either
// time. Each shape is pinned here because the failure is silence.
func TestModifiersFollowEveryTypeArgument(t *testing.T) {
	for name, c := range map[string]struct {
		spec  string
		check func(t *testing.T, f Field)
	}{
		"belongs_to with a model": {
			"owner:belongs_to:User:optional",
			func(t *testing.T, f Field) {
				if f.RelatedModel != "User" || f.Required {
					t.Errorf("model=%q required=%v", f.RelatedModel, f.Required)
				}
			},
		},
		"belongs_to with a modifier in the model's place": {
			"client:belongs_to:optional",
			func(t *testing.T, f Field) {
				if f.RelatedModel != "" || f.Required {
					t.Errorf("model=%q required=%v", f.RelatedModel, f.Required)
				}
			},
		},
		"belongs_to with nothing after it": {
			"client:belongs_to",
			func(t *testing.T, f Field) {
				if !f.Required {
					t.Error("a plain belongs_to stopped being required")
				}
			},
		},
		"many_to_many": {
			"tags:many_to_many:Tag:required",
			func(t *testing.T, f Field) {
				if f.RelatedModel != "Tag" || !f.Required {
					t.Errorf("model=%q required=%v", f.RelatedModel, f.Required)
				}
			},
		},
		"file": {
			"photo:file:image:required",
			func(t *testing.T, f Field) {
				if len(f.FileAccepts) != 1 || f.FileAccepts[0] != "image" || !f.Required {
					t.Errorf("accepts=%v required=%v", f.FileAccepts, f.Required)
				}
			},
		},
		"slug with a source": {
			"code:slug:title:required",
			func(t *testing.T, f Field) {
				if f.SlugSource != "title" || !f.Required {
					t.Errorf("source=%q required=%v", f.SlugSource, f.Required)
				}
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			def, err := ParseFieldsUnvalidated("Thing", c.spec)
			if err != nil {
				t.Fatalf("%s: %v", c.spec, err)
			}
			c.check(t, def.Fields[0])
		})
	}

	// And a condition on a relationship, which is where one is most often
	// wanted and where it did nothing at all.
	def, err := ParseInlineFields("Invoice",
		"kind:select:one_off=One-off|recurring=Recurring, plan:belongs_to:Plan:when=kind=recurring")
	if err != nil {
		t.Fatalf("when= on a belongs_to: %v", err)
	}
	if def.Fields[1].ShowWhen == nil {
		t.Fatal("when= was dropped on a belongs_to")
	}
}
