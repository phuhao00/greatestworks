package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"greatestworks/internal/infrastructure/persistence"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo *persistence.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *persistence.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Register 注册新用户
func (s *UserService) Register(ctx context.Context, username, password string) (int64, error) {
	// 检查用户名是否存在
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return 0, errors.New("username already exists")
	}

	// 密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// 生成用户ID（实际应用中应使用分布式ID生成器）
	userID := time.Now().UnixNano()

	// 创建用户
	user := &persistence.DbUser{
		UserID:       userID,
		Username:     username,
		PasswordHash: string(passwordHash),
		Status:       0,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, password string) (int64, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return 0, errors.New("invalid username or password")
	}

	// 检查用户状态
	if user.Status != 0 {
		return 0, errors.New("user account is banned")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return 0, errors.New("invalid username or password")
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, user.UserID); err != nil {
		// 记录日志但不影响登录
	}

	return user.UserID, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*persistence.DbUser, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// BanUser 封禁用户
func (s *UserService) BanUser(ctx context.Context, userID int64) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Status = 1
	return s.userRepo.Update(ctx, user)
}

// UnbanUser 解封用户
func (s *UserService) UnbanUser(ctx context.Context, userID int64) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Status = 0
	return s.userRepo.Update(ctx, user)
}
