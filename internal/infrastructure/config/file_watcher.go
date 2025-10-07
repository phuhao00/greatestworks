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

// FileWatcher 文件监听器
type FileWatcher struct {
	filePath    string
	callback    func(string)
	logger      logging.Logger
	watcher     *fsnotify.Watcher
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	isRunning   bool
	lastModTime time.Time
	debounce    time.Duration
}

// WatcherConfig 监听器配置
type WatcherConfig struct {
	DebounceInterval time.Duration `json:"debounce_interval" yaml:"debounce_interval"`
	WatchDirectory   bool          `json:"watch_directory" yaml:"watch_directory"`
	Recursive        bool          `json:"recursive" yaml:"recursive"`
	IgnorePatterns   []string      `json:"ignore_patterns" yaml:"ignore_patterns"`
}

// Watcher 文件监听器接口
type Watcher interface {
	// Start 启动监听
	Start() error

	// Stop 停止监听
	Stop() error

	// IsRunning 检查是否正在运行
	IsRunning() bool

	// AddPath 添加监听路径
	AddPath(path string) error

	// RemovePath 移除监听路径
	RemovePath(path string) error

	// SetCallback 设置回调函数
	SetCallback(callback func(string))
}

// NewFileWatcher 创建文件监听器
func NewFileWatcher(filePath string, callback func(string), logger logging.Logger) (*FileWatcher, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// 创建fsnotify监听器
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

	logger.Debug("File watcher created", logging.Fields{
		"file": filePath,
	})
	return fw, nil
}

// Start 启动监听
func (fw *FileWatcher) Start() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.isRunning {
		return fmt.Errorf("file watcher is already running")
	}

	// 添加文件到监听列表
	// 对于文件，我们监听其所在的目录
	dir := filepath.Dir(fw.filePath)
	if err := fw.watcher.Add(dir); err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	fw.isRunning = true

	// 启动监听协程
	go fw.watchLoop()

	fw.logger.Info("File watcher started", logging.Fields{
		"file": fw.filePath,
	})
	return nil
}

// Stop 停止监听
func (fw *FileWatcher) Stop() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if !fw.isRunning {
		return nil
	}

	// 取消上下文
	fw.cancel()

	// 关闭fsnotify监听器
	if err := fw.watcher.Close(); err != nil {
		fw.logger.Error("Failed to close fsnotify watcher", err)
	}

	fw.isRunning = false

	fw.logger.Info("File watcher stopped", logging.Fields{
		"file": fw.filePath,
	})
	return nil
}

// IsRunning 检查是否正在运行
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

	fw.logger.Debug("Path added to watcher", logging.Fields{
		"path": path,
	})
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

	fw.logger.Debug("Path removed from watcher", logging.Fields{
		"path": path,
	})
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

			fw.logger.Error("File watcher error", err)

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

	fw.logger.Debug("File event received", logging.Fields{
		"event": event.Op.String(),
		"file":  event.Name,
	})

	// 检查事件类型
	if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
		// 文件被写入或创建
		fw.handleFileChange(event.Name)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
		// 文件被删除或重命名
		fw.handleFileRemove(event.Name)
	}
}

// handleFileChange 处理文件变化
func (fw *FileWatcher) handleFileChange(filePath string) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fw.logger.Error("Failed to get file info after change", err, logging.Fields{
			"file": filePath,
		})
		return
	}

	// 检查修改时间，实现防抖
	fw.mu.Lock()
	lastModTime := fw.lastModTime
	fw.lastModTime = fileInfo.ModTime()
	fw.mu.Unlock()

	if fileInfo.ModTime().Sub(lastModTime) < fw.debounce {
		fw.logger.Debug("File change ignored due to debounce", logging.Fields{
			"file": filePath,
		})
		return
	}

	// 延迟一小段时间，确保文件写入完成
	time.Sleep(fw.debounce)

	// 再次检查文件是否存在（可能在写入过程中被删除）
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fw.logger.Debug("File no longer exists after change", logging.Fields{
			"file": filePath,
		})
		return
	}

	fw.logger.Info("File changed", logging.Fields{
		"file": filePath,
	})

	// 调用回调函数
	fw.mu.RLock()
	callback := fw.callback
	fw.mu.RUnlock()

	if callback != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fw.logger.Error("Panic in file change callback", fmt.Errorf("panic: %v", r), logging.Fields{
						"file": filePath,
					})
				}
			}()

			callback(filePath)
		}()
	}
}

// handleFileRemove 处理文件删除
func (fw *FileWatcher) handleFileRemove(filePath string) {
	fw.logger.Warn("File removed or renamed", logging.Fields{
		"file": filePath,
	})

	// 可以选择停止监听或等待文件重新创建
	// 这里我们选择继续监听，等待文件重新创建
}
