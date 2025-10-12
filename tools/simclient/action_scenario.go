package simclient

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"greatestworks/internal/infrastructure/logging"
	tcpProtocol "greatestworks/internal/interfaces/tcp/protocol"
)

// ActionScenario executes a predefined list of message actions for targeted testing.
type ActionScenario struct {
	cfg     ScenarioConfig
	logger  logging.Logger
	actions []resolvedAction
}

type resolvedAction struct {
	name           string
	messageType    uint32
	flags          uint16
	expectResponse bool
	pause          time.Duration
	repeat         int
}

// NewActionScenario constructs an action-driven scenario.
func NewActionScenario(cfg ScenarioConfig, logger logging.Logger) (*ActionScenario, error) {
	localCfg := cfg
	features := make([]string, 0, len(localCfg.Features))
	for _, feature := range localCfg.Features {
		if feature == "" {
			continue
		}
		features = append(features, normalizeKey(feature))
	}

	// Allow specifying a single feature via type when no explicit features are given.
	scenarioType := normalizeKey(localCfg.Type)
	if len(features) == 0 && scenarioType != "" && scenarioType != "basic" {
		features = append(features, scenarioType)
	}

	actions := make([]resolvedAction, 0)

	for _, feature := range features {
		steps, ok := featureLibrary[feature]
		if !ok {
			return nil, fmt.Errorf("unknown feature %q", feature)
		}
		for _, step := range steps {
			action, err := resolveAction(step, &localCfg)
			if err != nil {
				return nil, fmt.Errorf("feature %s: %w", feature, err)
			}
			actions = append(actions, action)
		}
	}

	for _, rawAction := range localCfg.Actions {
		action, err := resolveAction(rawAction, &localCfg)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	if len(actions) == 0 {
		return nil, fmt.Errorf("no feature actions configured")
	}

	return &ActionScenario{
		cfg:     localCfg,
		logger:  logger,
		actions: actions,
	}, nil
}

// Execute runs each configured action sequentially.
func (s *ActionScenario) Execute(ctx context.Context, client *SimulatorClient) (*ScenarioResult, error) {
	result := &ScenarioResult{ScenarioName: s.cfg.Name, StartedAt: time.Now()}
	defer func() {
		result.CompletedAt = time.Now()
	}()

	if token, err := authenticateAndRecord(ctx, client, result, s.logger); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
		s.logger.Warn("authentication failed but continuing", logging.Fields{"error": err})
	} else if token != "" {
		s.logger.Debug("received auth token", logging.Fields{"token_length": len(token)})
	}

	start := time.Now()
	conn, err := client.ConnectGateway(ctx)
	connectLatency := time.Since(start)
	fields := map[string]interface{}{}
	if err == nil && conn != nil {
		fields["remote"] = conn.RemoteAddr().String()
	}
	result.Record("gateway.connect", connectLatency, err, fields)
	if err != nil {
		return result, err
	}
	defer conn.Close()

	for _, action := range s.actions {
		for i := 0; i < action.repeat; i++ {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				result.Record("scenario.cancelled", 0, err, nil)
				return result, err
			default:
			}

			stepName := action.name
			if action.repeat > 1 {
				stepName = fmt.Sprintf("%s#%d", stepName, i+1)
			}

			if err := s.executeAction(result, client, conn, stepName, action); err != nil {
				if s.cfg.StopOnError {
					return result, err
				}
				s.logger.Warn("action execution failed", logging.Fields{
					"action": stepName,
					"error":  err,
				})
			}

			if action.pause > 0 {
				select {
				case <-ctx.Done():
					err := ctx.Err()
					result.Record("scenario.cancelled", 0, err, nil)
					return result, err
				case <-time.After(action.pause):
				}
			}
		}
	}

	return result, nil
}

// Name returns the scenario name for logging and metrics.
func (s *ActionScenario) Name() string {
	return s.cfg.Name
}

func (s *ActionScenario) executeAction(result *ScenarioResult, client *SimulatorClient, conn net.Conn, name string, action resolvedAction) error {
	start := time.Now()
	fields := map[string]interface{}{
		"message_type": action.messageType,
		"flags":        describeFlags(action.flags),
	}

	messageID, err := client.SendGatewayMessage(conn, action.messageType, action.flags)
	if err == nil {
		fields["message_id"] = messageID
		if action.expectResponse {
			received, respErr := client.TryRead(conn)
			switch {
			case respErr != nil:
				err = respErr
			case !received:
				err = fmt.Errorf("no response received")
			default:
				fields["response"] = "received"
			}
		}
	}

	result.Record(name, time.Since(start), err, fields)
	return err
}

func resolveAction(cfg ScenarioActionConfig, scenario *ScenarioConfig) (resolvedAction, error) {
	actionName := cfg.Name
	if actionName == "" {
		actionName = cfg.Message
	}
	if actionName == "" {
		return resolvedAction{}, fmt.Errorf("action name or message must be provided")
	}

	msgType, err := resolveMessageType(cfg.Message)
	if err != nil {
		return resolvedAction{}, fmt.Errorf("action %s: %w", actionName, err)
	}

	flags, err := resolveFlags(cfg.Flags)
	if err != nil {
		return resolvedAction{}, fmt.Errorf("action %s: %w", actionName, err)
	}

	expect := false
	if cfg.ExpectResponse != nil {
		expect = *cfg.ExpectResponse
	}

	pause := cfg.Pause.AsDuration()
	if pause == 0 {
		pause = scenario.ActionInterval.AsDuration()
	}

	repeat := cfg.Repeat
	if repeat <= 0 {
		repeat = 1
	}

	return resolvedAction{
		name:           actionName,
		messageType:    msgType,
		flags:          flags,
		expectResponse: expect,
		pause:          pause,
		repeat:         repeat,
	}, nil
}

func resolveMessageType(name string) (uint32, error) {
	key := normalizeKey(name)
	if key == "" {
		return 0, fmt.Errorf("message name is required")
	}

	if id, ok := messageNameToID[key]; ok {
		return id, nil
	}

	if strings.HasPrefix(key, "0x") {
		value, err := strconv.ParseUint(key, 0, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid hex message %q: %w", name, err)
		}
		return uint32(value), nil
	}

	if value, err := strconv.ParseUint(key, 10, 32); err == nil {
		return uint32(value), nil
	}

	return 0, fmt.Errorf("unknown message name %q", name)
}

func resolveFlags(flagNames []string) (uint16, error) {
	if len(flagNames) == 0 {
		return tcpProtocol.FlagRequest, nil
	}

	var mask uint16
	for _, name := range flagNames {
		key := normalizeKey(name)
		if key == "" {
			continue
		}
		value, ok := messageFlagNameToValue[key]
		if !ok {
			return 0, fmt.Errorf("unknown flag %q", name)
		}
		mask |= value
	}

	if mask == 0 {
		mask = tcpProtocol.FlagRequest
	}

	return mask, nil
}

func describeFlags(mask uint16) string {
	if mask == 0 {
		return ""
	}

	names := make([]string, 0, len(messageFlagNameToValue))
	for name, value := range messageFlagNameToValue {
		if mask&value != 0 {
			names = append(names, name)
		}
	}

	sort.Strings(names)
	return strings.Join(names, ",")
}

func normalizeKey(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

var messageNameToID = map[string]uint32{
	"system.heartbeat": tcpProtocol.MsgHeartbeat,
	"system.handshake": tcpProtocol.MsgHandshake,
	"system.auth":      tcpProtocol.MsgAuth,

	"player.login":  tcpProtocol.MsgPlayerLogin,
	"player.logout": tcpProtocol.MsgPlayerLogout,
	"player.move":   tcpProtocol.MsgPlayerMove,
	"player.info":   tcpProtocol.MsgPlayerInfo,
	"player.create": tcpProtocol.MsgPlayerCreate,
	"player.update": tcpProtocol.MsgPlayerUpdate,
	"player.delete": tcpProtocol.MsgPlayerDelete,
	"player.level":  tcpProtocol.MsgPlayerLevelUp,
	"player.exp":    tcpProtocol.MsgPlayerExpGain,
	"player.status": tcpProtocol.MsgPlayerStatus,
	"player.stats":  tcpProtocol.MsgPlayerStats,
	"player.sync":   tcpProtocol.MsgPlayerStatusSync,

	"battle.create": tcpProtocol.MsgCreateBattle,
	"battle.join":   tcpProtocol.MsgJoinBattle,
	"battle.start":  tcpProtocol.MsgStartBattle,
	"battle.action": tcpProtocol.MsgBattleAction,
	"battle.leave":  tcpProtocol.MsgLeaveBattle,
	"battle.result": tcpProtocol.MsgBattleResult,
	"battle.status": tcpProtocol.MsgBattleStatus,
	"battle.round":  tcpProtocol.MsgBattleRound,
	"battle.skill":  tcpProtocol.MsgBattleSkill,
	"battle.damage": tcpProtocol.MsgBattleDamage,

	"query.player_info":    tcpProtocol.MsgGetPlayerInfo,
	"query.online_players": tcpProtocol.MsgGetOnlinePlayers,
	"query.battle_info":    tcpProtocol.MsgGetBattleInfo,
	"query.player_stats":   tcpProtocol.MsgGetPlayerStats,
	"query.battle_list":    tcpProtocol.MsgGetBattleList,
	"query.rankings":       tcpProtocol.MsgGetRankings,
	"query.server_info":    tcpProtocol.MsgGetServerInfo,

	"pet.summon":    tcpProtocol.MsgPetSummon,
	"pet.dismiss":   tcpProtocol.MsgPetDismiss,
	"pet.info":      tcpProtocol.MsgPetInfo,
	"pet.move":      tcpProtocol.MsgPetMove,
	"pet.action":    tcpProtocol.MsgPetAction,
	"pet.level_up":  tcpProtocol.MsgPetLevelUp,
	"pet.evolution": tcpProtocol.MsgPetEvolution,
	"pet.train":     tcpProtocol.MsgPetTrain,
	"pet.feed":      tcpProtocol.MsgPetFeed,
	"pet.status":    tcpProtocol.MsgPetStatus,

	"building.create":  tcpProtocol.MsgBuildingCreate,
	"building.upgrade": tcpProtocol.MsgBuildingUpgrade,
	"building.destroy": tcpProtocol.MsgBuildingDestroy,
	"building.info":    tcpProtocol.MsgBuildingInfo,
	"building.produce": tcpProtocol.MsgBuildingProduce,
	"building.collect": tcpProtocol.MsgBuildingCollect,
	"building.repair":  tcpProtocol.MsgBuildingRepair,
	"building.status":  tcpProtocol.MsgBuildingStatus,

	"social.chat":           tcpProtocol.MsgChatMessage,
	"social.friend_request": tcpProtocol.MsgFriendRequest,
	"social.friend_accept":  tcpProtocol.MsgFriendAccept,
	"social.friend_reject":  tcpProtocol.MsgFriendReject,
	"social.friend_remove":  tcpProtocol.MsgFriendRemove,
	"social.friend_list":    tcpProtocol.MsgFriendList,
	"social.guild_create":   tcpProtocol.MsgGuildCreate,
	"social.guild_join":     tcpProtocol.MsgGuildJoin,
	"social.guild_leave":    tcpProtocol.MsgGuildLeave,
	"social.guild_info":     tcpProtocol.MsgGuildInfo,
	"social.team_create":    tcpProtocol.MsgTeamCreate,
	"social.team_join":      tcpProtocol.MsgTeamJoin,
	"social.team_leave":     tcpProtocol.MsgTeamLeave,
	"social.team_info":      tcpProtocol.MsgTeamInfo,

	"item.use":       tcpProtocol.MsgItemUse,
	"item.equip":     tcpProtocol.MsgItemEquip,
	"item.unequip":   tcpProtocol.MsgItemUnequip,
	"item.drop":      tcpProtocol.MsgItemDrop,
	"item.pickup":    tcpProtocol.MsgItemPickup,
	"item.trade":     tcpProtocol.MsgItemTrade,
	"item.inventory": tcpProtocol.MsgInventoryInfo,
	"item.info":      tcpProtocol.MsgItemInfo,
	"item.craft":     tcpProtocol.MsgItemCraft,
	"item.enhance":   tcpProtocol.MsgItemEnhance,

	"quest.accept":   tcpProtocol.MsgQuestAccept,
	"quest.complete": tcpProtocol.MsgQuestComplete,
	"quest.cancel":   tcpProtocol.MsgQuestCancel,
	"quest.progress": tcpProtocol.MsgQuestProgress,
	"quest.list":     tcpProtocol.MsgQuestList,
	"quest.info":     tcpProtocol.MsgQuestInfo,
	"quest.reward":   tcpProtocol.MsgQuestReward,
}

var messageFlagNameToValue = map[string]uint16{
	"request":    tcpProtocol.FlagRequest,
	"response":   tcpProtocol.FlagResponse,
	"error":      tcpProtocol.FlagError,
	"async":      tcpProtocol.FlagAsync,
	"broadcast":  tcpProtocol.FlagBroadcast,
	"encrypted":  tcpProtocol.FlagEncrypted,
	"compressed": tcpProtocol.FlagCompressed,
}
