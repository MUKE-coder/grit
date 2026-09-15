package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M24 and M25, in the files upgrade does not deliver whole.
//
// lib/query-client.ts, components/shared/providers.tsx, the form builder and the
// resource detail page are framework files and arrive whole. Two cases are left:
// the blog editor page, which is the demo resource's and so the developer's, and
// a query-client.ts or providers.tsx somebody edited, which the manifest guard
// keeps. The second matters because the two files change together: a new
// providers.tsx next to an old query-client.ts does not compile.

// ─── M25: the blog editor loads Tiptap on demand ─────────────────────────────

const (
	blogEditorOldReactImport = `import { useEffect, useRef, useState } from "react";`
	blogEditorNewReactImport = `import { Suspense, useEffect, useRef, useState } from "react";
import dynamic from "next/dynamic";`
	blogEditorOldStaticImport = `import { WordEditor } from "@/components/forms/word-editor";
`
	blogEditorAnchor = "\ninterface Blog {"

	// blogEditorDynamic sits between the imports and interface Blog.
	blogEditorDynamic = `// Tiptap and its extensions load after the page renders, behind a placeholder
// the editor's height, so the title, cover and excerpt do not wait for them.
// Suspense covers the Vite admin, where dynamic() is React.lazy.
const WordEditor = dynamic(
  () => import("@/components/forms/word-editor").then((m) => m.WordEditor),
  { ssr: false, loading: EditorPlaceholder },
);

function EditorPlaceholder() {
  return <div className="min-h-[500px] w-full animate-pulse rounded-lg border border-border bg-bg-hover/40" />;
}

`

	blogEditorOldUsage = `        <WordEditor
          value={content}
          onChange={setContent}
          placeholder="Start your article here. Use the toolbar to format headings, lists, tables, images, and more."
          minHeight={500}
          onBlur={() => { if (content !== blog.content) save.mutate({ content }); }}
        />`
	blogEditorNewUsage = `        <Suspense fallback={<EditorPlaceholder />}>
          <WordEditor
            value={content}
            onChange={setContent}
            placeholder="Start your article here. Use the toolbar to format headings, lists, tables, images, and more."
            minHeight={500}
            onBlur={() => { if (content !== blog.content) save.mutate({ content }); }}
          />
        </Suspense>`
)

// ─── M24: the admin's QueryClient is created per mount ───────────────────────

// adminQueryClientOldSource is lib/query-client.ts before M24.
const adminQueryClientOldSource = `import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
});
`

const adminQueryClientOldExport = "export const queryClient = new QueryClient({"

// adminPanelShape is one admin panel in a project: where its components and lib
// live, where its routes live, and the import alias its files use for the first.
type adminPanelShape struct {
	code, routes, alias string
}

// adminPanelShapes finds the Next.js admin panels a project has. The Vite admin
// is not upgraded, so it is not repaired either.
func adminPanelShapes(root string) []adminPanelShape {
	var shapes []adminPanelShape
	admin := filepath.Join(root, "apps", "admin")
	for _, cfg := range []string{"next.config.ts", "next.config.js", "next.config.mjs"} {
		if fileExists(filepath.Join(admin, cfg)) {
			shapes = append(shapes, adminPanelShape{code: admin, routes: filepath.Join(admin, "app"), alias: "@"})
			break
		}
	}
	web := filepath.Join(root, "apps", "web")
	if fileExists(filepath.Join(web, "admin-panel", "components", "shared", "providers.tsx")) {
		shapes = append(shapes, adminPanelShape{
			code:   filepath.Join(web, "admin-panel"),
			routes: filepath.Join(web, "app", "admin"),
			alias:  "@admin",
		})
	}
	return shapes
}

// withAlias points the admin's @/ imports in text at the panel's alias.
func (s adminPanelShape) withAlias(text string) string {
	if s.alias == "@" {
		return text
	}
	return strings.ReplaceAll(text, `"@/`, `"`+s.alias+`/`)
}

// repairAdminBundles brings a project's admin panel up to M24 and M25.
func repairAdminBundles(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, shape := range adminPanelShapes(root) {
		queryClient := filepath.Join(shape.code, "lib", "query-client.ts")
		if fileExists(queryClient) {
			if err := repairTextFile(root, m, queryClient, repairAdminQueryClientSource); err != nil {
				return err
			}
			warnModuleQueryClientImports(root, shape)
		}
		providers := filepath.Join(shape.code, "components", "shared", "providers.tsx")
		if fileExists(providers) {
			if err := repairTextFile(root, m, providers, func(src string) (string, []string, []string) {
				return repairAdminProvidersSource(src, shape)
			}); err != nil {
				return err
			}
		}
		blog := filepath.Join(shape.routes, "(dashboard)", "resources", "blogs", "[id]", "page.tsx")
		if fileExists(blog) {
			if err := repairTextFile(root, m, blog, func(src string) (string, []string, []string) {
				out, fixed, warn := repairBlogEditorSource(src, shape)
				out, routesFixed := repairBlogEditorRoutesSource(out)
				return out, append(fixed, routesFixed...), warn
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// repairAdminQueryClientSource turns the module-level client into a factory.
//
// The file Grit wrote becomes the template. An edited one keeps its options: the
// export becomes a function returning the same client, so staleTime or retry
// settings somebody tuned survive.
func repairAdminQueryClientSource(src string) (string, []string, []string) {
	if strings.Contains(src, "function makeQueryClient(") {
		return src, nil, nil
	}
	if src == adminQueryClientOldSource {
		return adminQueryClientSource, []string{"the admin creates its React Query client per mount, not once per server process"}, nil
	}
	body := strings.TrimRight(src, " \t\n")
	if strings.Count(src, adminQueryClientOldExport) != 1 || !strings.HasSuffix(body, "});") {
		return src, nil, []string{"lib/query-client.ts is not the file Grit wrote: export a makeQueryClient() that returns a new QueryClient, and create it once per mount in components/shared/providers.tsx, rather than a client at module scope"}
	}
	start := strings.Index(body, adminQueryClientOldExport)
	options := body[start+len(adminQueryClientOldExport) : len(body)-len("});")]
	out := body[:start] + `// A new client, for <Providers> to create once per mount. Never a client at
// module scope: on the server that is one cache for every request rendered.
export function makeQueryClient(): QueryClient {
  return new QueryClient({` + strings.ReplaceAll(options, "\n", "\n  ") + `});
}
`
	return out, []string{"the admin creates its React Query client per mount, with the options this file already had"}, nil
}

// repairAdminProvidersSource creates the client in the provider, per mount.
func repairAdminProvidersSource(src string, shape adminPanelShape) (string, []string, []string) {
	if strings.Contains(src, "makeQueryClient") {
		return src, nil, nil
	}
	oldImports := shape.withAlias(adminProvidersOldImports)
	if !strings.Contains(src, shape.withAlias(`import { queryClient } from "@/lib/query-client";`)) {
		return src, nil, nil
	}
	if strings.Count(src, oldImports) != 1 || strings.Count(src, adminProvidersOldBody) != 1 {
		return src, nil, []string{"components/shared/providers.tsx is not the file Grit wrote: replace the imported queryClient with const [queryClient] = useState(makeQueryClient), since lib/query-client.ts no longer exports a client"}
	}
	out := strings.Replace(src, oldImports, shape.withAlias(adminProvidersNewImports), 1)
	out = strings.Replace(out, adminProvidersOldBody, adminProvidersNewBody, 1)
	return out, []string{"the admin's provider creates its React Query client per mount"}, nil
}

// repairBlogEditorSource loads the Word-style editor on demand.
func repairBlogEditorSource(src string, shape adminPanelShape) (string, []string, []string) {
	staticImport := shape.withAlias(blogEditorOldStaticImport)
	if !strings.Contains(src, staticImport) {
		// Already lazy, or a page that no longer uses the editor.
		return src, nil, nil
	}
	if strings.Count(src, blogEditorOldReactImport) != 1 || strings.Count(src, blogEditorAnchor) != 1 ||
		strings.Count(src, blogEditorOldUsage) != 1 {
		return src, nil, []string{"the blog editor page is not the file Grit wrote: load WordEditor with dynamic(() => import(...), { ssr: false }) inside <Suspense>, so Tiptap is not in the page's first load"}
	}
	out := strings.Replace(src, staticImport, "", 1)
	out = strings.Replace(out, blogEditorOldReactImport, blogEditorNewReactImport, 1)
	out = strings.Replace(out, blogEditorAnchor, "\n"+shape.withAlias(blogEditorDynamic)+"interface Blog {", 1)
	out = strings.Replace(out, blogEditorOldUsage, blogEditorNewUsage, 1)
	return out, []string{"the blog editor loads Tiptap after the page renders"}, nil
}

// blogEditorOldRoute is where the blog editor saved, toggled and deleted a post.
// The API serves those under /api/admin/blogs; /api/blogs is the public,
// read-only list, so every save returned 404.
const blogEditorOldRoute = `"/api/blogs/" + params.id`

func repairBlogEditorRoutesSource(src string) (string, []string) {
	if !strings.Contains(src, blogEditorOldRoute) {
		return src, nil
	}
	return strings.ReplaceAll(src, blogEditorOldRoute, `"/api/admin/blogs/" + params.id`),
		[]string{"the blog editor saves, publishes and deletes through /api/admin/blogs, where the API serves them"}
}

// warnModuleQueryClientImports names any file still importing the module-level
// client, which lib/query-client.ts no longer exports.
func warnModuleQueryClientImports(root string, shape adminPanelShape) {
	needle := `from "` + shape.alias + `/lib/query-client"`
	for _, dir := range []string{shape.code, shape.routes} {
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if name := d.Name(); name == "node_modules" || name == ".next" {
					return filepath.SkipDir
				}
				return nil
			}
			ext := filepath.Ext(path)
			if (ext != ".ts" && ext != ".tsx") || strings.HasSuffix(filepath.ToSlash(path), "components/shared/providers.tsx") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil || !strings.Contains(string(data), needle) || !strings.Contains(string(data), "queryClient }") {
				return nil
			}
			shown := path
			if rel, err := filepath.Rel(root, path); err == nil {
				shown = filepath.ToSlash(rel)
			}
			fmt.Printf("  ⚠ %s: imports the module-level queryClient, which lib/query-client.ts no longer exports; call useQueryClient() inside the component instead\n", shown)
			return nil
		})
	}
}
