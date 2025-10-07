package auth

import (
	"net/http"
	"time"

	"greatestworks/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtService *JWTService
	logger     logging.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtService *JWTService, logger logging.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			m.logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// 验证token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.logger.Warn("Invalid token", logging.Fields{
				"error": err,
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("expires_at", claims.ExpiresAt)

		m.logger.Debug("User authenticated", logging.Fields{
			"user_id":  claims.UserID,
			"username": claims.Username,
		})
		c.Next()
	}
}

// OptionalAuth 可选认证的中间件
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}

		// 移除Bearer前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// 验证token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.logger.Debug("Invalid token in optional auth", logging.Fields{
				"error": err,
			})
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("expires_at", claims.ExpiresAt)

		m.logger.Debug("User authenticated (optional)", logging.Fields{
			"user_id":  claims.UserID,
			"username": claims.Username,
		})
		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先检查是否已认证
		userID, exists := c.Get("user_id")
		if !exists {
			m.logger.Warn("User not authenticated for role check")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户角色
		userRole, err := m.getUserRole(userID.(string))
		if err != nil {
			m.logger.Error("Failed to get user role", err, logging.Fields{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user role"})
			c.Abort()
			return
		}

		if userRole != role {
			m.logger.Warn("Insufficient permissions", logging.Fields{
				"user_id":       userID,
				"required_role": role,
				"user_role":     userRole,
			})
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		m.logger.Debug("Role check passed", logging.Fields{
			"user_id": userID,
			"role":    role,
		})
		c.Next()
	}
}

// RequireAnyRole 需要任意一个角色的中间件
func (m *AuthMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先检查是否已认证
		userID, exists := c.Get("user_id")
		if !exists {
			m.logger.Warn("User not authenticated for role check")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户角色
		userRole, err := m.getUserRole(userID.(string))
		if err != nil {
			m.logger.Error("Failed to get user role", err, logging.Fields{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user role"})
			c.Abort()
			return
		}

		// 检查用户是否有任意一个所需角色
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			m.logger.Warn("Insufficient permissions", logging.Fields{
				"user_id":        userID,
				"required_roles": roles,
				"user_role":      userRole,
			})
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		m.logger.Debug("Role check passed", logging.Fields{
			"user_id": userID,
			"roles":   roles,
		})
		c.Next()
	}
}

// RateLimit 限流中间件
func (m *AuthMiddleware) RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	// 这里应该实现一个简单的内存限流器
	// 实际项目中应该使用Redis等外部存储
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 这里应该检查限流逻辑
		// 简化实现，直接通过
		m.logger.Debug("Rate limit check", logging.Fields{
			"client_ip": clientIP,
			"requests":  requests,
			"window":    window,
		})
		c.Next()
	}
}

// CORS 跨域中间件
func (m *AuthMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger 请求日志中间件
func (m *AuthMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		m.logger.Info("HTTP Request", logging.Fields{
			"status":    statusCode,
			"latency":   latency,
			"client_ip": clientIP,
			"method":    method,
			"path":      path,
			"body_size": bodySize,
		})
	}
}

// 私有方法

// getUserRole 获取用户角色
func (m *AuthMiddleware) getUserRole(userID string) (string, error) {
	// 这里应该从数据库或缓存中获取用户角色
	// 简化实现，返回默认角色
	return "user", nil
}
