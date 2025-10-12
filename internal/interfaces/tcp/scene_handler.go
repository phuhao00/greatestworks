package tcp

import (
	"greatestworks/internal/application/services"
	"greatestworks/internal/infrastructure/logging"
)

// SceneHandler 场景TCP处理器
type SceneHandler struct {
	weatherService *services.WeatherService
	plantService   *services.PlantService
	logger         logging.Logger
}

// NewSceneHandler 创建场景处理器
func NewSceneHandler(weatherService *services.WeatherService, plantService *services.PlantService, logger logging.Logger) *SceneHandler {
	return &SceneHandler{
		weatherService: weatherService,
		plantService:   plantService,
		logger:         logger,
	}
}

// RegisterHandlers 注册处理器
func (h *SceneHandler) RegisterHandlers(server interface{}) error {
	// TODO: 实现场景处理器注册功能
	h.logger.Info("Scene handlers registration not implemented", map[string]interface{}{})
	return nil
}
