package generate

import "strings"

import "testing"

func TestPerfScriptDefaults(t *testing.T) {
	s := BuildPerfScript(PerfOptions{})
	// Defaults must materialize in the rendered script.
	for _, want := range []string{
		"target: 20",      // default VUs
		"duration: '30s'", // default hold duration
		"http://localhost:8080",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("default script missing %q", want)
		}
	}
}

func TestPerfScriptHasThresholds(t *testing.T) {
	s := BuildPerfScript(PerfOptions{})
	for _, want := range []string{
		"http_req_failed: ['rate<0.01']",
		"http_req_duration: ['p(95)<500']",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("script missing threshold %q — it wouldn't gate a regression", want)
		}
	}
}

func TestPerfScriptAuthSetup(t *testing.T) {
	s := BuildPerfScript(PerfOptions{})
	// setup() must register, log in, and pull the token from the real response
	// shape (data.tokens.access_token). A wrong path would leave every VU
	// unauthenticated and the load test meaningless.
	for _, want := range []string{
		"export function setup()",
		"/api/auth/register",
		"/api/auth/login",
		"res.json('data.tokens.access_token')",
		"Authorization: `Bearer ${data.token}`",
		"/api/health",
		"/api/profile",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("script missing auth/scenario element %q", want)
		}
	}
}

func TestPerfScriptResourceScenario(t *testing.T) {
	// Without a resource, no collection scenario.
	base := BuildPerfScript(PerfOptions{})
	if strings.Contains(base, "/api/blogs") {
		t.Errorf("no resource requested but a resource endpoint leaked in")
	}
	// With a resource, the pluralized, lowercased collection endpoint appears.
	s := BuildPerfScript(PerfOptions{Resource: "Blog"})
	if !strings.Contains(s, "${BASE}/api/blogs`") {
		t.Errorf("resource scenario missing GET /api/blogs")
	}
	if !strings.Contains(s, "'blog list is 200'") {
		t.Errorf("resource check missing")
	}
}

func TestPerfScriptCustomLoad(t *testing.T) {
	s := BuildPerfScript(PerfOptions{VUs: 100, Duration: "2m", BaseURL: "https://api.example.com"})
	for _, want := range []string{"target: 100", "duration: '2m'", "https://api.example.com"} {
		if !strings.Contains(s, want) {
			t.Errorf("custom option not applied: %q", want)
		}
	}
}

func TestPerfReadmeMatchesOptions(t *testing.T) {
	r := perfReadme(PerfOptions{Resource: "Order", VUs: 50, Duration: "1m"})
	for _, want := range []string{"k6 run perf/load.js", "/api/orders", "50 virtual users", "1m", "p95 under 500ms"} {
		if !strings.Contains(r, want) {
			t.Errorf("README missing %q", want)
		}
	}
}
