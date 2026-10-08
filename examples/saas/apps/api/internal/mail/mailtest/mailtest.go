// Package mailtest is a fake mail transport for tests: it keeps every message
// instead of sending it, and asserts on what was sent.
//
//	fake := mailtest.New()
//	handler := &handlers.AuthHandler{Mailer: fake.Mailer(), ...}
//	// ... exercise the handler ...
//	fake.AssertSent(t, "ada@example.com", "Reset your password")
package mailtest

import (
	"context"
	"net/mail"
	"strings"
	"sync"
	"testing"
	"time"

	appmail "saas/apps/api/internal/mail"
)

// Fake is a mail transport that records messages.
type Fake struct {
	mu   sync.Mutex
	sent []appmail.Message
	err  error
}

// New returns an empty Fake.
func New() *Fake {
	return &Fake{}
}

// Name implements mail.Transport.
func (f *Fake) Name() string { return "fake" }

// Send implements mail.Transport. It records a copy of the message, or returns
// the error set by FailWith.
func (f *Fake) Send(_ context.Context, m *appmail.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	msg := *m
	msg.To = append([]string(nil), m.To...)
	msg.Cc = append([]string(nil), m.Cc...)
	msg.Bcc = append([]string(nil), m.Bcc...)
	msg.Attachments = append([]appmail.Attachment(nil), m.Attachments...)
	if m.Headers != nil {
		msg.Headers = make(map[string]string, len(m.Headers))
		for k, v := range m.Headers {
			msg.Headers[k] = v
		}
	}
	f.sent = append(f.sent, msg)
	return nil
}

// Mailer returns a mail.Mailer that sends through this fake, from
// test@example.com.
func (f *Fake) Mailer() *appmail.Mailer {
	return appmail.NewWithTransport(f, "test@example.com")
}

// FailWith makes every later Send return err. FailWith(nil) undoes it.
func (f *Fake) FailWith(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// Sent returns every message recorded so far, oldest first.
func (f *Fake) Sent() []appmail.Message {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]appmail.Message(nil), f.sent...)
}

// Reset forgets every recorded message.
func (f *Fake) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sent = nil
}

// Find returns the first message to address (To, Cc or Bcc, compared without
// case or display name) whose subject contains subjectContains.
func (f *Fake) Find(to, subjectContains string) (appmail.Message, bool) {
	for _, m := range f.Sent() {
		if strings.Contains(m.Subject, subjectContains) && addressedTo(m, to) {
			return m, true
		}
	}
	return appmail.Message{}, false
}

// AssertSent fails the test unless a message to address with a subject
// containing subjectContains was sent, and returns it.
func (f *Fake) AssertSent(t testing.TB, to, subjectContains string) appmail.Message {
	t.Helper()
	m, ok := f.Find(to, subjectContains)
	if !ok {
		t.Fatalf("mailtest: no message to %s with a subject containing %q; sent: %s", to, subjectContains, f.summary())
	}
	return m
}

// AssertSentWithin is AssertSent for mail sent from a goroutine: it waits up
// to timeout for the message to arrive.
func (f *Fake) AssertSentWithin(t testing.TB, timeout time.Duration, to, subjectContains string) appmail.Message {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		if m, ok := f.Find(to, subjectContains); ok {
			return m
		}
		if time.Now().After(deadline) {
			t.Fatalf("mailtest: no message to %s with a subject containing %q within %s; sent: %s", to, subjectContains, timeout, f.summary())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// AssertNothingSent fails the test if any message was sent.
func (f *Fake) AssertNothingSent(t testing.TB) {
	t.Helper()
	if sent := f.Sent(); len(sent) > 0 {
		t.Fatalf("mailtest: expected no mail, got %d: %s", len(sent), f.summary())
	}
}

func (f *Fake) summary() string {
	sent := f.Sent()
	if len(sent) == 0 {
		return "nothing"
	}
	parts := make([]string, 0, len(sent))
	for _, m := range sent {
		parts = append(parts, strings.Join(m.To, ",")+" "+m.Subject)
	}
	return strings.Join(parts, "; ")
}

func addressedTo(m appmail.Message, want string) bool {
	want = bareAddress(want)
	for _, list := range [][]string{m.To, m.Cc, m.Bcc} {
		for _, addr := range list {
			if strings.EqualFold(bareAddress(addr), want) {
				return true
			}
		}
	}
	return false
}

func bareAddress(s string) string {
	if a, err := mail.ParseAddress(s); err == nil {
		return a.Address
	}
	return strings.TrimSpace(s)
}
