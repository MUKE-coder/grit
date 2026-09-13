package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairTokenSafety brings a project scaffolded before v3.244.0 up to the fix
// for H5 in the contact-app review.
//
// Access and refresh tokens were the same shape, so a seven-day refresh token
// was accepted as a bearer token and on the WebSocket. The auth middleware never
// read the sessions table, so logging out, signing out everywhere and changing a
// password left every access token working. And the refresh cookie was scoped
// to /api/auth while refresh and logout are mounted under /api/v1/auth.
//
// services/auth.go and services/session.go are framework code and arrive whole,
// behind the manifest guard. Everything that calls into them is the developer's
// and is edited only where it still reads as Grit wrote it, and only once both
// services are the new ones, or the edits would not compile.
func repairTokenSafety(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	services := filepath.Join(apiRoot, "internal", "services")
	authPath := filepath.Join(services, "auth.go")
	sessionPath := filepath.Join(services, "session.go")
	if !fileExists(authPath) || !fileExists(sessionPath) {
		return nil
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for path, content := range map[string]string{
		authPath:    apiAuthServiceGo(),
		sessionPath: apiSessionServiceGo(),
	} {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	if !fileContains(authPath, "func (s *AuthService) ValidateAccessToken(") || !fileContains(sessionPath, "func SessionIDForToken(") {
		fmt.Println("  ⚠ services/auth.go or services/session.go has been edited, so tokens are still not typed or tied to their session.\n" +
			"    See what the new versions change with grit upgrade --diff, take them, and run grit upgrade again.")
		return nil
	}

	steps := []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "middleware", "auth.go"), repairAuthMiddlewareTokenSource},
		{filepath.Join(apiRoot, "internal", "handlers", "realtime.go"), repairRealtimeTokenSource},
		{filepath.Join(apiRoot, "internal", "handlers", "auth.go"), repairRefreshHandlerSource},
		{filepath.Join(apiRoot, "internal", "handlers", "totp.go"), repairTOTPSessionSource},
		{filepath.Join(apiRoot, "internal", "handlers", "impersonate.go"), repairImpersonateSessionSource},
		{filepath.Join(apiRoot, "internal", "routes", "routes.go"), repairAuthServiceWiringSource},
	}
	for _, s := range steps {
		if !fileExists(s.path) {
			continue
		}
		if err := repairSourceFile(root, m, s.path, s.fn); err != nil {
			return err
		}
	}
	return nil
}

func repairAuthMiddlewareTokenSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func Auth(") || strings.Contains(src, "ValidateAccessToken(") {
		return src, nil, nil
	}
	const old = "authService.ValidateToken(token)"
	if !strings.Contains(src, old) {
		return src, nil, []string{"does not validate tokens the way Grit wrote it: call authService.ValidateAccessToken, or refresh tokens pass as access tokens"}
	}
	return strings.Replace(src, old, "authService.ValidateAccessToken(token)", 1),
		[]string{"only access tokens with a live session are accepted"}, nil
}

func repairRealtimeTokenSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "RealtimeHandler") || strings.Contains(src, "ValidateAccessToken(") {
		return src, nil, nil
	}
	const old = "h.Auth.ValidateToken(tokenStr)"
	if !strings.Contains(src, old) {
		return src, nil, []string{"does not validate the handshake token the way Grit wrote it: call h.Auth.ValidateAccessToken"}
	}
	return strings.Replace(src, old, "h.Auth.ValidateAccessToken(tokenStr)", 1),
		[]string{"the WebSocket accepts only access tokens with a live session"}, nil
}

const refreshSessionBlock = `	// Keep the session id through rotation, so the new access token names the
	// row the old one did. A refresh token from before token types carries none,
	// and its row is found by the token instead.
	sessionID := claims.SessionID
	if sessionID == "" {
		sessionID = services.SessionIDForToken(h.DB, refreshToken)
	}
	tokens, err := h.AuthService.GenerateSessionTokenPair(user.ID, user.Email, user.Role, sessionID)
`

func repairRefreshHandlerSource(src string) (string, []string, []string) {
	const decl = "func (h *AuthHandler) Refresh(c *gin.Context) {"
	start := strings.Index(src, decl)
	if start < 0 || strings.Contains(src, "GenerateSessionTokenPair(") {
		return src, nil, nil
	}
	end := strings.Index(src[start+len(decl):], "\nfunc ")
	if end < 0 {
		end = len(src) - start - len(decl)
	}
	end += start + len(decl)
	body := src[start:end]

	const validate = "h.AuthService.ValidateToken(refreshToken)"
	const mint = "\ttokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)\n"
	if strings.Count(body, validate) != 1 || strings.Count(body, mint) != 1 {
		return src, nil, []string{"Refresh is not the handler Grit wrote: validate with ValidateRefreshToken and mint with GenerateSessionTokenPair and the session id, or refreshed tokens name no session and are refused"}
	}
	body = strings.Replace(body, validate, "h.AuthService.ValidateRefreshToken(refreshToken)", 1)
	body = strings.Replace(body, mint, refreshSessionBlock, 1)
	return src[:start] + body + src[end:], []string{"refresh accepts only refresh tokens, and keeps the session id"}, nil
}

const totpSessionBlock = `	// Record the session. An access token names its session, and one whose
	// session was never recorded is refused on its first request.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("totp: failed to record session for %s: %v", user.ID, err)
	}
`

var totpMintRe = regexp.MustCompile(`(?s)\ttokens, err := h\.AuthService\.GenerateTokenPair\(user\.ID, user\.Email, user\.Role\)\n\tif err != nil \{\n.*?\n\t\treturn\n\t\}\n`)

func repairTOTPSessionSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "TOTPHandler") {
		return src, nil, nil
	}
	added := 0
	out := src
	for _, loc := range reverse(totpMintRe.FindAllStringIndex(src, -1)) {
		// Already recorded somewhere later in the same function.
		rest := out[loc[1]:]
		if next := strings.Index(rest, "\nfunc "); next >= 0 {
			rest = rest[:next]
		}
		if strings.Contains(rest, "CreateSession(") {
			continue
		}
		out = out[:loc[1]] + totpSessionBlock + out[loc[1]:]
		added++
	}
	if added == 0 {
		return src, nil, nil
	}
	var ok bool
	if !strings.Contains(out, "\t\"log\"\n") {
		if out, ok = addImportGroup(out, "log"); !ok {
			return src, nil, []string{"could not add the log import, so 2FA sign-in still records no session: add services.CreateSession after each GenerateTokenPair"}
		}
	}
	return out, []string{fmt.Sprintf("2FA sign-in records its session (%d %s)", added, plural(added, "place", "places"))}, nil
}

func reverse(locs [][]int) [][]int {
	for i, j := 0, len(locs)-1; i < j; i, j = i+1, j-1 {
		locs[i], locs[j] = locs[j], locs[i]
	}
	return locs
}

var impersonateMintRe = regexp.MustCompile(`(?s)\tpair, err := h\.Auth\.GenerateTokenPair\((\w+)\.(?:ID|UserID), \w+\.Email, \w+\.Role\)\n\tif err != nil \{\n\t\trespond\.Internal\(c, err\)\n\t\treturn\n\t\}\n`)

func repairImpersonateSessionSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "ImpersonateHandler") || strings.Contains(src, "CreateSession(") {
		return src, nil, nil
	}
	out := impersonateMintRe.ReplaceAllStringFunc(src, func(block string) string {
		owner := impersonateMintRe.FindStringSubmatch(block)[1]
		id := owner + ".ID"
		if owner == "claims" {
			id = "claims.UserID"
		}
		return block + "\tif _, err := services.CreateSession(h.DB, c, " + id + ", pair.RefreshToken); err != nil {\n\t\trespond.Internal(c, err)\n\t\treturn\n\t}\n"
	})
	out = strings.Replace(out, "h.Auth.ValidateToken(adminToken)", "h.Auth.ValidateAccessToken(adminToken)", 1)
	if out == src {
		return src, nil, []string{"is not the handler the plugin wrote: record a session with services.CreateSession after each GenerateTokenPair, or impersonation is refused"}
	}
	return out, []string{"impersonation records its sessions"}, nil
}

var authServiceLiteralRe = regexp.MustCompile(`(?s)(authService := &services\.AuthService\{[^}]*RefreshExpiry:\s*cfg\.JWTRefreshExpiry,[ \t]*\n)([ \t]*)\}\n`)

func repairAuthServiceWiringSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "authService := &services.AuthService{") || strings.Contains(src, "services.RefreshCookiePath") {
		return src, nil, nil
	}
	out := authServiceLiteralRe.ReplaceAllString(src,
		"${1}\t\tDB: db,\n${2}}\n\tservices.RefreshCookiePath = \"/api/\" + APIVersion + \"/auth\"\n")
	if out == src {
		return src, nil, []string{"builds the AuthService differently: give it DB: db and set services.RefreshCookiePath, or sessions are not checked and the refresh cookie misses its routes"}
	}
	return out, []string{"access tokens are checked against their session, and the refresh cookie reaches the auth routes"}, nil
}
