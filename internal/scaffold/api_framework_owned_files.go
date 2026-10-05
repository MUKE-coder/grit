package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeFrameworkOwnedFiles writes the code the framework owns.
//
// These are not your files. Nothing generates them, nothing in a resource
// definition describes them, and they only make sense as a set: the webhook
// model, the receiver that fills it and the dispatch package all read the same
// columns by name.
//
// They live in their own function because they have to travel on upgrade. A
// framework fix that changes a constraint or a field type is useless if half
// the cluster stays at the version the project was scaffolded with,
// and the failure is quiet: the code compiles, the migration reports nothing
// to do, and the constraint everything assumes is simply absent. That is
// exactly how webhook deduplication spent several versions being decoration.
//
// Your own models, in the same package, are never touched by this. So is
// models/user.go, which holds the registry people add to. The API key model
// is not here either: it travels with its seeder in writeMigrateSeedFiles,
// because the two break each other when they drift apart.
//
// writeFile is manifest-guarded, so a file someone has edited is reported as a
// conflict rather than overwritten.
func writeFrameworkOwnedFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		// The webhook cluster. The model, the receiver that builds it and the
		// dispatch package all read the same columns, so upgrading one without
		// the others gives a project that does not compile: the receiver
		// assigning a string to a *string is how this was found.
		filepath.Join(apiRoot, "internal", "models", "webhook_event.go"): apiWebhookEventModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "webhooks.go"):    apiWebhooksHandlerGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "webhooks.go"):    apiWebhooksGo(),
		// Public form sharing: the handler, its service and the two models, as one
		// set. Written once at scaffold time until v3.345.0, so the submission
		// count that lost a submission whenever two visitors posted together would
		// have been fixed for new projects only.
		filepath.Join(apiRoot, "internal", "models", "form_share.go"):      formShareModelGo(),
		filepath.Join(apiRoot, "internal", "models", "form_submission.go"): formSubmissionModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "form_share.go"):    formShareHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "form_share.go"):    apiFormShareServiceGo(),

		// The audit tables' read side: the request log, the semantic feed and the
		// OCSF export a collector polls, which had a query each in three handlers.
		// Writing to them stays with the LogX functions.
		filepath.Join(apiRoot, "internal", "services", "activity_reads.go"): apiActivityReadsServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "activity.go"):       apiActivityHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "user_activity.go"):  userActivityHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "ocsf.go"):           apiOCSFHandlerGo(),

		// The auth endpoints that read or write one account: the provider callback,
		// the password reset, the verification resend and the recovery contacts.
		// Here since v3.347.0, when they stopped running those queries themselves,
		// because a fix in any of them reached new projects only.
		filepath.Join(apiRoot, "internal", "handlers", "auth_oauth.go"):              apiAuthOAuthGo(),
		filepath.Join(apiRoot, "internal", "handlers", "auth_password_reset.go"):     apiAuthPasswordResetGo(),
		filepath.Join(apiRoot, "internal", "handlers", "auth_email_verification.go"): apiAuthEmailVerificationGo(),
		filepath.Join(apiRoot, "internal", "handlers", "recovery.go"):                recoveryHandlerGo(),

		// Failed-login counting and the admin unlock. Here since v3.346.0, when the
		// counters became AuthService's: an existing project that kept this file
		// would keep its own copy of the three writes the second factor shares.
		filepath.Join(apiRoot, "internal", "handlers", "auth_lockout.go"): apiAuthLockoutGo(),

		// The offline sync protocol's row store, which its handler stopped
		// querying directly in v3.346.0.
		filepath.Join(apiRoot, "internal", "services", "sync.go"): apiSyncServiceGo(),

		// The four states /api/health answers in, and the registry a plugin adds
		// its own probe to. Here rather than in the API map because routes.go is
		// repaired in place on upgrade, and a repair that writes health.OK into a
		// project without the package is a project that does not compile.
		filepath.Join(apiRoot, "internal", "health", "health.go"):      apiHealthGo(),
		filepath.Join(apiRoot, "internal", "health", "health_test.go"): apiHealthTestGo(),

		// llms.txt and llms-full.txt, which an agent pointed at a running service
		// reads before it calls anything. Here for the reason health is: routes.go
		// is repaired in place, and a repair that mounts a package the project has
		// not got is a project that does not compile.
		filepath.Join(apiRoot, "internal", "llms", "llms.go"):      apiLLMSGo(),
		filepath.Join(apiRoot, "internal", "llms", "llms_test.go"): apiLLMSTestGo(),

		// What each route demands of a caller. The table beside this is
		// generated from routes.go after the files are written, because the
		// running process cannot work it out: gin's RouteInfo carries a route's
		// last handler and nothing about the middleware in front of it.
		filepath.Join(apiRoot, "internal", "access", "access.go"): apiAccessGo(),

		// The shared test setup and assertions. Framework-owned so an upgrade
		// delivers a new assertion to a project that already has tests, and so a
		// generated resource's test can rely on it existing.
		// Policies: rules about a record and a user together, which neither
		// the permission catalogue nor the owner column can express.
		// OpenTelemetry. Inert until OTEL_EXPORTER_OTLP_ENDPOINT is set, and
		// framework-owned so an existing project gets it from an upgrade.
		filepath.Join(apiRoot, "internal", "tracing", "tracing.go"):    apiTracingGo(),
		filepath.Join(apiRoot, "internal", "tracing", "middleware.go"): apiTracingMiddlewareGo(),

		// The SSE fallback for platforms that do not pass a WebSocket upgrade
		// through. Framework-owned so an upgrade delivers it.
		filepath.Join(apiRoot, "internal", "handlers", "realtime_sse.go"): apiRealtimeSSEGo(),

		filepath.Join(apiRoot, "internal", "authz", "gate.go"):      apiGateGo(),
		filepath.Join(apiRoot, "internal", "authz", "gate_test.go"): apiGateTestGo(),

		filepath.Join(apiRoot, "internal", "testkit", "testkit.go"): apiTestKitGo(),
		filepath.Join(apiRoot, "internal", "testkit", "http.go"):    apiTestKitHTTPGo(),

		filepath.Join(apiRoot, "internal", "services", "webhook_event.go"):            apiWebhookEventServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "webhook_event_test.go"):       apiWebhookEventServiceTestGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "dedup_test.go"):               apiWebhookDedupTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "webhooks_redelivery_test.go"): apiWebhookRedeliveryTestGo(),
		filepath.Join(apiRoot, "internal", "services", "oauth_hooks.go"):              apiOAuthHooksGo(),
		filepath.Join(apiRoot, "internal", "services", "oauth_hooks_test.go"):         apiOAuthHooksTestGo(),

		filepath.Join(apiRoot, "internal", "models", "outbox_message.go"): apiOutboxModelGo(),

		// Sagas: the engine, its tests, and the two tables. A set like the
		// webhook cluster, and here rather than in the API map for the reason
		// the outbox model is: a fix to the engine that reached new projects
		// only would be no fix at all for the one app that has money moving
		// through it.
		filepath.Join(apiRoot, "internal", "saga", "saga.go"):                apiSagaGo(),
		filepath.Join(apiRoot, "internal", "saga", "saga_test.go"):           apiSagaTestGo(),
		filepath.Join(apiRoot, "internal", "models", "saga_run.go"):          apiSagaModelsGo(),
		filepath.Join(apiRoot, "internal", "services", "saga.go"):            apiSagaServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "saga_test.go"):       apiSagaServiceTestGo(),
		filepath.Join(apiRoot, "internal", "handlers", "saga.go"):            apiSagaHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "routes_explorer.go"): apiRouteExplorerHandlerGo(),

		// The migration history, and the rollback computed from it. Here rather
		// than with cmd/migrate because an existing project needs the package
		// before its rewritten cmd/migrate can import it: split across the two
		// sets, grit upgrade would deliver a main.go importing a package that
		// never arrived.
		filepath.Join(apiRoot, "internal", "migrate", "history.go"):      apiMigrateHistoryGo(),
		filepath.Join(apiRoot, "internal", "migrate", "history_test.go"): apiMigrateHistoryTestGo(),

		// The actor on the context, which generated services scope owned rows
		// by. A new file beside authz.go, so a project whose authz.go was
		// edited still gets it.
		filepath.Join(apiRoot, "internal", "authz", "actor.go"):      authzActorGo(),
		filepath.Join(apiRoot, "internal", "authz", "actor_test.go"): authzActorTestGo(),

		// Feature flags: the rules model and the engine that reads it, a
		// set like the webhook cluster. In the API map they were written once
		// and never again, so no flag fix reached an existing project.
		filepath.Join(apiRoot, "internal", "models", "feature_flag.go"): apiFeatureFlagModelGo(),
		filepath.Join(apiRoot, "internal", "flags", "flags.go"):         apiFlagsGo(),
		filepath.Join(apiRoot, "internal", "flags", "flags_test.go"):    apiFlagsTestGo(),
		// The flag and access-review endpoints, whose handlers joined this set in
		// v3.344.0 for the reason user.go did: a fix in one of them used to reach
		// new projects only, and the access-review list ran four counts per
		// campaign with every error dropped.
		filepath.Join(apiRoot, "internal", "handlers", "flags.go"):              apiFlagsHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "access_review.go"):      apiAccessReviewHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "access_review.go"):      apiAccessReviewServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "access_review_test.go"): apiAccessReviewTestGo(),
		filepath.Join(apiRoot, "internal", "services", "feature_flag.go"):       apiFeatureFlagServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "feature_flag_test.go"):  apiFeatureFlagServiceTestGo(),

		// The one insert of a user row, which registration and the admin's
		// Create User both call. Here because the handlers that call it are
		// repaired in place on upgrade, and a repair that calls a function the
		// project has not got is a project that does not compile.
		filepath.Join(apiRoot, "internal", "services", "user_write.go"):      apiUserWriteServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "user_write_test.go"): apiUserWriteServiceTestGo(),

		// The user endpoints: the handler, the service behind it and the
		// service's tests, as one set.
		//
		// The handler moved here from the scaffold-only map in v3.340.0, which is
		// the release that gave it a service. Before that, a fix to it reached new
		// projects and no existing one, so the hand-rolled list that read none of
		// the admin's filters had to be repaired into place with text
		// substitutions, three separate times. writeFile is manifest-guarded: a
		// user.go somebody has edited is reported as a conflict and kept, and the
		// repairs below are still there for it.
		filepath.Join(apiRoot, "internal", "handlers", "user.go"):      apiUserHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "user.go"):      apiUserServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "user_test.go"): apiUserServiceTestGo(),

		// The roles API, the same way, since v3.341.0: the handler, its test and
		// the service its queries moved to, which have to arrive together.
		filepath.Join(apiRoot, "internal", "handlers", "role.go"):      roleHandlerGo(),
		filepath.Join(apiRoot, "internal", "handlers", "role_test.go"): roleHandlerTestGo(),
		filepath.Join(apiRoot, "internal", "services", "role.go"):      apiRoleServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "role_test.go"): apiRoleServiceTestGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
