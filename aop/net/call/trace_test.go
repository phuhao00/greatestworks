package call

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

func TestTraceSerialization(t *testing.T) {
	// Create a random trace context.
	rndBytes := func() []byte {
		b := uuid.New()
		return b[:]
	}
	span := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID(uuid.New()),
		SpanID:     *(*trace.SpanID)(rndBytes()[:8]),
		TraceFlags: trace.TraceFlags(rndBytes()[0]),
	})

	// Serialize the trace context.
	var b [25]byte
	writeTraceContext(
		trace.ContextWithSpanContext(context.Background(), span), b[:])

	// Deserialize the trace context.
	actual := readTraceContext(b[:])
	expect := span.WithRemote(true)
	if !expect.Equal(actual) {
		want, _ := json.Marshal(expect)
		got, _ := json.Marshal(actual)
		t.Errorf("span context diff, want %q, got %q", want, got)
	}
}
