package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SESTransport sends through Amazon SES API v2 SendEmail over HTTPS
// (MAIL_MAILER=ses), signed with SigV4 and no AWS SDK:
// https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html
//
// The message goes as Raw MIME, which carries cc, reply-to, a text part,
// attachments and custom headers in one shape.
type SESTransport struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	// SessionToken is for temporary credentials (AWS_SESSION_TOKEN).
	SessionToken string
	// Endpoint defaults to https://email.<region>.amazonaws.com.
	Endpoint string
	Client   *http.Client

	now func() time.Time
}

// NewSES returns an SES transport for a region and credentials.
func NewSES(region, accessKeyID, secretAccessKey, sessionToken string) *SESTransport {
	return &SESTransport{
		Region:          region,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
		Endpoint:        "https://email." + region + ".amazonaws.com",
	}
}

// Name implements Transport.
func (t *SESTransport) Name() string { return "ses" }

// Send implements Transport.
func (t *SESTransport) Send(ctx context.Context, m *Message) error {
	raw, err := buildMIME(m, time.Now())
	if err != nil {
		return fmt.Errorf("ses: %w", err)
	}
	destination := map[string]interface{}{}
	for field, list := range map[string][]string{"ToAddresses": m.To, "CcAddresses": m.Cc, "BccAddresses": m.Bcc} {
		if len(list) == 0 {
			continue
		}
		addrs, err := addressesOf(list)
		if err != nil {
			return fmt.Errorf("ses: %w", err)
		}
		destination[field] = addrs
	}
	payload := map[string]interface{}{
		"FromEmailAddress": m.From,
		"Destination":      destination,
		"Content": map[string]interface{}{
			"Raw": map[string]string{"Data": base64.StdEncoding.EncodeToString(raw)},
		},
	}
	if m.ReplyTo != "" {
		payload["ReplyToAddresses"] = []string{m.ReplyTo}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ses: encoding the message: %w", err)
	}

	creds := awsCredentials{AccessKeyID: t.AccessKeyID, SecretAccessKey: t.SecretAccessKey, SessionToken: t.SessionToken}
	endpoint := strings.TrimRight(t.Endpoint, "/") + "/v2/email/outbound-emails"
	return postWithRetry(ctx, t.Client, "ses", func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		now := time.Now
		if t.now != nil {
			now = t.now
		}
		// Signed per attempt, so a retry carries a fresh date.
		signV4(req, body, creds, t.Region, "ses", now())
		return req, nil
	})
}
