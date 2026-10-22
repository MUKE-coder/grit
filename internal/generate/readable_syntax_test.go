package generate

import "testing"

// The readable field syntax, against the positional one it rewrites to.
func TestReadableFieldSyntax(t *testing.T) {
	for name, c := range map[string]struct {
		readable   string
		positional string
	}{
		"options":                {"cycle:select[monthly, annually]", "cycle:select:monthly|annually"},
		"modifier after":         {"cycle:select[a, b]:unique", "cycle:select:a|b:unique"},
		"modifier and condition": {"cycle:select[a, b]:optional if kind == r", "cycle:select:a|b:optional:when=kind=r"},
		"options with label":     {`plan:select[basic, pro="Pro plan"]`, "plan:select:basic|pro=Pro plan"},
		"condition":              {"cycle:select[monthly, annually] if kind == recurring", "cycle:select:monthly|annually:when=kind=recurring"},
		"condition alone":        {"paid_on:date if status == paid", "paid_on:date:when=status=paid"},
		"relationship":           {"author:belongs_to[User]", "author:belongs_to:User"},
		"files":                  {"docs:files[pdf, doc]", "docs:files:[pdf,doc]"},
		"quoted value":           {`cycle:select[a, b] if kind == "one off"`, "cycle:select:a|b:when=kind=one off"},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := normalizeFieldSpec(c.readable)
			if err != nil {
				t.Fatalf("%s: %v", c.readable, err)
			}
			if got != c.positional {
				t.Errorf("normalised to %q, want %q", got, c.positional)
			}
		})
	}
}

// Both forms have to produce the same field, or there are two syntaxes rather
// than one with two spellings.
func TestBothSyntaxesAgree(t *testing.T) {
	readable, err := ParseInlineFields("Project",
		"kind:select[one_off, recurring], cycle:select[monthly, annually] if kind == recurring")
	if err != nil {
		t.Fatalf("readable: %v", err)
	}
	positional, err := ParseInlineFields("Project",
		"kind:select:one_off|recurring, cycle:select:monthly|annually:when=kind=recurring")
	if err != nil {
		t.Fatalf("positional: %v", err)
	}
	for i := range readable.Fields {
		a, b := readable.Fields[i], positional.Fields[i]
		if a.Name != b.Name || a.Type != b.Type || len(a.Options) != len(b.Options) {
			t.Fatalf("field %d differs: %+v vs %+v", i, a, b)
		}
		if (a.ShowWhen == nil) != (b.ShowWhen == nil) {
			t.Fatalf("field %d: conditions differ", i)
		}
		if a.ShowWhen != nil && *a.ShowWhen != *b.ShowWhen {
			t.Fatalf("field %d: %+v vs %+v", i, *a.ShowWhen, *b.ShowWhen)
		}
	}
	cycle := readable.Fields[1]
	if cycle.ShowWhen == nil || cycle.ShowWhen.Field != "kind" || cycle.ShowWhen.Value != "recurring" {
		t.Fatalf("condition = %+v", cycle.ShowWhen)
	}
	// Labels are derived from the values when none is given.
	if cycle.Options[0].Label != "Monthly" {
		t.Errorf("label = %q, want Monthly", cycle.Options[0].Label)
	}
}

func TestReadableSyntaxErrors(t *testing.T) {
	for name, spec := range map[string]string{
		"no operator":     "cycle:select[a, b] if kind recurring",
		"empty value":     "cycle:select[a, b] if kind ==",
		"unclosed":        "cycle:select[a, b",
		"two in brackets": "author:belongs_to[User, Admin]",
		"empty brackets":  "email:email[]",
	} {
		if _, err := normalizeFieldSpec(spec); err == nil {
			t.Errorf("%s: accepted %q", name, spec)
		}
	}
}
