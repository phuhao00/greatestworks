package call

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

const traceHeaderLen = 25

// writeTraceContext serializes the trace context (if any) contained in ctx
// into b.
// REQUIRES: len(b) >= traceHeaderLen
func writeTraceContext(ctx context.Context, b []byte) {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return
	}

	// Send trace information in the header.
	// TODO(spetrovic): Confirm that we don't need to bother with TraceState,
	// which seems to be used for storing vendor-specific information.
	traceID := sc.TraceID()
	spanID := sc.SpanID()
	copy(b, traceID[:])
	copy(b[16:], spanID[:])
	b[24] = byte(sc.TraceFlags())
}

// readTraceContext returns a span context with tracing information stored in b.
// REQUIRES: len(b) >= traceHeaderLen
func readTraceContext(b []byte) trace.SpanContext {
	cfg := trace.SpanContextConfig{
		TraceID:    *(*trace.TraceID)(b[:16]),
		SpanID:     *(*trace.SpanID)(b[16:24]),
		TraceFlags: trace.TraceFlags(b[24]),
		Remote:     true,
	}
	return trace.NewSpanContext(cfg)
}
