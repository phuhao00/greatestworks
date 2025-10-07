package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"greatestworks/internal/infrastructure/logging"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtSecret string
	logger    logging.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtSecret string, logger logging.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing authorization header")
			c.JSON(401, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.logger.Warn("Invalid authorization header format")
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			m.logger.Warn("Invalid token", logging.Fields{
				"error": err,
			})
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 检查token是否过期
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				m.logger.Warn("Token expired")
				c.JSON(401, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		}

		// 将用户信息存储到上下文中
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if exp, ok := claims["exp"].(float64); ok {
			c.Set("expires_at", time.Unix(int64(exp), 0))
		}

		m.logger.Debug("User authenticated", logging.Fields{
			"user_id":  claims["user_id"],
			"username": claims["username"],
		})
		c.Next()
	}
}

// OptionalAuth 可选认证的中间件
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			m.logger.Debug("Invalid token in optional auth", logging.Fields{
				"error": err,
			})
			c.Next()
			return
		}

		// 检查token是否过期
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				m.logger.Debug("Token expired in optional auth")
				c.Next()
				return
			}
		}

		// 将用户信息存储到上下文中
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if exp, ok := claims["exp"].(float64); ok {
			c.Set("expires_at", time.Unix(int64(exp), 0))
		}

		m.logger.Debug("User authenticated (optional)", logging.Fields{
			"user_id":  claims["user_id"],
			"username": claims["username"],
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
			c.JSON(401, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户角色
		userRole, err := m.getUserRole(userID.(string))
		if err != nil {
			m.logger.Error("Failed to get user role", err, logging.Fields{
				"user_id": userID,
			})
			c.JSON(500, gin.H{"error": "Failed to check user role"})
			c.Abort()
			return
		}

		if userRole != role {
			m.logger.Warn("Insufficient permissions", logging.Fields{
				"user_id":       userID,
				"required_role": role,
				"user_role":     userRole,
			})
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
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
			c.JSON(401, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户角色
		userRole, err := m.getUserRole(userID.(string))
		if err != nil {
			m.logger.Error("Failed to get user role", err, logging.Fields{
				"user_id": userID,
			})
			c.JSON(500, gin.H{"error": "Failed to check user role"})
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
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
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

// 私有方法

// validateToken 验证JWT token
func (m *AuthMiddleware) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// getUserRole 获取用户角色
func (m *AuthMiddleware) getUserRole(userID string) (string, error) {
	// 这里应该从数据库或缓存中获取用户角色
	// 简化实现，返回默认角色
	return "user", nil
}
