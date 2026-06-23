package expose

import (
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/generate"
)

// autoFields is the critical security boundary for grit expose form —
// any field it leaks through ends up as a labelled input on a
// customer-facing page. Tests below pin the exact filter behavior so
// future refactors can't accidentally widen the surface area.

func TestAutoFields_DropsFrameworkColumns(t *testing.T) {
	s := &generate.GoStruct{
		Name: "Contact",
		Fields: []generate.GoField{
			{Name: "ID", GoType: "string", JSONName: "id"},
			{Name: "Name", GoType: "string", JSONName: "name"},
			{Name: "Version", GoType: "int", JSONName: "version"},
			{Name: "CreatedAt", GoType: "time.Time", JSONName: "created_at"},
			{Name: "UpdatedAt", GoType: "time.Time", JSONName: "updated_at"},
			{Name: "DeletedAt", GoType: "gorm.DeletedAt", JSONName: "deleted_at"},
		},
	}
	got := autoFields(s)
	if len(got) != 1 {
		t.Fatalf("want exactly 1 field (Name); got %d (%v)", len(got), got)
	}
	if got[0].JSONName != "name" {
		t.Errorf("want field 'name'; got %q", got[0].JSONName)
	}
}

func TestAutoFields_DropsPointerAndValueAssociations(t *testing.T) {
	// The contact-app Contact model has both:
	//   GroupID string  (FK column — should appear)
	//   Group   *Group  OR  Group Group  (association — must be dropped)
	s := &generate.GoStruct{
		Name: "Contact",
		Fields: []generate.GoField{
			{Name: "Name", GoType: "string", JSONName: "name"},
			{Name: "GroupID", GoType: "string", JSONName: "group_id"},
			{Name: "Group", GoType: "*Group", JSONName: "group"},
			{Name: "Author", GoType: "User", JSONName: "author"},
		},
	}
	got := autoFields(s)
	if len(got) != 2 {
		t.Fatalf("want 2 fields (Name + GroupID); got %d:\n%v", len(got), got)
	}
	for _, f := range got {
		if f.JSONName == "group" || f.JSONName == "author" {
			t.Errorf("association field %q leaked through autoFields", f.JSONName)
		}
	}
}

func TestAutoFields_DropsSliceAssociations(t *testing.T) {
	// Many-to-many associations come back as []SomeModel — these can't
	// be rendered as one <input>; they need a multi-select widget that
	// expose doesn't ship.
	s := &generate.GoStruct{
		Name: "Post",
		Fields: []generate.GoField{
			{Name: "Title", GoType: "string", JSONName: "title"},
			{Name: "Tags", GoType: "[]Tag", JSONName: "tags"},
			{Name: "Comments", GoType: "[]Comment", JSONName: "comments"},
		},
	}
	got := autoFields(s)
	if len(got) != 1 || got[0].JSONName != "title" {
		t.Errorf("want only 'title'; got %v", got)
	}
}

func TestAutoFields_KeepsPrimitives(t *testing.T) {
	// Sanity check — every primitive type isPrimitiveGoType recognises
	// should make it through.
	s := &generate.GoStruct{
		Name: "Mix",
		Fields: []generate.GoField{
			{Name: "S", GoType: "string", JSONName: "s"},
			{Name: "I", GoType: "int", JSONName: "i"},
			{Name: "U", GoType: "uint64", JSONName: "u"},
			{Name: "F", GoType: "float64", JSONName: "f"},
			{Name: "B", GoType: "bool", JSONName: "b"},
			{Name: "T", GoType: "time.Time", JSONName: "t"},
			{Name: "TP", GoType: "*time.Time", JSONName: "tp"},
		},
	}
	got := autoFields(s)
	if len(got) != 7 {
		t.Fatalf("want 7 primitive fields; got %d:\n%v", len(got), got)
	}
}

func TestAutoFields_LabelHandlesAcronyms(t *testing.T) {
	// "GroupID" was previously labelled "Group i d" because the
	// camelToSnake split treated I and D as separate words. Labels are
	// now derived from JSONName directly, which preserves acronyms.
	s := &generate.GoStruct{
		Name: "Contact",
		Fields: []generate.GoField{
			{Name: "GroupID", GoType: "string", JSONName: "group_id"},
		},
	}
	got := autoFields(s)
	if len(got) != 1 {
		t.Fatalf("want 1 field; got %d", len(got))
	}
	if got[0].Label != "Group id" {
		t.Errorf("want label 'Group id'; got %q", got[0].Label)
	}
}

// ── pluralPascal / pluralKebab regression tests ──────────────────────────

func TestPluralPascal(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Contact", "Contacts"},
		{"User", "Users"},
		{"Category", "Categories"},
		{"Address", "Addresses"},
	}
	for _, c := range cases {
		got := pluralPascal(c.in)
		if got != c.want {
			t.Errorf("pluralPascal(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}

func TestPluralKebab(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Contact", "contacts"},
		{"BlogPost", "blog_posts"},
		{"FormShare", "form_shares"},
	}
	for _, c := range cases {
		got := pluralKebab(c.in)
		if got != c.want {
			t.Errorf("pluralKebab(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}

// ── path validation ──────────────────────────────────────────────────────

func TestResolveTarget_RejectsNonTSX(t *testing.T) {
	_, err := resolveTarget(Opts{To: "apps/web/app/contact-us/page.ts", Root: "/proj"})
	if err == nil || !strings.Contains(err.Error(), ".tsx") {
		t.Errorf("want error mentioning .tsx; got %v", err)
	}
}

func TestResolveTarget_RequiresTo(t *testing.T) {
	_, err := resolveTarget(Opts{Root: "/proj"})
	if err == nil || !strings.Contains(err.Error(), "--to") {
		t.Errorf("want error mentioning --to; got %v", err)
	}
}

func TestResolveTarget_JoinsRelative(t *testing.T) {
	out, err := resolveTarget(Opts{To: "apps/web/app/page.tsx", Root: "/proj"})
	if err != nil {
		t.Fatal(err)
	}
	// filepath.Join uses OS-native separator; check both contains the
	// root and ends with the expected suffix.
	if !strings.Contains(out, "proj") || !strings.HasSuffix(out, "page.tsx") {
		t.Errorf("unexpected target path: %q", out)
	}
}
