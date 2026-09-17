package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review L7, L8 and L9, the frontend half of the Low security
// findings. The new text lives here once: the templates splice it in, and
// repairFrontendLinkSafety puts it into an existing project, anchored on what
// Grit generated before.

// L7: img-src named every https: host.
const (
	cspImgSrcTight = "\"img-src 'self' data: blob: \" + API_ORIGIN + \" \" + STORAGE_ORIGIN + \" \" + IMAGE_ORIGINS,"

	nextIsDevLine    = "const isDev = process.env.NODE_ENV !== \"production\";\n"
	nextImageOrigins = `// Hosts images load from besides this site, the API and storage: the avatars
// Google and GitHub sign-in save, plus NEXT_PUBLIC_IMAGE_ORIGINS (comma-separated
// origins, for a CDN or an identity provider's avatars). img-src used to allow
// every https: host, so markup that got into a page could load, and report to,
// any server.
const IMAGE_ORIGINS = ["https://lh3.googleusercontent.com", "https://avatars.githubusercontent.com"]
  .concat((process.env.NEXT_PUBLIC_IMAGE_ORIGINS || "").split(",").map((s) => s.trim()).filter(Boolean).map((s) => toOrigin(s)))
  .concat(isDev ? ["https://picsum.photos"] : [])
  .join(" ");
`

	viteIsDevLine    = "const isDev = viteMode !== 'production'\n"
	viteImageOrigins = `// Hosts images load from besides this site, the API and storage: the avatars
// Google and GitHub sign-in save, plus VITE_IMAGE_ORIGINS (comma-separated
// origins). img-src used to allow every https: host.
const IMAGE_ORIGINS = ['https://lh3.googleusercontent.com', 'https://avatars.githubusercontent.com']
  .concat((viteEnv.VITE_IMAGE_ORIGINS || '').split(',').map((s) => s.trim()).filter(Boolean).map((s) => toOrigin(s)))
  .concat(isDev ? ['https://picsum.photos'] : [])
  .join(' ')
`
)

// L8: links built from stored values.

// adminSafeHrefTS is lib/safe-href.ts in the admin panel.
func adminSafeHrefTS() string {
	return `// safeHref returns a stored value as a link target only when following it is
// safe: an http or https URL, a mailto: or tel: address, or a path on this site.
// Anything else (javascript:, data:, vbscript:, a protocol-relative //host)
// comes back as "#".
//
// React refuses javascript: URLs, but not the others, and the values that reach
// these links are written by people outside the team: public form submissions,
// notification payloads.
const SAFE_HREF = /^(?:https?:\/\/|mailto:|tel:|\/(?![/\\]))/i;

export function safeHref(value: string | null | undefined): string {
  const href = typeof value === "string" ? value.trim() : "";
  return SAFE_HREF.test(href) ? href : "#";
}
`
}

// adminSafeHrefTest is __tests__/safe-href.test.ts in the admin panel.
func adminSafeHrefTest(libImport string) string {
	return `import { describe, it, expect } from "vitest";
import { safeHref } from "` + libImport + `/safe-href";

describe("safeHref", () => {
  it("keeps web, mail and phone links and paths on this site", () => {
    for (const href of ["https://example.com/a", "http://example.com", "mailto:a@example.com", "tel:+256700000000", "/system/jobs"]) {
      expect(safeHref(href)).toBe(href);
    }
  });

  it("refuses every other scheme and protocol-relative hosts", () => {
    for (const href of ["javascript:alert(1)", " JavaScript:alert(1)", "java\tscript:alert(1)", "data:text/html,<script>alert(1)</script>", "vbscript:msgbox(1)", "//evil.example", "/\\evil.example", "", null, undefined]) {
      expect(safeHref(href)).toBe("#");
    }
  });
});
`
}

const (
	linkCellOld = `function LinkCell({ value }: { value: string }) {
  let hostname = value;
  try {
    hostname = new URL(value).hostname;
  } catch {
    // use raw value if not a valid URL
  }
  return (
    <a
      href={value}
`
	linkCellNew = `function LinkCell({ value }: { value: string }) {
  // Public form submissions reach this cell, so a value is a link only with a
  // scheme that is safe to follow. The rest show as text.
  const href = safeHref(value);
  if (href === "#") {
    return <span className="text-sm text-text-secondary">{value}</span>;
  }
  let hostname = value;
  try {
    hostname = new URL(value).hostname;
  } catch {
    // use raw value if not a valid URL
  }
  return (
    <a
      href={href}
`
	fileRefHrefOld      = "      href={value.url}\n"
	fileRefHrefNew      = "      href={safeHref(value.url)}\n"
	notificationHrefOld = `href={n.link || "#"}`
	notificationHrefNew = `href={safeHref(n.link)}`
)

// adminIconsImportRe finds the icons import, which every file L8 touches has,
// with the alias the file uses (@admin in the embedded panel).
var adminIconsImportRe = regexp.MustCompile(`(?m)^import \{[^}\n]*\} from "(@|@admin)/lib/icons";\n`)

// addSafeHrefImport puts the safeHref import under the icons import.
func addSafeHrefImport(src string) (string, bool) {
	if strings.Contains(src, "/lib/safe-href\"") {
		return src, true
	}
	loc := adminIconsImportRe.FindStringSubmatchIndex(src)
	if loc == nil {
		return src, false
	}
	alias := src[loc[2]:loc[3]]
	return src[:loc[1]] + `import { safeHref } from "` + alias + `/lib/safe-href";` + "\n" + src[loc[1]:], true
}

// L9: the SSO redirect.
const (
	ssoImportOld = `import { apiUrl } from "@/lib/api-core";`
	ssoImportNew = `import { API_URL, apiUrl } from "@/lib/api-core";`

	ssoAnchor   = "// SSOSignIn is the enterprise entry point."
	ssoRedirect = `// apiRedirect turns a path the API answered with into an address on the API,
// or null when it is not a path or resolves to another origin. Appending the
// value to the API's address let "@evil.example" become
// http://api.example.com@evil.example, which is evil.example.
function apiRedirect(path: unknown): string | null {
  if (typeof path !== "string" || !path.startsWith("/") || path.startsWith("//") || path.startsWith("/\\")) {
    return null;
  }
  try {
    const api = new URL(API_URL, window.location.href);
    const target = new URL(apiUrl(path), api);
    return target.origin === api.origin ? target.href : null;
  } catch {
    return null;
  }
}

`
	ssoFollowOld = `      if (res.ok && body?.data?.sso && body.data.redirect_url) {
        window.location.href = apiUrl(body.data.redirect_url);
        return;
      }
`
	ssoFollowNew = `      const target = res.ok && body?.data?.sso ? apiRedirect(body.data.redirect_url) : null;
      if (target) {
        window.location.href = target;
        return;
      }
`
)

// repairFrontendLinkSafety brings L7, L8 and L9 to an existing project. None of
// the files involved are delivered whole by upgrade.
func repairFrontendLinkSafety(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	// L7: the CSP in each frontend's config.
	for _, app := range []string{"admin", "web", "docs"} {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js", "vite.config.ts"} {
			if path := filepath.Join(root, "apps", app, name); fileExists(path) {
				if err := repairTextFile(root, m, path, repairCSPImageOriginsSource); err != nil {
					return err
				}
			}
		}
	}

	// L8: the admin panel in each of its shapes.
	for _, lib := range []string{
		filepath.Join(root, "apps", "admin", "lib"),
		filepath.Join(root, "apps", "admin", "src", "lib"),
		filepath.Join(root, "apps", "web", "admin-panel", "lib"),
	} {
		if err := repairSafeHrefPanel(root, m, lib); err != nil {
			return err
		}
	}

	// L9: the SSO sign-in, in the admin and the web app.
	for _, path := range []string{
		filepath.Join(root, "apps", "admin", "components", "auth", "SocialAuthButtons.tsx"),
		filepath.Join(root, "apps", "admin", "src", "components", "auth", "SocialAuthButtons.tsx"),
		filepath.Join(root, "apps", "web", "components", "auth", "SocialAuthButtons.tsx"),
		filepath.Join(root, "apps", "web", "admin-panel", "components", "auth", "SocialAuthButtons.tsx"),
	} {
		if fileExists(path) {
			if err := repairTextFile(root, m, path, repairSSORedirectSource); err != nil {
				return err
			}
		}
	}
	return nil
}

// repairSafeHrefPanel adds lib/safe-href.ts to one admin panel and uses it in
// the cells and notification links.
func repairSafeHrefPanel(root string, m *manifest.Manifest, lib string) error {
	if !fileExists(filepath.Join(lib, "icons.ts")) {
		return nil
	}
	panel := filepath.Dir(lib)
	if err := writeFile(filepath.Join(lib, "safe-href.ts"), adminSafeHrefTS()); err != nil {
		return err
	}
	for _, path := range []string{
		filepath.Join(panel, "components", "tables", "cell-renderers.tsx"),
		filepath.Join(panel, "components", "chrome", "NotificationBell.tsx"),
		filepath.Join(panel, "app", "(dashboard)", "system", "notifications", "page.tsx"),
		filepath.Join(panel, "pages", "system", "notifications.tsx"),
		filepath.Join(filepath.Dir(panel), "app", "admin", "(dashboard)", "system", "notifications", "page.tsx"),
	} {
		if fileExists(path) {
			if err := repairTextFile(root, m, path, repairSafeHrefSource); err != nil {
				return err
			}
		}
	}
	return nil
}

func repairCSPImageOriginsSource(src string) (string, []string, []string) {
	if strings.Contains(src, "IMAGE_ORIGINS") || !strings.Contains(src, "img-src") {
		return src, nil, nil
	}
	const warning = "the Content-Security-Policy is not the one Grit wrote: drop https: from img-src and name the image origins you use, or injected markup can load any server's images"
	isDev, origins := nextIsDevLine, nextImageOrigins
	if strings.Count(src, viteIsDevLine) == 1 {
		isDev, origins = viteIsDevLine, viteImageOrigins
	}
	old := ""
	for _, candidate := range []string{cspImgSrcNew, cspImgSrcOld} {
		if strings.Count(src, candidate) == 1 {
			old = candidate
		}
	}
	if old == "" || strings.Count(src, isDev) != 1 {
		return src, nil, []string{warning}
	}
	out := strings.Replace(src, old, cspImgSrcTight, 1)
	out = strings.Replace(out, isDev, isDev+origins, 1)
	return out, []string{"img-src names the API, storage and sign-in avatar origins instead of every https: host"}, nil
}

func repairSafeHrefSource(src string) (string, []string, []string) {
	if strings.Contains(src, "safeHref(") {
		return src, nil, nil
	}
	out := src
	if strings.Count(out, linkCellOld) == 1 {
		out = strings.Replace(out, linkCellOld, linkCellNew, 1)
		out = strings.Replace(out, fileRefHrefOld, fileRefHrefNew, 1)
	}
	out = strings.ReplaceAll(out, notificationHrefOld, notificationHrefNew)
	if out == src {
		if strings.Contains(src, "function LinkCell(") || strings.Contains(src, "n.link") {
			return src, nil, []string{"this file is not the one Grit wrote: pass stored links through safeHref from lib/safe-href, or a javascript: or //other-host value becomes a clickable link"}
		}
		return src, nil, nil
	}
	out, ok := addSafeHrefImport(out)
	if !ok {
		return src, nil, []string{"this file has no icons import to add the safeHref import under: pass stored links through safeHref from lib/safe-href"}
	}
	return out, []string{"stored links are clickable only with an http, https, mailto or tel scheme, or a path on this site"}, nil
}

func repairSSORedirectSource(src string) (string, []string, []string) {
	if strings.Contains(src, "function apiRedirect(") || !strings.Contains(src, "redirect_url") {
		return src, nil, nil
	}
	importOld, importNew := ssoImportOld, ssoImportNew
	if strings.Contains(src, `"@admin/lib/api-core"`) {
		importOld = strings.Replace(importOld, "@/", "@admin/", 1)
		importNew = strings.Replace(importNew, "@/", "@admin/", 1)
	}
	if strings.Count(src, importOld) != 1 || strings.Count(src, ssoAnchor) != 1 || strings.Count(src, ssoFollowOld) != 1 {
		return src, nil, []string{"SocialAuthButtons.tsx is not the file Grit wrote: resolve redirect_url with new URL(apiUrl(path), API_URL) and follow it only when its origin is the API's, or a crafted value redirects to another site"}
	}
	out := strings.Replace(src, importOld, importNew, 1)
	out = strings.Replace(out, ssoAnchor, ssoRedirect+ssoAnchor, 1)
	out = strings.Replace(out, ssoFollowOld, ssoFollowNew, 1)
	return out, []string{"the SSO redirect is followed only when it stays on the API's origin"}, nil
}
