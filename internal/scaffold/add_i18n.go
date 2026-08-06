package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	}

	module, err := detectModule(layout.APIRoot)
	if err != nil {
		return nil, err
	}

	for path, body := range files {
		if !force && fileExists(path) {
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
		res.Written = append(res.Written, path)
	}

	if err := wireI18n(layout, res); err != nil {
		return nil, err
	}
	return res, nil
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
			`"next-intl": "^3.26.0",`, &res.Wired, name+" dependency"); err != nil {
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
	}
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
			if err := os.WriteFile(path, []byte(strings.Join(out, "\n")), 0644); err != nil {
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
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
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
		return nil
	}
	marker := `"dependencies": {`
	idx := strings.Index(content, marker)
	if idx == -1 {
		return nil
	}
	at := idx + len(marker)
	content = content[:at] + "\n    " + dep + content[at:]
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
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

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	*wired = append(*wired, label+" layout provider")
	return nil
}
