package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MailOptions describes one `grit generate mail` request.
type MailOptions struct {
	Name string
}

// GenerateMail writes internal/mail/templates/<name>.go: a typed data struct,
// an html/template body inside the shared layout, a text alternative, and
// Render, Send and Queue helpers. The file registers the email for the admin's
// Mail Preview.
//
// A new email used to mean a string constant in internal/mail/templates.go, a
// map entry, a map[string]interface{} of data whose keys nothing checked, and
// a JSX copy in the admin so the preview showed something like it.
func GenerateMail(opts MailOptions) error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	arch, _ := readGritJSON(root)
	apiRoot := root
	if arch != "single" {
		apiRoot = filepath.Join(root, "apps", "api")
	}
	module, err := readModulePath(root, arch)
	if err != nil {
		return err
	}
	return generateMailAt(apiRoot, module, opts)
}

func generateMailAt(apiRoot, module string, opts MailOptions) error {
	if strings.TrimSpace(opts.Name) == "" {
		return fmt.Errorf("an email name is required (e.g. grit generate mail OrderShipped)")
	}
	names := MakeNames(opts.Name)
	if names.Pascal == "" || !isGoIdentifier(names.Pascal) {
		return fmt.Errorf("%q is not a usable email name: use letters and digits, starting with a letter (e.g. OrderShipped)", opts.Name)
	}

	mailDir := filepath.Join(apiRoot, "internal", "mail")
	if !fileHas(filepath.Join(mailDir, "preview.go"), "func Register(") {
		return fmt.Errorf("internal/mail has no template registry (mail.Register), so the email would have nowhere to register: run grit upgrade first")
	}
	path := filepath.Join(mailDir, "templates", names.Snake+".go")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("internal/mail/templates/%s.go already exists. It holds your email, so it is not overwritten", names.Snake)
	}

	fmt.Printf("\n  Generating email: %s (template %q)\n\n", names.Pascal, names.Kebab)

	doc := filepath.Join(mailDir, "templates", "templates.go")
	if _, err := os.Stat(doc); err != nil {
		if err := writeFileWithDirs(doc, mailTemplatesPackageDoc); err != nil {
			return fmt.Errorf("writing the templates package: %w", err)
		}
	}
	if err := writeFileWithDirs(path, mailTemplateFileGo(names, module)); err != nil {
		return fmt.Errorf("writing the email: %w", err)
	}
	fmt.Printf("  ✓ internal/mail/templates/%s.go\n", names.Snake)
	fmt.Printf("  ✓ Registered for the admin's Mail Preview as %q\n", names.Kebab)

	fmt.Printf("\n  ✅ Email generated. Write it in %sHTML and %sText, then send it:\n\n", names.Camel, names.Camel)
	fmt.Printf("      templates.Send%s(ctx, svc.Mailer, user.Email, templates.%sData{Name: user.FirstName})\n", names.Pascal, names.Pascal)
	fmt.Printf("      templates.Queue%s(ctx, svc.Jobs, user.Email, data) // sent by the worker, retried on failure\n\n", names.Pascal)
	return nil
}

func isGoIdentifier(s string) bool {
	for i, r := range s {
		letter := r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		if !letter && (i == 0 || r < '0' || r > '9') {
			return false
		}
	}
	return s != ""
}

func fileHas(path, needle string) bool {
	raw, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(raw), needle)
}

const mailTemplatesPackageDoc = `// Package templates holds the application's own emails. grit generate mail
// writes one file here per email: a typed data struct, an html/template body
// inside the shared layout, a text alternative, and Send and Queue helpers.
//
// Each file registers its email with mail.Register, so the admin's Mail
// Preview lists it. internal/handlers/mail_preview.go imports this package for
// that reason.
package templates
`

// mailTemplateFileGo is the email's own file. It is the developer's from the
// moment it is written, which is why a second generate refuses.
func mailTemplateFileGo(names Names, module string) string {
	words := splitPascal(names.Pascal)
	for i := 1; i < len(words); i++ {
		if strings.ToUpper(words[i]) != words[i] {
			words[i] = strings.ToLower(words[i])
		}
	}
	human := strings.Join(words, " ")
	return strings.NewReplacer(
		"{{P}}", names.Pascal,
		"{{C}}", names.Camel,
		"{{K}}", names.Kebab,
		"{{H}}", human,
		"{{MODULE}}", module,
		"{{BT}}", "`",
	).Replace(`package templates

import (
	"bytes"
	"context"
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"

	"{{MODULE}}/internal/mail"
)

// {{P}}Data is what the {{H}} email shows. Add the fields it needs: the
// compiler then checks every place that sends it.
type {{P}}Data struct {
	// AppName is shown above the card and in the footer. Empty uses APP_NAME.
	AppName   string
	Name      string
	ActionURL string
}

// {{P}}Subject is the subject line.
const {{P}}Subject = "{{H}}"

// {{C}}HTML is the body. It is html/template, so every value is escaped, and
// it goes inside the shared layout, which has the .btn and .code styles.
var {{C}}HTML = htmltemplate.Must(htmltemplate.New("{{K}}.html").Parse({{BT}}<h1>{{H}}</h1>
<p>Hi {{.Name}},</p>
<p>Write the email here.</p>
{{if .ActionURL}}<p style="text-align: center; margin-top: 24px;">
  <a href="{{.ActionURL}}" class="btn">Open</a>
</p>{{end}}{{BT}}))

// {{C}}Text is the plain-text part, for mail clients that do not show HTML.
// Spam filters also score a message without one lower.
var {{C}}Text = texttemplate.Must(texttemplate.New("{{K}}.txt").Parse({{BT}}Hi {{.Name}},

Write the email here.
{{if .ActionURL}}
Open: {{.ActionURL}}
{{end}}{{BT}}))

func init() {
	mail.Register(mail.Template{
		Name:        "{{K}}",
		Description: "{{H}}",
		Render: func() (*mail.Message, error) {
			return Render{{P}}(sample{{P}}())
		},
	})
}

// sample{{P}} is the data the admin's Mail Preview renders the email with.
func sample{{P}}() {{P}}Data {
	return {{P}}Data{Name: "Ada Lovelace", ActionURL: "https://example.com"}
}

// Render{{P}} builds the message without sending it: the subject, the HTML
// inside the shared layout and the text part. The recipients are the caller's.
func Render{{P}}(data {{P}}Data) (*mail.Message, error) {
	var body bytes.Buffer
	if err := {{C}}HTML.Execute(&body, data); err != nil {
		return nil, fmt.Errorf("rendering the {{K}} email: %w", err)
	}
	// #nosec G203 -- body was rendered by html/template, which escaped every value.
	html, err := mail.RenderLayout(data.AppName, htmltemplate.HTML(body.String()))
	if err != nil {
		return nil, err
	}
	var text bytes.Buffer
	if err := {{C}}Text.Execute(&text, data); err != nil {
		return nil, fmt.Errorf("rendering the {{K}} text part: %w", err)
	}
	return &mail.Message{Subject: {{P}}Subject, HTML: html, Text: text.String()}, nil
}

// Send{{P}} renders the email and sends it now, through the Mailer main.go
// built.
func Send{{P}}(ctx context.Context, m *mail.Mailer, to string, data {{P}}Data) error {
	if m == nil {
		return fmt.Errorf("sending the {{K}} email: no mail driver is configured")
	}
	msg, err := Render{{P}}(data)
	if err != nil {
		return err
	}
	msg.To = []string{to}
	return m.SendMessage(ctx, msg)
}

// Queue{{P}} renders the email and hands it to the background worker, which
// sends it and retries a provider that is briefly down. Pass svc.Jobs.
func Queue{{P}}(ctx context.Context, q mail.Enqueuer, to string, data {{P}}Data) error {
	msg, err := Render{{P}}(data)
	if err != nil {
		return err
	}
	msg.To = []string{to}
	return mail.Queue(ctx, q, msg)
}
`)
}
