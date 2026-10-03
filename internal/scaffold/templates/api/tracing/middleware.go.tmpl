package tracing

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Middleware starts a span for every request, continuing one that arrived.
//
// Written here rather than taken from otelgin, for two reasons. It names the
// span after the route pattern rather than the URL, so /api/v1/widgets/:id is
// one span name and not one per id, which is the difference between a readable
// trace list and forty thousand entries. And it ties the trace id to this
// project's request id, which a general-purpose middleware cannot know about.
//
// Mounted unconditionally. With no endpoint configured the global provider is a
// no-op, the span costs an allocation, and the alternative is a conditional
// mount that gets the order wrong the first time somebody reorders the chain.
func Middleware() gin.HandlerFunc {
	tracer := otel.Tracer("http")
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		// Continue the caller's trace when it sent one. Without this each
		// service starts its own and the trace shows one hop.
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// FullPath is the route pattern; it is empty for a 404, where the raw
		// path is all there is and is also what you want to see.
		name := c.FullPath()
		if name == "" {
			name = c.Request.Method + " (no route)"
		}

		ctx, span := tracer.Start(ctx, name,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(c.Request.Method),
				semconv.URLPath(c.Request.URL.Path),
				semconv.HTTPRoute(c.FullPath()),
				semconv.UserAgentOriginal(c.Request.UserAgent()),
			),
		)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		// One id for the log line and the span.
		//
		// RequestID has already run and set a generated id. When the request is
		// sampled, the trace id replaces it, so searching the logs for the id on
		// a slow span finds the lines for that exact request. Two ids that do
		// not match are two searches and a guess.
		if traceID := TraceIDFrom(ctx); traceID != "" {
			c.Set("request_id", traceID)
			c.Header("X-Request-ID", traceID)
			c.Header("X-Trace-Id", traceID)
		}

		c.Next()

		status := c.Writer.Status()
		span.SetAttributes(semconv.HTTPResponseStatusCode(status))

		// Only 5xx marks the span as an error. A 404 or a 422 is the API
		// working: marking those red makes the error rate a number nobody can
		// act on, and a trace list where everything is red is one nobody reads.
		if status >= 500 {
			span.SetStatus(codes.Error, strconv.Itoa(status))
		}
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last())
		}

		// The authenticated caller, when there is one. Enough to answer "whose
		// requests are slow", and deliberately not the email: a trace goes to a
		// third-party collector, and an id is enough to join back to a user
		// here while being nothing on its own over there.
		if userID, ok := c.Get("user_id"); ok {
			if id, isString := userID.(string); isString && id != "" {
				span.SetAttributes(attribute.String("enduser.id", id))
			}
		}
	}
}
