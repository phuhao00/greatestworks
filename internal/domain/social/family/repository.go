package family

import "context"

// FamilyRepository 家族仓储接口
type FamilyRepository interface {
	// 家族相关
	SaveFamily(ctx context.Context, family *Family) error
	GetFamilyByID(ctx context.Context, familyID string) (*Family, error)
	GetFamilyByName(ctx context.Context, name string) (*Family, error)
	FamilyExistsByName(ctx context.Context, name string) (bool, error)
	DeleteFamily(ctx context.Context, familyID string) error
	GetFamiliesByLevel(ctx context.Context, level int) ([]*Family, error)

	// 成员相关
	GetFamilyByPlayerID(ctx context.Context, playerID string) (*Family, error)
	PlayerHasFamily(ctx context.Context, playerID string) (bool, error)
	GetFamilyMembers(ctx context.Context, familyID string) ([]*FamilyMember, error)
	UpdateMember(ctx context.Context, familyID string, member *FamilyMember) error

	// 查询相关
	GetTopFamiliesByLevel(ctx context.Context, limit int) ([]*Family, error)
	GetTopFamiliesByMemberCount(ctx context.Context, limit int) ([]*Family, error)
	SearchFamiliesByName(ctx context.Context, keyword string, limit int) ([]*Family, error)
}
