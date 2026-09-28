package generate

import (
	"os"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// HasRichtext reports whether any field is rendered as HTML, which is what makes
// the resource's write routes interesting to Sentinel's WAF.
//
// Richtext only, not text: every resource has a description or a note, and
// excluding all of them from body inspection would leave the WAF inspecting
// nothing worth inspecting.
func (d *ResourceDefinition) HasRichtext() bool {
	for i := range d.Fields {
		if FieldType(d.Fields[i].Type) == FieldRichtext {
			return true
		}
	}
	return false
}

// wafRichtextEntries is the resource's exclusion, as it goes into
// wafExcludedRoutes' paths. The text comes from internal/scaffold so that grit
// upgrade, which adds the same block to a project generated before v3.335.0,
// writes something grit remove can match.
func wafRichtextEntries(names Names) string {
	return scaffold.WAFRichtextExclusion(names.Plural)
}

// injectWAFRichtext excludes the resource's routes from WAF body inspection.
//
// Silently skipped when the marker is absent, which is every project scaffolded
// before v3.335.0: failing a generate over a WAF entry would be a poor trade,
// and grit upgrade adds both the marker and the entry.
func injectWAFRichtext(routesFile string, names Names) bool {
	if !fileExists(routesFile) {
		return false
	}
	data, err := os.ReadFile(routesFile)
	if err != nil || !strings.Contains(string(data), scaffold.WAFRichtextMarker) {
		return false
	}
	return injectBefore(routesFile, scaffold.WAFRichtextMarker, wafRichtextEntries(names)) == nil
}

// removeWAFRichtext takes the resource's exclusion back out of
// wafExcludedRoutes.
//
// Matching is on the whole block, comment included, so a list somebody has
// edited by hand is left alone rather than half-cut. grit upgrade writes the
// block with the line endings the file already had, so both forms are tried.
func removeWAFRichtext(routesFile string, names Names) bool {
	if !fileExists(routesFile) {
		return false
	}
	data, err := os.ReadFile(routesFile)
	if err != nil {
		return false
	}
	content := string(data)
	lf := wafRichtextEntries(names) + "\n"
	for _, block := range []string{lf, strings.ReplaceAll(lf, "\n", "\r\n")} {
		if !strings.Contains(content, block) {
			continue
		}
		return os.WriteFile(routesFile, []byte(strings.Replace(content, block, "", 1)), 0o644) == nil
	}
	return false
}
