package friend

import (
	"context"
	"fmt"
)

// FriendService 好友领域服务
type FriendService struct {
	friendRepo FriendRepository
}

// NewFriendService 创建好友服务
func NewFriendService(friendRepo FriendRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
	}
}

// SendFriendRequest 发送好友请求
func (s *FriendService) SendFriendRequest(ctx context.Context, fromPlayerID, toPlayerID, message string) (*FriendRequest, error) {
	// 验证不能给自己发送好友请求
	if fromPlayerID == toPlayerID {
		return nil, ErrCannotAddSelf
	}
	
	// 检查是否已经是好友
	friendship, err := s.friendRepo.GetFriendship(ctx, fromPlayerID, toPlayerID)
	if err != nil {
		return nil, fmt.Errorf("检查好友关系失败: %w", err)
	}
	if friendship != nil && friendship.IsActive() {
		return nil, ErrAlreadyFriends
	}
	
	// 检查是否已有待处理的请求
	existingRequest, err := s.friendRepo.GetPendingRequest(ctx, fromPlayerID, toPlayerID)
	if err != nil {
		return nil, fmt.Errorf("检查待处理请求失败: %w", err)
	}
	if existingRequest != nil {
		return nil, ErrRequestAlreadyExists
	}
	
	// 检查好友数量限制
	friendCount, err := s.friendRepo.GetFriendCount(ctx, fromPlayerID)
	if err != nil {
		return nil, fmt.Errorf("获取好友数量失败: %w", err)
	}
	if friendCount >= MaxFriendsPerPlayer {
		return nil, ErrTooManyFriends
	}
	
	// 创建好友请求
	request := NewFriendRequest(fromPlayerID, toPlayerID, message)
	
	// 保存请求
	if err := s.friendRepo.SaveFriendRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("保存好友请求失败: %w", err)
	}
	
	return request, nil
}

// AcceptFriendRequest 接受好友请求
func (s *FriendService) AcceptFriendRequest(ctx context.Context, requestID, playerID string) (*Friendship, error) {
	// 获取请求
	request, err := s.friendRepo.GetFriendRequestByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("获取好友请求失败: %w", err)
	}
	if request == nil {
		return nil, ErrRequestNotFound
	}
	
	// 验证权限（只有接收者可以接受请求）
	if request.ToPlayerID != playerID {
		return nil, ErrInsufficientPermission
	}
	
	// 接受请求
	if err := request.Accept(); err != nil {
		return nil, err
	}
	
	// 创建好友关系
	friendship := NewFriendship(request.FromPlayerID, request.ToPlayerID, request.FromPlayerID)
	if err := friendship.Accept(); err != nil {
		return nil, fmt.Errorf("创建好友关系失败: %w", err)
	}
	
	// 保存好友关系
	if err := s.friendRepo.SaveFriendship(ctx, friendship); err != nil {
		return nil, fmt.Errorf("保存好友关系失败: %w", err)
	}
	
	// 更新请求状态
	if err := s.friendRepo.SaveFriendRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("更新请求状态失败: %w", err)
	}
	
	return friendship, nil
}

// RejectFriendRequest 拒绝好友请求
func (s *FriendService) RejectFriendRequest(ctx context.Context, requestID, playerID string) error {
	// 获取请求
	request, err := s.friendRepo.GetFriendRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("获取好友请求失败: %w", err)
	}
	if request == nil {
		return ErrRequestNotFound
	}
	
	// 验证权限
	if request.ToPlayerID != playerID {
		return ErrInsufficientPermission
	}
	
	// 拒绝请求
	if err := request.Reject(); err != nil {
		return err
	}
	
	// 保存请求
	if err := s.friendRepo.SaveFriendRequest(ctx, request); err != nil {
		return fmt.Errorf("保存请求失败: %w", err)
	}
	
	return nil
}

// RemoveFriend 删除好友
func (s *FriendService) RemoveFriend(ctx context.Context, playerID, friendID string) error {
	// 获取好友关系
	friendship, err := s.friendRepo.GetFriendship(ctx, playerID, friendID)
	if err != nil {
		return fmt.Errorf("获取好友关系失败: %w", err)
	}
	if friendship == nil {
		return ErrFriendshipNotFound
	}
	
	// 删除好友关系
	if err := friendship.Delete(); err != nil {
		return err
	}
	
	// 保存更改
	if err := s.friendRepo.SaveFriendship(ctx, friendship); err != nil {
		return fmt.Errorf("保存好友关系失败: %w", err)
	}
	
	return nil
}

// BlockFriend 屏蔽好友
func (s *FriendService) BlockFriend(ctx context.Context, playerID, friendID string) error {
	// 获取好友关系
	friendship, err := s.friendRepo.GetFriendship(ctx, playerID, friendID)
	if err != nil {
		return fmt.Errorf("获取好友关系失败: %w", err)
	}
	if friendship == nil {
		return ErrFriendshipNotFound
	}
	
	// 屏蔽好友
	if err := friendship.Block(); err != nil {
		return err
	}
	
	// 保存更改
	if err := s.friendRepo.SaveFriendship(ctx, friendship); err != nil {
		return fmt.Errorf("保存好友关系失败: %w", err)
	}
	
	return nil
}

// GetFriendList 获取好友列表
func (s *FriendService) GetFriendList(ctx context.Context, playerID string) ([]*Friendship, error) {
	return s.friendRepo.GetFriendsByPlayerID(ctx, playerID)
}

// GetPendingRequests 获取待处理的好友请求
func (s *FriendService) GetPendingRequests(ctx context.Context, playerID string) ([]*FriendRequest, error) {
	return s.friendRepo.GetPendingRequestsByPlayerID(ctx, playerID)
}

const (
	MaxFriendsPerPlayer = 100 // 每个玩家最大好友数量
)