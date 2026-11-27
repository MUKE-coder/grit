package scaffold

import (
	"strings"
	"testing"
)

// The list page says how long rows are kept, from the window --retain writes
// into the definition (grit#260). Both halves travel together: a page reading
// a field the type does not declare is a build error, and a declared field the
// page never reads is a policy nobody sees.
func TestTheListPageShowsTheRetentionWindow(t *testing.T) {
	types := adminResourceTypes()
	for _, want := range []string{"retention?: string;", "export function retentionPhrase(window: string): string"} {
		if !strings.Contains(types, want) {
			t.Errorf("lib/resource.ts is missing %q", want)
		}
	}

	for name, page := range map[string]string{"next": adminResourcePage(), "tanstack": nextToTanStack(adminResourcePage())} {
		if !strings.Contains(page, `import { retentionPhrase } from "@/lib/resource";`) {
			t.Errorf("%s: the page does not import retentionPhrase", name)
		}
		// Behind the field, so a resource with no window shows no sentence
		// rather than one that says "undefined".
		if !strings.Contains(page, "{resource.retention && (") ||
			!strings.Contains(page, "{retentionPhrase(resource.retention)}") {
			t.Errorf("%s: the page does not show the window, or shows it unguarded", name)
		}
	}
}
