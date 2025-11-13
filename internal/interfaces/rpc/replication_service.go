// Package rpc 副本RPC服务
package rpc

import (
	"context"
	"greatestworks/internal/application/services"
	"greatestworks/internal/infrastructure/logging"
	"time"
)

// ReplicationRPCService 副本RPC服务
type ReplicationRPCService struct {
	app    *services.ReplicationService
	logger logging.Logger
}

// NewReplicationRPCService 构造
func NewReplicationRPCService(app *services.ReplicationService, logger logging.Logger) *ReplicationRPCService {
	return &ReplicationRPCService{app: app, logger: logger}
}

// CreateInstance
type CreateInstanceRPCRequest struct {
	TemplateID    string
	InstanceType  int
	OwnerPlayerID string
	OwnerName     string
	OwnerLevel    int
	MaxPlayers    int
	Difficulty    int
	LifetimeSec   int64
}

type CreateInstanceRPCResponse struct {
	Instance *services.InstanceInfoDTO
	Error    string
}

func (s *ReplicationRPCService) CreateInstance(req CreateInstanceRPCRequest, resp *CreateInstanceRPCResponse) error {
	ctx := context.Background()
	dto, err := s.app.CreateInstance(ctx, &services.CreateInstanceCommand{
		TemplateID:    req.TemplateID,
		InstanceType:  req.InstanceType,
		OwnerPlayerID: req.OwnerPlayerID,
		OwnerName:     req.OwnerName,
		OwnerLevel:    req.OwnerLevel,
		MaxPlayers:    req.MaxPlayers,
		Difficulty:    req.Difficulty,
		Lifetime:      time.Duration(req.LifetimeSec) * time.Second,
	})
	if err != nil {
		resp.Error = err.Error()
		return nil
	}
	resp.Instance = dto
	return nil
}

// JoinInstance
type JoinInstanceRPCRequest struct {
	InstanceID string
	PlayerID   string
	PlayerName string
	Level      int
	Role       string
}

type SimpleRPCResponse struct{ Error string }

func (s *ReplicationRPCService) JoinInstance(req JoinInstanceRPCRequest, resp *SimpleRPCResponse) error {
	ctx := context.Background()
	err := s.app.JoinInstance(ctx, &services.JoinInstanceCommand{
		InstanceID: req.InstanceID,
		PlayerID:   req.PlayerID,
		PlayerName: req.PlayerName,
		Level:      req.Level,
		Role:       req.Role,
	})
	if err != nil {
		resp.Error = err.Error()
	}
	return nil
}

// LeaveInstance
type LeaveInstanceRPCRequest struct{ InstanceID, PlayerID string }

func (s *ReplicationRPCService) LeaveInstance(req LeaveInstanceRPCRequest, resp *SimpleRPCResponse) error {
	ctx := context.Background()
	err := s.app.LeaveInstance(ctx, &services.LeaveInstanceCommand{InstanceID: req.InstanceID, PlayerID: req.PlayerID})
	if err != nil {
		resp.Error = err.Error()
	}
	return nil
}

// GetInstanceInfo
type GetInstanceInfoRPCRequest struct{ InstanceID string }

type GetInstanceInfoRPCResponse struct {
	Instance *services.InstanceInfoDTO
	Error    string
}

func (s *ReplicationRPCService) GetInstanceInfo(req GetInstanceInfoRPCRequest, resp *GetInstanceInfoRPCResponse) error {
	ctx := context.Background()
	dto, err := s.app.GetInstanceInfo(ctx, req.InstanceID)
	if err != nil {
		resp.Error = err.Error()
		return nil
	}
	resp.Instance = dto
	return nil
}

// ListActiveInstances

type ListActiveInstancesRPCRequest struct{}

type ListActiveInstancesRPCResponse struct {
	Instances []*services.InstanceInfoDTO
	Error     string
}

func (s *ReplicationRPCService) ListActiveInstances(_ ListActiveInstancesRPCRequest, resp *ListActiveInstancesRPCResponse) error {
	ctx := context.Background()
	dtos, err := s.app.ListActiveInstances(ctx)
	if err != nil {
		resp.Error = err.Error()
		return nil
	}
	resp.Instances = dtos
	return nil
}

// CleanupExpiredInstances

type CleanupExpiredInstancesRPCRequest struct{}

type CleanupExpiredInstancesRPCResponse struct {
	Count int
	Error string
}

func (s *ReplicationRPCService) CleanupExpiredInstances(_ CleanupExpiredInstancesRPCRequest, resp *CleanupExpiredInstancesRPCResponse) error {
	ctx := context.Background()
	count, err := s.app.CleanupExpiredInstances(ctx)
	if err != nil {
		resp.Error = err.Error()
		return nil
	}
	resp.Count = count
	return nil
}
