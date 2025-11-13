package replication

import "time"

// InstanceSnapshot 用于持久化/重建的快照
type InstanceSnapshot struct {
	InstanceID, TemplateID, SceneID, OwnerPlayerID string
	InstanceType                                   InstanceType
	Status                                         InstanceStatus
	MaxPlayers                                     int
	MinPlayers                                     int
	Difficulty                                     int
	CreatedAt                                      time.Time
	StartedAt                                      time.Time
	ExpireAt                                       time.Time
	ClosedAt                                       time.Time
	Lifetime                                       time.Duration
	Progress                                       int
	Completed                                      []string
	Metadata                                       map[string]string
	Players                                        []PlayerInfo
}

// NewInstanceFromSnapshot 通过快照重建实例
func NewInstanceFromSnapshot(s InstanceSnapshot) *Instance {
	inst := &Instance{
		instanceID:      s.InstanceID,
		templateID:      s.TemplateID,
		instanceType:    s.InstanceType,
		sceneID:         s.SceneID,
		players:         make(map[string]*PlayerInfo),
		maxPlayers:      s.MaxPlayers,
		minPlayers:      s.MinPlayers,
		ownerPlayerID:   s.OwnerPlayerID,
		status:          s.Status,
		difficulty:      s.Difficulty,
		createdAt:       s.CreatedAt,
		startedAt:       s.StartedAt,
		expireAt:        s.ExpireAt,
		closedAt:        s.ClosedAt,
		lifetime:        s.Lifetime,
		progress:        s.Progress,
		completedTasks:  append([]string(nil), s.Completed...),
		metadata:        map[string]string{},
		scoreMultiplier: 1.0,
		events:          nil,
	}
	if s.Metadata != nil {
		for k, v := range s.Metadata {
			inst.metadata[k] = v
		}
	}
	for _, p := range s.Players {
		cp := p // copy
		inst.players[p.PlayerID] = &cp
	}
	return inst
}
