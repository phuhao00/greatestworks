package player

import "time"

// PlayerQuery 玩家查询条件
type PlayerQuery struct {
	// 基础查询条件
	ID       *PlayerID     `json:"id,omitempty"`
	Username string        `json:"username,omitempty"`
	Nickname string        `json:"nickname,omitempty"`
	Status   *PlayerStatus `json:"status,omitempty"`

	// 等级范围
	MinLevel int `json:"min_level,omitempty"`
	MaxLevel int `json:"max_level,omitempty"`

	// 时间范围
	CreatedAfter    *time.Time `json:"created_after,omitempty"`
	CreatedBefore   *time.Time `json:"created_before,omitempty"`
	LastLoginAfter  *time.Time `json:"last_login_after,omitempty"`
	LastLoginBefore *time.Time `json:"last_login_before,omitempty"`

	// 分页
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// 排序
	OrderBy string `json:"order_by,omitempty"` // "created_at", "level", "vip_level"
	Order   string `json:"order,omitempty"`    // "asc", "desc"
}

// NewPlayerQuery 创建新的玩家查询
func NewPlayerQuery() *PlayerQuery {
	return &PlayerQuery{}
}

// WithID 设置ID查询条件
func (q *PlayerQuery) WithID(id PlayerID) *PlayerQuery {
	q.ID = &id
	return q
}

// WithUsername 设置用户名查询条件
func (q *PlayerQuery) WithUsername(username string) *PlayerQuery {
	q.Username = username
	return q
}

// WithNickname 设置昵称查询条件
func (q *PlayerQuery) WithNickname(nickname string) *PlayerQuery {
	q.Nickname = nickname
	return q
}

// WithStatus 设置状态查询条件
func (q *PlayerQuery) WithStatus(status PlayerStatus) *PlayerQuery {
	q.Status = &status
	return q
}

// WithLevelRange 设置等级范围查询条件
func (q *PlayerQuery) WithLevelRange(minLevel, maxLevel int) *PlayerQuery {
	q.MinLevel = minLevel
	q.MaxLevel = maxLevel
	return q
}

// WithCreatedTimeRange 设置创建时间范围查询条件
func (q *PlayerQuery) WithCreatedTimeRange(after, before time.Time) *PlayerQuery {
	q.CreatedAfter = &after
	q.CreatedBefore = &before
	return q
}

// WithLastLoginTimeRange 设置最后登录时间范围查询条件
func (q *PlayerQuery) WithLastLoginTimeRange(after, before time.Time) *PlayerQuery {
	q.LastLoginAfter = &after
	q.LastLoginBefore = &before
	return q
}

// WithPagination 设置分页查询条件
func (q *PlayerQuery) WithPagination(limit, offset int) *PlayerQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}

// WithOrder 设置排序查询条件
func (q *PlayerQuery) WithOrder(orderBy, order string) *PlayerQuery {
	q.OrderBy = orderBy
	q.Order = order
	return q
}
