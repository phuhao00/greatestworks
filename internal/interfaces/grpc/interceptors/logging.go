package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"greatestworks/internal/infrastructure/logger"
)

// LoggingUnaryInterceptor 日志拦截器（一元RPC）
func LoggingUnaryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 获取请求信息
		requestInfo := extractRequestInfo(ctx, info.FullMethod)

		// 记录请求开始
		logger.Info("gRPC request started",
			"method", info.FullMethod,
			"user_id", requestInfo.UserID,
			"client_ip", requestInfo.ClientIP,
			"user_agent", requestInfo.UserAgent,
			"request_id", requestInfo.RequestID)

		// 执行处理器
		resp, err := handler(ctx, req)

		// 计算处理时间
		duration := time.Since(start)

		// 获取状态码
		statusCode := codes.OK
		errorMessage := ""
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
				errorMessage = st.Message()
			} else {
				statusCode = codes.Internal
				errorMessage = err.Error()
			}
		}

		// 记录请求完成
		logFields := []interface{}{
			"method", info.FullMethod,
			"status_code", statusCode.String(),
			"duration", duration.String(),
			"duration_ms", duration.Milliseconds(),
			"user_id", requestInfo.UserID,
			"client_ip", requestInfo.ClientIP,
			"request_id", requestInfo.RequestID,
		}

		if err != nil {
			logFields = append(logFields, "error", errorMessage)
			logger.Error("gRPC request failed", logFields...)
		} else {
			logger.Info("gRPC request completed", logFields...)
		}

		// 记录慢请求
		if duration > 5*time.Second {
			logger.Warn("Slow gRPC request detected",
				"method", info.FullMethod,
				"duration", duration.String(),
				"user_id", requestInfo.UserID)
		}

		return resp, err
	}
}

// LoggingStreamInterceptor 日志拦截器（流式RPC）
func LoggingStreamInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// 获取请求信息
		requestInfo := extractRequestInfo(stream.Context(), info.FullMethod)

		// 记录流开始
		logger.Info("gRPC stream started",
			"method", info.FullMethod,
			"is_client_stream", info.IsClientStream,
			"is_server_stream", info.IsServerStream,
			"user_id", requestInfo.UserID,
			"client_ip", requestInfo.ClientIP,
			"request_id", requestInfo.RequestID)

		// 创建包装的流用于统计
		wrappedStream := &loggingServerStream{
			ServerStream: stream,
			logger:       logger,
			requestInfo:  requestInfo,
			methodName:   info.FullMethod,
			messagesSent: 0,
			messagesRecv: 0,
		}

		// 执行处理器
		err := handler(srv, wrappedStream)

		// 计算处理时间
		duration := time.Since(start)

		// 获取状态码
		statusCode := codes.OK
		errorMessage := ""
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
				errorMessage = st.Message()
			} else {
				statusCode = codes.Internal
				errorMessage = err.Error()
			}
		}

		// 记录流完成
		logFields := []interface{}{
			"method", info.FullMethod,
			"status_code", statusCode.String(),
			"duration", duration.String(),
			"duration_ms", duration.Milliseconds(),
			"messages_sent", wrappedStream.messagesSent,
			"messages_recv", wrappedStream.messagesRecv,
			"user_id", requestInfo.UserID,
			"client_ip", requestInfo.ClientIP,
			"request_id", requestInfo.RequestID,
		}

		if err != nil {
			logFields = append(logFields, "error", errorMessage)
			logger.Error("gRPC stream failed", logFields...)
		} else {
			logger.Info("gRPC stream completed", logFields...)
		}

		// 记录慢流
		if duration > 30*time.Second {
			logger.Warn("Long-running gRPC stream detected",
				"method", info.FullMethod,
				"duration", duration.String(),
				"user_id", requestInfo.UserID)
		}

		return err
	}
}

// RequestInfo 请求信息
type RequestInfo struct {
	UserID    string
	ClientIP  string
	UserAgent string
	RequestID string
	TraceID   string
}

// extractRequestInfo 提取请求信息
func extractRequestInfo(ctx context.Context, method string) *RequestInfo {
	info := &RequestInfo{
		UserID:    GetUserIDFromContext(ctx),
		RequestID: generateRequestID(),
	}

	// 从metadata中获取客户端信息
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// 获取客户端IP
		if clientIPs := md.Get("x-forwarded-for"); len(clientIPs) > 0 {
			info.ClientIP = clientIPs[0]
		} else if clientIPs := md.Get("x-real-ip"); len(clientIPs) > 0 {
			info.ClientIP = clientIPs[0]
		} else if clientIPs := md.Get("remote-addr"); len(clientIPs) > 0 {
			info.ClientIP = clientIPs[0]
		}

		// 获取User-Agent
		if userAgents := md.Get("user-agent"); len(userAgents) > 0 {
			info.UserAgent = userAgents[0]
		}

		// 获取请求ID
		if requestIDs := md.Get("x-request-id"); len(requestIDs) > 0 {
			info.RequestID = requestIDs[0]
		}

		// 获取追踪ID
		if traceIDs := md.Get("x-trace-id"); len(traceIDs) > 0 {
			info.TraceID = traceIDs[0]
		}
	}

	return info
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单的请求ID生成（实际应该使用UUID或其他唯一标识符）
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// loggingServerStream 包装的服务器流用于日志记录
type loggingServerStream struct {
	grpc.ServerStream
	logger       logger.Logger
	requestInfo  *RequestInfo
	methodName   string
	messagesSent int
	messagesRecv int
}

// SendMsg 发送消息
func (s *loggingServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.messagesSent++
		s.logger.Debug("gRPC stream message sent",
			"method", s.methodName,
			"messages_sent", s.messagesSent,
			"user_id", s.requestInfo.UserID,
			"request_id", s.requestInfo.RequestID)
	} else {
		s.logger.Error("Failed to send gRPC stream message",
			"method", s.methodName,
			"error", err,
			"user_id", s.requestInfo.UserID,
			"request_id", s.requestInfo.RequestID)
	}
	return err
}

// RecvMsg 接收消息
func (s *loggingServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.messagesRecv++
		s.logger.Debug("gRPC stream message received",
			"method", s.methodName,
			"messages_recv", s.messagesRecv,
			"user_id", s.requestInfo.UserID,
			"request_id", s.requestInfo.RequestID)
	} else if err.Error() != "EOF" {
		s.logger.Error("Failed to receive gRPC stream message",
			"method", s.methodName,
			"error", err,
			"user_id", s.requestInfo.UserID,
			"request_id", s.requestInfo.RequestID)
	}
	return err
}

// LogMethodCall 记录方法调用（用于手动日志记录）
func LogMethodCall(ctx context.Context, logger logger.Logger, method string, args ...interface{}) {
	requestInfo := extractRequestInfo(ctx, method)

	logFields := []interface{}{
		"method", method,
		"user_id", requestInfo.UserID,
		"client_ip", requestInfo.ClientIP,
		"request_id", requestInfo.RequestID,
	}

	// 添加额外参数
	logFields = append(logFields, args...)

	logger.Info("Method called", logFields...)
}

// LogMethodResult 记录方法结果（用于手动日志记录）
func LogMethodResult(ctx context.Context, logger logger.Logger, method string, duration time.Duration, err error, args ...interface{}) {
	requestInfo := extractRequestInfo(ctx, method)

	logFields := []interface{}{
		"method", method,
		"duration", duration.String(),
		"duration_ms", duration.Milliseconds(),
		"user_id", requestInfo.UserID,
		"client_ip", requestInfo.ClientIP,
		"request_id", requestInfo.RequestID,
	}

	// 添加额外参数
	logFields = append(logFields, args...)

	if err != nil {
		logFields = append(logFields, "error", err.Error())
		logger.Error("Method failed", logFields...)
	} else {
		logger.Info("Method completed", logFields...)
	}
}

// GetRequestID 从上下文获取请求ID
func GetRequestID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if requestIDs := md.Get("x-request-id"); len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}
	return ""
}

// GetTraceID 从上下文获取追踪ID
func GetTraceID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if traceIDs := md.Get("x-trace-id"); len(traceIDs) > 0 {
			return traceIDs[0]
		}
	}
	return ""
}

// AddRequestIDToContext 将请求ID添加到上下文
func AddRequestIDToContext(ctx context.Context, requestID string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set("x-request-id", requestID)
	return metadata.NewIncomingContext(ctx, md)
}

// AddTraceIDToContext 将追踪ID添加到上下文
func AddTraceIDToContext(ctx context.Context, traceID string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set("x-trace-id", traceID)
	return metadata.NewIncomingContext(ctx, md)
}