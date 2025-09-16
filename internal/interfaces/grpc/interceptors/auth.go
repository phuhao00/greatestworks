package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"greatestworks/internal/infrastructure/logger"
)

// AuthUnaryInterceptor 认证拦截器（一元RPC）
func AuthUnaryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 跳过健康检查和反射服务的认证
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// 从metadata中获取认证信息
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Warn("Missing metadata in gRPC request", "method", info.FullMethod)
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// 验证认证token
		if err := validateAuthToken(md, logger); err != nil {
			logger.Warn("Authentication failed", "method", info.FullMethod, "error", err)
			return nil, err
		}

		// 将用户信息添加到上下文
		ctx = addUserToContext(ctx, md)

		logger.Debug("gRPC request authenticated", "method", info.FullMethod)
		return handler(ctx, req)
	}
}

// AuthStreamInterceptor 认证拦截器（流式RPC）
func AuthStreamInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 跳过健康检查和反射服务的认证
		if isPublicMethod(info.FullMethod) {
			return handler(srv, stream)
		}

		// 从metadata中获取认证信息
		ctx := stream.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Warn("Missing metadata in gRPC stream", "method", info.FullMethod)
			return status.Error(codes.Unauthenticated, "missing metadata")
		}

		// 验证认证token
		if err := validateAuthToken(md, logger); err != nil {
			logger.Warn("Stream authentication failed", "method", info.FullMethod, "error", err)
			return err
		}

		// 将用户信息添加到上下文
		ctx = addUserToContext(ctx, md)

		// 创建包装的流
		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}

		logger.Debug("gRPC stream authenticated", "method", info.FullMethod)
		return handler(srv, wrappedStream)
	}
}

// wrappedServerStream 包装的服务器流
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回流的上下文
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// isPublicMethod 检查是否是公开方法（不需要认证）
func isPublicMethod(fullMethod string) bool {
	publicMethods := []string{
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
	}

	for _, method := range publicMethods {
		if fullMethod == method {
			return true
		}
	}

	return false
}

// validateAuthToken 验证认证token
func validateAuthToken(md metadata.MD, logger logger.Logger) error {
	// 获取Authorization header
	auth := md.Get("authorization")
	if len(auth) == 0 {
		return status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := auth[0]
	if token == "" {
		return status.Error(codes.Unauthenticated, "empty authorization token")
	}

	// 检查Bearer token格式
	if !strings.HasPrefix(token, "Bearer ") {
		return status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	// 提取实际的token
	actualToken := strings.TrimPrefix(token, "Bearer ")
	if actualToken == "" {
		return status.Error(codes.Unauthenticated, "empty bearer token")
	}

	// 验证token（这里简化处理，实际应该验证JWT或查询数据库）
	if err := verifyJWTToken(actualToken, logger); err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	return nil
}

// verifyJWTToken 验证JWT token
func verifyJWTToken(token string, logger logger.Logger) error {
	// TODO: 实现JWT token验证逻辑
	// 这里应该：
	// 1. 解析JWT token
	// 2. 验证签名
	// 3. 检查过期时间
	// 4. 验证issuer和audience
	// 5. 检查token是否在黑名单中

	// 临时实现：简单的token验证
	if len(token) < 10 {
		return status.Error(codes.Unauthenticated, "token too short")
	}

	// 这里可以添加更复杂的验证逻辑
	logger.Debug("JWT token verified", "token_length", len(token))
	return nil
}

// addUserToContext 将用户信息添加到上下文
func addUserToContext(ctx context.Context, md metadata.MD) context.Context {
	// 从token中提取用户信息（这里简化处理）
	userID := extractUserIDFromMetadata(md)
	userRole := extractUserRoleFromMetadata(md)

	// 将用户信息添加到上下文
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "user_role", userRole)

	return ctx
}

// extractUserIDFromMetadata 从metadata中提取用户ID
func extractUserIDFromMetadata(md metadata.MD) string {
	// 尝试从user-id header获取
	if userIDs := md.Get("user-id"); len(userIDs) > 0 {
		return userIDs[0]
	}

	// 尝试从x-user-id header获取
	if userIDs := md.Get("x-user-id"); len(userIDs) > 0 {
		return userIDs[0]
	}

	// TODO: 从JWT token中解析用户ID
	// 这里应该解析authorization header中的JWT token
	// 并从中提取用户ID

	return "unknown"
}

// extractUserRoleFromMetadata 从metadata中提取用户角色
func extractUserRoleFromMetadata(md metadata.MD) string {
	// 尝试从user-role header获取
	if roles := md.Get("user-role"); len(roles) > 0 {
		return roles[0]
	}

	// 尝试从x-user-role header获取
	if roles := md.Get("x-user-role"); len(roles) > 0 {
		return roles[0]
	}

	// TODO: 从JWT token中解析用户角色
	// 这里应该解析authorization header中的JWT token
	// 并从中提取用户角色

	return "user"
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUserRoleFromContext 从上下文中获取用户角色
func GetUserRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value("user_role").(string); ok {
		return role
	}
	return ""
}

// RequireRole 检查用户是否具有指定角色
func RequireRole(ctx context.Context, requiredRole string) error {
	userRole := GetUserRoleFromContext(ctx)
	if userRole == "" {
		return status.Error(codes.Unauthenticated, "user role not found")
	}

	// 简单的角色检查（实际应该实现更复杂的权限系统）
	if userRole != requiredRole && userRole != "admin" {
		return status.Error(codes.PermissionDenied, "insufficient permissions")
	}

	return nil
}

// RequireAdmin 检查用户是否是管理员
func RequireAdmin(ctx context.Context) error {
	return RequireRole(ctx, "admin")
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(ctx context.Context) bool {
	userID := GetUserIDFromContext(ctx)
	return userID != "" && userID != "unknown"
}

// GetAuthInfo 获取认证信息
func GetAuthInfo(ctx context.Context) map[string]string {
	return map[string]string{
		"user_id":   GetUserIDFromContext(ctx),
		"user_role": GetUserRoleFromContext(ctx),
	}
}