// Package sms is the seam an SMS provider plugs into.
//
// There is no default implementation, deliberately. Every provider is a paid
// account with its own credentials, and the right one depends on where your
// users are: Twilio is not the sensible choice in Kampala, and Africa's Talking
// is not the sensible choice in Berlin. Baking one in would make everybody
// carry a dependency most of them cannot use.
//
// So the framework defines the interface and the features that need SMS ask
// whether anything is registered. Nothing that needs a text message is offered
// in the UI until something is.
//
// Wire one in main.go:
//
//	sms.Register(sms.SenderFunc(func(ctx context.Context, to, body string) error {
//	    // call your provider here
//	    return nil
//	}))
package sms

import (
	"context"
	"errors"
	"sync"
)

// ErrNoProvider is returned by Send when nothing has been registered. Callers
// should check Configured() first and not offer the feature at all, so this is
// the backstop rather than the expected path.
var ErrNoProvider = errors.New("no SMS provider is configured")

// Sender delivers one text message.
type Sender interface {
	Send(ctx context.Context, to, body string) error
}

// SenderFunc adapts a plain function to Sender.
type SenderFunc func(ctx context.Context, to, body string) error

func (f SenderFunc) Send(ctx context.Context, to, body string) error { return f(ctx, to, body) }

var (
	mu      sync.RWMutex
	current Sender
)

// Register installs the provider. Call it once, from main, before serving.
func Register(s Sender) {
	mu.Lock()
	defer mu.Unlock()
	current = s
}

// Configured reports whether a provider is installed.
//
// The security overview endpoint reports this to the client, so the admin can
// hide phone recovery rather than offering a button that cannot work. A
// disabled control with no explanation is worse than no control.
func Configured() bool {
	mu.RLock()
	defer mu.RUnlock()
	return current != nil
}

// Send delivers a message through the registered provider.
func Send(ctx context.Context, to, body string) error {
	mu.RLock()
	s := current
	mu.RUnlock()
	if s == nil {
		return ErrNoProvider
	}
	return s.Send(ctx, to, body)
}

// Reset clears the provider. Tests only.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	current = nil
}
