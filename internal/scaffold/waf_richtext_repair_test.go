package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The bug in grit#90: a richtext resource was excluded from WAF body inspection
// at its public path and not at the /admin path the admin panel writes through,
// so a post whose body quoted a shell command could be published once and then
// never edited.
func TestScaffoldedRoutesExcludeTheAdminWrites(t *testing.T) {
	src := apiRoutesGo()
	for _, want := range []string{
		`"/blogs", "/blogs/*",`,
		`"/admin/blogs", "/admin/blogs/*",`,
		`"/admin/posts", "/admin/posts/*",`,
		`"/admin/articles", "/admin/articles/*",`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("wafExcludedRoutes is missing %s", want)
		}
	}
	// Every admin write route of a richtext resource must be covered, or the
	// exclusion is decoration.
	for _, route := range []string{
		`staff.POST("/admin/blogs",`,
		`staff.PUT("/admin/blogs/:id",`,
	} {
		if !strings.Contains(src, route) {
			t.Errorf("the scaffolded blog no longer registers %s, so this test guards nothing", route)
		}
	}
	if !strings.Contains(src, WAFRichtextMarker) {
		t.Error("routes.go has no grit:waf:richtext marker, so grit generate cannot add a richtext resource")
	}
}

func TestRepairWAFRichtextAddsTheAdminTwins(t *testing.T) {
	fresh := apiRoutesGo()
	old := fresh
	for _, plural := range []string{"blogs", "posts", "articles"} {
		old = strings.Replace(old, "\t\t\"/admin/"+plural+"\", \"/admin/"+plural+"/*\",\n", "", 1)
	}
	old = strings.Replace(old, "\t\t"+WAFRichtextMarker+"\n", "", 1)
	if old == fresh {
		t.Fatal("the fixture is the fresh template, so the repair is not being exercised")
	}

	out, fixed, warn := repairWAFRichtextSource(old, nil)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{
		`"/admin/blogs", "/admin/blogs/*",`,
		`"/admin/posts", "/admin/posts/*",`,
		`"/admin/articles", "/admin/articles/*",`,
		WAFRichtextMarker,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repair did not add %s", want)
		}
	}
	// Uploads is excluded for the body cap, not for markup, and has no admin
	// route: an entry for one would be dead config.
	if strings.Contains(out, `"/admin/uploads"`) {
		t.Error("the repair invented an admin route for uploads")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("the repaired routes.go is not valid Go: %v", err)
	}
	if again, fixed, _ := repairWAFRichtextSource(out, nil); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed routes.go again")
	}
	if out, fixed, _ := repairWAFRichtextSource(fresh, nil); out != fresh || len(fixed) > 0 {
		t.Error("a fresh routes.go still needs the repair")
	}
}

func TestRepairWAFRichtextAddsAGeneratedResource(t *testing.T) {
	out, fixed, warn := repairWAFRichtextSource(apiRoutesGo(), []string{"stories"})
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if !strings.Contains(fixed[0], "stories") {
		t.Errorf("the report does not name the resource it added: %q", fixed[0])
	}
	for _, want := range []string{
		`"/stories", "/stories/*",`,
		`"/admin/stories", "/admin/stories/*",`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repair did not add %s", want)
		}
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairWAFRichtextSource(out, []string{"stories"}); again != out || len(fixed) > 0 {
		t.Error("a second upgrade added the resource twice")
	}
}

func TestRepairWAFRichtextLeavesAnEditedListAlone(t *testing.T) {
	src := strings.Replace(apiRoutesGo(), "paths := []string{\n", "paths := wafPaths()\n\t_ = paths\n\t{\n", 1)
	out, fixed, warn := repairWAFRichtextSource(src, []string{"stories"})
	if out != src || len(fixed) > 0 {
		t.Error("a hand-written exclusion list was edited")
	}
	if len(warn) != 1 || !strings.Contains(warn[0], "/admin/<resource>") {
		t.Errorf("the developer was not told what to add by hand: %v", warn)
	}
}

func TestStructsWithHTMLSanitize(t *testing.T) {
	src := `package models

type Story struct {
	Title string ` + "`" + `json:"title"` + "`" + `
	Body  string ` + "`" + `json:"body" sanitize:"html"` + "`" + `
}

type Tag struct {
	Name string ` + "`" + `json:"name"` + "`" + `
}
`
	got := structsWithHTMLSanitize(src)
	if len(got) != 1 || got[0] != "Story" {
		t.Errorf("got %v, want [Story]", got)
	}
}

func TestRoutePluralsForBothRouteShapes(t *testing.T) {
	inline := map[string]string{"routes.go": "\t\tm.Staff.PUT(\"/admin/stories/:id\", middleware.RequireRole(\"ADMIN\"), storyHandler.Update)\n" +
		"\t\tm.Staff.PUT(\"/admin/tags/:id\", tagHandler.Update)\n"}
	if got := routePluralsFor(inline, "Story"); len(got) != 1 || got[0] != "stories" {
		t.Errorf("inline routes: got %v, want [stories]", got)
	}

	// A per-resource routes file: every route in it belongs to the one
	// resource, and the handler is called h.
	split := map[string]string{"blog_post_routes.go": "\t\tm.Admin.PUT(\"/blog-posts/:id\", h.Update)\n"}
	if got := routePluralsFor(split, "BlogPost"); len(got) != 1 || got[0] != "blog-posts" {
		t.Errorf("split routes: got %v, want [blog-posts]", got)
	}
}
