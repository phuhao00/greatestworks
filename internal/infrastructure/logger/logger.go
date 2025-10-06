// Package logger 日志系统
// Author: MMO Server Team
// Created: 2024

package logger

import "github.com/phuhao00/spoor/v2"

// Logger 日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Trace(msg string, args ...interface{})
	WithFields(fields map[string]interface{}) spoor.Logger
	WithField(key string, value interface{}) spoor.Logger
	WithError(err error) spoor.Logger
	SetLevel(level spoor.LogLevel)
	GetLevel() spoor.LogLevel
	SetFormatter(formatter spoor.Formatter)
	AddHook(hook spoor.Hook)
	Close() error
}

// Config 日志配置
type Config struct {
	Filename   string            `json:"filename"`
	MaxSize    int               `json:"max_size"`
	MaxBackups int               `json:"max_backups"`
	MaxAge     int               `json:"max_age"`
	Compress   bool              `json:"compress"`
	Prefix     string            `json:"prefix"`
	Async      bool              `json:"async"`
	BufferSize int               `json:"buffer_size"`
	Fields     map[string]string `json:"fields"`
	Level      string            `json:"level"`
	Format     string            `json:"format"`
}

// NewConfig 创建默认配置
func NewConfig() *Config {
	return &Config{
		Filename:   "app.log",
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   true,
		Prefix:     "mmo",
		Async:      false,
		BufferSize: 1000,
		Fields:     make(map[string]string),
		Level:      "info",
		Format:     "text",
	}
}

// NewLogger 创建新的日志记录器（使用spoor单例）
func NewLogger() Logger {
	return GetInstance()
}

// NewLoggerWithConfig 根据配置创建日志记录器
func NewLoggerWithConfig(config *Config) Logger {
	instance := GetInstance()

	// 设置日志级别
	if config.Level != "" {
		level := parseLevel(config.Level)
		instance.SetLevel(level)
	}

	// 设置格式化器
	if config.Format != "" {
		formatter := parseFormatter(config.Format)
		instance.SetFormatter(formatter)
	}

	return instance
}

// parseLevel 解析日志级别
func parseLevel(level string) spoor.LogLevel {
	switch level {
	case "trace":
		return spoor.LevelDebug // spoor v2没有TraceLevel，使用DebugLevel
	case "debug":
		return spoor.LevelDebug
	case "info":
		return spoor.LevelInfo
	case "warn":
		return spoor.LevelWarn
	case "error":
		return spoor.LevelError
	case "fatal":
		return spoor.LevelFatal
	case "panic":
		return spoor.LevelFatal // spoor v2没有PanicLevel，使用FatalLevel
	default:
		return spoor.LevelInfo
	}
}

// parseFormatter 解析格式化器
func parseFormatter(format string) spoor.Formatter {
	switch format {
	case "json":
		return &spoor.JSONFormatter{}
	case "text":
		return &spoor.TextFormatter{}
	default:
		return &spoor.TextFormatter{}
	}
}
