package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// The docs app's next.config.mjs is JavaScript, and has to parse as JavaScript.
//
// The shared security-header block is written for a TypeScript config, which
// apps/web and apps/admin both have. The docs app's is .mjs, because that is
// what fumadocs expects, and one type annotation in a .mjs file is a syntax
// error that stops Next loading the config at all. The documentation site in a
// --full project had never built: "Unexpected token ':'", from a line nobody had
// written by hand.
func TestDocsNextConfigIsValidJavaScript(t *testing.T) {
	cfg := docsNextConfig()

	// A parameter or return annotation: "(value: string)" or "): string {".
	annotation := regexp.MustCompile(`\(\s*\w+\s*:\s*\w+|\)\s*:\s*\w+\s*\{`)
	if m := annotation.FindString(cfg); m != "" {
		t.Errorf("next.config.mjs carries a TypeScript annotation (%q), which Next cannot parse:\n%s", m, cfg)
	}
	for _, ts := range []string{"interface ", "satisfies ", " as string"} {
		if strings.Contains(cfg, ts) {
			t.Errorf("next.config.mjs contains TypeScript syntax %q", ts)
		}
	}

	// And it is still the same policy: a JS flavour that dropped the headers
	// would parse perfectly and protect nothing.
	for _, want := range []string{"Content-Security-Policy", "toOrigin", "Strict-Transport-Security"} {
		if !strings.Contains(cfg, want) {
			t.Errorf("the docs config no longer sets %s", want)
		}
	}
}

// Every construct the transform claims to handle, and the shape of the output.
//
// The exact-string version of this transform was wrong four times in a row:
// each fix revealed the next construct, and the last was found by a syntax
// checker rather than by reading. The cases below are those four, plus the one
// that made the output worse than the input: a parameter list whose parentheses
// the replacement dropped, turning a function signature into nonsense that
// parsed as far as the colon.
func TestTheTransformHandlesEveryConstructItMeets(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"function toOrigin(value: string): string {", "function toOrigin(value) {"},
		{"const pluginOrigins: [string, string][] = [", "const pluginOrigins = ["},
		{`{ protocol: "https" as const }`, `{ protocol: "https" }`},
		{`protocol: x.replace(":", "") as "http" | "https",`, `protocol: x.replace(":", ""),`},
		{"function f(a) {", "function f(a) {"},
	} {
		if got := toPlainJavaScript(tc.in); got != tc.want {
			t.Errorf("toPlainJavaScript(%q)\n  = %q\n want %q", tc.in, got, tc.want)
		}
	}
}

// The transform touches the docs config and nothing else.
//
// apps/web and apps/admin have TypeScript configs, and stripping their
// annotations would lose the checking they exist for.
func TestTheTypeScriptConfigsKeepTheirTypes(t *testing.T) {
	if !strings.Contains(webNextConfig(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext}), "function toOrigin(value: string): string {") {
		t.Error("the web app's config lost its typed helper")
	}
	if !strings.Contains(adminNextConfig(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext}), "function toOrigin(value: string): string {") {
		t.Error("the admin's config lost its typed helper")
	}

	// And the transform really is a transform: given TypeScript, it returns
	// JavaScript that says the same thing.
	got := toPlainJavaScript(`const pluginOrigins: [string, string][] = []
const x = { protocol: "https" as const }`)
	if strings.Contains(got, ": [string, string][]") || strings.Contains(got, "as const") {
		t.Errorf("the transform left TypeScript behind:\n%s", got)
	}
	if !strings.Contains(got, "const pluginOrigins = []") || !strings.Contains(got, `protocol: "https"`) {
		t.Errorf("the transform changed what the code does:\n%s", got)
	}
}
