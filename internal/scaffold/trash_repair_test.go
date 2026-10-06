package scaffold

import (
	"strings"
	"testing"
)

// routes.go as it is before the bin exists: the two lines the repair anchors
// on, and nothing else it cares about.
const routesBeforeTrash = `package routes

func Setup() {
	syncHandler := handlers.NewSyncHandler(db, syncRegistry)

	{
		staff.GET("/admin/observability/summary", middleware.RequireRole("ADMIN", "perm:system.view"), observabilityHandler.Summary)
	}
}
`

// The bin has to reach a project that already exists.
//
// internal/routes/routes.go is the developer's file: the generator injects
// resource routes into it and an upgrade never rewrites it, so a page whose
// routes live only in the template would 404 for every project but a new one.
func TestTrashRepairAddsTheRoutesAndTheHandler(t *testing.T) {
	out, fixed, warnings := repairTrashRoutesSource(routesBeforeTrash)
	if len(warnings) > 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(fixed) != 2 {
		t.Errorf("reported %v, want the bin and the closed accounts", fixed)
	}

	for _, want := range []string{
		"trashHandler := handlers.NewTrashHandler(db, syncRegistry)",
		`staff.GET("/admin/trash"`,
		`staff.POST("/admin/trash/:table/:id/restore"`,
		`staff.DELETE("/admin/trash/:table/:id"`,
		`staff.GET("/admin/deleted-accounts"`,
		`staff.POST("/admin/deleted-accounts/:id/restore"`,
		`staff.DELETE("/admin/deleted-accounts/:id"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repair did not add %s", want)
		}
	}

	// Reading the bin is a system view; restoring or purging is ADMIN and
	// nothing less. A resource's delete permission says you may remove a row,
	// not that you may undo somebody else's removal.
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "/admin/trash/:table/:id") || strings.Contains(line, "deleted-accounts/:id") {
			if strings.Contains(line, "perm:system.view") {
				t.Errorf("a restore or purge is reachable with a view permission: %s", strings.TrimSpace(line))
			}
		}
	}
}

// A second upgrade changes nothing, because a repair that keeps repairing grows
// the file every time.
func TestTrashRepairIsIdempotent(t *testing.T) {
	once, _, _ := repairTrashRoutesSource(routesBeforeTrash)
	twice, fixed, _ := repairTrashRoutesSource(once)
	if twice != once {
		t.Error("a second run changed the file again")
	}
	if len(fixed) > 0 {
		t.Errorf("a second run reported %v", fixed)
	}
}

// Each block is checked for separately.
//
// A single "has this repair run" guard is true the moment the first half lands,
// and everything added afterwards reaches new projects only. That is exactly how
// v3.374.0 shipped a fix that skipped every project already running, so the
// half-applied file is the case worth a test.
func TestTrashRepairFinishesAHalfAppliedFile(t *testing.T) {
	// Only the bin, as an earlier build of this repair would have left it.
	half, _, _ := repairTrashRoutesSource(routesBeforeTrash)
	half = strings.Replace(half, deletedAccountRouteLines, "", 1)
	if strings.Contains(half, "deleted-accounts") {
		t.Fatal("the fixture still has the closed-account routes")
	}

	out, fixed, warnings := repairTrashRoutesSource(half)
	if len(warnings) > 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if !strings.Contains(out, `staff.GET("/admin/deleted-accounts"`) {
		t.Error("the closed-account routes never reached a file that already had the bin")
	}
	if strings.Count(out, "trashHandler := handlers.NewTrashHandler") != 1 {
		t.Error("the handler was added a second time")
	}
	if len(fixed) != 1 {
		t.Errorf("reported %v, want only the closed accounts", fixed)
	}
}

// A file Grit did not write is left alone, with a warning that says what to do.
func TestTrashRepairLeavesAnUnfamiliarFileAlone(t *testing.T) {
	src := "package routes\n\nfunc Setup() {}\n"
	out, fixed, warnings := repairTrashRoutesSource(src)
	if out != src {
		t.Error("the repair edited a file it does not recognise")
	}
	if len(fixed) > 0 {
		t.Errorf("reported %v", fixed)
	}
	if len(warnings) == 0 {
		t.Error("it edited nothing and said nothing, so the page will 404 with no explanation")
	}
}
