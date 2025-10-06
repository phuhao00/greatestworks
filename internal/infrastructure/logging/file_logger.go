// Package logging 文件日志器
// Author: MMO Server Team
// Created: 2024

package logging

import (
	"context"
	"fmt"

	// "io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileLogger 文件日志器
type FileLogger struct {
	file       *os.File
	filePath   string
	level      Level
	fields     Fields
	ctx        context.Context
	mutex      sync.RWMutex
	config     *Config
	formatter  Formatter
	rotator    *FileRotator
	buffer     []byte
	bufferSize int
	lastFlush  time.Time
	flushTimer *time.Timer
	closed     bool
}

// FileRotator 文件轮转器
type FileRotator struct {
	dir        string
	filename   string
	maxSize    int64
	maxBackups int
	maxAge     time.Duration
	compress   bool
	mutex      sync.Mutex
}

// NewFileLogger 创建文件日志器
func NewFileLogger(config *Config) (*FileLogger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 确保日志目录存在
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 构建日志文件路径
	filename := config.Filename
	if filename == "" {
		filename = "app.log"
	}
	filePath := filepath.Join(config.Dir, filename)

	// 打开日志文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// 创建文件轮转器
	rotator := &FileRotator{
		dir:        config.Dir,
		filename:   filename,
		maxSize:    int64(config.MaxSize) * 1024 * 1024, // MB to bytes
		maxBackups: config.MaxBackups,
		maxAge:     time.Duration(config.MaxAge) * 24 * time.Hour, // days to duration
		compress:   config.Compress,
	}

	// 创建格式化器
	formatter := NewTextFormatter()
	if config.Format == "json" {
		formatter = NewJSONFormatter()
	}

	// 创建缓冲区
	bufferSize := config.BufferSize
	if bufferSize <= 0 {
		bufferSize = 1024
	}

	logger := &FileLogger{
		file:       file,
		filePath:   filePath,
		level:      config.Level,
		fields:     make(Fields),
		config:     config,
		formatter:  formatter,
		rotator:    rotator,
		buffer:     make([]byte, 0, bufferSize),
		bufferSize: bufferSize,
		lastFlush:  time.Now(),
	}

	// 启动定时刷新
	logger.startFlushTimer()

	return logger, nil
}

// 基础日志方法实现

func (fl *FileLogger) Trace(msg string, args ...interface{}) {
	fl.log(TraceLevel, fmt.Sprintf(msg, args...), nil)
}

func (fl *FileLogger) Debug(msg string, args ...interface{}) {
	fl.log(DebugLevel, fmt.Sprintf(msg, args...), nil)
}

func (fl *FileLogger) Info(msg string, args ...interface{}) {
	fl.log(InfoLevel, fmt.Sprintf(msg, args...), nil)
}

func (fl *FileLogger) Warn(msg string, args ...interface{}) {
	fl.log(WarnLevel, fmt.Sprintf(msg, args...), nil)
}

func (fl *FileLogger) Error(msg string, args ...interface{}) {
	fl.log(ErrorLevel, fmt.Sprintf(msg, args...), nil)
}

func (fl *FileLogger) Fatal(msg string, args ...interface{}) {
	fl.log(FatalLevel, fmt.Sprintf(msg, args...), nil)
	os.Exit(1)
}

func (fl *FileLogger) Panic(msg string, args ...interface{}) {
	fl.log(PanicLevel, fmt.Sprintf(msg, args...), nil)
	panic(fmt.Sprintf(msg, args...))
}

// 带字段的日志方法实现

func (fl *FileLogger) TraceWithFields(msg string, fields Fields) {
	fl.log(TraceLevel, msg, fields)
}

func (fl *FileLogger) DebugWithFields(msg string, fields Fields) {
	fl.log(DebugLevel, msg, fields)
}

func (fl *FileLogger) InfoWithFields(msg string, fields Fields) {
	fl.log(InfoLevel, msg, fields)
}

func (fl *FileLogger) WarnWithFields(msg string, fields Fields) {
	fl.log(WarnLevel, msg, fields)
}

func (fl *FileLogger) ErrorWithFields(msg string, fields Fields) {
	fl.log(ErrorLevel, msg, fields)
}

func (fl *FileLogger) FatalWithFields(msg string, fields Fields) {
	fl.log(FatalLevel, msg, fields)
	os.Exit(1)
}

func (fl *FileLogger) PanicWithFields(msg string, fields Fields) {
	fl.log(PanicLevel, msg, fields)
	panic(msg)
}

// 带错误的日志方法实现

func (fl *FileLogger) TraceWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	fl.log(TraceLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) DebugWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	fl.log(DebugLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) InfoWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	fl.log(InfoLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) WarnWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	fl.log(WarnLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) ErrorWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	fl.log(ErrorLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) FatalWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	fl.log(FatalLevel, fmt.Sprintf(msg, args...), fields)
	os.Exit(1)
}

func (fl *FileLogger) PanicWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	fl.log(PanicLevel, fmt.Sprintf(msg, args...), fields)
	panic(fmt.Sprintf(msg, args...))
}

// 上下文日志方法实现

func (fl *FileLogger) TraceContext(ctx context.Context, msg string, args ...interface{}) {
	fields := fl.extractContextFields(ctx)
	fl.log(TraceLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	fields := fl.extractContextFields(ctx)
	fl.log(DebugLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	fields := fl.extractContextFields(ctx)
	fl.log(InfoLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	fields := fl.extractContextFields(ctx)
	fl.log(WarnLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	fields := fl.extractContextFields(ctx)
	fl.log(ErrorLevel, fmt.Sprintf(msg, args...), fields)
}

func (fl *FileLogger) FatalContext(ctx context.Context, msg string, args ...interface{}) {
	fields := fl.extractContextFields(ctx)
	fl.log(FatalLevel, fmt.Sprintf(msg, args...), fields)
	os.Exit(1)
}

func (fl *FileLogger) PanicContext(ctx context.Context, msg string, args ...interface{}) {
	fields := fl.extractContextFields(ctx)
	fl.log(PanicLevel, fmt.Sprintf(msg, args...), fields)
	panic(fmt.Sprintf(msg, args...))
}

// 配置方法实现

func (fl *FileLogger) SetLevel(level Level) {
	fl.mutex.Lock()
	defer fl.mutex.Unlock()
	fl.level = level
}

func (fl *FileLogger) GetLevel() Level {
	fl.mutex.RLock()
	defer fl.mutex.RUnlock()
	return fl.level
}

func (fl *FileLogger) WithField(key string, value interface{}) Logger {
	newFields := fl.fields.Clone()
	newFields[key] = value
	return &FileLogger{
		file:       fl.file,
		filePath:   fl.filePath,
		level:      fl.level,
		fields:     newFields,
		ctx:        fl.ctx,
		config:     fl.config,
		formatter:  fl.formatter,
		rotator:    fl.rotator,
		buffer:     fl.buffer,
		bufferSize: fl.bufferSize,
		lastFlush:  fl.lastFlush,
	}
}

func (fl *FileLogger) WithFields(fields Fields) Logger {
	newFields := fl.fields.Clone().Merge(fields)
	return &FileLogger{
		file:       fl.file,
		filePath:   fl.filePath,
		level:      fl.level,
		fields:     newFields,
		ctx:        fl.ctx,
		config:     fl.config,
		formatter:  fl.formatter,
		rotator:    fl.rotator,
		buffer:     fl.buffer,
		bufferSize: fl.bufferSize,
		lastFlush:  fl.lastFlush,
	}
}

func (fl *FileLogger) WithError(err error) Logger {
	return fl.WithField(FieldError, err.Error())
}

func (fl *FileLogger) WithContext(ctx context.Context) Logger {
	return &FileLogger{
		file:       fl.file,
		filePath:   fl.filePath,
		level:      fl.level,
		fields:     fl.fields,
		ctx:        ctx,
		config:     fl.config,
		formatter:  fl.formatter,
		rotator:    fl.rotator,
		buffer:     fl.buffer,
		bufferSize: fl.bufferSize,
		lastFlush:  fl.lastFlush,
	}
}

// 生命周期方法实现

func (fl *FileLogger) Flush() error {
	fl.mutex.Lock()
	defer fl.mutex.Unlock()

	if fl.closed {
		return ErrWriterClosed
	}

	if len(fl.buffer) > 0 {
		if _, err := fl.file.Write(fl.buffer); err != nil {
			return err
		}
		fl.buffer = fl.buffer[:0]
	}

	fl.lastFlush = time.Now()
	return fl.file.Sync()
}

func (fl *FileLogger) Close() error {
	fl.mutex.Lock()
	defer fl.mutex.Unlock()

	if fl.closed {
		return nil
	}

	fl.closed = true

	// 停止刷新定时器
	if fl.flushTimer != nil {
		fl.flushTimer.Stop()
	}

	// 刷新缓冲区
	if len(fl.buffer) > 0 {
		fl.file.Write(fl.buffer)
	}

	// 关闭文件
	return fl.file.Close()
}

// 核心日志方法

func (fl *FileLogger) log(level Level, msg string, fields Fields) {
	// 检查日志级别
	if level < fl.level {
		return
	}

	// 合并字段
	allFields := fl.fields.Clone()
	if fields != nil {
		allFields.Merge(fields)
	}

	// 创建日志条目
	entry := &Entry{
		Level:     level,
		Message:   msg,
		Fields:    allFields,
		Timestamp: time.Now(),
	}

	// 格式化日志
	data, err := fl.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to format log entry: %v\n", err)
		return
	}

	// 写入日志
	fl.write(data)
}

// write 写入日志数据
func (fl *FileLogger) write(data []byte) {
	fl.mutex.Lock()
	defer fl.mutex.Unlock()

	if fl.closed {
		return
	}

	// 检查是否需要轮转
	if err := fl.rotator.ShouldRotate(fl.filePath); err == nil {
		fl.rotate()
	}

	// 如果缓冲区启用
	if fl.bufferSize > 0 {
		// 检查缓冲区是否有足够空间
		if len(fl.buffer)+len(data) > fl.bufferSize {
			// 刷新缓冲区
			fl.file.Write(fl.buffer)
			fl.buffer = fl.buffer[:0]
		}

		// 添加到缓冲区
		fl.buffer = append(fl.buffer, data...)

		// 如果是高级别日志，立即刷新
		if len(data) > 0 && (data[len(data)-1] == '\n') {
			entry := string(data)
			if contains(entry, "ERROR") || contains(entry, "FATAL") || contains(entry, "PANIC") {
				fl.file.Write(fl.buffer)
				fl.buffer = fl.buffer[:0]
				fl.file.Sync()
			}
		}
	} else {
		// 直接写入文件
		fl.file.Write(data)
		fl.file.Sync()
	}
}

// rotate 轮转日志文件
func (fl *FileLogger) rotate() {
	// 关闭当前文件
	fl.file.Close()

	// 执行轮转
	fl.rotator.Rotate(fl.filePath)

	// 重新打开文件
	file, err := os.OpenFile(fl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to reopen log file after rotation: %v\n", err)
		return
	}

	fl.file = file
}

// extractContextFields 从上下文提取字段
func (fl *FileLogger) extractContextFields(ctx context.Context) Fields {
	fields := make(Fields)

	if ctx == nil {
		return fields
	}

	// 提取常用字段
	if requestID := ctx.Value(FieldRequestID); requestID != nil {
		fields[FieldRequestID] = requestID
	}

	if userID := ctx.Value(FieldUserID); userID != nil {
		fields[FieldUserID] = userID
	}

	if sessionID := ctx.Value(FieldSessionID); sessionID != nil {
		fields[FieldSessionID] = sessionID
	}

	if traceID := ctx.Value(FieldTraceID); traceID != nil {
		fields[FieldTraceID] = traceID
	}

	return fields
}

// startFlushTimer 启动刷新定时器
func (fl *FileLogger) startFlushTimer() {
	if fl.bufferSize <= 0 {
		return
	}

	fl.flushTimer = time.AfterFunc(5*time.Second, func() {
		fl.Flush()
		fl.startFlushTimer() // 重新启动定时器
	})
}

// FileRotator 方法实现

// ShouldRotate 检查是否应该轮转
func (fr *FileRotator) ShouldRotate(filePath string) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if stat.Size() >= fr.maxSize {
		return nil // 需要轮转
	}

	return fmt.Errorf("no rotation needed")
}

// Rotate 执行轮转
func (fr *FileRotator) Rotate(filePath string) error {
	fr.mutex.Lock()
	defer fr.mutex.Unlock()

	// 生成备份文件名
	backupPath := fr.generateBackupPath(filePath)

	// 重命名当前文件
	if err := os.Rename(filePath, backupPath); err != nil {
		return err
	}

	// 清理旧文件
	go fr.cleanup()

	return nil
}

// generateBackupPath 生成备份文件路径
func (fr *FileRotator) generateBackupPath(filePath string) string {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	timestamp := time.Now().Format("20060102-150405")
	backupName := fmt.Sprintf("%s-%s%s", name, timestamp, ext)

	return filepath.Join(dir, backupName)
}

// cleanup 清理旧文件
func (fr *FileRotator) cleanup() {
	// 清理超过数量限制的文件
	if fr.maxBackups > 0 {
		fr.cleanupByCount()
	}

	// 清理超过时间限制的文件
	if fr.maxAge > 0 {
		fr.cleanupByAge()
	}
}

// cleanupByCount 按数量清理
func (fr *FileRotator) cleanupByCount() {
	// 实现按数量清理逻辑
	// 这里简化实现
}

// cleanupByAge 按时间清理
func (fr *FileRotator) cleanupByAge() {
	// 实现按时间清理逻辑
	// 这里简化实现
}

// 辅助函数

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
