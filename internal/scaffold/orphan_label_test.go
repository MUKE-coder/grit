package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	labelFor   = regexp.MustCompile(`htmlFor=\{([\w.]+)\}`)
	controlID  = regexp.MustCompile(`\bid=\{([\w.]+)\}`)
	tmplHeader = regexp.MustCompile(`(?m)^func (\w+)\([^)]*\) string \{`)
)

// A label has to point at a control that exists.
//
// TextareaField and SelectField both computed a fieldId, put it on their label
// and never on the control, so every generated form announced its Description
// as an unlabelled text area and its pickers as unlabelled buttons. This is the
// same fault the project has fixed before, which is the reason to check for it
// rather than to fix it once.
func TestNoLabelPointsAtNothing(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		src := strings.ReplaceAll(string(raw), "\r\n", "\n")

		starts := tmplHeader.FindAllStringSubmatchIndex(src, -1)
		for i, loc := range starts {
			end := len(src)
			if i+1 < len(starts) {
				end = starts[i+1][0]
			}
			body := src[loc[0]:end]
			fn := src[loc[2]:loc[3]]

			ids := map[string]bool{}
			for _, m := range controlID.FindAllStringSubmatch(body, -1) {
				ids[m[1]] = true
			}
			for _, m := range labelFor.FindAllStringSubmatch(body, -1) {
				checked++
				if !ids[m[1]] {
					t.Errorf("%s: %s has a label with htmlFor={%s} and no control carrying that id, so the label names nothing",
						name, fn, m[1])
				}
			}
		}
	}
	if checked == 0 {
		t.Fatal("no htmlFor was examined, so this test proves nothing")
	}
	t.Logf("checked %d labels", checked)
}
