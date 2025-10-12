package simclient

import (
	"context"
	"os"
	"testing"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

func TestBasicScenarioSmoke(t *testing.T) {
	if os.Getenv("SIMCLIENT_E2E") != "1" {
		t.Skip("set SIMCLIENT_E2E=1 to run simulator integration smoke test")
	}

	cfg := DefaultConfig()
	cfg.Auth.Enabled = false
	cfg.Scenario.Duration = NewDuration(3 * time.Second)
	cfg.Scenario.ActionInterval = NewDuration(500 * time.Millisecond)

	logger := logging.NewBaseLogger(logging.DebugLevel)
	runner, err := NewRunner(cfg, logger)
	if err != nil {
		t.Fatalf("failed to build scenario: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := runner.RunOnce(ctx)
	if err != nil {
		t.Fatalf("scenario execution returned error: %v", err)
	}
	if result == nil {
		t.Fatalf("scenario returned nil result")
	}
	if !result.Success() {
		t.Fatalf("scenario actions reported errors: %+v", result.Errors())
	}
}
