package tcp

import (
	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/infrastructure/network"
)

// SceneHandler 场景TCP处理器
type SceneHandler struct {
	weatherService *services.WeatherService
	plantService   *services.PlantService
	logger         logger.Logger
}

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
	// TODO: 实现场景处理器注册功能
	h.logger.Info("Scene handlers registration not implemented")
	return nil
}
