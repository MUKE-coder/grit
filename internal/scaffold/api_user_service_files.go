package scaffold

// The user service: every read and write behind the user endpoints, and its
// tests.
//
// The handler used to run these fourteen queries itself, which is the thing the
// generator forbids for every resource it writes. Moving them means a job, a
// console command and PUT /api/users/:id apply the same rules, instead of the
// rules living on the HTTP route.
//
// Both are template files rather than Go string literals, for the reason the
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
