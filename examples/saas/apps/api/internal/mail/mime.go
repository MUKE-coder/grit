package mail

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
