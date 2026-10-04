package factory

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/jsontime"
	"library/apps/api/internal/models"
	"library/apps/api/internal/money"
)

// NewBook builds a valid Book without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
//
// AuthorID is left empty: the factory cannot know which parent you mean, and
// inserting one would make every test that touches this resource write rows in
// another table too. Pass it in:
//
//	book := factory.CreateBook(t, db, func(m *models.Book) {
//		m.AuthorID = parent.ID
//	})
func NewBook(overrides ...func(*models.Book)) *models.Book {
	n := next()

	record := &models.Book{
		Title:     fmt.Sprintf("Book %d", n),
		Isbn:      fmt.Sprintf("Isbn %d", n),
		Summary:   fmt.Sprintf("Summary %d", n),
		Published: &jsontime.Date{Time: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)},
		Price:     money.New(1999+int64(n), "USD"),
		Genre:     "fiction",
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateBook builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateBook(tb testing.TB, db *gorm.DB, overrides ...func(*models.Book)) *models.Book {
	tb.Helper()
	record := NewBook(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a book fixture: %v", err)
	}
	return record
}

// CreateBooks writes n of them, for a test about paging or sorting.
func CreateBooks(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Book)) []*models.Book {
	tb.Helper()
	out := make([]*models.Book, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateBook(tb, db, overrides...))
	}
	return out
}
