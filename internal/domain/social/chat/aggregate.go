package chat

import (
	"context"
	"time"
)

// ChatChannel 聊天频道聚合根
type ChatChannel struct {
	ID          string
	Name        string
	Type        ChannelType
	Description string
	MaxMembers  int
	members     map[string]*Member
	messages    []*Message
	createdAt   time.Time
	updatedAt   time.Time
	version     int64
}

// ChannelType 频道类型
type ChannelType int

const (
	ChannelTypeWorld  ChannelType = iota // 世界频道
	ChannelTypeGuild                     // 公会频道
	ChannelTypeTeam                      // 队伍频道
	ChannelTypePrivate                   // 私聊频道
	ChannelTypeSystem                    // 系统频道
)

// NewChatChannel 创建新的聊天频道
func NewChatChannel(id, name string, channelType ChannelType) *ChatChannel {
	return &ChatChannel{
		ID:        id,
		Name:      name,
		Type:      channelType,
		members:   make(map[string]*Member),
		messages:  make([]*Message, 0),
		createdAt: time.Now(),
		updatedAt: time.Now(),
		version:   1,
	}
}

// AddMember 添加成员
func (c *ChatChannel) AddMember(member *Member) error {
	if len(c.members) >= c.MaxMembers && c.MaxMembers > 0 {
		return ErrChannelFull
	}
	
	if _, exists := c.members[member.PlayerID]; exists {
		return ErrMemberAlreadyExists
	}
	
	c.members[member.PlayerID] = member
	c.updatedAt = time.Now()
	c.version++
	
	return nil
}

// RemoveMember 移除成员
func (c *ChatChannel) RemoveMember(playerID string) error {
	if _, exists := c.members[playerID]; !exists {
		return ErrMemberNotFound
	}
	
	delete(c.members, playerID)
	c.updatedAt = time.Now()
	c.version++
	
	return nil
}

// SendMessage 发送消息
func (c *ChatChannel) SendMessage(ctx context.Context, message *Message) error {
	// 验证发送者是否在频道中
	if _, exists := c.members[message.SenderID]; !exists && c.Type != ChannelTypeSystem {
		return ErrSenderNotInChannel
	}
	
	// 验证消息内容
	if err := message.Validate(); err != nil {
		return err
	}
	
	// 添加消息到频道
	c.messages = append(c.messages, message)
	c.updatedAt = time.Now()
	c.version++
	
	// 限制消息历史数量
	if len(c.messages) > MaxMessagesPerChannel {
		c.messages = c.messages[len(c.messages)-MaxMessagesPerChannel:]
	}
	
	return nil
}

// GetMembers 获取所有成员
func (c *ChatChannel) GetMembers() []*Member {
	members := make([]*Member, 0, len(c.members))
	for _, member := range c.members {
		members = append(members, member)
	}
	return members
}

// GetRecentMessages 获取最近的消息
func (c *ChatChannel) GetRecentMessages(limit int) []*Message {
	if limit <= 0 || limit > len(c.messages) {
		limit = len(c.messages)
	}
	
	start := len(c.messages) - limit
	if start < 0 {
		start = 0
	}
	
	return c.messages[start:]
}

// IsMember 检查是否为成员
func (c *ChatChannel) IsMember(playerID string) bool {
	_, exists := c.members[playerID]
	return exists
}

// GetVersion 获取版本号
func (c *ChatChannel) GetVersion() int64 {
	return c.version
}

const (
	MaxMessagesPerChannel = 100 // 每个频道最大消息数
)