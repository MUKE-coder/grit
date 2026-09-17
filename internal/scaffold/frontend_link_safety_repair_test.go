package scaffold

import (
	"strings"
	"testing"
)

const safeHrefImport = `import { safeHref } from "@/lib/safe-href";` + "\n"

func TestFrontendCSPNamesImageOrigins(t *testing.T) {
	for name, src := range map[string]string{"next": nextSecurityHeaders(), "vite": viteSecurityHeaders()} {
		if strings.Contains(src, "blob: https:") {
			t.Errorf("%s: img-src still allows every https: host", name)
		}
		if strings.Count(src, cspImgSrcTight) != 1 || !strings.Contains(src, "const IMAGE_ORIGINS") {
			t.Errorf("%s: img-src does not name its origins", name)
		}
		for _, directive := range []string{"object-src 'none'", "base-uri 'self'", "frame-ancestors 'none'"} {
			if !strings.Contains(src, directive) {
				t.Errorf("%s: CSP lacks %s", name, directive)
			}
		}
	}
}

func TestCSPImageOriginsRepair(t *testing.T) {
	next := nextSecurityHeaders()
	oldNext := strings.Replace(next, nextImageOrigins, "", 1)
	checkRepair(t, "next.config", repairCSPImageOriginsSource, strings.Replace(oldNext, cspImgSrcTight, cspImgSrcNew, 1), next)
	// Before the storage repair added API_ORIGIN.
	checkRepair(t, "next.config (pre-storage)", repairCSPImageOriginsSource, strings.Replace(oldNext, cspImgSrcTight, cspImgSrcOld, 1), next)

	vite := viteSecurityHeaders()
	oldVite := strings.Replace(strings.Replace(vite, viteImageOrigins, "", 1), cspImgSrcTight, cspImgSrcNew, 1)
	checkRepair(t, "vite.config", repairCSPImageOriginsSource, oldVite, vite)

	// The storage repair, which runs later, leaves the result alone.
	if out, fixed, _ := repairCSPImageSource(next); out != next || len(fixed) != 0 {
		t.Error("the storage CSP repair rewrites the tightened img-src")
	}
}

func TestStoredLinksPassThroughSafeHref(t *testing.T) {
	cells := adminCellRenderers()
	if !strings.Contains(cells, linkCellNew) || !strings.Contains(cells, fileRefHrefNew) || !strings.Contains(cells, safeHrefImport) {
		t.Error("cell-renderers.tsx does not pass stored links through safeHref")
	}
	for name, src := range map[string]string{"bell": adminNotificationBellComponent(), "page": adminNotificationsPage()} {
		if strings.Contains(src, notificationHrefOld) || !strings.Contains(src, notificationHrefNew) || !strings.Contains(src, safeHrefImport) {
			t.Errorf("%s: notification links do not pass through safeHref", name)
		}
	}
	if !strings.Contains(adminSafeHrefTS(), "export function safeHref(") {
		t.Error("lib/safe-href.ts does not export safeHref")
	}
}

func TestSafeHrefRepair(t *testing.T) {
	cells := adminCellRenderers()
	oldCells := strings.Replace(cells, linkCellNew, linkCellOld, 1)
	oldCells = strings.Replace(oldCells, fileRefHrefNew, fileRefHrefOld, 1)
	oldCells = strings.Replace(oldCells, safeHrefImport, "", 1)
	checkRepair(t, "cell-renderers", repairSafeHrefSource, oldCells, cells)

	for name, src := range map[string]string{"bell": adminNotificationBellComponent(), "page": adminNotificationsPage()} {
		old := strings.Replace(strings.ReplaceAll(src, notificationHrefNew, notificationHrefOld), safeHrefImport, "", 1)
		checkRepair(t, name, repairSafeHrefSource, old, src)
	}

	// The embedded panel imports through @admin.
	embedded := embedAdminContent(cells, nil)
	oldEmbedded := embedAdminContent(oldCells, nil)
	checkRepair(t, "embedded cell-renderers", repairSafeHrefSource, oldEmbedded, embedded)
	if !strings.Contains(embedded, `"@admin/lib/safe-href"`) {
		t.Error("the embedded panel does not import safe-href through @admin")
	}
}

func TestSSORedirectStaysOnTheAPI(t *testing.T) {
	src := adminAuthSocialButtons()
	if strings.Contains(src, "apiUrl(body.data.redirect_url)") || !strings.Contains(src, "target.origin === api.origin") {
		t.Error("SocialAuthButtons follows redirect_url without checking its origin")
	}
	old := strings.Replace(src, ssoImportNew, ssoImportOld, 1)
	old = strings.Replace(old, ssoRedirect, "", 1)
	old = strings.Replace(old, ssoFollowNew, ssoFollowOld, 1)
	checkRepair(t, "SocialAuthButtons", repairSSORedirectSource, old, src)
	checkRepair(t, "embedded SocialAuthButtons", repairSSORedirectSource, embedAdminContent(old, nil), embedAdminContent(src, nil))
}
