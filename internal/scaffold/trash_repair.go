package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The bin, for a project that already exists.
//
// internal/routes/routes.go is the developer's file: the generator injects
// resource routes into it and an upgrade never rewrites it. So the handler and
// the five routes that serve /admin/trash reach a new project through the
// template and reach every other project through this.
//
// Both edits anchor on lines Grit wrote and nothing else touches: the sync
// handler, which the bin's registry comes from, and the observability route,
// which sits in the same staff group the bin belongs to.

const (
	trashHandlerAnchor = "\tsyncHandler := handlers.NewSyncHandler(db, syncRegistry)\n"
	trashHandlerLine   = "\t// The bin reads the same registry: every resource that can be synced is a\n" +
		"\t// resource whose rows soft-delete, which is the same list.\n" +
		"\ttrashHandler := handlers.NewTrashHandler(db, syncRegistry)\n"

	trashRouteAnchor = "\t\tstaff.GET(\"/admin/observability/summary\", middleware.RequireRole(\"ADMIN\", \"perm:system.view\"), observabilityHandler.Summary)\n"
	trashRouteLines  = "\n\t\t// The bin. Reading it is a system view; restoring somebody else's\n" +
		"\t\t// delete, or making one permanent, is ADMIN and nothing less: a\n" +
		"\t\t// resource's own delete permission says you may remove a row, not that\n" +
		"\t\t// you may undo another person's removal or put it beyond recovery.\n" +
		"\t\tstaff.GET(\"/admin/trash\", middleware.RequireRole(\"ADMIN\", \"perm:system.view\"), trashHandler.Buckets)\n" +
		"\t\tstaff.GET(\"/admin/trash/:table\", middleware.RequireRole(\"ADMIN\", \"perm:system.view\"), trashHandler.List)\n" +
		"\t\tstaff.POST(\"/admin/trash/:table/:id/restore\", middleware.RequireRole(\"ADMIN\"), trashHandler.Restore)\n" +
		"\t\tstaff.DELETE(\"/admin/trash/:table/:id\", middleware.RequireRole(\"ADMIN\"), trashHandler.Purge)\n" +
		"\t\tstaff.DELETE(\"/admin/trash/:table\", middleware.RequireRole(\"ADMIN\"), trashHandler.Empty)\n"

	// The last trash route, which the closed-account block is anchored on.
	trashEmptyRouteLine = "\t\tstaff.DELETE(\"/admin/trash/:table\", middleware.RequireRole(\"ADMIN\"), trashHandler.Empty)\n"

	deletedAccountRouteLines = "\n\t\t// Closed accounts, which are a soft delete the bin cannot show: users\n" +
		"\t\t// are deliberately not in the sync registry it reads.\n" +
		"\t\tstaff.GET(\"/admin/deleted-accounts\", middleware.RequireRole(\"ADMIN\"), userHandler.DeletedAccounts)\n" +
		"\t\tstaff.POST(\"/admin/deleted-accounts/:id/restore\", middleware.RequireRole(\"ADMIN\"), userHandler.RestoreAccount)\n" +
		"\t\tstaff.DELETE(\"/admin/deleted-accounts/:id\", middleware.RequireRole(\"ADMIN\"), userHandler.PurgeAccount)\n"
)

// repairTrashRoutes wires /admin/trash into a routes.go that predates it.
func repairTrashRoutes(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	// The handler has to exist first. It is written by
	// writeFrameworkOwnedFiles, which runs earlier in the same upgrade, but a
	// shape without an API gets neither.
	if !fileExists(filepath.Join(apiRoot, "internal", "handlers", "trash.go")) {
		return nil
	}
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, routes, repairTrashRoutesSource)
}

// repairTrashRoutesSource adds the handler and the two route blocks.
//
// Each block is checked for separately. A single "has this repair run" guard
// would be true the moment the first half landed, and everything added to it
// afterwards would reach new projects only: that is precisely how v3.374.0
// shipped a fix that skipped every project already running.
func repairTrashRoutesSource(src string) (string, []string, []string) {
	out := src
	var fixed []string
	var warnings []string

	if !strings.Contains(out, "trashHandler :=") {
		// Without the sync handler there is no registry to read, and the bin
		// would list nothing while looking like it worked.
		if !strings.Contains(out, trashHandlerAnchor) {
			return src, nil, []string{"has no sync handler to put the bin beside, so /system/trash will not load: see the v3.378.0 changelog for the routes to add"}
		}
		out = strings.Replace(out, trashHandlerAnchor, trashHandlerAnchor+trashHandlerLine, 1)
	}

	if !strings.Contains(out, `staff.GET("/admin/trash"`) {
		if !strings.Contains(out, trashRouteAnchor) {
			warnings = append(warnings, "has no observability route to put the bin's routes beside, so /system/trash will not load: see the v3.378.0 changelog for the routes to add")
		} else {
			out = strings.Replace(out, trashRouteAnchor, trashRouteAnchor+trashRouteLines, 1)
			fixed = append(fixed, "deleted records are listed, restorable and purgeable at /system/trash")
		}
	}

	if !strings.Contains(out, `staff.GET("/admin/deleted-accounts"`) {
		anchor := trashEmptyRouteLine
		if !strings.Contains(out, anchor) {
			warnings = append(warnings, "has no trash routes to put the closed-account routes beside, so /system/deleted-accounts will not load")
		} else {
			out = strings.Replace(out, anchor, anchor+deletedAccountRouteLines, 1)
			fixed = append(fixed, "closed accounts are listed, restorable and purgeable at /system/deleted-accounts")
		}
	}

	if out == src {
		return src, nil, warnings
	}
	return out, fixed, warnings
}
