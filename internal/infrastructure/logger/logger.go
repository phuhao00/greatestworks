// Package logger 日志系统
// Author: MMO Server Team
// Created: 2024

package logger

import "fmt"

// import "github.com/phuhao00/spoor/v2" // TODO: 暂时注释掉有问题的依赖

// Logger 日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Trace(msg string, args ...interface{})
	WithFields(fields map[string]interface{}) Logger
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	SetLevel(level LogLevel)
	GetLevel() LogLevel
	SetFormatter(formatter Formatter)
	AddHook(hook Hook)
	Close() error
}

// LogLevel 日志级别
type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Formatter 格式化器接口
type Formatter interface {
	Format(level LogLevel, msg string, fields map[string]interface{}) []byte
}

// Hook 钩子接口
type Hook interface {
	Fire(level LogLevel, msg string, fields map[string]interface{}) error
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

// simpleLogger 简单日志实现
type simpleLogger struct {
	level LogLevel
}

// NewLogger 创建新的日志记录器
func NewLogger() Logger {
	return &simpleLogger{
		level: LevelInfo,
	}
}

// GetInstance 获取单例实例
func GetInstance() Logger {
	return &simpleLogger{
		level: LevelInfo,
	}
}

// 实现Logger接口
func (l *simpleLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

func (l *simpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
}

func (l *simpleLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+msg+"\n", args...)
}

func (l *simpleLogger) Warn(msg string, args ...interface{}) {
	fmt.Printf("[WARN] "+msg+"\n", args...)
}

func (l *simpleLogger) Fatal(msg string, args ...interface{}) {
	fmt.Printf("[FATAL] "+msg+"\n", args...)
}

func (l *simpleLogger) Panic(msg string, args ...interface{}) {
	fmt.Printf("[PANIC] "+msg+"\n", args...)
}

func (l *simpleLogger) Trace(msg string, args ...interface{}) {
	fmt.Printf("[TRACE] "+msg+"\n", args...)
}

func (l *simpleLogger) WithFields(fields map[string]interface{}) Logger {
	return l
}

func (l *simpleLogger) WithField(key string, value interface{}) Logger {
	return l
}

func (l *simpleLogger) WithError(err error) Logger {
	return l
}

func (l *simpleLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *simpleLogger) GetLevel() LogLevel {
	return l.level
}

func (l *simpleLogger) SetFormatter(formatter Formatter) {
	// 简单实现，不做任何操作
}

func (l *simpleLogger) AddHook(hook Hook) {
	// 简单实现，不做任何操作
}

func (l *simpleLogger) Close() error {
	return nil
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
func parseLevel(level string) LogLevel {
	switch level {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	case "panic":
		return LevelFatal
	default:
		return LevelInfo
	}
}

// parseFormatter 解析格式化器
func parseFormatter(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "text":
		return &TextFormatter{}
	default:
		return &TextFormatter{}
	}
}

// JSONFormatter JSON格式化器
type JSONFormatter struct{}

func (f *JSONFormatter) Format(level LogLevel, msg string, fields map[string]interface{}) []byte {
	// TODO: 实现JSON格式化
	return []byte(fmt.Sprintf(`{"level":"%s","msg":"%s"}`, level.String(), msg))
}

// TextFormatter 文本格式化器
type TextFormatter struct{}

func (f *TextFormatter) Format(level LogLevel, msg string, fields map[string]interface{}) []byte {
	// TODO: 实现文本格式化
	return []byte(fmt.Sprintf("[%s] %s", level.String(), msg))
}

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "INFO"
	}
}
