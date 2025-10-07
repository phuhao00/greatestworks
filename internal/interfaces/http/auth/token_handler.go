package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"greatestworks/internal/infrastructure/logging"
)

// TokenHandler Token管理处理器
type TokenHandler struct {
	jwtSecret string
	logger    logging.Logger
}

// NewTokenHandler 创建Token处理器
func NewTokenHandler(jwtSecret string, logger logging.Logger) *TokenHandler {
	return &TokenHandler{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// GenerateToken 生成JWT令牌
func (h *TokenHandler) GenerateToken(userID, username string, expiresIn time.Duration) (string, time.Time, error) {
	// 计算过期时间
	expiresAt := time.Now().Add(expiresIn)

	// 创建声明
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "greatestworks",
		"aud":      "greatestworks-users",
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		h.logger.Error("Failed to sign token", err, logging.Fields{
			"user_id": userID,
		})
		return "", time.Time{}, err
	}

	h.logger.Debug("Token generated", logging.Fields{
		"user_id":    userID,
		"username":   username,
		"expires_at": expiresAt,
	})
	return tokenString, expiresAt, nil
}

// ValidateToken 验证JWT令牌
func (h *TokenHandler) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		h.logger.Warn("Token validation failed", logging.Fields{
			"error": err,
		})
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		h.logger.Debug("Token validated successfully", logging.Fields{
			"user_id": claims["user_id"],
		})
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshToken 刷新令牌
func (h *TokenHandler) RefreshToken(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User not authenticated for token refresh")
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	username, _ := c.Get("username")

	// 生成新的令牌
	token, expiresAt, err := h.GenerateToken(userID.(string), username.(string), 24*time.Hour)
	if err != nil {
		h.logger.Error("Failed to refresh token", err, logging.Fields{
			"user_id": userID,
		})
		c.JSON(500, gin.H{"error": "Failed to refresh token"})
		return
	}

	response := gin.H{
		"token":      token,
		"expires_at": expiresAt,
		"type":       "Bearer",
	}

	h.logger.Info("Token refreshed", logging.Fields{
		"user_id": userID,
	})
	c.JSON(200, response)
}

// RevokeToken 撤销令牌
func (h *TokenHandler) RevokeToken(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User not authenticated for token revocation")
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	// 这里应该实现令牌黑名单逻辑
	// 简化实现，实际项目中应该将令牌添加到黑名单
	h.logger.Info("Token revoked", logging.Fields{
		"user_id": userID,
	})
	c.JSON(200, gin.H{"message": "Token revoked successfully"})
}

// GetTokenInfo 获取令牌信息
func (h *TokenHandler) GetTokenInfo(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User not authenticated for token info")
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	username, _ := c.Get("username")
	expiresAt, _ := c.Get("expires_at")

	response := gin.H{
		"user_id":    userID,
		"username":   username,
		"expires_at": expiresAt,
	}

	c.JSON(200, response)
}

// ValidateTokenMiddleware 验证令牌中间件
func (h *TokenHandler) ValidateTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			h.logger.Warn("Missing authorization header")
			c.JSON(401, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			h.logger.Warn("Invalid authorization header format")
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证token
		claims, err := h.ValidateToken(tokenString)
		if err != nil {
			h.logger.Warn("Token validation failed", logging.Fields{
				"error": err,
			})
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
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

		h.logger.Debug("Token validated", logging.Fields{
			"user_id": claims["user_id"],
		})
		c.Next()
	}
}
