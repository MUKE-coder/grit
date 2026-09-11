package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// AddI18n installs internationalisation into an existing Grit project.
//
// This is the only implementation. `grit add i18n` calls it, and `grit new
// --i18n` calls it immediately after scaffolding rather than carrying a second
// set of conditional templates. Two code paths that must agree is how the two
// admin variants drifted, and one of them is enough.
//
// Everything here is additive. The base scaffold has no next-intl import, no
// locale middleware and no catalogues, so a project that never asks for i18n
// carries none of it: no unused dependency in package.json, no dead JSON, and
// no half-translated admin implying a feature that is not there.
//
// Idempotent. Files that exist are left alone and reported, and every injection
// checks for itself first, so running it twice is a no-op rather than a
// duplicate import.
type I18nResult struct {
	Written []string
	Skipped []string
	Wired   []string
}

// AddI18n writes the i18n files and wires them into the API and the admin.
func AddI18n(projectRoot string, force bool) (*I18nResult, error) {
	layout, err := detectLayout(projectRoot)
	if err != nil {
		return nil, err
	}

	res := &I18nResult{}

	// Record what this writes and edits, so grit upgrade can tell Grit's files
	// from yours. Before v3.222.0 nothing was recorded, and every upgrade of an
	// i18n project reported six files as edited by you and never updated them.
	// grit new --i18n calls this after the scaffold's own recording has
	// stopped; inside grit upgrade it joins that recording.
	release, err := manifest.Start(projectRoot, DefaultVersion, "i18n")
	if err != nil {
		return nil, err
	}
	released := false
	defer func() {
		if !released {
			_ = release()
		}
	}()

	files := map[string]string{}

	// ── API side ────────────────────────────────────────────────────
	// Translating what the API says is the half a frontend cannot do for
	// itself: it receives a message, not a reason.
	if layout.APIRoot != "" {
		files[filepath.Join(layout.APIRoot, "internal", "i18n", "i18n.go")] = i18nGo()
		files[filepath.Join(layout.APIRoot, "internal", "i18n", "locales", "en.json")] = i18nLocaleEN()
		files[filepath.Join(layout.APIRoot, "internal", "i18n", "locales", "fr.json")] = i18nLocaleFR()
		files[filepath.Join(layout.APIRoot, "internal", "i18n", "locales", "sw.json")] = i18nLocaleSW()
		files[filepath.Join(layout.APIRoot, "internal", "middleware", "locale.go")] = i18nMiddlewareGo()
		files[filepath.Join(layout.APIRoot, "internal", "response", "response.go")] = i18nResponseGo()
		files[filepath.Join(layout.APIRoot, "internal", "handlers", "i18n.go")] = i18nHandlerGo()
		files[filepath.Join(layout.APIRoot, "internal", "i18n", "i18n_test.go")] = i18nTestGo()
	}

	// ── Frontend side ───────────────────────────────────────────────
	for _, root := range layout.NextRoots {
		files[filepath.Join(root, "i18n", "request.ts")] = i18nRequestTS()
		files[filepath.Join(root, "lib", "locale.ts")] = i18nLocaleLibTS()
		files[filepath.Join(root, "components", "language-switcher.tsx")] = i18nSwitcherTSX()
		files[filepath.Join(root, "messages", "en.json")] = i18nMessagesEN()
		files[filepath.Join(root, "messages", "fr.json")] = i18nMessagesFR()
		files[filepath.Join(root, "messages", "sw.json")] = i18nMessagesSW()
		// The admin's translation seam, for a project older than it.
		if filepath.Base(root) == "admin" {
			files[filepath.Join(root, "lib", "i18n.tsx")] = adminI18nLib()
		}
	}

	module, err := detectModule(layout.APIRoot)
	if err != nil {
		return nil, err
	}

	for path, body := range files {
		if !force && fileExists(path) && !brokenSwitcher(path) {
			// A catalogue is yours to edit, so it is never replaced, but keys
			// added in a later release are merged in so they are not shown in
			// English in the middle of a translated page.
			if isCatalogue(path) {
				added, err := mergeCatalogue(path, strings.ReplaceAll(body, "{{MODULE}}", module))
				if err != nil {
					return nil, err
				}
				if added {
					res.Wired = append(res.Wired, filepath.Base(filepath.Dir(filepath.Dir(path)))+"/"+filepath.Base(path)+" new keys")
				}
			}
			res.Skipped = append(res.Skipped, path)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("creating %s: %w", filepath.Dir(path), err)
		}
		body = strings.ReplaceAll(body, "{{MODULE}}", module)
		if err := os.WriteFile(path, []byte(body), 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", path, err)
		}
		manifest.Note(path, body)
		res.Written = append(res.Written, path)
	}

	if err := wireI18n(layout, res); err != nil {
		return nil, err
	}
	released = true
	if err := release(); err != nil {
		return nil, fmt.Errorf("recording the i18n files: %w", err)
	}
	return res, nil
}

// brokenSwitcher reports a language switcher from before v3.222.0, which
// imported a dropdown-menu component neither app has and so failed the type
// check of every --i18n project. Replacing it can only fix things: while that
// component is missing, the file cannot compile, customised or not.
func brokenSwitcher(path string) bool {
	return filepath.Base(path) == "language-switcher.tsx" &&
		fileContains(path, "@/components/ui/dropdown-menu") &&
		!fileExists(filepath.Join(filepath.Dir(path), "ui", "dropdown-menu.tsx"))
}

// writeEdit saves an edit to a file Grit may track. The edited file is adopted
// into the manifest only when it was exactly what Grit last wrote, so it stays
// Grit's and the next upgrade can refresh it. A file somebody changed stays
// theirs, and upgrade goes on leaving it alone.
func writeEdit(path, content string) error {
	pristine := manifest.IsUnchanged(path)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	if pristine {
		manifest.Refresh(path)
	}
	return nil
}

// i18nLayout is the set of places i18n has to reach in whichever architecture
// this project was scaffolded as.
type i18nLayout struct {
	APIRoot   string
	NextRoots []string
}

func detectLayout(root string) (*i18nLayout, error) {
	l := &i18nLayout{}

	for _, candidate := range []string{
		filepath.Join(root, "apps", "api"), // triple / double / mobile
		root,                               // api-only and single
	} {
		if fileExists(filepath.Join(candidate, "go.mod")) {
			l.APIRoot = candidate
			break
		}
	}
	if l.APIRoot == "" {
		return nil, fmt.Errorf("no Go API found: run this from the root of a Grit project")
	}

	// Only Next.js apps get next-intl. The Vite and Expo clients have their own
	// i18n stories and wiring them here would be guessing.
	for _, candidate := range []string{
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "apps", "web"),
	} {
		if fileExists(filepath.Join(candidate, "next.config.ts")) {
			l.NextRoots = append(l.NextRoots, candidate)
		}
	}

	return l, nil
}

func detectModule(apiRoot string) (string, error) {
	data, err := os.ReadFile(filepath.Join(apiRoot, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("reading go.mod: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("no module line in go.mod")
}

// wireI18n performs the edits that turn the written files into a working
// feature: the middleware that resolves a locale, the route that lists them,
// and on the Next side the plugin, the provider and the dependency.
func wireI18n(l *i18nLayout, res *I18nResult) error {
	// ── API: register the middleware early, before anything that reports
	// an error, so those errors are already translatable.
	routes := filepath.Join(l.APIRoot, "internal", "routes", "routes.go")
	if fileExists(routes) {
		if err := injectAfterLine(routes,
			"r.Use(middleware.RequestID())",
			"\tr.Use(middleware.Locale())",
			&res.Wired, "locale middleware"); err != nil {
			return err
		}
	}

	for _, root := range l.NextRoots {
		name := filepath.Base(root)

		// ── package.json: the dependency
		if err := injectJSONDep(filepath.Join(root, "package.json"),
			`"next-intl": "^4.14.0",`, &res.Wired, name+" dependency"); err != nil {
			return err
		}

		// ── next.config.ts: the plugin. Without this next-intl never reads
		// i18n/request.ts and useTranslations throws at render.
		cfg := filepath.Join(root, "next.config.ts")
		if err := injectAfterLine(cfg,
			`import type { NextConfig } from "next";`,
			"import createNextIntlPlugin from \"next-intl/plugin\";\n\nconst withNextIntl = createNextIntlPlugin();",
			&res.Wired, name+" next.config plugin"); err != nil {
			return err
		}
		if err := replaceOnce(cfg,
			"export default nextConfig;",
			"export default withNextIntl(nextConfig);",
			&res.Wired, name+" config export"); err != nil {
			return err
		}

		// ── layout: the provider that lets client components translate.
		if err := wireLayout(filepath.Join(root, "app", "layout.tsx"), &res.Wired, name); err != nil {
			return err
		}

		// ── the admin's own components read the catalogue through lib/i18n,
		// which cannot import next-intl, so the layout hands it the messages.
		if fileExists(filepath.Join(root, "lib", "i18n.tsx")) {
			if err := wireAdminMessages(filepath.Join(root, "app", "layout.tsx"), &res.Wired, name); err != nil {
				return err
			}
		}

		// ── the switcher, where people can reach it. It was written and never
		// mounted, so nothing a user could click changed the language.
		header := filepath.Join(root, "components", "chrome", "PageHeader.tsx")
		if err := injectAfterLine(header, `import { DarkModeToggle } from "./DarkModeToggle";`,
			`import { LanguageSwitcher } from "@/components/language-switcher";`, &res.Wired, name+" switcher import"); err != nil {
			return err
		}
		if err := injectAfterLine(header, "<DarkModeToggle />", "          <LanguageSwitcher />",
			&res.Wired, name+" switcher"); err != nil {
			return err
		}
		navbar := filepath.Join(root, "components", "navbar.tsx")
		if err := injectAfterLine(navbar, `from "lucide-react";`,
			`import { LanguageSwitcher } from "@/components/language-switcher";`, &res.Wired, name+" switcher import"); err != nil {
			return err
		}
		if err := injectBeforeLine(navbar, "Admin CTA", "          <LanguageSwitcher />",
			&res.Wired, name+" switcher"); err != nil {
			return err
		}
	}
	return nil
}

// injectBeforeLine adds code on the line before the first line containing marker.
func injectBeforeLine(path, marker, code string, wired *[]string, label string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil // the file is optional for this architecture
	}
	content := string(data)
	norm := func(x string) string { return strings.Join(strings.Fields(x), " ") }
	if strings.Contains(norm(content), norm(code)) {
		return nil
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, marker) {
			out := append([]string{}, lines[:i]...)
			out = append(out, code)
			out = append(out, lines[i:]...)
			if err := writeEdit(path, strings.Join(out, "\n")); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}
			*wired = append(*wired, label)
			return nil
		}
	}
	return nil
}

// wireAdminMessages gives the admin's lib/i18n the same catalogue next-intl
// has, inside the provider wireLayout added.
func wireAdminMessages(path string, wired *[]string, label string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	const open = "<NextIntlClientProvider locale={locale} messages={messages}>"
	if strings.Contains(content, "<I18nProvider") || !strings.Contains(content, open) {
		return nil
	}
	content = strings.Replace(content, open, open+"\n          <I18nProvider messages={messages}>", 1)
	content = strings.Replace(content, "</NextIntlClientProvider>", "</I18nProvider>\n        </NextIntlClientProvider>", 1)
	content = strings.Replace(content, `import { NextIntlClientProvider } from "next-intl";`,
		`import { NextIntlClientProvider } from "next-intl";`+"\n"+`import { I18nProvider } from "@/lib/i18n";`, 1)
	if err := writeEdit(path, content); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	*wired = append(*wired, label+" admin catalogue")
	return nil
}

// injectAfterLine adds code on the line after the first line containing marker.
func injectAfterLine(path, marker, code string, wired *[]string, label string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil // the file is optional for this architecture
	}
	content := string(data)
	// Whitespace-insensitive, so a reformatted file still counts as wired and
	// running this twice does not produce a duplicate registration.
	norm := func(x string) string { return strings.Join(strings.Fields(x), " ") }
	if strings.Contains(norm(content), norm(code)) {
		return nil
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, marker) {
			out := append([]string{}, lines[:i+1]...)
			out = append(out, code)
			out = append(out, lines[i+1:]...)
			if err := writeEdit(path, strings.Join(out, "\n")); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}
			*wired = append(*wired, label)
			return nil
		}
	}
	return nil
}

func replaceOnce(path, old, new string, wired *[]string, label string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.Contains(content, new) {
		return nil
	}
	if !strings.Contains(content, old) {
		return nil
	}
	content = strings.Replace(content, old, new, 1)
	if err := writeEdit(path, content); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	*wired = append(*wired, label)
	return nil
}

// injectJSONDep adds a dependency line to package.json.
//
// Textual rather than a parse-and-reserialise, because rewriting the whole file
// reorders keys and reformats the parts nobody touched, turning a one-line
// change into an unreviewable diff.
func injectJSONDep(path, dep string, wired *[]string, label string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.Contains(content, `"next-intl"`) {
		// next-intl 3 predates Next 16, whose builds run on Turbopack, and its
		// plugin never registers its config there: every production build of
		// an --i18n app stopped with "Couldn't find next-intl config file".
		// 4 supports Next 16 and needs nothing else changed here.
		if bumped := nextIntl3.ReplaceAllString(content, `"next-intl": "^4.14.0"`); bumped != content {
			if err := writeEdit(path, bumped); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}
			*wired = append(*wired, label+" next-intl 4")
		}
		return nil
	}
	marker := `"dependencies": {`
	idx := strings.Index(content, marker)
	if idx == -1 {
		return nil
	}
	at := idx + len(marker)
	content = content[:at] + "\n    " + dep + content[at:]
	if err := writeEdit(path, content); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	*wired = append(*wired, label)
	return nil
}

// wireLayout wraps the root layout in NextIntlClientProvider.
//
// The layout also becomes async and its <html lang> starts carrying the real
// locale, which matters for screen readers and for the browser's own offer to
// translate the page.
func wireLayout(path string, wired *[]string, label string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.Contains(content, "NextIntlClientProvider") {
		return nil
	}

	content = strings.Replace(content,
		`import type { Metadata } from "next";`,
		`import type { Metadata } from "next";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";`, 1)

	content = strings.Replace(content,
		`export default function RootLayout({`,
		`export default async function RootLayout({`, 1)

	// Resolve the locale just inside the function body.
	if i := strings.Index(content, "}) {"); i != -1 {
		at := i + len("}) {")
		content = content[:at] + "\n  const locale = await getLocale();\n  const messages = await getMessages();" + content[at:]
	}

	content = strings.Replace(content, `<html lang="en"`, `<html lang={locale}`, 1)

	// Wrap whatever the body contains, rather than matching a specific tree.
	//
	// The admin layout renders <Providers>{children}</Providers> and the web
	// layout renders <Providers><AppChrome>{children}</AppChrome></Providers>.
	// Matching either exactly means the other silently keeps its import and
	// gains no provider, and useTranslations then throws at runtime in one app
	// and not the other. <body> is the one anchor both share.
	openIdx := strings.Index(content, "<body")
	if openIdx == -1 {
		return fmt.Errorf("%s: no <body> to wrap; layout not in the expected shape", path)
	}
	openEnd := strings.Index(content[openIdx:], ">")
	if openEnd == -1 {
		return fmt.Errorf("%s: unterminated <body> tag", path)
	}
	openEnd += openIdx + 1
	closeIdx := strings.LastIndex(content, "</body>")
	if closeIdx == -1 || closeIdx < openEnd {
		return fmt.Errorf("%s: no closing </body>", path)
	}

	inner := strings.TrimRight(content[openEnd:closeIdx], " \t\n")
	wrapped := "\n        <NextIntlClientProvider locale={locale} messages={messages}>" +
		strings.ReplaceAll(inner, "\n", "\n  ") +
		"\n        </NextIntlClientProvider>\n      "
	content = content[:openEnd] + wrapped + content[closeIdx:]

	// Silent failure is what this whole function is guarding against, so check.
	if !strings.Contains(content, "<NextIntlClientProvider") {
		return fmt.Errorf("%s: provider wrap produced nothing", path)
	}

	if err := writeEdit(path, content); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	*wired = append(*wired, label+" layout provider")
	return nil
}
