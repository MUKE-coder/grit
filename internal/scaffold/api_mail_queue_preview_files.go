package scaffold

// Queued mail and the template registry (issue 004, B3, B5, B6).
//
// mail.Queue puts a whole Message on the email:send job the worker already
// ran, which nothing enqueued. mail.Register is the list the admin's Mail
// Preview reads, so the page renders the real templates instead of JSX copies
// of them, and grit generate mail registers each email it writes.

// mailQueueGo writes internal/mail/queue.go: Queue and the email:send payload.
func mailQueueGo() string {
	return `package mail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// TaskSend is the background job queued mail travels on. internal/jobs routes
// the same task type to the worker that sends it.
const TaskSend = "email:send"

// MaxQueuedAttachmentBytes caps the attachments a queued message may carry,
// 256 KB in all. The message waits in Redis until the worker sends it, and
// again for every retry, so a large file belongs in storage with a link to it
// in the message.
const MaxQueuedAttachmentBytes = 256 << 10

// ErrAttachmentsTooLarge is returned by Queue for a message whose attachments
// come to more than MaxQueuedAttachmentBytes.
var ErrAttachmentsTooLarge = errors.New("mail: attachments too large to queue")

// Enqueuer is the part of internal/jobs.Client that Queue uses. It is declared
// here because internal/jobs imports this package.
type Enqueuer interface {
	EnqueuePayload(ctx context.Context, taskType string, payload []byte) error
}

// QueuedMessage is the email:send payload Queue writes and the worker reads.
type QueuedMessage struct {
	Message *Message ` + "`" + `json:"message"` + "`" + `
}

// Queue hands msg to the background worker, which sends it with the Mailer
// main.go built and retries a provider that is briefly down. The request that
// queues it does not wait for the provider.
//
// Mail somebody is waiting for, a password reset or a sign-in code, is better
// sent directly with Mailer.SendMessage, so a failure shows while they are
// still there.
func Queue(ctx context.Context, q Enqueuer, msg *Message) error {
	if msg == nil {
		return errors.New("mail: nil message")
	}
	if q == nil {
		return errors.New("mail: no job queue to send through (Redis is not configured); send with Mailer.SendMessage instead")
	}
	total := 0
	for _, a := range msg.Attachments {
		total += len(a.Content)
	}
	if total > MaxQueuedAttachmentBytes {
		return fmt.Errorf("%w: %q has %d bytes of attachments and a queued message may carry %d (256 KB). Store the file with internal/storage and send a link to it instead",
			ErrAttachmentsTooLarge, msg.Subject, total, MaxQueuedAttachmentBytes)
	}
	// Checked now, while the caller can still see the error. The worker would
	// find out after the request had already answered. From may be empty here:
	// the worker's Mailer fills it in.
	check := *msg
	if check.From == "" {
		check.From = "queued@example.com"
	}
	if err := check.Validate(); err != nil {
		return err
	}
	payload, err := json.Marshal(QueuedMessage{Message: msg})
	if err != nil {
		return fmt.Errorf("mail: encoding the queued message: %w", err)
	}
	if err := q.EnqueuePayload(ctx, TaskSend, payload); err != nil {
		return fmt.Errorf("mail: queueing %q: %w", msg.Subject, err)
	}
	return nil
}
`
}

// mailPreviewGo writes internal/mail/preview.go: the template registry, the
// shared layout and the built-in templates' sample data.
func mailPreviewGo() string {
	return `package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"sort"
	"sync"
	"time"
)

// Template is an email the admin's Mail Preview lists. Render builds the
// message with sample data, through the same code that sends the real one.
type Template struct {
	// Name is the :template in /api/admin/mail/preview/:template.
	Name        string
	Description string
	Render      func() (*Message, error)
}

var (
	registryMu sync.RWMutex
	registry   = map[string]Template{}
)

// Register adds a template to the Mail Preview. Each email grit generate mail
// writes in internal/mail/templates registers itself from an init function.
// Registering a name again replaces the first.
func Register(t Template) {
	if t.Name == "" || t.Render == nil {
		panic("mail.Register: a template needs a Name and a Render function")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[t.Name] = t
}

// Templates lists the registered templates by name.
func Templates() []Template {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]Template, 0, len(registry))
	for _, t := range registry {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// LookupTemplate returns the registered template called name.
func LookupTemplate(name string) (Template, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	t, ok := registry[name]
	return t, ok
}

var layout = template.Must(template.New("layout").Parse(baseLayout))

// RenderLayout puts content inside the shared layout: the app name over the
// card, the year and the app name under it. content must already be safe
// HTML, which is what html/template renders. An empty appName uses APP_NAME.
func RenderLayout(appName string, content template.HTML) (string, error) {
	var buf bytes.Buffer
	err := layout.Execute(&buf, map[string]interface{}{
		"AppName": appNameOr(appName),
		"Content": content,
		"Year":    time.Now().Year(),
	})
	if err != nil {
		return "", fmt.Errorf("rendering the mail layout: %w", err)
	}
	return buf.String(), nil
}

func appNameOr(name string) string {
	if name != "" {
		return name
	}
	if env := os.Getenv("APP_NAME"); env != "" {
		return env
	}
	return "App"
}

// builtIn registers one of the templates in EmailTemplates with sample data.
func builtIn(name, description, subject string, sample func() map[string]interface{}) {
	Register(Template{
		Name:        name,
		Description: description,
		Render: func() (*Message, error) {
			data := sample()
			data["AppName"] = appNameOr("")
			data["Year"] = time.Now().Year()
			var m Mailer
			html, err := m.renderTemplate(name, data)
			if err != nil {
				return nil, err
			}
			return &Message{Subject: subject, HTML: html}, nil
		},
	})
}

func init() {
	builtIn("welcome", "Sent when a new user registers", "Welcome", func() map[string]interface{} {
		return map[string]interface{}{"Name": "Ada Lovelace", "DashboardURL": "https://example.com/dashboard"}
	})
	builtIn("password-reset", "Sent when a user asks to reset their password", "Reset your password", func() map[string]interface{} {
		return map[string]interface{}{"ResetURL": "https://example.com/reset-password?token=sample"}
	})
	builtIn("email-verification", "Sent to confirm a user's email address", "Confirm your email address", func() map[string]interface{} {
		return map[string]interface{}{"VerifyURL": "https://example.com/verify-email?token=sample"}
	})
	builtIn("notification", "A general notification with an optional button", "New activity", func() map[string]interface{} {
		return map[string]interface{}{
			"Title":      "New activity",
			"Message":    "Someone commented on your post.",
			"ActionURL":  "https://example.com/activity",
			"ActionText": "View activity",
		}
	})
}
`
}

// mailTemplatesDocGo writes internal/mail/templates/templates.go, the package
// grit generate mail writes into. The preview handler imports it, which is
// what runs each template's Register.
func mailTemplatesDocGo() string {
	return `// Package templates holds the application's own emails. grit generate mail
// writes one file here per email: a typed data struct, an html/template body
// inside the shared layout, a text alternative, and Send and Queue helpers.
//
// Each file registers its email with mail.Register, so the admin's Mail
// Preview lists it. internal/handlers/mail_preview.go imports this package for
// that reason.
package templates
`
}

// mailQueuePreviewTestGo writes internal/mail/queue_test.go.
func mailQueuePreviewTestGo() string {
	return `package mail

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
`
}

// mailPreviewHandlerGo writes internal/handlers/mail_preview.go.
func mailPreviewHandlerGo() string {
	return `package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/mail"
	// For its init functions: every email grit generate mail writes registers
	// itself for this preview.
	_ "{{MODULE}}/internal/mail/templates"
)

// MailPreviewHandler serves the admin's Mail Preview: the registered email
// templates, rendered with sample data by the same code that sends them.
type MailPreviewHandler struct {
	Mailer *mail.Mailer
}

// MailTemplateView is one template in the preview's list.
type MailTemplateView struct {
	Name        string ` + "`" + `json:"name"` + "`" + `
	Description string ` + "`" + `json:"description"` + "`" + `
	Subject     string ` + "`" + `json:"subject"` + "`" + `
	HasText     bool   ` + "`" + `json:"has_text"` + "`" + `
}

// Templates handles GET /api/admin/mail/templates. driver names the transport
// mail goes out through, empty when none is configured.
func (h *MailPreviewHandler) Templates(c *gin.Context) {
	registered := mail.Templates()
	out := make([]MailTemplateView, 0, len(registered))
	for _, t := range registered {
		view := MailTemplateView{Name: t.Name, Description: t.Description}
		if msg, err := t.Render(); err == nil {
			view.Subject = msg.Subject
			view.HasText = msg.Text != ""
		}
		out = append(out, view)
	}
	driver := ""
	if h.Mailer != nil {
		driver = h.Mailer.Driver()
	}
	c.JSON(http.StatusOK, gin.H{"data": out, "driver": driver})
}

// Preview handles GET /api/admin/mail/preview/:template: the rendered HTML, or
// the text part with ?part=text.
func (h *MailPreviewHandler) Preview(c *gin.Context) {
	t, ok := mail.LookupTemplate(c.Param("template"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "No mail template by that name"}})
		return
	}
	msg, err := t.Render()
	if err != nil {
		log.Printf("mail preview: rendering %s: %v", t.Name, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "The template could not be rendered"}})
		return
	}
	// The admin shows this in a sandboxed iframe. Opened on its own, the page
	// still runs no script and loads nothing but images.
	c.Header("Content-Security-Policy", "default-src 'none'; img-src data: https:; style-src 'unsafe-inline'; sandbox")
	c.Header("Cache-Control", "no-store")
	if c.Query("part") == "text" {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(msg.Text))
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg.HTML))
}
`
}

// jobsEnqueuePayloadFunc is the jobs.Client method mail.Queue sends through.
const jobsEnqueuePayloadFunc = `
// EnqueuePayload queues a task whose payload is already encoded, with the
// default options. mail.Queue sends through it: internal/mail cannot import
// this package, so it asks for this one method.
func (c *Client) EnqueuePayload(ctx context.Context, taskType string, payload []byte) error {
	if c == nil || c.client == nil {
		return errors.New("jobs: no queue client, Redis is not configured")
	}
	err := c.Enqueue(ctx, taskType, payload)
	if errors.Is(err, ErrDuplicateTask) {
		return nil
	}
	return err
}
`

// jobsWorkerMailOld is the start of handleEmailSend as Grit wrote it.
const jobsWorkerMailOld = `		if deps.Mailer == nil {
			return fmt.Errorf("mailer not configured")
		}

		var payload EmailPayload
`

const jobsWorkerMailNew = `		if deps.Mailer == nil {
			return fmt.Errorf("mailer not configured")
		}

		// mail.Queue carries the whole message. EnqueueSendEmail carries a
		// template name and its data, handled below.
		var queued mail.QueuedMessage
		if err := json.Unmarshal(task.Payload(), &queued); err != nil {
			return fmt.Errorf("unmarshaling email payload: %v: %w", err, asynq.SkipRetry)
		}
		if queued.Message != nil {
			log.Printf("Sending queued email: %s", queued.Message.Subject)
			return deps.Mailer.SendMessage(ctx, queued.Message)
		}

		var payload EmailPayload
`

// routesMailPreviewAnchor is the staff route the Mail Preview routes follow.
const routesMailPreviewAnchor = "\t\tstaff.GET(\"/admin/cron/tasks\", middleware.RequireRole(\"ADMIN\", \"perm:jobs.view\"), cronHandler.ListTasks)\n"

const routesMailPreview = `		// Mail Preview: the registered email templates, rendered with sample
		// data by the code that sends them.
		mailPreview := &handlers.MailPreviewHandler{Mailer: svc.Mailer}
		staff.GET("/admin/mail/templates", middleware.RequireRole("ADMIN", "perm:system.view"), mailPreview.Templates)
		staff.GET("/admin/mail/preview/:template", middleware.RequireRole("ADMIN", "perm:system.view"), mailPreview.Preview)
`

// ticketMailFuncOpen opens SendTicketCreatedEmail.
const ticketMailFuncOpen = "func SendTicketCreatedEmail(m *mail.Mailer, t *models.Ticket, creator *models.User) error {\n"

const ticketMailSettingCheck = `	// notifications.email_enabled is the admin's switch for notification mail
	// (Settings, Notifications). Off, the ticket is still saved and admins
	// still get the in-app notification.
	if !settings.Bool(context.Background(), "notifications.email_enabled") {
		log.Printf("ticket-mail: notification emails are turned off in settings, skipping ticket %s", t.ID)
		return nil
	}
`
