package mail

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
