package mail

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
