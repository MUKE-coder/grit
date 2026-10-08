package factory

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/jsontime"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/money"
)

// NewInvoice builds a valid Invoice without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
//
// SubscriptionID and UserID is left empty: the factory cannot know which parent you mean, and
// inserting one would make every test that touches this resource write rows in
// another table too. Pass it in:
//
//	invoice := factory.CreateInvoice(t, db, func(m *models.Invoice) {
//		m.SubscriptionID = parent.ID
//	})
func NewInvoice(overrides ...func(*models.Invoice)) *models.Invoice {
	n := next()

	record := &models.Invoice{
		Number:   fmt.Sprintf("Number %d", n),
		Amount:   money.New(1999+int64(n), "USD"),
		Status:   "draft",
		IssuedOn: &jsontime.Date{Time: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)},
		PaidOn:   &jsontime.Date{Time: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateInvoice builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateInvoice(tb testing.TB, db *gorm.DB, overrides ...func(*models.Invoice)) *models.Invoice {
	tb.Helper()
	record := NewInvoice(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a invoice fixture: %v", err)
	}
	return record
}

// CreateInvoices writes n of them, for a test about paging or sorting.
func CreateInvoices(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Invoice)) []*models.Invoice {
	tb.Helper()
	out := make([]*models.Invoice, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateInvoice(tb, db, overrides...))
	}
	return out
}
