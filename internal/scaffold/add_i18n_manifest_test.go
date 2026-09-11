package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

func writeTestFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

const i18nTestLayout = `import type { Metadata } from "next";
import { Providers } from "@/components/shared/providers";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
`

const i18nTestHeader = `import { DarkModeToggle } from "./DarkModeToggle";

export function PageHeader() {
  return (
    <div>
          <DarkModeToggle />
    </div>
  );
}
`

// An --i18n project failed its type check (the switcher imported a
// dropdown-menu neither app has), never mounted the switcher, gave the admin's
// own components nothing to translate with, and left six files every upgrade
// reported as edited by you.
func TestAddI18nMountsRecordsAndRepairs(t *testing.T) {
	root := t.TempDir()
	admin := filepath.Join(root, "apps", "admin")
	writeTestFile(t, filepath.Join(root, "apps", "api", "go.mod"), "module shop/apps/api\n")
	files := map[string]string{
		filepath.Join(admin, "next.config.ts"):                        "import type { NextConfig } from \"next\";\n\nconst nextConfig: NextConfig = {};\n\nexport default nextConfig;\n",
		filepath.Join(admin, "package.json"):                          "{\n  \"dependencies\": {\n    \"next\": \"15.0.0\"\n  }\n}\n",
		filepath.Join(admin, "app", "layout.tsx"):                     i18nTestLayout,
		filepath.Join(admin, "components", "chrome", "PageHeader.tsx"): i18nTestHeader,
		filepath.Join(admin, "lib", "i18n.tsx"):                        adminI18nLib(),
	}

	// Recorded as Grit wrote them, then package.json edited by hand.
	release, err := manifest.Start(root, "3.221.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	for path, body := range files {
		writeTestFile(t, path, body)
		manifest.Note(path, body)
	}
	if err := release(); err != nil {
		t.Fatal(err)
	}
	pkg := filepath.Join(admin, "package.json")
	writeTestFile(t, pkg, "{\n  \"dependencies\": {\n    \"next\": \"15.0.0\",\n    \"zod\": \"3.0.0\"\n  }\n}\n")

	// The switcher from before v3.222.0, which cannot compile here.
	switcher := filepath.Join(admin, "components", "language-switcher.tsx")
	writeTestFile(t, switcher, "import { DropdownMenu } from '@/components/ui/dropdown-menu'\n")

	if _, err := AddI18n(root, false); err != nil {
		t.Fatalf("AddI18n: %v", err)
	}

	got, _ := os.ReadFile(switcher)
	if strings.Contains(string(got), "dropdown-menu") || !strings.Contains(string(got), "<select") {
		t.Error("the broken switcher was left in place")
	}
	header, _ := os.ReadFile(filepath.Join(admin, "components", "chrome", "PageHeader.tsx"))
	if !strings.Contains(string(header), "<LanguageSwitcher />") ||
		!strings.Contains(string(header), `import { LanguageSwitcher } from "@/components/language-switcher";`) {
		t.Errorf("the switcher is not mounted:\n%s", header)
	}
	layout, _ := os.ReadFile(filepath.Join(admin, "app", "layout.tsx"))
	for _, want := range []string{"<I18nProvider messages={messages}>", `import { I18nProvider } from "@/lib/i18n";`, "</I18nProvider>"} {
		if !strings.Contains(string(layout), want) {
			t.Errorf("the admin layout does not feed lib/i18n: missing %s\n%s", want, layout)
		}
	}

	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	status := func(path string) manifest.Status {
		rel, _ := manifest.Rel(root, path)
		return m.StatusOf(root, rel)
	}
	// Grit's own edits to a file Grit wrote keep it Grit's, so the next upgrade
	// can refresh it rather than reporting it as edited by you.
	for _, path := range []string{filepath.Join(admin, "app", "layout.tsx"), filepath.Join(admin, "next.config.ts")} {
		if s := status(path); s != manifest.Unchanged {
			t.Errorf("%s reads as %v after i18n edited it", filepath.Base(path), s)
		}
	}
	// A file someone had edited stays theirs.
	if s := status(pkg); s == manifest.Unchanged {
		t.Error("an edited package.json was adopted as Grit's, so an upgrade would overwrite the edit")
	}
	// The files i18n wrote are recorded.
	if s := status(filepath.Join(admin, "messages", "fr.json")); s != manifest.Unchanged {
		t.Errorf("messages/fr.json reads as %v", s)
	}

	// Idempotent: a second run changes nothing and mounts nothing twice.
	if _, err := AddI18n(root, false); err != nil {
		t.Fatal(err)
	}
	header, _ = os.ReadFile(filepath.Join(admin, "components", "chrome", "PageHeader.tsx"))
	if n := strings.Count(string(header), "<LanguageSwitcher />"); n != 1 {
		t.Errorf("the switcher is mounted %d times", n)
	}
}

// next-intl 3 supports Next up to 15. On the scaffold's Next 16 every
// production build of an --i18n app stopped with "Couldn't find next-intl
// config file", so new projects get 4 and an existing 3 is moved up.
func TestNextIntlIsVersionFour(t *testing.T) {
	dir := t.TempDir()
	fresh := filepath.Join(dir, "fresh.json")
	old := filepath.Join(dir, "old.json")
	writeTestFile(t, fresh, "{\n  \"dependencies\": {\n    \"next\": \"^16.1.6\"\n  }\n}\n")
	writeTestFile(t, old, "{\n  \"dependencies\": {\n    \"next-intl\": \"^3.26.0\",\n    \"next\": \"^16.1.6\"\n  }\n}\n")

	var wired []string
	for _, p := range []string{fresh, old} {
		if err := injectJSONDep(p, `"next-intl": "^4.14.0",`, &wired, "app"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(p)
		if !strings.Contains(string(got), `"next-intl": "^4.14.0"`) || strings.Contains(string(got), `"^3.`) {
			t.Errorf("%s does not depend on next-intl 4:\n%s", filepath.Base(p), got)
		}
		if n := strings.Count(string(got), `"next-intl"`); n != 1 {
			t.Errorf("%s names next-intl %d times", filepath.Base(p), n)
		}
	}
}

// The guard reads the manifest once, at the start of an upgrade, so a file
// adopted during the run was still held back and reported as edited by you.
func TestGuardAdoptReleasesAFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "apps", "admin", "package.json")
	writeTestFile(t, path, "original\n")
	release, err := manifest.Start(root, "3.221.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	manifest.Note(path, "original\n")
	if err := release(); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, path, "wired\n")

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	defer stopGuard()
	if guardAllows(path, "template\n") {
		t.Fatal("an edited file was not held back")
	}
	guardAdopt(path, "wired\n")
	if !guardAllows(path, "template\n") {
		t.Error("an adopted file is still held back")
	}
}

// The admin translates through lib/i18n, which cannot import next-intl.
func TestAdminComponentsTranslate(t *testing.T) {
	for what, src := range map[string]string{
		"sidebar":     adminCollapsibleSidebarComponent(Options{ProjectName: "app"}),
		"data table":  adminDataTable(),
		"toolbar":     adminTableToolbar(),
		"pagination":  adminTablePagination(),
		"form":        adminFormBuilder(),
		"form modal":  adminFormModal(),
		"form page":   adminFormPage(),
		"empty state": adminTableEmptyState(),
	} {
		if !strings.Contains(src, `from "@/lib/i18n"`) || !strings.Contains(src, "t(\"") {
			t.Errorf("the %s does not translate", what)
		}
		if strings.Contains(src, `from "next-intl"`) {
			t.Errorf("the %s imports next-intl, which a project without i18n does not have", what)
		}
	}
	if lib := adminI18nLib(); strings.Contains(lib, `from "next-intl"`) || !strings.Contains(lib, "export function useLocalizedResource(") {
		t.Error("lib/i18n.tsx must stand alone and localize resources")
	}
}
