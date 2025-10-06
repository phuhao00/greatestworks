// Package logger 日志系统
// Author: MMO Server Team
// Created: 2024

package logger

// "context" // 未使用
// "fmt" // 未使用

// Logger 日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
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
	}
}
