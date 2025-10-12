package team

import "context"

// TeamRepository 队伍仓储接口
type TeamRepository interface {
	// 队伍相关
	SaveTeam(ctx context.Context, team *Team) error
	GetTeamByID(ctx context.Context, teamID string) (*Team, error)
	DeleteTeam(ctx context.Context, teamID string) error
	GetPublicTeams(ctx context.Context, limit int) ([]*Team, error)

	// 成员相关
	GetTeamByPlayerID(ctx context.Context, playerID string) (*Team, error)
	PlayerHasTeam(ctx context.Context, playerID string) (bool, error)
	GetTeamMembers(ctx context.Context, teamID string) ([]*TeamMember, error)
	UpdateMember(ctx context.Context, teamID string, member *TeamMember) error

	// 查询相关
	SearchTeamsByName(ctx context.Context, keyword string, limit int) ([]*Team, error)
	GetTeamsByLevelRange(ctx context.Context, minLevel, maxLevel int, limit int) ([]*Team, error)
}
