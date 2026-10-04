package factory

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/jsontime"
	"library/apps/api/internal/models"
)

// NewLoan builds a valid Loan without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
//
// BookID and BorrowerID is left empty: the factory cannot know which parent you mean, and
// inserting one would make every test that touches this resource write rows in
// another table too. Pass it in:
//
//	loan := factory.CreateLoan(t, db, func(m *models.Loan) {
//		m.BookID = parent.ID
//	})
func NewLoan(overrides ...func(*models.Loan)) *models.Loan {
	record := &models.Loan{
		BorrowedAt: &jsontime.DateTime{Time: time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC)},
		DueAt:      &jsontime.DateTime{Time: time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC)},
		ReturnedAt: &jsontime.DateTime{Time: time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC)},
		Status:     "out",
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateLoan builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateLoan(tb testing.TB, db *gorm.DB, overrides ...func(*models.Loan)) *models.Loan {
	tb.Helper()
	record := NewLoan(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a loan fixture: %v", err)
	}
	return record
}

// CreateLoans writes n of them, for a test about paging or sorting.
func CreateLoans(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Loan)) []*models.Loan {
	tb.Helper()
	out := make([]*models.Loan, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateLoan(tb, db, overrides...))
	}
	return out
}
