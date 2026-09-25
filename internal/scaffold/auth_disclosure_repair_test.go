package scaffold

import (
	"encoding/base64"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var envKeyLine = regexp.MustCompile(`(?m)^FIELD_ENCRYPTION_KEY=(.*)$`)

// Every new project gets its own key, and .env.example never carries it.
func TestNewProjectsGetAFieldEncryptionKey(t *testing.T) {
	for _, opts := range []Options{
		{ProjectName: "triple", Architecture: ArchTriple},
		{ProjectName: "single", Architecture: ArchSingle},
		{ProjectName: "api", Architecture: ArchAPI},
	} {
		env := envFile(opts)
		m := envKeyLine.FindStringSubmatch(env)
		if m == nil {
			t.Fatalf("%s: .env has no uncommented FIELD_ENCRYPTION_KEY line", opts.ProjectName)
		}
		raw, err := base64.StdEncoding.DecodeString(m[1])
		if len(m[1]) != 44 || err != nil || len(raw) != 32 {
			t.Errorf("%s: FIELD_ENCRYPTION_KEY is not 32 bytes of base64 (%d chars, %v)", opts.ProjectName, len(m[1]), err)
		}
		example := envKeyLine.FindStringSubmatch(envExampleFile(opts))
		if example == nil || example[1] != "CHANGE_ME" {
			t.Errorf("%s: .env.example does not carry the placeholder: %v", opts.ProjectName, example)
		}
		if again := envKeyLine.FindStringSubmatch(envFile(opts)); again[1] == m[1] {
			t.Errorf("%s: two scaffolds got the same key", opts.ProjectName)
		}
		if !strings.Contains(env, fieldEncryptionKeyComment) || strings.Contains(fieldEncryptionKeyComment, "—") {
			t.Errorf("%s: the key's comment is missing or has an em dash", opts.ProjectName)
		}
	}
	if m := envKeyLine.FindStringSubmatch(envCloudExampleFile(Options{ProjectName: "x"})); m == nil || m[1] != "CHANGE_ME" {
		t.Error("the cloud example should carry a placeholder key")
	}
	if !strings.Contains(apiCryptoFieldGo(), `strings.EqualFold(b64, "CHANGE_ME")`) {
		t.Error("the API should name the placeholder when .env.example is copied unchanged")
	}
}

// An upgrade advises a key and never writes one.
func TestUpgradeAdvisesAFieldEncryptionKey(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchTriple}
	model := filepath.Join(opts.APIRoot(root), "internal", "models", "two_factor.go")
	if err := os.MkdirAll(filepath.Dir(model), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(model, []byte("package models\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	env := filepath.Join(root, ".env")
	for _, c := range []struct {
		body   string
		advise bool
	}{
		{"JWT_SECRET=x\n# FIELD_ENCRYPTION_KEY=\n", true},
		{"FIELD_ENCRYPTION_KEY=\n", true},
		{"FIELD_ENCRYPTION_KEY=" + base64.StdEncoding.EncodeToString(make([]byte, 32)) + "\n", false},
	} {
		if err := os.WriteFile(env, []byte(c.body), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := fieldKeyAdviceFor(root, opts) != ""; got != c.advise {
			t.Errorf("%q: advice %v, want %v", c.body, got, c.advise)
		}
		if b, _ := os.ReadFile(env); string(b) != c.body {
			t.Errorf("the advice changed .env")
		}
	}
	if err := os.Remove(model); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(env, []byte("JWT_SECRET=x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if fieldKeyAdviceFor(root, opts) != "" {
		t.Error("a project without two-factor sign-in was advised a key")
	}
}

const disclosureModule = "example.com/app"

func formatted(t *testing.T, name, src string) string {
	t.Helper()
	out, err := format.Source([]byte(src))
	if err != nil {
		t.Fatalf("%s is not valid Go: %v", name, err)
	}
	return string(out)
}

// withModule is a template as a project holds it.
func withModule(src string) string {
	return strings.ReplaceAll(src, "{{MODULE}}", disclosureModule)
}

// L1: nothing about the account is reported before the password matches.
func TestLoginReportsNothingBeforeThePassword(t *testing.T) {
	src := apiAuthHandlerGo()
	for _, want := range []string{
		"func (h *AuthHandler) checkCredentials(",
		"spendPasswordCheck(req.Password)",
		"if !h.checkCredentials(c, &user, req.Password) {",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("auth.go is missing %q", want)
		}
	}
	for _, gone := range []string{`"SOCIAL_AUTH_ONLY"`, `"ACCOUNT_LOCKED"`, "This account uses %s"} {
		if strings.Contains(src, gone) {
			t.Errorf("auth.go still answers sign-in with %s", gone)
		}
	}
	// The password is checked before any refusal that names the account's state.
	if strings.Index(src, "h.checkCredentials(c, &user, req.Password)") > strings.Index(src, "h.loginRefusal(&user)") {
		t.Error("Login reports a refusal before it checks the password")
	}
	mustFormatGo(t, "handlers/auth.go", src)
}

func TestLoginDisclosureRepair(t *testing.T) {
	fresh := withModule(apiAuthHandlerGo())
	respondImport := "\t\"" + disclosureModule + "/internal/respond\"\n"
	// An auth.go as each release wrote it: the Login block, fmt for the
	// refusal messages, and neither sync nor bcrypt.
	shape := func(block string) string {
		old := strings.Replace(fresh, loginBlockNew, block, 1)
		old = strings.Replace(old, "\t\"errors\"\n", "\t\"errors\"\n\t\"fmt\"\n", 1)
		old = strings.Replace(old, "\t\"sync\"\n", "", 1)
		return strings.Replace(old, "\t\"golang.org/x/crypto/bcrypt\"\n", "", 1)
	}
	for _, c := range []struct{ name, old string }{
		{"3.283.0", shape(loginBlockOld283)},
		{"3.286.0", shape(loginBlockOld286)},
	} {
		if c.old == fresh || !strings.Contains(c.old, "SOCIAL_AUTH_ONLY") {
			t.Fatalf("%s: could not reconstruct the old auth.go", c.name)
		}
		formatted(t, c.name+" auth.go", c.old)

		out, fixed, warn := repairLoginDisclosureSource(c.old)
		if len(fixed) != 1 || len(warn) != 0 {
			t.Fatalf("%s: auth.go was not repaired: %v %v", c.name, fixed, warn)
		}
		if formatted(t, "repaired auth.go", out) != formatted(t, "auth.go", fresh) {
			t.Errorf("%s: the repaired auth.go is not the template:\n%s", c.name, out)
		}
		if again, fixed, _ := repairLoginDisclosureSource(out); again != out || len(fixed) != 0 {
			t.Errorf("%s: the Login repair is not idempotent", c.name)
		}
		edited := strings.Replace(c.old, "Invalid email or password", "Wrong details", 1)
		if out, _, warn := repairLoginDisclosureSource(edited); out != edited || len(warn) != 1 {
			t.Errorf("%s: an edited Login was changed, or not warned about", c.name)
		}
	}

	// Before 3.285.0 auth.go answered by hand and did not import respond, which
	// the new block uses.
	bare := strings.Replace(shape(loginBlockOld283), respondImport, "", 1)
	out, fixed, _ := repairLoginDisclosureSource(bare)
	if len(fixed) != 1 || !strings.Contains(out, respondImport) {
		t.Errorf("a 3.283.0 auth.go without the respond import did not get it:\n%s", out)
	}
	formatted(t, "repaired 3.283.0 auth.go", out)
}

// L2 and L3: the secret is encrypted at rest, and the trusted-device cookie is
// Secure on HTTPS and SameSite=Lax.
func TestTwoFactorSecretAndCookie(t *testing.T) {
	model := twoFactorModelsGo()
	if !strings.Contains(model, "Secret crypto.EncryptedString") || !strings.Contains(model, "/internal/crypto\"") {
		t.Error("TwoFactorConfig.Secret is not an EncryptedString")
	}
	mustFormatGo(t, "models/two_factor.go", model)

	handler := totpHandlerGo()
	for _, want := range []string{
		"config.Secret = crypto.EncryptedString(req.Secret)",
		// Validated against the device's own clock offset since v3.328.0; what
		// matters here is that the secret is decrypted at the call.
		"totp.ValidateCodeOffset(string(config.Secret), req.Code, config.StepOffset)",
		"h.sealSecret(c, config)",
		"c.SetSameSite(http.SameSiteLaxMode)",
		`c.SetCookie("totp_trusted", deviceToken, int(totp.TrustedDeviceDuration.Seconds()), "/", "", secure, true)`,
	} {
		if !strings.Contains(handler, want) {
			t.Errorf("handlers/totp.go is missing %q", want)
		}
	}
	if strings.Contains(handler, "false, // secure") {
		t.Error("the trusted-device cookie is still never Secure")
	}
	if strings.Count(handler, "h.sealSecret(c, config)") != 2 {
		t.Error("a secret should be sealed after both a code and a backup code verify")
	}
	mustFormatGo(t, "handlers/totp.go", handler)

	if !strings.Contains(apiCryptoFieldGo(), "func EncryptExisting(db *gorm.DB, models ...interface{}) (int64, error) {") {
		t.Error("internal/crypto has no EncryptExisting")
	}
	mustFormatGo(t, "crypto/field.go", apiCryptoFieldGo())
	mustFormatGo(t, "crypto/map_update_test.go", apiCryptoMapUpdateTestGo())
	migrate := apiMigrateMainGo()
	if !strings.Contains(migrate, "crypto.EncryptExisting(db, models.Models()...)") {
		t.Error("grit migrate does not encrypt existing plaintext values")
	}
	mustFormatGo(t, "cmd/migrate/main.go", migrate)
}

// L4: a refused handshake token is not explained to the caller.
func TestRealtimeTokenErrorIsNotEchoed(t *testing.T) {
	fresh := withModule(apiRealtimeHandlerGo())
	if strings.Contains(fresh, `"message": err.Error()`) || !strings.Contains(fresh, realtimeInvalidTokenNew) {
		t.Fatal("the realtime handler still sends err.Error() with INVALID_TOKEN")
	}
	if !strings.Contains(fresh, "realtime.CheckOrigin(c.Request)") {
		t.Error("the realtime handler no longer checks the origin of a cookie handshake")
	}
	respondImport := "\t\"" + disclosureModule + "/internal/respond\"\n"
	for _, c := range []struct{ name, block string }{
		{"3.283.0", realtimeInvalidTokenOld283},
		{"3.286.0", realtimeInvalidTokenOld286},
	} {
		// Neither release imported respond in the realtime handler.
		old := strings.Replace(fresh, realtimeInvalidTokenNew, c.block, 1)
		old = strings.Replace(old, respondImport, "", 1)
		formatted(t, c.name+" realtime.go", old)
		out, fixed, warn := repairRealtimeTokenErrorSource(old)
		if len(fixed) != 1 || len(warn) != 0 || out != fresh {
			t.Fatalf("%s: the realtime handler was not repaired to the template (%v %v):\n%s", c.name, fixed, warn, out)
		}
		if again, fixed, _ := repairRealtimeTokenErrorSource(out); again != out || len(fixed) != 0 {
			t.Errorf("%s: the realtime repair is not idempotent", c.name)
		}
	}
}

// L5: an import job records who started it, and only they or an admin read it.
func TestImportJobOwnerRepair(t *testing.T) {
	model := withModule(importJobModelGo())
	handler := withModule(importJobHandlerGo())
	if !strings.Contains(model, "CreatedBy string") || !strings.Contains(handler, `"created_by = ? AND created_by <> ''"`) {
		t.Fatal("the import job is not scoped to whoever started it")
	}

	oldModel := strings.Replace(model, importJobCreatedByField, "", 1)
	out, fixed, warn := repairImportJobModelSource(oldModel)
	if len(fixed) != 1 || len(warn) != 0 || formatted(t, "model", out) != formatted(t, "model", model) {
		t.Fatalf("the ImportJob model was not repaired to the template (%v %v):\n%s", fixed, warn, out)
	}
	if again, _, _ := repairImportJobModelSource(out); again != out {
		t.Error("the model repair is not idempotent")
	}

	// 3.286.0 answers the missing job through respond.Fail; 3.283.0 by hand,
	// without the respond import. The repair changes only the query.
	authzImport := "\t\"" + disclosureModule + "/internal/authz\"\n"
	respondImport := "\t\"" + disclosureModule + "/internal/respond\"\n"
	failLine := "\t\trespond.Fail(c, respond.CodeNotFound, \"Import job not found\")\n"
	handBuilt := "\t\tc.JSON(http.StatusNotFound, gin.H{\n\t\t\t\"error\": gin.H{\"code\": \"NOT_FOUND\", \"message\": \"Import job not found\"},\n\t\t})\n"
	if !strings.Contains(handler, importJobGetNew+failLine) {
		t.Fatal("the import job handler template no longer answers a missing job through respond.Fail")
	}
	handler283 := strings.Replace(strings.Replace(handler, failLine, handBuilt, 1), respondImport, "", 1)
	for _, c := range []struct{ name, want string }{
		{"3.286.0", handler},
		{"3.283.0", handler283},
	} {
		old := strings.Replace(c.want, importJobGetNew, importJobGetOld, 1)
		old = strings.Replace(old, authzImport, "", 1)
		out, fixed, warn := repairImportJobHandlerSource(old, disclosureModule)
		if len(fixed) != 1 || len(warn) != 0 || formatted(t, "handler", out) != formatted(t, "handler", c.want) {
			t.Fatalf("%s: the import job handler was not repaired (%v %v):\n%s", c.name, fixed, warn, out)
		}
		if again, _, _ := repairImportJobHandlerSource(out, disclosureModule); again != out {
			t.Errorf("%s: the handler repair is not idempotent", c.name)
		}
	}

	importer := "package services\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\n\t\"gorm.io/gorm/clause\"\n\n\t\"" + disclosureModule + "/internal/imports\"\n\t\"" + disclosureModule + "/internal/models\"\n)\n\n" +
		"var _ = clause.OnConflict{}\nvar _ = imports.Wait\n\n" +
		"func (s *ProductService) StartImport(ctx context.Context, total int) (*models.ImportJob, error) {\n" +
		"\tjob := models.ImportJob{Resource: \"products\", Status: \"processing\", Total: total}\n" +
		"\tif err := s.db(ctx).Create(&job).Error; err != nil {\n\t\treturn nil, fmt.Errorf(\"starting: %w\", err)\n\t}\n\treturn &job, nil\n}\n"
	out, fixed, warn = repairImportStarterSource(importer, disclosureModule)
	if len(fixed) != 1 || len(warn) != 0 ||
		!strings.Contains(out, "CreatedBy: authz.UserIDFrom(ctx)}") ||
		!strings.Contains(out, "\t\""+disclosureModule+"/internal/authz\"\n\t\""+disclosureModule+"/internal/imports\"\n") {
		t.Fatalf("the importer was not repaired (%v %v):\n%s", fixed, warn, out)
	}
	formatted(t, "importer", out)
	if again, _, _ := repairImportStarterSource(out, disclosureModule); again != out {
		t.Error("the importer repair is not idempotent")
	}
}

// loginContactApp is Login from a real project created long before v3.283.0 and
// upgraded release by release to v3.286.0. It holds no secrets: the function
// only, as that project has it.
const loginContactApp = `func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("email = ?", req.Email).First(&user).Error; err != nil {
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		services.LogLoginFailed(h.DB, c, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	if !user.Active {
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_blocked",
			Severity:     "warn",
			Summary:      "Sign-in blocked for disabled account " + user.Email,
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_DISABLED",
				"message": "Your account has been disabled",
			},
		})
		return
	}

	// Locked accounts are refused before the password is even compared, so a
	// lockout cannot be probed by timing the comparison.
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		remaining := time.Until(*user.LockedUntil).Round(time.Minute)
		if remaining < time.Minute {
			remaining = time.Minute
		}
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_locked",
			Severity:     "warn",
			Summary:      "Sign-in refused: account is temporarily locked",
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_LOCKED",
				"message": fmt.Sprintf("Too many failed attempts. Try again in about %d minute(s), or reset your password.", int(remaining.Minutes())),
			},
		})
		return
	}

	// Opt-in gate. Social and SSO sign-ins are unaffected: the IdP already
	// proved the address, and those paths set EmailVerifiedAt on first login.
	if h.Config.RequireEmailVerification && user.EmailVerifiedAt == nil && user.Password != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "EMAIL_NOT_VERIFIED",
				"message": "Confirm your email address before signing in. Check your inbox for the link.",
			},
		})
		return
	}

	if user.Password == "" {
		provider := user.Provider
		if provider == "" || provider == "local" {
			provider = "social login"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "SOCIAL_AUTH_ONLY",
				"message": fmt.Sprintf("This account uses %s. Please sign in with your social account.", provider),
			},
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		// Wrong password on a real account: distinct from "unknown email"
		// because Sentinel's brute-force heuristics weight these higher.
		services.LogLoginFailed(h.DB, c, req.Email)
		h.registerFailedLogin(&user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Check if user has TOTP enabled
	var totpConfig models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err == nil {
		// TOTP is enabled: check for trusted device
		if !IsTrustedDevice(c, h.DB, user.ID) {
			// Generate a short-lived pending token for TOTP verification
			pendingToken, err := totp.GeneratePendingToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to create verification session"},
				})
				return
			}

			// Store hashed pending token in DB
			h.DB.WithContext(c.Request.Context()).Create(&models.TOTPPendingToken{
				UserID:    user.ID,
				TokenHash: totp.HashToken(pendingToken),
				ExpiresAt: time.Now().Add(totp.PendingTokenExpiry),
			})

			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"totp_required": true,
					"pending_token": pendingToken,
				},
				"message": "Two-factor authentication required",
			})
			return
		}
	}

	// The failure count is cleared when the sign-in completes. Cleared at the
	// password, with 2FA still ahead, it handed a fresh set of guesses at the
	// code to anyone who knew the password and simply signed in again.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error; err != nil {
			log.Printf("lockout: clearing the failure count for %s: %v", user.ID, err)
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate tokens",
			},
		})
		return
	}

	// Set HttpOnly auth cookies for browser clients. Native mobile/desktop
	// clients ignore them and continue to use the Bearer header from the
	// tokens object below: both flows work.
	//
	// Record the refresh token as a server-side session so this device can be
	// listed and revoked later.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: successful sign-in lands in /system/activity at info
	// severity. IP + user-agent come from the request context inside
	// LogLogin so brute-force investigation has the full pair.
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}
`

// loginTemplate240 is Login as the v3.240.0 template wrote it: no request
// context on its queries, and the failure count cleared at the password.
const loginTemplate240 = `func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		services.LogLoginFailed(h.DB, c, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	if !user.Active {
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_blocked",
			Severity:     "warn",
			Summary:      "Sign-in blocked for disabled account " + user.Email,
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_DISABLED",
				"message": "Your account has been disabled",
			},
		})
		return
	}

	// Locked accounts are refused before the password is even compared, so a
	// lockout cannot be probed by timing the comparison.
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		remaining := time.Until(*user.LockedUntil).Round(time.Minute)
		if remaining < time.Minute {
			remaining = time.Minute
		}
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_locked",
			Severity:     "warn",
			Summary:      "Sign-in refused: account is temporarily locked",
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_LOCKED",
				"message": fmt.Sprintf("Too many failed attempts. Try again in about %d minute(s), or reset your password.", int(remaining.Minutes())),
			},
		})
		return
	}

	// Opt-in gate. Social and SSO sign-ins are unaffected: the IdP already
	// proved the address, and those paths set EmailVerifiedAt on first login.
	if h.Config.RequireEmailVerification && user.EmailVerifiedAt == nil && user.Password != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "EMAIL_NOT_VERIFIED",
				"message": "Confirm your email address before signing in. Check your inbox for the link.",
			},
		})
		return
	}

	if user.Password == "" {
		provider := user.Provider
		if provider == "" || provider == "local" {
			provider = "social login"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "SOCIAL_AUTH_ONLY",
				"message": fmt.Sprintf("This account uses %s. Please sign in with your social account.", provider),
			},
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		// Wrong password on a real account: distinct from "unknown email"
		// because Sentinel's brute-force heuristics weight these higher.
		services.LogLoginFailed(h.DB, c, req.Email)
		h.registerFailedLogin(&user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Any successful password check clears the counter, including one that
	// still has 2FA ahead of it: the password was correct, which is what this
	// counter measures.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		h.DB.Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil})
	}

	// Check if user has TOTP enabled
	var totpConfig models.TwoFactorConfig
	if err := h.DB.Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err == nil {
		// TOTP is enabled: check for trusted device
		if !IsTrustedDevice(c, h.DB, user.ID) {
			// Generate a short-lived pending token for TOTP verification
			pendingToken, err := totp.GeneratePendingToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to create verification session"},
				})
				return
			}

			// Store hashed pending token in DB
			h.DB.Create(&models.TOTPPendingToken{
				UserID:    user.ID,
				TokenHash: totp.HashToken(pendingToken),
				ExpiresAt: time.Now().Add(totp.PendingTokenExpiry),
			})

			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"totp_required": true,
					"pending_token": pendingToken,
				},
				"message": "Two-factor authentication required",
			})
			return
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate tokens",
			},
		})
		return
	}

	// Set HttpOnly auth cookies for browser clients. Native mobile/desktop
	// clients ignore them and continue to use the Bearer header from the
	// tokens object below: both flows work.
	//
	// Record the refresh token as a server-side session so this device can be
	// listed and revoked later.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: successful sign-in lands in /system/activity at info
	// severity. IP + user-agent come from the request context inside
	// LogLogin so brute-force investigation has the full pair.
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}
`

// Projects created before v3.283.0 hold one monolithic Login, which the exact
// block anchors do not match. Each known copy is repaired to the template.
func TestMonolithicLoginRepair(t *testing.T) {
	fresh := withModule(apiAuthHandlerGo())
	// An auth.go around a monolithic Login: the comment Grit wrote twice, fmt
	// for the refusal messages, and neither sync nor bcrypt.
	shape := func(login string) string {
		old := strings.Replace(fresh, loginBlockNew,
			"// Login authenticates a user and returns tokens.\n\n// Login authenticates a user and returns tokens.\n"+login+"\n\n", 1)
		old = strings.Replace(old, "\t\"errors\"\n", "\t\"errors\"\n\t\"fmt\"\n", 1)
		old = strings.Replace(old, "\t\"sync\"\n", "", 1)
		return strings.Replace(old, "\t\"golang.org/x/crypto/bcrypt\"\n", "", 1)
	}
	want := canonicalLogin(loginMonolithic)
	for _, c := range []struct{ name, login string }{
		{"v3.282.0 upgraded to v3.286.0", loginMonolithic},
		{"contact-app", loginContactApp},
		{"v3.240.0 template", loginTemplate240},
	} {
		if got := canonicalLogin(c.login); got != want {
			t.Errorf("%s: the canonical Login differs from the known shape", c.name)
			continue
		}
		old := shape(c.login)
		formatted(t, c.name+" auth.go", old)
		out, fixed, warn := repairLoginDisclosureSource(old)
		if len(fixed) != 1 || len(warn) != 0 {
			t.Fatalf("%s: auth.go was not repaired: %v %v", c.name, fixed, warn)
		}
		if formatted(t, "repaired auth.go", out) != formatted(t, "auth.go", fresh) {
			t.Errorf("%s: the repaired auth.go is not the template:\n%s", c.name, out)
		}
		if !strings.Contains(out, "if !h.checkCredentials(c, &user, req.Password) {") || strings.Contains(out, "SOCIAL_AUTH_ONLY") {
			t.Errorf("%s: the repaired Login does not answer once until the password matches", c.name)
		}
		if again, fixed, _ := repairLoginDisclosureSource(out); again != out || len(fixed) != 0 {
			t.Errorf("%s: the repair is not idempotent", c.name)
		}
	}

	// A Login with a real change is left alone, and says why.
	for name, edit := range map[string][2]string{
		"an extra refusal":  {"\tif !user.Active {\n", "\tif user.Role == \"BANNED\" {\n\t\treturn\n\t}\n\n\tif !user.Active {\n"},
		"a changed message": {"\"Your account has been disabled\"", "\"Account suspended, contact support\""},
		"a changed lookup":  {"Where(\"email = ?\", req.Email)", "Where(\"LOWER(email) = LOWER(?)\", req.Email)"},
	} {
		if !strings.Contains(loginContactApp, edit[0]) {
			t.Fatalf("%s: the fixture no longer contains %q", name, edit[0])
		}
		old := shape(strings.Replace(loginContactApp, edit[0], edit[1], 1))
		if out, fixed, warn := repairLoginDisclosureSource(old); out != old || len(fixed) != 0 || len(warn) != 1 {
			t.Errorf("%s: an edited Login was changed, or not warned about", name)
		}
	}

	// So is one beside a helper of the new block's name, which would not compile.
	crowded := shape(loginContactApp) + "\nfunc invalidCredentials(c *gin.Context) {}\n"
	if out, _, warn := repairLoginDisclosureSource(crowded); out != crowded || len(warn) != 1 {
		t.Error("a Login beside a helper named like the new block's was changed")
	}
}
