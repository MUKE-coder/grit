package scaffold

import (
	"encoding/json"
	"strings"
	"testing"
)

// packages/shared has to be type-checked on its own, not only through whatever
// an app happens to import.
//
// grit#213: schemas/field-formats.ts shipped for two releases with zod 3
// syntax against the zod 4 pin, and the root `pnpm type-check` passed the
// whole time. The reason is structural rather than an oversight. tsc only
// compiles what it is pointed at, and schemas/index.ts does not re-export
// field-formats or phone, so those two files are reached only when a resource
// with a formatted field imports them. A project with no such resource never
// compiled them at all.
//
// The package already had a tsconfig and a local typescript. What it had no
// script, so `turbo type-check` walked straight past it. The CI workflow even
// documented the gap in a comment and worked around it rather than closing it.
func TestSharedPackageTypeChecksItself(t *testing.T) {
	var pkg struct {
		Scripts      map[string]string `json:"scripts"`
		DevDeps      map[string]string `json:"devDependencies"`
		Dependencies map[string]string `json:"dependencies"`
	}
	raw := sharedPackageJSON(tripleOptions())
	if err := json.Unmarshal([]byte(raw), &pkg); err != nil {
		t.Fatalf("packages/shared/package.json is not valid JSON: %v\n%s", err, raw)
	}

	check, ok := pkg.Scripts["type-check"]
	if !ok {
		t.Fatal("packages/shared declares no type-check script, so turbo type-check skips the package " +
			"and every file nothing imports goes uncompiled (grit#213)")
	}
	if !strings.Contains(check, "tsc") {
		t.Errorf("type-check is %q, which does not run the compiler", check)
	}
	// The script needs a compiler in the package, or the recursive form fails
	// with "Command tsc not found" and says nothing about the code.
	if _, ok := pkg.DevDeps["typescript"]; !ok {
		t.Error("packages/shared declares a type-check script and no typescript to run it")
	}

	// The tsconfig has to reach every file, not just the ones with an importer.
	cfg := sharedTSConfig()
	if !strings.Contains(cfg, `"include": ["**/*.ts"]`) {
		t.Errorf("packages/shared/tsconfig.json does not include every .ts file, so the script cannot "+
			"see an unimported one:\n%s", cfg)
	}
}

// The zod 4 move raised the pin and modernised the call sites it could see.
// These are the spellings zod 4 removed outright rather than deprecated: each
// one is a type error, not a warning, and the generator must emit none of them.
//
// z.record is the one that bit. It took a value type in zod 3 and takes a key
// type as well in zod 4, so the single-argument form is TS2554 rather than
// anything a reader would notice.
func TestSharedSchemasUseNoRemovedZodSyntax(t *testing.T) {
	files := map[string]string{
		"schemas/field-formats.ts": sharedFieldFormatsSchema(),
		"schemas/phone.ts":         sharedPhoneSchema(),
		"schemas/user.ts":          sharedUserSchema(),
		"schemas/file-ref.ts":      sharedFileRefSchema(),
	}
	// Spelling, and what zod 4 wants instead.
	removed := []struct{ bad, instead string }{
		{"invalid_type_error", "the error function"},
		{"required_error", "the error function"},
		{".nonempty(", ".min(1)"},
		{"z.string().ip(", "z.ipv4() or z.ipv6()"},
		{"z.string().cidr(", "z.cidrv4() or z.cidrv6()"},
		{"z.ZodType<JsonValue>", "z.json()"},
	}
	for name, body := range files {
		for _, r := range removed {
			if strings.Contains(body, r.bad) {
				t.Errorf("%s uses %s, which zod 4 removed: use %s", name, r.bad, r.instead)
			}
		}
		// A one-argument z.record is the arity change, and it has to be
		// matched on shape rather than on the literal text.
		for _, call := range recordCalls(body) {
			if !strings.Contains(call, ",") {
				t.Errorf("%s calls %s with one argument; zod 4 wants a key type first, "+
					"such as z.record(z.string(), X)", name, call)
			}
		}
	}
}

// recordCalls pulls each z.record(...) call out of a source file, balanced so a
// nested call does not truncate it.
func recordCalls(src string) []string {
	var out []string
	const marker = "z.record("
	for i := 0; ; {
		at := strings.Index(src[i:], marker)
		if at < 0 {
			return out
		}
		start := i + at
		depth := 0
		for j := start + len(marker) - 1; j < len(src); j++ {
			switch src[j] {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					out = append(out, src[start:j+1])
					i = j + 1
					goto next
				}
			}
		}
		return out
	next:
	}
}
