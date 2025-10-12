package simclient

import (
	"sort"
	"sync"
	"time"
)

// MetricsRecorder stores aggregated action timings across many scenarios.
type MetricsRecorder struct {
	mu      sync.Mutex
	totals  map[string]*metricAccumulator
	overall overallStats
}

type metricAccumulator struct {
	count     int
	failures  int
	successes int
	total     time.Duration
	min       time.Duration
	max       time.Duration
	samples   []time.Duration
}

type overallStats struct {
	scenarios int
	successes int
	failures  int
}

// MetricSummary exposes derived statistics for reporting.
type MetricSummary struct {
	Action    string
	Count     int
	Successes int
	Failures  int
	Min       time.Duration
	Max       time.Duration
	Avg       time.Duration
	P95       time.Duration
}

// OverallSummary provides coarse scenario-level counts.
type OverallSummary struct {
	Scenarios int
	Successes int
	Failures  int
}

// NewMetricsRecorder creates an empty recorder.
func NewMetricsRecorder() *MetricsRecorder {
	return &MetricsRecorder{
		totals: make(map[string]*metricAccumulator),
	}
}

// AddScenario captures metrics for a scenario execution.
func (r *MetricsRecorder) AddScenario(result *ScenarioResult) {
	if result == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.overall.scenarios++
	if result.Success() {
		r.overall.successes++
	} else {
		r.overall.failures++
	}

	for _, action := range result.Actions {
		acc := r.totals[action.Name]
		if acc == nil {
			acc = &metricAccumulator{}
			r.totals[action.Name] = acc
		}
		acc.count++
		if action.Err != nil {
			acc.failures++
		} else {
			acc.successes++
		}
		acc.total += action.Duration
		acc.samples = append(acc.samples, action.Duration)
		if acc.count == 1 || action.Duration < acc.min {
			acc.min = action.Duration
		}
		if action.Duration > acc.max {
			acc.max = action.Duration
		}
	}
}

// Snapshot returns the aggregated metrics summaries.
func (r *MetricsRecorder) Snapshot() ([]MetricSummary, OverallSummary) {
	r.mu.Lock()
	defer r.mu.Unlock()

	summaries := make([]MetricSummary, 0, len(r.totals))
	for action, acc := range r.totals {
		if acc.count == 0 {
			continue
		}
		summary := MetricSummary{
			Action:    action,
			Count:     acc.count,
			Successes: acc.successes,
			Failures:  acc.failures,
		}
		if len(acc.samples) > 0 {
			sort.Slice(acc.samples, func(i, j int) bool { return acc.samples[i] < acc.samples[j] })
			summary.Min = acc.min
			summary.Max = acc.max
			summary.Avg = time.Duration(0)
			divisor := acc.successes + acc.failures
			if divisor > 0 {
				summary.Avg = acc.total / time.Duration(divisor)
			}
			summary.P95 = percentile(acc.samples, 0.95)
		}
		summaries = append(summaries, summary)
	}

	sort.Slice(summaries, func(i, j int) bool { return summaries[i].Action < summaries[j].Action })

	return summaries, OverallSummary{
		Scenarios: r.overall.scenarios,
		Successes: r.overall.successes,
		Failures:  r.overall.failures,
	}
}

func percentile(data []time.Duration, p float64) time.Duration {
	if len(data) == 0 {
		return 0
	}
	if p <= 0 {
		return data[0]
	}
	if p >= 1 {
		return data[len(data)-1]
	}
	idx := int(float64(len(data)-1) * p)
	return data[idx]
}
