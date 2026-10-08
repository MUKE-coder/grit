package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"saas/apps/api/internal/models"
	"saas/apps/api/internal/money"
)

// NewPlan builds a valid Plan without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
func NewPlan(overrides ...func(*models.Plan)) *models.Plan {
	n := next()

	record := &models.Plan{
		Name:     fmt.Sprintf("Plan %d", n),
		Price:    money.New(1999+int64(n), "USD"),
		Interval: "monthly",
		Seats:    n,
		Blurb:    fmt.Sprintf("Blurb %d", n),
		Active:   true,
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreatePlan builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreatePlan(tb testing.TB, db *gorm.DB, overrides ...func(*models.Plan)) *models.Plan {
	tb.Helper()
	record := NewPlan(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a plan fixture: %v", err)
	}
	return record
}

// CreatePlans writes n of them, for a test about paging or sorting.
func CreatePlans(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Plan)) []*models.Plan {
	tb.Helper()
	out := make([]*models.Plan, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreatePlan(tb, db, overrides...))
	}
	return out
}
