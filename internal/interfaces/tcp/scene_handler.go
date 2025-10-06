package tcp

import (
	"context"
	"encoding/json"
	"fmt"

	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/infrastructure/network"
	// "github.com/netcore-go/netcore" // TODO: 实现netcore-go集成
)

// SceneHandler 场景TCP处理器
type SceneHandler struct {
	weatherService *services.WeatherService
	plantService   *services.PlantService
	logger         logger.Logger
}

// SceneRequest 场景请求
type SceneRequest struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// SceneResponse 场景响应
type SceneResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 消息类型常量
const (
	// 天气系统
	MsgTypeWeatherInfo     uint32 = 2001
	MsgTypeWeatherForecast uint32 = 2002
	MsgTypeWeatherUpdate   uint32 = 2003

	// 种植系统
	MsgTypeFarmInfo       uint32 = 2101
	MsgTypePlantSeed      uint32 = 2102
	MsgTypeHarvestCrop    uint32 = 2103
	MsgTypeWaterPlant     uint32 = 2104
	MsgTypeFertilizePlant uint32 = 2105
	MsgTypeCropStatus     uint32 = 2106
	MsgTypeFarmUpgrade    uint32 = 2107
)

// NewSceneHandler 创建场景处理器
func NewSceneHandler(weatherService *services.WeatherService, plantService *services.PlantService, logger logger.Logger) *SceneHandler {
	return &SceneHandler{
		weatherService: weatherService,
		plantService:   plantService,
		logger:         logger,
	}
}

// RegisterHandlers 注册处理器
func (h *SceneHandler) RegisterHandlers(server network.Server) error {
	// 注册天气相关处理器
	if err := server.RegisterHandler(&WeatherInfoHandler{h}); err != nil {
		return fmt.Errorf("failed to register weather info handler: %w", err)
	}

	if err := server.RegisterHandler(&WeatherForecastHandler{h}); err != nil {
		return fmt.Errorf("failed to register weather forecast handler: %w", err)
	}

	if err := server.RegisterHandler(&WeatherUpdateHandler{h}); err != nil {
		return fmt.Errorf("failed to register weather update handler: %w", err)
	}

	// 注册种植相关处理器
	if err := server.RegisterHandler(&FarmInfoHandler{h}); err != nil {
		return fmt.Errorf("failed to register farm info handler: %w", err)
	}

	if err := server.RegisterHandler(&PlantSeedHandler{h}); err != nil {
		return fmt.Errorf("failed to register plant seed handler: %w", err)
	}

	if err := server.RegisterHandler(&HarvestCropHandler{h}); err != nil {
		return fmt.Errorf("failed to register harvest crop handler: %w", err)
	}

	if err := server.RegisterHandler(&WaterPlantHandler{h}); err != nil {
		return fmt.Errorf("failed to register water plant handler: %w", err)
	}

	if err := server.RegisterHandler(&FertilizePlantHandler{h}); err != nil {
		return fmt.Errorf("failed to register fertilize plant handler: %w", err)
	}

	if err := server.RegisterHandler(&CropStatusHandler{h}); err != nil {
		return fmt.Errorf("failed to register crop status handler: %w", err)
	}

	if err := server.RegisterHandler(&FarmUpgradeHandler{h}); err != nil {
		return fmt.Errorf("failed to register farm upgrade handler: %w", err)
	}

	h.logger.Info("Scene handlers registered successfully")
	return nil
}

// 天气信息处理器
type WeatherInfoHandler struct {
	*SceneHandler
}

func (h *WeatherInfoHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal weather info request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取区域ID
	regionID, ok := req.Data["region_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing region_id")
	}

	// 调用服务层获取天气信息
	weather, err := h.weatherService.GetCurrentWeather(ctx, regionID)
	if err != nil {
		h.logger.Error("Failed to get weather info", "error", err, "region_id", regionID)
		return h.sendErrorResponse(conn, "Failed to get weather info: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Weather info retrieved successfully",
		Data:    weather,
	}

	return h.sendResponse(conn, MsgTypeWeatherInfo, response)
}

func (h *WeatherInfoHandler) GetMessageType() uint32 {
	return MsgTypeWeatherInfo
}

func (h *WeatherInfoHandler) GetHandlerName() string {
	return "WeatherInfoHandler"
}

// 天气预报处理器
type WeatherForecastHandler struct {
	*SceneHandler
}

func (h *WeatherForecastHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal weather forecast request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取区域ID和天数
	regionID, ok := req.Data["region_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing region_id")
	}

	days := 7 // 默认7天
	if d, ok := req.Data["days"].(float64); ok {
		days = int(d)
	}

	// 调用服务层获取天气预报
	forecast, err := h.weatherService.GetWeatherForecast(ctx, regionID, days)
	if err != nil {
		h.logger.Error("Failed to get weather forecast", "error", err, "region_id", regionID, "days", days)
		return h.sendErrorResponse(conn, "Failed to get weather forecast: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Weather forecast retrieved successfully",
		Data:    forecast,
	}

	return h.sendResponse(conn, MsgTypeWeatherForecast, response)
}

func (h *WeatherForecastHandler) GetMessageType() uint32 {
	return MsgTypeWeatherForecast
}

func (h *WeatherForecastHandler) GetHandlerName() string {
	return "WeatherForecastHandler"
}

// 天气更新处理器
type WeatherUpdateHandler struct {
	*SceneHandler
}

func (h *WeatherUpdateHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal weather update request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取区域ID
	regionID, ok := req.Data["region_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing region_id")
	}

	// 调用服务层更新天气
	if err := h.weatherService.UpdateWeather(ctx, regionID); err != nil {
		h.logger.Error("Failed to update weather", "error", err, "region_id", regionID)
		return h.sendErrorResponse(conn, "Failed to update weather: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Weather updated successfully",
	}

	h.logger.Info("Weather updated successfully", "region_id", regionID)
	return h.sendResponse(conn, MsgTypeWeatherUpdate, response)
}

func (h *WeatherUpdateHandler) GetMessageType() uint32 {
	return MsgTypeWeatherUpdate
}

func (h *WeatherUpdateHandler) GetHandlerName() string {
	return "WeatherUpdateHandler"
}

// 农场信息处理器
type FarmInfoHandler struct {
	*SceneHandler
}

func (h *FarmInfoHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal farm info request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 调用服务层获取农场信息
	farmInfo, err := h.plantService.GetFarmInfo(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to get farm info", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Failed to get farm info: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Farm info retrieved successfully",
		Data:    farmInfo,
	}

	return h.sendResponse(conn, MsgTypeFarmInfo, response)
}

func (h *FarmInfoHandler) GetMessageType() uint32 {
	return MsgTypeFarmInfo
}

func (h *FarmInfoHandler) GetHandlerName() string {
	return "FarmInfoHandler"
}

// 种植种子处理器
type PlantSeedHandler struct {
	*SceneHandler
}

func (h *PlantSeedHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal plant seed request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	seedID, ok := req.Data["seed_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing seed_id")
	}

	plotID, ok := req.Data["plot_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing plot_id")
	}

	// 调用服务层种植种子
	crop, err := h.plantService.PlantSeed(ctx, playerID, seedID, plotID)
	if err != nil {
		h.logger.Error("Failed to plant seed", "error", err, "player_id", playerID, "seed_id", seedID, "plot_id", plotID)
		return h.sendErrorResponse(conn, "Failed to plant seed: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Seed planted successfully",
		Data:    crop,
	}

	h.logger.Info("Seed planted successfully", "player_id", playerID, "seed_id", seedID, "plot_id", plotID)
	return h.sendResponse(conn, MsgTypePlantSeed, response)
}

func (h *PlantSeedHandler) GetMessageType() uint32 {
	return MsgTypePlantSeed
}

func (h *PlantSeedHandler) GetHandlerName() string {
	return "PlantSeedHandler"
}

// 收获作物处理器
type HarvestCropHandler struct {
	*SceneHandler
}

func (h *HarvestCropHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal harvest crop request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	cropID, ok := req.Data["crop_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing crop_id")
	}

	// 调用服务层收获作物
	rewards, err := h.plantService.HarvestCrop(ctx, playerID, cropID)
	if err != nil {
		h.logger.Error("Failed to harvest crop", "error", err, "player_id", playerID, "crop_id", cropID)
		return h.sendErrorResponse(conn, "Failed to harvest crop: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Crop harvested successfully",
		Data:    rewards,
	}

	h.logger.Info("Crop harvested successfully", "player_id", playerID, "crop_id", cropID)
	return h.sendResponse(conn, MsgTypeHarvestCrop, response)
}

func (h *HarvestCropHandler) GetMessageType() uint32 {
	return MsgTypeHarvestCrop
}

func (h *HarvestCropHandler) GetHandlerName() string {
	return "HarvestCropHandler"
}

// 浇水处理器
type WaterPlantHandler struct {
	*SceneHandler
}

func (h *WaterPlantHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal water plant request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	cropID, ok := req.Data["crop_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing crop_id")
	}

	// 调用服务层浇水
	if err := h.plantService.WaterPlant(ctx, playerID, cropID); err != nil {
		h.logger.Error("Failed to water plant", "error", err, "player_id", playerID, "crop_id", cropID)
		return h.sendErrorResponse(conn, "Failed to water plant: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Plant watered successfully",
	}

	h.logger.Info("Plant watered successfully", "player_id", playerID, "crop_id", cropID)
	return h.sendResponse(conn, MsgTypeWaterPlant, response)
}

func (h *WaterPlantHandler) GetMessageType() uint32 {
	return MsgTypeWaterPlant
}

func (h *WaterPlantHandler) GetHandlerName() string {
	return "WaterPlantHandler"
}

// 施肥处理器
type FertilizePlantHandler struct {
	*SceneHandler
}

func (h *FertilizePlantHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal fertilize plant request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	cropID, ok := req.Data["crop_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing crop_id")
	}

	fertilizerID, ok := req.Data["fertilizer_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing fertilizer_id")
	}

	// 调用服务层施肥
	if err := h.plantService.FertilizePlant(ctx, playerID, cropID, fertilizerID); err != nil {
		h.logger.Error("Failed to fertilize plant", "error", err, "player_id", playerID, "crop_id", cropID, "fertilizer_id", fertilizerID)
		return h.sendErrorResponse(conn, "Failed to fertilize plant: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Plant fertilized successfully",
	}

	h.logger.Info("Plant fertilized successfully", "player_id", playerID, "crop_id", cropID, "fertilizer_id", fertilizerID)
	return h.sendResponse(conn, MsgTypeFertilizePlant, response)
}

func (h *FertilizePlantHandler) GetMessageType() uint32 {
	return MsgTypeFertilizePlant
}

func (h *FertilizePlantHandler) GetHandlerName() string {
	return "FertilizePlantHandler"
}

// 作物状态处理器
type CropStatusHandler struct {
	*SceneHandler
}

func (h *CropStatusHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal crop status request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	cropID, ok := req.Data["crop_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing crop_id")
	}

	// 调用服务层获取作物状态
	status, err := h.plantService.GetCropStatus(ctx, playerID, cropID)
	if err != nil {
		h.logger.Error("Failed to get crop status", "error", err, "player_id", playerID, "crop_id", cropID)
		return h.sendErrorResponse(conn, "Failed to get crop status: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Crop status retrieved successfully",
		Data:    status,
	}

	return h.sendResponse(conn, MsgTypeCropStatus, response)
}

func (h *CropStatusHandler) GetMessageType() uint32 {
	return MsgTypeCropStatus
}

func (h *CropStatusHandler) GetHandlerName() string {
	return "CropStatusHandler"
}

// 农场升级处理器
type FarmUpgradeHandler struct {
	*SceneHandler
}

func (h *FarmUpgradeHandler) Handle(ctx context.Context, conn network.Connection, packet network.Packet) error {
	var req SceneRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal farm upgrade request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	upgradeType, ok := req.Data["upgrade_type"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing upgrade_type")
	}

	// 调用服务层升级农场
	result, err := h.plantService.UpgradeFarm(ctx, playerID, upgradeType)
	if err != nil {
		h.logger.Error("Failed to upgrade farm", "error", err, "player_id", playerID, "upgrade_type", upgradeType)
		return h.sendErrorResponse(conn, "Failed to upgrade farm: "+err.Error())
	}

	// 发送成功响应
	response := SceneResponse{
		Success: true,
		Message: "Farm upgraded successfully",
		Data:    result,
	}

	h.logger.Info("Farm upgraded successfully", "player_id", playerID, "upgrade_type", upgradeType)
	return h.sendResponse(conn, MsgTypeFarmUpgrade, response)
}

func (h *FarmUpgradeHandler) GetMessageType() uint32 {
	return MsgTypeFarmUpgrade
}

func (h *FarmUpgradeHandler) GetHandlerName() string {
	return "FarmUpgradeHandler"
}

// 辅助方法

// sendResponse 发送响应
func (h *SceneHandler) sendResponse(conn network.Connection, msgType uint32, response SceneResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal response", "error", err)
		return err
	}

	packet := network.NewPacket(msgType, data)
	return conn.Send(packet)
}

// sendErrorResponse 发送错误响应
func (h *SceneHandler) sendErrorResponse(conn network.Connection, errorMsg string) error {
	response := SceneResponse{
		Success: false,
		Message: "Request failed",
		Error:   errorMsg,
	}

	data, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal error response", "error", err)
		return err
	}

	// 使用通用错误消息类型
	packet := network.NewPacket(9999, data)
	return conn.Send(packet)
}
