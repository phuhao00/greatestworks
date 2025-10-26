package simclient

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// Runner orchestrates scenarios for integration or load testing modes.
type Runner struct {
	cfg      Config
	logger   logging.Logger
	scenario Scenario
}

// NewRunner constructs a runner using the basic scenario by default.
func NewRunner(cfg Config, logger logging.Logger) (*Runner, error) {
	cfg.Normalize()
	scenarioLogger := logger.WithField("component", "scenario")

	scenarioCfg := cfg.Scenario
	scenario, err := buildScenario(&scenarioCfg, scenarioLogger)
	if err != nil {
		return nil, err
	}

	cfg.Scenario = scenarioCfg

	return &Runner{
		cfg:      cfg,
		logger:   logger,
		scenario: scenario,
	}, nil
}

// Config returns a copy of the runner configuration.
func (r *Runner) Config() Config {
	return r.cfg
}

// RunOnce executes the scenario a single time for integration testing.
func (r *Runner) RunOnce(ctx context.Context) (*ScenarioResult, error) {
	client := NewSimulatorClient(1, &r.cfg, r.logger.WithField("mode", "integration"))
	result, err := r.scenario.Execute(ctx, client)
	return result, err
}

// LoadReport summarises the outcome of a load test run.
type LoadReport struct {
	Scenario    string
	StartedAt   time.Time
	CompletedAt time.Time
	Users       int
	Concurrency int
	Iterations  int
	Metrics     []MetricSummary
	Overall     OverallSummary
	Errors      []string
}

// RunLoad executes multiple scenario instances concurrently.
func (r *Runner) RunLoad(ctx context.Context) (*LoadReport, error) {
	if !r.cfg.Load.Enabled {
		return nil, fmt.Errorf("load testing disabled in configuration")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	recorder := NewMetricsRecorder()
	report := &LoadReport{
		Scenario:    r.scenario.Name(),
		StartedAt:   time.Now(),
		Users:       r.cfg.Load.VirtualUsers,
		Concurrency: r.cfg.Load.Concurrency,
		Iterations:  r.cfg.Load.Iterations,
	}

	semaphore := make(chan struct{}, r.cfg.Load.Concurrency)
	var wg sync.WaitGroup
	var errOnce atomic.Bool
	var errorsMu sync.Mutex

outer:
	for user := 0; user < r.cfg.Load.VirtualUsers; user++ {
		select {
		case <-ctx.Done():
			r.logger.Warn("load test cancelled before completion", logging.Fields{"error": ctx.Err()})
			break outer
		default:
		}

		semaphore <- struct{}{}
		wg.Add(1)

		go func(userIndex int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			clientLogger := r.logger.WithFields(logging.Fields{
				"mode":       "load",
				"user_index": userIndex,
			})
			client := NewSimulatorClient(userIndex+1, &r.cfg, clientLogger)

			for iteration := 0; iteration < r.cfg.Load.Iterations; iteration++ {
				if ctx.Err() != nil {
					return
				}
				result, err := r.scenario.Execute(ctx, client)
				recorder.AddScenario(result)

				if err != nil {
					errorsMu.Lock()
					if len(report.Errors) < 50 {
						report.Errors = append(report.Errors, fmt.Sprintf("user %d iteration %d: %v", userIndex+1, iteration+1, err))
					}
					errorsMu.Unlock()

					if r.cfg.Load.StopOnError && !errOnce.Load() {
						errOnce.Store(true)
						cancel()
					}
					if r.cfg.Load.StopOnError {
						return
					}
				}
			}
		}(user)

		if ramp := r.cfg.Load.RampUp.AsDuration(); ramp > 0 {
			perUserDelay := ramp / time.Duration(r.cfg.Load.VirtualUsers)
			if perUserDelay > 0 {
				select {
				case <-ctx.Done():
					break outer
				case <-time.After(perUserDelay):
				}
			}
		}
	}

	wg.Wait()
	report.CompletedAt = time.Now()
	report.Metrics, report.Overall = recorder.Snapshot()
	return report, nil
}

func buildScenario(cfg *ScenarioConfig, logger logging.Logger) (Scenario, error) {
	scenarioType := strings.TrimSpace(strings.ToLower(cfg.Type))
	hasCustomActions := len(cfg.Actions) > 0 || len(cfg.Features) > 0

	// E2E 场景优先
	if scenarioType == "e2e" {
		cfg.Type = "e2e"
		return NewE2EScenario(*cfg, logger), nil
	}

	if (scenarioType == "" || scenarioType == "basic") && !hasCustomActions {
		cfg.Type = "basic"
		return NewBasicScenario(*cfg, logger), nil
	}

	if scenarioType == "" || scenarioType == "basic" {
		cfg.Type = "feature"
	}

	if len(cfg.Features) == 0 {
		if _, ok := featureLibrary[scenarioType]; ok {
			cfg.Features = append(cfg.Features, scenarioType)
		}
	}

	scenario, err := NewActionScenario(*cfg, logger)
	if err != nil {
		return nil, err
	}

	return scenario, nil
}
