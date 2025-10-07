package auth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"greatestworks/application/handlers"
	playerQueries "greatestworks/application/queries/player"
	"greatestworks/internal/infrastructure/logging"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	logger     logging.Logger
	jwtSecret  string
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// NewLoginHandler 创建登录处理器
func NewLoginHandler(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logging.Logger, jwtSecret string) *LoginHandler {
	return &LoginHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
		jwtSecret:  jwtSecret,
	}
}

// Login 处理登录请求
func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", logging.Fields{
			"error": err,
		})
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// 验证用户凭据
	user, err := h.authenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Warn("Authentication failed", logging.Fields{
			"error":    err,
			"username": req.Username,
		})
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成JWT令牌
	token, expiresAt, err := h.generateToken(user.ID, user.Username)
	if err != nil {
		h.logger.Error("Failed to generate token", err, logging.Fields{
			"user_id": user.ID,
		})
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	// 返回登录响应
	response := LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	h.logger.Info("User logged in successfully", logging.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	})
	c.JSON(200, response)
}

// Logout 处理登出请求
func (h *LoginHandler) Logout(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	// 这里可以实现令牌黑名单逻辑
	h.logger.Info("User logged out", logging.Fields{
		"user_id": userID,
	})
	c.JSON(200, gin.H{"message": "Logged out successfully"})
}

// RefreshToken 刷新令牌
func (h *LoginHandler) RefreshToken(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	username, _ := c.Get("username")

	// 生成新的令牌
	token, expiresAt, err := h.generateToken(userID.(string), username.(string))
	if err != nil {
		h.logger.Error("Failed to refresh token", err, logging.Fields{
			"user_id": userID,
		})
		c.JSON(500, gin.H{"error": "Failed to refresh token"})
		return
	}

	response := LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	h.logger.Info("Token refreshed", logging.Fields{
		"user_id": userID,
	})
	c.JSON(200, response)
}

// 私有方法

// authenticateUser 验证用户凭据
func (h *LoginHandler) authenticateUser(ctx context.Context, username, password string) (*UserInfo, error) {
	// 这里应该实现实际的用户认证逻辑
	// 简化实现，实际项目中应该查询数据库

	// 模拟用户查询
	query := &playerQueries.GetPlayerQuery{
		PlayerID: "user_123", // 使用PlayerID而不是Username
	}

	// 这里应该调用查询总线
	// result, err := h.queryBus.Execute(ctx, query)
	// 简化实现
	_ = query

	// 模拟用户信息
	user := &UserInfo{
		ID:       "user_123",
		Username: username,
		Email:    username + "@example.com",
	}

	return user, nil
}

// generateToken 生成JWT令牌
func (h *LoginHandler) generateToken(userID, username string) (string, time.Time, error) {
	// 设置过期时间
	expiresAt := time.Now().Add(24 * time.Hour)

	// 创建声明
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
