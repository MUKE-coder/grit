package scaffold

func apiAuthTestGo() string { return tmpl("api/handlers/auth_test.go") }

func apiUserTestGo() string { return tmpl("api/handlers/user_test.go") }

func apiBenchTestGo() string { return tmpl("api/handlers/bench_test.go") }
