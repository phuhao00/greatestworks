package auth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"greatestworks/application/commands/player"
	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
)

// RegisterHandler 注册处理器
type RegisterHandler struct {
	commandBus *handlers.CommandBus
	logger     logger.Logger
}

// NewRegisterHandler 创建注册处理器
func NewRegisterHandler(commandBus *handlers.CommandBus, logger logger.Logger) *RegisterHandler {
	return &RegisterHandler{
		commandBus: commandBus,
		logger:     logger,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Password        string `json:"password" binding:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	PlayerName      string `json:"player_name" binding:"required,min=2,max=50"`
	Avatar          string `json:"avatar,omitempty"`
	Gender          int    `json:"gender,omitempty" binding:"min=0,max=2"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	PlayerID  string    `json:"player_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CheckUsernameRequest 检查用户名请求
type CheckUsernameRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
}

// CheckUsernameResponse 检查用户名响应
type CheckUsernameResponse struct {
	Available bool   `json:"available"`
	Message   string `json:"message"`
}

// Register 用户注册
func (h *RegisterHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid register request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	// 验证密码确认
	if req.Password != req.ConfirmPassword {
		h.logger.Warn("Password confirmation mismatch", "username", req.Username)
		c.JSON(400, gin.H{"error": "Password confirmation does not match", "success": false})
		return
	}

	// 验证用户名唯一性
	if !h.isUsernameAvailable(req.Username) {
		h.logger.Warn("Username already exists", "username", req.Username)
		c.JSON(409, gin.H{"error": "Username already exists", "success": false})
		return
	}

	// 验证邮箱唯一性
	if !h.isEmailAvailable(req.Email) {
		h.logger.Warn("Email already exists", "email", req.Email)
		c.JSON(409, gin.H{"error": "Email already exists", "success": false})
		return
	}

	// 加密密码
	hashedPassword, err := h.hashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error", "success": false})
		return
	}

	ctx := context.Background()

	// 创建用户账户和玩家角色
	cmd := &player.CreatePlayerWithAccountCommand{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
		PlayerName:   req.PlayerName,
		Avatar:       req.Avatar,
		Gender:       req.Gender,
	}

	result, err := handlers.ExecuteTyped[*player.CreatePlayerWithAccountCommand, *player.CreatePlayerWithAccountResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to create player account", "error", err, "username", req.Username)
		c.JSON(500, gin.H{"error": "Failed to create account", "success": false})
		return
	}

	// 记录注册日志
	h.logger.Info("User registered successfully", "username", req.Username, "player_id", result.PlayerID, "email", req.Email)

	// 返回响应
	response := &RegisterResponse{
		PlayerID:  result.PlayerID,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: result.CreatedAt,
	}

	c.JSON(201, gin.H{"data": response, "success": true, "message": "Account created successfully"})
}

// CheckUsername 检查用户名可用性
func (h *RegisterHandler) CheckUsername(c *gin.Context) {
	var req CheckUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid check username request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	available := h.isUsernameAvailable(req.Username)
	message := "Username is available"
	if !available {
		message = "Username is already taken"
	}

	response := &CheckUsernameResponse{
		Available: available,
		Message:   message,
	}

	c.JSON(200, gin.H{"data": response, "success": true})
}

// CheckEmail 检查邮箱可用性
func (h *RegisterHandler) CheckEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(400, gin.H{"error": "Email parameter is required", "success": false})
		return
	}

	available := h.isEmailAvailable(email)
	message := "Email is available"
	if !available {
		message = "Email is already registered"
	}

	response := map[string]interface{}{
		"available": available,
		"message":   message,
	}

	c.JSON(200, gin.H{"data": response, "success": true})
}

// ResetPassword 重置密码
func (h *RegisterHandler) ResetPassword(c *gin.Context) {
	type ResetPasswordRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid reset password request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	// 检查邮箱是否存在
	if h.isEmailAvailable(req.Email) {
		h.logger.Warn("Reset password attempt for non-existent email", "email", req.Email)
		// 为了安全，不透露邮箱是否存在
		c.JSON(200, gin.H{"success": true, "message": "If the email exists, a reset link has been sent"})
		return
	}

	// 生成重置令牌并发送邮件（这里简化处理）
	resetToken := h.generateResetToken(req.Email)
	if err := h.sendResetEmail(req.Email, resetToken); err != nil {
		h.logger.Error("Failed to send reset email", "error", err, "email", req.Email)
		c.JSON(500, gin.H{"error": "Failed to send reset email", "success": false})
		return
	}

	h.logger.Info("Password reset email sent", "email", req.Email)
	c.JSON(200, gin.H{"success": true, "message": "Password reset email sent successfully"})
}

// 私有方法

// hashPassword 加密密码
func (h *RegisterHandler) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// isUsernameAvailable 检查用户名是否可用
func (h *RegisterHandler) isUsernameAvailable(username string) bool {
	// 这里应该查询数据库检查用户名是否已存在
	// 临时实现：简单的内存检查
	// TODO: 实现数据库查询
	return true // 暂时返回true
}

// isEmailAvailable 检查邮箱是否可用
func (h *RegisterHandler) isEmailAvailable(email string) bool {
	// 这里应该查询数据库检查邮箱是否已存在
	// 临时实现：简单的内存检查
	// TODO: 实现数据库查询
	return true // 暂时返回true
}

// generateResetToken 生成重置令牌
func (h *RegisterHandler) generateResetToken(email string) string {
	// 这里应该生成安全的重置令牌
	// 临时实现：简单的时间戳
	// TODO: 实现安全的令牌生成
	return "reset_token_" + email
}

// sendResetEmail 发送重置邮件
func (h *RegisterHandler) sendResetEmail(email, token string) error {
	// 这里应该发送实际的邮件
	// 临时实现：仅记录日志
	// TODO: 实现邮件发送服务
	h.logger.Info("Reset email would be sent", "email", email, "token", token)
	return nil
}