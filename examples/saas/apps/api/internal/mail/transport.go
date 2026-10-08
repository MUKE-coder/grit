package mail

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	netmail "net/mail"
	"path/filepath"
	"strings"
	"time"
)

// Attachment is a file sent with a message. ContentType may be left empty: it
// is then guessed from the file name.
type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
}

// Message is one email, whichever transport carries it. From may be left
// empty to send from the Mailer's MAIL_FROM.
type Message struct {
	From        string
	ReplyTo     string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	HTML        string
	Text        string
	Attachments []Attachment
	Headers     map[string]string
}

// Transport delivers a message. Every driver implements it, so the code that
// sends mail never knows which provider is behind it.
type Transport interface {
	Send(ctx context.Context, m *Message) error
	Name() string
}

// Validate reports a message no transport should be handed: no sender, no
// recipient, no body, an address that does not parse, or a line break in a
// header, which would let a value taken from a form add headers of its own.
func (m *Message) Validate() error {
	if m.From == "" {
		return errors.New("mail: the message has no From address")
	}
	if len(m.To)+len(m.Cc)+len(m.Bcc) == 0 {
		return errors.New("mail: the message has no recipients")
	}
	if m.HTML == "" && m.Text == "" {
		return errors.New("mail: the message has no body")
	}
	fields := []string{m.From, m.ReplyTo, m.Subject}
	for k, v := range m.Headers {
		fields = append(fields, k, v)
	}
	for _, a := range m.Attachments {
		if a.Filename == "" {
			return errors.New("mail: an attachment has no file name")
		}
		fields = append(fields, a.Filename, a.ContentType)
	}
	addresses := append(append(append([]string{m.From}, m.To...), m.Cc...), m.Bcc...)
	if m.ReplyTo != "" {
		addresses = append(addresses, m.ReplyTo)
	}
	for _, f := range append(fields, addresses...) {
		if strings.ContainsAny(f, "\r\n") {
			return errors.New("mail: a header value contains a line break")
		}
	}
	for _, a := range addresses {
		if _, err := parseAddress(a); err != nil {
			return err
		}
	}
	return nil
}

func parseAddress(s string) (*netmail.Address, error) {
	a, err := netmail.ParseAddress(s)
	if err != nil {
		return nil, fmt.Errorf("mail: %q is not an email address: %w", s, err)
	}
	return a, nil
}

// addressOf returns the bare address of "Name <addr>" or "addr".
func addressOf(s string) (string, error) {
	a, err := parseAddress(s)
	if err != nil {
		return "", err
	}
	return a.Address, nil
}

func addressesOf(list []string) ([]string, error) {
	out := make([]string, 0, len(list))
	for _, s := range list {
		a, err := addressOf(s)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func attachmentType(a Attachment) string {
	if a.ContentType != "" {
		return a.ContentType
	}
	if t := mime.TypeByExtension(filepath.Ext(a.Filename)); t != "" {
		return t
	}
	return "application/octet-stream"
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("mail: reading random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// ProviderError is a request a provider refused (4xx) or failed (5xx).
type ProviderError struct {
	Driver string
	Status int
	Body   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("%s: the provider answered %d: %s", e.Driver, e.Status, e.Body)
}

// defaultHTTPClient is shared by the HTTP transports. The timeout covers the
// whole exchange, so a provider that stops answering cannot hold a send open.
var defaultHTTPClient = &http.Client{Timeout: 15 * time.Second}

// retryDelay is the pause before the one retry. A variable so tests need not wait.
var retryDelay = 500 * time.Millisecond

// postWithRetry sends the request build returns, and once more when the first
// try hit a network error or a 5xx. A 4xx is the provider saying the message
// itself is wrong, and sending it again would get the same answer.
func postWithRetry(ctx context.Context, client *http.Client, driver string, build func() (*http.Request, error)) error {
	if client == nil {
		client = defaultHTTPClient
	}
	var last error
	for attempt := 1; attempt <= 2; attempt++ {
		if attempt == 2 {
			timer := time.NewTimer(retryDelay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return last
			case <-timer.C:
			}
		}
		req, err := build()
		if err != nil {
			return fmt.Errorf("%s: building the request: %w", driver, err)
		}
		resp, err := client.Do(req)
		if err != nil {
			last = fmt.Errorf("%s: sending: %w", driver, err)
			if ctx.Err() != nil {
				return last
			}
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 2048))
		if closeErr := resp.Body.Close(); closeErr != nil && readErr == nil {
			readErr = closeErr
		}
		if resp.StatusCode < 300 {
			return nil
		}
		text := strings.TrimSpace(string(body))
		if readErr != nil {
			text = "reading the response: " + readErr.Error()
		}
		perr := &ProviderError{Driver: driver, Status: resp.StatusCode, Body: text}
		if resp.StatusCode < 500 {
			return perr
		}
		last = perr
	}
	return last
}
