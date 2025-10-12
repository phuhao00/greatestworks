package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"greatestworks/internal/infrastructure/logging"
	simclient "greatestworks/tools/simclient"
)

func main() {
	var (
		configPath  string
		mode        string
		users       int
		concurrency int
		iterations  int
		enableAuth  bool
		disableAuth bool
		debug       bool
	)

	flag.StringVar(&configPath, "config", "", "Path to simulator YAML configuration")
	flag.StringVar(&mode, "mode", "integration", "Mode: integration or load")
	flag.IntVar(&users, "users", 0, "Override virtual user count for load mode")
	flag.IntVar(&concurrency, "concurrency", 0, "Override concurrency for load mode")
	flag.IntVar(&iterations, "iterations", 0, "Override scenario iterations per user in load mode")
	flag.BoolVar(&enableAuth, "auth", false, "Force enable auth flow")
	flag.BoolVar(&disableAuth, "no-auth", false, "Force disable auth flow")
	flag.BoolVar(&debug, "debug", false, "Enable verbose debug logging")
	flag.Parse()

	cfg, err := simclient.LoadConfigFromFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if enableAuth {
		cfg.Auth.Enabled = true
	}
	if disableAuth {
		cfg.Auth.Enabled = false
	}
	if strings.EqualFold(mode, "load") && !cfg.Load.Enabled {
		cfg.Load.Enabled = true
	}
	if users > 0 {
		cfg.Load.VirtualUsers = users
	}
	if concurrency > 0 {
		cfg.Load.Concurrency = concurrency
	}
	if iterations > 0 {
		cfg.Load.Iterations = iterations
	}

	level := logging.InfoLevel
	if debug {
		level = logging.DebugLevel
	}
	logger := logging.NewBaseLogger(level)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	runner, err := simclient.NewRunner(cfg, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build scenario: %v\n", err)
		os.Exit(1)
	}

	switch strings.ToLower(mode) {
	case "integration", "single", "once":
		result, err := runner.RunOnce(ctx)
		printScenarioResult(result, err)
		if err != nil {
			os.Exit(1)
		}
	case "load", "stress":
		report, err := runner.RunLoad(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load test failed: %v\n", err)
			os.Exit(1)
		}
		printLoadReport(report)
		if report.Overall.Failures > 0 || len(report.Errors) > 0 {
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q\n", mode)
		os.Exit(2)
	}
}

func printScenarioResult(result *simclient.ScenarioResult, execErr error) {
	if result == nil {
		fmt.Println("scenario did not produce a result")
		if execErr != nil {
			fmt.Printf("error: %v\n", execErr)
		}
		return
	}

	duration := result.CompletedAt.Sub(result.StartedAt)
	fmt.Printf("Scenario: %s\n", result.ScenarioName)
	fmt.Printf("Duration: %s\n", duration)
	fmt.Printf("Success: %t\n", result.Success() && execErr == nil)
	if execErr != nil {
		fmt.Printf("Error: %v\n", execErr)
	}

	fmt.Println("Actions:")
	for _, action := range result.Actions {
		status := "OK"
		if action.Err != nil {
			status = "ERR"
		}
		fmt.Printf("  %-26s %8s  %s", action.Name, action.Duration, status)
		if action.Err != nil {
			fmt.Printf(" (%v)", action.Err)
		}
		fmt.Println()
	}
}

func printLoadReport(report *simclient.LoadReport) {
	if report == nil {
		fmt.Println("load test produced no report")
		return
	}

	totalDuration := report.CompletedAt.Sub(report.StartedAt)
	fmt.Printf("Load Scenario: %s\n", report.Scenario)
	fmt.Printf("Users: %d  Concurrency: %d  Iterations/User: %d\n", report.Users, report.Concurrency, report.Iterations)
	fmt.Printf("Total Duration: %s\n", totalDuration)
	fmt.Printf("Scenarios: %d (success: %d, failures: %d)\n", report.Overall.Scenarios, report.Overall.Successes, report.Overall.Failures)

	fmt.Println("Action Metrics:")
	for _, metric := range report.Metrics {
		fmt.Printf("  %-24s count=%4d success=%4d fail=%3d", metric.Action, metric.Count, metric.Successes, metric.Failures)
		if metric.Count > 0 {
			fmt.Printf("  min=%8s avg=%8s p95=%8s max=%8s", metric.Min, metric.Avg, metric.P95, metric.Max)
		}
		fmt.Println()
	}

	if len(report.Errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range report.Errors {
			fmt.Printf("  %s\n", err)
		}
	}

	fmt.Printf("Completed at %s\n", report.CompletedAt.Format(time.RFC3339))
}
