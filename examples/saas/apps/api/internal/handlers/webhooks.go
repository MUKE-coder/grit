package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"saas/apps/api/internal/models"
	"saas/apps/api/internal/paginate"
	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
	"saas/apps/api/internal/webhooks"
)

// WebhookHandler is the universal entry point for inbound webhooks.
// One handler, one route — POST /webhooks/:provider routes by the
// provider path param and dispatches to whatever was registered.
type WebhookHandler struct {
	DB *gorm.DB
}

func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	return &WebhookHandler{DB: db}
}

// Receive is mounted at POST /webhooks/:provider. It:
//
//  1. Looks up the provider in the registry (404 if unknown).
//  2. Reads the raw body + collects headers.
//  3. Calls Provider.Verify — 401 on signature mismatch.
//  4. Calls Provider.Extract to get (event_type, external_id).
//  5. Inserts a WebhookEvent (unique on provider+external_id). A second
//     delivery of an event already processed is skipped; one that failed, or
//     that a dead process left pending, is claimed and run again.
//  6. Calls webhooks.Dispatch in the request context.
//  7. Updates status=processed or status=failed with HandlerError.
//
// Always returns 200 to the provider on a verified+stored event so
// they don't retry forever — handler failures are surfaced via the
// admin replay endpoint.
func (h *WebhookHandler) Receive(c *gin.Context) {
	providerName := c.Param("provider")
	provider, ok := webhooks.LookupProvider(providerName)
	if !ok {
		respond.Fail(c, respond.CodeUnknownProvider, "no webhook provider registered for "+providerName)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respond.Fail(c, respond.CodeReadBodyFailed, err.Error())
		return
	}

	headers := flattenHeaders(c.Request.Header)
	secret := os.Getenv(provider.SecretEnv)
	if err := provider.Verify(secret, body, headers); err != nil {
		respond.Fail(c, respond.CodeInvalidSignature, err.Error())
		return
	}

	eventType, externalID, err := provider.Extract(body, headers)
	if err != nil {
		respond.Fail(c, respond.CodeExtractFailed, err.Error())
		return
	}

	// A provider that supplies no event id cannot be deduplicated, and every
	// such delivery is a distinct row. NULL says that; "" would make two
	// unrelated anonymous events collide on the unique index and drop the
	// second one as a duplicate it is not.
	var externalRef *string
	if externalID != "" {
		externalRef = &externalID
	}

	event := models.WebhookEvent{
		Provider:   providerName,
		EventType:  eventType,
		ExternalID: externalRef,
		Payload:    datatypes.JSON(body),
		Status:     "pending",
	}
	if err := h.events().Store(c.Request.Context(), &event); err != nil {
		// Already stored. Whether that is the end of it depends on what
		// happened the first time, which redelivered decides.
		if webhooks.IsDuplicateError(err) {
			h.redelivered(c, providerName, externalID)
			return
		}
		respond.ServerError(c, "PERSIST_FAILED", err, "Internal server error")
		return
	}

	h.runHandler(c, &event)
}

// events is the service every method below reads and writes through. Whether an
// event was stored once, claimed by this request and recorded as processed is a
// decision about data, and the claim is why: see
// internal/services/webhook_event.go.
func (h *WebhookHandler) events() *services.WebhookEventService {
	return &services.WebhookEventService{DB: h.DB}
}

// webhookInFlight is how long a pending event is assumed to be somebody's
// request still running. Past it, a redelivery may take the event over: the
// process that stored it is gone, and nothing else will ever run it.
const webhookInFlight = 5 * time.Minute

// redelivered answers a delivery of an event already in the table.
//
// Calling it a duplicate and stopping is right for an event that was
// processed, and wrong for the two cases that actually bring a provider back:
// an event whose handler failed, and an event still marked pending because the
// process died between storing it and running it. Both used to be answered
// "skipped", so the retry designed to recover exactly this was thrown away,
// and a subscription that granted nothing went on granting nothing until
// somebody thought to read the webhook_events table.
//
// So a redelivery re-runs the handler, but only when nobody else is on it. The
// row is claimed with a conditional update that exactly one caller wins, and a
// pending row is only claimed once it is too old to be a request in flight.
// Handlers have to be idempotent to survive a provider's retries at all, and
// this is the case they are idempotent for.
func (h *WebhookHandler) redelivered(c *gin.Context, provider, externalID string) {
	ctx := c.Request.Context()

	event, err := h.events().ByExternalID(ctx, provider, externalID)
	if err != nil {
		// The unique index just said this row exists, so failing to read it is
		// a database problem, not a duplicate. Nothing the provider can do
		// about it, so take the delivery and say so in the log.
		log.Printf("webhooks: %s event %s collided with a row that will not load: %v", provider, externalID, err)
		c.JSON(http.StatusOK, gin.H{"status": "skipped", "reason": "duplicate"})
		return
	}

	if event.Status == "processed" {
		c.JSON(http.StatusOK, gin.H{"status": "skipped", "reason": "duplicate", "id": event.ID})
		return
	}

	claimed, err := h.events().Claim(ctx, event, webhookInFlight)
	if err != nil {
		respond.ServerError(c, "PERSIST_FAILED", err, "Internal server error")
		return
	}
	if !claimed {
		// Another request is running this one right now. Two handlers for one
		// event is the thing worth avoiding; the provider will come back.
		c.JSON(http.StatusOK, gin.H{"status": "skipped", "reason": "in flight", "id": event.ID})
		return
	}

	log.Printf("webhooks: %s event %s came back as %s, running it again", provider, externalID, event.Status)
	h.runHandler(c, event)
}

// runHandler dispatches a stored event and records what happened.
//
// Dispatch runs in the request context so handlers can attach DB timeouts and
// cancellation. A handler failure is recorded and never bubbles up: the
// provider has a 200 the moment the event is stored, and a failure is
// recovered from the admin replay endpoint or from the provider's own
// redelivery.
func (h *WebhookHandler) runHandler(c *gin.Context, event *models.WebhookEvent) {
	if dispatchErr := webhooks.Dispatch(c.Request.Context(), event); dispatchErr != nil {
		if err := h.events().RecordOutcome(c.Request.Context(), event, dispatchErr); err != nil {
			log.Printf("webhooks: recording that event %s failed: %v", event.ID, err)
		}
		c.JSON(http.StatusOK, gin.H{"status": "received", "id": event.ID, "handler": "failed"})
		return
	}
	if err := h.events().RecordOutcome(c.Request.Context(), event, nil); err != nil {
		log.Printf("webhooks: recording that event %s was processed: %v", event.ID, err)
	}
	c.JSON(http.StatusOK, gin.H{"status": "processed", "id": event.ID})
}

// List returns the recent webhook events with the standard paginate envelope.
//
//	GET /api/admin/webhooks?provider=stripe&status=failed
func (h *WebhookHandler) List(c *gin.Context) {
	params := paginate.Bind(c).
		With("provider", c.Query("provider")).
		With("status", c.Query("status"))

	res, err := h.events().List(c.Request.Context(), params)
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Replay re-runs the handler for an existing webhook event. Used to
// recover from a transient handler failure or a deploy that fixed a
// bug. Increments retry_count + records the new outcome.
//
//	POST /api/admin/webhooks/:id/replay
func (h *WebhookHandler) Replay(c *gin.Context) {
	event, err := h.events().ByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.Fail(c, respond.CodeNotFound, "webhook event not found")
			return
		}
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}

	dispatchErr := webhooks.Dispatch(c.Request.Context(), event)
	// The count comes back from the write: two replays of one event each add 1
	// rather than both reading the same number, so the one this request holds is
	// stale the moment the update lands. See RecordReplay.
	retries, err := h.events().RecordReplay(c.Request.Context(), event, dispatchErr)
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	if dispatchErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":        "failed",
			"handler_error": dispatchErr.Error(),
			"retry_count":   retries,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "processed", "retry_count": retries})
}

// flattenHeaders turns http.Header (multi-value) into a single-value
// map for the Verify / Extract callbacks. Keeps the framework API
// simple — nobody needs the multi-value form for webhook signing.
func flattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}
