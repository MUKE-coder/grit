package scaffold

// The internal/mail package. Every file here is framework-owned: written on
// scaffold, and delivered on upgrade by repairMailDrivers.

// mailPackageFiles maps each file's path under internal/mail to its template.
func mailPackageFiles() map[string]string {
	return map[string]string{
		"mailer.go":                 mailerServiceGo(),
		"transport.go":              mailTransportGo(),
		"from_config.go":            mailFromConfigGo(),
		"resend.go":                 mailResendGo(),
		"smtp.go":                   mailSMTPGo(),
		"mime.go":                   mailMIMEGo(),
		"mailgun.go":                mailMailgunGo(),
		"postmark.go":               mailPostmarkGo(),
		"sendgrid.go":               mailSendGridGo(),
		"ses.go":                    mailSESGo(),
		"sigv4.go":                  mailSigV4Go(),
		"log_transport.go":          mailLogTransportGo(),
		"failover.go":               mailFailoverGo(),
		"mail_test.go":              mailTransportsTestGo(),
		"mailtest/mailtest.go":      mailtestGo(),
		"mailtest/mailtest_test.go": mailtestTestGo(),
	}
}

// mailerServiceGo writes internal/mail/mailer.go: the Mailer every call site uses, over a Transport.
func mailerServiceGo() string {
	return `package mail

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
)

// Mailer sends email through a Transport. FromConfig picks the transport once
// from configuration, so the code that sends mail does not change when the
// provider does.
type Mailer struct {
	transport Transport
	from      string
}

// New returns a Mailer that sends through Resend. It keeps the signature it
// has always had, so code written against it compiles unchanged. FromConfig is
// the constructor that honours MAIL_MAILER.
func New(apiKey, from string) *Mailer {
	return NewWithTransport(NewResend(apiKey), from)
}

// NewWithTransport returns a Mailer over any transport, including the fake in
// internal/mail/mailtest.
func NewWithTransport(t Transport, from string) *Mailer {
	return &Mailer{transport: t, from: from}
}

// Driver names the transport: "smtp", "resend", "failover(smtp,log)" and so on.
func (m *Mailer) Driver() string {
	return m.transport.Name()
}

// From is the address a message without its own From is sent from.
func (m *Mailer) From() string {
	return m.from
}

// SendOptions configures an email to send.
type SendOptions struct {
	To       string
	Subject  string
	Template string
	Data     map[string]interface{}
}

// Send renders a template from EmailTemplates and sends it.
func (m *Mailer) Send(ctx context.Context, opts SendOptions) error {
	htmlBody, err := m.renderTemplate(opts.Template, opts.Data)
	if err != nil {
		return fmt.Errorf("rendering template %q: %w", opts.Template, err)
	}
	return m.SendMessage(ctx, &Message{To: []string{opts.To}, Subject: opts.Subject, HTML: htmlBody})
}

// SendRaw sends an email with raw HTML content (no template rendering).
func (m *Mailer) SendRaw(ctx context.Context, to, subject, htmlBody string) error {
	return m.SendMessage(ctx, &Message{To: []string{to}, Subject: subject, HTML: htmlBody})
}

// SendMessage sends a message using any of Message's fields: several
// recipients, cc, bcc, reply-to, a text part, attachments and headers.
func (m *Mailer) SendMessage(ctx context.Context, msg *Message) error {
	if msg == nil {
		return errors.New("mail: nil message")
	}
	out := *msg
	if out.From == "" {
		out.From = m.from
	}
	if err := out.Validate(); err != nil {
		return err
	}
	return m.transport.Send(ctx, &out)
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
`
}

// mailTransportGo writes internal/mail/transport.go: Message, Attachment, Transport and the shared HTTP retry.
func mailTransportGo() string {
	return `package mail

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
`
}

// mailFromConfigGo writes internal/mail/from_config.go: FromConfig: MAIL_MAILER to a transport.
func mailFromConfigGo() string {
	return `package mail

import (
	"errors"
	"fmt"
	netmail "net/mail"
	"strings"

	"{{MODULE}}/internal/config"
)

// Drivers lists every value MAIL_MAILER accepts.
var Drivers = []string{"smtp", "resend", "mailgun", "postmark", "sendgrid", "ses", "log", "failover"}

// FromConfig builds the Mailer the configuration asks for.
//
// MAIL_MAILER names the driver. Left empty, a real RESEND_API_KEY means
// Resend, which is how every project sent mail before there were drivers, so
// one that set only that key keeps sending the same way. With neither,
// development sends to Mailhog from docker compose and falls back to the log
// when Mailhog is not running, so a reset link is never lost. Production gets
// no mailer (nil, nil), and the code that sends mail says so.
//
// A driver that is named but not configured is an error rather than a quiet
// fallback: an app that believes it sends through SES should not be writing
// its mail to a log.
func FromConfig(cfg *config.Config) (*Mailer, error) {
	driver := strings.ToLower(strings.TrimSpace(cfg.Mail.Mailer))
	var (
		t   Transport
		err error
	)
	switch {
	case driver != "":
		t, err = transportFor(cfg, driver)
	case resendKeySet(cfg.ResendAPIKey):
		t = NewResend(cfg.ResendAPIKey)
	case cfg.AppEnv == "production":
		return nil, nil
	default:
		t, err = failoverFor(cfg, []string{"smtp", "log"})
	}
	if err != nil {
		return nil, err
	}
	return NewWithTransport(t, fromAddress(cfg)), nil
}

func transportFor(cfg *config.Config, driver string) (Transport, error) {
	mc := cfg.Mail
	switch driver {
	case "resend":
		if !resendKeySet(cfg.ResendAPIKey) {
			return nil, errors.New("MAIL_MAILER=resend needs RESEND_API_KEY")
		}
		return NewResend(cfg.ResendAPIKey), nil
	case "smtp":
		return NewSMTP(mc.SMTPHost, mc.SMTPPort, mc.SMTPUsername, mc.SMTPPassword, mc.SMTPEncryption), nil
	case "mailgun":
		if mc.MailgunDomain == "" || mc.MailgunSecret == "" {
			return nil, errors.New("MAIL_MAILER=mailgun needs MAILGUN_DOMAIN and MAILGUN_SECRET")
		}
		return NewMailgun(mc.MailgunDomain, mc.MailgunSecret, mc.MailgunEndpoint), nil
	case "postmark":
		if mc.PostmarkToken == "" {
			return nil, errors.New("MAIL_MAILER=postmark needs POSTMARK_TOKEN")
		}
		return NewPostmark(mc.PostmarkToken, mc.PostmarkMessageStream), nil
	case "sendgrid":
		if mc.SendGridAPIKey == "" {
			return nil, errors.New("MAIL_MAILER=sendgrid needs SENDGRID_API_KEY")
		}
		return NewSendGrid(mc.SendGridAPIKey), nil
	case "ses":
		if mc.SESRegion == "" || mc.SESAccessKeyID == "" || mc.SESSecretAccessKey == "" {
			return nil, errors.New("MAIL_MAILER=ses needs AWS_SES_REGION, AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY")
		}
		return NewSES(mc.SESRegion, mc.SESAccessKeyID, mc.SESSecretAccessKey, mc.SESSessionToken), nil
	case "log":
		if cfg.AppEnv == "production" && !mc.AllowLogInProduction {
			return nil, errors.New("MAIL_MAILER=log writes every message, reset links included, to the log instead of sending it, so production refuses it. Set MAIL_ALLOW_LOG_IN_PRODUCTION=true if that is really what you want")
		}
		return NewLog(mc.LogPath), nil
	case "failover":
		if len(mc.Failover) == 0 {
			return nil, errors.New("MAIL_MAILER=failover needs MAIL_FAILOVER, a comma-separated list of drivers such as smtp,log")
		}
		return failoverFor(cfg, mc.Failover)
	default:
		return nil, fmt.Errorf("MAIL_MAILER=%q is not a mail driver; use one of %s", driver, strings.Join(Drivers, ", "))
	}
}

func failoverFor(cfg *config.Config, names []string) (Transport, error) {
	transports := make([]Transport, 0, len(names))
	for _, name := range names {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		if name == "failover" {
			return nil, errors.New("MAIL_FAILOVER cannot name failover itself")
		}
		t, err := transportFor(cfg, name)
		if err != nil {
			return nil, fmt.Errorf("MAIL_FAILOVER: %w", err)
		}
		transports = append(transports, t)
	}
	return NewFailover(transports...), nil
}

// resendKeySet is false for an empty key and for the placeholders the .env
// files ship with.
func resendKeySet(key string) bool {
	return key != "" && !strings.HasPrefix(key, "re_your_api_key")
}

// fromAddress is MAIL_FROM with MAIL_FROM_NAME as its display name.
func fromAddress(cfg *config.Config) string {
	from := strings.TrimSpace(cfg.MailFrom)
	if cfg.Mail.FromName == "" || strings.Contains(from, "<") {
		return from
	}
	return (&netmail.Address{Name: cfg.Mail.FromName, Address: from}).String()
}
`
}

// mailResendGo writes internal/mail/resend.go: the Resend driver.
func mailResendGo() string {
	return `package mail

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
`
}

// mailSMTPGo writes internal/mail/smtp.go: the SMTP driver.
func mailSMTPGo() string {
	return `package mail

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// SMTPTransport sends over SMTP (MAIL_MAILER=smtp). With nothing configured
// it talks to localhost:1025, which is Mailhog from docker compose.
type SMTPTransport struct {
	Host     string
	Port     string
	Username string
	Password string
	// Encryption is "tls" (implicit TLS, usually port 465), "starttls" or
	// "none". Empty picks tls on port 465, none for localhost and starttls
	// for anything else, so a password never crosses the network in clear.
	Encryption string
	// DialTimeout bounds connecting. The whole exchange is bounded by the
	// context's deadline, or by one minute without one.
	DialTimeout time.Duration
	// TLSConfig replaces the default, which verifies the host's certificate
	// and requires TLS 1.2.
	TLSConfig *tls.Config
}

// NewSMTP returns an SMTP transport.
func NewSMTP(host, port, username, password, encryption string) *SMTPTransport {
	return &SMTPTransport{Host: host, Port: port, Username: username, Password: password, Encryption: encryption}
}

// Name implements Transport.
func (t *SMTPTransport) Name() string { return "smtp" }

func (t *SMTPTransport) hostPort() (string, string) {
	host, port := t.Host, t.Port
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "1025"
	}
	return host, port
}

func (t *SMTPTransport) mode() (string, error) {
	host, port := t.hostPort()
	switch strings.ToLower(strings.TrimSpace(t.Encryption)) {
	case "tls", "ssl":
		return "tls", nil
	case "starttls":
		return "starttls", nil
	case "none":
		return "none", nil
	case "":
		if port == "465" {
			return "tls", nil
		}
		if isLocalHost(host) {
			return "none", nil
		}
		return "starttls", nil
	default:
		return "", fmt.Errorf("smtp: SMTP_ENCRYPTION %q is not tls, starttls or none", t.Encryption)
	}
}

func isLocalHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// Send implements Transport.
func (t *SMTPTransport) Send(ctx context.Context, m *Message) error {
	mode, err := t.mode()
	if err != nil {
		return err
	}
	from, err := addressOf(m.From)
	if err != nil {
		return fmt.Errorf("smtp: %w", err)
	}
	recipients, err := addressesOf(append(append(append([]string(nil), m.To...), m.Cc...), m.Bcc...))
	if err != nil {
		return fmt.Errorf("smtp: %w", err)
	}
	raw, err := buildMIME(m, time.Now())
	if err != nil {
		return fmt.Errorf("smtp: %w", err)
	}

	host, port := t.hostPort()
	addr := net.JoinHostPort(host, port)
	tlsConfig := t.TLSConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
	}
	dialTimeout := t.DialTimeout
	if dialTimeout == 0 {
		dialTimeout = 10 * time.Second
	}
	dialer := &net.Dialer{Timeout: dialTimeout}
	var conn net.Conn
	if mode == "tls" {
		conn, err = (&tls.Dialer{NetDialer: dialer, Config: tlsConfig}).DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("smtp: connecting to %s: %w", addr, err)
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(time.Minute)
	}
	if err := conn.SetDeadline(deadline); err != nil {
		return closeAfter(fmt.Errorf("smtp: setting a deadline: %w", err), conn)
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return closeAfter(fmt.Errorf("smtp: greeting from %s: %w", addr, err), conn)
	}
	if err := t.deliver(client, mode, tlsConfig, host, from, recipients, raw); err != nil {
		return closeAfter(err, client)
	}
	return nil
}

func (t *SMTPTransport) deliver(c *smtp.Client, mode string, tlsConfig *tls.Config, host, from string, recipients []string, raw []byte) error {
	if mode == "starttls" {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return errors.New("smtp: the server does not offer STARTTLS; set SMTP_ENCRYPTION to tls or none")
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("smtp: STARTTLS: %w", err)
		}
	}
	if t.Username != "" {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp: SMTP_USERNAME is set but the server offers no AUTH")
		}
		// PlainAuth refuses to send the password over a connection that is
		// neither TLS nor to localhost.
		if err := c.Auth(smtp.PlainAuth("", t.Username, t.Password, host)); err != nil {
			return fmt.Errorf("smtp: authenticating: %w", err)
		}
	}
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp: MAIL FROM %s: %w", from, err)
	}
	for _, r := range recipients {
		if err := c.Rcpt(r); err != nil {
			return fmt.Errorf("smtp: RCPT TO %s: %w", r, err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp: DATA: %w", err)
	}
	if _, err := w.Write(raw); err != nil {
		return fmt.Errorf("smtp: writing the message: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp: the server refused the message: %w", err)
	}
	if err := c.Quit(); err != nil {
		return fmt.Errorf("smtp: QUIT: %w", err)
	}
	return nil
}

// closeAfter closes c after a failure and returns the failure, with the close
// error joined when there is one worth reporting.
func closeAfter(err error, c io.Closer) error {
	if cerr := c.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
		return errors.Join(err, cerr)
	}
	return err
}
`
}

// mailMIMEGo writes internal/mail/mime.go: the MIME builder SMTP and SES share.
func mailMIMEGo() string {
	return `package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"sort"
	"strings"
	"time"
)

// reservedHeaders are written by buildMIME itself. A custom header by one of
// these names is ignored rather than written twice.
var reservedHeaders = map[string]bool{
	"From": true, "To": true, "Cc": true, "Bcc": true, "Reply-To": true, "Subject": true,
	"Date": true, "Message-Id": true, "Mime-Version": true,
	"Content-Type": true, "Content-Transfer-Encoding": true,
}

// buildMIME renders a message as an RFC 5322 email: the text and HTML bodies
// as multipart/alternative, inside multipart/mixed when there are
// attachments. Bcc is never written as a header, or every recipient would see
// it; the transport passes those addresses in the envelope instead.
func buildMIME(m *Message, now time.Time) ([]byte, error) {
	var out bytes.Buffer
	from, err := parseAddress(m.From)
	if err != nil {
		return nil, err
	}
	writeHeader(&out, "From", from.String())
	for _, list := range []struct {
		name  string
		addrs []string
	}{{"To", m.To}, {"Cc", m.Cc}} {
		if len(list.addrs) == 0 {
			continue
		}
		formatted, err := formatAddressList(list.addrs)
		if err != nil {
			return nil, err
		}
		writeHeader(&out, list.name, formatted)
	}
	if m.ReplyTo != "" {
		replyTo, err := parseAddress(m.ReplyTo)
		if err != nil {
			return nil, err
		}
		writeHeader(&out, "Reply-To", replyTo.String())
	}
	writeHeader(&out, "Subject", mime.QEncoding.Encode("utf-8", m.Subject))
	writeHeader(&out, "Date", now.Format(time.RFC1123Z))
	id, err := messageID(from.Address)
	if err != nil {
		return nil, err
	}
	writeHeader(&out, "Message-ID", id)
	writeHeader(&out, "MIME-Version", "1.0")

	names := make([]string, 0, len(m.Headers))
	for name := range m.Headers {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		canonical := textproto.CanonicalMIMEHeaderKey(name)
		if reservedHeaders[canonical] {
			continue
		}
		writeHeader(&out, canonical, mime.QEncoding.Encode("utf-8", m.Headers[name]))
	}

	bodyHeader, body, err := bodyPart(m)
	if err != nil {
		return nil, err
	}
	if len(m.Attachments) == 0 {
		writeHeader(&out, "Content-Type", bodyHeader.Get("Content-Type"))
		if cte := bodyHeader.Get("Content-Transfer-Encoding"); cte != "" {
			writeHeader(&out, "Content-Transfer-Encoding", cte)
		}
		out.WriteString("\r\n")
		out.Write(body)
		return out.Bytes(), nil
	}

	var mixed bytes.Buffer
	mw := multipart.NewWriter(&mixed)
	part, err := mw.CreatePart(bodyHeader)
	if err != nil {
		return nil, fmt.Errorf("mail: writing the body: %w", err)
	}
	if _, err := part.Write(body); err != nil {
		return nil, fmt.Errorf("mail: writing the body: %w", err)
	}
	for _, a := range m.Attachments {
		if err := writeAttachment(mw, a); err != nil {
			return nil, err
		}
	}
	if err := mw.Close(); err != nil {
		return nil, fmt.Errorf("mail: closing the message: %w", err)
	}
	writeHeader(&out, "Content-Type", "multipart/mixed; boundary="+mw.Boundary())
	out.WriteString("\r\n")
	out.Write(mixed.Bytes())
	return out.Bytes(), nil
}

func writeHeader(out *bytes.Buffer, name, value string) {
	out.WriteString(name + ": " + value + "\r\n")
}

func formatAddressList(list []string) (string, error) {
	formatted := make([]string, 0, len(list))
	for _, s := range list {
		a, err := parseAddress(s)
		if err != nil {
			return "", err
		}
		formatted = append(formatted, a.String())
	}
	return strings.Join(formatted, ", "), nil
}

func messageID(from string) (string, error) {
	domain := "localhost"
	if at := strings.LastIndex(from, "@"); at >= 0 && at < len(from)-1 {
		domain = from[at+1:]
	}
	id, err := randomHex(16)
	if err != nil {
		return "", err
	}
	return "<" + id + "@" + domain + ">", nil
}

// bodyPart returns the headers and encoded content of the message body: one
// text or HTML part, or both as multipart/alternative.
func bodyPart(m *Message) (textproto.MIMEHeader, []byte, error) {
	h := make(textproto.MIMEHeader)
	var buf bytes.Buffer
	if m.Text != "" && m.HTML != "" {
		aw := multipart.NewWriter(&buf)
		for _, p := range []struct{ contentType, body string }{
			{"text/plain; charset=utf-8", m.Text},
			{"text/html; charset=utf-8", m.HTML},
		} {
			ph := make(textproto.MIMEHeader)
			ph.Set("Content-Type", p.contentType)
			ph.Set("Content-Transfer-Encoding", "quoted-printable")
			w, err := aw.CreatePart(ph)
			if err != nil {
				return nil, nil, fmt.Errorf("mail: writing the body: %w", err)
			}
			if err := writeQuotedPrintable(w, p.body); err != nil {
				return nil, nil, err
			}
		}
		if err := aw.Close(); err != nil {
			return nil, nil, fmt.Errorf("mail: writing the body: %w", err)
		}
		h.Set("Content-Type", "multipart/alternative; boundary="+aw.Boundary())
		return h, buf.Bytes(), nil
	}

	contentType, body := "text/html; charset=utf-8", m.HTML
	if m.HTML == "" {
		contentType, body = "text/plain; charset=utf-8", m.Text
	}
	if err := writeQuotedPrintable(&buf, body); err != nil {
		return nil, nil, err
	}
	h.Set("Content-Type", contentType)
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	return h, buf.Bytes(), nil
}

func writeQuotedPrintable(w io.Writer, s string) error {
	qw := quotedprintable.NewWriter(w)
	if _, err := qw.Write([]byte(s)); err != nil {
		return fmt.Errorf("mail: encoding the body: %w", err)
	}
	if err := qw.Close(); err != nil {
		return fmt.Errorf("mail: encoding the body: %w", err)
	}
	return nil
}

func writeAttachment(mw *multipart.Writer, a Attachment) error {
	mediaType, params, err := mime.ParseMediaType(attachmentType(a))
	if err != nil {
		mediaType, params = "application/octet-stream", map[string]string{}
	}
	params["name"] = a.Filename
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", mime.FormatMediaType(mediaType, params))
	h.Set("Content-Transfer-Encoding", "base64")
	h.Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": a.Filename}))
	part, err := mw.CreatePart(h)
	if err != nil {
		return fmt.Errorf("mail: attaching %s: %w", a.Filename, err)
	}
	encoded := base64.StdEncoding.EncodeToString(a.Content)
	for len(encoded) > 76 {
		if _, err := io.WriteString(part, encoded[:76]+"\r\n"); err != nil {
			return fmt.Errorf("mail: attaching %s: %w", a.Filename, err)
		}
		encoded = encoded[76:]
	}
	if _, err := io.WriteString(part, encoded+"\r\n"); err != nil {
		return fmt.Errorf("mail: attaching %s: %w", a.Filename, err)
	}
	return nil
}
`
}

// mailMailgunGo writes internal/mail/mailgun.go: the Mailgun driver.
func mailMailgunGo() string {
	return `package mail

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
`
}

// mailPostmarkGo writes internal/mail/postmark.go: the Postmark driver.
func mailPostmarkGo() string {
	return `package mail

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
`
}

// mailSendGridGo writes internal/mail/sendgrid.go: the SendGrid driver.
func mailSendGridGo() string {
	return `package mail

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
`
}

// mailSESGo writes internal/mail/ses.go: the Amazon SES driver.
func mailSESGo() string {
	return `package mail

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
`
}

// mailSigV4Go writes internal/mail/sigv4.go: the SigV4 signer SES uses.
func mailSigV4Go() string {
	return `package mail

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// awsCredentials are what a SigV4 signature is made with.
type awsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// signV4 signs req with AWS Signature Version 4, following
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html
// Every header already on the request is signed, with host. payload must be
// the request body.
func signV4(req *http.Request, payload []byte, creds awsCredentials, region, service string, now time.Time) {
	amzDate := now.UTC().Format("20060102T150405Z")
	req.Header.Set("X-Amz-Date", amzDate)
	if creds.SessionToken != "" {
		req.Header.Set("X-Amz-Security-Token", creds.SessionToken)
	}
	canonical, signedHeaders := canonicalRequest(req, payload)
	scope := amzDate[:8] + "/" + region + "/" + service + "/aws4_request"
	key := signingKey(creds.SecretAccessKey, amzDate[:8], region, service)
	signature := hex.EncodeToString(hmacSHA256(key, stringToSign(amzDate, scope, canonical)))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+creds.AccessKeyID+"/"+scope+
		", SignedHeaders="+signedHeaders+", Signature="+signature)
}

// canonicalRequest returns the canonical request and its signed-headers list.
func canonicalRequest(req *http.Request, payload []byte) (string, string) {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	values := map[string][]string{"host": {host}}
	for name, v := range req.Header {
		lower := strings.ToLower(name)
		if lower == "authorization" {
			continue
		}
		values[lower] = v
	}
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)

	var headers strings.Builder
	for _, name := range names {
		trimmed := make([]string, 0, len(values[name]))
		for _, v := range values[name] {
			trimmed = append(trimmed, strings.Join(strings.Fields(v), " "))
		}
		headers.WriteString(name + ":" + strings.Join(trimmed, ",") + "\n")
	}
	signedHeaders := strings.Join(names, ";")

	path := req.URL.EscapedPath()
	if path == "" {
		path = "/"
	}
	canonical := strings.Join([]string{
		req.Method,
		path,
		canonicalQuery(req.URL),
		headers.String(),
		signedHeaders,
		sha256Hex(payload),
	}, "\n")
	return canonical, signedHeaders
}

func canonicalQuery(u *url.URL) string {
	query := u.Query()
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var pairs []string
	for _, k := range keys {
		vals := append([]string(nil), query[k]...)
		sort.Strings(vals)
		for _, v := range vals {
			pairs = append(pairs, awsEscape(k)+"="+awsEscape(v))
		}
	}
	return strings.Join(pairs, "&")
}

// awsEscape percent-encodes everything but A-Z a-z 0-9 - _ . ~, with a space
// as %20 rather than +.
func awsEscape(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}

func stringToSign(amzDate, scope, canonical string) string {
	return "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + sha256Hex([]byte(canonical))
}

func signingKey(secret, date, region, service string) []byte {
	k := hmacSHA256([]byte("AWS4"+secret), date)
	k = hmacSHA256(k, region)
	k = hmacSHA256(k, service)
	return hmacSHA256(k, "aws4_request")
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
`
}

// mailLogTransportGo writes internal/mail/log_transport.go: the log driver.
func mailLogTransportGo() string {
	return `package mail

import (
	"context"
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogTransport writes each message to the log and to an HTML file under Dir
// instead of sending it (MAIL_MAILER=log). It is for development: the log
// holds every link a message carries, password resets included, which is why
// production refuses it unless MAIL_ALLOW_LOG_IN_PRODUCTION=true.
type LogTransport struct {
	Dir string
	now func() time.Time
}

// NewLog returns a log transport that saves messages under dir, or
// storage/mail when dir is empty.
func NewLog(dir string) *LogTransport {
	if dir == "" {
		dir = filepath.Join("storage", "mail")
	}
	return &LogTransport{Dir: dir, now: time.Now}
}

// Name implements Transport.
func (t *LogTransport) Name() string { return "log" }

// Send implements Transport.
func (t *LogTransport) Send(_ context.Context, m *Message) error {
	now := time.Now
	if t.now != nil {
		now = t.now
	}
	if err := os.MkdirAll(t.Dir, 0o750); err != nil {
		return fmt.Errorf("log mailer: creating %s: %w", t.Dir, err)
	}
	path := filepath.Join(t.Dir, now().UTC().Format("20060102-150405.000000000")+"-"+slugify(m.Subject)+".html")

	body := m.HTML
	if body == "" {
		body = "<pre>" + html.EscapeString(m.Text) + "</pre>"
	}
	header := fmt.Sprintf("<!--\nFrom: %s\nTo: %s\nCc: %s\nBcc: %s\nReply-To: %s\nSubject: %s\nAttachments: %d\n-->\n",
		commentSafe(m.From), commentSafe(strings.Join(m.To, ", ")), commentSafe(strings.Join(m.Cc, ", ")),
		commentSafe(strings.Join(m.Bcc, ", ")), commentSafe(m.ReplyTo), commentSafe(m.Subject), len(m.Attachments))
	if err := os.WriteFile(path, []byte(header+body), 0o600); err != nil {
		return fmt.Errorf("log mailer: writing %s: %w", path, err)
	}

	text := m.Text
	if text == "" {
		text = m.HTML
	}
	log.Printf("mail (log driver): from=%s to=%s subject=%q saved to %s\n%s",
		m.From, strings.Join(append(append(append([]string(nil), m.To...), m.Cc...), m.Bcc...), ", "), m.Subject, path, text)
	return nil
}

func commentSafe(s string) string {
	return strings.ReplaceAll(s, "--", "- -")
}

func slugify(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
			continue
		}
		if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
		if b.Len() >= 50 {
			break
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "message"
	}
	return slug
}
`
}

// mailFailoverGo writes internal/mail/failover.go: the failover driver.
func mailFailoverGo() string {
	return `package mail

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

// FailoverTransport tries its transports in order and stops at the first that
// delivers (MAIL_MAILER=failover, with MAIL_FAILOVER naming them). Each failure
// is logged, so a provider that is quietly down still shows up.
type FailoverTransport struct {
	Transports []Transport
}

// NewFailover returns a transport that tries each of transports in turn.
func NewFailover(transports ...Transport) *FailoverTransport {
	return &FailoverTransport{Transports: transports}
}

// Name implements Transport: "failover(smtp,log)".
func (t *FailoverTransport) Name() string {
	names := make([]string, 0, len(t.Transports))
	for _, tr := range t.Transports {
		names = append(names, tr.Name())
	}
	return "failover(" + strings.Join(names, ",") + ")"
}

// Send implements Transport.
func (t *FailoverTransport) Send(ctx context.Context, m *Message) error {
	if len(t.Transports) == 0 {
		return errors.New("failover: no transports to try")
	}
	var errs []error
	for i, tr := range t.Transports {
		err := tr.Send(ctx, m)
		if err == nil {
			if i > 0 {
				log.Printf("mail: %s delivered %q after %d failed", tr.Name(), m.Subject, i)
			}
			return nil
		}
		log.Printf("mail: %s could not send %q: %v", tr.Name(), m.Subject, err)
		errs = append(errs, fmt.Errorf("%s: %w", tr.Name(), err))
		if ctx.Err() != nil {
			break
		}
	}
	return fmt.Errorf("failover: every transport failed: %w", errors.Join(errs...))
}
`
}

// mailTransportsTestGo writes internal/mail/mail_test.go: tests for every driver.
func mailTransportsTestGo() string {
	return `package mail

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	netmail "net/mail"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"{{MODULE}}/internal/config"
)

// captured is one request a fake provider received.
type captured struct {
	Method string
	Path   string
	Header http.Header
	Body   []byte
}

// fakeProvider records every request and answers them with statuses in turn,
// repeating the last one; 200 when none are given.
func fakeProvider(t *testing.T, statuses ...int) (*httptest.Server, func() []captured) {
	t.Helper()
	var mu sync.Mutex
	var got []captured
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading the request body: %v", err)
		}
		mu.Lock()
		got = append(got, captured{Method: r.Method, Path: r.URL.Path, Header: r.Header.Clone(), Body: body})
		n := len(got)
		mu.Unlock()
		status := http.StatusOK
		if len(statuses) > 0 {
			i := n - 1
			if i >= len(statuses) {
				i = len(statuses) - 1
			}
			status = statuses[i]
		}
		w.WriteHeader(status)
		if _, err := w.Write([]byte("{}")); err != nil {
			t.Errorf("writing the response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, func() []captured {
		mu.Lock()
		defer mu.Unlock()
		return append([]captured(nil), got...)
	}
}

func noRetryDelay(t *testing.T) {
	saved := retryDelay
	retryDelay = 0
	t.Cleanup(func() { retryDelay = saved })
}

func sampleMessage() *Message {
	return &Message{
		From:        "Grit App <noreply@example.com>",
		ReplyTo:     "support@example.com",
		To:          []string{"Ada <ada@example.com>", "bob@example.com"},
		Cc:          []string{"carol@example.com"},
		Bcc:         []string{"dan@example.com"},
		Subject:     "Your invoice",
		HTML:        "<p>Hello</p>",
		Text:        "Hello",
		Attachments: []Attachment{{Filename: "invoice.pdf", ContentType: "application/pdf", Content: []byte("PDF-1.4 test")}},
		Headers:     map[string]string{"X-Entity-Ref": "inv-42"},
	}
}

func onlyRequest(t *testing.T, requests func() []captured) captured {
	t.Helper()
	got := requests()
	if len(got) != 1 {
		t.Fatalf("the provider got %d requests, want 1", len(got))
	}
	return got[0]
}

func decodeJSON(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()
	var out map[string]interface{}
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("the body is not JSON: %v\n%s", err, body)
	}
	return out
}

func asMap(t *testing.T, name string, v interface{}) map[string]interface{} {
	t.Helper()
	m, ok := v.(map[string]interface{})
	if !ok {
		t.Fatalf("%s = %#v, want an object", name, v)
	}
	return m
}

func asList(t *testing.T, name string, v interface{}) []interface{} {
	t.Helper()
	l, ok := v.([]interface{})
	if !ok {
		t.Fatalf("%s = %#v, want a list", name, v)
	}
	return l
}

func expect(t *testing.T, name string, got, want interface{}) {
	t.Helper()
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func expectBase64(t *testing.T, name string, got interface{}, want string) {
	t.Helper()
	s, ok := got.(string)
	if !ok {
		t.Fatalf("%s = %#v, want a base64 string", name, got)
	}
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil || string(raw) != want {
		t.Errorf("%s decodes to %q (%v), want %q", name, raw, err, want)
	}
}

func TestResendRequestShape(t *testing.T) {
	srv, requests := fakeProvider(t)
	tr := NewResend("re_test")
	tr.Endpoint = srv.URL + "/emails"
	if err := tr.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	r := onlyRequest(t, requests)
	expect(t, "method", r.Method, http.MethodPost)
	expect(t, "path", r.Path, "/emails")
	expect(t, "Authorization", r.Header.Get("Authorization"), "Bearer re_test")
	expect(t, "Content-Type", r.Header.Get("Content-Type"), "application/json")
	if len(r.Header.Get("Idempotency-Key")) != 32 {
		t.Errorf("Idempotency-Key = %q, want 32 hex characters", r.Header.Get("Idempotency-Key"))
	}

	body := decodeJSON(t, r.Body)
	expect(t, "from", body["from"], "Grit App <noreply@example.com>")
	expect(t, "to", body["to"], []interface{}{"Ada <ada@example.com>", "bob@example.com"})
	expect(t, "cc", body["cc"], []interface{}{"carol@example.com"})
	expect(t, "bcc", body["bcc"], []interface{}{"dan@example.com"})
	expect(t, "reply_to", body["reply_to"], "support@example.com")
	expect(t, "subject", body["subject"], "Your invoice")
	expect(t, "html", body["html"], "<p>Hello</p>")
	expect(t, "text", body["text"], "Hello")
	expect(t, "headers", asMap(t, "headers", body["headers"])["X-Entity-Ref"], "inv-42")
	att := asMap(t, "attachment", asList(t, "attachments", body["attachments"])[0])
	expect(t, "attachment filename", att["filename"], "invoice.pdf")
	expect(t, "attachment content_type", att["content_type"], "application/pdf")
	expectBase64(t, "attachment content", att["content"], "PDF-1.4 test")
}

func TestHTTPTransportsRetryOnceOn5xxAndNeverOn4xx(t *testing.T) {
	noRetryDelay(t)
	ctx := context.Background()

	srv, requests := fakeProvider(t, http.StatusServiceUnavailable, http.StatusOK)
	tr := NewResend("k")
	tr.Endpoint = srv.URL
	if err := tr.Send(ctx, sampleMessage()); err != nil {
		t.Fatalf("a 503 then a 200 should succeed: %v", err)
	}
	got := requests()
	if len(got) != 2 {
		t.Fatalf("%d requests after a 503, want 2", len(got))
	}
	if got[0].Header.Get("Idempotency-Key") != got[1].Header.Get("Idempotency-Key") {
		t.Error("the retry carried a different Idempotency-Key, so Resend could send the email twice")
	}

	srv, requests = fakeProvider(t, http.StatusInternalServerError)
	tr.Endpoint = srv.URL
	var perr *ProviderError
	if err := tr.Send(ctx, sampleMessage()); !errors.As(err, &perr) || perr.Status != http.StatusInternalServerError {
		t.Errorf("two 500s gave %v, want a ProviderError with status 500", err)
	}
	if n := len(requests()); n != 2 {
		t.Errorf("%d requests for a provider that keeps failing, want 2", n)
	}

	srv, requests = fakeProvider(t, http.StatusUnprocessableEntity, http.StatusOK)
	tr.Endpoint = srv.URL
	if err := tr.Send(ctx, sampleMessage()); !errors.As(err, &perr) || perr.Status != http.StatusUnprocessableEntity {
		t.Errorf("a 422 gave %v, want a ProviderError with status 422", err)
	}
	if n := len(requests()); n != 1 {
		t.Errorf("a 422 was sent %d times, want once", n)
	}

	// A connection dropped before any answer is a network error: retried.
	var calls int32
	dropping := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Error("the test server cannot hijack connections")
				return
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				t.Errorf("hijacking: %v", err)
				return
			}
			if err := conn.Close(); err != nil {
				t.Errorf("closing: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer dropping.Close()
	pm := NewPostmark("token", "")
	pm.Endpoint = dropping.URL
	if err := pm.Send(ctx, sampleMessage()); err != nil {
		t.Errorf("a dropped connection then a 200 should succeed: %v", err)
	}
	if n := atomic.LoadInt32(&calls); n != 2 {
		t.Errorf("%d requests after a dropped connection, want 2", n)
	}
}

func TestPostmarkRequestShape(t *testing.T) {
	srv, requests := fakeProvider(t)
	tr := NewPostmark("server-token", "outbound")
	tr.Endpoint = srv.URL + "/email"
	if err := tr.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	r := onlyRequest(t, requests)
	expect(t, "method", r.Method, http.MethodPost)
	expect(t, "path", r.Path, "/email")
	expect(t, "X-Postmark-Server-Token", r.Header.Get("X-Postmark-Server-Token"), "server-token")
	expect(t, "Accept", r.Header.Get("Accept"), "application/json")
	expect(t, "Content-Type", r.Header.Get("Content-Type"), "application/json")

	body := decodeJSON(t, r.Body)
	expect(t, "From", body["From"], "Grit App <noreply@example.com>")
	expect(t, "To", body["To"], "Ada <ada@example.com>,bob@example.com")
	expect(t, "Cc", body["Cc"], "carol@example.com")
	expect(t, "Bcc", body["Bcc"], "dan@example.com")
	expect(t, "ReplyTo", body["ReplyTo"], "support@example.com")
	expect(t, "Subject", body["Subject"], "Your invoice")
	expect(t, "HtmlBody", body["HtmlBody"], "<p>Hello</p>")
	expect(t, "TextBody", body["TextBody"], "Hello")
	expect(t, "MessageStream", body["MessageStream"], "outbound")
	header := asMap(t, "header", asList(t, "Headers", body["Headers"])[0])
	expect(t, "header Name", header["Name"], "X-Entity-Ref")
	expect(t, "header Value", header["Value"], "inv-42")
	att := asMap(t, "attachment", asList(t, "Attachments", body["Attachments"])[0])
	expect(t, "attachment Name", att["Name"], "invoice.pdf")
	expect(t, "attachment ContentType", att["ContentType"], "application/pdf")
	expectBase64(t, "attachment Content", att["Content"], "PDF-1.4 test")
}

func TestSendGridRequestShape(t *testing.T) {
	srv, requests := fakeProvider(t, http.StatusAccepted)
	tr := NewSendGrid("SG.test")
	tr.Endpoint = srv.URL + "/v3/mail/send"
	if err := tr.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	r := onlyRequest(t, requests)
	expect(t, "method", r.Method, http.MethodPost)
	expect(t, "path", r.Path, "/v3/mail/send")
	expect(t, "Authorization", r.Header.Get("Authorization"), "Bearer SG.test")
	expect(t, "Content-Type", r.Header.Get("Content-Type"), "application/json")

	body := decodeJSON(t, r.Body)
	from := asMap(t, "from", body["from"])
	expect(t, "from email", from["email"], "noreply@example.com")
	expect(t, "from name", from["name"], "Grit App")
	p := asMap(t, "personalization", asList(t, "personalizations", body["personalizations"])[0])
	to := asList(t, "to", p["to"])
	expect(t, "to[0] email", asMap(t, "to[0]", to[0])["email"], "ada@example.com")
	expect(t, "to[0] name", asMap(t, "to[0]", to[0])["name"], "Ada")
	expect(t, "to[1] email", asMap(t, "to[1]", to[1])["email"], "bob@example.com")
	expect(t, "cc", asMap(t, "cc[0]", asList(t, "cc", p["cc"])[0])["email"], "carol@example.com")
	expect(t, "bcc", asMap(t, "bcc[0]", asList(t, "bcc", p["bcc"])[0])["email"], "dan@example.com")
	expect(t, "reply_to", asMap(t, "reply_to", body["reply_to"])["email"], "support@example.com")
	expect(t, "subject", body["subject"], "Your invoice")
	content := asList(t, "content", body["content"])
	expect(t, "content[0] type", asMap(t, "content[0]", content[0])["type"], "text/plain")
	expect(t, "content[1] type", asMap(t, "content[1]", content[1])["type"], "text/html")
	expect(t, "content[1] value", asMap(t, "content[1]", content[1])["value"], "<p>Hello</p>")
	expect(t, "headers", asMap(t, "headers", body["headers"])["X-Entity-Ref"], "inv-42")
	att := asMap(t, "attachment", asList(t, "attachments", body["attachments"])[0])
	expect(t, "attachment filename", att["filename"], "invoice.pdf")
	expect(t, "attachment type", att["type"], "application/pdf")
	expect(t, "attachment disposition", att["disposition"], "attachment")
	expectBase64(t, "attachment content", att["content"], "PDF-1.4 test")
}

func TestMailgunRequestShape(t *testing.T) {
	srv, requests := fakeProvider(t)
	tr := NewMailgun("mg.example.com", "key-test", srv.URL)
	if err := tr.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	r := onlyRequest(t, requests)
	expect(t, "method", r.Method, http.MethodPost)
	expect(t, "path", r.Path, "/v3/mg.example.com/messages")
	expect(t, "Authorization", r.Header.Get("Authorization"), "Basic "+base64.StdEncoding.EncodeToString([]byte("api:key-test")))

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/form-data" {
		t.Fatalf("Content-Type = %q (%v), want multipart/form-data", r.Header.Get("Content-Type"), err)
	}
	form, err := multipart.NewReader(bytes.NewReader(r.Body), params["boundary"]).ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("reading the form: %v", err)
	}
	expect(t, "from", form.Value["from"], []string{"Grit App <noreply@example.com>"})
	expect(t, "to", form.Value["to"], []string{"Ada <ada@example.com>", "bob@example.com"})
	expect(t, "cc", form.Value["cc"], []string{"carol@example.com"})
	expect(t, "bcc", form.Value["bcc"], []string{"dan@example.com"})
	expect(t, "subject", form.Value["subject"], []string{"Your invoice"})
	expect(t, "text", form.Value["text"], []string{"Hello"})
	expect(t, "html", form.Value["html"], []string{"<p>Hello</p>"})
	expect(t, "h:Reply-To", form.Value["h:Reply-To"], []string{"support@example.com"})
	expect(t, "h:X-Entity-Ref", form.Value["h:X-Entity-Ref"], []string{"inv-42"})
	files := form.File["attachment"]
	if len(files) != 1 {
		t.Fatalf("%d attachment files, want 1", len(files))
	}
	expect(t, "attachment filename", files[0].Filename, "invoice.pdf")
	expect(t, "attachment Content-Type", files[0].Header.Get("Content-Type"), "application/pdf")
	f, err := files[0].Open()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, "attachment content", string(content), "PDF-1.4 test")

	expect(t, "default endpoint", NewMailgun("mg.example.com", "k", "").url(), "https://api.mailgun.net/v3/mg.example.com/messages")
	expect(t, "EU endpoint", NewMailgun("mg.example.com", "k", "api.eu.mailgun.net").url(), "https://api.eu.mailgun.net/v3/mg.example.com/messages")
}

func TestSESRequestShape(t *testing.T) {
	srv, requests := fakeProvider(t)
	tr := NewSES("eu-west-1", "AKIDEXAMPLE", "not-a-real-secret", "session-token")
	tr.Endpoint = srv.URL
	tr.now = func() time.Time { return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC) }
	if err := tr.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	r := onlyRequest(t, requests)
	expect(t, "method", r.Method, http.MethodPost)
	expect(t, "path", r.Path, "/v2/email/outbound-emails")
	expect(t, "X-Amz-Date", r.Header.Get("X-Amz-Date"), "20260102T030405Z")
	expect(t, "X-Amz-Security-Token", r.Header.Get("X-Amz-Security-Token"), "session-token")
	auth := r.Header.Get("Authorization")
	prefix := "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20260102/eu-west-1/ses/aws4_request, SignedHeaders=content-type;host;x-amz-date;x-amz-security-token, Signature="
	if !strings.HasPrefix(auth, prefix) || len(auth) != len(prefix)+64 {
		t.Errorf("Authorization = %q, want %s followed by 64 hex characters", auth, prefix)
	}

	body := decodeJSON(t, r.Body)
	expect(t, "FromEmailAddress", body["FromEmailAddress"], "Grit App <noreply@example.com>")
	dest := asMap(t, "Destination", body["Destination"])
	expect(t, "ToAddresses", dest["ToAddresses"], []interface{}{"ada@example.com", "bob@example.com"})
	expect(t, "CcAddresses", dest["CcAddresses"], []interface{}{"carol@example.com"})
	expect(t, "BccAddresses", dest["BccAddresses"], []interface{}{"dan@example.com"})
	expect(t, "ReplyToAddresses", body["ReplyToAddresses"], []interface{}{"support@example.com"})
	data, ok := asMap(t, "Raw", asMap(t, "Content", body["Content"])["Raw"])["Data"].(string)
	if !ok {
		t.Fatal("Content.Raw.Data is not a string")
	}
	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		t.Fatalf("Content.Raw.Data is not base64: %v", err)
	}
	msg, err := netmail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("Content.Raw.Data is not a MIME message: %v", err)
	}
	expect(t, "raw Subject", msg.Header.Get("Subject"), "Your invoice")
	if msg.Header.Get("Bcc") != "" {
		t.Error("the raw message carries a Bcc header every recipient would see")
	}
}

// The expected values are from the AWS Signature Version 4 test suite
// (get-vanilla and post-x-www-form-urlencoded), as published with botocore in
// tests/unit/auth/aws4_testsuite. The key is the suite's documented example
// key, not a credential.
func TestSigV4MatchesTheAWSTestSuite(t *testing.T) {
	creds := awsCredentials{AccessKeyID: "AKIDEXAMPLE", SecretAccessKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"} // #nosec G101 -- AWS test-suite example key
	now := time.Date(2015, 8, 30, 12, 36, 0, 0, time.UTC)
	cases := []struct {
		name, method, contentType, body, creq, sts, authz string
	}{
		{
			name:   "get-vanilla",
			method: http.MethodGet,
			creq:   "GET\n/\n\nhost:example.amazonaws.com\nx-amz-date:20150830T123600Z\n\nhost;x-amz-date\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			sts:    "AWS4-HMAC-SHA256\n20150830T123600Z\n20150830/us-east-1/service/aws4_request\nbb579772317eb040ac9ed261061d46c1f17a8133879d6129b6e1c25292927e63",
			authz:  "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/service/aws4_request, SignedHeaders=host;x-amz-date, Signature=5fa00fa31553b73ebf1942676e86291e8372ff2a2260956d9b8aae1d763fbf31",
		},
		{
			name:        "post-x-www-form-urlencoded",
			method:      http.MethodPost,
			contentType: "application/x-www-form-urlencoded",
			body:        "Param1=value1",
			creq:        "POST\n/\n\ncontent-type:application/x-www-form-urlencoded\nhost:example.amazonaws.com\nx-amz-date:20150830T123600Z\n\ncontent-type;host;x-amz-date\n9095672bbd1f56dfc5b65f3e153adc8731a4a654192329106275f4c7b24d0b6e",
			sts:         "AWS4-HMAC-SHA256\n20150830T123600Z\n20150830/us-east-1/service/aws4_request\n42a5e5bb34198acb3e84da4f085bb7927f2bc277ca766e6d19c73c2154021281",
			authz:       "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/service/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=ff11897932ad3f4e8b18135d722051e5ac45fc38421b1da7b9d196a0fe09473a",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "https://example.amazonaws.com/", strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			signV4(req, []byte(tc.body), creds, "us-east-1", "service", now)
			creq, _ := canonicalRequest(req, []byte(tc.body))
			expect(t, "canonical request", creq, tc.creq)
			expect(t, "string to sign", stringToSign("20150830T123600Z", "20150830/us-east-1/service/aws4_request", creq), tc.sts)
			expect(t, "Authorization", req.Header.Get("Authorization"), tc.authz)
		})
	}
}

func TestMIMEMessage(t *testing.T) {
	raw, err := buildMIME(sampleMessage(), time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := netmail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("not a valid message: %v", err)
	}
	expect(t, "From", msg.Header.Get("From"), "\"Grit App\" <noreply@example.com>")
	expect(t, "To", msg.Header.Get("To"), "\"Ada\" <ada@example.com>, <bob@example.com>")
	expect(t, "Cc", msg.Header.Get("Cc"), "<carol@example.com>")
	expect(t, "Reply-To", msg.Header.Get("Reply-To"), "<support@example.com>")
	expect(t, "Subject", msg.Header.Get("Subject"), "Your invoice")
	expect(t, "X-Entity-Ref", msg.Header.Get("X-Entity-Ref"), "inv-42")
	if msg.Header.Get("Bcc") != "" {
		t.Error("the message carries a Bcc header every recipient would see")
	}
	if !strings.HasSuffix(msg.Header.Get("Message-ID"), "@example.com>") {
		t.Errorf("Message-ID = %q", msg.Header.Get("Message-ID"))
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/mixed" {
		t.Fatalf("Content-Type = %q, want multipart/mixed", msg.Header.Get("Content-Type"))
	}
	mixed := multipart.NewReader(msg.Body, params["boundary"])
	body, err := mixed.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	altType, altParams, err := mime.ParseMediaType(body.Header.Get("Content-Type"))
	if err != nil || altType != "multipart/alternative" {
		t.Fatalf("first part = %q, want multipart/alternative", body.Header.Get("Content-Type"))
	}
	alt := multipart.NewReader(body, altParams["boundary"])
	for _, want := range []struct{ contentType, body string }{{"text/plain; charset=utf-8", "Hello"}, {"text/html; charset=utf-8", "<p>Hello</p>"}} {
		part, err := alt.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		content, err := io.ReadAll(part)
		if err != nil {
			t.Fatal(err)
		}
		expect(t, want.contentType+" part", string(content), want.body)
	}
	att, err := mixed.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	expect(t, "attachment filename", att.FileName(), "invoice.pdf")
	encoded, err := io.ReadAll(att)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(string(encoded), "\r\n", ""))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, "attachment content", string(decoded), "PDF-1.4 test")
}

func TestValidateRejectsHeaderInjection(t *testing.T) {
	for name, mutate := range map[string]func(*Message){
		"subject":    func(m *Message) { m.Subject = "Hi\r\nBcc: everyone@example.com" },
		"recipient":  func(m *Message) { m.To = []string{"ada@example.com\nBcc: x@example.com"} },
		"header":     func(m *Message) { m.Headers = map[string]string{"X-Ref": "a\r\nBcc: x@example.com"} },
		"no body":    func(m *Message) { m.HTML, m.Text = "", "" },
		"no address": func(m *Message) { m.To, m.Cc, m.Bcc = nil, nil, nil },
		"bad from":   func(m *Message) { m.From = "not an address" },
	} {
		m := sampleMessage()
		mutate(m)
		if err := m.Validate(); err == nil {
			t.Errorf("%s: Validate accepted the message", name)
		}
	}
	if err := sampleMessage().Validate(); err != nil {
		t.Errorf("the sample message is valid, got %v", err)
	}
}

// smtpSink is just enough of an SMTP server to receive one message.
type smtpSink struct {
	addr string
	done chan struct{}

	mu   sync.Mutex
	from string
	rcpt []string
	data string
}

func startSMTPSink(t *testing.T) *smtpSink {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := &smtpSink{addr: ln.Addr().String(), done: make(chan struct{})}
	t.Cleanup(func() {
		if err := ln.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Errorf("closing the sink: %v", err)
		}
	})
	go s.serve(ln)
	return s
}

func (s *smtpSink) serve(ln net.Listener) {
	conn, err := ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	reply := func(lines ...string) bool {
		for _, line := range lines {
			if _, err := rw.WriteString(line + "\r\n"); err != nil {
				return false
			}
		}
		return rw.Flush() == nil
	}
	if !reply("220 sink ready") {
		return
	}
	var data strings.Builder
	inData := false
	for {
		line, err := rw.ReadString('\n')
		if err != nil {
			return
		}
		if inData {
			if line == ".\r\n" {
				inData = false
				s.mu.Lock()
				s.data = data.String()
				s.mu.Unlock()
				if !reply("250 queued") {
					return
				}
				continue
			}
			data.WriteString(line)
			continue
		}
		cmd := strings.ToUpper(strings.TrimSpace(line))
		ok := true
		switch {
		case strings.HasPrefix(cmd, "EHLO"):
			ok = reply("250 sink")
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			s.mu.Lock()
			s.from = strings.TrimSpace(line[len("MAIL FROM:"):])
			s.mu.Unlock()
			ok = reply("250 ok")
		case strings.HasPrefix(cmd, "RCPT TO:"):
			s.mu.Lock()
			s.rcpt = append(s.rcpt, strings.TrimSpace(line[len("RCPT TO:"):]))
			s.mu.Unlock()
			ok = reply("250 ok")
		case cmd == "DATA":
			inData = true
			ok = reply("354 go ahead")
		case cmd == "QUIT":
			reply("221 bye")
			close(s.done)
			return
		default:
			ok = reply("250 ok")
		}
		if !ok {
			return
		}
	}
}

func TestSMTPTransportDelivers(t *testing.T) {
	sink := startSMTPSink(t)
	host, port, err := net.SplitHostPort(sink.addr)
	if err != nil {
		t.Fatal(err)
	}
	tr := NewSMTP(host, port, "", "", "")
	if err := tr.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-sink.done:
	case <-time.After(5 * time.Second):
		t.Fatal("the sink never saw QUIT")
	}
	sink.mu.Lock()
	defer sink.mu.Unlock()
	expect(t, "MAIL FROM", sink.from, "<noreply@example.com>")
	expect(t, "RCPT TO", sink.rcpt, []string{"<ada@example.com>", "<bob@example.com>", "<carol@example.com>", "<dan@example.com>"})
	if !strings.Contains(sink.data, "Subject: Your invoice\r\n") {
		t.Errorf("the message has no Subject header:\n%s", sink.data)
	}
	if strings.Contains(sink.data, "dan@example.com") {
		t.Error("the Bcc address is in the message itself, where every recipient sees it")
	}
}

func TestSMTPRefusesToSkipAskedForSTARTTLS(t *testing.T) {
	sink := startSMTPSink(t)
	host, port, err := net.SplitHostPort(sink.addr)
	if err != nil {
		t.Fatal(err)
	}
	tr := NewSMTP(host, port, "", "", "starttls")
	err = tr.Send(context.Background(), sampleMessage())
	if err == nil || !strings.Contains(err.Error(), "STARTTLS") {
		t.Errorf("a server without STARTTLS gave %v, want an error naming STARTTLS", err)
	}
	if _, err := NewSMTP("mail.example.com", "587", "", "", "").mode(); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ host, port, encryption, want string }{
		{"localhost", "1025", "", "none"},
		{"127.0.0.1", "1025", "", "none"},
		{"smtp.example.com", "587", "", "starttls"},
		{"smtp.example.com", "465", "", "tls"},
		{"smtp.example.com", "25", "none", "none"},
	} {
		got, err := NewSMTP(tc.host, tc.port, "", "", tc.encryption).mode()
		if err != nil || got != tc.want {
			t.Errorf("%s:%s encryption %q = %q (%v), want %q", tc.host, tc.port, tc.encryption, got, err, tc.want)
		}
	}
}

func TestLogTransportWritesTheMessage(t *testing.T) {
	dir := t.TempDir()
	tr := NewLog(dir)
	if err := tr.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || !strings.HasSuffix(entries[0].Name(), "-your-invoice.html") {
		t.Fatalf("files = %v, want one ending -your-invoice.html", entries)
	}
	content, err := os.ReadFile(dir + string(os.PathSeparator) + entries[0].Name())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "<p>Hello</p>") || !strings.Contains(string(content), "Subject: Your invoice") {
		t.Errorf("the saved message is missing its body or subject:\n%s", content)
	}
}

type stubTransport struct {
	name string
	err  error
	mu   sync.Mutex
	sent []*Message
}

func (s *stubTransport) Name() string { return s.name }

func (s *stubTransport) Send(_ context.Context, m *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err != nil {
		return s.err
	}
	s.sent = append(s.sent, m)
	return nil
}

func TestFailoverTriesTheNextTransport(t *testing.T) {
	bad := &stubTransport{name: "bad", err: errors.New("connection refused")}
	good := &stubTransport{name: "good"}
	f := NewFailover(bad, good)
	expect(t, "name", f.Name(), "failover(bad,good)")
	if err := f.Send(context.Background(), sampleMessage()); err != nil {
		t.Fatal(err)
	}
	if len(good.sent) != 1 {
		t.Errorf("the second transport got %d messages, want 1", len(good.sent))
	}

	worse := &stubTransport{name: "worse", err: errors.New("401 unauthorized")}
	err := NewFailover(bad, worse).Send(context.Background(), sampleMessage())
	if err == nil || !strings.Contains(err.Error(), "connection refused") || !strings.Contains(err.Error(), "401 unauthorized") {
		t.Errorf("every transport failing gave %v, want both failures", err)
	}
}

func TestFromConfigPicksTheDriver(t *testing.T) {
	cases := []struct {
		name    string
		cfg     config.Config
		driver  string
		wantErr string
	}{
		{name: "an existing project with only RESEND_API_KEY keeps Resend", cfg: config.Config{AppEnv: "production", ResendAPIKey: "re_live"}, driver: "resend"},
		{name: "the placeholder key in development sends to Mailhog, then the log", cfg: config.Config{AppEnv: "development", ResendAPIKey: "re_your_api_key"}, driver: "failover(smtp,log)"},
		{name: "the cloud placeholder key is not a key", cfg: config.Config{AppEnv: "development", ResendAPIKey: "re_your_api_key_here"}, driver: "failover(smtp,log)"},
		{name: "nothing configured in production", cfg: config.Config{AppEnv: "production"}, driver: ""},
		{name: "MAIL_MAILER wins over a Resend key", cfg: config.Config{AppEnv: "production", ResendAPIKey: "re_live", Mail: config.MailConfig{Mailer: "SMTP", SMTPHost: "smtp.example.com", SMTPPort: "587"}}, driver: "smtp"},
		{name: "resend without a key", cfg: config.Config{Mail: config.MailConfig{Mailer: "resend"}}, wantErr: "RESEND_API_KEY"},
		{name: "log is refused in production", cfg: config.Config{AppEnv: "production", Mail: config.MailConfig{Mailer: "log"}}, wantErr: "MAIL_ALLOW_LOG_IN_PRODUCTION"},
		{name: "log in production when allowed", cfg: config.Config{AppEnv: "production", Mail: config.MailConfig{Mailer: "log", AllowLogInProduction: true}}, driver: "log"},
		{name: "mailgun without a secret", cfg: config.Config{Mail: config.MailConfig{Mailer: "mailgun", MailgunDomain: "mg.example.com"}}, wantErr: "MAILGUN_SECRET"},
		{name: "mailgun", cfg: config.Config{Mail: config.MailConfig{Mailer: "mailgun", MailgunDomain: "mg.example.com", MailgunSecret: "key"}}, driver: "mailgun"},
		{name: "postmark without a token", cfg: config.Config{Mail: config.MailConfig{Mailer: "postmark"}}, wantErr: "POSTMARK_TOKEN"},
		{name: "postmark", cfg: config.Config{Mail: config.MailConfig{Mailer: "postmark", PostmarkToken: "t"}}, driver: "postmark"},
		{name: "sendgrid without a key", cfg: config.Config{Mail: config.MailConfig{Mailer: "sendgrid"}}, wantErr: "SENDGRID_API_KEY"},
		{name: "sendgrid", cfg: config.Config{Mail: config.MailConfig{Mailer: "sendgrid", SendGridAPIKey: "SG.x"}}, driver: "sendgrid"},
		{name: "ses without credentials", cfg: config.Config{Mail: config.MailConfig{Mailer: "ses", SESRegion: "us-east-1"}}, wantErr: "AWS_ACCESS_KEY_ID"},
		{name: "ses", cfg: config.Config{Mail: config.MailConfig{Mailer: "ses", SESRegion: "us-east-1", SESAccessKeyID: "AKID", SESSecretAccessKey: "s"}}, driver: "ses"},
		{name: "failover without a list", cfg: config.Config{Mail: config.MailConfig{Mailer: "failover"}}, wantErr: "MAIL_FAILOVER"},
		{name: "failover", cfg: config.Config{ResendAPIKey: "re_live", Mail: config.MailConfig{Mailer: "failover", Failover: []string{"resend", " smtp "}}}, driver: "failover(resend,smtp)"},
		{name: "failover cannot carry log into production", cfg: config.Config{AppEnv: "production", Mail: config.MailConfig{Mailer: "failover", Failover: []string{"smtp", "log"}}}, wantErr: "MAIL_ALLOW_LOG_IN_PRODUCTION"},
		{name: "an unknown driver", cfg: config.Config{Mail: config.MailConfig{Mailer: "pigeon"}}, wantErr: "not a mail driver"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tc.cfg
			cfg.MailFrom = "noreply@example.com"
			m, err := FromConfig(&cfg)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want one naming %s", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tc.driver == "" {
				if m != nil {
					t.Fatalf("got a %s mailer, want none", m.Driver())
				}
				return
			}
			if m == nil {
				t.Fatalf("got no mailer, want %s", tc.driver)
			}
			expect(t, "driver", m.Driver(), tc.driver)
		})
	}

	named := config.Config{AppEnv: "development", MailFrom: "noreply@example.com", Mail: config.MailConfig{FromName: "Grit App"}}
	m, err := FromConfig(&named)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, "from with MAIL_FROM_NAME", m.From(), "\"Grit App\" <noreply@example.com>")
}

func TestMailerKeepsItsAPI(t *testing.T) {
	ctx := context.Background()
	rec := &stubTransport{name: "rec"}
	m := NewWithTransport(rec, "noreply@example.com")

	err := m.Send(ctx, SendOptions{
		To:       "ada@example.com",
		Subject:  "Reset your password",
		Template: "password-reset",
		Data:     map[string]interface{}{"AppName": "Grit", "ResetURL": "https://app.example.com/reset?token=abc", "Year": 2026},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.SendRaw(ctx, "bob@example.com", "Your code", "<p>123456</p>"); err != nil {
		t.Fatal(err)
	}
	if len(rec.sent) != 2 {
		t.Fatalf("%d messages, want 2", len(rec.sent))
	}
	expect(t, "template from", rec.sent[0].From, "noreply@example.com")
	if !strings.Contains(rec.sent[0].HTML, "token=abc") {
		t.Errorf("the rendered template lost the reset link:\n%s", rec.sent[0].HTML)
	}
	expect(t, "raw html", rec.sent[1].HTML, "<p>123456</p>")

	if err := m.SendMessage(ctx, &Message{Subject: "nobody", HTML: "<p>x</p>"}); err == nil {
		t.Error("a message with no recipients was sent")
	}
	expect(t, "messages after an invalid one", len(rec.sent), 2)
	expect(t, "New is Resend", New("re_key", "noreply@example.com").Driver(), "resend")
}
`
}

// mailtestGo writes internal/mail/mailtest/mailtest.go: the fake transport for tests.
func mailtestGo() string {
	return `// Package mailtest is a fake mail transport for tests: it keeps every message
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

	appmail "{{MODULE}}/internal/mail"
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
`
}

// mailtestTestGo writes internal/mail/mailtest/mailtest_test.go: tests for the fake.
func mailtestTestGo() string {
	return `package mailtest

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
`
}
