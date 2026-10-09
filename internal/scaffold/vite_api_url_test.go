package scaffold

import (
	"strings"
	"testing"
)

// A Vite frontend could not reach its own API on any project that had moved
// APP_PORT, and it is the mirror image of the fault the Next frontends had.
//
// There the policy was stale and the client was right. Here the policy is
// right and the client is stale: vite.config.ts resolves the API's address
// from VITE_API_URL or APP_PORT and authorises exactly that origin, while the
// browser bundle reads import.meta.env.VITE_API_URL, which nothing sets, and
// falls back to localhost:8080. Every request blocked, reported as a console
// violation with no HTTP status.
//
// Found by signing in to a generated TanStack admin.

func viteConfigs() map[string]string {
	opts := Options{ProjectName: "app", Frontend: FrontendTanStack}
	return map[string]string{
		"apps/admin/vite.config.ts": adminTanStackViteConfig(),
		"apps/web/vite.config.ts":   webTanStackViteConfig(opts),
	}
}

// The config hands the browser bundle the address it just resolved.
func TestAViteConfigDefinesTheResolvedAPIURL(t *testing.T) {
	for name, src := range viteConfigs() {
		if !strings.Contains(src, "define: clientEnv") {
			t.Errorf("%s does not define clientEnv, so the browser bundle never sees the "+
				"address the config resolved and the policy authorised: every request goes "+
				"to the 8080 fallback and is blocked", name)
		}
		if !strings.Contains(src, "const apiTarget") {
			t.Errorf("%s does not resolve apiTarget", name)
		}
	}
}

// And the value it defines is the one the policy is built from, not a second
// reading of the environment.
func TestTheDefineAndThePolicyComeFromOneValue(t *testing.T) {
	src := adminTanStackViteConfig()

	if !strings.Contains(src, "const API_ORIGIN = toOrigin(apiTarget)") {
		t.Error("the policy's origin is not derived from apiTarget")
	}
	if !strings.Contains(src, "'import.meta.env.VITE_API_URL': JSON.stringify(apiTarget)") {
		t.Error("the client's VITE_API_URL is not derived from apiTarget, so the two can " +
			"disagree about where the API is")
	}
}

// Vite does not polyfill process, so a client module reading process.env
// throws the moment anything imports it. vite build uses esbuild and does no
// type checking, so nothing catches it before a user sees a blank screen.
func TestAViteRealtimeClientDoesNotReadProcessEnv(t *testing.T) {
	src := realtimeClientTS(frontendVite)

	if strings.Contains(src, "process.env") {
		for _, line := range strings.Split(src, "\n") {
			if strings.Contains(line, "process.env") {
				t.Errorf("the Vite realtime client reads process.env, which is not defined "+
					"in the browser:\n  %s", strings.TrimSpace(line))
			}
		}
	}
	if !strings.Contains(src, "import.meta.env.VITE_API_URL") {
		t.Error("the Vite realtime client does not read VITE_API_URL")
	}
}

// The other two targets keep the spelling their bundler understands.
func TestTheOtherRealtimeClientsKeepTheirOwnSpelling(t *testing.T) {
	if next := realtimeClientTS(frontendNext); !strings.Contains(next, "process.env.NEXT_PUBLIC_API_URL") {
		t.Error("the Next realtime client no longer reads NEXT_PUBLIC_API_URL")
	}
	native := realtimeClientTS(frontendNative)
	if !strings.Contains(native, "process.env.EXPO_PUBLIC_API_URL") {
		t.Error("the Expo realtime client no longer reads EXPO_PUBLIC_API_URL")
	}
	if strings.Contains(native, "import.meta.env") {
		t.Error("the Expo client reads import.meta.env, which Metro does not provide")
	}
}

// frontendFor is the one place that asks which bundler a project uses, so a
// new client cannot answer it differently.
func TestFrontendForFollowsTheFrontendOption(t *testing.T) {
	if got := frontendFor(Options{Frontend: FrontendTanStack}); got != frontendVite {
		t.Errorf("a TanStack project resolved to %v, not frontendVite", got)
	}
	if got := frontendFor(Options{Frontend: FrontendNext}); got != frontendNext {
		t.Errorf("a Next project resolved to %v, not frontendNext", got)
	}
}
