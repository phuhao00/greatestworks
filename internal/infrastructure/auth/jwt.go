package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"greatestworks/internal/infrastructure/logger"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string
	Issuer           string
	Audience         string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	SigningMethod    jwt.SigningMethod
}

// DefaultJWTConfig 默认JWT配置
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:           generateRandomSecret(),
		Issuer:           "greatestworks",
		Audience:         "greatestworks-users",
		AccessTokenTTL:   1 * time.Hour,
		RefreshTokenTTL:  24 * time.Hour,
		SigningMethod:    jwt.SigningMethodHS256,
	}
}

// UserClaims 用户声明
type UserClaims struct {
	UserID    string            `json:"user_id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Role      string            `json:"role"`
	PlayerID  string            `json:"player_id,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// JWTService JWT服务
type JWTService struct {
	config *JWTConfig
	logger logger.Logger
}

// NewJWTService 创建JWT服务
func NewJWTService(config *JWTConfig, logger logger.Logger) *JWTService {
	if config == nil {
		config = DefaultJWTConfig()
	}

	return &JWTService{
		config: config,
		logger: logger,
	}
}

// GenerateTokenPair 生成令牌对
func (s *JWTService) GenerateTokenPair(userID, username, email, role string, metadata map[string]string) (*TokenPair, error) {
	now := time.Now()
	accessExpiresAt := now.Add(s.config.AccessTokenTTL)
	refreshExpiresAt := now.Add(s.config.RefreshTokenTTL)

	// 生成会话ID
	sessionID := s.generateSessionID()

	// 创建访问令牌声明
	accessClaims := &UserClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Role:      role,
		SessionID: sessionID,
		Metadata:  metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Audience:  jwt.ClaimStrings{s.config.Audience},
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        s.generateJTI(),
		},
	}

	// 创建刷新令牌声明
	refreshClaims := &UserClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Audience:  jwt.ClaimStrings{s.config.Audience},
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        s.generateJTI(),
		},
	}

	// 生成访问令牌
	accessToken := jwt.NewWithClaims(s.config.SigningMethod, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		s.logger.Error("Failed to sign access token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken := jwt.NewWithClaims(s.config.SigningMethod, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		s.logger.Error("Failed to sign refresh token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	tokenPair := &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.AccessTokenTTL.Seconds()),
		ExpiresAt:    accessExpiresAt,
	}

	s.logger.Info("Token pair generated", "user_id", userID, "session_id", sessionID, "expires_at", accessExpiresAt)
	return tokenPair, nil
}

// ValidateToken 验证令牌
func (s *JWTService) ValidateToken(tokenString string) (*UserClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if token.Method != s.config.SigningMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		s.logger.Warn("Token validation failed", "error", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 检查令牌是否有效
	if !token.Valid {
		s.logger.Warn("Token is not valid")
		return nil, errors.New("token is not valid")
	}

	// 提取声明
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		s.logger.Error("Failed to extract claims from token")
		return nil, errors.New("failed to extract claims")
	}

	// 验证声明
	if err := s.validateClaims(claims); err != nil {
		s.logger.Warn("Claims validation failed", "error", err, "user_id", claims.UserID)
		return nil, err
	}

	s.logger.Debug("Token validated successfully", "user_id", claims.UserID, "session_id", claims.SessionID)
	return claims, nil
}

// RefreshToken 刷新令牌
func (s *JWTService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// 验证刷新令牌
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		s.logger.Warn("Refresh token validation failed", "error", err)
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 生成新的令牌对
	newTokenPair, err := s.GenerateTokenPair(
		claims.UserID,
		claims.Username,
		claims.Email,
		claims.Role,
		claims.Metadata,
	)
	if err != nil {
		s.logger.Error("Failed to generate new token pair", "error", err, "user_id", claims.UserID)
		return nil, err
	}

	s.logger.Info("Token refreshed", "user_id", claims.UserID, "old_session_id", claims.SessionID)
	return newTokenPair, nil
}

// ExtractTokenFromBearer 从Bearer字符串中提取令牌
func (s *JWTService) ExtractTokenFromBearer(bearerToken string) string {
	const bearerPrefix = "Bearer "
	if len(bearerToken) > len(bearerPrefix) && bearerToken[:len(bearerPrefix)] == bearerPrefix {
		return bearerToken[len(bearerPrefix):]
	}
	return bearerToken
}

// GetTokenInfo 获取令牌信息（不验证有效性）
func (s *JWTService) GetTokenInfo(tokenString string) (*UserClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &UserClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	return claims, nil
}

// IsTokenExpired 检查令牌是否过期
func (s *JWTService) IsTokenExpired(tokenString string) bool {
	claims, err := s.GetTokenInfo(tokenString)
	if err != nil {
		return true
	}

	if claims.ExpiresAt == nil {
		return true
	}

	return claims.ExpiresAt.Time.Before(time.Now())
}

// RevokeToken 撤销令牌（这里需要配合黑名单实现）
func (s *JWTService) RevokeToken(tokenString string) error {
	claims, err := s.GetTokenInfo(tokenString)
	if err != nil {
		return err
	}

	// TODO: 将令牌添加到黑名单
	// 这里可以使用Redis或数据库存储被撤销的令牌

	s.logger.Info("Token revoked", "user_id", claims.UserID, "jti", claims.ID)
	return nil
}

// validateClaims 验证声明
func (s *JWTService) validateClaims(claims *UserClaims) error {
	now := time.Now()

	// 验证过期时间
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(now) {
		return errors.New("token has expired")
	}

	// 验证生效时间
	if claims.NotBefore != nil && claims.NotBefore.Time.After(now) {
		return errors.New("token is not yet valid")
	}

	// 验证签发者
	if claims.Issuer != s.config.Issuer {
		return fmt.Errorf("invalid issuer: expected %s, got %s", s.config.Issuer, claims.Issuer)
	}

	// 验证受众
	validAudience := false
	for _, aud := range claims.Audience {
		if aud == s.config.Audience {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return fmt.Errorf("invalid audience: expected %s", s.config.Audience)
	}

	// 验证必需字段
	if claims.UserID == "" {
		return errors.New("user_id is required")
	}
	if claims.Username == "" {
		return errors.New("username is required")
	}
	if claims.Role == "" {
		return errors.New("role is required")
	}

	return nil
}

// generateSessionID 生成会话ID
func (s *JWTService) generateSessionID() string {
	return fmt.Sprintf("sess_%d_%s", time.Now().Unix(), s.generateRandomString(8))
}

// generateJTI 生成JWT ID
func (s *JWTService) generateJTI() string {
	return fmt.Sprintf("jti_%d_%s", time.Now().UnixNano(), s.generateRandomString(16))
}

// generateRandomString 生成随机字符串
func (s *JWTService) generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateRandomSecret 生成随机密钥
func generateRandomSecret() string {
	bytes := make([]byte, 32) // 256 bits
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// UpdateConfig 更新配置
func (s *JWTService) UpdateConfig(config *JWTConfig) {
	if config != nil {
		s.config = config
		s.logger.Info("JWT configuration updated")
	}
}

// GetConfig 获取配置
func (s *JWTService) GetConfig() *JWTConfig {
	// 返回配置副本，隐藏密钥
	return &JWTConfig{
		Secret:           "[HIDDEN]",
		Issuer:           s.config.Issuer,
		Audience:         s.config.Audience,
		AccessTokenTTL:   s.config.AccessTokenTTL,
		RefreshTokenTTL:  s.config.RefreshTokenTTL,
		SigningMethod:    s.config.SigningMethod,
	}
}

// GetStats 获取统计信息
func (s *JWTService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"issuer":             s.config.Issuer,
		"audience":           s.config.Audience,
		"access_token_ttl":   s.config.AccessTokenTTL.String(),
		"refresh_token_ttl":  s.config.RefreshTokenTTL.String(),
		"signing_method":     s.config.SigningMethod.Alg(),
	}
}