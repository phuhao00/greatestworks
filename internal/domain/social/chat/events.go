package chat

import "time"

// ChatEvent 聊天事件接口
type ChatEvent interface {
	GetEventType() string
	GetTimestamp() time.Time
	GetChannelID() string
}

// BaseEvent 基础事件
type BaseEvent struct {
	EventType string
	Timestamp time.Time
	ChannelID string
}

func (e BaseEvent) GetEventType() string {
	return e.EventType
}

func (e BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e BaseEvent) GetChannelID() string {
	return e.ChannelID
}

// ChannelCreatedEvent 频道创建事件
type ChannelCreatedEvent struct {
	BaseEvent
	ChannelName string
	ChannelType ChannelType
	CreatorID   string
}

// NewChannelCreatedEvent 创建频道创建事件
func NewChannelCreatedEvent(channelID, channelName string, channelType ChannelType, creatorID string) *ChannelCreatedEvent {
	return &ChannelCreatedEvent{
		BaseEvent: BaseEvent{
			EventType: "channel.created",
			Timestamp: time.Now(),
			ChannelID: channelID,
		},
		ChannelName: channelName,
		ChannelType: channelType,
		CreatorID:   creatorID,
	}
}

// MemberJoinedEvent 成员加入事件
type MemberJoinedEvent struct {
	BaseEvent
	PlayerID string
	Nickname string
}

// NewMemberJoinedEvent 创建成员加入事件
func NewMemberJoinedEvent(channelID, playerID, nickname string) *MemberJoinedEvent {
	return &MemberJoinedEvent{
		BaseEvent: BaseEvent{
			EventType: "member.joined",
			Timestamp: time.Now(),
			ChannelID: channelID,
		},
		PlayerID: playerID,
		Nickname: nickname,
	}
}

// MemberLeftEvent 成员离开事件
type MemberLeftEvent struct {
	BaseEvent
	PlayerID string
}

// NewMemberLeftEvent 创建成员离开事件
func NewMemberLeftEvent(channelID, playerID string) *MemberLeftEvent {
	return &MemberLeftEvent{
		BaseEvent: BaseEvent{
			EventType: "member.left",
			Timestamp: time.Now(),
			ChannelID: channelID,
		},
		PlayerID: playerID,
	}
}

// MessageSentEvent 消息发送事件
type MessageSentEvent struct {
	BaseEvent
	MessageID string
	SenderID  string
	Content   string
	Type      MessageType
}

// NewMessageSentEvent 创建消息发送事件
func NewMessageSentEvent(channelID, messageID, senderID, content string, msgType MessageType) *MessageSentEvent {
	return &MessageSentEvent{
		BaseEvent: BaseEvent{
			EventType: "message.sent",
			Timestamp: time.Now(),
			ChannelID: channelID,
		},
		MessageID: messageID,
		SenderID:  senderID,
		Content:   content,
		Type:      msgType,
	}
}

// MemberMutedEvent 成员被禁言事件
type MemberMutedEvent struct {
	BaseEvent
	PlayerID   string
	MutedBy    string
	Duration   time.Duration
	Reason     string
}

// NewMemberMutedEvent 创建成员禁言事件
func NewMemberMutedEvent(channelID, playerID, mutedBy string, duration time.Duration, reason string) *MemberMutedEvent {
	return &MemberMutedEvent{
		BaseEvent: BaseEvent{
			EventType: "member.muted",
			Timestamp: time.Now(),
			ChannelID: channelID,
		},
		PlayerID: playerID,
		MutedBy:  mutedBy,
		Duration: duration,
		Reason:   reason,
	}
}

// MemberUnmutedEvent 成员解除禁言事件
type MemberUnmutedEvent struct {
	BaseEvent
	PlayerID  string
	UnmutedBy string
}

// NewMemberUnmutedEvent 创建成员解除禁言事件
func NewMemberUnmutedEvent(channelID, playerID, unmutedBy string) *MemberUnmutedEvent {
	return &MemberUnmutedEvent{
		BaseEvent: BaseEvent{
			EventType: "member.unmuted",
			Timestamp: time.Now(),
			ChannelID: channelID,
		},
		PlayerID:  playerID,
		UnmutedBy: unmutedBy,
	}
}