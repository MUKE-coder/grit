package scaffold

import (
	"strings"
	"testing"
)

// New configs carry media-src; an old one gets it where a new one has it.
func TestCSPMediaRepairMatchesTheTemplates(t *testing.T) {
	for name, tmpl := range map[string]string{"next": nextSecurityHeaders(), "vite": viteSecurityHeaders()} {
		if !strings.Contains(tmpl, cspMediaSrc) {
			t.Fatalf("%s CSP has no media-src", name)
		}
		if out, changes, _ := repairCSPMediaSource(tmpl); out != tmpl || len(changes) != 0 {
			t.Errorf("%s: a current config was changed", name)
		}
		old := strings.Replace(tmpl, "  "+cspMediaSrc+"\n", "", 1)
		fixed, changes, warnings := repairCSPMediaSource(old)
		if fixed != tmpl || len(changes) != 1 || len(warnings) != 0 {
			t.Errorf("%s: the repair did not restore the template (%v, %v)", name, changes, warnings)
		}
	}
	conf := viteNginxConf("web")
	if !strings.Contains(conf, "media-src 'self' blob: https:") {
		t.Fatal("the nginx CSP has no media-src")
	}
	old := strings.Replace(conf, nginxImgMediaSrc, nginxImgSrc, 1)
	if fixed, changes, _ := repairCSPMediaNginxSource(old); fixed != conf || len(changes) != 1 {
		t.Error("the nginx repair did not restore the template")
	}
	if _, _, warnings := repairCSPMediaSource(`const csp = ["img-src *"]`); len(warnings) != 1 {
		t.Error("a hand-written CSP should get a warning")
	}
}
