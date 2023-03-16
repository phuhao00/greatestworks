package traceio

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	// Trace attribute keys for various Service Weaver identifiers. These
	// are attached to all exported traces by the weavelet, and displayed
	// in the UI by the Service Weaver visualization tools (e.g., dashboard).
	AppNameTraceKey             = attribute.Key("serviceweaver.app")
	VersionTraceKey             = attribute.Key("serviceweaver.version")
	ColocationGroupNameTraceKey = attribute.Key("serviceweaver.coloc_group")
	GroupReplicaIDTraceKey      = attribute.Key("serviceweaver.group_replica_id")
)

// TestTracer returns a simple tracer suitable for tests.
func TestTracer() trace.Tracer {
	exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	return sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter)).Tracer("test")
}
