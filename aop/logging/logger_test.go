package logging

import (
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestTestLogger(t *testing.T) {
	// Test plan: Launch a goroutine that continues to write to a TestLogger
	// after the test ends. The logger should stop logging when the test ends.
	t.Run("sub", func(t *testing.T) {
		logger := NewTestLogger(t)
		go func() {
			for {
				logger.Debug("Ping")
			}
		}()
		// Give the logger a chance to log something.
		time.Sleep(1 * time.Millisecond)
	})
	// Allow the goroutine to keep running, even though the test has finished.
	time.Sleep(50 * time.Millisecond)
}

func BenchmarkMakeEntry(b *testing.B) {
	opt := Options{
		App:        "app",
		Deployment: "dep",
		Component:  "comp",
		Weavelet:   "wlet",
	}
	for i := 0; i < b.N; i++ {
		e := makeEntry("info", "test", nil, 0, opt)
		proto.Marshal(e)
	}
}
