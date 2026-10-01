package scaffold

// apiSAMLTestGo emits internal/handlers/saml_test.go.
//
// These cover attribute extraction, which is where SAML integrations actually
// fail: every provider names its claims differently, so "your identity provider
// did not release an email address" is the error an operator hits when the
// fallbacks are wrong. Pinned names, the Okta-style short names and the Entra
// ID claim URIs are all exercised here.
func apiSAMLTestGo() string { return tmpl("api/handlers/saml_test.go") }
