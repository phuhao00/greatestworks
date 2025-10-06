// Package config 配置热重载
// Author: MMO Server Team
// Created: 2024

package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ConfigChangeEvent 配置变更事件
type ConfigChangeEvent struct {
	Type      ConfigChangeType `json:"type"`
	FilePath  string           `json:"file_path"`
	OldConfig *Config          `json:"old_config,omitempty"`
	NewConfig *Config          `json:"new_config,omitempty"`
	Error     error            `json:"error,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

// ConfigChangeType 配置变更类型
type ConfigChangeType string

const (
	// ConfigCreated 配置文件创建
	ConfigCreated ConfigChangeType = "created"
	// ConfigModified 配置文件修改
	ConfigModified ConfigChangeType = "modified"
	// ConfigDeleted 配置文件删除
	ConfigDeleted ConfigChangeType = "deleted"
	// ConfigRenamed 配置文件重命名
	ConfigRenamed ConfigChangeType = "renamed"
	// ConfigError 配置错误
	ConfigError ConfigChangeType = "error"
)

// ConfigChangeHandler 配置变更处理器
type ConfigChangeHandler func(event ConfigChangeEvent)

// HotReloader 配置热重载器
type HotReloader struct {
	configLoader  *ConfigLoader
	watcher       *fsnotify.Watcher
	handlers      []ConfigChangeHandler
	configPaths   []string
	currentConfig *Config
	mutex         sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	debounceTime  time.Duration
	lastReload    time.Time
	reloadChan    chan string
	errorChan     chan error
	logger        Logger
}

// Logger 日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// DefaultLogger 默认日志实现
type DefaultLogger struct{}

func (dl *DefaultLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (dl *DefaultLogger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg, args...)
}

func (dl *DefaultLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

func (dl *DefaultLogger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] "+msg, args...)
}

// NewHotReloader 创建热重载器
func NewHotReloader(configLoader *ConfigLoader) (*HotReloader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &HotReloader{
		configLoader: configLoader,
		watcher:      watcher,
		handlers:     make([]ConfigChangeHandler, 0),
		configPaths:  make([]string, 0),
		ctx:          ctx,
		cancel:       cancel,
		debounceTime: 1 * time.Second, // 默认防抖时间
		reloadChan:   make(chan string, 10),
		errorChan:    make(chan error, 10),
		logger:       &DefaultLogger{},
	}, nil
}

// SetLogger 设置日志器
func (hr *HotReloader) SetLogger(logger Logger) {
	hr.logger = logger
}

// SetDebounceTime 设置防抖时间
func (hr *HotReloader) SetDebounceTime(duration time.Duration) {
	hr.debounceTime = duration
}

// AddConfigPath 添加配置文件路径监听
func (hr *HotReloader) AddConfigPath(configPath string) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 检查路径是否已存在
	for _, path := range hr.configPaths {
		if path == absPath {
			return nil // 已存在，不重复添加
		}
	}

	// 添加到监听列表
	if err := hr.watcher.Add(absPath); err != nil {
		return fmt.Errorf("failed to watch config path %s: %w", absPath, err)
	}

	hr.configPaths = append(hr.configPaths, absPath)
	hr.logger.Info("Added config path to hot reload: %s", absPath)

	return nil
}

// AddConfigDir 添加配置目录监听
func (hr *HotReloader) AddConfigDir(configDir string) error {
	absDir, err := filepath.Abs(configDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute directory path: %w", err)
	}

	// 检查目录是否存在
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		return fmt.Errorf("config directory does not exist: %s", absDir)
	}

	// 监听目录
	if err := hr.watcher.Add(absDir); err != nil {
		return fmt.Errorf("failed to watch config directory %s: %w", absDir, err)
	}

	hr.logger.Info("Added config directory to hot reload: %s", absDir)

	return nil
}

// AddHandler 添加配置变更处理器
func (hr *HotReloader) AddHandler(handler ConfigChangeHandler) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()
	hr.handlers = append(hr.handlers, handler)
}

// RemoveHandler 移除配置变更处理器
func (hr *HotReloader) RemoveHandler(handler ConfigChangeHandler) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	for i, h := range hr.handlers {
		// 比较函数指针（简单实现）
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			hr.handlers = append(hr.handlers[:i], hr.handlers[i+1:]...)
			break
		}
	}
}

// Start 启动热重载
func (hr *HotReloader) Start() error {
	// 加载初始配置
	if err := hr.loadInitialConfig(); err != nil {
		return fmt.Errorf("failed to load initial config: %w", err)
	}

	// 启动文件监听协程
	go hr.watchFiles()

	// 启动重载处理协程
	go hr.handleReloads()

	hr.logger.Info("Hot reloader started")
	return nil
}

// Stop 停止热重载
func (hr *HotReloader) Stop() error {
	hr.cancel()

	if err := hr.watcher.Close(); err != nil {
		hr.logger.Error("Failed to close file watcher: %v", err)
		return err
	}

	close(hr.reloadChan)
	close(hr.errorChan)

	hr.logger.Info("Hot reloader stopped")
	return nil
}

// GetCurrentConfig 获取当前配置
func (hr *HotReloader) GetCurrentConfig() *Config {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()
	return hr.currentConfig
}

// ReloadConfig 手动重载配置
func (hr *HotReloader) ReloadConfig() error {
	select {
	case hr.reloadChan <- "manual":
		return nil
	default:
		return fmt.Errorf("reload channel is full")
	}
}

// loadInitialConfig 加载初始配置
func (hr *HotReloader) loadInitialConfig() error {
	config, err := hr.configLoader.Load()
	if err != nil {
		return err
	}

	hr.mutex.Lock()
	hr.currentConfig = config
	hr.mutex.Unlock()

	return nil
}

// watchFiles 监听文件变化
func (hr *HotReloader) watchFiles() {
	for {
		select {
		case <-hr.ctx.Done():
			return

		case event, ok := <-hr.watcher.Events:
			if !ok {
				return
			}

			hr.handleFileEvent(event)

		case err, ok := <-hr.watcher.Errors:
			if !ok {
				return
			}

			hr.logger.Error("File watcher error: %v", err)
			select {
			case hr.errorChan <- err:
			default:
				hr.logger.Warn("Error channel is full, dropping error: %v", err)
			}
		}
	}
}

// handleFileEvent 处理文件事件
func (hr *HotReloader) handleFileEvent(event fsnotify.Event) {
	// 检查是否为配置文件
	if !hr.isConfigFile(event.Name) {
		return
	}

	hr.logger.Debug("File event: %s %s", event.Op.String(), event.Name)

	// 防抖处理
	if time.Since(hr.lastReload) < hr.debounceTime {
		hr.logger.Debug("Debouncing file event: %s", event.Name)
		return
	}

	// 根据事件类型处理
	switch {
	case event.Op&fsnotify.Write == fsnotify.Write:
		hr.triggerReload(event.Name, ConfigModified)
	case event.Op&fsnotify.Create == fsnotify.Create:
		hr.triggerReload(event.Name, ConfigCreated)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		hr.triggerReload(event.Name, ConfigDeleted)
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		hr.triggerReload(event.Name, ConfigRenamed)
	}
}

// isConfigFile 检查是否为配置文件
func (hr *HotReloader) isConfigFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	validExts := []string{".yaml", ".yml", ".json", ".toml"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}

	return false
}

// triggerReload 触发重载
func (hr *HotReloader) triggerReload(filePath string, changeType ConfigChangeType) {
	select {
	case hr.reloadChan <- filePath:
		hr.lastReload = time.Now()
	default:
		hr.logger.Warn("Reload channel is full, dropping reload request for: %s", filePath)
	}
}

// handleReloads 处理重载请求
func (hr *HotReloader) handleReloads() {
	for {
		select {
		case <-hr.ctx.Done():
			return

		case filePath := <-hr.reloadChan:
			hr.performReload(filePath)

		case err := <-hr.errorChan:
			hr.notifyHandlers(ConfigChangeEvent{
				Type:      ConfigError,
				Error:     err,
				Timestamp: time.Now(),
			})
		}
	}
}

// performReload 执行重载
func (hr *HotReloader) performReload(filePath string) {
	hr.logger.Info("Reloading configuration from: %s", filePath)

	// 获取当前配置的副本
	hr.mutex.RLock()
	oldConfig := hr.currentConfig
	hr.mutex.RUnlock()

	// 加载新配置
	newConfig, err := hr.configLoader.Load()
	if err != nil {
		hr.logger.Error("Failed to reload config: %v", err)
		hr.notifyHandlers(ConfigChangeEvent{
			Type:      ConfigError,
			FilePath:  filePath,
			OldConfig: oldConfig,
			Error:     err,
			Timestamp: time.Now(),
		})
		return
	}

	// 更新当前配置
	hr.mutex.Lock()
	hr.currentConfig = newConfig
	hr.mutex.Unlock()

	// 通知处理器
	hr.notifyHandlers(ConfigChangeEvent{
		Type:      ConfigModified,
		FilePath:  filePath,
		OldConfig: oldConfig,
		NewConfig: newConfig,
		Timestamp: time.Now(),
	})

	hr.logger.Info("Configuration reloaded successfully")
}

// notifyHandlers 通知所有处理器
func (hr *HotReloader) notifyHandlers(event ConfigChangeEvent) {
	hr.mutex.RLock()
	handlers := make([]ConfigChangeHandler, len(hr.handlers))
	copy(handlers, hr.handlers)
	hr.mutex.RUnlock()

	for _, handler := range handlers {
		go func(h ConfigChangeHandler) {
			defer func() {
				if r := recover(); r != nil {
					hr.logger.Error("Config change handler panicked: %v", r)
				}
			}()
			h(event)
		}(handler)
	}
}

// ConfigWatcher 配置监听器（简化接口）
type ConfigWatcher struct {
	reloader *HotReloader
}

// Watch 监听配置文件
func (cw *ConfigWatcher) Watch(configPath string, handler ConfigChangeHandler) error {
	if err := cw.reloader.AddConfigPath(configPath); err != nil {
		return err
	}

	cw.reloader.AddHandler(handler)

	return cw.reloader.Start()
}

// Stop 停止监听
func (cw *ConfigWatcher) Stop() error {
	return cw.reloader.Stop()
}

// GetConfig 获取当前配置
func (cw *ConfigWatcher) GetConfig() *Config {
	return cw.reloader.GetCurrentConfig()
}
