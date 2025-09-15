// Package logging 控制台日志器
// Author: MMO Server Team
// Created: 2024

package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// ConsoleLogger 控制台日志器
type ConsoleLogger struct {
	writer    io.Writer
	level     Level
	fields    Fields
	ctx       context.Context
	mutex     sync.RWMutex
	config    *Config
	formatter Formatter
	colorized bool
}

// ANSI 颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"

	// 高亮颜色
	ColorBrightRed    = "\033[91m"
	ColorBrightGreen  = "\033[92m"
	ColorBrightYellow = "\033[93m"
	ColorBrightBlue   = "\033[94m"
	ColorBrightPurple = "\033[95m"
	ColorBrightCyan   = "\033[96m"
	ColorBrightWhite  = "\033[97m"

	// 背景颜色
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgPurple = "\033[45m"
	BgCyan   = "\033[46m"
	BgWhite  = "\033[47m"

	// 样式
	StyleBold      = "\033[1m"
	StyleDim       = "\033[2m"
	StyleItalic    = "\033[3m"
	StyleUnderline = "\033[4m"
	StyleBlink     = "\033[5m"
	StyleReverse   = "\033[7m"
	StyleStrike    = "\033[9m"
)

// NewConsoleLogger 创建控制台日志器
func NewConsoleLogger(config *Config) (*ConsoleLogger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 选择输出流
	var writer io.Writer
	switch config.Output {
	case "stdout", "console":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		writer = os.Stdout
	}

	// 创建格式化器
	formatter := NewConsoleFormatter()
	if config.Format == "json" {
		formatter = NewJSONFormatter()
	}

	// 检查是否支持颜色
	colorized := supportsColor()

	return &ConsoleLogger{
		writer:    writer,
		level:     config.Level,
		fields:    make(Fields),
		config:    config,
		formatter: formatter,
		colorized: colorized,
	}, nil
}

// supportsColor 检查终端是否支持颜色
func supportsColor() bool {
	// Windows 检查
	if runtime.GOOS == "windows" {
		// Windows 10 及以上版本支持 ANSI 颜色
		return os.Getenv("TERM") != "" || os.Getenv("ConEmuANSI") == "ON"
	}

	// Unix-like 系统检查
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}

	// 检查是否在 CI 环境中
	if os.Getenv("CI") != "" {
		return false
	}

	// 检查 NO_COLOR 环境变量
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	return true
}

// 基础日志方法实现

func (cl *ConsoleLogger) Trace(msg string, args ...interface{}) {
	cl.log(TraceLevel, fmt.Sprintf(msg, args...), nil)
}

func (cl *ConsoleLogger) Debug(msg string, args ...interface{}) {
	cl.log(DebugLevel, fmt.Sprintf(msg, args...), nil)
}

func (cl *ConsoleLogger) Info(msg string, args ...interface{}) {
	cl.log(InfoLevel, fmt.Sprintf(msg, args...), nil)
}

func (cl *ConsoleLogger) Warn(msg string, args ...interface{}) {
	cl.log(WarnLevel, fmt.Sprintf(msg, args...), nil)
}

func (cl *ConsoleLogger) Error(msg string, args ...interface{}) {
	cl.log(ErrorLevel, fmt.Sprintf(msg, args...), nil)
}

func (cl *ConsoleLogger) Fatal(msg string, args ...interface{}) {
	cl.log(FatalLevel, fmt.Sprintf(msg, args...), nil)
	os.Exit(1)
}

func (cl *ConsoleLogger) Panic(msg string, args ...interface{}) {
	cl.log(PanicLevel, fmt.Sprintf(msg, args...), nil)
	panic(fmt.Sprintf(msg, args...))
}

// 带字段的日志方法实现

func (cl *ConsoleLogger) TraceWithFields(msg string, fields Fields) {
	cl.log(TraceLevel, msg, fields)
}

func (cl *ConsoleLogger) DebugWithFields(msg string, fields Fields) {
	cl.log(DebugLevel, msg, fields)
}

func (cl *ConsoleLogger) InfoWithFields(msg string, fields Fields) {
	cl.log(InfoLevel, msg, fields)
}

func (cl *ConsoleLogger) WarnWithFields(msg string, fields Fields) {
	cl.log(WarnLevel, msg, fields)
}

func (cl *ConsoleLogger) ErrorWithFields(msg string, fields Fields) {
	cl.log(ErrorLevel, msg, fields)
}

func (cl *ConsoleLogger) FatalWithFields(msg string, fields Fields) {
	cl.log(FatalLevel, msg, fields)
	os.Exit(1)
}

func (cl *ConsoleLogger) PanicWithFields(msg string, fields Fields) {
	cl.log(PanicLevel, msg, fields)
	panic(msg)
}

// 带错误的日志方法实现

func (cl *ConsoleLogger) TraceWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	cl.log(TraceLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) DebugWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	cl.log(DebugLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) InfoWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	cl.log(InfoLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) WarnWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	cl.log(WarnLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) ErrorWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	cl.log(ErrorLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) FatalWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	cl.log(FatalLevel, fmt.Sprintf(msg, args...), fields)
	os.Exit(1)
}

func (cl *ConsoleLogger) PanicWithError(err error, msg string, args ...interface{}) {
	fields := Fields{FieldError: err.Error()}
	cl.log(PanicLevel, fmt.Sprintf(msg, args...), fields)
	panic(fmt.Sprintf(msg, args...))
}

// 上下文日志方法实现

func (cl *ConsoleLogger) TraceContext(ctx context.Context, msg string, args ...interface{}) {
	fields := cl.extractContextFields(ctx)
	cl.log(TraceLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	fields := cl.extractContextFields(ctx)
	cl.log(DebugLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	fields := cl.extractContextFields(ctx)
	cl.log(InfoLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	fields := cl.extractContextFields(ctx)
	cl.log(WarnLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	fields := cl.extractContextFields(ctx)
	cl.log(ErrorLevel, fmt.Sprintf(msg, args...), fields)
}

func (cl *ConsoleLogger) FatalContext(ctx context.Context, msg string, args ...interface{}) {
	fields := cl.extractContextFields(ctx)
	cl.log(FatalLevel, fmt.Sprintf(msg, args...), fields)
	os.Exit(1)
}

func (cl *ConsoleLogger) PanicContext(ctx context.Context, msg string, args ...interface{}) {
	fields := cl.extractContextFields(ctx)
	cl.log(PanicLevel, fmt.Sprintf(msg, args...), fields)
	panic(fmt.Sprintf(msg, args...))
}

// 配置方法实现

func (cl *ConsoleLogger) SetLevel(level Level) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.level = level
}

func (cl *ConsoleLogger) GetLevel() Level {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()
	return cl.level
}

func (cl *ConsoleLogger) WithField(key string, value interface{}) Logger {
	newFields := cl.fields.Clone()
	newFields[key] = value
	return &ConsoleLogger{
		writer:    cl.writer,
		level:     cl.level,
		fields:    newFields,
		ctx:       cl.ctx,
		config:    cl.config,
		formatter: cl.formatter,
		colorized: cl.colorized,
	}
}

func (cl *ConsoleLogger) WithFields(fields Fields) Logger {
	newFields := cl.fields.Clone().Merge(fields)
	return &ConsoleLogger{
		writer:    cl.writer,
		level:     cl.level,
		fields:    newFields,
		ctx:       cl.ctx,
		config:    cl.config,
		formatter: cl.formatter,
		colorized: cl.colorized,
	}
}

func (cl *ConsoleLogger) WithError(err error) Logger {
	return cl.WithField(FieldError, err.Error())
}

func (cl *ConsoleLogger) WithContext(ctx context.Context) Logger {
	return &ConsoleLogger{
		writer:    cl.writer,
		level:     cl.level,
		fields:    cl.fields,
		ctx:       ctx,
		config:    cl.config,
		formatter: cl.formatter,
		colorized: cl.colorized,
	}
}

// 生命周期方法实现

func (cl *ConsoleLogger) Flush() error {
	// 控制台日志器不需要刷新
	return nil
}

func (cl *ConsoleLogger) Close() error {
	// 控制台日志器不需要关闭
	return nil
}

// 核心日志方法

func (cl *ConsoleLogger) log(level Level, msg string, fields Fields) {
	// 检查日志级别
	if level < cl.level {
		return
	}

	// 合并字段
	allFields := cl.fields.Clone()
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

	// 如果启用颜色，添加颜色信息
	if cl.colorized {
		entry = cl.colorizeEntry(entry)
	}

	// 格式化日志
	data, err := cl.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to format log entry: %v\n", err)
		return
	}

	// 写入控制台
	cl.write(data)
}

// write 写入数据到控制台
func (cl *ConsoleLogger) write(data []byte) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.writer.Write(data)
}

// colorizeEntry 为日志条目添加颜色
func (cl *ConsoleLogger) colorizeEntry(entry *Entry) *Entry {
	if !cl.colorized {
		return entry
	}

	// 复制条目
	colorizedEntry := &Entry{
		Level:     entry.Level,
		Message:   entry.Message,
		Fields:    entry.Fields.Clone(),
		Timestamp: entry.Timestamp,
		Caller:    entry.Caller,
		Error:     entry.Error,
	}

	// 根据级别添加颜色
	color := cl.getLevelColor(entry.Level)
	levelText := cl.getColorizedLevel(entry.Level)

	// 添加颜色字段
	colorizedEntry.Fields["_color"] = color
	colorizedEntry.Fields["_level_colored"] = levelText

	return colorizedEntry
}

// getLevelColor 获取级别对应的颜色
func (cl *ConsoleLogger) getLevelColor(level Level) string {
	switch level {
	case TraceLevel:
		return ColorGray
	case DebugLevel:
		return ColorCyan
	case InfoLevel:
		return ColorGreen
	case WarnLevel:
		return ColorYellow
	case ErrorLevel:
		return ColorRed
	case FatalLevel:
		return ColorBrightRed + StyleBold
	case PanicLevel:
		return BgRed + ColorBrightWhite + StyleBold
	default:
		return ColorWhite
	}
}

// getColorizedLevel 获取彩色的级别文本
func (cl *ConsoleLogger) getColorizedLevel(level Level) string {
	color := cl.getLevelColor(level)
	levelText := level.String()
	return color + levelText + ColorReset
}

// extractContextFields 从上下文提取字段
func (cl *ConsoleLogger) extractContextFields(ctx context.Context) Fields {
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

// SetColorized 设置是否启用颜色
func (cl *ConsoleLogger) SetColorized(colorized bool) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.colorized = colorized
}

// IsColorized 检查是否启用颜色
func (cl *ConsoleLogger) IsColorized() bool {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()
	return cl.colorized
}

// SetWriter 设置输出流
func (cl *ConsoleLogger) SetWriter(writer io.Writer) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.writer = writer
}

// GetWriter 获取输出流
func (cl *ConsoleLogger) GetWriter() io.Writer {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()
	return cl.writer
}

// 便捷方法

// Colorize 为文本添加颜色
func Colorize(text, color string) string {
	if !supportsColor() {
		return text
	}
	return color + text + ColorReset
}

// ColorizeLevel 为级别添加颜色
func ColorizeLevel(level Level) string {
	if !supportsColor() {
		return level.String()
	}

	var color string
	switch level {
	case TraceLevel:
		color = ColorGray
	case DebugLevel:
		color = ColorCyan
	case InfoLevel:
		color = ColorGreen
	case WarnLevel:
		color = ColorYellow
	case ErrorLevel:
		color = ColorRed
	case FatalLevel:
		color = ColorBrightRed + StyleBold
	case PanicLevel:
		color = BgRed + ColorBrightWhite + StyleBold
	default:
		color = ColorWhite
	}

	return color + level.String() + ColorReset
}

// StripColors 移除颜色代码
func StripColors(text string) string {
	// 简单的颜色代码移除实现
	// 在实际项目中可能需要更复杂的正则表达式
	result := text
	colors := []string{
		ColorReset, ColorRed, ColorGreen, ColorYellow, ColorBlue,
		ColorPurple, ColorCyan, ColorWhite, ColorGray,
		ColorBrightRed, ColorBrightGreen, ColorBrightYellow,
		ColorBrightBlue, ColorBrightPurple, ColorBrightCyan, ColorBrightWhite,
		BgRed, BgGreen, BgYellow, BgBlue, BgPurple, BgCyan, BgWhite,
		StyleBold, StyleDim, StyleItalic, StyleUnderline,
		StyleBlink, StyleReverse, StyleStrike,
	}

	for _, color := range colors {
		result = replaceAll(result, color, "")
	}

	return result
}

// replaceAll 简单的字符串替换实现
func replaceAll(s, old, new string) string {
	if old == new || old == "" {
		return s
	}

	result := ""
	for len(s) >= len(old) {
		if s[:len(old)] == old {
			result += new
			s = s[len(old):]
		} else {
			result += s[:1]
			s = s[1:]
		}
	}
	return result + s
}