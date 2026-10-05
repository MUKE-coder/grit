package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// NewCollection builds a valid Collection without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
func NewCollection(overrides ...func(*models.Collection)) *models.Collection {
	n := next()

	record := &models.Collection{
		Title:       fmt.Sprintf("Collection %d", n),
		Description: fmt.Sprintf("A description for fixture %d", n),
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateCollection builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateCollection(tb testing.TB, db *gorm.DB, overrides ...func(*models.Collection)) *models.Collection {
	tb.Helper()
	record := NewCollection(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a collection fixture: %v", err)
	}
	return record
}

// CreateCollections writes n of them, for a test about paging or sorting.
func CreateCollections(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Collection)) []*models.Collection {
	tb.Helper()
	out := make([]*models.Collection, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateCollection(tb, db, overrides...))
	}
	return out
}
