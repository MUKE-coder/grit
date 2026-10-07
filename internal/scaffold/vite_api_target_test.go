package scaffold

import (
	"strings"
	"testing"
)

// A Vite app has to find the project's .env before it can read anything from it.
//
// The single's config was fixed for this; the monorepo's two were not. They
// called loadEnv with process.cwd(), which in a monorepo is apps/web or
// apps/admin, where there is no .env. Every value came back undefined: the dev
// proxy stayed on 8080 and the CSP's connect-src with it, so moving APP_PORT
// broke the login and the fetches at once, both silently.
func TestViteConfigsFindTheProjectEnv(t *testing.T) {
	opts := Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendTanStack}
	for name, cfg := range map[string]string{
		"the web app's": webTanStackViteConfig(opts),
		"the admin's":   adminTanStackViteConfig(),
		"the single's":  singleFrontendViteConfig(),
	} {
		if !strings.Contains(cfg, "loadEnv(") {
			t.Errorf("%s vite config does not read .env", name)
			continue
		}
		if !strings.Contains(cfg, "apiTarget") {
			t.Errorf("%s vite config has no apiTarget", name)
		}
		if strings.Contains(cfg, "target: 'http://localhost:8080'") {
			t.Errorf("%s vite config still proxies to a written-in 8080", name)
		}
		if !strings.Contains(cfg, "APP_PORT") {
			t.Errorf("%s vite config ignores APP_PORT, the variable the Go binary reads", name)
		}
	}
	// A monorepo app runs two levels below the .env it needs.
	for name, cfg := range map[string]string{
		"the web app's": webTanStackViteConfig(opts),
		"the admin's":   adminTanStackViteConfig(),
	} {
		if !strings.Contains(cfg, "loadEnv(viteMode, path.resolve(__dirname, '../..'), '')") {
			t.Errorf("%s vite config looks for .env in its own directory", name)
		}
	}
	// The single's frontend IS the project root.
	if strings.Contains(singleFrontendViteConfig(), "'../..'") {
		t.Error("the single's vite config looks outside the project for .env")
	}
}
