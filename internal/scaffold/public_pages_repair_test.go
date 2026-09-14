package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

func TestPublicPagesRenderOnTheServer(t *testing.T) {
	pages := map[string]string{
		"home":      webLandingPage(Options{ProjectName: "demo", Architecture: ArchDouble, Frontend: FrontendNext}),
		"blog list": webBlogListPage(),
		"blog post": webBlogDetailPage(),
	}
	for name, src := range pages {
		if strings.Contains(src, `"use client"`) {
			t.Errorf("the %s page is a client component", name)
		}
		if strings.Contains(src, "@/hooks/use-blogs") || strings.Contains(src, "isLoading") {
			t.Errorf("the %s page still fetches in the browser", name)
		}
		if !strings.Contains(src, `from "@/lib/blog-api"`) {
			t.Errorf("the %s page does not read through lib/blog-api", name)
		}
	}
	if !strings.Contains(pages["blog post"], "export async function generateMetadata(") || !strings.Contains(pages["blog post"], "notFound()") {
		t.Error("a blog post has no metadata of its own, or no 404")
	}
	if !strings.Contains(pages["blog list"], "export const metadata") || !strings.Contains(pages["blog list"], "/blog?page=") {
		t.Error("the blog list has no metadata, or no page links")
	}
	// grit remove resource Blog cuts these out of the home page.
	for _, marker := range []string{"// grit:home:blog-import", "// grit:home:blog-hook-start", "// grit:home:blog-hook-end", "{/* grit:home:blog-start */}", "{/* grit:home:blog-end */}"} {
		if !strings.Contains(pages["home"], marker) {
			t.Errorf("the home page lost %s", marker)
		}
	}
	lib := webBlogAPILib()
	for _, want := range []string{"API_INTERNAL_URL", "revalidate: REVALIDATE_SECONDS", "cache(", "status === 404"} {
		if !strings.Contains(lib, want) {
			t.Errorf("lib/blog-api.ts is missing %s", want)
		}
	}
}

func TestRepairBlogSanitize(t *testing.T) {
	fresh := strings.ReplaceAll(blogHandlerGo(), "{{MODULE}}", "demo")
	old := strings.Replace(fresh, blogListSanitize, "", 1)
	old = strings.Replace(old, blogSlugSanitize, "", 1)
	old = strings.Replace(old, "\t\"demo/internal/sanitize\"\n", "", 1)
	if old == fresh {
		t.Fatal("could not rebuild blog_handler.go from before the repair")
	}
	out, fixed, warn := repairBlogSanitizeSource(old, "demo")
	if len(warn) > 0 || len(fixed) != 1 || strings.Count(out, "sanitize.HTML(") != 2 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairBlogSanitizeSource(out, "demo"); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed blog_handler.go again")
	}
	if out, fixed, _ := repairBlogSanitizeSource(fresh, "demo"); out != fresh || len(fixed) > 0 {
		t.Error("a fresh blog_handler.go still needs the repair")
	}
}

func TestComposeWebAPIEnv(t *testing.T) {
	next := dockerComposeProd(Options{ProjectName: "demo", Architecture: ArchDouble, Frontend: FrontendNext})
	if !strings.Contains(next, "API_INTERNAL_URL: http://api:8080") {
		t.Error("a Next.js web container is not told where the API is")
	}
	if strings.Contains(dockerComposeProd(Options{ProjectName: "demo", Architecture: ArchDouble, Frontend: FrontendTanStack}), "API_INTERNAL_URL") {
		t.Error("a TanStack SPA, which renders nothing on the server, got API_INTERNAL_URL")
	}
	old := strings.Replace(next, composeWebAPIEnv, "", 1)
	out, fixed, warn := repairComposeWebAPIEnvSource(old)
	if len(warn) > 0 || len(fixed) != 1 || out != next {
		t.Errorf("fixed %v, warned %v; the repaired compose file differs from a fresh one", fixed, warn)
	}
	if again, fixed, _ := repairComposeWebAPIEnvSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the compose file again")
	}
}
