package scaffold

import (
	"strings"
	"testing"
)

const oldSentinelProxiesLoad = `SentinelTrustedProxies: splitCSV(getEnv("SENTINEL_TRUSTED_PROXIES", "")),`

func TestTemplatesTrustOnlyConfiguredProxies(t *testing.T) {
	routes := apiRoutesGo()
	if !strings.Contains(routes, "middleware.TrustProxies(r, cfg.TrustedProxies)") {
		t.Error("routes.Setup does not restrict forwarded headers to TRUSTED_PROXIES")
	}
	// Registered before anything that reads the client address.
	if strings.Index(routes, "middleware.TrustProxies(") > strings.Index(routes, "mountGlobalMiddleware(r, cfg, svc)") {
		t.Error("TrustProxies is registered after the global middleware")
	}
	config := apiConfigGo()
	for _, want := range []string{"TrustedProxies []string", "func trustedProxies()", `os.Getenv("TRUSTED_PROXIES")`, `strings.Join(trustedProxies(), ",")`} {
		if !strings.Contains(config, want) {
			t.Errorf("config.go lacks %q", want)
		}
	}
	if !strings.Contains(middlewareProxiesGo(), "r.SetTrustedProxies(") || !strings.Contains(middlewareProxiesGo(), `"X-Forwarded-Proto"`) {
		t.Error("middleware/proxies.go does not set gin's trusted proxies and strip X-Forwarded-Proto")
	}
	if !strings.Contains(envFile(Options{ProjectName: "app"}), "# TRUSTED_PROXIES=") {
		t.Error(".env does not document TRUSTED_PROXIES")
	}
}

func TestTrustedProxiesRepairs(t *testing.T) {
	config := apiConfigGo()
	oldConfig := strings.Replace(config, configTrustedProxiesField, "", 1)
	oldConfig = strings.Replace(oldConfig, configSentinelProxiesNew, oldSentinelProxiesLoad, 1)
	oldConfig = strings.Replace(oldConfig, configTrustedProxiesFunc, "", 1)
	checkRepair(t, "config.go", repairTrustedProxiesConfigSource, oldConfig, config)

	routes := apiRoutesGo()
	checkRepair(t, "routes.go", repairTrustedProxiesRoutesSource,
		strings.Replace(routes, routesTrustProxiesBlock, "\tr := gin.New()\n\n\t// Global middleware\n", 1), routes)

	// gofmt aligns the struct literal in a generated project.
	aligned := strings.Replace(oldConfig, oldSentinelProxiesLoad, `SentinelTrustedProxies:   splitCSV(getEnv("SENTINEL_TRUSTED_PROXIES", "")),`, 1)
	if out, fixed, _ := repairTrustedProxiesConfigSource(aligned); len(fixed) == 0 || !strings.Contains(out, "TrustedProxies:         trustedProxies(),") {
		t.Error("the config repair does not match gofmt's alignment")
	}

	env := envFile(Options{ProjectName: "app"})
	checkRepair(t, ".env", repairTrustedProxiesEnvSource, strings.Replace(env, envTrustedProxies, "", 1), env)

	// A config.go that is not Grit's is named, not rewritten.
	odd := strings.Replace(oldConfig, oldSentinelProxiesLoad, `SentinelTrustedProxies: myProxies(),`, 1)
	if out, _, warnings := repairTrustedProxiesConfigSource(odd); out != odd || len(warnings) != 1 {
		t.Errorf("an unrecognised config.go was changed or not reported: %v", warnings)
	}
}
