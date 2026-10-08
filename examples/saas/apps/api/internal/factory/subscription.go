package factory

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/jsontime"
	"saas/apps/api/internal/models"
)

// NewSubscription builds a valid Subscription without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
//
// PlanID and UserID is left empty: the factory cannot know which parent you mean, and
// inserting one would make every test that touches this resource write rows in
// another table too. Pass it in:
//
//	subscription := factory.CreateSubscription(t, db, func(m *models.Subscription) {
//		m.PlanID = parent.ID
//	})
func NewSubscription(overrides ...func(*models.Subscription)) *models.Subscription {
	n := next()

	record := &models.Subscription{
		Status:           "trialing",
		Seats:            n,
		CurrentPeriodEnd: &jsontime.Date{Time: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)},
		CanceledOn:       &jsontime.Date{Time: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateSubscription builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateSubscription(tb testing.TB, db *gorm.DB, overrides ...func(*models.Subscription)) *models.Subscription {
	tb.Helper()
	record := NewSubscription(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a subscription fixture: %v", err)
	}
	return record
}

// CreateSubscriptions writes n of them, for a test about paging or sorting.
func CreateSubscriptions(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Subscription)) []*models.Subscription {
	tb.Helper()
	out := make([]*models.Subscription, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateSubscription(tb, db, overrides...))
	}
	return out
}
