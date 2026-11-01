package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanBlankMonorepoRemovesTanStackBlogRoutes(t *testing.T) {
	root := t.TempDir()
	writeStarterFile(t, root, "apps/web/src/routes/_site/blog/index.tsx", "export const Route = {}\n")
	writeStarterFile(t, root, "apps/web/src/routes/_site/blog/$slug.tsx", "export const Route = {}\n")
	writeStarterFile(t, root, "apps/web/src/routes/_site/index.tsx", `<Link
        to="/blog"
        className="button"
      >
        Read the Blog
      </Link>
`)
	writeStarterFile(t, root, "apps/web/src/lib/nav-menu.ts", `[
  {
    title: "Blog",
    href: "/blog",
  },
  {
    title: "About",
    href: "/about",
  },
]
`)
	writeStarterFile(t, root, "apps/web/__tests__/navbar.test.tsx", `expect(screen.getByRole("link", { name: /blog/i })).toBeInTheDocument()
expect(screen.getByRole("link", { name: /blog/i })).toHaveAttribute("href", "/blog")
`)

	if err := cleanBlankMonorepoLinks(root, "tanstack", false); err != nil {
		t.Fatalf("cleanBlankMonorepoLinks: %v", err)
	}

	for _, path := range []string{
		"apps/web/src/routes/_site/blog/index.tsx",
		"apps/web/src/routes/_site/blog/$slug.tsx",
	} {
		if _, err := os.Stat(filepath.Join(root, path)); !os.IsNotExist(err) {
			t.Errorf("%s still exists, want removed", path)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "apps/web/src/routes/_site/blog")); !os.IsNotExist(err) {
		t.Error("empty Blog route directory still exists, want removed")
	}

	home := readStarterFile(t, root, "apps/web/src/routes/_site/index.tsx")
	if strings.Contains(home, `to="/blog"`) {
		t.Error("landing page still links to the removed Blog route")
	}
	if !strings.Contains(home, `to="/admin/dashboard"`) {
		t.Error("landing page does not link to the embedded admin")
	}
	if !strings.Contains(home, "Open the admin panel") {
		t.Error("landing page still asks the reader to read the Blog")
	}

	nav := readStarterFile(t, root, "apps/web/src/lib/nav-menu.ts")
	if strings.Contains(nav, `title: "Blog"`) {
		t.Error("navbar still has a Blog link")
	}
	if !strings.Contains(nav, `title: "About"`) {
		t.Error("navbar lost a link that was not the Blog")
	}
}

func TestCleanBlankMonorepoSeparateAdminUsesAPlainLink(t *testing.T) {
	root := t.TempDir()
	writeStarterFile(t, root, "apps/web/src/routes/_site/index.tsx", `import { createFileRoute, Link } from '@tanstack/react-router'

function HomePage() {
  return (
    <Link
      to="/blog"
      className="button"
    >
      Read the Blog
    </Link>
  )
}
`)
	writeStarterFile(t, root, "apps/web/src/lib/nav-menu.ts", `[
  {
    title: "Blog",
    href: "/blog",
  },
]
`)
	writeStarterFile(t, root, "apps/admin/src/routes/_dashboard/resources/blogs/index.tsx", "import Page from '@/pages/resources/blogs'\n")
	writeStarterFile(t, root, "apps/admin/src/routes/_dashboard/resources/blogs/$id.tsx", "import { blogsResource } from '@/resources/blogs/blogs'\n")

	if err := cleanBlankMonorepoLinks(root, "tanstack", true); err != nil {
		t.Fatalf("cleanBlankMonorepoLinks: %v", err)
	}

	home := readStarterFile(t, root, "apps/web/src/routes/_site/index.tsx")
	if strings.Contains(home, "<Link") {
		t.Error("the landing page kept a router Link to a removed route")
	}
	if !strings.Contains(home, `href="http://localhost:3001"`) {
		t.Error("the landing page does not point at the separate admin app")
	}
	if strings.Contains(home, "createFileRoute, Link") {
		t.Error("the unused router Link import is still there")
	}
	if !strings.Contains(home, "Open the admin panel") {
		t.Error("the landing page still asks the reader to read the Blog")
	}

	for _, path := range []string{
		"apps/admin/src/routes/_dashboard/resources/blogs/index.tsx",
		"apps/admin/src/routes/_dashboard/resources/blogs/$id.tsx",
	} {
		if _, err := os.Stat(filepath.Join(root, path)); !os.IsNotExist(err) {
			t.Errorf("%s still exists, want removed", path)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "apps/admin/src/routes/_dashboard/resources/blogs")); !os.IsNotExist(err) {
		t.Error("empty admin Blog route directory still exists, want removed")
	}
}

func TestCleanBlankMonorepoRemovesExpoBlog(t *testing.T) {
	root := t.TempDir()
	writeStarterFile(t, root, "apps/expo/app/blogs/index.tsx", "export default function Blogs() {}\n")
	writeStarterFile(t, root, "apps/expo/app/blogs/new.tsx", "export default function NewBlog() {}\n")
	writeStarterFile(t, root, "apps/expo/hooks/use-blogs.ts", "export function useBlogs() {}\n")
	writeStarterFile(t, root, "apps/expo/app/(tabs)/explore.tsx", `const items = [
  { title: "Blogs", description: "Posts and articles", route: "/blogs" },
  { title: "Users", description: "People", route: "/users" },
]
`)
	writeStarterFile(t, root, "apps/expo/app/explore/content.tsx", `const items = [
  { icon: "newspaper-outline", title: "Posts", description: "Blog posts and announcements" },
  { icon: "people-outline", title: "Users", description: "People and roles" },
]
`)

	if err := cleanBlankMonorepoLinks(root, "next", false); err != nil {
		t.Fatalf("cleanBlankMonorepoLinks: %v", err)
	}

	for _, path := range []string{
		"apps/expo/app/blogs/index.tsx",
		"apps/expo/app/blogs/new.tsx",
		"apps/expo/hooks/use-blogs.ts",
	} {
		if _, err := os.Stat(filepath.Join(root, path)); !os.IsNotExist(err) {
			t.Errorf("%s still exists, want removed", path)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "apps/expo/app/blogs")); !os.IsNotExist(err) {
		t.Error("empty Expo Blog directory still exists, want removed")
	}

	explore := readStarterFile(t, root, "apps/expo/app/(tabs)/explore.tsx")
	if strings.Contains(explore, `title: "Blogs"`) {
		t.Error("Expo discover still has a Blogs entry")
	}
	if !strings.Contains(explore, `title: "Users"`) {
		t.Error("Expo discover lost an entry that was not the Blog")
	}

	content := readStarterFile(t, root, "apps/expo/app/explore/content.tsx")
	if strings.Contains(content, `title: "Posts"`) {
		t.Error("Expo content still has a Posts entry")
	}
	if !strings.Contains(content, `title: "Users"`) {
		t.Error("Expo content lost an entry that was not the Blog")
	}
}

func TestCleanBlankMonorepoRemovesDesktopBlogShortcut(t *testing.T) {
	root := t.TempDir()
	writeStarterFile(t, root, "apps/desktop/frontend/src/components/quick-access.tsx", `import {
  LayoutGrid, X, Home, RefreshCw, MessageSquare, FileText,
  Link as LinkIcon,
} from "lucide-react";

// Navigation shortcuts + system "new" flows shown alongside the per-resource
// cards. Support/blogs live under /system, so they get built-in entries.
const SYSTEM_ACTIONS = [
  { key: "sys:ticket", label: "New ticket", description: "Open a support ticket" },
  { key: "sys:blog", label: "New blog post", description: "Write and publish a post", to: "/app/system/blogs" },
];
`)

	if err := cleanBlankMonorepoLinks(root, "next", false); err != nil {
		t.Fatalf("cleanBlankMonorepoLinks: %v", err)
	}

	quick := readStarterFile(t, root, "apps/desktop/frontend/src/components/quick-access.tsx")
	if strings.Contains(quick, `"sys:blog"`) {
		t.Error("desktop quick access still opens the removed Blog screen")
	}
	if !strings.Contains(quick, `"sys:ticket"`) {
		t.Error("desktop quick access lost a shortcut that was not the Blog")
	}
	if strings.Contains(quick, "FileText") {
		t.Error("the unused Blog icon import is still there")
	}
	if !strings.Contains(quick, "LayoutGrid") {
		t.Error("the desktop icon import lost a symbol that was still in use")
	}
	if strings.Contains(quick, "Support/blogs") {
		t.Error("the desktop comment still mentions the removed Blog")
	}
}

// grit#173, and the reason this test exists at all.
//
// The starter cleanup edits files the scaffold owns, and the first version of
// this code recorded those edits in the manifest so the files would keep
// receiving framework fixes. That is what my review asked for and it was
// wrong: refreshing the hash hands the file back to the template, the
// template still contains the Blog, and the first `grit upgrade` put the Blog
// links back into nav-menu.ts pointing at a route the starter had removed.
//
// So the cleanup leaves the manifest alone, and the files it touches read as
// edited. The cost is that they stop receiving framework fixes; the
// alternative was an upgrade that silently undoes the starter.
//
// This asserts the choice rather than trusting a comment, because the
// mistake is an attractive one: refreshing looks obviously right until you
// run the upgrade.
func TestTheStarterCleanupLeavesTheManifestAlone(t *testing.T) {
	source, err := os.ReadFile("starter.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, call := range []string{"manifest.Refresh(", "manifest.Drop(", "manifest.DropTree("} {
		for _, line := range strings.Split(string(source), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") || !strings.Contains(line, call) {
				continue
			}
			t.Errorf("starter.go calls %s\n  %s\n\n"+
				"Refreshing an edited file hands it back to the template, which still has the "+
				"Blog in it, so the next upgrade undoes the starter. Dropping a deleted file's "+
				"entry does the same thing by making it untracked. See grit#173.", call, trimmed)
		}
	}
}

func writeStarterFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", rel, err)
	}
}

func readStarterFile(t *testing.T, root, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("reading %s: %v", rel, err)
	}
	return string(data)
}
