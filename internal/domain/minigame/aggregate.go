package minigame

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// MinigameAggregate 小游戏聚合根
type MinigameAggregate struct {
	// 基础信息
	ID          string       `json:"id" bson:"_id"`
	GameID      string       `json:"game_id" bson:"game_id"`
	GameType    GameType     `json:"game_type" bson:"game_type"`
	Category    GameCategory `json:"category" bson:"category"`
	Name        string       `json:"name" bson:"name"`
	Description string       `json:"description" bson:"description"`

	// 游戏配置
	Config   *GameConfig   `json:"config" bson:"config"`
	Rules    *GameRules    `json:"rules" bson:"rules"`
	Settings *GameSettings `json:"settings" bson:"settings"`

	// 游戏状态
	Status    GameStatus    `json:"status" bson:"status"`
	Phase     GamePhase     `json:"phase" bson:"phase"`
	IsActive  bool          `json:"is_active" bson:"is_active"`
	StartTime *time.Time    `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty" bson:"end_time,omitempty"`
	Duration  time.Duration `json:"duration" bson:"duration"`

	// 玩家信息
	CreatorID  uint64        `json:"creator_id" bson:"creator_id"`
	Players    []*GamePlayer `json:"players" bson:"players"`
	MaxPlayers int32         `json:"max_players" bson:"max_players"`
	MinPlayers int32         `json:"min_players" bson:"min_players"`

	// 游戏数据
	GameData *GameData     `json:"game_data" bson:"game_data"`
	Scores   []*GameScore  `json:"scores" bson:"scores"`
	Results  []*GameResult `json:"results" bson:"results"`

	// 奖励信息
	Rewards    []*GameReward `json:"rewards" bson:"rewards"`
	RewardPool *RewardPool   `json:"reward_pool" bson:"reward_pool"`

	// 统计信息
	Statistics *GameStatistics `json:"statistics" bson:"statistics"`
	PlayCount  int64           `json:"play_count" bson:"play_count"`
	WinCount   int64           `json:"win_count" bson:"win_count"`
	LoseCount  int64           `json:"lose_count" bson:"lose_count"`

	// 版本和时间戳
	Version   int64     `json:"version" bson:"version"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	// 内部状态
	mutex  sync.RWMutex    `json:"-" bson:"-"`
	dirty  bool            `json:"-" bson:"-"`
	events []MinigameEvent `json:"-" bson:"-"`
}

// NewMinigameAggregate 创建新的小游戏聚合
func NewMinigameAggregate(gameID string, gameType GameType, category GameCategory, creatorID uint64) *MinigameAggregate {
	now := time.Now()
	return &MinigameAggregate{
		ID:         generateMinigameID(gameID, now),
		GameID:     gameID,
		GameType:   gameType,
		Category:   category,
		CreatorID:  creatorID,
		Status:     GameStatusWaiting,
		Phase:      GamePhaseWaiting,
		IsActive:   true,
		Players:    make([]*GamePlayer, 0),
		MaxPlayers: DefaultMaxPlayers,
		MinPlayers: DefaultMinPlayers,
		Scores:     make([]*GameScore, 0),
		Results:    make([]*GameResult, 0),
		Rewards:    make([]*GameReward, 0),
		Statistics: NewGameStatistics(),
		PlayCount:  0,
		WinCount:   0,
		LoseCount:  0,
		Version:    1,
		CreatedAt:  now,
		UpdatedAt:  now,
		dirty:      true,
		events:     make([]MinigameEvent, 0),
	}
}

// StartGame 开始游戏
func (mg *MinigameAggregate) StartGame() error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status != GameStatusWaiting {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusWaiting, "start_game")
	}

	// 检查玩家数量
	if int32(len(mg.Players)) < mg.MinPlayers {
		return NewMinigameInsufficientPlayersError(mg.GameID, int32(len(mg.Players)), mg.MinPlayers)
	}

	// 更新游戏状态
	now := time.Now()
	mg.Status = GameStatusRunning
	mg.Phase = GamePhaseRunning
	mg.StartTime = &now
	mg.PlayCount++

	// 初始化游戏数据
	mg.initializeGameData()

	// 发布游戏开始事件
	mg.addEvent(NewGameStartedEvent(mg.GameID, mg.CreatorID))

	// 标记为脏数据
	mg.markDirty()

	return nil
}

// EndGame 结束游戏
func (mg *MinigameAggregate) EndGame(reason GameEndReason) error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status != GameStatusRunning {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusRunning, "end_game")
	}

	// 更新游戏状态
	now := time.Now()
	mg.Status = GameStatusFinished
	mg.Phase = GamePhaseFinished
	mg.EndTime = &now

	// 计算游戏时长
	if mg.StartTime != nil {
		mg.Duration = now.Sub(*mg.StartTime)
	}

	// 计算最终结果
	mg.calculateFinalResults(reason)

	// 分发奖励
	mg.distributeRewards()

	// 更新统计信息
	mg.updateStatistics()

	// 发布游戏结束事件
	mg.addEvent(NewGameEndedEvent(mg.GameID, reason, mg.CreatorID))

	// 标记为脏数据
	mg.markDirty()

	return nil
}

// AddPlayer 添加玩家
func (mg *MinigameAggregate) AddPlayer(playerID uint64, playerName string) error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status != GameStatusWaiting {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusWaiting, "add_player")
	}

	// 检查玩家是否已存在
	if mg.hasPlayer(playerID) {
		return NewPlayerAlreadyInGameError(playerID, mg.GameID)
	}

	// 检查玩家数量限制
	if int32(len(mg.Players)) >= mg.MaxPlayers {
		return NewGameFullError(mg.GameID, mg.MaxPlayers, int32(len(mg.Players)))
	}

	// 添加玩家
	player := &GamePlayer{
		PlayerID: fmt.Sprintf("%d", playerID),
		Username: playerName,
		JoinTime: time.Now(),
		IsActive: true,
	}
	mg.Players = append(mg.Players, player)

	// 发布玩家加入事件
	mg.addEvent(NewPlayerJoinedGameEvent(mg.GameID, playerID, playerName))

	// 标记为脏数据
	mg.markDirty()

	return nil
}

// RemovePlayer 移除玩家
func (mg *MinigameAggregate) RemovePlayer(playerID uint64, reason PlayerLeaveReason) error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status == GameStatusFinished {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusRunning, "remove_player")
	}

	// 查找并移除玩家
	for i, player := range mg.Players {
		if player.PlayerID == fmt.Sprintf("%d", playerID) {
			// 更新玩家状态
			player.IsActive = false

			// 从活跃玩家列表中移除
			mg.Players = append(mg.Players[:i], mg.Players[i+1:]...)

			// 发布玩家离开事件
			mg.addEvent(NewPlayerLeftGameEvent(mg.GameID, playerID, player.Username, reason.String()))

			// 检查是否需要结束游戏
			if mg.Status == GameStatusRunning && int32(len(mg.Players)) < mg.MinPlayers {
				mg.EndGame(GameEndReasonInsufficientPlayers)
			}

			// 标记为脏数据
			mg.markDirty()

			return nil
		}
	}

	return NewPlayerNotInGameError(mg.GameID, playerID)
}

// UpdatePlayerScore 更新玩家分数
func (mg *MinigameAggregate) UpdatePlayerScore(playerID uint64, score int64, scoreType ScoreType) error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status != GameStatusRunning {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusRunning, "update_score")
	}

	// 检查玩家是否在游戏中
	player := mg.findPlayer(playerID)
	if player == nil {
		return NewPlayerNotInGameError(mg.GameID, playerID)
	}

	// 更新玩家分数
	oldScore := player.Score
	player.Score = score

	// 记录分数历史
	scoreRecord := &GameScore{
		ID:         generateScoreID(),
		GameID:     mg.GameID,
		PlayerID:   playerID,
		ScoreType:  scoreType,
		Value:      score,
		FinalScore: score,
	}
	mg.Scores = append(mg.Scores, scoreRecord)

	// 发布分数更新事件
	mg.addEvent(NewPlayerScoreUpdatedEvent(mg.GameID, playerID, oldScore, score, scoreType))

	// 检查游戏结束条件
	mg.checkGameEndConditions()

	// 标记为脏数据
	mg.markDirty()

	return nil
}

// UpdateGameData 更新游戏数据
func (mg *MinigameAggregate) UpdateGameData(key string, value interface{}) error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status != GameStatusRunning {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusRunning, "update_data")
	}

	// 初始化游戏数据
	if mg.GameData == nil {
		mg.GameData = NewGameData()
	}

	// 更新数据
	mg.GameData.SetData(key, value)

	// 发布数据更新事件
	mg.addEvent(NewGameDataUpdatedEvent(mg.GameID, key, value))

	// 标记为脏数据
	mg.markDirty()

	return nil
}

// PauseGame 暂停游戏
func (mg *MinigameAggregate) PauseGame() error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status != GameStatusRunning {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusRunning, "pause_game")
	}

	// 更新游戏状态
	mg.Status = GameStatusPaused
	mg.Phase = GamePhasePaused

	// 发布游戏暂停事件
	mg.addEvent(NewGamePausedEvent(mg.GameID, mg.CreatorID))

	// 标记为脏数据
	mg.markDirty()

	return nil
}

// ResumeGame 恢复游戏
func (mg *MinigameAggregate) ResumeGame() error {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	// 检查游戏状态
	if mg.Status != GameStatusPaused {
		return NewMinigameInvalidStateError(mg.GameID, mg.Status, GameStatusPaused, "resume_game")
	}

	// 更新游戏状态
	mg.Status = GameStatusRunning
	mg.Phase = GamePhaseRunning

	// 发布游戏恢复事件
	mg.addEvent(NewGameResumedEvent(mg.GameID, mg.CreatorID))

	// 标记为脏数据
	mg.markDirty()

	return nil
}

// GetPlayer 获取玩家信息
func (mg *MinigameAggregate) GetPlayer(playerID uint64) *GamePlayer {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	return mg.findPlayer(playerID)
}

// GetPlayers 获取所有玩家
func (mg *MinigameAggregate) GetPlayers() []*GamePlayer {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	players := make([]*GamePlayer, len(mg.Players))
	copy(players, mg.Players)
	return players
}

// GetGameData 获取游戏数据
func (mg *MinigameAggregate) GetGameData() *GameData {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	if mg.GameData == nil {
		return nil
	}
	return mg.GameData.Clone()
}

// GetScores 获取分数记录
func (mg *MinigameAggregate) GetScores() []*GameScore {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	scores := make([]*GameScore, len(mg.Scores))
	copy(scores, mg.Scores)
	return scores
}

// GetResults 获取游戏结果
func (mg *MinigameAggregate) GetResults() []*GameResult {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	results := make([]*GameResult, len(mg.Results))
	copy(results, mg.Results)
	return results
}

// GetStatistics 获取游戏统计
func (mg *MinigameAggregate) GetStatistics() *GameStatistics {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	if mg.Statistics == nil {
		return nil
	}
	return mg.Statistics.Clone()
}

// GetEvents 获取领域事件
func (mg *MinigameAggregate) GetEvents() []MinigameEvent {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	events := make([]MinigameEvent, len(mg.events))
	copy(events, mg.events)
	return events
}

// ClearEvents 清除领域事件
func (mg *MinigameAggregate) ClearEvents() {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	mg.events = make([]MinigameEvent, 0)
}

// IsDirty 检查是否有未保存的更改
func (mg *MinigameAggregate) IsDirty() bool {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	return mg.dirty
}

// MarkClean 标记为已保存
func (mg *MinigameAggregate) MarkClean() {
	mg.mutex.Lock()
	defer mg.mutex.Unlock()

	mg.dirty = false
}

// 私有方法

// initializeGameData 初始化游戏数据
func (mg *MinigameAggregate) initializeGameData() {
	if mg.GameData == nil {
		mg.GameData = NewGameData()
	}

	// 根据游戏类型初始化特定数据
	switch mg.GameType {
	case GameTypeSaveDog:
		mg.initializeSaveDogData()
	case GameTypePuzzle:
		mg.initializePuzzleData()
	case GameTypeRacing:
		mg.initializeRacingData()
	default:
		mg.initializeDefaultData()
	}
}

// initializeSaveDogData 初始化拯救小狗游戏数据
func (mg *MinigameAggregate) initializeSaveDogData() {
	mg.GameData.SetData("dog_position", map[string]int{"x": 0, "y": 0})
	mg.GameData.SetData("obstacles", make([]map[string]interface{}, 0))
	mg.GameData.SetData("rescue_progress", 0)
	mg.GameData.SetData("time_limit", 300) // 5分钟
}

// initializePuzzleData 初始化拼图游戏数据
func (mg *MinigameAggregate) initializePuzzleData() {
	mg.GameData.SetData("puzzle_pieces", make([]map[string]interface{}, 0))
	mg.GameData.SetData("completed_pieces", 0)
	mg.GameData.SetData("total_pieces", 100)
}

// initializeRacingData 初始化赛车游戏数据
func (mg *MinigameAggregate) initializeRacingData() {
	mg.GameData.SetData("track_length", 1000)
	mg.GameData.SetData("laps", 3)
	mg.GameData.SetData("player_positions", make(map[uint64]map[string]interface{}))
}

// initializeDefaultData 初始化默认游戏数据
func (mg *MinigameAggregate) initializeDefaultData() {
	mg.GameData.SetData("initialized", true)
	mg.GameData.SetData("start_time", time.Now().Unix())
}

// hasPlayer 检查玩家是否存在
func (mg *MinigameAggregate) hasPlayer(playerID uint64) bool {
	return mg.findPlayer(playerID) != nil
}

// findPlayer 查找玩家
func (mg *MinigameAggregate) findPlayer(playerID uint64) *GamePlayer {
	for _, player := range mg.Players {
		if player.PlayerID == fmt.Sprintf("%d", playerID) {
			return player
		}
	}
	return nil
}

// getPlayerIDs 获取所有玩家ID
func (mg *MinigameAggregate) getPlayerIDs() []uint64 {
	playerIDs := make([]uint64, len(mg.Players))
	for i, player := range mg.Players {
		// Convert string PlayerID to uint64
		if id, err := strconv.ParseUint(player.PlayerID, 10, 64); err == nil {
			playerIDs[i] = id
		}
	}
	return playerIDs
}

// checkGameEndConditions 检查游戏结束条件
func (mg *MinigameAggregate) checkGameEndConditions() {
	// 根据游戏类型检查结束条件
	switch mg.GameType {
	case GameTypeSaveDog:
		mg.checkSaveDogEndConditions()
	case GameTypePuzzle:
		mg.checkPuzzleEndConditions()
	case GameTypeRacing:
		mg.checkRacingEndConditions()
	default:
		mg.checkDefaultEndConditions()
	}
}

// checkSaveDogEndConditions 检查拯救小狗游戏结束条件
func (mg *MinigameAggregate) checkSaveDogEndConditions() {
	if mg.GameData == nil {
		return
	}

	// 检查救援进度
	if progressData, exists := mg.GameData.GetData("rescue_progress"); exists {
		if progress, ok := progressData.(int); ok && progress >= 100 {
			mg.EndGame(GameEndReasonCompleted)
			return
		}
	}

	// 检查时间限制
	if mg.StartTime != nil {
		if timeLimitData, exists := mg.GameData.GetData("time_limit"); exists {
			if timeLimit, ok := timeLimitData.(int); ok {
				if time.Since(*mg.StartTime).Seconds() >= float64(timeLimit) {
					mg.EndGame(GameEndReasonTimeout)
				}
			}
		}
	}
}

// checkPuzzleEndConditions 检查拼图游戏结束条件
func (mg *MinigameAggregate) checkPuzzleEndConditions() {
	if mg.GameData == nil {
		return
	}

	// 检查完成的拼图块数
	var completedPieces, totalPieces int
	if completedData, exists := mg.GameData.GetData("completed_pieces"); exists {
		if val, ok := completedData.(int); ok {
			completedPieces = val
		}
	}
	if totalData, exists := mg.GameData.GetData("total_pieces"); exists {
		if val, ok := totalData.(int); ok {
			totalPieces = val
		}
	}

	if completedPieces >= totalPieces {
		mg.EndGame(GameEndReasonCompleted)
	}
}

// checkRacingEndConditions 检查赛车游戏结束条件
func (mg *MinigameAggregate) checkRacingEndConditions() {
	// 检查是否有玩家完成所有圈数
	for _, player := range mg.Players {
		if player.Score >= 3 { // 假设3圈为完成条件
			mg.EndGame(GameEndReasonCompleted)
			return
		}
	}
}

// checkDefaultEndConditions 检查默认游戏结束条件
func (mg *MinigameAggregate) checkDefaultEndConditions() {
	// 检查是否达到目标分数
	for _, player := range mg.Players {
		if player.Score >= 1000 { // 默认目标分数
			mg.EndGame(GameEndReasonCompleted)
			return
		}
	}
}

// calculateFinalResults 计算最终结果
func (mg *MinigameAggregate) calculateFinalResults(reason GameEndReason) {
	mg.Results = make([]*GameResult, 0, len(mg.Players))

	// 按分数排序玩家
	players := make([]*GamePlayer, len(mg.Players))
	copy(players, mg.Players)

	// 简单的分数排序
	for i := 0; i < len(players)-1; i++ {
		for j := i + 1; j < len(players); j++ {
			if players[i].Score < players[j].Score {
				players[i], players[j] = players[j], players[i]
			}
		}
	}

	// 生成结果
	for i, player := range players {
		result := &GameResult{
			GameID:      mg.GameID,
			WinnerID:    player.PlayerID,
			WinnerName:  player.Username,
			FinalScore:  player.Score,
			CompletedAt: time.Now(),
		}
		mg.Results = append(mg.Results, result)

		// 更新胜负统计
		if i == 0 && reason == GameEndReasonCompleted {
			mg.WinCount++
		} else {
			mg.LoseCount++
		}
	}
}

// distributeRewards 分发奖励
func (mg *MinigameAggregate) distributeRewards() {
	if mg.RewardPool == nil {
		return
	}

	for _, result := range mg.Results {
		rewards := mg.RewardPool.CalculateRewards(result.Rank, result.Score, result.IsWinner)
		for _, reward := range rewards {
			reward.PlayerID = result.PlayerID
			reward.GameID = mg.GameID
			reward.Timestamp = time.Now()
			mg.Rewards = append(mg.Rewards, reward)
		}
	}
}

// updateStatistics 更新统计信息
func (mg *MinigameAggregate) updateStatistics() {
	if mg.Statistics == nil {
		mg.Statistics = NewGameStatistics()
	}

	mg.Statistics.TotalGames++
	mg.Statistics.TotalPlayers += int64(len(mg.Players))

	if len(mg.Results) > 0 {
		mg.Statistics.AverageScore = mg.calculateAverageScore()
		mg.Statistics.HighestScore = mg.getHighestScore()
		mg.Statistics.LowestScore = mg.getLowestScore()
	}

	mg.Statistics.AverageGameDuration = mg.Duration
	mg.Statistics.LastPlayedAt = time.Now()
	mg.Statistics.UpdatedAt = time.Now()
}

// calculateAverageScore 计算平均分数
func (mg *MinigameAggregate) calculateAverageScore() float64 {
	if len(mg.Results) == 0 {
		return 0
	}

	totalScore := int64(0)
	for _, result := range mg.Results {
		totalScore += result.Score
	}

	return float64(totalScore) / float64(len(mg.Results))
}

// getHighestScore 获取最高分数
func (mg *MinigameAggregate) getHighestScore() int64 {
	if len(mg.Results) == 0 {
		return 0
	}

	highest := mg.Results[0].Score
	for _, result := range mg.Results {
		if result.Score > highest {
			highest = result.Score
		}
	}

	return highest
}

// getLowestScore 获取最低分数
func (mg *MinigameAggregate) getLowestScore() int64 {
	if len(mg.Results) == 0 {
		return 0
	}

	lowest := mg.Results[0].Score
	for _, result := range mg.Results {
		if result.Score < lowest {
			lowest = result.Score
		}
	}

	return lowest
}

// addEvent 添加领域事件
func (mg *MinigameAggregate) addEvent(event MinigameEvent) {
	mg.events = append(mg.events, event)
}

// markDirty 标记为脏数据
func (mg *MinigameAggregate) markDirty() {
	mg.dirty = true
	mg.UpdatedAt = time.Now()
	mg.Version++
}

// 辅助函数

// generateMinigameID 生成小游戏ID
func generateMinigameID(gameID string, timestamp time.Time) string {
	return fmt.Sprintf("minigame_%s_%d", gameID, timestamp.Unix())
}

// 常量定义

const (
	// 默认配置
	DefaultMaxPlayers   = 10
	DefaultMinPlayers   = 1
	DefaultGameDuration = 10 * time.Minute

	// 游戏限制
	MaxGameDuration = 60 * time.Minute
	MinGameDuration = 1 * time.Minute
	MaxPlayersLimit = 100
	MinPlayersLimit = 1
)

// 验证方法

// Validate 验证小游戏聚合
func (mg *MinigameAggregate) Validate() error {
	if mg.GameID == "" {
		return NewMinigameValidationError("game_id", mg.GameID, "game_id cannot be empty", "required")
	}

	if !mg.GameType.IsValid() {
		return NewMinigameValidationError("game_type", mg.GameType, "invalid game type", "enum")
	}

	if !mg.Category.IsValid() {
		return NewMinigameValidationError("category", mg.Category, "invalid category", "enum")
	}

	if mg.MaxPlayers < mg.MinPlayers {
		return NewMinigameValidationError("max_players", mg.MaxPlayers, "max_players must be greater than or equal to min_players", "range")
	}

	if mg.MaxPlayers > MaxPlayersLimit {
		return NewMinigameValidationError("max_players", mg.MaxPlayers, fmt.Sprintf("max_players cannot exceed %d", MaxPlayersLimit), "max")
	}

	if mg.MinPlayers < MinPlayersLimit {
		return NewMinigameValidationError("min_players", mg.MinPlayers, fmt.Sprintf("min_players cannot be less than %d", MinPlayersLimit), "min")
	}

	return nil
}

// Clone 克隆小游戏聚合
func (mg *MinigameAggregate) Clone() *MinigameAggregate {
	mg.mutex.RLock()
	defer mg.mutex.RUnlock()

	clone := &MinigameAggregate{
		ID:          mg.ID,
		GameID:      mg.GameID,
		GameType:    mg.GameType,
		Category:    mg.Category,
		Name:        mg.Name,
		Description: mg.Description,
		Status:      mg.Status,
		Phase:       mg.Phase,
		IsActive:    mg.IsActive,
		Duration:    mg.Duration,
		CreatorID:   mg.CreatorID,
		MaxPlayers:  mg.MaxPlayers,
		MinPlayers:  mg.MinPlayers,
		PlayCount:   mg.PlayCount,
		WinCount:    mg.WinCount,
		LoseCount:   mg.LoseCount,
		Version:     mg.Version,
		CreatedAt:   mg.CreatedAt,
		UpdatedAt:   mg.UpdatedAt,
		dirty:       mg.dirty,
	}

	// 深拷贝时间
	if mg.StartTime != nil {
		startTime := *mg.StartTime
		clone.StartTime = &startTime
	}
	if mg.EndTime != nil {
		endTime := *mg.EndTime
		clone.EndTime = &endTime
	}

	// 深拷贝配置
	if mg.Config != nil {
		clone.Config = mg.Config.Clone()
	}
	if mg.Rules != nil {
		clone.Rules = mg.Rules.Clone()
	}
	if mg.Settings != nil {
		clone.Settings = mg.Settings.Clone()
	}

	// 深拷贝玩家
	clone.Players = make([]*GamePlayer, len(mg.Players))
	for i, player := range mg.Players {
		clone.Players[i] = player.Clone()
	}

	// 深拷贝游戏数据
	if mg.GameData != nil {
		clone.GameData = mg.GameData.Clone()
	}

	// 深拷贝分数记录
	clone.Scores = make([]*GameScore, len(mg.Scores))
	for i, score := range mg.Scores {
		scoreCopy := *score
		clone.Scores[i] = &scoreCopy
	}

	// 深拷贝结果
	clone.Results = make([]*GameResult, len(mg.Results))
	for i, result := range mg.Results {
		resultCopy := *result
		clone.Results[i] = &resultCopy
	}

	// 深拷贝奖励
	clone.Rewards = make([]*GameReward, len(mg.Rewards))
	for i, reward := range mg.Rewards {
		rewardCopy := *reward
		clone.Rewards[i] = &rewardCopy
	}

	// 深拷贝奖励池
	if mg.RewardPool != nil {
		clone.RewardPool = mg.RewardPool.Clone()
	}

	// 深拷贝统计信息
	if mg.Statistics != nil {
		clone.Statistics = mg.Statistics.Clone()
	}

	// 深拷贝事件
	clone.events = make([]MinigameEvent, len(mg.events))
	copy(clone.events, mg.events)

	return clone
}
