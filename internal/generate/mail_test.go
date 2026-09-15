package generate

import (
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A new email was a string constant, a map entry, an untyped data map and a JSX
// copy in the admin. grit generate mail writes one typed, registered file.

func mailProject(t *testing.T) string {
	t.Helper()
	api := filepath.Join(t.TempDir(), "apps", "api")
	writeTestFile(t, filepath.Join(api, "internal", "mail", "preview.go"), "package mail\n\nfunc Register(t Template) {}\n")
	return api
}

func TestGenerateMailWritesARegisteredTemplate(t *testing.T) {
	api := mailProject(t)
	if err := generateMailAt(api, "example.com/app", MailOptions{Name: "OrderShipped"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(api, "internal", "mail", "templates", "order_shipped.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if _, err := parser.ParseFile(token.NewFileSet(), path, src, parser.AllErrors); err != nil {
		t.Fatalf("the generated email is not valid Go: %v\n%s", err, src)
	}
	if formatted, err := format.Source(raw); err != nil || string(formatted) != src {
		t.Errorf("the generated email is not gofmt-clean:\n%s", formatted)
	}
	for _, want := range []string{
		`"example.com/app/internal/mail"`,
		"type OrderShippedData struct",
		`const OrderShippedSubject = "Order shipped"`,
		`htmltemplate.New("order-shipped.html")`,
		`texttemplate.New("order-shipped.txt")`,
		"mail.RenderLayout(data.AppName,",
		`Name:        "order-shipped"`,
		"func RenderOrderShipped(data OrderShippedData) (*mail.Message, error)",
		"func SendOrderShipped(ctx context.Context, m *mail.Mailer, to string, data OrderShippedData) error",
		"func QueueOrderShipped(ctx context.Context, q mail.Enqueuer, to string, data OrderShippedData) error",
		"<p>Hi {{.Name}},</p>",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the generated email is missing %s", want)
		}
	}
	if strings.Contains(src, "—") || strings.Contains(src, "{{P}}") || strings.Contains(src, "{{BT}}") {
		t.Error("the generated email has an em dash or an unreplaced placeholder")
	}
	if _, err := os.Stat(filepath.Join(api, "internal", "mail", "templates", "templates.go")); err != nil {
		t.Errorf("the templates package doc was not written: %v", err)
	}

	err = generateMailAt(api, "example.com/app", MailOptions{Name: "OrderShipped"})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("a second generate gave %v, want a refusal to overwrite", err)
	}
}

func TestGenerateMailRefusesWithoutTheRegistryOrAName(t *testing.T) {
	api := filepath.Join(t.TempDir(), "apps", "api")
	err := generateMailAt(api, "example.com/app", MailOptions{Name: "OrderShipped"})
	if err == nil || !strings.Contains(err.Error(), "grit upgrade") {
		t.Errorf("a project without mail.Register gave %v, want a pointer to grit upgrade", err)
	}
	api = mailProject(t)
	for _, name := range []string{"", "  ", "123"} {
		if err := generateMailAt(api, "example.com/app", MailOptions{Name: name}); err == nil {
			t.Errorf("the name %q was accepted", name)
		}
	}
}
