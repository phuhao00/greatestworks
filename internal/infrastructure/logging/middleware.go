// Package logging 日志中间件
// Author: MMO Server Team
// Created: 2024

package logging

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// HTTPMiddleware HTTP日志中间件
type HTTPMiddleware struct {
	logger Logger
	config *MiddlewareConfig
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// 是否记录请求体
	LogRequestBody bool
	// 是否记录响应体
	LogResponseBody bool
	// 请求体最大长度
	MaxRequestBodySize int64
	// 响应体最大长度
	MaxResponseBodySize int64
	// 跳过的路径
	SkipPaths []string
	// 跳过的方法
	SkipMethods []string
	// 是否记录请求头
	LogHeaders bool
	// 敏感头部（不记录）
	SensitiveHeaders []string
	// 是否记录查询参数
	LogQueryParams bool
	// 敏感查询参数（不记录）
	SensitiveParams []string
	// 慢请求阈值
	SlowRequestThreshold time.Duration
	// 是否启用请求ID
	EnableRequestID bool
	// 请求ID头部名称
	RequestIDHeader string
	// 是否记录用户代理
	LogUserAgent bool
	// 是否记录客户端IP
	LogClientIP bool
	// 是否记录Referer
	LogReferer bool
}

// DefaultMiddlewareConfig 默认中间件配置
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		LogRequestBody:       false,
		LogResponseBody:      false,
		MaxRequestBodySize:   1024 * 1024, // 1MB
		MaxResponseBodySize:  1024 * 1024, // 1MB
		SkipPaths:           []string{"/health", "/metrics", "/favicon.ico"},
		SkipMethods:         []string{"OPTIONS"},
		LogHeaders:          true,
		SensitiveHeaders:    []string{"Authorization", "Cookie", "Set-Cookie", "X-Auth-Token"},
		LogQueryParams:      true,
		SensitiveParams:     []string{"password", "token", "secret", "key"},
		SlowRequestThreshold: 1 * time.Second,
		EnableRequestID:     true,
		RequestIDHeader:     "X-Request-ID",
		LogUserAgent:       true,
		LogClientIP:        true,
		LogReferer:         true,
	}
}

// NewHTTPMiddleware 创建HTTP日志中间件
func NewHTTPMiddleware(logger Logger, config *MiddlewareConfig) *HTTPMiddleware {
	if config == nil {
		config = DefaultMiddlewareConfig()
	}

	return &HTTPMiddleware{
		logger: logger,
		config: config,
	}
}

// Handler 中间件处理器
func (hm *HTTPMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查是否跳过
		if hm.shouldSkip(r) {
			next.ServeHTTP(w, r)
			return
		}

		// 生成请求ID
		requestID := hm.getOrGenerateRequestID(r)

		// 创建带有请求ID的上下文
		ctx := context.WithValue(r.Context(), FieldRequestID, requestID)
		r = r.WithContext(ctx)

		// 设置响应头
		if hm.config.EnableRequestID && hm.config.RequestIDHeader != "" {
			w.Header().Set(hm.config.RequestIDHeader, requestID)
		}

		// 记录请求开始
		start := time.Now()
		hm.logRequestStart(r, requestID)

		// 创建响应写入器包装器
		wrapper := &responseWriter{
			ResponseWriter: w,
			status:         200,
			size:           0,
		}

		// 处理panic
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				hm.logPanic(r, requestID, err, stack)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// 执行下一个处理器
		next.ServeHTTP(wrapper, r)

		// 记录请求完成
		duration := time.Since(start)
		hm.logRequestComplete(r, requestID, wrapper.status, wrapper.size, duration)
	})
}

// shouldSkip 检查是否应该跳过日志记录
func (hm *HTTPMiddleware) shouldSkip(r *http.Request) bool {
	// 检查路径
	for _, path := range hm.config.SkipPaths {
		if r.URL.Path == path {
			return true
		}
	}

	// 检查方法
	for _, method := range hm.config.SkipMethods {
		if r.Method == method {
			return true
		}
	}

	return false
}

// getOrGenerateRequestID 获取或生成请求ID
func (hm *HTTPMiddleware) getOrGenerateRequestID(r *http.Request) string {
	if !hm.config.EnableRequestID {
		return ""
	}

	// 尝试从头部获取
	if hm.config.RequestIDHeader != "" {
		if requestID := r.Header.Get(hm.config.RequestIDHeader); requestID != "" {
			return requestID
		}
	}

	// 生成新的请求ID
	return uuid.New().String()
}

// logRequestStart 记录请求开始
func (hm *HTTPMiddleware) logRequestStart(r *http.Request, requestID string) {
	fields := Fields{
		FieldRequestID: requestID,
		FieldMethod:    r.Method,
		FieldURL:       r.URL.String(),
		"path":         r.URL.Path,
		"proto":        r.Proto,
		"remote_addr":  r.RemoteAddr,
	}

	// 添加客户端IP
	if hm.config.LogClientIP {
		if clientIP := hm.getClientIP(r); clientIP != "" {
			fields[FieldIP] = clientIP
		}
	}

	// 添加用户代理
	if hm.config.LogUserAgent {
		if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
			fields[FieldUserAgent] = userAgent
		}
	}

	// 添加Referer
	if hm.config.LogReferer {
		if referer := r.Header.Get("Referer"); referer != "" {
			fields["referer"] = referer
		}
	}

	// 添加头部信息
	if hm.config.LogHeaders {
		headers := hm.filterHeaders(r.Header)
		if len(headers) > 0 {
			fields["headers"] = headers
		}
	}

	// 添加查询参数
	if hm.config.LogQueryParams {
		params := hm.filterQueryParams(r.URL.Query())
		if len(params) > 0 {
			fields["query_params"] = params
		}
	}

	// 添加请求体（如果启用）
	if hm.config.LogRequestBody && r.ContentLength > 0 && r.ContentLength <= hm.config.MaxRequestBodySize {
		// 注意：这里需要小心处理请求体，避免消费掉原始请求体
		// 在实际实现中，可能需要使用 io.TeeReader 或类似技术
		fields["content_length"] = r.ContentLength
		fields["content_type"] = r.Header.Get("Content-Type")
	}

	hm.logger.InfoWithFields("HTTP request started", fields)
}

// logRequestComplete 记录请求完成
func (hm *HTTPMiddleware) logRequestComplete(r *http.Request, requestID string, status, size int, duration time.Duration) {
	fields := Fields{
		FieldRequestID: requestID,
		FieldMethod:    r.Method,
		FieldURL:       r.URL.String(),
		"path":         r.URL.Path,
		FieldStatus:    status,
		"size":         size,
		FieldDuration:  duration,
		"duration_ms":  float64(duration.Nanoseconds()) / 1e6,
	}

	// 根据状态码和持续时间选择日志级别
	level := hm.getLogLevel(status, duration)
	message := fmt.Sprintf("HTTP request completed - %s %s %d", r.Method, r.URL.Path, status)

	// 记录日志
	switch level {
	case DebugLevel:
		hm.logger.DebugWithFields(message, fields)
	case InfoLevel:
		hm.logger.InfoWithFields(message, fields)
	case WarnLevel:
		hm.logger.WarnWithFields(message, fields)
	case ErrorLevel:
		hm.logger.ErrorWithFields(message, fields)
	default:
		hm.logger.InfoWithFields(message, fields)
	}
}

// logPanic 记录panic
func (hm *HTTPMiddleware) logPanic(r *http.Request, requestID string, err interface{}, stack []byte) {
	fields := Fields{
		FieldRequestID:  requestID,
		FieldMethod:     r.Method,
		FieldURL:        r.URL.String(),
		"path":          r.URL.Path,
		"panic":         err,
		FieldStackTrace: string(stack),
	}

	hm.logger.ErrorWithFields("HTTP request panicked", fields)
}

// getLogLevel 根据状态码和持续时间获取日志级别
func (hm *HTTPMiddleware) getLogLevel(status int, duration time.Duration) Level {
	// 错误状态码
	if status >= 500 {
		return ErrorLevel
	}
	if status >= 400 {
		return WarnLevel
	}

	// 慢请求
	if duration > hm.config.SlowRequestThreshold {
		return WarnLevel
	}

	// 正常请求
	return InfoLevel
}

// getClientIP 获取客户端IP
func (hm *HTTPMiddleware) getClientIP(r *http.Request) string {
	// 尝试从各种头部获取真实IP
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Client-IP",
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",   // Akamai
	}

	for _, header := range headers {
		if ip := r.Header.Get(header); ip != "" {
			// X-Forwarded-For 可能包含多个IP，取第一个
			if header == "X-Forwarded-For" {
				if ips := strings.Split(ip, ","); len(ips) > 0 {
					return strings.TrimSpace(ips[0])
				}
			}
			return ip
		}
	}

	// 回退到RemoteAddr
	if ip := r.RemoteAddr; ip != "" {
		// 移除端口号
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			return ip[:idx]
		}
		return ip
	}

	return ""
}

// filterHeaders 过滤敏感头部
func (hm *HTTPMiddleware) filterHeaders(headers http.Header) map[string]string {
	filtered := make(map[string]string)

	for name, values := range headers {
		// 检查是否为敏感头部
		if hm.isSensitiveHeader(name) {
			filtered[name] = "[REDACTED]"
		} else if len(values) > 0 {
			filtered[name] = values[0]
		}
	}

	return filtered
}

// filterQueryParams 过滤敏感查询参数
func (hm *HTTPMiddleware) filterQueryParams(params map[string][]string) map[string]string {
	filtered := make(map[string]string)

	for name, values := range params {
		// 检查是否为敏感参数
		if hm.isSensitiveParam(name) {
			filtered[name] = "[REDACTED]"
		} else if len(values) > 0 {
			filtered[name] = values[0]
		}
	}

	return filtered
}

// isSensitiveHeader 检查是否为敏感头部
func (hm *HTTPMiddleware) isSensitiveHeader(name string) bool {
	name = strings.ToLower(name)
	for _, sensitive := range hm.config.SensitiveHeaders {
		if strings.ToLower(sensitive) == name {
			return true
		}
	}
	return false
}

// isSensitiveParam 检查是否为敏感参数
func (hm *HTTPMiddleware) isSensitiveParam(name string) bool {
	name = strings.ToLower(name)
	for _, sensitive := range hm.config.SensitiveParams {
		if strings.ToLower(sensitive) == name {
			return true
		}
	}
	return false
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// WriteHeader 写入状态码
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Write 写入数据
func (rw *responseWriter) Write(data []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(data)
	rw.size += n
	return n, err
}

// Status 获取状态码
func (rw *responseWriter) Status() int {
	return rw.status
}

// Size 获取响应大小
func (rw *responseWriter) Size() int {
	return rw.size
}

// GameMiddleware 游戏日志中间件
type GameMiddleware struct {
	logger Logger
}

// NewGameMiddleware 创建游戏日志中间件
func NewGameMiddleware(logger Logger) *GameMiddleware {
	return &GameMiddleware{
		logger: logger,
	}
}

// LogPlayerAction 记录玩家行为
func (gm *GameMiddleware) LogPlayerAction(ctx context.Context, playerID string, action string, details Fields) {
	fields := Fields{
		FieldUserID:   playerID,
		FieldModule:   "game",
		"action":      action,
		"action_time": time.Now(),
	}

	// 合并详细信息
	if details != nil {
		fields.Merge(details)
	}

	// 从上下文提取请求ID
	if requestID := ctx.Value(FieldRequestID); requestID != nil {
		fields[FieldRequestID] = requestID
	}

	gm.logger.InfoWithFields(fmt.Sprintf("Player action: %s", action), fields)
}

// LogGameEvent 记录游戏事件
func (gm *GameMiddleware) LogGameEvent(ctx context.Context, eventType string, eventData Fields) {
	fields := Fields{
		FieldModule:   "game",
		"event_type":  eventType,
		"event_time":  time.Now(),
	}

	// 合并事件数据
	if eventData != nil {
		fields.Merge(eventData)
	}

	// 从上下文提取请求ID
	if requestID := ctx.Value(FieldRequestID); requestID != nil {
		fields[FieldRequestID] = requestID
	}

	gm.logger.InfoWithFields(fmt.Sprintf("Game event: %s", eventType), fields)
}

// LogPerformanceMetric 记录性能指标
func (gm *GameMiddleware) LogPerformanceMetric(ctx context.Context, metricName string, value interface{}, unit string) {
	fields := Fields{
		FieldModule:    "performance",
		"metric_name":  metricName,
		"metric_value": value,
		"metric_unit":  unit,
		"metric_time":  time.Now(),
	}

	// 从上下文提取请求ID
	if requestID := ctx.Value(FieldRequestID); requestID != nil {
		fields[FieldRequestID] = requestID
	}

	gm.logger.InfoWithFields(fmt.Sprintf("Performance metric: %s = %v %s", metricName, value, unit), fields)
}

// 便捷函数

// WithRequestID 为上下文添加请求ID
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, FieldRequestID, requestID)
}

// WithUserID 为上下文添加用户ID
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, FieldUserID, userID)
}

// WithSessionID 为上下文添加会话ID
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, FieldSessionID, sessionID)
}

// WithTraceID 为上下文添加跟踪ID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, FieldTraceID, traceID)
}

// GetRequestID 从上下文获取请求ID
func GetRequestID(ctx context.Context) string {
	if requestID := ctx.Value(FieldRequestID); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) string {
	if userID := ctx.Value(FieldUserID); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetSessionID 从上下文获取会话ID
func GetSessionID(ctx context.Context) string {
	if sessionID := ctx.Value(FieldSessionID); sessionID != nil {
		if id, ok := sessionID.(string); ok {
			return id
		}
	}
	return ""
}

// GetTraceID 从上下文获取跟踪ID
func GetTraceID(ctx context.Context) string {
	if traceID := ctx.Value(FieldTraceID); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}