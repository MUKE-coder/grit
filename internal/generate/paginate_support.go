package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// ensurePaginateSupport brings an older project's paginate package forward
// before a handler that depends on it is written.
//
// The specific failure: a public handler declares
// paginate.Config{RangeFilterable: ...} so a storefront can ask for
// ?price_min=50, and a paginate.go generated before that field existed makes
// the project fail to compile on "unknown field RangeFilterable". The command
// reports success and leaves a broken build, which is the worst outcome a
// generator has.
//
// Only the public surface needs this, so it runs only for --public. paginate.go
// is replaced only when the manifest proves nobody has edited it; a modified
// copy is left alone with a warning naming the one field to add, because
// silently overwriting somebody's pagination is worse than a build error they
// can read.
func (g *Generator) ensurePaginateSupport() error {
	if !g.Definition.Public {
		return nil
	}

	path := filepath.Join(g.APIRoot(), "internal", "paginate", "paginate.go")
	data, err := os.ReadFile(path)
	if err != nil {
		// No paginate package at all. Not this function's problem: the handler
		// templates have always depended on it, so the project is older than
		// anything --public can help with.
		return nil
	}
	if strings.Contains(string(data), "RangeFilterable") {
		return nil
	}

	body := scaffold.APIPaginateGo()
	if g.refreshIfUnchanged(path, body) {
		fmt.Println("  ✓ Updated internal/paginate/paginate.go with range filters")
		return nil
	}

	yellow := color.New(color.FgHiYellow)
	yellow.Printf("\n  ⚠ %s predates range filters and the new handler will not compile.\n",
		filepath.Join("internal", "paginate", "paginate.go"))
	fmt.Println("    That file has local edits, so it was left alone.")
	fmt.Println("    Add a RangeFilterable map[string]bool field to paginate.Config and")
	fmt.Println("    apply it in List() as a >= / <= pair on ?<col>_min / ?<col>_max.")
	fmt.Println()
	return nil
}
