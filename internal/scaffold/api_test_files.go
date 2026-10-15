package scaffold

func apiAuthTestGo() string { return tmpl("api/handlers/auth_test.go") }

func apiUserTestGo() string { return tmpl("api/handlers/user_test.go") }

func apiBenchTestGo() string { return tmpl("api/handlers/bench_test.go") }

// apiPerfTestGo emits internal/handlers/perf_test.go: the allocation budgets
// for the per-request hot path, and the benchmarks the CI perf job compares.
func apiPerfTestGo() string { return tmpl("api/handlers/perf_test.go") }
