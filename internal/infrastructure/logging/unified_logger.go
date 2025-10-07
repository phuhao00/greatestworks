// Package logging 统一日志系统
// Author: MMO Server Team
// Created: 2024

package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
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
	switch strings.ToLower(level) {
	case "trace":
		return TraceLevel, nil
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "panic":
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

// SimpleLogger 简单日志实现
type SimpleLogger struct {
	level     Level
	formatter Formatter
	writer    Writer
	fields    Fields
	mu        sync.RWMutex
}

// NewSimpleLogger 创建简单日志器
func NewSimpleLogger(config *Config) (Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 创建格式化器
	var formatter Formatter
	switch config.Format {
	case "json":
		formatter = &JSONFormatter{}
	case "text":
		formatter = &TextFormatter{}
	default:
		formatter = &TextFormatter{}
	}

	// 创建写入器
	var writer Writer
	switch config.Output {
	case "stdout":
		writer = &StdoutWriter{}
	case "stderr":
		writer = &StderrWriter{}
	case "file":
		fileWriter, err := NewFileWriter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create file writer: %w", err)
		}
		writer = fileWriter
	default:
		writer = &StdoutWriter{}
	}

	return &SimpleLogger{
		level:     config.Level,
		formatter: formatter,
		writer:    writer,
		fields:    make(Fields),
	}, nil
}

// 实现Logger接口
func (l *SimpleLogger) Trace(msg string, args ...interface{}) {
	l.log(TraceLevel, msg, args...)
}

func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, msg, args...)
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, msg, args...)
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, msg, args...)
}

func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	l.log(FatalLevel, msg, args...)
}

func (l *SimpleLogger) Panic(msg string, args ...interface{}) {
	l.log(PanicLevel, msg, args...)
}

func (l *SimpleLogger) TraceWithFields(msg string, fields Fields) {
	l.logWithFields(TraceLevel, msg, fields)
}

func (l *SimpleLogger) DebugWithFields(msg string, fields Fields) {
	l.logWithFields(DebugLevel, msg, fields)
}

func (l *SimpleLogger) InfoWithFields(msg string, fields Fields) {
	l.logWithFields(InfoLevel, msg, fields)
}

func (l *SimpleLogger) WarnWithFields(msg string, fields Fields) {
	l.logWithFields(WarnLevel, msg, fields)
}

func (l *SimpleLogger) ErrorWithFields(msg string, fields Fields) {
	l.logWithFields(ErrorLevel, msg, fields)
}

func (l *SimpleLogger) FatalWithFields(msg string, fields Fields) {
	l.logWithFields(FatalLevel, msg, fields)
}

func (l *SimpleLogger) PanicWithFields(msg string, fields Fields) {
	l.logWithFields(PanicLevel, msg, fields)
}

func (l *SimpleLogger) TraceWithError(err error, msg string, args ...interface{}) {
	l.logWithError(TraceLevel, err, msg, args...)
}

func (l *SimpleLogger) DebugWithError(err error, msg string, args ...interface{}) {
	l.logWithError(DebugLevel, err, msg, args...)
}

func (l *SimpleLogger) InfoWithError(err error, msg string, args ...interface{}) {
	l.logWithError(InfoLevel, err, msg, args...)
}

func (l *SimpleLogger) WarnWithError(err error, msg string, args ...interface{}) {
	l.logWithError(WarnLevel, err, msg, args...)
}

func (l *SimpleLogger) ErrorWithError(err error, msg string, args ...interface{}) {
	l.logWithError(ErrorLevel, err, msg, args...)
}

func (l *SimpleLogger) FatalWithError(err error, msg string, args ...interface{}) {
	l.logWithError(FatalLevel, err, msg, args...)
}

func (l *SimpleLogger) PanicWithError(err error, msg string, args ...interface{}) {
	l.logWithError(PanicLevel, err, msg, args...)
}

func (l *SimpleLogger) TraceContext(ctx context.Context, msg string, args ...interface{}) {
	l.logContext(TraceLevel, ctx, msg, args...)
}

func (l *SimpleLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	l.logContext(DebugLevel, ctx, msg, args...)
}

func (l *SimpleLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	l.logContext(InfoLevel, ctx, msg, args...)
}

func (l *SimpleLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	l.logContext(WarnLevel, ctx, msg, args...)
}

func (l *SimpleLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	l.logContext(ErrorLevel, ctx, msg, args...)
}

func (l *SimpleLogger) FatalContext(ctx context.Context, msg string, args ...interface{}) {
	l.logContext(FatalLevel, ctx, msg, args...)
}

func (l *SimpleLogger) PanicContext(ctx context.Context, msg string, args ...interface{}) {
	l.logContext(PanicLevel, ctx, msg, args...)
}

func (l *SimpleLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *SimpleLogger) GetLevel() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

func (l *SimpleLogger) WithField(key string, value interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newFields := make(Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &SimpleLogger{
		level:     l.level,
		formatter: l.formatter,
		writer:    l.writer,
		fields:    newFields,
	}
}

func (l *SimpleLogger) WithFields(fields Fields) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newFields := make(Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &SimpleLogger{
		level:     l.level,
		formatter: l.formatter,
		writer:    l.writer,
		fields:    newFields,
	}
}

func (l *SimpleLogger) WithError(err error) Logger {
	return l.WithField("error", err.Error())
}

func (l *SimpleLogger) WithContext(ctx context.Context) Logger {
	// 从上下文中提取字段
	fields := make(Fields)

	// 提取请求ID
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields["request_id"] = requestID
	}

	// 提取用户ID
	if userID := ctx.Value("user_id"); userID != nil {
		fields["user_id"] = userID
	}

	// 提取会话ID
	if sessionID := ctx.Value("session_id"); sessionID != nil {
		fields["session_id"] = sessionID
	}

	return l.WithFields(fields)
}

func (l *SimpleLogger) Flush() error {
	if l.writer != nil {
		return l.writer.Flush()
	}
	return nil
}

func (l *SimpleLogger) Close() error {
	if l.writer != nil {
		return l.writer.Close()
	}
	return nil
}

// 私有方法
func (l *SimpleLogger) log(level Level, msg string, args ...interface{}) {
	if level < l.GetLevel() {
		return
	}

	// 格式化消息
	formattedMsg := fmt.Sprintf(msg, args...)

	// 创建日志条目
	entry := &Entry{
		Level:     level,
		Message:   formattedMsg,
		Fields:    l.getFields(),
		Timestamp: time.Now(),
		Caller:    l.getCaller(),
	}

	l.writeEntry(entry)
}

func (l *SimpleLogger) logWithFields(level Level, msg string, fields Fields) {
	if level < l.GetLevel() {
		return
	}

	// 合并字段
	allFields := l.getFields()
	for k, v := range fields {
		allFields[k] = v
	}

	// 创建日志条目
	entry := &Entry{
		Level:     level,
		Message:   msg,
		Fields:    allFields,
		Timestamp: time.Now(),
		Caller:    l.getCaller(),
	}

	l.writeEntry(entry)
}

func (l *SimpleLogger) logWithError(level Level, err error, msg string, args ...interface{}) {
	if level < l.GetLevel() {
		return
	}

	// 格式化消息
	formattedMsg := fmt.Sprintf(msg, args...)

	// 创建日志条目
	entry := &Entry{
		Level:     level,
		Message:   formattedMsg,
		Fields:    l.getFields(),
		Timestamp: time.Now(),
		Caller:    l.getCaller(),
		Error:     err,
	}

	// 添加错误字段
	if entry.Fields == nil {
		entry.Fields = make(Fields)
	}
	entry.Fields["error"] = err.Error()

	l.writeEntry(entry)
}

func (l *SimpleLogger) logContext(level Level, ctx context.Context, msg string, args ...interface{}) {
	if level < l.GetLevel() {
		return
	}

	// 格式化消息
	formattedMsg := fmt.Sprintf(msg, args...)

	// 创建日志条目
	entry := &Entry{
		Level:     level,
		Message:   formattedMsg,
		Fields:    l.getFields(),
		Timestamp: time.Now(),
		Caller:    l.getCaller(),
	}

	// 从上下文中提取字段
	if requestID := ctx.Value("request_id"); requestID != nil {
		entry.Fields["request_id"] = requestID
	}
	if userID := ctx.Value("user_id"); userID != nil {
		entry.Fields["user_id"] = userID
	}
	if sessionID := ctx.Value("session_id"); sessionID != nil {
		entry.Fields["session_id"] = sessionID
	}

	l.writeEntry(entry)
}

func (l *SimpleLogger) getFields() Fields {
	l.mu.RLock()
	defer l.mu.RUnlock()

	fields := make(Fields)
	for k, v := range l.fields {
		fields[k] = v
	}
	return fields
}

func (l *SimpleLogger) getCaller() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func (l *SimpleLogger) writeEntry(entry *Entry) {
	// 格式化日志条目
	data, err := l.formatter.Format(entry)
	if err != nil {
		// 如果格式化失败，使用简单的文本格式
		data = []byte(fmt.Sprintf("[%s] %s %s\n", entry.Level.String(), entry.Timestamp.Format(time.RFC3339), entry.Message))
	}

	// 写入数据
	if l.writer != nil {
		l.writer.Write(data)
	}
}

// JSONFormatter JSON格式化器
type JSONFormatter struct{}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	// 简化的JSON格式化实现
	return []byte(fmt.Sprintf(`{"level":"%s","msg":"%s","timestamp":"%s"}`,
		entry.Level.String(),
		entry.Message,
		entry.Timestamp.Format(time.RFC3339))), nil
}

// TextFormatter 文本格式化器
type TextFormatter struct{}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] %s %s\n",
		entry.Level.String(),
		entry.Timestamp.Format(time.RFC3339),
		entry.Message)), nil
}

// StdoutWriter 标准输出写入器
type StdoutWriter struct{}

func (w *StdoutWriter) Write(data []byte) (int, error) {
	return os.Stdout.Write(data)
}

func (w *StdoutWriter) Flush() error {
	return nil
}

func (w *StdoutWriter) Close() error {
	return nil
}

// StderrWriter 标准错误写入器
type StderrWriter struct{}

func (w *StderrWriter) Write(data []byte) (int, error) {
	return os.Stderr.Write(data)
}

func (w *StderrWriter) Flush() error {
	return nil
}

func (w *StderrWriter) Close() error {
	return nil
}

// FileWriter 文件写入器
type FileWriter struct {
	file *os.File
}

func NewFileWriter(config *Config) (*FileWriter, error) {
	// 创建日志目录
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 打开日志文件
	filePath := filepath.Join(config.Dir, config.Filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &FileWriter{file: file}, nil
}

func (w *FileWriter) Write(data []byte) (int, error) {
	return w.file.Write(data)
}

func (w *FileWriter) Flush() error {
	return w.file.Sync()
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}

// 便捷函数
func NewLogger(config *Config) (Logger, error) {
	return NewSimpleLogger(config)
}

func NewDefaultLogger() (Logger, error) {
	return NewSimpleLogger(DefaultConfig())
}
