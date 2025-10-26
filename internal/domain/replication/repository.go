package replication

import "context"

// Repository 副本实例仓储接口
type Repository interface {
	// Save 保存实例
	Save(ctx context.Context, instance *Instance) error

	// FindByID 根据ID查找实例
	FindByID(ctx context.Context, instanceID string) (*Instance, error)

	// FindByTemplateID 根据模板ID查找实例列表
	FindByTemplateID(ctx context.Context, templateID string) ([]*Instance, error)

	// FindActiveInstances 查找所有活跃实例
	FindActiveInstances(ctx context.Context) ([]*Instance, error)

	// FindByPlayerID 根据玩家ID查找实例
	FindByPlayerID(ctx context.Context, playerID string) (*Instance, error)

	// Delete 删除实例
	Delete(ctx context.Context, instanceID string) error

	// UpdateStatus 更新实例状态
	UpdateStatus(ctx context.Context, instanceID string, status InstanceStatus) error

	// FindExpiredInstances 查找过期实例
	FindExpiredInstances(ctx context.Context) ([]*Instance, error)
}
