package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The anchors this repair edits around.
//
// They are declared here, once, and a test asserts each one appears exactly
// once in the sidebar the scaffold writes. Without that test an anchor can
// drift out of the template and the repair quietly stops applying: it warns
// into a log nobody reads and returns, so every upgraded project silently
// keeps the old sidebar. That has happened before.
const (
	contentNavMarker     = "// grit:nav:content"
	systemNavClose       = "  // grit:nav:system\n] as const;"
	contentNavGroupsLoop = "  const groups: Record<string, typeof resources> = { _root: [] };\n" +
		"  for (const r of visibleResources) {\n" +
		"    const key = r.group || \"_root\";\n" +
		"    if (!groups[key]) groups[key] = [];\n" +
		"    groups[key].push(r);\n" +
		"  }"
	contentNavHeadingOld = "{!collapsed && groups._root.length > 0 && ("
	contentNavHeadingNew = "{!collapsed && (groups._root.length > 0 || visibleContentNav.length > 0) && ("
	contentNavRootMapEnd = "                active={pathname.startsWith(\"/resources/\" + r.slug)}\n" +
		"                collapsed={collapsed}\n" +
		"                onClick={onMobileClose}\n" +
		"              />\n" +
		"            );\n" +
		"          })}\n"
)

// contentNavDeclaration is the array a plugin injects a working surface into.
const contentNavDeclaration = `

// Plugins inject a WORKING surface here, as opposed to an operational one.
//
// SYSTEM_NAV above feeds the System Hub's Extensions tab, which is three
// clicks from the dashboard. That is the right home for a webhook log or an
// impersonation screen: things you open while configuring. It is the wrong
// home for anything you open while working, and until this existed a plugin
// had nowhere else to put one.
//
// Entries here render in the rail's Content group, beside the resources,
// because that is what they are: a place the work lives. The marker MUST
// stay inside this array literal, or "grit plugin add" can't wire its link.
export const CONTENT_NAV: readonly NavEntry[] = [
  // grit:nav:content
] as const;`

const contentNavFilter = `

  // Plugin-injected working surfaces, filtered the same way resources are:
  // an entry marked adminOnly is for admins, and one that names a permission
  // waits for the permission check rather than flickering into view and out
  // again while it loads.
  const visibleContentNav = CONTENT_NAV.filter(
    (entry) =>
      (!entry.adminOnly || isAdmin) && (permsLoading || !entry.requires || can(entry.requires))
  );`

const contentNavRender = `
          {/* Working surfaces added by plugins: a board, a calendar, an inbox.
              They sit with the resources because that is what they are to the
              person using the app. */}
          {visibleContentNav.map((entry) => {
            const Icon = getIcon(entry.iconKey);
            return (
              <SidebarLink
                key={entry.href}
                href={entry.href}
                icon={<Icon className="h-5 w-5" />}
                label={t("nav." + entry.href.replace(/^\//, "").replace(/\//g, "."), entry.label)}
                active={pathname === entry.href || pathname.startsWith(entry.href + "/")}
                collapsed={collapsed}
                onClick={onMobileClose}
              />
            );
          })}
`

// repairContentNav gives an existing project the Content nav anchor.
//
// Before this, SYSTEM_NAV was the only place a plugin could put a link, and
// it renders in the System Hub's Extensions tab. A plugin that adds a board
// or a calendar had to file it under System, where nobody looks for a board.
func repairContentNav(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, shape := range adminPanelShapes(root) {
		sidebar := filepath.Join(shape.code, "components", "chrome", "CollapsibleSidebar.tsx")
		if !fileExists(sidebar) {
			continue
		}
		if err := repairTextFile(root, m, sidebar, repairContentNavSource); err != nil {
			return err
		}
	}
	return nil
}

// repairContentNavSource is the edit, separated from the file handling so it
// can be tested on a string.
func repairContentNavSource(src string) (string, []string, []string) {
	if strings.Contains(src, contentNavMarker) {
		return src, nil, nil // already has it
	}

	// Every anchor is checked before anything is written. A partial edit
	// here produces a sidebar that does not compile, which is worse than a
	// sidebar that is merely out of date.
	missing := []string{}
	if !strings.Contains(src, systemNavClose) {
		missing = append(missing, "the SYSTEM_NAV array")
	}
	if !strings.Contains(src, contentNavGroupsLoop) {
		missing = append(missing, "the resource grouping loop")
	}
	if !strings.Contains(src, contentNavHeadingOld) {
		missing = append(missing, "the Content heading")
	}
	if !strings.Contains(src, contentNavRootMapEnd) {
		missing = append(missing, "the resource link list")
	}
	if len(missing) > 0 {
		return src, nil, []string{
			"left alone, this sidebar has been edited and no longer has " +
				strings.Join(missing, ", ") +
				". To let plugins add a working surface to the rail, copy CONTENT_NAV " +
				"from a freshly scaffolded project.",
		}
	}

	out := strings.Replace(src, systemNavClose, systemNavClose+contentNavDeclaration, 1)
	out = strings.Replace(out, contentNavGroupsLoop, contentNavGroupsLoop+contentNavFilter, 1)
	out = strings.Replace(out, contentNavHeadingOld, contentNavHeadingNew, 1)
	out = strings.Replace(out, contentNavRootMapEnd, contentNavRootMapEnd+contentNavRender, 1)

	return out, []string{"plugins can now add a working surface to the Content group, not just System"}, nil
}
