package scaffold

// roleHandlerTestGo emits internal/handlers/role_test.go.
//
// Shipped with the app because these pin behaviour teams depend on and can
// easily break while editing the catalog:
//   - grants are STORED unexpanded (so wildcards keep inheriting) but SERVED
//     expanded (so the UI never reimplements matching)
//   - unknown permission keys are rejected, instead of stored and silently
//     never matching
//   - built-in roles cannot be renamed or deleted via the API, not merely
//     greyed out in the UI
//   - assigning roles keeps the legacy users.role string in step, or routes
//     still guarded by role name start returning spurious 403s
func roleHandlerTestGo() string { return tmpl("api/handlers/role_test.go") }
