package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// A generated hook must use the shared type, not a copy of it.
//
// "One backend, every client, from shared types" is the promise, and a local
// copy breaks it silently: `grit sync` rewrites packages/shared and cannot
// reach a duplicate declaration, so a field added to the Go model arrived in
// the shared type and nowhere else. Nothing errored, because each file stayed
// internally consistent. The web app simply did not know the field existed.
//
// Found by adding a column mid-project and running grit sync, which is the
// exact workflow the multi-client tutorial describes.
func TestGeneratedHookUsesTheSharedType(t *testing.T) {
	const module = "app/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Task", "title:string,done:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeReactQueryHooks(names, "web"); err != nil {
		t.Fatalf("hooks: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "web", "hooks", "use-tasks.ts"))

	if !strings.Contains(src, `import type { Task } from "@repo/shared/types";`) {
		t.Error("the hook does not import the shared Task type, so grit sync " +
			"cannot reach it and the client drifts from the backend")
	}
	// A local declaration is the bug. TasksResponse and UseTasksParams are
	// hook-shaped and belong here; the resource type does not.
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "interface Task {" || trimmed == "export interface Task {" {
			t.Error("the hook redeclares Task locally; a schema change would " +
				"update packages/shared and leave this copy stale")
		}
	}
}

// This is what the scaffold's own use-blogs.ts has always done. The generator
// wrote the same file a different way, which is how the two drifted.
func TestGeneratedHookMatchesTheScaffoldsOwnConvention(t *testing.T) {
	const module = "app/apps/api"
	root := setupMinimalProject(t, module)

	def, _ := ParseInlineFields("Task", "title:string")
	g := newTestGenerator(root, module, def)
	if err := g.writeReactQueryHooks(g.Names(), "web"); err != nil {
		t.Fatalf("hooks: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "web", "hooks", "use-tasks.ts"))

	if !strings.Contains(src, "@repo/shared/types") {
		t.Error("the generated hook does not reference the shared package at " +
			"all, while the scaffold's use-blogs.ts imports its type from it")
	}
}
