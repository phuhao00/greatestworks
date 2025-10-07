package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"greatestworks/internal/infrastructure/logging"
)

// JWTConfig JWT configuration
type JWTConfig struct {
	Secret          string
	Issuer          string
	Audience        string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Algorithm       string
	SigningMethod   jwt.SigningMethod
}

// JWTService JWT service
type JWTService struct {
	config *JWTConfig
	logger logging.Logger
}

// Claims JWT claims
type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
	Audience  string `json:"aud"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(config *JWTConfig, logger logging.Logger) *JWTService {
	if config == nil {
		config = &JWTConfig{
			Secret:          "default-secret",
			Issuer:          "greatestworks",
			Audience:        "greatestworks-users",
			AccessTokenTTL:  24 * time.Hour,
			RefreshTokenTTL: 7 * 24 * time.Hour,
			Algorithm:       "HS256",
			SigningMethod:   jwt.SigningMethodHS256,
		}
	}

	return &JWTService{
		config: config,
		logger: logger,
	}
}

// GenerateToken generates a new JWT token
func (j *JWTService) GenerateToken(userID, username, role string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.config.AccessTokenTTL)

	claims := &Claims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    j.config.Issuer,
		Audience:  j.config.Audience,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Audience:  []string{j.config.Audience},
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(j.config.SigningMethod, claims)
	tokenString, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		j.logger.Error("Failed to sign token", err, logging.Fields{
			"user_id": userID,
		})
		return "", time.Time{}, err
	}

	j.logger.Debug("Token generated", logging.Fields{
		"user_id":    userID,
		"username":   username,
		"expires_at": expiresAt,
	})

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTService) GenerateRefreshToken(userID string) (string, time.Time, error) {
	// Generate a random refresh token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}

	refreshToken := hex.EncodeToString(bytes)
	expiresAt := time.Now().Add(j.config.RefreshTokenTTL)

	j.logger.Debug("Refresh token generated", logging.Fields{
		"user_id":    userID,
		"expires_at": expiresAt,
	})

	return refreshToken, expiresAt, nil
}

// ValidateToken validates a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		j.logger.Warn("Token validation failed", logging.Fields{
			"error": err,
		})
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Verify issuer and audience
		if claims.Issuer != j.config.Issuer {
			return nil, errors.New("invalid issuer")
		}
		if claims.Audience != j.config.Audience {
			return nil, errors.New("invalid audience")
		}

		j.logger.Debug("Token validated successfully", logging.Fields{
			"user_id":  claims.UserID,
			"username": claims.Username,
		})

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken refreshes an access token using a refresh token
func (j *JWTService) RefreshToken(refreshToken, userID string) (string, time.Time, error) {
	// In a real implementation, you would validate the refresh token against a database
	// For now, we'll just generate a new token
	return j.GenerateToken(userID, "", "")
}

// RevokeToken revokes a token (adds it to a blacklist)
func (j *JWTService) RevokeToken(tokenString string) error {
	// In a real implementation, you would add the token to a blacklist
	// For now, we'll just log the revocation
	j.logger.Info("Token revoked", logging.Fields{
		"token": tokenString,
	})
	return nil
}

// IsTokenValid checks if a token is valid without parsing it
func (j *JWTService) IsTokenValid(tokenString string) bool {
	_, err := j.ValidateToken(tokenString)
	return err == nil
}

// GetTokenExpiration gets the expiration time of a token
func (j *JWTService) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(claims.ExpiresAt, 0), nil
}

// GetTokenClaims gets the claims from a token without validation
func (j *JWTService) GetTokenClaims(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTService) GenerateTokenPair(userID, username, role string) (accessToken, refreshToken string, accessExpires, refreshExpires time.Time, err error) {
	// Generate access token
	accessToken, accessExpires, err = j.GenerateToken(userID, username, role)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	// Generate refresh token
	refreshToken, refreshExpires, err = j.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	j.logger.Info("Token pair generated", logging.Fields{
		"user_id":         userID,
		"username":        username,
		"access_expires":  accessExpires,
		"refresh_expires": refreshExpires,
	})

	return accessToken, refreshToken, accessExpires, refreshExpires, nil
}

// ValidateTokenPair validates both access and refresh tokens
func (j *JWTService) ValidateTokenPair(accessToken, refreshToken string) (*Claims, error) {
	// Validate access token
	claims, err := j.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	// In a real implementation, you would also validate the refresh token
	// For now, we'll just return the access token claims
	return claims, nil
}

// GetConfig returns the JWT configuration
func (j *JWTService) GetConfig() *JWTConfig {
	return j.config
}

// UpdateConfig updates the JWT configuration
func (j *JWTService) UpdateConfig(config *JWTConfig) {
	j.config = config
	j.logger.Info("JWT configuration updated")
}
