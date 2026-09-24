package prompt

import (
	"os"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The picker and the --theme flag have to offer the same themes.
//
// They are separate lists, and when v3.322.0 added five themes only the flag
// learned about them. Nothing failed: `grit new` simply showed three options,
// so for everybody who does not pass --theme the new layouts had not shipped.
// A missing option is invisible in a way a rejected flag is not.
func TestThePickerOffersEveryTheme(t *testing.T) {
	source, err := os.ReadFile("prompt.go")
	if err != nil {
		t.Fatalf("reading prompt.go: %v", err)
	}
	picker := string(source)

	for _, theme := range scaffold.ValidThemes {
		if !strings.Contains(picker, `", "`+theme+`")`) {
			t.Errorf("theme %q is accepted by --theme but the interactive picker never offers it", theme)
		}
	}
}

// And the picker must not offer one the flag would then reject.
//
// Scoped to the theme select rather than the whole file: the same form picks
// the architecture, the frontend, the admin style and the database, and a
// skip-list of their values would need editing every time one of them gains
// an option.
func TestThePickerOffersNothingTheFlagRefuses(t *testing.T) {
	block := themeSelectBlock(t)

	valid := make(map[string]bool, len(scaffold.ValidThemes))
	for _, theme := range scaffold.ValidThemes {
		valid[theme] = true
	}

	for _, value := range optionValues(block) {
		if !valid[value] {
			t.Errorf("the picker offers theme %q, which --theme would reject", value)
		}
	}
}

// themeSelectBlock returns the source of the theme select, from its Key to the
// end of its Options call.
func themeSelectBlock(t *testing.T) string {
	t.Helper()
	source, err := os.ReadFile("prompt.go")
	if err != nil {
		t.Fatalf("reading prompt.go: %v", err)
	}
	s := string(source)

	start := strings.Index(s, `Key("theme")`)
	if start < 0 {
		t.Fatal("no theme select in prompt.go")
	}
	rest := s[start:]
	end := strings.Index(rest, "Value(&theme)")
	if end < 0 {
		t.Fatal("the theme select does not bind to &theme")
	}
	return rest[:end]
}

// optionValues pulls the value from each huh.NewOption(label, value) call.
func optionValues(block string) []string {
	var values []string
	for _, line := range strings.Split(block, "\n") {
		if !strings.Contains(line, "huh.NewOption(") {
			continue
		}
		open := strings.LastIndex(line, `", "`)
		if open < 0 {
			continue
		}
		rest := line[open+4:]
		if end := strings.Index(rest, `"`); end >= 0 {
			values = append(values, rest[:end])
		}
	}
	return values
}
