// Package erasure is what a right-to-erasure request deletes, and the deletion.
//
// A registry rather than a list, so a table joins by declaring itself: the
// framework's own tables from models/erasable.go, and every resource generated
// with --owned-by from its model file. It used to be nine framework tables
// named in the GDPR service, and erasing a user left every app record they
// owned in place while the deletion journal certified the erasure complete.
//
// It imports nothing from the app, because the models import it.
package erasure

import (
	"fmt"
	"regexp"
	"sync"

	"gorm.io/gorm"
)

type target struct {
	model  interface{}
	column string
}

var (
	mu      sync.RWMutex
	targets []target

	// identifier is what an owner column has to look like: it goes into the
	// WHERE clause as text, not as a bound parameter.
	identifier = regexp.MustCompile("^[a-z_][a-z0-9_]*$")
)

// Register marks rows of model as belonging to the user named in ownerColumn,
// so erasing that user deletes them. Call it from the model file's init().
func Register(model interface{}, ownerColumn string) {
	mu.Lock()
	defer mu.Unlock()
	targets = append(targets, target{model: model, column: ownerColumn})
}

// Scrub hard-deletes every registered row belonging to userID and anonymizes
// the user row in place, inside tx. It returns what was deleted, per table.
//
// Unscoped, because several models carry gorm.DeletedAt and a normal Delete
// only soft-deletes: the row, and the personal data in it, would remain.
// Tables that do not exist are skipped, so a database migrated with a subset of
// the models, a test's for one, can still run an erasure.
func Scrub(tx *gorm.DB, userID string) (map[string]int64, int, error) {
	mu.RLock()
	list := append([]target(nil), targets...)
	mu.RUnlock()

	counts := map[string]int64{}
	total := 0
	for _, t := range list {
		if !identifier.MatchString(t.column) {
			return nil, 0, fmt.Errorf("refusing the owner column %q", t.column)
		}
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(t.model); err != nil {
			return nil, 0, fmt.Errorf("resolving %T: %w", t.model, err)
		}
		table := stmt.Schema.Table
		if !tx.Migrator().HasTable(table) {
			continue
		}
		var n int64
		if err := tx.Unscoped().Model(t.model).Where(t.column+" = ?", userID).Count(&n).Error; err != nil {
			return nil, 0, fmt.Errorf("counting %s: %w", table, err)
		}
		if n > 0 {
			if err := tx.Unscoped().Where(t.column+" = ?", userID).Delete(t.model).Error; err != nil {
				return nil, 0, fmt.Errorf("deleting %s: %w", table, err)
			}
		}
		counts[table] += n
		total += int(n)
	}

	// Anonymize the user row in place: scrub every personal column, keep the id
	// so references resolve to a tombstone, and keep a unique, non-routable
	// email so the unique index stays satisfiable if the row is read again.
	// By table name, so no model hook runs over data that is being destroyed.
	if err := tx.Table("users").Where("id = ?", userID).Updates(map[string]interface{}{
		"first_name":  "Erased",
		"last_name":   "User",
		"email":       "erased-" + userID + "@deleted.invalid",
		"password":    "",
		"avatar":      "",
		"job_title":   "",
		"bio":         "",
		"ip_address":  "",
		"mac_address": "",
		"google_id":   "",
		"github_id":   "",
		"active":      false,
		"role":        "USER",
	}).Error; err != nil {
		return nil, 0, fmt.Errorf("anonymizing user: %w", err)
	}
	return counts, total, nil
}
