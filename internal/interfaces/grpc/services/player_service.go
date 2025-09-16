package services

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"greatestworks/application/commands/player"
	"greatestworks/application/queries/player"
	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
	pb "greatestworks/internal/interfaces/grpc/proto"
)

// PlayerServiceImpl 玩家服务实现
type PlayerServiceImpl struct {
	pb.UnimplementedPlayerServiceServer
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	logger     logger.Logger
}

// NewPlayerService 创建玩家服务
func NewPlayerService(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logger.Logger) *PlayerServiceImpl {
	return &PlayerServiceImpl{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

// CreatePlayer 创建玩家
func (s *PlayerServiceImpl) CreatePlayer(ctx context.Context, req *pb.CreatePlayerRequest) (*pb.CreatePlayerResponse, error) {
	s.logger.Info("Creating player", "name", req.Name, "email", req.Email)

	// 验证请求参数
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "player name is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "player email is required")
	}

	// 创建命令
	cmd := &player.CreatePlayerCommand{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// 添加初始数据
	if req.InitialData != nil {
		cmd.InitialData = req.InitialData
	}

	// 执行命令
	result, err := handlers.ExecuteTyped[*player.CreatePlayerCommand, *player.CreatePlayerResult](ctx, s.commandBus, cmd)
	if err != nil {
		s.logger.Error("Failed to create player", "error", err, "name", req.Name)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create player: %v", err))
	}

	if !result.Success {
		s.logger.Warn("Player creation failed", "reason", result.Message, "name", req.Name)
		return &pb.CreatePlayerResponse{
			Base: &pb.BaseResponse{
				Success:   false,
				Message:   result.Message,
				Code:      int32(result.ErrorCode),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	// 构造响应
	response := &pb.CreatePlayerResponse{
		Base: &pb.BaseResponse{
			Success:   true,
			Message:   "Player created successfully",
			Code:      0,
			Timestamp: time.Now().Unix(),
		},
		Player: s.convertToPlayerInfo(result.Player),
	}

	s.logger.Info("Player created successfully", "player_id", result.Player.ID, "name", req.Name)
	return response, nil
}

// GetPlayer 获取玩家信息
func (s *PlayerServiceImpl) GetPlayer(ctx context.Context, req *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	s.logger.Debug("Getting player", "player_id", req.PlayerId)

	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	// 创建查询
	query := &player.GetPlayerQuery{
		PlayerID:        req.PlayerId,
		IncludeStats:    req.IncludeStats,
		IncludeMetadata: req.IncludeMetadata,
	}

	// 执行查询
	result, err := handlers.ExecuteQueryTyped[*player.GetPlayerQuery, *player.GetPlayerResult](ctx, s.queryBus, query)
	if err != nil {
		s.logger.Error("Failed to get player", "error", err, "player_id", req.PlayerId)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get player: %v", err))
	}

	response := &pb.GetPlayerResponse{
		Base: &pb.BaseResponse{
			Success:   true,
			Message:   "Player retrieved successfully",
			Code:      0,
			Timestamp: time.Now().Unix(),
		},
		Found: result.Found,
	}

	if result.Found {
		response.Player = s.convertToPlayerInfo(result.Player)
		s.logger.Debug("Player found", "player_id", req.PlayerId, "name", result.Player.Name)
	} else {
		s.logger.Debug("Player not found", "player_id", req.PlayerId)
	}

	return response, nil
}

// UpdatePlayer 更新玩家信息
func (s *PlayerServiceImpl) UpdatePlayer(ctx context.Context, req *pb.UpdatePlayerRequest) (*pb.UpdatePlayerResponse, error) {
	s.logger.Info("Updating player", "player_id", req.PlayerId)

	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	// 创建命令
	cmd := &player.UpdatePlayerCommand{
		PlayerID: req.PlayerId,
	}

	// 设置更新字段
	if req.Name != nil {
		cmd.Name = *req.Name
		cmd.UpdateFields = append(cmd.UpdateFields, "name")
	}
	if req.Email != nil {
		cmd.Email = *req.Email
		cmd.UpdateFields = append(cmd.UpdateFields, "email")
	}
	if req.Level != nil {
		cmd.Level = *req.Level
		cmd.UpdateFields = append(cmd.UpdateFields, "level")
	}
	if req.Exp != nil {
		cmd.Exp = *req.Exp
		cmd.UpdateFields = append(cmd.UpdateFields, "exp")
	}
	if req.Status != nil {
		cmd.Status = *req.Status
		cmd.UpdateFields = append(cmd.UpdateFields, "status")
	}
	if req.Position != nil {
		cmd.Position = player.Position{
			X: req.Position.X,
			Y: req.Position.Y,
			Z: req.Position.Z,
		}
		cmd.UpdateFields = append(cmd.UpdateFields, "position")
	}
	if req.Stats != nil {
		cmd.Stats = player.Stats{
			HP:      int(req.Stats.Hp),
			MaxHP:   int(req.Stats.MaxHp),
			MP:      int(req.Stats.Mp),
			MaxMP:   int(req.Stats.MaxMp),
			Attack:  int(req.Stats.Attack),
			Defense: int(req.Stats.Defense),
			Speed:   int(req.Stats.Speed),
		}
		cmd.UpdateFields = append(cmd.UpdateFields, "stats")
	}
	if req.Metadata != nil {
		cmd.Metadata = req.Metadata
		cmd.UpdateFields = append(cmd.UpdateFields, "metadata")
	}

	// 执行命令
	result, err := handlers.ExecuteTyped[*player.UpdatePlayerCommand, *player.UpdatePlayerResult](ctx, s.commandBus, cmd)
	if err != nil {
		s.logger.Error("Failed to update player", "error", err, "player_id", req.PlayerId)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update player: %v", err))
	}

	if !result.Success {
		s.logger.Warn("Player update failed", "reason", result.Message, "player_id", req.PlayerId)
		return &pb.UpdatePlayerResponse{
			Base: &pb.BaseResponse{
				Success:   false,
				Message:   result.Message,
				Code:      int32(result.ErrorCode),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	response := &pb.UpdatePlayerResponse{
		Base: &pb.BaseResponse{
			Success:   true,
			Message:   "Player updated successfully",
			Code:      0,
			Timestamp: time.Now().Unix(),
		},
		Player: s.convertToPlayerInfo(result.Player),
	}

	s.logger.Info("Player updated successfully", "player_id", req.PlayerId)
	return response, nil
}

// DeletePlayer 删除玩家
func (s *PlayerServiceImpl) DeletePlayer(ctx context.Context, req *pb.DeletePlayerRequest) (*pb.DeletePlayerResponse, error) {
	s.logger.Info("Deleting player", "player_id", req.PlayerId, "soft_delete", req.SoftDelete)

	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	// 创建命令
	cmd := &player.DeletePlayerCommand{
		PlayerID:   req.PlayerId,
		SoftDelete: req.SoftDelete,
		Reason:     req.Reason,
	}

	// 执行命令
	result, err := handlers.ExecuteTyped[*player.DeletePlayerCommand, *player.DeletePlayerResult](ctx, s.commandBus, cmd)
	if err != nil {
		s.logger.Error("Failed to delete player", "error", err, "player_id", req.PlayerId)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete player: %v", err))
	}

	if !result.Success {
		s.logger.Warn("Player deletion failed", "reason", result.Message, "player_id", req.PlayerId)
		return &pb.DeletePlayerResponse{
			Base: &pb.BaseResponse{
				Success:   false,
				Message:   result.Message,
				Code:      int32(result.ErrorCode),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	response := &pb.DeletePlayerResponse{
		Base: &pb.BaseResponse{
			Success:   true,
			Message:   "Player deleted successfully",
			Code:      0,
			Timestamp: time.Now().Unix(),
		},
	}

	s.logger.Info("Player deleted successfully", "player_id", req.PlayerId)
	return response, nil
}

// QueryPlayers 查询玩家
func (s *PlayerServiceImpl) QueryPlayers(ctx context.Context, req *pb.QueryPlayersRequest) (*pb.QueryPlayersResponse, error) {
	s.logger.Debug("Querying players")

	// 创建查询
	query := &player.QueryPlayersQuery{
		IncludeStats:    req.IncludeStats,
		IncludeMetadata: req.IncludeMetadata,
	}

	// 设置查询条件
	if req.Condition != nil {
		query.Filters = s.convertFilters(req.Condition.Filters)
		query.Sorts = s.convertSorts(req.Condition.Sorts)
		if req.Condition.Pagination != nil {
			query.Page = int(req.Condition.Pagination.Page)
			query.PageSize = int(req.Condition.Pagination.PageSize)
		}
	}

	// 执行查询
	result, err := handlers.ExecuteQueryTyped[*player.QueryPlayersQuery, *player.QueryPlayersResult](ctx, s.queryBus, query)
	if err != nil {
		s.logger.Error("Failed to query players", "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to query players: %v", err))
	}

	// 转换玩家信息
	players := make([]*pb.PlayerInfo, len(result.Players))
	for i, p := range result.Players {
		players[i] = s.convertToPlayerInfo(p)
	}

	response := &pb.QueryPlayersResponse{
		Base: &pb.BaseResponse{
			Success:   true,
			Message:   "Players queried successfully",
			Code:      0,
			Timestamp: time.Now().Unix(),
		},
		Players: players,
		Pagination: &pb.PaginationInfo{
			Page:       int32(result.Page),
			PageSize:   int32(result.PageSize),
			Total:      int32(result.Total),
			TotalPages: int32(result.TotalPages),
		},
	}

	s.logger.Debug("Players queried successfully", "count", len(players), "total", result.Total)
	return response, nil
}

// PlayerLogin 玩家登录
func (s *PlayerServiceImpl) PlayerLogin(ctx context.Context, req *pb.PlayerLoginRequest) (*pb.PlayerLoginResponse, error) {
	s.logger.Info("Player login", "player_id", req.PlayerId, "session_id", req.SessionId)

	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	// 创建命令
	cmd := &player.PlayerLoginCommand{
		PlayerID:  req.PlayerId,
		SessionID: req.SessionId,
		ClientIP:  req.ClientIp,
		UserAgent: req.UserAgent,
		LoginData: req.LoginData,
	}

	// 执行命令
	result, err := handlers.ExecuteTyped[*player.PlayerLoginCommand, *player.PlayerLoginResult](ctx, s.commandBus, cmd)
	if err != nil {
		s.logger.Error("Failed to login player", "error", err, "player_id", req.PlayerId)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to login player: %v", err))
	}

	if !result.Success {
		s.logger.Warn("Player login failed", "reason", result.Message, "player_id", req.PlayerId)
		return &pb.PlayerLoginResponse{
			Base: &pb.BaseResponse{
				Success:   false,
				Message:   result.Message,
				Code:      int32(result.ErrorCode),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	response := &pb.PlayerLoginResponse{
		Base: &pb.BaseResponse{
			Success:   true,
			Message:   "Player logged in successfully",
			Code:      0,
			Timestamp: time.Now().Unix(),
		},
		Player:       s.convertToPlayerInfo(result.Player),
		SessionToken: result.SessionToken,
		ExpiresAt:    result.ExpiresAt.Unix(),
	}

	s.logger.Info("Player logged in successfully", "player_id", req.PlayerId)
	return response, nil
}

// PlayerLogout 玩家登出
func (s *PlayerServiceImpl) PlayerLogout(ctx context.Context, req *pb.PlayerLogoutRequest) (*pb.PlayerLogoutResponse, error) {
	s.logger.Info("Player logout", "player_id", req.PlayerId, "session_id", req.SessionId)

	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	// 创建命令
	cmd := &player.PlayerLogoutCommand{
		PlayerID:  req.PlayerId,
		SessionID: req.SessionId,
		Reason:    req.Reason,
	}

	// 执行命令
	result, err := handlers.ExecuteTyped[*player.PlayerLogoutCommand, *player.PlayerLogoutResult](ctx, s.commandBus, cmd)
	if err != nil {
		s.logger.Error("Failed to logout player", "error", err, "player_id", req.PlayerId)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to logout player: %v", err))
	}

	response := &pb.PlayerLogoutResponse{
		Base: &pb.BaseResponse{
			Success:   result.Success,
			Message:   result.Message,
			Code:      int32(result.ErrorCode),
			Timestamp: time.Now().Unix(),
		},
	}

	s.logger.Info("Player logged out successfully", "player_id", req.PlayerId)
	return response, nil
}

// PlayerMove 玩家移动
func (s *PlayerServiceImpl) PlayerMove(ctx context.Context, req *pb.PlayerMoveRequest) (*pb.PlayerMoveResponse, error) {
	s.logger.Debug("Player move", "player_id", req.PlayerId)

	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	// 创建命令
	cmd := &player.MovePlayerCommand{
		PlayerID: req.PlayerId,
		FromPosition: player.Position{
			X: req.FromPosition.X,
			Y: req.FromPosition.Y,
			Z: req.FromPosition.Z,
		},
		Position: player.Position{
			X: req.ToPosition.X,
			Y: req.ToPosition.Y,
			Z: req.ToPosition.Z,
		},
		Speed:     req.Speed,
		Timestamp: time.Unix(req.Timestamp, 0),
	}

	// 执行命令
	result, err := handlers.ExecuteTyped[*player.MovePlayerCommand, *player.MovePlayerResult](ctx, s.commandBus, cmd)
	if err != nil {
		s.logger.Error("Failed to move player", "error", err, "player_id", req.PlayerId)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to move player: %v", err))
	}

	response := &pb.PlayerMoveResponse{
		Base: &pb.BaseResponse{
			Success:   result.Success,
			Message:   result.Message,
			Code:      0,
			Timestamp: time.Now().Unix(),
		},
		OldPosition: &pb.Position{
			X: result.OldPosition.X,
			Y: result.OldPosition.Y,
			Z: result.OldPosition.Z,
		},
		NewPosition: &pb.Position{
			X: result.NewPosition.X,
			Y: result.NewPosition.Y,
			Z: result.NewPosition.Z,
		},
		Success: result.Success,
		Reason:  result.Reason,
	}

	s.logger.Debug("Player moved successfully", "player_id", req.PlayerId)
	return response, nil
}

// 辅助方法

// convertToPlayerInfo 转换为gRPC玩家信息
func (s *PlayerServiceImpl) convertToPlayerInfo(p *player.PlayerInfo) *pb.PlayerInfo {
	if p == nil {
		return nil
	}

	return &pb.PlayerInfo{
		Id:     p.ID,
		Name:   p.Name,
		Email:  p.Email,
		Level:  int32(p.Level),
		Exp:    p.Exp,
		Status: p.Status,
		Position: &pb.Position{
			X: p.Position.X,
			Y: p.Position.Y,
			Z: p.Position.Z,
		},
		Stats: &pb.Stats{
			Hp:      int32(p.Stats.HP),
			MaxHp:   int32(p.Stats.MaxHP),
			Mp:      int32(p.Stats.MP),
			MaxMp:   int32(p.Stats.MaxMP),
			Attack:  int32(p.Stats.Attack),
			Defense: int32(p.Stats.Defense),
			Speed:   int32(p.Stats.Speed),
		},
		CreatedAt:  p.CreatedAt.Unix(),
		UpdatedAt:  p.UpdatedAt.Unix(),
		LastLogin:  p.LastLogin.Unix(),
		Metadata:   p.Metadata,
	}
}

// convertFilters 转换过滤条件
func (s *PlayerServiceImpl) convertFilters(filters []*pb.FilterCondition) []player.FilterCondition {
	result := make([]player.FilterCondition, len(filters))
	for i, f := range filters {
		result[i] = player.FilterCondition{
			Field:    f.Field,
			Operator: f.Operator,
			Values:   f.Values,
		}
	}
	return result
}

// convertSorts 转换排序条件
func (s *PlayerServiceImpl) convertSorts(sorts []*pb.SortCondition) []player.SortCondition {
	result := make([]player.SortCondition, len(sorts))
	for i, s := range sorts {
		result[i] = player.SortCondition{
			Field:     s.Field,
			Direction: s.Direction,
		}
	}
	return result
}

// 其他方法的占位符实现

func (s *PlayerServiceImpl) UpdatePlayerStatus(ctx context.Context, req *pb.UpdatePlayerStatusRequest) (*pb.UpdatePlayerStatusResponse, error) {
	// TODO: 实现更新玩家状态
	return nil, status.Error(codes.Unimplemented, "method UpdatePlayerStatus not implemented")
}

func (s *PlayerServiceImpl) UpdatePlayerStats(ctx context.Context, req *pb.UpdatePlayerStatsRequest) (*pb.UpdatePlayerStatsResponse, error) {
	// TODO: 实现更新玩家统计
	return nil, status.Error(codes.Unimplemented, "method UpdatePlayerStats not implemented")
}

func (s *PlayerServiceImpl) GetOnlinePlayers(ctx context.Context, req *pb.GetOnlinePlayersRequest) (*pb.GetOnlinePlayersResponse, error) {
	// TODO: 实现获取在线玩家
	return nil, status.Error(codes.Unimplemented, "method GetOnlinePlayers not implemented")
}

func (s *PlayerServiceImpl) UpdatePlayerExp(ctx context.Context, req *pb.UpdatePlayerExpRequest) (*pb.UpdatePlayerExpResponse, error) {
	// TODO: 实现更新玩家经验
	return nil, status.Error(codes.Unimplemented, "method UpdatePlayerExp not implemented")
}

func (s *PlayerServiceImpl) UpdatePlayerLevel(ctx context.Context, req *pb.UpdatePlayerLevelRequest) (*pb.UpdatePlayerLevelResponse, error) {
	// TODO: 实现更新玩家等级
	return nil, status.Error(codes.Unimplemented, "method UpdatePlayerLevel not implemented")
}

func (s *PlayerServiceImpl) BatchGetPlayers(ctx context.Context, req *pb.BatchGetPlayersRequest) (*pb.BatchGetPlayersResponse, error) {
	// TODO: 实现批量获取玩家
	return nil, status.Error(codes.Unimplemented, "method BatchGetPlayers not implemented")
}

func (s *PlayerServiceImpl) SearchPlayers(ctx context.Context, req *pb.SearchPlayersRequest) (*pb.SearchPlayersResponse, error) {
	// TODO: 实现搜索玩家
	return nil, status.Error(codes.Unimplemented, "method SearchPlayers not implemented")
}

func (s *PlayerServiceImpl) GetPlayerRanking(ctx context.Context, req *pb.GetPlayerRankingRequest) (*pb.GetPlayerRankingResponse, error) {
	// TODO: 实现获取玩家排行榜
	return nil, status.Error(codes.Unimplemented, "method GetPlayerRanking not implemented")
}

func (s *PlayerServiceImpl) GetPlayerActivity(ctx context.Context, req *pb.GetPlayerActivityRequest) (*pb.GetPlayerActivityResponse, error) {
	// TODO: 实现获取玩家活动记录
	return nil, status.Error(codes.Unimplemented, "method GetPlayerActivity not implemented")
}

func (s *PlayerServiceImpl) RecordPlayerActivity(ctx context.Context, req *pb.RecordPlayerActivityRequest) (*pb.RecordPlayerActivityResponse, error) {
	// TODO: 实现记录玩家活动
	return nil, status.Error(codes.Unimplemented, "method RecordPlayerActivity not implemented")
}