package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ─── L34: dead admin files are not written ───────────────────────────────────

func adminMapsByShape(root string) map[string]map[string]string {
	return map[string]map[string]string{
		"triple (Next)":  adminFileMap(root, tripleOptions()),
		"triple (Vite)":  adminTanStackFileMap(root, viteTripleOptions()),
		"double (Next)":  embeddedAdminFileMap(root, doubleOptions()),
		"double (Vite)":  embeddedSingleAdminFileMap(root, viteDoubleOptions()),
		"single (SPA)":   embeddedSingleAdminFileMap(root, singleOptions()),
		"modern (Next)":  adminFileMap(root, Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext, Style: "modern"}),
		"glass (Vite)":   adminTanStackFileMap(root, Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendTanStack, Style: "glass"}),
		"minimal (Next)": adminFileMap(root, Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext, Style: "minimal"}),
	}
}

func TestAdminDoesNotWriteDeadFiles(t *testing.T) {
	root := t.TempDir()
	for shape, files := range adminMapsByShape(root) {
		widgets := 0
		importsWidgets := false
		for path, body := range files {
			slashed := filepath.ToSlash(path)
			for _, dead := range []string{"components/layout/sidebar.tsx", "components/resource/view-modal.tsx"} {
				if strings.HasSuffix(slashed, dead) {
					t.Errorf("%s still writes %s", shape, dead)
				}
			}
			if strings.Contains(slashed, "/components/widgets/") {
				widgets++
			}
			if strings.Contains(body, "components/widgets/") {
				importsWidgets = true
			}
		}
		usesWidgets := strings.HasPrefix(shape, "modern") || strings.HasPrefix(shape, "glass") || strings.HasPrefix(shape, "minimal")
		switch {
		case usesWidgets && (widgets != 4 || !importsWidgets):
			t.Errorf("%s: its dashboard needs the widgets: %d written, imported: %v", shape, widgets, importsWidgets)
		case !usesWidgets && (widgets != 0 || importsWidgets):
			t.Errorf("%s: %d widget files written for a dashboard that draws its own (imported: %v)", shape, widgets, importsWidgets)
		}
	}
}

func TestEmbeddedPanelSharesTheWebAppsRealtimeClient(t *testing.T) {
	cases := []struct {
		name       string
		opts       Options
		want, gone []string
	}{
		{
			name: "double (Next)",
			opts: doubleOptions(),
			want: []string{"apps/web/lib/realtime.ts", "apps/web/hooks/use-realtime.ts"},
			gone: []string{"apps/web/admin-panel/lib/realtime.ts", "apps/web/admin-panel/hooks/use-realtime.ts"},
		},
		{
			name: "double (Vite)",
			opts: viteDoubleOptions(),
			want: []string{"apps/web/src/lib/realtime.ts", "apps/web/src/hooks/use-realtime.ts"},
			gone: []string{"apps/web/src/admin-panel/lib/realtime.ts", "apps/web/src/admin-panel/hooks/use-realtime.ts"},
		},
		{
			name: "triple (Next)",
			opts: tripleOptions(),
			want: []string{"apps/web/lib/realtime.ts", "apps/admin/lib/realtime.ts", "apps/admin/hooks/use-realtime.ts"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := writeRealtimeClientFiles(root, tc.opts); err != nil {
				t.Fatal(err)
			}
			for _, rel := range tc.want {
				if !fileExists(filepath.Join(root, filepath.FromSlash(rel))) {
					t.Errorf("%s is missing", rel)
				}
			}
			for _, rel := range tc.gone {
				if fileExists(filepath.Join(root, filepath.FromSlash(rel))) {
					t.Errorf("%s is a second copy of the web app's client", rel)
				}
			}
		})
	}
}

// ─── L34: upgrade removes an existing project's copies ───────────────────────

func nextAdminPanel(root string) panelTree {
	admin := filepath.Join(root, "apps", "admin")
	return panelTree{code: admin, app: admin, aliases: map[string]string{"@/": admin}}
}

func TestPruneRemovesUnusedPristineFiles(t *testing.T) {
	root := t.TempDir()
	panel := nextAdminPanel(root)
	admin := panel.code
	for _, dead := range deadPanelFiles {
		writeFileT(t, filepath.Join(admin, filepath.FromSlash(dead.rel)), "export const x = 1;\n")
	}
	// The grid imports the other widgets, which must not keep them alive once
	// the grid itself is gone.
	writeFileT(t, filepath.Join(admin, "components", "widgets", "widget-grid.tsx"),
		"import { StatsCard } from \"./stats-card\";\nimport { ChartWidget } from \"./chart-widget\";\nimport { ActivityWidget } from \"./activity-widget\";\n")
	// A comment naming a file is not an import.
	writeFileT(t, filepath.Join(admin, "components", "ui", "Skeleton.tsx"), "// placeholder for stats-card rows, see components/layout/sidebar.tsx\n")

	if err := pruneDeadPanelFiles(root, []panelTree{panel}, func(string) bool { return true }); err != nil {
		t.Fatal(err)
	}
	for _, dead := range deadPanelFiles {
		if fileExists(filepath.Join(admin, filepath.FromSlash(dead.rel))) {
			t.Errorf("%s is still there", dead.rel)
		}
	}
	if dirExists(filepath.Join(admin, "components", "widgets")) {
		t.Error("the emptied widgets directory was left behind")
	}
	if !fileExists(filepath.Join(admin, "components", "ui", "Skeleton.tsx")) {
		t.Error("a file that is not on the list was removed")
	}
}

func TestPruneKeepsImportedAndEditedFiles(t *testing.T) {
	root := t.TempDir()
	panel := nextAdminPanel(root)
	admin := panel.code
	for _, dead := range deadPanelFiles {
		writeFileT(t, filepath.Join(admin, filepath.FromSlash(dead.rel)), "export const x = 1;\n")
	}
	writeFileT(t, filepath.Join(admin, "components", "widgets", "widget-grid.tsx"), "import { StatsCard } from \"./stats-card\";\n")
	// A modern dashboard renders the grid; someone's page kept the old modal,
	// imported relatively; and the sidebar was edited.
	writeFileT(t, filepath.Join(admin, "app", "(dashboard)", "dashboard", "page.tsx"),
		"import { WidgetGrid } from \"@/components/widgets/widget-grid\";\n")
	writeFileT(t, filepath.Join(admin, "components", "resource", "my-page.tsx"),
		"import { ViewModal } from \"./view-modal\";\n")
	sidebar := filepath.Join(admin, "components", "layout", "sidebar.tsx")
	pristine := func(path string) bool { return !sameFilePath(path, sidebar) }

	if err := pruneDeadPanelFiles(root, []panelTree{panel}, pristine); err != nil {
		t.Fatal(err)
	}
	for _, kept := range []string{
		"components/layout/sidebar.tsx",
		"components/resource/view-modal.tsx",
		"components/widgets/widget-grid.tsx",
		"components/widgets/stats-card.tsx",
	} {
		if !fileExists(filepath.Join(admin, filepath.FromSlash(kept))) {
			t.Errorf("%s was removed although it is imported or edited", kept)
		}
	}
	// Nothing imports these two: the grid imports only the stats card.
	for _, gone := range []string{"components/widgets/chart-widget.tsx", "components/widgets/activity-widget.tsx"} {
		if fileExists(filepath.Join(admin, filepath.FromSlash(gone))) {
			t.Errorf("%s is unused and still there", gone)
		}
	}
}

func TestPruneRemovesAnEmbeddedPanelsSecondRealtimeClient(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web")
	code := filepath.Join(web, "admin-panel")
	panel := panelTree{code: code, app: web, aliases: map[string]string{"@/": web, "@admin/": code}}
	hook := filepath.Join(code, "hooks", "use-realtime.ts")
	client := filepath.Join(code, "lib", "realtime.ts")
	writeFileT(t, hook, "import { subscribe } from \"@admin/lib/realtime\";\n")
	writeFileT(t, client, "export function subscribe() {}\n")
	all := func(string) bool { return true }

	// No client of the web app's own: the panel's copy is the only one.
	if err := pruneDeadPanelFiles(root, []panelTree{panel}, all); err != nil {
		t.Fatal(err)
	}
	if !fileExists(hook) || !fileExists(client) {
		t.Fatal("removed the only realtime client the app had")
	}

	writeFileT(t, filepath.Join(web, "lib", "realtime.ts"), "export function subscribe() {}\n")
	// A page of the web app importing the panel's copy keeps it.
	page := filepath.Join(web, "app", "live", "page.tsx")
	writeFileT(t, page, "import { useRealtime } from \"@admin/hooks/use-realtime\";\n")
	if err := pruneDeadPanelFiles(root, []panelTree{panel}, all); err != nil {
		t.Fatal(err)
	}
	if !fileExists(hook) || !fileExists(client) {
		t.Fatal("removed a realtime client the web app imports")
	}

	if err := os.Remove(page); err != nil {
		t.Fatal(err)
	}
	if err := pruneDeadPanelFiles(root, []panelTree{panel}, all); err != nil {
		t.Fatal(err)
	}
	if fileExists(hook) || fileExists(client) {
		t.Error("the panel's second realtime client is still there")
	}
	if !fileExists(filepath.Join(web, "lib", "realtime.ts")) {
		t.Error("the web app's own client was removed")
	}
}

func TestPruneFindsEveryPanelUpgradeMaintains(t *testing.T) {
	root := t.TempDir()
	writeFileT(t, filepath.Join(root, "apps", "admin", "next.config.ts"), "")
	writeFileT(t, filepath.Join(root, "apps", "web", "admin-panel", "components", "shared", "providers.tsx"), "")
	writeFileT(t, filepath.Join(root, "frontend", "src", "admin-panel", "lib", "utils.ts"), "")
	got := map[string]string{}
	for _, p := range deadCodePanels(root) {
		rel, err := filepath.Rel(root, p.code)
		if err != nil {
			t.Fatal(err)
		}
		appRel, err := filepath.Rel(root, p.app)
		if err != nil {
			t.Fatal(err)
		}
		got[filepath.ToSlash(rel)] = filepath.ToSlash(appRel)
	}
	want := map[string]string{
		"apps/admin":               "apps/admin",
		"apps/web/admin-panel":     "apps/web",
		"frontend/src/admin-panel": "frontend/src",
	}
	for code, app := range want {
		if got[code] != app {
			t.Errorf("panel %s: scanned %q for imports, want %q", code, got[code], app)
		}
	}
	if len(got) != len(want) {
		t.Errorf("found %d panels, want %d: %v", len(got), len(want), got)
	}
}

// ─── L35: no avoidable any ───────────────────────────────────────────────────

var explicitAnyRe = regexp.MustCompile(`:\s*any\b|<any>|\bany\[\]|\bas any\b|<any,|,\s*any>|=>\s*any\b`)

// allowedAny are the uses left on purpose, each explained where it stands.
var allowedAny = []string{
	// React's own constraint for a lazy component (React.lazy is declared this way).
	"export function dynamic<T extends React.ComponentType<any>>(",
	// A realtime payload is any JSON; unknown would break every handler that
	// annotates it or reads a field, in code that compiles today.
	"export type Handler = (payload: any, event: RealtimeEvent) => void;",
}

func TestFrontendTemplatesHaveNoAvoidableAny(t *testing.T) {
	root := t.TempDir()
	sources := map[string]map[string]string{
		"admin (Next)": adminFileMap(root, tripleOptions()),
		"admin (Vite)": adminTanStackFileMap(root, viteTripleOptions()),
		"web (Next)":   webFileMap(root, tripleOptions()),
	}
	sources["extras"] = map[string]string{
		"realtime.ts":     realtimeClientTS(false),
		"use-realtime.ts": useRealtimeTS(true),
		"blog (Vite)":     webTanStackBlogListRoute(),
	}
	for name, files := range sources {
		for path, body := range files {
			for i, line := range strings.Split(body, "\n") {
				code := line
				if idx := strings.Index(code, "//"); idx >= 0 && !strings.Contains(code[:idx], `"`) {
					code = code[:idx]
				}
				if !explicitAnyRe.MatchString(code) {
					continue
				}
				allowed := false
				for _, ok := range allowedAny {
					if strings.Contains(line, ok) {
						allowed = true
					}
				}
				if !allowed {
					t.Errorf("%s %s:%d uses any: %s", name, filepath.Base(path), i+1, strings.TrimSpace(line))
				}
			}
		}
	}
}

// ─── L36: named dropzones ────────────────────────────────────────────────────

func TestDropzoneIsComposedFromParts(t *testing.T) {
	src := adminDropzone()
	for _, want := range []string{
		"export function AvatarDropzone(props: AvatarDropzoneProps)",
		"export function InlineDropzone(props: DropzoneBaseProps)",
		"export function BoxDropzone(",
		"export function CompactDropzone(",
		"export function MinimalDropzone(props: DropzoneBaseProps)",
		"Root: DropzoneRoot,",
		"FileList: DropzoneFileList,",
		"Progress: DropzoneProgress,",
		"Target: DropzoneTarget,",
		"export function useDropzoneContext(): DropzoneState",
		// The old entry point still takes every prop it took.
		"export interface DropzoneProps extends DropzoneBaseProps {",
		"provider?: DropzoneProvider;",
		"variant?: DropzoneVariant;",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("dropzone.tsx is missing %q", want)
		}
	}
	// The avatar holds one picture, whatever it is given.
	if !strings.Contains(src, "<DropzoneRoot {...props} maxFiles={1}>") {
		t.Error("AvatarDropzone does not pin maxFiles to 1")
	}
	if explicitAnyRe.MatchString(src) {
		t.Error("dropzone.tsx still types its drop target props as any")
	}
}

func TestFileFieldsUseTheNamedDropzones(t *testing.T) {
	for name, tc := range map[string]struct {
		src, want string
	}{
		"image":  {adminImageField(), "<AvatarDropzone"},
		"images": {adminImagesField(), "<BoxDropzone"},
		"video":  {adminVideoField(), "<CompactDropzone"},
		"videos": {adminVideosField(), "<BoxDropzone"},
		// The file fields read the look from the resource definition, so they
		// keep the variant prop.
		"file":  {adminFileField(), `variant={field.dropzone ?? "default"}`},
		"files": {adminFilesField(), `variant={field.dropzone ?? "default"}`},
	} {
		if !strings.Contains(tc.src, tc.want) {
			t.Errorf("%s field: want %q", name, tc.want)
		}
		if name != "file" && name != "files" && strings.Contains(tc.src, "variant=") {
			t.Errorf("%s field still picks its look with a variant prop", name)
		}
	}
}

func TestNextAdminShipsADropzoneTest(t *testing.T) {
	root := t.TempDir()
	if err := writeFrontendTestFiles(root, tripleOptions()); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(root, "apps", "admin", "__tests__", "dropzone.test.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"AvatarDropzone", "InlineDropzone", "BoxDropzone", "CompactDropzone", "MinimalDropzone", "Dropzone.Root"} {
		if !strings.Contains(string(body), name) {
			t.Errorf("the dropzone test does not cover %s", name)
		}
	}
}
