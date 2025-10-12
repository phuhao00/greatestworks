package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	playerCommands "greatestworks/internal/application/commands/player"
	"greatestworks/internal/application/handlers"
	"greatestworks/internal/infrastructure/logging"
)

// RegisterHandler 注册处理器
type RegisterHandler struct {
	commandBus *handlers.CommandBus
	logger     logging.Logger
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50"`
	Password   string `json:"password" binding:"required,min=6,max=100"`
	Email      string `json:"email" binding:"required,email"`
	PlayerName string `json:"player_name" binding:"required,min=2,max=50"`
	Avatar     string `json:"avatar,omitempty"`
	Gender     int    `json:"gender,omitempty"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewRegisterHandler 创建注册处理器
func NewRegisterHandler(commandBus *handlers.CommandBus, logger logging.Logger) *RegisterHandler {
	return &RegisterHandler{
		commandBus: commandBus,
		logger:     logger,
	}
}

// Register 处理用户注册
func (h *RegisterHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid register request", logging.Fields{
			"error": err,
		})
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// 验证用户名是否已存在
	if h.isUsernameExists(c.Request.Context(), req.Username) {
		h.logger.Warn("Username already exists", logging.Fields{
			"username": req.Username,
		})
		c.JSON(409, gin.H{"error": "Username already exists"})
		return
	}

	// 验证邮箱是否已存在
	if h.isEmailExists(c.Request.Context(), req.Email) {
		h.logger.Warn("Email already exists", logging.Fields{
			"email": req.Email,
		})
		c.JSON(409, gin.H{"error": "Email already exists"})
		return
	}

	// 验证玩家名称是否已存在
	if h.isPlayerNameExists(c.Request.Context(), req.PlayerName) {
		h.logger.Warn("Player name already exists", logging.Fields{
			"player_name": req.PlayerName,
		})
		c.JSON(409, gin.H{"error": "Player name already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := h.hashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", err)
		c.JSON(500, gin.H{"error": "Failed to process password"})
		return
	}

	// 创建用户账户
	userID, err := h.createUserAccount(c.Request.Context(), req.Username, hashedPassword, req.Email)
	if err != nil {
		h.logger.Error("Failed to create user account", err, logging.Fields{
			"username": req.Username,
		})
		c.JSON(500, gin.H{"error": "Failed to create user account"})
		return
	}

	// 创建玩家角色
	playerID, err := h.createPlayerCharacter(c.Request.Context(), userID, req.PlayerName, req.Avatar, req.Gender)
	if err != nil {
		h.logger.Error("Failed to create player character", err, logging.Fields{
			"user_id": userID,
		})
		c.JSON(500, gin.H{"error": "Failed to create player character"})
		return
	}

	// 返回注册成功响应
	response := RegisterResponse{
		UserID:     userID,
		Username:   req.Username,
		Email:      req.Email,
		PlayerID:   playerID,
		PlayerName: req.PlayerName,
		CreatedAt:  time.Now(),
	}

	h.logger.Info("User registered successfully", logging.Fields{
		"user_id":   userID,
		"username":  req.Username,
		"player_id": playerID,
	})
	c.JSON(201, response)
}

// 私有方法

// isUsernameExists 检查用户名是否已存在
func (h *RegisterHandler) isUsernameExists(ctx context.Context, username string) bool {
	// 这里应该查询数据库检查用户名是否存在
	// 简化实现，实际项目中应该调用相应的服务
	return false
}

// isEmailExists 检查邮箱是否已存在
func (h *RegisterHandler) isEmailExists(ctx context.Context, email string) bool {
	// 这里应该查询数据库检查邮箱是否存在
	// 简化实现，实际项目中应该调用相应的服务
	return false
}

// isPlayerNameExists 检查玩家名称是否已存在
func (h *RegisterHandler) isPlayerNameExists(ctx context.Context, playerName string) bool {
	// 这里应该查询数据库检查玩家名称是否存在
	// 简化实现，实际项目中应该调用相应的服务
	return false
}

// hashPassword 加密密码
func (h *RegisterHandler) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// createUserAccount 创建用户账户
func (h *RegisterHandler) createUserAccount(ctx context.Context, username, hashedPassword, email string) (string, error) {
	// 这里应该调用用户服务创建用户账户
	// 简化实现，返回模拟的用户ID
	return "user_" + username, nil
}

// createPlayerCharacter 创建玩家角色
func (h *RegisterHandler) createPlayerCharacter(ctx context.Context, userID, playerName, avatar string, gender int) (string, error) {
	// 创建玩家命令
	cmd := &playerCommands.CreatePlayerCommand{
		Name:   playerName,
		Avatar: avatar,
		Gender: gender,
	}

	// 执行命令
	result, err := h.commandBus.Execute(ctx, cmd)
	if err != nil {
		return "", err
	}

	// 获取玩家ID
	createResult, ok := result.(*playerCommands.CreatePlayerResult)
	if !ok {
		return "", fmt.Errorf("unexpected result type")
	}

	return createResult.PlayerID, nil
}
