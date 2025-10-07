package player

import (
	"context"
	"time"
)

// GetPlayerQuery 获取玩家查询
type GetPlayerQuery struct {
	PlayerID string `json:"player_id" validate:"required"`
}

// PlayerDTO 玩家数据传输对象
type PlayerDTO struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Level        int          `json:"level"`
	Exp          int64        `json:"exp"`
	Status       string       `json:"status"`
	Position     PositionDTO  `json:"position"`
	Stats        StatsDTO     `json:"stats"`
	Username     string       `json:"username,omitempty"`
	Email        string       `json:"email,omitempty"`
	Avatar       string       `json:"avatar,omitempty"`
	Gender       int          `json:"gender,omitempty"`
	BanInfo      *BanInfoDTO  `json:"ban_info,omitempty"`
	LastLoginAt  *time.Time   `json:"last_login_at,omitempty"`
	LastLogoutAt *time.Time   `json:"last_logout_at,omitempty"`
	OnlineTime   int64        `json:"online_time,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// PositionDTO 位置数据传输对象
type PositionDTO struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// StatsDTO 属性数据传输对象
type StatsDTO struct {
	HP      int `json:"hp"`
	MaxHP   int `json:"max_hp"`
	MP      int `json:"mp"`
	MaxMP   int `json:"max_mp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
}

// BanInfoDTO 封禁信息数据传输对象
type BanInfoDTO struct {
	IsBanned  bool       `json:"is_banned"`
	BanType   string     `json:"ban_type,omitempty"`
	Reason    string     `json:"reason,omitempty"`
	BannedBy  string     `json:"banned_by,omitempty"`
	BannedAt  *time.Time `json:"banned_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// GetPlayerResult 获取玩家结果
type GetPlayerResult struct {
	Player *PlayerDTO `json:"player"`
	Found  bool       `json:"found"`
}

// GetPlayerByUsernameQuery 根据用户名获取玩家查询
type GetPlayerByUsernameQuery struct {
	Username string `json:"username" validate:"required"`
}

// QueryType 返回查询类型
func (query *GetPlayerByUsernameQuery) QueryType() string {
	return "GetPlayerByUsername"
}

// Validate 验证查询
func (query *GetPlayerByUsernameQuery) Validate() error {
	if query.Username == "" {
		return ErrInvalidUsername
	}
	return nil
}

// GetPlayerByUsernameResult 根据用户名获取玩家结果
type GetPlayerByUsernameResult struct {
	Player *PlayerWithAccountDTO `json:"player"`
	Found  bool                  `json:"found"`
}

// PlayerWithAccountDTO 带账户信息的玩家数据传输对象
type PlayerWithAccountDTO struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Level        int         `json:"level"`
	Exp          int64       `json:"exp"`
	Status       string      `json:"status"`
	Position     PositionDTO `json:"position"`
	Stats        StatsDTO    `json:"stats"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"password_hash"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// GetPlayerHandler 获取玩家查询处理器
type GetPlayerHandler struct {
	playerQueryService PlayerQueryService
}

// PlayerQueryService 玩家查询服务接口
type PlayerQueryService interface {
	GetPlayer(ctx context.Context, playerID string) (*PlayerDTO, error)
	GetPlayerByName(ctx context.Context, name string) (*PlayerDTO, error)
	GetOnlinePlayers(ctx context.Context, limit int) ([]*PlayerDTO, error)
	GetPlayersByLevel(ctx context.Context, minLevel, maxLevel int) ([]*PlayerDTO, error)
	ListPlayers(ctx context.Context, query *ListPlayersQuery) ([]*PlayerDTO, int64, error)
}

// NewGetPlayerHandler 创建查询处理器
func NewGetPlayerHandler(playerQueryService PlayerQueryService) *GetPlayerHandler {
	return &GetPlayerHandler{
		playerQueryService: playerQueryService,
	}
}

// Handle 处理获取玩家查询
func (h *GetPlayerHandler) Handle(ctx context.Context, query *GetPlayerQuery) (*GetPlayerResult, error) {
	player, err := h.playerQueryService.GetPlayer(ctx, query.PlayerID)
	if err != nil {
		if err == ErrPlayerNotFound {
			return &GetPlayerResult{Player: nil, Found: false}, nil
		}
		return nil, err
	}

	return &GetPlayerResult{Player: player, Found: true}, nil
}

// QueryType 返回查询类型
func (query *GetPlayerQuery) QueryType() string {
	return "GetPlayer"
}

// Validate 验证查询
func (query *GetPlayerQuery) Validate() error {
	if query.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	return nil
}
