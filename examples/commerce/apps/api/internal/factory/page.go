package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// NewPage builds a valid Page without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
func NewPage(overrides ...func(*models.Page)) *models.Page {
	n := next()

	record := &models.Page{
		Title: fmt.Sprintf("Page %d", n),
		Body:  fmt.Sprintf("A description for fixture %d", n),
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreatePage builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreatePage(tb testing.TB, db *gorm.DB, overrides ...func(*models.Page)) *models.Page {
	tb.Helper()
	record := NewPage(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a page fixture: %v", err)
	}
	return record
}

// CreatePages writes n of them, for a test about paging or sorting.
func CreatePages(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Page)) []*models.Page {
	tb.Helper()
	out := make([]*models.Page, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreatePage(tb, db, overrides...))
	}
	return out
}
