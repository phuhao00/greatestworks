// Package config 配置加载器
// Author: MMO Server Team
// Created: 2024

package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigValidator 配置验证器接口
type ConfigValidator interface {
	Validate(config *Config) error
}

// Load 加载配置
func (cl *ConfigLoader) Load() (*Config, error) {
	// 构建配置文件路径
	//todo complete  it
	configFile := cl.getConfigFilePath("")

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configFile)
	}

	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 环境变量替换
	configContent := cl.expandEnvVars(string(data))

	// 解析配置
	var config Config
	if err := yaml.Unmarshal([]byte(configContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 设置默认值
	cl.setDefaults(&config)

	// 验证配置
	if err := cl.validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// LoadFromFile 从指定文件加载配置
func (cl *ConfigLoader) LoadFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	configContent := cl.expandEnvVars(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(configContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cl.setDefaults(&config)

	if err := cl.validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// expandEnvVars 展开环境变量
func (cl *ConfigLoader) expandEnvVars(content string) string {
	return os.ExpandEnv(content)
}

// setDefaults 设置默认值
func (cl *ConfigLoader) setDefaults(config *Config) {
	setDefaults(config)
}

// validate 验证配置
func (cl *ConfigLoader) validate(config *Config) error {
	// 内置验证
	if err := validateConfig(config); err != nil {
		return err
	}

	// 自定义验证器
	for _, validator := range cl.validators {
		if err := validator.Validate(config); err != nil {
			return err
		}
	}

	return nil
}

// getEnvironment 获取当前环境
func getEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}
	if env == "" {
		env = "dev"
	}
	return strings.ToLower(env)
}

// setDefaults 设置配置默认值
func setDefaults(config *Config) {
	// 服务器默认配置
	if config.Server.Gateway.Port == 0 {
		config.Server.Gateway.Port = 8080
	}
	if config.Server.Gateway.Host == "" {
		config.Server.Gateway.Host = "0.0.0.0"
	}
	if config.Server.Gateway.MaxConnections == 0 {
		config.Server.Gateway.MaxConnections = 10000
	}
	if config.Server.Gateway.ReadTimeout == 0 {
		config.Server.Gateway.ReadTimeout = 30 * time.Second
	}
	if config.Server.Gateway.WriteTimeout == 0 {
		config.Server.Gateway.WriteTimeout = 30 * time.Second
	}
	if config.Server.Gateway.HeartbeatTime == 0 {
		config.Server.Gateway.HeartbeatTime = 30 * time.Second
	}
	if config.Server.Gateway.IdleTimeout == 0 {
		config.Server.Gateway.IdleTimeout = 300 * time.Second
	}

	// 场景服务器默认配置
	if config.Server.Scene.Port == 0 {
		config.Server.Scene.Port = 8081
	}
	if config.Server.Scene.Host == "" {
		config.Server.Scene.Host = "0.0.0.0"
	}
	if config.Server.Scene.MaxPlayers == 0 {
		config.Server.Scene.MaxPlayers = 1000
	}
	if config.Server.Scene.TickRate == 0 {
		config.Server.Scene.TickRate = 20
	}
	if config.Server.Scene.SyncInterval == 0 {
		config.Server.Scene.SyncInterval = 100 * time.Millisecond
	}
	if config.Server.Scene.ViewDistance == 0 {
		config.Server.Scene.ViewDistance = 100.0
	}

	// 战斗服务器默认配置
	if config.Server.Battle.Port == 0 {
		config.Server.Battle.Port = 8082
	}
	if config.Server.Battle.Host == "" {
		config.Server.Battle.Host = "0.0.0.0"
	}
	if config.Server.Battle.MaxBattles == 0 {
		config.Server.Battle.MaxBattles = 100
	}
	if config.Server.Battle.TickRate == 0 {
		config.Server.Battle.TickRate = 30
	}
	if config.Server.Battle.BattleTime == 0 {
		config.Server.Battle.BattleTime = 10 * time.Minute
	}
	if config.Server.Battle.MatchTimeout == 0 {
		config.Server.Battle.MatchTimeout = 30 * time.Second
	}

	// 活动服务器默认配置
	if config.Server.Activity.Port == 0 {
		config.Server.Activity.Port = 8083
	}
	if config.Server.Activity.Host == "" {
		config.Server.Activity.Host = "0.0.0.0"
	}
	if config.Server.Activity.MaxActivities == 0 {
		config.Server.Activity.MaxActivities = 50
	}
	if config.Server.Activity.UpdateInterval == 0 {
		config.Server.Activity.UpdateInterval = 1 * time.Minute
	}
	if config.Server.Activity.CacheTimeout == 0 {
		config.Server.Activity.CacheTimeout = 5 * time.Minute
	}

	// MongoDB默认配置
	if config.Database.MongoDB.URI == "" {
		config.Database.MongoDB.URI = "mongodb://localhost:27017"
	}
	if config.Database.MongoDB.Database == "" {
		config.Database.MongoDB.Database = "mmo_game"
	}
	if config.Database.MongoDB.MaxPoolSize == 0 {
		config.Database.MongoDB.MaxPoolSize = 100
	}
	if config.Database.MongoDB.MinPoolSize == 0 {
		config.Database.MongoDB.MinPoolSize = 5
	}
	if config.Database.MongoDB.MaxIdleTime == 0 {
		config.Database.MongoDB.MaxIdleTime = 10 * time.Minute
	}
	if config.Database.MongoDB.ConnectTimeout == 0 {
		config.Database.MongoDB.ConnectTimeout = 10 * time.Second
	}
	if config.Database.MongoDB.SocketTimeout == 0 {
		config.Database.MongoDB.SocketTimeout = 30 * time.Second
	}
	if config.Database.MongoDB.ReadPreference == "" {
		config.Database.MongoDB.ReadPreference = "primary"
	}

	// Redis默认配置
	if config.Database.Redis.Addr == "" {
		config.Database.Redis.Addr = "localhost:6379"
	}
	if config.Database.Redis.PoolSize == 0 {
		config.Database.Redis.PoolSize = 10
	}
	if config.Database.Redis.MinIdleConns == 0 {
		config.Database.Redis.MinIdleConns = 2
	}
	if config.Database.Redis.MaxIdleConns == 0 {
		config.Database.Redis.MaxIdleConns = 5
	}
	if config.Database.Redis.ConnMaxAge == 0 {
		config.Database.Redis.ConnMaxAge = 30 * time.Minute
	}
	if config.Database.Redis.DialTimeout == 0 {
		config.Database.Redis.DialTimeout = 5 * time.Second
	}
	if config.Database.Redis.ReadTimeout == 0 {
		config.Database.Redis.ReadTimeout = 3 * time.Second
	}
	if config.Database.Redis.WriteTimeout == 0 {
		config.Database.Redis.WriteTimeout = 3 * time.Second
	}

	// 缓存默认配置
	if config.Cache.DefaultTTL == 0 {
		config.Cache.DefaultTTL = 1 * time.Hour
	}
	if config.Cache.CleanupTime == 0 {
		config.Cache.CleanupTime = 10 * time.Minute
	}
	if config.Cache.MaxSize == 0 {
		config.Cache.MaxSize = 1024 * 1024 * 100 // 100MB
	}
	if config.Cache.EvictionPolicy == "" {
		config.Cache.EvictionPolicy = "lru"
	}

	// JWT默认配置
	if config.Security.JWT.TokenDuration == 0 {
		config.Security.JWT.TokenDuration = 24 * time.Hour
	}
	if config.Security.JWT.RefreshTime == 0 {
		config.Security.JWT.RefreshTime = 7 * 24 * time.Hour
	}
	if config.Security.JWT.Issuer == "" {
		config.Security.JWT.Issuer = "mmo-server"
	}

	// 加密默认配置
	if config.Security.Encryption.Algorithm == "" {
		config.Security.Encryption.Algorithm = "AES-256-GCM"
	}
	if config.Security.Encryption.KeySize == 0 {
		config.Security.Encryption.KeySize = 32
	}

	// 限流默认配置
	if config.Security.RateLimit.RPS == 0 {
		config.Security.RateLimit.RPS = 100
	}
	if config.Security.RateLimit.Burst == 0 {
		config.Security.RateLimit.Burst = 200
	}
	if config.Security.RateLimit.WindowSize == 0 {
		config.Security.RateLimit.WindowSize = 1 * time.Minute
	}

	// 日志默认配置
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}
	if config.Logging.Dir == "" {
		config.Logging.Dir = "./logs"
	}
	if config.Logging.MaxSize == 0 {
		config.Logging.MaxSize = 100 // MB
	}
	if config.Logging.MaxBackups == 0 {
		config.Logging.MaxBackups = 7
	}
	if config.Logging.MaxAge == 0 {
		config.Logging.MaxAge = 30 // days
	}

	// 游戏默认配置
	if config.Game.MaxLevel == 0 {
		config.Game.MaxLevel = 100
	}
	if config.Game.ExpMultiplier == 0 {
		config.Game.ExpMultiplier = 1.0
	}
	if config.Game.GoldMultiplier == 0 {
		config.Game.GoldMultiplier = 1.0
	}
	if config.Game.DropRate == 0 {
		config.Game.DropRate = 0.1
	}

	// 网络默认配置
	if config.Network.Protocol == "" {
		config.Network.Protocol = "tcp"
	}
	if config.Network.BufferSize == 0 {
		config.Network.BufferSize = 4096
	}
	if config.Network.MaxPacketSize == 0 {
		config.Network.MaxPacketSize = 65536
	}
	if config.Network.CompressionType == "" {
		config.Network.CompressionType = "gzip"
	}
	if config.Network.EncryptionType == "" {
		config.Network.EncryptionType = "aes"
	}

	// NSQ默认配置
	if config.Messaging.NSQ.NSQDAddress == "" {
		config.Messaging.NSQ.NSQDAddress = "localhost:4150"
	}
	if len(config.Messaging.NSQ.LookupdHTTP) == 0 {
		config.Messaging.NSQ.LookupdHTTP = []string{"localhost:4161"}
	}
	if config.Messaging.NSQ.MaxInFlight == 0 {
		config.Messaging.NSQ.MaxInFlight = 200
	}

	// RabbitMQ默认配置
	if config.Messaging.RabbitMQ.URL == "" {
		config.Messaging.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	}
	if config.Messaging.RabbitMQ.Exchange == "" {
		config.Messaging.RabbitMQ.Exchange = "mmo_exchange"
	}
	if config.Messaging.RabbitMQ.Queue == "" {
		config.Messaging.RabbitMQ.Queue = "mmo_queue"
	}

	// 监控默认配置
	if config.Monitoring.Port == 0 {
		config.Monitoring.Port = 9090
	}
	if config.Monitoring.Path == "" {
		config.Monitoring.Path = "/metrics"
	}
	if config.Monitoring.Prometheus.Namespace == "" {
		config.Monitoring.Prometheus.Namespace = "mmo"
	}
	if config.Monitoring.Prometheus.Subsystem == "" {
		config.Monitoring.Prometheus.Subsystem = "server"
	}

	// Excel默认配置
	if config.Excel.Path == "" {
		config.Excel.Path = "./excel"
	}
}
