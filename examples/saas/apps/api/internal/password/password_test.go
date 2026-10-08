package password

import (
	"strings"
	"testing"
)

func TestTheRulesAndTheChecklistAgree(t *testing.T) {
	ids := map[string]bool{}
	for _, rule := range Rules() {
		if rule.ID == "" || rule.Label == "" {
			t.Errorf("rule %+v has no id or no label, so the checklist cannot show it", rule)
		}
		if ids[rule.ID] {
			t.Errorf("two rules share the id %q, so the checklist cannot tell them apart", rule.ID)
		}
		ids[rule.ID] = true
	}
	// Every id Check can return has to be a rule somebody can read.
	for _, id := range Check("a", "ada@example.com") {
		if !ids[id] {
			t.Errorf("Check returned %q, which is not in Rules()", id)
		}
	}
}

func TestCheck(t *testing.T) {
	cases := []struct {
		name      string
		password  string
		about     []string
		wantFails []string
	}{
		{"a good one", "harbour-lamp-97", nil, nil},
		{"too short", "ab1", nil, []string{"length"}},
		{"letters only", "abcdefghij", nil, []string{"variety"}},
		{"digits only", "1234567890", nil, []string{"variety", "not-common"}},
		{"a space counts as variety", "correct horse staple", nil, nil},
		{"the one everybody tries", "password1", nil, []string{"not-common"}},
		{"case does not hide a common one", "PassWord1", nil, []string{"not-common"}},
		{"their email is in it", "ada.okello-2026", []string{"ada.okello@example.com"}, []string{"not-personal"}},
		{"their company is in it", "example-street-12", []string{"ada@example.com"}, []string{"not-personal"}},
		{"their surname is in it", "okello-was-here", []string{"Ada", "Okello"}, []string{"not-personal"}},
		{"a short name is a coincidence", "ng-harbour-lamp", []string{"Ng"}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Check(tc.password, tc.about...)
			if len(got) != len(tc.wantFails) {
				t.Fatalf("Check(%q) = %v, want %v", tc.password, got, tc.wantFails)
			}
			for i := range got {
				if got[i] != tc.wantFails[i] {
					t.Errorf("Check(%q)[%d] = %q, want %q", tc.password, i, got[i], tc.wantFails[i])
				}
			}
		})
	}
}

// The message is what the person actually reads when the save is refused, so
// it has to name the rules rather than say "invalid password".
func TestMessageNamesWhatFailed(t *testing.T) {
	msg := Message(Check("abc", "ada@example.com"))
	if !strings.Contains(msg, "8 characters") {
		t.Errorf("message = %q, want it to say what is wrong", msg)
	}
	if Message(nil) != "" {
		t.Error("a password that passed should produce no message")
	}
}

// "a" is short AND all letters AND contains nothing personal: the order of the
// failures is the order of the checklist, so the UI can line them up.
func TestFailuresComeBackInChecklistOrder(t *testing.T) {
	got := Check("a")
	want := []string{"length", "variety"}
	if len(got) != len(want) {
		t.Fatalf("Check(\"a\") = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Check(\"a\") = %v, want %v", got, want)
		}
	}
}
