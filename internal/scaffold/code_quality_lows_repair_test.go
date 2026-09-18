package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// L29: no handler Grit writes asserts a context value without checking it.
func TestNoUncheckedContextAssertions(t *testing.T) {
	// An assertion to a caller's type (user.(models.User), userID.(string)) is
	// only allowed in the comma-ok form, which cannot panic.
	assertion := regexp.MustCompile(`\.\((string|uint|models\.User|\*models\.User)\)`)
	for name, src := range map[string]string{
		"handlers/dashboard_layout.go": dashboardLayoutHandlerGo(),
		"handlers/upload.go":           uploadHandlerGo(),
		"handlers/totp.go":             totpHandlerGo(),
	} {
		for _, line := range strings.Split(src, "\n") {
			if assertion.MatchString(line) && !strings.Contains(line, ", ok :=") && !strings.Contains(line, ", _ :=") {
				t.Errorf("%s asserts a context value without checking it: %q", name, strings.TrimSpace(line))
			}
		}
	}
	if !strings.Contains(uploadHandlerGo(), uploadCreateGuardNew) {
		t.Error("the upload handler does not check the caller before reading the body")
	}
	if strings.Count(dashboardLayoutHandlerGo(), dashboardUserIDNew) != 2 {
		t.Error("the dashboard layout handler does not read the caller through authz.CurrentUserID")
	}
}

func TestDashboardLayoutUserRepair(t *testing.T) {
	current := dashboardLayoutHandlerGo()
	old := strings.ReplaceAll(current, dashboardUserIDNew, dashboardUserIDOld)
	old = strings.Replace(old, "UserID:          userID,", "UserID:          userID.(string),", 1)
	old = strings.Replace(old, "models.DashboardLayout{UserID: userID}", "models.DashboardLayout{UserID: userID.(string)}", 1)
	old = strings.Replace(old, "\t\"{{MODULE}}/internal/authz\"\n", "", 1)
	if strings.Count(old, "userID.(string)") != 2 {
		t.Fatalf("the reconstruction has %d assertions, want 2", strings.Count(old, "userID.(string)"))
	}
	checkPerfLowRepair(t, "handlers/dashboard_layout.go", current, strings.ReplaceAll(old, "{{MODULE}}", "example.com/app"), repairDashboardLayoutUserSource)
}

func TestUploadUserRepair(t *testing.T) {
	current := uploadHandlerGo()
	old := strings.Replace(current, uploadCreateGuardNew, uploadCreateGuardOld, 1)
	old = strings.Replace(old, "\tupload := models.Upload{\n", uploadUserIDOld, 1)
	old = strings.Replace(old, "\t\tUserID:       userID,\n", "\t\tUserID:       userID.(string),\n", 1)
	checkPerfLowRepair(t, "handlers/upload.go", current, old, repairUploadUserSource)

	// Somebody's own version is reported, not rewritten.
	edited := strings.Replace(old, uploadCreateGuardOld, "func (h *UploadHandler) Create(c *gin.Context) {\n", 1)
	if out, fixed, warn := repairUploadUserSource(edited); out != edited || len(fixed) != 0 || len(warn) != 1 {
		t.Errorf("an edited upload handler was not left alone with a warning (%v %v)", fixed, warn)
	}
}

// L32: the ticket service and handler name statuses, priorities and the label
// cap, and the model declares them.
func TestTicketValuesAreNamed(t *testing.T) {
	for name, src := range map[string]string{
		"handlers/ticket.go": ticketHandlerGo(),
		"services/ticket.go": ticketServiceGo(),
		"models/ticket.go":   ticketModelGo(),
	} {
		for _, literal := range []string{`transitionStatus(c, "`, `status == "`, `t.Status = "`, `t.Priority = "`, `case "`, `, 8)`, "max int"} {
			if strings.Contains(src, literal) {
				t.Errorf("%s still spells %s", name, literal)
			}
		}
		formatPerfLowGo(t, name, src)
	}
	if !strings.Contains(ticketModelGo(), "TicketPriorityCritical = \"critical\"") ||
		!strings.Contains(ticketServiceGo(), "const MaxTicketLabels = 8") {
		t.Error("the ticket constants are not declared")
	}
}

func TestTicketConstantsRepair(t *testing.T) {
	current := ticketModelGo()
	old := strings.Replace(current, ticketConstantsBlock, "", 1)
	old = strings.Replace(old, ticketDefaultsNew, ticketDefaultsOld, 1)
	old = strings.Replace(old, "// comma-separated, capped by the ticket service", "// comma-separated up to 8", 1)
	out, fixed, warn := repairTicketConstantsSource(formatPerfLowGo(t, "models/ticket.go", old))
	if len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("models/ticket.go was not repaired (%v %v)", fixed, warn)
	}
	// The one difference left is a comment the repair does not touch.
	want := strings.Replace(current, "// comma-separated, capped by the ticket service", "// comma-separated up to 8", 1)
	if formatPerfLowGo(t, "models/ticket.go", out) != formatPerfLowGo(t, "models/ticket.go", want) {
		t.Errorf("repairing models/ticket.go did not produce the template:\n%s", out)
	}
	again, fixed, warn := repairTicketConstantsSource(out)
	if again != out || len(fixed) != 0 || len(warn) != 0 {
		t.Error("the ticket constants repair is not idempotent")
	}

	// An edited model still gets the constants the new service needs.
	edited := strings.Replace(old, ticketDefaultsOld, "\tt.Status = \"new\"\n", 1)
	out, fixed, _ = repairTicketConstantsSource(formatPerfLowGo(t, "models/ticket.go", edited))
	if len(fixed) != 1 || !strings.Contains(out, "TicketStatusClosed") || !strings.Contains(out, "t.Status = \"new\"") {
		t.Error("an edited ticket model did not get the constants, or lost its edit")
	}
}

// L30: nothing kept alive by a blank identifier, and no doc comment twice.
func TestNoLeftoverBlankUses(t *testing.T) {
	for name, src := range map[string]string{
		"handlers/webhooks.go": apiWebhooksHandlerGo(),
		"handlers/ocsf.go":     apiOCSFHandlerGo(),
		"handlers/sync.go":     apiSyncHandlerGo(),
		"handlers/auth.go":     apiAuthHandlerGo(),
	} {
		for _, leftover := range []string{"var _ = context.Background", "var _ = fmt.Sprint", "_ = since", "func getTimeField("} {
			if strings.Contains(src, leftover) {
				t.Errorf("%s still has %q", name, leftover)
			}
		}
		formatPerfLowGo(t, name, src)
	}
	// Declarations nothing uses, as golangci-lint's unused reported them.
	for name, pair := range map[string][2]string{
		"services/sso.go":     {apiSSOServiceGo(), "type providerSession"},
		"handlers/passkey.go": {passkeyHandlerGo(), "type finishRegistrationRequest"},
		"authz/authz_test.go": {authzTestGo(), "type fakeUserScoped"},
	} {
		if strings.Contains(pair[0], pair[1]) {
			t.Errorf("%s still declares %s, which nothing uses", name, pair[1])
		}
	}
	// A doc comment written twice: a comment block, a blank line, and the same
	// block again directly above a declaration.
	for name, src := range map[string]string{
		"handlers/auth.go":         apiAuthHandlerGo(),
		"handlers/auth_oauth.go":   apiAuthOAuthGo(),
		"handlers/auth_lockout.go": apiAuthLockoutGo(),
	} {
		blocks := strings.Split(src, "\n\n")
		for i := 1; i < len(blocks); i++ {
			prev := blocks[i-1]
			if !strings.HasPrefix(prev, "//") || strings.Contains(prev, "\n\t") {
				continue
			}
			if rest, ok := strings.CutPrefix(blocks[i], prev+"\n"); ok && !strings.HasPrefix(rest, "//") {
				t.Errorf("%s has a doc comment written twice:\n%s", name, prev)
			}
		}
	}
}

// L31: the upload allowlist is built from config, not by init from the
// environment.
func TestUploadAllowlistComesFromConfig(t *testing.T) {
	src := uploadHandlerGo()
	if strings.Contains(src, "func init()") || strings.Contains(src, "os.Getenv") || strings.Contains(src, "AllowedMimeTypes") {
		t.Error("handlers/upload.go still builds its allowlist from the environment at init")
	}
	if !strings.Contains(src, "func UploadMIMEAllowlist(extra []string) map[string]bool {") {
		t.Error("handlers/upload.go has no UploadMIMEAllowlist")
	}
}

// L31, in an existing project: config reads the variable and routes.go hands
// the list to the upload handler, which no longer reads it itself.
func TestUploadAllowlistRepairs(t *testing.T) {
	config := apiConfigGo()
	if !strings.Contains(config, configCORSFieldAnchor+configUploadMIMEField+"\n") ||
		!strings.Contains(config, configCORSLoadAnchor+configUploadMIMELoad+"\n") {
		t.Fatal("config.go does not read UPLOAD_ALLOWED_MIME where the repair puts it")
	}
	oldConfig := strings.Replace(config, configUploadMIMEField+"\n", "", 1)
	oldConfig = strings.Replace(oldConfig, configUploadMIMELoad+"\n", "", 1)
	checkPerfLowRepair(t, "config/config.go", config, oldConfig, repairUploadMIMEConfigSource)

	routes := apiRoutesGo()
	if !strings.Contains(routes, routesUploadHandlerNew) {
		t.Fatal("routes.go does not build the upload handler with its allowlist")
	}
	oldRoutes := strings.Replace(routes, routesUploadHandlerNew, routesUploadHandlerOld, 1)
	if out, fixed, warn := repairUploadMIMERoutesSource(oldRoutes); len(fixed) != 1 || len(warn) != 0 || out != routes {
		t.Errorf("routes.go was not repaired to the template (%v %v)", fixed, warn)
	}
	if again, fixed, warn := repairUploadMIMERoutesSource(routes); again != routes || len(fixed) != 0 || len(warn) != 0 {
		t.Error("the routes repair is not idempotent")
	}
}

// L33: the SSO refusals are errors, and the sentences are chosen in one place.
func TestSSORefusalsAreErrors(t *testing.T) {
	src := apiSSOHandlerGo()
	if regexp.MustCompile(`fmt\.Errorf\("[A-Z]`).MatchString(src) {
		t.Error("handlers/sso.go still returns a sentence as an error")
	}
	if strings.Contains(src, "h.failLogin(c, err.Error())") || strings.Contains(apiSAMLHandlerGo(), "h.failLogin(c, err.Error())") {
		t.Error("an SSO callback still shows an error's text on the login page")
	}
	for _, want := range []string{"errSSOLookupFailed, err)", "errSSOProvisionFailed, err)", "func ssoRefusalMessage(err error, email string) string {"} {
		if !strings.Contains(src, want) {
			t.Errorf("handlers/sso.go is missing %q", want)
		}
	}
}
