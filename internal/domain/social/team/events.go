package team

import "time"

// TeamEvent 队伍事件接口
type TeamEvent interface {
	GetEventType() string
	GetTimestamp() time.Time
	GetTeamID() string
}

// BaseTeamEvent 基础队伍事件
type BaseTeamEvent struct {
	EventType string
	Timestamp time.Time
	TeamID    string
}

func (e BaseTeamEvent) GetEventType() string {
	return e.EventType
}

func (e BaseTeamEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e BaseTeamEvent) GetTeamID() string {
	return e.TeamID
}

// TeamCreatedEvent 队伍创建事件
type TeamCreatedEvent struct {
	BaseTeamEvent
	TeamName string
	LeaderID string
}

// NewTeamCreatedEvent 创建队伍创建事件
func NewTeamCreatedEvent(teamID, teamName, leaderID string) *TeamCreatedEvent {
	return &TeamCreatedEvent{
		BaseTeamEvent: BaseTeamEvent{
			EventType: "team.created",
			Timestamp: time.Now(),
			TeamID:    teamID,
		},
		TeamName: teamName,
		LeaderID: leaderID,
	}
}

// MemberJoinedTeamEvent 成员加入队伍事件
type MemberJoinedTeamEvent struct {
	BaseTeamEvent
	PlayerID string
	Nickname string
}

// NewMemberJoinedTeamEvent 创建成员加入队伍事件
func NewMemberJoinedTeamEvent(teamID, playerID, nickname string) *MemberJoinedTeamEvent {
	return &MemberJoinedTeamEvent{
		BaseTeamEvent: BaseTeamEvent{
			EventType: "team.member.joined",
			Timestamp: time.Now(),
			TeamID:    teamID,
		},
		PlayerID: playerID,
		Nickname: nickname,
	}
}

// MemberLeftTeamEvent 成员离开队伍事件
type MemberLeftTeamEvent struct {
	BaseTeamEvent
	PlayerID string
}

// NewMemberLeftTeamEvent 创建成员离开队伍事件
func NewMemberLeftTeamEvent(teamID, playerID string) *MemberLeftTeamEvent {
	return &MemberLeftTeamEvent{
		BaseTeamEvent: BaseTeamEvent{
			EventType: "team.member.left",
			Timestamp: time.Now(),
			TeamID:    teamID,
		},
		PlayerID: playerID,
	}
}

// TeamLeaderChangedEvent 队长变更事件
type TeamLeaderChangedEvent struct {
	BaseTeamEvent
	OldLeaderID string
	NewLeaderID string
}

// NewTeamLeaderChangedEvent 创建队长变更事件
func NewTeamLeaderChangedEvent(teamID, oldLeaderID, newLeaderID string) *TeamLeaderChangedEvent {
	return &TeamLeaderChangedEvent{
		BaseTeamEvent: BaseTeamEvent{
			EventType: "team.leader.changed",
			Timestamp: time.Now(),
			TeamID:    teamID,
		},
		OldLeaderID: oldLeaderID,
		NewLeaderID: newLeaderID,
	}
}