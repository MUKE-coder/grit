package mail

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"sort"
	"strings"
)

// MailgunTransport sends through Mailgun's messages API (MAIL_MAILER=mailgun):
// https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun/messages
type MailgunTransport struct {
	Domain string
	Secret string
	// Endpoint is the API host: api.mailgun.net (the default) or
	// api.eu.mailgun.net for a domain in the EU region. A full URL is used as
	// it is.
	Endpoint string
	Client   *http.Client
}

// NewMailgun returns a Mailgun transport for a sending domain.
func NewMailgun(domain, secret, endpoint string) *MailgunTransport {
	return &MailgunTransport{Domain: domain, Secret: secret, Endpoint: endpoint}
}

// Name implements Transport.
func (t *MailgunTransport) Name() string { return "mailgun" }

func (t *MailgunTransport) url() string {
	base := strings.TrimRight(strings.TrimSpace(t.Endpoint), "/")
	if base == "" {
		base = "api.mailgun.net"
	}
	if !strings.Contains(base, "://") {
		base = "https://" + base
	}
	return base + "/v3/" + url.PathEscape(t.Domain) + "/messages"
}

// Send implements Transport.
func (t *MailgunTransport) Send(ctx context.Context, m *Message) error {
	fields := [][2]string{{"from", m.From}}
	for _, list := range []struct {
		name  string
		addrs []string
	}{{"to", m.To}, {"cc", m.Cc}, {"bcc", m.Bcc}} {
		for _, a := range list.addrs {
			fields = append(fields, [2]string{list.name, a})
		}
	}
	fields = append(fields, [2]string{"subject", m.Subject})
	if m.Text != "" {
		fields = append(fields, [2]string{"text", m.Text})
	}
	if m.HTML != "" {
		fields = append(fields, [2]string{"html", m.HTML})
	}
	if m.ReplyTo != "" {
		fields = append(fields, [2]string{"h:Reply-To", m.ReplyTo})
	}
	names := make([]string, 0, len(m.Headers))
	for name := range m.Headers {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fields = append(fields, [2]string{"h:" + name, m.Headers[name]})
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for _, f := range fields {
		if err := w.WriteField(f[0], f[1]); err != nil {
			return fmt.Errorf("mailgun: encoding %s: %w", f[0], err)
		}
	}
	for _, a := range m.Attachments {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", mime.FormatMediaType("form-data", map[string]string{"name": "attachment", "filename": a.Filename}))
		h.Set("Content-Type", attachmentType(a))
		part, err := w.CreatePart(h)
		if err != nil {
			return fmt.Errorf("mailgun: encoding %s: %w", a.Filename, err)
		}
		if _, err := part.Write(a.Content); err != nil {
			return fmt.Errorf("mailgun: encoding %s: %w", a.Filename, err)
		}
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mailgun: encoding the message: %w", err)
	}

	body := buf.Bytes()
	contentType := w.FormDataContentType()
	return postWithRetry(ctx, t.Client, "mailgun", func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.url(), bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth("api", t.Secret)
		req.Header.Set("Content-Type", contentType)
		return req, nil
	})
}
