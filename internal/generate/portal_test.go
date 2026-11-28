package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// saasProject is a project that said yes to customers, with the customer area
// grit add web-auth writes.
//
// The layout comes from the scaffold template rather than a stand-in. The
// generator injects into it at two markers, and a hand-written fixture would
// go on passing after somebody edited the template and took a marker out,
// which is the failure these tests exist to catch.
func saasProject(t *testing.T, module string, saas bool) string {
	t.Helper()
	root := setupMinimalProject(t, module)

	writeTestFile(t, filepath.Join(root, "grit.json"), `{
  "architecture": "triple",
  "frontend": "next",
  "saas": `+boolWord(saas)+`
}
`)
	writeTestFile(t,
		filepath.Join(root, "apps", "web", "app", "(app)", "layout.tsx"),
		scaffold.WebAccountLayout(scaffold.Options{ProjectName: "keep"}))
	// The hooks directory is what makes the generator write web hooks at all,
	// and the pages import them.
	writeTestFile(t, filepath.Join(root, "apps", "web", "hooks", "use-resource.ts"), "export {};\n")
	return root
}

func boolWord(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func ownedNoteGenerator(t *testing.T, root, module string) *Generator {
	t.Helper()
	def := mustFields(t, "Note", "title:string,body:text,user:belongs_to:User")
	def.OwnedBy = "user"
	g := newTestGenerator(root, module, def)
	g.Architecture = "triple"
	g.Frontend = "next"
	return g
}

// The case the whole thing exists for: a resource that belongs to a customer,
// in a project that has customers, gets the customer's screens.
func TestAnOwnedResourceInASaaSGetsItsCustomerScreens(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	g := ownedNoteGenerator(t, root, module)

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}

	web := filepath.Join(root, "apps", "web")
	for _, rel := range []string{
		"app/(app)/notes/page.tsx",
		"app/(app)/notes/new/page.tsx",
		"app/(app)/notes/[id]/page.tsx",
		"components/portal/note-form.tsx",
		"components/portal/portal-page.tsx",
		"components/portal/portal-field.tsx",
		"lib/portal/notes.ts",
		"lib/portal/types.ts",
	} {
		if _, err := os.Stat(filepath.Join(web, filepath.FromSlash(rel))); err != nil {
			t.Errorf("missing %s: the customer has a hook and no screen, which is the bug", rel)
		}
	}

	layout := readTestFile(t, filepath.Join(web, "app", "(app)", "layout.tsx"))
	if !strings.Contains(layout, `{ href: "/notes", label: "Notes"`) {
		t.Error("the screen is not in the customer sidebar, so nothing reaches it")
	}
	// Both halves of the injection or neither: an entry naming an icon the
	// file does not import takes the whole customer area down, not one page.
	icon := guessLucideIcon("notes")
	if !strings.Contains(layout, "import { "+icon+" }") {
		t.Errorf("the sidebar entry uses %s and the layout does not import it, so it will not compile", icon)
	}
}

// The default is unchanged. A project that never said it had customers gets
// exactly what every version before this one wrote.
func TestANonSaaSProjectGetsNoPortal(t *testing.T) {
	const module = "shop/apps/api"
	root := saasProject(t, module, false)
	g := ownedNoteGenerator(t, root, module)

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "app", "(app)", "notes")); err == nil {
		t.Error("a project that is not a SaaS got customer screens it did not ask for")
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "lib", "portal")); err == nil {
		t.Error("a project that is not a SaaS got a portal descriptor")
	}
}

// An explicit answer is not second-guessed by a file.
func TestTheFlagBeatsTheProject(t *testing.T) {
	const module = "shop/apps/api"
	root := saasProject(t, module, false)
	g := ownedNoteGenerator(t, root, module)
	g.Portal = PortalScreens

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "app", "(app)", "notes", "page.tsx")); err != nil {
		t.Error("--portal screens was ignored because grit.json says this is not a SaaS")
	}
}

// Headless is the level for the screen you are writing yourself: the data and
// the shape, and no JSX to delete first.
func TestHeadlessWritesTheDescriptorAndNoPages(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	g := ownedNoteGenerator(t, root, module)
	g.Portal = PortalHeadless

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "lib", "portal", "notes.ts")); err != nil {
		t.Error("headless wrote no descriptor, which is the only thing it was asked for")
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "app", "(app)", "notes")); err == nil {
		t.Error("headless wrote pages: the point is that the screen is yours")
	}
	layout := readTestFile(t, filepath.Join(root, "apps", "web", "app", "(app)", "layout.tsx"))
	if strings.Contains(layout, `href: "/notes"`) {
		t.Error("headless put a sidebar link on a page it did not write")
	}
}

// A resource nobody owns is not the customer's, whatever the project is.
func TestASharedResourceGetsNoCustomerScreen(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	g := newTestGenerator(root, module, mustFields(t, "Plan", "name:string,price:int"))
	g.Architecture, g.Frontend = "triple", "next"

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "app", "(app)", "plans")); err == nil {
		t.Error("a shared table got a customer screen, which would list everybody's rows")
	}
}

// The descriptor is what the headless level is for, so what it leaves out
// matters as much as what it carries.
func TestTheDescriptorLeavesOutWhatIsNotTheCustomersToSet(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	def := mustFields(t, "Note", "title:string,user:belongs_to:User")
	def.OwnedBy = "user"
	g := newTestGenerator(root, module, def)
	g.Architecture, g.Frontend = "triple", "next"

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "web", "lib", "portal", "notes.ts"))
	if strings.Contains(src, `"user_id"`) {
		t.Error("the owner column is on the customer's form, so they could give the row away")
	}
	if !strings.Contains(src, `scope: "user"`) {
		t.Error("the descriptor does not say who the rows belong to")
	}
	if !strings.Contains(src, `name: "title"`) {
		t.Error("the descriptor is missing an ordinary column")
	}
}

// A second run must not throw away the edit somebody made to the screen. The
// pages are theirs from the moment they are written.
func TestASecondRunKeepsAnEditedPage(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	g := ownedNoteGenerator(t, root, module)

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("first run: %v", err)
	}
	page := filepath.Join(root, "apps", "web", "app", "(app)", "notes", "page.tsx")
	writeTestFile(t, page, "// my masonry grid\n")

	if err := g.writePortalFiles(g.Names()); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if got := readTestFile(t, page); !strings.Contains(got, "my masonry grid") {
		t.Error("a second generate overwrote the customer screen, so adding a field costs you the page")
	}
	// And the sidebar gains nothing the second time.
	layout := readTestFile(t, filepath.Join(root, "apps", "web", "app", "(app)", "layout.tsx"))
	if n := strings.Count(layout, `href: "/notes"`); n != 1 {
		t.Errorf("the sidebar has %d entries for one resource, want 1", n)
	}
}

// Removing the resource removes the customer's half with it, the sidebar
// entry included: a link to a page that no longer exists is the one thing a
// removal must not leave behind.
func TestRemovingTakesTheCustomerScreensToo(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	g := ownedNoteGenerator(t, root, module)
	names := g.Names()

	if err := g.writePortalFiles(names); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	g.removePortalFiles(names)

	for _, rel := range []string{
		"app/(app)/notes",
		"components/portal/note-form.tsx",
		"lib/portal/notes.ts",
	} {
		if _, err := os.Stat(filepath.Join(root, "apps", "web", filepath.FromSlash(rel))); err == nil {
			t.Errorf("%s survived the removal", rel)
		}
	}
	layout := readTestFile(t, filepath.Join(root, "apps", "web", "app", "(app)", "layout.tsx"))
	if strings.Contains(layout, `href: "/notes"`) {
		t.Error("the sidebar still links to a page that is gone")
	}
}

// The level resolves from the project and the resource together, and --portal
// off means off even in a SaaS.
func TestPortalLevelResolution(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)

	g := ownedNoteGenerator(t, root, module)
	if got := g.portalLevel(); got != PortalScreens {
		t.Errorf("an owned resource in a SaaS resolved to %q, want screens", got)
	}
	g.Portal = PortalOff
	if got := g.portalLevel(); got != PortalOff {
		t.Errorf("--portal off resolved to %q", got)
	}

	// Organizations own the rows instead of a person, and that is still the
	// customer's data: a tenant-owned resource gets a screen.
	def := mustFields(t, "Doc", "title:string")
	def.TenantOwned = true
	tenant := newTestGenerator(root, module, def)
	tenant.Architecture, tenant.Frontend = "triple", "next"
	if got := tenant.portalLevel(); got != PortalScreens {
		t.Errorf("a tenant-owned resource resolved to %q, want screens", got)
	}
}

func TestValidPortalLevel(t *testing.T) {
	for _, in := range []string{"screens", "headless", "off", "", "SCREENS"} {
		if _, err := ValidPortalLevel(in); err != nil {
			t.Errorf("ValidPortalLevel(%q): %v", in, err)
		}
	}
	if _, err := ValidPortalLevel("pages"); err == nil {
		t.Error("an unknown level was accepted, so a typo writes nothing and says nothing")
	}
}
