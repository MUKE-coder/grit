package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// A resource generated with --owned-by holds rows that belong to a user, so a
// right-to-erasure request for that user has to delete them.
//
// Found building a patient-records app. Erasure deleted rows from a fixed list
// of nine framework tables and anonymized the user; a provider's notes, owned
// by that provider, stayed exactly where they were while the deletion journal
// recorded the erasure as done.

// erasesWithOwner reports whether this resource's rows go when their owner is
// erased. Not for append-only resources: their rows cannot be deleted by
// design, and which of retention or erasure wins is for the project to decide,
// not for the generator to decide quietly.
func (g *Generator) erasesWithOwner() bool {
	return g.Definition.OwnerField() != nil && !g.Definition.AppendOnly
}

// erasureModelInit registers the model with the erasure package from its own
// file, so removing the resource removes the registration with it.
func (g *Generator) erasureModelInit(names Names) string {
	if !g.erasesWithOwner() {
		return ""
	}
	owner := g.Definition.OwnerField()
	return fmt.Sprintf(`
// init registers %[1]s for right-to-erasure: erasing a user deletes the rows
// they own here, by %[2]s. See internal/erasure.
func init() { erasure.Register(&%[1]s{}, %[2]q) }
`, names.Pascal, owner.FKColumnName())
}

// prepareErasure makes sure the erasure package is there for the model to
// import, and says so when the project's GDPR service is too old to use it.
func (g *Generator) prepareErasure() error {
	api := g.APIRoot()
	if err := scaffold.WriteErasurePackage(api, g.Module, false); err != nil {
		return err
	}
	gdpr := filepath.Join(api, "internal", "services", "gdpr.go")
	if data, err := os.ReadFile(gdpr); err == nil && !strings.Contains(string(data), "erasure.Scrub") {
		fmt.Println("  ⚠ This project's right-to-erasure predates owned resources: erasing a user")
		fmt.Println("    will not delete their rows here until you run grit upgrade.")
	}
	return nil
}
