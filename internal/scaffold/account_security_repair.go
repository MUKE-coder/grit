package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M1, M3, M4 and M5, in the files writeAPIFiles owns. The
// new text lives here once: the templates splice these constants in, and
// repairAccountSecurity puts the same text into an existing project, anchored on
// what Grit generated before.

// M1: an idempotent replay is bound to the credential that made the request.
const (
	idempotencyKeyOld = "\t\tcacheKey := \"idem:\" + c.Request.Method + \":\" + c.FullPath() + \":\" + key\n"
	idempotencyKeyNew = `		// The replay is bound to the credential the request carries. A key on its
		// own is not a secret: it turns up in logs and proxies, and the stored body
		// can hold tokens or a new API key. A request with no credential is not
		// replayed at all, and neither are sign-in or API key routes.
		scope := idempotencyScope(c)
		if scope == "" {
			c.Next()
			return
		}
		cacheKey := "idem:" + scope + ":" + c.Request.Method + ":" + c.FullPath() + ":" + key
`
	idempotencyScopeAnchor = "\ntype idempotentResponse struct {\n"
	idempotencyScopeFunc   = `
// idempotencyScope names who is asking, as a hash of the credential on the
// request, or returns "" when the request must not be replayed. It runs before
// authentication, so it cannot use the user id; the credential stands in for it.
func idempotencyScope(c *gin.Context) string {
	path := c.Request.URL.Path
	if strings.Contains(path, "/auth/") || strings.Contains(path, "/api-keys") {
		return ""
	}
	credential := strings.TrimSpace(c.GetHeader("Authorization"))
	if credential == "" {
		credential = strings.TrimSpace(c.GetHeader("X-API-Key"))
	}
	if credential == "" {
		if cookie, err := c.Cookie("grit_access"); err == nil {
			credential = cookie
		}
	}
	if credential == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(credential))
	return hex.EncodeToString(sum[:16])
}
`
)

// M3: a provider sign-in cannot inherit an account someone else registered.
const (
	oauthEmailAnchor = "\t// Find or create user by email\n"
	oauthEmailGuard  = `	// A provider that shares no address cannot be matched to an account.
	if gothUser.Email == "" {
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Your provider did not share an email address."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

`
	oauthLinkAnchor     = "\t} else {\n\t\t// Link OAuth provider to existing account\n"
	oauthUnverifiedLink = `	} else {
		// An account whose address was never confirmed may have been registered by
		// someone else, waiting for the owner to sign in with the provider. The
		// provider has now shown who owns the address, so the unconfirmed password
		// stops working and its sessions end before the provider is linked.
		if user.EmailVerifiedAt == nil {
			verifiedAt := time.Now()
			if err := h.DB.Model(&user).Updates(map[string]interface{}{"password": "", "email_verified_at": verifiedAt}).Error; err != nil {
				log.Printf("oauth: securing unverified account %s: %v", user.ID, err)
				redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Something went wrong."))
				c.Redirect(http.StatusTemporaryRedirect, redirectURL)
				return
			}
			if err := services.RevokeAllUserSessions(h.DB, user.ID, ""); err != nil {
				log.Printf("oauth: revoking sessions of unverified account %s: %v", user.ID, err)
			}
			user.Password = ""
			user.EmailVerifiedAt = &verifiedAt
		}
		// Link OAuth provider to existing account
`
	gothicStoreOld = "\tgothic.Store = sessions.NewCookieStore([]byte(cfg.JWTSecret))\n"
	gothicStoreNew = `	// Gothic keeps the OAuth handshake in a cookie of its own. Its key is derived
	// from the JWT secret rather than being it, so the cookie and the tokens never
	// share a key, and the cookie is HttpOnly, SameSite=Lax, and Secure when the
	// API is served over https.
	gothicKey := sha256.Sum256([]byte("gothic-cookie-store:" + cfg.JWTSecret))
	gothicStore := sessions.NewCookieStore(gothicKey[:])
	gothicStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   strings.HasPrefix(cfg.AppURL, "https://"),
		SameSite: http.SameSiteLaxMode,
	}
	gothic.Store = gothicStore
`
)

// M4: reading another user's record takes users.view.
const (
	userByIDRouteOld   = "\t\tprotected.GET(\"/users/:id\", userHandler.GetByID)\n"
	userListRoute      = "\t\tstaff.GET(\"/users\", middleware.RequireRole(\"ADMIN\", \"perm:users.view\"), userHandler.List)\n"
	userByIDStaffRoute = "\t\t// Reading a user takes the same grant as listing them. Your own record is\n" +
		"\t\t// GET /profile.\n" +
		"\t\tstaff.GET(\"/users/:id\", middleware.RequireRole(\"ADMIN\", \"perm:users.view\"), userHandler.GetByID)\n"
)

// M5: changing the email or the password takes the current password.
const (
	profileRequestAnchor = "\tBio       string `json:\"bio\"`\n}\n"
	profileRequestNew    = "\tBio       string `json:\"bio\"`\n\n" +
		"\t// Required to change the email or the password.\n" +
		"\tCurrentPassword string `json:\"current_password\"`\n}\n"
	profileChecksAnchor = "\tupdates := map[string]interface{}{}\n\tpasswordChanged := false\n"
	profileChecks       = `	// Changing the address or the password takes the current password. Without
	// it a stolen session was enough to take the account for good. An account
	// with no password, one made through a provider, proves the address with a
	// password reset instead.
	emailChanged := req.Email != "" && !strings.EqualFold(req.Email, user.Email)
	if emailChanged || req.Password != "" {
		if user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "NO_PASSWORD", "message": "Set a password with a password reset before changing your email or password"},
			})
			return
		}
		if req.CurrentPassword == "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "Enter your current password to change your email or password",
					"details": gin.H{"current_password": "Required"},
				},
			})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "INVALID_PASSWORD", "message": "Your current password is not correct"},
			})
			return
		}
	}
	if emailChanged {
		var taken int64
		if err := h.DB.Model(&models.User{}).Where("LOWER(email) = LOWER(?) AND id <> ?", req.Email, user.ID).Count(&taken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to check the email address"},
			})
			return
		}
		if taken > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{"code": "EMAIL_EXISTS", "message": "Another account already uses that email address"},
			})
			return
		}
	}

`
	profileEmailOld = "\tif req.Email != \"\" {\n\t\tupdates[\"email\"] = req.Email\n\t}\n"
	profileEmailNew = "\tif emailChanged {\n" +
		"\t\t// A new address is unconfirmed until its owner confirms it.\n" +
		"\t\tupdates[\"email\"] = req.Email\n" +
		"\t\tupdates[\"email_verified_at\"] = nil\n" +
		"\t}\n"
	profileRevokeOld = "\tif passwordChanged {\n\t\tif err := services.RevokeAllUserSessions("
	profileRevokeNew = "\t// A new address ends every other session too: the same stolen-session\n" +
		"\t// takeover, by way of the password reset the new address would receive.\n" +
		"\tif passwordChanged || emailChanged {\n\t\tif err := services.RevokeAllUserSessions("
)

// repairAccountSecurity applies M1, M3, M4 and M5 to an existing project.
func repairAccountSecurity(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	files := []struct {
		path string
		fix  func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "middleware", "idempotency.go"), repairIdempotencySource},
		{filepath.Join(apiRoot, "internal", "handlers", "auth_oauth.go"), repairOAuthLinkSource},
		{filepath.Join(apiRoot, "cmd", "server", "main.go"), repairGothicStoreSource},
		{filepath.Join(apiRoot, "internal", "routes", "routes.go"), repairUserByIDRouteSource},
		{filepath.Join(apiRoot, "internal", "handlers", "user.go"), repairProfileChangeSource},
	}
	for _, f := range files {
		if !fileExists(f.path) {
			continue
		}
		if err := repairSourceFile(root, m, f.path, f.fix); err != nil {
			return err
		}
	}
	return nil
}

func withImports(src string, paths ...string) (string, bool) {
	for _, p := range paths {
		var ok bool
		if src, ok = addImportGroup(src, p); !ok {
			return src, false
		}
	}
	return src, true
}

func repairIdempotencySource(src string) (string, []string, []string) {
	if strings.Contains(src, "idempotencyScope(c)") {
		return src, nil, nil
	}
	if strings.Count(src, idempotencyKeyOld) != 1 || strings.Count(src, idempotencyScopeAnchor) != 1 {
		return src, nil, []string{"idempotency.go is not the file Grit wrote: key the cache on the caller's credential as well as the method, path and Idempotency-Key, and skip /auth/ and /api-keys, or one user's key replays another user's response"}
	}
	out := strings.Replace(src, idempotencyKeyOld, idempotencyKeyNew, 1)
	out = strings.Replace(out, idempotencyScopeAnchor, idempotencyScopeFunc+idempotencyScopeAnchor, 1)
	out, ok := withImports(out, "crypto/sha256", "encoding/hex", "strings")
	if !ok {
		return src, nil, []string{"could not add imports to idempotency.go"}
	}
	return out, []string{"idempotent replays are bound to the caller's credential and never cover sign-in or API keys"}, nil
}

func repairOAuthLinkSource(src string) (string, []string, []string) {
	if strings.Contains(src, "user.EmailVerifiedAt == nil") {
		return src, nil, nil
	}
	if strings.Count(src, oauthEmailAnchor) != 1 || strings.Count(src, oauthLinkAnchor) != 1 {
		return src, nil, []string{"auth_oauth.go is not the file Grit wrote: before linking a provider to an account whose email is unverified, clear its password and revoke its sessions, or an account registered in advance with the victim's address keeps password access"}
	}
	out := strings.Replace(src, oauthEmailAnchor, oauthEmailGuard+oauthEmailAnchor, 1)
	out = strings.Replace(out, oauthLinkAnchor, oauthUnverifiedLink, 1)
	return out, []string{"a provider sign-in no longer leaves a pre-registered password working on the account"}, nil
}

func repairGothicStoreSource(src string) (string, []string, []string) {
	if !strings.Contains(src, gothicStoreOld) {
		return src, nil, nil
	}
	out := strings.Replace(src, gothicStoreOld, gothicStoreNew, 1)
	out, ok := withImports(out, "crypto/sha256", "net/http", "strings")
	if !ok {
		return src, nil, []string{"could not add imports to main.go"}
	}
	return out, []string{"the OAuth cookie store has its own key and a HttpOnly, SameSite=Lax cookie"}, nil
}

func repairUserByIDRouteSource(src string) (string, []string, []string) {
	if !strings.Contains(src, userByIDRouteOld) {
		return src, nil, nil
	}
	if strings.Count(src, userListRoute) != 1 {
		return src, nil, []string{"routes.go is not the file Grit wrote: move GET /users/:id behind middleware.RequireRole(\"ADMIN\", \"perm:users.view\"), or any signed-in user can read any user's record"}
	}
	out := strings.Replace(src, userByIDRouteOld, "", 1)
	out = strings.Replace(out, userListRoute, userListRoute+userByIDStaffRoute, 1)
	return out, []string{"reading another user's record takes users.view"}, nil
}

func repairProfileChangeSource(src string) (string, []string, []string) {
	const head = "func (h *UserHandler) UpdateProfile(c *gin.Context) {\n"
	if strings.Contains(src, "CurrentPassword string") || !strings.Contains(src, head) {
		return src, nil, nil
	}
	start := strings.Index(src, head)
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		return src, nil, nil
	}
	body := src[start : start+end+3]
	warn := []string{"user.go is not the file Grit wrote: make UpdateProfile require current_password to change the email or the password, check the new email is free, and reset email_verified_at, or a stolen session can take the account"}
	if strings.Count(src, profileRequestAnchor) != 1 || strings.Count(body, profileChecksAnchor) != 1 ||
		strings.Count(body, profileEmailOld) != 1 || strings.Count(body, profileRevokeOld) != 1 {
		return src, nil, warn
	}
	newBody := strings.Replace(body, profileChecksAnchor, profileChecks+profileChecksAnchor, 1)
	newBody = strings.Replace(newBody, profileEmailOld, profileEmailNew, 1)
	newBody = strings.Replace(newBody, profileRevokeOld, profileRevokeNew, 1)
	out := src[:start] + newBody + src[start+len(body):]
	out = strings.Replace(out, profileRequestAnchor, profileRequestNew, 1)
	out, ok := withImports(out, "strings")
	if !ok {
		return src, nil, []string{"could not add the strings import to user.go"}
	}
	return out, []string{"changing the email or password takes the current password, and a new email is checked and unconfirmed"}, nil
}
