package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// Mailer sends emails via the Resend API.
type Mailer struct {
	apiKey string
	from   string
	client *http.Client
}

// New creates a new Mailer instance.
func New(apiKey, from string) *Mailer {
	return &Mailer{
		apiKey: apiKey,
		from:   from,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// SendOptions configures an email to send.
type SendOptions struct {
	To       string
	Subject  string
	Template string
	Data     map[string]interface{}
}

// Send renders a template and sends the email via Resend.
func (m *Mailer) Send(ctx context.Context, opts SendOptions) error {
	// Render the email template
	htmlBody, err := m.renderTemplate(opts.Template, opts.Data)
	if err != nil {
		return fmt.Errorf("rendering template %q: %w", opts.Template, err)
	}

	// Build Resend API payload
	payload := map[string]interface{}{
		"from":    m.from,
		"to":      []string{opts.To},
		"subject": opts.Subject,
		"html":    htmlBody,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling email payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("resend API error (%d): %v", resp.StatusCode, errResp)
	}

	return nil
}

// SendRaw sends an email with raw HTML content (no template rendering).
func (m *Mailer) SendRaw(ctx context.Context, to, subject, htmlBody string) error {
	payload := map[string]interface{}{
		"from":    m.from,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling email payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("resend API error (%d): %v", resp.StatusCode, errResp)
	}

	return nil
}

// SendInvitation sends a staff invitation email.
func (m *Mailer) SendInvitation(to, businessName, role, inviteURL string) error {
	return m.Send(context.Background(), SendOptions{
		To:       to,
		Subject:  "You've been invited to join " + businessName + " on Grit Motors",
		Template: "invitation",
		Data: map[string]interface{}{
			"BusinessName": businessName,
			"Role":         role,
			"InviteURL":    inviteURL,
		},
	})
}

// SendLowStockAlert sends a low stock notification email to the admin.
func (m *Mailer) SendLowStockAlert(to, productTitle, branchName string, currentStock, threshold int) error {
	return m.Send(context.Background(), SendOptions{
		To:       to,
		Subject:  "Low Stock Alert: " + productTitle + " at " + branchName,
		Template: "low-stock",
		Data: map[string]interface{}{
			"ProductTitle": productTitle,
			"BranchName":   branchName,
			"CurrentStock": currentStock,
			"Threshold":    threshold,
		},
	})
}

func (m *Mailer) renderTemplate(name string, data map[string]interface{}) (string, error) {
	tmplStr, ok := EmailTemplates[name]
	if !ok {
		return "", fmt.Errorf("template %q not found", name)
	}

	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %q: %w", name, err)
	}

	return buf.String(), nil
}
