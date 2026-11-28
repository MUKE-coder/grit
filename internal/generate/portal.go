package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The customer's half of a resource.
//
// A generated resource has always produced a complete screen for the operator
// and a React Query hook for the customer. That asymmetry is the whole of what
// an indie SaaS complains about: the admin is declarative and generated, the
// product is hand-written, and the product is the part people pay for.
//
// It is also the odd one out inside Grit rather than a missing idea. The Expo
// generator writes a list screen, a detail screen, a form and a create and an
// edit screen, and injects a link into the app's navigation. The Wails client
// does the same. Only the browser, where almost every Grit app actually sells
// something, got the hook and nothing else.
//
// So this is that pattern again, for apps/web/app/(app): the signed-in customer
// area that `grit add web-auth` writes.

// PortalLevel is how much of the customer's half the generator writes.
type PortalLevel string

const (
	// PortalScreens is pages, a form and a sidebar entry. The default in a SaaS
	// project, because somebody building one in two days wants a screen.
	PortalScreens PortalLevel = "screens"
	// PortalHeadless is the data and permission layer with no JSX: the hooks,
	// a descriptor of the resource and its form schema. For the screen you are
	// going to write yourself, which for the product's own surface is most of
	// them. Google Keep's note grid is nobody's generated table.
	PortalHeadless PortalLevel = "headless"
	// PortalOff is the hook alone, which is what every version before v3.584.0
	// wrote.
	PortalOff PortalLevel = "off"
)

// ValidPortalLevel turns the flag into a level, or says what the answers are.
func ValidPortalLevel(s string) (PortalLevel, error) {
	switch PortalLevel(strings.ToLower(strings.TrimSpace(s))) {
	case PortalScreens:
		return PortalScreens, nil
	case PortalHeadless:
		return PortalHeadless, nil
	case PortalOff:
		return PortalOff, nil
	case "":
		return "", nil
	}
	return "", fmt.Errorf("unknown --portal %q: use screens, headless or off", s)
}

// portalLevel is what this run writes for the customer, and why.
//
// The flag wins when it was given, because an explicit answer should not be
// second-guessed by a file. Otherwise the project decides: a SaaS writes
// screens for a resource that belongs to somebody, and nothing otherwise.
//
// "Belongs to somebody" is --owned-by or --tenant-owned, and the distinction
// matters: with organizations the customer's rows are the organization's, not
// the person's, and a portal screen that listed only your own rows inside a
// shared workspace would be wrong in the direction nobody checks.
func (g *Generator) portalLevel() PortalLevel {
	if g.Portal != "" {
		return g.Portal
	}
	if !g.projectIsSaaS() {
		return PortalOff
	}
	if g.Definition.OwnerField() == nil && !g.Definition.TenantOwned {
		return PortalOff
	}
	return PortalScreens
}

// portalRoot is apps/web/app/(app), or empty when there is nowhere to write.
//
// The directory has to exist already. It is written by grit add web-auth, with
// the layout, the guards and the sign-in pages that make a customer area more
// than a folder, and creating a page under a group that has no layout would
// put an unprotected screen on the internet.
func (g *Generator) portalRoot() string {
	if g.Architecture != "triple" && g.Architecture != "double" && g.Architecture != "full" {
		return ""
	}
	if g.Frontend == "tanstack" {
		// The Vite web app keeps routes under src/routes and has no (app)
		// group. Worth doing, and not by pretending the Next shape fits it.
		return ""
	}
	dir := filepath.Join(g.Root, "apps", "web", "app", "(app)")
	if !dirExists(dir) {
		return ""
	}
	return dir
}

// projectIsSaaS reads the one bit grit new --saas writes.
func (g *Generator) projectIsSaaS() bool {
	return readGritJSONBool(g.Root, "saas")
}

// readGritJSONBool reads a true/false field out of grit.json.
//
// Separate from extractJSONValue, which finds the next quoted string after the
// key and would answer a neighbouring key's name for a boolean. Same reason
// there is no encoding/json here: the file is three lines and the package
// reads it on every command.
func readGritJSONBool(root, key string) bool {
	data, err := os.ReadFile(filepath.Join(root, "grit.json"))
	if err != nil {
		return false
	}
	content := string(data)
	idx := strings.Index(content, `"`+key+`"`)
	if idx < 0 {
		return false
	}
	rest := content[idx+len(key)+2:]
	colon := strings.Index(rest, ":")
	if colon < 0 {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(rest[colon+1:]), "true")
}

// writePortalFiles writes the customer's half, at whatever level applies.
func (g *Generator) writePortalFiles(names Names) error {
	level := g.portalLevel()
	if level == PortalOff {
		return nil
	}
	root := g.portalRoot()
	if root == "" {
		// Asked for explicitly and impossible: say so rather than succeed
		// quietly, because the next thing that happens is somebody looking for
		// a page that was never written.
		if g.Portal != "" {
			fmt.Println("  ⚠ --portal needs the signed-in customer area, which this project does not have.")
			fmt.Println("    Run grit add web-auth first. It writes apps/web/app/(app) and the sign-in pages.")
		}
		return nil
	}

	if err := g.ensurePortalKit(); err != nil {
		return err
	}
	if err := g.writePortalDescriptor(names); err != nil {
		return err
	}
	fmt.Printf("  ✓ apps/web/lib/portal/%s.ts\n", names.PluralKebab)

	if level == PortalHeadless {
		// Deliberately no page and no nav entry. The descriptor, the schema
		// and the hooks are the contract; the screen is yours, and a generated
		// one you have to delete first is worse than none.
		return nil
	}

	dir := filepath.Join(root, names.PluralKebab)
	if err := g.writePortalForm(names); err != nil {
		return err
	}
	fmt.Printf("  ✓ apps/web/components/portal/%s-form.tsx\n", names.Kebab)

	for _, page := range []struct {
		path string
		body string
	}{
		{filepath.Join(dir, "page.tsx"), g.portalListPage(names)},
		{filepath.Join(dir, "new", "page.tsx"), g.portalNewPage(names)},
		{filepath.Join(dir, "[id]", "page.tsx"), g.portalDetailPage(names)},
	} {
		// Never overwritten. The list screen is the first file anybody edits,
		// and a second `grit generate resource` to add a field must not throw
		// the edit away. Adding a field reaches the form through the
		// descriptor, which is generated and says so at the top.
		if fileExists(page.path) {
			continue
		}
		if err := writeFileWithDirs(page.path, page.body); err != nil {
			return err
		}
	}
	fmt.Printf("  ✓ apps/web/app/(app)/%s/ (list, new, detail)\n", names.PluralKebab)

	if err := g.injectPortalNav(names); err != nil {
		return err
	}
	return nil
}

// injectPortalNav adds the screen to the customer sidebar.
//
// Both halves or neither: an entry naming an icon the file does not import is
// a layout that does not compile, which would take the whole customer area
// down rather than one screen.
func (g *Generator) injectPortalNav(names Names) error {
	layout := filepath.Join(g.Root, "apps", "web", "app", "(app)", "layout.tsx")
	data, err := os.ReadFile(layout)
	if err != nil {
		return nil
	}
	content := string(data)
	href := "/" + names.PluralKebab
	if strings.Contains(content, `href: "`+href+`"`) {
		return nil // already there, so a repeat run changes nothing
	}
	if !strings.Contains(content, "// grit:nav") {
		fmt.Printf("  ⚠ apps/web/app/(app)/layout.tsx has no // grit:nav marker, so %s is not in the sidebar.\n", portalLabel(names))
		fmt.Println("    Add it to the NAV array by hand, or run grit upgrade to get the marker.")
		return nil
	}

	icon := guessLucideIcon(names.Plural)
	if !strings.Contains(content, "// grit:nav-icons") {
		// An older layout with the entry marker and not the icon one. Rather
		// than fail, use an icon the file is known to import already.
		icon = "LayoutDashboard"
	} else if !strings.Contains(content, "import { "+icon+" }") {
		if err := injectBefore(layout, "// grit:nav-icons",
			fmt.Sprintf("import { %s } from \"lucide-react\";", icon)); err != nil {
			return err
		}
	}
	entry := fmt.Sprintf("  { href: %q, label: %q, icon: %s },", href, portalLabel(names), icon)
	if err := injectBefore(layout, "// grit:nav", entry); err != nil {
		return err
	}
	fmt.Printf("  ✓ %s in the customer sidebar\n", portalLabel(names))
	return nil
}

// refreshPortalDescriptor rewrites lib/portal/<resource>.ts after a column
// was added to the model.
//
// Only for a resource that already has one. A resource generated with
// --portal off does not acquire a descriptor by having a field added: that
// would be a new file nobody asked for, in a project that may have no web app
// at all.
//
// One row injected, not the file rewritten.
//
// Rewriting it looked simpler and was wrong. The definition a rebuild works
// from is read back out of the Go model, and that round trip is lossy: a
// text column comes back as a string, a required column comes back optional,
// and tenant ownership comes back as nothing at all. A rewrite therefore
// changed three fields that nobody had touched and flipped the scope of the
// resource from organization to user, which is a lie about who the rows
// belong to. The TypeScript type and the Zod schemas beside it are injected
// into for the same reason.
func (g *Generator) refreshPortalDescriptor(names Names, added Field) error {
	path := filepath.Join(g.Root, "apps", "web", "lib", "portal", names.PluralKebab+".ts")
	if !fileExists(path) {
		return nil
	}

	row := portalFieldRow(portalFormField{
		Name:     toPascalCase(added.Name),
		JSON:     toSnakeCase(added.Name),
		Label:    humanLabel(added.Name),
		Kind:     portalKind(added),
		Required: added.Required,
		Options:  portalOptions(added),
		Relation: portalRelation(added),
	})

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	if strings.Contains(content, fmt.Sprintf("name: %q", toSnakeCase(added.Name))) {
		return nil // already there, so a repeat run changes nothing
	}
	const closer = "  ] satisfies PortalField[],"
	if !strings.Contains(content, closer) {
		return fmt.Errorf("the field list in %s does not end the way it was written", relativeToRoot(g.Root, path))
	}
	content = strings.Replace(content, closer, row+"\n"+closer, 1)
	if err := writeFileWithDirs(path, content); err != nil {
		return err
	}
	fmt.Printf("  + apps/web/lib/portal/%s.ts\n", names.PluralKebab)
	return nil
}

// removePortalFiles takes the customer's half away with the resource.
func (g *Generator) removePortalFiles(names Names) []string {
	var removed []string
	web := filepath.Join(g.Root, "apps", "web")
	for _, path := range []string{
		filepath.Join(web, "app", "(app)", names.PluralKebab),
		filepath.Join(web, "components", "portal", names.Kebab+"-form.tsx"),
		filepath.Join(web, "lib", "portal", names.PluralKebab+".ts"),
	} {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		if err := os.RemoveAll(path); err == nil {
			removed = append(removed, relativeToRoot(g.Root, path))
		}
	}
	// The nav entry and its icon import, which are injections rather than
	// files: a removal that leaves them behind leaves a sidebar link to a 404.
	layout := filepath.Join(web, "app", "(app)", "layout.tsx")
	data, err := os.ReadFile(layout)
	if err != nil {
		return removed
	}
	href := `{ href: "/` + names.PluralKebab + `"`
	if !strings.Contains(string(data), href) {
		return removed
	}
	if err := removeLinesContaining(layout, href); err != nil {
		return removed
	}
	// The icon import only when nothing else uses it: two resources can guess
	// the same icon, and dropping it would break the other one's entry.
	icon := guessLucideIcon(names.Plural)
	if after, err := os.ReadFile(layout); err == nil && !strings.Contains(string(after), "icon: "+icon) {
		_ = removeLinesContaining(layout, "import { "+icon+" } from \"lucide-react\";")
	}
	return append(removed, relativeToRoot(g.Root, layout)+" (sidebar entry)")
}

// portalLabel is what the sidebar and the page heading call this resource:
// "Blog posts", not "BlogPosts" and not "blog_posts".
func portalLabel(names Names) string {
	return humanLabel(names.PluralPascal)
}
