package factory

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/jsontime"
	"saas/apps/api/internal/models"
)

// NewUsageRecord builds a valid UsageRecord without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
//
// SubscriptionID and UserID is left empty: the factory cannot know which parent you mean, and
// inserting one would make every test that touches this resource write rows in
// another table too. Pass it in:
//
//	usageRecord := factory.CreateUsageRecord(t, db, func(m *models.UsageRecord) {
//		m.SubscriptionID = parent.ID
//	})
func NewUsageRecord(overrides ...func(*models.UsageRecord)) *models.UsageRecord {
	n := next()

	record := &models.UsageRecord{
		Metric:     fmt.Sprintf("Metric %d", n),
		Quantity:   n,
		RecordedOn: &jsontime.Date{Time: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateUsageRecord builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateUsageRecord(tb testing.TB, db *gorm.DB, overrides ...func(*models.UsageRecord)) *models.UsageRecord {
	tb.Helper()
	record := NewUsageRecord(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a usagerecord fixture: %v", err)
	}
	return record
}

// CreateUsageRecords writes n of them, for a test about paging or sorting.
func CreateUsageRecords(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.UsageRecord)) []*models.UsageRecord {
	tb.Helper()
	out := make([]*models.UsageRecord, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateUsageRecord(tb, db, overrides...))
	}
	return out
}
