package mailtest

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFakeRecordsAndAsserts(t *testing.T) {
	ctx := context.Background()
	fake := New()
	mailer := fake.Mailer()

	if err := mailer.SendRaw(ctx, "Ada <ada@example.com>", "Reset your password", "<p>hi</p>"); err != nil {
		t.Fatalf("sending: %v", err)
	}
	m := fake.AssertSent(t, "ADA@example.com", "Reset")
	if m.From != "test@example.com" || m.HTML != "<p>hi</p>" {
		t.Errorf("recorded message = %+v", m)
	}
	if mailer.Driver() != "fake" {
		t.Errorf("driver = %q", mailer.Driver())
	}

	go func() {
		if err := mailer.SendRaw(ctx, "bob@example.com", "Confirm your email", "<p>hi</p>"); err != nil {
			t.Errorf("sending from a goroutine: %v", err)
		}
	}()
	fake.AssertSentWithin(t, 2*time.Second, "bob@example.com", "Confirm")

	if _, ok := fake.Find("carol@example.com", ""); ok {
		t.Error("found a message to an address nothing was sent to")
	}

	fake.Reset()
	fake.AssertNothingSent(t)

	fake.FailWith(errors.New("provider down"))
	if err := mailer.SendRaw(ctx, "ada@example.com", "Hi", "<p>hi</p>"); err == nil {
		t.Error("FailWith did not make Send fail")
	}
	fake.AssertNothingSent(t)
}
