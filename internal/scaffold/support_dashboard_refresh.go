package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// refreshSupportAndDashboards brings four framework-owned handlers forward.
//
// These are not files anybody edits to configure something: there is no
// resource definition behind them and nothing is generated into them. They are
// Grit's own screens, and they carried Grit's own defects.
//
//   - handlers/ticket.go wrote its visibility rule out four times, answered
//     somebody else's ticket with 403 rather than 404, spelled ADMIN and EDITOR
//     as string literals, and pre-applied an ORDER BY that made ?sort_by= a
//     tie-breaker.
//   - handlers/notification.go wrote its scope out three times and left it out
//     of MarkRead, where any signed-in account could clear any notification.
//   - handlers/security.go and handlers/observability.go discarded the error
//     from every call to Sentinel and Pulse, so a dead upstream rendered as a
//     200 full of zeros.
//
// writeFile is manifest-guarded: a copy somebody has edited is reported as a
// conflict with a diff, not overwritten. Each file is written only where it is
// already present, so a project without the support desk or without the
// dashboards does not suddenly acquire one.
func refreshSupportAndDashboards(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	handlers := filepath.Join(apiRoot, "internal", "handlers")

	files := map[string]string{}
	add := func(name, content string) {
		path := filepath.Join(handlers, name)
		if fileExists(path) {
			files[path] = content
		}
	}
	add("ticket.go", ticketHandlerGo())
	add("notification.go", notificationHandlerGo())
	add("security.go", securityHandlerGo())
	add("observability.go", observabilityHandlerGo())

	module := opts.Module()
	for path, content := range files {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
