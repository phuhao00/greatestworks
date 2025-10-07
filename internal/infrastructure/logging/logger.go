// Package logging 提供统一的日志记录接口
package logging

import (
	"context"
	"time"
)

// Level 日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
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
	default:
		return "UNKNOWN"
	}
}

// Fields 日志字段类型
type Fields map[string]interface{}

// Logger 日志记录器接口
type Logger interface {
	// 基础日志方法
	Debug(msg string, fields ...Fields)
	Info(msg string, fields ...Fields)
	Warn(msg string, fields ...Fields)
	Error(msg string, err error, fields ...Fields)
	Fatal(msg string, err error, fields ...Fields)

	// 带上下文的日志方法
	DebugWithContext(ctx context.Context, msg string, fields ...Fields)
	InfoWithContext(ctx context.Context, msg string, fields ...Fields)
	WarnWithContext(ctx context.Context, msg string, fields ...Fields)
	ErrorWithContext(ctx context.Context, msg string, err error, fields ...Fields)
	FatalWithContext(ctx context.Context, msg string, err error, fields ...Fields)

	// 结构化日志方法
	WithFields(fields Fields) Logger
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger

	// 设置日志级别
	SetLevel(level Level)
	GetLevel() Level

	// 关闭日志记录器
	Close() error
}

// LogEntry 日志条目
type LogEntry struct {
	Level     Level                  `json:"level"`
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Error     string                 `json:"error,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
}

// Formatter 日志格式化器接口
type Formatter interface {
	Format(entry *LogEntry) ([]byte, error)
}

// Writer 日志写入器接口
type Writer interface {
	Write(entry *LogEntry) error
	Close() error
}

// Config 日志配置
type Config struct {
	Level      Level  `json:"level"`
	Format     string `json:"format"`      // json, text
	Output     string `json:"output"`      // stdout, stderr, file
	FilePath   string `json:"file_path"`   // 文件路径（当output为file时）
	MaxSize    int    `json:"max_size"`    // 文件最大大小（MB）
	MaxBackups int    `json:"max_backups"` // 最大备份文件数
	MaxAge     int    `json:"max_age"`     // 文件最大保存天数
	Compress   bool   `json:"compress"`    // 是否压缩备份文件
}

// NewLogger 创建新的日志记录器
func NewLogger(config *Config) (Logger, error) {
	// 这里应该根据配置创建具体的日志记录器实现
	// 为了简化，这里返回一个默认实现
	return &defaultLogger{
		level:  config.Level,
		fields: make(Fields),
	}, nil
}

// defaultLogger 默认日志记录器实现
type defaultLogger struct {
	level  Level
	fields Fields
}

func (l *defaultLogger) Debug(msg string, fields ...Fields) {
	l.log(DebugLevel, msg, nil, fields...)
}

func (l *defaultLogger) Info(msg string, fields ...Fields) {
	l.log(InfoLevel, msg, nil, fields...)
}

func (l *defaultLogger) Warn(msg string, fields ...Fields) {
	l.log(WarnLevel, msg, nil, fields...)
}

func (l *defaultLogger) Error(msg string, err error, fields ...Fields) {
	l.log(ErrorLevel, msg, err, fields...)
}

func (l *defaultLogger) Fatal(msg string, err error, fields ...Fields) {
	l.log(FatalLevel, msg, err, fields...)
}

func (l *defaultLogger) DebugWithContext(ctx context.Context, msg string, fields ...Fields) {
	l.logWithContext(ctx, DebugLevel, msg, nil, fields...)
}

func (l *defaultLogger) InfoWithContext(ctx context.Context, msg string, fields ...Fields) {
	l.logWithContext(ctx, InfoLevel, msg, nil, fields...)
}

func (l *defaultLogger) WarnWithContext(ctx context.Context, msg string, fields ...Fields) {
	l.logWithContext(ctx, WarnLevel, msg, nil, fields...)
}

func (l *defaultLogger) ErrorWithContext(ctx context.Context, msg string, err error, fields ...Fields) {
	l.logWithContext(ctx, ErrorLevel, msg, err, fields...)
}

func (l *defaultLogger) FatalWithContext(ctx context.Context, msg string, err error, fields ...Fields) {
	l.logWithContext(ctx, FatalLevel, msg, err, fields...)
}

func (l *defaultLogger) WithFields(fields Fields) Logger {
	newFields := make(Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &defaultLogger{
		level:  l.level,
		fields: newFields,
	}
}

func (l *defaultLogger) WithField(key string, value interface{}) Logger {
	return l.WithFields(Fields{key: value})
}

func (l *defaultLogger) WithError(err error) Logger {
	return l.WithField("error", err.Error())
}

func (l *defaultLogger) WithContext(ctx context.Context) Logger {
	// 这里可以从上下文中提取跟踪信息等
	return l
}

func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

func (l *defaultLogger) GetLevel() Level {
	return l.level
}

func (l *defaultLogger) Close() error {
	return nil
}

func (l *defaultLogger) log(level Level, msg string, err error, fields ...Fields) {
	if level < l.level {
		return
	}

	entry := &LogEntry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   msg,
		Fields:    make(map[string]interface{}),
	}

	// 合并字段
	for k, v := range l.fields {
		entry.Fields[k] = v
	}
	for _, f := range fields {
		for k, v := range f {
			entry.Fields[k] = v
		}
	}

	if err != nil {
		entry.Error = err.Error()
	}

	// 这里应该使用实际的格式化器和写入器
	// 为了简化，这里只是打印到控制台
	println(entry.String())
}

func (l *defaultLogger) logWithContext(ctx context.Context, level Level, msg string, err error, fields ...Fields) {
	// 从上下文中提取跟踪信息等
	// 这里简化处理
	l.log(level, msg, err, fields...)
}

func (e *LogEntry) String() string {
	// 简化的字符串表示
	return e.Timestamp.Format(time.RFC3339) + " [" + e.Level.String() + "] " + e.Message
}

// NewBaseLogger creates a new base logger
func NewBaseLogger(level Level) Logger {
	return &defaultLogger{
		level:  level,
		fields: make(map[string]interface{}),
	}
}
