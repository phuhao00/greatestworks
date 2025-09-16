package auth

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"greatestworks/internal/infrastructure/logger"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtSecret string
	logger    logger.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtSecret string, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			m.logger.Warn("Missing authorization token", "path", c.Request.URL.Path)
			c.JSON(401, gin.H{"error": "Authorization token required", "success": false})
			c.Abort()
			return
		}

		claims, err := m.validateToken(token)
		if err != nil {
			m.logger.Warn("Invalid authorization token", "error", err, "path", c.Request.URL.Path)
			c.JSON(401, gin.H{"error": "Invalid authorization token", "success": false})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("player_id", claims.PlayerID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token_claims", claims)

		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先检查是否已认证
		claims, exists := c.Get("token_claims")
		if !exists {
			m.logger.Warn("No token claims found in context", "path", c.Request.URL.Path)
			c.JSON(401, gin.H{"error": "Authentication required", "success": false})
			c.Abort()
			return
		}

		jwtClaims, ok := claims.(*JWTClaims)
		if !ok {
			m.logger.Error("Invalid token claims type", "path", c.Request.URL.Path)
			c.JSON(500, gin.H{"error": "Internal server error", "success": false})
			c.Abort()
			return
		}

		// 检查用户角色
		userRole := jwtClaims.Role
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		m.logger.Warn("Insufficient permissions", "user_role", userRole, "required_roles", roles, "path", c.Request.URL.Path)
		c.JSON(403, gin.H{"error": "Insufficient permissions", "success": false})
		c.Abort()
	}
}

// RequireGM GM权限中间件
func (m *AuthMiddleware) RequireGM() gin.HandlerFunc {
	return m.RequireRole("gm", "admin", "super_admin")
}

// RequireAdmin 管理员权限中间件
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin", "super_admin")
}

// RequireSuperAdmin 超级管理员权限中间件
func (m *AuthMiddleware) RequireSuperAdmin() gin.HandlerFunc {
	return m.RequireRole("super_admin")
}

// OptionalAuth 可选认证中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token != "" {
			claims, err := m.validateToken(token)
			if err == nil {
				// 认证成功，存储用户信息
				c.Set("player_id", claims.PlayerID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Set("token_claims", claims)
				c.Set("authenticated", true)
			} else {
				m.logger.Debug("Optional auth failed", "error", err)
				c.Set("authenticated", false)
			}
		} else {
			c.Set("authenticated", false)
		}

		c.Next()
	}
}

// RefreshTokenMiddleware 刷新令牌中间件
func (m *AuthMiddleware) RefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("token_claims")
		if !exists {
			c.Next()
			return
		}

		jwtClaims, ok := claims.(*JWTClaims)
		if !ok {
			c.Next()
			return
		}

		// 检查令牌是否即将过期（剩余时间少于1小时）
		if time.Until(jwtClaims.ExpiresAt.Time) < time.Hour {
			// 生成新令牌
			newToken, expiresAt, err := m.generateJWT(jwtClaims.PlayerID, jwtClaims.Username, jwtClaims.Role)
			if err != nil {
				m.logger.Error("Failed to generate refresh token", "error", err)
			} else {
				// 在响应头中返回新令牌
				c.Header("X-New-Token", newToken)
				c.Header("X-Token-Expires-At", expiresAt.Format(time.RFC3339))
				m.logger.Info("Token refreshed", "player_id", jwtClaims.PlayerID)
			}
		}

		c.Next()
	}
}

// CORSMiddleware CORS中间件
func (m *AuthMiddleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 允许的源列表（在生产环境中应该配置具体的域名）
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://yourdomain.com",
		}

		// 检查是否为允许的源
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed || origin == "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "X-New-Token, X-Token-Expires-At")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 私有方法

// extractToken 从请求中提取令牌
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	// 从Authorization头提取
	auth := c.GetHeader("Authorization")
	if auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
	}

	// 从查询参数提取
	token := c.Query("token")
	if token != "" {
		return token
	}

	// 从Cookie提取
	cookie, err := c.Cookie("auth_token")
	if err == nil && cookie != "" {
		return cookie
	}

	return ""
}

// validateToken 验证令牌
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// generateJWT 生成JWT令牌
func (m *AuthMiddleware) generateJWT(playerID, username, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour) // 24小时过期

	claims := &JWTClaims{
		PlayerID: playerID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "greatestworks",
			Subject:   playerID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GetCurrentUser 获取当前用户信息的辅助函数
func GetCurrentUser(c *gin.Context) (*JWTClaims, bool) {
	claims, exists := c.Get("token_claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*JWTClaims)
	return jwtClaims, ok
}

// IsAuthenticated 检查是否已认证的辅助函数
func IsAuthenticated(c *gin.Context) bool {
	authenticated, exists := c.Get("authenticated")
	if !exists {
		return false
	}

	auth, ok := authenticated.(bool)
	return ok && auth
}

// GetPlayerID 获取当前玩家ID的辅助函数
func GetPlayerID(c *gin.Context) (string, bool) {
	playerID, exists := c.Get("player_id")
	if !exists {
		return "", false
	}

	id, ok := playerID.(string)
	return id, ok
}