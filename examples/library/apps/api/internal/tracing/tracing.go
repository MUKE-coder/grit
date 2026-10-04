// Package tracing exports OpenTelemetry traces, when you point it somewhere.
//
// # Why this is here, and why it is off by default
//
// This project already answers "what is my server doing" better than most: Pulse
// profiles every request, the activity log records who did what, and every
// response carries an X-Request-ID that ties a log line to a request. All of
// that stops at the edge of this process.
//
// A trace does not. When this API calls a payments service, which calls a
// ledger, a trace is the one artifact that shows the whole thing as a single
// timeline with the slow span highlighted. You cannot reconstruct that from
// three sets of logs with three different request ids, which is what everybody
// tries first.
//
// It is off until OTEL_EXPORTER_OTLP_ENDPOINT is set, because a tracer with
// nowhere to send spans is a background goroutine and a growing buffer that
// pays for nothing.
//
// # No new configuration
//
// The variables are OpenTelemetry's own, not Grit's:
//
//	OTEL_EXPORTER_OTLP_ENDPOINT   http://localhost:4318   turns it on
//	OTEL_SERVICE_NAME             defaults to APP_NAME
//	OTEL_TRACES_SAMPLER_ARG       0.0 to 1.0, default 1.0 in dev, 0.1 otherwise
//
// Anyone who has run a collector before already knows them, and anyone who has
// not can paste a line from a vendor's quickstart and have it work. A
// GRIT_TRACING_URL would have been one more thing to look up.
//
// # What it does to a request id
//
// When a trace is sampled, the request id becomes the trace id. That is the
// whole point of having both: a log line and a span that share an id are one
// search away from each other, and two ids that do not match are two searches
// and a guess. X-Request-ID still goes out, so nothing that reads it breaks.
package tracing

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"library/apps/api/internal/config"
)

// EndpointEnv is the variable that turns tracing on.
const EndpointEnv = "OTEL_EXPORTER_OTLP_ENDPOINT"

// Enabled reports whether an endpoint is configured.
func Enabled() bool {
	return strings.TrimSpace(os.Getenv(EndpointEnv)) != ""
}

// Setup starts the exporter and returns the function that drains it.
//
// Always returns a usable shutdown, including when tracing is off and when
// setup failed, so a caller never has to decide whether it has one. A server
// that exits without draining loses the spans for the request that was being
// served when somebody looked at the trace, which is the request they cared
// about.
func Setup(ctx context.Context, cfg *config.Config) (shutdown func(context.Context) error, err error) {
	noop := func(context.Context) error { return nil }

	if !Enabled() {
		return noop, nil
	}

	endpoint := strings.TrimSpace(os.Getenv(EndpointEnv))
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		// Reported, not fatal. A collector that is down must not stop the API
		// from serving: observability is how you find out something is wrong,
		// and it refusing to start is not something being wrong.
		return noop, fmt.Errorf("tracing to %s: %w", endpoint, err)
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName(cfg)),
		attribute.String("deployment.environment", cfg.AppEnv),
	))
	if err != nil {
		res = resource.Default()
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler(cfg)),
	)

	otel.SetTracerProvider(provider)
	// W3C tracecontext, plus baggage. Without a propagator every service starts
	// its own trace and the distributed part of distributed tracing is lost,
	// which is the single most common way a working collector still shows
	// nothing useful.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		return provider.Shutdown(ctx)
	}, nil
}

// serviceName is what the trace is labelled with in the collector.
func serviceName(cfg *config.Config) string {
	if name := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME")); name != "" {
		return name
	}
	if cfg != nil && cfg.AppName != "" {
		return cfg.AppName
	}
	return "grit-api"
}

// sampler decides how much is kept.
//
// Everything in development, because a developer looking for one request wants
// that request and not a tenth of it. A tenth in production, because a trace
// per request at a thousand requests a second is a bill rather than a tool.
// OTEL_TRACES_SAMPLER_ARG overrides both.
func sampler(cfg *config.Config) sdktrace.Sampler {
	ratio := 0.1
	if cfg != nil && cfg.IsDevelopment() {
		ratio = 1.0
	}
	if raw := strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER_ARG")); raw != "" {
		if parsed, err := strconv.ParseFloat(raw, 64); err == nil && parsed >= 0 && parsed <= 1 {
			ratio = parsed
		}
	}
	// ParentBased so a request that arrives already sampled stays sampled.
	// Sampling each service independently produces traces with holes in them,
	// which read as a service that did not respond.
	return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))
}

// Tracer is how application code starts a span of its own.
//
//	ctx, span := tracing.Tracer("billing").Start(ctx, "charge card")
//	defer span.End()
//
// Safe to call when tracing is off: the global provider is a no-op then, and
// the span costs an allocation and nothing else.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// TraceIDFrom is the current trace id, or "" when nothing is being traced.
//
// Used to make the request id and the trace id the same string, which is the
// difference between one search and two.
func TraceIDFrom(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() || !sc.TraceID().IsValid() {
		return ""
	}
	return sc.TraceID().String()
}
