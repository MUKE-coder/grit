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
