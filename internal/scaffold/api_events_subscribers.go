package scaffold

// apiEventsSubscribersGo emits internal/events/subscribers.go: the default
// wiring that turns three built-but-unfired systems into subscribers.
//
// Audit was already called directly by every generated handler, so it keeps
// working exactly as before, just from the other side of the bus. Webhooks and
// realtime were never fired by generated code at all: dispatching a webhook
// meant hand-writing the call in every handler for every operation, and
// forgetting one produced no error, just a subscription that never heard
// anything.
func apiEventsSubscribersGo() string {
	return `package events

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/realtime"
	"{{MODULE}}/internal/services"
)

// WebhookDispatcher is the shape the webhooks plugin's DispatchWebhook has.
// Declared here as a function type so internal/events does not import a
// package that may not be installed: a project without the webhooks plugin
// simply never sets it.
type WebhookDispatcher func(eventType string, payload interface{}) error

// RegisterDefaults wires the subscribers every project gets.
//
// Called once from routes.Setup, after Init. hub and dispatch may be nil, for
// a project that has no realtime or has not installed the webhooks plugin.
func RegisterDefaults(db *gorm.DB, hub *realtime.Hub, dispatch WebhookDispatcher) {
	registerAudit(db)
	registerRealtime(hub)
	registerWebhooks(dispatch)
}

// registerAudit records every event in the activity feed.
//
// Sync, and deliberately so. The activity row should exist before the caller
// is told the write succeeded, and this is the one subscriber that legitimately
// needs the request context: the feed records IP and user agent.
func registerAudit(db *gorm.DB) {
	if db == nil {
		return
	}
	On("*", Sync, "audit", func(e Event) error {
		if e.C == nil {
			return nil // emitted outside a request; nothing to attribute it to
		}
		switch verb(e.Name) {
		case "created":
			services.LogCreate(db, e.C, e.Entity, e.Label, e.ID, e.Detail)
		case "deleted":
			services.LogDelete(db, e.C, e.Entity, e.Label, e.ID)
		default:
			// Everything else reads as a change to the record, including
			// workflow transitions, which is what you want in a feed: "Maria
			// sent invoice INV-0007" rather than a row saying "updated".
			services.LogUpdate(db, e.C, e.Entity, e.Label, e.ID, e.Detail)
		}
		return nil
	})
}

// registerRealtime pushes every event to connected clients.
//
// Async: a websocket write to a client on a bad connection must not slow the
// request that caused the event.
//
// Broadcast rather than per-user, matching what the hub can do today. Once
// rooms exist this becomes a send to the room watching the record, which is
// the point at which presence and live editing become possible.
func registerRealtime(hub *realtime.Hub) {
	if hub == nil {
		return
	}
	On("*", Async, "realtime", func(e Event) error {
		hub.Broadcast(realtime.Event{
			Type: e.Name,
			Payload: map[string]interface{}{
				"resource": e.Resource,
				"id":       e.ID,
				"label":    e.Label,
				"actor":    e.Actor,
				"at":       e.At,
			},
		})
		return nil
	})
}

// registerWebhooks fans every event out to matching subscriptions.
//
// Async, because a subscriber's endpoint is somebody else's server and may be
// slow or down. The webhooks plugin owns retries and the delivery log from
// there; the bus only has to hand it over.
func registerWebhooks(dispatch WebhookDispatcher) {
	if dispatch == nil {
		return
	}
	On("*", Async, "webhooks", func(e Event) error {
		return dispatch(e.Name, map[string]interface{}{
			"resource": e.Resource,
			"id":       e.ID,
			"label":    e.Label,
			"actor":    e.Actor,
			"at":       e.At,
			"data":     e.After,
		})
	})
}

// verb pulls the trailing segment off an event name: "invoices.created" gives
// "created". A name with no dot is its own verb.
func verb(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
	}
	return name
}

// Emitted is a small helper for the common case: a CRUD event on a resource.
//
// Generated handlers call this rather than building an Event literal, so the
// name and resource cannot drift apart and every resource emits the same
// shape.
func Emitted(c *gin.Context, resource, entity, action, id, label, detail string, before, after interface{}) {
	Emit(c, Event{
		Name:     resource + "." + action,
		Resource: resource,
		Entity:   entity,
		ID:       id,
		Label:    label,
		Detail:   detail,
		Before:   before,
		After:    after,
	})
}
`
}
