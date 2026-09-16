package scaffold

import (
	"strings"
	"testing"
)

// M24: a QueryClient at module scope is one cache for every request the server
// renders. M25: Tiptap and the form stack rode in the first load of every form,
// detail page and the blog editor.

func TestAdminQueryClientIsCreatedPerMount(t *testing.T) {
	src := adminQueryClient()
	if strings.Contains(src, "export const queryClient") || strings.Contains(src, "= new QueryClient(") {
		t.Error("lib/query-client.ts still creates a client at module scope")
	}
	if !strings.Contains(src, "export function makeQueryClient(): QueryClient {") {
		t.Error("lib/query-client.ts has no makeQueryClient factory")
	}

	for name, providers := range map[string]string{
		"Next.js admin":  adminProviders(),
		"embedded admin": embedAdminContent(adminProviders(), adminRoutePrefixes),
		"Vite admin":     nextToTanStack(adminProviders()),
	} {
		if !strings.Contains(providers, "const [queryClient] = useState(makeQueryClient);") {
			t.Errorf("the %s's provider does not create its client per mount", name)
		}
		if strings.Contains(providers, "import { queryClient }") {
			t.Errorf("the %s's provider still imports a module-level client", name)
		}
	}
	if !strings.Contains(embedAdminContent(adminProviders(), adminRoutePrefixes), `import { makeQueryClient } from "@admin/lib/query-client";`) {
		t.Error("the embedded provider's import was not repointed at admin-panel")
	}
}

func TestWebProvidersLeaveTheAdminSectionToItsOwnClient(t *testing.T) {
	double := webProvidersFor(Options{Architecture: ArchDouble})
	for _, want := range []string{
		`useSelectedLayoutSegment() === "admin"`,
		"<QueryClientContext.Provider value={active}>",
		"inAdmin ? null : makeQueryClient()",
		"active.mount();",
		"return () => active.unmount();",
	} {
		if !strings.Contains(double, want) {
			t.Errorf("the double's web providers are missing %q", want)
		}
	}
	if webProvidersFor(Options{Architecture: ArchTriple}) != webProviders() {
		t.Error("a triple's web app has no admin section, so its providers must not change")
	}
}

func TestFormsLoadTheEditorAndFormStackOnDemand(t *testing.T) {
	cases := []struct {
		name, src, static, lazy string
	}{
		{"form builder", adminFormBuilder(), `import { RichTextField } from`, `import("./fields/rich-text-field")`},
		{"resource detail page", adminResourceDetailPage(), `import { FormSheet } from`, `import("@/components/forms/form-sheet")`},
		{"blog editor page", legacyAdminBlogDetailPage(), `import { WordEditor } from`, `import("@/components/forms/word-editor")`},
	}
	for _, c := range cases {
		if strings.Contains(c.src, c.static) {
			t.Errorf("the %s still imports statically: %s", c.name, c.static)
		}
		if !strings.Contains(c.src, c.lazy) {
			t.Errorf("the %s does not load %s", c.name, c.lazy)
		}
		if !strings.Contains(c.src, `import dynamic from "next/dynamic";`) {
			t.Errorf("the %s does not import dynamic", c.name)
		}
		if !strings.Contains(c.src, "<Suspense") {
			t.Errorf("the %s renders a lazy component without Suspense; the Vite admin's dynamic() is React.lazy", c.name)
		}
		if vite := nextToTanStack(c.src); !strings.Contains(vite, `import { dynamic } from "@/lib/next-compat";`) {
			t.Errorf("the Vite %s does not take dynamic from the compat shim", c.name)
		}
	}
	if !strings.Contains(adminFormBuilder(), "{ ssr: false, loading: RichTextPlaceholder }") ||
		!strings.Contains(legacyAdminBlogDetailPage(), "{ ssr: false, loading: EditorPlaceholder }") {
		t.Error("the editors must not render on the server: Tiptap's DOM never matches React's")
	}
}

// oldBlogEditorPage is the blog editor page as Grit wrote it before M25.
func oldBlogEditorPage() string {
	src := legacyAdminBlogDetailPage()
	src = strings.Replace(src, blogEditorNewReactImport, blogEditorOldReactImport, 1)
	src = strings.Replace(src, `import { PageHeader } from "@/components/chrome/PageHeader";
`, `import { PageHeader } from "@/components/chrome/PageHeader";
`+blogEditorOldStaticImport, 1)
	src = strings.Replace(src, blogEditorDynamic, "", 1)
	return strings.Replace(src, blogEditorNewUsage, blogEditorOldUsage, 1)
}

func TestRepairBlogEditorProducesTheTemplate(t *testing.T) {
	triple := adminPanelShape{alias: "@"}
	old := oldBlogEditorPage()
	if !strings.Contains(old, blogEditorOldStaticImport) || strings.Contains(old, "dynamic") {
		t.Fatal("the reconstructed old page is not the old page")
	}
	out, fixed, warn := repairBlogEditorSource(old, triple)
	if len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("not repaired: %v %v", fixed, warn)
	}
	if out != legacyAdminBlogDetailPage() {
		t.Errorf("the repaired page is not the template:\n%s", out)
	}
	if again, fixed, _ := repairBlogEditorSource(out, triple); again != out || len(fixed) != 0 {
		t.Error("the blog editor repair is not idempotent")
	}

	// Inside the web app the page imports through @admin.
	double := adminPanelShape{alias: "@admin"}
	embeddedOld := embedAdminContent(old, adminRoutePrefixes)
	out, fixed, warn = repairBlogEditorSource(embeddedOld, double)
	if len(fixed) != 1 || len(warn) != 0 || out != embedAdminContent(legacyAdminBlogDetailPage(), adminRoutePrefixes) {
		t.Errorf("the embedded page was not repaired to its template (%v %v):\n%s", fixed, warn, out)
	}

	edited := strings.Replace(old, "minHeight={500}", "minHeight={600}", 1)
	if out, _, warn := repairBlogEditorSource(edited, triple); out != edited || len(warn) != 1 {
		t.Error("an edited editor block must be left alone with a warning")
	}
}

func TestRepairAdminQueryClientProducesTheTemplate(t *testing.T) {
	out, fixed, warn := repairAdminQueryClientSource(adminQueryClientOldSource)
	if len(fixed) != 1 || len(warn) != 0 || out != adminQueryClient() {
		t.Fatalf("the query client was not repaired to the template (%v %v):\n%s", fixed, warn, out)
	}
	if again, fixed, _ := repairAdminQueryClientSource(out); again != out || len(fixed) != 0 {
		t.Error("the query client repair is not idempotent")
	}

	// Somebody tuned the options: they survive, in a factory.
	tuned := strings.Replace(adminQueryClientOldSource, "staleTime: 5 * 60 * 1000", "staleTime: 30 * 1000", 1)
	out, fixed, warn = repairAdminQueryClientSource(tuned)
	if len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("the tuned query client was not repaired: %v %v", fixed, warn)
	}
	for _, want := range []string{"export function makeQueryClient(): QueryClient {", "staleTime: 30 * 1000", "  return new QueryClient({\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("the tuned repair is missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "export const queryClient") || !strings.HasSuffix(out, "  });\n}\n") {
		t.Errorf("the tuned repair left a module-level client or a broken ending:\n%s", out)
	}

	other := "import { QueryClient } from \"@tanstack/react-query\";\nexport const client = new QueryClient();\n"
	if out, _, warn := repairAdminQueryClientSource(other); out != other || len(warn) != 1 {
		t.Error("an unrecognised query-client.ts must be left alone with a warning")
	}
}

func TestRepairAdminProvidersProducesTheTemplate(t *testing.T) {
	old := strings.Replace(adminProviders(), adminProvidersNewImports, adminProvidersOldImports, 1)
	old = strings.Replace(old, adminProvidersNewBody, adminProvidersOldBody, 1)
	if !strings.Contains(old, `import { queryClient } from "@/lib/query-client";`) {
		t.Fatal("the reconstructed old provider is not the old provider")
	}

	triple := adminPanelShape{alias: "@"}
	out, fixed, warn := repairAdminProvidersSource(old, triple)
	if len(fixed) != 1 || len(warn) != 0 || out != adminProviders() {
		t.Fatalf("the provider was not repaired to the template (%v %v):\n%s", fixed, warn, out)
	}
	if again, fixed, _ := repairAdminProvidersSource(out, triple); again != out || len(fixed) != 0 {
		t.Error("the provider repair is not idempotent")
	}

	double := adminPanelShape{alias: "@admin"}
	embeddedOld := embedAdminContent(old, adminRoutePrefixes)
	out, fixed, _ = repairAdminProvidersSource(embeddedOld, double)
	if len(fixed) != 1 || out != embedAdminContent(adminProviders(), adminRoutePrefixes) {
		t.Errorf("the embedded provider was not repaired to its template:\n%s", out)
	}

	edited := strings.Replace(old, "export function Providers({ children }: { children: React.ReactNode }) {\n  return (", "export function Providers({ children }: { children: React.ReactNode }) {\n  const x = 1;\n  return (", 1)
	if out, _, warn := repairAdminProvidersSource(edited, triple); out != edited || len(warn) != 1 {
		t.Error("an edited provider must be left alone with a warning")
	}
}
