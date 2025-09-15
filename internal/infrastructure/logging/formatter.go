// Package logging 日志格式化器
// Author: MMO Server Team
// Created: 2024

package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// JSONFormatter JSON格式化器
type JSONFormatter struct {
	TimestampFormat string
	PrettyPrint     bool
	FieldMap        map[string]string
	DisableHTMLEscape bool
}

// NewJSONFormatter 创建JSON格式化器
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
		FieldMap:        make(map[string]string),
		DisableHTMLEscape: true,
	}
}

// Format 格式化日志条目为JSON
func (jf *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	// 创建数据映射
	data := make(map[string]interface{})

	// 添加基础字段
	data[jf.getFieldKey("timestamp")] = entry.Timestamp.Format(jf.TimestampFormat)
	data[jf.getFieldKey("level")] = entry.Level.String()
	data[jf.getFieldKey("message")] = entry.Message

	// 添加调用者信息
	if entry.Caller != "" {
		data[jf.getFieldKey("caller")] = entry.Caller
	}

	// 添加错误信息
	if entry.Error != nil {
		data[jf.getFieldKey("error")] = entry.Error.Error()
	}

	// 添加自定义字段
	for k, v := range entry.Fields {
		// 跳过内部字段
		if strings.HasPrefix(k, "_") {
			continue
		}
		data[k] = v
	}

	// 序列化为JSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(!jf.DisableHTMLEscape)

	if jf.PrettyPrint {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}

	return buf.Bytes(), nil
}

// getFieldKey 获取字段键名（支持字段映射）
func (jf *JSONFormatter) getFieldKey(key string) string {
	if mapped, exists := jf.FieldMap[key]; exists {
		return mapped
	}
	return key
}

// TextFormatter 文本格式化器
type TextFormatter struct {
	TimestampFormat string
	FullTimestamp   bool
	DisableColors   bool
	DisableQuote    bool
	QuoteEmptyFields bool
	FieldMap        map[string]string
	CallerPrettyfier func(*Entry) (function string, file string)
	SortingFunc     func([]string)
}

// NewTextFormatter 创建文本格式化器
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		DisableColors:   false,
		DisableQuote:    false,
		QuoteEmptyFields: true,
		FieldMap:        make(map[string]string),
		SortingFunc:     sort.Strings,
	}
}

// Format 格式化日志条目为文本
func (tf *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	// 添加时间戳
	if tf.FullTimestamp {
		buf.WriteString(entry.Timestamp.Format(tf.TimestampFormat))
		buf.WriteString(" ")
	}

	// 添加级别
	levelText := entry.Level.String()
	if !tf.DisableColors {
		if coloredLevel, exists := entry.Fields["_level_colored"]; exists {
			levelText = fmt.Sprintf("%v", coloredLevel)
		} else {
			levelText = ColorizeLevel(entry.Level)
		}
	}
	buf.WriteString(fmt.Sprintf("[%s]", levelText))
	buf.WriteString(" ")

	// 添加调用者信息
	if entry.Caller != "" {
		if tf.CallerPrettyfier != nil {
			function, file := tf.CallerPrettyfier(entry)
			if function != "" {
				buf.WriteString(fmt.Sprintf("%s() ", function))
			}
			if file != "" {
				buf.WriteString(fmt.Sprintf("%s ", file))
			}
		} else {
			buf.WriteString(fmt.Sprintf("%s ", entry.Caller))
		}
	}

	// 添加消息
	buf.WriteString(entry.Message)

	// 添加字段
	if len(entry.Fields) > 0 {
		buf.WriteString(" ")
		tf.appendFields(&buf, entry)
	}

	// 添加错误信息
	if entry.Error != nil {
		buf.WriteString(fmt.Sprintf(" error=%q", entry.Error.Error()))
	}

	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// appendFields 添加字段到缓冲区
func (tf *TextFormatter) appendFields(buf *bytes.Buffer, entry *Entry) {
	// 收集字段键
	keys := make([]string, 0, len(entry.Fields))
	for k := range entry.Fields {
		// 跳过内部字段
		if strings.HasPrefix(k, "_") {
			continue
		}
		keys = append(keys, k)
	}

	// 排序字段
	if tf.SortingFunc != nil {
		tf.SortingFunc(keys)
	}

	// 添加字段
	for i, key := range keys {
		if i > 0 {
			buf.WriteString(" ")
		}

		value := entry.Fields[key]
		mappedKey := tf.getFieldKey(key)

		// 格式化字段
		tf.appendKeyValue(buf, mappedKey, value)
	}
}

// appendKeyValue 添加键值对
func (tf *TextFormatter) appendKeyValue(buf *bytes.Buffer, key string, value interface{}) {
	buf.WriteString(key)
	buf.WriteString("=")

	// 格式化值
	formattedValue := tf.formatValue(value)

	// 决定是否需要引号
	needsQuote := tf.needsQuote(formattedValue)
	if needsQuote {
		buf.WriteString(strconv.Quote(formattedValue))
	} else {
		buf.WriteString(formattedValue)
	}
}

// formatValue 格式化值
func (tf *TextFormatter) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return v.Format(time.RFC3339)
	case time.Duration:
		return v.String()
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%+v", v)
	}
}

// needsQuote 检查是否需要引号
func (tf *TextFormatter) needsQuote(value string) bool {
	if tf.DisableQuote {
		return false
	}

	if value == "" {
		return tf.QuoteEmptyFields
	}

	// 检查是否包含空格或特殊字符
	for _, r := range value {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '"' || r == '\\' {
			return true
		}
	}

	return false
}

// getFieldKey 获取字段键名（支持字段映射）
func (tf *TextFormatter) getFieldKey(key string) string {
	if mapped, exists := tf.FieldMap[key]; exists {
		return mapped
	}
	return key
}

// ConsoleFormatter 控制台格式化器
type ConsoleFormatter struct {
	*TextFormatter
	Colorize bool
}

// NewConsoleFormatter 创建控制台格式化器
func NewConsoleFormatter() *ConsoleFormatter {
	textFormatter := NewTextFormatter()
	textFormatter.DisableColors = false

	return &ConsoleFormatter{
		TextFormatter: textFormatter,
		Colorize:      true,
	}
}

// Format 格式化日志条目为控制台输出
func (cf *ConsoleFormatter) Format(entry *Entry) ([]byte, error) {
	// 如果禁用颜色，使用文本格式化器
	if !cf.Colorize || cf.DisableColors {
		return cf.TextFormatter.Format(entry)
	}

	var buf bytes.Buffer

	// 添加彩色时间戳
	if cf.FullTimestamp {
		timestamp := entry.Timestamp.Format(cf.TimestampFormat)
		buf.WriteString(Colorize(timestamp, ColorGray))
		buf.WriteString(" ")
	}

	// 添加彩色级别
	levelText := ColorizeLevel(entry.Level)
	buf.WriteString(fmt.Sprintf("[%s]", levelText))
	buf.WriteString(" ")

	// 添加调用者信息
	if entry.Caller != "" {
		caller := Colorize(entry.Caller, ColorCyan)
		buf.WriteString(fmt.Sprintf("%s ", caller))
	}

	// 添加消息（根据级别着色）
	messageColor := cf.getMessageColor(entry.Level)
	message := Colorize(entry.Message, messageColor)
	buf.WriteString(message)

	// 添加字段
	if len(entry.Fields) > 0 {
		buf.WriteString(" ")
		cf.appendColorizedFields(&buf, entry)
	}

	// 添加错误信息
	if entry.Error != nil {
		errorText := Colorize(entry.Error.Error(), ColorRed)
		buf.WriteString(fmt.Sprintf(" error=%s", errorText))
	}

	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// getMessageColor 获取消息颜色
func (cf *ConsoleFormatter) getMessageColor(level Level) string {
	switch level {
	case ErrorLevel, FatalLevel, PanicLevel:
		return ColorRed
	case WarnLevel:
		return ColorYellow
	case InfoLevel:
		return ColorGreen
	case DebugLevel:
		return ColorCyan
	case TraceLevel:
		return ColorGray
	default:
		return ColorWhite
	}
}

// appendColorizedFields 添加彩色字段
func (cf *ConsoleFormatter) appendColorizedFields(buf *bytes.Buffer, entry *Entry) {
	// 收集字段键
	keys := make([]string, 0, len(entry.Fields))
	for k := range entry.Fields {
		// 跳过内部字段
		if strings.HasPrefix(k, "_") {
			continue
		}
		keys = append(keys, k)
	}

	// 排序字段
	if cf.SortingFunc != nil {
		cf.SortingFunc(keys)
	}

	// 添加彩色字段
	for i, key := range keys {
		if i > 0 {
			buf.WriteString(" ")
		}

		value := entry.Fields[key]
		mappedKey := cf.getFieldKey(key)

		// 彩色键名
		colorizedKey := Colorize(mappedKey, ColorBlue)
		buf.WriteString(colorizedKey)
		buf.WriteString("=")

		// 格式化并着色值
		formattedValue := cf.formatValue(value)
		colorizedValue := cf.colorizeValue(value, formattedValue)

		// 决定是否需要引号
		needsQuote := cf.needsQuote(formattedValue)
		if needsQuote {
			buf.WriteString(strconv.Quote(colorizedValue))
		} else {
			buf.WriteString(colorizedValue)
		}
	}
}

// colorizeValue 为值添加颜色
func (cf *ConsoleFormatter) colorizeValue(value interface{}, formatted string) string {
	switch value.(type) {
	case string:
		return Colorize(formatted, ColorGreen)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return Colorize(formatted, ColorYellow)
	case float32, float64:
		return Colorize(formatted, ColorYellow)
	case bool:
		if formatted == "true" {
			return Colorize(formatted, ColorGreen)
		}
		return Colorize(formatted, ColorRed)
	case time.Time, time.Duration:
		return Colorize(formatted, ColorCyan)
	case error:
		return Colorize(formatted, ColorRed)
	default:
		return Colorize(formatted, ColorWhite)
	}
}

// CompactFormatter 紧凑格式化器
type CompactFormatter struct {
	TimestampFormat string
	DisableColors   bool
}

// NewCompactFormatter 创建紧凑格式化器
func NewCompactFormatter() *CompactFormatter {
	return &CompactFormatter{
		TimestampFormat: "15:04:05",
		DisableColors:   false,
	}
}

// Format 格式化日志条目为紧凑格式
func (cf *CompactFormatter) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	// 时间戳
	timestamp := entry.Timestamp.Format(cf.TimestampFormat)
	if !cf.DisableColors {
		timestamp = Colorize(timestamp, ColorGray)
	}
	buf.WriteString(timestamp)
	buf.WriteString(" ")

	// 级别（单字符）
	levelChar := cf.getLevelChar(entry.Level)
	if !cf.DisableColors {
		levelChar = cf.colorizeLevelChar(entry.Level, levelChar)
	}
	buf.WriteString(levelChar)
	buf.WriteString(" ")

	// 消息
	message := entry.Message
	if !cf.DisableColors {
		messageColor := cf.getMessageColor(entry.Level)
		message = Colorize(message, messageColor)
	}
	buf.WriteString(message)

	// 重要字段
	if requestID, exists := entry.Fields[FieldRequestID]; exists {
		buf.WriteString(fmt.Sprintf(" [%v]", requestID))
	}

	if userID, exists := entry.Fields[FieldUserID]; exists {
		buf.WriteString(fmt.Sprintf(" user:%v", userID))
	}

	// 错误信息
	if entry.Error != nil {
		errorText := entry.Error.Error()
		if !cf.DisableColors {
			errorText = Colorize(errorText, ColorRed)
		}
		buf.WriteString(fmt.Sprintf(" err:%s", errorText))
	}

	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// getLevelChar 获取级别字符
func (cf *CompactFormatter) getLevelChar(level Level) string {
	switch level {
	case TraceLevel:
		return "T"
	case DebugLevel:
		return "D"
	case InfoLevel:
		return "I"
	case WarnLevel:
		return "W"
	case ErrorLevel:
		return "E"
	case FatalLevel:
		return "F"
	case PanicLevel:
		return "P"
	default:
		return "?"
	}
}

// colorizeLevelChar 为级别字符添加颜色
func (cf *CompactFormatter) colorizeLevelChar(level Level, char string) string {
	switch level {
	case TraceLevel:
		return Colorize(char, ColorGray)
	case DebugLevel:
		return Colorize(char, ColorCyan)
	case InfoLevel:
		return Colorize(char, ColorGreen)
	case WarnLevel:
		return Colorize(char, ColorYellow)
	case ErrorLevel:
		return Colorize(char, ColorRed)
	case FatalLevel:
		return Colorize(char, ColorBrightRed+StyleBold)
	case PanicLevel:
		return Colorize(char, BgRed+ColorBrightWhite+StyleBold)
	default:
		return Colorize(char, ColorWhite)
	}
}

// getMessageColor 获取消息颜色
func (cf *CompactFormatter) getMessageColor(level Level) string {
	switch level {
	case ErrorLevel, FatalLevel, PanicLevel:
		return ColorRed
	case WarnLevel:
		return ColorYellow
	default:
		return ColorWhite
	}
}