package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/models"
)

func webhookEventDB(t *testing.T) (*gorm.DB, *WebhookEventService) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.WebhookEvent{}))
	return db, &WebhookEventService{DB: db}
}

func storedEvent(t *testing.T, events *WebhookEventService, externalID, status string, age time.Duration) *models.WebhookEvent {
	t.Helper()
	ref := externalID
	event := models.WebhookEvent{
		Provider:   "stripe",
		EventType:  "invoice.paid",
		ExternalID: &ref,
		Payload:    datatypes.JSON([]byte(`{}`)),
		Status:     status,
	}
	require.NoError(t, events.Store(context.Background(), &event))
	if age > 0 {
		require.NoError(t, events.DB.Model(&event).
			UpdateColumn("created_at", time.Now().Add(-age)).Error)
	}
	return &event
}

// A provider that redelivers a failed event hands it to whichever request claims
// it, and only one may: two handlers for one event is the thing worth avoiding.
func TestClaimIsWonByOneRedelivery(t *testing.T) {
	_, events := webhookEventDB(t)
	ctx := context.Background()
	event := storedEvent(t, events, "evt_failed", "failed", 0)

	claimed, err := events.Claim(ctx, event, 5*time.Minute)
	require.NoError(t, err)
	assert.True(t, claimed)

	// A second redelivery arriving while the first runs finds it pending and new.
	claimed, err = events.Claim(ctx, event, 5*time.Minute)
	require.NoError(t, err)
	assert.False(t, claimed, "two redeliveries both ran the handler for one event")
}

// A pending event is somebody's request still running, until it is old enough
// that the process which stored it must be gone.
func TestClaimWaitsOutARequestInFlight(t *testing.T) {
	_, events := webhookEventDB(t)
	ctx := context.Background()
	fresh := storedEvent(t, events, "evt_fresh", "pending", 0)
	stale := storedEvent(t, events, "evt_stale", "pending", 10*time.Minute)

	claimed, err := events.Claim(ctx, fresh, 5*time.Minute)
	require.NoError(t, err)
	assert.False(t, claimed, "an event somebody is still processing was taken over")

	claimed, err = events.Claim(ctx, stale, 5*time.Minute)
	require.NoError(t, err)
	assert.True(t, claimed, "an event whose process is gone was never picked up again")
}

// The claim counts the attempt, so an event that keeps coming back shows how
// many times it has.
func TestClaimCountsTheAttempt(t *testing.T) {
	db, events := webhookEventDB(t)
	ctx := context.Background()
	event := storedEvent(t, events, "evt_count", "failed", 0)

	for want := 1; want <= 2; want++ {
		// Back to failed, as a failing handler would leave it.
		require.NoError(t, events.RecordOutcome(ctx, event, errors.New("handler said no")))
		claimed, err := events.Claim(ctx, event, 5*time.Minute)
		require.NoError(t, err)
		require.True(t, claimed)
		var fresh models.WebhookEvent
		require.NoError(t, db.First(&fresh, "id = ?", event.ID).Error)
		assert.Equal(t, want, fresh.RetryCount)
	}
}

// Recording success clears the error from the attempt that failed: an event that
// was fixed and ran again should not still carry the old message.
func TestRecordOutcomeClearsAPreviousFailure(t *testing.T) {
	db, events := webhookEventDB(t)
	ctx := context.Background()
	event := storedEvent(t, events, "evt_outcome", "pending", 0)

	require.NoError(t, events.RecordOutcome(ctx, event, errors.New("boom")))
	var failed models.WebhookEvent
	require.NoError(t, db.First(&failed, "id = ?", event.ID).Error)
	assert.Equal(t, "failed", failed.Status)
	assert.Equal(t, "boom", failed.HandlerError)
	require.NotNil(t, failed.ProcessedAt)

	require.NoError(t, events.RecordOutcome(ctx, event, nil))
	var processed models.WebhookEvent
	require.NoError(t, db.First(&processed, "id = ?", event.ID).Error)
	assert.Equal(t, "processed", processed.Status)
	assert.Equal(t, "", processed.HandlerError, "a fixed event still carries the message from the attempt that failed")
}

// A replay returns the count the database holds, not the one the caller read:
// two replays of one event each add 1.
func TestRecordReplayReturnsTheStoredCount(t *testing.T) {
	_, events := webhookEventDB(t)
	ctx := context.Background()
	event := storedEvent(t, events, "evt_replay", "failed", 0)

	first, err := events.RecordReplay(ctx, event, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, first)

	// The caller still holds the event as it was read, with no retries on it.
	second, err := events.RecordReplay(ctx, event, errors.New("still broken"))
	require.NoError(t, err)
	assert.Equal(t, 2, second, "the second replay reported the first one's count")
}

// The row a redelivery collided with is found by provider and external id, and
// one provider's id does not match another's.
func TestByExternalIDIsScopedToItsProvider(t *testing.T) {
	_, events := webhookEventDB(t)
	ctx := context.Background()
	storedEvent(t, events, "evt_shared", "processed", 0)

	found, err := events.ByExternalID(ctx, "stripe", "evt_shared")
	require.NoError(t, err)
	assert.Equal(t, "processed", found.Status)

	_, err = events.ByExternalID(ctx, "paystack", "evt_shared")
	assert.Error(t, err, "another provider's event with the same id was returned")
}
