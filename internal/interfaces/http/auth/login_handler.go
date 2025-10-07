package auth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	playerQueries "greatestworks/application/queries/player"
	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	logger     logger.Logger
	jwtSecret  string
}

// NewLoginHandler 创建登录处理器
func NewLoginHandler(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logger.Logger, jwtSecret string) *LoginHandler {
	return &LoginHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
		jwtSecret:  jwtSecret,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`
	PlayerID  string    `json:"player_id"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// JWTClaims JWT声明
type JWTClaims struct {
	PlayerID string `json:"player_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Login 用户登录
func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	ctx := context.Background()

	// 验证用户凭据（这里简化处理，实际应该查询用户表）
	query := &playerQueries.GetPlayerByUsernameQuery{Username: req.Username}
	result, err := handlers.ExecuteQueryTyped[*playerQueries.GetPlayerByUsernameQuery, *playerQueries.GetPlayerByUsernameResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player by username", "error", err, "username", req.Username)
		c.JSON(401, gin.H{"error": "Invalid credentials", "success": false})
		return
	}

	if !result.Found {
		h.logger.Warn("Player not found", "username", req.Username)
		c.JSON(401, gin.H{"error": "Invalid credentials", "success": false})
		return
	}

	// 验证密码（这里简化处理，实际应该使用bcrypt等加密）
	if !h.validatePassword(req.Password, result.Player.PasswordHash) {
		h.logger.Warn("Invalid password", "username", req.Username)
		c.JSON(401, gin.H{"error": "Invalid credentials", "success": false})
		return
	}

	// 生成JWT令牌
	token, expiresAt, err := h.generateJWT(result.Player.ID, req.Username, "player")
	if err != nil {
		h.logger.Error("Failed to generate JWT token", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error", "success": false})
		return
	}

	// 记录登录日志
	h.logger.Info("User logged in successfully", "username", req.Username, "player_id", result.Player.ID)

	// 返回响应
	response := &LoginResponse{
		Token:     token,
		PlayerID:  result.Player.ID,
		Username:  req.Username,
		ExpiresAt: expiresAt,
	}

	c.JSON(200, gin.H{"data": response, "success": true, "message": "Login successful"})
}

// RefreshToken 刷新令牌
func (h *LoginHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid refresh token request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	// 验证现有令牌
	claims, err := h.validateJWT(req.Token)
	if err != nil {
		h.logger.Error("Invalid JWT token", "error", err)
		c.JSON(401, gin.H{"error": "Invalid token", "success": false})
		return
	}

	// 生成新令牌
	newToken, expiresAt, err := h.generateJWT(claims.PlayerID, claims.Username, claims.Role)
	if err != nil {
		h.logger.Error("Failed to generate new JWT token", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error", "success": false})
		return
	}

	// 返回新令牌
	response := &LoginResponse{
		Token:     newToken,
		PlayerID:  claims.PlayerID,
		Username:  claims.Username,
		ExpiresAt: expiresAt,
	}

	c.JSON(200, gin.H{"data": response, "success": true, "message": "Token refreshed successfully"})
}

// Logout 用户登出
func (h *LoginHandler) Logout(c *gin.Context) {
	// 从请求头获取令牌
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(400, gin.H{"error": "Authorization header required", "success": false})
		return
	}

	// 移除Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 验证令牌
	claims, err := h.validateJWT(token)
	if err != nil {
		h.logger.Error("Invalid JWT token on logout", "error", err)
		c.JSON(401, gin.H{"error": "Invalid token", "success": false})
		return
	}

	// 记录登出日志
	h.logger.Info("User logged out", "username", claims.Username, "player_id", claims.PlayerID)

	// 在实际应用中，这里应该将令牌加入黑名单
	// TODO: 实现令牌黑名单机制

	c.JSON(200, gin.H{"success": true, "message": "Logout successful"})
}

// 私有方法

// generateJWT 生成JWT令牌
func (h *LoginHandler) generateJWT(playerID, username, role string) (string, time.Time, error) {
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
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// validateJWT 验证JWT令牌
func (h *LoginHandler) validateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// validatePassword 验证密码
func (h *LoginHandler) validatePassword(password, hash string) bool {
	// 这里简化处理，实际应该使用bcrypt.CompareHashAndPassword
	// 临时实现：直接比较明文（仅用于演示）
	return password == hash
}