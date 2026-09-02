package scaffold

// apiWebhookDedupTestGo emits internal/webhooks/dedup_test.go.
//
// This exists because the previous version of webhook deduplication was
// entirely documentation. The model's comment described a unique index on
// (provider, external_id), the handler treated a duplicate-key error as
// "already processed", and the index itself lived in an Indexes() method that
// nothing ever called. The column had a plain index, the second INSERT
// succeeded, and every retried delivery ran the handler again.
//
// A payment provider retrying a webhook is not an edge case; it is the normal
// operation of every provider worth integrating. So the constraint gets a test
// that fails loudly if it ever stops being enforced.
func apiWebhookDedupTestGo() string {
	return `package webhooks

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
)

func newDedupDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.WebhookEvent{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func ref(s string) *string { return &s }

func event(provider, external string) models.WebhookEvent {
	e := models.WebhookEvent{
		Provider:  provider,
		EventType: "invoice.paid",
		Payload:   datatypes.JSON(` + "`" + `{"ok":true}` + "`" + `),
		Status:    "pending",
	}
	if external != "" {
		e.ExternalID = ref(external)
	}
	return e
}

// The whole point: a provider retrying the same event must not produce a
// second row, because a second row means the handler runs a second time and
// the payment is recorded twice.
func TestDuplicateExternalIDIsRejected(t *testing.T) {
	db := newDedupDB(t)

	first := event("stripe", "evt_123")
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("first delivery: %v", err)
	}

	second := event("stripe", "evt_123")
	err := db.Create(&second).Error
	if err == nil {
		t.Fatal("a repeated delivery of evt_123 was accepted; the handler would run twice")
	}
	if !IsDuplicateError(err) {
		t.Errorf("error = %v; IsDuplicateError did not recognise it, so the handler "+
			"would answer 500 instead of skipping the duplicate", err)
	}

	var n int64
	db.Model(&models.WebhookEvent{}).Count(&n)
	if n != 1 {
		t.Errorf("rows = %d, want 1", n)
	}
}

// The same event id from a different provider is a different event. Deduping
// on the id alone would drop it.
func TestSameIDFromDifferentProvidersBothStored(t *testing.T) {
	db := newDedupDB(t)

	for _, p := range []string{"stripe", "paystack"} {
		e := event(p, "evt_123")
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("%s: %v", p, err)
		}
	}

	var n int64
	db.Model(&models.WebhookEvent{}).Count(&n)
	if n != 2 {
		t.Errorf("rows = %d, want 2: two providers sharing an event id collided", n)
	}
}

// A provider that supplies no event id cannot be deduplicated, so each
// delivery is its own row. Storing "" instead of NULL would make the second
// one look like a duplicate of the first and silently discard it.
func TestEventsWithoutAnExternalIDAreAllKept(t *testing.T) {
	db := newDedupDB(t)

	for i := 0; i < 3; i++ {
		e := event("partner", "")
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("delivery %d: %v", i, err)
		}
	}

	var n int64
	db.Model(&models.WebhookEvent{}).Count(&n)
	if n != 3 {
		t.Errorf("rows = %d, want 3: anonymous events collided on the unique index", n)
	}
}

// The constraint has to be in the database, not in a check the handler does
// first. Two deliveries arriving at once both pass a SELECT and only the
// index stops the second INSERT.
func TestTheConstraintIsInTheDatabase(t *testing.T) {
	db := newDedupDB(t)

	indexes, err := db.Migrator().GetIndexes(&models.WebhookEvent{})
	if err != nil {
		t.Fatalf("GetIndexes: %v", err)
	}
	for _, idx := range indexes {
		unique, ok := idx.Unique()
		if !ok || !unique {
			continue
		}
		cols := idx.Columns()
		if len(cols) == 2 && cols[0] == "provider" && cols[1] == "external_id" {
			return
		}
	}
	t.Error("no unique index on (provider, external_id); deduplication is decoration")
}
`
}
