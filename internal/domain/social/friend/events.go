package friend

import "time"

// FriendEvent 好友事件接口
type FriendEvent interface {
	GetEventType() string
	GetTimestamp() time.Time
	GetPlayerID() string
}

// BaseFriendEvent 基础好友事件
type BaseFriendEvent struct {
	EventType string
	Timestamp time.Time
	PlayerID  string
}

func (e BaseFriendEvent) GetEventType() string {
	return e.EventType
}

func (e BaseFriendEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e BaseFriendEvent) GetPlayerID() string {
	return e.PlayerID
}

// FriendRequestSentEvent 好友请求发送事件
type FriendRequestSentEvent struct {
	BaseFriendEvent
	RequestID  string
	ToPlayerID string
	Message    string
}

// NewFriendRequestSentEvent 创建好友请求发送事件
func NewFriendRequestSentEvent(playerID, requestID, toPlayerID, message string) *FriendRequestSentEvent {
	return &FriendRequestSentEvent{
		BaseFriendEvent: BaseFriendEvent{
			EventType: "friend.request.sent",
			Timestamp: time.Now(),
			PlayerID:  playerID,
		},
		RequestID:  requestID,
		ToPlayerID: toPlayerID,
		Message:    message,
	}
}

// FriendRequestAcceptedEvent 好友请求接受事件
type FriendRequestAcceptedEvent struct {
	BaseFriendEvent
	RequestID    string
	FromPlayerID string
	FriendshipID string
}

// NewFriendRequestAcceptedEvent 创建好友请求接受事件
func NewFriendRequestAcceptedEvent(playerID, requestID, fromPlayerID, friendshipID string) *FriendRequestAcceptedEvent {
	return &FriendRequestAcceptedEvent{
		BaseFriendEvent: BaseFriendEvent{
			EventType: "friend.request.accepted",
			Timestamp: time.Now(),
			PlayerID:  playerID,
		},
		RequestID:    requestID,
		FromPlayerID: fromPlayerID,
		FriendshipID: friendshipID,
	}
}

// FriendRequestRejectedEvent 好友请求拒绝事件
type FriendRequestRejectedEvent struct {
	BaseFriendEvent
	RequestID    string
	FromPlayerID string
}

// NewFriendRequestRejectedEvent 创建好友请求拒绝事件
func NewFriendRequestRejectedEvent(playerID, requestID, fromPlayerID string) *FriendRequestRejectedEvent {
	return &FriendRequestRejectedEvent{
		BaseFriendEvent: BaseFriendEvent{
			EventType: "friend.request.rejected",
			Timestamp: time.Now(),
			PlayerID:  playerID,
		},
		RequestID:    requestID,
		FromPlayerID: fromPlayerID,
	}
}

// FriendAddedEvent 好友添加事件
type FriendAddedEvent struct {
	BaseFriendEvent
	FriendID     string
	FriendshipID string
}

// NewFriendAddedEvent 创建好友添加事件
func NewFriendAddedEvent(playerID, friendID, friendshipID string) *FriendAddedEvent {
	return &FriendAddedEvent{
		BaseFriendEvent: BaseFriendEvent{
			EventType: "friend.added",
			Timestamp: time.Now(),
			PlayerID:  playerID,
		},
		FriendID:     friendID,
		FriendshipID: friendshipID,
	}
}

// FriendRemovedEvent 好友删除事件
type FriendRemovedEvent struct {
	BaseFriendEvent
	FriendID     string
	FriendshipID string
}

// NewFriendRemovedEvent 创建好友删除事件
func NewFriendRemovedEvent(playerID, friendID, friendshipID string) *FriendRemovedEvent {
	return &FriendRemovedEvent{
		BaseFriendEvent: BaseFriendEvent{
			EventType: "friend.removed",
			Timestamp: time.Now(),
			PlayerID:  playerID,
		},
		FriendID:     friendID,
		FriendshipID: friendshipID,
	}
}

// FriendBlockedEvent 好友屏蔽事件
type FriendBlockedEvent struct {
	BaseFriendEvent
	BlockedPlayerID string
	FriendshipID    string
}

// NewFriendBlockedEvent 创建好友屏蔽事件
func NewFriendBlockedEvent(playerID, blockedPlayerID, friendshipID string) *FriendBlockedEvent {
	return &FriendBlockedEvent{
		BaseFriendEvent: BaseFriendEvent{
			EventType: "friend.blocked",
			Timestamp: time.Now(),
			PlayerID:  playerID,
		},
		BlockedPlayerID: blockedPlayerID,
		FriendshipID:    friendshipID,
	}
}

// FriendUnblockedEvent 好友解除屏蔽事件
type FriendUnblockedEvent struct {
	BaseFriendEvent
	UnblockedPlayerID string
	FriendshipID      string
}

// NewFriendUnblockedEvent 创建好友解除屏蔽事件
func NewFriendUnblockedEvent(playerID, unblockedPlayerID, friendshipID string) *FriendUnblockedEvent {
	return &FriendUnblockedEvent{
		BaseFriendEvent: BaseFriendEvent{
			EventType: "friend.unblocked",
			Timestamp: time.Now(),
			PlayerID:  playerID,
		},
		UnblockedPlayerID: unblockedPlayerID,
		FriendshipID:      friendshipID,
	}
}
