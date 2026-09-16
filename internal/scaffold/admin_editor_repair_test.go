package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// ─── M43 and Tiptap 3: the templates ─────────────────────────────────────────

func TestOneEditorOnOneTiptap3Schema(t *testing.T) {
	editor := adminWordEditor()
	for _, want := range []string{
		`import { richTextExtensions } from "@/lib/tiptap-extensions";`,
		"useEditorState(",
		"setContent(value || \"\", { emitUpdate: false })",
		"immediatelyRender: false",
	} {
		if !strings.Contains(editor, want) {
			t.Errorf("word-editor.tsx is missing %q", want)
		}
	}
	field := adminRichTextField()
	if !strings.Contains(field, "<WordEditor") || strings.Contains(field, "@tiptap/") {
		t.Error("the rich text field should render the Word-style editor, not build a Tiptap editor of its own")
	}

	schema := adminTiptapExtensions()
	for _, want := range []string{"StarterKit.configure(", "TextAlign", "Color", "Highlight", "Image", "TableKit", "Placeholder"} {
		if !strings.Contains(schema, want) {
			t.Errorf("lib/tiptap-extensions.ts is missing %s", want)
		}
	}

	// Nothing in the panel may still import a Tiptap 2 package or API.
	stale := regexp.MustCompile(`@tiptap/extension-(table-row|table-cell|table-header|color|underline|link|placeholder)"|import TextStyle from|setContent\([^)]*, (false|true)\)`)
	for path, body := range adminFileMap(t.TempDir(), tripleOptions()) {
		if m := stale.FindString(body); m != "" {
			t.Errorf("%s still uses Tiptap 2: %s", filepath.Base(path), m)
		}
	}
}

func TestEveryPanelPackageJSONPinsTiptap3(t *testing.T) {
	sources := map[string]string{
		"admin (Next)":    adminPackageJSON(tripleOptions()),
		"admin (Vite)":    adminTanStackPackageJSON(viteTripleOptions()),
		"single frontend": singleFrontendPackageJSON(Options{ProjectName: "demo", Architecture: ArchSingle}),
		"double web deps": "{\"dependencies\": {\"a\": \"1\"" + webAdminDependencies(doubleOptions()) + "}}",
	}
	for name, src := range sources {
		for _, pkg := range tiptapPackages {
			if !strings.Contains(src, `"`+pkg+`": "`+tiptapVersion+`"`) {
				t.Errorf("%s does not pin %s to %s", name, pkg, tiptapVersion)
			}
		}
		if regexp.MustCompile(`"@tiptap/[a-z-]+": "\^2`).MatchString(src) {
			t.Errorf("%s still has a Tiptap 2 range", name)
		}
	}
}

func TestTiptapBelowPin(t *testing.T) {
	cases := map[string]bool{
		"^2.1.0":            true,
		"2.27.3":            true,
		"~3.30.4":           true,
		tiptapVersion:       false,
		"^" + tiptapVersion: false,
		"3.40.0":            false,
		"4.0.0-beta.1":      false,
		"latest":            false,
		"workspace:*":       false,
	}
	for version, want := range cases {
		if got := tiptapBelowPin(version); got != want {
			t.Errorf("tiptapBelowPin(%q) = %v, want %v", version, got, want)
		}
	}
}

// oldAdminPackageJSON is the admin package.json's Tiptap block as Grit wrote it
// before Tiptap 3.
const oldTiptapBlock = `    "@tiptap/extension-link": "^2.1.0",
    "@tiptap/extension-text-align": "^2.1.0",
    "@tiptap/extension-text-style": "^2.1.0",
    "@tiptap/extension-color": "^2.1.0",
    "@tiptap/extension-highlight": "^2.1.0",
    "@tiptap/extension-underline": "^2.1.0",
    "@tiptap/extension-image": "^2.1.0",
    "@tiptap/extension-table": "^2.1.0",
    "@tiptap/extension-table-row": "^2.1.0",
    "@tiptap/extension-table-cell": "^2.1.0",
    "@tiptap/extension-table-header": "^2.1.0",
    "@tiptap/extension-placeholder": "^2.1.0",
    "@tiptap/pm": "^2.1.0",
    "@tiptap/react": "^2.1.0",
    "@tiptap/starter-kit": "^2.1.0",
`

func TestRepairTiptapPackageJSON(t *testing.T) {
	old := "{\n  \"dependencies\": {\n    \"axios\": \"^1.6.0\",\n" + oldTiptapBlock + "    \"zod\": \"^3.22.0\"\n  }\n}\n"
	out, fixed, warn := repairTiptapPackageJSON(old)
	if len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("fixed %v, warnings %v", fixed, warn)
	}
	if strings.Contains(out, "^2.1.0") {
		t.Error("a Tiptap 2 range survived")
	}
	for _, pkg := range append(tiptapPackages, "@tiptap/extension-table-row", "@tiptap/extension-link") {
		if !strings.Contains(out, `"`+pkg+`": "`+tiptapVersion+`"`) {
			t.Errorf("%s is not at %s", pkg, tiptapVersion)
		}
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("the repaired package.json is not JSON: %v\n%s", err, out)
	}
	if again, fixed, _ := repairTiptapPackageJSON(out); again != out || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}

	// The Tiptap block last in its object: new lines must not leave a
	// trailing comma.
	last := "{\n  \"dependencies\": {\n    \"@tiptap/react\": \"^2.1.0\"\n  }\n}\n"
	out, _, _ = repairTiptapPackageJSON(last)
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("repairing a last-line block broke the JSON: %v\n%s", err, out)
	}

	// A project already past the pin keeps its version.
	ahead := "{\"dependencies\": {\"@tiptap/react\": \"3.40.0\"}}"
	out, _, _ = repairTiptapPackageJSON(ahead)
	if !strings.Contains(out, `"@tiptap/react": "3.40.0"`) {
		t.Error("the repair lowered a newer Tiptap")
	}
}

func TestRepairTiptapLeavesAnUnmigratedEditorAlone(t *testing.T) {
	root := t.TempDir()
	admin := filepath.Join(root, "apps", "admin")
	pkg := "{\n  \"dependencies\": {\n" + oldTiptapBlock + "    \"zod\": \"^3.22.0\"\n  }\n}\n"
	writeFileT(t, filepath.Join(admin, "package.json"), pkg)
	writeFileT(t, filepath.Join(admin, "components", "forms", "word-editor.tsx"), `import Table from "@tiptap/extension-table";`)
	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := repairTiptap(root, m); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(admin, "package.json")); string(got) != pkg {
		t.Error("Tiptap was raised under an editor still written for Tiptap 2")
	}

	writeFileT(t, filepath.Join(admin, "components", "forms", "word-editor.tsx"), adminWordEditor())
	writeFileT(t, filepath.Join(admin, "components", "forms", "fields", "rich-text-field.tsx"), adminRichTextField())
	if err := repairTiptap(root, m); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(admin, "package.json")); strings.Contains(string(got), "^2.1.0") {
		t.Error("Tiptap was not raised once the editor was on Tiptap 3")
	}
}

// ─── M42: the blog on its resource definition ────────────────────────────────

func TestBlogRunsOnItsResourceDefinition(t *testing.T) {
	files := adminFileMap(t.TempDir(), tripleOptions())
	for path, body := range files {
		slashed := filepath.ToSlash(path)
		if strings.HasSuffix(slashed, "resources/blogs/page.tsx") && !strings.Contains(body, "<ResourcePage resource={blogsResource} />") {
			t.Error("the blog list page is not <ResourcePage>")
		}
		if strings.HasSuffix(slashed, "resources/blogs/[id]/page.tsx") && !strings.Contains(body, "<ResourceDetailPage resource={blogsResource}") {
			t.Error("the blog detail page is not <ResourceDetailPage>")
		}
		if strings.Contains(body, legacyBlogsListMarker) {
			t.Errorf("%s still fetches a single page of 100 blogs", slashed)
		}
	}
	if !strings.Contains(adminBlogsResource(), `formView: "page",`) {
		t.Error("the blog resource does not give its editor a page")
	}
	for path := range adminTanStackFileMap(t.TempDir(), viteTripleOptions()) {
		if strings.HasSuffix(filepath.ToSlash(path), "pages/resources/blog-detail.tsx") {
			t.Error("the Vite admin still writes the hand-written blog detail page")
		}
	}
}

func TestRepairBlogsResourceSource(t *testing.T) {
	old := strings.Replace(adminBlogsResource(), blogsResourceFormView, "", 1)
	out, fixed, _ := repairBlogsResourceSource(old)
	if out != adminBlogsResource() || len(fixed) != 1 {
		t.Fatalf("the repaired definition is not the template:\n%s", out)
	}
	if again, fixed, _ := repairBlogsResourceSource(out); again != out || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
}

func TestRepairBlogPages(t *testing.T) {
	root := t.TempDir()
	admin := filepath.Join(root, "apps", "admin")
	shape := adminPanelShape{code: admin, routes: filepath.Join(admin, "app"), alias: "@"}
	list := filepath.Join(shape.routes, "(dashboard)", "resources", "blogs", "page.tsx")
	detail := filepath.Join(shape.routes, "(dashboard)", "resources", "blogs", "[id]", "page.tsx")

	// Edited: named, left alone.
	writeFileT(t, list, legacyAdminBlogsListPage())
	writeFileT(t, detail, legacyAdminBlogDetailPage())
	if err := repairBlogPages(root, shape, func(string) bool { return false }); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(list); string(got) != legacyAdminBlogsListPage() {
		t.Error("an edited blog list was replaced")
	}

	// Grit's: replaced.
	if err := repairBlogPages(root, shape, func(string) bool { return true }); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(list); string(got) != adminBlogsPage() {
		t.Errorf("the blog list is not <ResourcePage>:\n%s", got)
	}
	if got, _ := os.ReadFile(detail); string(got) != adminResourceDetailRoute("blogs", "blogs", "Blogs") {
		t.Errorf("the blog detail page is not <ResourceDetailPage>:\n%s", got)
	}

	// Idempotent: nothing left to recognise.
	if err := repairBlogPages(root, shape, func(string) bool { t.Error("asked about a page already replaced"); return true }); err != nil {
		t.Fatal(err)
	}
}

func TestRepairBlogPagesInsideTheWebApp(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web")
	shape := adminPanelShape{code: filepath.Join(web, "admin-panel"), routes: filepath.Join(web, "app", "admin"), alias: "@admin"}
	list := filepath.Join(shape.routes, "(dashboard)", "resources", "blogs", "page.tsx")
	writeFileT(t, list, embedAdminContent(legacyAdminBlogsListPage(), adminRoutePrefixes))
	if err := repairBlogPages(root, shape, func(string) bool { return true }); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(list)
	if string(got) != embedAdminContent(adminBlogsPage(), adminRoutePrefixes) || !strings.Contains(string(got), `"@admin/resources/blogs/blogs"`) {
		t.Errorf("the embedded blog list is not the embedded template:\n%s", got)
	}
}

// ─── M44: one page header ─────────────────────────────────────────────────────

func TestOnePageHeader(t *testing.T) {
	for path, body := range adminFileMap(t.TempDir(), tripleOptions()) {
		if strings.HasSuffix(filepath.ToSlash(path), legacyPageHeaderRel) {
			t.Error("the scaffold still writes components/layout/page-header.tsx")
		}
		if strings.Contains(body, legacyPageHeaderImport) {
			t.Errorf("%s imports the old page header", filepath.Base(path))
		}
	}
	header := adminPageHeaderComponent()
	if !strings.Contains(header, "stats?: StatCard[];") || !strings.Contains(header, "<StatCards stats={stats} />") {
		t.Error("chrome/PageHeader.tsx does not take stat cards")
	}
	if !strings.Contains(adminResourcePage(), `import { PageHeader } from "@/components/chrome/PageHeader";`) {
		t.Error("resource-page.tsx does not use the chrome page header")
	}
}

func TestResourceControllerIsComposed(t *testing.T) {
	controller := adminUseResourceController()
	for _, hook := range []string{"useResourceURLState(", "useResourceSelection<T>(", "useResourceDialogs<T>("} {
		if !strings.Contains(controller, hook) {
			t.Errorf("the controller does not use %s", hook)
		}
	}
	if n := strings.Count(controller, "useState("); n > 1 {
		t.Errorf("the controller holds %d useStates of its own; its state belongs in the hooks it is built from", n)
	}
	files := adminFileMap(t.TempDir(), tripleOptions())
	for _, rel := range []string{"use-resource-url-state.ts", "use-resource-selection.ts", "use-resource-dialogs.ts"} {
		found := false
		for path := range files {
			if strings.HasSuffix(filepath.ToSlash(path), "hooks/"+rel) {
				found = true
			}
		}
		if !found {
			t.Errorf("the scaffold does not write hooks/%s", rel)
		}
	}
}

func TestPruneLegacyPageHeader(t *testing.T) {
	root := t.TempDir()
	admin := filepath.Join(root, "apps", "admin")
	shape := adminPanelShape{code: admin, routes: filepath.Join(admin, "app"), alias: "@"}
	header := filepath.Join(admin, filepath.FromSlash(legacyPageHeaderRel))
	importer := filepath.Join(admin, "app", "(dashboard)", "custom", "page.tsx")
	writeFileT(t, header, "export function PageHeader() { return null; }\n")
	writeFileT(t, importer, `import { PageHeader } from "@/components/layout/page-header";`+"\n")

	pruneLegacyPageHeader(root, shape, func(string) bool { return true })
	if !fileExists(header) {
		t.Fatal("the old header was removed while a page still imports it")
	}
	if err := os.Remove(importer); err != nil {
		t.Fatal(err)
	}
	pruneLegacyPageHeader(root, shape, func(string) bool { return false })
	if !fileExists(header) {
		t.Fatal("an edited old header was removed")
	}
	pruneLegacyPageHeader(root, shape, func(string) bool { return true })
	if fileExists(header) {
		t.Error("the unused old header is still there")
	}
}

// ─── M41: shared model types and profile schemas ─────────────────────────────

func TestAdminImportsSharedModelTypes(t *testing.T) {
	for name, src := range map[string]string{
		"support list":      adminSupportListPage(),
		"ticket thread":     adminTicketThreadPage(),
		"use-notifications": adminUseNotifications(),
	} {
		if !strings.Contains(src, `from "@repo/shared/types";`) {
			t.Errorf("%s does not import its row type from @repo/shared/types", name)
		}
		if regexp.MustCompile(`(?m)^(export )?interface (Ticket|Notification|Reply) \{`).MatchString(src) {
			t.Errorf("%s still declares its own model type", name)
		}
	}
	profile := adminCaptivatingProfile()
	if strings.Contains(profile, "z.object(") || !strings.Contains(profile, `} from "@repo/shared/schemas";`) {
		t.Error("the profile page still declares its own Zod schemas")
	}
	index := sharedTypesIndex()
	for _, typ := range SharedModelTypes {
		if !strings.Contains(index, "export type { "+typ.Name+" } from \"./"+typ.Kebab+"\";") {
			t.Errorf("types/index.ts does not export %s", typ.Name)
		}
	}
	if !strings.Contains(sharedSchemasIndex(), sharedProfileSchemaExport) {
		t.Error("schemas/index.ts does not export the profile schemas")
	}
}

func TestRepairSharedTypes(t *testing.T) {
	root := t.TempDir()
	shared := filepath.Join(root, "packages", "shared")
	oldTypes := strings.Replace(sharedTypesIndex(), sharedModelTypeExports(), "", 1)
	oldSchemas := strings.Replace(sharedSchemasIndex(), sharedProfileSchemaExport, "", 1)
	writeFileT(t, filepath.Join(shared, "types", "index.ts"), oldTypes)
	writeFileT(t, filepath.Join(shared, "schemas", "index.ts"), oldSchemas)
	// A type grit sync already wrote, perhaps from a model the project changed.
	synced := "export interface Ticket {\n  id: string;\n  extra: string;\n}\n"
	writeFileT(t, filepath.Join(shared, "types", "ticket.ts"), synced)

	if err := repairSharedTypes(root); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(shared, "types", "index.ts")); string(got) != sharedTypesIndex() {
		t.Errorf("types/index.ts is not the template:\n%s", got)
	}
	if got, _ := os.ReadFile(filepath.Join(shared, "schemas", "index.ts")); string(got) != sharedSchemasIndex() {
		t.Errorf("schemas/index.ts is not the template:\n%s", got)
	}
	if got, _ := os.ReadFile(filepath.Join(shared, "types", "ticket.ts")); string(got) != synced {
		t.Error("an existing shared type was overwritten")
	}
	if got, _ := os.ReadFile(filepath.Join(shared, "types", "api-key.ts")); string(got) != SharedModelTypes[0].Source {
		t.Error("a missing shared type was not added")
	}

	if err := repairSharedTypes(root); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(shared, "types", "index.ts")); string(got) != sharedTypesIndex() {
		t.Error("the repair is not idempotent")
	}
}
