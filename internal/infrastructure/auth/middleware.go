package auth

import (
	"context"
	"errors"
	"net/http"

	// "strings"
	"time"

	"github.com/gin-gonic/gin"

	"greatestworks/internal/infrastructure/logging"
)

// AuthMiddleware 认证中间�?
type AuthMiddleware struct {
	jwtService *JWTService
	logger     logger.Logger
}

// NewAuthMiddleware 创建认证中间�?
func NewAuthMiddleware(jwtService *JWTService, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// RequireAuth HTTP认证中间�?
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		token := m.extractTokenFromRequest(c.Request)
		if token == "" {
			m.logger.Warn("Missing authorization token", "path", c.Request.URL.Path, "method", c.Request.Method)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Missing authorization token",
			})
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.logger.Warn("Token validation failed", "error", err, "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 将用户信息添加到上下�?
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("player_id", claims.PlayerID)
		c.Set("session_id", claims.SessionID)
		c.Set("user_claims", claims)

		m.logger.Debug("User authenticated", "user_id", claims.UserID, "path", c.Request.URL.Path)
		c.Next()
	}
}

// RequireRole 角色验证中间�?
func (m *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行认�?
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// 获取用户角色
		userRole, exists := c.Get("role")
		if !exists {
			m.logger.Error("User role not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "User role not found",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			m.logger.Error("Invalid user role type")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Invalid user role",
			})
			c.Abort()
			return
		}

		// 检查角色权限（管理员拥有所有权限）
		if role != requiredRole && role != "admin" {
			userID, _ := c.Get("user_id")
			m.logger.Warn("Insufficient permissions", "user_id", userID, "required_role", requiredRole, "user_role", role)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 管理员验证中间件
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// OptionalAuth 可选认证中间件
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		token := m.extractTokenFromRequest(c.Request)
		if token == "" {
			// 没有令牌，继续处理但不设置用户信�?
			c.Next()
			return
		}

		// 验证令牌
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			// 令牌无效，记录警告但继续处理
			m.logger.Warn("Optional auth token validation failed", "error", err)
			c.Next()
			return
		}

		// 将用户信息添加到上下�?
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("player_id", claims.PlayerID)
		c.Set("session_id", claims.SessionID)
		c.Set("user_claims", claims)
		c.Set("authenticated", true)

		m.logger.Debug("Optional auth successful", "user_id", claims.UserID)
		c.Next()
	}
}

// RateLimitByUser 按用户限流中间件
func (m *AuthMiddleware) RateLimitByUser(requestsPerMinute int) gin.HandlerFunc {
	// 简化的限流实现，实际应该使用Redis或其他存�?
	userRequests := make(map[string][]time.Time)

	return func(c *gin.Context) {
		// 先进行认�?
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userIDStr := userID.(string)
		now := time.Now()
		oneMinuteAgo := now.Add(-time.Minute)

		// 清理过期的请求记�?
		if requests, exists := userRequests[userIDStr]; exists {
			validRequests := make([]time.Time, 0)
			for _, reqTime := range requests {
				if reqTime.After(oneMinuteAgo) {
					validRequests = append(validRequests, reqTime)
				}
			}
			userRequests[userIDStr] = validRequests
		}

		// 检查是否超过限�?
		if len(userRequests[userIDStr]) >= requestsPerMinute {
			m.logger.Warn("Rate limit exceeded", "user_id", userIDStr, "requests", len(userRequests[userIDStr]))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		// 记录当前请求
		userRequests[userIDStr] = append(userRequests[userIDStr], now)
		c.Next()
	}
}

// extractTokenFromRequest 从请求中提取令牌
func (m *AuthMiddleware) extractTokenFromRequest(r *http.Request) string {
	// 从Authorization头获�?
	auth := r.Header.Get("Authorization")
	if auth != "" {
		return m.jwtService.ExtractTokenFromBearer(auth)
	}

	// 从查询参数获�?
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	// 从Cookie获取
	if cookie, err := r.Cookie("access_token"); err == nil {
		return cookie.Value
	}

	return ""
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(c *gin.Context) (*UserClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*UserClaims)
	return userClaims, ok
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetUserRoleFromContext 从上下文获取用户角色
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}

	roleStr, ok := role.(string)
	return roleStr, ok
}

// IsAuthenticated 检查是否已认证
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// IsAdmin 检查是否是管理�?
func IsAdmin(c *gin.Context) bool {
	role, exists := GetUserRoleFromContext(c)
	return exists && role == "admin"
}

// HasRole 检查是否具有指定角�?
func HasRole(c *gin.Context, requiredRole string) bool {
	role, exists := GetUserRoleFromContext(c)
	if !exists {
		return false
	}
	return role == requiredRole || role == "admin"
}

// TCPAuthValidator TCP认证验证�?
type TCPAuthValidator struct {
	jwtService *JWTService
	logger     logger.Logger
}

// NewTCPAuthValidator 创建TCP认证验证�?
func NewTCPAuthValidator(jwtService *JWTService, logger logger.Logger) *TCPAuthValidator {
	return &TCPAuthValidator{
		jwtService: jwtService,
		logger:     logger,
	}
}

// ValidateToken 验证TCP令牌
func (v *TCPAuthValidator) ValidateToken(token string) (*UserClaims, error) {
	return v.jwtService.ValidateToken(token)
}

// ValidateTokenWithContext 带上下文的令牌验�?
func (v *TCPAuthValidator) ValidateTokenWithContext(ctx context.Context, token string) (*UserClaims, error) {
	// 检查上下文是否已取�?
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	claims, err := v.jwtService.ValidateToken(token)
	if err != nil {
		v.logger.Warn("TCP token validation failed", "error", err)
		return nil, err
	}

	v.logger.Debug("TCP token validated", "user_id", claims.UserID, "session_id", claims.SessionID)
	return claims, nil
}

// CreateAuthContext 创建认证上下�?
func (v *TCPAuthValidator) CreateAuthContext(ctx context.Context, claims *UserClaims) context.Context {
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "username", claims.Username)
	ctx = context.WithValue(ctx, "email", claims.Email)
	ctx = context.WithValue(ctx, "role", claims.Role)
	ctx = context.WithValue(ctx, "player_id", claims.PlayerID)
	ctx = context.WithValue(ctx, "session_id", claims.SessionID)
	ctx = context.WithValue(ctx, "user_claims", claims)
	return ctx
}

// GetUserFromTCPContext 从TCP上下文获取用户信�?
func GetUserFromTCPContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value("user_claims").(*UserClaims)
	return claims, ok
}

// GetUserIDFromTCPContext 从TCP上下文获取用户ID
func GetUserIDFromTCPContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}

// GetUserRoleFromTCPContext 从TCP上下文获取用户角�?
func GetUserRoleFromTCPContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value("role").(string)
	return role, ok
}

// RequireTCPRole 检查TCP上下文中的用户角�?
func RequireTCPRole(ctx context.Context, requiredRole string) error {
	role, exists := GetUserRoleFromTCPContext(ctx)
	if !exists {
		return errors.New("user role not found")
	}

	if role != requiredRole && role != "admin" {
		return errors.New("insufficient permissions")
	}

	return nil
}
