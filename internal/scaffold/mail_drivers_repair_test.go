package scaffold

import (
	"strings"
	"testing"
)

func TestMailDriversInTemplates(t *testing.T) {
	opts := Options{ProjectName: "app", Architecture: ArchDouble}

	for name, main := range map[string]string{
		"cmd/server/main.go": apiMainGo(opts),
		"single main.go":     singleMainGo(Options{ProjectName: "app", Architecture: ArchSingle}),
	} {
		if strings.Contains(main, "mail.New(") || !strings.Contains(main, mailerInitNew) {
			t.Errorf("%s still builds the mailer with mail.New instead of mail.FromConfig", name)
		}
	}

	cfg := apiConfigGo()
	for _, want := range []string{configMailFieldsNew, configMailLoadNew, configMailFuncs} {
		if !strings.Contains(cfg, want) {
			t.Errorf("config.go is missing:\n%s", want)
		}
	}
	mustFormatGo(t, "config.go", cfg)

	routes := apiRoutesGo()
	if strings.Contains(routes, "cfg.ResendAPIKey != \"\"") || !strings.Contains(routes, healthMailStatusNew) || !strings.Contains(routes, healthDriverFieldNew) {
		t.Error("/api/health still reports email from RESEND_API_KEY rather than the mailer's driver")
	}

	env := envExampleFile(opts)
	if !strings.Contains(env, envMailHeadNew+"MAIL_FROM=noreply@app.dev\n"+envMailDriversNew+envSupportAnchor) {
		t.Error(".env.example does not document MAIL_MAILER and the driver keys")
	}

	files := mailPackageFiles()
	for _, required := range []string{"mailer.go", "transport.go", "from_config.go", "smtp.go", "resend.go", "mailgun.go", "postmark.go", "sendgrid.go", "ses.go", "log_transport.go", "failover.go", "mailtest/mailtest.go"} {
		if _, ok := files[required]; !ok {
			t.Errorf("internal/mail/%s is not written", required)
		}
	}
	for rel, content := range files {
		if strings.Contains(content, "—") {
			t.Errorf("internal/mail/%s contains an em dash", rel)
		}
		mustFormatGo(t, rel, strings.ReplaceAll(content, "{{MODULE}}", "example.com/app"))
	}
	for _, api := range []string{"func New(apiKey, from string) *Mailer", "func (m *Mailer) Send(ctx context.Context, opts SendOptions) error", "func (m *Mailer) SendRaw(ctx context.Context, to, subject, htmlBody string) error", "func (m *Mailer) SendMessage(ctx context.Context, msg *Message) error"} {
		if !strings.Contains(files["mailer.go"], api) {
			t.Errorf("mailer.go no longer has %s, which existing call sites use", api)
		}
	}
}

func TestMailDriverRepairs(t *testing.T) {
	cfg := apiConfigGo()
	oldCfg := strings.Replace(cfg, configMailFieldsNew, configMailFieldsOld, 1)
	oldCfg = strings.Replace(oldCfg, configMailLoadNew, configMailLoadOld, 1)
	oldCfg = strings.Replace(oldCfg, "\n"+configMailFuncs, "", 1)
	if strings.Contains(oldCfg, "MailConfig") {
		t.Fatal("could not rebuild the old config.go")
	}
	out, changes, warnings := repairMailConfigSource(oldCfg)
	if out != cfg || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("repairing the old config.go did not give the template (changes %v, warnings %v)", changes, warnings)
	}
	if again, changes, warnings := repairMailConfigSource(out); again != out || len(changes) != 0 || len(warnings) != 0 {
		t.Error("the config.go repair is not idempotent")
	}
	if _, _, warnings := repairMailConfigSource("package config\n\ntype Config struct{}\n"); len(warnings) != 1 {
		t.Error("an edited config.go should get a note")
	}

	for name, tc := range map[string]struct{ fresh, old string }{
		"api":    {apiMainGo(Options{ProjectName: "app", Architecture: ArchDouble}), mailerInitOld},
		"single": {singleMainGo(Options{ProjectName: "app", Architecture: ArchSingle}), mailerInitOldSingle},
	} {
		oldMain := strings.Replace(tc.fresh, mailerInitNew, tc.old, 1)
		if oldMain == tc.fresh {
			t.Fatalf("%s: could not rebuild the old main.go", name)
		}
		out, changes, warnings := repairMailerInitSource(oldMain)
		if out != tc.fresh || len(changes) != 1 || len(warnings) != 0 {
			t.Errorf("%s: repairing the old main.go did not give the template (changes %v, warnings %v)", name, changes, warnings)
		}
		if again, changes, _ := repairMailerInitSource(out); again != out || len(changes) != 0 {
			t.Errorf("%s: the main.go repair is not idempotent", name)
		}
	}
	if _, _, warnings := repairMailerInitSource("package main\n\nfunc main() { m := mail.New(key, from) }\n"); len(warnings) != 1 {
		t.Error("an edited main.go should get a note")
	}

	routes := apiRoutesGo()
	oldRoutes := strings.Replace(routes, healthMailStatusNew, healthMailStatusOld, 1)
	oldRoutes = strings.Replace(oldRoutes, healthDriverFieldNew, healthDriverFieldOld, 1)
	out, changes, warnings = repairMailHealthSource(oldRoutes)
	if out != routes || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("repairing the old routes.go did not give the template (changes %v, warnings %v)", changes, warnings)
	}
	if again, changes, _ := repairMailHealthSource(out); again != out || len(changes) != 0 {
		t.Error("the routes.go repair is not idempotent")
	}

	env := envExampleFile(Options{ProjectName: "app", Architecture: ArchDouble})
	oldEnv := strings.Replace(env, envMailHeadNew, envMailHeadOld, 1)
	oldEnv = strings.Replace(oldEnv, "\n"+envMailDriversNew+envSupportAnchor, "\n"+envSupportAnchor, 1)
	if strings.Contains(oldEnv, "MAIL_MAILER") {
		t.Fatal("could not rebuild the old .env.example")
	}
	out, changes, _ = repairEnvMailSource(oldEnv)
	if out != env || len(changes) != 1 {
		t.Errorf("repairing the old .env.example did not give the template (changes %v)", changes)
	}
	if again, changes, _ := repairEnvMailSource(out); again != out || len(changes) != 0 {
		t.Error("the .env.example repair is not idempotent")
	}
}
