package logger

import (
	"sync"

	"github.com/phuhao00/spoor/v2"
)

// SpoorLogger spoor日志单例管理器
type SpoorLogger struct {
	logger spoor.Logger
	mu     sync.RWMutex
}

var (
	instance *SpoorLogger
	once     sync.Once
)

// GetInstance 获取spoor日志单例实例
func GetInstance() *SpoorLogger {
	once.Do(func() {
		// 创建默认的spoor logger
		logger := spoor.NewWithDefaults()
		instance = &SpoorLogger{
			logger: logger,
		}
	})
	return instance
}

// GetLogger 获取spoor logger实例
func (s *SpoorLogger) GetLogger() spoor.Logger {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.logger
}

// SetLogger 设置spoor logger实例
func (s *SpoorLogger) SetLogger(logger spoor.Logger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logger = logger
}

// Info 记录信息日志
func (s *SpoorLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		s.GetLogger().Infof(msg, args...)
	} else {
		s.GetLogger().Info(msg)
	}
}

// Error 记录错误日志
func (s *SpoorLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		s.GetLogger().Errorf(msg, args...)
	} else {
		s.GetLogger().Error(msg)
	}
}

// Debug 记录调试日志
func (s *SpoorLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		s.GetLogger().Debugf(msg, args...)
	} else {
		s.GetLogger().Debug(msg)
	}
}

// Warn 记录警告日志
func (s *SpoorLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		s.GetLogger().Warnf(msg, args...)
	} else {
		s.GetLogger().Warn(msg)
	}
}

// Fatal 记录致命错误日志
func (s *SpoorLogger) Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		s.GetLogger().Fatalf(msg, args...)
	} else {
		s.GetLogger().Fatal(msg)
	}
}

// Panic 记录恐慌日志（spoor v2没有Panic方法，使用Fatal代替）
func (s *SpoorLogger) Panic(msg string, args ...interface{}) {
	if len(args) > 0 {
		s.GetLogger().Fatalf(msg, args...)
	} else {
		s.GetLogger().Fatal(msg)
	}
}

// Trace 记录跟踪日志（spoor v2没有Trace方法，使用Debug代替）
func (s *SpoorLogger) Trace(msg string, args ...interface{}) {
	if len(args) > 0 {
		s.GetLogger().Debugf(msg, args...)
	} else {
		s.GetLogger().Debug(msg)
	}
}

// WithFields 创建带字段的日志记录器
func (s *SpoorLogger) WithFields(fields map[string]interface{}) spoor.Logger {
	return s.GetLogger().WithFields(fields)
}

// WithField 创建带单个字段的日志记录器
func (s *SpoorLogger) WithField(key string, value interface{}) spoor.Logger {
	return s.GetLogger().WithField(key, value)
}

// WithError 创建带错误的日志记录器
func (s *SpoorLogger) WithError(err error) spoor.Logger {
	return s.GetLogger().WithError(err)
}

// SetLevel 设置日志级别
func (s *SpoorLogger) SetLevel(level spoor.LogLevel) {
	spoor.SetLevel(level)
}

// GetLevel 获取当前日志级别
func (s *SpoorLogger) GetLevel() spoor.LogLevel {
	return spoor.GetLevel()
}

// SetFormatter 设置日志格式化器（spoor v2通过配置设置）
func (s *SpoorLogger) SetFormatter(formatter spoor.Formatter) {
	// spoor v2的格式化器通过配置设置，这里暂时不实现
}

// AddHook 添加日志钩子（spoor v2通过配置设置）
func (s *SpoorLogger) AddHook(hook spoor.Hook) {
	// spoor v2的钩子通过配置设置，这里暂时不实现
}

// Close 关闭日志记录器
func (s *SpoorLogger) Close() error {
	// spoor v2的logger没有Close方法，返回nil
	return nil
}
