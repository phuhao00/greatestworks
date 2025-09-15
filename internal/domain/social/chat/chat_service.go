package chat

import (
	"context"
	"fmt"
)

// ChatService 聊天领域服务
type ChatService struct {
	chatRepo ChatRepository
}

// NewChatService 创建聊天服务
func NewChatService(chatRepo ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

// CreateChannel 创建聊天频道
func (s *ChatService) CreateChannel(ctx context.Context, name string, channelType ChannelType, creatorID string) (*ChatChannel, error) {
	// 验证频道名称
	if err := s.validateChannelName(name); err != nil {
		return nil, err
	}
	
	// 检查频道是否已存在
	exists, err := s.chatRepo.ChannelExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("检查频道是否存在失败: %w", err)
	}
	if exists {
		return nil, ErrChannelAlreadyExists
	}
	
	// 创建频道
	channelID := generateChannelID()
	channel := NewChatChannel(channelID, name, channelType)
	
	// 添加创建者为所有者
	creator := NewMember(creatorID, "")
	creator.SetRole(MemberRoleOwner)
	if err := channel.AddMember(creator); err != nil {
		return nil, fmt.Errorf("添加创建者失败: %w", err)
	}
	
	// 保存频道
	if err := s.chatRepo.SaveChannel(ctx, channel); err != nil {
		return nil, fmt.Errorf("保存频道失败: %w", err)
	}
	
	return channel, nil
}

// JoinChannel 加入频道
func (s *ChatService) JoinChannel(ctx context.Context, channelID, playerID, nickname string) error {
	// 获取频道
	channel, err := s.chatRepo.GetChannelByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("获取频道失败: %w", err)
	}
	if channel == nil {
		return ErrChannelNotFound
	}
	
	// 检查是否已经是成员
	if channel.IsMember(playerID) {
		return ErrMemberAlreadyExists
	}
	
	// 创建新成员
	member := NewMember(playerID, nickname)
	
	// 添加成员到频道
	if err := channel.AddMember(member); err != nil {
		return err
	}
	
	// 保存频道
	if err := s.chatRepo.SaveChannel(ctx, channel); err != nil {
		return fmt.Errorf("保存频道失败: %w", err)
	}
	
	return nil
}

// LeaveChannel 离开频道
func (s *ChatService) LeaveChannel(ctx context.Context, channelID, playerID string) error {
	// 获取频道
	channel, err := s.chatRepo.GetChannelByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("获取频道失败: %w", err)
	}
	if channel == nil {
		return ErrChannelNotFound
	}
	
	// 移除成员
	if err := channel.RemoveMember(playerID); err != nil {
		return err
	}
	
	// 保存频道
	if err := s.chatRepo.SaveChannel(ctx, channel); err != nil {
		return fmt.Errorf("保存频道失败: %w", err)
	}
	
	return nil
}

// SendMessage 发送消息
func (s *ChatService) SendMessage(ctx context.Context, channelID, senderID, content string, msgType MessageType) (*Message, error) {
	// 获取频道
	channel, err := s.chatRepo.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("获取频道失败: %w", err)
	}
	if channel == nil {
		return nil, ErrChannelNotFound
	}
	
	// 创建消息
	message := NewMessage(channelID, senderID, content, msgType)
	
	// 发送消息到频道
	if err := channel.SendMessage(ctx, message); err != nil {
		return nil, err
	}
	
	// 保存消息
	if err := s.chatRepo.SaveMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("保存消息失败: %w", err)
	}
	
	// 保存频道（更新版本）
	if err := s.chatRepo.SaveChannel(ctx, channel); err != nil {
		return nil, fmt.Errorf("保存频道失败: %w", err)
	}
	
	return message, nil
}

// GetChannelMessages 获取频道消息
func (s *ChatService) GetChannelMessages(ctx context.Context, channelID string, limit int) ([]*Message, error) {
	return s.chatRepo.GetMessagesByChannelID(ctx, channelID, limit)
}

// validateChannelName 验证频道名称
func (s *ChatService) validateChannelName(name string) error {
	if len(name) < 2 {
		return ErrChannelNameTooShort
	}
	if len(name) > 50 {
		return ErrChannelNameTooLong
	}
	return nil
}

// generateChannelID 生成频道ID
func generateChannelID() string {
	return "ch_" + randomString(16)
}