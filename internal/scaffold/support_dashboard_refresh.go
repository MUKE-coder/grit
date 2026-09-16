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
//     tie-breaker. Since M29 the rules live in services/ticket.go, which is
//     written with it, along with the ticket_mail.go the service calls.
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
	add("notification.go", notificationHandlerGo())
	add("security.go", securityHandlerGo())
	add("observability.go", observabilityHandlerGo())

	module := opts.Module()

	// The ticket handler calls services.TicketService (M29), which calls
	// QueueTicketCreatedEmail in ticket_mail.go (M30). All three go together,
	// mail and service first: a new handler over an old service does not
	// compile. When ticket_mail.go is somebody's own and keeps its old shape,
	// the handler is left as it is and the reason is printed.
	if handler := filepath.Join(handlers, "ticket.go"); fileExists(handler) {
		services := filepath.Join(apiRoot, "internal", "services")
		mailFile := filepath.Join(services, "ticket_mail.go")
		if err := writeFile(mailFile, strings.ReplaceAll(ticketMailGo(), "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", mailFile, err)
		}
		if fileContains(mailFile, "func QueueTicketCreatedEmail(") {
			for _, f := range []struct{ path, content string }{
				{filepath.Join(services, "ticket.go"), ticketServiceGo()},
				{handler, ticketHandlerGo()},
			} {
				if err := writeFile(f.path, strings.ReplaceAll(f.content, "{{MODULE}}", module)); err != nil {
					return fmt.Errorf("writing %s: %w", f.path, err)
				}
			}
		} else {
			fmt.Printf("  ⚠ internal/services/ticket_mail.go is yours, so the ticket handler was left alone: add QueueTicketCreatedEmail, which builds the same message and hands it to mail.Queue, and upgrade again\n")
		}
	}

	for path, content := range files {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
