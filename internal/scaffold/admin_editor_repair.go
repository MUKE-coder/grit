package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M42, M43 and M44, and Tiptap 3, in the files upgrade does
// not deliver whole.
//
// The editor (components/forms/word-editor.tsx), the rich text field, the new
// lib/tiptap-extensions.ts and the page header arrive whole with the admin
// panel. Left here: the blog's pages and resource definition, which belong to
// the demo resource and so to the developer; the package.json that has to move
// to Tiptap 3 together with the editor; and the old page header, which nothing
// imports once resource-page.tsx is replaced.
//
// This runs after the admin panel is written, not with the API repairs: the
// package.json it bumps is one the panel's own writer may just have replaced.

// ─── M42: the blog runs on its resource definition ────────────────────────────

// blogsResourceFormView gives the Blog resource a form page of its own.
const blogsResourceFormView = `  // A blog is long-form writing, so its form is a page of its own rather than
  // a sheet, and the editor has the width to work in.
  formView: "page",
`

// blogsResourceLabelAnchor is where the form view goes in resources/blogs/blogs.ts.
const blogsResourceLabelAnchor = `  label: { singular: "Blog", plural: "Blogs" },
`

// Text only the hand-written blog pages had, which tells them apart from the
// <ResourcePage> wrappers and from a page somebody wrote from scratch.
const (
	legacyBlogsListMarker  = `"/api/admin/blogs?page_size=100"`
	legacyBlogDetailMarker = `<WordEditor`
)

// repairBlogsResourceSource adds the form view to the blog's definition.
func repairBlogsResourceSource(src string) (string, []string, []string) {
	if strings.Contains(src, "formView:") {
		return src, nil, nil
	}
	if strings.Count(src, blogsResourceLabelAnchor) != 1 {
		// Somebody renamed or rewrote the resource: its form view is theirs.
		return src, nil, nil
	}
	out := strings.Replace(src, blogsResourceLabelAnchor, blogsResourceLabelAnchor+blogsResourceFormView, 1)
	return out, []string{"blog posts are written on a form page of their own"}, nil
}

// blogPageRepair is one of the two hand-written blog pages and what replaces it.
type blogPageRepair struct {
	path    string
	marker  string
	content string
	fixed   string
	advice  string
}

func blogPageRepairs(shape adminPanelShape) []blogPageRepair {
	dir := filepath.Join(shape.routes, "(dashboard)", "resources", "blogs")
	return []blogPageRepair{
		{
			path:    filepath.Join(dir, "page.tsx"),
			marker:  legacyBlogsListMarker,
			content: adminBlogsPage(),
			fixed:   "the blog list is <ResourcePage>: paginated, with the columns, filters, bulk actions and export resources/blogs/blogs.ts declares",
			advice:  "this blog list is hand-written and fetches only the first 100 posts: replace it with <ResourcePage resource={blogsResource} />, as resources/users does",
		},
		{
			path:    filepath.Join(dir, "[id]", "page.tsx"),
			marker:  legacyBlogDetailMarker,
			content: adminResourceDetailRoute("blogs", "blogs", "Blogs"),
			fixed:   "the blog detail page is <ResourceDetailPage>, and posts are edited in the resource's own form",
			advice:  "this blog page is hand-written: replace it with <ResourceDetailPage resource={blogsResource} id={id} />, as resources/users/[id] does; the rich text field is the same Word-style editor",
		},
	}
}

// repairBlogPages replaces the blog's hand-written pages with the resource
// page, when Grit wrote them and nobody has edited them since. An edited one is
// named, with what to change, and left alone.
func repairBlogPages(root string, shape adminPanelShape, pristine func(string) bool) error {
	for _, page := range blogPageRepairs(shape) {
		data, err := os.ReadFile(page.path)
		if err != nil {
			continue
		}
		if !strings.Contains(string(data), page.marker) {
			continue
		}
		shown := page.path
		if rel, err := filepath.Rel(root, page.path); err == nil {
			shown = filepath.ToSlash(rel)
		}
		if !pristine(page.path) {
			fmt.Printf("  ⚠ %s: %s\n", shown, page.advice)
			continue
		}
		if err := writeFile(page.path, page.content); err != nil {
			return err
		}
		fmt.Printf("  ✓ %s: %s\n", shown, page.fixed)
	}
	return nil
}

// ─── M43: one editor on one schema ────────────────────────────────────────────

// tiptapExtensionsImport is how the Tiptap 3 editor names the shared schema.
// An editor without it is still on Tiptap 2 and its own extension list.
const tiptapExtensionsImport = `/lib/tiptap-extensions"`

// panelEditorFiles are the editor files a panel rooted at code has, where they exist.
func panelEditorFiles(code string) []string {
	var files []string
	for _, rel := range []string{"components/forms/word-editor.tsx", "components/forms/fields/rich-text-field.tsx"} {
		path := filepath.Join(code, filepath.FromSlash(rel))
		if fileExists(path) {
			files = append(files, path)
		}
	}
	return files
}

// panelCodeRoots are the directories an admin panel's components can live in,
// for the app whose package.json is in app.
func panelCodeRoots(app string) []string {
	return []string{
		app,
		filepath.Join(app, "src"),
		filepath.Join(app, "admin-panel"),
		filepath.Join(app, "src", "admin-panel"),
	}
}

// ─── Tiptap 3 ─────────────────────────────────────────────────────────────────

var tiptapDependencyLine = regexp.MustCompile(`"(@tiptap/[a-z0-9-]+)":\s*"([^"]*)"`)

// tiptapBelowPin reports whether a package.json version is older than
// tiptapVersion. A range is read by its floor; anything unparseable (a tag, a
// URL, a workspace link) is somebody's deliberate choice and is left alone.
func tiptapBelowPin(version string) bool {
	parse := func(v string) ([3]int, bool) {
		var out [3]int
		v = strings.TrimLeft(strings.TrimSpace(v), "^~=v")
		parts := strings.SplitN(v, ".", 3)
		if len(parts) != 3 {
			return out, false
		}
		for i, p := range parts {
			if i == 2 {
				if cut := strings.IndexAny(p, "-+"); cut >= 0 {
					p = p[:cut]
				}
			}
			n, err := strconv.Atoi(p)
			if err != nil {
				return out, false
			}
			out[i] = n
		}
		return out, true
	}
	have, ok := parse(version)
	if !ok {
		return false
	}
	want, _ := parse(tiptapVersion)
	for i := range have {
		if have[i] != want[i] {
			return have[i] < want[i]
		}
	}
	return false
}

// repairTiptapPackageJSON moves every @tiptap dependency below the pin to it,
// and adds the packages the Tiptap 3 editor imports that the file lacks.
//
// Old packages the editor no longer imports (extension-link, extension-color,
// the table row, cell and header) are raised rather than removed: they all
// exist in Tiptap 3, and code of the project's own may import them. What
// cannot stay is a Tiptap 2 package next to a Tiptap 3 one, since each wants
// its siblings at its own exact version.
func repairTiptapPackageJSON(src string) (string, []string, []string) {
	if !strings.Contains(src, `"@tiptap/`) {
		return src, nil, nil
	}
	raised := 0
	out := tiptapDependencyLine.ReplaceAllStringFunc(src, func(line string) string {
		m := tiptapDependencyLine.FindStringSubmatch(line)
		if !tiptapBelowPin(m[2]) {
			return line
		}
		raised++
		return fmt.Sprintf("%q: %q", m[1], tiptapVersion)
	})

	var missing []string
	for _, pkg := range tiptapPackages {
		if !strings.Contains(out, `"`+pkg+`"`) {
			missing = append(missing, pkg)
		}
	}
	if len(missing) > 0 {
		// After the last @tiptap line, so the block stays together and the
		// comma rules of JSON hold: that line is never the last in its object
		// unless it has no comma, and then the new lines take one each.
		locs := tiptapDependencyLine.FindAllStringIndex(out, -1)
		last := locs[len(locs)-1]
		lineStart := strings.LastIndex(out[:last[0]], "\n") + 1
		indent := out[lineStart:last[0]]
		cut := last[1]
		hasComma := strings.HasPrefix(out[cut:], ",")
		var b strings.Builder
		if hasComma {
			cut++
			for _, pkg := range missing {
				fmt.Fprintf(&b, "\n%s%q: %q,", indent, pkg, tiptapVersion)
			}
		} else {
			for _, pkg := range missing {
				fmt.Fprintf(&b, ",\n%s%q: %q", indent, pkg, tiptapVersion)
			}
		}
		out = out[:cut] + b.String() + out[cut:]
	}
	if out == src {
		return src, nil, nil
	}
	return out, []string{fmt.Sprintf("Tiptap is %s (%d package(s) raised, %d added), clearing GHSA-cp6q-959q-f8rh", tiptapVersion, raised, len(missing))}, nil
}

// repairTiptap raises Tiptap in every app whose editor is on Tiptap 3. An app
// whose editor is still the Tiptap 2 one, because it was edited, keeps its
// versions and is told what to change: raising them under that code would stop
// the panel compiling.
func repairTiptap(root string, m *manifest.Manifest) error {
	apps := []string{
		root,
		filepath.Join(root, "frontend"),
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "apps", "web"),
	}
	for _, app := range apps {
		pkg := filepath.Join(app, "package.json")
		if !fileContains(pkg, `"@tiptap/`) {
			continue
		}
		var editors []string
		for _, code := range panelCodeRoots(app) {
			editors = append(editors, panelEditorFiles(code)...)
		}
		if len(editors) == 0 {
			continue
		}
		var stale []string
		for _, editor := range editors {
			// Tiptap imported, and not through the shared list: the Tiptap 2 editor.
			if fileContains(editor, `from "@tiptap/`) && !fileContains(editor, tiptapExtensionsImport) {
				stale = append(stale, editor)
			}
		}
		if len(stale) > 0 {
			for _, editor := range stale {
				shown := editor
				if rel, err := filepath.Rel(root, editor); err == nil {
					shown = filepath.ToSlash(rel)
				}
				fmt.Printf("  ⚠ %s: still builds its own Tiptap 2 extension list, so Tiptap stays on 2.x in this app (GHSA-cp6q-959q-f8rh): build the editor from richTextExtensions() in lib/tiptap-extensions.ts, or take Grit's copy with grit upgrade --force, then raise every @tiptap package to %s\n", shown, tiptapVersion)
			}
			continue
		}
		if err := repairTextFile(root, m, pkg, repairTiptapPackageJSON); err != nil {
			return err
		}
	}
	return nil
}

// ─── M44: one page header ─────────────────────────────────────────────────────

const (
	legacyPageHeaderRel    = "components/layout/page-header.tsx"
	legacyPageHeaderImport = `/components/layout/page-header"`
)

// pruneLegacyPageHeader deletes components/layout/page-header.tsx once nothing
// imports it, if Grit wrote it and it has not been edited. resource-page.tsx
// was its only user, and now renders components/chrome/PageHeader.tsx.
func pruneLegacyPageHeader(root string, shape adminPanelShape, pristine func(string) bool) {
	path := filepath.Join(shape.code, filepath.FromSlash(legacyPageHeaderRel))
	if !fileExists(path) || !pristine(path) {
		return
	}
	importer := ""
	for _, dir := range []string{shape.code, shape.routes} {
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || importer != "" {
				return nil
			}
			if d.IsDir() {
				if name := d.Name(); name == "node_modules" || name == ".next" {
					return filepath.SkipDir
				}
				return nil
			}
			if ext := filepath.Ext(p); ext != ".ts" && ext != ".tsx" {
				return nil
			}
			if fileContains(p, legacyPageHeaderImport) {
				importer = p
			}
			return nil
		})
	}
	if importer != "" {
		return
	}
	if err := os.Remove(path); err != nil {
		return
	}
	manifest.Drop(path)
	shown := path
	if rel, err := filepath.Rel(root, path); err == nil {
		shown = filepath.ToSlash(rel)
	}
	fmt.Printf("  ✓ %s: removed, since components/chrome/PageHeader.tsx is the one page header\n", shown)
}

// repairAdminScreens applies M41 to M44 and Tiptap 3 to a project's panels.
func repairAdminScreens(root string, opts Options) error {
	// M41 first: the pages just written import these.
	if err := repairSharedTypes(root); err != nil {
		return err
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	// The recording, not the manifest on disk: the panel was written earlier
	// in this same upgrade, and the file on disk is not saved until it ends.
	pristine := manifest.IsUnchanged
	for _, shape := range adminPanelShapes(root) {
		blogs := filepath.Join(shape.code, "resources", "blogs", "blogs.ts")
		if fileExists(blogs) {
			if err := repairTextFile(root, m, blogs, repairBlogsResourceSource); err != nil {
				return err
			}
		}
		if err := repairBlogPages(root, shape, pristine); err != nil {
			return err
		}
		pruneLegacyPageHeader(root, shape, pristine)
	}
	return repairTiptap(root, m)
}
