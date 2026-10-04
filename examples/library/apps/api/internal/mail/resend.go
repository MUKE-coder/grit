package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// ResendEndpoint is Resend's send-email API:
// https://resend.com/docs/api-reference/emails/send-email
const ResendEndpoint = "https://api.resend.com/emails"

// ResendTransport sends through Resend's HTTP API (MAIL_MAILER=resend).
type ResendTransport struct {
	APIKey   string
	Endpoint string
	Client   *http.Client
}

// NewResend returns a Resend transport for an API key.
func NewResend(apiKey string) *ResendTransport {
	return &ResendTransport{APIKey: apiKey, Endpoint: ResendEndpoint}
}

// Name implements Transport.
func (t *ResendTransport) Name() string { return "resend" }

// Send implements Transport.
func (t *ResendTransport) Send(ctx context.Context, m *Message) error {
	payload := map[string]interface{}{
		"from":    m.From,
		"to":      m.To,
		"subject": m.Subject,
	}
	if len(m.Cc) > 0 {
		payload["cc"] = m.Cc
	}
	if len(m.Bcc) > 0 {
		payload["bcc"] = m.Bcc
	}
	if m.ReplyTo != "" {
		payload["reply_to"] = m.ReplyTo
	}
	if m.HTML != "" {
		payload["html"] = m.HTML
	}
	if m.Text != "" {
		payload["text"] = m.Text
	}
	if len(m.Headers) > 0 {
		payload["headers"] = m.Headers
	}
	if len(m.Attachments) > 0 {
		list := make([]map[string]string, 0, len(m.Attachments))
		for _, a := range m.Attachments {
			list = append(list, map[string]string{
				"filename":     a.Filename,
				"content":      base64.StdEncoding.EncodeToString(a.Content),
				"content_type": attachmentType(a),
			})
		}
		payload["attachments"] = list
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("resend: encoding the message: %w", err)
	}
	// One key for both tries: when the first reached Resend and only its answer
	// was lost, Resend recognises the retry instead of sending the email twice.
	key, err := randomHex(16)
	if err != nil {
		return err
	}
	return postWithRetry(ctx, t.Client, "resend", func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.Endpoint, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+t.APIKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", key)
		return req, nil
	})
}
