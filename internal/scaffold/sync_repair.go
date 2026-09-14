package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// syncConnectHook installs the soft-delete hook sync pull relies on.
const syncConnectHook = `	// A soft delete also moves updated_at, so sync pull pages on one indexed
	// column and still carries deletes. See internal/sync.
	if err := sync.Install(db); err != nil {
		return nil, fmt.Errorf("installing the sync soft-delete hook: %w", err)
	}

`

// migrateSyncBackfillHook follows the search indexes in cmd/migrate.
const migrateSyncBackfillHook = `	// Rows soft-deleted before a delete also moved updated_at get it now, so
	// sync pull, which pages on updated_at, carries those deletes too.
	if err := sync.BackfillUpdatedAt(db, models.Models()...); err != nil {
		log.Printf("Some soft-deleted rows were not brought up to date for sync: %v", err)
	}
`

var unindexedUpdatedAtRe = regexp.MustCompile("(?m)^(\\tUpdatedAt[ \\t]+time\\.Time[ \\t]+)`json:\"updated_at\"`$")

// repairSyncPull brings a project up to the fix for H20 in the contact-app
// review. Sync pull ordered by the later of updated_at and deleted_at, which
// sorted the whole table on every pull and failed outright on MySQL, and its
// cursor was a bare time, so rows sharing a timestamp were lost at a page
// boundary: of 1,200 such rows, 500 reached the client.
//
// handlers/sync.go arrives whole, as it has since v3.241.0. internal/sync gets
// the soft-delete hook the new pull relies on, database.go installs it, and
// generated and blog models index updated_at where the line is still what Grit
// wrote. grit migrate builds the indexes and backfills older deletes.
func repairSyncPull(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	syncDir := filepath.Join(apiRoot, "internal", "sync")
	if !fileExists(filepath.Join(syncDir, "registry.go")) {
		return nil
	}
	for path, content := range map[string]string{
		filepath.Join(syncDir, "softdelete.go"):      syncSoftDeleteGo(),
		filepath.Join(syncDir, "softdelete_test.go"): syncSoftDeleteTestGo(),
	} {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if database := filepath.Join(apiRoot, "internal", "database", "database.go"); fileExists(database) {
		module := opts.Module()
		if err := repairSourceFile(root, m, database, func(src string) (string, []string, []string) {
			return repairSyncWiringSource(src, module)
		}); err != nil {
			return err
		}
	}
	models, err := goSources(filepath.Join(apiRoot, "internal", "models"))
	if err != nil {
		return err
	}
	for _, path := range models {
		base := filepath.Base(path)
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return repairUpdatedAtIndexSource(src, base)
		}); err != nil {
			return err
		}
	}
	return nil
}

func repairSyncWiringSource(src, module string) (string, []string, []string) {
	if strings.Contains(src, "sync.Install(db)") || !strings.Contains(src, "func Connect(") {
		return src, nil, nil
	}
	const anchor = "\tsqlDB, err := db.DB()\n"
	if strings.Count(src, anchor) != 1 || strings.Contains(src, "\t\"sync\"\n") {
		return src, nil, []string{"database.go is not the file Grit wrote: call sync.Install(db) after connecting, or deletes stop reaching offline clients"}
	}
	out := strings.Replace(src, anchor, syncConnectHook+anchor, 1)
	var ok bool
	if out, ok = addImportGroup(out, module+"/internal/sync"); !ok {
		return src, nil, []string{"could not add the sync import to database.go: call sync.Install(db) after connecting"}
	}
	return out, []string{"a soft delete also moves updated_at, which sync pull pages on"}, nil
}

func repairUpdatedAtIndexSource(src, base string) (string, []string, []string) {
	generated := generatedArchivedRe.MatchString(src)
	blog := base == "blog.go" && strings.Contains(src, "type Blog struct {")
	if !generated && !blog {
		return src, nil, nil
	}
	if !unindexedUpdatedAtRe.MatchString(src) {
		return src, nil, nil
	}
	return unindexedUpdatedAtRe.ReplaceAllString(src, "${1}`gorm:\"index\" json:\"updated_at\"`"),
		[]string{"updated_at is indexed, so sync pull reads changes in order from the index"}, nil
}
