package scaffold

// The auth handler, split by concern.
//
// auth.go passed a thousand lines covering sessions, password reset, email
// verification, OAuth and account lockout at once, and it is a file people
// customise. Each flow now has its own file, so changing one means opening one.
// auth.go keeps the session core: register, login, refresh, logout, me.

// apiAuthPasswordResetGo emits handlers/auth_password_reset.go.
//
// Carved out of auth.go, which was a thousand lines holding five unrelated
// flows. Changing the reset email or the token lifetime meant scrolling past
// OAuth and session refresh to find them.
func apiAuthPasswordResetGo() string { return tmpl("api/handlers/auth_password_reset.go") }

// apiAuthEmailVerificationGo emits handlers/auth_email_verification.go.
func apiAuthEmailVerificationGo() string { return tmpl("api/handlers/auth_email_verification.go") }

// apiAuthOAuthGo emits handlers/auth_oauth.go. The provider registry lives in
// services; these are the two endpoints a browser actually visits.
func apiAuthOAuthGo() string { return tmpl("api/handlers/auth_oauth.go") }

// apiAuthLockoutGo emits handlers/auth_lockout.go: the failed-login counter
// and the admin unlock that clears it.
func apiAuthLockoutGo() string { return tmpl("api/handlers/auth_lockout.go") }
