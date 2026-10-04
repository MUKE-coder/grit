package services

import (
	"gorm.io/gorm"

	"library/apps/api/internal/events"
	"library/apps/api/internal/realtime"
)

// WebhookDispatcher is the shape the webhooks plugin's DispatchWebhook has.
// Declared here as a function type so internal/events does not import a
// package that may not be installed: a project without the webhooks plugin
// simply never sets it.
type WebhookDispatcher func(eventType string, payload interface{}) error

// RegisterEventSubscribers wires the subscribers every project gets.
//
// Called once from routes.Setup, after Init. hub and dispatch may be nil, for
// a project that has no realtime or has not installed the webhooks plugin.
func RegisterEventSubscribers(db *gorm.DB, hub *realtime.Hub, dispatch WebhookDispatcher) {
	registerAudit(db)
	registerRealtime(hub)
	registerWebhooks(dispatch)
}

// registerAudit records every event in the activity feed.
//
// Sync, because this is the one subscriber that needs the request context: the
// feed records IP and user agent. The row is built here and written by a
// batching writer (services/activity_writer.go), so the request does not wait
// on the insert.
func registerAudit(db *gorm.DB) {
	if db == nil {
		return
	}
	events.On("*", events.Sync, "audit", func(e events.Event) error {
		if e.C == nil {
			return nil // emitted outside a request; nothing to attribute it to
		}
		switch verb(e.Name) {
		case "created":
			LogCreate(db, e.C, e.Entity, e.Label, e.ID, e.Detail)
		case "deleted":
			LogDelete(db, e.C, e.Entity, e.Label, e.ID)
		default:
			// Everything else reads as a change to the record, including
			// workflow transitions, which is what you want in a feed: "Maria
			// sent invoice INV-0007" rather than a row saying "updated".
			LogUpdate(db, e.C, e.Entity, e.Label, e.ID, e.Detail)
		}
		return nil
	})
}

// RealtimeAudience decides which users see a resource event.
//
// The default is the actor and nobody else, so a user's other devices stay in
// sync and no one learns about a record they may not be allowed to read.
//
// This is deliberately narrow. The payload carries e.Label, which is the row's
// human-readable title: a conversation subject, a document name, a customer.
// Broadcasting that to every connected session hands every authenticated user
// a live feed of every title in the system, across tenants, whatever the REST
// layer permits. There is no per-row authorization in the event bus to lean
// on, so the only safe default is the one person already known to have seen
// the row: whoever just wrote it.
//
// Widen it deliberately when the app knows who is allowed to look. A chat app
// sends to the conversation's participants; an admin dashboard sends to staff:
//
//	services.RealtimeAudience = func(e events.Event) []string {
//	    if e.Resource == "messages" {
//	        return participantIDs(db, e)
//	    }
//	    return staffIDs(db)
//	}
//
// Returning nil drops the event, which is what an unauthenticated or
// system-initiated write does by default.
var RealtimeAudience = func(e events.Event) []string {
	if e.Actor == "" {
		return nil
	}
	return []string{e.Actor}
}

// RealtimeChannels sends a resource event to channels as well as to users.
//
// Nil by default, so resource events reach only RealtimeAudience. Set it when a
// page subscribes to a record rather than waiting for events addressed to its
// user:
//
//	services.RealtimeChannels = func(e events.Event) []string {
//	    if e.Resource == "invoices" {
//	        return []string{"private-invoices." + e.ID}
//	    }
//	    return nil
//	}
//
// Every subscriber of a returned channel receives the event, label included, so
// return a private- or presence- channel whose authorizer (realtime.Channel)
// admits only users allowed to read the record. A public- channel reaches any
// connection that asks for it.
var RealtimeChannels func(e events.Event) []string

// registerRealtime pushes resource events to the users RealtimeAudience picks.
//
// Async: a websocket write to a client on a bad connection must not slow the
// request that caused the event.
func registerRealtime(hub *realtime.Hub) {
	if hub == nil {
		return
	}
	events.On("*", events.Async, "realtime", func(e events.Event) error {
		evt := realtime.Event{
			Type: e.Name,
			Payload: map[string]interface{}{
				"resource": e.Resource,
				"id":       e.ID,
				"label":    e.Label,
				"actor":    e.Actor,
				"at":       e.At,
			},
		}
		if audience := RealtimeAudience(e); len(audience) > 0 {
			hub.SendToUsers(audience, evt)
		}
		if RealtimeChannels != nil {
			for _, channel := range RealtimeChannels(e) {
				hub.Publish(channel, evt)
			}
		}
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
	events.On("*", events.Async, "webhooks", func(e events.Event) error {
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
