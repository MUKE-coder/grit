package scaffold

func apiAuthTestGo() string { return tmpl("api/handlers/auth_test.go") }

// apiAuthSessionTestGo emits internal/handlers/auth_session_test.go: Refresh,
// rotation replay and Logout, which auth_test.go does not reach (grit#163).
func apiAuthSessionTestGo() string { return tmpl("api/handlers/auth_session_test.go") }

// apiAuthMiddlewareTestGo emits internal/middleware/auth_test.go: the accept and
// reject paths of the function every protected route runs first (grit#163).
func apiAuthMiddlewareTestGo() string { return tmpl("api/middleware/auth_test.go") }

func apiUserTestGo() string { return tmpl("api/handlers/user_test.go") }

func apiBenchTestGo() string { return tmpl("api/handlers/bench_test.go") }

// apiPerfTestGo emits internal/handlers/perf_test.go: the allocation budgets
// for the per-request hot path, and the benchmarks the CI perf job compares.
func apiPerfTestGo() string { return tmpl("api/handlers/perf_test.go") }
