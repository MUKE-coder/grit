package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M26: the admin painted the wrong theme on every load, and
// carried a second theme system nothing rendered.
//
// The theme was applied in an effect inside DarkModeToggle, which sits in the
// dashboard chrome and so mounts only after /auth/me answers. A dark dashboard
// showed a light one first, for as long as that request took. Alongside it,
// components/shared/theme-provider.tsx stored a "grit-theme" key nothing else
// read, and its only consumer was components/layout/navbar.tsx, which no layout
// has rendered since v3.29.
//
// A new project gets the fixed files from the scaffold. An upgrade rewrites the
// three framework files whole when they are untouched; this repairs the copies
// somebody edited, and removes the two orphans, which upgrade would otherwise
// leave sitting in the tree forever.

const (
	themeProviderImportDefault = `import { ThemeProvider } from "./theme-provider";`
	themeScriptMarker          = "const themeScript ="
	darkToggleFixedMarker      = `getAttribute("data-theme-mode")`
	darkToggleOldEffect        = `  useEffect(() => {
    const stored = (typeof window !== "undefined"
      ? (window.localStorage.getItem("grit-theme-mode") as Mode | null)
      : null);
    // Prefer stored choice; fall back to OS-level preference; default light.
    const osDark = typeof window !== "undefined"
      && window.matchMedia
      && window.matchMedia("(prefers-color-scheme: dark)").matches;
    const initial: Mode = stored || (osDark ? "dark" : "light");
    setMode(initial);
    applyMode(initial);
    setHydrated(true);
  }, []);`
	darkToggleNewEffect = `  useEffect(() => {
    // The layout script already wrote data-theme-mode. Read it rather than
    // deciding again: two places deciding is how they disagree.
    const applied = document.documentElement.getAttribute("data-theme-mode");
    if (applied === "dark" || applied === "light") {
      setMode(applied);
      setHydrated(true);
      return;
    }
    // Fallback for a host document that carries no theme script (a panel
    // mounted inside an app Grit did not write the shell for). Same decision
    // the script makes, just late enough to flash.
    const stored = window.localStorage.getItem("grit-theme-mode") as Mode | null;
    const osDark = window.matchMedia
      && window.matchMedia("(prefers-color-scheme: dark)").matches;
    const initial: Mode = stored === "dark" || stored === "light" ? stored : (osDark ? "dark" : "light");
    setMode(initial);
    applyMode(initial);
    setHydrated(true);
  }, []);`
)

// repairThemeFlash applies M26 to every admin panel in the project.
func repairThemeFlash(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, shape := range adminPanelShapes(root) {
		providers := filepath.Join(shape.code, "components", "shared", "providers.tsx")
		if fileExists(providers) {
			if err := repairTextFile(root, m, providers, repairProvidersThemeSource); err != nil {
				return err
			}
		}
		layout := filepath.Join(shape.routes, "layout.tsx")
		if fileExists(layout) {
			if err := repairTextFile(root, m, layout, repairAdminLayoutThemeSource); err != nil {
				return err
			}
		}
		toggle := filepath.Join(shape.code, "components", "chrome", "DarkModeToggle.tsx")
		if fileExists(toggle) {
			if err := repairTextFile(root, m, toggle, repairDarkModeToggleSource); err != nil {
				return err
			}
		}
		// Removed only once nothing imports them, so a panel somebody wired the
		// old provider into keeps compiling.
		//
		// The navbar goes first, and it matters: the navbar is the only thing
		// that ever imported the theme provider, so checking the provider while
		// the navbar is still on disk finds a reference and keeps both. That
		// left the provider behind until a second upgrade.
		for _, orphan := range []struct{ path, needle, shown string }{
			{filepath.Join(shape.code, "components", "layout", "navbar.tsx"), "layout/navbar", "components/layout/navbar.tsx"},
			{filepath.Join(shape.code, "components", "shared", "theme-provider.tsx"), "theme-provider", "components/shared/theme-provider.tsx"},
		} {
			if removeUnreferencedAdminFile(shape, orphan.path, orphan.needle) {
				fmt.Printf("  ✓ %s: removed, nothing imported it\n", orphan.shown)
			}
		}
	}
	return nil
}

// repairProvidersThemeSource unwraps <ThemeProvider> and drops its import.
func repairProvidersThemeSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "ThemeProvider") {
		return src, nil, nil
	}
	var kept []string
	dropped := 0
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "theme-provider") && strings.HasPrefix(trimmed, "import ") {
			dropped++
			continue
		}
		if trimmed == "<ThemeProvider>" || trimmed == "</ThemeProvider>" {
			dropped++
			continue
		}
		kept = append(kept, line)
	}
	out := strings.Join(kept, "\n")
	if dropped == 0 || strings.Contains(out, "ThemeProvider") {
		return src, nil, []string{"components/shared/providers.tsx still wraps the app in a ThemeProvider: remove it, the theme is applied by the script in the layout now"}
	}
	return out, []string{"the second theme system is gone from the provider tree"}, nil
}

// repairAdminLayoutThemeSource puts the blocking theme script in the layout.
func repairAdminLayoutThemeSource(src string) (string, []string, []string) {
	if strings.Contains(src, themeScriptMarker) {
		return src, nil, nil
	}
	anchor := "\nexport default function "
	if !strings.Contains(src, anchor) {
		return src, nil, nil
	}
	decl := "\n" + adminThemeScriptSource() + "\n"
	out := strings.Replace(src, anchor, decl+anchor, 1)

	script := `<script dangerouslySetInnerHTML={{ __html: themeScript }} />`
	switch {
	case strings.Contains(out, "<head>"):
		out = strings.Replace(out, "<head>", "<head>\n        "+script, 1)
	case strings.Contains(out, "<html "):
		// The document layout Grit writes has no <head> of its own.
		idx := strings.Index(out, "\n      <body")
		if idx < 0 {
			return src, nil, []string{"the admin's root layout is not the file Grit wrote: add a <head> containing <script dangerouslySetInnerHTML={{ __html: themeScript }} /> so the stored theme is applied before the first paint"}
		}
		head := "\n      <head>\n        " + script + "\n      </head>"
		out = out[:idx] + head + out[idx:]
	case strings.Contains(out, "return <Providers>{children}</Providers>;"):
		// The panel embedded in the web app: the web app's root layout owns the
		// document, so the script rides at the top of this layout instead.
		out = strings.Replace(out, "return <Providers>{children}</Providers>;",
			`return (
    <>
      `+script+`
      <Providers>{children}</Providers>
    </>
  );`, 1)
	default:
		return src, nil, []string{"the admin's layout is not the file Grit wrote: render <script dangerouslySetInnerHTML={{ __html: themeScript }} /> before the panel, so the stored theme is applied before the first paint"}
	}
	return out, []string{"the stored theme is applied before the first paint, not in an effect"}, nil
}

// repairDarkModeToggleSource makes the toggle read the applied theme.
func repairDarkModeToggleSource(src string) (string, []string, []string) {
	if strings.Contains(src, darkToggleFixedMarker) {
		return src, nil, nil
	}
	if !strings.Contains(src, darkToggleOldEffect) {
		return src, nil, nil
	}
	return strings.Replace(src, darkToggleOldEffect, darkToggleNewEffect, 1),
		[]string{"the toggle reads the theme the layout already applied"}, nil
}

// removeUnreferencedAdminFile deletes path when no other file in the panel
// mentions needle. Reports whether it removed anything.
func removeUnreferencedAdminFile(shape adminPanelShape, path, needle string) bool {
	if !fileExists(path) {
		return false
	}
	for _, dir := range []string{shape.code, shape.routes} {
		referenced := false
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || referenced {
				return nil //nolint:nilerr // an unreadable tree is not a reference
			}
			switch filepath.Ext(p) {
			case ".ts", ".tsx":
			default:
				return nil
			}
			if sameFilePath(p, path) {
				return nil
			}
			body, readErr := os.ReadFile(p)
			if readErr == nil && strings.Contains(string(body), needle) {
				referenced = true
			}
			return nil
		})
		if referenced {
			return false
		}
	}
	return os.Remove(path) == nil
}

// sameFilePath compares two paths as the filesystem would.
func sameFilePath(a, b string) bool {
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}
