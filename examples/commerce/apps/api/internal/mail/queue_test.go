package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"strings"
	"testing"
)

type recordingQueue struct {
	taskType string
	payload  []byte
	calls    int
}

func (q *recordingQueue) EnqueuePayload(_ context.Context, taskType string, payload []byte) error {
	q.calls++
	q.taskType = taskType
	q.payload = payload
	return nil
}

func TestQueueCarriesTheMessage(t *testing.T) {
	q := &recordingQueue{}
	msg := &Message{
		To:          []string{"ada@example.com"},
		Subject:     "Your invoice",
		HTML:        "<p>Attached</p>",
		Text:        "Attached",
		Attachments: []Attachment{{Filename: "invoice.txt", Content: []byte("total 10")}},
	}
	if err := Queue(context.Background(), q, msg); err != nil {
		t.Fatal(err)
	}
	if q.calls != 1 || q.taskType != TaskSend {
		t.Fatalf("enqueued %d tasks of type %q, want one %s", q.calls, q.taskType, TaskSend)
	}
	var got QueuedMessage
	if err := json.Unmarshal(q.payload, &got); err != nil {
		t.Fatal(err)
	}
	if got.Message == nil || got.Message.Subject != "Your invoice" || got.Message.Text != "Attached" ||
		len(got.Message.Attachments) != 1 || !bytes.Equal(got.Message.Attachments[0].Content, []byte("total 10")) {
		t.Errorf("the payload lost part of the message: %+v", got.Message)
	}
}

func TestQueueRefusesLargeAttachmentsAndBadMessages(t *testing.T) {
	q := &recordingQueue{}
	big := &Message{
		To:          []string{"ada@example.com"},
		Subject:     "Report",
		HTML:        "<p>Report</p>",
		Attachments: []Attachment{{Filename: "report.pdf", Content: make([]byte, 300<<10)}},
	}
	err := Queue(context.Background(), q, big)
	if !errors.Is(err, ErrAttachmentsTooLarge) || !strings.Contains(err.Error(), "link") {
		t.Errorf("a 300 KB attachment gave %v, want ErrAttachmentsTooLarge saying to send a link", err)
	}
	if err := Queue(context.Background(), q, &Message{Subject: "nobody", HTML: "<p>x</p>"}); err == nil {
		t.Error("a message with no recipient was queued")
	}
	if err := Queue(context.Background(), nil, &Message{To: []string{"ada@example.com"}, Subject: "x", HTML: "x"}); err == nil {
		t.Error("queueing with no queue did not fail")
	}
	if q.calls != 0 {
		t.Errorf("%d refused messages reached the queue", q.calls)
	}
}

func TestBuiltInTemplatesArePreviewable(t *testing.T) {
	names := map[string]bool{}
	for _, tmpl := range Templates() {
		names[tmpl.Name] = true
		msg, err := tmpl.Render()
		if err != nil {
			t.Errorf("%s: %v", tmpl.Name, err)
			continue
		}
		if msg.Subject == "" || !strings.Contains(msg.HTML, "</html>") {
			t.Errorf("%s rendered without a subject or a whole document", tmpl.Name)
		}
	}
	for _, want := range []string{"welcome", "password-reset", "email-verification", "notification"} {
		if !names[want] {
			t.Errorf("%s is not registered for the preview", want)
		}
	}
	if _, ok := LookupTemplate("no-such-template"); ok {
		t.Error("LookupTemplate found a template that was never registered")
	}
}

func TestRenderLayout(t *testing.T) {
	out, err := RenderLayout("<Acme>", template.HTML("<h1>Shipped</h1>"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<h1>Shipped</h1>") || !strings.Contains(out, "&lt;Acme&gt;") || strings.Contains(out, "<Acme>") {
		t.Errorf("the layout did not keep the content and escape the app name:\n%s", out)
	}
}
