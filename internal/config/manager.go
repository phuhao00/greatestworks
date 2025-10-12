package config

import (
	"context"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatcherFunc receives an updated configuration snapshot.
type WatcherFunc func(*Config)

// Manager caches configuration and optionally hot-reloads when files change.
type Manager struct {
	loader           *Loader
	mu               sync.RWMutex
	cfg              *Config
	watchers         []WatcherFunc
	fsWatcher        *fsnotify.Watcher
	watchedFiles     map[string]struct{}
	debounceInterval time.Duration
}

// ManagerOption configures a Manager instance.
type ManagerOption func(*Manager)

// WithDebounce sets the debounce interval applied before reloading after file events.
func WithDebounce(interval time.Duration) ManagerOption {
	return func(m *Manager) {
		if interval > 0 {
			m.debounceInterval = interval
		}
	}
}

// NewManager constructs a Manager using the provided loader and options.
func NewManager(loader *Loader, opts ...ManagerOption) (*Manager, error) {
	if loader == nil {
		loader = NewLoader()
	}

	cfg, files, err := loader.Load()
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		loader:           loader,
		cfg:              cfg,
		watchers:         make([]WatcherFunc, 0),
		watchedFiles:     make(map[string]struct{}, len(files)),
		debounceInterval: 250 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(manager)
	}

	for _, file := range files {
		manager.watchedFiles[file] = struct{}{}
	}

	return manager, nil
}

// Config returns a clone of the current configuration for safe concurrent use.
func (m *Manager) Config() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg.Clone()
}

// OnChange registers a watcher callback and immediately invokes it asynchronously with the current config.
func (m *Manager) OnChange(callback WatcherFunc) {
	if callback == nil {
		return
	}

	m.mu.Lock()
	m.watchers = append(m.watchers, callback)
	snapshot := m.cfg.Clone()
	m.mu.Unlock()

	go callback(snapshot)
}

// Reload forces the manager to reload configuration files and notify watchers on success.
func (m *Manager) Reload() error {
	cfg, files, err := m.loader.Load()
	if err != nil {
		return err
	}

	newSet := make(map[string]struct{}, len(files))
	for _, file := range files {
		newSet[file] = struct{}{}
	}

	m.mu.Lock()
	oldSet := m.watchedFiles
	m.cfg = cfg
	m.watchedFiles = newSet
	watchers := append([]WatcherFunc(nil), m.watchers...)
	watcher := m.fsWatcher
	m.mu.Unlock()

	if watcher != nil {
		for file := range oldSet {
			if _, ok := newSet[file]; !ok {
				_ = watcher.Remove(file)
			}
		}
		for file := range newSet {
			if _, ok := oldSet[file]; !ok {
				_ = watcher.Add(file)
			}
		}
	}

	for _, cb := range watchers {
		if cb != nil {
			cb(cfg.Clone())
		}
	}

	return nil
}

// StartWatching begins watching configuration files for changes until the context is cancelled.
func (m *Manager) StartWatching(ctx context.Context) error {
	m.mu.Lock()
	if m.fsWatcher != nil {
		m.mu.Unlock()
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		m.mu.Unlock()
		return err
	}

	for file := range m.watchedFiles {
		_ = watcher.Add(file)
	}

	m.fsWatcher = watcher
	debounce := m.debounceInterval
	m.mu.Unlock()

	go m.watchLoop(ctx, watcher, debounce)
	return nil
}

// Close stops file watching if it is active.
func (m *Manager) Close() error {
	m.mu.Lock()
	watcher := m.fsWatcher
	m.fsWatcher = nil
	m.mu.Unlock()

	if watcher != nil {
		return watcher.Close()
	}
	return nil
}

func (m *Manager) watchLoop(ctx context.Context, watcher *fsnotify.Watcher, debounce time.Duration) {
	defer watcher.Close()

	var (
		pending bool
		timer   *time.Timer
	)

	if debounce > 0 {
		timer = time.NewTimer(debounce)
		if !timer.Stop() {
			<-timer.C
		}
	}

	triggerReload := func() {
		_ = m.Reload()
	}

	for {
		var timerCh <-chan time.Time
		if timer != nil {
			timerCh = timer.C
		}

		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}

			if debounce <= 0 {
				triggerReload()
				continue
			}

			if !pending {
				pending = true
				timer.Reset(debounce)
				continue
			}

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(debounce)
		case <-timerCh:
			pending = false
			triggerReload()
		case <-watcher.Errors:
			// Errors are ignored; loader reload will surface issues when necessary.
		}
	}
}
