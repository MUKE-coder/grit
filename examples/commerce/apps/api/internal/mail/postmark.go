package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// PostmarkEndpoint is Postmark's single-email API:
// https://postmarkapp.com/developer/api/email-api
const PostmarkEndpoint = "https://api.postmarkapp.com/email"

// PostmarkTransport sends through Postmark (MAIL_MAILER=postmark).
type PostmarkTransport struct {
	Token string
	// MessageStream is the stream to send on. Empty means Postmark's default,
	// "outbound".
	MessageStream string
	Endpoint      string
	Client        *http.Client
}

// NewPostmark returns a Postmark transport for a server token.
func NewPostmark(token, messageStream string) *PostmarkTransport {
	return &PostmarkTransport{Token: token, MessageStream: messageStream, Endpoint: PostmarkEndpoint}
}

// Name implements Transport.
func (t *PostmarkTransport) Name() string { return "postmark" }

// Send implements Transport.
func (t *PostmarkTransport) Send(ctx context.Context, m *Message) error {
	payload := map[string]interface{}{
		"From":    m.From,
		"Subject": m.Subject,
	}
	if len(m.To) > 0 {
		payload["To"] = strings.Join(m.To, ",")
	}
	if len(m.Cc) > 0 {
		payload["Cc"] = strings.Join(m.Cc, ",")
	}
	if len(m.Bcc) > 0 {
		payload["Bcc"] = strings.Join(m.Bcc, ",")
	}
	if m.ReplyTo != "" {
		payload["ReplyTo"] = m.ReplyTo
	}
	if m.HTML != "" {
		payload["HtmlBody"] = m.HTML
	}
	if m.Text != "" {
		payload["TextBody"] = m.Text
	}
	if t.MessageStream != "" {
		payload["MessageStream"] = t.MessageStream
	}
	if len(m.Headers) > 0 {
		names := make([]string, 0, len(m.Headers))
		for name := range m.Headers {
			names = append(names, name)
		}
		sort.Strings(names)
		headers := make([]map[string]string, 0, len(names))
		for _, name := range names {
			headers = append(headers, map[string]string{"Name": name, "Value": m.Headers[name]})
		}
		payload["Headers"] = headers
	}
	if len(m.Attachments) > 0 {
		list := make([]map[string]string, 0, len(m.Attachments))
		for _, a := range m.Attachments {
			list = append(list, map[string]string{
				"Name":        a.Filename,
				"Content":     base64.StdEncoding.EncodeToString(a.Content),
				"ContentType": attachmentType(a),
			})
		}
		payload["Attachments"] = list
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("postmark: encoding the message: %w", err)
	}
	return postWithRetry(ctx, t.Client, "postmark", func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.Endpoint, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Postmark-Server-Token", t.Token)
		return req, nil
	})
}
