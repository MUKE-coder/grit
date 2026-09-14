package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Production is the default, and the strict mode (H7 and H9 in the contact-app
// review).
//
//   - APP_ENV unset meant development: no rate limits, the WAF only logging,
//     password-reset links in the log.
//   - Outside development nothing refused a default dashboard password or a
//     short JWT secret.
//   - GORM Studio, a writable SQL console, was on in production, and /docs
//     published a map of every route with a console to call them.
//   - Sentinel's login and register limits, AuthShield and the CSRF bootstrap
//     list were keyed to /api/auth/..., while the router mounts /api/v1, so not
//     one of them ever matched a request.
//
// The snippets below are spliced into the scaffold templates and applied by the
// upgrade repair, so a new project and an upgraded one get the same code.

const configProductionBlock = `	// Whether the API reference at /docs is served in production. It maps every
	// route, admin, backup, SSO and GDPR ones included, with a console to call
	// them, so it is not published unless asked for.
	cfg.APIDocsPublic = getEnv("API_DOCS_PUBLIC", "false") == "true"

	if cfg.AppEnv == "production" {
		// GORM Studio is a browser SQL console over every table. A development
		// .env says GORM_STUDIO_ENABLED=true, and .env files get copied to servers
		// whole, so in production it takes its own switch, and is read-only with
		// no SQL editor even then.
		cfg.GORMStudioEnabled = getEnv("GORM_STUDIO_IN_PRODUCTION", "false") == "true"
		cfg.GORMStudioReadOnly = true
` + sqliteProductionCheck + `	if cfg.AppEnv != "development" {
		if err := checkSecrets(cfg); err != nil {
			return nil, err
		}
	}

`

const configCheckSecretsFunc = `// weakSecrets are values a scaffold, a README or a tutorial once put in these
// settings.
var weakSecrets = map[string]bool{
	"studio": true, "sentinel": true, "pulse": true, "admin": true, "password": true, "secret": true,
	"change-me": true, "change_me": true, "change-me-in-prod": true, "sentinel-secret-change-me": true,
	"minioadmin": true,
}

// checkSecrets refuses to start anywhere but development with a secret that can
// be guessed: a JWT secret short enough to brute-force, which forges any
// account's token, or a dashboard password from a default.
func checkSecrets(cfg *Config) error {
	var weak []string
	if len(cfg.JWTSecret) < 32 || weakSecrets[strings.ToLower(cfg.JWTSecret)] {
		weak = append(weak, "JWT_SECRET")
	}
	for _, s := range []struct {
		name, value string
		used        bool
	}{
		{"GORM_STUDIO_PASSWORD", cfg.GORMStudioPassword, cfg.GORMStudioEnabled},
		{"SENTINEL_PASSWORD", cfg.SentinelPassword, cfg.SentinelEnabled},
		{"SENTINEL_SECRET_KEY", cfg.SentinelSecretKey, cfg.SentinelEnabled},
		{"PULSE_PASSWORD", cfg.PulsePassword, cfg.PulseEnabled},
		// MinIO's root password, when MinIO is the store. S3, R2 and B2 keys are
		// issued by the provider and are not guessable defaults.
		{"MINIO_SECRET_KEY", cfg.Storage.SecretKey, cfg.StorageDriver == "minio"},
	} {
		if s.used && (len(s.value) < 16 || weakSecrets[strings.ToLower(s.value)]) {
			weak = append(weak, s.name)
		}
	}
	if len(weak) > 0 {
		return fmt.Errorf("refusing to start with APP_ENV=%s, because these are missing, too short or a known default: %s. "+
			"Generate each with: openssl rand -hex 32. On a machine only you can reach, APP_ENV=development skips this check",
			cfg.AppEnv, strings.Join(weak, ", "))
	}
	return nil
}
`

const routesLimitsBlock = `		// Built from the prefix the router mounts, because Sentinel matches the
		// request path exactly. Written as /api/auth/login these never matched
		// a request, so the login limit and AuthShield never fired.
		routeLimits := authRouteLimits(false)
		if isDev {
			ipLimit = &sentinel.Limit{Requests: 1000, Window: 1 * time.Minute}
			routeLimits = authRouteLimits(true)
		}
`

const routesAuthLimitsFunc = `// authRouteLimits are the per-route limits on the routes that take a password
// or a one-time code, under the live API prefix. Development gets room to test.
func authRouteLimits(dev bool) map[string]sentinel.Limit {
	v := "/api/" + APIVersion + "/auth"
	if dev {
		relaxed := sentinel.Limit{Requests: 100, Window: time.Minute}
		return map[string]sentinel.Limit{v + "/login": relaxed, v + "/register": relaxed}
	}
	return map[string]sentinel.Limit{
		v + "/login":           {Requests: 5, Window: 15 * time.Minute},
		v + "/register":        {Requests: 3, Window: 15 * time.Minute},
		v + "/forgot-password": {Requests: 3, Window: 15 * time.Minute},
		v + "/reset-password":  {Requests: 5, Window: 15 * time.Minute},
		v + "/totp/verify":     {Requests: 5, Window: 5 * time.Minute},
	}
}
`

const routesDocsBlock = `	// In production only with API_DOCS_PUBLIC=true: the reference maps every
	// route, admin, backup, SSO and GDPR ones included, with a console to call
	// them.
	if cfg.AppEnv != "production" || cfg.APIDocsPublic {
		registerAPIDocs(r, db, cfg)
	}
`

const csrfBootstrapMap = `bootstrap := map[string]bool{
		"login":                    true,
		"register":                 true,
		"refresh":                  true,
		"forgot-password":          true,
		"reset-password":           true,
		"totp/verify":              true,
		"totp/backup-codes/verify": true,
	}`

const csrfAuthRouteFunc = `// authRoute is what follows /api/auth/ or /api/<version>/auth/ in path, or ""
// for anything else. The bootstrap list was written as /api/auth/login while the
// router mounts /api/v1, so not one of its entries ever matched a request.
func authRoute(path string) string {
	rest, ok := strings.CutPrefix(path, "/api/")
	if !ok {
		return ""
	}
	if i := strings.IndexByte(rest, '/'); i > 1 && rest[0] == 'v' && strings.Trim(rest[1:i], "0123456789") == "" {
		rest = rest[i+1:]
	}
	name, ok := strings.CutPrefix(rest, "auth/")
	if !ok {
		return ""
	}
	return name
}
`

// serverDashboardLogs announces only what is mounted. Printed unconditionally,
// a production log said Studio and the docs were available when neither was.
const serverDashboardLogs = `		if cfg.GORMStudioEnabled {
			log.Printf("GORM Studio available at http://localhost:%s/studio", cfg.Port)
		}
		if cfg.AppEnv != "production" || cfg.APIDocsPublic {
			log.Printf("API Documentation at http://localhost:%s/docs", cfg.Port)
		}
`

var serverLogsRe = regexp.MustCompile(`(?m)^[ \t]*log\.Printf\("GORM Studio available at http://localhost:%s/studio", cfg\.Port\)\n[ \t]*log\.Printf\("API Documentation at http://localhost:%s/docs", cfg\.Port\)\n`)

// repairServerLogsSource makes the startup log name only what is mounted.
func repairServerLogsSource(src string) (string, []string, []string) {
	next := serverLogsRe.ReplaceAllString(src, serverDashboardLogs)
	if next == src {
		return src, nil, nil
	}
	return next, []string{"the startup log names GORM Studio and /docs only when they are mounted"}, nil
}

// repairProductionSafety applies the production defaults to a project
// scaffolded before v3.243.0. config.go, routes.go and the CSRF middleware are
// the developer's, so each edit anchors on the text Grit generated.
func repairProductionSafety(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	config := filepath.Join(apiRoot, "internal", "config", "config.go")
	if fileExists(config) {
		if err := repairSourceFile(root, m, config, repairConfigProductionSource); err != nil {
			return err
		}
	}
	// The docs switch reads a field only the config repair adds.
	docsSwitch := fileContains(config, "APIDocsPublic bool")
	steps := []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "routes", "routes.go"), func(src string) (string, []string, []string) {
			return repairRoutesProductionSource(src, docsSwitch)
		}},
		{filepath.Join(apiRoot, "internal", "middleware", "csrf.go"), repairCSRFSource},
	}
	if docsSwitch {
		steps = append(steps, struct {
			path string
			fn   func(string) (string, []string, []string)
		}{filepath.Join(apiRoot, "cmd", "server", "main.go"), repairServerLogsSource})
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

var (
	appEnvDevRe  = regexp.MustCompile(`(?m)^([ \t]*)AppEnv:([ \t]*)getEnv\("APP_ENV", "development"\),`)
	oauthFieldRe = regexp.MustCompile(`(?m)^([ \t]*OAuthFrontendURL[ \t]+string[^\n]*\n)\}`)
	loginRouteRe = regexp.MustCompile(`LoginRoute:([ \t]*)"/api/auth/login",`)
	docsCallRe   = regexp.MustCompile(`(?m)^[ \t]*registerAPIDocs\(r, db, cfg\)[ \t]*\n`)
	// The two route-limit maps as the scaffold wrote them, whitespace aside.
	oldRouteLimitsRe = regexp.MustCompile(`(?s)[ \t]*routeLimits := map\[string\]sentinel\.Limit\{\s*"/api/auth/login":\s*\{Requests: 5, Window: 15 \* time\.Minute\},\s*"/api/auth/register":\s*\{Requests: 3, Window: 15 \* time\.Minute\},\s*\}\s*if isDev \{\s*ipLimit = &sentinel\.Limit\{Requests: 1000, Window: 1 \* time\.Minute\}\s*routeLimits = map\[string\]sentinel\.Limit\{\s*"/api/auth/login":\s*\{Requests: 100, Window: 1 \* time\.Minute\},\s*"/api/auth/register":\s*\{Requests: 100, Window: 1 \* time\.Minute\},\s*\}\s*\}\n`)
	csrfBootstrapRe  = regexp.MustCompile(`bootstrap := map\[string\]bool\{[^}]*"/api/auth/login"[^}]*\}`)
)

func repairConfigProductionSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func Load() (*Config, error)") {
		return src, nil, nil
	}
	out := src
	var fixed []string
	if appEnvDevRe.MatchString(out) {
		out = appEnvDevRe.ReplaceAllString(out,
			"${1}// Production unless told otherwise: a server that forgot APP_ENV is strict.\n${1}AppEnv:${2}getEnv(\"APP_ENV\", \"production\"),")
		fixed = append(fixed, "APP_ENV defaults to production")
	}
	if !strings.Contains(out, "func checkSecrets(") {
		const anchor = "\tif cfg.JWTSecret == \"\" {"
		withField := oauthFieldRe.ReplaceAllString(out,
			"${1}\n\t// Serve the API reference at /docs in production (API_DOCS_PUBLIC).\n\tAPIDocsPublic bool\n}")
		if withField == out || !strings.Contains(withField, anchor) {
			return src, nil, []string{"is not the file Grit generated, so production still accepts default secrets and runs GORM Studio: see the v3.243.0 changelog for the checks to add"}
		}
		next := strings.Replace(withField, anchor, configProductionBlock+anchor, 1)
		next, ok := insertBeforeDecl(next, "func (c *Config) IsDevelopment() bool {", configCheckSecretsFunc)
		if !ok {
			return src, nil, []string{"has no IsDevelopment method to place checkSecrets beside, so it was left alone"}
		}
		out = next
		fixed = append(fixed, "outside development it refuses default or short secrets; in production GORM Studio is off unless GORM_STUDIO_IN_PRODUCTION=true, and read-only")
	}
	return out, fixed, nil
}

func repairRoutesProductionSource(src string, docsSwitch bool) (string, []string, []string) {
	if !strings.Contains(src, "sentinel.Config{") {
		return src, nil, nil
	}
	out := src
	var fixed, warn []string
	if !strings.Contains(out, "authRouteLimits(") && strings.Contains(out, `"/api/auth/login"`) {
		const decl = "func wafExcludedRoutes() []string {"
		next := oldRouteLimitsRe.ReplaceAllString(out, routesLimitsBlock)
		if next != out && strings.Contains(next, decl) {
			next, _ = insertBeforeDecl(next, decl, routesAuthLimitsFunc)
			out = next
			fixed = append(fixed, "the login, register, password-reset and 2FA rate limits now match the /api/v1 routes")
		} else {
			warn = append(warn, `its Sentinel route limits are not the generated ones, so they were left alone: key them by "/api/" + APIVersion + "/auth/login", or they never match a request`)
		}
	}
	if next := loginRouteRe.ReplaceAllString(out, `LoginRoute:${1}"/api/" + APIVersion + "/auth/login",`); next != out {
		out = next
		fixed = append(fixed, "AuthShield watches the real login route")
	}
	if docsSwitch && !strings.Contains(out, "cfg.APIDocsPublic") {
		if next := docsCallRe.ReplaceAllString(out, routesDocsBlock); next != out {
			out = next
			fixed = append(fixed, "/docs is served in production only with API_DOCS_PUBLIC=true")
		}
	}
	return out, fixed, warn
}

func repairCSRFSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func AutoCSRF()") || strings.Contains(src, "func authRoute(") {
		return src, nil, nil
	}
	next := csrfBootstrapRe.ReplaceAllString(src, csrfBootstrapMap)
	if next == src || !strings.Contains(next, "if bootstrap[path] {") {
		return src, nil, []string{"its CSRF bootstrap list is not the generated one, so it was left alone: /api/auth/... entries never match the /api/v1 routes"}
	}
	next = strings.Replace(next, "if bootstrap[path] {", "if bootstrap[authRoute(path)] {", 1)
	next, ok := insertBeforeDecl(next, "func AutoCSRF() gin.HandlerFunc {", csrfAuthRouteFunc)
	if !ok {
		return src, nil, nil
	}
	return next, []string{"sign-in routes are exempt from CSRF on their real /api/v1 paths"}, nil
}

// insertBeforeDecl puts text before the top-level declaration that starts with
// decl, above its doc comment, followed by a blank line.
func insertBeforeDecl(src, decl, text string) (string, bool) {
	i := strings.Index(src, "\n"+decl)
	if i < 0 {
		return src, false
	}
	start := i + 1
	for start > 0 {
		prevEnd := start - 1
		prevStart := strings.LastIndex(src[:prevEnd], "\n") + 1
		if !strings.HasPrefix(strings.TrimSpace(src[prevStart:prevEnd]), "//") {
			break
		}
		start = prevStart
	}
	return src[:start] + text + "\n" + src[start:], true
}
