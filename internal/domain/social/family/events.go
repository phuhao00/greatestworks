package family

import "time"

// FamilyEvent 家族事件接口
type FamilyEvent interface {
	GetEventType() string
	GetTimestamp() time.Time
	GetFamilyID() string
}

// BaseFamilyEvent 基础家族事件
type BaseFamilyEvent struct {
	EventType string
	Timestamp time.Time
	FamilyID  string
}

func (e BaseFamilyEvent) GetEventType() string {
	return e.EventType
}

func (e BaseFamilyEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e BaseFamilyEvent) GetFamilyID() string {
	return e.FamilyID
}

// FamilyCreatedEvent 家族创建事件
type FamilyCreatedEvent struct {
	BaseFamilyEvent
	FamilyName string
	LeaderID   string
}

// NewFamilyCreatedEvent 创建家族创建事件
func NewFamilyCreatedEvent(familyID, familyName, leaderID string) *FamilyCreatedEvent {
	return &FamilyCreatedEvent{
		BaseFamilyEvent: BaseFamilyEvent{
			EventType: "family.created",
			Timestamp: time.Now(),
			FamilyID:  familyID,
		},
		FamilyName: familyName,
		LeaderID:   leaderID,
	}
}

// MemberJoinedFamilyEvent 成员加入家族事件
type MemberJoinedFamilyEvent struct {
	BaseFamilyEvent
	PlayerID string
	Nickname string
}

// NewMemberJoinedFamilyEvent 创建成员加入家族事件
func NewMemberJoinedFamilyEvent(familyID, playerID, nickname string) *MemberJoinedFamilyEvent {
	return &MemberJoinedFamilyEvent{
		BaseFamilyEvent: BaseFamilyEvent{
			EventType: "family.member.joined",
			Timestamp: time.Now(),
			FamilyID:  familyID,
		},
		PlayerID: playerID,
		Nickname: nickname,
	}
}

// MemberLeftFamilyEvent 成员离开家族事件
type MemberLeftFamilyEvent struct {
	BaseFamilyEvent
	PlayerID string
}

// NewMemberLeftFamilyEvent 创建成员离开家族事件
func NewMemberLeftFamilyEvent(familyID, playerID string) *MemberLeftFamilyEvent {
	return &MemberLeftFamilyEvent{
		BaseFamilyEvent: BaseFamilyEvent{
			EventType: "family.member.left",
			Timestamp: time.Now(),
			FamilyID:  familyID,
		},
		PlayerID: playerID,
	}
}

// LeadershipTransferredEvent 族长转让事件
type LeadershipTransferredEvent struct {
	BaseFamilyEvent
	OldLeaderID string
	NewLeaderID string
}

// NewLeadershipTransferredEvent 创建族长转让事件
func NewLeadershipTransferredEvent(familyID, oldLeaderID, newLeaderID string) *LeadershipTransferredEvent {
	return &LeadershipTransferredEvent{
		BaseFamilyEvent: BaseFamilyEvent{
			EventType: "family.leadership.transferred",
			Timestamp: time.Now(),
			FamilyID:  familyID,
		},
		OldLeaderID: oldLeaderID,
		NewLeaderID: newLeaderID,
	}
}

// FamilyLevelUpEvent 家族升级事件
type FamilyLevelUpEvent struct {
	BaseFamilyEvent
	OldLevel int
	NewLevel int
}

// NewFamilyLevelUpEvent 创建家族升级事件
func NewFamilyLevelUpEvent(familyID string, oldLevel, newLevel int) *FamilyLevelUpEvent {
	return &FamilyLevelUpEvent{
		BaseFamilyEvent: BaseFamilyEvent{
			EventType: "family.level.up",
			Timestamp: time.Now(),
			FamilyID:  familyID,
		},
		OldLevel: oldLevel,
		NewLevel: newLevel,
	}
}
