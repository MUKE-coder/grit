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

// hasFileColumn reports whether rows of this resource point at stored
// objects, which is what makes erasing one of them more than a DELETE.
func (g *Generator) hasFileColumn() bool {
	for _, f := range g.Definition.Fields {
		if f.IsFileField() {
			return true
		}
	}
	return false
}

// erasesOwnedFiles reports whether the generated erasure registration carries
// a cleanup that deletes the files on the rows it deletes.
//
// Three conditions, and the third is about the project rather than the
// resource: the cleanup calls files.EraseOwnedFiles, and a copy of
// internal/files from before v3.583.0 does not declare it. prepareErasure
// tries to bring it up to date first and says so when it cannot. Registering
// without the cleanup then keeps the old leak, which is the lesser of the two:
// a generated resource that does not compile would be worse.
func (g *Generator) erasesOwnedFiles() bool {
	return g.erasesWithOwner() && g.hasFileColumn() && g.filesPackageCanErase()
}

// filesPackageCanErase reports whether this project's internal/files declares
// what the generated cleanup calls. Matched on the declaration rather than the
// name, which also appears in the prose around it.
func (g *Generator) filesPackageCanErase() bool {
	lifecycle := filepath.Join(g.APIRoot(), "internal", "files", "lifecycle.go")
	return fileContainsText(lifecycle, "func EraseOwnedFiles(")
}

// erasureModelInit registers the model with the erasure package from its own
// file, so removing the resource removes the registration with it.
func (g *Generator) erasureModelInit(names Names) string {
	if !g.erasesWithOwner() {
		return ""
	}
	owner := g.Definition.OwnerField()
	if g.erasesOwnedFiles() {
		return fmt.Sprintf(`
// init registers %[1]s for right-to-erasure: erasing a user deletes the rows
// they own here, by %[2]s, and the files on them. See internal/erasure.
func init() { erasure.RegisterWithCleanup(&%[1]s{}, %[2]q, erase%[1]sFiles) }

// erase%[1]sFiles deletes the stored objects of the %[1]s rows an erasure is
// about to delete, whoever uploaded them.
//
// Erasing a user already takes their own uploads and every object those own. A
// file somebody else uploaded onto one of these rows, an administrator above
// all, has an uploads row the erasure does not touch, and once this row is
// gone nothing points at the object and the orphan sweep cannot see it either:
// it stayed in the bucket for good while the deletion journal recorded the
// erasure as complete (grit#311).
//
// Inside the erasure transaction and before the rows go, because the file
// columns live on the rows.
func erase%[1]sFiles(tx *gorm.DB, ownerColumn, userID string) error {
	// A disk that is not configured is not an error here. The rows still name
	// keys, but there is no bucket they could be in, and failing the erasure
	// would leave personal data in the database over an object store nobody
	// set up.
	var disk files.Storage
	if configured := storage.Disks.Default(); configured != nil {
		disk = configured
	}
	return files.EraseOwnedFiles(context.Background(), tx, disk, &%[1]s{}, ownerColumn, userID)
}
`, names.Pascal, owner.FKColumnName())
	}
	return fmt.Sprintf(`
// init registers %[1]s for right-to-erasure: erasing a user deletes the rows
// they own here, by %[2]s. See internal/erasure.
func init() { erasure.Register(&%[1]s{}, %[2]q) }
`, names.Pascal, owner.FKColumnName())
}

// confidentialModelInit marks the model confidential from its own file, so
// the places that read below the service layer, a backup archive above all,
// can ask whether this caller may see the table at all.
//
// In the model file rather than a central list, for the same reason the
// erasure registration is: removing the resource removes the registration
// with it, and a central list is a thing to forget.
func (g *Generator) confidentialModelInit(names Names) string {
	if !g.Definition.Confidential {
		return ""
	}
	return fmt.Sprintf(`
// init marks %[1]s confidential: an administrator reads nothing here. Only
// the owner, or a holder of %[3]s, and every such read is
// audited. See internal/authz.ScopeConfidential.
func init() { confidential.Register(%[2]q, %[3]q) }
`, names.Pascal, names.PluralSnake, g.Definition.ConfidentialPermission())
}

// prepareErasure makes sure the erasure package is there for the model to
// import, and says so when the project's GDPR service is too old to use it.
func (g *Generator) prepareErasure() error {
	api := g.APIRoot()
	if err := scaffold.WriteErasurePackage(api, g.Module, false); err != nil {
		return err
	}
	// The confidential registry too, for the same reason and on the same
	// terms: a project generated before it existed has no such package, and
	// a --confidential model would not compile without one.
	if err := scaffold.WriteConfidentialPackage(api, g.Module, false); err != nil {
		return err
	}
	gdpr := filepath.Join(api, "internal", "services", "gdpr.go")
	if data, err := os.ReadFile(gdpr); err == nil && !strings.Contains(string(data), "erasure.Scrub") {
		fmt.Println("  ⚠ This project's right-to-erasure predates owned resources: erasing a user")
		fmt.Println("    will not delete their rows here until you run grit upgrade.")
	}
	if g.hasFileColumn() {
		g.ensureFilesCanErase()
	}
	return nil
}

// ensureFilesCanErase brings internal/files up to the version the generated
// cleanup needs, when the manifest proves nobody has edited it.
//
// Two files, because EraseOwnedFiles reads each row's objects through
// FileRef.OwnedKeys and that lives in the other one: refreshing the half that
// calls it, on a project from before v3.544.0, would leave a package that does
// not compile. So the pair moves together or neither does.
//
// A project that has edited either keeps its copy and is told. The resource is
// still generated, and registers for erasure without the cleanup: that is the
// behaviour of every version before this one, so the fallback is the old leak
// rather than a build nobody asked for.
func (g *Generator) ensureFilesCanErase() {
	dir := filepath.Join(g.APIRoot(), "internal", "files")
	lifecycle := filepath.Join(dir, "lifecycle.go")
	if !fileExists(lifecycle) || g.filesPackageCanErase() {
		return
	}

	ref := filepath.Join(dir, "file_ref.go")
	refReady := fileContainsText(ref, "func (f FileRef) OwnedKeys()")
	if !refReady {
		refReady = g.refreshIfUnchanged(ref, strings.ReplaceAll(scaffold.FilesFileRefGo(), "{{MODULE}}", g.Module))
	}
	if refReady && g.refreshIfUnchanged(lifecycle, strings.ReplaceAll(scaffold.FilesLifecycleGo(), "{{MODULE}}", g.Module)) {
		fmt.Println("  ✓ Updated internal/files: erasing a user now deletes the files on the rows erased with them")
		return
	}
	fmt.Println("  ⚠ This project's internal/files predates files.EraseOwnedFiles, so erasing a user")
	fmt.Println("    will leave the files on these rows in storage. Run grit upgrade, or replace")
	fmt.Println("    internal/files with the current template.")
}
