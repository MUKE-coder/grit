package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// The web navbar test shipped to every project asserted a "Components" link
// that its own stub never rendered, left over from the /components page removed
// in v3.31.78. It could not pass, so `pnpm test` in apps/web failed on every
// fresh project: visibly on Next, and hidden on Vite only because that suite
// could not start at all.
//
// Every link the test asks for has to be one the stub it renders contains.
func TestWebNavbarTestAssertsOnlyWhatItRenders(t *testing.T) {
	src := webNavbarTest()
	asked := regexp.MustCompile(`getByRole\("link", \{ name: /([a-z]+)/i \}\)`).FindAllStringSubmatch(src, -1)
	if len(asked) == 0 {
		t.Fatal("the navbar test asserts no links at all")
	}
	for _, m := range asked {
		label := m[1]
		if !strings.Contains(strings.ToLower(src), ">"+label+"</a>") {
			t.Errorf("the test asserts a %q link its stub never renders", label)
		}
	}
}
