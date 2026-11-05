package doctor

import (
	"strings"
	"testing"
)

// grit#179: a bool with a database default of true can never be saved as
// false. GORM leaves a zero value out of the INSERT, the database fills in
// true, and GORM writes that back into the struct. A form that saves
// "unticked" stores "ticked", and nothing errors.
//
// Found three times in this project's own models and twice in a real
// application, once on marketing consent: the default was true, so a person
// who declined was recorded as having agreed.
func TestABoolThatCanNeverBeFalseIsAnError(t *testing.T) {
	const src = "\tActive bool `gorm:\"default:true\" json:\"active\"`\n" +
		"\tConsent bool `gorm:\"index;default:TRUE;not null\" json:\"consent\"`\n"

	p := &project{
		resources: []string{"Contact"},
		models:    map[string]*model{"Contact": {Name: "Contact", Source: src}},
	}

	findings := checkBoolDefaultTrue(p)
	if len(findings) != 2 {
		t.Fatalf("got %d findings, want 2 (one per field); the check reads the gorm tag, "+
			"and a different spelling of the default must not slip past", len(findings))
	}
	for _, f := range findings {
		if f.Level != "error" {
			t.Errorf("%s is a %s; this stores the wrong value with no error anywhere, "+
				"which is worse than most warnings", f.Resource, f.Level)
		}
		if !strings.Contains(f.Fix, "*bool") {
			t.Errorf("the fix for %s does not mention *bool, which is the option that "+
				"keeps a default: %q", f.Resource, f.Fix)
		}
	}
}

// The shapes that are fine, because a check that fires on them is one people
// switch off.
func TestTheSafeBoolShapesAreNotFlagged(t *testing.T) {
	cases := map[string]string{
		"no default":        "\tActive bool `gorm:\"\" json:\"active\"`\n",
		"default false":     "\tActive bool `gorm:\"default:false\" json:\"active\"`\n",
		"default zero":      "\tCount int `gorm:\"default:0\" json:\"count\"`\n",
		"pointer with true": "\tActive *bool `gorm:\"default:true\" json:\"active\"`\n",
		"a string default":  "\tRole string `gorm:\"default:USER\" json:\"role\"`\n",
	}
	for name, src := range cases {
		p := &project{
			resources: []string{"Thing"},
			models:    map[string]*model{"Thing": {Name: "Thing", Source: src}},
		}
		if findings := checkBoolDefaultTrue(p); len(findings) != 0 {
			t.Errorf("%s was flagged: %q", name, findings[0].Message)
		}
	}
}

// A *bool is the one case worth saying out loud: nil and false are different
// values to GORM, so the zero-value omission does not apply.
func TestAPointerBoolKeepsItsDefault(t *testing.T) {
	src := "\tActive *bool `gorm:\"default:true\" json:\"active\"`\n"
	p := &project{
		resources: []string{"Thing"},
		models:    map[string]*model{"Thing": {Name: "Thing", Source: src}},
	}
	if findings := checkBoolDefaultTrue(p); len(findings) != 0 {
		t.Error("a *bool with default:true was flagged. nil and false are different " +
			"values to GORM, so a pointer has no zero-value problem, and it is one of " +
			"the two fixes this check suggests.")
	}
}
