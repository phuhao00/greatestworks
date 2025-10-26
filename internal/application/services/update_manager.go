package services

import (
	"context"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// Updatable represents any object that can be updated on each tick.
type Updatable interface {
	Update(ctx context.Context, delta time.Duration) error
}

// UpdateFunc is a functional adapter to register plain functions.
type UpdateFunc func(ctx context.Context, delta time.Duration) error

func (f UpdateFunc) Update(ctx context.Context, delta time.Duration) error { return f(ctx, delta) }

// UpdateManager runs a game-loop style tick and calls registered updatables.
type UpdateManager struct {
	logger   logging.Logger
	interval time.Duration
	mu       sync.RWMutex
	running  bool
	cancel   context.CancelFunc
	updaters map[string]Updatable
}

// NewUpdateManager creates an update manager with a default interval (50ms ~ 20 TPS) when interval<=0.
func NewUpdateManager(logger logging.Logger, interval time.Duration) *UpdateManager {
	if interval <= 0 {
		interval = 50 * time.Millisecond
	}
	return &UpdateManager{
		logger:   logger,
		interval: interval,
		updaters: make(map[string]Updatable),
	}
}

// Register adds an updatable by id (last write wins).
func (m *UpdateManager) Register(id string, u Updatable) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updaters[id] = u
}

// Unregister removes an updatable.
func (m *UpdateManager) Unregister(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.updaters, id)
}

// Start begins the tick loop in a goroutine; calling Start again has no effect.
func (m *UpdateManager) Start(parent context.Context) {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(parent)
	m.cancel = cancel
	m.running = true
	m.mu.Unlock()

	go m.loop(ctx)
}

// Stop cancels the loop and waits for it to end.
func (m *UpdateManager) Stop() {
	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		return
	}
	cancel := m.cancel
	m.running = false
	m.cancel = nil
	m.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (m *UpdateManager) loop(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	last := time.Now()

	m.logger.Info("UpdateManager loop started", logging.Fields{"interval_ms": m.interval.Milliseconds()})
	for {
		select {
		case <-ctx.Done():
			m.logger.Info("UpdateManager loop stopped")
			return
		case <-ticker.C:
			now := time.Now()
			delta := now.Sub(last)
			last = now
			m.tick(ctx, delta)
		}
	}
}

func (m *UpdateManager) tick(ctx context.Context, delta time.Duration) {
	m.mu.RLock()
	// copy to avoid holding lock during updates
	updaters := make([]Updatable, 0, len(m.updaters))
	for _, u := range m.updaters {
		updaters = append(updaters, u)
	}
	m.mu.RUnlock()

	for _, u := range updaters {
		// Best-effort: isolate panics per updater
		func() {
			defer func() {
				if r := recover(); r != nil {
					m.logger.Error("updater panic", nil, logging.Fields{"recover": r})
				}
			}()
			if err := u.Update(ctx, delta); err != nil {
				m.logger.Warn("updater error", logging.Fields{"error": err.Error()})
			}
		}()
	}
}
