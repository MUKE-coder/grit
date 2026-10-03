package scaffold

import (
	"path/filepath"
	"strings"
)

// internal/saga: a multi-step process that has to either finish or be undone.
//
// internal/workflow turns a status column into a state machine, which is the
// right tool when the states are a record's own and the transitions are a
// person clicking a button. It is no help at all when the steps are in four
// systems: charge a card, reserve stock, book a courier, send the receipt. The
// card does not roll back when the courier refuses, and the process holding all
// of that in its head is exactly the thing that crashes.
//
// So this is the other half: each step has a Do and an Undo, where the run got
// to is a row rather than a stack frame, and a failure past its retries undoes
// the completed steps newest first. A process that dies between the charge and
// the reservation resumes rather than leaving a customer charged for nothing.
//
// The models are separate from the engine because AutoMigrate reads
// internal/models and nothing else, and the engine imports them rather than the
// other way round.

// apiSagaGo emits internal/saga/saga.go.
func apiSagaGo() string { return tmpl("api/saga/saga.go") }

// apiSagaTestGo emits internal/saga/saga_test.go.
func apiSagaTestGo() string { return tmpl("api/saga/saga_test.go") }

// apiSagaModelsGo emits internal/models/saga_run.go, which holds both tables.
func apiSagaModelsGo() string { return tmpl("api/models/saga_run.go") }

// sagaModelRegistration adds the two tables to the AutoMigrate list.
//
// Registered like every other table so the backup writer includes them: a saga
// missing from a backup loses the record of what was half-done, which is the
// one thing you cannot reconstruct.
const sagaModelRegistration = `		// Sagas: the multi-step processes that have to finish or be undone, and
		// the step-by-step record of each run. Here like every other table, so
		// AutoMigrate creates them and a backup includes them.
		&SagaRun{},
		&SagaStep{},
`

// sagaModelAnchor is the line the registration goes before.
const sagaModelAnchor = "\t\t// grit:models\n"

// apiSagasRegistryGo emits internal/sagas/sagas.go, the project's own package
// of sagas. Empty until `grit generate workflow` writes one.
func apiSagasRegistryGo() string { return tmpl("api/sagas/sagas.go") }

// sagaRuntimeBlock registers the sagas and starts the runner, beside the event
// bus relay for the same reason: both are background delivery that every
// replica runs and that row claims keep from doubling up.
//
// The runner is not started when nothing is registered. A project with no sagas
// would otherwise poll an empty table every five seconds for the life of the
// process, which is a cost with no benefit and a line in the slow query log
// that somebody has to rule out.
const sagaRuntimeBlock = `
	// Sagas: the multi-step processes that have to finish or be undone. Every
	// replica runs a runner; a run is claimed before it is touched, so two of
	// them cannot advance the same one and charge a card twice.
	sagas.Register()
	if len(saga.Registered()) > 0 {
		(&saga.Runner{DB: db}).Start(context.Background())
		log.Printf("Sagas: %d registered, runner started", len(saga.Registered()))
	}
`

// writeSagaRegistry creates internal/sagas/sagas.go if it is not there.
//
// Created once and never rewritten. It holds the developer's sagas and the
// generator injects into it, so an upgrade that replaced it would delete every
// saga in the project. The engine beside it travels on upgrade; this does not,
// and that split is the same one models/user.go is on.
func writeSagaRegistry(root string, opts Options) (bool, error) {
	path := filepath.Join(opts.APIRoot(root), "internal", "sagas", "sagas.go")
	return createIfMissing(path, strings.ReplaceAll(apiSagasRegistryGo(), "{{MODULE}}", opts.Module()))
}

// apiSagaServiceGo emits internal/services/saga.go, the read side of the two
// tables plus the one write an operator needs.
func apiSagaServiceGo() string { return tmpl("api/services/saga.go") }

// apiSagaHandlerGo emits internal/handlers/saga.go.
func apiSagaHandlerGo() string { return tmpl("api/handlers/saga.go") }

// apiSagaServiceTestGo emits internal/services/saga_test.go.
func apiSagaServiceTestGo() string { return tmpl("api/services/saga_test.go") }

// sagaHandlerBlock builds the handler, beside the other system handlers.
const sagaHandlerBlock = `	// The saga runs screen. A run that goes stuck needs a person, and without a
	// screen the only way to find one is SQL against saga_runs.
	sagaHandler := &handlers.SagaHandler{DB: db}
`

// sagaRoutesBlock mounts the screen's endpoints, under the same guard as the
// jobs screen: these rows carry what a saga was started with, which for a
// checkout is an order id and for an onboarding is somebody's email.
const sagaRoutesBlock = `		staff.GET("/admin/sagas", middleware.RequireRole("ADMIN", "perm:jobs.view"), sagaHandler.List)
		staff.GET("/admin/sagas/definitions", middleware.RequireRole("ADMIN", "perm:jobs.view"), sagaHandler.Definitions)
		staff.GET("/admin/sagas/:id", middleware.RequireRole("ADMIN", "perm:jobs.view"), sagaHandler.GetByID)
		staff.POST("/admin/sagas/:id/retry", middleware.RequireRole("ADMIN", "perm:jobs.edit"), sagaHandler.Retry)
`
