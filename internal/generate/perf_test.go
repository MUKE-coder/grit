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

// The generated load test registers itself with Pulse. Without this the Test
// Runs page stays empty in every project, which is how the feature sat unused:
// Pulse has had POST /pulse/api/test-runs since v1.0 and nothing called it.
func TestPerfScriptRegistersTheRunWithPulse(t *testing.T) {
	s := BuildPerfScript(PerfOptions{VUs: 50, Duration: "2m"})
	for _, want := range []string{
		"import { pulseBase, pulseLogin, pulseStartRun, pulseEndRun } from './pulse.js'",
		"const pulseToken = pulseLogin(pulse)",
		"pulseStartRun(pulse, pulseToken, {",
		"type: 'k6.average-load'",
		"export function teardown(data)",
		"pulseEndRun(data.pulse, data.pulseToken, data.run",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("the script does not record the run in Pulse: missing %q", want)
		}
	}
	// The run carries the shape of the load, so two runs can be told apart.
	for _, want := range []string{"50 VUs for 2m", "vus_peak: 50", "hold: '2m'"} {
		if !strings.Contains(s, want) {
			t.Errorf("the run is not labelled with the load it applied: missing %q", want)
		}
	}
	// setup() must hand teardown() what it needs; k6 passes only the return value.
	if !strings.Contains(s, "return { token, staff, pulse, pulseToken, run };") {
		t.Error("setup does not return the Pulse handles, so teardown cannot close the run")
	}
}

// An observability hiccup must never fail a load test, so every call gives up
// quietly rather than throwing. Pulse's own reference bridge throws when the
// URL is missing; a generated project runs in places where Pulse is off.
func TestPerfPulseBridgeFailsSoft(t *testing.T) {
	b := BuildPerfPulseBridge()
	for _, want := range []string{
		"if (!password) {",
		"return null;",
		"if (!token) return null;",
		"if (!token || !run || !run.id) return;",
		"console.warn",
	} {
		if !strings.Contains(b, want) {
			t.Errorf("the bridge is not fail-soft: missing %q", want)
		}
	}
	if strings.Contains(b, "throw new Error") {
		t.Error("the bridge throws, so an unreachable dashboard would fail the load test")
	}
	// Pulse is on the API's own origin in a Grit project, so there is no second
	// URL to configure and PULSE_URL is only an override.
	if !strings.Contains(b, "__ENV.PULSE_URL || apiBase") {
		t.Error("the bridge does not default the Pulse origin to the API's own")
	}
	// The links are the payoff: a recorded run nobody opens is a row in a table.
	for _, want := range []string{"/pulse/ui/test-runs", "/compare", "/profile?format=svg"} {
		if !strings.Contains(b, want) {
			t.Errorf("the bridge does not print where to read the result: missing %q", want)
		}
	}
	// But the profile link only when there is a profile. Printed
	// unconditionally it is a 404 whenever profiling is off, and a tool that
	// prints dead links stops being read.
	if !strings.Contains(b, "if (run.profiling) {") {
		t.Error("the bridge prints the flame graph link even when no profile was taken")
	}
}

func TestPerfReadmeDocumentsPulse(t *testing.T) {
	r := BuildPerfReadme(PerfOptions{})
	for _, want := range []string{
		"It records itself in Pulse",
		"set -a; source .env; set +a",
		"PULSE_PROFILING=true",
		"profile?format=svg",
	} {
		if !strings.Contains(r, want) {
			t.Errorf("the runbook does not explain the Pulse half: missing %q", want)
		}
	}
}

// The Pulse calls must not share the VU cookie jar. setup() logs into the API
// first, which sets a session cookie; k6 keeps one jar per VU, so a shared jar
// puts that cookie on the Pulse POST and the API's CSRF middleware answers 403
// CSRF_INVALID. The run then goes unrecorded and nothing says so loudly.
// Found by running the generated test against a real project.
func TestPerfPulseBridgeUsesItsOwnCookieJar(t *testing.T) {
	b := BuildPerfPulseBridge()
	if !strings.Contains(b, "const jar = new http.CookieJar();") {
		t.Error("the bridge shares the VU cookie jar, so the API session cookie will 403 the Pulse POST")
	}
	if !strings.Contains(b, "return { headers, jar, tags: { pulse: 'true' } };") {
		t.Error("pulseParams does not attach the dedicated jar")
	}
	// Every call goes through pulseParams rather than building headers inline,
	// so a fourth call cannot quietly miss the jar.
	if n := strings.Count(b, "pulseParams("); n < 4 {
		t.Errorf("only %d uses of pulseParams; each Pulse request must go through it", n)
	}
	if strings.Contains(b, "headers: { 'Content-Type': 'application/json', Authorization:") {
		t.Error("a Pulse request still builds its params inline, bypassing the jar")
	}
}
