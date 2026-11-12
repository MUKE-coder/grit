package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Giving packages/shared a type-check script in a project that already exists.
//
// The upgrade rewrites packages/shared/package.json from the template, so a
// project that has not touched it gets the script that way. A project that
// added one dependency of its own has edited the file, the upgrade guard holds
// it back for exactly the right reason, and the package stays uncompiled.
//
// That is the same argument raiseHookformResolvers makes, and the same answer:
// edit the one thing, which is a change no reader will mind and which they
// cannot make for themselves without knowing why it matters.
//
// Why it matters: without a script, `turbo type-check` walks past the package
// entirely, and a file nothing imports is never compiled. schemas/
// field-formats.ts shipped zod 3 syntax against a zod 4 pin for two releases
// that way, and the project's own type-check passed throughout (grit#213).

const sharedTypeCheckScript = "tsc --noEmit"

// addSharedTypeCheckScript puts a type-check script in packages/shared's
// package.json when it has none.
func addSharedTypeCheckScript(root string) error {
	path := filepath.Join(root, "packages", "shared", "package.json")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	src := string(raw)

	// Read it rather than pattern-match it: a package.json with a scripts
	// block that has no type-check needs an entry added to that block, and one
	// with no scripts block needs the block. Telling those apart by regex is
	// how an edit lands in the wrong object.
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(raw, &pkg); err != nil {
		fmt.Println("  ⚠ packages/shared/package.json is not valid JSON, so its type-check script was left alone")
		return nil
	}
	if _, ok := pkg.Scripts["type-check"]; ok {
		return nil
	}

	crlf := strings.Contains(src, "\r\n")
	body := strings.ReplaceAll(src, "\r\n", "\n")

	var out string
	switch {
	case pkg.Scripts != nil:
		// A scripts block exists: add one line to it, keeping the rest.
		anchor := "\"scripts\": {"
		at := strings.Index(body, anchor)
		if at < 0 {
			return nil
		}
		insert := at + len(anchor)
		out = body[:insert] + "\n    \"type-check\": \"" + sharedTypeCheckScript + "\"," + body[insert:]
	default:
		// No scripts block. Put one before dependencies, which every generated
		// shared package has, so the result reads the way the template does.
		anchor := "  \"dependencies\": {"
		at := strings.Index(body, anchor)
		if at < 0 {
			return nil
		}
		out = body[:at] + "  \"scripts\": {\n    \"type-check\": \"" + sharedTypeCheckScript + "\"\n  },\n" + body[at:]
	}

	// It has to still parse, and it has to actually have the script. A
	// generated edit to JSON that produces invalid JSON is worse than no edit:
	// pnpm refuses the whole workspace.
	var check struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal([]byte(out), &check); err != nil {
		fmt.Printf("  ⚠ packages/shared/package.json: adding the type-check script did not produce valid JSON (%v), so it was left alone\n", err)
		return nil
	}
	if check.Scripts["type-check"] != sharedTypeCheckScript {
		fmt.Println("  ⚠ packages/shared/package.json: the type-check script did not land where expected, so it was left alone")
		return nil
	}

	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	if err := writeRepaired(path, out); err != nil {
		return err
	}
	fmt.Println("  ✓ packages/shared: type-check runs over the package itself now, so a schema nothing imports is still compiled")
	return nil
}
