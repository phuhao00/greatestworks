// Package logging 统一日志系统
// Author: MMO Server Team
// Created: 2024

package logging

import (
	"context"
	"fmt"
	"time"
)

// Level 日志级别
type Level int8

const (
	// TraceLevel 跟踪级别
	TraceLevel Level = iota - 2
	// DebugLevel 调试级别
	DebugLevel
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命错误级别
	FatalLevel
	// PanicLevel 恐慌级别
	PanicLevel
)

// String 返回级别字符串
func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel 解析日志级别
func ParseLevel(level string) (Level, error) {
	switch level {
	case "trace", "TRACE":
		return TraceLevel, nil
	case "debug", "DEBUG":
		return DebugLevel, nil
	case "info", "INFO":
		return InfoLevel, nil
	case "warn", "WARN", "warning", "WARNING":
		return WarnLevel, nil
	case "error", "ERROR":
		return ErrorLevel, nil
	case "fatal", "FATAL":
		return FatalLevel, nil
	case "panic", "PANIC":
		return PanicLevel, nil
	default:
		return InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// Fields 日志字段集合
type Fields map[string]interface{}

// Entry 日志条目
type Entry struct {
	Level     Level
	Message   string
	Fields    Fields
	Timestamp time.Time
	Caller    string
	Error     error
}

// Logger 日志接口
type Logger interface {
	// 基础日志方法
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})

	// 带字段的日志方法
	TraceWithFields(msg string, fields Fields)
	DebugWithFields(msg string, fields Fields)
	InfoWithFields(msg string, fields Fields)
	WarnWithFields(msg string, fields Fields)
	ErrorWithFields(msg string, fields Fields)
	FatalWithFields(msg string, fields Fields)
	PanicWithFields(msg string, fields Fields)

	// 带错误的日志方法
	TraceWithError(err error, msg string, args ...interface{})
	DebugWithError(err error, msg string, args ...interface{})
	InfoWithError(err error, msg string, args ...interface{})
	WarnWithError(err error, msg string, args ...interface{})
	ErrorWithError(err error, msg string, args ...interface{})
	FatalWithError(err error, msg string, args ...interface{})
	PanicWithError(err error, msg string, args ...interface{})

	// 上下文日志方法
	TraceContext(ctx context.Context, msg string, args ...interface{})
	DebugContext(ctx context.Context, msg string, args ...interface{})
	InfoContext(ctx context.Context, msg string, args ...interface{})
	WarnContext(ctx context.Context, msg string, args ...interface{})
	ErrorContext(ctx context.Context, msg string, args ...interface{})
	FatalContext(ctx context.Context, msg string, args ...interface{})
	PanicContext(ctx context.Context, msg string, args ...interface{})

	// 配置方法
	SetLevel(level Level)
	GetLevel() Level
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger

	// 生命周期方法
	Flush() error
	Close() error
}

// ContextLogger 上下文日志器
type ContextLogger interface {
	Logger
	WithRequestID(requestID string) Logger
	WithUserID(userID string) Logger
	WithSessionID(sessionID string) Logger
	WithTraceID(traceID string) Logger
}

// StructuredLogger 结构化日志器
type StructuredLogger interface {
	Logger
	LogEntry(entry *Entry)
	LogWithLevel(level Level, msg string, fields Fields)
}

// AsyncLogger 异步日志器
type AsyncLogger interface {
	Logger
	StartAsync() error
	StopAsync() error
	IsAsync() bool
}

// Hook 日志钩子接口
type Hook interface {
	Levels() []Level
	Fire(entry *Entry) error
}

// Formatter 日志格式化器接口
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// Writer 日志写入器接口
type Writer interface {
	Write(data []byte) (int, error)
	Flush() error
	Close() error
}

// Config 日志配置
type Config struct {
	Level      Level             `yaml:"level" json:"level"`
	Format     string            `yaml:"format" json:"format"`
	Output     string            `yaml:"output" json:"output"`
	Dir        string            `yaml:"dir" json:"dir"`
	Filename   string            `yaml:"filename" json:"filename"`
	MaxSize    int               `yaml:"max_size" json:"max_size"`
	MaxBackups int               `yaml:"max_backups" json:"max_backups"`
	MaxAge     int               `yaml:"max_age" json:"max_age"`
	Compress   bool              `yaml:"compress" json:"compress"`
	Prefix     string            `yaml:"prefix" json:"prefix"`
	Async      bool              `yaml:"async" json:"async"`
	BufferSize int               `yaml:"buffer_size" json:"buffer_size"`
	Fields     map[string]string `yaml:"fields" json:"fields"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      InfoLevel,
		Format:     "json",
		Output:     "stdout",
		Dir:        "./logs",
		Filename:   "app.log",
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   true,
		Prefix:     "mmo",
		Async:      false,
		BufferSize: 1000,
		Fields:     make(map[string]string),
	}
}

// Factory 日志工厂接口
type Factory interface {
	CreateLogger(config *Config) (Logger, error)
	CreateContextLogger(config *Config) (ContextLogger, error)
	CreateStructuredLogger(config *Config) (StructuredLogger, error)
	CreateAsyncLogger(config *Config) (AsyncLogger, error)
}

// Manager 日志管理器接口
type Manager interface {
	GetLogger(name string) Logger
	CreateLogger(name string, config *Config) (Logger, error)
	RegisterLogger(name string, logger Logger)
	RemoveLogger(name string)
	GetAllLoggers() map[string]Logger
	SetGlobalLevel(level Level)
	FlushAll() error
	CloseAll() error
}

// 便捷函数类型
type (
	// LogFunc 日志函数类型
	LogFunc func(msg string, args ...interface{})
	// LogWithFieldsFunc 带字段的日志函数类型
	LogWithFieldsFunc func(msg string, fields Fields)
	// LogWithErrorFunc 带错误的日志函数类型
	LogWithErrorFunc func(err error, msg string, args ...interface{})
	// LogContextFunc 上下文日志函数类型
	LogContextFunc func(ctx context.Context, msg string, args ...interface{})
)

// 常用字段名常量
const (
	FieldRequestID  = "request_id"
	FieldUserID     = "user_id"
	FieldSessionID  = "session_id"
	FieldTraceID    = "trace_id"
	FieldSpanID     = "span_id"
	FieldComponent  = "component"
	FieldModule     = "module"
	FieldFunction   = "function"
	FieldFile       = "file"
	FieldLine       = "line"
	FieldError      = "error"
	FieldStackTrace = "stack_trace"
	FieldDuration   = "duration"
	FieldStatus     = "status"
	FieldMethod     = "method"
	FieldURL        = "url"
	FieldIP         = "ip"
	FieldUserAgent  = "user_agent"
)

// 预定义错误
var (
	ErrLoggerNotFound    = fmt.Errorf("logger not found")
	ErrInvalidConfig     = fmt.Errorf("invalid logger config")
	ErrLoggerExists      = fmt.Errorf("logger already exists")
	ErrAsyncNotSupported = fmt.Errorf("async logging not supported")
	ErrWriterClosed      = fmt.Errorf("writer is closed")
)

// NewField 创建日志字段
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// NewFields 创建日志字段集合
func NewFields() Fields {
	return make(Fields)
}

// Add 添加字段
func (f Fields) Add(key string, value interface{}) Fields {
	f[key] = value
	return f
}

// With 链式添加字段
func (f Fields) With(key string, value interface{}) Fields {
	newFields := make(Fields)
	for k, v := range f {
		newFields[k] = v
	}
	newFields[key] = value
	return newFields
}

// Merge 合并字段
func (f Fields) Merge(other Fields) Fields {
	for k, v := range other {
		f[k] = v
	}
	return f
}

// Clone 克隆字段
func (f Fields) Clone() Fields {
	newFields := make(Fields)
	for k, v := range f {
		newFields[k] = v
	}
	return newFields
}

// NewEntry 创建日志条目
func NewEntry(level Level, msg string) *Entry {
	return &Entry{
		Level:     level,
		Message:   msg,
		Fields:    make(Fields),
		Timestamp: time.Now(),
	}
}

// WithField 添加字段到条目
func (e *Entry) WithField(key string, value interface{}) *Entry {
	if e.Fields == nil {
		e.Fields = make(Fields)
	}
	e.Fields[key] = value
	return e
}

// WithFields 添加多个字段到条目
func (e *Entry) WithFields(fields Fields) *Entry {
	if e.Fields == nil {
		e.Fields = make(Fields)
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

// WithError 添加错误到条目
func (e *Entry) WithError(err error) *Entry {
	e.Error = err
	if e.Fields == nil {
		e.Fields = make(Fields)
	}
	e.Fields[FieldError] = err.Error()
	return e
}

// WithCaller 添加调用者信息到条目
func (e *Entry) WithCaller(caller string) *Entry {
	e.Caller = caller
	return e
}