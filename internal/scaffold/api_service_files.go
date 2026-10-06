package scaffold

// The services behind the framework's own endpoints, and their tests.
//
// Each of these handlers used to run its own queries, which is the thing the
// generator forbids for every resource it writes. Moving them means a job, a
// console command and the HTTP route apply the same rules, instead of the rules
// living on the route.
//
// They are template files rather than Go string literals, for the reason the
// handlers moved in v3.338.0: a test can parse them, an editor can highlight
// them, and a missing brace is a unit-test failure naming the line.

// apiUserServiceGo emits internal/services/user.go.
func apiUserServiceGo() string { return tmpl("api/services/user.go") }

// apiUserServiceTestGo emits internal/services/user_test.go.
//
// It carries the two role-assignment cases that used to sit in
// handlers/user_role_sync_test.go, beside the code they are about, plus the one
// nothing tested: that Update writes the row and the assignment together or
// leaves both alone.
func apiUserServiceTestGo() string { return tmpl("api/services/user_test.go") }

// apiRoleServiceGo emits internal/services/role.go.
//
// The roles API's seventeen queries, moved in v3.341.0. The handler keeps the
// decisions about the request: who may change a built-in role, that a role hands
// out no more than its author holds, that an unknown permission key is a 400.
func apiRoleServiceGo() string { return tmpl("api/services/role.go") }

// apiRoleServiceTestGo emits internal/services/role_test.go.
func apiRoleServiceTestGo() string { return tmpl("api/services/role_test.go") }

// apiTwoFactorServiceGo emits internal/services/two_factor.go.
//
// The two-factor tables, moved out of an 870-line handler in v3.342.0. Three of
// its methods are compare-and-sets: a code, a backup code and a pending token
// each spent by one request and not another.
func apiTwoFactorServiceGo() string { return tmpl("api/services/two_factor.go") }

// apiTwoFactorServiceTestGo emits internal/services/two_factor_test.go.
func apiTwoFactorServiceTestGo() string { return tmpl("api/services/two_factor_test.go") }

// apiSSOServiceTestGo emits internal/services/sso_connections_test.go.
//
// The connection and identity store, which the SSO handler stopped querying
// itself in v3.343.0: one domain claimed twice, a disabled connection, a SAML
// assertion arriving for a connection switched to OIDC, and the subject match a
// returning sign-in depends on.
func apiSSOServiceTestGo() string { return tmpl("api/services/sso_connections_test.go") }

// apiFeatureFlagServiceGo emits internal/services/feature_flag.go.
//
// The flag rows and the exposure counts drawn from them. The engine in
// internal/flags still decides what a flag answers and when its cache is stale,
// which is the one thing a service cannot know.
func apiFeatureFlagServiceGo() string { return tmpl("api/services/feature_flag.go") }

// The bin: what a soft delete left behind, and the four things worth doing
// to it. Built on the sync registry, which already names every resource.
func apiTrashServiceGo() string { return tmpl("api/services/trash.go") }

func apiTrashHandlerGo() string { return tmpl("api/handlers/trash.go") }

// apiFeatureFlagServiceTestGo emits internal/services/feature_flag_test.go.
func apiFeatureFlagServiceTestGo() string { return tmpl("api/services/feature_flag_test.go") }

// apiWebhookEventServiceGo emits internal/services/webhook_event.go.
//
// The webhook_events table, including the claim that decides which of two
// redeliveries of one event may run its handler.
func apiWebhookEventServiceGo() string { return tmpl("api/services/webhook_event.go") }

// apiWebhookEventServiceTestGo emits internal/services/webhook_event_test.go.
func apiWebhookEventServiceTestGo() string { return tmpl("api/services/webhook_event_test.go") }

// apiFormShareServiceGo emits internal/services/form_share.go.
//
// Public form shares and their submission rows. The submission count is an
// increment in SQL, because a public form is where two writes arrive at once.
func apiFormShareServiceGo() string { return tmpl("api/services/form_share.go") }

// apiUploadServiceGo emits internal/services/upload.go.
//
// The uploads table, scoped by owner. The storage bucket stays with the handler:
// a row and an object are two systems, and only the caller knows which order to
// fail in.
func apiUploadServiceGo() string { return tmpl("api/services/upload.go") }

// apiSyncServiceGo emits internal/services/sync.go.
//
// The row store behind the offline sync protocol. It takes destinations rather
// than models, because a project syncs whatever it registered and the type is the
// registry's: what the service owns is the keyset pagination a pull depends on,
// and what a save from a phone may touch.
func apiSyncServiceGo() string { return tmpl("api/services/sync.go") }

// apiActivityReadsServiceGo emits internal/services/activity_reads.go.
//
// The read side of the two audit tables, which had a query each in three
// handlers. Writing to them is still what the LogX functions in that package do:
// a caller logging an action should not have to build a service for it.
func apiActivityReadsServiceGo() string { return tmpl("api/services/activity_reads.go") }
