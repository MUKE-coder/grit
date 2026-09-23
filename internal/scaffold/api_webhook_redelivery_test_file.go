package scaffold

// apiWebhookRedeliveryTestGo emits internal/handlers/webhooks_redelivery_test.go.
//
// A provider that retries is the only thing standing between a handler that
// failed and a customer who paid for nothing. The receiver used to answer
// every repeat delivery "skipped: duplicate", including the repeats of an
// event whose handler had failed and of an event a dead process left pending,
// so the retry was thrown away and the failure was permanent, recorded only in
// a table nobody reads.
//
// These tests hold the three cases apart: processed stays skipped, failed runs
// again, and an event still in flight is not run twice.
func apiWebhookRedeliveryTestGo() string {
	return `package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/webhooks"
)

// A provider with no signature to check, so these tests are about the
// delivery and not about HMAC.
func registerTestProvider(t *testing.T, name string) {
	t.Helper()
	webhooks.Register(name, webhooks.Provider{
		SecretEnv: "UNUSED_SECRET",
		Verify:    func(secret string, body []byte, headers map[string]string) error { return nil },
		Extract: func(body []byte, headers map[string]string) (string, string, error) {
			var payload struct {
				ID   string ` + "`json:\"id\"`" + `
				Type string ` + "`json:\"type\"`" + `
			}
			if err := json.Unmarshal(body, &payload); err != nil {
				return "", "", err
			}
			return payload.Type, payload.ID, nil
		},
	})
}

func redeliveryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, db.AutoMigrate(&models.WebhookEvent{}))
	webhooks.Setup(db)
	return db
}

func deliver(t *testing.T, db *gorm.DB, provider, id string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/webhooks/:provider", NewWebhookHandler(db).Receive)

	body, err := json.Marshal(map[string]string{"id": id, "type": "invoice.paid"})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/"+provider, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func answer(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	return out
}

// The ordinary case, and the one the old code got right: an event that was
// handled is not handled again.
func TestRedeliveryOfAProcessedEventIsSkipped(t *testing.T) {
	db := redeliveryDB(t)
	registerTestProvider(t, "processed_provider")

	runs := 0
	webhooks.On("processed_provider", "", func(ctx context.Context, e *models.WebhookEvent) error {
		runs++
		return nil
	})

	require.Equal(t, http.StatusOK, deliver(t, db, "processed_provider", "evt_ok").Code)
	body := answer(t, deliver(t, db, "processed_provider", "evt_ok"))

	assert.Equal(t, 1, runs, "the handler ran twice for one event")
	assert.Equal(t, "skipped", body["status"])
	assert.Equal(t, "duplicate", body["reason"])
}

// The case that cost a customer their purchase: the handler failed, the
// provider came back, and the retry was answered "duplicate" and dropped.
func TestRedeliveryOfAFailedEventRunsAgain(t *testing.T) {
	db := redeliveryDB(t)
	registerTestProvider(t, "failing_provider")

	runs := 0
	webhooks.On("failing_provider", "", func(ctx context.Context, e *models.WebhookEvent) error {
		runs++
		if runs == 1 {
			return assert.AnError
		}
		return nil
	})

	first := answer(t, deliver(t, db, "failing_provider", "evt_retry"))
	require.Equal(t, "failed", first["handler"], "the first delivery was supposed to fail")

	second := answer(t, deliver(t, db, "failing_provider", "evt_retry"))
	assert.Equal(t, 2, runs, "the retry was dropped; a failed webhook could only be fixed by hand")
	assert.Equal(t, "processed", second["status"])

	var stored models.WebhookEvent
	require.NoError(t, db.Where("external_id = ?", "evt_retry").First(&stored).Error)
	assert.Equal(t, "processed", stored.Status)
	assert.Equal(t, "", stored.HandlerError, "the error from the first attempt outlived it")
	assert.Equal(t, 1, stored.RetryCount, "the second attempt was not counted")

	// Still exactly one row: the retry recovered the event, it did not store
	// a second copy of it.
	var n int64
	db.Model(&models.WebhookEvent{}).Count(&n)
	assert.Equal(t, int64(1), n)
}

// A pending row means somebody's request is running the handler right now.
// Taking it over would run one event twice at once, which is the thing the
// unique index exists to prevent.
func TestRedeliveryWhileTheFirstIsStillRunningIsSkipped(t *testing.T) {
	db := redeliveryDB(t)
	registerTestProvider(t, "slow_provider")
	webhooks.On("slow_provider", "", func(ctx context.Context, e *models.WebhookEvent) error { return nil })

	require.NoError(t, db.Create(&models.WebhookEvent{
		Provider:   "slow_provider",
		EventType:  "invoice.paid",
		ExternalID: ptr("evt_in_flight"),
		Status:     "pending",
		CreatedAt:  time.Now(),
	}).Error)

	body := answer(t, deliver(t, db, "slow_provider", "evt_in_flight"))
	assert.Equal(t, "skipped", body["status"])
	assert.Equal(t, "in flight", body["reason"])
}

// And the reason the in-flight rule has a clock in it: a process that died
// between storing the event and running it leaves a pending row forever, and
// without this the provider's retries bounce off it until somebody notices.
func TestRedeliveryOfAnAbandonedEventRunsIt(t *testing.T) {
	db := redeliveryDB(t)
	registerTestProvider(t, "crashed_provider")

	runs := 0
	webhooks.On("crashed_provider", "", func(ctx context.Context, e *models.WebhookEvent) error {
		runs++
		return nil
	})

	require.NoError(t, db.Create(&models.WebhookEvent{
		Provider:   "crashed_provider",
		EventType:  "invoice.paid",
		ExternalID: ptr("evt_abandoned"),
		Status:     "pending",
		CreatedAt:  time.Now().Add(-time.Hour),
	}).Error)

	body := answer(t, deliver(t, db, "crashed_provider", "evt_abandoned"))
	assert.Equal(t, 1, runs, "an event nobody is working on was left pending forever")
	assert.Equal(t, "processed", body["status"])
}

func ptr(s string) *string { return &s }
`
}
