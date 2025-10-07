package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"

	"github.com/fsnotify/fsnotify"
)

// FileWatcher 文件监听�?
type FileWatcher struct {
	filePath    string
	callback    func(string)
	logger      logger.Logger
	watcher     *fsnotify.Watcher
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	isRunning   bool
	lastModTime time.Time
	debounce    time.Duration
}

// WatcherConfig 监听器配�?
type WatcherConfig struct {
	DebounceInterval time.Duration `json:"debounce_interval" yaml:"debounce_interval"`
	WatchDirectory   bool          `json:"watch_directory" yaml:"watch_directory"`
	Recursive        bool          `json:"recursive" yaml:"recursive"`
	IgnorePatterns   []string      `json:"ignore_patterns" yaml:"ignore_patterns"`
}

// Watcher 文件监听器接�?
type Watcher interface {
	// Start 启动监听
	Start() error

	// Stop 停止监听
	Stop() error

	// IsRunning 检查是否正在运�?
	IsRunning() bool

	// AddPath 添加监听路径
	AddPath(path string) error

	// RemovePath 移除监听路径
	RemovePath(path string) error

	// SetCallback 设置回调函数
	SetCallback(callback func(string))
}

// NewFileWatcher 创建文件监听�?
func NewFileWatcher(filePath string, callback func(string), logger logger.Logger) (*FileWatcher, error) {
	// 检查文件是否存�?
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// 创建fsnotify监听�?
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 获取文件修改时间
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fw := &FileWatcher{
		filePath:    filePath,
		callback:    callback,
		logger:      logger,
		watcher:     watcher,
		ctx:         ctx,
		cancel:      cancel,
		isRunning:   false,
		lastModTime: fileInfo.ModTime(),
		debounce:    100 * time.Millisecond, // 默认防抖时间
	}

	logger.Debug("File watcher created", "file", filePath)
	return fw, nil
}

// Start 启动监听
func (fw *FileWatcher) Start() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.isRunning {
		return fmt.Errorf("file watcher is already running")
	}

	// 添加文件到监听列�?
	// 对于文件，我们监听其所在的目录
	dir := filepath.Dir(fw.filePath)
	if err := fw.watcher.Add(dir); err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	fw.isRunning = true

	// 启动监听协程
	go fw.watchLoop()

	fw.logger.Info("File watcher started", "file", fw.filePath)
	return nil
}

// Stop 停止监听
func (fw *FileWatcher) Stop() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if !fw.isRunning {
		return nil
	}

	// 取消上下�?
	fw.cancel()

	// 关闭fsnotify监听�?
	if err := fw.watcher.Close(); err != nil {
		fw.logger.Error("Failed to close fsnotify watcher", "error", err)
	}

	fw.isRunning = false

	fw.logger.Info("File watcher stopped", "file", fw.filePath)
	return nil
}

// IsRunning 检查是否正在运�?
func (fw *FileWatcher) IsRunning() bool {
	fw.mu.RLock()
	defer fw.mu.RUnlock()
	return fw.isRunning
}

// AddPath 添加监听路径
func (fw *FileWatcher) AddPath(path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if !fw.isRunning {
		return fmt.Errorf("file watcher is not running")
	}

	if err := fw.watcher.Add(path); err != nil {
		return fmt.Errorf("failed to add path to watcher: %w", err)
	}

	fw.logger.Debug("Path added to watcher", "path", path)
	return nil
}

// RemovePath 移除监听路径
func (fw *FileWatcher) RemovePath(path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if !fw.isRunning {
		return fmt.Errorf("file watcher is not running")
	}

	if err := fw.watcher.Remove(path); err != nil {
		return fmt.Errorf("failed to remove path from watcher: %w", err)
	}

	fw.logger.Debug("Path removed from watcher", "path", path)
	return nil
}

// SetCallback 设置回调函数
func (fw *FileWatcher) SetCallback(callback func(string)) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.callback = callback
	fw.logger.Debug("Callback function updated")
}

// 私有方法

// watchLoop 监听循环
func (fw *FileWatcher) watchLoop() {
	fw.logger.Debug("File watcher loop started")

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				fw.logger.Debug("File watcher events channel closed")
				return
			}

			fw.handleEvent(event)

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				fw.logger.Debug("File watcher errors channel closed")
				return
			}

			fw.logger.Error("File watcher error", "error", err)

		case <-fw.ctx.Done():
			fw.logger.Debug("File watcher context cancelled")
			return
		}
	}
}

// handleEvent 处理文件事件
func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
	// 只处理我们关心的文件
	if event.Name != fw.filePath {
		return
	}

	fw.logger.Debug("File event received", "event", event.Op.String(), "file", event.Name)

	// 检查事件类�?
	if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
		// 文件被写入或创建
		fw.handleFileChange(event.Name)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
		// 文件被删除或重命�?
		fw.handleFileRemove(event.Name)
	}
}

// handleFileChange 处理文件变化
func (fw *FileWatcher) handleFileChange(filePath string) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fw.logger.Error("Failed to get file info after change", "error", err, "file", filePath)
		return
	}

	// 检查修改时间，实现防抖
	fw.mu.Lock()
	lastModTime := fw.lastModTime
	fw.lastModTime = fileInfo.ModTime()
	fw.mu.Unlock()

	if fileInfo.ModTime().Sub(lastModTime) < fw.debounce {
		fw.logger.Debug("File change ignored due to debounce", "file", filePath)
		return
	}

	// 延迟一小段时间，确保文件写入完�?
	time.Sleep(fw.debounce)

	// 再次检查文件是否存在（可能在写入过程中被删除）
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fw.logger.Debug("File no longer exists after change", "file", filePath)
		return
	}

	fw.logger.Info("File changed", "file", filePath)

	// 调用回调函数
	fw.mu.RLock()
	callback := fw.callback
	fw.mu.RUnlock()

	if callback != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fw.logger.Error("Panic in file change callback", "panic", r, "file", filePath)
				}
			}()

			callback(filePath)
		}()
	}
}

// handleFileRemove 处理文件删除
func (fw *FileWatcher) handleFileRemove(filePath string) {
	fw.logger.Warn("File removed or renamed", "file", filePath)

	// 可以选择停止监听或等待文件重新创�?
	// 这里我们选择继续监听，等待文件重新创�?
}

// DirectoryWatcher 目录监听�?
type DirectoryWatcher struct {
	dirPath   string
	callback  func(string, fsnotify.Op)
	logger    logger.Logger
	watcher   *fsnotify.Watcher
	config    *WatcherConfig
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
	isRunning bool
	ignoreMap map[string]bool
}

// NewDirectoryWatcher 创建目录监听�?
func NewDirectoryWatcher(dirPath string, callback func(string, fsnotify.Op), config *WatcherConfig, logger logger.Logger) (*DirectoryWatcher, error) {
	// 检查目录是否存�?
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dirPath)
	}

	if config == nil {
		config = &WatcherConfig{
			DebounceInterval: 100 * time.Millisecond,
			WatchDirectory:   true,
			Recursive:        false,
			IgnorePatterns:   []string{".git", ".svn", "node_modules", ".DS_Store"},
		}
	}

	// 创建fsnotify监听�?
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 构建忽略模式映射
	ignoreMap := make(map[string]bool)
	for _, pattern := range config.IgnorePatterns {
		ignoreMap[pattern] = true
	}

	dw := &DirectoryWatcher{
		dirPath:   dirPath,
		callback:  callback,
		logger:    logger,
		watcher:   watcher,
		config:    config,
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
		ignoreMap: ignoreMap,
	}

	logger.Debug("Directory watcher created", "directory", dirPath)
	return dw, nil
}

// Start 启动目录监听
func (dw *DirectoryWatcher) Start() error {
	dw.mu.Lock()
	defer dw.mu.Unlock()

	if dw.isRunning {
		return fmt.Errorf("directory watcher is already running")
	}

	// 添加目录到监听列�?
	if err := dw.addDirectory(dw.dirPath); err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	dw.isRunning = true

	// 启动监听协程
	go dw.watchLoop()

	dw.logger.Info("Directory watcher started", "directory", dw.dirPath)
	return nil
}

// Stop 停止目录监听
func (dw *DirectoryWatcher) Stop() error {
	dw.mu.Lock()
	defer dw.mu.Unlock()

	if !dw.isRunning {
		return nil
	}

	// 取消上下�?
	dw.cancel()

	// 关闭fsnotify监听�?
	if err := dw.watcher.Close(); err != nil {
		dw.logger.Error("Failed to close fsnotify watcher", "error", err)
	}

	dw.isRunning = false

	dw.logger.Info("Directory watcher stopped", "directory", dw.dirPath)
	return nil
}

// IsRunning 检查是否正在运�?
func (dw *DirectoryWatcher) IsRunning() bool {
	dw.mu.RLock()
	defer dw.mu.RUnlock()
	return dw.isRunning
}

// 私有方法

// addDirectory 添加目录到监�?
func (dw *DirectoryWatcher) addDirectory(dirPath string) error {
	if err := dw.watcher.Add(dirPath); err != nil {
		return err
	}

	// 如果启用递归监听
	if dw.config.Recursive {
		return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && path != dirPath {
				// 检查是否应该忽略这个目�?
				if dw.shouldIgnore(filepath.Base(path)) {
					return filepath.SkipDir
				}

				if err := dw.watcher.Add(path); err != nil {
					dw.logger.Error("Failed to add subdirectory to watcher", "error", err, "path", path)
				}
			}

			return nil
		})
	}

	return nil
}

// shouldIgnore 检查是否应该忽略文件或目录
func (dw *DirectoryWatcher) shouldIgnore(name string) bool {
	return dw.ignoreMap[name]
}

// watchLoop 监听循环
func (dw *DirectoryWatcher) watchLoop() {
	dw.logger.Debug("Directory watcher loop started")

	for {
		select {
		case event, ok := <-dw.watcher.Events:
			if !ok {
				dw.logger.Debug("Directory watcher events channel closed")
				return
			}

			dw.handleDirectoryEvent(event)

		case err, ok := <-dw.watcher.Errors:
			if !ok {
				dw.logger.Debug("Directory watcher errors channel closed")
				return
			}

			dw.logger.Error("Directory watcher error", "error", err)

		case <-dw.ctx.Done():
			dw.logger.Debug("Directory watcher context cancelled")
			return
		}
	}
}

// handleDirectoryEvent 处理目录事件
func (dw *DirectoryWatcher) handleDirectoryEvent(event fsnotify.Event) {
	// 检查是否应该忽略这个文�?
	if dw.shouldIgnore(filepath.Base(event.Name)) {
		return
	}

	dw.logger.Debug("Directory event received", "event", event.Op.String(), "file", event.Name)

	// 如果是目录创建事件且启用了递归监听
	if event.Op&fsnotify.Create == fsnotify.Create && dw.config.Recursive {
		if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
			if err := dw.addDirectory(event.Name); err != nil {
				dw.logger.Error("Failed to add new directory to watcher", "error", err, "path", event.Name)
			}
		}
	}

	// 调用回调函数
	dw.mu.RLock()
	callback := dw.callback
	dw.mu.RUnlock()

	if callback != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					dw.logger.Error("Panic in directory change callback", "panic", r, "file", event.Name)
				}
			}()

			callback(event.Name, event.Op)
		}()
	}
}
