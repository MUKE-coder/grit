package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// M33: one visibility rule for tickets, answering 404, and role constants
// rather than string literals.
func TestTicketVisibilityIsOneRuleAnswering404(t *testing.T) {
	src := ticketHandlerGo()
	service := ticketServiceGo()

	// Since M29 the rule lives in services.TicketService: Reply, SetStatus and
	// Assign load through Visible, and Visible and Query both apply scope.
	if n := strings.Count(service, "s.Visible(ctx, actor, id, false)"); n != 3 {
		t.Errorf("Reply, SetStatus and Assign should all load the ticket through Visible; found %d", n)
	}
	if !strings.Contains(service, "s.scope(q, actor).First(&t,") {
		t.Error("the visibility rule is checked after the query instead of being part of it")
	}
	if strings.Contains(src+service, "not your ticket") || strings.Contains(service, "ErrTicketNotYours") {
		t.Error("somebody else's ticket is still answered with 403, which confirms the id exists")
	}
	if strings.Contains(src, "h.DB.WithContext") {
		t.Error("the ticket handler runs a query of its own again")
	}
	if strings.Count(service, `"ADMIN"`) != 0 || !strings.Contains(service, "models.RoleAdmin, true") {
		t.Error("the admin notification fan-out spells ADMIN as a literal")
	}
	if n := strings.Count(src, `role == "ADMIN"`) + strings.Count(src, `role != "ADMIN"`) + strings.Count(src, `"ADMIN", true`); n != 0 {
		t.Errorf("the ticket handler still spells ADMIN as a literal %d times; models.RoleAdmin exists", n)
	}
	if !strings.Contains(src, "models.RoleAdmin || role == models.RoleEditor") {
		t.Error("ticketStaff does not use the role constants")
	}
	// M32, the ticket half: the COALESCE order only when nothing was asked for.
	if !strings.Contains(src, "if !ticketListConfig.Sortable[params.SortBy] {") {
		t.Error("the ticket list still pre-applies its order, so ?sort_by= only breaks ties")
	}
	if strings.Contains(src, `q.Where("subject LIKE ?`) {
		t.Error("the ticket search is still a case-sensitive LIKE")
	}
	if !strings.Contains(src, `Searchable:   []string{"subject", "description"}`) {
		t.Error("the ticket list does not search through paginate")
	}
	assertTemplateIsGo(t, "handlers/ticket.go", src)
}

// M33: the notification scope in one place, MarkRead included.
func TestNotificationScopeIsOneRule(t *testing.T) {
	src := notificationHandlerGo()
	if n := strings.Count(src, "scopeNotifications("); n != 5 {
		t.Errorf("the scope should be defined once and applied four times; found %d mentions", n)
	}
	if strings.Contains(src, `role == "ADMIN"`) {
		t.Error("the notification handler still spells ADMIN as a literal")
	}
	if !strings.Contains(src, "res.RowsAffected == 0") {
		t.Error("MarkRead still marks any notification read by id, whoever it belongs to")
	}
	assertTemplateIsGo(t, "handlers/notification.go", src)
}

// M34: a dashboard says when its upstream is down instead of rendering zeros.
func TestDashboardsReportADeadUpstream(t *testing.T) {
	security := securityHandlerGo()
	observability := observabilityHandlerGo()

	for name, src := range map[string]string{"security.go": security, "observability.go": observability} {
		if strings.Contains(src, "_ = h.Bridge.") {
			t.Errorf("%s still discards the error from a bridge call", name)
		}
		if !strings.Contains(src, `"degraded": up.names()`) {
			t.Errorf("%s does not report which panels are missing", name)
		}
		assertTemplateIsGo(t, name, src)
	}
	if !strings.Contains(security, `"SENTINEL_UNAVAILABLE"`) || !strings.Contains(security, "up.allFailed(3)") {
		t.Error("the security dashboard still answers 200 when Sentinel answers nothing")
	}
	if !strings.Contains(observability, `"PULSE_UNAVAILABLE"`) || !strings.Contains(observability, "up.allFailed(4)") {
		t.Error("the performance dashboard still answers 200 when Pulse answers nothing")
	}
	// The anonymous structs became named types, so the two handlers read as
	// bridge shapes rather than as forty lines of inline struct.
	for _, want := range []string{"type sentinelBlocked struct", "type sentinelStats struct", "type sentinelThreats struct"} {
		if !strings.Contains(security, want) {
			t.Errorf("security.go is missing %q", want)
		}
	}
	for _, want := range []string{"type pulseOverview struct", "type pulseRuntime struct", "type pulseN1 struct", "type pulseErrors struct"} {
		if !strings.Contains(observability, want) {
			t.Errorf("observability.go is missing %q", want)
		}
	}
}

func assertTemplateIsGo(t *testing.T, name, src string) {
	t.Helper()
	if _, err := format.Source([]byte(strings.ReplaceAll(src, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Errorf("%s is not valid Go: %v", name, err)
	}
}
