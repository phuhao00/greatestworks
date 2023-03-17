package codegen

import (
	"testing"
	"time"
)

func BenchmarkMetrics(b *testing.B) {
	metrics := MethodMetricsFor(MethodLabels{
		Caller:    "caller",
		Component: "component",
		Method:    "method",
	})

	b.Run("Everything", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// NOTE(mwhittaker): We don't include metrics.Error.Add(1) because
			// the metric is updated infrequently.
			start := time.Now()
			metrics.Count.Add(1)
			metrics.Latency.Put(float64(time.Since(start).Microseconds()))
			metrics.BytesRequest.Put(100)
			metrics.BytesReply.Put(100)
		}
	})

	b.Run("Time", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			start := time.Now()
			time.Since(start).Microseconds()
		}
	})

	b.Run("Counter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metrics.Count.Add(1)
		}
	})

	b.Run("Latency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metrics.Latency.Put(100)
		}
	})

	b.Run("Bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metrics.Latency.Put(100)
		}
	})
}
