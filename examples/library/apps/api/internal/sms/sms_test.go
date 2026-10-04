package sms

import (
	"context"
	"errors"
	"testing"
)

func TestUnconfiguredReportsItselfAndRefuses(t *testing.T) {
	Reset()
	if Configured() {
		t.Fatal("nothing is registered, so Configured must be false")
	}
	// The backstop: a caller that skipped the Configured check gets a named
	// error rather than a nil dereference.
	if err := Send(context.Background(), "+256700000000", "hi"); !errors.Is(err, ErrNoProvider) {
		t.Errorf("expected ErrNoProvider, got %v", err)
	}
}

func TestRegisteredProviderReceivesTheMessage(t *testing.T) {
	Reset()
	defer Reset()

	var gotTo, gotBody string
	Register(SenderFunc(func(_ context.Context, to, body string) error {
		gotTo, gotBody = to, body
		return nil
	}))

	if !Configured() {
		t.Fatal("Configured must be true once a provider is registered")
	}
	if err := Send(context.Background(), "+256700000000", "code 123456"); err != nil {
		t.Fatal(err)
	}
	if gotTo != "+256700000000" || gotBody != "code 123456" {
		t.Errorf("provider got %q / %q", gotTo, gotBody)
	}
}

func TestProviderErrorIsReturned(t *testing.T) {
	Reset()
	defer Reset()

	boom := errors.New("provider rejected the number")
	Register(SenderFunc(func(_ context.Context, _, _ string) error { return boom }))

	if err := Send(context.Background(), "x", "y"); !errors.Is(err, boom) {
		t.Errorf("the provider's error should reach the caller, got %v", err)
	}
}
