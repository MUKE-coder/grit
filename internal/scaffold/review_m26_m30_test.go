package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Contact-app review M26 to M30.

// M26: the layout carries the blocking theme script and the provider tree has no
// second theme system.
func TestAdminThemeAppliedBeforeFirstPaint(t *testing.T) {
	layout := adminRootLayout(Options{ProjectName: "demo"})
	if !strings.Contains(layout, "<head>") || !strings.Contains(layout, "__html: themeScript") {
		t.Error("the admin root layout does not apply the stored theme from <head>")
	}
	if strings.Contains(adminProviders(), "ThemeProvider") {
		t.Error("providers.tsx still wraps the app in the second theme system")
	}
	if strings.Contains(adminThemeScript, "`") {
		t.Error("the theme script contains a backtick, which breaks the Go raw strings it travels through")
	}
	if strings.Contains(adminDarkModeToggleComponent(), `window.localStorage.getItem("grit-theme-mode") as Mode | null)
      : null);`) {
		t.Error("DarkModeToggle still decides the theme in an effect instead of reading what the script applied")
	}
}

func TestThemeFlashRepair(t *testing.T) {
	oldProviders := "\"use client\";\n\nimport { ThemeProvider } from \"./theme-provider\";\n\nexport function Providers() {\n  return (\n    <ThemeProvider>\n      <div />\n    </ThemeProvider>\n  );\n}\n"
	out, fixed, warn := repairProvidersThemeSource(oldProviders)
	if len(fixed) != 1 || len(warn) != 0 || strings.Contains(out, "ThemeProvider") {
		t.Errorf("providers: fixed %v warn %v\n%s", fixed, warn, out)
	}

	oldLayout := "import \"./globals.css\";\n\nexport default function RootLayout() {\n  return (\n    <html lang=\"en\">\n      <body>\n      </body>\n    </html>\n  );\n}\n"
	out, fixed, _ = repairAdminLayoutThemeSource(oldLayout)
	if len(fixed) != 1 || !strings.Contains(out, "<head>") || !strings.Contains(out, themeScriptMarker) {
		t.Errorf("layout was not given the script:\n%s", out)
	}
	if again, fixed, _ := repairAdminLayoutThemeSource(out); again != out || len(fixed) != 0 {
		t.Error("the layout repair is not idempotent")
	}

	// The navbar is the only importer of the provider, so it has to go first or
	// the provider survives until a second upgrade.
	root := t.TempDir()
	code := filepath.Join(root, "admin")
	for path, body := range map[string]string{
		"components/layout/navbar.tsx":         `import { useTheme } from "@/components/shared/theme-provider";`,
		"components/shared/theme-provider.tsx": "export function ThemeProvider() {}",
		"components/shared/providers.tsx":      "export function Providers() {}",
	} {
		full := filepath.Join(code, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	shape := adminPanelShape{code: code, routes: filepath.Join(code, "app"), alias: "@"}
	if !removeUnreferencedAdminFile(shape, filepath.Join(code, "components", "layout", "navbar.tsx"), "layout/navbar") {
		t.Fatal("the unreferenced navbar was not removed")
	}
	if !removeUnreferencedAdminFile(shape, filepath.Join(code, "components", "shared", "theme-provider.tsx"), "theme-provider") {
		t.Error("the provider was kept after its only importer was removed")
	}
}

// M27: the share is read in a server component, not in a useEffect.
func TestPublicFormReadsTheShareOnTheServer(t *testing.T) {
	page := webPublicFormPage()
	if strings.Contains(page, `"use client"`) || strings.Contains(page, "useEffect") || strings.Contains(page, "axios") {
		t.Error("the public form page is still a client component fetching in an effect")
	}
	if !strings.Contains(page, `cache: "no-store"`) {
		t.Error("the share fetch can be cached, so a disabled link keeps working")
	}
	client := webPublicFormClient()
	if !strings.HasPrefix(client, `"use client"`) || strings.Contains(client, `from "axios"`) {
		t.Error("public-form.tsx is not a client component posting through the app's api client")
	}
	if !strings.Contains(page, `from "./public-form"`) {
		t.Error("the page does not render the client form")
	}
}

// M28: the named helpers carry no status of their own.
func TestRespondHelpersReadTheCatalogue(t *testing.T) {
	src := apiRespondGo()
	for _, stale := range []string{"func fail(", `"BAD_REQUEST"`, `"NOT_FOUND"`, `"VALIDATION_ERROR"`, "http.StatusBadRequest", "http.StatusNotFound"} {
		if strings.Contains(src, stale) {
			t.Errorf("respond.go still hardcodes %s next to the catalogue", stale)
		}
	}
	for _, helper := range []string{"Fail(c, CodeBadRequest", "Fail(c, CodeNotFound", "Fail(c, CodeValidationError"} {
		if !strings.Contains(src, helper) {
			t.Errorf("respond.go is missing %s", helper)
		}
	}
}

// M29: the ticket handler runs no query of its own.
func TestTicketHandlerCallsTheService(t *testing.T) {
	handler := ticketHandlerGo()
	if strings.Contains(handler, "h.DB.WithContext") || strings.Contains(handler, "h.DB.Where") {
		t.Error("the ticket handler still queries the database itself")
	}
	service := ticketServiceGo()
	if strings.Contains(service, `"github.com/gin-gonic/gin"`) {
		t.Error("the ticket service imports gin")
	}
	for _, sig := range []string{
		"func (s *TicketService) Open(ctx context.Context, actor TicketActor",
		"func (s *TicketService) Visible(ctx context.Context, actor TicketActor",
	} {
		if !strings.Contains(service, sig) {
			t.Errorf("the ticket service is missing %s", sig)
		}
	}
	meta := servicesRequestMetaGo()
	if !strings.Contains(meta, "func LogActivityCtx(ctx context.Context, db *gorm.DB, args ActivityArgs)") {
		t.Error("there is no gin-free LogActivity")
	}
}

// M30: no email leaves a handler from a bare goroutine.
func TestNoMailFromBareGoroutines(t *testing.T) {
	for name, src := range map[string]string{
		"auth.go":                    apiAuthHandlerGo(),
		"auth_password_reset.go":     apiAuthPasswordResetGo(),
		"auth_email_verification.go": apiAuthEmailVerificationGo(),
		"ticket.go":                  ticketHandlerGo(),
		"services/ticket.go":         ticketServiceGo(),
	} {
		for _, bare := range []string{"go h.deliverVerificationEmail(", "go h.deliverPasswordReset(", "go h.emitTicketCreated("} {
			if strings.Contains(src, bare) {
				t.Errorf("%s still sends mail from a bare goroutine: %s", name, bare)
			}
		}
	}
}

func TestQueuedAuthMailRepair(t *testing.T) {
	old := "\tif h.Mailer != nil {\n\t\tctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)\n\t\tdefer cancel()\n\t\tif err := h.Mailer.Send(ctx, mail.SendOptions{\n\t\t\tTo: user.Email,\n\t\t}); err != nil {\n\t\t\tlog.Printf(\"x\")\n\t\t}\n\t\treturn\n\t}\n"
	out := replaceMailerSend(old, `"verify:"+user.ID+":"+token`)
	if strings.Contains(out, "h.Mailer.Send(") || !strings.Contains(out, `dispatchMail(ctx, h.Mailer, h.Jobs, "verify:"+user.ID+":"+token, mail.SendOptions{`) {
		t.Errorf("the inline send was not queued:\n%s", out)
	}
	if !strings.Contains(out, "\t\t\tTo: user.Email,\n\t\t}); err != nil {") {
		t.Error("the message between the two anchors was changed")
	}
}
