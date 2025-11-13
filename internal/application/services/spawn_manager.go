package services

import (
	"context"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// SpawnTask is a unit of spawning work executed asynchronously.
type SpawnTask func(ctx context.Context)

// SpawnManager provides a simple asynchronous queue to run spawn tasks.
type SpawnManager struct {
	logger logging.Logger
	ch     chan SpawnTask
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewSpawnManager(logger logging.Logger, queueSize int) *SpawnManager {
	if queueSize <= 0 {
		queueSize = 1024
	}
	return &SpawnManager{
		logger: logger,
		ch:     make(chan SpawnTask, queueSize),
	}
}

func (m *SpawnManager) Start(parent context.Context, workers int) {
	if workers <= 0 {
		workers = 2
	}
	ctx, cancel := context.WithCancel(parent)
	m.cancel = cancel
	m.logger.Info("SpawnManager starting", logging.Fields{"workers": workers})
	for i := 0; i < workers; i++ {
		m.wg.Add(1)
		go func(id int) {
			defer m.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-m.ch:
					safeRun(ctx, task, m.logger)
				}
			}
		}(i)
	}
}

func (m *SpawnManager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
	close(m.ch)
	done := make(chan struct{})
	go func() { m.wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		// timeout best-effort
	}
}

func (m *SpawnManager) Enqueue(task SpawnTask) {
	select {
	case m.ch <- task:
	default:
		m.logger.Warn("Spawn queue full, dropping task", logging.Fields{})
	}
}

func safeRun(ctx context.Context, task SpawnTask, logger logging.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("spawn task panic", nil, logging.Fields{"recover": r})
		}
	}()
	task(ctx)
}
