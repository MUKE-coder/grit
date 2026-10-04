package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendGridEndpoint is SendGrid's v3 Mail Send API:
// https://www.twilio.com/docs/sendgrid/api-reference/mail-send/mail-send
const SendGridEndpoint = "https://api.sendgrid.com/v3/mail/send"

// SendGridTransport sends through SendGrid (MAIL_MAILER=sendgrid).
type SendGridTransport struct {
	APIKey   string
	Endpoint string
	Client   *http.Client
}

// NewSendGrid returns a SendGrid transport for an API key.
func NewSendGrid(apiKey string) *SendGridTransport {
	return &SendGridTransport{APIKey: apiKey, Endpoint: SendGridEndpoint}
}

// Name implements Transport.
func (t *SendGridTransport) Name() string { return "sendgrid" }

func sendgridAddress(s string) (map[string]string, error) {
	a, err := parseAddress(s)
	if err != nil {
		return nil, err
	}
	out := map[string]string{"email": a.Address}
	if a.Name != "" {
		out["name"] = a.Name
	}
	return out, nil
}

func sendgridAddresses(list []string) ([]map[string]string, error) {
	out := make([]map[string]string, 0, len(list))
	for _, s := range list {
		a, err := sendgridAddress(s)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// Send implements Transport.
func (t *SendGridTransport) Send(ctx context.Context, m *Message) error {
	from, err := sendgridAddress(m.From)
	if err != nil {
		return err
	}
	personalization := map[string]interface{}{}
	for field, list := range map[string][]string{"to": m.To, "cc": m.Cc, "bcc": m.Bcc} {
		if len(list) == 0 {
			continue
		}
		addrs, err := sendgridAddresses(list)
		if err != nil {
			return err
		}
		personalization[field] = addrs
	}
	// SendGrid wants text/plain before text/html.
	content := []map[string]string{}
	if m.Text != "" {
		content = append(content, map[string]string{"type": "text/plain", "value": m.Text})
	}
	if m.HTML != "" {
		content = append(content, map[string]string{"type": "text/html", "value": m.HTML})
	}
	payload := map[string]interface{}{
		"personalizations": []interface{}{personalization},
		"from":             from,
		"subject":          m.Subject,
		"content":          content,
	}
	if m.ReplyTo != "" {
		replyTo, err := sendgridAddress(m.ReplyTo)
		if err != nil {
			return err
		}
		payload["reply_to"] = replyTo
	}
	if len(m.Headers) > 0 {
		payload["headers"] = m.Headers
	}
	if len(m.Attachments) > 0 {
		list := make([]map[string]string, 0, len(m.Attachments))
		for _, a := range m.Attachments {
			list = append(list, map[string]string{
				"content":     base64.StdEncoding.EncodeToString(a.Content),
				"filename":    a.Filename,
				"type":        attachmentType(a),
				"disposition": "attachment",
			})
		}
		payload["attachments"] = list
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sendgrid: encoding the message: %w", err)
	}
	return postWithRetry(ctx, t.Client, "sendgrid", func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.Endpoint, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+t.APIKey)
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})
}
