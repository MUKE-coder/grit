package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Wiring the saga engine into a project that predates it.
//
// writeFrameworkOwnedFiles delivers internal/saga and the two models, and
// writeSagaRegistry creates internal/sagas. What does not arrive is the two
// lines in files the developer owns: the tables in the AutoMigrate list, and
// the registration and runner in routes.go.
//
// Without those an upgraded project has an engine, no tables and nothing
// advancing a run, which is the shape of a half-delivered feature: it compiles,
// the migration reports nothing to do, and the first saga anybody writes sits in
// a table that does not exist. That has happened here before, to media, to
// recovery contacts and to passkeys, which is why this file exists.

// repairSagaWiring adds the tables and the runner to an existing project.
func repairSagaWiring(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	// The engine has to be there first: these two edits reference it, and a
	// repair that writes a call to a package the project has not got is a
	// project that does not compile.
	if !fileContains(filepath.Join(apiRoot, "internal", "saga", "saga.go"), "func Register(d Definition)") {
		return nil
	}
	if !fileContains(filepath.Join(apiRoot, "internal", "sagas", "sagas.go"), "// grit:sagas") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	models := filepath.Join(apiRoot, "internal", "models", "user.go")
	if fileExists(models) {
		if err := repairSourceFile(root, m, models, repairSagaModelsSource); err != nil {
			return err
		}
	}
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	return repairSourceFile(root, m, routes, repairSagaRuntimeSource(opts.Module()))
}

// repairSagaModelsSource adds the two tables to the AutoMigrate list.
func repairSagaModelsSource(src string) (string, []string, []string) {
	if strings.Contains(src, "&SagaRun{}") {
		return src, nil, nil
	}
	const marker = "\t\t// grit:models\n"
	if strings.Count(src, marker) != 1 {
		return src, nil, []string{"the model list has no // grit:models marker, so the saga tables were not registered: " +
			"add &SagaRun{} and &SagaStep{} to it, or a saga will write to a table AutoMigrate never created"}
	}
	out := strings.Replace(src, marker, sagaModelRegistration+marker, 1)
	return out, []string{"the saga tables are in the model list, so migrate creates them and a backup includes them"}, nil
}

// repairSagaRuntimeSource registers the sagas and starts the runner.
func repairSagaRuntimeSource(module string) func(string) (string, []string, []string) {
	return func(src string) (string, []string, []string) {
		if strings.Contains(src, "sagas.Register()") {
			return src, nil, nil
		}
		// Anchored on the event bus relay, which is the line this belongs beside:
		// both are background delivery that every replica runs and that row
		// claims keep from doubling up.
		const anchor = "\tevents.StartRelay(db)\n"
		if strings.Count(src, anchor) != 1 {
			return src, nil, []string{"routes.go does not start the event bus relay the way Grit wrote it, so the saga runner was not started: " +
				"call sagas.Register() and start a saga.Runner, or a saga will be registered and never advance"}
		}
		out := strings.Replace(src, anchor, anchor+sagaRuntimeBlock, 1)
		var ok bool
		for _, path := range []string{module + "/internal/saga", module + "/internal/sagas"} {
			if out, ok = addImportGroup(out, path); !ok {
				return src, nil, []string{"could not add the saga imports to routes.go, so the runner was not started"}
			}
		}
		fixed := []string{"sagas are registered and a runner advances them, so a generated workflow actually runs"}

		// The screen's endpoints, in the same pass. A project with the handler
		// file and no routes has an admin page that answers 404, which is a
		// worse state than not having the page: it looks broken rather than
		// absent.
		const handlerAnchor = "\timportJobHandler := &handlers.ImportJobHandler{DB: db}\n"
		const routeAnchor = "\t\tstaff.DELETE(\"/admin/jobs/queue/:queue\", middleware.RequireRole(\"ADMIN\", \"perm:jobs.edit\"), jobsHandler.ClearQueue)\n"
		switch {
		case strings.Contains(out, "sagaHandler :="):
			// Already there.
		case strings.Count(out, handlerAnchor) != 1 || strings.Count(out, routeAnchor) != 1:
			fixed = append(fixed, "the saga runs screen has no endpoints: routes.go is not where Grit left it, so mount sagaHandler by hand or /system/sagas will answer 404")
		default:
			out = strings.Replace(out, handlerAnchor, handlerAnchor+sagaHandlerBlock, 1)
			out = strings.Replace(out, routeAnchor, routeAnchor+sagaRoutesBlock, 1)
			fixed = append(fixed, "the saga runs screen can read them at /api/v1/admin/sagas")
		}

		// The routes screen, on the same two anchors. It is a different feature
		// and it is mounted in the same two places, so doing it in a second pass
		// would mean matching text this one has already moved.
		if !strings.Contains(out, "routeExplorerHandler :=") && strings.Contains(out, sagaHandlerBlock) {
			out = strings.Replace(out, sagaHandlerBlock, sagaHandlerBlock+routeExplorerHandlerBlock, 1)
			out = strings.Replace(out, sagaRoutesBlock, sagaRoutesBlock+routeExplorerRoutesBlock, 1)
			fixed = append(fixed, "the routes screen can read the router's table at /api/v1/admin/routes")
		}
		return out, fixed, nil
	}
}
