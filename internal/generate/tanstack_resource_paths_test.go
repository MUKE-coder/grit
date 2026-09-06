package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tanStackAdmin makes a project look like a TanStack one to the generator,
// which decides by looking for apps/admin/src/resources.
func tanStackAdmin(t *testing.T, root string) {
	t.Helper()
	for _, dir := range []string{
		filepath.Join("apps", "admin", "src", "resources"),
		filepath.Join("apps", "admin", "src", "routes", "_dashboard", "resources"),
	} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}
	writeTestFile(t, filepath.Join(root, "apps", "admin", "src", "resources", "index.ts"),
		"// grit:resources\nexport {}\n")
}

// The definition and the routes that import it must agree about where it is.
//
// They did not: the definition was written flat as src/resources/posts.ts while
// both route writers imported '@/resources/posts/posts', so every generated
// resource in a Vite project failed to typecheck on the line reaching its own
// definition. Reported as "Cannot find module '@/resources/posts/posts'".
func TestTanStackResourceDefinitionMatchesItsRouteImports(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	tanStackAdmin(t, root)

	def, err := ParseInlineFields("Post", "title:string,body:text,published:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()

	if err := g.writeResourceDefinitionTanStack(names); err != nil {
		t.Fatalf("definition: %v", err)
	}
	if err := g.writeResourcePageTanStack(names); err != nil {
		t.Fatalf("list route: %v", err)
	}
	if err := g.writeResourceDetailPageTanStack(names); err != nil {
		t.Fatalf("detail route: %v", err)
	}

	// The folder layout, matching the Next generator and the built-in blogs and
	// users. The .custom.tsx overlay has to sit beside what it overlays.
	adminSrc := filepath.Join(root, "apps", "admin", "src")
	for _, rel := range []string{
		filepath.Join("resources", "posts", "posts.ts"),
		filepath.Join("resources", "posts", "posts.custom.tsx"),
	} {
		if _, err := os.Stat(filepath.Join(adminSrc, rel)); err != nil {
			t.Errorf("missing %s", filepath.ToSlash(rel))
		}
	}

	// And every route import has to resolve to a file that is actually there.
	routes := filepath.Join(adminSrc, "routes", "_dashboard", "resources", "posts")
	for _, name := range []string{"index.tsx", "$id.tsx"} {
		body := readTestFile(t, filepath.Join(routes, name))
		spec := resourceImportSpec(t, body)
		if spec == "" {
			t.Errorf("%s imports no resource definition", name)
			continue
		}
		rel := strings.TrimPrefix(spec, "@/")
		if _, err := os.Stat(filepath.Join(adminSrc, rel+".ts")); err != nil {
			t.Errorf("%s imports %q, and there is no such file", name, spec)
		}
	}
}

// A project still holding a flat definition from before the fix keeps it, and
// its routes import the flat path rather than one that is not there.
func TestTanStackFlatDefinitionKeepsWorking(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	tanStackAdmin(t, root)

	// The pre-fix layout.
	flat := filepath.Join(root, "apps", "admin", "src", "resources", "posts.ts")
	writeTestFile(t, flat, "export const postResource = {} as never\n")

	def, err := ParseInlineFields("Post", "title:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()

	if err := g.writeResourceDefinitionTanStack(names); err != nil {
		t.Fatalf("definition: %v", err)
	}
	if err := g.writeResourceDetailPageTanStack(names); err != nil {
		t.Fatalf("detail route: %v", err)
	}

	// Left where it was: a second copy in a folder would leave two definitions
	// and a registry pointing at whichever it found first.
	if _, err := os.Stat(flat); err != nil {
		t.Error("the flat definition was moved out from under the project")
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "admin", "src", "resources", "posts", "posts.ts")); err == nil {
		t.Error("wrote a second definition in a folder beside the flat one")
	}

	body := readTestFile(t, filepath.Join(root, "apps", "admin", "src", "routes",
		"_dashboard", "resources", "posts", "$id.tsx"))
	if spec := resourceImportSpec(t, body); spec != "@/resources/posts" {
		t.Errorf("route imports %q, want @/resources/posts to match the flat file", spec)
	}
}

// resourceImportSpec pulls the module specifier off the line importing the
// resource definition.
func resourceImportSpec(t *testing.T, body string) string {
	t.Helper()
	for _, line := range strings.Split(body, "\n") {
		if !strings.Contains(line, "Resource } from") {
			continue
		}
		i := strings.Index(line, "'")
		j := strings.LastIndex(line, "'")
		if i >= 0 && j > i {
			return line[i+1 : j]
		}
	}
	return ""
}
