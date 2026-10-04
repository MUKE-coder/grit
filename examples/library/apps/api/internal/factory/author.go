package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
)

// NewAuthor builds a valid Author without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
func NewAuthor(overrides ...func(*models.Author)) *models.Author {
	n := next()

	record := &models.Author{
		Name:    fmt.Sprintf("Author %d", n),
		Bio:     fmt.Sprintf("Bio %d", n),
		Website: fmt.Sprintf("https://example.com/%d", n),
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateAuthor builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateAuthor(tb testing.TB, db *gorm.DB, overrides ...func(*models.Author)) *models.Author {
	tb.Helper()
	record := NewAuthor(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a author fixture: %v", err)
	}
	return record
}

// CreateAuthors writes n of them, for a test about paging or sorting.
func CreateAuthors(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Author)) []*models.Author {
	tb.Helper()
	out := make([]*models.Author, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateAuthor(tb, db, overrides...))
	}
	return out
}
