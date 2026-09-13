package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// configV241 is the shape of internal/config/config.go before v3.243.0, trimmed
// to the parts the repair touches.
const configV241 = `package config

import (
	"fmt"
	"strings"
)

type Config struct {
	AppEnv               string
	JWTSecret            string
	GORMStudioEnabled    bool
	GORMStudioPassword   string
	GORMStudioReadOnly   bool
	GORMStudioDisableSQL bool
	SentinelEnabled      bool
	SentinelPassword     string
	SentinelSecretKey    string
	PulseEnabled         bool
	PulsePassword        string
	OAuthFrontendURL   string // Where to redirect after OAuth callback
}

func getEnv(key, fallback string) string { return fallback }

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		OAuthFrontendURL:   strings.TrimSpace(getEnv("OAUTH_FRONTEND_URL", "")),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	return cfg, nil
}

// IsDevelopment returns true if the app is running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}
`

const routesV241 = `package routes

func wafNothing() {}

// wafExcludedRoutes lists the paths Sentinel's WAF steps aside for.
func wafExcludedRoutes() []string {
	return nil
}

func Setup() {
	if cfg.SentinelEnabled {
		isDev := cfg.AppEnv == "development"
		ipLimit := &sentinel.Limit{Requests: 100, Window: 1 * time.Minute}
		routeLimits := map[string]sentinel.Limit{
			"/api/auth/login":    {Requests: 5, Window: 15 * time.Minute},
			"/api/auth/register": {Requests: 3, Window: 15 * time.Minute},
		}
		if isDev {
			ipLimit = &sentinel.Limit{Requests: 1000, Window: 1 * time.Minute}
			routeLimits = map[string]sentinel.Limit{
				"/api/auth/login":    {Requests: 100, Window: 1 * time.Minute},
				"/api/auth/register": {Requests: 100, Window: 1 * time.Minute},
			}
		}
		_ = sentinel.Config{
			RateLimit:  sentinel.RateLimitConfig{ByIP: ipLimit, ByRoute: routeLimits},
			AuthShield: sentinel.AuthShieldConfig{
				Enabled:    !isDev,
				LoginRoute: "/api/auth/login",
			},
		}
	}

	// API Documentation
	registerAPIDocs(r, db, cfg)
}
`

const csrfV241 = `package middleware

func isSAML(path string) bool {
	return strings.HasPrefix(path, "/api/auth/saml/")
}

func AutoCSRF() gin.HandlerFunc {
	bootstrap := map[string]bool{
		"/api/auth/login":                    true,
		"/api/auth/register":                 true,
		"/api/auth/totp/backup-codes/verify": true,
	}
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if bootstrap[path] {
			c.Next()
			return
		}
	}
}
`

func assertRepaired(t *testing.T, what, out string, fixed []string, want ...string) {
	t.Helper()
	if len(fixed) == 0 {
		t.Errorf("%s: nothing reported fixed", what)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("%s is not valid Go after the repair: %v\n%s", what, err, out)
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("%s is missing %s", what, w)
		}
	}
}

func TestRepairConfigMakesProductionStrict(t *testing.T) {
	out, fixed, warn := repairConfigProductionSource(configV241)
	if len(warn) > 0 {
		t.Fatalf("warnings: %v", warn)
	}
	assertRepaired(t, "config.go", out, fixed,
		`getEnv("APP_ENV", "production")`,
		"APIDocsPublic bool",
		`cfg.GORMStudioEnabled = getEnv("GORM_STUDIO_IN_PRODUCTION", "false") == "true"`,
		"func checkSecrets(cfg *Config) error {",
		`if cfg.AppEnv != "development" {`,
	)
	if strings.Index(out, "func checkSecrets") > strings.Index(out, "// IsDevelopment returns") {
		t.Error("checkSecrets was not placed above IsDevelopment's doc comment")
	}
	if again, fixed, _ := repairConfigProductionSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed config.go again")
	}
}

func TestRepairRoutesKeysLimitsByTheLivePrefix(t *testing.T) {
	out, fixed, warn := repairRoutesProductionSource(routesV241, true)
	if len(warn) > 0 {
		t.Fatalf("warnings: %v", warn)
	}
	assertRepaired(t, "routes.go", out, fixed,
		"routeLimits := authRouteLimits(false)",
		"func authRouteLimits(dev bool) map[string]sentinel.Limit {",
		`LoginRoute: "/api/" + APIVersion + "/auth/login",`,
		`if cfg.AppEnv != "production" || cfg.APIDocsPublic {`,
	)
	if strings.Contains(out, `"/api/auth/`) {
		t.Errorf("an unversioned auth path survived:\n%s", out)
	}
	if again, fixed, _ := repairRoutesProductionSource(out, true); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed routes.go again")
	}
}

// Without the config field the docs switch would not compile, so it waits.
func TestRepairRoutesLeavesDocsWithoutTheConfigField(t *testing.T) {
	out, _, _ := repairRoutesProductionSource(routesV241, false)
	if strings.Contains(out, "cfg.APIDocsPublic") {
		t.Error("the docs switch was added without the config field it reads")
	}
}

func TestRepairCSRFMatchesVersionedPaths(t *testing.T) {
	out, fixed, warn := repairCSRFSource(csrfV241)
	if len(warn) > 0 {
		t.Fatalf("warnings: %v", warn)
	}
	assertRepaired(t, "csrf.go", out, fixed,
		"if bootstrap[authRoute(path)] {",
		`"totp/backup-codes/verify": true,`,
		"func authRoute(path string) string {",
	)
}

// A fresh project and an upgraded one must end up with the same code, so the
// templates must need no repair at all.
func TestFreshTemplatesNeedNoProductionRepair(t *testing.T) {
	for name, c := range map[string]struct {
		src string
		fn  func(string) (string, []string, []string)
	}{
		"config.go": {apiConfigGo(), repairConfigProductionSource},
		"routes.go": {apiRoutesGo(), func(s string) (string, []string, []string) { return repairRoutesProductionSource(s, true) }},
		"csrf.go":   {csrfMiddlewareGo(), repairCSRFSource},
	} {
		out, fixed, warn := c.fn(c.src)
		if out != c.src || len(fixed) > 0 || len(warn) > 0 {
			t.Errorf("%s: the scaffold template still needs the repair (fixed %v, warn %v)", name, fixed, warn)
		}
	}
}

// A production log said Studio and the docs were available when neither was.
func TestRepairServerLogsNameOnlyWhatIsMounted(t *testing.T) {
	src := "package main\n\nfunc main() {\n\tgo func() {\n" +
		"\t\tlog.Printf(\"GORM Studio available at http://localhost:%s/studio\", cfg.Port)\n" +
		"\t\tlog.Printf(\"API Documentation at http://localhost:%s/docs\", cfg.Port)\n\t}()\n}\n"
	out, fixed, _ := repairServerLogsSource(src)
	assertRepaired(t, "main.go", out, fixed, "if cfg.GORMStudioEnabled {", `if cfg.AppEnv != "production" || cfg.APIDocsPublic {`)
	if again, fixed, _ := repairServerLogsSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed main.go again")
	}
	if out, fixed, _ := repairServerLogsSource(apiMainGo(Options{ProjectName: "x"})); len(fixed) > 0 || !strings.Contains(out, "if cfg.GORMStudioEnabled {") {
		t.Error("the scaffold's main.go still logs Studio and the docs unconditionally")
	}
}

func TestInsertBeforeDeclSkipsTheDocComment(t *testing.T) {
	src := "package x\n\nvar a = 1\n\n// F does a thing.\n// Really.\nfunc F() {}\n"
	out, ok := insertBeforeDecl(src, "func F() {", "func G() {}\n")
	if !ok || !strings.Contains(out, "var a = 1\n\nfunc G() {}\n\n// F does a thing.") {
		t.Errorf("got %q", out)
	}
}
