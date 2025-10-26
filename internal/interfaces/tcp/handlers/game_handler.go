package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	appServices "greatestworks/internal/application/services"
	"greatestworks/internal/domain/character"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/interfaces/tcp/connection"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// GameHandler 游戏处理器
type GameHandler struct {
	logger           logging.Logger
	connManager      *connection.Manager
	mapService       *appServices.MapService
	fightService     *appServices.FightService
	characterService *appServices.CharacterService
}

// NewGameHandler 创建游戏处理器
func NewGameHandler(commandBus interface{}, queryBus interface{}, connManager *connection.Manager, logger logging.Logger) *GameHandler {
	return &GameHandler{
		logger:      logger,
		connManager: connManager,
	}
}

// SetMapService 注入地图服务
func (h *GameHandler) SetMapService(ms *appServices.MapService) { h.mapService = ms }

// SetFightService 注入战斗服务
func (h *GameHandler) SetFightService(fs *appServices.FightService) { h.fightService = fs }

// SetCharacterService 注入角色服务
func (h *GameHandler) SetCharacterService(cs *appServices.CharacterService) { h.characterService = cs }

// HandleMessage 处理消息
func (h *GameHandler) HandleMessage(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理游戏消息", map[string]interface{}{
		"message_type": message.Header.MessageType,
		"player_id":    message.Header.PlayerID,
	})

	switch message.Header.MessageType {
	case protocol.MsgPlayerLogin:
		return h.handlePlayerLogin(session, message)
	case protocol.MsgPlayerLogout:
		return h.handlePlayerLogout(session, message)
	case protocol.MsgPlayerMove:
		return h.handlePlayerMove(session, message)
	case protocol.MsgBattleSkill:
		return h.handleSkillCast(session, message)
	case protocol.MsgPlayerStatus:
		return h.handlePlayerChat(session, message)
	case protocol.MsgPlayerStats:
		return h.handlePlayerAction(session, message)
	case protocol.MsgChatMessage:
		return h.handleChatMessage(session, message)
	case protocol.MsgTeamCreate:
		return h.handleTeamCreate(session, message)
	case protocol.MsgTeamJoin:
		return h.handleTeamJoin(session, message)
	case protocol.MsgTeamLeave:
		return h.handleTeamLeave(session, message)
	case protocol.MsgTeamInfo:
		return h.handleTeamInfo(session, message)
	default:
		return fmt.Errorf("未知的消息类型: %d", message.Header.MessageType)
	}
}

// handlePlayerLogin 处理玩家登录
func (h *GameHandler) handlePlayerLogin(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理玩家登录", map[string]interface{}{
		"player_id": message.Header.PlayerID,
	})

	// 解析请求负载
	var req protocol.PlayerLoginRequest
	if payloadMap, ok := message.Payload.(map[string]interface{}); ok {
		if b, err := json.Marshal(payloadMap); err == nil {
			_ = json.Unmarshal(b, &req)
		}
	}

	// 简化绑定：假设协议中的 PlayerID 为实体ID或可转换的数字ID
	var entityID int32
	var characterID int64
	if req.PlayerID != "" {
		if id64, err := strconv.ParseInt(req.PlayerID, 10, 64); err == nil {
			entityID = int32(id64)
			characterID = id64
		}
	}
	if entityID == 0 && message.Header.PlayerID != 0 {
		entityID = int32(message.Header.PlayerID)
		characterID = int64(message.Header.PlayerID)
	}

	// 绑定会话与玩家
	if entityID != 0 && h.connManager != nil {
		h.connManager.BindPlayer(entityID, session)
	}
	session.SetUserID(req.PlayerID)

	// 推断地图ID与位置：优先从角色服务加载持久化位置
	var mapID int32 = 1
	var x, y, z float32 = 0, 0, 0
	if h.characterService != nil && characterID != 0 {
		if dbChar, err := h.characterService.GetCharacter(context.Background(), characterID); err == nil && dbChar != nil {
			mapID = dbChar.MapID
			x, y, z = dbChar.PositionX, dbChar.PositionY, dbChar.PositionZ
		}
	}
	// 允许客户端覆盖map_id（可选协议字段）
	if payloadMap, ok := message.Payload.(map[string]interface{}); ok {
		if v, ok := payloadMap["map_id"]; ok {
			switch vv := v.(type) {
			case float64:
				mapID = int32(vv)
			case string:
				if id64, err := strconv.ParseInt(vv, 10, 32); err == nil {
					mapID = int32(id64)
				}
			}
		}
	}
	session.SetGroupID(fmt.Sprintf("map:%d", mapID))

	// 确保地图加载并注册入地图（以便后续移动/AOI广播可用）
	if h.mapService != nil && entityID != 0 {
		_ = h.mapService.LoadMap(context.Background(), mapID)
		e := character.NewEntity(character.EntityID(entityID), character.EntityTypePlayer, 1, character.NewVector3(0, 0, 0), character.NewVector3(0, 0, 1))
		_ = h.mapService.EnterMap(context.Background(), e, mapID, x, y, z)
	}

	// 构造登录响应
	resp := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   message.Header.MessageID,
			MessageType: uint32(protocol.MsgPlayerLogin),
			Flags:       protocol.FlagResponse,
			PlayerID:    message.Header.PlayerID,
			Timestamp:   time.Now().Unix(),
			Sequence:    0,
		},
		Payload: protocol.PlayerLoginResponse{
			BaseResponse: protocol.NewBaseResponse(true, "login ok"),
			SessionID:    session.ID,
			ServerTime:   time.Now().Unix(),
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("序列化登录响应失败: %w", err)
	}
	return session.Send(data)

}

// handlePlayerLogout 处理玩家登出
func (h *GameHandler) handlePlayerLogout(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理玩家登出", map[string]interface{}{
		"player_id": message.Header.PlayerID,
	})

	// 清理绑定
	if h.connManager != nil {
		// 尝试从会话反查实体ID
		if entityID, ok := h.connManager.GetPlayerBySession(session.ID); ok {
			// 尝试从GroupID解析地图，并从地图移除实体
			var mapID int32 = 1
			if h.mapService != nil {
				if gid := session.GetGroupID(); gid != "" {
					if len(gid) > 4 && gid[:4] == "map:" {
						if v, err := strconv.ParseInt(gid[4:], 10, 32); err == nil {
							mapID = int32(v)
						}
					}
				}
				// 从地图中获取最终位置并保存
				if h.characterService != nil && mapID > 0 {
					if m, err := h.mapService.GetMap(mapID); err == nil && m != nil {
						if e := m.GetEntity(character.EntityID(entityID)); e != nil {
							pos := e.Position()
							_ = h.characterService.UpdateLastLocation(
								context.Background(), int64(entityID), mapID, pos.X, pos.Y, pos.Z,
							)
						}
					}
				}
				_ = h.mapService.LeaveMapByID(context.Background(), mapID, entityID)
			}
			h.connManager.UnbindPlayer(entityID)
		} else if message.Header.PlayerID != 0 {
			h.connManager.UnbindPlayer(int32(message.Header.PlayerID))
		}
	}

	return nil
}

// handlePlayerMove 处理玩家移动
func (h *GameHandler) handlePlayerMove(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理玩家移动", map[string]interface{}{
		"player_id": message.Header.PlayerID,
	})

	if h.mapService == nil || h.connManager == nil {
		return fmt.Errorf("map service or connection manager not ready")
	}

	// 解析请求
	var req protocol.PlayerMoveRequest
	if payloadMap, ok := message.Payload.(map[string]interface{}); ok {
		if b, err := json.Marshal(payloadMap); err == nil {
			_ = json.Unmarshal(b, &req)
		}
	}

	// 获取玩家绑定的实体ID
	entityID, ok := h.connManager.GetPlayerBySession(session.ID)
	if !ok {
		// 回退到Header PlayerID
		entityID = int32(message.Header.PlayerID)
	}
	if entityID == 0 {
		return fmt.Errorf("no bound entity for session")
	}

	// 推断mapID: 从会话GroupID形如"map:<id>"中解析；否则默认1
	var mapID int32 = 1
	if gid := session.GetGroupID(); gid != "" {
		if len(gid) > 4 && gid[:4] == "map:" {
			if v, err := strconv.ParseInt(gid[4:], 10, 32); err == nil {
				mapID = int32(v)
			}
		}
	}

	// 执行位置更新
	if err := h.mapService.UpdatePositionByID(
		context.Background(),
		mapID, entityID, float32(req.Position.X), float32(req.Position.Y), float32(req.Position.Z),
	); err != nil {
		return err
	}

	// 回执
	resp := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   message.Header.MessageID,
			MessageType: uint32(protocol.MsgPlayerMove),
			Flags:       protocol.FlagResponse,
			PlayerID:    message.Header.PlayerID,
			Timestamp:   time.Now().Unix(),
		},
		Payload: protocol.PlayerMoveResponse{BaseResponse: protocol.NewBaseResponse(true, "move ok")},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("序列化移动响应失败: %w", err)
	}
	return session.Send(data)
}

// handlePlayerChat 处理玩家聊天
func (h *GameHandler) handlePlayerChat(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理玩家聊天", map[string]interface{}{
		"player_id": message.Header.PlayerID,
	})

	// TODO: 实现玩家聊天逻辑
	// 1. 验证聊天内容
	// 2. 过滤敏感词
	// 3. 广播聊天消息

	return nil
}

// handlePlayerAction 处理玩家动作
func (h *GameHandler) handlePlayerAction(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理玩家动作", logging.Fields{
		"player_id": message.Header.PlayerID,
	})

	// TODO: 实现玩家动作逻辑
	// 1. 验证动作合法性
	// 2. 执行动作
	// 3. 更新游戏状态

	return nil
}

// handleSkillCast 处理技能释放（占位实现，待接FightService）
func (h *GameHandler) handleSkillCast(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理技能释放", logging.Fields{
		"player_id":    message.Header.PlayerID,
		"message_id":   message.Header.MessageID,
		"message_type": message.Header.MessageType,
	})

	// 解析payload允许简单格式: { "skill_id": string|number, "target_id": string|number }
	var skillIDStr string
	var targetIDStr string
	var skillID int32
	var targetID int32
	if payloadMap, ok := message.Payload.(map[string]interface{}); ok {
		if v, ok := payloadMap["skill_id"]; ok {
			switch vv := v.(type) {
			case string:
				skillIDStr = vv
			case float64:
				skillID = int32(vv)
			}
		}
		if v, ok := payloadMap["target_id"]; ok {
			switch vv := v.(type) {
			case string:
				targetIDStr = vv
			case float64:
				targetID = int32(vv)
			}
		}
	}

	if skillID == 0 && skillIDStr != "" {
		if id64, err := strconv.ParseInt(skillIDStr, 10, 32); err == nil {
			skillID = int32(id64)
		}
	}
	if targetID == 0 && targetIDStr != "" {
		if id64, err := strconv.ParseInt(targetIDStr, 10, 32); err == nil {
			targetID = int32(id64)
		}
	}

	// 获取施法者实体ID
	casterID, ok := h.connManager.GetPlayerBySession(session.ID)
	if !ok {
		casterID = int32(message.Header.PlayerID)
	}

	// 调用战斗服务计算伤害
	var castResult *appServices.SkillCastResult
	var castErr error
	if h.fightService != nil {
		castResult, castErr = h.fightService.CastSkillByID(context.Background(), casterID, targetID, skillID)
	}

	// 回执
	resp := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   message.Header.MessageID,
			MessageType: uint32(protocol.MsgBattleSkill),
			Flags:       protocol.FlagResponse,
			PlayerID:    message.Header.PlayerID,
			Timestamp:   time.Now().Unix(),
		},
		Payload: map[string]interface{}{
			"result": func() interface{} {
				if castErr != nil {
					return protocol.NewBaseResponse(false, castErr.Error())
				}
				return protocol.NewBaseResponse(true, "cast ok")
			}(),
			"skill_id": func() interface{} {
				if skillID != 0 {
					return skillID
				}
				return skillIDStr
			}(),
			"target_id": func() interface{} {
				if targetID != 0 {
					return targetID
				}
				return targetIDStr
			}(),
			"caster_id": casterID,
			"damage": func() int32 {
				if castResult != nil {
					return castResult.Damage
				}
				return 0
			}(),
			"is_critical": func() bool {
				if castResult != nil {
					return castResult.IsCritical
				}
				return false
			}(),
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("序列化技能响应失败: %w", err)
	}
	if err := session.Send(data); err != nil {
		return err
	}

	// 施法广播给AOI内玩家
	if h.mapService != nil {
		// 推断mapID
		var mapID int32 = 1
		if gid := session.GetGroupID(); gid != "" {
			if len(gid) > 4 && gid[:4] == "map:" {
				if v, err := strconv.ParseInt(gid[4:], 10, 32); err == nil {
					mapID = int32(v)
				}
			}
		}
		payload := map[string]interface{}{
			"caster_id": casterID,
			"skill_id": func() interface{} {
				if skillID != 0 {
					return skillID
				}
				return skillIDStr
			}(),
			"target_id": func() interface{} {
				if targetID != 0 {
					return targetID
				}
				return targetIDStr
			}(),
			"ts": time.Now().UnixMilli(),
		}
		if m, err := h.mapService.GetMap(mapID); err == nil {
			ents := m.GetAllEntities()
			recvs := make([]character.EntityID, 0, len(ents))
			for _, e := range ents {
				recvs = append(recvs, e.ID())
			}
			m.BroadcastTo(recvs, "skill_cast", payload)
		}
	}

	return nil
}

// (removed unused helper)

// handleChatMessage 处理聊天消息
func (h *GameHandler) handleChatMessage(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理聊天消息", logging.Fields{
		"player_id":    message.Header.PlayerID,
		"message_id":   message.Header.MessageID,
		"message_type": message.Header.MessageType,
	})

	// 简单回执响应
	resp := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   message.Header.MessageID,
			MessageType: uint32(protocol.MsgChatMessage),
			Flags:       protocol.FlagResponse,
			PlayerID:    message.Header.PlayerID,
			Timestamp:   message.Header.Timestamp,
			Sequence:    0,
		},
		Payload: protocol.NewBaseResponse(true, "chat received"),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("序列化聊天响应失败: %w", err)
	}
	return session.Send(data)
}

// handleTeamCreate 处理创建队伍
func (h *GameHandler) handleTeamCreate(session *connection.Session, message *protocol.Message) error {
	return h.replyOK(session, message, uint32(protocol.MsgTeamCreate), "team created")
}

// handleTeamJoin 处理加入队伍
func (h *GameHandler) handleTeamJoin(session *connection.Session, message *protocol.Message) error {
	return h.replyOK(session, message, uint32(protocol.MsgTeamJoin), "team joined")
}

// handleTeamLeave 处理离开队伍
func (h *GameHandler) handleTeamLeave(session *connection.Session, message *protocol.Message) error {
	return h.replyOK(session, message, uint32(protocol.MsgTeamLeave), "team left")
}

// handleTeamInfo 处理队伍信息
func (h *GameHandler) handleTeamInfo(session *connection.Session, message *protocol.Message) error {
	return h.replyOK(session, message, uint32(protocol.MsgTeamInfo), "team info")
}

// replyOK 通用成功回执
func (h *GameHandler) replyOK(session *connection.Session, message *protocol.Message, msgType uint32, text string) error {
	resp := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   message.Header.MessageID,
			MessageType: msgType,
			Flags:       protocol.FlagResponse,
			PlayerID:    message.Header.PlayerID,
			Timestamp:   message.Header.Timestamp,
			Sequence:    0,
		},
		Payload: protocol.NewBaseResponse(true, text),
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("序列化响应失败: %w", err)
	}
	return session.Send(data)
}

// SendResponse 发送响应
func (h *GameHandler) SendResponse(playerID uint64, responseType uint32, data interface{}) error {
	_ = &protocol.Message{
		Header: protocol.MessageHeader{
			MessageType: responseType,
			PlayerID:    playerID,
		},
		Payload: data,
	}

	// TODO: 实现发送响应逻辑
	h.logger.Info("发送响应", logging.Fields{
		"player_id":     playerID,
		"response_type": responseType,
	})

	return nil
}
