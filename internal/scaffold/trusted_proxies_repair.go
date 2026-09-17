package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review L6: the config and routes text, spliced into the templates
// and put into an existing project by repairTrustedProxies.
const (
	configTrustedProxiesField = "\t// TrustedProxies are the reverse proxies (IPs or CIDRs) whose X-Forwarded-For,\n" +
		"\t// X-Real-IP and X-Forwarded-Proto the API believes. Read from TRUSTED_PROXIES;\n" +
		"\t// see trustedProxies for the default.\n" +
		"\tTrustedProxies []string\n"

	// Sentinel follows TRUSTED_PROXIES unless SENTINEL_TRUSTED_PROXIES says
	// otherwise, so its rate limits and the API's client IP name one address.
	configSentinelProxiesNew = `SentinelTrustedProxies: splitCSV(getEnv("SENTINEL_TRUSTED_PROXIES", strings.Join(trustedProxies(), ","))),` + "\n" +
		"\t\tTrustedProxies:         trustedProxies(),"

	configTrustedProxiesAnchor = "// splitCSV trims and splits a comma-separated env var."
	configTrustedProxiesFunc   = `// defaultTrustedProxies are the addresses a reverse proxy on this machine or on
// a private network (a Docker network, a VPC) connects from. A client on the
// internet cannot connect from one, so the proxies of the documented deployments
// are believed and nobody else is.
var defaultTrustedProxies = []string{"127.0.0.0/8", "::1/128", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7"}

// trustedProxies reads TRUSTED_PROXIES: a comma-separated list of IPs and CIDRs,
// "none" to trust no proxy, or unset for defaultTrustedProxies. Behind a CDN
// that connects to the API directly (Cloudflare with no proxy of your own), list
// the CDN's ranges, or every request carries the CDN's address.
func trustedProxies() []string {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	if raw == "" {
		return append([]string(nil), defaultTrustedProxies...)
	}
	if strings.EqualFold(raw, "none") {
		return []string{}
	}
	return splitCSV(raw)
}

`

	routesTrustProxiesBlock = "\tr := gin.New()\n\n" +
		"\t// X-Forwarded-For, X-Real-IP and X-Forwarded-Proto are believed only from the\n" +
		"\t// proxies TRUSTED_PROXIES names. First, so every middleware after it keys rate\n" +
		"\t// limits, sessions and audit rows on the real client address.\n" +
		"\tif err := middleware.TrustProxies(r, cfg.TrustedProxies); err != nil {\n" +
		"\t\tlog.Printf(\"%v\", err)\n" +
		"\t}\n\n" +
		"\t// Global middleware\n"

	envTrustedProxies = "# The reverse proxies whose X-Forwarded-For and X-Forwarded-Proto the API\n" +
		"# believes, as IPs or CIDRs. Unset trusts loopback and the private ranges, which\n" +
		"# covers a proxy on this machine or on the Docker network. \"none\" trusts no\n" +
		"# proxy. Behind a CDN that connects to the API directly, list its ranges.\n" +
		"# TRUSTED_PROXIES=127.0.0.0/8,::1/128,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7\n"
)

var (
	configSentinelProxiesFieldRe = regexp.MustCompile(`(?m)^\tSentinelTrustedProxies\s+\[\]string\n`)
	configSentinelProxiesLoadRe  = regexp.MustCompile(`SentinelTrustedProxies:\s+splitCSV\(getEnv\("SENTINEL_TRUSTED_PROXIES", ""\)\),`)
	routesGinNewRe               = regexp.MustCompile(`\tr := gin\.New\(\)\n\n\t// Global middleware\n`)
	envSentinelAuditKeyRe        = regexp.MustCompile(`(?m)^SENTINEL_AUDIT_KEY=[^\n]*\n`)
	envTrustedProxiesRe          = regexp.MustCompile(`(?m)^(# )?TRUSTED_PROXIES=`)
)

// repairTrustedProxies brings L6 to an existing project: middleware/proxies.go,
// TRUSTED_PROXIES in config.go, the call at the top of routes.Setup, and the
// setting documented in .env.example.
func repairTrustedProxies(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	config := filepath.Join(apiRoot, "internal", "config", "config.go")
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(config) || !fileExists(routes) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if err := repairSourceFile(root, m, config, repairTrustedProxiesConfigSource); err != nil {
		return err
	}
	// The middleware and its call arrive only once config has the field they
	// read, so a config.go too changed to patch leaves a project that builds.
	if !fileContains(config, "func trustedProxies()") {
		return nil
	}
	dir := filepath.Join(apiRoot, "internal", "middleware")
	for path, content := range map[string]string{
		filepath.Join(dir, "proxies.go"):      middlewareProxiesGo(),
		filepath.Join(dir, "proxies_test.go"): middlewareProxiesTestGo(),
	} {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	if err := repairSourceFile(root, m, routes, repairTrustedProxiesRoutesSource); err != nil {
		return err
	}
	if env := filepath.Join(root, ".env.example"); fileExists(env) {
		if err := repairTextFile(root, m, env, repairTrustedProxiesEnvSource); err != nil {
			return err
		}
	}
	return nil
}

func repairTrustedProxiesConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func trustedProxies()") {
		return src, nil, nil
	}
	if len(configSentinelProxiesFieldRe.FindAllStringIndex(src, -1)) != 1 ||
		len(configSentinelProxiesLoadRe.FindAllStringIndex(src, -1)) != 1 ||
		strings.Count(src, configTrustedProxiesAnchor) != 1 {
		return src, nil, []string{"config.go is not the file Grit wrote: add TrustedProxies read from TRUSTED_PROXIES (default loopback and the private ranges) and pass it to middleware.TrustProxies, or the API believes X-Forwarded-For from any client"}
	}
	out := configSentinelProxiesFieldRe.ReplaceAllStringFunc(src, func(line string) string {
		return line + configTrustedProxiesField
	})
	out = configSentinelProxiesLoadRe.ReplaceAllLiteralString(out, configSentinelProxiesNew)
	out = strings.Replace(out, configTrustedProxiesAnchor, configTrustedProxiesFunc+configTrustedProxiesAnchor, 1)
	return out, []string{"TRUSTED_PROXIES names the proxies whose forwarded headers are believed"}, nil
}

func repairTrustedProxiesRoutesSource(src string) (string, []string, []string) {
	if strings.Contains(src, "middleware.TrustProxies(") {
		return src, nil, nil
	}
	if len(routesGinNewRe.FindAllStringIndex(src, -1)) != 1 {
		return src, nil, []string{"routes.go is not the file Grit wrote: call middleware.TrustProxies(r, cfg.TrustedProxies) right after gin.New(), or any client can set the IP its sessions, audit rows and rate limits are keyed on"}
	}
	out := routesGinNewRe.ReplaceAllLiteralString(src, routesTrustProxiesBlock)
	return out, []string{"X-Forwarded-For and X-Forwarded-Proto are believed only from TRUSTED_PROXIES"}, nil
}

func repairTrustedProxiesEnvSource(src string) (string, []string, []string) {
	if envTrustedProxiesRe.MatchString(src) {
		return src, nil, nil
	}
	fixed := []string{"TRUSTED_PROXIES is documented"}
	loc := envSentinelAuditKeyRe.FindStringIndex(src)
	if loc == nil {
		return strings.TrimRight(src, "\n") + "\n\n" + envTrustedProxies, fixed, nil
	}
	return src[:loc[1]] + envTrustedProxies + src[loc[1]:], fixed, nil
}
