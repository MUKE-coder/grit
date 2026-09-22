package scaffold

import (
	"strings"
	"testing"
)

// A new config lets next/image read localhost in development; an old one gets
// the same line, in the same place.
func TestNextImageLocalRepairMatchesTheTemplate(t *testing.T) {
	tmpl := nextSecurityHeadersConfig()
	if !strings.Contains(tmpl, "dangerouslyAllowLocalIP: isDev") {
		t.Fatal("the Next.js image config still refuses localhost")
	}
	if out, changes, _ := repairNextImageLocalSource(tmpl); out != tmpl || len(changes) != 0 {
		t.Error("a current config was changed")
	}
	old := strings.Replace(tmpl, nextImagesNew, nextImagesOld, 1)
	fixed, changes, warnings := repairNextImageLocalSource(old)
	if fixed != tmpl || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("the repair did not restore the template (%v, %v)", changes, warnings)
	}
	if _, _, warnings := repairNextImageLocalSource("images: { remotePatterns: [] }"); len(warnings) != 1 {
		t.Error("a hand-written image config should get a warning")
	}
}
