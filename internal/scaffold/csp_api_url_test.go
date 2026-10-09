package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// The CSP authorises what the app is about to call, so the two have to agree
// about where the API is.
//
// They did not. The client's base URL read NEXT_PUBLIC_API_URL || API_URL,
// while the Content-Security-Policy read only the public name. The generated
// root .env sets API_URL, so a developer who followed the comment above
// APP_PORT and moved the API to another port got a policy still naming 8080
// and an app calling the new one: every request blocked, and reported as a
// console violation rather than as a status, so the network tab showed nothing
// that looked like a failure.
//
// Found by moving APP_PORT to 8099 on a generated triple-tier project and
// watching the admin's sign-in do nothing.

// envNamesIn pulls the process.env.X names out of one expression.
var envNamesIn = regexp.MustCompile(`process\.env\.([A-Z_][A-Z0-9_]*)`)

// originExpression returns the right-hand side of the API_ORIGIN declaration.
func originExpression(t *testing.T, name, src string) string {
	t.Helper()
	const marker = "const API_ORIGIN = toOrigin("
	start := strings.Index(src, marker)
	if start < 0 {
		t.Fatalf("%s declares no API_ORIGIN", name)
	}
	rest := src[start+len(marker):]
	end := strings.Index(rest, ");")
	if end < 0 {
		t.Fatalf("%s: the API_ORIGIN declaration does not end", name)
	}
	return rest[:end]
}

// clientURLExpression returns the NEXT_PUBLIC_API_URL value the env block
// bakes into the browser bundle.
func clientURLExpression(t *testing.T, name, src string) string {
	t.Helper()
	const marker = "NEXT_PUBLIC_API_URL:"
	start := strings.Index(src, marker)
	if start < 0 {
		t.Fatalf("%s sets no NEXT_PUBLIC_API_URL in its env block", name)
	}
	rest := src[start+len(marker):]
	end := strings.Index(rest, ",\n")
	if end < 0 {
		t.Fatalf("%s: the NEXT_PUBLIC_API_URL entry does not end", name)
	}
	return rest[:end]
}

func nextConfigs() map[string]string {
	opts := Options{ProjectName: "app", Frontend: FrontendNext}
	return map[string]string{
		"apps/admin/next.config.ts": adminNextConfig(opts),
		"apps/web/next.config.ts":   webNextConfig(opts),
	}
}

// The policy and the client read the same environment variables, in the same
// order.
func TestTheCSPAndTheClientAgreeOnTheAPIOrigin(t *testing.T) {
	for name, src := range nextConfigs() {
		origin := originExpression(t, name, src)
		client := clientURLExpression(t, name, src)

		var originNames, clientNames []string
		for _, m := range envNamesIn.FindAllStringSubmatch(origin, -1) {
			originNames = append(originNames, m[1])
		}
		for _, m := range envNamesIn.FindAllStringSubmatch(client, -1) {
			clientNames = append(clientNames, m[1])
		}

		if strings.Join(originNames, ",") != strings.Join(clientNames, ",") {
			t.Errorf("%s: the CSP origin reads %v and the client reads %v. "+
				"Whichever one the project sets, the policy has to authorise the host "+
				"the app will actually call, or every request is blocked with a console "+
				"violation and no HTTP status.",
				name, originNames, clientNames)
		}
	}
}

// API_URL specifically, because that is the one the generated .env sets.
func TestTheCSPOriginReadsAPIURL(t *testing.T) {
	for name, src := range nextConfigs() {
		origin := originExpression(t, name, src)
		if !strings.Contains(origin, "process.env.API_URL") {
			t.Errorf("%s: the CSP origin does not read API_URL, which is the name the "+
				"generated root .env sets. Moving APP_PORT then blocks every request "+
				"from this app.\n  got: %s", name, strings.TrimSpace(origin))
		}
	}
}

// The root .env names the port in more than one place, and they have to match
// out of the box or a new project is broken before anybody touches it.
func TestTheGeneratedEnvAgreesWithItselfAboutThePort(t *testing.T) {
	env := envFile(Options{ProjectName: "app", Frontend: FrontendNext})

	port := envValue(t, env, "APP_PORT")
	for _, key := range []string{"APP_URL", "API_URL"} {
		value := envValue(t, env, key)
		if !strings.Contains(value, ":"+port) {
			t.Errorf(".env sets APP_PORT=%s but %s=%s. They name one address in "+
				"several places and have to move together.", port, key, value)
		}
	}
}

// envValue reads one KEY=value out of a rendered .env.
func envValue(t *testing.T, env, key string) string {
	t.Helper()
	for _, line := range strings.Split(strings.ReplaceAll(env, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, key+"=") {
			return strings.TrimSpace(strings.TrimPrefix(line, key+"="))
		}
	}
	t.Fatalf("the generated .env sets no %s", key)
	return ""
}

// An existing project gets the fallback from an upgrade, and one that already
// has it is left alone.
func TestTheAPIURLRepairUpgradesAnOldConfig(t *testing.T) {
	old := "x\n" + nextAPIOriginLineOld + "y\n"

	out, changed, warnings := repairCSPAPIURLFallbackSource(old)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(changed) == 0 {
		t.Fatal("the repair reported no change on a config that needed one")
	}
	if !strings.Contains(out, "process.env.API_URL") {
		t.Fatalf("the repair did not add the API_URL fallback:\n%s", out)
	}

	// Run again: nothing left to do.
	again, changed, _ := repairCSPAPIURLFallbackSource(out)
	if len(changed) != 0 {
		t.Errorf("the repair changed an already-repaired config: %v", changed)
	}
	if again != out {
		t.Error("the repair is not idempotent")
	}
}

// The websocket repair still recognises both spellings of the origin line,
// because a project may be on either one.
func TestTheWebSocketRepairAcceptsBothOriginLines(t *testing.T) {
	for name, origin := range map[string]string{
		"old": nextAPIOriginLineOld,
		"new": nextAPIOriginLine,
	} {
		src := origin + "  " + nextConnectSrcOld + " ? \"\" : \"\")\n"
		out, changed, warnings := repairCSPWebSocketNextSource(src)
		if len(warnings) != 0 {
			t.Errorf("%s origin line: the repair refused it: %v", name, warnings)
			continue
		}
		if len(changed) == 0 {
			t.Errorf("%s origin line: the repair did nothing", name)
			continue
		}
		if !strings.Contains(out, "API_WS_ORIGIN") {
			t.Errorf("%s origin line: no socket origin was added", name)
		}
	}
}
