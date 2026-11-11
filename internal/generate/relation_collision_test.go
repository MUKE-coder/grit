package generate

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// grit#209: two belongs_to fields pointing at the SAME model generated a
// desktop form that did not compile.
//
//	grants-form.tsx(25,9): error TS2451: Cannot redeclare block-scoped variable 'usersOpts'
//	grants-form.tsx(26,9): error TS2451: Cannot redeclare block-scoped variable 'usersOpts'
//
// The options hook was declared once per relation and named after the TARGET
// model, so the second declaration collided with the first. "A credit granted
// to one user by another" is an ordinary shape, and `--owned-by user` makes it
// easier to hit by accident because the ownership column is a third relation
// to the same model.
//
// pnpm type-check fails for the whole monorepo, so it fails CI too.

// twoRelationsToOneModel is the reporter's resource.
func twoRelationsToOneModel(t *testing.T) *Generator {
	t.Helper()
	def, err := ParseInlineFields("Grant",
		"amount:uint,approved_by:belongs_to[User]:optional,requested_by:belongs_to[User]:optional")
	if err != nil {
		t.Fatalf("parsing the fields: %v", err)
	}
	// The shape they actually had: plus --owned-by user, a third relation to
	// the same model.
	def.ApplyFlags(false, false, "user", false, false, false, false)
	def.EnsureOwnerField()

	var toUser int
	for _, f := range def.Fields {
		if f.IsBelongsTo() && f.RelatedModelName() == "User" {
			toUser++
		}
	}
	if toUser < 3 {
		t.Fatalf("the fixture has %d relations to User, want 3: with fewer than two there is "+
			"nothing to collide and this test cannot see grit#209", toUser)
	}
	return &Generator{Definition: def, Module: "app"}
}

// declaredConsts is every `const x =` declared at one level of a TSX body.
var declaredConst = regexp.MustCompile(`(?m)^\s*const (\w+)\s*=`)

// duplicates is the names declared more than once, which in a TSX function
// body is TS2451 and a build that fails.
func duplicates(src string) []string {
	count := map[string]int{}
	for _, m := range declaredConst.FindAllStringSubmatch(src, -1) {
		count[m[1]]++
	}
	var out []string
	for name, n := range count {
		if n > 1 {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func TestTheDesktopFormDeclaresEachOptionsListOnce(t *testing.T) {
	g := twoRelationsToOneModel(t)
	src := g.desktopClientForm(g.Names())

	if dup := duplicates(src); len(dup) > 0 {
		t.Errorf("the desktop form redeclares %v, which is TS2451 and fails pnpm type-check for "+
			"the whole monorepo (grit#209)", dup)
	}
}

// Named after the field, not the target model. The reporter's preferred fix,
// and the readable one: two selects both reading usersOpts says nothing about
// which is which.
//
// Deduplicating by target model would also compile and is the wrong fix: both
// selects would share one hook result, which happens to be right for User and
// is silently wrong the moment the two relations want different filtering.
func TestTheDesktopOptionsAreNamedAfterTheField(t *testing.T) {
	g := twoRelationsToOneModel(t)
	src := g.desktopClientForm(g.Names())

	for _, want := range []string{"approvedByOpts", "requestedByOpts", "userOpts"} {
		if !strings.Contains(src, "const "+want+" =") {
			t.Errorf("no %s: the options list is not named after its field", want)
		}
	}
	// One hook import for the three relations, which was already right.
	if n := strings.Count(src, `from "@/hooks/use-users"`); n != 1 {
		t.Errorf("the users hook is imported %d times, want 1", n)
	}
	// And three call sites, so each select has its own list.
	if n := strings.Count(src, "useUsers()"); n != 3 {
		t.Errorf("useUsers() is called %d times, want 3: one per relation", n)
	}
}

// The mobile form has the same generator shape and the same bug, with two
// variables per relation rather than one. Not in the report: found by grepping
// for the pattern the desktop bug was an instance of.
func TestTheMobileFormDeclaresEachOptionsListOnce(t *testing.T) {
	g := twoRelationsToOneModel(t)
	src := g.mobileFormSource(g.Names())

	if dup := duplicates(src); len(dup) > 0 {
		t.Errorf("the Expo form redeclares %v (grit#209, the same cause on the other client)", dup)
	}
	for _, want := range []string{"approvedByOpts", "requestedByOpts", "userOpts"} {
		if !strings.Contains(src, "const "+want+" =") {
			t.Errorf("no %s in the Expo form", want)
		}
	}
}

// One relation per target is the common case and must be unchanged, because
// every existing project is in it.
func TestOneRelationPerModelIsUnchanged(t *testing.T) {
	def, err := ParseInlineFields("Invoice", "amount:uint,customer:belongs_to[User]")
	if err != nil {
		t.Fatal(err)
	}
	g := &Generator{Definition: def, Module: "app"}
	src := g.desktopClientForm(g.Names())

	if dup := duplicates(src); len(dup) > 0 {
		t.Errorf("a single relation now collides with itself: %v", dup)
	}
	if !strings.Contains(src, "const customerOpts =") {
		t.Error("the options list is not named after the field")
	}
}
