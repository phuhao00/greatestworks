// Package config 环境管理器
// Author: MMO Server Team
// Created: 2024

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Environment 环境类型
type Environment string

const (
	// EnvDevelopment 开发环境
	EnvDevelopment Environment = "dev"
	// EnvTesting 测试环境
	EnvTesting Environment = "test"
	// EnvStaging 预发布环境
	EnvStaging Environment = "staging"
	// EnvProduction 生产环境
	EnvProduction Environment = "prod"
	// EnvLocal 本地环境
	EnvLocal Environment = "local"
)

// EnvManager 环境管理器
type EnvManager struct {
	currentEnv Environment
	envVars    map[string]string
	mutex      sync.RWMutex
}

// NewEnvManager 创建环境管理器
func NewEnvManager() *EnvManager {
	return &EnvManager{
		currentEnv: detectEnvironment(),
		envVars:    make(map[string]string),
	}
}

// GetCurrentEnvironment 获取当前环境
func (em *EnvManager) GetCurrentEnvironment() Environment {
	em.mutex.RLock()
	defer em.mutex.RUnlock()
	return em.currentEnv
}

// SetEnvironment 设置环境
func (em *EnvManager) SetEnvironment(env Environment) {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	em.currentEnv = env
}

// IsProduction 是否为生产环境
func (em *EnvManager) IsProduction() bool {
	return em.GetCurrentEnvironment() == EnvProduction
}

// IsDevelopment 是否为开发环境
func (em *EnvManager) IsDevelopment() bool {
	env := em.GetCurrentEnvironment()
	return env == EnvDevelopment || env == EnvLocal
}

// IsTesting 是否为测试环境
func (em *EnvManager) IsTesting() bool {
	return em.GetCurrentEnvironment() == EnvTesting
}

// IsStaging 是否为预发布环境
func (em *EnvManager) IsStaging() bool {
	return em.GetCurrentEnvironment() == EnvStaging
}

// LoadEnvFile 加载环境变量文件
func (em *EnvManager) LoadEnvFile(envDir string) error {
	envFile := em.getEnvFilePath(envDir)
	if envFile == "" {
		return nil // 没有找到环境文件，不是错误
	}

	return em.loadEnvFromFile(envFile)
}

// LoadEnvFromFile 从指定文件加载环境变量
func (em *EnvManager) LoadEnvFromFile(filePath string) error {
	return em.loadEnvFromFile(filePath)
}

// GetEnvVar 获取环境变量
func (em *EnvManager) GetEnvVar(key string) string {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	// 优先从缓存获取
	if value, exists := em.envVars[key]; exists {
		return value
	}

	// 从系统环境变量获取
	return os.Getenv(key)
}

// SetEnvVar 设置环境变量
func (em *EnvManager) SetEnvVar(key, value string) {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	em.envVars[key] = value
	os.Setenv(key, value)
}

// GetConfigFileName 获取配置文件名
func (em *EnvManager) GetConfigFileName() string {
	return fmt.Sprintf("config.%s.yaml", em.GetCurrentEnvironment())
}

// GetConfigFilePath 获取配置文件路径
func (em *EnvManager) GetConfigFilePath(configDir string) string {
	// 优先级：环境变量指定 > 环境特定文件 > 默认文件
	if configFile := em.GetEnvVar("CONFIG_FILE"); configFile != "" {
		if filepath.IsAbs(configFile) {
			return configFile
		}
		return filepath.Join(configDir, configFile)
	}

	// 环境特定配置文件
	envConfigFile := em.GetConfigFileName()
	envConfigPath := filepath.Join(configDir, envConfigFile)
	if _, err := os.Stat(envConfigPath); err == nil {
		return envConfigPath
	}

	// 默认配置文件
	return filepath.Join(configDir, "config.yaml")
}

// ValidateEnvironment 验证环境配置
func (em *EnvManager) ValidateEnvironment() error {
	env := em.GetCurrentEnvironment()

	// 验证环境类型
	validEnvs := []Environment{EnvDevelopment, EnvTesting, EnvStaging, EnvProduction, EnvLocal}
	valid := false
	for _, validEnv := range validEnvs {
		if env == validEnv {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid environment: %s", env)
	}

	// 生产环境特殊验证
	if em.IsProduction() {
		requiredVars := []string{
			"JWT_SECRET_KEY",
			"MONGODB_URI",
			"REDIS_ADDR",
		}

		for _, varName := range requiredVars {
			if em.GetEnvVar(varName) == "" {
				return fmt.Errorf("required environment variable %s is not set for production", varName)
			}
		}
	}

	return nil
}

// GetEnvironmentSpecificConfig 获取环境特定配置
func (em *EnvManager) GetEnvironmentSpecificConfig() map[string]interface{} {
	config := make(map[string]interface{})

	switch em.GetCurrentEnvironment() {
	case EnvDevelopment, EnvLocal:
		config["debug"] = true
		config["log_level"] = "debug"
		config["enable_pprof"] = true
		config["enable_metrics"] = true

	case EnvTesting:
		config["debug"] = true
		config["log_level"] = "info"
		config["enable_pprof"] = false
		config["enable_metrics"] = true
		config["test_mode"] = true

	case EnvStaging:
		config["debug"] = false
		config["log_level"] = "info"
		config["enable_pprof"] = false
		config["enable_metrics"] = true
		config["rate_limit_enabled"] = true

	case EnvProduction:
		config["debug"] = false
		config["log_level"] = "warn"
		config["enable_pprof"] = false
		config["enable_metrics"] = true
		config["rate_limit_enabled"] = true
		config["tls_enabled"] = true
		config["security_enhanced"] = true
	}

	return config
}

// ApplyEnvironmentOverrides 应用环境覆盖配置
func (em *EnvManager) ApplyEnvironmentOverrides(config *Config) {
	envConfig := em.GetEnvironmentSpecificConfig()

	// 应用日志级别覆盖
	if logLevel, exists := envConfig["log_level"]; exists {
		if config.Logging.Level == "" {
			config.Logging.Level = logLevel.(string)
		}
	}

	// 应用安全配置覆盖
	if tlsEnabled, exists := envConfig["tls_enabled"]; exists && tlsEnabled.(bool) {
		config.Server.Gateway.TLSEnabled = true
	}

	if rateLimitEnabled, exists := envConfig["rate_limit_enabled"]; exists && rateLimitEnabled.(bool) {
		config.Security.RateLimit.Enabled = true
	}

	// 应用监控配置覆盖
	if enableMetrics, exists := envConfig["enable_metrics"]; exists {
		config.Monitoring.Enabled = enableMetrics.(bool)
	}

	// 应用环境变量覆盖
	em.applyEnvVarOverrides(config)
}

// applyEnvVarOverrides 应用环境变量覆盖
func (em *EnvManager) applyEnvVarOverrides(config *Config) {
	// 服务器配置覆盖
	if port := em.GetEnvVar("GATEWAY_PORT"); port != "" {
		if p := parseInt(port); p > 0 {
			config.Server.Gateway.Port = p
		}
	}

	if host := em.GetEnvVar("GATEWAY_HOST"); host != "" {
		config.Server.Gateway.Host = host
	}

	// 数据库配置覆盖
	if mongoURI := em.GetEnvVar("MONGODB_URI"); mongoURI != "" {
		config.Database.MongoDB.URI = mongoURI
	}

	if mongoDatabase := em.GetEnvVar("MONGODB_DATABASE"); mongoDatabase != "" {
		config.Database.MongoDB.Database = mongoDatabase
	}

	if redisAddr := em.GetEnvVar("REDIS_ADDR"); redisAddr != "" {
		config.Database.Redis.Addr = redisAddr
	}

	if redisPassword := em.GetEnvVar("REDIS_PASSWORD"); redisPassword != "" {
		config.Database.Redis.Password = redisPassword
	}

	// 安全配置覆盖
	if jwtSecret := em.GetEnvVar("JWT_SECRET_KEY"); jwtSecret != "" {
		config.Security.JWT.SecretKey = jwtSecret
	}

	if encryptionSalt := em.GetEnvVar("ENCRYPTION_SALT"); encryptionSalt != "" {
		config.Security.Encryption.Salt = encryptionSalt
	}

	// 日志配置覆盖
	if logLevel := em.GetEnvVar("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}

	if logDir := em.GetEnvVar("LOG_DIR"); logDir != "" {
		config.Logging.Dir = logDir
	}

	// 消息队列配置覆盖
	if nsqdAddr := em.GetEnvVar("NSQD_ADDRESS"); nsqdAddr != "" {
		config.Messaging.NSQ.NSQDAddress = nsqdAddr
	}

	if rabbitmqURL := em.GetEnvVar("RABBITMQ_URL"); rabbitmqURL != "" {
		config.Messaging.RabbitMQ.URL = rabbitmqURL
	}
}

// detectEnvironment 检测当前环境
func detectEnvironment() Environment {
	// 优先级：APP_ENV > ENVIRONMENT > GO_ENV > 默认dev
	envVars := []string{"APP_ENV", "ENVIRONMENT", "GO_ENV"}

	for _, envVar := range envVars {
		if env := os.Getenv(envVar); env != "" {
			return Environment(strings.ToLower(env))
		}
	}

	return EnvDevelopment
}

// getEnvFilePath 获取环境文件路径
func (em *EnvManager) getEnvFilePath(envDir string) string {
	// 优先级：.env.{environment} > .env
	envFiles := []string{
		fmt.Sprintf(".env.%s", em.GetCurrentEnvironment()),
		".env",
	}

	for _, envFile := range envFiles {
		filePath := filepath.Join(envDir, envFile)
		if _, err := os.Stat(filePath); err == nil {
			return filePath
		}
	}

	return ""
}

// loadEnvFromFile 从文件加载环境变量
func (em *EnvManager) loadEnvFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read env file %s: %w", filePath, err)
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid env file format at line %d: %s", i+1, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除引号
		if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		em.SetEnvVar(key, value)
	}

	return nil
}

// parseInt 解析整数
func parseInt(s string) int {
	var result int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		result = result*10 + int(r-'0')
	}
	return result
}

// 全局环境管理器实例
var (
	globalEnvManager *EnvManager
	envManagerOnce   sync.Once
)

// GetEnvManager 获取全局环境管理器
func GetEnvManager() *EnvManager {
	envManagerOnce.Do(func() {
		globalEnvManager = NewEnvManager()
	})
	return globalEnvManager
}

// InitEnvManager 初始化环境管理器
func InitEnvManager(envDir string) error {
	manager := GetEnvManager()
	return manager.LoadEnvFile(envDir)
}

// GetCurrentEnv 获取当前环境（便捷函数）
func GetCurrentEnv() Environment {
	return GetEnvManager().GetCurrentEnvironment()
}

// IsProduction 是否为生产环境（便捷函数）
func IsProduction() bool {
	return GetEnvManager().IsProduction()
}

// IsDevelopment 是否为开发环境（便捷函数）
func IsDevelopment() bool {
	return GetEnvManager().IsDevelopment()
}

// GetEnv 获取环境变量（便捷函数）
func GetEnv(key string) string {
	return GetEnvManager().GetEnvVar(key)
}

// SetEnv 设置环境变量（便捷函数）
func SetEnv(key, value string) {
	GetEnvManager().SetEnvVar(key, value)
}
