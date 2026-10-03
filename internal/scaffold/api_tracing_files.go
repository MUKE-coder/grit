package scaffold

// apiTracingGo and apiTracingMiddlewareGo emit internal/tracing: OpenTelemetry
// spans, exported when OTEL_EXPORTER_OTLP_ENDPOINT is set and inert otherwise.
//
// Pulse profiles every request, the activity log records who did what, and
// X-Request-ID ties a log line to a request. All of that stops at the edge of
// this process. A trace does not: when this API calls a payments service which
// calls a ledger, a trace is the one artifact that shows the whole thing as a
// single timeline, and that cannot be reconstructed from three sets of logs
// with three different request ids.
func apiTracingGo() string { return tmpl("api/tracing/tracing.go") }

// apiTracingMiddlewareGo is the per-request span, which also makes the request
// id and the trace id the same string.
func apiTracingMiddlewareGo() string { return tmpl("api/tracing/middleware.go") }
