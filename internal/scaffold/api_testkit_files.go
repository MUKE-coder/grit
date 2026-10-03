package scaffold

// apiTestKitGo and apiTestKitHTTPGo emit internal/testkit: the shared setup and
// assertions every handler test in a generated project needs.
//
// Before this, each test file opened its own SQLite, wired its own config, and
// unmarshalled the body into a map[string]any before it could look at anything.
// That is about fifteen lines before the first assertion, and the fifteenth copy
// is the one that gets a detail wrong. It also meant every generated test ran
// against SQLite with no way to point the same suite at Postgres, which is where
// the timestamp precision and JSONB differences live.
func apiTestKitGo() string { return tmpl("api/testkit/testkit.go") }

// apiTestKitHTTPGo is the request client and the response assertions.
func apiTestKitHTTPGo() string { return tmpl("api/testkit/http.go") }
