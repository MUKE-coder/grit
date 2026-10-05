package scaffold

// apiRealtimeSSEGo emits internal/handlers/realtime_sse.go: the Server-Sent
// Events fallback for the realtime hub.
//
// Some platforms do not pass a WebSocket upgrade through to the application.
// Laravel Cloud strips Connection: upgrade before the container (issue #92),
// and it is not alone: a corporate proxy, an older load balancer or an API
// gateway configured for plain HTTP all do the same. The handshake fails, the
// client's reconnect loop retries forever, and live updates silently never
// arrive.
//
// SSE is plain HTTP with a response that never ends, so it survives all of it.
// The hub needs no changes: Client.Conn was already allowed to be nil and
// delivery already goes through Send.
func apiRealtimeSSEGo() string { return tmpl("api/handlers/realtime_sse.go") }
