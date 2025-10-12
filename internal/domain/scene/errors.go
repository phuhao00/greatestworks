package scene

import "errors"

var (
	// 场景相关错误
	ErrSceneNotFound      = errors.New("scene not found")
	ErrSceneNotActive     = errors.New("scene is not active")
	ErrSceneFull          = errors.New("scene is full")
	ErrSceneClosed        = errors.New("scene is closed")
	ErrSceneMaintenance   = errors.New("scene is under maintenance")
	ErrInvalidSceneType   = errors.New("invalid scene type")
	ErrInvalidSceneStatus = errors.New("invalid scene status")
	ErrSceneAlreadyExists = errors.New("scene already exists")

	// 玩家相关错误
	ErrPlayerNotInScene     = errors.New("player not in scene")
	ErrPlayerAlreadyInScene = errors.New("player already in scene")
	ErrPlayerDead           = errors.New("player is dead")
	ErrPlayerAFK            = errors.New("player is AFK")
	ErrPlayerInCombat       = errors.New("player is in combat")
	ErrPlayerTrading        = errors.New("player is trading")

	// 实体相关错误
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityNotActive     = errors.New("entity is not active")
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrInvalidEntityType   = errors.New("invalid entity type")
	ErrEntityDead          = errors.New("entity is dead")

	// 位置相关错误
	ErrInvalidPosition  = errors.New("invalid position")
	ErrPositionOccupied = errors.New("position is occupied")
	ErrOutOfBounds      = errors.New("position is out of bounds")
	ErrTooFarAway       = errors.New("target is too far away")
	ErrCannotMove       = errors.New("cannot move to target position")

	// 怪物相关错误
	ErrMonsterNotFound      = errors.New("monster not found")
	ErrMonsterAlreadyExists = errors.New("monster already exists")
	ErrMonsterDead          = errors.New("monster is dead")
	ErrMonsterRespawning    = errors.New("monster is respawning")
	ErrInvalidMonsterType   = errors.New("invalid monster type")

	// NPC相关错误
	ErrNPCNotFound     = errors.New("npc not found")
	ErrNPCNotAvailable = errors.New("npc is not available")
	ErrNPCBusy         = errors.New("npc is busy")
	ErrNPCDead         = errors.New("npc is dead")
	ErrInvalidNPCType  = errors.New("invalid npc type")

	// 物品相关错误
	ErrItemNotFound       = errors.New("item not found")
	ErrItemAlreadyExists  = errors.New("item already exists")
	ErrItemExpired        = errors.New("item has expired")
	ErrItemNotPickable    = errors.New("item is not pickable")
	ErrItemOwnershipError = errors.New("item ownership error")

	// 传送门相关错误
	ErrPortalNotFound       = errors.New("portal not found")
	ErrPortalNotActive      = errors.New("portal is not active")
	ErrPortalLocked         = errors.New("portal is locked")
	ErrInsufficientLevel    = errors.New("insufficient level for portal")
	ErrMissingRequiredItems = errors.New("missing required items for portal")
	ErrInsufficientGold     = errors.New("insufficient gold for portal")

	// 刷新点相关错误
	ErrSpawnPointNotFound = errors.New("spawn point not found")
	ErrSpawnPointFull     = errors.New("spawn point is full")
	ErrSpawnPointInactive = errors.New("spawn point is inactive")
	ErrSpawnCooldown      = errors.New("spawn point is on cooldown")
	ErrInvalidSpawnType   = errors.New("invalid spawn type")

	// AOI相关错误
	ErrAOINotInitialized = errors.New("aoi manager not initialized")
	ErrInvalidAOIRadius  = errors.New("invalid aoi radius")
	ErrAOIEntityNotFound = errors.New("aoi entity not found")
	ErrAOIGridNotFound   = errors.New("aoi grid not found")

	// AI相关错误
	ErrAINotInitialized  = errors.New("ai behavior not initialized")
	ErrInvalidAIBehavior = errors.New("invalid ai behavior type")
	ErrAITargetNotFound  = errors.New("ai target not found")
	ErrAIPathNotFound    = errors.New("ai path not found")

	// 战斗相关错误
	ErrNotInCombat      = errors.New("not in combat")
	ErrAlreadyInCombat  = errors.New("already in combat")
	ErrCannotAttack     = errors.New("cannot attack target")
	ErrAttackOnCooldown = errors.New("attack is on cooldown")
	ErrInvalidTarget    = errors.New("invalid attack target")

	// 权限相关错误
	ErrInsufficientPermission = errors.New("insufficient permission")
	ErrAccessDenied           = errors.New("access denied")
	ErrNotAuthorized          = errors.New("not authorized")

	// 配置相关错误
	ErrInvalidSceneConfig  = errors.New("invalid scene configuration")
	ErrSceneConfigNotFound = errors.New("scene configuration not found")
	ErrInvalidEntityConfig = errors.New("invalid entity configuration")
	ErrConfigLoadFailed    = errors.New("failed to load configuration")
)
