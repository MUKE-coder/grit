package scaffold

import (
	"strings"
	"testing"
)

// grit#173: the four framework files that hard-coded the Blog example.
//
// Removing the Blog meant editing them, and `grit upgrade` leaves an edited
// file alone, so the project kept its removal and stopped receiving framework
// fixes to those files for ever. Recording the edits instead, which was tried
// and reverted in v3.462.0, had the upgrade rewrite them from a template that
// still contained the Blog: the nav links came back as 404s and seed.go stopped
// compiling on `undefined: SeedBlogs`.
//
// The templates ask now, so the file the removal leaves behind IS what the
// template produces, the upgrade guard adopts it, and both properties hold.

func TestTheNavMenuHasNoBlogLinksWithoutABlog(t *testing.T) {
	opts := Options{ProjectName: "shop", Architecture: ArchTriple, Frontend: FrontendNext}

	with := webNavMenuConfig(opts, true)
	if !strings.Contains(with, `href: "/blog"`) {
		t.Error("a project WITH a blog has no blog links, so this test is the wrong way round")
	}

	without := webNavMenuConfig(opts, false)
	if strings.Contains(without, "/blog") {
		t.Error("the nav menu still points at /blog in a project with no blog: four 404s " +
			"in the navbar, which is what the blank starter used to edit this file to fix")
		for _, line := range strings.Split(without, "\n") {
			if strings.Contains(line, "/blog") {
				t.Logf("    %s", strings.TrimSpace(line))
			}
		}
	}
	if strings.Contains(without, `title: "Blog"`) {
		t.Error("a menu entry is still titled Blog")
	}
	// An unused import fails the frontend lint, which a generated project
	// should never ship.
	for _, icon := range []string{"Rss", "BookOpen"} {
		if strings.Contains(without, icon) {
			t.Errorf("%s is imported and nothing uses it now the blog links are gone", icon)
		}
	}
	// And the menu still has to be a menu.
	for _, want := range []string{"Admin panel", "Accounts", "Support", "export const navEntries"} {
		if !strings.Contains(without, want) {
			t.Errorf("%q went with the blog links", want)
		}
	}
}

func TestTheLandingPageHasNoBlogWithoutABlog(t *testing.T) {
	opts := Options{ProjectName: "shop", Architecture: ArchTriple, Frontend: FrontendNext}

	with := webLandingPage(opts, true)
	if !strings.Contains(with, "getPublishedBlogs") {
		t.Error("a project WITH a blog has no blog section, so this test is the wrong way round")
	}

	without := webLandingPage(opts, false)
	for _, gone := range []string{
		"getPublishedBlogs", // the data hook
		"@/lib/blog-api",    // its import
		"Read the Blog",     // the hero's call to action, a 404 without a blog
		"Recent Posts",      // the section
		"grit:home:blog",    // and the markers with it
	} {
		if strings.Contains(without, gone) {
			t.Errorf("the landing page still has %q with no blog to show", gone)
		}
	}
	for _, want := range []string{"Open the admin panel", "<DevLinks />", "export default async function HomePage"} {
		if !strings.Contains(without, want) {
			t.Errorf("%q went with the blog section", want)
		}
	}
}

func TestSeedGoCallsSeedBlogsOnlyWhenThereIsABlog(t *testing.T) {
	with := apiSeedGo(true)
	if !strings.Contains(with, "SeedBlogs(db)") {
		t.Error("a project WITH a blog does not seed it")
	}

	without := apiSeedGo(false)
	if strings.Contains(without, "SeedBlogs") {
		t.Error("seed.go still calls SeedBlogs with no blogs_seeder.go to define it, " +
			"which is `undefined: SeedBlogs` and a project that does not compile")
	}
	for _, want := range []string{"SeedUsers(db)", "// grit:seeders", "func Seed(db *gorm.DB) error {"} {
		if !strings.Contains(without, want) {
			t.Errorf("%q went missing from seed.go", want)
		}
	}
	// The spacing matters more than it looks. The upgrade guard adopts a file
	// whose contents already match what the template would write and holds back
	// one that differs, so a stray blank line where the Blog block was is the
	// difference between seed.go receiving framework fixes again and never
	// receiving another one.
	if strings.Contains(without, "\n\n\n") {
		t.Error("seed.go has a double blank line where the blog block was; the removal's " +
			"output will not match this template and the upgrade will hold the file back")
	}
}

// The trash screen used to be edited over a word in a doc comment, which cost
// it every later framework fix. The example names a table nobody removes now.
func TestTheTrashScreenDoesNotNameARemovableResource(t *testing.T) {
	src := adminTrashPage()
	if strings.Contains(src, "blog_posts") {
		t.Error("the trash screen's humanise() example still names blog_posts, so the " +
			"blank starter has to edit this file, and an edited framework file never " +
			"receives another framework fix")
	}
	if !strings.Contains(src, "function humanise(") {
		t.Fatal("humanise() is gone, so this test is checking the wrong file")
	}
}

func TestWithoutRegionLeavesAnUnmatchedTemplateAlone(t *testing.T) {
	src := "a\n// start\nb\n// end\nc"
	if got := withoutRegion(src, "// start", "// end"); got != "a\nc" {
		t.Errorf("withoutRegion = %q, want %q", got, "a\nc")
	}
	// A template whose shape has changed must not be half-cut: better to leave
	// the blog section in than to delete from a marker to the end of the file.
	if got := withoutRegion(src, "// start", "// nope"); got != src {
		t.Errorf("a missing end marker cut anyway: %q", got)
	}
	if got := withoutRegion(src, "// nope", "// end"); got != src {
		t.Errorf("a missing start marker cut anyway: %q", got)
	}
}
