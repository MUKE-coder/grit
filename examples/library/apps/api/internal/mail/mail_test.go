package mail

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

	"library/apps/api/internal/config"
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
		var ok bool
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
