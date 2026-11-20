package plugin

import (
	"strings"
	"testing"
)

// Nothing the skeleton writes contains a leaked Go concatenation.
//
// These files are Go raw strings with ` + pkg + ` spliced through them, and a
// backticked markdown span inside one has to leave the raw string for its
// backtick and leave it again for the variable. Get the nesting wrong and the
// concatenation becomes literal text: the README said
//
//	So `files/api/` + pkg + `.go` and ...
//
// which compiles, passes every other test, and is only visible by generating a
// plugin and reading the file. Found that way; this is cheaper.
func TestTheSkeletonWritesNoLeakedConcatenation(t *testing.T) {
	for path, body := range SkeletonFiles("greeter") {
		for _, leak := range []string{"` + ", " + `", `+ pkg +`, `+ name +`} {
			if strings.Contains(body, leak) {
				t.Errorf("%s contains %q, which is Go source that should have been "+
					"evaluated. Check the backtick nesting in the template.", path, leak)
			}
		}
		// The placeholders the installer substitutes are a different thing and
		// are meant to survive: {{MODULE}}, {{API_ROOT}}, {{PROJECT}}.
		if strings.Contains(body, "{{") && !strings.Contains(body, "}}") {
			t.Errorf("%s has an unclosed placeholder", path)
		}
	}
}
