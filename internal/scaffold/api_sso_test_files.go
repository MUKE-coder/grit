package scaffold

// apiSSOTestGo emits internal/handlers/sso_test.go.
//
// These cover the part of SSO that a browser test can't reach: how a login
// resolves to an account, and how IdP groups become roles. They exist because
// the first version of this code had two real bugs that only a test caught —
// a `default:true` bool meant an operator turning OFF just-in-time
// provisioning silently got it ON, and the same for a disabled connection.
func apiSSOTestGo() string { return tmpl("api/handlers/sso_test.go") }
