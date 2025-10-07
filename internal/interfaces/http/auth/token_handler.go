package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"greatestworks/internal/infrastructure/logging"
)

// TokenHandler Token管理处理�?
type TokenHandler struct {
	jwtSecret string
	logger    logger.Logger
}

// NewTokenHandler 创建Token处理�?
func NewTokenHandler(jwtSecret string, logger logger.Logger) *TokenHandler {
	return &TokenHandler{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// ValidateTokenRequest 验证令牌请求
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ValidateTokenResponse 验证令牌响应
type ValidateTokenResponse struct {
	Valid     bool      `json:"valid"`
	PlayerID  string    `json:"player_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Role      string    `json:"role,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Message   string    `json:"message,omitempty"`
}

// RevokeTokenRequest 撤销令牌请求
type RevokeTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// TokenInfoResponse 令牌信息响应
type TokenInfoResponse struct {
	PlayerID  string    `json:"player_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Issuer    string    `json:"issuer"`
}

// ValidateToken 验证令牌
func (h *TokenHandler) ValidateToken(c *gin.Context) {
	var req ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid validate token request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	claims, err := h.validateJWT(req.Token)
	if err != nil {
		h.logger.Debug("Token validation failed", "error", err)
		response := &ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid or expired token",
		}
		c.JSON(200, gin.H{"data": response, "success": true})
		return
	}

	// 检查令牌是否即将过�?
	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	message := "Token is valid"
	if timeUntilExpiry < time.Hour {
		message = "Token is valid but will expire soon"
	}

	response := &ValidateTokenResponse{
		Valid:     true,
		PlayerID:  claims.PlayerID,
		Username:  claims.Username,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
		Message:   message,
	}

	c.JSON(200, gin.H{"data": response, "success": true})
}

// RefreshToken 刷新令牌
func (h *TokenHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid refresh token request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	// 验证现有令牌
	claims, err := h.validateJWT(req.Token)
	if err != nil {
		h.logger.Error("Invalid JWT token for refresh", "error", err)
		c.JSON(401, gin.H{"error": "Invalid token", "success": false})
		return
	}

	// 检查令牌是否可以刷新（距离过期时间不超�?天）
	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	if timeUntilExpiry > 7*24*time.Hour {
		h.logger.Warn("Token refresh attempted too early", "player_id", claims.PlayerID, "time_until_expiry", timeUntilExpiry)
		c.JSON(400, gin.H{"error": "Token does not need refresh yet", "success": false})
		return
	}

	// 生成新令�?
	newToken, expiresAt, err := h.generateJWT(claims.PlayerID, claims.Username, claims.Role)
	if err != nil {
		h.logger.Error("Failed to generate new JWT token", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error", "success": false})
		return
	}

	// 记录令牌刷新日志
	h.logger.Info("Token refreshed successfully", "player_id", claims.PlayerID, "username", claims.Username)

	// 返回新令�?
	response := &LoginResponse{
		Token:     newToken,
		PlayerID:  claims.PlayerID,
		Username:  claims.Username,
		ExpiresAt: expiresAt,
	}

	c.JSON(200, gin.H{"data": response, "success": true, "message": "Token refreshed successfully"})
}

// RevokeToken 撤销令牌
func (h *TokenHandler) RevokeToken(c *gin.Context) {
	var req RevokeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid revoke token request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	// 验证令牌
	claims, err := h.validateJWT(req.Token)
	if err != nil {
		h.logger.Debug("Revoke attempt with invalid token", "error", err)
		// 即使令牌无效，也返回成功（避免信息泄露）
		c.JSON(200, gin.H{"success": true, "message": "Token revoked successfully"})
		return
	}

	// 在实际应用中，这里应该将令牌加入黑名�?
	// TODO: 实现令牌黑名单机制（Redis存储�?
	h.addToBlacklist(req.Token, claims.ExpiresAt.Time)

	// 记录撤销日志
	h.logger.Info("Token revoked", "player_id", claims.PlayerID, "username", claims.Username)

	c.JSON(200, gin.H{"success": true, "message": "Token revoked successfully"})
}

// GetTokenInfo 获取令牌信息
func (h *TokenHandler) GetTokenInfo(c *gin.Context) {
	// 从当前上下文获取令牌信息
	claims, exists := c.Get("token_claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Authentication required", "success": false})
		return
	}

	jwtClaims, ok := claims.(*JWTClaims)
	if !ok {
		h.logger.Error("Invalid token claims type")
		c.JSON(500, gin.H{"error": "Internal server error", "success": false})
		return
	}

	response := &TokenInfoResponse{
		PlayerID:  jwtClaims.PlayerID,
		Username:  jwtClaims.Username,
		Role:      jwtClaims.Role,
		IssuedAt:  jwtClaims.IssuedAt.Time,
		ExpiresAt: jwtClaims.ExpiresAt.Time,
		Issuer:    jwtClaims.Issuer,
	}

	c.JSON(200, gin.H{"data": response, "success": true})
}

// RevokeAllTokens 撤销用户的所有令�?
func (h *TokenHandler) RevokeAllTokens(c *gin.Context) {
	// 获取当前用户信息
	claims, exists := c.Get("token_claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Authentication required", "success": false})
		return
	}

	jwtClaims, ok := claims.(*JWTClaims)
	if !ok {
		h.logger.Error("Invalid token claims type")
		c.JSON(500, gin.H{"error": "Internal server error", "success": false})
		return
	}

	// 在实际应用中，这里应该撤销用户的所有令�?
	// TODO: 实现用户令牌全部撤销机制
	h.revokeAllUserTokens(jwtClaims.PlayerID)

	// 记录撤销日志
	h.logger.Info("All tokens revoked for user", "player_id", jwtClaims.PlayerID, "username", jwtClaims.Username)

	c.JSON(200, gin.H{"success": true, "message": "All tokens revoked successfully"})
}

// 私有方法

// validateJWT 验证JWT令牌
func (h *TokenHandler) validateJWT(tokenString string) (*JWTClaims, error) {
	// 首先检查令牌是否在黑名单中
	if h.isTokenBlacklisted(tokenString) {
		return nil, jwt.ErrTokenInvalidClaims
	}

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

// generateJWT 生成JWT令牌
func (h *TokenHandler) generateJWT(playerID, username, role string) (string, time.Time, error) {
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

// addToBlacklist 将令牌加入黑名单
func (h *TokenHandler) addToBlacklist(token string, expiresAt time.Time) {
	// TODO: 实现Redis黑名单存�?
	// 这里应该将令牌存储到Redis中，设置过期时间为令牌的过期时间
	h.logger.Debug("Token added to blacklist", "token_hash", h.hashToken(token), "expires_at", expiresAt)
}

// isTokenBlacklisted 检查令牌是否在黑名单中
func (h *TokenHandler) isTokenBlacklisted(token string) bool {
	// TODO: 实现Redis黑名单检�?
	// 这里应该从Redis中检查令牌是否存�?
	return false // 临时返回false
}

// revokeAllUserTokens 撤销用户的所有令�?
func (h *TokenHandler) revokeAllUserTokens(playerID string) {
	// TODO: 实现用户令牌全部撤销
	// 可以通过在Redis中设置用户的令牌版本号来实现
	h.logger.Debug("All tokens revoked for user", "player_id", playerID)
}

// hashToken 对令牌进行哈希处理（用于存储�?
func (h *TokenHandler) hashToken(token string) string {
	// TODO: 实现令牌哈希
	// 可以使用SHA256等哈希算�?
	return "hashed_" + token[:10] // 临时实现
}