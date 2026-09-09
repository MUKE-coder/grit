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

// Every client hook uses the shared type, not a copy of it.
//
// "One backend, every client, from shared types" only holds if the clients
// import them. Expo declared its own interface and drifted after grit sync;
// desktop declared `Record<string, unknown> & { id: string }`, which is not
// drift but an absent type: task.titel typechecked.
func TestEveryClientHookUsesTheSharedType(t *testing.T) {
	const module = "app/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Task", "title:string,done:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()

	if err := g.writeMobileHook(names); err != nil {
		t.Fatalf("mobile hook: %v", err)
	}
	mobile := readTestFile(t, filepath.Join(root, "apps", "expo", "hooks", "use-tasks.ts"))
	if !strings.Contains(mobile, `import type { Task } from "@repo/shared/types";`) {
		t.Error("the Expo hook does not import the shared type, so grit sync " +
			"cannot reach it")
	}
	if strings.Contains(mobile, "export interface Task {") {
		t.Error("the Expo hook still declares its own Task, which goes stale " +
			"the moment the Go model changes")
	}
	// The screens import the type from the hook, so it has to be re-exported.
	if !strings.Contains(mobile, "export type { Task };") {
		t.Error("the Expo hook does not re-export Task; the generated screens " +
			"import it from here")
	}

	desktop := g.desktopClientHook(names)
	if !strings.Contains(desktop, `import type { Task } from "@repo/shared/types";`) {
		t.Error("the desktop hook does not import the shared type")
	}
	if strings.Contains(desktop, "= Record<string, unknown> & { id: string }") {
		t.Error("the desktop hook still types the resource as an open record, " +
			"which accepts anything and checks nothing")
	}
}
